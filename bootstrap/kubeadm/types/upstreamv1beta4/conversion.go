/*
Copyright 2024 The Kubernetes Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package upstreamv1beta4

import (
	corev1 "k8s.io/api/core/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sort"
	"unsafe"

	"github.com/pkg/errors"
	apimachineryconversion "k8s.io/apimachinery/pkg/conversion"
	"sigs.k8s.io/controller-runtime/pkg/conversion"

	"sigs.k8s.io/cluster-api/bootstrap/kubeadm/api/v1beta1"
)

func (src *ClusterConfiguration) ConvertTo(dstRaw conversion.Hub) error {
	dst := dstRaw.(*v1beta1.ClusterConfiguration)
	return Convert_upstreamv1beta4_ClusterConfiguration_To_v1beta1_ClusterConfiguration(src, dst, nil)
}

func (dst *ClusterConfiguration) ConvertFrom(srcRaw conversion.Hub) error {
	src := srcRaw.(*v1beta1.ClusterConfiguration)
	return Convert_v1beta1_ClusterConfiguration_To_upstreamv1beta4_ClusterConfiguration(src, dst, nil)
}

func (src *InitConfiguration) ConvertTo(dstRaw conversion.Hub) error {
	dst := dstRaw.(*v1beta1.InitConfiguration)
	return Convert_upstreamv1beta4_InitConfiguration_To_v1beta1_InitConfiguration(src, dst, nil)
}

func (dst *InitConfiguration) ConvertFrom(srcRaw conversion.Hub) error {
	src := srcRaw.(*v1beta1.InitConfiguration)
	return Convert_v1beta1_InitConfiguration_To_upstreamv1beta4_InitConfiguration(src, dst, nil)
}

func (src *JoinConfiguration) ConvertTo(dstRaw conversion.Hub) error {
	dst := dstRaw.(*v1beta1.JoinConfiguration)
	return Convert_upstreamv1beta4_JoinConfiguration_To_v1beta1_JoinConfiguration(src, dst, nil)
}

func (dst *JoinConfiguration) ConvertFrom(srcRaw conversion.Hub) error {
	src := srcRaw.(*v1beta1.JoinConfiguration)
	return Convert_v1beta1_JoinConfiguration_To_upstreamv1beta4_JoinConfiguration(src, dst, nil)
}

// Custom conversion from this API, kubeadm v1beta4, to the hub version, CABPK v1beta1.

func Convert_upstreamv1beta4_ClusterConfiguration_To_v1beta1_ClusterConfiguration(in *ClusterConfiguration, out *v1beta1.ClusterConfiguration, s apimachineryconversion.Scope) error {
	// Following fields do not exist in CABPK v1beta1 version:
	// - Proxy (Not supported yet)
	// - EncryptionAlgorithm (Not supported yet)
	// - CertificateValidityPeriod (Not supported yet)
	// - CACertificateValidityPeriod (Not supported yet)
	return autoConvert_upstreamv1beta4_ClusterConfiguration_To_v1beta1_ClusterConfiguration(in, out, s)
}

func Convert_upstreamv1beta4_ControlPlaneComponent_To_v1beta1_ControlPlaneComponent(in *ControlPlaneComponent, out *v1beta1.ControlPlaneComponent, s apimachineryconversion.Scope) error {
	// Following fields do not exist in CABPK v1beta1 version:
	// - ExtraEnvs (Not supported yet)

	// Following fields exists in CABPK v1beta1 but they need a custom conversions.
	// Note: there is a potential info loss when there are two values for the same arg but this is not an issue because the CAPBK v1beta1 does not allow this use case.
	out.ExtraArgs = convertFromArgs(in.ExtraArgs)
	return autoConvert_upstreamv1beta4_ControlPlaneComponent_To_v1beta1_ControlPlaneComponent(in, out, s)
}

func Convert_upstreamv1beta4_LocalEtcd_To_v1beta1_LocalEtcd(in *LocalEtcd, out *v1beta1.LocalEtcd, s apimachineryconversion.Scope) error {
	// Following fields do not exist in CABPK v1beta1 version:
	// - ExtraEnvs (Not supported yet)

	// Following fields require a custom conversions.
	// Note: there is a potential info loss when there are two values for the same arg but this is not an issue because the CAPBK v1beta1 does not allow this use case.
	out.ExtraArgs = convertFromArgs(in.ExtraArgs)
	return autoConvert_upstreamv1beta4_LocalEtcd_To_v1beta1_LocalEtcd(in, out, s)
}

func Convert_upstreamv1beta4_DNS_To_v1beta1_DNS(in *DNS, out *v1beta1.DNS, s apimachineryconversion.Scope) error {
	// Following fields do not exist in CABPK v1beta1 version:
	// - Disabled (Not supported yet)
	return autoConvert_upstreamv1beta4_DNS_To_v1beta1_DNS(in, out, s)
}

func Convert_upstreamv1beta4_InitConfiguration_To_v1beta1_InitConfiguration(in *InitConfiguration, out *v1beta1.InitConfiguration, s apimachineryconversion.Scope) error {
	// Following fields do not exist in CABPK v1beta1 version:
	// - DryRun (Does not make sense for CAPBK)
	// - CertificateKey (CABPK does not use automatic copy certs)
	// - Timeouts (Not supported yet)
	return autoConvert_upstreamv1beta4_InitConfiguration_To_v1beta1_InitConfiguration(in, out, s)
}

func Convert_upstreamv1beta4_JoinConfiguration_To_v1beta1_JoinConfiguration(in *JoinConfiguration, out *v1beta1.JoinConfiguration, s apimachineryconversion.Scope) error {
	// Following fields do not exist in CABPK v1beta1 version:
	// - DryRun (Does not make sense for CAPBK)
	// - Timeouts (Not supported yet)
	err := autoConvert_upstreamv1beta4_JoinConfiguration_To_v1beta1_JoinConfiguration(in, out, s)

	// Handle migration of JoinConfiguration.Timeouts.TLSBootstrap to Discovery.Timeout.
	if in.Timeouts != nil && in.Timeouts.TLSBootstrap != nil {
		out.Discovery.Timeout = in.Timeouts.TLSBootstrap
	}

	return err
}

func Convert_upstreamv1beta4_NodeRegistrationOptions_To_v1beta1_NodeRegistrationOptions(in *NodeRegistrationOptions, out *v1beta1.NodeRegistrationOptions, s apimachineryconversion.Scope) error {
	// Following fields do not exist in CABPK v1beta1 version:
	// - ImagePullSerial (Not supported yet)

	// Following fields require a custom conversions.
	// Note: there is a potential info loss when there are two values for the same arg but this is not an issue because the CAPBK v1beta1 does not allow this use case.
	out.KubeletExtraArgs = convertFromArgs(in.KubeletExtraArgs)
	return autoConvert_upstreamv1beta4_NodeRegistrationOptions_To_v1beta1_NodeRegistrationOptions(in, out, s)
}

func Convert_upstreamv1beta4_JoinControlPlane_To_v1beta1_JoinControlPlane(in *JoinControlPlane, out *v1beta1.JoinControlPlane, s apimachineryconversion.Scope) error {
	// Following fields do not exist in CABPK v1beta1 version:
	// - CertificateKey (CABPK does not use automatic copy certs)
	return autoConvert_upstreamv1beta4_JoinControlPlane_To_v1beta1_JoinControlPlane(in, out, s)
}

// Custom conversion from the hub version, CABPK v1beta1, to this API, kubeadm v1beta4.

func Convert_v1beta1_ControlPlaneComponent_To_upstreamv1beta4_ControlPlaneComponent(in *v1beta1.ControlPlaneComponent, out *ControlPlaneComponent, s apimachineryconversion.Scope) error {
	// Following fields require a custom conversions.
	out.ExtraArgs = convertToArgs(in.ExtraArgs)
	return autoConvert_v1beta1_ControlPlaneComponent_To_upstreamv1beta4_ControlPlaneComponent(in, out, s)
}

func Convert_v1beta1_APIServer_To_upstreamv1beta4_APIServer(in *v1beta1.APIServer, out *APIServer, s apimachineryconversion.Scope) error {
	// Following fields do not exist in kubeadm v1beta4 version:
	// - TimeoutForControlPlane (this field has been migrated to Init/JoinConfiguration; migration is handled by ConvertFromClusterConfiguration custom converters.
	return autoConvert_v1beta1_APIServer_To_upstreamv1beta4_APIServer(in, out, s)
}

func Convert_v1beta1_LocalEtcd_To_upstreamv1beta4_LocalEtcd(in *v1beta1.LocalEtcd, out *LocalEtcd, s apimachineryconversion.Scope) error {
	// Following fields require a custom conversions.
	out.ExtraArgs = convertToArgs(in.ExtraArgs)
	return autoConvert_v1beta1_LocalEtcd_To_upstreamv1beta4_LocalEtcd(in, out, s)
}

func Convert_v1beta1_JoinConfiguration_To_upstreamv1beta4_JoinConfiguration(in *v1beta1.JoinConfiguration, out *JoinConfiguration, s apimachineryconversion.Scope) error {
	err := autoConvert_v1beta1_JoinConfiguration_To_upstreamv1beta4_JoinConfiguration(in, out, s)

	// Handle migration of Discovery.Timeout to JoinConfiguration.Timeouts.TLSBootstrap.
	if in.Discovery.Timeout != nil {
		if out.Timeouts == nil {
			out.Timeouts = &Timeouts{}
		}
		out.Timeouts.TLSBootstrap = in.Discovery.Timeout
	}
	return err
}

func Convert_v1beta1_NodeRegistrationOptions_To_upstreamv1beta4_NodeRegistrationOptions(in *v1beta1.NodeRegistrationOptions, out *NodeRegistrationOptions, s apimachineryconversion.Scope) error {
	// Following fields exists in kubeadm v1beta4 types and can be converted to CAPBK v1beta1.
	out.KubeletExtraArgs = convertToArgs(in.KubeletExtraArgs)
	return autoConvert_v1beta1_NodeRegistrationOptions_To_upstreamv1beta4_NodeRegistrationOptions(in, out, s)
}

func Convert_v1beta1_Discovery_To_upstreamv1beta4_Discovery(in *v1beta1.Discovery, out *Discovery, s apimachineryconversion.Scope) error {
	// Following fields do not exist in kubeadm v1beta4 version:
	// - Timeout (this field has been migrated to JoinConfiguration.Timeouts.TLSBootstrap, the conversion is handled in Convert_v1beta1_JoinConfiguration_To_upstreamv1beta4_JoinConfiguration)
	return autoConvert_v1beta1_Discovery_To_upstreamv1beta4_Discovery(in, out, s)
}

// convertToArgs takes a argument map and converts it to a slice of arguments.
// Te resulting argument slice is sorted alpha-numerically.
func convertToArgs(in map[string]string) []Arg {
	if in == nil {
		return nil
	}
	args := make([]Arg, 0, len(in))
	for k, v := range in {
		args = append(args, Arg{Name: k, Value: v})
	}
	sort.Slice(args, func(i, j int) bool {
		if args[i].Name == args[j].Name {
			return args[i].Value < args[j].Value
		}
		return args[i].Name < args[j].Name
	})
	return args
}

// convertFromArgs takes a slice of arguments and returns an argument map.
// Duplicate argument keys will be de-duped, where later keys will take precedence.
func convertFromArgs(in []Arg) map[string]string {
	if in == nil {
		return nil
	}
	args := make(map[string]string, len(in))
	for _, arg := range in {
		args[arg.Name] = arg.Value
	}
	return args
}

// Custom conversions to handle fields migrated from ClusterConfiguration to Init and JoinConfiguration in the kubeadm v1beta4 API version.

func (dst *InitConfiguration) ConvertFromClusterConfiguration(clusterConfiguration *v1beta1.ClusterConfiguration) error {
	if clusterConfiguration == nil || clusterConfiguration.APIServer.TimeoutForControlPlane == nil {
		return nil
	}

	if dst.Timeouts == nil {
		dst.Timeouts = &Timeouts{}
	}
	dst.Timeouts.ControlPlaneComponentHealthCheck = clusterConfiguration.APIServer.TimeoutForControlPlane
	return nil
}

func (dst *JoinConfiguration) ConvertFromClusterConfiguration(clusterConfiguration *v1beta1.ClusterConfiguration) error {
	if clusterConfiguration == nil || clusterConfiguration.APIServer.TimeoutForControlPlane == nil {
		return nil
	}

	if dst.Timeouts == nil {
		dst.Timeouts = &Timeouts{}
	}
	dst.Timeouts.ControlPlaneComponentHealthCheck = clusterConfiguration.APIServer.TimeoutForControlPlane
	return nil
}

func (src *InitConfiguration) ConvertToClusterConfiguration(clusterConfiguration *v1beta1.ClusterConfiguration) error {
	if src.Timeouts == nil || src.Timeouts.ControlPlaneComponentHealthCheck == nil {
		return nil
	}

	if clusterConfiguration == nil {
		return errors.New("cannot convert InitConfiguration to a nil ClusterConfiguration")
	}
	clusterConfiguration.APIServer.TimeoutForControlPlane = src.Timeouts.ControlPlaneComponentHealthCheck
	return nil
}

func (src *JoinConfiguration) ConvertToClusterConfiguration(clusterConfiguration *v1beta1.ClusterConfiguration) error {
	if src.Timeouts == nil || src.Timeouts.ControlPlaneComponentHealthCheck == nil {
		return nil
	}

	if clusterConfiguration == nil {
		return errors.New("cannot convert JoinConfiguration to a nil ClusterConfiguration")
	}
	clusterConfiguration.APIServer.TimeoutForControlPlane = src.Timeouts.ControlPlaneComponentHealthCheck
	return nil
}

func autoConvert_upstreamv1beta4_APIEndpoint_To_v1beta1_APIEndpoint(in *APIEndpoint, out *v1beta1.APIEndpoint, s apimachineryconversion.Scope) error {
	out.AdvertiseAddress = in.AdvertiseAddress
	out.BindPort = in.BindPort
	return nil
}

// Convert_upstreamv1beta4_APIEndpoint_To_v1beta1_APIEndpoint is an autogenerated conversion function.
func Convert_upstreamv1beta4_APIEndpoint_To_v1beta1_APIEndpoint(in *APIEndpoint, out *v1beta1.APIEndpoint, s apimachineryconversion.Scope) error {
	return autoConvert_upstreamv1beta4_APIEndpoint_To_v1beta1_APIEndpoint(in, out, s)
}

func autoConvert_v1beta1_APIEndpoint_To_upstreamv1beta4_APIEndpoint(in *v1beta1.APIEndpoint, out *APIEndpoint, s apimachineryconversion.Scope) error {
	out.AdvertiseAddress = in.AdvertiseAddress
	out.BindPort = in.BindPort
	return nil
}

// Convert_v1beta1_APIEndpoint_To_upstreamv1beta4_APIEndpoint is an autogenerated conversion function.
func Convert_v1beta1_APIEndpoint_To_upstreamv1beta4_APIEndpoint(in *v1beta1.APIEndpoint, out *APIEndpoint, s apimachineryconversion.Scope) error {
	return autoConvert_v1beta1_APIEndpoint_To_upstreamv1beta4_APIEndpoint(in, out, s)
}

func autoConvert_upstreamv1beta4_APIServer_To_v1beta1_APIServer(in *APIServer, out *v1beta1.APIServer, s apimachineryconversion.Scope) error {
	if err := Convert_upstreamv1beta4_ControlPlaneComponent_To_v1beta1_ControlPlaneComponent(&in.ControlPlaneComponent, &out.ControlPlaneComponent, s); err != nil {
		return err
	}
	out.CertSANs = *(*[]string)(unsafe.Pointer(&in.CertSANs))
	return nil
}

// Convert_upstreamv1beta4_APIServer_To_v1beta1_APIServer is an autogenerated conversion function.
func Convert_upstreamv1beta4_APIServer_To_v1beta1_APIServer(in *APIServer, out *v1beta1.APIServer, s apimachineryconversion.Scope) error {
	return autoConvert_upstreamv1beta4_APIServer_To_v1beta1_APIServer(in, out, s)
}

func autoConvert_v1beta1_APIServer_To_upstreamv1beta4_APIServer(in *v1beta1.APIServer, out *APIServer, s apimachineryconversion.Scope) error {
	if err := Convert_v1beta1_ControlPlaneComponent_To_upstreamv1beta4_ControlPlaneComponent(&in.ControlPlaneComponent, &out.ControlPlaneComponent, s); err != nil {
		return err
	}
	out.CertSANs = *(*[]string)(unsafe.Pointer(&in.CertSANs))
	// WARNING: in.TimeoutForControlPlane requires manual conversion: does not exist in peer-type
	return nil
}

func autoConvert_upstreamv1beta4_BootstrapToken_To_v1beta1_BootstrapToken(in *BootstrapToken, out *v1beta1.BootstrapToken, s apimachineryconversion.Scope) error {
	out.Token = (*v1beta1.BootstrapTokenString)(unsafe.Pointer(in.Token))
	out.Description = in.Description
	out.TTL = (*v1.Duration)(unsafe.Pointer(in.TTL))
	out.Expires = (*v1.Time)(unsafe.Pointer(in.Expires))
	out.Usages = *(*[]string)(unsafe.Pointer(&in.Usages))
	out.Groups = *(*[]string)(unsafe.Pointer(&in.Groups))
	return nil
}

// Convert_upstreamv1beta4_BootstrapToken_To_v1beta1_BootstrapToken is an autogenerated conversion function.
func Convert_upstreamv1beta4_BootstrapToken_To_v1beta1_BootstrapToken(in *BootstrapToken, out *v1beta1.BootstrapToken, s apimachineryconversion.Scope) error {
	return autoConvert_upstreamv1beta4_BootstrapToken_To_v1beta1_BootstrapToken(in, out, s)
}

func autoConvert_v1beta1_BootstrapToken_To_upstreamv1beta4_BootstrapToken(in *v1beta1.BootstrapToken, out *BootstrapToken, s apimachineryconversion.Scope) error {
	out.Token = (*BootstrapTokenString)(unsafe.Pointer(in.Token))
	out.Description = in.Description
	out.TTL = (*v1.Duration)(unsafe.Pointer(in.TTL))
	out.Expires = (*v1.Time)(unsafe.Pointer(in.Expires))
	out.Usages = *(*[]string)(unsafe.Pointer(&in.Usages))
	out.Groups = *(*[]string)(unsafe.Pointer(&in.Groups))
	return nil
}

// Convert_v1beta1_BootstrapToken_To_upstreamv1beta4_BootstrapToken is an autogenerated conversion function.
func Convert_v1beta1_BootstrapToken_To_upstreamv1beta4_BootstrapToken(in *v1beta1.BootstrapToken, out *BootstrapToken, s apimachineryconversion.Scope) error {
	return autoConvert_v1beta1_BootstrapToken_To_upstreamv1beta4_BootstrapToken(in, out, s)
}

func autoConvert_upstreamv1beta4_BootstrapTokenDiscovery_To_v1beta1_BootstrapTokenDiscovery(in *BootstrapTokenDiscovery, out *v1beta1.BootstrapTokenDiscovery, s apimachineryconversion.Scope) error {
	out.Token = in.Token
	out.APIServerEndpoint = in.APIServerEndpoint
	out.CACertHashes = *(*[]string)(unsafe.Pointer(&in.CACertHashes))
	out.UnsafeSkipCAVerification = in.UnsafeSkipCAVerification
	return nil
}

// Convert_upstreamv1beta4_BootstrapTokenDiscovery_To_v1beta1_BootstrapTokenDiscovery is an autogenerated conversion function.
func Convert_upstreamv1beta4_BootstrapTokenDiscovery_To_v1beta1_BootstrapTokenDiscovery(in *BootstrapTokenDiscovery, out *v1beta1.BootstrapTokenDiscovery, s apimachineryconversion.Scope) error {
	return autoConvert_upstreamv1beta4_BootstrapTokenDiscovery_To_v1beta1_BootstrapTokenDiscovery(in, out, s)
}

func autoConvert_v1beta1_BootstrapTokenDiscovery_To_upstreamv1beta4_BootstrapTokenDiscovery(in *v1beta1.BootstrapTokenDiscovery, out *BootstrapTokenDiscovery, s apimachineryconversion.Scope) error {
	out.Token = in.Token
	out.APIServerEndpoint = in.APIServerEndpoint
	out.CACertHashes = *(*[]string)(unsafe.Pointer(&in.CACertHashes))
	out.UnsafeSkipCAVerification = in.UnsafeSkipCAVerification
	return nil
}

// Convert_v1beta1_BootstrapTokenDiscovery_To_upstreamv1beta4_BootstrapTokenDiscovery is an autogenerated conversion function.
func Convert_v1beta1_BootstrapTokenDiscovery_To_upstreamv1beta4_BootstrapTokenDiscovery(in *v1beta1.BootstrapTokenDiscovery, out *BootstrapTokenDiscovery, s apimachineryconversion.Scope) error {
	return autoConvert_v1beta1_BootstrapTokenDiscovery_To_upstreamv1beta4_BootstrapTokenDiscovery(in, out, s)
}

func autoConvert_upstreamv1beta4_BootstrapTokenString_To_v1beta1_BootstrapTokenString(in *BootstrapTokenString, out *v1beta1.BootstrapTokenString, s apimachineryconversion.Scope) error {
	out.ID = in.ID
	out.Secret = in.Secret
	return nil
}

// Convert_upstreamv1beta4_BootstrapTokenString_To_v1beta1_BootstrapTokenString is an autogenerated conversion function.
func Convert_upstreamv1beta4_BootstrapTokenString_To_v1beta1_BootstrapTokenString(in *BootstrapTokenString, out *v1beta1.BootstrapTokenString, s apimachineryconversion.Scope) error {
	return autoConvert_upstreamv1beta4_BootstrapTokenString_To_v1beta1_BootstrapTokenString(in, out, s)
}

func autoConvert_v1beta1_BootstrapTokenString_To_upstreamv1beta4_BootstrapTokenString(in *v1beta1.BootstrapTokenString, out *BootstrapTokenString, s apimachineryconversion.Scope) error {
	out.ID = in.ID
	out.Secret = in.Secret
	return nil
}

// Convert_v1beta1_BootstrapTokenString_To_upstreamv1beta4_BootstrapTokenString is an autogenerated conversion function.
func Convert_v1beta1_BootstrapTokenString_To_upstreamv1beta4_BootstrapTokenString(in *v1beta1.BootstrapTokenString, out *BootstrapTokenString, s apimachineryconversion.Scope) error {
	return autoConvert_v1beta1_BootstrapTokenString_To_upstreamv1beta4_BootstrapTokenString(in, out, s)
}

func autoConvert_upstreamv1beta4_ClusterConfiguration_To_v1beta1_ClusterConfiguration(in *ClusterConfiguration, out *v1beta1.ClusterConfiguration, s apimachineryconversion.Scope) error {
	if err := Convert_upstreamv1beta4_Etcd_To_v1beta1_Etcd(&in.Etcd, &out.Etcd, s); err != nil {
		return err
	}
	if err := Convert_upstreamv1beta4_Networking_To_v1beta1_Networking(&in.Networking, &out.Networking, s); err != nil {
		return err
	}
	out.KubernetesVersion = in.KubernetesVersion
	out.ControlPlaneEndpoint = in.ControlPlaneEndpoint
	if err := Convert_upstreamv1beta4_APIServer_To_v1beta1_APIServer(&in.APIServer, &out.APIServer, s); err != nil {
		return err
	}
	if err := Convert_upstreamv1beta4_ControlPlaneComponent_To_v1beta1_ControlPlaneComponent(&in.ControllerManager, &out.ControllerManager, s); err != nil {
		return err
	}
	if err := Convert_upstreamv1beta4_ControlPlaneComponent_To_v1beta1_ControlPlaneComponent(&in.Scheduler, &out.Scheduler, s); err != nil {
		return err
	}
	if err := Convert_upstreamv1beta4_DNS_To_v1beta1_DNS(&in.DNS, &out.DNS, s); err != nil {
		return err
	}
	// WARNING: in.Proxy requires manual conversion: does not exist in peer-type
	out.CertificatesDir = in.CertificatesDir
	out.ImageRepository = in.ImageRepository
	out.FeatureGates = *(*map[string]bool)(unsafe.Pointer(&in.FeatureGates))
	out.ClusterName = in.ClusterName
	// WARNING: in.EncryptionAlgorithm requires manual conversion: does not exist in peer-type
	// WARNING: in.CertificateValidityPeriod requires manual conversion: does not exist in peer-type
	// WARNING: in.CACertificateValidityPeriod requires manual conversion: does not exist in peer-type
	return nil
}

func autoConvert_v1beta1_ClusterConfiguration_To_upstreamv1beta4_ClusterConfiguration(in *v1beta1.ClusterConfiguration, out *ClusterConfiguration, s apimachineryconversion.Scope) error {
	if err := Convert_v1beta1_Etcd_To_upstreamv1beta4_Etcd(&in.Etcd, &out.Etcd, s); err != nil {
		return err
	}
	if err := Convert_v1beta1_Networking_To_upstreamv1beta4_Networking(&in.Networking, &out.Networking, s); err != nil {
		return err
	}
	out.KubernetesVersion = in.KubernetesVersion
	out.ControlPlaneEndpoint = in.ControlPlaneEndpoint
	if err := Convert_v1beta1_APIServer_To_upstreamv1beta4_APIServer(&in.APIServer, &out.APIServer, s); err != nil {
		return err
	}
	if err := Convert_v1beta1_ControlPlaneComponent_To_upstreamv1beta4_ControlPlaneComponent(&in.ControllerManager, &out.ControllerManager, s); err != nil {
		return err
	}
	if err := Convert_v1beta1_ControlPlaneComponent_To_upstreamv1beta4_ControlPlaneComponent(&in.Scheduler, &out.Scheduler, s); err != nil {
		return err
	}
	if err := Convert_v1beta1_DNS_To_upstreamv1beta4_DNS(&in.DNS, &out.DNS, s); err != nil {
		return err
	}
	out.CertificatesDir = in.CertificatesDir
	out.ImageRepository = in.ImageRepository
	out.FeatureGates = *(*map[string]bool)(unsafe.Pointer(&in.FeatureGates))
	out.ClusterName = in.ClusterName
	return nil
}

// Convert_v1beta1_ClusterConfiguration_To_upstreamv1beta4_ClusterConfiguration is an autogenerated conversion function.
func Convert_v1beta1_ClusterConfiguration_To_upstreamv1beta4_ClusterConfiguration(in *v1beta1.ClusterConfiguration, out *ClusterConfiguration, s apimachineryconversion.Scope) error {
	return autoConvert_v1beta1_ClusterConfiguration_To_upstreamv1beta4_ClusterConfiguration(in, out, s)
}

func autoConvert_upstreamv1beta4_ControlPlaneComponent_To_v1beta1_ControlPlaneComponent(in *ControlPlaneComponent, out *v1beta1.ControlPlaneComponent, s apimachineryconversion.Scope) error {
	// WARNING: in.ExtraArgs requires manual conversion: inconvertible types ([]sigs.k8s.io/cluster-api/bootstrap/kubeadm/types/upstreamv1beta4.Arg vs map[string]string)
	out.ExtraVolumes = *(*[]v1beta1.HostPathMount)(unsafe.Pointer(&in.ExtraVolumes))
	//out.ExtraEnvs = *(*[]v1beta1.EnvVar)(unsafe.Pointer(&in.ExtraEnvs))
	return nil
}

func autoConvert_v1beta1_ControlPlaneComponent_To_upstreamv1beta4_ControlPlaneComponent(in *v1beta1.ControlPlaneComponent, out *ControlPlaneComponent, s apimachineryconversion.Scope) error {
	// WARNING: in.ExtraArgs requires manual conversion: inconvertible types (map[string]string vs []sigs.k8s.io/cluster-api/bootstrap/kubeadm/types/upstreamv1beta4.Arg)
	out.ExtraVolumes = *(*[]HostPathMount)(unsafe.Pointer(&in.ExtraVolumes))
	//out.ExtraEnvs = *(*[]EnvVar)(unsafe.Pointer(&in.ExtraEnvs))
	return nil
}

func autoConvert_upstreamv1beta4_DNS_To_v1beta1_DNS(in *DNS, out *v1beta1.DNS, s apimachineryconversion.Scope) error {
	if err := Convert_upstreamv1beta4_ImageMeta_To_v1beta1_ImageMeta(&in.ImageMeta, &out.ImageMeta, s); err != nil {
		return err
	}
	// WARNING: in.Disabled requires manual conversion: does not exist in peer-type
	return nil
}

func autoConvert_v1beta1_DNS_To_upstreamv1beta4_DNS(in *v1beta1.DNS, out *DNS, s apimachineryconversion.Scope) error {
	if err := Convert_v1beta1_ImageMeta_To_upstreamv1beta4_ImageMeta(&in.ImageMeta, &out.ImageMeta, s); err != nil {
		return err
	}
	return nil
}

// Convert_v1beta1_DNS_To_upstreamv1beta4_DNS is an autogenerated conversion function.
func Convert_v1beta1_DNS_To_upstreamv1beta4_DNS(in *v1beta1.DNS, out *DNS, s apimachineryconversion.Scope) error {
	return autoConvert_v1beta1_DNS_To_upstreamv1beta4_DNS(in, out, s)
}

func autoConvert_upstreamv1beta4_Discovery_To_v1beta1_Discovery(in *Discovery, out *v1beta1.Discovery, s apimachineryconversion.Scope) error {
	out.BootstrapToken = (*v1beta1.BootstrapTokenDiscovery)(unsafe.Pointer(in.BootstrapToken))
	if in.File != nil {
		in, out := &in.File, &out.File
		*out = new(v1beta1.FileDiscovery)
		if err := Convert_upstreamv1beta4_FileDiscovery_To_v1beta1_FileDiscovery(*in, *out, s); err != nil {
			return err
		}
	} else {
		out.File = nil
	}
	out.TLSBootstrapToken = in.TLSBootstrapToken
	return nil
}

// Convert_upstreamv1beta4_Discovery_To_v1beta1_Discovery is an autogenerated conversion function.
func Convert_upstreamv1beta4_Discovery_To_v1beta1_Discovery(in *Discovery, out *v1beta1.Discovery, s apimachineryconversion.Scope) error {
	return autoConvert_upstreamv1beta4_Discovery_To_v1beta1_Discovery(in, out, s)
}

func autoConvert_v1beta1_Discovery_To_upstreamv1beta4_Discovery(in *v1beta1.Discovery, out *Discovery, s apimachineryconversion.Scope) error {
	out.BootstrapToken = (*BootstrapTokenDiscovery)(unsafe.Pointer(in.BootstrapToken))
	if in.File != nil {
		in, out := &in.File, &out.File
		*out = new(FileDiscovery)
		if err := Convert_v1beta1_FileDiscovery_To_upstreamv1beta4_FileDiscovery(*in, *out, s); err != nil {
			return err
		}
	} else {
		out.File = nil
	}
	out.TLSBootstrapToken = in.TLSBootstrapToken
	// WARNING: in.Timeout requires manual conversion: does not exist in peer-type
	return nil
}

func autoConvert_upstreamv1beta4_Etcd_To_v1beta1_Etcd(in *Etcd, out *v1beta1.Etcd, s apimachineryconversion.Scope) error {
	if in.Local != nil {
		in, out := &in.Local, &out.Local
		*out = new(v1beta1.LocalEtcd)
		if err := Convert_upstreamv1beta4_LocalEtcd_To_v1beta1_LocalEtcd(*in, *out, s); err != nil {
			return err
		}
	} else {
		out.Local = nil
	}
	out.External = (*v1beta1.ExternalEtcd)(unsafe.Pointer(in.External))
	return nil
}

// Convert_upstreamv1beta4_Etcd_To_v1beta1_Etcd is an autogenerated conversion function.
func Convert_upstreamv1beta4_Etcd_To_v1beta1_Etcd(in *Etcd, out *v1beta1.Etcd, s apimachineryconversion.Scope) error {
	return autoConvert_upstreamv1beta4_Etcd_To_v1beta1_Etcd(in, out, s)
}

func autoConvert_v1beta1_Etcd_To_upstreamv1beta4_Etcd(in *v1beta1.Etcd, out *Etcd, s apimachineryconversion.Scope) error {
	if in.Local != nil {
		in, out := &in.Local, &out.Local
		*out = new(LocalEtcd)
		if err := Convert_v1beta1_LocalEtcd_To_upstreamv1beta4_LocalEtcd(*in, *out, s); err != nil {
			return err
		}
	} else {
		out.Local = nil
	}
	out.External = (*ExternalEtcd)(unsafe.Pointer(in.External))
	return nil
}

// Convert_v1beta1_Etcd_To_upstreamv1beta4_Etcd is an autogenerated conversion function.
func Convert_v1beta1_Etcd_To_upstreamv1beta4_Etcd(in *v1beta1.Etcd, out *Etcd, s apimachineryconversion.Scope) error {
	return autoConvert_v1beta1_Etcd_To_upstreamv1beta4_Etcd(in, out, s)
}

func autoConvert_upstreamv1beta4_ExternalEtcd_To_v1beta1_ExternalEtcd(in *ExternalEtcd, out *v1beta1.ExternalEtcd, s apimachineryconversion.Scope) error {
	out.Endpoints = *(*[]string)(unsafe.Pointer(&in.Endpoints))
	out.CAFile = in.CAFile
	out.CertFile = in.CertFile
	out.KeyFile = in.KeyFile
	return nil
}

// Convert_upstreamv1beta4_ExternalEtcd_To_v1beta1_ExternalEtcd is an autogenerated conversion function.
func Convert_upstreamv1beta4_ExternalEtcd_To_v1beta1_ExternalEtcd(in *ExternalEtcd, out *v1beta1.ExternalEtcd, s apimachineryconversion.Scope) error {
	return autoConvert_upstreamv1beta4_ExternalEtcd_To_v1beta1_ExternalEtcd(in, out, s)
}

func autoConvert_v1beta1_ExternalEtcd_To_upstreamv1beta4_ExternalEtcd(in *v1beta1.ExternalEtcd, out *ExternalEtcd, s apimachineryconversion.Scope) error {
	out.Endpoints = *(*[]string)(unsafe.Pointer(&in.Endpoints))
	out.CAFile = in.CAFile
	out.CertFile = in.CertFile
	out.KeyFile = in.KeyFile
	return nil
}

// Convert_v1beta1_ExternalEtcd_To_upstreamv1beta4_ExternalEtcd is an autogenerated conversion function.
func Convert_v1beta1_ExternalEtcd_To_upstreamv1beta4_ExternalEtcd(in *v1beta1.ExternalEtcd, out *ExternalEtcd, s apimachineryconversion.Scope) error {
	return autoConvert_v1beta1_ExternalEtcd_To_upstreamv1beta4_ExternalEtcd(in, out, s)
}

func autoConvert_upstreamv1beta4_FileDiscovery_To_v1beta1_FileDiscovery(in *FileDiscovery, out *v1beta1.FileDiscovery, s apimachineryconversion.Scope) error {
	out.KubeConfigPath = in.KubeConfigPath
	return nil
}

// Convert_upstreamv1beta4_FileDiscovery_To_v1beta1_FileDiscovery is an autogenerated conversion function.
func Convert_upstreamv1beta4_FileDiscovery_To_v1beta1_FileDiscovery(in *FileDiscovery, out *v1beta1.FileDiscovery, s apimachineryconversion.Scope) error {
	return autoConvert_upstreamv1beta4_FileDiscovery_To_v1beta1_FileDiscovery(in, out, s)
}

func autoConvert_v1beta1_FileDiscovery_To_upstreamv1beta4_FileDiscovery(in *v1beta1.FileDiscovery, out *FileDiscovery, s apimachineryconversion.Scope) error {
	out.KubeConfigPath = in.KubeConfigPath
	// WARNING: in.KubeConfig requires manual conversion: does not exist in peer-type
	return nil
}

func autoConvert_upstreamv1beta4_HostPathMount_To_v1beta1_HostPathMount(in *HostPathMount, out *v1beta1.HostPathMount, s apimachineryconversion.Scope) error {
	out.Name = in.Name
	out.HostPath = in.HostPath
	out.MountPath = in.MountPath
	out.ReadOnly = in.ReadOnly
	out.PathType = corev1.HostPathType(in.PathType)
	return nil
}

// Convert_upstreamv1beta4_HostPathMount_To_v1beta1_HostPathMount is an autogenerated conversion function.
func Convert_upstreamv1beta4_HostPathMount_To_v1beta1_HostPathMount(in *HostPathMount, out *v1beta1.HostPathMount, s apimachineryconversion.Scope) error {
	return autoConvert_upstreamv1beta4_HostPathMount_To_v1beta1_HostPathMount(in, out, s)
}

func autoConvert_v1beta1_HostPathMount_To_upstreamv1beta4_HostPathMount(in *v1beta1.HostPathMount, out *HostPathMount, s apimachineryconversion.Scope) error {
	out.Name = in.Name
	out.HostPath = in.HostPath
	out.MountPath = in.MountPath
	out.ReadOnly = in.ReadOnly
	out.PathType = corev1.HostPathType(in.PathType)
	return nil
}

// Convert_v1beta1_HostPathMount_To_upstreamv1beta4_HostPathMount is an autogenerated conversion function.
func Convert_v1beta1_HostPathMount_To_upstreamv1beta4_HostPathMount(in *v1beta1.HostPathMount, out *HostPathMount, s apimachineryconversion.Scope) error {
	return autoConvert_v1beta1_HostPathMount_To_upstreamv1beta4_HostPathMount(in, out, s)
}

func autoConvert_upstreamv1beta4_ImageMeta_To_v1beta1_ImageMeta(in *ImageMeta, out *v1beta1.ImageMeta, s apimachineryconversion.Scope) error {
	out.ImageRepository = in.ImageRepository
	out.ImageTag = in.ImageTag
	return nil
}

// Convert_upstreamv1beta4_ImageMeta_To_v1beta1_ImageMeta is an autogenerated conversion function.
func Convert_upstreamv1beta4_ImageMeta_To_v1beta1_ImageMeta(in *ImageMeta, out *v1beta1.ImageMeta, s apimachineryconversion.Scope) error {
	return autoConvert_upstreamv1beta4_ImageMeta_To_v1beta1_ImageMeta(in, out, s)
}

func autoConvert_v1beta1_ImageMeta_To_upstreamv1beta4_ImageMeta(in *v1beta1.ImageMeta, out *ImageMeta, s apimachineryconversion.Scope) error {
	out.ImageRepository = in.ImageRepository
	out.ImageTag = in.ImageTag
	return nil
}

// Convert_v1beta1_ImageMeta_To_upstreamv1beta4_ImageMeta is an autogenerated conversion function.
func Convert_v1beta1_ImageMeta_To_upstreamv1beta4_ImageMeta(in *v1beta1.ImageMeta, out *ImageMeta, s apimachineryconversion.Scope) error {
	return autoConvert_v1beta1_ImageMeta_To_upstreamv1beta4_ImageMeta(in, out, s)
}

func autoConvert_upstreamv1beta4_InitConfiguration_To_v1beta1_InitConfiguration(in *InitConfiguration, out *v1beta1.InitConfiguration, s apimachineryconversion.Scope) error {
	out.BootstrapTokens = *(*[]v1beta1.BootstrapToken)(unsafe.Pointer(&in.BootstrapTokens))
	// WARNING: in.DryRun requires manual conversion: does not exist in peer-type
	if err := Convert_upstreamv1beta4_NodeRegistrationOptions_To_v1beta1_NodeRegistrationOptions(&in.NodeRegistration, &out.NodeRegistration, s); err != nil {
		return err
	}
	if err := Convert_upstreamv1beta4_APIEndpoint_To_v1beta1_APIEndpoint(&in.LocalAPIEndpoint, &out.LocalAPIEndpoint, s); err != nil {
		return err
	}
	// WARNING: in.CertificateKey requires manual conversion: does not exist in peer-type
	out.SkipPhases = *(*[]string)(unsafe.Pointer(&in.SkipPhases))
	out.Patches = (*v1beta1.Patches)(unsafe.Pointer(in.Patches))
	// WARNING: in.Timeouts requires manual conversion: does not exist in peer-type
	return nil
}

func autoConvert_v1beta1_InitConfiguration_To_upstreamv1beta4_InitConfiguration(in *v1beta1.InitConfiguration, out *InitConfiguration, s apimachineryconversion.Scope) error {
	out.BootstrapTokens = *(*[]BootstrapToken)(unsafe.Pointer(&in.BootstrapTokens))
	if err := Convert_v1beta1_NodeRegistrationOptions_To_upstreamv1beta4_NodeRegistrationOptions(&in.NodeRegistration, &out.NodeRegistration, s); err != nil {
		return err
	}
	if err := Convert_v1beta1_APIEndpoint_To_upstreamv1beta4_APIEndpoint(&in.LocalAPIEndpoint, &out.LocalAPIEndpoint, s); err != nil {
		return err
	}
	out.SkipPhases = *(*[]string)(unsafe.Pointer(&in.SkipPhases))
	out.Patches = (*Patches)(unsafe.Pointer(in.Patches))
	return nil
}

// Convert_v1beta1_InitConfiguration_To_upstreamv1beta4_InitConfiguration is an autogenerated conversion function.
func Convert_v1beta1_InitConfiguration_To_upstreamv1beta4_InitConfiguration(in *v1beta1.InitConfiguration, out *InitConfiguration, s apimachineryconversion.Scope) error {
	return autoConvert_v1beta1_InitConfiguration_To_upstreamv1beta4_InitConfiguration(in, out, s)
}

func autoConvert_upstreamv1beta4_JoinConfiguration_To_v1beta1_JoinConfiguration(in *JoinConfiguration, out *v1beta1.JoinConfiguration, s apimachineryconversion.Scope) error {
	// WARNING: in.DryRun requires manual conversion: does not exist in peer-type
	if err := Convert_upstreamv1beta4_NodeRegistrationOptions_To_v1beta1_NodeRegistrationOptions(&in.NodeRegistration, &out.NodeRegistration, s); err != nil {
		return err
	}
	out.CACertPath = in.CACertPath
	if err := Convert_upstreamv1beta4_Discovery_To_v1beta1_Discovery(&in.Discovery, &out.Discovery, s); err != nil {
		return err
	}
	if in.ControlPlane != nil {
		in, out := &in.ControlPlane, &out.ControlPlane
		*out = new(v1beta1.JoinControlPlane)
		if err := Convert_upstreamv1beta4_JoinControlPlane_To_v1beta1_JoinControlPlane(*in, *out, s); err != nil {
			return err
		}
	} else {
		out.ControlPlane = nil
	}
	out.SkipPhases = *(*[]string)(unsafe.Pointer(&in.SkipPhases))
	out.Patches = (*v1beta1.Patches)(unsafe.Pointer(in.Patches))
	// WARNING: in.Timeouts requires manual conversion: does not exist in peer-type
	return nil
}

func autoConvert_v1beta1_JoinConfiguration_To_upstreamv1beta4_JoinConfiguration(in *v1beta1.JoinConfiguration, out *JoinConfiguration, s apimachineryconversion.Scope) error {
	if err := Convert_v1beta1_NodeRegistrationOptions_To_upstreamv1beta4_NodeRegistrationOptions(&in.NodeRegistration, &out.NodeRegistration, s); err != nil {
		return err
	}
	out.CACertPath = in.CACertPath
	if err := Convert_v1beta1_Discovery_To_upstreamv1beta4_Discovery(&in.Discovery, &out.Discovery, s); err != nil {
		return err
	}
	if in.ControlPlane != nil {
		in, out := &in.ControlPlane, &out.ControlPlane
		*out = new(JoinControlPlane)
		if err := Convert_v1beta1_JoinControlPlane_To_upstreamv1beta4_JoinControlPlane(*in, *out, s); err != nil {
			return err
		}
	} else {
		out.ControlPlane = nil
	}
	out.SkipPhases = *(*[]string)(unsafe.Pointer(&in.SkipPhases))
	out.Patches = (*Patches)(unsafe.Pointer(in.Patches))
	return nil
}

func autoConvert_upstreamv1beta4_JoinControlPlane_To_v1beta1_JoinControlPlane(in *JoinControlPlane, out *v1beta1.JoinControlPlane, s apimachineryconversion.Scope) error {
	if err := Convert_upstreamv1beta4_APIEndpoint_To_v1beta1_APIEndpoint(&in.LocalAPIEndpoint, &out.LocalAPIEndpoint, s); err != nil {
		return err
	}
	// WARNING: in.CertificateKey requires manual conversion: does not exist in peer-type
	return nil
}

func autoConvert_v1beta1_JoinControlPlane_To_upstreamv1beta4_JoinControlPlane(in *v1beta1.JoinControlPlane, out *JoinControlPlane, s apimachineryconversion.Scope) error {
	if err := Convert_v1beta1_APIEndpoint_To_upstreamv1beta4_APIEndpoint(&in.LocalAPIEndpoint, &out.LocalAPIEndpoint, s); err != nil {
		return err
	}
	return nil
}

// Convert_v1beta1_JoinControlPlane_To_upstreamv1beta4_JoinControlPlane is an autogenerated conversion function.
func Convert_v1beta1_JoinControlPlane_To_upstreamv1beta4_JoinControlPlane(in *v1beta1.JoinControlPlane, out *JoinControlPlane, s apimachineryconversion.Scope) error {
	return autoConvert_v1beta1_JoinControlPlane_To_upstreamv1beta4_JoinControlPlane(in, out, s)
}

func autoConvert_upstreamv1beta4_LocalEtcd_To_v1beta1_LocalEtcd(in *LocalEtcd, out *v1beta1.LocalEtcd, s apimachineryconversion.Scope) error {
	if err := Convert_upstreamv1beta4_ImageMeta_To_v1beta1_ImageMeta(&in.ImageMeta, &out.ImageMeta, s); err != nil {
		return err
	}
	out.DataDir = in.DataDir
	// WARNING: in.ExtraArgs requires manual conversion: inconvertible types ([]sigs.k8s.io/cluster-api/bootstrap/kubeadm/types/upstreamv1beta4.Arg vs map[string]string)
	//out.ExtraEnvs = *(*[]v1beta1.EnvVar)(unsafe.Pointer(&in.ExtraEnvs))
	out.ServerCertSANs = *(*[]string)(unsafe.Pointer(&in.ServerCertSANs))
	out.PeerCertSANs = *(*[]string)(unsafe.Pointer(&in.PeerCertSANs))
	return nil
}

func autoConvert_v1beta1_LocalEtcd_To_upstreamv1beta4_LocalEtcd(in *v1beta1.LocalEtcd, out *LocalEtcd, s apimachineryconversion.Scope) error {
	if err := Convert_v1beta1_ImageMeta_To_upstreamv1beta4_ImageMeta(&in.ImageMeta, &out.ImageMeta, s); err != nil {
		return err
	}
	out.DataDir = in.DataDir
	// WARNING: in.ExtraArgs requires manual conversion: inconvertible types (map[string]string vs []sigs.k8s.io/cluster-api/bootstrap/kubeadm/types/upstreamv1beta4.Arg)
	//out.ExtraEnvs = *(*[]EnvVar)(unsafe.Pointer(&in.ExtraEnvs))
	out.ServerCertSANs = *(*[]string)(unsafe.Pointer(&in.ServerCertSANs))
	out.PeerCertSANs = *(*[]string)(unsafe.Pointer(&in.PeerCertSANs))
	return nil
}

func autoConvert_upstreamv1beta4_Networking_To_v1beta1_Networking(in *Networking, out *v1beta1.Networking, s apimachineryconversion.Scope) error {
	out.ServiceSubnet = in.ServiceSubnet
	out.PodSubnet = in.PodSubnet
	out.DNSDomain = in.DNSDomain
	return nil
}

// Convert_upstreamv1beta4_Networking_To_v1beta1_Networking is an autogenerated conversion function.
func Convert_upstreamv1beta4_Networking_To_v1beta1_Networking(in *Networking, out *v1beta1.Networking, s apimachineryconversion.Scope) error {
	return autoConvert_upstreamv1beta4_Networking_To_v1beta1_Networking(in, out, s)
}

func autoConvert_v1beta1_Networking_To_upstreamv1beta4_Networking(in *v1beta1.Networking, out *Networking, s apimachineryconversion.Scope) error {
	out.ServiceSubnet = in.ServiceSubnet
	out.PodSubnet = in.PodSubnet
	out.DNSDomain = in.DNSDomain
	return nil
}

// Convert_v1beta1_Networking_To_upstreamv1beta4_Networking is an autogenerated conversion function.
func Convert_v1beta1_Networking_To_upstreamv1beta4_Networking(in *v1beta1.Networking, out *Networking, s apimachineryconversion.Scope) error {
	return autoConvert_v1beta1_Networking_To_upstreamv1beta4_Networking(in, out, s)
}

func autoConvert_upstreamv1beta4_NodeRegistrationOptions_To_v1beta1_NodeRegistrationOptions(in *NodeRegistrationOptions, out *v1beta1.NodeRegistrationOptions, s apimachineryconversion.Scope) error {
	out.Name = in.Name
	out.CRISocket = in.CRISocket
	out.Taints = *(*[]corev1.Taint)(unsafe.Pointer(&in.Taints))
	// WARNING: in.KubeletExtraArgs requires manual conversion: inconvertible types ([]sigs.k8s.io/cluster-api/bootstrap/kubeadm/types/upstreamv1beta4.Arg vs map[string]string)
	out.IgnorePreflightErrors = *(*[]string)(unsafe.Pointer(&in.IgnorePreflightErrors))
	//out.ImagePullPolicy = string(in.ImagePullPolicy)
	//out.ImagePullSerial = (*bool)(unsafe.Pointer(in.ImagePullSerial))
	return nil
}

func autoConvert_v1beta1_NodeRegistrationOptions_To_upstreamv1beta4_NodeRegistrationOptions(in *v1beta1.NodeRegistrationOptions, out *NodeRegistrationOptions, s apimachineryconversion.Scope) error {
	out.Name = in.Name
	out.CRISocket = in.CRISocket
	out.Taints = *(*[]corev1.Taint)(unsafe.Pointer(&in.Taints))
	// WARNING: in.KubeletExtraArgs requires manual conversion: inconvertible types (map[string]string vs []sigs.k8s.io/cluster-api/bootstrap/kubeadm/types/upstreamv1beta4.Arg)
	out.IgnorePreflightErrors = *(*[]string)(unsafe.Pointer(&in.IgnorePreflightErrors))
	//out.ImagePullPolicy = corev1.PullPolicy(in.ImagePullPolicy)
	//out.ImagePullSerial = (*bool)(unsafe.Pointer(in.ImagePullSerial))
	return nil
}

func autoConvert_upstreamv1beta4_Patches_To_v1beta1_Patches(in *Patches, out *v1beta1.Patches, s apimachineryconversion.Scope) error {
	out.Directory = in.Directory
	return nil
}

// Convert_upstreamv1beta4_Patches_To_v1beta1_Patches is an autogenerated conversion function.
func Convert_upstreamv1beta4_Patches_To_v1beta1_Patches(in *Patches, out *v1beta1.Patches, s apimachineryconversion.Scope) error {
	return autoConvert_upstreamv1beta4_Patches_To_v1beta1_Patches(in, out, s)
}

func autoConvert_v1beta1_Patches_To_upstreamv1beta4_Patches(in *v1beta1.Patches, out *Patches, s apimachineryconversion.Scope) error {
	out.Directory = in.Directory
	return nil
}

// Convert_v1beta1_Patches_To_upstreamv1beta4_Patches is an autogenerated conversion function.
func Convert_v1beta1_Patches_To_upstreamv1beta4_Patches(in *v1beta1.Patches, out *Patches, s apimachineryconversion.Scope) error {
	return autoConvert_v1beta1_Patches_To_upstreamv1beta4_Patches(in, out, s)
}

func Convert_v1beta1_FileDiscovery_To_upstreamv1beta4_FileDiscovery(in *v1beta1.FileDiscovery, out *FileDiscovery, s apimachineryconversion.Scope) error {
	// JoinConfiguration.Discovery.File.KubeConfig does not exist in kubeadm because it's internal to Cluster API, dropping those info.
	return autoConvert_v1beta1_FileDiscovery_To_upstreamv1beta4_FileDiscovery(in, out, s)
}
