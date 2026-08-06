---
agent: kimi-1
idea: skill-sync-cli-1-39
round: 1
date: 2026-08-06
---

## Summary

**Ruling on the facilitator's claim: CONFIRMED (PRIMARY)** — the skill's "Generic CLI
Invocation Contract" and its JSON config shape do teach a two-list model the CLI does not have,
and `AutonomousWrite.Args` is indeed never consumed at launch. One refinement, not a reversal:
the precise shared flaw is *where validation points* (the declaration, not the effective argv),
and the practical danger is sharper than the claim states — the CLI does not parse the skill's
JSON files at all, so the two lists must be hand-flattened into one `headless_args` with zero
guidance. Full evidence in "Proposed approach", item A.

Per-item positions, argued below:

- **opencode autonomous-write row: carry.** It changes the flags a facilitator passes and the
  prompt-delivery mode (argv, not stdin).
- **kimi row: already covered** — SKILL.md:249 already states `-p` is the only autonomous
  headless shape and `--yolo`/`--auto` cannot combine with it. Matches the shipped spec.
- **Fail-closed wording: carry**, generalized to the effective-argv check; exact sentences
  proposed in item C. The existing confinement sentence stays.
- **"Source of truth is the spec's `autonomous_write` field" (SKILL.md:251): carry a rewrite.**
  1.39.0 inverts it; the brief missed it.
- **`compatibility.json`: bump plus one assertion in the existing test suite.** Bump alone
  re-drifts; the mechanism of the silence is proven below (item D).
- **ACP-catalog history, the hermes incident narrative, display internals, CLI struct names:
  leave out** — none changes what a facilitator or participant does.
- **Bundled `references/COOPERATION.md`: untouched.** I did not re-verify the currency diff
  (owner: facilitator); no position I take depends on it, and the byte-identical constraint
  binds regardless.

**Scope checked (§15 declaration).** Read in full: `00-prompt.md`; `COOPERATION.md` §15
(lines 1176–1316); `SKILL.md` (all 840 lines); `references/compatibility.json`;
`references/WORKED_EXAMPLES.md`; `package.json`; `CHANGELOG.md` v1.39.0 entry;
`internal/agents/discover.go` (all 694 lines); `internal/config/runtime.go` (agentOverride
struct 97–132, configLayers 137–154, applyOverride 514–662); `internal/runner/runner.go`
1094–1134 (`buildAgentInvocation`). Repo-wide greps in the CLI: `AutonomousWrite`,
`AutonomousEffective`, `MissingFrom`, `HeadlessArgs`, `writeModeArgs|write_mode_args|
WriteModeArgs`, `PrintRuntimeMatrix`, `row.Auto`. Repo-wide grep in the skill repo:
`skillVersion`, `compatibility.json` consumers. Executed: `shasum -a 256
skills/parley-deck/references/compatibility.json` against the hash pinned in
`parley-addon.json`. Not done: running the parley CLI, running the skill test suite, probing
live `kimi`/`opencode` binaries, or independently diffing v1.38.0..v1.39.0. Disclosure: the
repo-wide greps incidentally returned a few lines of another participant's round-01 file
(grep scans the whole tree); I did not open or read any round-01 file, and nothing below
derives from them.

## Proposed approach

### A. Ruling on the facilitator's claim — CONFIRMED (PRIMARY), with one refinement

The claim has four checkable parts. Each verdict is PRIMARY; I executed the greps and read
the files myself.

1. **"The launch uses `HeadlessArgs` alone" — CONFIRMED (PRIMARY).**
   `internal/runner/runner.go:1097-1108`, `buildAgentInvocation` — the launch argv is built
   exclusively from `agent.HeadlessArgs` with `{root}`/`{prompt}` substitution:
   ```go
   args = make([]string, 0, len(agent.HeadlessArgs))
   for _, arg := range agent.HeadlessArgs {
       switch arg {
       case "{root}":  args = append(args, root)
       case "{prompt}": args = append(args, prompt)
       default:         args = append(args, arg)
       }
   }
   ```
   No other arg source exists in the launch path.

2. **"`AutonomousWrite.Args` is read only by the two AUTO display sites and tests" —
   CONFIRMED (PRIMARY).** Repo-wide grep for `AutonomousWrite|AutonomousEffective|
   MissingFrom`: `.Args` is consumed only inside `MissingFrom` (discover.go:115–130), which
   feeds `AutonomousEffective` (discover.go:138–140). Its non-test consumers are exactly two
   display surfaces: the `agents list`/`discover` runtime matrix (`AUTO` column,
   discover.go:506, plus the missing-args warning line, discover.go:538–541) and the roster
   table (`row.Auto = spec.AutonomousEffective()`, internal/app/roster.go:155, commented
   "Fail-closed"). Tests: `autonomous_test.go`, `acp_specs_test.go`, `runtime_test.go`.
   No launch-path consumer exists.

3. **"A config layer replaces `headless_args` wholesale and never touches
   `autonomous_write`" — CONFIRMED (PRIMARY).** `agentOverride` (runtime.go:97–132) has no
   autonomous-write field at all; `applyOverride` (runtime.go:514–662) never references
   `spec.AutonomousWrite`; the replacement at runtime.go:542–545 is wholesale:
   `spec.HeadlessArgs = expandSlice(override.HeadlessArgs, root, tempdir)`.

4. **"No `writeModeArgs` exists in the CLI" — CONFIRMED (PRIMARY).** Grep for
   `writeModeArgs|write_mode_args|WriteModeArgs` across the CLI repo: zero hits in Go source
   (only deck log files quoting the skill text). Meanwhile the skill documents
   `writeModeArgs` as a distinct config field (SKILL.md:363) and as a distinct
   launch-construction step (SKILL.md:811, step 3, after step 2 "Add `headlessArgs`").

**Refinement (where I'd scope the claim's interpretation).** The skill's two-list JSON is
not the defect's *literal* mechanism — the live defect was a TOML `headless_args` override
against a built-in-only declaration no config layer can touch. And the skill's JSON files
are not what the CLI parses: `PARLEY_HEADLESS_AGENT_CONFIG` and every config layer are
TOML (`toml.Unmarshal`, runtime.go:440), and `configLayers` (runtime.go:137–154) reads
`parley-deck/agents.toml` / `agents.local.toml` — never the skill's
`parley-deck/meta/headless-agents.local.json`. So the reproduction path is not "the skill's
field leaks into the CLI"; it is that a facilitator holding the skill's two-list mental
model must **hand-flatten two lists into one `headless_args` array with no stated
invariant**, while the skill's own text points validation at the wrong target: "The source
of truth for each agent's mode is the spec's `autonomous_write` field" (SKILL.md:251).
That is the mental model that made the defect invisible — trust the declaration, not the
effective launch — and on that substance the facilitator is right. The claim survives
scrutiny; it should be carried into the fix as "validate the effective argv," not narrated
as "the skill's field caused the bug."

### B. Per-item positions

1. **opencode row in the Autonomous Execution table — carry.** It changes what a
   facilitator DOES: which command and flags to pass (`opencode run --auto <prompt>`,
   spec at discover.go:344–375) and how to deliver the prompt (argv positional, not stdin —
   `PromptMode: PromptArg`, discover.go:359; the spec note states "the message is argv, not
   stdin", discover.go:372). One motivating clause is worth carrying because it instructs
   future behavior: `opencode run` writes unattended even without `--auto`; `--auto` is
   passed explicitly because an implicit default is what a vendor may change between
   versions (spec comment, discover.go:350–354). That clause tells a facilitator *not to
   drop the flag when editing config* — it is operational, not changelog.

2. **kimi row — already covered.** SKILL.md:249: "plain `-p` — its print mode already
   auto-approves in-workspace writes. NOTE: `--yolo`/`--auto` are mutually exclusive with
   `-p`, so `-p` IS kimi's yolo-equivalent." This matches the shipped spec (discover.go:
   312–343, `AutonomousWrite{Mode: "prompt", Args: ["-p"]}`). No 1.39.0 edit needed.

3. **Fail-closed coverage — carry the generalization** (exact wording in item C). The
   current sentence covers only workspace confinement (SKILL.md:251); the new case — a
   declared mode whose enabling flag the effective launch never passes — is now the
   demonstrated failure class and must be named.

4. **"Source of truth" sentence — carry a rewrite** (see item E; the brief missed it).

5. **Leave out, with reasons:**
   - *ACP-catalog-stub history / "both keep ACP as an alternative launch mode"* — the
     skill's contract is one-shot headless invocation (SKILL.md:805) and mentions ACP
     nowhere; launch-mode internals change nothing a facilitator does.
   - *The hermes incident narrative* — changelog material. The rule it motivates is what
     carries; the war story does not change behavior. (One subordinate "a config override
     can silently drop the flag" clause stays, because it tells the facilitator what to
     look for.)
   - *kimi's installer not adding the binary to PATH* (new note, discover.go:339) — the
     remedy is already generic in the skill: Startup Flow step 6, "verify … with
     `command -v <cli>` or an explicit configured path" (SKILL.md:152). A vendor-specific
     install note does not belong in vendor-neutral instruction text.
   - *opencode telemetry description, `headless:`-line/AUTO display internals,
     `AutonomousWrite.MissingFrom`, spec structs* — CLI internals; the brief's non-goal.
   - *Bundled `references/COOPERATION.md`* — the brief asserts it current (facilitator
     owns that claim; I did not re-run the diff) and the constraint requires
     byte-identical. Do not touch.

### C. The fail-closed wording — exact proposal

Replace the two sentences at SKILL.md:251 ("The source of truth … fail-closed) …") with:

> The source of truth for an agent's autonomous capability is the **effective launch
> argv**, not the declared mode. Before trusting that a participant can write its own
> artifact, confirm that the flags enabling its auto-approve mode are present in the exact
> command that will run — a config override can replace the launch arguments and silently
> drop them. When the parley CLI drives the agents, `agents list` performs this check and
> fails closed: it reports `AUTO=no`, names the missing flags, and prints the effective
> argv; `parley roster show` applies the same reading. If the enabling flags are absent,
> treat the agent as non-autonomous until the config is fixed — do not assume the declared
> mode applies. If workspace confinement cannot be demonstrated for an agent, treat its
> autonomous bit as unset (fail-closed) rather than escalating to a full-filesystem
> bypass.

Design notes: the check is stated vendor-neutrally first (inspect the effective command),
with the CLI surfaces named second — the skill already names CLI commands where they are
the check (`parley roster show`, SKILL.md:257), so this sets no new precedent. The rule
covers kimi without naming it: `-p` is its enabling flag, so "flags enabling its
auto-approve mode present in the command" applies verbatim. Keep the two-list JSON shape
(`headlessArgs` + `writeModeArgs`) — the defect is the missing invariant and the missing
check, not the field split; deleting a documented field churns the shape without adding
the guard. Instead, add one line to the invocation contract after step 3 (SKILL.md:811):

> The assembled command — not any single config field — must contain every write-mode
> flag; after any config change, verify the effective argv before trusting the launch to
> write unattended.

### D. `compatibility.json` — bump plus one assertion; the silence mechanism is proven

Facts (PRIMARY): `references/compatibility.json:4` says `"skillVersion": "1.4.3"`;
`package.json:3` says `"version": "2.3.0"`. Repo-wide grep for `skillVersion` /
`compatibility.json` consumers across the skill repo: no code in `bin/`, `lib/`, or
`scripts/` reads the `skillVersion` value — the installer references the file only as a
required payload entry (installer.js:130, 136, 2379–2383). The one existing guard,
`parley-addon.json`, pins a sha256 of the file — and I executed the comparison:

```text
$ shasum -a 256 skills/parley-deck/references/compatibility.json
b2465b207a16ac4241040e05c64b9e9f4ea79c1c780b63149c4e8c0d18267abd
parley-addon.json:9: "references/compatibility.json": "sha256:b2465b207a16ac4241040e05c64b9e9f4ea79c1c780b63149c4e8c0d18267abd"
```

The pinned hash **matches the stale content**. The existing guard verifies byte integrity,
not version freshness — that is exactly why four releases of drift passed silently.

Position: bumping alone re-drifts — `skillVersion` is a hand-maintained duplicate of
`package.json`'s `version`, and hand-maintained duplicates drift; this one demonstrably
did. The proportionate guard needs no new tooling: **one assertion in the existing
`node --test` suite** (e.g. in `test/manifest-coverage.test.js`, which already gates
freshness of a generated artifact and already runs on `npm test` and `prepack`) that
`references/compatibility.json`'s `skillVersion` === `package.json`'s `version`. That
converts the owner's standing rule ("a CLI release always ships the skill in the same
turn") into a mechanical fact for the skill's own version field. The rejected alternative:
delete `skillVersion` (nothing consumes it) — but the file is a published schema
(`schemaVersion: 1`) and external tooling may read it; keeping the field correct is safer
than removing it. Note the bump target is whatever version *this* release ships, not
necessarily 2.3.0 — which is precisely why the assertion, not the bump, is the fix.

### E. What 1.39.0 invalidates that the brief missed

1. **SKILL.md:251, first sentence: "The source of truth for each agent's mode is the
   spec's `autonomous_write` field."** 1.39.0's entire fail-closed change inverts this:
   the declaration is now *checked against* the effective argv and fails closed. Left
   standing, this sentence teaches exactly the pre-1.39.0 mental model inside the fixed
   text. The brief lists the fail-closed sentence's *second* half (confinement) but not
   this one. Rewritten in item C.
2. **`references/WORKED_EXAMPLES.md:33` carries the same two-list shape** —
   `"writeModeArgs": ["--workspace-write"]` — and the brief's stale list names only
   SKILL.md. Whatever lands in the JSON shape must be mirrored there; it is labeled
   non-authoritative but is the file a facilitator copies from. Secondary defect in the
   same line: `--workspace-write` as a single flag matches no shipped adapter (codex's is
   the two tokens `--sandbox workspace-write`, discover.go:196); a copied example should
   not model a flag shape no CLI accepts.
3. **Nothing to remove about ACP-only status.** Grep of SKILL.md: the skill never claimed
   kimi/opencode were ACP-only stubs, so the gap is purely additive (the opencode row).
   Stated so the drafter does not hunt for text that does not exist.

## Concerns / open questions

- **Adjacent, pre-existing, explicitly out of scope for this idea:** the skill documents
  `PARLEY_HEADLESS_AGENT_CONFIG` as "a JSON config file" (SKILL.md:142 and 341), but the
  CLI parses that env var's file as TOML (runtime.go:440), and the CLI never reads the
  skill's `parley-deck/meta/headless-agents.local.json` path (configLayers,
  runtime.go:137–154). 1.39.0 did not invalidate this — it predates the release — but it
  is the same two-worlds mismatch one layer down, and it will confuse the next facilitator
  who hands the CLI the file the skill describes. Worth an inbox note or a follow-up idea;
  it should not enlarge this one's diff.
- **Naming CLI surfaces in vendor-neutral text:** item C names `agents list` and
  `parley roster show`. Precedent exists (SKILL.md:257), and the alternative — a check
  described with no tool that performs it — is what the pre-1.39.0 text had. If consensus
  prefers fully vendor-neutral wording, the generic first two sentences of the proposal
  stand alone; the CLI sentence can be dropped without losing the invariant.
- **Guard-test coupling:** an equality assertion means any `package.json` bump without the
  `compatibility.json` bump fails `npm test`/`prepack`. That is the intended fail-closed
  behavior, but it does make the two files a single atomic edit forever; if the project
  ever wants them decoupled, the test should be deleted in the same change, not weakened.

## Risks

- **Over-scoping the diff.** The tempting move is rewriting the whole Autonomous Execution
  section or the whole invocation contract. The positions above are deliberately four
  surgical edits (one table row, one sentence pair, one contract line, one mirrored
  example) plus the version guard; anything more reopens text that 1.39.0 did not
  invalidate.
- **Wording rot.** If the CLI later renames `agents list`'s AUTO column or the
  `headless:` line, the skill sentence goes stale silently — the same failure class this
  idea is fixing. Mitigation: the generic invariant is stated first and carries the rule;
  the tool names are illustrative. Residual risk accepted; a skill that names no tool
  teaches no check.
- **My verdict's refinement could be misread as a partial refutation.** It is not: every
  falsifiable part of the facilitator's claim checked out against primary sources. The
  refinement only narrows the *mechanism story* so the fix targets the invariant and the
  validation target, not a field rename that would leave the hole open under a new name.
- **Unverified-by-me dependency:** the claim that bundled `references/COOPERATION.md` is
  current is the facilitator's (verdict quoted in the brief: empty `git diff
  v1.38.0..v1.39.0` and a clean normalized diff). I did not re-run it. My "do not touch"
  position rests on the byte-identical constraint, not on that claim, so a surprise there
  would change the release mechanics, not any position in this file.
