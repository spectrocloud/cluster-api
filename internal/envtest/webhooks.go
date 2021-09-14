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

package envtest

import (
	"os"
	"path"
	"path/filepath"
	goruntime "runtime"
	"strings"
	"time"

	"k8s.io/klog/v2"
	utilyaml "sigs.k8s.io/cluster-api/util/yaml"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/envtest"
)

const (
	mutatingWebhookKind   = "MutatingWebhookConfiguration"
	validatingWebhookKind = "ValidatingWebhookConfiguration"
	mutatingwebhook       = "mutating-webhook-configuration"
	validatingwebhook     = "validating-webhook-configuration"
)

func initWebhookInstallOptions() envtest.WebhookInstallOptions {
	validatingWebhooks := []client.Object{}
	mutatingWebhooks := []client.Object{}

	// Get the root of the current file to use in CRD paths.
	_, filename, _, _ := goruntime.Caller(0) //nolint
	root := path.Join(path.Dir(filename), "..", "..")
	configyamlFile, err := os.ReadFile(filepath.Join(root, "config", "webhook", "manifests.yaml"))
	if err != nil {
		klog.Fatalf("Failed to read core webhook configuration file: %v ", err)
	}
	if err != nil {
		klog.Fatalf("failed to parse yaml")
	}
	// append the webhook with suffix to avoid clashing webhooks. repeated for every webhook
	mutatingWebhooks, validatingWebhooks, err = appendWebhookConfiguration(mutatingWebhooks, validatingWebhooks, configyamlFile, "config")
	if err != nil {
		klog.Fatalf("Failed to append core controller webhook config: %v", err)
	}

	bootstrapyamlFile, err := os.ReadFile(filepath.Join(root, "bootstrap", "kubeadm", "config", "webhook", "manifests.yaml"))
	if err != nil {
		klog.Fatalf("Failed to get bootstrap yaml file: %v", err)
	}
	mutatingWebhooks, validatingWebhooks, err = appendWebhookConfiguration(mutatingWebhooks, validatingWebhooks, bootstrapyamlFile, "bootstrap")

	if err != nil {
		klog.Fatalf("Failed to append bootstrap controller webhook config: %v", err)
	}
	controlplaneyamlFile, err := os.ReadFile(filepath.Join(root, "controlplane", "kubeadm", "config", "webhook", "manifests.yaml"))
	if err != nil {
		klog.Fatalf(" Failed to get controlplane yaml file err: %v", err)
	}
	mutatingWebhooks, validatingWebhooks, err = appendWebhookConfiguration(mutatingWebhooks, validatingWebhooks, controlplaneyamlFile, "cp")
	if err != nil {
		klog.Fatalf("Failed to append controlplane controller webhook config: %v", err)
	}
	return envtest.WebhookInstallOptions{
		MaxTime:                      20 * time.Second,
		PollInterval:                 time.Second,
		ValidatingWebhooks:           validatingWebhooks,
		MutatingWebhooks:             mutatingWebhooks,
		LocalServingHostExternalName: os.Getenv("CAPI_WEBHOOK_HOSTNAME"),
	}
}

// Mutate the name of each webhook, because kubebuilder generates the same name for all controllers.
// In normal usage, kustomize will prefix the controller name, which we have to do manually here.
func appendWebhookConfiguration(mutatingWebhooks []client.Object, validatingWebhooks []client.Object, configyamlFile []byte, tag string) ([]client.Object, []client.Object, error) {
	objs, err := utilyaml.ToUnstructured(configyamlFile)
	if err != nil {
		klog.Fatalf("failed to parse yaml")
	}
	// look for resources of kind MutatingWebhookConfiguration
	for i := range objs {
		o := objs[i]
		if o.GetKind() == mutatingWebhookKind {
			// update the name in metadata
			if o.GetName() == mutatingwebhook {
				o.SetName(strings.Join([]string{mutatingwebhook, "-", tag}, ""))
				mutatingWebhooks = append(mutatingWebhooks, &o)
			}
		}
		if o.GetKind() == validatingWebhookKind {
			// update the name in metadata
			if o.GetName() == validatingwebhook {
				o.SetName(strings.Join([]string{validatingwebhook, "-", tag}, ""))
				validatingWebhooks = append(validatingWebhooks, &o)
			}
		}
	}
	return mutatingWebhooks, validatingWebhooks, err
}
