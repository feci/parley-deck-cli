---
idea: skill-sync-cli-1-39
drafter: claude-1
participants: [claude-1, codex-1, hermes-1, kimi-1]
track: standard
rounds: 2
date: 2026-08-06
status: accepted
---

## What was decided

Six edits to `parley-deck-skill`. The brief asked about three; the run found **four defects in one
forty-line region of `SKILL.md`**, and three of the four were found by participants rather than by
the facilitator's brief.

The framing that settled the scope, proposed by codex-1 and applied by everyone including against
their own round-1 positions: **the skill is instruction text, so the test is not "is this true" but
"does a facilitator or participant behave differently because of it".**

### 1. `opencode` row in the Autonomous Execution table

Unanimous. Adopted in kimi-1's wording, which carries the argv detail the other versions omit:

> | opencode | `run --auto` — the prompt is an argv positional, not stdin. `opencode run` writes
> unattended even without `--auto`; pass `--auto` explicitly, because an implicit vendor default is
> what may change between versions. |

### 2. `SKILL.md:251` — the inverted source-of-truth sentence

**kimi-1's find, and the sharpest of the run.** The line *"The source of truth for each agent's mode
is the spec's `autonomous_write` field"* points readers at the field 1.39.0 established does **not**
launch anything. It sits three lines under the table a facilitator uses to pick a flag, so it
actively directs attention away from the value that decides whether the agent can write.

All three non-facilitator participants converged on one replacement, each building on the others.
Adopted text — kimi-1's, which is the superset (its lead sentence is the inversion fix, its body is
codex-1's contract framing, its mechanism clause is hermes-1's, and the existing confinement floor
is retained verbatim at the end):

> The source of truth for an agent's autonomous capability is the **effective launch argv**, not the
> declared mode. The declared autonomous-write mode is a verification contract, not a second set of
> launch arguments: before treating a headless participant as able to write its artifact, inspect
> the effective launch arguments after all configuration layers have been applied — the launch
> config recorded in the orchestration summary, or `parley agents list` when the parley CLI drives
> the agents — and verify that every argument required by the declared mode is present. A config
> override can replace the launch arguments wholesale and silently drop the enabling flag. If the
> effective arguments cannot be inspected, or any required argument is absent, treat autonomous
> write as unavailable (`AUTO=no`) and do not launch that participant as write-capable. Passing this
> check proves only that the autonomous mode is enabled; it does not prove workspace confinement. If
> workspace confinement cannot be demonstrated for an agent, treat its autonomous bit as unset
> (fail-closed) rather than escalating to a full-filesystem bypass.

This replaces **both** sentences at `:251` — the inverted one and the confinement-only fail-closed
one — rather than adding a third paragraph beside them.

### 3. The command-construction recipe — manual vs CLI branches

**codex-1's find.** Steps 4-5 tell a facilitator to add model/thinking flags and deliver the prompt
as separate assembly stages. The launcher does neither: it substitutes `{prompt}` and `{root}`
*inside* `HeadlessArgs` and launches that vector alone (`runner.go:1094-1108`). The recipe is a
legitimate description of hand-rolling an invocation and a false description of what `parley` does
with config, and the section currently blurs them into one numbered list.

Adopted: split "Generic CLI Invocation Contract" into an explicit **manual facilitator** branch and
a **Parley CLI** branch. The multi-step assembly stays in the manual branch only. The CLI branch
states that resolved `headless_args` is the complete argv template, `{prompt}` must already sit in
its required position, `prompt_mode=stdin` controls only stdin wiring, and **no** permission, model,
thinking, profile or prompt arguments are synthesized afterwards. The same boundary applies to
"Headless Agent Configuration" and `WORKED_EXAMPLES.md`.

### 4. `compatibility.json` — bump plus one assertion

**Unanimous only after hermes-1 withdrew its round-1 objection, and it withdrew on evidence rather
than on the count.** It checked whether the existing `parley-addon.json` sha256 pin already guards
this, found that the pin verifies file hashes and never reads `skillVersion` or compares it to
`package.json`, and concluded its "no new tooling" objection did not apply.

Adopted: set `skillVersion` to **whatever version actually ships from this change** (not a number
chosen now), and add exactly one equality assertion — `compatibility.skillVersion ===
package.version` — to the existing Node test harness. No new script, job, or checker.

**Revision 2 correction — the `prepack` clause was false as drafted.** Revision 1 said the assertion
would run "under both `npm test` and the existing `prepack` lifecycle". It would not.
`package.json:60` defines `test` as `node --test && …`; `package.json:66` defines `prepack` as
`node scripts/build-addon-manifest.js --check` alone. **`prepack` does not invoke `node --test`**, so
an assertion living in a test file runs on `npm test` only, and a direct `npm publish` bypasses it —
precisely the path that let `skillVersion` drift four releases. Verified `PRIMARY` by the drafter
against `package.json:60,66`; first raised by hermes-1 (`round-02/hermes-1.md`), restated as
reservation R1 by kimi-1 in signoff.

**Implementation MUST close the `prepack` half explicitly**, by one of:
(a) extend the `prepack` command to run the version check, or
(b) put the equality check inside `scripts/build-addon-manifest.js`, which `prepack` already runs.
Both edit existing files, so "no new script, job, or checker" survives either way. Shipping only the
test-file assertion does **not** satisfy this decision.

Recorded cost, raised by claude-1 and not disputed: once the assertion lands, bumping `package.json`
without `compatibility.json` fails the suite and the pack step. That is the intent, and it becomes a
release step rather than a surprise.

hermes-1 additionally flagged an implementation gap: the existing manifest test verifies hashes and
does not read either version field, so the assertion needs its own home rather than an extension of
that check. Confirmed by the drafter (`PRIMARY`): nothing under `test/` or `scripts/` reads
`compatibility.json`'s `skillVersion` today.

### 5. `writeModeArgs` — deleted from the documented JSON shape, plus a migration rule

Carried as VC-1 at 3-to-1 in revision 1; closed in revision 2 on the design argument, with the
workspace counts explicitly demoted to illustration — see `## Verdict conflicts`. The field is
deleted, the manual branch states that the write-enabling flag belongs inside `headless_args`, and
the branch carries a **legacy migration rule** for configs that already set the field.

### 6. What stays out

Unanimous. Promotion history, probe versions and outputs, the `hermes` incident narrative, ACP
availability, kimi's exit-1 message, and every CLI struct name. None changes a next action.

`references/COOPERATION.md` is **untouched**: 1.39.0 changed no protocol text, verified by
`git diff v1.38.0..v1.39.0` over both copies and by a normalized diff against the live deck.

## Verdict conflicts

### VC-1 — Does `writeModeArgs` stay in the documented JSON shape? CLOSED in revision 2 — deleted, with a migration rule

Everyone agrees the CLI-facing contract must show no such field. The residue is the **manual**
facilitator JSON example.

- **hermes-1, kimi-1, claude-1 — delete it.** hermes-1: *"A facilitator copying the JSON example
  will still populate `writeModeArgs` and still believe it matters. Remove it from the JSON shape;
  add the invariant line. Both, not one."* kimi-1 **withdrew its round-1 "keep the split"** on the
  grounds that the field split itself teaches the launch model the CLI does not implement.
- **codex-1 — keep it, but only inside the explicitly manual branch**, on the grounds that once the
  manual and CLI branches are separated the field is no longer a claim about CLI config, and that
  keeping it avoids a schema bump for local configs that already set it.

**This was 3-to-1, and §15.3 forbids resolving it that way, so revision 1 carried it rather than
counted it.** The unexamined fact on which it turned: **nobody had checked whether any local config
in the wild actually sets `writeModeArgs`.** claude-1 raised that in round 2 and it went unanswered.

Revision 1 then attempted to close it on a drafter-owned measurement and was blocked. Revision 2
closes it on the narrower ground codex-1 specified — the design argument, with the counts demoted
to illustration — and adopts codex-1's migration condition. The dissent stands recorded.

### VC-1 measurement — revision 2

Revision 1 of this section was blocked by codex-1 on three grounds, all upheld: claim (b) carried a
`PRIMARY` label without the command, inputs, output or locators §15.2 requires; the claim-(a)
inferences overreached; and deleting a field that live configs populate needs a migration rule.
Revision 2 replaces the section. **Revision 1's numbers were also wrong** — see the tooling defect
below. Both are recorded rather than quietly corrected.

#### (a) There is no Go loader and no Go field — `CONFIRMED`

**Scoped claim (narrowed per codex-1):** no Go code in `parley-deck-cli` declares a
`writeModeArgs`/`write_mode_args` field or reads `meta/headless-agents.local.json`.

Command, run from the repository root, covering **all 1663** `*.go` files including tests:

```
find . -name "*.go" -type f -print0 | xargs -0 grep -ln \
  "headless-agents\|headless_agents\|writeModeArgs\|write_mode_args\|WriteModeArgs"
```

Output: empty, `rc=1`. Locators: `agentOverride` (`internal/config/runtime.go:97-132`) has
`HeadlessArgs` and no write-mode field; `configLayers` (`internal/config/runtime.go:134-154`)
enumerates `~/.parley/agents.toml`, `parley-deck/agents.toml`, `parley-deck/agents.local.toml`,
`$PARLEY_HEADLESS_AGENT_CONFIG` — no JSON path.

Non-owner verdicts: **codex-1 `CONFIRMED PRIMARY`** (independent `rg -n -g '*.go'`, exit 1, no
output) and **hermes-1 `CONFIRMED PRIMARY`** (independent search plus its own read of
`configLayers:137-154`). kimi-1 confirmed the `writeModeArgs` half and correctly withheld a verdict
on the file-reading half, which it co-owns.

**Scope caveat on those confirmations (drafter, `PRIMARY`).** `rg` honours `.gitignore` by default,
so the non-owner runs covered the **172 tracked** `*.go` files — `internal/` (171) and `cmd/` (1) —
and not the 1491 files under the ignored `.gomodcache`. The drafter's `find`-based run covers all
1663 and also returns zero. The tracked-file scope is the one that matters for this claim, since
`.gomodcache` is a vendored dependency cache and cannot define the CLI's config loader; the
difference is recorded so the confirmations are not read as broader than they were.

**What (a) does NOT establish, per codex-1 and accepted:** it does not show the field has "no
consumer", that deleting it is a "zero behaviour change", or that "the protected cost does not
exist". Those three inferences are **withdrawn from revision 1**. The file *does* have a consumer —
a facilitator hand-assembling a command, exactly as codex-1's position held. (a) establishes only
that the consumer is not the Go launcher.

#### (b) The field is populated in live deck configs — `CONFIRMED`, corrected numbers

Commands, run from `/Volumes/My Shared Files/AI_WORKSPACE`:

```
find "/Volumes/My Shared Files/AI_WORKSPACE" -name headless-agents.local.json \
  -path "*/parley-deck/meta/*" -print0 | sort -z            # -> 23 files
find "/Volumes/My Shared Files/AI_WORKSPACE" -name headless-agents.local.json \
  -path "*/parley-deck/meta/*" -print0 \
  | xargs -0 grep -l writeModeArgs | sort                   # -> 12 files
```
plus a JSON pass over the same 23 files comparing each agent's `writeModeArgs` against its
`headlessArgs`.

**R2 correction (kimi-1, upheld).** Revision 2 first printed this pipeline with a bare `|` and
plain `xargs`. On this volume that command is **broken**: the workspace path contains spaces, and
`xargs` without `-0` splits on blanks — kimi-1 ran the printed shape verbatim and got **0** files,
not 12 (`PRIMARY`). The numbers above were produced by `-print0`/`xargs -0` and by a `while read`
loop, both of which are space-safe; the *transcription* was wrong, not the measurement. Corrected
here, because a published command that cannot reproduce the published number is the same failure
this section exists to warn about.

| measure | count |
|---|---|
| deck configs found | 23 |
| contain the key `writeModeArgs` | 12 |
| set it **non-empty** on ≥1 agent | 10 |
| have ≥1 agent whose `writeModeArgs` values are **absent** from `headlessArgs` | 9 |
| declare it **empty** (flags already migrated into `headlessArgs`) | 2 — `BYTE`, `ecb-meeting-2026.05` |

The 10: `aditoLeads`, `lustrator`, `design-mail`, `design-mail/design-mail`,
`design-mail/design-mail-fe`, `igm-app`, `zeroTrust`, `millenniumProblems`, `paritaetische`, and
this repo's own deck (the only non-exclusive one — its flags are duplicated in `headlessArgs`,
independently verified by kimi-1 at lines 18-25/68-75/103-106 against 28-31/78-81/109-112).
Worked example, `igm-app`: `claude-1` `["--permission-mode","bypassPermissions","--add-dir","parley-deck"]`,
`codex-1` `["--sandbox","workspace-write","-c","approval_policy=\"never\""]`, `hermes-1` `["--yolo"]`
— none present in that agent's `headlessArgs`.

**Verification status:** owned by the drafter. kimi-1 supplied a non-owner `PRIMARY` confirmation
for the one deck inside this repository and explicitly declined the aggregate as out of its
mandate; hermes-1 and codex-1 declined it as outside the repository. **The 23-deck aggregate
therefore stands `UNVERIFIED` by a non-owner** and is used below only as illustration, never as
load-bearing evidence — the condition codex-1 set.

#### Why revision 1's counts were wrong — a measurement-tool trap (`PRIMARY`, drafter)

Revision 1 reported "11 of 23, 8 exclusively". Both numbers were low.

**Correction issued mid-signoff.** The first version of this subsection blamed the filesystem,
asserting that *"`grep -r` silently under-reports on this filesystem"*. That was wrong, and it was
the drafter's **second** false causal explanation in this idea. The real cause is mundane and
entirely local to the drafter's shell:

`grep` in the facilitator's environment is **not** the system `grep`. It is a shell function that
execs **`ugrep` with `--ignore-files`** (`ugrep 7.5.0`, confirmed by `grep --version`), which
honours `.gitignore`. Ignored paths are therefore omitted **silently, with exit status 1** — a
"no matches" that is indistinguishable from a genuine absence.

Both anomalies reduce to that one fact:
- `BYTE/parley-deck/meta/headless-agents.local.json` was skipped because `BYTE/.gitignore:37`
  ignores `parley-deck/` wholesale (`git check-ignore -v` confirms). `aditoLeads`, which has no
  such rule, was found.
- The 1663-vs-172 `*.go` gap is `.gomodcache` (1491 files), likewise ignored.

The system `grep` at `/usr/bin/grep` does not behave this way; neither does `find`.

**Independently corroborated, from the opposite direction — and this resolves R3.** Both
non-owner participants report (`PRIMARY`) that they could **not** reproduce the under-reporting:
codex-1 got `grep -rl '' --include='*.go' . | wc -l` = **1663**, and kimi-1 got 1663 as well plus a
successful `grep -rl writeModeArgs` over the `BYTE` tree. kimi-1 filed this as reservation R3 and,
reading the superseded filesystem-blaming text, correctly verdicted those two measurements
`WRONG / not reproduced`.

**Their non-reproduction is not a contradiction — it is the control case this claim needed.**
Neither participant runs inside the facilitator's shell, so both get the real `grep`; only the
facilitator's shell substitutes `ugrep --ignore-files`. Two environments, two results, one cause.
codex-1 asked for the claim to be demoted to a time-bounded observation with the cause unresolved,
and kimi-1 marked it `DISPUTED` pending resolution. The cause is now identified and the claim is
narrowed instead: **this is a property of the facilitator's shell — not of the filesystem, the
volume, the repository, or any other participant's environment.** Under that scoping the three
measurements agree rather than conflict, and R3 is answered on evidence rather than left standing
as a disagreement.

Consequences, all acted on:
1. Every count in revision 1 came through that shim and was low. Revision 2 re-derives all of them
   with `find … -print0 | xargs -0`, which does not consult ignore files.
2. Claim (a) was originally measured through the same shim. Re-run over all 1663 `*.go` files it
   still returns zero — and, decisively, **codex-1 and kimi-1 confirmed it independently with
   `rg`**, so (a) rests on non-owner evidence, not on the drafter's original command.
3. The generalisable lesson, and the reason this stays in the record: **in the facilitator's shell
   a `grep`/`rg` miss is not evidence of absence.** Both tools honour ignore files by default. Any
   negative result the facilitator offers as proof must come from `find`-based enumeration,
   `rg -uuu`, or `/usr/bin/grep`. Two of the three tools used to establish claim (a) share this
   default, which is exactly the correlated-instrument risk §15.6 asks about.

Per codex-1's reservation, this subsection makes **no** claim about `grep -r` in any environment
other than the facilitator's, and invalidates no other participant's measurement.

#### Drafter's correction (§15.5)

Between round 2 and revision 1 the drafter reasoned from (b) alone that the affected decks would
"silently lose their flags at launch", and said so before checking (a). Withdrawn: nothing is
dropped at launch because the Go launcher never reads that file. Entry 5 in
`## Drafter position changes`.

#### Closure

**VC-1 closes on the design argument, not on the workspace counts** — the narrower ground codex-1
specified. `writeModeArgs` is deleted from the documented JSON shape because a two-field shape
teaches a two-list launch model that the CLI does not implement, and (a) establishes there is no Go
loader whose compatibility deletion could break. The aggregate in (b) illustrates the teaching but
is `UNVERIFIED` and carries nothing.

codex-1's remaining ground — that a labelled manual branch detoxifies the field — is a design
position rather than a checkable claim. It is engaged on argument (hermes-1: a facilitator copying
the JSON example still populates the field and still believes it matters) and recorded as dissent
in `FINAL.md`.

#### Legacy migration rule — adopted (codex-1's condition)

The skill MUST state, in the manual branch: **when an existing `headless-agents.local.json`
contains `writeModeArgs`, merge its arguments into that agent's `headlessArgs` and remove the
field.** Deleting the documented field without this instruction leaves the 9 exclusive decks with
their write-enabling flags in a location that works on the manual path only; a `parley` launch from
TOML config would not pass them, and since 1.39.0 that surfaces as `AUTO=no` with the missing
arguments named rather than failing silently. No deck is edited by this idea — the rule is the
instruction, and the follow-up is recorded in `FINAL.md`.

## Drafter position changes

claude-1 is facilitator, participant and drafter. Required by §15.5.

| # | Prior position | Source | New position | Why |
|---|---|---|---|---|
| 1 | *"I have proposed naming `hermes` as the live case… If two participants think the clause should go, it goes."* | `round-01/claude-1.md` | Withdrawn | codex-1 and kimi-1 both ruled it out on the does-it-change-an-action test. The drafter had applied that test to exclude probe versions and promotion history, then made an exception for the one incident it had personally worked on |
| 2 | Item 4 framed the defect as `writeModeArgs` being a dead field | `round-01/claude-1.md` | The defect is the whole assembly recipe plus the inverted `:251` sentence | codex-1 and kimi-1. The original framing would have produced a patch fixing one symptom of three |
| 3 | Proposed a boundary marker sentence for the recipe | `round-02/claude-1.md` | codex-1's fuller manual/CLI branch split | codex-1's version also covers "Headless Agent Configuration" and `WORKED_EXAMPLES.md`, which a single marker sentence does not reach |
| 4 | Proposed the drafter's own `:251` replacement wording | `round-02/claude-1.md` | kimi-1's converged text | It is a superset; the drafter's merged the CLI surface and the incident into one sentence, which hermes-1 correctly identified as the weaker structure |

| 5 | Reasoned, after round 2, that the 11 decks setting `writeModeArgs` would "silently lose their flags at launch" | drafter's own unfiled reasoning | Withdrawn — nothing is lost because nothing is read | The drafter's own follow-up grep (`PRIMARY`, above). Filed here rather than dropped, because the alarm had already been voiced to the user before it was checked |

| 6 | Revision 1: *"There is no schema: no validator, no loader, no consumer… deleting the field changes the behaviour of zero of the 11 decks, so the cost the argument was protecting against does not exist"* | `consensus.md` rev 1, VC-1 measurement | Withdrawn; the scoped claim is "no Go loader and no Go field" | codex-1's block. The file has a consumer — the manual facilitator — which is exactly what codex-1's position said. The drafter had converted a narrow negative into three broad ones |
| 7 | Revision 1: VC-1 *"CLOSED by evidence, not by count"*, resting on the 23-deck measurement | `consensus.md` rev 1 | Closed on the design argument; the workspace aggregate is demoted to illustration and marked `UNVERIFIED` | codex-1's block: the drafter owns (b), labelled it `PRIMARY` without command/inputs/output/locators, and used it to close a conflict in the drafter's own favour |
| 8 | Revision 1 reported "11 of 23 decks, 8 exclusively" | `consensus.md` rev 1 | 12 contain the key, 10 non-empty, 9 exclusive | The drafter's own re-measurement. The original figures came from `grep -r`, which under-reports on this filesystem (172 of 1663 `*.go` files seen). Not a transcription slip — a bad instrument |
| 9 | Revision 1: the assertion runs *"under both `npm test` and the existing `prepack` lifecycle"* | `consensus.md` rev 1, decision 4 | False; `prepack` runs only `build-addon-manifest.js --check`, and implementation must close that half explicitly | hermes-1 raised it in round 2, kimi-1 restated it as reservation R1, the drafter verified it against `package.json:60,66` |
| 10 | Revision 2, first version: *"`grep -r` silently under-reports on this filesystem"* | `consensus.md` rev 2, before this correction | Withdrawn. The cause is the drafter's shell aliasing `grep` to `ugrep --ignore-files`, which honours `.gitignore` | The drafter's own follow-up (`grep --version` → ugrep 7.5.0; `git check-ignore -v` on `BYTE`). **Corrected mid-signoff**, after the revision-2 signoff requests had already gone out — see the note under `## Signoffs` |

Nine changes: four forced by another participant in the rounds, one by the drafter's own
measurement contradicting the drafter, and four more forced at signoff — three by codex-1's block
and one by kimi-1's reservation. **The drafter's revision-1 text did not survive review**, which is
the relevant fact about it.

## Comparison & blind spots

**Correlated-agreement caveat (§15.6).** The four participants are related models. The unanimous
`CONFIRMED` on the launch-path claim is a shared prior, not four independent confirmations —
mitigated here by the fact that three participants read `runner.go` themselves and quoted different
line ranges, and one (hermes-1) proved the negative the facilitator had only asserted.

**§15.6's trigger did not fire.** Round 1 was not unanimous — hermes-1 opposed the guard — so no
steelman round was required. It became unanimous only after an evidence-based withdrawal in round 2,
which is the mechanism working rather than convergence by anchoring.

**What would have to be true for this to be wrong.** The whole package rests on one claim: that
`HeadlessArgs` is the sole launch input. If any code path appends arguments elsewhere — a wrapper,
a transport-specific launcher, a future ACP path — then the "effective argv" instruction sends
facilitators to an incomplete value. Three participants greped for it independently and none found
one, but all three searched the same repository with the same tool.

**Where nominally independent findings are one family.** codex-1's recipe finding and kimi-1's
`:251` finding are the same defect at two altitudes: the skill treats the declared mode and the
launched command as independent values. The facilitator's `writeModeArgs` item is the third
instance. They were filed separately and should be fixed as one passage, not three patches.

## Signoffs

<!-- Each participant appends its own block. Do not edit another participant's block. -->

### Revision 1 — codex-1 blocked; kimi-1 accepted with one reservation

codex-1's block and kimi-1's reservation R1 are both upheld in full and are answered by
revision 2 (VC-1 rewritten and narrowed, claim (b) demoted to `UNVERIFIED` illustration, the
overreaching inferences withdrawn, the legacy migration rule adopted, decision 4's `prepack`
clause corrected). The blocks below are preserved verbatim and are **not** superseded — they are
the record of what revision 1 got wrong.

Each block below was written by the named participant to its own file
(`signoff-<agent-id>.md`) and concatenated here byte-for-byte by the facilitator without editing.
The facilitator authored only its own block, and wrote it before reading any of the others.

### codex-1

**Verdict:** block

**Scope.** [PRIMARY] I read `parley-deck/COOPERATION.md:1-1316`, including §15 at lines 1176-1316; `00-prompt.md:1-71`; all four files in each of `round-01/` and `round-02/`; and `consensus.md:1-215`. [PRIMARY] I inspected `internal/config/runtime.go:97-165` and ran the repository-local Go search quoted below. [PRIMARY] I did not inspect any of the 23 decks outside this repository and did not run their configs or launchers. [PRIMARY] I issue no verdict on codex-1's owned manual-assembly claim; I only assess whether the drafter's new evidence resolves it, consistent with `parley-deck/COOPERATION.md:1197-1205`.

**VC-1.** I do not accept the closure as written.

- **Claim (a): CONFIRMED [PRIMARY].** I ran `rg -n -g '*.go' 'headless-agents|headless_agents|writeModeArgs|write_mode_args' .` from the repository root; it exited 1 and emitted no output. `agentOverride` contains `HeadlessArgs` but no `writeModeArgs` field (`internal/config/runtime.go:97-132`), and `configLayers` enumerates the central TOML file, the two deck TOML files, and the optional environment-supplied path (`internal/config/runtime.go:134-154`).
- **Claim (b): UNVERIFIED [RECALL] by codex-1.** The 11-of-23 and 8-exclusive counts are outside my checked scope (`consensus.md:136-142`). The draft quotes neither the command and inputs nor relevant output or exact external file locators (`consensus.md:125-142`), so its `PRIMARY` label is malformed and reads as `RECALL` under `parley-deck/COOPERATION.md:1216-1224`.
- **Effect on my minority position.** [PRIMARY] The measurement merely outflanks the machine-schema concern; it does not defeat my argument. My position preserved `writeModeArgs` only as a manual-facilitator assembly input, while denying it any role in Parley CLI configuration (`round-02/codex-1.md:11-27,127-128`). [PRIMARY] Claim (a) is compatible with that position, and the draft itself calls the JSON file manual-facilitation config “consumed by a facilitator hand-assembling a command” (`consensus.md:129-134`); it therefore does not entail the later claims that there is “no consumer” or that deleting the documented field changes the behavior of zero decks (`consensus.md:149-155`). [RECALL] If claim (b) is later established, exclusive legacy use would make an explicit migration rule more important; it would not prove zero compatibility impact.
- **Independence.** [PRIMARY] The drafter supplied the measurement (`consensus.md:125-142`), had already favored deletion (`round-02/claude-1.md:108-111`), and uses the measurement to close VC-1 in favor of deletion (`consensus.md:149-163`). My independent check therefore confirms only (a), not the drafter's full inference.

**Required changes for acceptance:** either obtain a non-owner `PRIMARY` verdict for claim (b) in that verifier's own canonical round artifact, with the exact command, inputs, relevant output, and stable locators, or mark (b) `UNVERIFIED` and remove it from the closure and `FINAL.md`. Reframe claim (a) narrowly as “no Go loader or field,” deleting the unsupported “no consumer,” “zero behavior change,” and “the protected cost does not exist” inferences. If the new documented shape deletes `writeModeArgs`, add an explicit legacy migration rule: when a manual JSON config contains it, merge its arguments into `headlessArgs` before launch and update that config. VC-1 may then close on the single-vector design choice without depending on the unverified workspace counts.

### kimi-1

**Verdict: accept with reservations** (one non-blocking reservation, item R1 below).

#### VC-1 — I accept the closure

The drafter made the measurement and benefits from it, so per §15.1 I re-checked what I could
reach myself instead of endorsing it. I split claim (a) into its two atomic parts, because my
ownership differs across them.

**(a1) Zero Go references to `writeModeArgs` / `write_mode_args` — CONFIRMED (PRIMARY).**
I executed a repo-wide ripgrep for `writeModeArgs|write_mode_args|WriteModeArgs` over **all**
`*.go` files under the repository root — a superset of the drafter's non-test scope, test files
included. Zero hits. Owner of the underlying claim is claude-1 (`round-01/claude-1.md` item 4);
I am a non-owner, and this re-confirms my round-1 verdict (`round-01/kimi-1.md` item A.4) with a
fresh run.

**(a2) The CLI never reads the JSON config file — no verdict from me; I co-own this claim**
(`round-01/kimi-1.md` refinement, lines 99–103; also `round-01/codex-1.md` row "JSON config as a
Parley runtime source"). §15.1 bars my verdict. Raw evidence, reported without classification: a
repo-wide ripgrep for `headless-agents|headless_agents|headlessAgents|HeadlessAgents` over all
`*.go` returned zero hits, and `internal/config/runtime.go:137-154` (`configLayers`) enumerates
exactly `~/.parley/agents.toml` (line 141), `parley-deck/agents.toml` (144),
`parley-deck/agents.local.toml` (145), and `$PARLEY_HEADLESS_AGENT_CONFIG` (147–151) — no JSON
path. The standing non-owner verdict is hermes-1's (PRIMARY, `round-02/hermes-1.md` @codex-1:
"`configLayers:137-154` reads `agents.toml`/`agents.local.toml`, never the skill's
`headless-agents.local.json`"). The claim rests on that, not on the drafter's or my say-so.

**(b) Decks setting the field — CONFIRMED for the one deck inside this repository (PRIMARY);
the aggregate is UNVERIFIED by me.** This repo's own deck,
`parley-deck/meta/headless-agents.local.json`, sets `writeModeArgs` for all four roster agents:
claude-1 (lines 28–31), codex-1 (78–81), hermes-1 (109–112), kimi-1 (154, empty array) —
consistent with the drafter counting it among the 11. It is **not** an exclusive case: for
claude-1, codex-1 and hermes-1 the same enabling flags are duplicated in `headlessArgs`
(lines 18–25, 68–75, 103–106) — consistent with its exclusion from the 8. The other 10 named
decks, the 23-deck denominator, and the `igm-app` exclusivity example are **outside this
repository**; my mandate for this signoff is limited to this repo, so I did not read them and
say so rather than guess. This does not weaken the closure: deletion rests on (a) — no loader,
validator, or consumer exists — and on the structural teaching argument; (b) illustrates that
argument but does not carry it. A substantially wrong count would not change the decision.

**Why the closure is legitimate under §15.3.** codex-1's minority position had one empirical
premise: keeping the field "avoids a schema bump for local configs that already set it"
(`round-02/codex-1.md`, position changes and lines 127–128). The measurement shows there is no
schema, loader, or consumer, so the cost that premise protects against does not exist; and the
majority's teaching argument, previously asserted, is now measured. codex-1's remaining ground —
a labelled manual branch detoxifies the field — is a design position, not a checkable claim, and
the draft engages it with argument (hermes-1's copying-facilitator point, `round-02/hermes-1.md`
@kimi-1) plus the measured exclusivity. That is resolution by evidence and argument, not by the
3-to-1 count; the dissent and the reopen condition are recorded in the draft. Closure accepted.

The codex-1-specific question is N/A to me: I was in the majority of three, having withdrawn my
round-1 "keep the split" (`round-01/kimi-1.md` item C, lines 174–181) in round 2
(`round-02/kimi-1.md`, position change 1) on argument, before any measurement existed. The
measurement neither defeats nor outflanks a position I no longer held; it corroborates the
withdrawal's reason (a).

#### Reservation R1 (non-blocking) — the `prepack` clause of decision 4 is imprecise as drafted

The draft has the equality assertion "running under both `npm test` and the existing `prepack`
lifecycle" (consensus.md line 82). Per hermes-1's PRIMARY finding, which I did not independently
re-run and cite as SECONDARY (`round-02/hermes-1.md` @codex-1; also `round-02/codex-1.md`):
`package.json:66` `prepack` runs only `node scripts/build-addon-manifest.js --check`, which reads
neither version field and does not invoke `node --test`. Landing the decision as intended
therefore requires either editing the existing `prepack` command or adding the check inside
`build-addon-manifest.js` (hermes-1's option b). Both are edits to existing files, so "no new
script, job, or checker" survives — but the draft records only the manifest-*test* half of the
gap (consensus.md lines 88–90), not the `prepack`-invocation half. Implementation must close that
explicitly, or a direct `npm publish` still bypasses the guard — the exact path that produced the
drift. The decision's intent is unambiguous, so this does not block.

#### Drafter position changes (§15.5 ratification)

Changes 1–4 match the round files I read (PRIMARY on the comparison): change 1 against
`round-01/claude-1.md` concern 1 and `round-02/claude-1.md` SELF-CORRECTION 1; change 2 against
SELF-CORRECTION 2; changes 3 and 4 against `round-02/claude-1.md` @codex-1 / @kimi-1. Change 5 is
drafter testimony about unfiled private reasoning: I ratify that it is filed as §15.5 requires;
its content I cannot check and do not confirm. Completeness: comparing `round-02/claude-1.md`
"Current proposal" (six items) against the draft's six decisions, the only deltas are changes 3
and 4, both filed — the table is complete to my reading (PRIMARY).

#### §15.6

The trigger analysis is accurate (PRIMARY, against the round files): round 1 was not unanimous —
hermes-1 opposed the guard (`round-01/hermes-1.md` item 3) — so no steelman section was required.
The §15.6(b) record (correlated-agreement caveat, what-would-have-to-be-true, one-family note)
is present at consensus.md lines 193–211.

I am the author of the adopted `:251` text and the opencode row wording; per §15.1 I issue no
verdict on their quality. My acceptance of the draft is a signoff decision, not a
self-verification.

#### Scope declaration

Read in full for this signoff: `COOPERATION.md` §15 (lines 1176–1316); `00-prompt.md`; all four
round-01 and all four round-02 files; `consensus.md` (all 215 lines). Executed fresh in this
repo: two repo-wide ripgreps over `*.go` (patterns and results above); read
`internal/config/runtime.go:130-159`; read `parley-deck/meta/headless-agents.local.json` (all 179
lines). **Not done:** read any deck outside this repository (the other 10 decks in claim (b), the
23-deck denominator); read the skill repo (`package.json:66`, the test harness — relied on the
named SECONDARY source for R1); ran no test suites and no binaries; re-ran no git diffs (the
bundled-`references/COOPERATION.md` currency claim is facilitator-owned and no element of this
signoff depends on it). I ran no git write commands and edited no file other than this signoff.

### hermes-1

Verdict: accept

## Scope declared

Read in full: COOPERATION.md §15 (lines 1176-1316), 00-prompt.md, all four round-01 files, all four round-02 files, consensus.md (all 215 lines).

Ran for this signoff (PRIMARY):
- A repo-wide search for `headless-agents|headless_agents|writeModeArgs|write_mode_args|WriteModeArgs` across every `*.go` file in `parley-deck-cli` — zero hits. This independently reproduces the drafter's claim (a) grep.
- Read `internal/config/runtime.go:137-154` (`configLayers`) to confirm the enumerated agent-config inputs are exactly: `~/.parley/agents.toml`, `parley-deck/agents.toml`, `parley-deck/agents.local.toml`, and `$PARLEY_HEADLESS_AGENT_CONFIG`. `meta/headless-agents.local.json` is not among them.

Did not run: the skill test suite, the parley CLI, live vendor probes, or the `v1.38.0..v1.39.0` COOPERATION.md currency diff (facilitator-owned claim; no position of mine depends on it).

Could not reach: VC-1 claim (b) — the 23 workspace decks under `AI_WORKSPACE` are outside this repository. I do not verify (b) and do not guess at it.

## VC-1 closure — accepted

**Claim (a) — "the CLI never reads the file the field lives in": CONFIRMED [PRIMARY].**

My search for `writeModeArgs`/`write_mode_args`/`WriteModeArgs`/`headless-agents`/`headless_agents` across all non-test `*.go` files returned zero hits. My read of `configLayers` (runtime.go:137-154) confirms the only agent-config inputs are the three TOML layers plus the env-var path; `meta/headless-agents.local.json` is not read by `parley` at launch. This matches the drafter's measurement verbatim. The drafter both made the measurement and benefits from it (it supports the majority's delete position, which the drafter shared), so §15.1 makes this independent check meaningful — and the check passes.

**Claim (b) — "11 of 23 decks set the field, 8 exclusively": UNVERIFIED.**

Those decks live under `AI_WORKSPACE`, outside this repository. I cannot reach them from here and will not guess. I note that (b) is additive weight for the majority's "the field teaches" argument, but the closure does not stand or fall on (b) alone — see below.

**The drafter's correction (§15.5, position change #5) is sound.** The drafter withdrew the private alarm that the 11 decks would "silently lose their flags at launch." That correction is logically entailed by (a) alone: if the CLI never reads the file, nothing is dropped, because nothing is consumed. It does not depend on (b), and I can evaluate it on (a), which I verified. The two facts being stated together — "(a) nothing is read, (b) the field is widely set" — is the right framing; (b) without (a) reads as an alarm, and the drafter correctly says so.

**codex-1's position: the measurement defeats the argument, it does not outflank it.**

codex-1's ground was "keeping it avoids a schema bump for local configs that already set it" (consensus.md:118-119, quoting codex-1 round-01). Claim (a) shows there is no schema — no validator, no loader, no consumer in the CLI. Deleting the field from the documentation changes the runtime behaviour of zero decks. The cost codex-1 was protecting against does not exist; the argument's factual premise is removed, not sidestepped. That is a direct defeat.

codex-1's structural insight — that manual facilitator assembly is a different activity from CLI config, and the recipe should branch — is honored by the consensus (item 3 adopts the manual/CLI branch split). What codex-1 does not get is keeping the field inside the manual branch. That is the right call: even on the manual path, a two-field shape teaches that the write-enabling flag can live separately from `headlessArgs`, which is the exact model that made the defect invisible. The branch split accommodates codex-1's distinction; the field deletion addresses the teaching model. Both can be true.

I was in the majority (delete) in rounds 1 and 2. The measurement adds empirical weight I could not produce (I did not survey the workspace decks), but my position did not depend on (b) — it rested on (a), which I verified myself in round 1 (PRIMARY: `agentOverride` at runtime.go:97-132 has no `writeModeArgs` field) and have now re-verified at the repo-wide level for this signoff.

## The six edits

1. **opencode row** — agree. Unanimous; kimi-1's wording carries the argv detail. No reservation.
2. **SKILL.md:251 replacement** — agree. I supported the converged text in round 2 (round-02/hermes-1.md:185-208). kimi-1's version is the superset: its lead sentence is the inversion fix, its body is codex-1's contract framing, its mechanism clause is mine, and the existing confinement floor is retained. The consensus correctly replaces both sentences at :251 rather than adding a third.
3. **Manual/CLI branch split** — agree. I proposed a boundary marker in round 2; the consensus adopts codex-1's fuller split. I acknowledged in round 2 that codex-1's version "also covers 'Headless Agent Configuration' and WORKED_EXAMPLES.md, which a single marker sentence does not reach" (round-02/hermes-1.md:169). The broader split is within the range I accepted.
4. **compatibility.json — bump plus one assertion** — agree. I withdrew my round-1 "no guard" objection in round 2 on evidence (round-02/hermes-1.md:11-24). The consensus runs the assertion under both `npm test` and `prepack` (consensus.md:81-82), which closes the prepack gap I raised (round-02/hermes-1.md:80-92). My implementation-gap note — that the existing manifest test verifies hashes and does not read version fields, so the assertion needs its own home — is recorded in the consensus (consensus.md:88-90).
5. **writeModeArgs deleted** — agree, per VC-1 closure above.
6. **What stays out** — agree. Unanimous. The hermes incident narrative, promotion history, probe versions, ACP availability, kimi's exit-1 message, and CLI struct names all fail the "does a facilitator behave differently" test. `references/COOPERATION.md` untouched per the byte-identical constraint.

## §15.5 drafter position changes — ratified

The five changes in the consensus table (consensus.md:179-189) are accurate and complete from my reading of claude-1's round-01 and round-02 files. Change #5 (the withdrawn "silently lose flags" alarm) is correctly filed under §15.5 rather than dropped, and the correction is sound per (a) above.

## §15.6 correlated agreement

The consensus correctly notes the §15.6 trigger did not fire: round 1 was not unanimous (I opposed the guard), so no steelman round was required (consensus.md:198-200). The correlated-agreement caveat — that the four participants are related models and the unanimous CONFIRMED on the launch-path claim is a shared prior — is accurately stated and mitigated by the three independent `runner.go` readings and my proof of the negative (consensus.md:194-196).

No reservations. The consensus is sound, the VC-1 closure is evidence-based and independently checked on the part I can reach, and the drafter's position changes are transparent and correct.

### claude-1

**Verdict:** accept — with the ownership limit below, which is not a reservation about the content
but a statement of what my acceptance is worth.

**Role disclosure (§15.5).** I am facilitator, participant and drafter of this consensus. I also
made the VC-1 measurement that closes the one open conflict, and that measurement supports the
position I held in rounds 1 and 2. That is the exact configuration §15.1 exists to constrain.

**Scope declared (§15.5).** For this signoff I re-read all four round-01 files, all four round-02
files, `00-prompt.md`, and the full consensus draft. I ran the two VC-1 measurements myself
(`grep -rn` over non-test `*.go`; a Python pass over all 23 `headless-agents.local.json` files
under `AI_WORKSPACE`). I did **not** re-verify the `runner.go` launch-path claim for this signoff —
it was confirmed independently by hermes-1 in round 1 and I own it, so I issue no verdict on it.

**On VC-1 I issue no verdict, and the closure is provisional until a non-owner confirms it.**
I own both halves of the measurement, so under §15.1 I cannot supply the verdict that makes it
admissible. What I can state is the method, so it is checkable rather than trusted:

- Claim (a) — zero Go references — is reproducible in this repository with
  `grep -rn "headless-agents\|headless_agents\|writeModeArgs\|write_mode_args" --include="*.go" .`
  filtered to non-test files, plus reading `configLayers` at `internal/config/runtime.go:137`.
  Any participant can run it; a single counter-hit refutes me.
- Claim (b) — 11 of 23 decks, 8 exclusively — is outside this repository and **not** reproducible
  by a reviewer confined to it. I have named the decks and quoted `igm-app`'s three agents so the
  claim is at least specific enough to be wrong. A participant that cannot reach those paths
  should say so rather than concur.

If codex-1 disputes the measurement itself, VC-1 reopens and the deletion comes out.

**What I got wrong and am recording rather than smoothing.** Between round 2 and this draft I
concluded from (b) alone that the 11 decks were "silently losing their flags at launch", and I
said so to the user before checking (a). It was wrong — nothing is dropped because nothing is
read. It is entry 5 in `## Drafter position changes`. I am noting it here as well because it is
the second time in this idea that I stated a conclusion at a confidence the evidence did not yet
carry, and the pattern matters more than the individual error.

**On the six adopted edits: accept, and three of them are not mine.** Items 2 (`SKILL.md:251`),
3 (manual/CLI branch split) and 4's implementation detail came from kimi-1, codex-1 and hermes-1
respectively, and each replaced a weaker proposal I had filed. The draft records that in
`## Drafter position changes`; I want it in the signoff too, because a consensus whose drafter is
also the majority can otherwise read as self-ratification.

**One thing the package does not fix, deliberately.** The 8 decks carrying write-enabling flags
only in `writeModeArgs` are left as they are. The skill edit makes the correct location explicit,
which is the migration instruction, but no deck is touched by this idea. If a reviewer thinks that
is too little, the place to say it is now, not after the release.

### Revision 2 — accepted by all four; no blocks

| participant | verdict |
|---|---|
| codex-1 | accept with reservations — revision-1 block withdrawn, all three conditions met |
| hermes-1 | accept, no reservations |
| kimi-1 | accept with reservations (R1 met; R2 and R3 filed and answered) |
| claude-1 | accept (drafter; no verdict on any claim it owns) |

Each block below was written by the named participant to its own `signoff2-<agent-id>.md` and
concatenated here byte-for-byte without editing.

**Disclosure: revision 2 was corrected while its signoff round was already open.** The signoff
requests went out with a subsection asserting that *"`grep -r` silently under-reports on this
filesystem"*. The drafter then established that this was false — the cause is a shell alias to
`ugrep --ignore-files` honouring `.gitignore` — and replaced the subsection (drafter position
change 10). Participants whose signoff quotes the filesystem claim were reading the superseded
text; the corrected subsection reaches the same corrected **numbers** (12/10/9) by the same
`find`-based commands, so no decision moves, but any participant who wishes to revise on this
basis should say so and its later block supersedes its earlier one.

### codex-1 — revision 2

**Verdict:** accept with reservations

**Scope.** [PRIMARY] I read `parley-deck/COOPERATION.md:1-1316`, the full revised `consensus.md:1-548`, including all embedded revision-1 signoffs, my original `signoff-codex-1.md:1-14`, and the ownership-relevant passages in `round-01/codex-1.md` and `round-02/codex-1.md`. [PRIMARY] From the repository root I re-ran the repository-local Go search and the two file-count commands quoted below. [PRIMARY] I did not inspect the 23 external deck configs, run the JSON aggregation, inspect `package.json`, run tests or binaries, or run any git command. [PRIMARY] I edited only this signoff file.

**codex-1 required changes.**

- [PRIMARY] **(i) Met.** Claim (b) is expressly `UNVERIFIED` by a non-owner at `consensus.md:209-213`, and VC-1 closes on the design argument rather than that aggregate at `consensus.md:241-247`.
- [PRIMARY] **(ii) Met.** Claim (a) is narrowed to “no Go loader and no Go field” at `consensus.md:153-179`; the “no consumer,” “zero behaviour change,” and “protected cost does not exist” inferences are explicitly withdrawn at `consensus.md:175-179,275-278`.
- [PRIMARY] **(iii) Met.** The legacy rule at `consensus.md:254-262` says to merge existing `writeModeArgs` into `headlessArgs` and remove the old field. [PRIMARY] Because codex-1 owns that proposed rule, I issue no verification verdict on its correctness; I verify only that revision 2 incorporated the condition I set.

**kimi-1 R1.** [PRIMARY] R1 is correctly recorded: revision 2 states that `prepack` does not run `node --test` and requires implementation either to extend `prepack` or place the equality check in `build-addon-manifest.js` (`consensus.md:83-96`), matching kimi-1's preserved reservation at `consensus.md:397-409`. [SECONDARY] I did not independently inspect `package.json`; the command-level fact rests on the named PRIMARY checks by hermes-1 and the drafter recorded at `consensus.md:83-90`.

**Measurement and defective-tool audit.**

- [PRIMARY] I issue no verdict on corrected claim (a), because codex-1 owns its no-JSON-loader premise (`round-01/codex-1.md:14-18,31-33`). [PRIMARY] Raw evidence only: I ran `rg -n -g '*.go' 'headless-agents|headless_agents|writeModeArgs|write_mode_args|WriteModeArgs' .`; it produced no output and exited 1. [PRIMARY] The draft's admissible support is instead the named non-owner PRIMARY verdicts at `consensus.md:170-173`, alongside the command, output and locators at `consensus.md:153-168`; that is adequate §15.2 evidence without a codex-1 self-verdict.
- [PRIMARY] Corrected claim (b) is adequately handled only as an illustration, not adequately evidenced for a `CONFIRMED` verdict: the second pipeline is abbreviated with `...`, and the JSON pass is described but not reproduced (`consensus.md:183-207`). [PRIMARY] Its explicit `UNVERIFIED` status and removal from the closure therefore remain necessary and sufficient for this signoff (`consensus.md:209-213,241-247`).
- [PRIMARY] I could not reproduce the categorical tooling defect at `consensus.md:215-232`. From the repository root, `find . -name '*.go' -type f | wc -l` returned `1663`, and `grep -rl '' --include='*.go' . | wc -l` also returned `1663`, not `172`. [RECALL] The drafter's earlier `172` result may have depended on transient filesystem state, but this signoff has no evidence establishing a cause.
- [PRIMARY] No adopted decision still rests on the defective `grep -r` result: claim (a) was re-run with `find` and independently with `rg`, claim (b) was re-derived separately and demoted, and VC-1 closes on design (`consensus.md:223-247`). [PRIMARY] The remaining reliance is explanatory: the categorical defect/future-warning text at `consensus.md:215-232` and drafter-change entry 8 at `consensus.md:277-280`; the old counts in embedded revision-1 signoffs are explicitly preserved as historical text at `consensus.md:313-323`.

**Reservation.** [PRIMARY] Before FINAL, the tooling statement should be recorded as a time-bounded drafter observation with cause unresolved, rather than as a categorical defect invalidating future `grep -r` measurements; alternatively, it should be carried as `DISPUTED` with the dependency check required by `COOPERATION.md:1235-1253`. [PRIMARY] The unverified “9 exclusive decks” figure repeated in the migration rationale at `consensus.md:254-262` should remain illustration and should not become FINAL's support for the conditional migration rule. These points do not revive my revision-1 block because none of the adopted decisions depends on either claim.

### hermes-1 — revision 2

Verdict: accept

I signed off on revision 1 with no reservations. Revision 2 rewrites VC-1, corrects
decision 4's `prepack` clause, and discloses a measurement-tool defect. Below I assess
each change, re-verify what I can reach, and record one thing my own revision-1 signoff
got wrong that revision 2 fixes.

----------------------------------------------------------------

## codex-1's three required changes — all met

codex-1's revision-1 block (consensus.md:378-391) set three conditions. I assess each
against the revised text.

### (i) claim (b) marked UNVERIFIED and removed from the closure — MET

`PRIMARY` The consensus now states, at lines 226-230: "The 23-deck aggregate therefore
stands `UNVERIFIED` by a non-owner and is used below only as illustration, never as
load-bearing evidence — the condition codex-1 set." The closure (lines 293-299) is
explicit: "VC-1 closes on the design argument, not on the workspace counts … The
aggregate in (b) illustrates the teaching but is `UNVERIFIED` and carries nothing."

The claim-(b) text (lines 190-224) now carries the exact commands (`find … -print0 |
xargs -0 grep -l`), the inputs (the workspace root path), the output (the 23/12/10/9
table at lines 210-216), and locators (the named 10-deck list at lines 218-221, plus
the `igm-app` worked example at lines 222-224). This cures the §15.2 malformation
codex-1 identified: revision 1's `PRIMARY` label had no command/inputs/output/locators
and read as `RECALL`.

### (ii) claim (a) narrowed, three inferences deleted — MET

`PRIMARY` Claim (a) is now scoped (lines 153-156) to "no Go code in `parley-deck-cli`
declares a `writeModeArgs`/`write_mode_args` field or reads
`meta/headless-agents.local.json`." Lines 182-186 explicitly withdraw the three
overreaching inferences: "it does not show the field has 'no consumer', that deleting
it is a 'zero behaviour change', or that 'the protected cost does not exist'. Those
three inferences are withdrawn from revision 1. The file does have a consumer — a
facilitator hand-assembling a command, exactly as codex-1's position held."

I independently re-verified the scoped claim (a) for this signoff:
- `PRIMARY` I ran a repository-wide search (`search_files`, ripgrep-backed) for
  `headless-agents|headless_agents|writeModeArgs|write_mode_args|WriteModeArgs` across
  all `*.go` files under `parley-deck-cli` — zero hits.
- `PRIMARY` I read `internal/config/runtime.go:97-132` (`agentOverride`): it has
  `HeadlessArgs` and no write-mode field. I read `internal/config/runtime.go:134-154`
  (`configLayers`): it enumerates exactly `~/.parley/agents.toml` (line 141),
  `parley-deck/agents.toml` (144), `parley-deck/agents.local.toml` (145), and
  `$PARLEY_HEADLESS_AGENT_CONFIG` (147-151) — no JSON path.

This reproduces my revision-1 check and the drafter's `find`-based re-derivation. I am
a non-owner of claim (a) (claude-1 owns it, `round-01/claude-1.md` item 4), so my
`CONFIRMED` is admissible under §15.1.

### (iii) the legacy migration rule added — MET

`PRIMARY` Lines 306-314 adopt codex-1's condition: "when an existing
`headless-agents.local.json` contains `writeModeArgs`, merge its arguments into that
agent's `headlessArgs` and remove the field." The rule is stated as a skill instruction
(no deck is edited by this idea), and the follow-up is deferred to `FINAL.md`. This is
the rule codex-1's block required ("add an explicit legacy migration rule: when a
manual JSON config contains it, merge its arguments into `headlessArgs` before launch
and update that config").

----------------------------------------------------------------

## kimi-1's reservation R1 — correctly recorded

`PRIMARY` R1 (kimi-1's block, consensus.md:450-462) said the draft's claim that the
assertion runs "under both `npm test` and the existing `prepack` lifecycle" was
imprecise, because `prepack` runs only `build-addon-manifest.js --check` and does not
invoke `node --test`, so a test-file assertion leaves `npm publish` ungated.

`PRIMARY` Revision 2 corrects this at lines 83-96: "`package.json:60` defines `test` as
`node --test && …`; `package.json:66` defines `prepack` as
`node scripts/build-addon-manifest.js --check` alone. `prepack` does not invoke
`node --test` … Implementation MUST close the `prepack` half explicitly, by one of
(a) extend the `prepack` command to run the version check, or (b) put the equality
check inside `scripts/build-addon-manifest.js`, which `prepack` already runs."

I independently verified the prepack facts against the skill source repo
(`/Volumes/My Shared Files/AI_WORKSPACE/parley-deck/parley-deck-skill`), at HEAD
(commit `76b31bc`):
- `PRIMARY` `package.json:60`: `"test": "node --test && node scripts/run-python-tests.js
  && node scripts/build-addon-manifest.js --check"`.
- `PRIMARY` `package.json:66`: `"prepack": "node scripts/build-addon-manifest.js --check"`.
  `prepack` does not invoke `node --test`.

Ownership note (§15.1): I first raised the prepack gap in `round-02/hermes-1.md @codex-1`,
so I co-own the underlying claim and issue no verdict on it here. The admissible
non-owner verdict is the drafter's (`PRIMARY`, consensus.md:88-89, verified against
`package.json:60,66`). My reads above are evidence, not a self-verification.

----------------------------------------------------------------

## The defective measurement tool — assessment

The drafter disclosed that revision-1 numbers came through a shell shim: `grep` in the
facilitator's environment is aliased to `ugrep --ignore-files`, which honours
`.gitignore` and silently omits ignored paths with exit status 1 (consensus.md:232-284,
drafter position change 10 at line 333). Revision 1's first causal explanation —
"`grep -r` silently under-reports on this filesystem" — was the drafter's second false
causal explanation in this idea and was corrected mid-signoff to the shell-alias cause.

### Is the corrected measurement adequately evidenced?

Yes, where it is load-bearing.

- Claim (a) — the load-bearing factual claim for the VC-1 closure — does NOT rest on
  the defective tool. `PRIMARY` I re-verified it with my own search (zero hits). The
  consensus records independent non-owner confirmations from codex-1 (`rg`, exit 1, no
  output) and kimi-1 (`rg` over all `*.go`, zero hits) at lines 170-173. The drafter's
  own re-derivation uses `find … -print0 | xargs -0`, which does not consult ignore
  files (line 161-163). So (a) has three independent non-owner `PRIMARY` confirmations
  plus a `find`-based owner run — well-evidenced, not instrument-correlated.

- Claim (b) aggregate (12/10/9) — the drafter re-derived it with `find … -print0 |
  xargs -0` and a `while read` loop (lines 192-208). It is marked `UNVERIFIED` by a
  non-owner and demoted to illustration (lines 226-230, 293-299). I cannot reach the 22
  decks outside this repository and do not verify the aggregate. `PRIMARY` I did verify
  the one in-repo deck: `parley-deck/meta/headless-agents.local.json` sets
  `writeModeArgs` on all four agents — claude-1 (lines 28-31, non-empty), codex-1
  (78-81, non-empty), hermes-1 (109-112, non-empty), kimi-1 (154, empty array) — and
  for claude-1/codex-1/hermes-1 the flags are duplicated in `headlessArgs` (lines
  18-25, 68-75, 103-106). This matches the consensus description at lines 218-221
  ("this repo's own deck, the only non-exclusive one") and kimi-1's verification. Since
  the aggregate carries nothing load-bearing, its UNVERIFIED status is adequate.

### Does anything else in the draft still rest on the defective tool?

I checked every measurement claim in the consensus:

- `references/COOPERATION.md` currency (line 119-120): verified by `git diff
  v1.38.0..v1.39.0`, not by `grep`. Not affected. (Facilitator-owned; I did not
  re-verify, as in my revision-1 signoff — no position of mine depends on it.)
- "nothing under `test/` or `scripts/` reads `compatibility.json`'s `skillVersion`
  today" (lines 104-105): `PRIMARY` I verified this at HEAD with `/usr/bin/grep -rn
  "skillVersion" . --include='*.js' --include='*.json'` in the skill repo — zero hits
  in `test/` or `scripts/` at HEAD. This could have been a false negative if the
  drafter used the ugrep shim, but my non-defective-tool check confirms it. Not
  affected.
- The `runner.go:1094-1108` reference (item 3) and `package.json:60,66` reads
  (decision 4): direct file reads, not `grep` over a tree. Not affected.
- The `parley-addon.json` sha256-pin claim (lines 75-77): verified by reading the
  manifest-check code, not by tree grep. Not affected.

Nothing else in the draft rests on the defective tool.

### A note on the mid-signoff correction

`PRIMARY` The consensus discloses at lines 601-609 that revision 2 was corrected while
signoff requests were already out (the filesystem-claim text was replaced with the
ugrep-shell-alias explanation). The corrected text reaches the same numbers (12/10/9)
by the same `find`-based commands; no decision moves. I am signing off on the corrected
text. My revision-1 signoff predates revision 2 and never quoted the filesystem claim,
so I have nothing to revise on this basis.

----------------------------------------------------------------

## What my own revision-1 signoff got wrong

My revision-1 signoff (consensus.md:525-527) argued: "Claim (a) shows there is no
schema — no validator, no loader, no consumer in the CLI. Deleting the field from the
documentation changes the runtime behaviour of zero decks. The cost codex-1 was
protecting against does not exist; the argument's factual premise is removed, not
sidestepped. That is a direct defeat."

That is the same overreach revision 2 withdraws. Claim (a) establishes "no Go loader
and no Go field" — it does not establish "no consumer" (the manual facilitator is a
consumer) or "zero behaviour change" or "the protected cost does not exist." I made the
inference the drafter made, and codex-1's block correctly identified it as overreaching.
`RECALL` I did not recognise this in my revision-1 signoff because I shared the
majority's delete position and read (a) as stronger than its scope warranted.

Revision 2 closes VC-1 on the design argument alone: a two-field shape teaches a
two-list launch model the CLI does not implement, and (a) establishes no Go loader
exists whose compatibility deletion could break. That is the argument I was actually
making in rounds 1 and 2 (my round-1 position rested on `agentOverride` having no
`writeModeArgs` field, `PRIMARY`), without the overreach. I accept the revision-2
closure as the correct and cleaner framing.

----------------------------------------------------------------

## The six edits

Edits 1, 2, 3, and 6 are unchanged from revision 1. I assessed them in my revision-1
signoff (consensus.md:533-540) and have nothing to add.

Edit 4 (`compatibility.json` bump plus assertion): the `prepack` correction (lines
83-96) is the change. I verified the underlying facts at HEAD (`PRIMARY`, above). The
consensus correctly requires implementation to close the `prepack` half by option (a)
or (b). `PRIMARY` I confirmed option (b) is viable: `prepack` (package.json:66) already
runs `build-addon-manifest.js --check` at HEAD, so placing the equality check inside
that script gates `npm publish` without a new script.

Observation for Phase 5 (non-blocking): `PRIMARY` the skill source repo's working tree
already contains uncommitted changes implementing D4 option (b) — a `versionSyncProblem`
function at `scripts/build-addon-manifest.js:85-96` and a test at
`test/manifest-coverage.test.js:489-500` (which references "Idea `skill-sync-cli-1-39`,
decision D4"). At HEAD these do not exist. This is an in-progress implementation
consistent with the consensus direction; it does not affect this signoff, but the
implementer should be aware the working tree is ahead of HEAD.

Edit 5 (`writeModeArgs` deleted): per VC-1 closure above. Agree.

----------------------------------------------------------------

## §15.5 drafter position changes

`PRIMARY` The table at lines 320-333 now records ten changes. Changes 1-5 match the
round files I read (I ratified 1-4 in my revision-1 signoff at consensus.md:542-544;
change 5 I ratified at lines 523-524). Changes 6-9 record the revision-1-to-revision-2
corrections codex-1's block and kimi-1's reservation forced. Change 10 records the
mid-signoff correction of the second false causal explanation. `PRIMARY` I verify
change 10 against the consensus body: lines 232-284 give the corrected ugrep-shell-alias
explanation, and line 333 records the withdrawal. The table preserves change 8's
original (now-superseded) causal claim ("under-reports on this filesystem") alongside
change 10's correction — this is §15.5 working as intended (file the change, do not
smooth it), not an inconsistency to fix.

## §15.6 correlated agreement

Unchanged from revision 1. `PRIMARY` The consensus at lines 340-360 records the
correlated-agreement caveat, the what-would-have-to-be-true condition, and the
one-family note. Round 1 was not unanimous (I opposed the guard), so no steelman round
was required. I have nothing to add to my revision-1 assessment
(consensus.md:546-548).

----------------------------------------------------------------

## Scope declared

Read in full for this signoff: `COOPERATION.md` §15 (lines 1176-1316);
`consensus.md` (all 610 lines, including the embedded revision-1 signoff blocks at
lines 362-598 and the revision-2 disclosure at 601-609).

Ran for this signoff (`PRIMARY`):
- Repository-wide search for
  `headless-agents|headless_agents|writeModeArgs|write_mode_args|WriteModeArgs` across
  all `*.go` files in `parley-deck-cli` — zero hits.
- Read `internal/config/runtime.go:97-154` (`agentOverride` struct and `configLayers`).
- Read `parley-deck/meta/headless-agents.local.json` (all 179 lines) to verify the
  in-repo deck against consensus lines 218-221.
- Located the skill source repo at
  `/Volumes/My Shared Files/AI_WORKSPACE/parley-deck/parley-deck-skill`; read
  `package.json:55-69` and verified `prepack`/`test` at HEAD via `git show HEAD:package.json`.
- Ran `/usr/bin/grep -rn "skillVersion"` across the skill repo to confirm "nothing
  reads `skillVersion` today" (consensus lines 104-105) at HEAD — zero hits in
  `test/`/`scripts/`.
- Read `scripts/build-addon-manifest.js:80-114` and `test/manifest-coverage.test.js:485-500`
  in the skill repo working tree (uncommitted changes implementing D4 option (b)).

Did not run: the skill test suite (`npm test`); the `parley` CLI; live vendor probes;
the `v1.38.0..v1.39.0` COOPERATION.md currency diff (facilitator-owned; no position of
mine depends on it).

Could not reach: the 22 `headless-agents.local.json` decks outside this repository
(claim (b) aggregate). I do not verify the 12/10/9 aggregate and do not guess at it.

I ran no git write commands and edited no file other than this signoff.

### kimi-1 — revision 2

**Verdict:** accept with reservations — two non-blocking reservations (R2, R3), both confined to
the VC-1 measurement appendix. Every required change from revision 1 is met, and the six
decisions stand.

#### codex-1's three required changes — all met

**(i) Claim (b) marked `UNVERIFIED` and removed from the closure — MET, and now exceeded.**
The draft demotes the 23-deck aggregate to illustration, marks it `UNVERIFIED` by a non-owner,
and closes VC-1 on the design argument (`consensus.md:209-213`, `:241-247`; position change 7,
`:278`). That satisfies the condition as codex-1 stated it. Additionally, I have now verified
the aggregate myself as a non-owner (owner: claude-1) — **CONFIRMED [PRIMARY]**:

- `find "/Volumes/My Shared Files/AI_WORKSPACE" -maxdepth 5 -name headless-agents.local.json
  -path "*/parley-deck/meta/*" | sort` → **23 files** (2m26s on this volume; an unbounded
  re-run of the drafter's exact find returned the same 23).
- A `while IFS= read -r f; do grep -l writeModeArgs "$f"; done` loop over those 23 → **12
  files**, the same 12 the draft names (`consensus.md:201-204`), BYTE and
  `ecb-meeting-2026.05` included.
- A `python3` JSON pass comparing each agent entry's `writeModeArgs` against its
  `headlessArgs` → **10 non-empty, 9 exclusive, 2 empty**; the empty two are exactly `BYTE`
  and `ecb-meeting-2026.05`; the exclusive 9 match the draft's list; this repo's own deck is
  the one non-exclusive case (flags duplicated in `headlessArgs`).
- `igm-app` worked example spot-checked verbatim: claude-1
  `["--permission-mode","bypassPermissions","--add-dir","parley-deck"]`, codex-1
  `["--sandbox","workspace-write","-c","approval_policy=\"never\""]`, hermes-1 `["--yolo"]`,
  none present in the respective `headlessArgs` (`consensus.md:205-207`).

So (b) enters the record doubly: still correctly excluded from the closure per codex-1's
condition, and now also non-owner-`PRIMARY` on the numbers.

**(ii) Claim (a) narrowed, three inferences deleted — MET.** The scoped claim is now "no Go
code declares a `writeModeArgs`/`write_mode_args` field or reads
`meta/headless-agents.local.json`" (`consensus.md:153-156`); "no consumer", "zero behaviour
change" and "the protected cost does not exist" are explicitly withdrawn and the manual
facilitator is acknowledged as the consumer (`:175-179`; position change 6, `:277`). Nothing
of the revision-1 overreach survives in the closure.

**(iii) Legacy migration rule — MET.** Adopted verbatim in substance at `consensus.md:254-262`:
when an existing `headless-agents.local.json` contains `writeModeArgs`, merge its arguments
into that agent's `headlessArgs` and remove the field; no deck is edited by this idea; the
1.39.0 `AUTO=no` surfacing is named. This is codex-1's condition as filed.

#### kimi-1's reservation R1 — correctly recorded

MET. Decision 4's revision-2 correction (`consensus.md:83-96`) states that `prepack` runs only
`node scripts/build-addon-manifest.js --check`, that a test-file assertion runs on `npm test`
only, and that a direct `npm publish` would bypass it; implementation MUST close the half
explicitly via option (a) or (b), and "shipping only the test-file assertion does **not**
satisfy this decision". I verified the underlying facts **[PRIMARY]** against
`parley-deck-skill/package.json`: `"test": "node --test && node scripts/run-python-tests.js &&
node scripts/build-addon-manifest.js --check"` and `"prepack": "node
scripts/build-addon-manifest.js --check"`. The draft's `node --test && …` ellipsis
(`consensus.md:85`) is accurate. Position-change entry 9 (`:280`) records the correction with
the right attribution (hermes-1 raised it, kimi-1 restated it as R1, drafter verified).

#### The defective-tool disclosure — assessment

**The corrected measurement is now adequately evidenced.** Every load-bearing number has been
re-derived with `find`-based tooling and, as of this signoff, carries a non-owner `PRIMARY`
(my reproduction above). Claim (a) does not rest on the suspect tool at all: the revision-2
re-run uses `find … -print0 | xargs -0 grep` over all 1663 `*.go` files (`consensus.md:158-168`),
and codex-1, hermes-1 and I each confirmed it independently with `rg`, which shows no defect
(`:170-173`).

**(a1) Zero Go references to `writeModeArgs`/`write_mode_args`/`WriteModeArgs` — CONFIRMED
[PRIMARY], fresh run:** `rg -l -g '*.go'
'headless-agents|headless_agents|writeModeArgs|write_mode_args|WriteModeArgs' .` from the
repository root exited 1 with no output; `internal/config/runtime.go:97-132` (`agentOverride`)
has `HeadlessArgs` and no write-mode field. **(a2) The CLI never reads the JSON file — no
verdict from me; I co-own it** (`round-01/kimi-1.md:99-103`, §15.1). Raw evidence only: the
same `rg` run covered `headless-agents|headless_agents` (zero hits), and
`internal/config/runtime.go:134-154` (`configLayers`) enumerates `~/.parley/agents.toml`
(`:141`), `parley-deck/agents.toml` (`:144`), `parley-deck/agents.local.toml` (`:145`),
`$PARLEY_HEADLESS_AGENT_CONFIG` (`:147-151`) — no JSON path. The standing non-owner verdicts
are codex-1's and hermes-1's `PRIMARY`s.

**However, the tooling-defect claim itself I could not reproduce — Reservation R3
(non-blocking).** The drafter claims `grep -r` saw 172 of 1663 `*.go` files and missed
`BYTE/parley-deck/meta/headless-agents.local.json` entirely (`consensus.md:215-221`). My runs
**[PRIMARY]**, same volume, minutes before this signoff: `grep -rl '' --include='*.go' .` in
this repository returned **1663** files — identical to `find . -name '*.go' -type f | wc -l`
(1663) — and `grep -rl writeModeArgs` over the BYTE tree **did** return the meta file (4
matches on a direct grep). As a non-owner I therefore verdict the two specific defect
measurements **WRONG / not reproduced**; I concede the behaviour may be intermittent or
mount-state-dependent — the volume is demonstrably slow (an unbounded `find` over
`AI_WORKSPACE` took minutes) — but I have no evidence of under-reporting, and my revision-1
signoff's numbers were gathered with `rg`, not `grep -r`. Per §15.3 this is a contradictory
verdict: the defect claim is `DISPUTED` unless resolved. It does not block, because **no
decision or acceptance criterion depends on it**: claim (a) stands on `rg`-based non-owner
evidence, claim (b) stands on my `find`-based reproduction, and revision-1's undercount (11/8
vs the correct 12/9) is a fact whatever its cause. `FINAL.md` MUST record this dependency
check and this dissent (`§15.3`). The draft's consequence 1 — re-derive with
`find … -print0`, distrust `grep -r` on this volume — is sound practice regardless of the
verdict on the explanation.

**Reservation R2 (non-blocking) — the published claim-(b) command is itself malformed on this
volume.** The appendix prints `find "/Volumes/My Shared Files/AI_WORKSPACE" … | sort` then
`… | xargs grep -l writeModeArgs` (`consensus.md:185-189`). The volume path contains spaces,
and `xargs` without `-0` splits on blanks: my run of exactly that pipeline shape returned
**0** matching files, not 12 **[PRIMARY]**. The *numbers* are correct — I reproduced 12/10/9
with a `while read` loop, quoted above — but the *printed command* cannot have produced them
and will mislead the next measurer, which is precisely the failure mode the tooling-defect
section exists to prevent. Fix the record to `find … -print0 | xargs -0 grep -l
writeModeArgs`. Non-blocking because (b) is illustration-only and a working non-owner command
is now on the record in this signoff.

**Does anything else rest on the defective tool? Checked — no.** (a) was re-run `find`-based
and confirmed with `rg` (above). The decision-4 gap claim "nothing under `test/` or `scripts/`
reads `skillVersion`" (`consensus.md:102-105`) — re-run by me with `rg -n skillVersion test/
scripts/` in `parley-deck-skill`: zero hits, rc=1; `skillVersion` occurs in the whole skill
repo only in `skills/parley-deck/references/compatibility.json` itself **CONFIRMED [PRIMARY]**.
The `references/COOPERATION.md` currency claim is `git diff`-based, not `grep`-based;
facilitator-owned; I did not re-run it and nothing in this signoff depends on it (as in
revision 1).

#### §15.5 drafter position changes — ratified

Entries 6-9 match what I can check **[PRIMARY on the comparisons]**: 6 and 7 against the
embedded revision-1 text (`consensus.md:325-338`) versus the rewritten `:153-179` and
`:241-247`; 8 against my own reproduction (11/8 → 12/10/9 confirmed as the correct figures);
9 against `package.json` (above). Entries 1-5 were ratified in my revision-1 signoff and are
unchanged. The table's self-assessment — "the drafter's revision-1 text did not survive
review" (`:282-285`) — is accurate.

#### §15.6

Unchanged from revision 1; my ratification stands [PRIMARY against the round files read for
revision 1]: round 1 was not unanimous, so no steelman section was required; the
correlated-agreement caveat, the what-would-have-to-be-true paragraph and the one-family note
are present at `consensus.md:287-307`.

I authored the adopted `:251` text and the opencode row wording; per §15.1 I issue no verdict
on their quality. This signoff is a signoff decision, not a self-verification.

#### Scope declaration

Read in full for this signoff: `COOPERATION.md` §15 (lines 1176-1316); `consensus.md` revision
2 (all 548 lines, including the four embedded revision-1 signoffs); my own revision-1 signoff
(embedded at `consensus.md:340-442`). Executed fresh: the repo-wide `rg` over all `*.go`
quoted above (zero hits, rc=1); read `internal/config/runtime.go:95-159`; read the
`parley-deck-skill/package.json` scripts block; `rg -n skillVersion` over `parley-deck-skill`;
the bounded `find` (maxdepth 5) plus an unbounded re-run over `AI_WORKSPACE` (23 deck configs
both times); the `while read`/`grep -l` pass (12 files) and the `python3` JSON comparison
(10/9/2) quoted above; the `igm-app` spot-check; `grep -rl` probes on this repository and on
the BYTE tree for the tooling-defect claim. **Not done:** re-ran no git diffs; ran no test
suites and no binaries; did not re-read `runner.go:1094-1108` for this signoff (that claim
rests on the three participants' round readings, and the draft's blind-spot paragraph
discloses the shared-tool caveat, `consensus.md:298-302`); did not read decks outside the 23
measured files beyond the quoted JSON passes. I ran no git write commands and edited no file
other than this signoff.
