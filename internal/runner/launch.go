package runner

import (
	"bytes"
	"context"
	"errors"
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
	evidence, err := beginLaunch(ctx, root, "one-shot", agent)
	if err != nil {
		return nil, nil, err
	}
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
	prompt, protocolContext, contextErr := prepareProtocolPrompt(opts.Root, opts.Prompt, info)
	info.Context = protocolContext
	ctx = WithLaunchInfo(ctx, info)
	if contextErr != nil {
		attempt, err := beginLaunch(ctx, opts.Root, "manual", opts.Agent)
		if err != nil {
			return record, err
		}
		if err := attempt.finish(contextErr, ctx.Err(), nil); err != nil {
			return record, err
		}
		return record, contextErr
	}
	cmd, cleanup, err := trackedCommandFor(ctx, opts.Root, opts.Agent, prompt)
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
