---
idea: meta-protocol-change-quota-auto-exclude
review-cycle: 5
outstanding_agreed_fixes: 4
blocked: false
drafted-by: codex-1
date: 2026-10-08
reviewed-commit: e04852ef36846b8ec18f581c4e8ae4f28b0c2262
skill-commit: b9596ddd5c37633026a26031f1bccaee76d2622f
---

## Scope and review basis

Full independent round-07 reports 0 CRITICAL, 0 MAJOR, 2 MINOR and 2 NIT.
Cycle 4 and G18/G18b/G19/G20 are independently verified there. This is the last
narrow fix-up cycle authorized by finish-now. Any findings remaining after its
full re-review require a blocking owner note; there is no cycle 6 authorization.
Role concentration (§15.5): codex-1 organizes, implements and drafts; claude-1
independently signs this plan and reviews the result in a separate process.
No code verdict, final signoff or close is inferred from this plan.

## Agreed fixes

- **G21 — explicit Windows limitation (R7-MINOR-1).** Choose the reviewer's
  disclosure-only alternative. CHANGELOG and release notes must plainly state
  that Windows CLI 1.51.0 cannot create new ideas with `parley run`: kickoff
  directory sync fails with "Access is denied", including policy-off launches.
  Windows assets remain experimental; CLI winget stays held. Route the runtime
  fix to the existing `windows-portability` track. Do not expand this cycle into
  Windows durability work or relabel the failed CI as untested/passing.
- **G22 — refuse aliased deck paths before lease writes (R7-MINOR-2).** Derive
  the deck from the lexical `parley-deck/ideas/<slug>` layout before resolving
  filesystem aliases. Explicitly reject a symlinked deck or idea-scope directory
  with a clear unsupported-alias diagnostic; never search resolved ancestors
  for an unrelated `parley-deck`. Resolve ordinary ancestors above the workspace
  normally (including macOS `/tmp`) and preserve one physical lease identity
  through supported ancestor aliases. Validate before Acquire can publish an
  integrity notice and before driver/projection lease creation. Documentation
  names symlinked decks/idea scopes unsupported. Test local/shared paths, a
  misleading outer `parley-deck` ancestor, direct deck/ideas/idea aliases,
  ordinary paths and ancestor aliases. Assert no runtime or inbox writes on
  refusal, same physical lease/conflict for supported aliases, and unchanged
  nested-context/off-scope behavior. No new schema or broader symlink support.
- **G23 — bind-time wording (R7-NIT-1).** Qualify the first contradictory
  live/archived-owner-answer rejection sentence with "at binding"; preserve the
  committed-object-only historical validation description. Add one concise
  organizer guidance sentence: do not roll status back to round-01 to bypass
  catch-up. This answers round-07 open question 1, without a new runtime gate.
- **G24 — skip unnecessary manual clarification (R7-NIT-2).** If publication
  inspection fails and no historical misleading manual notice was observed,
  end the correction branch with a single non-blocking diagnostic for that
  failed inspection. Do not inspect/publish the correction or create a
  `-manual-authority` receipt. Preserve actual historical clarification and
  genuine correction-receipt integrity checks. Meaningful tests cover unsafe
  inbox/current manual revisions, genuine historical labels, repeat recovery,
  and malformed/symlinked correction receipts; membership remains unaffected
  by publication-only failures. No new retry/queue/delivery mechanism.

## Deferred follow-ups

- Native AC2 / R5-MAJOR-2: NOT MET, expressly owner-waived for this release;
  accepted/deferred to `../quota-zcode-native-exhaustion-capture/00-prompt.md`.
  No capture or grammar relaxation now; no claim native zcode auto-exclusion
  works. Unrecognized failures retain the owner-confirmed path.
- Windows runtime fix: `../windows-portability/00-prompt.md`; G21 closes only
  the disclosure finding, with the broken runtime behavior retained explicitly.
- D6 legacy driver accounting and unsupported adapter provenance remain separate
  follow-ups (TBD where not opened); container/PID namespace behavior unverified.

## Dismissed findings

None. All four round-07 findings are addressed above, pending independent review.

## Coverage & blind spots

Round-07 covered the full implementation since FINAL, all fix-ups and skill diff,
with independent macOS full suite/build/vet/race/shared/local/skill execution.
Windows failure evidence is hosted CI, not local Windows execution. Re-review
must be full scope, with current-tree AC1–AC21 evidence, all dispositions freely
weighed, and required full Go/build/vet/gofmt/skill/packet checks on frozen source.
Current review heading is invalid to the literal validator; claude-1 must correct
its own `## Refutation attempts (AC1–AC21)` to `## Refutation attempts` before
transition. Original review/hash are preserved in source-context/cycle5-plan-20261008.
Codex does not edit the reviewer artifact. This formatting repair changes no verdict.

## Drafter position changes

The Windows limit is known failure of new-idea creation, beyond the prior general
experimental/unverified disclosure. Aliased deck support is explicitly refused,
consistent with the baseline evidence-scope limit, instead of accidentally using
an unrelated ancestor for storage. G18/G18b and smaller G19 remain unchanged.
There is no scope/FINAL change, waiver beyond the owner's AC2 ruling, or new gate.
The newest owner note supersedes the earlier cycle-4-only cap and close-request rule.

## User direction

> ## Owner direction: finish now (supersedes the wait in `…_long-quota-answer.md` and `…_provider-stop-answer.md`)
>
> Relayed by the owner's Claude Code session on 2026-10-07 at about 23:20 CEST. Verbatim (Slovak):
>
> > "sakra tak to fixni a dokonci a deployni cez vsetky kanaly, taha sa to dlho"
>
> Translation: "Damn, then fix it, finish it and deploy it through all channels, this is dragging on."
>
> ## Relay facts (PRIMARY, 23:18 CEST)
>
> - `ANTHROPIC_BASE_URL` points every Claude CLI call at the OmniRoute gateway, so there is no direct route.
> - Two probes ran, each with a 103,660-byte prompt (FINAL + consensus + review consensus): one with
>   `--model 'claude-opus-5-5[1m]'` (plain id) and one with `--model 'claude/claude-opus-5-5[1m]'`
>   (prefixed id). **Both returned `PONG`.** Large requests pass right now. The 22:01 failure came after
>   857 s of an agentic session, so the gateway pool is intermittent, not hard down.
> - The relay's own long sessions use the plain id and kept working through the pool errors that hit the
>   prefixed id.
>
> ## What the owner authorizes now
>
> 1. **Start immediately.** Do not wait for 2026-10-09. The relay killed the auto-resume daemon.
> 2. **claude-1 stays the reviewer and the model stays Opus 5.5.** For the rest of this idea, launch claude-1
>    with the plain model id `claude-opus-5-5[1m]` instead of `claude/claude-opus-5-5[1m]`. This is the
>    same model at the same max effort; only the gateway route id changes. Record it in
>    `organizer-notes.md`. Do not edit `agents.toml`.
> 3. **Provider errors no longer stop you.**
>    - If a claude-1 or codex-1 step fails with a 429 or 503 quota/unavailable error, whatever reset it
>      states, wait 15 minutes and relaunch the same step. Do this at most 8 times per step.
>    - Silent timeouts keep the earlier rule: relaunch with 2400 s, then 3600 s.
>    - Stop only on an auth or credit error, or when the attempts for a step run out.
> 4. **Fix-up.** Do cycle 4 now: the plan signoff, the implementation and the full re-review. If the cycle-4
>    re-review finds new findings, fix them narrowly in cycle 5 without asking. That is the last cycle the
>    protocol's deliberation cap allows. Escalate only if a CRITICAL needs a change of scope or FINAL, or
>    if findings remain after cycle 5.
> 5. **Close is pre-confirmed.** The owner's "dokonci a deployni" is the attended-close confirmation,
>    provided all of the following hold:
>    - the final re-review has no open CRITICAL or MAJOR;
>    - both review-consensus signoffs exist;
>    - current-tree evidence for AC1 to AC21 is recorded, with AC2 owner-waived as already decided.
>
>    Record this verbatim as the close authority. Do not write a separate close-request note.
> 6. **Release immediately after the close,** on all channels, as `IMPL-ORGANIZER-BRIEF.md` says:
>    - merge to `main`;
>    - CLI 1.51.0, skill 2.15.0 and core 2.15.0 staged;
>    - GitHub releases;
>    - both Homebrew formulae;
>    - a winget PR for the skill only;
>    - install the skill into all local runtimes and verify each by hash;
>    - claude-1 verifies every channel independently, using a short brief.
>
>    The two steps only the owner can run (`! npm publish …` and `! parley protocol publish …`) go in ONE
>    final note, `codex-1-to-user_meta-protocol-change-quota-auto-exclude_released.md`, with the exact
>    commands.
> 7. **Unchanged:** the owner's earlier answers (the zcode stderr rule, the AC2 waiver and the follow-up),
>    no other reviewer, no quorum change, and English artifacts.

## Signoffs

### Signoff: codex-1 — 2026-10-08
Status: ✅ ACCEPT
Notes: I accept this narrow plan and its explicit limits as implementer/organizer,
not as an independent code verdict. Full phase-7 context was attested with
context_mode=full, source_sha256=packet_sha256=73613f95aabff8ed57372a535913985554014b71911cc8126ec3f7890a9d4f4e,
fallback_reason absent. Current local canonical evidence governs; no new memory
claim is used. All four dispositions remain subject to claude-1's own signoff and
full re-review. Preserve genuine history/receipt gates and no out-of-root writes.

### Signoff: claude-1 — 2026-10-08
Status: 🟡 ACCEPT-WITH-RESERVATIONS
Notes: This signs the cycle-5 plan G21–G24 only. It is not a code verdict, accepts no code and grants no
close; my full-scope round-08 re-review follows the implementation. codex-1 organizes, implements and
drafts here (§15.5), so its procedural calls stay provisional and this is not an independent parent verdict.

Context attestation: `context_mode=full`,
`source_sha256=packet_sha256=73613f95aabff8ed57372a535913985554014b71911cc8126ec3f7890a9d4f4e`,
`fallback_reason` absent (the key is not in the JSON). I ran `parley protocol packet --dir . --phase 7
--track deliberation --idea meta-protocol-change-quota-auto-exclude --flag protocol_change --json` at
2026-10-07T22:22Z with my own build of `e04852e` and with the installed parley 1.50.0. Both exit 0 with
empty stderr and byte-identical 26,812 B JSON. The body (1,501 lines, 125,862 B) reproduces the hash and
equals deck `COOPERATION.md` (`cmp`) and the skill reference copy (hash). I read the whole body in chunks
1–300, 301–600, 601–900, 901–1200 and 1201–1501, with over-long lines 59, 452 and 1002 printed in full.
The shadow packet (`01ba034e…`) was not used.

Coverage this launch:
- Read in full: this file (`b6b56e56…` before this block); the archived cycle-4 plan
  `review/round-07/consensus.md` and my raw `review/round-07/claude-1.md`, whose original hashes match
  `cycle5-plan-20261008/original-hashes.json`; the finish-now (`caafe192…`), round05-answer and
  round06-answer notes; the release draft and the platform-CI note.
- Read in part: FINAL (`f90577f1…`, unchanged) §9, §13.6–13.7 and AC1–AC21; the uncommitted
  IMPLEMENTATION and organizer-notes cycle-5 entries; the source cited below; the `windows-portability`
  branch.
- Not reread: the other inbox answers, the rest of IMPLEMENTATION, and tests other than one runcontrol
  fixture.
- Executed: two scratch probes injected with `go test -overlay`, so no tree file was written, under my
  provider-guard PATH. `guard-denied.log` is absent and `git status` is unchanged except for my heading
  repair. Nothing was committed, and no provider, participant, worktree or OpenViking call was made.

User direction: unquoted, `## User direction` is byte-identical to the finish-now inbox body (checked).
Operative sentence, verbatim: "sakra tak to fixni a dokonci a deployni cez vsetky kanaly, taha sa to dlho".
Point 4 binds: "fix them narrowly in cycle 5 without asking. … Escalate only if a CRITICAL needs a change
of scope or FINAL, or if findings remain after cycle 5." Point 5 binds the close: no open CRITICAL/MAJOR,
both review-consensus signoffs and current-tree AC1–AC21 evidence, "with AC2 owner-waived as already
decided". AC2 stays NOT MET / owner-waived, never a PASS.

Heading repair (my own file, done): `review/round-07/claude-1.md` line 154 now reads
`## Refutation attempts`. `diff` against the preserved original shows only that line. SHA-256 changes
`3fe2eb29…` → `4705ad7c…` (25,258 → 25,245 B). PRIMARY: `ValidateReviewArtifact` rejects the original
("missing a non-empty '## Refutation attempts' section") and accepts the repaired file (`<nil>`). No
verdict, finding or evidence changed.

- **G21 (Windows, disclosure only): concur.** This is my alternative (b). Reservation 1:
  - (a) Disclose the full reach. The cause is a directory `Sync`, which every new durable path uses:
    `quota.DurableWrite`/`SyncPath` (`history.go:188–267`), `WriteKickoff` (`record.go:110`), the
    `ideas/` sync (`protocol/quota.go:300`), `store.AppendDurable`/`Sync` (`events.go:44,117`) and
    pidlease `publish`/`reap` (`lease.go:108–153`; `publish` links the lease before its failing sync).
    So on Windows, an idea with the default mid-idea scope also cannot be driven or signed, and any
    manual import, owner revision or transition fails the same way. Ideas without quota records skip
    these paths (`history.go:122`, `MidIdea` gates).
  - (b) Make the routing real (PRIMARY, git). `windows-portability` exists only on its unmerged branch
    (`3526b82`, 2026-09-28, IMPLEMENTATION `status: in-progress`), and
    `../windows-portability/00-prompt.md` exists neither on this branch nor on `main`. Its FINAL is
    frozen at base `868825f`, which has no `internal/quota/record.go` or `internal/pidlease/lease.go`.
    Its FINAL.md:164–166 keeps Windows `SyncDir` "the fail-closed named-refusal emitter for any site
    not carrying a proved mechanism, including future sites added without conversion". Hand the site
    list and the CI evidence (log `b36a9e47…`) to that track explicitly, cite it by branch and slug,
    and do not imply that it will restore idea creation.
  - (c) Recommended: carry the same line in the skill CHANGELOG. The skill is the only artifact this
    release sends to Windows through winget.
- **G22 (alias refusal): concur** with refusal rather than support. Supporting aliases safely needs a new
  physical lease home: two workspaces symlinking one deck would otherwise hold separate leases and lose
  AC12 serialization (reasoning, not executed). That is beyond narrow scope. Reservation 2:
  - (a) Scope the refusal to the points where a driving, projection, revision or manual-import lease is
    derived, or a scoped integrity escalation would be written. `app.go:1992` calls `Acquire` after
    every creation, so a refusal at the top of `Acquire` would newly break legacy, policy-off and
    kickoff-only runs on symlinked decks. 1.50.0 has no deck-alias refusal on that path (`git grep
    EvalSymlinks` at `27e42b8`), and the plan's "unchanged … off-scope behavior" requires that.
  - (b) Reuse the baseline predicate `verificationRefusalScope` (`evidence_refusals.go:22–40`): resolve
    the root, then require the resolved `root/parley-deck/ideas/<slug>` to equal its join. The two alias
    rules then agree, macOS `/tmp` works, and the workspace root itself is an ordinary ancestor.
  - (c) The default scope is `KickoffAndMidIdea` (`quota/quota.go:52`), and `parley run` creates the
    idea (`app.go:1972`) before `Acquire`. Validate before the kickoff writes when the resolved scope is
    mid-idea; otherwise disclose that such an idea is created but cannot be driven.
  - (d) Name the refusal in CHANGELOG and the release notes as well as the docs, with the workaround: a
    physical deck path, or `quota_auto_exclude = false`. The drafter's "no … new gate" holds only for
    membership gates. This is a new explicit refusal, replacing an out-of-root write and a local-path
    failure.
- **G23: concur**, on one wording condition. Attach "at binding" to the working-copy clause only, because
  `ValidateAuthority` (`authority.go:91–111`) re-checks commit, blob, digest, idea, attribution and quote
  on every read. Suggested: "At binding, a contradictory live or archived copy is rejected; a changed
  digest, wrong idea/author, missing object or fabricated path is rejected at binding and on every later
  read." The skill already states this correctly (`ROSTER_AND_PROTOCOL.md:285–286`). Place the organizer
  sentence next to the round-1 rule (skill `ROSTER_AND_PROTOCOL.md:299–300`, docs `:82–87`), not in
  `COOPERATION.md`, so the packet hash stays `73613f95…`.
- **G24: concur.** It is my R7-NIT-2 suggestion. With no receipt, the branch re-inspects at every mutation
  boundary, so an unsafe inbox yields one correction diagnostic per call (two while the ordinary notice is
  unapplied). That beats silence, because a historical label can hide behind an unsafe path. Tests should
  assert the per-call counts, no receipt, that `ReadApplied` errors still gate, and that a historical
  clarification is still published once the path is safe.
- **Reservation 3: a new observation outside the round-07 findings** (PRIMARY, executed). Kickoff notices
  bypass G18. `runcontrol.Create` writes them with a plain `O_CREATE|O_EXCL` open (`runcontrol.go:118`)
  and returns the error. `parley run` then prints `run create failed` and exits 1 (`app.go:1986–1988`),
  after the kickoff record, `run.created` and the manifest already exist. My overlay probe reuses the
  `TestQuotaAutomaticKickoffRecordsAndNotice` fixture:

  | `inbox/` | Create error | Notices | Kickoff records, event logs, manifests |
  | --- | --- | --- | --- |
  | present | none | 1 | 1, 1, 1 |
  | missing | `no such file or directory` | 0 | 1, 1, 1 |
  | read-only | `permission denied` | 0 | 1, 1, 1 |

  `parley init` creates `inbox/` without a `.gitkeep`, so a fresh clone of a deck with an empty inbox has
  none. Mid-idea `DurableWrite` would create it. The behavior conflicts with round06-answer ("the
  exclusion notice never blocks") and with the signed AC15 meaning (exactly one notice on a safe
  destination). Round-03 CRITICAL-1 fixed only this writer's sync call, not its blocking.
  - Recommended for cycle 5: the same non-blocking checked publication, creating a missing `inbox/`,
    with tests.
  - Otherwise I must raise it in round 08 as a MINOR, and under finish-now point 4 that remaining
    finding would require the blocking owner note.
- **Other dispositions: concur.**
  - R5-MAJOR-2 / AC2 stays deferred, NOT MET and owner-waived. It is not a PASS.
  - I concur with the D6, adapter-provenance and container/PID deferrals, with "Dismissed: none", and
    with the coverage section.
  - The drafter position changes are accurate, except for the gate wording noted in G22(d).

None of these reservations is a blocker, and none needs a scope or FINAL change. Each fits this final
narrow cycle and can be verified in round 08.
