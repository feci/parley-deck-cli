---
agent: claude
idea: roadmap-implementation-plan
review_round: 1
date: 2026-05-17
implementation-pr: https://github.com/feci/parley-deck-cli/pull/20
verdict: COMMENT
---

## Summary

The slice cleanly implements the two reviewable behaviors required by `FINAL.md`:
a context-only sanitizer (`SanitizeForContext`) and a deterministic, runner-owned
`_index.md` writer wired into `RunRoundOne`. Determinism, H2-only extraction,
the `(sanitized_bytes + 3) / 4` token heuristic, the index-failure-as-warning
path, and source-artifact non-mutation are all covered by tests. The deferral
of repo-map, secret redaction, and other later slices is honored.

No CRITICAL or MAJOR defects found. The notes below are MINOR/NIT improvements
that should be considered before the slice lands, but none block this PR.

## Findings

### MINOR — Unclosed opening fence silently drops the entire tail of the artifact

`removeTaggedBlocks` (internal/runner/round_index.go:21-40) returns immediately
when an opening fence (`<think>`, `<thought>`, or `<thinking>`) has no closing
counterpart, throwing away every byte that follows the opening tag. The behavior
is asserted by `TestSanitizeForContextDropsMalformedOpenFence`, so it appears
intentional, but it is a silent, lossy operation on malformed input that:

- contradicts the "degrade gracefully" spirit in `FINAL.md` for the index path;
- can swallow large amounts of legitimate content if a participant emits a stray
  `<think>` literal (e.g. inside a code block, or a partial stream).

Suggested fix: when no closing fence is found, either treat the open tag as
literal text (write `remaining` unchanged) or strip only the opening tag and
keep the trailing content. A unit test for "stray closing fence with no opener"
and for "fence literal inside a fenced code block" would also be worth adding.
Document the chosen behavior explicitly in the `SanitizeForContext` doc-comment.

### MINOR — Sanitizer index math assumes ASCII-only fold

`removeTaggedBlocks` calls `strings.ToLower(remaining)` to locate fences but
then indexes back into the original `remaining` (round_index.go:25, 30, 32).
`strings.ToLower` can change byte length for some Unicode runes (e.g. the
Turkish dotted I), which would shift indices and could either drop unrelated
content or panic on a slice bound. The fence tags are ASCII, so a cheaper and
safer approach is a manual case-insensitive scan over `remaining` directly, or
asserting `len(strings.ToLower(remaining)) == len(remaining)` once. Probability
of hitting this in practice is low, but it is a latent correctness bug.

### MINOR — `extractH2Sections` `### ` guard is dead code

`internal/runner/round_index.go` in `extractH2Sections`:

```
if strings.HasPrefix(trimmed, "## ") && !strings.HasPrefix(trimmed, "### ") {
```

`"### foo"` does not satisfy `HasPrefix("## ")` (third byte is `#`, not space),
so the second clause never fires. Either drop the redundant guard or, if the
intent was to also exclude H3+ headings (e.g. `##foo` is already excluded by
`## ` prefix, and `## ## bar` is genuinely an H2 with `## bar` as title), make
the intent explicit. Today it is misleading.

### MINOR — Summary line search treats `---` and `#` lines as opaque

The summary scan skips any line beginning with `#` or `---`. That correctly
ignores subsequent headings but also skips legitimate horizontal-rule-prefixed
or hash-prefixed content (e.g. a paragraph beginning with `#123` or a Markdown
HR). A section whose first non-empty content is one of these renders as
`no summary text` in the index. Not blocking, but worth a comment in the helper
documenting the intent so reviewers don't "fix" it later.

### MINOR — Token heuristic location vs. `FINAL.md` wording

`FINAL.md` says "Include a deterministic approximate token estimate in index
metadata using a simple documented heuristic." The estimate is included per
participant in the body table and per-section blocks, and the formula is
documented in the body prelude, but the YAML frontmatter does not carry the
heuristic identifier. Either reword `IMPLEMENTATION.md` to reflect that the
heuristic lives in the body, or add a `token-heuristic: bytes_div_4` key to the
frontmatter so future consumers can detect formula drift.

### MINOR — `"index"` is a reserved-feeling pseudo-agent ID

The warning result is appended with `AgentID: "index"` (runner.go:138). If any
real agent is ever registered with that ID, the CLI row in `printRunResults`
will collide and a real participant could also be classified as the index
writer. Cheap mitigation: use an unmistakable sentinel such as
`_runner/index` or expose a dedicated `Result.Kind` field; alternatively, add a
defensive check in agent discovery that rejects reserved IDs.

### MINOR — `writeRoundIndex` does not `MkdirAll` the round directory

If `RunRoundOne` is invoked with zero discovered/available agents (so no
participant ever creates `round-NN/`), `os.WriteFile` fails because the parent
directory does not exist. The fallback warning path then attaches a result with
`OutputPath` pointing inside a non-existent directory. Today this is masked
because typical runs always have at least one agent that creates the round dir,
but a defensive `os.MkdirAll(roundDir, 0o755)` before `WriteFile` would make
the warning path describe an actual on-disk artifact and avoid future
regressions.

### NIT — `escapeTable` collapses empty cells to literal `"none"`

`escapeTable("")` returns `"none"`, so an empty artifact name (e.g. when a
result has no `OutputPath`) shows the same string as a deliberate `(none)`
sentinel computed in `buildRoundIndexEntry`. This is harmless today but creates
two sources of truth for "empty" cells. Suggest using a single sentinel
(`"(none)"`) consistently or letting the table render an empty cell.

### NIT — `len([]byte(value))` in `approxTokens`

Go strings already report byte length via `len(value)`; the explicit
`[]byte(value)` conversion allocates a copy on every call. Replace with
`len(value)`.

### NIT — Sanitizer's trailing `TrimSpace` is undocumented behavior

`SanitizeForContext` trims surrounding whitespace (round_index.go:18). For the
current consumer (index summary extraction) that is fine, but a function named
`SanitizeForContext` that quietly mutates leading/trailing whitespace is a
foot-gun for the later `context-pack-wiring` slice that will reuse it for
prompt assembly. Document the trim, or move the trim into the index builder so
the helper stays a faithful "fence remover".

## Tests / verification reviewed

- Sanitizer:
  - `TestSanitizeForContextRemovesSupportedReasoningFences` — supported fence
    set (think/thought/thinking) including verifying inter-fence content
    survival. Good.
  - `TestSanitizeForContextDropsMalformedOpenFence` — locks in current
    drop-tail behavior; see MINOR above for the suggestion to revisit.
  - Missing: no explicit "content with no fences is byte-identical to input
    (modulo trim)" test, no nested or sequential fence test, no case-mixed
    fence test (`<Think>`), no fence-inside-code-block test.
- Index builder:
  - `TestBuildRoundIndexIsDeterministicAndExtractsH2Only` — covers determinism
    by double-call equality, H3 exclusion, sanitization of `<think>`, frontmatter
    pass-through, and `before == after` source non-mutation. Strong test.
  - `TestBuildRoundIndexIncludesSkippedWithoutRecognizedSections` — exercises
    the skipped row and "no recognized H2 sections" degrade path.
  - Missing: a "golden" byte-for-byte snapshot would be the most direct way to
    enforce determinism across refactors (the current test only proves two
    in-process calls return the same string).
- Runner wiring:
  - `TestRunRoundOneCreatesArtifactWithHeadlessAgent` — verifies `_index.md`
    exists, contains the `ok` row, and `round.index_written` is emitted.
  - `TestRunRoundOneRecordsAgentFailure` — verifies the failed row appears.
    Switching from positional event lookup to a search-by-type loop is a good
    defensive change.
  - `TestRunRoundOneIndexWriteFailureIsWarning` — clever use of pre-creating
    `_index.md` as a directory to force a write error; verifies the warning
    Result is appended and `round.completed` is still emitted last.
- Coverage matrix against `FINAL.md` "Tests Expected": all five bullets are
  represented (sanitizer behavior, determinism, skipped/failed participants,
  non-aborting on index failure, source-artifact non-mutation).

## Risks / open questions

1. **Forward compatibility of `Result.Warning` and the pseudo-agent row.** The
   current pattern of inflating `Result` with a synthetic `AgentID: "index"`
   row blurs the meaning of `[]Result`. As later slices (hooks, error
   classifier) extend this struct, a dedicated `RunnerEvent`/`RunnerWarning`
   type or a `Result.Kind` enum would scale better. Worth flagging now even
   though it is not a blocker for this slice.

2. **Determinism guarantees under sort instability.** `sort.SliceStable` by
   `AgentID` is sufficient today, but `RunRoundOne` allows duplicate AgentIDs
   in principle (nothing currently rejects them). If two results ever share an
   ID, body section ordering depends on input ordering — fine for stable sort,
   but write a regression test if this is later relied upon.

3. **Heuristic stability across `SanitizeForContext` changes.** The approx
   token estimate is computed on the sanitized body. Any future change to the
   sanitizer (e.g. adding more fence types, fixing the malformed-tail behavior
   above) will change historical index outputs and break golden tests in
   downstream slices. Pinning the token-heuristic identifier in frontmatter
   (see MINOR above) would let consumers detect this drift.

4. **`round.index_failed` is emitted before `round.completed`.** That ordering
   is correct, but downstream consumers that build state machines from the
   event log should know the round terminal event can be preceded by an
   index-related failure. Worth documenting in the events catalog when it
   exists.

5. **Concurrency.** `writeRoundIndex` is called serially after participant
   completion, so no race today. If the runner ever streams participants and
   writes a partial index incrementally, this implementation will need
   revisiting.
