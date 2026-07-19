package agents

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

// Composite agent display names (idea: composite-agent-naming-and-roster-reinit).
//
// A display name is the self-documenting form `family_model_effort` (optionally
// with a `_N` collision instance), e.g. `claude_opus-4.8-1m_max`,
// `codex_gpt-5.6-sol_xHigh`, `agy_gemini-3.5-flash_high`. Separators (user
// decision): `_` separates the three MEANINGS (family, model, effort); `-`
// separates WORDS within a section; `.` keeps version numbers natural. The name
// is DERIVED for rendering from a roster entry's family + model label + reasoning;
// it is never a stable identity (the roster ID is) and never stored as truth. All
// functions here are pure.
//
// Grammar (path-safe — dots only between alphanumerics, so never ".." or an edge dot):
//
//	display-name := family "_" model "_" effort [ "_" instance ]
//	family       := word ("-" word)*
//	model        := word ("-" word)*
//	effort       := one of the vocabulary tokens (rendered in its display form)
//	word         := [a-z0-9]+ ("." [a-z0-9]+)*
//	instance     := [2-9][0-9]*

// EffortVocabulary is the closed set of effort tokens (normalized, lowercase).
// Parsing is fail-closed: an effort outside this set is an error, never a guess.
var EffortVocabulary = []string{"low", "medium", "high", "xhigh", "max", "ultracode", "clidefault"}

// effortDisplay maps each normalized token to its camelCase display form (user
// decision: `xHigh`, `cliDefault`).
var effortDisplay = map[string]string{
	"low":        "low",
	"medium":     "medium",
	"high":       "high",
	"xhigh":      "xHigh",
	"max":        "max",
	"ultracode":  "ultracode",
	"clidefault": "cliDefault",
}

var (
	// section = one or more words joined by '-', each word alphanumeric with
	// optional internal version dots. No leading/trailing '-' or '.', no "..".
	sectionRe = regexp.MustCompile(`^[a-z0-9]+(\.[a-z0-9]+)*(-[a-z0-9]+(\.[a-z0-9]+)*)*$`)
	// wordSplit splits a raw label into words on any run of non-[a-z0-9.] chars,
	// so spaces, hyphens, brackets, parens and slashes all delimit words while
	// version dots are preserved inside a word.
	wordSplit = regexp.MustCompile(`[^a-z0-9.]+`)
	nonAlnum  = regexp.MustCompile(`[^a-z0-9]+`)
	dotRuns   = regexp.MustCompile(`\.{2,}`)
	tierRe    = regexp.MustCompile(`\s*\(([^)]*)\)\s*$`)
)

// ParsedName is the decomposition of a composite display name. Effort is the
// normalized (lowercase) token; use EffortDisplayForm to render it.
type ParsedName struct {
	Family   string
	Model    string
	Effort   string
	Instance int // 0 when no instance suffix is present; >= 2 otherwise
}

// SanitizeSection turns a human label into a section token: lowercase, split into
// words on any non-[a-z0-9.] run, drop empties, collapse dot runs and strip edge
// dots per word, join words with '-'. So "GPT-5.6 Sol" -> "gpt-5.6-sol",
// "Opus 4.8 1m" -> "opus-4.8-1m", "Gemini 3.5 Flash" -> "gemini-3.5-flash",
// "GLM 5.2" -> "glm-5.2", "K3" -> "k3". (Deriving from a raw model id like
// "claude-opus-4-8[1m]" yields garbage — always derive the model from a label.)
func SanitizeSection(s string) string {
	s = strings.ToLower(strings.TrimSpace(s))
	parts := wordSplit.Split(s, -1)
	words := make([]string, 0, len(parts))
	for _, p := range parts {
		p = strings.Trim(dotRuns.ReplaceAllString(p, "."), ".")
		if p != "" {
			words = append(words, p)
		}
	}
	return strings.Join(words, "-")
}

// NormalizeEffortToken maps a raw reasoning/effort value to a vocabulary token
// (e.g. "cli-default" -> "clidefault", "xHigh" -> "xhigh"); ok is false when the
// result is not in the vocabulary. It never invents a token.
func NormalizeEffortToken(s string) (string, bool) {
	t := nonAlnum.ReplaceAllString(strings.ToLower(strings.TrimSpace(s)), "")
	_, ok := effortDisplay[t]
	return t, ok
}

// EffortDisplayForm returns the camelCase display form for a normalized token
// (falling back to the token itself if unknown).
func EffortDisplayForm(token string) string {
	if d, ok := effortDisplay[token]; ok {
		return d
	}
	return token
}

// StripParenTier splits a trailing parenthesized qualifier off a model label,
// e.g. "Gemini 3.5 Flash (High)" -> ("Gemini 3.5 Flash", "High"). Used for agy,
// whose reasoning tier lives in the model label rather than a separate flag.
func StripParenTier(label string) (base, tier string) {
	if m := tierRe.FindStringSubmatch(label); m != nil {
		return strings.TrimSpace(tierRe.ReplaceAllString(label, "")), strings.TrimSpace(m[1])
	}
	return strings.TrimSpace(label), ""
}

// Compose builds and validates a composite display name from already-sanitized
// family and model sections plus a raw effort value. It never guesses: an
// out-of-grammar section or unknown effort is an error, not a repaired string.
func Compose(family, model, effortRaw string, instance int) (string, error) {
	if !sectionRe.MatchString(family) {
		return "", fmt.Errorf("naming: family %q is not a valid section", family)
	}
	if !sectionRe.MatchString(model) {
		return "", fmt.Errorf("naming: model %q is not a valid section", model)
	}
	tok, ok := NormalizeEffortToken(effortRaw)
	if !ok {
		return "", fmt.Errorf("naming: effort %q is not in the vocabulary %v", effortRaw, EffortVocabulary)
	}
	name := family + "_" + model + "_" + EffortDisplayForm(tok)
	if instance != 0 {
		if instance < 2 {
			return "", fmt.Errorf("naming: instance %d must be >= 2", instance)
		}
		name += "_" + strconv.Itoa(instance)
	}
	return name, nil
}

// Parse decomposes a composite display name, fail-closed. It splits on the
// section separator '_'; a trailing all-digit fourth section is the instance; the
// remaining three sections are family, model, effort. family and model must match
// the section grammar; effort must normalize to a vocabulary token. Legacy roster
// IDs (`claude`, `claude-1`) are NOT composite names and are rejected here.
func Parse(name string) (ParsedName, error) {
	var zero ParsedName
	sections := strings.Split(name, "_")
	instance := 0
	if len(sections) == 4 && isAllDigits(sections[3]) {
		n, err := strconv.Atoi(sections[3])
		if err != nil || n < 2 {
			return zero, fmt.Errorf("naming: bad instance suffix in %q", name)
		}
		instance = n
		sections = sections[:3]
	}
	if len(sections) != 3 {
		return zero, fmt.Errorf("naming: %q is not a composite display name (family_model_effort)", name)
	}
	family, model, effortSec := sections[0], sections[1], sections[2]
	if !sectionRe.MatchString(family) {
		return zero, fmt.Errorf("naming: family %q is not a valid section", family)
	}
	if !sectionRe.MatchString(model) {
		return zero, fmt.Errorf("naming: model %q is not a valid section", model)
	}
	tok, ok := NormalizeEffortToken(effortSec)
	if !ok {
		return zero, fmt.Errorf("naming: effort %q is not in the vocabulary", effortSec)
	}
	// Fail-closed canonical check (review MINOR): re-compose and require byte
	// equality, so non-canonical spellings that normalize to the same value
	// (`x-high` -> xhigh, instance `02` -> 2, lowercase `xhigh` vs `xHigh`) are
	// rejected rather than silently accepted with a lossy round-trip.
	canonical, cerr := Compose(family, model, tok, instance)
	if cerr != nil || canonical != name {
		return zero, fmt.Errorf("naming: %q is not in canonical form (want %q)", name, canonical)
	}
	return ParsedName{Family: family, Model: model, Effort: tok, Instance: instance}, nil
}

// RenderDisplayName derives the composite display name for a roster family and its
// resolved spec, from the model label (falling back to Model), the reasoning
// effort, and — for agy, whose tier lives in the model label — the parenthesized
// tier used as the effort when reasoning is cli-default. Returns an error only when
// the model/effort cannot be sanitized into the grammar; the caller then falls back
// to the raw roster ID. This is the single source of the display everywhere it is
// shown (§2 table, TUI, digests).
func RenderDisplayName(family string, spec Spec) (string, error) {
	label := strings.TrimSpace(spec.ModelLabel)
	if label == "" {
		label = spec.Model
	}
	base, effortRaw := label, spec.Reasoning
	// ONLY agy carries its reasoning tier inside the model label "(High)"; for any
	// other family a parenthesized qualifier stays part of the model, so we do not
	// strip/substitute it (review MINOR, codex-1: e.g. "Model (Preview)" must not
	// move "Preview" into the effort).
	if family == "agy" {
		var tier string
		if base, tier = StripParenTier(label); tier != "" {
			if norm, ok := NormalizeEffortToken(effortRaw); !ok || norm == "clidefault" {
				effortRaw = tier
			}
		}
	}
	return Compose(SanitizeSection(family), SanitizeSection(base), effortRaw, 0)
}

func isAllDigits(s string) bool {
	if s == "" {
		return false
	}
	for _, r := range s {
		if r < '0' || r > '9' {
			return false
		}
	}
	return true
}
