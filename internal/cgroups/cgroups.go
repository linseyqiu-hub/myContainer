// Package cgroups handles resource limiting via cgroups v2.
//
// Day 4 scope:
//   - Create a new cgroup (a directory under /sys/fs/cgroup/...)
//   - Write the target process's PID into that cgroup's `cgroup.procs`
//     file to assign it
//   - Write a limit into `memory.max` (or `cpu.max`) to cap resource
//     usage
//   - Demonstrate enforcement: run something memory-hungry inside and
//     observe an OOM-kill, or something CPU-bound and observe throttling
//
// Before writing any code: confirm you're actually on cgroups v2 —
//
//	cat /sys/fs/cgroup/cgroup.controllers
//
// If that file doesn't exist, you're on v1 and the paths/filenames
// differ (e.g. memory.limit_in_bytes instead of memory.max).
//
// Things to design yourself:
//   - Exported function shape(s) — probably something like "create a
//     cgroup with these limits" and "assign this PID to it"
//   - Where does cleanup happen? cgroups you create need to be removed
//     (rmdir the cgroup directory) after the process exits, or they'll
//     accumulate on your system.
package cgroups

import (
	errorHandlers "mycontainer/utils/errorHandlers"
	"os"
	"strconv"
)

// TODO: design your own exported function signature(s) here.
func Setup(id string, memMax string) error {
	// mkdir /sys/fs/cgroup/<id>, write memMax into memory.max if provided
	path := "/sys/fs/cgroup/" + id
	file := path + "/memory.max"
	return errorHandlers.RunSteps(
		func() error { return os.MkdirAll(path, 0755) },
		func() error {
			if len(memMax) == 0 {
				return nil
			}
			return os.WriteFile(file, []byte(memMax), 0644)
		},
	)
}

func AddProcess(id string, pid int) error {
	// write pid into /sys/fs/cgroup/<id>/cgroup.procs
	path := "/sys/fs/cgroup/" + id + "/cgroup.procs"
	return os.WriteFile(path, []byte(strconv.Itoa(pid)), 0644)

}

func Cleanup(id string) error {
	// os.Remove /sys/fs/cgroup/<id>  (must be empty of members — safe once process has exited)
	path := "/sys/fs/cgroup/" + id
	return os.Remove(path)

}
