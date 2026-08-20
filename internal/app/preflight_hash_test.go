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
