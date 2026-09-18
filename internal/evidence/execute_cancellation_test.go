//go:build !windows

package evidence

import (
	"context"
	"fmt"
	"io"
	"net"
	"os"
	"os/signal"
	"path/filepath"
	"strings"
	"syscall"
	"testing"
	"time"
)

func TestRunCriterionCancellationBeforeStart(t *testing.T) {
	root := t.TempDir()
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	rec := RunCriterion(ctx, root, "cancelled", "echo started > started", "codex-1")
	if rec.Status != StatusFail || rec.Command.ExitCode == 0 {
		t.Fatalf("cancelled launch must fail: %+v", rec)
	}
	if _, err := os.Stat(filepath.Join(root, "started")); !os.IsNotExist(err) {
		t.Fatalf("cancelled launch ran a command or marker is unreadable: %v", err)
	}
	// A failed Start must also return without waiting for a cancellation
	// callback or publishing successful execution evidence.
	rec = RunCriterion(context.Background(), filepath.Join(root, "absent"), "start-failure", "echo started", "codex-1")
	if rec.Status != StatusFail || rec.Command.ExitCode == 0 {
		t.Fatalf("failed Start must fail the criterion: %+v", rec)
	}
}

// The helper is a grandchild of the criterion shell, ignores TERM, and keeps
// the command's output pipes open. Its live TCP connection proves readiness
// before cancellation and EOF proves that cleanup reached the descendant.
func TestRunCriterionCancellationKillsDescendant(t *testing.T) {
	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	defer listener.Close()
	if err := listener.(*net.TCPListener).SetDeadline(time.Now().Add(10 * time.Second)); err != nil {
		t.Fatal(err)
	}
	bin, err := os.Executable()
	if err != nil {
		t.Fatal(err)
	}
	quote := func(s string) string { return "'" + strings.ReplaceAll(s, "'", "'\"'\"'") + "'" }
	command := fmt.Sprintf("%s -test.run '^TestCriterionCancellationHelper$' -- parley-cancellation-helper %s & wait", quote(bin), quote(listener.Addr().String()))
	ctx, cancel := context.WithTimeout(context.Background(), 12*time.Second)
	defer cancel()
	root := t.TempDir()
	result := make(chan CriterionRecord, 1)
	go func() {
		result <- RunCriterion(ctx, root, "descendant", command, "codex-1")
	}()
	conn, err := listener.Accept()
	if err != nil {
		t.Fatalf("helper did not become ready: %v", err)
	}
	defer conn.Close()
	if err := conn.SetReadDeadline(time.Now().Add(6 * time.Second)); err != nil {
		t.Fatal(err)
	}
	ready := make([]byte, 1)
	if _, err := io.ReadFull(conn, ready); err != nil || ready[0] != 'R' {
		t.Fatalf("missing helper readiness: %q, %v", ready, err)
	}
	cancel()
	if n, err := conn.Read(ready); n != 0 || err != io.EOF {
		t.Fatalf("cancellation left the TERM-resistant descendant alive: n=%d err=%v", n, err)
	}
	select {
	case rec := <-result:
		if rec.Status != StatusFail || rec.Command.ExitCode == 0 {
			t.Fatalf("cancelled criterion must be failed evidence: %+v", rec)
		}
	case <-time.After(3 * time.Second):
		t.Fatal("criterion did not finish after descendant cleanup")
	}
}

func TestCriterionCancellationHelper(t *testing.T) {
	args := os.Args
	if len(args) < 3 || args[len(args)-2] != "parley-cancellation-helper" {
		return
	}
	signal.Ignore(syscall.SIGTERM)
	conn, err := net.DialTimeout("tcp", args[len(args)-1], 3*time.Second)
	if err != nil {
		os.Exit(2)
	}
	defer conn.Close()
	if err := conn.SetDeadline(time.Now().Add(15 * time.Second)); err != nil {
		os.Exit(3)
	}
	if _, err := conn.Write([]byte("R")); err != nil {
		os.Exit(4)
	}
	// Independent safety bound if the parent test fails; normal cancellation
	// must kill this process, rather than waiting for its own timeout.
	_, _ = conn.Read(make([]byte, 1))
	os.Exit(5)
}
