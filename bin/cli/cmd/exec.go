package cmd

import (
	"os"

	"github.com/forestgagnon/kparanoid/bin/cli/cfg"
	"github.com/forestgagnon/kparanoid/bin/cli/cluster"
	"github.com/spf13/cobra"
)

var cmdExec = &cobra.Command{
	Use:   "exec clustername [anything]",
	Short: "exec an arbitrary command on an isolated session for a single cluster",
	Args:  cobra.MinimumNArgs(2),
	Run: func(c *cobra.Command, args []string) {
		requireConfigDir(c)
		requireConfigHostDir(c)
		cfg.Cfg.OnlyBuildRuntimeIfNeeded = true
		cfg.Cfg.SquelchDockerBuildOutput = true
		name := args[0]
		execArgs := args[1:]
		if execArgs[0] == "--" {
			execArgs = execArgs[1:]
		}
		if err := cluster.ExecClusterNonInteractive(name, execArgs); err != nil {
			c.PrintErrln(err.Error())
			os.Exit(1)
		}
	},
}

func init() {
	cmdRoot.AddCommand(cmdExec)
	kubectlVersionFlagParse(cmdExec)
	cmdExec.Flags().SetInterspersed(false)
}
