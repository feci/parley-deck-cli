package membership

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
	"time"

	"parley-deck-cli/internal/protocol"
	"parley-deck-cli/internal/quota"
	"parley-deck-cli/internal/runmanifest"
	"parley-deck-cli/internal/store"
	"parley-deck-cli/internal/telemetry"
)

func fixture(t *testing.T, p quota.Policy) (string, string, string) {
	t.Helper()
	root := t.TempDir()
	t.Setenv("PARLEY_HOME", t.TempDir())
	if err := protocol.InitWorkspace(root); err != nil {
		t.Fatal(err)
	}
	idea, k, err := protocol.CreateIdeaWithQuota(root, "membership fixture", []string{"a", "b", "c", "d"}, nil, "deliberation", "", "run-1", &p, nil)
	if err != nil {
		t.Fatal(err)
	}
	run := "run-1"
	m := runmanifest.New(runmanifest.Options{Root: root, RunID: run, IdeaSlug: idea.Slug, Participants: idea.Participants, QuotaKickoff: k})
	if err = runmanifest.Write(root, run, m); err != nil {
		t.Fatal(err)
	}
	s := store.New(filepath.Join(root, protocol.DeckDir, "runs", run))
	if err = s.AppendDurable(store.Event{Type: "run.created", Data: map[string]any{"idea": idea.Slug, "participants": idea.Participants, "quota_kickoff": k}}); err != nil {
		t.Fatal(err)
	}
	for _, id := range []string{"a", "b"} {
		validRound(t, idea.Path, idea.Slug, id)
	}
	return root, idea.Path, run
}
func validRound(t *testing.T, dir, slug, id string) {
	t.Helper()
	path := filepath.Join(dir, "round-01", id+".md")
	s := fmt.Sprintf("---\nagent: %s\nidea: %s\nround: 1\n---\n", id, slug)
	for _, h := range []string{"Summary", "Proposed approach", "Concerns / open questions", "Risks", "Existing alternatives"} {
		s += "\n## " + h + "\nConcrete fixture evidence.\n"
	}
	if err := os.WriteFile(path, []byte(s), 0600); err != nil {
		t.Fatal(err)
	}
}
func candidates() []quota.Member {
	now := time.Date(2026, 10, 4, 0, 0, 0, 0, time.UTC)
	reset := now.Add(2 * time.Hour)
	e := quota.Evidence{InvocationID: "inv-c", Adapter: "synthetic-test", Eligible: true, RuleID: "test", Provenance: "synthetic-test", ObservedAt: now, ResetAt: &reset, RawReset: "2h", Excerpt: "weekly limit exhausted"}
	f := e
	f.InvocationID = "inv-d"
	return []quota.Member{{ID: "a", Usable: true}, {ID: "b", Usable: true}, {ID: "c", Evidence: &e}, {ID: "d", Evidence: &f}}
}
func leaseFixture(t *testing.T, dir, run string) context.Context {
	t.Helper()
	ctx, release, err := Acquire(context.Background(), dir, run)
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(release)
	return ctx
}
func TestQuotaTransitionReplayHistoryAndIncompletePreservation(t *testing.T) {
	root, dir, run := fixture(t, quota.NewPolicy(nil, nil))
	ctx := leaseFixture(t, dir, run)
	stub := filepath.Join(dir, "round-01", "c.md")
	os.WriteFile(stub, []byte("---\nagent: c\n---\n"), 0600)
	eventPath := filepath.Join(root, protocol.DeckDir, "runs", run, "events.jsonl")
	s := store.New(filepath.Dir(eventPath))
	s.AppendDurable(store.Event{Type: "agent.failed", Data: map[string]any{"agent": "c", "invocation_id": "inv-c"}})
	s.AppendDurable(store.Event{Type: "round.incomplete", Data: map[string]any{"idea": filepath.Base(dir), "round": "round-01"}})
	before, _ := os.ReadFile(eventPath)
	artifact, _ := os.ReadFile(filepath.Join(dir, "round-01", "a.md"))
	b, err := Settle(ctx, root, dir, run, "round-01", []string{"a", "b", "c", "d"}, candidates())
	if err != nil || b == nil {
		t.Fatal(b, err)
	}
	for i := 0; i < 2; i++ {
		if err = quota.CommitBatch(dir, *b); err != nil {
			t.Fatal(err)
		}
		if _, err = Before(ctx, root, dir, run); err != nil {
			t.Fatal(err)
		}
	}
	h, err := quota.ReadHistory(dir)
	if err != nil || h.Revision != 1 || !reflect.DeepEqual(h.Current, []string{"a", "b"}) || len(h.Known) != 4 {
		t.Fatal(h, err)
	}
	after, _ := os.ReadFile(eventPath)
	if !bytes.HasPrefix(after, before) {
		t.Fatal("historical events rewritten")
	}
	a, _ := os.ReadFile(filepath.Join(dir, "round-01", "a.md"))
	if !bytes.Equal(a, artifact) {
		t.Fatal("artifact changed")
	}
	stubRaw, _ := os.ReadFile(stub)
	if string(stubRaw) != "---\nagent: c\n---\n" {
		t.Fatal("stub changed")
	}
	events, _ := s.Load()
	completed := 0
	for _, e := range events {
		if e.Type == "round.completed" {
			completed++
			if e.Data["quota_transition"] != b.ID {
				t.Fatal(e)
			}
		}
	}
	if completed != 1 {
		t.Fatal(completed)
	}
	notices, _ := filepath.Glob(filepath.Join(root, protocol.DeckDir, "inbox", "parley-to-user_batch-*.md"))
	if len(notices) != 1 {
		t.Fatal(notices)
	}
	view, err := runmanifest.Load(root, run)
	if err != nil || view.QuotaRevision != 1 || !reflect.DeepEqual(view.Participants, h.Current) {
		t.Fatal(view, err)
	}
}
func TestQuotaProjectionFaultsPendingReadOnlyAndRecovery(t *testing.T) {
	for _, stage := range []string{"prompt", "manifest", "evaluation", "applied"} {
		t.Run(stage, func(t *testing.T) {
			root, dir, run := fixture(t, quota.NewPolicy(nil, nil))
			ctx := leaseFixture(t, dir, run)
			projectionFault = func(s string) error {
				if s == stage {
					return fmt.Errorf("injected %s", s)
				}
				return nil
			}
			defer func() { projectionFault = func(string) error { return nil } }()
			b, err := Settle(ctx, root, dir, run, "round-01", []string{"a", "b", "c", "d"}, candidates())
			if err == nil || b == nil {
				t.Fatal(b, err)
			}
			if _, err = Before(ctx, root, dir, run); err == nil {
				t.Fatal("repeated projection fault accepted")
			}
			blocks, _ := filepath.Glob(filepath.Join(root, protocol.DeckDir, "inbox", "parley-to-user_quota-block-*.md"))
			if len(blocks) != 1 {
				t.Fatal("recovery duplicated blocking escalation", blocks)
			}
			prompt, _ := os.ReadFile(filepath.Join(dir, "00-prompt.md"))
			manifest, _ := os.ReadFile(runmanifest.Path(root, run))
			history, _ := os.ReadFile(filepath.Join(dir, quota.HistoryDir, "000001.json"))
			for i := 0; i < 2; i++ {
				v, e := protocol.InspectQuota(dir)
				if e != nil || v.Pending == "" || !reflect.DeepEqual(v.History.Current, []string{"a", "b"}) {
					t.Fatal(v, e)
				}
				if _, _, e = protocol.QuotaMembers(dir, nil); e == nil {
					t.Fatal("pending signoff/close passed")
				}
			}
			for path, want := range map[string][]byte{filepath.Join(dir, "00-prompt.md"): prompt, runmanifest.Path(root, run): manifest, filepath.Join(dir, quota.HistoryDir, "000001.json"): history} {
				got, _ := os.ReadFile(path)
				if !bytes.Equal(got, want) {
					t.Fatal("read repaired", path)
				}
			}
			projectionFault = func(string) error { return nil }
			if _, err = Before(ctx, root, dir, run); err != nil {
				t.Fatal(err)
			}
			if _, _, err = protocol.QuotaMembers(dir, nil); err != nil {
				t.Fatal(err)
			}
		})
	}
}
func TestQuotaStubSurvivorAndProtectedRolesBlockWholeBatch(t *testing.T) {
	for _, kind := range []string{"stub", "designee", "pin", "drafter"} {
		t.Run(kind, func(t *testing.T) {
			root, dir, run := fixture(t, quota.NewPolicy(nil, nil))
			ctx := leaseFixture(t, dir, run)
			switch kind {
			case "stub":
				os.WriteFile(filepath.Join(dir, "round-01", "a.md"), []byte("stub"), 0600)
			case "designee":
				p := filepath.Join(dir, "00-prompt.md")
				b, _ := os.ReadFile(p)
				os.WriteFile(p, []byte(strings.Replace(string(b), "status:", "implementer: c\nstatus:", 1)), 0600)
			case "pin":
				os.WriteFile(filepath.Join(dir, "IMPLEMENTATION.md"), []byte("---\nimplementer: c\n---\n"), 0600)
			case "drafter":
				os.WriteFile(filepath.Join(dir, "consensus.md"), []byte("---\ndrafted-by: c\n---\nStarted draft"), 0600)
			}
			for i := 0; i < 2; i++ {
				b, err := Settle(ctx, root, dir, run, "round-01", []string{"a", "b", "c", "d"}, candidates())
				if err == nil || b != nil {
					t.Fatal(b, err)
				}
			}
			h, _ := quota.ReadHistory(dir)
			if h.Revision != 0 || len(h.Current) != 4 {
				t.Fatal(h)
			}
			notes, _ := filepath.Glob(filepath.Join(root, protocol.DeckDir, "inbox", "parley-to-user_quota-block-*.md"))
			if len(notes) != 1 {
				t.Fatal(notes)
			}
		})
	}
}
func TestQuotaLifetimeLockDifferentRunAndOffScope(t *testing.T) {
	for _, p := range []quota.Policy{quota.NewPolicy(nil, nil), {Enabled: false, Scope: quota.KickoffAndMidIdea}, {Enabled: true, Scope: quota.KickoffOnly}} {
		_, dir, run := fixture(t, p)
		ctx := leaseFixture(t, dir, run)
		nested, release, err := Acquire(ctx, dir, run)
		if err != nil {
			t.Fatal(err)
		}
		release()
		_ = nested
		_, release, err = Acquire(context.Background(), dir, "run-2")
		if p.Enabled && p.Scope == quota.KickoffAndMidIdea {
			if err == nil {
				release()
				t.Fatal("different run admitted")
			}
		} else {
			if err != nil {
				t.Fatal(err)
			}
			release()
		}
	}
}
func TestQuotaCrossProcessLease(t *testing.T) {
	if dir := os.Getenv("PARLEY_QUOTA_LOCK_FIXTURE"); dir != "" {
		_, release, err := Acquire(context.Background(), dir, "child-run")
		if err == nil {
			release()
			os.Exit(19)
		}
		if !strings.Contains(err.Error(), "lock held") {
			os.Exit(20)
		}
		return
	}
	_, dir, run := fixture(t, quota.NewPolicy(nil, nil))
	leaseFixture(t, dir, run)
	cmd := exec.Command(os.Args[0], "-test.run=^TestQuotaCrossProcessLease$")
	cmd.Env = append(os.Environ(), "PARLEY_QUOTA_LOCK_FIXTURE="+dir)
	if raw, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("child lock check: %v %s", err, raw)
	}
}
func TestQuotaLaterRunReadsHistoryAndContradictionsBlock(t *testing.T) {
	root, dir, run := fixture(t, quota.NewPolicy(nil, nil))
	ctx, release, err := Acquire(context.Background(), dir, run)
	if err != nil {
		t.Fatal(err)
	}
	if _, err = Settle(ctx, root, dir, run, "round-01", []string{"a", "b", "c", "d"}, candidates()); err != nil {
		t.Fatal(err)
	}
	release()
	h, _ := quota.ReadHistory(dir)
	m := runmanifest.New(runmanifest.Options{Root: root, RunID: "run-2", IdeaSlug: filepath.Base(dir), Participants: h.Kickoff.Participants, QuotaKickoff: h.Kickoff})
	if err = runmanifest.Write(root, "run-2", m); err != nil {
		t.Fatal(err)
	}
	ctx = leaseFixture(t, dir, "run-2")
	if _, err = Before(ctx, root, dir, "run-2"); err != nil {
		t.Fatal(err)
	}
	m, _ = runmanifest.Load(root, "run-2")
	if m.QuotaRevision != 1 || len(m.Participants) != 2 {
		t.Fatal(m)
	}
	p := filepath.Join(dir, "00-prompt.md")
	raw, _ := os.ReadFile(p)
	os.WriteFile(p, []byte(strings.Replace(string(raw), "participants: [a, b]", "participants: [a, intruder]", 1)), 0600)
	if _, err = Before(ctx, root, dir, "run-2"); err == nil {
		t.Fatal("contradiction repaired")
	}
}

func TestQuotaUnsettledWriterFromPriorRunBlocksRecovery(t *testing.T) {
	root, dir, run := fixture(t, quota.NewPolicy(nil, nil))
	ctx := leaseFixture(t, dir, run)
	inv, err := telemetry.Begin(filepath.Join(root, ".parley-runtime", "invocations"), telemetry.Metadata{Idea: filepath.Base(dir), Agent: "c", RunID: "old-run", Adapter: "zcode"})
	if err != nil {
		t.Fatal(err)
	}
	if err = inv.Started(424242); err != nil {
		t.Fatal(err)
	}
	if _, err := Before(ctx, root, dir, run); err == nil || !strings.Contains(err.Error(), "unsettled writer") {
		t.Fatal(err)
	}
	if err := inv.Finish(telemetry.Outcome{Status: "failed"}); err != nil {
		t.Fatal(err)
	}
	if _, err := Before(ctx, root, dir, run); err != nil {
		t.Fatal(err)
	}
}

func TestQuotaReviewDiversityAndStrictGatesOnProspectiveMembers(t *testing.T) {
	for _, kind := range []string{"auto-two-reviewers", "model-diversity", "fast-diversity", "strict-retains-minor"} {
		t.Run(kind, func(t *testing.T) {
			root, dir, run := fixture(t, quota.NewPolicy(nil, nil))
			ctx := leaseFixture(t, dir, run)
			path := filepath.Join(dir, "00-prompt.md")
			raw, _ := os.ReadFile(path)
			field := "auto_implement: true\n"
			switch kind {
			case "model-diversity":
				field = "require_model_diversity: true\n"
			case "fast-diversity":
				field = ""
				raw = bytes.Replace(raw, []byte("track: deliberation"), []byte("track: fast"), 1)
			case "strict-retains-minor":
				field = "strict_gate: true\n"
			}
			raw = bytes.Replace(raw, []byte("status:"), []byte(field+"status:"), 1)
			os.WriteFile(path, raw, 0600)
			if kind == "strict-retains-minor" {
				os.WriteFile(filepath.Join(dir, "round-01", "c.md"), []byte("### [MINOR] Still requires disposition\n"), 0600)
			}
			b, err := Settle(ctx, root, dir, run, "round-01", []string{"a", "b", "c", "d"}, candidates())
			if kind == "strict-retains-minor" {
				if err != nil || b == nil || len(b.Retained) != 1 || b.Retained[0].Kind != "finding" {
					t.Fatal(b, err)
				}
			} else {
				if err == nil || b != nil || !IsBlocked(err) {
					t.Fatal(b, err)
				}
				h, _ := quota.ReadHistory(dir)
				if h.Revision != 0 {
					t.Fatal(h)
				}
			}
		})
	}
}
