package driver

// phasedigest.go extends the shipped round digest (digest.go) into a PhaseDigest:
// a deterministic, LLM-free, mechanically derived snapshot of one idea covering the
// latest design round, the latest review round, consensus signoff state, and
// implementation status (lean-organizer B).
//
// Guardrails (ratified): every field is derived from the filesystem or the shipped
// validators — no field is model-written; validity carries the shipped validator's
// verdict verbatim, with the failing check named; stance is keyword FLAGS, never a
// verdict; unparsed is fail-closed; the next-action line is a FIXED enumeration,
// never generated prose; every row carries its raw path; the digest is printed,
// never persisted as canonical (`runs/` telemetry capture excepted); raw artifacts
// stay canonical and any ❌ / DISPUTED / unparsed / adverse validity sends the
// organizer to the raw artifact.

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"parley-deck-cli/internal/consensus"
	"parley-deck-cli/internal/protocol"
	"parley-deck-cli/internal/runner"
)

// PhaseAgentRow is one agent's mechanically derived row.
type PhaseAgentRow struct {
	Agent       string `json:"agent"`
	Path        string `json:"path"`
	Filed       bool   `json:"filed"`
	Bytes       int    `json:"bytes"`
	Owner       string `json:"owner,omitempty"`    // frontmatter `agent:`
	Valid       bool   `json:"valid"`              // shipped validator verdict
	Validity    string `json:"validity,omitempty"` // verdict verbatim; "" = clean
	StanceFlags struct {
		Block    int `json:"block"`
		Counter  int `json:"counter"`
		Accept   int `json:"accept"`
		Escalate int `json:"escalate"`
	} `json:"stance_flags"`
	Unparsed bool `json:"unparsed"` // fail-closed: attribution or stance could not be parsed
	// FellBack reports degraded metadata derivation on the validator/ownership path:
	// the artifact frontmatter was unreadable, or carried no usable `agent:` owner so
	// attribution fell back to the filename-derived agent (fix-up F9 — it is NOT the
	// round-digest's "## Summary absent" extraction flag; PhaseDigest discards
	// position prose entirely, so that extraction says nothing about this digest).
	FellBack bool `json:"fell_back"`
}

// PhaseRoundSection is the latest design or review round.
type PhaseRoundSection struct {
	Label     string          `json:"label"`
	Review    bool            `json:"review"`
	Total     int             `json:"total"`
	Completed int             `json:"completed"`
	Rows      []PhaseAgentRow `json:"rows"`
}

// PhaseConsensusSection is the design-consensus signoff state, reusing the
// consensus-status parser (no forked stance parsing).
type PhaseConsensusSection struct {
	Present      bool     `json:"present"`
	Path         string   `json:"path,omitempty"`
	Triage       string   `json:"triage,omitempty"`
	Missing      []string `json:"missing,omitempty"`
	Signed       []string `json:"signed,omitempty"`
	Reservations []string `json:"reservations,omitempty"`
	Blocks       []string `json:"blocks,omitempty"`
	Errors       []string `json:"errors,omitempty"`
}

// PhaseImplSection is the implementation status.
type PhaseImplSection struct {
	Present        bool   `json:"present"`
	Status         string `json:"status,omitempty"`
	Implementer    string `json:"implementer,omitempty"`
	HeadCommit     string `json:"head_commit,omitempty"`
	ReadyForReview bool   `json:"ready_for_review"`
}

// PhaseDigest is the whole-idea snapshot. Deterministic and byte-identical over an
// unchanged tree: no timestamps, no absolute-clock values, rows sorted by agent.
type PhaseDigest struct {
	Idea           string                 `json:"idea"`
	Round          *PhaseRoundSection     `json:"round,omitempty"`
	Review         *PhaseRoundSection     `json:"review,omitempty"`
	Consensus      *PhaseConsensusSection `json:"consensus,omitempty"`
	Implementation *PhaseImplSection      `json:"implementation,omitempty"`
	Next           string                 `json:"next"`
}

// Next-action vocabulary — a fixed enumeration, never generated prose.
const (
	NextAwaitRoundArtifact  = "await round artifact"
	NextAwaitConsensus      = "await consensus signoff"
	NextAwaitReviewArtifact = "await review artifact"
	NextAwaitImplementation = "await implementation"
	NextAdjudicateRaw       = "open the raw artifact to adjudicate"
	NextFinalPublished      = "FINAL published"
)

var nextActionVocabulary = map[string]bool{
	NextAwaitRoundArtifact: true, NextAwaitConsensus: true, NextAwaitReviewArtifact: true,
	NextAwaitImplementation: true, NextAdjudicateRaw: true, NextFinalPublished: true,
}

// FixedNextVocabulary exposes the enumeration for tests and the brief.
func FixedNextVocabulary() []string {
	out := make([]string, 0, len(nextActionVocabulary))
	for k := range nextActionVocabulary {
		out = append(out, k)
	}
	sort.Strings(out)
	return out
}

// IsValidNextAction reports whether s is a member of the fixed enumeration.
func IsValidNextAction(s string) bool { return nextActionVocabulary[s] }

// BuildPhaseDigest reads the idea tree read-only and builds the digest. It never
// errors: an unreadable section is omitted; an unreadable artifact is a fail-closed
// unparsed row. participants is the idea's participant list.
func BuildPhaseDigest(root, ideaSlug, ideaDir string, participants []string) PhaseDigest {
	d := PhaseDigest{Idea: ideaSlug}
	if sec, ok := latestRoundSection(ideaDir, participants, false); ok {
		d.Round = sec
	}
	if sec, ok := latestRoundSection(ideaDir, participants, true); ok {
		d.Review = sec
	}
	if root != "" {
		if summary, err := consensus.Status(root, ideaSlug, false); err == nil {
			sec := &PhaseConsensusSection{
				Present: true,
				Path:    summary.Path,
				Triage:  string(summary.Triage),
				Missing: summary.Missing,
				Errors:  summary.Errors,
			}
			for _, s := range summary.Signoffs {
				sec.Signed = append(sec.Signed, s.Agent)
				switch s.Status {
				case consensus.StatusReservations:
					sec.Reservations = append(sec.Reservations, s.Agent)
				case consensus.StatusBlock:
					sec.Blocks = append(sec.Blocks, s.Agent)
				}
			}
			if summary.Triage == consensus.TriageMalformed {
				sec.Errors = append(sec.Errors, "consensus document malformed — open the raw artifact")
			}
			d.Consensus = sec
		}
	}
	d.Implementation = implSection(ideaDir)
	d.Next = nextAction(&d)
	return d
}

// nextAction picks from the fixed enumeration only. Any adverse signal
// (❌ block, malformed/error state, unparsed row, invalid artifact) routes the
// organizer to the raw artifact.
func nextAction(d *PhaseDigest) string {
	if d.Consensus != nil && (len(d.Consensus.Blocks) > 0 || len(d.Consensus.Errors) > 0) {
		return NextAdjudicateRaw
	}
	for _, sec := range []*PhaseRoundSection{d.Round, d.Review} {
		if sec == nil {
			continue
		}
		for _, row := range sec.Rows {
			if row.Unparsed || (row.Filed && !row.Valid) {
				return NextAdjudicateRaw
			}
		}
	}
	if d.Review != nil && d.Review.Completed < d.Review.Total {
		return NextAwaitReviewArtifact
	}
	if d.Round != nil && d.Round.Completed < d.Round.Total {
		return NextAwaitRoundArtifact
	}
	if d.Consensus != nil && len(d.Consensus.Missing) > 0 {
		return NextAwaitConsensus
	}
	if d.Implementation != nil && d.Implementation.Present {
		switch {
		case d.Implementation.Status == "complete":
			return NextFinalPublished
		case d.Implementation.ReadyForReview && d.Review == nil:
			// Published implementation, no review round opened yet (fix-up F10:
			// this state previously collapsed into "await implementation").
			return NextAwaitReviewArtifact
		}
	}
	return NextAwaitImplementation
}

func implSection(ideaDir string) *PhaseImplSection {
	path := filepath.Join(ideaDir, "IMPLEMENTATION.md")
	if _, err := os.ReadFile(path); err != nil {
		return &PhaseImplSection{}
	}
	meta, err := protocol.ReadFrontmatter(path)
	if err != nil {
		return &PhaseImplSection{Present: true, Status: "unparsed"}
	}
	status := strings.TrimSpace(meta["status"])
	sec := &PhaseImplSection{
		Present:     true,
		Status:      status,
		Implementer: strings.TrimSpace(meta["implementer"]),
		HeadCommit:  strings.TrimSpace(meta["head-commit"]),
	}
	// Ready-for-review is DERIVED from the protocol's closed vocabulary
	// (protocol.ValidImplementationStatus) minus the in-progress states, not from a
	// second hand-written list — the hand list omitted `ready-for-review`, a status
	// four files in this deck live at (fix-up F3, claude-1 MAJ-3).
	switch status {
	case "", "unparsed", "in-progress":
		sec.ReadyForReview = false
	default:
		sec.ReadyForReview = protocol.ValidImplementationStatus(status)
	}
	return sec
}

// latestRoundSection finds the highest-N round dir under round-*/ (design) or
// review/round-*/ (review) that exists, and validates each expected participant's
// artifact with the SHIPPED validators.
func latestRoundSection(ideaDir string, participants []string, review bool) (*PhaseRoundSection, bool) {
	var dirs []string
	root := ideaDir
	if review {
		root = filepath.Join(ideaDir, "review")
	}
	entries, err := os.ReadDir(root)
	if err != nil {
		return nil, false
	}
	for _, e := range entries {
		if !e.IsDir() || !strings.HasPrefix(e.Name(), "round-") {
			continue
		}
		dirs = append(dirs, e.Name())
	}
	if len(dirs) == 0 {
		return nil, false
	}
	sort.Sort(sort.Reverse(sort.StringSlice(dirs)))
	label := dirs[0]
	roundDir := filepath.Join(root, label)
	round := roundNumberFromLabel(label)

	expected := participants
	if review {
		expected = consensus.ExpectedRoundParticipants(ideaDir, participants, true)
	}
	sec := &PhaseRoundSection{Label: label, Review: review, Total: len(expected)}
	for _, p := range expected {
		row := agentRow(filepath.Join(roundDir, p+".md"), p, ideaDir, round, review)
		sec.Rows = append(sec.Rows, row)
		if row.Filed && row.Valid {
			sec.Completed++
		}
	}
	return sec, true
}

func roundNumberFromLabel(label string) int {
	var n int
	if _, err := fmt.Sscanf(strings.TrimPrefix(label, "round-"), "%d", &n); err != nil {
		return 0
	}
	return n
}

func agentRow(path, agent, ideaDir string, round int, review bool) PhaseAgentRow {
	row := PhaseAgentRow{Agent: agent, Path: path}
	data, err := os.ReadFile(path)
	if err != nil {
		return row
	}
	row.Filed = true
	row.Bytes = len(data)
	meta, merr := protocol.ReadFrontmatter(path)
	row.Owner = strings.Trim(strings.TrimSpace(meta["agent"]), `"'`)
	if merr != nil || row.Owner == "" || row.Owner != agent {
		row.Unparsed = true
	}
	// fell_back: the metadata path degraded — frontmatter unreadable, or no usable
	// `agent:` owner so attribution falls back to the filename-derived agent (F9).
	row.FellBack = merr != nil || row.Owner == ""
	if review {
		if verr := protocol.ValidateReviewArtifact(path, agent, filepath.Base(ideaDir), round); verr != nil {
			row.Validity = verr.Error()
		} else {
			row.Valid = true
		}
	} else {
		if verr := runner.ValidateRoundArtifact(path, agent, filepath.Base(ideaDir), round); verr != nil {
			row.Validity = verr.Error()
		} else {
			row.Valid = true
		}
	}
	b, c, a, e := stanceFlags(string(data))
	row.StanceFlags.Block, row.StanceFlags.Counter, row.StanceFlags.Accept, row.StanceFlags.Escalate = b, c, a, e
	return row
}
