package app

import (
	"bytes"
	"context"
	"encoding/json"
	"os"
	"parley-deck-cli/internal/budget"
	"parley-deck-cli/internal/runner"
	"parley-deck-cli/internal/trajectory"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

// Drive the actual replacement helper to completion but obstruct publication.
// Recovery must publish its retained evidence without relaunching any process.
func TestAppVerifierParentPublicationRecovery(t *testing.T) {
	f := refusedVerifierAppFixture(t)
	ctx := context.Background()
	recovery := applyVerifierRecovery(t, f)
	code, plan, msg := runVerifierRelaunchCLI(t, ctx, f.binary, "--dir", f.root, "--idea", "idea-x", "--run", f.runID)
	if code != 0 {
		t.Fatalf("preview: %d %s", code, msg)
	}
	observation := filepath.Join(f.runBase(), "parent-recovered.json")
	if err := os.Mkdir(observation, 0700); err != nil {
		t.Fatal(err)
	}
	spend := budget.Store{Dir: filepath.Join(t.TempDir(), "launch"), Scope: "publication-interruption"}
	authorized := runner.WithLaunchBudget(ctx, runner.LaunchBudget{Store: spend})
	code, _, msg = runVerifierRelaunchCLI(t, authorized, f.binary, "--dir", f.root, "--idea", "idea-x", "--run", f.runID, "--sha256", plan.SHA256, "--timeout", (90 * time.Second).String(), "--yes")
	if code == 0 {
		t.Fatal("publication unexpectedly succeeded despite obstruction")
	}
	if f.verifierStarts(t) != recovery.Preview.InvocationID+"\n" {
		t.Fatalf("replacement did not really run: %s", msg)
	}
	if err := os.Remove(observation); err != nil {
		t.Fatal(err)
	}
	f.assertUnchanged(t)
	// Existing routes cannot restart the consumed replacement or reinterpret the
	// retained failed parent as missing. Keep those refusals when adding recovery.
	code, _, msg = runVerifierRelaunchCLI(t, authorized, f.binary, "--dir", f.root, "--idea", "idea-x", "--run", f.runID, "--sha256", plan.SHA256, "--yes")
	if code == 0 || !strings.Contains(msg, "recovered independent trajectory verifier has no matching successful observed terminal") {
		t.Fatalf("duplicate relaunch did not fail at the consumed-invocation boundary: %d %s", code, msg)
	}
	if _, err := trajectory.PreviewParentRecovery(ctx, f.root, "idea-x", f.runID); err == nil || !strings.Contains(err.Error(), "retained parent contains conflicting identity, failure or assessment") {
		t.Fatalf("missing-parent route did not refuse the retained failure: %v", err)
	}
	if _, err := trajectory.PreviewReconciliation(ctx, f.root, "idea-x", f.runID); err == nil || !strings.Contains(err.Error(), "parent did not retain the independently derived successful comparison") {
		t.Fatalf("unpublished observation did not refuse the conflicting original parent: %v", err)
	}
	before, err := spend.Inspect(ctx)
	if err != nil || len(before.Entries) != 1 {
		t.Fatalf("launch charge: %+v %v", before, err)
	}
	var stdout, stderr bytes.Buffer
	args := []string{"recover-verifier-parent", "--dir", f.root, "--idea", "idea-x", "--run", f.runID}
	if code := runTrajectory(ctx, args, &stdout, &stderr); code != 0 {
		t.Fatalf("retained complete replacement is stranded after publication failure: %d %s", code, stderr.String())
	}
	var output struct {
		Preview trajectory.RecoveredParentPreview `json:"preview"`
		SHA256  string                            `json:"sha256"`
		Applied bool                              `json:"applied"`
	}
	if err := json.Unmarshal(stdout.Bytes(), &output); err != nil || output.Applied || output.SHA256 == "" {
		t.Fatalf("preview: %s %v", stdout.String(), err)
	}
	if _, err := os.Stat(observation); !os.IsNotExist(err) {
		t.Fatalf("preview wrote observation: %v", err)
	}
	stale := append(append([]string{}, args...), "--sha256", strings.Repeat("0", 64), "--yes")
	stdout.Reset()
	stderr.Reset()
	if runTrajectory(ctx, stale, &stdout, &stderr) == 0 {
		t.Fatal("stale approval accepted")
	}
	applied := append(append([]string{}, args...), "--sha256", output.SHA256, "--yes")
	var retained []byte
	for n := 0; n < 2; n++ {
		stdout.Reset()
		stderr.Reset()
		if code := runTrajectory(ctx, applied, &stdout, &stderr); code != 0 {
			t.Fatalf("apply/replay: %d %s", code, stderr.String())
		}
		data, err := os.ReadFile(observation)
		if err != nil {
			t.Fatal(err)
		}
		if n == 0 {
			retained = data
		} else if !bytes.Equal(data, retained) {
			t.Fatal("replay rewrote observation")
		}
	}
	if f.verifierStarts(t) != recovery.Preview.InvocationID+"\n" {
		t.Fatal("publication recovery relaunched verifier")
	}
	after, err := spend.Inspect(ctx)
	if err != nil {
		t.Fatal(err)
	}
	b1, _ := json.Marshal(before)
	b2, _ := json.Marshal(after)
	if !bytes.Equal(b1, b2) {
		t.Fatal("publication recovery changed launch accounting")
	}
	f.assertUnchanged(t)
	reconciliation, err := trajectory.PreviewReconciliation(ctx, f.root, "idea-x", f.runID)
	if err != nil {
		t.Fatal(err)
	}
	if err = trajectory.Reconcile(ctx, f.root, "idea-x", f.runID, reconciliation.SHA256()); err != nil {
		t.Fatal(err)
	}
}
