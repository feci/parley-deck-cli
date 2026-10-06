package quota

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strings"
	"time"
)

// Revision records either an explicit owner decision or a manual policy-off edit. Automatic batches never add members
// or alter policy. Catchup snapshots are canonical late round-1 artifacts, not
// replacement agents or a timer-based rejoin.
type Revision struct {
	Authority *Authority        `json:"authority,omitempty"`
	Catchup   map[string]string `json:"catchup,omitempty"`
	// ManualPrompt preserves an existing §9.0 knob-off confirmation without
	// requiring the new CLI. Protocol validates it before committing the snapshot.
	ManualPrompt string `json:"manual_prompt,omitempty"`
}

type RevisionDirective struct {
	Participants []string `json:"participants"`
	Policy       Policy   `json:"policy"`
}

func (d RevisionDirective) Text() string {
	b, _ := json.Marshal(d)
	return "Quota revision: " + string(b)
}
func (h *History) Policy() Policy {
	if h == nil || h.Kickoff == nil {
		return Policy{}
	}
	p := h.Kickoff.Policy
	for _, b := range h.Batches {
		if b.Owner != nil {
			p = b.Policy
		}
	}
	return p
}
func NewRevision(h *History, run string, ids []string, policy Policy, owner Revision, now time.Time) Batch {
	b := Batch{Version: 1, Idea: h.Kickoff.Idea, RunID: run, PriorRevision: h.Revision, Policy: policy, RecordedAt: now.UTC(), Owner: &owner, Decision: Decision{Before: append([]string(nil), h.Current...), After: append([]string(nil), ids...)}}
	b.ID = b.digest()
	return b
}
func (b Batch) validateRevision(h *History) error {
	if b.Version != 1 || b.Idea != h.Kickoff.Idea || b.RunID == "" || b.ID != b.digest() || b.PriorRevision != h.Revision || b.RecordedAt.IsZero() || b.Policy.Validate() != nil || !reflect.DeepEqual(b.Decision.Before, h.Current) || b.Decision.Applied || len(b.Decision.Candidates) != 0 || len(b.Decision.After) == 0 || len(Unique(b.Decision.After)) != len(b.Decision.After) {
		return fmt.Errorf("invalid owner membership revision")
	}
	for _, id := range b.Decision.After {
		if !identifier.MatchString(id) {
			return fmt.Errorf("invalid member identity")
		}
	}
	if b.Owner.Authority == nil {
		if h.Policy().Enabled || b.Policy != h.Policy() || b.Owner.ManualPrompt == "" {
			return fmt.Errorf("missing owner authority")
		}
		if err := validateManualSnapshot(b, h); err != nil {
			return err
		}
	} else if b.Owner.ManualPrompt != "" || !HasDirective(b.Owner.Authority.Quote, (RevisionDirective{b.Decision.After, b.Policy}).Text()) {
		return fmt.Errorf("owner quote does not authorize this membership/policy revision")
	}
	for _, id := range b.Decision.After {
		if !contains(h.Known, id) {
			if b.Owner.Authority == nil && ManualRoundOneReturn(b.Owner.ManualPrompt, id) {
				continue
			}
			validate := ValidateCatchupSnapshot
			if b.Owner.Authority == nil {
				validate = ValidateManualCatchupSnapshot
			}
			if err := validate(b.Owner.Catchup[id], b.Idea, id); err != nil {
				return fmt.Errorf("new identity %s needs a valid catch-up artifact: %w", id, err)
			}
		}
	}
	for id := range b.Owner.Catchup {
		if contains(h.Known, id) || !contains(b.Decision.After, id) {
			return fmt.Errorf("unrelated catch-up snapshot %s", id)
		}
	}
	return nil
}

func HasDirective(quote, directive string) bool {
	for _, line := range strings.Split(quote, "\n") {
		if strings.TrimSpace(line) == directive {
			return true
		}
	}
	return false
}

func (b Batch) validateOwnerAuthority(ideaDir string) error {
	if b.Owner == nil {
		return nil
	}
	if b.Owner.Authority != nil {
		return ValidateAuthority(IdeaRoot(ideaDir), b.Idea, *b.Owner.Authority)
	}
	return nil
}
