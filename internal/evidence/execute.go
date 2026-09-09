package evidence

import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
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
	start := time.Now()
	cmd := exec.CommandContext(ctx, "sh", "-c", command)
	cmd.Dir = root
	procctl.SetNewProcessGroup(cmd)
	var spawned procctl.Spawned
	cmd.Cancel = func() error {
		if spawned.PID > 0 {
			return procctl.KillGroup(spawned)
		}
		return nil
	}
	cmd.WaitDelay = waitDelay
	hasher := sha256.New()
	capbuf := &cappedWriter{max: captureMax}
	sink := io.MultiWriter(hasher, capbuf)
	cmd.Stdout = sink
	cmd.Stderr = sink
	runErr := cmd.Start()
	if runErr == nil {
		spawned = procctl.Capture(cmd, "evidence-criterion:"+name)
		runErr = cmd.Wait()
	}
	out := capbuf.buf.Bytes()

	ce := CommandEvidence{
		Command:        ScrubAndTruncate(command), // safe representation only; raw text is never persisted
		CommandSHA256:  sha256Hex([]byte(command)), // hash binds the EXACT raw input
		OutputSHA256:   hex.EncodeToString(hasher.Sum(nil)),
		ExitCode:       exitCodeOf(runErr),
		DurationMillis: time.Since(start).Milliseconds(),
		Format:         FormatShell,
		ExecutedCases:  -1,
		FailedCases:    -1,
		SkippedCases:   -1,
		Diagnostics:    ScrubAndTruncate(string(out)),
	}
	if capbuf.overflow {
		ce.Diagnostics = "output exceeded the evidence capture bound — partial output is not a semantic basis\n" + ce.Diagnostics
	}
	envelopeInvalid := ""
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
	} else if execCount, failCount, skipCount, ok := ParseGoTestJSON(out); ok {
		ce.Format = FormatGoTestJSON
		ce.ExecutedCases = execCount
		ce.FailedCases = failCount
		ce.SkippedCases = skipCount
	}

	status := StatusPass
	switch {
	case runErr != nil:
		status = StatusFail
	case capbuf.overflow:
		status = StatusFail // unbounded output was truncated: not fully captured evidence
	case envelopeInvalid != "":
		status = StatusFail
	case ce.Format != FormatShell && ce.FailedCases > 0:
		status = StatusFail // structured evidence of failing cases overrides exit 0
	case ce.Format != FormatShell && ce.ExecutedCases == 0 && ce.SkippedCases > 0:
		status = StatusSkipped // structured proof that cases were skipped, none executed
	case ce.Format != FormatShell && ce.ExecutedCases == 0:
		status = StatusNotRun // structured proof that nothing executed
	}
	return CriterionRecord{
		Name:       name,
		Status:     status,
		Command:    ce,
		Provenance: Provenance{Executor: executor},
	}
}

// cappedWriter retains at most max bytes of a stream while consuming it all,
// so a flooding command cannot exhaust memory. overflow records truncation.
type cappedWriter struct {
	buf      bytes.Buffer
	max      int
	overflow bool
}

func (w *cappedWriter) Write(p []byte) (int, error) {
	if room := w.max - w.buf.Len(); room > 0 {
		n := len(p)
		if n > room {
			n = room
		}
		w.buf.Write(p[:n])
	}
	if len(p) > w.max-w.buf.Len() {
		w.overflow = true
	}
	return len(p), nil
}

// testEvent is one test2json output line.
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
func ParseGoTestJSON(out []byte) (executed, failed, skipped int, recognized bool) {
	sawEvent := false
	for _, line := range bytes.Split(out, []byte("\n")) {
		line = bytes.TrimSpace(line)
		if len(line) == 0 || line[0] != '{' {
			continue
		}
		var ev testEvent
		if err := json.Unmarshal(line, &ev); err != nil || ev.Action == "" {
			continue
		}
		sawEvent = true
		if ev.Test == "" {
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
	return executed, failed, skipped, sawEvent
}

// ParseEnvelope extracts the explicit evidence envelope, failing closed.
// Returns (envelope, present, err): present is true when ANY PARLEY-EVIDENCE
// line exists, and err is non-nil whenever that envelope cannot be trusted —
// malformed JSON, null/partial counts, negative counts, failed>executed, or
// duplicate envelope lines. A caller MUST NOT recover an earlier valid line
// when the last one is malformed, and MUST NOT fall through to another output
// parser when present is true.
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
	var e Envelope
	if err := json.Unmarshal([]byte(lines[0]), &e); err != nil {
		return Envelope{}, true, fmt.Errorf("malformed evidence envelope: %v", err)
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
