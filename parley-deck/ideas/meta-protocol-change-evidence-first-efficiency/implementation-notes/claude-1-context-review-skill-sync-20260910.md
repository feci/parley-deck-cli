---
agent: claude-1
idea: meta-protocol-change-evidence-first-efficiency
date: 2026-09-10
reviewed-commit: f84560b
reviewed-commit-provenance: launcher-supplied worktree HEAD; SECONDARY, not independently confirmed (no shell)
artifact-kind: owner change record + partial independent review
not-a-signoff: true
not-a-phase-6-review: true
not-an-acceptance: true
owner-files-changed:
  - parley-deck/COOPERATION.md (§9 item 1)
  - internal/protocol/defaults/COOPERATION.md (§9 item 1, identical)
  - <skill-worktree>/skills/parley-deck/references/COOPERATION.md (§4.0.1 + §4 Phase-8 LE-7 mirror; §9 item 1)
  - <skill-worktree>/skills/parley-deck/SKILL.md (Required Protocol Context, 2 paragraphs)
  - parley-deck/meta/protocol-changelog.md (one new UNRELEASED entry, prepended)
excluded-as-owner: internal/protocolpacket/** — read for cross-owner interaction only, no independent verdicts on its internals
---

# Part 1 — Owner changes (Task A)

**A1. LE-7 mirrored into the skill reference copy.** Both edits from the live
`COOPERATION.md` now exist verbatim in the skill worktree's
`references/COOPERATION.md`: the §4.0.1 LE-7/LE-11 line ("can only withhold a close, never
establish one …") and the §4 Phase-8 *Close-decision integrity* paragraph replacing
"fail-open on its own error". That copy's own transport/roster/layout differences are
untouched — only those two passages changed, and the file keeps its own line offsets.
This supersedes the "not yet mirrored" status note in my previous changelog entry.

**A2. `refused` is now explicitly a stop, in all four owned instruction sources.**
§9 item 1 previously read "Without an attestation, or on `refused`, read all of
`parley-deck/COOPERATION.md` and record `context_mode=full-fallback`" — one clause covering
two different outcomes, readable as permission to downgrade a refusal into a fallback.
It now separates them, identically in the live deck protocol, the embedded default mirror
and the skill reference:

- `full-fallback` is unchanged and stays valid: read the live authority in full, record the
  mode with its reason, proceed.
- `refused` (unprovable authority, detected secret) is a **stop** — never emit the refused
  content, never substitute another authority (bundled snapshot, cached or stale copy,
  hand-assembled excerpt), never continue that launch on unattested text; resolve at the
  renderer and re-render, or report the blocker.
- A protocol task launch carrying **no** attestation is unresolved the same way and obtains
  one from the renderer before the task starts; reading the full live source applies only
  where no renderer is reachable.

`SKILL.md` carries the matching wording in both places that previously implied a missing
attestation is simply recorded as `full-fallback`.

This states existing FINAL D4 rules at the point of action. **No authority was broadened:**
no new context mode, `full` still the default, `packet` still the ratified trial's explicit
input, no applicability classification touched, and no heading edited (the change is prose
inside the existing `### 9.0 …` block, so every `meta/packet-applicability.yaml` locator
still resolves). Changelog entry added, marked UNRELEASED.

**Not done here, by instruction:** no add-on hash-manifest regeneration, no build, no
`go test`, no drift-guard run. No shell exists in this session; the facilitator owns those.

# Part 2 — Independent review of f84560b

Method: PRIMARY source reads with file:line locators. I executed nothing. The
facilitator-reported runner/app timings and package passes are SECONDARY testimony and
verify nothing below.

## Refutation attempts that failed (the implementation held)

- **"A caller can certify forged bytes."** `protocol_context.go:117-120` re-derives the
  context and overwrites `info.Context` unconditionally; a caller-supplied `Mode:"full"`
  cannot survive. `protocol_context_test.go:127-159` asserts this for command/exec/ACP with
  no child spawned.
- **"Some task boundary still spawns unattested."** `beginLaunch` has exactly three callers
  (`launch.go:237` probe, `protocol_context.go:121`, `handoff.go:54`), and every task
  `exec.Command*` sits downstream of `beginProtocolLaunch` (`launch.go:42`, `runner.go:1039`,
  `acp.go:37`). In `execAgentProcess` the prepared prompt shadows the raw one before both
  `buildAgentInvocation` (`:1055`) and stdin (`:1089`).
- **"RunMeasured double-wraps."** One chain only: `launch.go:201` → `:42`. `acp_test.go:198`
  asserts exactly one `<parley-protocol>` envelope reaches a real ACP child.
- **"An app caller overwrites the attested stdin."** The only remaining `.Stdin =` in
  `internal/app` is `consensus_request_signoffs.go:566` (`os.Stdin` on a raw interactive
  editor spawn — not a prompt channel). `runHeadlessSignoffAgent:441` now goes through
  `runner.CommandFor` and sets no stdin, closing the signoff asymmetry I filed earlier.
- **"A CLI user can select probe-only."** `agents exec --phase` flows to `RunMeasured` →
  `trackedCommandFor`, never `ProbeCommandFor`; the only `Phase:"preflight"/"runtime-probe"`
  setters are `preflight.go:846` and `app.go:2171`, both with CLI-constructed probe prompts.
- **"A refusal still writes a prompt."** `handoff.go:63-65` returns before `MkdirAllResilient`
  and before either file write.

**Earlier findings of mine now resolved (PRIMARY):** `parseGoalVerdict` is fail-closed —
sticky FAIL, any unknown marker blocks PASS (`driver_impl.go:422-452`), retiring my
mixed-verdict MAJOR; refusal reasons survive via `SafeLabel` (`protocol_context.go:84`);
the closing-envelope guard exists (`:92`); the phase map covers consensus/final/review/
round-N and `signoffContext:859-867` supplies idea+phase; the Darwin raw-fd lifetime NIT is
fixed with `SyscallConn().Control` (`sync_darwin.go:15-22`); and my Windows CRITICAL is
addressed by `replace_windows.go:14-26` (`MoveFileEx` + `MOVEFILE_WRITE_THROUGH` instead of
a directory handle).

## Findings

**[MAJOR] Non-darwin directory-sync barrier still has no unsupported-errno rescue.**
`replace_unix.go:22-27` requires `SyncFile(dir)` on the parent directory, and
`sync_other.go:8` is a bare `file.Sync()` rescuing nothing; darwin rescues only `ENOTTY`
(`sync_darwin.go:14`). Any `EINVAL`/`ENOTSUP`/`EOPNOTSUPP` from a directory fsync fails the
whole publication — for `cursor.go:101` that becomes `ActionEscalated`, for `handoff.go:159`
a refused handoff. `fsutil` documents itself as hardening for "virtio-fs, NFS, SMB", and this
repo runs on such a mount. Downgraded from my earlier CRITICAL because Windows is now fixed;
still an unvalidated hard dependency. Fix: enumerate an unsupported-errno set treated as
satisfied while every I/O error stays fatal, plus a fixture forcing that path.

**[MINOR] The `ProbeCommandFor` boundary is a label guard and is untested.** `launch.go:232`
gates on `info.Phase` only, not on prompt provenance; any future in-module caller setting
`Phase:"runtime-probe"` with a task prompt gets an unattested spawn labelled `probe-only`.
No test references `ProbeCommandFor` or `probe-only`. Suggest a probe-prompt constructor (or
a package-private entry point) plus a negative test on the guard.

**[MINOR] The envelope guard checks the closing tag only.** `protocol_context.go:92` rejects
`</parley-protocol>` but not a stray `<parley-protocol>`, which would leave a recipient's
first-marker scan ambiguous and would break `acp_test.go:198`'s exact-count assertion. One
extra clause.

**[MINOR] The spawn-TTY child runs outside launch evidence.**
`consensus_request_signoffs.go:564-569` spawns via raw `exec.CommandContext`, so the process
that actually consumes the attested prompt has no invocation record or procctl marker; only
the handoff packet write is recorded.

**[NIT]** Per-invocation `handoff-prompt-<id>.md` files are never pruned (good for audit,
unbounded on disk). **[NIT]** `writeHandoffPrompt:144-159` double-closes the staged file and
post-rename `os.Remove` always fails; both ignored. **[NIT]** The protocol is re-rendered per
launch, so agents in one round can carry different `source_sha256` if the source changes
mid-round — detectable per record, but nothing compares them at run level.

## On the `ProbeCommandFor` D4 interpretation (asked explicitly)

I concur with the disposition. D4 governs protocol *task* context; §9.0 defines the liveness
ping as a bounded round-trip through the agent's real configured invocation, and that ping
must work on a host with no deck — requiring one just to test PONG is a requirement neither
D4 nor §9.0 states, and it would break preflight on a fresh install. The integrity property
that matters holds: `probe-only` is a distinct mode with no hashes and an explicit
`no-protocol-task` reason, and an empty mode still records `unattested`
(`telemetry.go:60-61`), so nothing reads as attested that is not. My reservation is the
guard's shape, filed as the MINOR above — not the interpretation.

## Limitations

No execution of any kind: no build, no tests, no drift guard, no git. `f84560b` is the
launcher's word. Windows remains untested by anyone as far as I can see; the platform
replacement resolves the compile-and-design half of my CRITICAL, not the runtime half.
`internal/protocolpacket/**` is mine and excluded from independent verdicts. Nothing here
assesses AC-P1/P2, AC-E1/E2, AC-T*, AC-B*, AC-X* or the packet trial, and the full suite is
reported running, not passing. **This is not a Phase-6 signoff and not an acceptance.**
