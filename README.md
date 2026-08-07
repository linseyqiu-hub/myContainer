# mycontainer

A minimal container runtime built from scratch in Go, to understand what
`docker run` actually does underneath — namespaces, cgroups, and rootfs
isolation. Not a Docker clone: no image layers, no registry client, no
networking, no polished CLI. Just enough to see the real primitives work.

## Scope by day

- Day 1 — `internal/namespaces`: PID + UTS namespace isolation
- Day 2 — `internal/rootfs`: chroot into a minimal rootfs
- Day 3 — `internal/rootfs`: upgrade to pivot_root + mount namespace, /proc remount
- Day 4 — `internal/cgroups`: memory + CPU limits via cgroups v2
- Day 5 — `internal/container`: wire everything into one `mycontainer run <cmd>`
- Day 6 — buffer / real Docker usage practice (Dockerfile, docker-compose)

## Layout

- `cmd/mycontainer/main.go` — CLI entrypoint only, no core logic
- `internal/namespaces/` — namespace + hostname isolation
- `internal/rootfs/` — filesystem isolation (chroot/pivot_root)
- `internal/cgroups/` — resource limits
- `internal/container/` — orchestration, ties the above together
- `testdata/rootfs/` — placeholder for a minimal root filesystem (Day 2+)

## Known gotcha to watch for on Day 1

The classic "PID 1 problem": a plain re-exec of a target command inside a
new PID namespace often needs a two-step child re-exec pattern — the first
process clones into the new namespace, then the actual target command is
exec'd as PID 1 of that namespace, not launched as a grandchild. Worth
reading about before assuming something's broken.
