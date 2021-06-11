package cmd

import (
	"os"
	"strings"

	"github.com/forestgagnon/kparanoid/bin/cli/cfg"
	"github.com/forestgagnon/kparanoid/bin/cli/cluster"
	"github.com/spf13/cobra"
)

var cmdAddCluster = &cobra.Command{
	Use:       "add-cluster",
	Short:     "add a config for a new kubernetes cluster",
	ValidArgs: []string{"gke", "vanilla"},
	Args:      cobra.OnlyValidArgs,
	Run: func(c *cobra.Command, args []string) {
		if len(args) == 0 {
			c.Help()
			os.Exit(1)
		}
	},
}

var cmdAddClusterGKE = &cobra.Command{
	Use:   "gke clustername",
	Short: "add a Google Kubernetes Engine cluster which uses gcloud authentication",
	Args:  cobra.ExactArgs(1),
	Run: func(c *cobra.Command, args []string) {
		requireConfigDir(c)
		requireConfigHostDir(c)
		name := args[0]
		if err := cluster.AddGKECluster(name); err != nil {
			c.PrintErrln(err.Error())
			os.Exit(1)
		}
	},
}

var cmdAddClusterVanilla = &cobra.Command{
	Use:   "vanilla clustername",
	Short: "add a plain cluster with a kubeconfig file",
	Args:  cobra.ExactArgs(1),
	Run: func(c *cobra.Command, args []string) {
		requireConfigDir(c)
		requireConfigHostDir(c)
		name := args[0]
		if err := cluster.AddVanillaCluster(name); err != nil {
			c.PrintErrln(err.Error())
			os.Exit(1)
		}
	},
}

var cmdListClusters = &cobra.Command{
	Use:     "list-clusters",
	Aliases: []string{"ls"},
	Short:   "list clusters",
	Args:    cobra.NoArgs,
	Run: func(c *cobra.Command, args []string) {
		requireConfigDir(c)
		list, err := cluster.ListClusters()
		if err != nil {
			c.PrintErrln(err.Error())
			os.Exit(1)
		}
		b := strings.Builder{}
		for _, s := range list {
			_, _ = b.WriteString(s + "\n")
		}
		c.Print(b.String())
	},
}

var cmdRemoveCluster = &cobra.Command{
	Use:   "remove-cluster clustername",
	Short: "permanently remove a cluster config from kparanoid",
	Args:  cobra.ExactArgs(1),
	Run: func(c *cobra.Command, args []string) {
		requireConfigDir(c)
		name := args[0]
		if err := cluster.RemoveClusterConfig(name); err != nil {
			c.PrintErrln(err.Error())
			os.Exit(1)
		}
	},
}

func init() {
	cmdRoot.AddCommand(cmdAddCluster)
	cmdAddCluster.AddCommand(cmdAddClusterVanilla)
	cmdAddCluster.AddCommand(cmdAddClusterGKE)

	cmdAddClusterGKE.Flags().StringVar(
		&cfg.Cfg.Flags.AddClusterGKE.GcloudClusterCredentialsCommand,
		"gcloud-cluster-creds-cmd", "",
		`the gcloud command which sets up a kubeconfig for a GKE cluster.
Typically, this is a form of 'gcloud container clusters get-credentials'`,
	)

	cmdAddClusterGKE.MarkFlagRequired("gcloud-cluster-creds-cmd")

	cmdRoot.AddCommand(cmdListClusters)
	cmdRoot.AddCommand(cmdRemoveCluster)
}
