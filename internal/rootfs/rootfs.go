// Package rootfs handles filesystem isolation for the container process.
//
// Day 2 scope:
//   - syscall.Chroot into a minimal prepared root directory
//     (testdata/rootfs/ is the placeholder location — populate it with
//     a tiny set of binaries, e.g. a static busybox, before testing)
//
// Day 3 scope (upgrade Day 2's approach):
//   - Add CLONE_NEWNS (mount namespace) to the namespace flags in
//     internal/namespaces
//   - Replace chroot with pivot_root — swap the root mount entirely
//     instead of just remapping the visible path, then unmount the
//     old root so it's unreachable
//   - Remount /proc inside the new root so `ps` behaves correctly
//     from inside the container (this is what unblocks the "real ps"
//     verification you deferred from Day 1)
//
// Things to design yourself:
//   - Exported function shape(s) for "prepare and enter this rootfs"
//   - How do you sequence this relative to namespace creation in
//     internal/namespaces — does rootfs setup happen in the parent
//     before Start(), or does the child do it after it starts running
//     inside its new namespaces? (Hint: chroot/pivot_root need to run
//     from inside the target namespaces, not before they're created.)
package rootfs

import (
	"os"
	"syscall"
)

// TODO: design your own exported function signature(s) here.
