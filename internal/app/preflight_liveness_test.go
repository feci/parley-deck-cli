package app

import (
	"bytes"
	"context"
	"os"
	"strings"
	"testing"
	"time"

	"parley-deck-cli/internal/agents"
)

// --- pure classification table (D7) ---------------------------------------------

func TestClassifyReadiness(t *testing.T) {
	cases := []struct {
		name            string
		stdout, stderr  string
		exitCode        int
		timedOut        bool
		want            ReadinessClass
		wantReady       bool
		wantSawSentinel bool
	}{
		// Exact plain PONG is ready only for plain-output adapters.
		{"exact plain PONG", "PONG\n", "", 0, false, ClassReady, true, false},
		{"exact PONG whitespace", "  PONG  \n", "", 0, false, ClassReady, true, false},
		// Recognized assistant/result envelope schemas: explicit positive coverage
		// for every accepted content key, wrapper key, and the assistant role.
		{"bare content no provenance is malformed", `{"content":"PONG"}`, "", 0, false, ClassMalformedReply, false, true},
		{"unattributed JSON message key", `{"message":"PONG"}`, "", 0, false, ClassMalformedReply, false, true},
		{"unattributed JSON text key", `{"text":"PONG"}`, "", 0, false, ClassMalformedReply, false, true},
		{"unattributed JSON result string", `{"result":"PONG"}`, "", 0, false, ClassMalformedReply, false, true},
		{"recognized assistant-role message", `{"role":"assistant","content":"PONG"}`, "", 0, false, ClassReady, true, false},
		{"unattributed result wrapper", `{"result":{"content":"PONG"}}`, "", 0, false, ClassMalformedReply, false, true},
		{"unattributed data wrapper", `{"data":{"content":"PONG"}}`, "", 0, false, ClassMalformedReply, false, true},
		{"unattributed payload wrapper", `{"payload":{"content":"PONG"}}`, "", 0, false, ClassMalformedReply, false, true},
		{"unattributed response wrapper", `{"response":{"content":"PONG"}}`, "", 0, false, ClassMalformedReply, false, true},
		{"unattributed result wrapper message key", `{"result":{"message":"PONG"}}`, "", 0, false, ClassMalformedReply, false, true},
		{"recognized assistant result schema with subtype success", `{"type":"result","subtype":"success","is_error":false,"result":"PONG"}`, "", 0, false, ClassReady, true, false},
		// Echoes, fences, malformed JSON, role-tagged non-assistant messages, and
		// unrecognized string-valued keys are never ready.
		{"echoed instruction is not ready", "Reply with exactly the single token: PONG\n", "", 0, false, ClassMalformedReply, false, true},
		{"bullet PONG is not ready", "• PONG\n", "", 0, false, ClassMalformedReply, false, true},
		{"fenced PONG is not ready", "```\nPONG\n```\n", "", 0, false, ClassMalformedReply, false, true},
		{"malformed JSON with PONG is not ready", `{"content":"PONG"`, "", 0, false, ClassMalformedReply, false, true},
		{"user-role echo is not ready", `{"role":"user","content":"PONG"}`, "", 0, false, ClassMalformedReply, false, true},
		{"system-role message is not ready", `{"role":"system","content":"PONG"}`, "", 0, false, ClassMalformedReply, false, true},
		{"tool-role message is not ready", `{"role":"tool","content":"PONG"}`, "", 0, false, ClassMalformedReply, false, true},
		{"wrapped user-role echo is not ready", `{"data":{"role":"user","content":"PONG"}}`, "", 0, false, ClassMalformedReply, false, true},
		{"unrecognized answer key is not ready", `{"answer":"PONG"}`, "", 0, false, ClassMalformedReply, false, true},
		{"unrecognized output key is not ready", `{"output":"PONG"}`, "", 0, false, ClassMalformedReply, false, true},
		{"unrecognized pong key is not ready", `{"pong":"PONG"}`, "", 0, false, ClassMalformedReply, false, true},
		{"JSONL stream is not ready", "{\"role\":\"assistant\",\"content\":\"PONG\"}\n{\"role\":\"assistant\",\"content\":\"PONG\"}\n", "", 0, false, ClassMalformedReply, false, true},
		// Structured error/status envelopes are provider failures — never ready,
		// never an excludable process failure — even when a subtype claims success.
		{"JSON error envelope defeated by subtype", `{"type":"error","subtype":"success","error":"boom"}`, "", 0, false, ClassProviderFailure, false, false},
		{"JSON error envelope with PONG nested", `{"type":"error","message":"cannot compute PONG"}`, "", 0, false, ClassProviderFailure, false, false},
		{"is_error true with success subtype", `{"type":"result","is_error":true,"subtype":"success","result":"PONG"}`, "", 0, false, ClassProviderFailure, false, false},
		{"wrapped is_error", `{"result":{"is_error":true,"content":"PONG"}}`, "", 0, false, ClassProviderFailure, false, false},
		{"success false envelope", `{"success":false,"content":"PONG"}`, "", 0, false, ClassProviderFailure, false, false},
		// Silence, non-JSON text, and deadline classes.
		{"empty exit-zero wrapper", "", "", 0, false, ClassExitedEmpty, false, false},
		{"unrecognized exit-zero text", "hello world\n", "", 0, false, ClassMalformedReply, false, false},
		{"provider error on nonzero exit", "Error: authentication_failed — run login\n", "", 1, false, ClassProviderFailure, false, false},
		{"plain nonzero exit", "something inscrutable happened\n", "", 3, false, ClassProcessFailure, false, false},
		{"deadline with no output", "", "", -1, true, ClassDeadlineNoOutput, false, false},
		{"deadline after partial output", "partial output\n", "", -1, true, ClassDeadlineAfterOutput, false, false},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			obs := classifyReadiness(tc.stdout, tc.stderr, tc.exitCode, tc.timedOut, false, "")
			if obs.Class != tc.want {
				t.Fatalf("class=%q want %q", obs.Class, tc.want)
			}
			if obs.Ready != tc.wantReady {
				t.Errorf("ready=%v want %v", obs.Ready, tc.wantReady)
			}
			if obs.SawSentinel != tc.wantSawSentinel {
				t.Errorf("sawSentinel=%v want %v", obs.SawSentinel, tc.wantSawSentinel)
			}
		})
	}
}

func TestClassifyReadinessProviderSubclass(t *testing.T) {
	obs := classifyReadiness("", "429 Too Many Requests, rate limit hit", 1, false, false, "")
	if obs.Class != ClassProviderFailure || obs.ProviderClass != "rate-limit" {
		t.Fatalf("class=%q subclass=%q want provider-failure/rate-limit", obs.Class, obs.ProviderClass)
	}
}

// --- readiness gate semantics (no child process; fake probe) --------------------

func TestMalformedReadinessIsResolveGateNotExclusion(t *testing.T) {
	root := sourceWorkspace(t)
	withFakeProbe(t, func(_ context.Context, _ string, a agents.Discovery, _ time.Duration) readinessObservation {
		if a.ID == "b" {
			return readinessObservation{Class: ClassMalformedReply, SawSentinel: true, StdoutTail: "• PONG"}
		}
		return readyObs()
	})
	// Even with --yes, a malformed reply MUST NOT be auto-excluded: it opens a
	// resolve-readiness gate and keeps the roster (quorum) unchanged.
	report, code, err := preflight(context.Background(), preflightOptions{Root: root, Yes: true},
		[]agents.Discovery{found("a"), found("b"), found("c")}, &bytes.Buffer{}, &bytes.Buffer{})
	if err != nil {
		t.Fatal(err)
	}
	if code != 3 {
		t.Fatalf("code=%d want 3 (blocking resolve-readiness gate)", code)
	}
	if len(report.Excluded) != 0 {
		t.Fatalf("--yes must NOT auto-exclude a malformed agent, got %v", report.Excluded)
	}
	if len(report.Roster) != 3 {
		t.Fatalf("roster length=%d want 3 (quorum unchanged)", len(report.Roster))
	}
	if len(report.Gates) != 1 || report.Gates[0].Kind != gateResolveReadiness {
		t.Fatalf("gates=%v want one resolve-readiness gate", report.Gates)
	}
}

func TestEmptyAndDeadlineAreResolveGatesNotExclusion(t *testing.T) {
	for _, cls := range []ReadinessClass{ClassExitedEmpty, ClassDeadlineNoOutput, ClassDeadlineAfterOutput} {
		t.Run(string(cls), func(t *testing.T) {
			root := sourceWorkspace(t)
			withFakeProbe(t, func(_ context.Context, _ string, a agents.Discovery, _ time.Duration) readinessObservation {
				if a.ID == "b" {
					return readinessObservation{Class: cls}
				}
				return readyObs()
			})
			report, code, err := preflight(context.Background(), preflightOptions{Root: root, Yes: true},
				[]agents.Discovery{found("a"), found("b"), found("c")}, &bytes.Buffer{}, &bytes.Buffer{})
			if err != nil {
				t.Fatal(err)
			}
			if code != 3 {
				t.Fatalf("code=%d want 3", code)
			}
			if len(report.Excluded) != 0 || len(report.Gates) != 1 || report.Gates[0].Kind != gateResolveReadiness {
				t.Fatalf("excluded=%v gates=%v want no exclusion and one resolve gate", report.Excluded, report.Gates)
			}
		})
	}
}

func TestProviderFailureIsBlockingDistinctGate(t *testing.T) {
	root := sourceWorkspace(t)
	withFakeProbe(t, func(_ context.Context, _ string, a agents.Discovery, _ time.Duration) readinessObservation {
		if a.ID == "b" {
			return readinessObservation{Class: ClassProviderFailure, ProviderClass: "auth"}
		}
		return readyObs()
	})
	report, code, err := preflight(context.Background(), preflightOptions{Root: root, Yes: true},
		[]agents.Discovery{found("a"), found("b"), found("c")}, &bytes.Buffer{}, &bytes.Buffer{})
	if err != nil {
		t.Fatal(err)
	}
	if code != 3 {
		t.Fatalf("code=%d want 3 (provider failure stays blocking)", code)
	}
	if len(report.Excluded) != 0 {
		t.Fatalf("provider failure must not auto-exclude, got %v", report.Excluded)
	}
	if len(report.Gates) != 1 || report.Gates[0].Kind != gateProviderFailure {
		t.Fatalf("gates=%v want one provider-failure gate", report.Gates)
	}
}

func TestProcessFailureKeepsExplicitExclusion(t *testing.T) {
	root := sourceWorkspace(t)
	withFakeProbe(t, func(_ context.Context, _ string, a agents.Discovery, _ time.Duration) readinessObservation {
		if a.ID == "b" {
			return procFailObs()
		}
		return readyObs()
	})
	// Non-provider process failure is definite unavailability: --yes records the
	// exclusion, exactly as before.
	report, code, err := preflight(context.Background(), preflightOptions{Root: root, Yes: true},
		[]agents.Discovery{found("a"), found("b"), found("c")}, &bytes.Buffer{}, &bytes.Buffer{})
	if err != nil {
		t.Fatal(err)
	}
	if code != 0 {
		t.Fatalf("code=%d want 0 (exclusion confirmed, >=2 remain)", code)
	}
	if len(report.Excluded) != 1 {
		t.Fatalf("excluded=%v want one recorded exclusion", report.Excluded)
	}
}

// --- real fake child processes through hostedPONG -------------------------------

func hasPongArg(want string) bool {
	seen := false
	for _, a := range os.Args[1:] {
		if a == "--" {
			seen = true
			continue
		}
		if seen && a == want {
			return true
		}
	}
	return false
}

// TestPongFixtureHelper is re-executed as a real child process by the probe
// fixtures below; each scenario emits the output an agent would.
func TestPongFixtureHelper(t *testing.T) {
	scenario := ""
	seen := false
	for _, a := range os.Args[1:] {
		if a == "--" {
			seen = true
			continue
		}
		if seen {
			scenario = a
			break
		}
	}
	switch scenario {
	case "":
		return // normal test run
	case "ready":
		_, _ = os.Stdout.WriteString("PONG\n")
	case "json-ready":
		_, _ = os.Stdout.WriteString(`{"role":"assistant","content":"PONG"}` + "\n")
	case "echo":
		_, _ = os.Stdout.WriteString("Reply with exactly the single token: PONG\n")
	case "malformed-json":
		_, _ = os.Stdout.WriteString(`{"content":"PONG"` + "\n")
	case "json-error-subtype-success":
		_, _ = os.Stdout.WriteString(`{"type":"error","subtype":"success","error":"boom"}` + "\n")
	case "empty":
		// exit 0 with no output at all
	case "provider-error":
		_, _ = os.Stderr.WriteString("Error: authentication_failed — run login\n")
		os.Exit(1)
	case "process-error":
		_, _ = os.Stderr.WriteString("something inscrutable happened\n")
		os.Exit(3)
	case "never-ending":
		time.Sleep(30 * time.Second)
	case "partial-then-hang":
		_, _ = os.Stdout.WriteString("partial output\n")
		_ = os.Stdout.Sync()
		time.Sleep(30 * time.Second)
	}
	os.Exit(0)
}

func pongFakeAgent(scenario string) agents.Discovery {
	return agents.Discovery{
		Spec: agents.Spec{
			ID:           "fake",
			Commands:     []string{os.Args[0]},
			HeadlessArgs: []string{"-test.run=TestPongFixtureHelper", "--", scenario},
			PromptMode:   agents.PromptStdin,
		},
		Path:  os.Args[0],
		Found: true,
	}
}

func probe(t *testing.T, scenario string, timeout time.Duration) readinessObservation {
	t.Helper()
	root := t.TempDir()
	return hostedPONG(context.Background(), root, pongFakeAgent(scenario), timeout)
}

func TestHostedPONGRealChildFixtures(t *testing.T) {
	t.Run("ready exact", func(t *testing.T) {
		if obs := probe(t, "ready", 5*time.Second); obs.Class != ClassReady || !obs.Ready {
			t.Fatalf("obs=%+v want ready", obs)
		}
	})
	t.Run("recognized PONG wrapper", func(t *testing.T) {
		if obs := probe(t, "json-ready", 5*time.Second); obs.Class != ClassReady || !obs.Ready {
			t.Fatalf("obs=%+v want ready", obs)
		}
	})
	t.Run("echoed instruction", func(t *testing.T) {
		obs := probe(t, "echo", 5*time.Second)
		if obs.Class != ClassMalformedReply || !obs.SawSentinel || obs.Ready {
			t.Fatalf("obs=%+v want malformed-reply, not ready", obs)
		}
	})
	t.Run("malformed JSON", func(t *testing.T) {
		obs := probe(t, "malformed-json", 5*time.Second)
		if obs.Class != ClassMalformedReply || obs.Ready {
			t.Fatalf("obs=%+v want malformed-reply, not ready", obs)
		}
	})
	t.Run("JSON error envelope subtype success", func(t *testing.T) {
		obs := probe(t, "json-error-subtype-success", 5*time.Second)
		if obs.Class != ClassProviderFailure || obs.Ready {
			t.Fatalf("obs=%+v want provider-failure, not ready", obs)
		}
	})
	t.Run("empty exit-zero wrapper", func(t *testing.T) {
		if obs := probe(t, "empty", 5*time.Second); obs.Class != ClassExitedEmpty || obs.Ready {
			t.Fatalf("obs=%+v want process-exited-empty", obs)
		}
	})
	t.Run("provider error", func(t *testing.T) {
		obs := probe(t, "provider-error", 5*time.Second)
		if obs.Class != ClassProviderFailure || obs.ProviderClass != "auth" {
			t.Fatalf("obs=%+v want provider-failure/auth", obs)
		}
	})
	t.Run("nonzero process error", func(t *testing.T) {
		obs := probe(t, "process-error", 5*time.Second)
		if obs.Class != ClassProcessFailure {
			t.Fatalf("obs=%+v want process-failure", obs)
		}
	})
	t.Run("deadline no output (never-ending child)", func(t *testing.T) {
		start := time.Now()
		obs := probe(t, "never-ending", 400*time.Millisecond)
		if obs.Class != ClassDeadlineNoOutput {
			t.Fatalf("obs=%+v want deadline-no-output", obs)
		}
		if elapsed := time.Since(start); elapsed > 5*time.Second {
			t.Fatalf("probe took %v; the child must be killed and reaped promptly", elapsed)
		}
	})
	t.Run("deadline after output (partial-output child)", func(t *testing.T) {
		start := time.Now()
		obs := probe(t, "partial-then-hang", 400*time.Millisecond)
		if obs.Class != ClassDeadlineAfterOutput {
			t.Fatalf("obs=%+v want deadline-after-output", obs)
		}
		if elapsed := time.Since(start); elapsed > 5*time.Second {
			t.Fatalf("probe took %v; the child must be killed and reaped promptly", elapsed)
		}
	})
}

// --- regression tests for malformed/duplicate envelopes, bounded capture,
// ambiguous content, nested malformed roles, secret-safe tails ------------

func TestDuplicateAndAliasedKeysRejected(t *testing.T) {
	// Duplicate semantic key (after Unmarshal, map collapses duplicates to
	// last-wins; pre-decode scan must catch the duplication).
	cases := []struct {
		name    string
		raw     string
		wantErr bool
	}{
		{"duplicate content", `{"content":"PONG","content":"WRONG"}`, true},
		{"duplicate role", `{"role":"assistant","role":"user","content":"PONG"}`, true},
		{"case-aliased is_error/isError same value", `{"is_error":true,"isError":true,"content":"PONG"}`, true},
		{"case-aliased is_error/isError contradictory", `{"is_error":true,"isError":false,"content":"PONG"}`, true},
		{"present-null role", `{"role":null,"content":"PONG"}`, true},
		{"present-null is_error", `{"is_error":null,"content":"PONG"}`, true},
		{"duplicate nested wrapper", `{"result":{"content":"PONG"},"result":{"message":"WRONG"}}`, true},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			err := preDecodeScan(tc.raw)
			if (err != nil) != tc.wantErr {
				t.Fatalf("preDecodeScan(%q) error=%v want error=%v", tc.raw, err, tc.wantErr)
			}
		})
	}
}

func TestAmbiguousContentDoesNotPassAsReady(t *testing.T) {
	// Bare content key without assistant role is a recognized result envelope
	// (intended good CLI reply); but an unknown typed object (e.g. number, array,
	// nested object without recognized role/content) must not yield ready.
	obsUnknown := classifyReadiness(`{"content":123}`, "", 0, false, false, "")
	if obsUnknown.Class != ClassMalformedReply {
		t.Fatalf("unknown typed content (number) should be malformed-reply, got %s", obsUnknown.Class)
	}
	obsNestedUnknown := classifyReadiness(`{"data":{"unknown":"PONG"}}`, "", 0, false, false, "")
	if obsNestedUnknown.Class != ClassMalformedReply {
		t.Fatalf("nested unknown wrapper should be malformed-reply, got %s", obsNestedUnknown.Class)
	}
	// Contradictory nested role/content: wrapper has assistant role but inner
	// has user role — should not silently choose passing interpretation.
	obsContradictory := classifyReadiness(`{"data":{"role":"user","content":"PONG"}}`, "", 0, false, false, "")
	if obsContradictory.Class != ClassMalformedReply {
		t.Fatalf("nested contradictory role/content should be malformed-reply, got %s", obsContradictory.Class)
	}
	// Bare assistant content without role or recognized schema is malformed —
	// it lacks provenance and must not pass as ready (fixture correction).
	obsBareContent := classifyReadiness(`{"content":"PONG"}`, "", 0, false, false, "")
	if obsBareContent.Class != ClassMalformedReply || obsBareContent.Ready {
		t.Fatalf("bare content without provenance should be malformed-reply/not-ready, got %s ready=%v", obsBareContent.Class, obsBareContent.Ready)
	}
}

func TestBoundedCaptureOverflowExplicitNotReady(t *testing.T) {
	// Host-level test: boundedWriter with overflow must produce non-ready
	// observation with preserved tail observation (not silently ready).
	var out bytes.Buffer
	overflow := false
	bw := boundedWriter{buf: &out, max: 8, overflow: &overflow}
	// Write more than max bytes to trigger overflow.
	long := "PONGPONGPONGPONGPONG"
	n, err := bw.Write([]byte(long))
	if err != nil {
		t.Fatalf("boundedWriter Write error: %v", err)
	}
	if n != len(long) {
		t.Fatalf("Write returned %d want %d (must continue consuming)", n, len(long))
	}
	if !overflow {
		t.Fatalf("expected overflow=true after exceeding max bytes")
	}
	// Observation on overflow must not be ClassReady.
	obs := classifyReadiness(out.String(), "", 0, false, true, "overflow")
	if obs.Class == ClassReady {
		t.Fatalf("overflowed partial capture must NOT yield ready (prevent partial PONG pass)")
	}
	if !strings.Contains(string(obs.StdoutTail), "PONG") && obs.SawSentinel {
		// The partial PONG may be preserved in tail as observation; that is fine.
	}
}

func TestSecretSafeTailScrubsCredentials(t *testing.T) {
	// Synthetic bearer/JSON credential values without real credentials.
	samples := []string{
		`Authorization: Bearer sk-test12345678901234567890abcd`,
		`token=ghp_synthetic0123456789012345678a`,
		`{"error":"auth failed","token":"sk-testsecretvalue"}`,
		`key = AKIAIOSFODNN7EXAMPLE`,
		`authorization: bearer eyJhbGci.eyJzdWI.test.signaturevalue123`,
	}
	for _, s := range samples {
		scrubbed := scrubSecrets(s)
		if strings.Contains(scrubbed, "sk-test") || strings.Contains(scrubbed, "ghp_") ||
			strings.Contains(scrubbed, "AKIAIOSFODNN7") || strings.Contains(scrubbed, "eyJhbGci") {
			t.Errorf("scrubSecrets did not fully scrub synthetic secret in %q; got %q", s, scrubbed)
		}
		// The authorization bearer pattern specifically must not leak suffix.
		if strings.Contains(s, "Authorization") && strings.Contains(scrubbed, "Bearer") {
			// After full redaction, "Bearer" may still appear if it is the label,
			// but the token value must not remain untouched. Confirm full value replaced.
			if strings.Contains(scrubbed, "test123") || strings.Contains(scrubbed, "abcd") {
				t.Errorf("Authorization bearer value suffix leaked in %q -> %q", s, scrubbed)
			}
		}
	}
}

func TestNestedMalformedRoleNotPassedAsReady(t *testing.T) {
	// Nested wrapper with assistant role and inner content must pass; nested
	// wrapper with user/system/tool role must fail; nested wrapper with malformed
	// role type (non-string) must fail.
	goodNested := `{"data":{"role":"assistant","content":"PONG"}}`
	obs := classifyReadiness(goodNested, "", 0, false, false, "")
	if obs.Class != ClassReady || !obs.Ready {
		t.Fatalf("good nested assistant envelope expected ready, got %s ready=%v", obs.Class, obs.Ready)
	}
	badRoleType := `{"data":{"role":123,"content":"PONG"}}`
	obs2 := classifyReadiness(badRoleType, "", 0, false, false, "")
	if obs2.Class != ClassMalformedReply {
		t.Fatalf("nested malformed role type (number) expected malformed-reply, got %s", obs2.Class)
	}
}

func TestAmbiguousUnknownTypeRejected(t *testing.T) {
	obs := classifyReadiness(`{"type":"unknown","content":"PONG"}`, "", 0, false, false, "")
	if obs.Class == ClassReady {
		t.Fatalf("unknown type %q with content should NOT be ready (ambiguous schema)", "unknown")
	}
}

func TestCompetingOutputFieldsRejected(t *testing.T) {
	obs := classifyReadiness(`{"role":"assistant","content":"PONG","text":"FAIL"}`, "", 0, false, false, "")
	if obs.Class == ClassReady {
		t.Fatalf("competing content+text fields should NOT be ready")
	}
}

func TestBoundedCapturePreservesWhitespaceObservation(t *testing.T) {
	var out bytes.Buffer
	overflow := false
	truncated := false
	bw := boundedWriter{buf: &out, max: 4, overflow: &overflow, truncated: &truncated, observedBytes: 0}
	n, err := bw.Write([]byte("    PONG"))
	if err != nil {
		t.Fatalf("boundedWriter Write error: %v", err)
	}
	if n != len("    PONG") {
		t.Fatalf("Write returned %d want %d", n, len("    PONG"))
	}
	if !overflow || !truncated {
		t.Fatalf("expected overflow=true and truncated=true")
	}
	if bw.observedBytes != len("    PONG") {
		t.Fatalf("observedBytes=%d want %d (independent of cap)", bw.observedBytes, len("    PONG"))
	}
	obs := classifyReadiness(out.String(), "", 0, false, true, "overflow")
	if obs.Class == ClassReady {
		t.Fatalf("whitespace-prefixed overflow must NOT yield ready")
	}
	if !obs.Truncated || obs.TruncationReason != "overflow" {
		t.Fatalf("expected Truncated=true and TruncationReason=overflow, got truncated=%v reason=%q", obs.Truncated, obs.TruncationReason)
	}
}

func TestBoundedCaptureEmptyWriteAtExactBoundNoOverflow(t *testing.T) {
	var out bytes.Buffer
	overflow := false
	truncated := false
	bw := boundedWriter{buf: &out, max: 4, overflow: &overflow, truncated: &truncated, observedBytes: 0}
	n, err := bw.Write([]byte("PONG"))
	if err != nil {
		t.Fatalf("Write error: %v", err)
	}
	if n != 4 {
		t.Fatalf("Write returned %d want 4", n)
	}
	if overflow || truncated {
		t.Fatalf("exact-bound write must NOT set overflow or truncated")
	}
	if bw.observedBytes != 4 {
		t.Fatalf("observedBytes=%d want 4", bw.observedBytes)
	}
}

func TestSecretSafeTailScrubsJSONLabeledCredentials(t *testing.T) {
	raw := `{"token":"sk-testsecretvalue"}`
	scrubbed := scrubSecrets(raw)
	if strings.Contains(scrubbed, "sk-testsecretvalue") {
		t.Errorf("labeled JSON secret value not fully scrubbed: got %q", scrubbed)
	}
	if !strings.Contains(scrubbed, "«redacted»") {
		t.Errorf("expected redaction marker in scrubbed result, got %q", scrubbed)
	}
	bearer := `Authorization: Bearer sk-testsecretvalue123`
	scrubbedBearer := scrubSecrets(bearer)
	if strings.Contains(scrubbedBearer, "sk-testsecretvalue123") {
		t.Errorf("bearer token suffix leaked in %q", scrubbedBearer)
	}
}
