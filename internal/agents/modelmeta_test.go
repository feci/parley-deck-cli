package agents

import "testing"

// Golden table. The two category errors this registry exists to prevent are the last
// two groups: a gateway must never become the company, and the CLI vendor must never
// become the company.
func TestDeriveModelMetaGolden(t *testing.T) {
	cases := []struct {
		ref     string
		family  string
		company string
		route   string
		known   bool
	}{
		// Plain ids, company derived from the id itself.
		{"claude-opus-5[1m]", "Claude Opus", "Anthropic", "", true},
		{"claude-sonnet-5", "Claude Sonnet", "Anthropic", "", true},
		{"gpt-5.6-sol", "GPT", "OpenAI", "", true},
		{"glm-5p2", "GLM", "Zhipu AI", "", true},
		{"Gemini 3.6 Flash (High)", "Gemini", "Google", "", true},

		// Qualified by an explicit producer namespace.
		{"kimi-code/k3", "Kimi K", "Moonshot AI", "", true},
		{"xai/grok-4.5", "Grok", "xAI", "", true},

		// Gateway-routed: the gateway is the ROUTE, never the company.
		{"litellm/xai/grok-4.5", "Grok", "xAI", "LiteLLM", true},
		{"openrouter/anthropic/claude-opus-5", "Claude Opus", "Anthropic", "OpenRouter", true},

		// Unresolvable: say so rather than guess.
		{"internal-model-7", Unknown, Unknown, "", false},
		{"", Unknown, Unknown, "", false},
		{CLIDefault, Unknown, Unknown, "", false},
	}
	for _, c := range cases {
		got := DeriveModelMeta(c.ref)
		if got.Family != c.family || got.Company != c.company || got.Route != c.route || got.Known != c.known {
			t.Errorf("DeriveModelMeta(%q) = %+v, want family=%q company=%q route=%q known=%v",
				c.ref, got, c.family, c.company, c.route, c.known)
		}
	}
}

// The company must come from the model reference, never from the adapter that launches it.
// hermes launches GLM (Zhipu AI) and opencode launches Grok (xAI); inferring from the CLI
// name would attribute both to the wrong maker.
func TestDeriveModelMetaNeverInfersCompanyFromAdapter(t *testing.T) {
	if got := DeriveModelMeta("glm-5p2").Company; got != "Zhipu AI" {
		t.Errorf("hermes runs glm-5p2; company=%q want Zhipu AI", got)
	}
	if got := DeriveModelMeta("litellm/xai/grok-4.5").Company; got != "xAI" {
		t.Errorf("opencode runs grok via litellm; company=%q want xAI", got)
	}
	if got := DeriveModelMeta("litellm/xai/grok-4.5").Route; got != "LiteLLM" {
		t.Errorf("route=%q want LiteLLM", got)
	}
}

// Every model this machine's roster can actually launch must resolve, or the registry
// has fallen behind the adapters it ships with.
func TestDeriveModelMetaCoversBuiltinDefaults(t *testing.T) {
	for _, s := range DefaultSpecs() {
		if s.Model == "" || s.Model == CLIDefault {
			continue
		}
		if meta := DeriveModelMeta(s.Model); !meta.Known {
			t.Errorf("built-in %s ships model %q which modelmeta cannot resolve", s.ID, s.Model)
		}
	}
}

// zcode emits model ids as `zai/<model>` (no hyphen), which the z-ai/zhipu producers did
// not cover — roster show reported metadata-unknown for it (idea zcode-adapter).
func TestDeriveModelMetaResolvesZaiNamespace(t *testing.T) {
	meta := DeriveModelMeta("zai/glm-5.3")
	if !meta.Known {
		t.Fatal("zai/glm-5.3 did not resolve")
	}
	if meta.Family != "GLM" || meta.Company != "Zhipu AI" {
		t.Fatalf("zai/glm-5.3 -> family=%q company=%q, want GLM / Zhipu AI", meta.Family, meta.Company)
	}
}
