package runner

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"path/filepath"
	"strconv"
	"strings"

	"parley-deck-cli/internal/agents"
	"parley-deck-cli/internal/config"
	"parley-deck-cli/internal/protocol"
	"parley-deck-cli/internal/protocolcore"
	"parley-deck-cli/internal/protocolpacket"
	"parley-deck-cli/internal/telemetry"
)

type protocolContextError struct{ reason string }

func (e *protocolContextError) Error() string { return "protocol context refused: " + e.reason }

// prepareProtocolPrompt is launch glue around the one shared renderer. It
// passes the rendered in-memory bytes, rather than re-reading a mutable path
// after attestation. Full context remains the only enabled launch treatment.
// Unknown phase/track metadata remains unknown in the shadow audit.
func prepareProtocolPrompt(root, prompt string, info LaunchInfo) (string, telemetry.Context, error) {
	refuse := func(reason string) (string, telemetry.Context, error) {
		return "", telemetry.Context{Mode: protocolpacket.ModeRefused, FallbackReason: telemetry.String(reason)},
			&protocolContextError{reason: reason}
	}
	home := config.CentralHome()
	if strings.TrimSpace(home) == "" {
		return refuse("authority-unavailable")
	}
	phase := -1
	switch info.Phase {
	case "preflight":
		phase = 0
	case "round-01":
		phase = 1
	case "consensus":
		phase = 3
	case "final":
		phase = 4
	case "implementation":
		phase = 5
	case "review":
		phase = 6
	case "review-consensus":
		phase = 7
	case "fixup":
		phase = 8
	default:
		if strings.HasPrefix(info.Phase, "round-") {
			if round, err := strconv.Atoi(strings.TrimPrefix(info.Phase, "round-")); err == nil && round > 0 && info.Phase == roundLabel(round) {
				phase = 2
				if round == 1 {
					phase = 1
				}
			}
		}
	}
	track := "unknown"
	if info.Idea != "" && info.Idea != "." && info.Idea != ".." && filepath.Base(info.Idea) == info.Idea && !strings.ContainsAny(info.Idea, `/\\`) {
		if fields, err := protocol.ReadFrontmatter(filepath.Join(root, protocol.DeckDir, "ideas", info.Idea, "00-prompt.md")); err == nil {
			switch fields["track"] {
			case "fast", "standard", "deliberation":
				track = fields["track"]
			}
		}
	}
	c, err := protocolpacket.Render(root, protocolcore.StoreAt(home), protocolpacket.Request{
		Phase: phase, Track: track, IdeaSlug: info.Idea,
	})
	if err != nil {
		if errors.Is(err, protocolpacket.ErrAuthority) {
			return refuse("authority-unavailable")
		}
		return refuse("context-publication-failed")
	}
	if c.ContextMode == protocolpacket.ModeRefused {
		if reason := telemetry.SafeLabel(c.FallbackReason); reason != nil {
			return refuse(*reason)
		}
		return refuse("renderer-refused")
	}
	if c.Body == "" || protocolpacket.Hash(c.Body) != c.PacketSHA256 {
		return refuse("context-hash-mismatch")
	}
	if strings.Contains(c.Body, "</parley-protocol>") || strings.Contains(c.Body, "<parley-protocol>") {
		return refuse("protocol-envelope-collision")
	}
	attestation, err := json.Marshal(c.Attestation)
	if err != nil {
		return refuse("attestation-encoding-failed")
	}
	// A full source and a failed shadow optimization are separate observations.
	// Preserve the full shadow audit in the private launch prompt, without
	// passing raw source diagnostics into public telemetry labels.
	shadow, err := json.Marshal(c.Shadow)
	if err != nil {
		return refuse("shadow-encoding-failed")
	}
	ctx := telemetry.Context{Mode: c.ContextMode, SourceSHA256: telemetry.String(c.SourceSHA256),
		PacketSHA256: telemetry.String(c.PacketSHA256), FallbackReason: telemetry.SafeLabel(c.FallbackReason)}
	if c.FallbackReason != "" && ctx.FallbackReason == nil {
		ctx.FallbackReason = telemetry.String("renderer-fallback")
	}
	return fmt.Sprintf("Protocol context attestation: %s\nShadow packet audit: %s\nThis is an unapplied diagnostic only. Shadow included/omitted block counts do not describe the supplied full protocol.\n\nThe following is the resolved live protocol, supplied verbatim for this launch.\n<parley-protocol>\n%s\n</parley-protocol>\n\nLaunch task:\n%s", attestation, shadow, c.Body, prompt), ctx, nil
}

// beginProtocolLaunch owns attestation at each actual launch boundary. Callers
// cannot inject a precomputed Context to certify different prompt bytes.
func beginProtocolLaunch(ctx context.Context, root, runID string, agent agents.Discovery, prompt string) (context.Context, string, *launchEvidence, error) {
	info, _ := ctx.Value(launchInfoKey{}).(LaunchInfo)
	prepared, protocolContext, contextErr := prepareProtocolPrompt(root, prompt, info)
	info.Context = protocolContext
	ctx = WithLaunchInfo(ctx, info)
	if contextErr != nil {
		// An already refused prompt cannot execute. Retain its request/terminal
		// before any fresh cycle, step, launch or helper-ticket reservation.
		origin, requestInfo := launchRequestInfo(ctx, root, runID)
		evidence, err := recordLaunchRequest(origin, agent, requestInfo)
		if err != nil {
			return ctx, "", nil, err
		}
		if err := evidence.finish(contextErr, ctx.Err(), nil); err != nil {
			return ctx, "", nil, err
		}
		return ctx, "", nil, contextErr
	}
	evidence, err := beginLaunch(ctx, root, runID, agent)
	if err != nil {
		return ctx, "", nil, err
	}
	return ctx, prepared, evidence, nil
}

func protocolLaunchPhase(opts Options) string {
	if opts.Phase == "" || opts.Phase == "deliberation" {
		round := opts.Round
		if round < 1 {
			round = 1
		}
		return roundLabel(round)
	}
	return opts.Phase
}
