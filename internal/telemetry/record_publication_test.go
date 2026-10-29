package telemetry

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"runtime"
	"testing"
)

func TestInvocationPublicationFaults(t *testing.T) {
	for _, stage := range []string{"write", "file-sync", "close", "link", "directory-sync"} {
		t.Run(stage, func(t *testing.T) {
			store := t.TempDir()
			injected := errors.New("injected " + stage)
			prior := invocationWriteFault
			defer func() { invocationWriteFault = prior }()
			seen := false
			invocationWriteFault = func(at, path string) error {
				if at != stage {
					return nil
				}
				seen = true
				if at == "directory-sync" {
					raw, err := os.ReadFile(path)
					if err != nil || !json.Valid(raw) {
						t.Fatalf("published record is not complete: %q %v", raw, err)
					}
				} else if _, err := os.Lstat(path); !os.IsNotExist(err) {
					t.Fatalf("record visible before publication at %s: %v", at, err)
				}
				staged, err := filepath.Glob(filepath.Join(filepath.Dir(path), ".record-*"))
				if err != nil || len(staged) != 1 {
					t.Fatalf("staging files: %v %v", staged, err)
				}
				info, err := os.Stat(staged[0])
				if err != nil {
					t.Fatal(err)
				}
				if runtime.GOOS != "windows" && info.Mode().Perm() != 0o600 {
					t.Fatalf("staging mode = %o", info.Mode().Perm())
				}
				return injected
			}
			inv, err := BeginBound(store, "fault-attempt", Metadata{Agent: "alpha"})
			if inv != nil || !errors.Is(err, injected) || !seen {
				t.Fatalf("failed request admitted: inv=%v err=%v checkpoint=%v", inv, err, seen)
			}
			dir := filepath.Join(store, "fault-attempt")
			entries, err := os.ReadDir(dir)
			if err != nil {
				t.Fatal(err)
			}
			if stage == "directory-sync" {
				if len(entries) != 1 || entries[0].Name() != "requested.json" {
					t.Fatalf("published state or staging cleanup lost: %v", entries)
				}
				if r := readRecord(t, filepath.Join(dir, "requested.json")); r.Type != "invocation.requested" || r.StartedAt != nil {
					t.Fatalf("failed publication invented a child: %+v", r)
				}
			} else if len(entries) != 0 {
				t.Fatalf("failed staging left files: %v", entries)
			}
			invocationWriteFault = prior
			if _, err = BeginBound(store, "fault-attempt", Metadata{}); err == nil {
				t.Fatal("failed invocation directory was reused")
			}
		})
	}
}

func TestInvocationPublicationRefusesExistingNames(t *testing.T) {
	for _, kind := range []string{"file", "symlink", "directory"} {
		t.Run(kind, func(t *testing.T) {
			dir := t.TempDir()
			path := filepath.Join(dir, "terminal.json")
			preserved := path
			switch kind {
			case "symlink":
				preserved = filepath.Join(t.TempDir(), "target")
				if err := os.Symlink(preserved, path); err != nil {
					t.Fatal(err)
				}
			case "directory":
				if err := os.Mkdir(path, 0o700); err != nil {
					t.Fatal(err)
				}
				preserved = filepath.Join(path, "original")
			}
			if err := os.WriteFile(preserved, []byte("preserve exactly"), 0o600); err != nil {
				t.Fatal(err)
			}
			inv := &Invocation{Dir: dir, ID: "collision", record: Record{Type: "invocation.terminal"}}
			if err := inv.write("terminal.json"); !os.IsExist(err) {
				t.Fatalf("existing %s did not refuse: %v", kind, err)
			}
			raw, err := os.ReadFile(preserved)
			if err != nil || string(raw) != "preserve exactly" {
				t.Fatalf("existing %s changed: %q %v", kind, raw, err)
			}
			entries, err := os.ReadDir(dir)
			if err != nil || len(entries) != 1 || entries[0].Name() != "terminal.json" {
				t.Fatalf("collision left staging files: %v %v", entries, err)
			}
			if kind == "symlink" {
				if target, err := os.Readlink(path); err != nil || target != preserved {
					t.Fatalf("symlink changed: %q %v", target, err)
				}
			}
		})
	}
}

func TestInvocationTerminalPublicationFaults(t *testing.T) {
	for _, stage := range []string{"write", "file-sync", "close", "link", "directory-sync"} {
		t.Run(stage, func(t *testing.T) {
			inv, err := Begin(t.TempDir(), Metadata{Agent: "alpha"})
			if err != nil {
				t.Fatal(err)
			}
			injected := errors.New("injected terminal " + stage)
			prior := invocationWriteFault
			defer func() { invocationWriteFault = prior }()
			invocationWriteFault = func(at, path string) error {
				if at == stage && filepath.Base(path) == "terminal.json" {
					return injected
				}
				return nil
			}
			for attempt := 0; attempt < 2; attempt++ {
				if err := inv.Finish(Outcome{Status: "finished"}); !errors.Is(err, injected) {
					t.Fatalf("terminal failure %d lost: %v", attempt, err)
				}
			}
			path := filepath.Join(inv.Dir, "terminal.json")
			if stage == "directory-sync" {
				r := readRecord(t, path)
				if r.Type != "invocation.terminal" || r.Outcome.Status != "finished" {
					t.Fatalf("published terminal lost: %+v", r)
				}
			} else if _, err := os.Lstat(path); !os.IsNotExist(err) {
				t.Fatalf("failed staging published a terminal: %v", err)
			}
			staged, err := filepath.Glob(filepath.Join(inv.Dir, ".record-*"))
			if err != nil || len(staged) != 0 {
				t.Fatalf("terminal left staging files: %v %v", staged, err)
			}
		})
	}
}

func TestInvocationStoreSyncFailurePreventsReservation(t *testing.T) {
	store := filepath.Join(t.TempDir(), "new-runtime", "invocations")
	injected := errors.New("injected store sync")
	prior := invocationWriteFault
	defer func() { invocationWriteFault = prior }()
	invocationWriteFault = func(at, path string) error {
		if at == "store-sync" && path == store {
			return injected
		}
		return nil
	}
	inv, err := BeginBound(store, "not-started", Metadata{})
	if inv != nil || !errors.Is(err, injected) {
		t.Fatalf("store failure admitted a request: %v %v", inv, err)
	}
	entries, err := os.ReadDir(store)
	if err != nil || len(entries) != 0 {
		t.Fatalf("store failure reserved an invocation: %v %v", entries, err)
	}
}

func TestInvocationStagingResidueDoesNotEnterAdmission(t *testing.T) {
	root := t.TempDir()
	store := filepath.Join(root, ".parley-runtime", "invocations")
	inv, err := Begin(store, Metadata{Idea: "sample", Agent: "alpha"})
	if err != nil {
		t.Fatal(err)
	}
	for _, dir := range []string{inv.Dir, filepath.Join(store, "interrupted")} {
		if err := os.MkdirAll(dir, 0o700); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(filepath.Join(dir, ".record-crash"), []byte(`{"partial":`), 0o600); err != nil {
			t.Fatal(err)
		}
	}
	if err := ExhaustionAdmission(root, "sample", "alpha"); err != nil {
		t.Fatalf("inert crash residue entered admission: %v", err)
	}
	if err := os.WriteFile(filepath.Join(inv.Dir, "requested.json"), []byte(`{"partial":`), 0o600); err != nil {
		t.Fatal(err)
	}
	if err := ExhaustionAdmission(root, "sample", "alpha"); err == nil {
		t.Fatal("corrupt published request was ignored")
	}
}
