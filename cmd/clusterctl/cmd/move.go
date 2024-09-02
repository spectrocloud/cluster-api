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

package cmd

import (
	"context"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"

	"sigs.k8s.io/cluster-api/cmd/clusterctl/client"
)

type moveOptions struct {
	fromKubeconfig        string
	fromKubeconfigContext string
	toKubeconfig          string
	toKubeconfigContext   string
	namespace             string
	fromDirectory         string
	toDirectory           string
	dryRun                bool

	// palette specific options
	toNamespace        string
	clusterName        string
	ignoreClusterClass bool
	toPaletteCRD       string
}

var mo = &moveOptions{}

var moveCmd = &cobra.Command{
	Use:     "move",
	GroupID: groupManagement,
	Short:   "Move Cluster API objects and all dependencies between management clusters",
	Long: LongDesc(`
		Move Cluster API objects and all dependencies between management/workload clusters.

		Note: The destination cluster MUST have the required provider components installed.`),

	Example: Examples(`
		Move Cluster API objects and all dependencies between management clusters.
		palettectl move -n source-namespace --clusterName cluster-to-be-moved --to-namespace cluster-87wekfg858qwefgqliufgp578 --to-kubeconfig=target-kubeconfig.yaml

		Write Cluster API objects in palette consumable template format from a management cluster to directory.
		palettectl move -n source-namespace --clusterName cluster-to-be-moved --to-template-directory /tmp/backup-directory
	`),
	Args: cobra.NoArgs,
	RunE: func(*cobra.Command, []string) error {
		return runMove()
	},
}

func init() {
	moveCmd.Flags().StringVar(&mo.fromKubeconfig, "kubeconfig", "",
		"Path to the kubeconfig file for the source management cluster. If unspecified, default discovery rules apply.")
	moveCmd.Flags().StringVar(&mo.toKubeconfig, "to-kubeconfig", "",
		"Path to the kubeconfig file to use for the destination management cluster.")
	moveCmd.Flags().StringVar(&mo.fromKubeconfigContext, "kubeconfig-context", "",
		"Context to be used within the kubeconfig file for the source management cluster. If empty, current context will be used.")
	moveCmd.Flags().StringVar(&mo.toKubeconfigContext, "to-kubeconfig-context", "",
		"Context to be used within the kubeconfig file for the destination management cluster. If empty, current context will be used.")
	moveCmd.Flags().StringVarP(&mo.namespace, "namespace", "n", "",
		"The namespace where the workload cluster is hosted. If unspecified, the current context's namespace is used.")
	moveCmd.Flags().BoolVar(&mo.dryRun, "dry-run", false,
		"Enable dry run, don't really perform the move actions")
	//moveCmd.Flags().StringVar(&mo.toDirectory, "to-directory", "",
	//	"Write Cluster API objects and all dependencies from a management cluster to directory.")
	//moveCmd.Flags().StringVar(&mo.fromDirectory, "from-directory", "",
	//	"Read Cluster API objects and all dependencies from a directory into a management cluster.")

	// palette specific flags
	moveCmd.Flags().StringVar(&mo.toPaletteCRD, "to-template-directory", "",
		"Write Cluster API objects output in palette template format to directory")
	moveCmd.Flags().StringVar(&mo.toNamespace, "to-namespace", "",
		"Namespace where the cluster needs to be moved on destination cluster")
	moveCmd.Flags().StringVar(&mo.clusterName, "clusterName", "",
		"The cluster that needs to be moved")
	//moveCmd.Flags().BoolVar(&mo.ignoreClusterClass, "ignore-cluster-class", false,
	//	"Ignore cluster class while moving the resources to target cluster")

	moveCmd.MarkFlagsMutuallyExclusive("to-template-directory", "to-kubeconfig")
	//moveCmd.MarkFlagsMutuallyExclusive("from-directory", "to-directory")
	//moveCmd.MarkFlagsMutuallyExclusive("from-directory", "kubeconfig")

	RootCmd.AddCommand(moveCmd)
}

func runMove() error {
	ctx := context.Background()

	if mo.toDirectory == "" &&
		mo.fromDirectory == "" &&
		mo.toKubeconfig == "" &&
		mo.toPaletteCRD == "" &&
		!mo.dryRun {
		return errors.New("please specify a target cluster using the --to-kubeconfig flag when not using --dry-run, --to-directory or --from-directory")
	}

	c, err := client.New(ctx, cfgFile)
	if err != nil {
		return err
	}

	return c.Move(ctx, client.MoveOptions{
		FromKubeconfig: client.Kubeconfig{Path: mo.fromKubeconfig, Context: mo.fromKubeconfigContext},
		ToKubeconfig:   client.Kubeconfig{Path: mo.toKubeconfig, Context: mo.toKubeconfigContext},
		//FromDirectory:  mo.fromDirectory,
		//ToDirectory:    mo.toDirectory,
		Namespace:    mo.namespace,
		DryRun:       mo.dryRun,
		ClusterName:  mo.clusterName,
		ToNamespace:  mo.toNamespace,
		ToPaletteCRD: mo.toPaletteCRD,
	})
}
