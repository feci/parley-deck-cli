package app

import (
	"bytes"
	"context"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
	"time"

	"parley-deck-cli/internal/agents"
	"parley-deck-cli/internal/consensus"
	"parley-deck-cli/internal/driver"
	"parley-deck-cli/internal/evidence"
	"parley-deck-cli/internal/protocol"
	"parley-deck-cli/internal/runner"
	"parley-deck-cli/internal/store"
)

// Only the prior review consensus is a fixture. RunChecks, the independent
// agent process, its real CLI helper, all criterion subprocesses and Complete
// use production implementations through Driver.Advance.
type evidenceCloseFixtureOps struct{ driverImplOps }

func (o evidenceCloseFixtureOps) ReviewRoundComplete(int) (bool, error) { return true, nil }
func (o evidenceCloseFixtureOps) ReviewStatus() (driver.ReviewStatus, error) {
	return driver.ReviewStatus{Summary: consensus.Summary{Triage: consensus.TriageReady}, ReviewerCount: 2}, nil
}

func TestEvidenceVerifierProductionClosure(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("POSIX fixture agent; Windows runtime requires its own execution host")
	}
	t.Setenv("PARLEY_HOME", t.TempDir())
	t.Setenv("PARLEY_HEADLESS_AGENT_CONFIG", "")
	t.Setenv("PARLEY_AGENT_ID", "")
	t.Setenv("PARLEY_TEST_VERIFY_MODE", "")
	binary := filepath.Join(t.TempDir(), "parley fixture")
	build := exec.Command("go", "build", "-o", binary, "./cmd/parley")
	build.Dir = filepath.Join("..", "..")
	if output, err := build.CombinedOutput(); err != nil {
		t.Fatalf("build actual helper CLI: %v\n%s", err, output)
	}
	cases := []struct {
		name, beforeHelper, afterHelper, mode string
		self, wantPass                        bool
	}{
		{name: "real-independent-execution", wantPass: true},
		{name: "text-pass-without-helper", beforeHelper: "printf 'GOAL-CHECK: PASS\\n'; exit 0"},
		{name: "self", self: true},
		{name: "missing-process", beforeHelper: "exit 7"},
		{name: "failed-process-after-helper", afterHelper: "exit 7"},
		{name: "stale-code-after-helper", afterHelper: "printf '\\n// changed after verification\\n' >> src/a.go"},
		{name: "changed-code-before-helper", beforeHelper: "printf '\\n// changed before verification\\n' >> src/a.go"},
		{name: "changed-implementation-after-helper", afterHelper: "printf '\\n## Unverified extra scope\\nExtra work claimed\\n' >> parley-deck/ideas/idea-x/IMPLEMENTATION.md"},
		{name: "skipped-independent-test", mode: "skip"},
		{name: "masked-package-failure", mode: "package-fail"},
		{name: "unwritable-report-path", beforeHelper: "mv parley-deck/ideas/idea-x/EVIDENCE.json .parley-runtime/original-report.json; mkdir parley-deck/ideas/idea-x/EVIDENCE.json"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			// Pipeline deliberately masks go test's exit code. Structured package
			// failures still must veto closure even if one real test passed.
			root, ideaDir := gateScratchRepo(t, "checks:\n  - name: unit\n    command: go test -count=1 -json ./src | cat\n")
			if err := protocol.InitWorkspace(root); err != nil {
				t.Fatal(err)
			}
			declareAppTestSource(t, root)
			if err := os.WriteFile(filepath.Join(root, ".gitignore"), []byte(".parley-runtime/\nparley-deck/runs/\n"), 0o644); err != nil {
				t.Fatal(err)
			}
			if err := os.MkdirAll(filepath.Join(root, ".parley-runtime"), 0o700); err != nil {
				t.Fatal(err)
			}
			fixtureTest := `package src
import ("fmt"; "os"; "testing")
func TestActualExecution(t *testing.T) {
 f, err := os.OpenFile("../.parley-runtime/executions.txt", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0600)
 if err != nil { t.Fatal(err) }
 fmt.Fprintf(f, "%s:%d\n", os.Getenv("PARLEY_AGENT_ID"), os.Getpid()); f.Close()
 if os.Getenv("PARLEY_TEST_VERIFY_MODE") == "skip" { t.Skip("real skipped verifier fixture") }
}
func TestMain(m *testing.M) {
 code := m.Run()
 if os.Getenv("PARLEY_TEST_VERIFY_MODE") == "package-fail" { os.Exit(1) }
 os.Exit(code)
}
`
			if err := os.WriteFile(filepath.Join(root, "src/a_test.go"), []byte(fixtureTest), 0o644); err != nil {
				t.Fatal(err)
			}
			reviewDir := filepath.Join(ideaDir, "review", "round-01")
			if err := os.MkdirAll(reviewDir, 0o755); err != nil {
				t.Fatal(err)
			}
			if err := os.WriteFile(filepath.Join(ideaDir, "review", "consensus.md"), []byte("Prior review consensus is supplied by the fixture adapter.\n"), 0o644); err != nil {
				t.Fatal(err)
			}
			script := "#!/bin/sh\ncommand=$(awk '/^Verifier command: / {sub(/^Verifier command: /, \"\"); value=$0} END {print value}')\n" + tc.beforeHelper + "\n" +
				"export PARLEY_TEST_VERIFY_MODE='" + tc.mode + "'\nsh -c \"$command\"\nresult=$?\n" + tc.afterHelper + "\nprintf 'helper exit %s\\n' \"$result\"\nexit \"$result\"\n"
			path := filepath.Join(root, "fixture-verifier")
			if err := os.WriteFile(path, []byte(script), 0o700); err != nil {
				t.Fatal(err)
			}
			gateGit(t, root, "add", "-A")
			gateGit(t, root, "-c", "user.email=t@t", "-c", "user.name=t", "commit", "-qm", "fixture inputs")
			runDir := filepath.Join(root, protocol.DeckDir, "runs", "verification-fixture")
			var progress bytes.Buffer
			op := driverImplOps{root: root, ideaDir: ideaDir, ideaSlug: "idea-x", implementer: "implementer", drafter: "reviewer", reviewers: []string{"reviewer"}, out: &progress, verificationCLI: binary,
				base: runner.Options{Root: root, RunID: "verification-fixture", Store: store.New(runDir), Timeout: 20 * time.Second,
					Agents: []agents.Discovery{{Spec: agents.Spec{ID: "reviewer", PromptMode: agents.PromptStdin, LaunchMode: agents.LaunchHeadless, HeadlessArgs: []string{"--fixture"}}, Found: true, Path: path}}}}
			if tc.self {
				op.drafter = op.implementer
			}
			d := driver.New(driver.Config{Root: root, IdeaDir: ideaDir, IdeaSlug: "idea-x", RunDir: runDir, Participants: []string{"implementer", "reviewer"}, Auto: true, Events: store.New(runDir), Impl: evidenceCloseFixtureOps{op}, Out: io.Discard}, nil)
			action, _, err := d.Advance(context.Background())
			impl, readErr := os.ReadFile(filepath.Join(ideaDir, "IMPLEMENTATION.md"))
			if readErr != nil {
				t.Fatal(readErr)
			}
			if tc.wantPass {
				if err != nil || action != driver.ActionComplete || !strings.Contains(string(impl), "status: complete") {
					logs, _ := filepath.Glob(filepath.Join(root, ".parley-runtime", "evidence-verification", "*", "agent.*.log"))
					for _, p := range logs {
						data, _ := os.ReadFile(p)
						t.Logf("fixture log %s: %s", filepath.Base(p), data)
					}
					t.Fatalf("real verifier did not close: %s %v\n%s", action, err, progress.String())
				}
				executions, err := os.ReadFile(filepath.Join(root, ".parley-runtime/executions.txt"))
				if err != nil {
					t.Fatal(err)
				}
				lines := strings.Fields(string(executions))
				if len(lines) != 2 || !strings.HasPrefix(lines[0], ":") || !strings.HasPrefix(lines[1], "reviewer:") || strings.TrimPrefix(lines[0], ":") == strings.TrimPrefix(lines[1], "reviewer:") {
					t.Fatalf("not two distinct actual test processes: %q", executions)
				}
			} else if err == nil || action != driver.ActionEscalated || strings.Contains(string(impl), "status: complete") {
				t.Fatalf("invalid verification closed: %s %v\n%s", action, err, progress.String())
			}
			if !tc.self && tc.name != "text-pass-without-helper" && tc.name != "missing-process" {
				paths, _ := filepath.Glob(filepath.Join(root, ".parley-runtime", "evidence-verification", "*", "result.json"))
				if len(paths) != 1 {
					t.Fatalf("expected actual helper receipt, found %v; failure could be unrelated: %v\n%s", paths, err, progress.String())
				}
				var receipt evidenceVerificationReceipt
				if _, err := readVerificationJSON(paths[0], &receipt); err != nil {
					t.Fatal(err)
				}
				if receipt.Agent != "reviewer" || receipt.RunID != "verification-fixture" || receipt.HelperPID == os.Getpid() || receipt.ProcessMarker == "" {
					t.Fatalf("helper was not a distinct attributed process: %+v", receipt)
				}
				if tc.mode != "" {
					if len(receipt.Executions) != 1 || receipt.Error == "" {
						t.Fatalf("no actual failed rerun: %+v", receipt)
					}
					rec := receipt.Executions[0]
					if tc.mode == "skip" && (rec.Status != evidence.StatusSkipped || rec.Command.ExecutedCases != 0 || rec.Command.SkippedCases != 1) {
						t.Fatalf("not the skipped-test failure: %+v", rec)
					}
					if tc.mode == "package-fail" && (rec.Status != evidence.StatusFail || rec.Command.ExitCode != 0 || rec.Command.ExecutedCases != 1 || rec.Command.FailedPackages <= 0) {
						t.Fatalf("not the masked package failure: %+v", rec)
					}
				}
			}
		})
	}
}

func TestEvidenceHelperRefusesMissingRuntimeIdentity(t *testing.T) {
	t.Setenv("PARLEY_RUN_ID", "")
	t.Setenv("PARLEY_AGENT_ID", "")
	t.Setenv("PARLEY_PROC_MARKER", "")
	root, ideaDir := gateScratchRepo(t, twoCriterionContract())
	canonical, err := filepath.EvalSymlinks(root)
	if err != nil {
		t.Fatal(err)
	}
	dir := filepath.Join(canonical, ".parley-runtime", "evidence-verification", "attempt-x")
	if err := os.MkdirAll(dir, 0o700); err != nil {
		t.Fatal(err)
	}
	path := filepath.Join(dir, "request.json")
	req := evidenceVerificationRequest{Version: 1, Root: canonical, Idea: "idea-x", RunID: "run", Verifier: "reviewer"}
	if err := writeVerificationJSON(path, req); err != nil {
		t.Fatal(err)
	}
	raw, _ := os.ReadFile(path)
	var out, stderr bytes.Buffer
	if code := runEvidenceVerify(context.Background(), []string{"verify", "--request", path, "--request-sha256", sha256Hex(string(raw))}, &out, &stderr); code != 1 || !strings.Contains(stderr.String(), "inside the selected independent verifier") {
		t.Fatalf("unattributed helper accepted: %d %s", code, stderr.String())
	}
	if _, err := os.Stat(filepath.Join(ideaDir, "EVIDENCE.json")); !os.IsNotExist(err) {
		t.Fatal("unattributed helper wrote report")
	}
	if _, err := os.Stat(filepath.Join(dir, "result.json")); !os.IsNotExist(err) {
		t.Fatal("unattributed helper wrote receipt")
	}
}
