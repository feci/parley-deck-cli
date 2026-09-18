package budget

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"
)

// Identity classification for one historical run directory. Only an absent
// identity is declarable: a conflicting, malformed or unreadable identity is
// evidence of identity, so declaring it away would discard evidence instead of
// admitting an unknown. None of these values is a count and none asserts zero.
const (
	RunIdentityBound       = "bound"
	RunIdentityAbsent      = "run-identity-absent"
	RunIdentityConflicting = "run-identity-conflicting"
	RunIdentityUnreadable  = "run-identity-unreadable"
)

const (
	runIdentitySchema   = 1
	maxRunManifestFiles = 4096
	maxRunManifestBytes = 64 << 20
	maxRunManifestFile  = 16 << 20
)

// RunManifestEntry is one regular file inside a run directory, named by its
// slash-relative path within that directory, so an identical file set digests
// identically in every copy whatever the copy's absolute location is.
type RunManifestEntry struct {
	Path   string `json:"path"`
	SHA256 string `json:"sha256"`
	Bytes  int64  `json:"bytes"`
}

// RunDirectoryManifest lists every regular file under one run directory, sorted
// by its relative path. It is read-only and bounded. Symlinks and other
// non-regular entries are refused rather than followed or skipped: a
// declaration that binds a file set has to bind every byte of it, and an entry
// this walk cannot hash is not an entry it may ignore.
func RunDirectoryManifest(dir string) ([]RunManifestEntry, error) {
	info, err := os.Lstat(dir)
	if err != nil {
		return nil, err
	}
	if !info.IsDir() {
		return nil, fmt.Errorf("historical run path is not a real directory: %s", dir)
	}
	entries := []RunManifestEntry{}
	total := int64(0)
	err = filepath.WalkDir(dir, func(path string, e os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if e.Type()&os.ModeSymlink != 0 {
			return fmt.Errorf("historical run directory contains a symlink: %s", path)
		}
		if e.IsDir() {
			return nil
		}
		if !e.Type().IsRegular() {
			return fmt.Errorf("historical run directory contains a non-regular entry: %s", path)
		}
		relative, err := filepath.Rel(dir, path)
		if err != nil {
			return err
		}
		data, err := readStepHistoryFile(path, maxRunManifestFile)
		if err != nil {
			return err
		}
		total += int64(len(data))
		if len(entries) >= maxRunManifestFiles || total > maxRunManifestBytes {
			return errors.New("historical run directory exceeds the bounded manifest")
		}
		entries = append(entries, RunManifestEntry{filepath.ToSlash(relative), key(string(data)), int64(len(data))})
		return nil
	})
	if err != nil {
		return nil, err
	}
	sort.Slice(entries, func(a, b int) bool { return entries[a].Path < entries[b].Path })
	return entries, nil
}

// RunDirectoryManifestDigest is the value a declaration binds: the exact
// recursive file set of one run directory. It is a property of the contents
// alone, so every visible copy of the same run must produce the same digest,
// and an added, deleted or changed file in any copy changes it.
func RunDirectoryManifestDigest(dir string) (string, error) {
	entries, err := RunDirectoryManifest(dir)
	if err != nil {
		return "", err
	}
	return migrationDigest(entries), nil
}

// RunIdentityRow is one distinct (relative run directory, manifest digest) pair
// and the roots that hold exactly that file set. Copies counts copies of one
// file set; it is never an action count. Declaration is the verbatim
// --declare-unscoped-run value and is present only for a declarable row: an
// absent identity that the strict declared scan would also admit, re-validated
// per copy. Where the tolerant read and the strict scan disagree, Declaration
// is withheld and the disagreement is recorded in the report's Uncertainty.
type RunIdentityRow struct {
	Path        string   `json:"path"`
	SHA256      string   `json:"sha256"`
	Observation string   `json:"observation"`
	Idea        string   `json:"idea,omitempty"`
	History     string   `json:"history"`
	Copies      int      `json:"copies"`
	Roots       []string `json:"roots"`
	Declaration string   `json:"declaration,omitempty"`
	Detail      string   `json:"detail,omitempty"`
}

// RunIdentityReport is a read-only classification of historical run identities.
// It grants nothing: no apply, attest or confirm path in this package consumes
// it. It contains no floor, no total and no action count of any kind, and it
// never reports an unclassifiable run as absent history.
type RunIdentityReport struct {
	Schema       int              `json:"schema"`
	Root         string           `json:"root"`
	Idea         string           `json:"idea"`
	Scope        string           `json:"scope"`
	Coverage     string           `json:"coverage"`
	Roots        []string         `json:"roots"`
	Rows         []RunIdentityRow `json:"rows"`
	Uncertainty  []string         `json:"uncertainty"`
	ReportSHA256 string           `json:"report_sha256"`
}

// observedRunIdentities is the digest input. Root is excluded because it
// describes where the caller stood, not what the repository holds.
type observedRunIdentities struct {
	Schema   int              `json:"schema"`
	Scope    string           `json:"scope"`
	Coverage string           `json:"coverage"`
	Rows     []RunIdentityRow `json:"rows"`
}

// OperatorRecord renders the exact bytes a later operator record would quote.
// It is a pure function of this observation and grants nothing.
func (r RunIdentityReport) OperatorRecord() ([]byte, error) {
	return json.Marshal(observedRunIdentities{runIdentitySchema, r.Scope, r.Coverage, r.Rows})
}

// InspectRunIdentities enumerates every visible historical run directory and
// classifies the idea identity each one can prove about itself. It exists
// because the migration scanner refuses on the first unrecoverable identity, so
// without it no read-only surface can tell an operator what there is to declare.
//
// It is deliberately tolerant of an absent identity — that is the fact it
// reports — and deliberately strict about everything else: a broken repository,
// an unreadable root or a non-regular entry still refuses. It computes no
// floor, no total and no count, and it writes nothing.
//
// Tolerance is for enumeration only. A declaration is offered for an absent
// identity solely when the strict declared scan would admit that exact value,
// re-validated per copy under the scan's own bounds and validation
// (runDeclarationEligibility). Where the two reads disagree — a missing
// events.jsonl, an unusable event time, a strict-JSON failure, an empty or
// oversized cursor — the row is kept, the offer is withheld and the reason is
// recorded in Uncertainty, so a value printed here is never one the apply path
// deterministically refuses.
func InspectRunIdentities(ctx context.Context, root, idea string, declaredUnavailable []string) (RunIdentityReport, error) {
	_, scope, roots, unavailable, err := launchScopeDeclared(ctx, root, idea, true, declaredUnavailable)
	if err != nil {
		return RunIdentityReport{}, err
	}
	sort.Strings(roots)
	r := RunIdentityReport{Schema: runIdentitySchema, Root: root, Idea: idea, Scope: scope, Coverage: "visible-roots", Roots: roots, Rows: []RunIdentityRow{}, Uncertainty: []string{
		"this report classifies run identity only; no field of it is an action count, a floor or a total, and no row asserts that a run performed zero actions",
	}}
	if len(unavailable) > 0 {
		r.Coverage = DeclaredIncomplete
		for _, row := range unavailable {
			r.Uncertainty = append(r.Uncertainty, fmt.Sprintf("registration %q was declared unavailable (%s); its runs are not enumerated here and its history stays unknown, not absent", row.Path, row.Observation))
		}
	}
	rows := map[string]*RunIdentityRow{}
	digests := map[string]map[string]bool{}
	for _, origin := range roots {
		if err := ctx.Err(); err != nil {
			return RunIdentityReport{}, err
		}
		base := filepath.Join(origin, "parley-deck", "runs")
		entries, err := stepHistoryDirs(base)
		if err != nil {
			return RunIdentityReport{}, err
		}
		for _, entry := range entries {
			dir := filepath.Join(base, entry.Name())
			digest, err := RunDirectoryManifestDigest(dir)
			if err != nil {
				return RunIdentityReport{}, err
			}
			path := unscopedRunPrefix + entry.Name()
			if digests[path] == nil {
				digests[path] = map[string]bool{}
			}
			digests[path][digest] = true
			id := path + "=" + digest
			row := rows[id]
			if row == nil {
				if len(rows) >= maxMigrationItems {
					return RunIdentityReport{}, errors.New("too many historical run directories")
				}
				observation, name, detail := classifyRunIdentity(dir)
				row = &RunIdentityRow{Path: path, SHA256: digest, Observation: observation, Idea: name, History: UnknownHistory, Detail: detail}
				if observation == RunIdentityBound {
					row.History = "scoped"
				}
				if observation == RunIdentityAbsent {
					row.Declaration = id
				}
				rows[id] = row
			}
			// An absent identity is this report's tolerant answer; a declaration
			// must also survive the strict declared scan. Re-validate the offer
			// against every copy under the scan's exact bounds and validation
			// and withdraw it where the two reads disagree, recording the reason
			// — an offered value the apply path deterministically refuses is a
			// false remedy, not a remedy. The row and its unknown history stay.
			if row.Declaration != "" {
				if reason, eligible := runDeclarationEligibility(origin, entry.Name()); !eligible {
					row.Declaration = ""
					r.Uncertainty = append(r.Uncertainty, fmt.Sprintf("run %q has no recoverable idea identity, but the strict declared scan would refuse its declaration: %s; reconcile the underlying files and re-inspect rather than declare", row.Path, reason))
				}
			}
			row.Roots = append(row.Roots, origin)
			row.Copies = len(row.Roots)
		}
	}
	for _, row := range rows {
		// A path whose copies disagree cannot be declared by content: the
		// declaration would refuse. Say so here rather than emit a value that
		// cannot be used.
		if len(digests[row.Path]) > 1 {
			row.Declaration = ""
			r.Uncertainty = append(r.Uncertainty, fmt.Sprintf("run %q has copies with different file sets; reconcile them before any declaration, which binds one exact file set", row.Path))
		}
		sort.Strings(row.Roots)
		r.Rows = append(r.Rows, *row)
	}
	sort.Slice(r.Rows, func(a, b int) bool {
		x, y := r.Rows[a], r.Rows[b]
		if x.Path == y.Path {
			return x.SHA256 < y.SHA256
		}
		return x.Path < y.Path
	})
	sort.Strings(r.Uncertainty[1:])
	raw, err := r.OperatorRecord()
	if err != nil {
		return RunIdentityReport{}, err
	}
	r.ReportSHA256 = key(string(raw))
	return r, nil
}

// classifyRunIdentity recovers the idea identity one historical run can prove
// about itself. An absent identity is an answer, not an error — enumerating
// absent identities is this function's purpose. Malformed, unreadable and
// conflicting inputs are classified as themselves and never read as absent.
func classifyRunIdentity(dir string) (observation, idea, detail string) {
	identity, present := "", false
	bind := func(raw json.RawMessage) (string, bool) {
		if raw == nil {
			return "", true
		}
		var name string
		if json.Unmarshal(raw, &name) != nil || strings.TrimSpace(name) == "" {
			return "run reports an empty or non-string idea identity", false
		}
		if identity != "" && identity != name {
			return "run reports two different idea identities: " + identity + " and " + name, false
		}
		identity, present = name, true
		return "", true
	}
	events, eventsErr := readStepHistoryFile(filepath.Join(dir, "events.jsonl"), maxRunManifestFile)
	if eventsErr != nil && !os.IsNotExist(eventsErr) {
		return RunIdentityUnreadable, "", eventsErr.Error()
	}
	cursor, cursorErr := readStepHistoryFile(filepath.Join(dir, "driver.json"), maxRunManifestFile)
	if cursorErr != nil && !os.IsNotExist(cursorErr) {
		return RunIdentityUnreadable, "", cursorErr.Error()
	}
	if len(cursor) > 0 {
		var fields map[string]json.RawMessage
		if err := migrationJSON(cursor, &fields); err != nil {
			return RunIdentityUnreadable, "", "driver.json is malformed: " + err.Error()
		}
		for _, name := range []string{"idea", "idea_slug"} {
			if reason, ok := bind(fields[name]); !ok {
				return RunIdentityConflicting, "", reason
			}
		}
	}
	scan := bufio.NewScanner(bytes.NewReader(events))
	scan.Buffer(make([]byte, 4096), 1<<20)
	lines := 0
	for scan.Scan() {
		line := bytes.TrimSpace(scan.Bytes())
		if len(line) == 0 {
			continue
		}
		lines++
		if lines > maxMigrationItems {
			return RunIdentityUnreadable, "", "too many historical run events"
		}
		var e struct {
			Type string                     `json:"type"`
			Data map[string]json.RawMessage `json:"data,omitempty"`
		}
		if err := json.Unmarshal(line, &e); err != nil {
			return RunIdentityUnreadable, "", "events.jsonl is malformed: " + err.Error()
		}
		for _, name := range []string{"idea", "idea_slug"} {
			if reason, ok := bind(e.Data[name]); !ok {
				return RunIdentityConflicting, "", reason
			}
		}
	}
	if err := scan.Err(); err != nil {
		return RunIdentityUnreadable, "", err.Error()
	}
	if present {
		return RunIdentityBound, identity, ""
	}
	if len(events) == 0 && len(cursor) == 0 {
		return RunIdentityAbsent, "", "run directory carries no readable events and no cursor"
	}
	return RunIdentityAbsent, "", "no event and no cursor field carries an idea identity"
}

// runDeclarationEligibility answers the second question a declaration must
// survive: classifyRunIdentity asks, tolerantly, whether a run carries an idea
// identity; this asks, strictly, whether the declared scan's own machinery
// would admit a declaration for this exact run. It mirrors, for one directory
// and nothing else, the declared branch of the migration scanner's run walk —
// the same parent-chain and bounded-regular-file checks, the same per-file
// bounds (events 16MiB, cursor 1MiB), the same migrationJSON validation
// (duplicate or non-lowercase fields, required-field shape, no trailing data),
// the same untyped-event and identity-bind refusals, and the same
// missing-or-future event-time refusal. Nothing here duplicates the tolerant
// read, widens what is declarable, or relaxes the scanner: an eligible answer
// means the strict scan already accepts this run declared, and any divergence
// is returned as the reason the offer must be withheld. Eligibility is a
// property of this run's own files under this one root, so one un-declarable
// run never hides another row's valid offer.
func runDeclarationEligibility(root, name string) (string, bool) {
	relative := unscopedRunPrefix + name
	read := func(file string, limit int64) ([]byte, error) {
		// The same aliased-parent and bounded-read checks the migration scanner
		// applies to every history file, without its inventory bookkeeping.
		parent := root
		parts := strings.Split(relative+"/"+file, "/")
		for _, part := range parts[:len(parts)-1] {
			parent = filepath.Join(parent, part)
			st, err := os.Lstat(parent)
			if err != nil {
				return nil, err
			}
			if !st.IsDir() {
				return nil, errors.New("migration history contains an aliased directory")
			}
		}
		return readStepHistoryFile(filepath.Join(root, filepath.FromSlash(relative+"/"+file)), limit)
	}
	events, err := read("events.jsonl", 16<<20)
	if err != nil {
		if os.IsNotExist(err) {
			if _, e := os.Lstat(filepath.Join(root, filepath.FromSlash(relative), "driver.json")); e == nil || !os.IsNotExist(e) {
				return "its events.jsonl is missing while its driver.json exists, and the strict declared scan refuses a cursor with no readable run identity", false
			}
			return "its events.jsonl is missing, so the strict declared scan retains no row for it and refuses the declaration as not present in any visible root", false
		}
		return "its events.jsonl cannot be read under the strict declared scan's bounds: " + err.Error(), false
	}
	identity := ""
	bind := func(raw json.RawMessage) (string, bool) {
		if raw == nil {
			return "", true
		}
		var idea string
		if json.Unmarshal(raw, &idea) != nil || strings.TrimSpace(idea) == "" || identity != "" && identity != idea {
			return "a malformed or conflicting idea identity", false
		}
		identity = idea
		return "", true
	}
	cursor, err := read("driver.json", 1<<20)
	if err == nil {
		var fields map[string]json.RawMessage
		if e := migrationJSON(cursor, &fields); e != nil {
			return "its driver.json is not valid strict historical JSON: " + e.Error(), false
		}
		for _, key := range []string{"idea", "idea_slug"} {
			if what, ok := bind(fields[key]); !ok {
				return "its driver.json carries " + what, false
			}
		}
	} else if !os.IsNotExist(err) {
		return "its driver.json cannot be read under the strict declared scan's 1 MiB cursor bound: " + err.Error(), false
	}
	type event struct {
		Time time.Time                  `json:"time"`
		Type string                     `json:"type"`
		Data map[string]json.RawMessage `json:"data,omitempty"`
	}
	history := []event{}
	scan := bufio.NewScanner(bytes.NewReader(events))
	scan.Buffer(make([]byte, 4096), 1<<20)
	for scan.Scan() {
		line := bytes.TrimSpace(scan.Bytes())
		if len(line) == 0 {
			continue
		}
		var e event
		if err := migrationJSON(line, &e); err != nil {
			return "an events.jsonl line is not valid strict historical JSON: " + err.Error(), false
		}
		if e.Type == "" {
			return "an events.jsonl event carries no type", false
		}
		for _, key := range []string{"idea", "idea_slug"} {
			if what, ok := bind(e.Data[key]); !ok {
				return "an events.jsonl event carries " + what, false
			}
		}
		history = append(history, e)
		if len(history) > maxMigrationItems {
			return "its events.jsonl exceeds the strict declared scan's event bound", false
		}
	}
	if err := scan.Err(); err != nil {
		return "its events.jsonl cannot be scanned under the strict declared scan's line bound: " + err.Error(), false
	}
	// A declaration admits an unknown; it never overrides a recovered identity.
	if identity != "" {
		return "the strict declared scan recovers an idea identity from it: " + identity, false
	}
	for _, e := range history {
		if e.Time.IsZero() || e.Time.After(time.Now().UTC()) {
			return "an events.jsonl event time is missing or in the future", false
		}
	}
	return "", true
}
