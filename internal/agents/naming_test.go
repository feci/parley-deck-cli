package agents

import "testing"

func TestSanitizeModelToken(t *testing.T) {
	cases := map[string]string{
		"Opus 4.8":            "opus4.8",
		"GPT-5.5":             "gpt5.5",
		"GLM 5.2":             "glm5.2",
		"K3":                  "k3",
		"Gemini 3.5 Flash":    "gemini3.5flash",
		"  ..weird..  ":       "weird",          // edge dots stripped
		"a....b":              "a.b",            // dot runs collapsed -> never ".."
		"claude-opus-4-8[1m]": "claudeopus481m", // raw id: hyphens+brackets deleted, no dots -> garbage (proves we must derive from a label, not the raw id)
	}
	for in, want := range cases {
		if got := SanitizeModelToken(in); got != want {
			t.Errorf("SanitizeModelToken(%q) = %q, want %q", in, got, want)
		}
	}
}

func TestNormalizeEffortToken(t *testing.T) {
	ok := map[string]string{"max": "max", "xhigh": "xhigh", "High": "high", "cli-default": "clidefault", "  MAX ": "max"}
	for in, want := range ok {
		got, valid := NormalizeEffortToken(in)
		if !valid || got != want {
			t.Errorf("NormalizeEffortToken(%q) = %q,%v want %q,true", in, got, valid, want)
		}
	}
	for _, bad := range []string{"turbo", "", "fast", "reasoning"} {
		if _, valid := NormalizeEffortToken(bad); valid {
			t.Errorf("NormalizeEffortToken(%q) should be invalid", bad)
		}
	}
}

func TestStripParenTier(t *testing.T) {
	base, tier := StripParenTier("Gemini 3.5 Flash (High)")
	if base != "Gemini 3.5 Flash" || tier != "High" {
		t.Fatalf("got (%q,%q)", base, tier)
	}
	base, tier = StripParenTier("K3")
	if base != "K3" || tier != "" {
		t.Fatalf("got (%q,%q)", base, tier)
	}
}

func TestComposeTheFiveRoster(t *testing.T) {
	cases := []struct{ family, model, effort, want string }{
		{"claude", "opus4.8", "max", "claude-opus4.8-max"},
		{"codex", "gpt5.5", "xhigh", "codex-gpt5.5-xhigh"},
		{"hermes", "glm5.2", "high", "hermes-glm5.2-high"},
		{"agy", "gemini3.5flash", "high", "agy-gemini3.5flash-high"},
		{"kimi", "k3", "max", "kimi-k3-max"},
	}
	for _, c := range cases {
		got, err := Compose(c.family, c.model, c.effort, 0)
		if err != nil || got != c.want {
			t.Errorf("Compose(%q,%q,%q) = %q,%v want %q", c.family, c.model, c.effort, got, err, c.want)
		}
	}
}

func TestComposeRejects(t *testing.T) {
	// bad effort
	if _, err := Compose("claude", "opus4.8", "turbo", 0); err == nil {
		t.Error("expected error for effort not in vocabulary")
	}
	// family with a dot
	if _, err := Compose("cla.ude", "opus4.8", "max", 0); err == nil {
		t.Error("expected error for dotted family")
	}
	// model with edge/double dot (path-unsafe) must not validate
	if _, err := Compose("claude", ".opus", "max", 0); err == nil {
		t.Error("expected error for leading-dot model")
	}
	if _, err := Compose("claude", "opus..8", "max", 0); err == nil {
		t.Error("expected error for double-dot model")
	}
	// instance < 2
	if _, err := Compose("claude", "opus4.8", "max", 1); err == nil {
		t.Error("expected error for instance 1")
	}
}

func TestComposeCollisionInstance(t *testing.T) {
	got, err := Compose("claude", "opus4.8", "max", 2)
	if err != nil || got != "claude-opus4.8-max-2" {
		t.Fatalf("got %q,%v", got, err)
	}
}

func TestParseRoundTrip(t *testing.T) {
	for _, name := range []string{"claude-opus4.8-max", "codex-gpt5.5-xhigh", "agy-gemini3.5flash-high", "kimi-k3-max"} {
		p, err := Parse(name)
		if err != nil {
			t.Fatalf("Parse(%q) error: %v", name, err)
		}
		back, err := Compose(p.Family, p.Model, p.Effort, p.Instance)
		if err != nil || back != name {
			t.Errorf("round-trip %q -> %+v -> %q,%v", name, p, back, err)
		}
	}
}

func TestParseInstance(t *testing.T) {
	p, err := Parse("claude-opus4.8-max-2")
	if err != nil || p.Instance != 2 || p.Effort != "max" || p.Model != "opus4.8" {
		t.Fatalf("got %+v err=%v", p, err)
	}
}

func TestParseAllDigitModelIsUnambiguous(t *testing.T) {
	// A hypothetical all-digit model like "530" must still parse: the effort token
	// (never all-digits) anchors the right side, so "530" is read as the model.
	p, err := Parse("codex-530-xhigh")
	if err != nil || p.Model != "530" || p.Effort != "xhigh" || p.Instance != 0 {
		t.Fatalf("got %+v err=%v", p, err)
	}
}

func TestParseFailsClosed(t *testing.T) {
	for _, bad := range []string{
		"claude",               // bare family (legacy roster/spec id), not a composite
		"claude-1",             // legacy roster id, not a composite
		"claude-opus4.8",       // missing effort
		"claude-opus4.8-turbo", // effort not in vocabulary
		"claude-opus4.8-max-1", // instance must be >= 2
		"",                     // empty
	} {
		if _, err := Parse(bad); err == nil {
			t.Errorf("Parse(%q) should fail closed", bad)
		}
	}
}
