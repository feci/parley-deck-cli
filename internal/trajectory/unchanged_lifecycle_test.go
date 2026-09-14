package trajectory

import (
	"context"
	"encoding/json"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"parley-deck-cli/internal/telemetry"
)

// Each case presents a fresh, internally matching set of lifecycle records.
// No old preview digest can mask the qualification predicate being exercised.
func TestUnchangedPreviewRejectsConsistentAbnormalLifecycle(t *testing.T) {
	for _, kind := range []string{"timeout", "cancelled", "stalled", "provider-error", "signal", "requested-before-epoch", "start-before-charge", "start-before-request", "trajectory-terminal-before-start", "completion-before-trajectory-terminal"} {
		t.Run(kind, func(t *testing.T) {
			root, b := unchangedFixture(t, "true")
			command := "exit 7"
			if kind == "signal" {
				command = "kill -TERM $$"
			}
			r := unchangedProcess(t, root, b, command, false)
			s, _, err := readState(statePath(*b))
			if err != nil {
				t.Fatal(err)
			}
			a := &s.Attempts[0]
			base := filepath.Join(root, ".parley-runtime", "invocations", r.InvocationID)
			records := map[string]*telemetry.Record{}
			for _, name := range []string{"requested.json", "started.json", "terminal.json"} {
				var record telemetry.Record
				if err := json.Unmarshal(snapshotRead(t, filepath.Join(base, name)), &record); err != nil {
					t.Fatal(err)
				}
				records[name] = &record
			}
			terminal := records["terminal.json"]
			requestAt := func(at time.Time) {
				for _, record := range records {
					record.RequestedAt = at
				}
			}
			switch kind {
			case "timeout", "cancelled", "stalled", "provider-error":
				terminal.Outcome.FailureClass = telemetry.String(kind)
			case "signal":
				if *terminal.Outcome.ExitCode != -1 || *a.Terminal.ExitCode != -1 {
					t.Fatal("fixture did not actually exit on a signal")
				}
			case "requested-before-epoch":
				requestAt(s.Policy.StartedAt.Add(-time.Nanosecond))
			case "start-before-charge":
				at := a.Charge.ReservedAt.Add(-time.Nanosecond)
				if at.Before(s.Policy.StartedAt) {
					t.Fatal("fixture cannot isolate start/charge order")
				}
				requestAt(s.Policy.StartedAt)
				records["started.json"].StartedAt, terminal.StartedAt = &at, &at
			case "start-before-request":
				requestAt(terminal.StartedAt.Add(time.Nanosecond))
			case "trajectory-terminal-before-start":
				at := terminal.StartedAt.Add(-time.Nanosecond)
				if at.Before(a.Charge.ReservedAt) {
					t.Fatal("fixture cannot isolate terminal/start order")
				}
				a.Terminal.At = at
			case "completion-before-trajectory-terminal":
				at := a.Terminal.At.Add(-time.Nanosecond)
				if at.Before(*terminal.StartedAt) {
					t.Fatal("fixture cannot isolate completion/trajectory order")
				}
				terminal.CompletedAt = &at
			}
			for name, record := range records {
				raw, err := canonical(record)
				if err != nil {
					t.Fatal(err)
				}
				snapshotWrite(t, base, name, raw, 0600)
			}
			if err := writeState(statePath(*b), s); err != nil {
				t.Fatal(err)
			}
			if _, err := Inspect(context.Background(), root, "fixture"); err != nil {
				t.Fatal("fixture failed before the lifecycle qualification boundary", err)
			}
			_, err = PreviewUnchanged(context.Background(), root, "fixture", 1)
			if err == nil || !strings.Contains(err.Error(), "matching normally exited process lifecycle") {
				t.Fatal("abnormal lifecycle was not refused by the qualification predicate", err)
			}
		})
	}
}
