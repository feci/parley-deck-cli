package track

import "testing"

func TestNormalize(t *testing.T) {
	cases := []struct {
		raw     string
		want    Track
		present bool
	}{
		{"fast", Fast, true},
		{"standard", Standard, true},
		{"deliberation", Deliberation, true},
		{"FAST", Fast, true},
		{"  `standard` ", Standard, true},
		{`"deliberation"`, Deliberation, true},
		{"", Standard, false},
		{"speedy", Standard, false},
		{"fast | standard | deliberation", Standard, false}, // the literal template placeholder
	}
	for _, c := range cases {
		got, present := Normalize(c.raw)
		if got != c.want || present != c.present {
			t.Errorf("Normalize(%q) = (%q,%v), want (%q,%v)", c.raw, got, present, c.want, c.present)
		}
	}
}

func TestClassifyDeliberationFirst(t *testing.T) {
	// Every deliberation trigger must win even when the size looks fast.
	small := Inputs{Files: 1, LOC: 10, Reversible: true, MechVerifiable: true}
	triggers := []func(*Inputs){
		func(i *Inputs) { i.ProtocolChange = true },
		func(i *Inputs) { i.Security = true },
		func(i *Inputs) { i.Irreversible = true },
		func(i *Inputs) { i.DataMigration = true },
		func(i *Inputs) { i.StrictGate = true },
		func(i *Inputs) { i.AutoImplement = true },
		func(i *Inputs) { i.Pipeline = true },
		func(i *Inputs) { i.APIBreak = true },
		func(i *Inputs) { i.SchemaBreak = true },
	}
	for idx, apply := range triggers {
		in := small
		apply(&in)
		if got, _ := Classify(in); got != Deliberation {
			t.Errorf("trigger %d: Classify = %q, want deliberation", idx, got)
		}
	}
	if got, _ := Classify(Inputs{Files: 16, Reversible: true, MechVerifiable: true}); got != Deliberation {
		t.Errorf(">15 files: got %q, want deliberation", got)
	}
	if got, _ := Classify(Inputs{LOC: 1001, Reversible: true, MechVerifiable: true}); got != Deliberation {
		t.Errorf(">1000 LOC: got %q, want deliberation", got)
	}
}

func TestClassifyFastAndStandard(t *testing.T) {
	fastable := Inputs{Files: 5, LOC: 300, FilesKnown: true, LOCKnown: true, Reversible: true, MechVerifiable: true}
	if got, _ := Classify(fastable); got != Fast {
		t.Errorf("boundary fast inputs: got %q, want fast", got)
	}
	// Fail-safe: unknown reversibility / verifiability keeps it out of fast.
	if got, _ := Classify(Inputs{Files: 2, LOC: 20, FilesKnown: true, LOCKnown: true, Reversible: false, MechVerifiable: true}); got != Standard {
		t.Errorf("not-known-reversible: got %q, want standard", got)
	}
	if got, _ := Classify(Inputs{Files: 2, LOC: 20, FilesKnown: true, LOCKnown: true, Reversible: true, MechVerifiable: false}); got != Standard {
		t.Errorf("not-mech-verifiable: got %q, want standard", got)
	}
	// Fail-safe: UNKNOWN size is never fast (review-01 F2).
	if got, _ := Classify(Inputs{Reversible: true, MechVerifiable: true}); got != Standard {
		t.Errorf("unknown size: got %q, want standard", got)
	}
	if got, _ := Classify(Inputs{Files: 1, LOC: 10, FilesKnown: true, Reversible: true, MechVerifiable: true}); got != Standard {
		t.Errorf("LOC unknown: got %q, want standard", got)
	}
	// Fail-safe: negative counts are never fast.
	if got, _ := Classify(Inputs{Files: 1, LOC: -1, FilesKnown: true, LOCKnown: true, Reversible: true, MechVerifiable: true}); got != Standard {
		t.Errorf("negative LOC: got %q, want standard", got)
	}
	// Boundary bands fall to standard.
	if got, _ := Classify(Inputs{Files: 10, LOC: 500, FilesKnown: true, LOCKnown: true, Reversible: true, MechVerifiable: true}); got != Standard {
		t.Errorf("6-15 file band: got %q, want standard", got)
	}
}

func TestPolicyForDeliberationNonSolo(t *testing.T) {
	// review-01 F1: explicit deliberation with no independent reviewer must error.
	if _, err := PolicyFor(Deliberation, true, 0, false, false); err == nil {
		t.Error("explicit deliberation with 0 reviewers must error (non-solo)")
	}
	// But absent-track (legacy) must NOT error (preflight enforces non-solo there).
	if _, err := PolicyFor(Standard, false, 0, false, false); err != nil {
		t.Errorf("absent track must not non-solo-error (legacy path): %v", err)
	}
}

func TestPolicyForStandardCapsCrossReview(t *testing.T) {
	p, _ := PolicyFor(Standard, true, 3, false, false)
	if p.CapCrossReviewRounds != 2 {
		t.Errorf("standard CapCrossReviewRounds = %d, want 2", p.CapCrossReviewRounds)
	}
}

func TestPolicyForAbsentIsLegacy(t *testing.T) {
	p, err := PolicyFor(Standard, false, 3, false, false)
	if err != nil {
		t.Fatalf("absent track error: %v", err)
	}
	if p.ApplyOverrides {
		t.Error("absent track must NOT apply overrides (legacy behaviour)")
	}
	if p.CrossReviewRounds != -1 || p.MaxReviewers != 0 || p.MinReviewers != 0 || p.MaxFixupCycles != 0 {
		t.Errorf("absent track must leave all knobs untouched, got %+v", p)
	}
}

func TestPolicyForFast(t *testing.T) {
	p, err := PolicyFor(Fast, true, 3, false, false)
	if err != nil {
		t.Fatalf("fast policy error: %v", err)
	}
	if !p.ApplyOverrides || p.MaxReviewers != 1 || p.MinReviewers != 1 || p.CrossReviewRounds != 0 || p.MaxFixupCycles != 1 {
		t.Errorf("fast policy = %+v, want overrides 1/1/0/1", p)
	}
	// Contradictions: fast + auto_implement / strict_gate must error.
	if _, err := PolicyFor(Fast, true, 3, true, false); err == nil {
		t.Error("fast + auto_implement must error")
	}
	if _, err := PolicyFor(Fast, true, 3, false, true); err == nil {
		t.Error("fast + strict_gate must error")
	}
	// Non-solo: fast with 0 available reviewers must error.
	if _, err := PolicyFor(Fast, true, 0, false, false); err == nil {
		t.Error("fast with 0 reviewers must error (non-solo)")
	}
}

func TestPolicyForStandard(t *testing.T) {
	p, err := PolicyFor(Standard, true, 3, false, false)
	if err != nil {
		t.Fatalf("standard policy error: %v", err)
	}
	if !p.ApplyOverrides || p.MaxReviewers != 2 || p.MinReviewers != 2 || p.MaxFixupCycles != 2 {
		t.Errorf("standard policy = %+v, want overrides 2/2/-/2", p)
	}
	if p.CrossReviewRounds != -1 {
		t.Errorf("standard must leave CrossReviewRounds (-1), got %d", p.CrossReviewRounds)
	}
	// Two-participant degradation: 1 available reviewer → MinReviewers 1.
	deg, err := PolicyFor(Standard, true, 1, false, false)
	if err != nil {
		t.Fatalf("standard degradation error: %v", err)
	}
	if deg.MinReviewers != 1 {
		t.Errorf("standard with 1 reviewer must degrade MinReviewers to 1, got %d", deg.MinReviewers)
	}
	// Non-solo.
	if _, err := PolicyFor(Standard, true, 0, false, false); err == nil {
		t.Error("standard with 0 reviewers must error (non-solo)")
	}
}

func TestPolicyForDeliberationIsLegacy(t *testing.T) {
	p, err := PolicyFor(Deliberation, true, 3, true, true)
	if err != nil {
		t.Fatalf("deliberation policy error: %v", err)
	}
	if p.ApplyOverrides {
		t.Error("deliberation must NOT override (it equals today's full lifecycle)")
	}
	if p.Track != Deliberation {
		t.Errorf("deliberation track = %q", p.Track)
	}
}
