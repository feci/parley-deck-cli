// Package quota holds shared policy and evidence, independent of dispatch and adapters.
// Exclusion markers are display records; they are never membership authority.
package quota

import (
	"encoding/json"
	"fmt"
	"sort"
	"strings"
	"time"
)

const Floor = 2
const MinimumReset = time.Hour
const KickoffOnly = "kickoff-only"
const KickoffAndMidIdea = "kickoff-and-mid-idea"
const ParticipantFailure = "participant-failure-v1"
const ParticipantFailureRule = "participant-failure.v1"

type Policy struct {
	Enabled bool   `json:"enabled"`
	Scope   string `json:"scope"`
	Trigger string `json:"trigger,omitempty"`
}

// Presence is authority: an omitted false value cannot silently disable a saved policy.
func (p *Policy) UnmarshalJSON(raw []byte) error {
	var fields map[string]json.RawMessage
	if err := json.Unmarshal(raw, &fields); err != nil {
		return err
	}
	if !uniqueKeys(raw) || (len(fields) != 2 && len(fields) != 3) || fields["enabled"] == nil || fields["scope"] == nil {
		return fmt.Errorf("incomplete quota policy")
	}
	var enabled *bool
	var scope *string
	if err := json.Unmarshal(fields["enabled"], &enabled); err != nil || enabled == nil {
		return fmt.Errorf("invalid quota policy boolean")
	}
	if err := json.Unmarshal(fields["scope"], &scope); err != nil || scope == nil {
		return fmt.Errorf("invalid quota policy scope")
	}
	p.Enabled, p.Scope = *enabled, *scope
	p.Trigger = ""
	if len(fields) == 3 {
		var trigger *string
		if fields["trigger"] == nil || json.Unmarshal(fields["trigger"], &trigger) != nil || trigger == nil || *trigger != ParticipantFailure {
			return fmt.Errorf("invalid automatic exclusion trigger")
		}
		p.Trigger = *trigger
	}
	return p.Validate()
}

// NewPolicy retains the legacy constructor for saved quota-only callers.
func NewPolicy(defaultValue, ideaValue *bool) Policy {
	enabled := true
	if defaultValue != nil {
		enabled = *defaultValue
	}
	if ideaValue != nil {
		enabled = *ideaValue
	}
	return Policy{Enabled: enabled, Scope: KickoffAndMidIdea}
}

// NewParticipantPolicy is used only when creating a new idea, never on resume.
func NewParticipantPolicy(defaultValue, ideaValue *bool) Policy {
	p := NewPolicy(defaultValue, ideaValue)
	p.Trigger = ParticipantFailure
	return p
}
func (p Policy) Dropout() bool { return p.Enabled && p.Trigger == ParticipantFailure }
func (p Policy) Validate() error {
	if p.Scope != KickoffOnly && p.Scope != KickoffAndMidIdea {
		return fmt.Errorf("ambiguous quota policy scope %q", p.Scope)
	}
	if p.Trigger != "" && p.Trigger != ParticipantFailure {
		return fmt.Errorf("unknown automatic exclusion trigger %q", p.Trigger)
	}
	return nil
}

// Evidence contains only scrubbed, normalized observations, never full logs or environments.
// Eligible is meaningful only after invocation, artifact and batch-success gates pass.
type Evidence struct {
	InvocationID string      `json:"invocation_id"`
	Adapter      string      `json:"adapter"`
	RuleID       string      `json:"rule_id,omitempty"`
	Provenance   string      `json:"provenance,omitempty"`
	Eligible     bool        `json:"eligible"`
	Reason       string      `json:"reason"`
	Excerpt      string      `json:"excerpt,omitempty"`
	RawReset     string      `json:"raw_reset,omitempty"`
	ObservedAt   time.Time   `json:"observed_at"`
	ResetAt      *time.Time  `json:"reset_at,omitempty"`
	Failure      *FailedStep `json:"failure,omitempty"`
}

func (e Evidence) ResetHint() string {
	if e.ResetAt == nil {
		return "unknown"
	}
	return e.ResetAt.UTC().Format(time.RFC3339Nano)
}
func (e Evidence) RelaunchHint() string {
	if e.ResetAt == nil || e.RuleID == ParticipantFailureRule {
		return ""
	}
	return "owner-authorized relaunch after " + e.ResetAt.Add(5*time.Minute).UTC().Format(time.RFC3339Nano) + " (provider estimate)"
}

type Member struct {
	ID            string
	Usable        bool
	ValidArtifact bool
	LaterSuccess  bool
	Evidence      *Evidence
}
type Roles struct {
	Facilitator       string
	Designee          string
	PinnedImplementer string
	StartedDrafters   []string
	// Global defaults intentionally do not protect a participant before a pin.
}

func (r Roles) Protected(id string) bool {
	return id != "" && (id == r.Facilitator || id == r.Designee || id == r.PinnedImplementer || contains(r.StartedDrafters, id))
}

type Candidate struct {
	Agent    string   `json:"agent"`
	Evidence Evidence `json:"evidence"`
}
type Decision struct {
	Before          []string    `json:"before"`
	After           []string    `json:"after"`
	Candidates      []Candidate `json:"candidates,omitempty"`
	UsableSurvivors int         `json:"usable_non_facilitators"`
	Applied         bool        `json:"applied"`
	Block           string      `json:"block,omitempty"`
}

// Evaluate settles the WHOLE batch once. Sorting makes evidence and decisions stable
// across probe completion order; duplicate identities cannot inflate the floor.
func Evaluate(policy Policy, proposed []string, members []Member, roles Roles) Decision {
	d := Decision{Before: FilterConfirmed(proposed, nil)}
	d.After = append([]string(nil), d.Before...)
	if err := policy.Validate(); err != nil {
		d.Block = err.Error()
		return d
	}
	if !policy.Enabled {
		return d
	}
	byID := map[string][]Member{}
	usable := map[string]bool{}
	for _, m := range members {
		byID[m.ID] = append(byID[m.ID], m)
	}
	for _, id := range d.Before {
		observations := byID[id]
		success, unresolved := false, len(observations) == 0
		var candidates []Candidate
		for _, m := range observations {
			if m.Usable || m.ValidArtifact || m.LaterSuccess {
				success = true
			}
			if m.Evidence != nil && m.Evidence.Eligible && m.Evidence.InvocationID != "" && m.Evidence.Provenance != "" && m.Evidence.RuleID != "" && (!policy.Dropout() || ValidateFailureEvidence(*m.Evidence, "", id) == nil) {
				candidates = append(candidates, Candidate{Agent: id, Evidence: *m.Evidence})
			} else if !m.Usable && !m.ValidArtifact && !m.LaterSuccess {
				unresolved = true
			}
		}
		if success {
			usable[id] = true
			if id != roles.Facilitator {
				d.UsableSurvivors++
			}
			continue
		}
		if len(candidates) > 0 && !unresolved {
			d.Candidates = append(d.Candidates, candidates...)
		}
	}
	sort.Slice(d.Candidates, func(i, j int) bool {
		a, b := d.Candidates[i], d.Candidates[j]
		if a.Agent != b.Agent {
			return a.Agent < b.Agent
		}
		return a.Evidence.InvocationID < b.Evidence.InvocationID
	})
	if policy.Dropout() {
		for _, id := range Unique(d.Before) {
			if roles.Protected(id) && len(byID[id]) > 0 && !usable[id] {
				d.Block = fmt.Sprintf("protected role %s failed or has unresolved usability; re-designate, record implementer_waived, or set implementer: none with owner confirmation; no exclusions applied", id)
				return d
			}
		}
	}
	if len(d.Candidates) == 0 {
		return d
	}
	if policy.Dropout() {
		for _, id := range []string{roles.Designee, roles.PinnedImplementer} {
			if id != "" && !usable[id] {
				d.Block = fmt.Sprintf("protected implementer %s lacks positive usability evidence; no exclusions applied", id)
				return d
			}
		}
	}
	protected := map[string]bool{roles.Facilitator: true, roles.Designee: true, roles.PinnedImplementer: true}
	for _, id := range roles.StartedDrafters {
		protected[id] = true
	}
	for _, c := range d.Candidates {
		if protected[c.Agent] {
			d.Block = fmt.Sprintf("protected role %s; re-designate, record implementer_waived, or set implementer: none with owner confirmation", c.Agent)
			return d
		}
	}
	if d.UsableSurvivors < Floor {
		d.Block = fmt.Sprintf("quota floor: %d proposed distinct participants, %d candidate identities, %d usable non-facilitator survivors < %d", len(d.Before), len(CandidateIDs(d.Candidates)), d.UsableSurvivors, Floor)
		return d
	}
	excluded := map[string]bool{}
	for _, c := range d.Candidates {
		excluded[c.Agent] = true
	}
	d.After = nil
	for _, id := range d.Before {
		if !excluded[id] {
			d.After = append(d.After, id)
		}
	}
	d.Applied = true
	return d
}
func Unique(ids []string) []string {
	seen := map[string]bool{}
	out := []string{}
	for _, id := range ids {
		if id != "" && !seen[id] {
			seen[id] = true
			out = append(out, id)
		}
	}
	sort.Strings(out)
	return out
}
func CandidateIDs(cs []Candidate) []string {
	ids := []string{}
	for _, c := range cs {
		ids = append(ids, c.Agent)
	}
	return Unique(ids)
}

// FilterConfirmed is used ONCE at creation for the explicit preflight decision,
// never to reconstruct later membership from prompt excluded: display lines.
func FilterConfirmed(proposed, records []string) []string {
	excluded := map[string]bool{}
	for _, record := range records {
		id, _, ok := strings.Cut(record, " — ")
		if ok {
			excluded[strings.TrimSpace(id)] = true
		}
	}
	out := []string{}
	seen := map[string]bool{}
	for _, id := range proposed {
		if !excluded[id] && !seen[id] {
			out = append(out, id)
			seen[id] = true
		}
	}
	return out
}
