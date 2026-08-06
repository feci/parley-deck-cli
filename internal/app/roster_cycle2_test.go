package app

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"parley-deck-cli/internal/agents"
	"parley-deck-cli/internal/config"
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

// A5's residual: at the real continuation boundary the discoveries are still keyed by
// ADAPTER, so per-roster-ID pins collapsed even though applyRosterSnapshot itself was
// correct. The cycle-2 test proved the function; this one proves the call site.
func TestSnapshotPinsSurviveParticipantResolution(t *testing.T) {
	frozen := []runmanifest.RosterSnapshotEntry{
		{Agent: "claude-1", Adapter: "claude", Model: "opus-a"},
		{Agent: "claude-2", Adapter: "claude", Model: "opus-b"},
	}
	// What the continuation actually holds: one adapter-level discovery.
	discovered := []agents.Discovery{{Found: true, Spec: agents.Spec{ID: "claude", Model: "drifted"}}}
	mapping := map[string]string{"claude-1": "claude", "claude-2": "claude"}

	out := applyRosterSnapshotToParticipants([]string{"claude-1", "claude-2"}, discovered, mapping, frozen, nil)

	got := map[string]string{}
	for _, d := range out {
		got[d.Spec.ID] = d.Spec.Model
	}
	if got["claude-1"] != "opus-a" || got["claude-2"] != "opus-b" {
		t.Fatalf("per-ID pins collapsed at the continuation boundary: %v", got)
	}
}

// The run revision must change when only the frozen launch args change, or autonomy drift
// reports `current` and `stale-snapshot` never fires for it.
func TestRosterRevisionCoversLaunchArgs(t *testing.T) {
	a := []runmanifest.RosterSnapshotEntry{{Agent: "h-1", Adapter: "hermes", Auto: true, LaunchArgs: []string{"--yolo"}}}
	b := []runmanifest.RosterSnapshotEntry{{Agent: "h-1", Adapter: "hermes", Auto: true, LaunchArgs: []string{}}}
	if runmanifest.RosterRevisionOf(a) == runmanifest.RosterRevisionOf(b) {
		t.Fatal("revision ignores launch args; autonomy drift would report `current`")
	}
}

// A valid legacy §2 table is the deck's compatibility membership. A machine roster must
// not be inherited over it — that was ratified in the consensus and omitted in cycle 2.
func TestLegacySection2BeatsTheMachineRoster(t *testing.T) {
	root := deckWith(t, "", "[roster.claude-1]\nadapter=\"claude\"\n[roster.codex-1]\nadapter=\"codex\"\n"+
		"[roster.hermes-1]\nadapter=\"hermes\"\n[roster.kimi-1]\nadapter=\"kimi\"\n[roster.opencode-1]\nadapter=\"opencode\"\n")
	coop := "## 2. Active agents (roster)\n\n| Agent ID | Workspace dir | Role |\n| --- | --- | --- |\n" +
		"| `claude-1` | . | participant |\n| `kimi-1` | . | reviewer |\n"
	if err := os.WriteFile(filepath.Join(root, "parley-deck", "COOPERATION.md"), []byte(coop), 0o644); err != nil {
		t.Fatal(err)
	}
	active, inactive, ok := RosterMembership(root)
	if !ok {
		t.Fatal("membership not resolved")
	}
	if len(active)+len(inactive) != 2 || !active["claude-1"] || !active["kimi-1"] {
		t.Fatalf("legacy §2 membership overridden by the machine roster: active=%v inactive=%v", active, inactive)
	}
	rows, err := resolveRoster(root, nil, rosterViewOpts{})
	if err != nil {
		t.Fatal(err)
	}
	for _, r := range rows {
		if hasStatus(r, "inherited-roster") {
			t.Errorf("%s marked inherited-roster; a legacy deck declares its own roster", r.Agent)
		}
		if !hasStatus(r, "legacy-roster") {
			t.Errorf("%s missing legacy-roster status: %v", r.Agent, r.Status)
		}
	}
}

// A machine-scope write must not warn that the file it just wrote is masking itself.
func TestMachineScopeWriteIsNotReportedAsMasked(t *testing.T) {
	root := deckWith(t, "", "[roster.claude-1]\nadapter=\"claude\"\n")
	target := filepath.Join(os.Getenv("PARLEY_HOME"), "agents.toml")
	if _, masked := rosterFieldMaskedBy(root, "claude-1", "adapter", target); masked {
		t.Error("machine-scope write reported itself as masked")
	}
}

// `roster init` writes [roster.*] blocks, which is a membership change. Gating only
// `roster set` left the confirmation bypassable through the deprecated alias.
func TestRosterInitRequiresConfirmBreaking(t *testing.T) {
	root := deckWith(t, "", "")
	coop := "## 2. Active agents (roster)\n\n| Agent ID | Workspace dir | Role |\n| --- | --- | --- |\n" +
		"| `claude-1` | . | participant |\n"
	if err := os.WriteFile(filepath.Join(root, "parley-deck", "COOPERATION.md"), []byte(coop), 0o644); err != nil {
		t.Fatal(err)
	}
	var out, errb strings.Builder
	if code := runRoster([]string{"init", "--dir", root, "--yes"}, &out, &errb); code != 2 {
		t.Fatalf("init --yes alone wrote membership; exit=%d out=%s", code, out.String())
	}
	if _, err := os.Stat(filepath.Join(root, "parley-deck", "agents.toml")); err == nil {
		t.Error("init wrote agents.toml despite refusing")
	}
	out.Reset()
	errb.Reset()
	if code := runRoster([]string{"init", "--dir", root, "--yes", "--confirm-breaking"}, &out, &errb); code != 0 {
		t.Fatalf("init with --confirm-breaking exit=%d: %s", code, errb.String())
	}
}

// A participant spelled exactly like its adapter family (`claude`, not `claude-1`) must
// still get the frozen values. Cycle 3's append-based helper left the live discovery first,
// and the runner resolves by FIRST exact ID match — so the freeze never reached the launch.
func TestSnapshotFreezeReachesBareFamilyParticipants(t *testing.T) {
	frozen := []runmanifest.RosterSnapshotEntry{{Agent: "claude", Adapter: "claude", Model: "frozen-model"}}
	discovered := []agents.Discovery{{Found: true, Spec: agents.Spec{ID: "claude", Model: "drifted"}}}

	out := applyRosterSnapshotToParticipants([]string{"claude"}, discovered, map[string]string{}, frozen, nil)

	// Resolve exactly as the runner does, and assert the record it would launch.
	picked, err := agents.ResolveParticipant("claude", out, map[string]string{})
	if err != nil {
		t.Fatalf("resolve: %v", err)
	}
	if picked.Spec.Model != "frozen-model" {
		t.Fatalf("runner would launch %q, not the frozen %q", picked.Spec.Model, "frozen-model")
	}
}

// A value-only layer must not be able to retire or revive a committed deck member:
// membership state belongs to the record that grants membership.
func TestValueLayersCannotChangeMembershipState(t *testing.T) {
	root := deckWith(t,
		"[roster.claude-1]\nadapter=\"claude\"\n[roster.kimi-1]\nadapter=\"kimi\"\n",
		"[roster.claude-1]\nadapter=\"claude\"\nactive = false\n")
	// The gitignored local layer also tries to retire a committed member.
	if err := os.WriteFile(filepath.Join(root, "parley-deck", "agents.local.toml"),
		[]byte("[roster.kimi-1]\nactive = false\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	active, inactive, ok := RosterMembership(root)
	if !ok {
		t.Fatal("membership not resolved")
	}
	for _, id := range []string{"claude-1", "kimi-1"} {
		if inactive[id] || !active[id] {
			t.Errorf("%s was retired by a value-only layer; the deck file says active", id)
		}
	}
}

// The frozen v1 contract is eleven columns in BOTH renderings. The earlier test only
// checked that status was non-null, so re-adding display_name/note — or any twelfth field
// — would not have failed it. This asserts the exact key set, so a regression must fail.
func TestJSONRowHasExactlyTheFrozenColumns(t *testing.T) {
	root := deckWith(t, "[roster.claude-1]\nadapter = \"claude\"\nmodel = \"m\"\n", "")
	var out, errb strings.Builder
	if code := rosterShow(root, true, rosterViewOpts{}, &out, &errb); code != 0 {
		t.Fatalf("exit=%d: %s", code, errb.String())
	}
	var payload struct {
		Roster []map[string]any `json:"roster"`
	}
	if err := json.Unmarshal([]byte(out.String()), &payload); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}
	want := map[string]bool{
		"agent": true, "adapter": true, "state": true, "installed": true, "model": true,
		"model_family": true, "model_company": true, "effort": true, "speed": true,
		"autonomous": true, "status": true,
	}
	if len(payload.Roster) == 0 {
		t.Fatal("no rows")
	}
	for _, row := range payload.Roster {
		for k := range row {
			if !want[k] {
				t.Errorf("row carries %q, which is outside the frozen eleven columns", k)
			}
		}
		for k := range want {
			if _, ok := row[k]; !ok {
				t.Errorf("row is missing frozen column %q", k)
			}
		}
	}
}

// `active` follows the membership authority, so its PROVENANCE and its masking warning
// must name that authority too. Cycle 4 fixed the behavior and left both reporting
// surfaces attributing state to a layer the resolver ignores — the diagnostics
// contradicting the result they describe.
func TestActiveProvenanceAndMaskingFollowTheAuthority(t *testing.T) {
	root := deckWith(t, "[roster.claude-1]\nadapter=\"claude\"\n", "[roster.claude-1]\nadapter=\"claude\"\nactive = false\n")
	env := filepath.Join(t.TempDir(), "env.toml")
	if err := os.WriteFile(env, []byte("[roster.claude-1]\nactive = false\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	t.Setenv("PARLEY_HEADLESS_AGENT_CONFIG", env)

	src, err := config.RosterStateSource(root)
	if err != nil {
		t.Fatal(err)
	}
	if src != "parley-deck/agents.toml" {
		t.Errorf("state provenance = %q, want the membership authority", src)
	}
	// A write to the authority must not be reported as masked by an ignored layer.
	target := filepath.Join(root, "parley-deck", "agents.toml")
	if who, masked := rosterFieldMaskedBy(root, "claude-1", "active", target); masked {
		t.Errorf("write to the authority reported as masked by %q", who)
	}
}

// The gate fires on a real state FLIP, not on writing a value the member already has.
func TestMembershipGateIgnoresNoOpStateWrites(t *testing.T) {
	if got := membershipChange([]string{"+ active = true"}, true, true); got != "" {
		t.Errorf("writing active=true to an already-active member gated as %q", got)
	}
	if got := membershipChange([]string{"+ active = false"}, true, true); got == "" {
		t.Error("retiring an active member was not gated")
	}
	if got := membershipChange([]string{"+ active = true"}, true, false); got == "" {
		t.Error("reviving a retired member was not gated")
	}
}

// State provenance must name the authority for the SCOPE being asked about. Cycle 5's
// override re-derived the deck's authority unconditionally, so `--scope machine --explain`
// named parley-deck/agents.toml while its own membership header named the machine file.
func TestActiveProvenanceIsScopeAware(t *testing.T) {
	root := deckWith(t,
		"[roster.claude-1]\nadapter=\"claude\"\nactive = false\n",
		"[roster.claude-1]\nadapter=\"claude\"\nactive = true\n")

	for _, tc := range []struct{ scope, wantState string }{
		{"deck", "inactive"},
		{"machine", "active"},
	} {
		var out, errb strings.Builder
		if code := rosterExplain(root, "claude-1", rosterViewOpts{scope: tc.scope}, &out, &errb); code != 0 {
			t.Fatalf("%s: exit=%d %s", tc.scope, code, errb.String())
		}
		text := out.String()
		header := strings.SplitN(text, "\n", 2)[0]
		var activeLine string
		for _, l := range strings.Split(text, "\n") {
			if strings.HasPrefix(l, "active") {
				activeLine = l
			}
		}
		if !strings.Contains(activeLine, tc.wantState) {
			t.Errorf("%s scope: active line %q, want state %q", tc.scope, activeLine, tc.wantState)
		}
		// The provenance must agree with the membership header of the same output.
		src := strings.TrimSpace(strings.TrimPrefix(header, "claude-1 — membership from "))
		src = strings.TrimSuffix(src, " (INHERITED — this deck declares no roster of its own)")
		if !strings.Contains(activeLine, src) {
			t.Errorf("%s scope: provenance %q contradicts its own membership header %q", tc.scope, activeLine, src)
		}
	}
}
