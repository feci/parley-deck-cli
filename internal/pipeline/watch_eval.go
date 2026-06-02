package pipeline

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"gopkg.in/yaml.v3"
)

// SignalSource provides current values for monitored signals. It is an
// interface so the watcher stays vendor-neutral; the default reads a JSON file.
type SignalSource interface {
	Values() (map[string]float64, error)
}

// FileSignalSource reads a JSON object of {signal_name: number} from disk.
type FileSignalSource struct{ Path string }

func (f FileSignalSource) Values() (map[string]float64, error) {
	data, err := os.ReadFile(f.Path)
	if err != nil {
		return nil, fmt.Errorf("read signals: %w", err)
	}
	var out map[string]float64
	if err := json.Unmarshal(data, &out); err != nil {
		return nil, fmt.Errorf("parse signals: %w", err)
	}
	return out, nil
}

// LoadMonitoring reads a structured watcher spec (YAML) from disk.
func LoadMonitoring(path string) (Monitoring, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return Monitoring{}, fmt.Errorf("read monitoring spec: %w", err)
	}
	var m Monitoring
	if err := yaml.Unmarshal(data, &m); err != nil {
		return Monitoring{}, fmt.Errorf("parse monitoring spec: %w", err)
	}
	return m, nil
}

// breaches reports whether an observed value satisfies a threshold's BREACH
// condition. Supported operators: >=, <=, >, <, ==, != followed by a number.
func thresholdBreached(value float64, threshold string) (bool, error) {
	t := strings.TrimSpace(threshold)
	op := ""
	for _, candidate := range []string{">=", "<=", "==", "!=", ">", "<"} {
		if strings.HasPrefix(t, candidate) {
			op = candidate
			break
		}
	}
	if op == "" {
		return false, fmt.Errorf("threshold %q must start with one of >= <= == != > <", threshold)
	}
	bound, err := strconv.ParseFloat(strings.TrimSpace(t[len(op):]), 64)
	if err != nil {
		return false, fmt.Errorf("threshold %q has no numeric bound: %w", threshold, err)
	}
	switch op {
	case ">=":
		return value >= bound, nil
	case "<=":
		return value <= bound, nil
	case ">":
		return value > bound, nil
	case "<":
		return value < bound, nil
	case "==":
		return value == bound, nil
	case "!=":
		return value != bound, nil
	}
	return false, nil
}

// EvaluateBreaches returns one Breach per signal whose observed value satisfies
// its threshold breach condition. Signals absent from values are skipped.
func EvaluateBreaches(m Monitoring, values map[string]float64, now time.Time) ([]Breach, error) {
	var out []Breach
	for _, s := range m.Signals {
		v, ok := values[s.Name]
		if !ok {
			continue
		}
		breached, err := thresholdBreached(v, s.Threshold)
		if err != nil {
			return nil, fmt.Errorf("signal %q: %w", s.Name, err)
		}
		if breached {
			out = append(out, Breach{
				Signal:    s.Name,
				Target:    s.Target,
				Threshold: s.Threshold,
				Observed:  strconv.FormatFloat(v, 'f', -1, 64),
				Class:     s.Class,
				At:        now,
			})
		}
	}
	return out, nil
}

// BreachRecordStatus is the persisted lifecycle of a breach fingerprint.
type BreachRecordStatus string

const (
	BreachOpen     BreachRecordStatus = "open"
	BreachNotified BreachRecordStatus = "notified"
	BreachResolved BreachRecordStatus = "resolved"
)

// BreachRecord persists a breach across watch passes so an ongoing breach is
// deduplicated and a recovered one can be marked resolved (§12.11).
type BreachRecord struct {
	Fingerprint string             `json:"fingerprint"`
	PipelineSlug string            `json:"pipeline_slug"`
	Breach      Breach             `json:"breach"`
	Status      BreachRecordStatus `json:"status"`
	RemediationIdea string         `json:"remediation_idea,omitempty"`
	FirstSeen   time.Time          `json:"first_seen"`
	LastSeen    time.Time          `json:"last_seen"`
}

func breachDir(deckDir, slug string) string {
	return filepath.Join(PipelineDir(deckDir, slug), "breaches")
}

func breachPath(deckDir, slug, fingerprint string) string {
	return filepath.Join(breachDir(deckDir, slug), fingerprint+".json")
}

// LoadBreachRecord reads a persisted breach record; the bool reports existence.
func LoadBreachRecord(deckDir, slug, fingerprint string) (BreachRecord, bool, error) {
	data, err := os.ReadFile(breachPath(deckDir, slug, fingerprint))
	if os.IsNotExist(err) {
		return BreachRecord{}, false, nil
	}
	if err != nil {
		return BreachRecord{}, false, fmt.Errorf("read breach record: %w", err)
	}
	var r BreachRecord
	if err := json.Unmarshal(data, &r); err != nil {
		return BreachRecord{}, false, fmt.Errorf("parse breach record: %w", err)
	}
	return r, true, nil
}

// SaveBreachRecord writes a breach record atomically.
func SaveBreachRecord(deckDir string, r BreachRecord) error {
	if err := os.MkdirAll(breachDir(deckDir, r.PipelineSlug), 0o755); err != nil {
		return fmt.Errorf("create breaches dir: %w", err)
	}
	data, err := json.MarshalIndent(r, "", "  ")
	if err != nil {
		return fmt.Errorf("encode breach record: %w", err)
	}
	data = append(data, '\n')
	path := breachPath(deckDir, r.PipelineSlug, r.Fingerprint)
	tmp := path + ".tmp"
	if err := os.WriteFile(tmp, data, 0o644); err != nil {
		return fmt.Errorf("write breach record: %w", err)
	}
	if err := os.Rename(tmp, path); err != nil {
		return fmt.Errorf("commit breach record: %w", err)
	}
	return nil
}

// ListOpenBreaches returns the fingerprints of currently open/notified records,
// so a watch pass can mark recovered ones resolved.
func ListOpenBreaches(deckDir, slug string) (map[string]BreachRecord, error) {
	out := map[string]BreachRecord{}
	entries, err := os.ReadDir(breachDir(deckDir, slug))
	if os.IsNotExist(err) {
		return out, nil
	}
	if err != nil {
		return nil, err
	}
	for _, e := range entries {
		if e.IsDir() || !strings.HasSuffix(e.Name(), ".json") {
			continue
		}
		fp := strings.TrimSuffix(e.Name(), ".json")
		r, ok, err := LoadBreachRecord(deckDir, slug, fp)
		if err != nil {
			return nil, err
		}
		if ok && r.Status != BreachResolved {
			out[fp] = r
		}
	}
	return out, nil
}
