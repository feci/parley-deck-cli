package runner

// Kimi Code's structured `--output-format stream-json` mode emits role-tagged JSON
// envelope lines on stdout (verified live 2026-09-23, kimi 0.42.0, lean-organizer):
//
//	{"role":"meta","type":"system.version","version":"0.42.0"}
//	{"role":"assistant","content":"PONG"}
//	{"role":"meta","type":"session.resume_hint",...}
//
// Text consumers (the preflight PONG sentinel, consult answers, the stdout artifact
// fallback) must read the ASSISTANT CONTENT, not the envelopes. UnwrapKimiStreamJSON
// concatenates assistant contents verbatim and drops meta lines; non-JSON lines pass
// through unchanged (mixed output degrades to the plain-text path). When no
// assistant line is found the input is returned unchanged — never a guess.

import (
	"encoding/json"
	"io"
	"strings"
)

// UnwrapKimiStreamJSON extracts the readable reply from a kimi stream-json stdout
// capture. It never errors and never fabricates content.
func UnwrapKimiStreamJSON(out string) string {
	trimmed := strings.TrimSpace(out)
	if trimmed == "" || trimmed[0] != '{' {
		return out
	}
	var contents []string
	sawAssistant := false
	for _, raw := range strings.Split(out, "\n") {
		line := strings.TrimSpace(raw)
		if len(line) == 0 || line[0] != '{' {
			continue
		}
		decoder := json.NewDecoder(strings.NewReader(line))
		decoder.UseNumber()
		var event map[string]any
		if decoder.Decode(&event) != nil {
			continue
		}
		var trailing any
		if decoder.Decode(&trailing) != io.EOF {
			continue
		}
		role, _ := event["role"].(string)
		switch role {
		case "assistant":
			sawAssistant = true
			if content, ok := event["content"].(string); ok {
				contents = append(contents, content)
			}
		case "meta":
			// envelopes: dropped
		}
	}
	if !sawAssistant || len(contents) == 0 {
		return out
	}
	return strings.Join(contents, "\n")
}
