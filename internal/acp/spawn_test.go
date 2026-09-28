package acp

import (
	"bufio"
	"context"
	"errors"
	"io"
	"os"
	"strings"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

type byteObserver struct{ bytes atomic.Int64 }

func (b *byteObserver) Write(data []byte) (int, error) {
	b.bytes.Add(int64(len(data)))
	return len(data), nil
}

func TestSpawnObservesAllStderrAndRealExit(t *testing.T) {
	observer := &byteObserver{}
	p, err := Spawn(context.Background(), SpawnOptions{Command: os.Args[0],
		Args: []string{"-test.run=TestStderrObserverChild", "--", "stderr-observer-child"}, StderrObserver: observer})
	if err != nil {
		t.Fatal(err)
	}
	if p.ExitCode() != nil {
		t.Fatal("invented exit before wait")
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := p.Stop(ctx); err == nil {
		t.Fatal("nonzero child passed")
	}
	if p.ExitCode() == nil || *p.ExitCode() != 7 {
		t.Fatalf("exit: %v", p.ExitCode())
	}
	if observer.bytes.Load() != 16384 {
		t.Fatalf("observed bytes: %d", observer.bytes.Load())
	}
	if len(p.Stderr()) != 8192 {
		t.Fatalf("ring size: %d", len(p.Stderr()))
	}
}

func TestStderrObserverChild(t *testing.T) {
	if len(os.Args) == 0 || os.Args[len(os.Args)-1] != "stderr-observer-child" {
		return
	}
	_, _ = os.Stderr.WriteString(strings.Repeat("x", 16384))
	os.Exit(7)
}

// gatingObserver parks the stderr copier inside its first Write so a test
// can control exactly when the drain resumes. bytes counts every byte the
// copier actually delivered, which is what a Wait-before-drain reap closes
// off.
type gatingObserver struct {
	first   chan struct{} // closed when the copier delivers its first chunk
	release chan struct{} // closed by the test to let the copier proceed
	bytes   atomic.Int64
	started sync.Once
}

func (g *gatingObserver) Write(data []byte) (int, error) {
	g.bytes.Add(int64(len(data)))
	g.started.Do(func() { close(g.first) })
	<-g.release
	return len(data), nil
}

// spawnGatedChild launches the re-executed child that writes 8192 bytes of
// stderr, waits for a go signal on stdin, writes 8192 more, announces READY
// on stdout and exits 7. READY proves all 16384 bytes are already in the
// pipes and the child is on its way out, so a reap that closes the stderr
// pipe mid-drain is the only way bytes can go missing.
//
// The 1 KiB in-flight volume (§D.7 a) sits inside every observed pipe
// capacity, including Windows' ~4 KiB anonymous-pipe buffer — the historical
// 16 KiB volume stalled the child before READY there (row 17).
// spawnGatedChild launches the §D.7 (a) capacity-independent ordering probe:
// the child's total stderr is 1 KiB, so every observed pipe capacity (Windows
// ~4 KiB included) holds the in-flight bytes and the READY handshake cannot
// stall on a full pipe. The copier still parks inside its first Write, so
// the drain-ordering invariant under test — Stop/Wait must not reap before
// the drain completes — is unchanged.
func spawnGatedChild(t *testing.T, observer *gatingObserver) *Process {
	t.Helper()
	p, err := Spawn(context.Background(), SpawnOptions{Command: os.Args[0],
		Args:           []string{"-test.run=TestGatedStderrChild", "--", "gated-stderr-child"},
		StderrObserver: observer})
	if err != nil {
		t.Fatal(err)
	}
	<-observer.first
	if _, err := io.WriteString(p.Stdin(), "go\n"); err != nil {
		t.Fatal(err)
	}
	// The READY read is bounded: on a platform whose pipe buffer cannot hold
	// the in-flight bytes the child stalls before READY, and an unguarded
	// ReadString would surface that only as the binary-wide global timeout.
	// The guard fails at a named site instead (D1). The parked reader
	// goroutine only outlives this call on the already-failing path.
	type readyResult struct {
		line string
		err  error
	}
	ready := make(chan readyResult, 1)
	go func() {
		line, err := bufio.NewReader(p.Stdout()).ReadString('\n')
		ready <- readyResult{line, err}
	}()
	select {
	case r := <-ready:
		if r.err != nil || r.line != "READY\n" {
			t.Fatalf("child ready handshake failed: line=%q err=%v", r.line, r.err)
		}
	case <-time.After(10 * time.Second):
		t.Fatal("child READY handshake did not complete within 10s: child stalled before READY (pipe capacity below in-flight bytes) or died without announcing")
	}
	return p
}

func assertGatedChildDrained(t *testing.T, p *Process, observer *gatingObserver) {
	t.Helper()
	if p.ExitCode() == nil || *p.ExitCode() != 7 {
		t.Fatalf("exit: %v", p.ExitCode())
	}
	if got := observer.bytes.Load(); got != 1024 {
		t.Fatalf("observed bytes: %d", got)
	}
	if len(p.Stderr()) != 1024 {
		t.Fatalf("ring size: %d", len(p.Stderr()))
	}
}

// TestSpawnStopDrainsStderrBeforeReaping pins the Stop ordering: the stderr
// copier must finish before cmd.Wait, because Wait closes the pipe. The
// sleeps only give a pre-fix Stop time to reap the exited child and close
// the pipe under the parked copier; the fixed ordering is correct regardless
// of when the release lands.
func TestSpawnStopDrainsStderrBeforeReaping(t *testing.T) {
	observer := &gatingObserver{first: make(chan struct{}), release: make(chan struct{})}
	p := spawnGatedChild(t, observer)
	time.Sleep(500 * time.Millisecond)
	stopped := make(chan error, 1)
	go func() { stopped <- p.Stop(context.Background()) }()
	time.Sleep(500 * time.Millisecond)
	close(observer.release)
	select {
	case err := <-stopped:
		if err == nil {
			t.Fatal("nonzero child passed")
		}
	case <-time.After(10 * time.Second):
		t.Fatal("Stop did not return after the copier was released")
	}
	assertGatedChildDrained(t, p, observer)
}

// TestSpawnWaitDrainsStderrBeforeReaping pins the same ordering for Wait.
func TestSpawnWaitDrainsStderrBeforeReaping(t *testing.T) {
	observer := &gatingObserver{first: make(chan struct{}), release: make(chan struct{})}
	p := spawnGatedChild(t, observer)
	time.Sleep(500 * time.Millisecond)
	waited := make(chan error, 1)
	go func() { waited <- p.Wait() }()
	time.Sleep(500 * time.Millisecond)
	close(observer.release)
	select {
	case err := <-waited:
		if err == nil {
			t.Fatal("nonzero child passed")
		}
	case <-time.After(10 * time.Second):
		t.Fatal("Wait did not return after the copier was released")
	}
	assertGatedChildDrained(t, p, observer)
}

// TestSpawnStopTimeoutIsBoundedAfterKill guards the restructured shutdown:
// once the deadline fires, Stop kills the group, interrupts the copier and
// reaps, then returns the context error instead of waiting on a child that
// ignores shutdown.
func TestSpawnStopTimeoutIsBoundedAfterKill(t *testing.T) {
	p, err := Spawn(context.Background(), SpawnOptions{Command: os.Args[0],
		Args: []string{"-test.run=TestParkingStderrChild", "--", "parking-stderr-child"}})
	if err != nil {
		t.Fatal(err)
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
		t.Fatal("Stop did not return after the kill deadline")
	}
	if p.ExitCode() == nil {
		t.Fatal("exit code unavailable after killed shutdown")
	}
}

// TestGatedStderrChild is re-executed by the drain-ordering tests: half the
// stderr payload, a go signal on stdin, the second half, READY on stdout,
// exit 7.
func TestGatedStderrChild(t *testing.T) {
	if len(os.Args) == 0 || os.Args[len(os.Args)-1] != "gated-stderr-child" {
		return
	}
	// §D.7 (a) split: the gated drain-ordering probe uses a ≤2 KiB total
	// payload so it is pipe-capacity-INDEPENDENT — Windows anonymous pipes
	// buffer ~4 KiB, where the historical 16 KiB in-flight volume stalled the
	// child before READY (row 17). The capacity-dependent 16384-write volume
	// stays in the UNGATED ring-cap test (TestSpawnObservesAllStderrAndRealExit).
	if _, err := os.Stderr.WriteString(strings.Repeat("x", 512)); err != nil {
		os.Exit(11)
	}
	if _, err := bufio.NewReader(os.Stdin).ReadString('\n'); err != nil {
		os.Exit(12)
	}
	if _, err := os.Stderr.WriteString(strings.Repeat("y", 512)); err != nil {
		os.Exit(13)
	}
	if _, err := os.Stdout.WriteString("READY\n"); err != nil {
		os.Exit(14)
	}
	os.Exit(7)
}

// TestParkingStderrChild is re-executed by the bounded-shutdown test: it
// stays alive holding the stderr pipe open until killed.
func TestParkingStderrChild(t *testing.T) {
	if len(os.Args) == 0 || os.Args[len(os.Args)-1] != "parking-stderr-child" {
		return
	}
	time.Sleep(2 * time.Minute)
	os.Exit(0)
}
