//go:build !windows

package fsacl

import "os"

// EnsurePrivateStore creates dir if missing and enforces the Unix privacy
// contract byte-identically to the previous inline guard: a real (non-symlink)
// directory with no group/other permission bits.
func EnsurePrivateStore(dir string) error {
	if err := os.MkdirAll(dir, 0o700); err != nil {
		return err
	}
	return VerifyPrivateStore(dir)
}

// VerifyPrivateStore checks the contract without mutating anything.
func VerifyPrivateStore(dir string) error {
	return verifyRealPrivateDir(dir)
}

func verifyRealPrivateDir(dir string) error {
	info, err := os.Lstat(dir)
	if err != nil {
		return err
	}
	if !info.IsDir() || info.Mode()&os.ModeSymlink != 0 || info.Mode().Perm()&0o077 != 0 {
		return ErrNotPrivate
	}
	return nil
}

// DenyRead makes path unreadable for tests: mode 0 on Unix (§D.9).
func DenyRead(path string) error {
	if err := os.Chmod(path, 0o000); err != nil {
		return err
	}
	return nil
}
