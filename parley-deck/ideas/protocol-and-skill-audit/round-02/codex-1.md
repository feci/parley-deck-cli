---
agent: codex-1
idea: protocol-and-skill-audit
round: 2
date: 2026-08-20
responding-to: [codex-1/round-01, zcode-1/round-01, kimi-1/round-01, claude-1/round-01]
---

## Finding verdicts

[PRIMARY] CLI setup for all CLI verdicts below: from the untouched repository copy at
`/tmp/codex-pdaudit-round2.tbjJDC/repo`, I ran
`go build -o /tmp/codex-pdaudit-round2.tbjJDC/parley-src ./cmd/parley` and then
`/tmp/codex-pdaudit-round2.tbjJDC/parley-src --version`; it printed `parley 1.45.0`.
`$BIN` below denotes that source-built binary. Mutating cases used separate subcopies under
`/tmp/codex-pdaudit-round2.tbjJDC/cases/`. For compact exact commands below, `$COPY` denotes
`/tmp/codex-pdaudit-round2.tbjJDC/repo`, `$CASES` denotes
`/tmp/codex-pdaudit-round2.tbjJDC/cases`, `$PH` denotes
`/tmp/codex-pdaudit-round2.tbjJDC/parley-home`, and `$SKILL` denotes
`/tmp/codex-pdaudit-round2.tbjJDC/skill`.

### zcode-1/F1 — PARTIAL
what I ran: [PRIMARY] Under isolated `PARLEY_HOME=$PH`, I ran `$BIN roster set <id> --scope
machine --adapter <adapter> --yes --confirm-breaking --dir $COPY` for
`claude-1:claude`, `codex-1:codex`, `hermes-1:hermes`, `kimi-1:kimi`,
`opencode-1:opencode`, and `zcode-1:zcode`; then `PARLEY_HOME=$PH $BIN roster show --dir $COPY`.
what I got: Six rows, and every row included `inherited-roster`; the copied deck's
`parley-deck/agents.toml` has no `[roster.*]` blocks.
verdict reasoning: The observed inherited state is real, but the claim that it disproves the
authority model is too broad. The same §2 says machine scope writes `~/.parley/agents.toml` and
"every deck inherits it"; an empty deck roster deliberately selects that fallback. A declared
deck roster still controls membership. The defect is narrower documentation drift: §2's opening
sentence is unqualified and its live generated table does not explain this deck's deliberate
inherited state. The CLI behavior is correct; the documentation needs qualification.
corrected severity (if different): MINOR, not MAJOR.
is this a duplicate of another finding? Closely overlaps zcode-1/F2; both arise from this deck's
intentional inherited-roster configuration.

### zcode-1/F2 — REFUTED
what I ran: [PRIMARY] On a fresh deck copy, first
`PARLEY_HOME=$PH $BIN roster render --dir $CASES/zf2`, then
`PARLEY_HOME=$PH $BIN roster render --dir $CASES/zf2 --adopt-inherited --yes`, followed by
`grep -A8 -n '^The generated view:' $CASES/zf2/parley-deck/COOPERATION.md`.
what I got: The unqualified command exited 1 and explained that adopting machine-local rows would
commit them into a shared file. The qualified command printed `Regenerated §2` and the table then
contained all six active agents.
verdict reasoning: The round-01 report says exit 0 and calls the table "permanently empty"; neither
is true for the current source. The refusal is a safety gate with an explicit, functioning escape
hatch. Appendix A also invokes render only *after* declaring deck roster blocks, where no adoption
flag is needed. The owner's decision not to adopt is project policy, not a CLI defect.
corrected severity (if different): No finding.
is this a duplicate of another finding? Same inherited-roster root state as zcode-1/F1.

### zcode-1/F3 — REFUTED
what I ran: [PRIMARY] `sed -n '1,9p' parley-deck/COOPERATION.md`; `sed -n '1,80p'
parley-deck/meta/version.json`; `find parley-deck/meta -name 'protocol-sync_*'`; and
`git log -S'parley-deck-skill 2.9.0' -- parley-deck/COOPERATION.md` in the repository copy.
what I got: The header says 2.9.0, metadata says 2.8.0, and only the 2026-06-13 sync record exists.
Git attributes the 2.9.0 header change to commit `2da38c5` (`zcode-adapter: roster reads
MODEL/EFFORT ...`), not to a §9.0 consumer auto-sync.
verdict reasoning: The state mismatch is reproducible, but the alleged contradiction assumes the
header edit was the §9.0 consumer-sync operation. The cited requirement to create a
`protocol-sync_*` record belongs specifically to that operation. This is the upstream/source repo,
and the history shows an ordinary source commit updated the version label. No run demonstrated a
consumer sync that omitted its record. Stale metadata is addressed more directly by kimi-1/F2.
corrected severity (if different): No finding as stated.
is this a duplicate of another finding? Its metadata symptom overlaps zcode-1/F5 and kimi-1/F2.

### zcode-1/F4 — PARTIAL
what I ran: [PRIMARY] In `$COPY`, `git branch -a --format='%(refname:short)' | rg -c
'^(.*/)?idea/'`; `git log --all --format='%H %s' | rg -c 'FINAL\.md \+ close idea'`; the same
log with `--merges`; then both counts on `git log --first-parent ... main`.
what I got: There is currently one `origin/idea/...` ref, 34 matching close-message commits across
all refs, and 9 matching merge commits. On `main`'s first-parent history there are only 14 matching
close transactions, 9 of them merges: five, not 25, are demonstrably non-merge first-parent closes.
Recent Phase-0 files also appear as single-parent first-parent commits.
verdict reasoning: Some history does violate the declared GitHub-PR/merge-commit mechanics, but
round 1's subtraction across all refs mistakes ordinary commits *inside merged branches* for
direct-to-main closes. Deleted idea branches are required by §11.B, so their absence proves
nothing about past PR use. The narrower five-close mismatch survives; the headline measurement
and "never been used" claim do not.
corrected severity (if different): MINOR, absent a demonstrated lost review or bad merge.
is this a duplicate of another finding? No.

### zcode-1/F5 — REFUTED
what I ran: [PRIMARY] `rg -n 'protocolRole' parley-deck/meta/version.json parley-deck/meta || true`.
what I got: No matches.
verdict reasoning: Absence alone is not a contradiction: §9.0 explicitly defines the
missing/unknown state and says to ask once and backfill. Round 1 did not run `preflight` and show
that fallback failing. The reason the field can disappear is independently testable under
kimi-1/F2; merely observing the anticipated pre-remediation state is not a separate defect.
corrected severity (if different): No finding.
is this a duplicate of another finding? Duplicate symptom of kimi-1/F2; also overlaps zcode-1/F3.

### zcode-1/F6 — PARTIAL
what I ran: [PRIMARY] `find parley-deck/ideas -maxdepth 1` with `*pipeline*`/`*end-to-end*`
filters, `git log --all -- '*end-to-end-pipeline*'`, and
`rg -n 'full automatic idea-to-monitoring pipeline' parley-deck/ideas/*/00-prompt.md`.
what I got: The literal cited slug has no path or historical path. The ratification material is
at `ideas/2026-06-02T12-07-14-meta-protocol-ch/`, whose prompt contains the exact pipeline title.
verdict reasoning: The literal citation is stale, so the documentation locator should be fixed.
The stronger claim that the ratification trail "exists nowhere" is false: it is readily located
by title and date and contains the artifact. This is a legacy naming/citation problem, not a lost
decision or CLI defect.
corrected severity (if different): NIT, not MINOR.
is this a duplicate of another finding? No.

### zcode-1/F7 — PARTIAL
what I ran: [PRIMARY] `rg -n '^## (8|9|10|11|12|13|14|15)\.'
parley-deck/COOPERATION.md` and `sed -n '34,41p' parley-deck/COOPERATION.md`.
what I got: §10 precedes §9; the progressive-disclosure list names §9 and §11–§14 but not §10 or
§15.
verdict reasoning: Omitting §15 from the map is a real navigation defect because §15 binds on
every track. Omitting §10 is not: it is a redundant TL;DR rather than required reference material,
and numeric placement has no semantic effect. Moreover, the core Phase 1/2/3/6 text explicitly
points readers to §15, limiting the consequence.
corrected severity (if different): NIT, not MINOR.
is this a duplicate of another finding? No.

### zcode-1/F8 — PARTIAL
what I ran: [PRIMARY] `find parley-deck -maxdepth 1 -mindepth 1 -print | sort` and compared the
result with the §3 tree.
what I got: Disk has `COOPERATION.md`, `agents.toml`, `ideas/`, `inbox/`, `meta/`, and `runs/`;
the illustrative tree omits `agents.toml` and `runs/`.
verdict reasoning: The tree is stale and should show both entries. Nothing in the command shows a
runtime or audit failure, and both paths are documented where they matter (§2 and §12.12), so the
claimed consequence is overstated.
corrected severity (if different): NIT, not MINOR.
is this a duplicate of another finding? No.

### zcode-1/F9 — REFUTED
what I ran: [PRIMARY] `sed -n '586,609p' parley-deck/COOPERATION.md`.
what I got: The LE-4 closing paragraph appears between the template-introducing colon and the
indented `## Fix-up cycle N` template.
verdict reasoning: This is inelegant placement, but the paragraph and template remain complete and
unambiguous. No command failure, contradictory instruction, or unreachable state was reproduced.
The audit prompt expressly excludes "this could be confusing" findings.
corrected severity (if different): No finding.
is this a duplicate of another finding? No.

### zcode-1/F10 — REFUTED
what I ran: [PRIMARY] `$BIN learn --help`; `$BIN learn automation-outer-loop --dir
$CASES/zf10/parley-deck`; the same command with `--dir $CASES/zf10`; status inspection of
that idea's prompt, FINAL, and IMPLEMENTATION; and `$BIN learn
2026-06-02T12-07-14-meta-protocol-ch --dir $CASES/zf10`.
what I got: Help explicitly says `--dir` is the **workspace root**. Passing the deck directory
failed with the doubled path; passing the documented root wrote one advisory playbook. Although
the prompt is stale at `round-01`, `FINAL.md` is `final` and `IMPLEMENTATION.md` is `complete`.
The uppercase timestamp name was rejected as not lowercase kebab-case.
verdict reasoning: All three alleged deviations dissolve. The reporter passed a deck path against
the command's explicit root contract; the chosen idea is completed despite stale kickoff metadata;
and the legacy directory violates the protocol's current slug grammar. The stale kickoff status is
a record-quality issue under zcode-1/F14, not evidence that `learn` accepts unfinished work.
corrected severity (if different): No finding.
is this a duplicate of another finding? Subclaim (b) duplicates/misattributes zcode-1/F14.

### zcode-1/F11 — PARTIAL
what I ran: [PRIMARY] `sed -n '977,982p;1009,1016p' parley-deck/COOPERATION.md`.
what I got: Phase 5 says to commit `IMPLEMENTATION.md` directly; the later recommended branch
protection says to require PRs for all `ideas/` changes.
verdict reasoning: There is a conditional documentation gap if a project adopts the recommended
protection: Phase 5 needs the same small-PR escape already stated for Phase 6. The branch-protection
block is explicitly recommended, not a universal invariant, and no protected push was attempted,
so this is not a demonstrated blocked implementation.
corrected severity (if different): NIT, not MINOR.
is this a duplicate of another finding? No.

### zcode-1/F12 — REFUTED
what I ran: [PRIMARY] `sed -n '3,9p;1087,1094p' parley-deck/COOPERATION.md`.
what I got: The live header has six fields; Appendix A asks adopters to fill in two labels not
present there verbatim.
verdict reasoning: Markdown frontmatter-like header fields are not a closed schema. "Fill in the
header" can instruct an adopter to add project-specific fields, and `Parley deck:` already carries
the relevant shared path. No parser rejection or contradictory field contract was shown.
corrected severity (if different): No finding.
is this a duplicate of another finding? No.

### zcode-1/F13 — REFUTED
what I ran: [PRIMARY] `rg -n 'excluded|readiness|preflight|liveness|available|unavailable'
parley-deck/ideas/{agents-verify-hermes-probe,preflight-liveness-false-negative,protocol-and-skill-audit}/00-prompt.md`.
what I got: No kickoff readiness record; matches are only prose describing audit/preflight topics.
verdict reasoning: This reproduces noncompliant artifacts created by their initiators, not a
protocol or CLI contradiction. The rule tells the facilitator to run and record the check; round 1
did not show a source-built `preflight` losing a result after it was run. Agent-created bad state
must not be converted into a product defect.
corrected severity (if different): No finding.
is this a duplicate of another finding? No; it is process noncompliance rather than a software
finding.

### zcode-1/F14 — PARTIAL
what I ran: [PRIMARY] I enumerated `ideas/*/FINAL.md`, parsed the first `status:` from each existing
`00-prompt.md`, and printed non-final/complete/abandoned values.
what I got: 78 directories have `FINAL.md`; 77 have a prompt. Of those, 20 prompts are stale, and
`launch-orphan-hardening` has no prompt, giving the reporter's 21 anomalous records only if that
missing file is counted as an empty status.
verdict reasoning: The stale-record measurement survives and is worth cleanup. The stated
consequence does not: no command reproduced an agent actually writing into a closed round, and
other commands such as `learn` correctly use `FINAL.md`/`IMPLEMENTATION.md` completion evidence.
This is historical protocol noncompliance, not proof of a CLI guard failure.
corrected severity (if different): NIT, not MINOR.
is this a duplicate of another finding? Its `automation-outer-loop` instance is the stale-state
subclaim incorrectly counted as a `learn` defect in zcode-1/F10(b).

### zcode-1/F15 — CONFIRMED
what I ran: [PRIMARY] `$BIN --help | rg -i 'learn|preset' || true`, then `$BIN learn --help` and
`$BIN preset list --dir /tmp/codex-pdaudit-round2.tbjJDC/repo`.
what I got: Top-level help produced no `learn`/`preset` match. Both subcommands parsed and ran;
`preset list` printed `No roster presets defined...` and exited 0.
verdict reasoning: The source-built CLI has functional command routes that its primary discovery
surface omits. The finding is narrow and correctly classified.
corrected severity (if different): NIT (unchanged).
is this a duplicate of another finding? No.

### kimi-1/F1 — CONFIRMED
what I ran: [PRIMARY] From the copied 2.9.0 package I ran
`KIMI_CODE_HOME=/tmp/codex-pdaudit-round2.tbjJDC/kimi-code-home node
bin/parley-deck-skill.js install --target kimi`, changed only `references/COOPERATION.md` in the
installed managed core (`Multi-Agent` → `Multi-agent`) with `apply_patch`, then ran `doctor
--target kimi` under the same `KIMI_CODE_HOME`.
what I got: Install reported all six units installed. After the core protocol edit, doctor still
reported `kimi/parley-deck: valid` and all five add-ons valid, exit 0.
verdict reasoning: A current, managed core protocol can change bytes without the advertised health
gate detecting it. The test targets the highest-value installed protocol file and is not explained
by a stale binary or foreign-install semantics. The behavior, not merely the README, should be
fixed because add-ons and foreign cores already establish byte verification as the intended guard.
corrected severity (if different): MAJOR (unchanged): a corrupted installed protocol is trusted.
is this a duplicate of another finding? kimi-1/F5 is a subordinate manifestation of the same
managed-core manifest omission and should not become a separate fix.

### kimi-1/F2 — CONFIRMED
what I ran: [PRIMARY] In a copied project I added `"protocolRole": "source"`, ran
`node $SKILL/bin/parley-deck-skill.js sync-project --project $CASES/kf2 --yes`, and inspected the
file. I then ran source-built `$BIN preflight --dir $CASES/kf2 --no-ping --json`, followed by
the same command with `--yes`.
what I got: `sync-project` printed `wrote $CASES/kf2/parley-deck/meta/version.json`; the role disappeared while
`deckVersion` became 2.9.0. Preflight then exited 3 with classification `unknown-role`. With
`--yes`, it exited 0, reported `protocolRole was absent; backfilled protocolRole=consumer`, and
wrote `consumer` into what had been a source-role project.
verdict reasoning: This is direct cross-package semantic corruption with a demonstrated
consequence. The metadata refresher must preserve the CLI-owned role (or share a schema), not
silently erase it. Merely changing documentation would leave the source/consumer safety gate
unstable.
corrected severity (if different): MAJOR (unchanged).
is this a duplicate of another finding? It explains the symptoms in zcode-1/F3 and F5. kimi-1/F4
is another message-level manifestation of the same missing role awareness.

### kimi-1/F3 — CONFIRMED
what I ran: [PRIMARY] `node -e` over `require('./lib/installer').TARGETS`, `paths --target all
--include-undetected | wc -l`, `rg -n 'fourteen|--target auto' README.md`, and the generated CLI
help.
what I got: The implementation lists 15 targets including `zcode`; paths printed 90 lines (15 ×
6). README has four `fourteen` claims and its target list stops at `aionrs`; generated help includes
`zcode`.
verdict reasoning: This is reproducible README drift. The installer is correct; the fix direction
is documentation only—change the count and enumerate `zcode`.
corrected severity (if different): MINOR (unchanged).
is this a duplicate of another finding? No.

### kimi-1/F4 — PARTIAL
what I ran: [PRIMARY] In a fresh project copy whose metadata explicitly contained
`"protocolRole": "source"`, I ran `node $SKILL/bin/parley-deck-skill.js status --project
$CASES/kf4 --target kimi`.
what I got: Status called metadata valid but printed `Review the local COOPERATION.md changes
before adopting packaged protocol updates`, plus a `sync-project --yes` action.
verdict reasoning: The role-insensitive wording is reproduced and should be made source-aware.
The consequence is overstated: `status` remains advisory, does not write `COOPERATION.md`, and
tells the operator to review rather than blindly adopt. More importantly, this is not an
independent defect: it is the same missing `protocolRole` ownership/read path as kimi-1/F2.
corrected severity (if different): NIT, not MINOR; fold its message correction into F2.
is this a duplicate of another finding? Yes—kimi-1/F2.

### kimi-1/F5 — REFUTED
what I ran: [PRIMARY] I inspected the markers written by the sandbox install and
`sed -n '1,24p' lib/installer.js`.
what I got: The core marker has `markerSchema: 2`, `addon: false`, and no `manifest`; the add-on
marker has `addon: true` and a manifest. The comment says a schema-2 marker omitting `manifest` is
malformed.
verdict reasoning: The comment is under-scoped, but round 1 explicitly concedes there is no runtime
effect. Core and add-on markers are discriminated by `addon`; the only material defect is that the
managed core is not integrity-checked, already reproduced as kimi-1/F1. A comment wording nit must
not be counted as a second defect for the same omission.
corrected severity (if different): No independent finding.
is this a duplicate of another finding? Yes—fully subsumed by kimi-1/F1.

### claude-1/F1 — REFUTED
what I ran: [PRIMARY] The reported loop of case-insensitive counts over the packaged
`skills/parley-deck/SKILL.md`, plus `rg` for its protocol-loading instructions.
what I got: All six command-name counts are zero. The same skill says `Always read
parley-deck/COOPERATION.md first`, has a `Required Protocol Context` gate, and declares the live
protocol canonical.
verdict reasoning: The counts are true but do not establish the claim. The skill intentionally
does not duplicate the full command catalog; it requires loading the document that names those
commands. Its readiness instructions can also be performed manually, so absence of the optional
`preflight` shortcut does not make the obligation impossible to discharge. Treating the skill as
standalone contradicts its explicit loading contract.
corrected severity (if different): No finding.
is this a duplicate of another finding? No.

### claude-1/F2 — CONFIRMED
what I ran: [PRIMARY] In a fresh full repository copy I ran source-built `$BIN roster set codex-1
--scope deck --adapter codex --yes --confirm-breaking`, then `$BIN roster render --yes`, then
`go test ./internal/protocol/...`. I separately repeated the test after inherited-roster adoption.
what I got: Render printed `Regenerated §2`; the test failed:
`live deck: anchor "| Agent ID       | Workspace dir                       | Role          |"
appears 0 times, want exactly 1 (drift guard fails closed)`. This occurs with an explicitly
declared deck roster too, not only with `--adopt-inherited`. Control: the identical test in the
untouched repository copy passed (`ok parley-deck-cli/internal/protocol`).
verdict reasoning: The current source command mandated after normal roster declaration emits the
four-column generated table, while the current source test requires the old padded three-column
anchor. The ordinary documented update path makes this repository's guard red. This is a genuine
renderer/guard contract defect.
corrected severity (if different): MAJOR (unchanged).
is this a duplicate of another finding? It shares the roster-render surface with zcode-1/F2, but
that finding concerns inherited refusal/empty display; this is a distinct test-breaking output
contract.

### claude-1/F3 — CONFIRMED
what I ran: [PRIMARY] In a copied deck I supplied a higher-precedence
`PARLEY_HEADLESS_AGENT_CONFIG` with `[roster.codex-1] model = "env-model"`, then source-built
`roster set codex-1 --scope deck --model new-deck-model --yes`, `roster show`, and `roster show
--explain codex-1`.
what I got: `roster set` warned that the value was `MASKED` and explicitly said `(status
masked-by-env; see roster show --explain codex-1)`. Both show forms reported only
`effort-unknown,metadata-unknown`; neither emitted `masked-by-env`, while explain correctly named
the environment file as the effective model source.
verdict reasoning: This is stronger than the round-01 source grep: the promised status is absent
in a real masked write. Either emit it in the frozen STATUS column or remove/reclassify it in both
the vocabulary and warning.
corrected severity (if different): MINOR (unchanged).
is this a duplicate of another finding? No.

## Findings I could not assess, and why

None of the 23 permitted findings was left unassessed. I did not inspect
`round-01/codex-1.md`, as instructed. I also did not mutate a GitHub/GitLab remote; where a claim
depended on optional branch protection rather than reproducible local behavior (zcode-1/F11), I
recorded only the narrower documentation tension and did not confirm the hypothetical push failure.

## Duplicates I found across authors

- `kimi-1/F2` is the material `protocolRole` defect. `zcode-1/F5` is its absent-field symptom,
  `zcode-1/F3` partly observes the resulting stale metadata, and `kimi-1/F4` is the same missing
  role-awareness expressed through generic status advice. They should produce one primary fix,
  with message/metadata tests, not four defect tickets.
- `kimi-1/F5` is wholly subsumed by `kimi-1/F1`: the schema-2 core-marker comment and the absent
  managed-core integrity check are the same manifest omission. Only F1 has a runtime consequence.
- `zcode-1/F10(b)` is not a `learn` defect; it is one instance of the stale kickoff statuses
  measured in `zcode-1/F14`.
- `zcode-1/F1` and `zcode-1/F2` share the inherited-roster configuration and should not both be
  described as roster-authority failures. F1 leaves a narrow docs qualification; F2 is refuted.
- `zcode-1/F2` and `claude-1/F2` touch `roster render` but are **not** duplicates: the former
  alleges an inherited-roster refusal/permanently empty view; the latter reproduces a distinct
  renderer-versus-drift-guard test failure.

## Anything round 1 missed that these findings led me to

### Additional verification — standalone `preflight` ignores the canonical roster

[PRIMARY] In a source-role deck copy I declared exactly two active members with source-built CLI:
`codex-1` (`adapter=codex`, installed) and `qwen-1` (`adapter=qwen`, not installed). I then ran
`$BIN roster show --dir $CASES/new-preflight` followed by
`$BIN preflight --dir $CASES/new-preflight --no-ping`.

The relevant `roster show` rows were:

```text
codex-1  codex  active  yes  gpt-5.6-sol  GPT  OpenAI  max      deep  yes  ok
qwen-1   qwen   active  no   unknown      unknown unknown unknown  deep  no   not-installed,model-unbound,effort-unknown,metadata-unknown
```

But standalone preflight discarded both stable roster IDs, ignored unavailable `qwen-1`, expanded
to seven installed adapter families (`codex`, `claude`, `agy`, `hermes`, `kimi`, `opencode`,
`zcode`), and ended:

```text
Ready: no pending gates.
```

This contradicts §9.0's requirement to probe every rostered participant and can certify readiness
while an actual quorum member is unavailable. Corrected scope: this is the standalone `preflight`
command; the source-built `parley run` path separately maps its selected participant set before its
shared precheck. Suggested severity: **MAJOR** because the documented facilitator command gives a
false readiness result against the wrong quorum. This issue was absent from the 23 findings I was
allowed to inspect; I did not compare it with the prohibited `codex-1/round-01` file.
