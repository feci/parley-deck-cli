---
agent: claude-1
idea: meta-protocol-change-lean-organizer
review-round: 3
date: 2026-09-24
reviewed-commit: 913f8ba
responding-to: [claude-1/review/round-02, kimi-1/review/round-02]
---

**Protocol context attestation (Phase-6 review round 3):**

```json
{"context_mode": "full", "source_sha256": "8ce83cde0b4f1c3019229869433df0c32451fdbaa13a186952fdf9805c3a9db7", "packet_sha256": "8ce83cde0b4f1c3019229869433df0c32451fdbaa13a186952fdf9805c3a9db7", "fallback_reason": null, "body_path": "/Volumes/My Shared Files/AI_WORKSPACE/parley-deck/worktrees/lean-organizer/.parley-runtime/protocol-packets/full-phase6-deliberation-8ce83cde0b4f1c3019229869433df0c32451fdbaa13a186952fdf9805c3a9db7.md"}
```

Attestation cross-check (PRIMARY): `shasum -a 256` over the packet body returns
`8ce83cde…a9db7` (109,928 B), and `parley protocol packet --audience facilitator` run in
my own review tree reports the same `source_sha256` for
`parley-deck/COOPERATION.md` — the attested packet IS the live authority at the tree I
reviewed.

**Provenance and isolation.** My own paired review worktrees, reused and moved to the
review commits, both clean (`git status --porcelain` empty at start and end):

- CLI: `worktrees/lean-organizer-review-claude-1` at **913f8ba** (record tip).
  `git diff --stat 998346c 913f8ba` touches `IMPLEMENTATION.md` only (+319/-14), so the
  source I read and built is the reviewed source `998346c`.
- Skill: `worktrees/lean-organizer-review-claude-1-skill` at **b06a65a**.
- Binary: my own `go build -o /tmp/claude1-r3/parley ./cmd/parley` (go1.27.1 darwin/arm64),
  sha256 `c5420ff6efc6631a772eaec003409e96b653a5cc14a1f02fbaf670e99b4e3c07`. This differs
  from the task binary `/tmp/parley-lean-organizer-fixup2/parley`
  (`cefd8c0b…159e`) by hash and by 16,512 B of size; both facts are explained by DF-4's
  embedded build path (my tree path is 85 characters longer), and the source identity is
  established by the `git diff` above rather than by hash equality. I did not rely on the
  task binary for any verdict.
- Inputs read in full: the live packet above; frozen `FINAL.md`
  (`git diff 120a9bf 913f8ba -- FINAL.md` empty — re-verified); `IMPLEMENTATION.md` at
  `913f8ba` including all of fix-up cycle 2; both complete round-02 review files; the
  entire signed `review/consensus.md` (all three signoff blocks); `review/consensus-cycle-01.md`;
  `00-prompt.md`; `organizer-notes.md`; the core-publish and wait-observation-02 inbox notes.
- No source edit, commit, release, tag, install or publish. `parley protocol publish` was
  not run; `~/.parley/protocol/core/` still holds only `2.10.0`. No FINAL, peer-review,
  consensus, organizer or inbox artifact was modified. This file is my only output.
  Evidence preserved under `/tmp/claude1-r3/`, `/tmp/c1r3-g1/`, `/tmp/c1r3-g2/`,
  `/tmp/c1r3-g10/`, `/tmp/claude1-r3-test-full.log`, `/tmp/claude1-r3-skill.log`.

## Summary

All ten signed fixes G1–G10 are present and, on my own independent measurement, correct;
the two frozen-FINAL acceptance elements I named in round 02 as unmet — C.3/R-2's
floor-recording half and cross-cutting "both suites green" at the command FINAL names —
are **both now met**, so I withdraw R2-MAJ-1, R2-MAJ-2 and R2-MAJ-3 as fixed and my
round-02 NOT-READY rationale no longer holds on its own terms. I reproduced the G1 floor
of **52,295 B** by two methods that do not use the implementer's test code, and the
omission totals **30,262 B / 23,230 B / 19,736 B** exactly. This round I file **no
CRITICAL and no MAJOR**: three MINOR and four NIT, of which one is behavioral (G4(a) does
not engage when file mtimes are equal, which is the state of every fresh clone or
checkout), one is a missing Phase-8-required section, one is an unstated close condition
that `auto_implement: true` imposes on the next consensus, and four are record-accuracy
items in the same class this idea has been disciplined about throughout.

## Refutation attempts

Per FINAL criterion, what I tried to break and the result. Every check below is PRIMARY —
run by me this session in my own trees, commands and output quoted.

**C.3 — named sets, floor, guardrail (R-2 / open item 14). ATTACKED HARDEST; HOLDS.**
I refused to read the floor out of the implementer's test. I rendered the packet myself —
`parley protocol packet --phase 1 --track deliberation --transport github-pr --audience
facilitator --json` — and computed the retained total two independent ways:

- *Method A, from the rendered artifact:* body bytes minus the preamble (up to and
  including the first `\n\n---\n\n`) minus the omission-index tail (from
  `---\n\n## Packet omission index\n\n` to EOF) = **52,295 B**.
- *Method B, from the source:* for every `included` record in the attestation index, take
  its `start_line..end_line` slice of `parley-deck/COOPERATION.md` and sum
  `len(rstrip(text,"\n")) + 2` = **52,295 B**.

Both agree with the recorded floor. The `+2` convention is not an assertion I accepted: I
read `renderPacket` (`internal/protocolpacket/packet.go:471`), which writes
`strings.TrimRight(r.Text,"\n")` followed by `"\n\n"` for each included block — the
convention mirrors the renderer exactly. Body − floor = 6,911 B = preamble (762 B) +
index table (6,234 B), which closes the arithmetic.

I could not break the omission companion either. My first pass produced 30,292 B (+30 over
the record) and 23,258 B (+28); reading `Parse` (`packet.go:175-226`) showed `Text` is
`strings.Join(lines[start:end], "\n")` with **no** trailing newline, so `len(Text)+1` is
exactly the block's source bytes and my reconstruction had double-counted one newline per
block. Recomputed under the correct convention: named omission set **30,262 B**, seven
top-level sections **23,230 B**, `11.A`+`11.C` **7,032 B** — all exact. I also verified the
index tiles the source exactly once in order (block 1 starts at line 1, last ends at 1387 =
the file's line count), so nothing is double-counted or skipped.

I attacked the *criterion*, not just the number. FINAL's acceptance row names "the measured
floor (**≈ 42 KB with §2**)" and the record now logs 52,295 B — a +10.2 KB / +24 %
divergence, so I checked whether the record is describing the same quantity. FINAL's own
figures reconstruct the old derivation exactly: 65,516 (optimize, phase 2) − 27,420
(omission set) = 38,096 B = 38.1 KB, matching FINAL's "≈ 38.1 KB at phase 2"; + §2
(≈ 4,004 B) = 42,100 B = **42.1 KB**, matching FINAL's "≈ 42.1 KB with §2 retained"
(decimal KB throughout). That derivation subtracts a whole-section quantity from an
**already within-section-optimized** body, so it cannot be the whole-block retained total;
the implementer's account of where ≈ 42 KB came from is arithmetically confirmed, and
52,295 B is the conceptually correct realization of "the measured whole-section floor"
R-1 framed the guardrail against (headroom to 70,000 B = 17,705 B — still "conservative
headroom above"). **R-2's operative duty is discharged**: body, floor, omission total,
guardrail and `--optimize` baseline are emitted by one `t.Logf` in one run, and the floor
is now load-bearing — `if len(c.Body) < floor { t.Errorf }` did not exist before. I
reproduced the log line verbatim from my own run:

    R-2 measurement (phase 1 / deliberation / github-pr): facilitator body=59206 B; floor (retained whole-block source total: sum over the request's included blocks of len(text)+2, the renderer's exact layout)=52295 B; named-omission-set whole-section bytes (each block len(text)+1, subsections included)=30262 B; guardrail=70000 B; --optimize baseline=65750 B

I tried to break the body figure: my own render measured **59,291 B**, not 59,206 B. The
85 B delta is exactly the embedded source-path length difference (my absolute path is 117
chars, the test's `../../parley-deck/COOPERATION.md` is 32). Not a defect, and consistent
with the record's own test-path/live-path distinction in DF-3. Floor is path-independent;
`body ≥ floor` holds at both paths.

**C.3 gating half — named-omission-set absence.** Checked against the real rendered body,
not the test: all nine named omission headings **absent**, all seven retention headings
(Quickstart, §4, §5, §9, §2, §15, active `### 11.B`) **present**, `## Packet omission
index` present, and the never-cut §15.1/15.2/15.3/15.4/15.7 blocks all present. `parley
protocol packet check` → `ok`, exit 0, 69 blocks.

**B.3 — `--json` exit codes an organizer can branch on (G2). ATTACKED; HOLDS ON ALL FOUR
PATHS.** I ran every path with streams separated and decoded stdout through `json.load`:

| path | invocation | exit | stdout | stderr |
|---|---|---|---|---|
| 0 | live deck, `--for review --timeout 3s --json` (codex-1's observation-02 shape) | 0 | 3,398 B, decodes, keys `['notes','digest']` | 47 B `wait: boundary reached (review round complete)` |
| 3 | fixture, `--for consensus --timeout 3s --json` | 3 | 1,429 B, decodes, keys `['digest']` | 60 B `wait: timeout after 3s; outstanding: consensus.md not filed` |
| 4 | fixture, escalation touched mid-wait | 4 | 1,333 B, decodes, keys `['notes','digest']` | 108 B `wait: blocking escalation …` |
| 1 | three shapes: unknown `--idea`; `--timeout 99h`; non-deck `--dir` | 1 | **0 B** in all three | error on stderr |

`notes` is correctly omitted when empty (`omitempty`), which is why the exit-3 envelope has
one key. The skill core documents exactly this contract (`SKILL.md:139-144`). Nothing in
`internal/driver/` or `internal/runner/` consumes `wait`'s stdout, so the stream move has no
internal consumer to break.

**Cross-cutting — both suites green at the command FINAL names (G3). ATTACKED; HOLDS, and
G3 is load-bearing today.** My own full run in my own tree:
`go build ./... && go test ./... -count=1 -timeout 2400s` → **exit 0, 31/31 packages ok, 0
FAIL** (10 m 50 s wall). Crucially, `internal/trajectory` took **644.652 s** on this run —
**above Go's 600 s per-package default** — and `internal/app` 576.755 s. Without `-timeout`
this suite would fail on my machine right now, which is stronger evidence than my round-02
589.7 s measurement. `.github/workflows/tests.yml:35` is `go test ./... -count=1 -timeout
45m` on the tri-platform matrix, the motivation is in the step comment, and no `go test` in
`.github/workflows/` lacks a `-timeout`. Skill: `npm test` → **exit 0, 399 passing
assertions, 0 fail**, 54 python tests, all addon manifests `ok`. Skill core **17,802 B ≤
20,000 B**.

**C — brief and audience packet.** `parley organizer brief` = **2,631 B ≤ 8,192 B**,
byte-identical across two runs (`cmp` silent), and wrote no file (`.parley-runtime` file
count 12 → 12).

**A — preflight fail-closed (G10). ATTACKED; HOLDS, no usability regression.** kimi-1's
fixture shape (agents.toml + `meta/version.json` + an idea declaring `facilitator: claude-1`
inside `participants:`), no `COOPERATION.md`: exit **1**, stdout **empty**, stderr
`preflight failed: cannot read workspace status — the facilitator-declaration gate could not
run (is parley-deck/COOPERATION.md present?): open … no such file or directory`. Dropping
`COOPERATION.md` in flips it to exit **3** with `[facilitator-declaration]` firing — same
deck, same conflict, the only variable being whether the status read succeeds. I then tried
to break it as a regression: a directory with no `parley-deck/` at all still gives the
pre-existing `no parley-deck workspace found; run \`parley init\` first`, and `parley init`
followed by `parley preflight` gives exit 3, not 1. The new hard failure fires only on a
deck that has `parley-deck/` but is genuinely unreadable.

**G5 — unevaluable `to-user` notes.** Both shapes annotated in human output *and* in the
`--json` `notes` list. I probed branch reachability: malformed-YAML, binary and
no-frontmatter files all land in the "no `idea:`" branch (because `ReadFrontmatter` only
errors on open/scanner failure, never on parse failure), and a `chmod 000` note does reach
the "frontmatter unreadable/unparseable" branch — so both branches are live and **every**
unevaluable note I could construct was annotated rather than dropped. A note with
unterminated frontmatter still parses and is treated as a real qualifying escalation;
that is pre-existing and fail-safe.

**G6 — core-publish note.** Note now reads 109,772 B, sha256
`fc907e59…62c9f`, composition verified 2026-09-24, with the verify-before-publish `shasum`
line. Live staged file: `shasum -a 256 ~/.parley/staging/COOPERATION-2.13.0.md` →
`fc907e5914a072d1a6afe249fc39401e1f8761cc1d67f2ce002dfde210762c9f`, `wc -c` → **109,772**,
mtime `2026-09-24 03:35:27`. Exact match on every fact.

**Staged core and the three protocol copies (unchanged-recheck).** Staged core vs skill
copy: `cmp` **silent — byte-identical**. CLI deck vs staged core: 13 differing lines, all
inside the project-specific allowlist (header `Workspace`/`Transport`/`Created`/`Protocol
synced`, six §2 roster rows). §15 region under the newly-stated convention `sed -n '/^## 15\.
Verification integrity/,$p' | wc -c`: **8,056 B in all three copies**, all three sharing one
region sha256 `29ca8a1eb554733a…`. G9's figure and its extraction boundary both check out.

**G7 / G8 — record and comment.** The A.4 gate-side-parity deviation bullet is present under
`## Deviations from FINAL.md`. The G8 comment sits above `arrivedBlockingEscalation`
(`internal/app/wait.go:373-384`), records the mtime approximation and its fail-loud
direction, and cites `internal/driver/loop.go:353-354`; I read those exact lines and the
citation is accurate (`name := fmt.Sprintf("claude-to-user_%s_%s.md", …)` then an
unconditional `os.WriteFile`).

**G4 — where refutation succeeded.** See R3-MIN-1. (b) and (c) hold: on a fixture with
round-01 complete 2/2 and no `consensus.md`, `next` is now `await consensus signoff` and
`outstanding:` names `consensus.md not filed`. (a) is where the attack landed.

**Named regression tests.** All seven exist and pass at my HEAD:
`TestLiveDeckFacilitatorPacketNamedSetsAndGuardrail`,
`TestWaitJSONStdoutCarriesOnlyTheEnvelopeOnAllExitPaths`,
`TestPhaseDigestNextActionFixUpPublishedAwaitsReview`,
`TestPhaseDigestNextActionRoundsCompleteNoConsensusAwaitsConsensus`,
`TestWaitOutstandingNamesConsensusNotFiled`,
`TestWaitUnevaluableToUserNoteIsAnnotatedNotSilent`,
`TestPreflightFailsClosedWhenWorkspaceStatusUnreadable`.

## Findings

### [MINOR] G4(a) does not engage when mtimes are equal — every fresh clone or checkout still reads `await implementation`

`implementationNewerThanLatestReview` (`internal/driver/phasedigest.go:187-205`) ends in
`return implInfo.ModTime().After(newest)` — a strict comparison. A `git clone`, `git
checkout`, worktree creation, archive extraction or any `rsync` without `-t` gives
`IMPLEMENTATION.md` and the review artifacts **the same** mtime, so the predicate is false
and `nextAction` falls through to `NextAwaitImplementation` — the exact wrong output
R2-MIN-1(a) named.

Reproduced in my own review tree at `913f8ba`, which has review round-02 complete 2/2,
`implementation: present=true status=fix-up-cycle-2 ready=true`, and no `round-03/`
directory — i.e. precisely the fix-up-published state G4(a) targets:

    $ stat -f '%Sm' IMPLEMENTATION.md review/round-02/claude-1.md
    2026-09-24 06:04:18      2026-09-24 06:04:18          # equal: fresh checkout
    next  = "await implementation"                         # WRONG
    $ touch IMPLEMENTATION.md                              # mtime strictly newer
    next  = "await review artifact"                        # correct

Both organizer-facing surfaces are affected: the `--json` digest and `parley organizer
brief`, which prints `## Next action (fixed enumeration)` / `await implementation`. The
live deck avoids it only because its empty `review/round-03/` directory triggers the older
directory-based path — the very path K2-F4 established was already working.

The regression test cannot catch this: `TestPhaseDigestNextActionFixUpPublishedAwaitsReview`
explicitly pins the ordering with `os.Chtimes` ("Pin the arrival order: the fix-up publish is
NEWER than every review artifact"), so the equal-mtime case is never exercised.

Why it matters beyond the wrong string: G8 was raised and accepted precisely because mtime
is an approximation, and its resolution turned on the failure being **loud** — "a false
'arrived' costs one premature exit 4 … never a missed escalation". G4(a) introduces a second
mtime dependence whose failure direction is the opposite: **silent**, and wrong in the
organizer's branch surface. That asymmetry is not recorded anywhere.

The `## Fix-up progress` record (`IMPLEMENTATION.md:503-508`) states the residual is
"retired" and that the state "now reads `await review artifact` even with no next-round
directory" — unconditional wording for behavior that is conditional on mtime ordering. It
does name the mtime keying in the same sentence, so this is imprecision, not concealment.

**Suggested fix.** Prefer a content signal over mtime: `d.Implementation.Status` is
`fix-up-cycle-N` and `d.Review.Label` is `round-0M`; when a fix-up cycle is published and
`M <= N`, a further review round is awaited regardless of timestamps (an incomplete next
round is already caught earlier by `Completed < Total`). If mtime is kept, use
`!implInfo.ModTime().Before(newest)` and add an equal-mtime test case. Either way, soften
the `:503-508` claim to state the condition.

### [MINOR] Fix-up cycle 2 has no `### Deviations from agreed fixes` section, and there is something that belongs in it

Phase 8 of the packet specifies the fix-up section shape as `## Fix-up cycle N` → `### Fixes
applied` → `### Deviations from agreed fixes`, with the explicit note that "`None` is a valid
answer" — i.e. the subsection is required, not optional. `## Fix-up cycle 1` has one
(`IMPLEMENTATION.md:494`). `## Fix-up cycle 2` (lines 14-289) has none; `grep -i deviat` over
that span returns only the G7 bullet's reference to the separate `## Deviations from FINAL.md`
section. Cycle 2 dropped a section cycle 1 carried, in the idea whose whole subject is this
protocol.

It is not an empty formality here. The signed G1 plan text says "compute the floor as the
**retained whole-section total** (**≈ 42 KB with §2**, recomputed at the fix-up-2 HEAD)", and
FINAL's acceptance row says "the measured floor (**≈ 42 KB with §2**)". The delivered figure
is **52,295 B** — +10.2 KB, +24 % against both. I hold (see Refutation attempts) that this is
the *correct* number and that R-2's duty is discharged; but a signed fix landing a value a
quarter away from the one its own plan and the frozen acceptance row name is exactly the
"turned out to … require a different approach, with rationale" case the subsection exists for.

This is disclosed, not hidden: the G1 bullet carries a "Recorded-number note for round 3"
and the `### Cycle-2 closure conditions check` residuals close with "The floor's recorded
figure is now the retained whole-block total with its convention stated (52,295 B) rather
than FINAL's round-1 ≈ 42 KB derivation". The finding is about the required section's
absence and the fact that the divergence is recorded in two narrative places rather than the
one place a decision-maker looks.

**Suggested fix.** Add `### Deviations from agreed fixes` to the cycle-2 section with one
bullet for the G1 figure (plan said ≈ 42 KB, delivered 52,295 B, why the old derivation
under-stated it, why FINAL stays frozen) and `None` for the other nine.

### [MINOR] `auto_implement: true` adds two close conditions that no round-02 artifact or the signed consensus accounts for

`00-prompt.md:8` sets **`auto_implement: true`** (and `strict_gate` is absent, as kimi-1
correctly noted). kimi-1's round-02 close rule — "default close rule — `strict_gate` is
absent from `00-prompt.md`" — is right about `strict_gate` but incomplete for this deck,
and the signed consensus's "Cycle mechanics" paragraph carries the same reading ("Phase 8
closes the idea only on a Phase-7 consensus listing **zero** Agreed fixes").

Under `auto_implement`, LE-7/LE-11 make zero Agreed fixes necessary but **not sufficient**.
This is enforced in code, not only in prose (`internal/driver/impl.go`):

- `:283-286` — `if d.cfg.AutoImplement { if rs.Summary.Triage == consensus.TriageReserved {
  return ActionEscalated, …"review consensus is ACCEPT-WITH-RESERVATIONS; under
  auto_implement, reservations need human review before completion (LE-11)" } }`
- `:287-289` — fewer than `MinReviewers` independent reviewers escalates (we have two,
  so this one is satisfied).
- `:293-297` — `if d.cfg.AutoImplement || d.cfg.StrictGate { if ok, detail :=
  d.cfg.Impl.GoalCheck(ctx); !ok { return ActionEscalated, …"goal-done gate …(LE-7)" } }`

Both independent reviewers signed cycle 2 **🟡 ACCEPT-WITH-RESERVATIONS**. If the cycle-3
consensus lists zero Agreed fixes but any signoff is again 🟡, the close **escalates to a
human instead of completing**. Separately, the LE-7 goal-done check — a fresh
non-implementer verifying FINAL's observable acceptance criteria — appears nowhere in this
idea's artifacts except as a role FINAL A forbids the facilitator from taking; I find no
record of one having been commissioned or run.

I am not asserting a code defect: the gates are correct and are this idea's own safety
design. The finding is that the closing consensus is currently being planned against an
incomplete statement of what closes it, which is the kind of surprise that surfaces at the
worst moment.

**Suggested fix.** Record the real close conditions in the cycle-3 consensus: zero Agreed
fixes **plus** either ✅ signoffs or an explicit recorded operator ruling accepting a 🟡
close, **plus** a goal-done check by a non-implementer (or a recorded operator decision that
this manually-driven close does not route through the driver). Naming it now costs a
paragraph; discovering it at close costs a cycle.

### [NIT] The 23,230 ↔ 23,223 reconciliation does not reconcile arithmetically

`IMPLEMENTATION.md` G1 says the seven top-level sections "subtotal 23,230 B — claude-1's
independently measured 23,223 B **plus this convention's per-block newline**". Those seven
sections contain **28** blocks (§1:4, §3:1, §8:2, §10:1, §12:13, §13:6, Appendix A:1), so a
per-block newline would add 28, not 7. The actual delta is 7 = one per *section*, i.e. my
round-02 extraction was one byte short at each section boundary. The stated cause would
produce a different number than the one recorded beside it.

For the record, and as a **SELF-CORRECTION (§15.1)** of my own round-02 claim: my round-02
sentence "My own whole-section measurement over the source gives 23,223 B for the seven
top-level omitted sections alone" (`review/round-02/claude-1.md:355`) is **superseded** —
the exact whole-section source total is **23,230 B**, which I have now measured twice.
This is a weakening of my own figure and takes effect immediately; it does not change any
finding, because R2-MAJ-1 never turned on that number.

**Suggested fix.** Replace the clause with the measured cause, e.g. "claude-1's round-02
23,223 B was one byte short per section at the extraction boundary; the exact source total
is 23,230 B", or simply drop the reconciliation and state the measured value.

### [NIT] "the old undercounted sum was 19,736 B mislabelled `floor`" attaches the label to the wrong number

Both numbers are real; the attribution is crossed. At `118b245` the variable named `floor`
held the heading-block-only sum and was **logged** as `named-omission-set bytes`; the
`+§2(…) reference` figure was the larger one. I measured both from the same source under
the old code's own `len(b.Text)` convention:

| quantity at `118b245` | value | how it was labelled |
|---|---|---|
| `floor` variable = Σ `byHeading[heading]` | **15,759 B** | logged `named-omission-set bytes` |
| §2 block | **3,977 B** | logged inside `+§2(…)` |
| `floor + section2` | **19,736 B** | logged `+§2(…) reference` |

So 19,736 B is the number that *stood in for* FINAL's "floor with §2" and is the right one
to compare against ≈ 42,100 B — the record's substantive point is sound — but it is not the
number the `floor` identifier was attached to. The `~11.9 KB` / 17-subsection parenthetical
is fine: 11,903 B is correct in the old sum's own `len(Text)` convention (I measure 11,920 B
under the new `len(Text)+1` convention, and the text marks it approximate).

**Suggested fix.** One clause: "the `floor` variable held 15,759 B, logged as
`named-omission-set bytes`; the `+§2 reference` it was compared against was 19,736 B".

### [NIT] The "no driver transition since F7 / zero-width manifests" residual is stale and mischaracterises the run set

The cycle-2 residuals say "no driver transition occurred during this cycle either — the
deck's `runs/` records are still the historical zero-width set", carrying forward the
consensus blind-spot 3 wording "the newest run record `20260924T003301…Z` predates it
[F7 at `64a622c`, 2026-09-24T01:48Z]". Two facts do not match the live deck:

- F7 is `64a622c`, committed `2026-09-24T03:48:45+02:00` = **01:48:45Z** (confirmed). The
  newest run record is **`20260924T025208.421874000Z`** = **02:52:08Z**, which **postdates
  F7 by ~63 minutes** and was created during this cycle (it is one of the untracked `runs/`
  directories in the working tree). Its `events.jsonl` opens with
  `{"type":"run.created","data":{"idea":"meta-protocol-change-lean-organizer","mode":"consensus-signoff",…}}`.
- "still the … zero-width set" understates the state. Of the four recent runs, only
  `20260923T202501Z` has a `run.json` at all (and its `created_at == updated_at`, so that
  one is genuinely zero-width); `20260923T211817Z`, `20260924T003301Z` and
  `20260924T025208Z` carry **only `events.jsonl` — no manifest at all**.

The conclusion these facts support — live attribution windows remain test-proven only — is
**unchanged and if anything reinforced** (a run with no manifest proves even less than a
zero-width one). Whether a `consensus-signoff` runner launch counts as a "driver transition"
is arguable; the dated claim about which record is newest is not.

**Suggested fix.** Restate as: "the newest run record is `20260924T025208Z` (02:52Z,
`mode: consensus-signoff`), which postdates F7; it and the two before it write no
`run.json` at all, and the one manifest that exists (`20260923T202501Z`) is zero-width — so
no live attribution window has yet been produced."

### [NIT] The cycle-2 record overstates what the G4(a) test asserts

`IMPLEMENTATION.md:112-113` credits `TestPhaseDigestNextActionFixUpPublishedAwaitsReview`
with "(both the fix-up-published state **and the newer-review-artifact relaxation**)". The
first half is genuinely asserted (`if d.Next != driver.NextAwaitReviewArtifact { t.Errorf }`).
The second half writes a newer review artifact and then only **logs** the outcome —

    if d2.Next == driver.NextAwaitReviewArtifact && d2.Review.Completed == d2.Review.Total {
        t.Logf("after the newer review artifact: next=%q", d2.Next)
    }
    if !driver.IsValidNextAction(d2.Next) { t.Errorf(...) }

— so the only assertion on the relaxation path is that the value stays inside the
enumeration, which every branch satisfies by construction. The in-test comment is honest
about this; the IMPLEMENTATION summary of it is not.

**Suggested fix.** Either assert the expected `d2.Next` value, or describe the second half
as "enumeration-validity only" in the record.

## Open questions

1. **Does the quorum want the ≈ 42 KB figure in `## Deviations from FINAL.md` as well as in
   `### Deviations from agreed fixes`?** The signed plan ruled "no deviation" for G1, and
   both reviewers confirmed branch 1 — but neither of us knew at signoff that the recomputed
   floor would be 52,295 B rather than the "≈ 42 KB" the plan's own G1 text predicted. I do
   **not** re-litigate the branch (it is correct); I ask only where the divergence is
   recorded. My own lean: the agreed-fixes deviation bullet (R3-MIN-2) is sufficient, and
   the FINAL-deviations section is for mechanism changes like G7's.
2. **Is a goal-done check going to be run, or is this close being taken outside the driver?**
   (R3-MIN-3.) Either answer is fine; the closing consensus should say which, because the
   two produce different records.
3. **For R3-MIN-1, does the quorum prefer the content signal or the `>=` relaxation?** I
   lean content signal — `status: fix-up-cycle-N` vs the latest review-round label is
   deterministic, survives any filesystem operation, and removes a mtime dependence rather
   than widening one. If the quorum judges R3-MIN-1 a deferral instead, the honest minimum is
   to soften the "retired" wording at `IMPLEMENTATION.md:503-508`.
4. **Should the four NITs be fixed in a cycle 3, or dispositioned by the closing consensus?**
   All four are record-accuracy items on this idea's own artifacts. `strict_gate` is absent,
   so NITs do not block a close by rule. I would accept either route (see Readiness).

## Position changes since prior review round

- **R2-MAJ-1 (floor not recorded / mislabelled / undercounted) — WITHDRAWN AS FIXED.** G1
  records a genuinely measured retained floor with its convention stated, adds a real
  `body ≥ floor` assertion, and logs the omission companion whole-section-correct. I verified
  52,295 / 30,262 / 23,230 independently of the implementer's test. **My round-02 position
  that C.3/R-2's floor half was an unmet frozen-FINAL acceptance element no longer holds.**
- **R2-MAJ-2 (`--json` trailer) — WITHDRAWN AS FIXED.** All four exit paths verified by me.
- **R2-MAJ-3 (CI timeout) — WITHDRAWN AS FIXED**, and strengthened: `internal/trajectory`
  took 644.7 s in my run at this HEAD, i.e. it would now fail Go's 600 s default outright.
  **The second unmet acceptance element from my round-02 is therefore also met.**
- **R2-MIN-2/-3/-4 and R2-NIT-1/-2/-4 — WITHDRAWN AS FIXED** (G5, G6, G7, G4(c), G8, G9).
  R2-NIT-3 stays retired-by-owner.
- **R2-MIN-1 — PARTIALLY FIXED; re-filed narrowly as R3-MIN-1.** (b) and (c) are fixed and
  verified. (a) is fixed only where mtime ordering happens to hold. I keep MINOR — the same
  severity I gave it in round 02 — rather than escalating: nothing gating is broken.
- **VC-4 — my "NOT READY" rationale is spent.** I argued NOT-READY on exactly two unmet
  frozen-FINAL acceptance elements; both are now met and I say so plainly. My round-03
  position is **not-yet-closing**, but on a materially smaller and different basis: a
  behavioral edge, a missing required section, an unstated close condition, and four record
  items — **no CRITICAL, no MAJOR, no unmet acceptance criterion**. This is a weaker
  objection than round 02's and I record it as such.
- **VC-5, VC-6 — closed; no residual position.** Both landed at the shape I picked (stderr;
  G4 whole). I verified the stderr shape rather than restating my preference for it.
- **SELF-CORRECTION (§15.1):** my round-02 figure "23,223 B for the seven top-level omitted
  sections" is superseded by the measured **23,230 B** (see R3-NIT-4). A weakening of my own
  claim; effective immediately.

## Responses to other reviewers

### @kimi-1

**Your K2-F2 was the round's most valuable catch, and G10 lands it cleanly.** I reproduced
your fixture again at this HEAD: exit 1 with the read failure named, exit 3 once
`COOPERATION.md` is present. I also went looking for the usability regression your
fail-closed leaning could have caused and did not find one — a directory with no
`parley-deck/` still gets the old `run \`parley init\` first` message, and a freshly
`init`ed deck gets exit 3, so only a genuinely broken deck hard-fails. Your leaning was
right and I have nothing to add to it.

**Your K2-F4 precision holds and is now correctly recorded** — and it turns out to be the
reason my own tree exposed R3-MIN-1 while the live deck hides it. The live deck's empty
`review/round-03/` directory routes through the directory-based path you isolated, which
already worked; my fresh checkout has no such directory and falls into the mtime path, where
G4(a) silently does not engage. Your distinction between the directory trigger and the file
trigger is what made that diagnosis legible, so I want it on record that the record half of
K2-F4 earned more than a record correction.

**K2-F1 / R2-NIT-4 — we measured the same 8,056 B and we were both right.** I re-measured in
all three copies this round; all three agree and share one region sha256. G9 also states the
extraction convention beside the number, which was the part neither of us had pinned.

**On your VC-4 self-correction.** You withdrew "ready for a zero-fix closing consensus" on
the strength of my R2-MAJ-1 and R2-MAJ-3, writing that you "hold no counter-evidence against
either". I should tell you plainly that **both of those elements are now met**, verified by
me independently this round — so the specific basis on which you withdrew has been
discharged. Your original instinct that this implementation was close to ready was better
calibrated than the two MAJORs made it look, and I do not think you should read my continued
not-yet-closing position as vindication of the withdrawal. My remaining objections are
smaller than the ones you deferred to, and a reasonable reviewer could disposition all seven
without another code cycle.

**Where I would value your independent check.** R3-MIN-1 is the one finding I would most
like a second pair of eyes on, because my reproduction depends on a filesystem property
rather than on program logic: clone or check out the repo fresh, confirm
`IMPLEMENTATION.md` and `review/round-02/*.md` share an mtime, and run
`parley organizer brief --idea meta-protocol-change-lean-organizer`. If you see `await
review artifact` where I see `await implementation`, my finding is wrong and I withdraw it.
I would also value your read on R3-MIN-3 — you are the reviewer who stated the close rule
most explicitly, and `auto_implement: true` changes it in a way neither of us caught.

**On the G2 envelope-freeze point you carried:** confirmed. `parley version` reports
`1.48.0` and `v1.48.0` is the newest tag, so both the round-1→fix-up-1 envelope change and
the G2 stream split are genuinely inside one unreleased 1.49.0 window. Nothing shipped
between them.

**@zcode-1** (implementer, no verdict owed to me): the fix-up is accurate work. Ten fixes,
all present, all measurably correct, and the two numbers most likely to be quietly wrong
(52,295 and 30,262) reproduce exactly under conventions you stated well enough for me to
check them without reading your test. My remaining items are four record clauses, one
required section, one mtime edge, and a close condition that is not yours to have caught.

## Updated findings

| id | severity | title | status |
|---|---|---|---|
| R3-MIN-1 | MINOR | G4(a) inert on equal mtimes (fresh clone/checkout) | **new**, narrows R2-MIN-1(a) |
| R3-MIN-2 | MINOR | Fix-up cycle 2 missing `### Deviations from agreed fixes` (G1 figure belongs there) | **new** |
| R3-MIN-3 | MINOR | `auto_implement: true` close conditions (LE-7/LE-11) unaccounted for | **new** |
| R3-NIT-1 | NIT | 23,230 ↔ 23,223 reconciliation arithmetically unexplained (+ my own SELF-CORRECTION) | **new** |
| R3-NIT-2 | NIT | 19,736 B labelled as the `floor` variable's value (it held 15,759 B) | **new** |
| R3-NIT-3 | NIT | run-record residual stale: newer post-F7 run exists; manifests absent, not zero-width | **new** |
| R3-NIT-4 | NIT | record overstates the G4(a) test's second assertion | **new** |
| R2-MAJ-1 | — | floor not recorded / mislabelled / undercounted | **withdrawn — fixed (G1), verified independently** |
| R2-MAJ-2 | — | `--json` non-JSON trailer on stdout | **withdrawn — fixed (G2), all four paths** |
| R2-MAJ-3 | — | CI leg without `-timeout` | **withdrawn — fixed (G3), reconfirmed at 644.7 s** |
| R2-MIN-1 | — | next-action wrong in two states | **(b) fixed; (a) re-filed as R3-MIN-1** |
| R2-MIN-2 | — | fail-open on unevaluable `to-user` notes | **withdrawn — fixed (G5), both branches live** |
| R2-MIN-3 | — | core-publish note stale facts | **withdrawn — fixed (G6), exact match** |
| R2-MIN-4 | — | A.4 mechanism removal absent from deviations | **withdrawn — fixed (G7)** |
| R2-NIT-1 | — | no "consensus.md not filed" outstanding line | **withdrawn — fixed (G4c), observed live** |
| R2-NIT-2 | — | mtime arrival approximation undocumented | **withdrawn — fixed (G8), citation accurate** |
| R2-NIT-3 | — | literal `\n` in organizer-notes.md | **retired by owner (confirmed again)** |
| R2-NIT-4 | — | §15 region figure 8,041 B | **withdrawn — fixed (G9), 8,056 B ×3** |

**Severity counts this round: 0 CRITICAL, 0 MAJOR, 3 MINOR, 4 NIT** (7 findings; all seven
are new, and 11 of my 11 round-02 findings are withdrawn as fixed except R2-MIN-1, whose (a)
half is narrowed into R3-MIN-1).

**Unresolved acceptance criteria: none.** Every row of FINAL's A–D acceptance table that I
re-attacked this round passes, including the two I named in round 02 as unmet — C.3/R-2's
floor-recording half (now a measured, asserted floor of 52,295 B recorded beside the
guardrail and `--optimize` baseline in one run) and cross-cutting "both suites green" at the
command FINAL names (CLI exit 0, 31/31 packages; skill exit 0, 399 pass). C.3's gating half,
the never-cut floor, `packet check`, the ≤ 20,000 B skill core, and the ≤ 8,192 B
byte-identical read-only brief all hold on my own measurements.

**Readiness for a zero-fix closing consensus: not yet, but the gap is now dispositional
rather than substantive.** No CRITICAL, no MAJOR, no unmet acceptance criterion, no broken
gating property, and the trajectory is the converging shape §4's stopping judgment describes
(round 1: 25 findings incl. 1 CRITICAL → round 2: 15, 0 CRITICAL, 3 MAJOR → round 3: 7, 0
CRITICAL, 0 MAJOR, all confined to fix-up code and its record). Cycle 3 of 5 on the
deliberation budget.

Concretely, I would sign a zero-fix closing consensus that does **either** of the following,
and I do not insist on the first:

1. lands R3-MIN-1 and R3-MIN-2 as a short cycle 3 and dispositions the four NITs; **or**
2. dispositions all seven by written signoff decision — provided R3-MIN-1's "retired"
   wording at `IMPLEMENTATION.md:503-508` is softened to match the mtime condition, since
   leaving an unconditional claim over conditional behavior is the record-accuracy class
   this idea has corrected in every prior cycle.

Either route additionally requires R3-MIN-3 answered in the record, because under
`auto_implement: true` a zero-Agreed-fixes consensus signed 🟡 does not close this idea —
`internal/driver/impl.go:285` escalates it — and the LE-7 goal-done check has not been run.
That is a statement about what will happen at close, not a new fix.

I am **not** calling this implementation complete. I have released, published, merged,
tagged and installed nothing; `parley protocol publish` was not run; no source file, FINAL,
peer artifact, consensus, organizer note or inbox note was modified in this invocation.
