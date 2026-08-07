// Package namespaces handles process isolation via Linux namespaces.
//
// Day 1 scope:
//   - Launch a child process with new PID and UTS namespaces
//     (syscall.CLONE_NEWPID | syscall.CLONE_NEWUTS via SysProcAttr.Cloneflags)
//   - Set a custom hostname inside the new UTS namespace
//     (syscall.Sethostname — must run from inside the new namespace,
//     i.e. in the child, not the parent)
//
// Things to design yourself:
//   - What does the exported entrypoint look like? (e.g. does it take a
//     target command + args and return an *exec.Cmd ready to Start(),
//     or does it also handle Start()/Wait() itself?)
//   - How do you verify from inside the child that the namespace took
//     effect? (os.Getpid() should report 1; reading /proc/self/status
//     is another option since real `ps` needs /proc remounted, which
//     isn't in scope until Day 3)
//   - Where does hostname get set — before exec, or does the child
//     need to call Sethostname itself after it starts running?
package namespaces

import (
	"os/exec"
	"syscall"
)

// TODO: design your own exported function signature(s) here.
// Likely need something like a constructor that returns a configured
// *exec.Cmd, given a target command and args, with Cloneflags set.
