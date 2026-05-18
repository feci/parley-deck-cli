---
agent: claude
idea: session-resume-cache-plan
review-round: 1
date: 2026-05-18
scope: slice-1
---

## Summary

Slice 1 lands the read-only foundation we sketched in round-02: a repo-local `run.json`, two `sessions` subcommands, and the existing `~/.parley-deck/sessions.json` left alone. The code is small, surgical, and additive — no orchestration semantics change, no schema change to the global session index, and the existing TUI workspace sessions console keeps working. The main issues are (a) the new `run.json` schema diverges from the round-02 proposal without being justified in `IMPLEMENTATION.md`, even though `schema_version: 1` is now a durable contract, and (b) test coverage stops at the happy path despite round-02 enumerating several explicit edge cases.

## Findings

### [MAJOR] `run.json` schema diverges from round-02 proposal without rationale

`internal/runmanifest/manifest.go:17-28` defines:

```
SchemaVersion, RunID, WorkspaceRoot, IdeaSlug, Task, Mode,
Transport, Participants, CreatedAt, UpdatedAt
```

Round-02 (`parley-deck/ideas/session-resume-cache-plan/round-02/claude.md`) explicitly proposed:

```
schema_version, run_id, idea, workspace, transport,
participants, started_at, status
```

Differences that now ship as `schema_version: 1`:

- Renames: `idea` → `idea_slug`, `workspace` → `workspace_root`, `started_at` → `created_at`.
- Additions: `task`, `mode`, `updated_at` (not in round-02 spec).
- Omission: **no `status` field**, even though round-02 explicitly carried `"running | completed | failed | abandoned"` and said it "may be updated on terminal transitions." Without it, the manifest cannot express "run finished" without a schema bump in a later slice.

These are defensible choices in isolation (Task/Mode are nice for `inspect`), but `schema_version: 1` becomes the canonical disk contract the moment a single run is created. Renames after the fact require a real migrator; additions/removals require either tolerance code or a `schema_version: 2`. Suggested fixes:

1. Add `status` (string, one of `running|completed|failed|abandoned`, default `running`) now, even if no writer flips it past `running` in slice 1. It costs one field and removes a forced schema bump later.
2. Either rename `idea_slug`/`workspace_root`/`created_at` back to the round-02 names, **or** explicitly log the deviation in `IMPLEMENTATION.md` under a "Deviations from round-02 proposal" section so reviewers and follow-up slices have a paper trail.
3. Decide whether `task`/`mode`/`updated_at` are slice-1 contract or implementation detail; if contract, document them.

### [MAJOR] Test coverage stops at the happy path

Round-02 enumerated concrete test cases for slice 1. Current state (`internal/app/app_test.go:481-548`, `internal/runcontrol/runcontrol_test.go:15-69`):

| Round-02 test case                                            | Covered? |
| ------------------------------------------------------------- | -------- |
| `sessions list` empty index                                   | no       |
| `sessions list` malformed index (treated as empty + warning)  | no       |
| `sessions list` populated (single session)                    | yes      |
| `sessions list --json` stable shape                           | no       |
| `sessions inspect` run with `run.json`                        | yes      |
| `sessions inspect` legacy run without `run.json`              | no       |
| `sessions inspect` run-id only resolvable via `--dir`         | no       |
| `sessions inspect` run-id not found anywhere (non-zero exit)  | no       |
| `sessions inspect` workspace path missing (degraded, no error)| no       |
| `sessions inspect` run-id ambiguous across workspaces         | no (yet code exists in `resolveIndexedSession:805-810`) |
| `run.json` writer: schema_version, RFC3339 round-trip         | partial (via runcontrol test only; no direct `runmanifest` unit test) |

The error-path branches in `resolveIndexedSession` (multiple matches, target not in index without `--dir`, target not in index with `--dir`-only fallback) and the `errors.Is(err, os.ErrNotExist)` legacy-manifest branch in `inspectSession` are unreached by tests. These are the exact paths a real user is most likely to hit first (no index entry yet, legacy run created before slice 1). Add the cases listed above; the test scaffolding in `app_test.go` is already set up for `PARLEY_HOME` and synthetic run dirs, so each test is small.

### [MINOR] `sessions inspect <run-id>` requires `--dir` for off-index fallback

Round-02 said: *"Exit non-zero only if `<run-id>` is not in the global index AND not found on disk by direct path lookup."*

`runSessionsInspect` (`internal/app/app.go:733-762`) falls back to `runstate.ResolveRun` only when `*root != ""`. Without `--dir` the user gets `run X is not in the local session index` even if they are currently `cd`'d into the workspace that owns the run. Two reasonable fixes:

- Default `*root` to `"."` (matching `parley resume`/`parley status`).
- Or document explicitly in `parley sessions inspect --help` and the top-level help block that `--dir` is required for legacy/unindexed runs.

This is a UX corner, not a correctness bug.

### [MINOR] `sessions list` output schema differs from round-02 proposal

Round-02 proposed one row per session: `run-id | idea | workspace | started-at | status`. The implementation (`internal/app/app.go:717-728`) prints `run-id  idea=  state=  last=` on one line and `workspace=` / `participants=` on indented continuation lines, where:

- `state` is `active|terminal` (binary), not `running|completed|failed|abandoned` from the proposal.
- `last` is `LastEventAt` (or fallback), not `started-at`.

The richer layout is fine in itself, but combined with the missing `status` in `run.json` it means slice 1 has no path to surface "this run failed" vs "this run completed" — both are just `terminal`. If you accept the `status` field fix above, also surface it here.

### [MINOR] Legacy-run message in `printSessionDetail` drops the "(legacy run)" qualifier

`internal/app/app.go:882` prints `Manifest: missing (%s)` with the expected path. Round-02 specified `manifest: missing (legacy run)`. The current message is honest (it shows where the manifest was expected) but loses the cue that this is the expected state for runs created before slice 1. A user seeing `Manifest: missing (/abs/path/run.json)` may file it as a bug. Suggest `Manifest: missing (legacy run; expected at %s)`.

### [MINOR] `--json` flag advertised as "unstable" while tests treat it as stable

`internal/app/app.go:686` and `:737` register the flag with help text `"print unstable JSON output"`. The committed test `TestSessionsCLIListAndInspect` (lines 528-547) asserts on `session.run_id`, `session.idea_slug`, `manifest.run_id`, and `manifest.mode` — i.e. it treats the shape as stable. Round-02 said the `--json` shape **is** stable for slice 1. Either drop the "unstable" word from the help text or weaken the assertions; the current mix sends mixed signals to downstream tooling.

### [MINOR] `--help` does not record the single-process / no-locking assumption

Round-02 explicitly said the single-process assumption should be documented in `--help`. The help block (`internal/app/app.go:134-137`) describes the `sessions` command but does not mention that concurrent `parley` processes against the same run are out of scope, and `runmanifest.Write` (lines 73-107) has no file lock. The behavior is fine for slice 1 — but the caveat needs to be visible to users.

### [NIT] `sessions inspect --dir` is undocumented in the top-level help

The `Usage:` block (`internal/app/app.go:101-102`) shows `sessions inspect [--dir DIR] [--json] RUN_ID`, but the longer `sessions` description (lines 134-137) doesn't mention what `--dir` is for. Worth one line: "Use `--dir` to inspect an unindexed run by direct workspace path."

### [NIT] `runmanifest.Path` does not normalize `root`

`internal/runmanifest/manifest.go:69-71` returns `filepath.Join(root, ...)` verbatim. If `root` is relative, the returned path is relative, but `New()` (lines 51-54) already stores the absolute path in `WorkspaceRoot`. Inconsistency between "what's in the JSON" and "what `Path()` returns" is a small footgun if anyone uses `Path()` to derive other paths later. Make `Path()` use `filepath.Abs(root)` (best-effort, same pattern as `New`) for symmetry.

### [NIT] `registerSession` error-swallowing is consistent with prior code, but worth a comment

`internal/runcontrol/runcontrol.go:142-157` returns silently on `sessionstore.Default()` error and ignores `Upsert` error. This matches `registerWorkspaceSessions` in `internal/app/app.go:1693-1719`, so it's not a regression. But silently dropping session-index writes means `sessions list` may report nothing after a successful run — a single `// best-effort; sessions index is an advisory cache` comment would prevent a future reviewer flagging this.

## Open questions

- Should slice 1 also persist a stub `~/.parley-deck/runs/<run-id>/` directory now (even empty), so slice 2's per-attempt cache can land additively without changing the global-index→local-cache layout? Round-02 deliberately deferred the per-run cache; if we keep that deferral, no change needed — but the layered model in claude/round-01 assumed it exists.
- The IMPLEMENTATION.md frontmatter does not yet record a `head-commit:` and the slice is currently uncommitted in the working tree. Per §4 Phase 5 the implementer fills this on completion. Worth confirming the commit boundary for slice-1 before this review is signed off, so the `reviewed-commit:` field of subsequent review-round files is unambiguous.
- Quorum: `gemini` and `hermes` round-01 are still missing per `inbox/codex-to-all_session-resume-cache-plan_agent-blockers.md`. Slice review proceeds under the user-authorized incremental implementation, but a recorded solo exception (or a fresh retry) should land before consensus/FINAL on the full plan.

## Residual risks (if all of the above stays as-is)

- Locking `schema_version: 1` with the current field set forces a real migration to add `status` later.
- No concurrency story means two simultaneous `parley run` invocations against the same `run-id` would have undefined ordering on `run.json` (last-writer-wins, no error). Acceptable for slice 1; must be addressed before slice 3 introduces resume.
- Legacy runs are degraded by code but not by test — regressions there will be silent.
- Workspace moves (user `mv`s the repo) will leave `sessions.json` pointing at the new path while `run.json`'s `workspace_root` still points at the old absolute path. Out of scope for slice 1 (`rebind` is slice 3), but document the limitation in `--help` for `sessions inspect`.
