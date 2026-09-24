package protocolpacket

// The applicability map: parley-deck/meta/packet-applicability.yaml.
//
// It is a human-authored protocol artifact — changing it is a §7 protocol change — and no
// generator can prove its semantic correctness. What the code CAN prove is structural: every
// source block is classified exactly once, every conditional entry names its trigger, every
// dependency resolves, and the ratified never-cut list is honored. That is what Check does.

import (
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"

	"gopkg.in/yaml.v3"
)

// MapSchema is the only accepted map schema.
const MapSchema = "parley.packet-applicability/v1"

// Include kinds.
const (
	IncludeAlways   = "always"
	IncludeWhen     = "when"
	IncludeOnDemand = "on-demand"
)

// ErrNoMap reports that the deck has no applicability map. Callers fall back to full context.
var ErrNoMap = errors.New("protocolpacket: no applicability map")

// Map is a parsed, structurally validated applicability map.
type Map struct {
	Schema string
	Blocks []Rule
	// Audiences (lean-organizer C) is the optional per-audience omission overlay.
	// The audience dimension may only OMIT blocks; it can never include a block the
	// phase/track/transport rules cut, and never below the never-cut floor.
	Audiences map[string]AudienceSpec
	index     map[string]int
}

// AudienceSpec is one audience's omission set: whole blocks, with the trigger that
// says when to read the full source instead.
type AudienceSpec struct {
	Omit []OmitRule `yaml:"omit,omitempty" json:"omit,omitempty"`
}

// OmitRule removes one block for this audience, carrying the omission-index trigger.
type OmitRule struct {
	Locator string `yaml:"locator" json:"locator"`
	Trigger string `yaml:"trigger" json:"trigger"`
}

// KnownAudiences is the closed audience vocabulary. `participant` is the default and
// carries no omissions; `facilitator` is the pure-organizer view.
var KnownAudiences = []string{"participant", "facilitator"}

// Rule classifies one source block.
type Rule struct {
	Locator    string   `yaml:"locator" json:"locator"`
	Include    string   `yaml:"include" json:"include"`
	Phases     []int    `yaml:"phases,omitempty" json:"phases,omitempty"`
	Tracks     []string `yaml:"tracks,omitempty" json:"tracks,omitempty"`
	Transports []string `yaml:"transports,omitempty" json:"transports,omitempty"`
	Flags      []string `yaml:"flags,omitempty" json:"flags,omitempty"`
	Requires   []string `yaml:"requires,omitempty" json:"requires,omitempty"`
	Trigger    string   `yaml:"trigger,omitempty" json:"trigger,omitempty"`
}

type mapDoc struct {
	Schema    string                  `yaml:"schema"`
	Source    string                  `yaml:"source"`
	Blocks    []Rule                  `yaml:"blocks"`
	Audiences map[string]AudienceSpec `yaml:"audiences,omitempty"`
}

// ParseMap validates a map document. Unknown keys, a second document, an unknown include
// kind, a conditional rule without a trigger or without any condition, and duplicate
// locators are refused rather than ignored.
func ParseMap(raw string) (*Map, error) {
	if strings.TrimSpace(raw) == "" {
		return nil, errors.New("applicability map: file is empty")
	}
	var doc mapDoc
	dec := yaml.NewDecoder(strings.NewReader(raw))
	dec.KnownFields(true)
	if err := dec.Decode(&doc); err != nil {
		return nil, fmt.Errorf("applicability map: %w", err)
	}
	var extra yaml.Node
	if err := dec.Decode(&extra); err == nil {
		return nil, errors.New("applicability map: file contains more than one YAML document")
	}
	if doc.Schema != MapSchema {
		return nil, fmt.Errorf("applicability map: schema is %q, want %q", doc.Schema, MapSchema)
	}
	m := &Map{Schema: doc.Schema, index: map[string]int{}}
	for i, r := range doc.Blocks {
		r.Locator = strings.Join(strings.Fields(r.Locator), " ")
		if r.Locator == "" {
			return nil, fmt.Errorf("applicability map: entry %d has an empty locator", i+1)
		}
		if _, dup := m.index[r.Locator]; dup {
			return nil, fmt.Errorf("applicability map: duplicate locator %q", r.Locator)
		}
		hasCond := len(r.Phases) > 0 || len(r.Tracks) > 0 || len(r.Transports) > 0 || len(r.Flags) > 0
		switch r.Include {
		case IncludeAlways:
			if hasCond {
				return nil, fmt.Errorf("applicability map: %q is `always` but carries conditions", r.Locator)
			}
		case IncludeWhen:
			if !hasCond {
				return nil, fmt.Errorf("applicability map: %q is `when` but names no phase/track/transport/flag", r.Locator)
			}
			if strings.TrimSpace(r.Trigger) == "" {
				return nil, fmt.Errorf("applicability map: %q is conditional but has no trigger", r.Locator)
			}
		case IncludeOnDemand:
			if hasCond {
				return nil, fmt.Errorf("applicability map: %q is `on-demand` but carries conditions", r.Locator)
			}
			if strings.TrimSpace(r.Trigger) == "" {
				return nil, fmt.Errorf("applicability map: %q is on-demand but has no trigger", r.Locator)
			}
		default:
			return nil, fmt.Errorf("applicability map: %q has include %q; want always|when|on-demand", r.Locator, r.Include)
		}
		for _, p := range r.Phases {
			if p < 0 || p > 8 {
				return nil, fmt.Errorf("applicability map: %q names phase %d; want 0..8", r.Locator, p)
			}
		}
		for _, t := range r.Tracks {
			if !contains(KnownTracks, t) {
				return nil, fmt.Errorf("applicability map: %q names unknown track %q", r.Locator, t)
			}
		}
		for _, t := range r.Transports {
			if !contains(KnownTransports, t) {
				return nil, fmt.Errorf("applicability map: %q names unknown transport %q", r.Locator, t)
			}
		}
		for _, f := range r.Flags {
			if !contains(KnownFlags, f) {
				return nil, fmt.Errorf("applicability map: %q names unknown flag %q", r.Locator, f)
			}
		}
		m.index[r.Locator] = len(m.Blocks)
		m.Blocks = append(m.Blocks, r)
	}
	if len(m.Blocks) == 0 {
		return nil, errors.New("applicability map: declares no blocks")
	}
	for name, aud := range doc.Audiences {
		if !contains(KnownAudiences, name) {
			return nil, fmt.Errorf("applicability map: unknown audience %q; want one of %s", name, strings.Join(KnownAudiences, "|"))
		}
		seenOmit := map[string]bool{}
		for i, or := range aud.Omit {
			or.Locator = strings.Join(strings.Fields(or.Locator), " ")
			if or.Locator == "" {
				return nil, fmt.Errorf("applicability map: audiences.%s omit entry %d has an empty locator", name, i+1)
			}
			if seenOmit[or.Locator] {
				return nil, fmt.Errorf("applicability map: audiences.%s omits %q twice", name, or.Locator)
			}
			seenOmit[or.Locator] = true
			if strings.TrimSpace(or.Trigger) == "" {
				return nil, fmt.Errorf("applicability map: audiences.%s omits %q without a trigger naming when to read the full source", name, or.Locator)
			}
			aud.Omit[i] = or
		}
	}
	m.Audiences = doc.Audiences
	return m, nil
}

// audienceOmit returns the omission trigger for a locator under the named audience.
func (m *Map) audienceOmit(audience, locator string) (string, bool) {
	if m == nil || m.Audiences == nil {
		return "", false
	}
	spec, ok := m.Audiences[audience]
	if !ok {
		return "", false
	}
	for _, or := range spec.Omit {
		if or.Locator == locator {
			return or.Trigger, true
		}
	}
	return "", false
}

// neverCutForRequest reports whether a locator is pinned by the ratified never-cut
// floor FOR THIS REQUEST (always; its flag set; its transport active; its phase
// listed). The audience pass may not omit such a block — this is the hard safety
// property that the audience dimension can never cut below the floor.
func neverCutForRequest(locator string, req Request) bool {
	for _, nc := range neverCut {
		if !strings.HasPrefix(locator, nc.prefix) {
			continue
		}
		switch {
		case nc.always:
			return true
		case nc.flag != "":
			if contains(req.Flags, nc.flag) {
				return true
			}
		case nc.transport != "":
			if req.Transport == nc.transport {
				return true
			}
		case len(nc.phases) > 0:
			for _, p := range nc.phases {
				if p == req.Phase {
					return true
				}
			}
		}
	}
	return false
}

// LoadMap reads and parses the map at path. A missing file is ErrNoMap.
func LoadMap(path string) (*Map, error) {
	b, err := os.ReadFile(path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return nil, ErrNoMap
		}
		return nil, err
	}
	return ParseMap(string(b))
}

func (m *Map) rule(locator string) (Rule, bool) {
	i, ok := m.index[locator]
	if !ok {
		return Rule{}, false
	}
	return m.Blocks[i], true
}

// applies reports why a conditional rule applies to the request, or "" when it does not.
func (r Rule) applies(req Request) string {
	for _, p := range r.Phases {
		if p == req.Phase {
			return "phase " + strconv.Itoa(p)
		}
	}
	if contains(r.Tracks, req.Track) {
		return "track " + req.Track
	}
	if contains(r.Transports, req.Transport) {
		return "transport " + req.Transport
	}
	for _, f := range req.Flags {
		if contains(r.Flags, f) {
			return "flag " + f
		}
	}
	return ""
}

// neverCutRule is one ratified invariant the map must honor
// (phase-packet FINAL.md "Never cut"; evidence-first FINAL.md "Packet Never-Cut").
type neverCutRule struct {
	prefix    string
	always    bool
	phases    []int
	flag      string
	transport string
}

var kernel = []int{1, 2, 3, 5, 6, 7, 8}
var fullFifteen = []int{1, 2, 3, 6, 7}

var neverCut = []neverCutRule{
	{prefix: "## 0.", always: true},                              // universal invariants: files canonical, non-solo
	{prefix: "### Non-solo execution requirement", always: true}, // non-solo close guard
	{prefix: "### 4.0 —", always: true},                          // §4.0 overrides and invariants
	{prefix: "### 4.0.1", always: true},
	{prefix: "### Escalation to user", always: true},
	{prefix: "## 6.", always: true}, // rule 3, status re-read, English-only, no-secrets
	{prefix: "## 14.", always: true},
	{prefix: "### 14.1", always: true},
	{prefix: "### 14.2", always: true},
	{prefix: "### 14.3", always: true},
	{prefix: "## 11.", always: true},
	{prefix: "### 11.A", transport: "local-dir"},
	{prefix: "### 11.B", transport: "github-pr"},
	{prefix: "### 11.C", transport: "gitlab-mr"},
	{prefix: "## 7.", flag: "protocol_change"},
	{prefix: "## 15.", phases: kernel},
	{prefix: "### 15.1", phases: kernel},
	{prefix: "### 15.2", phases: kernel},
	{prefix: "### 15.3", phases: kernel},
	{prefix: "### 15.4", phases: kernel},
	{prefix: "### 15.7", phases: kernel},
	{prefix: "### 15.5", phases: fullFifteen},
	{prefix: "### 15.6", phases: fullFifteen},
	{prefix: "### Phase 0 —", phases: []int{0}},
	{prefix: "### Phase 1 —", phases: []int{1}},
	{prefix: "### Phase 2 —", phases: []int{2}},
	{prefix: "### Phase 3 —", phases: []int{3}},
	{prefix: "### Phase 4 —", phases: []int{4}},
	{prefix: "### Phase 5 —", phases: []int{5}},
	{prefix: "### Phase 6 —", phases: []int{6}},
	{prefix: "### Phase 7 —", phases: []int{7}},
	{prefix: "### Phase 8 —", phases: []int{8}},
	{prefix: "#### Strict review gate", phases: []int{8}},
	{prefix: "#### Stopping judgment", phases: []int{8}},
}

// Report is the result of `packet check`.
type Report struct {
	OK           bool     `json:"ok"`
	SourceSHA256 string   `json:"source_sha256"`
	Blocks       int      `json:"blocks"`
	Unclassified []string `json:"unclassified,omitempty"`
	Stale        []string `json:"stale,omitempty"`
	NeverCut     []string `json:"never_cut_violations,omitempty"`
	Fallbacks    []string `json:"fallbacks,omitempty"`
	Secret       string   `json:"secret,omitempty"`
}

// Check verifies map completeness against the live source, the ratified never-cut list, and
// that an optimized packet renders for every phase and track on the source's transport. It
// proves structure, not semantics: a wrongly classified block that is present in the map is
// invisible here and is why the default stays full/shadow.
func Check(src Source, m *Map) Report {
	rep := Report{SourceSHA256: Hash(src.Raw)}
	if name := DetectSecret(src.Raw); name != "" {
		rep.Secret = name
	}
	blocks := Parse(src.Raw)
	rep.Blocks = len(blocks)
	have := map[string]bool{}
	for _, b := range blocks {
		have[b.Locator] = true
		if b.Locator == PreambleLocator {
			continue
		}
		if m == nil {
			rep.Unclassified = append(rep.Unclassified, b.Locator)
			continue
		}
		if _, ok := m.rule(b.Locator); !ok {
			rep.Unclassified = append(rep.Unclassified, b.Locator)
		}
	}
	if m != nil {
		for _, r := range m.Blocks {
			if !have[r.Locator] {
				rep.Stale = append(rep.Stale, r.Locator)
			}
			for _, dep := range r.Requires {
				if !have[dep] {
					rep.Stale = append(rep.Stale, dep+" (required by "+r.Locator+")")
				}
			}
		}
		for _, b := range blocks {
			for _, nc := range neverCut {
				if !strings.HasPrefix(b.Locator, nc.prefix) {
					continue
				}
				r, ok := m.rule(b.Locator)
				if !ok {
					continue // already reported as unclassified
				}
				if v := nc.violation(r); v != "" {
					rep.NeverCut = append(rep.NeverCut, b.Locator+": "+v)
				}
			}
		}
		for name, aud := range m.Audiences {
			for _, or := range aud.Omit {
				if !have[or.Locator] {
					rep.Stale = append(rep.Stale, "audiences."+name+" omits unknown block "+or.Locator)
					continue
				}
				// An audience may never name an UNCONDITIONAL never-cut block
				// (pinned at every request) nor a PHASE- or FLAG-PINNED one
				// (§15.x at the kernel phases, `### Phase N`, §7 under
				// protocol_change): the omission rule is a map-level floor breach
				// even though Build would refuse to apply it — `packet check` must
				// PROVE the floor for exactly the sections FINAL C.1 singles out
				// (fix-up F6, claude-1 MAJ-6: the phase-pinned §15 family was
				// previously allowed here, so the negative test could not fire for
				// it). The allowance is narrowed to TRANSPORT-conditional blocks —
				// the non-active §11 subsections the ratified facilitator set
				// legitimately omits.
				for _, nc := range neverCut {
					if nc.transport != "" {
						continue
					}
					if strings.HasPrefix(or.Locator, nc.prefix) {
						rep.NeverCut = appendUnique(rep.NeverCut, "audiences."+name+" omits never-cut block "+or.Locator)
						break
					}
				}
			}
		}
		transport := src.Transport
		if transport == "" {
			transport = "local-dir"
		}
		for phase := 0; phase <= 8; phase++ {
			for _, track := range KnownTracks {
				c := Build(src, m, Request{Phase: phase, Track: track, Transport: transport, Optimize: true})
				if c.ContextMode != ModePacket {
					rep.Fallbacks = appendUnique(rep.Fallbacks, fmt.Sprintf("phase %d/%s: %s (%s)", phase, track, c.ContextMode, c.FallbackReason))
				}
				// Audience safety proof: for every declared audience, the scoped
				// build may never omit a never-cut block for that request. This is
				// the `packet check` negative surface for the audience dimension.
				for audience := range m.Audiences {
					ca := Build(src, m, Request{Phase: phase, Track: track, Transport: transport, Optimize: true, Audience: audience})
					req := Request{Phase: phase, Track: track, Transport: transport, Flags: ca.Request.Flags}
					for _, rec := range ca.Index {
						if rec.Included {
							continue
						}
						if neverCutForRequest(rec.Locator, req) {
							rep.NeverCut = appendUnique(rep.NeverCut, fmt.Sprintf("phase %d/%s audience %s omits never-cut block %s", phase, track, audience, rec.Locator))
						}
					}
				}
			}
		}
	}
	rep.OK = m != nil && rep.Secret == "" && len(rep.Unclassified) == 0 && len(rep.Stale) == 0 && len(rep.NeverCut) == 0 && len(rep.Fallbacks) == 0
	return rep
}

func (nc neverCutRule) violation(r Rule) string {
	if r.Include == IncludeAlways {
		return ""
	}
	switch {
	case nc.always:
		return "must be `always`"
	case nc.flag != "":
		if !contains(r.Flags, nc.flag) {
			return "must include flag " + nc.flag
		}
	case nc.transport != "":
		if !contains(r.Transports, nc.transport) {
			return "must include transport " + nc.transport
		}
	case len(nc.phases) > 0:
		var missing []string
		for _, p := range nc.phases {
			found := false
			for _, q := range r.Phases {
				if q == p {
					found = true
					break
				}
			}
			if !found {
				missing = append(missing, strconv.Itoa(p))
			}
		}
		if len(missing) > 0 {
			return "must include phases " + strings.Join(missing, ",")
		}
	}
	return ""
}
