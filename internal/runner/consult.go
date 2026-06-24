package runner

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"parley-deck-cli/internal/agents"
	"parley-deck-cli/internal/protocol"
)

// Consult (runner-hardening-kindly D10): a lightweight advisory cross-agent
// question with repo context — no run, no protocol state. The facilitator
// captures stdout into the durable artifact, so the consulted agent stays
// read-only-ish. Consults run under the same supervision (watchdogs) and
// failure classification as protocol agents.

type ConsultOptions struct {
	Root       string
	Agent      agents.Discovery
	Question   string
	Timeout    time.Duration
	StdoutPath string
	StderrPath string
	// Prompt, when set, is used verbatim instead of BuildConsultPrompt(Question). The
	// goal-done gate (LE-7) sets it to a verdict-style prompt; ordinary consults leave
	// it empty and get the advisory framing.
	Prompt string
	// Progress receives content-free liveness lines (heartbeats, watchdog
	// notices) — callers point it at stderr so stdout stays redirectable.
	Progress io.Writer
}

type ConsultResult struct {
	Answer       string
	ExitError    string
	FailureClass string
	RecoveryHint string
	AgentExit    int
	Duration     time.Duration
	// EffectiveTimeout is the timeout RunConsult actually enforced (the flag
	// or the agent's timeout_ms default) — recorded in the artifact provenance.
	EffectiveTimeout time.Duration
	// SessionID is the consulted CLI's session id when discoverable; one-shot
	// CLIs expose none today, so it is usually empty (written regardless).
	SessionID string
}

// BuildConsultPrompt frames the advisory contract (adapted from kindly's
// consult prompt): recommendation first, grounded in what the agent read,
// explicitly NOT a pass/fail gate.
func BuildConsultPrompt(agent agents.Discovery, question string) string {
	return fmt.Sprintf(`You are %s, consulted as a fresh outside advisor with read-only access to this repository (an advisory consult, NOT a review gate).

Rules:
- Read whatever in the repository bears on the question, then answer candidly and practically.
- Lead with your recommendation; surface risks, alternatives, simplifications, and concrete next steps, grounding each point in what you read.
- Do NOT create, modify, or delete any files. Print your entire answer to stdout as markdown.
- No severity tables, no pass/fail verdict — this is advisory thinking.

Question:
%s
`, agent.ID, question)
}

// BuildGoalCheckPrompt frames the LE-7 goal-done gate: a fresh agent checks the FINAL.md
// observable acceptance criteria against the implementation and ends with a structured
// verdict line the driver parses. Unlike the advisory consult, this one DOES ask for a
// pass/fail verdict — but it is still read-only and runs once, before close.
func BuildGoalCheckPrompt(agent agents.Discovery, idea protocol.IdeaStatus) string {
	final := filepath.Join(idea.Path, "FINAL.md")
	impl := filepath.Join(idea.Path, "IMPLEMENTATION.md")
	return fmt.Sprintf(`You are %s, the goal-done checker for Parley Deck idea %s (read-only).

Task: decide whether the implementation meets the OBSERVABLE ACCEPTANCE CRITERIA in
%s, using %s and the actual repository state. Check each criterion concretely (run or
trace it); do not assume.

Output: a short rationale, then a FINAL line that is EXACTLY one of:
  GOAL-CHECK: PASS
  GOAL-CHECK: FAIL — <the criteria that are not met>
Print to stdout only; do not modify any file. Answer PASS only if every observable
criterion is met; if any is unmet or unverifiable, answer FAIL and name it.
`, agent.ID, idea.Slug, final, impl)
}

// RunConsult invokes the agent once (no retry) with supervision and captures
// stdout as the answer.
func RunConsult(ctx context.Context, opts ConsultOptions) ConsultResult {
	started := time.Now()
	prompt := BuildConsultPrompt(opts.Agent, opts.Question)
	if strings.TrimSpace(opts.Prompt) != "" {
		prompt = opts.Prompt
	}
	hardTimeout := timeoutForAgent(opts.Timeout, opts.Agent)
	result := ConsultResult{EffectiveTimeout: hardTimeout}
	cctx, cancel := context.WithTimeout(ctx, hardTimeout)
	defer cancel()

	progress := func(format string, args ...any) {
		if opts.Progress != nil {
			fmt.Fprintf(opts.Progress, "[%s] "+format+"\n", append([]any{time.Now().Format("15:04:05")}, args...)...)
		}
	}
	act := &activityTracker{}
	cfg := supervisionForAgent(opts.Agent, hardTimeout)
	hooks := supervisionHooks{
		onHeartbeat: func(snap activitySnapshot, elapsed time.Duration) {
			progress("%s thinking (elapsed: %ds, bytes: %d+%d)", opts.Agent.ID, int(elapsed.Seconds()), snap.StdoutBytes, snap.StderrBytes)
		},
		onWatchdog: func(kind string, snap activitySnapshot, elapsed time.Duration) {
			progress("%s %s after %ds — killing process tree", opts.Agent.ID, kind, int(elapsed.Seconds()))
		},
	}
	_, err := execAgentProcess(cctx, opts.Root, "consult", opts.Agent.ID, "consult:"+opts.Agent.ID, opts.Agent, prompt, opts.StdoutPath, opts.StderrPath, nil, act, cfg, hooks)
	result.Duration = time.Since(started)

	watchdog := ""
	switch {
	case errors.Is(err, errNoFirstOutput):
		watchdog = "no_first_output"
	case errors.Is(err, errStalled):
		watchdog = "stalled"
	}
	answer, readErr := os.ReadFile(opts.StdoutPath)
	if readErr == nil {
		result.Answer = string(answer)
	}
	if err != nil {
		// Artifact-wins, consult flavor: a non-empty answer overrides an
		// ordinary nonzero exit; watchdogs and the hard timeout always fail.
		var exitErr *exec.ExitError
		ordinaryExit := watchdog == "" && cctx.Err() == nil && errors.As(err, &exitErr)
		if ordinaryExit && len(result.Answer) > 0 {
			result.AgentExit = exitErr.ExitCode()
		} else {
			result.ExitError = err.Error()
			if cctx.Err() != nil {
				result.ExitError = cctx.Err().Error()
			}
			result.FailureClass, result.RecoveryHint = terminalFailureClass(watchdog, cctx.Err(), false, opts.StderrPath, opts.StdoutPath, result.ExitError)
		}
	}
	if result.ExitError == "" && len(result.Answer) == 0 {
		result.ExitError = "the agent produced no answer on stdout"
		result.FailureClass, result.RecoveryHint = classifyFailure(opts.StderrPath, opts.StdoutPath, result.ExitError)
	}
	return result
}
