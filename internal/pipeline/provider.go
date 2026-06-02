package pipeline

import "errors"

// Capability is a generic, provider-agnostic action name (§12.4). Concrete
// providers (Vercel, Atlassian, ...) sit behind this vocabulary; the protocol
// never depends on a specific vendor.
type Capability string

const (
	CapDeployPreview    Capability = "deploy.preview"
	CapDeployProduction Capability = "deploy.production"
	CapRuntimeRollback  Capability = "runtime.rollback"
	CapMonitorAlert     Capability = "monitor.alert"
	CapIssueCreate      Capability = "issue.create"
	CapNotifySend       Capability = "notify.send"
)

// IsProduction reports whether a capability mutates production. Such actions
// are non-bypassable: the driver may only Plan them for real execution behind
// an approved production gate (§12.8).
func (c Capability) IsProduction() bool {
	return c == CapDeployProduction || c == CapRuntimeRollback
}

var (
	// ErrUnsupportedCapability means the provider cannot perform the capability;
	// the driver must halt the block before consensus (§12.9), not degrade.
	ErrUnsupportedCapability = errors.New("provider does not support capability")
	// ErrProductionGateRequired means a production action was requested for real
	// execution without an approved gate.
	ErrProductionGateRequired = errors.New("production capability requires an approved gate")
)

// ActionRequest is a generic, provider-neutral request to perform a side effect.
type ActionRequest struct {
	Capability   Capability
	Target       string
	Params       map[string]string
	DryRun       bool
	GateApproved bool // true only when the driver holds an approved gate for this action
}

// ProviderCall is the concrete external invocation the driver/harness performs
// (typically an MCP tool call). The Go CLI PLANS it; it does not execute it —
// the agents-write-markdown / driver-executes boundary (§12.4) means the actual
// call is made by the orchestrating harness, never by this process.
type ProviderCall struct {
	Provider string            `json:"provider"`
	Tool     string            `json:"tool"`
	Args     map[string]string `json:"args"`
	DryRun   bool              `json:"dry_run"`
}

// Provider maps generic capabilities to concrete external calls.
type Provider interface {
	Name() string
	Supports(c Capability) bool
	Plan(req ActionRequest) (ProviderCall, error)
}

// guardProduction enforces the non-bypassable production rule shared by all
// providers: a production capability may only be planned for real (non-dry-run)
// execution when the request carries an approved gate.
func guardProduction(req ActionRequest) error {
	if req.Capability.IsProduction() && !req.DryRun && !req.GateApproved {
		return ErrProductionGateRequired
	}
	return nil
}

// VercelProvider is the first concrete deploy/runtime provider. It maps generic
// capabilities to the Vercel MCP tool the harness will invoke.
type VercelProvider struct{}

func (VercelProvider) Name() string { return "vercel" }

var vercelTools = map[Capability]string{
	CapDeployPreview:    "mcp__claude_ai_Vercel__deploy_to_vercel",
	CapDeployProduction: "mcp__claude_ai_Vercel__deploy_to_vercel",
	CapRuntimeRollback:  "mcp__claude_ai_Vercel__deploy_to_vercel",
	CapMonitorAlert:     "mcp__claude_ai_Vercel__get_runtime_logs",
}

func (VercelProvider) Supports(c Capability) bool {
	_, ok := vercelTools[c]
	return ok
}

func (p VercelProvider) Plan(req ActionRequest) (ProviderCall, error) {
	tool, ok := vercelTools[req.Capability]
	if !ok {
		return ProviderCall{}, ErrUnsupportedCapability
	}
	if err := guardProduction(req); err != nil {
		return ProviderCall{}, err
	}
	args := map[string]string{"target": req.Target}
	for k, v := range req.Params {
		args[k] = v
	}
	switch req.Capability {
	case CapDeployProduction:
		args["production"] = "true"
	case CapRuntimeRollback:
		args["rollback"] = "true"
	}
	return ProviderCall{Provider: p.Name(), Tool: tool, Args: args, DryRun: req.DryRun}, nil
}

// NoopProvider supports every capability and plans an inert call. It is used for
// tests and for dry-run pipelines that should never touch a real backend.
type NoopProvider struct{}

func (NoopProvider) Name() string             { return "noop" }
func (NoopProvider) Supports(Capability) bool { return true }

func (p NoopProvider) Plan(req ActionRequest) (ProviderCall, error) {
	if err := guardProduction(req); err != nil {
		return ProviderCall{}, err
	}
	return ProviderCall{Provider: p.Name(), Tool: "noop", Args: map[string]string{"capability": string(req.Capability), "target": req.Target}, DryRun: req.DryRun}, nil
}
