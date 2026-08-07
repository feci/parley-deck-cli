package protocolcore

import (
	"fmt"
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
}

// Render materializes a deck's COOPERATION.md from a core release plus the deck's identity slots.
//
// Deliberately a pure function of (release, prior deck body): no filesystem, no clock. The
// synced-stamp is derived from the release, so two machines holding the same release render
// byte-identical output — which is what makes G1's idempotence testable at all.
func Render(rel Release, priorDeckBody string) (RenderResult, error) {
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
	res.Removed = droppedContent(priorDeckBody, body)
	if crlf {
		body = strings.ReplaceAll(body, "\n", "\r\n")
	}
	res.Body = body
	return res, nil
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
// Two earlier attempts were wrong in instructive ways. Comparing HEADINGS missed content living
// under a heading the core also has. Comparing a flat SET of trimmed lines then lost section
// context and multiplicity: an identical line elsewhere in the render made a dropped line look
// kept, and three dropped lines that happened to be equal reported as one. Both leave semantic
// erasure silent, which is what G1 exists to prevent — and what destroyed a deck's local section
// during the 2026-08-06 fleet sync.
//
// So the comparison is per-section and multiplicity-aware: within the section a line belongs to,
// a deck line counts as carried only if the render still has an unconsumed copy of it.
func droppedContent(priorBody, renderedBody string) []string {
	if strings.TrimSpace(priorBody) == "" {
		return nil
	}
	renderCounts := map[string]int{}
	for _, l := range strings.Split(renderedBody, "\n") {
		if t := strings.TrimSpace(l); t != "" {
			renderCounts[t]++
		}
	}
	type group struct {
		heading string
		lines   int
	}
	var groups []group
	idx := map[string]int{}
	current := "(document header)"
	for _, l := range strings.Split(priorBody, "\n") {
		t := strings.TrimSpace(l)
		if h := heading(l); h != "" {
			current = h
		}
		if t == "" {
			continue
		}
		// The synced stamp is regenerated on every render by design; reporting the old one as
		// "lost project content" would cry wolf on every single version bump.
		if strings.HasPrefix(t, syncedPrefix) {
			continue
		}
		if renderCounts[t] > 0 {
			renderCounts[t]--
			continue
		}
		if i, ok := idx[current]; ok {
			groups[i].lines++
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

func heading(line string) string {
	t := strings.TrimRight(line, " \t")
	if strings.HasPrefix(t, "## ") || strings.HasPrefix(t, "### ") {
		return t
	}
	return ""
}
