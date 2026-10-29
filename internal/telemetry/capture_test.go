package telemetry

import (
	"bytes"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
	"time"
)

func TestZcodeCaptureScrubAndTail(t *testing.T) {
	raw := recordedZcode(t).Stderr
	raw = strings.Replace(raw, "messages: []", "messages: ['PRIVATE_PROMPT']", 1)
	raw = strings.Replace(raw, "'content-type': 'application/json'", "'authorization': 'Bearer PRIVATE_CREDENTIAL', 'cookie': 'PRIVATE_COOKIE'", 1)
	out, err := ScrubZcodeCapture([]byte(raw), false)
	if err != nil {
		t.Fatal(err)
	}
	for _, secret := range []string{"PRIVATE_", "provider.invalid", "offlineFixture", "fixture-model"} {
		if bytes.Contains(out, []byte(secret)) {
			t.Fatalf("unscrubbed %s", secret)
		}
	}
	if !bytes.Contains(out, []byte("[REDACTED]")) || !bytes.Contains(out, []byte("Weekly/Monthly Limit Exhausted")) || !bytes.Contains(out, []byte("Error: Turn execution failed")) {
		t.Fatal(string(out))
	}
	for _, bad := range []string{raw + "\nprompt: secret", strings.ReplaceAll(raw, "Weekly/Monthly Limit Exhausted", "PRIVATE_PROMPT"), strings.ReplaceAll(raw, `"retry_after":176930`, `"unknown":"PRIVATE_CREDENTIAL"`)} {
		if b, err := ScrubZcodeCapture([]byte(bad), false); err == nil || len(b) != 0 {
			t.Fatal("uncertain scrub persisted")
		}
	}
	tail := strings.Repeat("private", CaptureFileLimit) + "\nError: Turn execution failed (traceId: safe)\n"
	out, err = ScrubZcodeCapture([]byte(tail[len(tail)-CaptureFileLimit:]), true)
	if err != nil || !bytes.HasPrefix(out, []byte("[TRUNCATED]\n[REDACTED]")) || bytes.Contains(out, []byte("private")) {
		t.Fatalf("tail: %s %v", out, err)
	}
}

func captureTestDir(t *testing.T) string {
	t.Helper()
	root, err := filepath.EvalSymlinks(t.TempDir())
	if err != nil {
		t.Fatal(err)
	}
	return filepath.Join(root, "captures", "zcode")
}

func TestZcodeCapturePrivateBoundsAndCleanup(t *testing.T) {
	dir := captureTestDir(t)
	now := time.Now()
	raw := []byte(recordedZcode(t).Stderr)
	if runtime.GOOS == "windows" {
		made, err := CaptureZcodeFailure(dir, "inv-1", raw, false, now)
		if made || err == nil || !strings.Contains(err.Error(), "mode 0700") {
			t.Fatalf("private-mode refusal: %v %v", made, err)
		}
		if _, err := os.Stat(filepath.Join(dir, "inv-1.log")); !os.IsNotExist(err) {
			t.Fatalf("capture exists after private-mode refusal: %v", err)
		}
		return
	}
	if made, err := CaptureZcodeFailure(dir, "inv-1", raw, false, now); err != nil || !made {
		t.Fatalf("%v %v", made, err)
	}
	if made, err := CaptureZcodeFailure(dir, "inv-1", raw, false, now); err != nil || made {
		t.Fatalf("replay: %v %v", made, err)
	}
	for _, path := range []string{filepath.Dir(dir), dir} {
		info, _ := os.Stat(path)
		if info.Mode().Perm() != 0700 {
			t.Fatal(info.Mode())
		}
	}
	info, _ := os.Stat(filepath.Join(dir, "inv-1.log"))
	if info.Mode().Perm() != 0600 {
		t.Fatal(info.Mode())
	}
	os.Chtimes(filepath.Join(dir, "inv-1.log"), now.Add(-8*24*time.Hour), now.Add(-8*24*time.Hour))
	for i := 0; i < 9; i++ {
		path := filepath.Join(dir, string(rune('a'+i))+".log")
		os.WriteFile(path, bytes.Repeat([]byte("x"), CaptureFileLimit), 0600)
		os.Chtimes(path, now.Add(time.Duration(i-10)*time.Hour), now.Add(time.Duration(i-10)*time.Hour))
	}
	state, err := CaptureStatus(dir, now, false)
	if err != nil || state.ExpiredRemoved != 1 || state.CapacityRemoved != 2 || state.Bytes >= CaptureTotalLimit {
		t.Fatalf("%+v %v", state, err)
	}
	receipt, _ := os.ReadFile(filepath.Join(dir, "cleanup.json"))
	if len(receipt) > 2048 || bytes.Contains(receipt, []byte("Limit Exhausted")) {
		t.Fatal(string(receipt))
	}
	state, err = CaptureStatus(dir, now, true)
	if err != nil || state.Purged != 7 || state.Files != 0 {
		t.Fatalf("purge: %+v %v", state, err)
	}
}

func TestZcodeCaptureRejectsUnsafeStorage(t *testing.T) {
	for _, name := range []string{"capture.lease", "cleanup.json", "inv.log", "directory-mode", "file-mode", "directory-link"} {
		t.Run(name, func(t *testing.T) {
			dir := captureTestDir(t)
			if err := os.MkdirAll(dir, 0700); err != nil {
				t.Fatal(err)
			}
			switch name {
			case "directory-mode":
				if err := os.Chmod(dir, 0755); err != nil {
					t.Fatal(err)
				}
			case "file-mode":
				path := filepath.Join(dir, "inv.log")
				if err := os.WriteFile(path, []byte("x"), 0644); err != nil {
					t.Fatal(err)
				}
				if err := os.Chmod(path, 0644); err != nil {
					t.Fatal(err)
				}
			case "directory-link":
				if err := os.Remove(dir); err != nil {
					t.Fatal(err)
				}
				if err := os.Symlink(t.TempDir(), dir); err != nil {
					t.Fatal(err)
				}
			default:
				target := filepath.Join(t.TempDir(), "target")
				if err := os.WriteFile(target, []byte("preserve"), 0600); err != nil {
					t.Fatal(err)
				}
				if err := os.Symlink(target, filepath.Join(dir, name)); err != nil {
					t.Fatal(err)
				}
			}
			if _, err := CaptureStatus(dir, time.Now(), true); err == nil {
				t.Fatal("unsafe status/purge accepted")
			}
			if _, err := CaptureZcodeFailure(dir, "new", []byte(recordedZcode(t).Stderr), false, time.Now()); err == nil {
				t.Fatal("unsafe capture accepted")
			}
		})
	}
}
