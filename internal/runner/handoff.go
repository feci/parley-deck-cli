package runner

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"parley-deck-cli/internal/agents"
	"parley-deck-cli/internal/fsutil"
	"parley-deck-cli/internal/protocol"
)

const UsageCaveat = "This is a user-driven interactive handoff. Provider billing and usage accounting are determined by the provider and your account. Headless mode is programmatic execution."

type HandoffOptions struct {
	Root               string
	RunID              string
	Agent              agents.Discovery
	Prompt             string
	TargetPath         string
	ValidationContract string
	ResumeCommand      string
}

type HandoffPacket struct {
	Dir              string
	PromptPath       string
	InstructionsPath string
}

func WriteHandoffPacket(opts HandoffOptions) (HandoffPacket, error) {
	if opts.RunID == "" {
		return HandoffPacket{}, fmt.Errorf("run id is required for handoff packet")
	}
	agentDir := filepath.Join(opts.Root, protocol.DeckDir, "runs", opts.RunID, "agents", opts.Agent.ID)
	if err := fsutil.MkdirAllResilient(agentDir, 0o755); err != nil {
		return HandoffPacket{}, err
	}

	packet := HandoffPacket{
		Dir:              agentDir,
		PromptPath:       filepath.Join(agentDir, "handoff-prompt.md"),
		InstructionsPath: filepath.Join(agentDir, "handoff.md"),
	}
	if err := os.WriteFile(packet.PromptPath, []byte(opts.Prompt), 0o644); err != nil {
		return HandoffPacket{}, err
	}
	if err := os.WriteFile(packet.InstructionsPath, []byte(handoffInstructions(opts, packet)), 0o644); err != nil {
		return HandoffPacket{}, err
	}
	return packet, nil
}

func handoffInstructions(opts HandoffOptions, packet HandoffPacket) string {
	command := agents.InteractiveCommandOrDefault(opts.Agent.Spec)
	args := ExpandInteractiveArgs(opts.Agent.InteractiveArgs, opts.Root, packet.PromptPath, opts.TargetPath)
	if command == "" {
		command = "<configure interactive_command>"
	}
	if len(args) > 0 {
		command += " " + strings.Join(quoteShellArgs(args), " ")
	}

	var b strings.Builder
	fmt.Fprintf(&b, "# Interactive handoff: %s\n\n", opts.Agent.ID)
	fmt.Fprintf(&b, "Agent: %s\n", opts.Agent.ID)
	fmt.Fprintf(&b, "Launch mode: %s\n", agents.LaunchModeOrDefault(opts.Agent.LaunchMode))
	fmt.Fprintf(&b, "Invoke: %s\n", agents.InteractiveInvokeOrDefault(opts.Agent.InteractiveInvoke))
	fmt.Fprintf(&b, "Prompt mode: %s\n", agents.InteractivePromptModeOrDefault(opts.Agent.InteractivePromptMode))
	fmt.Fprintf(&b, "Prompt file: %s\n", packet.PromptPath)
	fmt.Fprintf(&b, "Target artifact: %s\n", opts.TargetPath)
	fmt.Fprintf(&b, "Suggested command: %s\n", command)
	if opts.ValidationContract != "" {
		fmt.Fprintf(&b, "\nValidation contract:\n%s\n", opts.ValidationContract)
	}
	fmt.Fprintf(&b, "\nUsage note:\n%s\n", UsageCaveat)
	if strings.TrimSpace(opts.Agent.InteractiveNotes) != "" {
		fmt.Fprintf(&b, "\nOperator notes:\n%s\n", opts.Agent.InteractiveNotes)
	}
	if opts.ResumeCommand != "" {
		fmt.Fprintf(&b, "\nResume:\n%s\n", opts.ResumeCommand)
	}
	return b.String()
}

func quoteShellArgs(args []string) []string {
	quoted := make([]string, len(args))
	for i, arg := range args {
		if strings.ContainsAny(arg, " \t\n\"'") {
			quoted[i] = fmt.Sprintf("%q", arg)
			continue
		}
		quoted[i] = arg
	}
	return quoted
}

func ExpandInteractiveArgs(args []string, root, promptPath, targetPath string) []string {
	out := make([]string, len(args))
	for i, arg := range args {
		arg = strings.ReplaceAll(arg, "{root}", root)
		arg = strings.ReplaceAll(arg, "{prompt_path}", promptPath)
		arg = strings.ReplaceAll(arg, "{target_path}", targetPath)
		out[i] = arg
	}
	return out
}
