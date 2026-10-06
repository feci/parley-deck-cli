package telemetry

import (
	"os"
	"strings"
	"testing"
	"time"
)

func TestQuotaFixupSDKAllowlistAndAllReviewerWrongPositives(t *testing.T) {
	in := recordedZcode(t)
	termAt := strings.LastIndex(in.Stderr, "Error: Turn execution failed")
	record := in.Stderr[:termAt]
	term := in.Stderr[termAt:]
	for _, bad := range []string{
		"The earlier probe said:", "$ cat zcode-incident-2.stderr", "```assistant", "{\"role\":\"tool\",\"content\":\"quoted\"}", "API Error: 429 Weekly Limit Exhausted",
		"Error: ENOENT open x", "TypeError: x is not a function", "RangeError: Maximum call stack size exceeded", "SyntaxError: Unexpected token < in JSON", "AbortError: This operation was aborted", "FetchError: request to https://x failed, reason: socket hang up", "Unhandled promise rejection: ZodError: invalid_type", "AI_RetryError: Failed after 3 attempts. Last error: Internal Server Error", "statusCode: 500", "AI_APICallError: dangling", "unknown ordinary text",
	} {
		t.Run(bad, func(t *testing.T) {
			for _, raw := range []string{bad + "\n" + in.Stderr, record + bad + "\n" + term, in.Stderr + bad + "\n"} {
				x := in
				x.Stderr = raw
				if e := ClassifyQuota(x); e.Eligible {
					t.Fatalf("wrong-positive %q: %+v", bad, e)
				}
			}
		})
	}
	cases := map[string]string{
		"missing-wrapper":         in.Stderr[strings.Index(in.Stderr, "\n")+1:],
		"missing-close":           strings.Replace(in.Stderr, "\n}\n", "\n", 1),
		"missing-symbol":          strings.Replace(in.Stderr, "  Symbol(vercel.ai.error.AI_APICallError): true\n", "", 1),
		"other-cause":             strings.Replace(in.Stderr, "cause: undefined", "cause: 'FetchError: socket hang up'", 1),
		"other-data":              strings.Replace(in.Stderr, "data: undefined", "data: { error: { message: 'Internal Server Error' } }", 1),
		"duplicate-field":         strings.Replace(in.Stderr, "statusCode: 429,", "statusCode: 429, statusCode: 429,", 1),
		"header-status-mismatch":  strings.Replace(in.Stderr, "statusCode: 429", "statusCode: 503", 1),
		"header-message-mismatch": strings.Replace(in.Stderr, "Weekly/Monthly Limit Exhausted", "Internal Server Error", 1),
		"unknown-field":           strings.Replace(in.Stderr, "data: undefined", "invented: true, data: undefined", 1),
		"truncated-request":       strings.Replace(in.Stderr, "messages: []", "messages: [Object]", 1),
		"duplicate-header":        strings.Replace(in.Stderr, "'retry-after': '176930'", "'retry-after': '176930', 'Retry-After': '3600'", 1),
		"contradictory-header":    strings.Replace(in.Stderr, "'retry-after': '176930'", "'retry-after': '3600'", 1),
		"double-terminal":         in.Stderr + term,
	}
	for name, raw := range cases {
		t.Run(name, func(t *testing.T) {
			x := in
			x.Stderr = raw
			if e := ClassifyQuota(x); e.Eligible {
				t.Fatal(e)
			}
		})
	}
	// The historical excerpt's body/reset survives in the complete source-derived
	// fixture. The four-line excerpt alone cannot demonstrate complete framing.
	excerpt, e := os.ReadFile("testdata/quota/zcode-incident-2.stderr")
	if e != nil {
		t.Fatal(e)
	}
	x := in
	x.Stderr = string(excerpt)
	if ClassifyQuota(x).Eligible {
		t.Fatal("incomplete native excerpt accepted")
	}
	if e := ClassifyQuota(in); !e.Eligible || e.ResetHint() != "2026-10-04T22:14:57.003Z" {
		t.Fatal("recorded body/reset lost", e)
	}
}
func TestQuotaFixupRetryFramingAndNestedConsistency(t *testing.T) {
	raw, e := os.ReadFile("testdata/quota/zcode-sdk-retry-source-derived.stderr")
	if e != nil {
		t.Fatal(e)
	}
	in := recordedZcode(t)
	in.Stderr = string(raw)
	if e := ClassifyQuota(in); !e.Eligible {
		t.Fatal("source-derived SDK retry rejected", e)
	}
	for name, bad := range map[string]string{
		"wrong-count":         strings.Replace(in.Stderr, "after 3 attempts", "after 4 attempts", 1),
		"wrong-reason":        strings.Replace(in.Stderr, "maxRetriesExceeded", "errorNotRetryable", 1),
		"wrong-last-message":  strings.Replace(in.Stderr, "Last error: [glm/glm-5.3] [429]: Weekly/Monthly Limit Exhausted", "Last error: Internal Server Error", 1),
		"nested-mixed-status": strings.Replace(in.Stderr, "statusCode: 429", "statusCode: 503", 1),
		"missing-last":        in.Stderr[:strings.Index(in.Stderr, "lastError:")] + "}\nError: Turn execution failed\n",
		"incomplete-inspect":  strings.Replace(in.Stderr, "requestBodyValues: { model: 'fixture-model', messages: [] }", "requestBodyValues: [Object]", 1),
	} {
		t.Run(name, func(t *testing.T) {
			x := in
			x.Stderr = bad
			if e := ClassifyQuota(x); e.Eligible {
				t.Fatal(e)
			}
		})
	}
}
func TestQuotaFixupPerRecordReceiptTimes(t *testing.T) {
	in := recordedZcode(t)
	end := strings.LastIndex(in.Stderr, "Error: Turn execution failed")
	record := in.Stderr[:end]
	term := in.Stderr[end:]
	second := strings.ReplaceAll(strings.ReplaceAll(record, "176930", "176890"), "49h 8m 50s", "49h 8m 10s")
	c := NewCollector("zcode", false)
	c.quotaClock = func() time.Time { return in.ObservedAt }
	c.Writer("stderr").Write([]byte(record))
	c.quotaClock = func() time.Time { return in.ObservedAt.Add(40 * time.Second) }
	c.Writer("stderr").Write([]byte(second + term))
	delayed := in
	delayed.ObservedAt = in.ObservedAt.Add(10 * time.Minute)
	if e := c.QuotaEvidence(delayed); !e.Eligible {
		t.Fatal(e)
	}
	delayed.Stderr = record + second + term
	if e := ClassifyQuota(delayed); !e.Eligible || !e.ObservedAt.Equal(delayed.ObservedAt) {
		t.Fatal("complete aggregate should use its terminal receipt", e)
	}
}

func completeSyntheticSDKBody(message, body string) string {
	return "APICallError [AI_APICallError]: " + message + "\n    at fixture (offline.js:1:1) {\n cause: undefined, url: 'https://provider.invalid', requestBodyValues: {}, statusCode: 429, responseHeaders: {}, responseBody: '" + body + "', isRetryable: true, data: undefined, Symbol(vercel.ai.error): true, Symbol(vercel.ai.error.AI_APICallError): true\n}\nError: Turn execution failed\n"
}
func TestQuotaFixupCompleteFramingSemanticBoundaries(t *testing.T) {
	for _, tc := range []struct {
		message, fields string
		want            bool
	}{
		{"Weekly limit exhausted", `"retry_after":3600`, true},
		{"Weekly limit exhausted", `"retry_after":3599`, false},
		{"5-hour usage limit reached", `"retry_after":7200`, true},
		{"Hourly limit exhausted", `"retry_after":7200`, true},
		{"Too Many Requests", `"retry_after":7200`, false},
		{"quota exceeded", `"retry_after":7200`, false},
		{"credit balance", `"retry_after":7200`, false},
		{"Your weekly limit has not been reached", `"retry_after":7200`, false},
		{"Rate limited; weekly limit reached soon? no: the weekly limit reached warning threshold is 80%", `"retry_after":7200`, false},
		{"Weekly limit exhausted", `"reset_at":null,"retry_after":7200`, false},
		{"Weekly limit exhausted", `"reset_at":"2026-10-05 00:00:00"`, false},
		{"Weekly limit exhausted", `"reset_at":"2026-10-01T00:00:00Z"`, false},
	} {
		t.Run(tc.message+tc.fields, func(t *testing.T) {
			in := recordedZcode(t)
			body := `{"error":{"message":"` + tc.message + `",` + tc.fields + `}}`
			in.Stderr = completeSyntheticSDKBody(tc.message, body)
			got := ClassifyQuota(in)
			if got.Eligible != tc.want {
				t.Fatal(got)
			}
		})
	}
	in := recordedZcode(t)
	for _, raw := range []string{strings.Replace(in.Stderr, "cause: undefined", "cause: 'undefined'", 1), strings.Replace(in.Stderr, "    at offlineFixture (sdk-framing-harness.cjs:1:1) {", "    at arbitrary explanatory prose {", 1), strings.Replace(in.Stderr, "cause: undefined,", "", 1)} {
		bad := in
		bad.Stderr = raw
		if e := ClassifyQuota(bad); e.Eligible {
			t.Fatal("invalid framing accepted", e)
		}
	}
	raw, e := os.ReadFile("testdata/quota/zcode-sdk-retry-source-derived.stderr")
	if e != nil {
		t.Fatal(e)
	}
	at := strings.Index(string(raw), "lastError:")
	changed := string(raw[:at]) + strings.Replace(string(raw[at:]), "fixture-model", "different-request", 1)
	in.Stderr = changed
	if e := ClassifyQuota(in); e.Eligible {
		t.Fatal("lastError differed from final retry", e)
	}
}
