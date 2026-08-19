---
agent: codex-1
idea: roster-membership-overlay
round: 3
date: 2026-08-19
responding-to: [user/path-c-ruling, claude-1/path-c-measurement, codex-1/signoff-block]
---

## E1 result

[PRIMARY] **Result: yes, the resolver can separate a values-only block from membership. The minimal
content-keyed prototype satisfies E1(a), but it is not shippable.** It changes 35 of the 38 active
member sets found in the current workspace fleet, including nine decks that resolve to zero active
members and seven that resolve to partial sets. That is stronger evidence for an explicit schema
boundary than the earlier estimate that every declared deck would merely gain one member.

### Isolation and prototype

[PRIMARY] I created the required copy before changing code. The shared tree was not an experiment
root. The actual command and output were:

```text
$ experiment_dir=$(mktemp -d /tmp/codex-r3-path-c.XXXXXX)
$ cp -R internal cmd go.mod go.sum parley-deck "$experiment_dir/"
$ chmod -R u+w "$experiment_dir"
$ find "$experiment_dir" -maxdepth 1 -mindepth 1 -print | sort
/tmp/codex-r3-path-c.Psdv91/cmd
/tmp/codex-r3-path-c.Psdv91/go.mod
/tmp/codex-r3-path-c.Psdv91/go.sum
/tmp/codex-r3-path-c.Psdv91/internal
/tmp/codex-r3-path-c.Psdv91/parley-deck
/tmp/codex-r3-path-c.Psdv91
```

[PRIMARY] The mandated copy set omitted `VERSION`, which one baseline app test reads through
`../../VERSION`. The first unmodified-copy run therefore failed only as follows:

```text
$ go test ./...
--- FAIL: TestVersionFileMatchesBinaryVersion (0.00s)
    app_test.go:208: open ../../VERSION: no such file or directory
FAIL
FAIL    parley-deck-cli/internal/app
```

[PRIMARY] I copied that one additional read-only source file into the temp module and reran the
unmodified baseline. It was green:

```text
$ cp '/Volumes/My Shared Files/AI_WORKSPACE/parley-deck/parley-deck-cli/VERSION' \
    /tmp/codex-r3-path-c.Psdv91/VERSION
$ go test ./...
?       parley-deck-cli/cmd/parley                 [no test files]
ok      parley-deck-cli/internal/acp               (cached)
ok      parley-deck-cli/internal/agents            (cached)
ok      parley-deck-cli/internal/app               50.425s
ok      parley-deck-cli/internal/config            (cached)
ok      parley-deck-cli/internal/consensus         (cached)
ok      parley-deck-cli/internal/driver            (cached)
ok      parley-deck-cli/internal/fsutil            (cached)
ok      parley-deck-cli/internal/hitl              (cached)
ok      parley-deck-cli/internal/loop              (cached)
ok      parley-deck-cli/internal/pipeline          (cached)
ok      parley-deck-cli/internal/procctl           (cached)
ok      parley-deck-cli/internal/protocol          (cached)
ok      parley-deck-cli/internal/protocolcore      (cached)
ok      parley-deck-cli/internal/repomap           (cached)
ok      parley-deck-cli/internal/retro             (cached)
ok      parley-deck-cli/internal/runaction         (cached)
ok      parley-deck-cli/internal/runcontrol        (cached)
ok      parley-deck-cli/internal/runmanifest       (cached)
ok      parley-deck-cli/internal/runner            (cached)
ok      parley-deck-cli/internal/runplan           (cached)
ok      parley-deck-cli/internal/runstate          (cached)
ok      parley-deck-cli/internal/sessionstore      (cached)
ok      parley-deck-cli/internal/steer             (cached)
ok      parley-deck-cli/internal/store             (cached)
ok      parley-deck-cli/internal/track             (cached)
ok      parley-deck-cli/internal/tui               (cached)
```

[PRIMARY] I patched only `/tmp/codex-r3-path-c.Psdv91/internal/config/runtime.go`. For this
experiment, an explicit `active` key was the only existing block field treated as membership
intent; `adapter`, `model`, `effort`, `speed`, and render metadata were values. I also prevented a
modern value-only deck file from falling into its generated §2 projection before reaching the
machine roster. The operative diff was:

```diff
@@ func LoadRosterScoped(root string) (RosterScope, error) {
     deckMembers := map[string]bool{}
     machineMembers := map[string]bool{}
+    deckHasRosterBlocks := false
@@
+    deckHasRosterBlocks = true
     for _, id := range ids {
+        if entries[id].Active == nil {
+            continue
+        }
         deckMembers[id] = true
-        deckActive[id] = entries[id].Active == nil || *entries[id].Active
+        deckActive[id] = *entries[id].Active
     }
@@
-if active, inactive, ok := protocol.ReadRosterIDs(root); ok {
+if !deckHasRosterBlocks {
+    if active, inactive, ok := protocol.ReadRosterIDs(root); ok {
         // unchanged legacy-§2 resolution
+    }
 }
```

This discriminator was deliberately cheap so the fleet experiment could falsify it. It is not my
proposed file contract.

### E1(a): six inherited members plus the deck speed override

[PRIMARY] I built the patched binary and made a copied fixture whose only config block was exactly
the requested one:

```text
$ go build -o /tmp/codex-r3-path-c.Psdv91/parley-path-c ./cmd/parley
$ cat /tmp/codex-r3-path-c.Psdv91/e1-deck/parley-deck/agents.toml
[roster.kimi-1]
speed = "fast"

$ parley --version
parley 1.45.0
$ /tmp/codex-r3-path-c.Psdv91/parley-path-c --version
parley 1.45.0
```

[PRIMARY] The shipped binary is the control: it reports one member.

```text
$ env -u PARLEY_HEADLESS_AGENT_CONFIG parley roster show \
    --dir /tmp/codex-r3-path-c.Psdv91/e1-deck
AGENT        ADAPTER    STATE    INSTALLED MODEL                  MODEL-FAMILY   MODEL-COMPANY EFFORT   SPEED    AUTO STATUS
kimi-1       kimi       active   yes       kimi-code/k3           Kimi K         Moonshot AI   max      fast     yes  effort-from-config
```

[PRIMARY] The patched binary reports all six machine members and keeps `kimi-1` at `fast`:

```text
$ env -u PARLEY_HEADLESS_AGENT_CONFIG \
    /tmp/codex-r3-path-c.Psdv91/parley-path-c roster show \
    --dir /tmp/codex-r3-path-c.Psdv91/e1-deck
AGENT        ADAPTER    STATE    INSTALLED MODEL                  MODEL-FAMILY   MODEL-COMPANY EFFORT   SPEED    AUTO STATUS
claude-1     claude     active   yes       claude-opus-5[1m]      Claude Opus    Anthropic     max      deep     yes  inherited-roster
codex-1      codex      active   yes       gpt-5.6-sol            GPT            OpenAI        max      deep     yes  inherited-roster
hermes-1     hermes     active   yes       fireworks/inkling      Inkling        Thinking Machines Lab high     deep     yes  inherited-roster
kimi-1       kimi       active   yes       kimi-code/k3           Kimi K         Moonshot AI   max      fast     yes  inherited-roster,effort-from-config
opencode-1   opencode   active   yes       litellm/xai/grok-4.6   Grok           xAI           xhigh    deep     yes  inherited-roster,effort-from-config
zcode-1      zcode      active   yes       zai/glm-5.3            GLM            Zhipu AI      max      deep     yes  inherited-roster,model-from-config,effort-from-config
```

[PRIMARY] `--explain` shows that membership and speed came from different properties/layers:

```text
$ env -u PARLEY_HEADLESS_AGENT_CONFIG \
    /tmp/codex-r3-path-c.Psdv91/parley-path-c roster show \
    --dir /tmp/codex-r3-path-c.Psdv91/e1-deck --explain kimi-1
kimi-1 — membership from ~/.parley/agents.toml (INHERITED — this deck declares no roster of its own)

FIELD          EFFECTIVE                SET BY
adapter        kimi                     ~/.parley/agents.toml
model          kimi-code/k3             ~/.parley/agents.toml
effort         max                      ~/.kimi-code/config.toml (agent's own, read at launch)
speed          fast                     parley-deck/agents.toml
active         active                   ~/.parley/agents.toml

status: inherited-roster,effort-from-config
effort source: read live from ~/.kimi-code/config.toml -> [thinking] effort — kimi has no per-invocation effort flag, so parley passes none and kimi reads this file itself at launch.
```

### E1(b): full test result and classification

[PRIMARY] This is the raw output from the exact requested command against the patched copy:

```text
$ go test ./...
?       parley-deck-cli/cmd/parley [no test files]
ok      parley-deck-cli/internal/acp       (cached)
ok      parley-deck-cli/internal/agents    (cached)
--- FAIL: TestZcodeRowReportsModelAndEffortReadFromAgentConfig (1.14s)
    roster_configread_test.go:61: no zcode-1 row in:
        AGENT        ADAPTER    STATE    INSTALLED MODEL                  MODEL-FAMILY   MODEL-COMPANY EFFORT   SPEED    AUTO STATUS
        antigravity-1 agy        active   yes       Gemini 3.6 Flash (High) Gemini         Google        unknown  balanced yes  legacy-roster,unmapped,effort-unknown
          ⚠ unmapped — declare the adapter with `parley roster set antigravity-1 --scope deck --adapter <family> --yes --confirm-breaking`
        claude-1     claude     active   yes       claude-opus-4-8[1m]    Claude Opus    Anthropic     max      balanced yes  legacy-roster,unmapped
          ⚠ unmapped — declare the adapter with `parley roster set claude-1 --scope deck --adapter <family> --yes --confirm-breaking`
        codex-1      codex      active   yes       unknown                unknown        unknown       unknown  balanced yes  legacy-roster,unmapped,model-unbound,effort-unknown,metadata-unknown
          ⚠ unmapped — declare the adapter with `parley roster set codex-1 --scope deck --adapter <family> --yes --confirm-breaking`
--- FAIL: TestDeckDeclaredModelNeverOverridesAgentConfigForUnbindableAdapter (1.23s)
    roster_configread_test.go:79: no zcode-1 row in:
        AGENT        ADAPTER    STATE    INSTALLED MODEL                  MODEL-FAMILY   MODEL-COMPANY EFFORT   SPEED    AUTO STATUS
        antigravity-1 agy        active   yes       Gemini 3.6 Flash (High) Gemini         Google        unknown  balanced yes  legacy-roster,unmapped,effort-unknown
          ⚠ unmapped — declare the adapter with `parley roster set antigravity-1 --scope deck --adapter <family> --yes --confirm-breaking`
        claude-1     claude     active   yes       claude-opus-4-8[1m]    Claude Opus    Anthropic     max      balanced yes  legacy-roster,unmapped
          ⚠ unmapped — declare the adapter with `parley roster set claude-1 --scope deck --adapter <family> --yes --confirm-breaking`
        codex-1      codex      active   yes       unknown                unknown        unknown       unknown  balanced yes  legacy-roster,unmapped,model-unbound,effort-unknown,metadata-unknown
          ⚠ unmapped — declare the adapter with `parley roster set codex-1 --scope deck --adapter <family> --yes --confirm-breaking`
--- FAIL: TestUnreadableAgentConfigFallsBackToUnknown (1.05s)
    roster_configread_test.go:95: no zcode-1 row in:
        AGENT        ADAPTER    STATE    INSTALLED MODEL                  MODEL-FAMILY   MODEL-COMPANY EFFORT   SPEED    AUTO STATUS
        antigravity-1 agy        active   yes       Gemini 3.6 Flash (High) Gemini         Google        unknown  balanced yes  legacy-roster,unmapped,effort-unknown
          ⚠ unmapped — declare the adapter with `parley roster set antigravity-1 --scope deck --adapter <family> --yes --confirm-breaking`
        claude-1     claude     active   yes       claude-opus-4-8[1m]    Claude Opus    Anthropic     max      balanced yes  legacy-roster,unmapped
          ⚠ unmapped — declare the adapter with `parley roster set claude-1 --scope deck --adapter <family> --yes --confirm-breaking`
        codex-1      codex      active   yes       unknown                unknown        unknown       unknown  balanced yes  legacy-roster,unmapped,model-unbound,effort-unknown,metadata-unknown
          ⚠ unmapped — declare the adapter with `parley roster set codex-1 --scope deck --adapter <family> --yes --confirm-breaking`
--- FAIL: TestJSONStatusMatchesTextForAHealthyRow (1.08s)
    roster_cycle2_test.go:102: roster show --json exit=1: roster show: no roster: declare [roster.<id>] in parley-deck/agents.toml (or keep a legacy §2 table in COOPERATION.md)
--- FAIL: TestValueLayersCannotChangeMembershipState (0.00s)
    roster_cycle2_test.go:257: claude-1 was retired by a value-only layer; the deck file says active
    roster_cycle2_test.go:257: kimi-1 was retired by a value-only layer; the deck file says active
--- FAIL: TestJSONRowHasExactlyTheFrozenColumns (0.94s)
    roster_cycle2_test.go:269: exit=1: roster show: no roster: declare [roster.<id>] in parley-deck/agents.toml (or keep a legacy §2 table in COOPERATION.md)
--- FAIL: TestActiveProvenanceAndMaskingFollowTheAuthority (0.00s)
    roster_cycle2_test.go:316: state provenance = "~/.parley/agents.toml", want the membership authority
    roster_cycle2_test.go:321: write to the authority reported as masked by "~/.parley/agents.toml"
--- FAIL: TestDeckMembershipIsTheDeckFileNotTheLayeredUnion (0.00s)
    roster_membership_test.go:48: membership = 4 agents, want 2 (deck declares claude-1, kimi-1); active=map[claude-1:true codex-1:true hermes-1:true opencode-1:true] inactive=map[]
--- FAIL: TestRosterRenderIsIdempotent (0.00s)
    roster_render_test.go:50: first render exit=1 stderr=roster render: no [roster.*] entries in this deck's config — nothing to render
--- FAIL: TestRosterRenderPreservesLegacyCells (0.00s)
    roster_render_test.go:75: exit=1 stderr=roster render: no [roster.*] entries in this deck's config — nothing to render
--- FAIL: TestRosterRenderTouchesOnlyTheTable (0.00s)
    roster_render_test.go:96: render did not add the new member:
        # Protocol

        ## 2. Active agents (roster)

        Prose explaining that this table is generated.

        | Agent ID | Workspace dir | Role | State |
        | -------- | ------------- | ---- | ----- |
        | `claude-1` | ../claude/ | facilitator | active |

        ## 3. Next section

        Body.
--- FAIL: TestRosterRenderOrderingIsDeterministic (0.00s)
    roster_render_test.go:113: want active alpha-1 < zulu-1 < inactive retired-1, got:
        | Agent ID | Workspace dir | Role | State |
        | -------- | ------------- | ---- | ----- |
        | `retired-1` | – | participant | inactive |
--- FAIL: TestDefaultRosterParticipants (0.00s)
    roster_test.go:76: case1 (inactive-filter): ids=[] had=false err=<nil>
FAIL
FAIL    parley-deck-cli/internal/app       44.079s
ok      parley-deck-cli/internal/config    0.584s
ok      parley-deck-cli/internal/consensus (cached)
ok      parley-deck-cli/internal/driver    (cached)
ok      parley-deck-cli/internal/fsutil    (cached)
ok      parley-deck-cli/internal/hitl      (cached)
ok      parley-deck-cli/internal/loop      (cached)
ok      parley-deck-cli/internal/pipeline  (cached)
ok      parley-deck-cli/internal/procctl   (cached)
ok      parley-deck-cli/internal/protocol  (cached)
ok      parley-deck-cli/internal/protocolcore      (cached)
ok      parley-deck-cli/internal/repomap   (cached)
ok      parley-deck-cli/internal/retro     (cached)
ok      parley-deck-cli/internal/runaction (cached)
ok      parley-deck-cli/internal/runcontrol        (cached)
ok      parley-deck-cli/internal/runmanifest       (cached)
ok      parley-deck-cli/internal/runner    (cached)
ok      parley-deck-cli/internal/runplan   (cached)
ok      parley-deck-cli/internal/runstate  (cached)
ok      parley-deck-cli/internal/sessionstore      (cached)
ok      parley-deck-cli/internal/steer     (cached)
ok      parley-deck-cli/internal/store     (cached)
ok      parley-deck-cli/internal/track     (cached)
ok      parley-deck-cli/internal/tui       (cached)
FAIL
```

[PRIMARY] All 13 failures are in `internal/app`; every other package stays green. Here is the
classification of each failure:

| Failing test | Classification under Path C |
| --- | --- |
| `TestZcodeRowReportsModelAndEffortReadFromAgentConfig` | Old-rule fixture coupling. Its unrelated config-read test creates membership through a values-only zcode block. The assertion remains valid after the fixture declares `members`. |
| `TestDeckDeclaredModelNeverOverridesAgentConfigForUnbindableAdapter` | Same old-rule fixture coupling; the model-binding assertion remains valid. |
| `TestUnreadableAgentConfigFallsBackToUnknown` | Same old-rule fixture coupling; the honest-unknown assertion remains valid. |
| `TestJSONStatusMatchesTextForAHealthyRow` | Old-rule fixture coupling. It has no parent membership and assumes a model/value block creates a member. The frozen JSON/text contract is not contradicted. |
| `TestValueLayersCannotChangeMembershipState` | Direct old authority-rule assertion: adapter-only deck blocks are treated as the membership authority. Its still-valid requirement—that uncommitted layers cannot silently change quorum state—must be restated against explicit `members` plus layered `active`. |
| `TestJSONRowHasExactlyTheFrozenColumns` | Old-rule fixture coupling; the eleven-column assertion remains valid. |
| `TestActiveProvenanceAndMaskingFollowTheAuthority` | Direct old-rule assertion that an adapter-only deck block owns membership/state provenance. Path C needs a new provenance expectation. |
| `TestDeckMembershipIsTheDeckFileNotTheLayeredUnion` | The exact old rule Path C rejects. Its two-member expectation is intentionally obsolete unless the child explicitly sets `members = ["claude-1", "kimi-1"]`. |
| `TestRosterRenderIsIdempotent` | Genuine incompleteness of the minimal prototype: the production renderer needs to read and render the new explicit `members` property. |
| `TestRosterRenderPreservesLegacyCells` | Genuine prototype incompleteness at the membership/render boundary; legacy-cell preservation remains required. |
| `TestRosterRenderTouchesOnlyTheTable` | Genuine prototype incompleteness; a values block no longer adds a member, and no explicit membership parser was added. |
| `TestRosterRenderOrderingIsDeterministic` | Genuine defect in the content heuristic: only the block with `active = false` is classified as membership, so two adapter-only IDs disappear. The fleet census shows the same mixed-content failure. |
| `TestDefaultRosterParticipants` | Genuine defect in the content heuristic/missing `members` implementation: one explicit inactive state plus two adapter-only blocks becomes an empty active roster. |

The old-rule failures do not argue against Path C. The five genuine prototype gaps show why the
resolver-only patch is an experiment, not a production patch: a real implementation needs the
versioned membership property, renderer, participant selection, provenance, and migrated tests in
one contract change.

### E1(c): complete current fleet census

[PRIMARY] I defined a valid deck as a root containing `parley-deck/COOPERATION.md`, then used
`--hidden --no-ignore` so ignored project roots were not skipped. The current volume contains 38
such decks; 37 have at least one `[roster.*]` block:

```text
$ rg --files --hidden --no-ignore '/Volumes/My Shared Files/AI_WORKSPACE' \
    -g 'COOPERATION.md' 2>/dev/null \
    | rg '/parley-deck/COOPERATION\.md$' | sort | wc -l
      38

$ rg -l --hidden --no-ignore '^\[roster\.' '/Volumes/My Shared Files/AI_WORKSPACE' \
    -g 'agents.toml' 2>/dev/null \
    | rg '/parley-deck/agents\.toml$' | sort | wc -l
      37
```

[PRIMARY] The census command ran the installed `/opt/homebrew/bin/parley` and the patched binary
against every one of those 38 roots, parsed the active rows from each JSON result, and compared the
sorted sets. The exact loop was:

```bash
workspace='/Volumes/My Shared Files/AI_WORKSPACE'
shipped='/opt/homebrew/bin/parley'
patched='/tmp/codex-r3-path-c.Psdv91/parley-path-c'
scanned=0; changed=0; errors=0

active_set() {
  local binary="$1" root="$2" payload
  if ! payload=$(env -u PARLEY_HEADLESS_AGENT_CONFIG \
      "$binary" roster show --dir "$root" --json 2>&1); then
    printf 'ERROR:%s' "$(printf '%s' "$payload" | tr '\n' ' ')"
    return 1
  fi
  printf '%s' "$payload" \
    | jq -r '[.roster[] | select(.state == "active") | .agent] | sort | join(",")'
}

while IFS= read -r cooperation; do
  root=${cooperation%/parley-deck/COOPERATION.md}
  rel=${root#"$workspace"/}
  scanned=$((scanned + 1))
  old=$(active_set "$shipped" "$root"); old_rc=$?
  new=$(active_set "$patched" "$root"); new_rc=$?
  verdict=SAME
  if [ "$old_rc" -ne 0 ] || [ "$new_rc" -ne 0 ]; then
    verdict=ERROR; errors=$((errors + 1))
  elif [ "$old" != "$new" ]; then
    verdict=CHANGED; changed=$((changed + 1))
  fi
  printf '%s\t%s\tshipped=[%s]\tpatched=[%s]\n' "$verdict" "$rel" "$old" "$new"
done < <(
  rg --files --hidden --no-ignore "$workspace" -g 'COOPERATION.md' 2>/dev/null \
    | rg '/parley-deck/COOPERATION\.md$' | sort
)
printf 'SUMMARY scanned=%d changed=%d errors=%d\n' "$scanned" "$changed" "$errors"
```

[PRIMARY] Raw output:

```text
CHANGED BYTE shipped=[claude-1,codex-1,hermes-1,kimi-1,opencode-1] patched=[claude-1,codex-1,hermes-1,kimi-1,opencode-1,zcode-1]
CHANGED Finance shipped=[claude-1,codex-1,hermes-1,kimi-1,opencode-1] patched=[]
CHANGED IGBCE shipped=[claude-1,codex-1,hermes-1,kimi-1,opencode-1] patched=[claude-1,codex-1,hermes-1,kimi-1,opencode-1,zcode-1]
CHANGED IHK_PFALZ shipped=[claude-1,codex-1,hermes-1,kimi-1,opencode-1] patched=[claude-1,codex-1,hermes-1,kimi-1,opencode-1,zcode-1]
CHANGED SU-Group-Prompt_library shipped=[claude-1,codex-1,hermes-1,kimi-1,opencode-1] patched=[claude-1,codex-1,hermes-1,kimi-1,opencode-1,zcode-1]
CHANGED adito-outlook-plugin shipped=[claude-1,codex-1,hermes-1,kimi-1,opencode-1] patched=[claude-1,codex-1,hermes-1,kimi-1,opencode-1,zcode-1]
CHANGED aditoLeads shipped=[claude-1,codex-1,hermes-1,kimi-1,opencode-1] patched=[claude-1,codex-1,hermes-1,kimi-1,opencode-1,zcode-1]
CHANGED adito_jvm _issue shipped=[claude-1,codex-1,hermes-1,kimi-1,opencode-1] patched=[claude-1,codex-1,hermes-1,kimi-1,opencode-1,zcode-1]
CHANGED ai_prezz shipped=[claude-1,codex-1,hermes-1,kimi-1,opencode-1] patched=[claude-1,codex-1]
SAME altfins/altfins-marketData shipped=[claude-1,codex-1,hermes-1,kimi-1,opencode-1] patched=[claude-1,codex-1,hermes-1,kimi-1,opencode-1]
CHANGED altfins/altfins-patterns shipped=[claude-1,codex-1,hermes-1,kimi-1,opencode-1] patched=[claude-1,codex-1,hermes-1,kimi-1,opencode-1,zcode-1]
CHANGED altfins shipped=[claude-1,codex-1,hermes-1,kimi-1,opencode-1] patched=[claude-1,codex-1,hermes-1,kimi-1,opencode-1,zcode-1]
CHANGED auftra shipped=[claude-1,codex-1,hermes-1,kimi-1,opencode-1] patched=[]
CHANGED cms shipped=[claude-1,codex-1,hermes-1,kimi-1,opencode-1] patched=[]
CHANGED design-mail/design-mail-fe shipped=[claude-1,codex-1,hermes-1,kimi-1,opencode-1] patched=[claude-1,codex-1,hermes-1,kimi-1,opencode-1,zcode-1]
CHANGED design-mail/design-mail shipped=[claude-1,codex-1,hermes-1,kimi-1,opencode-1] patched=[claude-1,codex-1,hermes-1,kimi-1,opencode-1,zcode-1]
CHANGED design-mail shipped=[claude-1,codex-1,hermes-1,kimi-1,opencode-1] patched=[claude-1,codex-1,hermes-1,kimi-1,opencode-1,zcode-1]
CHANGED ecb-ai-prezz shipped=[claude-1,codex-1,hermes-1,kimi-1,opencode-1] patched=[]
CHANGED ecb-api shipped=[claude-1,codex-1,hermes-1,kimi-1,opencode-1] patched=[claude-1,codex-1,hermes-1,kimi-1]
CHANGED ecb-meeting-2026.05 shipped=[claude-1,codex-1,hermes-1,kimi-1,opencode-1] patched=[]
CHANGED igm-app/igm-app-node shipped=[claude-1,codex-1,hermes-1,kimi-1,opencode-1] patched=[claude-1,codex-1,hermes-1,kimi-1,opencode-1,zcode-1]
CHANGED igm-app shipped=[claude-1,codex-1,hermes-1,kimi-1,opencode-1] patched=[claude-1,codex-1,hermes-1,kimi-1,opencode-1,zcode-1]
CHANGED ldx-wt-mail-fixups shipped=[claude-1,codex-1,hermes-1,kimi-1,opencode-1] patched=[]
CHANGED ldx shipped=[claude-1,codex-1,hermes-1,kimi-1,opencode-1] patched=[]
CHANGED librade-algoTrader shipped=[claude-1,codex-1,hermes-1,kimi-1,opencode-1] patched=[claude-1,codex-1,hermes-1,kimi-1]
CHANGED lustrator shipped=[claude-1,codex-1,hermes-1,kimi-1,opencode-1] patched=[claude-1,codex-1,hermes-1,kimi-1,opencode-1,zcode-1]
SAME millenniumProblems shipped=[claude-1,codex-1,hermes-1,kimi-1,opencode-1,zcode-1] patched=[claude-1,codex-1,hermes-1,kimi-1,opencode-1,zcode-1]
CHANGED paritaetische shipped=[claude-1,codex-1,hermes-1,kimi-1,opencode-1] patched=[claude-1,codex-1,hermes-1,kimi-1,opencode-1,zcode-1]
SAME parley-deck/parley-deck-cli shipped=[claude-1,codex-1,hermes-1,kimi-1,opencode-1,zcode-1] patched=[claude-1,codex-1,hermes-1,kimi-1,opencode-1,zcode-1]
CHANGED parley-deck/wt-editor-composer shipped=[claude-1,codex-1,hermes-1,kimi-1,opencode-1] patched=[claude-1,codex-1,hermes-1]
CHANGED parley-deck/wt-learn-playbooks shipped=[claude-1,codex-1,hermes-1,kimi-1,opencode-1] patched=[claude-1,codex-1,hermes-1]
CHANGED parley-deck/wt-roster-presets shipped=[claude-1,codex-1,hermes-1,kimi-1,opencode-1] patched=[claude-1,codex-1,hermes-1]
CHANGED parley-deck/wt-round-summary shipped=[claude-1,codex-1,hermes-1,kimi-1,opencode-1] patched=[claude-1,codex-1,hermes-1]
CHANGED rev-kimi-scratch shipped=[claude-1,codex-1,hermes-1,kimi-1,opencode-1] patched=[]
CHANGED scaleup shipped=[claude-1,codex-1,hermes-1,kimi-1,opencode-1] patched=[claude-1,codex-1,hermes-1,kimi-1,opencode-1,zcode-1]
CHANGED servers shipped=[claude-1,codex-1,hermes-1,kimi-1,opencode-1] patched=[]
CHANGED test-nextjs shipped=[claude-1,codex-1,hermes-1,kimi-1,opencode-1] patched=[claude-1,codex-1,hermes-1,kimi-1,opencode-1,zcode-1]
CHANGED zeroTrust shipped=[claude-1,codex-1,hermes-1,kimi-1,opencode-1] patched=[claude-1,codex-1,hermes-1,kimi-1,opencode-1,zcode-1]
SUMMARY scanned=38 changed=35 errors=0
```

[PRIMARY] Exact answer to E1(c): **35 decks change their active member set**, and every changed
deck is named above. Nineteen gain `zcode-1`; nine become empty; one becomes a two-member set; two
become four-member sets; and four become three-member sets. Three do not change:
`altfins/altfins-marketData`, `millenniumProblems`, and `parley-deck/parley-deck-cli`. The last is
the sole roster-block-free deck; the other two still undergo a semantic/provenance reinterpretation
even though their active set happens to compare equal. Thus all 37 declared files would change
meaning, while 35 would change active quorum under this particular discriminator.

### E1(d): the rule I would ship

[PRIMARY] The fleet output is a counterexample to block-content inference. Historical `active`
keys are not a reliable witness for an entire `members` list: the heuristic turns real five-member
decks into empty or partial sets. No selection of `adapter`, `active`, or “nonempty block” can
recover the intent that old files never recorded.

I would ship a versioned, presence-aware file contract:

```toml
schema = 2

# Omitted: inherit the parent's complete members property.
# Present: replace the parent's list as one ordinary property.
# members = ["claude-1", "codex-1"]

[roster.kimi-1]
speed = "fast"
```

The rules would be:

1. `schema` absent or `schema = 1` retains the shipped block-presence/legacy semantics byte for
   byte. It is compatibility input only; new files are not written in that form.
2. `schema = 2` activates Path C. The parent is resolved first. Every declared deck property then
   replaces that one property. A `[roster.<id>]` block contains member values only and never
   creates membership.
3. In schema 2, omitted `members` inherits the machine list; present `members = [...]` replaces it.
   `active` is a per-member state property and does not manufacture an ID absent from the effective
   `members` list.
4. V1 has no `super.members +` / `-` operations. Whole-list replacement is the owner's stated
   object-property rule. Add/remove expressions can be a later ergonomic feature if measured use
   justifies them; they are not a second authority.
5. A schema-2 projection of membership into §2 is marked non-authoritative so deleting a TOML
   property cannot promote rendered prose back into membership.

This beats content inspection because it never guesses old intent, beats a standalone
`[membership]` authority because `members` is an ordinary peer property, and beats an unmarked
semantic flip because every old file remains stable until an attended migration. Schema 2 is the
default for new decks and the destination for migrated decks; the marker is a compatibility
boundary, not a competing long-term model.

## E2 result

[PRIMARY] I added one exhaustive audit test only in the temp copy. It reflects over
`agentOverride`, so a newly added TOML field would make the coverage check fail. For each of the
34 current `[agents.*]` keys, the machine parent supplied that property plus an unrelated sentinel,
the deck supplied only a new value for the tested property, and the test checked both the child
winner and the unmentioned parent sentinel.

```text
$ go test -run '^TestPathCE2EveryAgentPropertyLayersIndependently$' -v ./internal/config
=== RUN   TestPathCE2EveryAgentPropertyLayersIndependently
    path_c_e2_test.go:93: COVERAGE 34/34 [agents.*] TOML properties
    path_c_e2_test.go:132: PASS command                    child="child-command"; unmentioned external_backend inherited
    path_c_e2_test.go:132: PASS path                       child="/child/path"; unmentioned external_backend inherited
    path_c_e2_test.go:132: PASS commands                   child="child"; unmentioned external_backend inherited
    path_c_e2_test.go:132: PASS version_args               child="--child"; unmentioned external_backend inherited
    path_c_e2_test.go:132: PASS launch_mode                child="manual"; unmentioned external_backend inherited
    path_c_e2_test.go:132: PASS headless_mode              child="child"; unmentioned external_backend inherited
    path_c_e2_test.go:132: PASS headless_args              child="--child"; unmentioned external_backend inherited
    path_c_e2_test.go:132: PASS acp_args                   child="--child-acp"; unmentioned external_backend inherited
    path_c_e2_test.go:132: PASS interactive_mode           child="child"; unmentioned external_backend inherited
    path_c_e2_test.go:132: PASS interactive_command        child="child-command"; unmentioned external_backend inherited
    path_c_e2_test.go:132: PASS interactive_args           child="--child"; unmentioned external_backend inherited
    path_c_e2_test.go:132: PASS interactive_prompt_mode    child="child"; unmentioned external_backend inherited
    path_c_e2_test.go:132: PASS interactive_invoke         child="child"; unmentioned external_backend inherited
    path_c_e2_test.go:132: PASS interactive_timeout_ms     child="202"; unmentioned external_backend inherited
    path_c_e2_test.go:132: PASS interactive_poll_ms        child="204"; unmentioned external_backend inherited
    path_c_e2_test.go:132: PASS interactive_notes          child="child"; unmentioned external_backend inherited
    path_c_e2_test.go:132: PASS prompt_mode                child="arg"; unmentioned external_backend inherited
    path_c_e2_test.go:132: PASS sandbox_mode               child="child"; unmentioned external_backend inherited
    path_c_e2_test.go:132: PASS approval_policy            child="child"; unmentioned external_backend inherited
    path_c_e2_test.go:132: PASS model                      child="child-model"; unmentioned external_backend inherited
    path_c_e2_test.go:132: PASS model_label                child="child-label"; unmentioned external_backend inherited
    path_c_e2_test.go:132: PASS reasoning                  child="child"; unmentioned external_backend inherited
    path_c_e2_test.go:132: PASS profile                    child="child"; unmentioned external_backend inherited
    path_c_e2_test.go:132: PASS speed                      child="fast"; unmentioned external_backend inherited
    path_c_e2_test.go:132: PASS timeout_ms                 child="206"; unmentioned external_backend inherited
    path_c_e2_test.go:132: PASS first_event_timeout_ms     child="208"; unmentioned external_backend inherited
    path_c_e2_test.go:132: PASS stall_timeout_ms           child="210"; unmentioned external_backend inherited
    path_c_e2_test.go:132: PASS heartbeat_ms               child="212"; unmentioned external_backend inherited
    path_c_e2_test.go:132: PASS isolate_home               child="false"; unmentioned external_backend inherited
    path_c_e2_test.go:132: PASS buffers_stdout             child="false"; unmentioned external_backend inherited
    path_c_e2_test.go:132: PASS isolated_home_env          child="CHILD=yes"; unmentioned external_backend inherited
    path_c_e2_test.go:132: PASS external_backend           child="hosted"; unmentioned notes inherited
    path_c_e2_test.go:132: PASS telemetry                  child="child"; unmentioned external_backend inherited
    path_c_e2_test.go:132: PASS notes                      child="child"; unmentioned external_backend inherited
--- PASS: TestPathCE2EveryAgentPropertyLayersIndependently (0.03s)
PASS
ok      parley-deck-cli/internal/config    0.358s
```

[PRIMARY] Therefore all 34 accepted nonempty/valid child values already layer property by
property across built-in → machine → deck. None of the 34 fails that ordinary case. The complete
set is:

- launch/discovery: `command`, `path`, `commands`, `version_args`, `launch_mode`,
  `headless_mode`, `headless_args`, `acp_args`;
- interactive: `interactive_mode`, `interactive_command`, `interactive_args`,
  `interactive_prompt_mode`, `interactive_invoke`, `interactive_timeout_ms`,
  `interactive_poll_ms`, `interactive_notes`;
- invocation/policy/model: `prompt_mode`, `sandbox_mode`, `approval_policy`, `model`,
  `model_label`, `reasoning`, `profile`, `speed`, `timeout_ms`;
- supervision/isolation/metadata: `first_event_timeout_ms`, `stall_timeout_ms`, `heartbeat_ms`,
  `isolate_home`, `buffers_stdout`, `isolated_home_env`, `external_backend`, `telemetry`, `notes`.

[PRIMARY] Path C does expose one defect beyond membership: “replace a property” is not fully
presence-aware when the replacement is empty or zero. I measured the boundary separately:

```text
$ go test -run '^TestPathCE2EmptyReplacementAudit$' -v ./internal/config
=== RUN   TestPathCE2EmptyReplacementAudit
    path_c_e2_test.go:183: EMPTY-CLEAR commands=[parent-command] version_args=[--parent-version] headless_args=[--parent-headless] acp_args=[] interactive_args=[--parent-interactive] isolated_home_env=map[PARENT:yes] model_label="parent-label" sandbox_mode="parent-sandbox" speed="deep" timeout_ms=901 interactive_timeout_ms=902 interactive_poll_ms=903 notes="parent-notes"
--- PASS: TestPathCE2EmptyReplacementAudit (0.00s)
PASS
ok      parley-deck-cli/internal/config    0.264s
```

`acp_args = []` replaces the parent with an empty list. In contrast, `commands = []`,
`version_args = []`, `headless_args = []`, `interactive_args = []`, and
`isolated_home_env = {}` are treated as absent; empty strings such as `model_label = ""`,
`sandbox_mode = ""`, `speed = ""`, and `notes = ""` are treated as absent; and zero for
`timeout_ms`, `interactive_timeout_ms`, and `interactive_poll_ms` is treated as absent. Explicit
`false` works for both boolean properties, and explicit zero works for the pointer-typed
`first_event_timeout_ms`, `stall_timeout_ms`, and `heartbeat_ms` supervision properties (where it
means disabled). A uniform Path C implementation must document invalid zero values and make every
valid empty/reset value presence-aware rather than silently inheriting the parent.

## Position under Path C

[PRIMARY] My round-02 position was (c): repair the three gestures, then add a separate explicit
overlay with value overrides plus `add` and `remove`. The binding ruling rejects the premise that
membership needs its own authority. I accept that ruling and withdraw the overlay as the v1 model.

[PRIMARY] E1 shows what I now build: schema-2 parent/child resolution, with `members` as an
ordinary optional list property and every `[roster.<id>]` field as an ordinary child override.
Omitted `members` inherits; present `members` replaces. I would not ship `add`/`remove` in v1.
Those are possible future expressions over `super.members`, not a reason to create a second
authority today.

[PRIMARY] Path C does not retire D-A, D-B, or D-C. `set` must preview the resolver's exact
before/after member sets; renderer/parser/§2/drift-guard changes must remain atomic; and `sync`
must say that it rebased fields rather than claiming membership inheritance. The end state changes,
not the need for truthful gestures.

[PRIMARY] I have two plain technical objections to an incomplete Path C implementation:

1. Enabling it on unmarked files is unsafe: this run measured 35 active-set changes, including
   nine empty quorums. The schema boundary is mandatory.
2. Treating empty/zero as “not declared” is not ordinary property replacement. The resolver must
   distinguish absence from an explicit reset wherever reset is valid.

Neither objection conflicts with the direction. Both are requirements for implementing it
faithfully.

## Migration

[PRIMARY] The affected population is the 35 `CHANGED` rows in E1(c). The two declared decks whose
active set compared equal must still be reviewed because equality did not prove intent and their
membership provenance changes. The roster-block-free CLI deck already exercises schema-2-like
inheritance but still needs an explicit marker when the new contract ships.

I would migrate as follows:

1. Ship a dual reader first. Schema absent/1 is frozen compatibility behavior; schema 2 is Path C.
   `parley init` and every new writer emit schema 2. Nothing rewrites an existing file merely
   because the binary was upgraded.
2. Add an attended `roster migrate --schema 2` preview. For each deck it prints the shipped active
   and total member sets, machine-parent sets, candidate schema-2 sets, every joining/leaving ID,
   value overrides retained, legacy-§2 disposition, dirty/non-Git state, and the candidate hash.
   Preview is the default and writes nothing.
3. Require an explicit per-deck choice:
   - **inherit members**: write `schema = 2`, omit `members`, retain only actual value overrides;
   - **preserve current set**: write `schema = 2` plus `members = [...]` equal to the shipped
     complete set, retaining explicit inactive states.
   There is no inferred choice from omissions, block contents, or active-set equality.
4. If candidate and shipped active sets differ, apply additionally requires the existing breaking
   confirmation and repeats the exact diff. Dirty or non-Git decks require backups and remain
   attended; no fleet loop supplies `--yes` on their behalf.
5. Apply atomically, then reload through `LoadRosterScoped`; assert the resolved active and total
   sets equal the accepted preview, assert schema/provenance, render §2 as a marked projection, and
   run the drift guard. A mismatch rolls back from the backup and reports the deck by name.
6. Migrate open ideas only between runs. The current code already freezes a run roster snapshot;
   a machine change or migration affects future ideas, not a locked in-flight quorum.

This lets all 35 affected decks move without a silent quorum change: preserving is exact by
construction, while adopting machine membership is an explicit reviewed breaking decision.

## What I would sign

I would sign a Path C FINAL with all of these conditions:

- `schema = 2` is the versioned compatibility boundary; unmarked/schema-1 files never change
  meaning merely by upgrading the binary.
- In schema 2, `members` is one ordinary replace-on-presence property; omission inherits. Roster
  blocks contain values only. No block-content intent inference remains.
- New decks use schema 2. Existing decks move only through the attended preserve-vs-inherit
  migration above, with exact before/after sets and post-write resolver proof.
- V1 does not add `super.members +/-`, `add`, `remove`, or tombstone syntax. Whole-list replacement
  is sufficient for the ruled model; ergonomics can follow measured demand.
- D-A, D-B, and D-C remain independently releasable prerequisites and their truthful-effect tests
  are green before schema 2 is enabled for new writes.
- Every existing membership/render/participant-selection test is migrated to explicit `members`,
  and `go test ./...` returns green; old-rule assertions are removed only when an equivalent
  schema-2 safety assertion replaces them.
- All 34 `[agents.*]` keys retain their measured per-property layering, and valid empty/reset values
  become presence-aware. `--explain` names the actual source of membership, state, and each value.
- Fleet acceptance records `38 scanned / 37 declared / 35 active-set changes under the rejected
  content heuristic`, and a migration dry-run demonstrates zero unapproved active-set changes.

I would not sign the old (a)/(c) split, an unversioned Path C rollout, the `active`-key heuristic
used only for E1, or a FINAL that calls the current field-only `roster sync` a membership migration.
