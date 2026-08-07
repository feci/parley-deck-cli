//go:build !darwin && !freebsd && !netbsd && !openbsd && !linux

package app

// hasTTYSupported reports whether this platform can PROVE a controlling terminal.
//
// Windows and any other platform without the termios ioctl cannot, and the weaker
// character-device test is not an acceptable substitute: /dev/null is a character device, which is
// precisely how an agent run redirects stdin. Rather than accept a check that does not check, the
// publisher fails closed here — an attended publish is unavailable on this platform.
//
// Deleting this file (as an earlier revision did) does not make the platform stricter; it makes the
// build fail, which is how a broken Windows binary nearly shipped.
const hasTTYSupported = false

func platformHasTTY() bool { return false }
