//go:build linux

// Tests for the bounded /proc/<pid>/cmdline poll (the execve argv-publication
// window behind hosted U2). The poll has no context parameter by design: the
// probe interface is frozen and portable, so external cancellation reaches it
// only as process death — a read error — which returns immediately instead of
// polling (pinned below). Cancellation via KillGroup lands as SIGKILL within
// the existing 1500ms grace, far beyond the 100ms bound, so no caller's
// deadline can be pushed past its previous envelope: the bound can only extend
// paths that already failed.
package procctl

import (
	"bytes"
	"errors"
	"os"
	"os/exec"
	"strconv"
	"sync"
	"testing"
	"time"
)

// fakeCmdlineReader serves scripted read results in order (last one repeats).
type fakeCmdlineReader struct {
	results []func() ([]byte, error)
	calls   int
}

func (f *fakeCmdlineReader) read() ([]byte, error) {
	i := f.calls
	if i >= len(f.results) {
		i = len(f.results) - 1
	}
	f.calls++
	return f.results[i]()
}

func empty() func() ([]byte, error) {
	return func() ([]byte, error) { return nil, nil }
}

func value(b []byte) func() ([]byte, error) {
	return func() ([]byte, error) { return b, nil }
}

func gone() func() ([]byte, error) {
	return func() ([]byte, error) { return nil, errors.New("open /proc/…/cmdline: no such file or directory") }
}

// The healthy path is untouched: one read, no sleeps.
func TestPollCmdlineHealthyFastPath(t *testing.T) {
	f := &fakeCmdlineReader{results: []func() ([]byte, error){value([]byte("sleep\x0030\x00"))}}
	begin := time.Now()
	got, ok := pollCmdline(f.read)
	if !ok {
		t.Fatal("first non-empty read must succeed")
	}
	if got != "sleep 30" {
		t.Fatalf("NUL-separated argv must parse to %q, got %q", "sleep 30", got)
	}
	if f.calls != 1 {
		t.Fatalf("healthy path must read exactly once, read %d times", f.calls)
	}
	if elapsed := time.Since(begin); elapsed > cmdlinePublishBound {
		t.Fatalf("healthy path must not wait, took %s", elapsed)
	}
}

// The evidenced defect: zero bytes for a live process, then argv publishes.
func TestPollCmdlineRecoversTransientEmpty(t *testing.T) {
	f := &fakeCmdlineReader{results: []func() ([]byte, error){
		empty(), empty(), value([]byte("sh\x00-c\x00wait\x00")),
	}}
	got, ok := pollCmdline(f.read)
	if !ok {
		t.Fatal("transient emptiness within the bound must recover")
	}
	if got != "sh -c wait" {
		t.Fatalf("unexpected parse: %q", got)
	}
	if f.calls != 3 {
		t.Fatalf("must return on the first non-empty read (no re-reads), calls=%d", f.calls)
	}
}

// A permanently empty cmdline fails closed after exactly the bound.
func TestPollCmdlinePermanentEmptyFailsClosedWithinBound(t *testing.T) {
	f := &fakeCmdlineReader{results: []func() ([]byte, error){empty()}}
	begin := time.Now()
	got, ok := pollCmdline(f.read)
	elapsed := time.Since(begin)
	if ok || got != "" {
		t.Fatalf("permanent emptiness must fail closed with empty value, got (%q,%v)", got, ok)
	}
	if elapsed < cmdlinePublishBound-15*time.Millisecond {
		t.Fatalf("must poll the full bound before refusing, gave up after %s", elapsed)
	}
	if elapsed > 5*time.Second {
		t.Fatalf("bound exceeded: %s", elapsed)
	}
	if max := int(cmdlinePublishBound/cmdlinePublishPoll) + 5; f.calls > max {
		t.Fatalf("too many attempts: %d > %d", f.calls, max)
	}
}

// A dead process (read error) stops the poll immediately — cancellation
// arrives here as process death, never as a retry loop over a gone pid.
func TestPollCmdlineDeadProcessStopsPolling(t *testing.T) {
	f := &fakeCmdlineReader{results: []func() ([]byte, error){
		empty(), empty(), gone(),
	}}
	begin := time.Now()
	got, ok := pollCmdline(f.read)
	if ok || got != "" {
		t.Fatalf("read error must fail closed, got (%q,%v)", got, ok)
	}
	if f.calls != 3 {
		t.Fatalf("must stop at the error, calls=%d", f.calls)
	}
	if elapsed := time.Since(begin); elapsed >= cmdlinePublishBound {
		t.Fatalf("dead process must not consume the bound, took %s", elapsed)
	}
}

func TestPollCmdlineImmediateError(t *testing.T) {
	f := &fakeCmdlineReader{results: []func() ([]byte, error){gone()}}
	got, ok := pollCmdline(f.read)
	if ok || got != "" {
		t.Fatalf("immediate error must fail closed, got (%q,%v)", got, ok)
	}
	if f.calls != 1 {
		t.Fatalf("must read exactly once, calls=%d", f.calls)
	}
}

// spawnCaptureAssert is the claude-1 validation shape: a just-started process
// must never record an empty Command, and must attribute immediately. `sleep`
// never execs again after Start, so its argv is stable for the comparison.
func spawnCaptureAssert(t *testing.T) {
	t.Helper()
	cmd := exec.Command("sleep", "30")
	SetNewProcessGroup(cmd)
	if err := cmd.Start(); err != nil {
		t.Fatalf("start: %v", err)
	}
	sp := Capture(cmd, "cmdline-poll")
	stop := func() { _ = cmd.Process.Kill(); _ = cmd.Wait() }
	if sp.Command == "" {
		stop()
		t.Fatalf("capture recorded an empty command for just-started pid %d", sp.PID)
	}
	if sp.BootID == "" || sp.ProcStart == "" || sp.PGID == 0 {
		stop()
		t.Fatalf("incomplete identity captured: %+v", sp)
	}
	if ok, reason := Attributed(sp); !ok {
		stop()
		t.Fatalf("just-started process must be attributed, refused: %s", reason)
	}
	stop()
}

// Post-fix guarantee over many spawns (sequential, then concurrent to contend
// the exec window): Capture never records an empty Command. Pre-fix this was
// the probabilistic U2 race; with the poll it is structural — the only ways
// out are publication within the bound or a read error, both pinned above.
func TestLinuxCaptureOfStartedProcessRecordsCommandEveryTime(t *testing.T) {
	for i := 0; i < 40; i++ {
		spawnCaptureAssert(t)
	}
	var wg sync.WaitGroup
	errs := make(chan error, 4)
	for w := 0; w < 4; w++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for i := 0; i < 20; i++ {
				cmd := exec.Command("sleep", "30")
				SetNewProcessGroup(cmd)
				if err := cmd.Start(); err != nil {
					errs <- err
					return
				}
				sp := CaptureByPID(cmd.Process.Pid, "cmdline-poll")
				_ = cmd.Process.Kill()
				_ = cmd.Wait()
				if sp.Command == "" {
					errs <- errors.New("empty command recorded")
					return
				}
			}
		}()
	}
	wg.Wait()
	close(errs)
	for err := range errs {
		t.Fatal(err)
	}
}

// An unreaped killed child is a zombie: alive for signal-0, stat/pgid still
// readable (so every facet before Command still matches), but its cmdline
// reads as PERMANENTLY empty — the deterministic real-/proc stand-in for a
// command that never publishes. The probe must spend the bound polling and
// then fail closed, and Attributed must refuse at the live-command facet
// exactly as before the poll existed.
func TestLinuxZombieCmdlinePollsBoundThenFailsClosed(t *testing.T) {
	cmd := exec.Command("sleep", "30")
	SetNewProcessGroup(cmd)
	if err := cmd.Start(); err != nil {
		t.Fatalf("start: %v", err)
	}
	sp := Capture(cmd, "zombie-poll")
	t.Cleanup(func() { _ = cmd.Process.Kill(); _ = cmd.Wait() }) // reap at exit
	if sp.Command == "" {
		t.Skipf("identity probe unavailable in this environment: %+v", sp)
	}
	if ok, reason := Attributed(sp); !ok {
		t.Fatalf("live process must attribute before the kill, refused: %s", reason)
	}
	if err := cmd.Process.Kill(); err != nil {
		t.Fatalf("kill: %v", err)
	}
	// Wait for the zombie state (no Wait: reaping would remove /proc/<pid>).
	// Signal-0 answers for zombies too, so parse the stat state field: the
	// char after the closing paren of comm must be 'Z'.
	zombieDeadline := time.Now().Add(3 * time.Second)
	for time.Now().Before(zombieDeadline) {
		data, err := os.ReadFile("/proc/" + strconv.Itoa(sp.PID) + "/stat")
		if err == nil {
			if i := bytes.LastIndexByte(data, ')'); i >= 0 && i+2 < len(data) && data[i+2] == 'Z' {
				break
			}
		}
		time.Sleep(10 * time.Millisecond)
	}
	begin := time.Now()
	got, ok := linuxProbe{}.command(sp.PID)
	elapsed := time.Since(begin)
	if ok || got != "" {
		t.Fatalf("zombie cmdline must stay unreadable, got (%q,%v)", got, ok)
	}
	if elapsed < cmdlinePublishBound-15*time.Millisecond {
		t.Fatalf("poll must engage for the bound on real /proc before refusing, took %s", elapsed)
	}
	if elapsed > 5*time.Second {
		t.Fatalf("bound exceeded on real /proc: %s", elapsed)
	}
	if ok, reason := Attributed(sp); ok {
		t.Fatal("a zombie (dead, unreaped) must never be attributed")
	} else if reason != "cannot read live command" {
		t.Fatalf("refusal must stay at the live-command facet, got %q", reason)
	}
}

// The poll never weakens the comparison: a value it returns is still compared
// exactly, and earlier facets refuse first — PID reuse stays rejected.
func TestLinuxAttributedRefusesIdentityChangesUnderPollingProbe(t *testing.T) {
	cmd := exec.Command("sleep", "30")
	SetNewProcessGroup(cmd)
	if err := cmd.Start(); err != nil {
		t.Fatalf("start: %v", err)
	}
	sp := Capture(cmd, "strict-compare")
	t.Cleanup(func() { _ = cmd.Process.Kill(); _ = cmd.Wait() })
	if sp.Command == "" {
		t.Skipf("identity probe unavailable in this environment: %+v", sp)
	}
	bad := sp
	bad.Command = "/definitely/not/the/recorded/command"
	if ok, reason := Attributed(bad); ok || reason != "command mismatch (pid was reused)" {
		t.Fatalf("changed command must refuse at the command facet, got (%v,%q)", ok, reason)
	}
	bad = sp
	bad.ProcStart = "1"
	if ok, reason := Attributed(bad); ok || reason != "process start time mismatch (pid was reused)" {
		t.Fatalf("changed start time must refuse first, got (%v,%q)", ok, reason)
	}
}

// A reaped pid short-circuits: the ENOENT read stops the poll at once, so
// capture of a gone pid stays prompt (cancellation arrives as death).
func TestLinuxCaptureByPIDOfReapedProcessIsPrompt(t *testing.T) {
	cmd := exec.Command("sleep", "30")
	SetNewProcessGroup(cmd)
	if err := cmd.Start(); err != nil {
		t.Fatalf("start: %v", err)
	}
	pid := cmd.Process.Pid
	_ = cmd.Process.Kill()
	_ = cmd.Wait() // reap: /proc/<pid> disappears
	begin := time.Now()
	sp := CaptureByPID(pid, "reaped")
	if elapsed := time.Since(begin); elapsed >= cmdlinePublishBound {
		t.Fatalf("capture of a reaped pid must not wait the bound, took %s", elapsed)
	}
	if ok, _ := Attributed(sp); ok {
		t.Fatalf("a reaped pid must never be attributed: %+v", sp)
	}
}
