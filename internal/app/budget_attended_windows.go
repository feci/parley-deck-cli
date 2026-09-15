//go:build windows

package app

import "golang.org/x/sys/windows"

func budgetAttendance() (supported, attended bool) {
	input, err := windows.GetStdHandle(windows.STD_INPUT_HANDLE)
	if err != nil {
		return true, false
	}
	var mode uint32
	return true, windows.GetConsoleMode(input, &mode) == nil
}
