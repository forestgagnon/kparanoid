package main

import (
	"syscall"

	"github.com/forestgagnon/kparanoid/bin/cli/cmd"
)

func init() {
	// Hack to reduce the chance of leaving problematic root-owned files
	// on the host, since this runs from a container.
	syscall.Umask(0)
}

func main() {
	cmd.Execute()
}
