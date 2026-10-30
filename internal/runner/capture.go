package runner

import (
	"fmt"
	"io"
	"os"
	"time"

	"parley-deck-cli/internal/telemetry"
)

// Capture is best-effort diagnostic data. It supplies no membership authority.
func captureValidatedFailure(r telemetry.Record, valid bool, path, directory string) error {
	if valid || !r.Metadata.CaptureZcodeFailures || r.Metadata.Adapter != "zcode" || path == "" || r.Outcome == nil || r.CompletedAt == nil || r.StartedAt == nil || r.Outcome.ExitCode == nil || *r.Outcome.ExitCode == 0 {
		return nil
	}
	info, err := os.Lstat(path)
	if err != nil {
		return err
	}
	if !info.Mode().IsRegular() {
		return fmt.Errorf("capture input is not a regular file")
	}
	f, err := os.Open(path)
	if err != nil {
		return err
	}
	defer f.Close()
	opened, err := f.Stat()
	if err != nil || !os.SameFile(info, opened) {
		return fmt.Errorf("capture input changed")
	}
	truncated := info.Size() > telemetry.CaptureFileLimit
	if truncated {
		if _, err = f.Seek(-telemetry.CaptureFileLimit, io.SeekEnd); err != nil {
			return err
		}
	}
	raw, err := io.ReadAll(io.LimitReader(f, telemetry.CaptureFileLimit))
	if err != nil {
		return err
	}
	_, err = telemetry.CaptureZcodeFailure(directory, r.InvocationID, raw, truncated, time.Now())
	return err
}
