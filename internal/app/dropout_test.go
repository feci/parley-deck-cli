package app

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"parley-deck-cli/internal/agents"
	"parley-deck-cli/internal/consensus"
	"parley-deck-cli/internal/quota"
	"parley-deck-cli/internal/telemetry"
)

func TestDropoutSignoffInvalidOwnSuffixAndTampering(t *testing.T) {
	for _, kind := range []string{"invalid-own", "tampering", "block"} {
		t.Run(kind, func(t *testing.T) {
			root, idea, _ := fixupAppFixture(t, quota.NewParticipantPolicy(nil, nil))
			promptPath := filepath.Join(idea.Path, "00-prompt.md")
			prompt, _ := os.ReadFile(promptPath)
			writeConsensusIdea(t, root, idea.Slug, idea.Participants, false, nil)
			os.WriteFile(promptPath, prompt, 0600)
			path := filepath.Join(idea.Path, "consensus.md")
			raw, _ := os.ReadFile(path)
			raw = bytes.Replace(raw, []byte("---\n"), []byte("---\ndrafted-by: alpha\n"), 1)
			os.WriteFile(path, raw, 0600)
			bin := t.TempDir()
			alpha := writeFakeSignoffCLI(t, bin, "alpha", "accept", 0)
			beta := writeFakeSignoffCLI(t, bin, "beta", "accept", 0)
			gamma := filepath.Join(bin, "gamma")
			target := "'" + strings.ReplaceAll(path, "'", "'\\''") + "'"
			script := "#!/bin/sh\ncat >/dev/null\nprintf '\\n### Signoff: gamma — 2026-10-08\\n\\nStatus: nonsense\\n' >> " + target + "\nexit 7\n"
			if kind == "tampering" {
				script = "#!/bin/sh\ncat >/dev/null\nprintf 'overwrote all signoffs' > " + target + "\nexit 7\n"
			}
			if kind == "block" {
				gamma = writeFakeSignoffCLI(t, bin, "gamma", "block", 1)
			} else {
				os.WriteFile(gamma, []byte(script), 0700)
			}
			writeAgentsLocalConfig(t, root, fakeAgentConfig{ID: "alpha", Path: alpha}, fakeAgentConfig{ID: "beta", Path: beta}, fakeAgentConfig{ID: "gamma", Path: gamma})
			var out, errs bytes.Buffer
			err := requestConsensusSignoffs(context.Background(), requestSignoffsOptions{Root: root, IdeaSlug: idea.Slug, Yes: true}, &out, &errs)
			h, he := quota.ReadHistory(idea.Path)
			if he != nil {
				t.Fatal(he)
			}
			if kind == "invalid-own" {
				if err != nil || !h.Dropped("gamma") {
					t.Fatalf("%v history=%+v\n%s\n%s", err, h, out.String(), errs.String())
				}
				summary, e := consensus.Status(root, idea.Slug, false)
				if e != nil || summary.Triage != consensus.TriageReady {
					t.Fatalf("%+v %v", summary, e)
				}
				copies, _ := filepath.Glob(filepath.Join(root, ".parley-runtime", "invocations", "*", "participant-files", "00-consensus.md"))
				invalid := 0
				for _, p := range copies {
					b, _ := os.ReadFile(p)
					if strings.Contains(string(b), "Status: nonsense") {
						invalid++
					}
				}
				if invalid != 2 {
					t.Fatal("invalid suffix copies", invalid)
				}
			} else {
				if err == nil || h.Revision != 0 {
					t.Fatalf("%s did not stop: %+v %v", kind, h, err)
				}
				if kind == "block" {
					summary, e := consensus.Status(root, idea.Slug, false)
					if e != nil || summary.Triage != consensus.TriageBlocked {
						t.Fatalf("BLOCK lost: %+v %v", summary, e)
					}
				}
			}
			paths, _ := filepath.Glob(filepath.Join(root, ".parley-runtime", "invocations", "*", "terminal.json"))
			count := 0
			for _, p := range paths {
				b, _ := os.ReadFile(p)
				var r telemetry.Record
				json.Unmarshal(b, &r)
				if r.Metadata.Agent == "gamma" {
					count++
				}
			}
			want := 1
			if kind == "invalid-own" {
				want = 2
			}
			if count != want {
				t.Fatalf("attempts=%d want=%d", count, want)
			}
		})
	}
}

func TestDropoutReadinessRetriesBeforeMembership(t *testing.T) {
	root, _, _ := fixupAppFixture(t, quota.NewParticipantPolicy(nil, nil))
	p := quota.NewParticipantPolicy(nil, nil)
	var discovered []agents.Discovery
	for _, id := range []string{"alpha", "beta", "gamma"} {
		script := "cat >/dev/null; printf PONG"
		if id == "gamma" {
			script = "cat >/dev/null; printf 'HTTP 503' >&2; exit 9"
		}
		discovered = append(discovered, agents.Discovery{Spec: agents.Spec{ID: id, HeadlessArgs: []string{"-c", script}, PromptMode: agents.PromptStdin}, Found: true, Path: "/bin/sh"})
	}
	start := time.Now()
	entries := checkRoster(context.Background(), preflightOptions{Root: root, QuotaPolicy: &p, PingTimeout: time.Second}, discovered)
	if time.Since(start) < 5*time.Second {
		t.Fatal("retry delay was skipped")
	}
	ms := []quota.Member{}
	for _, e := range entries {
		ms = append(ms, quota.Member{ID: e.RosterID, Usable: e.Available, Evidence: e.QuotaEvidence})
	}
	d := quota.Evaluate(p, []string{"alpha", "beta", "gamma"}, ms, quota.Roles{})
	if !d.Applied || len(d.Candidates) != 1 || len(d.Candidates[0].Evidence.Failure.Attempts) != 2 {
		t.Fatalf("%+v entries=%+v", d, entries)
	}
	// Standalone reporting has no policy and therefore no automatic retry/reduction.
	standalone := checkRoster(context.Background(), preflightOptions{Root: root, PingTimeout: time.Second}, discovered)
	if standalone[2].QuotaEvidence != nil && standalone[2].QuotaEvidence.Failure != nil {
		t.Fatal("standalone authorized dropout")
	}
}

func TestDropoutKickoffBlockCreatesSafeInbox(t *testing.T) {
	root := t.TempDir()
	d := quota.Decision{Before: []string{"a", "b", "c"}, UsableSurvivors: 1, Block: "usable floor"}
	if err := writeQuotaBlock(root, d); err != nil {
		t.Fatal(err)
	}
	files, _ := filepath.Glob(filepath.Join(root, "parley-deck", "inbox", "*.md"))
	if len(files) != 1 {
		t.Fatal(files)
	}
	raw, _ := os.ReadFile(files[0])
	if !strings.Contains(string(raw), quota.OwnerOptions) || !strings.Contains(string(raw), "fixed floor: 2") {
		t.Fatal(string(raw))
	}
	if err := writeQuotaBlock(root, d); err != nil {
		t.Fatal(err)
	}
	os.Remove(files[0])
	os.Remove(filepath.Dir(files[0]))
	os.Symlink(t.TempDir(), filepath.Dir(files[0]))
	if err := writeQuotaBlock(root, d); err == nil {
		t.Fatal(fmt.Sprint("unsafe inbox accepted"))
	}
}
