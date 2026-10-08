package protocol

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"time"

	"parley-deck-cli/internal/quota"
)

// Existing knob-off §9.0 confirmations stay usable without a new CLI command.
// The prompt records a manual change, never owner-confirmed authority. We read the explicit current set;
// excluded markers only validate a requested change, never subtract membership.
func manualQuotaRevision(dir string, h *quota.History, raw string, ids []string) (quota.Batch, error) {
	if err := h.CheckReturn(ids); err != nil {
		return quota.Batch{}, err
	}
	if len(ids) == 0 || len(quota.Unique(ids)) != len(ids) {
		return quota.Batch{}, fmt.Errorf("invalid manual membership")
	}
	for _, id := range ids {
		if !quotaManualID.MatchString(id) {
			return quota.Batch{}, fmt.Errorf("invalid manual member identity %q", id)
		}
	}
	meta, err := ReadFrontmatter(filepath.Join(dir, "00-prompt.md"))
	if err != nil {
		return quota.Batch{}, err
	}
	if meta["status"] == "final" || meta["status"] == "closed" {
		return quota.Batch{}, fmt.Errorf("closed idea membership is frozen")
	}
	has := func(xs []string, id string) bool {
		for _, x := range xs {
			if x == id {
				return true
			}
		}
		return false
	}
	confirmed := func(key, id string) bool {
		return quota.ConfirmedInPrompt(raw, key, id)
	}
	for _, id := range h.Current {
		if !has(ids, id) && !confirmed("excluded", id) {
			return quota.Batch{}, quota.ManualExclusionError(id)
		}
	}
	catchup := map[string]string{}
	pending := &manualCatchupPending{}
	for _, id := range ids {
		if has(h.Current, id) {
			continue
		}
		if !has(h.Known, id) && !quota.ManualRoundOneReturn(raw) {
			path := filepath.Join(dir, "round-01", id+".md")
			if !manualCatchupPathSafe(dir, path) {
				return quota.Batch{}, fmt.Errorf("ambiguous catch-up path: %s", path)
			}
			data, err := readQuotaCatchup(path)
			if err != nil && !os.IsNotExist(err) {
				return quota.Batch{}, err
			}
			if err == nil {
				st, e := os.Lstat(path)
				if e != nil || !st.Mode().IsRegular() {
					return quota.Batch{}, fmt.Errorf("catch-up path must be a regular file: %s", path)
				}
			}
			if err != nil || ValidateParticipantRoundArtifact(path, id, h.Kickoff.Idea, 1) != nil || quota.ValidateManualCatchupSnapshot(data, h.Kickoff.Idea, id) != nil {
				pending.ids = append(pending.ids, id)
				pending.commands = append(pending.commands, fmt.Sprintf("parley agents exec --agent %s --artifact %q --prompt-file <catch-up-prompt> --yes", id, filepath.ToSlash(filepath.Join(DeckDir, "ideas", h.Kickoff.Idea, "round-01", id+".md"))))
				continue
			}
			catchup[id] = data
		}
	}
	if len(pending.ids) > 0 {
		return quota.Batch{}, pending
	}
	return quota.NewRevision(h, h.Kickoff.RunID, ids, h.Policy(), quota.Revision{ManualPrompt: raw, Catchup: catchup}, time.Now()), nil
}

// Catch-up must be an attributable late round-1 that explicitly records reading
// prior rounds and joining from round 2, under the existing §5 conditions.
func ValidateQuotaCatchup(path, slug, id string) error {
	if err := ValidateParticipantRoundArtifact(path, id, slug, 1); err != nil {
		return err
	}
	raw, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	if err = quota.ValidateCatchupSnapshot(string(raw), slug, id); err != nil {
		return err
	}
	priors, err := quota.CatchupPriors(string(raw))
	if err != nil {
		return err
	}
	for _, prior := range priors {
		full := filepath.Join(filepath.Dir(filepath.Dir(path)), filepath.FromSlash(prior))
		parts := strings.Split(prior, "/")
		round, e := strconv.Atoi(strings.TrimPrefix(parts[len(parts)-2], "round-"))
		if e != nil {
			return e
		}
		if e = ValidateParticipantRoundArtifact(full, strings.TrimSuffix(filepath.Base(prior), ".md"), slug, round); e != nil {
			return fmt.Errorf("catch-up prior unavailable or invalid: %w", e)
		}
	}
	return nil
}

func readQuotaCatchup(path string) (string, error) { b, e := os.ReadFile(path); return string(b), e }

var quotaManualID = regexp.MustCompile(`^[A-Za-z0-9][A-Za-z0-9._-]*$`)

type manualCatchupPending struct{ ids, commands []string }

func (p *manualCatchupPending) Error() string {
	return "pending catch-up for " + strings.Join(p.ids, ", ") + "; complete each own late round-1 artifact, then the manual edit can be imported: " + strings.Join(p.commands, "; ")
}

// ManualCatchupTarget is only a launch permission for an unknown policy-off
// joiner's own late round-1. It grants neither membership nor owner authority.
// A known excluded member, or a kickoff exclusion without a participant edit,
// cannot use the exception. An incomplete own stub may be retried in place.
func ManualCatchupTarget(dir string, v QuotaView, id, path string) bool {
	h := v.History
	if h == nil || h.Dropped(id) || h.Policy().Enabled || !quotaManualID.MatchString(id) || (v.Pending != "" && len(v.Catchup) == 0) {
		return false
	}
	for _, known := range h.Known {
		if id == known {
			return false
		}
	}
	want, err := filepath.Abs(filepath.Join(dir, "round-01", id+".md"))
	target, e := filepath.Abs(path)
	if err != nil || e != nil || target != want {
		return false
	}
	meta, err := ReadFrontmatter(filepath.Join(dir, "00-prompt.md"))
	if err != nil || meta["status"] == "final" || meta["status"] == "closed" {
		return false
	}
	raw, err := os.ReadFile(filepath.Join(dir, "00-prompt.md"))
	if err != nil {
		return false
	}
	listed := false
	for _, member := range parseList(meta["participants"]) {
		listed = listed || member == id
	}
	if quota.ConfirmedInPrompt(string(raw), "excluded", id) && !listed {
		return false
	}
	if !manualCatchupPathSafe(dir, target) {
		return false
	}
	if st, err := os.Lstat(path); err == nil {
		if !st.Mode().IsRegular() || ValidateParticipantRoundArtifact(path, id, h.Kickoff.Idea, 1) == nil {
			return false
		}
		m, err := readRoundOneFrontmatter(path)
		if err != nil {
			return false
		}
		for key, want := range map[string]string{"agent": id, "idea": h.Kickoff.Idea, "round": "1"} {
			if got := strings.Trim(m[key], "\"'"); got != "" && got != want {
				return false
			}
		}
	} else if !os.IsNotExist(err) {
		return false
	}
	return true
}

func manualCatchupPathSafe(dir, path string) bool {
	target, err := filepath.Abs(path)
	if err != nil {
		return false
	}
	idea, err := filepath.Abs(dir)
	if err != nil {
		return false
	}
	for p := target; ; p = filepath.Dir(p) {
		st, err := os.Lstat(p)
		if err != nil && !os.IsNotExist(err) {
			return false
		}
		if err == nil && st.Mode()&os.ModeSymlink != 0 {
			return false
		}
		if p == filepath.Dir(idea) || p == filepath.Dir(p) {
			break
		}
	}
	return true
}
