package app

import (
	"context"
	"fmt"
	"io"
	"net"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"testing"
	"time"

	"parley-deck-cli/internal/driver"
	"parley-deck-cli/internal/evidence"
)

// gateScratchRepo builds a real git repo with a Go file, an idea directory
// whose 00-prompt.md carries a named (list-form) checks contract, and an
// IMPLEMENTATION.md. A bare t.TempDir() inside this worktree is swallowed by
// the parent repo's ignore rules, so the digest/gate tests need a scratch repo.
func gateScratchRepo(t *testing.T, checksYAML string) (root, ideaDir string) {
	t.Helper()
	root = t.TempDir()
	gateGit(t, root, "init", "-q")
	ideaDir = filepath.Join(root, "parley-deck", "ideas", "idea-x")
	if err := os.MkdirAll(ideaDir, 0o755); err != nil {
		t.Fatal(err)
	}
	files := map[string]string{
		"go.mod":   "module scratch\n\ngo 1.26\n",
		"src/a.go": "package src\n",
		filepath.Join("parley-deck/ideas/idea-x", "00-prompt.md"):      "---\nidea: idea-x\nstatus: implemented\n" + checksYAML + "---\n\n## Problem\n",
		filepath.Join("parley-deck/ideas/idea-x", "IMPLEMENTATION.md"): "---\nidea: idea-x\nstatus: implemented\n---\n\n## Summary of work\n\ndone\n\n## Validation evidence\n\n(pending)\n",
	}
	for rel, content := range files {
		abs := filepath.Join(root, rel)
		if err := os.MkdirAll(filepath.Dir(abs), 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(abs, []byte(content), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	gateGit(t, root, "add", "-A")
	gateGit(t, root, "-c", "user.email=t@t", "-c", "user.name=t", "commit", "-qm", "init")
	return root, ideaDir
}

func gateGit(t *testing.T, dir string, args ...string) {
	t.Helper()
	cmd := exec.Command("git", append([]string{"-C", dir}, args...)...)
	if out, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("git %v: %v\n%s", args, err, out)
	}
}

// passJSONCmd prints a minimal passing test2json stream for one test case.
const passJSONCmd = `printf '%s\n' '{"Action":"run","Test":"TestA"}' '{"Action":"pass","Test":"TestA"}' '{"Action":"pass","Package":"x"}'`

func twoCriterionContract() string {
	return "checks:\n  - name: unit\n    command: >\n      " + passJSONCmd + "\n  - name: lint\n    command: >\n      " + passJSONCmd + "\n"
}

func gateOps(root, ideaDir string) driverImplOps {
	return driverImplOps{root: root, ideaSlug: "idea-x", ideaDir: ideaDir, implementer: "kimi-1", out: io.Discard}
}

// Positive: a full run of the contract followed by a REAL independent
// verifier re-execution — each criterion re-run by the verifier against the
// tested tree, bound before/after, attested and persisted — allows closure.
func TestEvidenceCloseGatePositive(t *testing.T) {
	root, ideaDir := gateScratchRepo(t, twoCriterionContract())
	o := gateOps(root, ideaDir)
	criteria := mustContract(t, ideaDir)
	if ok, detail := o.runChecksContract(context.Background(), criteria); !ok {
		t.Fatalf("contract run should pass: %s", detail)
	}
	report, err := evidence.Load(ideaDir)
	if err != nil {
		t.Fatalf("typed report not persisted: %v", err)
	}
	// The independent verifier re-executes every criterion itself and binds the
	// attestation to the tested tree digests taken immediately before/after.
	excl, err := definedEvidenceArtifacts(root, ideaDir)
	if err != nil {
		t.Fatal(err)
	}
	for _, c := range criteria {
		before, err := evidence.TreeDigest(root, excl...)
		if err != nil {
			t.Fatal(err)
		}
		rerun := evidence.RunCriterion(context.Background(), root, c.Name, c.Command, "codex-1")
		after, err := evidence.TreeDigest(root, excl...)
		if err != nil {
			t.Fatal(err)
		}
		if before != after || before != report.TreeSHA256 {
			t.Fatalf("verifier tree binding broken: before=%s after=%s report=%s", before[:12], after[:12], report.TreeSHA256[:12])
		}
		err = evidence.AttestExecution(report, c.Name, "codex-1", evidence.VerifierExecution{
			Command:          rerun.Command,
			TreeBeforeSHA256: before,
			TreeAfterSHA256:  after,
		})
		if err != nil {
			t.Fatalf("real independent attestation must succeed: %v", err)
		}
	}
	if err := evidence.Save(ideaDir, report); err != nil {
		t.Fatalf("persisting the attested report: %v", err)
	}
	gate := o.EvidenceCloseGate("codex-1")
	if !gate.Allowed {
		t.Fatalf("fresh independent evidence must allow closure: %v", gate.Reasons)
	}
}

func mustContract(t *testing.T, ideaDir string) []driver.CheckCriterion {
	t.Helper()
	criteria, isList, err := driver.ReadChecksContract(ideaDir)
	if err != nil || !isList {
		t.Fatalf("contract unreadable: %v", err)
	}
	return criteria
}

// Adversarial: no typed report at all → denied.
func TestEvidenceCloseGateNoReport(t *testing.T) {
	root, ideaDir := gateScratchRepo(t, twoCriterionContract())
	if gate := gateOps(root, ideaDir).EvidenceCloseGate("codex-1"); gate.Allowed {
		t.Fatal("missing report must not close")
	}
}

// Adversarial: the closing verifier is the executor — a self verdict.
func TestEvidenceCloseGateSelfVerdict(t *testing.T) {
	root, ideaDir := gateScratchRepo(t, twoCriterionContract())
	o := gateOps(root, ideaDir)
	if ok, _ := o.runChecksContract(context.Background(), mustContract(t, ideaDir)); !ok {
		t.Fatal("setup contract run failed")
	}
	gate := o.EvidenceCloseGate("kimi-1") // same identity as the executor
	if gate.Allowed || !gateHas(gate, "self verdict") {
		t.Fatalf("self verdict must be rejected: %+v", gate)
	}
}

// Adversarial: code changed after the checks ran — the tested tree is stale.
func TestEvidenceCloseGateStaleTree(t *testing.T) {
	root, ideaDir := gateScratchRepo(t, twoCriterionContract())
	o := gateOps(root, ideaDir)
	if ok, _ := o.runChecksContract(context.Background(), mustContract(t, ideaDir)); !ok {
		t.Fatal("setup contract run failed")
	}
	if err := os.WriteFile(filepath.Join(root, "src", "a.go"), []byte("package src // changed after evidence\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	gate := o.EvidenceCloseGate("codex-1")
	if gate.Allowed || !gateHas(gate, "stale") {
		t.Fatalf("stale tree must be rejected: %+v", gate)
	}
}

// Adversarial: an all-skip structured run is typed skipped and cannot close,
// even though the command exits 0.
func TestEvidenceCloseGateAllSkip(t *testing.T) {
	skipCmd := `printf '%s\n' '{"Action":"run","Test":"TestA"}' '{"Action":"skip","Test":"TestA"}' '{"Action":"pass","Package":"x"}'`
	yaml := "checks:\n  - name: unit\n    command: >\n      " + skipCmd + "\n  - name: lint\n    command: >\n      " + passJSONCmd + "\n"
	root, ideaDir := gateScratchRepo(t, yaml)
	o := gateOps(root, ideaDir)
	if ok, _ := o.runChecksContract(context.Background(), mustContract(t, ideaDir)); ok {
		t.Fatal("all-skip criterion must fail the cycle gate")
	}
	gate := o.EvidenceCloseGate("codex-1")
	if gate.Allowed || !gateHas(gate, "skipped") {
		t.Fatalf("skipped evidence must not close: %+v", gate)
	}
}

// Adversarial: structured proof of zero executed cases (exit 0) cannot close.
func TestEvidenceCloseGateZeroExecution(t *testing.T) {
	zeroCmd := `printf '%s\n' '{"Action":"start","Package":"x"}' '{"Action":"pass","Package":"x"}'`
	yaml := "checks:\n  - name: unit\n    command: >\n      " + zeroCmd + "\n  - name: lint\n    command: >\n      " + passJSONCmd + "\n"
	root, ideaDir := gateScratchRepo(t, yaml)
	o := gateOps(root, ideaDir)
	if ok, _ := o.runChecksContract(context.Background(), mustContract(t, ideaDir)); ok {
		t.Fatal("zero-execution criterion must fail the cycle gate")
	}
	gate := o.EvidenceCloseGate("codex-1")
	if gate.Allowed || !gateHas(gate, "no execution recorded") {
		t.Fatalf("zero-execution evidence must not close: %+v", gate)
	}
}

// Adversarial: partial scope — a record was dropped after the run.
func TestEvidenceCloseGatePartialScope(t *testing.T) {
	root, ideaDir := gateScratchRepo(t, twoCriterionContract())
	o := gateOps(root, ideaDir)
	if ok, _ := o.runChecksContract(context.Background(), mustContract(t, ideaDir)); !ok {
		t.Fatal("setup contract run failed")
	}
	report, err := evidence.Load(ideaDir)
	if err != nil {
		t.Fatal(err)
	}
	report.Records = report.Records[:1] // drop the lint record
	if err := evidence.Save(ideaDir, report); err != nil {
		t.Fatal(err)
	}
	gate := o.EvidenceCloseGate("codex-1")
	if gate.Allowed || !gateHas(gate, "partial scope") {
		t.Fatalf("partial scope must not close: %+v", gate)
	}
}

// Adversarial: opaque shell output passes the per-cycle gate (exit 0) but is
// NOT semantically certified — the close gate requires structured formats.
func TestEvidenceCloseGateShellNotCertified(t *testing.T) {
	yaml := "checks:\n  - name: unit\n    command: >\n      echo 'PASS all green'\n  - name: lint\n    command: >\n      " + passJSONCmd + "\n"
	root, ideaDir := gateScratchRepo(t, yaml)
	o := gateOps(root, ideaDir)
	if ok, _ := o.runChecksContract(context.Background(), mustContract(t, ideaDir)); !ok {
		t.Fatal("shell exit-0 should pass the cycle gate")
	}
	gate := o.EvidenceCloseGate("codex-1")
	if gate.Allowed || !gateHas(gate, "not semantically certified") {
		t.Fatalf("opaque shell output must not close: %+v", gate)
	}
}

// Adversarial: a corrupt typed report denies closure (fail closed).
func TestEvidenceCloseGateCorruptReport(t *testing.T) {
	root, ideaDir := gateScratchRepo(t, twoCriterionContract())
	o := gateOps(root, ideaDir)
	if ok, _ := o.runChecksContract(context.Background(), mustContract(t, ideaDir)); !ok {
		t.Fatal("setup contract run failed")
	}
	if err := os.WriteFile(evidence.ReportPath(ideaDir), []byte("{corrupt"), 0o644); err != nil {
		t.Fatal(err)
	}
	if gate := o.EvidenceCloseGate("codex-1"); gate.Allowed {
		t.Fatal("corrupt report must not close")
	}
}

// Concurrency fixture — serial versus barrier (rerunnable by a non-owner):
//
//	go test ./internal/app/ -run TestSerialVsBarrierConcurrencyFixture -v
//
// Portable: no netcat, no FIFOs (this host's /usr/bin/nc proved unreliable and
// the shared volume cannot host FIFOs). Both barrier halves are the test
// binary re-executing itself as a helper (TestHelperProcess, gated by
// PARLEY_BARRIER_HELPER). The server listens on a loopback port, publishes it
// atomically, and blocks in Accept; the client polls for the port file and
// fails on its own when no listener appears. Executed SERIALLY, both halves
// fail — no serial check runner can certify this criterion. Executed
// CONCURRENTLY (barrier), both pass with structured executed-case proof. The
// evidence records make the distinction auditable: the serial attempt leaves
// fail-typed records, the barrier attempt leaves pass records with 1 executed
// case each.
func TestSerialVsBarrierConcurrencyFixture(t *testing.T) {
	bin, err := os.Executable()
	if err != nil {
		t.Fatal(err)
	}
	helper := func(mode, dir string) string {
		return fmt.Sprintf("PARLEY_BARRIER_HELPER=1 %q -test.run '^TestHelperProcess$' -- %s %q", bin, mode, dir)
	}

	// Serial: client first (no port file → explicit failure after its poll
	// window), then the server bounded to 500ms (no client → deadline kill).
	serialDir := t.TempDir()
	recC := evidence.RunCriterion(context.Background(), serialDir, "barrier-client", helper("dial", serialDir), "kimi-1")
	serialCtx, cancel := context.WithTimeout(context.Background(), 500*time.Millisecond)
	recS := evidence.RunCriterion(serialCtx, serialDir, "barrier-server", helper("serve", serialDir), "kimi-1")
	cancel()
	if recC.Status != evidence.StatusFail || recS.Status != evidence.StatusFail {
		t.Fatalf("serial execution of a barrier pair must fail both halves, got client=%q server=%q", recC.Status, recS.Status)
	}

	// Barrier: both concurrently, 10s safety deadline.
	barrierDir := t.TempDir()
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	ch := make(chan evidence.CriterionRecord, 2)
	go func() {
		ch <- evidence.RunCriterion(ctx, barrierDir, "barrier-server", helper("serve", barrierDir), "kimi-1")
	}()
	go func() {
		ch <- evidence.RunCriterion(ctx, barrierDir, "barrier-client", helper("dial", barrierDir), "kimi-1")
	}()
	for range 2 {
		r := <-ch
		if r.Status != evidence.StatusPass || r.Command.Format != evidence.FormatGoTestJSON || r.Command.ExecutedCases != 1 {
			t.Fatalf("barrier execution must pass both halves with structured proof, got %q (%+v)", r.Name, r.Command)
		}
	}
	fmt.Printf("fixture: serial=(client:%s,server:%s) barrier=(pass,pass) — serial execution cannot certify a barrier criterion\n",
		recC.Status, recS.Status)
}

// TestHelperProcess is not a test: with PARLEY_BARRIER_HELPER=1 it runs the
// barrier helper (`serve`/`dial <dir>`) and exits, so the concurrency fixture
// needs nothing but this test binary. Both modes print a test2json pass event
// on success so the records carry structured executed-case proof.
func TestHelperProcess(t *testing.T) {
	if os.Getenv("PARLEY_BARRIER_HELPER") != "1" {
		return
	}
	args := os.Args
	for i, a := range args {
		if a == "--" {
			args = args[i+1:]
			break
		}
	}
	if len(args) != 2 {
		os.Exit(2)
	}
	mode, dir := args[0], args[1]
	portFile := filepath.Join(dir, "barrier-port")
	pass := func(test string) {
		fmt.Printf("{\"Action\":\"run\",\"Test\":%q}\n{\"Action\":\"pass\",\"Test\":%q}\n", test, test)
	}
	switch mode {
	case "serve":
		ln, err := net.Listen("tcp", "127.0.0.1:0")
		if err != nil {
			os.Exit(1)
		}
		// Publish the port atomically so the client never reads a partial write.
		tmp := portFile + ".tmp"
		if err := os.WriteFile(tmp, []byte(strconv.Itoa(ln.Addr().(*net.TCPAddr).Port)), 0o644); err != nil {
			os.Exit(1)
		}
		if err := os.Rename(tmp, portFile); err != nil {
			os.Exit(1)
		}
		conn, err := ln.Accept() // blocks until the client arrives or the deadline kill
		if err != nil {
			os.Exit(1)
		}
		_ = conn.Close()
		pass("TestBarrierServer")
		os.Exit(0)
	case "dial":
		deadline := time.Now().Add(2 * time.Second)
		var port []byte
		for {
			var err error
			port, err = os.ReadFile(portFile)
			if err == nil {
				break
			}
			if time.Now().After(deadline) {
				os.Exit(1) // no listener appeared — the serial half fails here
			}
			time.Sleep(25 * time.Millisecond)
		}
		conn, err := net.Dial("tcp", "127.0.0.1:"+strings.TrimSpace(string(port)))
		if err != nil {
			os.Exit(1)
		}
		_ = conn.Close()
		pass("TestBarrierClient")
		os.Exit(0)
	}
	os.Exit(2)
}

func gateHas(g EvidenceGateResult, sub string) bool {
	for _, r := range g.Reasons {
		if strings.Contains(r, sub) {
			return true
		}
	}
	return false
}

// attestContract runs the full contract and has verifier independently
// re-execute and attest every criterion (the same real-rerun discipline as
// TestEvidenceCloseGatePositive), returning the loaded, unsaved report.
func attestContract(t *testing.T, root, ideaDir string, o driverImplOps, verifier string) *evidence.Report {
	t.Helper()
	criteria := mustContract(t, ideaDir)
	if ok, detail := o.runChecksContract(context.Background(), criteria); !ok {
		t.Fatalf("contract run should pass: %s", detail)
	}
	report, err := evidence.Load(ideaDir)
	if err != nil {
		t.Fatalf("typed report not persisted: %v", err)
	}
	excl, err := definedEvidenceArtifacts(root, ideaDir)
	if err != nil {
		t.Fatal(err)
	}
	for _, c := range criteria {
		before, err := evidence.TreeDigest(root, excl...)
		if err != nil {
			t.Fatal(err)
		}
		rerun := evidence.RunCriterion(context.Background(), root, c.Name, c.Command, verifier)
		after, err := evidence.TreeDigest(root, excl...)
		if err != nil {
			t.Fatal(err)
		}
		if before != after || before != report.TreeSHA256 {
			t.Fatalf("verifier tree binding broken: before=%s after=%s report=%s", before[:12], after[:12], report.TreeSHA256[:12])
		}
		if err := evidence.AttestExecution(report, c.Name, verifier, evidence.VerifierExecution{
			Command:          rerun.Command,
			TreeBeforeSHA256: before,
			TreeAfterSHA256:  after,
		}); err != nil {
			t.Fatalf("real independent attestation must succeed: %v", err)
		}
	}
	return report
}

// authorizeCompletion has the attesting verifier authorize the completion
// transition from the exact bound rest bytes and persists the report — the
// adapter call the independently invoked verifier makes BEFORE Save.
func authorizeCompletion(t *testing.T, root, ideaDir string, report *evidence.Report, verifier string) {
	t.Helper()
	restContent, implRel, err := implementationRestContent(root, ideaDir)
	if err != nil {
		t.Fatal(err)
	}
	if err := evidence.AuthorizeCompletionTransition(report, implRel, restContent, verifier); err != nil {
		t.Fatalf("verifier authorization must succeed: %v", err)
	}
	if err := evidence.Save(ideaDir, report); err != nil {
		t.Fatal(err)
	}
}

// flipStatusComplete applies the deterministic completion transformation to
// IMPLEMENTATION.md — the same bytes Complete must write.
func flipStatusComplete(t *testing.T, ideaDir string) {
	t.Helper()
	implPath := filepath.Join(ideaDir, "IMPLEMENTATION.md")
	raw, err := os.ReadFile(implPath)
	if err != nil {
		t.Fatal(err)
	}
	completed, from, err := evidence.TransitionStatusToComplete(raw)
	if err != nil {
		t.Fatalf("deterministic completion transformation: %v", err)
	}
	if from != "implemented" {
		t.Fatalf("unexpected source status %q", from)
	}
	if err := os.WriteFile(implPath, completed, 0o644); err != nil {
		t.Fatal(err)
	}
}

// The facilitator's counterexample, fixed: independently attested evidence,
// verifier-authorized status transition, deterministic flip to status:
// complete — the close gate still allows, because the transition is verified
// against the current content and recomputed back to the recorded state.
func TestEvidenceCloseGateCompletionTransition(t *testing.T) {
	root, ideaDir := gateScratchRepo(t, twoCriterionContract())
	o := gateOps(root, ideaDir)
	report := attestContract(t, root, ideaDir, o, "codex-1")
	authorizeCompletion(t, root, ideaDir, report, "codex-1")
	flipStatusComplete(t, ideaDir)
	raw, err := os.ReadFile(filepath.Join(ideaDir, "IMPLEMENTATION.md"))
	if err != nil || !strings.Contains(string(raw), "\nstatus: complete\n") {
		t.Fatalf("the deterministic flip must land status: complete: %v", err)
	}
	if gate := o.EvidenceCloseGate("codex-1"); !gate.Allowed {
		t.Fatalf("an authorized, exactly-applied completion transition must close: %v", gate.Reasons)
	}
}

// Counterexample regression: the SAME status flip WITHOUT a recorded
// verifier-authorized transition is still denied — status/frontmatter scope
// stays bound, never silently excluded.
func TestEvidenceCloseGateStatusFlipWithoutTransitionDenied(t *testing.T) {
	root, ideaDir := gateScratchRepo(t, twoCriterionContract())
	o := gateOps(root, ideaDir)
	report := attestContract(t, root, ideaDir, o, "codex-1")
	if err := evidence.Save(ideaDir, report); err != nil {
		t.Fatal(err)
	}
	flipStatusComplete(t, ideaDir)
	gate := o.EvidenceCloseGate("codex-1")
	if gate.Allowed || !gateHas(gate, "non-evidence content changed") {
		t.Fatalf("an unauthorized status flip must deny: %+v", gate)
	}
}

// A transition authorized by one verifier does not serve a different selected
// verifier.
func TestEvidenceCloseGateCompletionTransitionWrongVerifier(t *testing.T) {
	root, ideaDir := gateScratchRepo(t, twoCriterionContract())
	o := gateOps(root, ideaDir)
	report := attestContract(t, root, ideaDir, o, "codex-1")
	authorizeCompletion(t, root, ideaDir, report, "codex-1")
	flipStatusComplete(t, ideaDir)
	gate := o.EvidenceCloseGate("opencode-1")
	if gate.Allowed || !gateHas(gate, "authorized by codex-1, not the selected verifier opencode-1") {
		t.Fatalf("a transition must not authorize a different verifier: %+v", gate)
	}
}

// The authorized transition covers ONLY the status flip: any extra
// non-evidence edit alongside it is still denied.
func TestEvidenceCloseGateCompletionTransitionExtraEditDenied(t *testing.T) {
	root, ideaDir := gateScratchRepo(t, twoCriterionContract())
	o := gateOps(root, ideaDir)
	report := attestContract(t, root, ideaDir, o, "codex-1")
	authorizeCompletion(t, root, ideaDir, report, "codex-1")
	flipStatusComplete(t, ideaDir)
	implPath := filepath.Join(ideaDir, "IMPLEMENTATION.md")
	raw, err := os.ReadFile(implPath)
	if err != nil {
		t.Fatal(err)
	}
	// An edit OUTSIDE the generated evidence section (the summary section).
	if err := os.WriteFile(implPath, []byte(strings.Replace(string(raw), "done", "done differently", 1)), 0o644); err != nil {
		t.Fatal(err)
	}
	gate := o.EvidenceCloseGate("codex-1")
	if gate.Allowed || !gateHas(gate, "not the authorized post-completion state") {
		t.Fatalf("an extra edit must not hide inside the authorized transition: %+v", gate)
	}
}

// An authorized-but-never-applied transition (the completion write failed or
// was vetoed after authorization) leaves the evidence in its original, still
// valid state — closure remains possible and completion can be retried.
func TestEvidenceCloseGateCompletionTransitionUnapplied(t *testing.T) {
	root, ideaDir := gateScratchRepo(t, twoCriterionContract())
	o := gateOps(root, ideaDir)
	report := attestContract(t, root, ideaDir, o, "codex-1")
	authorizeCompletion(t, root, ideaDir, report, "codex-1")
	if gate := o.EvidenceCloseGate("codex-1"); !gate.Allowed {
		t.Fatalf("the un-applied authorized transition must leave the original state valid: %v", gate.Reasons)
	}
}
