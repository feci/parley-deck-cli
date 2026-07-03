---
agent: hermes-1
idea: track-aware-driver
round: 1
date: 2026-07-03
---

# track-aware-driver — hermes-1 round-01

Independent analysis. I read 00-prompt.md, the §4.0 track spec in COOPERATION.md
(lines 172–227), and the cited source: transport.go, driver.go, impl.go,
driver_impl.go, app.go, runtime.go, plus cursor.go, loop.go, consensus.go,
runner.go, and runcontrol.go for threading context. I did not read other agents'
round-01 files.

## 0. Mental model — what exists today

The driver today is effectively **single-track** (what §4.0 calls `deliberation`
with some knobs). Every `driver.Config` (driver.go:41-75) is built identically in
three sites in app.go (lines 1154, 1827, 1881 — `continueAuto`, `run --no-tui`,
and `startAutoDrive`), reading per-idea frontmatter fields:

- `CrossReviewRounds` ← `ReadCrossReviewRounds` (transport.go:33, default 1)
- `StrictGate` ← `ReadStrictGate` (transport.go:55, default false)
- `AutoImplement` ← `ReadAutoImplement` (transport.go:45, default false)
- `MaxFixupCycles` ← `New` default 3 (driver.go:91-93, not read from frontmatter)
- Reviewer count ← implicit: all non-implementer participants
  (driver_impl.go:38-64, `newDriverImplOps`)
- Per-agent timeout ← `runner.Options.Timeout` (0 → falls back to per-agent
  `TimeoutMS`, runner.go:1105-1117) — **not set** in the three Config sites; the
  `runOpts` in `continueAuto` (app.go:1146) omits Timeout entirely, and
  `runcontrol.Create` (runcontrol.go:103-110) also omits it. So every agent gets
  its `agents.DefaultTimeoutMS = 1_800_000` (30 min) default today.
- Readiness ping ← `runTaskPreflight` (preflight.go:231) runs before idea
  creation; gated by `--no-ping` / `centralPingSkips`.

Key insight: the three Config construction sites are copy-pasted with minor
variations. Any track→Config derivation must be a **single shared function**
called from all three, not three inline blocks — otherwise drift is guaranteed.

---

## 1. Classifier input model

**Recommendation: hybrid (c) primary + (a) optional — author declares `track:`,
CLI validates + warns on obvious under-tier; `parley classify` is a later slice.**

### Why not (a) flags-first for the MVP

(a) `parley classify --files N --loc N --security ...` is the cleanest, most
deterministic option, and I want it eventually. But for the MVP it has a boot
problem: the track must be known **before idea creation** (it seeds
`00-prompt.md` frontmatter, timeouts, and the preflight ping decision), and at
that point the diff/LOC/security surface is a free-text task description, not
structured data. Forcing the user to pass `--files N --loc N` at `parley run`
time for every idea is ceremony that §4.0's `fast` track was designed to
reduce. The classifier inputs are also genuinely uncertain pre-implementation
(§4.0 says "on doubt → stricter"), so a flag-based classifier that demands exact
numbers creates a false-precision hazard.

### Why not (b) git-diff inference

(b) `git diff --stat` inference is attractive for `auto_implement` ideas (the
diff exists post-implementation), but the track is **binding once Phase 0
closes** (§4.0 line 218) — before implementation. Inferring from a diff that
doesn't exist yet is impossible; inferring post-hoc violates binding. It's
useful only for the mid-idea upgrade signal (§4.0: "if implementation later
reveals a higher-risk surface"), which is a later slice.

### The hybrid (c)+(a) MVP

**Phase 0 (now):** the author declares `track:` in `00-prompt.md` (the template
already ships it per v1.32.0). `ReadTrack` (transport.go, new) reads it,
defaulting to `standard` on absent/unknown/empty — exactly mirroring
`ReadStrictGate`/`ReadAutoImplement`.

**Validation (MVP):** the CLI performs a **lightweight under-tier warn** — NOT a
full classifier. It checks the `00-prompt.md` frontmatter for **objective
deliberation triggers that are already machine-readable**:

- `auto_implement: true` → must be `deliberation` (§4.0 table: auto_implement is
  a deliberation trigger)
- `strict_gate: true` → must be `deliberation` (§4.0 table: strict_gate is a
  deliberation trigger)
- `track: fast` with either of those set → **hard reject** (not just warn):
  these are self-contradictory declarations.

Anything requiring judgment (files touched, LOC, security surface, protocol
change) is **not** machine-checkable from frontmatter alone and is left to the
author + the existing human review. This is deterministic and fail-safe: the
only automated action is **escalating** a contradiction, never downgrading.

**Later slice (a):** `parley classify` with explicit flags, usable as a
pre-flight check before writing `00-prompt.md`. This is where the full §4.0
classifier table (files/LOC/security/reversibility) gets encoded. It prints the
recommended track and the trigger that fired. The author pastes it into
`track:`. This is high-value but not MVP-blocking — the binding declaration +
contradiction reject covers the safety-critical path.

### `parley classify` interface (for the later slice, specified now)

```
parley classify [--files N] [--loc N] [--security] [--irreversible]
                [--protocol-change] [--auto-implement] [--strict-gate]
                [--pipeline] [--api-break] [--schema-break]
                [--json]
```

Deterministic rules (§4.0 lines 179-190, first-match-wins, fail-closed):

1. Evaluate all `deliberation` triggers; if ANY is set → `deliberation`.
2. Else evaluate `fast` conditions: ALL of (no --security, no --irreversible,
   --files ≤ 5, --loc ≤ 300, no --protocol-change, no --auto-implement, no
   --strict-gate, no --pipeline, no --api-break, no --schema-break) → `fast`.
3. Else → `standard`.
4. Boundary cases (files 6-14, LOC 301-1000, any unconfirmed trigger) → fail
   closed to the stricter track.

Output: `track: standard\ntrigger: default (no deliberation or fast conditions met)\n`
(or `--json`: `{"track":"standard","trigger":"default",...}`).

---

## 2. Track → Config mapping

### 2.1 The mapping table

§4.0 lines 194-203 defines 8 per-track aspects. Here's how each maps to the
existing `driver.Config` (driver.go:41-75) and related structures:

| §4.0 aspect | fast | standard | deliberation | Config field / threading |
|---|---|---|---|---|
| Readiness ping (§9.0) | skipped | full | full | **NO Config field** — preflight.go:231; see §2.3 |
| Cross-review rounds (Phase 2) | 0 | ≤2 | current (1, from frontmatter) | `CrossReviewRounds` (driver.go:48) |
| Consensus + FINAL (Phase 3-4) | collapsed: one FINAL.md with embedded signoffs | separate, drafted simultaneously | separate | **NO Config field** — consensus.go:39-127; see §2.4 |
| Reviewers (Phase 6) | 1 (model-diverse) | 2 | all non-implementers | **NO Config field** — driver_impl.go:38-64; see §2.5 |
| Review consensus (Phase 7) | 1 reviewer's ✅ = consensus | reviewers who reviewed sign off | all participants sign off | **NO Config field** — consensus.Status; see §2.4 |
| Fix-up (Phase 8) | cap 1 cycle | cap 2 cycles | unbounded; strict_gate available | `MaxFixupCycles` (driver.go:61) + `StrictGate` (driver.go:66) |
| Timeout per agent | ~5 min | ~15 min | ~30 min | **NO Config field** — runner.Options.Timeout; see §2.6 |
| Auto-advance | full (pause only for one signoff) | auto-advance; human gate at FINAL→impl | human gate at each transition | **Partially covered** — `Auto` (driver.go:50), `AutoImplement` (driver.go:59); see §2.7 |

### 2.2 Where the derivation plugs in

A single function in `internal/app` (not `internal/driver` — the driver core
must stay app-agnostic per driver.go:40 comment):

```go
// internal/app/track_config.go (new file)
func applyTrack(cfg *driver.Config, ideaDir string, runOpts *runner.Options)
```

Called from all three Config sites (app.go:1154, 1827, 1881) right after the
struct literal, before `driver.New`. It reads `driver.ReadTrack(ideaDir)` and
mutates `cfg` (and `runOpts` for timeout). This centralizes the mapping and
prevents the three-site drift.

### 2.3 Readiness ping — no Config field

The ping runs in `runTaskPreflight` (preflight.go:231) **before** idea creation,
so it can't read `driver.Config`. For `fast`, the ping should be skipped. The
cleanest MVP approach: `runTaskPreflight` reads the track from the **task
description / proposed 00-prompt** — but that doesn't exist yet at preflight
time (the idea is created by `runcontrol.Create` at app.go:1796, after
preflight at app.go:1785).

**MVP decision:** defer ping-skipping to a later slice. The ping is cheap
(presence-only with `--no-ping`), and §4.0 says `fast` skips it for speed — but
the safety cost of running it is near-zero (a few seconds). The value of
skipping it is low relative to the wiring complexity. Flag this as a §4.0
behavior with no current field; implement in slice 2 by having `runcontrol.Create`
write `track:` to `00-prompt.md` first, then `runTaskPreflight` reads it —
requires reordering preflight after idea creation, which is a bigger change.

### 2.4 Collapsed consensus/FINAL for fast — no Config field

`advanceConsensus` (consensus.go:39-127) drafts consensus.md → FINAL.md as
separate steps. §4.0 says `fast` collapses this to "one FINAL.md with embedded
signoffs." This is a **semantic change to the consensus gate**, not a Config
field.

**MVP decision:** defer. `fast` track in the MVP runs the same consensus/FINAL
flow as `standard`. The collapsed-FINAL optimization is pure ceremony-reduction
(no safety impact — the invariants are enforced elsewhere). Implement in a
later slice by adding a `CollapsedFinal bool` to `Config` (or a `Track` field
that `advanceConsensus` switches on), and a fast-path in `advanceConsensus` that
skips consensus.md and goes straight to FINAL.md with embedded signoffs. This
requires consensus-package changes (the signoff validation expects
consensus.md), so it's non-trivial.

### 2.5 Reviewer count — no Config field, but existing logic to amend

Today `newDriverImplOps` (driver_impl.go:38-64) sets `reviewers` to ALL
non-implementer participants. `ReviewStatus.ReviewerCount` (impl.go:52) is
`len(o.reviewers)`. The LE-11 guard (impl.go:240) rejects `< 2` under
`AutoImplement`.

For tracks:
- `fast`: 1 reviewer → need to **limit** the reviewer set to 1 (the first
  non-implementer), and the LE-11 `< 2` guard must be relaxed to `< 1` for fast.
- `standard`: 2 reviewers → need to **limit** to 2 (with the §4.0 degradation:
  "with only two participants, standard's 2 reviewers degrades to 1").
- `deliberation`: all non-implementers → **today's behavior, no change**.

**Minimal code:** add `MinReviewers int` to `driver.Config` (driver.go). In
`newDriverImplOps`, accept a `track` parameter (or `MinReviewers`) and truncate
`reviewers` to `MinReviewers` if `len(reviewers) > MinReviewers`. In
`advanceReview` (impl.go:240), replace the hardcoded `< 2` with
`< d.cfg.MinReviewers` (where `MinReviewers` is 0 for fast, 2 for standard, 2
for deliberation — deliberation keeps the existing LE-11 guard since
`auto_implement` already implies deliberation).

Wait — `MinReviewers` is the **required** count, not the **available** count.
The §4.0 table says fast *uses* 1, standard *uses* 2, deliberation *uses* all.
The constraint is: the reviewer set is **capped** at the track's count (fast=1,
standard=2, deliberation=∞), and the LE-11 **minimum** is the same number
(fast≥1, standard≥2 but degrades to ≥1 with 2 participants, deliberation≥2
under auto_implement).

Cleaner design: two fields.

```go
// driver.Config (new fields)
MaxReviewers int  // cap on the reviewer set; 0 = unlimited (deliberation)
MinReviewers int  // LE-11 minimum for auto-complete; 0 = no minimum
```

- fast: `MaxReviewers=1, MinReviewers=1`
- standard: `MaxReviewers=2, MinReviewers=2` (but see degradation below)
- deliberation: `MaxReviewers=0, MinReviewers=2` (today's effective behavior)

**Two-participant degradation (§4.0 line 225):** "with only two participants,
standard's 2 reviewers degrades to 1." This means if `len(participants)==2`
(1 implementer + 1 reviewer), standard's `MinReviewers` drops to 1. This logic
lives in `applyTrack`: if track==standard and `len(participants)<=2`,
`MinReviewers=1`.

### 2.6 Per-agent timeout — no Config field

`runner.Options.Timeout` (runner.go:28) is the per-agent timeout override;
`timeoutForAgent` (runner.go:1105) uses it if > 0, else falls back to the
agent's `TimeoutMS`. Today it's 0 in all three Config sites, so every agent gets
`DefaultTimeoutMS = 1_800_000` (30 min).

§4.0: fast ~5 min, standard ~15 min, deliberation ~30 min.

**Minimal code:** `applyTrack` sets `runOpts.Timeout` based on track:
- fast: 5 * time.Minute
- standard: 15 * time.Minute
- deliberation: 30 * time.Minute (or 0 to keep the existing per-agent fallback,
  which is also 30 min — equivalent, but explicit is clearer and survives
  DefaultTimeoutMS changes).

This threads through `runner.Options` → `timeoutForAgent` →
`context.WithTimeout` (runner.go:507-508). No Config field needed — it's a
runner option, not a driver config. But `applyTrack` needs a pointer to
`runOpts` to mutate it.

**Also:** the `roundDeadline` constant (loop.go:18, 30 min) bounds how long a
tick waits for an incomplete round. For fast, this should probably also shrink
(to ~10 min, since the per-agent timeout is 5 min and there's 1 round + 1
review). But `roundDeadline` is a package const, not per-Config. **MVP decision:**
leave `roundDeadline` at 30 min for all tracks — it's a tick-level deadline, not
per-agent, and shrinking it for fast risks escalating prematurely if agents are
slow. Flag for later: make `roundDeadline` a Config field.

### 2.7 Auto-advance behavior

§4.0: fast = "full (pause only for the one signoff)", standard = "auto-advance;
human gate at FINAL→implementation", deliberation = "human gate at each
transition."

Today: `Auto` (driver.go:50) is the master switch; `AutoImplement`
(driver.go:59) gates the FINAL→impl transition. The "human gate at each
transition" for deliberation is **not** currently enforced — the driver
auto-advances through all phases under `--auto`. This is a §4.0 behavior with no
current field.

**MVP decision:** defer the "human gate at each transition" for deliberation.
Today's behavior (full auto-advance under `--auto`) is the **deliberation**
behavior per the constraint "deliberation behavior must remain EXACTLY today's
full lifecycle" (00-prompt.md line 56-57). So the §4.0 table's "human gate at
each transition" for deliberation is **aspirational relative to the current
code** — today's auto-drive IS the deliberation behavior. Do not change it.

For `fast`'s "pause only for the one signoff" — the fast track in the MVP still
runs through the consensus signoff path (since collapsed FINAL is deferred).
The existing signoff gate (consensus.go:74-88, `TriagePartial` →
`RequestSignoffs` → escalate if still missing) already pauses for signoffs. So
fast's auto-advance is effectively the same as standard's in the MVP, just with
fewer rounds and reviewers.

### 2.8 Summary of Config changes

New `driver.Config` fields:
```go
Track        string  // "fast" | "standard" | "deliberation" (for advanceReview branching)
MaxReviewers int     // cap on reviewer set; 0 = unlimited
MinReviewers int     // LE-11 minimum for auto-complete; 0 = no minimum
```

Existing fields whose values change per track:
```go
CrossReviewRounds int  // fast=0, standard=min(2, frontmatter), deliberation=frontmatter
MaxFixupCycles    int  // fast=1, standard=2, deliberation=3 (current default)
StrictGate        bool // fast=false, standard=false, deliberation=frontmatter (today's)
```

Threaded via `applyTrack(cfg, ideaDir, runOpts)` in app.go, called from all
three Config sites. Timeout via `runOpts.Timeout`.

### 2.9 Fields NOT changed per track (to preserve deliberation = today)

- `MaxRounds` (driver.go:49, default 4) — stays the same; it's a circuit
  breaker, not a track knob.
- `MaxDriverSteps`, `MaxWallClock`, `MaxCostUSD` (driver.go:71-73) — LE-5 loop
  budgets, not track-specific. Stay from `loopBudget(root)`.
- `RequireModelDiversity` — read per-idea from frontmatter
  (transport.go:65-70), not track-specific. But §4.0 says fast requires
  "model-diverse" reviewer. **MVP decision:** for fast, `applyTrack` sets
  `RequireModelDiversity=true` if not already set in frontmatter (fast's single
  reviewer MUST be model-diverse per §4.0). This is a behavioral addition, but
  it's fail-safe (stricter, not weaker).

---

## 3. Invariant enforcement

§4.0 lines 211-216 defines invariants that hold on **every** track: non-solo
(at least one independent non-facilitator artifact), refutation-default review,
round-1 independence, append-only signoffs, files-canonical audit trail, §14
human brake, English-only, no-secrets.

The two the driver must **hard-reject** (00-prompt.md deliverable 3):
1. **0 independent reviewers (non-solo)** — a config that would drop the
   reviewer count to 0.
2. **A review path with no refutation** — a config that would skip review
   entirely for an auto_implement idea.

### 3.1 Where: `applyTrack` validation (pre-construction)

`applyTrack` runs after reading `track:` and before `driver.New`. It validates
the derived config and **returns an error** (which the caller surfaces as a hard
stop) if:

**Non-solo violation:**
- `MaxReviewers` is set (fast=1, standard=2) AND `len(participants) -
  1(implementer) < 1` → error: "track fast/standard requires at least 1
  independent reviewer; only N participant(s) declared."
- More precisely: `availableReviewers = len(participants) - 1` (the
  implementer). If `availableReviewers < 1` → hard reject regardless of track.
  This is the §1 non-solo invariant, already enforced in preflight
  (preflight.go:245-246), but `applyTrack` is a second line of defense.

**Refutation violation:**
- `auto_implement: true` AND `track: fast` with `MaxReviewers=0` → impossible by
  construction (fast sets MaxReviewers=1), but `applyTrack` asserts
  `MinReviewers >= 1` when `AutoImplement` is true. If somehow
  `MinReviewers < 1` and `AutoImplement` → error.
- `auto_implement: true` AND `track` would skip review entirely → there's no
  "skip review" config today; review always runs if `Impl` is wired. But if a
  future "collapsed" fast path skips Phase 6-7, this guard must fire. For the
  MVP, the review path always runs, so this is structurally enforced.

### 3.2 Where: `advanceReview` runtime guard (impl.go:236-243)

The existing LE-11 guard (impl.go:240) is the **runtime** enforcement:
```go
if d.cfg.AutoImplement {
    if rs.ReviewerCount < 2 {  // ← hardcoded 2
        return ActionEscalated, ...
    }
}
```

This must be generalized to `rs.ReviewerCount < d.cfg.MinReviewers`:
```go
if d.cfg.AutoImplement && d.cfg.MinReviewers > 0 {
    if rs.ReviewerCount < d.cfg.MinReviewers {
        return ActionEscalated, c, fmt.Errorf(
            "only %d independent reviewer(s); track %s requires at least %d (LE-11) — add a reviewer or sign off manually",
            rs.ReviewerCount, d.cfg.Track, d.cfg.MinReviewers)
    }
}
```

For `fast` (MinReviewers=1): a solo reviewer is allowed (it's the track's
design), but 0 reviewers is still rejected (the `< 1` check). This preserves
non-solo while allowing fast's single-reviewer model.

For `deliberation` (MinReviewers=2 under AutoImplement): identical to today.

### 3.3 The contradiction reject (from §1)

`applyTrack` also hard-rejects self-contradictory declarations:
- `track: fast` + `auto_implement: true` → error (auto_implement is a
  deliberation trigger per §4.0).
- `track: fast` + `strict_gate: true` → error (strict_gate is a deliberation
  trigger).

This is the **fail-safe classifier** for the MVP: the only automated
classification action is rejecting contradictions. Everything else is
author-declared.

### 3.4 What about "no review path with no refutation"?

The refutation-default review (LE-1) is a **reviewer behavior** requirement,
not a driver-config field. The driver enforces it structurally: review always
runs (Phase 6) when `Impl` is wired, and the review-consensus drafter certifies
findings. The driver can't verify "did the reviewer actually attempt
refutation" — that's a content quality issue. The structural invariant the
driver enforces is: **review happens** (OpenReviewRound is called, reviewers
produce artifacts, DraftReviewConsensus runs). The `MinReviewers >= 1` guard
ensures at least one reviewer exists. That's the driver's contribution to
refutation: it can't force the reviewer to refute, but it can ensure the review
path isn't skipped.

---

## 4. MVP slicing

Ordered by (safety + value) ÷ risk. Each slice is independently shippable.

### Slice 1: ReadTrack + Config derivation + invariant gate (HIGH value, LOW risk)

**What:**
- `ReadTrack(ideaDir)` in transport.go (mirrors ReadStrictGate, default
  `standard`).
- `applyTrack(cfg, ideaDir, runOpts)` in `internal/app/track_config.go` (new).
- New `Config` fields: `Track`, `MaxReviewers`, `MinReviewers`.
- `applyTrack` sets: CrossReviewRounds (fast=0, standard=min(2,frontmatter),
  deliberation=frontmatter), MaxFixupCycles (fast=1, standard=2,
  deliberation=3/default), MaxReviewers (fast=1, standard=2, deliberation=0),
  MinReviewers (fast=1, standard=2→1 if 2 participants, deliberation=2),
  runOpts.Timeout (5/15/30 min).
- Contradiction reject: fast + auto_implement → error; fast + strict_gate →
  error.
- Generalize LE-11 guard (impl.go:240) to `MinReviewers`.
- `newDriverImplOps` truncates reviewers to `MaxReviewers` (if > 0).
- Call `applyTrack` from all three Config sites (app.go:1154, 1827, 1881).

**Safety:** enforces non-solo and contradiction rejects. **Value:** fast/standard
tracks actually work with reduced ceremony. **Risk:** low — additive fields,
default `standard` preserves today's behavior, existing tests stay green.

**Verify:** `go test ./internal/driver/...` passes; new tests for
`ReadTrack`, `applyTrack` per-track values, contradiction reject, MinReviewers
guard with fast=1.

### Slice 2: `parley classify` command (MEDIUM value, LOW risk)

**What:** the `parley classify` flags command from §1. Pure read-only
computation, no side effects. Prints the track + trigger. Author pastes into
`00-prompt.md`.

**Safety:** none directly (it's advisory). **Value:** closes the
"deterministic + script-checkable" requirement. **Risk:** zero — new command,
no changes to existing paths.

**Verify:** unit tests for the classifier function with all §4.0 triggers and
boundary cases.

### Slice 3: Ping skip for fast (LOW value, MEDIUM risk)

**What:** reorder `runTaskPreflight` to run after idea creation (or read track
from the task/proposed frontmatter pre-creation), skip hosted-PONG ping for
`fast`.

**Safety:** minimal (ping is advisory). **Value:** speed for fast. **Risk:**
medium — reordering preflight relative to idea creation changes a
well-established flow; needs careful handling to not create half-open ideas on
gate failure.

**Verify:** preflight tests with track=fast skip ping; track=standard runs ping.

### Slice 4: Collapsed consensus/FINAL for fast (MEDIUM value, MEDIUM risk)

**What:** fast-path in `advanceConsensus` (consensus.go:39) for `fast` track:
skip consensus.md, draft FINAL.md directly with embedded signoffs. Requires
consensus-package changes (signoff validation expects consensus.md).

**Safety:** low (ceremony reduction, invariants unchanged). **Value:** major
speed win for fast. **Risk:** medium — touches the consensus gate, which is
safety-critical; needs new validation for "embedded signoffs in FINAL.md."

**Verify:** consensus tests with track=fast; FINAL.md signoff validation.

### Slice 5: Mid-idea upgrade signal (LOW value, MEDIUM risk)

**What:** post-implementation, scan the git diff for deliberation triggers
(security files, >15 files, >1000 LOC) and warn if track < deliberation. §4.0's
"if implementation later reveals a higher-risk surface, force-upgrade."

**Safety:** high (catches under-tiering). **Value:** medium. **Risk:** medium —
diff scanning is heuristic, false positives annoy.

**Verify:** upgrade-signal tests with mock diffs.

---

## 5. Backward-compat & test plan

### 5.1 The invariant: standard/absent ≡ today

`ReadTrack` defaults to `standard` on absent/empty/unknown track. `applyTrack`
with `standard` must produce a `Config` identical to today's effective defaults:

| Field | Today's effective value | standard must set |
|---|---|---|
| CrossReviewRounds | frontmatter (default 1) | min(2, frontmatter) — **DIFFERENT if frontmatter > 2** |
| MaxFixupCycles | 3 (New default) | 2 — **DIFFERENT** |
| MaxReviewers | ∞ (all non-implementers) | 2 — **DIFFERENT** |
| MinReviewers | 2 (LE-11 hardcoded) | 2 — same |
| Timeout | 0 → 30 min fallback | 15 min — **DIFFERENT** |
| StrictGate | frontmatter | frontmatter — same |

**PROBLEM:** `standard` per §4.0 is NOT today's defaults. §4.0 standard caps
cross-review at 2, fix-up at 2, reviewers at 2, timeout at 15 min. Today's
defaults are cross-review=1 (but unbounded via frontmatter), fix-up=3, reviewers
=all, timeout=30 min.

**Resolution:** The constraint (00-prompt.md line 58-59) says "standard must
reproduce today's effective defaults (no silent behavior change for existing
ideas)." But §4.0 (COOPERATION.md line 205-206) says "This table is the single
authoritative per-track gate. It OVERRIDES the full-lifecycle defaults."

These conflict. The resolution is: **existing ideas (no `track:` field) get
today's defaults, NOT §4.0 standard.** An idea that explicitly declares
`track: standard` gets §4.0 standard. This means:

- `track:` absent/empty/unknown → **today's defaults** (the backward-compat
  path). `ReadTrack` returns `"standard"` as the normalized value, but
  `applyTrack` treats absent differently from explicit `standard`.

**Implementation:** `ReadTrack` returns the raw string + a `present` bool (like
`readFrontmatterField`). `applyTrack` checks: if `track` is absent → set
`Config.Track = "standard"` but apply **today's defaults** (no cap changes, no
timeout change, MaxReviewers=0, MinReviewers=2 under AutoImplement). If `track`
is explicitly `standard` → apply §4.0 standard values.

Alternatively: `ReadTrack` returns `"standard"` for absent, and `applyTrack`
treats `"standard"` as "today's defaults" (not §4.0 standard). Then §4.0
standard is only activated by... what? This doesn't work — §4.0 says standard
IS the default track.

**Cleaner resolution:** `applyTrack` distinguishes "absent" from "explicit
standard" via the `present` bool from `readFrontmatterField`. Absent → today's
defaults (backward compat). Explicit `standard` → §4.0 standard. This means
existing ideas (pre-v1.32.0, no track:) keep today's behavior, and new ideas
that explicitly declare `track: standard` get the §4.0 standard caps. The
00-prompt template (v1.32.0) already includes `track:` so new ideas will have it
explicit.

**This is the key backward-compat decision.** I recommend the `present` bool
approach. It means:

- `track:` absent → `Config.Track = "standard"`, but CrossReviewRounds,
  MaxFixupCycles, MaxReviewers, Timeout keep today's values. MinReviewers=2
  (today's LE-11). This is **byte-for-byte today's behavior.**
- `track: standard` (explicit) → §4.0 standard: CrossReviewRounds=min(2,fm),
  MaxFixupCycles=2, MaxReviewers=2, MinReviewers=2 (→1 if 2 participants),
  Timeout=15min.
- `track: fast` → §4.0 fast.
- `track: deliberation` → §4.0 deliberation = today's defaults (the constraint
  says deliberation = today's full lifecycle). So deliberation and absent
  produce the same Config, except deliberation sets `Track="deliberation"`
  explicitly (for future branching).

### 5.2 Test plan

**Existing tests (must stay green):**
- `go test ./internal/driver/...` — all existing tests in driver_test.go,
  strict_gate_test.go, close_integrity_test.go, impl_test.go. These construct
  Config directly (not via applyTrack), so they're unaffected by the track
  derivation. The LE-11 generalization (impl.go:240 → MinReviewers) must keep
  them green: existing tests use `ReviewerCount: 2` and the drivers set
  `AutoImplement: true` with no `MinReviewers` → `New` defaults MinReviewers to
  2 (backward-compat default in `New`).
- `go test ./internal/app/...` — app_test.go, preflight_test.go. These use
  `--no-preflight` and don't set `track:`, so `applyTrack` hits the absent path
  → today's defaults.

**New tests:**

1. `transport_test.go` — `TestReadTrack`:
   - absent track → ("standard", false)
   - `track: standard` → ("standard", true)
   - `track: fast` → ("fast", true)
   - `track: deliberation` → ("deliberation", true)
   - `track: unknown` → ("standard", false) (unknown normalizes to standard,
     but `present` is... hmm, the field IS present, just unknown. Return
     ("standard", true) for present-but-unknown? No — unknown should be treated
     as absent for backward compat. Return ("standard", false).)
   - `track: FAST` → ("fast", true) (case-insensitive)

2. `track_config_test.go` — `TestApplyTrack`:
   - absent track → Config matches today's defaults (CrossReviewRounds=1,
     MaxFixupCycles=3, MaxReviewers=0, MinReviewers=2, Timeout unchanged)
   - explicit standard → §4.0 standard values
   - fast → §4.0 fast values (CrossReviewRounds=0, MaxFixupCycles=1,
     MaxReviewers=1, MinReviewers=1, Timeout=5min)
   - deliberation → today's defaults + Track="deliberation"
   - fast + auto_implement → error
   - fast + strict_gate → error
   - standard + 2 participants → MinReviewers=1 (degradation)

3. `impl_test.go` — `TestMinReviewersGuard`:
   - fast (MinReviewers=1) with ReviewerCount=1 → completes (not escalated)
   - fast with ReviewerCount=0 → escalated (non-solo)
   - standard (MinReviewers=2) with ReviewerCount=1 → escalated
   - absent track (MinReviewers=2, today's default) with ReviewerCount=1 →
     escalated (today's behavior preserved)

4. `driver_impl_test.go` — `TestReviewerTruncation`:
   - fast (MaxReviewers=1) with 3 non-implementers → reviewers truncated to 1
   - standard (MaxReviewers=2) with 3 → truncated to 2
   - deliberation (MaxReviewers=0) with 3 → all 3 (unchanged)
   - absent (MaxReviewers=0) with 3 → all 3 (today's behavior)

**`New` defaults (driver.go:84-98):** add `MinReviewers` default: if
`MinReviewers <= 0 && AutoImplement` → `MinReviewers = 2` (preserve LE-11
default for configs that don't go through `applyTrack` — i.e., existing tests).

---

## 6. Risks

### R1: The absent-vs-explicit-standard distinction is subtle

The `present` bool from `readFrontmatterField` determines whether §4.0 standard
caps apply. If `readFrontmatterField` has a parsing edge case (e.g., `track:`
with trailing whitespace, `track: "standard"` with quotes), an explicit standard
could be misread as absent, silently applying today's looser defaults instead of
§4.0's caps. **Mitigation:** `ReadTrack` uses the same `strings.Trim` +
`strings.EqualFold` pattern as `ReadStrictGate` (transport.go:56-60), which
already handles quotes/whitespace. Test all variations.

### R2: MaxReviewers truncation changes the reviewer set for standard

Today, all non-implementers review. With `track: standard`, only 2 review. If
an idea has 4 participants (1 impl + 3 reviewers), standard drops 1 reviewer.
That reviewer's signoff is no longer requested, which could surprise a
participant who expected to review. **Mitigation:** this is §4.0's design
(ceremony scales to risk). The `applyTrack` output should log which reviewers
were selected/truncated. Not a safety risk — the invariant is ≥2, and 2 is
met.

### R3: Timeout reduction for standard (30→15 min) could cause premature kills

Today's 30-min timeout is generous. Reducing to 15 min for explicit `standard`
could kill agents that need 20 min for a complex round. **Mitigation:** 15 min
is §4.0's spec. If an agent needs more, the idea should be `deliberation` (30
min). The timeout is per-agent, not per-round, and the `roundDeadline` (30 min)
is unchanged. An agent killed at 15 min escalates (doesn't silently fail).

### R4: The three Config sites could drift if applyTrack isn't called from all three

app.go has three `driver.Config{...}` literals (lines 1154, 1827, 1881). If
`applyTrack` is added to two but not the third (e.g., the `startAutoDrive`
closure at 1881), that path would use today's defaults regardless of track.
**Mitigation:** the `applyTrack` function is the single point of truth; a grep
for `driver.Config{` should find exactly three sites, all calling `applyTrack`.
A lint test could assert this.

### R5: `fast` + `auto_implement` contradiction reject could break a valid workflow

§4.0 says auto_implement is a deliberation trigger, so `fast` + auto_implement
is contradictory. But an author might legitimately want a fast implementation
with 1 reviewer. **Mitigation:** this is §4.0's normative rule, not our
invention. The reject message should say "auto_implement is a deliberation
trigger per §4.0; use track: deliberation or remove auto_implement." The author
can declare `deliberation` (which is today's behavior). Not a regression —
today auto_implement ideas already run the full lifecycle.

### R6: MaxFixupCycles reduction (3→2 for standard, 3→1 for fast) could escalate more

Today's default is 3 fix-up cycles. Reducing to 2 (standard) or 1 (fast) means
ideas that need 3 fix-up cycles will escalate instead of completing. **Mitigation:**
§4.0 specifies these caps. Escalation is the correct behavior at the cap (not
silent completion). The author can declare `deliberation` for the full 3-cycle
budget. Existing ideas (no track) keep 3.

### R7: `MinReviewers` default in `New` could mask missing applyTrack

If `New` defaults `MinReviewers` to 2 (for backward compat with tests that
build Config directly), a production path that forgets `applyTrack` would get
MinReviewers=2 — today's behavior, not a safety regression. But it would also
not get the fast/standard caps. **Mitigation:** this is the same as R4 — the
single-point-of-truth pattern. The `New` default is a safety net, not the
primary path.

### R8: `ReadTrack` for unknown values → standard (not error)

An author writing `track: speedy` gets `standard` silently. §4.0 says
"unknown/empty → standard" (00-prompt.md line 23). This is correct per spec but
could mask a typo. **Mitigation:** `applyTrack` could warn (not error) on
unknown track values: "warning: track 'speedy' is not fast/standard/deliberation;
defaulting to standard." This is fail-safe (stricter-or-equal, never weaker).
