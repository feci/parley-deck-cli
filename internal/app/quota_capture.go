package app

import (
	"encoding/json"
	"fmt"
	"io"
	"time"

	"parley-deck-cli/internal/config"
	"parley-deck-cli/internal/telemetry"
)

func runQuotaCapture(args []string, out, errs io.Writer) int {
	if len(args) != 1 || (args[0] != "status" && args[0] != "purge") {
		fmt.Fprintln(errs, "usage: parley quota capture status|purge")
		return 2
	}
	enabled, dir, err := config.MachineZcodeCapture()
	if err != nil {
		fmt.Fprintln(errs, "capture configuration:", err)
		return 1
	}
	state, err := telemetry.CaptureStatus(dir, time.Now(), args[0] == "purge")
	if err != nil {
		fmt.Fprintln(errs, "capture storage:", err)
		return 1
	}
	result := struct {
		Enabled       bool                       `json:"enabled"`
		Directory     string                     `json:"directory"`
		FileLimit     int                        `json:"file_limit_bytes"`
		TotalLimit    int                        `json:"total_limit_bytes"`
		RetentionDays int                        `json:"retention_days"`
		Inventory     telemetry.CaptureInventory `json:"inventory"`
	}{enabled, dir, telemetry.CaptureFileLimit, telemetry.CaptureTotalLimit, 7, state}
	if err = json.NewEncoder(out).Encode(result); err != nil {
		fmt.Fprintln(errs, err)
		return 1
	}
	return 0
}
