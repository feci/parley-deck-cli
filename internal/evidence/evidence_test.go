package evidence

import (
	"strings"
	"testing"
	"time"
)

// positiveRecord returns a fully valid passing record for criterion `name`,
// including a retained independent attestation that reconciles with it.
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
			SkippedCases:   0,
		},
		Provenance: Provenance{
			Executor: executor,
			Verifier: "codex-1",
			VerifierRerun: &VerifierExecution{
				Command: CommandEvidence{
					CommandSHA256: "abc",
					OutputSHA256:  "rerun-out",
					ExitCode:      0,
					Format:        FormatGoTestJSON,
					ExecutedCases: 2,
					FailedCases:   0,
					SkippedCases:  0,
				},
				TreeBeforeSHA256: "treehash",
				TreeAfterSHA256:  "treehash",
			},
		},
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
		r.Records[i].Provenance.VerifierRerun = nil
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

// Adversarial: an attestation whose retained re-run executed against a
// different tree than the report's is not evidence for this report.
func TestEvaluateAttestationTreeMismatchRejected(t *testing.T) {
	r := positiveReport()
	r.Records[0].Provenance.VerifierRerun.TreeBeforeSHA256 = "other-tree"
	r.Records[0].Provenance.VerifierRerun.TreeAfterSHA256 = "other-tree"
	if reasons := Evaluate(r, positiveOpts()); !containsAny(reasons, "different tree than the report tested") {
		t.Fatalf("attestation bound to another tree must be rejected, got: %v", reasons)
	}
}

// Adversarial: a verifier name with NO retained independent execution is not a
// verdict — the rerun record must be present and must reconcile.
func TestEvaluateAttestationWithoutRetainedRerunRejected(t *testing.T) {
	r := positiveReport()
	r.Records[1].Provenance.VerifierRerun = nil
	if reasons := Evaluate(r, positiveOpts()); !containsAny(reasons, "retains no independent execution record") {
		t.Fatalf("name-only attestation must be rejected, got: %v", reasons)
	}
}

// Adversarial: a retained re-run of a DIFFERENT command (unrelated output) is
// rejected at close, even when the verifier name matches.
func TestEvaluateUnrelatedRerunRejected(t *testing.T) {
	r := positiveReport()
	r.Records[0].Provenance.VerifierRerun.Command.CommandSHA256 = "unrelated-command"
	r.Records[0].Provenance.VerifierRerun.Command.OutputSHA256 = "unrelated-output"
	if reasons := Evaluate(r, positiveOpts()); !containsAny(reasons, "not the same exact command") {
		t.Fatalf("unrelated re-run command must be rejected, got: %v", reasons)
	}
}

// Adversarial: a retained re-run whose tree changed mid-verification, whose
// counts diverge from the original record, or that failed, is rejected.
func TestEvaluateRerunReconciliationRejected(t *testing.T) {
	r := positiveReport()
	r.Records[0].Provenance.VerifierRerun.TreeAfterSHA256 = "changed-during-rerun"
	if reasons := Evaluate(r, positiveOpts()); !containsAny(reasons, "changed during the independent re-run") {
		t.Fatalf("unstable verification tree must be rejected, got: %v", reasons)
	}
	r = positiveReport()
	r.Records[0].Provenance.VerifierRerun.Command.ExecutedCases = 7
	if reasons := Evaluate(r, positiveOpts()); !containsAny(reasons, "executed-case count differs") {
		t.Fatalf("diverging executed counts must be rejected, got: %v", reasons)
	}
	r = positiveReport()
	r.Records[0].Provenance.VerifierRerun.Command.FailedCases = 1
	if reasons := Evaluate(r, positiveOpts()); !containsAny(reasons, "independent re-run recorded 1 failed cases") {
		t.Fatalf("failed re-run must be rejected, got: %v", reasons)
	}
	r = positiveReport()
	r.Records[0].Provenance.VerifierRerun.Command.OutputSHA256 = ""
	if reasons := Evaluate(r, positiveOpts()); !containsAny(reasons, "binds no output hash") {
		t.Fatalf("re-run without output binding must be rejected, got: %v", reasons)
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
// retained, bound evidence of an independent execution.
func TestAttestExecution(t *testing.T) {
	fresh := func() *Report {
		r := positiveReport()
		for i := range r.Records {
			r.Records[i].Provenance.Verifier = ""
			r.Records[i].Provenance.VerifierRerun = nil
		}
		return r
	}
	goodRerun := VerifierExecution{
		Command: CommandEvidence{
			CommandSHA256: "abc", // matches positiveRecord's command hash
			OutputSHA256:  "rerun-out",
			ExitCode:      0,
			Format:        FormatGoTestJSON,
			ExecutedCases: 2,
			FailedCases:   0,
			SkippedCases:  0,
		},
		TreeBeforeSHA256: "treehash",
		TreeAfterSHA256:  "treehash",
	}
	r := fresh()
	if err := AttestExecution(r, "unit", "codex-1", goodRerun); err != nil {
		t.Fatalf("valid attestation must succeed: %v", err)
	}
	got := r.Records[0].Provenance
	if got.Verifier != "codex-1" || got.VerifierRerun == nil ||
		got.VerifierRerun.TreeBeforeSHA256 != r.TreeSHA256 || got.VerifierRerun.Command.CommandSHA256 != "abc" {
		t.Fatalf("attestation must retain verifier, command and tree binding, got %+v", got)
	}
	if err := AttestExecution(r, "integration", "kimi-1", goodRerun); err == nil {
		t.Fatal("self attestation must be refused")
	}
	fail := goodRerun
	fail.Command.ExitCode = 1
	if err := AttestExecution(r, "integration", "codex-1", fail); err == nil {
		t.Fatal("failed re-run must not attest")
	}
	shell := goodRerun
	shell.Command.Format = FormatShell
	shell.Command.ExecutedCases = -1
	shell.Command.FailedCases = -1
	shell.Command.SkippedCases = -1
	if err := AttestExecution(r, "integration", "codex-1", shell); err == nil {
		t.Fatal("opaque shell re-run must not attest")
	}
	zero := goodRerun
	zero.Command.ExecutedCases = 0
	if err := AttestExecution(r, "integration", "codex-1", zero); err == nil {
		t.Fatal("zero-execution re-run must not attest")
	}
	otherCmd := goodRerun
	otherCmd.Command.CommandSHA256 = "unrelated-command"
	if err := AttestExecution(r, "integration", "codex-1", otherCmd); err == nil {
		t.Fatal("re-run of a different command must not attest")
	}
	noOut := goodRerun
	noOut.Command.OutputSHA256 = ""
	if err := AttestExecution(r, "integration", "codex-1", noOut); err == nil {
		t.Fatal("re-run without an output hash must not attest")
	}
	otherTree := goodRerun
	otherTree.TreeBeforeSHA256 = "other-tree"
	otherTree.TreeAfterSHA256 = "other-tree"
	if err := AttestExecution(r, "integration", "codex-1", otherTree); err == nil {
		t.Fatal("re-run against a different tree must not attest")
	}
	unstable := goodRerun
	unstable.TreeAfterSHA256 = "changed-during-rerun"
	if err := AttestExecution(r, "integration", "codex-1", unstable); err == nil {
		t.Fatal("re-run on a tree that changed mid-verification must not attest")
	}
	diverge := goodRerun
	diverge.Command.ExecutedCases = 3
	if err := AttestExecution(r, "integration", "codex-1", diverge); err == nil {
		t.Fatal("re-run with diverging executed-case count must not attest")
	}
	if err := AttestExecution(r, "nonexistent", "codex-1", goodRerun); err == nil {
		t.Fatal("attesting an unknown criterion must fail")
	}
	if err := AttestExecution(r, "integration", "", goodRerun); err == nil {
		t.Fatal("empty verifier identity must fail")
	}
	failedRec := fresh()
	failedRec.Records[1].Status = StatusFail
	if err := AttestExecution(failedRec, "integration", "codex-1", goodRerun); err == nil {
		t.Fatal("attesting a non-pass record must fail")
	}
	noTree := &Report{}
	if err := AttestExecution(noTree, "unit", "codex-1", goodRerun); err == nil {
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

// Adversarial (independent probe, codex-1 2026-09-11): a retained gotest-json
// re-run with a missing (-1) or malformed (<-1) skipped count is refused at
// attestation — a real `go test -json` parse always yields a skipped count, so
// a negative one proves the retained record is not one. Only the envelope
// format may legitimately carry skipped_cases=-1 (it has no such field), and
// valid-but-diverging skipped counts still fail reconciliation.
func TestAttestExecutionRerunCountValidity(t *testing.T) {
	fresh := func() *Report {
		r := positiveReport()
		for i := range r.Records {
			r.Records[i].Provenance.Verifier = ""
			r.Records[i].Provenance.VerifierRerun = nil
		}
		return r
	}
	goodRerun := VerifierExecution{
		Command: CommandEvidence{
			CommandSHA256: "abc", // matches positiveRecord's command hash
			OutputSHA256:  "rerun-out",
			ExitCode:      0,
			Format:        FormatGoTestJSON,
			ExecutedCases: 2,
			FailedCases:   0,
			SkippedCases:  0,
		},
		TreeBeforeSHA256: "treehash",
		TreeAfterSHA256:  "treehash",
	}
	for _, n := range []int{-1, -2, -100} {
		r := fresh()
		bad := goodRerun
		bad.Command.SkippedCases = n
		if err := AttestExecution(r, "unit", "codex-1", bad); err == nil {
			t.Fatalf("gotest-json re-run with skipped_cases=%d must not attest", n)
		}
	}
	// Envelope re-runs legitimately lack a skipped count: -1 is retained.
	r := fresh()
	env := goodRerun
	env.Command.Format = FormatEnvelope
	env.Command.SkippedCases = -1
	if err := AttestExecution(r, "unit", "codex-1", env); err != nil {
		t.Fatalf("envelope re-run without a skipped count must attest: %v", err)
	}
	// ...but a malformed count below -1 is refused in every format.
	r = fresh()
	envBad := env
	envBad.Command.SkippedCases = -2
	if err := AttestExecution(r, "unit", "codex-1", envBad); err == nil {
		t.Fatal("envelope re-run with skipped_cases=-2 must not attest")
	}
	// Valid-but-diverging skipped counts still fail reconciliation.
	r = fresh()
	diverge := goodRerun
	diverge.Command.SkippedCases = 1
	if err := AttestExecution(r, "unit", "codex-1", diverge); err == nil {
		t.Fatal("re-run with diverging skipped-case count must not attest")
	}
}

// Adversarial (independent probe, codex-1 2026-09-11): a PERSISTED retained
// re-run whose skipped count is invalid fails closed at Evaluate, exactly as
// an invalid original record does — format compatibility cannot bypass
// reconciliation.
func TestEvaluateInvalidRerunSkippedCountsRejected(t *testing.T) {
	for _, n := range []int{-1, -2, -100} {
		r := positiveReport()
		r.Records[0].Provenance.VerifierRerun.Command.SkippedCases = n
		if reasons := Evaluate(r, positiveOpts()); len(reasons) == 0 {
			t.Fatalf("retained re-run with skipped_cases=%d must not close", n)
		}
	}
	// The legitimate envelope absence (-1) still closes.
	r := positiveReport()
	r.Records[0].Provenance.VerifierRerun.Command.Format = FormatEnvelope
	r.Records[0].Provenance.VerifierRerun.Command.SkippedCases = -1
	if reasons := Evaluate(r, positiveOpts()); len(reasons) != 0 {
		t.Fatalf("envelope re-run without a skipped count must still close, got: %v", reasons)
	}
}

// Adversarial (independent MAJOR probe, codex-1 2026-09-11): package-level
// fail/build-fail events are structured proof of failure at package scope
// (a failed build runs no test cases), so a pass-claimed gotest-json record
// carrying any, a gotest-json record missing the count (-1), or a malformed
// count (<-1) each fail closed at Evaluate — as does a persisted retained
// re-run whose masked exit laundered a package failure.
func TestEvaluatePackageFailuresRejected(t *testing.T) {
	r := positiveReport()
	r.Records[0].Command.FailedPackages = 2
	if reasons := Evaluate(r, positiveOpts()); !containsAny(reasons, "failed package-level event") {
		t.Fatalf("failed package-level events must be rejected, got: %v", reasons)
	}
	r = positiveReport()
	r.Records[0].Command.FailedPackages = -1
	if reasons := Evaluate(r, positiveOpts()); !containsAny(reasons, "failed-package count absent") {
		t.Fatalf("gotest-json record without a failed-package count must be rejected, got: %v", reasons)
	}
	r = positiveReport()
	r.Records[0].Command.FailedPackages = -2
	if reasons := Evaluate(r, positiveOpts()); !containsAny(reasons, "invalid (negative) case count") {
		t.Fatalf("malformed failed-package count must be rejected, got: %v", reasons)
	}
	r = positiveReport()
	r.Records[0].Provenance.VerifierRerun.Command.FailedPackages = 1
	if reasons := Evaluate(r, positiveOpts()); !containsAny(reasons, "failed package-level event") {
		t.Fatalf("retained re-run with a package failure must not close, got: %v", reasons)
	}
	r = positiveReport()
	r.Records[0].Provenance.VerifierRerun.Command.FailedPackages = -1
	if reasons := Evaluate(r, positiveOpts()); !containsAny(reasons, "carries no failed-package count") {
		t.Fatalf("retained gotest-json re-run without a failed-package count must not close, got: %v", reasons)
	}
	// The legitimate envelope absence (-1) still closes.
	r = positiveReport()
	r.Records[0].Provenance.VerifierRerun.Command.Format = FormatEnvelope
	r.Records[0].Provenance.VerifierRerun.Command.SkippedCases = -1
	r.Records[0].Provenance.VerifierRerun.Command.FailedPackages = -1
	if reasons := Evaluate(r, positiveOpts()); len(reasons) != 0 {
		t.Fatalf("envelope re-run without package counts must still close, got: %v", reasons)
	}
}

// Adversarial (independent MAJOR probe, codex-1 2026-09-11): a verifier's
// independent re-run whose output masks a package/build failure behind exit 0
// is NOT an attestation basis — the retained execution carries the structured
// failure, and a missing or malformed failed-package count on a gotest-json
// re-run proves it did not come from a real parse. Only the envelope format
// may legitimately carry failed_packages=-1.
func TestAttestExecutionRerunPackageFailureRejected(t *testing.T) {
	fresh := func() *Report {
		r := positiveReport()
		for i := range r.Records {
			r.Records[i].Provenance.Verifier = ""
			r.Records[i].Provenance.VerifierRerun = nil
		}
		return r
	}
	goodRerun := VerifierExecution{
		Command: CommandEvidence{
			CommandSHA256: "abc", // matches positiveRecord's command hash
			OutputSHA256:  "rerun-out",
			ExitCode:      0,
			Format:        FormatGoTestJSON,
			ExecutedCases: 2,
			FailedCases:   0,
			SkippedCases:  0,
		},
		TreeBeforeSHA256: "treehash",
		TreeAfterSHA256:  "treehash",
	}
	for _, n := range []int{1, 2} {
		r := fresh()
		bad := goodRerun
		bad.Command.FailedPackages = n
		if err := AttestExecution(r, "unit", "codex-1", bad); err == nil {
			t.Fatalf("re-run with failed_packages=%d (masked package failure) must not attest", n)
		}
	}
	for _, n := range []int{-1, -2} {
		r := fresh()
		bad := goodRerun
		bad.Command.FailedPackages = n
		if err := AttestExecution(r, "unit", "codex-1", bad); err == nil {
			t.Fatalf("gotest-json re-run with failed_packages=%d must not attest", n)
		}
	}
	// Envelope re-runs legitimately lack package counts: -1 is retained.
	r := fresh()
	env := goodRerun
	env.Command.Format = FormatEnvelope
	env.Command.SkippedCases = -1
	env.Command.FailedPackages = -1
	if err := AttestExecution(r, "unit", "codex-1", env); err != nil {
		t.Fatalf("envelope re-run without package counts must attest: %v", err)
	}
}
