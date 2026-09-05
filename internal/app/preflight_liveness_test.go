package app

import (
	"bytes"
	"context"
	"os"
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
		{"recognized JSON content key", `{"content":"PONG"}`, "", 0, false, ClassReady, true, false},
		{"recognized JSON message key", `{"message":"PONG"}`, "", 0, false, ClassReady, true, false},
		{"recognized JSON text key", `{"text":"PONG"}`, "", 0, false, ClassReady, true, false},
		{"recognized JSON result string", `{"result":"PONG"}`, "", 0, false, ClassReady, true, false},
		{"recognized assistant-role message", `{"role":"assistant","content":"PONG"}`, "", 0, false, ClassReady, true, false},
		{"recognized result wrapper", `{"result":{"content":"PONG"}}`, "", 0, false, ClassReady, true, false},
		{"recognized data wrapper", `{"data":{"content":"PONG"}}`, "", 0, false, ClassReady, true, false},
		{"recognized payload wrapper", `{"payload":{"content":"PONG"}}`, "", 0, false, ClassReady, true, false},
		{"recognized response wrapper", `{"response":{"content":"PONG"}}`, "", 0, false, ClassReady, true, false},
		{"recognized result wrapper message key", `{"result":{"message":"PONG"}}`, "", 0, false, ClassReady, true, false},
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
			obs := classifyReadiness(tc.stdout, tc.stderr, tc.exitCode, tc.timedOut)
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
	obs := classifyReadiness("", "429 Too Many Requests, rate limit hit", 1, false)
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
		_, _ = os.Stdout.WriteString(`{"content":"PONG"}` + "\n")
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
