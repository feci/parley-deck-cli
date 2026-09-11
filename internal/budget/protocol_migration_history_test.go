package budget

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

func protocolRunFixture(t *testing.T, root, name, events, cursor string) {
	t.Helper()
	migrationRunFixture(t, root, events, cursor)
	if name != "historical" {
		base := filepath.Join(root, "parley-deck", "runs")
		if err := os.Rename(filepath.Join(base, "historical"), filepath.Join(base, name)); err != nil {
			t.Fatal(err)
		}
	}
}

func TestProtocolMigrationHistoryDistinctStepsAndLifetimeCounters(t *testing.T) {
	ctx := context.Background()
	root := t.TempDir()
	for run, count := range []int{2, 3} {
		name := fmt.Sprintf("run-%d", run)
		events := `{"time":"TIME","type":"run.created","data":{"idea":"idea"}}` + "\n"
		for n := 0; n < count; n++ {
			events += fmt.Sprintf(`{"time":"TIME","type":"run.phase","data":{"run_id":%q,"action":"promoted","current_round":%d}}`+"\n", name, n+2)
		}
		protocolRunFixture(t, root, name, events, "")
	}
	i, err := InspectProtocolMigration(ctx, root, "idea", DriverStep, "")
	if err != nil || i.LowerBound != 5 {
		t.Fatalf("distinct published steps lost: floor=%d err=%v", i.LowerBound, err)
	}
	r := protocolMigrationFixture(t, root, DriverStep)
	r.TotalActions = 4
	if _, err := MigrateProtocolBudget(ctx, root, "idea", DriverStep, r); err == nil {
		t.Fatal("operator total below independently known actions accepted")
	}
	// A later lifetime counter can overlap all five actions; do not add it.
	protocolRunFixture(t, root, "counter", `{"time":"TIME","type":"loop.budget","data":{"idea":"idea","steps":8,"elapsed_ms":7200000}}`, "")
	i, err = InspectProtocolMigration(ctx, root, "idea", DriverStep, "")
	if err != nil || i.LowerBound != 8 || i.Earliest == nil || i.Earliest.After(time.Now().Add(-2*time.Hour)) {
		t.Fatalf("counter/origin: %+v %v", i, err)
	}
}

func TestProtocolMigrationHistoryGroupedChildrenAndRefusals(t *testing.T) {
	for _, kind := range []Kind{DriverStep, Fixup, CrossReview} {
		t.Run(string(kind), func(t *testing.T) {
			root := t.TempDir()
			phase := "fixup"
			if kind == CrossReview {
				phase = "round-02"
			}
			historicalProtocolCall(t, root, phase, false)
			i, err := InspectProtocolMigration(context.Background(), root, "idea", kind, "")
			if err != nil || i.LowerBound != 0 || len(i.UngroupedInvocations) != 1 {
				t.Fatalf("refusal accounting: %+v %v", i, err)
			}
			for n := 0; n < 3; n++ {
				historicalProtocolCall(t, root, phase, true)
			}
			i, err = InspectProtocolMigration(context.Background(), root, "idea", kind, "")
			if err != nil || i.LowerBound != 1 || len(i.UngroupedInvocations) != 4 {
				t.Fatalf("children became separate protocol charges: %+v %v", i, err)
			}
		})
	}
}

func TestProtocolMigrationHistoryRejectsAmbiguousCountersAndIdentity(t *testing.T) {
	created := `{"time":"TIME","type":"run.created","data":{"idea":"idea"}}` + "\n"
	for _, tc := range []struct {
		name, event, cursor string
		kind                Kind
	}{
		{"missing-step", `{"time":"TIME","type":"loop.budget","data":{}}`, "", DriverStep},
		{"null-step", `{"time":"TIME","type":"loop.budget","data":{"steps":null}}`, "", DriverStep},
		{"alias-step", `{"time":"TIME","type":"loop.budget","data":{"Steps":1}}`, "", DriverStep},
		{"negative", `{"time":"TIME","type":"loop.budget","data":{"steps":-1}}`, "", DriverStep},
		{"fraction", `{"time":"TIME","type":"loop.budget","data":{"steps":1.5}}`, "", DriverStep},
		{"elapsed-null", `{"time":"TIME","type":"loop.budget","data":{"steps":1,"elapsed_ms":null}}`, "", DriverStep},
		{"elapsed-overflow", `{"time":"TIME","type":"loop.budget","data":{"steps":1,"elapsed_ms":9223372036854775807}}`, "", DriverStep},
		{"run-id", `{"time":"TIME","type":"run.phase","data":{"run_id":"different","action":"fixup"}}`, "", Fixup},
		{"missing-action", `{"time":"TIME","type":"run.phase","data":{}}`, "", DriverStep},
		{"unknown-action", `{"time":"TIME","type":"run.phase","data":{"action":"unknown"}}`, "", Fixup},
		{"cursor-null", "", `{"schema_version":2,"fixup_cycles_published":null}`, Fixup},
		{"cursor-missing", "", `{"schema_version":2}`, Fixup},
		{"schema-new", "", `{"schema_version":3,"rounds_run":1}`, CrossReview},
		{"cursor-identity", "", `{"idea":"different","rounds_run":1}`, CrossReview},
	} {
		t.Run(tc.name, func(t *testing.T) {
			root := t.TempDir()
			migrationRunFixture(t, root, created+tc.event, tc.cursor)
			if _, err := InspectProtocolMigration(context.Background(), root, "idea", tc.kind, ""); err == nil {
				t.Fatal("ambiguous history accepted")
			}
		})
	}
}

func TestProtocolMigrationHistoryCursorIdentityAndCycleRecovery(t *testing.T) {
	ctx := context.Background()
	root := t.TempDir()
	migrationRunFixture(t, root, `{"time":"TIME","type":"run.phase","data":{"action":"fixup"}}`, `{"idea":"idea","schema_version":2,"fixup_cycles_published":4,"rounds_run":3}`)
	for _, tc := range []struct {
		kind Kind
		want int
	}{{DriverStep, 1}, {Fixup, 4}, {CrossReview, 2}} {
		i, err := InspectProtocolMigration(ctx, root, "idea", tc.kind, "")
		if err != nil || i.LowerBound != tc.want {
			t.Fatalf("cursor-only identity excluded for %s: %d %v", tc.kind, i.LowerBound, err)
		}
	}
	// A fixup phase event may finish a previously charged cycle after a crash.
	// With no operation IDs, two such events cannot prove two distinct fixups.
	other := t.TempDir()
	event := `{"time":"TIME","type":"run.phase","data":{"idea":"idea","action":"fixup","current_round":1}}` + "\n"
	migrationRunFixture(t, other, event+strings.Replace(event, `"current_round":1`, `"current_round":2`, 1), "")
	i, err := InspectProtocolMigration(ctx, other, "idea", Fixup, "")
	if err != nil || i.LowerBound != 1 {
		t.Fatalf("recovery invented another fixup: %d %v", i.LowerBound, err)
	}
}

func TestProtocolMigrationHistoryNestedArtifactsAndReadOnly(t *testing.T) {
	root := t.TempDir()
	ctx := context.Background()
	relative := "parley-deck/pipelines/p/blocks/b/idea"
	i, err := InspectProtocolMigration(ctx, root, "idea", Fixup, relative)
	if err != nil {
		t.Fatal(err)
	}
	entries, err := os.ReadDir(root)
	if err != nil || len(entries) != 0 {
		t.Fatal("read-only inventory mutated root", err)
	}
	marker := filepath.Join(root, filepath.FromSlash(relative), "review", "round-01", ".fixup-done")
	if err := os.MkdirAll(filepath.Dir(marker), 0700); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(marker, []byte("retained"), 0600); err != nil {
		t.Fatal(err)
	}
	j, err := InspectProtocolMigration(ctx, root, "idea", Fixup, relative)
	if err != nil || j.LowerBound != 1 || j.HistorySHA256 == i.HistorySHA256 || len(j.History.Sources) != 1 {
		t.Fatalf("nested source omitted: %+v %v", j, err)
	}
	for _, path := range []string{"../escape", "parley-deck/../escape", "parley-deck/ideas/", "parley-deck\\ideas\\idea"} {
		if _, err := InspectProtocolMigration(ctx, root, "idea", Fixup, path); err == nil {
			t.Fatal("noncanonical path accepted", path)
		}
	}
}

func TestProtocolMigrationHistoryWorktreeCopiesAndDisjointMarkers(t *testing.T) {
	ctx, root := context.Background(), t.TempDir()
	git := func(args ...string) {
		t.Helper()
		c := exec.Command("git", args...)
		c.Dir = root
		if out, err := c.CombinedOutput(); err != nil {
			t.Fatalf("git %v: %s %v", args, out, err)
		}
	}
	git("init", "-q")
	git("-c", "user.name=Fixture", "-c", "user.email=fixture@example.invalid", "commit", "-q", "--allow-empty", "-m", "fixture")
	other := filepath.Join(t.TempDir(), "other")
	git("worktree", "add", "-q", "-b", "other", other)
	migrationRunFixture(t, root, `{"time":"TIME","type":"run.phase","data":{"idea":"idea","action":"promoted"}}`, `{"schema_version":2,"fixup_cycles_published":1,"rounds_run":2}`)
	base := "parley-deck/runs/historical"
	if err := os.MkdirAll(filepath.Join(other, filepath.FromSlash(base)), 0700); err != nil {
		t.Fatal(err)
	}
	for _, name := range []string{"events.jsonl", "driver.json"} {
		raw, err := os.ReadFile(filepath.Join(root, filepath.FromSlash(base), name))
		if err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(filepath.Join(other, filepath.FromSlash(base), name), raw, 0600); err != nil {
			t.Fatal(err)
		}
	}
	for n, origin := range []string{root, other} {
		marker := filepath.Join(origin, "parley-deck", "ideas", "idea", "review", fmt.Sprintf("round-%02d", n+1), ".fixup-done")
		if err := os.MkdirAll(filepath.Dir(marker), 0700); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(marker, []byte("fixture"), 0600); err != nil {
			t.Fatal(err)
		}
	}
	for _, tc := range []struct {
		kind Kind
		want int
	}{{DriverStep, 1}, {Fixup, 2}, {CrossReview, 1}} {
		i, err := InspectProtocolMigration(ctx, root, "idea", tc.kind, "")
		if err != nil || i.LowerBound != tc.want {
			t.Fatalf("copy/marker accounting: floor=%d %v", i.LowerBound, err)
		}
		j, err := InspectProtocolMigration(ctx, other, "idea", tc.kind, "")
		if err != nil || j.HistorySHA256 != i.HistorySHA256 {
			t.Fatal("worktree-dependent digest", err)
		}
	}
	migrationRewrite(t, filepath.Join(other, filepath.FromSlash(base), "driver.json"), func(v map[string]any) { v["rounds_run"] = 3 })
	if _, err := InspectProtocolMigration(ctx, root, "idea", DriverStep, ""); err == nil {
		t.Fatal("divergent run copies accepted")
	}
}

func TestProtocolMigrationHistoryMandatoryFieldsAndCountConsistency(t *testing.T) {
	ctx, root := context.Background(), t.TempDir()
	historicalProtocolCall(t, root, "fixup", true)
	i, err := InspectProtocolMigration(ctx, root, "idea", Fixup, "")
	if err != nil {
		t.Fatal(err)
	}
	r := protocolMigrationFixture(t, root, Fixup)
	r.IdeaPath = i.IdeaPath
	for _, mutation := range []string{"floor", "evidence-count", "evidence-rule", "earliest"} {
		t.Run(mutation, func(t *testing.T) {
			j := i
			j.Evidence = append([]ProtocolCountEvidence{}, i.Evidence...)
			switch mutation {
			case "floor":
				j.LowerBound = 0
			case "evidence-count":
				j.Evidence[0].Floor = 2
				j.LowerBound = 2
			case "evidence-rule":
				j.Evidence[0].Rule = "unknown"
			case "earliest":
				j.Earliest = nil
			}
			j.HistorySHA256 = j.digest()
			q := r
			q.ExpectedHistorySHA256 = j.HistorySHA256
			q.TotalActions = 3
			if _, err := protocolMigrationInitial(j, q); err == nil {
				t.Fatal("inconsistent frozen inventory accepted")
			}
		})
	}
}
