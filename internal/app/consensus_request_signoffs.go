package app

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"reflect"
	"sort"
	"strconv"
	"strings"
	"time"

	"parley-deck-cli/internal/agents"
	"parley-deck-cli/internal/consensus"
	"parley-deck-cli/internal/protocol"
	"parley-deck-cli/internal/runner"
	"parley-deck-cli/internal/store"
)

func runConsensusRequestSignoffs(ctx context.Context, args []string, stdout, stderr io.Writer) int {
	fs := flag.NewFlagSet("consensus request-signoffs", flag.ContinueOnError)
	fs.SetOutput(stderr)
	root := fs.String("dir", ".", "workspace directory")
	review := fs.Bool("review", false, "use review consensus")
	participantsFlag := fs.String("participants", "", "comma-separated participant IDs to request")
	yes := fs.Bool("yes", false, "confirm hosted/non-local backend launches")
	dryRun := fs.Bool("dry-run", false, "print planned signoff requests without invoking agents")
	var modeOverrides launchModeOverrides
	fs.Var(&modeOverrides, "mode", "launch mode override: headless|interactive|manual or AGENT=MODE; repeatable")
	if err := fs.Parse(args); err != nil {
		return 2
	}
	if fs.NArg() != 1 {
		fmt.Fprintln(stderr, "usage: parley consensus request-signoffs [--dir DIR] [--review] [--participants IDS] [--yes] [--dry-run] IDEA")
		return 2
	}

	opts := requestSignoffsOptions{
		Root:            *root,
		IdeaSlug:        fs.Arg(0),
		Review:          *review,
		ParticipantsRaw: *participantsFlag,
		Yes:             *yes,
		DryRun:          *dryRun,
		ModeOverrides:   modeOverrides,
	}
	if err := requestConsensusSignoffs(ctx, opts, stdout, stderr); err != nil {
		fmt.Fprintf(stderr, "consensus request-signoffs failed: %v\n", err)
		if errors.Is(err, errRequestNeedsConfirmation) || errors.Is(err, errRequestUsage) {
			return 2
		}
		return 1
	}
	return 0
}

type requestSignoffsOptions struct {
	Root            string
	IdeaSlug        string
	Review          bool
	ParticipantsRaw string
	Yes             bool
	DryRun          bool
	ModeOverrides   launchModeOverrides
}

var (
	errRequestNeedsConfirmation = errors.New("confirmation required")
	errRequestUsage             = errors.New("invalid request")
)

func requestConsensusSignoffs(ctx context.Context, opts requestSignoffsOptions, stdout, stderr io.Writer) error {
	rootAbs, err := filepath.Abs(opts.Root)
	if err != nil {
		return err
	}
	summary, err := consensus.Status(opts.Root, opts.IdeaSlug, opts.Review)
	if err != nil {
		return err
	}
	if len(summary.Errors) > 0 {
		return fmt.Errorf("target consensus is malformed: %s", strings.Join(summary.Errors, "; "))
	}
	if summary.Triage == consensus.TriageBlocked {
		return errors.New("target consensus is blocked; resolve the BLOCK before requesting more signoffs")
	}

	targets, explicit, err := requestSignoffTargets(summary, opts.ParticipantsRaw)
	if err != nil {
		return err
	}
	if len(targets) == 0 {
		if explicit {
			return fmt.Errorf("%w: no participants selected", errRequestUsage)
		}
		fmt.Fprintf(stdout, "No missing signoffs for %s.\n", opts.IdeaSlug)
		return nil
	}

	discovered, err := discoverConfigured(ctx, opts.Root)
	if err != nil {
		return err
	}
	selected, err := requestSignoffAgents(targets, discovered)
	if err != nil {
		return err
	}
	selected, err = applyLaunchModeOverrides(selected, opts.ModeOverrides)
	if err != nil {
		return err
	}
	if err := validateLaunchModes(selected); err != nil {
		return err
	}
	hosted := hostedHeadlessAgentIDs(selected)
	requiresYes := len(hosted) > 0

	if opts.DryRun {
		printRequestSignoffsDryRun(stdout, rootAbs, summary, selected, requiresYes)
		return nil
	}
	if requiresYes && !opts.Yes {
		return fmt.Errorf("%w: selected participants use hosted/non-local backends (%s); rerun with --yes or --dry-run", errRequestNeedsConfirmation, strings.Join(hosted, ","))
	}

	successes := make([]string, 0, len(selected))
	pending := make([]string, 0)
	runID := store.NewRunID(time.Now())
	if hasNonHeadlessLaunch(selected) {
		runStore := store.New(filepath.Join(rootAbs, protocol.DeckDir, "runs", runID))
		if err := runStore.Append(store.Event{
			Time: time.Now().UTC(),
			Type: "run.created",
			Data: map[string]any{
				"idea":         opts.IdeaSlug,
				"mode":         "consensus-signoff",
				"participants": agentIDs(selected),
				"review":       opts.Review,
			},
		}); err != nil {
			return err
		}
	}
	for _, agent := range selected {
		fmt.Fprintf(stdout, "Requesting signoff from %s (%s)...\n", agent.ID, agents.LaunchModeOrDefault(agent.LaunchMode))
		before, err := consensus.Status(opts.Root, opts.IdeaSlug, opts.Review)
		if err != nil {
			return err
		}
		if len(before.Errors) > 0 {
			return fmt.Errorf("target consensus became malformed before %s: %s", agent.ID, strings.Join(before.Errors, "; "))
		}
		beforeRaw, err := os.ReadFile(before.Path)
		if err != nil {
			return err
		}

		prompt := buildConsensusSignoffPrompt(rootAbs, summary.Path, opts.Review, agent.ID)
		runResult, runErr := runSignoffAgent(ctx, rootAbs, runID, agent, prompt, summary.Path, stdout, stderr)
		if runResult.Pending {
			pending = append(pending, agent.ID)
			fmt.Fprintf(stdout, "Pending signoff from %s. Handoff: %s\n", agent.ID, runResult.InstructionsPath)
			continue
		}

		after, statusErr := consensus.Status(opts.Root, opts.IdeaSlug, opts.Review)
		if statusErr != nil {
			printPartialProgress(stdout, successes)
			return statusErr
		}
		afterRaw, err := os.ReadFile(after.Path)
		if err != nil {
			printPartialProgress(stdout, successes)
			return err
		}
		signoff, validateErr := validateRequestedSignoff(before, after, agent.ID, string(beforeRaw), string(afterRaw))
		if validateErr != nil {
			printPartialProgress(stdout, successes)
			return validateErr
		}
		if runErr != nil {
			printPartialProgress(stdout, successes)
			return fmt.Errorf("%s exited with error after appending valid signoff: %w", agent.ID, runErr)
		}

		successes = append(successes, agent.ID)
		fmt.Fprintf(stdout, "Accepted signoff from %s: %s\n", agent.ID, signoff.Status)
	}

	finalSummary, err := consensus.Status(opts.Root, opts.IdeaSlug, opts.Review)
	if err != nil {
		return err
	}
	if len(pending) > 0 {
		fmt.Fprintf(stdout, "Requested signoffs pending: %s\n", strings.Join(pending, ","))
		printConsensusSummary(stdout, finalSummary)
		return nil
	}
	fmt.Fprintf(stdout, "Requested signoffs complete: %s\n", strings.Join(successes, ","))
	printConsensusSummary(stdout, finalSummary)
	return nil
}

func requestSignoffTargets(summary consensus.Summary, raw string) ([]string, bool, error) {
	if strings.TrimSpace(raw) == "" {
		return append([]string(nil), summary.Missing...), false, nil
	}

	participants := make(map[string]bool, len(summary.Participants))
	for _, participant := range summary.Participants {
		participants[participant] = true
	}
	signed := make(map[string]bool, len(summary.Signoffs))
	for _, signoff := range summary.Signoffs {
		signed[signoff.Agent] = true
	}

	var targets []string
	seen := map[string]bool{}
	for _, part := range strings.Split(raw, ",") {
		id := strings.TrimSpace(part)
		if id == "" {
			continue
		}
		if seen[id] {
			return nil, true, fmt.Errorf("%w: duplicate participant %s", errRequestUsage, id)
		}
		if !participants[id] {
			return nil, true, fmt.Errorf("%w: unknown participant %s", errRequestUsage, id)
		}
		if signed[id] {
			return nil, true, fmt.Errorf("%w: participant %s already signed", errRequestUsage, id)
		}
		targets = append(targets, id)
		seen[id] = true
	}
	return targets, true, nil
}

func requestSignoffAgents(targets []string, discovered []agents.Discovery) ([]agents.Discovery, error) {
	byID := make(map[string]agents.Discovery, len(discovered))
	for _, result := range discovered {
		byID[result.ID] = result
	}
	selected := make([]agents.Discovery, 0, len(targets))
	for _, target := range targets {
		agent, ok := byID[target]
		if !ok {
			return nil, fmt.Errorf("%w: participant %s has no configured runner entry", errRequestUsage, target)
		}
		if !agent.Found {
			return nil, fmt.Errorf("%w: participant %s runner is not installed", errRequestUsage, target)
		}
		selected = append(selected, agent)
	}
	return selected, nil
}

type launchModeOverrides struct {
	global string
	byID   map[string]string
}

func (o *launchModeOverrides) String() string {
	if o == nil {
		return ""
	}
	var parts []string
	if o.global != "" {
		parts = append(parts, o.global)
	}
	for id, mode := range o.byID {
		parts = append(parts, id+"="+mode)
	}
	sort.Strings(parts)
	return strings.Join(parts, ",")
}

func (o *launchModeOverrides) Set(value string) error {
	value = strings.TrimSpace(value)
	if value == "" {
		return fmt.Errorf("empty launch mode override")
	}
	if o.byID == nil {
		o.byID = map[string]string{}
	}
	if id, mode, ok := strings.Cut(value, "="); ok {
		id = strings.TrimSpace(id)
		mode = strings.TrimSpace(mode)
		if id == "" {
			return fmt.Errorf("empty agent in launch mode override %q", value)
		}
		if !validLaunchMode(mode) {
			return fmt.Errorf("invalid launch mode %q", mode)
		}
		o.byID[id] = mode
		return nil
	}
	if !validLaunchMode(value) {
		return fmt.Errorf("invalid launch mode %q", value)
	}
	o.global = value
	return nil
}

func applyLaunchModeOverrides(selected []agents.Discovery, overrides launchModeOverrides) ([]agents.Discovery, error) {
	out := make([]agents.Discovery, len(selected))
	copy(out, selected)
	for i := range out {
		if overrides.global != "" {
			out[i].LaunchMode = overrides.global
		}
		if mode, ok := overrides.byID[out[i].ID]; ok {
			out[i].LaunchMode = mode
		}
	}
	for id := range overrides.byID {
		found := false
		for _, agent := range out {
			if agent.ID == id {
				found = true
				break
			}
		}
		if !found {
			return nil, fmt.Errorf("%w: launch mode override references unselected participant %s", errRequestUsage, id)
		}
	}
	return out, nil
}

func validateLaunchModes(selected []agents.Discovery) error {
	for _, agent := range selected {
		mode := agents.LaunchModeOrDefault(agent.LaunchMode)
		if !validLaunchMode(mode) {
			return fmt.Errorf("%w: %s has invalid launch_mode %q", errRequestUsage, agent.ID, agent.LaunchMode)
		}
		if mode == agents.LaunchHeadless {
			continue
		}
		promptMode := agents.InteractivePromptModeOrDefault(agent.InteractivePromptMode)
		if promptMode != agents.InteractivePromptNone && promptMode != agents.InteractivePromptFile && promptMode != agents.InteractivePromptArg {
			return fmt.Errorf("%w: %s has invalid interactive_prompt_mode %q", errRequestUsage, agent.ID, agent.InteractivePromptMode)
		}
		invoke := agents.InteractiveInvokeOrDefault(agent.InteractiveInvoke)
		if invoke != agents.InteractiveInvokePrintOnly && invoke != agents.InteractiveInvokeSpawnTTY {
			return fmt.Errorf("%w: %s has invalid interactive_invoke %q", errRequestUsage, agent.ID, agent.InteractiveInvoke)
		}
	}
	return nil
}

func validLaunchMode(mode string) bool {
	switch strings.TrimSpace(strings.ToLower(mode)) {
	case agents.LaunchHeadless, agents.LaunchInteractive, agents.LaunchManual:
		return true
	default:
		return false
	}
}

func hostedHeadlessAgentIDs(selected []agents.Discovery) []string {
	var hosted []string
	for _, agent := range selected {
		if agents.LaunchModeOrDefault(agent.LaunchMode) == agents.LaunchHeadless && requiresExplicitConfirmation(agent.ExternalBackend) {
			hosted = append(hosted, agent.ID)
		}
	}
	return hosted
}

func hasNonHeadlessLaunch(selected []agents.Discovery) bool {
	for _, agent := range selected {
		if agents.LaunchModeOrDefault(agent.LaunchMode) != agents.LaunchHeadless {
			return true
		}
	}
	return false
}

func requiresExplicitConfirmation(externalBackend string) bool {
	backend := strings.TrimSpace(strings.ToLower(externalBackend))
	return backend != "" && backend != agents.ExternalLocal
}

func printRequestSignoffsDryRun(stdout io.Writer, rootAbs string, summary consensus.Summary, selected []agents.Discovery, requiresYes bool) {
	fmt.Fprintln(stdout, "Consensus signoff request dry-run")
	fmt.Fprintf(stdout, "Target: %s\n", summary.Path)
	fmt.Fprintf(stdout, "Current status: %s\n", summary.Triage)
	fmt.Fprintf(stdout, "Participants: %s\n", strings.Join(agentIDs(selected), ","))
	fmt.Fprintf(stdout, "Requires --yes: %s\n", yesNo(requiresYes))
	fmt.Fprintln(stdout, "Launch order:")
	for i, agent := range selected {
		fmt.Fprintf(stdout, "  %d. %s backend=%s mode=%s command=%s\n", i+1, agent.ID, valueOr(agent.ExternalBackend, agents.ExternalLocal), agents.LaunchModeOrDefault(agent.LaunchMode), commandPreview(rootAbs, agent))
		if agents.LaunchModeOrDefault(agent.LaunchMode) != agents.LaunchHeadless {
			fmt.Fprintf(stdout, "     usage: %s\n", runner.UsageCaveat)
		}
	}
}

type signoffRunResult struct {
	Pending          bool
	InstructionsPath string
}

func runSignoffAgent(ctx context.Context, rootAbs, runID string, agent agents.Discovery, prompt, consensusPath string, stdout, stderr io.Writer) (signoffRunResult, error) {
	switch agents.LaunchModeOrDefault(agent.LaunchMode) {
	case agents.LaunchHeadless:
		return signoffRunResult{}, runHeadlessSignoffAgent(ctx, rootAbs, agent, prompt, stdout, stderr)
	case agents.LaunchInteractive:
		return runInteractiveSignoffAgent(ctx, rootAbs, runID, agent, prompt, consensusPath, stdout, stderr)
	case agents.LaunchManual:
		return runManualSignoffAgent(rootAbs, runID, agent, prompt, consensusPath, stdout)
	default:
		return signoffRunResult{}, fmt.Errorf("%s has invalid launch_mode %q", agent.ID, agent.LaunchMode)
	}
}

func runHeadlessSignoffAgent(ctx context.Context, rootAbs string, agent agents.Discovery, prompt string, stdout, stderr io.Writer) error {
	agentCtx, cancel := context.WithTimeout(ctx, requestSignoffTimeout(agent))
	defer cancel()

	cmd, cleanup, err := runner.CommandFor(agentCtx, rootAbs, agent, prompt)
	if cleanup != nil {
		defer cleanup()
	}
	if err != nil {
		return err
	}
	cmd.Dir = rootAbs
	cmd.Stdout = stdout
	cmd.Stderr = stderr
	if agent.PromptMode == agents.PromptStdin {
		cmd.Stdin = strings.NewReader(prompt)
	}
	if err := cmd.Run(); err != nil {
		if agentCtx.Err() != nil {
			return agentCtx.Err()
		}
		return err
	}
	return nil
}

func runInteractiveSignoffAgent(ctx context.Context, rootAbs, runID string, agent agents.Discovery, prompt, consensusPath string, stdout, stderr io.Writer) (signoffRunResult, error) {
	packet, err := writeSignoffHandoff(rootAbs, runID, agent, prompt, consensusPath)
	if err != nil {
		return signoffRunResult{}, err
	}
	fmt.Fprintf(stdout, "Interactive handoff for %s\n", agent.ID)
	fmt.Fprintf(stdout, "  Prompt: %s\n", packet.PromptPath)
	fmt.Fprintf(stdout, "  Instructions: %s\n", packet.InstructionsPath)
	fmt.Fprintf(stdout, "  Target: %s\n", consensusPath)
	fmt.Fprintf(stdout, "  Usage: %s\n", runner.UsageCaveat)
	if err := appendSignoffEvent(rootAbs, runID, "agent.handoff.started", map[string]any{
		"agent":        agent.ID,
		"artifact":     consensusPath,
		"prompt":       packet.PromptPath,
		"instructions": packet.InstructionsPath,
		"launch_mode":  agents.LaunchInteractive,
	}); err != nil {
		return signoffRunResult{}, err
	}

	if agents.InteractiveInvokeOrDefault(agent.InteractiveInvoke) == agents.InteractiveInvokeSpawnTTY {
		if err := runInteractiveTTY(ctx, rootAbs, agent, packet, consensusPath); err != nil {
			return signoffRunResult{}, err
		}
	}

	agentCtx, cancel := context.WithTimeout(ctx, requestInteractiveSignoffTimeout(agent))
	defer cancel()
	ticker := time.NewTicker(time.Duration(agents.InteractivePollMSOrDefault(agent.InteractivePollMS)) * time.Millisecond)
	defer ticker.Stop()
	for {
		if consensusContainsSignoff(consensusPath, agent.ID) {
			if err := appendSignoffEvent(rootAbs, runID, "agent.handoff.completed", map[string]any{
				"agent":       agent.ID,
				"artifact":    consensusPath,
				"launch_mode": agents.LaunchInteractive,
			}); err != nil {
				return signoffRunResult{}, err
			}
			return signoffRunResult{}, nil
		}
		select {
		case <-agentCtx.Done():
			return signoffRunResult{}, agentCtx.Err()
		case <-ticker.C:
		}
	}
}

func runManualSignoffAgent(rootAbs, runID string, agent agents.Discovery, prompt, consensusPath string, stdout io.Writer) (signoffRunResult, error) {
	packet, err := writeSignoffHandoff(rootAbs, runID, agent, prompt, consensusPath)
	if err != nil {
		return signoffRunResult{}, err
	}
	fmt.Fprintf(stdout, "Manual handoff for %s\n", agent.ID)
	fmt.Fprintf(stdout, "  Prompt: %s\n", packet.PromptPath)
	fmt.Fprintf(stdout, "  Instructions: %s\n", packet.InstructionsPath)
	fmt.Fprintf(stdout, "  Target: %s\n", consensusPath)
	fmt.Fprintf(stdout, "  Resume after the signoff is appended: parley consensus status %s\n", ideaSlugForPath(consensusPath))
	if err := appendSignoffEvent(rootAbs, runID, "agent.handoff.pending", map[string]any{
		"agent":        agent.ID,
		"artifact":     consensusPath,
		"prompt":       packet.PromptPath,
		"instructions": packet.InstructionsPath,
		"launch_mode":  agents.LaunchManual,
	}); err != nil {
		return signoffRunResult{}, err
	}
	return signoffRunResult{Pending: true, InstructionsPath: packet.InstructionsPath}, nil
}

func appendSignoffEvent(rootAbs, runID, eventType string, data map[string]any) error {
	runStore := store.New(filepath.Join(rootAbs, protocol.DeckDir, "runs", runID))
	return runStore.Append(store.Event{Time: time.Now().UTC(), Type: eventType, Data: data})
}

func writeSignoffHandoff(rootAbs, runID string, agent agents.Discovery, prompt, consensusPath string) (runner.HandoffPacket, error) {
	return runner.WriteHandoffPacket(runner.HandoffOptions{
		Root:               rootAbs,
		RunID:              runID,
		Agent:              agent,
		Prompt:             prompt,
		TargetPath:         consensusPath,
		ValidationContract: "Append exactly one signoff block for the agent and do not edit existing consensus content.",
		ResumeCommand:      fmt.Sprintf("parley consensus status %s", ideaSlugForPath(consensusPath)),
	})
}

func runInteractiveTTY(ctx context.Context, rootAbs string, agent agents.Discovery, packet runner.HandoffPacket, targetPath string) error {
	if !isTerminal(os.Stdin) || !isTerminal(os.Stdout) {
		return fmt.Errorf("%s interactive_invoke=spawn-tty requires a terminal", agent.ID)
	}
	command := strings.TrimSpace(agent.InteractiveCommand)
	if command == "" {
		command = agent.Path
	}
	if command == "" {
		command = agents.InteractiveCommandOrDefault(agent.Spec)
	}
	args := expandInteractiveArgs(agent.InteractiveArgs, rootAbs, packet.PromptPath, targetPath)
	cmd := exec.CommandContext(ctx, command, args...)
	cmd.Dir = rootAbs
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		if ctx.Err() != nil {
			return ctx.Err()
		}
		return err
	}
	return nil
}

func expandInteractiveArgs(args []string, rootAbs, promptPath, targetPath string) []string {
	out := make([]string, len(args))
	for i, arg := range args {
		arg = strings.ReplaceAll(arg, "{root}", rootAbs)
		arg = strings.ReplaceAll(arg, "{prompt_path}", promptPath)
		arg = strings.ReplaceAll(arg, "{target_path}", targetPath)
		out[i] = arg
	}
	return out
}

func consensusContainsSignoff(path, agentID string) bool {
	data, err := os.ReadFile(path)
	if err != nil {
		return false
	}
	return len(signoffHeaderAgents(string(data))) > 0 && signoffsContainAgent(string(data), agentID)
}

func signoffsContainAgent(raw, agentID string) bool {
	for _, agent := range signoffHeaderAgents(raw) {
		if agent == agentID {
			return true
		}
	}
	return false
}

func requestInteractiveSignoffTimeout(agent agents.Discovery) time.Duration {
	return time.Duration(agents.InteractiveTimeoutMSOrDefault(agent.InteractiveTimeoutMS)) * time.Millisecond
}

func isTerminal(file *os.File) bool {
	info, err := file.Stat()
	return err == nil && (info.Mode()&os.ModeCharDevice) != 0
}

func requestSignoffTimeout(agent agents.Discovery) time.Duration {
	timeoutMS := agent.TimeoutMS
	if timeoutMS <= 0 {
		timeoutMS = agents.DefaultTimeoutMS
	}
	return time.Duration(timeoutMS) * time.Millisecond
}

func validateRequestedSignoff(before, after consensus.Summary, agentID, beforeRaw, afterRaw string) (consensus.Signoff, error) {
	if len(after.Errors) > 0 {
		return consensus.Signoff{}, fmt.Errorf("target consensus is malformed after %s: %s", agentID, strings.Join(after.Errors, "; "))
	}
	if err := validateAppendOnlyContent(beforeRaw, afterRaw, agentID); err != nil {
		return consensus.Signoff{}, err
	}

	beforeByAgent := signoffsByAgent(before.Signoffs)
	afterByAgent := signoffsByAgent(after.Signoffs)
	if _, ok := beforeByAgent[agentID]; ok {
		return consensus.Signoff{}, fmt.Errorf("participant %s was already signed before invocation", agentID)
	}
	if newAgents := newSignoffAgents(beforeByAgent, afterByAgent); len(newAgents) != 1 || newAgents[0] != agentID {
		return consensus.Signoff{}, fmt.Errorf("%s appended unexpected signoff set: %s", agentID, strings.Join(newAgents, ","))
	}
	signoff, ok := afterByAgent[agentID]
	if !ok {
		return consensus.Signoff{}, fmt.Errorf("%s did not append a signoff", agentID)
	}
	for agent, beforeSignoff := range beforeByAgent {
		if agent == agentID {
			continue
		}
		afterSignoff, ok := afterByAgent[agent]
		if !ok {
			return consensus.Signoff{}, fmt.Errorf("%s removed signoff for %s", agentID, agent)
		}
		if !reflect.DeepEqual(beforeSignoff, afterSignoff) {
			return consensus.Signoff{}, fmt.Errorf("%s changed signoff for %s", agentID, agent)
		}
	}
	status, err := consensus.CanonicalStatus(signoff.Status)
	if err != nil {
		return consensus.Signoff{}, err
	}
	if status == consensus.StatusBlock {
		return consensus.Signoff{}, fmt.Errorf("%s appended BLOCK signoff", agentID)
	}
	return signoff, nil
}

func validateAppendOnlyContent(beforeRaw, afterRaw, agentID string) error {
	if !strings.HasPrefix(afterRaw, beforeRaw) {
		return fmt.Errorf("%s changed existing consensus content outside the append-only suffix", agentID)
	}
	suffix := strings.TrimSpace(strings.TrimPrefix(afterRaw, beforeRaw))
	if suffix == "" {
		return nil
	}
	headerAgents := signoffHeaderAgents(suffix)
	if len(headerAgents) != 1 {
		return fmt.Errorf("%s appended %d signoff blocks; expected exactly one", agentID, len(headerAgents))
	}
	if headerAgents[0] != agentID {
		return fmt.Errorf("%s appended signoff for %s", agentID, headerAgents[0])
	}
	return nil
}

func signoffHeaderAgents(raw string) []string {
	var agents []string
	for _, line := range strings.Split(raw, "\n") {
		line = strings.TrimSpace(line)
		if !strings.HasPrefix(line, "### Signoff:") {
			continue
		}
		rest := strings.TrimSpace(strings.TrimPrefix(line, "### Signoff:"))
		if agent, _, ok := strings.Cut(rest, "—"); ok {
			agents = append(agents, strings.TrimSpace(agent))
			continue
		}
		if idx := strings.LastIndex(rest, " - "); idx >= 0 {
			agents = append(agents, strings.TrimSpace(rest[:idx]))
			continue
		}
		agents = append(agents, strings.TrimSpace(rest))
	}
	return agents
}

func newSignoffAgents(beforeByAgent, afterByAgent map[string]consensus.Signoff) []string {
	var agents []string
	for agent := range afterByAgent {
		if _, ok := beforeByAgent[agent]; !ok {
			agents = append(agents, agent)
		}
	}
	sort.Strings(agents)
	return agents
}

func signoffsByAgent(signoffs []consensus.Signoff) map[string]consensus.Signoff {
	byAgent := make(map[string]consensus.Signoff, len(signoffs))
	for _, signoff := range signoffs {
		byAgent[signoff.Agent] = signoff
	}
	return byAgent
}

func buildConsensusSignoffPrompt(rootAbs, consensusPath string, review bool, agentID string) string {
	consensusAbs, err := filepath.Abs(consensusPath)
	if err != nil {
		consensusAbs = consensusPath
	}
	ideaDir := ideaDirFromConsensusPath(consensusAbs, review)
	rounds := roundDirsForPrompt(ideaDir, review)
	return fmt.Sprintf(`You are Parley Deck participant %s in the consensus signoff phase.

Repository root: %s
Consensus file to sign: %s

Protocol rules:
- Edit exactly one file: the consensus file above.
- Append exactly one signoff block for %s under ## Signoffs.
- Do not edit or rewrite any existing line, and do not edit any other file.
- Read enough context to make your own signoff decision.
- Use one canonical status: %s, %s, or %s.
- If you cannot accept, include required notes and, for BLOCK, a counter-proposal.
- Do not let the facilitator write your signoff.

Relevant paths:
- Cooperation protocol: %s
- Idea prompt: %s
- Discussion rounds: %s

Canonical block shape:

### Signoff: %s — %s
Status: %s
Notes: <brief note>
Counter-proposal (required if ❌): <counter-proposal if status is ❌ BLOCK>

Return a short confirmation after appending.
`, agentID,
		rootAbs,
		consensusAbs,
		agentID,
		consensus.StatusAccept,
		consensus.StatusReservations,
		consensus.StatusBlock,
		filepath.Join(rootAbs, protocol.DeckDir, "COOPERATION.md"),
		filepath.Join(ideaDir, "00-prompt.md"),
		strings.Join(rounds, ", "),
		agentID,
		time.Now().Format("2006-01-02"),
		consensus.StatusAccept)
}

func printPartialProgress(stdout io.Writer, successes []string) {
	if len(successes) == 0 {
		return
	}
	fmt.Fprintf(stdout, "Signoffs accepted before failure: %s\n", strings.Join(successes, ","))
}

func ideaDirFromConsensusPath(consensusPath string, review bool) string {
	if review {
		return filepath.Dir(filepath.Dir(consensusPath))
	}
	return filepath.Dir(consensusPath)
}

func roundDirsForPrompt(ideaDir string, review bool) []string {
	base := ideaDir
	if review {
		base = filepath.Join(ideaDir, "review")
	}
	entries, err := os.ReadDir(base)
	if err != nil {
		return []string{"none found"}
	}
	var rounds []string
	for _, entry := range entries {
		if entry.IsDir() && strings.HasPrefix(entry.Name(), "round-") {
			rounds = append(rounds, filepath.Join(base, entry.Name()))
		}
	}
	sort.Strings(rounds)
	if len(rounds) == 0 {
		return []string{"none found"}
	}
	return rounds
}

func ideaSlugForPath(consensusPath string) string {
	dir := filepath.Dir(consensusPath)
	if filepath.Base(dir) == "review" {
		return filepath.Base(filepath.Dir(dir))
	}
	return filepath.Base(dir)
}

func agentIDs(selected []agents.Discovery) []string {
	ids := make([]string, len(selected))
	for i, agent := range selected {
		ids[i] = agent.ID
	}
	return ids
}

func commandPreview(rootAbs string, agent agents.Discovery) string {
	if agents.LaunchModeOrDefault(agent.LaunchMode) != agents.LaunchHeadless {
		command := strings.TrimSpace(agent.InteractiveCommand)
		if command == "" {
			command = agents.InteractiveCommandOrDefault(agent.Spec)
		}
		parts := []string{command}
		parts = append(parts, expandInteractiveArgs(agent.InteractiveArgs, rootAbs, "<handoff-prompt>", "<target-artifact>")...)
		for i, part := range parts {
			if strings.ContainsAny(part, " \t\n") {
				parts[i] = strconv.Quote(part)
			}
		}
		return strings.Join(parts, " ")
	}
	parts := []string{agent.Path}
	for _, arg := range agent.HeadlessArgs {
		switch arg {
		case "{root}":
			parts = append(parts, rootAbs)
		case "{prompt}":
			parts = append(parts, "<signoff-prompt>")
		default:
			parts = append(parts, arg)
		}
	}
	for i, part := range parts {
		if strings.ContainsAny(part, " \t\n") {
			parts[i] = strconv.Quote(part)
		}
	}
	return strings.Join(parts, " ")
}

func yesNo(value bool) string {
	if value {
		return "yes"
	}
	return "no"
}
