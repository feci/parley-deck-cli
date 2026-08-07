package app

import (
	"encoding/json"
	"os"
	"path/filepath"
	"regexp"
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

// The removal report must be CONTENT-based, not heading-based. Content under a heading the core
// ALSO has was the gap: it vanished with nothing reported — the silent-erasure class G1 exists to
// prevent, and the one that destroyed a real deck's local section during the 2026-08-06 sync.
func TestProtocolRenderReportsContentLostUnderASharedHeading(t *testing.T) {
	deck := strings.Replace(testDeck,
		"STALE core text that must be replaced.",
		"STALE core text that must be replaced.\n\nPROJECT RULE: deployments need two approvals.", 1)
	root := protocolFixture(t, deck)

	var out, errb strings.Builder
	if code := runProtocol([]string{"render", "--dir", root, "--dry-run"}, &out, &errb); code != 0 {
		t.Fatalf("exit=%d: %s", code, errb.String())
	}
	// "## 3. Phases" exists in BOTH deck and core, so a heading diff reports nothing here.
	if !strings.Contains(out.String(), "## 3. Phases") {
		t.Errorf("content lost under a shared heading was not reported:\n%s", out.String())
	}
}

// A CRLF deck must not report every section as removed, and must not produce mixed endings.
func TestProtocolRenderHandlesCRLFDecks(t *testing.T) {
	root := protocolFixture(t, strings.ReplaceAll(testDeck, "\n", "\r\n"))
	path := filepath.Join(root, "parley-deck", "COOPERATION.md")

	var out, errb strings.Builder
	if code := runProtocol([]string{"render", "--dir", root, "--dry-run"}, &out, &errb); code != 0 {
		t.Fatalf("exit=%d: %s", code, errb.String())
	}
	if strings.Contains(out.String(), "## 2. Active agents") {
		t.Errorf("CRLF deck produced a false removal report:\n%s", out.String())
	}
	out.Reset()
	if code := runProtocol([]string{"render", "--dir", root, "--yes"}, &out, &errb); code != 0 {
		t.Fatalf("apply exit=%d: %s", code, errb.String())
	}
	body, _ := os.ReadFile(path)
	if strings.Contains(strings.ReplaceAll(string(body), "\r\n", ""), "\n") {
		t.Error("output has mixed line endings")
	}
}

// D8/G8: the version comes from a COMMITTED lock any contributor can edit, and it is joined onto
// the store root. It must never escape the store.
//
// The first version of this test was VACUOUS — kimi-1 caught it: it asserted only a non-zero exit,
// which render returns anyway because no release exists at the traversal target. It now asserts the
// REASON, and exercises Load directly against a store where the escape would otherwise succeed.
func TestProtocolRejectsPathTraversalInTheLock(t *testing.T) {
	root := protocolFixture(t, testDeck)
	lock := filepath.Join(root, "parley-deck", "meta", "protocol-lock.yaml")
	// NOTE: a leading space is NOT in this list — the lock parser trims the value before it reaches
	// ValidVersion, which is correct for a YAML-ish field. Untrimmed input is rejected at the API
	// level instead; see TestPublishRejectsUnsafeVersions.
	for _, bad := range []string{"../../../etc", "..", ".", "a/b"} {
		if err := os.WriteFile(lock, []byte("core-version: "+bad+"\n"), 0o644); err != nil {
			t.Fatal(err)
		}
		var out, errb strings.Builder
		code := runProtocol([]string{"render", "--dir", root, "--yes"}, &out, &errb)
		if code == 0 {
			t.Errorf("accepted unsafe version %q", bad)
		}
		if !strings.Contains(errb.String(), "unsafe version") {
			t.Errorf("version %q was rejected for the wrong reason: %s", bad, errb.String())
		}
	}
}

// The escape must be refused even when the traversal target really exists — otherwise the test
// passes for the wrong reason and would survive reverting the fix.
func TestLoadRefusesToEscapeTheStore(t *testing.T) {
	home := t.TempDir()
	store := protocolcore.StoreAt(home)
	outside := filepath.Join(home, "outside")
	if err := os.MkdirAll(outside, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(outside, protocolcore.CoreFileName), []byte("PLANTED"), 0o644); err != nil {
		t.Fatal(err)
	}
	// ../../outside resolves to a real, readable release body outside the store.
	rel, err := store.Load("../../outside")
	if err == nil {
		t.Fatalf("Load escaped the store and read %q", rel.Body)
	}
	if !strings.Contains(err.Error(), "unsafe version") {
		t.Errorf("rejected for the wrong reason: %v", err)
	}
}

// Write-once is per RELEASE, not per file: an existing release directory whose body is missing is
// still an existing release, and populating it would let a second publish complete a half-written
// one. A symlink at the release path must never be written through.
func TestPublishRefusesExistingReleaseDirAndSymlinks(t *testing.T) {
	home := t.TempDir()
	store := protocolcore.StoreAt(home)

	if err := os.MkdirAll(store.ReleaseDir("2.0.0"), 0o755); err != nil {
		t.Fatal(err)
	}
	if _, err := store.Publish("2.0.0", testCore); err == nil {
		t.Error("published into an existing release directory")
	}

	victim := filepath.Join(t.TempDir(), "victim.md")
	if err := os.WriteFile(victim, []byte("original"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.MkdirAll(store.ReleaseDir("3.0.0"), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.Symlink(victim, filepath.Join(store.ReleaseDir("3.0.0"), protocolcore.CoreFileName)); err != nil {
		t.Skipf("cannot create symlink: %v", err)
	}
	_, _ = store.Publish("3.0.0", "attacker body")
	got, _ := os.ReadFile(victim)
	if string(got) != "original" {
		t.Error("publish wrote THROUGH a symlink to a file outside the store")
	}
}

// An unsafe version must be rejected by Publish as well as Load.
func TestPublishRejectsUnsafeVersions(t *testing.T) {
	store := protocolcore.StoreAt(t.TempDir())
	for _, bad := range []string{"", ".", "..", "../x", "a/b", ".hidden"} {
		if _, err := store.Publish(bad, testCore); err == nil {
			t.Errorf("accepted unsafe version %q", bad)
		}
	}
}

// G7b demands the REAL command entry point. The tests above call runProtocol, which is one layer
// below dispatch; codex-1 was right that this can pass while the command is unreachable. This one
// goes through Run, exactly as `parley protocol …` does.
func TestProtocolIsReachableThroughProductionDispatch(t *testing.T) {
	root := protocolFixture(t, testDeck)
	var out, errb strings.Builder
	if code := Run([]string{"protocol", "status", "--dir", root, "--json"}, &out, &errb); code != 0 {
		t.Fatalf("dispatch exit=%d: out=%s err=%s", code, out.String(), errb.String())
	}
	if !strings.Contains(out.String(), `"pinned"`) {
		t.Errorf("dispatch did not reach protocol status: %s", out.String())
	}
	out.Reset()
	errb.Reset()
	if code := Run([]string{"protocol", "render", "--dir", root, "--yes"}, &out, &errb); code != 0 {
		t.Fatalf("render via dispatch exit=%d: %s", code, errb.String())
	}
	body, _ := os.ReadFile(filepath.Join(root, "parley-deck", "COOPERATION.md"))
	if !strings.Contains(string(body), "Core rules live here") {
		t.Error("render via dispatch did not apply the core")
	}
}

// G7b: the publisher's refusal is a stated guarantee, so it gets an end-to-end test through the
// real entry point rather than a claim in prose (hermes-1).
func TestPublishRefusesWithoutATerminal(t *testing.T) {
	home := t.TempDir()
	t.Setenv("PARLEY_HOME", home)
	src := filepath.Join(t.TempDir(), "core.md")
	if err := os.WriteFile(src, []byte(testCore), 0o644); err != nil {
		t.Fatal(err)
	}
	var out, errb strings.Builder
	// `go test` runs with stdin/stdout not a tty, which is exactly an agent run's shape.
	code := runProtocol([]string{"publish", "--version", "1.0.0", "--from", src}, &out, &errb)
	if code == 0 {
		t.Fatal("published without a controlling terminal")
	}
	if _, err := protocolcore.StoreAt(home).Load("1.0.0"); err == nil {
		t.Error("a release was written despite the refusal")
	}
	msg := errb.String()
	if !strings.Contains(msg, "attended") {
		t.Errorf("refusal does not explain itself: %s", msg)
	}
	// The refusal must not overclaim: a pty defeats it, and the text says so.
	if !strings.Contains(msg, "pty") && !strings.Contains(msg, "unavailable on this platform") {
		t.Errorf("refusal overclaims its own strength: %s", msg)
	}
}

// G7b: `protocol status` must report a read failure, not render it as an empty store (hermes-1).
func TestProtocolStatusReportsReadErrors(t *testing.T) {
	root := protocolFixture(t, testDeck)
	lock := filepath.Join(root, "parley-deck", "meta", "protocol-lock.yaml")
	if err := os.Remove(lock); err != nil {
		t.Fatal(err)
	}
	if err := os.Mkdir(lock, 0o755); err != nil { // a directory where a file is expected
		t.Fatal(err)
	}
	var out, errb strings.Builder
	if code := runProtocol([]string{"status", "--dir", root}, &out, &errb); code == 0 {
		t.Errorf("status reported success on an unreadable lock: %s", out.String())
	}
}

// Cross-section erasure: a line the deck has under §2 that the render has ONLY under §3.
//
// Per-section matching must report it as lost from §2. A GLOBAL multiset finds the identical line
// in §3 and reports nothing — which is exactly what cycle 2 shipped while documenting the fix as
// per-section. The first version of this test was itself tautological (hermes-1): its fixture put
// the line in the same section on both sides, so it passed with the fix reverted. Verified: with
// matching reverted to global, this version FAILS.
func TestProtocolRenderDetectsCrossSectionErasure(t *testing.T) {
	// "Core rules live here." lives under ## 3. in the core. Put it under ## 2. in the deck.
	deck := strings.Replace(testDeck,
		"## 3. Phases",
		"Core rules live here.\n\n## 3. Phases", 1)
	root := protocolFixture(t, deck)

	var out, errb strings.Builder
	if code := runProtocol([]string{"render", "--dir", root, "--dry-run"}, &out, &errb); code != 0 {
		t.Fatalf("exit=%d: %s", code, errb.String())
	}
	if !strings.Contains(out.String(), "## 2. Active agents") {
		t.Errorf("a line dropped from §2 was masked by an identical line surviving in §3:\n%s", out.String())
	}
}

// Load must refuse a symlinked store just as Publish does — a read that silently comes from
// elsewhere is how a deck ends up governed by a protocol nobody published.
func TestLoadRefusesASymlinkedStore(t *testing.T) {
	home := t.TempDir()
	real := t.TempDir()
	store := protocolcore.StoreAt(real)
	if _, err := store.Publish("1.0.0", testCore); err != nil {
		t.Fatal(err)
	}
	// ~/.parley/protocol -> <real>/protocol
	if err := os.MkdirAll(filepath.Join(home, "protocol"), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.RemoveAll(filepath.Join(home, "protocol")); err != nil {
		t.Fatal(err)
	}
	if err := os.Symlink(filepath.Join(real, "protocol"), filepath.Join(home, "protocol")); err != nil {
		t.Skipf("cannot symlink: %v", err)
	}
	if _, err := protocolcore.StoreAt(home).Load("1.0.0"); err == nil {
		t.Fatal("Load read through a symlinked store component")
	}
}

// Load must refuse a symlinked RELEASE DIRECTORY, not only a symlinked store component. kimi-1
// showed the two halves are separate code paths and only one was pinned by a test.
func TestLoadRefusesASymlinkedReleaseDirectory(t *testing.T) {
	home := t.TempDir()
	store := protocolcore.StoreAt(home)
	elsewhere := t.TempDir()
	if err := os.WriteFile(filepath.Join(elsewhere, protocolcore.CoreFileName), []byte("PLANTED"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.MkdirAll(store.Root, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.Symlink(elsewhere, store.ReleaseDir("1.0.0")); err != nil {
		t.Skipf("cannot symlink: %v", err)
	}
	rel, err := store.Load("1.0.0")
	if err == nil {
		t.Fatalf("Load read through a symlinked release directory and got %q", rel.Body)
	}
}

// Publish must refuse a symlinked store component too. The read half was pinned; the write half
// was not, while the changelog asserted both (kimi-1).
func TestPublishRefusesASymlinkedStoreComponent(t *testing.T) {
	home := t.TempDir()
	real := t.TempDir()
	if err := os.MkdirAll(filepath.Join(real, "protocol"), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.Symlink(filepath.Join(real, "protocol"), filepath.Join(home, "protocol")); err != nil {
		t.Skipf("cannot symlink: %v", err)
	}
	if _, err := protocolcore.StoreAt(home).Publish("1.0.0", testCore); err == nil {
		t.Fatal("published through a symlinked store component")
	}
}

// `protocol render` must report an unreadable deck file, not treat it as empty — rendering as
// though the deck had no content would report nothing lost, the silent-erasure shape again.
func TestProtocolRenderReportsDeckReadErrors(t *testing.T) {
	root := protocolFixture(t, testDeck)
	path := filepath.Join(root, "parley-deck", "COOPERATION.md")
	if err := os.Remove(path); err != nil {
		t.Fatal(err)
	}
	if err := os.Mkdir(path, 0o755); err != nil { // a directory where the file is expected
		t.Fatal(err)
	}
	// --dry-run, so the failure cannot come from the WRITE also failing: with the read error
	// swallowed, a dry run would happily report a full regeneration and exit 0.
	var out, errb strings.Builder
	code := runProtocol([]string{"render", "--dir", root, "--dry-run"}, &out, &errb)
	if code == 0 {
		t.Errorf("render treated an unreadable deck as empty: %s", out.String())
	}
	if !strings.Contains(errb.String(), "protocol render:") {
		t.Errorf("render did not report the read failure: %s", errb.String())
	}
}

// codex-1's three probes, as permanent tests: each is a way content changed meaning or vanished
// while the report stayed silent.
func TestDroppedContentDetectsIndentationLoss(t *testing.T) {
	// The core has unindented prose; the deck has it as an indented code block. TrimSpace called
	// them identical, so the meaning-changing edit passed as "carried forward".
	deck := strings.Replace(testDeck, "STALE core text that must be replaced.",
		"    Core rules live here.", 1)
	root := protocolFixture(t, deck)
	var out, errb strings.Builder
	if code := runProtocol([]string{"render", "--dir", root, "--dry-run"}, &out, &errb); code != 0 {
		t.Fatalf("exit=%d: %s", code, errb.String())
	}
	if !strings.Contains(out.String(), "## 3. Phases") {
		t.Errorf("indentation-significant content was lost without a report:\n%s", out.String())
	}
}

func TestDroppedContentDetectsSyncedPrefixOutsideTheHeader(t *testing.T) {
	// Genuine project prose that happens to start with the stamp prefix must not be skipped just
	// because the header's stamp is regenerated.
	deck := strings.Replace(testDeck, "STALE core text that must be replaced.",
		"**Protocol synced:** this deck documents its own sync policy here.", 1)
	root := protocolFixture(t, deck)
	var out, errb strings.Builder
	if code := runProtocol([]string{"render", "--dir", root, "--dry-run"}, &out, &errb); code != 0 {
		t.Fatalf("exit=%d: %s", code, errb.String())
	}
	if !strings.Contains(out.String(), "## 3. Phases") {
		t.Errorf("project prose beginning with the stamp prefix vanished unreported:\n%s", out.String())
	}
}

func TestDroppedContentTreatsDeeperHeadingsAsSections(t *testing.T) {
	// A #### subsection was not a section boundary, so content lost from it was attributed to —
	// and masked by — its parent.
	deck := strings.Replace(testDeck, "STALE core text that must be replaced.",
		"#### Parent\n\nfirst deck line\nsecond deck line", 1)
	root := protocolFixture(t, deck)
	var out, errb strings.Builder
	if code := runProtocol([]string{"render", "--dir", root, "--dry-run"}, &out, &errb); code != 0 {
		t.Fatalf("exit=%d: %s", code, errb.String())
	}
	// The point is ATTRIBUTION: the loss is charged to #### Parent, not merged into its parent
	// section. The exact count also includes the heading and the blank lines around it, which are
	// content too — an earlier version silently normalized whitespace-only lines away.
	if !strings.Contains(out.String(), "#### Parent — ") {
		t.Errorf("a level-4 subsection was not treated as its own section:\n%s", out.String())
	}
}

// codex-1's round-6 cases, as permanent tests.
//
// The fixtures matter as much as the assertions: an earlier version of these tests kept an extra
// dropped line in the section, so the report appeared under a multiset too and all three passed
// with the fix reverted (hermes-1). Each fixture below differs from the core ONLY by the
// transformation under test, so a multiset sees no difference and only a sequence diff does.
// Verified in both directions.

// protocolFixtureWith builds a deck against an arbitrary core body.
func protocolFixtureWith(t *testing.T, core, deck string) string {
	t.Helper()
	home := t.TempDir()
	t.Setenv("PARLEY_HOME", home)
	if _, err := protocolcore.StoreAt(home).Publish("1.0.0", core); err != nil {
		t.Fatal(err)
	}
	root := t.TempDir()
	meta := filepath.Join(root, "parley-deck", "meta")
	if err := os.MkdirAll(meta, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(root, "parley-deck", "COOPERATION.md"), []byte(deck), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(meta, "protocol-lock.yaml"), []byte("core-version: 1.0.0\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	return root
}

const minimalCore = `# T

**Workspace:** ` + "`<w>`" + `
**Transport:** ` + "`<t>`" + `
**Created:** ` + "`<d>`" + `

## 3. Phases

alpha
beta
gamma
`

func renderPreview(t *testing.T, root string) string {
	t.Helper()
	var out, errb strings.Builder
	if code := runProtocol([]string{"render", "--dir", root, "--dry-run"}, &out, &errb); code != 0 {
		t.Fatalf("exit=%d: %s", code, errb.String())
	}
	return out.String()
}

// A Markdown hard break (two trailing spaces) is meaning-bearing. Trimmed comparison equates the
// two lines; exact comparison does not. The ONLY difference from the core is those two spaces.
func TestDroppedContentDetectsMarkdownHardBreakLoss(t *testing.T) {
	deck := strings.Replace(minimalCore, "alpha", "alpha  ", 1)
	got := renderPreview(t, protocolFixtureWith(t, minimalCore, deck))
	if !strings.Contains(got, "## 3. Phases") {
		t.Errorf("a lost Markdown hard break was not reported:\n%s", got)
	}
}

// Same lines, different ORDER. A multiset is identical; a sequence diff is not.
func TestDroppedContentDetectsReorderedContent(t *testing.T) {
	deck := strings.Replace(minimalCore, "alpha\nbeta\ngamma", "gamma\nbeta\nalpha", 1)
	got := renderPreview(t, protocolFixtureWith(t, minimalCore, deck))
	if !strings.Contains(got, "## 3. Phases") {
		t.Errorf("reordering that changes meaning was not reported:\n%s", got)
	}
}

// Exactly ONE stamp is forgiven — the regenerated one. A second stamp-prefixed line is genuine
// project prose and must be reported.
func TestDroppedContentForgivesExactlyOneStamp(t *testing.T) {
	core := strings.Replace(minimalCore, "**Created:** `<d>`", "**Created:** `<d>`\n**Protocol synced:** placeholder", 1)

	// Converge once so the deck holds a REAL generated stamp, then re-render: the regenerated
	// stamp alone must not be reported. Only a recognizably-generated stamp is forgiven — a
	// legacy stamp in some other format is genuinely being replaced and is reported.
	root := protocolFixtureWith(t, core, core)
	var out, errb strings.Builder
	if code := runProtocol([]string{"render", "--dir", root, "--yes"}, &out, &errb); code != 0 {
		t.Fatalf("initial render exit=%d: %s", code, errb.String())
	}
	path := filepath.Join(root, "parley-deck", "COOPERATION.md")
	body, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	// Age the stamp: same generated FORMAT, different content — exactly what a version bump does.
	aged := regexp.MustCompile(`\*\*Protocol synced:\*\* core \S+ \([0-9a-f]+\)`).
		ReplaceAllString(string(body), "**Protocol synced:** core 0.9.0 (deadbeefcafe)")
	if err := os.WriteFile(path, []byte(aged), 0o644); err != nil {
		t.Fatal(err)
	}
	if got := renderPreview(t, root); strings.Contains(got, "not carried forward") {
		t.Errorf("a regenerated stamp was reported as a loss:\n%s", got)
	}
}

// codex-1's round-8 case: a deck that ALREADY carries the current generated stamp, plus genuine
// stamp-prefixed prose. Nothing was replaced, so nothing may be forgiven — and the prose must be
// reported. Verified to fail when the exemption reverts to a prefix or position rule.
func TestStampExemptionForgivesNothingWhenNothingWasReplaced(t *testing.T) {
	core := strings.Replace(minimalCore, "**Created:** `<d>`",
		"**Created:** `<d>`\n**Protocol synced:** placeholder", 1)
	root := protocolFixtureWith(t, core, core)

	// Converge the deck first, so it holds the exact generated stamp.
	var out, errb strings.Builder
	if code := runProtocol([]string{"render", "--dir", root, "--yes"}, &out, &errb); code != 0 {
		t.Fatalf("initial render exit=%d: %s", code, errb.String())
	}
	path := filepath.Join(root, "parley-deck", "COOPERATION.md")
	body, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	// Insert genuine prose immediately BEFORE the current stamp.
	lines := strings.Split(string(body), "\n")
	var withProse []string
	for _, l := range lines {
		if strings.HasPrefix(l, "**Protocol synced:**") {
			withProse = append(withProse, "**Protocol synced:** genuine project prose")
		}
		withProse = append(withProse, l)
	}
	if err := os.WriteFile(path, []byte(strings.Join(withProse, "\n")), 0o644); err != nil {
		t.Fatal(err)
	}

	got := renderPreview(t, root)
	if !strings.Contains(got, "not carried forward") {
		t.Errorf("genuine stamp-prefixed prose was silently forgiven:\n%s", got)
	}
}

// Hirschberg must agree with a straightforward full-DP LCS, not merely be small. Randomized
// comparison, since a linear-space rewrite is exactly where an off-by-one hides.
func TestHirschbergAgreesWithFullDP(t *testing.T) {
	alphabet := []string{"a", "b", "c", "d"}
	seed := 12345
	next := func(n int) int { seed = (seed*1103515245 + 12345) & 0x7fffffff; return seed % n }
	for trial := 0; trial < 200; trial++ {
		a := make([]string, next(12))
		for i := range a {
			a[i] = alphabet[next(len(alphabet))]
		}
		b := make([]string, next(12))
		for i := range b {
			b[i] = alphabet[next(len(alphabet))]
		}
		got := protocolcore.LCSMaskForTest(a, b)
		wantLen := fullDPLCSLength(a, b)
		gotLen := 0
		for _, m := range got {
			if m {
				gotLen++
			}
		}
		if gotLen != wantLen {
			t.Fatalf("trial %d: Hirschberg matched %d lines, full DP says LCS is %d\na=%v\nb=%v",
				trial, gotLen, wantLen, a, b)
		}
	}
}

func fullDPLCSLength(a, b []string) int {
	n, m := len(a), len(b)
	dp := make([][]int, n+1)
	for i := range dp {
		dp[i] = make([]int, m+1)
	}
	for i := n - 1; i >= 0; i-- {
		for j := m - 1; j >= 0; j-- {
			if a[i] == b[j] {
				dp[i][j] = dp[i+1][j+1] + 1
			} else if dp[i+1][j] >= dp[i][j+1] {
				dp[i][j] = dp[i+1][j]
			} else {
				dp[i][j] = dp[i][j+1]
			}
		}
	}
	return dp[0][0]
}

// codex-1's round-9 case: a deck that never had a generated stamp, with genuine stamp-prefixed
// prose occupying the slot after **Created:**. Nothing was generated there, so nothing may be
// forgiven — the prose must be reported.
func TestStampExemptionRequiresProofTheSlotWasGenerated(t *testing.T) {
	core := strings.Replace(minimalCore, "**Created:** `<d>`",
		"**Created:** `<d>`\n**Protocol synced:** placeholder", 1)
	deck := strings.Replace(core, "**Protocol synced:** placeholder",
		"**Protocol synced:** genuine project prose, never a generated stamp", 1)
	got := renderPreview(t, protocolFixtureWith(t, core, deck))
	if !strings.Contains(got, "not carried forward") {
		t.Errorf("prose in the stamp slot of a deck with no generated stamp was silently forgiven:\n%s", got)
	}
}

// The terminal-CR fix needs its own pin: a lone CR on the last line is content, not line-ending
// noise, and stripping it deleted real characters.
func TestTerminalCarriageReturnIsNotErased(t *testing.T) {
	// The core's LAST line is "trailing"; the deck's is "trailing\r". Stripping the lone CR makes
	// them compare equal and the change vanishes from the report. The earlier version of this test
	// appended a line the core lacked entirely, so it was reported either way — a tautology.
	core := minimalCore + "trailing"
	deck := minimalCore + "trailing\r"
	got := renderPreview(t, protocolFixtureWith(t, core, deck))
	if !strings.Contains(got, "not carried forward") {
		t.Errorf("a lone terminal carriage return was erased without a report:\n%s", got)
	}
}

// Whitespace-only lines are content; normalizing them away hid real changes.
func TestWhitespaceOnlyLinesAreCounted(t *testing.T) {
	deck := strings.Replace(minimalCore, "alpha\nbeta", "alpha\n   \nbeta", 1)
	got := renderPreview(t, protocolFixtureWith(t, minimalCore, deck))
	if !strings.Contains(got, "not carried forward") {
		t.Errorf("a whitespace-only line was normalized away instead of reported:\n%s", got)
	}
}

// `parley protocol --help` must print the group usage and succeed. It exited 2 with an
// unknown-subcommand error, so the documented way to discover the command failed (deploy
// verification, codex-1).
func TestProtocolHelpSucceeds(t *testing.T) {
	for _, arg := range []string{"--help", "-h", "help"} {
		var out, errb strings.Builder
		if code := runProtocol([]string{arg}, &out, &errb); code != 0 {
			t.Errorf("protocol %s exit=%d: %s", arg, code, errb.String())
		}
		if !strings.Contains(out.String(), "parley protocol status") {
			t.Errorf("protocol %s did not print the group usage: %s", arg, out.String())
		}
	}
}
