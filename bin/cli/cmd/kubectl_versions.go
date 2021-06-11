package cmd

import (
	"os"

	"github.com/forestgagnon/kparanoid/bin/cli/cluster"
	"github.com/spf13/cobra"
)

var cmdListKubectlVersions = &cobra.Command{
	Use:   "kubectl-versions",
	Short: "list available kubectl versions",
	Args:  cobra.NoArgs,
	Run: func(c *cobra.Command, args []string) {
		if len(args) == 0 {
			c.Help()
			os.Exit(1)
		}
	},
}

var cmdListKubectlVersionsGKE = &cobra.Command{
	Use:   "gke",
	Short: "list available kubectl versions for GKE sessions",
	Args:  cobra.NoArgs,
	Run: func(c *cobra.Command, args []string) {
		if err := cluster.PrintAvailableKubectlVersionsApt(); err != nil {
			c.PrintErrln(err.Error())
			os.Exit(1)
		}
	},
}

var cmdListKubectlVersionsVanilla = &cobra.Command{
	Use:   "vanilla",
	Short: "list available kubectl versions for vanilla sessions",
	Args:  cobra.NoArgs,
	Run: func(c *cobra.Command, args []string) {
		if err := cluster.PrintAvailableKubectlVersionsApt(); err != nil {
			c.PrintErrln(err.Error())
			os.Exit(1)
		}
	},
}

func init() {
	cmdRoot.AddCommand(cmdListKubectlVersions)
	cmdListKubectlVersions.AddCommand(cmdListKubectlVersionsGKE)
	cmdListKubectlVersions.AddCommand(cmdListKubectlVersionsVanilla)
}
