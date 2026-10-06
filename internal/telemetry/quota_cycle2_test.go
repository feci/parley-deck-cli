package telemetry

import (
	"os"
	"strings"
	"testing"
	"time"
)

func TestQuotaCycle2SourceDerivedRetriesAndSingleReceipt(t *testing.T) {
	for _, name := range []string{"retained-retries", "retained-top-level", "backoff-retries"} {
		t.Run(name, func(t *testing.T) {
			in := recordedZcode(t)
			raw, err := os.ReadFile("testdata/quota/zcode-" + name + "-source-derived.stderr")
			if err != nil {
				t.Fatal(err)
			}
			in.Stderr = string(raw)
			if strings.HasPrefix(name, "retained") {
				in.ObservedAt = in.ObservedAt.Add(40 * time.Second)
			}
			in.StartedAt = in.ObservedAt.Add(-time.Minute)
			got := ClassifyQuota(in)
			if !got.Eligible || !got.ObservedAt.Equal(in.ObservedAt) || got.ResetHint() != "2026-10-04T22:14:57.003Z" {
				t.Fatal(got)
			}
			c := NewCollector("zcode", false)
			c.quotaClock = func() time.Time { return in.ObservedAt }
			c.Writer("stderr").Write(raw)
			delayed := in
			delayed.ObservedAt = delayed.ObservedAt.Add(time.Hour)
			if e := c.QuotaEvidence(delayed); !e.Eligible || !e.ObservedAt.Equal(in.ObservedAt) {
				t.Fatal(e)
			}
			for _, kind := range []string{"future", "before-start", "class", "header"} {
				bad := in
				switch kind {
				case "future":
					bad.ObservedAt = bad.ObservedAt.Add(-2 * time.Second)
				case "before-start":
					bad.StartedAt = bad.ObservedAt
				case "class":
					bad.Stderr = strings.ReplaceAll(bad.Stderr, "Weekly/Monthly Limit Exhausted", "Daily Limit Exhausted")
					at := strings.Index(bad.Stderr, "\n}")
					if at > 0 {
						bad.Stderr = in.Stderr[:at] + bad.Stderr[at:]
					} else {
						continue
					}
				case "header":
					bad.Stderr = strings.Replace(bad.Stderr, "'retry-after': '", "'retry-after': '1", 1)
				}
				if e := ClassifyQuota(bad); e.Eligible {
					t.Fatalf("%s: %+v", kind, e)
				}
			}
		})
	}
}

func TestQuotaCycle2TerminalThresholdAndObservationOrder(t *testing.T) {
	in := recordedZcode(t)
	reset := in.ObservedAt.Add(time.Hour)
	template := in.Stderr
	makeRecord := func(seconds string) string {
		// Complete framing from the existing SDK fixture, with all duration sources
		// changed together; observations are derived from one absolute reset.
		raw := strings.ReplaceAll(template, "2026-10-04T22:14:57.003Z", reset.Format(time.RFC3339Nano))
		raw = strings.ReplaceAll(raw, "2026-10-05 06:14:57", reset.In(time.FixedZone("UTC+8", 8*3600)).Format("2006-01-02 15:04:05"))
		raw = strings.ReplaceAll(raw, "176930", seconds)
		raw = strings.ReplaceAll(raw, "49h 8m 50s", seconds+"s")
		return raw[:strings.LastIndex(raw, "Error: Turn execution failed")]
	}
	in.Stderr = makeRecord("3660") + makeRecord("3540") + "Error: Turn execution failed\n"
	in.ObservedAt = reset.Add(-59 * time.Minute)
	if got := ClassifyQuota(in); got.Eligible || !strings.Contains(got.Reason, "below 60") {
		t.Fatal("61-minute first / 59-minute terminal", got)
	}
	in.ObservedAt = reset.Add(-time.Hour)
	in.Stderr = makeRecord("3660") + makeRecord("3600") + "Error: Turn execution failed\n"
	if got := ClassifyQuota(in); !got.Eligible {
		t.Fatal("exact 60 minutes", got)
	}
	in.Stderr = makeRecord("3600") + makeRecord("3660") + "Error: Turn execution failed\n"
	if got := ClassifyQuota(in); got.Eligible || !strings.Contains(got.Reason, "observation") {
		t.Fatal("backwards observations", got)
	}
}

func TestQuotaCycle2DefaultConsoleDepthStaysRejected(t *testing.T) {
	in := recordedZcode(t)
	raw, err := os.ReadFile("testdata/quota/zcode-default-console-source-derived.stderr")
	if err != nil {
		t.Fatal(err)
	}
	in.Stderr = string(raw)
	if !strings.Contains(in.Stderr, "[Object]") {
		t.Fatal("offline default sink no longer reproduces truncation")
	}
	got := ClassifyQuota(in)
	if got.Eligible || !strings.Contains(got.Reason, "framing:") || strings.Contains(got.Reason, "unsupported") {
		t.Fatal(got)
	}
}
