package acp

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"strconv"
	"sync"
)

// Handler receives notifications and incoming requests sent by the agent
// to the client (session/update, session/request_permission, fs/*).
// All methods are called from a single dispatcher goroutine; implementations
// should not block for long. ReadTextFile/WriteTextFile may return errors
// that propagate to the agent as JSON-RPC error responses.
type Handler interface {
	SessionUpdate(update SessionUpdate) error
	RequestPermission(req PermissionRequest) (PermissionResult, error)
	ReadTextFile(req FSReadRequest) (FSReadResult, error)
	WriteTextFile(req FSWriteRequest) error
}

// NoopHandler returns minimal-but-valid responses for every Handler method.
// Useful when the client does not advertise fs capabilities and wants every
// permission auto-allowed.
type NoopHandler struct{}

func (NoopHandler) SessionUpdate(SessionUpdate) error { return nil }
func (NoopHandler) RequestPermission(req PermissionRequest) (PermissionResult, error) {
	for _, opt := range req.Options {
		if opt.Kind == PermissionAllowAlways {
			return PermissionResult{Outcome: PermissionOutcome{Outcome: "selected", OptionID: opt.OptionID}}, nil
		}
	}
	for _, opt := range req.Options {
		if opt.Kind == PermissionAllowOnce {
			return PermissionResult{Outcome: PermissionOutcome{Outcome: "selected", OptionID: opt.OptionID}}, nil
		}
	}
	return PermissionResult{Outcome: PermissionOutcome{Outcome: "cancelled"}}, nil
}
func (NoopHandler) ReadTextFile(FSReadRequest) (FSReadResult, error) {
	return FSReadResult{}, errors.New("acp: fs/read_text_file not supported by client")
}
func (NoopHandler) WriteTextFile(FSWriteRequest) error {
	return errors.New("acp: fs/write_text_file not supported by client")
}

// pendingResp pairs an outbound request id with the channel awaiting its reply.
type pendingResp struct {
	result chan Message
}

// Client owns the ACP session: it sends requests on a Transport and routes
// incoming messages to either pending request channels or the Handler.
type Client struct {
	transport  *Transport
	handler    Handler
	clientInfo ClientInfo

	mu       sync.Mutex
	nextID   int64
	pending  map[int64]*pendingResp
	closed   bool
	readErr  error
	readDone chan struct{}
}

// NewClient wires a Transport into a Client. The handler must be non-nil;
// use NoopHandler when no callbacks are needed.
func NewClient(transport *Transport, handler Handler, clientInfo ClientInfo) *Client {
	if handler == nil {
		handler = NoopHandler{}
	}
	return &Client{
		transport:  transport,
		handler:    handler,
		clientInfo: clientInfo,
		pending:    make(map[int64]*pendingResp),
		readDone:   make(chan struct{}),
	}
}

// Start launches the dispatcher goroutine that reads from the transport and
// dispatches messages. Call exactly once per Client.
func (c *Client) Start() {
	go c.readLoop()
}

// Close marks the client as closed and rejects all pending requests.
// It does not stop the dispatcher — that exits when the transport EOFs.
func (c *Client) Close() {
	c.mu.Lock()
	c.closed = true
	pending := c.pending
	c.pending = make(map[int64]*pendingResp)
	c.mu.Unlock()
	for _, p := range pending {
		close(p.result)
	}
}

// Wait blocks until the dispatcher exits (transport EOF or error).
// It returns the error that ended the read loop, or nil on clean EOF.
func (c *Client) Wait() error {
	<-c.readDone
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.readErr
}

// Initialize sends the initialize request. Must be the first call.
func (c *Client) Initialize(ctx context.Context, caps ClientCapabilities) (InitializeResult, error) {
	params := InitializeParams{
		ClientInfo:         c.clientInfo,
		ProtocolVersion:    ProtocolVersion,
		ClientCapabilities: caps,
	}
	var result InitializeResult
	if err := c.call(ctx, MethodInitialize, params, &result); err != nil {
		return InitializeResult{}, err
	}
	return result, nil
}

// NewSession opens a fresh session in the given working directory.
func (c *Client) NewSession(ctx context.Context, params NewSessionParams) (NewSessionResult, error) {
	if params.MCPServers == nil {
		params.MCPServers = []map[string]any{}
	}
	var result NewSessionResult
	if err := c.call(ctx, MethodSessionNew, params, &result); err != nil {
		return NewSessionResult{}, err
	}
	return result, nil
}

// Prompt sends a user turn. Blocks until the agent finishes (or ctx cancels).
// During the call the dispatcher streams SessionUpdate notifications to the
// Handler.
func (c *Client) Prompt(ctx context.Context, sessionID, text string) (PromptResult, error) {
	params := PromptParams{
		SessionID: sessionID,
		Prompt:    []PromptContent{{Type: "text", Text: text}},
	}
	var result PromptResult
	if err := c.call(ctx, MethodSessionPrompt, params, &result); err != nil {
		return PromptResult{}, err
	}
	return result, nil
}

// Cancel notifies the agent that the in-flight prompt for sessionID should
// stop. The agent may still send a final SessionUpdate and a prompt reply
// with stopReason="cancelled".
func (c *Client) Cancel(ctx context.Context, sessionID string) error {
	return c.call(ctx, MethodSessionCancel, CancelParams{SessionID: sessionID}, nil)
}

// call sends a JSON-RPC request and waits for its matching response.
// When result is non-nil the response.Result is decoded into it.
func (c *Client) call(ctx context.Context, method string, params any, result any) error {
	id := c.allocID()
	idNum := json.Number(strconv.FormatInt(id, 10))

	rawParams, err := json.Marshal(params)
	if err != nil {
		return fmt.Errorf("acp: marshal %s params: %w", method, err)
	}

	pending := &pendingResp{result: make(chan Message, 1)}
	c.mu.Lock()
	if c.closed {
		c.mu.Unlock()
		return errors.New("acp: client closed")
	}
	c.pending[id] = pending
	c.mu.Unlock()

	err = c.transport.Write(Message{
		JSONRPC: JSONRPCVersion,
		ID:      &idNum,
		Method:  method,
		Params:  rawParams,
	})
	if err != nil {
		c.removePending(id)
		return fmt.Errorf("acp: send %s: %w", method, err)
	}

	select {
	case <-ctx.Done():
		c.removePending(id)
		return ctx.Err()
	case msg, ok := <-pending.result:
		if !ok {
			return errors.New("acp: connection closed before response")
		}
		if msg.Error != nil {
			return fmt.Errorf("acp: %s rpc error: %s", method, msg.Error.Message)
		}
		if result == nil {
			return nil
		}
		if len(msg.Result) == 0 {
			return nil
		}
		if err := json.Unmarshal(msg.Result, result); err != nil {
			return fmt.Errorf("acp: decode %s result: %w", method, err)
		}
		return nil
	}
}

func (c *Client) allocID() int64 {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.nextID++
	return c.nextID
}

func (c *Client) removePending(id int64) {
	c.mu.Lock()
	delete(c.pending, id)
	c.mu.Unlock()
}

func (c *Client) readLoop() {
	defer close(c.readDone)
	for {
		msg, err := c.transport.Read()
		if err != nil {
			c.mu.Lock()
			c.readErr = err
			pending := c.pending
			c.pending = make(map[int64]*pendingResp)
			c.mu.Unlock()
			for _, p := range pending {
				close(p.result)
			}
			if errors.Is(err, io.EOF) {
				c.mu.Lock()
				c.readErr = nil
				c.mu.Unlock()
			}
			return
		}
		c.dispatch(msg)
	}
}

func (c *Client) dispatch(msg Message) {
	if msg.ID != nil && msg.Method == "" {
		c.dispatchResponse(msg)
		return
	}
	if msg.Method == "" {
		return
	}
	if msg.ID == nil {
		c.dispatchNotification(msg)
		return
	}
	c.dispatchRequest(msg)
}

func (c *Client) dispatchResponse(msg Message) {
	if msg.ID == nil {
		return
	}
	id, err := msg.ID.Int64()
	if err != nil {
		return
	}
	c.mu.Lock()
	pending, ok := c.pending[id]
	delete(c.pending, id)
	c.mu.Unlock()
	if !ok {
		return
	}
	pending.result <- msg
	close(pending.result)
}

func (c *Client) dispatchNotification(msg Message) {
	switch msg.Method {
	case MethodSessionUpdate:
		var update SessionUpdate
		if err := json.Unmarshal(msg.Params, &update); err != nil {
			return
		}
		update.Update.Raw = msg.Params
		_ = c.handler.SessionUpdate(update)
	}
}

func (c *Client) dispatchRequest(msg Message) {
	switch msg.Method {
	case MethodRequestPermission:
		var req PermissionRequest
		if err := json.Unmarshal(msg.Params, &req); err != nil {
			c.respondError(msg.ID, -32602, "invalid params")
			return
		}
		result, err := c.handler.RequestPermission(req)
		if err != nil {
			c.respondError(msg.ID, -32000, err.Error())
			return
		}
		c.respond(msg.ID, result)
	case MethodFSReadTextFile:
		var req FSReadRequest
		if err := json.Unmarshal(msg.Params, &req); err != nil {
			c.respondError(msg.ID, -32602, "invalid params")
			return
		}
		result, err := c.handler.ReadTextFile(req)
		if err != nil {
			c.respondError(msg.ID, -32000, err.Error())
			return
		}
		c.respond(msg.ID, result)
	case MethodFSWriteTextFile:
		var req FSWriteRequest
		if err := json.Unmarshal(msg.Params, &req); err != nil {
			c.respondError(msg.ID, -32602, "invalid params")
			return
		}
		if err := c.handler.WriteTextFile(req); err != nil {
			c.respondError(msg.ID, -32000, err.Error())
			return
		}
		c.respond(msg.ID, struct{}{})
	default:
		c.respondError(msg.ID, -32601, "method not implemented: "+msg.Method)
	}
}

func (c *Client) respond(id *json.Number, result any) {
	encoded, err := json.Marshal(result)
	if err != nil {
		c.respondError(id, -32603, "encode error: "+err.Error())
		return
	}
	_ = c.transport.Write(Message{
		JSONRPC: JSONRPCVersion,
		ID:      id,
		Result:  encoded,
	})
}

func (c *Client) respondError(id *json.Number, code int, message string) {
	_ = c.transport.Write(Message{
		JSONRPC: JSONRPCVersion,
		ID:      id,
		Error:   &RPCError{Code: code, Message: message},
	})
}
