package trajectory

import (
	"context"
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

func TestSnapshotTimestampTransitionRequiresOneFullStableRead(t *testing.T) {
	for _, kind := range []string{"stable", "timestamp", "material", "replacement", "mode", "size", "second-transition", "recheck-failure"} {
		t.Run(kind, func(t *testing.T) {
			root, dir := snapshotFixture(t), snapshotStoreFixture(t)
			source, ref := captureFixture(t, root, dir)
			file := snapshotPath(dir, ref)
			f, err := os.Open(file)
			if err != nil {
				t.Fatal(err)
			}
			defer f.Close()
			initial, err := os.Lstat(file)
			if err != nil {
				t.Fatal(err)
			}
			opened, err := f.Stat()
			if err != nil {
				t.Fatal(err)
			}
			switch kind {
			case "material":
				raw := snapshotRead(t, file)
				raw[len(raw)-1] = 1
				if err = os.WriteFile(file, raw, 0600); err != nil {
					t.Fatal(err)
				}
			case "replacement":
				raw := snapshotRead(t, file)
				if err = os.Rename(file, file+".original"); err != nil {
					t.Fatal(err)
				}
				if err = os.WriteFile(file, raw, 0600); err != nil {
					t.Fatal(err)
				}
			case "mode":
				if err = os.Chmod(file, 0640); err != nil {
					t.Fatal(err)
				}
			case "size":
				w, err := os.OpenFile(file, os.O_APPEND|os.O_WRONLY, 0)
				if err != nil {
					t.Fatal(err)
				}
				_, err = w.Write([]byte{0})
				w.Close()
				if err != nil {
					t.Fatal(err)
				}
			}
			if kind != "stable" {
				if err = os.Chtimes(file, time.Now(), initial.ModTime().Add(2*time.Second)); err != nil {
					t.Fatal(err)
				}
			}
			calls := 0
			fault := errors.New("injected strict-reread failure")
			err = checkSnapshotFileStability(f, file, initial, opened, kind != "second-transition", func() error {
				calls++
				if kind == "recheck-failure" {
					return fault
				}
				return inspectSnapshotFileAttempt(context.Background(), file, ref, source, nil, nil, initial, false)
			})
			switch kind {
			case "stable":
				if err != nil || calls != 0 {
					t.Fatalf("stable archive reread or refused: %d %v", calls, err)
				}
			case "timestamp":
				if err != nil || calls != 1 {
					t.Fatalf("timestamp transition did not receive one full stable read: %d %v", calls, err)
				}
			case "material":
				if err == nil || calls != 1 {
					t.Fatalf("changed material bypassed strict full reread: %d %v", calls, err)
				}
			case "recheck-failure":
				if !errors.Is(err, fault) || calls != 1 {
					t.Fatalf("strict reread failure swallowed: %d %v", calls, err)
				}
			case "second-transition":
				if err == nil || calls != 0 {
					t.Fatalf("second timestamp transition received another allowance: %d %v", calls, err)
				}
			default:
				if err == nil || calls != 0 {
					t.Fatalf("identity/mode/size drift received reread permission: %d %v", calls, err)
				}
			}
		})
	}
}
func TestSnapshotStableRereadPinsFirstFileIdentity(t *testing.T) {
	root, dir := snapshotFixture(t), snapshotStoreFixture(t)
	source, ref := captureFixture(t, root, dir)
	original := snapshotPath(dir, ref)
	info, err := os.Lstat(original)
	if err != nil {
		t.Fatal(err)
	}
	other := filepath.Join(dir, "replacement.tar")
	if err = os.WriteFile(other, snapshotRead(t, original), 0600); err != nil {
		t.Fatal(err)
	}
	if err := inspectSnapshotFileAttempt(context.Background(), other, ref, source, nil, nil, info, false); err == nil || !strings.Contains(err.Error(), "origin changed") {
		t.Fatalf("stable reread accepted another inode with identical bytes: %v", err)
	}
}
