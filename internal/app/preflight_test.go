package app

import (
	"bytes"
	"context"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"parley-deck-cli/internal/agents"
	"parley-deck-cli/internal/protocol"
)

// --- mergePreservingZones (pure, unit-tested in isolation) ---

func TestMergePreservingZones(t *testing.T) {
	project := strings.Join([]string{
		"# COOPERATION.md — Multi-Agent Cooperation Protocol",
		"**Protocol synced:** 2026-01-01 — old",
		"",
		"## 0. Choose the transport",
		"PROJECT-TRANSPORT",
		"## 1. Scope and purpose",
		"PROJECT-SCOPE",
		"## 2. Active agents (roster)",
		"| Agent ID | Workspace | Role |",
		"| -------- | --------- | ---- |",
		"| claude-1 | ../claude | participant |",
		"## 3. Directory layout",
		"PROJECT-OLD-LAYOUT",
		"## 4. Protocol",
		"PROJECT-OLD-PHASES",
	}, "\n")

	packaged := strings.Join([]string{
		"# COOPERATION.md — Multi-Agent Cooperation Protocol",
		"",
		"## 0. Choose the transport",
		"PACKAGED-TRANSPORT-SHOULD-BE-DROPPED",
		"## 1. Scope and purpose",
		"PACKAGED-SCOPE-SHOULD-BE-DROPPED",
		"## 2. Active agents (roster)",
		"| Agent ID | Workspace | Role |",
		"| -------- | --------- | ---- |",
		"## 3. Directory layout",
		"PACKAGED-NEW-LAYOUT",
		"## 4. Protocol",
		"PACKAGED-NEW-PHASES",
		"## 5. New section",
		"BRAND-NEW",
	}, "\n")

	merged, err := mergePreservingZones(project, packaged)
	if err != nil {
		t.Fatalf("merge failed: %v", err)
	}

	// Project zones (§0/§1/§2 + header) preserved.
	for _, want := range []string{"PROJECT-TRANSPORT", "PROJECT-SCOPE", "| claude-1 | ../claude | participant |"} {
		if !strings.Contains(merged, want) {
			t.Fatalf("merged dropped preserved project zone %q:\n%s", want, merged)
		}
	}
	// §3-onward body replaced from packaged.
	for _, want := range []string{"PACKAGED-NEW-LAYOUT", "PACKAGED-NEW-PHASES", "BRAND-NEW"} {
		if !strings.Contains(merged, want) {
			t.Fatalf("merged missing packaged body %q:\n%s", want, merged)
		}
	}
	// Old project §3+ body and packaged §0-§2 zones dropped.
	for _, notWant := range []string{"PROJECT-OLD-LAYOUT", "PROJECT-OLD-PHASES", "PACKAGED-TRANSPORT-SHOULD-BE-DROPPED", "PACKAGED-SCOPE-SHOULD-BE-DROPPED"} {
		if strings.Contains(merged, notWant) {
			t.Fatalf("merged kept zone it should have replaced %q:\n%s", notWant, merged)
		}
	}
	// The §2/§3 boundary is intact and ordered.
	if i2, i3 := strings.Index(merged, "## 2."), strings.Index(merged, "## 3."); i2 < 0 || i3 < 0 || i2 > i3 {
		t.Fatalf("section ordering broken: §2 at %d, §3 at %d", i2, i3)
	}
}

func TestMergePreservingZonesMissingSectionThree(t *testing.T) {
	body := "# Title\n## 0. transport\n## 1. scope\n## 2. roster\n"
	if _, err := mergePreservingZones(body, body); err == nil {
		t.Fatal("expected an error when no '## 3.' section line is present")
	}
}

func TestRefreshProtocolSyncedLineReplacesExisting(t *testing.T) {
	body := "# Title\n**Protocol synced:** 2020-01-01 — old\n## 0. transport\n"
	out := refreshProtocolSyncedLine(body, "2.0.0", "2.0.0")
	if strings.Contains(out, "2020-01-01") {
		t.Fatalf("old sync line not replaced:\n%s", out)
	}
	if n := strings.Count(out, "**Protocol synced:**"); n != 1 {
		t.Fatalf("want exactly one Protocol synced line, got %d:\n%s", n, out)
	}
}

// --- classifyBump (semver of deckVersion) ---

func TestClassifyBump(t *testing.T) {
	cases := []struct {
		name             string
		project, packaged string
		want             bumpKind
	}{
		{"minor is additive", "1.3.1", "1.4.0", bumpPatchMinor},
		{"patch is additive", "1.3.1", "1.3.2", bumpPatchMinor},
		{"same is additive", "1.3.1", "1.3.1", bumpPatchMinor},
		{"major is breaking", "1.3.1", "2.0.0", bumpMajor},
		{"unparseable degrades to additive", "garbage", "1.4.0", bumpPatchMinor},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if got := classifyBump(tc.project, tc.packaged); got != tc.want {
				t.Fatalf("classifyBump(%q,%q)=%v want %v", tc.project, tc.packaged, got, tc.want)
			}
		})
	}
}

// --- freshness classifier (source advisory / unknown role gate / no version.json) ---

func writeVersionJSON(t *testing.T, root string, meta map[string]any) {
	t.Helper()
	dir := filepath.Join(root, protocol.DeckDir, "meta")
	if err := os.MkdirAll(dir, 0o755); err != nil {
		t.Fatal(err)
	}
	data, err := json.Marshal(meta)
	if err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(dir, "version.json"), data, 0o644); err != nil {
		t.Fatal(err)
	}
}

func TestClassifyFreshnessSourceIsAdvisoryNoWrite(t *testing.T) {
	root := t.TempDir()
	deckPath := filepath.Join(root, protocol.DeckDir, "COOPERATION.md")
	if err := os.MkdirAll(filepath.Dir(deckPath), 0o755); err != nil {
		t.Fatal(err)
	}
	original := "# COOPERATION.md\n## 3. layout\nORIGINAL\n"
	if err := os.WriteFile(deckPath, []byte(original), 0o644); err != nil {
		t.Fatal(err)
	}
	writeVersionJSON(t, root, map[string]any{
		"protocolRole":           "source",
		"protocolSha256":         "aaa",
		"packagedProtocolSha256": "bbb", // differs, but source never writes
		"deckVersion":            "1.3.1",
	})

	fr, gates, err := classifyAndSyncFreshness(context.Background(), preflightOptions{Root: root})
	if err != nil {
		t.Fatal(err)
	}
	if fr.Classification != "source-advisory" {
		t.Fatalf("classification=%q want source-advisory", fr.Classification)
	}
	if len(gates) != 0 {
		t.Fatalf("source role must produce no gate, got %v", gates)
	}
	if got, _ := os.ReadFile(deckPath); string(got) != original {
		t.Fatalf("source role wrote COOPERATION.md; must never write")
	}
}

func TestClassifyFreshnessUnknownRoleGate(t *testing.T) {
	root := t.TempDir()
	writeVersionJSON(t, root, map[string]any{
		"protocolSha256":         "aaa",
		"packagedProtocolSha256": "bbb",
		"deckVersion":            "1.3.1",
		// no protocolRole
	})
	fr, gates, err := classifyAndSyncFreshness(context.Background(), preflightOptions{Root: root})
	if err != nil {
		t.Fatal(err)
	}
	if fr.Classification != "unknown-role" {
		t.Fatalf("classification=%q want unknown-role", fr.Classification)
	}
	if len(gates) != 1 || gates[0].Kind != gateUnknownRole {
		t.Fatalf("want one unknown-role gate, got %v", gates)
	}
}

func TestClassifyFreshnessConsumerBreakingGate(t *testing.T) {
	root := t.TempDir()
	writeVersionJSON(t, root, map[string]any{
		"protocolRole":           "consumer",
		"protocolSha256":         "aaa",
		"packagedProtocolSha256": "bbb", // drift
		"deckVersion":            "1.3.1",
	})
	// No parley-deck-skill on PATH in tests, so packaged deck version is unknown
	// and classifyBump degrades to additive. To exercise the breaking path
	// deterministically, classifyBump is unit-tested directly above. Here we
	// assert the consumer-drift path at least does NOT report in-sync.
	fr, _, err := classifyAndSyncFreshness(context.Background(), preflightOptions{Root: root})
	if err != nil {
		t.Fatal(err)
	}
	if fr.Classification == "in-sync" {
		t.Fatalf("consumer drift must not be reported as in-sync")
	}
	if fr.Classification != "additive" && fr.Classification != "breaking" {
		t.Fatalf("consumer drift classification=%q want additive or breaking", fr.Classification)
	}
}

func TestClassifyFreshnessConsumerInSync(t *testing.T) {
	root := t.TempDir()
	writeVersionJSON(t, root, map[string]any{
		"protocolRole":           "consumer",
		"protocolSha256":         "same",
		"packagedProtocolSha256": "same",
		"deckVersion":            "1.3.1",
	})
	fr, gates, err := classifyAndSyncFreshness(context.Background(), preflightOptions{Root: root})
	if err != nil {
		t.Fatal(err)
	}
	if fr.Classification != "in-sync" {
		t.Fatalf("classification=%q want in-sync", fr.Classification)
	}
	if len(gates) != 0 {
		t.Fatalf("in-sync must produce no gate, got %v", gates)
	}
}

func TestClassifyFreshnessNoVersionJSON(t *testing.T) {
	root := t.TempDir()
	fr, gates, err := classifyAndSyncFreshness(context.Background(), preflightOptions{Root: root})
	if err != nil {
		t.Fatalf("missing version.json must not error: %v", err)
	}
	if fr.Classification != "no-version-json" {
		t.Fatalf("classification=%q want no-version-json", fr.Classification)
	}
	if len(gates) != 0 {
		t.Fatalf("missing version.json must produce no gate, got %v", gates)
	}
}

func TestSyncConsumerProtocolWritesAndRecords(t *testing.T) {
	root := t.TempDir()
	deckPath := filepath.Join(root, protocol.DeckDir, "COOPERATION.md")
	if err := os.MkdirAll(filepath.Dir(deckPath), 0o755); err != nil {
		t.Fatal(err)
	}
	project := "# COOPERATION.md\n## 0. transport\nKEEP-ME\n## 1. scope\n## 2. roster\n## 3. layout\nOLD-BODY\n"
	if err := os.WriteFile(deckPath, []byte(project), 0o644); err != nil {
		t.Fatal(err)
	}
	packaged := "# COOPERATION.md\n## 0. transport\nDROP\n## 1. scope\n## 2. roster\n## 3. layout\nNEW-BODY\n"
	skillStatus := map[string]any{
		"installer": map[string]any{"version": "1.4.0"},
		"packaged":  map[string]any{"protocol": packaged, "deckVersion": "1.4.0"},
	}
	meta := versionMeta{ProtocolRole: "consumer", ProtocolSha256: "old", PackagedProtocolSha256: "new", DeckVersion: "1.3.1"}

	synced, summary, err := syncConsumerProtocol(context.Background(), root, meta, skillStatus)
	if err != nil {
		t.Fatal(err)
	}
	if !synced {
		t.Fatalf("expected a sync write, got no-op: %s", summary)
	}
	got, _ := os.ReadFile(deckPath)
	if !strings.Contains(string(got), "KEEP-ME") || strings.Contains(string(got), "OLD-BODY") || !strings.Contains(string(got), "NEW-BODY") {
		t.Fatalf("sync did not preserve zones / replace body:\n%s", got)
	}
	if !strings.Contains(string(got), "**Protocol synced:**") {
		t.Fatalf("sync did not refresh the Protocol synced line:\n%s", got)
	}
	// A sync record was dropped under meta/.
	entries, _ := os.ReadDir(filepath.Join(root, protocol.DeckDir, "meta"))
	foundRecord := false
	for _, e := range entries {
		if strings.HasPrefix(e.Name(), "protocol-sync_") {
			foundRecord = true
		}
	}
	if !foundRecord {
		t.Fatalf("no protocol-sync_* record written under meta/")
	}
}

// --- exit-code paths via the shared preflight helper (fake probe) ---

func sourceWorkspace(t *testing.T) string {
	t.Helper()
	root := t.TempDir()
	writeVersionJSON(t, root, map[string]any{
		"protocolRole":           "source",
		"protocolSha256":         "aaa",
		"packagedProtocolSha256": "aaa",
		"deckVersion":            "1.3.1",
	})
	return root
}

func found(id string) agents.Discovery {
	return agents.Discovery{Spec: agents.Spec{ID: id, Commands: []string{id}}, Found: true, Version: "v1"}
}

func withFakeProbe(t *testing.T, fn probeFunc) {
	t.Helper()
	prev := pingProbe
	pingProbe = fn
	t.Cleanup(func() { pingProbe = prev })
}

func TestPreflightReadyAllAvailable(t *testing.T) {
	root := sourceWorkspace(t)
	withFakeProbe(t, func(context.Context, string, agents.Discovery, time.Duration) (bool, string) {
		return true, ""
	})
	_, code, err := preflight(context.Background(), preflightOptions{Root: root}, []agents.Discovery{found("a"), found("b")}, &bytes.Buffer{}, &bytes.Buffer{})
	if err != nil {
		t.Fatal(err)
	}
	if code != 0 {
		t.Fatalf("code=%d want 0", code)
	}
}

func TestPreflightUnavailableAgentIsGateExit3(t *testing.T) {
	root := sourceWorkspace(t)
	withFakeProbe(t, func(_ context.Context, _ string, a agents.Discovery, _ time.Duration) (bool, string) {
		if a.ID == "b" {
			return false, "unavailable:no-pong"
		}
		return true, ""
	})
	// Three available + one unavailable: excluding the one still leaves >= 2, so
	// the unavailable agent is a pending gate (exit 3), not a hard failure.
	report, code, err := preflight(context.Background(), preflightOptions{Root: root},
		[]agents.Discovery{found("a"), found("b"), found("c")}, &bytes.Buffer{}, &bytes.Buffer{})
	if err != nil {
		t.Fatal(err)
	}
	if code != 3 {
		t.Fatalf("code=%d want 3", code)
	}
	if len(report.Gates) != 1 || report.Gates[0].Kind != gateExcludeAgent {
		t.Fatalf("want one exclude-agent gate, got %v", report.Gates)
	}
}

func TestPreflightLessThanTwoParticipantsIsHardFailExit1(t *testing.T) {
	root := sourceWorkspace(t)
	withFakeProbe(t, func(_ context.Context, _ string, a agents.Discovery, _ time.Duration) (bool, string) {
		// Only "a" is available; excluding the rest would leave 1 participant.
		return a.ID == "a", "unavailable:no-pong"
	})
	_, code, err := preflight(context.Background(), preflightOptions{Root: root},
		[]agents.Discovery{found("a"), found("b")}, &bytes.Buffer{}, &bytes.Buffer{})
	if err != nil {
		t.Fatal(err)
	}
	if code != 1 {
		t.Fatalf("code=%d want 1 (§1 non-solo hard-stop)", code)
	}
}

func TestPreflightNoPingPresenceOnly(t *testing.T) {
	root := sourceWorkspace(t)
	withFakeProbe(t, func(context.Context, string, agents.Discovery, time.Duration) (bool, string) {
		t.Fatal("probe must not be called in --no-ping mode")
		return false, ""
	})
	missing := agents.Discovery{Spec: agents.Spec{ID: "gone", Commands: []string{"gone"}}, Found: false}
	report, code, err := preflight(context.Background(), preflightOptions{Root: root, NoPing: true},
		[]agents.Discovery{found("a"), found("b"), missing}, &bytes.Buffer{}, &bytes.Buffer{})
	if err != nil {
		t.Fatal(err)
	}
	// Present agents are available by presence; the missing one is a gate.
	if code != 3 {
		t.Fatalf("code=%d want 3 (missing agent gate)", code)
	}
	var missingReason string
	for _, e := range report.Roster {
		if e.RosterID == "gone" {
			missingReason = e.Reason
		}
	}
	if missingReason != "unavailable:missing" {
		t.Fatalf("missing agent reason=%q want unavailable:missing", missingReason)
	}
}

// --- command-level exit codes (usage / ready) ---

func TestPreflightUsageError(t *testing.T) {
	var stdout, stderr bytes.Buffer
	code := Run([]string{"preflight", "extra-arg"}, &stdout, &stderr)
	if code != 2 {
		t.Fatalf("code=%d want 2 (usage)", code)
	}
}
