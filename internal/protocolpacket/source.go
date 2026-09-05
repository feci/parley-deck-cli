package protocolpacket

// Authority resolution.
//
// Source-role decks (this repository) are the protocol's upstream: the live deck file IS the
// authority, and the expected drift against an older global core must never substitute stale
// core text (ALT-07). Consumer decks render from the pinned, hash-attested core plus lock plus
// overlay; that chain must verify BEFORE any optimization, and a verified render that differs
// from the on-disk view is drift that cannot pass as a valid optimized packet. Missing or
// unprovable authority is an error, never a bundled snapshot.

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"parley-deck-cli/internal/protocol"
	"parley-deck-cli/internal/protocolcore"
)

// ErrAuthority wraps every "cannot prove what the protocol is" failure.
var ErrAuthority = errors.New("protocolpacket: authority")

// Source is the resolved authoritative protocol text.
type Source struct {
	Path      string
	Role      string
	Authority string
	Raw       string
	Transport string // parsed from the `**Transport:**` header line
	Drift     string // consumer only: non-empty when the on-disk view is not the verified render
}

// DeckProtocolPath is the deck's protocol view.
func DeckProtocolPath(root string) string {
	return filepath.Join(root, protocol.DeckDir, "COOPERATION.md")
}

// MapPath is the applicability map location.
func MapPath(root string) string {
	return filepath.Join(root, protocol.DeckDir, "meta", "packet-applicability.yaml")
}

// RuntimeDir is where generated bodies go. The path is ignored by git (the ignore entry is
// integration-owned); a packet is never committed because a committed packet is a stale copy.
func RuntimeDir(root string) string {
	return filepath.Join(root, ".parley-runtime", "protocol-packets")
}

// ResolveSource resolves the deck's authoritative protocol text according to its role.
func ResolveSource(root string, store protocolcore.Store) (Source, error) {
	role, err := protocolRole(root)
	if err != nil {
		return Source{}, err
	}
	path := DeckProtocolPath(root)
	raw, err := os.ReadFile(path)
	if err != nil {
		return Source{}, fmt.Errorf("%w: cannot read %s: %v", ErrAuthority, path, err)
	}
	src := Source{Path: path, Role: role, Raw: string(raw), Transport: headerTransport(string(raw))}
	switch role {
	case "source":
		src.Authority = "live source file (protocolRole: source); global core drift is expected and never substituted"
		return src, nil
	case "consumer":
		lockPath := filepath.Join(root, protocol.DeckDir, "meta", "protocol-lock.yaml")
		lb, err := os.ReadFile(lockPath)
		if err != nil {
			return Source{}, fmt.Errorf("%w: consumer deck has no readable %s; the deck view cannot be attested to a core release", ErrAuthority, lockPath)
		}
		lock, err := protocolcore.ParseLock(string(lb))
		if err != nil {
			return Source{}, fmt.Errorf("%w: %v", ErrAuthority, err)
		}
		rel, err := store.Load(lock.CoreVersion)
		if err != nil {
			return Source{}, fmt.Errorf("%w: %w", ErrAuthority, err)
		}
		if rel.SHA256 != lock.CoreBodySHA256 {
			return Source{}, fmt.Errorf("%w: lock attests core %s body %s but the installed release hashes %s",
				ErrAuthority, lock.CoreVersion, protocolcore.ShortHash(lock.CoreBodySHA256), protocolcore.ShortHash(rel.SHA256))
		}
		ovPath := filepath.Join(root, protocol.DeckDir, protocolcore.OverlayFileName)
		var ov *protocolcore.Overlay
		ob, oerr := os.ReadFile(ovPath)
		switch {
		case oerr == nil:
			ov, err = protocolcore.ReconcileOverlay(lock, true, string(ob), nil)
		case errors.Is(oerr, os.ErrNotExist):
			ov, err = protocolcore.ReconcileOverlay(lock, false, "", nil)
		default:
			ov, err = protocolcore.ReconcileOverlay(lock, true, "", oerr)
		}
		if err != nil {
			return Source{}, fmt.Errorf("%w: %v", ErrAuthority, err)
		}
		res, err := protocolcore.Render(rel, src.Raw, ov)
		if err != nil {
			return Source{}, fmt.Errorf("%w: %v", ErrAuthority, err)
		}
		src.Authority = fmt.Sprintf("core %s (%s) + lock + overlay, verified", lock.CoreVersion, protocolcore.ShortHash(rel.SHA256))
		if res.Body != src.Raw {
			src.Drift = "deck view differs from the verified core render (protocol check: hand-edited-or-stale)"
		}
		return src, nil
	default:
		return Source{}, fmt.Errorf("%w: meta/version.json protocolRole %q is unknown; confirm source|consumer (parley preflight)", ErrAuthority, role)
	}
}

func protocolRole(root string) (string, error) {
	path := filepath.Join(root, protocol.DeckDir, "meta", "version.json")
	b, err := os.ReadFile(path)
	if err != nil {
		return "", fmt.Errorf("%w: cannot read %s: %v", ErrAuthority, path, err)
	}
	var meta struct {
		ProtocolRole string `json:"protocolRole"`
	}
	if err := json.Unmarshal(b, &meta); err != nil {
		return "", fmt.Errorf("%w: %s is not valid JSON: %v", ErrAuthority, path, err)
	}
	return strings.ToLower(strings.TrimSpace(meta.ProtocolRole)), nil
}

// headerTransport reads the deck's `**Transport:** `value“ line. An absent or placeholder
// value is returned as-is and rejected by Build as unknown-transport.
func headerTransport(body string) string {
	for _, l := range strings.Split(body, "\n") {
		if !strings.HasPrefix(l, "**Transport:**") {
			continue
		}
		v := strings.TrimSpace(strings.TrimPrefix(l, "**Transport:**"))
		return strings.Trim(v, "` ")
	}
	return ""
}

// WriteBody writes the emitted body under RuntimeDir and records the path on the context.
// A refused context has no body and nothing is written.
func WriteBody(root string, ctx *Context) error {
	if ctx.ContextMode == ModeRefused {
		return fmt.Errorf("protocolpacket: refused context has no body (%s)", ctx.FallbackReason)
	}
	dir := RuntimeDir(root)
	if err := os.MkdirAll(dir, 0o700); err != nil {
		return err
	}
	name := fmt.Sprintf("%s-phase%d-%s-%s.md", ctx.ContextMode, ctx.Request.Phase, safeName(ctx.Request.Track), protocolcore.ShortHash(ctx.PacketSHA256))
	path := filepath.Join(dir, name)
	if err := os.WriteFile(path, []byte(ctx.Body), 0o600); err != nil {
		return err
	}
	ctx.BodyPath = path
	return nil
}

func safeName(s string) string {
	var b strings.Builder
	for _, r := range s {
		switch {
		case r >= 'a' && r <= 'z', r >= '0' && r <= '9', r == '-':
			b.WriteRune(r)
		default:
			b.WriteRune('_')
		}
	}
	if b.Len() == 0 {
		return "unknown"
	}
	return b.String()
}

// Render is the one-call entry: resolve authority, load the map, build, and write the body.
// The returned error is an authority error (nothing emitted) or an I/O error; a fallback or a
// refusal is reported in the context, not as an error, so callers can attest it.
func Render(root string, store protocolcore.Store, req Request) (Context, error) {
	src, err := ResolveSource(root, store)
	if err != nil {
		return Context{}, err
	}
	m, err := LoadMap(MapPath(root))
	if err != nil && !errors.Is(err, ErrNoMap) {
		// A malformed map is not "no map": it is reported as the fallback reason so the
		// operator sees it, while the body stays full context.
		m = nil
		req.Flags = append([]string(nil), req.Flags...)
		c := Build(src, nil, req)
		if c.ContextMode != ModeRefused {
			c.Problems = append([]string{"malformed-map:" + err.Error()}, c.Problems...)
			if c.ContextMode == ModeFullFallback {
				c.FallbackReason = "malformed-map:" + err.Error() + "; " + c.FallbackReason
			} else if c.Shadow != nil {
				c.Shadow.WouldFallbackReason = "malformed-map:" + err.Error() + "; " + c.Shadow.WouldFallbackReason
			}
			if werr := WriteBody(root, &c); werr != nil {
				return c, werr
			}
		}
		return c, nil
	}
	c := Build(src, m, req)
	if c.ContextMode == ModeRefused {
		return c, nil
	}
	if err := WriteBody(root, &c); err != nil {
		return c, err
	}
	return c, nil
}
