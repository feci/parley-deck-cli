//go:build !windows

package protocolcore

import "syscall"

// noFollow refuses to open through a symlink. Combined with O_EXCL it makes "write-once" mean the
// release path itself, not whatever a planted link points at.
const noFollow = syscall.O_NOFOLLOW
