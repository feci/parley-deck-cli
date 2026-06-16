---
idea: rho-retro-tooling
review-cycle: 1
drafted-by: claude
date: 2026-06-16
reviewed-commit: 984c757
outstanding_agreed_fixes: 6
---

## Agreed fixes

hermes ACCEPT (no findings). codex + agy ACCEPT-WITH-FIXES; their findings are
agreed (several overlap and are merged):

1. **[MAJOR] `propose` write-boundary hardening (codex MAJOR-1 + agy MINOR + agy
   slug NIT).** `propose` must fail closed if ANYTHING already exists at
   `ideas/<slug>` (not just `00-prompt.md`), and must not be coercible to write
   outside the new idea dir. Fixes: validate the slug as strict kebab-case;
   `Lstat` `ideas/<slug>` and fail if it exists (covers a symlinked entry);
   create the final dir with `os.Mkdir` (not `MkdirAll`); write `00-prompt.md`
   with `os.OpenFile(O_CREATE|O_EXCL|O_WRONLY)`. Tests: existing slug dir without
   `00-prompt.md`, symlinked slug entry, non-kebab slug.
2. **[MAJOR] design-churn classification gap (agy MAJOR-A).** `classify` uses
   `s.Rounds > 2` while `score` adds friction at `s.Rounds > 1`; a 2-design-round
   idea scores >0 but buckets as `low-friction` and is then dropped by `Select`.
   Fix: `s.Rounds > 2` → `s.Rounds > 1` in `classify`.
3. **[MAJOR] blocker detection misses review-file verdicts (agy MAJOR-B).**
   `reBlocker` only matches `Status: ❌` (consensus signoffs); reviewer files use
   `Verdict: BLOCK`. Fix: also match `Verdict:\s*(?:❌|BLOCK)`.
4. **[MAJOR] missing D4 failure signals (codex MAJOR-2).** D4 names
   blocked/**abandoned** work, **drift-guard failures**, and **watchdog /
   `agent.failed`** events; the scanner reads none of the run-side ones. Fix:
   parse `status: abandoned` from IMPLEMENTATION.md/00-prompt.md frontmatter, and
   scan structured `parley-deck/runs/*/events.jsonl` (structured logs, NOT raw
   session transcripts) for `agent.failed` / `agent.no_first_output` /
   `agent.stalled` / `driver.error`, attributed to the idea via the run's
   `run.created` idea slug. Extend scoring + tests.
5. **[MINOR] generated `author:` hard-codes claude (codex MINOR).** Phase-4
   ownership depends on `author:`; a user or another agent may run `retro propose`.
   Fix: write a neutral `author: <fill: author>` placeholder (matching the
   existing `created: <fill: date>` placeholder).
6. **[NIT] test helper `itoa` overflows past 9 (agy NIT).** Replace with
   `strconv.Itoa`.

## Deferred follow-ups

- DPP/embeddings, raw-JSONL ingestion, live re-rollout, best-of-N, auto-apply,
  persistent quarantine registry — per the parent FINAL (cut from v1). Reading
  `runs/*/events.jsonl` (fix 4) is structured-log evidence, explicitly distinct
  from raw session transcripts, and is in v1 scope per D4.

## Dismissed findings

- None.

## Signoffs

<!-- Each participant APPENDS their own signoff block. Do NOT edit others' blocks. -->

### Signoff: claude — 2026-06-16
Status: ✅ ACCEPT
Notes: I drafted this; six agreed fixes (codex+agy, deduped). All are real; I will apply them in fix-up cycle 1. Fix 4 reads structured run event logs (in-scope per D4), not raw transcripts.

### Signoff: codex — 2026-06-16
Status: ✅ ACCEPT
Notes: Consensus preserves my three required fixes: propose hardening, D4 signal coverage, and neutral generated authorship.

### Signoff: hermes — 2026-06-16
Status: ✅ ACCEPT
Notes: Consensus correctly records my ACCEPT with no findings while capturing the six fixes from codex+agy.

### Signoff: agy — 2026-06-16
Status: ✅ ACCEPT
Notes: The consensus accurately merges and reflects my findings, including the six agreed fixes.

