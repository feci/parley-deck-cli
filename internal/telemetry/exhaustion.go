package telemetry

import (
	"encoding/json"
	"regexp"
	"strings"

	"parley-deck-cli/internal/quota"
)

// CodexContractVersion recognizes the discovered version, not a model label.
func CodexContractVersion(version string) string {
	switch strings.TrimSpace(version) {
	case "0.165.0", "codex-cli 0.165.0":
		return "0.165.0"
	}
	return ""
}

// ExplicitCodexJSON reads the effective argv. A value or prompt is not a flag.
func ExplicitCodexJSON(args []string) bool {
	execMode, jsonMode := false, false
	values := map[string]bool{"-c": true, "--config": true, "-m": true, "--model": true,
		"-C": true, "--cd": true, "-s": true, "--sandbox": true, "-o": true,
		"--output-last-message": true, "-i": true, "--image": true, "--add-dir": true,
		"-p": true, "--profile": true, "--color": true, "--output-schema": true,
		"--enable": true, "--disable": true}
	for i := 0; i < len(args); i++ {
		a := args[i]
		if a == "--" {
			break
		}
		if values[a] {
			i++
			continue
		}
		if a == "exec" && !execMode {
			execMode = true
			continue
		}
		if a == "--json" && execMode {
			jsonMode = true
			continue
		}
		if a == "resume" && execMode {
			continue
		}
		if !strings.HasPrefix(a, "-") || a == "-" {
			break
		}
	}
	return execMode && jsonMode
}

// ClassifyExhaustion is independent of the legacy 60-minute quota predicate.
// It returns native classification only. Validation binds the final step proof.
func ClassifyExhaustion(in QuotaInput) quota.Evidence {
	e := quotaPreconditions(in)
	if in.ValidArtifact || in.LaterSuccess || in.Watchdog != "" || in.Truncated || in.InvocationID == "" || in.ObservedAt.IsZero() {
		return e
	}
	switch in.Adapter {
	case "codex":
		if !in.ExplicitJSON || CodexContractVersion(in.CLIVersion) == "" {
			e.Reason = "diagnostic-only: explicit Codex JSON and exact CLI 0.165.0 required"
			return e
		}
		return classifyCodexExhaustion(in, e)
	case "zcode":
		if in.ExitCode == nil || *in.ExitCode == 0 {
			return e
		}
		return classifyZcodeExhaustion(in, e)
	}
	return e
}

func classifyCodexExhaustion(in QuotaInput, e quota.Evidence) quota.Evidence {
	thread, turn, terminal := false, false, false
	message := ""
	var errorsSeen []string
	for _, line := range strings.Split(strings.TrimSpace(in.Stdout), "\n") {
		var event map[string]json.RawMessage
		if terminal || !uniqueJSONKeys([]byte(line)) || json.Unmarshal([]byte(line), &event) != nil {
			e.Reason = "Codex JSON framing or terminal sequence rejected"
			return e
		}
		var kind string
		if json.Unmarshal(event["type"], &kind) != nil {
			return e
		}
		switch kind {
		case "thread.started":
			var id string
			if thread || turn || len(event) != 2 || json.Unmarshal(event["thread_id"], &id) != nil || id == "" {
				return e
			}
			thread = true
		case "turn.started":
			if !thread || turn || len(event) != 1 {
				return e
			}
			turn = true
		case "item.started", "item.updated", "item.completed":
			var item map[string]json.RawMessage
			if !thread || len(event) != 2 || json.Unmarshal(event["item"], &item) != nil || item == nil {
				return e
			}
		case "error":
			var text string
			if !thread || len(event) != 2 || json.Unmarshal(event["message"], &text) != nil {
				return e
			}
			errorsSeen = append(errorsSeen, text)
		case "turn.failed":
			var body map[string]json.RawMessage
			if !thread || !turn || len(event) != 2 || json.Unmarshal(event["error"], &body) != nil || len(body) != 1 || json.Unmarshal(body["message"], &message) != nil || !quota.CodexExhaustionMessage(message) {
				return e
			}
			terminal = true
		default:
			return e
		}
	}
	if !terminal {
		return e
	}
	for _, text := range errorsSeen {
		if text != message {
			return e
		}
	}
	// The source-backed terminal envelope itself establishes a failed turn.
	e.Eligible, e.RuleID, e.Provenance = true, quota.ProviderExhaustionRule, quota.CodexExhaustionProvenance
	e.Excerpt, e.Reason = message, "explicit provider exhaustion; source-backed Codex JSON"
	return e
}

var exhaustionResetTail = regexp.MustCompile(`(?i)^[.;]?\s*(?:(?:your limit will |will )?reset(?:s)? (?:at|after)\b[^\r\n]*)?[.!]?$`)
var accountTokenExhaustion = regexp.MustCompile(`(?i)^account (?:tokens?|token allowance|usage allowance) (?:is |are |has been )?(?:exhausted|depleted|exceeded)\b`)

func zcodeExhaustionClass(message string) string {
	message = zcodeProviderPrefix.ReplaceAllString(message, "")
	message = strings.TrimPrefix(message, "Your ")
	for _, rule := range []*regexp.Regexp{accountExhaustion, allowanceExhaustion, accountTokenExhaustion} {
		m := rule.FindStringIndex(message)
		if m != nil && m[0] == 0 && exhaustionResetTail.MatchString(strings.TrimSpace(message[m[1]:])) {
			return message[:m[1]]
		}
	}
	return ""
}

func classifyZcodeExhaustion(in QuotaInput, e quota.Evidence) quota.Evidence {
	records, err := parseSDKCapture(in.Stderr)
	if err != nil {
		e.Reason = "zcode evidence rejected: framing: " + err.Error()
		return e
	}
	class := ""
	for _, record := range records {
		var envelope map[string]json.RawMessage
		var body map[string]json.RawMessage
		if json.Unmarshal([]byte(record.body), &envelope) != nil || len(envelope) != 1 || json.Unmarshal(envelope["error"], &body) != nil {
			return e
		}
		var message string
		if json.Unmarshal(body["message"], &message) != nil {
			return e
		}
		current := zcodeExhaustionClass(message)
		if current == "" || (class != "" && !strings.EqualFold(current, class)) {
			return e
		}
		class = current
		for key, value := range body {
			switch key {
			case "message", "reset_at", "retry_after":
			case "status", "statusCode", "code":
				if string(value) != "429" && string(value) != `"429"` {
					return e
				}
			default:
				return e
			}
		}
	}
	if class == "" {
		return e
	}
	e.Eligible, e.RuleID, e.Provenance = true, quota.ProviderExhaustionRule, quota.ZcodeExhaustionProvenance
	e.Excerpt, e.Reason = class, "explicit provider exhaustion; reset unknown"
	terminalEnd := len(strings.TrimRight(in.Stderr, "\r\n"))
	for _, ob := range in.stderrObservations {
		if ob.End >= terminalEnd {
			e.ObservedAt = ob.At
			break
		}
	}
	// Only the unchanged legacy validator can supply a long, consistent reset.
	// Failure of that separate hint never negates proven exhaustion.
	if hint := ClassifyQuota(in); hint.Eligible && hint.ResetAt != nil {
		e.RawReset, e.ResetAt, e.ObservedAt = hint.RawReset, hint.ResetAt, hint.ObservedAt
		e.Reason = "explicit provider exhaustion; validated reset hint"
	}
	return e
}

func (c *Collector) ExhaustionEvidence(in QuotaInput) quota.Evidence {
	c.mu.Lock()
	defer c.mu.Unlock()
	in.Adapter, in.Stdout, in.Stderr = c.adapter, string(c.quotaStdout), string(c.quotaStderr)
	in.stderrObservations = c.quotaStderrObservations
	in.Truncated = in.Truncated || c.quotaTruncated
	return ClassifyExhaustion(in)
}
