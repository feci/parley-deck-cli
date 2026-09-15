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

	"parley-deck-cli/internal/agents"
	"parley-deck-cli/internal/budget"
	"parley-deck-cli/internal/runner"
	"parley-deck-cli/internal/trajectory"
)

func reconciledHelper(t *testing.T, root, binary string, agent agents.Discovery) (trajectoryVerificationResult, trajectory.ReconciliationPreview) {
	t.Helper()
	r, err := verifyTrajectoryWithAgent(context.Background(), root, "idea-x", agent, 60*time.Second, binary, io.Discard)
	if err != nil {
		t.Fatalf("independent helper: %+v %v", r, err)
	}
	p, err := trajectory.PreviewReconciliation(context.Background(), root, "idea-x", r.RunID)
	if err != nil {
		t.Fatal("reconciliation preview", err)
	}
	if err = trajectory.Reconcile(context.Background(), root, "idea-x", r.RunID, p.SHA256()); err != nil {
		t.Fatal("reconciliation apply", err)
	}
	return r, p
}

func commitTrajectorySource(t *testing.T, root string) {
	t.Helper()
	gateGit(t, root, "add", "source")
	gateGit(t, root, "-c", "user.email=t@t", "-c", "user.name=t", "commit", "-qm", "Test-owned source promotion")
}

func continueTrajectory(t *testing.T, root, id string, review, inconclusive bool) {
	t.Helper()
	p, err := trajectory.PreviewContinuation(context.Background(), root, "idea-x")
	if err != nil {
		t.Fatal(err)
	}
	args := []string{"continue", "--dir", root, "--idea", "idea-x", "--sha256", p.SHA256(), "--decision-id", id, "--reason", "test-owned explicit continuation", "--yes"}
	if review {
		args = append(args, "--acknowledge-review")
	}
	if inconclusive {
		args = append(args, "--acknowledge-inconclusive")
	}
	var out, stderr bytes.Buffer
	for i := 0; i < 2; i++ {
		if code := runTrajectoryReconciliation(context.Background(), args, &out, &stderr, true); code != 0 {
			t.Fatalf("attended fixture continuation/replay: %d %s", code, &stderr)
		}
	}
}

func measuredTrajectoryPatch(t *testing.T, root, runID, source string) (bool, error) {
	t.Helper()
	command := "printf '%s\\n' '" + strings.ReplaceAll(source, "'", "'\\''") + "' > source"
	agent := agents.Discovery{Spec: agents.Spec{ID: "builder", PromptMode: agents.PromptStdin, LaunchMode: agents.LaunchHeadless, HeadlessArgs: []string{"-c", command}}, Found: true, Path: "/bin/sh"}
	r, err := runner.RunMeasured(context.Background(), runner.ExecOptions{Root: root, Agent: agent, Prompt: "synthetic next charged patch", Timeout: 30 * time.Second, Info: runner.LaunchInfo{RunID: runID, Idea: "idea-x", Phase: "fixup"}})
	return r.StartedAt != nil, err
}

func TestTrajectoryReconciliationRunsWholeChargedHistoryAndRetainsReview(t *testing.T) {
	binary := trajectoryHelperBinary(t)
	secondary := trajectory.Criterion{Name: "secondary", Command: `if [ "$(cat source)" = broken2 ]; then printf 'PARLEY-EVIDENCE {"executed_cases":1,"failed_cases":1}\n'; exit 1; else printf 'PARLEY-EVIDENCE {"executed_cases":1,"failed_cases":0}\n'; fi`}
	root, _, _, agent := trajectoryHelperFixtureCriteria(t, "actual-helper", []trajectory.Criterion{secondary})
	_, first := reconciledHelper(t, root, binary, agent)
	if first.Assessment.Outcome != trajectory.Regression {
		t.Fatal("first regression absent")
	}
	s, err := trajectory.Inspect(context.Background(), root, "idea-x")
	if err != nil {
		t.Fatal(err)
	}
	original, _ := json.Marshal(s.Attempts[0])
	if started, err := measuredTrajectoryPatch(t, root, "dirty-refusal", "forbidden"); err == nil || started {
		t.Fatal("dirty before-state silently promoted")
	}
	commitTrajectorySource(t, root)
	continueTrajectory(t, root, "promote-first", false, false)
	if started, err := measuredTrajectoryPatch(t, root, "second-patch", "broken2"); err != nil || !started {
		t.Fatalf("second charged patch refused: %v", err)
	}
	_, second := reconciledHelper(t, root, binary, agent)
	if second.Assessment.Outcome != trajectory.Regression {
		t.Fatal("new secondary regression absent")
	}
	h, err := trajectory.InspectHistory(context.Background(), root, "idea-x")
	if err != nil || !h.ReviewPending || h.Decision.TriggerSequence != 2 || h.Decision.Consecutive != 2 || len(h.Decision.Assessments) != 2 {
		t.Fatalf("full ordered history did not trigger review: %+v %v", h, err)
	}
	if started, err := measuredTrajectoryPatch(t, root, "review-refusal", "forbidden"); err == nil || started {
		t.Fatal("review gate allowed another model patch")
	}
	commitTrajectorySource(t, root)
	p, err := trajectory.PreviewContinuation(context.Background(), root, "idea-x")
	if err != nil {
		t.Fatal(err)
	}
	if _, err = trajectory.Continue(context.Background(), root, "idea-x", p.SHA256(), "missing-ack", "test-owned missing acknowledgment", false, false); err == nil {
		t.Fatal("review was implicitly acknowledged")
	}
	continueTrajectory(t, root, "review-second", true, false)
	if started, err := measuredTrajectoryPatch(t, root, "corrective-patch", "original"); err != nil || !started {
		t.Fatalf("explicit continuation did not admit the corrective patch: %v", err)
	}
	_, third := reconciledHelper(t, root, binary, agent)
	if third.Assessment.Outcome != trajectory.NoRegression {
		t.Fatal("corrected material outcome not derived")
	}
	h, err = trajectory.InspectHistory(context.Background(), root, "idea-x")
	if err != nil || h.ReviewPending || !h.Decision.ReviewRequired || h.Decision.TriggerSequence != 2 || h.Decision.Consecutive != 0 || len(h.Decision.Assessments) != 3 {
		t.Fatalf("clean patch erased trigger or charged coverage: %+v %v", h, err)
	}
	if err = trajectory.RequireResolved(context.Background(), root, "idea-x"); err != nil {
		t.Fatal("clean reconciled current source still blocked", err)
	}
	promptPath := filepath.Join(root, "parley-deck", "ideas", "idea-x", "00-prompt.md")
	prompt, err := os.ReadFile(promptPath)
	if err != nil {
		t.Fatal(err)
	}
	for _, replacement := range [][2]string{{"name: material", "name: changed"}, {"[builder, reviewer]", "[builder]"}} {
		changed := bytes.ReplaceAll(prompt, []byte(replacement[0]), []byte(replacement[1]))
		if err = os.WriteFile(promptPath, changed, 0600); err != nil {
			t.Fatal(err)
		}
		_, readErr := trajectory.InspectHistory(context.Background(), root, "idea-x")
		if err = os.WriteFile(promptPath, prompt, 0600); err != nil {
			t.Fatal(err)
		}
		if readErr == nil {
			t.Fatal("changed current criterion/quorum was accepted after reconciliation")
		}
	}
	s, err = trajectory.Inspect(context.Background(), root, "idea-x")
	if err != nil {
		t.Fatal(err)
	}
	retained, _ := json.Marshal(s.Attempts[0])
	if !bytes.Equal(retained, original) || s.Attempts[0].After.Clean || len(s.Attempts) != 3 || len(s.Resolutions) != 3 {
		t.Fatal("original dirty charged history was rewritten")
	}
	b, err := budget.LoadCycleBinding(context.Background(), root, "idea-x", budget.Fixup)
	if err != nil {
		t.Fatal(err)
	}
	ledger, err := b.Store.Inspect(context.Background())
	if err != nil || b.Count(ledger) != 3 {
		t.Fatal("refusals reset or duplicated fixup charges", err)
	}
	// A fresh CLI must reconstruct the complete history from retained evidence.
	cmd := exec.Command(binary, "trajectory", "history", "--dir", root, "--idea", "idea-x")
	raw, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("fresh history CLI: %v %s", err, raw)
	}
	var fresh trajectory.History
	if err = json.Unmarshal(raw, &fresh); err != nil || !sameVerificationJSON(fresh, h) {
		t.Fatal("fresh process changed history", err)
	}
	// Replaying an older decision must return that original receipt even when
	// the latest attempt has a different source and pending promotion preview.
	old := s.Continuations[0]
	beforeReplay, _ := json.Marshal(s)
	var replayOut, replayErr bytes.Buffer
	args := []string{"continue", "--dir", root, "--idea", "idea-x", "--sha256", old.SHA256, "--decision-id", old.DecisionID, "--reason", old.Reason, "--yes"}
	if code := runTrajectoryReconciliation(context.Background(), args, &replayOut, &replayErr, true); code != 0 {
		t.Fatalf("historical continuation replay consulted a later attempt: %d %s", code, &replayErr)
	}
	var replay struct {
		Preview trajectory.ContinuationPreview `json:"preview"`
		SHA256  string                         `json:"sha256"`
		Applied bool                           `json:"applied"`
	}
	if err = json.Unmarshal(replayOut.Bytes(), &replay); err != nil || !replay.Applied || replay.SHA256 != old.SHA256 || !sameVerificationJSON(replay.Preview, old.Preview) {
		t.Fatal("historical continuation replay returned another decision", err)
	}
	afterReplay, err := trajectory.Inspect(context.Background(), root, "idea-x")
	if err != nil {
		t.Fatal(err)
	}
	afterReplayBytes, _ := json.Marshal(afterReplay)
	if !bytes.Equal(beforeReplay, afterReplayBytes) {
		t.Fatal("historical replay changed charged state")
	}
	statePath := filepath.Join(filepath.Dir(b.Store.Dir), "trajectory.json")
	stateBytes, err := os.ReadFile(statePath)
	if err != nil {
		t.Fatal(err)
	}
	for _, kind := range []string{"omit-resolution", "erase-review-ack", "erase-promotion", "rewrite-dirty-after", "downgrade-version"} {
		t.Run(kind, func(t *testing.T) {
			var changed trajectory.State
			if err := json.Unmarshal(stateBytes, &changed); err != nil {
				t.Fatal(err)
			}
			switch kind {
			case "omit-resolution":
				changed.Resolutions = append(changed.Resolutions[:1], changed.Resolutions[2:]...)
			case "erase-review-ack":
				changed.Continuations[1].ReviewThrough = 0
			case "erase-promotion":
				changed.Continuations = changed.Continuations[1:]
			case "rewrite-dirty-after":
				changed.Attempts[0].After.Clean = true
			case "downgrade-version":
				changed.Version = 2
			}
			raw, err := json.MarshalIndent(changed, "", "  ")
			if err != nil {
				t.Fatal(err)
			}
			if err = os.WriteFile(statePath, append(raw, '\n'), 0600); err != nil {
				t.Fatal(err)
			}
			defer os.WriteFile(statePath, stateBytes, 0600)
			if _, err = trajectory.InspectHistory(context.Background(), root, "idea-x"); err == nil {
				t.Fatal("rewritten charged history accepted")
			}
			if err = trajectory.RequireResolved(context.Background(), root, "idea-x"); err == nil {
				t.Fatal("rewritten history admitted completion")
			}
		})
	}
}

func TestTrajectoryReconciliationRechecksEvidenceAndExactReplay(t *testing.T) {
	binary := trajectoryHelperBinary(t)
	root, _, journal, agent := trajectoryHelperFixture(t, "actual-helper")
	r, p := reconciledHelper(t, root, binary, agent)
	s, err := trajectory.Inspect(context.Background(), root, "idea-x")
	if err != nil {
		t.Fatal(err)
	}
	original, _ := json.Marshal(s)
	if err = trajectory.Reconcile(context.Background(), root, "idea-x", r.RunID, p.SHA256()); err != nil {
		t.Fatal("exact reconciliation replay failed", err)
	}
	if err = trajectory.Reconcile(context.Background(), root, "idea-x", r.RunID, strings.Repeat("0", 64)); err == nil {
		t.Fatal("changed reconciliation replay accepted")
	}
	for _, path := range []string{
		filepath.Join(filepath.Dir(r.RequestPath), "parent-result.json"), r.RequestPath,
		filepath.Join(root, ".parley-runtime", "invocations", r.InvocationID, "terminal.json"),
		filepath.Join(journal, "receipt.json"), filepath.Join(journal, "process-001.json"), filepath.Join(journal, "step-001.json"),
	} {
		t.Run(filepath.Base(path), func(t *testing.T) {
			bytesBefore, err := os.ReadFile(path)
			if err != nil {
				t.Fatal(err)
			}
			if err = os.WriteFile(path, append(bytesBefore, '\n'), 0600); err != nil {
				t.Fatal(err)
			}
			defer os.WriteFile(path, bytesBefore, 0600)
			if _, err = trajectory.InspectHistory(context.Background(), root, "idea-x"); err == nil {
				t.Fatal("changed accepted evidence became trusted cached state")
			}
			if _, err = trajectory.PreviewContinuation(context.Background(), root, "idea-x"); err == nil {
				t.Fatal("lost evidence allowed continuation preview")
			}
		})
	}
	s, err = trajectory.Inspect(context.Background(), root, "idea-x")
	if err != nil {
		t.Fatal(err)
	}
	after, _ := json.Marshal(s)
	if !bytes.Equal(original, after) {
		t.Fatal("replay or failed read rewrote state")
	}
	commitTrajectorySource(t, root)
	preview, err := trajectory.PreviewContinuation(context.Background(), root, "idea-x")
	if err != nil {
		t.Fatal(err)
	}
	args := []string{"continue", "--dir", root, "--idea", "idea-x", "--sha256", preview.SHA256(), "--decision-id", "operator-only", "--reason", "fixture", "--yes"}
	if code := runTrajectoryReconciliation(context.Background(), args, io.Discard, io.Discard, false); code != 2 {
		t.Fatal("unattended process issued continuation")
	}
	if err = os.WriteFile(filepath.Join(root, "source"), []byte("unrecorded\n"), 0600); err != nil {
		t.Fatal(err)
	}
	if _, err = trajectory.Continue(context.Background(), root, "idea-x", preview.SHA256(), "changed-source", "fixture", false, false); err == nil {
		t.Fatal("stale promotion accepted changed material source")
	}
}

func TestTrajectoryReconciliationRetainsInconclusiveAndNeedsExplicitDecision(t *testing.T) {
	binary := trajectoryHelperBinary(t)
	root, _, _, agent := trajectoryHelperFixture(t, "inconclusive")
	_, p := reconciledHelper(t, root, binary, agent)
	if p.Assessment.Outcome != trajectory.Inconclusive {
		t.Fatalf("fixture did not produce inconclusive paired evidence: %+v", p.Assessment)
	}
	commitTrajectorySource(t, root)
	if started, err := measuredTrajectoryPatch(t, root, "inconclusive-refusal", "forbidden"); err == nil || started {
		t.Fatal("inconclusive comparison granted another patch")
	}
	continueTrajectory(t, root, "retain-inconclusive", false, true)
	h, err := trajectory.InspectHistory(context.Background(), root, "idea-x")
	if err != nil || len(h.InconclusivePending) != 0 || !h.Decision.Pending || h.Decision.Assessments[0].Outcome != trajectory.Inconclusive {
		t.Fatal("explicit decision rewrote uncertain evidence", err)
	}
	if err = trajectory.RequireResolved(context.Background(), root, "idea-x"); err == nil {
		t.Fatal("explicit continuation became completion")
	}
}

func TestTrajectoryReconciliationCLIFromFreshProcesses(t *testing.T) {
	binary := trajectoryHelperBinary(t)
	root, _, _, agent := trajectoryHelperFixture(t, "actual-helper")
	r, err := verifyTrajectoryWithAgent(context.Background(), root, "idea-x", agent, 60*time.Second, binary, io.Discard)
	if err != nil {
		t.Fatal(err)
	}
	args := []string{"trajectory", "reconcile", "--dir", root, "--idea", "idea-x", "--run", r.RunID}
	raw, err := exec.Command(binary, args...).CombinedOutput()
	if err != nil {
		t.Fatalf("fresh preview failed: %v %s", err, raw)
	}
	var preview struct {
		Preview trajectory.ReconciliationPreview `json:"preview"`
		SHA256  string                           `json:"sha256"`
		Applied bool                             `json:"applied"`
	}
	if err = json.Unmarshal(raw, &preview); err != nil || preview.Applied || preview.SHA256 != preview.Preview.SHA256() {
		t.Fatal("invalid fresh CLI preview", err)
	}
	apply := append(append([]string{}, args...), "--sha256", preview.SHA256, "--yes")
	type outcome struct {
		raw []byte
		err error
	}
	done := make(chan outcome, 2)
	for i := 0; i < 2; i++ {
		go func() { raw, err := exec.Command(binary, apply...).CombinedOutput(); done <- outcome{raw, err} }()
	}
	var failures []string
	for i := 0; i < 2; i++ {
		result := <-done
		if result.err != nil {
			failures = append(failures, result.err.Error()+": "+string(result.raw))
		}
	}
	if len(failures) > 0 {
		t.Fatalf("competing fresh apply/replay failed: %s", strings.Join(failures, "\n"))
	}
	s, err := trajectory.Inspect(context.Background(), root, "idea-x")
	if err != nil || len(s.Resolutions) != 1 {
		t.Fatal("CLI duplicated or lost reconciliation", err)
	}
	commitTrajectorySource(t, root)
	p, err := trajectory.PreviewContinuation(context.Background(), root, "idea-x")
	if err != nil {
		t.Fatal(err)
	}
	args = []string{"trajectory", "continue", "--dir", root, "--idea", "idea-x", "--sha256", p.SHA256(), "--decision-id", "unattended-cli", "--reason", "synthetic process", "--yes"}
	raw, err = exec.Command(binary, args...).CombinedOutput()
	if err == nil || !bytes.Contains(raw, []byte("attended operator")) {
		t.Fatalf("fresh unattended CLI issued continuation: %v %s", err, raw)
	}
}
