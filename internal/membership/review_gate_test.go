package membership

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"parley-deck-cli/internal/protocol"
	"parley-deck-cli/internal/quota"
	"parley-deck-cli/internal/quotatest"
	"parley-deck-cli/internal/runmanifest"
)

func reviewGateFixture(t *testing.T) (string, string, string) {
	t.Helper()
	root, dir, run := fixture(t, quota.NewParticipantPolicy(nil, nil))
	path := filepath.Join(dir, "00-prompt.md")
	raw, _ := os.ReadFile(path)
	os.WriteFile(path, []byte(strings.Replace(string(raw), "status:", "auto_implement: true\nimplementer: a\nrequire_model_diversity: false\nstatus:", 1)), 0600)
	reviewModels(t, root, run, "Model-A", "Model-B")
	return root, dir, run
}

func reviewModels(t *testing.T, root, run, a, b string) {
	t.Helper()
	m, err := runmanifest.Load(root, run)
	if err != nil {
		t.Fatal(err)
	}
	m.RosterSnapshot = []runmanifest.RosterSnapshotEntry{{Agent: "a", Model: a}, {Agent: "b", Model: b}}
	if err = runmanifest.Write(root, run, m); err != nil {
		t.Fatal(err)
	}
}

func TestReviewGateProspectiveCommitAndPolicyRevision(t *testing.T) {
	root, dir, run := reviewGateFixture(t)
	ctx := leaseFixture(t, dir, run)
	b, err := Settle(ctx, root, dir, run, "round-01", []string{"a", "b", "c", "d"}, dropoutMembers(t, filepath.Base(dir)))
	if err != nil || b == nil {
		t.Fatalf("prospective reduction: %+v %v", b, err)
	}
	if err := CheckGates(root, dir, run, []string{"a", "b"}); err != nil {
		t.Fatal(err)
	}
	off := quota.Policy{Enabled: false, Scope: quota.KickoffAndMidIdea}
	a := quotatest.Authority(t, root, filepath.Base(dir), "off", (quota.RevisionDirective{Participants: []string{"a", "b"}, Policy: off}).Text())
	if _, err := Revise(ctx, root, dir, run, []string{"a", "b"}, off, a, nil); err != nil {
		t.Fatal(err)
	}
	if ok, err := SingleReviewerAfterDropout(root, dir, run, "a", []string{"b", "a"}, nil); !ok || err != nil {
		t.Fatalf("policy-only revision lost cause: %v %v", ok, err)
	}
	// A subsequent owner membership edit cannot spend an older automatic loss.
	validRound(t, dir, filepath.Base(dir), "new")
	catchupPath := filepath.Join(dir, "round-01", "new.md")
	catchupRaw, _ := os.ReadFile(catchupPath)
	os.WriteFile(catchupPath, []byte(strings.Replace(string(catchupRaw), "round: 1\n", "round: 1\ncatch-up: true\njoin-from: round-02\nread-priors: [round-01/a.md]\n", 1)), 0600)
	a = quotatest.Authority(t, root, filepath.Base(dir), "replace", (quota.RevisionDirective{Participants: []string{"a", "new"}, Policy: off}).Text())
	if _, err := Revise(ctx, root, dir, run, []string{"a", "new"}, off, a, map[string]string{"new": "round-01/new.md"}); err != nil {
		t.Fatal(err)
	}
	if ok, err := SingleReviewerAfterDropout(root, dir, run, "a", []string{"a", "new"}, nil); ok || err != nil {
		t.Fatalf("manual edit earned credit: %v %v", ok, err)
	}
}

func TestReviewGateSnapshotModelsAndNoEvidence(t *testing.T) {
	for _, model := range []string{"", "unknown", "cli-default", " model-A ", "mOdEl-a", " model-B "} {
		t.Run(model, func(t *testing.T) {
			root, dir, run := reviewGateFixture(t)
			reviewModels(t, root, run, " Model-A ", model)
			ctx := leaseFixture(t, dir, run)
			_, err := Settle(ctx, root, dir, run, "round-01", []string{"a", "b", "c", "d"}, dropoutMembers(t, filepath.Base(dir)))
			if (err == nil) != (model == " model-B ") {
				t.Fatalf("model %q: %v", model, err)
			}
		})
	}
	root, dir, run := reviewGateFixture(t)
	if ok, err := SingleReviewerAfterDropout(root, dir, run, "a", []string{"a", "b"}, nil); ok || err == nil {
		t.Fatalf("prospective IDs without decision accepted: %v %v", ok, err)
	}
	// A two-person-by-design legacy prompt and confirmed marker supply no history.
	os.Remove(filepath.Join(dir, quota.KickoffFile))
	os.WriteFile(filepath.Join(dir, "00-prompt.md"), []byte("---\nidea: "+filepath.Base(dir)+"\nparticipants: [a, b]\nauto_implement: true\nimplementer: a\nexcluded: c — provider failed — confirmed 2026-10-09\n---\n"), 0600)
	if ok, err := SingleReviewerAfterDropout(root, dir, run, "a", []string{"a", "b"}, nil); ok || err != nil {
		t.Fatalf("marker accepted: %v %v", ok, err)
	}
	if err := CheckGates(root, dir, run, []string{"a", "b"}); err == nil {
		t.Fatal("unrecorded auto single reviewer passed")
	}
}

func TestReviewGateProspectiveAdversarialEvidence(t *testing.T) {
	for _, mutation := range []string{"arbitrary-rule", "wrong-agent", "missing-attempt", "duplicate-candidate", "added-id", "wrong-before", "not-applied", "protected", "unknown-snapshot", "pending", "corrupt"} {
		t.Run(mutation, func(t *testing.T) {
			root, dir, run := reviewGateFixture(t)
			d := quota.Evaluate(quota.NewParticipantPolicy(nil, nil), []string{"a", "b", "c", "d"}, dropoutMembers(t, filepath.Base(dir)), quota.Roles{})
			switch mutation {
			case "arbitrary-rule":
				d.Candidates[0].Evidence.RuleID = "test"
			case "wrong-agent":
				d.Candidates[0].Evidence.Failure.Agent = "other"
			case "missing-attempt":
				d.Candidates[0].Evidence.Failure.Attempts = d.Candidates[0].Evidence.Failure.Attempts[:1]
			case "duplicate-candidate":
				d.Candidates[1] = d.Candidates[0]
			case "added-id":
				d.After = []string{"a", "other"}
			case "wrong-before":
				d.Before = []string{"a", "b", "c"}
			case "not-applied":
				d.Applied = false
			case "protected":
				os.WriteFile(filepath.Join(dir, "consensus.md"), []byte("---\ndrafted-by: c\n---\n"), 0600)
			case "unknown-snapshot":
				reviewModels(t, root, run, "", "model-b")
			case "pending":
				h, _ := quota.ReadHistory(dir)
				b := quota.NewBatch(h, run, "round-01", h.Current, d, time.Now())
				if err := quota.CommitBatch(dir, b); err != nil {
					t.Fatal(err)
				}
			case "corrupt":
				os.WriteFile(filepath.Join(dir, quota.KickoffFile), []byte("{"), 0600)
			}
			if ok, err := SingleReviewerAfterDropout(root, dir, run, "a", d.After, &d); ok {
				t.Fatalf("%s earned exception (%v)", mutation, err)
			}
		})
	}
}

func TestReviewGateKickoffAndLegacyRules(t *testing.T) {
	for _, kind := range []string{"dropout", legacyResetRule, legacyAccountRule, legacyAllowanceRule, "test"} {
		t.Run(kind, func(t *testing.T) {
			root, dir, run := reviewGateFixture(t)
			p := quota.NewParticipantPolicy(nil, nil)
			ms := dropoutMembers(t, "readiness-batch")
			for i := range ms {
				if e := ms[i].Evidence; e != nil {
					if kind == "dropout" {
						e.Failure.Step = "readiness"
					} else {
						p = quota.NewPolicy(nil, nil)
						e.RuleID, e.Provenance, e.Adapter, e.Failure = kind, legacyProvenance, "zcode", nil
						if kind == legacyResetRule {
							reset := e.ObservedAt.Add(time.Hour)
							e.ResetAt, e.RawReset = &reset, "3600"
						}
					}
				}
			}
			d := quota.Evaluate(p, []string{"a", "b", "c", "d"}, ms, quota.Roles{})
			os.RemoveAll(dir)
			idea, k, err := protocol.CreateIdeaWithQuota(root, "kickoff review gate", d.After, nil, "deliberation", "", run, &p, &d)
			if err != nil {
				t.Fatal(err)
			}
			m := runmanifest.New(runmanifest.Options{Root: root, RunID: run, IdeaSlug: idea.Slug, Participants: k.Participants, QuotaKickoff: k})
			if err := runmanifest.Write(root, run, m); err != nil {
				t.Fatal(err)
			}
			reviewModels(t, root, run, "model-a", "model-b")
			ok, err := SingleReviewerAfterDropout(root, idea.Path, run, "a", d.After, nil)
			if err != nil || ok != (kind != "test") {
				t.Fatalf("%s: %v %v", kind, ok, err)
			}
		})
	}
}
