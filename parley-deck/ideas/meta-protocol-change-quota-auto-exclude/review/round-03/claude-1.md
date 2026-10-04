---
agent: claude-1
idea: meta-protocol-change-quota-auto-exclude
review-round: 3
date: 2026-10-04
reviewed-commit: 78ac5367cff0d9a2f9ae68448eb670489ec8ed29
skill-commit: dc85b5362b4872cc719a8ed745ba2a6cc89cfad6
responding-to: [review/round-01/claude-1.md, review/round-02/claude-1.md]
---

## Summary

This is the first full-scope Phase 6 review. Rounds 01 and 02 were early provenance feedback only. It covers
both stages at CLI `78ac536` (product diff since FINAL `27e42b8`) and the skill at `dc85b53` (diff since
`a5664d8`). The design core is well built and its tests pass on local disk:

- the whole-batch floor and role guards;
- immutable hash-bound history;
- current-versus-known signers;
- read-only status, wait and organizer brief views.

The implementation is **not ready to merge**. Two CRITICAL defects reproduce on this deck's own volume
(`/Volumes/My Shared Files`, AppleVirtIOFS) and in the default configuration:

1. Every new idea fails to be created there, with the knob on or off.
2. Every new idea, again with the knob on or off, can no longer take an owner-confirmed manual §9.0
   exclusion or a catch-up join without becoming "contradictory".

Four MAJOR findings follow:

- Owner-confirmed re-inclusion, veto withdrawal and scope widening have no path.
- The idea lease is not exclusive on this volume.
- The zcode recognizer accepts raw-quoted records and several mixed other-error forms.
- Manual-mode signoff handoffs lose the lease.

Counts: CRITICAL 2, MAJOR 4, MINOR 5, NIT 4. This review signs nothing.

### Protocol context and evidence base

- **Protocol packet.** `parley protocol packet --dir . --phase 6 --track deliberation --idea
  meta-protocol-change-quota-auto-exclude --flag protocol_change --json` gave:
  - `context_mode=full`;
  - `source_sha256=73613f95aabff8ed57372a535913985554014b71911cc8126ec3f7890a9d4f4e`;
  - `packet_sha256=73613f95aabff8ed57372a535913985554014b71911cc8126ec3f7890a9d4f4e`;
  - `fallback_reason` absent.

  I read the emitted body's Phase 6 and §15 sections; its sha256 matches. FINAL.md was read in full, and
  IMPLEMENTATION.md, the producer evidence and both diffs were read. The CLI diff was read in full for
  protocol, quota, telemetry, membership, consensus and the signoff and lock paths, and the remaining
  callers were spot-checked.
- **Provenance labels.**
  - PRIMARY: my own executed check, or my direct reading at the cited locator on `78ac536`.
  - SECONDARY: producer logs or tests I did not rerun, or unverified leads from my own read-only search
    helpers.
  - RECALL: none used.
- **Helpers.** Two read-only search helpers supplied leads. Every finding below was re-checked at its
  locator unless it is labeled SECONDARY.
- **Probe programs.** All temporary probes live under the git-ignored `.parley-runtime/claude1-r3/`. No
  tracked file, implementation, protocol, signoff or other participant file was edited. No provider was
  invoked.

## User direction

Quoted verbatim from `parley-deck/inbox/user-to-codex-1_meta-protocol-change-quota-auto-exclude_scope-reset-answer.md`
(body below its frontmatter):

> ## Owner answers to `codex-1-to-user_meta-protocol-change-quota-auto-exclude_scope-reset.md`
>
> Relayed by the owner's Claude Code session on 2026-10-04 at about 02:35 CEST. The relay asked both questions
> in Slovak and gave three options for question 1: A (build a reliable zcode channel, codex-1's
> recommendation), B (defer, no supported adapter) and a pragmatic option. The selected options are below,
> verbatim. A translation follows each one.
>
> **Question 1, AC2 and zcode support.** Selected: **"Pragmaticky: stderr zcode stačí"** ("Pragmatic: zcode's
> stderr is enough"). The option read (Slovak, verbatim): "Moje odporúčanie, ide o najrýchlejšiu funkčnú
> cestu. Pri zcode stačí JSON `responseBody` v stderr, ak proces skončil chybou, nenechal výstup a každý
> záznam o chybe je 429 „Limit Exhausted“ s `reset_at`. Je to výslovná výnimka z FINAL. Riziko: chyba
> sub-agenta by mohla vyradiť agenta, ktorý zlyhal z iného dôvodu. Tlmia to minimum 2 a notifikácia."
>
> Translation: "My recommendation, the fastest path that works. For zcode, the `responseBody` JSON in stderr is
> enough when the process ended with an error, left no output, and every error record is a 429 'Limit
> Exhausted' with `reset_at`. This is an explicit exception to FINAL. Risk: a sub-agent's error could exclude an
> agent that failed for a different reason. The minimum of 2 and the notice soften that."
>
> **Question 2, display time without a timezone.** Selected: **"Áno, podľa návrhu claude-1 (Recommended)"**.
> This adopts the bounded interpretation exactly as your note states it. A display clock counts only when it
> appears together with a machine reset value (`reset_at` or `retry_after`) in the same record and agrees with
> it within one second. Missing, contradictory or display-only resets still gate.
>
> ## What this authorizes (an owner-directed deviation from FINAL; record it in IMPLEMENTATION.md)
>
> - The zcode adapter becomes a **supported** recognizer, but only under this owner-defined evidence rule. All
>   of the following must hold:
>   - the zcode process exited non-zero;
>   - it produced no valid artifact, and no later attempt in the batch succeeded;
>   - stderr contains at least one provider error record whose `responseBody` JSON is a 429 with explicit
>     exhaustion text ("Limit Exhausted", or allowance semantics as in FINAL §4.4) and a machine reset value;
>   - **every** provider error record in that stderr agrees: each is that same exhaustion class, and their
>     reset values agree within the tolerance. A mixed or contradictory record set gates;
>   - the run ends with zcode's turn-failure line;
>   - the reset clears FINAL's 60-minute threshold.
> - Positive fixtures are this run's recorded zcode stderr (evidence incidents 1 and 2). Adversarial fixtures
>   include a quoted 429 in assistant or tool text, a mixed 429 plus other-error stderr, a reset under 60
>   minutes, a display-clock-only reset and a success after a 429.
> - Everything else in FINAL is unchanged: the floor of 2, the role guards, fail-closed, the record and
>   notice, both stages and claude-1's binding full review. Other adapters stay diagnostic-only unless they
>   have native evidence.
> - Record the deviation in `IMPLEMENTATION.md` with this note quoted, and carry the protocol wording into the
>   §9.0 hunk under §7. The claude-1 review checks the deviation as implemented, not whether to make it.

I take the bounded zcode stderr channel and the display-clock rule as owner decisions. I review only whether
the code follows them. The accepted residual risk covers a subagent's provider error that comes with an
unrelated root failure. It does not cover quoted assistant or tool text, or a mixed other-error stderr. The
owner listed both of those as adversarial fixtures that must gate.

## Refutation attempts

**Checks I ran on `78ac536` (PRIMARY):**

- `go build ./...` and `go vet ./internal/...` passed (`/tmp/claude1-r3/build.log`, `vet.log`).
- `go test -count=1` passed for quota, telemetry, membership, consensus, runcontrol, runstate, runmanifest,
  store, protocol and config. The protocol package run includes `TestEmbeddedDefaultMatchesLiveDeck`.
- In driver, runner and app, all `TestQuota*` tests passed, as did
  `TestLiveDeckFacilitatorPacketNamedSetsAndGuardrail`: facilitator body 69,966 B against the unchanged
  70,000 B guard.
- `cmp` of the deck `COOPERATION.md` against the skill `references/COOPERATION.md` returned exit 0.

**Producer runs I did not repeat (SECONDARY):**

- The organizer's final full host run (`both-stages-full-host-final.log`) lists 33 `ok` packages and no
  `FAIL`.
- The skill suite exited 0 (`full-review/host-validation-result.json`).

**Probe programs (PRIMARY, all under `.parley-runtime/claude1-r3/`):**

- `zadv/`: 54 zcode recognizer cases.
- `reinclude/`: owner re-inclusion and scope edits.
- `knoboff/`: knob-off manual exclusion.
- `create/`: idea creation on AppleVirtIOFS against local `/tmp`.
- `lease/`: two-process idea lease on AppleVirtIOFS against `/tmp`.

### Per-criterion results

**AC1, protocol text. PASS, scoped. PRIMARY.**

- The §9.0, §5, Phase 0, Phase 3/5/6/7 and §9-checklist hunks are present
  (`git diff --word-diff 27e42b8 78ac536 -- parley-deck/COOPERATION.md`).
- The drift test passes, the skill copy is byte-identical (`cmp` exit 0), and `meta/protocol-changelog.md:2`
  names the idea.
- Packets for phases 0, 5 and 8 render with `context_mode=full`, packet `73613f95…` and no fallback. Each
  contains `**Quota auto-exclusion**`, `quota_auto_exclude: false`, the Phase 5 pin cross-reference and the
  zcode exception.
- Two caveats:
  - The packets are full-context, so inclusion is trivially satisfied.
  - "Not yet in force" is absent. That is consistent only because both stages ship in one change.
- Bootstrap: the CLI embedded template keeps its generic header zones and the drift test passes. The
  skill-based bootstrap now relies on the new Appendix A instruction (NIT-4).

**AC2, positive classification. PASS for the retained excerpt only. PRIMARY.**

- `zadv` case P1 is the 4-line incident-2 stderr with exit 1. It is eligible with:
  - `rule=quota.named-reset-ge-60m.v1`;
  - excerpt `Weekly/Monthly Limit Exhausted`;
  - `reset=2026-10-04T22:14:57.003Z`;
  - raw `2026-10-04T22:14:57.003Z; 176930; 49h 8m 50s; 2026-10-05 06:14:57`.
- The invocation id and observation time are carried on the evidence.
- Limits:
  - The positive is an excerpt (`testdata/quota/README.md`), and incident 1 is a labeled paired replay.
  - The real SDK dump framing is unverified against the recognizer. That framing includes the wrapper
    line, multi-line `responseHeaders`/`data`, and the aggregation of four or more retries.

**AC3, negative classification. PARTIAL FAIL. PRIMARY via `zadv`.**

Rejected, as required:

- exit 0, a valid artifact, a later success, a watchdog or a truncated capture;
- no terminal line, or trailing text after it;
- the codex, claude and kimi adapters;
- a markdown fence, a JSON `role` envelope, or an `API Error:` quote;
- 429 with an `Error:` line, a `TypeError`, a 503 record, a status-only 500, or a dangling wrapper;
- a 30 m reset and a 59 m 59 s reset;
- display-only, a display off by 2 s, a display at a non-15-minute offset, and a display at +14:15;
- `reset_at` contradicting `retry_after`, or a header contradicting the body;
- past, zone-less, null or duplicate-key resets;
- hourly or 5-hour limits with no reset;
- a bare 429 with a long `retry_after`, a generic `quota exceeded`, a bare `credit balance`, and "limit has
  not been reached";
- two records with different allowance text.

Accepted, wrongly:

- a raw (unenveloped) quoted record after prose or a `$ cat` line;
- a 429 together with `RangeError`, `SyntaxError`, `AbortError`, `FetchError … socket hang up`,
  `AI_RetryError … Last error: Internal Server Error`, or an `Unhandled promise rejection` line.

See MAJOR-3. The codex 503s and the codex 401 are diagnostic-only by adapter. The zcode path rejects any
status other than 429 (`quota_zcode.go:48`).

**AC4, unsupported adapters. PASS. PRIMARY.**

- `ClassifyQuota` short-circuits every non-zcode adapter (`telemetry/quota.go:68-74`).
- The support table lists codex, kimi and claude/text as diagnostic-only (`:18-25`).
- I executed this for codex, claude and kimi.

**AC5, the floor. PASS. PRIMARY.**

- Reading `quota/quota.go:118-192`:
  - dedup goes through `FilterConfirmed`;
  - a member that is unresolved, or that has mixed observations, is not usable;
  - the facilitator is excluded from `UsableSurvivors`;
  - the candidates are sorted;
  - all-or-nothing applies `Block`.
- `internal/quota` tests (permutations, 4→2, 3→1) pass in my run. A two-non-facilitator idea can reach at
  most 1 survivor, which is under 2, so it can never apply.

**AC6, kickoff (C1). PASS on local disk only; BLOCKED on this deck's volume by CRITICAL-1. PRIMARY.**

- `runcontrol.go:56` filters before `CreateIdeaWithQuota`, `run.created` and the manifest are written.
- The runcontrol and app quota tests pass. On AppleVirtIOFS, creation fails before any of these are
  written.

**AC7, bare 503. PASS. PRIMARY.**

- `(^|[^0-9])503([^0-9]|$)` is added to both the runner (`failclass.go:31-32`) and the preflight classifier.
- `TestQuotaBare503RunnerGate` and `TestQuotaPreflightReportsOnlyAndBare503` pass in my run. The pattern
  is broad (NIT-3).

**AC8, protected roles. PASS. PRIMARY.**

- Before the floor is checked, `Evaluate` blocks on the designee, the pin, the facilitator or a started
  drafter (`quota/quota.go:166-175`).
- `membership.Roles` reads the designation, the `IMPLEMENTATION.md` pin and the drafters
  (`membership.go:318-353`).
- The global default is used only in `CheckGates`, as fall-through (`gates.go:47-55`).
- The blocking note carries the three exits and the evidence JSON (`membership.go:354-390`).
- Over-protection is noted as NIT-1.

**AC9, filed artifacts survive. PARTIAL. PRIMARY.**

- Known signers come from history and required signers from the current set (`consensus.go:125-135`).
- Unresolved retained obligations force `TriageBlocked` (`:139-147`).
- An excluded agent cannot append a signoff (`:296-302`).
- Kickoff-excluded ids are never known, because `Known` is the filtered kickoff list
  (`history.go:124`).
- Tests pass in my consensus run.
- Not met: the FINAL §7 release of a veto, "withdraws it after an owner-confirmed re-inclusion", cannot
  happen (MAJOR-1).

**AC10, gates. PASS for re-evaluation at transition time. PRIMARY.**

- `CheckGates` (`gates.go:17-112`) re-evaluates:
  - the track reviewer count;
  - the auto-implement floor of 2 (LE-7/LE-11);
  - the goal checker under auto or strict;
  - model diversity from the manifest roster snapshot.
- A shortfall becomes a `Block` (`membership.go:283-288`). `Finalize` re-checks gates on the current set
  (`consensus.go:333-352`).

**AC11, durability. PASS on local disk for batch and projection faults; FAIL for kickoff durability on this
volume (CRITICAL-1).**

- PRIMARY: `DurableWrite`, `syncDir` and `SyncPath` use `fsutil.SyncFile`.
- PRIMARY: `CommitBatch` replays are idempotent by id (`history.go:278-302`).
- PRIMARY: `InspectQuota` marks a missing receipt pending (`protocol/quota.go:78-89`).
- PRIMARY: pending blocks launch (`telemetry.go:93-95`), and `QuotaMembers` blocks signoff and close
  (`protocol/quota.go:117-119`).
- SECONDARY: the producer's fault-injection tests pass. My package run includes them.
- A truncated receipt cannot be repaired (MINOR-5).

**AC12, serialization. FAIL on this deck's volume; PASS on local `/tmp`. PRIMARY executed (MAJOR-2).**

**AC13, partial artifacts. PASS. PRIMARY.**

- `ValidateRound` checks every survivor before commit (`membership.go:172-211`, called at `:289-296`).
- Partial and failed files are untouched by transitions. That is a reading; the producer tests pass.

**AC14, history. PASS. PRIMARY.**

- `ReadHistory` enforces strict decoding, sequence, digest and revision checks, and never repairs
  (`history.go:108-156`).
- In my `reinclude` probe, a contradictory projection or policy returned an error and blocked.

**AC15, records and surfaces. PASS, with NIT-2. PRIMARY.**

- The marker carries the automatic label, rule, reset or `unknown`, raw reset, transition id and date. It
  never says "confirmed" (`record.go:128-137`).
- Notices are created with `O_EXCL` and named by transition id.
- `TestQuotaStatusWaitBriefAgreePendingAndAppliedReadOnly` passes in my run.

**AC16, configuration. FAIL (CRITICAL-2; widening part of MAJOR-1).**

- PASS, PRIMARY: presence-aware layering and fail-closed malformed values (`config/runtime_test.go`, which
  passes).
- PASS, PRIMARY: legacy ideas without a kickoff record stay on confirmation (`protocol/quota.go:54-58`).
- PASS, PRIMARY: recorded policy is reused and immutable.

**AC17, return. PASS. PRIMARY.**

- No product code writes `agents.toml`. The only additions are test fixtures in `config/runtime_test.go`.
- `RelaunchHint` is the reset plus 5 minutes, labeled a provider estimate, and empty when the reset is
  unknown (`quota/quota.go:82-87`).
- There is no timer or rejoin code.

**AC18, integrity. CONDITIONAL. PRIMARY.**

- By design the only trigger is invocation evidence. Retained capture runs after the decision
  (`membership.go:300-303`).
- Raw-quoted stderr acceptance (MAJOR-3) is a possible content trigger if zcode ever writes model or tool
  text unenveloped to stderr. That is UNVERIFIED.

**AC19, completion. CONDITIONAL, currently honored.** `IMPLEMENTATION.md` says `status: in-progress`.

**AC20, checks. PASS as executed or reported, but not representative.**

- PRIMARY: build, vet, the targeted packages, drift, `cmp` and the packets.
- SECONDARY: the full host suite and the skill suite.
- None of these suites runs on the deck's AppleVirtIOFS volume. CRITICAL-1 and MAJOR-2 are invisible to
  them.

**AC21, close. NOT MET. Conditional on the owner.** It awaits fixes, review consensus, both signoffs and the
attended owner close.

### Prior dispositions (my independent evaluation)

- **round-02 MAJOR-1 (bounded stderr channel) and MAJOR-2 (display clock).** Closed by the operator ruling
  quoted above.
  - The display-clock rule is implemented as ruled. PRIMARY: agreement within 1 s, a machine value
    required, offsets −12:00 to +14:00 in 15-minute steps, and display-only resets gating
    (`telemetry/quota.go:195-249`).
  - The stderr channel follows most clauses. Its quoted and mixed gates are incomplete, which is new
    MAJOR-3 and not a re-raise.
- **round-02 MINOR-1 (support table).** Concur resolved. `QuotaSupport` and `testdata/quota/README.md` state
  the limits.
- **round-02 MINOR-2 (parenthesis overcapture).** Concur resolved. `resetAtText` stops at `(`, and the raw
  value shows `2026-10-05 06:14:57` separately (PRIMARY).
- **Historical AC1 bootstrap discrepancy.** Concur that AC1's literal byte-identity is met while the CLI
  bootstrap zones remain. The residual risk is NIT-4.
- **No other reviewers.** codex-1 implements, and no other non-implementer reviewer is active.

## Findings

### [CRITICAL] CRITICAL-1: New-idea creation fails on this deck's AppleVirtIOFS volume, with the knob on or off

**What is wrong.**

- `quota.WriteKickoff` calls raw `f.Sync()` and `dir.Sync()` (`internal/quota/record.go:97,109`). On
  Darwin that is `F_FULLFSYNC`.
- The same raw calls appear in:
  - `protocol.CreateIdeaWithQuota` (`internal/protocol/quota.go:237,258`);
  - the kickoff blocking-note writer (`internal/app/quota.go:35`);
  - the kickoff notice writer (`internal/runcontrol/runcontrol.go:123`).
- The repository already has `fsutil.SyncFile` for exactly this case
  (`internal/fsutil/sync_darwin.go:9-25`). Its comment reads: "Some shared filesystems reject that
  device-specific operation with ENOTTY but support ordinary fsync".
- `parley run` (`internal/app/app.go:1966`) and the TUI launcher (`:2472`) always pass a non-nil policy.

**Executed (PRIMARY, `.parley-runtime/claude1-r3/create`).** On `/Volumes/My Shared Files/…`
(`mount`: AppleVirtIOFS), `CreateIdeaWithQuota` fails for both knob values with:

`sync …/.parley-runtime/quota-kickoff-…/quota-kickoff.json: inappropriate ioctl for device`

The legacy nil-policy API succeeds on the same volume. On local `/tmp`, both knob values succeed.

**Why it blocks.**

- After delivery, `parley run` cannot create any idea in this workspace. That is a deck-wide, knob-off
  regression outside C1 and the bare 503.
- Host tests use `t.TempDir()` on APFS, so they cannot see it.

**Fix.**

- Route every new sync through `fsutil.SyncFile`.
- Add a seam test that injects `ENOTTY`.
- Run one acceptance check of `parley run --no-preflight` against a scratch deck on this volume. Any
  provider dispatch can be stubbed.

### [CRITICAL] CRITICAL-2: Every new idea locks its membership into quota history, which breaks owner-confirmed manual exclusion and catch-up, even with the knob off

**What is wrong.**

- Every idea created by `parley run` or the TUI gets an immutable kickoff record, whatever its policy
  (`app.go:1966`, `runcontrol.go:57`, `protocol/quota.go:226`).
- `InspectQuota` then treats any `participants:` that does not equal a recorded revision as follows
  (`protocol/quota.go:66-77`):
  - "contradictory quota membership projection", or
  - "pending" when the list equals an earlier recorded revision.
- `QuotaMembers` errors on both (`:109-123`), so these all fail:
  - consensus `Status`, `AppendSignoff`, `Draft`, `Finalize` and `Reopen` (`consensus.go:125,296,529`);
  - the launch gate (`telemetry.go:87-95`).
- The batch record type can only shrink membership (`history.go:67-95`). Nothing can record an
  owner-confirmed change.

**Executed (PRIMARY, `knoboff`).**

- `NewPolicy(false)` creates the kickoff record.
- An owner-confirmed manual §9.0 exclusion (participants `[a, b, c]` plus an `excluded:` line) gives
  `QuotaMembers` → `contradictory quota membership projection`.
- A catch-up join of `e` gives the same error.

**Why it blocks.**

- FINAL §11 and the §10 basis allow exactly two knob-off behavior changes, C1 and the bare 503.
- R2 says ideas that do not auto-exclude "stay on confirmation". The confirmation path is now impossible
  for every new idea on every deck.

**Fix (either option).**

- Option A: write no kickoff record when the policy is off. Keep only C1 filtering, so knob-off ideas
  stay on the legacy prompt authority.
- Option B: add a recorded, owner-confirmed membership-revision type that the manual §9.0 and catch-up
  paths use, with a CLI entry point. It must cite the archived owner answer.

Add regression tests for a manual exclusion and for a catch-up join on knob-off ideas.

### [MAJOR] MAJOR-1: Owner-confirmed re-inclusion, veto withdrawal and scope widening have no path in policy-on ideas, although the protocol and skill promise them

**What is wrong.**

- FINAL makes owner-confirmed re-inclusion part of the design:
  - §5: re-inclusion "stays owner-confirmed, under the existing catch-up rules";
  - §7: a ❌ "stands until that agent withdraws it after an owner-confirmed re-inclusion, or until the
    owner rules";
  - §9: an agent cannot sign "until an owner-confirmed re-inclusion puts it back".
- FINAL §11: widening goes through the existing owner-confirmed path.
- The new §9.0 and §5 hunks, `SKILL.md` ("Re-inclusion is owner-confirmed") and `ROSTER_AND_PROTOCOL.md`
  repeat these promises.
- Nothing in the code implements them:
  - no command or record type for any of the three;
  - batches can only shrink membership (`history.go:72-94`);
  - a retained veto can be disposed only by `operator-ruling` (`quota_retained.go:148`), and never by its
    own author (`:171`). So author withdrawal after re-inclusion can never release it.

**Executed (PRIMARY, `reinclude`).**

- Restoring `participants: [a, b, c, d]` after `d` was excluded gives `pending="quota membership projection
  pending"`, and `QuotaMembers` errors.
- `ReconcileQuotaPrompt`, which the next driving `Before` calls, rewrites the line back to `[a, b, c]`. The
  owner's re-inclusion is silently reverted.
- Adding a new agent, or editing the scope or policy line, gives "contradictory".

**Fix.**

- Add an owner-confirmed revision record. It cites the owner-answer path and may only add ids that are
  already in `Known`.
- Give the reconcile step a rule for it.
- Let a re-included author's later recorded withdrawal dispose of its retained veto.
- Either define the widening path, or state in FINAL-compatible protocol text that none exists. Today
  `NewPolicy` never produces `kickoff-only`.

### [MAJOR] MAJOR-2: The idea lease and projection lock are not exclusive on this deck's volume, so the AC12 serialization silently does nothing

**What is wrong.**

- `tryLock` is `unix.Flock(LOCK_EX|LOCK_NB)` (`internal/membership/lock_unix.go:11-17`) on
  `<idea>/quota-driver.lock` (`lock.go:35-53`).
- On AppleVirtIOFS, flock is not exclusive across processes.

**Executed (PRIMARY, `lease` probe).** I built a mid-idea-scoped idea and used the real
`membership.Acquire`:

- On the shared volume, run-A held the lease and run-B's `Acquire` also returned `<nil>`.
- The `/tmp` control rejected run-B: `idea driving lock held: demo (competing run run-B)`.
- A Python `fcntl.flock` probe gave the same result.

**Why it matters.**

- The FINAL §10 guarantee fails without any error on the owner's deck: a competing run is neither rejected
  nor serialized. FINAL's own list of what would prove the position wrong includes "a consumer that
  dispatches from a stale list".
- FINAL §10 named the driver's PID primitive. The flock deviation is not portable to this volume.

**Fix.**

- Use an `O_EXCL` lease file carrying the run id and PID, with liveness checking and stale takeover. Hard
  links and `O_EXCL` work here (PRIMARY: `ln` succeeded on the volume).
- Alternatively, self-test exclusivity at acquire time and fail closed.
- Test on this volume.

### [MAJOR] MAJOR-3: The zcode recognizer accepts raw-quoted records and most mixed other-error stderr forms the owner rule requires to gate

**What is wrong.**

- `classifyZcodeQuota` (`internal/telemetry/quota_zcode.go:24-158`) skips any unrecognized line unless it
  contains `responseBody`, `statusCode`, `retry-after`, `"role"` or `API Error:`, or starts with `Error:`,
  `TypeError:`, `[tool`, `[assistant` or a code fence (`:149`).
- **Raw quotes are accepted.** A record reproduced after free prose ("The earlier probe said:" or
  `$ cat zcode-incident-2.stderr`) is eligible (PRIMARY, `zadv` cases Q-raw).
- **Mixed errors are accepted.** A valid 429 followed by any of these is eligible:
  - `RangeError:`, `SyntaxError:`, `AbortError:`;
  - `FetchError: … socket hang up`;
  - `AI_RetryError: … Last error: Internal Server Error`;
  - `Unhandled promise rejection: ZodError` (PRIMARY, `zadv` cases M-*).
- The owner rule requires that every provider error record agrees, that a "mixed … record set gates", and
  that "a quoted 429 in assistant or tool text" be an adversarial fixture.
- The README says zcode's stderr sink is untagged. The fixtures cover only JSON-enveloped quotes. Whether
  zcode can emit model or tool text raw on stderr is UNVERIFIED either way.
- The accepted residual risk is a subagent attribution, not a quotation or a different provider failure.

**Fix.**

- Gate on any JS error-class header (`[\w$]*(Error|Exception):`) other than the accepted
  `AI_APICallError` wrapper.
- Accept `AI_RetryError` only when its "Last error" equals the record message.
- Require each record to sit inside the SDK dump framing, and reject free prose between framing and the
  terminal line.
- Add fixtures for each case above.
- Obtain one complete native capture so the positive is tested against real framing.

### [MAJOR] MAJOR-4: Manual-mode signoff handoffs drop the idea lease and always fail on mid-idea ideas

**What is wrong.**

- `writeSignoffHandoff` (`internal/app/consensus_request_signoffs.go:613-626`) does not pass `Context`.
- `runner.WriteHandoffPacket` falls back to `context.Background()` (`internal/runner/handoff.go:54-58`).
- `beginLaunch` then requires the lease for MidIdea ideas (`internal/runner/telemetry.go:100-103`), and
  that check fails with "quota mutation requires idea driving lease".

**Provenance.** PRIMARY reading. Not executed.

**Impact.** Interactive or manual signoff requests fail on every default-on new idea.

**Fix.** Thread the leased `ctx` into `HandoffOptions.Context`, and add a handoff test under a mid-idea
history.

### [MINOR] MINOR-1: The projection lock fails instead of waiting

`ProjectionLock` uses `LOCK_NB` (`internal/membership/lock.go:96-102`). Signers that append at the same
moment as each other, or during a commit, get "quota projection lock held" instead of waiting a bounded
time. Suggest a bounded wait with retry.

### [MINOR] MINOR-2: The disposition gate checks structure only

`quotaDisposition` (`internal/protocol/quota_retained.go:127-188`) accepts:

- a self-declared frontmatter author;
- the string `Authority: owner`;
- any `## User direction` heading;
- any `Evidence:` or `PRIMARY` token anywhere in the file.

The producer discloses this. Suggest binding owner rulings to an archived `inbox/user-to-*` answer path.

### [MINOR] MINOR-3: Launches without an idea skip the quota gate

`beginLaunch` gates only when `info.Idea != ""` and the phase is not preflight (`telemetry.go:87`). So
`parley agents exec` without `--idea` can launch an excluded id, and `RequireStopped` does not see it.

### [MINOR] MINOR-4: Lock files land in the tracked idea directory

`quota-driver.lock` and `quota-projection.lock` are created in the canonical idea directory and are not
git-ignored (PRIMARY: `git check-ignore` printed nothing).

### [MINOR] MINOR-5: A truncated receipt has no repair path

`reconcileLocked` treats an existing `quota-applied/<id>` as applied (`membership.go:138-146`), while
`InspectQuota` rejects a receipt whose content differs (`protocol/quota.go:87-88`). A truncated receipt
therefore blocks permanently. This fails closed, but it needs a documented owner recovery.

### [NIT] NIT-1: Completed drafts still protect their drafters

`membership.Roles` protects the drafter of any existing `consensus.md`, `FINAL.md` or `review/consensus.md`,
including closed ones (`membership.go:332-351`). This fails closed, but it reduces the benefit.

### [NIT] NIT-2: The notice names no specific gates

The notice says "Existing review, diversity and close gates remain in force" (`record.go:147`). FINAL §9
asks for the remaining gates by name.

### [NIT] NIT-3: The bare-503 pattern is broad

The bare `503` token also matches line numbers and byte counts (`failclass.go:31-32`).

### [NIT] NIT-4: The skill reference carries this deck's header

The skill `references/COOPERATION.md` now carries this deck's header and host table. Skill-based bootstraps
depend on the new Appendix A replacement instruction.

## Open questions

1. Can zcode 3.7.7 write assistant or tool text unenveloped to stderr in `--prompt` headless mode?
   - The answer sets the severity of the raw-quote half of MAJOR-3.
   - Two leads came from reading `vendor/zcode.cjs` (sha `3e3433d9…`): the AI SDK default
     `onError: console.error(error)`, and a WASI `fd_write` path to stderr. Both are UNVERIFIED.
2. Does the owner want CRITICAL-2 fixed by option A (no kickoff record when the policy is off) or by option B
   (a recorded owner-confirmed revision type)? Option B also resolves MAJOR-1.
3. Two unverified leads from my read-only helper (SECONDARY):
   - The TUI path may emit a spurious integrity-block note after its lease is released (`app.go:1990`).
   - A pipeline holds the lease only per call (`pipeline_cmd.go:791-857`).

   Please confirm or rebut both in Phase 7.
4. Not executed by me:
   - the full host `go test ./...` (the producer reports 33 `ok`);
   - the full skill suite (the producer reports exit 0);
   - Windows runtime (compile-only, reported by the producer);
   - fault-injection internals beyond the passing package tests.
