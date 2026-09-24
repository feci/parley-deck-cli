package protocolpacket

// lean-organizer C tests: the audience dimension. Hard safety property — the
// audience may only omit blocks already omittable and never below the never-cut
// floor; unknown audiences and facilitator_participates fall back to full with a
// stated reason; the attestation gains an additive `audience` field only.

import (
	"encoding/json"
	"strings"
	"testing"
)

const audienceSource = `# COOPERATION.md — Multi-Agent Cooperation Protocol

**Workspace:** ` + "`<workspace-name>`" + `

## Quickstart — start here

quickstart body

## 0. Choose the transport

transport body

## 1. Scope and purpose

scope body

### Non-solo execution requirement

non-solo body

## 4. Protocol — phases of an idea

phases body

### Phase 1 — Round 1

round-1 body

## 5. Quorum and async participation

quorum body

## 11. Transport mechanics

transport mechanics body

### 11.A — Local directory

local body

### 11.B — GitHub Pull Requests

github body

## 15. Verification integrity

fifteen body

### 15.1 Scope, ownership, location

15.1 body

### 15.7 Per-track binding

15.7 body
`

const audienceMap = `schema: parley.packet-applicability/v1
source: COOPERATION.md
blocks:
  - locator: "# COOPERATION.md — Multi-Agent Cooperation Protocol"
    include: always
  - locator: "## Quickstart — start here"
    include: always
  - locator: "## 0. Choose the transport"
    include: always
  - locator: "## 1. Scope and purpose"
    include: always
  - locator: "### Non-solo execution requirement"
    include: always
  - locator: "## 4. Protocol — phases of an idea"
    include: always
  - locator: "### Phase 1 — Round 1"
    include: when
    phases: [1]
    trigger: "writing a round-01 artifact"
  - locator: "## 5. Quorum and async participation"
    include: always
  - locator: "## 11. Transport mechanics"
    include: always
  - locator: "### 11.A — Local directory"
    include: when
    transports: [local-dir]
    trigger: "Transport: local-dir"
  - locator: "### 11.B — GitHub Pull Requests"
    include: when
    transports: [github-pr]
    trigger: "Transport: github-pr"
  - locator: "## 15. Verification integrity"
    include: when
    phases: [1, 2, 3, 5, 6, 7, 8]
    trigger: "any verification verdict"
  - locator: "### 15.1 Scope, ownership, location"
    include: when
    phases: [1, 2, 3, 5, 6, 7, 8]
    trigger: "any verification verdict"
  - locator: "### 15.7 Per-track binding"
    include: when
    phases: [1, 2, 3, 5, 6, 7, 8]
    trigger: "any verification verdict"
audiences:
  facilitator:
    omit:
      - locator: "## 1. Scope and purpose"
        trigger: "briefing a participant on §1 duties"
      - locator: "### 11.A — Local directory"
        trigger: "deck header names another transport"
      - locator: "### 11.B — GitHub Pull Requests"
        trigger: "deck header names another transport"
`

func audienceFixture(t *testing.T) (Source, *Map) {
	t.Helper()
	src := Source{Path: "COOPERATION.md", Role: "source", Authority: "test", Raw: audienceSource, Transport: "github-pr"}
	m, err := ParseMap(audienceMap)
	if err != nil {
		t.Fatalf("ParseMap: %v", err)
	}
	return src, m
}

func facilitatorBody(t *testing.T, src Source, m *Map, phase int) string {
	t.Helper()
	c := Build(src, m, Request{Phase: phase, Track: "deliberation", Transport: "github-pr", Audience: "facilitator"})
	if c.ContextMode != ModePacket {
		t.Fatalf("facilitator build fell back: %s (%s)", c.ContextMode, c.FallbackReason)
	}
	return c.Body
}

// headingVerbatim reports whether a heading line is reproduced VERBATIM in the
// packet body (line-anchored). The omission index necessarily names omitted blocks
// inside table rows — `^Heading$` matches only the verbatim reproduction, so the
// named-omission-set assertion cannot be fooled by the index itself.
func headingVerbatim(body, heading string) bool {
	for _, line := range strings.Split(body, "\n") {
		if strings.TrimSpace(line) == heading {
			return true
		}
	}
	return false
}

func TestAudienceOmitsNamedSetKeepsRetentionVerbatim(t *testing.T) {
	src, m := audienceFixture(t)
	body := facilitatorBody(t, src, m, 1)
	// Named omission set absent (the GATING check, R-2): the heading is not
	// reproduced verbatim anywhere and the block's body text is gone entirely.
	for _, bad := range []struct{ heading, text string }{
		{"## 1. Scope and purpose", "scope body"},
		{"### 11.A — Local directory", "local body"},
	} {
		if headingVerbatim(body, bad.heading) {
			t.Errorf("facilitator body must not reproduce %q verbatim", bad.heading)
		}
		if strings.Contains(body, bad.text) {
			t.Errorf("facilitator body must omit the body text of %q", bad.heading)
		}
	}
	// Retention set verbatim present (line-anchored headings + body text).
	for _, want := range []struct{ heading, text string }{
		{"## Quickstart — start here", "quickstart body"},
		{"## 0. Choose the transport", "transport body"},
		{"### Non-solo execution requirement", "non-solo body"},
		{"## 4. Protocol — phases of an idea", "phases body"},
		{"### Phase 1 — Round 1", "round-1 body"},
		{"## 5. Quorum and async participation", "quorum body"},
		{"## 11. Transport mechanics", "transport mechanics body"},
		{"### 11.B — GitHub Pull Requests", "github body"},
		{"## 15. Verification integrity", "fifteen body"},
		{"### 15.1 Scope, ownership, location", "15.1 body"},
		{"### 15.7 Per-track binding", "15.7 body"},
	} {
		if !headingVerbatim(body, want.heading) || !strings.Contains(body, want.text) {
			t.Errorf("facilitator body must retain %q verbatim", want.heading)
		}
	}
	// Complete omission index naming every omitted block.
	if !strings.Contains(body, "## Packet omission index") {
		t.Errorf("body must carry the complete omission index")
	}
	for _, omitted := range []string{"`## 1. Scope and purpose`", "`### 11.A — Local directory`"} {
		if !strings.Contains(body, omitted) {
			t.Errorf("omission index must list %s", omitted)
		}
	}
}

func TestAudienceNeverOmitsNeverCutAtKernelPhases(t *testing.T) {
	// §15 is never-cut at kernel phases {1,2,3,5,6,7,8}; a hostile map naming it
	// must not be able to cut it at phase 1.
	// Hostile copy of the full map that additionally omits §15 for facilitators.
	hostile := audienceMap + "      - locator: \"## 15. Verification integrity\"\n        trigger: \"hostile attempt to cut the floor\"\n"
	hm, err := ParseMap(hostile)
	if err != nil {
		t.Fatalf("ParseMap hostile: %v", err)
	}
	src := Source{Path: "COOPERATION.md", Role: "source", Authority: "test", Raw: audienceSource, Transport: "github-pr"}
	_ = src
	hsrc := Source{Path: "COOPERATION.md", Role: "source", Authority: "test", Raw: audienceSource, Transport: "github-pr"}
	body := facilitatorBody(t, hsrc, hm, 1)
	if !headingVerbatim(body, "## 15. Verification integrity") || !strings.Contains(body, "fifteen body") {
		t.Errorf("never-cut floor breached: §15 omitted at a kernel phase")
	}
}

func TestPacketCheckFailsAudienceOmittingAlwaysNeverCut(t *testing.T) {
	src, m := audienceFixture(t)
	// Rewrite the audience to omit an always-pinned never-cut block.
	bad := strings.Replace(audienceMap, `      - locator: "## 1. Scope and purpose"
        trigger: "briefing a participant on §1 duties"`, `      - locator: "## 0. Choose the transport"
        trigger: "hostile"`, 1)
	bm, err := ParseMap(bad)
	if err != nil {
		t.Fatalf("ParseMap bad: %v", err)
	}
	rep := Check(src, bm)
	if rep.OK || len(rep.NeverCut) == 0 {
		t.Fatalf("packet check must fail an audience omitting a never-cut block; ok=%v neverCut=%v", rep.OK, rep.NeverCut)
	}
	found := false
	for _, v := range rep.NeverCut {
		if strings.Contains(v, "audiences.facilitator omits never-cut block") {
			found = true
		}
	}
	if !found {
		t.Errorf("never-cut violation must name the audience rule; got %v", rep.NeverCut)
	}
	// The clean map passes, audiences included.
	if rep := Check(src, m); !rep.OK {
		t.Errorf("clean map with audiences must pass check: %+v", rep)
	}
}

// TestPacketCheckFailsAudienceOmittingPhasePinnedNeverCut (fix-up F6, claude-1
// MAJ-6): the FINAL C.1 proof obligation — hostile `audiences.facilitator.omit`
// entries naming the phase-pinned never-cut §15 family must FAIL `packet check` at
// map level (they previously passed: only `always` blocks were rejected, so the
// negative test could not fire for exactly the section FINAL singles out). The
// transport-conditional §11 subsections remain nameable (the ratified facilitator
// set omits the non-active ones).
func TestPacketCheckFailsAudienceOmittingPhasePinnedNeverCut(t *testing.T) {
	src, _ := audienceFixture(t)
	hostile := audienceMap
	for _, locator := range []string{
		"## 15. Verification integrity",
		"### 15.1 Scope, ownership, location",
		"### 15.7 Per-track binding",
	} {
		hostile += "      - locator: \"" + locator + "\"\n        trigger: \"hostile attempt to cut the phase-pinned floor\"\n"
	}
	hm, err := ParseMap(hostile)
	if err != nil {
		t.Fatalf("ParseMap hostile: %v", err)
	}
	rep := Check(src, hm)
	if rep.OK || len(rep.NeverCut) == 0 {
		t.Fatalf("packet check must fail an audience naming phase-pinned never-cut blocks; ok=%v neverCut=%v", rep.OK, rep.NeverCut)
	}
	for _, locator := range []string{
		"## 15. Verification integrity",
		"### 15.1 Scope, ownership, location",
		"### 15.7 Per-track binding",
	} {
		found := false
		for _, v := range rep.NeverCut {
			if strings.Contains(v, locator) {
				found = true
			}
		}
		if !found {
			t.Errorf("never-cut violation must name %q; got %v", locator, rep.NeverCut)
		}
	}
}

// TestPacketCheckAllowsTransportConditionalOmit (fix-up F6): the narrowed allowance —
// a non-active §11 subsection (transport-conditional never-cut) may still be NAMED by
// an audience without a map-level violation.
func TestPacketCheckAllowsTransportConditionalOmit(t *testing.T) {
	src, m := audienceFixture(t)
	rep := Check(src, m)
	if !rep.OK {
		t.Fatalf("the clean map omits non-active §11 subsections and must pass; got %+v", rep)
	}
	for _, v := range rep.NeverCut {
		if strings.Contains(v, "11.") {
			t.Errorf("transport-conditional §11 omission must not be a violation; got %q", v)
		}
	}
}

// TestOptimizeAudienceBranchSharesReasonsSet (fix-up F18, claude-1 NIT-3): the
// merged `--optimize`/audience branch claims "every guard that governs --optimize
// governs it identically" — pinned by asserting the two paths yield the IDENTICAL
// fallback-reason set for the same request shape across the kernel phases.
func TestOptimizeAudienceBranchSharesReasonsSet(t *testing.T) {
	src, m := audienceFixture(t)
	for _, phase := range []int{1, 5, 8} {
		opt := Build(src, m, Request{Phase: phase, Track: "deliberation", Transport: "github-pr", Optimize: true})
		aud := Build(src, m, Request{Phase: phase, Track: "deliberation", Transport: "github-pr", Optimize: true, Audience: "facilitator"})
		if strings.Join(opt.Problems, ";") != strings.Join(aud.Problems, ";") {
			t.Errorf("phase %d: --optimize and audience paths must produce the same reasons set; optimize=%v audience=%v", phase, opt.Problems, aud.Problems)
		}
	}
}

// TestOptimizeUnknownAudienceDoesNotStampRejectedAudience (fix-up F12, claude-1
// MIN-4): a body built under `--optimize` with an unrecognized audience must not
// stamp `audience=banana` into its header while the attestation says it fell back.
func TestOptimizeUnknownAudienceDoesNotStampRejectedAudience(t *testing.T) {
	src, m := audienceFixture(t)
	c := Build(src, m, Request{Phase: 1, Track: "deliberation", Transport: "github-pr", Optimize: true, Audience: "banana"})
	// The optimizer path still builds a participant-shaped packet; the rejected
	// audience is recorded as the audience-dimension fallback reason, and (F12) the
	// body header must carry the RESOLVED audience — "-" — never `banana`.
	if c.ContextMode != ModePacket {
		t.Fatalf("optimize with unknown audience still builds a packet, got %s", c.ContextMode)
	}
	if c.AudienceFallbackReason != "unknown-audience:banana" {
		t.Errorf("fallback reason must name the unknown audience, got %q", c.AudienceFallbackReason)
	}
	if c.Audience != "" {
		t.Errorf("resolved audience must stay empty on fallback, got %q", c.Audience)
	}
	if strings.Contains(c.Body, "audience=banana") {
		t.Errorf("the body header must carry the RESOLVED audience, never the rejected one; body starts:\n%s", firstLine(c.Body))
	}
}

func firstLine(s string) string {
	if i := strings.IndexByte(s, '\n'); i >= 0 {
		return s[:i]
	}
	return s
}

func TestUnknownAudienceFallsBackToFullWithReason(t *testing.T) {
	src, m := audienceFixture(t)
	c := Build(src, m, Request{Phase: 1, Track: "deliberation", Transport: "github-pr", Audience: "banana"})
	if c.ContextMode != ModeFull || c.Body != src.Raw {
		t.Fatalf("unknown audience must fall back to FULL context, got mode=%s", c.ContextMode)
	}
	if c.AudienceFallbackReason != "unknown-audience:banana" {
		t.Errorf("fallback reason must state the unknown audience, got %q", c.AudienceFallbackReason)
	}
	if c.Audience != "" {
		t.Errorf("resolved audience must stay empty on fallback, got %q", c.Audience)
	}
}

func TestFacilitatorParticipatesKeepsFullContext(t *testing.T) {
	src, m := audienceFixture(t)
	c := Build(src, m, Request{Phase: 1, Track: "deliberation", Transport: "github-pr", Audience: "facilitator", FacilitatorParticipates: true})
	if c.ContextMode != ModeFull || c.Body != src.Raw {
		t.Fatalf("facilitator_participates must keep full context, got mode=%s", c.ContextMode)
	}
	if c.AudienceFallbackReason != "facilitator-participates" {
		t.Errorf("fallback reason must name facilitator-participates, got %q", c.AudienceFallbackReason)
	}
}

func TestAttestationAudienceFieldAdditiveOnly(t *testing.T) {
	src, m := audienceFixture(t)
	c := Build(src, m, Request{Phase: 1, Track: "deliberation", Transport: "github-pr", Audience: "facilitator"})
	if c.Audience != "facilitator" {
		t.Errorf("attestation must carry the resolved audience, got %q", c.Audience)
	}
	raw, _ := json.Marshal(c.Attestation)
	if !strings.Contains(string(raw), `"audience":"facilitator"`) {
		t.Errorf("attestation JSON must carry the additive audience field; got %s", raw)
	}
	if strings.Contains(string(raw), `"role"`) {
		t.Errorf("attestation must NEVER carry a role key (protocolRole collision); got %s", raw)
	}
}
