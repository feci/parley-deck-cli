package telemetry

import (
	"encoding/json"
	"strings"
	"testing"

	"parley-deck-cli/internal/quota"
)

func codexFailure(message string) string {
	b, _ := json.Marshal(map[string]any{"type": "turn.failed", "error": map[string]string{"message": message}})
	return "{\"type\":\"thread.started\",\"thread_id\":\"primary\"}\n{\"type\":\"turn.started\"}\n" + string(b) + "\n"
}

func TestExhaustionCodexExactSourceContract(t *testing.T) {
	in := quotaInput()
	in.Adapter, in.CLIVersion, in.ExplicitJSON = "codex", "0.165.0", true
	messages := []string{
		"Your workspace is out of credits. Add credits to continue.",
		"Your workspace is out of credits. Ask your workspace owner to refill in order to continue.",
		"You hit your spend cap set in your workspace. Increase your spend cap to continue.",
		"You hit your spend cap set by the owner of your workspace. Ask an owner to increase your spend cap to continue.",
	}
	for _, message := range messages {
		in.Stdout = codexFailure(message)
		e := ClassifyExhaustion(in)
		if !e.Eligible || e.Excerpt != message || e.RuleID != quota.ProviderExhaustionRule || e.Provenance != quota.CodexExhaustionProvenance {
			t.Fatal(e)
		}
		if ClassifyQuota(in).Eligible {
			t.Fatal("legacy recognition widened")
		}
	}
	for _, message := range []string{
		messages[0] + " (extra detail)", "Quota exceeded. Check your plan and billing details.",
		"You've hit your usage limit.", "You’ve hit your usage limit.", "turn failed",
		"To use Codex with your ChatGPT plan, upgrade to Plus: https://chatgpt.com/explore/plus.",
		"shared rollout token budget exhausted",
		"Codex ran out of room in the model's context window. Start a new thread or clear earlier history before retrying.",
	} {
		in.Stdout = codexFailure(message)
		if e := ClassifyExhaustion(in); e.Eligible {
			t.Fatal(message, e)
		}
	}
	in.Stdout = codexFailure(messages[0])
	for name, change := range map[string]func(*QuotaInput){
		"text":            func(i *QuotaInput) { i.ExplicitJSON = false },
		"unknown-version": func(i *QuotaInput) { i.CLIVersion = "" },
		"future-version":  func(i *QuotaInput) { i.CLIVersion = "0.166.0" },
		"truncated":       func(i *QuotaInput) { i.Truncated = true },
		"watchdog":        func(i *QuotaInput) { i.Watchdog = "timeout" },
		"valid-block":     func(i *QuotaInput) { i.ValidArtifact = true },
		"later-success":   func(i *QuotaInput) { i.LaterSuccess = true },
		"mixed-terminal":  func(i *QuotaInput) { i.Stdout += "{\"type\":\"turn.completed\",\"usage\":{}}\n" },
		"second-thread":   func(i *QuotaInput) { i.Stdout = "{\"type\":\"thread.started\",\"thread_id\":\"other\"}\n" + i.Stdout },
		"missing-turn":    func(i *QuotaInput) { i.Stdout = strings.Replace(i.Stdout, "{\"type\":\"turn.started\"}\n", "", 1) },
		"duplicate-key": func(i *QuotaInput) {
			i.Stdout = strings.Replace(i.Stdout, `"turn.failed"`, `"turn.failed","type":"turn.failed"`, 1)
		},
		"assistant-quote": func(i *QuotaInput) {
			b, _ := json.Marshal(map[string]any{"type": "item.completed", "item": map[string]string{"type": "agent_message", "text": i.Stdout}})
			i.Stdout = string(b)
		},
		"tool-quote": func(i *QuotaInput) {
			b, _ := json.Marshal(map[string]any{"type": "item.completed", "item": map[string]string{"type": "command_execution", "aggregated_output": i.Stdout}})
			i.Stdout = string(b)
		},
	} {
		t.Run(name, func(t *testing.T) {
			bad := in
			change(&bad)
			if e := ClassifyExhaustion(bad); e.Eligible {
				t.Fatal(e)
			}
		})
	}
}

func TestExhaustionExplicitEffectiveArgs(t *testing.T) {
	for _, args := range [][]string{{"exec", "--json", "-"}, {"-c", "x=y", "exec", "-m", "model", "--json", "prompt"}} {
		if !ExplicitCodexJSON(args) {
			t.Fatal(args)
		}
	}
	for _, args := range [][]string{{"exec", "prompt --json"}, {"exec", "-m", "--json", "prompt"}, {"exec", "--", "--json"}, {"exec", "-", "--json"}, {"--json", "exec"}, {"exec", "--config", "--json"}, {"exec", "--json=false"}} {
		if ExplicitCodexJSON(args) {
			t.Fatal(args)
		}
	}
}

func TestExhaustionZcodeResetCannotNegateClass(t *testing.T) {
	base := recordedZcode(t)
	for name, change := range map[string]func(string) string{
		"long": func(s string) string { return s },
		"short": func(s string) string {
			s = strings.ReplaceAll(s, "49h 8m 50s", "5m")
			s = strings.ReplaceAll(s, "176930", "300")
			return strings.ReplaceAll(s, "2026-10-04T22:14:57.003Z", "2026-10-02T21:11:07.003Z")
		},
		"absent": func(s string) string {
			s = strings.ReplaceAll(s, ". Your limit will reset at 2026-10-05 06:14:57 (reset after 49h 8m 50s)", "")
			s = strings.ReplaceAll(s, `,"retry_after":176930,"reset_at":"2026-10-04T22:14:57.003Z"`, "")
			return strings.ReplaceAll(s, ", 'retry-after': '176930'", "")
		},
		"malformed": func(s string) string { return strings.ReplaceAll(s, "2026-10-04T22:14:57.003Z", "not-a-date") },
		"wrong-type": func(s string) string {
			return strings.ReplaceAll(s, `"reset_at":"2026-10-04T22:14:57.003Z"`, `"reset_at":true`)
		},
		"contradictory": func(s string) string { return strings.ReplaceAll(s, "'retry-after': '176930'", "'retry-after': '10'") },
	} {
		t.Run(name, func(t *testing.T) {
			in := base
			in.Stderr = change(in.Stderr)
			e := ClassifyExhaustion(in)
			if !e.Eligible || e.RuleID != quota.ProviderExhaustionRule {
				t.Fatal(e)
			}
			if name != "long" && (e.ResetAt != nil || e.RawReset != "") {
				t.Fatal("uncertain reset persisted", e)
			}
		})
	}
	for _, suffix := range []string{"\nquoted prose", "\n{\"type\":\"turn.failed\"}", "\nAPICallError: incomplete"} {
		in := base
		in.Stderr += suffix
		if ClassifyExhaustion(in).Eligible {
			t.Fatal(suffix)
		}
	}
	in := base
	in.Stderr = strings.ReplaceAll(in.Stderr, "Weekly/Monthly Limit Exhausted", "Context window limit exhausted")
	if ClassifyExhaustion(in).Eligible {
		t.Fatal("local context exhaustion accepted")
	}
}
