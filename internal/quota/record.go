package quota

import (
	"crypto/sha256"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"parley-deck-cli/internal/fsutil"
	"path/filepath"
	"reflect"
	"time"
)

const KickoffFile = "quota-kickoff.json"

// Kickoff is immutable revision zero, including the policy frozen for all later
// runs. Participants is the first quorum: kickoff-excluded identities are not
// historical members. Before in Transition is only the proposed readiness set.
type Kickoff struct {
	Version      int         `json:"version"`
	Idea         string      `json:"idea"`
	RunID        string      `json:"run_id"`
	Revision     int         `json:"revision"`
	Policy       Policy      `json:"policy"`
	Participants []string    `json:"participants"`
	Transition   *Transition `json:"transition,omitempty"`
}
type Transition struct {
	ID            string    `json:"id"`
	PriorRevision int       `json:"prior_revision"`
	Policy        Policy    `json:"policy"`
	Decision      Decision  `json:"decision"`
	RecordedAt    time.Time `json:"recorded_at"`
}

func NewKickoff(idea, run string, p Policy, participants []string, d *Decision, now time.Time) Kickoff {
	k := Kickoff{Version: 1, Idea: idea, RunID: run, Policy: p, Participants: append([]string(nil), participants...)}
	if d != nil && d.Applied {
		k.Transition = &Transition{PriorRevision: -1, Policy: p, Decision: *d, RecordedAt: now.UTC()}
		payload, _ := json.Marshal(struct {
			Idea, Run string
			Decision  Decision
		}{idea, run, *d})
		k.Transition.ID = fmt.Sprintf("kickoff-%x", sha256.Sum256(payload))
	}
	return k
}
func (k Kickoff) Validate() error {
	if k.Version != 1 || k.Idea == "" || k.RunID == "" || k.Revision != 0 || len(k.Participants) == 0 {
		return fmt.Errorf("invalid quota kickoff history")
	}
	if err := k.Policy.Validate(); err != nil {
		return err
	}
	if len(Unique(k.Participants)) != len(k.Participants) {
		return fmt.Errorf("duplicate kickoff membership")
	}
	if tr := k.Transition; tr != nil {
		if tr.ID == "" || tr.PriorRevision != -1 || tr.RecordedAt.IsZero() || !tr.Decision.Applied || !k.Policy.Enabled || tr.Policy != k.Policy || !reflect.DeepEqual(Unique(tr.Decision.After), Unique(k.Participants)) {
			return fmt.Errorf("contradictory quota kickoff transition")
		}
		if tr.Decision.UsableSurvivors < Floor || len(tr.Decision.Candidates) == 0 {
			return fmt.Errorf("invalid quota kickoff floor or evidence")
		}
		readinessBatch := ""
		for _, c := range tr.Decision.Candidates {
			if k.Policy.Dropout() || c.Evidence.RuleID == ParticipantFailureRule {
				// A readiness batch is bound before the final idea slug exists.
				if !k.Policy.Dropout() || contains(k.Participants, c.Agent) || ValidateFailureEvidence(c.Evidence, "", c.Agent) != nil {
					return fmt.Errorf("invalid kickoff participant failure evidence")
				}
				if c.Evidence.Failure.Step != "readiness" || !contains(tr.Decision.Before, c.Agent) {
					return fmt.Errorf("kickoff failure is not a proposed readiness step")
				}
				if readinessBatch != "" && readinessBatch != c.Evidence.Failure.Idea {
					return fmt.Errorf("mixed kickoff readiness batches")
				}
				readinessBatch = c.Evidence.Failure.Idea
			}
		}
	}
	return nil
}

// WriteKickoff uses exclusive creation and checks file and directory durability.
// A crash-truncated record remains visible and unreadable: it cannot turn into a
// legacy policy or be silently overwritten on a later run.
func WriteKickoff(ideaDir string, k Kickoff) error {
	if err := k.Validate(); err != nil {
		return err
	}
	b, err := json.MarshalIndent(k, "", "  ")
	if err != nil {
		return err
	}
	b = append(b, '\n')
	path := filepath.Join(ideaDir, KickoffFile)
	f, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_EXCL, 0600)
	if errors.Is(err, os.ErrExist) {
		old, e := ReadKickoff(ideaDir)
		if e == nil && reflect.DeepEqual(old, &k) {
			return nil
		}
		return fmt.Errorf("quota kickoff history already exists or is corrupt")
	}
	if err != nil {
		return err
	}
	if _, err = f.Write(b); err != nil {
		f.Close()
		return err
	}
	if err = fsutil.SyncFile(f); err != nil {
		f.Close()
		return err
	}
	if err = f.Close(); err != nil {
		return err
	}
	dir, err := os.Open(ideaDir)
	if err != nil {
		return err
	}
	defer dir.Close()
	return fsutil.SyncFile(dir)
}
func ReadKickoff(ideaDir string) (*Kickoff, error) {
	b, err := os.ReadFile(filepath.Join(ideaDir, KickoffFile))
	if errors.Is(err, os.ErrNotExist) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	var k Kickoff
	if err = strictDecode(b, &k); err != nil {
		return nil, fmt.Errorf("invalid quota history: %w", err)
	}
	if err = k.Validate(); err != nil {
		return nil, err
	}
	return &k, nil
}
func (k Kickoff) Markers() []string {
	if k.Transition == nil {
		return nil
	}
	out := []string{}
	for _, c := range k.Transition.Decision.Candidates {
		e := c.Evidence
		label := "automatic quota"
		if e.RuleID == ParticipantFailureRule {
			label = "permanent participant dropout"
		}
		out = append(out, fmt.Sprintf("%s — %s — rule=%s — reset=%s — raw_reset=%q — transition=%s — recorded %s", c.Agent, label, e.RuleID, e.ResetHint(), e.RawReset, k.Transition.ID, k.Transition.RecordedAt.Format("2006-01-02")))
	}
	return out
}
func (k Kickoff) Notice() string {
	if k.Transition == nil {
		return ""
	}
	tr := k.Transition
	body := fmt.Sprintf("---\nfrom: parley\nto: user\nidea: %s\nphase: kickoff\nblocking: no\ntransition: %s\ndate: %s\n---\n\nAutomatic quota exclusion for this idea only.\n\n", k.Idea, tr.ID, tr.RecordedAt.Format("2006-01-02"))
	for _, c := range tr.Decision.Candidates {
		if c.Evidence.RuleID == ParticipantFailureRule {
			body += "Permanent participant dropout for this idea; no same-idea return, even after opt-out or policy revision. The next idea probes afresh. Retained vetoes, disputes and findings remain in force.\n\n"
		}
		body += fmt.Sprintf("- %s: %s (%s), reset %s. %s\n", c.Agent, c.Evidence.Excerpt, c.Evidence.RuleID, c.Evidence.ResetHint(), c.Evidence.RelaunchHint())
	}
	return body + fmt.Sprintf("\nSurvivors: %v. Remaining gates: reviewer count for the selected track; LE-7/LE-11 two-reviewer close (auto_implement); independent goal-checker eligibility (auto_implement or strict_gate); require_model_diversity (when enabled, and fast track); strict_gate (when enabled); retained vetoes, DISPUTED claims and findings.\n", k.Participants)
}
