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
//   cat /sys/fs/cgroup/cgroup.controllers
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
	"os"
	"path/filepath"
	"strconv"
)

// TODO: design your own exported function signature(s) here.
