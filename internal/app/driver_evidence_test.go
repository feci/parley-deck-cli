package app

import (
	"context"
	"fmt"
	"io"
	"net"
	"os"
	"os/exec"
	"path/filepath"
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
		filepath.Join("parley-deck/ideas/idea-x", "00-prompt.md"): "---\nidea: idea-x\nstatus: implemented\n" + checksYAML + "---\n\n## Problem\n",
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

// Positive: a full run of the contract followed by an independent verifier's
// gate evaluation allows closure.
func TestEvidenceCloseGatePositive(t *testing.T) {
	root, ideaDir := gateScratchRepo(t, twoCriterionContract())
	o := gateOps(root, ideaDir)
	if ok, detail := o.runChecksContract(context.Background(), mustContract(t, ideaDir)); !ok {
		t.Fatalf("contract run should pass: %s", detail)
	}
	if _, err := os.Stat(evidence.ReportPath(ideaDir)); err != nil {
		t.Fatalf("typed report not persisted: %v", err)
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
// A client and a server rendezvous over a TCP loopback connection (the
// filesystem here cannot host FIFOs; /usr/bin/nc is a macOS built-in, not a
// new dependency). The server blocks until a client connects; the client
// fails immediately with no listener. Executed SERIALLY, both halves fail —
// no serial check runner can certify this criterion. Executed CONCURRENTLY
// (barrier), both pass. The evidence records make the distinction auditable:
// the serial attempt leaves fail-typed records, the barrier attempt leaves
// pass records with matching command digests.
func TestSerialVsBarrierConcurrencyFixture(t *testing.T) {
	root := t.TempDir()
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Skipf("loopback unavailable: %v", err)
	}
	port := ln.Addr().(*net.TCPAddr).Port
	_ = ln.Close()
	server := fmt.Sprintf("nc -l %d >/dev/null", port)
	client := fmt.Sprintf("printf hello | nc 127.0.0.1 %d", port)

	// Serial: client first (no listener → connection refused), then the
	// server bounded to 500ms (no client → deadline kill).
	recC := evidence.RunCriterion(context.Background(), root, "barrier-client", client, "kimi-1")
	serialCtx, cancel := context.WithTimeout(context.Background(), 500*time.Millisecond)
	recS := evidence.RunCriterion(serialCtx, root, "barrier-server", server, "kimi-1")
	cancel()
	if recC.Status != evidence.StatusFail || recS.Status != evidence.StatusFail {
		t.Fatalf("serial execution of a barrier pair must fail both halves, got client=%q server=%q", recC.Status, recS.Status)
	}

	// Barrier: both concurrently, 5s safety deadline.
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	ch := make(chan evidence.CriterionRecord, 2)
	go func() { ch <- evidence.RunCriterion(ctx, root, "barrier-server", server, "kimi-1") }()
	go func() { ch <- evidence.RunCriterion(ctx, root, "barrier-client", client, "kimi-1") }()
	for range 2 {
		if r := <-ch; r.Status != evidence.StatusPass {
			t.Fatalf("barrier execution must pass both halves, got %q (%+v)", r.Name, r)
		}
	}
	fmt.Printf("fixture: serial=(client:%s,server:%s) barrier=(pass,pass) — serial execution cannot certify a barrier criterion\n",
		recC.Status, recS.Status)
}

func gateHas(g EvidenceGateResult, sub string) bool {
	for _, r := range g.Reasons {
		if strings.Contains(r, sub) {
			return true
		}
	}
	return false
}
