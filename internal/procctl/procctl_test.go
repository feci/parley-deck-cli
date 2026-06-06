//go:build !windows

package procctl

import (
	"bufio"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"testing"
	"time"
)

func waitGone(t *testing.T, pid int, d time.Duration) {
	t.Helper()
	deadline := time.Now().Add(d)
	for time.Now().Before(deadline) {
		if !unixAlive(pid) {
			return
		}
		time.Sleep(20 * time.Millisecond)
	}
	t.Fatalf("pid %d still alive after %s", pid, d)
}

// A self-spawned process is attributed; tampering with any identity field is
// refused (the PID-reuse safety gate); KillGroup reaps it.
func TestAttributedSelfAndRefusals(t *testing.T) {
	cmd := exec.Command("sleep", "30")
	SetNewProcessGroup(cmd)
	if err := cmd.Start(); err != nil {
		t.Fatalf("start: %v", err)
	}
	sp := Capture(cmd, "test-marker")
	go func() { _ = cmd.Wait() }() // reap so a killed leader isn't left a zombie
	defer func() { _ = cmd.Process.Kill() }()

	if sp.PID == 0 {
		t.Fatal("no pid captured")
	}
	if sp.BootID == "" || sp.ProcStart == "" || sp.Command == "" || sp.PGID == 0 {
		t.Skipf("process-identity probe unavailable in this environment: %+v", sp)
	}
	if ok, reason := Attributed(sp); !ok {
		t.Fatalf("self process must be attributed, refused: %s", reason)
	}
	// Each tampered field must cause a refusal.
	for _, bad := range []Spawned{
		func() Spawned { b := sp; b.ProcStart = "BOGUS-START"; return b }(),
		func() Spawned { b := sp; b.BootID = "different-boot"; return b }(),
		func() Spawned { b := sp; b.PGID = sp.PGID + 99999; return b }(),
		func() Spawned { b := sp; b.Command = "/usr/bin/totally-different"; return b }(),
	} {
		if ok, _ := Attributed(bad); ok {
			t.Fatalf("tampered identity must be refused: %+v", bad)
		}
	}
	if err := KillGroup(sp); err != nil {
		t.Fatalf("kill group: %v", err)
	}
	waitGone(t, sp.PID, 3*time.Second)
}

// Killing the process group reaps a backgrounded grandchild, not just the leader.
func TestKillGroupReapsChild(t *testing.T) {
	// sh prints the backgrounded child's PID, then both sleep.
	cmd := exec.Command("sh", "-c", "sleep 30 & echo $!; sleep 30")
	SetNewProcessGroup(cmd)
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		t.Fatal(err)
	}
	if err := cmd.Start(); err != nil {
		t.Fatalf("start: %v", err)
	}
	sp := Capture(cmd, "m")
	go func() { _ = cmd.Wait() }() // reap the leader after kill
	defer func() { _ = KillGroup(sp) }()

	line, err := bufio.NewReader(stdout).ReadString('\n')
	if err != nil {
		t.Fatalf("read child pid: %v", err)
	}
	childPID, err := strconv.Atoi(strings.TrimSpace(line))
	if err != nil {
		t.Fatalf("parse child pid %q: %v", line, err)
	}
	if !unixAlive(childPID) {
		t.Fatalf("child %d should be alive before kill", childPID)
	}
	if err := KillGroup(sp); err != nil {
		t.Fatalf("kill group: %v", err)
	}
	waitGone(t, sp.PID, 3*time.Second)
	waitGone(t, childPID, 3*time.Second) // the grandchild is reaped too
}

// A dead PID is never attributed (so it is never signaled).
func TestAttributedDeadPidRefused(t *testing.T) {
	sp := Spawned{PID: 2147480000, PGID: 2147480000, BootID: CurrentBootID(), ProcStart: "x", Command: "y"}
	if ok, _ := Attributed(sp); ok {
		t.Fatal("a dead pid must not be attributed")
	}
}

// Attributed FAILS CLOSED: an incomplete recorded identity is never attributed,
// even if the PID is alive — guards against killing on partial metadata.
func TestAttributedFailsClosedOnMissingFields(t *testing.T) {
	if !SupportsDurableKill() {
		t.Skip("durable kill unsupported here")
	}
	live := os.Getpid() // a guaranteed-alive pid
	for _, sp := range []Spawned{
		{PID: live},                                                                 // nothing recorded
		{PID: live, ProcStart: "x", PGID: live, Command: "y"},                       // no boot id
		{PID: live, BootID: CurrentBootID(), PGID: live, Command: "y"},              // no start time
		{PID: live, BootID: CurrentBootID(), ProcStart: "x", Command: "y"},          // no pgid
		{PID: live, BootID: CurrentBootID(), ProcStart: "x", PGID: live},            // no command
	} {
		if ok, _ := Attributed(sp); ok {
			t.Fatalf("incomplete identity must be refused: %+v", sp)
		}
	}
}
