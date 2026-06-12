package app

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"parley-deck-cli/internal/agents"
	"parley-deck-cli/internal/fsutil"
	"parley-deck-cli/internal/protocol"
	"parley-deck-cli/internal/runner"
)

// parley consult (runner-hardening-kindly D10): an advisory cross-agent
// question with repo context. The answer is captured into a durable artifact
// under parley-deck/consults/ with a provenance line in index.jsonl. Consult
// artifacts are advisory and non-canonical — never quorum evidence (see the
// protocol's Consults note).

func runConsult(ctx context.Context, args []string, stdout, stderr io.Writer) int {
	fs := flag.NewFlagSet("consult", flag.ContinueOnError)
	fs.SetOutput(stderr)
	dir := fs.String("dir", ".", "workspace root")
	timeout := fs.Duration("timeout", 0, "per-agent timeout (default: the agent's timeout_ms)")
	if err := fs.Parse(args); err != nil {
		return 2
	}
	rest := fs.Args()
	if len(rest) < 1 {
		fmt.Fprintln(stderr, "usage: parley consult [--dir DIR] [--timeout D] AGENT [QUESTION] (question falls back to stdin)")
		return 2
	}
	agentID := rest[0]
	question := strings.TrimSpace(strings.Join(rest[1:], " "))
	if question == "" {
		data, err := io.ReadAll(os.Stdin)
		if err != nil {
			fmt.Fprintf(stderr, "consult: read question from stdin: %v\n", err)
			return 1
		}
		question = strings.TrimSpace(string(data))
	}
	if question == "" {
		fmt.Fprintln(stderr, "consult: a question is required (argument or stdin)")
		return 2
	}

	root, err := filepath.Abs(*dir)
	if err != nil {
		fmt.Fprintf(stderr, "consult: %v\n", err)
		return 1
	}
	discovered, err := discoverConfigured(ctx, root)
	if err != nil {
		fmt.Fprintf(stderr, "consult: %v\n", err)
		return 1
	}
	var agent *agents.Discovery
	for i := range discovered {
		if discovered[i].ID == agentID {
			agent = &discovered[i]
			break
		}
	}
	if agent == nil || !agent.Found {
		fmt.Fprintf(stderr, "consult: agent %q is not installed/configured (try: parley agents list)\n", agentID)
		return 1
	}

	consultsDir := filepath.Join(root, protocol.DeckDir, "consults")
	created := time.Now().UTC()
	slug := consultSlug(question)
	name := fmt.Sprintf("%s-%s-%s", created.Format("20060102T150405Z"), agentID, slug)
	artifactPath := filepath.Join(consultsDir, name+".md")
	logsDir := filepath.Join(consultsDir, "logs")
	if err := fsutil.MkdirAllResilient(logsDir, 0o755); err != nil {
		fmt.Fprintf(stderr, "consult: %v\n", err)
		return 1
	}
	stdoutLog := filepath.Join(logsDir, name+".stdout.log")
	stderrLog := filepath.Join(logsDir, name+".stderr.log")

	fmt.Fprintf(stderr, "Consulting %s: %q\n", agentID, truncateQuestion(question, 100))
	res := runner.RunConsult(ctx, runner.ConsultOptions{
		Root:       root,
		Agent:      *agent,
		Question:   question,
		Timeout:    *timeout,
		StdoutPath: stdoutLog,
		StderrPath: stderrLog,
		Progress:   stderr,
	})

	outcome := "answered"
	exitCode := 0
	if res.ExitError != "" {
		outcome = "failed"
		exitCode = 1
	}
	// The artifact is written even on failure (with the diagnostics) so a slow,
	// expensive consult is never silently lost.
	var body strings.Builder
	fmt.Fprintf(&body, `---
artifact: consult
agent: %s
model: %s
created: %s
question_slug: %s
question: %q
workspace_root: %s
timeout_ms: %d
exit_code: %d
agent_exit: %d
stdout_log: %s
stderr_log: %s
quorum: false
---

`, agentID, valueOr(agent.Model, "cli-default"), created.Format(time.RFC3339), slug, question,
		root, timeout.Milliseconds(), exitCode, res.AgentExit,
		relOrSelf(root, stdoutLog), relOrSelf(root, stderrLog))
	if res.ExitError != "" {
		fmt.Fprintf(&body, "## Consult failed\n\n- error: %s\n- class: %s\n- hint: %s\n", res.ExitError, res.FailureClass, res.RecoveryHint)
	} else {
		body.WriteString(strings.TrimSpace(res.Answer))
		body.WriteString("\n")
	}
	if err := os.WriteFile(artifactPath, []byte(body.String()), 0o644); err != nil {
		fmt.Fprintf(stderr, "consult: write artifact: %v\n", err)
		return 1
	}

	ledger, _ := json.Marshal(map[string]any{
		"ts": created.Format(time.RFC3339), "agent": agentID, "model": valueOr(agent.Model, "cli-default"),
		"question": question, "slug": slug, "path": relOrSelf(root, artifactPath),
		"outcome": outcome, "exit_code": exitCode, "agent_exit": res.AgentExit,
		"duration_ms": res.Duration.Milliseconds(),
	})
	if err := fsutil.AppendLine(filepath.Join(consultsDir, "index.jsonl"), ledger); err != nil {
		fmt.Fprintf(stderr, "consult: warning: index append failed: %v\n", err)
	}

	if res.ExitError != "" {
		fmt.Fprintf(stderr, "consult: %s failed: %s (class: %s)\n", agentID, res.ExitError, res.FailureClass)
		if res.RecoveryHint != "" {
			fmt.Fprintf(stderr, "consult: hint: %s\n", res.RecoveryHint)
		}
		fmt.Fprintf(stderr, "consult: diagnostics saved to %s\n", artifactPath)
		return 1
	}
	// The answer goes to stdout (redirectable); provenance to stderr.
	fmt.Fprintln(stdout, strings.TrimSpace(res.Answer))
	if res.AgentExit != 0 {
		fmt.Fprintf(stderr, "consult: note: %s exited %d after producing the answer\n", agentID, res.AgentExit)
	}
	fmt.Fprintf(stderr, "consult: saved to %s\n", relOrSelf(root, artifactPath))
	return 0
}

// runConsults implements `parley consults list`.
func runConsults(args []string, stdout, stderr io.Writer) int {
	fs := flag.NewFlagSet("consults", flag.ContinueOnError)
	fs.SetOutput(stderr)
	dir := fs.String("dir", ".", "workspace root")
	if err := fs.Parse(args); err != nil {
		return 2
	}
	if rest := fs.Args(); len(rest) == 0 || rest[0] != "list" {
		fmt.Fprintln(stderr, "usage: parley consults list [--dir DIR]")
		return 2
	}
	root, err := filepath.Abs(*dir)
	if err != nil {
		fmt.Fprintf(stderr, "consults: %v\n", err)
		return 1
	}
	indexPath := filepath.Join(root, protocol.DeckDir, "consults", "index.jsonl")
	data, err := os.ReadFile(indexPath)
	if err != nil {
		if os.IsNotExist(err) {
			fmt.Fprintln(stdout, "no consults yet — run: parley consult AGENT \"question\"")
			return 0
		}
		fmt.Fprintf(stderr, "consults: %v\n", err)
		return 1
	}
	type row struct{ ts, agent, question, file string }
	var rows []row
	for _, line := range strings.Split(string(data), "\n") {
		if strings.TrimSpace(line) == "" {
			continue
		}
		var entry map[string]any
		if json.Unmarshal([]byte(line), &entry) != nil {
			continue
		}
		ts, _ := entry["ts"].(string)
		agent, _ := entry["agent"].(string)
		question, _ := entry["question"].(string)
		path, _ := entry["path"].(string)
		rows = append(rows, row{ts: ts, agent: agent, question: question, file: filepath.Base(path)})
	}
	sort.Slice(rows, func(i, j int) bool { return rows[i].ts < rows[j].ts })
	// FILE shows just the filename — the shared prefix would wrap at 100 cols
	// (consensus D10).
	fmt.Fprintf(stdout, "%-20s %-8s %-38s %s\n", "DATE", "AGENT", "QUESTION", "FILE")
	for _, r := range rows {
		when := r.ts
		if t, err := time.Parse(time.RFC3339, r.ts); err == nil {
			when = t.Local().Format("2006-01-02 15:04:05")
		}
		fmt.Fprintf(stdout, "%-20s %-8s %-38s %s\n", when, r.agent, truncateQuestion(r.question, 38), r.file)
	}
	return 0
}

// consultSlug derives a filename-safe slug from the question's leading words.
func consultSlug(question string) string {
	var b strings.Builder
	words := 0
	for _, r := range strings.ToLower(question) {
		switch {
		case r >= 'a' && r <= 'z', r >= '0' && r <= '9':
			b.WriteRune(r)
		case b.Len() > 0 && !strings.HasSuffix(b.String(), "-"):
			b.WriteByte('-')
			words++
		}
		if words >= 6 || b.Len() >= 48 {
			break
		}
	}
	slug := strings.Trim(b.String(), "-")
	if slug == "" {
		return "question"
	}
	return slug
}

func truncateQuestion(q string, n int) string {
	q = strings.Join(strings.Fields(q), " ")
	if len(q) <= n {
		return q
	}
	return q[:n-1] + "…"
}

func relOrSelf(root, path string) string {
	if rel, err := filepath.Rel(root, path); err == nil && !strings.HasPrefix(rel, "..") {
		return rel
	}
	return path
}
