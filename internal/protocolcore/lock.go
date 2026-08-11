package protocolcore

// The deck lock, schema v2.
//
// This is the load-bearing compatibility control of the overlay slice, and it is here rather than
// in a follow-up because of what codex-1 caught during consensus: a release-format marker is read
// only by NEW code, but the binary that has to be stopped is an OLD one, and what an old binary
// reads is the LOCK.
//
// The pre-v2 lock was a flat `core-version:` line, and the old parser scans for exactly that prefix
// and ignores every other line. So a lock carrying overlay state would be read by a stale binary as
// an ordinary lock and the deck rendered with the overlay silently absent. Nesting the version
// under `core:` removes that prefix from the file, so a stale binary finds no pinned version and
// refuses — verified against the shipped parley 1.42.1, which exits non-zero on both `check` and
// `render` rather than rendering a partial view.

import (
	"errors"
	"fmt"
	"strings"

	"gopkg.in/yaml.v3"
)

// LockSchemaV2 is the only accepted lock schema.
const LockSchemaV2 = "parley.protocol-lock/v2"

// ResolverV1 is the exact resolver literal v1 writes and accepts. It is a literal rather than an
// integer so that a future resolver change is a visibly different string in a diff.
const ResolverV1 = "overlay-v1"

// OverlayNone is the literal recorded when the deck has no overlay file. It is a required positive
// statement: a missing key would be indistinguishable from an older writer that did not know about
// overlays.
const OverlayNone = "none"

// Lock is a parsed deck lock.
type Lock struct {
	Schema          string
	CoreVersion     string
	CoreBodySHA256  string
	Overlay         string // OverlayNone or 64 lowercase hex
	ResolverVersion string
}

type lockDoc struct {
	Schema          string      `yaml:"schema"`
	Core            lockCoreDoc `yaml:"core"`
	Overlay         string      `yaml:"overlay"`
	ResolverVersion string      `yaml:"resolver-version"`
}

type lockCoreDoc struct {
	Version    string `yaml:"version"`
	BodySHA256 string `yaml:"body-sha256"`
}

// ParseLock validates and decodes a deck lock. It is pure: the caller reads the file.
func ParseLock(raw string) (Lock, error) {
	norm := normalizeLF(raw)
	if strings.TrimSpace(norm) == "" {
		return Lock{}, errors.New("lock: file is empty")
	}
	// A pre-v2 flat lock is named explicitly rather than failing with a generic schema error,
	// because the fix is a migration and the operator should be told so.
	for _, l := range strings.Split(norm, "\n") {
		if strings.HasPrefix(strings.TrimSpace(l), "core-version:") {
			return Lock{}, fmt.Errorf(
				"lock: this is a pre-v2 flat lock (`core-version:`). Migrate it to %s:\n"+
					"  schema: %s\n  core:\n    version: <v>\n    body-sha256: <64 hex>\n  overlay: none | <64 hex>\n  resolver-version: %s",
				LockSchemaV2, LockSchemaV2, ResolverV1)
		}
	}

	var ld lockDoc
	dec := yaml.NewDecoder(strings.NewReader(norm))
	dec.KnownFields(true)
	if err := dec.Decode(&ld); err != nil {
		return Lock{}, fmt.Errorf("lock: %w", err)
	}
	var extra yaml.Node
	if err := dec.Decode(&extra); err == nil {
		return Lock{}, errors.New("lock: file contains more than one YAML document")
	}

	if ld.Schema != LockSchemaV2 {
		return Lock{}, fmt.Errorf("lock: schema is %q, want %q", ld.Schema, LockSchemaV2)
	}
	if !ValidVersion(ld.Core.Version) {
		// The same guard the release store applies. The lock is a COMMITTED file any contributor can
		// edit, and its version is joined onto the store root.
		return Lock{}, fmt.Errorf("lock: core.version %q is missing or unsafe", ld.Core.Version)
	}
	if !hex64Re.MatchString(ld.Core.BodySHA256) {
		return Lock{}, fmt.Errorf("lock: core.body-sha256 %q; want 64 lowercase hex characters", ld.Core.BodySHA256)
	}
	if ld.Overlay != OverlayNone && !hex64Re.MatchString(ld.Overlay) {
		return Lock{}, fmt.Errorf("lock: overlay %q; want %q or 64 lowercase hex characters", ld.Overlay, OverlayNone)
	}
	if ld.ResolverVersion != ResolverV1 {
		return Lock{}, fmt.Errorf("lock: resolver-version %q, want %q", ld.ResolverVersion, ResolverV1)
	}

	return Lock{
		Schema:          ld.Schema,
		CoreVersion:     ld.Core.Version,
		CoreBodySHA256:  ld.Core.BodySHA256,
		Overlay:         ld.Overlay,
		ResolverVersion: ld.ResolverVersion,
	}, nil
}

// DeclaresOverlay reports whether the lock says an overlay exists, and its expected hash.
func (l Lock) DeclaresOverlay() (string, bool) {
	if l.Overlay == OverlayNone || l.Overlay == "" {
		return "", false
	}
	return l.Overlay, true
}

// Render serializes a lock. Used by tests and by any writer; the field order is fixed so that a
// rewritten lock diffs cleanly.
func (l Lock) Render() string {
	var b strings.Builder
	fmt.Fprintf(&b, "schema: %s\n", LockSchemaV2)
	b.WriteString("core:\n")
	fmt.Fprintf(&b, "  version: %s\n", l.CoreVersion)
	fmt.Fprintf(&b, "  body-sha256: %s\n", l.CoreBodySHA256)
	overlay := l.Overlay
	if overlay == "" {
		overlay = OverlayNone
	}
	fmt.Fprintf(&b, "overlay: %s\n", overlay)
	fmt.Fprintf(&b, "resolver-version: %s\n", ResolverV1)
	return b.String()
}

// ReconcileOverlay decides whether a deck's overlay state matches what its lock attests, and
// returns the overlay to compose with.
//
// Every mismatch BLOCKS before composition. That is the whole point: the failure this design exists
// to prevent is a render that quietly proceeds while the deck's local content is absent from the
// output. Rendering "what we could resolve" is exactly the 2026-08-06 silent-erasure incident with
// a version number attached.
//
// present reports whether the overlay FILE exists; raw is its contents when it does; readErr is any
// error encountered reading it (an unreadable file is never treated as an absent one).
func ReconcileOverlay(l Lock, present bool, raw string, readErr error) (*Overlay, error) {
	want, declared := l.DeclaresOverlay()

	if readErr != nil {
		return nil, fmt.Errorf("overlay: %s is unreadable: %w", OverlayFileName, readErr)
	}
	if !present {
		if declared {
			return nil, fmt.Errorf(
				"overlay: the lock declares an overlay (%s) but %s is absent; refusing to render without it",
				ShortHash(want), OverlayFileName)
		}
		return nil, nil // absent and undeclared: the canonical "no customization" state
	}
	if !declared {
		return nil, fmt.Errorf(
			"overlay: %s exists but the lock records `overlay: %s`; adopt it explicitly rather than rendering without it",
			OverlayFileName, OverlayNone)
	}

	ov, err := ParseOverlay(raw)
	if err != nil {
		return nil, err
	}
	if ov.SHA256 != want {
		return nil, fmt.Errorf(
			"overlay: %s has hash %s but the lock attests %s; the overlay changed without the lock being updated",
			OverlayFileName, ShortHash(ov.SHA256), ShortHash(want))
	}
	return &ov, nil
}
