package pipeline

import (
	"crypto/sha256"
	"encoding/hex"
	"strings"
	"time"
)

// SLO is a service-level objective tracked by a watcher block (§12.11).
type SLO struct {
	Name      string `yaml:"name" json:"name"`
	Objective string `yaml:"objective" json:"objective"`
	Target    string `yaml:"target" json:"target"`
}

// Signal is a monitored input source and its breach threshold. Threshold
// expresses the BREACH condition directly (e.g. ">500", "<0.99", ">=5"): a
// breach is raised when the observed value satisfies the comparison. Class
// names the BreachClass used for remediation policy.
type Signal struct {
	Name      string `yaml:"name" json:"name"`
	Source    string `yaml:"source" json:"source"`
	Threshold string `yaml:"threshold" json:"threshold"`
	Target    string `yaml:"target" json:"target"`
	Class     string `yaml:"class" json:"class"`
}

// BreachClass categorizes a breach and declares its remediation autonomy.
// Auto-open is honored only for predeclared low-risk classes; everything else
// notifies and opens a human gate.
type BreachClass struct {
	Name     string `yaml:"name" json:"name"`
	Risk     Risk   `yaml:"risk" json:"risk"`
	AutoOpen bool   `yaml:"auto_open" json:"auto_open"`
}

// Monitoring is the MONITORING.md watcher spec (§12.11).
type Monitoring struct {
	PipelineSlug  string        `yaml:"pipeline_slug" json:"pipeline_slug"`
	SLOs          []SLO         `yaml:"slos" json:"slos"`
	Signals       []Signal      `yaml:"signals" json:"signals"`
	Destinations  []string      `yaml:"destinations" json:"destinations"`
	BreachClasses []BreachClass `yaml:"breach_classes" json:"breach_classes"`
	DedupeWindow  time.Duration `yaml:"dedupe_window" json:"dedupe_window"`
}

// Breach is one observed SLO/threshold violation.
type Breach struct {
	Signal    string    `json:"signal"`
	Target    string    `json:"target"`
	Threshold string    `json:"threshold"`
	Observed  string    `json:"observed"`
	Class     string    `json:"class"`
	At        time.Time `json:"at"`
}

// Fingerprint is the stable identity of a breach for deduplication: one ongoing
// breach cannot spawn duplicate remediation ideas (§12.11).
func (b Breach) Fingerprint() string {
	sum := sha256.Sum256([]byte(strings.Join([]string{b.Signal, b.Target, b.Threshold, b.Class}, "|")))
	return hex.EncodeToString(sum[:])[:16]
}

// DedupeBreaches collapses breaches sharing a fingerprint, keeping the first.
func DedupeBreaches(breaches []Breach) []Breach {
	seen := make(map[string]bool, len(breaches))
	out := make([]Breach, 0, len(breaches))
	for _, b := range breaches {
		fp := b.Fingerprint()
		if seen[fp] {
			continue
		}
		seen[fp] = true
		out = append(out, b)
	}
	return out
}

// CanAutoOpen reports whether a breach may auto-open a remediation idea without
// a human gate. This is allowed ONLY for a predeclared low-risk, non-production
// breach class; everything else notifies and gates (§12.11).
func (m Monitoring) CanAutoOpen(b Breach) bool {
	for _, c := range m.BreachClasses {
		if c.Name != b.Class {
			continue
		}
		if c.Risk == RiskProduction {
			return false
		}
		return c.AutoOpen && c.Risk == RiskLow
	}
	return false
}
