package quota

import (
	"fmt"
	"regexp"
)

// ExhaustedStep binds native evidence to one failed, validated child invocation.
// The terminal and validation bytes remain in the existing invocation ledger.
type ExhaustedStep struct {
	Idea             string        `json:"idea"`
	Agent            string        `json:"agent"`
	Step             string        `json:"step"`
	Attempt          FailedAttempt `json:"attempt"`
	TerminalSHA256   string        `json:"terminal_sha256"`
	ValidationSHA256 string        `json:"validation_sha256"`
	CLIVersion       string        `json:"cli_version,omitempty"`
	ExplicitJSON     bool          `json:"explicit_json,omitempty"`
}

const CodexExhaustionProvenance = "codex.exec-json.turn-failed.v1"
const ZcodeExhaustionProvenance = "zcode.stderr.owner-deviation.2026-10-04"

var proofHash = regexp.MustCompile(`^[0-9a-f]{64}$`)

func ValidateExhaustionEvidence(e Evidence, idea, agent string) error {
	f := e.Exhaustion
	if e.RuleID != ProviderExhaustionRule || !e.Eligible || e.Failure != nil || f == nil ||
		f.Idea == "" || f.Agent != agent || f.Step == "" || (idea != "" && f.Idea != idea) ||
		e.InvocationID == "" || e.InvocationID != f.Attempt.InvocationID || e.ObservedAt.IsZero() ||
		e.Excerpt == "" || !proofHash.MatchString(f.TerminalSHA256) || !proofHash.MatchString(f.ValidationSHA256) {
		return fmt.Errorf("invalid provider exhaustion proof")
	}
	a := f.Attempt
	if a.ObservedAt.IsZero() || e.ObservedAt.After(a.ObservedAt) || a.Excerpt == "" ||
		(a.Status != "failed" && a.Status != "process-exited") ||
		(a.FailureClass != "" && !ParticipantFailureClass(a.FailureClass)) ||
		(a.FailureClass == "" && (a.ExitCode == nil || *a.ExitCode == 0)) {
		return fmt.Errorf("exhaustion proof lacks failed execution")
	}
	if e.ResetAt != nil && (e.RawReset == "" || e.ResetAt.Sub(e.ObservedAt) < MinimumReset) {
		return fmt.Errorf("invalid exhaustion reset hint")
	}
	switch e.Adapter {
	case "codex":
		if e.Provenance != CodexExhaustionProvenance || f.CLIVersion != "0.165.0" || !f.ExplicitJSON || !CodexExhaustionMessage(e.Excerpt) || e.ResetAt != nil || e.RawReset != "" {
			return fmt.Errorf("unsupported Codex exhaustion contract")
		}
	case "zcode":
		if e.Provenance != ZcodeExhaustionProvenance || a.ExitCode == nil || *a.ExitCode == 0 {
			return fmt.Errorf("unsupported Zcode exhaustion contract")
		}
	default:
		return fmt.Errorf("unsupported provider exhaustion adapter")
	}
	return nil
}

func CodexExhaustionMessage(message string) bool {
	switch message {
	case "Your workspace is out of credits. Add credits to continue.",
		"Your workspace is out of credits. Ask your workspace owner to refill in order to continue.",
		"You hit your spend cap set in your workspace. Increase your spend cap to continue.",
		"You hit your spend cap set by the owner of your workspace. Ask an owner to increase your spend cap to continue.":
		return true
	}
	return false
}

func ValidatePolicyEvidence(p Policy, e Evidence, idea, agent string) error {
	if e.RuleID == ProviderExhaustionRule {
		if !p.Exhaustion() {
			return fmt.Errorf("exhaustion evidence contradicts saved trigger")
		}
		return ValidateExhaustionEvidence(e, idea, agent)
	}
	if p.Dropout() || e.RuleID == ParticipantFailureRule {
		if !p.Dropout() {
			return fmt.Errorf("participant failure evidence contradicts saved trigger")
		}
		return ValidateFailureEvidence(e, idea, agent)
	}
	return nil
}
