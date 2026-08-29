package runner

import (
	"bufio"
	"context"
	"encoding/json"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"testing"
	"time"

	"parley-deck-cli/internal/acp"
	"parley-deck-cli/internal/agents"
	"parley-deck-cli/internal/protocol"
	"parley-deck-cli/internal/store"
)

func TestRunRoundOneRoutesACPAgent(t *testing.T) {
	root := t.TempDir()
	if err := protocol.InitWorkspace(root); err != nil {
		t.Fatal(err)
	}
	idea, err := protocol.CreateIdea(root, "ACP runner test", []string{"fake-acp"})
	if err != nil {
		t.Fatal(err)
	}
	runStore := store.New(filepath.Join(root, protocol.DeckDir, "runs", "acp-run"))

	results := RunRoundOne(context.Background(), Options{
		Root:  root,
		RunID: "acp-run",
		Idea:  idea,
		Task:  "ACP runner test",
		Agents: []agents.Discovery{
			{
				Spec: agents.Spec{
					ID:         "fake-acp",
					LaunchMode: agents.LaunchACP,
					ACPArgs:    []string{"-test.run=TestFakeACPAgentHelper", "--", "parley-fake-acp-agent"},
					PromptMode: agents.PromptStdin,
				},
				Path:  os.Args[0],
				Found: true,
			},
		},
		Timeout: 10 * time.Second,
		Store:   runStore,
	})

	if len(results) != 1 {
		t.Fatalf("got %d results, want 1", len(results))
	}
	if !results[0].ArtifactOK || results[0].ExitError != "" {
		t.Fatalf("unexpected result: %+v", results[0])
	}
	if _, statErr := os.Stat(filepath.Join(idea.Path, "round-01", "fake-acp.md")); statErr != nil {
		t.Fatalf("artifact missing: %v", statErr)
	}

	events, err := runStore.Load()
	if err != nil {
		t.Fatal(err)
	}
	wantTypes := []string{"agent.acp.initialized", "agent.acp.session_opened", "agent.acp.prompt_completed", "agent.finished"}
	seen := map[string]bool{}
	for _, ev := range events {
		seen[ev.Type] = true
	}
	for _, want := range wantTypes {
		if !seen[want] {
			t.Errorf("event %q missing; got types=%v", want, eventTypes(events))
		}
	}
}

func TestRunRoundOneACPAgentMissingArgsFails(t *testing.T) {
	root := t.TempDir()
	if err := protocol.InitWorkspace(root); err != nil {
		t.Fatal(err)
	}
	idea, err := protocol.CreateIdea(root, "ACP missing args", []string{"misconfigured"})
	if err != nil {
		t.Fatal(err)
	}
	runStore := store.New(filepath.Join(root, protocol.DeckDir, "runs", "acp-bad"))

	results := RunRoundOne(context.Background(), Options{
		Root:  root,
		RunID: "acp-bad",
		Idea:  idea,
		Task:  "ACP missing args",
		Agents: []agents.Discovery{
			{
				Spec: agents.Spec{
					ID:         "misconfigured",
					LaunchMode: agents.LaunchACP,
					ACPArgs:    nil,
				},
				Path:  os.Args[0],
				Found: true,
			},
		},
		Timeout: 5 * time.Second,
		Store:   runStore,
	})

	if len(results) != 1 {
		t.Fatalf("got %d results, want 1", len(results))
	}
	if results[0].ExitError == "" {
		t.Fatalf("expected error for missing ACPArgs, got: %+v", results[0])
	}
	if !strings.Contains(results[0].ExitError, "ACPArgs is empty") {
		t.Fatalf("error mismatch: %q", results[0].ExitError)
	}
}

func eventTypes(events []store.Event) []string {
	out := make([]string, 0, len(events))
	for _, ev := range events {
		out = append(out, ev.Type)
	}
	return out
}

// TestFakeACPAgentHelper acts as a minimal ACP-speaking child process when
// invoked with the parley-fake-acp-agent marker in argv. It mirrors the
// fake headless helper but talks JSON-RPC 2.0 + NDJSON on stdio.
func TestFakeACPAgentHelper(t *testing.T) {
	if !hasArg("parley-fake-acp-agent") {
		return
	}
	reader := bufio.NewReader(os.Stdin)
	writer := os.Stdout
	for {
		line, err := readJSONLine(reader)
		if err != nil {
			if err == io.EOF {
				return
			}
			t.Fatalf("read: %v", err)
		}
		if len(line) == 0 {
			continue
		}
		var msg acp.Message
		if jsonErr := json.Unmarshal(line, &msg); jsonErr != nil {
			continue
		}
		switch msg.Method {
		case acp.MethodInitialize:
			respondACPHelper(t, writer, msg.ID, acp.InitializeResult{
				ProtocolVersion: acp.ProtocolVersion,
				AgentInfo:       &acp.AgentInfo{Name: "fake-acp", Version: "0.0.1"},
			})
		case acp.MethodSessionNew:
			respondACPHelper(t, writer, msg.ID, acp.NewSessionResult{SessionID: "fake-session"})
		case acp.MethodSessionPrompt:
			var params acp.PromptParams
			if jsonErr := json.Unmarshal(msg.Params, &params); jsonErr != nil {
				t.Fatalf("decode prompt params: %v", jsonErr)
			}
			text := ""
			for _, block := range params.Prompt {
				if block.Type == "text" {
					text += block.Text
				}
			}
			outRe := regexp.MustCompile(`(?m)^- Create exactly this file and no other protocol artifact: (.+)$`)
			outMatch := outRe.FindStringSubmatch(text)
			if len(outMatch) != 2 {
				t.Fatalf("missing output path in prompt:\n%s", text)
			}
			ideaRe := regexp.MustCompile(`(?m)^idea: (.+)$`)
			ideaMatch := ideaRe.FindStringSubmatch(text)
			if len(ideaMatch) != 2 {
				t.Fatalf("missing idea slug in prompt:\n%s", text)
			}
			body := `---
agent: fake-acp
idea: ` + ideaMatch[1] + `
round: 1
date: 2026-05-23
---

## Summary
Fake ACP artifact.

## Proposed approach
Speak ACP.

## Existing alternatives
Searched the stdlib and the lockfile; nothing ships this. Hand-built route is correct.

## Concerns / open questions
None.

## Risks
None.
`
			if writeErr := os.WriteFile(outMatch[1], []byte(body), 0o644); writeErr != nil {
				t.Fatalf("write artifact: %v", writeErr)
			}
			respondACPHelper(t, writer, msg.ID, acp.PromptResult{StopReason: "complete"})
			os.Exit(0)
		default:
			// Unknown methods are ignored — keep the conversation alive.
		}
	}
}

func readJSONLine(r *bufio.Reader) ([]byte, error) {
	var collected []byte
	for {
		chunk, isPrefix, err := r.ReadLine()
		if err != nil {
			if err == io.EOF && len(collected) > 0 {
				return collected, nil
			}
			return nil, err
		}
		collected = append(collected, chunk...)
		if !isPrefix {
			return collected, nil
		}
	}
}

func respondACPHelper(t *testing.T, w io.Writer, id *json.Number, result any) {
	encoded, err := json.Marshal(result)
	if err != nil {
		t.Fatalf("marshal result: %v", err)
	}
	msg := acp.Message{
		JSONRPC: acp.JSONRPCVersion,
		ID:      id,
		Result:  encoded,
	}
	out, err := json.Marshal(msg)
	if err != nil {
		t.Fatalf("marshal message: %v", err)
	}
	if _, err := w.Write(append(out, '\n')); err != nil {
		t.Fatalf("write: %v", err)
	}
}
