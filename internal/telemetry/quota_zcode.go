package telemetry

import (
	"bytes"
	"encoding/json"
	"io"
	"regexp"
	"strconv"
	"strings"
	"time"

	"parley-deck-cli/internal/quota"
)

var zcodeProviderWrapper = regexp.MustCompile(`^\s*(?:(?:lastError|\[cause\]):\s*)?(?:AI_APICallError|APICallError)(?:\s+\[AI_APICallError\])?:`)
var zcodeTerminal = regexp.MustCompile(`^Error: Turn execution failed(?: \(traceId: [A-Za-z0-9_-]+\))?$`)
var zcodeStatus = regexp.MustCompile(`^\s*statusCode: ([0-9]{3}),?\s*$`)
var zcodeBody = regexp.MustCompile(`^\s*responseBody: '(.*)',?\s*$`)
var zcodeRetry = regexp.MustCompile(`^\s*'retry-after': '([0-9]+)',?\s*$`)

// This is deliberately a stderr adapter, not a native root-channel claim. The
// owner accepted subagent-source ambiguity. Quotes/envelopes are not unwrapped;
// malformed/mixed provider records cannot be skipped to find a positive later.
func classifyZcodeQuota(in QuotaInput, fallback quota.Evidence) quota.Evidence {
	lines := strings.Split(strings.TrimRight(in.Stderr, "\r\n "), "\n")
	if len(lines) == 0 || !zcodeTerminal.MatchString(strings.TrimSpace(lines[len(lines)-1])) {
		return fallback
	}
	status, bodies, header := 0, 0, ""
	wrapperPending := false
	var result quota.Evidence
	var first, last time.Time
	offset := 0
	for _, line := range lines[:len(lines)-1] {
		offset += len(line) + 1
		if zcodeProviderWrapper.MatchString(line) {
			if wrapperPending || status != 0 {
				return fallback
			}
			wrapperPending = true
			continue
		}
		if m := zcodeStatus.FindStringSubmatch(line); m != nil {
			if status != 0 {
				return fallback
			}
			status, _ = strconv.Atoi(m[1])
			if status != 429 {
				return fallback
			}
			continue
		}
		if m := zcodeRetry.FindStringSubmatch(line); m != nil {
			if header != "" {
				return fallback
			}
			header = m[1]
			continue
		}
		if m := zcodeBody.FindStringSubmatch(line); m != nil {
			if status != 429 {
				return fallback
			}
			if !uniqueJSONKeys([]byte(m[1])) {
				return fallback
			}
			recordInput := in
			for _, ob := range in.stderrObservations {
				if ob.End >= offset-1 {
					recordInput.ObservedAt = ob.At
					break
				}
			}
			var body struct {
				Error struct {
					Message    string      `json:"message"`
					ResetAt    string      `json:"reset_at"`
					RetryAfter json.Number `json:"retry_after"`
				} `json:"error"`
			}
			if json.Unmarshal([]byte(m[1]), &body) != nil {
				return fallback
			}
			b := body.Error
			if b.ResetAt == "" && b.RetryAfter == "" {
				return fallback
			}
			// A present null/empty reset is unavailable, not an absent reset.
			var fields map[string]json.RawMessage
			var envelope map[string]json.RawMessage
			if json.Unmarshal([]byte(m[1]), &envelope) != nil || json.Unmarshal(envelope["error"], &fields) != nil || len(envelope) != 1 {
				return fallback
			}
			for key, value := range fields {
				switch key {
				case "message", "reset_at", "retry_after":
				case "status", "statusCode", "code":
					if string(value) != "429" && string(value) != `"429"` {
						return fallback
					}
				default:
					return fallback
				}
			}
			if _, ok := fields["reset_at"]; ok && b.ResetAt == "" {
				return fallback
			}
			if _, ok := fields["retry_after"]; ok && b.RetryAfter == "" {
				return fallback
			}
			n := nativeQuotaError{Status: status, Message: b.Message, ResetAt: b.ResetAt, RetryAfter: string(b.RetryAfter), Provenance: "zcode.stderr.owner-deviation.2026-10-04"}
			e := refineNativeQuota(recordInput, n)
			if !e.Eligible || e.ResetAt == nil {
				return e
			}
			if header != "" {
				n.RetryAfter = header
				h := refineNativeQuota(recordInput, n)
				if !h.Eligible || h.ResetAt == nil || absReset(h.ResetAt.Sub(*e.ResetAt)) > time.Second {
					return fallback
				}
			}
			if bodies == 0 {
				result = e
				first = *e.ResetAt
				last = first
			} else {
				if e.RuleID != result.RuleID || !strings.EqualFold(e.Excerpt, result.Excerpt) {
					return fallback
				}
				if e.ResetAt.Before(first) {
					first = *e.ResetAt
				}
				if e.ResetAt.After(last) {
					last = *e.ResetAt
				}
				if last.Sub(first) > time.Second {
					return fallback
				}
			}
			bodies++
			wrapperPending = false
			status = 0
			header = ""
			continue
		}
		// These fields outside the supported SDK record framing, or content
		// envelopes which merely quote them, invalidate the complete capture.
		if strings.Contains(line, "responseBody") || strings.Contains(line, "statusCode") || strings.Contains(line, "retry-after") || strings.Contains(line, `"role"`) || strings.Contains(line, "API Error:") || strings.HasPrefix(strings.TrimSpace(line), "Error:") || strings.HasPrefix(strings.TrimSpace(line), "TypeError:") || strings.HasPrefix(strings.TrimSpace(line), "[tool") || strings.HasPrefix(strings.TrimSpace(line), "[assistant") || strings.HasPrefix(strings.TrimSpace(line), "```") {
			return fallback
		}
	}
	if bodies == 0 || status != 0 || header != "" || wrapperPending {
		return fallback
	}
	result.ResetAt = &first
	return result
}

func absReset(d time.Duration) time.Duration {
	if d < 0 {
		return -d
	}
	return d
}

// Duplicate keys are contradictory evidence, including a hidden overwritten reset.
func uniqueJSONKeys(raw []byte) bool {
	d := json.NewDecoder(bytes.NewReader(raw))
	var value func() bool
	value = func() bool {
		t, err := d.Token()
		if err != nil {
			return false
		}
		delim, ok := t.(json.Delim)
		if !ok {
			return true
		}
		switch delim {
		case '{':
			seen := map[string]bool{}
			for d.More() {
				k, err := d.Token()
				key, ok := k.(string)
				if err != nil || !ok || seen[key] {
					return false
				}
				seen[key] = true
				if !value() {
					return false
				}
			}
			end, err := d.Token()
			return err == nil && end == json.Delim('}')
		case '[':
			for d.More() {
				if !value() {
					return false
				}
			}
			end, err := d.Token()
			return err == nil && end == json.Delim(']')
		}
		return false
	}
	if !value() {
		return false
	}
	_, err := d.Token()
	return err == io.EOF
}
