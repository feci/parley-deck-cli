package quota

import (
	"bytes"
	"crypto/sha256"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"reflect"
	"sort"
	"strings"
	"time"

	"parley-deck-cli/internal/fsutil"
)

const HistoryDir = "quota-history"

type Obligation struct {
	ID     string `json:"id"`
	Kind   string `json:"kind"`
	Agent  string `json:"agent"`
	Path   string `json:"path"`
	SHA256 string `json:"sha256"`
	Text   string `json:"text"`
}

type Batch struct {
	Owner         *Revision    `json:"owner,omitempty"`
	Retained      []Obligation `json:"retained,omitempty"`
	Version       int          `json:"version"`
	ID            string       `json:"id"`
	Idea          string       `json:"idea"`
	RunID         string       `json:"run_id"`
	PriorRevision int          `json:"prior_revision"`
	Policy        Policy       `json:"policy"`
	Decision      Decision     `json:"decision"`
	RecordedAt    time.Time    `json:"recorded_at"`
	Round         string       `json:"round,omitempty"`
	// Expected is the dispatched artifact set, not the whole quorum in a review.
	Expected []string `json:"expected,omitempty"`
}

type History struct {
	Kickoff  *Kickoff `json:"kickoff,omitempty"`
	Batches  []Batch  `json:"batches,omitempty"`
	Current  []string `json:"current"`
	Known    []string `json:"known"`
	Revision int      `json:"revision"`
}

func (h *History) MidIdea() bool {
	return h != nil && h.Kickoff != nil && h.Policy().Enabled && h.Policy().Scope == KickoffAndMidIdea
}
func (b Batch) digest() string {
	b.ID = ""
	raw, _ := json.Marshal(b)
	return fmt.Sprintf("batch-%x", sha256.Sum256(raw))
}
func NewBatch(h *History, run, round string, expected []string, d Decision, now time.Time) Batch {
	b := Batch{Version: 1, Idea: h.Kickoff.Idea, RunID: run, PriorRevision: h.Revision, Policy: h.Policy(), Decision: d, RecordedAt: now.UTC(), Round: round, Expected: append([]string(nil), expected...)}
	b.ID = b.digest()
	return b
}
func (b Batch) Validate(h *History) error {
	if b.Owner != nil {
		return b.validateRevision(h)
	}
	if !h.MidIdea() || b.Version != 1 || b.Idea != h.Kickoff.Idea || b.RunID == "" || b.ID != b.digest() || b.PriorRevision != h.Revision || b.Policy != h.Policy() || b.RecordedAt.IsZero() {
		return fmt.Errorf("contradictory quota batch identity/policy/revision")
	}
	d := b.Decision
	if !d.Applied || d.Block != "" || d.UsableSurvivors < Floor || len(d.Candidates) == 0 || !reflect.DeepEqual(d.Before, h.Current) {
		return fmt.Errorf("invalid quota batch decision")
	}
	removed := map[string]bool{}
	for _, c := range d.Candidates {
		e := c.Evidence
		if !contains(h.Current, c.Agent) || !e.Eligible || e.InvocationID == "" || e.Provenance == "" || e.RuleID == "" || e.ObservedAt.IsZero() || e.Excerpt == "" {
			return fmt.Errorf("invalid quota batch evidence")
		}
		if e.ResetAt != nil && e.ResetAt.Sub(e.ObservedAt) < MinimumReset {
			return fmt.Errorf("invalid quota batch reset")
		}
		if b.Policy.Dropout() || e.RuleID == ParticipantFailureRule {
			if !b.Policy.Dropout() {
				return fmt.Errorf("participant failure evidence contradicts saved trigger")
			}
			if err := ValidateFailureEvidence(e, b.Idea, c.Agent); err != nil {
				return err
			}
		}
		removed[c.Agent] = true
	}
	after := []string{}
	for _, id := range h.Current {
		if !removed[id] {
			after = append(after, id)
		}
	}
	if !reflect.DeepEqual(after, d.After) || len(after) < Floor {
		return fmt.Errorf("contradictory quota batch membership")
	}
	return nil
}
func contains(ids []string, id string) bool {
	for _, v := range ids {
		if v == id {
			return true
		}
	}
	return false
}

// ReadHistory never repairs. Sequence gaps, unknown files and truncated records
// block: absence cannot be interpreted as a policy opt-out or a legacy idea.
func ReadHistory(ideaDir string) (*History, error) {
	k, err := ReadKickoff(ideaDir)
	if err != nil {
		return nil, err
	}
	dir := filepath.Join(ideaDir, HistoryDir)
	entries, err := os.ReadDir(dir)
	if err != nil && !errors.Is(err, os.ErrNotExist) {
		return nil, err
	}
	if k == nil {
		if len(entries) > 0 {
			return nil, fmt.Errorf("missing immutable quota kickoff history")
		}
		return nil, nil
	}
	h := &History{Kickoff: k, Current: append([]string(nil), k.Participants...), Known: append([]string(nil), k.Participants...)}
	var names []string
	for _, e := range entries {
		if strings.HasPrefix(e.Name(), ".stage-") {
			continue
		}
		if e.IsDir() || !strings.HasSuffix(e.Name(), ".json") {
			return nil, fmt.Errorf("unknown quota history record %s", e.Name())
		}
		names = append(names, e.Name())
	}
	sort.Strings(names)
	for i, name := range names {
		if name != fmt.Sprintf("%06d.json", i+1) {
			return nil, fmt.Errorf("missing quota history revision %d", i+1)
		}
		raw, e := os.ReadFile(filepath.Join(dir, name))
		if e != nil {
			return nil, e
		}
		var b Batch
		if e = strictDecode(raw, &b); e != nil {
			return nil, fmt.Errorf("invalid quota history %s: %w", name, e)
		}
		if e = b.Validate(h); e != nil {
			return nil, e
		}
		if e = b.validateOwnerAuthority(ideaDir); e != nil {
			return nil, e
		}
		h.Batches = append(h.Batches, b)
		h.Known = FilterConfirmed(append(h.Known, b.Decision.After...), nil)
		h.Current = append([]string(nil), b.Decision.After...)
		h.Revision++
	}
	return h, nil
}
func strictDecode(raw []byte, v any) error {
	if !uniqueKeys(raw) {
		return fmt.Errorf("malformed or duplicate quota history keys")
	}
	d := json.NewDecoder(bytes.NewReader(raw))
	d.DisallowUnknownFields()
	if err := d.Decode(v); err != nil {
		return err
	}
	var extra any
	if d.Decode(&extra) != io.EOF {
		return fmt.Errorf("trailing quota history")
	}
	return nil
}

// DecodeRevisionRequest applies the same no-duplicate/no-trailing-data grammar
// as immutable records before any authority or filesystem mutation.
func DecodeRevisionRequest(raw []byte, v any) error { return strictDecode(raw, v) }

// DurableWrite publishes complete bytes with checked file and directory sync.
// Exclusive records are immutable. Failures after publication leave readable
// pending state; failures before publication leave no committed record.
func DurableWrite(path string, data []byte, exclusive bool) error {
	if err := durableMkdir(filepath.Dir(path)); err != nil {
		return err
	}
	if err := writeFault("create", path); err != nil {
		return err
	}
	f, err := os.CreateTemp(filepath.Dir(path), ".stage-")
	if err != nil {
		return err
	}
	defer os.Remove(f.Name())
	if err = writeFault("write", path); err == nil {
		var n int
		n, err = f.Write(data)
		if err == nil && n != len(data) {
			err = io.ErrShortWrite
		}
	}
	if err == nil {
		err = writeFault("sync", path)
	}
	if err == nil {
		err = fsutil.SyncFile(f)
	}
	cerr := f.Close()
	if err == nil {
		err = cerr
	}
	if err != nil {
		return err
	}
	if err = writeFault("publish", path); err != nil {
		return err
	}
	if exclusive {
		err = os.Link(f.Name(), path)
		if errors.Is(err, os.ErrExist) {
			old, e := os.ReadFile(path)
			if e == nil && bytes.Equal(old, data) {
				return syncDir(filepath.Dir(path))
			}
			return fmt.Errorf("immutable quota record differs: %s", path)
		}
	} else {
		err = os.Rename(f.Name(), path)
	}
	if err != nil {
		return err
	}
	if err = writeFault("directory-sync", path); err != nil {
		return err
	}
	return syncDir(filepath.Dir(path))
}

func syncDir(path string) error {
	d, err := os.Open(path)
	if err != nil {
		return err
	}
	defer d.Close()
	return fsutil.SyncFile(d)
}

// SyncPath rechecks durability when replay observes a previously published write.
func SyncPath(path string) error {
	f, err := os.Open(path)
	if err != nil {
		return err
	}
	err = fsutil.SyncFile(f)
	closeErr := f.Close()
	if err != nil {
		return err
	}
	if closeErr != nil {
		return closeErr
	}
	return syncDir(filepath.Dir(path))
}

// Sync each newly created directory in its parent before publishing a record.
func durableMkdir(path string) error {
	if _, err := os.Stat(path); err == nil {
		return nil
	} else if !os.IsNotExist(err) {
		return err
	}
	parent := filepath.Dir(path)
	if err := durableMkdir(parent); err != nil {
		return err
	}
	if err := os.Mkdir(path, 0700); err != nil && !os.IsExist(err) {
		return err
	}
	return syncDir(parent)
}

// Fault injection stays private to this package's persistence tests.
var writeFault = func(stage, path string) error { return nil }

func CommitBatch(ideaDir string, b Batch) error {
	h, err := ReadHistory(ideaDir)
	if err != nil {
		return err
	}
	if h == nil {
		return fmt.Errorf("missing quota history")
	}
	if b.PriorRevision < h.Revision {
		for _, old := range h.Batches {
			if old.ID == b.ID && b.ID == b.digest() {
				return syncDir(filepath.Join(ideaDir, HistoryDir))
			}
		}
		return fmt.Errorf("conflicting quota replay")
	}
	if err = b.Validate(h); err != nil {
		return err
	}
	if err = b.validateOwnerAuthority(ideaDir); err != nil {
		return err
	}
	raw, err := json.MarshalIndent(b, "", "  ")
	if err != nil {
		return err
	}
	return DurableWrite(filepath.Join(ideaDir, HistoryDir, fmt.Sprintf("%06d.json", b.PriorRevision+1)), append(raw, '\n'), true)
}
func (b Batch) Markers() []string {
	k := Kickoff{Transition: &Transition{ID: b.ID, RecordedAt: b.RecordedAt, Decision: b.Decision}}
	return k.Markers()
}
func (b Batch) Notice() string {
	if b.Owner != nil && b.Owner.Authority == nil {
		return fmt.Sprintf("---\nfrom: parley\nto: user\nidea: %s\nblocking: no\ntransition: %s\n---\n\nManual policy-off membership revision (not owner-confirmed authority). Current participants: %v. Retained obligations remain in force.\n", b.Idea, b.ID, b.Decision.After)
	}
	if b.Owner != nil {
		return fmt.Sprintf("---\nfrom: parley\nto: user\nidea: %s\nblocking: no\ntransition: %s\n---\n\nOwner-confirmed membership/policy revision. Current participants: %v. Policy: %+v.\n", b.Idea, b.ID, b.Decision.After, b.Policy)
	}
	k := Kickoff{Idea: b.Idea, Participants: b.Decision.After, Transition: &Transition{ID: b.ID, RecordedAt: b.RecordedAt, Decision: b.Decision}}
	return strings.Replace(k.Notice(), "phase: kickoff", "phase: "+b.Round, 1) + fmt.Sprintf("\nArithmetic: %d before - %d candidates = %d current; %d usable non-facilitators >= %d.\n", len(b.Decision.Before), len(CandidateIDs(b.Decision.Candidates)), len(b.Decision.After), b.Decision.UsableSurvivors, Floor)
}

// BindRetained freezes the filed obligations inside the same immutable batch.
func (b Batch) BindRetained(items []Obligation) Batch {
	b.Retained = items
	b.ID = b.digest()
	return b
}

func uniqueKeys(raw []byte) bool {
	d := json.NewDecoder(bytes.NewReader(raw))
	var value func() bool
	value = func() bool {
		t, err := d.Token()
		if err != nil {
			return false
		}
		delim, ok := t.(json.Delim)
		if !ok {
			return true
		}
		switch delim {
		case '{':
			seen := map[string]bool{}
			for d.More() {
				k, err := d.Token()
				key, ok := k.(string)
				if err != nil || !ok || seen[key] {
					return false
				}
				seen[key] = true
				if !value() {
					return false
				}
			}
			end, err := d.Token()
			return err == nil && end == json.Delim('}')
		case '[':
			for d.More() {
				if !value() {
					return false
				}
			}
			end, err := d.Token()
			return err == nil && end == json.Delim(']')
		}
		return false
	}
	if !value() {
		return false
	}
	_, err := d.Token()
	return err == io.EOF
}
