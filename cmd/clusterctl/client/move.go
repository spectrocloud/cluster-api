/*
Copyright 2020 The Kubernetes Authors.

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

package client

import (
	"context"
	"fmt"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"os"

	"github.com/pkg/errors"

	"sigs.k8s.io/cluster-api/cmd/clusterctl/client/cluster"
)

// MoveOptions carries the options supported by move.
type MoveOptions struct {
	// FromKubeconfig defines the kubeconfig to use for accessing the source management cluster. If empty,
	// default rules for kubeconfig discovery will be used.
	FromKubeconfig Kubeconfig

	// ToKubeconfig defines the kubeconfig to use for accessing the target management cluster. If empty,
	// default rules for kubeconfig discovery will be used.
	ToKubeconfig Kubeconfig

	// Convert to palette CRD input
	ToPaletteCRD string

	// Namespace where the objects describing the workload cluster exists. If unspecified, the current
	// namespace will be used.
	Namespace string

	// ExperimentalResourceMutatorFn accepts any number of resource mutator functions that are applied on all resources being moved.
	// This is an experimental feature and is exposed only from the library and not (yet) through the CLI.
	ExperimentalResourceMutators []cluster.ResourceMutatorFunc

	// FromDirectory apply configuration from directory.
	FromDirectory string

	// ToDirectory save configuration to directory.
	ToDirectory string

	// DryRun means the move action is a dry run, no real action will be performed.
	DryRun bool

	// palette specific options
	IgnoreClusterClass bool
	ClusterName        string
	ToNamespace        string
}

func (c *clusterctlClient) Move(ctx context.Context, options MoveOptions) error {
	// Both backup and restore makes no sense. It's a complete move.
	if options.FromDirectory != "" && options.ToDirectory != "" {
		return errors.Errorf("can't set both FromDirectory and ToDirectory")
	}

	if !options.DryRun &&
		options.FromDirectory == "" &&
		options.ToDirectory == "" &&
		options.ToPaletteCRD == "" &&
		options.ToKubeconfig == (Kubeconfig{}) {
		return errors.Errorf("at least one of FromDirectory, ToDirectory and ToKubeconfig must be set")
	}

	if options.ToPaletteCRD != "" {
		return c.toPaletteCRD(ctx, options)
	}

	if options.ToDirectory != "" {
		return c.toDirectory(ctx, options)
	} else if options.FromDirectory != "" {
		return c.fromDirectory(ctx, options)
	}

	return c.move(ctx, options)
}

func (c *clusterctlClient) move(ctx context.Context, options MoveOptions) error {
	// Get the client for interacting with the source management cluster.
	fromCluster, err := c.getClusterClient(ctx, options.FromKubeconfig)
	if err != nil {
		return err
	}

	// If the option specifying the Namespace is empty, try to detect it.
	if options.Namespace == "" {
		currentNamespace, err := fromCluster.Proxy().CurrentNamespace()
		if err != nil {
			return err
		}
		options.Namespace = currentNamespace
	}

	var toCluster cluster.Client
	if !options.DryRun {
		// Get the client for interacting with the target management cluster.
		if toCluster, err = c.getClusterClient(ctx, options.ToKubeconfig); err != nil {
			return err
		}
	}

	mutators := getPaletteMutators(options.Namespace, options.ClusterName, options.ToNamespace)

	return fromCluster.ObjectMover().Move(ctx, options.Namespace, options.ClusterName, toCluster, options.DryRun, mutators...)
}

func (c *clusterctlClient) fromDirectory(ctx context.Context, options MoveOptions) error {
	toCluster, err := c.getClusterClient(ctx, options.ToKubeconfig)
	if err != nil {
		return err
	}

	if _, err := os.Stat(options.FromDirectory); os.IsNotExist(err) {
		return err
	}

	return toCluster.ObjectMover().FromDirectory(ctx, toCluster, options.FromDirectory)
}

func (c *clusterctlClient) toDirectory(ctx context.Context, options MoveOptions) error {
	fromCluster, err := c.getClusterClient(ctx, options.FromKubeconfig)
	if err != nil {
		return err
	}

	// If the option specifying the Namespace is empty, try to detect it.
	if options.Namespace == "" {
		currentNamespace, err := fromCluster.Proxy().CurrentNamespace()
		if err != nil {
			return err
		}
		options.Namespace = currentNamespace
	}

	if _, err := os.Stat(options.ToDirectory); os.IsNotExist(err) {
		return err
	}

	mutators := getPaletteMutators(options.Namespace, options.ClusterName, options.ToNamespace)

	return fromCluster.ObjectMover().ToDirectory(ctx, options.Namespace, options.ClusterName, options.ToDirectory, mutators...)
}

func (c *clusterctlClient) toPaletteCRD(ctx context.Context, options MoveOptions) error {
	fromCluster, err := c.getClusterClient(ctx, options.FromKubeconfig)
	if err != nil {
		return err
	}

	// If the option specifying the Namespace is empty, try to detect it.
	if options.Namespace == "" {
		currentNamespace, err := fromCluster.Proxy().CurrentNamespace()
		if err != nil {
			return err
		}
		options.Namespace = currentNamespace
	}

	if _, err := os.Stat(options.ToPaletteCRD); os.IsNotExist(err) {
		return err
	}

	mutators := GetClusterTemplateMutator()
	return fromCluster.ObjectMover().ToPaletteCRD(ctx, options.Namespace, options.ClusterName, options.ToPaletteCRD, mutators)
}

func (c *clusterctlClient) getClusterClient(ctx context.Context, kubeconfig Kubeconfig) (cluster.Client, error) {
	cluster, err := c.clusterClientFactory(ClusterClientFactoryInput{Kubeconfig: kubeconfig})
	if err != nil {
		return nil, err
	}

	// Ensure this command only runs against management clusters with the current Cluster API contract.
	if err := cluster.ProviderInventory().CheckCAPIContract(ctx); err != nil {
		return nil, err
	}

	// Ensures the custom resource definitions required by clusterctl are in place.
	if err := cluster.ProviderInventory().EnsureCustomResourceDefinitions(ctx); err != nil {
		return nil, err
	}
	return cluster, nil
}

func getPaletteMutators(currentNamespace, clusterName, targetNamespace string) []cluster.ResourceMutatorFunc {
	fmt.Println("Applying palette namespace mutators")
	return []cluster.ResourceMutatorFunc{getNamespaceMutator(targetNamespace)}
}

func getTemplateMutatorKinds() map[string][][]string {
	kindsToUpdate := map[string][][]string{
		"Cluster": {
			//{"metadata", "annotations", "kubectl.kubernetes.io/last-applied-configuration"},
			//{"metadata", "annotations", "TKGOperationInfo"},
			//{"metadata", "annotations", "TKGOperationLastObservedTimestamp"},
			{"metadata", "annotations"},
			{"metadata", "creationTimestamp"},
			{"metadata", "finalizers"},
			{"metadata", "generation"},
			{"metadata", "managedFields"},
			{"metadata", "namespace"},
			{"metadata", "resourceVersion"},
			{"metadata", "uid"},
			{"spec", "controlPlaneRef", "namespace"},
			{"spec", "infrastructureRef", "namespace"},
			{"status"},
		},
		"AWSCluster": {
			{"metadata", "annotations"},
			{"metadata", "creationTimestamp"},
			{"metadata", "ownerReferences"},
			{"metadata", "finalizers"},
			{"metadata", "generation"},
			{"metadata", "managedFields"},
			{"metadata", "namespace"},
			{"metadata", "resourceVersion"},
			{"metadata", "uid"},
			{"status"},
		},
		"AWSMachineTemplate": {
			{"metadata", "annotations"},
			{"metadata", "creationTimestamp"},
			{"metadata", "ownerReferences"},
			{"metadata", "generation"},
			{"metadata", "managedFields"},
			{"metadata", "namespace"},
			{"metadata", "resourceVersion"},
			{"metadata", "uid"},
			{"status"},
		},
		"KubeadmControlPlane": {
			{"metadata", "annotations"},
			{"metadata", "creationTimestamp"},
			{"metadata", "ownerReferences"},
			{"metadata", "finalizers"},
			{"metadata", "generation"},
			{"metadata", "managedFields"},
			{"metadata", "namespace"},
			{"metadata", "resourceVersion"},
			{"metadata", "uid"},
			{"spec", "machineTemplate", "infrastructureRef", "namespace"},
			{"status"},
		},
		"MachineDeployment": {
			{"metadata", "annotations", "kubectl.kubernetes.io/last-applied-configuration"},
			{"metadata", "annotations", "machinedeployment.clusters.x-k8s.io/revision"},
			{"metadata", "creationTimestamp"},
			{"metadata", "ownerReferences"},
			{"metadata", "generation"},
			{"metadata", "managedFields"},
			{"metadata", "namespace"},
			{"metadata", "resourceVersion"},
			{"metadata", "uid"},
			{"status"},
		},
		"KubeadmConfigTemplate": {
			{"metadata", "annotations"},
			{"metadata", "creationTimestamp"},
			{"metadata", "ownerReferences"},
			{"metadata", "generation"},
			{"metadata", "managedFields"},
			{"metadata", "namespace"},
			{"metadata", "resourceVersion"},
			{"metadata", "uid"},
		},
	}
	return kindsToUpdate
}

func GetClusterTemplateMutator() cluster.ResourceMutatorFunc {
	kinds := getTemplateMutatorKinds()
	var mutator cluster.ResourceMutatorFunc = func(u *unstructured.Unstructured) error {
		if u == nil || u.Object == nil {
			return nil
		}
		if fields, knownKind := kinds[u.GetKind()]; knownKind {
			for _, nsField := range fields {
				_, exists, err := unstructured.NestedFieldNoCopy(u.Object, nsField...)
				if err != nil {
					fmt.Println("Failed to get field")
					return err
				}
				if exists {
					unstructured.RemoveNestedField(u.Object, nsField...)
				}
			}
		}
		return nil
	}
	return mutator
}

func getNamespaceMutator(targetNamespace string) cluster.ResourceMutatorFunc {
	kinds := getNamespaceFieldsToBeUpdated()
	var namespaceMutator cluster.ResourceMutatorFunc = func(u *unstructured.Unstructured) error {
		if u == nil || u.Object == nil {
			return nil
		}
		if u.GetNamespace() != "" {
			u.SetNamespace(targetNamespace)
		}
		if fields, knownKind := kinds[u.GetKind()]; knownKind {
			for _, nsField := range fields {
				_, exists, err := unstructured.NestedFieldNoCopy(u.Object, nsField...)
				if err != nil {
					fmt.Println("Failed to get field")
					return err
				}
				if exists {
					err := unstructured.SetNestedField(u.Object, targetNamespace, nsField...)
					if err != nil {
						fmt.Println("Failed to set field")
						return err
					}
				}
			}
		}
		return nil
	}
	return namespaceMutator
}

func getNamespaceFieldsToBeUpdated() map[string][][]string {
	kindsToUpdate := map[string][][]string{
		"Cluster": {
			{"metadata", "namespace"},
			{"spec", "controlPlaneRef", "namespace"},
			{"spec", "infrastructureRef", "namespace"},
		},
		"KubeadmControlPlane": {
			{"spec", "machineTemplate", "infrastructureRef", "namespace"},
		},
		"Machine": {
			{"spec", "bootstrap", "configRef", "namespace"},
			{"spec", "infrastructureRef", "namespace"},
		},
	}
	return kindsToUpdate
}
