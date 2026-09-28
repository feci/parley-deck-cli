//go:build windows

package pipeline

import (
	"errors"

	"golang.org/x/sys/windows"
)

// errInvalidName reports whether err means the path itself is unrepresentable
// on this OS (ERROR_INVALID_NAME): a legacy raw gate name carrying '>' can
// never exist on Windows, so the legacy read fallback treats that as
// "not found" rather than a read failure.
func errInvalidName(err error) bool {
	return errors.Is(err, windows.ERROR_INVALID_NAME)
}
