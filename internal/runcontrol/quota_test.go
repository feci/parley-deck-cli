package runcontrol

import (
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	"parley-deck-cli/internal/protocol"
	"parley-deck-cli/internal/quota"
	"parley-deck-cli/internal/runmanifest"
	"parley-deck-cli/internal/runstate"
)

func TestQuotaC1KickoffReplayAndFrozenScope(t *testing.T) {
	root := t.TempDir()
	t.Setenv("PARLEY_HOME", t.TempDir())
	if err := protocol.InitWorkspace(root); err != nil {
		t.Fatal(err)
	}
	off := false
	p := quota.NewPolicy(&off, nil)
	p.Scope = quota.KickoffOnly
	run, err := Create(CreateOptions{Root: root, Task: "quota kickoff", Participants: []string{"a", "b", "c"}, Excluded: []string{"c — missing — confirmed 2026-10-04"}, QuotaPolicy: &p})
	if err != nil {
		t.Fatal(err)
	}
	want := []string{"a", "b"}
	if !reflect.DeepEqual(run.Idea.Participants, want) || !reflect.DeepEqual(run.RunOptions.Idea.Participants, want) {
		t.Fatal(run.Idea)
	}
	meta, err := protocol.ReadFrontmatter(filepath.Join(run.Idea.Path, "00-prompt.md"))
	if err != nil || meta["participants"] != "[a, b]" {
		t.Fatal(meta, err)
	}
	m, err := runmanifest.Load(root, run.RunID)
	if err != nil || !reflect.DeepEqual(m.Participants, want) || len(m.ActiveSteps) != 2 || m.QuotaKickoff.Policy.Scope != quota.KickoffOnly {
		t.Fatal(m, err)
	}
	events, err := run.Store.Load()
	if err != nil {
		t.Fatal(err)
	}
	for _, e := range events {
		if e.Type == "run.created" {
			if !reflect.DeepEqual(e.Data["participants"], []any{"a", "b"}) {
				t.Fatal(e)
			}
		}
	}
	for i := 0; i < 2; i++ {
		replay, err := runstate.LoadRun(root, run.RunID)
		if err != nil || replay.QuotaPending != "" || !reflect.DeepEqual(replay.Participants, want) || replay.QuotaKickoff.Policy.Enabled {
			t.Fatal(replay, err)
		}
	}
	prompt := filepath.Join(run.Idea.Path, "00-prompt.md")
	b, _ := os.ReadFile(prompt)
	os.WriteFile(prompt, []byte(strings.Replace(string(b), quota.KickoffOnly, quota.KickoffAndMidIdea, 1)), 0644)
	replay, err := runstate.LoadRun(root, run.RunID)
	if err != nil || replay.QuotaPending == "" {
		t.Fatal("scope upgrade silently accepted", replay, err)
	}
}
func TestQuotaLegacyC1NoPolicyAndNoHistoricalExcludedSigner(t *testing.T) {
	root := t.TempDir()
	t.Setenv("PARLEY_HOME", t.TempDir())
	protocol.InitWorkspace(root)
	run, err := Create(CreateOptions{Root: root, Task: "legacy c1", Participants: []string{"a", "b", "c"}, Excluded: []string{"c — offline — confirmed 2026-10-04"}})
	if err != nil {
		t.Fatal(err)
	}
	k, err := protocol.ReadQuotaState(run.Idea.Path)
	if err != nil || k != nil {
		t.Fatal(k, err)
	}
	r, err := runstate.LoadRun(root, run.RunID)
	if err != nil || !reflect.DeepEqual(r.Participants, []string{"a", "b"}) {
		t.Fatal(r, err)
	}
}

func TestQuotaAutomaticKickoffRecordsAndNotice(t *testing.T) {
	root := t.TempDir()
	t.Setenv("PARLEY_HOME", t.TempDir())
	if err := protocol.InitWorkspace(root); err != nil {
		t.Fatal(err)
	}
	policy := quota.Policy{Enabled: true, Scope: quota.KickoffOnly}
	ev := quota.Evidence{InvocationID: "fixture", Adapter: "synthetic-test", RuleID: "synthetic-test", Provenance: "synthetic-test", Eligible: true}
	decision := quota.Evaluate(policy, []string{"a", "b", "c"}, []quota.Member{{ID: "a", Usable: true}, {ID: "b", Usable: true}, {ID: "c", Evidence: &ev}}, quota.Roles{})
	run, err := Create(CreateOptions{Root: root, Task: "automatic fixture", Participants: decision.After, QuotaPolicy: &policy, QuotaDecision: &decision})
	if err != nil {
		t.Fatal(err)
	}
	k, err := protocol.ReadQuotaState(run.Idea.Path)
	if err != nil || k.Transition == nil {
		t.Fatal(k, err)
	}
	if len(k.Markers()) != 1 || !strings.Contains(k.Markers()[0], "automatic") || strings.Contains(k.Markers()[0], "confirmed") {
		t.Fatal(k.Markers())
	}
	notes, _ := filepath.Glob(filepath.Join(root, "parley-deck", "inbox", "parley-to-user_kickoff-*.md"))
	if len(notes) != 1 {
		t.Fatal(notes)
	}
	for i := 0; i < 2; i++ {
		r, err := runstate.LoadRun(root, run.RunID)
		if err != nil || r.QuotaPending != "" || !reflect.DeepEqual(r.Participants, []string{"a", "b"}) {
			t.Fatal(r, err)
		}
	}
	notes2, _ := filepath.Glob(filepath.Join(root, "parley-deck", "inbox", "parley-to-user_kickoff-*.md"))
	if !reflect.DeepEqual(notes, notes2) {
		t.Fatal(notes2)
	}
}
