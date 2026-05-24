// Package acp implements a minimal client for the Agent Client Protocol
// (https://agentclientprotocol.dev) — JSON-RPC 2.0 over NDJSON on a child
// process's stdio. It mirrors the subset used by AionUi's ProcessAcpClient.
package acp

import "encoding/json"

// ProtocolVersion is the ACP protocol version we declare during initialize.
// Bumping this should be a deliberate decision tied to handler updates.
const ProtocolVersion = 1

// JSONRPCVersion is the only JSON-RPC version we speak.
const JSONRPCVersion = "2.0"

// Message is the union of JSON-RPC request/response/notification.
// Presence/absence of ID and Method disambiguates the variant.
type Message struct {
	JSONRPC string          `json:"jsonrpc"`
	ID      *json.Number    `json:"id,omitempty"`
	Method  string          `json:"method,omitempty"`
	Params  json.RawMessage `json:"params,omitempty"`
	Result  json.RawMessage `json:"result,omitempty"`
	Error   *RPCError       `json:"error,omitempty"`
}

// RPCError is the standard JSON-RPC 2.0 error object.
type RPCError struct {
	Code    int             `json:"code"`
	Message string          `json:"message"`
	Data    json.RawMessage `json:"data,omitempty"`
}

func (e *RPCError) Error() string {
	if e == nil {
		return ""
	}
	return e.Message
}

// InitializeParams is sent by the client to the agent right after spawn.
type InitializeParams struct {
	ClientInfo         ClientInfo         `json:"clientInfo"`
	ProtocolVersion    int                `json:"protocolVersion"`
	ClientCapabilities ClientCapabilities `json:"clientCapabilities"`
}

type ClientInfo struct {
	Name    string `json:"name"`
	Version string `json:"version"`
}

// ClientCapabilities advertises which features the client implements.
// Leaving fs unset signals the agent should use its own filesystem tools
// rather than routing fs/read_text_file and fs/write_text_file through us.
type ClientCapabilities struct {
	FS *FSCapabilities `json:"fs,omitempty"`
}

type FSCapabilities struct {
	ReadTextFile  bool `json:"readTextFile"`
	WriteTextFile bool `json:"writeTextFile"`
}

// InitializeResult is the agent's response to initialize.
type InitializeResult struct {
	ProtocolVersion   int                    `json:"protocolVersion"`
	AgentCapabilities map[string]any         `json:"agentCapabilities,omitempty"`
	AgentInfo         *AgentInfo             `json:"agentInfo,omitempty"`
	AuthMethods       []map[string]any       `json:"authMethods,omitempty"`
	Meta              map[string]any         `json:"_meta,omitempty"`
	Modes             map[string]any         `json:"modes,omitempty"`
}

type AgentInfo struct {
	Name    string `json:"name"`
	Version string `json:"version"`
	Title   string `json:"title,omitempty"`
}

// NewSessionParams creates a fresh session in the agent's working directory.
type NewSessionParams struct {
	CWD                  string         `json:"cwd"`
	MCPServers           []map[string]any `json:"mcpServers"`
	AdditionalDirectories []string      `json:"additionalDirectories,omitempty"`
}

type NewSessionResult struct {
	SessionID string `json:"sessionId"`
}

// PromptParams sends user content to an existing session.
// The prompt is a slice of content blocks per the ACP spec; for parley
// we only need the "text" block.
type PromptParams struct {
	SessionID string          `json:"sessionId"`
	Prompt    []PromptContent `json:"prompt"`
}

type PromptContent struct {
	Type     string `json:"type"`
	Text     string `json:"text,omitempty"`
	Data     string `json:"data,omitempty"`
	MimeType string `json:"mimeType,omitempty"`
	URI      string `json:"uri,omitempty"`
}

// PromptResult is what the agent returns after the turn is complete.
type PromptResult struct {
	StopReason string         `json:"stopReason,omitempty"`
	Meta       map[string]any `json:"_meta,omitempty"`
}

// CancelParams cancels an in-flight prompt for the given session.
type CancelParams struct {
	SessionID string `json:"sessionId"`
}

// SessionUpdate is the notification the agent sends while answering a prompt.
// The discriminator is the nested "sessionUpdate" field on Update.
type SessionUpdate struct {
	SessionID string         `json:"sessionId"`
	Update    UpdatePayload  `json:"update"`
}

type UpdatePayload struct {
	SessionUpdate string          `json:"sessionUpdate"`
	Content       *UpdateContent  `json:"content,omitempty"`
	ToolCallID    string          `json:"toolCallId,omitempty"`
	Status        string          `json:"status,omitempty"`
	Title         string          `json:"title,omitempty"`
	Kind          string          `json:"kind,omitempty"`
	RawInput      json.RawMessage `json:"rawInput,omitempty"`
	Entries       []PlanEntry     `json:"entries,omitempty"`
	Used          *int            `json:"used,omitempty"`
	Size          *int            `json:"size,omitempty"`
	Raw           json.RawMessage `json:"-"`
}

type UpdateContent struct {
	Type     string `json:"type"`
	Text     string `json:"text,omitempty"`
	Data     string `json:"data,omitempty"`
	MimeType string `json:"mimeType,omitempty"`
	URI      string `json:"uri,omitempty"`
}

type PlanEntry struct {
	Content  string `json:"content"`
	Status   string `json:"status"`
	Priority string `json:"priority,omitempty"`
}

// PermissionRequest is the params for session/request_permission.
type PermissionRequest struct {
	SessionID string             `json:"sessionId"`
	Options   []PermissionOption `json:"options"`
	ToolCall  PermissionToolCall `json:"toolCall"`
}

type PermissionOption struct {
	OptionID string `json:"optionId"`
	Name     string `json:"name"`
	Kind     string `json:"kind"`
}

type PermissionToolCall struct {
	ToolCallID string          `json:"toolCallId"`
	Title      string          `json:"title,omitempty"`
	Kind       string          `json:"kind,omitempty"`
	Status     string          `json:"status,omitempty"`
	RawInput   json.RawMessage `json:"rawInput,omitempty"`
}

type PermissionResult struct {
	Outcome PermissionOutcome `json:"outcome"`
}

type PermissionOutcome struct {
	Outcome  string `json:"outcome"`
	OptionID string `json:"optionId,omitempty"`
}

// FSReadRequest is the params for fs/read_text_file.
type FSReadRequest struct {
	SessionID string `json:"sessionId"`
	Path      string `json:"path"`
}

type FSReadResult struct {
	Content string `json:"content"`
}

// FSWriteRequest is the params for fs/write_text_file.
type FSWriteRequest struct {
	SessionID string `json:"sessionId"`
	Path      string `json:"path"`
	Content   string `json:"content"`
}

// Standard ACP method names.
const (
	MethodInitialize        = "initialize"
	MethodSessionNew        = "session/new"
	MethodSessionPrompt     = "session/prompt"
	MethodSessionCancel     = "session/cancel"
	MethodSessionUpdate     = "session/update"
	MethodRequestPermission = "session/request_permission"
	MethodFSReadTextFile    = "fs/read_text_file"
	MethodFSWriteTextFile   = "fs/write_text_file"
)

// Standard session/update discriminators.
const (
	UpdateAgentMessageChunk = "agent_message_chunk"
	UpdateAgentThoughtChunk = "agent_thought_chunk"
	UpdateToolCall          = "tool_call"
	UpdateToolCallUpdate    = "tool_call_update"
	UpdatePlan              = "plan"
	UpdateAvailableCommands = "available_commands_update"
	UpdateUserMessageChunk  = "user_message_chunk"
	UpdateUsage             = "usage_update"
)

// Standard permission option kinds (mirrors AionUi / ACP spec).
const (
	PermissionAllowOnce    = "allow_once"
	PermissionAllowAlways  = "allow_always"
	PermissionRejectOnce   = "reject_once"
	PermissionRejectAlways = "reject_always"
)
