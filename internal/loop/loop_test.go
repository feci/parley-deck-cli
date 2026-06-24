package loop

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

func fixedNow() time.Time { return time.Date(2026, 6, 24, 12, 0, 0, 0, time.UTC) }

// hasFrontmatterKey reports whether any line is a top-level YAML key `<key>:` (a
// real frontmatter claim) — as opposed to a prose mention of the token mid-line.
func hasFrontmatterKey(body, key string) bool {
	for _, line := range strings.Split(body, "\n") {
		if strings.HasPrefix(strings.TrimSpace(line), key+":") {
			return true
		}
	}
	return false
}

func readPrompt(t *testing.T, deck, slug string) string {
	t.Helper()
	data, err := os.ReadFile(filepath.Join(deck, "ideas", slug, "00-prompt.md"))
	if err != nil {
		t.Fatalf("read prompt for %s: %v", slug, err)
	}
	return string(data)
}

// Disabled-by-default: a tick with Enabled=false writes nothing and reports disabled.
func TestTickDisabledWritesNothing(t *testing.T) {
	deck := filepath.Join(t.TempDir(), "parley-deck")
	sigs := []Candidate{{Source: "commit", ID: "abc123", Title: "fix x"}}
	res, err := Tick(deck, Config{Enabled: false}, sigs, fixedNow())
	if err != nil {
		t.Fatal(err)
	}
	if res.Enabled || len(res.Created) != 0 {
		t.Fatalf("disabled tick must create nothing; got %+v", res)
	}
	if _, err := os.Stat(filepath.Join(deck, "ideas")); !os.IsNotExist(err) {
		t.Fatal("disabled tick must not create the ideas dir")
	}
}

// Enabled: drafts a status: candidate prompt with no quorum claim and a Promotion note.
func TestTickEnabledDraftsCandidate(t *testing.T) {
	deck := filepath.Join(t.TempDir(), "parley-deck")
	sigs := []Candidate{{Source: "ci", ID: "build-42", Title: "flaky test", Detail: "TestFoo retries"}}
	res, err := Tick(deck, Config{Enabled: true}, sigs, fixedNow())
	if err != nil {
		t.Fatal(err)
	}
	if len(res.Created) != 1 {
		t.Fatalf("expected 1 candidate, got %+v", res)
	}
	body := readPrompt(t, deck, res.Created[0])
	if !strings.Contains(body, "status: candidate") {
		t.Error("prompt must be status: candidate")
	}
	if hasFrontmatterKey(body, "participants") {
		t.Error("a loop-drafted candidate must NOT claim a participants: quorum (§14)")
	}
	if !strings.Contains(body, "## Promotion") {
		t.Error("prompt must carry a ## Promotion note")
	}
	if !strings.Contains(body, "build-42") || !strings.Contains(body, "TestFoo retries") {
		t.Error("prompt must carry the signal provenance (id + detail)")
	}
}

// Dedupe: a second tick over the same signal skips the existing candidate.
func TestTickDedupes(t *testing.T) {
	deck := filepath.Join(t.TempDir(), "parley-deck")
	sigs := []Candidate{{Source: "issue", ID: "GH-7", Title: "bug"}}
	if _, err := Tick(deck, Config{Enabled: true}, sigs, fixedNow()); err != nil {
		t.Fatal(err)
	}
	res, err := Tick(deck, Config{Enabled: true}, sigs, fixedNow())
	if err != nil {
		t.Fatal(err)
	}
	if len(res.Created) != 0 || len(res.Skipped) != 1 {
		t.Fatalf("second tick must skip the existing candidate; got %+v", res)
	}
}

// Fingerprint defaults to a stable hash of source+id when none is supplied, and an
// explicit fingerprint is honored (so the same logical signal dedupes across ID churn).
func TestSlugFingerprint(t *testing.T) {
	a := SlugFor(Candidate{Source: "commit", ID: "abc"})
	b := SlugFor(Candidate{Source: "commit", ID: "abc"})
	if a != b {
		t.Fatalf("same signal must yield same slug: %s vs %s", a, b)
	}
	if SlugFor(Candidate{Source: "commit", ID: "xyz"}) == a {
		t.Fatal("different IDs must yield different slugs")
	}
	exp := SlugFor(Candidate{Source: "commit", ID: "1", Fingerprint: "feature-login"})
	if !strings.Contains(exp, "feature-login") {
		t.Fatalf("explicit fingerprint must drive the slug; got %s", exp)
	}
}

// ReadSignals: a missing file is empty + no error (a cron tick with nothing queued).
func TestReadSignalsMissing(t *testing.T) {
	sigs, err := ReadSignals(filepath.Join(t.TempDir(), "nope.json"))
	if err != nil || len(sigs) != 0 {
		t.Fatalf("missing signals file must be empty/no-error; got %v / %v", sigs, err)
	}
}

func TestReadSignalsParses(t *testing.T) {
	p := filepath.Join(t.TempDir(), "signals.json")
	os.WriteFile(p, []byte(`[{"source":"commit","id":"a1","title":"t"}]`), 0o644)
	sigs, err := ReadSignals(p)
	if err != nil || len(sigs) != 1 || sigs[0].ID != "a1" {
		t.Fatalf("expected one parsed signal; got %v / %v", sigs, err)
	}
}

// ReadConfig: absent → disabled (no error); malformed → fails closed (disabled + error).
func TestReadConfig(t *testing.T) {
	deck := filepath.Join(t.TempDir(), "parley-deck")
	cfg, err := ReadConfig(deck)
	if err != nil || cfg.Enabled {
		t.Fatalf("absent config must be disabled/no-error; got %+v / %v", cfg, err)
	}
	os.MkdirAll(filepath.Join(deck, "loop"), 0o755)
	os.WriteFile(ConfigPath(deck), []byte(`{"enabled": true}`), 0o644)
	if cfg, _ := ReadConfig(deck); !cfg.Enabled {
		t.Fatal("enabled config must parse to Enabled=true")
	}
	os.WriteFile(ConfigPath(deck), []byte(`not json`), 0o644)
	if cfg, err := ReadConfig(deck); err == nil || cfg.Enabled {
		t.Fatal("malformed config must fail closed (disabled + error)")
	}
}
