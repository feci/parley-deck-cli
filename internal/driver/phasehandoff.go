package driver

// phasehandoff.go (lean-organizer D.1): the driver-written per-phase handoff record
// under runs/<run-id>/handoff-phase-<phase>.md. Non-canonical, advisory driver state
// — the recomputed view (`parley organizer brief` / `parley status`) is AUTHORITATIVE
// on any disagreement; no agent is ever obliged to author or refresh it by hand.

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"parley-deck-cli/internal/fsutil"
)

// PhaseHandoffPath is the runs/-relative file name for one phase's record.
func PhaseHandoffPath(phase Phase) string {
	return "handoff-phase-" + phaseSlug(string(phase)) + ".md"
}

var phaseSlugNonAlnum = regexp.MustCompile(`[^a-z0-9]+`)

func phaseSlug(phase string) string {
	s := phaseSlugNonAlnum.ReplaceAllString(strings.ToLower(phase), "-")
	return strings.Trim(s, "-")
}

// PhaseHandoffRecord is the schema of one handoff record. All fields are mechanically
// derived (cursor + tree); none is model-written.
type PhaseHandoffRecord struct {
	RunID         string      `json:"run_id"`
	Idea          string      `json:"idea"`
	Phase         string      `json:"phase"`
	PreviousPhase string      `json:"previous_phase"`
	Action        string      `json:"action"`
	RoundLabel    string      `json:"round_label"`
	IdeaStatus    string      `json:"idea_status"`
	WrittenAt     time.Time   `json:"written_at"`
	Next          string      `json:"next"` // fixed enumeration, shared with PhaseDigest
	Digest        PhaseDigest `json:"digest"`
}

// PhaseHandoffSchemaDoc is stated IN the record: the recomputed view is authoritative.
const PhaseHandoffSchemaDoc = "This record is non-canonical driver state (advisory). On any disagreement the\n" +
	"recomputed view is authoritative: run `parley organizer brief --idea <slug>` or\n" +
	"`parley status --idea <slug> --json`. No agent is obliged to author or refresh\n" +
	"this file by hand; the driver rewrites it at each phase transition."

// BuildPhaseHandoffRecord derives the record from the cursor and the idea tree,
// sharing the PhaseDigest computation. The run dir is a PARAMETER (fix-up F13,
// claude-1 MIN-5): the previous derivation evaluated to "parley-deck" for any
// standard layout and was dead weight every caller had to overwrite.
func BuildPhaseHandoffRecord(root, runDir, ideaSlug, ideaDir string, participants []string, c Cursor, action Action, previous Phase) PhaseHandoffRecord {
	digest := BuildPhaseDigest(root, ideaSlug, ideaDir, participants)
	return PhaseHandoffRecord{
		RunID:         filepath.Base(runDir),
		Idea:          ideaSlug,
		Phase:         string(c.Phase),
		PreviousPhase: string(previous),
		Action:        string(action),
		RoundLabel:    roundLabel(c.CurrentRound),
		IdeaStatus:    c.IdeaStatus,
		WrittenAt:     time.Now().UTC(),
		Next:          digest.Next,
		Digest:        digest,
	}
}

// storeHandoffRecord writes the record atomically (same write discipline as the
// shipped handoff machinery) with the schema doc embedded.
func storeHandoffRecord(path string, record PhaseHandoffRecord) error {
	var b strings.Builder
	fmt.Fprintf(&b, "---\nrun_id: %s\nidea: %s\nphase: %s\nprevious_phase: %s\naction: %s\nround_label: %s\nidea_status: %s\nwritten_at: %s\nnext: %s\n---\n\n",
		record.RunID, record.Idea, record.Phase, record.PreviousPhase, record.Action, record.RoundLabel, record.IdeaStatus, record.WrittenAt.Format(time.RFC3339Nano), record.Next)
	fmt.Fprintf(&b, "# Phase handoff — %s → %s\n\n", record.PreviousPhase, record.Phase)
	fmt.Fprintf(&b, "%s\n\n", PhaseHandoffSchemaDoc)
	digestJSON, err := json.MarshalIndent(record.Digest, "", "  ")
	if err != nil {
		return err
	}
	fmt.Fprintf(&b, "## PhaseDigest (mechanically derived at transition)\n\n```json\n%s\n```\n", string(digestJSON))
	if err := fsutil.MkdirAllResilient(filepath.Dir(path), 0o755); err != nil {
		return err
	}
	return fsutil.WriteFileAtomic(path, []byte(b.String()), 0o644)
}

// LoadPhaseHandoffRecord parses a record back (round-trip verification).
func LoadPhaseHandoffRecord(path string) (PhaseHandoffRecord, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return PhaseHandoffRecord{}, err
	}
	raw := string(data)
	get := func(key string) string {
		re := regexp.MustCompile(`(?m)^` + key + `: (.*)$`)
		if m := re.FindStringSubmatch(raw); len(m) == 2 {
			return strings.TrimSpace(m[1])
		}
		return ""
	}
	// Fix-up F13 (claude-1 MIN-5): parse ALL nine frontmatter fields so the
	// round-trip is complete — the loader previously restored only 4 of 9.
	rec := PhaseHandoffRecord{
		RunID:         get("run_id"),
		Idea:          get("idea"),
		Phase:         get("phase"),
		PreviousPhase: get("previous_phase"),
		Action:        get("action"),
		RoundLabel:    get("round_label"),
		IdeaStatus:    get("idea_status"),
		Next:          get("next"),
	}
	if written := get("written_at"); written != "" {
		if when, werr := time.Parse(time.RFC3339Nano, written); werr == nil {
			rec.WrittenAt = when
		}
	}
	if rec.RunID == "" || rec.Idea == "" || rec.Phase == "" {
		return rec, fmt.Errorf("handoff record %s is missing required fields", path)
	}
	return rec, nil
}
