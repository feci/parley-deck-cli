//go:build windows

package pidlease

// Windows boot/liveness attribution is not established. Fail closed to owner
// recovery; cross-compilation is not evidence of Windows runtime exclusivity.
func definitelyDead(pid int) bool { return false }
