// claude-1 round-04 probe (read-only use of exported telemetry API; no provider launched).
// Exercises the historical changing-countdown input: two separately emitted native records
// (176930 / 49h 8m 50s and 176890 / 49h 8m 10s, same absolute reset) and an aggregated
// AI SDK RetryError wrapper whose nested messages differ only in the countdown.
package main

import (
	"fmt"
	"os"
	"strings"
	"time"

	"parley-deck-cli/internal/telemetry"
)

const resetAt = "2026-10-04T22:14:57.003Z"

func msg(countdown string) string {
	return "[glm/glm-5.3] [429]: Weekly/Monthly Limit Exhausted. Your limit will reset at 2026-10-05 06:14:57 (reset after " + countdown + ")"
}

func api(indent, countdown, retry, reset, display string) string {
	m := "[glm/glm-5.3] [429]: Weekly/Monthly Limit Exhausted. Your limit will reset at " + display + " (reset after " + countdown + ")"
	body := `{"error":{"message":"` + m + `","retry_after":` + retry + `,"reset_at":"` + reset + `"}}`
	lines := []string{
		"APICallError [AI_APICallError]: " + m,
		"    at offlineFixture (sdk-framing-harness.cjs:1:1) {",
		"  cause: undefined,",
		"  url: 'https://provider.invalid/v1/messages',",
		"  requestBodyValues: { model: 'fixture-model', messages: [] },",
		"  statusCode: 429,",
		"  responseHeaders: { 'content-type': 'application/json', 'retry-after': '" + retry + "' },",
		"  responseBody: '" + body + "',",
		"  isRetryable: true,",
		"  data: undefined,",
		"  Symbol(vercel.ai.error): true,",
		"  Symbol(vercel.ai.error.AI_APICallError): true",
		"}",
	}
	for i := range lines {
		lines[i] = indent + lines[i]
	}
	return strings.Join(lines, "\n")
}

func retryWrapper(nested []string, lastCountdown, lastRetry string) string {
	var b strings.Builder
	b.WriteString("RetryError [AI_RetryError]: Failed after " + fmt.Sprint(len(nested)) + " attempts. Last error: " + msg(lastCountdown) + "\n")
	b.WriteString("    at offlineFixture (sdk-framing-harness.cjs:1:1) {\n  cause: undefined,\n  reason: 'maxRetriesExceeded',\n  errors: [\n")
	for i, n := range nested {
		b.WriteString(n)
		if i < len(nested)-1 {
			b.WriteString(",")
		}
		b.WriteString("\n")
	}
	last := api("  ", lastCountdown, lastRetry, resetAt, "2026-10-05 06:14:57")
	last = strings.TrimPrefix(last, "  ")
	b.WriteString("  ],\n  lastError: " + last + ",\n  Symbol(vercel.ai.error): true,\n  Symbol(vercel.ai.error.AI_RetryError): true\n}\n")
	return b.String()
}

const term = "Error: Turn execution failed (traceId: claude1-r4-probe)\n"

func classify(label, stderr string, observed time.Time) {
	code := 1
	e := telemetry.ClassifyQuota(telemetry.QuotaInput{Adapter: "zcode", InvocationID: "claude1-r4-" + label, Stderr: stderr, ObservedAt: observed, ExitCode: &code})
	fmt.Printf("%-58s eligible=%-5v rule=%q reason=%q reset=%s raw=%q\n", label, e.Eligible, e.RuleID, e.Reason, e.ResetHint(), e.RawReset)
}

func main() {
	fixture, err := os.ReadFile("internal/telemetry/testdata/quota/zcode-sdk-retry-source-derived.stderr")
	if err != nil {
		panic(err)
	}
	obs50, _ := time.Parse(time.RFC3339Nano, "2026-10-02T21:06:07.003Z") // reset_at - 176930 s
	obs10, _ := time.Parse(time.RFC3339Nano, "2026-10-02T21:06:47.003Z") // reset_at - 176890 s
	// A. Producer source-derived retry fixture (three identical nested errors).
	classify("A source-derived retry fixture (identical retries)", string(fixture), obs50)
	// B. Aggregated RetryError, countdown changes per attempt, same absolute reset, lastError == final.
	n1 := api("    ", "49h 8m 50s", "176930", resetAt, "2026-10-05 06:14:57")
	n2 := api("    ", "49h 8m 30s", "176910", resetAt, "2026-10-05 06:14:57")
	n3 := api("    ", "49h 8m 10s", "176890", resetAt, "2026-10-05 06:14:57")
	agg := retryWrapper([]string{n1, n2, n3}, "49h 8m 10s", "176890")
	os.WriteFile(".parley-runtime/claude1-r4/zretry/B-aggregated-changing-countdown.stderr", []byte(agg+term), 0600)
	classify("B aggregated retry, changing countdown, obs=final", agg+term, obs10)
	classify("B aggregated retry, changing countdown, obs=first", agg+term, obs50)
	// B2. Two attempts only (the two retained native values).
	agg2 := retryWrapper([]string{n1, n3}, "49h 8m 10s", "176890")
	classify("B2 aggregated retry of the two retained countdowns", agg2+term, obs10)
	// B3. Control: identical nested messages built by this probe (should match fixture behaviour).
	agg3 := retryWrapper([]string{n3, n3, n3}, "49h 8m 10s", "176890")
	classify("B3 control: probe-built identical retries", agg3+term, obs10)
	// C. Two separately emitted top-level records, single observation clock (one write).
	top1 := api("", "49h 8m 50s", "176930", resetAt, "2026-10-05 06:14:57")
	top2 := api("", "49h 8m 10s", "176890", resetAt, "2026-10-05 06:14:57")
	classify("C two top-level records, one clock (obs=first)", top1+"\n"+top2+"\n"+term, obs50)
	classify("C two top-level records, one clock (obs=final)", top1+"\n"+top2+"\n"+term, obs10)
	classify("C1 first record alone (obs=first)", top1+"\n"+term, obs50)
	classify("C2 second record alone (obs=final)", top2+"\n"+term, obs10)
	// D. Two separately emitted records through the real Collector clock (2 s apart, real time).
	start := time.Now().UTC().Truncate(time.Second)
	reset := start.Add(2 * time.Hour)
	disp := reset.In(time.FixedZone("UTC+8", 8*3600)).Format("2006-01-02 15:04:05")
	ra := reset.Format(time.RFC3339)
	d1 := api("", "2h", "7200", ra, disp)
	c := telemetry.NewCollector("zcode", false)
	for time.Now().UTC().Before(start) {
		time.Sleep(10 * time.Millisecond)
	}
	time.Sleep(time.Until(start.Add(200 * time.Millisecond)))
	c.Writer("stderr").Write([]byte(d1 + "\n"))
	time.Sleep(2 * time.Second)
	d2 := api("", "1h 59m 58s", "7198", ra, disp)
	c.Writer("stderr").Write([]byte(d2 + "\n" + term))
	code := 1
	e := c.QuotaEvidence(telemetry.QuotaInput{InvocationID: "claude1-r4-D", ObservedAt: time.Now().UTC().Add(5 * time.Minute), ExitCode: &code})
	fmt.Printf("%-58s eligible=%-5v rule=%q reason=%q reset=%s observed=%s\n", "D two records, real collector clock 2 s apart", e.Eligible, e.RuleID, e.Reason, e.ResetHint(), e.ObservedAt.Format(time.RFC3339Nano))
	// E. Aggregated wrapper with changing countdown through the collector in one write (real time).
	start2 := time.Now().UTC().Truncate(time.Second).Add(time.Second)
	reset2 := start2.Add(2 * time.Hour)
	disp2 := reset2.In(time.FixedZone("UTC+8", 8*3600)).Format("2006-01-02 15:04:05")
	ra2 := reset2.Format(time.RFC3339)
	e1 := api("    ", "2h 0m 6s", "7206", ra2, disp2)
	e2 := api("    ", "2h 0m 4s", "7204", ra2, disp2)
	e3 := api("    ", "2h", "7200", ra2, disp2)
	w := strings.Replace(retryWrapper([]string{e1, e2, e3}, "2h", "7200"), resetAt, ra2, -1)
	w = strings.Replace(w, "2026-10-05 06:14:57", disp2, -1)
	c2 := telemetry.NewCollector("zcode", false)
	time.Sleep(time.Until(start2.Add(200 * time.Millisecond)))
	c2.Writer("stderr").Write([]byte(w + term))
	e = c2.QuotaEvidence(telemetry.QuotaInput{InvocationID: "claude1-r4-E", ObservedAt: time.Now().UTC(), ExitCode: &code})
	fmt.Printf("%-58s eligible=%-5v rule=%q reason=%q reset=%s\n", "E aggregated, backoff-like countdowns 6s/4s/0s, one write", e.Eligible, e.RuleID, e.Reason, e.ResetHint())
	// E-control: same wrapper with identical nested errors.
	wc := strings.Replace(retryWrapper([]string{e3, e3, e3}, "2h", "7200"), resetAt, ra2, -1)
	wc = strings.Replace(wc, "2026-10-05 06:14:57", disp2, -1)
	classify("E-control identical nested, same clock", wc+term, start2)
}
