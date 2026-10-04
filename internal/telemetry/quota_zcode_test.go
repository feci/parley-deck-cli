package telemetry

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"parley-deck-cli/internal/quota"
)

func recordedZcode(t *testing.T) QuotaInput {
	t.Helper()
	raw, err := os.ReadFile("testdata/quota/zcode-sdk-source-derived.stderr")
	if err != nil {
		t.Fatal(err)
	}
	in := quotaInput()
	in.Adapter = "zcode"
	in.Stderr = string(raw)
	in.ObservedAt = time.Date(2026, 10, 2, 21, 6, 7, 3000000, time.UTC)
	return in
}

func TestQuotaIncidentOneReadinessWithRecordedProviderReplay(t *testing.T) {
	raw, err := os.ReadFile("testdata/quota/zcode-incident-1-preflight.json")
	if err != nil {
		t.Fatal(err)
	}
	var f struct {
		ProviderFixture string `json:"provider_fixture"`
		Observed        string `json:"provider_observed_at"`
		Facilitator     string `json:"facilitator"`
		Readiness       []struct {
			ID        string `json:"id"`
			Available bool   `json:"available"`
		} `json:"readiness"`
	}
	if err = json.Unmarshal(raw, &f); err != nil {
		t.Fatal(err)
	}
	in := recordedZcode(t)
	body, err := os.ReadFile(filepath.Join("testdata/quota", "zcode-sdk-source-derived.stderr"))
	if err != nil {
		t.Fatal(err)
	}
	in.Stderr = string(body)
	in.ObservedAt, err = time.Parse(time.RFC3339Nano, f.Observed)
	if err != nil {
		t.Fatal(err)
	}
	ev := ClassifyQuota(in)
	if !ev.Eligible {
		t.Fatal(ev)
	}
	ids := []string{}
	members := []quota.Member{}
	for _, r := range f.Readiness {
		ids = append(ids, r.ID)
		m := quota.Member{ID: r.ID, Usable: r.Available}
		if r.ID == "zcode-1" {
			m.Evidence = &ev
		}
		members = append(members, m)
	}
	d := quota.Evaluate(quota.NewPolicy(nil, nil), ids, members, quota.Roles{Facilitator: f.Facilitator})
	if d.Applied || d.UsableSurvivors != 1 || len(d.Candidates) != 1 || d.Block == "" {
		t.Fatal(d)
	}
}

func TestQuotaZcodeEveryRecordAndObservedClock(t *testing.T) {
	in := recordedZcode(t)
	record := in.Stderr[:strings.LastIndex(in.Stderr, "Error: Turn execution failed")]
	in.Stderr = record + record + "Error: Turn execution failed\n"
	if e := ClassifyQuota(in); !e.Eligible {
		t.Fatal(e)
	}
	for _, extra := range []string{"AI_APICallError: second incomplete record\n", "TypeError: other failure\n", "```assistant\n"} {
		bad := in
		bad.Stderr = record + extra + "Error: Turn execution failed\n"
		if e := ClassifyQuota(bad); e.Eligible {
			t.Fatal(extra, e)
		}
	}
	bad := in
	bad.Stderr = strings.Replace(in.Stderr, `"reset_at":"2026-10-04T22:14:57.003Z"`, `"reset_at":"2026-10-04T22:14:57.003Z","reset_at":"2026-10-05T22:14:57.003Z"`, 1)
	if e := ClassifyQuota(bad); e.Eligible {
		t.Fatal("duplicate reset accepted", e)
	}
	// Process teardown can be minutes later. Retry-after is relative to record
	// receipt; treating it as relative to teardown wrongly makes records conflict.
	c := NewCollector("zcode", false)
	c.quotaClock = func() time.Time { return in.ObservedAt }
	c.Writer("stderr").Write([]byte(in.Stderr))
	later := in
	later.ObservedAt = later.ObservedAt.Add(10 * time.Minute)
	if e := c.QuotaEvidence(later); !e.Eligible || !e.ObservedAt.Equal(in.ObservedAt) {
		t.Fatal(e)
	}
}
func TestQuotaZcodeRecordedBodySourceDerivedFraming(t *testing.T) {
	in := recordedZcode(t)
	e := ClassifyQuota(in)
	if !e.Eligible || e.InvocationID == "" || e.RuleID == "" || e.Provenance == "" || e.ResetAt == nil || e.ResetHint() != "2026-10-04T22:14:57.003Z" || !strings.Contains(e.RawReset, "2026-10-05 06:14:57") {
		t.Fatalf("%+v", e)
	}
	// Incident 1 preflight used the same provider failure shape. Its recorded
	// readiness JSON lacks a full dump, so do not invent a second native capture.
	c := NewCollector("zcode", false)
	c.quotaClock = func() time.Time { return in.ObservedAt }
	c.Writer("stderr").Write([]byte(in.Stderr))
	if got := c.QuotaEvidence(in); !got.Eligible {
		t.Fatal(got)
	}
}
func TestQuotaZcodeAdversarialFixtures(t *testing.T) {
	for _, name := range []string{"assistant-quoted", "tool-quoted", "mixed-error", "malformed", "short-reset", "display-only", "contradictory"} {
		t.Run(name, func(t *testing.T) {
			in := recordedZcode(t)
			raw, err := os.ReadFile("testdata/quota/zcode-" + name + ".stderr")
			if err != nil {
				t.Fatal(err)
			}
			in.Stderr = string(raw)
			if e := ClassifyQuota(in); e.Eligible {
				t.Fatal(e)
			}
		})
	}
	for _, name := range []string{"valid-artifact", "later-success", "zero-exit", "stdout-only", "no-terminal", "unfinished-record"} {
		t.Run(name, func(t *testing.T) {
			in := recordedZcode(t)
			switch name {
			case "valid-artifact":
				in.ValidArtifact = true
			case "later-success":
				in.LaterSuccess = true
			case "zero-exit":
				*in.ExitCode = 0
			case "stdout-only":
				in.Stdout = in.Stderr
				in.Stderr = "Error: Turn execution failed"
			case "no-terminal":
				in.Stderr = strings.ReplaceAll(in.Stderr, "Error: Turn execution failed", "")
			case "unfinished-record":
				in.Stderr = strings.ReplaceAll(in.Stderr, "Error: Turn execution failed", "statusCode: 429\nError: Turn execution failed")
			}
			if e := ClassifyQuota(in); e.Eligible {
				t.Fatal(e)
			}
		})
	}
}
func TestQuotaDisplayClockBoundsAndParentheses(t *testing.T) {
	in := quotaInput()
	cases := []struct {
		message, reset, retry string
		want                  bool
	}{
		{"weekly limit exhausted; reset at 2026-10-03T02:00:00Z (reset after 2h)", "2026-10-03T02:00:00Z", "7200", true},
		{"weekly limit exhausted; reset at 2026-10-03 16:00:00 (reset after 2h)", "2026-10-03T02:00:00Z", "7200", true},
		{"weekly limit exhausted; reset at 2026-10-02 14:00:00 (reset after 2h)", "2026-10-03T02:00:00Z", "7200", true},
		{"weekly limit exhausted; reset at 2026-10-03 16:15:00 (reset after 2h)", "2026-10-03T02:00:00Z", "7200", false},
		{"weekly limit exhausted; reset at 2026-10-03 02:01:00 (reset after 2h)", "2026-10-03T02:00:00Z", "7200", false},
		{"weekly limit exhausted; reset at 2026-10-03 02:00:00 (reset after 2h)", "", "", false},
		{"weekly limit exhausted; reset at 2026-10-03 02:00:00 (reset after 2h)", "2026-10-03T02:00:02Z", "7200", false},
	}
	for _, c := range cases {
		e := refineNativeQuota(in, nativeQuotaError{Status: 429, Message: c.message, ResetAt: c.reset, RetryAfter: c.retry, Provenance: "policy-test"})
		if e.Eligible != c.want {
			t.Errorf("%s: %+v", c.message, e)
		}
	}
}
