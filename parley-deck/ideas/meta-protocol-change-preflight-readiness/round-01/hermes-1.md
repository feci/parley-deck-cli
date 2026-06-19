---
agent: hermes-1
idea: meta-protocol-change-preflight-readiness
round: 1
lens: performance & simplicity
date: 2026-06-19
---

# Summary

The three locked decisions are mechanically sound. Through the performance/simplicity
lens the per-idea cost can collapse from the brief's implied "4 hosted PONGs + a
freshness sync at every idea start" to **one bounded 10s subprocess + the free CLI
discovery that `parley run` already performs**, by reusing two primitives that exist
today and letting the runner's supervised round-01 dispatch *be* the hosted-liveness
probe. The two new user-confirmed gates never block on stdin in the auto/CI path, so
unattended `parley run` cannot deadlock. The protocol text changes are small; the CLI
addition is one command wrapping one shared function. Concrete file:line references are
to this repo so the panel can verify the reuse claims.

# Findings & refinements

## A. Freshness check — reuse `parleyDeckSkillStatus`, compute once, no file cache

The expensive part of freshness already exists and is already bounded:
`internal/app/version_status.go:70-100` (`parleyDeckSkillStatus`) shells out to
`parley-deck-skill status --target all --project <root> --json` under a **10s
`context.WithTimeout`**, with a 5s `--version` legacy fallback. `parley preflight`
should call this SAME function — no second subprocess, no second code path. The
drift decision itself is then one string equality on `meta/version.json`
(`protocolSha256` vs the status payload's packaged hash) — effectively free.

Cost model that matters: in `parley run`, freshness is paid **once per process**. The
auto-driver advancing across rounds/ideas is the same OS process, so the 10s is a
one-time startup cost, not a per-idea tax. That makes a cross-process/file TTL cache
unnecessary for v1 — an in-process "compute once into a variable" is enough. I'd defer
any `meta/version.json`-persisted `packagedProtocolCheckedAt` TTL to v2, and only if
`parley preflight` gets hammered in a shell loop. A `--no-preflight` (or
`PREFLIGHT_SKIP`) escape hatch is the cheaper answer for CI that already vetted
freshness upstream.

Source vs consumer (brief Q1): prefer the explicit `meta/version.json.protocolRole:
"source"|"consumer"` field over deriving "which side is newer" from semver. It's one
string read vs a heuristic engine that has to handle pre-release/non-semver local
builds, and the source repo already knows it is the source. In a `source` repo the
freshness path is a cheap early-return: run status once (10s), print "advisory only —
you are ahead of the published skill", **write nothing**. The source-vs-consumer
inversion pitfall is then structurally impossible, not guarded by a fragile comparison.

Additive vs breaking (brief Q2): for v1 use the **semver of `deckVersion`** as the
single signal — major bump = breaking (pause for confirm); minor/patch = additive
(auto-sync). Do NOT build a structural protocol-section diff engine for v1. The safety
net for the residual risk (a minor that silently changed a rule) is the existing
allowlist zone-preservation on the write + the `meta/protocol-sync_<date>.md` record +
diff summary that the locked decision already mandates — a human can audit post-hoc,
and the §7 carve-out means this is not a protocol-change idea. A section-hash diff is a
clean v2 add-on; it is not needed to satisfy the locked decision.

Write surface in consumers is one file: the "two-copy drift-guard lockstep" (decision 3)
is a source-repo concern, and source = advisory/off. A consumer repo has only its
project `parley-deck/COOPERATION.md`; the embedded default lives in the installed
skill, not the consumer tree. So auto-sync writes a single file preserving the
project-specific zones (header, §2 roster table at COOPERATION.md:68, §0 transport) via
the same allowlist discipline. Simpler than "two copies" sounds.

## B. Roster ping — Tier-0 (free) by default; do NOT pre-pay a hosted PONG

This is the biggest win. The brief and `00-prompt.md` recorded "real PONG round-trip,
bounded" for 4 agents = 4 hosted round-trips per idea. That is 4× latency + 4× tokens
before any real work, every idea, compounding across an auto-drive session. It is also
paying twice for one signal:

- `parley run` already calls `discoverConfigured` (`internal/app/app.go:1693`) →
  `agents.Discover` (`internal/agents/discover.go:261`), which is `exec.LookPath`
  (the `command -v`) + a `--version` probe under a **4s cap** (`discover.go:446`). Zero
  tokens, local only. `Found && Version != "" && Error == ""` = available; missing CLI
  or probe timeout/error = unavailable. **This is the roster ping for v1.** No new
  subprocess, no new code path.
- The actual hosted liveness is tested by the real round-01 dispatch, which is already
  supervised by the three-layer watchdog (`internal/runner/supervision.go:31-65`:
  first-event 120s, stall 30m, heartbeat 60s) and classified by
  `internal/runner/failclass.go:24-53` (`no_first_output` / `stalled` / `timeout` /
  `unknown` with recovery hints). A hosted hang during round-01 is therefore already
  detected, classified, and surfaced — a separate pre-flight hosted PONG would buy the
  same answer at extra token cost.

So the per-idea roster cost collapses to **0 hosted round-trips by default**. A hosted
PONG becomes opt-in only, mirroring the existing `agents verify --full --yes` gate
(`internal/app/app.go:1916`, `runFullVerification`): `parley preflight --full --yes`
adds Tier-1; `parley run`'s auto path NEVER does Tier-1. The protocol §9.0 text should
state this explicitly: the roster ping is Tier-0 (CLI presence + version) by default;
hosted liveness is implicitly covered by the supervised round-01 dispatch.

Parallelism (contained perf win): `agents.Discover` currently loops sequentially
(`discover.go:261-286`), so worst case is N × 4s if several version probes stall. Since
`--version` is local it's usually fast, but the worst case is cheaply removed by
fanning the version probes out concurrently behind a single shared deadline (a
`sync.WaitGroup` + per-result slot, ~15 lines). Worst-case wall-clock then ≈ 4s, not
4×N×4s. Recommended; low complexity, real tail-latency win. No file TTL cache — reuse
the in-process discovery result (computed once at `runTask:1693`).

Hung vs slow (brief Q3): don't try to distinguish in the preflight. Tier-0's 4s
version-probe cap answers "is the CLI installed and does `--version` respond in 4s";
timeout = unavailable. The slow-vs-hung distinction for the hosted call is far better
answered by the runner's 120s/30m watchdog during real work than by any pre-flight
PONG. Keep preflight's job narrow.

Skip conditions / re-include (brief Q4): the exclusion record stays in `00-prompt.md`
(`excluded: [<id> — reason — confirmed <date>]`) per the locked decision. For
re-include detection, preflight scans the most recent ideas' `excluded:` lines for an
agent id that is now Tier-0 green — a directory scan over a small number of idea
prompts, no new state file. A derived `meta/exclusions.json` index is a v2 option if
the idea count grows; for v1, parsing the existing `excluded:` line keeps the single
source of truth where the locked decision put it and avoids a new artifact to keep in
sync. When `--participants` pins a subset, ping only that subset (`selectedParticipantIDs`
already does this at `app.go:1700`).

## C. `parley preflight` — smallest useful surface + non-deadlocking wiring

Smallest surface: one command, reusing the existing exit-code semantics
(`internal/app/app.go:290-294`):

```
parley preflight [--dir DIR] [--json] [--full] [--yes]
  default  = Tier-0 freshness + Tier-0 roster ping; print readiness report; no irreversible action
  --full --yes = add Tier-1 hosted PONG (mirrors `agents verify --full --yes`); off by default
```

Exit codes (no new semantics invented):
- `0` ready — no gate (freshness OK or source-advisory; all rostered agents Tier-0 green).
- `3` pending manual handoff — a gate requires confirmation (breaking freshness bump,
  an agent to exclude, or a re-include). Reuses the existing exit-3 "pending
  manual/interactive handoff" meaning already used by the consensus path. The report
  names the gate and prints the exact confirm command.
- `1` hard failure — no workspace, or resolving a gate would violate §1 non-solo / §5
  quorum (excluding leaves <2 participants → the §1 block the locked decision names).
- `2` usage error.

Output: a compact readiness block — freshness line (source/consumer, drift y/n,
additive/breaking), roster table (agent, installed, version, available y/n, reason),
and a `Gates:` section. Reuse `printSkillStatusSummary` + the `PrintRuntimeMatrix`
patterns already in the codebase — minimal new rendering. `--json` mirrors the same
fields for tooling.

Implementation shape: ONE `preflight(root, opts) (report, gates)` function. `parley
preflight` (standalone) and `parley run` (in-process step) both call it — no duplicated
logic. The standalone command is for operators/CI who want to check readiness without
starting a run.

Wiring into `parley run` without deadlocking unattended runs — the critical invariant:
`runTask` calls `preflight(...)` BEFORE `runcontrol.Create` (`app.go:1713`), reusing the
`discovered` slice already computed at `app.go:1693` (no second discovery pass).

- **No gates → proceed.** Added latency is just the one 10s freshness subprocess
  (skipped entirely under `--no-preflight`).
- **Gates non-empty — attended** (TTY present, or `--no-auto`): route the gate through
  the EXISTING HITL machinery (`runcontrol.StartAutoAnswerer` / `parley answer`,
  `app.go:1732`) — print the gate + the confirm command and STOP, letting the operator
  confirm via `parley answer` or re-run with the decision. This matches the existing
  `confirmLaunch` + `--yes` pattern (`app.go:1709-1712`); no new blocking-read code.
- **Gates non-empty — unattended** (`--auto` and no TTY / CI): the locked decision
  requires user confirmation for exclude/re-include, and there is no user present. So
  the rule is: **do NOT auto-exclude, do NOT read stdin, do NOT block.** Hard-stop with
  a clear non-zero exit and the message "roster/freshness gate unmet in unattended
  mode; re-run attended or pre-select with --participants." The operator re-runs
  attended. This preserves §1/§5 (never silently solo, never silent roster change) and
  can never hang a CI pipeline.
- **Auto-answerer boundary:** the two new gates (breaking-freshness, exclude, re-include)
  MUST be excluded from `StartAutoAnswerer`'s auto-answer set. They are
  non-auto-answerable by design. In unattended mode they hard-stop; in attended mode
  they surface as HITL questions. This single rule keeps the locked "user-confirmed"
  invariant intact in every mode and is the whole deadlock-safety argument in one
  sentence.
- **§1 interaction:** if excluding an unavailable agent leaves <2 participants, that is
  NOT a confirmable gate — it is a hard-stop exit 1 ("excluding <id> leaves solo; §1
  requires a user-authorized solo exception or a roster fix"), matching the existing
  `len(participants)==0` hard-stop at `app.go:1705-1708`.

Freshness auto-sync inside `parley run`: an additive consumer drift is the ONE write
preflight does without confirmation (it is the locked "auto additive" behavior), gated
by the source-role check (source = no write). It writes the single project
`COOPERATION.md` preserving zones + drops `meta/protocol-sync_<date>.md` + prints the
diff summary. A breaking bump never auto-adopts: attended → HITL gate; unattended →
hard-stop exit 3→1 path above.

Scope guard: preflight is a `parley run` step (idea creation is when "start of every
idea" fires). `resume` of an interrupted run does only the Tier-0 roster re-check (no
freshness, no gate) to warn if a participant CLI vanished mid-run — reusing discovery,
~0 added cost. Do not wire gates into `continue`/`resume`; the per-idea gate already
fired at idea start.

# Risks

- **Tier-0 can't see a hosted backend outage.** A CLI can be installed and `--version`
  fine while the hosted endpoint is down (auth/billing/rate-limit). By design we defer
  that to the supervised round-01 dispatch + `failclass.go`, which classifies it
  precisely. Trade-off: the first idea may launch a round that fails fast on a hosted
  error instead of being caught at preflight. Acceptable: the failure is classified and
  retryable, and the operator saves 4 hosted PONGs per idea in the common case. If a
  hosted pre-check is ever required, `--full --yes` is already there.
- **Semver-only additive/breaking can misclassify a sneaky minor.** Mitigated by
  zone-preservation + the sync record + diff summary (human-auditable), and the §7
  carve-out means a correction is a normal sync, not a protocol-change idea. A
  section-hash diff is the v2 hardening if this bites.
- **Source-role field is a new `meta/version.json` key** — adding it is a one-field
  schema bump, but it must be backfilled on existing decks (source repo sets
  `"source"`, consumers default to `"consumer"` when absent, with a one-time advisory
  print). Missing the backfill would make a consumer silently run advisory-only; the
  default-on-absence must be `"consumer"` (the safe, sync-capable direction), not
  `"source"`.
- **Re-include detection by scanning `00-prompt.md` `excluded:` lines** is O(ideas) and
  regex-over-prose-fragile. Fine at current idea counts; if the deck grows, a derived
  `meta/exclusions.json` index should replace it. Flagged as a known v1 ceiling.
- **`--no-preflight` escape hatch** trades safety for speed in CI. Document that it
  also skips the roster gate, so a CI run with `--no-preflight --auto` that hits an
  unavailable agent will fail mid-round-01 (supervised) rather than at preflight — the
  failure is still caught, just later. Acceptable, but the flag's help text must say so.
