package runner

import (
	"context"
	"os"
	"parley-deck-cli/internal/agents"
	"parley-deck-cli/internal/protocol"
	"parley-deck-cli/internal/store"
	"path/filepath"
	"testing"
)

func TestPrecheckMirrorsActualArtifactSkipAndOverwrite(t *testing.T) {
	for _, tc := range []struct {
		name, phase string
		call        func(context.Context, Options) error
		paths       []string
		wantAgent   string
		overwrite   bool
	}{
		{"round-first-skipped", "round-02", PrecheckRound, []string{"round-02/first.md"}, "second", false},
		{"round-all-skipped", "round-02", PrecheckRound, []string{"round-02/first.md", "round-02/second.md"}, "", false},
		{"round-overwrite", "round-02", PrecheckRound, []string{"round-02/first.md"}, "first", true},
		{"review-first-skipped", "review", PrecheckReviewRound, []string{"review/round-02/first.md"}, "second", false},
		{"implementation-first-only-skipped", "implementation", PrecheckImplementation, []string{"IMPLEMENTATION.md"}, "", false},
		{"review-consensus-always-overwrites", "review-consensus", PrecheckReviewConsensus, []string{"review/consensus.md"}, "first", false},
	} {
		t.Run(tc.name, func(t *testing.T) {
			root := t.TempDir()
			writeLaunchProtocol(t, root)
			first, second := telemetryShell("exit 99", false), telemetryShell("exit 99", false)
			first.ID, second.ID, first.Found, second.Found = "first", "second", true, true
			idea := protocol.IdeaStatus{Slug: "demo", Path: filepath.Join(root, protocol.DeckDir, "ideas", "demo"), Participants: []string{"first", "second"}}
			for _, rel := range tc.paths {
				p := filepath.Join(idea.Path, rel)
				if err := os.MkdirAll(filepath.Dir(p), 0755); err != nil {
					t.Fatal(err)
				}
				if err := os.WriteFile(p, []byte("owned existing bytes"), 0600); err != nil {
					t.Fatal(err)
				}
			}
			opts := Options{Root: root, RunID: "skip-control", Idea: idea, Round: 2, Agents: []agents.Discovery{first, second}, Store: store.New(t.TempDir()), Overwrite: tc.overwrite}
			restore := refusePreparedProtocolCache(t, root, LaunchInfo{Idea: "demo", Phase: tc.phase})
			defer restore()
			err := tc.call(context.Background(), opts)
			records := terminalRecords(t, root)
			if tc.wantAgent == "" {
				if err != nil || len(records) != 0 {
					t.Fatalf("skipped launch was prechecked: err=%v terminals=%+v", err, records)
				}
				if tc.name == "round-all-skipped" {
					results := RunRound(context.Background(), opts)
					if len(results) != 2 {
						t.Fatalf("results: %+v", results)
					}
					for _, result := range results {
						if !result.Skipped {
							t.Fatalf("actual runner did not skip: %+v", result)
						}
					}
					if len(terminalRecords(t, root)) != 0 {
						t.Fatal("skipped actual round created an invocation")
					}
				}
			} else {
				if err == nil || len(records) != 1 || records[0].StartedAt != nil || records[0].Metadata.Agent != tc.wantAgent || records[0].Metadata.Phase != tc.phase {
					t.Fatalf("wrong pending launch prechecked: err=%v terminals=%+v", err, records)
				}
			}
			for _, rel := range tc.paths {
				raw, err := os.ReadFile(filepath.Join(idea.Path, rel))
				if err != nil || string(raw) != "owned existing bytes" {
					t.Fatalf("existing artifact changed: %s %v", rel, err)
				}
			}
		})
	}
}
