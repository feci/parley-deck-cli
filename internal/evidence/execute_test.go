package evidence

import (
	"context"
	"strings"
	"testing"
	"time"
)

// Positive: real `go test -json` output is counted from structured events.
func TestParseGoTestJSONCountsRealCases(t *testing.T) {
	out := `{"Time":"2026-09-10T00:00:00Z","Action":"start","Package":"x"}
{"Time":"2026-09-10T00:00:01Z","Action":"run","Test":"TestA"}
{"Time":"2026-09-10T00:00:01Z","Action":"pass","Test":"TestA"}
{"Time":"2026-09-10T00:00:02Z","Action":"run","Test":"TestB"}
{"Time":"2026-09-10T00:00:02Z","Action":"fail","Test":"TestB"}
{"Time":"2026-09-10T00:00:02Z","Action":"skip","Test":"TestC"}
{"Time":"2026-09-10T00:00:03Z","Action":"fail","Package":"x"}
`
	executed, failed, skipped, ok := ParseGoTestJSON([]byte(out))
	if !ok || executed != 2 || failed != 1 || skipped != 1 {
		t.Fatalf("got executed=%d failed=%d skipped=%d ok=%v", executed, failed, skipped, ok)
	}
}

// Adversarial: plain text mentioning "PASS" or "ok" is NOT test2json and must
// never be counted as executed cases.
func TestParseGoTestJSONRejectsOpaqueText(t *testing.T) {
	for _, out := range []string{
		"PASS\nok  \texample.com/x\t0.01s\n",
		"all tests passed\n",
		"{\"Action\":\"pass\"}\n", // event without a Test: package-level only
		"not json at all",
		"",
	} {
		executed, _, _, ok := ParseGoTestJSON([]byte(out))
		if executed != 0 {
			t.Fatalf("opaque output %q counted %d executed cases", out, executed)
		}
		_ = ok
	}
}

// Adversarial: the word PASS inside opaque shell output yields a shell record
// with unknown counts — never a structured pass.
func TestRunCriterionUnknownOutputIsShellNotCertified(t *testing.T) {
	rec := RunCriterion(context.Background(), t.TempDir(), "c", "echo 'PASS: everything green'", "kimi-1")
	if rec.Status != StatusPass {
		t.Fatalf("exit 0 shell command should be pass-typed, got %q", rec.Status)
	}
	if rec.Command.Format != FormatShell || rec.Command.ExecutedCases != -1 {
		t.Fatalf("opaque output must stay FormatShell with unknown counts, got %+v", rec.Command)
	}
}

// Positive: structured zero-execution proof overrides an exit-0 pass.
func TestRunCriterionZeroExecutionIsNotRun(t *testing.T) {
	cmd := `printf '%s\n' '{"Action":"start","Package":"x"}' '{"Action":"pass","Package":"x"}'`
	rec := RunCriterion(context.Background(), t.TempDir(), "c", cmd, "kimi-1")
	if rec.Status != StatusNotRun || rec.Command.ExecutedCases != 0 {
		t.Fatalf("zero-case test2json must be NotRun, got %q %+v", rec.Status, rec.Command)
	}
}

// Positive: an all-skip structured run is typed skipped, not pass.
func TestRunCriterionAllSkipIsSkipped(t *testing.T) {
	cmd := `printf '%s\n' '{"Action":"run","Test":"TestX"}' '{"Action":"skip","Test":"TestX"}' '{"Action":"pass","Package":"x"}'`
	rec := RunCriterion(context.Background(), t.TempDir(), "c", cmd, "kimi-1")
	if rec.Status != StatusSkipped || rec.Command.SkippedCases != 1 || rec.Command.ExecutedCases != 0 {
		t.Fatalf("all-skip must be Skipped, got %q %+v", rec.Status, rec.Command)
	}
}

// Positive: structured failing cases override exit 0.
func TestRunCriterionFailedCasesOverrideExitZero(t *testing.T) {
	cmd := `printf '%s\n' '{"Action":"run","Test":"TestX"}' '{"Action":"fail","Test":"TestX"}' '{"Action":"pass","Package":"x"}'`
	rec := RunCriterion(context.Background(), t.TempDir(), "c", cmd, "kimi-1")
	if rec.Status != StatusFail || rec.Command.FailedCases != 1 {
		t.Fatalf("failed cases must override exit 0, got %q %+v", rec.Status, rec.Command)
	}
}

// Positive: a real passing go test -json run counts executed cases.
func TestRunCriterionGoTestJSONPass(t *testing.T) {
	cmd := `printf '%s\n' '{"Action":"run","Test":"TestA"}' '{"Action":"pass","Test":"TestA"}' '{"Action":"pass","Package":"x"}'`
	rec := RunCriterion(context.Background(), t.TempDir(), "c", cmd, "kimi-1")
	if rec.Status != StatusPass || rec.Command.Format != FormatGoTestJSON || rec.Command.ExecutedCases != 1 || rec.Command.FailedCases != 0 {
		t.Fatalf("structured pass expected, got %q %+v", rec.Status, rec.Command)
	}
	if rec.Provenance.Executor != "kimi-1" || rec.Provenance.Verifier != "" {
		t.Fatalf("provenance must carry the asserted executor only, got %+v", rec.Provenance)
	}
}

// Adversarial: a failing command is fail-typed even with no structured output.
func TestRunCriterionExitNonzeroFails(t *testing.T) {
	rec := RunCriterion(context.Background(), t.TempDir(), "c", "exit 3", "kimi-1")
	if rec.Status != StatusFail || rec.Command.ExitCode != 3 {
		t.Fatalf("exit 3 must fail, got %q exit=%d", rec.Status, rec.Command.ExitCode)
	}
}

// Envelope: a single well-formed envelope is used.
func TestRunCriterionEnvelopeSingleValid(t *testing.T) {
	cmd := `printf '%s\n' 'PARLEY-EVIDENCE {"executed_cases":4,"failed_cases":0}'`
	rec := RunCriterion(context.Background(), t.TempDir(), "c", cmd, "kimi-1")
	if rec.Command.Format != FormatEnvelope || rec.Command.ExecutedCases != 4 || rec.Status != StatusPass {
		t.Fatalf("valid envelope must pass, got %q %+v", rec.Status, rec.Command)
	}
}

// Adversarial: a single envelope line containing a DUPLICATE JSON field
// (executed_cases 0 then 1) fails closed — last-wins decoding is not evidence.
func TestRunCriterionEnvelopeDuplicateFieldFailsClosed(t *testing.T) {
	cmd := `printf '%s\n' 'PARLEY-EVIDENCE {"executed_cases":0,"executed_cases":1,"failed_cases":0}'`
	rec := RunCriterion(context.Background(), t.TempDir(), "c", cmd, "kimi-1")
	if rec.Status != StatusFail || rec.Command.ExecutedCases != -1 {
		t.Fatalf("duplicate envelope fields must fail closed with no counts, got %q %+v", rec.Status, rec.Command)
	}
}

// Adversarial: duplicate envelope lines fail closed — an earlier valid pass
// must NEVER shadow a later claim (the old last-wins recovery is gone).
func TestRunCriterionEnvelopeDuplicateFailsClosed(t *testing.T) {
	cmd := `printf '%s\n' 'PARLEY-EVIDENCE {"executed_cases":1,"failed_cases":0}' 'PARLEY-EVIDENCE {"executed_cases":4,"failed_cases":0}'`
	rec := RunCriterion(context.Background(), t.TempDir(), "c", cmd, "kimi-1")
	if rec.Status != StatusFail || rec.Command.ExecutedCases != -1 {
		t.Fatalf("duplicate envelopes must fail closed with no counts, got %q %+v", rec.Status, rec.Command)
	}
}

// Adversarial: a trailing MALFORMED envelope after an earlier valid pass fails
// closed and does not fall through to the test2json parser.
func TestRunCriterionMalformedFinalEnvelopeFailsClosed(t *testing.T) {
	cmd := `printf '%s\n' '{"Action":"run","Test":"TestA"}' '{"Action":"pass","Test":"TestA"}' 'PARLEY-EVIDENCE {not json'`
	rec := RunCriterion(context.Background(), t.TempDir(), "c", cmd, "kimi-1")
	if rec.Status != StatusFail || rec.Command.Format != FormatEnvelope || rec.Command.ExecutedCases != -1 {
		t.Fatalf("malformed final envelope must fail closed, got %q %+v", rec.Status, rec.Command)
	}
}

// Adversarial: null/partial/negative/conflicting envelope counts all fail closed,
// as do duplicate, case-aliased, unknown and trailing fields — permissive
// decoding is never evidence of execution.
func TestParseEnvelopeFailClosed(t *testing.T) {
	cases := map[string]string{
		"malformed":       "PARLEY-EVIDENCE {not json}",
		"not json":        "PARLEY-EVIDENCE hello",
		"null counts":     `PARLEY-EVIDENCE {"executed_cases":null,"failed_cases":null}`,
		"partial":         `PARLEY-EVIDENCE {"executed_cases":3}`,
		"negative":        `PARLEY-EVIDENCE {"executed_cases":-2,"failed_cases":0}`,
		"conflicting":     `PARLEY-EVIDENCE {"executed_cases":1,"failed_cases":2}`,
		"empty object":    `PARLEY-EVIDENCE {}`,
		"duplicate field": `PARLEY-EVIDENCE {"executed_cases":0,"executed_cases":1,"failed_cases":0}`,
		"aliased field":   `PARLEY-EVIDENCE {"Executed_Cases":1,"failed_cases":0}`,
		"unknown field":   `PARLEY-EVIDENCE {"executed_cases":1,"failed_cases":0,"extra":9}`,
		"string count":    `PARLEY-EVIDENCE {"executed_cases":"3","failed_cases":0}`,
		"fractional":      `PARLEY-EVIDENCE {"executed_cases":1.5,"failed_cases":0}`,
		"trailing object": `PARLEY-EVIDENCE {"executed_cases":1,"failed_cases":0} {"executed_cases":9}`,
		"trailing junk":   `PARLEY-EVIDENCE {"executed_cases":1,"failed_cases":0} trailing`,
		"array payload":   `PARLEY-EVIDENCE [1,0]`,
	}
	for name, out := range cases {
		env, present, err := ParseEnvelope(out)
		if !present || err == nil {
			t.Fatalf("%s: must be present and fail closed, got present=%v err=%v env=%+v", name, present, err, env)
		}
	}
	if _, present, err := ParseEnvelope("no envelope here\n"); present || err != nil {
		t.Fatalf("no envelope line must be absent, got present=%v err=%v", present, err)
	}
	if env, present, err := ParseEnvelope(`PARLEY-EVIDENCE {"executed_cases":2,"failed_cases":1}`); !present || err != nil || *env.ExecutedCases != 2 || *env.FailedCases != 1 {
		t.Fatalf("valid envelope must parse, got present=%v err=%v env=%+v", present, err, env)
	}
}

// Concurrency fixture (serial-vs-barrier): two sequential runs of the same
// stateful command must observe each other's effects — proving reruns are
// serialized through a barrier, not isolated. The records must still bind
// DISTINCT output hashes, so no cross-contamination between records.
func TestRunCriterionSerialVsBarrierFixture(t *testing.T) {
	dir := t.TempDir()
	cmd := `n=$(cat counter 2>/dev/null || echo 0); n=$((n+1)); echo "$n" > counter; printf '%s\n' "{\"Action\":\"run\",\"Test\":\"TestA$n\"}" "{\"Action\":\"pass\",\"Test\":\"TestA$n\"}"`
	first := RunCriterion(context.Background(), dir, "c", cmd, "kimi-1")
	second := RunCriterion(context.Background(), dir, "c", cmd, "kimi-1")
	if first.Status != StatusPass || second.Status != StatusPass {
		t.Fatalf("both serial runs must pass, got %q then %q", first.Status, second.Status)
	}
	if first.Command.OutputSHA256 == second.Command.OutputSHA256 {
		t.Fatal("serial reruns through a barrier must bind distinct outputs")
	}
	if !strings.Contains(second.Command.Diagnostics, "TestA2") {
		t.Fatalf("second run must observe the first run's state (barrier, not isolation): %q", second.Command.Diagnostics)
	}
}

// Adversarial: secret-shaped tokens in command output are scrubbed from the
// persisted diagnostics.
func TestRunCriterionScrubsSecrets(t *testing.T) {
	rec := RunCriterion(context.Background(), t.TempDir(), "c",
		"echo 'token=sk-abcdefghijklmnop1234567890'; echo 'Authorization: Bearer abcREALtoken1234567890'", "kimi-1")
	for _, leak := range []string{"sk-abcdefghijklmnop", "REALtoken"} {
		if strings.Contains(rec.Command.Diagnostics, leak) {
			t.Fatalf("diagnostics leaked %q: %q", leak, rec.Command.Diagnostics)
		}
		if strings.Contains(rec.Command.Command, leak) {
			t.Fatalf("persisted command representation leaked %q: %q", leak, rec.Command.Command)
		}
	}
}

// Positive: command and output hashes are always populated.
func TestRunCriterionHashesBound(t *testing.T) {
	rec := RunCriterion(context.Background(), t.TempDir(), "c", "echo hi", "kimi-1")
	if rec.Command.CommandSHA256 == "" || rec.Command.OutputSHA256 == "" {
		t.Fatalf("hash binding missing: %+v", rec.Command)
	}
	if rec.Command.DurationMillis < 0 || rec.Command.DurationMillis > int64((10*time.Second)/time.Millisecond) {
		t.Fatalf("implausible duration: %d", rec.Command.DurationMillis)
	}
}
