package telemetry

// lean-organizer D.4: the kimi telemetry case, built from fixtures captured live
// this run (wire usage.record shape), the honest coverage-none path, and the
// structured-argv detection for the new default.

import (
	"os"
	"testing"
)

func TestKimiUsageRecordParsesWireFixture(t *testing.T) {
	data, err := os.ReadFile("testdata/kimi-stream.txt")
	if err != nil {
		t.Skipf("fixture missing: %v", err)
	}
	// The live-captured fixture carries both emitted shapes: the plain wire
	// record (inputOther 357 / cacheRead 24064 / output 27) and the
	// role-envelope-wrapped record (129 / 30111 / 38).
	type expected struct{ input, cacheRead, output int64 }
	want := []expected{{357 + 24064, 24064, 27}, {129 + 30111, 30111, 38}}
	saw := 0
	for _, line := range splitLines(string(data)) {
		c := NewCollector("kimi", true)
		c.consume([]byte(line))
		u := c.usage
		if u.InputTokens == nil {
			continue
		}
		if saw >= len(want) {
			t.Fatalf("more usage lines parsed than the fixture documents")
		}
		w := want[saw]
		saw++
		if *u.InputTokens != w.input {
			t.Errorf("input = inputOther+inputCacheRead, got %d want %d", *u.InputTokens, w.input)
		}
		if u.CacheReadTokens == nil || *u.CacheReadTokens != w.cacheRead {
			t.Errorf("cache read verbatim, got %v want %d", u.CacheReadTokens, w.cacheRead)
		}
		if u.OutputTokens == nil || *u.OutputTokens != w.output {
			t.Errorf("output verbatim, got %v want %d", u.OutputTokens, w.output)
		}
		if u.CostBasis != "unavailable" || u.Coverage != "reported" {
			t.Errorf("tokens only, no dollar claim; got basis=%s coverage=%s", u.CostBasis, u.Coverage)
		}
	}
	if saw != len(want) {
		t.Fatalf("expected %d kimi usage.record lines to parse, got %d", len(want), saw)
	}
}

func TestKimiCoverageNoneWhenNoUsageEmitted(t *testing.T) {
	// kimi 0.42.0 emits usage to its wire log but NOT to stdout: a stream with no
	// usage.record must downgrade to coverage "none" — never a fabricated number.
	c := NewCollector("kimi", true)
	for _, line := range []string{
		`{"role":"meta","type":"system.version","version":"0.42.0"}`,
		`{"role":"assistant","content":"PONG"}`,
		`{"role":"meta","type":"session.resume_hint","session_id":"s","command":"kimi -r s"}`,
	} {
		c.consume([]byte(line))
	}
	u, _, _ := c.Result()
	if u.Coverage == "reported" {
		t.Errorf("no usage record on stdout must not claim reported coverage; got %+v", u)
	}
}

func TestKimiAdapterUsesStructuredArgvDetection(t *testing.T) {
	args := []string{"--output-format", "stream-json", "-m", "kimi-code/k3", "-p", "hello"}
	if !StructuredArgs(args) {
		t.Errorf("kimi's structured argv must be detected as structured")
	}
}

func TestKimiUsageIgnoresForeignShapes(t *testing.T) {
	c := NewCollector("kimi", true)
	c.consume([]byte(`{"role":"assistant","content":"hi"}`))
	if c.usage.InputTokens != nil {
		t.Errorf("an assistant content line is not usage")
	}
	c2 := NewCollector("claude", true)
	c2.consume([]byte(`{"type":"usage.record","usage":{"inputOther":1}}`))
	if c2.usage.InputTokens != nil {
		t.Errorf("usage.record is kimi-shaped; the claude case must not adopt it")
	}
}

func splitLines(s string) []string {
	var out []string
	start := 0
	for i := 0; i < len(s); i++ {
		if s[i] == '\n' {
			out = append(out, s[start:i])
			start = i + 1
		}
	}
	if start < len(s) {
		out = append(out, s[start:])
	}
	return out
}
