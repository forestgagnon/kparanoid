package cmd

import (
	"fmt"
	"os"
	"path"

	"github.com/forestgagnon/kparanoid/bin/cli/cfg"
	"github.com/forestgagnon/kparanoid/bin/cli/util"
	"github.com/spf13/cobra"
)

var cmdInstall = &cobra.Command{
	Use:       "install",
	Short:     "generate install scripts for kparanoid which can be executed",
	ValidArgs: []string{"bash", "bootstrap"},
	Args:      cobra.OnlyValidArgs,
	Run: func(c *cobra.Command, args []string) {
		if len(args) == 0 {
			c.Help()
			os.Exit(1)
		}
	},
}

var cmdInstallBash = &cobra.Command{
	Use:   "bash",
	Short: "generate install script which can be piped to bash",
	Run: func(c *cobra.Command, args []string) {
		requireConfigDir(c)
		c.Print(fmt.Sprintf(bashInstallScript, cfg.Cfg.Flags.ConfigDir))
	},
}

var cmdInstallBootstrap = &cobra.Command{
	Use:    "bootstrap",
	Hidden: true,
	Short:  "adds essential files to the kparanoid config folder",
	Run: func(c *cobra.Command, args []string) {
		requireConfigDir(c)
		if os.Getenv("KPARANOID_INSTALL_FLOW") != "yes" {
			c.PrintErrln("using bootstrap outside of the installation flow doesn't make sense.")
			os.Exit(1)
		}

		if err := os.MkdirAll(cfg.Cfg.ContainmentDockerfilesPath(), 0777); err != nil {
			panic(err)
		}

		if err := os.MkdirAll(cfg.Cfg.ClustersPath(), 0777); err != nil {
			panic(err)
		}

		util.CopyFileInefficientlyOrPanic("/app/installation/containment-dockerfiles/debian-common.Dockerfile", cfg.Cfg.ContainmentDockerfile("debian-common.Dockerfile"))
		util.CopyFileInefficientlyOrPanic("/app/installation/.bash_profile", path.Join(cfg.Cfg.EnvPath(), ".bash_profile"))
	},
}

var cmdInstallPATH = &cobra.Command{
	Use:   "path",
	Short: "add the output of this to one of your shell configuration files to put 'kparanoid' in your PATH",
	Run: func(c *cobra.Command, args []string) {
		requireConfigDir(c)
		requireConfigHostDir(c)
		c.Print(fmt.Sprintf("export PATH=\"$PATH:%s/bin/\"\n", cfg.Cfg.ConfigDirOnHost))
	},
}

func init() {
	cmdRoot.AddCommand(cmdInstall)
	cmdInstall.AddCommand(cmdInstallBash)
	cmdInstall.AddCommand(cmdInstallBootstrap)
	cmdInstall.AddCommand(cmdInstallPATH)
}

var bashInstallScript = `
set -euo pipefail

config_dir="%s"

echo "Installing kparanoid to '$config_dir'"

mkdir -p "$config_dir"
mkdir -p "$config_dir/bin"

kparanoid_launcher="$config_dir/bin/kparanoid"

docker run --rm \
	-v "$config_dir:/kparanoid/conf" \
	--env KPARANOID_INSTALL_FLOW=yes \
	forestgagnon/kparanoid:1 install bootstrap --config-dir=/kparanoid/conf

cat <<EOF > "$kparanoid_launcher"
#!/usr/bin/env bash

set -euo pipefail

ttyflag=""
if [ -t 1 ] && [ -t 0 ] ; then ttyflag="--tty"; fi

docker run --rm -i \$ttyflag \
	-v "$config_dir:/kparanoid/conf" \
	-v /var/run/docker.sock:/var/run/docker.sock \
	--env KPARANOID_CONFIG_DIR_ON_HOST="$config_dir" \
	forestgagnon/kparanoid:1 --config-dir=/kparanoid/conf "\$@"
EOF
chmod +x "$kparanoid_launcher"

printf "Installation complete\n\n"
printf "To add kparanoid to your PATH, run the command below to get a line you can add to your shell configuration:\n\n"
echo "  $config_dir/bin/kparanoid install path"
`
