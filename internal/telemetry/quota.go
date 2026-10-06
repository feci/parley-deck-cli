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
		{"zcode", "supported-owner-deviation", "owner scope-reset answer 2026-10-04; zcode-app-cli 3.7.7-13 stderr responseBody", "subagent-source ambiguity explicitly accepted by owner; every provider record must agree"},
		{"codex", "diagnostic-only", "codex-cli 0.159.3; native source not established in this worktree", "turn.failed alone does not establish terminal provider-error provenance"},
		{"kimi", "diagnostic-only", "kimi 0.42.0 installed binary; stream-json role envelopes", "no located terminal provider-error discriminator"},
		{"claude", "diagnostic-only", "configured --output-format text", "assistant and error prose are not distinguishable"},
	}
}

// QuotaInput carries execution facts from the supervisor, not model claims. Raw
// streams are consumed transiently; Evidence contains no transcript or environment.
type quotaStreamObservation struct {
	End int
	At  time.Time
}

type QuotaInput struct {
	stderrObservations          []quotaStreamObservation
	Adapter, InvocationID       string
	Stdout, Stderr              string
	ObservedAt                  time.Time
	StartedAt                   time.Time // optional supervisor-recorded invocation start
	ExitCode                    *int
	StructuredFailure           bool
	ValidArtifact, LaterSuccess bool
	Watchdog                    string
	Truncated                   bool
}

// ClassifyQuota deliberately fails closed until an adapter's native error channel
// is established. Do not wire generic failure classifiers or regex tails into
// refineNativeQuota: a quote followed by failure is still a quote.
func quotaPreconditions(in QuotaInput) quota.Evidence {
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

// ClassifyQuota supports only the owner's bounded zcode stderr exception. The
// native structured channels remain diagnostic-only: their reset values are lost.
func ClassifyQuota(in QuotaInput) quota.Evidence {
	e := quotaPreconditions(in)
	if in.Adapter != "zcode" || in.ExitCode == nil || *in.ExitCode == 0 || in.ValidArtifact || in.LaterSuccess || in.Watchdog != "" || in.Truncated {
		return e
	}
	return classifyZcodeQuota(in, e)
}

// nativeQuotaError is deliberately private. Only a provenance-established adapter
// recognizer may construct it in production. Semantics are
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
	resetAfter          = regexp.MustCompile(`(?i)\breset after ([^();\n]+)`)
	resetAtText         = regexp.MustCompile(`(?i)\breset(?:s| will reset)? at ([^();\n]+)`)
	durationParts       = regexp.MustCompile(`(?i)^([0-9]+h)?\s*([0-9]+m)?\s*([0-9]+s)?$`)
)

func refineNativeQuota(in QuotaInput, native nativeQuotaError) quota.Evidence {
	e := quotaPreconditions(in)
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
	var displays []time.Time
	var raw []string
	invalid := func() (string, *time.Time, bool, bool) { return "[unparseable stated reset]", nil, true, false }
	absolute := func(s string) bool {
		t, err := time.Parse(time.RFC3339Nano, strings.TrimSpace(s))
		if err != nil {
			return false
		}
		times = append(times, t.UTC())
		raw = append(raw, s)
		return true
	}
	machine := false
	if n.ResetAt != "" {
		if !absolute(n.ResetAt) {
			return invalid()
		}
		machine = true
	}
	if n.RetryAfter != "" {
		seconds, err := strconv.ParseInt(n.RetryAfter, 10, 64)
		if err != nil || seconds <= 0 || seconds > 315360000 {
			return invalid()
		}
		times = append(times, observed.Add(time.Duration(seconds)*time.Second).UTC())
		raw = append(raw, n.RetryAfter)
		machine = true
	}
	for _, m := range resetAfter.FindAllStringSubmatch(n.Message, -1) {
		s := strings.TrimSpace(m[1])
		if s == "" || !durationParts.MatchString(s) {
			return invalid()
		}
		d, err := time.ParseDuration(strings.ReplaceAll(s, " ", ""))
		if err != nil || d <= 0 {
			return invalid()
		}
		times = append(times, observed.Add(d).UTC())
		raw = append(raw, s)
	}
	for _, m := range resetAtText.FindAllStringSubmatch(n.Message, -1) {
		s := strings.TrimSpace(m[1])
		if absolute(s) {
			continue
		}
		// Owner's bounded interpretation, never the machine's local timezone.
		t, err := time.Parse("2006-01-02 15:04:05", s)
		if err != nil || !machine {
			return invalid()
		}
		displays = append(displays, t)
		raw = append(raw, s)
	}
	rest := resetAfter.ReplaceAllString(n.Message, "")
	rest = resetAtText.ReplaceAllString(rest, "")
	if strings.Contains(strings.ToLower(rest), "reset") {
		return invalid()
	}
	if len(times) == 0 {
		return "", nil, false, true
	}
	earliest, latest := times[0], times[0]
	for _, t := range times {
		if t.Before(earliest) {
			earliest = t
		}
		if t.After(latest) {
			latest = t
		}
	}
	if !earliest.After(observed) || latest.Sub(earliest) > time.Second {
		return strings.Join(raw, "; "), nil, true, false
	}
	canonical := earliest
	for _, display := range displays {
		agrees := false
		for minutes := -12 * 60; minutes <= 14*60; minutes += 15 {
			instant := display.Add(-time.Duration(minutes) * time.Minute)
			lo, hi := earliest, latest
			if instant.Before(lo) {
				lo = instant
			}
			if instant.After(hi) {
				hi = instant
			}
			if hi.Sub(lo) <= time.Second {
				agrees = true
				earliest, latest = lo, hi
				break
			}
		}
		if !agrees {
			return strings.Join(raw, "; "), nil, true, false
		}
	}
	// Keep the earliest machine instant as canonical; display clocks only constrain agreement.
	return strings.Join(raw, "; "), &canonical, true, true
}

// QuotaEvidence uses the collector's bounded captures without unwrapping any
// assistant/tool envelope. It is called once the supervised invocation has stopped.
func (c *Collector) QuotaEvidence(in QuotaInput) quota.Evidence {
	c.mu.Lock()
	defer c.mu.Unlock()
	in.Adapter = c.adapter
	in.Stdout = string(c.quotaStdout)
	in.Stderr = string(c.quotaStderr)
	in.stderrObservations = c.quotaStderrObservations
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
