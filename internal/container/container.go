// Package container orchestrates namespaces, rootfs, and cgroups into a
// single `mycontainer run <cmd>` flow.
//
// Day 5 scope:
//   - One exported entrypoint that cmd/mycontainer/main.go calls
//   - Wires together, in the correct order:
//       1. internal/namespaces  — create the namespaced child process
//       2. internal/rootfs      — chroot/pivot_root once inside the
//                                 new namespaces
//       3. internal/cgroups     — apply resource limits before/as the
//                                 target command actually execs
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

// TODO: design your own exported entrypoint and any config struct here.
