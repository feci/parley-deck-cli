package telemetry

import (
	"encoding/json"
	"regexp"
	"strconv"
	"strings"
	"time"

	"parley-deck-cli/internal/quota"
)

// AdapterQuotaSupport describes recognizer authority, not a claim that an adapter
// name or an is_error bit authenticates provider errors. Source details and fixture
// limitations live in testdata/quota/README.md.
type AdapterQuotaSupport struct{ Adapter, Status, Source, Limitation string }

func QuotaSupport() []AdapterQuotaSupport {
	return []AdapterQuotaSupport{
		{"zcode", "diagnostic-only", "zcode-app-cli 3.7.7-13/vendor/zcode.cjs: streamText onError; runPrompt", "SDK error stream also accepts tool callback failures; recorded excerpt omits its native envelope"},
		{"codex", "diagnostic-only", "codex-cli 0.159.3; native source not established in this worktree", "turn.failed alone does not establish terminal provider-error provenance"},
		{"kimi", "diagnostic-only", "kimi 0.42.0 installed binary; stream-json role envelopes", "no located terminal provider-error discriminator"},
		{"claude", "diagnostic-only", "configured --output-format text", "assistant and error prose are not distinguishable"},
	}
}

// QuotaInput carries execution facts from the supervisor, not model claims. Raw
// streams are consumed transiently; Evidence contains no transcript or environment.
type QuotaInput struct {
	Adapter, InvocationID       string
	Stdout, Stderr              string
	ObservedAt                  time.Time
	ExitCode                    *int
	StructuredFailure           bool
	ValidArtifact, LaterSuccess bool
	Watchdog                    string
	Truncated                   bool
}

// ClassifyQuota deliberately fails closed until an adapter's native error channel
// is established. Do not wire generic failure classifiers or regex tails into
// refineNativeQuota: a quote followed by failure is still a quote.
func ClassifyQuota(in QuotaInput) quota.Evidence {
	e := quota.Evidence{Adapter: in.Adapter, InvocationID: in.InvocationID, ObservedAt: in.ObservedAt.UTC()}
	switch {
	case in.ValidArtifact || in.LaterSuccess:
		e.Reason = "completed artifact or later success"
	case in.Watchdog != "":
		e.Reason = "watchdog or cancellation"
	case in.Truncated:
		e.Reason = "incomplete capture"
	case (in.ExitCode == nil || *in.ExitCode == 0) && !in.StructuredFailure:
		e.Reason = "no failed invocation"
	default:
		e.Reason = "diagnostic-only: native terminal provider-error provenance unsupported"
	}
	return e
}

// nativeQuotaError is deliberately private. Only a provenance-established adapter
// recognizer may construct it in production; at present none does. Semantics are
// tested separately so a future recognizer cannot weaken the policy predicate.
type nativeQuotaError struct {
	Status                                   int
	Message, ResetAt, RetryAfter, Provenance string
}

var (
	accountExhaustion   = regexp.MustCompile(`(?i)\b(account (credit|quota)|account[- ]credit balance|credits? balance) (is |has been )?(exhausted|depleted|insufficient)|\binsufficient (account )?credits?\b|\baccount (credit|quota) (limit )?(exceeded|exhausted)\b`)
	allowanceExhaustion = regexp.MustCompile(`(?i)\b(daily|weekly|monthly|weekly/monthly|hourly|[0-9]+[- ]hour) (usage |token |spending )?(limit|allowance) (has been |is )?(exhausted|reached|exceeded)\b`)
	longAllowance       = regexp.MustCompile(`(?i)^(daily|weekly|monthly|weekly/monthly)\b`)
	hourAllowance       = regexp.MustCompile(`(?i)^([0-9]+)[- ]hour\b`)
	resetAfter          = regexp.MustCompile(`(?i)\breset after ([^);\n]+)`)
	resetAtText         = regexp.MustCompile(`(?i)\breset(?:s| will reset)? at ([^);\n]+)`)
	durationParts       = regexp.MustCompile(`(?i)^([0-9]+h)?\s*([0-9]+m)?\s*([0-9]+s)?$`)
)

func refineNativeQuota(in QuotaInput, native nativeQuotaError) quota.Evidence {
	e := ClassifyQuota(in)
	if in.ValidArtifact || in.LaterSuccess || in.Watchdog != "" || in.Truncated || in.InvocationID == "" || in.ObservedAt.IsZero() || native.Provenance == "" || (in.ExitCode == nil && !in.StructuredFailure) || (in.ExitCode != nil && *in.ExitCode == 0 && !in.StructuredFailure) {
		return e
	}
	e.Provenance = native.Provenance
	if native.Status == 400 || native.Status == 401 || native.Status == 403 || native.Status >= 500 {
		e.Reason = "non-authorizing HTTP status"
		return e
	}
	account := accountExhaustion.FindString(native.Message)
	allowance := allowanceExhaustion.FindString(native.Message)
	if account == "" && allowance == "" {
		e.Reason = "no explicit exhausted account or named allowance"
		return e
	}
	// Persist the matched decisive phrase, not arbitrary message text. This allowlist
	// avoids copying a credential embedded elsewhere in the provider error.
	e.Excerpt = account
	if e.Excerpt == "" {
		e.Excerpt = allowance
	}
	raw, reset, stated, ok := quotaReset(native, in.ObservedAt)
	e.RawReset = raw
	e.ResetAt = reset
	if !ok {
		e.Reason = "contradictory, past or unparseable stated reset"
		return e
	}
	if stated {
		if reset.Sub(in.ObservedAt) < quota.MinimumReset {
			e.Reason = "known reset below 60 minutes"
			return e
		}
		e.RuleID = "quota.named-reset-ge-60m.v1"
	} else {
		long := longAllowance.MatchString(allowance)
		if m := hourAllowance.FindStringSubmatch(allowance); m != nil {
			n, _ := strconv.Atoi(m[1])
			long = n >= 24
		}
		if account == "" && !long {
			e.Reason = "no stated reset for allowance shorter than 24 hours"
			return e
		}
		if account != "" {
			e.RuleID = "quota.account-exhausted-no-reset.v1"
		} else {
			e.RuleID = "quota.allowance-ge-24h-no-reset.v1"
		}
	}
	e.Eligible = true
	e.Reason = "explicit provider exhaustion"
	return e
}

// Every stated reset is validated. Invalid strings never fall through to the
// unknown branch. Absolute values require a timezone; no machine-local timezone
// inference. Multiple independent reset expressions must agree within one second.
func quotaReset(n nativeQuotaError, observed time.Time) (string, *time.Time, bool, bool) {
	var times []time.Time
	var raw []string
	addAbsolute := func(s string) bool {
		t, err := time.Parse(time.RFC3339Nano, strings.TrimSpace(s))
		if err != nil {
			return false
		}
		times = append(times, t.UTC())
		raw = append(raw, s)
		return true
	}
	addDuration := func(s string) bool {
		if !durationParts.MatchString(strings.TrimSpace(s)) || strings.TrimSpace(s) == "" {
			return false
		}
		d, err := time.ParseDuration(strings.ReplaceAll(strings.TrimSpace(s), " ", ""))
		if err != nil || d <= 0 {
			return false
		}
		times = append(times, observed.Add(d).UTC())
		raw = append(raw, s)
		return true
	}
	invalid := func() (string, *time.Time, bool, bool) { return "[unparseable stated reset]", nil, true, false }
	if n.ResetAt != "" {
		if !addAbsolute(n.ResetAt) {
			return invalid()
		}
	}
	for _, m := range resetAfter.FindAllStringSubmatch(n.Message, -1) {
		if !addDuration(m[1]) {
			return invalid()
		}
	}
	for _, m := range resetAtText.FindAllStringSubmatch(n.Message, -1) {
		if !addAbsolute(strings.TrimSpace(m[1])) {
			return invalid()
		}
	}
	// Any reset word left unparsed is a stated but unsupported reset, never unknown.
	low := strings.ToLower(n.Message)
	if strings.Contains(low, "reset") && len(times) == 0 {
		return invalid()
	}
	if n.RetryAfter != "" { // Retry-After alone never supplies exhaustion semantics.
		seconds, err := strconv.ParseInt(n.RetryAfter, 10, 64)
		if err != nil || seconds <= 0 || seconds > 315360000 {
			return invalid()
		}
		times = append(times, observed.Add(time.Duration(seconds)*time.Second).UTC())
		raw = append(raw, n.RetryAfter)
	}
	if len(times) == 0 {
		return "", nil, false, true
	}
	for _, t := range times {
		delta := t.Sub(times[0])
		if !t.After(observed) || delta > time.Second || delta < -time.Second {
			return strings.Join(raw, "; "), nil, true, false
		}
	}
	reset := times[0]
	return strings.Join(raw, "; "), &reset, true, true
}

// QuotaEvidence uses the collector's bounded captures without unwrapping any
// assistant/tool envelope. It is called once the supervised invocation has stopped.
func (c *Collector) QuotaEvidence(in QuotaInput) quota.Evidence {
	c.mu.Lock()
	defer c.mu.Unlock()
	in.Adapter = c.adapter
	in.Stdout = string(c.quotaStdout)
	in.Stderr = string(c.quotaStderr)
	in.Truncated = in.Truncated || c.quotaTruncated
	return ClassifyQuota(in)
}

// DecodeQuotaEvidence is shared by event projections; malformed evidence is never
// interpreted as an authorization or as successful membership reconciliation.
func DecodeQuotaEvidence(v any) (*quota.Evidence, error) {
	b, err := json.Marshal(v)
	if err != nil {
		return nil, err
	}
	var e quota.Evidence
	err = json.Unmarshal(b, &e)
	return &e, err
}
