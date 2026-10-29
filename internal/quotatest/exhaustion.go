package quotatest

import (
	"strings"
	"time"

	"parley-deck-cli/internal/quota"
)

// Exhaustion supplies synthetic bound native evidence without a provider call.
func Exhaustion(idea, agent, step string) quota.Evidence {
	now := time.Date(2026, 10, 29, 20, 0, 0, 0, time.UTC)
	exit := 1
	message := "Your workspace is out of credits. Add credits to continue."
	return quota.Evidence{InvocationID: agent + "-inv", Adapter: "codex", RuleID: quota.ProviderExhaustionRule,
		Provenance: quota.CodexExhaustionProvenance, Eligible: true, Excerpt: message, ObservedAt: now,
		Exhaustion: &quota.ExhaustedStep{Idea: idea, Agent: agent, Step: step, CLIVersion: "0.165.0", ExplicitJSON: true,
			TerminalSHA256: strings.Repeat("a", 64), ValidationSHA256: strings.Repeat("b", 64),
			Attempt: quota.FailedAttempt{InvocationID: agent + "-inv", Status: "failed", ExitCode: &exit, ObservedAt: now, Excerpt: message}}}
}
