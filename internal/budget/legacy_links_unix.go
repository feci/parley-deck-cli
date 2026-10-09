//go:build !windows

package budget

import (
	"errors"
	"os"
	"syscall"
)

func legacySingleLink(path string) error {
	st, err := os.Lstat(path)
	if err != nil {
		return err
	}
	native, ok := st.Sys().(*syscall.Stat_t)
	if !ok || native.Nlink != 1 {
		return errors.New("legacy history must not have hard-link aliases")
	}
	return nil
}
