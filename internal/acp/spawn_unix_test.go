//go:build !windows

package acp

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"syscall"
	"testing"
	"time"
)

// TestSpawnStopIsBoundedWhenStderrWriterEscapesGroup discriminates the
// ctx-branch p.stderr.Close() line in Stop, which no other test covers:
// the parking child of TestSpawnStopTimeoutIsBoundedAfterKill dies with
// the group kill, so the copier reaches EOF whether or not the read end
// was ever closed. Here the child spawns a grandchild that detaches into
// its own session (its own Setsid) and parks holding the inherited stderr
// pipe, so the group kill cannot reach the only remaining writer and
// closing the read end is the sole interrupt for the parked copier:
// without that line Stop never returns and the guard below fails; with it
// Stop returns the context error in bounded time and still reaps the
// killed child. Unix-only by construction — the escape needs POSIX
// process groups and the interrupt needs a pollable pipe fd (a blocked
// Read unblocks when the read end closes); Windows anonymous pipes are
// not pollable, Windows is owner-blocked, and this test has never run
// there, so it implies nothing about Windows behavior.
func TestSpawnStopIsBoundedWhenStderrWriterEscapesGroup(t *testing.T) {
	p, err := Spawn(context.Background(), SpawnOptions{Command: os.Args[0],
		Args: []string{"-test.run=^TestEscapedStderrChild$", "--", "escaped-stderr-child"}})
	if err != nil {
		t.Fatal(err)
	}
	reader := bufio.NewReader(p.Stdout())
	pidLine, err := reader.ReadString('\n')
	if err != nil {
		t.Fatalf("grandchild pid handshake failed: line=%q err=%v", pidLine, err)
	}
	var grandPID int
	if _, err := fmt.Sscanf(pidLine, "GRAND %d", &grandPID); err != nil || grandPID <= 0 {
		t.Fatalf("grandchild pid not reported: line=%q err=%v", pidLine, err)
	}
	t.Cleanup(func() { _ = syscall.Kill(grandPID, syscall.SIGKILL) })
	if line, err := reader.ReadString('\n'); err != nil || line != "READY\n" {
		t.Fatalf("child ready handshake failed: line=%q err=%v", line, err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 300*time.Millisecond)
	defer cancel()
	done := make(chan error, 1)
	go func() { done <- p.Stop(ctx) }()
	select {
	case err := <-done:
		if !errors.Is(err, context.DeadlineExceeded) {
			t.Fatalf("stop: %v", err)
		}
	case <-time.After(10 * time.Second):
		t.Fatal("Stop did not return after interrupting the group-escaping stderr writer")
	}
	if p.ExitCode() == nil {
		t.Fatal("exit code unavailable after killed shutdown")
	}
}

// TestEscapedStderrChild is re-executed by the escaping-writer test: it
// writes a small stderr marker, spawns a grandchild that detaches into its
// own session and parks holding the inherited stderr pipe, reports the
// grandchild pid and READY, then parks until the group kill takes it.
func TestEscapedStderrChild(t *testing.T) {
	if len(os.Args) == 0 || os.Args[len(os.Args)-1] != "escaped-stderr-child" {
		return
	}
	if _, err := os.Stderr.WriteString(strings.Repeat("e", 64)); err != nil {
		os.Exit(21)
	}
	grand := exec.Command(os.Args[0], "-test.run=^TestEscapedStderrGrandchild$", "--", "escaped-stderr-grandchild")
	grand.SysProcAttr = &syscall.SysProcAttr{Setsid: true}
	grand.Stderr = os.Stderr
	if err := grand.Start(); err != nil {
		os.Exit(22)
	}
	if _, err := fmt.Fprintf(os.Stdout, "GRAND %d\nREADY\n", grand.Process.Pid); err != nil {
		os.Exit(23)
	}
	time.Sleep(2 * time.Minute)
	os.Exit(0)
}

// TestEscapedStderrGrandchild is re-executed by the escaping-writer test's
// child: it parks in its own session doing nothing but holding the
// inherited stderr pipe open.
func TestEscapedStderrGrandchild(t *testing.T) {
	if len(os.Args) == 0 || os.Args[len(os.Args)-1] != "escaped-stderr-grandchild" {
		return
	}
	time.Sleep(2 * time.Minute)
	os.Exit(0)
}
