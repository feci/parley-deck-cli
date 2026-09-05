package app

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strings"
	"time"
)

// ReadinessClass is the typed readiness/liveness observation (D7): what the
// hosted-PONG probe actually OBSERVED, kept distinct from the readiness verdict
// the probe implies. Only ClassReady is "available"; every other class is a
// distinct observation that must not be silently collapsed into "dead" or
// "ready". A timeout is an observation, not a diagnosis of a hang.
type ReadinessClass string

const (
	// ClassReady — an exact PONG assistant response extracted from a recognized
	// envelope. The sole available verdict.
	ClassReady ReadinessClass = "ready"
	// ClassMalformedReply — the process exited 0 and produced output we could
	// not parse into a ready shape (echoed instruction, bullet, fence, unknown
	// text, or malformed JSON).
	ClassMalformedReply ReadinessClass = "malformed-reply"
	// ClassExitedEmpty — the process exited 0 with no parseable output.
	ClassExitedEmpty ReadinessClass = "process-exited-empty"
	// ClassDeadlineNoOutput — the probe deadline elapsed with no output/activity.
	ClassDeadlineNoOutput ReadinessClass = "deadline-no-output"
	// ClassDeadlineAfterOutput — the probe deadline elapsed after some output.
	ClassDeadlineAfterOutput ReadinessClass = "deadline-after-output"
	// ClassProviderFailure — a classified provider-side failure (auth, billing,
	// rate-limit, overload, model-not-found). Environmental, not the agent being
	// dead. Blocking but never an exclusion.
	ClassProviderFailure ReadinessClass = "provider-failure"
	// ClassProcessFailure — a non-provider failure: non-zero exit, command-build
	// error, or start error.
	ClassProcessFailure ReadinessClass = "process-failure"
)

// readinessObservation is the typed probe result carried into the roster table
// and (through it) into the readiness gates. Ready is true only for ClassReady.
type readinessObservation struct {
	Class         ReadinessClass
	Ready         bool
	ExitCode      int    // -1 when the process never ran or was cut off by the probe deadline
	SawSentinel   bool   // "PONG" present in output but not an exact ready shape
	ProviderClass string // provider sub-class (auth/rate-limit/...), empty otherwise
	StdoutTail    string // sanitized, <= readinessTailBytes
	StderrTail    string // sanitized, <= readinessTailBytes
	BuffersStdout bool   // agent declared Spec.BuffersStdout
	Duration      time.Duration
}

const readinessTailBytes = 256

// classifyReadiness maps one observed probe outcome to a typed class. It is pure
// so the full fixture table can be unit-tested without a child process.
//
// Order matters on exit 0: a JSON error/status envelope is a provider failure
// and wins over any embedded "PONG" substring (an error is a failure even when
// its subtype claims success, and even on a clean exit-zero wrapper); then the
// recognized assistant/result envelope's exact payload; then the plain-text
// exact sentinel. A "PONG" that only survives as a substring (echo, bullet,
// fence, malformed JSON, a role-tagged user/system/tool message, or an
// unrecognized string-valued key) is never ready.
func classifyReadiness(stdout, stderr string, exitCode int, timedOut bool) readinessObservation {
	obs := readinessObservation{ExitCode: exitCode}
	corpus := stdout + "\n" + stderr
	switch {
	case timedOut:
		if strings.TrimSpace(stdout) == "" && strings.TrimSpace(stderr) == "" {
			obs.Class = ClassDeadlineNoOutput
		} else {
			obs.Class = ClassDeadlineAfterOutput
		}
	case exitCode != 0:
		obs.Class = ClassProcessFailure
		if pc := providerFailureClass(corpus); pc != "" {
			obs.Class = ClassProviderFailure
			obs.ProviderClass = pc
		}
	default:
		isJSON, isErr, payload := recognizeEnvelope(stdout)
		if isJSON && isErr {
			// A structured error/status envelope is a provider-side failure
			// signal, never readiness, whatever a nested success/subtype flag
			// claims. It must never be auto-excluded: a JSON error is a
			// provider status field, not an arbitrary string of process text,
			// so it is classified as a provider failure and does NOT fall back
			// to an excludable process failure.
			obs.Class = ClassProviderFailure
			if pc := providerFailureClass(corpus); pc != "" {
				obs.ProviderClass = pc
			}
		} else if isJSON && strings.TrimSpace(payload) == pongSentinel {
			obs.Class = ClassReady
		} else if isExactPONG(stdout) {
			obs.Class = ClassReady
		} else if containsSentinel(stdout) {
			obs.Class = ClassMalformedReply
			obs.SawSentinel = true
		} else if strings.TrimSpace(stdout) == "" && strings.TrimSpace(stderr) == "" {
			obs.Class = ClassExitedEmpty
		} else {
			obs.Class = ClassMalformedReply
		}
	}
	if obs.Class == ClassReady {
		obs.Ready = true
	}
	obs.StdoutTail = sanitizeTail(stdout)
	obs.StderrTail = sanitizeTail(stderr)
	return obs
}

// readinessReason renders the human-facing reason string for a non-ready entry.
func readinessReason(obs readinessObservation) string {
	switch obs.Class {
	case ClassProviderFailure:
		if obs.ProviderClass != "" {
			return "provider-failure:" + obs.ProviderClass
		}
		return string(ClassProviderFailure)
	case ClassProcessFailure:
		if obs.ExitCode >= 0 {
			return fmt.Sprintf("process-failure:exit-%d", obs.ExitCode)
		}
		return string(ClassProcessFailure)
	default:
		return string(obs.Class)
	}
}

// isDefiniteUnavailable reports whether a non-ready class is a definite
// unavailability that keeps the explicit operator-exclusion path (a missing CLI
// or a plain non-provider process failure). Everything else — provider failures
// and the ambiguous readiness classes — must never be auto-excluded.
func isDefiniteUnavailable(class string) bool {
	return class == "" || class == "missing" || class == string(ClassProcessFailure)
}

// readinessGateFor maps a non-ready roster entry to the gate it raises.
// Definite unavailability keeps the explicit exclusion gate (confirmCommand);
// provider failures and ambiguous readiness open a blocking readiness-resolution
// gate that is NOT an automatic exclusion (see D7).
func readinessGateFor(entry rosterEntry, root string) (kind gateKind, detail string, confirm string) {
	switch {
	case entry.Class == string(ClassProviderFailure):
		return gateProviderFailure,
			fmt.Sprintf("%s reports a provider failure (%s) — resolve the provider/auth state; the agent is not excluded", entry.RosterID, entry.Reason),
			recheckCommand(root)
	case entry.Class == string(ClassMalformedReply),
		entry.Class == string(ClassExitedEmpty),
		entry.Class == string(ClassDeadlineNoOutput),
		entry.Class == string(ClassDeadlineAfterOutput):
		return gateResolveReadiness,
			fmt.Sprintf("%s readiness is unresolved (%s) — investigate and re-check; the agent is not excluded", entry.RosterID, entry.Reason),
			recheckCommand(root)
	default:
		// missing / empty class (presence check) / process failure: definite.
		return gateExcludeAgent,
			fmt.Sprintf("%s unavailable (%s) — confirm excluding it from this idea", entry.RosterID, entry.Reason),
			confirmCommand(root)
	}
}

// recheckCommand is the resolution a blocking readiness gate advertises: the
// operator investigates, then re-runs the check. It deliberately omits --yes,
// which would only confirm an exclusion and does not resolve readiness.
func recheckCommand(root string) string {
	return fmt.Sprintf("parley preflight --dir %s", root)
}

// --- envelope recognition -----------------------------------------------------

// recognizedEnvelopeContentFields are the assistant/result content field names
// accepted in a recognized structured (JSON) envelope, in priority order. The set
// is deliberately explicit and narrow: only field names a known assistant/result
// envelope uses for its output are trusted. Arbitrary string-valued keys (e.g.
// "answer", "pong", "output", "reply", "response", or any other bare key) are NOT
// assistant output and are intentionally absent — an object carrying only such a
// key is an unrecognized schema and stays non-ready. ("result" is both a content
// key and a recognized wrapper key; see wrapperKeys.)
var recognizedEnvelopeContentFields = []string{"content", "message", "text", "result"}

// recognizeEnvelope parses raw as a single JSON object and reports:
//
//	isJSON   — raw is one JSON object (so structured handling applies)
//	isError  — the object is an error/failure envelope (never ready)
//	payload  — the assistant/result content string, or "" when the object is
//	           not a recognized assistant/result envelope
//
// Only a single, standalone JSON object is recognized. Streaming JSONL — several
// objects on one line, or line-delimited objects — is NOT parsed or reassembled;
// a JSONL stream falls through to the plain-text path and is never ready off an
// echoed or streamed input line. The exact supported schema set is documented in
// the handoff note and pinned by the positive fixtures below.
func recognizeEnvelope(raw string) (isJSON, isError bool, payload string) {
	trimmed := strings.TrimSpace(raw)
	if trimmed == "" || trimmed[0] != '{' {
		return false, false, ""
	}
	var obj map[string]any
	if err := json.Unmarshal([]byte(trimmed), &obj); err != nil {
		return false, false, "" // malformed JSON → falls through to the text path
	}
	if errorEnvelope(obj) {
		return true, true, ""
	}
	if s, ok := assistantPayload(obj); ok {
		return true, false, s
	}
	return true, false, "" // well-formed JSON, but not a recognized assistant envelope
}

// wrapperKeys are the recognized envelope wrapper keys that may nest the
// assistant/result content, one level deep.
var wrapperKeys = []string{"result", "data", "payload", "response"}

// errorEnvelope reports whether obj is an error/failure envelope, at the top
// level or one wrapper level deep. The signal is an explicit provider status
// field — a non-empty "error" value, is_error/isError == true, a type/status in
// the error set, or success/ok == false. A nested "success"/"subtype" flag does
// NOT override an error signal, and is_error:true wins whatever the subtype
// claims.
func errorEnvelope(obj map[string]any) bool {
	if signalsError(obj) {
		return true
	}
	for _, wrapper := range wrapperKeys {
		if inner, ok := obj[wrapper].(map[string]any); ok && signalsError(inner) {
			return true
		}
	}
	return false
}

// signalsError checks one object for the explicit error/status fields that make
// it a provider/result-level failure rather than a success envelope.
func signalsError(obj map[string]any) bool {
	if v, ok := obj["error"]; ok && !isEmptySignal(v) {
		return true
	}
	for _, key := range []string{"is_error", "isError"} {
		if b, ok := obj[key].(bool); ok && b {
			return true
		}
	}
	for _, key := range []string{"type", "status"} {
		if s, ok := obj[key].(string); ok {
			switch strings.ToLower(strings.TrimSpace(s)) {
			case "error", "failure", "fault", "exception", "failed":
				return true
			}
		}
	}
	for _, key := range []string{"success", "ok"} {
		if b, ok := obj[key].(bool); ok && !b {
			return true
		}
	}
	return false
}

// assistantPayload extracts the assistant content string from a recognized
// non-error envelope. It returns (payload, false) when obj is NOT a recognized
// assistant/result envelope: a role-tagged user/system/tool message, a wrapped
// non-assistant message, or an object with no recognized content key. Only an
// explicit "assistant" role (or a role-less result object) may supply content.
func assistantPayload(obj map[string]any) (string, bool) {
	if !assistantRoleOrAbsent(obj) {
		return "", false
	}
	if s, ok := contentField(obj); ok {
		return s, true
	}
	for _, wrapper := range wrapperKeys {
		if inner, ok := obj[wrapper].(map[string]any); ok {
			if !assistantRoleOrAbsent(inner) {
				continue
			}
			if s, ok := contentField(inner); ok {
				return s, true
			}
		}
	}
	return "", false
}

// assistantRoleOrAbsent reports whether obj presents no role at all (a raw
// assistant/result object) or an explicit "assistant" role. A user/system/tool
// role is an echoed input line, not assistant output.
func assistantRoleOrAbsent(obj map[string]any) bool {
	role, ok := obj["role"].(string)
	if !ok {
		return true
	}
	return strings.EqualFold(strings.TrimSpace(role), "assistant")
}

func contentField(obj map[string]any) (string, bool) {
	for _, key := range recognizedEnvelopeContentFields {
		if s, ok := obj[key].(string); ok {
			return s, true
		}
	}
	return "", false
}

func isEmptySignal(v any) bool {
	switch t := v.(type) {
	case nil:
		return true
	case string:
		return strings.TrimSpace(t) == ""
	case bool:
		return !t
	case map[string]any:
		return len(t) == 0
	case []any:
		return len(t) == 0
	default:
		return false
	}
}

// containsSentinel reports whether the raw output contains the PONG token as a
// substring. It feeds only the malformed-vs-silent classification — NEVER a
// ready PASS (an echoed instruction, bullet, or fence is not readiness).
func containsSentinel(s string) bool {
	return strings.Contains(s, pongSentinel)
}

// --- provider vs process failure split ------------------------------------------

// providerFailureRules is the minimal provider-side classifier for the probe. It
// mirrors the provider classes of internal/runner/failclass.go (auth, billing,
// rate-limit, overloaded, model-not-found); everything else is a plain process
// failure. The runner's full classifier stays authoritative for run-level
// classification — the probe only needs the provider/process split.
var providerFailureRules = []struct {
	class string
	re    *regexp.Regexp
}{
	{"rate-limit", regexp.MustCompile(`(?i)rate[-_ ]?limit|rate_limit_error|usageLimitExceeded|usage limit|session limit|too many requests|(^|[^0-9])429([^0-9]|$)`)},
	{"auth", regexp.MustCompile(`(?i)authentication[_ ](failed|error|required)|unauthorized|forbidden|permission_error|oauth[_ ](org|error|failed|not allowed)|oauth_org_not_allowed|api key not valid|invalid api key|(^|[^0-9])401([^0-9]|$)`)},
	{"billing", regexp.MustCompile(`(?i)billing[_ ](error|failed|required)|payment|credit (error|exhausted|balance)|quota exceeded`)},
	{"overloaded", regexp.MustCompile(`(?i)overloaded|overloaded_error|serverOverloaded|server[_ ]error|internalServerError`)},
	{"model-not-found", regexp.MustCompile(`(?i)model[_ ]not[_ ]found|model not found|not_found_error|unknown model`)},
}

// providerFailureClass returns the provider sub-class matching text, or "" when
// the text is not a recognized provider failure.
func providerFailureClass(text string) string {
	for _, r := range providerFailureRules {
		if r.re.MatchString(text) {
			return r.class
		}
	}
	return ""
}

// sanitizeTail bounds a diagnosis tail to readinessTailBytes and strips control
// characters other than newline/tab so a hostile or ANSI-laden blob cannot
// smuggle escapes into a report.
func sanitizeTail(s string) string {
	t := strings.TrimSpace(s)
	if len(t) > readinessTailBytes {
		t = t[len(t)-readinessTailBytes:]
	}
	return strings.Map(func(r rune) rune {
		if r == '\n' || r == '\t' {
			return r
		}
		if r < 0x20 || r == 0x7f {
			return ' '
		}
		return r
	}, t)
}
