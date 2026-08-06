---
idea: roster-operations-standard
phase: 7 — review consensus signoff
agent: hermes-1
date: 2026-08-06
verdict: ACCEPT
---
# Signoff — hermes-1

## Verdict

ACCEPT. The draft reproduces every finding I raised (and several I missed), resolves the one
reviewer-vs-reviewer conflict against the drafter's own earlier work by reading the normative
artifact rather than counting, and escalates A1 against the implementer's own implementation.
Nothing of mine was dropped, downgraded, or misrepresented.

## Answers to 1-6

### 1. Are my findings accounted for?

Every finding I filed in `review/round-01/hermes-1.md` is either agreed in §2 or carried as a
residual there. I traced each one:

- My CRITICAL G1 (snapshot written but not consumed; `continueAuto` still calls
  `discoverConfigured`, app.go:1152/1163; acceptance test absent) — `PRIMARY` at v1.40.0. The
  v1.40.1 fix landed consumption + a test (kimi-1 C1 status, `SECONDARY` from kimi-1's review
  §1/codex-1 CRITICAL-1). The residual gaps I could not have seen at v1.40.0 — auto-args not
  pinned (A4) and per-adapter keying collapsing per-roster-ID pins (A5) — are filed as MAJOR by
  the drafter from kimi-1 M1/M2. I accept the deferral shape: the gate's core (consume the
  snapshot) shipped; the remaining two are hardening, not a gate re-violation.
- My CRITICAL G4 (no §2 generator; protocol text claims generation that doesn't exist; drift
  guard masks the live/embedded divergence by emptying both bodies) — `PRIMARY` at v1.40.0. The
  generator shipped in v1.40.1 (`roster_render.go`, `SECONDARY` from kimi-1 §1 C2 / codex-1
  MAJOR). Residual A2 (§2-only IDs silently dropped and erased by render) is the unimplemented
  conflict rule from the normative field table (`consensus.md:355`), which I had not separately
  raised; the drafter adopted it from kimi-1 M3. Correctly MAJOR.
- My MAJOR D5 (`--yes` alone not refused for membership; `roster set ... --yes` wrote a new
  block, roster_set.go:126-139) — `PRIMARY` at v1.40.0. The v1.40.1 fix added
  `--confirm-breaking`, but the gate is bypassed by shape: a block created with only
  `--model`/`--effort`/`--speed`/`--state` slips it (A3). I verified this myself now —
  `membershipChange` (`roster_set.go:236-248`) gates on `strings.Contains(c, "+ adapter = ")`,
  so a new `[roster.<id>]` block appended with no `adapter` key is not detected as a membership
  change. `PRIMARY`. A3 is the right fix site: gate on "block did not exist before", not on
  which key was written.
- My MAJOR D7 (legacy normalizer absent; `applyOverride` does wholesale replacement,
  runtime.go:622-625) — `PRIMARY` at v1.40.0 and unchanged at v1.40.1. Agreed as A8, MAJOR, all
  three reviewers independently. Correct.
- My MINOR D5 (`--explain`/`--all` not implemented; flag parser defines neither,
  roster.go:87-100) — `PRIMARY` at v1.40.0. Upgraded to MAJOR as A6 via C-2. I accept the
  upgrade: I scored only the *missing* flags (correctly MINOR — a rejected flag is a visible
  failure), while kimi-1/codex-1 additionally reproduced `--scope machine` parsing, being
  advertised in help, and silently changing nothing. A silently wrong answer is a different
  defect class from an absent one. The severity resolution in C-2 is by defect class, not by
  count, which is the §15.3-correct way to do it.
- My MINOR D6 (`sessions inspect` never reports `stale-snapshot`; only an aspirational comment
  at manifest.go:56) — `PRIMARY` at v1.40.0. Upgraded to MAJOR as A10. I accept: D6 is a
  two-part decision (persist+consume the snapshot; report staleness) and only the first part
  shipped. The frozen STATUS vocabulary ships a code no surface can emit.
- My MINOR (three COOPERATION.md copies not identical in §2; live has 4 rows, embedded/skill
  empty; drift guard masks it) — `PRIMARY` at v1.40.0. Folded into A9, which is broader: the
  drafter corroborated it with codex-1's finding that every protocol surface still contains
  §2-as-a-store instructions ("mirrored in the §2 roster", "Fill in §2 roster", "Modify the
  active roster (§2)"). My finding was the symptom (divergent table content); A9 captures the
  cause (contradictory instructions remain). Correctly MAJOR.
- My NIT D1 (`agents list` not relabelled "adapter/runtime inventory — not the roster";
  runAgentsList unchanged) — `PRIMARY` at v1.40.0. Agreed as A15 (N-min-2). Correct.
- My NIT (docs contain zero `roster` mentions; `grep -n "roster" docs/cli-reference.md
  docs/agent-runtime-configuration.md` → 0) — `PRIMARY` at v1.40.0. Agreed as A13, split
  MAJOR-for-`--help` / MINOR-for-docs. Correct.

Nothing of mine was dropped or downgraded below what I filed. Two of my MINORs were *upgraded*
(A6, A10), which I accept on the defect-class reasoning in C-2 and the two-part-decision
reasoning for D6.

### 2. C-1: G5 changelog-format conflict — is the resolution sound?

Sound. I verified it myself (`PRIMARY`).

The §7 template (`COOPERATION.md:745-748`) requires four fields:
```
## YYYY-MM-DD — <short description>
Idea: ideas/meta-protocol-change-<topic>/
Drafted by: <agent-id>
Summary: <1–2 sentences>
```

The entry at `parley-deck/meta/protocol-changelog.md:119-139` uses:
- `## 2026-08-06 — §2 roster authority moves to parley-deck/agents.toml` — heading matches ✓
- `**Idea:** \`roster-operations-standard\`` — bold label, bare slug, no `ideas/.../` path ✗
- No `Drafted by:` line at all ✗
- `**Change:**` / `**Why:**` instead of `Summary:` ✗

kimi-1's "in §7 format" verdict (kimi-1 §5, G5: "Satisfied ... in §7 format") read the entry's
substantive completeness (it IS genuinely thorough — it names the idea, the track, the venue
deviation, and the one-off) rather than the template's field names. codex-1's "WRONG on format"
verdict (codex-1 MINOR-2) read the template. The template is a normative artifact in the
repository; it decides. The resolution does not rest on counting reviewers (2-vs-1 for codex-1's
reading would have been meaningless) — it rests on reading `COOPERATION.md:745-748` and comparing
field names. A11 is the correct fix: add `Idea:` as a path, `Drafted by:`, and a `Summary:` line;
keep the existing prose below. MINOR is the correct severity (format defect on an otherwise
substantively complete entry).

### 3. A1 escalation (membership layering) — right resolution, right scope?

Right resolution, right scope. I verified the mechanism myself (`PRIMARY`).

`RosterMembership` (`roster.go:643-657`) and `renderRosterTable` (`roster_render.go:29-30`) both
call `config.LoadRoster(root)` — the *layered* machine+deck view — while `rosterSync`
(`roster_sync.go:46-50`) operates on the deck file alone. So the release implements two
different meanings of "deck membership" in different verbs, exactly as kimi-1 M4 reported and
the consensus §2 A1 records. I could not have seen this at v1.40.0 because `RosterMembership`
and `roster render` are v1.40.1 additions; at my baseline the failure mode was the *opposite*
(selection on §2 while show read TOML). The v1.40.1 fix moved selection onto `RosterMembership`
but through the layered view, which swaps one two-sources-of-truth defect for another: a deck
declaring two members will run five (the three machine-only IDs leak in), and `roster render`
then commits that inheritance into `COOPERATION.md`, where it goes stale on the next machine
change — re-creating the drift vector inside committed files.

The adopted resolution (kimi-1's position) is the correct one and is consistent with the
ratified design: D9/FINAL decided "agents.toml is the deck authority" and `roster sync`'s rebase
model already treats the machine layer as seeding *values*, not membership. Making deck
membership = the deck file, with inherited rows marked for display only and `roster render`
refusing to bake an inherited roster into §2 without an explicit flag, closes the gap without
contradicting D9. The scope is right: it touches `RosterMembership`, `renderRosterTable`, the
participant-selection call sites, and render's behavior — the four places the layered view leaks
membership. CRITICAL is the correct severity: a deck that declares two participants will run
five, which can defeat the protocol's non-solo invariant (§1) by silently expanding quorum, and
render commits the leak into the file the drift guard watches.

The drafter escalating this against their own implementation, and adopting a reviewer's position
rather than their own earlier work, is exactly the §15.1 self-verdict boundary respected: the
drafter records the finding and the fix, the reviewers own the verdict.

### 4. DF-1 deferral (roster migrate contract deviations) — acceptable?

Acceptable, with the guard. I could not independently verify the fleet-run outcome (I did not
enumerate the 40-deck fleet — stated in my review §Scope); the "24 applied, 9 skipped, 3
unchanged, 0 failed, with per-deck backups and post-write validation with rollback" is
`SECONDARY` (kimi-1 §0 records the baseline shift; the drafter's consensus §3 DF-1 records the
counts). I accept the deferral reasoning on its own terms: kimi-1 filed M10 as "fix before it
bites 40 repositories," but the attended fleet run already happened before the review landed, so
the mitigation kimi-1 asked for arrived as the tool's own validate-and-rollback path rather than
as the ratified compare-and-swap/per-deck-confirmation design. The finding is upheld in full
(DF-1 lists every deviation: no CAS, bulk `--yes`, no foreign-deck version gate, thin inventory,
`workspace_dir`/`role`/`host_handle` not carried into committed TOML) — what changed is only
that it is now a hardening task rather than a pre-flight blocker, because the flight already
happened and landed clean.

The cheap guard shipping in this cycle (dirty-tree skip + `--confirm-breaking` on `roster
migrate`) is the right minimal scope: it prevents a *second* unattended run from compounding the
exposure while the full contract is outstanding. I would not accept deferral without that guard.
With it, deferring the full contract to a follow-up is the correct call — the residual risk is
to future fleet operations, not to the 24 already-migrated decks (which have backups and passed
post-write validation), and re-running migrate on an already-migrated deck is a no-op by the
resumability property.

### 5. Should anything be dismissed?

No. I reviewed §4 ("None") against all three review files and the agreed/deferred lists. Every
finding from all three reviews is either agreed (A1–A15) or deferred with a stated reason
(DF-1–DF-6). I do not see a finding that should have been dismissed as invalid. The one
candidate — kimi-1 NIT-5 (`RosterSnapshot` drops `Display`/provenance) — is correctly deferred
(DF-4) as "never decided either way" rather than dismissed, which is the honest classification.
codex-1's MAJOR on `PARLEY_HOME` machine-scope writes writing to the wrong nested path is not in
the agreed list, but I checked the current source: `rosterScopeFile` (`roster_set.go:89-107`)
now calls `config.CentralAgentsPath()` for the machine scope rather than composing
`$PARLEY_HOME/.parley/agents.toml`, with a comment naming the exact defect codex-1 reported.
`PRIMARY`. That fix landed in fix-up cycle 1 (v1.40.1) before the consensus was drafted, so it
is correctly absent from the cycle-2 list — it is not a dropped finding, it is an already-fixed
one. The consensus does not claim it as a cycle-2 item.

### 6. Anything the three reviews collectively missed?

Two items, neither rising to BLOCK:

(a) The §15.5 `## Drafter position changes` requirement. The review consensus draft (consensus.md
frontmatter: `drafter: claude-1`, `status: draft`) does not contain a `## Drafter position
changes` section. §15.5 requires it whenever the facilitator is also a participant and drafts
`consensus.md` — which is the case here (claude-1 is both drafter and a reviewed participant).
The design-phase `consensus.md` has it (15 entries, `consensus.md:248-274`); the review-phase
consensus does not. This is a §15.5 procedural gap, not a finding about the implementation. I
flag it as a condition below, not a BLOCK: the drafter should append a `## Drafter position
changes` section (or `None`) before the signoffs close. I own this observation `PRIMARY` (I read
the review consensus in full and the §15.5 rule at `COOPERATION.md:1299-1305`).

(b) A2's fix description says `render` "never deletes a row it did not generate without reporting
it" but does not name the mechanism for the `unmapped` row in the generated table. The normative
field table (`consensus.md:355`) says a §2-only ID is "reported `unmapped`, never auto-added" —
i.e. it surfaces as a row in `roster show` but is not written into the generated §2 by `render`.
The fix should make explicit whether the `unmapped` row appears in the generated §2 table (with a
status marker) or only in `roster show`/`--explain`. This is a fix-clarity note for Phase 8, not
a cycle-1 gap. `PRIMARY` (I read the field table and the A2 fix text).

Neither is a missed *implementation* defect; both are drafting-completeness notes. The three
reviews collectively covered the implementation surface well — the two CRITICALs and the
membership-layering escalation were independently corroborated, and the severity conflicts were
resolved by evidence.

## Conditions

1. §15.5 drafter position changes: the review consensus should append a `## Drafter position
   changes` section (the drafter is claude-1, a participant) before signoffs close, recording any
   material position change since its most recent round/review artifact, or `None` if there were
   none. This is a §15.5 procedural requirement (`COOPERATION.md:1299-1305`), not a finding about
   the code under review. It does not block my ACCEPT, but it should be corrected in the
   consensus file before the cycle closes.

2. A2 fix clarity (Phase 8): when implementing the §2-only-ID conflict rule, make explicit
   whether the `unmapped` row appears in the generated §2 table or only in `roster show`. The
   normative field table (`consensus.md:355`) decides "reported `unmapped`, never auto-added";
   the fix text should name the surface.
