package budget

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"testing"
)

func TestStepHistoryIncludesFailedInvocationBeforeCursorPublication(t *testing.T) {
	for _, phase := range []string{"preflight", "round-01", "round-02", "fixup", "consensus", "unspecified"} {
		t.Run(phase, func(t *testing.T) {
			root := t.TempDir()
			dir := filepath.Join(root, ".parley-runtime", "invocations", "failed-before-cursor")
			if err := os.MkdirAll(dir, 0700); err != nil {
				t.Fatal(err)
			}
			body := fmt.Sprintf(`{"metadata":{"idea":"idea","phase":%q}}`, phase)
			if err := os.WriteFile(filepath.Join(dir, "requested.json"), []byte(body), 0600); err != nil {
				t.Fatal(err)
			}
			_, err := EnsureStepBinding(context.Background(), root, "idea", 1, 0)
			allowed := phase == "preflight" || phase == "round-01"
			if (err == nil) != allowed {
				t.Fatalf("allowed=%v, err=%v", allowed, err)
			}
		})
	}
}

func TestStepHistoryResolvesCursorFromConsistentRunEvents(t *testing.T) {
	for _, tc := range []struct {
		name, events string
		allow        bool
	}{
		{"unrelated", `{"type":"run.created","data":{"idea":"other"}}` + "\n" + `{"type":"run.phase","data":{"phase":"review"}}`, true},
		{"same-idea", `{"type":"run.created","data":{"idea":"idea"}}`, false},
		{"unknown-cursor", "", false},
		{"reassigned-run", `{"type":"run.created","data":{"idea":"idea"}}` + "\n" + `{"type":"run.created","data":{"idea":"other"}}`, false},
		{"conflicting-phase", `{"type":"run.created","data":{"idea":"other"}}` + "\n" + `{"type":"run.phase","data":{"idea":"idea"}}`, false},
		{"null-identity", `{"type":"run.created","data":{"idea":null}}`, false},
		{"aliased-identity", `{"type":"run.created","data":{"Idea":"other"}}`, false},
		{"duplicate-identity", `{"type":"run.created","data":{"idea":"idea","idea":"other"}}`, false},
		{"partial-json", `{"type":"run.created","data":{"idea":"other"}`, false},
	} {
		t.Run(tc.name, func(t *testing.T) {
			root := t.TempDir()
			dir := filepath.Join(root, "parley-deck", "runs", "old")
			if err := os.MkdirAll(dir, 0700); err != nil {
				t.Fatal(err)
			}
			for name, body := range map[string]string{"driver.json": `{"phase":"review","current_round":2,"fixup_cycles_published":1}`, "events.jsonl": tc.events} {
				if err := os.WriteFile(filepath.Join(dir, name), []byte(body), 0600); err != nil {
					t.Fatal(err)
				}
			}
			_, err := EnsureStepBinding(context.Background(), root, "idea", 1, 0)
			if (err == nil) != tc.allow {
				t.Fatalf("allow=%v: %v", tc.allow, err)
			}
		})
	}
}

func TestStepHistoryAllowsRoundOneBeforeFirstDriver(t *testing.T) {
	root := t.TempDir()
	dir := filepath.Join(root, "parley-deck", "runs", "first")
	if err := os.MkdirAll(dir, 0700); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(dir, "events.jsonl"), []byte("{\"type\":\"run.created\",\"data\":{\"idea\":\"idea\"}}\n{\"type\":\"round.completed\",\"data\":{\"round\":\"round-01\"}}\n"), 0600); err != nil {
		t.Fatal(err)
	}
	if _, err := EnsureStepBinding(context.Background(), root, "idea", 1, 0); err != nil {
		t.Fatal(err)
	}
}

func TestStepHistoryBoundedReaderRejectsAliasesAndOversize(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "events.jsonl")
	if err := os.WriteFile(path, []byte("12345"), 0600); err != nil {
		t.Fatal(err)
	}
	if _, err := readStepHistoryFile(path, 4); err == nil {
		t.Fatal("oversized history accepted")
	}
	if data, err := readStepHistoryFile(path, 5); err != nil || string(data) != "12345" {
		t.Fatalf("bounded regular file: %s %v", data, err)
	}
	alias := filepath.Join(dir, "alias")
	if err := os.Symlink(path, alias); err != nil {
		t.Skipf("symlinks unavailable: %v", err)
	}
	if _, err := readStepHistoryFile(alias, 5); err == nil {
		t.Fatal("aliased history accepted")
	}
}
