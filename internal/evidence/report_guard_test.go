package evidence

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"testing"
	"time"
)

func TestReportCASChild(t *testing.T) {
	dir := os.Getenv("PARLEY_REPORT_CAS_DIR")
	if dir == "" {
		return
	}
	id := os.Getenv("PARLEY_REPORT_CAS_ID")
	original, err := os.ReadFile(ReportPath(dir))
	if err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(dir, "ready-"+id), nil, 0600); err != nil {
		t.Fatal(err)
	}
	deadline := time.Now().Add(5 * time.Second)
	for {
		if _, err := os.Stat(filepath.Join(dir, "release")); err == nil {
			break
		}
		if time.Now().After(deadline) {
			t.Fatal("parent barrier not released")
		}
		time.Sleep(5 * time.Millisecond)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	err = WithReportWriter(ctx, dir, func(w *ReportWriter) error {
		r := positiveReport()
		r.Idea = id
		_, err := w.SaveIfUnchanged(original, r)
		return err
	})
	if errors.Is(err, ErrReportChanged) {
		os.Exit(3)
	}
	if err != nil {
		t.Fatal(err)
	}
}

func TestReportCASSeparateProcesses(t *testing.T) {
	dir := t.TempDir()
	if err := Save(dir, positiveReport()); err != nil {
		t.Fatal(err)
	}
	var children []*exec.Cmd
	for i := 0; i < 2; i++ {
		cmd := exec.Command(os.Args[0], "-test.run=^TestReportCASChild$", "-test.count=1")
		cmd.Env = append(os.Environ(), "PARLEY_REPORT_CAS_DIR="+dir, fmt.Sprintf("PARLEY_REPORT_CAS_ID=writer-%d", i))
		if err := cmd.Start(); err != nil {
			t.Fatal(err)
		}
		children = append(children, cmd)
		t.Cleanup(func() { _ = cmd.Process.Kill() })
	}
	deadline := time.Now().Add(5 * time.Second)
	for {
		ready, _ := filepath.Glob(filepath.Join(dir, "ready-*"))
		if len(ready) == 2 {
			break
		}
		if time.Now().After(deadline) {
			t.Fatal("children did not snapshot the same report")
		}
		time.Sleep(5 * time.Millisecond)
	}
	if err := os.WriteFile(filepath.Join(dir, "release"), nil, 0600); err != nil {
		t.Fatal(err)
	}
	winners, refused := 0, 0
	for _, cmd := range children {
		err := cmd.Wait()
		if err == nil {
			winners++
		} else if cmd.ProcessState.ExitCode() == 3 {
			refused++
		} else {
			t.Fatalf("child failed for an unrelated reason: %v", err)
		}
	}
	if winners != 1 || refused != 1 {
		t.Fatalf("CAS writers: accepted=%d refused=%d", winners, refused)
	}
	r, err := Load(dir)
	if err != nil || (r.Idea != "writer-0" && r.Idea != "writer-1") {
		t.Fatalf("lost or corrupt winner: %+v %v", r, err)
	}
}

func TestReportGuardSharedWithOrdinarySaveAndFinalization(t *testing.T) {
	dir := t.TempDir()
	if err := Save(dir, positiveReport()); err != nil {
		t.Fatal(err)
	}
	original, _ := os.ReadFile(ReportPath(dir))
	var escaped *ReportWriter
	var finished chan error
	err := WithReportWriter(context.Background(), dir, func(w *ReportWriter) error {
		escaped = w
		started := make(chan struct{})
		finished = make(chan error, 1)
		go func() { close(started); finished <- Save(dir, positiveReport()) }()
		<-started
		select {
		case err := <-finished:
			return fmt.Errorf("Save ignored publication guard: %v", err)
		case <-time.After(60 * time.Millisecond):
		}
		// Simulate the final authorized write while the same publication guard
		// is held. The delayed Save must see complete after it acquires the lock.
		return os.WriteFile(filepath.Join(dir, "IMPLEMENTATION.md"), []byte("---\nstatus: complete\n---\n"), 0600)
	})
	if err != nil {
		t.Fatal(err)
	}
	if err := <-finished; !errors.Is(err, ErrReportFinalized) {
		t.Fatalf("delayed writer replaced completed evidence: %v", err)
	}
	current, _ := os.ReadFile(ReportPath(dir))
	if string(original) != string(current) {
		t.Fatal("final report changed")
	}
	if _, err := escaped.Save(positiveReport()); err == nil {
		t.Fatal("writer escaped callback lifetime")
	}
}

func TestReportGuardAliasAndCancellation(t *testing.T) {
	dir := t.TempDir()
	alias := filepath.Join(t.TempDir(), "alias")
	if err := os.Symlink(dir, alias); err != nil {
		t.Skipf("symlink unavailable: %v", err)
	}
	err := WithReportWriter(context.Background(), dir, func(_ *ReportWriter) error {
		ctx, cancel := context.WithTimeout(context.Background(), 60*time.Millisecond)
		defer cancel()
		called := false
		err := WithReportWriter(ctx, alias, func(_ *ReportWriter) error { called = true; return nil })
		if called || !errors.Is(err, context.DeadlineExceeded) {
			return fmt.Errorf("alias bypassed held guard: callback=%t err=%v", called, err)
		}
		return nil
	})
	if err != nil {
		t.Fatal(err)
	}
}
