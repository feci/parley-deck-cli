// Package driver advances a Parley Deck deliberation past round-01: it owns the
// "next tick" that the one-shot runner never had. The deliberation phase graph is
// degenerate-linear (round → cross-review → consensus → final → impl →
// review/fix-up) with one back-edge (consensus BLOCK → reopen round), so the
// driver is a small ordered switch over a disk-derived cursor — NOT a DAG,
// scheduler, worker pool, or claim-locked task kernel. Disk is the single source
// of truth; the cursor is a rebuildable cache. See
// parley-deck/ideas/deliberation-driver/FINAL.md (consensus D1-D15).
//
// Slice 1 (this file set) implements ONLY the round phase: promote a completed
// round-01 to round-02 via the injected RoundRunner. Consensus/final/impl phases
// are later slices.
package driver

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"parley-deck-cli/internal/fsutil"
)

// Phase is the deliberation phase the cursor currently rests in.
type Phase string

const (
	PhaseRound     Phase = "round"
	PhaseConsensus Phase = "consensus"
	PhaseFinal     Phase = "final"
	PhaseImpl      Phase = "impl"
	PhaseReview    Phase = "review"
	PhaseDone      Phase = "done"
	PhaseBlocked   Phase = "blocked"
)

// Cursor is a rebuildable cache of the deliberation phase. Disk is authoritative:
// Rebuild derives every field from on-disk artifacts and the persisted Phase is
// never trusted over disk. Only MaxRounds is config-derived.
//
// One field is deliberately NOT rebuildable: FixupCyclesPublished is a safety count that
// must not be recoverable by editing artifacts, so Advance carries it forward from the
// persisted cursor instead of deriving it.
type Cursor struct {
	Phase        Phase  `json:"phase"`
	CurrentRound int    `json:"current_round"`
	IdeaStatus   string `json:"idea_status"`
	RoundsRun    int    `json:"rounds_run"`
	MaxRounds    int    `json:"max_rounds"`
	UpdatedAt    string `json:"updated_at"`
	// FixupCyclesPublished is the driver's own monotonic count of CHARGED fix-up
	// attempts — reserved before the code-writing call, so an attempt that errors is
	// counted too. (The name predates the corrected unit; `.fixup-done` markers are the
	// record of cycles that actually completed.) It lives in the run directory, NOT in the idea directory, so no
	// participant artifact can lower it. The budget takes the MAXIMUM of this and the
	// on-disk `.fixup-done` markers: deleting markers cannot buy a cycle, and forging
	// one can only raise the count, which escalates sooner. Review round-03 showed the
	// marker-only count was still editable state — the class had moved, not closed.
	FixupCyclesPublished int `json:"fixup_cycles_published,omitempty"`
}

// Save writes the cursor atomically (tmp + rename, same dir) so a crash mid-write
// cannot corrupt the durable state. Mirrors internal/pipeline/run.go Save.
func (c Cursor) Save(path string) error {
	if err := fsutil.MkdirAllResilient(filepath.Dir(path), 0o755); err != nil {
		return fmt.Errorf("create driver dir: %w", err)
	}
	data, err := json.MarshalIndent(c, "", "  ")
	if err != nil {
		return fmt.Errorf("encode cursor: %w", err)
	}
	data = append(data, '\n')
	tmp := path + ".tmp"
	if err := os.WriteFile(tmp, data, 0o644); err != nil {
		return fmt.Errorf("write cursor: %w", err)
	}
	if err := os.Rename(tmp, path); err != nil {
		return fmt.Errorf("commit cursor: %w", err)
	}
	return nil
}

// LoadCursor best-effort reads a persisted cursor. A missing or corrupt file is
// non-fatal: callers Rebuild from disk, which is authoritative.
func LoadCursor(path string) (Cursor, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return Cursor{}, err
	}
	var c Cursor
	if err := json.Unmarshal(data, &c); err != nil {
		return Cursor{}, err
	}
	return c, nil
}

// Rebuild derives the phase purely from disk (consensus D2/D3). ideaDir is the
// idea directory (…/ideas/<slug>). maxRounds is the config circuit-breaker bound.
func Rebuild(ideaDir string, maxRounds int) Cursor {
	detail, _ := RebuildDetail(ideaDir, maxRounds)
	return detail.Cursor
}

// PhaseDetail augments the rebuilt Cursor with the display-only disk evidence a
// consumer needs to disambiguate the review steps (6 review vs 7 review-consensus
// vs 8 fix-up/complete) without re-probing the idea dir
// (consensus tui-protocol-visibility D2). It is derived, never persisted.
type PhaseDetail struct {
	Cursor                Cursor
	HighestReviewRound    int
	ReviewConsensusExists bool
	ImplementationStatus  string
	FinalScaffoldReason   string // "" when FINAL.md is absent or acceptable
}

// RebuildDetail derives the cursor and its display evidence in one disk pass.
// Missing artifacts are normal zero values; an unreadable idea/review directory
// or an unexpected stat/read error returns the partial detail with a non-nil
// error so callers can keep their previous snapshot instead of trusting a
// half-read state (consensus D2; review cycle-1 fix 1).
func RebuildDetail(ideaDir string, maxRounds int) (PhaseDetail, error) {
	var firstErr error
	keepErr := func(err error) {
		if err != nil && firstErr == nil {
			firstErr = err
		}
	}
	c := Cursor{Phase: PhaseRound, CurrentRound: 1, MaxRounds: maxRounds}
	status, _, err := readFrontmatterFieldErr(filepath.Join(ideaDir, "00-prompt.md"), "status")
	keepErr(err)
	c.IdeaStatus = status
	highest, err := highestRoundErr(ideaDir)
	keepErr(err)
	if highest >= 1 {
		c.CurrentRound = highest
		c.RoundsRun = highest
	}
	finalPath := filepath.Join(ideaDir, "FINAL.md")
	implPath := filepath.Join(ideaDir, "IMPLEMENTATION.md")
	reviewConsensus := filepath.Join(ideaDir, "review", "consensus.md")
	var detail PhaseDetail
	reviewConsensusExists, err := statRegular(reviewConsensus)
	keepErr(err)
	detail.ReviewConsensusExists = reviewConsensusExists
	reviewRound, err := highestReviewRoundErr(ideaDir)
	keepErr(err)
	detail.HighestReviewRound = reviewRound
	implExists, err := statRegular(implPath)
	keepErr(err)
	if implExists {
		implStatus, _, err := readFrontmatterFieldErr(implPath, "status")
		keepErr(err)
		detail.ImplementationStatus = implStatus
	}
	finalExists, err := statRegular(finalPath)
	keepErr(err)
	if finalExists {
		detail.FinalScaffoldReason = finalScaffoldReason(finalPath)
		if detail.FinalScaffoldReason == "FINAL.md is missing" {
			// The stat saw the file but the read failed — surface it instead of
			// classifying the idea as pre-final.
			keepErr(fmt.Errorf("FINAL.md exists but could not be read"))
		}
	}
	consensusExists, err := statRegular(filepath.Join(ideaDir, "consensus.md"))
	keepErr(err)
	switch {
	// Most-terminal-first (D2): implementation/review artifacts win over FINAL/
	// consensus so a valid FINAL.md never hides later phases.
	case implExists && detail.ImplementationStatus == "complete":
		c.Phase = PhaseDone
	case detail.ReviewConsensusExists || detail.HighestReviewRound >= 1:
		c.Phase = PhaseReview
	case implExists:
		c.Phase = PhaseImpl
	case finalExists && detail.FinalScaffoldReason == "":
		// Only a VALID (non-scaffold) FINAL.md is truly final. A scaffold FINAL.md
		// from a failed/partial draft must NOT strand the idea at PhaseFinal — it
		// stays in the consensus phase so the gate re-drafts it (slice-2 AF1).
		c.Phase = PhaseFinal
	case consensusExists:
		c.Phase = PhaseConsensus
	case finalExists:
		// Scaffold FINAL.md with no consensus.md to re-drive: treat as final to
		// avoid a phantom round phase; the surface-only stop surfaces it to a human.
		c.Phase = PhaseFinal
	default:
		c.Phase = PhaseRound
	}
	c.UpdatedAt = time.Now().UTC().Format(time.RFC3339)
	detail.Cursor = c
	return detail, firstErr
}

// implementationStatus returns the status: frontmatter of an IMPLEMENTATION.md.
func implementationStatus(implPath string) string {
	v, _ := readFrontmatterField(implPath, "status")
	return v
}

// highestReviewRound returns the largest N for which review/round-NN/ exists, or 0.
func highestReviewRound(ideaDir string) int {
	n, _ := highestReviewRoundErr(ideaDir)
	return n
}

// highestReviewRoundErr is highestReviewRound, surfacing non-NotExist ReadDir
// errors (a half-readable review dir must not silently read as "no rounds").
func highestReviewRoundErr(ideaDir string) (int, error) {
	entries, err := os.ReadDir(filepath.Join(ideaDir, "review"))
	if err != nil {
		if errors.Is(err, fs.ErrNotExist) {
			return 0, nil
		}
		return 0, err
	}
	highest := 0
	for _, e := range entries {
		if e.IsDir() && strings.HasPrefix(e.Name(), "round-") {
			if n := roundOrdinal(e.Name()); n > highest {
				highest = n
			}
		}
	}
	return highest, nil
}

// highestRound returns the largest N for which a round-NN/ directory exists, or 0.
func highestRound(ideaDir string) int {
	n, _ := highestRoundErr(ideaDir)
	return n
}

// highestRoundErr is highestRound, surfacing non-NotExist ReadDir errors.
func highestRoundErr(ideaDir string) (int, error) {
	entries, err := os.ReadDir(ideaDir)
	if err != nil {
		if errors.Is(err, fs.ErrNotExist) {
			return 0, nil
		}
		return 0, err
	}
	highest := 0
	for _, e := range entries {
		if !e.IsDir() || !strings.HasPrefix(e.Name(), "round-") {
			continue
		}
		if n := roundOrdinal(e.Name()); n > highest {
			highest = n
		}
	}
	return highest, nil
}

// roundOrdinal parses "round-NN" → N, or 0 when the label is not a round dir.
func roundOrdinal(label string) int {
	if !strings.HasPrefix(label, "round-") {
		return 0
	}
	n, err := strconv.Atoi(strings.TrimLeft(strings.TrimPrefix(label, "round-"), "0"))
	if err != nil {
		return 0
	}
	return n
}

func roundLabel(n int) string { return fmt.Sprintf("round-%02d", n) }

func nowRFC3339() string { return time.Now().UTC().Format(time.RFC3339) }

// readIdeaStatus reads the status: field from the idea 00-prompt.md frontmatter.
func readIdeaStatus(ideaDir string) string {
	field, _ := readFrontmatterField(filepath.Join(ideaDir, "00-prompt.md"), "status")
	return field
}

// statRegular reports whether a regular file exists at path; a missing file is
// (false, nil), any other stat error is surfaced (consensus D2).
func statRegular(path string) (bool, error) {
	info, err := os.Stat(path)
	if err != nil {
		if errors.Is(err, fs.ErrNotExist) {
			return false, nil
		}
		return false, err
	}
	return info.Mode().IsRegular(), nil
}

// readFrontmatterField returns the trimmed value of a top-level frontmatter key.
func readFrontmatterField(path, key string) (string, bool) {
	value, ok, _ := readFrontmatterFieldErr(path, key)
	return value, ok
}

// readFrontmatterFieldErr is readFrontmatterField with non-NotExist read errors
// surfaced (a missing file stays a normal zero value).
func readFrontmatterFieldErr(path, key string) (string, bool, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		if errors.Is(err, fs.ErrNotExist) {
			return "", false, nil
		}
		return "", false, err
	}
	inFrontmatter := false
	prefix := key + ":"
	for _, line := range strings.Split(string(data), "\n") {
		trimmed := strings.TrimSpace(line)
		if trimmed == "---" {
			if !inFrontmatter {
				inFrontmatter = true
				continue
			}
			break
		}
		if inFrontmatter && strings.HasPrefix(trimmed, prefix) {
			return strings.TrimSpace(strings.TrimPrefix(trimmed, prefix)), true, nil
		}
	}
	return "", false, nil
}

func fileExists(path string) bool {
	info, err := os.Stat(path)
	return err == nil && info.Mode().IsRegular()
}

func dirExists(path string) bool {
	info, err := os.Stat(path)
	return err == nil && info.IsDir()
}
