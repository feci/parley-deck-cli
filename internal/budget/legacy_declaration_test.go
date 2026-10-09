package budget

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

const legacyRun = "parley-deck/runs/20260510T194003Z"

func legacyRequest(t *testing.T, root string) LegacyDeclarationRequest {
	t.Helper()
	p, err := InspectLegacyDeclaration(context.Background(), root, legacyRun)
	if err != nil {
		t.Fatal(err)
	}
	if p.History != "unknown-history" || p.ManifestSHA256 == "" {
		t.Fatalf("lost uncertainty: %+v", p)
	}
	return LegacyDeclarationRequest{Path: legacyRun, DecisionID: "fixture-only", Reason: "fixture explicit unknown history", ExpectedPreviewSHA256: p.PreviewSHA256, WritersStopped: true, AcknowledgeUnknown: true}
}

func TestLegacyDeclarationUnstallsSuccessiveIdeasWithoutResettingCaps(t *testing.T) {
	ctx := context.Background()
	root, linked, _ := worktreeRepoFixture(t)
	for _, r := range []string{root, linked} {
		writeRunDirectory(t, r, "20260510T194003Z", unscopedRunEvent)
	}
	before, _ := os.ReadFile(filepath.Join(root, legacyRun, "events.jsonl"))
	if _, err := EnsureCycleBinding(ctx, root, "first", CrossReview, 3, 2, "", ""); err == nil || !strings.Contains(err.Error(), "lacks idea identity") {
		t.Fatalf("undeclared history passed: %v", err)
	}
	dir, _, _, _ := legacyScope(ctx, root)
	if _, err := os.Stat(dir); !os.IsNotExist(err) {
		t.Fatal("driver authored declaration state")
	}
	q := legacyRequest(t, root)
	if _, err := os.Stat(dir); !os.IsNotExist(err) {
		t.Fatal("inspect mutated authority")
	}
	applied, err := ApplyLegacyDeclaration(ctx, root, q)
	if err != nil {
		t.Fatal(err)
	}
	again, err := ApplyLegacyDeclaration(ctx, linked, q)
	if err != nil || migrationDigest(again) != migrationDigest(applied) {
		t.Fatalf("replay not identical: %v", err)
	}
	conflict := q
	conflict.Reason = "different decision"
	if _, err := ApplyLegacyDeclaration(ctx, root, conflict); err == nil {
		t.Fatal("conflicting ID accepted")
	}
	for _, idea := range []string{"first", "second"} {
		b, err := EnsureCycleBinding(ctx, linked, idea, CrossReview, 3, 2, "", "")
		if err != nil {
			t.Fatal(err)
		}
		if n, err := b.Reserve(ctx, "third"); err != nil || n != 3 {
			t.Fatalf("count changed: %d %v", n, err)
		}
		if _, err := b.Reserve(ctx, "fourth"); !errors.Is(err, ErrLimit) {
			t.Fatalf("cap reset: %v", err)
		}
		st, err := InspectCycleBudget(ctx, root, idea, CrossReview)
		if err != nil || st.LegacyHistory == nil || st.Policy.LegacyHistorySHA256 == "" || st.Spent != 3 {
			t.Fatalf("missing frozen uncertainty: %+v %v", st, err)
		}
		raw, _ := json.Marshal(st.LegacyHistory)
		if !bytes.Contains(raw, []byte("unknown-history")) || !bytes.Contains(raw, []byte(q.DecisionID)) || !bytes.Contains(raw, []byte(applied.Preview.ManifestSHA256)) {
			t.Fatalf("missing provenance: %s", raw)
		}
	}
	after, _ := os.ReadFile(filepath.Join(root, legacyRun, "events.jsonl"))
	if !bytes.Equal(before, after) {
		t.Fatal("historical bytes rewritten")
	}
	// Request-scoped migration stays a distinct authority; it does not silently
	// consume the new durable declaration or change its old unknown semantics.
	if _, err := InspectProtocolMigration(ctx, root, "third", CrossReview, ""); err == nil {
		t.Fatal("new ledger changed existing per-idea migration semantics")
	}
}

func TestLegacyDeclarationHostileChangesRefuseNewAndCachedBindings(t *testing.T) {
	for _, mode := range []string{"mutated", "missing-copy", "all-deleted", "cursor", "run-manifest", "charge", "identity", "extra-unknown", "missing-record", "corrupt-record", "hidden-file", "empty-directory", "symlink", "hardlink", "missing-adoption", "changed-policy-reference", "unavailable-root"} {
		t.Run(mode, func(t *testing.T) {
			ctx := context.Background()
			root, linked, _ := worktreeRepoFixture(t)
			for _, r := range []string{root, linked} {
				writeRunDirectory(t, r, "20260510T194003Z", unscopedRunEvent)
			}
			q := legacyRequest(t, root)
			if _, err := ApplyLegacyDeclaration(ctx, root, q); err != nil {
				t.Fatal(err)
			}
			b, err := EnsureCycleBinding(ctx, root, "first", CrossReview, 3, 1, "", "")
			if err != nil {
				t.Fatal(err)
			}
			dir := filepath.Join(linked, legacyRun)
			events := filepath.Join(dir, "events.jsonl")
			ledger, _, _, _ := legacyScope(ctx, root)
			switch mode {
			case "mutated":
				err = os.WriteFile(events, []byte(unscopedRunEvent+"\n\n"), 0600)
			case "missing-copy":
				err = os.RemoveAll(dir)
			case "all-deleted":
				err = os.RemoveAll(dir)
				if err == nil {
					err = os.RemoveAll(filepath.Join(root, legacyRun))
				}
			case "cursor":
				err = os.WriteFile(filepath.Join(dir, "driver.json"), []byte(`{"rounds_run":2}`), 0600)
			case "run-manifest":
				err = os.WriteFile(filepath.Join(dir, "run.json"), []byte(`{}`), 0600)
			case "charge":
				err = os.WriteFile(events, []byte(unscopedRunEvent+"\n"+`{"time":"2026-05-10T19:40:03Z","type":"run.phase","data":{"action":"promoted"}}`+"\n"), 0600)
			case "identity":
				err = os.WriteFile(events, []byte(strings.Replace(unscopedRunEvent, `"mode"`, `"idea":"real","mode"`, 1)+"\n"), 0600)
			case "extra-unknown":
				writeRunDirectory(t, linked, "extra", unscopedRunEvent)
			case "missing-record":
				err = os.Remove(filepath.Join(ledger, "records", key(q.DecisionID)+".json"))
			case "corrupt-record":
				err = os.WriteFile(filepath.Join(ledger, "records", key(q.DecisionID)+".json"), []byte(`{}`), 0600)
			case "hidden-file":
				err = os.WriteFile(filepath.Join(dir, ".attempt"), []byte("charged"), 0600)
			case "empty-directory":
				err = os.Mkdir(filepath.Join(dir, "agents"), 0700)
			case "symlink":
				target := filepath.Join(t.TempDir(), "events")
				if err = os.Rename(events, target); err == nil {
					err = os.Symlink(target, events)
				}
				if err != nil {
					t.Skipf("symlink unavailable: %v", err)
				}
			case "hardlink":
				err = os.Link(events, filepath.Join(t.TempDir(), "alias"))
				if err != nil {
					t.Skipf("hardlink unavailable: %v", err)
				}
			case "missing-adoption":
				err = os.Remove(filepath.Join(filepath.Dir(b.Store.Dir), "legacy.json"))
			case "changed-policy-reference":
				p := b.Policy
				p.LegacyHistorySHA256 = ""
				raw, _ := json.Marshal(p)
				err = os.WriteFile(filepath.Join(filepath.Dir(b.Store.Dir), "policy.json"), raw, 0600)
			case "unavailable-root":
				err = os.RemoveAll(linked)
			}
			if err != nil {
				t.Fatal(err)
			}
			if _, err := b.Reserve(ctx, "next"); err == nil {
				t.Fatal("cached binding admitted hostile mutation")
			}
			st, err := b.Store.Inspect(ctx)
			if err != nil || b.Count(st) != 1 {
				t.Fatalf("refusal changed charges: %v %+v", err, st)
			}
			if mode != "missing-adoption" && mode != "changed-policy-reference" {
				if _, err := EnsureCycleBinding(ctx, root, "new", CrossReview, 3, 0, "", ""); err == nil {
					t.Fatal("new idea admitted hostile history")
				}
			}
		})
	}
}

func TestLegacyDeclarationClosedEligibility(t *testing.T) {
	for _, body := range []string{"", `{}`, `{"time":"2026-05-10T19:40:03Z","type":"agent.started","data":{}}`, `{"time":"2026-05-10T19:40:03Z","type":"run.phase","data":{"action":"fixup"}}`, `{"time":"2026-05-10T19:40:03Z","type":"run.created","data":{"rounds_run":"0"}}`, `{"time":"2026-05-10T19:40:03Z","type":"run.created","data":{"idea":null}}`, `{"time":"2026-05-10T19:40:03Z","type":"run.created","data":{"mode":"auto","mode":"manual"}}`} {
		root := t.TempDir()
		writeRunDirectory(t, root, "20260510T194003Z", body)
		if _, err := InspectLegacyDeclaration(context.Background(), root, legacyRun); err == nil {
			t.Fatalf("ineligible history accepted: %s", body)
		}
	}
	root := t.TempDir()
	writeRunDirectory(t, root, "20260510T194003Z", unscopedRunEvent)
	q := legacyRequest(t, root)
	q.ExpectedPreviewSHA256 = strings.Repeat("a", 64)
	if _, err := ApplyLegacyDeclaration(context.Background(), root, q); err == nil {
		t.Fatal("stale preview accepted")
	}
	for _, path := range []string{"../runs/old", "parley-deck/runs/../old", legacyRun + "/", filepath.Join(root, legacyRun)} {
		if _, err := InspectLegacyDeclaration(context.Background(), root, path); err == nil {
			t.Fatalf("bad path accepted: %s", path)
		}
	}
}
