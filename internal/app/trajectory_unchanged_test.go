package app

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
	"time"

	"parley-deck-cli/internal/agents"
	"parley-deck-cli/internal/budget"
	"parley-deck-cli/internal/protocol"
	"parley-deck-cli/internal/runner"
	"parley-deck-cli/internal/telemetry"
	"parley-deck-cli/internal/trajectory"
)

func unchangedAppFixture(t *testing.T) string {
	t.Helper()
	if runtime.GOOS == "windows" {
		t.Skip("actual POSIX execution required")
	}
	t.Setenv("PARLEY_HOME", t.TempDir())
	t.Setenv("PARLEY_HEADLESS_AGENT_CONFIG", "")
	root, idea := gateScratchRepo(t, "participants: [builder, reviewer]\ntrack: deliberation\nchecks:\n  - name: material\n    command: 'true'\n")
	if err := protocol.InitWorkspace(root); err != nil {
		t.Fatal(err)
	}
	declareAppTestSource(t, root)
	for name, body := range map[string]string{".gitignore": ".parley-runtime/\nparley-deck/runs/\n", "source": "original\n"} {
		if err := os.WriteFile(filepath.Join(root, name), []byte(body), 0600); err != nil {
			t.Fatal(err)
		}
	}
	gateGit(t, root, "add", "-A")
	gateGit(t, root, "-c", "user.email=t@t", "-c", "user.name=t", "commit", "-qm", "Unchanged launch fixture")
	if _, err := budget.EnsureCycleBinding(context.Background(), root, "idea-x", budget.Fixup, 5, 0, "", idea); err != nil {
		t.Fatal(err)
	}
	p, expected, err := trajectory.NewPolicy(context.Background(), root, "idea-x", "builder", []trajectory.Criterion{{Name: "material", Command: "true"}})
	if err != nil {
		t.Fatal(err)
	}
	if err = trajectory.Activate(context.Background(), root, expected, p); err != nil {
		t.Fatal(err)
	}
	return root
}

func measuredUnchanged(t *testing.T, root, command string) telemetry.Record {
	t.Helper()
	agent := agents.Discovery{Spec: agents.Spec{ID: "builder", PromptMode: agents.PromptStdin, LaunchMode: agents.LaunchHeadless, HeadlessArgs: []string{"-c", "printf 'executed\\n' >> .parley-runtime/unchanged-executions; " + command}}, Found: true, Path: "/bin/sh"}
	r, err := runner.RunMeasured(context.Background(), runner.ExecOptions{Root: root, Agent: agent, Prompt: "synthetic charged unchanged fixture", Timeout: 30 * time.Second, Info: runner.LaunchInfo{RunID: "unchanged-fixture", Idea: "idea-x", Phase: "fixup"}})
	if r.StartedAt == nil || r.CompletedAt == nil || r.Outcome == nil || r.Outcome.ExitCode == nil || (*r.Outcome.ExitCode == 0 && err != nil) || (*r.Outcome.ExitCode != 0 && err == nil) {
		t.Fatalf("actual unchanged process unavailable: %+v %v", r, err)
	}
	return r
}

func applyUnchanged(t *testing.T, root string, sequence int) trajectory.ReconciliationPreview {
	t.Helper()
	p, err := trajectory.PreviewUnchanged(context.Background(), root, "idea-x", sequence)
	if err != nil {
		t.Fatal(err)
	}
	result, err := trajectory.ReconcileUnchanged(context.Background(), root, "idea-x", sequence, p.SHA256())
	if err != nil || result.SHA256() != p.SHA256() {
		t.Fatal("unchanged apply", err)
	}
	return p
}

func TestTrajectoryUnchangedCLIOutputFailureAndHistoricalReplay(t *testing.T) {
	ctx := context.Background()
	root := unchangedAppFixture(t)
	r := measuredUnchanged(t, root, "exit 7")
	args := []string{"reconcile-unchanged", "--dir", root, "--idea", "idea-x", "--sequence", "1"}
	var out, errout bytes.Buffer
	if code := runTrajectory(ctx, args, &out, &errout); code != 0 {
		t.Fatal("CLI preview", code, &errout)
	}
	var p struct {
		Preview trajectory.ReconciliationPreview `json:"preview"`
		SHA256  string                           `json:"sha256"`
		Applied bool                             `json:"applied"`
	}
	if err := json.Unmarshal(out.Bytes(), &p); err != nil || p.Applied || p.SHA256 != p.Preview.SHA256() || p.Preview.Unchanged == nil || p.Preview.Unchanged.InvocationID != r.InvocationID {
		t.Fatal("invalid unchanged preview", err)
	}
	args = append(args, "--sha256", p.SHA256, "--yes")
	if code := runTrajectory(ctx, args, migrationOutputFailure{}, &errout); code != 1 {
		t.Fatal("output error hidden", code, &errout)
	}
	s, err := trajectory.Inspect(ctx, root, "idea-x")
	if err != nil || len(s.Resolutions) != 1 || s.Resolutions[0].Preview.Assessment.Outcome != trajectory.Inconclusive {
		t.Fatal("output failure lost observation", err)
	}
	continueTrajectory(t, root, "ack-first-unchanged", false, true)
	measuredUnchanged(t, root, "exit 0")
	before, _ := json.Marshal(s.Attempts)
	executions := parentRecoveryBytes(t, filepath.Join(root, ".parley-runtime", "unchanged-executions"))
	s, err = trajectory.Inspect(ctx, root, "idea-x")
	if err != nil {
		t.Fatal(err)
	}
	history, _ := json.Marshal(s)
	for i := 0; i < 2; i++ {
		out.Reset()
		if code := runTrajectory(ctx, args, &out, &errout); code != 0 {
			t.Fatal("historical CLI replay", code, &errout)
		}
		var applied struct {
			Preview trajectory.ReconciliationPreview `json:"preview"`
			SHA256  string                           `json:"sha256"`
			Applied bool                             `json:"applied"`
		}
		if err = json.Unmarshal(out.Bytes(), &applied); err != nil || !applied.Applied || applied.SHA256 != p.SHA256 || applied.Preview.SHA256() != p.SHA256 {
			t.Fatal("historical CLI replay changed preview", err)
		}
	}
	after, err := trajectory.Inspect(ctx, root, "idea-x")
	if err != nil {
		t.Fatal(err)
	}
	retained, _ := json.Marshal(after)
	first, _ := json.Marshal(after.Attempts[:1])
	if !bytes.Equal(history, retained) || !bytes.Equal(before, first) || !bytes.Equal(executions, parentRecoveryBytes(t, filepath.Join(root, ".parley-runtime", "unchanged-executions"))) {
		t.Fatal("replay rewrote evidence or executed work")
	}
	if err = trajectory.RequireResolved(ctx, root, "idea-x"); err == nil {
		t.Fatal("incomplete history accepted")
	}
}

func TestTrajectoryUnchangedConcurrentProductionCLI(t *testing.T) {
	binary := trajectoryHelperBinary(t)
	root := unchangedAppFixture(t)
	measuredUnchanged(t, root, "exit 0")
	args := []string{"trajectory", "reconcile-unchanged", "--dir", root, "--idea", "idea-x", "--sequence", "1"}
	raw, err := exec.Command(binary, args...).CombinedOutput()
	if err != nil {
		t.Fatalf("production CLI preview: %v %s", err, raw)
	}
	var preview struct {
		Preview trajectory.ReconciliationPreview `json:"preview"`
		SHA256  string                           `json:"sha256"`
		Applied bool                             `json:"applied"`
	}
	if err = json.Unmarshal(raw, &preview); err != nil || preview.Applied || preview.SHA256 != preview.Preview.SHA256() {
		t.Fatal("invalid production preview", err)
	}
	apply := append(args, "--sha256", preview.SHA256, "--yes")
	commands := []*exec.Cmd{exec.Command(binary, apply...), exec.Command(binary, apply...)}
	output := make([]bytes.Buffer, 2)
	for i, cmd := range commands {
		cmd.Stdout = &output[i]
		cmd.Stderr = &output[i]
		if err = cmd.Start(); err != nil {
			t.Fatal(err)
		}
	}
	for i, cmd := range commands {
		if err = cmd.Wait(); err != nil {
			t.Fatalf("concurrent production apply: %v %s", err, &output[i])
		}
	}
	if !bytes.Equal(output[0].Bytes(), output[1].Bytes()) {
		t.Fatal("concurrent output differs")
	}
	s, err := trajectory.Inspect(context.Background(), root, "idea-x")
	if err != nil || len(s.Attempts) != 1 || len(s.Resolutions) != 1 || len(s.Continuations) != 0 {
		t.Fatal("concurrent apply changed charge or continuation", err)
	}
	if string(parentRecoveryBytes(t, filepath.Join(root, ".parley-runtime", "unchanged-executions"))) != "executed\n" {
		t.Fatal("CLI replay executed another process")
	}
	if err = trajectory.RequireResolved(context.Background(), root, "idea-x"); err == nil {
		t.Fatal("production observation granted completion")
	}
}

func TestTrajectoryUnchangedBetweenRealRegressionsAndParentRecovery(t *testing.T) {
	ctx := context.Background()
	binary := trajectoryHelperBinary(t)
	secondary := trajectory.Criterion{Name: "secondary", Command: `if [ "$(cat source)" = broken2 ]; then printf 'PARLEY-EVIDENCE {"executed_cases":1,"failed_cases":1}\n'; exit 1; else printf 'PARLEY-EVIDENCE {"executed_cases":1,"failed_cases":0}\n'; fi`}
	root, trace, _, agent := trajectoryHelperFixtureCriteria(t, "actual-helper", []trajectory.Criterion{secondary})
	first, _ := reconciledHelper(t, root, binary, agent)
	commitTrajectorySource(t, root)
	continueTrajectory(t, root, "promote-first", false, false)
	measuredUnchanged(t, root, "exit 7")
	u := applyUnchanged(t, root, 2)
	if len(u.Assessment.Unresolved) != 2 {
		t.Fatal("unchanged observation reduced original material scope")
	}
	continueTrajectory(t, root, "ack-unchanged-between", false, true)
	if started, err := measuredTrajectoryPatch(t, root, "second-regression", "broken2"); err != nil || !started {
		t.Fatal(err)
	}
	_, third := reconciledHelper(t, root, binary, agent)
	if third.Assessment.Outcome != trajectory.Regression {
		t.Fatal("actual second regression not observed")
	}
	h, err := trajectory.InspectHistory(ctx, root, "idea-x")
	if err != nil || h.Decision.Consecutive != 2 || !h.ReviewPending || h.Decision.TriggerSequence != 3 || h.RequiredReviewSequence != 3 {
		t.Fatal("unchanged interrupted real regression streak", err, h)
	}
	commitTrajectorySource(t, root)
	continueTrajectory(t, root, "review-second-regression", true, false)
	measuredUnchanged(t, root, "exit 0")
	applyUnchanged(t, root, 4)
	h, err = trajectory.InspectHistory(ctx, root, "idea-x")
	if err != nil || h.Decision.Consecutive != 2 || h.ReviewPending || h.RequiredReviewSequence != 3 || h.Decision.TriggerSequence != 3 || len(h.InconclusivePending) != 1 || h.InconclusivePending[0] != 4 {
		t.Fatal("unchanged reopened prior acknowledged review", err, h)
	}
	// Recovery of the first lost parent must revalidate both typed unchanged
	// resolutions, without inventing a verifier run for either one.
	executions := parentRecoveryBytes(t, trace)
	if err = os.Remove(filepath.Join(filepath.Dir(first.RequestPath), "parent-result.json")); err != nil {
		t.Fatal(err)
	}
	recovery, err := trajectory.PreviewParentRecovery(ctx, root, "idea-x", first.RunID)
	if err != nil {
		t.Fatal("parent recovery could not dispatch unchanged evidence", err)
	}
	if _, err = trajectory.RecoverParent(ctx, root, "idea-x", first.RunID, recovery.SHA256()); err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(executions, parentRecoveryBytes(t, trace)) {
		t.Fatal("parent recovery re-executed criteria")
	}
	if _, err = trajectory.ReconcileUnchanged(ctx, root, "idea-x", 2, u.SHA256()); err != nil {
		t.Fatal("unchanged replay after later parent recovery", err)
	}
	continueTrajectory(t, root, "ack-final-unchanged", false, true)
	if err = trajectory.RequireResolved(ctx, root, "idea-x"); err == nil {
		t.Fatal("acknowledged unchanged observation closed trajectory")
	}
	if started, err := measuredTrajectoryPatch(t, root, "corrective-fifth", "original"); err != nil || !started {
		t.Fatal(err)
	}
	_, fifth := reconciledHelper(t, root, binary, agent)
	if fifth.Assessment.Outcome != trajectory.NoRegression {
		t.Fatal("final allowed corrective attempt lost")
	}
	if err = trajectory.RequireResolved(ctx, root, "idea-x"); err != nil {
		t.Fatal("final independently verified material result refused", err)
	}
	s, err := trajectory.Inspect(ctx, root, "idea-x")
	if err != nil || len(s.Attempts) != 5 || len(s.Resolutions) != 5 {
		t.Fatal("unchanged attempts did not stay spent", err)
	}
}

func TestTrajectoryUnchangedCLIRejectsInvalidAndChangedSource(t *testing.T) {
	root := unchangedAppFixture(t)
	base := []string{"reconcile-unchanged", "--dir", root, "--idea", "idea-x"}
	for _, flags := range [][]string{nil, {"--sequence", "0"}, {"--sequence", "-1"}, {"--sequence", "129"}, {"--sequence", "1", "--yes"}, {"--sequence", "1", "--sha256", "stale"}, {"--sequence", "1", "--run", "invented"}, {"--sequence", "1", "--acknowledge-inconclusive"}, {"--sequence", "1", "extra"}} {
		var out bytes.Buffer
		args := append(append([]string{}, base...), flags...)
		if code := runTrajectory(context.Background(), args, &out, &out); code != 2 {
			t.Fatal("invalid unchanged CLI accepted", flags, code, &out)
		}
	}
	measuredUnchanged(t, root, "printf 'changed\\n' > source")
	var errout bytes.Buffer
	if code := runTrajectory(context.Background(), append(base, "--sequence", "1"), io.Discard, &errout); code != 1 || !strings.Contains(errout.String(), "changed material source") {
		t.Fatal("changed source accepted as unchanged", code, &errout)
	}
}
