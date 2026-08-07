package app

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"parley-deck-cli/internal/protocolcore"
)

const testCore = `# Parley Deck Cooperation Protocol

**Workspace:** ` + "`<workspace-name>`" + `
**Transport:** ` + "`<transport-choice>`" + `
**Created:** ` + "`<YYYY-MM-DD>`" + `
**Protocol synced:** placeholder-from-core

## 2. Active agents (roster)

| Agent ID       | Workspace dir                       | Role          |
| -------------- | ----------------------------------- | ------------- |
| ` + "`example-1`" + `    | ` + "`./`" + `                                | participant   |

## 3. Phases

Core rules live here.
`

const testDeck = `# Parley Deck Cooperation Protocol

**Workspace:** ` + "`my-project`" + `
**Transport:** ` + "`local-dir`" + `
**Created:** 2026-05-09
**Protocol synced:** 2026-06-13 — old

## 2. Active agents (roster)

| Agent ID       | Workspace dir                       | Role          |
| -------------- | ----------------------------------- | ------------- |
| ` + "`claude-1`" + `     | ` + "`./`" + `                                | facilitator   |
| ` + "`codex-1`" + `      | ` + "`./`" + `                                | participant   |

## 3. Phases

STALE core text that must be replaced.

## 99. Project-specific rule

This deck added its own section.
`

// protocolFixture builds an isolated Parley home with one published core release and a deck that
// pins it.
func protocolFixture(t *testing.T, deckBody string) (root string) {
	t.Helper()
	home := t.TempDir()
	t.Setenv("PARLEY_HOME", home)
	if _, err := protocolcore.StoreAt(home).Publish("1.0.0", testCore); err != nil {
		t.Fatal(err)
	}
	root = t.TempDir()
	meta := filepath.Join(root, "parley-deck", "meta")
	if err := os.MkdirAll(meta, 0o755); err != nil {
		t.Fatal(err)
	}
	if deckBody != "" {
		if err := os.WriteFile(filepath.Join(root, "parley-deck", "COOPERATION.md"), []byte(deckBody), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	if err := os.WriteFile(filepath.Join(meta, "protocol-lock.yaml"), []byte("core-version: 1.0.0\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	return root
}

// G1: render is idempotent — a second apply produces byte-identical output and reports nothing to
// do. A generator that is not idempotent makes every render a diff, and the drift returns wearing
// a new name.
func TestProtocolRenderIsIdempotent(t *testing.T) {
	root := protocolFixture(t, testDeck)
	path := filepath.Join(root, "parley-deck", "COOPERATION.md")

	var out, errb strings.Builder
	if code := runProtocol([]string{"render", "--dir", root, "--yes"}, &out, &errb); code != 0 {
		t.Fatalf("first render exit=%d: %s", code, errb.String())
	}
	first, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}

	out.Reset()
	errb.Reset()
	if code := runProtocol([]string{"render", "--dir", root, "--yes"}, &out, &errb); code != 0 {
		t.Fatalf("second render exit=%d: %s", code, errb.String())
	}
	second, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	if string(first) != string(second) {
		t.Error("render is not idempotent: the second apply changed the file")
	}
	if !strings.Contains(out.String(), "nothing to do") {
		t.Errorf("second render did not report a no-op: %s", out.String())
	}
}

// G1: the identity slots survive and the stale core text is replaced. This is the whole point —
// the deck keeps WHO IT IS, the core supplies the rules.
func TestProtocolRenderPreservesIdentityAndReplacesCoreText(t *testing.T) {
	root := protocolFixture(t, testDeck)
	var out, errb strings.Builder
	if code := runProtocol([]string{"render", "--dir", root, "--yes"}, &out, &errb); code != 0 {
		t.Fatalf("exit=%d: %s", code, errb.String())
	}
	body, err := os.ReadFile(filepath.Join(root, "parley-deck", "COOPERATION.md"))
	if err != nil {
		t.Fatal(err)
	}
	got := string(body)

	for _, want := range []string{"`my-project`", "`local-dir`", "2026-05-09", "`claude-1`", "`codex-1`"} {
		if !strings.Contains(got, want) {
			t.Errorf("identity value %q was lost", want)
		}
	}
	if strings.Contains(got, "STALE core text") {
		t.Error("stale core text survived the render")
	}
	if !strings.Contains(got, "Core rules live here") {
		t.Error("core rules were not applied")
	}
	if strings.Contains(got, "`example-1`") {
		t.Error("the core's placeholder roster row leaked into the deck")
	}
	if !strings.Contains(got, "core 1.0.0") {
		t.Error("the synced stamp does not name the core release")
	}
}

// G1: a section the deck has and the core does not MUST be reported before it is removed —
// in preview AND on apply. Silent erasure of project content is a defect this project has
// already shipped once, against a real deck.
func TestProtocolRenderReportsRemovedSections(t *testing.T) {
	root := protocolFixture(t, testDeck)

	var preview, errb strings.Builder
	if code := runProtocol([]string{"render", "--dir", root, "--dry-run"}, &preview, &errb); code != 0 {
		t.Fatalf("preview exit=%d: %s", code, errb.String())
	}
	if !strings.Contains(preview.String(), "## 99. Project-specific rule") {
		t.Errorf("preview did not report the section it would remove:\n%s", preview.String())
	}
	if !strings.Contains(preview.String(), "Nothing was written") {
		t.Error("preview wrote without --yes")
	}

	var apply, errb2 strings.Builder
	if code := runProtocol([]string{"render", "--dir", root, "--yes"}, &apply, &errb2); code != 0 {
		t.Fatalf("apply exit=%d: %s", code, errb2.String())
	}
	if !strings.Contains(apply.String(), "## 99. Project-specific rule") {
		t.Errorf("apply removed a section without reporting it:\n%s", apply.String())
	}
}

// G7b: `protocol check` reports a hand-edit rather than silently overwriting, and says so through
// the real entry point with a non-zero exit a script can act on.
func TestProtocolCheckReportsHandEditAndNeverWrites(t *testing.T) {
	root := protocolFixture(t, testDeck)
	path := filepath.Join(root, "parley-deck", "COOPERATION.md")

	var out, errb strings.Builder
	if code := runProtocol([]string{"render", "--dir", root, "--yes"}, &out, &errb); code != 0 {
		t.Fatalf("render exit=%d: %s", code, errb.String())
	}
	out.Reset()
	if code := runProtocol([]string{"check", "--dir", root}, &out, &errb); code != 0 {
		t.Fatalf("check on a freshly rendered deck exit=%d: %s", code, out.String())
	}

	rendered, _ := os.ReadFile(path)
	edited := string(rendered) + "\n## 100. Sneaked in by hand\n"
	if err := os.WriteFile(path, []byte(edited), 0o644); err != nil {
		t.Fatal(err)
	}

	out.Reset()
	code := runProtocol([]string{"check", "--dir", root}, &out, &errb)
	if code == 0 {
		t.Error("check accepted a hand-edited protocol")
	}
	if !strings.Contains(out.String(), "does NOT match core") {
		t.Errorf("check did not name the mismatch: %s", out.String())
	}
	after, _ := os.ReadFile(path)
	if string(after) != edited {
		t.Error("check MODIFIED the file; it must report, never overwrite")
	}
}

// D8: a deck pinning a version that is not installed BLOCKS. It must never fall back to another
// installed release — substituting a same-named or current version is the failure this rule exists
// to prevent.
func TestProtocolBlocksWhenPinnedReleaseIsMissing(t *testing.T) {
	root := protocolFixture(t, testDeck)
	lock := filepath.Join(root, "parley-deck", "meta", "protocol-lock.yaml")
	if err := os.WriteFile(lock, []byte("core-version: 9.9.9\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	var out, errb strings.Builder
	if code := runProtocol([]string{"render", "--dir", root, "--yes"}, &out, &errb); code == 0 {
		t.Fatal("render succeeded with a missing pinned release")
	}
	if !strings.Contains(errb.String(), "9.9.9") {
		t.Errorf("error does not name the missing version: %s", errb.String())
	}
	body, _ := os.ReadFile(filepath.Join(root, "parley-deck", "COOPERATION.md"))
	if !strings.Contains(string(body), "STALE core text") {
		t.Error("the deck file was modified despite the block")
	}
}

// D8: an unpinned deck is reported, never guessed for. Picking "the newest installed" is exactly
// the substitution the lock exists to forbid.
func TestProtocolRefusesToGuessAVersion(t *testing.T) {
	root := protocolFixture(t, testDeck)
	if err := os.Remove(filepath.Join(root, "parley-deck", "meta", "protocol-lock.yaml")); err != nil {
		t.Fatal(err)
	}
	var out, errb strings.Builder
	if code := runProtocol([]string{"render", "--dir", root, "--yes"}, &out, &errb); code == 0 {
		t.Fatal("render picked a version for an unpinned deck")
	}
	if !strings.Contains(errb.String(), "never guessed") {
		t.Errorf("refusal does not explain itself: %s", errb.String())
	}
}

// D1/G2: releases are write-once. This is what makes the global immutable by construction rather
// than by rule — a change becomes a new version, it cannot overwrite an old one.
func TestCoreReleasesAreWriteOnce(t *testing.T) {
	home := t.TempDir()
	store := protocolcore.StoreAt(home)
	if _, err := store.Publish("1.0.0", testCore); err != nil {
		t.Fatal(err)
	}
	if _, err := store.Publish("1.0.0", "different bytes"); err == nil {
		t.Fatal("republishing an existing version was allowed")
	}
	rel, err := store.Load("1.0.0")
	if err != nil {
		t.Fatal(err)
	}
	if rel.Body != testCore {
		t.Error("an existing release was modified")
	}
}

// G7b: `protocol status --json` is a machine surface, so its shape is asserted, not assumed.
func TestProtocolStatusJSON(t *testing.T) {
	root := protocolFixture(t, testDeck)
	var out, errb strings.Builder
	if code := runProtocol([]string{"status", "--dir", root, "--json"}, &out, &errb); code != 0 {
		t.Fatalf("exit=%d: %s", code, errb.String())
	}
	var payload struct {
		Installed []string `json:"installed"`
		Pinned    string   `json:"pinned"`
		DeckSHA   string   `json:"deck_sha256"`
	}
	if err := json.Unmarshal([]byte(out.String()), &payload); err != nil {
		t.Fatalf("invalid JSON: %v\n%s", err, out.String())
	}
	if payload.Pinned != "1.0.0" || len(payload.Installed) != 1 || payload.DeckSHA == "" {
		t.Errorf("unexpected status payload: %+v", payload)
	}
}

// A deck with no COOPERATION.md at all renders one from the core: the core's own identity
// placeholders stand, and nothing is reported as removed.
func TestProtocolRenderOnAFreshDeck(t *testing.T) {
	root := protocolFixture(t, "")
	var out, errb strings.Builder
	if code := runProtocol([]string{"render", "--dir", root, "--yes"}, &out, &errb); code != 0 {
		t.Fatalf("exit=%d: %s", code, errb.String())
	}
	body, err := os.ReadFile(filepath.Join(root, "parley-deck", "COOPERATION.md"))
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(body), "Core rules live here") {
		t.Error("fresh deck did not receive the core body")
	}
	if strings.Contains(out.String(), "will be REMOVED") {
		t.Errorf("fresh deck reported removals: %s", out.String())
	}
}
