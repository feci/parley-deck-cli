---
agent: zcode-1
idea: protocol-and-skill-audit
review-round: 1
reviewed-commit: 0bb99031aa1ddc1027605424848daa3a696e9f00
date: 2026-08-21
---

## Summary

I built from source (`go build ./cmd/parley`) and ran both suites myself: `go test ./...` green
across 26 packages, `npm test` in ../parley-deck-skill 388 pass / 0 fail (PRIMARY). I also built
the PRE-FIX binary (from d9aa26d, the last commit before the fix batch) and diffed its verdicts
against the new binary across every consensus in this deck — 87 design consensuses and 65 review
consensuses (PRIMARY).

Most fixes hold under attack: the F20 signoff-position rule and F21 slug check changed zero design
verdicts deck-wide; the F1/F17/F18 content floors reject no in-flight artifact; F13's new refusal
breaks no idea in this deck or the five fleet decks; the F24 unknown-freshness gate is honest
friction with clear instructions (verified on a freshly initialized scratch deck); the fixture
edits in existing tests encode the new contract rather than weakening assertions; the five
deferrals all have factual, verifiable reasons.

But three fixes are broken in the too-strict direction, and all three share one shape: a gate was
tightened without updating the thing that produces the artifact the gate checks.

1. **F3 (implementer exclusion) malforms every review consensus in which the implementer signed**
   — which is what COOPERATION.md:591 tells it to do. 23+ real consensuses flipped ready/partial →
   malformed, and the auto-driver's Phase 8 escalates on malformed, so a protocol-compliant
   artifact can never complete.
2. **The driver's FINAL drafting prompt still asks for one section while the driver's own new
   FINAL gate demands seven plus a matching slug** — auto-drive cannot close any idea whose
   drafter follows the prompt.
3. **The pipeline auto-drive ignores the two-step finalize contract** — it prints "finalized" on
   the scaffold step and its block-completion predicate accepts the scaffold (which carries
   `status: final`), so F5's close-around-a-scaffold defect survives intact in that path.

## Refutation attempts

What I tried that FAILED to break the fixes (all PRIMARY unless noted):

- **Old-vs-new binary diff, design consensuses (87 ideas).** Ran `consensus status --json` with
  both binaries on every idea; parsed triage, signoff sets, missing, errors. Zero differences.
  The F20 position rule, F21 frontmatter check, and F4 do not misread any real design consensus.
- **F20 against real signoff shapes.** The implementation notes claim 405/405 real signoffs sit at
  or after a `## Signoff*` heading; the binary diff independently confirms no signoff set changed.
  Heading variants (`## Sign-offs`, `### Signoffs`) would zero out all signoffs, but no file in
  the deck uses them.
- **F1 floor (blank / no-heading).** Scanned every round artifact in the deck with the same
  logic: 0 blank, 4 headingless (all early `hermes.md` files, all in closed ideas). No open idea
  (the 8 non-terminal ideas checked individually) is blocked by the floor.
- **F17/F18.** Scanned all round-01 artifacts for empty required sections (0 hits) and all review
  artifacts for missing `reviewed-commit` (191, all historical; the review prompt now asks for it,
  and this idea's own review round carries it). No in-flight idea is rejected.
- **F13's new refusal.** Scanned every 00-prompt.md in this deck and in parley-deck-test,
  wt-editor-composer, wt-learn-playbooks, wt-roster-presets, wt-round-summary: no idea declares
  explicit `track: standard` together with `auto_implement` or `strict_gate`. The refusal fires
  only on new runs and matches `classify`'s existing verdict.
- **F22/F5 retroactivity.** 74 of 78 historical FINAL.md files would fail the new seven-section
  gate — but all belong to closed ideas; `Finalize` refuses `status: final` early and the driver
  gates only newly drafted finals. No retroactive breakage.
- **F24 on a fresh deck.** Initialized a scratch deck and ran preflight: exit 3 with an
  `[unknown-freshness]` gate and an explicit `--yes` confirm command. One-time, honest friction
  ("in sync" was previously asserted from two missing hashes); not a blocker.
- **Fixture-edit audit.** Read every edit to existing tests (app_test.go, consensus_test.go,
  driver/consensus_test.go, track_test.go, drift_test.go, phase58_test.go, phase58_le_test.go).
  All update fixtures to the new contract (seven sections / `reviewed-commit` / renumbered
  two-step expectations) while keeping the assertions; `TestPolicyForAbsentIsLegacy` re-pins the
  deferred F14 behaviour explicitly. None weakened a test to make a change pass.
- **New tests are not tautological.** They call the real functions (`missingRoundArtifacts`,
  `parseDocument`, `classifyAndSyncFreshness`, `PolicyFor`, `NormalizeStrict`,
  `nextCrossReviewRound`, `runaction.Command`) and pair every negative case with a positive
  control (`TestWellFormedArtifactCounts`, `TestCompleteFinalIsAccepted`, `TestSectionsMayBeNA`,
  `TestSignoffsUnderLaterHeadingsStillCount`).
- **Deferral facts.** kimi-1/F1: confirmed the installed core (~/.claude/skills/parley-deck)
  carries `plugin.json` and `gemini-extension.json` that `skills/parley-deck/parley-addon.json`'s
  file map does not describe — byte-verifying today would report them unexpected. codex-1/F14:
  confirmed in track.go that `ApplyOverrides` sets exact values (a configured MaxFixupCycles=5
  would be forced to 2), so a real default needs per-knob clamp semantics. F6/F8 are protocol /
  feature design decisions. None of the five reasons is an excuse for an easy fix.
- **Protocol text fixes.** zcode-1/F6's new citation `2026-06-02T12-07-14-meta-protocol-ch` names
  a real idea dir; present in all three copies (defaults, deck, skill snapshot). F1/F7/F8/F11
  wording verified in parley-deck/COOPERATION.md (§2 resolution order with `inherited-roster`,
  Quickstart §15/§10, §3 `agents.toml`, §11.B rewording with the correction note).
- **Hermes probe (claude-1/F2).** `parley agents verify --full --agent hermes --yes` →
  `hermes: headless probe passed`, probe dir under the repo (PRIMARY).

## Findings

### [MAJOR] F3's implementer exclusion malforms every review consensus the implementer signed — which the protocol tells it to append

What is wrong: `Status(root, slug, true)` now passes `expectedRoundParticipants(...)` (implementer
removed) as the `participants` argument to `validateDocument` (internal/consensus/consensus.go:113).
But that list is used for TWO different checks: the quorum (`Missing`) AND signer validity — any
signoff from an agent not in the list becomes `line N: unknown participant X` → `Errors` → triage
`malformed`. COOPERATION.md:591's Phase-7 template says the opposite of "not in the list":
`<!-- Each active participant (implementer included) APPENDS their signoff block. -->`.

Why it matters (PRIMARY): I diffed old vs new binaries over all 65 review consensuses. Six flipped
partial→ready (the intended fix) — and **23 flipped ready/partial→malformed** solely because the
implementer's signoff is present, e.g. tui-editor-composer (implementer claude-1 signed with note
"Implementer. All agreed fixes applied", reviewer codex-1 accepted): old `ready`, new `malformed —
line 27: unknown participant claude-1`. Same for driver-impl-phase, verification-honesty,
named-roster-presets, roster-operations-standard, loop-budgets, and ~18 more. These artifacts were
written exactly per the protocol's own template. Worse, the consumer that matters is the driver:
`driver/impl.go`'s Phase-8 path escalates with "review consensus triage=malformed; cannot
complete" (SECONDARY — read from internal/driver/impl.go:190-210 and
internal/app/driver_impl.go:328-330), so any mixed manual/auto flow in which the implementer signs
as instructed can never close. The error text is also false — claude-1 IS a participant.

Why the tests missed it: `TestReviewConsensusDoesNotAwaitTheImplementersSignoff` only exercises
the shape where the implementer did NOT sign.

Concrete fix: in `validateDocument`, validate signers against the FULL participant list (keeping
the "unknown participant" error for genuinely foreign agents) and compute `Missing`/quorum against
the reduced reviewer list. Pass both lists (or a struct) from `Status`. Add a test where the
implementer signs and the consensus is `ready` with zero missing reviewers.

### [MAJOR] The driver's FINAL drafting prompt cannot satisfy the driver's own new FINAL gate

What is wrong: `buildFinalDraftPrompt` (internal/app/driver_consensus.go:131) instructs the
drafter to write FINAL.md with frontmatter "that includes `status: final`" and "a populated
section: `## Final plan / specification` … no placeholders" — one section, and no `idea:` slug
asked. The same driver then validates with `finalScaffoldReason` (internal/driver/consensus.go),
which after F22 requires ALL seven protocol sections and, after the slug fix, a frontmatter
`idea:` equal to the directory name.

Why it matters: in `advanceConsensus` (internal/driver/consensus.go:43-58) a FINAL written exactly
as the prompt asks fails with "FINAL.md is not acceptable after drafting: missing required
section(s): Purpose / user-visible outcome, …", and if the drafter only keeps `status: final` it
fails `finalSlugMismatch` first. Auto-drive escalates on its own instruction-following output;
the idea can only close if the drafter volunteers structure the prompt never named. This is the
audit's own defect class — the implementer fixed it for F10 (review-consensus template) and F18
(review prompt), and wrote in IMPLEMENTATION.md that F18's lesson was "the prompt was fixed
first"; this prompt was not. No test covers `buildFinalDraftPrompt`. (SECONDARY — code reading;
the mismatch is deterministic.)

Concrete fix: make the prompt name all of `protocol.RequiredFinalSections` (content may be `N/A`),
require `idea: <slug>` in the frontmatter, and add a test asserting the prompt contains every
required section heading.

### [MAJOR] The pipeline auto-drive ignores the two-step finalize: prints "finalized" on the scaffold step and completes blocks around unwritten FINAL.md

What is wrong: `autoDriveDeliberationBlock` (internal/app/pipeline_cmd.go:736) calls
`consensus.Finalize(...)` once, checks only `err`, then prints `auto: block %q finalized.` After
codex-1/F5 the first call writes the scaffold and returns `Scaffolded: true` with NO error — no
second finalize ever runs in this path and no drafter fills the scaffold.

Why it matters: the printed success statement is now false (the idea stays open at
`status: consensus`). Worse, block completion is decided by `blockCompleteFunc` →
`isFinalized(content)`, which just greps `status: final` — and `finalTemplate` writes
`status: final` into the scaffold's own frontmatter. So the wave logic marks the block complete,
unblocks downstream blocks, and the pipeline advances around an unwritten FINAL.md — F5's exact
defect ("closed around an empty scaffold while reporting success"), transposed from idea status
to pipeline block status. (SECONDARY — code reading of pipeline_cmd.go:736, 1254-1278, 1323-1330;
the scaffold's `status: final` line is visible in `finalTemplate`.)

Concrete fix: after `Finalize`, branch on `summary.Scaffolded`: stop with a message telling the
operator to fill FINAL.md and re-run (mirroring the CLI's two truthful sentences). Independently,
make `isFinalized` (or the block predicate) reject bodies for which
`protocol.FinalIsScaffold(body) != ""`, since the scaffold's `status: final` frontmatter is
indistinguishable from a written final by that grep. Consider having the scaffold not claim
`status: final` until it is written (that is the root of both shapes of the bug).

### [MINOR] IMPLEMENTATION.md accounting gap: codex-1/F12 and F23 appear in neither the fixed nor the deferred list

What is wrong: the fix list and the "NOT fixed, with the reason" section account for the findings
that were addressed, but codex-1/F12 (init reports completed bootstrap without the mandatory
roster/model/effort confirmation) and codex-1/F23 (implementation gate accepts `status: banana`
and an empty one-heading artifact) appear nowhere. Both were PARTIAL/contested — but so was F15,
which WAS fixed, so PARTIAL status alone does not explain the split, and nothing records the
choice. I verified F23's validator is unchanged: `ValidateImplementationArtifact`
(internal/runner/phase58.go:394-413) still accepts any non-empty status and only greps for the
`## Summary of work` substring (SECONDARY).

Why it matters: this review's own audit trail is the deliverable; a reader of IMPLEMENTATION.md
cannot tell whether F12/F23 were forgotten or deliberately left. F23 is the same
presence-only-gate class as the fixed F17/F18/F22 and guards Phase 5→6.

Concrete fix: add both to the deferral list with one-line reasons (or fix F23's status
enumeration, which is a small allowlist check).

### [MINOR] F14's repair missed two provably-closed ideas still declaring `status: implementation`

What is wrong: the repair (6df4703) touched exactly the 20 closed ideas that carry a FINAL.md.
`kimi-opencode-full-adapters` and `rho-retro-tooling` have no FINAL.md, but their IMPLEMENTATION.md
says `status: complete` with closed review consensuses and merged commits
("kimi-opencode-full-adapters: IMPLEMENTATION.md fix-up cycle 1 — complete",
"rho-retro-tooling: IMPLEMENTATION.md — complete" in git log) — both still feed §6 rule 5's
stale-round guard false liveness data, which is what F14 was about (PRIMARY: read the files and
git log).

Concrete fix: set both 00-prompt.md statuses to the terminal value chosen for the other repairs
(`final`), or record why implementation-complete ideas were scoped out.

### [MINOR] The seven-section list now exists twice: the driver redefines what protocol/ was created to own

What is wrong: `protocol.RequiredFinalSections` (internal/protocol/finalsections.go) exists,
per its own comment, because "two independent gates need the same list … They disagreed before,
which is how an idea could be closed around a scaffold by one path and rejected by the other."
Yet `internal/driver/consensus.go:166` defines its own `requiredFinalSections` copy instead of
importing protocol's — while the driver package already imports internal/protocol
(transport.go:8). The driver's `finalScaffoldReason` also keeps its own inline placeholder list
rather than `protocol.FinalScaffoldPlaceholders`. The lists agree today (verified by reading
both), but the drift risk the fix's own rationale names is rebuilt into the tree.

Concrete fix: delete the driver-local copy and use `protocol.RequiredFinalSections` /
`protocol.MissingFinalSections` (and the placeholder list) in `finalScaffoldReason`.

### [MINOR] F16's `responding-to` gate rejects the multiline YAML list form that real deck artifacts use

What is wrong: `hasRespondingTo` (internal/driver/driver.go:481) reads only the same-line value of
the frontmatter key. Six real artifacts (skills-cli-install-path/review/round-17..21) write the
key as a multiline YAML list (`responding-to:` followed by `  - agent/review/round-NN` lines);
the parser returns `""` for these, so the gate answers "names nobody" and would mark a round
incomplete (PRIMARY: read the files and the parser).

Why it matters only mildly: those six are review-round artifacts and the gate fires on design
round-02+ under the driver, so nothing live is currently rejected — the fix's measurement
("162 of 162 design artifacts carry names") holds for the population it measured. But the
multiline form is legitimate YAML an agent can naturally produce in a design round, and refusing
it produces a confusing error ("must name somebody" when it names three).

Concrete fix: when the same-line value is empty, continue scanning frontmatter lines and treat
`- <entry>` continuation lines as list entries.

### [NIT] F4's `--by` check is stricter than Phase 4's drafter rule

Phase 4 says the FINAL drafter is "the initiator or an agreed participant". A human initiator who
is not in `participants:` (agents only) is now refused by `--by <name>`, while omitting `--by`
silently records the anonymous `user`. Consider allowing the workspace owner/initiator identity,
or documenting that human initiators must omit `--by`. (SECONDARY — protocol text vs
consensus.go:216-221.)

### [NIT] kimi-1/F2 preserves two NAMED keys, not "the parts it does not understand"

The commit's principle ("a tool that rewrites a shared file must preserve the parts it does not
understand") is implemented as `PRESERVED_PROJECT_METADATA_KEYS = ["protocolRole", "created"]` in
lib/installer.js. That covers exactly what the CLI writes today (verified against
`writeInitVersionMeta` and the CLI's versionMeta fields), but any future CLI-owned key gets
silently deleted again. Either invert to "preserve every key not in the owned set" or note the
allowlist's maintenance duty next to it. (SECONDARY.)

## Open questions

- The driver's own review-consensus prompt (internal/runner/phase58.go:351-383) still shows
  `outstanding_agreed_fixes: 0` in its "Required file shape" example, while the manual template
  (after F10) deliberately seeds a placeholder so the count is "never a silent 0". The driver
  prompt does carry an explicit "set to 0 only when nothing remains to fix" rule, so this may be
  intentional — but the two templates now teach different defaults for the same field.
- `runaction.Command`'s no-run-ID fallback prints `parley continue <idea-slug>`, which errors with
  "no runs found" for an idea that was driven manually without `parley run` (PRIMARY: tried on a
  scratch deck). The run-scoped form (RunID present) is the one actually surfaced in practice; is
  the slug fallback ever reachable for a runless idea, and if so should it name a command that
  works?
- The F3 fix's `resolveImplementer` reads IMPLEMENTATION.md `implementer`, else FINAL.md
  `implementer`/`drafted-by`, and ignores names that are not participants (fails closed to the
  full list). That ordering matches Phase 5's "default implementer is the FINAL drafter" — worth
  one sentence in the code or protocol confirming the fallback is a deliberate reading of that
  rule, since it decides who may be silently dropped from a quorum.
