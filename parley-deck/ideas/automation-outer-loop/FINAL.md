---
idea: automation-outer-loop
status: final
drafter: claude-1
created: 2026-06-24
participants: [claude-1, codex-1, hermes-1, antigravity-1]
---

## Decision

Implement Tier 4 (the outer loop) as the human-braked discovery layer described in the
research FINAL.md. Two deliverables, both bound by a single new invariant.

## LE-8 — Human-brake invariant (protocol §14)

Add `## 14. Automated outer loop (loop engineering) — the human brake` to BOTH
`COOPERATION.md` copies (live `parley-deck/COOPERATION.md` + embedded
`internal/protocol/defaults/COOPERATION.md`; the drift guard requires byte-identical text).
The section states, normatively:

- An **automated / standing / scheduled loop** (anything not driven by a human in the
  session — cron, CI, an MCP trigger, `parley loop tick`) may **discover candidates and
  draft Phase 0/1 prompts only**.
- It must **never**, without a recorded human gate OR a full-quorum consensus gate:
  push a candidate to quorum (staff `participants:` / flip `status: candidate` → `round-01`),
  start a deliberation / `parley run`, implement, land / merge / push, finalize (write
  `FINAL.md` / close), modify the roster, or override / bypass / reopen consensus.
- A loop-created idea is therefore always a non-active **`status: candidate`** with no
  `participants:` claim and a `## Promotion` note (consistent with the §12.11 monitoring
  watcher). Promotion to `round-01` is a human or manifest action.
- The brake is **fail-safe**: when in doubt, the automated loop does less (draft a candidate,
  escalate) — never more.

## LE-9 — `parley loop tick`

New `internal/loop` package + `parley loop tick` command. One-shot, scheduler-friendly,
**not a daemon**.

### `internal/loop`

- `type Candidate struct { Source, ID, Title, Detail, Fingerprint string }` — a discovered
  signal. `Source` ∈ {commit, ci, issue, manual}. `Fingerprint` dedupes; when empty it
  defaults to a stable hash of `Source + ":" + ID`.
- `type Config struct { Enabled bool }` — read from `parley-deck/loop/config.json`
  (`{"enabled": false}` is the default when the file is absent). Disabled-by-default.
- `type TickResult struct { Enabled bool; Created, Skipped []string; ... }`.
- `func ReadSignals(path string) ([]Candidate, error)` — reads a JSON array of signals from
  the given path (default `parley-deck/loop/signals.json`). Missing file → empty slice, no error.
- `func Tick(deck string, cfg Config, signals []Candidate, now time.Time) (TickResult, error)`:
  1. If `!cfg.Enabled` → return `{Enabled:false}` immediately, write nothing.
  2. For each signal, compute its candidate idea slug `loop-<source>-<fingerprint8>`.
  3. **Dedupe**: if an idea dir with that slug already exists, skip it (record in `Skipped`).
  4. Else write `ideas/<slug>/00-prompt.md` as a `status: candidate` prompt (no `participants:`,
     a `## Promotion` note, provenance: source/id/detail/fingerprint). Record in `Created`.
  5. Never staff a quorum, never run, never push. Return the result.

### `parley loop tick` (in `internal/app`)

- Usage: `parley loop tick [--dir DIR] [--signals PATH] [--enable] [--json]`.
- Loads `loop/config.json`; `--enable` forces `Enabled=true` for this one-off run (a human
  explicitly running the tick is the human gate — but it STILL only drafts candidates).
- Reads signals (`--signals` or default path), calls `loop.Tick`, prints a human or `--json`
  summary: enabled?, N created (slugs), N skipped (dups).
- Returns 0 on success even when disabled or when zero candidates (idempotent / cron-safe).
- It does **not** call `parley run`, push, merge, or finalize — LE-8 enforced in code.

## Tests

- `internal/loop/loop_test.go`: disabled-by-default writes nothing; enabled writes a
  `status: candidate` prompt (assert `status: candidate`, no `participants:`, `## Promotion`);
  dedupe (second tick over the same signal skips); fingerprint default; `ReadSignals` missing
  file → empty.
- `internal/app/loop_cmd_test.go`: `loop tick` disabled exits 0 + says disabled; `--enable`
  creates candidates; `--json` shape; never writes a non-candidate idea.

## Verification

`gofmt`, `go build ./...`, `go vet`, `go test -count=1 ./...` (incl. drift guard) green.

## Out of scope (follow-ups)

- Live connectors (GitHub/CI APIs) — MVP reads a signals file.
- Optional human-confirmed `parley run` from a promoted candidate (kept fully manual here).
- LE-12 durable goal state (contested, separate idea).
