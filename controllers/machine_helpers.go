/*
Copyright 2019 The Kubernetes Authors.

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

package controllers

import (
	"context"

	"github.com/pkg/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	clusterv1 "sigs.k8s.io/cluster-api/api/v1alpha3"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// getActiveMachinesInCluster returns all of the active Machine objects
// that belong to the cluster with given namespace/name
func getActiveMachinesInCluster(ctx context.Context, c client.Client, namespace, name string) ([]*clusterv1.Machine, error) {
	if name == "" {
		return nil, nil
	}

	machineList := &clusterv1.MachineList{}
	labels := map[string]string{clusterv1.ClusterLabelName: name}

	if err := c.List(ctx, machineList, client.InNamespace(namespace), client.MatchingLabels(labels)); err != nil {
		return nil, errors.Wrap(err, "failed to list machines")
	}

	machines := []*clusterv1.Machine{}
	for i := range machineList.Items {
		m := &machineList.Items[i]
		if m.DeletionTimestamp.IsZero() {
			machines = append(machines, m)
		}
	}
	return machines, nil
}

// hasMatchingLabels verifies that the Label Selector matches the given Labels
func hasMatchingLabels(matchSelector metav1.LabelSelector, matchLabels map[string]string) bool {
	// This should never fail, validating webhook should catch this first
	selector, err := metav1.LabelSelectorAsSelector(&matchSelector)
	if err != nil {
		return false
	}
	// If a nil or empty selector creeps in, it should match nothing, not everything.
	if selector.Empty() {
		return false
	}
	if !selector.Matches(labels.Set(matchLabels)) {
		return false
	}
	return true
}

func getAllMachinesCountForMS(ctx context.Context, c client.Client, ms *clusterv1.MachineSet) (int, error) {
	if ms == nil {
		return 0, nil
	}

	allMachines := &clusterv1.MachineList{}
	selectorMap, err := metav1.LabelSelectorAsMap(&ms.Spec.Selector)

	if err != nil {
		return 0, errors.Wrapf(err, "failed to convert MachineSet %q label selector to a map", ms.Name)
	}
	if err := c.List(ctx, allMachines, client.InNamespace(ms.Namespace), client.MatchingLabels(selectorMap)); err != nil {
		return 0, errors.Wrap(err, "failed to list machines")
	}

	return len(allMachines.Items), nil
}
