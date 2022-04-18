//go:build !ignore_autogenerated
// +build !ignore_autogenerated

/*
Copyright The Kubernetes Authors.

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

// Code generated by controller-gen. DO NOT EDIT.

package v1alpha4

import (
	runtime "k8s.io/apimachinery/pkg/runtime"
	cluster_apiapiv1alpha4 "sigs.k8s.io/cluster-api/api/v1alpha4"
	apiv1alpha4 "sigs.k8s.io/cluster-api/test/infrastructure/docker/api/v1alpha4"
)

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *DockerMachinePool) DeepCopyInto(out *DockerMachinePool) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	in.Spec.DeepCopyInto(&out.Spec)
	in.Status.DeepCopyInto(&out.Status)
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new DockerMachinePool.
func (in *DockerMachinePool) DeepCopy() *DockerMachinePool {
	if in == nil {
		return nil
	}
	out := new(DockerMachinePool)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *DockerMachinePool) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *DockerMachinePoolInstanceStatus) DeepCopyInto(out *DockerMachinePoolInstanceStatus) {
	*out = *in
	if in.Addresses != nil {
		in, out := &in.Addresses, &out.Addresses
		*out = make([]cluster_apiapiv1alpha4.MachineAddress, len(*in))
		copy(*out, *in)
	}
	if in.ProviderID != nil {
		in, out := &in.ProviderID, &out.ProviderID
		*out = new(string)
		**out = **in
	}
	if in.Version != nil {
		in, out := &in.Version, &out.Version
		*out = new(string)
		**out = **in
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new DockerMachinePoolInstanceStatus.
func (in *DockerMachinePoolInstanceStatus) DeepCopy() *DockerMachinePoolInstanceStatus {
	if in == nil {
		return nil
	}
	out := new(DockerMachinePoolInstanceStatus)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *DockerMachinePoolList) DeepCopyInto(out *DockerMachinePoolList) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ListMeta.DeepCopyInto(&out.ListMeta)
	if in.Items != nil {
		in, out := &in.Items, &out.Items
		*out = make([]DockerMachinePool, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new DockerMachinePoolList.
func (in *DockerMachinePoolList) DeepCopy() *DockerMachinePoolList {
	if in == nil {
		return nil
	}
	out := new(DockerMachinePoolList)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *DockerMachinePoolList) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *DockerMachinePoolMachineTemplate) DeepCopyInto(out *DockerMachinePoolMachineTemplate) {
	*out = *in
	if in.EnvironmentVariables != nil {
		in, out := &in.EnvironmentVariables, &out.EnvironmentVariables
		*out = make(map[string]string, len(*in))
		for key, val := range *in {
			(*out)[key] = val
		}
	}
	if in.PreLoadImages != nil {
		in, out := &in.PreLoadImages, &out.PreLoadImages
		*out = make([]string, len(*in))
		copy(*out, *in)
	}
	if in.ExtraMounts != nil {
		in, out := &in.ExtraMounts, &out.ExtraMounts
		*out = make([]apiv1alpha4.Mount, len(*in))
		copy(*out, *in)
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new DockerMachinePoolMachineTemplate.
func (in *DockerMachinePoolMachineTemplate) DeepCopy() *DockerMachinePoolMachineTemplate {
	if in == nil {
		return nil
	}
	out := new(DockerMachinePoolMachineTemplate)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *DockerMachinePoolSpec) DeepCopyInto(out *DockerMachinePoolSpec) {
	*out = *in
	in.Template.DeepCopyInto(&out.Template)
	if in.ProviderIDList != nil {
		in, out := &in.ProviderIDList, &out.ProviderIDList
		*out = make([]string, len(*in))
		copy(*out, *in)
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new DockerMachinePoolSpec.
func (in *DockerMachinePoolSpec) DeepCopy() *DockerMachinePoolSpec {
	if in == nil {
		return nil
	}
	out := new(DockerMachinePoolSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *DockerMachinePoolStatus) DeepCopyInto(out *DockerMachinePoolStatus) {
	*out = *in
	if in.Instances != nil {
		in, out := &in.Instances, &out.Instances
		*out = make([]DockerMachinePoolInstanceStatus, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
	if in.Conditions != nil {
		in, out := &in.Conditions, &out.Conditions
		*out = make(cluster_apiapiv1alpha4.Conditions, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new DockerMachinePoolStatus.
func (in *DockerMachinePoolStatus) DeepCopy() *DockerMachinePoolStatus {
	if in == nil {
		return nil
	}
	out := new(DockerMachinePoolStatus)
	in.DeepCopyInto(out)
	return out
}
