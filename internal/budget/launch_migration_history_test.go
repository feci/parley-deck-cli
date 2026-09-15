package budget

import (
	"context"
	"encoding/json"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

func migrationRewrite(t *testing.T, path string, change func(map[string]any)) {
	t.Helper()
	raw, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	var v map[string]any
	if err := json.Unmarshal(raw, &v); err != nil {
		t.Fatal(err)
	}
	change(v)
	raw, err = json.Marshal(v)
	if err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, raw, 0600); err != nil {
		t.Fatal(err)
	}
}

func migrationRunFixture(t *testing.T, root, events, cursor string) {
	t.Helper()
	dir := filepath.Join(root, "parley-deck", "runs", "historical")
	if err := os.MkdirAll(dir, 0700); err != nil {
		t.Fatal(err)
	}
	if events != "" {
		events = strings.ReplaceAll(events, "TIME", time.Now().UTC().Add(-time.Minute).Format(time.RFC3339Nano))
		if err := os.WriteFile(filepath.Join(dir, "events.jsonl"), []byte(events), 0600); err != nil {
			t.Fatal(err)
		}
	}
	if cursor != "" {
		if err := os.WriteFile(filepath.Join(dir, "driver.json"), []byte(cursor), 0600); err != nil {
			t.Fatal(err)
		}
	}
}

func TestLaunchMigrationHistoryReadOnlyAndCanonicalHash(t *testing.T) {
	ctx := context.Background()
	root := t.TempDir()
	i, err := InspectLaunchMigration(ctx, root, "idea")
	if err != nil || len(i.Sources) != 0 || len(i.Launches) != 0 {
		t.Fatalf("%+v %v", i, err)
	}
	files, err := os.ReadDir(root)
	if err != nil || len(files) != 0 {
		t.Fatal("inspection mutated absent state", err)
	}
	base := filepath.Join(root, "parley-deck", "ideas", "idea")
	if err := os.MkdirAll(base, 0700); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(base, "FINAL.md"), []byte("Original canonical source"), 0600); err != nil {
		t.Fatal(err)
	}
	j, err := InspectLaunchMigration(ctx, root, "idea")
	if err != nil || !j.HasCanonicalHistory || j.HistorySHA256 == i.HistorySHA256 {
		t.Fatalf("%+v %v", j, err)
	}
	if err := os.WriteFile(filepath.Join(base, "FINAL.md"), []byte("Changed canonical source"), 0600); err != nil {
		t.Fatal(err)
	}
	k, err := InspectLaunchMigration(ctx, root, "idea")
	if err != nil || k.HistorySHA256 == j.HistorySHA256 {
		t.Fatal("canonical edit not bound", err)
	}
}

func TestLaunchMigrationHistoryRejectsContradictoryRunIdentity(t *testing.T) {
	created := `{"time":"TIME","type":"run.created","data":{"idea":"other"}}` + "\n"
	for _, tc := range []struct{ name, event, cursor string }{
		{"phase", `{"time":"TIME","type":"run.phase","data":{"idea":"idea"}}`, ""},
		{"start", `{"time":"TIME","type":"agent.started","data":{"idea":"idea"}}`, ""},
		{"cursor", "", `{"idea":"idea","phase":"review"}`},
		{"cursor-slug", "", `{"idea_slug":"idea"}`},
		{"null-cursor", "", `{"idea":null}`},
		{"aliased-cursor", "", `{"Idea":"idea"}`},
		{"malformed-cursor", "", `{`},
		{"aliased-phase", `{"time":"TIME","type":"run.phase","data":{"Idea":"idea"}}`, ""},
		{"duplicate-phase", `{"time":"TIME","type":"run.phase","data":{"idea":"other","idea":"idea"}}`, ""},
	} {
		t.Run(tc.name, func(t *testing.T) {
			root := t.TempDir()
			migrationRunFixture(t, root, created+tc.event, tc.cursor)
			if _, err := InspectLaunchMigration(context.Background(), root, "idea"); err == nil {
				t.Fatal("conflicting history was excluded")
			}
		})
	}
}

func TestLaunchMigrationHistoryLegacyFloorAndRefusalContradiction(t *testing.T) {
	ctx := context.Background()
	root := t.TempDir()
	inv := historicalMigrationCall(t, root, "failed", "budget_refused", false, nil)
	created := `{"time":"TIME","type":"run.created","data":{"idea":"idea"}}` + "\n"
	start := `{"time":"TIME","type":"agent.started","data":{"invocation_id":"ID"}}` + "\n"
	migrationRunFixture(t, root, created+strings.ReplaceAll(start, "ID", inv.ID), "")
	if _, err := InspectLaunchMigration(ctx, root, "idea"); err == nil {
		t.Fatal("start contradicting exempt refusal accepted")
	}
	migrationRunFixture(t, root, created+strings.ReplaceAll(start, "ID", "older-attempt"), `{"phase":"review"}`)
	i, err := InspectLaunchMigration(ctx, root, "idea")
	if err != nil || i.LegacyStartFloor != 1 {
		t.Fatalf("%+v %v", i, err)
	}
	r := migrationFixtureRequest(t, root)
	if _, err := MigrateLaunchBudget(ctx, root, "idea", r); err == nil {
		t.Fatal("operator count below observed floor accepted")
	}
	r.AdditionalLaunches = 1
	s, err := MigrateLaunchBudget(ctx, root, "idea", r)
	if err != nil || s.Spent != 1 || s.ExposureMicros != nil {
		t.Fatalf("legacy work lost: %+v %v", s, err)
	}
}

func TestLaunchMigrationHistoryRejectsIncompleteAndContradictoryRecords(t *testing.T) {
	for _, mutation := range []string{"no-terminal", "no-start", "pid", "metadata", "future", "tokens", "bytes", "activity", "missing-usage", "null-usage", "missing-cost", "null-status", "missing-observation", "duplicate", "alias", "oversized", "symlink"} {
		t.Run(mutation, func(t *testing.T) {
			root := t.TempDir()
			started := mutation == "no-start" || mutation == "pid"
			inv := historicalMigrationCall(t, root, "failed", "budget_refused", started, nil)
			path := filepath.Join(inv.Dir, "terminal.json")
			switch mutation {
			case "no-terminal":
				if err := os.Remove(path); err != nil {
					t.Fatal(err)
				}
			case "no-start":
				if err := os.Remove(filepath.Join(inv.Dir, "started.json")); err != nil {
					t.Fatal(err)
				}
			case "duplicate", "alias":
				raw, _ := os.ReadFile(path)
				field := `"outcome": null,`
				if mutation == "alias" {
					field = `"Outcome": null,`
				}
				if err := os.WriteFile(path, append([]byte("{"+field), raw[1:]...), 0600); err != nil {
					t.Fatal(err)
				}
			case "oversized":
				if err := os.WriteFile(path, []byte(strings.Repeat(" ", (1<<20)+1)), 0600); err != nil {
					t.Fatal(err)
				}
			case "symlink":
				target := filepath.Join(t.TempDir(), "terminal.json")
				if err := os.Rename(path, target); err != nil {
					t.Fatal(err)
				}
				if err := os.Symlink(target, path); err != nil {
					t.Skip(err)
				}
			default:
				migrationRewrite(t, path, func(v map[string]any) {
					o := v["outcome"].(map[string]any)
					switch mutation {
					case "pid":
						v["pid"] = 99
					case "metadata":
						v["metadata"].(map[string]any)["idea"] = "other"
					case "future":
						v["completed_at"] = time.Now().Add(time.Hour).UTC().Format(time.RFC3339Nano)
					case "tokens":
						o["usage"].(map[string]any)["input_tokens"] = 1
					case "bytes":
						o["observation"].(map[string]any)["stdout_bytes"] = 1
					case "activity":
						o["observation"].(map[string]any)["first_activity_ms"] = 0
					case "missing-usage":
						delete(o, "usage")
					case "null-usage":
						o["usage"] = nil
					case "missing-cost":
						delete(o["usage"].(map[string]any), "cost_usd")
					case "null-status":
						o["status"] = nil
					case "missing-observation":
						delete(o, "observation")
					}
				})
			}
			if _, err := InspectLaunchMigration(context.Background(), root, "idea"); err == nil {
				t.Fatal("invalid source accepted")
			}
		})
	}
}

func TestLaunchMigrationHistoryMissingZeroAuthorityRefuses(t *testing.T) {
	for _, mutation := range []string{"additional", "canonical", "floor", "cost", "reserve", "refusal-cost"} {
		t.Run(mutation, func(t *testing.T) {
			root := t.TempDir()
			historicalMigrationCall(t, root, "failed", "budget_refused", false, nil)
			r := migrationFixtureRequest(t, root)
			if _, err := MigrateLaunchBudget(context.Background(), root, "idea", r); err != nil {
				t.Fatal(err)
			}
			dir, _, _, err := launchScope(context.Background(), root, "idea", false)
			if err != nil {
				t.Fatal(err)
			}
			migrationRewrite(t, filepath.Join(dir, "migration.json"), func(v map[string]any) {
				i := v["inventory"].(map[string]any)
				req := v["request"].(map[string]any)
				switch mutation {
				case "additional":
					delete(req, "additional_launches")
				case "canonical":
					delete(i, "has_canonical_history")
				case "floor":
					delete(i, "legacy_start_floor")
				case "cost":
					delete(req["policy"].(map[string]any), "max_cost_micros")
				case "reserve":
					delete(req["policy"].(map[string]any), "reserve_micros")
				case "refusal-cost":
					delete(i["launches"].([]any)[0].(map[string]any), "observed_micros")
				}
			})
			if _, err := LoadLaunchBinding(context.Background(), root, "idea"); err == nil {
				t.Fatal("missing zero-valued authority normalized into valid hash")
			}
		})
	}
}

func TestLaunchMigrationHistorySharedWorktreesPreserveOneAttempt(t *testing.T) {
	ctx := context.Background()
	root := t.TempDir()
	git := func(args ...string) {
		t.Helper()
		if out, err := exec.Command("git", append([]string{"-C", root}, args...)...).CombinedOutput(); err != nil {
			t.Fatalf("git: %s %v", out, err)
		}
	}
	git("init", "-q")
	git("-c", "user.name=fixture", "-c", "user.email=fixture@example.invalid", "commit", "--allow-empty", "-qm", "fixture")
	other := filepath.Join(t.TempDir(), "other")
	git("worktree", "add", "-q", "-b", "other", other)
	cost := 0.1250001
	inv := historicalMigrationCall(t, root, "process-exited", "", true, &cost)
	copyDir := filepath.Join(other, ".parley-runtime", "invocations", inv.ID)
	if err := os.MkdirAll(copyDir, 0700); err != nil {
		t.Fatal(err)
	}
	for _, name := range []string{"requested.json", "started.json", "terminal.json"} {
		raw, err := os.ReadFile(filepath.Join(inv.Dir, name))
		if err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(filepath.Join(copyDir, name), raw, 0600); err != nil {
			t.Fatal(err)
		}
	}
	i, err := InspectLaunchMigration(ctx, root, "idea")
	if err != nil || len(i.Roots) != 2 || len(i.Launches) != 1 || *i.Launches[0].ObservedMicros != 125001 {
		t.Fatalf("%+v %v", i, err)
	}
	j, err := InspectLaunchMigration(ctx, other, "idea")
	if err != nil || j.HistorySHA256 != i.HistorySHA256 {
		t.Fatal("worktree-dependent inventory", err)
	}
	migrationRewrite(t, filepath.Join(copyDir, "terminal.json"), func(v map[string]any) { v["outcome"].(map[string]any)["usage"].(map[string]any)["cost_usd"] = 2.0 })
	if _, err := InspectLaunchMigration(ctx, root, "idea"); err == nil {
		t.Fatal("conflicting invocation copies accepted")
	}
	// Remove only this disposable fixture worktree, leaving its Git registration.
	if err := os.RemoveAll(other); err != nil {
		t.Fatal(err)
	}
	if _, err := InspectLaunchMigration(ctx, root, "idea"); err == nil {
		t.Fatal("unavailable registered worktree ignored")
	}
}

func TestLaunchMigrationHistoryDirectoryAliasesAndBounds(t *testing.T) {
	for _, relative := range []string{".parley-runtime", "parley-deck", "parley-deck/ideas", "parley-deck/ideas/idea", "parley-deck/runs"} {
		t.Run(relative, func(t *testing.T) {
			root := t.TempDir()
			path := filepath.Join(root, filepath.FromSlash(relative))
			if err := os.MkdirAll(filepath.Dir(path), 0700); err != nil {
				t.Fatal(err)
			}
			if err := os.Symlink(t.TempDir(), path); err != nil {
				t.Skip(err)
			}
			if _, err := InspectLaunchMigration(context.Background(), root, "idea"); err == nil {
				t.Fatal("empty aliased directory accepted")
			}
		})
	}
	root := t.TempDir()
	if err := os.WriteFile(filepath.Join(root, "source"), []byte("x"), 0600); err != nil {
		t.Fatal(err)
	}
	for _, count := range []bool{false, true} {
		s := migrationScanner{ctx: context.Background()}
		if count {
			s.i.Sources = make([]MigrationSource, maxMigrationItems)
		} else {
			s.bytes = 64 << 20
		}
		if _, err := s.file(root, "source", 1); err == nil {
			t.Fatal("inventory bound ignored")
		}
	}
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	if _, err := InspectLaunchMigration(ctx, root, "idea"); err == nil {
		t.Fatal("cancelled inventory accepted")
	}
}
