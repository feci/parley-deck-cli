//go:build !windows

package app

func budgetAttendance() (supported, attended bool) {
	return hasTTYSupported, hasTTYSupported && platformHasTTY()
}
