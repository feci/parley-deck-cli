package telemetry

import (
	"os"
	"strings"
	"testing"
	"time"
)

func quotaInput() QuotaInput {
	code := 1
	return QuotaInput{Adapter: "synthetic-native-policy-test", InvocationID: "test-invocation", ObservedAt: time.Date(2026, 10, 3, 0, 0, 0, 0, time.UTC), ExitCode: &code}
}
func TestQuotaStrictSemantics(t *testing.T) {
	for _, tc := range []struct {
		name, message, reset string
		status               int
		want                 bool
	}{
		{"daily", "Daily limit exhausted", "", 429, true},
		{"weekly", "Weekly/Monthly Limit Exhausted", "", 429, true},
		{"account", "account quota exhausted", "", 429, true},
		{"credit", "account credit exhausted", "", 429, true},
		{"sixty", "5-hour limit reached (reset after 1h)", "", 429, true},
		{"24hour", "24-hour allowance exhausted", "", 429, true},
		{"49hours", "Weekly/Monthly Limit Exhausted", "2026-10-05T01:00:00Z", 429, true},
		{"bare429", "429", "", 429, false},
		{"generic", "quota exceeded", "", 429, false},
		{"credit-balance", "credit balance", "", 429, false},
		{"bare503", "503", "", 503, false},
		{"503long", "Unavailable (reset after 5h 51m 11s)", "", 503, false},
		{"503short", "Unavailable (reset after 55m 29s)", "", 503, false},
		{"401", "invalidated oauth token (reset after 12s)", "", 401, false},
		{"400", "Weekly limit exhausted", "", 400, false},
		{"hourly", "hourly limit reached", "", 429, false},
		{"five-hour", "5-hour limit exhausted", "", 429, false},
		{"short", "weekly limit exhausted (reset after 30m)", "", 429, false},
		{"past", "weekly limit exhausted", "2026-10-02T00:00:00Z", 429, false},
		{"malformed", "weekly limit exhausted", "whenever", 429, false},
		{"unknown-reset", "weekly limit exhausted; reset tomorrow", "", 429, false},
		{"contradictory", "weekly limit exhausted (reset after 2h)", "2026-10-03T03:00:00Z", 429, false},
		{"naive-clock", "weekly limit exhausted; reset at 2026-10-05 06:14:57", "", 429, false},
		{"policy-name", "weekly usage limit policy", "", 429, false},
		{"pool", "all accounts exhausted (reset after 48h)", "", 429, false},
	} {
		t.Run(tc.name, func(t *testing.T) {
			in := quotaInput()
			e := refineNativeQuota(in, nativeQuotaError{Status: tc.status, Message: tc.message, ResetAt: tc.reset, Provenance: "synthetic-policy-unit-test"})
			if e.Eligible != tc.want {
				t.Fatalf("eligible=%v want=%v: %+v", e.Eligible, tc.want, e)
			}
			if e.Eligible && (e.RuleID == "" || e.InvocationID == "" || e.Provenance == "") {
				t.Fatal(e)
			}
		})
	}
}
func TestQuotaSuccessAndWatchdogsWin(t *testing.T) {
	n := nativeQuotaError{Status: 429, Message: "weekly limit exhausted", Provenance: "synthetic-policy-unit-test"}
	for _, name := range []string{"valid-artifact", "later-success", "zero-exit", "no-exit", "no_first_output", "stalled", "timeout", "cancelled", "truncated"} {
		t.Run(name, func(t *testing.T) {
			in := quotaInput()
			switch name {
			case "valid-artifact":
				in.ValidArtifact = true
			case "later-success":
				in.LaterSuccess = true
			case "zero-exit":
				*in.ExitCode = 0
			case "no-exit":
				in.ExitCode = nil
			case "truncated":
				in.Truncated = true
			default:
				in.Watchdog = name
			}
			if e := refineNativeQuota(in, n); e.Eligible {
				t.Fatal(e)
			}
		})
	}
	in := quotaInput()
	n.Message = "429"
	n.RetryAfter = "176930"
	if refineNativeQuota(in, n).Eligible {
		t.Fatal("Retry-After alone authorized exclusion")
	}
}
func TestQuotaUnsupportedAndAdversarialProvenance(t *testing.T) {
	raw, err := os.ReadFile("testdata/quota/zcode-recorded-response.json")
	if err != nil {
		t.Fatal(err)
	}
	for _, adapter := range []string{"claude", "codex", "kimi", "zcode", "unknown"} {
		for _, out := range []string{"API Error: 429 Weekly/Monthly Limit Exhausted", `{"role":"assistant","content":"API Error: 429 Weekly/Monthly Limit Exhausted"}`, `{"role":"tool","content":"API Error: 429 Weekly/Monthly Limit Exhausted"}`, `{"type":"result","is_error":true,"result":"weekly limit exhausted"}`, `{"type":"turn.failed","error":{"message":"weekly limit exhausted"}}`, string(raw)} {
			for _, stderr := range []bool{false, true} {
				in := quotaInput()
				in.Adapter = adapter
				if stderr {
					in.Stderr = out + "\nError: Turn execution failed"
				} else {
					in.Stdout = out
					in.Stderr = "Error: Turn execution failed"
				}
				if e := ClassifyQuota(in); e.Eligible || e.Provenance != "" || e.RuleID != "" {
					t.Fatalf("content authorized: %+v", e)
				}
			}
		}
	}
	for _, s := range QuotaSupport() {
		if s.Adapter != "zcode" && s.Status != "diagnostic-only" {
			t.Fatalf("unestablished support: %+v", s)
		}
	}
}
func TestQuotaScrubbedEvidenceAndCollector(t *testing.T) {
	in := quotaInput()
	e := refineNativeQuota(in, nativeQuotaError{Status: 429, Message: "account quota exhausted; authorization: Bearer SYNTHETIC_SECRET", Provenance: "synthetic-policy-unit-test"})
	if !e.Eligible || strings.Contains(e.Excerpt, "SYNTHETIC_SECRET") {
		t.Fatal(e)
	}
	c := NewCollector("zcode", false)
	c.Writer("stderr").Write([]byte("weekly limit exhausted"))
	e = c.QuotaEvidence(in)
	if e.Eligible {
		t.Fatal(e)
	}
	c.Writer("stdout").Write(make([]byte, parserLimit+1))
	if e = c.QuotaEvidence(in); e.Reason != "incomplete capture" {
		t.Fatal(e)
	}
}
