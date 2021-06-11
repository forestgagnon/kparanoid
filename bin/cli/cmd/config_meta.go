package cmd

import (
	"os"
	"path"

	"github.com/forestgagnon/kparanoid/bin/cli/cfg"
	"github.com/spf13/cobra"
)

var cmdClusterConf = &cobra.Command{
	Use:   "cluster-config",
	Short: "manage kparanoid configuration for a cluster",
	Args:  cobra.NoArgs,
	Run: func(c *cobra.Command, args []string) {
		if len(args) == 0 {
			c.Help()
			os.Exit(1)
		}
	},
}

var cmdClusterConfFilepath = &cobra.Command{
	Use:   "get-filepath",
	Short: "returns the filepath of the kparanoid config file for a given cluster.",
	Args:  cobra.ExactArgs(1),
	Run: func(c *cobra.Command, args []string) {
		c.Println(path.Join(cfg.Cfg.ClusterPathHost(args[0]), "cluster.json"))
	},
}

func init() {
	cmdRoot.AddCommand(cmdClusterConf)
	cmdClusterConf.AddCommand(cmdClusterConfFilepath)
}
