package evidence

import (
	"strings"
	"testing"
	"time"
)

// positiveRecord returns a fully valid passing record for criterion `name`.
func positiveRecord(name, executor string) CriterionRecord {
	return CriterionRecord{
		Name:   name,
		Status: StatusPass,
		Command: CommandEvidence{
			Command:        "go test -json ./x",
			CommandSHA256:  "abc",
			OutputSHA256:   "def",
			ExitCode:       0,
			DurationMillis: 3,
			Format:         FormatGoTestJSON,
			ExecutedCases:  2,
			FailedCases:    0,
		},
		Provenance: Provenance{Executor: executor, Verifier: "codex-1", VerifierTreeSHA256: "treehash"},
	}
}

func positiveReport() *Report {
	return &Report{
		Idea:           "idea-x",
		ReviewedCommit: "0123456789abcdef",
		TreeSHA256:     "treehash",
		GeneratedAt:    time.Now(),
		Records: []CriterionRecord{
			positiveRecord("unit", "kimi-1"),
			positiveRecord("integration", "kimi-1"),
		},
	}
}

func positiveOpts() ClosureOptions {
	return ClosureOptions{
		RequiredScope:     []string{"unit", "integration"},
		CurrentTreeSHA256: "treehash",
		Verifier:          "codex-1",
		RequireStructured: true,
	}
}

// Positive: a complete, fresh, independently verified report closes.
func TestEvaluatePositiveCloses(t *testing.T) {
	if reasons := Evaluate(positiveReport(), positiveOpts()); len(reasons) != 0 {
		t.Fatalf("valid evidence must close, got: %v", reasons)
	}
}

func TestEvaluateNilReportFailsClosed(t *testing.T) {
	if reasons := Evaluate(nil, positiveOpts()); len(reasons) == 0 {
		t.Fatal("nil report must not close")
	}
}

// Adversarial: an empty verifier identity must not close.
func TestEvaluateNoVerifier(t *testing.T) {
	opts := positiveOpts()
	opts.Verifier = ""
	if reasons := Evaluate(positiveReport(), opts); len(reasons) == 0 {
		t.Fatal("missing verifier must not close")
	}
}

// Adversarial: self verdict — the closing verifier is the executor.
func TestEvaluateSelfVerdictRejected(t *testing.T) {
	opts := positiveOpts()
	opts.Verifier = "kimi-1" // same as executor
	reasons := Evaluate(positiveReport(), opts)
	if !containsAny(reasons, "self verdict") {
		t.Fatalf("self verdict must be rejected, got: %v", reasons)
	}
}

// Adversarial: stale tree — code changed after the checks ran.
func TestEvaluateStaleTreeRejected(t *testing.T) {
	opts := positiveOpts()
	opts.CurrentTreeSHA256 = "different-tree"
	if reasons := Evaluate(positiveReport(), opts); !containsAny(reasons, "stale") {
		t.Fatalf("stale tree must be rejected, got: %v", reasons)
	}
}

// Adversarial: partial scope — one required criterion has no record at all.
func TestEvaluatePartialScopeRejected(t *testing.T) {
	r := positiveReport()
	r.Records = r.Records[:1]
	if reasons := Evaluate(r, positiveOpts()); !containsAny(reasons, "partial scope") {
		t.Fatalf("partial scope must be rejected, got: %v", reasons)
	}
}

// Adversarial: every non-pass typed status blocks closure.
func TestEvaluateNonPassStatusesRejected(t *testing.T) {
	for _, st := range []Status{StatusFail, StatusSkipped, StatusNotRun} {
		r := positiveReport()
		r.Records[0].Status = st
		if reasons := Evaluate(r, positiveOpts()); len(reasons) == 0 {
			t.Fatalf("status %q must not close", st)
		}
	}
}

// Adversarial: zero executed cases in a structured format is a no-execution
// report even when the criterion claims pass.
func TestEvaluateZeroExecutionRejected(t *testing.T) {
	r := positiveReport()
	r.Records[0].Command.ExecutedCases = 0
	if reasons := Evaluate(r, positiveOpts()); !containsAny(reasons, "zero executed cases") {
		t.Fatalf("zero executed cases must be rejected, got: %v", reasons)
	}
}

// Adversarial: failed cases recorded in a pass-claimed record still block.
func TestEvaluateFailedCasesRejected(t *testing.T) {
	r := positiveReport()
	r.Records[0].Command.FailedCases = 1
	if reasons := Evaluate(r, positiveOpts()); !containsAny(reasons, "failed cases") {
		t.Fatalf("failed cases must be rejected, got: %v", reasons)
	}
}

// Adversarial: opaque shell output is NOT semantically certified — under a
// structured-evidence contract a shell exit 0 cannot close.
func TestEvaluateShellNotCertifiedUnderStructuredContract(t *testing.T) {
	r := positiveReport()
	r.Records[0].Command.Format = FormatShell
	r.Records[0].Command.ExecutedCases = -1
	r.Records[0].Command.FailedCases = -1
	if reasons := Evaluate(r, positiveOpts()); !containsAny(reasons, "not semantically certified") {
		t.Fatalf("shell output must not close under RequireStructured, got: %v", reasons)
	}
	// But without the structured requirement, a pass-typed shell record with
	// hash binding is admissible (the contract decides, not the regex).
	opts := positiveOpts()
	opts.RequireStructured = false
	if reasons := Evaluate(r, opts); len(reasons) != 0 {
		t.Fatalf("shell record should close when structured not required, got: %v", reasons)
	}
}

// Adversarial: unknown/unrecognized format strings are not certified either.
func TestEvaluateUnknownFormatRejected(t *testing.T) {
	r := positiveReport()
	r.Records[0].Command.Format = "junit-xml-v9"
	r.Records[0].Command.ExecutedCases = -1
	r.Records[0].Command.FailedCases = -1
	if reasons := Evaluate(r, positiveOpts()); !containsAny(reasons, "unsupported evidence format") {
		t.Fatalf("unknown format must be rejected, got: %v", reasons)
	}
}

// Adversarial: a record with no executor provenance cannot close.
func TestEvaluateNoExecutorRejected(t *testing.T) {
	r := positiveReport()
	r.Records[1].Provenance.Executor = ""
	if reasons := Evaluate(r, positiveOpts()); !containsAny(reasons, "no executor provenance") {
		t.Fatalf("missing executor must be rejected, got: %v", reasons)
	}
}

// Adversarial: a record missing its command/output hash binding cannot close.
func TestEvaluateMissingHashBindingRejected(t *testing.T) {
	r := positiveReport()
	r.Records[0].Command.OutputSHA256 = ""
	if reasons := Evaluate(r, positiveOpts()); !containsAny(reasons, "hash binding") {
		t.Fatalf("missing hash binding must be rejected, got: %v", reasons)
	}
}

// Adversarial: an envelope-formatted record with zero executed cases is also
// no-execution (envelopes assert, they do not get a free pass).
func TestEvaluateEnvelopeZeroCasesRejected(t *testing.T) {
	r := positiveReport()
	r.Records[0].Command.Format = FormatEnvelope
	r.Records[0].Command.ExecutedCases = 0
	r.Records[0].Command.FailedCases = 0
	if reasons := Evaluate(r, positiveOpts()); !containsAny(reasons, "zero executed cases") {
		t.Fatalf("envelope zero cases must be rejected, got: %v", reasons)
	}
}

// Adversarial (probe 1): a passing record whose Executor is the implementer
// with EMPTY persisted Verifier must NOT close merely because the caller
// supplies a different ClosureOptions.Verifier name.
func TestEvaluateUnattestedVerifierRejected(t *testing.T) {
	r := positiveReport()
	for i := range r.Records {
		r.Records[i].Provenance.Verifier = ""
		r.Records[i].Provenance.VerifierTreeSHA256 = ""
	}
	reasons := Evaluate(r, positiveOpts())
	if !containsAny(reasons, "no persisted independent verifier attestation") {
		t.Fatalf("unattested record must be rejected, got: %v", reasons)
	}
}

// Adversarial: the closing identity must match the PERSISTED attestation — a
// caller naming a different verifier creates nothing.
func TestEvaluateVerifierMismatchRejected(t *testing.T) {
	opts := positiveOpts()
	opts.Verifier = "someone-else"
	if reasons := Evaluate(positiveReport(), opts); !containsAny(reasons, "does not match persisted attestation") {
		t.Fatalf("foreign verifier name must be rejected, got: %v", reasons)
	}
}

// Adversarial: an attestation bound to a different tree than the report's is
// not evidence for this report.
func TestEvaluateAttestationTreeMismatchRejected(t *testing.T) {
	r := positiveReport()
	r.Records[0].Provenance.VerifierTreeSHA256 = "other-tree"
	if reasons := Evaluate(r, positiveOpts()); !containsAny(reasons, "bound to a different tree") {
		t.Fatalf("attestation bound to another tree must be rejected, got: %v", reasons)
	}
}

// Adversarial: empty scope, absent current digest, duplicate records,
// contradictory exit/status, and invalid counts each fail closed.
func TestEvaluateFailClosedSurface(t *testing.T) {
	opts := positiveOpts()
	opts.RequiredScope = nil
	if reasons := Evaluate(positiveReport(), opts); !containsAny(reasons, "empty required scope") {
		t.Fatalf("empty scope must be rejected, got: %v", reasons)
	}
	opts = positiveOpts()
	opts.CurrentTreeSHA256 = ""
	if reasons := Evaluate(positiveReport(), opts); !containsAny(reasons, "no current tree digest") {
		t.Fatalf("absent current digest must be rejected, got: %v", reasons)
	}
	r := positiveReport()
	r.Records = append(r.Records, r.Records[0])
	if reasons := Evaluate(r, positiveOpts()); !containsAny(reasons, "duplicate evidence record") {
		t.Fatalf("duplicate records must be rejected, got: %v", reasons)
	}
	r = positiveReport()
	r.Records[0].Command.ExitCode = 1
	if reasons := Evaluate(r, positiveOpts()); !containsAny(reasons, "contradictory record") {
		t.Fatalf("pass with non-zero exit must be rejected, got: %v", reasons)
	}
	r = positiveReport()
	r.Records[0].Command.ExecutedCases = -7
	if reasons := Evaluate(r, positiveOpts()); !containsAny(reasons, "invalid (negative) case count") {
		t.Fatalf("invalid counts must be rejected, got: %v", reasons)
	}
	r = positiveReport()
	r.Records[0].Command.FailedCases = -1
	if reasons := Evaluate(r, positiveOpts()); !containsAny(reasons, "failed-case count absent") {
		t.Fatalf("absent structured failed count must be rejected, got: %v", reasons)
	}
}

// AttestExecution: the only path to a verifier name on a record runs through
// evidence of an independent execution.
func TestAttestExecution(t *testing.T) {
	r := positiveReport()
	for i := range r.Records {
		r.Records[i].Provenance.Verifier = ""
		r.Records[i].Provenance.VerifierTreeSHA256 = ""
	}
	rerun := CommandEvidence{ExitCode: 0, Format: FormatGoTestJSON, ExecutedCases: 2, FailedCases: 0}
	if err := AttestExecution(r, "unit", "codex-1", rerun); err != nil {
		t.Fatalf("valid attestation must succeed: %v", err)
	}
	if r.Records[0].Provenance.Verifier != "codex-1" || r.Records[0].Provenance.VerifierTreeSHA256 != r.TreeSHA256 {
		t.Fatalf("attestation must bind verifier and tree, got %+v", r.Records[0].Provenance)
	}
	if err := AttestExecution(r, "integration", "kimi-1", rerun); err == nil {
		t.Fatal("self attestation must be refused")
	}
	if err := AttestExecution(r, "integration", "codex-1", CommandEvidence{ExitCode: 1, Format: FormatGoTestJSON, ExecutedCases: 2}); err == nil {
		t.Fatal("failed re-run must not attest")
	}
	if err := AttestExecution(r, "integration", "codex-1", CommandEvidence{ExitCode: 0, Format: FormatShell}); err == nil {
		t.Fatal("opaque shell re-run must not attest")
	}
	if err := AttestExecution(r, "integration", "codex-1", CommandEvidence{ExitCode: 0, Format: FormatGoTestJSON, ExecutedCases: 0, FailedCases: 0}); err == nil {
		t.Fatal("zero-execution re-run must not attest")
	}
	if err := AttestExecution(r, "nonexistent", "codex-1", rerun); err == nil {
		t.Fatal("attesting an unknown criterion must fail")
	}
	if err := AttestExecution(r, "integration", "", rerun); err == nil {
		t.Fatal("empty verifier identity must fail")
	}
	noTree := &Report{}
	if err := AttestExecution(noTree, "unit", "codex-1", rerun); err == nil {
		t.Fatal("attestation without a tested tree digest must fail")
	}
}

func containsAny(reasons []string, sub string) bool {
	for _, r := range reasons {
		if strings.Contains(r, sub) {
			return true
		}
	}
	return false
}
