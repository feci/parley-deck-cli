package budget

import (
	"context"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestStepPolicyRejectsMalformedOrMissingLimits(t *testing.T) {
	root := t.TempDir()
	ctx := context.Background()
	b, err := EnsureStepBinding(ctx, root, "idea", 1, 0)
	if err != nil {
		t.Fatal(err)
	}
	path := filepath.Join(filepath.Dir(b.Store.Dir), "policy.json")
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	for _, bad := range []string{strings.Replace(string(data), "\"max_steps\": 1,", "", 1), strings.Replace(string(data), "\"max_steps\": 1", "\"max_steps\": null", 1), strings.Replace(string(data), "\"max_steps\": 1", "\"max_steps\": 1, \"max_steps\": 0", 1), strings.Replace(string(data), "\"max_steps\": 1", "\"Max_Steps\": 0", 1)} {
		if err = os.WriteFile(path, []byte(bad), 0600); err != nil {
			t.Fatal(err)
		}
		if _, err = EnsureStepBinding(ctx, root, "idea", 0, 0); err == nil {
			t.Fatal("malformed frozen ceiling allowed reset")
		}
	}
	if err = os.Remove(path); err != nil {
		t.Fatal(err)
	}
	if _, err = EnsureStepBinding(ctx, root, "idea", 1, 0); err == nil {
		t.Fatal("missing policy silently recreated")
	}
}
func TestStepPolicyRefusesHistoricalDriverCharges(t *testing.T) {
	for _, event := range []string{"run.phase", "loop.budget", "driver.action"} {
		t.Run(event, func(t *testing.T) {
			root := t.TempDir()
			path := filepath.Join(root, "parley-deck", "runs", "old", "events.jsonl")
			if err := os.MkdirAll(filepath.Dir(path), 0700); err != nil {
				t.Fatal(err)
			}
			if err := os.WriteFile(path, []byte("{\"type\":\""+event+"\",\"data\":{\"idea\":\"idea\"}}\n"), 0600); err != nil {
				t.Fatal(err)
			}
			if _, err := EnsureStepBinding(context.Background(), root, "idea", 1, 0); err == nil {
				t.Fatal("old driver counts were zeroed")
			}
			if b, err := LoadStepBinding(context.Background(), root, "idea"); err != nil || b != nil {
				t.Fatalf("legacy refusal activated partial policy: %+v %v", b, err)
			}
		})
	}
}
