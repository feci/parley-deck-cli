package app

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"parley-deck-cli/internal/agents"
	"parley-deck-cli/internal/runmanifest"
)

// Two roster IDs may share an adapter and still run different models. Keying the frozen
// map by adapter made the second entry overwrite the first, so BOTH continuations
// launched the last entry's model — the freeze was wrong exactly where it mattered most.
func TestSnapshotPinsPerRosterIDNotPerAdapter(t *testing.T) {
	frozen := []runmanifest.RosterSnapshotEntry{
		{Agent: "claude-1", Adapter: "claude", Model: "opus-a", Effort: "max"},
		{Agent: "claude-2", Adapter: "claude", Model: "opus-b", Effort: "low"},
	}
	discovered := []agents.Discovery{
		{Spec: agents.Spec{ID: "claude-1", Model: "drifted"}},
		{Spec: agents.Spec{ID: "claude-2", Model: "drifted"}},
	}
	out := applyRosterSnapshot(discovered, frozen, nil)
	got := map[string]string{}
	for _, d := range out {
		got[d.Spec.ID] = d.Spec.Model
	}
	if got["claude-1"] != "opus-a" || got["claude-2"] != "opus-b" {
		t.Fatalf("per-ID pins collapsed: %v", got)
	}
}

// G1's acceptance clause names auto-args, not just model/effort/speed. A machine-config
// change that drops an auto-approve flag must not alter a running idea's autonomy.
func TestSnapshotPinsAutonomousLaunchArgs(t *testing.T) {
	frozen := []runmanifest.RosterSnapshotEntry{{
		Agent: "hermes-1", Adapter: "hermes", Model: "glm", Auto: true,
		LaunchArgs: []string{"--yolo", "--oneshot", "{prompt}"},
	}}
	discovered := []agents.Discovery{{Spec: agents.Spec{
		ID: "hermes-1", Model: "glm",
		HeadlessArgs: []string{"--oneshot", "{prompt}"}, // --yolo removed since the freeze
	}}}
	out := applyRosterSnapshot(discovered, frozen, nil)
	if !strings.Contains(strings.Join(out[0].Spec.HeadlessArgs, " "), "--yolo") {
		t.Fatalf("frozen auto-args not restored: %v", out[0].Spec.HeadlessArgs)
	}
}

// A --keep token that matches nothing is a typo. Silently accepting it meant
// `--keep kimi-1.modle --yes` protected nothing and removed kimi-1.model anyway.
func TestSyncRejectsUnmatchedKeepTokens(t *testing.T) {
	root := deckWith(t,
		"[roster.kimi-1]\nadapter = \"kimi\"\nmodel = \"deck-pin\"\n",
		"[roster.kimi-1]\nadapter = \"kimi\"\nmodel = \"machine\"\n")
	var out, errb strings.Builder
	code := rosterSync(root, multiFlag{"kimi-1.modle"}, false, true, &out, &errb)
	if code == 0 {
		t.Fatalf("typoed --keep accepted; exit=%d out=%s", code, out.String())
	}
	if !strings.Contains(errb.String(), "kimi-1.modle") {
		t.Errorf("error does not name the unmatched token: %s", errb.String())
	}
	// And nothing was written: the pin survives.
	b, _ := os.ReadFile(filepath.Join(root, "parley-deck", "agents.toml"))
	if !strings.Contains(string(b), "deck-pin") {
		t.Error("the pin the operator tried to protect was removed anyway")
	}
}

// The legacy normalizer is D7's second half: a config layer that hardcodes a model
// literal in headless_args must not outrank the `model` field beside it.
func TestLegacyModelArgsAreNormalizedToPlaceholders(t *testing.T) {
	got, changed := agents.NormalizeLegacyModelArgs(
		[]string{"-p", "--model", "ancient-literal", "--effort", "low", "--output-format", "text"})
	if !changed {
		t.Fatal("normalizer reported no change")
	}
	joined := strings.Join(got, " ")
	if strings.Contains(joined, "ancient-literal") {
		t.Errorf("model literal survived: %v", got)
	}
	if !strings.Contains(joined, agents.ModelPlaceholder) || !strings.Contains(joined, agents.EffortPlaceholder) {
		t.Errorf("placeholders not installed: %v", got)
	}
	// A boolean flag followed by another flag must be left alone.
	if out, changed := agents.NormalizeLegacyModelArgs([]string{"--model", "--verbose"}); changed {
		t.Errorf("rewrote a flag that takes no model value: %v", out)
	}
}

// The frozen contract claims text and JSON are the same table. A healthy row printed `ok`
// in text and marshalled to `null` in JSON — the same row contradicting itself.
func TestJSONStatusMatchesTextForAHealthyRow(t *testing.T) {
	root := deckWith(t, "[roster.claude-1]\nadapter = \"claude\"\nmodel = \"m\"\neffort = \"max\"\n", "")
	var out, errb strings.Builder
	if code := rosterShow(root, true, rosterViewOpts{}, &out, &errb); code != 0 {
		t.Fatalf("roster show --json exit=%d: %s", code, errb.String())
	}
	var payload struct {
		SchemaVersion int      `json:"schema_version"`
		Columns       []string `json:"columns"`
		Roster        []struct {
			Agent  string   `json:"agent"`
			Status []string `json:"status"`
		} `json:"roster"`
	}
	if err := json.Unmarshal([]byte(out.String()), &payload); err != nil {
		t.Fatalf("invalid JSON: %v\n%s", err, out.String())
	}
	if payload.SchemaVersion != RosterSchemaVersion || len(payload.Columns) != len(RosterColumns) {
		t.Errorf("contract header drifted: version=%d columns=%d", payload.SchemaVersion, len(payload.Columns))
	}
	for _, r := range payload.Roster {
		if r.Status == nil {
			t.Errorf("%s marshalled status as null; text renders `ok` for the same row", r.Agent)
		}
	}
}
