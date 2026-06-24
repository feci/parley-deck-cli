// Package loop implements LE-9: the one-shot, human-braked outer loop
// (`parley loop tick`). It discovers candidate signals and drafts non-active
// `status: candidate` idea prompts ONLY — it never staffs a quorum, runs a
// deliberation, implements, pushes, merges, or finalizes. Those steps are the
// human-gated promotion described by COOPERATION.md §14 (LE-8, the human brake).
//
// The loop is disabled by default and fail-safe: with no config it writes
// nothing and returns cleanly, so a cron/CI/MCP scheduler can call it idempotently.
package loop

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"time"

	"parley-deck-cli/internal/protocol"
)

// Candidate is a discovered signal that may become a candidate idea.
type Candidate struct {
	Source      string `json:"source"` // commit | ci | issue | manual
	ID          string `json:"id"`
	Title       string `json:"title,omitempty"`
	Detail      string `json:"detail,omitempty"`
	Fingerprint string `json:"fingerprint,omitempty"`
}

// Config is the loop's persisted on/off state (parley-deck/loop/config.json).
// Disabled by default: an absent or unparseable file means Enabled=false.
type Config struct {
	Enabled bool `json:"enabled"`
}

// TickResult reports what a single tick did. Created/Skipped hold idea slugs.
type TickResult struct {
	Enabled bool     `json:"enabled"`
	Created []string `json:"created"`
	Skipped []string `json:"skipped"` // already-existing candidate ideas (dedupe)
}

// ConfigPath / SignalsPath are the default locations under the deck.
func ConfigPath(deck string) string  { return filepath.Join(deck, "loop", "config.json") }
func SignalsPath(deck string) string { return filepath.Join(deck, "loop", "signals.json") }

// ReadConfig reads loop/config.json. A missing file is the disabled default
// (no error); a present-but-malformed file fails closed (disabled + error) so a
// broken config can never silently enable the loop.
func ReadConfig(deck string) (Config, error) {
	data, err := os.ReadFile(ConfigPath(deck))
	if errors.Is(err, fs.ErrNotExist) {
		return Config{Enabled: false}, nil
	}
	if err != nil {
		return Config{Enabled: false}, err
	}
	var c Config
	if err := json.Unmarshal(data, &c); err != nil {
		return Config{Enabled: false}, fmt.Errorf("loop config %s: %w", ConfigPath(deck), err)
	}
	return c, nil
}

// ReadSignals reads a JSON array of Candidate from path. A missing file yields
// an empty slice and no error (a cron tick with nothing queued is normal).
func ReadSignals(path string) ([]Candidate, error) {
	data, err := os.ReadFile(path)
	if errors.Is(err, fs.ErrNotExist) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	var sigs []Candidate
	if err := json.Unmarshal(data, &sigs); err != nil {
		return nil, fmt.Errorf("loop signals %s: %w", path, err)
	}
	return sigs, nil
}

// fingerprintOf returns the candidate's fingerprint, deriving a stable short hash
// of source+id when none was supplied.
func fingerprintOf(c Candidate) string {
	if fp := strings.TrimSpace(c.Fingerprint); fp != "" {
		return sanitize(fp)
	}
	sum := sha256.Sum256([]byte(c.Source + ":" + c.ID))
	return hex.EncodeToString(sum[:])[:8]
}

// sanitize lowercases and keeps only [a-z0-9-] so a value is slug-safe.
func sanitize(s string) string {
	s = strings.ToLower(strings.TrimSpace(s))
	var b strings.Builder
	for _, r := range s {
		switch {
		case r >= 'a' && r <= 'z', r >= '0' && r <= '9', r == '-':
			b.WriteRune(r)
		case r == ' ' || r == '_' || r == '/' || r == ':':
			b.WriteRune('-')
		}
	}
	out := strings.Trim(b.String(), "-")
	if out == "" {
		return "x"
	}
	return out
}

// SlugFor is the deterministic candidate-idea slug for a signal. An auto-derived
// fingerprint is already an 8-char hash; an explicit, human-supplied fingerprint is
// kept whole (sanitized) so the slug stays readable.
func SlugFor(c Candidate) string {
	src := sanitize(c.Source)
	if src == "" {
		src = "signal"
	}
	return "loop-" + src + "-" + fingerprintOf(c)
}

// Tick is one-shot LE-9. When disabled it writes nothing. When enabled it drafts
// a `status: candidate` idea for each not-yet-seen signal (dedupe by slug) and
// returns. It NEVER staffs a quorum, runs, pushes, merges, or finalizes — the
// §14 human brake is enforced here by only ever writing candidate prompts.
func Tick(deck string, cfg Config, signals []Candidate, now time.Time) (TickResult, error) {
	res := TickResult{Enabled: cfg.Enabled}
	if !cfg.Enabled {
		return res, nil
	}
	for _, sig := range signals {
		slug := SlugFor(sig)
		dir := filepath.Join(deck, "ideas", slug)
		if _, err := os.Stat(dir); err == nil {
			res.Skipped = append(res.Skipped, slug) // dedupe: candidate already exists
			continue
		} else if !errors.Is(err, fs.ErrNotExist) {
			return res, err // fail closed on an unexpected stat error
		}
		if err := writeCandidate(deck, slug, sig, now); err != nil {
			return res, err
		}
		res.Created = append(res.Created, slug)
	}
	return res, nil
}

// writeCandidate writes ideas/<slug>/00-prompt.md as a non-active candidate (§14 /
// §12.11 shape): status: candidate, NO `participants:` claim, a `## Promotion` note.
func writeCandidate(deck, slug string, c Candidate, now time.Time) error {
	dir := filepath.Join(deck, "ideas", slug)
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return err
	}
	title := strings.TrimSpace(c.Title)
	if title == "" {
		title = c.Source + " signal " + c.ID
	}
	detail := strings.TrimSpace(c.Detail)
	if detail == "" {
		detail = "(no detail provided)"
	}
	prompt := fmt.Sprintf(`---
idea: %s
author: loop
created: %s
status: candidate
source: %s
source_id: %s
fingerprint: %s
---

## Problem / idea

Auto-drafted by `+"`parley loop tick`"+` (COOPERATION.md §14, the outer-loop human brake)
from a discovered %s signal.

- title: %s
- source: %s
- source id: %s
- detail: %s

## Promotion

This is a non-active CANDIDATE (status: candidate), NOT an open round-01 idea. The automated
loop may discover and draft only (§14); it does not staff a quorum, run a deliberation,
implement, land, merge, or finalize. Before any deliberation starts, a human or the manifest
sets `+"`participants:`"+` (at least one non-facilitator participant) and flips status to
round-01 (the non-solo Phase-0 invariant).

## Constraints
## Non-goals
`, slug, now.Format("2006-01-02"), c.Source, c.ID, fingerprintOf(c),
		c.Source, title, c.Source, c.ID, detail)
	return os.WriteFile(filepath.Join(dir, "00-prompt.md"), []byte(prompt), 0o644)
}

// DeckDir is re-exported for callers that resolve the deck path.
const DeckDir = protocol.DeckDir
