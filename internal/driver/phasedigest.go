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
	"time"

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

	// consensusAbsent records that consensus.md does not exist on disk (drives the
	// G4(b) next-action; unexported, so the marshalled digest shape is unchanged).
	consensusAbsent bool
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
	d.consensusAbsent = false
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
		} else if _, statErr := os.Stat(filepath.Join(ideaDir, "consensus.md")); statErr != nil {
			// Status errors and the file is genuinely absent: rounds-complete +
			// no-consensus.md is a real deck state (Phase 2→3), not an absent one.
			// Recorded unexported — the digest JSON shape is unchanged (fix-up G4(b)).
			d.consensusAbsent = true
		}
	}
	d.Implementation = implSection(ideaDir)
	implNewer := d.Implementation != nil && d.Implementation.ReadyForReview &&
		implementationNewerThanLatestReview(ideaDir, d.Review)
	d.Next = nextAction(&d, implNewer, d.consensusAbsent)
	return d
}

// implementationNewerThanLatestReview reports the fix-up-published state (fix-up
// G4(a), claude-1 R2-MIN-1(a) ≡ kimi-1 K2-F4): the implementation is present and
// ready, the latest review round is complete, and IMPLEMENTATION.md changed after
// every artifact of that round. mtime is the arrival signal — the same signal
// `parley wait` uses for escalations (F2); the digest stays byte-identical over an
// unchanged tree because none of those mtimes move.
func implementationNewerThanLatestReview(ideaDir string, review *PhaseRoundSection) bool {
	if review == nil || review.Total == 0 || review.Completed != review.Total {
		return false
	}
	implInfo, err := os.Stat(filepath.Join(ideaDir, "IMPLEMENTATION.md"))
	if err != nil {
		return false
	}
	newest := time.Time{}
	for _, row := range review.Rows {
		if !row.Filed {
			continue
		}
		if info, ierr := os.Stat(row.Path); ierr == nil && info.ModTime().After(newest) {
			newest = info.ModTime()
		}
	}
	return implInfo.ModTime().After(newest)
}

// fixUpCycleNumberFromStatus parses the N of a `fix-up-cycle-N` implementation
// status; 0 for anything else. The vocabulary is the closed `^fix-up-cycle-\d+$`
// set internal/protocol gates ReadyForReview with; the Sscanf idiom matches
// roundNumberFromLabel's for round-0M.
func fixUpCycleNumberFromStatus(status string) int {
	var n int
	if _, err := fmt.Sscanf(strings.TrimPrefix(status, "fix-up-cycle-"), "%d", &n); err != nil {
		return 0
	}
	return n
}

// fixUpAwaitingReviewRound reports the fix-up-published state by CONTENT (fix-up
// cycle 3 H1, claude-1 R3-MIN-1): IMPLEMENTATION.md's status names fix-up cycle
// N and the latest review round is round-0M with M <= N, so no complete round
// has reviewed cycle N yet (the round that reviews it is round-0(N+1) or later)
// and a further review round is awaited REGARDLESS of timestamps. The mtime
// branch above stays as the arrival signal, but keyed alone it read
// `await implementation` silently on every fresh clone, checkout, worktree
// creation, archive extraction or rsync without -t — operations that give
// IMPLEMENTATION.md and the review artifacts ONE mtime, so the strict After
// never fires (three independent reproductions, review round 3). Content is
// deterministic under every filesystem operation. A complete round-0M with
// M > N reviewed this cycle and falls through exactly as before; an incomplete
// next round never reaches here — Completed < Total returns earlier.
func fixUpAwaitingReviewRound(impl *PhaseImplSection, review *PhaseRoundSection) bool {
	if impl == nil {
		return false
	}
	n := fixUpCycleNumberFromStatus(impl.Status)
	if n <= 0 {
		return false
	}
	if review == nil {
		return true // no review round exists at all: nothing has reviewed cycle N
	}
	return roundNumberFromLabel(review.Label) <= n
}

// nextAction picks from the fixed enumeration only. Any adverse signal
// (❌ block, malformed/error state, unparsed row, invalid artifact) routes the
// organizer to the raw artifact. implNewerThanReview is the G4(a) fix-up-published
// mtime signal; consensusAbsent is the G4(b) no-consensus.md signal.
func nextAction(d *PhaseDigest, implNewerThanReview, consensusAbsent bool) string {
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
	if d.Consensus == nil && consensusAbsent && roundsComplete(d) {
		// Rounds complete but no consensus.md exists at all (fix-up G4(b), claude-1
		// R2-MIN-1(b)): the deck sits at Phase 2→3 awaiting the consensus draft —
		// previously this collapsed into "await implementation", skipping Phases 3–4.
		return NextAwaitConsensus
	}
	if d.Consensus != nil && len(d.Consensus.Missing) > 0 {
		return NextAwaitConsensus
	}
	if d.Implementation != nil && d.Implementation.Present {
		switch {
		case d.Implementation.Status == "complete":
			return NextFinalPublished
		case d.Implementation.ReadyForReview && (d.Review == nil || implNewerThanReview ||
			fixUpAwaitingReviewRound(d.Implementation, d.Review)):
			// Published implementation awaiting review: no review round exists
			// yet (fix-up F10), or the latest round is complete and the
			// implementation is newer than it (mtime arrival signal), or the
			// implementation publishes fix-up cycle N while the latest complete
			// round is round-0M with M <= N (content signal, fix-up cycle 3 H1 —
			// equal mtimes on a fresh checkout must not mask the publication) —
			// the fix-up-published state (fix-up G4(a)): all previously
			// collapsed into "await implementation".
			return NextAwaitReviewArtifact
		}
	}
	return NextAwaitImplementation
}

// roundsComplete: every filed round section is fully complete (a nil review
// section — no review rounds yet — is complete for this purpose).
func roundsComplete(d *PhaseDigest) bool {
	if d.Round != nil && !(d.Round.Total > 0 && d.Round.Completed == d.Round.Total) {
		return false
	}
	if d.Review != nil && !(d.Review.Total > 0 && d.Review.Completed == d.Review.Total) {
		return false
	}
	return d.Round != nil
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
