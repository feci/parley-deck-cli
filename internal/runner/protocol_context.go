package runner

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"parley-deck-cli/internal/config"
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
	case "implementation":
		phase = 5
	case "review":
		phase = 6
	case "review-consensus":
		phase = 7
	case "fixup":
		phase = 8
	}
	c, err := protocolpacket.Render(root, protocolcore.StoreAt(home), protocolpacket.Request{
		Phase: phase, Track: "unknown", IdeaSlug: info.Idea,
	})
	if err != nil {
		if errors.Is(err, protocolpacket.ErrAuthority) {
			return refuse("authority-unavailable")
		}
		return refuse("context-publication-failed")
	}
	if c.ContextMode == protocolpacket.ModeRefused {
		return refuse("renderer-refused")
	}
	if c.Body == "" || protocolpacket.Hash(c.Body) != c.PacketSHA256 {
		return refuse("context-hash-mismatch")
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
	return fmt.Sprintf("Protocol context attestation: %s\nShadow packet audit: %s\n\nThe following is the resolved live protocol, supplied verbatim for this launch.\n<parley-protocol>\n%s\n</parley-protocol>\n\nLaunch task:\n%s", attestation, shadow, c.Body, prompt), ctx, nil
}
