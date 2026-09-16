package budget

import (
	"context"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

// The MINOR-1 round-trip property, behaviorally pinned: for every run the
// tolerant report classifies as identity-absent, the offered remedy must agree
// with the strict declared scan. Either the report offers a Declaration and
// InspectProtocolMigrationDeclarations accepts exactly that value, or the
// report withholds it, explains why in Uncertainty, and the strict scan keeps
// its original refusal verbatim. Each fixture below is one divergence the
// tolerant read used to emit a guaranteed-refused flag value for; the shared
// helpers (runIdentityRoot, writeRunDirectory, writeRunFile, declareRun,
// rowsByPath, unscopedRunEvent) live in the package's existing test files.
func TestRunDeclarationEligibilityWithholdsWhatTheStrictScanRefuses(t *testing.T) {
	ctx := context.Background()
	future := time.Now().UTC().Add(48 * time.Hour).Format(time.RFC3339Nano)
	oversized := `{"cycle":1,"padding":"` + strings.Repeat("x", 1<<20) + `"}`
	for _, c := range []struct {
		name      string
		events    []string // events.jsonl lines; nil means the file does not exist
		cursor    string   // driver.json body, written only when hasCursor
		hasCursor bool
		refusal   string // the strict declared scan's original refusal, preserved verbatim
		reason    string // the explanation the report records in Uncertainty
	}{
		// (a) A cursor without events: the strict scan refuses it before the
		// declared branch is ever reached, so the row can never be recorded.
		{"cursor-only", nil, `{"cycle":1}`, true,
			"historical driver cursor has no readable run identity",
			"its events.jsonl is missing while its driver.json exists"},
		// (a') Both files missing: the scan skips the directory and the
		// declaration dies as stale — present in no visible root.
		{"empty-run-directory", nil, "", false,
			"not present in any visible root",
			"retains no row for it"},
		// (b) An event line without its required time field.
		{"missing-event-time", []string{`{"type":"run.phase","data":{"action":"fixup"}}`}, "", false,
			"missing required historical field time",
			"an events.jsonl line is not valid strict historical JSON"},
		// (b') An event stamped in the future.
		{"future-event-time", []string{`{"time":"` + future + `","type":"run.phase","data":{"action":"fixup"}}`}, "", false,
			"historical accounting timestamp is missing or in the future",
			"an events.jsonl event time is missing or in the future"},
		// (c) Duplicate keys: the tolerant read takes last-wins, the strict
		// scan refuses the bytes.
		{"duplicate-key", []string{`{"time":"2026-05-10T19:40:03.126637Z","type":"run.created","type":"run.phase","data":{}}`}, "", false,
			"malformed or duplicate historical fields",
			"an events.jsonl line is not valid strict historical JSON"},
		// (c') Non-lowercase keys: encoding/json matches fields
		// case-insensitively, the strict scan does not.
		{"uppercase-key", []string{`{"Time":"2026-05-10T19:40:03.126637Z","Type":"run.created","Data":{}}`}, "", false,
			"malformed or duplicate historical fields",
			"an events.jsonl line is not valid strict historical JSON"},
		// (d) A zero-byte crash-artifact cursor: the tolerant read treats it as
		// no cursor, the strict scan fails to parse it.
		{"empty-cursor", []string{unscopedRunEvent}, "", true,
			"malformed or duplicate historical fields",
			"its driver.json is not valid strict historical JSON"},
		// (e) A cursor inside the report's 16MiB read but past the scan's 1MiB
		// cursor bound.
		{"oversized-cursor", []string{unscopedRunEvent}, oversized, true,
			"driver history is not a bounded regular file",
			"1 MiB cursor bound"},
	} {
		root := runIdentityRoot(t)
		if c.events == nil {
			if err := os.MkdirAll(filepath.Join(root, "parley-deck", "runs", c.name), 0o700); err != nil {
				t.Fatal(err)
			}
		} else {
			writeRunDirectory(t, root, c.name, c.events...)
		}
		if c.hasCursor {
			writeRunFile(t, root, c.name, "driver.json", c.cursor)
		}
		// The report stays tolerant: it enumerates and classifies the run
		// without refusing. The row and its unknown history are preserved.
		r, err := InspectRunIdentities(ctx, root, "idea", nil)
		if err != nil {
			t.Fatalf("%s: the tolerant report refused the history it exists to enumerate: %v", c.name, err)
		}
		row, ok := rowsByPath(t, r)[unscopedRunPrefix+c.name]
		if !ok {
			t.Fatalf("%s: the run was dropped from the enumeration: %+v", c.name, r.Rows)
		}
		if row.Observation != RunIdentityAbsent || row.History != UnknownHistory || row.Idea != "" {
			t.Fatalf("%s: classification changed instead of only the offer: %+v", c.name, row)
		}
		if row.Declaration != "" {
			t.Fatalf("%s: a declaration the strict scan deterministically refuses was offered: %q", c.name, row.Declaration)
		}
		uncertainty := strings.Join(r.Uncertainty, "\n")
		if !strings.Contains(uncertainty, row.Path) || !strings.Contains(uncertainty, c.reason) {
			t.Fatalf("%s: the withheld declaration was not explained to the operator: %+v", c.name, r.Uncertainty)
		}
		// The strict scan's original refusal is preserved, verbatim, against
		// the exact value the report would previously have printed.
		declaration := declareRun(t, root, c.name)
		if _, err := InspectProtocolMigrationDeclarations(ctx, root, "idea", CrossReview, "", nil, []string{declaration}); err == nil || !strings.Contains(err.Error(), c.refusal) {
			t.Fatalf("%s: the strict declared scan's original refusal changed: %v", c.name, err)
		}
	}
}

// The valid control is the known 117-byte May 10 run.created event this remedy
// was scoped against: it must remain declarable, and the exact value the
// report prints must be the value the strict declared scan accepts — retaining
// the row as unknown history, on its actual root, counting nothing.
func TestRunDeclarationEligibilityAdmitsTheValidControl(t *testing.T) {
	if len(unscopedRunEvent+"\n") != 117 {
		t.Fatalf("the known-valid fixture this remedy was scoped against changed: %d bytes", len(unscopedRunEvent+"\n"))
	}
	ctx := context.Background()
	root := runIdentityRoot(t)
	writeRunDirectory(t, root, "20260510T194003Z", unscopedRunEvent)
	r, err := InspectRunIdentities(ctx, root, "idea", nil)
	if err != nil {
		t.Fatal(err)
	}
	row := rowsByPath(t, r)[unscopedRunPrefix+"20260510T194003Z"]
	if row.Observation != RunIdentityAbsent || row.Declaration != declareRun(t, root, "20260510T194003Z") {
		t.Fatalf("the valid control lost its exact declarable offer: %+v", row)
	}
	if strings.Contains(strings.Join(r.Uncertainty, "\n"), row.Path) {
		t.Fatalf("a withholding was recorded against the valid control: %+v", r.Uncertainty)
	}
	i, err := InspectProtocolMigrationDeclarations(ctx, root, "idea", CrossReview, "", nil, []string{row.Declaration})
	if err != nil {
		t.Fatalf("the strict declared scan refused the offered value: %v", err)
	}
	if len(i.History.UnscopedRuns) != 1 {
		t.Fatalf("the declared row was dropped or duplicated: %+v", i.History.UnscopedRuns)
	}
	declared := i.History.UnscopedRuns[0]
	if declared.Path+"="+declared.SHA256 != row.Declaration || declared.History != UnknownHistory || declared.Copies != 1 || len(declared.Roots) != 1 || declared.Roots[0] != root {
		t.Fatalf("the accepted row lost its binding, unknown history or actual copy: %+v", declared)
	}
	if i.History.HistoryCoverage != DeclaredIncomplete || i.LowerBound != 0 {
		t.Fatalf("the accepted declaration miscounted or lost its coverage statement: %q %d", i.History.HistoryCoverage, i.LowerBound)
	}
}

// Eligibility is a property of one run's own file set under one root, never of
// a whole-root inspection: an un-declarable run on a root does not hide,
// poison or rename the valid declarable row beside it, and the explanatory
// Uncertainty names only the run it applies to.
func TestRunDeclarationEligibilityJudgesEachRunOnItsOwnFiles(t *testing.T) {
	ctx := context.Background()
	root := runIdentityRoot(t)
	writeRunDirectory(t, root, "20260510T194003Z", unscopedRunEvent)
	writeRunFile(t, root, "cursor-only", "driver.json", `{"cycle":1}`)
	r, err := InspectRunIdentities(ctx, root, "idea", nil)
	if err != nil {
		t.Fatal(err)
	}
	rows := rowsByPath(t, r)
	valid := rows[unscopedRunPrefix+"20260510T194003Z"]
	withheld := rows[unscopedRunPrefix+"cursor-only"]
	if valid.Declaration != declareRun(t, root, "20260510T194003Z") {
		t.Fatalf("a valid row was hidden by an unrelated un-declarable neighbour: %+v", valid)
	}
	if withheld.Observation != RunIdentityAbsent || withheld.Declaration != "" {
		t.Fatalf("the un-declarable run lost its classification or kept its offer: %+v", withheld)
	}
	uncertainty := strings.Join(r.Uncertainty, "\n")
	if !strings.Contains(uncertainty, withheld.Path) {
		t.Fatalf("the withholding was not recorded against its own run: %+v", r.Uncertainty)
	}
	if strings.Contains(uncertainty, valid.Path) {
		t.Fatalf("the withholding named an unrelated valid row: %+v", r.Uncertainty)
	}
}
