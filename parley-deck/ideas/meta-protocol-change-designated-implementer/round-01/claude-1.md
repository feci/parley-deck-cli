---
agent: claude-1
idea: meta-protocol-change-designated-implementer
round: 1
date: 2026-09-24
---

## Protocol context attestation

Obtained with `parley protocol packet --dir . --phase 1 --track deliberation --idea
meta-protocol-change-designated-implementer --flag protocol_change --json`. Body read in full
(1,386 lines / 109,928 bytes).

```json
{
  "context_mode": "full",
  "source_sha256": "8ce83cde0b4f1c3019229869433df0c32451fdbaa13a186952fdf9805c3a9db7",
  "packet_sha256": "8ce83cde0b4f1c3019229869433df0c32451fdbaa13a186952fdf9805c3a9db7",
  "fallback_reason": null,
  "body_path": ".parley-runtime/protocol-packets/full-phase1-deliberation-8ce83cde0b4f1c3019229869433df0c32451fdbaa13a186952fdf9805c3a9db7.md"
}
```

`fallback_reason` is **absent** from the renderer's JSON (not null-valued); rendered as `null` above.
The renderer did not refuse. `packet_sha256 == source_sha256` and both equal the kickoff authority SHA
and `shasum -a 256 parley-deck/COOPERATION.md` at HEAD `9134c7a` — the packet is the live authority
verbatim, not an optimized subset. The envelope also carries a `shadow` block describing the packet
that was *not* used (`7ec29a81…`, 68,811 B, 36 of 69 blocks); I read the full body, not the shadow.

---

## Summary

Yes — the owner should be able to designate the implementer, and the mechanism should copy the
`facilitator:` precedent that shipped in 1.49.0 almost exactly. But the verification changed what I
think this idea is *for*. The brief's hypothesis is that with `author: user` the implementer is
decided by round-01 speed; at HEAD that is not what happens. Under the driver the implementer is the
**first entry of the `participants:` list**, because the FINAL-drafter inheritance rule is effectively
inert: the driver's own FINAL-draft prompt never asks for a frontmatter key that the resolver reads,
and 51 of 85 FINAL.md files on disk carry no readable key either. So today the role is selected by
list order — an undocumented rule nobody chose — and the documented claim fallback has **zero lines of
code** behind it. My proposal therefore does two things: make the existing selection explicit,
observable and consistent (repair), and add the owner's designation on top (feature), with the global
default shipped UNSET so the no-designation path stays byte-identical to today.

---

## Verification of the brief's locators at HEAD

HEAD `9134c7a4e1346c5a54cb2439c7ebcdf3b8e5324d`, branch `designated-implementer`.
The brief marks every item below as testimony; I am a verifier, not an owner, of these claims (§15.1).

| # | Brief's locator / claim | Verdict | Tag |
|---|---|---|---|
| 1 | Phase 4 ~L407-409: FINAL drafter is `author:`; with `author: user` default drafter is first round-01 submitter; volunteer with `Drafter: yes` | `CONFIRMED` — L407 and L409 exactly | `PRIMARY` |
| 2 | Phase 5 ~L443: default implementer is the FINAL drafter; claim via `inbox/<from>-to-all_<slug>_impl-claim.md`; facilitator never implements/verifies unless `facilitator_participates: true` | `CONFIRMED` — all three clauses on L443 | `PRIMARY` |
| 3 | Same facilitator rule also at ~L883 | `CONFIRMED` — L883 (§9.0 tail) | `PRIMARY` |
| 4 | `roles:` ~L311 advisory only; does not change drafter eligibility | `CONFIRMED` — L311, and the same sentence at L97 (§1) | `PRIMARY` |
| 5 | Reviewers are all non-implementers ~L234 | `CONFIRMED` — L234, the §4.0 per-track table, `deliberation` column | `PRIMARY` |
| 6 | No merge without an invokable non-implementer reviewer ~L538 | `CONFIRMED` — L538 | `PRIMARY` |
| 7 | Goal-done check by a fresh non-implementer ~L689 | `CONFIRMED` — L689 | `PRIMARY` |
| 8 | Two separate unexported `resolveImplementer` at `internal/app/driver_impl.go:104` and `internal/consensus/consensus.go:751` | `CONFIRMED` — both at those exact lines | `PRIMARY` |
| 9 | "Check whether they agree." | **They do not.** See the diff below. | `PRIMARY` |
| 10 | `facilitator:` / `facilitator_participates:` shipped in 1.49.0 | `CONFIRMED` — `CHANGELOG.md` §`1.49.0 — 2026-09-24` L26-29; `internal/protocol/facilitator.go` added in `86d028b`; `VERSION` = `1.49.0`; tag `v1.49.0` exists | `PRIMARY` |
| 11 | In that run kimi-1 claimed drafting, zcode-1 claimed implementation | `CONFIRMED` — `inbox/kimi-1-to-all_…_drafter-claim.md`, `inbox/zcode-1-to-all_…_impl-claim.md`, `IMPLEMENTATION.md` frontmatter `implementer: zcode-1` | `PRIMARY` |
| 12 | claude-1's first code review found 1 CRITICAL and 8 MAJOR | `CONFIRMED` — `…/review/round-01/claude-1.md`: 1 `### [CRITICAL]`, 8 `### [MAJOR]` (also 7 MINOR, 3 NIT) | `PRIMARY` |
| 13 | Organizer-token study measured self-implementation/self-verification at ~30-40% of an organizer's input | `CONFIRMED` for the copied study's own wording (`source-context/organizer-token-study.md` §2 item 1, incl. its 49.8% upper bound and its caveat that the session's organizer was also author, drafter, implementer and release manager) | `PRIMARY` |
| 14 | "implementation is the heaviest token work of a run" | `UNVERIFIED` — the same study puts post-compaction re-orientation at ~33-39% of input (§2 item 2) and compaction at 42% of all uncached input (§2 item 3). Nothing in the copied study ranks implementation first, and nothing measures implementation as a share of a *run's* total across all agents. | `PRIMARY` (on the study text; the ranking claim has no support in it) |
| 15 | Brief's hypothesis: with `author: user` the default drafter and therefore implementer is decided by who files round-01 first | **`WRONG` for the auto-driven path.** No code anywhere reads round-01 submission time to select a drafter or implementer. Evidence below. It remains a fair reading of the *prose* for a hand-run idea. | `PRIMARY` |

I did not open the external study at `/Volumes/…/organizer-token-study/2026-09-23/README.md`; item 13 is
`PRIMARY` on the in-deck §6-rule-4 copy only. Any claim about the *original* file is testimony I do not own.

### The two resolvers do not agree — exact delta

Both function bodies are byte-identical for 26 lines (the participant check and the
`IMPLEMENTATION.md` → `FINAL.md` frontmatter walk). `diff` of `driver_impl.go:104-131` against
`consensus.go:751-776` is a single hunk — the terminal fallback:

```go
// internal/app/driver_impl.go:129-131   — positional fallback
	if len(participants) > 0 {
		return participants[0]
	}
	return ""

// internal/consensus/consensus.go:775   — no fallback
	return ""
```

Both behaviours are deliberately test-pinned, on opposite sides:

- `internal/app/app_test.go:1573-1577` (`TestResolveImplementerFromRoleMetadata`) asserts the
  positional fallback (`want codex (fallback; hermes not a participant)`).
- `internal/consensus/roundgate_test.go:100-107` (`TestUnresolvableImplementerExpectsEveryone`)
  asserts the `""` path, whose consumer `expectedRoundParticipants` then **expects everyone** —
  "fails closed toward asking for more rather than silently accepting a short round"
  (`consensus.go:723-729`).

So "make them one function" is not a free win. A naive merge breaks one of the two pinned
behaviours. §D10 below proposes the shape that does not.

---

## What HEAD actually does

Everything in this section is a claim I assert first canonically, so I am its **owner** and issue no
verdict on it (§15.1). Each is supplied with `PRIMARY`-grade evidence — locator plus quoted source or
quoted command output — for a non-owner to verdict in round 2.

**F-1. Under the driver, the implementer is `participants[0]` minus the declared facilitator — list
order, not speed, not the drafter.** `newDriverImplOps` (`driver_impl.go:45-98`) filters the
declared facilitator out of `participants` into `eligible`, then calls
`resolveImplementer(ideaDir, eligible)`. The value is stored in a struct field and **frozen at driver
construction** (`driver_impl.go:93-97`); `Implement()` launches exactly that id
(`driver_impl.go:213-223`). The three production call sites pass the idea's `participants:` list
(`app.go:1252`, `app.go:1984`, `app.go:2038`), and `parseList` (`workspace.go:399-414`) preserves
declaration order. On a fresh `parley run` neither `FINAL.md` nor `IMPLEMENTATION.md` exists, so the
frontmatter walk finds nothing and the positional branch fires.

**F-2. The FINAL-drafter inheritance rule is inert in practice, because the resolver reads keys the
protocol's own FINAL template does not produce.** `resolveImplementer` reads only `implementer` and
`drafted-by` (both copies, lines 117-118 / 763-765). The Phase 4 FINAL.md template in COOPERATION.md
(L414-420) prescribes `author:`, not `drafted-by:`. Measured across the deck:

```
total FINAL.md files: 85
  carries `implementer:`  -> 11
  carries `drafted-by:`   -> 25
  carries `author:`       -> 32
  carries `drafter:`      -> 23
  carries `finalized-by:` -> 4

FINAL.md files the resolver CANNOT read a drafter/implementer from
(no `implementer:` and no `drafted-by:`): 51
  of those, carrying `author:`       -> 32
  of those, carrying `drafter:`      -> 14
  of those, carrying `finalized-by:` -> 4
```

60% of historical FINAL.md files are unreadable to the rule that is supposed to consume them — and
the single most common key among them is `author:`, the one the protocol itself prescribes. The
precedent run confirms it end to end: `ideas/meta-protocol-change-lean-organizer/FINAL.md` frontmatter
is `idea, status, author: kimi-1, consensus-date, participants` — no `implementer:`, no `drafted-by:`.

**F-3. The driver instructs its own drafter to omit the key its own resolver needs.**
`buildFinalDraftPrompt` (`driver_consensus.go:174-199`) says verbatim:

```
YAML frontmatter MUST include:
  idea: %s
  status: final
```

That is the complete frontmatter requirement. Neither `implementer:` nor `drafted-by:` nor `author:`
is requested, and `ValidateImplementationArtifact` (`runner/phase58.go:426-460`) likewise requires
only `idea`, `status` and a non-empty `## Summary of work` — never `implementer:`. The loop that is
meant to carry the drafter's identity into Phase 5 is open at both ends.

**F-4. The FINAL drafter and the implementer are chosen by two different functions with different
eligibility filters, so they can disagree with each other.** The drafter is
`firstEligibleHeadlessAgent` (`driver_consensus.go:109-123`): first participant, in list order, that
is **discovered, `Found`, and headless**, skipping the declared facilitator. The implementer is
`resolveImplementer`: first participant in list order, skipping the declared facilitator, with **no
discovery or launch-mode filter at all**. An uninstalled or interactive-mode `participants[0]` is
skipped for drafting and selected for implementing.

**F-5. The `impl-claim` fallback has no implementation.** `grep -rn "impl-claim" --include="*.go" .`
returns nothing. The string exists only in prose (`parley-deck/COOPERATION.md`,
`internal/protocol/defaults/COOPERATION.md`), in hand-written inbox files, and in the lean-organizer
artifacts. No driver, preflight, consensus or wait path reads, validates or honours a claim file. In
the precedent run the claim took effect only because a human read it and zcode-1 then wrote
`implementer: zcode-1` into `IMPLEMENTATION.md` itself — i.e. after Phase 5 had already begun.

**F-6. The divergence is user-visible through `parley wait`.** `internal/driver/phasedigest.go:364`
calls `consensus.ExpectedRoundParticipants` (`consensus.go:1089-1091`), which routes to the
**consensus** copy. With a protocol-conformant `FINAL.md` (`author:` only) and no `IMPLEMENTATION.md`
yet, the driver launches reviewers `eligible[1..]` while `parley wait`'s digest reports the review
round as owing a file from **all** participants — including the agent the driver picked to implement,
whom §6/Phase 6 forbids from reviewing its own work. A `parley consensus draft --review` in the same
window fails with `… is incomplete; missing <implementer>.md` (`consensus.go:176-181`).

**F-7. `parley run --participants` is an undocumented implementer switch.** `parley run --help`
lists `-participants string  comma-separated agent IDs to run` and `-preset`. Given F-1, the *order*
of that flag's value silently decides who implements. Nothing in the help text, the protocol, or the
skill says so.

**F-8. Preflight's freshness block reports a remembered SHA, not a measured one.**
`preflight.go:486-487` assigns `fr.LiveSha = meta.ProtocolSha256`, read from
`parley-deck/meta/version.json` (`readVersionMeta`, `preflight.go:457`). That file records
`"protocolSha256": "12e4b31c…"`, `"updatedAt": "2026-09-18T09:32:02.255Z"`. The live file hashes to
`8ce83cde…`. `12e4b31c…` is exactly `parley-deck/COOPERATION.md` at commit `b4831d6`
("release: finish evidence-first CLI 1.48.0") — i.e. the pre-1.49.0 text, stale by the two
lean-organizer protocol commits `86d028b` and `64a622c`. For this `source`-role deck it is currently
harmless (`preflight.go:489-492` short-circuits to `source-advisory` before any drift comparison), and
that is why it has gone unnoticed; but it is a direct constraint on D11: **new validation must read
live idea frontmatter, never `meta/version.json`.** The copied `source-context/preflight.json` carries
the same stale value, so that snapshot does not describe the tree this run is deliberating on.

**F-9. The three protocol copies are in lockstep today.** `parley-deck/COOPERATION.md` `8ce83cde…`
(109,928 B); `internal/protocol/defaults/COOPERATION.md` `cce5d7d9…`, differing by 16 lines that are
all project-specific zones (workspace name, created date, `Protocol synced:` line, §2 host-handle
table); `…-skill/skills/parley-deck/references/COOPERATION.md` `fc907e59…` (109,772 B), differing by
18 lines of the same zones, and byte-identical to the published core 2.13.0 recorded in
`source-context/prior-release-handoff.md`. All three carry the Phase 4 rule at L409/L402/L402 and the
Phase 5 rule at L443/L436/L436. Any hunk this idea agrees must land in all three, identically.

### What this changes about the problem statement

The owner asked "who does it now? is it prescribed in the protocol?" The honest answer that falls out
of F-1 … F-5 is: **the protocol prescribes one rule, the code implements a different one, and a third
rule is documented with no implementation.** The prose says FINAL-drafter-inherits with a claim
fallback; the driver does list-order; the claim fallback does not exist. Adding a designation field on
top of that without repairing it would give the owner control of one path while the unset path stays
silently wrong — which is the worse outcome, because the unset path is what every other deck runs.

---

## Proposed approach

### D1 — The answer

**Yes, and it should be an explicit designation rather than a preference.** Three reasons drawn from
the verification, not from the brief: (a) the role is already being assigned today by an accident of
list order (F-1) — designation replaces an implicit rule with an explicit one rather than adding a new
axis of control; (b) the protocol already has a working precedent for exactly this shape of per-idea
role metadata in `facilitator:`, including preflight gating and driver eligibility filtering, shipped
one version ago; (c) the alternative mechanisms that already exist are all either advisory by
construction (`roles:`, L311) or negative (`facilitator:` excludes, it does not appoint).

### D2 — Per-idea field

`implementer: <agent-id>` in `00-prompt.md`, optional. Absent → today's resolution, unchanged.
Parsed by a new `internal/protocol/implementer.go` mirroring `facilitator.go` almost line for line:

```go
const ImplementerKey = "implementer"

type ImplementerDesignation struct {
	Declared bool
	Agent    string
}

func ReadImplementerDesignation(ideaDir string) ImplementerDesignation
func ImplementerDesignationFromMeta(meta map[string]string) ImplementerDesignation
func (d ImplementerDesignation) Conflict(participants []string, role FacilitatorRole) string
```

`Conflict` returns a non-empty fail-closed message, naming both fields, when the designated agent is
(a) not in `participants:`, or (b) the declared facilitator of a run without
`facilitator_participates: true`. This is the same contract as `FacilitatorRole.Conflict`
(`facilitator.go:59-71`) and reuses its `IneligibleForRoles` predicate rather than restating it.

### D3 — Global default, shipped UNSET, and precedence

`[defaults].default_implementer` in `~/.parley/agents.toml`, **absent on ship**. `parley init` writes
it as a commented line, the way `ping_tier` / `preferred_transport` / `roster_change_policy` are
already emitted (`internal/config/runtime.go:637-639`).

The precedence question has a better answer than a ranking table: **materialize the global default at
Phase 0 instead of consulting it at Phase 5.** When `default_implementer` is set and `00-prompt.md`
has no `implementer:`, kickoff writes the concrete agent id into the idea with a provenance comment;
from then on there is exactly one rule in force and one place to read it. This is not a new idea — it
is precisely what `ResolveRoster` already does for participants (`internal/config/roster.go:73-95`:
explicit preset → track default → "no preset applies … the caller keeps today's behavior", with a
`Provenance` one-liner). Adopting the same pattern means the owner can always see, in the idea file,
who was designated and why, and a later change to the machine config cannot retroactively reinterpret
an in-flight idea.

Resulting resolution order, total and observable:

1. `IMPLEMENTATION.md` `implementer:` — the re-entry pin (existing; must never be overridden mid-flight).
2. `00-prompt.md` `implementer:` — the designation (owner-written, or materialized from the global default at Phase 0).
3. A perfected claim file — **only when (2) is absent** (see D7).
4. `FINAL.md` drafter key — the inheritance rule (see D10 for the key-set repair).
5. Positional `participants[0]` — retained, but **loud** (see D11).

### D4 — Designated agent unavailable: fail closed, never silently reassign

Three windows, three answers, all matching machinery that already exists:

- **Kickoff / §9.0 readiness.** If the designated agent pings unavailable, this is the case §9.0
  already governs: exclusion requires *explicit user confirmation* and is recorded in `00-prompt.md`.
  The designation follows the same rule and records the same way —
  `implementer_waived: <id> — <reason> — confirmed <date>` naming the replacement. No silent
  fall-through to list order.
- **Driver construction.** Copy the `roleErr` pattern verbatim (`driver_impl.go:52-63`): when the
  designation cannot be honoured, the ops carry an error that **every role action escalates with**,
  and the driver "NEVER silently falls back". That comment is already in the file; the new field
  should be inside its scope, not beside it.
- **Mid-Phase-5 death.** Escalate via `inbox/` and halt. Consistent with LE-5: hitting a ceiling
  "escalates … never marks an idea complete" (L697-702). Reassignment after a partial implementation
  is a human decision, not a driver decision, because the partial tree is now evidence.

### D5 — Scope of the role: code *and* non-code

No new vocabulary. Phase 5 already covers non-code work ("If the idea is design-only (no code
artifact), Phase 5 may be reduced to a brief `IMPLEMENTATION.md` describing where the design output
was applied", L508), and the owner's own words are "writes the code if needed, or does whatever work
is needed". The implementer is therefore defined as *the participant that executes FINAL*, whatever
the artifact. **Exactly one at a time** — the lean-organizer claim file already states the invariant
("Exactly one drafter and exactly one implementer at any time; no parallel writers") and every
resolver returns a single id. I explicitly do not propose split scopes, co-implementers, or
per-file ownership; nothing in the owner's request asks for it and each would fork the reviewer set.

### D6 — Drafter versus implementer, and self-contained FINAL

The protocol already requires FINAL to be implementable by a non-drafter: "`FINAL.md` plus
`IMPLEMENTATION.md` MUST be self-contained enough that a fresh agent or the auto-drive driver can
implement or resume **from them alone**, without session transcripts" (L433). A designation makes
drafter≠implementer the *common* case rather than the exception, which strengthens that requirement
rather than weakening it. Two proposals, one soft and one mechanical:

- **Soft:** when the designated implementer is also the FINAL drafter, `FINAL.md` records the role
  concentration in one line. This is §15.5's existing device applied to a second role, not a new gate —
  §15.5 already mandates exactly this for the facilitator-drafter case, and existing signoffs ratify it.
- **Mechanical:** `firstEligibleHeadlessAgent` gains one more skip — prefer a drafter that is *not*
  the designated implementer, **falling back to the implementer when no other eligible drafter
  exists**. Ordering-only, fail-open in the degenerate case, and it uses the filter slot the function
  already has (`role.IneligibleForRoles(p)` at `driver_consensus.go:113-115`). On a 3-participant deck
  this alone buys drafter separation with no new obligation.

I do not propose forbidding drafter==implementer. With two participants it would make a designation
unsatisfiable, and the protocol's own two-participant rule (L466-468) says the rules apply unchanged.

### D7 — The claim mechanism stays, as a genuine fallback, and gets the code it never had

Keep it, subordinate it, and make it real. Concretely:

- A claim binds only when `00-prompt.md` declares no `implementer:`. An owner designation is the
  owner's instruction; a participant does not override it by filing a file. A participant that
  disagrees escalates (§4 escalation) — which the protocol already provides and which produces a
  record the owner can act on.
- Give the claim file a reader. At minimum, `parley preflight` and the driver should *detect* a
  well-formed `inbox/<from>-to-all_<slug>_impl-claim.md` and surface it, so the documented fallback
  stops being prose-only (F-5). Whether the driver should *honour* it automatically, or only report it
  for the organizer to materialize into `00-prompt.md`, is the one place I would accept either answer —
  see the open questions.

### D8 — Reviewer model diversity when one agent implements everything

The machinery already exists and already escalates: `checkModelDiversity`
(`driver_impl.go:178-210`) emits an `agent.model_diversity` event and either warns or, with
`require_model_diversity: true`, refuses to open review; `fast` forces the gate on regardless of the
flag (`driver_impl.go:187-191`). For this roster a persistent designation is structurally safe —
`parley roster show --scope machine` reports three distinct model companies among the active agents
(Anthropic / Moonshot AI / Zhipu AI), so designating any one leaves two model-diverse reviewers.

The real exposure is a **two-participant deck**, where designating one agent leaves exactly one
reviewer, possibly same-model. Proposal: move the existing gate earlier — if a designation would leave
no model-diverse reviewer and `require_model_diversity: true`, refuse at Phase 0 rather than at
Phase 6. This is relocating a gate that already exists, not adding one; it costs the owner a Phase-0
error message instead of a wasted design cycle.

### D9 — Does the implementer take part in design rounds and signoffs? Yes; change nothing

The owner's words are "the others evaluate **it**" — the work, i.e. Phase 6, which already excludes
the implementer (L234, L512, `consensus.go:723-746`). Nothing asks to remove the implementer from
design. Removing it would shrink quorum (§5: "Quorum = all agents listed in `participants:`"), and
Phase 7 explicitly includes the implementer in review-consensus signoffs ("Each active participant
(implementer included) APPENDS their signoff block", L578). The designation should therefore change
**who executes Phase 5 and Phase 8, and nothing else.** Because the question was asked explicitly, it
is worth one sentence of protocol text saying so, so the next reader does not have to re-derive it.

### D10 — Resolver consistency: share the reading, keep the two fallbacks

One exported resolver in `internal/protocol`, returning the *source* as well as the id, with each call
site keeping its own terminal policy:

```go
type ImplementerSource int // SourceImplementationMD, SourcePrompt, SourceClaim, SourceFinalMD, SourceNone

func ResolveImplementer(ideaDir string, participants []string) (id string, src ImplementerSource, ok bool)
```

- `internal/app/driver_impl.go`: `if !ok { id, src = eligible[0], SourcePositional }` — byte-equivalent
  behaviour, pinned by `TestResolveImplementerFromRoleMetadata`.
- `internal/consensus/consensus.go`: `if !ok { return participants }` — byte-equivalent behaviour,
  pinned by `TestUnresolvableImplementerExpectsEveryone`.

This removes the duplicated *reading* logic — the actual drift risk, since both copies must learn the
new precedence order — while preserving both intentional fallbacks. `parley wait` inherits the fix
for free through `ExpectedRoundParticipants` (`consensus.go:1089-1091`), closing F-6.

**Prior art I must address, because it is mine.** In the lean-organizer review I filed
`[MINOR] MIN-7` against a `consensus.ResolveImplementer` export
(`…/review/round-01/claude-1.md:607-612`), and commit `64a622c` dropped it ("ResolveImplementer export
dropped (F15)"). That finding was that the export had **zero callers** and a doc describing behaviour
it did not have — not that sharing the resolver is wrong. This proposal has two real callers plus a
third consumer through `phasedigest.go:364`, and the doc would describe the behaviour it has. If the
group judges that re-litigated, the fallback is an unexported helper in `internal/protocol` with thin
per-package wrappers; I have no attachment to the export specifically.

**The key-set repair, flagged as a behaviour change.** Fixing F-2/F-3 means (i) `buildFinalDraftPrompt`
requires `implementer: <drafter-id>`, and (ii) the resolver's FINAL.md key set widens to
`implementer, drafted-by, author, drafter, finalized-by`. (i) is inert for existing decks. (ii) is
**not** — it makes the drafter-inherits rule start working on 51 historical FINAL.md files. It is safe
in shape (every candidate is still validated against `participants:`, so `author: user` simply does not
match) but it is a real change, it is *not* covered by "default UNSET means behaviour is exactly
today's", and it must be decided explicitly and recorded in `## Agreed trade-offs`. My position: do it.
Shipping designation on top of an inheritance rule that silently does not run leaves the unset path —
the path every other deck is on — wrong.

### D11 — Driver and preflight, fail-closed

- `IdeaStatus` gains `Implementer ImplementerDesignation`, populated in `ReadWorkspaceStatus`
  alongside `FacilitatorRole` (`workspace.go:353-359`).
- `preflight.go` gains `implementerDesignationGates(root)` beside `facilitatorConflictGates(root)`
  (`preflight.go:323-340`): same `gate{Kind, Detail, Confirm}` shape, same `confirmCommand`, same
  fail-closed treatment of an unreadable workspace (the F-fix already commented at `preflight.go:319-322`).
  Gate kind `implementer-designation`.
- **Read live idea frontmatter, never `meta/version.json`** — F-8. `facilitatorConflictGates` already
  does the right thing via `protocol.ReadWorkspaceStatus`; the new gate must copy that, not the
  freshness path.
- A durable `agent.implementer_resolved` event `{idea, implementer, source}` on the existing store,
  modelled on the `agent.model_diversity` event (`driver_impl.go:196-206`), and a one-line stdout
  notice. This is what makes F-1 and F-7 visible without changing them: a deck that never sets the
  field still runs positional order, but now says so.
- `parley run --participants` help text and the skill gain a sentence noting that order currently
  decides the implementer when nothing else does.

### D12 — Protocol text

Hunks in §4 Phase 4 (drafter-separation sentence), §4 Phase 5 (the designation, precedence,
unavailability, claim subordination, and the "designation changes Phase 5/8 only" sentence), §4.0
Phase 0 frontmatter template (`implementer:` line), §9.0 (designation checked in the readiness check),
and §10 TL;DR item 6. Identical bytes in all three copies (F-9), plus
`meta/protocol-changelog.md`. `parley protocol publish` stays the owner's attended action; no
participant stages, publishes or installs a core.

### D13 — What I deliberately do not propose

No change to quorum, roster, `roles:` semantics, model or effort. No co-implementers or split scopes.
No removal of the positional fallback (it would break every deck that does not opt in). No
`strict_implementer`-style new flag — `require_model_diversity` plus the gate relocation in D8 covers
the diversity case, and a second flag would be configurability nobody asked for. No selection of the
owner's default agent (§Process notes).

---

## Existing alternatives

Per §15.6(a): for each mechanism this proposal would build by hand, what the toolchain already ships,
with a locator.

| Mechanism I propose | What already ships | Locator | Why it is not enough |
|---|---|---|---|
| Per-idea `implementer:` role field | `facilitator:` / `facilitator_participates:` — identical shape: optional key, `Declared` struct, `Conflict()`, `IneligibleForRoles()` | `internal/protocol/facilitator.go` (whole file); `COOPERATION.md:443`, `:883` | It **excludes** an agent from roles; it cannot appoint one. `IneligibleForRoles` is a negative predicate. |
| Per-idea advisory role hints | `roles:` map in `00-prompt.md` | `COOPERATION.md:311`, and the same sentence at `:97` | Both sentences state it "do[es] not change quorum, signoff weight, artifact ownership, drafter eligibility, or roster membership" — advisory by construction; using it would contradict ratified text. |
| Owner-set global default with per-idea override | `[defaults].track_rosters` + `[rosters.<slug>]` presets, resolved by `ResolveRoster` (explicit → track default → none) with fail-closed validation against the §2 active roster and a `Provenance` line | `internal/config/runtime.go:277-291`; `internal/config/roster.go:35-95` | Selects *who participates*, not *who executes*. `parley preset list` reports "No roster presets defined" on this deck. I propose copying its precedence and provenance pattern rather than inventing one. |
| Claim-based selection | `inbox/<from>-to-all_<slug>_impl-claim.md` | `COOPERATION.md:443`; used at `inbox/zcode-1-to-all_meta-protocol-change-lean-organizer_impl-claim.md` | Prose only — `grep -rn "impl-claim" --include="*.go" .` is empty (F-5). It works today only because a human reads it. |
| Durable per-idea record of the implementer | `IMPLEMENTATION.md` `implementer:`, read first by both resolvers; written by the Phase-5 prompt | `internal/runner/phase58.go:248`; `driver_impl.go:117`, `consensus.go:764` | Written *after* Phase 5 starts — it records the outcome, it cannot select the input. |
| FINAL-drafter inheritance | `FINAL.md` `implementer:` / `drafted-by:` in both resolvers | `driver_impl.go:118`, `consensus.go:765` | Inert in practice: the driver's FINAL prompt requires only `idea:`/`status:` (F-3) and 51 of 85 FINAL.md files carry no readable key (F-2). |
| Reviewer diversity gate for a concentrated implementer | `require_model_diversity:` → `checkModelDiversity`, warn by default, escalate when set, forced on `fast` | `driver_impl.go:178-210`; `internal/driver/transport.go:63-70`; `COOPERATION.md` LE-3 | Fires at Phase 6, after a design cycle is spent. I propose relocating it to Phase 0 for the designation case, not building a second gate. |
| Fail-closed role validation at kickoff | `facilitatorConflictGates` → blocking `gate{Kind,Detail,Confirm}`, live-frontmatter sourced | `internal/app/preflight.go:314-340` | Covers the facilitator field only. Directly extensible; this is the template, not a competitor. |
| Observability of a resolved role | `agent.model_diversity` durable store event | `driver_impl.go:196-206` | No equivalent event for implementer resolution, which is why F-1 went unobserved. |
| Shared resolution across packages | `consensus.ExpectedRoundParticipants` exported precisely so `parley wait` "reuse[s] the SAME resolution instead of forking it" | `consensus.go:1085-1091`; consumer `internal/driver/phasedigest.go:364` | The stated principle already exists; it just was not applied one level down, to `resolveImplementer` itself. |

Sources consulted for this scan: the full packet body; `internal/protocol/`, `internal/app/`,
`internal/consensus/`, `internal/config/`, `internal/driver/`, `internal/runner/`;
`parley run --help`; `parley preset list`; `parley roster show --scope machine`;
`~/.parley/agents.toml`; `CHANGELOG.md`; the idea's `source-context/`.

---

## Concerns / open questions

1. **Does a claim override the global default?** D3 makes the question mostly disappear by
   materializing the default into the idea at Phase 0 — after which a claim faces a per-idea field and
   loses per D7. But that is a *design choice*, and the alternative (global default is weaker than a
   claim) is defensible: it keeps a standing preference from silencing an agent who knows it is better
   placed for one specific idea. I lean to materialization because the owner's motivation is precisely
   that the role should not be decided by whoever moves first. Worth one explicit decision.

2. **Should the driver honour a claim file automatically, or only report it?** (D7.) Automatic
   honouring restores a speed race — first claim wins — which is the failure mode the owner is trying
   to remove. Report-only keeps a human or the organizer in the loop but leaves the documented
   fallback still not self-executing. My weak preference is report-only plus a preflight gate when a
   claim exists and no designation does, but I would sign either.

3. **Is the D10 key-set widening in or out of scope?** It is the one behaviour change to existing
   decks in this proposal. In-scope reading: the rule this idea lets the owner override does not
   currently run, so shipping the override without it is shipping half. Out-of-scope reading: it is a
   separate defect with its own blast radius and belongs in its own idea. I argue in-scope (D10) but
   this needs an explicit `## Agreed trade-offs` entry either way, and if the group splits it out, the
   split must be recorded rather than assumed.

4. **Does designation interact with `no-implement` and design-only ideas?** `parley run
   --no-implement` stops at FINAL. A designation on a design-only idea then names an agent that never
   acts. Harmless, but preflight should probably not gate on availability of an implementer for an
   idea that will not reach Phase 5. I have not traced whether `--no-implement` is visible at preflight
   time; someone should.

5. **`IMPLEMENTATION.md` re-entry precedence versus a corrected designation.** Rank 1 pins the
   re-entry to whoever actually wrote the artifact. If the owner corrects a wrong designation *after*
   a failed Phase 5 attempt left an `IMPLEMENTATION.md` behind, the pin wins and the correction is
   ignored. I think that is right (the partial tree is evidence, D4) but it should be stated, not
   discovered.

6. **Unverified in the brief, still unverified here.** "Implementation is the heaviest token work of
   a run" (item 14) has no support in the copied study, whose own numbers put re-orientation and
   compaction at comparable or larger shares. If that claim is load-bearing for anyone's argument, it
   needs its own measurement; it must not enter `FINAL.md` as established (§15.2, `RECALL`-only
   material stays `UNVERIFIED`).

7. **`preflight.json` in `source-context/` is stale** (F-8): its `protocolSha256` is the v1.48.0-era
   text, not this tree's. Its roster and gate rows look current and I relied on them only as testimony
   corroborated by my own `parley roster show` run. Anyone citing that file for protocol state should
   not.

8. **Inherited release limitations, assessed against this idea's delivery criteria** (per the §Release
   boundary constraint). From `source-context/prior-release-handoff.md`: CLI winget not submitted, npm
   not published (2.12.1 latest), a 1.49.1 candidate staged but untagged, and an unresolved Windows
   scope decision. None of these block *designing* or *implementing* this change. Two touch delivery:
   (a) this idea's protocol hunks must be staged on top of the published core 2.13.0 (`fc907e59…`),
   which the handoff reports as live and which I verified is byte-identical to the skill worktree's
   bundled copy (F-9) — so the base is real, not assumed; (b) a new CLI minor must be chosen above
   whatever is actually released at staging time, and with 1.49.1 staged-but-unpublished that number
   is not knowable from here. I record it as a release-time check, not a design input. I inherit no
   authorization to touch Windows platform architecture and propose none.

---

## Risks

- **R-1 (highest) — the repair is bigger than the feature.** F-1 … F-6 are defects in shipped code, not
  gaps this idea invents. A reviewer could reasonably call D10's key-set widening scope creep, and a
  reviewer could equally reasonably call shipping without it negligent. Whichever way it goes, the
  decision must be explicit; the failure mode is deciding it by silence. Mitigation: it is already
  named as a required `## Agreed trade-offs` entry (Q3).
- **R-2 — a designation concentrates implementation in one model.** Over many ideas, one model's blind
  spots become the codebase's blind spots, and reviewers who never implement lose the tacit knowledge
  that makes review sharp. The existing LE-3 gate protects the *reviewer set*, not this. It is a real
  cost of what the owner asked for, and the honest mitigation is the owner setting and revisiting the
  default deliberately — which is exactly why shipping UNSET with a post-release recommendation is the
  right sequencing, not a hedge.
- **R-3 — fail-closed gates make kickoff brittle.** Every new blocking gate is a new way for
  `parley run` to stop. This deck already carries one blocking `facilitator-declaration` gate from an
  unrelated historical idea (`source-context/preflight.json`; also recorded in `00-prompt.md`), which
  had to be worked around for this very run. A second gate class of the same shape will do the same
  thing to somebody. Mitigation: gate only on the *idea being run*, not on every idea in the deck —
  and note that `facilitatorConflictGates` currently does the latter (`preflight.go:327-339` iterates
  `status.Ideas`), which is why the unrelated idea blocks this one.
- **R-4 — three copies, one hunk.** F-9 shows they are in lockstep now. A partial application would
  leave decks disagreeing about who implements, which is worse than no change. Mitigation: a
  byte-comparison of the changed region across all three as an acceptance criterion.
- **R-5 — the `--participants` order dependency becomes load-bearing documentation.** Writing F-7 into
  the help text hardens an accident into a contract. Mitigation: word it as current behaviour with the
  designation as the supported way to control it, not as a guarantee.
- **R-6 — my verification is static.** I read source, tests, artifacts and CLI help; I did not execute
  an end-to-end `parley run` to observe the implementer selection empirically, because the run
  instruction forbids running the workflow and spawning roster identities. F-1's claim that the
  positional branch fires on a fresh run is a reading of `newDriverImplOps` + `Implement()` + the three
  call sites, not an observed launch. A non-owner should verdict it, ideally against a fixture —
  `internal/app/facilitator_autodrive_test.go:146` already constructs `newDriverImplOps` directly and
  is the cheapest place to pin it.
- **R-7 — evidence asymmetry on the token study.** My item-13 verdict is `PRIMARY` on the in-deck copy
  only; the external original was not opened (§6 rule 4 disclosure). If anyone's argument depends on
  the original's numbers rather than the copy's, that asymmetry has to be disclosed rather than assumed
  away (§6 rule 4, second paragraph).

---

## Process notes

**Willingness (current protocol, voluntary).** Under today's rules I am willing to draft `FINAL.md`
and/or to implement, if the group wants either. I have **not** perfected a claim: under the current
protocol a drafter volunteer is declared in the consensus signoff block (`Drafter: yes`, L409) and an
implementer claim in `inbox/<from>-to-all_<slug>_impl-claim.md` before Phase-5 work begins (L443) —
neither is a round-1 artifact, and I will not pre-empt those phases here. This is a statement of
availability under the protocol in force, and it is **not** a bid for the owner's future global default
implementer; the two are different questions and conflating them would be exactly the kind of quiet
self-selection this idea exists to remove.

**Owner-decision boundary.** I name no agent as the owner's default. The only implementation-quality
evidence in this deck is n=1 (zcode-1 implemented lean-organizer; my review found 1 CRITICAL and 8
MAJOR, all fixed across 3 fix-up cycles and 4 review rounds, closing on a zero-fix consensus). One run
with a good ending is not a basis for a persistent default. A defensible recommendation after release
would need, at minimum: findings-per-implementation and fix-up-cycle counts by implementer across
several ideas, and per-implementer cost from the 1.49.0 usage ledger — the instrument the prior release
explicitly declined to make claims from yet.

**Evidence ownership (§15).** Verdicts in the locator table are mine as a *verifier* of the organizer's
testimony; the brief marks all of it unverified, so none of it is mine to own. Claims F-1 … F-9 are
asserted here first canonically, so I own them and issue **no** verdict on any of them; each carries
`PRIMARY`-grade evidence for a non-owner to verdict in round 2. Where I quote the copied study or the
copied preflight snapshot, I am transcribing testimony and have marked it as such rather than relying
on it as established.

**Independence.** Written before reading any other round-01 file. `round-01/zcode-1.md` exists on disk
at the time of writing; I did not open it.

<!-- Original language of the owner's quoted words: Slovak; translations as supplied verbatim in 00-prompt.md and source-context/owner-brief.md. -->
