package main

import (
	"fmt"
	"os"
	"time"

	"parley-deck-cli/internal/telemetry"
)

func main() {
	raw, _ := os.ReadFile(".parley-runtime/claude1-r4/tail/retained-native-tail.stderr")
	obs, _ := time.Parse(time.RFC3339Nano, "2026-10-02T21:06:47.003Z")
	code := 1
	e := telemetry.ClassifyQuota(telemetry.QuotaInput{Adapter: "zcode", InvocationID: "claude1-r4-tail", Stderr: string(raw), ObservedAt: obs, ExitCode: &code})
	fmt.Printf("retained native tail alone: eligible=%v reason=%q\n", e.Eligible, e.Reason)
}
