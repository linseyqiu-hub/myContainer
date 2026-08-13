// Package container orchestrates namespaces, rootfs, and cgroups into a
// single `mycontainer run <cmd>` flow.
//
// Day 5 scope:
//   - One exported entrypoint that cmd/mycontainer/main.go calls
//   - Wires together, in the correct order:
//     1. internal/namespaces  — create the namespaced child process
//     2. internal/rootfs      — chroot/pivot_root once inside the
//     new namespaces
//     3. internal/cgroups     — apply resource limits before/as the
//     target command actually execs
//   - Handles the "PID 1 problem" re-exec pattern if you haven't
//     already resolved it inside internal/namespaces
//
// Things to design yourself:
//   - What does the exported entrypoint's signature look like? What
//     does it need as input (target command? a config struct bundling
//     rootfs path + resource limits + hostname?), and what does it
//     return (an error? nothing, and it just blocks until the child
//     exits?)
//   - Where does config for a "run" live — a struct in this package,
//     or does main.go build one from parsed flags and pass it in?
package container

import (
	"flag"
	"mycontainer/internal/cgroups"
	"mycontainer/internal/namespaces"
	"mycontainer/internal/rootfs"
	"mycontainer/utils/IDGeneration"
	"mycontainer/utils/errorHandlers"
	"mycontainer/utils/types"
	"os"
	"syscall"
)

// child : 1(yes) 0(not) int; targetCommands : string; targetArgs: [] string; hostname: string; memMax:string; rootPath: string

// TODO: design your own exported entrypoint and any config struct here.

func preflightParent(config *types.Config) error {
	memMax := flag.String("mem", "", "...")
	hostName := flag.String("mem", "", "...")
	flag.Parse()
	remaining := flag.Args()
	targetCmd := remaining[1]
	targetArgs := remaining[2:]
	if len(*hostName) == 0 {
		var err error
		*hostName, err = IDGeneration.GenerateID(32)
		if err != nil {
			return err
		}

	}
	config.Child = 0
	config.Hostname = *hostName
	var errGen error
	config.ID, errGen = IDGeneration.GenerateID(32)
	if errGen != nil {
		return errGen
	}
	config.MemMax = *memMax
	config.TargetCmd = targetCmd
	config.TargetArgs = targetArgs

	return nil
}

// ["/proc/self/exe", "child", <id>, <hostname>, <memmax>, <targetCmd>, <targetArg0>, <targetArg1>, <targetArg2>,...<targetArgn>]
func preflightChild(config *types.Config) {
	config.Child = 1
	config.ID = os.Args[2]
	config.Hostname = os.Args[3]
	config.MemMax = os.Args[4]
	config.TargetCmd = os.Args[5]
	config.TargetArgs = os.Args[6:]
}

func preflight(config *types.Config) error {
	if os.Args[1] == "child" {
		preflightChild(config)
		return nil
	} else {
		return preflightParent(config)
	}
}

func Parent(config *types.Config) error {

	cmd, w, errNS := namespaces.Namespace(config)

	if errNS != nil {
		return errNS
	}
	pid := cmd.Process.Pid
	id := config.ID
	memMax := config.MemMax
	errorHandlers.RunSteps(
		func() error { return cgroups.Setup(id, memMax) },
		func() error { return cgroups.AddProcess(id, pid) },
		func() error {
			_, err := w.Write([]byte{0})
			return err
		},
	)
	cmd.Wait()
	cgroups.Cleanup(id)
	return nil
}

func Child(config *types.Config) error {
	return errorHandlers.RunSteps(
		func() error { return namespaces.SetHostname(config.Hostname) },
		func() error { return rootfs.RootFSPivot("testdata/rootfs") },
		func() error {
			pipeFromParent := os.NewFile(3, "pipe")
			buf := make([]byte, 1)
			_, err := pipeFromParent.Read(buf)
			return err
		},
		func() error {
			return syscall.Exec(config.TargetCmd, config.TargetArgs, os.Environ())
		},
	)
}

func Container() error {
	var config types.Config
	return errorHandlers.RunSteps(
		func() error {
			return preflight(&config)
		},
		func() error {
			if config.Child == 0 {
				return Parent(&config)
			} else {
				return Child(&config)
			}
		},
	)
}
