/*
Copyright 2021 The Kubernetes Authors.

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

package utils

import (
	"testing"
	"time"

	"github.com/blang/semver"
	"github.com/google/go-cmp/cmp"
	. "github.com/onsi/gomega"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"

	bootstrapv1 "sigs.k8s.io/cluster-api/bootstrap/kubeadm/api/v1beta1"
	"sigs.k8s.io/cluster-api/bootstrap/kubeadm/types/upstreamv1beta2"
	"sigs.k8s.io/cluster-api/bootstrap/kubeadm/types/upstreamv1beta3"
	"sigs.k8s.io/cluster-api/bootstrap/kubeadm/types/upstreamv1beta4"
)

func TestKubeVersionToKubeadmAPIGroupVersion(t *testing.T) {
	type args struct {
		version semver.Version
	}
	tests := []struct {
		name    string
		args    args
		want    schema.GroupVersion
		wantErr bool
	}{
		{
			name: "fails when kubernetes version is too old",
			args: args{
				version: semver.MustParse("1.12.0"),
			},
			want:    schema.GroupVersion{},
			wantErr: true,
		},
		{
			name: "fails with kubernetes version for kubeadm API v1beta1",
			args: args{
				version: semver.MustParse("1.14.99"),
			},
			wantErr: true,
		},
		{
			name: "pass with minimum kubernetes alpha version for kubeadm API v1beta2",
			args: args{
				version: semver.MustParse("1.15.0-alpha.0.734+ba502ee555924a"),
			},
			want:    upstreamv1beta2.GroupVersion,
			wantErr: false,
		},
		{
			name: "pass with minimum kubernetes version for kubeadm API v1beta2",
			args: args{
				version: semver.MustParse("1.15.0"),
			},
			want:    upstreamv1beta2.GroupVersion,
			wantErr: false,
		},
		{
			name: "pass with kubernetes version for kubeadm API v1beta2",
			args: args{
				version: semver.MustParse("1.20.99"),
			},
			want:    upstreamv1beta2.GroupVersion,
			wantErr: false,
		},
		{
			name: "pass with minimum kubernetes alpha version for kubeadm API v1beta3",
			args: args{
				version: semver.MustParse("1.22.0-alpha.0.734+ba502ee555924a"),
			},
			want:    upstreamv1beta3.GroupVersion,
			wantErr: false,
		},
		{
			name: "pass with minimum kubernetes version for kubeadm API v1beta3",
			args: args{
				version: semver.MustParse("1.22.0"),
			},
			want:    upstreamv1beta3.GroupVersion,
			wantErr: false,
		},
		{
			name: "pass with kubernetes version for kubeadm API v1beta3",
			args: args{
				version: semver.MustParse("1.23.99"),
			},
			want:    upstreamv1beta3.GroupVersion,
			wantErr: false,
		},
		{
			name: "pass with minimum kubernetes version for kubeadm API v1beta4",
			args: args{
				version: semver.MustParse("1.31.0"),
			},
			want:    upstreamv1beta4.GroupVersion,
			wantErr: false,
		},
		{
			name: "pass with kubernetes version for kubeadm API v1beta4",
			args: args{
				version: semver.MustParse("1.32.99"),
			},
			want:    upstreamv1beta4.GroupVersion,
			wantErr: false,
		},
		{
			name: "pass with future kubernetes version",
			args: args{
				version: semver.MustParse("99.99.99"),
			},
			want:    upstreamv1beta4.GroupVersion,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := NewWithT(t)

			got, err := KubeVersionToKubeadmAPIGroupVersion(tt.args.version)
			if tt.wantErr {
				g.Expect(err).To(HaveOccurred())
				return
			}
			g.Expect(err).ToNot(HaveOccurred())
			g.Expect(got).To(BeComparableTo(tt.want))
		})
	}
}

func TestMarshalClusterConfigurationForVersion(t *testing.T) {
	type args struct {
		capiObj *bootstrapv1.ClusterConfiguration
		version semver.Version
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "Generates a v1beta2 kubeadm configuration",
			args: args{
				capiObj: &bootstrapv1.ClusterConfiguration{},
				version: semver.MustParse("1.15.0"),
			},
			want: "apiServer: {}\n" +
				"apiVersion: kubeadm.k8s.io/v1beta2\n" +
				"controllerManager: {}\n" +
				"dns: {}\n" +
				"etcd: {}\n" +
				"kind: ClusterConfiguration\n" +
				"networking: {}\n" +
				"scheduler: {}\n",
			wantErr: false,
		},
		{
			name: "Generates a v1beta3 kubeadm configuration",
			args: args{
				capiObj: &bootstrapv1.ClusterConfiguration{},
				version: semver.MustParse("1.22.0"),
			},
			want: "apiServer: {}\n" +
				"apiVersion: kubeadm.k8s.io/v1beta3\n" +
				"controllerManager: {}\n" +
				"dns: {}\n" +
				"etcd: {}\n" +
				"kind: ClusterConfiguration\n" +
				"networking: {}\n" +
				"scheduler: {}\n",
			wantErr: false,
		},
		{
			name: "Generates a v1beta4 kubeadm configuration",
			args: args{
				capiObj: &bootstrapv1.ClusterConfiguration{},
				version: semver.MustParse("1.31.0"),
			},
			want: "apiServer: {}\n" +
				"apiVersion: kubeadm.k8s.io/v1beta4\n" +
				"controllerManager: {}\n" +
				"dns: {}\n" +
				"etcd: {}\n" +
				"kind: ClusterConfiguration\n" +
				"networking: {}\n" +
				"proxy: {}\n" +
				"scheduler: {}\n",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := NewWithT(t)

			got, err := MarshalClusterConfigurationForVersion(tt.args.capiObj, tt.args.version)
			if tt.wantErr {
				g.Expect(err).To(HaveOccurred())
				return
			}
			g.Expect(err).ToNot(HaveOccurred())
			g.Expect(got).To(Equal(tt.want), cmp.Diff(tt.want, got))
		})
	}
}

func TestMarshalClusterStatusForVersion(t *testing.T) {
	type args struct {
		capiObj *bootstrapv1.ClusterStatus
		version semver.Version
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "Generates a v1beta2 kubeadm status",
			args: args{
				capiObj: &bootstrapv1.ClusterStatus{},
				version: semver.MustParse("1.15.0"),
			},
			want: "apiEndpoints: null\n" +
				"apiVersion: kubeadm.k8s.io/v1beta2\n" + "" +
				"kind: ClusterStatus\n",
			wantErr: false,
		},
		{
			name: "Fails generating a v1beta3 kubeadm status",
			args: args{
				capiObj: &bootstrapv1.ClusterStatus{},
				version: semver.MustParse("1.22.0"),
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := NewWithT(t)

			got, err := MarshalClusterStatusForVersion(tt.args.capiObj, tt.args.version)
			if tt.wantErr {
				g.Expect(err).To(HaveOccurred())
				return
			}
			g.Expect(err).ToNot(HaveOccurred())
			g.Expect(got).To(Equal(tt.want), cmp.Diff(tt.want, got))
		})
	}
}

func TestMarshalInitConfigurationForVersion(t *testing.T) {
	timeout := metav1.Duration{Duration: 10 * time.Second}
	type args struct {
		clusterConfiguration *bootstrapv1.ClusterConfiguration
		initConfiguration    *bootstrapv1.InitConfiguration
		version              semver.Version
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "Generates a v1beta2 kubeadm init configuration",
			args: args{
				initConfiguration: &bootstrapv1.InitConfiguration{
					NodeRegistration: bootstrapv1.NodeRegistrationOptions{
						IgnorePreflightErrors: []string{"some-preflight-check"},
					},
				},
				version: semver.MustParse("1.15.0"),
			},
			want: "apiVersion: kubeadm.k8s.io/v1beta2\n" +
				"kind: InitConfiguration\n" +
				"localAPIEndpoint: {}\n" +
				"nodeRegistration:\n" +
				"  ignorePreflightErrors:\n" +
				"  - some-preflight-check\n",
			wantErr: false,
		},
		{
			name: "Generates a v1beta3 kubeadm init configuration",
			args: args{
				initConfiguration: &bootstrapv1.InitConfiguration{
					NodeRegistration: bootstrapv1.NodeRegistrationOptions{
						IgnorePreflightErrors: []string{"some-preflight-check"},
					},
				},
				version: semver.MustParse("1.22.0"),
			},
			want: "apiVersion: kubeadm.k8s.io/v1beta3\n" +
				"kind: InitConfiguration\n" +
				"localAPIEndpoint: {}\n" +
				"nodeRegistration:\n" +
				"  ignorePreflightErrors:\n" +
				"  - some-preflight-check\n" +
				"  taints: null\n",
			wantErr: false,
		},
		{
			name: "Generates a v1beta4 kubeadm init configuration",
			args: args{
				initConfiguration: &bootstrapv1.InitConfiguration{
					NodeRegistration: bootstrapv1.NodeRegistrationOptions{
						IgnorePreflightErrors: []string{"some-preflight-check"},
					},
				},
				version: semver.MustParse("1.31.0"),
			},
			want: "apiVersion: kubeadm.k8s.io/v1beta4\n" +
				"kind: InitConfiguration\n" +
				"localAPIEndpoint: {}\n" +
				"nodeRegistration:\n" +
				"  ignorePreflightErrors:\n" +
				"  - some-preflight-check\n" +
				"  taints: null\n",
			wantErr: false,
		},
		{
			name: "Generates a v1beta4 kubeadm init configuration with data from cluster configuration",
			args: args{
				clusterConfiguration: &bootstrapv1.ClusterConfiguration{
					APIServer: bootstrapv1.APIServer{TimeoutForControlPlane: &timeout},
				},
				initConfiguration: &bootstrapv1.InitConfiguration{
					NodeRegistration: bootstrapv1.NodeRegistrationOptions{
						IgnorePreflightErrors: []string{"some-preflight-check"},
					},
				},
				version: semver.MustParse("1.31.0"),
			},
			want: "apiVersion: kubeadm.k8s.io/v1beta4\n" +
				"kind: InitConfiguration\n" +
				"localAPIEndpoint: {}\n" +
				"nodeRegistration:\n" +
				"  ignorePreflightErrors:\n" +
				"  - some-preflight-check\n" +
				"  taints: null\n" +
				"timeouts:\n" +
				"  controlPlaneComponentHealthCheck: 10s\n",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := NewWithT(t)

			got, err := MarshalInitConfigurationForVersion(tt.args.clusterConfiguration, tt.args.initConfiguration, tt.args.version)
			if tt.wantErr {
				g.Expect(err).To(HaveOccurred())
				return
			}
			g.Expect(err).ToNot(HaveOccurred())
			g.Expect(got).To(Equal(tt.want), cmp.Diff(tt.want, got))
		})
	}
}

func TestMarshalJoinConfigurationForVersion(t *testing.T) {
	timeout := metav1.Duration{Duration: 10 * time.Second}
	type args struct {
		clusterConfiguration *bootstrapv1.ClusterConfiguration
		joinConfiguration    *bootstrapv1.JoinConfiguration
		version              semver.Version
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "Generates a v1beta2 kubeadm join configuration",
			args: args{
				joinConfiguration: &bootstrapv1.JoinConfiguration{
					NodeRegistration: bootstrapv1.NodeRegistrationOptions{
						IgnorePreflightErrors: []string{"some-preflight-check"},
					},
				},
				version: semver.MustParse("1.15.0"),
			},
			want: "apiVersion: kubeadm.k8s.io/v1beta2\n" + "" +
				"discovery: {}\n" +
				"kind: JoinConfiguration\n" +
				"nodeRegistration:\n" +
				"  ignorePreflightErrors:\n" +
				"  - some-preflight-check\n",
			wantErr: false,
		},
		{
			name: "Generates a v1beta3 kubeadm join configuration",
			args: args{
				joinConfiguration: &bootstrapv1.JoinConfiguration{
					NodeRegistration: bootstrapv1.NodeRegistrationOptions{
						IgnorePreflightErrors: []string{"some-preflight-check"},
					},
				},
				version: semver.MustParse("1.22.0"),
			},
			want: "apiVersion: kubeadm.k8s.io/v1beta3\n" + "" +
				"discovery: {}\n" +
				"kind: JoinConfiguration\n" +
				"nodeRegistration:\n" +
				"  ignorePreflightErrors:\n" +
				"  - some-preflight-check\n" +
				"  taints: null\n",
			wantErr: false,
		},
		{
			name: "Generates a v1beta4 kubeadm join configuration",
			args: args{
				joinConfiguration: &bootstrapv1.JoinConfiguration{
					NodeRegistration: bootstrapv1.NodeRegistrationOptions{
						IgnorePreflightErrors: []string{"some-preflight-check"},
					},
				},
				version: semver.MustParse("1.31.0"),
			},
			want: "apiVersion: kubeadm.k8s.io/v1beta4\n" + "" +
				"discovery: {}\n" +
				"kind: JoinConfiguration\n" +
				"nodeRegistration:\n" +
				"  ignorePreflightErrors:\n" +
				"  - some-preflight-check\n" +
				"  taints: null\n",
			wantErr: false,
		},
		{
			name: "Generates a v1beta4 kubeadm join configuration with data from cluster configuration",
			args: args{
				clusterConfiguration: &bootstrapv1.ClusterConfiguration{
					APIServer: bootstrapv1.APIServer{TimeoutForControlPlane: &timeout},
				},
				joinConfiguration: &bootstrapv1.JoinConfiguration{
					NodeRegistration: bootstrapv1.NodeRegistrationOptions{
						IgnorePreflightErrors: []string{"some-preflight-check"},
					},
				},
				version: semver.MustParse("1.31.0"),
			},
			want: "apiVersion: kubeadm.k8s.io/v1beta4\n" + "" +
				"discovery: {}\n" +
				"kind: JoinConfiguration\n" +
				"nodeRegistration:\n" +
				"  ignorePreflightErrors:\n" +
				"  - some-preflight-check\n" +
				"  taints: null\n" +
				"timeouts:\n" +
				"  controlPlaneComponentHealthCheck: 10s\n",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := NewWithT(t)

			got, err := MarshalJoinConfigurationForVersion(tt.args.clusterConfiguration, tt.args.joinConfiguration, tt.args.version)
			if tt.wantErr {
				g.Expect(err).To(HaveOccurred())
				return
			}
			g.Expect(err).ToNot(HaveOccurred())
			g.Expect(got).To(Equal(tt.want), cmp.Diff(tt.want, got))
		})
	}
}

func TestUnmarshalClusterConfiguration(t *testing.T) {
	type args struct {
		yaml string
	}
	tests := []struct {
		name    string
		args    args
		want    *bootstrapv1.ClusterConfiguration
		wantErr bool
	}{
		{
			name: "Parses a v1beta2 kubeadm configuration",
			args: args{
				yaml: "apiServer: {}\n" +
					"apiVersion: kubeadm.k8s.io/v1beta2\n" + "" +
					"controllerManager: {}\n" +
					"dns: {}\n" +
					"etcd: {}\n" +
					"kind: ClusterConfiguration\n" +
					"networking: {}\n" +
					"scheduler: {}\n",
			},
			want:    &bootstrapv1.ClusterConfiguration{},
			wantErr: false,
		},
		{
			name: "Parses a v1beta3 kubeadm configuration",
			args: args{
				yaml: "apiServer: {}\n" +
					"apiVersion: kubeadm.k8s.io/v1beta3\n" + "" +
					"controllerManager: {}\n" +
					"dns: {}\n" +
					"etcd: {}\n" +
					"kind: ClusterConfiguration\n" +
					"networking: {}\n" +
					"scheduler: {}\n",
			},
			want:    &bootstrapv1.ClusterConfiguration{},
			wantErr: false,
		},
		{
			name: "Parses a v1beta4 kubeadm configuration",
			args: args{
				yaml: "    apiVersion: kubeadm.k8s.io/v1beta4\n    kind: ClusterConfiguration\n    apiServer:\n      extraArgs:\n      - name: admission-control-config-file\n        value: /etc/kubernetes/pod-security-standard.yaml\n      - name: anonymous-auth\n        value: \"true\"\n      - name: audit-log-maxage\n        value: \"30\"\n      - name: audit-log-maxbackup\n        value: \"10\"\n      - name: audit-log-maxsize\n        value: \"100\"\n      - name: audit-log-path\n        value: /var/log/apiserver/audit.log\n      - name: audit-policy-file\n        value: /etc/kubernetes/audit-policy.yaml\n      - name: authorization-mode\n        value: RBAC,Node\n      - name: cloud-provider\n        value: external\n      - name: default-not-ready-toleration-seconds\n        value: \"60\"\n      - name: default-unreachable-toleration-seconds\n        value: \"60\"\n      - name: disable-admission-plugins\n        value: AlwaysAdmit\n      - name: enable-admission-plugins\n        value: AlwaysPullImages,NamespaceLifecycle,ServiceAccount,NodeRestriction,PodSecurity\n      - name: kubelet-certificate-authority\n        value: /etc/kubernetes/pki/ca.crt\n      - name: profiling\n        value: \"false\"\n      - name: secure-port\n        value: \"6443\"\n      - name: tls-cipher-suites\n        value: TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256,TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256,TLS_ECDHE_ECDSA_WITH_CHACHA20_POLY1305,TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,TLS_ECDHE_RSA_WITH_CHACHA20_POLY1305,TLS_ECDHE_ECDSA_WITH_AES_256_GCM_SHA384,TLS_RSA_WITH_AES_256_GCM_SHA384,TLS_RSA_WITH_AES_128_GCM_SHA256\n      extraVolumes:\n      - hostPath: /var/log/apiserver\n        mountPath: /var/log/apiserver\n        name: audit-log\n        pathType: DirectoryOrCreate\n      - hostPath: /etc/kubernetes/audit-policy.yaml\n        mountPath: /etc/kubernetes/audit-policy.yaml\n        name: audit-policy\n        pathType: File\n        readOnly: true\n      - hostPath: /etc/kubernetes/pod-security-standard.yaml\n        mountPath: /etc/kubernetes/pod-security-standard.yaml\n        name: pod-security-standard\n        pathType: File\n        readOnly: true\n    caCertificateValidityPeriod: 87600h0m0s\n    certificateValidityPeriod: 8760h0m0s\n    certificatesDir: /etc/kubernetes/pki\n    clusterName: am-kcp-cluster2\n    controlPlaneEndpoint: am-kcp-cluster2-or8i.spectrocloud.dev:6443\n    controllerManager:\n      extraArgs:\n      - name: cloud-provider\n        value: external\n      - name: feature-gates\n        value: RotateKubeletServerCertificate=true\n      - name: profiling\n        value: \"false\"\n      - name: terminated-pod-gc-threshold\n        value: \"25\"\n      - name: use-service-account-credentials\n        value: \"true\"\n    encryptionAlgorithm: RSA-2048\n    etcd:\n      local:\n        dataDir: /var/lib/etcd\n    featureGates:\n      ControlPlaneKubeletLocalMode: true\n    imageRepository: registry.k8s.io\n    kubernetesVersion: v1.31.4\n    networking:\n      dnsDomain: cluster.local\n      podSubnet: 192.168.0.0/16\n      serviceSubnet: 10.96.0.0/12\n    scheduler:\n      extraArgs:\n      - name: profiling\n        value: \"false\"",
			},
			want:    &bootstrapv1.ClusterConfiguration{},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := NewWithT(t)

			_, err := UnmarshalClusterConfiguration(tt.args.yaml)
			if tt.wantErr {
				g.Expect(err).To(HaveOccurred())
				return
			}
			g.Expect(err).ToNot(HaveOccurred())
			//g.Expect(got).To(BeComparableTo(tt.want), cmp.Diff(tt.want, got))
		})
	}
}

func TestUnmarshalClusterStatus(t *testing.T) {
	type args struct {
		yaml string
	}
	tests := []struct {
		name    string
		args    args
		want    *bootstrapv1.ClusterStatus
		wantErr bool
	}{
		{
			name: "Parses a v1beta2 kubeadm configuration",
			args: args{
				yaml: "apiEndpoints: null\n" +
					"apiVersion: kubeadm.k8s.io/v1beta2\n" + "" +
					"kind: ClusterStatus\n",
			},
			want:    &bootstrapv1.ClusterStatus{},
			wantErr: false,
		},
		{
			name: "Fails parsing a v1beta3 kubeadm configuration",
			args: args{
				yaml: "apiEndpoints: null\n" +
					"apiVersion: kubeadm.k8s.io/v1beta3\n" + "" +
					"kind: ClusterStatus\n",
			},
			wantErr: true,
		},
		{
			name: "Fails parsing a v1beta4 kubeadm configuration",
			args: args{
				yaml: "apiEndpoints: null\n" +
					"apiVersion: kubeadm.k8s.io/v1beta4\n" + "" +
					"kind: ClusterStatus\n",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := NewWithT(t)

			got, err := UnmarshalClusterStatus(tt.args.yaml)
			if tt.wantErr {
				g.Expect(err).To(HaveOccurred())
				return
			}
			g.Expect(err).ToNot(HaveOccurred())
			g.Expect(got).To(BeComparableTo(tt.want), cmp.Diff(tt.want, got))
		})
	}
}

func TestUnmarshalInitConfiguration(t *testing.T) {
	timeout := metav1.Duration{Duration: 10 * time.Second}
	type args struct {
		yaml string
	}
	tests := []struct {
		name                     string
		args                     args
		want                     *bootstrapv1.InitConfiguration
		wantClusterConfiguration *bootstrapv1.ClusterConfiguration
		wantErr                  bool
	}{
		{
			name: "Parses a v1beta2 kubeadm configuration",
			args: args{
				yaml: "apiVersion: kubeadm.k8s.io/v1beta2\n" + "" +
					"kind: InitConfiguration\n" +
					"localAPIEndpoint: {}\n" +
					"nodeRegistration: {}\n",
			},
			want:                     &bootstrapv1.InitConfiguration{},
			wantClusterConfiguration: &bootstrapv1.ClusterConfiguration{},
			wantErr:                  false,
		},
		{
			name: "Parses a v1beta3 kubeadm configuration",
			args: args{
				yaml: "apiVersion: kubeadm.k8s.io/v1beta3\n" + "" +
					"kind: InitConfiguration\n" +
					"localAPIEndpoint: {}\n" +
					"nodeRegistration: {}\n",
			},
			want:                     &bootstrapv1.InitConfiguration{},
			wantClusterConfiguration: &bootstrapv1.ClusterConfiguration{},
			wantErr:                  false,
		},
		{
			name: "Parses a v1beta4 kubeadm configuration",
			args: args{
				yaml: "apiVersion: kubeadm.k8s.io/v1beta4\n" + "" +
					"kind: InitConfiguration\n" +
					"localAPIEndpoint: {}\n" +
					"nodeRegistration: {}\n",
			},
			want:                     &bootstrapv1.InitConfiguration{},
			wantClusterConfiguration: &bootstrapv1.ClusterConfiguration{},
			wantErr:                  false,
		},
		{
			name: "Parses a v1beta4 kubeadm configuration with data to cluster configuration",
			args: args{
				yaml: "apiVersion: kubeadm.k8s.io/v1beta4\n" + "" +
					"kind: InitConfiguration\n" +
					"localAPIEndpoint: {}\n" +
					"nodeRegistration: {}\n" +
					"timeouts:\n" +
					"  controlPlaneComponentHealthCheck: 10s\n",
			},
			want: &bootstrapv1.InitConfiguration{},
			wantClusterConfiguration: &bootstrapv1.ClusterConfiguration{
				APIServer: bootstrapv1.APIServer{TimeoutForControlPlane: &timeout},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := NewWithT(t)

			gotClusterConfiguration := &bootstrapv1.ClusterConfiguration{}
			got, err := UnmarshalInitConfiguration(tt.args.yaml, gotClusterConfiguration)
			if tt.wantErr {
				g.Expect(err).To(HaveOccurred())
				return
			}
			g.Expect(err).ToNot(HaveOccurred())
			g.Expect(got).To(BeComparableTo(tt.want), cmp.Diff(tt.want, got))
			g.Expect(gotClusterConfiguration).To(BeComparableTo(tt.wantClusterConfiguration), cmp.Diff(tt.wantClusterConfiguration, gotClusterConfiguration))
		})
	}
}

func TestUnmarshalJoinConfiguration(t *testing.T) {
	timeout := metav1.Duration{Duration: 10 * time.Second}
	type args struct {
		yaml string
	}
	tests := []struct {
		name                     string
		args                     args
		want                     *bootstrapv1.JoinConfiguration
		wantClusterConfiguration *bootstrapv1.ClusterConfiguration
		wantErr                  bool
	}{
		{
			name: "Parses a v1beta2 kubeadm configuration",
			args: args{
				yaml: "apiVersion: kubeadm.k8s.io/v1beta2\n" + "" +
					"caCertPath: \"\"\n" +
					"discovery: {}\n" +
					"kind: JoinConfiguration\n",
			},
			want:                     &bootstrapv1.JoinConfiguration{},
			wantClusterConfiguration: &bootstrapv1.ClusterConfiguration{},
			wantErr:                  false,
		},
		{
			name: "Parses a v1beta3 kubeadm configuration",
			args: args{
				yaml: "apiVersion: kubeadm.k8s.io/v1beta3\n" + "" +
					"caCertPath: \"\"\n" +
					"discovery: {}\n" +
					"kind: JoinConfiguration\n",
			},
			want:                     &bootstrapv1.JoinConfiguration{},
			wantClusterConfiguration: &bootstrapv1.ClusterConfiguration{},
			wantErr:                  false,
		},
		{
			name: "Parses a v1beta4 kubeadm configuration",
			args: args{
				yaml: "apiVersion: kubeadm.k8s.io/v1beta4\n" + "" +
					"caCertPath: \"\"\n" +
					"discovery: {}\n" +
					"kind: JoinConfiguration\n",
			},
			want:                     &bootstrapv1.JoinConfiguration{},
			wantClusterConfiguration: &bootstrapv1.ClusterConfiguration{},
			wantErr:                  false,
		},
		{
			name: "Parses a v1beta4 kubeadm configuration with data to cluster configuration",
			args: args{
				yaml: "apiVersion: kubeadm.k8s.io/v1beta4\n" + "" +
					"caCertPath: \"\"\n" +
					"discovery: {}\n" +
					"kind: JoinConfiguration\n" +
					"timeouts:\n" +
					"  controlPlaneComponentHealthCheck: 10s\n",
			},
			want: &bootstrapv1.JoinConfiguration{},
			wantClusterConfiguration: &bootstrapv1.ClusterConfiguration{
				APIServer: bootstrapv1.APIServer{TimeoutForControlPlane: &timeout},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := NewWithT(t)

			gotClusterConfiguration := &bootstrapv1.ClusterConfiguration{}
			got, err := UnmarshalJoinConfiguration(tt.args.yaml, gotClusterConfiguration)
			if tt.wantErr {
				g.Expect(err).To(HaveOccurred())
				return
			}
			g.Expect(err).ToNot(HaveOccurred())
			g.Expect(got).To(BeComparableTo(tt.want), cmp.Diff(tt.want, got))
			g.Expect(gotClusterConfiguration).To(BeComparableTo(tt.wantClusterConfiguration), cmp.Diff(tt.wantClusterConfiguration, gotClusterConfiguration))
		})
	}
}
