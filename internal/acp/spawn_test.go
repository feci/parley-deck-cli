package acp

import (
	"context"
	"os"
	"strings"
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
