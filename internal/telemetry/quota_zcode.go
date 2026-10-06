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
	reject := func(reason string) quota.Evidence {
		fallback.Reason = "zcode evidence rejected: " + reason
		return fallback
	}
	records, err := parseSDKCapture(in.Stderr)
	if err != nil {
		return reject("framing: " + err.Error())
	}
	if in.InvocationID == "" || in.ObservedAt.IsZero() {
		return reject("missing invocation identity or receipt")
	}
	// Use receipt of the terminal line, not process teardown or an inferred retry
	// observation. Multiple records in one write share this receipt.
	receipt := in.ObservedAt
	terminalEnd := len(strings.TrimRight(in.Stderr, "\r\n"))
	for _, ob := range in.stderrObservations {
		if ob.End >= terminalEnd {
			receipt = ob.At
			break
		}
	}
	var result quota.Evidence
	var first, last, previous time.Time
	for index, record := range records {
		recordReceipt := receipt
		for _, ob := range in.stderrObservations {
			if ob.End >= record.end {
				recordReceipt = ob.At
				break
			}
		}
		var envelope map[string]json.RawMessage
		var fields map[string]json.RawMessage
		if json.Unmarshal([]byte(record.body), &envelope) != nil || len(envelope) != 1 || json.Unmarshal(envelope["error"], &fields) != nil {
			return reject("malformed provider envelope")
		}
		var body struct {
			Message    string      `json:"message"`
			ResetAt    string      `json:"reset_at"`
			RetryAfter json.Number `json:"retry_after"`
		}
		if json.Unmarshal(envelope["error"], &body) != nil || body.ResetAt == "" && body.RetryAfter == "" || !zcodeExhaustionMessage(body.Message) {
			return reject("missing machine reset or explicit exhaustion class")
		}
		for key, value := range fields {
			switch key {
			case "message", "reset_at", "retry_after":
			case "status", "statusCode", "code":
				if string(value) != "429" && string(value) != `"429"` {
					return reject("mixed provider status")
				}
			default:
				return reject("unknown provider body field")
			}
		}
		if _, ok := fields["reset_at"]; ok && body.ResetAt == "" {
			return reject("empty machine reset")
		}
		if _, ok := fields["retry_after"]; ok && body.RetryAfter == "" {
			return reject("empty retry duration")
		}
		n := nativeQuotaError{Status: 429, Message: body.Message, ResetAt: body.ResetAt, RetryAfter: string(body.RetryAfter), Provenance: "zcode.stderr.owner-deviation.2026-10-04"}
		observed := recordReceipt
		// Absolute reset plus a relative duration identifies this attempt's
		// observation, even when an SDK aggregates retries into one dump.
		if n.ResetAt != "" {
			absolute, e := time.Parse(time.RFC3339Nano, n.ResetAt)
			if e != nil {
				return reject("unparseable absolute reset")
			}
			duration, present, valid := zcodeObservationDuration(n, record.header)
			if !valid {
				return reject("unparseable retry duration")
			}
			if present {
				observed = absolute.Add(-duration)
			}
		}
		if observed.After(recordReceipt.Add(time.Second)) || observed.After(receipt.Add(time.Second)) || !in.StartedAt.IsZero() && observed.Before(in.StartedAt.Add(-time.Second)) || !previous.IsZero() && observed.Before(previous.Add(-time.Second)) {
			return reject("contradictory inferred observation order/start/receipt")
		}
		if observed.After(previous) {
			previous = observed
		}
		recordInput := in
		recordInput.ObservedAt = observed
		e := refineNativeQuota(recordInput, n)
		if !e.Eligible || e.ResetAt == nil {
			return reject(e.Reason)
		}
		if record.header != "" {
			n.RetryAfter = record.header
			check := refineNativeQuota(recordInput, n)
			if !check.Eligible || check.ResetAt == nil || absReset(check.ResetAt.Sub(*e.ResetAt)) > time.Second {
				return reject("contradictory retry header/countdown/reset")
			}
		}
		if index == 0 {
			result = e
			first = *e.ResetAt
			last = first
		} else {
			if e.RuleID != result.RuleID || !strings.EqualFold(e.Excerpt, result.Excerpt) {
				return reject("mixed exhaustion class")
			}
			if e.ResetAt.Before(first) {
				first = *e.ResetAt
			}
			if e.ResetAt.After(last) {
				last = *e.ResetAt
			}
			if last.Sub(first) > time.Second {
				return reject("contradictory absolute resets")
			}
		}
	}
	if first.Sub(receipt) < quota.MinimumReset {
		return reject("known reset below 60 minutes at terminal receipt")
	}
	result.ResetAt = &first
	result.ObservedAt = receipt.UTC()
	return result
}

// Select an anchor only; quotaReset subsequently checks EVERY stated duration,
// display and machine value against it. No contradictory value is discarded.
func zcodeObservationDuration(n nativeQuotaError, header string) (time.Duration, bool, bool) {
	value := n.RetryAfter
	if value == "" {
		value = header
	}
	if value != "" {
		seconds, err := strconv.ParseInt(value, 10, 64)
		return time.Duration(seconds) * time.Second, true, err == nil && seconds > 0 && seconds <= 315360000
	}
	if m := resetAfter.FindStringSubmatch(n.Message); m != nil {
		s := strings.TrimSpace(m[1])
		d, err := time.ParseDuration(strings.ReplaceAll(s, " ", ""))
		return d, true, durationParts.MatchString(s) && err == nil && d > 0
	}
	return 0, false, true
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
