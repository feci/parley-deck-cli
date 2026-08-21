package app

import (
	"context"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"parley-deck-cli/internal/protocol"
)

// codex-1/F24: two MISSING hashes compare equal, so a fresh consumer deck that recorded neither
// hash was reported "in sync" whatever its protocol actually said — including an altered one.
// An equality test is evidence only when both sides exist.
func TestFreshnessWithNoHashToCompareIsNotInSync(t *testing.T) {
	root := t.TempDir()
	if err := protocol.InitWorkspace(root); err != nil {
		t.Fatal(err)
	}
	metaPath := filepath.Join(root, protocol.DeckDir, "meta", "version.json")
	body, err := json.Marshal(map[string]string{
		"protocolRole": "consumer",
		"deckVersion":  "1.0.0",
		// no protocolSha256, no packagedProtocolSha256 — exactly the fresh-deck shape
	})
	if err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(metaPath, body, 0o644); err != nil {
		t.Fatal(err)
	}

	fr, gates, err := classifyAndSyncFreshness(context.Background(), preflightOptions{Root: root})
	if err != nil {
		t.Fatalf("classifyAndSyncFreshness: %v", err)
	}
	if fr.Classification == "in-sync" {
		t.Fatalf("reported in-sync with nothing to compare: %+v", fr)
	}
	if !strings.Contains(fr.Summary, "unknown") {
		t.Errorf("summary should say freshness is unknown, got %q", fr.Summary)
	}
	var found bool
	for _, g := range gates {
		if g.Kind == gateUnknownFreshness {
			found = true
		}
	}
	if !found {
		t.Errorf("no unknown-freshness gate raised: %+v", gates)
	}
}

// codex-1, MAJOR (review round-02): the F24 gate above had no `--yes` branch, so a freshly
// initialized deck — which records no packaged hash — was trapped: the gate demanded a
// confirmation, and no flag could give it. `preflight --yes` could never report a new deck ready.
func TestConfirmedFreshnessClearsTheGateOnAFreshDeck(t *testing.T) {
	root := t.TempDir()
	if err := protocol.InitWorkspace(root); err != nil {
		t.Fatal(err)
	}
	metaPath := filepath.Join(root, protocol.DeckDir, "meta", "version.json")
	body, err := json.Marshal(map[string]string{
		"protocolRole": "consumer",
		"deckVersion":  "1.0.0",
	})
	if err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(metaPath, body, 0o644); err != nil {
		t.Fatal(err)
	}

	// Without --yes the gate stands (that is TestFreshnessWithNoHashToCompareIsNotInSync).
	fr, gates, err := classifyAndSyncFreshness(context.Background(), preflightOptions{Root: root, Yes: true})
	if err != nil {
		t.Fatalf("classifyAndSyncFreshness: %v", err)
	}
	if len(gates) != 0 {
		t.Fatalf("--yes left a gate standing, so the deck is still unreadyable: %+v", gates)
	}
	if fr.Classification != "freshness-confirmed" {
		t.Errorf("classification=%q, want freshness-confirmed", fr.Classification)
	}
	// The confirmation must not claim a comparison it did not make.
	if strings.Contains(fr.Summary, "in sync") && fr.PackagedSha == "" {
		t.Errorf("claimed sync with no packaged hash: %q", fr.Summary)
	}

	// The hash it recorded must be the hash of the protocol actually on disk.
	deckBody, err := os.ReadFile(filepath.Join(root, protocol.DeckDir, "COOPERATION.md"))
	if err != nil {
		t.Fatal(err)
	}
	want := sha256Hex(string(deckBody))
	if fr.LiveSha != want {
		t.Errorf("recorded protocolSha256=%q, want the live protocol's %q", fr.LiveSha, want)
	}

	// And it must be persisted, or the next run re-raises the same gate.
	persisted, err := readVersionMeta(root)
	if err != nil {
		t.Fatal(err)
	}
	if persisted.ProtocolSha256 != want {
		t.Errorf("version.json protocolSha256=%q, want %q", persisted.ProtocolSha256, want)
	}
	if persisted.DeckVersion != "1.0.0" || persisted.ProtocolRole != "consumer" {
		t.Errorf("the write clobbered pre-existing keys: %+v", persisted)
	}
}

// `parley init` must record the hash of the protocol it just wrote — the metadata should describe
// the deck from its first minute — and must NOT invent a packaged hash it never computed.
func TestInitRecordsTheLiveProtocolHashAndNoPackagedHash(t *testing.T) {
	root := t.TempDir()
	if err := protocol.InitWorkspace(root); err != nil {
		t.Fatal(err)
	}
	meta, err := readVersionMeta(root)
	if err != nil {
		t.Fatal(err)
	}
	deckBody, err := os.ReadFile(filepath.Join(root, protocol.DeckDir, "COOPERATION.md"))
	if err != nil {
		t.Fatal(err)
	}
	if want := sha256Hex(string(deckBody)); meta.ProtocolSha256 != want {
		t.Errorf("init recorded protocolSha256=%q, want %q", meta.ProtocolSha256, want)
	}
	if meta.PackagedProtocolSha256 != "" {
		t.Errorf("init invented a packaged hash it never computed: %q", meta.PackagedProtocolSha256)
	}
}
