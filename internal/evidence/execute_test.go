package evidence

import (
	"context"
	"strings"
	"testing"
	"time"
)

// Positive: real `go test -json` output is counted from structured events.
// The stream also fails package x at package scope (real shape: TestB's
// test-level fail fails the package) — the strict parser keeps the EXACT
// test-case counts and surfaces the package failure separately, while the
// lossy exported wrapper zeroes its counts (fail closed, never a pass).
func TestParseGoTestJSONCountsRealCases(t *testing.T) {
	out := `{"Time":"2026-09-10T00:00:00Z","Action":"start","Package":"x"}
{"Time":"2026-09-10T00:00:01Z","Action":"run","Test":"TestA"}
{"Time":"2026-09-10T00:00:01Z","Action":"pass","Test":"TestA"}
{"Time":"2026-09-10T00:00:02Z","Action":"run","Test":"TestB"}
{"Time":"2026-09-10T00:00:02Z","Action":"fail","Test":"TestB"}
{"Time":"2026-09-10T00:00:02Z","Action":"skip","Test":"TestC"}
{"Time":"2026-09-10T00:00:03Z","Action":"fail","Package":"x"}
`
	executed, failed, skipped, failedPackages, recognized, err := parseGoTestJSONEvidence([]byte(out))
	if err != nil || !recognized || executed != 2 || failed != 1 || skipped != 1 || failedPackages != 1 {
		t.Fatalf("strict parse: got executed=%d failed=%d skipped=%d failedPackages=%d recognized=%v err=%v",
			executed, failed, skipped, failedPackages, recognized, err)
	}
	executed, failed, skipped, ok := ParseGoTestJSON([]byte(out))
	if !ok || executed != 0 || failed != 0 || skipped != 0 {
		t.Fatalf("package failure must zero the wrapper counts (fail closed), got executed=%d failed=%d skipped=%d ok=%v",
			executed, failed, skipped, ok)
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

// Adversarial (independent probe, codex-1 2026-09-11): the capture overflow
// decision uses PRE-WRITE remaining capacity — an exact-bound write (one 16B
// chunk or split 8+8) and any under-bound write (9B) are retained in full
// without overflow; only a chunk larger than the room left overflows.
func TestCappedWriterBoundaryAccounting(t *testing.T) {
	cases := []struct {
		name     string
		max      int
		chunks   []string
		overflow bool
		retained int
	}{
		{"exact bound single chunk", 16, []string{strings.Repeat("x", 16)}, false, 16},
		{"exact bound split chunks", 16, []string{strings.Repeat("x", 8), strings.Repeat("x", 8)}, false, 16},
		{"under bound single chunk", 16, []string{strings.Repeat("x", 9)}, false, 9},
		{"one byte over single chunk", 16, []string{strings.Repeat("x", 17)}, true, 16},
		{"one byte over split chunks", 16, []string{strings.Repeat("x", 16), "y"}, true, 16},
		{"split chunk crossing the bound", 16, []string{strings.Repeat("x", 8), strings.Repeat("y", 9)}, true, 16},
		{"empty write at the bound", 16, []string{strings.Repeat("x", 16), ""}, false, 16},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			w := &cappedWriter{max: tc.max}
			for _, chunk := range tc.chunks {
				n, err := w.Write([]byte(chunk))
				if n != len(chunk) || err != nil {
					t.Fatalf("Write must consume the full chunk, got n=%d err=%v", n, err)
				}
			}
			if w.overflow != tc.overflow {
				t.Fatalf("overflow=%v, want %v", w.overflow, tc.overflow)
			}
			if w.buf.Len() != tc.retained {
				t.Fatalf("retained %d bytes, want %d", w.buf.Len(), tc.retained)
			}
		})
	}
}

// Adversarial (independent probe, codex-1 2026-09-11): duplicate or
// case-aliased test2json semantic fields, non-string semantic values, trailing
// content, and malformed event lines are RECOGNIZED-but-INVALID — the parser
// yields an error and no usable counts, never a collapsed pass and never an
// opaque shell fallback.
func TestParseGoTestJSONStrictFailClosed(t *testing.T) {
	cases := map[string]string{
		"duplicate action fail-to-pass": `{"Action":"fail","Action":"pass","Test":"TestX"}`,
		"duplicate action pass-to-fail": `{"Action":"pass","Action":"fail","Test":"TestX"}`,
		"duplicate test field":          `{"Action":"pass","Test":"TestA","Test":"TestB"}`,
		"case-aliased action":           `{"action":"pass","Test":"TestX"}`,
		"case-aliased test":             `{"Action":"pass","TEST":"TestX"}`,
		"non-string action":             `{"Action":1,"Test":"TestX"}`,
		"non-string test":               `{"Action":"pass","Test":2}`,
		"trailing second object":        `{"Action":"pass","Test":"TestX"} {"Action":"fail","Test":"TestY"}`,
		"trailing junk":                 `{"Action":"pass","Test":"TestX"} junk`,
		"malformed braced line":         `{not json}`,
		"malformed line amid events":    "{\"Action\":\"run\",\"Test\":\"TestX\"}\n{broken\n{\"Action\":\"pass\",\"Test\":\"TestX\"}",
	}
	for name, out := range cases {
		_, _, _, failedPackages, recognized, err := parseGoTestJSONEvidence([]byte(out))
		if !recognized || err == nil || failedPackages != 0 {
			t.Fatalf("%s: must be recognized-but-invalid with zeroed package failures, got recognized=%v err=%v failedPackages=%d", name, recognized, err, failedPackages)
		}
		executed, failed, _, ok := ParseGoTestJSON([]byte(out))
		if ok != recognized || executed != 0 || failed != 0 {
			t.Fatalf("%s: strict rejection must yield zero counts, got executed=%d failed=%d recognized=%v", name, executed, failed, ok)
		}
	}
}

// Positive + adversarial: normal real test2json streams are preserved — the
// full real field set (Time, Package, Output, Elapsed), duplicate NON-semantic
// fields (no case effect), unknown future fields with nested values, and plain
// noise lines around the events. But the package-level build-fail and fail
// events for example.com/y are NOT harmless noise: the strict parser keeps the
// EXACT test-case counts (a failed build runs no test cases — none are
// invented) while surfacing 2 package-level failure events, and the lossy
// exported wrapper zeroes its counts so such a stream can never be a pass
// basis through it.
func TestParseGoTestJSONRealStreamsPreserved(t *testing.T) {
	out := `# example.com/x
{"Time":"2026-09-11T00:00:00Z","Action":"start","Package":"example.com/x"}
{"Time":"2026-09-11T00:00:01Z","Action":"run","Package":"example.com/x","Test":"TestA"}
{"Time":"2026-09-11T00:00:01Z","Action":"output","Package":"example.com/x","Test":"TestA","Output":"=== RUN   TestA\n"}
{"Time":"2026-09-11T00:00:02Z","Action":"pass","Package":"example.com/x","Test":"TestA","Elapsed":0.01}
{"Time":"2026-09-11T00:00:02Z","Action":"output","Package":"example.com/x","Test":"TestA","Output":"a\n","Output":"b\n"}
{"Time":"2026-09-11T00:00:03Z","Action":"skip","Package":"example.com/x","Test":"TestB","Elapsed":0}
{"Time":"2026-09-11T00:00:04Z","Action":"pass","Package":"example.com/x","Test":"TestC","FutureField":{"nested":[1,{"x":2}]}}
{"Time":"2026-09-11T00:00:05Z","Action":"pass","Package":"example.com/x","Elapsed":0.02}
FAIL	example.com/y [build failed]
{"Time":"2026-09-11T00:00:06Z","Action":"build-fail","Package":"example.com/y"}
{"Time":"2026-09-11T00:00:06Z","Action":"fail","Package":"example.com/y","Elapsed":0}
`
	executed, failed, skipped, failedPackages, recognized, err := parseGoTestJSONEvidence([]byte(out))
	if err != nil || !recognized {
		t.Fatalf("real test2json stream must parse cleanly, got recognized=%v err=%v", recognized, err)
	}
	if executed != 2 || failed != 0 || skipped != 1 {
		t.Fatalf("exact test-case counts must be preserved, got executed=%d failed=%d skipped=%d, want 2/0/1", executed, failed, skipped)
	}
	if failedPackages != 2 {
		t.Fatalf("package-level fail + build-fail events must surface as 2 failed-package events, got %d", failedPackages)
	}
	if executed, failed, skipped, ok := ParseGoTestJSON([]byte(out)); !ok || executed != 0 || failed != 0 || skipped != 0 {
		t.Fatalf("wrapper must fail closed (zeroed counts) on package failure, got executed=%d failed=%d skipped=%d ok=%v", executed, failed, skipped, ok)
	}
}

// Adversarial (independent probe, codex-1 2026-09-11): a test2json stream with
// a duplicated Action field is recognized-but-invalid — RunCriterion types it
// fail with the gotest-json format and UNKNOWN counts; it never collapses
// into a shell PASS.
func TestRunCriterionMalformedGoTestJSONFailsClosed(t *testing.T) {
	cmd := `printf '%s\n' '{"Action":"fail","Action":"pass","Test":"TestX"}'`
	rec := RunCriterion(context.Background(), t.TempDir(), "c", cmd, "kimi-1")
	if rec.Status != StatusFail || rec.Command.Format != FormatGoTestJSON {
		t.Fatalf("tampered test2json must fail closed as gotest-json, got %q %+v", rec.Status, rec.Command)
	}
	if rec.Command.ExecutedCases != -1 || rec.Command.FailedCases != -1 || rec.Command.SkippedCases != -1 {
		t.Fatalf("tampered stream must keep counts unknown, got %+v", rec.Command)
	}
	if !strings.Contains(rec.Command.Diagnostics, "invalid test2json event stream") {
		t.Fatalf("diagnostics must record the strict rejection, got %q", rec.Command.Diagnostics)
	}
}

// Positive: package-level events and plain noise lines around real test2json
// events (the build-failure shape) do not disturb case accounting.
func TestRunCriterionGoTestJSONNoiseTolerated(t *testing.T) {
	cmd := `printf '%s\n' '# example.com/x' '{"Action":"start","Package":"x"}' '{"Action":"run","Test":"TestA"}' '{"Action":"pass","Test":"TestA"}' '{"Action":"pass","Package":"x"}'`
	rec := RunCriterion(context.Background(), t.TempDir(), "c", cmd, "kimi-1")
	if rec.Status != StatusPass || rec.Command.ExecutedCases != 1 || rec.Command.FailedCases != 0 {
		t.Fatalf("noise around real events must not disturb accounting, got %q %+v", rec.Status, rec.Command)
	}
	if rec.Command.FailedPackages != 0 {
		t.Fatalf("a clean stream must carry a zero failed-package count, got %+v", rec.Command)
	}
}

// Adversarial (independent MAJOR probe, codex-1 2026-09-11): a zero-exit
// command emitting a passing test from one package PLUS a package-level fail
// or build-fail event from another became an overall PASS — a masked pipeline
// exit laundered a package/build failure. Package-level failure events are
// structured proof of failure at package scope and must fail closed, while
// the exact real test-case counts are preserved (no executed test is invented
// for a failed build).
func TestRunCriterionPackageFailureFailsClosed(t *testing.T) {
	cases := []struct {
		name           string
		cmd            string
		executed       int
		failed         int
		failedPackages int
	}{
		{
			name:           "masked exit package fail",
			cmd:            `printf '%s\n' '{"Action":"pass","Test":"TestGood","Package":"good"}' '{"Action":"fail","Package":"bad"}'`,
			executed:       1,
			failed:         0,
			failedPackages: 1,
		},
		{
			name:           "masked exit build-fail",
			cmd:            `printf '%s\n' '{"Action":"pass","Test":"TestGood","Package":"good"}' '{"Action":"build-fail","Package":"bad"}'`,
			executed:       1,
			failed:         0,
			failedPackages: 1,
		},
		{
			// The full real `go test -json ./...` failure shape (captured from
			// go1.27.1): build-output noise, an ImportPath build-fail, a
			// package-level fail carrying FailedBuild, and the good package
			// passing — with the failing producer's non-zero exit MASKED by a
			// pipe (the pipeline's exit status is cat's 0).
			name: "masked pipeline full real stream",
			cmd: `(printf '%s\n' ` +
				`'{"ImportPath":"bad [bad.test]","Action":"build-output","Output":"# bad [bad.test]\n"}' ` +
				`'{"ImportPath":"bad [bad.test]","Action":"build-fail"}' ` +
				`'{"Action":"start","Package":"bad"}' ` +
				`'{"Action":"fail","Package":"bad","Elapsed":0.001,"FailedBuild":"bad [bad.test]"}' ` +
				`'{"Action":"start","Package":"good"}' ` +
				`'{"Action":"run","Package":"good","Test":"TestGood"}' ` +
				`'{"Action":"pass","Package":"good","Test":"TestGood","Elapsed":0}' ` +
				`'{"Action":"pass","Package":"good","Elapsed":0.3}' ` +
				`'FAIL' ; exit 1) | cat`,
			executed:       1,
			failed:         0,
			failedPackages: 2, // build-fail + the package-level fail it causes
		},
		{
			name:           "package failure with no test events is fail, not not-run",
			cmd:            `printf '%s\n' '{"Action":"start","Package":"bad"}' '{"Action":"fail","Package":"bad","Elapsed":0}'`,
			executed:       0,
			failed:         0,
			failedPackages: 1,
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			rec := RunCriterion(context.Background(), t.TempDir(), "c", tc.cmd, "kimi-1")
			if rec.Status != StatusFail || rec.Command.Format != FormatGoTestJSON {
				t.Fatalf("package/build failure must fail closed as gotest-json, got %q %+v", rec.Status, rec.Command)
			}
			if rec.Command.ExitCode != 0 {
				t.Fatalf("fixture must mask the exit code (exit 0), got %d", rec.Command.ExitCode)
			}
			if rec.Command.FailedPackages != tc.failedPackages {
				t.Fatalf("failedPackages=%d, want %d: %+v", rec.Command.FailedPackages, tc.failedPackages, rec.Command)
			}
			if rec.Command.ExecutedCases != tc.executed || rec.Command.FailedCases != tc.failed {
				t.Fatalf("exact test-case counts must be preserved: got %d/%d, want %d/%d",
					rec.Command.ExecutedCases, rec.Command.FailedCases, tc.executed, tc.failed)
			}
			if !strings.Contains(rec.Command.Diagnostics, "package-level failure event") {
				t.Fatalf("diagnostics must record the package failure, got %q", rec.Command.Diagnostics)
			}
		})
	}
}

// Positive: normal multi-package success is preserved — interleaved package
// start/output/pass events across two packages are not failures, and the
// record carries an exact zero failed-package count.
func TestRunCriterionMixedPassingPackagesPreserved(t *testing.T) {
	cmd := `printf '%s\n' ` +
		`'{"Action":"start","Package":"good"}' ` +
		`'{"Action":"start","Package":"other"}' ` +
		`'{"Action":"run","Package":"good","Test":"TestGood"}' ` +
		`'{"Action":"output","Package":"good","Test":"TestGood","Output":"=== RUN   TestGood\n"}' ` +
		`'{"Action":"pass","Package":"good","Test":"TestGood","Elapsed":0}' ` +
		`'{"Action":"pass","Package":"good","Elapsed":0.1}' ` +
		`'{"Action":"run","Package":"other","Test":"TestOther"}' ` +
		`'{"Action":"pass","Package":"other","Test":"TestOther","Elapsed":0.02}' ` +
		`'{"Action":"pass","Package":"other","Elapsed":0.2}'`
	rec := RunCriterion(context.Background(), t.TempDir(), "c", cmd, "kimi-1")
	if rec.Status != StatusPass || rec.Command.Format != FormatGoTestJSON {
		t.Fatalf("mixed passing packages must stay a structured pass, got %q %+v", rec.Status, rec.Command)
	}
	if rec.Command.ExecutedCases != 2 || rec.Command.FailedCases != 0 || rec.Command.FailedPackages != 0 {
		t.Fatalf("exact counts 2/0 with zero failed packages expected, got %+v", rec.Command)
	}
}
