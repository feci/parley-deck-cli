package agents

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

// Composite agent display names (idea: composite-agent-naming-and-roster-reinit).
//
// A display name is the self-documenting form `family-model-effort` (optionally
// with a `-N` collision instance), e.g. `claude-opus4.8-max`, `codex-gpt5.5-xhigh`,
// `agy-gemini3.5flash-high`. It is DERIVED for rendering from a roster entry's
// family + model label + reasoning; it is never a stable identity (the roster ID
// is), and never stored as truth. All functions here are pure.
//
// Grammar (path-safe: dots only between alphanumerics, so never ".." or an edge dot):
//
//	display-name := family "-" model "-" effort [ "-" instance ]
//	family       := [a-z0-9]+
//	model        := [a-z0-9]+ ("." [a-z0-9]+)*
//	effort       := low|medium|high|xhigh|max|ultracode|clidefault
//	instance     := [2-9][0-9]*

// EffortVocabulary is the closed set of effort tokens. Parsing is fail-closed:
// a name whose effort token is outside this set is an error, never a guess. An
// effort token can never be all-digits, which is what lets the parser tell an
// instance suffix apart from the effort unambiguously.
var EffortVocabulary = []string{"low", "medium", "high", "xhigh", "max", "ultracode", "clidefault"}

var effortSet = func() map[string]bool {
	m := make(map[string]bool, len(EffortVocabulary))
	for _, e := range EffortVocabulary {
		m[e] = true
	}
	return m
}()

var (
	familyRe = regexp.MustCompile(`^[a-z0-9]+$`)
	modelRe  = regexp.MustCompile(`^[a-z0-9]+(\.[a-z0-9]+)*$`)
	tierRe   = regexp.MustCompile(`\s*\(([^)]*)\)\s*$`)
	nonAlnum = regexp.MustCompile(`[^a-z0-9]+`)
	dotRuns  = regexp.MustCompile(`\.{2,}`)
)

// ParsedName is the decomposition of a composite display name.
type ParsedName struct {
	Family   string
	Model    string
	Effort   string
	Instance int // 0 when no instance suffix is present; >= 2 otherwise
}

// SanitizeFamily lowercases and keeps only [a-z0-9] (families never carry dots).
func SanitizeFamily(s string) string {
	return nonAlnum.ReplaceAllString(strings.ToLower(strings.TrimSpace(s)), "")
}

// SanitizeModelToken lowercases, deletes every character outside [a-z0-9.],
// collapses dot runs to one, and strips leading/trailing dots — so a human model
// label becomes a path-safe token that preserves version dots: "Opus 4.8" ->
// "opus4.8", "GPT-5.5" -> "gpt5.5", "GLM 5.2" -> "glm5.2", "K3" -> "k3". The
// caller strips a parenthesized tier first (see StripParenTier) when relevant.
func SanitizeModelToken(s string) string {
	s = strings.ToLower(strings.TrimSpace(s))
	// delete every char not in [a-z0-9.]
	var b strings.Builder
	for _, r := range s {
		if (r >= 'a' && r <= 'z') || (r >= '0' && r <= '9') || r == '.' {
			b.WriteRune(r)
		}
	}
	out := dotRuns.ReplaceAllString(b.String(), ".")
	return strings.Trim(out, ".")
}

// NormalizeEffortToken maps a raw reasoning/effort value to a vocabulary token
// (e.g. "cli-default" -> "clidefault", "High" -> "high"); ok is false when the
// result is not in EffortVocabulary. It never invents a token.
func NormalizeEffortToken(s string) (string, bool) {
	t := nonAlnum.ReplaceAllString(strings.ToLower(strings.TrimSpace(s)), "")
	return t, effortSet[t]
}

// StripParenTier splits a trailing parenthesized qualifier off a model label,
// e.g. "Gemini 3.5 Flash (High)" -> ("Gemini 3.5 Flash", "High"). Used for agy,
// whose reasoning tier lives in the model label rather than a separate flag. When
// there is no trailing "(...)", tier is "".
func StripParenTier(label string) (base, tier string) {
	if m := tierRe.FindStringSubmatch(label); m != nil {
		return strings.TrimSpace(tierRe.ReplaceAllString(label, "")), strings.TrimSpace(m[1])
	}
	return strings.TrimSpace(label), ""
}

// Compose builds and validates a composite display name from already-derived
// tokens. family and model are validated against their grammars, effort must be
// a vocabulary member, and instance (0 = none) must be >= 2. It never guesses:
// an out-of-charset token is an error, not a silently-repaired string.
func Compose(family, model, effort string, instance int) (string, error) {
	if !familyRe.MatchString(family) {
		return "", fmt.Errorf("naming: family %q is not [a-z0-9]+", family)
	}
	if !modelRe.MatchString(model) {
		return "", fmt.Errorf("naming: model %q is not a dotted [a-z0-9] token", model)
	}
	if !effortSet[effort] {
		return "", fmt.Errorf("naming: effort %q is not in the vocabulary %v", effort, EffortVocabulary)
	}
	name := family + "-" + model + "-" + effort
	if instance != 0 {
		if instance < 2 {
			return "", fmt.Errorf("naming: instance %d must be >= 2", instance)
		}
		name += "-" + strconv.Itoa(instance)
	}
	return name, nil
}

// Parse decomposes a composite display name, fail-closed (right-to-left): a
// trailing all-digit token is the instance; the next-from-right must be a
// vocabulary effort; the single remaining middle token is the model (dots
// allowed); the first token is the family. Anything else is an error — legacy
// roster IDs (`claude`, `claude-1`) are NOT composite names and are rejected here.
func Parse(name string) (ParsedName, error) {
	var zero ParsedName
	tokens := strings.Split(name, "-")
	instance := 0
	// A trailing all-digit token is the instance suffix (effort is never all-digit).
	if len(tokens) >= 4 && isAllDigits(tokens[len(tokens)-1]) {
		n, err := strconv.Atoi(tokens[len(tokens)-1])
		if err != nil || n < 2 {
			return zero, fmt.Errorf("naming: bad instance suffix in %q", name)
		}
		instance = n
		tokens = tokens[:len(tokens)-1]
	}
	if len(tokens) != 3 {
		return zero, fmt.Errorf("naming: %q is not a composite display name (family-model-effort)", name)
	}
	family, model, effort := tokens[0], tokens[1], tokens[2]
	if !familyRe.MatchString(family) {
		return zero, fmt.Errorf("naming: family %q is not [a-z0-9]+", family)
	}
	if !modelRe.MatchString(model) {
		return zero, fmt.Errorf("naming: model %q is not a dotted [a-z0-9] token", model)
	}
	if !effortSet[effort] {
		return zero, fmt.Errorf("naming: effort %q is not in the vocabulary", effort)
	}
	return ParsedName{Family: family, Model: model, Effort: effort, Instance: instance}, nil
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
