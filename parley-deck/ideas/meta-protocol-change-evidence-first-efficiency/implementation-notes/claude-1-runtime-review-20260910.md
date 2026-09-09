---
agent: claude-1
idea: meta-protocol-change-evidence-first-efficiency
date: 2026-09-10
reviewed-commit: 8b11f1d
reviewed-commit-provenance: launcher-supplied git snapshot; SECONDARY, not independently confirmed (no shell in native restricted mode)
artifact-kind: partial-slice review handoff
not-a-signoff: true
not-a-phase-6-review: true
scope:
  - internal/fsutil/sync_darwin.go, sync_other.go
  - internal/driver/cursor.go (+ its reservation call site internal/driver/impl.go)
  - internal/runner/protocol_context.go, launch.go, handoff.go, telemetry.go
  - internal/app/driver_impl.go (GoalCheck), consensus_request_signoffs.go
excluded-as-owner:
  - internal/protocolpacket/** (claude-1-owned; read only for cross-owner API interaction, no verdicts on its internals)
unfinished-full-scope-gates:
  - AC-P1 (renderer wiring covers only manual agents exec + interactive handoff; round/ACP/consult/headless-signoff paths unwired)
  - AC-T1/T2/T3, AC-E1/E2, AC-B1/B2, AC-L1, AC-P2, AC-X1/X2, AC-H1/F1 — none assessed here
  - Kimi typed-criterion gate not integrated; textual PASS is not criterion evidence
  - No execution performed by this reviewer; facilitator-reported green results are treated as reported, not verified
---

# Partial-slice runtime review — Codex recovery changes

Method: PRIMARY source consultation of the files above at the current worktree
state. I ran no commands and executed no tests; every verdict below rests on
quoted code with a file:line locator. Facilitator-reported PASS results are
`SECONDARY` testimony and are not treated as verification of any claim here.

## Refutation attempts that failed (the implementation held)

- **"Full-fallback spuriously refuses the launch."** `protocol_context.go:59`
  refuses when `Hash(c.Body) != c.PacketSHA256`. I expected `ModeFullFallback`
  to leave `PacketSHA256` empty (it is `omitempty`) and hard-refuse every
  legitimate fallback. `packet.go:397` sets `ctx.PacketSHA256 = Hash(ctx.Body)`
  unconditionally on every non-refused path, so the check passes. **The
  attestation does name the bytes actually passed** (PRIMARY).
- **"Double `finish` duplicates a terminal record."** `launch.go:216-230`
  finishes on an `openPrivateLog` failure without setting `started`/`waited`,
  so `cleanup` (`launch.go:75-87`) calls `finish` a second time.
  `telemetry.go:103` wraps the body in `sync.Once`; the second call is inert
  (PRIMARY).
- **"`checker == o.implementer` can strand a real idea."**
  `driver_impl.go:381` is genuinely unreachable: the only way to get
  `drafter == implementer` is `len(reviewers) == 0` (`driver_impl.go:66-69`),
  and `ReviewRoundComplete` returns `len(o.reviewers) > 0`
  (`driver_impl.go:297`), so close is never reached. The CF6 comment's claim
  holds (PRIMARY).
- **"Refusal spawns anyway."** Both `RunMeasured` (`launch.go:196-208`) and
  `WriteHandoffPacket` (`handoff.go:45-59`) record a failed invocation and
  return before any `exec` and before any prompt file is written (PRIMARY).
- **"A post-rename dir-sync failure loses the charge."** `cursor.go:101-113`
  renames first, so a failure leaves the *conservative* (charged) state on
  disk, and `impl.go:320-323` escalates without running `Fixup` (PRIMARY).
- `LoadCursor` duplicate-key, `null`, negative-counter, trailing-JSON and
  symlink-swap rejection all hold as written (`cursor.go:144-196`) (PRIMARY).

## Findings

### [CRITICAL] Unconditional directory fsync makes the cursor unsavable on a release-gated platform

`cursor.go:106-113` requires `fsutil.SyncFile(dir)` on the parent directory and
returns an error otherwise; `sync_other.go:8` is a bare `file.Sync()` with no
fallback for any non-darwin platform. `impl.go:321-322` turns any `Save` error
into `ActionEscalated`, so a platform that rejects directory sync cannot
reserve a fix-up cycle at all — the driver is not degraded there, it is dead.

The implementer's own manifest concedes the exposure: "platforms that reject it
fail closed and remain unvalidated (including Windows)"
(`IMPLEMENTATION.md:266-267`). Windows is not out of scope. It is a
user-mandated release target — `inbox/user-to-all_parley-deck-cli-plan_priorities.md:14`
("all three including windows") and `:25` ("Release gate includes macOS, Linux,
and Windows") — with shipped binaries (`parley-v1.47.0-windows-x64.exe`,
`windows-arm64.exe`) and live Windows code (`internal/driver/proclive_windows.go`,
`internal/procctl/procctl_windows.go`). **Answering the question the launch task
asks: no, the stated caveat is not acceptable under this project's supported
scope.** An unvalidated hard dependency introduced into a release-gated platform
is a release blocker whichever way the platform call actually behaves.

The same shape applies on Linux/macOS over the network mounts this repo
explicitly targets: `fsutil.go:1-4` documents the package as hardening against
"virtio-fs, NFS, SMB", and this change adds an unconditional directory-fsync
requirement to a code path on exactly those mounts. `sync_darwin.go:14` only
rescues `ENOTTY`; `ENOTSUP`/`EINVAL`/`EOPNOTSUPP` still fail closed, and
`sync_other.go` rescues nothing.

Suggested fix: keep the file sync mandatory, and treat a directory sync that the
platform/filesystem *cannot support* (a distinct, enumerated errno set, or a
platform-specific no-op on Windows) as satisfied, while every I/O error stays
fatal. Add a fixture that forces the unsupported-directory-sync path.

### [MAJOR] Fail-closed goal check contradicts the protocol text in force, and is not recorded as a deviation

`driver_impl.go:375-416` now returns `false` for an unresolvable checker (`:386`),
a mkdir failure (`:391`), any process error or nonzero exit (`:406`) and an
inconclusive verdict (`:415`), and `impl.go:273-275` escalates on `!ok`.

The live `COOPERATION.md` §4 "Close-decision integrity (LE-7/LE-11)" — supplied
verbatim to this launch — states the opposite: "The goal-check is
defense-in-depth on top of the review consensus and **fail-open on its own
error (a broken or inconclusive checker never blocks a review-clean idea)**."
`impl.go:271` still carries the matching stale comment: "a checker error is
advisory (fail-open inside GoalCheck)", while the interface doc six lines away
(`impl.go:32-36`) was updated to say the reverse. `IMPLEMENTATION.md:132` records
"Deviations from FINAL.md: None" and no §7 protocol-text change accompanies it.

Concrete regression, reachable today: `TestStrictDesignOnlyCompletesAndRunsGoalCheck`
(`close_integrity_test.go:87-105`) shows LE-7 fires on `strict_gate: true` alone,
with one reviewer and no `auto_implement`. On such an idea, a reviewer whose CLI
is not resolvable at close time (`agents.ResolveParticipant` failure), or a
checker that exceeds the hard 2-minute consult deadline (`driver_impl.go:400`),
now escalates on every tick and the review-clean idea can never complete. That
is precisely the case the protocol reserved as fail-open.

I am not arguing fail-closed is the wrong end state — D3 plainly wants it. The
objection is that the code now ships against the ratified text without the §7
change, without a recorded deviation, and with a contradicting in-tree comment.
Fix: update `impl.go:271`, record the deviation, and carry the §4 LE-7 sentence
through the protocol change this idea already is.

### [MAJOR] `parseGoalVerdict` last-wins converts a stated FAIL into PASS on a trailing template echo

`driver_impl.go:432-445` resets the verdict on every matched line so the last
one wins. A checker that states its verdict and then reprints the output
template on its own line —

```
GOAL-CHECK: FAIL
...
Required format:
GOAL-CHECK: PASS
```

— yields `PASS`, because the template line survives the wrapper strip at `:428`
and `:439` and matches `rest == "PASS"` at `:441`. Every other tightening in
this change is fail-closed; this one is fail-open, in the single function that
decides an auto-close. The `rest == "PASS"` exact match (correctly rejecting
"PASS with reservations") does not help here.

Suggested fix: make the aggregation fail-closed rather than positional — any
matched `FAIL` line wins over any `PASS`, and more than one *distinct* verdict
is ambiguous. Add the trailing-template-echo case as a fixture.

### [MINOR] No delimiter guard on the `<parley-protocol>` envelope

`protocol_context.go:78` interpolates `c.Body` between literal
`<parley-protocol>` / `</parley-protocol>` markers with no check that the body
is free of the closing marker. The body is the protocol document itself, which
already documents its own launch attestation (§9 item 1), so a future source
edit that quotes the envelope silently truncates the protocol region a
recipient reads. Cheap fix: refuse (or use a per-launch nonce delimiter) when
`strings.Contains(c.Body, "</parley-protocol>")`.

### [MINOR] Renderer refusals lose their reason

`protocol_context.go:56-58` flattens every `ModeRefused` to `renderer-refused`,
discarding `c.FallbackReason` — which for a detected credential is
`secret-detected:<shape>` (`packet.go:252`), a shape class, not a secret. The
fallback path immediately below (`:74-77`) does preserve its reason through
`telemetry.SafeLabel`. An operator hitting the refusal gets no diagnosis.

### [MINOR] Phase mapping is partial, so most launches attest an unknown phase

`protocol_context.go:33-46` maps only `preflight`, `round-01`, `implementation`,
`review`, `review-consensus`, `fixup`. Consensus, final, `round-02`+, steer,
consult and signoff all fall to `phase = -1`, which `packet.go:267-269` records
as `unknown-phase:-1`. Full context is unaffected (correctly), but the shadow
record — the only per-launch data this boundary produces for the D5/AC-P2
packet work — is measuring an unknown phase for most real launches.

### [MINOR] Shadow block counts are surfaced into the agent prompt next to `packet_bytes: 0`

`protocol_context.go:69-78` embeds the whole `Shadow` verbatim in the launch
prompt. When a would-fallback reason is present, `packet.go:381-394` still fills
`included_blocks`/`omitted_blocks` from a classification pass that was **not**
applied, while `packet_bytes` stays `0`. This launch's own header reads
`{"would_fallback_reason":"unknown-track:unknown","packet_bytes":0,
"included_blocks":37,"omitted_blocks":32}` — readable as "32 blocks were
omitted from your context" when nothing was omitted. *Ownership note (§15.1): I
own the `Shadow` field semantics in `protocolpacket`, so this is an owner
observation, not an independent verdict; the verdict-eligible part is the
runner's decision to surface it unlabelled to the agent.*

### [MINOR] New within-command asymmetry on the signoff path

`consensus_request_signoffs.go:463-466` and `:514-517` correctly propagate the
`WriteHandoffPacket` refusal unchanged (the claimed behaviour holds, PRIMARY).
But `runHeadlessSignoffAgent` (`:436-460`) still goes through `runner.CommandFor`
with no protocol context, so within one `request-signoffs` invocation a manual
participant signs with the full attested protocol and a headless participant
signs with none, on the same gate. Telemetry records this honestly as
`Context.Mode = "unattested"` (`telemetry.go:60-62`) — the gap is disclosed, not
misreported — but the asymmetry is new and should be closed with the remaining
launch paths, not after them.

### [NIT] Non-atomic, unsynced handoff prompt write

`handoff.go:71` writes the attested prompt with a plain `os.WriteFile` to a
stable path that each attempt overwrites, with no sync and no post-write hash
check — in the same change set that gives `Cursor.Save` staging + file sync +
atomic rename + directory sync. `fsutil.WriteFileAtomic` already exists.

### [NIT] Unreachable unknown-field switch in `LoadCursor`

`cursor.go:176-182` rejects unrecognized keys, but `cursor.go:169`
(`strict.DisallowUnknownFields()`) has already failed the decode for exactly
that input, so the `default:` branch cannot be reached.

### [NIT] `syscall.Fsync(int(file.Fd()))` drops the `*os.File` lifetime guard

`sync_darwin.go:15` uses a raw descriptor. Safe in the current single-goroutine
call sites, but `SyncFile` is now a shared helper across owners; `f.SyncfdControl`
via `RawConn.Control` keeps the file alive for the call.

## Open questions

1. Is any CI job building or testing `GOOS=windows` today? If not, CRITICAL-1
   will not be caught before the release gate.
2. Does the `strict_gate` fail-closed goal check reach existing decks, or only
   ideas opened after this change? That decides whether MAJOR-1 is a migration
   hazard or only a forward-looking one.
