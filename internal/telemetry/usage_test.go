package telemetry

import (
	"fmt"
	"strings"
	"testing"
)

func feed(c *Collector, text string) { _, _ = c.Writer("stdout").Write([]byte(text)) }

func TestClaudeAuthoritativeSummaryAndReportedIdentity(t *testing.T) {
	c := NewCollector("claude", true)
	feed(c, `{"type":"result","usage":{"input_tokens":10,"output_tokens":20,"cache_read_input_tokens":5},"total_cost_usd":0.125,"modelUsage":{"actually-reported-model":{"costUSD":9}}}`+"\n")
	u, observation, failure := c.Result()
	if failure != "" || u.CostUSD == nil || *u.CostUSD != .125 || *u.InputTokens != 10 || *u.OutputTokens != 20 || *u.ReportedModel != "actually-reported-model" {
		t.Fatalf("bad summary: %+v / %s", u, failure)
	}
	if u.CostBasis != "cli-estimate" || u.TotalTokens != nil || observation.FirstActivityMS == nil {
		t.Fatal("invalid provenance or observation")
	}
	second, _, _ := c.Result()
	if *second.CostUSD != .125 {
		t.Fatal("result read double-counted modelUsage cost")
	}
}

func TestStructuredUsagePreservesModelContextSuffix(t *testing.T) {
	for adapter, envelope := range map[string]string{
		"claude":  `{"type":"result","usage":{"input_tokens":1},"modelUsage":{"claude-opus-5[1m]":{}}}`,
		"codex":   `{"type":"turn.completed","usage":{"input_tokens":1},"model":"claude-opus-5[1m]"}`,
		"generic": `{"type":"result","usage":{"input_tokens":1},"model":"claude-opus-5[1m]"}`,
	} {
		t.Run(adapter, func(t *testing.T) {
			c := NewCollector(adapter, true)
			feed(c, envelope+"\n")
			u, _, _ := c.Result()
			if u.ReportedModel == nil || *u.ReportedModel != "claude-opus-5[1m]" {
				t.Fatalf("provider identity lost: %+v", u)
			}
		})
	}
}

func TestPrettyJSONAndChunkBoundaries(t *testing.T) {
	c := NewCollector("claude", true)
	text := "{\n  \"type\": \"result\",\n  \"usage\": {\"input_tokens\": 12},\n  \"total_cost_usd\": 0\n}\n"
	for _, character := range []byte(text) {
		_, _ = c.Writer("stdout").Write([]byte{character})
	}
	u, _, _ := c.Result()
	if u.InputTokens == nil || *u.InputTokens != 12 || u.CostUSD == nil || *u.CostUSD != 0 {
		t.Fatalf("chunked pretty JSON lost: %+v", u)
	}
}

func TestCodexTokensAreNotInventedCost(t *testing.T) {
	c := NewCollector("codex", true)
	feed(c, `{"type":"turn.completed","usage":{"input_tokens":100,"cached_input_tokens":30,"output_tokens":40}}`+"\n")
	u, _, _ := c.Result()
	if u.InputTokens == nil || *u.CacheReadTokens != 30 || u.CostUSD != nil || u.ReportedModel != nil || u.TotalTokens != nil {
		t.Fatalf("bad Codex accounting: %+v", u)
	}
}

func TestACPContextUtilizationNeverBecomesBilling(t *testing.T) {
	c := NewCollector("kimi", true)
	feed(c, `{"type":"usage","used":5000,"size":128000}`+"\n")
	u, _, _ := c.Result()
	if u.InputTokens != nil || u.TotalTokens != nil || u.CostUSD != nil || u.Coverage != "unavailable" {
		t.Fatalf("invented billed usage: %+v", u)
	}
}

func TestUnstructuredPromptLookingLikeUsageIsNotProviderEvidence(t *testing.T) {
	c := NewCollector("claude", false)
	feed(c, `{"type":"result","usage":{"input_tokens":100},"total_cost_usd":99}`)
	u, observation, _ := c.Result()
	if u.InputTokens != nil || u.CostUSD != nil || observation.StdoutBytes == 0 {
		t.Fatal("plain model output treated as provider usage")
	}
}

func TestOpencodeStepsDeduplicateAndUnknownCostStaysUnknown(t *testing.T) {
	c := NewCollector("opencode", true)
	step := `{"type":"step_finish","part":{"id":"part-one","tokens":{"input":10,"output":4,"cache":{"read":2}},"cost":0.25}}` + "\n"
	feed(c, step+step)
	feed(c, `{"type":"step_finish","part":{"id":"part-two","tokens":{"input":20,"output":5},"cost":0.5}}`+"\n")
	u, _, _ := c.Result()
	if u.CostUSD == nil || *u.CostUSD != .75 || *u.InputTokens != 30 || *u.OutputTokens != 9 || u.CacheReadTokens != nil {
		t.Fatalf("bad step accounting: %+v", u)
	}
	feed(c, `{"type":"step_finish","part":{"id":"part-three","tokens":{"input":1,"output":1}}}`+"\n")
	u, _, _ = c.Result()
	if u.CostUSD != nil {
		t.Fatal("missing step cost treated as zero")
	}
}

func TestConflictingOrMissingStepIdentityCannotUndercount(t *testing.T) {
	for _, second := range []string{`{"type":"step_finish","part":{"id":"one","cost":2}}`, `{"type":"step_finish","part":{"cost":2}}`} {
		c := NewCollector("opencode", true)
		feed(c, `{"type":"step_finish","part":{"id":"one","cost":1}}`+"\n"+second+"\n")
		u, _, _ := c.Result()
		if u.CostUSD != nil || u.Coverage != "ambiguous-step-identity" {
			t.Fatalf("ambiguous total accepted: %+v", u)
		}
	}
}

func TestProviderErrorOverridesMisleadingSuccessSubtype(t *testing.T) {
	c := NewCollector("claude", true)
	feed(c, `{"type":"result","subtype":"success","is_error":true,"api_error_status":429,"total_cost_usd":0}`)
	u, _, failure := c.Result()
	if failure != "rate-limit" || u.CostUSD == nil || *u.CostUSD != 0 {
		t.Fatalf("error was success: %+v / %s", u, failure)
	}
}

func TestRecoveredIntermediateErrorDoesNotPoisonFinalSuccess(t *testing.T) {
	c := NewCollector("codex", true)
	feed(c, `{"type":"error","message":"temporary reconnect"}`+"\n"+`{"type":"turn.completed","usage":{"input_tokens":1}}`+"\n")
	_, _, failure := c.Result()
	if failure != "" {
		t.Fatal("transient error treated as final failure")
	}
}

func TestOversizedLinesBoundMemoryAndRecoverAtNextEnvelope(t *testing.T) {
	c := NewCollector("codex", true)
	feed(c, strings.Repeat("x", parserLimit+100)+"\n")
	feed(c, `{"type":"turn.completed","usage":{"input_tokens":7}}`+"\n")
	u, observation, _ := c.Result()
	if !observation.TruncatedInput || *u.InputTokens != 7 || len(c.line) > parserLimit || len(c.tail) > parserLimit {
		t.Fatal("bounded parser did not recover")
	}
	for index := 0; index < 4100; index++ {
		feed(c, fmt.Sprintf(`{"type":"ignored","id":"%d"}`+"\n", index))
	}
	if len(c.line) > parserLimit || len(c.tail) > parserLimit {
		t.Fatal("unbounded collector")
	}
}

func TestStructuredArgsAreObservedNotModified(t *testing.T) {
	for _, args := range [][]string{{"--json"}, {"--output-format", "json"}, {"--output-format=stream-json"}, {"--format", "json"}} {
		if !StructuredArgs(args) {
			t.Fatalf("missed structured mode: %v", args)
		}
	}
	for _, args := range [][]string{{"--output-format", "text"}, {"--prompt", "say --json"}, {"--output-format"}} {
		if StructuredArgs(args) {
			t.Fatalf("invented structured mode: %v", args)
		}
	}
}
