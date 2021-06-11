package cmd

import (
	"fmt"
	"os"
	"path"

	"github.com/spf13/cobra"

	"github.com/forestgagnon/kparanoid/bin/cli/cfg"
)

var cmdRoot = &cobra.Command{
	Use:   "kparanoid",
	Short: "kparanoid makes kubectl safer to use",
	Run: func(c *cobra.Command, args []string) {
		if len(args) == 0 {
			c.Help()
			os.Exit(1)
		}
	},
}

func init() {
	cfg.Cfg.ConfigDirOnHost = os.Getenv("KPARANOID_CONFIG_DIR_ON_HOST")
	cmdRoot.PersistentFlags().StringVar(
		&cfg.Cfg.Flags.ConfigDir,
		"config-dir", "",
		"kparanoid config directory. Generally, should not be specified manually, and 'kparanoid install' should be used instead.",
	)
	cmdRoot.MarkPersistentFlagRequired("config-dir")
	cmdRoot.SetOut(os.Stdout)
	cmdRoot.SetErr(os.Stderr)
}

func Execute() {
	if err := cmdRoot.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func requireConfigDir(c *cobra.Command) {
	if cfg.Cfg.Flags.ConfigDir == "" {
		c.PrintErrln("--config-dir must be specified")
		os.Exit(1)
	}

	if !path.IsAbs(cfg.Cfg.Flags.ConfigDir) {
		c.PrintErrln("--config-dir must be an absolute path")
		os.Exit(1)
	}
}

func requireConfigHostDir(c *cobra.Command) {
	if cfg.Cfg.ConfigDirOnHost == "" {
		c.PrintErrln("KPARANOID_CONFIG_DIR_ON_HOST must be specified")
		os.Exit(1)
	}

	if !path.IsAbs(cfg.Cfg.ConfigDirOnHost) {
		c.PrintErrln("KPARANOID_CONFIG_DIR_ON_HOST must be an absolute path")
		os.Exit(1)
	}
}

func kubectlVersionFlagParse(c *cobra.Command) {
	c.Flags().StringVar(
		&cfg.Cfg.Flags.KubectlVersion,
		"kubectl-version", "",
		`kubectl version to use. Must be in the format returned by on of the 'kparanoid kubectl-versions' commands`,
	)
}
