package acp

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os/exec"
	"runtime"
	"sync"
)

// SpawnOptions describes the child process to launch.
type SpawnOptions struct {
	// Command is the resolved binary path (output of exec.LookPath).
	Command string
	// Args are the arguments appended after Command, including ACPArgs.
	Args []string
	// WorkingDir becomes the spawned process's cwd.
	WorkingDir string
	// Env is the full environment for the child. Pass empty for os.Environ.
	Env []string
}

// Process is a spawned ACP child plus the wired stdio streams.
type Process struct {
	cmd       *exec.Cmd
	stdin     io.WriteCloser
	stdout    io.ReadCloser
	stderr    io.ReadCloser
	stderrBuf *RingBuffer

	wg sync.WaitGroup
}

// Spawn launches the child with stdio piped. The caller must Stop() when
// done to clean up the process and goroutines.
func Spawn(ctx context.Context, opts SpawnOptions) (*Process, error) {
	if opts.Command == "" {
		return nil, errors.New("acp: SpawnOptions.Command is required")
	}
	cmd := exec.CommandContext(ctx, opts.Command, opts.Args...)
	cmd.Dir = opts.WorkingDir
	if len(opts.Env) > 0 {
		cmd.Env = opts.Env
	}
	setSysProcAttr(cmd)

	stdin, err := cmd.StdinPipe()
	if err != nil {
		return nil, fmt.Errorf("acp: stdin pipe: %w", err)
	}
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		_ = stdin.Close()
		return nil, fmt.Errorf("acp: stdout pipe: %w", err)
	}
	stderr, err := cmd.StderrPipe()
	if err != nil {
		_ = stdin.Close()
		_ = stdout.Close()
		return nil, fmt.Errorf("acp: stderr pipe: %w", err)
	}

	if err := cmd.Start(); err != nil {
		_ = stdin.Close()
		_ = stdout.Close()
		_ = stderr.Close()
		return nil, fmt.Errorf("acp: spawn %s: %w", opts.Command, err)
	}

	p := &Process{
		cmd:       cmd,
		stdin:     stdin,
		stdout:    stdout,
		stderr:    stderr,
		stderrBuf: NewRingBuffer(8 * 1024),
	}
	p.wg.Add(1)
	go func() {
		defer p.wg.Done()
		_, _ = io.Copy(p.stderrBuf, stderr)
	}()
	return p, nil
}

// Stdin returns the pipe used to send NDJSON to the agent.
func (p *Process) Stdin() io.Writer { return p.stdin }

// Stdout returns the pipe used to receive NDJSON from the agent.
func (p *Process) Stdout() io.Reader { return p.stdout }

// Stderr returns the captured stderr ring buffer (last 8 KiB).
func (p *Process) Stderr() string { return p.stderrBuf.String() }

// PID returns the OS process id.
func (p *Process) PID() int {
	if p.cmd.Process == nil {
		return 0
	}
	return p.cmd.Process.Pid
}

// Stop closes stdin (giving the agent a chance to exit), waits up to
// ctx.Done()-or-process-exit, and then kills hard if needed.
func (p *Process) Stop(ctx context.Context) error {
	if p.stdin != nil {
		_ = p.stdin.Close()
	}
	done := make(chan error, 1)
	go func() { done <- p.cmd.Wait() }()
	select {
	case err := <-done:
		p.wg.Wait()
		return err
	case <-ctx.Done():
		if p.cmd.Process != nil {
			_ = killProcessGroup(p.cmd.Process.Pid) // reap the whole tree, not just the child
		}
		<-done
		p.wg.Wait()
		return ctx.Err()
	}
}

// Wait blocks until the process exits without sending a signal.
func (p *Process) Wait() error {
	err := p.cmd.Wait()
	p.wg.Wait()
	return err
}

// Kill sends SIGKILL to the whole process group (reaping grandchildren).
func (p *Process) Kill() error {
	if p.cmd.Process == nil {
		return nil
	}
	return killProcessGroup(p.cmd.Process.Pid)
}

// PlatformIsWindows reports whether the runtime is Windows. Exposed so
// callers can branch on detection/shell semantics.
func PlatformIsWindows() bool {
	return runtime.GOOS == "windows"
}
