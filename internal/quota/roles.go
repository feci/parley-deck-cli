package quota

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"path/filepath"
	"reflect"
	"strings"
)

// RoleSnapshot freezes identity and order. It does not assert availability.
type RoleSnapshot struct {
	Version     int      `json:"version"`
	Order       []string `json:"order"`
	Organizer   string   `json:"organizer,omitempty"`
	Facilitator string   `json:"facilitator,omitempty"`
}

type DraftArchive struct {
	Path    string `json:"path"`
	Author  string `json:"author"`
	SHA256  string `json:"sha256"`
	Archive string `json:"archive"`
}

type RoleAssignment struct {
	Organizer         string `json:"organizer,omitempty"`
	Implementer       string `json:"implementer,omitempty"`
	ImplementerSource string `json:"implementer_source,omitempty"`
}

// RoleTransition is part of the membership transaction, not a separate ledger.
type RoleTransition struct {
	Rule       string         `json:"rule"`
	ID         string         `json:"id"`
	Before     RoleAssignment `json:"before"`
	After      RoleAssignment `json:"after"`
	Successor  string         `json:"successor"`
	Roles      []string       `json:"roles"`
	Usable     []string       `json:"usable"`
	Phase      string         `json:"phase,omitempty"`
	Base       string         `json:"base"`
	BaseSHA256 string         `json:"base_sha256,omitempty"`
	Branch     string         `json:"branch,omitempty"`
	Head       string         `json:"head,omitempty"`
}

func (s RoleSnapshot) Validate() error {
	if s.Version != 1 || len(s.Order) == 0 || len(Unique(s.Order)) != len(s.Order) {
		return fmt.Errorf("invalid frozen role order")
	}
	for _, id := range append(append([]string{}, s.Order...), s.Organizer, s.Facilitator) {
		if id != "" && !identifier.MatchString(id) {
			return fmt.Errorf("invalid frozen role identity")
		}
	}
	if s.Organizer != "" && !contains(s.Order, s.Organizer) && s.Organizer != s.Facilitator {
		return fmt.Errorf("external organizer must be the frozen facilitator")
	}
	return nil
}

func (h *History) RoleSnapshot() *RoleSnapshot {
	if h == nil || h.Kickoff == nil {
		return nil
	}
	s := h.Kickoff.Roles
	for _, b := range h.Batches {
		if b.RoleSnapshot != nil {
			s = b.RoleSnapshot
		}
	}
	return s
}

func (h *History) LastRoleTransition() *RoleTransition {
	if h == nil {
		return nil
	}
	for i := len(h.Batches) - 1; i >= 0; i-- {
		if r := h.Batches[i].Decision.RoleTransition; r != nil {
			return r
		}
	}
	if h.Kickoff != nil && h.Kickoff.Transition != nil {
		return h.Kickoff.Transition.Decision.RoleTransition
	}
	return nil
}

func (r RoleTransition) digest() string {
	r.ID = ""
	b, _ := json.Marshal(r)
	return fmt.Sprintf("role-%x", sha256.Sum256(b))
}

func (r RoleTransition) Path() string { return "role-handover/" + r.ID + ".md" }

func (d Decision) BindHandover(phase, base, hash, branch, head string) Decision {
	if d.RoleTransition == nil {
		return d
	}
	r := *d.RoleTransition
	r.Phase, r.Base, r.BaseSHA256, r.Branch, r.Head = phase, base, hash, branch, head
	r.ID = r.digest()
	d.RoleTransition = &r
	return d
}

// SelectSuccession applies the owner order once to the complete loss set.
// It never searches past an unusable selected identity.
func SelectSuccession(s RoleSnapshot, before RoleAssignment, current []string, lost, usable []string) (*RoleTransition, error) {
	if err := s.Validate(); err != nil {
		return nil, err
	}
	organizerLost := before.Organizer != "" && contains(lost, before.Organizer)
	implementerLost := before.Implementer != "" && contains(lost, before.Implementer)
	if !organizerLost && !implementerLost {
		return nil, nil
	}
	after := before
	roles := []string{}
	successor := ""
	switch {
	case organizerLost && implementerLost:
		for _, id := range s.Order {
			if contains(current, id) && !contains(lost, id) {
				successor = id
				break
			}
		}
		after.Organizer, after.Implementer = successor, successor
		roles = []string{"organizer", "implementer"}
	case organizerLost:
		successor = before.Implementer
		after.Organizer = successor
		roles = []string{"organizer"}
	case implementerLost:
		successor = before.Organizer
		after.Implementer = successor
		roles = []string{"implementer"}
	}
	if successor == "" || contains(lost, successor) || !contains(usable, successor) {
		return nil, fmt.Errorf("selected successor %q lacks positive usability; no substitute is authorized", successor)
	}
	if !contains(current, successor) && successor != s.Facilitator {
		return nil, fmt.Errorf("selected successor is outside frozen membership")
	}
	if before.Implementer == "" || before.ImplementerSource == "" {
		return nil, fmt.Errorf("succession requires a resolved implementer and source")
	}
	if implementerLost {
		after.ImplementerSource = "role-transition"
	}
	r := &RoleTransition{Rule: "role_transition.v1", Before: before, After: after, Successor: successor,
		Roles: roles, Usable: Unique(usable), Base: "none"}
	r.ID = r.digest()
	return r, nil
}

func validateRoleDecision(s *RoleSnapshot, previous *RoleTransition, d Decision) error {
	r := d.RoleTransition
	if r == nil {
		if s != nil && previous == nil && contains(CandidateIDs(d.Candidates), s.Organizer) {
			return fmt.Errorf("frozen organizer lost without succession")
		}
		if previous != nil {
			for _, id := range CandidateIDs(d.Candidates) {
				if id == previous.After.Organizer || id == previous.After.Implementer {
					return fmt.Errorf("role holder lost without succession")
				}
			}
		}
		return nil
	}
	if s == nil || r.Rule != "role_transition.v1" || r.ID != r.digest() || len(r.Usable) != len(Unique(r.Usable)) {
		return fmt.Errorf("invalid role transition authority")
	}
	for _, id := range []string{r.Before.Organizer, r.Before.Implementer, r.Successor} {
		if id != "" && (!identifier.MatchString(id) || (!contains(d.Before, id) && id != s.Facilitator)) {
			return fmt.Errorf("role holder is outside frozen authority")
		}
	}
	if r.Before.ImplementerSource == "" || r.Before.Implementer == "" {
		return fmt.Errorf("missing implementer resolution")
	}
	positive := 0
	for _, id := range r.Usable {
		if (!contains(d.Before, id) && id != s.Facilitator) || contains(CandidateIDs(d.Candidates), id) {
			return fmt.Errorf("invalid positive usability binding")
		}
		if id != s.Facilitator {
			positive++
		}
	}
	if positive < Floor || positive != d.UsableSurvivors {
		return fmt.Errorf("role transition lacks the positive floor")
	}
	if previous != nil && r.Before != previous.After || previous == nil && r.Before.Organizer != s.Organizer {
		return fmt.Errorf("role transition contradicts prior holders")
	}
	for _, c := range d.Candidates {
		if c.Agent == r.Before.Organizer || c.Agent == r.Before.Implementer {
			if c.Evidence.RuleID != ProviderExhaustionRule {
				return fmt.Errorf("protected succession requires provider exhaustion")
			}
		}
	}
	want, err := SelectSuccession(*s, r.Before, d.Before, CandidateIDs(d.Candidates), r.Usable)
	if err != nil || want == nil || want.After != r.After || want.Successor != r.Successor || !reflect.DeepEqual(want.Roles, r.Roles) {
		return fmt.Errorf("role transition violates owner order")
	}
	if r.Base == "none" {
		if previous != nil || r.BaseSHA256 != "" {
			return fmt.Errorf("role transition has a missing predecessor")
		}
	} else if !proofHash.MatchString(r.BaseSHA256) || filepath.IsAbs(r.Base) || filepath.ToSlash(filepath.Clean(r.Base)) != r.Base || strings.HasPrefix(r.Base, "../") || (previous == nil && r.Base != "IMPLEMENTATION.md") || (previous != nil && r.Base != previous.Path()) {
		return fmt.Errorf("invalid handover base binding")
	}
	return nil
}
