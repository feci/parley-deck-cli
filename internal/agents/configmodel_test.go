package agents

import (
	"os"
	"path/filepath"
	"testing"
)

// write creates path with content, making parent directories as needed.
func write(t *testing.T, path, content string) {
	t.Helper()
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}
}

func fakeHome(t *testing.T) string {
	t.Helper()
	home := t.TempDir()
	t.Setenv("HOME", home)
	t.Setenv("XDG_CONFIG_HOME", filepath.Join(home, ".config"))
	return home
}

func TestZcodeResolvesModelAndEffortFromItsOwnConfig(t *testing.T) {
	home := fakeHome(t)
	write(t, filepath.Join(home, ".zcode", "cli", "config.json"),
		`{"model":{"main":"zai/glm-5.3","lite":"zai/glm-5-turbo"}}`)
	write(t, filepath.Join(home, ".zcode", "v2", "config.json"),
		`{"provider":{"builtin:zai":{"models":{"GLM-5.3":{"reasoning":{"enabled":true,"variants":["high","low","max"],"defaultVariant":"max"}}}}}}`)

	model, effort, ok := ConfigResolvedModel("zcode", "")
	if !ok || model != "zai/glm-5.3" || effort != "max" {
		t.Fatalf("got model=%q effort=%q ok=%v; want zai/glm-5.3 / max / true", model, effort, ok)
	}
}

// The catalogue is keyed by display name (GLM-5.3) while model.main is a qualified id
// (zai/glm-5.3). Matching the whole string finds nothing and would silently report no effort.
func TestZcodeEffortMatchesCatalogueByLastSegment(t *testing.T) {
	home := fakeHome(t)
	write(t, filepath.Join(home, ".zcode", "cli", "config.json"), `{"model":{"main":"zai/glm-5.3"}}`)
	write(t, filepath.Join(home, ".zcode", "v2", "config.json"),
		`{"provider":{"p":{"models":{"zai/glm-5.3":{"reasoning":{"defaultVariant":"low"}}}}}}`)
	if _, effort, _ := ConfigResolvedModel("zcode", ""); effort != "" {
		t.Fatalf("full-id catalogue key must not match; got effort=%q", effort)
	}
}

func TestMissingAgentConfigResolvesNothing(t *testing.T) {
	fakeHome(t)
	for _, id := range []string{"zcode", "kimi", "opencode"} {
		if model, effort, ok := ConfigResolvedModel(id, "litellm/xai/grok-4.6"); ok {
			t.Fatalf("%s: absent config must not resolve; got model=%q effort=%q", id, model, effort)
		}
	}
}

func TestUnknownAdapterNeverResolves(t *testing.T) {
	home := fakeHome(t)
	write(t, filepath.Join(home, ".zcode", "cli", "config.json"), `{"model":{"main":"zai/glm-5.3"}}`)
	if _, _, ok := ConfigResolvedModel("claude", ""); ok {
		t.Fatal("only adapters with a known config location may resolve")
	}
}

func TestKimiThinkingEffort(t *testing.T) {
	home := fakeHome(t)
	path := filepath.Join(home, ".kimi-code", "config.toml")

	write(t, path, "default_model = \"kimi-code/k3\"\n\n[thinking]\nenabled = true\neffort = \"max\"\n")
	if _, effort, ok := ConfigResolvedModel("kimi", "kimi-code/k3"); !ok || effort != "max" {
		t.Fatalf("got effort=%q ok=%v; want max/true", effort, ok)
	}

	// Disabled thinking must report nothing rather than an effort the run will not apply.
	write(t, path, "[thinking]\nenabled = false\neffort = \"max\"\n")
	if _, effort, ok := ConfigResolvedModel("kimi", "kimi-code/k3"); ok || effort != "" {
		t.Fatalf("disabled thinking resolved effort=%q ok=%v", effort, ok)
	}
}

// kimi's model IS argv-bound, so the resolver must not also claim it — otherwise the roster
// could show a config model where the launch passes a different one.
func TestKimiResolverReportsNoModel(t *testing.T) {
	home := fakeHome(t)
	write(t, filepath.Join(home, ".kimi-code", "config.toml"),
		"default_model = \"kimi-code/other\"\n[thinking]\nenabled = true\neffort = \"high\"\n")
	if model, _, _ := ConfigResolvedModel("kimi", "kimi-code/k3"); model != "" {
		t.Fatalf("kimi resolver must leave model to argv; got %q", model)
	}
}

func TestOpencodeReasoningEffortForBoundModel(t *testing.T) {
	home := fakeHome(t)
	write(t, filepath.Join(home, ".config", "opencode", "opencode.jsonc"), `{
  // a line comment, and a value containing a double slash below
  "provider": {
    "litellm": {
      "options": { "baseURL": "https://gw.example.com/v1" },
      "models": {
        "xai/grok-4.6": { "name": "Grok 4.6", "options": { "reasoningEffort": "xhigh" } },
        "xai/grok-4.5": { "options": { "reasoningEffort": "high" } }
      }
    }
  }
}`)
	if _, effort, ok := ConfigResolvedModel("opencode", "litellm/xai/grok-4.6"); !ok || effort != "xhigh" {
		t.Fatalf("got effort=%q ok=%v; want xhigh/true", effort, ok)
	}
	// A model with no declaration must stay unknown, not inherit a sibling's effort.
	if _, effort, ok := ConfigResolvedModel("opencode", "litellm/xai/grok-9"); ok || effort != "" {
		t.Fatalf("undeclared model resolved effort=%q ok=%v", effort, ok)
	}
	// Without a bound model there is no per-model entry to read.
	if _, effort, _ := ConfigResolvedModel("opencode", ""); effort != "" {
		t.Fatalf("unbound model resolved effort=%q", effort)
	}
}

// A URL in a string value contains "//". A comment stripper that ignores string state deletes
// the rest of that line and the config stops parsing.
func TestStripJSONCommentsPreservesDoubleSlashInsideStrings(t *testing.T) {
	in := []byte(`{"url":"https://x/y","a":1 /* block */, "b":2 // trailing
}`)
	got := string(stripJSONComments(in))
	if !contains(got, `"https://x/y"`) {
		t.Fatalf("string value corrupted: %s", got)
	}
	if contains(got, "block") || contains(got, "trailing") {
		t.Fatalf("comments survived: %s", got)
	}
}

func contains(s, sub string) bool {
	return len(s) >= len(sub) && (func() bool {
		for i := 0; i+len(sub) <= len(s); i++ {
			if s[i:i+len(sub)] == sub {
				return true
			}
		}
		return false
	})()
}
