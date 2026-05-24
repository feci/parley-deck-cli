package acp

import (
	"context"
	"encoding/json"
	"io"
	"strings"
	"sync"
	"testing"
	"time"
)

// pipePair returns two connected in-memory bidirectional pipes:
// clientToAgent + agentToClient.
type pipePair struct {
	clientReader io.Reader
	clientWriter io.Writer
	agentReader  io.Reader
	agentWriter  io.Writer
}

func newPipePair() *pipePair {
	// Pipe 1: agent writes → client reads (agent → client direction).
	clientReadEnd, agentWriteEnd := io.Pipe()
	// Pipe 2: client writes → agent reads (client → agent direction).
	agentReadEnd, clientWriteEnd := io.Pipe()
	return &pipePair{
		clientReader: clientReadEnd,
		clientWriter: clientWriteEnd,
		agentReader:  agentReadEnd,
		agentWriter:  agentWriteEnd,
	}
}

// stubAgent reads NDJSON requests from the client and replies with the
// configured handler. It runs until ctx is cancelled or the stream closes.
type stubAgent struct {
	transport *Transport
}

func (s *stubAgent) run(ctx context.Context, handle func(msg Message) []Message) {
	for {
		msg, err := s.transport.Read()
		if err != nil {
			return
		}
		if ctx.Err() != nil {
			return
		}
		replies := handle(msg)
		if len(replies) == 0 {
			continue
		}
		// Writes happen in a goroutine so the read loop can drain any
		// client-side responses (e.g. permission replies) that the
		// dispatcher writes while we are still emitting our reply batch.
		go func(replies []Message) {
			for _, reply := range replies {
				_ = s.transport.Write(reply)
			}
		}(replies)
	}
}

func TestInitializeRoundTrip(t *testing.T) {
	pipes := newPipePair()
	clientTransport := NewTransport(pipes.clientReader, pipes.clientWriter)
	agentTransport := NewTransport(pipes.agentReader, pipes.agentWriter)

	agent := &stubAgent{transport: agentTransport}
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	go agent.run(ctx, func(msg Message) []Message {
		if msg.Method != MethodInitialize {
			return nil
		}
		result := InitializeResult{
			ProtocolVersion: ProtocolVersion,
			AgentInfo:       &AgentInfo{Name: "stub", Version: "0.0.1"},
		}
		raw, _ := json.Marshal(result)
		return []Message{{
			JSONRPC: JSONRPCVersion,
			ID:      msg.ID,
			Result:  raw,
		}}
	})

	client := NewClient(clientTransport, NoopHandler{}, ClientInfo{Name: "parley", Version: "test"})
	client.Start()
	defer client.Close()

	res, err := client.Initialize(ctx, ClientCapabilities{})
	if err != nil {
		t.Fatalf("initialize: %v", err)
	}
	if res.AgentInfo == nil || res.AgentInfo.Name != "stub" {
		t.Fatalf("unexpected agent info: %+v", res.AgentInfo)
	}
}

func TestPromptStreamsSessionUpdates(t *testing.T) {
	pipes := newPipePair()
	clientTransport := NewTransport(pipes.clientReader, pipes.clientWriter)
	agentTransport := NewTransport(pipes.agentReader, pipes.agentWriter)

	agent := &stubAgent{transport: agentTransport}
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	go agent.run(ctx, func(msg Message) []Message {
		switch msg.Method {
		case MethodSessionNew:
			raw, _ := json.Marshal(NewSessionResult{SessionID: "sess-1"})
			return []Message{{JSONRPC: JSONRPCVersion, ID: msg.ID, Result: raw}}
		case MethodSessionPrompt:
			updates := []Message{}
			for _, chunk := range []string{"Hello ", "world."} {
				notif := SessionUpdate{
					SessionID: "sess-1",
					Update: UpdatePayload{
						SessionUpdate: UpdateAgentMessageChunk,
						Content:       &UpdateContent{Type: "text", Text: chunk},
					},
				}
				params, _ := json.Marshal(notif)
				updates = append(updates, Message{JSONRPC: JSONRPCVersion, Method: MethodSessionUpdate, Params: params})
			}
			result, _ := json.Marshal(PromptResult{StopReason: "complete"})
			updates = append(updates, Message{JSONRPC: JSONRPCVersion, ID: msg.ID, Result: result})
			return updates
		}
		return nil
	})

	captured := &captureHandler{}
	client := NewClient(clientTransport, captured, ClientInfo{Name: "parley", Version: "test"})
	client.Start()
	defer client.Close()

	sess, err := client.NewSession(ctx, NewSessionParams{CWD: "/tmp"})
	if err != nil {
		t.Fatalf("new session: %v", err)
	}
	if sess.SessionID != "sess-1" {
		t.Fatalf("unexpected sessionId %q", sess.SessionID)
	}

	res, err := client.Prompt(ctx, sess.SessionID, "Hi")
	if err != nil {
		t.Fatalf("prompt: %v", err)
	}
	if res.StopReason != "complete" {
		t.Fatalf("unexpected stop reason %q", res.StopReason)
	}
	captured.waitForChunks(t, 2)
	if got := captured.text(); got != "Hello world." {
		t.Fatalf("unexpected streamed text: %q", got)
	}
}

func TestPermissionAutoAllow(t *testing.T) {
	pipes := newPipePair()
	clientTransport := NewTransport(pipes.clientReader, pipes.clientWriter)
	agentTransport := NewTransport(pipes.agentReader, pipes.agentWriter)

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	// Agent sends a permission request before responding to prompt.
	agent := &stubAgent{transport: agentTransport}
	go func() {
		// Wait until prompt arrives, then issue a permission request and
		// finally respond to the original prompt.
		agent.run(ctx, func(msg Message) []Message {
			if msg.Method != MethodSessionPrompt {
				return nil
			}
			permParams, _ := json.Marshal(PermissionRequest{
				SessionID: "sess",
				Options: []PermissionOption{
					{OptionID: "yes", Name: "allow", Kind: PermissionAllowAlways},
				},
				ToolCall: PermissionToolCall{ToolCallID: "tc1"},
			})
			permID := json.Number("99")
			permReq := Message{JSONRPC: JSONRPCVersion, ID: &permID, Method: MethodRequestPermission, Params: permParams}
			result, _ := json.Marshal(PromptResult{StopReason: "complete"})
			return []Message{
				permReq,
				{JSONRPC: JSONRPCVersion, ID: msg.ID, Result: result},
			}
		})
	}()

	client := NewClient(clientTransport, NoopHandler{}, ClientInfo{Name: "parley", Version: "test"})
	client.Start()
	defer client.Close()

	if _, err := client.Prompt(ctx, "sess", "Hi"); err != nil {
		t.Fatalf("prompt: %v", err)
	}
}

type captureHandler struct {
	mu     sync.Mutex
	chunks []string
}

func (c *captureHandler) SessionUpdate(update SessionUpdate) error {
	if update.Update.SessionUpdate != UpdateAgentMessageChunk || update.Update.Content == nil {
		return nil
	}
	c.mu.Lock()
	c.chunks = append(c.chunks, update.Update.Content.Text)
	c.mu.Unlock()
	return nil
}

func (c *captureHandler) RequestPermission(req PermissionRequest) (PermissionResult, error) {
	return NoopHandler{}.RequestPermission(req)
}

func (c *captureHandler) ReadTextFile(req FSReadRequest) (FSReadResult, error) {
	return NoopHandler{}.ReadTextFile(req)
}

func (c *captureHandler) WriteTextFile(req FSWriteRequest) error {
	return NoopHandler{}.WriteTextFile(req)
}

func (c *captureHandler) waitForChunks(t *testing.T, count int) {
	t.Helper()
	deadline := time.Now().Add(2 * time.Second)
	for time.Now().Before(deadline) {
		c.mu.Lock()
		ready := len(c.chunks) >= count
		c.mu.Unlock()
		if ready {
			return
		}
		time.Sleep(10 * time.Millisecond)
	}
	c.mu.Lock()
	got := len(c.chunks)
	c.mu.Unlock()
	t.Fatalf("expected %d chunks, got %d", count, got)
}

func (c *captureHandler) text() string {
	c.mu.Lock()
	defer c.mu.Unlock()
	return strings.Join(c.chunks, "")
}
