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

// briefPhase derives the protocol phase pointer for an idea status.
func briefPhase(status string) int {
	switch {
	case strings.HasPrefix(status, "round-01"):
		return 1
	case strings.HasPrefix(status, "round-"):
		return 2
	case status == "consensus":
		return 3
	case status == "final":
		return 4
	case strings.HasPrefix(status, "review"):
		return 6
	default:
		return 0
	}
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

	phase := briefPhase(st.Status)
	track := "standard"
	if t, present := driver.ReadTrack(st.Path); present && t != "" {
		track = string(t)
	}

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

	digest := driver.BuildPhaseDigest(root, st.Slug, st.Path, st.Participants)
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
	b.WriteString("Computed view — never stored. Recompute with `parley organizer brief`;\n")
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
