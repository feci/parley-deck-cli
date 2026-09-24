package evidence

import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"regexp"
	"strings"
	"time"

	"parley-deck-cli/internal/procctl"
)

// execute.go runs one criterion command and produces a typed CommandEvidence.
//
// Semantic attestation is format-aware:
//   - `go test -json` output is parsed as structured test2json events; executed
//     and failed case counts come from real Pass/Fail test events, not regexes.
//     The event decoder is strict: duplicate or case-aliased semantic fields
//     and malformed event lines fail closed (never collapse into a shell pass).
//     Package-level fail/build-fail events (a failed build runs no test cases)
//     fail closed through FailedPackages without touching the exact test-case
//     counts, so a masked pipeline exit cannot launder a package failure into
//     an overall pass.
//   - An explicit evidence envelope (a `PARLEY-EVIDENCE {...}` JSON line) may
//     assert executed_cases/failed_cases for claims test2json cannot express.
//   - Anything else is FormatShell: exit code and output hash only, no case
//     semantics — closure decides whether that is acceptable.

const (
	evidenceMaxLines = 100
	evidenceMaxBytes = 4096
	// captureMax bounds the output bytes RETAINED for parsing and diagnostics.
	// The full stream is still hashed (streaming), so the output binding covers
	// everything the command emitted; but a command that floods past the bound
	// cannot claim a semantic pass — overflow fails closed.
	captureMax = 4 << 20 // 4 MiB
	// waitDelay bounds how long Wait tolerates descendants holding our pipes
	// after the shell exits (a backgrounded grandchild must not hang evidence
	// collection or defeat context cancellation).
	waitDelay = 10 * time.Second
)

// envelopePrefix marks an explicit evidence envelope line in command output.
const envelopePrefix = "PARLEY-EVIDENCE "

// Envelope is the explicit evidence format for claims structured test JSON
// cannot prove. A command prints `PARLEY-EVIDENCE {"executed_cases":N,...}`.
type Envelope struct {
	ExecutedCases *int `json:"executed_cases"`
	FailedCases   *int `json:"failed_cases"`
}

// RunCriterion executes `command` via sh -c with cwd=root and returns the
// typed record. Executor is the asserted runtime identity of whoever ran it;
// verifier attestation happens later and is NOT stamped here.
//
// Hardening: the shell runs in its own process group and context cancellation
// kills the whole group (procctl), with WaitDelay so a grandchild holding our
// pipes cannot hang collection; captured output is bounded (overflow fails
// closed) while the hash still covers the full stream; and the persisted
// Command field is the SCRUBBED safe representation — the raw command text is
// never persisted, only its hash (commands routinely carry secrets in env
// assignments).
func RunCriterion(ctx context.Context, root, name, command, executor string) CriterionRecord {
	return RunCriterionDetailed(ctx, root, name, command, executor).Record
}

// CriterionExecution distinguishes a complete observed failing test from an
// interrupted process, truncated capture or malformed structured output. The
// latter remain failures for closure, but cannot prove a patch-induced
// regression. Complete does not certify opaque shell output or authenticity.
type CriterionExecution struct {
	Record   CriterionRecord `json:"record"`
	Complete bool            `json:"complete"`
}

// RunCriterionDetailed uses the same execution and parsing path as RunCriterion
// and adds typed observation completeness without changing persisted reports.
func RunCriterionDetailed(ctx context.Context, root, name, command, executor string) CriterionExecution {
	return RunCriterionControlled(ctx, root, name, command, executor, nil)
}

// CriterionStartControl serializes a durable stop check, process creation and
// identity publication. start creates a waiting supervisor, not the material
// command. The controller calls release only after publishing identity, while
// still holding its guard. On any error this executor reaps its owned group.
type CriterionStartControl func(start func() (procctl.Spawned, error), release func() error) error

// The supervisor keeps a stable command and session identity even when the
// material shell execs another program. Raw commands travel through a private
// environment value, removed before the material shell starts, never its argv.
// A caught TERM keeps the leader alive while its child is being terminated.
const criterionSupervisor = `trap ':' TERM
IFS= read -r ready || exit 125
[ "$ready" = go ] || exit 125
parley_command=$PARLEY_CAPTURED_COMMAND
unset PARLEY_CAPTURED_COMMAND
sh -c "$parley_command" </dev/null &
child=$!
wait "$child"
code=$?
while kill -0 "$child" 2>/dev/null; do
    wait "$child"
    code=$?
done
exit "$code"`

func RunCriterionControlled(ctx context.Context, root, name, command, executor string, control CriterionStartControl) CriterionExecution {
	start := time.Now()
	cmd := exec.CommandContext(ctx, "sh", "-c", command)
	var release, gate *os.File
	var controlErr error
	if control != nil {
		cmd = exec.CommandContext(ctx, "sh", "-c", criterionSupervisor)
		cmd.Env = append(os.Environ(), "PARLEY_CAPTURED_COMMAND="+command)
		gate, release, controlErr = os.Pipe()
		if controlErr == nil {
			defer gate.Close()
			defer release.Close()
			cmd.Stdin = gate
		}
	}
	cmd.Dir = root
	procctl.SetNewProcessGroup(cmd)
	cmd.Cancel = func() error {
		if cmd.Process == nil {
			return os.ErrProcessDone
		}
		// exec publishes Process before starting its context watcher. Setsid
		// makes this PID the group ID, even if the shell has already exited.
		// No post-Start identity assignment may race with cancellation.
		return procctl.KillGroup(procctl.Spawned{PID: cmd.Process.Pid, PGID: cmd.Process.Pid})
	}
	cmd.WaitDelay = waitDelay
	hasher := sha256.New()
	capbuf := &cappedWriter{max: captureMax}
	sink := io.MultiWriter(hasher, capbuf)
	cmd.Stdout = sink
	cmd.Stderr = sink
	runErr := controlErr
	if runErr == nil {
		if control == nil {
			runErr = cmd.Start()
		} else {
			started, released := false, false
			runErr = control(func() (procctl.Spawned, error) {
				if started {
					return procctl.Spawned{}, errors.New("criterion supervisor already started")
				}
				started = true
				if err := cmd.Start(); err != nil {
					return procctl.Spawned{}, err
				}
				sp := procctl.Capture(cmd, "captured-criterion")
				if ok, reason := procctl.Attributed(sp); !ok {
					return sp, fmt.Errorf("criterion supervisor identity unavailable: %s", reason)
				}
				return sp, nil
			}, func() error {
				if cmd.Process == nil || released {
					return errors.New("criterion release requires one started supervisor")
				}
				if err := ctx.Err(); err != nil {
					return err
				}
				released = true
				_, err := release.WriteString("go\n")
				return err
			})
			if runErr == nil && !released {
				runErr = errors.New("criterion control did not release its supervisor")
			}
			_ = release.Close()
		}
	}
	if cmd.Process != nil {
		if runErr != nil {
			_ = procctl.KillGroup(procctl.Spawned{PID: cmd.Process.Pid, PGID: cmd.Process.Pid})
			_ = cmd.Wait()
		} else {
			runErr = cmd.Wait()
		}
	}
	out := capbuf.buf.Bytes()

	ce := CommandEvidence{
		Command:        ScrubAndTruncate(command),  // safe representation only; raw text is never persisted
		CommandSHA256:  sha256Hex([]byte(command)), // hash binds the EXACT raw input
		OutputSHA256:   hex.EncodeToString(hasher.Sum(nil)),
		ExitCode:       exitCodeOf(runErr),
		DurationMillis: time.Since(start).Milliseconds(),
		Format:         FormatShell,
		ExecutedCases:  -1,
		FailedCases:    -1,
		SkippedCases:   -1,
		FailedPackages: -1,
		Diagnostics:    ScrubAndTruncate(string(out)),
	}
	if capbuf.overflow {
		ce.Diagnostics = "output exceeded the evidence capture bound — partial output is not a semantic basis\n" + ce.Diagnostics
	}
	// A non-ExitError run failure with zero captured output would otherwise
	// persist an unexplained exit_code of -1 with empty diagnostics: the run
	// error's text names the failing branch (attribution refusal, artifact
	// write, control refusal, start failure) and was discarded. Persist the
	// scrubbed, bounded reason so retained evidence identifies the branch.
	// ExitError keeps empty diagnostics — the exit code already tells that
	// story. OutputSHA256 still hashes only the stream the command emitted
	// (here: the empty string); the label marks this text as executor-side,
	// never as command output.
	if runErr != nil && len(out) == 0 {
		var exitErr *exec.ExitError
		if !errors.As(runErr, &exitErr) {
			ce.Diagnostics = ScrubAndTruncate("run error (no command output): " + runErr.Error())
		}
	}
	envelopeInvalid := ""
	goTestInvalid := ""
	if env, present, envErr := ParseEnvelope(string(out)); present {
		// An explicit envelope line is authoritative for case semantics: never
		// fall through to the unrelated test2json parser, and never recover an
		// earlier pass when the envelope is malformed — fail closed.
		ce.Format = FormatEnvelope
		if envErr != nil {
			envelopeInvalid = envErr.Error()
			ce.Diagnostics = "invalid evidence envelope (fails closed): " + envErr.Error() + "\n" + ce.Diagnostics
		} else {
			ce.ExecutedCases = *env.ExecutedCases
			ce.FailedCases = *env.FailedCases
		}
	} else if execCount, failCount, skipCount, pkgFailCount, recognized, gerr := parseGoTestJSONEvidence(out); recognized {
		// A recognized test2json stream is authoritative for case semantics:
		// when it is malformed or carries duplicate/case-aliased semantic
		// fields there is NO fall-through to an opaque shell pass — the record
		// keeps the gotest-json format with unknown counts and fails closed,
		// mirroring the envelope path.
		ce.Format = FormatGoTestJSON
		if gerr != nil {
			goTestInvalid = gerr.Error()
			ce.Diagnostics = "invalid test2json event stream (fails closed): " + gerr.Error() + "\n" + ce.Diagnostics
		} else {
			ce.ExecutedCases = execCount
			ce.FailedCases = failCount
			ce.SkippedCases = skipCount
			ce.FailedPackages = pkgFailCount
			if pkgFailCount > 0 {
				ce.Diagnostics = fmt.Sprintf("test2json stream reports %d package-level failure event(s) (fail/build-fail — a failed build runs no test cases; a masked exit code does not clear this)\n", pkgFailCount) + ce.Diagnostics
			}
		}
	}

	status := StatusPass
	switch {
	case runErr != nil:
		status = StatusFail
	case capbuf.overflow:
		status = StatusFail // unbounded output was truncated: not fully captured evidence
	case envelopeInvalid != "":
		status = StatusFail
	case goTestInvalid != "":
		status = StatusFail
	case ce.Format == FormatGoTestJSON && ce.FailedPackages > 0:
		status = StatusFail // structured evidence of a package/build failure overrides a masked exit 0
	case ce.Format != FormatShell && ce.FailedCases > 0:
		status = StatusFail // structured evidence of failing cases overrides exit 0
	case ce.Format != FormatShell && ce.ExecutedCases == 0 && ce.SkippedCases > 0:
		status = StatusSkipped // structured proof that cases were skipped, none executed
	case ce.Format != FormatShell && ce.ExecutedCases == 0:
		status = StatusNotRun // structured proof that nothing executed
	}
	complete := !capbuf.overflow && envelopeInvalid == "" && goTestInvalid == "" && ctx.Err() == nil
	// POSIX wait encodes a signalled material child as 128+signal. The
	// supervisor cannot distinguish that from an explicit high exit code;
	// neither is complete evidence of a patch-induced test regression.
	if control != nil && ce.ExitCode >= 128 {
		complete = false
	}
	if runErr != nil {
		var exitErr *exec.ExitError
		complete = complete && errors.As(runErr, &exitErr) && exitErr.ExitCode() >= 0
	}
	return CriterionExecution{Complete: complete, Record: CriterionRecord{
		Name:       name,
		Status:     status,
		Command:    ce,
		Provenance: Provenance{Executor: executor},
	}}
}

// cappedWriter retains at most max bytes of a stream while consuming it all,
// so a flooding command cannot exhaust memory. overflow records truncation.
type cappedWriter struct {
	buf      bytes.Buffer
	max      int
	overflow bool
}

func (w *cappedWriter) Write(p []byte) (int, error) {
	// Overflow is decided from the PRE-WRITE remaining capacity: a chunk that
	// fits exactly (alone or as the last of several) is retained in full and
	// is not truncation; only a chunk larger than what is left overflows.
	room := w.max - w.buf.Len()
	if len(p) > room {
		w.overflow = true
	}
	if room > 0 {
		n := len(p)
		if n > room {
			n = room
		}
		w.buf.Write(p[:n])
	}
	return len(p), nil
}

// testEvent is one test2json output line. Action and Test are the SEMANTIC
// fields the counters rely on; the strict parser below rejects a line where
// either is duplicated, case-aliased, or not a string.
type testEvent struct {
	Action string `json:"Action"`
	Test   string `json:"Test"`
}

// ParseGoTestJSON parses `go test -json` output. It returns
// (executed, failed, skipped, recognized): recognized is false when the stream
// does not look like test2json at all, so opaque text is never silently
// counted. Only terminal test-level actions count: pass/fail mark executed
// cases; skip is reported separately — an all-skip run proves no executed
// cases and is typed StatusSkipped, which cannot close.
//
// A stream that IS recognized as test2json but is malformed, carries
// duplicate/case-aliased semantic fields, or reports a PACKAGE-LEVEL failure
// (a fail or build-fail event with no Test field — a failed build runs no
// test cases, so such a failure is invisible to the case counts this
// signature returns) is rejected as a pass basis: this wrapper reports it as
// recognized with all counts zeroed (no usable evidence, never a pass).
// Callers that must distinguish those cases from clean recognition — and fail
// closed with the reason and the exact counts instead of degrading — use
// parseGoTestJSONEvidence; RunCriterion does.
func ParseGoTestJSON(out []byte) (executed, failed, skipped int, recognized bool) {
	executed, failed, skipped, failedPackages, recognized, err := parseGoTestJSONEvidence(out)
	if err != nil || failedPackages > 0 {
		return 0, 0, 0, recognized
	}
	return executed, failed, skipped, recognized
}

// parseGoTestJSONEvidence is the strict test2json parser behind RunCriterion.
// recognized=true means the stream claimed test2json structure; err is non-nil
// when that claim cannot be trusted — a line that decodes as an event object
// with a duplicated, case-aliased or non-string Action/Test field, a `{` line
// that is not one complete JSON object, or trailing content after an event.
// recognized=true with err != nil MUST fail closed: it is never a basis for
// counts and never falls back to opaque shell semantics. Plain non-`{` noise
// lines (build-failure text, summaries) and well-formed objects without an
// Action claim remain noise and are skipped, as before.
//
// Package-level events (no Test field) are never test cases and never touch
// the exact test-case counts: start/run/output/pass/skip carry no failure
// semantics, but a package-level fail or build-fail event is structured proof
// the PACKAGE failed at package scope (a failed build runs no test cases; a
// TestMain/panic failure may run none of them). Each such EVENT increments
// failedPackages — one failed build normally contributes two (its build-fail
// plus its package-level fail), so failedPackages is a failure SIGNAL to be
// read as >0, not a distinct-package tally. failedPackages == 0 guarantees the
// stream reports no package-scope failure.
func parseGoTestJSONEvidence(out []byte) (executed, failed, skipped, failedPackages int, recognized bool, err error) {
	sawEvent := false
	for _, line := range bytes.Split(out, []byte("\n")) {
		line = bytes.TrimSpace(line)
		if len(line) == 0 || line[0] != '{' {
			continue
		}
		ev, isEvent, lerr := parseTestEventStrict(line)
		if lerr != nil {
			return 0, 0, 0, 0, true, lerr
		}
		if !isEvent {
			continue
		}
		sawEvent = true
		if ev.Test == "" {
			if ev.Action == "fail" || ev.Action == "build-fail" {
				failedPackages++
			}
			continue // package-level event, not a test case
		}
		switch ev.Action {
		case "pass":
			executed++
		case "fail":
			executed++
			failed++
		case "skip":
			skipped++
		}
	}
	return executed, failed, skipped, failedPackages, sawEvent, nil
}

// parseTestEventStrict decodes one trimmed `{`-prefixed test2json output line
// with a token walk so that semantics a permissive json.Unmarshal would hide
// are rejected (same discipline as parseEnvelopeStrict): duplicate Action or
// Test fields (encoding/json's last-wins collapse would decode
// {"Action":"fail","Action":"pass"} as a pass — that is NEVER accepted as
// evidence), case-aliased semantic fields (encoding/json matches keys
// case-insensitively — "action" would silently set Action), non-string
// semantic values, and trailing content after the closing brace. Other real
// test2json fields (Time, Package, Output, Elapsed) and unknown future fields
// are skipped without inspection, including nested values.
//
// isEvent is false (with err nil) for a well-formed object that makes no
// Action claim — noise as far as case semantics are concerned. err non-nil
// means the line purported to carry event structure but cannot be trusted.
func parseTestEventStrict(line []byte) (ev testEvent, isEvent bool, err error) {
	dec := json.NewDecoder(bytes.NewReader(line))
	tok, err := dec.Token()
	if err != nil {
		return ev, false, fmt.Errorf("malformed test2json event line: %v", err)
	}
	if d, ok := tok.(json.Delim); !ok || d != '{' {
		return ev, false, fmt.Errorf("test2json event line is not a single JSON object")
	}
	seenAction, seenTest := false, false
	for dec.More() {
		kt, err := dec.Token()
		if err != nil {
			return ev, false, fmt.Errorf("malformed test2json event line: %v", err)
		}
		key, ok := kt.(string)
		if !ok {
			return ev, false, fmt.Errorf("malformed test2json event line: non-string field name")
		}
		semantic := key == "Action" || key == "Test"
		if !semantic {
			if strings.EqualFold(key, "Action") || strings.EqualFold(key, "Test") {
				return ev, false, fmt.Errorf("case-aliased test2json field %q — exact field names are required", key)
			}
			if err := skipJSONValue(dec); err != nil {
				return ev, false, fmt.Errorf("malformed test2json event line: %v", err)
			}
			continue
		}
		if (key == "Action" && seenAction) || (key == "Test" && seenTest) {
			return ev, false, fmt.Errorf("duplicate test2json field %q — last-wins collapse is not evidence", key)
		}
		if key == "Action" {
			seenAction = true
		} else {
			seenTest = true
		}
		vt, err := dec.Token()
		if err != nil {
			return ev, false, fmt.Errorf("malformed test2json event line: %v", err)
		}
		s, ok := vt.(string)
		if !ok {
			return ev, false, fmt.Errorf("test2json field %q must be a string — non-string semantic values are rejected", key)
		}
		if key == "Action" {
			ev.Action = s
		} else {
			ev.Test = s
		}
	}
	if _, err := dec.Token(); err != nil { // the closing '}'
		return ev, false, fmt.Errorf("malformed test2json event line: %v", err)
	}
	if _, err := dec.Token(); err != io.EOF {
		if err == nil {
			return ev, false, fmt.Errorf("trailing content after the test2json event object")
		}
		return ev, false, fmt.Errorf("malformed trailing content after the test2json event: %v", err)
	}
	return ev, ev.Action != "", nil
}

// skipJSONValue consumes one complete JSON value from dec, including nested
// objects and arrays. Used for event fields the counters do not rely on.
func skipJSONValue(dec *json.Decoder) error {
	tok, err := dec.Token()
	if err != nil {
		return err
	}
	if d, ok := tok.(json.Delim); ok && (d == '{' || d == '[') {
		depth := 1
		for depth > 0 {
			t, err := dec.Token()
			if err != nil {
				return err
			}
			if dd, ok := t.(json.Delim); ok {
				switch dd {
				case '{', '[':
					depth++
				case '}', ']':
					depth--
				}
			}
		}
	}
	return nil
}

// ParseEnvelope extracts the explicit evidence envelope, failing closed.
// Returns (envelope, present, err): present is true when ANY PARLEY-EVIDENCE
// line exists, and err is non-nil whenever that envelope cannot be trusted —
// malformed JSON, a non-object payload, duplicate fields (encoding/json's
// last-wins collapse is NEVER accepted as evidence), case-aliased or unknown
// fields, null/non-integer counts, trailing content after the object, negative
// counts, failed>executed, or duplicate envelope lines. A caller MUST NOT
// recover an earlier valid line when the last one is malformed, and MUST NOT
// fall through to another output parser when present is true.
func ParseEnvelope(out string) (Envelope, bool, error) {
	var lines []string
	for _, line := range strings.Split(out, "\n") {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, envelopePrefix) {
			lines = append(lines, strings.TrimPrefix(line, envelopePrefix))
		}
	}
	if len(lines) == 0 {
		return Envelope{}, false, nil
	}
	if len(lines) > 1 {
		return Envelope{}, true, fmt.Errorf("duplicate evidence envelope lines (%d) — ambiguous claim", len(lines))
	}
	e, err := parseEnvelopeStrict(lines[0])
	if err != nil {
		return Envelope{}, true, err
	}
	if e.ExecutedCases == nil || e.FailedCases == nil {
		return Envelope{}, true, fmt.Errorf("partial evidence envelope: executed_cases and failed_cases are both required")
	}
	if *e.ExecutedCases < 0 || *e.FailedCases < 0 {
		return Envelope{}, true, fmt.Errorf("negative case counts in evidence envelope")
	}
	if *e.FailedCases > *e.ExecutedCases {
		return Envelope{}, true, fmt.Errorf("conflicting evidence envelope: %d failed > %d executed", *e.FailedCases, *e.ExecutedCases)
	}
	return e, true, nil
}

// parseEnvelopeStrict decodes one envelope body with a token walk so that
// semantics a permissive json.Unmarshal would hide are rejected: duplicate
// fields (a zero executed count followed by a one must NOT decode as one),
// case-aliased fields (encoding/json matches keys case-insensitively —
// "Executed_Cases" would silently set ExecutedCases), unknown fields, null or
// non-integer values, and any trailing content after the closing brace.
func parseEnvelopeStrict(line string) (Envelope, error) {
	dec := json.NewDecoder(strings.NewReader(line))
	dec.UseNumber()
	tok, err := dec.Token()
	if err != nil {
		return Envelope{}, fmt.Errorf("malformed evidence envelope: %v", err)
	}
	if d, ok := tok.(json.Delim); !ok || d != '{' {
		return Envelope{}, fmt.Errorf("evidence envelope must be a single JSON object")
	}
	var e Envelope
	for dec.More() {
		kt, err := dec.Token()
		if err != nil {
			return Envelope{}, fmt.Errorf("malformed evidence envelope: %v", err)
		}
		key, ok := kt.(string)
		if !ok {
			return Envelope{}, fmt.Errorf("malformed evidence envelope: non-string field name")
		}
		var target **int
		switch key {
		case "executed_cases":
			target = &e.ExecutedCases
		case "failed_cases":
			target = &e.FailedCases
		default:
			if strings.EqualFold(key, "executed_cases") || strings.EqualFold(key, "failed_cases") {
				return Envelope{}, fmt.Errorf("aliased evidence envelope field %q — exact lowercase field names are required", key)
			}
			return Envelope{}, fmt.Errorf("unknown evidence envelope field %q", key)
		}
		if *target != nil {
			return Envelope{}, fmt.Errorf("duplicate evidence envelope field %q", key)
		}
		vt, err := dec.Token()
		if err != nil {
			return Envelope{}, fmt.Errorf("malformed evidence envelope: %v", err)
		}
		num, ok := vt.(json.Number)
		if !ok {
			return Envelope{}, fmt.Errorf("evidence envelope field %q must be a JSON integer — null, booleans, strings and nested values are rejected", key)
		}
		n, err := num.Int64()
		if err != nil {
			return Envelope{}, fmt.Errorf("evidence envelope field %q must be an integer: %v", key, err)
		}
		v := int(n)
		*target = &v
	}
	if _, err := dec.Token(); err != nil { // the closing '}'
		return Envelope{}, fmt.Errorf("malformed evidence envelope: %v", err)
	}
	if _, err := dec.Token(); err != io.EOF {
		if err == nil {
			return Envelope{}, fmt.Errorf("trailing content after the evidence envelope object")
		}
		return Envelope{}, fmt.Errorf("malformed trailing content after the evidence envelope: %v", err)
	}
	return e, nil
}

// secretPatterns redact credential-shaped tokens from recorded output before
// it is persisted. Same layering as the existing driver evidence writer:
// labeled key/value pairs AND standalone provider token shapes.
var secretPatterns = []*regexp.Regexp{
	regexp.MustCompile(`(?i)(authorization\s*[:=]\s*)(bearer\s+)?\S+`),
	regexp.MustCompile(`(?i)\bbearer\s+[A-Za-z0-9._\-]+`),
	regexp.MustCompile(`(?i)(token|secret|password|passwd|api[_-]?key|access[_-]?key|private[_-]?key)(\s*[:=]\s*)\S+`),
	regexp.MustCompile(`\bsk-[A-Za-z0-9_\-]{16,}`),
	regexp.MustCompile(`\bgh[pousr]_[A-Za-z0-9]{20,}`),
	regexp.MustCompile(`\bxox[baprs]-[A-Za-z0-9\-]{10,}`),
	regexp.MustCompile(`\bAKIA[0-9A-Z]{16}\b`),
	regexp.MustCompile(`\beyJ[A-Za-z0-9_\-]{10,}\.[A-Za-z0-9_\-]{10,}\.[A-Za-z0-9_\-]{5,}`), // JWT
}

// ScrubAndTruncate redacts credential-shaped tokens and caps output to a
// bounded tail (last evidenceMaxLines lines, then evidenceMaxBytes bytes).
func ScrubAndTruncate(s string) string {
	s = secretPatterns[0].ReplaceAllString(s, "$1«redacted»")
	s = secretPatterns[2].ReplaceAllString(s, "$1$2«redacted»")
	for i, re := range secretPatterns {
		if i == 0 || i == 2 {
			continue
		}
		s = re.ReplaceAllString(s, "«redacted»")
	}
	lines := strings.Split(strings.TrimRight(s, "\n"), "\n")
	if len(lines) > evidenceMaxLines {
		lines = lines[len(lines)-evidenceMaxLines:]
	}
	out := strings.Join(lines, "\n")
	if len(out) > evidenceMaxBytes {
		out = "…" + out[len(out)-evidenceMaxBytes:]
	}
	return out
}

func exitCodeOf(err error) int {
	if err == nil {
		return 0
	}
	if ee, ok := err.(*exec.ExitError); ok {
		return ee.ExitCode()
	}
	return -1
}

func sha256Hex(b []byte) string {
	sum := sha256.Sum256(b)
	return hex.EncodeToString(sum[:])
}
