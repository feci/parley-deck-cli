package budget

import (
	"context"
	"encoding/json"
	"os"
	"path/filepath"
	"reflect"
	"sort"
	"strings"
	"testing"
)

// A root with no Git registration: launchScopeDeclared resolves symlinks, so a
// row's roots are the evaluated path and a comparison against the raw temporary
// directory would pass or fail for the wrong reason.
func runIdentityRoot(t *testing.T) string {
	t.Helper()
	root, err := filepath.EvalSymlinks(t.TempDir())
	if err != nil {
		t.Fatal(err)
	}
	return root
}

func writeRunFile(t *testing.T, root, name, file, body string) {
	t.Helper()
	dir := filepath.Join(root, "parley-deck", "runs", name)
	if err := os.MkdirAll(filepath.Dir(filepath.Join(dir, file)), 0o700); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(dir, file), []byte(body), 0o600); err != nil {
		t.Fatal(err)
	}
}

func rowsByPath(t *testing.T, r RunIdentityReport) map[string]RunIdentityRow {
	t.Helper()
	byPath := map[string]RunIdentityRow{}
	for _, row := range r.Rows {
		if _, ok := byPath[row.Path]; ok {
			t.Fatalf("a path was reported twice with one file set: %q", row.Path)
		}
		byPath[row.Path] = row
	}
	return byPath
}

// The load-bearing distinction of the discovery surface: an absent identity is
// an answer and is the only declarable class, while conflicting and unreadable
// identities are evidence of identity and must never be reported as absent or
// offered a declaration. The inspection reports all of them without refusing,
// which is precisely what the migration scanner cannot do.
func TestInspectRunIdentitiesSeparatesAbsentFromEveryOtherClass(t *testing.T) {
	ctx := context.Background()
	root := runIdentityRoot(t)
	named := func(idea string) string {
		return `{"time":"2026-05-10T19:40:03.126637Z","type":"run.created","data":{"idea":"` + idea + `"}}`
	}
	writeRunDirectory(t, root, "absent", unscopedRunEvent)
	writeRunDirectory(t, root, "conflicting", named("idea"), `{"time":"2026-05-10T19:40:04.126637Z","type":"run.phase","data":{"idea":"other"}}`)
	writeRunDirectory(t, root, "malformed", `{"time":"2026-05-10T19:40:03.126637Z"`)
	writeRunDirectory(t, root, "empty-identity", `{"time":"2026-05-10T19:40:03.126637Z","type":"run.created","data":{"idea":"   "}}`)
	writeRunDirectory(t, root, "bound", named("idea"))

	r, err := InspectRunIdentities(ctx, root, "idea", nil)
	if err != nil {
		t.Fatalf("discovery refused on the identities it exists to enumerate: %v", err)
	}
	rows := rowsByPath(t, r)
	if len(rows) != 5 {
		t.Fatalf("a visible run directory was dropped from the enumeration: %+v", r.Rows)
	}
	for _, c := range []struct{ name, observation string }{
		{"absent", RunIdentityAbsent},
		{"conflicting", RunIdentityConflicting},
		{"malformed", RunIdentityUnreadable},
		// An empty-string identity is a malformed claim of identity, not the
		// absence of one, so it is classified with the conflicting class and is
		// never declarable.
		{"empty-identity", RunIdentityConflicting},
		{"bound", RunIdentityBound},
	} {
		row, ok := rows[unscopedRunPrefix+c.name]
		if !ok {
			t.Fatalf("run %q was not enumerated: %+v", c.name, r.Rows)
		}
		if row.Observation != c.observation {
			t.Fatalf("run %q classified as %q, want %q (detail %q)", c.name, row.Observation, c.observation, row.Detail)
		}
		declarable := c.observation == RunIdentityAbsent
		if (row.Declaration != "") != declarable {
			t.Fatalf("run %q declarability is %q for observation %q", c.name, row.Declaration, row.Observation)
		}
		if declarable && row.Declaration != declareRun(t, root, c.name) {
			t.Fatalf("run %q offered a declaration the accounting boundary does not build: %q", c.name, row.Declaration)
		}
		// Every unbound row says unknown, never zero and never a count; only a
		// run that proved its own identity is reported as scoped, with a name.
		if c.observation == RunIdentityBound {
			if row.History != "scoped" || row.Idea != "idea" {
				t.Fatalf("a run that proved its identity lost it: %+v", row)
			}
		} else if row.History != UnknownHistory || row.Idea != "" {
			t.Fatalf("an unrecovered identity was reported as known: %+v", row)
		}
		if row.Copies != 1 || !reflect.DeepEqual(row.Roots, []string{root}) {
			t.Fatalf("run %q lost its actual single copy: %+v", c.name, row)
		}
	}
	// Nothing in the report may be read as an action count, and the report is a
	// pure function of the observation: the root the caller stood in is excluded
	// from the digest, so a second inspection of the same tree is identical.
	if !strings.Contains(r.Uncertainty[0], "no field of it is an action count") {
		t.Fatalf("the report dropped its own non-count disclaimer: %+v", r.Uncertainty)
	}
	if r.Coverage != "visible-roots" || r.Schema != runIdentitySchema {
		t.Fatalf("report coverage or schema changed: %+v", r)
	}
	again, err := InspectRunIdentities(ctx, root, "idea", nil)
	if err != nil {
		t.Fatal(err)
	}
	if again.ReportSHA256 != r.ReportSHA256 {
		t.Fatalf("one unchanged observation produced two digests: %q %q", r.ReportSHA256, again.ReportSHA256)
	}
	// The quoted record is the observation alone: the caller's own root is not a
	// property of what the repository holds, so it is not a field of the bytes a
	// later operator record would quote.
	record, err := r.OperatorRecord()
	if err != nil {
		t.Fatal(err)
	}
	var fields map[string]json.RawMessage
	if err := json.Unmarshal(record, &fields); err != nil {
		t.Fatal(err)
	}
	names := []string{}
	for name := range fields {
		names = append(names, name)
	}
	sort.Strings(names)
	if !reflect.DeepEqual(names, []string{"coverage", "rows", "schema", "scope"}) {
		t.Fatalf("the operator record changed shape: %v", names)
	}
}

// Copies of one file set group into one row; copies that disagree cannot be
// declared by content, so the report withholds the flag value rather than
// emitting one the accounting boundary would refuse.
func TestInspectRunIdentitiesGroupsCopiesAndWithholdsDivergentDeclarations(t *testing.T) {
	ctx := context.Background()
	root, linked, _ := worktreeRepoFixture(t)
	for _, origin := range []string{root, linked} {
		writeRunDirectory(t, origin, "20260510T194003Z", unscopedRunEvent)
	}
	r, err := InspectRunIdentities(ctx, root, "idea", nil)
	if err != nil {
		t.Fatal(err)
	}
	rows := rowsByPath(t, r)
	row := rows[unscopedRunPrefix+"20260510T194003Z"]
	if row.Copies != 2 || !reflect.DeepEqual(row.Roots, []string{linked, root}) {
		t.Fatalf("identical copies were not grouped onto their actual roots: %+v", row)
	}
	// Copies counts copies of one file set. It is never an action count, so one
	// declaration covers every copy.
	if row.Declaration != declareRun(t, root, "20260510T194003Z") {
		t.Fatalf("grouped copies did not yield the single covering declaration: %q", row.Declaration)
	}
	writeRunDirectory(t, linked, "20260510T194003Z", unscopedRunEvent+" ")
	divergent, err := InspectRunIdentities(ctx, root, "idea", nil)
	if err != nil {
		t.Fatalf("divergent copies refused instead of being reported: %v", err)
	}
	if len(divergent.Rows) != 2 {
		t.Fatalf("divergent copies were collapsed into one file set: %+v", divergent.Rows)
	}
	for _, split := range divergent.Rows {
		if split.Path != unscopedRunPrefix+"20260510T194003Z" || split.Copies != 1 {
			t.Fatalf("a divergent copy lost its own identity row: %+v", split)
		}
		if split.Declaration != "" {
			t.Fatalf("a declaration was offered for a path whose copies disagree: %+v", split)
		}
	}
	if !strings.Contains(strings.Join(divergent.Uncertainty, "\n"), "copies with different file sets") {
		t.Fatalf("divergence was not disclosed to the operator: %+v", divergent.Uncertainty)
	}
}

// The digest binds contents, not location: that is what makes one declaration
// cover every copy. Entries the walk cannot hash are refused rather than
// followed or skipped, because a declaration that binds a file set has to bind
// every byte of it.
func TestRunDirectoryManifestBindsContentNotLocation(t *testing.T) {
	first, second := runIdentityRoot(t), runIdentityRoot(t)
	for _, root := range []string{first, second} {
		writeRunDirectory(t, root, "run", unscopedRunEvent)
		writeRunFile(t, root, "run", filepath.Join("agents", "b.log"), "second\n")
		writeRunFile(t, root, "run", "a.json", "{}\n")
	}
	entries, err := RunDirectoryManifest(filepath.Join(first, "parley-deck", "runs", "run"))
	if err != nil {
		t.Fatal(err)
	}
	paths := []string{}
	for _, e := range entries {
		paths = append(paths, e.Path)
	}
	// Sorted, slash-relative and recursive: a nested file is bound under the
	// same rule as a top-level one, on every platform's separator.
	if !reflect.DeepEqual(paths, []string{"a.json", "agents/b.log", "events.jsonl"}) {
		t.Fatalf("manifest paths are not sorted slash-relative entries: %+v", paths)
	}
	if entries[1].Bytes != int64(len("second\n")) || entries[1].SHA256 != key("second\n") {
		t.Fatalf("a nested entry lost its actual bytes or digest: %+v", entries[1])
	}
	one, err := RunDirectoryManifestDigest(filepath.Join(first, "parley-deck", "runs", "run"))
	if err != nil {
		t.Fatal(err)
	}
	two, err := RunDirectoryManifestDigest(filepath.Join(second, "parley-deck", "runs", "run"))
	if err != nil {
		t.Fatal(err)
	}
	if one != two {
		t.Fatalf("the same file set digested differently in two locations: %q %q", one, two)
	}
	// A file the scanner never reads is still bound, so an unrelated addition
	// changes the digest. Deliberate and fail-closed.
	writeRunFile(t, second, "run", filepath.Join("agents", "c.log"), "")
	if changed, err := RunDirectoryManifestDigest(filepath.Join(second, "parley-deck", "runs", "run")); err != nil || changed == one {
		t.Fatalf("an added file did not change the bound file set: %q %v", changed, err)
	}
}

func TestRunDirectoryManifestRefusesUnhashableEntriesAndRoots(t *testing.T) {
	root := runIdentityRoot(t)
	runs := filepath.Join(root, "parley-deck", "runs")
	writeRunDirectory(t, root, "run", unscopedRunEvent)
	dir := filepath.Join(runs, "run")
	if _, err := RunDirectoryManifest(dir); err != nil {
		t.Fatalf("the plain fixture is not hashable: %v", err)
	}
	// A symlink is refused, never followed and never silently skipped: following
	// it would bind bytes outside the directory, skipping it would leave a named
	// entry out of a file set the declaration claims to bind whole.
	outside := filepath.Join(root, "outside.log")
	if err := os.WriteFile(outside, []byte("elsewhere\n"), 0o600); err != nil {
		t.Fatal(err)
	}
	for _, target := range []string{outside, runs} {
		link := filepath.Join(dir, "link")
		if err := os.Symlink(target, link); err != nil {
			t.Skipf("symlinks unavailable: %v", err)
		}
		if _, err := RunDirectoryManifest(dir); err == nil || !strings.Contains(err.Error(), "contains a symlink") {
			t.Fatalf("a symlink to %q was admitted into a bound file set: %v", target, err)
		}
		if err := os.Remove(link); err != nil {
			t.Fatal(err)
		}
	}
	// The root itself must be a real directory. A symlinked run directory is not
	// one: resolving it would let a declaration bind a file set held elsewhere.
	alias := filepath.Join(runs, "alias")
	if err := os.Symlink(dir, alias); err != nil {
		t.Skipf("symlinks unavailable: %v", err)
	}
	for _, bad := range []string{alias, filepath.Join(dir, "events.jsonl")} {
		if _, err := RunDirectoryManifest(bad); err == nil || !strings.Contains(err.Error(), "is not a real directory") {
			t.Fatalf("%q was hashed as a run directory: %v", bad, err)
		}
	}
	if _, err := RunDirectoryManifest(filepath.Join(runs, "absent")); !os.IsNotExist(err) {
		t.Fatalf("a missing run directory was not reported as missing: %v", err)
	}
}

// Discovery preserves cursor identity without overriding strict migration;
// for cursor-only unknown history it withholds an unusable declaration.
func TestRunIdentityDiscoveryPreservesIdentityAndWithholdsUnusableDeclaration(t *testing.T) {
	ctx := context.Background()

	// A run.created that carries no idea key, in a run whose cursor names one.
	// Discovery recovers the identity from the cursor and reports the run as
	// bound, so it offers no declaration; the scanner still refuses, because the
	// idea key is required on run.created whatever the cursor proves. The
	// operator is then left with a refusal and no flag value that resolves it.
	cursorRoot := runIdentityRoot(t)
	writeRunDirectory(t, cursorRoot, "cursor", unscopedRunEvent)
	writeRunFile(t, cursorRoot, "cursor", "driver.json", `{"idea":"idea","cycle":1}`)
	r, err := InspectRunIdentities(ctx, cursorRoot, "idea", nil)
	if err != nil {
		t.Fatalf("discovery refused a cursor-identified run: %v", err)
	}
	row := rowsByPath(t, r)[unscopedRunPrefix+"cursor"]
	if row.Observation != RunIdentityBound || row.Idea != "idea" || row.History != "scoped" || row.Declaration != "" {
		t.Fatalf("a cursor-identified run was not reported as bound and undeclarable: %+v", row)
	}
	if _, err := InspectProtocolMigration(ctx, cursorRoot, "idea", CrossReview, ""); err == nil || !strings.Contains(err.Error(), "historical run identity is missing or conflicting") {
		t.Fatalf("default nil-declaration scanning of a keyless run.created changed: %v", err)
	}

	// A cursor with no events remains unknown. Discovery withholds the
	// unusable declaration and explains why; manually supplying its manifest
	// must retain the original strict scanner refusal.
	bareRoot := runIdentityRoot(t)
	writeRunFile(t, bareRoot, "bare", "driver.json", `{"cycle":1}`)
	bare, err := InspectRunIdentities(ctx, bareRoot, "idea", nil)
	if err != nil {
		t.Fatalf("discovery refused a cursor-only run: %v", err)
	}
	bareRow := rowsByPath(t, bare)[unscopedRunPrefix+"bare"]
	if bareRow.Observation != RunIdentityAbsent || bareRow.Declaration != "" || bareRow.History != UnknownHistory {
		t.Fatalf("a cursor-only unknown run offered an unusable declaration: %+v", bareRow)
	}
	if !strings.Contains(strings.Join(bare.Uncertainty, "\n"), bareRow.Path) {
		t.Fatalf("withheld declaration lacks a run-specific explanation: %+v", bare.Uncertainty)
	}
	if _, err := InspectProtocolMigrationDeclarations(ctx, bareRoot, "idea", CrossReview, "", nil, []string{declareRun(t, bareRoot, "bare")}); err == nil || !strings.Contains(err.Error(), "historical driver cursor has no readable run identity") {
		t.Fatalf("declared cursor-only scanning changed: %v", err)
	}
}

// Ordinary callers pass no declaration and must see exactly what they saw
// before this surface existed: discovery is read-only and grants nothing, so
// running it changes no refusal and creates no accounting state.
func TestInspectRunIdentitiesGrantsNothingToOrdinaryCallers(t *testing.T) {
	ctx := context.Background()
	root := runIdentityRoot(t)
	writeRunDirectory(t, root, "20260510T194003Z", unscopedRunEvent)
	if _, err := InspectRunIdentities(ctx, root, "idea", nil); err != nil {
		t.Fatal(err)
	}
	// Checked here, before any other call: the discovery surface writes nothing,
	// so no accounting state exists yet on a repository it has inspected.
	for _, path := range []string{filepath.Join(root, ".git", "parley-launch-budgets"), filepath.Join(root, ".parley-runtime", "launch-budgets")} {
		if _, err := os.Stat(path); !os.IsNotExist(err) {
			t.Fatalf("a read-only inspection created accounting state at %s: %v", path, err)
		}
	}
	if _, err := InspectProtocolMigration(ctx, root, "idea", CrossReview, ""); err == nil || !strings.Contains(err.Error(), "historical run identity is missing or conflicting") {
		t.Fatalf("inspecting identities changed the undeclared refusal: %v", err)
	}
	if _, err := EnsureCycleBinding(ctx, root, "idea", CrossReview, 3, 0, "", ""); err == nil {
		t.Fatal("ordinary cycle bootstrap was unblocked by a read-only inspection")
	}
	// An empty repository is enumerated as empty without inventing a row, and
	// without claiming that absence is an observation about any run.
	empty, err := InspectRunIdentities(ctx, runIdentityRoot(t), "idea", nil)
	if err != nil {
		t.Fatal(err)
	}
	if len(empty.Rows) != 0 || len(empty.Uncertainty) != 1 {
		t.Fatalf("an empty repository produced rows or extra claims: %+v", empty)
	}
}
