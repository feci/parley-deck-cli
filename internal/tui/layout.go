package tui

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
)

const (
	minTUIWidth            = 80
	defaultDashboardHeight = 36
	defaultLiveHeight      = 36
	compactDashboardHeight = 30
	compactLiveHeight      = 30
)

var (
	sectionStyle = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("39"))
	successBadge = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("42"))
	activeBadge  = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("45"))
	warnBadge    = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("214"))
	dangerBadge  = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("196"))
	neutralBadge = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("8"))
)

func tuiWidth(width int) int {
	if width < minTUIWidth {
		return minTUIWidth
	}
	return width
}

func tuiHeight(height, fallback int) int {
	if height <= 0 {
		return fallback
	}
	return height
}

func clampInt(value, minValue, maxValue int) int {
	if value < minValue {
		return minValue
	}
	if value > maxValue {
		return maxValue
	}
	return value
}

func clipLines(value string, maxLines int) string {
	if maxLines <= 0 {
		return ""
	}
	lines := strings.Split(strings.TrimRight(value, "\n"), "\n")
	if len(lines) <= maxLines {
		return strings.Join(lines, "\n")
	}
	return strings.Join(lines[:maxLines], "\n")
}

func sectionTitle(title string) string {
	return sectionStyle.Render(title)
}

func attentionBadge(value string) string {
	label := badgeLabel(value, "IDLE")
	switch strings.ToLower(label) {
	case "action", "blocked", "needs-answer", "question":
		return badge(label, warnBadge)
	case "running", "pending":
		return badge(label, activeBadge)
	case "failed", "error":
		return badge(label, dangerBadge)
	case "idle", "done", "complete", "finished":
		return badge(label, successBadge)
	default:
		return badge(label, neutralBadge)
	}
}

func stateBadge(value string) string {
	label := badgeLabel(value, "UNKNOWN")
	switch strings.ToLower(label) {
	case "finished", "complete", "answered", "found", "ok", "success":
		return badge(label, successBadge)
	case "running", "pending", "starting":
		return badge(label, activeBadge)
	case "failed", "error", "missing":
		return badge(label, dangerBadge)
	case "skipped", "unknown", "idle":
		return badge(label, neutralBadge)
	default:
		return badge(label, warnBadge)
	}
}

func riskBadge(value string) string {
	label := badgeLabel(value, "-")
	switch strings.ToLower(label) {
	case "high", "critical", "danger":
		return badge(label, dangerBadge)
	case "normal", "medium":
		return badge(label, warnBadge)
	case "low", "-":
		return badge(label, neutralBadge)
	default:
		return badge(label, neutralBadge)
	}
}

func badge(value string, style lipgloss.Style) string {
	return style.Render("[" + badgeLabel(value, "UNKNOWN") + "]")
}

func badgeLabel(value, fallback string) string {
	label := strings.TrimSpace(value)
	if label == "" {
		label = fallback
	}
	label = strings.ToUpper(strings.ReplaceAll(label, " ", "-"))
	if len(label) > 12 {
		label = label[:12]
	}
	return label
}

func oneLine(value string) string {
	return strings.Join(strings.Fields(value), " ")
}

func firstVisibleLine(value string) string {
	for _, line := range strings.Split(value, "\n") {
		line = strings.TrimSpace(line)
		if line != "" {
			return line
		}
	}
	return ""
}
