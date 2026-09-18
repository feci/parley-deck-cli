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
	"strconv"
	"strings"
	"testing"
	"time"

	"parley-deck-cli/internal/agents"
	"parley-deck-cli/internal/budget"
	"parley-deck-cli/internal/procctl"
	"parley-deck-cli/internal/protocol"
	"parley-deck-cli/internal/runner"
	"parley-deck-cli/internal/telemetry"
	"parley-deck-cli/internal/trajectory"
)

func trajectoryHelperBinary(t *testing.T) string {
	t.Helper()
	if runtime.GOOS == "windows" {
		t.Skip("POSIX execution host required")
	}
	binary := filepath.Join(t.TempDir(), "parley helper fixture")
	cmd := exec.Command("go", "build", "-o", binary, "./cmd/parley")
	_, source, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("cannot locate the fixture's original source repository")
	}
	cmd.Dir = filepath.Join(filepath.Dir(source), "..", "..")
	if out, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("build real trajectory CLI: %v %s", err, out)
	}
	return binary
}

func trajectoryHelperFixture(t *testing.T, mode string) (root, trace, journal string, agent agents.Discovery) {
	return trajectoryHelperFixtureCriteria(t, mode, nil)
}

func trajectoryHelperFixtureCriteria(t *testing.T, mode string, extra []trajectory.Criterion) (root, trace, journal string, agent agents.Discovery) {
	return trajectoryHelperFixtureScope(t, mode, extra, "[builder, reviewer]")
}

func trajectoryHelperFixtureScope(t *testing.T, mode string, extra []trajectory.Criterion, members string) (root, trace, journal string, agent agents.Discovery) {
	t.Helper()
	t.Setenv("PARLEY_HOME", t.TempDir())
	t.Setenv("PARLEY_HEADLESS_AGENT_CONFIG", "")
	t.Setenv("PARLEY_TRAJECTORY_MODE", mode)
	trace = filepath.Join(t.TempDir(), "actual-executions")
	t.Setenv("PARLEY_TRAJECTORY_TRACE", trace)
	command := `printf '%s:%s\n' "$PARLEY_AGENT_ID" "$(cat source)" >> "$PARLEY_TRAJECTORY_TRACE"; if [ "$(cat source)" = original ]; then printf 'PARLEY-EVIDENCE {"executed_cases":1,"failed_cases":0}\n'; else printf 'PARLEY-EVIDENCE {"executed_cases":1,"failed_cases":1}\n'; exit 1; fi`
	if mode == "criterion-cancel" {
		command = `trap '' TERM; printf '%s\n' "$$" > "$PARLEY_TRAJECTORY_CRITERION_SIGNAL"; exec sleep 40`
		t.Setenv("PARLEY_TRAJECTORY_CRITERION_SIGNAL", filepath.Join(filepath.Dir(trace), "criterion.pid"))
	}
	if mode == "inconclusive" {
		command = `printf '%s:%s\n' "$PARLEY_AGENT_ID" "$(cat source)" >> "$PARLEY_TRAJECTORY_TRACE"; if [ "$(cat source)" = original ] || [ "$(wc -l < "$PARLEY_TRAJECTORY_TRACE" | tr -d ' ')" = 3 ]; then printf 'PARLEY-EVIDENCE {"executed_cases":1,"failed_cases":0}\n'; else printf 'PARLEY-EVIDENCE {"executed_cases":1,"failed_cases":1}\n'; exit 1; fi`
	}
	checks := append([]trajectory.Criterion{{Name: "material", Command: command}}, extra...)
	frontmatter := "participants: " + members + "\ntrack: deliberation\nchecks:\n"
	for _, c := range checks {
		frontmatter += "  - name: " + c.Name + "\n    command: >\n      " + strings.ReplaceAll(c.Command, "\n", "\n      ") + "\n"
	}
	root, ideaDir := gateScratchRepo(t, frontmatter)
	if err := protocol.InitWorkspace(root); err != nil {
		t.Fatal(err)
	}
	declareAppTestSource(t, root)
	for name, body := range map[string]string{".gitignore": ".parley-runtime/\nparley-deck/runs/\n", "source": "original\n"} {
		if err := os.WriteFile(filepath.Join(root, name), []byte(body), 0600); err != nil {
			t.Fatal(err)
		}
	}
	script := `#!/bin/sh
if [ "$1" = --version ]; then printf 'fixture\n'; exit 0; fi
command=$(awk '/^Verifier command: / {sub(/^Verifier command: /, ""); value=$0} END {print value}')
printf '%s\n' "$PARLEY_PROC_MARKER" >> .parley-runtime/verifier-starts
case "$PARLEY_TRAJECTORY_MODE" in
 text-only) printf 'PASS\n'; exit 0 ;;
 wrong-identity) export PARLEY_AGENT_ID=builder ;;
 request-tamper) for f in .parley-runtime/trajectory-verification/*/request.json; do printf '\n' >> "$f"; done ;;
 slow) exec sleep 40 ;;
 parent-publication) for f in .parley-runtime/trajectory-verification/*; do mkdir "$f/parent-result.json"; done ;;
esac
sh -c "$command"
code=$?
case "$PARLEY_TRAJECTORY_MODE" in
 failed-after-helper) exit 7 ;;
 missing-receipt) rm "$PARLEY_TRAJECTORY_JOURNAL/receipt.json" ;;
 scope-change) sed 's/name: material/name: changed/' parley-deck/ideas/idea-x/00-prompt.md > .parley-runtime/changed-prompt; mv .parley-runtime/changed-prompt parley-deck/ideas/idea-x/00-prompt.md ;;
esac
printf 'helper exit %s\n' "$code"
exit "$code"
`
	path := filepath.Join(root, "fixture-verifier")
	if err := os.WriteFile(path, []byte(script), 0700); err != nil {
		t.Fatal(err)
	}
	writeAgentsLocalConfig(t, root, fakeAgentConfig{ID: "reviewer", Path: path})
	gateGit(t, root, "add", "-A")
	gateGit(t, root, "-c", "user.email=t@t", "-c", "user.name=t", "commit", "-qm", "Frozen helper fixture")
	if _, err := budget.EnsureCycleBinding(context.Background(), root, "idea-x", budget.Fixup, 5, 0, "", ideaDir); err != nil {
		t.Fatal(err)
	}
	p, expected, err := trajectory.NewPolicy(context.Background(), root, "idea-x", "builder", checks)
	if err != nil {
		t.Fatal(err)
	}
	if err = trajectory.Activate(context.Background(), root, expected, p); err != nil {
		t.Fatal(err)
	}
	builder := agents.Discovery{Spec: agents.Spec{ID: "builder", PromptMode: agents.PromptStdin, LaunchMode: agents.LaunchHeadless, HeadlessArgs: []string{"-c", "printf 'broken\\n' > source; exit 7"}}, Found: true, Path: "/bin/sh"}
	if rec, err := runner.RunMeasured(context.Background(), runner.ExecOptions{Root: root, Agent: builder, Prompt: "synthetic patch fixture", Timeout: 30 * time.Second, Info: runner.LaunchInfo{RunID: "patch-fixture", Idea: "idea-x", Phase: "fixup"}}); err == nil || rec.StartedAt == nil {
		t.Fatalf("actual charged patch did not execute: %+v %v", rec, err)
	}
	binding, err := budget.LoadCycleBinding(context.Background(), root, "idea-x", budget.Fixup)
	if err != nil {
		t.Fatal(err)
	}
	state, err := trajectory.Inspect(context.Background(), root, "idea-x")
	if err != nil || state == nil || len(state.Attempts) != 1 {
		t.Fatal("actual patch was not captured", err)
	}
	journal = filepath.Join(filepath.Dir(binding.Store.Dir), "trajectory-verifications", state.Attempts[0].Charge.EntryKey)
	t.Setenv("PARLEY_TRAJECTORY_JOURNAL", journal)
	agent = agents.Discovery{Spec: agents.Spec{ID: "reviewer", PromptMode: agents.PromptStdin, LaunchMode: agents.LaunchHeadless, HeadlessArgs: []string{"--fixture"}}, Found: true, Path: path}
	return root, trace, journal, agent
}

func TestTrajectoryIndependentHelperProductionPath(t *testing.T) {
	binary := trajectoryHelperBinary(t)
	for _, mode := range []string{"actual-helper", "text-only", "wrong-identity", "request-tamper", "failed-after-helper", "missing-receipt", "scope-change", "parent-publication", "failed-start"} {
		t.Run(mode, func(t *testing.T) {
			root, trace, journal, agent := trajectoryHelperFixture(t, mode)
			if mode == "failed-start" {
				if err := os.Remove(agent.Path); err != nil {
					t.Fatal(err)
				}
			}
			var progress bytes.Buffer
			result, err := verifyTrajectoryWithAgent(context.Background(), root, "idea-x", agent, 90*time.Second, binary, &progress)
			want := mode == "actual-helper"
			if (err == nil) != want || (result.Assessment != nil) != want || !result.TrajectoryPending {
				t.Fatalf("mode %s: %+v %v\n%s", mode, result, err, progress.String())
			}
			if result.InvocationID == "" || result.TerminalSHA256 == "" {
				t.Fatalf("attempt lost actual terminal: %+v %v", result, err)
			}
			if want {
				if result.Assessment.Outcome != trajectory.Regression || result.ReceiptSHA256 == "" {
					t.Fatal("actual regression helper outcome not bound")
				}
				data, err := os.ReadFile(trace)
				if err != nil || string(data) != "reviewer:original\nreviewer:broken\nreviewer:broken\nreviewer:original\n" {
					t.Fatalf("independent helper did not execute original AB/BA criteria: %s %v", data, err)
				}
				var receipt trajectory.VerificationReceipt
				raw, err := readVerificationJSON(filepath.Join(journal, "receipt.json"), &receipt)
				if err != nil || sha256Hex(string(raw)) != result.ReceiptSHA256 || receipt.HelperPID == os.Getpid() || receipt.InvocationID != result.InvocationID {
					t.Fatalf("receipt not produced by observed child: %+v %v", receipt, err)
				}
			}
			// Neither a successful observation nor a failed launch closes the
			// trajectory or permits replay of the original charged patch.
			if err := trajectory.RequireResolved(context.Background(), root, "idea-x"); err == nil {
				t.Fatal("independent helper auto-resolved trajectory")
			}
			before, _ := os.ReadFile(filepath.Join(root, ".parley-runtime/verifier-starts"))
			if retry, err := verifyTrajectoryWithAgent(context.Background(), root, "idea-x", agent, 90*time.Second, binary, &progress); err == nil || retry.InvocationID != "" {
				t.Fatal("second request silently launched a new verifier", err)
			}
			after, _ := os.ReadFile(filepath.Join(root, ".parley-runtime/verifier-starts"))
			if !bytes.Equal(before, after) {
				t.Fatal("retry started another verifier")
			}
		})
	}
}

func TestTrajectoryIndependentHelperCLIAndIdentityRefusal(t *testing.T) {
	binary := trajectoryHelperBinary(t)
	root, trace, _, _ := trajectoryHelperFixture(t, "actual-helper")
	cmd := exec.Command(binary, "trajectory", "verify", "--dir", root, "--idea", "idea-x", "--verifier", "reviewer", "--timeout", "90s", "--yes")
	var out, errout bytes.Buffer
	cmd.Stdout, cmd.Stderr = &out, &errout
	if err := cmd.Run(); err != nil {
		t.Fatalf("actual parent/helper CLI failed: %v\n%s\n%s", err, out.String(), errout.String())
	}
	var result trajectoryVerificationResult
	if err := json.Unmarshal(out.Bytes(), &result); err != nil || result.Assessment == nil || result.Assessment.Outcome != trajectory.Regression {
		t.Fatalf("CLI did not return actual paired observation: %+v %v", result, err)
	}
	before, err := os.ReadFile(trace)
	if err != nil {
		t.Fatal(err)
	}
	for _, args := range [][]string{
		{"trajectory", "verify-helper", "--request", result.RequestPath, "--request-sha256", result.RequestSHA256},
		{"trajectory", "verify-helper", "--request", result.RequestPath, "--request-sha256", strings.Repeat("0", 64)},
		{"trajectory", "verify", "--dir", root, "--idea", "idea-x", "--verifier", "builder", "--yes"},
		{"trajectory", "verify", "--dir", root, "--idea", "idea-x", "--verifier", "reviewer"},
	} {
		cmd := exec.Command(binary, args...)
		cmd.Env = append(os.Environ(), "PARLEY_RUN_ID=", "PARLEY_AGENT_ID=", "PARLEY_PROC_MARKER=")
		if output, err := cmd.CombinedOutput(); err == nil {
			t.Fatalf("unauthorized/unbound helper or self-verifier accepted: %v %s", args, output)
		}
	}
	after, err := os.ReadFile(trace)
	if err != nil || !bytes.Equal(before, after) {
		t.Fatal("refused CLI entrypoint repeated executions", err)
	}
}

func TestTrajectoryIndependentVerifierCancellationRetainsLaunch(t *testing.T) {
	binary := trajectoryHelperBinary(t)
	root, _, journal, agent := trajectoryHelperFixture(t, "slow")
	ctx, cancel := context.WithTimeout(context.Background(), 90*time.Second)
	defer cancel()
	type outcome struct {
		r   trajectoryVerificationResult
		err error
	}
	done := make(chan outcome, 1)
	go func() {
		result, err := verifyTrajectoryWithAgent(ctx, root, "idea-x", agent, 90*time.Second, binary, io.Discard)
		done <- outcome{result, err}
	}()
	ticker := time.NewTicker(10 * time.Millisecond)
	defer ticker.Stop()
started:
	for {
		select {
		case early := <-done:
			t.Fatalf("verifier ended before actual start: %+v %v", early.r, early.err)
		case <-ctx.Done():
			t.Fatal("no actual verifier start before fixture bound")
		case <-ticker.C:
			if _, err := os.Stat(filepath.Join(root, ".parley-runtime/verifier-starts")); err == nil {
				break started
			}
		}
	}
	cancel()
	actual := <-done
	if actual.err == nil || actual.r.Assessment != nil || actual.r.InvocationID == "" || actual.r.TerminalSHA256 == "" {
		t.Fatalf("cancelled verifier result: %+v %v", actual.r, actual.err)
	}
	var terminal telemetry.Record
	if _, err := readVerificationJSON(filepath.Join(root, ".parley-runtime/invocations", actual.r.InvocationID, "terminal.json"), &terminal); err != nil || terminal.Outcome.FailureClass == nil || *terminal.Outcome.FailureClass != "cancelled" {
		t.Fatalf("actual cancellation terminal missing: %+v %v", terminal, err)
	}
	if _, err := os.Stat(filepath.Join(journal, "launch.json")); err != nil {
		t.Fatal("cancelled invocation reservation not retained", err)
	}
}

func TestTrajectoryVerificationRejectsUnknownParticipantBeforeLaunch(t *testing.T) {
	root, _, _, agent := trajectoryHelperFixture(t, "actual-helper")
	agent.ID = "unrostered"
	result, err := verifyTrajectoryWithAgent(context.Background(), root, "idea-x", agent, time.Second, "/unneeded", &bytes.Buffer{})
	if err == nil || result.InvocationID != "" {
		t.Fatalf("unknown participant was launched: %+v %v", result, err)
	}
	if _, err := os.Stat(filepath.Join(root, ".parley-runtime/verifier-starts")); !os.IsNotExist(err) {
		t.Fatalf("unrostered participant process started: %v", err)
	}
}

func TestTrajectoryIndependentCancellationReachesActiveCriterion(t *testing.T) {
	binary := trajectoryHelperBinary(t)
	root, trace, journal, agent := trajectoryHelperFixture(t, "criterion-cancel")
	ctx, cancel := context.WithTimeout(context.Background(), 90*time.Second)
	defer cancel()
	type outcome struct {
		result trajectoryVerificationResult
		err    error
	}
	done := make(chan outcome, 1)
	go func() {
		r, err := verifyTrajectoryWithAgent(ctx, root, "idea-x", agent, 90*time.Second, binary, io.Discard)
		done <- outcome{r, err}
	}()
	ticker := time.NewTicker(10 * time.Millisecond)
	defer ticker.Stop()
	var child procctl.Spawned
started:
	for {
		select {
		case early := <-done:
			t.Fatalf("helper ended before actual criterion: %+v %v", early.result, early.err)
		case <-ctx.Done():
			t.Fatal("criterion did not start before fixture deadline")
		case <-ticker.C:
			if raw, err := os.ReadFile(filepath.Join(filepath.Dir(trace), "criterion.pid")); err == nil {
				if pid, err := strconv.Atoi(strings.TrimSpace(string(raw))); err == nil && pid > 0 {
					child = procctl.CaptureByPID(pid, "criterion-cancellation-fixture")
					break started
				}
			}
		}
	}
	// The material shell may exec while a stable supervisor leads its group.
	// Cleanup remains attributed to the actual test-owned group if assertions fail.
	group := procctl.CaptureByPID(child.PGID, "criterion-cancellation-fixture-group")
	defer func() { _, _, _ = procctl.KillTreeAttributed(group) }()
	cancel()
	actual := <-done
	if actual.err == nil || actual.result.Assessment != nil {
		t.Fatal("interrupted helper produced an accepted assessment")
	}
	for _, name := range []string{"process-001.json", "stop.json"} {
		if info, err := os.Stat(filepath.Join(journal, name)); err != nil || !info.Mode().IsRegular() {
			t.Fatalf("cancelled verifier lost %s: %v", name, err)
		}
	}
	deadline := time.NewTimer(5 * time.Second)
	defer deadline.Stop()
	for procctl.Alive(child) {
		select {
		case <-deadline.C:
			t.Fatal("verifier cancellation left its TERM-resistant criterion process alive")
		case <-ticker.C:
		}
	}
}
