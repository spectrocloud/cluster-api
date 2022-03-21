/*
Copyright 2018 The Kubernetes Authors.

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

package docker

import (
	"testing"

	. "github.com/onsi/gomega"
)

func TestMachineToInfrastructureMapFunc(t *testing.T) {
	g := NewWithT(t)

	var testcases = []struct {
		name   string
		input  string
		output bool
	}{
		{
			name: "Should return true as the the string is present",
			input: `[WARNING DirAvailable--etc-kubernetes-manifests[]: /etc/kubernetes/manifests is not empty
			[WARNING FileAvailable--etc-kubernetes-kubelet.conf[]: /etc/kubernetes/kubelet.conf already exists
			[WARNING Swap[]: running with swap on is not supported. Please disable swap
			[WARNING Port-10250[]: Port 10250 is in use
			[WARNING Port-6443[]: Port 6443 is in use
			[WARNING Port-10259[]: Port 10259 is in use
			[WARNING Port-10257[]: Port 10257 is in use
			[WARNING FileAvailable--etc-kubernetes-manifests-kube-apiserver.yaml[]: /etc/kubernetes/manifests/kube-apiserver.yaml already exists
			[WARNING FileAvailable--etc-kubernetes-manifests-kube-controller-manager.yaml[]: /etc/kubernetes/manifests/kube-controller-manager.yaml already exists
			[WARNING FileAvailable--etc-kubernetes-manifests-kube-scheduler.yaml[]: /etc/kubernetes/manifests/kube-scheduler.yaml already exists
			[WARNING FileAvailable--etc-kubernetes-manifests-etcd.yaml[]: /etc/kubernetes/manifests/etcd.yaml already exists
			[WARNING Port-2379[]: Port 2379 is in use
			[WARNING Port-2380[]: Port 2380 is in use
			[WARNING DirAvailable--var-lib-etcd[]: /var/lib/etcd is not empty error execution phase kubelet-start: a 
			Node with name "spectrocluster-sample-cp-tj6sq" and status "Ready" already exists in the cluster. You must delete the existing Node or change the name of this new joining Node
		To see the stack trace of this error execute with --v=5 or higher`,
			output: true,
		},
		{
			name:   "Should return true as the the string is present",
			input:  `Node with name "spectrocluster-sample-cp-tj6sq" and status "Ready" already exists in the cluster. You must delete the existing Node or change the name of this new joining Node\n`,
			output: true,
		},
		{
			name:   "Should return false as the the string is not present",
			input:  `Node with name abcasd \"spectrocluster-sample-cp-tj6sq\" and status \"Ready\" already exists in the cluster. You must delete the existing Node or change the name of this new joining Node\n`,
			output: false,
		},
		{
			name:   "Should return false as the the string is not present",
			input:  `should not match string`,
			output: false,
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			out := nodeAlreadyExists(tc.input)
			g.Expect(out).To(Equal(tc.output))
		})
	}
}
