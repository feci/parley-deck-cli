package app

import (
	"errors"
	"flag"
	"fmt"
	"io"
	"path/filepath"
	"strings"

	"parley-deck-cli/internal/protocolpacket"
)

// `parley protocol packet` — the phase-scoped protocol context renderer and its check.
const protocolPacketUsage = `usage:
  parley protocol packet [--dir DIR] --phase N [--track T] [--transport X] [--idea SLUG]
                         [--flag F]... [--optimize] [--json] [--print]
  parley protocol packet check [--dir DIR] [--json]

  Renders the protocol context for one launch from the LIVE resolved authority (a source-role
  deck's own file; a consumer deck's verified core + lock + overlay) — never a bundled
  snapshot — and prints the launch attestation: context_mode, source_sha256, packet_sha256,
  fallback_reason. The default is full context plus a shadow packet audit record. --optimize is
  the explicit experimental input of the ratified packet trial, not an enabled release; unknown
  applicability, dependencies, phase/track/transport/flag values or drift fall back to full
  context with the reason. Missing authority or a detected secret refuses. Bodies are written
  under .parley-runtime/protocol-packets/ and are never committed.

  packet check verifies parley-deck/meta/packet-applicability.yaml against the live source:
  every heading classified once, every omission has a trigger, dependencies resolve, the
  ratified never-cut list holds, and an optimized packet renders for every phase and track.`

type packetFlagList []string

func (l *packetFlagList) String() string     { return strings.Join(*l, ",") }
func (l *packetFlagList) Set(v string) error { *l = append(*l, strings.TrimSpace(v)); return nil }

func runProtocolPacket(args []string, stdout, stderr io.Writer) int {
	if len(args) > 0 {
		switch args[0] {
		case "check":
			return protocolPacketCheck(args[1:], stdout, stderr)
		case "--help", "-h", "help":
			fmt.Fprintln(stdout, protocolPacketUsage)
			return 0
		}
	}
	fs := flag.NewFlagSet("protocol packet", flag.ContinueOnError)
	fs.SetOutput(stderr)
	dir := fs.String("dir", ".", "workspace directory")
	phase := fs.Int("phase", -1, "protocol phase 0..8 (required)")
	track := fs.String("track", "standard", "track: fast|standard|deliberation")
	transport := fs.String("transport", "", "transport (default: the deck header's Transport: value)")
	idea := fs.String("idea", "", "idea slug (a meta-protocol-change-* slug sets the protocol_change flag)")
	optimize := fs.Bool("optimize", false, "EXPERIMENTAL: emit the optimized packet when every guard passes (default: full context + shadow audit)")
	jsonOut := fs.Bool("json", false, "machine-readable attestation and index")
	print := fs.Bool("print", false, "also write the emitted body to stdout")
	var flags packetFlagList
	fs.Var(&flags, "flag", "launch flag: strict_gate|auto_implement|pipeline|protocol_change (repeatable)")
	if err := fs.Parse(args); err != nil {
		return 2
	}
	if *phase < 0 {
		fmt.Fprintln(stderr, "protocol packet: --phase is required\n"+protocolPacketUsage)
		return 2
	}
	root, err := filepath.Abs(*dir)
	if err != nil {
		fmt.Fprintf(stderr, "protocol packet: %v\n", err)
		return 1
	}
	c, err := BuildProtocolContext(root, protocolpacket.Request{
		Phase: *phase, Track: *track, Transport: *transport, IdeaSlug: *idea, Flags: flags, Optimize: *optimize,
	})
	if err != nil {
		fmt.Fprintf(stderr, "protocol packet: %v\n", err)
		return 1
	}
	if c.ContextMode == protocolpacket.ModeRefused {
		fmt.Fprintf(stderr, "protocol packet: REFUSED — %s; no context emitted\n", c.FallbackReason)
		if *jsonOut {
			_ = writeJSON(stdout, stderr, c)
		}
		return 1
	}
	if *jsonOut {
		if code := writeJSON(stdout, stderr, c); code != 0 {
			return code
		}
	} else {
		fmt.Fprintf(stdout, "context_mode    : %s\n", c.ContextMode)
		fmt.Fprintf(stdout, "source          : %s (%s, %s)\n", c.Source.Path, c.Source.Role, c.Source.Authority)
		fmt.Fprintf(stdout, "source_sha256   : %s\n", c.SourceSHA256)
		fmt.Fprintf(stdout, "packet_sha256   : %s\n", c.PacketSHA256)
		if c.FallbackReason != "" {
			fmt.Fprintf(stdout, "fallback_reason : %s\n", c.FallbackReason)
		}
		if c.Shadow != nil {
			fmt.Fprintf(stdout, "shadow packet   : %d/%d blocks included, %d of %d bytes, sha256 %s%s\n",
				c.Shadow.IncludedBlocks, c.Shadow.IncludedBlocks+c.Shadow.OmittedBlocks, c.Shadow.PacketBytes, c.Shadow.SourceBytes,
				orDash(c.Shadow.PacketSHA256), suffixIf(c.Shadow.WouldFallbackReason != "", " (would fall back: "+c.Shadow.WouldFallbackReason+")"))
		}
		fmt.Fprintf(stdout, "body            : %s\n", c.BodyPath)
	}
	if *print {
		fmt.Fprint(stdout, c.Body)
	}
	return 0
}

func suffixIf(cond bool, s string) string {
	if cond {
		return s
	}
	return ""
}

func protocolPacketCheck(args []string, stdout, stderr io.Writer) int {
	fs := flag.NewFlagSet("protocol packet check", flag.ContinueOnError)
	fs.SetOutput(stderr)
	dir := fs.String("dir", ".", "workspace directory")
	jsonOut := fs.Bool("json", false, "machine-readable report")
	if err := fs.Parse(args); err != nil {
		return 2
	}
	root, err := filepath.Abs(*dir)
	if err != nil {
		fmt.Fprintf(stderr, "protocol packet check: %v\n", err)
		return 1
	}
	store, err := coreStore()
	if err != nil {
		fmt.Fprintf(stderr, "protocol packet check: %v\n", err)
		return 1
	}
	src, err := protocolpacket.ResolveSource(root, store)
	if err != nil {
		fmt.Fprintf(stderr, "protocol packet check: %v\n", err)
		return 1
	}
	m, err := protocolpacket.LoadMap(protocolpacket.MapPath(root))
	if err != nil && !errors.Is(err, protocolpacket.ErrNoMap) {
		fmt.Fprintf(stderr, "protocol packet check: %v\n", err)
		return 1
	}
	rep := protocolpacket.Check(src, m)
	if *jsonOut {
		_ = writeJSON(stdout, stderr, rep)
	} else {
		fmt.Fprintf(stdout, "source   : %s (%s)\n", src.Path, src.Authority)
		fmt.Fprintf(stdout, "sha256   : %s\n", rep.SourceSHA256)
		fmt.Fprintf(stdout, "blocks   : %d\n", rep.Blocks)
		if m == nil {
			fmt.Fprintf(stdout, "map      : ABSENT (%s) — every launch falls back to full context\n", protocolpacket.MapPath(root))
		}
		list := func(label string, items []string) {
			if len(items) == 0 {
				return
			}
			fmt.Fprintf(stdout, "%s:\n", label)
			for _, it := range items {
				fmt.Fprintf(stdout, "  - %s\n", it)
			}
		}
		list("unclassified headings (unknown applicability; included, and optimization falls back)", rep.Unclassified)
		list("stale map entries (no such heading in the live source)", rep.Stale)
		list("never-cut violations", rep.NeverCut)
		list("fallbacks when rendering every phase/track", rep.Fallbacks)
		if rep.Secret != "" {
			fmt.Fprintf(stdout, "secret   : detected (%s) — the renderer refuses every launch until it is removed\n", rep.Secret)
		}
		if rep.OK {
			fmt.Fprintln(stdout, "packet check: ok (structure only — a wrong classification present in the map is not detectable here)")
		} else {
			fmt.Fprintln(stdout, "packet check: FAILED")
		}
	}
	if rep.OK {
		return 0
	}
	return 1
}

// BuildProtocolContext is the app-level entry every builder and handoff should call: resolve
// the deck's authority, render, and write the body under .parley-runtime/protocol-packets/.
// It returns an error only when nothing can be attested (authority or I/O); a fallback or a
// refusal is carried in the context so the launch record can attest it. A refused context has
// no body and MUST NOT be launched with substituted protocol text.
func BuildProtocolContext(root string, req protocolpacket.Request) (protocolpacket.Context, error) {
	store, err := coreStore()
	if err != nil {
		return protocolpacket.Context{}, err
	}
	return protocolpacket.Render(root, store, req)
}
