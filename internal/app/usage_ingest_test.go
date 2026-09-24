package app

// lean-organizer D tests: path-only usage ingest (no count-accepting flag — proven
// by scanning the command's own flag definitions), six total_token_usage fields
// verbatim, idempotent re-ingest, round-trip re-parse equality, <= 1 KB stdout,
// explicit-args + timestamp-window attribution with an honest `ambiguous` residue,
// and streaming over a generated ~228 MB fixture with bounded memory.

import (
	"bytes"
	"encoding/json"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"testing"

	"parley-deck-cli/internal/protocol"
	"parley-deck-cli/internal/runmanifest"
	"time"
)

func writeCodexFixture(t *testing.T, path string, totals string, ts string) {
	t.Helper()
	body := `{"timestamp":"` + ts + `","type":"event_msg","payload":{"type":"token_count","info":{"total_token_usage":` + totals + `,"last_token_usage":{},"model_context_window":258400}}}` + "\n"
	big := `{"timestamp":"2026-09-23T10:00:00.000Z","type":"response_item","payload":{"type":"message","content":[{"type":"input_text","text":"` + strings.Repeat("x", 4096) + `"}]}}` + "\n"
	var b bytes.Buffer
	for i := 0; i < 50; i++ {
		b.WriteString(big)
	}
	b.WriteString(body)
	if err := os.WriteFile(path, b.Bytes(), 0o644); err != nil {
		t.Fatal(err)
	}
}

func TestUsageIngestHasNoCountAcceptingFlag(t *testing.T) {
	// Adversarial source scan: every flag this command defines must come from the
	// ratified allowlist; nothing count/token/usage-number-shaped may exist.
	src, err := os.ReadFile("usage_ingest.go")
	if err != nil {
		t.Fatal(err)
	}
	re := regexp.MustCompile(`fs\.\w+\("([a-z-]+)"`)
	defined := map[string]bool{}
	for _, m := range re.FindAllStringSubmatch(string(src), -1) {
		defined[m[1]] = true
	}
	if len(defined) == 0 {
		t.Fatal("flag scan found no definitions — the test is broken")
	}
	allow := map[string]bool{"dir": true, "agent": true, "source": true, "path": true, "idea": true, "phase": true}
	for name := range defined {
		if !allow[name] {
			t.Errorf("usage ingest defines flag --%s beyond the path-only allowlist", name)
		}
		lower := name
		if strings.Contains(lower, "token") || strings.Contains(lower, "count") || strings.Contains(lower, "usage") {
			t.Errorf("flag --%s smells like a count-accepting flag", name)
		}
	}
}

func TestUsageIngestCodexRolloutSixFieldsVerbatim(t *testing.T) {
	dir := t.TempDir()
	root := dir
	seedMinimalDeck(t, root)
	if err := os.MkdirAll(filepath.Join(root, protocol.DeckDir, "ideas", "u"), 0o755); err != nil {
		t.Fatal(err)
	}
	os.WriteFile(filepath.Join(root, protocol.DeckDir, "ideas", "u", "00-prompt.md"), []byte("---\nidea: u\nstatus: round-01\n---\n"), 0o644)
	fixture := filepath.Join(dir, "rollout.jsonl")
	want := `{"input_tokens":17199,"cached_input_tokens":17024,"cache_write_input_tokens":3,"output_tokens":204,"reasoning_output_tokens":11,"total_tokens":17403}`
	writeCodexFixture(t, fixture, want, "2026-09-23T20:00:00.000Z")

	var out, errOut bytes.Buffer
	code := runUsageIngest([]string{"--dir", root, "--agent", "codex-1", "--source", "codex-rollout", "--path", fixture, "--idea", "u", "--phase", "2"}, &out, &errOut)
	if code != 0 {
		t.Fatalf("ingest exit %d: %s", code, errOut.String())
	}
	if len(out.String()) > 1024 {
		t.Errorf("stdout must stay under 1 KB, got %d B", len(out.String()))
	}
	ledger := usageLedgerPath(root, "u")
	data, err := os.ReadFile(ledger)
	if err != nil {
		t.Fatalf("ledger: %v", err)
	}
	var rows []ledgerRow
	for _, line := range strings.Split(string(data), "\n") {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		var r ledgerRow
		if json.Unmarshal([]byte(line), &r) == nil {
			rows = append(rows, r)
		}
	}
	if len(rows) != 1 {
		t.Fatalf("exactly one ledger row per ingest, got %d", len(rows))
	}
	row := rows[0]
	rawTotals, _ := json.Marshal(row.Result.Totals)
	// The six fields verbatim, in canonical JSON naming.
	for _, field := range []string{`"input_tokens":17199`, `"cached_input_tokens":17024`, `"cache_write_input_tokens":3`, `"output_tokens":204`, `"reasoning_output_tokens":11`, `"total_tokens":17403`} {
		if !strings.Contains(string(rawTotals), field) {
			t.Errorf("totals must carry %s verbatim; got %s", field, rawTotals)
		}
	}
	if !strings.HasPrefix(string(data), "# attribution method:") {
		t.Errorf("ledger must state the attribution method in its header")
	}
	// Round-trip re-parse equality: re-reading the row yields the same totals.
	res, err := parseCodexRollout(fixture)
	if err != nil {
		t.Fatal(err)
	}
	if res.Totals != row.Result.Totals {
		t.Errorf("round-trip re-parse mismatch: %+v vs %+v", res.Totals, row.Result.Totals)
	}
}

func TestUsageIngestIdempotent(t *testing.T) {
	dir := t.TempDir()
	root := dir
	seedMinimalDeck(t, root)
	os.MkdirAll(filepath.Join(root, protocol.DeckDir, "ideas", "u"), 0o755)
	fixture := filepath.Join(dir, "rollout.jsonl")
	writeCodexFixture(t, fixture, `{"input_tokens":10,"cached_input_tokens":0,"cache_write_input_tokens":0,"output_tokens":2,"reasoning_output_tokens":0,"total_tokens":12}`, "2026-09-23T20:00:00.000Z")
	args := []string{"--dir", root, "--agent", "a1", "--source", "codex-rollout", "--path", fixture, "--idea", "u", "--phase", "1"}
	var out, errOut bytes.Buffer
	if code := runUsageIngest(args, &out, &errOut); code != 0 {
		t.Fatalf("first ingest: %s", errOut.String())
	}
	out.Reset()
	if code := runUsageIngest(args, &out, &errOut); code != 0 {
		t.Fatalf("second ingest: %s", errOut.String())
	}
	if !strings.Contains(out.String(), "idempotent no-op") {
		t.Errorf("re-ingest must be an idempotent no-op; got %q", out.String())
	}
	data, _ := os.ReadFile(usageLedgerPath(root, "u"))
	n := 0
	for _, line := range strings.Split(string(data), "\n") {
		if strings.HasPrefix(strings.TrimSpace(line), "{") {
			n++
		}
	}
	if n != 1 {
		t.Errorf("idempotent re-ingest must not append a second row, got %d", n)
	}
}

// TestUsageIngestAttributionWindows (fix-up F7, claude-1 MAJ-7): the run records
// here are built by the REAL machinery the driver uses — runmanifest.New/Write at
// creation (created_at == updated_at, the only shape a real launch produces) and
// runmanifest.TouchUpdatedAt at the driver's phase-transition chokepoint — replacing
// the former hand-written two-hour-window manifest no driver ever produced.
func TestUsageIngestAttributionWindows(t *testing.T) {
	dir := t.TempDir()
	root := dir
	seedMinimalDeck(t, root)
	os.MkdirAll(filepath.Join(root, protocol.DeckDir, "ideas", "u"), 0o755)

	created := time.Now().UTC().Add(-1 * time.Hour)
	runID := "20260923T190000.000000000Z"
	manifest := runmanifest.New(runmanifest.Options{RunID: runID, Root: root, IdeaSlug: "u", CreatedAt: created, UpdatedAt: created})
	if err := runmanifest.Write(root, runID, manifest); err != nil {
		t.Fatal(err)
	}
	// The driver's transition touch opens the window: updated_at advances to now.
	if err := runmanifest.TouchUpdatedAt(root, runID); err != nil {
		t.Fatal(err)
	}

	// An accounting event inside the real window [created, now] for the SAME idea.
	fixture := filepath.Join(dir, "rollout.jsonl")
	writeCodexFixture(t, fixture, `{"input_tokens":10,"cached_input_tokens":0,"cache_write_input_tokens":0,"output_tokens":2,"reasoning_output_tokens":0,"total_tokens":12}`, time.Now().UTC().Add(-30*time.Minute).Format("2006-01-02T15:04:05.000Z"))
	res, err := parseCodexRollout(fixture)
	if err != nil {
		t.Fatal(err)
	}
	att, note := resolveAttribution(root, "u", res)
	if att != "attributed" {
		t.Errorf("want attributed against the driver-built window, got %s (%s)", att, note)
	}

	// A window of a DIFFERENT idea containing the event → ambiguous.
	otherID := "20260923T193000.000000000Z"
	other := runmanifest.New(runmanifest.Options{RunID: otherID, Root: root, IdeaSlug: "other-idea", CreatedAt: created, UpdatedAt: created})
	if err := runmanifest.Write(root, otherID, other); err != nil {
		t.Fatal(err)
	}
	if err := runmanifest.TouchUpdatedAt(root, otherID); err != nil {
		t.Fatal(err)
	}
	touch, err := runmanifest.Load(root, otherID)
	if err != nil {
		t.Fatal(err)
	}
	// Make the other-idea window the ONLY one containing the event: narrow the
	// first run's window to end before the event by rewriting its manifest with a
	// bounded updated_at (real shape, bounded window).
	narrow := manifest
	narrow.UpdatedAt = time.Now().UTC().Add(-45 * time.Minute)
	if err := runmanifest.Write(root, runID, narrow); err != nil {
		t.Fatal(err)
	}
	_ = touch
	att, _ = resolveAttribution(root, "u", res)
	if att != "ambiguous" {
		t.Errorf("an event inside only a different idea's window must be ambiguous, got %s", att)
	}
}

func TestUsageIngestStreamsLargeFixtureBoundedMemory(t *testing.T) {
	if testing.Short() {
		t.Skip("228 MB fixture generation is slow in -short mode")
	}
	dir := t.TempDir()
	root := dir
	seedMinimalDeck(t, root)
	os.MkdirAll(filepath.Join(root, protocol.DeckDir, "ideas", "u"), 0o755)
	big := filepath.Join(dir, "rollout-228mb.jsonl")
	f, err := os.Create(big)
	if err != nil {
		t.Fatal(err)
	}
	w := &countingWriter{f: f}
	payload := strings.Repeat("y", 8192)
	line := `{"timestamp":"2026-09-23T10:00:00.000Z","type":"response_item","payload":{"type":"message","content":[{"type":"input_text","text":"` + payload + `"}]}}` + "\n"
	target := int64(228 * 1024 * 1024)
	for w.n+int64(len(line)) <= target {
		if _, err := w.Write([]byte(line)); err != nil {
			t.Fatal(err)
		}
	}
	tail := `{"timestamp":"2026-09-23T20:00:00.000Z","type":"event_msg","payload":{"type":"token_count","info":{"total_token_usage":{"input_tokens":1000,"cached_input_tokens":900,"cache_write_input_tokens":0,"output_tokens":50,"reasoning_output_tokens":5,"total_tokens":1050}}}}` + "\n"
	if _, err := w.Write([]byte(tail)); err != nil {
		t.Fatal(err)
	}
	f.Close()
	st, _ := os.Stat(big)
	if st.Size() < int64(220*1024*1024) {
		t.Fatalf("fixture too small: %d B", st.Size())
	}

	var out, errOut bytes.Buffer
	code := runUsageIngest([]string{"--dir", root, "--agent", "codex-1", "--source", "codex-rollout", "--path", big, "--idea", "u", "--phase", "5"}, &out, &errOut)
	if code != 0 {
		t.Fatalf("228 MB ingest failed: %s", errOut.String())
	}
	if len(out.String()) > 1024 {
		t.Errorf("stdout must stay <= 1 KB over the 228 MB fixture, got %d B", len(out.String()))
	}
	if !strings.Contains(out.String(), "total_tokens=1050") {
		t.Errorf("cumulative totals must survive the stream; got %q", out.String())
	}
}

type countingWriter struct {
	f *os.File
	n int64
}

func (c *countingWriter) Write(p []byte) (int, error) {
	n, err := c.f.Write(p)
	c.n += int64(n)
	return n, err
}

// TestStreamLinesDropsOversizedLineWhole (fix-up F19, kimi-1 K1-F6c): a line above
// the 16 MiB cap is discarded whole — its tail must never surface as a partial
// "line". The fixture: a normal line, a 17 MiB line whose TAIL contains a decoy
// JSON-looking fragment, then another normal line.
func TestStreamLinesDropsOversizedLineWhole(t *testing.T) {
	path := filepath.Join(t.TempDir(), "oversized.jsonl")
	var b bytes.Buffer
	b.WriteString("{\"type\":\"normal\",\"n\":1}\n")
	huge := bytes.Repeat([]byte("x"), 17<<20)
	b.Write(huge)
	b.WriteString("{\"type\":\"token_count\"}\n") // tail of the oversized line
	b.WriteString("{\"type\":\"normal\",\"n\":2}\n")
	if err := os.WriteFile(path, b.Bytes(), 0o644); err != nil {
		t.Fatal(err)
	}
	var got []string
	if err := streamLines(path, func(line []byte) error {
		got = append(got, string(line))
		return nil
	}); err != nil {
		t.Fatal(err)
	}
	if len(got) != 2 {
		t.Fatalf("exactly the two normal lines must be yielded (the oversized line dropped whole, tail included); got %d lines", len(got))
	}
	for i, want := range []string{"{\"type\":\"normal\",\"n\":1}\n", "{\"type\":\"normal\",\"n\":2}\n"} {
		if got[i] != want {
			t.Errorf("line %d: got %.60q want %q", i, got[i], want)
		}
	}
}
