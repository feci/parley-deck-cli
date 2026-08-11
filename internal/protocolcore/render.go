package protocolcore

import (
	"fmt"
	"regexp"
	"strings"
)

// IdentitySlots are the six per-deck values the render carries over from the deck's own file.
//
// They are the ONLY genuinely per-deck content the fleet produced: a 36-deck measurement found
// exactly one local protocol section, and it was sync governance that belongs in the core. So the
// renderer preserves values, not prose.
type IdentitySlots struct {
	Workspace   string // "**Workspace:** …"
	Created     string // "**Created:** …"
	Transport   string // "**Transport:** …"
	RosterTable []string
	HandleTable []string
}

// slot prefixes, matched at line start.
const (
	workspacePrefix = "**Workspace:** "
	createdPrefix   = "**Created:** "
	transportPrefix = "**Transport:** "
	syncedPrefix    = "**Protocol synced:**"
)

// RenderResult is what a render produced and what it had to change to produce it.
//
// Removed is not a detail: `roster render` erasing rows without saying so was a MAJOR finding in
// the roster idea, and G1 requires this renderer to report every block it replaces or removes,
// in preview and on apply.
type RenderResult struct {
	Body      string
	Removed   []string
	Preserved []string
	// Events is the typed, source-attributed change report required by the overlay design.
	//
	// Removed/Preserved are kept because they are what the existing CLI output and tests read, but
	// they cannot express a second content source: an overlay that injects a whole section produced
	// an EMPTY report, because the old shape had no way to say "added". A second content source with
	// no way to report its own contribution is the silent-erasure class wearing a different hat.
	Events []ChangeEvent
}

// Change event kinds and sources.
const (
	EventRemoved   = "removed"
	EventRelocated = "relocated"
	EventAdded     = "added"
	EventPreserved = "preserved"

	SourceCore     = "core"
	SourceIdentity = "identity"
	SourceOverlay  = "overlay"
)

// ChangeEvent is one attributed change in a render.
type ChangeEvent struct {
	Kind   string
	Source string
	Detail string
	// OpID is the overlay operation this event is attributed to, when Source is SourceOverlay.
	OpID  string
	Lines int
}

// Render materializes a deck's COOPERATION.md from a core release plus the deck's identity slots.
//
// Deliberately a pure function of (release, prior deck body): no filesystem, no clock. The
// synced-stamp is derived from the release, so two machines holding the same release render
// byte-identical output — which is what makes G1's idempotence testable at all.
// ov may be nil, which is the canonical "this deck has no local customization" state. It is passed
// in rather than read here so that Render stays pure — the caller owns the filesystem.
func Render(rel Release, priorDeckBody string, ov *Overlay) (RenderResult, error) {
	if strings.TrimSpace(rel.Body) == "" {
		return RenderResult{}, fmt.Errorf("protocolcore: release %s has an empty body", rel.Version)
	}
	res := RenderResult{}
	// A CRLF deck must not produce a false removal for every section (each line differing only by
	// \r) nor mixed endings in the output. Normalize in, restore the deck's convention out.
	crlf := strings.Contains(priorDeckBody, "\r\n")
	priorDeckBody = strings.ReplaceAll(priorDeckBody, "\r\n", "\n")
	// A CRLF CORE must be normalized as well, or the render emits mixed endings, never converges
	// across two runs, and a CRLF deck restore turns \r\n into \r\r\n.
	rel.Body = strings.ReplaceAll(rel.Body, "\r\n", "\n")
	slots := ExtractIdentity(priorDeckBody)

	out := make([]string, 0, 512)
	lines := strings.Split(rel.Body, "\n")
	for i := 0; i < len(lines); i++ {
		line := lines[i]

		if strings.HasPrefix(line, syncedPrefix) {
			continue // regenerated below, never carried from the core
		}
		if v, ok := substituteSlot(line, slots, &res); ok {
			out = append(out, v)
			continue
		}
		if isTableHeader(line) {
			body, consumed := tableBodyFor(line, slots)
			out = append(out, line)
			i++
			if i < len(lines) && strings.HasPrefix(strings.TrimSpace(lines[i]), "| -") {
				out = append(out, lines[i])
				i++
			}
			// Drop the core's own placeholder rows; the deck's data replaces them.
			for i < len(lines) && strings.HasPrefix(strings.TrimSpace(lines[i]), "|") {
				i++
			}
			i--
			out = append(out, body...)
			if consumed != "" {
				res.Preserved = append(res.Preserved, consumed)
			}
			continue
		}
		out = append(out, line)
	}

	body := strings.Join(out, "\n")
	stamp := fmt.Sprintf("%s core %s (%s)", syncedPrefix, rel.Version, ShortHash(rel.SHA256))
	if created := findLine(body, createdPrefix); created != "" {
		body = strings.Replace(body, created, created+"\n"+stamp, 1)
	}
	// ext-1 — the terminal core/overlay boundary. Its position is the END of the normalized core
	// body BY DEFINITION rather than by lookup, which is why v1 places it without a registry and why
	// it survives a core that grows new sections: new core content lands before the boundary, so the
	// deck's own content does not move relative to it.
	if ov != nil {
		for _, op := range ov.Operations {
			body = strings.TrimRight(body, "\n") + "\n\n" + strings.TrimRight(op.Markdown, "\n") + "\n"
			res.Events = append(res.Events, ChangeEvent{
				Kind:   EventAdded,
				Source: SourceOverlay,
				OpID:   op.ID,
				Detail: "payload rendered at the terminal core/overlay boundary",
				Lines:  len(splitLines(strings.TrimRight(op.Markdown, "\n"))),
			})
		}
	}

	removed, events := changeReport(priorDeckBody, body, ov)
	res.Removed = removed
	res.Events = append(res.Events, events...)
	for _, p := range res.Preserved {
		res.Events = append(res.Events, ChangeEvent{Kind: EventPreserved, Source: SourceIdentity, Detail: p})
	}
	if crlf {
		body = strings.ReplaceAll(body, "\n", "\r\n")
	}
	res.Body = body
	return res, nil
}

// changeReport is droppedContent plus the relocation witness.
//
// A block the overlay now carries has usually MOVED — it sat mid-document in the deck and is
// rendered at the terminal boundary. Reporting that as a loss trains operators to ignore the report,
// which is worse than not having one. But the exemption has to be a PROOF of relocation, not a
// similarity heuristic: revision 2 of the consensus rejected the loose "its lines reappear
// somewhere" formulation precisely because the diff was deliberately made order-sensitive.
//
// A removed run is reclassified `relocated` only when all four hold, after LF normalization:
//  1. it is byte-identical to exactly one complete overlay payload;
//  2. those bytes occur exactly once in the prior deck;
//  3. those bytes occur exactly once in the composed output;
//  4. the output occurrence is attributed to that same operation.
//
// Anything duplicated, partial, edited or interleaved stays `removed`. The invariant is unchanged:
// an empty report means no line disappeared; it does NOT mean no meaning was lost.
func changeReport(priorBody, renderedBody string, ov *Overlay) ([]string, []ChangeEvent) {
	prior := splitLines(priorBody)
	rendered := splitLines(renderedBody)
	if len(prior) == 0 {
		return nil, nil
	}
	keep := lcsMask(prior, rendered)

	renderStamp := headerStamp(rendered)
	deckStamp := headerStamp(prior)
	forgive := renderStamp != "" && deckStamp != "" && !containsLine(prior, renderStamp)

	// Mark the one forgiven stamp line as kept so it neither counts as a loss nor breaks a run.
	if forgive {
		for i, l := range prior {
			if !keep[i] && l == deckStamp {
				keep[i] = true
				break
			}
		}
	}

	var events []ChangeEvent
	// A run is relocated only as a whole. Walk contiguous non-kept runs and test each against the
	// payloads; a matched run is marked kept so it drops out of the loss grouping entirely.
	for i := 0; i < len(prior); {
		if keep[i] {
			i++
			continue
		}
		j := i
		for j < len(prior) && !keep[j] {
			j++
		}
		// Blank-line padding at the run's edges is excluded before matching.
		//
		// This is a deliberate, narrow relaxation of "the run is byte-identical to one payload", and
		// it is recorded as such. Without it the witness essentially never fires: a block sitting
		// mid-document is bounded by blank lines that belong to the SURROUNDING sections, so the
		// contiguous non-kept run is almost always payload-plus-padding and never byte-equal to the
		// payload. A rule that cannot fire would silently demote the design back to the noisy report
		// that consensus rejected.
		//
		// What is NOT relaxed: the payload bytes themselves are still compared byte for byte, the
		// match must still be to exactly one complete payload, and it must still occur exactly once
		// on each side. Only blank lines — which carry no content — are allowed at the edges.
		lo, hi := trimBlankEdges(prior, i, j)
		run := strings.Join(prior[lo:hi], "\n")
		if op := relocationWitness(run, priorBody, renderedBody, ov); op != nil {
			for k := i; k < j; k++ {
				keep[k] = true
			}
			events = append(events, ChangeEvent{
				Kind:   EventRelocated,
				Source: SourceOverlay,
				OpID:   op.ID,
				Detail: "moved from its position in the deck to the terminal boundary; bytes identical",
				Lines:  j - i,
			})
		}
		i = j
	}

	removed := groupLosses(prior, keep)
	for _, r := range removed {
		events = append(events, ChangeEvent{Kind: EventRemoved, Source: SourceCore, Detail: r})
	}
	return removed, events
}

// relocationWitness returns the operation a removed run is proven to have moved into, or nil.
func relocationWitness(run, priorBody, renderedBody string, ov *Overlay) *Operation {
	if ov == nil || strings.TrimSpace(run) == "" {
		return nil
	}
	var match *Operation
	for i := range ov.Operations {
		payload := strings.TrimRight(ov.Operations[i].Markdown, "\n")
		if payload != strings.TrimRight(run, "\n") {
			continue
		}
		if match != nil {
			return nil // identical to more than one payload: not "exactly one"
		}
		match = &ov.Operations[i]
	}
	if match == nil {
		return nil
	}
	body := strings.TrimRight(match.Markdown, "\n")
	if strings.Count(normalizeLF(priorBody), body) != 1 || strings.Count(normalizeLF(renderedBody), body) != 1 {
		return nil // duplicated on either side: the move cannot be attributed uniquely
	}
	return match
}

// substituteSlot replaces an identity line with the deck's own value.
func substituteSlot(line string, slots IdentitySlots, res *RenderResult) (string, bool) {
	for _, s := range []struct {
		prefix string
		value  string
		name   string
	}{
		{workspacePrefix, slots.Workspace, "Workspace"},
		{createdPrefix, slots.Created, "Created"},
		{transportPrefix, slots.Transport, "Transport"},
	} {
		if strings.HasPrefix(line, s.prefix) {
			if s.value == "" {
				return line, true // deck has none; the core's own line stands
			}
			res.Preserved = append(res.Preserved, s.name)
			return s.value, true
		}
	}
	return "", false
}

func isTableHeader(line string) bool {
	t := strings.TrimSpace(line)
	return strings.HasPrefix(t, "| Agent ID") &&
		(strings.Contains(t, "Workspace") || strings.Contains(t, "Host handle"))
}

func tableBodyFor(header string, slots IdentitySlots) ([]string, string) {
	if strings.Contains(header, "Host handle") {
		if len(slots.HandleTable) > 0 {
			return slots.HandleTable, "host-handle table"
		}
		return nil, ""
	}
	if len(slots.RosterTable) > 0 {
		return slots.RosterTable, "§2 roster table"
	}
	return nil, ""
}

// ExtractIdentity pulls the six per-deck values out of an existing deck body. On a deck that has
// none (a fresh init), every slot is empty and the core's own lines stand.
func ExtractIdentity(body string) IdentitySlots {
	var s IdentitySlots
	s.Workspace = findLine(body, workspacePrefix)
	s.Created = findLine(body, createdPrefix)
	s.Transport = findLine(body, transportPrefix)
	s.RosterTable = tableRows(body, "Workspace")
	s.HandleTable = tableRows(body, "Host handle")
	return s
}

func findLine(body, prefix string) string {
	for _, l := range strings.Split(body, "\n") {
		if strings.HasPrefix(l, prefix) {
			return l
		}
	}
	return ""
}

// tableRows returns the data rows of the first "| Agent ID |" table whose header contains marker.
func tableRows(body, marker string) []string {
	lines := strings.Split(body, "\n")
	for i, l := range lines {
		t := strings.TrimSpace(l)
		if !strings.HasPrefix(t, "| Agent ID") || !strings.Contains(t, marker) {
			continue
		}
		j := i + 1
		if j < len(lines) && strings.HasPrefix(strings.TrimSpace(lines[j]), "| -") {
			j++
		}
		var rows []string
		for j < len(lines) && strings.HasPrefix(strings.TrimSpace(lines[j]), "|") {
			rows = append(rows, lines[j])
			j++
		}
		return rows
	}
	return nil
}

// droppedContent lists the deck content the render does NOT carry forward.
//
// This is a real SEQUENCE diff, not a multiset comparison, and the difference is the whole point.
// Five successive multiset variants were each defeated by a different Markdown construct — content
// under a shared heading, a line surviving in another section, a claimed-but-not-implemented
// per-section index, an indented code block equated with its unindented twin, and finally
// reordering, duplicate headings, hard line breaks and HTML-comment boundaries, where the same
// lines in a different ORDER change the meaning while every multiset stays equal.
//
// A multiset cannot see order, so no amount of further patching would have made it correct. An
// LCS diff over exact lines can: a moved line is a removal plus an addition, a lost hard break is a
// changed line, and whitespace is compared as written.
//
// What it reports, precisely: lines present in the deck that the render does not contain at that
// position. It does NOT interpret Markdown semantics.
//
// KNOWN LIMITS, stated rather than implied. This function was wrong in nine distinct ways across
// nine review cycles, so the honest posture is that it is a line-level diff and nothing more:
//   - a construct whose meaning depends on lines it does not span (a fenced block whose fence
//     survives while its content moves, a reference-style link whose definition moves) can change
//     meaning with every individual line still matched somewhere;
//   - attribution to a section is by the nearest preceding ATX heading, so a setext heading
//     (underlined with === or ---) is not a boundary;
//   - an empty report means "no line disappeared", NOT "no meaning was lost".
//
// The operator reads the diff; this report tells them where to look.
// groupLosses turns a per-line "the render still carries this" mask into the per-heading report.
//
// The stamp exemption and the relocation witness are applied by the caller, which marks the
// exempted lines as kept. Keeping that OUT of here matters: two earlier stamp rules were wrong in
// ways that only showed up because the exemption logic and the grouping logic were entangled —
// "first unmatched stamp-prefixed line" let genuine prose inherit the forgiveness when the deck
// already carried the current stamp, and "the line after **Created:**" did the same when prose sat
// in that slot.
func groupLosses(prior []string, keep []bool) []string {
	type group struct {
		heading string
		lines   int
	}
	var groups []group
	idx := map[string]int{}
	current := headerSection
	for i, l := range prior {
		if h := heading(l); h != "" {
			current = h
		}
		if keep[i] {
			continue
		}
		if j, ok := idx[current]; ok {
			groups[j].lines++
			continue
		}
		idx[current] = len(groups)
		groups = append(groups, group{heading: current, lines: 1})
	}
	out := make([]string, 0, len(groups))
	for _, g := range groups {
		unit := "lines"
		if g.lines == 1 {
			unit = "line"
		}
		out = append(out, fmt.Sprintf("%s — %d %s not carried forward", g.heading, g.lines, unit))
	}
	return out
}

const headerSection = "(document header)"

// splitLines splits on newlines and strips only the CR of a CRLF pair. Nothing else is normalized:
// trailing spaces are a Markdown hard line break, and leading spaces make a code block.
func splitLines(s string) []string {
	if s == "" {
		return nil
	}
	out := strings.Split(s, "\n")
	for i := range out {
		// Only a CRLF pair's carriage return is line-ending noise. A lone terminal CR on the LAST
		// element has no following newline, so stripping it would erase real content silently.
		if i < len(out)-1 {
			out[i] = strings.TrimSuffix(out[i], "\r")
		}
	}
	return out
}

// lcsMask returns, for each line of a, whether it belongs to a longest common subsequence with b
// — i.e. whether the render still carries it, in order.
//
// Hirschberg's algorithm: linear space. The obvious full DP table is O(n*m) MEMORY, which codex-1
// correctly flagged as unbounded — a 100k-line document would allocate tens of gigabytes and take
// the process down. A protocol is normally ~1-2k lines, but "normally" is not a bound, and a
// correctness-critical report must not be the thing that crashes on a big input.
func lcsMask(a, b []string) []bool {
	keep := make([]bool, len(a))
	if len(a) == 0 || len(b) == 0 {
		return keep
	}
	hirschberg(a, b, 0, keep)
	return keep
}

// hirschberg marks the LCS of a against b, offsetting marks by base into keep.
func hirschberg(a, b []string, base int, keep []bool) {
	switch {
	case len(a) == 0 || len(b) == 0:
		return
	case len(a) == 1:
		for _, x := range b {
			if x == a[0] {
				keep[base] = true
				return
			}
		}
		return
	}
	mid := len(a) / 2
	left := lcsRow(a[:mid], b)
	right := lcsRowReverse(a[mid:], b)
	var best int32 = -1
	split := 0
	for j := 0; j <= len(b); j++ {
		if v := left[j] + right[len(b)-j]; v > best {
			best, split = v, j
		}
	}
	hirschberg(a[:mid], b[:split], base, keep)
	hirschberg(a[mid:], b[split:], base+mid, keep)
}

// lcsRow is the last row of the LCS DP for a against b, in O(len(b)) space.
func lcsRow(a, b []string) []int32 {
	prev := make([]int32, len(b)+1)
	cur := make([]int32, len(b)+1)
	for i := 0; i < len(a); i++ {
		cur[0] = 0
		for j := 1; j <= len(b); j++ {
			if a[i] == b[j-1] {
				cur[j] = prev[j-1] + 1
			} else if prev[j] >= cur[j-1] {
				cur[j] = prev[j]
			} else {
				cur[j] = cur[j-1]
			}
		}
		prev, cur = cur, prev
	}
	return prev
}

// lcsRowReverse is lcsRow over the reversed sequences.
func lcsRowReverse(a, b []string) []int32 {
	ra := make([]string, len(a))
	for i := range a {
		ra[i] = a[len(a)-1-i]
	}
	rb := make([]string, len(b))
	for i := range b {
		rb[i] = b[len(b)-1-i]
	}
	return lcsRow(ra, rb)
}

// generatedStampRe matches the stamp this renderer PRODUCES: `**Protocol synced:** core <v> (<h>)`.
//
// The exemption needs proof that the deck's slot line really was a generated stamp, not merely
// that something sits after `**Created:**`. Without that proof, genuine stamp-prefixed prose
// occupying the slot in a deck that never had a generated stamp is forgiven and deleted with no
// removal report (codex-1, round 9).
//
// A legacy stamp in some older format does NOT match, so it is reported as not carried forward.
// That is the intended side: a data-loss report should err loud, and a legacy stamp genuinely is
// being replaced.
var generatedStampRe = regexp.MustCompile(`^\*\*Protocol synced:\*\* core \S+ \([0-9a-f]+\)\s*$`)

// headerStamp returns the deck's generated stamp line — the slot immediately after
// `**Created:**`, and only when that line is recognizably a generated stamp.
func headerStamp(lines []string) string {
	for i, l := range lines {
		if !strings.HasPrefix(l, createdPrefix) {
			continue
		}
		if i+1 < len(lines) && generatedStampRe.MatchString(strings.TrimRight(lines[i+1], " \t")) {
			return lines[i+1]
		}
		return ""
	}
	return ""
}

// containsLine reports whether an exact line is present.
func containsLine(lines []string, want string) bool {
	for _, l := range lines {
		if l == want {
			return true
		}
	}
	return false
}

// heading recognizes ANY ATX heading level, so a #### subsection is its own section boundary and
// its losses are not masked by the parent.
func heading(line string) string {
	t := strings.TrimRight(line, " \t")
	for _, p := range []string{"# ", "## ", "### ", "#### ", "##### ", "###### "} {
		if strings.HasPrefix(t, p) {
			return t
		}
	}
	return ""
}

// LCSMaskForTest exposes the diff to package-external tests. The algorithm is correctness-critical
// and a linear-space rewrite is exactly where an off-by-one hides, so it is verified against an
// independent full-DP implementation rather than only through the command surface.
func LCSMaskForTest(a, b []string) []bool { return lcsMask(a, b) }

// trimBlankEdges narrows [i,j) to exclude leading and trailing blank lines.
//
// "Blank" is a line that is empty or whitespace only. Nothing inside the range is touched: an
// interior blank line is part of the content and must match.
func trimBlankEdges(lines []string, i, j int) (int, int) {
	for i < j && strings.TrimSpace(lines[i]) == "" {
		i++
	}
	for j > i && strings.TrimSpace(lines[j-1]) == "" {
		j--
	}
	return i, j
}
