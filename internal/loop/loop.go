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

// dedupeDigest is the collision-resistant identity of a signal (AF2): a 32-hex
// (128-bit) sha256 prefix over an UNAMBIGUOUS canonical key. strconv.Quote makes the
// field boundaries unforgeable, so `a/b` vs `a:b` (lossy-sanitize collision) and
// `"ci:"+"build"` vs `"ci"+":build"` (separator-shift collision) produce distinct
// digests, and a non-ASCII / emoji fingerprint hashes to a stable value instead of
// collapsing to a shared slug. An explicit fingerprint is the dedupe key when
// present; otherwise source+id is.
//
// 128 bits, not the original 8 hex / 32 bits (AF9): a 32-bit suffix was collidable
// by ordinary birthday search — a malicious signal could pick a fingerprint that
// collided with another's slug and silently suppress that distinct candidate. A
// 128-bit prefix makes a deliberate second-preimage collision infeasible.
func dedupeDigest(c Candidate) string {
	var key string
	if fp := strings.TrimSpace(c.Fingerprint); fp != "" {
		key = "fp:" + strconv.Quote(fp)
	} else {
		key = "si:" + strconv.Quote(c.Source) + strconv.Quote(c.ID)
	}
	sum := sha256.Sum256([]byte(key))
	return hex.EncodeToString(sum[:])[:32]
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
// prompt's FRONTMATTER (AF1): C0 control characters (0x00–0x1F) — newlines above all —
// are flattened to spaces so an untrusted signal cannot inject extra YAML frontmatter
// keys (e.g. a `participants:` quorum claim or a second `status:`). It also flattens
// the Unicode line/paragraph/next-line separators U+2028/U+2029/U+0085 (AF8): the
// repo's current line scanners split only on `\n`, so these are not a live bypass
// today, but they ARE YAML 1.1 line breaks — flattening them keeps the contract
// ("no line break can inject a key") true even if a real YAML parser is adopted later.
func cleanField(s string) string {
	return strings.Map(func(r rune) rune {
		if r == '\n' || r == '\r' || r == '\t' || r < 0x20 ||
			r == '\u2028' || r == '\u2029' || r == '\u0085' {
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

// safeMkdir ensures `dir` exists as a REAL directory without following a symlink at
// `dir` itself (AF10/AF14). os.Mkdir never creates through a symlink, so on ErrExist we
// Lstat (which does not follow links) and reject a symlink or non-directory. An existing
// real directory (e.g. an empty dir healed per AF7) is accepted.
func safeMkdir(dir string) error {
	if err := os.Mkdir(dir, 0o755); err != nil {
		if !errors.Is(err, fs.ErrExist) {
			return err
		}
		fi, lerr := os.Lstat(dir)
		if lerr != nil {
			return lerr
		}
		if fi.Mode()&os.ModeSymlink != 0 || !fi.IsDir() {
			return fmt.Errorf("loop: refusing %s: a symlink or non-directory", dir)
		}
	}
	return nil
}

// assertInsideDeck rejects a slug dir that resolves — through any symlink at any path
// component — to a location outside the deck (AF14). It is the depth-complete companion
// to safeMkdir's per-level leaf guards.
func assertInsideDeck(deck, dir string) error {
	realDeck, err := filepath.EvalSymlinks(deck)
	if err != nil {
		return err
	}
	realDir, err := filepath.EvalSymlinks(dir)
	if err != nil {
		return err
	}
	rel, err := filepath.Rel(realDeck, realDir)
	if err != nil || rel == ".." || strings.HasPrefix(rel, ".."+string(filepath.Separator)) {
		return fmt.Errorf("loop: refusing %s: resolves outside the deck %s", realDir, realDeck)
	}
	return nil
}

// writeCandidate writes ideas/<slug>/00-prompt.md as a non-active candidate (§14 /
// §12.11 shape): status: candidate, NO `participants:` claim, a `## Promotion` note.
// It returns created=false (no error) when the candidate already exists.
//
// AF7: the atomic claim is the PROMPT FILE, not the directory. An O_CREATE|O_EXCL
// open of 00-prompt.md serializes concurrent ticks (exactly one create, the rest
// skip), and — unlike a directory claim — a previously-crashed tick that left an
// empty ideas/<slug>/ no longer suppresses the signal forever (the absent file is
// re-created). A failed write removes the partial file so the next tick retries.
//
// AF1: the frontmatter scalars (Source, ID) and the one-line Title bullet go through
// cleanField so an untrusted signal cannot inject extra YAML keys. AF6: Detail is the
// free-form body field — it is rendered as an indented block with its newlines
// PRESERVED (a multi-line stack trace / log stays readable); the closed frontmatter
// above it means body newlines cannot inject a frontmatter key.
func writeCandidate(deck, slug string, c Candidate, now time.Time) (bool, error) {
	// The deck root is the user's trusted workspace; create it permissively.
	if err := os.MkdirAll(deck, 0o755); err != nil {
		return false, err
	}
	// AF10/AF14: a loop must only ever write inside parley-deck/ideas/<slug>/. Guard
	// BOTH the `ideas/` parent and the `ideas/<slug>` leaf against a planted symlink —
	// a symlink at EITHER would otherwise let the prompt land outside the deck.
	ideasDir := filepath.Join(deck, "ideas")
	if err := safeMkdir(ideasDir); err != nil {
		return false, err
	}
	dir := filepath.Join(ideasDir, slug)
	if err := safeMkdir(dir); err != nil {
		return false, err
	}
	// AF14: belt-and-suspenders containment — the slug dir, resolved through every
	// symlink, must stay inside the resolved deck. The per-level safeMkdir guards above
	// catch a symlink at ideas/ or ideas/<slug>; this one closes the whole class at ANY
	// ancestor depth in a single place (e.g. a future symlink higher up the path).
	if err := assertInsideDeck(deck, dir); err != nil {
		return false, err
	}
	promptPath := filepath.Join(dir, "00-prompt.md")
	f, err := os.OpenFile(promptPath, os.O_WRONLY|os.O_CREATE|os.O_EXCL, 0o644)
	if err != nil {
		if errors.Is(err, fs.ErrExist) {
			return false, nil // candidate already drafted → skip (dedupe)
		}
		return false, err
	}
	wrote := false
	defer func() {
		f.Close()
		if !wrote {
			os.Remove(promptPath) // AF7: never leave a poisoned partial prompt behind
		}
	}()

	source := cleanField(c.Source)
	id := cleanField(c.ID)
	title := cleanField(c.Title)
	if title == "" {
		title = source + " signal " + id
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

## Signal detail

%s

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
		source, title, source, id, indentDetail(c.Detail))

	if _, err := f.WriteString(prompt); err != nil {
		return false, err
	}
	wrote = true
	return true, nil
}

// indentDetail renders the free-form Detail as a markdown indented code block with
// newlines preserved (AF6), so multi-line provenance (a stack trace, a log excerpt)
// stays readable. Every line carries a four-space prefix, so no Detail content can
// reach column 0 to masquerade as a heading, a `---` fence, or a frontmatter key.
func indentDetail(s string) string {
	// AF11/AF15: normalize EVERY line-break-like separator to "\n" before splitting, so
	// each logical line gets the four-space prefix under ANY line splitter (LF scanners,
	// CommonMark, Python splitlines): CR/CRLF, vertical tab, form feed, the C0 info
	// separators U+001C/U+001D/U+001E, and the Unicode/YAML breaks U+0085/U+2028/U+2029.
	// Other C0 controls become spaces (they carry no line/structure semantics). \t is kept.
	s = strings.ReplaceAll(s, "\r\n", "\n")
	var b strings.Builder
	for _, r := range s {
		switch {
		case r == '\n' || r == '\r' || r == '\v' || r == '\f' ||
			r == 0x1c || r == 0x1d || r == 0x1e ||
			r == 0x85 || r == 0x2028 || r == 0x2029:
			b.WriteByte('\n')
		case r == '\t':
			b.WriteByte('\t')
		case r < 0x20:
			b.WriteByte(' ')
		default:
			b.WriteRune(r)
		}
	}
	out := strings.TrimSpace(b.String()) // AF12: drop leading/trailing blank lines
	if out == "" {
		return "    (no detail provided)"
	}
	lines := strings.Split(out, "\n")
	for i := range lines {
		lines[i] = "    " + lines[i]
	}
	return strings.Join(lines, "\n")
}

// DeckDir is re-exported for callers that resolve the deck path.
const DeckDir = protocol.DeckDir
