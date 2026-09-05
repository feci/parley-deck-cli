// Package protocolpacket renders the phase-scoped protocol packet ratified by
// ideas/meta-protocol-change-phase-packet-and-fixup-budget/FINAL.md and refined by
// ideas/meta-protocol-change-evidence-first-efficiency/FINAL.md (D4).
//
// One renderer, called by every prompt builder. It never embeds, bundles or freezes protocol
// text: the body is rendered from the live resolved authority (see source.go) and bound to
// its hash. Blocks are verbatim; the omission index is complete; unknown applicability,
// unresolved dependencies, malformed input and unexpected phase/track/transport/flag values
// fall back VISIBLY to the full protocol. Missing authority and detected secrets refuse.
//
// The default is full context with a shadow packet computed for audit only. An optimized
// packet body is produced only on explicit request (Request.Optimize), which is the
// experimental input of the ratified trial, not an enabled release.
package protocolpacket

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

// Context modes carried in the launch attestation.
const (
	ModeFull         = "full"          // default: the complete source is the body
	ModePacket       = "packet"        // explicit experimental optimized packet
	ModeFullFallback = "full-fallback" // optimization requested, refused for a stated reason
	ModeRefused      = "refused"       // no context emitted: secret detected or authority missing
)

// Index classifications.
const (
	ClassAlways      = "always"
	ClassConditional = "conditional"
	ClassOnDemand    = "on-demand"
	ClassUnknown     = "unknown"
)

// Known request vocabularies. Anything outside them is an unexpected value and falls back.
var (
	KnownTracks     = []string{"fast", "standard", "deliberation"}
	KnownTransports = []string{"local-dir", "github-pr", "gitlab-mr"}
	KnownFlags      = []string{"protocol_change", "strict_gate", "auto_implement", "pipeline"}
)

// PreambleLocator names content that precedes the first heading. It is identity text and is
// always included; it is not a map entry because it has no heading to key on.
const PreambleLocator = "(preamble)"

// Request describes the launch the packet is for.
type Request struct {
	Phase     int      `json:"phase"`
	Track     string   `json:"track"`
	Transport string   `json:"transport"` // empty: taken from the source header
	IdeaSlug  string   `json:"idea_slug,omitempty"`
	Flags     []string `json:"flags,omitempty"`
	// Optimize is the explicit experimental input. False (the default) renders full context
	// and a shadow packet record; true renders the optimized packet when every guard passes.
	Optimize bool `json:"optimize"`
}

// Attestation is what a launch record carries about the protocol context it was given.
type Attestation struct {
	ContextMode    string `json:"context_mode"`
	SourceSHA256   string `json:"source_sha256"`
	PacketSHA256   string `json:"packet_sha256,omitempty"`
	FallbackReason string `json:"fallback_reason,omitempty"`
}

// Block is one verbatim heading-delimited slice of the source.
type Block struct {
	Locator   string `json:"locator"`
	Level     int    `json:"level"`
	StartLine int    `json:"start_line"`
	EndLine   int    `json:"end_line"`
	SHA256    string `json:"sha256"`
	Text      string `json:"-"`
}

// BlockRecord is one row of the complete included/omitted index.
type BlockRecord struct {
	Block
	Classification string `json:"classification"`
	Included       bool   `json:"included"`
	Why            string `json:"why"`
	Trigger        string `json:"trigger,omitempty"`
}

// Shadow is the audit record computed when the body is full context.
type Shadow struct {
	PacketSHA256        string `json:"packet_sha256,omitempty"`
	WouldFallbackReason string `json:"would_fallback_reason,omitempty"`
	PacketBytes         int    `json:"packet_bytes"`
	SourceBytes         int    `json:"source_bytes"`
	IncludedBlocks      int    `json:"included_blocks"`
	OmittedBlocks       int    `json:"omitted_blocks"`
}

// Context is the renderer's complete output.
type Context struct {
	Attestation
	Source   SourceInfo    `json:"source"`
	Request  Request       `json:"request"`
	Index    []BlockRecord `json:"index"`
	Shadow   *Shadow       `json:"shadow,omitempty"`
	Problems []string      `json:"problems,omitempty"`
	BodyPath string        `json:"body_path,omitempty"`
	Body     string        `json:"-"`
}

// SourceInfo is the JSON-visible description of the resolved authority.
type SourceInfo struct {
	Path      string `json:"path"`
	Role      string `json:"role"`
	Authority string `json:"authority"`
	Transport string `json:"transport,omitempty"`
	Drift     string `json:"drift,omitempty"`
	Bytes     int    `json:"bytes"`
}

// Hash is the content address used for source and packet bytes.
func Hash(s string) string {
	sum := sha256.Sum256([]byte(s))
	return hex.EncodeToString(sum[:])
}

// secretShapes are detection-only. A match REFUSES emission; nothing is redacted, because a
// redacted body presented under the original source hash is a false attestation (ALT-05).
var secretShapes = []struct {
	name string
	re   *regexp.Regexp
}{
	{"bearer-token", regexp.MustCompile(`(?i)\bbearer\s+[A-Za-z0-9._\-]{16,}`)},
	{"labeled-secret", regexp.MustCompile(`(?i)\b(api[_-]?key|access[_-]?key|secret[_-]?key|private[_-]?key|passwd)\s*[:=]\s*['"]?[A-Za-z0-9._\-/+]{12,}`)},
	{"openai-key", regexp.MustCompile(`\bsk-[A-Za-z0-9_\-]{16,}`)},
	{"github-token", regexp.MustCompile(`\bgh[pousr]_[A-Za-z0-9]{20,}`)},
	{"slack-token", regexp.MustCompile(`\bxox[baprs]-[A-Za-z0-9\-]{10,}`)},
	{"aws-access-key", regexp.MustCompile(`\bAKIA[0-9A-Z]{16}\b`)},
	{"jwt", regexp.MustCompile(`\beyJ[A-Za-z0-9_\-]{10,}\.[A-Za-z0-9_\-]{10,}\.[A-Za-z0-9_\-]{5,}`)},
	{"private-key-block", regexp.MustCompile(`-----BEGIN [A-Z ]*PRIVATE KEY-----`)},
}

// DetectSecret returns the name of the first credential-shaped token found, or "".
func DetectSecret(body string) string {
	for _, s := range secretShapes {
		if s.re.MatchString(body) {
			return s.name
		}
	}
	return ""
}

// Parse splits the source into heading-delimited blocks. Every line belongs to exactly one
// block; headings inside fenced code are not headings. Locators are the heading line with
// whitespace collapsed; a repeated heading gets " #2", " #3" so locators stay unique.
func Parse(raw string) []Block {
	lines := strings.Split(raw, "\n")
	var blocks []Block
	seen := map[string]int{}
	start := 0
	locator := PreambleLocator
	level := 0
	flush := func(end int) {
		if end <= start && locator == PreambleLocator {
			return // no preamble content
		}
		text := strings.Join(lines[start:end], "\n")
		if locator == PreambleLocator && strings.TrimSpace(text) == "" {
			return
		}
		blocks = append(blocks, Block{
			Locator: locator, Level: level, StartLine: start + 1, EndLine: end, SHA256: Hash(text), Text: text,
		})
	}
	inFence := false
	for i, line := range lines {
		t := strings.TrimSpace(line)
		if strings.HasPrefix(t, "```") || strings.HasPrefix(t, "~~~") {
			inFence = !inFence
			continue
		}
		if inFence {
			continue
		}
		lvl := headingLevel(line)
		if lvl == 0 {
			continue
		}
		flush(i)
		start = i
		level = lvl
		locator = strings.Join(strings.Fields(line), " ")
		seen[locator]++
		if n := seen[locator]; n > 1 {
			locator = fmt.Sprintf("%s #%d", locator, n)
		}
	}
	flush(len(lines))
	return blocks
}

func headingLevel(line string) int {
	if !strings.HasPrefix(line, "#") {
		return 0
	}
	n := 0
	for n < len(line) && line[n] == '#' {
		n++
	}
	if n > 6 || n >= len(line) || line[n] != ' ' {
		return 0
	}
	return n
}

// Build renders the context for one launch. It is pure: the caller resolved the source and
// loaded the map. A nil map is a missing map and falls back.
func Build(src Source, m *Map, req Request) Context {
	ctx := Context{Request: req}
	ctx.SourceSHA256 = Hash(src.Raw)
	ctx.Source = SourceInfo{Path: src.Path, Role: src.Role, Authority: src.Authority, Transport: src.Transport, Drift: src.Drift, Bytes: len(src.Raw)}

	if strings.TrimSpace(src.Raw) == "" {
		ctx.ContextMode = ModeRefused
		ctx.FallbackReason = "authority:empty-source"
		return ctx
	}
	if name := DetectSecret(src.Raw); name != "" {
		// No body, no packet hash: refusing is the only emission that does not disclose.
		ctx.ContextMode = ModeRefused
		ctx.FallbackReason = "secret-detected:" + name
		return ctx
	}
	if req.Transport == "" {
		req.Transport = src.Transport
		ctx.Request.Transport = src.Transport
	}
	// §7 binds in every phase of a protocol-change idea; derive the flag from the slug so a
	// caller cannot forget it.
	if strings.HasPrefix(req.IdeaSlug, "meta-protocol-change-") && !contains(req.Flags, "protocol_change") {
		req.Flags = append(append([]string(nil), req.Flags...), "protocol_change")
		ctx.Request.Flags = req.Flags
	}

	var reasons []string
	if req.Phase < 0 || req.Phase > 8 {
		reasons = append(reasons, "unknown-phase:"+strconv.Itoa(req.Phase))
	}
	if !contains(KnownTracks, req.Track) {
		reasons = append(reasons, "unknown-track:"+req.Track)
	}
	if !contains(KnownTransports, req.Transport) {
		reasons = append(reasons, "unknown-transport:"+req.Transport)
	}
	for _, f := range req.Flags {
		if !contains(KnownFlags, f) {
			reasons = append(reasons, "unknown-flag:"+f)
		}
	}
	if src.Drift != "" {
		reasons = append(reasons, "drift:"+src.Drift)
	}

	blocks := Parse(src.Raw)
	byLocator := map[string]int{}
	for i, b := range blocks {
		byLocator[b.Locator] = i
	}
	if m == nil {
		reasons = append(reasons, "no-applicability-map")
	} else {
		for _, r := range m.Blocks {
			if _, ok := byLocator[r.Locator]; !ok {
				reasons = append(reasons, "stale-map:"+r.Locator)
			}
		}
	}

	// Classification and direct applicability.
	included := make([]bool, len(blocks))
	records := make([]BlockRecord, len(blocks))
	for i, b := range blocks {
		rec := BlockRecord{Block: b}
		switch {
		case b.Locator == PreambleLocator:
			rec.Classification, rec.Included, rec.Why = ClassAlways, true, "identity preamble"
		case m == nil:
			rec.Classification, rec.Included, rec.Why = ClassUnknown, true, "no applicability map"
		default:
			rule, ok := m.rule(b.Locator)
			if !ok {
				rec.Classification, rec.Included, rec.Why = ClassUnknown, true, "no applicability entry"
				reasons = append(reasons, "unclassified:"+b.Locator)
				break
			}
			rec.Trigger = rule.Trigger
			switch rule.Include {
			case IncludeAlways:
				rec.Classification, rec.Included, rec.Why = ClassAlways, true, "never cut"
			case IncludeWhen:
				rec.Classification = ClassConditional
				if why := rule.applies(req); why != "" {
					rec.Included, rec.Why = true, why
				} else {
					rec.Why = "not applicable to this launch"
				}
			case IncludeOnDemand:
				rec.Classification, rec.Why = ClassOnDemand, "read on demand"
			}
		}
		included[i] = rec.Included
		records[i] = rec
	}

	// Dependency closure: a block pulls in what it requires; an unresolvable requirement is an
	// unresolved dependency and falls back.
	if m != nil {
		changed := true
		for changed {
			changed = false
			for i, b := range blocks {
				if !included[i] {
					continue
				}
				rule, ok := m.rule(b.Locator)
				if !ok {
					continue
				}
				for _, dep := range rule.Requires {
					j, ok := byLocator[dep]
					if !ok {
						reasons = appendUnique(reasons, "dependency-unresolved:"+dep+" (required by "+b.Locator+")")
						continue
					}
					if !included[j] {
						included[j] = true
						records[j].Included = true
						records[j].Why = "required by " + b.Locator
						changed = true
					}
				}
			}
		}
	}
	ctx.Index = records

	packetBody := renderPacket(src, ctx.SourceSHA256, req, records)
	if req.Optimize {
		if len(reasons) == 0 {
			ctx.ContextMode = ModePacket
			ctx.Body = packetBody
		} else {
			ctx.ContextMode = ModeFullFallback
			ctx.Body = src.Raw
			ctx.FallbackReason = strings.Join(reasons, "; ")
		}
	} else {
		ctx.ContextMode = ModeFull
		ctx.Body = src.Raw
		ctx.Shadow = &Shadow{SourceBytes: len(src.Raw)}
		if len(reasons) == 0 {
			ctx.Shadow.PacketSHA256 = Hash(packetBody)
			ctx.Shadow.PacketBytes = len(packetBody)
		} else {
			ctx.Shadow.WouldFallbackReason = strings.Join(reasons, "; ")
		}
		for _, r := range records {
			if r.Included {
				ctx.Shadow.IncludedBlocks++
			} else {
				ctx.Shadow.OmittedBlocks++
			}
		}
	}
	ctx.Problems = reasons
	ctx.PacketSHA256 = Hash(ctx.Body)
	return ctx
}

// renderPacket lays out the verbatim included blocks in source order, then the complete
// omission index. Nothing here paraphrases protocol text.
func renderPacket(src Source, sourceHash string, req Request, records []BlockRecord) string {
	var b strings.Builder
	flags := strings.Join(req.Flags, ",")
	if flags == "" {
		flags = "-"
	}
	fmt.Fprintf(&b, "<!-- parley-protocol-packet context_mode=%s source_sha256=%s phase=%d track=%s transport=%s flags=%s -->\n",
		ModePacket, sourceHash, req.Phase, req.Track, req.Transport, flags)
	fmt.Fprintf(&b, "# Protocol packet — phase %d, track %s, transport %s\n\n", req.Phase, req.Track, req.Transport)
	fmt.Fprintf(&b, "Excerpt of the authoritative protocol `%s` (sha256 `%s`). Every block below is verbatim\n", src.Path, sourceHash)
	b.WriteString("source text in source order. The omission index at the end lists every block not reproduced\n")
	b.WriteString("here, with the trigger that requires reading the full source. On any doubt, read the full source.\n\n---\n\n")
	for _, r := range records {
		if !r.Included {
			continue
		}
		b.WriteString(strings.TrimRight(r.Text, "\n"))
		b.WriteString("\n\n")
	}
	b.WriteString("---\n\n## Packet omission index\n\n")
	b.WriteString("Every source block appears exactly once: verbatim above, or in this table.\n\n")
	b.WriteString("| Locator | Lines | sha256 | Classification | Trigger |\n| --- | --- | --- | --- | --- |\n")
	for _, r := range records {
		if r.Included {
			continue
		}
		fmt.Fprintf(&b, "| `%s` | %d-%d | %s | %s | %s |\n", r.Locator, r.StartLine, r.EndLine, r.SHA256[:12], r.Classification, cell(r.Trigger))
	}
	return b.String()
}

func cell(s string) string {
	s = strings.ReplaceAll(s, "|", "\\|")
	return strings.Join(strings.Fields(s), " ")
}

func contains(list []string, v string) bool {
	for _, x := range list {
		if x == v {
			return true
		}
	}
	return false
}

func appendUnique(list []string, v string) []string {
	if contains(list, v) {
		return list
	}
	return append(list, v)
}

// Locators lists the source's block locators in order. Used by `packet check` and by the
// map author to enumerate what must be classified.
func Locators(raw string) []string {
	blocks := Parse(raw)
	out := make([]string, 0, len(blocks))
	for _, b := range blocks {
		if b.Locator != PreambleLocator {
			out = append(out, b.Locator)
		}
	}
	return out
}
