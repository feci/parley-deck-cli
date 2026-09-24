---
agent: claude-1
idea: meta-protocol-change-lean-organizer
review-round: 1
date: 2026-09-24
reviewed-commit: 3c97f44
---

## Summary

The implementation ships all four scope items as enforced tooling, both suites are green when I
rebuild and run them myself, and several load-bearing properties survive adversarial probing: the
never-cut floor holds at runtime against a hostile audience map, the staged core 2.13.0 contains
exactly the three ratified change sets with a placeholder header and a stubbed §2, the three
COOPERATION.md copies carry byte-identical hunks, the SKILL.md relocation drops nothing, and
`usage ingest` refuses every count-accepting flag I could invent. Against that, one shipped change
is a release blocker: a **new hard consensus-section gate** (`internal/consensus/consensus.go:547`)
that was not part of the ratified A.4 scope now rejects the protocol's **own published Phase-3
template** and **79 of 80** existing consensus artifacts in this repo's deck, contradicting two
explicit sentences in the very COOPERATION.md this change edits. Eight further MAJOR findings sit
mostly in B and C: `parley wait` cannot reach exit 0 on this deck at all (an unrelated six-week-old
escalation and the implementer's own P.2 note both trip it), and every R-2 byte figure recorded in
`IMPLEMENTATION.md` is **pre-hunk and non-reproducible** at either committed state — including both
numbers the phase-8 quorum disposition rests on.

**Provenance (§15.2).** Every verdict below is `PRIMARY` unless tagged otherwise: each is a check I
executed in my own isolated worktrees, with the command and relevant output quoted. I ran no
release, version, publish or global-install action and did not run `protocol publish`.

**Tree provenance.**

| Repo | Review worktree | Branch | Commit | Tree |
|---|---|---|---|---|
| CLI | `../worktrees/lean-organizer-review-claude-1` | `review/meta-protocol-change-lean-organizer/claude-1-20260924` | `3c97f445336a70bd7bd57eb15a0da78144bcd255` | `dcc3ce2357f8c924ebf34136f663452c68edc793` |
| Skill | `../worktrees/lean-organizer-review-claude-1-skill` | same branch name | `0f513f6c0fc78523d11321e60868942a08079d1d` | `21004a477829b1afbbdf97f9a1a69a15a45b9aa4` |

CLI bases: `b37f7ef9dd0941b21c3ba146c91d1a118259cbed` (pre-implementation) →
`86d028b5558a60858534ac72f25537e52ed395fc` (implementation) → `3c97f44` (record).
Skill base: `d1e57d5c8a56ac40fc0c8b4ac035069d2c0b9283` → `0f513f6`.
My own build: `go build -o /tmp/claude1-review/parley ./cmd/parley` → sha256
`abaad23610c874ec827cf00255285d1f03212544c6994367f305562cc846701b`, `go version go1.27.1 darwin/arm64`.
Comparison baseline binary built from `git archive b37f7ef` → `/tmp/claude1-review/parley-base`.
Protocol packet attestation retained:
`{"context_mode": "full", "source_sha256": "3ad8a7ebbe954b5299c38a28d8a270ae66cf360265c54f8bd6e70e857bb26086", "packet_sha256": "3ad8a7ebbe954b5299c38a28d8a270ae66cf360265c54f8bd6e70e857bb26086", "fallback_reason": null, "body_path": "/Volumes/My Shared Files/AI_WORKSPACE/parley-deck/worktrees/lean-organizer/.parley-runtime/protocol-packets/full-phase6-deliberation-3ad8a7ebbe954b5299c38a28d8a270ae66cf360265c54f8bd6e70e857bb26086.md"}`
— I verified the file hashes to that value (`shasum -a 256`, 109,663 B).

## Refutation attempts

Refutation-default (LE-1): for each FINAL acceptance criterion I tried to construct a failing case.
Unsuccessful attempts are recorded as such — they are evidence the property held.

### A — pure organizer default

- **A.1 preflight fail-closed naming both fields.** Read `internal/protocol/facilitator.go:59-71`
  and `internal/app/preflight.go:314-334`; the message names `facilitator:`, `participants:` and
  `facilitator_participates: true`. Re-ran the implementer's suite independently:
  `go test ./internal/app/ -run 'TestPreflightFacilitator|TestFacilitatorConflictGate|...' -count=1`
  → `ok  parley-deck-cli/internal/app  0.506s`. **Could not break. PRIMARY.**
- **A.3 declared facilitator never selected for any code role.** I attacked the predicate from the
  side the tests do not: `IMPLEMENTATION.md` frontmatter naming the facilitator as `implementer:`.
  `resolveImplementer` (`internal/consensus/consensus.go:1061-1089`) gates every candidate through
  `isParticipant(id)` against the list it is given, and `newDriverImplOps`
  (`internal/app/driver_impl.go:64`) passes the **`eligible`** slice, not `participants`. Implementer,
  reviewers and drafter all derive from `eligible` (`:72-93`). The bypass does not exist.
  **Could not break. PRIMARY.**
- **A.2 escalate-not-fallback.** Verified the deadlock error reaches every role entry point
  (`driver_impl.go:214, 278, 407, 504`) and that `internal/driver/impl.go:87-89` returns
  `ActionEscalated`. The shipped test asserts the literal phrase the implementation writes, so I
  treated it as mirroring and re-derived the property from the code paths. **Held. PRIMARY.**
- **A.4 prompt/gate parity — FAILED to hold as a no-new-obligation change.** The prompt repair
  itself is correct and genuinely single-sourced. The **gate half** is where the refutation
  succeeded: see CRIT-1. I built a deck whose `consensus.md` is verbatim the Phase-3 template from
  `parley-deck/COOPERATION.md:375-384` and ran the shipped binary. **Broken.**
- **A.5 absent-field byte-identical run plan.** `go test ./internal/runplan/ -count=1` green inside
  the full suite; `FacilitatorRoleFromMeta` returns the zero value on a missing key
  (`facilitator.go:44-48`), so `IneligibleForRoles` is `false` for every id. **Could not break.**
- **A witness `implementer:` ≠ facilitator.** `IMPLEMENTATION.md:4` `implementer: zcode-1`;
  `00-prompt.md:3` `facilitator: codex-1`. **Holds. PRIMARY.**

### B — `parley wait` + PhaseDigest

- **B exit 0 on a reached boundary — BROKEN on the live deck.** `parley wait --dir <review-worktree>
  --idea meta-protocol-change-lean-organizer --for round --timeout 20s` printed a correct
  `round-02: 3/3 filed-and-valid` digest and then exited **4** on
  `claude-1-to-user_fixup-budget_cap-exceeded-trajectory.md` — an escalation for a *different* idea
  (`meta-protocol-change-phase-packet-and-fixup-budget`), dated **2026-08-12**. See MAJ-1.
- **B missing ≠ invalid.** `invalidArtifact` gates on `row.Filed && !row.Valid`
  (`internal/app/wait.go:240`); a not-yet-filed row has `Valid=false, Validity=""` and is skipped.
  I could not construct a missing-file case that exits 4. **Could not break. PRIMARY.**
- **B present-but-invalid → exit 4 with the verbatim validator string.** Confirmed indirectly via
  the consensus path (T3 below) where the exact validator text
  `missing required consensus section(s): Drafter position changes, Alternatives disposition`
  was reproduced verbatim on stderr. **Holds.**
- **B timeout ceiling.** Read `trackTimeoutCeiling` (`wait.go:53-62`) → 5/15/30 min, which matches
  the §4.0 "Timeout per agent" row (`COOPERATION.md` §4.0 per-track table) exactly. Explicit
  `--timeout` above the ceiling is rejected (`:126-129`); a configured default above it is silently
  clamped (`:132-135`), which is the safe direction. `default < 30m && configurable` holds
  (`waitDefaultMS = 25*60*1000`, `wait.go:45`; `[defaults.timeouts] wait_ms`).
  **Could not break. PRIMARY.**
- **B `driver.error` → exit 4 "on arrival" — BROKEN.** Fixture with one historical `driver.error`
  followed by recovery and round completion → exit 4 forever. See MAJ-2.
- **B `--for implementation` — BROKEN for `ready-for-review`.** See MAJ-3.
- **B digest never rewrites a round file.** I ran `wait` repeatedly against the live deck and
  `git status --short` stayed empty throughout. **Could not break. PRIMARY.**
- **B byte-identical digest over an unchanged tree.** Two `organizer brief` runs (which embed the
  digest) were byte-identical (`cmp` silent). `PhaseDigest` carries no timestamp field
  (`phasedigest.go:81-88`). **Could not break. PRIMARY.**
- **B no model-written field.** `agentRow` discards the extracted position
  (`phasedigest.go:298` — `_, row.FellBack = extractPosition(...)`), so no prose survives into the
  digest. The property holds; the residue is MIN-1.

### C — audience packet, slim skill, brief

- **C.1 never-cut floor unbreachable by the audience dimension — HELD.** I wrote a hostile map
  appending `## 15. Verification integrity`, `### 15.1`, `### 15.2`, `### 15.7` and `## 6.` to
  `audiences.facilitator.omit`, then rebuilt at every kernel phase with **line-anchored** heading
  checks (an earlier `strings.Contains` pass of mine gave a false positive off the omission index;
  corrected). Every block stayed present. `neverCutForRequest` (`applicability.go:181-210`) is the
  guard. **Could not break the runtime floor. PRIMARY.**
- **C.1 `packet check` "proves it with a negative test" — PARTIALLY BROKEN.** Same hostile maps
  through `Check`: `## 6.` and `## 14.` are caught; `## 15.`, `### 15.1`, `### 15.7` are **not**
  (`OK:true`, `NeverCut:[]`). See MAJ-6.
- **C.1 unknown audience → full fallback with reason.** `Audience: "banana"` → `ContextMode=full`,
  body 109,663 B = the full source, `audience_fallback_reason=unknown-audience:banana`.
  **Could not break. PRIMARY.**
- **C.1 `facilitator_participates: true` → full context.** Same shape: `full`, 109,663 B, reason
  `facilitator-participates`. **Could not break. PRIMARY.**
- **C.1 attestation additive `audience`, never top-level `role`.** `Attestation` gains `audience` /
  `audience_fallback_reason` only (`packet.go:77-87`); no `role` key. **Holds.** Residual: MIN-4.
- **C.2/C.3 named-omission-set absence (the R-2 gating check) and retention set.** Re-ran
  `TestLiveDeckFacilitatorPacketNamedSetsAndGuardrail` at both `86d028b` and `3c97f44` — PASS. I
  separately confirmed the assertion is line-anchored (`headingLinePresent`), so it is immune to
  omission-index mentions — the right shape. **Could not break the gating check. PRIMARY.**
- **C.3 measurements — BROKEN as recorded.** The figures in `IMPLEMENTATION.md` reproduce at
  *neither* committed state. See MAJ-4 for the byte-exact cause.
- **C.3 `--optimize` baseline captured in the same run.** It is (`facilitator_packet_live_test.go:112-119`),
  but the recorded value is stale by the same defect. **Broken as recorded.**
- **C.4 SKILL.md core.** Measured independently: 17,038 B ≤ 20,000 B; all six driver commands named
  (`run`/`continue`/`wait`/`status`/`consensus`/`preflight`); all six files in `references/` linked
  from core; frontmatter block and `## Core Rule` byte-identical to `d1e57d5`. I re-derived the
  relocation myself over **all** heading levels (not just `##`): 76 original headings → 9 retained
  in core, 6 → `HEADLESS_LAUNCH.md`, 49 → `ARTIFACT_TEMPLATES.md`, 12 → `ROSTER_AND_PROTOCOL.md`,
  **0 dropped, 0 duplicated across references**. `npm ci && npm test` → exit 0.
  **Could not break. PRIMARY.**
- **C.5 brief ≤ 8,192 B and byte-identical ×2.** Live deck: 1,925 B; two runs `cmp`-identical.
  **Could not break. PRIMARY.**
- **C.5 brief "writes no file" — partially broken.** It writes a packet body outside the deck on
  every call. See MIN-3.
- **C.5 brief content correctness — BROKEN.** The phase it resolves is wrong for every phase ≥ 5.
  See MAJ-5.

### D — handoff, ledger, telemetry

- **D.2 "no flag accepts a token count".** I tried `--tokens --input-tokens --total-tokens --count
  --usage --input --output`; every one was rejected with `flag provided but not defined`.
  **Could not break. PRIMARY.**
- **D.2 one row, six fields verbatim, ≤ 1 KB stdout, idempotent, last-cumulative-wins.** Fixture
  rollout with two `token_count` events (12 then 120 total): stdout 107 B; ledger has exactly 1 row
  after two ingests, second reports `idempotent no-op`; `total_token_usage` keys are exactly
  `cache_write_input_tokens, cached_input_tokens, input_tokens, output_tokens,
  reasoning_output_tokens, total_tokens`; `total_tokens=120` (the last cumulative event, not the
  first, not a sum). **Could not break. PRIMARY.**
- **D.3 attribution by timestamp windows — BROKEN in practice.** See MAJ-7.
- **D.3 never slug scanning.** Neither parser reads the idea slug from the session file; attribution
  comes only from `--idea`/`--phase` plus `run.json` windows. **Could not break. PRIMARY.**
- **D.4 kimi telemetry honest degradation.** `kimiUsageRecord` (`telemetry/usage.go:163-188`)
  derives input as `inputOther + inputCacheRead`, sets `CostBasis: "unavailable"`, and returns the
  zero `Usage` when no bucket parses — no fabricated number. The three text consumers that the
  stream-json argv change breaks all unwrap (`preflight.go:928`, `runner.go:672`, `consult.go:153`).
  **Could not break the honesty property. PRIMARY.**
- **D.1 handoff record per transition.** `commitCursor` → `writePhaseHandoff`
  (`internal/driver/driver.go:222, 234-251`); schema doc embedded and states recomputation is
  authoritative (`phasehandoff.go:47-51`); write is atomic and non-fatal on failure.
  **Holds.** Residual: MIN-5.

### Cross-cutting

- **Three copies carry identical hunks.** I extracted all six idea hunks and asserted each is
  present in all three copies (`parley-deck/COOPERATION.md`,
  `internal/protocol/defaults/COOPERATION.md`, skill `references/COOPERATION.md`): 6/6 in 3/3.
  **Could not break. PRIMARY.**
- **Staged core `~/.parley/staging/COOPERATION-2.13.0.md` (109,507 B).** Verified independently and
  in full:
  - *Placeholder header*: lines 1-7 are byte-identical to the 2.10.0 template
    (`<workspace-name>`, `<transport-choice>`, `<YYYY-MM-DD>`), and the deck view's
    `**Protocol synced:**` line and concrete values are **absent**.
  - *Stub §2*: the handle table is the generated-view stub; `grep -c 'claude-1|kimi-1|zcode-1|codex-1'`
    over the whole staged core → **0**. No deck-view roster leaked.
  - *Exactly three change sets*: I diffed 2.10.0 → 2.13.0 into 9 hunks and attributed each by
    normalised signature against (a) `diff(core 2.10.0, staged 2.11.0)`, (b)
    `diff(v1.47.0, v1.48.0 internal/protocol/defaults/COOPERATION.md)`, (c)
    `diff(b37f7ef, 3c97f44 parley-deck/COOPERATION.md)`. Result: **2 from 2.11.0, 2 from 1.48.0,
    5 from this idea, 0 unexplained.** I then confirmed all **three** 1.48.0 changes are present by
    content (the LE-7/LE-11 bullet, the goal-check "can only withhold a close" paragraph, and §9
    item 1 — the last appears inside this idea's superseding §9 hunk, which is why it attributes to
    the idea rather than going missing).
  - *Matches current reviewed text*: the §15 region is **byte-identical** to the deck view
    (8,056 B both).
  **Could not break. PRIMARY.**
- **Attended publish gate.** The escalation note carries the exact command
  (`inbox/zcode-1-to-user_..._core-publish.md:36`:
  `parley protocol publish --version 2.13.0 --from ~/.parley/staging/COOPERATION-2.13.0.md`).
  I did not run it and did not manufacture a TTY. **Holds.**
- **Both suites green.** CLI: `go test ./... -count=1 -timeout 2400s` in my own worktree →
  **31/31 packages ok, exit 0**, 0 FAIL. Skill: `npm ci && npm test` → **exit 0**.
  Independently reproduced. **Holds.**

### On the implementer's phase-8 disposition (my position, as asked)

The disposition is: *"phase-8 facilitator packet 71,168 B exceeds 70,000 B, while FINAL's ratified
hard guardrail names phase 1 (58,190 B reported); map/cap unchanged and measurement brought to
quorum."*

**I concur with the scope half and reject the measurement half.**

- **The scope reading is correct.** FINAL C.3 states the guardrail as "≤ 70,000 B for the **phase-1
  / deliberation / github-pr** facilitator body", and the acceptance table repeats the phase-1
  framing. The map and the ceiling were not touched, and open item 2's required action on a
  non-fit — "show the bytes to the quorum before changing either the map or the ceiling" — was
  performed in `IMPLEMENTATION.md ## Deviations from FINAL.md`. Phase 8 exceeding 70,000 B is
  therefore **not by itself a FINAL violation**, and I do not file it as one.
- **The numbers are wrong, and both numbers in the disposition are wrong.** Reproducing
  `TestLiveDeckFacilitatorAcrossPhases` at *both* `86d028b` and `3c97f44` gives phase 1 =
  **58,941 B** (not 58,190) and phase 8 = **72,431 B** (not 71,168). A quorum asked to rule on a
  byte overage is being shown figures that do not reproduce at the commit it is ruling on. Filed as
  MAJ-4 with the byte-exact cause.
- **Two things the disposition does not say, which the quorum should have.** (i) Phase 7 measures
  **69,821 B** — **179 B** under the guardrail; the phase-1-scoped framing is one protocol sentence
  away from a second overage. (ii) `TestLiveDeckFacilitatorAcrossPhases` only `t.Logf`s
  (`facilitator_packet_live_test.go:160`); it cannot fail. Open item 2 asks whether the map fits
  "**at each phase**", and that question is currently answered by a log line no CI gate reads.

## Findings

### [CRITICAL] CRIT-1 — New consensus-section gate rejects the protocol's own Phase-3 template and 79 of 80 existing consensus artifacts

`internal/consensus/consensus.go:541-549` adds a hard gate requiring all seven of
`protocol.RequiredConsensusSections` (`internal/protocol/consensussections.go:17-25`) in every
design `consensus.md`. No such gate existed at `b37f7ef`. Three separate problems:

**1. It rejects the protocol's own published template.** I built a deck whose `consensus.md` is
verbatim the Phase-3 template at `parley-deck/COOPERATION.md:375-384`, with both participants at
`✅ ACCEPT`:

```
$ parley consensus status --dir /tmp/claude1-review/fx-protocol-template tpl-deck
Consensus: malformed
Errors:
  missing required consensus section(s): Drafter position changes, Alternatives disposition
Signoffs:
  alpha-1    ✅ ACCEPT
  beta-1     ✅ ACCEPT
```

The shipped template does not contain either heading. The same fixture against the baseline binary
built from `b37f7ef` returns `Consensus: ready`.

**2. It contradicts two explicit sentences in the same file, which were not updated.** Both survive
verbatim at HEAD:
- `parley-deck/COOPERATION.md:401` — "The `## Comparison & blind spots` section is an **advisory
  drafting discipline**, not a gate: append-only signoffs remain the only gate".
- `parley-deck/COOPERATION.md:1360` (§15.6) — "Only (a) is machine-validated today — (b) and (c)
  bind by discipline, **and this sentence says so rather than implying a gate that does not exist**."
  `## Alternatives disposition` is §15.6(c); it is now gated.

Additionally §15.5 conditions `## Drafter position changes` on "when the facilitator **is also a
participant** and drafts `consensus.md`"; the gate requires it unconditionally.

**3. Blast radius, measured with the real binary over this repo's own deck:**

| | base `b37f7ef` | HEAD `3c97f44` |
|---|---|---|
| protocol's own Phase-3 template | `ready` | `malformed` |
| repo deck `consensus.md` files | **9 / 80** malformed | **79 / 80** malformed |

The single survivor is this idea's own `consensus.md`. `malformed` is blocking:
`internal/driver/consensus.go:50` and `internal/runplan/runplan.go:188` advance only on
`Ready`/`Reserved`, and it propagates into `parley wait` (see MAJ-1 context / T3).

This is a **new mandatory obligation**, while FINAL's owner stop-rule table records **"No"** for
row A, and it is broader than FINAL A.4, which asked that the *prompt* emit the sections with "a
parity test proving prompt and gate read the same value" — not that the gate begin hard-requiring
the §15 duty sections and the explicitly-advisory `## Comparison & blind spots`.

The repo states the correct discipline for exactly this situation two files away, and it was not
followed here — `internal/protocol/reviewartifact.go:48-49` ("**Measured before enforcing**: 348 of
539 review artifacts already carry it. The 191 that do not are historical and are not revalidated;
this binds new reviews") and `:92-95` ("a gate that rejects live work is a worse defect than the one
it fixes").

**Suggested fix (smallest change that keeps A.4's real property):** keep
`RequiredConsensusSections` as the single source for the **prompt and the scaffold**, and split the
**gate** list so it stays at the pre-existing set. If the quorum does want the §15 duty sections
gated, that is a §7 protocol change: it must amend `COOPERATION.md:401` and `:1360`, land in all
three copies and the staged core, and adopt the `reviewartifact.go` pattern (measure first, bind
new artifacts only). Either way `MissingConsensusSections`
(`consensussections.go:39-47`) should use `HasHeadingLine` rather than `strings.Contains` — today
`### Agreed decisions` satisfies a `## Agreed decisions` requirement by substring.

### [MAJOR] MAJ-1 — `parley wait` treats any `-to-user_` filename anywhere in the inbox as a blocking escalation, so it can never return 0 on this deck

`internal/app/wait.go:250-266`. `blockingEscalation(root)` scans the whole deck inbox and returns
true for any non-archived `*.md` whose **name** contains `-to-user_`. It does not scope to the idea,
does not read `blocking:`, does not read `status:`, and has no notion of "new".

```
$ parley wait --dir <review-worktree> --idea meta-protocol-change-lean-organizer --for round --timeout 20s
round-02: 3/3 filed-and-valid           # boundary WAS reached
wait: blocking escalation (unanswered to-user inbox note): claude-1-to-user_fixup-budget_cap-exceeded-trajectory.md
$ echo $?
4
```

That note is `idea: meta-protocol-change-phase-packet-and-fixup-budget`, `date: 2026-08-12` — a
different idea, six weeks old. Nine `-to-user_` notes sit in this inbox. Worse, `86d028b` itself
adds `inbox/zcode-1-to-user_meta-protocol-change-lean-organizer_core-publish.md`, whose frontmatter
says **`blocking: no`** — so the commit that ships `parley wait` also guarantees `wait` returns 4 for
this idea permanently. FINAL B.3 scopes this to a "**new unanswered**" escalation and the acceptance
table to an "unanswered `to-user` escalation"; none of the three qualifiers is implemented.

*Independent corroboration and my position on codex-1's observation.* I reached this finding by my
own execution before reading
`inbox/codex-1-to-all_meta-protocol-change-lean-organizer_wait-observation.md` (dated 2026-09-24),
which reports the organizer hitting the identical failure from normal use —
`parley wait --idea meta-protocol-change-lean-organizer --for review --timeout 25m --json` → exit 4
on the same unrelated `claude-1-to-user_fixup-budget_...` note. My verdict stays **PRIMARY** on my
own runs quoted above; codex-1's note is independent operational testimony that converges on the
same defect, not the basis for my verdict. **My position, as that note asks: the behaviour does not
meet FINAL B's contract.** FINAL B.3 scopes the early return to a "new unanswered" escalation, and
the acceptance table to an "unanswered `to-user` escalation"; a six-week-old note belonging to a
different idea is neither new, nor unanswered-for-this-idea, nor this idea's. The note correctly
suppresses nothing, and I treat it as concurring evidence rather than a disposition to weigh.

**Suggested fix:** filter to notes whose frontmatter `idea:` matches the awaited slug **and**
`blocking:` is not `no` **and** `status:` is not answered/resolved; and only fire for notes whose
mtime is after `wait` started (the "arrives"/"new" semantics), so a pre-existing note is reported in
the digest rather than being a hard exit.

### [MAJOR] MAJ-2 — a single historical `driver.error` poisons `parley wait` forever

`internal/app/wait.go:281-299`. `driverErrorEvent` loads the entire event log and returns the
**first** `driver.error`/`run.error` it finds, with no "since this wait began" boundary. Fixture: one
`driver.error` at the top of the log, then `run.phase` and `round.completed`:

```
$ parley wait --dir /tmp/claude1-review/fx-wait --idea wfx --for round --timeout 11s
wait: driver.error event: transient: context canceled
$ echo $?
4
```

The round is complete and the run recovered. Any run that ever hit a transient error can never use
`wait` again. This exact situation occurred in this run
(`inbox/claude-to-user_..._driver-error.md`: "draft FINAL.md: context canceled", after which the run
continued). FINAL B.3 says the exit fires when an event "**arrives**".

**Suggested fix:** record the event-log length (or newest event time) at wait start and only treat
`driver.error`/`run.error` entries after that point as blocking; surface earlier ones as a digest
annotation.

### [MAJOR] MAJ-3 — `wait --for implementation` never returns for `status: ready-for-review`, and the timeout names nobody

`internal/driver/phasedigest.go:211-216`. The `ReadyForReview` switch lists `implemented`,
`complete` and `fix-up-cycle-1..5` but **not** `ready-for-review` — even though
`protocol.ValidImplementationStatus` accepts it and `internal/protocol/reviewartifact.go:89-101`
records it as a status "4 files" in this deck actually use.

```
$ parley wait --dir /tmp/claude1-review/fx-wait --idea wfx --for implementation --timeout 11s
implementation: present=true status=ready-for-review implementer=alpha-1
next: await implementation
wait: timeout after 11s; outstanding: none named — inspect the digest
$ echo $?      # 3
```

"outstanding: none named — inspect the digest" (`wait.go:358`) is precisely the *status-conceals-the-
reason* class that FINAL cites B as fixing.

**Suggested fix:** derive `ReadyForReview` from `protocol.ValidImplementationStatus` minus the
in-progress states rather than a second hand-written list, and have `outstandingAgents` name the
blocking condition (e.g. "implementation status `ready-for-review` is not a recognised
ready state") instead of "none named".

### [MAJOR] MAJ-4 — every R-2 byte measurement recorded in `IMPLEMENTATION.md` is pre-hunk and does not reproduce at either committed state

`IMPLEMENTATION.md:98-105, 242-245` and Decision Log `:148-150`. Re-running the implementer's own
tests at **both** `86d028b` and `3c97f44` (identical results; the only diff between those commits is
one line of `IMPLEMENTATION.md`):

| phase | recorded | reproduced at `86d028b` and `3c97f44` | delta |
|---|---:|---:|---:|
| 0 | 52,854 | 53,605 | +751 |
| **1 (the hard-asserted guardrail scope)** | **58,190** | **58,941** | **+751** |
| 2 | 58,291 | 59,042 | +751 |
| 3 | 59,551 | 60,302 | +751 |
| 4 | 62,469 | 63,220 | +751 |
| 5 | 59,791 | 60,864 | +1,073 |
| 6 | 61,005 | 61,946 | +941 |
| 7 | 68,880 | 69,821 | +941 |
| **8 (the disposition)** | **71,168** | **72,431** | **+1,263** |
| `--optimize` baseline | 64,734 | 65,485 | +751 |

The cause is byte-exact. The five P.1 protocol hunks added to `parley-deck/COOPERATION.md` (which
the test reads as its authority, `facilitator_packet_live_test.go:20`) measure:
Quickstart facilitator row **+202**, §4 Phase 5 **+322**, §4 Phase 6 **+190**, §9 item 1 **+346**,
§11 advisory **+203** (total +1,263, matching `108,400 → 109,663 B`). The facilitator view always
retains Quickstart, §9 and §11 → **202+346+203 = 751** for every phase; phase 5 additionally pins the
Phase-5 subsection → **+322 = 1,073**; phases 6-7 pin the Phase-6 subsection → **+190 = 941**; phase 8
pins both → **+1,263**. Every observed delta is accounted for exactly. **The measurements were taken
before the P.1 hunks landed in the deck copy.**

FINAL open item 3 requires the baseline be "captured in the same test run, **not transcribed**", and
R-2 requires the body be "recorded against both the floor and the guardrail in the same test run".
The test does capture correctly; the transcription into `IMPLEMENTATION.md` is stale, and it is what
the quorum is reading.

**Suggested fix:** re-run `go test ./internal/app/ -run TestLiveDeckFacilitator -count=1 -v` at the
fix-up HEAD and replace every figure in `## Deviations from FINAL.md`, `## Validation evidence` and
the Decision Log with the reproduced values (phase 1 **58,941 B**, phase 8 **72,431 B**, `--optimize`
**65,485 B**), noting that phase 7 sits 179 B under the guardrail.

### [MAJOR] MAJ-5 — `parley organizer brief` resolves the wrong protocol phase for every phase ≥ 5, and yields a §15-free facilitator view for two live status values

`internal/app/organizer.go:29-44`. `briefPhase` maps the idea's `00-prompt.md` `status:` to a packet
phase and covers only `round-01`, `round-*`, `consensus`, `final`, `review*`. Measured against the
`status:` vocabulary actually present in this deck:

| `status:` | count in deck | packet phase |
|---|---:|---:|
| `final` | 80 | 4 |
| `complete` | 5 | **0** |
| `open` | 3 | 0 |
| `round-01` / `round-02` | 3 | 1 / 2 |
| `implementation` | 2 | **0** |
| `abandoned` | 1 | **0** |
| `review` (the only phase-6 trigger) | **0** | 6 |

So the organizer re-orienting during Phases 5-8 — exactly the phases the pure-organizer default
exists for — is handed the **phase-4** packet (status `final`, 80 ideas) or the **phase-0** packet.
Phase 0 is not a kernel phase, so its facilitator body carries **no §15 at all**. Reproduced:

```
$ parley organizer brief --dir /tmp/claude1-review/fx-brief --idea impl-idea   # status: implementation
facilitator body: .../packet-phase0-deliberation-a4aaf412....md
  ABSENT   ## 15. Verification integrity
  ABSENT   ### 15.1 Scope, ownership, location
  ABSENT   ### 15.7 Per-track binding
  PRESENT  ## 4. Protocol — phases of an idea
```

The never-cut floor is not breached (phase 0 is outside the kernel set), so this is a phase-
resolution defect, not a floor defect — but FINAL C.2 names §15 in the retention set and C.5 makes
this brief the documented post-compaction re-orientation surface. The current live deck also shows
it: the brief for this idea renders `packet-phase4-...` while the idea is in Phase 6.

**Suggested fix:** derive the phase from the driver cursor / run record (`runPhasePointer` already
reads `run.json` `phase`, `organizer.go:129-144`) and fall back to `status:` only when no run exists;
extend the mapping to `implementation`→5, `review*`→6, `fix-up-cycle-*`→8, `complete`→8.

### [MAJOR] MAJ-6 — `packet check`'s audience negative test cannot fire for the phase-pinned never-cut blocks, i.e. for §15

`internal/protocolpacket/applicability.go:382-402`. The map-level scan rejects an audience omitting a
never-cut block only when `nc.always` (`:396`). The runtime scan (`:413-426`) only reports blocks
`Build` actually omitted — and `Build` refuses, so it never reports. Net effect, measured by
appending one omit entry at a time to the live map:

| omit entry added to `audiences.facilitator` | `packet check` catches | runtime floor holds |
|---|---|---|
| `## 6. Conflict-avoidance mechanics` (always) | **yes** | yes |
| `## 14. …human brake` (always) | **yes** | yes |
| `## 15. Verification integrity` (phases) | **no** | yes |
| `### 15.1 Scope, ownership, location` (phases) | **no** | yes |
| `### 15.7 Per-track binding` (phases) | **no** | yes |

FINAL C.1 states the property as "it can never cut below the ratified never-cut floor …, and
`packet check` **proves it with a negative test**". The runtime guarantee is intact — this is the
proof obligation that is not met, for exactly the section FINAL singled out ("not a preference: the
ratified never-cut floor pins `## 15.` …"). The code comment at `:388-393` justifies allowing
conditionally-pinned blocks to be named because the facilitator set must name the non-active §11
subsections; that reasoning applies to §11 and does not need to extend to §15.

**Suggested fix:** narrow the allowance to transport-conditional blocks — reject an audience omit
entry naming a `phases:`-pinned never-cut block (§15.x) at map level, the way `always` blocks are
already rejected.

### [MAJOR] MAJ-7 — client-accounting attribution can only ever return `ambiguous`: every real run record has a zero-width window

`internal/app/usage_ingest.go:264-307`. `resolveAttribution` requires the accounting event to satisfy
`!when.Before(created) && !when.After(updated)` over `run.json`'s `created_at..updated_at`. Measured
across **every** run record in this deck:

```
20260602T100714…  window=0.0s      20260602T191449…  window=0.0s
20260602T195452…  window=0.0s      20260923T202501…  window=0.0s   (this run)
```

`created_at == updated_at` in all four. The lean-organizer run spans 20:25→23:44 per
`organizer-usage.md`, yet its `run.json` still reports `updated_at = 2026-09-23T20:25:01.377412Z`.
`internal/runstate/runstate.go:209-210` only **reads** `UpdatedAt`; `internal/runmanifest/manifest.go:175,197`
writes it once at creation; nothing rewrites `run.json` with a later value (the many
`c.UpdatedAt = nowRFC3339()` sites are the **cursor**, not the run manifest). So the corroboration
arm of D.3 is structurally dead and *all* usage is residue.

The shipped test does not catch this because it hand-writes a manifest shape the driver never
produces — `usage_ingest_test.go:161`:
`{"idea_slug":"u","created_at":"2026-09-23T19:00:00Z","updated_at":"2026-09-23T21:00:00Z"}` (a
2-hour window). The implementer's live smoke test reporting `attribution=ambiguous` is consistent
with this; the stated reason ("no run window covers the final event") understates the cause.

To be fair to the design: `ambiguous` is the **honest** fail-safe FINAL asked for, and no number is
fabricated. The finding is that the mechanism cannot succeed against real artifacts, which makes the
label uninformative rather than selective.

**Suggested fix:** either advance `run.json`'s `updated_at` at each `commitCursor` (the handoff write
is already there), or widen the window to `created_at .. max(updated_at, newest event time in
events.jsonl)`. Then replace the synthetic-manifest test with one built by the driver itself.

### [MAJOR] MAJ-8 — FINAL's required blind-spot (i) measurement was not performed; measured here, §15.5/§15.6 are absent from the facilitator view at phases 5 and 8

FINAL `## Context & orientation`, "Known blind spots flagged for implementation (**measure, do not
silently assume**): (i) §15.5/§15.6 are never-cut at phases {1,2,3,6,7} but **not pinned at phases 5
and 8** — check what a facilitator view loses there and whether it matters." `IMPLEMENTATION.md`
records no such measurement.

Measured (line-anchored, live map + live protocol):

```
phase 0 facilitator: absent §15 blocks = [15. … 15.1 … 15.2 … 15.3 … 15.4 … 15.5 … 15.6 … 15.7]
phase 1,2,3,4,6,7    absent = []
phase 5 facilitator: absent = [15.5 Role concentration, 15.6 Alternatives and correlated agreement]
phase 8 facilitator: absent = [15.5 Role concentration, 15.6 Alternatives and correlated agreement]
```

The live map gives `### 15.5` / `### 15.6` `phases=[1 2 3 4 6 7]`. So at Phase 8 — fix-up, where the
facilitator adjudicates — the facilitator's reading set omits §15.5 ("procedural calls … are
**provisional** until the corresponding signoff gate passes. The signoffs, not the facilitator's
judgment, are the close") while **retaining** §15.7, whose table asserts both rules bind on every
track. The facilitator is shown a table of duties whose text it cannot read.

**Suggested fix:** add phases 5 and 8 to the `### 15.5` / `### 15.6` map entries (a §7 change to
`packet-applicability.yaml`, +~2.4 KB at those phases — note this pushes phase 8 further above the
guardrail, which is quorum information), or record the measurement with an explicit quorum decision
that the loss is acceptable. Either way the required check should appear in `IMPLEMENTATION.md`.

### [MINOR] MIN-1 — `fell_back` reports the degradation of a computation the digest discards, and is `true` for every valid round-2 artifact

`internal/driver/phasedigest.go:298`. `extractPosition` returns `fell=true` when no `## Summary`
heading exists; PhaseDigest discards the position (correctly — no model-written field) but keeps the
flag. Round-2 artifacts legitimately use other headings, so the live digest shows
`fell_back=true` for all three `valid=true` round-02 rows, suggesting degradation where none exists.
`grep -c "^## Summary$"` → 0 for each of the three files. FINAL B.1 requires the column, so it must
stay; its meaning should change. **Fix:** set `fell_back` from the validator/ownership path (e.g.
frontmatter unreadable, validator fell back) rather than from the discarded summary extraction, and
document it in the column legend.

### [MINOR] MIN-2 — the next-action line says "await implementation" when the implementation is already published

`internal/driver/phasedigest.go:184-192`: the `if imp := …` block returns only for `complete`, and
both remaining branches return `NextAwaitImplementation`, so the two are indistinguishable. Live
output shows `implementation: present=true status=implemented` beside `next: await implementation`.
It stays inside the fixed enumeration, so the enumeration test passes. **Fix:** return
`NextAwaitReviewArtifact` when the implementation is present and ready for review but no review
round exists yet, and delete the dead branch at `:189-192`.

### [MINOR] MIN-3 — the brief writes a packet body on every invocation, so "computed, never stored / writes no file" is not literally true

`internal/app/organizer.go:106` calls `protocolpacket.Render`, which writes
`<root>/.parley-runtime/protocol-packets/packet-phase<N>-<track>-<sha>.md`. Reproduced: deleting
`.parley-runtime` and running the brief once recreates it. The acceptance criterion is scoped to a
read-only **deck** and the shipped test (`organizer_test.go:50-81`) snapshots only
`parley-deck/`, so the criterion as written is met — but the brief prints "Computed view — never
stored" (`organizer.go:146`) while leaving a file behind. **Fix:** either state the runtime cache in
that line and in FINAL's wording, or add a render mode that returns the body without persisting it.

### [MINOR] MIN-4 — `--optimize` with an unrecognised audience stamps the rejected audience into the packet header

`internal/protocolpacket/packet.go:469-478` builds the header from `req.Audience` (requested) rather
than the resolved audience. With `Optimize: true, Audience: "banana"` the build is `ModePacket`,
`ctx.Audience` is correctly `""` and `audience_fallback_reason=unknown-audience:banana`, but the body's
first line reads `… flags=- audience=banana -->` though no audience omission was applied. The body is
what the agent reads and what `packet_sha256` covers, so the attestation embedded in the artifact
disagrees with the attestation in the JSON. **Fix:** pass the resolved audience into `renderPacket`.

### [MINOR] MIN-5 — `BuildPhaseHandoffRecord` returns `RunID: "parley-deck"`, and the round-trip loader covers 4 of 9 fields

`internal/driver/phasehandoff.go:58` —
`filepath.Base(filepath.Dir(filepath.Clean(ideaDir + "/..")))` evaluates to `parley-deck` for any
standard layout; the comment concedes it is "replaced by caller context below when possible". The
shipped caller does patch it (`internal/driver/driver.go:236`), so the written records are correct,
but the exported constructor is wrong for any other caller and the test only asserts `RunID != ""`
(`phasehandoff_test.go:36`). Separately, `LoadPhaseHandoffRecord` (`phasehandoff.go:91-115`) restores
only `run_id/idea/phase/action/next` — `previous_phase`, `round_label`, `idea_status`, `written_at`
and the embedded `Digest` are dropped, so the "round-trip" is partial. **Fix:** take the run dir as a
parameter and drop the derivation; parse the remaining frontmatter keys in the loader.

### [MINOR] MIN-6 — the recorded task-local binary hash was built from a dirty tree at the pre-implementation base and is not reproducible

`IMPLEMENTATION.md:289-291` records `/tmp/parley-lean-organizer/parley` sha256
`8d3e4c1c684023c03aaf26d2f47d52be27eaefa3693e5bc45db5aca0feca1dcb`. `go version -m` on that file
reports `vcs.revision=b37f7ef9dd0941b21c3ba146c91d1a118259cbed`, `vcs.time=2026-09-23T21:47:22Z`,
`vcs.modified=true` — i.e. the **pre-implementation** commit with uncommitted changes. My build from
clean `3c97f44` reports `vcs.revision=3c97f44…`, `vcs.modified=false` and sha256
`abaad23610c874ec827cf00255285d1f03212544c6994367f305562cc846701b`. Under §15.2 a recorded hash
should be reproducible from a named tree; this one is not. **Fix:** rebuild from the fix-up HEAD in a
clean tree and record that hash with `vcs.revision`/`vcs.modified`.

### [MINOR] MIN-7 — `consensus.ResolveImplementer` is exported with zero callers, and its doc and the Decision Log both describe behaviour it does not have

`internal/consensus/consensus.go:1097`. `grep -rn "consensus.ResolveImplementer" internal/` returns
nothing (`ExpectedRoundParticipants` does have one caller, `phasedigest.go:249`). `IMPLEMENTATION.md`
Decision Log `:151-153` states "PhaseDigest calls `consensus.Status` and the newly exported
`consensus.ExpectedRoundParticipants` / `consensus.ResolveImplementer`" — `implSection` reads
`meta["implementer"]` from frontmatter directly instead (`phasedigest.go:208`). The doc comment at
`:1094-1096` also claims an "else `participants[0]`" fallback that `resolveImplementer` does not
implement (it returns `""`). **Fix:** either use it in `implSection` or drop the export, and correct
the Decision Log entry and the doc comment.

### [NIT] NIT-1 — hand-rolled byte search instead of the stdlib

`internal/app/usage_ingest.go:227-243`: `containsBytes`/`indexOfBytes` reimplement `bytes.Contains`
with a naive O(n·m) loop, run over every line of a 228 MB file. `bytes.Contains` is both shorter and
substantially faster. (§15.6(a) discipline: name what the toolchain already ships.)

### [NIT] NIT-2 — a dead import kept alive by a sentinel

`internal/app/wait.go:371-372`: `var _ = strconv.Itoa` with the comment "keeps the import list stable
for future numeric options". `strconv` is otherwise unused; remove both.

### [NIT] NIT-3 — the `--optimize` guard comment and the audience path share one branch without a test seam

`internal/protocolpacket/packet.go:428-431`: `if req.Optimize || audience != ""` merges the
experimental optimizer path and the ratified audience path. The comment claims "every guard that
governs `--optimize` governs it identically", which is true today; a short assertion that the two
produce the same `reasons` set would keep it true. (Related to MIN-4, which is the one place they
already diverge.)

## Open questions

1. **CRIT-1 disposition.** Is the hard consensus-section gate intended scope at all? FINAL A.4 asks
   only for prompt/gate **parity** over one constant. If the quorum wants the §15.5/§15.6 duty
   sections gated, that is a §7 protocol change and must amend `COOPERATION.md:401` and `:1360` in
   all three copies **and** the staged core before release — which would also change the staged
   core's "exactly three change sets" property I verified. Which way does the quorum want this
   resolved?
2. **MAJ-8 vs the guardrail.** Pinning §15.5/§15.6 at phases 5 and 8 is the protocol-correct fix, but
   it adds ~2.4 KB at exactly the phase already over 70,000 B (72,431 B → ~74.8 KB). FINAL says the
   bytes come back to the quorum before any map or ceiling change. Does the quorum prefer (a) accept
   the larger phase-8 body, (b) re-scope the guardrail explicitly per phase, or (c) accept the §15.5/
   §15.6 loss at 5 and 8 with a recorded rationale?
3. **Phase-7 headroom.** Phase 7 is 69,821 B — 179 B under the guardrail. Should
   `TestLiveDeckFacilitatorAcrossPhases` assert per-phase (currently `t.Logf` only,
   `facilitator_packet_live_test.go:160`), so open item 2's "at each phase" question has a gate
   rather than a log line?
4. **MAJ-7 direction.** Is the right repair to make `run.json` maintain a real `updated_at`
   (touching driver state shape), or to widen the window from `events.jsonl`? The first is more
   correct; the second is contained inside D.
5. **`parley wait` usability bar.** With MAJ-1/MAJ-2 fixed, `wait` still exits 4 on any deck holding
   an unanswered escalation for the awaited idea — including this idea's own core-publish note,
   which stays in the inbox until the owner publishes. Is "blocking" frontmatter the intended
   discriminator, or should `wait` report it in the digest and keep waiting?
6. **Scope boundary for the fix-ups.** MAJ-3 (`ready-for-review`) and MIN-2 (next action) are
   arguably pre-existing-vocabulary bugs surfaced by B rather than B defects. I filed them against B
   because B is the first consumer that makes them user-visible; the quorum may prefer to defer them.

---

**Canonical path:**
`parley-deck/ideas/meta-protocol-change-lean-organizer/review/round-01/claude-1.md`

**Severity counts:** CRITICAL 1 · MAJOR 8 · MINOR 7 · NIT 3 (19 findings).

**Review verdict:** ❌ **BLOCK** — CRIT-1 must be resolved before this can be called complete or
released. It is a shipped gate that invalidates 79 of 80 existing consensus artifacts and the
protocol's own published Phase-3 template, and it contradicts two sentences of the COOPERATION.md in
the same commit. MAJ-1 and MAJ-2 additionally make `parley wait` — scope B's headline deliverable —
unable to return 0 on the deck it was built for. I am not calling the implementation complete and I
am not releasing anything. Scope A's role predicate, C's runtime never-cut floor, C.4's skill
relocation, D.2's ledger contract and the staged core all held under my attempts to break them.

**Tested commits:** CLI `3c97f445336a70bd7bd57eb15a0da78144bcd255` (with comparison builds at
`86d028b5558a60858534ac72f25537e52ed395fc` and `b37f7ef9dd0941b21c3ba146c91d1a118259cbed`); skill
`0f513f6c0fc78523d11321e60868942a08079d1d` (base `d1e57d5c8a56ac40fc0c8b4ac035069d2c0b9283`).
