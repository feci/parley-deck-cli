package agents

import "strings"

// modelmeta derives the model's FAMILY and COMPANY from an effective model reference.
//
// Two category errors this must never make:
//   - calling the CLI vendor the model company (hermes is not the maker of GLM),
//   - calling the routing gateway the model company (litellm is not the maker of Grok).
//
// So the reference is peeled: recognized gateway prefixes come off first and become the
// route, then the remaining producer namespace or model prefix decides family/company.
// No deck ever hand-writes these values; unknown returns "unknown" rather than a guess,
// so registry lag is visible instead of silently wrong.

// ModelMeta is the derived metadata for one effective model reference.
type ModelMeta struct {
	Family  string // e.g. "Claude Opus", "GPT", "GLM", "Grok", "Kimi K"
	Company string // the maker of the model, never the CLI vendor or the gateway
	Route   string // outermost routing layer, e.g. "LiteLLM"; empty means direct
	Known   bool   // false when nothing matched — callers add STATUS=metadata-unknown
}

// Unknown is what an unmatched reference yields.
const Unknown = "unknown"

// gateways are routing prefixes that wrap a real provider/model reference.
var gateways = map[string]string{
	"litellm":    "LiteLLM",
	"openrouter": "OpenRouter",
}

// producers maps an explicit vendor namespace (the segment before the model id in a
// qualified reference such as "xai/grok-4.5") to its company.
var producers = map[string]string{
	"anthropic": "Anthropic",
	"openai":    "OpenAI",
	"xai":       "xAI",
	"google":    "Google",
	"moonshot":  "Moonshot AI",
	"kimi-code": "Moonshot AI",
	"zhipu":     "Zhipu AI",
	"z-ai":      "Zhipu AI",
}

// prefixRule maps an unqualified model-id prefix to its family and company. Order
// matters: the first match wins, so longer/more specific prefixes come first.
type prefixRule struct {
	prefix  string
	family  string
	company string
}

var prefixRules = []prefixRule{
	{"claude-opus", "Claude Opus", "Anthropic"},
	{"claude-sonnet", "Claude Sonnet", "Anthropic"},
	{"claude-haiku", "Claude Haiku", "Anthropic"},
	{"claude", "Claude", "Anthropic"},
	{"gpt", "GPT", "OpenAI"},
	{"o3", "OpenAI o-series", "OpenAI"},
	{"o4", "OpenAI o-series", "OpenAI"},
	{"gemini", "Gemini", "Google"},
	{"glm", "GLM", "Zhipu AI"},
	{"grok", "Grok", "xAI"},
	{"kimi", "Kimi", "Moonshot AI"},
	// Bare Kimi codenames are k<digit> (k2, k3) — matched by kimiCodename below rather
	// than by a plain "k" prefix, which swallowed every id starting with k, made the
	// "kimi" rule unreachable, and misclassified unrelated models.
}

// DeriveModelMeta resolves family/company/route for an effective model reference.
func DeriveModelMeta(ref string) ModelMeta {
	ref = strings.TrimSpace(ref)
	if ref == "" || ref == CLIDefault {
		return ModelMeta{Family: Unknown, Company: Unknown}
	}

	meta := ModelMeta{}
	segments := strings.Split(ref, "/")

	// Peel a recognized gateway prefix into the route.
	if len(segments) > 1 {
		if route, ok := gateways[strings.ToLower(segments[0])]; ok {
			meta.Route = route
			segments = segments[1:]
		}
	}

	// An explicit producer namespace decides the company outright.
	if len(segments) > 1 {
		if company, ok := producers[strings.ToLower(segments[0])]; ok {
			meta.Company = company
			meta.Family = familyFromID(segments[len(segments)-1])
			if meta.Family == "" {
				meta.Family = Unknown
			}
			meta.Known = true
			return meta
		}
	}

	// Otherwise derive both from the model id itself.
	id := segments[len(segments)-1]
	for _, r := range prefixRules {
		if strings.HasPrefix(strings.ToLower(id), r.prefix) {
			return ModelMeta{Family: r.family, Company: r.company, Route: meta.Route, Known: true}
		}
	}
	if kimiCodename(id) {
		return ModelMeta{Family: "Kimi K", Company: "Moonshot AI", Route: meta.Route, Known: true}
	}

	return ModelMeta{Family: Unknown, Company: Unknown, Route: meta.Route}
}

// familyFromID resolves only the family, for references whose company came from an
// explicit namespace. Returns "" when nothing matches, so the caller can decide.
func familyFromID(id string) string {
	for _, r := range prefixRules {
		if strings.HasPrefix(strings.ToLower(id), r.prefix) {
			return r.family
		}
	}
	if kimiCodename(id) {
		return "Kimi K"
	}
	return ""
}

// kimiCodename matches a BARE Kimi codename: `k` followed by a digit (k2, k3, k2-0711).
// The former `{"k", ...}` prefix rule matched any id starting with k, which is why an
// unrelated model could come back as "Kimi K".
func kimiCodename(id string) bool {
	id = strings.ToLower(strings.TrimSpace(id))
	return len(id) >= 2 && id[0] == 'k' && id[1] >= '0' && id[1] <= '9'
}
