package budget

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"math"
	"math/big"
	"os"
	"path/filepath"
	"reflect"
	"sort"
	"strconv"
	"strings"
	"time"

	"parley-deck-cli/internal/telemetry"
)

const maxMigrationItems = 10000

// MigrationSource binds metadata and canonical artifacts without copying their
// contents (which may contain private prompts or diagnostics) into a decision.
type MigrationSource struct {
	Root   string `json:"root"`
	Path   string `json:"path"`
	SHA256 string `json:"sha256"`
}

type HistoricalLaunch struct {
	ID               string    `json:"id"`
	RequestedAt      time.Time `json:"requested_at"`
	TerminalSHA256   string    `json:"terminal_sha256"`
	Classification   string    `json:"classification"`
	ObservedMicros   *int64    `json:"observed_micros"`
	ObservationBasis string    `json:"observation_basis"`
}

type LaunchMigrationInventory struct {
	Version             int                `json:"version"`
	Scope               string             `json:"scope"`
	Idea                string             `json:"idea"`
	Roots               []string           `json:"roots"`
	Sources             []MigrationSource  `json:"sources"`
	Launches            []HistoricalLaunch `json:"launches"`
	Earliest            *time.Time         `json:"earliest"`
	LegacyStartFloor    int                `json:"legacy_start_floor"`
	HasCanonicalHistory bool               `json:"has_canonical_history"`
	HistorySHA256       string             `json:"history_sha256"`
}

func migrationDigest(v any) string {
	raw, _ := json.Marshal(v)
	return key(string(raw))
}

func (i LaunchMigrationInventory) digest() string {
	i.HistorySHA256 = ""
	return migrationDigest(i)
}

type migrationScanner struct {
	i      LaunchMigrationInventory
	bytes  int64
	seen   map[string]HistoricalLaunch
	starts map[string]time.Time
	ctx    context.Context
}

func migrationDirectory(root, relative string) error {
	path := root
	for _, part := range strings.Split(relative, "/") {
		path = filepath.Join(path, part)
		st, err := os.Lstat(path)
		if os.IsNotExist(err) {
			return nil
		}
		if err != nil {
			return err
		}
		if !st.IsDir() {
			return errors.New("migration history contains an aliased or non-directory parent")
		}
	}
	return nil
}

func (s *migrationScanner) file(root, relative string, limit int64) ([]byte, error) {
	if err := s.ctx.Err(); err != nil {
		return nil, err
	}
	// Lstat the intermediate directories too: a regular leaf does not prove
	// that an aliased parent stayed inside this worktree's history.
	parent := root
	parts := strings.Split(relative, "/")
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
	data, err := readStepHistoryFile(filepath.Join(root, filepath.FromSlash(relative)), limit)
	if err != nil {
		return nil, err
	}
	s.bytes += int64(len(data))
	if len(s.i.Sources) >= maxMigrationItems || s.bytes > 64<<20 {
		return nil, errors.New("migration history exceeds the bounded inventory; preserve it for explicit archival migration")
	}
	s.i.Sources = append(s.i.Sources, MigrationSource{root, relative, key(string(data))})
	return data, nil
}

func (s *migrationScanner) earlier(at time.Time) error {
	if at.IsZero() || at.After(time.Now().UTC()) {
		return errors.New("historical accounting timestamp is missing or in the future")
	}
	if s.i.Earliest == nil || at.Before(*s.i.Earliest) {
		t := at.UTC()
		s.i.Earliest = &t
	}
	return nil
}

// InspectLaunchMigration has no mutation path. Complete metadata for matching
// invocations is mandatory; absence of a terminal is not proof of a dead process.
func InspectLaunchMigration(ctx context.Context, root, idea string) (LaunchMigrationInventory, error) {
	_, scope, roots, err := launchScope(ctx, root, idea, true)
	if err != nil {
		return LaunchMigrationInventory{}, err
	}
	sort.Strings(roots)
	s := migrationScanner{i: LaunchMigrationInventory{Version: 1, Scope: scope, Idea: idea, Roots: roots, Sources: []MigrationSource{}, Launches: []HistoricalLaunch{}}, seen: map[string]HistoricalLaunch{}, starts: map[string]time.Time{}, ctx: ctx}
	for _, origin := range roots {
		paths := []string{".parley-runtime/invocations", "parley-deck/runs"}
		if idea != "" {
			paths = append(paths, "parley-deck/ideas/"+idea)
		}
		for _, path := range paths {
			if err := migrationDirectory(origin, path); err != nil {
				return s.i, err
			}
		}
		if err := s.invocations(origin, idea); err != nil {
			return s.i, err
		}
		if err := s.runs(origin, idea); err != nil {
			return s.i, err
		}
		if idea != "" {
			if err := s.artifacts(origin, idea); err != nil {
				return s.i, err
			}
		}
	}
	for id := range s.starts {
		if launch, exists := s.seen[id]; exists && launch.Classification == "retained-pre-start-refusal" {
			return s.i, errors.New("historical process start contradicts a pre-start refusal")
		} else if !exists {
			s.i.LegacyStartFloor++
		}
	}
	for _, launch := range s.seen {
		s.i.Launches = append(s.i.Launches, launch)
	}
	sort.Slice(s.i.Launches, func(a, b int) bool { return s.i.Launches[a].ID < s.i.Launches[b].ID })
	sort.Slice(s.i.Sources, func(a, b int) bool {
		x, y := s.i.Sources[a], s.i.Sources[b]
		if x.Root == y.Root {
			return x.Path < y.Path
		}
		return x.Root < y.Root
	})
	s.i.HistorySHA256 = s.i.digest()
	return s.i, nil
}

func migrationJSON(data []byte, out any) error {
	if err := validateHistoryJSON(json.NewDecoder(bytes.NewReader(data)), 0); err != nil {
		return errors.New("malformed or duplicate historical fields")
	}
	if err := migrationShape(data, reflect.TypeOf(out).Elem()); err != nil {
		return err
	}
	d := json.NewDecoder(bytes.NewReader(data))
	d.DisallowUnknownFields()
	if err := d.Decode(out); err != nil {
		return errors.New("unsupported historical accounting schema")
	}
	if d.Decode(new(any)) != io.EOF {
		return errors.New("trailing historical data")
	}
	return nil
}

// A normalized digest must not turn deleted zero-valued authority into an
// explicit assertion. Match required fields recursively, using the persisted
// schema's JSON tags. Nullable observations remain distinct from missing keys.
func migrationShape(data []byte, schema reflect.Type) error {
	if schema == reflect.TypeOf(json.RawMessage{}) {
		return nil
	}
	null := bytes.Equal(bytes.TrimSpace(data), []byte("null"))
	if schema.Kind() == reflect.Pointer {
		if null {
			return nil
		}
		return migrationShape(data, schema.Elem())
	}
	if null && schema.Kind() != reflect.Slice && schema.Kind() != reflect.Map {
		return errors.New("non-null historical field is null")
	}
	if schema == reflect.TypeOf(time.Time{}) {
		return nil
	}
	switch schema.Kind() {
	case reflect.Struct:
		var fields map[string]json.RawMessage
		if err := json.Unmarshal(data, &fields); err != nil {
			return err
		}
		for n := 0; n < schema.NumField(); n++ {
			field := schema.Field(n)
			tag := strings.Split(field.Tag.Get("json"), ",")
			if field.PkgPath != "" || tag[0] == "-" {
				continue
			}
			raw, exists := fields[tag[0]]
			if !exists {
				if len(tag) > 1 && tag[1] == "omitempty" {
					continue
				}
				return fmt.Errorf("missing required historical field %s", tag[0])
			}
			if err := migrationShape(raw, field.Type); err != nil {
				return err
			}
		}
	case reflect.Slice:
		var elements []json.RawMessage
		if err := json.Unmarshal(data, &elements); err != nil {
			return err
		}
		for _, raw := range elements {
			if err := migrationShape(raw, schema.Elem()); err != nil {
				return err
			}
		}
	case reflect.Map:
		var elements map[string]json.RawMessage
		if err := json.Unmarshal(data, &elements); err != nil {
			return err
		}
		for _, raw := range elements {
			if err := migrationShape(raw, schema.Elem()); err != nil {
				return err
			}
		}
	}
	return nil
}

func migrationRecord(data []byte, id, kind string) (telemetry.Record, error) {
	var r telemetry.Record
	if err := migrationJSON(data, &r); err != nil {
		return r, err
	}
	var fields map[string]json.RawMessage
	_ = json.Unmarshal(data, &fields)
	for _, name := range []string{"schema_version", "type", "invocation_id", "metadata", "requested_at", "started_at", "completed_at", "pid", "outcome"} {
		if _, ok := fields[name]; !ok {
			return r, errors.New("incomplete historical invocation")
		}
	}
	var meta map[string]json.RawMessage
	_ = json.Unmarshal(fields["metadata"], &meta)
	for _, name := range []string{"idea", "phase", "run_id", "segment_id"} {
		v, ok := meta[name]
		if !ok || bytes.Equal(v, []byte("null")) {
			return r, errors.New("missing historical invocation identity")
		}
	}
	if r.SchemaVersion != 1 || r.InvocationID != id || r.Type != kind || r.RequestedAt.IsZero() || len(id) == 0 || len(id) > 128 {
		return r, errors.New("historical invocation identity or schema mismatch")
	}
	return r, nil
}

func (s *migrationScanner) invocations(root, idea string) error {
	parent := filepath.Join(root, ".parley-runtime")
	if st, err := os.Lstat(parent); err == nil && !st.IsDir() {
		return errors.New("historical runtime must be a real directory")
	} else if err != nil && !os.IsNotExist(err) {
		return err
	}
	entries, err := stepHistoryDirs(filepath.Join(parent, "invocations"))
	if err != nil {
		return err
	}
	for _, entry := range entries {
		base := ".parley-runtime/invocations/" + entry.Name() + "/"
		raw, err := s.file(root, base+"requested.json", 1<<20)
		if err != nil {
			return fmt.Errorf("cannot inventory historical request: %w", err)
		}
		req, err := migrationRecord(raw, entry.Name(), "invocation.requested")
		if err != nil {
			return err
		}
		if req.Metadata.Idea != idea {
			continue
		}
		if req.StartedAt != nil || req.CompletedAt != nil || req.PID != nil || req.Outcome != nil {
			return errors.New("requested record carries later execution state")
		}
		if err := s.earlier(req.RequestedAt); err != nil {
			return err
		}
		raw, err = s.file(root, base+"terminal.json", 1<<20)
		if err != nil {
			return fmt.Errorf("matching historical invocation is not terminal or cannot be read; stop and reconcile writers before migration: %w", err)
		}
		terminal, err := migrationRecord(raw, entry.Name(), "invocation.terminal")
		if err != nil {
			return err
		}
		if !reflect.DeepEqual(terminal.Metadata, req.Metadata) || !terminal.RequestedAt.Equal(req.RequestedAt) || terminal.CompletedAt == nil || terminal.CompletedAt.Before(req.RequestedAt) || terminal.CompletedAt.After(time.Now().UTC()) || terminal.Outcome == nil {
			return errors.New("terminal record differs from its historical request")
		}
		var started *telemetry.Record
		startRaw, err := s.file(root, base+"started.json", 1<<20)
		if err == nil {
			v, err := migrationRecord(startRaw, entry.Name(), "invocation.started")
			if err != nil {
				return err
			}
			started = &v
		} else if !os.IsNotExist(err) {
			return err
		}
		if (terminal.StartedAt != nil) != (started != nil) || (terminal.PID != nil) != (started != nil) {
			return errors.New("incomplete historical process-start evidence")
		}
		if started != nil && (terminal.PID == nil || *terminal.PID <= 0 || started.PID == nil || *started.PID != *terminal.PID || started.StartedAt == nil || !started.StartedAt.Equal(*terminal.StartedAt) || terminal.StartedAt.Before(req.RequestedAt) || terminal.CompletedAt.Before(*terminal.StartedAt) || !reflect.DeepEqual(started.Metadata, req.Metadata) || !started.RequestedAt.Equal(req.RequestedAt) || started.CompletedAt != nil || started.Outcome != nil) {
			return errors.New("conflicting historical process-start evidence")
		}
		cost, err := migrationObservedMicros(terminal.Outcome.Usage.CostUSD)
		if err != nil {
			return err
		}
		classification := "attempt"
		outcome := terminal.Outcome
		failure := ""
		if outcome.FailureClass != nil {
			failure = *outcome.FailureClass
		}
		// These typed failures happen before spawn. Neither an unobserved
		// handoff nor an ordinary failed start receives this exemption.
		if started == nil && outcome.Status == "failed" && (failure == "budget_refused" || failure == "protocol_context_refused") {
			for _, tokens := range []*int64{outcome.Usage.InputTokens, outcome.Usage.OutputTokens, outcome.Usage.TotalTokens, outcome.Usage.CacheReadTokens, outcome.Usage.CacheWriteTokens} {
				if tokens != nil && *tokens != 0 {
					return errors.New("pre-start refusal conflicts with observed token usage")
				}
			}
			if cost != nil && *cost != 0 || outcome.ExitCode != nil || outcome.Observation.TruncatedInput || outcome.Observation.FirstActivityMS != nil || outcome.Observation.StdoutBytes != nil && *outcome.Observation.StdoutBytes != 0 || outcome.Observation.StderrBytes != nil && *outcome.Observation.StderrBytes != 0 {
				return errors.New("pre-start refusal conflicts with observed work")
			}
			classification = "retained-pre-start-refusal"
		}
		if started == nil && outcome.Status == "unobserved-handoff" {
			classification, cost = "unobserved-handoff", nil
		}
		launch := HistoricalLaunch{req.InvocationID, req.RequestedAt, key(string(raw)), classification, cost, outcome.Usage.CostBasis}
		if old, ok := s.seen[launch.ID]; ok && !reflect.DeepEqual(old, launch) {
			return errors.New("conflicting historical invocation copies across worktrees")
		}
		s.seen[launch.ID] = launch
		if len(s.seen) > maxMigrationItems {
			return errors.New("too many historical invocations")
		}
	}
	return nil
}

// Round observations upward to whole microdollars without recasting a CLI
// estimate as an invoice. Its exact basis and terminal hash stay in the import.
func migrationObservedMicros(usd *float64) (*int64, error) {
	if usd == nil {
		return nil, nil
	}
	if math.IsNaN(*usd) || math.IsInf(*usd, 0) || *usd < 0 {
		return nil, errors.New("invalid historical monetary observation")
	}
	r, ok := new(big.Rat).SetString(strconv.FormatFloat(*usd, 'g', -1, 64))
	if !ok {
		return nil, errors.New("unrepresentable historical cost")
	}
	r.Mul(r, big.NewRat(1_000_000, 1))
	q, rem := new(big.Int), new(big.Int)
	q.QuoRem(r.Num(), r.Denom(), rem)
	if rem.Sign() > 0 {
		q.Add(q, big.NewInt(1))
	}
	if !q.IsInt64() {
		return nil, ErrCostOverflow
	}
	v := q.Int64()
	return &v, nil
}

func (s *migrationScanner) artifacts(root, idea string) error {
	relative := filepath.ToSlash(filepath.Join("parley-deck", "ideas", idea))
	base := filepath.Join(root, filepath.FromSlash(relative))
	return filepath.WalkDir(base, func(path string, e os.DirEntry, err error) error {
		if os.IsNotExist(err) && path == base {
			return nil
		}
		if err != nil {
			return err
		}
		if e.IsDir() {
			return nil
		}
		rel, err := filepath.Rel(root, path)
		if err != nil {
			return err
		}
		if _, err := s.file(root, filepath.ToSlash(rel), 16<<20); err != nil {
			return err
		}
		if filepath.Base(path) != "00-prompt.md" {
			s.i.HasCanonicalHistory = true
		}
		return nil
	})
}

func (s *migrationScanner) runs(root, idea string) error {
	base := filepath.Join(root, "parley-deck", "runs")
	entries, err := stepHistoryDirs(base)
	if err != nil {
		return err
	}
	for _, entry := range entries {
		rel := "parley-deck/runs/" + entry.Name() + "/"
		raw, err := s.file(root, rel+"events.jsonl", 16<<20)
		if os.IsNotExist(err) {
			if _, e := os.Lstat(filepath.Join(base, entry.Name(), "driver.json")); e == nil || !os.IsNotExist(e) {
				return errors.New("historical driver cursor has no readable run identity")
			}
			continue
		}
		if err != nil {
			return err
		}
		identity := ""
		bind := func(raw json.RawMessage, required bool) error {
			if raw == nil && !required {
				return nil
			}
			var name string
			if json.Unmarshal(raw, &name) != nil || strings.TrimSpace(name) == "" || identity != "" && identity != name {
				return errors.New("historical run identity is missing or conflicting")
			}
			identity = name
			return nil
		}
		cursor, err := s.file(root, rel+"driver.json", 1<<20)
		if err == nil {
			var fields map[string]json.RawMessage
			if err := migrationJSON(cursor, &fields); err != nil {
				return err
			}
			for _, name := range []string{"idea", "idea_slug"} {
				if err := bind(fields[name], false); err != nil {
					return err
				}
			}
		} else if !os.IsNotExist(err) {
			return err
		}
		type event struct {
			Time time.Time                  `json:"time"`
			Type string                     `json:"type"`
			Data map[string]json.RawMessage `json:"data,omitempty"`
		}
		var events []event
		scan := bufio.NewScanner(bytes.NewReader(raw))
		scan.Buffer(make([]byte, 4096), 1<<20)
		for scan.Scan() {
			line := bytes.TrimSpace(scan.Bytes())
			if len(line) == 0 {
				continue
			}
			var e event
			if err := migrationJSON(line, &e); err != nil {
				return err
			}
			if e.Type == "" {
				return errors.New("historical event has no type")
			}
			for _, name := range []string{"idea", "idea_slug"} {
				if err := bind(e.Data[name], name == "idea" && e.Type == "run.created"); err != nil {
					return err
				}
			}
			events = append(events, e)
			if len(events) > maxMigrationItems {
				return errors.New("too many historical run events")
			}
		}
		if err := scan.Err(); err != nil {
			return err
		}
		if identity == "" && (len(events) > 0 || len(cursor) > 0) {
			return errors.New("historical run needs an explicit recoverable idea identity")
		}
		if identity != idea {
			continue
		}
		for n, e := range events {
			if err := s.earlier(e.Time); err != nil {
				return err
			}
			if e.Type != "agent.started" {
				continue
			}
			id := ""
			for _, name := range []string{"invocation_id", "proc_marker"} {
				if v, ok := e.Data[name]; ok {
					if bytes.Equal(bytes.TrimSpace(v), []byte("null")) || json.Unmarshal(v, &id) != nil {
						return errors.New("invalid historical process identity")
					}
					if id != "" {
						break
					}
				}
			}
			if id == "" {
				id = "legacy-event:" + entry.Name() + ":" + strconv.Itoa(n) + ":" + migrationDigest(e)
			}
			if old, ok := s.starts[id]; ok && !old.Equal(e.Time) {
				return errors.New("historical process identity was reused")
			}
			s.starts[id] = e.Time
			if len(s.starts) > maxMigrationItems {
				return errors.New("too many historical starts")
			}
		}
	}
	return nil
}
