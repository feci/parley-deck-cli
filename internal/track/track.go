// Package track implements the §4.0 conditional-rigor track model: a pure,
// deterministic classifier (objective inputs → track) and the per-track policy
// (reviewer caps, fix-up caps, cross-review) the driver enforces. It has no
// dependencies so it is trivially unit-testable and importable by both the
// driver and the CLI (`parley classify`).
//
// Design invariants (idea track-aware-driver, shipped after
// meta-protocol-change-devx-speed):
//   - Classification is deliberation-first, then fast, else standard, and fails
//     safe to the stricter track on doubt (§4.0).
//   - An ABSENT track reproduces today's behaviour byte-for-byte (ApplyOverrides
//     false); only an EXPLICIT track applies §4.0's reduced ceremony.
//   - No track may drop the non-solo (≥1 independent reviewer) invariant; fast is
//     rejected when combined with the deliberation triggers auto_implement /
//     strict_gate.
package track

import (
	"fmt"
	"strings"
)

// Track is one of the three §4.0 rigor tracks.
type Track string

const (
	Fast         Track = "fast"
	Standard     Track = "standard"
	Deliberation Track = "deliberation"
)

// Normalize maps a raw frontmatter value to a Track. It returns (track, true)
// for a recognized value and (Standard, false) for absent/empty/unknown — so an
// unrecognized or missing track fails safe to standard AND is reported as
// "not explicitly declared" (the caller uses that to preserve legacy behaviour).
func Normalize(raw string) (Track, bool) {
	v := strings.ToLower(strings.Trim(strings.TrimSpace(raw), "`'\"* "))
	switch v {
	case string(Fast):
		return Fast, true
	case string(Standard):
		return Standard, true
	case string(Deliberation):
		return Deliberation, true
	default:
		return Standard, false
	}
}

// Inputs are the objective signals the §4.0 classifier reads. Booleans that are
// not known should be left false (fail-safe: unknown reversibility/verifiability
// simply keeps a task out of fast).
type Inputs struct {
	Files          int
	LOC            int
	FilesKnown     bool // the file count was actually supplied (unknown size is never fast)
	LOCKnown       bool // the LOC count was actually supplied
	Reversible     bool // the change is known to be fully reversible
	MechVerifiable bool // the change is mechanically verifiable (lint/type/test)

	// Deliberation triggers (any one forces deliberation).
	ProtocolChange bool
	Security       bool // security/auth/secrets/payments/privacy/production-infra
	Irreversible   bool // irreversible/destructive op
	DataMigration  bool
	StrictGate     bool
	AutoImplement  bool
	Pipeline       bool // pipeline or action block
	APIBreak       bool // public-API break
	SchemaBreak    bool // persisted-schema break
}

// Classify implements §4.0: evaluate deliberation triggers first (first match
// wins); else fast only when ALL fast conditions hold; else standard. It returns
// the track and a short reason (the trigger that fired), suitable for `--json`.
func Classify(in Inputs) (Track, string) {
	for _, d := range []struct {
		on  bool
		why string
	}{
		{in.ProtocolChange, "protocol-change"},
		{in.Security, "security/secrets/production surface"},
		{in.Irreversible, "irreversible/destructive"},
		{in.DataMigration, "data-migration"},
		{in.StrictGate, "strict_gate"},
		{in.AutoImplement, "auto_implement"},
		{in.Pipeline, "pipeline/action block"},
		{in.APIBreak, "public-api-break"},
		{in.SchemaBreak, "persisted-schema-break"},
		{in.Files > 15, "more than 15 files"},
		{in.LOC > 1000, "more than 1000 LOC"},
	} {
		if d.on {
			return Deliberation, d.why
		}
	}
	// Fast requires ALL fast conditions to be POSITIVELY known. Unknown or invalid
	// (negative) size is not proof of "small" — it fails safe to standard (§4.0).
	if in.Reversible && in.MechVerifiable &&
		in.FilesKnown && in.Files >= 0 && in.Files <= 5 &&
		in.LOCKnown && in.LOC >= 0 && in.LOC <= 300 {
		return Fast, "all fast conditions met (reversible, mechanically verifiable, small)"
	}
	return Standard, "default (neither forced to deliberation nor fully fast)"
}

// Policy is the per-track driver configuration derived from a declared track.
// A zero/"leave" value means "do not change the driver default":
//   - MaxReviewers 0        → unlimited (all non-implementers review)
//   - MinReviewers 0        → use the driver's default (2)
//   - CrossReviewRounds -1  → leave as configured
//   - MaxFixupCycles 0      → leave the driver default
type Policy struct {
	Track                Track
	ApplyOverrides       bool // false = reproduce today's behaviour (absent/deliberation)
	MaxReviewers         int
	MinReviewers         int
	CrossReviewRounds    int // exact value to set; -1 = leave as configured
	CapCrossReviewRounds int // upper bound to clamp the configured value to; 0 = no cap
	MaxFixupCycles       int
}

// PolicyFor derives the per-track policy. `present` distinguishes an explicit
// track (apply §4.0) from an absent one (preserve legacy behaviour).
// availableReviewers is the number of distinct non-implementer participants; it
// drives the non-solo hard-reject and the §4.0 two-participant degradation.
// It returns an error for the §4.0 contradictions (fast + auto_implement, fast +
// strict_gate) and for a non-solo violation, so the caller can escalate.
func PolicyFor(t Track, present bool, availableReviewers int, autoImplement, strictGate bool) (Policy, error) {
	// Absent track → legacy behaviour, nothing overridden. (The non-solo floor for
	// legacy ideas is enforced at the app/preflight layer; this path preserves
	// today's driver behaviour byte-for-byte.)
	if !present {
		return Policy{Track: Standard, ApplyOverrides: false, CrossReviewRounds: -1}, nil
	}
	// Non-solo (§1) is an all-track invariant: every EXPLICIT track requires at
	// least one independent reviewer — including deliberation (review-01 fix).
	if availableReviewers < 1 {
		return Policy{}, fmt.Errorf("track: %s requires at least 1 independent reviewer (non-solo, §1); none available", t)
	}
	switch t {
	case Fast:
		if autoImplement {
			return Policy{}, fmt.Errorf("track: fast is invalid with auto_implement — auto_implement is a §4.0 deliberation trigger; use track: deliberation or remove auto_implement")
		}
		if strictGate {
			return Policy{}, fmt.Errorf("track: fast is invalid with strict_gate — strict_gate is a §4.0 deliberation trigger; use track: deliberation or remove strict_gate")
		}
		return Policy{Track: Fast, ApplyOverrides: true, MaxReviewers: 1, MinReviewers: 1, CrossReviewRounds: 0, MaxFixupCycles: 1}, nil
	case Deliberation:
		// Deliberation == today's full lifecycle (backward-compat constraint), but
		// still subject to the non-solo floor checked above.
		return Policy{Track: Deliberation, ApplyOverrides: false, CrossReviewRounds: -1}, nil
	default: // explicit Standard
		min := 2
		if availableReviewers <= 1 { // §4.0 two-participant degradation
			min = 1
		}
		return Policy{Track: Standard, ApplyOverrides: true, MaxReviewers: 2, MinReviewers: min, CrossReviewRounds: -1, CapCrossReviewRounds: 2, MaxFixupCycles: 2}, nil
	}
}
