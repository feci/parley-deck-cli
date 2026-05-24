package agents

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
	// ACPArgs are the arguments that put the CLI into ACP mode
	// (e.g. ["--experimental-acp"], ["acp"], ["--acp"], or [] when the
	// binary speaks ACP by default).
	ACPArgs []string
	// Notes captures backend-specific behavior worth surfacing.
	Notes string
}

// ACPCatalog returns the table of well-known ACP-capable CLIs.
// Order matches AionUi's ACP_BACKENDS_ALL so new entries land in one place.
// Entries whose primary parley-deck ID already lives in DefaultSpecs (codex,
// claude, gemini, hermes) intentionally use distinct catalog IDs ("hermes-acp",
// "claude-acp", "codex-acp") so the existing headless text-mode launches
// remain the default for those four. Users can opt in via the ACP entry.
func ACPCatalog() []ACPBackend {
	return []ACPBackend{
		{ID: "claude-acp", Name: "Claude Code (ACP)", Command: "claude", ACPArgs: []string{"--experimental-acp"},
			Notes: "ACP mode for claude CLI; coexists with built-in `claude` headless spec"},
		{ID: "qwen", Name: "Qwen Code", Command: "qwen", ACPArgs: []string{"--acp"}},
		{ID: "codex-acp", Name: "Codex (ACP)", Command: "codex", ACPArgs: []string{},
			Notes: "codex-acp bridge speaks ACP by default; coexists with built-in `codex` headless spec"},
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
		{ID: "hermes-acp", Name: "Hermes Agent (ACP)", Command: "hermes", ACPArgs: []string{"acp"},
			Notes: "ACP mode for hermes CLI; coexists with built-in `hermes` headless spec"},
		{ID: "snow", Name: "Snow CLI", Command: "snow", ACPArgs: []string{"--acp"}},
	}
}

// ACPSpecs converts the ACP catalog into Spec entries suitable for merging
// into DefaultSpecs. Each spec is marked with LaunchACP so the runner picks
// the JSON-RPC code path; HeadlessArgs is intentionally left empty.
func ACPSpecs() []Spec {
	catalog := ACPCatalog()
	specs := make([]Spec, 0, len(catalog))
	for _, backend := range catalog {
		notes := backend.Notes
		if notes == "" {
			notes = "ACP backend; spawned with " + backend.Command + " " + joinACPArgs(backend.ACPArgs)
		}
		specs = append(specs, withBuiltinSources(Spec{
			ID:                    backend.ID,
			Commands:              []string{backend.Command},
			VersionArgs:           []string{"--version"},
			LaunchMode:            LaunchACP,
			HeadlessMode:          backend.Command + " " + joinACPArgs(backend.ACPArgs),
			ACPArgs:               append([]string(nil), backend.ACPArgs...),
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
		}))
	}
	return specs
}

func joinACPArgs(args []string) string {
	if len(args) == 0 {
		return "(no args)"
	}
	out := ""
	for i, a := range args {
		if i > 0 {
			out += " "
		}
		out += a
	}
	return out
}
