package app

import (
	"bytes"
	"context"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"parley-deck-cli/internal/agents"
	"parley-deck-cli/internal/consensus"
	"parley-deck-cli/internal/membership"
	"parley-deck-cli/internal/protocol"
	"parley-deck-cli/internal/quota"
	"parley-deck-cli/internal/runmanifest"
	"parley-deck-cli/internal/runner"
)

// Exercise the actual standalone signoff entry: it creates a new run, preserves
// its prior model snapshot, supervises a stalled signer, settles the entire batch
// and lets an already-constructed implementation adapter observe the new quorum.
func TestReviewGateSignoffWatchdogAndSameTickRebind(t *testing.T) {
	root, idea, run := fixupAppFixture(t, quota.NewParticipantPolicy(nil, nil))
	pp := filepath.Join(idea.Path, "00-prompt.md")
	pr, _ := os.ReadFile(pp)
	pr = bytes.Replace(pr, []byte("status:"), []byte("auto_implement: true\nimplementer: alpha\nrequire_model_diversity: false\nstatus:"), 1)
	writeConsensusIdea(t, root, idea.Slug, idea.Participants, false, nil)
	os.WriteFile(pp, pr, 0600)
	cp := filepath.Join(idea.Path, "consensus.md")
	cr, _ := os.ReadFile(cp)
	os.WriteFile(cp, bytes.Replace(cr, []byte("---\n"), []byte("---\ndrafted-by: alpha\n"), 1), 0600)
	m, _ := runmanifest.Load(root, run)
	m.RosterSnapshot = []runmanifest.RosterSnapshotEntry{{Agent: "alpha", Model: "model-a"}, {Agent: "beta", Model: "model-b"}, {Agent: "gamma", Model: "model-c"}}
	if err := runmanifest.Write(root, run, m); err != nil {
		t.Fatal(err)
	}
	bin := t.TempDir()
	a := writeFakeSignoffCLI(t, bin, "alpha", "accept", 0)
	b := writeFakeSignoffCLI(t, bin, "beta", "accept", 0)
	c := filepath.Join(bin, "gamma")
	os.WriteFile(c, []byte("#!/bin/sh\ncat >/dev/null\nprintf ready\nexec sleep 30\n"), 0700)
	writeAgentsLocalConfig(t, root, fakeAgentConfig{ID: "alpha", Path: a}, fakeAgentConfig{ID: "beta", Path: b}, fakeAgentConfig{ID: "gamma", Path: c})
	cfg := filepath.Join(root, protocol.DeckDir, "agents.local.toml")
	raw, _ := os.ReadFile(cfg)
	os.WriteFile(cfg, []byte(strings.Replace(string(raw), "[agents.gamma]\n", "[agents.gamma]\nfirst_event_timeout_ms = 50\nstall_timeout_ms = 50\nheartbeat_ms = 50\n", 1)), 0600)
	base := runner.Options{Root: root, RunID: run, Idea: idea, Agents: []agents.Discovery{{Spec: agents.Spec{ID: "alpha", Model: "model-a"}, Found: true}, {Spec: agents.Spec{ID: "beta", Model: "model-b"}, Found: true}, {Spec: agents.Spec{ID: "gamma", Model: "model-c"}, Found: true}}}
	var out, errs bytes.Buffer
	o := newDriverImplOps(base, root, idea.Slug, idea.Path, idea.Participants, &out).(driverImplOps)
	if err := requestConsensusSignoffs(context.Background(), requestSignoffsOptions{Root: root, IdeaSlug: idea.Slug, Yes: true}, &out, &errs); err != nil {
		t.Fatalf("%v\n%s\n%s", err, &out, &errs)
	}
	h, err := quota.ReadHistory(idea.Path)
	if err != nil || !h.Dropped("gamma") || h.Revision != 1 {
		t.Fatalf("history %+v %v", h, err)
	}
	summary, err := consensus.Status(root, idea.Slug, false)
	if err != nil || summary.Triage != consensus.TriageReady {
		t.Fatalf("remaining signers: %+v %v", summary, err)
	}
	for _, r := range dropoutAppRecords(t, root) {
		if r.Metadata.Agent == "gamma" && (r.Outcome.FailureClass == nil || *r.Outcome.FailureClass != "stalled") {
			t.Fatalf("generic instead of watchdog terminal: %+v", r.Outcome)
		}
	}
	if len(h.Batches[0].Decision.Candidates[0].Evidence.Failure.Attempts) != 2 {
		t.Fatal("watchdog did not use exactly two slots")
	}
	o, err = o.quotaCurrent()
	if err != nil || len(o.reviewers) != 1 || o.reviewers[0] != "beta" {
		t.Fatalf("stale adapter: %+v %v", o.reviewers, err)
	}
	if ok, err := o.singleReviewerAfterDropout(); err != nil || !ok {
		t.Fatalf("derived gate: %v %v", ok, err)
	}
	newRun := h.Batches[0].RunID
	if newRun == run {
		t.Fatal("fixture did not exercise standalone run")
	}
	if ok, err := membership.SingleReviewerAfterDropout(root, idea.Path, newRun, "alpha", h.Current, nil); err != nil || !ok {
		t.Fatalf("new signoff run lost snapshot: %v %v", ok, err)
	}
}

// A missing snapshot is compatible with old runs. A missing manifest is not:
// immutable history has always required the kickoff and every touched run.
func TestReviewGateSignoffSnapshotAndManifestIntegrity(t *testing.T) {
	for _, kind := range []string{"absent-snapshot", "missing-manifest", "corrupt-manifest", "foreign-manifest"} {
		t.Run(kind, func(t *testing.T) {
			root, idea, run := fixupAppFixture(t, quota.NewParticipantPolicy(nil, nil))
			prompt := filepath.Join(idea.Path, "00-prompt.md")
			original, err := os.ReadFile(prompt)
			if err != nil {
				t.Fatal(err)
			}
			writeConsensusIdea(t, root, idea.Slug, idea.Participants, false, nil)
			if err := os.WriteFile(prompt, original, 0600); err != nil {
				t.Fatal(err)
			}
			consensusPath := filepath.Join(idea.Path, "consensus.md")
			consensusRaw, err := os.ReadFile(consensusPath)
			if err != nil {
				t.Fatal(err)
			}
			if err := os.WriteFile(consensusPath, bytes.Replace(consensusRaw, []byte("---\n"), []byte("---\ndrafted-by: alpha\n"), 1), 0600); err != nil {
				t.Fatal(err)
			}
			bin := t.TempDir()
			var config []fakeAgentConfig
			for _, id := range idea.Participants {
				config = append(config, fakeAgentConfig{ID: id, Path: writeFakeSignoffCLI(t, bin, id, "accept", 0)})
			}
			writeAgentsLocalConfig(t, root, config...)
			path := runmanifest.Path(root, run)
			switch kind {
			case "missing-manifest":
				if err := os.Remove(path); err != nil {
					t.Fatal(err)
				}
			case "corrupt-manifest":
				if err := os.WriteFile(path, []byte("{"), 0600); err != nil {
					t.Fatal(err)
				}
			case "foreign-manifest":
				m, err := runmanifest.Load(root, run)
				if err != nil {
					t.Fatal(err)
				}
				m.IdeaSlug = "another-idea"
				if err := runmanifest.Write(root, run, m); err != nil {
					t.Fatal(err)
				}
			}
			var out, errs bytes.Buffer
			err = requestConsensusSignoffs(context.Background(), requestSignoffsOptions{Root: root, IdeaSlug: idea.Slug, Yes: true}, &out, &errs)
			if kind != "absent-snapshot" {
				if err == nil {
					t.Fatal("invalid prior manifest admitted signoffs")
				}
				if kind == "missing-manifest" && !os.IsNotExist(err) ||
					kind == "corrupt-manifest" && !strings.Contains(err.Error(), "unexpected end of JSON input") ||
					kind == "foreign-manifest" && !strings.Contains(err.Error(), "contradictory quota manifest") {
					t.Fatalf("blocked for the wrong reason: %v", err)
				}
				if got := dropoutAppRecords(t, root); len(got) != 0 {
					t.Fatalf("child launched before integrity gate: %+v", got)
				}
				return
			}
			if err != nil {
				t.Fatalf("absent snapshot blocked signoffs: %v\n%s", err, &errs)
			}
			if !strings.Contains(errs.String(), "no roster snapshot") {
				t.Fatalf("missing snapshot diagnostic: %s", &errs)
			}
			summary, err := consensus.Status(root, idea.Slug, false)
			if err != nil || summary.Triage != consensus.TriageReady {
				t.Fatalf("signoffs: %+v %v", summary, err)
			}
			paths, err := filepath.Glob(filepath.Join(root, protocol.DeckDir, "runs", "*", "run.json"))
			if err != nil || len(paths) != 2 {
				t.Fatalf("standalone run missing: %v %v", paths, err)
			}
			for _, path := range paths {
				m, err := runmanifest.Load(root, filepath.Base(filepath.Dir(path)))
				if err != nil || len(m.RosterSnapshot) != 0 {
					t.Fatalf("invented snapshot: %+v %v", m.RosterSnapshot, err)
				}
			}
		})
	}
}
