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
	RoleSnapshot *RoleSnapshot `json:"role_snapshot,omitempty"`
	Participants []string      `json:"participants"`
	Policy       Policy        `json:"policy"`
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
func NewRevision(h *History, run string, ids []string, policy Policy, owner Revision, now time.Time, roles ...*RoleSnapshot) Batch {
	b := Batch{Version: 1, Idea: h.Kickoff.Idea, RunID: run, PriorRevision: h.Revision, Policy: policy, RecordedAt: now.UTC(), Owner: &owner, Decision: Decision{Before: append([]string(nil), h.Current...), After: append([]string(nil), ids...)}}
	if len(roles) > 0 {
		b.RoleSnapshot = roles[0]
	}
	b.ID = b.digest()
	return b
}
func (b Batch) validateRevision(h *History) error {
	if err := h.CheckReturn(b.Decision.After); err != nil {
		return err
	}
	if b.Version != 1 || b.Idea != h.Kickoff.Idea || b.RunID == "" || b.ID != b.digest() || b.PriorRevision != h.Revision || b.RecordedAt.IsZero() || b.Policy.Validate() != nil || !reflect.DeepEqual(b.Decision.Before, h.Current) || b.Decision.Applied || len(b.Decision.Candidates) != 0 || len(b.Decision.After) == 0 || len(Unique(b.Decision.After)) != len(b.Decision.After) {
		return fmt.Errorf("invalid owner membership revision")
	}
	if b.Policy.rolesFrozen() && h.RoleSnapshot() == nil {
		if b.RoleSnapshot == nil || b.RoleSnapshot.Validate() != nil || !reflect.DeepEqual(b.RoleSnapshot.Order, h.Current) || b.Owner.Authority == nil {
			return fmt.Errorf("v2 widening requires owner-bound frozen role snapshot")
		}
	} else if b.RoleSnapshot != nil {
		return fmt.Errorf("frozen role snapshot cannot be replaced")
	}
	if b.Decision.RoleTransition != nil {
		return fmt.Errorf("owner policy revision cannot fabricate automatic succession")
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
	} else if b.Owner.ManualPrompt != "" || !HasDirective(b.Owner.Authority.Quote, (RevisionDirective{Participants: b.Decision.After, Policy: b.Policy, RoleSnapshot: b.RoleSnapshot}).Text()) {
		return fmt.Errorf("owner quote does not authorize this membership/policy revision")
	}
	for _, id := range b.Decision.After {
		if !contains(h.Known, id) {
			if b.Owner.Authority == nil && ManualRoundOneReturn(b.Owner.ManualPrompt) {
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
