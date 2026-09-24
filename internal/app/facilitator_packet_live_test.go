package app

// lean-organizer C.3 against the LIVE deck protocol + map: the named-omission-set
// absence is the GATING check (R-2); the facilitator body is recorded against BOTH
// the measured floor (the retained whole-block source total) and the <= 70,000 B
// guardrail in this same run; the --optimize baseline is captured here too
// (measured, not transcribed from any round file). Fix-up G1 (claude-1 R2-MAJ-1):
// the floor is what the body RETAINS; the omission-set total is a separate
// quantity under its own label — an omission sum is never called a floor.

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"parley-deck-cli/internal/protocolpacket"
)

const (
	facilitatorGuardrailBytes = 70_000
	deckProtocolRel           = "../../parley-deck/COOPERATION.md"
	deckMapRel                = "../../parley-deck/meta/packet-applicability.yaml"
)

// namedOmissionSet is the ratified facilitator omission set (FINAL C.2), top-level
// blocks, verified by heading-line absence below.
var namedOmissionSet = []string{
	"## 1. Scope and purpose",
	"## 3. Directory layout",
	"## 8. Inbox (lightweight channel)",
	"## 10. TL;DR",
	"## 12. Pipeline blocks & action stages",
	"## 13. Retrospective optimization",
	"## Appendix A — Adopting this protocol in a new project",
	"### 11.A — Local directory",
	"### 11.C — GitLab Merge Requests",
}

// namedRetentionSet is the ratified facilitator retention set (FINAL C.2).
var namedRetentionSet = []string{
	"## Quickstart — start here (developers & first-timers)",
	"## 2. Active agents (roster)",
	"## 4. Protocol — phases of an idea",
	"## 5. Quorum and async participation",
	"## 9. Session-start checklist for every agent",
	"### 11.B — GitHub Pull Requests",
	"## 15. Verification integrity",
}

func headingLinePresent(body, heading string) bool {
	for _, line := range strings.Split(body, "\n") {
		if strings.TrimSpace(line) == heading {
			return true
		}
	}
	return false
}

func TestLiveDeckFacilitatorPacketNamedSetsAndGuardrail(t *testing.T) {
	raw, err := os.ReadFile(deckProtocolRel)
	if err != nil {
		t.Fatalf("live deck protocol: %v", err)
	}
	mapRaw, err := os.ReadFile(deckMapRel)
	if err != nil {
		t.Fatalf("live deck map: %v", err)
	}
	m, err := protocolpacket.ParseMap(string(mapRaw))
	if err != nil {
		t.Fatalf("ParseMap: %v", err)
	}
	src := protocolpacket.Source{Path: deckProtocolRel, Role: "source", Authority: "test", Raw: string(raw), Transport: "github-pr"}

	// GATING check (R-2): each named-omission-set block verified absent — a body
	// under 70,000 B that retains any of them fails C.3 regardless of size.
	req := protocolpacket.Request{Phase: 1, Track: "deliberation", Transport: "github-pr", Audience: "facilitator"}
	c := protocolpacket.Build(src, m, req)
	if c.ContextMode != protocolpacket.ModePacket {
		t.Fatalf("facilitator build fell back: %s (%s)", c.ContextMode, c.FallbackReason)
	}
	for _, heading := range namedOmissionSet {
		if headingLinePresent(c.Body, heading) {
			t.Errorf("GATING: facilitator body retains named-omission-set block %q", heading)
		}
	}
	// Retention set verbatim.
	for _, heading := range namedRetentionSet {
		if !headingLinePresent(c.Body, heading) {
			t.Errorf("facilitator body must retain %q verbatim", heading)
		}
	}
	// Complete omission index with triggers.
	if !strings.Contains(c.Body, "## Packet omission index") {
		t.Errorf("facilitator body must carry the complete omission index")
	}

	// Guardrail, floor and omission total, all measured in this same run (R-2 /
	// open item 14; fix-up G1 — claude-1 R2-MAJ-1). The FLOOR is the retained
	// whole-block source total: the bytes the facilitator body itself carries for
	// the blocks this request INCLUDES, summed exactly as renderPacket lays them
	// out (each block's text plus its blank-line separator). Parse splits `###`
	// subsections into their own blocks, so every retained subsection is counted —
	// nothing is silently excluded. A floor is what the body RETAINS; the omission
	// set is the opposite quantity and is logged under its own label below.
	floor := 0
	for _, r := range c.Index {
		if !r.Included {
			continue
		}
		floor += len(strings.TrimRight(r.Text, "\n")) + 2
	}
	if len(c.Body) < floor {
		t.Errorf("facilitator body %d B is below its own retained-source floor %d B — the renderer dropped retained bytes", len(c.Body), floor)
	}

	// The named-omission-set total, whole-section-correct (fix-up G1): a
	// `##`-level section's span runs from its heading block through every deeper
	// block below it (its `###`/`####` subsections — Parse splits them out, so a
	// byHeading lookup alone silently drops them; claude-1 measured 17 blocks /
	// 11,903 B of §12/§13 subsection content the old sum missed). Convention: each
	// block contributes len(Text)+1 (its text plus the newline joining it to the
	// next block of the section).
	blocks := protocolpacket.Parse(string(raw))
	omissionWholeSections := 0
	for _, heading := range namedOmissionSet {
		for i, b := range blocks {
			if b.Locator != heading {
				continue
			}
			span := len(b.Text) + 1
			for j := i + 1; j < len(blocks) && blocks[j].Level > b.Level; j++ {
				span += len(blocks[j].Text) + 1
			}
			omissionWholeSections += span
			break
		}
	}

	// --optimize baseline (participant audience) captured in the same run.
	base := protocolpacket.Build(src, m, protocolpacket.Request{Phase: 1, Track: "deliberation", Transport: "github-pr", Optimize: true, Audience: "participant"})
	optimizeBytes := len(base.Body)
	if base.ContextMode == protocolpacket.ModePacket && base.Shadow != nil {
		optimizeBytes = base.Shadow.PacketBytes
	}

	t.Logf("R-2 measurement (phase 1 / deliberation / github-pr): facilitator body=%d B; floor (retained whole-block source total: sum over the request's included blocks of len(text)+2, the renderer's exact layout)=%d B; named-omission-set whole-section bytes (each block len(text)+1, subsections included)=%d B; guardrail=%d B; --optimize baseline=%d B",
		len(c.Body), floor, omissionWholeSections, facilitatorGuardrailBytes, optimizeBytes)

	if len(c.Body) > facilitatorGuardrailBytes {
		t.Errorf("facilitator body %d B exceeds the 70,000 B guardrail — bytes go back to the quorum before any map/ceiling change", len(c.Body))
	}

	// packet check stays green with the audiences key (structure incl. the audience
	// never-cut proof across every phase x track).
	rep := protocolpacket.Check(src, m)
	if !rep.OK {
		t.Errorf("packet check with audiences must stay green: %+v", rep)
	}
}

// TestLiveDeckFacilitatorAcrossPhases renders the facilitator packet at EVERY phase
// (deliberation/github-pr): no build may fall back (structural), and every body is
// measured and logged. The <= 70,000 B guardrail's ratified scope is the phase-1
// body (hard-asserted above); other phases are MEASURED here — a phase whose body
// exceeds the guardrail is shown to the quorum before any map/ceiling change (open
// item 2), never silently relaxed or silently asserted.
func TestLiveDeckFacilitatorAcrossPhases(t *testing.T) {
	raw, err := os.ReadFile(deckProtocolRel)
	if err != nil {
		t.Fatalf("live deck protocol: %v", err)
	}
	mapRaw, err := os.ReadFile(deckMapRel)
	if err != nil {
		t.Fatalf("live deck map: %v", err)
	}
	m, _ := protocolpacket.ParseMap(string(mapRaw))
	src := protocolpacket.Source{Path: deckProtocolRel, Role: "source", Authority: "test", Raw: string(raw), Transport: "github-pr"}
	for phase := 0; phase <= 8; phase++ {
		c := protocolpacket.Build(src, m, protocolpacket.Request{Phase: phase, Track: "deliberation", Transport: "github-pr", Audience: "facilitator"})
		if c.ContextMode != protocolpacket.ModePacket {
			t.Errorf("phase %d: facilitator build fell back: %s (%s)", phase, c.ContextMode, c.FallbackReason)
			continue
		}
		note := ""
		if len(c.Body) > facilitatorGuardrailBytes {
			note = " (above the 70,000 B guardrail — guardrail scope is phase 1; bytes recorded for the quorum per open item 2)"
		}
		t.Logf("phase %d facilitator body: %d B%s", phase, len(c.Body), note)
	}
}

var _ = filepath.Join
