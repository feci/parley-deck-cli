package runner

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"parley-deck-cli/internal/agents"
	"parley-deck-cli/internal/protocol"
	"parley-deck-cli/internal/store"
)

type Options struct {
	Root       string
	RunID      string
	Idea       protocol.IdeaStatus
	Task       string
	Agents     []agents.Discovery
	Timeout    time.Duration
	Store      store.Store
	Round      int
	RoundLabel string
}

type Result struct {
	AgentID     string
	OutputPath  string
	StdoutPath  string
	StderrPath  string
	StartedAt   time.Time
	CompletedAt time.Time
	ExitError   string
	ArtifactOK  bool
	Skipped     bool
	SkipReason  string
	Warning     string
	Duration    time.Duration
}

type Handle struct {
	RunID  string
	RunDir string

	done    chan struct{}
	results []Result
	mu      sync.Mutex
}

func RunRoundOneAsync(ctx context.Context, opts Options) *Handle {
	handle := &Handle{
		RunID:  opts.RunID,
		RunDir: filepath.Join(opts.Root, protocol.DeckDir, "runs", opts.RunID),
		done:   make(chan struct{}),
	}
	go func() {
		defer close(handle.done)
		defer func() {
			if recovered := recover(); recovered != nil {
				now := time.Now().UTC()
				result := Result{
					AgentID:     "runner",
					CompletedAt: now,
					ExitError:   fmt.Sprintf("runner panic: %v", recovered),
				}
				_ = opts.Store.Append(store.Event{
					Time: now,
					Type: "run.failed",
					Data: map[string]any{"error": result.ExitError},
				})
				handle.setResults([]Result{result})
			}
		}()
		handle.setResults(RunRoundOne(ctx, opts))
	}()
	return handle
}

func (h *Handle) Done() <-chan struct{} {
	return h.done
}

func (h *Handle) Wait() []Result {
	<-h.done
	return h.Results()
}

func (h *Handle) Results() []Result {
	h.mu.Lock()
	defer h.mu.Unlock()
	results := make([]Result, len(h.results))
	copy(results, h.results)
	return results
}

func (h *Handle) setResults(results []Result) {
	h.mu.Lock()
	defer h.mu.Unlock()
	h.results = append([]Result(nil), results...)
}

func RunRoundOne(ctx context.Context, opts Options) []Result {
	if opts.Round == 0 {
		opts.Round = 1
	}
	if opts.RoundLabel == "" {
		opts.RoundLabel = "round-01"
	}

	selected := selectedAgents(opts.Idea.Participants, opts.Agents)
	results := make([]Result, len(selected))
	var wg sync.WaitGroup
	for i, agent := range selected {
		i, agent := i, agent
		wg.Add(1)
		go func() {
			defer wg.Done()
			results[i] = runAgent(ctx, opts, agent)
		}()
	}
	wg.Wait()

	eventType := "round.completed"
	okCount := 0
	for _, result := range results {
		if result.ExitError == "" && (result.ArtifactOK || result.Skipped) {
			okCount++
			continue
		}
		eventType = "round.incomplete"
	}
	if indexPath, err := writeRoundIndex(opts.Idea, opts.RoundLabel, results); err != nil {
		now := time.Now().UTC()
		warning := "round index write failed: " + err.Error()
		results = append(results, Result{
			AgentID:     "runner/index",
			OutputPath:  indexPath,
			CompletedAt: now,
			Warning:     warning,
		})
		_ = opts.Store.Append(store.Event{
			Time: now,
			Type: "round.index_failed",
			Data: map[string]any{
				"idea":     opts.Idea.Slug,
				"round":    opts.RoundLabel,
				"artifact": indexPath,
				"error":    err.Error(),
			},
		})
	} else {
		_ = opts.Store.Append(store.Event{
			Time: time.Now().UTC(),
			Type: "round.index_written",
			Data: map[string]any{
				"idea":     opts.Idea.Slug,
				"round":    opts.RoundLabel,
				"artifact": indexPath,
			},
		})
	}
	if err := opts.Store.Append(store.Event{
		Time: time.Now().UTC(),
		Type: eventType,
		Data: map[string]any{
			"idea":      opts.Idea.Slug,
			"round":     opts.RoundLabel,
			"completed": okCount,
			"total":     len(results),
		},
	}); err != nil {
		results = append(results, Result{
			AgentID:     "runner",
			CompletedAt: time.Now().UTC(),
			ExitError:   "round event append failed: " + err.Error(),
		})
	}
	return results
}

func selectedAgents(participants []string, discovered []agents.Discovery) []agents.Discovery {
	byID := make(map[string]agents.Discovery, len(discovered))
	for _, agent := range discovered {
		if agent.Found {
			byID[agent.ID] = agent
		}
	}

	var selected []agents.Discovery
	for _, participant := range participants {
		if agent, ok := byID[participant]; ok {
			selected = append(selected, agent)
		}
	}
	return selected
}

func runAgent(parent context.Context, opts Options, agent agents.Discovery) Result {
	now := time.Now().UTC()
	agentDir := filepath.Join(opts.Root, protocol.DeckDir, "runs", opts.RunID, "agents", agent.ID)
	stdoutPath := filepath.Join(agentDir, "stdout.log")
	stderrPath := filepath.Join(agentDir, "stderr.log")
	outputPath := filepath.Join(opts.Idea.Path, opts.RoundLabel, agent.ID+".md")
	result := Result{
		AgentID:    agent.ID,
		OutputPath: outputPath,
		StdoutPath: stdoutPath,
		StderrPath: stderrPath,
		StartedAt:  now,
	}

	if _, err := os.Stat(outputPath); err == nil {
		result.Skipped = true
		result.SkipReason = "artifact already exists"
		result.CompletedAt = time.Now().UTC()
		result.Duration = result.CompletedAt.Sub(result.StartedAt)
		if err := opts.Store.Append(store.Event{
			Time: result.CompletedAt,
			Type: "agent.skipped",
			Data: map[string]any{"agent": agent.ID, "reason": result.SkipReason, "artifact": outputPath},
		}); err != nil {
			result.ExitError = "event append failed: " + err.Error()
		}
		return result
	}

	if err := os.MkdirAll(agentDir, 0o755); err != nil {
		return failEarly(opts, result, err)
	}
	if err := os.MkdirAll(filepath.Dir(outputPath), 0o755); err != nil {
		return failEarly(opts, result, err)
	}

	questionsDir := filepath.Join(opts.Root, protocol.DeckDir, "runs", opts.RunID, "questions")
	prompt, err := BuildRoundOnePrompt(agent, opts.Idea, opts.Task, outputPath, questionsDir)
	if err != nil {
		return failEarly(opts, result, err)
	}

	if agents.LaunchModeOrDefault(agent.LaunchMode) == agents.LaunchACP {
		return runACPAgent(parent, opts, agent, result, outputPath, stdoutPath, stderrPath, prompt)
	}

	ctx, cancel := context.WithTimeout(parent, timeoutForAgent(opts.Timeout, agent))
	defer cancel()

	cmd, cleanup, err := CommandFor(ctx, opts.Root, agent, prompt)
	if cleanup != nil {
		defer cleanup()
	}
	if err != nil {
		return failEarly(opts, result, err)
	}
	cmd.Dir = opts.Root

	stdoutFile, err := os.Create(stdoutPath)
	if err != nil {
		return failEarly(opts, result, err)
	}
	defer stdoutFile.Close()
	stderrFile, err := os.Create(stderrPath)
	if err != nil {
		return failEarly(opts, result, err)
	}
	defer stderrFile.Close()
	cmd.Stdout = stdoutFile
	cmd.Stderr = stderrFile
	if agent.PromptMode == agents.PromptStdin {
		cmd.Stdin = strings.NewReader(prompt)
	}

	if err := opts.Store.Append(store.Event{
		Time: now,
		Type: "agent.started",
		Data: map[string]any{
			"agent":    agent.ID,
			"artifact": outputPath,
			"stdout":   stdoutPath,
			"stderr":   stderrPath,
		},
	}); err != nil {
		return failEarly(opts, result, fmt.Errorf("event append failed: %w", err))
	}

	err = cmd.Run()
	result.CompletedAt = time.Now().UTC()
	result.Duration = result.CompletedAt.Sub(result.StartedAt)
	if err != nil {
		result.ExitError = err.Error()
		if ctx.Err() != nil {
			result.ExitError = ctx.Err().Error()
		}
	}

	if _, statErr := os.Stat(outputPath); statErr == nil {
		if validateErr := ValidateRoundOneArtifact(outputPath, agent.ID, opts.Idea.Slug); validateErr != nil {
			result.ExitError = combineError(result.ExitError, validateErr)
		} else {
			result.ArtifactOK = true
		}
	}

	eventType := "agent.finished"
	if result.ExitError != "" || !result.ArtifactOK {
		eventType = "agent.failed"
	}
	if err := opts.Store.Append(store.Event{
		Time: result.CompletedAt,
		Type: eventType,
		Data: map[string]any{
			"agent":       agent.ID,
			"artifact":    outputPath,
			"artifact_ok": result.ArtifactOK,
			"duration_ms": result.Duration.Milliseconds(),
			"error":       result.ExitError,
		},
	}); err != nil {
		result.ExitError = combineError(result.ExitError, fmt.Errorf("event append failed: %w", err))
	}

	return result
}

func failEarly(opts Options, result Result, err error) Result {
	result.CompletedAt = time.Now().UTC()
	result.Duration = result.CompletedAt.Sub(result.StartedAt)
	result.ExitError = err.Error()
	if eventErr := opts.Store.Append(store.Event{
		Time: result.CompletedAt,
		Type: "agent.failed",
		Data: map[string]any{
			"agent":       result.AgentID,
			"artifact":    result.OutputPath,
			"artifact_ok": false,
			"duration_ms": result.Duration.Milliseconds(),
			"error":       result.ExitError,
		},
	}); eventErr != nil {
		result.ExitError = combineError(result.ExitError, fmt.Errorf("event append failed: %w", eventErr))
	}
	return result
}

func combineError(primary string, err error) string {
	if err == nil {
		return primary
	}
	if primary == "" {
		return err.Error()
	}
	return primary + "; " + err.Error()
}

func BuildRoundOnePrompt(agent agents.Discovery, idea protocol.IdeaStatus, task, outputPath, questionsDir string) (string, error) {
	promptData, err := os.ReadFile(filepath.Join(idea.Path, "00-prompt.md"))
	if err != nil {
		return "", err
	}
	return fmt.Sprintf(`You are %s, a participant in a Parley Deck cooperation round.

Rules:
- Create exactly this file and no other protocol artifact: %s
- Do not edit any other agent's file.
- Do not overwrite the file if it already exists; report a blocker instead.
- Do not read or reference other agents' round-01 answers.
- Write the complete file, including YAML frontmatter. The first line of the file must be exactly "---".
- Return only a short confirmation with the path written.
- Be concrete, concise, and state trade-offs.
- If you are blocked by missing human input, create one JSON question file under: %s
- Question files use this shape: {"id":"<unique-id>","agent":"%s","prompt":"<question>","details":"<context>","default_answer":"<safe default if any>","risk":"low|normal|high","status":"open","answer":"","created_at":"<RFC3339 time>","answered_at":"0001-01-01T00:00:00Z"}
- If you choose to wait for an answer, poll your question file until status is answered or auto_answered. Otherwise proceed with an explicit assumption in your artifact.

Effective launch config:
- model: %s
- thinking/reasoning/effort/profile: %s
- speed: %s
- sandbox: %s
- approval: %s
- timeoutMs: %d

Idea prompt:
%s

Required file shape:
---
agent: %s
idea: %s
round: 1
date: %s
---

## Summary
## Proposed approach
## Concerns / open questions
## Risks
`, agent.ID, outputPath, questionsDir, agent.ID,
		runtimeValue(agent.Model),
		runtimeValue(firstNonEmpty(agent.Reasoning, agent.Profile)),
		runtimeValue(firstNonEmpty(agent.Speed, agents.DefaultSpeed)),
		runtimeValue(agent.SandboxMode),
		runtimeValue(agent.ApprovalPolicy),
		timeoutMSForAgent(agent),
		string(promptData), agent.ID, idea.Slug, time.Now().Format("2006-01-02")), nil
}

func CommandFor(ctx context.Context, root string, agent agents.Discovery, prompt string) (*exec.Cmd, func(), error) {
	args := make([]string, 0, len(agent.HeadlessArgs))
	for _, arg := range agent.HeadlessArgs {
		switch arg {
		case "{root}":
			args = append(args, root)
		case "{prompt}":
			args = append(args, prompt)
		default:
			args = append(args, arg)
		}
	}
	cmd := exec.CommandContext(ctx, agent.Path, args...)
	cleanup := func() {}
	if agent.IsolateHome {
		env, remove, err := isolatedAgentHome(agent)
		if err != nil {
			return nil, nil, err
		}
		cleanup = remove
		cmd.Env = append(os.Environ(), env...)
		if agent.ID == "hermes" {
			cmd.Env = append(cmd.Env, "HERMES_ACCEPT_HOOKS=1", "HERMES_SESSION_SOURCE=parley")
		}
	}
	return cmd, cleanup, nil
}

func timeoutForAgent(override time.Duration, agent agents.Discovery) time.Duration {
	if override > 0 {
		return override
	}
	return time.Duration(timeoutMSForAgent(agent)) * time.Millisecond
}

func timeoutMSForAgent(agent agents.Discovery) int {
	if agent.TimeoutMS > 0 {
		return agent.TimeoutMS
	}
	return agents.DefaultTimeoutMS
}

func runtimeValue(value string) string {
	if strings.TrimSpace(value) == "" {
		return agents.CLIDefault
	}
	return value
}

func firstNonEmpty(values ...string) string {
	for _, value := range values {
		if strings.TrimSpace(value) != "" {
			return value
		}
	}
	return ""
}

func isolatedAgentHome(agent agents.Discovery) ([]string, func(), error) {
	switch agent.ID {
	case "gemini":
		home, err := isolatedGeminiHome()
		if err != nil {
			return nil, nil, err
		}
		return isolatedHomeEnv(agent, home, "GEMINI_CLI_HOME"), func() { _ = os.RemoveAll(home) }, nil
	case "hermes":
		home, err := isolatedHermesHome()
		if err != nil {
			return nil, nil, err
		}
		return isolatedHomeEnv(agent, home, "HERMES_HOME"), func() { _ = os.RemoveAll(home) }, nil
	default:
		if len(agent.IsolatedHomeEnv) == 0 {
			return nil, nil, fmt.Errorf("no isolated home strategy for %s", agent.ID)
		}
		home, err := os.MkdirTemp("", "parley-"+agent.ID+"-home.*")
		if err != nil {
			return nil, nil, err
		}
		return isolatedHomeEnv(agent, home, ""), func() { _ = os.RemoveAll(home) }, nil
	}
}

func isolatedHomeEnv(agent agents.Discovery, home, fallbackKey string) []string {
	if len(agent.IsolatedHomeEnv) == 0 && fallbackKey != "" {
		return []string{fallbackKey + "=" + home}
	}
	env := make([]string, 0, len(agent.IsolatedHomeEnv))
	for key, template := range agent.IsolatedHomeEnv {
		value := strings.ReplaceAll(template, "{tempdir}", home)
		env = append(env, key+"="+value)
	}
	return env
}

func isolatedGeminiHome() (string, error) {
	sourceHome, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	base, err := os.MkdirTemp("", "parley-gemini-home.*")
	if err != nil {
		return "", err
	}
	geminiDir := filepath.Join(base, ".gemini")
	if err := os.MkdirAll(geminiDir, 0o700); err != nil {
		_ = os.RemoveAll(base)
		return "", err
	}
	copied := 0
	for _, name := range []string{"oauth_creds.json", "google_accounts.json"} {
		source := filepath.Join(sourceHome, ".gemini", name)
		data, err := os.ReadFile(source)
		if err != nil {
			continue
		}
		if err := os.WriteFile(filepath.Join(geminiDir, name), data, 0o600); err != nil {
			_ = os.RemoveAll(base)
			return "", err
		}
		copied++
	}
	if copied == 0 && os.Getenv("GEMINI_API_KEY") == "" && os.Getenv("GOOGLE_API_KEY") == "" {
		_ = os.RemoveAll(base)
		return "", fmt.Errorf("no Gemini OAuth files found in %s and no Gemini API key is set", filepath.Join(sourceHome, ".gemini"))
	}
	settings := []byte(`{"security":{"auth":{"selectedType":"oauth-personal"}}}` + "\n")
	if err := os.WriteFile(filepath.Join(geminiDir, "settings.json"), settings, 0o600); err != nil {
		_ = os.RemoveAll(base)
		return "", err
	}
	return base, nil
}

func isolatedHermesHome() (string, error) {
	sourceHome, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	sourceDir := filepath.Join(sourceHome, ".hermes")
	base, err := os.MkdirTemp("", "parley-hermes-home.*")
	if err != nil {
		return "", err
	}
	for _, dir := range []string{"logs", "sessions", "home"} {
		if err := os.MkdirAll(filepath.Join(base, dir), 0o700); err != nil {
			_ = os.RemoveAll(base)
			return "", err
		}
	}

	copied := 0
	for _, name := range []string{"config.yaml", ".env", "auth.json", "SOUL.md"} {
		source := filepath.Join(sourceDir, name)
		data, err := os.ReadFile(source)
		if err != nil {
			continue
		}
		if err := os.WriteFile(filepath.Join(base, name), data, 0o600); err != nil {
			_ = os.RemoveAll(base)
			return "", err
		}
		copied++
	}
	if copied == 0 {
		_ = os.RemoveAll(base)
		return "", fmt.Errorf("no Hermes config files found in %s", sourceDir)
	}
	return base, nil
}
