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
//     from inside the container (this is what unblocks the "real ps"9	eazsfxcv
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
	"fmt"
	"mycontainer/utils/errorHandlers"
	"os"
	"syscall"
)

// TODO: design your own exported function signature(s) here.
func RootFSChroot(path string) error {
	err := syscall.Chroot(path)

	syscall.Chdir("/")
	return err
}

func RootFSPivot(path string) error {

	return errorHandlers.RunSteps(
		func() error { return os.MkdirAll(path+"/.oldroot", 0700) },
		func() error { return syscall.Mount(path, path, "", syscall.MS_BIND, "") },
		func() error {
			fmt.Println("=========Pivot root==============")
			return syscall.PivotRoot(path, path+"/.oldroot")
		},
		func() error {
			fmt.Println("=========Chdir==============")
			return syscall.Chdir("/")
		},
		func() error {
			fmt.Println("=========Unmount==============")
			return syscall.Unmount(".oldroot", syscall.MNT_DETACH)
		},
		func() error {
			fmt.Println("=========Remove old root==============")
			return os.Remove(".oldroot")
		},
	)
}

func ProcRemount() error {
	return errorHandlers.RunSteps(
		func() error { return os.MkdirAll("/proc", 0755) },
		func() error { return syscall.Mount("proc", "/proc", "proc", 0, "") },
	)
}
