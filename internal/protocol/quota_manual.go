package protocol

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"parley-deck-cli/internal/quota"
)

// Existing knob-off §9.0 confirmations stay usable without a new CLI command.
// The prompt records a manual change, never owner-confirmed authority. We read the explicit current set;
// excluded markers only validate a requested change, never subtract membership.
func manualQuotaRevision(dir string, h *quota.History, raw string, ids []string) (quota.Batch, error) {
	if len(ids) == 0 || len(quota.Unique(ids)) != len(ids) {
		return quota.Batch{}, fmt.Errorf("invalid manual membership")
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
	for _, id := range ids {
		if has(h.Current, id) {
			continue
		}
		if !has(h.Known, id) {
			path := filepath.Join(dir, "round-01", id+".md")
			if err := ValidateParticipantRoundArtifact(path, id, h.Kickoff.Idea, 1); err != nil {
				return quota.Batch{}, err
			}
			// Capture bytes via the caller below; immutable history retains the proof.
			data, err := readQuotaCatchup(path)
			if err != nil {
				return quota.Batch{}, err
			}
			catchup[id] = data
		}
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
