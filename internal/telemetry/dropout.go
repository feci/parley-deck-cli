package telemetry

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strings"

	"parley-deck-cli/internal/quota"
)

// ParticipantAttempt uses settled supervisor facts plus the caller's structural
// validation. It never searches provider/model/tool prose for dropout triggers.
// A nil attempt means a valid result; an error is a control/integrity stop.
func ParticipantAttempt(r Record, valid bool, validatorReason string) (*quota.FailedAttempt, error) {
	if (r.Type != "invocation.terminal" && r.Type != "invocation.crash-settled") || r.InvocationID == "" || r.CompletedAt == nil || r.Outcome == nil {
		return nil, fmt.Errorf("participant attempt has no settled terminal evidence")
	}
	o := r.Outcome
	class := ""
	if o.FailureClass != nil {
		class = *o.FailureClass
	}
	if class != "" && !quota.ParticipantFailureClass(class) {
		return nil, fmt.Errorf("control-plane participant outcome: %s", class)
	}
	if !o.DispatchAttempted && r.StartedAt == nil || (o.Status != "failed" && o.Status != "process-exited") {
		return nil, fmt.Errorf("participant child dispatch was not observed")
	}
	if valid {
		return nil, nil
	}
	if class == "" && (o.ExitCode == nil || *o.ExitCode == 0) && validatorReason == "" {
		return nil, fmt.Errorf("no participant failure or structural validation result")
	}
	a := &quota.FailedAttempt{InvocationID: r.InvocationID, Status: o.Status, FailureClass: class,
		ExitCode: o.ExitCode, ObservedAt: *r.CompletedAt}
	if validatorReason != "" {
		a.ValidatorReason = participantDiagnostic(validatorReason)
	}
	if r.Metadata.RetryOf != nil {
		a.RetryOf = *r.Metadata.RetryOf
	}
	// Quote the decisive, content-free terminal fields, not an unsafe transcript
	// tail. These are the actual recorded supervisor values, not provider guesses.
	raw, _ := json.Marshal(struct {
		Status       string `json:"status"`
		ExitCode     *int   `json:"exit_code"`
		FailureClass string `json:"failure_class"`
		Validator    string `json:"validator_reason,omitempty"`
	}{a.Status, a.ExitCode, a.FailureClass, a.ValidatorReason})
	a.Excerpt = string(raw)
	return a, nil
}

func ParticipantEvidence(idea, agent, step, adapter string, attempts []quota.FailedAttempt) (quota.Evidence, error) {
	if len(attempts) != 2 {
		return quota.Evidence{}, fmt.Errorf("participant dropout requires exactly two settled failures")
	}
	last := attempts[1]
	e := quota.Evidence{InvocationID: last.InvocationID, Adapter: adapter, RuleID: quota.ParticipantFailureRule,
		Provenance: "supervisor-terminal", Eligible: true, Reason: "two settled child failures",
		Excerpt: attempts[0].Excerpt + "\n" + last.Excerpt, ObservedAt: last.ObservedAt,
		Failure: &quota.FailedStep{Idea: idea, Agent: agent, Step: step, Attempts: append([]quota.FailedAttempt(nil), attempts...)}}
	return e, quota.ValidateFailureEvidence(e, idea, agent)
}

// Keep the actual validator reason, redacting credential shapes before canonical
// evidence. Raw diagnostics remain only in the private per-invocation receipt.
var participantSecrets = []*regexp.Regexp{
	regexp.MustCompile(`(?i)(authorization\s*[:=]\s*)(bearer\s+)?\S+`),
	regexp.MustCompile(`(?i)(token|secret|password|passwd|api[_-]?key|access[_-]?key|private[_-]?key)(\s*[:=]\s*)\S+`),
	regexp.MustCompile(`(?i)\bbearer\s+[A-Za-z0-9._\-]+`),
	regexp.MustCompile(`\bsk-[A-Za-z0-9_\-]{16,}|\bgh[pousr]_[A-Za-z0-9]{20,}|\bxox[baprs]-[A-Za-z0-9\-]{10,}|\bAKIA[0-9A-Z]{16}\b|\beyJ[A-Za-z0-9_\-]{10,}\.[A-Za-z0-9_\-]{10,}\.[A-Za-z0-9_\-]{5,}`),
}

func participantDiagnostic(s string) string {
	s = participantSecrets[0].ReplaceAllString(s, "${1}[REDACTED]")
	s = participantSecrets[1].ReplaceAllString(s, "${1}${2}[REDACTED]")
	for _, re := range participantSecrets[2:] {
		s = re.ReplaceAllString(s, "[REDACTED]")
	}
	if len(s) > 4096 {
		s = s[:4096] + " [TRUNCATED; full diagnostic retained privately]"
	}
	return strings.ToValidUTF8(s, "�")
}
