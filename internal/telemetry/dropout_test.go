package telemetry

import (
	"strings"
	"testing"
	"time"
)

func TestDropoutSupervisorOutcomes(t *testing.T) {
	now := time.Now().UTC()
	one, zero := 1, 0
	for _, class := range []string{"process_failure", "start_failure", "no_first_output", "stalled", "timeout", "crash", "provider-error", "auth-error", "rate-limit"} {
		r := Record{Type: "invocation.terminal", InvocationID: "observed", CompletedAt: &now, Outcome: &Outcome{Status: "failed", ExitCode: &one, FailureClass: String(class), DispatchAttempted: true}}
		a, err := ParticipantAttempt(r, false, "own output missing")
		if err != nil || a == nil {
			t.Fatalf("%s: %+v %v", class, a, err)
		}
		if a, err = ParticipantAttempt(r, true, ""); err != nil || a != nil {
			t.Fatalf("valid BLOCK lost for %s: %+v %v", class, a, err)
		}
	}
	for _, class := range []string{"cancelled", "budget_refused", "protocol_context_refused", "telemetry_failure", "trajectory_failure", "future-unknown"} {
		r := Record{Type: "invocation.terminal", InvocationID: "refused", CompletedAt: &now, Outcome: &Outcome{Status: "failed", FailureClass: String(class), DispatchAttempted: true}}
		for _, valid := range []bool{false, true} {
			if a, err := ParticipantAttempt(r, valid, "own missing"); err == nil || a != nil {
				t.Fatalf("%s authorized dropout: %+v %v", class, a, err)
			}
		}
	}
	r := Record{Type: "invocation.terminal", InvocationID: "empty", CompletedAt: &now, Outcome: &Outcome{Status: "process-exited", ExitCode: &zero, DispatchAttempted: true}}
	if a, err := ParticipantAttempt(r, false, "missing Summary heading"); err != nil || a == nil || a.ValidatorReason != "missing Summary heading" {
		t.Fatalf("%+v %v", a, err)
	}
	r.Outcome.DispatchAttempted = false
	if _, err := ParticipantAttempt(r, false, "missing"); err == nil {
		t.Fatal("undispatched setup uncertainty counted")
	}
}

func TestDropoutDiagnosticRedactionAndVisibleTruncation(t *testing.T) {
	input := "invalid api_key=supersecret authorization: Bearer private-token sk-123456789012345678901234 "
	got := participantDiagnostic(input + strings.Repeat("x", 5000))
	for _, secret := range []string{"supersecret", "private-token", "sk-1234"} {
		if strings.Contains(got, secret) {
			t.Fatal("secret leaked")
		}
	}
	if !strings.Contains(got, "[TRUNCATED;") || !strings.Contains(got, "invalid api_key=[REDACTED]") {
		t.Fatal(got)
	}
}
