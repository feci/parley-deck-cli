package runner

import (
	"os"
	"regexp"
	"strings"
)

// Failure classification (runner-hardening-kindly D5): a bounded, ordered,
// data-driven regex table over the agent's stderr/stdout tails and the exit
// error text. Classes and recovery hints were converged in the idea's
// round-01 (agy's taxonomy, seeded from kindly-agent.sh's detail/hint scans).

const failTailBytes = 4096 // bounded: only the last 4 KiB of each log is scanned

type failureRule struct {
	class string
	hint  string
	re    *regexp.Regexp
}

// Order matters: the first matching rule wins. Auth outranks generic 4xx noise;
// rate-limit outranks overloaded (429 vs 5xx).
var failureRules = []failureRule{
	{"rate-limit", "Wait for reset or switch provider keys/endpoints.",
		regexp.MustCompile(`(?i)rate[-_ ]?limit|rate_limit_error|usageLimitExceeded|usage limit|session limit|too many requests|(^|[^0-9])429([^0-9]|$)`)},
	{"auth", "Run the agent CLI's auth command (e.g. 'claude login', 'hermes auth') to refresh credentials.",
		regexp.MustCompile(`(?i)authentication[_ ](failed|error|required)|unauthorized|forbidden|permission_error|oauth[_ ](org|error|failed|not allowed)|oauth_org_not_allowed|api key not valid|invalid api key|(^|[^0-9])401([^0-9]|$)`)},
	{"billing", "Check your API account balance and credit card status.",
		regexp.MustCompile(`(?i)billing[_ ](error|failed|required)|payment|credit (error|exhausted|balance)|quota exceeded`)},
	{"overloaded", "Retry in a few minutes or choose a less busy model.",
		regexp.MustCompile(`(?i)overloaded|overloaded_error|serverOverloaded|server[_ ]error|internalServerError|timeout_error|request timed out|deadline exceeded|(status|http|code|response|error)[^0-9]{0,40}5[0-9][0-9]|5[0-9][0-9][^0-9]{0,40}(status|http|code|response|error)`)},
	{"model-not-found", "Check the model spelling and access permissions in your API settings.",
		regexp.MustCompile(`(?i)model[_ ]not[_ ]found|model not found|not_found_error|unknown model`)},
	{"context-window", "Reduce the prompt size or prune file attachments/logs from scope.",
		regexp.MustCompile(`(?i)contextWindowExceeded|context (length|window)|request_too_large|max_output_tokens|max output tokens|prompt is too long`)},
	{"sandbox", "Adjust the local sandbox configuration or run with lower restriction.",
		regexp.MustCompile(`(?i)sandboxError|seatbelt|operation not permitted.*sandbox|sandbox.*denied|cyberPolicy`)},
	{"budget", "Increase the session budget limit (e.g. raise spend caps in settings).",
		regexp.MustCompile(`(?i)Exceeded USD budget|exceeded[[:space:][:alpha:]]*budget|max-budget`)},
	{"invalid-request", "Verify the prompt structure and system constraints in config.",
		regexp.MustCompile(`(?i)invalid_request(_error)?|invalid request|badRequest|bad request|(^|[^0-9])400([^0-9]|$)`)},
}

// Watchdog/timeout classes are assigned by the supervisor before regex
// classification ever runs (consensus D5); their hints live here so every
// class/hint pair stays in one table.
var watchdogHints = map[string]string{
	"no_first_output": "Verify the agent executable is not blocking or waiting for stdin.",
	"stalled":         "Check the process tree; the agent emitted no output within the stall window.",
	"timeout":         "The hard per-agent timeout elapsed; raise timeout_ms or split the task.",
	"unknown":         "Check the agent's stdout/stderr log files for details.",
}

// classifyFailure scans the bounded tails of stderr and stdout plus the exit
// error text and returns (failure_class, recovery_hint). Unmatched input is
// ("unknown", generic hint).
func classifyFailure(stderrPath, stdoutPath, exitError string) (string, string) {
	var b strings.Builder
	b.WriteString(exitError)
	b.WriteByte('\n')
	b.Write(tailOfFile(stderrPath, failTailBytes))
	b.WriteByte('\n')
	b.Write(tailOfFile(stdoutPath, failTailBytes))
	corpus := b.String()
	for _, rule := range failureRules {
		if rule.re.MatchString(corpus) {
			return rule.class, rule.hint
		}
	}
	return "unknown", watchdogHints["unknown"]
}

// tailOfFile reads at most n bytes from the end of path; missing files are
// empty (classification stays best-effort).
func tailOfFile(path string, n int64) []byte {
	f, err := os.Open(path)
	if err != nil {
		return nil
	}
	defer f.Close()
	info, err := f.Stat()
	if err != nil {
		return nil
	}
	size := info.Size()
	if size > n {
		if _, err := f.Seek(size-n, 0); err != nil {
			return nil
		}
	}
	buf := make([]byte, min64(size, n))
	read, _ := f.Read(buf)
	return buf[:read]
}

func min64(a, b int64) int64 {
	if a < b {
		return a
	}
	return b
}
