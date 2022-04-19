package ipam

import (
	"context"
	"fmt"
	"github.com/pkg/errors"
	corev1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"net"
	infrav1 "sigs.k8s.io/cluster-api/test/infrastructure/docker/api/v1alpha4"
	"sigs.k8s.io/cluster-api/test/infrastructure/docker/docker"
	"sigs.k8s.io/cluster-api/test/infrastructure/docker/utils"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"strings"
	"sync"
	"time"
)

var dockerNetworkCidr string

type IpamProvisioner struct {
	client.Client
}

const (
	configMapName = "ipam-config"
	CidrKey = "CIDR"
	LastClaimedKey   = "LAST_CLAIMED_IP"
	Available        = "AVAILABLE"

)

var p *IpamProvisioner
var ipamLatch sync.Mutex

func InitProvisioner(client client.Client) *IpamProvisioner {
	p = &IpamProvisioner{
		Client: client,
	}
	return p
}

/*
ClaimIP
* Do IP discovery to reclaim all released ip (mark released ips as AVAILABLE in ipam-config config map)

Before claiming an ip
1. First check if ip is already claimed for the current docker machine
2. If ip isn't already claimed then look for a released ip in the pool (if docker machine isn't present now)
3. if there is no available ip from claimed ip pool, release next ip. Next IP from pool will be generated.
*/

// ssh -i ~/work/spectro2020.pem ubuntu@10.10.149.168

const (
	defaultDockerNetwork = "kind"
)

func ClaimIP(namespace, machineName string) (string, error) {
	ipamLatch.Lock()
	defer ipamLatch.Unlock()
	_ = p.ipDiscovery(namespace)

	cidr, err := getCidrRange()
	if err != nil {
		return "", nil
	}

	if existingClaim, ip, err := p.findIpFromClaimedIpPool(namespace, machineName); err != nil {
		return "", err
	} else if existingClaim && isIpAvailableInNetwork(ip){
		return ip, utils.Retry(func() error { return p.cacheIPReclaim(namespace, cidr, ip, machineName) }, 2*time.Second, 10)
	} else {
		//if no ip can be found from current claimed pool, claim next available valid ip
		if ip, err := p.claimNextValidIP(ip, cidr); err != nil {
			return "", err
		} else {
			return ip, utils.Retry(func() error { return p.cacheAllocatedIP(namespace, cidr, ip, machineName) }, 2*time.Second, 10)
		}
	}
}

func getCidrRange() (string, error) {
	if len(dockerNetworkCidr) == 0 {
		if nw, err := docker.GetNetwork(defaultDockerNetwork); err != nil {
			return dockerNetworkCidr, nil
		} else if len(nw.Cidr()) > 0 {
			dockerNetworkCidr = nw.Cidr()
		}
	}
	return dockerNetworkCidr, nil
}

func (p IpamProvisioner) findIpFromClaimedIpPool(namespace, machineName string) (bool, string, error) {
	if data, err := p.getIpConfigData(namespace); err != nil {
		return false, "", err
	} else if isClaimed, ip := p.isIpAlreadyClaimedForMachine(machineName, data); isClaimed { //if ip is already claimed by machine
		return true, ip, nil
	} else if isAvailable, ip := p.isAnyClaimedIpAvailable(data); isAvailable { //if any claimed ip is available for reclaim
		return true, ip, nil
	} else {
		return false, data[LastClaimedKey], nil
	}
	return false, "", nil
}

func (p IpamProvisioner) claimNextValidIP(lastIP, cidr string) (string, error) {
	for {
		if ip, err := p.claimNextIP(lastIP, cidr); err != nil {
			return "", err
		} else if isIpValid(ip) {
			return ip, nil
		} else {
			lastIP = ip
		}
	}
}

//check if ip is valid and add more rules if required
func isIpValid(ip string) bool {
	if len(ip) == 0 || ip[len(ip)-2:] == ".0" || ip[len(ip)-2:] == ".1" {
		return false
	}
	return isIpAvailableInNetwork(ip)
}

func isIpAvailableInNetwork(ip string) bool {
	if nw, err := docker.GetNetwork(defaultDockerNetwork); err == nil && len(nw.Containers) > 0{
		for _, c := range nw.Containers {
			container := *c
			if strings.Split(container.Ipv4, "/")[0] == ip { //this means ip is already claimed by another container
				return false
			}
		}
		return true //if ip doesn't match with any existing container ip then this means ip is available
	}
	return false
}

func (p IpamProvisioner) claimNextIP(lastIP, cidr string) (string, error) {
	if len(lastIP) == 0 {
		if ip, err := p.incrementIPBy(cidr, 520); err != nil {
			return "", err
		} else {
			lastIP = ip
		}
	}

	ip := net.ParseIP(lastIP)
	_, ipNet, err := net.ParseCIDR(cidr)
	if err != nil {
		return lastIP, err
	}
	for i := len(ip) - 1; i >= 0; i-- {
		ip[i]++
		if ip[i] != 0 {
			break
		}
	}
	if !ipNet.Contains(ip) {
		return lastIP, errors.New("overflowed CIDR while incrementing IP")
	}
	return ip.String(), nil
}

func (p IpamProvisioner) incrementIPBy(cidr string, count int) (string, error) {
	lIp := strings.Split(cidr, "/")[0]
	for i := 0; i < count; i++ {
		if ip, err := p.claimNextIP(lIp, cidr); err != nil {
			return "", err
		} else {
			lIp = ip
		}
	}
	return lIp, nil
}

func (p IpamProvisioner) isIpAlreadyClaimedForMachine(machineName string, data map[string]string) (bool, string) {
	if len(data) > 0 {
		for k, v := range data {
			if v == machineName {
				return true, k
			}
		}
	}
	return false, ""
}

func (p IpamProvisioner) isAnyClaimedIpAvailable(data map[string]string) (bool, string) {
	if len(data) > 0 {
		for k, v := range data {
			if v == Available {
				return true, k
			}
		}
	}
	return false, ""
}

func (p IpamProvisioner) ipDiscovery(namespace string) error {
	if data, err := p.getIpConfigData(namespace); err != nil {
		return err
	} else if len(data) > 0 {
		existingMachines, err := p.getDockerMachineNames(namespace)
		if err != nil {
			return err
		}
		existingMachineNames := fmt.Sprintf("%s,", strings.Join(existingMachines, ","))
		claimData := make(map[string]string)
		for k, v := range data {
			if k != CidrKey && k != LastClaimedKey && v[len(v)-3:] != "-lb"{
				if strings.Contains(existingMachineNames, fmt.Sprintf("%s,", v)) {
					claimData[k] = v
				} else {
					claimData[k] = Available
				}
			}
		}
		return p.cacheIpConfigData(namespace, claimData)
	}
	return nil
}

func (p IpamProvisioner) getDockerMachineNames(namespace string) ([]string, error) {
	list := &infrav1.DockerMachineList{}
	if err := p.Client.List(context.TODO(), list, client.InNamespace(namespace)); err != nil {
		return nil, err
	} else if len(list.Items) > 0 {
		names := make([]string, 0, 1)
		for _, m := range list.Items {
			names = append(names, m.Name)
		}
		return names, nil
	}
	return nil, nil
}

// Caching allocated ip and docker machine name in config map
func (p IpamProvisioner) cacheAllocatedIP(namespace, cidr, ip, machineName string) error {
	return p.cacheIpConfigData(namespace, map[string]string{
		CidrKey:        cidr,
		LastClaimedKey: ip,
		ip:             machineName,
	})
}

func (p IpamProvisioner) cacheIPReclaim(namespace, cidr, ip, machineName string) error {
	return p.cacheIpConfigData(namespace, map[string]string{
		CidrKey: cidr,
		ip:      machineName,
	})
}

func (p IpamProvisioner) getIpConfigData(namespace string) (map[string]string, error) {
	c := &corev1.ConfigMap{}
	key := client.ObjectKey{Namespace: namespace, Name: configMapName}
	if err := p.Get(context.TODO(), key, c); err != nil && !apierrors.IsNotFound(err) {
		return nil, errors.Wrapf(err, "failed to retrieve ipam-config")
	}
	return c.Data, nil
}

func (p IpamProvisioner) cacheIpConfigData(namespace string, data map[string]string) error {
	c := &corev1.ConfigMap{}
	key := client.ObjectKey{Namespace: namespace, Name: configMapName}
	if err := p.Get(context.TODO(), key, c); err != nil {
		if apierrors.IsNotFound(err) {
			cm := &corev1.ConfigMap{
				ObjectMeta: metav1.ObjectMeta{
					Name:      configMapName,
					Namespace: namespace,
				},
				Data: data,
			}

			if err := p.Create(context.TODO(), cm); err != nil && !apierrors.IsAlreadyExists(err) {
				return err
			}
		}
		return err
	} else {
		for k, v := range data {
			c.Data[k] = v
		}
		return p.Update(context.TODO(), c)
	}
}
