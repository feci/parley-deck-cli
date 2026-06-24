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
	// An explicit fingerprint drives identity and is deterministic.
	fp1 := SlugFor(Candidate{Source: "commit", ID: "1", Fingerprint: "feature-login"})
	fp2 := SlugFor(Candidate{Source: "commit", ID: "2", Fingerprint: "feature-login"})
	if fp1 != fp2 {
		t.Fatalf("same explicit fingerprint must yield the same slug regardless of id: %s vs %s", fp1, fp2)
	}
	// AF2: lossy-sanitize collisions are gone — `a/b` and `a:b` are distinct identities.
	if SlugFor(Candidate{Source: "manual", ID: "x", Fingerprint: "a/b"}) ==
		SlugFor(Candidate{Source: "manual", ID: "x", Fingerprint: "a:b"}) {
		t.Fatal("distinct fingerprints a/b and a:b must NOT collide to one slug (AF2)")
	}
}

// AF2: the default source+id digest is unambiguous — a colon-boundary shift
// (`ci:`+`build` vs `ci`+`:build`) must not collide.
func TestColonBoundaryNoCollision(t *testing.T) {
	a := SlugFor(Candidate{Source: "ci:", ID: "build"})
	b := SlugFor(Candidate{Source: "ci", ID: ":build"})
	if a == b {
		t.Fatalf("colon-boundary-shifted signals must not collide: %s == %s", a, b)
	}
}

// AF1 (CRITICAL): a signal with newline-injected YAML in its source/id must NOT be
// able to write extra frontmatter keys — exactly one `status:` line (candidate), and
// no `participants:`/`checks:` keys, can appear in the drafted prompt.
func TestTickRejectsFrontmatterInjection(t *testing.T) {
	deck := filepath.Join(t.TempDir(), "parley-deck")
	evil := []Candidate{{
		Source: "commit", // valid source so it is not rejected outright
		ID:     "abc\nstatus: round-01\nparticipants: [evil]\nchecks: rm -rf /",
		Title:  "x\nparticipants: [evil2]",
	}}
	res, err := Tick(deck, Config{Enabled: true}, evil, fixedNow())
	if err != nil {
		t.Fatal(err)
	}
	if len(res.Created) != 1 {
		t.Fatalf("expected the candidate to be drafted (sanitized), got %+v", res)
	}
	body := readPrompt(t, deck, res.Created[0])
	// Frontmatter spans up to the second `---`. Count top-level keys in it.
	statusLines, partLines, checkLines := 0, 0, 0
	for _, line := range strings.Split(body, "\n") {
		t := strings.TrimSpace(line)
		switch {
		case strings.HasPrefix(t, "status:"):
			statusLines++
		case strings.HasPrefix(t, "participants:"):
			partLines++
		case strings.HasPrefix(t, "checks:"):
			checkLines++
		}
	}
	if statusLines != 1 {
		t.Fatalf("expected exactly one status: line (candidate), got %d:\n%s", statusLines, body)
	}
	if partLines != 0 || checkLines != 0 {
		t.Fatalf("injection leaked participants:/checks: keys:\n%s", body)
	}
	if !strings.Contains(body, "status: candidate") {
		t.Fatalf("the one status line must be candidate:\n%s", body)
	}
}

// AF7: an empty slug dir left behind by a crashed/failed prior tick must NOT suppress
// the signal forever — the next tick (re)writes the prompt (heals); a later tick skips.
func TestTickHealsPoisonedEmptyDir(t *testing.T) {
	deck := filepath.Join(t.TempDir(), "parley-deck")
	sig := []Candidate{{Source: "commit", ID: "heal-me"}}
	slug := SlugFor(sig[0])
	os.MkdirAll(filepath.Join(deck, "ideas", slug), 0o755) // simulate a crashed prior tick: dir, no prompt

	res, err := Tick(deck, Config{Enabled: true}, sig, fixedNow())
	if err != nil {
		t.Fatal(err)
	}
	if len(res.Created) != 1 {
		t.Fatalf("an empty (poisoned) slug dir must be healed, not skipped forever; got %+v", res)
	}
	if _, err := os.Stat(filepath.Join(deck, "ideas", slug, "00-prompt.md")); err != nil {
		t.Fatalf("00-prompt.md must exist after healing: %v", err)
	}
	res2, _ := Tick(deck, Config{Enabled: true}, sig, fixedNow())
	if len(res2.Created) != 0 || len(res2.Skipped) != 1 {
		t.Fatalf("once the prompt exists, the next tick must skip; got %+v", res2)
	}
}

// AF6: a multi-line Detail keeps its line breaks (rendered as an indented block in the
// body), instead of being flattened to one run-on line.
func TestTickPreservesMultilineDetail(t *testing.T) {
	deck := filepath.Join(t.TempDir(), "parley-deck")
	sig := []Candidate{{Source: "ci", ID: "log-1", Detail: "panic: boom\n\tat foo()\n\tat bar()"}}
	res, err := Tick(deck, Config{Enabled: true}, sig, fixedNow())
	if err != nil || len(res.Created) != 1 {
		t.Fatalf("expected one candidate; got %+v / %v", res, err)
	}
	body := readPrompt(t, deck, res.Created[0])
	if !strings.Contains(body, "## Signal detail") {
		t.Fatalf("expected a Signal detail section:\n%s", body)
	}
	// The newline must SURVIVE (not be flattened to a space): the first detail line is
	// indented and followed by a real line break, with the later lines also present.
	if !strings.Contains(body, "    panic: boom\n") || !strings.Contains(body, "at bar()") {
		t.Fatalf("multi-line detail must be preserved as an indented block (newlines kept):\n%s", body)
	}
}

// AF8: Unicode line/paragraph separators in a frontmatter field are flattened, so they
// cannot inject a frontmatter key even if a real YAML parser is adopted later.
func TestTickFlattensUnicodeSeparators(t *testing.T) {
	deck := filepath.Join(t.TempDir(), "parley-deck")
	ls, ps := string(rune(0x2028)), string(rune(0x2029)) // U+2028 line sep, U+2029 para sep
	sig := []Candidate{{Source: "commit", ID: "x" + ls + "status: round-01" + ps + "participants: [evil]"}}
	res, err := Tick(deck, Config{Enabled: true}, sig, fixedNow())
	if err != nil || len(res.Created) != 1 {
		t.Fatalf("expected one candidate; got %+v / %v", res, err)
	}
	body := readPrompt(t, deck, res.Created[0])
	statusLines, partLines := 0, 0
	for _, line := range strings.Split(body, "\n") {
		tl := strings.TrimSpace(line)
		if strings.HasPrefix(tl, "status:") {
			statusLines++
		}
		if strings.HasPrefix(tl, "participants:") {
			partLines++
		}
	}
	if statusLines != 1 || partLines != 0 {
		t.Fatalf("U+2028/U+2029 must not inject keys (status=%d participants=%d):\n%s", statusLines, partLines, body)
	}
}

// AF1: an unknown source is rejected (fail closed) — no idea is drafted for it.
func TestTickRejectsUnknownSource(t *testing.T) {
	deck := filepath.Join(t.TempDir(), "parley-deck")
	res, err := Tick(deck, Config{Enabled: true}, []Candidate{{Source: "evil", ID: "1"}}, fixedNow())
	if err != nil {
		t.Fatal(err)
	}
	if len(res.Created) != 0 || len(res.Rejected) != 1 {
		t.Fatalf("unknown source must be rejected, not drafted; got %+v", res)
	}
	if _, err := os.Stat(filepath.Join(deck, "ideas")); !os.IsNotExist(err) {
		t.Fatal("a rejected-only tick must not create any idea")
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
