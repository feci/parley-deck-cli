//go:build windows

package fsutil

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// AC-DUR-3 mechanics (executed hosted even though the product sites are
// dormant behind the trajectoryRuntime POSIX gate): det-stage publish,
// O_EXCL/no-REPLACE anti-clobber, stage-left-on-error anti-replay, stale
// stage blocking, EEXIST-tolerant recheck for directories.
func TestPublishFileDurableMechanics(t *testing.T) {
	dir := t.TempDir()
	final := filepath.Join(dir, "request.json")
	if err := PublishFileDurable(final, []byte("payload"), 0o600); err != nil {
		t.Fatalf("publish: %v", err)
	}
	data, err := os.ReadFile(final)
	if err != nil || string(data) != "payload" {
		t.Fatalf("published content wrong: %q err=%v", data, err)
	}
	if _, err := os.Lstat(final + StageSuffix); !os.IsNotExist(err) {
		t.Fatalf("stage survived a successful publication: %v", err)
	}
	// No-REPLACE: a completed publication is never clobbered, and the failed
	// attempt leaves its stage as the anti-replay blocker.
	if err := PublishFileDurable(final, []byte("second"), 0o600); err == nil {
		t.Fatal("republish over a completed publication must fail")
	}
	if data, _ := os.ReadFile(final); string(data) != "payload" {
		t.Fatalf("completed publication was clobbered: %q", data)
	}
	if _, err := os.Lstat(final + StageSuffix); err != nil {
		t.Fatalf("failed republish must leave the stage blocker: %v", err)
	}
	// A stale stage blocks a fresh publication of a not-yet-existing final.
	other := filepath.Join(dir, "parent-result.json")
	if err := os.Mkdir(other+StageSuffix, 0o700); err != nil { // any stale stage shape
		t.Fatal(err)
	}
	if err := PublishFileDurable(other, []byte("x"), 0o600); err == nil {
		t.Fatal("stale stage must block publication")
	}
}

func TestPublishDirDurableMechanics(t *testing.T) {
	dir := t.TempDir()
	final := filepath.Join(dir, "trajectory-verification")
	if err := PublishDirDurable(final, 0o700); err != nil {
		t.Fatalf("publish dir: %v", err)
	}
	if info, err := os.Lstat(final); err != nil || !info.IsDir() {
		t.Fatalf("final dir missing: %v", err)
	}
	// EEXIST-tolerant recheck: publishing an existing directory is a no-op.
	if err := PublishDirDurable(final, 0o700); err != nil {
		t.Fatalf("idempotent republish must be tolerated: %v", err)
	}
	// A file where the directory belongs fails loudly.
	fileFinal := filepath.Join(dir, "occupied")
	if err := os.WriteFile(fileFinal, []byte("x"), 0o600); err != nil {
		t.Fatal(err)
	}
	if err := PublishDirDurable(fileFinal, 0o700); err == nil || !strings.Contains(err.Error(), "not a directory") {
		t.Fatalf("file-at-final must fail loudly: %v", err)
	}
	// A stale stage with no final is the visible blocker: loud, never cleared.
	staleFinal := filepath.Join(dir, "run-dir")
	if err := os.Mkdir(staleFinal+StageSuffix, 0o700); err != nil {
		t.Fatal(err)
	}
	err := PublishDirDurable(staleFinal, 0o700)
	if err == nil || !strings.Contains(err.Error(), "blocks the publication") {
		t.Fatalf("stale stage must be the loud blocker: %v", err)
	}
	// The blocker is deliberately left in place.
	if _, serr := os.Lstat(staleFinal + StageSuffix); serr != nil {
		t.Fatalf("stale stage must remain as evidence: %v", serr)
	}
	if _, ferr := os.Lstat(staleFinal); !os.IsNotExist(ferr) {
		t.Fatalf("nothing must be published past a blocker: %v", ferr)
	}
}
