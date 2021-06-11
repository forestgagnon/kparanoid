package cmd

import (
	"os"

	"github.com/spf13/cobra"

	"github.com/forestgagnon/kparanoid/bin/cli/cluster"
)

var cmdOpen = &cobra.Command{
	Use:   "open clustername",
	Short: "open an isolated session with a single cluster",
	Args:  cobra.ExactArgs(1),
	Run: func(c *cobra.Command, args []string) {
		requireConfigDir(c)
		requireConfigHostDir(c)
		name := args[0]
		if err := cluster.OpenClusterInteractiveSession(name); err != nil {
			c.PrintErrln(err.Error())
			os.Exit(1)
		}
	},
}

func init() {
	cmdRoot.AddCommand(cmdOpen)
	kubectlVersionFlagParse(cmdOpen)
}
