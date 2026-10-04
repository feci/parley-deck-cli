package fsutil

import (
	"errors"
	"fmt"
	"syscall"
	"testing"
)

func TestSyncFileFallbackENOTTYAndFatalFailures(t *testing.T) {
	for _, tc := range []struct {
		name        string
		full, plain error
		calls       int
		want        error
	}{
		{"full-success", nil, syscall.EIO, 0, nil},
		{"ENOTTY-fallback", fmt.Errorf("shared volume: %w", syscall.ENOTTY), nil, 1, nil},
		{"fallback-failure", syscall.ENOTTY, syscall.EIO, 1, syscall.EIO},
		{"full-failure", syscall.EIO, nil, 0, syscall.EIO},
	} {
		t.Run(tc.name, func(t *testing.T) {
			n := 0
			err := syncWithFallback(func() error { return tc.full }, func() error { n++; return tc.plain })
			if !errors.Is(err, tc.want) || n != tc.calls {
				t.Fatalf("error=%v fallback calls=%d", err, n)
			}
			t.Logf("full=%v fallback_calls=%d final=%v", tc.full, n, err)
		})
	}
}
