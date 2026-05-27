package agents

import "strings"

// ACPBackend describes a CLI agent that supports the Agent Client Protocol
// (JSON-RPC 2.0 over NDJSON on stdio). The catalog mirrors AionUi's
// ACP_BACKENDS_ALL table so parley-deck can auto-detect the same set of
// agents on the user's PATH.
type ACPBackend struct {
	// ID is the parley-deck agent identifier, e.g. "qwen".
	ID string
	// Name is the human-readable display name.
	Name string
	// Command is the binary name resolved via exec.LookPath.
	Command string
	// ACPArgs are the arguments that put the CLI into ACP mode.
	// nil means Parley knows about the backend but has no safe default args.
	// [] means the binary speaks ACP by default.
	ACPArgs []string
	// Notes captures backend-specific behavior worth surfacing.
	Notes string
}

// ACPCatalog returns the table of well-known ACP-capable CLIs. Entries whose
// primary parley-deck ID already lives in DefaultSpecs use the same ID, so ACP
// becomes a selectable launch mode on that agent instead of a duplicate agent.
func ACPCatalog() []ACPBackend {
	return []ACPBackend{
		{ID: "claude", Name: "Claude Code (ACP)", Command: "claude", ACPArgs: []string{"--experimental-acp"},
			Notes: "ACP mode is available with --experimental-acp"},
		{ID: "qwen", Name: "Qwen Code", Command: "qwen", ACPArgs: []string{"--acp"}},
		{ID: "codex", Name: "Codex (ACP)", Command: "codex", ACPArgs: nil,
			Notes: "Codex ACP is not enabled by default because the local codex CLI exposes no stable ACP launch args; configure acp_args locally when available"},
		{ID: "codebuddy", Name: "CodeBuddy", Command: "codebuddy", ACPArgs: []string{"--acp"}},
		{ID: "goose", Name: "Goose", Command: "goose", ACPArgs: []string{"acp"}},
		{ID: "auggie", Name: "Augment Code", Command: "auggie", ACPArgs: []string{"--acp"}},
		{ID: "kimi", Name: "Kimi CLI", Command: "kimi", ACPArgs: []string{"acp"}},
		{ID: "opencode", Name: "OpenCode", Command: "opencode", ACPArgs: []string{"acp"}},
		{ID: "droid", Name: "Factory Droid", Command: "droid", ACPArgs: []string{"exec", "--output-format", "acp"}},
		{ID: "copilot", Name: "GitHub Copilot", Command: "copilot", ACPArgs: []string{"--acp", "--stdio"}},
		{ID: "qoder", Name: "Qoder CLI", Command: "qodercli", ACPArgs: []string{"--acp"}},
		{ID: "vibe", Name: "Mistral Vibe", Command: "vibe-acp", ACPArgs: []string{}},
		{ID: "cursor", Name: "Cursor Agent", Command: "agent", ACPArgs: []string{"acp"},
			Notes: "Cursor CLI binary is the generic name `agent`; PATH lookup may collide with other tools"},
		{ID: "kiro", Name: "Kiro", Command: "kiro-cli", ACPArgs: []string{"acp"}},
		{ID: "hermes", Name: "Hermes Agent (ACP)", Command: "hermes", ACPArgs: []string{"acp"},
			Notes: "ACP mode is available with hermes acp"},
		{ID: "snow", Name: "Snow CLI", Command: "snow", ACPArgs: []string{"--acp"}},
	}
}

func mergeACPCatalog(specs []Spec, catalog []ACPBackend) []Spec {
	out := make([]Spec, len(specs))
	copy(out, specs)
	byID := make(map[string]int, len(out))
	for i, spec := range out {
		byID[spec.ID] = i
	}
	for _, backend := range catalog {
		if index, ok := byID[backend.ID]; ok {
			out[index] = mergeACPBackend(out[index], backend)
			continue
		}
		out = append(out, specFromACPBackend(backend))
		byID[backend.ID] = len(out) - 1
	}
	return out
}

func mergeACPBackend(spec Spec, backend ACPBackend) Spec {
	if len(spec.Commands) == 0 {
		spec.Commands = []string{backend.Command}
	}
	if len(spec.VersionArgs) == 0 {
		spec.VersionArgs = []string{"--version"}
	}
	if backend.ACPArgs != nil {
		spec.ACPArgs = cloneStringSlice(backend.ACPArgs)
	}
	if backend.Notes != "" {
		spec.Notes = joinNotes(spec.Notes, backend.Notes)
	}
	if spec.Sources != nil {
		spec.Sources["acp_args"] = SourceBuiltIn
		if backend.Notes != "" {
			spec.Sources["notes"] = SourceBuiltIn
		}
	}
	return spec
}

// ACPSpecs converts the ACP catalog into Spec entries. Each spec is marked
// with LaunchACP so the runner picks the JSON-RPC code path; HeadlessArgs and
// HeadlessMode are intentionally left empty.
func ACPSpecs() []Spec {
	catalog := ACPCatalog()
	specs := make([]Spec, 0, len(catalog))
	for _, backend := range catalog {
		if backend.ACPArgs == nil {
			continue
		}
		specs = append(specs, specFromACPBackend(backend))
	}
	return specs
}

func specFromACPBackend(backend ACPBackend) Spec {
	notes := backend.Notes
	if notes == "" {
		notes = "ACP backend; spawned with " + backend.Command + " " + joinACPArgs(backend.ACPArgs)
	}
	return withBuiltinSources(Spec{
		ID:                    backend.ID,
		Commands:              []string{backend.Command},
		VersionArgs:           []string{"--version"},
		LaunchMode:            LaunchACP,
		ACPArgs:               cloneStringSlice(backend.ACPArgs),
		InteractivePromptMode: InteractivePromptNone,
		InteractiveInvoke:     InteractiveInvokePrintOnly,
		InteractivePollMS:     DefaultInteractivePollMS,
		PromptMode:            PromptStdin,
		SandboxMode:           CLIDefault,
		ApprovalPolicy:        CLIDefault,
		Model:                 CLIDefault,
		Reasoning:             CLIDefault,
		Profile:               CLIDefault,
		Speed:                 DefaultSpeed,
		TimeoutMS:             DefaultTimeoutMS,
		ExternalBackend:       ExternalUnknown,
		Telemetry:             "ACP session/update events",
		Notes:                 notes,
	})
}

func joinACPArgs(args []string) string {
	if args == nil {
		return "(not configured)"
	}
	if len(args) == 0 {
		return "(no args)"
	}
	return strings.Join(args, " ")
}

func joinNotes(current, next string) string {
	current = strings.TrimSpace(current)
	next = strings.TrimSpace(next)
	if current == "" {
		return next
	}
	if next == "" || strings.Contains(current, next) {
		return current
	}
	return current + "; " + next
}

func cloneStringSlice(values []string) []string {
	if values == nil {
		return nil
	}
	return append([]string{}, values...)
}
