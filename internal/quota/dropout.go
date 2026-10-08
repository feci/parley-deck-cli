package quota

import (
	"fmt"
	"time"
)

// FailedStep is the paired, supervisor-derived proof bound into the existing
// immutable membership transition. Private transcripts are never embedded.
type FailedStep struct {
	Idea     string          `json:"idea"`
	Agent    string          `json:"agent"`
	Step     string          `json:"step"`
	Attempts []FailedAttempt `json:"attempts"`
}

type FailedAttempt struct {
	InvocationID    string    `json:"invocation_id"`
	RetryOf         string    `json:"retry_of,omitempty"`
	Status          string    `json:"status"`
	FailureClass    string    `json:"failure_class,omitempty"`
	ExitCode        *int      `json:"exit_code"`
	ValidatorReason string    `json:"validator_reason,omitempty"`
	ObservedAt      time.Time `json:"observed_at"`
	Excerpt         string    `json:"excerpt"`
}

// ParticipantFailureClass admits only child outcomes. Policy, operator,
// transport-accounting and telemetry/integrity refusals are deliberately absent.
func ParticipantFailureClass(class string) bool {
	switch class {
	case "process_failure", "start_failure", "no_first_output", "stalled", "timeout",
		"auth", "billing", "quota", "rate_limit", "overload", "model_not_found", "provider_failure", "crash":
		return true
	case "provider-error", "auth-error", "rate-limit":
		return true
	}
	return false
}

func ValidateFailureEvidence(e Evidence, idea, agent string) error {
	f := e.Failure
	if e.RuleID != ParticipantFailureRule || e.Provenance != "supervisor-terminal" || !e.Eligible || e.ResetAt != nil || e.RawReset != "" || f == nil || f.Idea == "" || f.Agent != agent || f.Step == "" || (idea != "" && f.Idea != idea) || len(f.Attempts) != 2 {
		return fmt.Errorf("invalid participant failure proof")
	}
	for i, a := range f.Attempts {
		if a.InvocationID == "" || a.ObservedAt.IsZero() || a.Excerpt == "" || (a.Status != "failed" && a.Status != "process-exited") {
			return fmt.Errorf("invalid failed attempt identity/outcome")
		}
		if a.FailureClass != "" && !ParticipantFailureClass(a.FailureClass) || a.FailureClass == "" && (a.ExitCode == nil || *a.ExitCode == 0) && a.ValidatorReason == "" {
			return fmt.Errorf("attempt has no eligible child failure")
		}
		if i == 0 && a.RetryOf != "" || i == 1 && (a.RetryOf != f.Attempts[0].InvocationID || a.InvocationID == a.RetryOf || a.ObservedAt.Before(f.Attempts[0].ObservedAt)) {
			return fmt.Errorf("invalid participant retry linkage")
		}
	}
	last := f.Attempts[1]
	if e.InvocationID != last.InvocationID || !e.ObservedAt.Equal(last.ObservedAt) || e.Excerpt == "" {
		return fmt.Errorf("contradictory participant failure evidence")
	}
	return nil
}

// Dropped is history-derived and independent of the current policy, including
// a later owner-authorized opt-out/downgrade. Kickoff candidates never became
// known quorum members but are just as permanently ineligible in this idea.
func (h *History) Dropped(id string) bool {
	if h == nil || h.Kickoff == nil {
		return false
	}
	has := func(cs []Candidate) bool {
		for _, c := range cs {
			if c.Agent == id && c.Evidence.RuleID == ParticipantFailureRule {
				return true
			}
		}
		return false
	}
	if h.Kickoff.Transition != nil && has(h.Kickoff.Transition.Decision.Candidates) {
		return true
	}
	for _, b := range h.Batches {
		if b.Owner == nil && has(b.Decision.Candidates) {
			return true
		}
	}
	return false
}

func (h *History) CheckReturn(ids []string) error {
	for _, id := range ids {
		if h.Dropped(id) {
			return fmt.Errorf("participant %s permanently dropped from this idea; re-probe only in a new idea", id)
		}
	}
	return nil
}

const OwnerOptions = "Owner options: authorize a separate eligible model-diverse reviewer process (recommended); record an attended evidence-backed continuation with reduced reviewer count; or pause/abandon. No substitute or gate waiver is selected automatically."
