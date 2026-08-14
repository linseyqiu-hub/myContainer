// Command mycontainer is the CLI entrypoint for a minimal container runtime.
//
// This file should stay thin: parse os.Args, figure out which subcommand
// was requested, and hand off to internal/container. No namespace/cgroup/
// rootfs logic belongs here directly.
//
// Expected subcommands (design the actual dispatch yourself):
//   - `mycontainer run <cmd> [args...]`   — the user-facing entrypoint
//   - a hidden/internal re-exec subcommand for the PID-1 child process
//     (see README's "PID 1 problem" note — you'll likely need this by
//     the time namespaces.go is wired up)
package main

import (
	"fmt"
	"mycontainer/internal/container"
	"os"
)

func main() {
	// TODO: parse os.Args, dispatch to internal/container.
	// You'll import "mycontainer/internal/container" once you've
	// designed its exported entrypoint.

	err := container.Container()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	os.Exit(0)

}
