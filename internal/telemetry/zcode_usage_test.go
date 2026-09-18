package telemetry

import "testing"

func zcodeFeed(c *Collector, text string) { _, _ = c.Writer("stdout").Write([]byte(text)) }

// Verbatim copy of the stdout preserved at .parley-runtime/actual-zcode-stdout.txt
// whose terminal record reported every usage field unavailable: an AI SDK
// warning preamble line followed by the pretty-printed terminal envelope.
const zcodeActualEnvelope = `{
  "sessionId": "sess_eede3c76-3a87-4bb1-ad6e-efc6f2f7552b",
  "traceId": "c1fd5977-d25e-43f1-81b2-c707bf780d00",
  "turnId": "turn_c5f32c1c-a6ce-49dc-a03c-669d6cdf5dc0",
  "response": "The transcription of my authored artifact content is running now (this allocation gave me no direct file-write tool, so a helper is placing my exact text into my owned inbox file — disclosed inside the artifact itself). Once it completes I'll verify the file contents and report the artifact path with the three essential conclusions.",
  "usage": {
    "source": "provider",
    "modelRequestCount": 42,
    "inputTokens": 4119566,
    "outputTokens": 18299,
    "totalTokens": 4137865,
    "cacheReadTokens": 3996160,
    "cacheWriteTokens": 0,
    "reasoningTokens": 0,
    "webFetchRequests": 0,
    "webSearchRequests": 0
  },
  "eventCount": 11211,
  "projection": {
    "status": "idle",
    "turnCount": 1,
    "totalTokenCount": 4137865,
    "contextUsed": 122317,
    "contextWindow": 1000000
  }
}
`

const zcodeActualStdout = "AI SDK Warning System: To turn off warning logging, set the AI_SDK_LOG_WARNINGS global to false.\n" + zcodeActualEnvelope

// The discriminating regression: the actual CLI envelope must become reported
// usage. Run against the pre-change source this test fails, which is the
// counterexample showing the old parser missed the actual envelope.
func TestZcodeActualEnvelopeWithWarningPreamble(t *testing.T) {
	c := NewCollector("zcode", true)
	zcodeFeed(c, zcodeActualStdout)
	u, observation, failure := c.Result()
	if failure != "" {
		t.Fatalf("failure class invented: %s", failure)
	}
	if u.Source != "zcode.reported-usage" || u.Coverage != "reported" {
		t.Fatalf("actual envelope not recognized: %+v", u)
	}
	if u.InputTokens == nil || *u.InputTokens != 4119566 || u.OutputTokens == nil || *u.OutputTokens != 18299 ||
		u.TotalTokens == nil || *u.TotalTokens != 4137865 || u.CacheReadTokens == nil || *u.CacheReadTokens != 3996160 ||
		u.CacheWriteTokens == nil || *u.CacheWriteTokens != 0 {
		t.Fatalf("counters not verbatim from the usage block: %+v", u)
	}
	// totalTokens already contains the cache split, so re-adding cache tokens
	// would report 8134025; the projection child total would report 8275730;
	// projection.contextUsed would swap input down to 122317. Exact equality
	// on all five counters rules each of those out.
	if *u.TotalTokens != 4119566+18299 {
		t.Fatalf("totalTokens is not the reported parent total: %+v", u)
	}
	if u.CostUSD != nil || u.CostBasis != "unavailable" || u.ReportedModel != nil || len(u.ReportedModels) != 0 {
		t.Fatalf("cost or model identity invented for an envelope that reports neither: %+v", u)
	}
	if observation.StdoutBytes == nil || *observation.StdoutBytes != int64(len(zcodeActualStdout)) {
		t.Fatalf("observation lost: %+v", observation)
	}
}

// The same envelope without a preamble exercises the consume path directly.
func TestZcodeEnvelopeWithoutWarningPreamble(t *testing.T) {
	c := NewCollector("zcode", true)
	zcodeFeed(c, zcodeActualEnvelope)
	u, _, _ := c.Result()
	if u.Source != "zcode.reported-usage" || u.InputTokens == nil || *u.InputTokens != 4119566 {
		t.Fatalf("preamble-free envelope not recognized: %+v", u)
	}
}

// Write boundaries must not matter; the CLI's pipe may split anywhere.
func TestZcodeEnvelopeAcrossSingleByteWrites(t *testing.T) {
	c := NewCollector("zcode", true)
	for _, b := range []byte(zcodeActualStdout) {
		_, _ = c.Writer("stdout").Write([]byte{b})
	}
	u, _, _ := c.Result()
	if u.InputTokens == nil || *u.InputTokens != 4119566 || *u.TotalTokens != 4137865 {
		t.Fatalf("chunked envelope lost: %+v", u)
	}
}

// Preamble handling is robust: several warning lines, including one carrying a
// brace mid-line, still precede a recognized envelope.
func TestZcodeEnvelopeAfterRepeatedWarningLines(t *testing.T) {
	c := NewCollector("zcode", true)
	zcodeFeed(c, "warn one\nwarn two {not json}\nwarn three\n"+zcodeActualEnvelope)
	u, _, _ := c.Result()
	if u.Coverage != "reported" || u.InputTokens == nil || *u.InputTokens != 4119566 {
		t.Fatalf("repeated preamble defeated recognition: %+v", u)
	}
}

// projection.* is live session state and response is model echo; neither may
// become billed usage, and no price may be inferred from either.
func TestZcodeProjectionAndEchoTextNeverBecomeUsage(t *testing.T) {
	c := NewCollector("zcode", true)
	zcodeFeed(c, `{
  "sessionId": "sess_probe", "traceId": "trace_probe", "turnId": "turn_probe",
  "response": "Echo text mentioning 5 tokens and $9.99 is not telemetry.",
  "usage": {"source": "provider", "modelRequestCount": 3, "inputTokens": 1000, "outputTokens": 50, "totalTokens": 1050, "cacheReadTokens": 900, "cacheWriteTokens": 10},
  "eventCount": 7,
  "projection": {"status": "idle", "turnCount": 4, "totalTokenCount": 1050, "contextUsed": 987654, "contextWindow": 2000000}
}
`)
	u, _, _ := c.Result()
	if u.Coverage != "reported" || u.InputTokens == nil || *u.InputTokens != 1000 || *u.OutputTokens != 50 ||
		*u.TotalTokens != 1050 || *u.CacheReadTokens != 900 || *u.CacheWriteTokens != 10 {
		t.Fatalf("billed usage drifted from the usage block: %+v", u)
	}
	if *u.TotalTokens != 1000+50 {
		t.Fatalf("cache tokens double-counted into the total: %+v", u)
	}
	if u.CostUSD != nil {
		t.Fatalf("echo or projection text priced: %+v", u)
	}
	for _, field := range []*int64{u.InputTokens, u.OutputTokens, u.CacheReadTokens, u.CacheWriteTokens, u.TotalTokens} {
		if field != nil && (*field == 987654 || *field == 2000000 || *field == 7 || *field == 4 || *field == 5 || *field == 9) {
			t.Fatalf("projection or echo number leaked into usage: %+v", u)
		}
	}
}

// Recognition binds to the envelope signature; near-misses stay unavailable.
func TestZcodeEnvelopeShapeIsRequired(t *testing.T) {
	cases := map[string]string{
		"missing turnId":      `{"sessionId": "sess_x", "usage": {"inputTokens": 5, "outputTokens": 1}}`,
		"missing sessionId":   `{"turnId": "turn_x", "usage": {"inputTokens": 5, "outputTokens": 1}}`,
		"no usage object":     `{"sessionId": "sess_x", "turnId": "turn_x", "eventCount": 3}`,
		"string counters":     `{"sessionId": "sess_x", "turnId": "turn_x", "usage": {"inputTokens": "4119566", "outputTokens": "18299", "totalTokens": "4137865"}}`,
		"snake_case counters": `{"sessionId": "sess_x", "turnId": "turn_x", "usage": {"input_tokens": 5, "output_tokens": 1, "total_tokens": 6}}`,
		"negative counters":   `{"sessionId": "sess_x", "turnId": "turn_x", "usage": {"inputTokens": -5, "outputTokens": -1, "totalTokens": -6}}`,
	}
	for name, envelope := range cases {
		t.Run(name, func(t *testing.T) {
			c := NewCollector("zcode", true)
			zcodeFeed(c, envelope+"\n")
			u, _, _ := c.Result()
			if u.Coverage != "unavailable" || u.InputTokens != nil || u.OutputTokens != nil || u.TotalTokens != nil {
				t.Fatalf("loose shape accepted as usage: %+v", u)
			}
		})
	}
}

// The preamble-tolerant retry must accept only the zcode envelope: a generic
// typed event pretty-printed behind a preamble stays unparsed, exactly as
// before this change.
func TestZcodeWarningPreambleDoesNotUnlockGenericEvents(t *testing.T) {
	c := NewCollector("zcode", true)
	zcodeFeed(c, "AI SDK Warning System: preamble\n{\n  \"type\": \"result\",\n  \"usage\": {\"input_tokens\": 77}\n}\n")
	u, _, _ := c.Result()
	if u.Coverage != "unavailable" || u.InputTokens != nil {
		t.Fatalf("preamble retry widened generic parsing: %+v", u)
	}
}

// Cross-adapter negative control: the envelope shape binds to the zcode
// adapter. Fed to any other adapter — bare or behind a warning preamble —
// the envelope must not become that adapter's reported usage. Other adapters
// keep their own formats and the generic typed-event path only.
func TestZcodeEnvelopeNeverBecomesOtherAdapterUsage(t *testing.T) {
	for _, adapter := range []string{"claude", "codex", "opencode", "kimi", "unknown"} {
		for name, text := range map[string]string{
			"bare":     zcodeActualEnvelope,
			"preamble": zcodeActualStdout,
		} {
			t.Run(adapter+"/"+name, func(t *testing.T) {
				c := NewCollector(adapter, true)
				zcodeFeed(c, text)
				u, _, _ := c.Result()
				if u.Source != "unavailable" || u.Coverage != "unavailable" || u.InputTokens != nil || u.OutputTokens != nil ||
					u.TotalTokens != nil || u.CostUSD != nil {
					t.Fatalf("zcode envelope became %s usage: %+v", adapter, u)
				}
			})
		}
	}
}
