package app

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"parley-deck-cli/internal/budget"
	"parley-deck-cli/internal/trajectory"
)

func parentRecoveryBytes(t *testing.T, path string) []byte {
	t.Helper()
	raw, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	return raw
}
func parentRecoveryLedger(t *testing.T, root string) string {
	t.Helper()
	b, err := budget.LoadCycleBinding(context.Background(), root, "idea-x", budget.Fixup)
	if err != nil || b == nil {
		t.Fatal("missing charge", err)
	}
	return filepath.Join(b.Store.Dir, "ledger.json")
}

func TestTrajectoryParentRecoveryRetainsOriginalAndDoesNotExecute(t *testing.T) {
	binary := trajectoryHelperBinary(t)
	for _, mode := range []string{"missing", "incomplete", "empty-directory", "failed-publication"} {
		t.Run(mode, func(t *testing.T) {
			fixtureMode := "actual-helper"
			if mode == "empty-directory" {
				fixtureMode = "parent-publication"
			}
			root, trace, _, agent := trajectoryHelperFixture(t, fixtureMode)
			result, err := verifyTrajectoryWithAgent(context.Background(), root, "idea-x", agent, 60*time.Second, binary, io.Discard)
			if (err == nil) != (mode != "empty-directory") {
				t.Fatalf("unexpected original result: %+v %v", result, err)
			}
			parent := filepath.Join(filepath.Dir(result.RequestPath), "parent-result.json")
			var original []byte
			switch mode {
			case "missing":
				if err = os.Remove(parent); err != nil {
					t.Fatal(err)
				}
			case "incomplete":
				original = []byte(`{"version":1,"run_id":`)
				if err = os.WriteFile(parent, original, 0600); err != nil {
					t.Fatal(err)
				}
			case "failed-publication":
				failed := result
				failed.FailureStage = "parent-publication"
				failed.Assessment = nil
				original, _ = json.MarshalIndent(failed, "", "  ")
				if err = os.WriteFile(parent, original, 0600); err != nil {
					t.Fatal(err)
				}
			}
			ledger := parentRecoveryLedger(t, root)
			charged := parentRecoveryBytes(t, ledger)
			executions := parentRecoveryBytes(t, trace)
			starts := parentRecoveryBytes(t, filepath.Join(root, ".parley-runtime/verifier-starts"))
			if _, err = trajectory.PreviewReconciliation(context.Background(), root, "idea-x", result.RunID); err == nil {
				t.Fatal("unpublished parent was accepted before recovery")
			}
			p, err := trajectory.PreviewParentRecovery(context.Background(), root, "idea-x", result.RunID)
			if err != nil || p.Original.Kind != mode || p.Derived.Result.Assessment == nil || p.Derived.Result.Assessment.Outcome != trajectory.Regression {
				t.Fatalf("recovery preview %+v %v", p, err)
			}
			args := []string{"--dir", root, "--idea", "idea-x", "--run", result.RunID, "--sha256", p.SHA256(), "--yes"}
			// Publication survives an output failure. Exact replay must neither rewrite
			// its record nor execute the paid helper/criterion work again.
			if code := runTrajectoryParentRecovery(context.Background(), args, migrationOutputFailure{}, io.Discard); code != 1 {
				t.Fatal("lost output was ignored")
			}
			recoveryPath := filepath.Join(filepath.Dir(parent), "parent-recovery.json")
			record := parentRecoveryBytes(t, recoveryPath)
			var out, stderr bytes.Buffer
			if code := runTrajectoryParentRecovery(context.Background(), args, &out, &stderr); code != 0 {
				t.Fatalf("exact recovery replay: %d %s", code, &stderr)
			}
			var replay struct {
				Preview trajectory.ParentRecoveryPreview `json:"preview"`
				SHA256  string                           `json:"sha256"`
				Applied bool                             `json:"applied"`
			}
			if err = json.Unmarshal(out.Bytes(), &replay); err != nil || !replay.Applied || replay.SHA256 != p.SHA256() || !sameVerificationJSON(replay.Preview, p) {
				t.Fatal("replay lost original decision", err)
			}
			history, err := trajectory.InspectHistory(context.Background(), root, "idea-x")
			if err != nil || history.Unreconciled != 1 {
				t.Fatal("recovery auto-reconciled charged history", err)
			}
			rp, err := trajectory.PreviewReconciliation(context.Background(), root, "idea-x", result.RunID)
			if err != nil || rp.RecoverySHA256 == "" {
				t.Fatal("reconciliation omitted recovery provenance", err)
			}
			if err = trajectory.Reconcile(context.Background(), root, "idea-x", result.RunID, rp.SHA256()); err != nil {
				t.Fatal(err)
			}
			history, err = trajectory.InspectHistory(context.Background(), root, "idea-x")
			if err != nil || history.Unreconciled != 0 || len(history.Decision.Assessments) != 1 || history.Decision.Assessments[0].Outcome != trajectory.Regression {
				t.Fatal("recovered primary observations were not reconciled", err)
			}
			later, err := trajectory.RecoverParent(context.Background(), root, "idea-x", result.RunID, p.SHA256())
			if err != nil || !sameVerificationJSON(later, p) {
				t.Fatal("later reconciliation changed exact recovery replay", err)
			}
			if !bytes.Equal(record, parentRecoveryBytes(t, recoveryPath)) || !bytes.Equal(charged, parentRecoveryBytes(t, ledger)) || !bytes.Equal(executions, parentRecoveryBytes(t, trace)) || !bytes.Equal(starts, parentRecoveryBytes(t, filepath.Join(root, ".parley-runtime/verifier-starts"))) {
				t.Fatal("recovery/replay changed original charges, record or execution count")
			}
			switch mode {
			case "missing":
				if _, err = os.Lstat(parent); !os.IsNotExist(err) {
					t.Fatal("recovery synthesized the original parent file")
				}
			case "empty-directory":
				entries, err := os.ReadDir(parent)
				if err != nil || len(entries) != 0 {
					t.Fatal("original failed publication directory was changed", err)
				}
			default:
				if !bytes.Equal(original, parentRecoveryBytes(t, parent)) {
					t.Fatal("original failed/incomplete parent was overwritten")
				}
			}
			// A new resolution pins its recovery record even if someone copies the
			// derived success into the original parent path after removing provenance.
			if err = os.Remove(recoveryPath); err != nil {
				t.Fatal(err)
			}
			if _, err = trajectory.InspectHistory(context.Background(), root, "idea-x"); err == nil {
				t.Fatal("accepted recovery history disappeared silently")
			}
		})
	}
}

func TestTrajectoryParentRecoveryRestoresPreviouslyBoundFacts(t *testing.T) {
	binary := trajectoryHelperBinary(t)
	root, trace, _, agent := trajectoryHelperFixture(t, "actual-helper")
	result, old := reconciledHelper(t, root, binary, agent)
	before, err := trajectory.Inspect(context.Background(), root, "idea-x")
	if err != nil {
		t.Fatal(err)
	}
	parent := filepath.Join(filepath.Dir(result.RequestPath), "parent-result.json")
	original := parentRecoveryBytes(t, parent)
	if err = os.Remove(parent); err != nil {
		t.Fatal(err)
	}
	if _, err = trajectory.InspectHistory(context.Background(), root, "idea-x"); err == nil {
		t.Fatal("lost previously accepted parent remained trusted")
	}
	p, err := trajectory.PreviewParentRecovery(context.Background(), root, "idea-x", result.RunID)
	if err != nil || p.ParentSHA256 != old.ParentSHA256 {
		t.Fatal("recovery failed to reproduce previously pinned bytes", err)
	}
	executions := parentRecoveryBytes(t, trace)
	if _, err = trajectory.RecoverParent(context.Background(), root, "idea-x", result.RunID, p.SHA256()); err != nil {
		t.Fatal(err)
	}
	after, err := trajectory.Inspect(context.Background(), root, "idea-x")
	if err != nil || !sameVerificationJSON(before, after) || !bytes.Equal(executions, parentRecoveryBytes(t, trace)) {
		t.Fatal("recovery rewrote existing resolution or executed checks", err)
	}
	// A late writer must not silently replace the frozen missing-parent state.
	if err = os.WriteFile(parent, original, 0600); err != nil {
		t.Fatal(err)
	}
	if _, err = trajectory.InspectHistory(context.Background(), root, "idea-x"); err == nil {
		t.Fatal("late original publication overrode recovery provenance")
	}
}

func TestTrajectoryParentRecoveryRejectsChangedPrimaryEvidence(t *testing.T) {
	binary := trajectoryHelperBinary(t)
	root, trace, journal, agent := trajectoryHelperFixture(t, "actual-helper")
	result, err := verifyTrajectoryWithAgent(context.Background(), root, "idea-x", agent, 60*time.Second, binary, io.Discard)
	if err != nil {
		t.Fatal(err)
	}
	parent := filepath.Join(filepath.Dir(result.RequestPath), "parent-result.json")
	if err = os.Remove(parent); err != nil {
		t.Fatal(err)
	}
	p, err := trajectory.PreviewParentRecovery(context.Background(), root, "idea-x", result.RunID)
	if err != nil {
		t.Fatal(err)
	}
	executions := parentRecoveryBytes(t, trace)
	inv := filepath.Join(root, ".parley-runtime/invocations", result.InvocationID)
	paths := []string{result.RequestPath, filepath.Join(inv, "requested.json"), filepath.Join(inv, "started.json"), filepath.Join(inv, "terminal.json"), filepath.Join(journal, "request.json"), filepath.Join(journal, "launch.json"), filepath.Join(journal, "receipt.json"), filepath.Join(journal, "step-001.json"), filepath.Join(journal, "process-001.json")}
	for _, path := range paths {
		t.Run(filepath.Base(filepath.Dir(path))+"/"+filepath.Base(path), func(t *testing.T) {
			raw := parentRecoveryBytes(t, path)
			if err := os.WriteFile(path, append(raw, '\n'), 0600); err != nil {
				t.Fatal(err)
			}
			defer os.WriteFile(path, raw, 0600)
			if _, err := trajectory.RecoverParent(context.Background(), root, "idea-x", result.RunID, p.SHA256()); err == nil {
				t.Fatal("changed primary evidence passed original recovery preview")
			}
			if _, err := os.Lstat(filepath.Join(filepath.Dir(parent), "parent-recovery.json")); !os.IsNotExist(err) {
				t.Fatal("refused recovery published an observation")
			}
		})
	}
	if !bytes.Equal(executions, parentRecoveryBytes(t, trace)) {
		t.Fatal("evidence refusal reran the helper")
	}
	if _, err = trajectory.RecoverParent(context.Background(), root, "idea-x", result.RunID, strings.Repeat("0", 64)); err == nil {
		t.Fatal("wrong preview digest was accepted")
	}
}

func TestTrajectoryParentRecoveryConcurrentCLI(t *testing.T) {
	binary := trajectoryHelperBinary(t)
	root, trace, _, agent := trajectoryHelperFixture(t, "actual-helper")
	result, err := verifyTrajectoryWithAgent(context.Background(), root, "idea-x", agent, 60*time.Second, binary, io.Discard)
	if err != nil {
		t.Fatal(err)
	}
	if err = os.Remove(filepath.Join(filepath.Dir(result.RequestPath), "parent-result.json")); err != nil {
		t.Fatal(err)
	}
	args := []string{"trajectory", "recover-parent", "--dir", root, "--idea", "idea-x", "--run", result.RunID}
	raw, err := exec.Command(binary, args...).CombinedOutput()
	if err != nil {
		t.Fatalf("fresh CLI recovery preview: %v %s", err, raw)
	}
	var preview struct {
		Preview trajectory.ParentRecoveryPreview `json:"preview"`
		SHA256  string                           `json:"sha256"`
		Applied bool                             `json:"applied"`
	}
	if err = json.Unmarshal(raw, &preview); err != nil || preview.Applied || preview.SHA256 != preview.Preview.SHA256() {
		t.Fatal("invalid CLI preview", err)
	}
	executions := parentRecoveryBytes(t, trace)
	apply := append(args, "--sha256", preview.SHA256, "--yes")
	cmds := []*exec.Cmd{exec.Command(binary, apply...), exec.Command(binary, apply...)}
	output := make([]bytes.Buffer, 2)
	for i, c := range cmds {
		c.Stdout = &output[i]
		c.Stderr = &output[i]
		if err = c.Start(); err != nil {
			t.Fatal(err)
		}
	}
	for i, c := range cmds {
		if err = c.Wait(); err != nil {
			t.Fatalf("fresh concurrent recovery %v %s", err, &output[i])
		}
	}
	if !bytes.Equal(output[0].Bytes(), output[1].Bytes()) || !bytes.Equal(executions, parentRecoveryBytes(t, trace)) {
		t.Fatal("competing replay changed output or execution count")
	}
	if _, err = trajectory.PreviewReconciliation(context.Background(), root, "idea-x", result.RunID); err != nil {
		t.Fatal(err)
	}
}

func TestTrajectoryParentRecoveryNeverTurnsFailedExecutionIntoEvidence(t *testing.T) {
	binary := trajectoryHelperBinary(t)
	for _, mode := range []string{"text-only", "failed-after-helper", "missing-receipt"} {
		t.Run(mode, func(t *testing.T) {
			root, _, _, agent := trajectoryHelperFixture(t, mode)
			result, err := verifyTrajectoryWithAgent(context.Background(), root, "idea-x", agent, 60*time.Second, binary, io.Discard)
			if err == nil {
				t.Fatal("fixture unexpectedly succeeded")
			}
			if err = os.Remove(filepath.Join(filepath.Dir(result.RequestPath), "parent-result.json")); err != nil {
				t.Fatal(err)
			}
			if _, err = trajectory.PreviewParentRecovery(context.Background(), root, "idea-x", result.RunID); err == nil {
				t.Fatal("failed or incomplete independent execution became recoverable success")
			}
		})
	}
}

func TestTrajectoryParentRecoveryPinsScopeAndProvenanceAcrossLaterWork(t *testing.T) {
	binary := trajectoryHelperBinary(t)
	root, trace, _, agent := trajectoryHelperFixture(t, "actual-helper")
	r, err := verifyTrajectoryWithAgent(context.Background(), root, "idea-x", agent, 60*time.Second, binary, io.Discard)
	if err != nil {
		t.Fatal(err)
	}
	parent := filepath.Join(filepath.Dir(r.RequestPath), "parent-result.json")
	original := parentRecoveryBytes(t, parent)
	if err = os.Remove(parent); err != nil {
		t.Fatal(err)
	}
	p, err := trajectory.PreviewParentRecovery(context.Background(), root, "idea-x", r.RunID)
	if err != nil {
		t.Fatal(err)
	}
	scopePath := filepath.Join(root, "parley-deck/ideas/idea-x/00-prompt.md")
	scope := parentRecoveryBytes(t, scopePath)
	for _, replacement := range [][2]string{{"[builder, reviewer]", "[reviewer, builder]"}, {"name: material", "name: changed"}, {"cat source", "cat other"}} {
		changed := bytes.ReplaceAll(scope, []byte(replacement[0]), []byte(replacement[1]))
		if bytes.Equal(changed, scope) {
			t.Fatal("scope mutation fixture did not change scope")
		}
		if err = os.WriteFile(scopePath, changed, 0600); err != nil {
			t.Fatal(err)
		}
		_, refusal := trajectory.RecoverParent(context.Background(), root, "idea-x", r.RunID, p.SHA256())
		if err = os.WriteFile(scopePath, scope, 0600); err != nil {
			t.Fatal(err)
		}
		if refusal == nil {
			t.Fatal("changed original scope accepted recovery")
		}
	}
	if _, err = trajectory.RecoverParent(context.Background(), root, "idea-x", r.RunID, p.SHA256()); err != nil {
		t.Fatal(err)
	}
	rp, err := trajectory.PreviewReconciliation(context.Background(), root, "idea-x", r.RunID)
	if err != nil {
		t.Fatal(err)
	}
	if err = trajectory.Reconcile(context.Background(), root, "idea-x", r.RunID, rp.SHA256()); err != nil {
		t.Fatal(err)
	}
	commitTrajectorySource(t, root)
	continueTrajectory(t, root, "recover-then-continue", false, false)
	if started, err := measuredTrajectoryPatch(t, root, "later-correction", "original"); err != nil || !started {
		t.Fatal("later charged attempt refused", err)
	}
	_, second := reconciledHelper(t, root, binary, agent)
	if second.Assessment.Outcome != trajectory.NoRegression {
		t.Fatal("later correction was not observed")
	}
	state, err := trajectory.Inspect(context.Background(), root, "idea-x")
	if err != nil || len(state.Attempts) != 2 {
		t.Fatal("lost later history", err)
	}
	executions := parentRecoveryBytes(t, trace)
	ledger := parentRecoveryLedger(t, root)
	charges := parentRecoveryBytes(t, ledger)
	replay, err := trajectory.RecoverParent(context.Background(), root, "idea-x", r.RunID, p.SHA256())
	if err != nil || !sameVerificationJSON(replay, p) {
		t.Fatal("later legitimate work invalidated original exact replay", err)
	}
	if _, err = trajectory.RecoverParent(context.Background(), root, "idea-x", r.RunID, strings.Repeat("0", 64)); err == nil {
		t.Fatal("changed retained recovery preview accepted")
	}
	recordPath := filepath.Join(filepath.Dir(parent), "parent-recovery.json")
	record := parentRecoveryBytes(t, recordPath)
	var recovery trajectory.ParentRecovery
	if err = json.Unmarshal(record, &recovery); err != nil {
		t.Fatal(err)
	}
	// Even a structurally valid timestamp edit cannot change the record pinned
	// by an accepted resolution. Primary evidence remains untouched.
	recovery.RecoveredAt = recovery.RecoveredAt.Add(time.Second)
	changed, _ := json.MarshalIndent(recovery, "", "  ")
	if err = os.WriteFile(recordPath, append(changed, '\n'), 0600); err != nil {
		t.Fatal(err)
	}
	if _, err = trajectory.InspectHistory(context.Background(), root, "idea-x"); err == nil {
		t.Fatal("mutated accepted recovery provenance stayed trusted")
	}
	if _, err = trajectory.RecoverParent(context.Background(), root, "idea-x", r.RunID, p.SHA256()); err == nil {
		t.Fatal("replay repaired mutated accepted provenance")
	}
	if err = os.Remove(recordPath); err != nil {
		t.Fatal(err)
	}
	if err = os.WriteFile(parent, original, 0600); err != nil {
		t.Fatal(err)
	}
	if _, err = trajectory.InspectHistory(context.Background(), root, "idea-x"); err == nil {
		t.Fatal("copied parent replaced the accepted recovery record")
	}
	if err = os.Remove(parent); err != nil {
		t.Fatal(err)
	}
	if err = os.WriteFile(recordPath, record, 0600); err != nil {
		t.Fatal(err)
	}
	after, err := trajectory.Inspect(context.Background(), root, "idea-x")
	if err != nil || !sameVerificationJSON(after, state) || !bytes.Equal(executions, parentRecoveryBytes(t, trace)) || !bytes.Equal(charges, parentRecoveryBytes(t, ledger)) {
		t.Fatal("replay or provenance refusal mutated history or executed work", err)
	}
}

func TestTrajectoryParentRecoveryRejectsStalePreviewAfterOriginalReconciliation(t *testing.T) {
	binary := trajectoryHelperBinary(t)
	root, _, _, agent := trajectoryHelperFixture(t, "actual-helper")
	r, err := verifyTrajectoryWithAgent(context.Background(), root, "idea-x", agent, 60*time.Second, binary, io.Discard)
	if err != nil {
		t.Fatal(err)
	}
	parent := filepath.Join(filepath.Dir(r.RequestPath), "parent-result.json")
	// The alternate supported encoding must keep its previously pinned digest.
	original := append(parentRecoveryBytes(t, parent), '\n')
	if err = os.Remove(parent); err != nil {
		t.Fatal(err)
	}
	p, err := trajectory.PreviewParentRecovery(context.Background(), root, "idea-x", r.RunID)
	if err != nil {
		t.Fatal(err)
	}
	if err = os.WriteFile(parent, original, 0600); err != nil {
		t.Fatal(err)
	}
	rp, err := trajectory.PreviewReconciliation(context.Background(), root, "idea-x", r.RunID)
	if err != nil {
		t.Fatal(err)
	}
	if err = trajectory.Reconcile(context.Background(), root, "idea-x", r.RunID, rp.SHA256()); err != nil {
		t.Fatal(err)
	}
	if err = os.Remove(parent); err != nil {
		t.Fatal(err)
	}
	if _, err = trajectory.RecoverParent(context.Background(), root, "idea-x", r.RunID, p.SHA256()); err == nil {
		t.Fatal("stale preview applied after trajectory state changed")
	}
	fresh, err := trajectory.PreviewParentRecovery(context.Background(), root, "idea-x", r.RunID)
	if err != nil || fresh.ParentSHA256 != rp.ParentSHA256 || fresh.StateSHA256 == p.StateSHA256 {
		t.Fatal("fresh preview lost original alternate parent digest or changed state", err)
	}
	if _, err = trajectory.RecoverParent(context.Background(), root, "idea-x", r.RunID, fresh.SHA256()); err != nil {
		t.Fatal(err)
	}
	if _, err = trajectory.InspectHistory(context.Background(), root, "idea-x"); err != nil {
		t.Fatal(err)
	}
}
