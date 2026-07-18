package agents

import "testing"

func TestSanitizeSection(t *testing.T) {
	cases := map[string]string{
		"Opus 4.8 1m":         "opus-4.8-1m",
		"GPT-5.6 Sol":         "gpt-5.6-sol",
		"gpt-5.6-sol":         "gpt-5.6-sol", // already well-formed id round-trips
		"GLM 5.2":             "glm-5.2",
		"Gemini 3.5 Flash":    "gemini-3.5-flash",
		"K3":                  "k3",
		"  ..weird..  ":       "weird",
		"a....b":              "a.b",
		"claude-opus-4-8[1m]": "claude-opus-4-8-1m", // raw id: no version dot survives -> derive from a label
	}
	for in, want := range cases {
		if got := SanitizeSection(in); got != want {
			t.Errorf("SanitizeSection(%q) = %q, want %q", in, got, want)
		}
	}
}

func TestNormalizeEffortAndDisplay(t *testing.T) {
	ok := map[string]string{"max": "max", "xhigh": "xhigh", "xHigh": "xhigh", "High": "high", "cli-default": "clidefault", "  MAX ": "max"}
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
	// camelCase display forms (user decision).
	if EffortDisplayForm("xhigh") != "xHigh" || EffortDisplayForm("clidefault") != "cliDefault" || EffortDisplayForm("max") != "max" {
		t.Error("effort display forms wrong")
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
		{"claude", "opus-4.8-1m", "max", "claude_opus-4.8-1m_max"},
		{"codex", "gpt-5.6-sol", "xhigh", "codex_gpt-5.6-sol_xHigh"},
		{"hermes", "glm-5.2", "high", "hermes_glm-5.2_high"},
		{"agy", "gemini-3.5-flash", "high", "agy_gemini-3.5-flash_high"},
		{"kimi", "k3", "max", "kimi_k3_max"},
	}
	for _, c := range cases {
		got, err := Compose(c.family, c.model, c.effort, 0)
		if err != nil || got != c.want {
			t.Errorf("Compose(%q,%q,%q) = %q,%v want %q", c.family, c.model, c.effort, got, err, c.want)
		}
	}
}

func TestComposeRejects(t *testing.T) {
	if _, err := Compose("claude", "opus-4.8", "turbo", 0); err == nil {
		t.Error("expected error for effort not in vocabulary")
	}
	// path-unsafe model sections must not validate
	if _, err := Compose("claude", ".opus", "max", 0); err == nil {
		t.Error("expected error for leading-dot model")
	}
	if _, err := Compose("claude", "opus..8", "max", 0); err == nil {
		t.Error("expected error for double-dot model")
	}
	if _, err := Compose("claude", "opus-", "max", 0); err == nil {
		t.Error("expected error for trailing-hyphen model")
	}
	if _, err := Compose("claude", "opus-4.8", "max", 1); err == nil {
		t.Error("expected error for instance 1")
	}
}

func TestComposeCollisionInstance(t *testing.T) {
	got, err := Compose("claude", "opus-4.8-1m", "max", 2)
	if err != nil || got != "claude_opus-4.8-1m_max_2" {
		t.Fatalf("got %q,%v", got, err)
	}
}

func TestParseRoundTrip(t *testing.T) {
	for _, name := range []string{"claude_opus-4.8-1m_max", "codex_gpt-5.6-sol_xHigh", "agy_gemini-3.5-flash_high", "kimi_k3_max"} {
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
	p, err := Parse("codex_gpt-5.6-sol_xHigh_2")
	if err != nil || p.Instance != 2 || p.Effort != "xhigh" || p.Model != "gpt-5.6-sol" {
		t.Fatalf("got %+v err=%v", p, err)
	}
}

func TestParseAllDigitModelIsUnambiguous(t *testing.T) {
	// A single all-digit model word must still parse: instance is detected only at
	// four sections, so `codex_530_xHigh` reads 530 as the model.
	p, err := Parse("codex_530_xHigh")
	if err != nil || p.Model != "530" || p.Effort != "xhigh" || p.Instance != 0 {
		t.Fatalf("got %+v err=%v", p, err)
	}
}

func TestRenderDisplayName(t *testing.T) {
	cases := []struct {
		family string
		spec   Spec
		want   string
	}{
		{"claude", Spec{ModelLabel: "Opus 4.8 1m", Reasoning: "max"}, "claude_opus-4.8-1m_max"},
		{"codex", Spec{Model: "gpt-5.6-sol", Reasoning: "xhigh"}, "codex_gpt-5.6-sol_xHigh"},
		{"hermes", Spec{Model: "GLM 5.2", Reasoning: "high"}, "hermes_glm-5.2_high"},
		{"agy", Spec{Model: "Gemini 3.5 Flash (High)", Reasoning: "cli-default"}, "agy_gemini-3.5-flash_high"},
		{"kimi", Spec{Model: "k3", Reasoning: "max"}, "kimi_k3_max"},
	}
	for _, c := range cases {
		got, err := RenderDisplayName(c.family, c.spec)
		if err != nil || got != c.want {
			t.Errorf("RenderDisplayName(%s) = %q,%v want %q", c.family, got, err, c.want)
		}
	}
}

func TestParseFailsClosed(t *testing.T) {
	for _, bad := range []string{
		"claude",                   // bare family (legacy id), not a composite
		"claude-1",                 // legacy roster id, not a composite
		"claude_opus-4.8-1m",       // missing effort
		"claude_opus-4.8-1m_turbo", // effort not in vocabulary
		"claude_opus-4.8-1m_max_1", // instance must be >= 2
		"claude__max",              // empty model section
		"",                         // empty
		"codex_gpt-5.6-sol_x-high", // non-canonical effort spelling (normalizes to xhigh)
		"codex_gpt-5.6-sol_xhigh",  // lowercase effort is non-canonical (want xHigh)
		"codex_530_xHigh_02",       // leading-zero instance is non-canonical
	} {
		if _, err := Parse(bad); err == nil {
			t.Errorf("Parse(%q) should fail closed", bad)
		}
	}
}
