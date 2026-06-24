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
	"strconv"
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
	Enabled  bool     `json:"enabled"`
	Created  []string `json:"created"`
	Skipped  []string `json:"skipped"`  // already-existing candidate ideas (dedupe)
	Rejected []string `json:"rejected"` // signals with an unknown source (AF1, fail closed)
}

// validSources is the closed set a signal's source must belong to (FINAL.md LE-9).
// An unknown source is rejected rather than normalized — a defense-in-depth layer
// on top of the frontmatter sanitization (AF1).
var validSources = map[string]bool{"commit": true, "ci": true, "issue": true, "manual": true}

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

// dedupeDigest is the collision-resistant identity of a signal (AF2): an 8-char
// sha256 over an UNAMBIGUOUS canonical key. strconv.Quote makes the field
// boundaries unforgeable, so `a/b` vs `a:b` (lossy-sanitize collision) and
// `"ci:"+"build"` vs `"ci"+":build"` (separator-shift collision) produce distinct
// digests, and a non-ASCII / emoji fingerprint hashes to a stable value instead of
// collapsing to a shared slug. An explicit fingerprint is the dedupe key when
// present; otherwise source+id is.
func dedupeDigest(c Candidate) string {
	var key string
	if fp := strings.TrimSpace(c.Fingerprint); fp != "" {
		key = "fp:" + strconv.Quote(fp)
	} else {
		key = "si:" + strconv.Quote(c.Source) + strconv.Quote(c.ID)
	}
	sum := sha256.Sum256([]byte(key))
	return hex.EncodeToString(sum[:])[:8]
}

// sanitize lowercases and keeps only [a-z0-9-] so a value is slug-safe. Used only
// for the human-readable source hint in the slug; identity comes from dedupeDigest.
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

// cleanField neutralizes a signal value before it is written into a candidate
// prompt (AF1): control characters — newlines above all — are flattened to spaces
// so an untrusted signal cannot inject extra YAML frontmatter keys (e.g. a
// `participants:` quorum claim or a second `status:`) or break the markdown body.
func cleanField(s string) string {
	return strings.Map(func(r rune) rune {
		if r == '\n' || r == '\r' || r == '\t' || r < 0x20 {
			return ' '
		}
		return r
	}, strings.TrimSpace(s))
}

// SlugFor is the deterministic candidate-idea slug: a readable source hint plus the
// collision-resistant dedupe digest. Source is validated to the closed set before a
// signal reaches here, so the hint cannot collide across distinct valid sources.
func SlugFor(c Candidate) string {
	src := sanitize(c.Source)
	if src == "" {
		src = "signal"
	}
	return "loop-" + src + "-" + dedupeDigest(c)
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
		if !validSources[strings.TrimSpace(sig.Source)] {
			// AF1: reject an unknown source rather than normalize it (fail closed).
			res.Rejected = append(res.Rejected, cleanField(sig.Source)+":"+cleanField(sig.ID))
			continue
		}
		slug := SlugFor(sig)
		created, err := writeCandidate(deck, slug, sig, now)
		if err != nil {
			return res, err
		}
		if created {
			res.Created = append(res.Created, slug)
		} else {
			res.Skipped = append(res.Skipped, slug) // dedupe: slug already claimed
		}
	}
	return res, nil
}

// writeCandidate writes ideas/<slug>/00-prompt.md as a non-active candidate (§14 /
// §12.11 shape): status: candidate, NO `participants:` claim, a `## Promotion` note.
// It returns created=false (no error) when the slug is already claimed. The claim is
// an atomic os.Mkdir + O_EXCL write (AF4), so two concurrent ticks cannot both create
// or clobber the same prompt. Every signal-derived value is run through cleanField
// (AF1) so it can never inject extra YAML frontmatter keys or break the body.
func writeCandidate(deck, slug string, c Candidate, now time.Time) (bool, error) {
	ideasDir := filepath.Join(deck, "ideas")
	if err := os.MkdirAll(ideasDir, 0o755); err != nil {
		return false, err
	}
	dir := filepath.Join(ideasDir, slug)
	if err := os.Mkdir(dir, 0o755); err != nil {
		if errors.Is(err, fs.ErrExist) {
			return false, nil // slug already claimed → skip (dedupe)
		}
		return false, err // fail closed on an unexpected mkdir error
	}

	source := cleanField(c.Source)
	id := cleanField(c.ID)
	title := cleanField(c.Title)
	if title == "" {
		title = source + " signal " + id
	}
	detail := cleanField(c.Detail)
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

(to be filled on promotion)

## Non-goals

(to be filled on promotion)
`, slug, now.Format("2006-01-02"), source, id, dedupeDigest(c),
		source, title, source, id, detail)

	// AF4: O_EXCL so a concurrent writer that claimed the dir cannot clobber the prompt.
	f, err := os.OpenFile(filepath.Join(dir, "00-prompt.md"), os.O_WRONLY|os.O_CREATE|os.O_EXCL, 0o644)
	if err != nil {
		return false, err
	}
	defer f.Close()
	if _, err := f.WriteString(prompt); err != nil {
		return false, err
	}
	return true, nil
}

// DeckDir is re-exported for callers that resolve the deck path.
const DeckDir = protocol.DeckDir
