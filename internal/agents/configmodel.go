package agents

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"

	toml "github.com/pelletier/go-toml/v2"
)

// ConfigResolvedModel reports the model and/or effort an agent will actually use when its CLI
// exposes no flag for it, by reading the agent's OWN configuration — the same file the process
// reads at launch.
//
// This is deliberately NOT the same claim as an argv-bound cell, and the roster keeps the two
// apart (`model-from-config` / `effort-from-config`). The rule that a configured value the argv
// never carries must not fill the cell exists to stop a *parley-side* declaration being shown as
// if the launch enforced it. Reading the agent's own config is the opposite situation: no parley
// layer can bind the value, and the file being read is the one the process itself consults, so
// it answers "what will this agent run?" from the actual source of truth.
//
// Limitation, recorded rather than hidden: the file can change between this read and the launch,
// and none of these CLIs echo the resolved model back in their machine-readable output, so the
// value is not confirmable after a run.
//
// boundModel is the model the launch passes, when there is one; adapters whose effort is
// declared per-model need it to find the right entry.
func ConfigResolvedModel(specID, boundModel string) (model, effort string, ok bool) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", "", false
	}
	switch specID {
	case "zcode":
		// No --model and no effort flag at all: both come from zcode's own files.
		model = readJSONString(filepath.Join(home, ".zcode", "cli", "config.json"), "model", "main")
		lookup := model
		if lookup == "" {
			lookup = boundModel
		}
		effort = zcodeDefaultVariant(filepath.Join(home, ".zcode", "v2", "config.json"), lookup)
	case "kimi":
		// -m binds the model; thinking effort has no flag and lives in config.toml.
		effort = kimiThinkingEffort(filepath.Join(home, ".kimi-code", "config.toml"))
	case "opencode":
		// -m binds the model; reasoning effort is declared per-model in the config.
		effort = opencodeReasoningEffort(opencodeConfigPaths(home), boundModel)
	}
	return model, effort, model != "" || effort != ""
}

// readJSONString walks a JSON object by key path and returns a string leaf, or "".
func readJSONString(path string, keys ...string) string {
	data, err := os.ReadFile(path)
	if err != nil {
		return ""
	}
	var cur any
	if err := json.Unmarshal(data, &cur); err != nil {
		return ""
	}
	for _, k := range keys {
		m, isMap := cur.(map[string]any)
		if !isMap {
			return ""
		}
		cur = m[k]
	}
	s, _ := cur.(string)
	return s
}

// zcodeDefaultVariant finds the reasoning defaultVariant zcode records for the given model id.
// The catalogue is keyed by display name (GLM-5.3) under each provider, while model.main is a
// qualified id (zai/glm-5.3), so the last segment is matched case-insensitively.
func zcodeDefaultVariant(path, modelID string) string {
	data, err := os.ReadFile(path)
	if err != nil {
		return ""
	}
	var root map[string]any
	if err := json.Unmarshal(data, &root); err != nil {
		return ""
	}
	providers, _ := root["provider"].(map[string]any)
	want := lastSegment(modelID)
	if want == "" {
		return ""
	}
	for _, p := range providers {
		pm, _ := p.(map[string]any)
		models, _ := pm["models"].(map[string]any)
		for name, m := range models {
			if !strings.EqualFold(name, want) {
				continue
			}
			mm, _ := m.(map[string]any)
			reasoning, _ := mm["reasoning"].(map[string]any)
			if v, _ := reasoning["defaultVariant"].(string); v != "" {
				return v
			}
		}
	}
	return ""
}

// kimiThinkingEffort reads [thinking] effort from kimi's config.toml. A disabled thinking block
// reports nothing rather than an effort the run will not use.
func kimiThinkingEffort(path string) string {
	data, err := os.ReadFile(path)
	if err != nil {
		return ""
	}
	var cfg struct {
		Thinking struct {
			Enabled *bool  `toml:"enabled"`
			Effort  string `toml:"effort"`
		} `toml:"thinking"`
	}
	if err := toml.Unmarshal(data, &cfg); err != nil {
		return ""
	}
	if cfg.Thinking.Enabled != nil && !*cfg.Thinking.Enabled {
		return ""
	}
	return cfg.Thinking.Effort
}

func opencodeConfigPaths(home string) []string {
	base := os.Getenv("XDG_CONFIG_HOME")
	if base == "" {
		base = filepath.Join(home, ".config")
	}
	return []string{
		filepath.Join(base, "opencode", "opencode.jsonc"),
		filepath.Join(base, "opencode", "opencode.json"),
	}
}

// opencodeReasoningEffort finds provider.<p>.models.<model>.options.reasoningEffort for the model
// the launch binds. The model key carries the provider prefix opencode's -m uses
// (`litellm/xai/grok-4.6` -> provider `litellm`, key `xai/grok-4.6`), so both the full reference
// and the prefix-stripped remainder are tried.
func opencodeReasoningEffort(paths []string, boundModel string) string {
	if boundModel == "" {
		return ""
	}
	var root map[string]any
	for _, p := range paths {
		data, err := os.ReadFile(p)
		if err != nil {
			continue
		}
		if err := json.Unmarshal(stripJSONComments(data), &root); err == nil {
			break
		}
		root = nil
	}
	if root == nil {
		return ""
	}
	providers, _ := root["provider"].(map[string]any)
	wants := []string{boundModel}
	if i := strings.Index(boundModel, "/"); i >= 0 {
		wants = append(wants, boundModel[i+1:])
	}
	for _, p := range providers {
		pm, _ := p.(map[string]any)
		models, _ := pm["models"].(map[string]any)
		for _, want := range wants {
			m, hit := models[want]
			if !hit {
				continue
			}
			mm, _ := m.(map[string]any)
			opts, _ := mm["options"].(map[string]any)
			if v, _ := opts["reasoningEffort"].(string); v != "" {
				return v
			}
		}
	}
	return ""
}

// stripJSONComments removes // and /* */ comments from JSONC, leaving string literals intact.
// A naive regex would corrupt any value containing "//" — a base_url, for one.
func stripJSONComments(in []byte) []byte {
	out := make([]byte, 0, len(in))
	inStr, esc := false, false
	for i := 0; i < len(in); i++ {
		c := in[i]
		if inStr {
			out = append(out, c)
			switch {
			case esc:
				esc = false
			case c == '\\':
				esc = true
			case c == '"':
				inStr = false
			}
			continue
		}
		switch {
		case c == '"':
			inStr = true
			out = append(out, c)
		case c == '/' && i+1 < len(in) && in[i+1] == '/':
			for i < len(in) && in[i] != '\n' {
				i++
			}
			if i < len(in) {
				out = append(out, '\n')
			}
		case c == '/' && i+1 < len(in) && in[i+1] == '*':
			i += 2
			for i+1 < len(in) && !(in[i] == '*' && in[i+1] == '/') {
				i++
			}
			i++
		default:
			out = append(out, c)
		}
	}
	return out
}

func lastSegment(s string) string {
	if i := strings.LastIndex(s, "/"); i >= 0 {
		return s[i+1:]
	}
	return s
}
