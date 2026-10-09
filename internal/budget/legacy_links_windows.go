package budget

import (
	"errors"
	"os"
	"syscall"
)

func legacySingleLink(path string) error {
	f, err := os.Open(path)
	if err != nil {
		return err
	}
	defer f.Close()
	var info syscall.ByHandleFileInformation
	if err := syscall.GetFileInformationByHandle(syscall.Handle(f.Fd()), &info); err != nil {
		return err
	}
	if info.NumberOfLinks != 1 {
		return errors.New("legacy history must not have hard-link aliases")
	}
	return nil
}
