package app

import (
	"bytes"
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

	// Truncation metadata added for bounded capture observation (D7 correction).
	Truncated        bool   // true when boundedWriter overflowed (stream exceeded cap)
	TruncationReason string // e.g. "overflow" or "deadline-cut"; empty when not truncated
}

const readinessTailBytes = 256

type boundedWriter struct {
	buf           *bytes.Buffer
	max           int
	overflow      *bool
	observedBytes int   // total bytes actually written (before any cap) — preserved independently
	truncated     *bool // distinct from overflow: set when any truncation/dropping occurs
}

func (w *boundedWriter) Write(p []byte) (int, error) {
	if len(p) == 0 {
		return 0, nil
	}
	w.observedBytes += len(p)
	if w.buf.Len() >= w.max {
		if w.overflow != nil && !*w.overflow {
			*w.overflow = true
		}
		if w.truncated != nil && !*w.truncated {
			*w.truncated = true
		}
		return len(p), nil
	}
	remaining := w.max - w.buf.Len()
	if remaining >= len(p) {
		n, err := w.buf.Write(p)
		if err != nil {
			return n, err
		}
		return len(p), nil
	}
	w.buf.Write(p[:remaining])
	if w.overflow != nil {
		*w.overflow = true
	}
	if w.truncated != nil {
		*w.truncated = true
	}
	return len(p), nil
}

// classifyReadiness maps one observed probe outcome to a typed class. It is pure
// so the full fixture table can be unit-tested without a child process.
func classifyReadiness(stdout, stderr string, exitCode int, timedOut bool, truncated bool, truncationReason string) readinessObservation {
	obs := readinessObservation{
		ExitCode:         exitCode,
		Truncated:        truncated,
		TruncationReason: truncationReason,
	}
	corpus := stdout + "\n" + stderr
	switch {
	// Truncation is never READY — always a concrete NON-READY class.
	// Even an empty/truncated output yields a defined deadline/process class,
	// never the empty default (Class = ""). Whitespace-only or truncated
	// bytes are preserved as observations; do not use TrimSpace to erase
	// byte-observation evidence.
	case truncated:
		obs.Class = ClassMalformedReply
		if timedOut {
			obs.Class = ClassDeadlineAfterOutput
		}
		obs.SawSentinel = containsSentinel(stdout) || containsSentinel(stderr)
	case timedOut:
		// Actual byte presence (not TrimSpace) determines whether output existed.
		if stdout == "" && stderr == "" {
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

func isDefiniteUnavailable(class string) bool {
	return class == "" || class == "missing" || class == string(ClassProcessFailure)
}

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
		return gateExcludeAgent,
			fmt.Sprintf("%s unavailable (%s) — confirm excluding it from this idea", entry.RosterID, entry.Reason),
			confirmCommand(root)
	}
}

func recheckCommand(root string) string {
	return fmt.Sprintf("parley preflight --dir %s", root)
}

// --- envelope recognition -----------------------------------------------------

var recognizedEnvelopeContentFields = []string{"content", "message", "text", "result"}

func strictEnvelopeChecks(obj map[string]any) error {
	if err := rejectCaseAliasedFields(obj); err != nil {
		return err
	}
	if err := rejectContradictoryErrorSignals(obj); err != nil {
		return err
	}
	if err := rejectMalformedFieldTypes(obj); err != nil {
		return err
	}
	for _, value := range obj {
		if inner, ok := value.(map[string]any); ok {
			if err := strictEnvelopeChecks(inner); err != nil {
				return err
			}
		}
	}
	return nil
}

var semanticEnvelopeKeys = []string{
	"role", "type", "status", "subtype", "error", "is_error", "isError",
	"success", "ok", "content", "message", "text", "result", "data", "payload", "response",
}

// isSemanticKey reports whether a decoded key string is one of the interpreted
// envelope semantics (case-insensitive). It is used by the token-level preDecodeScan.
func isSemanticKey(key string) bool {
	lower := strings.ToLower(key)
	for _, sem := range semanticEnvelopeKeys {
		if lower == strings.ToLower(sem) {
			return true
		}
	}
	return false
}

// tokenScanState carries the per-object duplicate-key set during recursive JSON
// token scanning. A fresh set is created for each object encountered.
func newTokenScanState() *tokenScanState {
	return &tokenScanState{
		seen: make(map[string]bool),
	}
}

type tokenScanState struct {
	seen map[string]bool
}

// preDecodeScan performs a token-level walk over raw JSON using
// encoding/json.Decoder to detect: duplicate/case-aliased semantic keys
// (per object), present-null semantic values, and basic structural errors.
// It recurses into nested objects and arrays, creating a fresh seen-key set
// for each new object so that duplicates in distinct nested objects are NOT
// treated as duplicates. Only interpreted semantic keys are tracked.
func preDecodeScan(raw string) error {
	trimmed := strings.TrimSpace(raw)
	if trimmed == "" || trimmed[0] != '{' {
		return nil // non-JSON falls through; no envelope to validate
	}
	dec := json.NewDecoder(strings.NewReader(trimmed))
	// Consume opening '{' by reading first token (must be object delimiter)
	tok, err := dec.Token()
	if err != nil {
		return fmt.Errorf("invalid JSON start: %w", err)
	}
	if tok != json.Delim('{') {
		return fmt.Errorf("expected JSON object start, got %v", tok)
	}
	return scanObject(dec, newTokenScanState())
}

func scanObject(dec *json.Decoder, state *tokenScanState) error {
	for dec.More() {
		keyTok, err := dec.Token()
		if err != nil {
			return fmt.Errorf("reading object key token: %w", err)
		}
		keyStr, ok := keyTok.(string)
		if !ok {
			return fmt.Errorf("expected string key in object, got %v (%T)", keyTok, keyTok)
		}
		// Read value token directly (no ':' delimiter token — Token() never returns ':' or ',').
		valTok, err := dec.Token()
		if err != nil {
			return fmt.Errorf("reading value for key %q: %w", keyStr, err)
		}
		// Check duplicate/case-aliased semantic key and present-null.
		if isSemanticKey(keyStr) {
			lowerKey := strings.ToLower(keyStr)
			if keyStr != lowerKey || state.seen[lowerKey] {
				return fmt.Errorf("duplicate/case-aliased semantic key %q in envelope", keyStr)
			}
			if valTok == nil {
				return fmt.Errorf("semantic key %q is present with null value; treat as present (not absent)", keyStr)
			}
			state.seen[lowerKey] = true
		}
		// Handle nested structures by recursing; primitives consume themselves.
		switch val := valTok.(type) {
		case json.Delim:
			switch val {
			case '{':
				if err := scanObject(dec, newTokenScanState()); err != nil {
					return err
				}
			case '[':
				if err := scanArray(dec, newTokenScanState()); err != nil {
					return err
				}
			default:
				return fmt.Errorf("unexpected value delimiter %v for key %q", val, keyStr)
			}
		case nil:
			// Already handled for semantic keys; non-semantic null values allowed.
		case bool, float64, string:
			// Primitive value — nothing further.
		default:
			return fmt.Errorf("unexpected value token %v for key %q", valTok, keyStr)
		}
	}
	// Consume closing '}' delimiter.
	closeTok, err := dec.Token()
	if err != nil {
		return fmt.Errorf("reading closing delimiter for object: %w", err)
	}
	if d, ok := closeTok.(json.Delim); !ok || d != '}' {
		return fmt.Errorf("expected '}' to close object, got %v (%T)", closeTok, closeTok)
	}
	return nil
}

func scanArray(dec *json.Decoder, parentState *tokenScanState) error {
	for {
		tok, err := dec.Token()
		if err != nil {
			return fmt.Errorf("reading array token: %w", err)
		}
		switch v := tok.(type) {
		case json.Delim:
			switch v {
			case ']':
				return nil
			case '{':
				if err := scanObject(dec, newTokenScanState()); err != nil {
					return err
				}
				// After nested object, loop continues for ',' or ']'
			case '[':
				if err := scanArray(dec, newTokenScanState()); err != nil {
					return err
				}
			default:
				return fmt.Errorf("unexpected array delimiter %v", v)
			}
		case string, float64, bool, nil:
			// Primitive array values — nothing to scan further.
			continue
		default:
			return fmt.Errorf("unexpected array token %v", tok)
		}
	}
}

func rejectCaseAliasedFields(obj map[string]any) error {
	casePairs := [][2]string{{"is_error", "isError"}, {"success", "ok"}}
	for _, pair := range casePairs {
		v1, ok1 := obj[pair[0]]
		v2, ok2 := obj[pair[1]]
		if ok1 && ok2 {
			switch b1 := v1.(type) {
			case bool:
				switch b2 := v2.(type) {
				case bool:
					if b1 != b2 {
						return fmt.Errorf("contradictory case-aliased fields %q=%v vs %q=%v", pair[0], b1, pair[1], b2)
					}
				}
			}
		}
	}
	return nil
}

func rejectContradictoryErrorSignals(obj map[string]any) error {
	isErr := false
	if b, ok := obj["is_error"].(bool); ok && b {
		isErr = true
	}
	if b, ok := obj["isError"].(bool); ok && b {
		isErr = true
	}
	successTrue := false
	if b, ok := obj["success"].(bool); ok && b {
		successTrue = true
	}
	okTrue := false
	if b, ok := obj["ok"].(bool); ok && b {
		okTrue = true
	}
	if isErr && (successTrue || okTrue) {
		return fmt.Errorf("contradictory error+success signal: is_error/isError=true with success/ok=true")
	}
	return nil
}

func rejectMalformedFieldTypes(obj map[string]any) error {
	for _, key := range []string{"role", "type", "status", "subtype"} {
		if _, ok := obj[key]; ok {
			v, present := obj[key]
			if !present || v == nil {
				return fmt.Errorf("malformed envelope field %q is present with null value; treat as present (not absent)", key)
			}
			switch v.(type) {
			case string:
			default:
				return fmt.Errorf("malformed envelope field %q has non-string type %T", key, v)
			}
		}
	}
	for _, key := range []string{"is_error", "isError", "success", "ok"} {
		if _, ok := obj[key]; ok {
			v, present := obj[key]
			if !present || v == nil {
				return fmt.Errorf("malformed envelope field %q is present with null value; treat as present (not absent)", key)
			}
			switch v.(type) {
			case bool:
			default:
				return fmt.Errorf("malformed envelope field %q has non-bool type %T", key, v)
			}
		}
	}
	return nil
}

func recognizeEnvelope(raw string) (isJSON, isError bool, payload string) {
	trimmed := strings.TrimSpace(raw)
	if trimmed == "" || trimmed[0] != '{' {
		return false, false, ""
	}
	var obj map[string]any
	if err := json.Unmarshal([]byte(trimmed), &obj); err != nil {
		return false, false, "" // malformed JSON → falls through to the text path
	}
	if scanErr := preDecodeScan(trimmed); scanErr != nil {
		return true, false, "" // malformed/duplicate/case-aliased or present-null semantic key
	}
	if strictErr := strictEnvelopeChecks(obj); strictErr != nil {
		return true, false, ""
	}
	if errorEnvelope(obj) {
		return true, true, ""
	}
	if s, ok := assistantPayload(obj); ok {
		return true, false, s
	}
	return true, false, "" // well-formed JSON, but not a recognized assistant envelope
}

var wrapperKeys = []string{"result", "data", "payload", "response"}

func errorEnvelope(obj map[string]any) bool {
	if signalsError(obj) {
		return true
	}
	for _, wrapper := range wrapperKeys {
		if inner, ok := obj[wrapper].(map[string]any); ok && errorEnvelope(inner) {
			return true
		}
	}
	return false
}

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

func assistantPayload(obj map[string]any) (string, bool) {
	if !hasRecognizedAssistantSchema(obj) || errorEnvelope(obj) {
		return "", false
	}
	keys := []string{"content", "message", "text", "result", "data", "payload", "response"}
	count := 0
	for _, key := range keys {
		if _, exists := obj[key]; exists {
			count++
		}
	}
	if count != 1 {
		return "", false
	}
	if value, ok := contentField(obj); ok {
		return value, true
	}
	for _, key := range keys {
		if inner, ok := obj[key].(map[string]any); ok {
			return assistantPayload(inner)
		}
	}
	return "", false
}

// hasRecognizedAssistantSchema requires explicit assistant identity or a known
// assistant/result envelope, including a nested attributed wrapper. Generic
// content keys alone never establish provenance.
func hasRecognizedAssistantSchema(obj map[string]any) bool {
	if obj == nil {
		return false
	}
	if role, present := obj["role"]; present && role != "assistant" {
		return false
	}
	if kind, present := obj["type"]; present {
		// These are explicit assistant/result envelopes. Generic content keys
		// and unknown types never establish who produced the text.
		if kind != "assistant" && kind != "message" && kind != "result" {
			return false
		}
		if kind == "result" {
			if subtype, present := obj["subtype"]; present && subtype != "success" {
				return false
			}
		}
		if kind == "assistant" || kind == "result" || obj["role"] == "assistant" {
			return true
		}
	}
	if obj["role"] == "assistant" {
		return true
	}
	for _, wrapper := range wrapperKeys {
		if inner, ok := obj[wrapper].(map[string]any); ok && hasRecognizedAssistantSchema(inner) {
			return true
		}
	}
	return false
}

func assistantRoleOrAbsent(obj map[string]any) bool {
	if obj == nil {
		return false
	}
	roleVal := obj["role"]
	if roleVal == nil {
		return true
	}
	switch v := roleVal.(type) {
	case string:
		return strings.EqualFold(strings.TrimSpace(v), "assistant")
	default:
		return false
	}
}

func contentField(obj map[string]any) (string, bool) {
	if obj == nil {
		return "", false
	}
	for _, key := range recognizedEnvelopeContentFields {
		if s, ok := obj[key].(string); ok && s != "" {
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

func containsSentinel(s string) bool {
	return strings.Contains(s, pongSentinel)
}

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

func providerFailureClass(text string) string {
	for _, r := range providerFailureRules {
		if r.re.MatchString(text) {
			return r.class
		}
	}
	return ""
}

func sanitizeTail(s string) string {
	t := strings.TrimSpace(s)
	t = scrubSecrets(t)
	if len(t) > readinessTailBytes {
		t = t[len(t)-readinessTailBytes:]
	}
	t = strings.Map(func(r rune) rune {
		if r == '\n' || r == '\t' {
			return r
		}
		if r < 0x20 || r == 0x7f {
			return ' '
		}
		return r
	}, t)
	t = scrubSecrets(t)
	return t
}

func scrubSecrets(s string) string {
	var value any
	if json.Unmarshal([]byte(s), &value) == nil {
		var redact func(any) any
		redact = func(v any) any {
			switch node := v.(type) {
			case map[string]any:
				for key, child := range node {
					k := strings.ToLower(strings.NewReplacer("_", "", "-", "").Replace(key))
					if strings.Contains(k, "token") || strings.Contains(k, "secret") || strings.Contains(k, "password") || k == "passwd" || strings.Contains(k, "apikey") || strings.Contains(k, "accesskey") || strings.Contains(k, "privatekey") || k == "authorization" {
						node[key] = "«redacted»"
					} else {
						node[key] = redact(child)
					}
				}
			case []any:
				for i, child := range node {
					node[i] = redact(child)
				}
			}
			return v
		}
		if raw, err := json.Marshal(redact(value)); err == nil {
			s = string(raw)
		}
	} else if trimmed := strings.TrimSpace(s); strings.HasPrefix(trimmed, "{") || strings.HasPrefix(trimmed, "[") {
		// A truncated/invalid JSON string may hide an escaped credential key or
		// cut through its value. Preserve the observation, omit unsafe raw text.
		return "«redacted malformed JSON diagnostic»"
	}
	// Full authorization bearer redaction (label + value), labeled JSON secrets
	// (token/secret/password/api_key/access_key/private_key with value), standalone
	// token shapes, and JWT fragments. Keeps label where appropriate (e.g. the key
	// label for labeled secrets) but never leaks the value or bearer token suffix.
	patterns := []*regexp.Regexp{
		// Authorization header — keep label, replace label + bearer + value.
		regexp.MustCompile(`(?i)(authorization\s*[:=]\s*)(bearer\s+)?[A-Za-z0-9._\-]+`),
		// Standalone bearer token value (after the word bearer) — full token.
		regexp.MustCompile(`(?i)\bbearer\s+[A-Za-z0-9._\-]+`),
		// Labeled secrets: key label + separator + value. Keeps label and separator.
		regexp.MustCompile(`(?i)(token|secret|password|passwd|api[_-]?key|access[_-]?key|private[_-]?key)(\s*[:=]\s*)[^"]+`),
		// Standalone token shape sk-... (at least 16 chars after sk-)
		regexp.MustCompile(`\bsk-[A-Za-z0-9_\-]{16,}`),
		// GitHub-style token shapes
		regexp.MustCompile(`\bgh[pousr]_[A-Za-z0-9]{20,}`),
		// Slack/xox-style
		regexp.MustCompile(`\bxox[baprs]-[A-Za-z0-9\-]{10,}`),
		// AWS AKIA
		regexp.MustCompile(`\bAKIA[0-9A-Z]{16}\b`),
		// JWT shape (3 dot-separated base64url segments, min lengths)
		regexp.MustCompile(`\beyJ[A-Za-z0-9_\-]{10,}\.[A-Za-z0-9_\-]{10,}\.[A-Za-z0-9_\-]{5,}\b`),
	}
	result := s
	// Pattern 0: full authorization line — label preserved, bearer + value redacted.
	result = patterns[0].ReplaceAllString(result, "$1«redacted»")
	// Pattern 1: standalone bearer value removed entirely.
	result = patterns[1].ReplaceAllString(result, "bearer «redacted»")
	// Pattern 2: labeled secrets — keep label + separator; replace value.
	quoted := regexp.MustCompile(`(?i)((?:token|secret|password|passwd|api[_-]?key|access[_-]?key|private[_-]?key)\s*[:=]\s*)"(?:\\.|[^"\\])*(?:"|$)`)
	result = quoted.ReplaceAllString(result, "$1\"«redacted»\"")
	result = patterns[2].ReplaceAllString(result, "$1$2«redacted»")
	// JSON credential fields (e.g. {"token":"..."} or {"api_key":"..."}) — redact value, keep key.
	jsonSecretPattern := regexp.MustCompile(`(?i)("(?:token|secret|password|passwd|api[_-]?key|access[_-]?key|private[_-]?key)"\s*:\s*)"(?:\\.|[^"\\])*(?:"|$)`)
	result = jsonSecretPattern.ReplaceAllString(result, "$1\"«redacted»\"")
	// Patterns 3-7: standalone token shapes — full replacement.
	for i := 3; i < len(patterns); i++ {
		result = patterns[i].ReplaceAllString(result, "«redacted»")
	}
	return result
}
