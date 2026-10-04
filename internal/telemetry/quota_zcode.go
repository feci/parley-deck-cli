package telemetry

import (
	"bytes"
	"encoding/json"
	"io"
	"regexp"
	"strings"
	"time"

	"parley-deck-cli/internal/quota"
)

var zcodeTerminal = regexp.MustCompile(`^Error: Turn execution failed(?: \(traceId: [A-Za-z0-9_-]+\))?$`)
var zcodeProviderPrefix = regexp.MustCompile(`^\[[A-Za-z0-9_./-]+\] \[429\]: `)
var zcodeResetSuffix = regexp.MustCompile(`(?i)^[.;]?\s*(?:(?:your limit will |will )?reset(?:s)? (?:at [0-9T:.Z+ -]+(?:\s+\(reset after [0-9hms ]+\))?|after [0-9hms ]+))?[.!]?$`)

func zcodeExhaustionMessage(message string) bool {
	message = zcodeProviderPrefix.ReplaceAllString(message, "")
	message = strings.TrimPrefix(message, "Your ")
	for _, rule := range []*regexp.Regexp{accountExhaustion, allowanceExhaustion} {
		m := rule.FindStringIndex(message)
		if m != nil && m[0] == 0 && zcodeResetSuffix.MatchString(strings.TrimSpace(message[m[1]:])) {
			return true
		}
	}
	return false
}

// This is deliberately a stderr adapter, not a native root-channel claim. The
// owner accepted subagent-source ambiguity. Quotes/envelopes are not unwrapped;
// malformed/mixed provider records cannot be skipped to find a positive later.
func classifyZcodeQuota(in QuotaInput, fallback quota.Evidence) quota.Evidence {
	records, ok := parseSDKCapture(in.Stderr)
	if !ok {
		return fallback
	}
	var result quota.Evidence
	var first, last time.Time
	for index, record := range records {
		recordInput := in
		for _, ob := range in.stderrObservations {
			if ob.End >= record.end {
				recordInput.ObservedAt = ob.At
				break
			}
		}
		var envelope map[string]json.RawMessage
		var fields map[string]json.RawMessage
		if json.Unmarshal([]byte(record.body), &envelope) != nil || len(envelope) != 1 || json.Unmarshal(envelope["error"], &fields) != nil {
			return fallback
		}
		var body struct {
			Message    string      `json:"message"`
			ResetAt    string      `json:"reset_at"`
			RetryAfter json.Number `json:"retry_after"`
		}
		if json.Unmarshal(envelope["error"], &body) != nil || body.ResetAt == "" && body.RetryAfter == "" || !zcodeExhaustionMessage(body.Message) {
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
		if _, ok := fields["reset_at"]; ok && body.ResetAt == "" {
			return fallback
		}
		if _, ok := fields["retry_after"]; ok && body.RetryAfter == "" {
			return fallback
		}
		n := nativeQuotaError{Status: 429, Message: body.Message, ResetAt: body.ResetAt, RetryAfter: string(body.RetryAfter), Provenance: "zcode.stderr.owner-deviation.2026-10-04"}
		e := refineNativeQuota(recordInput, n)
		if !e.Eligible || e.ResetAt == nil {
			return fallback
		}
		if record.header != "" {
			n.RetryAfter = record.header
			check := refineNativeQuota(recordInput, n)
			if !check.Eligible || check.ResetAt == nil || absReset(check.ResetAt.Sub(*e.ResetAt)) > time.Second {
				return fallback
			}
		}
		if index == 0 {
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
