package runner

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"parley-deck-cli/internal/agents"
	"parley-deck-cli/internal/procctl"
	"parley-deck-cli/internal/telemetry"
)

// AgentCommand preserves exec.Cmd field configuration while owning Start/Wait
// instrumentation. Like exec.Cmd, callers must serialize lifecycle operations.
type AgentCommand struct {
	*exec.Cmd
	ctx             context.Context
	evidence        *launchEvidence
	started, waited bool
	marker          string
}

type lockedBuffer struct {
	mu sync.Mutex
	bytes.Buffer
}

func (b *lockedBuffer) Write(p []byte) (int, error) {
	b.mu.Lock()
	defer b.mu.Unlock()
	return b.Buffer.Write(p)
}

func trackedCommandFor(ctx context.Context, root string, agent agents.Discovery, prompt string) (*AgentCommand, func(), error) {
	ctx, prompt, evidence, err := beginProtocolLaunch(ctx, root, "one-shot", agent, prompt)
	if err != nil {
		return nil, nil, err
	}
	return commandForLaunch(ctx, root, agent, prompt, evidence)
}

func commandForLaunch(ctx context.Context, root string, agent agents.Discovery, prompt string, evidence *launchEvidence) (*AgentCommand, func(), error) {
	path, args, env, cleanup, err := buildAgentInvocation(root, agent, prompt)
	if err != nil {
		if cleanup != nil {
			cleanup()
		}
		if finalErr := evidence.finish(err, ctx.Err(), nil); finalErr != nil {
			return nil, nil, finalErr
		}
		return nil, nil, err
	}
	cmd := exec.CommandContext(ctx, path, args...)
	cmd.Dir = root
	cmd.WaitDelay = 2 * time.Second
	if env == nil {
		env = os.Environ()
	}
	marker := evidence.invocation.ID
	cmd.Env = append(cleanParticipantEnv(agent.Adapter(), env), procctl.MarkerEnv(evidence.info.RunID, agent.ID, marker)...)
	procctl.SetNewProcessGroup(cmd)
	cmd.Cancel = func() error {
		if cmd.Process == nil {
			return os.ErrProcessDone
		}
		return procctl.KillGroup(procctl.CaptureByPID(cmd.Process.Pid, marker))
	}
	if agent.PromptMode == agents.PromptStdin {
		cmd.Stdin = strings.NewReader(prompt)
	}
	tracked := &AgentCommand{Cmd: cmd, ctx: ctx, evidence: evidence, marker: marker}
	return tracked, func() {
		if !tracked.waited {
			if tracked.started {
				_ = cmd.Cancel()
				_ = tracked.Wait()
			} else {
				_ = evidence.finish(errors.New("command abandoned before start"), nil, nil)
			}
		}
		if cleanup != nil {
			cleanup()
		}
	}, nil
}

func (c *AgentCommand) Start() error {
	if c.started || c.waited {
		return errors.New("agent command already started or consumed")
	}
	out, errOut := c.Stdout, c.Stderr
	if out == nil {
		out = io.Discard
	}
	if errOut == nil {
		errOut = io.Discard
	}
	c.Stdout = io.MultiWriter(c.evidence.collector.Writer("stdout"), out)
	c.Stderr = io.MultiWriter(c.evidence.collector.Writer("stderr"), errOut)
	if err := c.Cmd.Start(); err != nil {
		c.waited = true
		if recordErr := c.evidence.finish(err, c.ctx.Err(), nil); recordErr != nil {
			return recordErr
		}
		return err
	}
	c.started = true
	if err := c.evidence.started(c.Process.Pid); err != nil {
		_ = c.Cancel()
		_ = c.Cmd.Wait()
		c.waited = true
		if recordErr := c.evidence.finish(err, c.ctx.Err(), commandExitCode(c.Cmd)); recordErr != nil {
			return recordErr
		}
		return err
	}
	return nil
}

func commandExitCode(cmd *exec.Cmd) *int {
	if cmd.ProcessState == nil {
		return nil
	}
	code := cmd.ProcessState.ExitCode()
	return &code
}

func (c *AgentCommand) Wait() error {
	if !c.started || c.waited {
		return errors.New("agent command has no unconsumed start")
	}
	err := c.Cmd.Wait()
	c.waited = true
	if recordErr := c.evidence.finish(err, c.ctx.Err(), commandExitCode(c.Cmd)); recordErr != nil {
		return recordErr
	}
	return err
}

func (c *AgentCommand) Run() error {
	if err := c.Start(); err != nil {
		return err
	}
	return c.Wait()
}

func (c *AgentCommand) Output() ([]byte, error) {
	if c.Stdout != nil {
		return nil, errors.New("exec: Stdout already set")
	}
	var out bytes.Buffer
	c.Stdout = &out
	err := c.Run()
	return out.Bytes(), err
}

func (c *AgentCommand) CombinedOutput() ([]byte, error) {
	if c.Stdout != nil || c.Stderr != nil {
		return nil, errors.New("exec: Stdout or Stderr already set")
	}
	// Separate stream writers are copied concurrently by os/exec. A shared
	// buffer must be synchronized even though Cmd normally recognizes equality.
	var out lockedBuffer
	c.Stdout, c.Stderr = &out, &out
	err := c.Run()
	return out.Bytes(), err
}

type ExecOptions struct {
	Root    string
	Agent   agents.Discovery
	Prompt  string
	Timeout time.Duration
	Info    LaunchInfo
}

// RunMeasured is the one-shot manual-facilitation entrypoint. Raw streams stay
// in the invocation's private directory; the return value is content-free.
func RunMeasured(parent context.Context, opts ExecOptions) (record telemetry.Record, returnedErr error) {
	if agents.LaunchModeOrDefault(opts.Agent.LaunchMode) != agents.LaunchHeadless {
		return record, errors.New("agents exec requires a headless agent; use the protocol runner for ACP or interactive handoff")
	}
	ctx, cancel := context.WithTimeout(parent, timeoutForAgent(opts.Timeout, opts.Agent))
	defer cancel()
	info := opts.Info
	observer := info.Observe
	info.Observe = func(r telemetry.Record) {
		record = r
		if observer != nil {
			observer(r)
		}
	}
	ctx = WithLaunchInfo(ctx, info)
	cmd, cleanup, err := trackedCommandFor(ctx, opts.Root, opts.Agent, opts.Prompt)
	if cleanup != nil {
		defer cleanup()
	}
	if err != nil {
		return record, err
	}
	for _, stream := range []string{"stdout", "stderr"} {
		file, err := openPrivateLog(filepath.Join(cmd.evidence.invocation.Dir, stream+".log"))
		if err != nil {
			if finalErr := cmd.evidence.finish(err, nil, nil); finalErr != nil {
				err = finalErr
			}
			return record, err
		}
		defer file.Close()
		if stream == "stdout" {
			cmd.Stdout = file
		} else {
			cmd.Stderr = file
		}
	}
	returnErr := cmd.Run()
	return record, returnErr
}

func interactivePlaceholder(agent agents.Discovery) string {
	placeholder := ""
	switch agents.InteractivePromptModeOrDefault(agent.InteractivePromptMode) {
	case agents.InteractivePromptFile:
		placeholder = "{prompt_path}"
	case agents.InteractivePromptArg:
		placeholder = "{prompt}"
	}
	return placeholder
}

// ValidateInteractiveDelivery is shared by selection and the actual spawn gate.
func ValidateInteractiveDelivery(agent agents.Discovery) error {
	if agents.LaunchModeOrDefault(agent.LaunchMode) != agents.LaunchInteractive {
		return errors.New("terminal spawning requires interactive launch mode")
	}
	placeholder := interactivePlaceholder(agent)
	for _, arg := range agent.InteractiveArgs {
		if placeholder != "" && strings.Contains(arg, placeholder) {
			return nil
		}
	}
	return &protocolContextError{reason: "interactive-prompt-delivery-unconfigured"}
}

// RunInteractive measures the process separately from a previously printed
// handoff. It renders a fresh prompt at the process boundary and passes real
// terminal descriptors directly; wrapping them in a writer would turn them
// into pipes and break terminal detection and interaction in the child.
// Stream usage remains explicitly unobserved. The child owns its process group
// so cancellation reaps its descendants and restores the caller's terminal.
func RunInteractive(parent context.Context, root string, agent agents.Discovery, prompt, targetPath string, stdin, stdout, stderr *os.File) (returnedErr error) {
	ctx, cancel := context.WithTimeout(parent, time.Duration(agents.InteractiveTimeoutMSOrDefault(agent.InteractiveTimeoutMS))*time.Millisecond)
	defer cancel()
	if agents.LaunchModeOrDefault(agent.LaunchMode) != agents.LaunchInteractive {
		return errors.New("terminal spawning requires interactive launch mode")
	}
	placeholder := interactivePlaceholder(agent)
	if err := ValidateInteractiveDelivery(agent); err != nil {
		info, _ := ctx.Value(launchInfoKey{}).(LaunchInfo)
		info.Context = telemetry.Context{Mode: "refused", FallbackReason: telemetry.String("interactive-prompt-delivery-unconfigured")}
		evidence, err := beginLaunch(WithLaunchInfo(ctx, info), root, "interactive", agent)
		if err != nil {
			return err
		}
		refusal := &protocolContextError{reason: "interactive-prompt-delivery-unconfigured"}
		if err := evidence.finish(refusal, ctx.Err(), nil); err != nil {
			return err
		}
		return refusal
	}
	ctx, prepared, evidence, err := beginProtocolLaunch(ctx, root, "interactive", agent, prompt)
	if err != nil {
		return err
	}
	evidence.directTerminal = true
	var cmd *exec.Cmd
	defer func() {
		var code *int
		if cmd != nil {
			code = commandExitCode(cmd)
		}
		if err := evidence.finish(returnedErr, ctx.Err(), code); err != nil {
			returnedErr = err
		}
	}()
	if stdin == nil || stdout == nil || stderr == nil {
		return errors.New("interactive launch requires input, output and error descriptors")
	}
	promptPath := filepath.Join(evidence.invocation.Dir, "interactive-prompt.md")
	if err := writeHandoffPrompt(promptPath, []byte(prepared)); err != nil {
		return err
	}
	command := strings.TrimSpace(agent.InteractiveCommand)
	if command == "" {
		command = agent.Path
	}
	args := ExpandInteractiveArgs(agent.InteractiveArgs, root, promptPath, targetPath)
	if placeholder == "{prompt}" {
		for i := range args {
			args[i] = strings.ReplaceAll(args[i], "{prompt}", prepared)
			// Conservative per-element bound, below Linux's common 128 KiB
			// limit; total argv/env can still cause a recorded failed start.
			if len(args[i]) > 120<<10 {
				return &protocolContextError{reason: "interactive-prompt-argument-too-large-use-file"}
			}
		}
	}
	cmd = exec.CommandContext(ctx, command, args...)
	cmd.WaitDelay = 2 * time.Second
	cmd.Dir = root
	cmd.Env = append(cleanParticipantEnv(agent.Adapter(), os.Environ()), procctl.MarkerEnv(evidence.info.RunID, agent.ID, evidence.invocation.ID)...)
	cmd.Stdin, cmd.Stdout, cmd.Stderr = stdin, stdout, stderr
	restore, err := procctl.AttachTerminal(cmd, stdin)
	if err != nil {
		return err
	}
	defer func() {
		if err := restore(); err != nil {
			returnedErr = errors.Join(returnedErr, fmt.Errorf("cannot restore terminal foreground process group: %w", err))
		}
	}()
	// The child PID is its newly created process-group ID even if it exits
	// before an identity probe can run. Never inherit the parent's group here.
	cmd.Cancel = func() error {
		if cmd.Process == nil {
			return os.ErrProcessDone
		}
		return procctl.KillGroup(procctl.Spawned{PID: cmd.Process.Pid, PGID: cmd.Process.Pid})
	}
	if err := cmd.Start(); err != nil {
		return err
	}
	if err := evidence.started(cmd.Process.Pid); err != nil {
		_ = cmd.Cancel()
		_ = cmd.Wait()
		return err
	}
	return cmd.Wait()
}

// ProbeCommandFor instruments a capability/readiness probe that can run before
// a protocol workspace exists. It carries no task protocol and cannot attest
// one. Protocol task callers must use CommandFor instead.
func ProbeCommandFor(ctx context.Context, root string, agent agents.Discovery, prompt string) (*AgentCommand, func(), error) {
	info, _ := ctx.Value(launchInfoKey{}).(LaunchInfo)
	if info.Phase != "preflight" && info.Phase != "runtime-probe" {
		return nil, nil, errors.New("probe command requires a readiness or runtime-probe phase")
	}
	info.Context = telemetry.Context{Mode: "probe-only", FallbackReason: telemetry.String("no-protocol-task")}
	ctx = WithLaunchInfo(ctx, info)
	evidence, err := beginLaunch(ctx, root, "probe", agent)
	if err != nil {
		return nil, nil, err
	}
	return commandForLaunch(ctx, root, agent, prompt, evidence)
}
