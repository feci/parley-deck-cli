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

// ProtectPrivateFile is a no-op on Unix: os.CreateTemp already creates the
// archive temp file with mode 0600, byte-identical to the previous behavior.
// Windows must set the owner-only DACL because mode bits are synthesized.
func ProtectPrivateFile(file string) error {
	_ = file
	return nil
}

// VerifyPrivateFile enforces the file half of the contract byte-identically
// to the previous inline read-back check: a regular file with no group/other
// permission bits. info is the caller's Lstat result.
func VerifyPrivateFile(file string, info os.FileInfo) error {
	_ = file
	if !info.Mode().IsRegular() || info.Mode().Perm()&0o077 != 0 {
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
