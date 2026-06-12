package fsutil

import (
	"errors"
	"io/fs"
	"os"
	"path/filepath"
	"testing"
	"time"
)

// saveSeams snapshots the package seams and restores them after the test.
func saveSeams(t *testing.T) {
	t.Helper()
	m, s, sl := mkdirAll, stat, sleep
	t.Cleanup(func() { mkdirAll, stat, sleep = m, s, sl })
}

var errBoom = errors.New("boom")

func notExist(_ string) (os.FileInfo, error) { return nil, fs.ErrNotExist }

// Test_CommonPath: first os.MkdirAll succeeds → nil with no Stat and no sleep.
func Test_CommonPath(t *testing.T) {
	saveSeams(t)
	mkdirCalls, statCalls, sleepCalls := 0, 0, 0
	mkdirAll = func(string, os.FileMode) error { mkdirCalls++; return nil }
	stat = func(p string) (os.FileInfo, error) { statCalls++; return notExist(p) }
	sleep = func(time.Duration) { sleepCalls++ }

	if err := MkdirAllResilient("/whatever", 0o755); err != nil {
		t.Fatalf("want nil, got %v", err)
	}
	if mkdirCalls != 1 || statCalls != 0 || sleepCalls != 0 {
		t.Fatalf("healthy path must be one mkdir, zero stat, zero sleep; got mkdir=%d stat=%d sleep=%d", mkdirCalls, statCalls, sleepCalls)
	}
}

// Test_TransientThenSuccess: mkdir fails once then succeeds on the immediate retry.
func Test_TransientThenSuccess(t *testing.T) {
	saveSeams(t)
	mkdirCalls, sleepCalls := 0, 0
	mkdirAll = func(string, os.FileMode) error {
		mkdirCalls++
		if mkdirCalls == 1 {
			return errBoom
		}
		return nil
	}
	stat = notExist
	sleep = func(time.Duration) { sleepCalls++ }

	if err := MkdirAllResilient("/whatever", 0o755); err != nil {
		t.Fatalf("want nil, got %v", err)
	}
	if mkdirCalls != 2 || sleepCalls != 0 {
		t.Fatalf("want 2 mkdir attempts and 0 sleeps (immediate first retry); got mkdir=%d sleep=%d", mkdirCalls, sleepCalls)
	}
}

// Test_HostSucceededGuestLied: mkdir ALWAYS errors but a fresh Stat shows the dir exists
// (the observed virtio-fs case) → nil via isDir, without exhausting retries or sleeping.
func Test_HostSucceededGuestLied(t *testing.T) {
	saveSeams(t)
	dir := t.TempDir() // a real, existing directory
	mkdirCalls, sleepCalls := 0, 0
	mkdirAll = func(string, os.FileMode) error { mkdirCalls++; return errBoom }
	stat = os.Stat // real
	sleep = func(time.Duration) { sleepCalls++ }

	if err := MkdirAllResilient(dir, 0o755); err != nil {
		t.Fatalf("want nil (dir exists), got %v", err)
	}
	if mkdirCalls != 1 || sleepCalls != 0 {
		t.Fatalf("must short-circuit on isDir without retrying/sleeping; got mkdir=%d sleep=%d", mkdirCalls, sleepCalls)
	}
}

// Test_AlreadyExists: real os.MkdirAll on a pre-existing dir → nil, zero sleeps.
func Test_AlreadyExists(t *testing.T) {
	saveSeams(t)
	dir := t.TempDir()
	sleepCalls := 0
	sleep = func(time.Duration) { sleepCalls++ }
	// mkdirAll and stat stay real.

	if err := MkdirAllResilient(dir, 0o755); err != nil {
		t.Fatalf("want nil, got %v", err)
	}
	if sleepCalls != 0 {
		t.Fatalf("existing dir must not sleep; got sleep=%d", sleepCalls)
	}
}

// Test_GenuineFailure: persistent failure, path never appears → returns the LAST mkdir
// error (distinct from earlier ones) after the bounded retries, sleeping exactly
// 15ms,35ms,100ms,250ms,500ms,1000ms (the d>0 entries of retryDelays).
func Test_GenuineFailure(t *testing.T) {
	saveSeams(t)
	errLast := errors.New("final-boom")
	mkdirCalls := 0
	var sleeps []time.Duration
	mkdirAll = func(string, os.FileMode) error {
		mkdirCalls++
		if mkdirCalls == 8 { // distinct error on the final attempt
			return errLast
		}
		return errBoom
	}
	stat = notExist
	sleep = func(d time.Duration) { sleeps = append(sleeps, d) }

	err := MkdirAllResilient("/nope", 0o755)
	if !errors.Is(err, errLast) {
		t.Fatalf("want the LAST mkdir error %v, got %v", errLast, err)
	}
	if errors.Is(err, errBoom) {
		t.Fatalf("must return the last error, not an earlier one; got %v", err)
	}
	if mkdirCalls != 8 {
		t.Fatalf("want 8 mkdir attempts (initial + 7 retries), got %d", mkdirCalls)
	}
	want := []time.Duration{15 * time.Millisecond, 35 * time.Millisecond, 100 * time.Millisecond, 250 * time.Millisecond, 500 * time.Millisecond, 1000 * time.Millisecond}
	if len(sleeps) != len(want) {
		t.Fatalf("want sleeps %v, got %v", want, sleeps)
	}
	for i := range want {
		if sleeps[i] != want[i] {
			t.Fatalf("sleep[%d] = %v, want %v", i, sleeps[i], want[i])
		}
	}
}

// Test_FailFastPermission: a permission error is not retried/slept.
func Test_FailFastPermission(t *testing.T) {
	saveSeams(t)
	mkdirCalls, sleepCalls := 0, 0
	mkdirAll = func(string, os.FileMode) error { mkdirCalls++; return fs.ErrPermission }
	stat = notExist
	sleep = func(time.Duration) { sleepCalls++ }

	err := MkdirAllResilient("/denied", 0o755)
	if !errors.Is(err, fs.ErrPermission) {
		t.Fatalf("want permission error, got %v", err)
	}
	if mkdirCalls != 1 || sleepCalls != 0 {
		t.Fatalf("permission must fail fast: one mkdir, zero sleeps; got mkdir=%d sleep=%d", mkdirCalls, sleepCalls)
	}
}

// Test_DirExistsBeatsPermission: a permission error does NOT win when the directory
// already exists — the isDir check is evaluated before the fail-fast, so success wins.
func Test_DirExistsBeatsPermission(t *testing.T) {
	saveSeams(t)
	dir := t.TempDir() // real, existing directory
	sleepCalls := 0
	mkdirAll = func(string, os.FileMode) error { return fs.ErrPermission }
	stat = os.Stat // real → reports the existing dir
	sleep = func(time.Duration) { sleepCalls++ }

	if err := MkdirAllResilient(dir, 0o755); err != nil {
		t.Fatalf("isDir must win before permission fail-fast; want nil, got %v", err)
	}
	if sleepCalls != 0 {
		t.Fatalf("dir-exists short-circuit must not sleep; got %d", sleepCalls)
	}
}

// Test_NonDirCollision: mkdir reports fs.ErrExist (EEXIST) but a regular FILE sits at the
// path → fs.ErrExist must NOT be trusted blindly; the fresh Stat (IsDir false) keeps it an
// error rather than a masked success.
func Test_NonDirCollision(t *testing.T) {
	saveSeams(t)
	file := filepath.Join(t.TempDir(), "collision")
	if err := os.WriteFile(file, []byte("x"), 0o644); err != nil {
		t.Fatal(err)
	}
	mkdirAll = func(string, os.FileMode) error { return fs.ErrExist }
	stat = os.Stat // real → reports a regular file (IsDir false)
	sleep = func(time.Duration) {}

	if err := MkdirAllResilient(file, 0o755); err == nil {
		t.Fatal("fs.ErrExist must not be trusted blindly when a regular file collides; want error, got nil")
	}
}

// Test_AppendLine: lines land newline-terminated, appends accumulate, and the
// claim directory never lingers.
func Test_AppendLine(t *testing.T) {
	saveSeams(t)
	sleep = func(time.Duration) {}
	path := filepath.Join(t.TempDir(), "ledger", "index.jsonl")
	if err := AppendLine(path, []byte(`{"a":1}`)); err != nil {
		t.Fatal(err)
	}
	if err := AppendLine(path, []byte(`{"b":2}`+"\n")); err != nil {
		t.Fatal(err)
	}
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	if string(data) != `{"a":1}`+"\n"+`{"b":2}`+"\n" {
		t.Fatalf("ledger content: %q", data)
	}
	if _, err := os.Stat(path + ".lock"); !os.IsNotExist(err) {
		t.Fatalf("claim dir must be removed, stat err=%v", err)
	}
}

// Test_AppendLine_StuckClaim: a wedged claim degrades to an unlocked append
// after the bounded wait instead of losing the record.
func Test_AppendLine_StuckClaim(t *testing.T) {
	saveSeams(t)
	sleep = func(time.Duration) {}
	path := filepath.Join(t.TempDir(), "index.jsonl")
	if err := os.Mkdir(path+".lock", 0o755); err != nil {
		t.Fatal(err)
	}
	if err := AppendLine(path, []byte("x")); err != nil {
		t.Fatal(err)
	}
	if _, err := os.Stat(path); err != nil {
		t.Fatalf("record lost behind a stuck claim: %v", err)
	}
}
