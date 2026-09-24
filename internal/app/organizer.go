package app

// `parley organizer brief --idea <slug>` (lean-organizer C.5): a COMPUTED,
// never-stored re-orientation surface for the organizer. Content: the packet
// attestation + facilitator body path, the idea/run state, the current PhaseDigest,
// the driver's next action (fixed enumeration), and the phase pointer from the run
// record. Contract: <= 8,192 B, byte-identical across two runs over an unchanged
// tree, writes no file inside the deck — a stored brief is a rejected sidecar
// pattern; any drift toward persisting it is a stop moment.

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"parley-deck-cli/internal/driver"
	"parley-deck-cli/internal/protocol"
	"parley-deck-cli/internal/protocolpacket"
	"parley-deck-cli/internal/runstate"
)

const organizerBriefByteCap = 8192

// briefPhaseFromStatus maps an idea/implementation status to a protocol phase
// (fix-up F5). Covers the deck's real vocabulary: the 00-prompt `status:` stops at
// `final` for most ideas (final×80 in this deck), so Phases 5–8 are read from the
// IMPLEMENTATION.md status and the review-round presence instead.
func briefPhaseFromStatus(status string) (int, bool) {
	switch {
	case strings.HasPrefix(status, "round-01"):
		return 1, true
	case strings.HasPrefix(status, "round-"):
		return 2, true
	case status == "consensus":
		return 3, true
	case status == "final":
		return 4, true
	case status == "implementation", status == "implemented", status == "in-progress":
		return 5, true
	case status == "ready-for-review", strings.HasPrefix(status, "review"):
		return 6, true
	case strings.HasPrefix(status, "fix-up-cycle-"), status == "complete":
		return 8, true
	}
	return 0, false
}

// briefPhaseFromRun reads the driver run cursor (run.json `phase`), the live pointer
// while a driver run is advancing the idea.
func briefPhaseFromRun(root, slug string) (int, bool) {
	summary, err := runstate.ResolveRun(root, slug)
	if err != nil || summary.RunID == "" {
		return 0, false
	}
	data, err := os.ReadFile(filepath.Join(summary.RunDir, "run.json"))
	if err != nil {
		return 0, false
	}
	var manifest struct {
		Phase        string `json:"phase"`
		CurrentRound string `json:"current_round"`
	}
	if json.Unmarshal(data, &manifest) != nil || manifest.Phase == "" {
		return 0, false
	}
	switch manifest.Phase {
	case "round":
		if manifest.CurrentRound == "round-01" {
			return 1, true
		}
		return 2, true
	case "consensus":
		return 3, true
	case "final":
		return 4, true
	case "impl":
		return 5, true
	case "review":
		return 6, true
	case "done":
		return 8, true
	}
	return 0, false
}

// briefPhase derives the protocol phase pointer (fix-up F5, claude-1 MAJ-5): the
// MOST ADVANCED of the available signals — the IMPLEMENTATION.md status (the only
// artifact that advances through Phases 5–8), the driver run cursor, the presence of
// review rounds, and the 00-prompt status. Every signal only ever advances when its
// phase is actually reached, so the maximum is the live pointer and a stale signal
// (a run cursor left at `round`, a prompt status left at `final`) can never drag the
// brief back to the §15-free phase-0 view. A live Phase 5–8 idea therefore never
// resolves to the phase-0 packet.
func briefPhase(root, slug, promptStatus string, digest driver.PhaseDigest) int {
	phase := 0
	if p, ok := briefPhaseFromStatus(promptStatus); ok && p > phase {
		phase = p
	}
	if p, ok := briefPhaseFromRun(root, slug); ok && p > phase {
		phase = p
	}
	if digest.Review != nil && phase < 6 {
		phase = 6 // review rounds exist — the idea is at Phase 6 or beyond
	}
	if digest.Implementation != nil && digest.Implementation.Present {
		if p, ok := briefPhaseFromStatus(digest.Implementation.Status); ok && p > phase {
			phase = p
		}
	}
	return phase
}

func runOrganizer(args []string, stdout, stderr io.Writer) int {
	if len(args) > 0 && args[0] == "brief" {
		return runOrganizerBrief(args[1:], stdout, stderr)
	}
	if len(args) == 0 {
		return runOrganizerBrief(nil, stdout, stderr)
	}
	fmt.Fprintln(stderr, "usage: parley organizer brief --idea <slug> [--dir DIR]")
	return 1
}

func runOrganizerBrief(args []string, stdout, stderr io.Writer) int {
	fs := flag.NewFlagSet("organizer brief", flag.ContinueOnError)
	fs.SetOutput(stderr)
	dir := fs.String("dir", ".", "workspace root")
	idea := fs.String("idea", "", "idea slug")
	if err := fs.Parse(args); err != nil {
		return 1
	}
	if strings.TrimSpace(*idea) == "" {
		fmt.Fprintln(stderr, "organizer brief: --idea is required")
		return 1
	}
	root, err := filepath.Abs(*dir)
	if err != nil {
		fmt.Fprintf(stderr, "organizer brief: %v\n", err)
		return 1
	}
	ws, err := protocol.ReadWorkspaceStatus(root)
	if err != nil {
		fmt.Fprintf(stderr, "organizer brief: %v\n", err)
		return 1
	}
	var st protocol.IdeaStatus
	for _, i := range ws.Ideas {
		if i.Slug == *idea {
			st = i
			break
		}
	}
	if st.Slug == "" {
		fmt.Fprintf(stderr, "organizer brief: idea %q not found\n", *idea)
		return 1
	}

	track := "standard"
	if t, present := driver.ReadTrack(st.Path); present && t != "" {
		track = string(t)
	}
	// The digest feeds the phase pointer (review-round presence, implementation
	// status) and is computed before the packet render (fix-up F5).
	digest := driver.BuildPhaseDigest(root, st.Slug, st.Path, st.Participants)
	phase := briefPhase(root, st.Slug, st.Status, digest)

	// Protocol context: the facilitator packet (computed; body cached by the
	// renderer under .parley-runtime/, never inside the deck). A signing
	// facilitator (facilitator_participates: true) keeps full context.
	req := protocolpacket.Request{
		Phase: phase, Track: track, Transport: ws.Transport, IdeaSlug: st.Slug,
		Audience:                "facilitator",
		FacilitatorParticipates: st.FacilitatorRole.Declared && st.FacilitatorRole.Participates,
	}
	coreSt, cerr := coreStore()
	ctx, rerr := protocolpacket.Render(root, coreSt, req)
	if rerr != nil || cerr != nil {
		reason := rerr.Error()
		if cerr != nil {
			reason = cerr.Error()
		}
		ctx = protocolpacket.Context{Attestation: protocolpacket.Attestation{ContextMode: "full-fallback", FallbackReason: reason}}
	}

	runPhase := runPhasePointer(root, st.Slug)

	brief := renderOrganizerBrief(st, ws.Transport, ctx, digest, runPhase)
	if len(brief) > organizerBriefByteCap {
		fmt.Fprintf(stderr, "organizer brief: %d B exceeds the %d B cap — refusing to emit a truncated view\n", len(brief), organizerBriefByteCap)
		return 1
	}
	fmt.Fprint(stdout, brief)
	return 0
}

// runPhasePointer reads the run manifest's phase pointer (the D handoff state). It is
// advisory: the recomputed view is authoritative on any disagreement.
func runPhasePointer(root, slug string) string {
	summary, err := runstate.ResolveRun(root, slug)
	if err != nil || summary.RunID == "" {
		return "no-run"
	}
	data, err := os.ReadFile(filepath.Join(summary.RunDir, "run.json"))
	if err == nil {
		var manifest struct {
			Phase string `json:"phase"`
		}
		if json.Unmarshal(data, &manifest) == nil && manifest.Phase != "" {
			return manifest.Phase
		}
	}
	return summary.CurrentRound
}

func renderOrganizerBrief(st protocol.IdeaStatus, transport string, ctx protocolpacket.Context, digest driver.PhaseDigest, runPhase string) string {
	var b strings.Builder
	fmt.Fprintf(&b, "# Organizer brief — %s\n\n", st.Slug)
	b.WriteString("Computed view — never stored in the deck. Recompute with `parley organizer brief`;\n")
	b.WriteString("the brief itself writes nothing into `parley-deck/` (the packet renderer caches its\n")
	b.WriteString("body under `.parley-runtime/protocol-packets/`, outside the deck).\n")
	b.WriteString("`parley status` / this brief are authoritative; `runs/` handoff records are advisory.\n\n")

	fmt.Fprintf(&b, "## Protocol context\n")
	fmt.Fprintf(&b, "context_mode: %s", ctx.ContextMode)
	if ctx.Audience != "" {
		fmt.Fprintf(&b, "  audience: %s", ctx.Audience)
	}
	if ctx.AudienceFallbackReason != "" {
		fmt.Fprintf(&b, "  (full fallback: %s)", ctx.AudienceFallbackReason)
	}
	b.WriteString("\n")
	fmt.Fprintf(&b, "source_sha256: %s\n", ctx.SourceSHA256)
	if ctx.PacketSHA256 != "" {
		fmt.Fprintf(&b, "packet_sha256: %s\n", ctx.PacketSHA256)
	}
	if ctx.FallbackReason != "" {
		fmt.Fprintf(&b, "fallback_reason: %s\n", ctx.FallbackReason)
	}
	if ctx.BodyPath != "" {
		fmt.Fprintf(&b, "facilitator body: %s\n", ctx.BodyPath)
	}
	fmt.Fprintf(&b, "transport: %s\n\n", transport)

	fmt.Fprintf(&b, "## Idea state\n")
	fmt.Fprintf(&b, "status: %s  participants: [%s]\n", st.Status, strings.Join(st.Participants, ", "))
	if st.FacilitatorRole.Declared {
		fmt.Fprintf(&b, "facilitator: %s (participates: %v)\n", st.FacilitatorRole.Facilitator, st.FacilitatorRole.Participates)
	}
	fmt.Fprintf(&b, "run phase pointer (advisory): %s\n\n", runPhase)

	fmt.Fprintf(&b, "## PhaseDigest (mechanically derived; raw artifacts are canonical)\n")
	if digest.Round != nil {
		fmt.Fprintf(&b, "%s: %d/%d filed-and-valid\n", digest.Round.Label, digest.Round.Completed, digest.Round.Total)
		for _, row := range digest.Round.Rows {
			writeBriefRow(&b, row)
		}
	}
	if digest.Review != nil {
		fmt.Fprintf(&b, "%s: %d/%d filed-and-valid\n", digest.Review.Label, digest.Review.Completed, digest.Review.Total)
		for _, row := range digest.Review.Rows {
			writeBriefRow(&b, row)
		}
	}
	if digest.Consensus != nil {
		fmt.Fprintf(&b, "consensus: triage=%s missing=%v signed=%v reservations=%v blocks=%v\n",
			digest.Consensus.Triage, digest.Consensus.Missing, digest.Consensus.Signed, digest.Consensus.Reservations, digest.Consensus.Blocks)
	}
	if digest.Implementation != nil {
		fmt.Fprintf(&b, "implementation: present=%v status=%s implementer=%s\n",
			digest.Implementation.Present, digest.Implementation.Status, digest.Implementation.Implementer)
	}

	fmt.Fprintf(&b, "\n## Next action (fixed enumeration)\n%s\n", digest.Next)
	if digest.Next == driver.NextAdjudicateRaw {
		b.WriteString("Open the raw artifact named above to adjudicate before acting.\n")
	}
	return b.String()
}

func writeBriefRow(b *strings.Builder, row driver.PhaseAgentRow) {
	validity := "ok"
	if !row.Valid {
		validity = row.Validity
	}
	fmt.Fprintf(b, "  %s: filed=%v bytes=%d owner=%s valid=%v(%s) unparsed=%v path=%s\n",
		row.Agent, row.Filed, row.Bytes, row.Owner, row.Valid, validity, row.Unparsed, row.Path)
}
