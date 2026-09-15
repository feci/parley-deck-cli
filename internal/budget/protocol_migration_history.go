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

// Protocol counters are not model invocations. Old telemetry does not carry
// synchronous operation IDs, so a group of participants cannot be charged once
// per child or assumed to represent a complete lifetime total.
type ProtocolCountEvidence struct {
	Root  string `json:"root"`
	Path  string `json:"path"`
	Rule  string `json:"rule"`
	Floor int    `json:"floor"`
}

type ProtocolMigrationInventory struct {
	Version              int                      `json:"version"`
	Scope                string                   `json:"scope"`
	Idea                 string                   `json:"idea"`
	IdeaPath             string                   `json:"idea_path"`
	Kind                 Kind                     `json:"kind"`
	History              LaunchMigrationInventory `json:"history"`
	Earliest             *time.Time               `json:"earliest"`
	LowerBound           int                      `json:"lower_bound"`
	Evidence             []ProtocolCountEvidence  `json:"evidence"`
	UngroupedInvocations []string                 `json:"ungrouped_invocations"`
	HistorySHA256        string                   `json:"history_sha256"`
}

func (i ProtocolMigrationInventory) digest() string {
	i.HistorySHA256 = ""
	return migrationDigest(i)
}

func protocolMigrationPath(idea, relative string) (string, error) {
	if relative == "" {
		relative = filepath.ToSlash(filepath.Join("parley-deck", "ideas", idea))
	}
	if !strings.HasPrefix(relative, "parley-deck/") || strings.ContainsAny(relative, "\\\x00\r\n") || filepath.IsAbs(relative) || filepath.ToSlash(filepath.Clean(relative)) != relative || strings.Contains(relative, "/../") {
		return "", errors.New("protocol migration idea path must be canonical and inside the deck")
	}
	return relative, nil
}

func protocolMigrationScope(ctx context.Context, root, idea string, kind Kind) (string, string, error) {
	if idea == "" {
		return "", "", errors.New("protocol migration requires an idea")
	}
	if kind == DriverStep {
		d, s, _, err := stepScope(ctx, root, idea, false)
		return d, s, err
	}
	d, s, _, err := cycleScope(ctx, root, idea, kind, false)
	return d, s, err
}

// InspectProtocolMigration freezes the source inventory and observed floors,
// not an inferred count of all past actions. TotalActions is a separate explicit
// operator reconciliation. The optional idea path supports nested pipeline ideas.
func InspectProtocolMigration(ctx context.Context, root, idea string, kind Kind, ideaPath string) (ProtocolMigrationInventory, error) {
	_, scope, err := protocolMigrationScope(ctx, root, idea, kind)
	if err != nil {
		return ProtocolMigrationInventory{}, err
	}
	relative, err := protocolMigrationPath(idea, ideaPath)
	if err != nil {
		return ProtocolMigrationInventory{}, err
	}
	history, err := InspectLaunchMigration(ctx, root, idea)
	if err != nil {
		return ProtocolMigrationInventory{}, err
	}
	i := ProtocolMigrationInventory{Version: 1, Scope: scope, Idea: idea, IdeaPath: relative, Kind: kind, History: history, Earliest: history.Earliest, Evidence: []ProtocolCountEvidence{}, UngroupedInvocations: []string{}}
	// Re-read bounded sources against their hashes. This detects changes while
	// deriving counts and bounds the total of ordinary plus nested-idea sources.
	scanner := migrationScanner{i: history, ctx: ctx}
	raw := map[string][]byte{}
	for _, source := range history.Sources {
		data, e := readStepHistoryFile(filepath.Join(source.Root, filepath.FromSlash(source.Path)), 16<<20)
		if e != nil {
			return i, e
		}
		if key(string(data)) != source.SHA256 {
			return i, fmt.Errorf("%w while deriving protocol history", errHistoryChanged)
		}
		scanner.bytes += int64(len(data))
		if scanner.bytes > 64<<20 {
			return i, errors.New("protocol history exceeds inventory byte bound")
		}
		raw[source.Root+"\x00"+source.Path] = data
	}
	if relative != filepath.ToSlash(filepath.Join("parley-deck", "ideas", idea)) {
		for _, origin := range history.Roots {
			if err := migrationDirectory(origin, relative); err != nil {
				return i, err
			}
			base := filepath.Join(origin, filepath.FromSlash(relative))
			err := filepath.WalkDir(base, func(path string, e os.DirEntry, err error) error {
				if os.IsNotExist(err) && path == base {
					return nil
				}
				if err != nil {
					return err
				}
				if e.IsDir() {
					return nil
				}
				rel, err := filepath.Rel(origin, path)
				if err != nil {
					return err
				}
				if _, exists := raw[origin+"\x00"+filepath.ToSlash(rel)]; exists {
					return nil
				}
				data, err := scanner.file(origin, filepath.ToSlash(rel), 16<<20)
				if err == nil {
					raw[origin+"\x00"+filepath.ToSlash(rel)] = data
				}
				return err
			})
			if err != nil {
				return i, err
			}
		}
	}
	sort.Slice(scanner.i.Sources, func(a, b int) bool {
		x, y := scanner.i.Sources[a], scanner.i.Sources[b]
		if x.Root == y.Root {
			return x.Path < y.Path
		}
		return x.Root < y.Root
	})
	scanner.i.HistorySHA256 = scanner.i.digest()
	i.History = scanner.i
	add := func(root, path, rule string, n int) error {
		if n < 0 || n > maxMigrationItems || len(i.Evidence) >= maxMigrationItems {
			return errors.New("protocol accounting floor exceeds import bound")
		}
		i.Evidence = append(i.Evidence, ProtocolCountEvidence{root, path, rule, n})
		return nil
	}
	launches := map[string]HistoricalLaunch{}
	for _, l := range history.Launches {
		launches[l.ID] = l
	}
	seenRequests := map[string]bool{}
	type run struct {
		root, path     string
		events, cursor []byte
	}
	runs := map[string]*run{}
	for _, source := range i.History.Sources {
		data := raw[source.Root+"\x00"+source.Path]
		if strings.HasPrefix(source.Path, ".parley-runtime/invocations/") && strings.HasSuffix(source.Path, "/requested.json") {
			id := filepath.Base(filepath.Dir(source.Path))
			l, matching := launches[id]
			if !matching || seenRequests[id] {
				continue
			}
			seenRequests[id] = true
			r, err := migrationRecord(data, id, "invocation.requested")
			if err != nil {
				return i, err
			}
			phase := r.Metadata.Phase
			matches, unknown := false, false
			if kind == DriverStep {
				_, cycle := CycleKindForPhase(phase)
				switch phase {
				case "implementation", "review", "review-consensus", "consensus", "final", "consensus-signoff", "review-signoff", "goal-check", "evidence-verification":
					matches = true
				default:
					matches, unknown = cycle, !knownNonCyclePhase(phase) && !cycle
				}
			} else if k, known := CycleKindForPhase(phase); known {
				matches = k == kind
			} else {
				unknown = !knownNonCyclePhase(phase)
			}
			if matches || unknown {
				i.UngroupedInvocations = append(i.UngroupedInvocations, id)
				// A pre-start launch refusal may follow a charged protocol action.
				// Retain it for the operator; do not assert either zero or one.
				if matches && l.Classification != "retained-pre-start-refusal" && l.Classification != "unobserved-handoff" {
					if err := add(source.Root, source.Path, "at-least-one-typed-action-not-a-per-child-count", 1); err != nil {
						return i, err
					}
				}
			}
		}
		if strings.HasPrefix(source.Path, "parley-deck/runs/") && (strings.HasSuffix(source.Path, "/events.jsonl") || strings.HasSuffix(source.Path, "/driver.json")) {
			path := filepath.ToSlash(filepath.Dir(source.Path))
			key := source.Root + "\x00" + path
			if runs[key] == nil {
				runs[key] = &run{root: source.Root, path: path}
			}
			if strings.HasSuffix(source.Path, "/events.jsonl") {
				runs[key].events = data
			} else {
				runs[key].cursor = data
			}
		}
	}
	// A copied run is one source, not another set of actions. Diverging copies
	// need reconciliation; choosing the smaller or summing both is not safe.
	seenRuns := map[string]string{}
	seenActions := map[string]bool{}
	keys := make([]string, 0, len(runs))
	for key := range runs {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	for _, key := range keys {
		r := runs[key]
		identity, events, err := protocolMigrationEvents(r.events)
		if err != nil {
			return i, err
		}
		if len(r.cursor) > 0 {
			var fields map[string]json.RawMessage
			if err := migrationJSON(r.cursor, &fields); err != nil {
				return i, err
			}
			for _, name := range []string{"idea", "idea_slug"} {
				if value, ok := fields[name]; ok {
					var name string
					if json.Unmarshal(value, &name) != nil || name == "" || identity != "" && identity != name {
						return i, errors.New("conflicting protocol cursor identity")
					}
					identity = name
				}
			}
		}
		if identity != idea {
			continue
		}
		digest := migrationDigest([]string{string(r.events), string(r.cursor)})
		if old, ok := seenRuns[r.path]; ok {
			if old != digest {
				return i, errors.New("conflicting historical protocol run copies")
			}
			continue
		}
		seenRuns[r.path] = digest
		published := 0
		for _, event := range events {
			if raw, ok := event.Data["run_id"]; ok {
				var id string
				if json.Unmarshal(raw, &id) != nil || id != filepath.Base(r.path) {
					return i, errors.New("protocol event run identity differs from its directory")
				}
			}
			if event.Type == "run.phase" {
				var action string
				if json.Unmarshal(event.Data["action"], &action) != nil || action == "" {
					return i, errors.New("protocol phase history lacks action identity")
				}
				switch action {
				case "promoted", "consensus-drafted", "signoffs-requested", "finalized", "reopened", "implemented", "review-opened", "review-drafted", "fixup":
				case "await", "complete", "consensus-ready", "surface-only", "escalated":
					continue
				default:
					return i, errors.New("unsupported historical driver action")
				}
				if kind == DriverStep {
					// Repeated copies of an identical observation are not proof of
					// another action. Distinct published transitions are additive,
					// including transitions in different runs of the same idea.
					digest := migrationDigest(event)
					if !seenActions[digest] {
						seenActions[digest] = true
						published++
					}
				} else if kind == Fixup && action == "fixup" || kind == CrossReview && (action == "promoted" || action == "reopened") {
					// A fixup event can finish an already charged cycle after a
					// crash. Without an operation identity, never count every
					// recovery publication as another code-writing attempt.
					published = 1
				}
			}
			if kind == DriverStep && event.Type == "loop.budget" {
				n, err := protocolCounter(event.Data, "steps", true)
				if err != nil {
					return i, err
				}
				if err := add(r.root, r.path+"/events.jsonl", "reported-lifetime-or-run-step-floor", n); err != nil {
					return i, err
				}
				if raw, ok := event.Data["elapsed_ms"]; ok {
					var elapsed *int64
					if json.Unmarshal(raw, &elapsed) != nil || elapsed == nil || *elapsed < 0 || *elapsed > int64((1<<63-1)/time.Millisecond) {
						return i, errors.New("invalid historical protocol elapsed time")
					}
					at := event.Time.Add(-time.Duration(*elapsed) * time.Millisecond)
					if i.Earliest == nil || at.Before(*i.Earliest) {
						i.Earliest = &at
					}
				}
			}
		}
		rule := "published-action-floor"
		if kind != DriverStep {
			rule = "at-least-one-published-cycle-or-recovery"
		}
		if err := add(r.root, r.path+"/events.jsonl", rule, published); err != nil {
			return i, err
		}
		if len(r.cursor) > 0 {
			var fields map[string]json.RawMessage
			if err := migrationJSON(r.cursor, &fields); err != nil {
				return i, err
			}
			version, err := protocolCounter(fields, "schema_version", false)
			if err != nil || version > 2 {
				return i, errors.New("unsupported historical cursor schema")
			}
			name := "fixup_cycles_published"
			if kind == CrossReview {
				name = "rounds_run"
			}
			if kind != DriverStep {
				n, err := protocolCounter(fields, name, kind == CrossReview || version > 0)
				if err != nil {
					return i, err
				}
				if kind == CrossReview && n > 0 {
					n--
				}
				if err := add(r.root, r.path+"/driver.json", "preserved-cursor-counter-floor", n); err != nil {
					return i, err
				}
			}
		}
	}
	if kind != DriverStep {
		for _, origin := range history.Roots {
			n, err := LegacyCycleFloor(filepath.Join(origin, filepath.FromSlash(relative)), kind)
			if err != nil {
				return i, err
			}
			if err := add(origin, relative, "structural-driver-marker-floor", n); err != nil {
				return i, err
			}
		}
	}
	if kind == Fixup {
		// Different worktrees may retain different completed-cycle markers.
		// Their canonical relative paths identify the same cycle across copies.
		seen := map[string]bool{}
		for _, source := range i.History.Sources {
			prefix := relative + "/review/"
			if !strings.HasPrefix(source.Path, prefix) {
				continue
			}
			parts := strings.Split(strings.TrimPrefix(source.Path, prefix), "/")
			if len(parts) != 2 || parts[1] != ".fixup-done" || seen[source.Path] {
				continue
			}
			if _, ok := cycleRound(parts[0]); !ok {
				continue
			}
			seen[source.Path] = true
			if err := add(source.Root, source.Path, "distinct-completed-fixup-marker", 1); err != nil {
				return i, err
			}
		}
	}
	i.LowerBound, err = protocolEvidenceFloor(kind, i.Evidence)
	if err != nil {
		return i, err
	}
	sort.Strings(i.UngroupedInvocations)
	i.HistorySHA256 = i.digest()
	return i, nil
}

// Only disjoint observed transitions and distinct completed-marker identities
// add. Cursor and loop counters can overlap those observations or other runs.
func protocolEvidenceFloor(kind Kind, evidence []ProtocolCountEvidence) (int, error) {
	floor, actions, markers := 0, 0, 0
	seen := map[string]bool{}
	for _, e := range evidence {
		id := e.Root + "\x00" + e.Path + "\x00" + e.Rule
		if e.Root == "" || e.Path == "" || e.Floor < 0 || e.Floor > maxMigrationItems {
			return 0, errors.New("invalid protocol count evidence")
		}
		switch e.Rule {
		case "published-action-floor":
			if kind != DriverStep || seen[id] {
				return 0, errors.New("invalid additive step evidence")
			}
			actions += e.Floor
		case "distinct-completed-fixup-marker":
			id = e.Path
			if kind != Fixup || e.Floor != 1 || seen[id] {
				return 0, errors.New("invalid distinct fixup evidence")
			}
			markers++
		case "at-least-one-typed-action-not-a-per-child-count", "at-least-one-published-cycle-or-recovery":
			if e.Floor > 1 {
				return 0, errors.New("ungrouped history is not a per-invocation count")
			}
		case "reported-lifetime-or-run-step-floor", "preserved-cursor-counter-floor", "structural-driver-marker-floor":
		default:
			return 0, errors.New("unsupported protocol count evidence")
		}
		seen[id] = true
		if e.Floor > floor {
			floor = e.Floor
		}
		if actions > maxMigrationItems || markers > maxMigrationItems {
			return 0, errors.New("protocol action history exceeds import bound")
		}
	}
	if actions > floor {
		floor = actions
	}
	if markers > floor {
		floor = markers
	}
	return floor, nil
}

type protocolMigrationEvent struct {
	Time time.Time                  `json:"time"`
	Type string                     `json:"type"`
	Data map[string]json.RawMessage `json:"data,omitempty"`
}

func protocolMigrationEvents(data []byte) (string, []protocolMigrationEvent, error) {
	identity := ""
	events := []protocolMigrationEvent{}
	scan := bufio.NewScanner(bytes.NewReader(data))
	scan.Buffer(make([]byte, 4096), 1<<20)
	for scan.Scan() {
		if len(bytes.TrimSpace(scan.Bytes())) == 0 {
			continue
		}
		var e protocolMigrationEvent
		if err := migrationJSON(scan.Bytes(), &e); err != nil {
			return "", nil, err
		}
		for _, name := range []string{"idea", "idea_slug"} {
			raw, ok := e.Data[name]
			if !ok {
				continue
			}
			var value string
			if json.Unmarshal(raw, &value) != nil || value == "" || identity != "" && identity != value {
				return "", nil, errors.New("conflicting protocol run idea")
			}
			identity = value
		}
		events = append(events, e)
		if len(events) > maxMigrationItems {
			return "", nil, errors.New("too many protocol events")
		}
	}
	return identity, events, scan.Err()
}

func protocolCounter(fields map[string]json.RawMessage, name string, required bool) (int, error) {
	raw, ok := fields[name]
	if !ok && !required {
		return 0, nil
	}
	var n *int
	if json.Unmarshal(raw, &n) != nil || n == nil || *n < 0 || *n > maxMigrationItems {
		return 0, fmt.Errorf("unknown or invalid historical counter %s", name)
	}
	return *n, nil
}
