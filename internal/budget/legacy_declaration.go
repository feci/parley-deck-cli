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
	"reflect"
	"sort"
	"strings"
	"time"
)

const legacyUnknownHistory = "unknown-history"
const maxLegacyDecisions = 256

// LegacyPreview is an observation, never a count or permission. The exact
// preview includes the visible copies and the existing append-only authority.
type LegacyPreview struct {
	Version        int      `json:"version"`
	Scope          string   `json:"scope"`
	Path           string   `json:"path"`
	ManifestSHA256 string   `json:"manifest_sha256"`
	History        string   `json:"history"`
	Roots          []string `json:"roots"`
	Decisions      []string `json:"decisions"`
	PreviewSHA256  string   `json:"preview_sha256"`
}

type LegacyDeclarationRequest struct {
	Path                  string `json:"path"`
	ExpectedPreviewSHA256 string `json:"expected_preview_sha256"`
	DecisionID            string `json:"decision_id"`
	Reason                string `json:"reason"`
	WritersStopped        bool   `json:"writers_stopped"`
	AcknowledgeUnknown    bool   `json:"acknowledge_unknown"`
}

// LegacyDeclaration records scope uncertainty, not a historical zero. Creation
// belongs only to the attended CLI control; this API cannot authenticate a human.
type LegacyDeclaration struct {
	Version    int                      `json:"version"`
	Scope      string                   `json:"scope"`
	RecordedAt time.Time                `json:"recorded_at"`
	Authority  string                   `json:"authority"`
	Request    LegacyDeclarationRequest `json:"request"`
	Preview    LegacyPreview            `json:"preview"`
}

func legacyScope(ctx context.Context, root string) (string, string, []string, error) {
	dir, scope, roots, err := launchScope(ctx, root, "", true)
	scope = "parley-legacy/v1:" + key(scope)
	dir = filepath.Join(filepath.Dir(dir), "legacy-"+key(scope))
	if err == nil {
		for path := dir; path != filepath.Dir(path); path = filepath.Dir(path) {
			st, e := os.Lstat(path)
			if os.IsNotExist(e) {
				continue
			}
			if e != nil {
				return "", "", nil, e
			}
			if !st.IsDir() {
				return "", "", nil, errors.New("legacy storage contains an aliased directory")
			}
		}
	}
	return dir, scope, roots, err
}

func legacyPreviewDigest(p LegacyPreview) string {
	p.PreviewSHA256 = ""
	return migrationDigest(p)
}

func validateLegacyRecord(r LegacyDeclaration, scope string) error {
	p, q := r.Preview, r.Request
	if r.Version != 1 || r.Scope != scope || r.Authority != "attended-operator" || r.RecordedAt.IsZero() || r.RecordedAt.After(time.Now().UTC()) || r.RecordedAt.Location() != time.UTC || !q.WritersStopped || !q.AcknowledgeUnknown || !validCycleDecision(q.DecisionID, q.Reason, q.ExpectedPreviewSHA256) || !validUnscopedRunPath(q.Path) || p.Version != 1 || p.Scope != scope || p.Path != q.Path || p.History != legacyUnknownHistory || p.PreviewSHA256 != q.ExpectedPreviewSHA256 || p.PreviewSHA256 != legacyPreviewDigest(p) || !validCycleDecision("manifest", "manifest", p.ManifestSHA256) || len(p.Roots) == 0 || len(p.Roots) > maxWorktreeRegistration || len(p.Decisions) > maxLegacyDecisions {
		return errors.New("invalid durable legacy declaration")
	}
	for i, root := range p.Roots {
		if !filepath.IsAbs(root) || filepath.Clean(root) != root || (i > 0 && p.Roots[i-1] >= root) {
			return errors.New("invalid legacy copy inventory")
		}
	}
	for i, digest := range p.Decisions {
		if !validCycleDecision("record", "record", digest) || (i > 0 && p.Decisions[i-1] >= digest) {
			return errors.New("invalid legacy decision inventory")
		}
	}
	return nil
}

// readLegacyRecords never initializes the store. Unknown files and orphaned
// partial publications are refused instead of silently dropping authority.
func readLegacyRecords(dir, scope string) ([]LegacyDeclaration, error) {
	if err := migrationDirectory(filepath.Dir(dir), filepath.Base(dir)+"/records"); err != nil {
		return nil, err
	}
	entries, err := os.ReadDir(filepath.Join(dir, "records"))
	if os.IsNotExist(err) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	if len(entries) > maxLegacyDecisions {
		return nil, errors.New("too many durable legacy declarations")
	}
	out := make([]LegacyDeclaration, 0, len(entries))
	paths := map[string]bool{}
	for _, e := range entries {
		data, err := readStepHistoryFile(filepath.Join(dir, "records", e.Name()), 1<<20)
		if err != nil {
			return nil, err
		}
		var r LegacyDeclaration
		if err := migrationJSON(data, &r); err != nil {
			return nil, err
		}
		if err := validateLegacyRecord(r, scope); err != nil {
			return nil, err
		}
		if e.Name() != key(r.Request.DecisionID)+".json" || paths[r.Request.Path] {
			return nil, errors.New("conflicting legacy decision identity or run path")
		}
		paths[r.Request.Path] = true
		out = append(out, r)
	}
	digests := map[string]LegacyDeclaration{}
	for _, r := range out {
		digests[migrationDigest(r)] = r
	}
	for _, r := range out {
		for _, digest := range r.Preview.Decisions {
			prior, ok := digests[digest]
			if !ok || !prior.RecordedAt.Before(r.RecordedAt) {
				return nil, errors.New("legacy decision lost its prior append-only authority")
			}
		}
	}
	return out, nil
}

// The closed structural region contains exactly events.jsonl. In particular,
// a cursor, execution log, hidden file or even an empty extra directory refuses.
func legacyManifest(root, relative string) (string, error) {
	if !validUnscopedRunPath(relative) {
		return "", errors.New("invalid canonical legacy run path")
	}
	if reason, ok := runDeclarationEligibility(root, strings.TrimPrefix(relative, unscopedRunPrefix)); !ok {
		return "", fmt.Errorf("legacy run is not declarable: %s", reason)
	}
	dir := filepath.Join(root, filepath.FromSlash(relative))
	entries, err := os.ReadDir(dir)
	if err != nil {
		return "", err
	}
	if len(entries) != 1 || entries[0].Name() != "events.jsonl" || !entries[0].Type().IsRegular() {
		return "", errors.New("durable legacy declaration requires exactly one regular events.jsonl and no other entry")
	}
	if err := legacySingleLink(filepath.Join(dir, "events.jsonl")); err != nil {
		return "", err
	}
	manifest, err := RunDirectoryManifest(dir)
	if err != nil {
		return "", err
	}
	data, err := readStepHistoryFile(filepath.Join(dir, "events.jsonl"), maxRunManifestFile)
	if err != nil {
		return "", err
	}
	scan := bufio.NewScanner(bytes.NewReader(data))
	scan.Buffer(make([]byte, 4096), 1<<20)
	count := 0
	for scan.Scan() {
		line := bytes.TrimSpace(scan.Bytes())
		if len(line) == 0 {
			continue
		}
		var e struct {
			Time time.Time                  `json:"time"`
			Type string                     `json:"type"`
			Data map[string]json.RawMessage `json:"data"`
		}
		if err := migrationJSON(line, &e); err != nil {
			return "", err
		}
		if e.Type != "run.created" && e.Type != "run.phase" {
			return "", errors.New("legacy run has execution or unsupported event history; use per-idea migration")
		}
		// A closed vocabulary prevents unfamiliar accounting fields from becoming
		// an implicit amnesty. Only descriptive creation data or a known non-cycle
		// phase can be declared; no idea identity or charge field is accepted.
		for name, value := range e.Data {
			var s string
			if json.Unmarshal(value, &s) != nil {
				return "", errors.New("non-string legacy event data")
			}
			switch {
			case e.Type == "run.created" && (name == "mode" || name == "task"):
			case e.Type == "run.phase" && name == "action" && knownNonCyclePhase(s):
			default:
				return "", errors.New("legacy run has identity, charge or unsupported event data; use per-idea migration")
			}
		}
		if e.Type == "run.phase" && e.Data["action"] == nil {
			return "", errors.New("legacy phase has no known non-cycle action")
		}
		count++
	}
	if err := scan.Err(); err != nil {
		return "", err
	}
	if count == 0 {
		return "", errors.New("empty legacy history is not a declaration witness")
	}
	// Refuse mutation between the strict structural read and manifest read.
	digest, err := RunDirectoryManifestDigest(dir)
	if err != nil {
		return "", err
	}
	if digest != migrationDigest(manifest) || len(manifest) != 1 || manifest[0].SHA256 != key(string(data)) {
		return "", errors.New("legacy history changed during inspection")
	}
	return digest, nil
}

func legacyCopies(roots []string, path, expected string, required []string) ([]string, string, error) {
	found := []string{}
	for _, root := range roots {
		if err := migrationDirectory(root, path); err != nil {
			return nil, "", err
		}
		_, err := os.Lstat(filepath.Join(root, filepath.FromSlash(path)))
		if os.IsNotExist(err) {
			continue
		}
		if err != nil {
			return nil, "", err
		}
		digest, err := legacyManifest(root, path)
		if err != nil {
			return nil, "", err
		}
		if expected != "" && digest != expected {
			return nil, "", errors.New("declared legacy run manifest changed or copies diverge")
		}
		expected = digest
		found = append(found, root)
	}
	sort.Strings(found)
	if len(found) == 0 {
		return nil, "", errors.New("declared legacy run is missing from all visible roots")
	}
	for _, prior := range required {
		n := sort.SearchStrings(found, prior)
		if n == len(found) || found[n] != prior {
			return nil, "", errors.New("a previously declared legacy copy is missing or no longer visible")
		}
	}
	return found, expected, nil
}

func validateLegacyCopies(roots []string, records []LegacyDeclaration) error {
	for _, r := range records {
		if _, _, err := legacyCopies(roots, r.Preview.Path, r.Preview.ManifestSHA256, r.Preview.Roots); err != nil {
			return err
		}
	}
	return nil
}

func InspectLegacyDeclaration(ctx context.Context, root, path string) (LegacyPreview, error) {
	dir, scope, roots, err := legacyScope(ctx, root)
	if err != nil {
		return LegacyPreview{}, err
	}
	records, err := readLegacyRecords(dir, scope)
	if err != nil {
		return LegacyPreview{}, err
	}
	if err := validateLegacyCopies(roots, records); err != nil {
		return LegacyPreview{}, err
	}
	copies, digest, err := legacyCopies(roots, path, "", nil)
	if err != nil {
		return LegacyPreview{}, err
	}
	p := LegacyPreview{Version: 1, Scope: scope, Path: path, ManifestSHA256: digest, History: legacyUnknownHistory, Roots: copies, Decisions: []string{}}
	for _, r := range records {
		p.Decisions = append(p.Decisions, migrationDigest(r))
	}
	sort.Strings(p.Decisions)
	p.PreviewSHA256 = legacyPreviewDigest(p)
	return p, nil
}

func ApplyLegacyDeclaration(ctx context.Context, root string, q LegacyDeclarationRequest) (LegacyDeclaration, error) {
	if !q.WritersStopped || !q.AcknowledgeUnknown || !validUnscopedRunPath(q.Path) || !validCycleDecision(q.DecisionID, q.Reason, q.ExpectedPreviewSHA256) {
		return LegacyDeclaration{}, errors.New("legacy declaration requires exact preview, decision/reason, stopped writers and acknowledgement of unknown history")
	}
	dir, scope, roots, err := legacyScope(ctx, root)
	if err != nil {
		return LegacyDeclaration{}, err
	}
	if err := migrationDirectory(filepath.Dir(dir), filepath.Base(dir)+"/records"); err != nil {
		return LegacyDeclaration{}, err
	}
	release, err := AcquireResourceGuard(ctx, dir)
	if err != nil {
		return LegacyDeclaration{}, err
	}
	defer release()
	records, err := readLegacyRecords(dir, scope)
	if err != nil {
		return LegacyDeclaration{}, err
	}
	if err := validateLegacyCopies(roots, records); err != nil {
		return LegacyDeclaration{}, err
	}
	for _, r := range records {
		if r.Request.DecisionID == q.DecisionID {
			if !reflect.DeepEqual(r.Request, q) {
				return LegacyDeclaration{}, errors.New("conflicting legacy decision ID reuse")
			}
			return r, nil
		}
		if r.Request.Path == q.Path {
			return LegacyDeclaration{}, errors.New("legacy run already has a decision; replay its exact request")
		}
	}
	if len(records) >= maxLegacyDecisions {
		return LegacyDeclaration{}, errors.New("too many durable legacy declarations")
	}
	preview, err := InspectLegacyDeclaration(ctx, root, q.Path)
	if err != nil {
		return LegacyDeclaration{}, err
	}
	if preview.PreviewSHA256 != q.ExpectedPreviewSHA256 {
		return LegacyDeclaration{}, errors.New("legacy preview changed; inspect again before applying")
	}
	r := LegacyDeclaration{Version: 1, Scope: scope, RecordedAt: time.Now().UTC(), Authority: "attended-operator", Request: q, Preview: preview}
	if err := validateLegacyRecord(r, scope); err != nil {
		return LegacyDeclaration{}, err
	}
	if err := os.MkdirAll(filepath.Join(dir, "records"), 0700); err != nil {
		return LegacyDeclaration{}, err
	}
	data, err := json.Marshal(r)
	if err != nil {
		return LegacyDeclaration{}, err
	}
	path := filepath.Join(dir, "records", key(q.DecisionID)+".json")
	if err := writeSynced(path, data); err != nil {
		return LegacyDeclaration{}, err
	}
	actual, err := readStepHistoryFile(path, 1<<20)
	if err != nil || !bytes.Equal(actual, data) {
		return LegacyDeclaration{}, errors.New("legacy decision readback failed; preserve state and replay exact decision")
	}
	return r, nil
}
