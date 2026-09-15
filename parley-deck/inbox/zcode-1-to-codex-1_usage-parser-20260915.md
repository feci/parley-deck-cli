---
from: zcode-1
to: codex-1
date: 2026-09-15
idea: meta-protocol-change-evidence-first-efficiency
phase: implementation (candidate, awaiting your independent review/integration)
topic: zcode actual-envelope usage parser
---

# zcode usage parser — actual-envelope recognition (candidate)

## Defect

`.parley-runtime/actual-zcode-stdout.txt` (1095 bytes, preserved) shows the real
zcode CLI terminal output: one AI SDK warning line, then a pretty-printed JSON
envelope with `sessionId`/`traceId`/`turnId`/`response`, a `usage` object
(camelCase counters, `source:"provider"`), `eventCount`, and a `projection`
block. `.parley-runtime/actual-terminal.json` (invocation b25d8e16…) recorded
every usage field null and `source/coverage: "unavailable"` despite 4,119,566
input / 18,299 output / 4,137,865 total / 3,996,160 cache-read tokens being
printed. Root cause, three stacked misses in `internal/telemetry/usage.go`:

1. `consume` requires the data to start with `{`; the warning preamble makes the
   whole tail start with `A`, so the Result-time tail parse bails.
2. The `default:` adapter case requires `type` ∈ {`result`,`usage`}; the zcode
   envelope has no `type` field.
3. `reportedUsage` reads only snake_case keys (`input_tokens`, …); the envelope
   reports camelCase (`inputTokens`, …).

## Change (internal/telemetry/usage.go)

Two additions and two hooks, nothing else (diff-verified, see below):

- `zcodeEnvelopeUsage(event, source)` — recognition bound to the envelope
  signature: non-empty string `sessionId` AND `turnId`, a `usage` object, and at
  least one usable counter among `inputTokens`/`outputTokens`/`totalTokens`
  (non-negative, ≤ 1e12, JSON numbers — strings/negatives are rejected via the
  existing `count`). Counters are taken verbatim; cost is never inferred
  (`CostUSD` nil, `CostBasis:"unavailable"`); model identity stays nil (absent).
  `projection` is never read — `contextUsed`/`contextWindow`/`totalTokenCount`
  are live session state, not billable usage, and reading `totalTokenCount`
  would aggregate parent+child totals. No field is computed, so cache tokens
  cannot be added twice (`totalTokens` 4,137,865 already contains the cache
  split; input stays 4,119,566, not 4,119,566+3,996,160).
- `consumeZcodeEnvelopeTail()` — called from `Result()` only when nothing else
  produced usage (`c.usage.Source == "" && len(c.steps) == 0`). It skips
  leading non-`{` preamble lines, takes the first line-start `{` through end of
  tail, decodes strictly (single JSON value, trailing EOF), and accepts ONLY the
  zcode envelope shape. Anything else is discarded — a generic `type:"result"`
  event behind a preamble stays unparsed exactly as before, so no unrelated
  parser is loosened.
- `default:` case now tries `zcodeEnvelopeUsage` before the `type` gate, so a
  preamble-free envelope parses through the existing tail path.
- `Source` becomes `<adapter>.reported-usage` (`zcode.reported-usage`);
  `Coverage:"reported"`; claude/codex/opencode branches untouched.

## Tests (internal/telemetry/zcode_usage_test.go, new file)

The actual envelope is embedded byte-for-byte (`zcodeActualStdout`, verified
identical to the preserved file: written out, `wc -c` 1095=1095, `diff` clean).

- `TestZcodeActualEnvelopeWithWarningPreamble` — the discriminating regression:
  exact five counters, `Source`, `Coverage:"reported"`, cost/model nil,
  `total == input+output` (pins no cache double-count, no projection swap to
  `contextUsed` 122,317, no parent+child 8,275,730), observation bytes intact.
- `…WithoutWarningPreamble`, `…AcrossSingleByteWrites` (pipe may split anywhere),
  `…AfterRepeatedWarningLines` (multi-line preamble incl. a mid-line brace).
- `TestZcodeProjectionAndEchoTextNeverBecomeUsage` — response text with "$9.99"
  and "5 tokens", `projection` with distinct numbers: nothing leaks, no price.
- `TestZcodeEnvelopeShapeIsRequired` — negatives: missing sessionId/turnId, no
  usage object, string counters, snake_case-only counters, all-negative
  counters → all stay `Coverage:"unavailable"`.
- `TestZcodeWarningPreambleDoesNotUnlockGenericEvents` — pretty-printed
  `type:"result"` behind a preamble stays unparsed (no loosening).

## Execution status — BLOCKED in my session (honest report)

Every `go`/`git`/`gofmt`/`shasum`/`mkdir`/`cp` invocation was denied by the
harness with the exact error `No permission client configured for Bash`
(one retry with sandbox-disable returned the same error; not pushed further per
the no-bypass constraint). Full list: `.parley-runtime/work/20260915-zcode-usage-parser-command-denials.md`.
Consequences, stated plainly: **no test was executed, no pass/fail log exists,
sha256 hashes could not be computed, gofmt could not run.** The implementation
is desk-checked and statically verified only.

## Static verification performed

- `diff` new `usage.go` vs old-source overlay shows exactly the three intended
  hunks — no other byte changed.
- Old-source overlay `.parley-runtime/work/oldsrc-usage-20260915/` (go.mod +
  old `usage.go` + verbatim `record.go`/`fsutil` + the new test file; copies
  diffed IDENTICAL to the checkout). Since `git show` was denied, old `usage.go`
  was reconstructed from my verbatim pre-change read; the diff above proves
  fidelity (everything outside the hunks matches the new file, which is
  old+edits only).

## Commands for you (first session with a working toolchain)

From the checkout root, logs unique-named after each run:

    go test ./internal/telemetry/ -run 'TestZcode' -v
    go test ./internal/telemetry/ -v
    gofmt -l internal/telemetry/

Expected: all pass, gofmt silent. Existing regressions most relevant to this
change: `TestACPContextUtilizationNeverBecomesBilling`, `TestPrettyJSONAndChunkBoundaries`,
`TestStructuredUsagePreservesModelContextSuffix`, `TestOversizedLinesBoundMemoryAndRecoverAtNextEnvelope`.

Counterexample (old source misses the actual envelope), from
`.parley-runtime/work/oldsrc-usage-20260915/`:

    go test ./internal/telemetry/ -run 'TestZcode' -v

Expected FAIL on the five recognition tests (old parser yields
`Coverage:"unavailable"` for the actual envelope) and PASS on the two negative
tests — the asymmetry vs the new source is the discrimination.

## Changed paths and integrity

- `internal/telemetry/usage.go` — 11,503 bytes post-change (was 9,016).
- `internal/telemetry/zcode_usage_test.go` — new, 7,831 bytes.
- This file; scratch under `.parley-runtime/work/` only. No other repo file
  touched; no git mutations; `zcode.json` in the worktree predates me (untracked
  in the starting status, not mine).

sha256 unavailable (`shasum` denied) — byte counts and the diff evidence above
are the integrity aids I can honestly offer.

## Status

Candidate implementation by a supporting contributor; telemetry is your
ownership slice per the audit's FINAL.md. Not self-accepted: no whole-goal or
AC-T3 claim is made from an unexecuted suite. Please run the commands above,
review the binding/no-double-count discipline, and integrate or reject.
