package runner

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"parley-deck-cli/internal/protocol"
)

var contextReasoningFenceTags = []string{"think", "thought", "thinking"}

// SanitizeForContext removes a closed set of complete hidden-reasoning fences
// before content is reused in prompts or derived context. It preserves malformed
// unclosed fences as literal text and is not secret redaction.
func SanitizeForContext(input string) string {
	output := input
	for _, tag := range contextReasoningFenceTags {
		output = removeTaggedBlocks(output, tag)
	}
	return output
}

func removeTaggedBlocks(input, tag string) string {
	open := "<" + tag + ">"
	close := "</" + tag + ">"
	var b strings.Builder
	offset := 0
	for {
		startRel := indexASCIIEqualFold(input[offset:], open)
		if startRel < 0 {
			b.WriteString(input[offset:])
			return b.String()
		}
		start := offset + startRel
		b.WriteString(input[offset:start])
		afterOpen := start + len(open)
		endRel := indexASCIIEqualFold(input[afterOpen:], close)
		if endRel < 0 {
			b.WriteString(input[start:])
			return b.String()
		}
		offset = afterOpen + endRel + len(close)
	}
}

func indexASCIIEqualFold(value, needle string) int {
	if needle == "" {
		return 0
	}
	for i := 0; i+len(needle) <= len(value); i++ {
		if asciiEqualFold(value[i:i+len(needle)], needle) {
			return i
		}
	}
	return -1
}

func asciiEqualFold(a, b string) bool {
	if len(a) != len(b) {
		return false
	}
	for i := 0; i < len(a); i++ {
		ca := asciiLower(a[i])
		cb := asciiLower(b[i])
		if ca != cb {
			return false
		}
	}
	return true
}

func asciiLower(value byte) byte {
	if value >= 'A' && value <= 'Z' {
		return value + ('a' - 'A')
	}
	return value
}

func writeRoundIndex(idea protocol.IdeaStatus, roundLabel string, results []Result) (string, error) {
	roundDir := filepath.Join(idea.Path, roundLabel)
	indexPath := filepath.Join(roundDir, "_index.md")
	if err := os.MkdirAll(roundDir, 0o755); err != nil {
		return indexPath, err
	}
	data := BuildRoundIndex(idea, roundLabel, results)
	if err := os.WriteFile(indexPath, []byte(data), 0o644); err != nil {
		return indexPath, err
	}
	return indexPath, nil
}

func BuildRoundIndex(idea protocol.IdeaStatus, roundLabel string, results []Result) string {
	entries := make([]roundIndexEntry, 0, len(results))
	for _, result := range results {
		if result.AgentID == "" {
			continue
		}
		entries = append(entries, buildRoundIndexEntry(result))
	}
	sort.SliceStable(entries, func(i, j int) bool {
		return entries[i].AgentID < entries[j].AgentID
	})

	var b strings.Builder
	fmt.Fprintf(&b, "---\nidea: %s\nround: %s\nartifact: round-index\nderived: true\ngenerated-by: parley\ntoken-heuristic: bytes_div_4\n---\n\n", idea.Slug, roundLabel)
	fmt.Fprintf(&b, "# Round Index: %s\n\n", roundLabel)
	fmt.Fprintln(&b, "This is a runner-owned derived artifact. Source participant artifacts are not modified.")
	fmt.Fprintln(&b)
	fmt.Fprintln(&b, "- Sanitizer: context-only hidden-reasoning fence removal, not secret redaction.")
	fmt.Fprintln(&b, "- Supported fences: `<think>`, `<thought>`, `<thinking>`.")
	fmt.Fprintln(&b, "- Approx tokens heuristic: `(sanitized_bytes + 3) / 4`.")
	fmt.Fprintln(&b)
	fmt.Fprintln(&b, "| Agent | Status | Approx tokens | H2 sections | Artifact |")
	fmt.Fprintln(&b, "| --- | --- | ---: | --- | --- |")
	for _, entry := range entries {
		fmt.Fprintf(&b, "| %s | %s | %d | %s | %s |\n",
			escapeTable(entry.AgentID),
			escapeTable(entry.Status),
			entry.ApproxTokens,
			escapeTable(strings.Join(entry.SectionNames(), "; ")),
			escapeTable(entry.ArtifactName),
		)
	}
	for _, entry := range entries {
		fmt.Fprintf(&b, "\n## %s\n\n", entry.AgentID)
		fmt.Fprintf(&b, "- Status: %s\n", entry.Status)
		fmt.Fprintf(&b, "- Artifact: `%s`\n", entry.ArtifactName)
		fmt.Fprintf(&b, "- Approx tokens: %d\n", entry.ApproxTokens)
		if entry.Note != "" {
			fmt.Fprintf(&b, "- Note: %s\n", entry.Note)
		}
		if len(entry.Sections) == 0 {
			fmt.Fprintln(&b, "- Sections: no recognized H2 sections")
			continue
		}
		fmt.Fprintln(&b, "- Sections:")
		for _, section := range entry.Sections {
			summary := section.Summary
			if summary == "" {
				summary = "no summary text"
			}
			fmt.Fprintf(&b, "  - %s: %s\n", section.Title, summary)
		}
	}
	b.WriteByte('\n')
	return b.String()
}

type roundIndexEntry struct {
	AgentID      string
	Status       string
	ArtifactName string
	ApproxTokens int
	Sections     []roundIndexSection
	Note         string
}

func (e roundIndexEntry) SectionNames() []string {
	if len(e.Sections) == 0 {
		return []string{"none"}
	}
	names := make([]string, len(e.Sections))
	for i, section := range e.Sections {
		names[i] = section.Title
	}
	return names
}

type roundIndexSection struct {
	Title   string
	Summary string
}

func buildRoundIndexEntry(result Result) roundIndexEntry {
	entry := roundIndexEntry{
		AgentID:      result.AgentID,
		Status:       roundIndexStatus(result),
		ArtifactName: filepath.Base(result.OutputPath),
	}
	if entry.ArtifactName == "." || entry.ArtifactName == string(filepath.Separator) {
		entry.ArtifactName = "(none)"
	}
	if result.ExitError != "" {
		entry.Note = result.ExitError
	}
	if result.Skipped && result.SkipReason != "" {
		entry.Note = result.SkipReason
	}
	if result.OutputPath == "" {
		return entry
	}

	data, err := os.ReadFile(result.OutputPath)
	if err != nil {
		if entry.Note == "" {
			entry.Note = "artifact read failed: " + err.Error()
		}
		return entry
	}
	sanitized := SanitizeForContext(string(data))
	entry.ApproxTokens = approxTokens(sanitized)
	entry.Sections = extractH2Sections(sanitized)
	return entry
}

func roundIndexStatus(result Result) string {
	switch {
	case result.Warning != "":
		return "warning"
	case result.ExitError != "":
		return "failed"
	case result.Skipped:
		return "skipped"
	case result.ArtifactOK:
		return "ok"
	default:
		return "missing"
	}
}

func extractH2Sections(markdown string) []roundIndexSection {
	lines := strings.Split(markdown, "\n")
	var sections []roundIndexSection
	current := -1
	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		if strings.HasPrefix(trimmed, "## ") {
			title := strings.TrimSpace(strings.TrimPrefix(trimmed, "## "))
			if title == "" {
				continue
			}
			sections = append(sections, roundIndexSection{Title: title})
			current = len(sections) - 1
			continue
		}
		if current < 0 || sections[current].Summary != "" || trimmed == "" {
			continue
		}
		if strings.HasPrefix(trimmed, "#") || strings.HasPrefix(trimmed, "---") {
			continue
		}
		sections[current].Summary = trimSummary(trimmed)
	}
	return sections
}

func approxTokens(value string) int {
	if value == "" {
		return 0
	}
	return (len(value) + 3) / 4
}

func trimSummary(value string) string {
	const maxSummaryChars = 180
	value = strings.Join(strings.Fields(value), " ")
	runes := []rune(value)
	if len(runes) <= maxSummaryChars {
		return value
	}
	return strings.TrimSpace(string(runes[:maxSummaryChars])) + "..."
}

func escapeTable(value string) string {
	value = strings.ReplaceAll(value, "\n", " ")
	value = strings.ReplaceAll(value, "|", `\|`)
	if strings.TrimSpace(value) == "" {
		return "none"
	}
	return value
}
