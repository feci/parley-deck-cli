package protocolcore

// The deck overlay: the one channel through which a deck adds project-local protocol content.
//
// From FINAL.md of idea protocol-overlay-local-extension. v1 is EXTEND-ONLY by user ruling —
// there is no replace operation, so nothing here addresses a block of the core. That is why v1
// ships no per-block registry: the only position the overlay needs is the terminal core/overlay
// boundary, which is defined as the end of the normalized core body rather than looked up.
//
// The grammar is deliberately smaller than YAML in general. Unknown keys, aliases, duplicate keys,
// a second document and any trailing body are refused rather than ignored, because a "strict"
// schema that silently drops what it does not understand is how a deck ends up governed by a
// protocol nobody wrote.

import (
	"errors"
	"fmt"
	"regexp"
	"strings"

	"gopkg.in/yaml.v3"
)

// OverlayFileName is the deck-relative name of the overlay. It is fixed, not configurable: a
// configurable path is one more thing a stale binary can disagree about.
const OverlayFileName = "protocol-overlay.md"

// OverlaySchema is the only accepted schema value in v1.
const OverlaySchema = "parley.protocol-overlay/v1"

// KindExtend is the only operation kind in v1. `replace` was dropped by user ruling and is the
// deferred `protocol-overlay-replace` follow-up.
const KindExtend = "extend"

// ErrNoOverlay reports that the deck has no overlay file.
//
// Callers must distinguish this from a parse error: absence is the ONLY canonical "no
// customization" state and is always legal, while an unreadable or invalid overlay blocks.
var ErrNoOverlay = errors.New("protocolcore: deck has no overlay")

// operationIDRe is the deck-namespaced operation id grammar. The file itself is the deck namespace,
// so the id is deck-local; deriving a namespace from Workspace or a remote URL would make identity
// depend on values the fleet measurement showed are unstable (16 of 29 decks disagree with their
// own directory name).
var operationIDRe = regexp.MustCompile(`^deck\.[a-z0-9](?:[a-z0-9-]{0,61}[a-z0-9])?$`)

var hex64Re = regexp.MustCompile(`^[0-9a-f]{64}$`)

// Operation is one overlay operation. In v1 there is at most one, and its Kind is always
// KindExtend.
type Operation struct {
	ID        string
	Kind      string
	Rationale string
	// Markdown is the decoded payload, LF-normalized. It is rendered verbatim at the terminal
	// boundary.
	Markdown string
	// CoreSHA256 is the core body hash this operation was written against. A mismatch is a
	// reconfirmation prompt, never an automatic migration.
	CoreSHA256 string
	// PayloadSHA256 hashes the DECODED payload and is used only in change reports. It is
	// deliberately NOT the lock hash — see Overlay.SHA256.
	PayloadSHA256 string
}

// Overlay is a parsed, validated overlay file.
type Overlay struct {
	// Raw is the entire file, LF-normalized. It is what SHA256 covers.
	Raw string
	// SHA256 is the LOCK hash: the SHA-256 of the UTF-8 bytes of the ENTIRE overlay file after
	// CRLF/CR to LF normalization.
	//
	// It must cover the whole file, not the decoded payload. codex-1 blocked consensus revision 2
	// over precisely this: hashing only the payload would let a change to an operation id,
	// rationale or compatibility metadata leave the lock unchanged, so the lock would attest to an
	// overlay that is no longer the one on disk.
	SHA256     string
	Operations []Operation
}

// yaml shapes. Fields are exhaustive on purpose: decoding with KnownFields(true) turns any key not
// listed here into an error rather than a silent drop.
type overlayDoc struct {
	Schema     string         `yaml:"schema"`
	Operations []operationDoc `yaml:"operations"`
}

type operationDoc struct {
	ID         string `yaml:"id"`
	Kind       string `yaml:"kind"`
	Rationale  string `yaml:"rationale"`
	Markdown   string `yaml:"markdown"`
	CoreSHA256 string `yaml:"core-sha256"`
}

// ParseOverlay validates and decodes an overlay file.
//
// raw is the file exactly as read from disk. Normalization happens here so that the hash the lock
// compares against is computed from the same bytes on every platform.
func ParseOverlay(raw string) (Overlay, error) {
	norm := normalizeLF(raw)
	if strings.TrimSpace(norm) == "" {
		// An EMPTY overlay is invalid, while an ABSENT one is the canonical "no customization"
		// state. Conflating them would make "this deck has no local content" indistinguishable from
		// "this deck's local content failed to save".
		return Overlay{}, errors.New("overlay: file is empty; absence of the file is the only 'no customization' state")
	}

	doc, err := yamlDocument(norm)
	if err != nil {
		return Overlay{}, err
	}

	var od overlayDoc
	dec := yaml.NewDecoder(strings.NewReader(doc))
	dec.KnownFields(true)
	if err := dec.Decode(&od); err != nil {
		return Overlay{}, fmt.Errorf("overlay: %w", err)
	}
	// A second document is refused rather than ignored: "the first document wins" is a silent
	// content drop with extra steps.
	var extra yaml.Node
	if err := dec.Decode(&extra); err == nil {
		return Overlay{}, errors.New("overlay: file contains more than one YAML document; exactly one is allowed")
	}

	if err := rejectAliases(doc); err != nil {
		return Overlay{}, err
	}

	if od.Schema != OverlaySchema {
		return Overlay{}, fmt.Errorf("overlay: schema is %q, want %q", od.Schema, OverlaySchema)
	}
	if len(od.Operations) == 0 {
		return Overlay{}, errors.New("overlay: declares no operations; delete the file instead — absence is the 'no customization' state")
	}
	if len(od.Operations) > 1 {
		return Overlay{}, fmt.Errorf("overlay: v1 allows at most one operation, found %d", len(od.Operations))
	}

	ops := make([]Operation, 0, len(od.Operations))
	seen := map[string]bool{}
	for _, o := range od.Operations {
		if o.Kind != KindExtend {
			return Overlay{}, fmt.Errorf(
				"overlay: operation %q has kind %q; v1 supports only %q (replace is the deferred protocol-overlay-replace follow-up)",
				o.ID, o.Kind, KindExtend)
		}
		if !operationIDRe.MatchString(o.ID) {
			return Overlay{}, fmt.Errorf("overlay: operation id %q is not of the form deck.<slug>", o.ID)
		}
		if seen[o.ID] {
			return Overlay{}, fmt.Errorf("overlay: duplicate operation id %q", o.ID)
		}
		seen[o.ID] = true
		if strings.TrimSpace(o.Rationale) == "" {
			return Overlay{}, fmt.Errorf("overlay: operation %q has an empty rationale; a rationale is required", o.ID)
		}
		if strings.TrimSpace(o.Markdown) == "" {
			return Overlay{}, fmt.Errorf("overlay: operation %q has an empty markdown payload", o.ID)
		}
		if !hex64Re.MatchString(o.CoreSHA256) {
			return Overlay{}, fmt.Errorf(
				"overlay: operation %q has core-sha256 %q; want 64 lowercase hex characters", o.ID, o.CoreSHA256)
		}
		payload := normalizeLF(o.Markdown)
		ops = append(ops, Operation{
			ID:            o.ID,
			Kind:          o.Kind,
			Rationale:     o.Rationale,
			Markdown:      payload,
			CoreSHA256:    o.CoreSHA256,
			PayloadSHA256: Hash(payload),
		})
	}

	return Overlay{Raw: norm, SHA256: Hash(norm), Operations: ops}, nil
}

// yamlDocument returns the YAML text of the overlay file, refusing any free-form body.
//
// The documented shape is frontmatter-style — a leading `---`, the document, a closing `---` — so
// that the file reads as Markdown-adjacent. A bare YAML document with no fences is also accepted.
// Anything after the closing fence is a free-form body, which D1 forbids: a second content channel
// inside the overlay would need its own delimiter language, and a payload documenting protocol
// syntax is exactly the payload most likely to contain that delimiter.
func yamlDocument(norm string) (string, error) {
	if !strings.HasPrefix(norm, "---\n") {
		return norm, nil
	}
	rest := norm[len("---\n"):]
	end := strings.Index(rest, "\n---")
	if end < 0 {
		// Opened a frontmatter fence and never closed it. Treating the remainder as a bare document
		// would silently accept a truncated file.
		return "", errors.New("overlay: opening '---' has no closing '---'")
	}
	body := rest[end+len("\n---"):]
	if strings.TrimSpace(strings.TrimPrefix(body, "\n")) != "" {
		return "", errors.New("overlay: content after the closing '---'; a free-form body is not allowed")
	}
	return rest[:end], nil
}

// rejectAliases refuses YAML anchors and aliases.
//
// An alias makes the same bytes mean different things depending on a definition elsewhere in the
// file, which defeats reviewing an overlay by reading it. The lock hash would still be stable, so
// nothing else would catch it.
func rejectAliases(doc string) error {
	var root yaml.Node
	if err := yaml.Unmarshal([]byte(doc), &root); err != nil {
		return fmt.Errorf("overlay: %w", err)
	}
	var walk func(n *yaml.Node) error
	walk = func(n *yaml.Node) error {
		if n == nil {
			return nil
		}
		if n.Kind == yaml.AliasNode {
			return errors.New("overlay: YAML aliases are not allowed")
		}
		if n.Anchor != "" {
			return errors.New("overlay: YAML anchors are not allowed")
		}
		for _, c := range n.Content {
			if err := walk(c); err != nil {
				return err
			}
		}
		return nil
	}
	return walk(&root)
}

// normalizeLF collapses CRLF and lone CR to LF.
//
// Both the lock hash and the composed output depend on this being applied to every source. The core
// renderer already learned this the hard way: a CRLF source that is not normalized produces mixed
// endings and a render that never converges across two runs.
func normalizeLF(s string) string {
	s = strings.ReplaceAll(s, "\r\n", "\n")
	return strings.ReplaceAll(s, "\r", "\n")
}
