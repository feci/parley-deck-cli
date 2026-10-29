---
idea: meta-protocol-change-credit-exhaustion-failover
review-cycle: 1
drafted-by: codex-1
date: 2026-10-29
reviewed-commit: 931550cb4c57f041f9a8a1d1337571cb271a873a
---

Codex-1 is the organizer, implementer, participant, and consensus drafter. Signoffs ratify this proposed fix list; they do not close implementation.
Skill target: `70d29c1aebd9a632f6f6239123fd2f52e3775a48`.
Both independent native processes exited 0. Their receipt hashes match their canonical review files.

## Agreed fixes

1. **F1 — Preserve transient behavior before an implementer pin.** Source: Kimi's first MAJOR in `round-01/kimi-1.md`.
   Separate succession assignment from protected-role status. Keep fallback and global-default transient steps on the original two-attempt path.
   Preserve ordinary exclusion and global-default fallthrough before a pin. Do not synthesize exhaustion succession from an ordinary fallback failure.
   Resolve and bind the initial fallback/default identity when exhaustion requires a transition.
   Preserve explicit designee, pin, organizer, facilitator, and started-drafter gates.
   Test readiness, ordinary settlement, exhaustion succession, and mixed simultaneous losses. Test immutable history validation, not only `Roles.Protected`.
2. **F2 — Match CLI version metadata.** Source: Kimi's second MAJOR in `round-01/kimi-1.md`.
   Set root `VERSION` to `1.56.0`. Run the version check and complete Go suite at the fixed source.
3. **F3 — Preserve observed implementation presence on read errors.** Source: Kimi's third MAJOR in `round-01/kimi-1.md`.
   Apply `ExecutionStarted`'s boolean only when its error is nil. Preserve positive stat evidence and return the integrity error.
   Retain early-handover behavior and the existing unreadable-file regression test.
4. **F4 — Make the unsafe-mode fixture independent of umask.** Source: both reviewers' MINOR capture findings.
   Force the fixture's `0644` mode with `os.Chmod`. Check setup errors in that test.
   Run it under umasks `0077` and `0022`. Product storage must continue to reject unsafe modes.
5. **F5 — Correct the wait message.** Source: Zcode's first NIT in `round-01/zcode-1.md`.
   Replace initial-file-only wording with active-execution-record wording. Cover absent records and handovers before implementation starts.
   Keep the digest API shape unchanged; this correction needs no new path field.
6. **F6 — Remove the empty changelog heading.** Source: Zcode's second NIT in `round-01/zcode-1.md`.
   Remove the empty `Unreleased` section above 1.56.0. Preserve all release entries.

These fixes preserve FINAL D1-D12. They do not alter the signed protocol hunks or the frozen owner order.
The implementer must update IMPLEMENTATION and record current-source test results. Both reviewers must inspect the fixed complete diff.

## Deferred follow-ups

- `quota-zcode-native-exhaustion-capture` retains native-positive provenance and completeness work beyond this release's diagnostic capture facility.
- `quota-kickoff-reporting-and-alias-guidance` retains its remaining reporting work. This release includes required account-alias guidance only.
- Peer account inference, machine exhaustion caches, other native adapters, and broader Codex versions remain deferred under FINAL.
- CLI WinGet remains held. npm and global-core publication remain owner-only.

## Dismissed findings

None. Every concrete finding in both round-01 files has a proposed fix above.
The reported capture failure concerns the fixture. Both peers observed that an actual unsafe file is rejected.

## Verdict conflicts

### AC-05 — Transient retry behavior

Zcode wrote in `round-01/zcode-1.md`:

> Attempt: I checked attempt limits, the five-second delay, and the frozen ceiling in `RunParticipantStep`. I ran `TestExhaustionValidOutputAndSingleAttemptTransient`. I checked `SingleAttempt` on implementation, fix-up, protected steps, and started drafters.
> Result: not broken. Ordinary steps keep two attempts, the delay, and the frozen ceiling. Code-mutating and protected steps run once. No new automatic retry exists for them.

This result has no explicit provenance tag. Section 15 therefore reads its verdict as `RECALL` and `UNVERIFIED`.
Kimi wrote in `round-01/kimi-1.md`:

> Result: `WRONG (PRIMARY)`. `internal/app/app.go:1965-1970` assigns either identity to `Assignment.Implementer`. `internal/quota/quota.go:143-145` treats that assignment as protected. `internal/app/preflight.go:964-972` then selects one attempt.
>
> This contradicts the legacy comment at `internal/quota/quota.go:140`. It can also block global-default fallthrough. See the first finding.

Proposed resolution: F1 addresses the specific pre-pin path identified by Kimi. Zcode's cited test exercises an explicit single-attempt option.
That test does not refute the eager protection counterexample. Preserve its narrower result; do not use it to close AC-05.
Both signers must assess this resolution. The next independent review must test the corrected path and record admissible provenance.

### AC-08 — Current execution readers

Zcode wrote in `round-01/zcode-1.md`:

> Result: not broken. The resolver selects only from validated typed history. Filename order supplies nothing. A changed base hash, a corrupt tip, or a symlink blocks without fallback. A pending handover admits only the selected successor's handover action (`internal/runner/telemetry.go`, `selectedHandover`). All consumers I audited route through the resolver: driver checks, evidence table and report guard, driver evidence, cursor and phase digest, `advanceFinal`, implementation status and completion, fix-up paths, goal-check prompt, signoff recovery, learn, pipeline, TUI, retro, and action identity. `ExecutionStarted` keeps a pre-Phase-5 handover from implying review readiness.

This result also lacks an explicit provenance tag. It reads as `RECALL` and `UNVERIFIED` for acceptance purposes.
Kimi wrote in `round-01/kimi-1.md`:

> Result: `WRONG (PRIMARY)` for the complete current-reader scope. Resolver and pending-action tests passed. The fixed-path consumers mostly use `ResolveActiveExecution`. However, the phase reader has a confirmed error-handling regression.
>
> Executed: `go test ./internal/driver -run '^TestRebuildDetailSurfacesReadErrors/unreadable_IMPLEMENTATION.md$' -count=1 -v`.
>
> Output: `partial detail phase=round, want impl (the stat still saw the file)`. See the third finding.

Proposed resolution: retain the successful resolver cases and fix the unreadable-file case through F3.
Resolver routing does not prove correct error handling. The specific failed test prevents a complete AC-08 pass.
No vote supplies the missing evidence. The next review must verify the fixed reader at current source.

## Drafter position changes

None in the frozen design. F1 restores D1 and D4; F3 restores D6.
The implementer withdraws any implication that its checked implementation checklist proves complete AC-05 or AC-08 satisfaction.
The prior checklist in `../IMPLEMENTATION.md` records implemented areas, not independent acceptance evidence.

## Alternatives disposition

- **ALT-F1: adopt separate succession and protection decisions.** This preserves ordinary retries and the owner order without a new membership store.
- **ALT-F2: reject universal one-attempt behavior.** It contradicts FINAL D1 and AC-05.
- **ALT-F3: reject a fixed-umask test requirement.** The fixture must construct the unsafe mode it claims to test.
- **ALT-F4: adopt active-record wording.** It corrects the wait text without another public digest field.

## Coverage & blind spots

Kimi found the transient protection defect, version mismatch, and phase-read regression. Both peers found the umask fixture defect.
Zcode found the two display/documentation issues. Both peers reconstructed the exact signed protocol hunks and checked packet sizes.
Their agreement is not independent proof by itself. Related models can share assumptions.
An untested launch path, transition boundary, or current-state consumer can invalidate the proposed correction.
Native exhaustion remains synthetic or source-backed within FINAL's limits. No provider exhaustion was induced.

The completed organizer full suite also failed the version and phase-read checks. Its log is private `review01-full-go.log`.
The older `resume3-full-go.log` compiled before the kickoff-order correction. It supplies no current-source result.
The fixed source must pass the full Go suite, relevant race checks, skill suite, and package/build checks before close.

## Release and operator questions

The required Go command is `go test ./... -count=1 -timeout 45m`. F4 removes the ambient-umask dependency.
A fresh Kimi or Zcode process must perform the pre-release goal check after review closes.
Then follow D12: merge commits, clean-source builds, authorized channels, runtime installs, and a separate channel-verification process.
AC-15 delivery remains open until the hosted and installed byte checks exist. This fix consensus grants no completion claim.
Verify Windows skill asset hashes after the release workflow completes, because it can replace those assets.
The Zcode follow-up remains deferred; no execution date is promised here.

A corrupt selected handover blocks without fallback. The operator must preserve its bytes and resolve the integrity defect.
Automatic removal or replacement is not authorized by a timeout. The current runbook already states that corrupt authority blocks.

The phase-7 driver refused with: `continue failed: idea "meta-protocol-change-credit-exhaustion-failover" has no runs yet`.
The owner-authorized native fallback applies. It supplies no driver accounting or membership authority.

## Protocol context

context_mode: full
source_sha256: b9284d7d7062ef49d59439a23117dd7434e494340b5b564a64651a12eb324709
packet_sha256: b9284d7d7062ef49d59439a23117dd7434e494340b5b564a64651a12eb324709
fallback_reason: none

## Signoffs

### Signoff: codex-1 — 2026-10-29
Status: ✅ ACCEPT
Notes: I accept F1-F6 and the stated evidence limits. I will apply the fixes after both independent signers accept.

### Signoff: zcode-1 — 2026-10-29
Status: ✅ ACCEPT
Notes: I accept F1-F6 and the stated evidence limits. This signoff agrees on fixes. It passes no criterion and releases nothing.

I inspected every F1-F6 site in the unchanged source myself. The eager fallback/default assignment sits near `internal/app/app.go:1965`. `Roles.Protected` includes `Assignment.Implementer` at `internal/quota/quota.go:143-149`, against its own comment. Readiness passes `SingleAttempt: opts.QuotaRoles.Protected(a.ID)` at `internal/app/preflight.go:971`. `VERSION` says `1.55.0`; `internal/app/version.go` says `1.56.0`. `RebuildDetail` overwrites the stat result at `internal/driver/cursor.go:268-273`. The fixture writes mode `0644` through `os.WriteFile` at `internal/telemetry/capture_test.go:98`. `outstandingAgents` prints "IMPLEMENTATION.md not filed". `CHANGELOG.md` keeps an empty `## Unreleased` heading.

F1 needs one binding clarification. Removing `Assignment.Implementer` from protection alone is insufficient. `Evaluate` passes every candidate ID to `SelectSuccession` at `internal/quota/quota.go:273`, whatever rule caused the loss. The fix must gate succession on exhaustion-rule losses there. It must also cover both protection sites: `Roles.Protected` and the inline `protected` map in `Evaluate`. Otherwise an ordinary fallback failure would synthesize a role transition. F1's wording already requires this. The whole-batch tests must include the mixed ordinary-plus-exhaustion loss case.

The AC-05 and AC-08 resolutions follow §15.3. They engage the contrary evidence and preserve the narrower passing results. No vote proves a criterion. My round-01 AC-05 and AC-08 results carried no provenance tags. I accept their RECALL reading. My current inspection confirms Kimi's counterexamples.

F5 resolves the reported display defect. A new digest field is not needed for a display-only defect. The new wording must stay correct for absent records and pre-implementation handovers.

The conditions in this consensus bind: apply F1-F6 first; record current-source results in IMPLEMENTATION; both reviewers inspect the fixed complete diff; the next review records current-source PRIMARY evidence; the full Go suite passes at the fixed HEAD before close.

### Signoff: kimi-1 — 2026-10-29
Status: ✅ ACCEPT
Notes: I accept F1-F6 as a fix plan only. This signoff gives no product pass, criterion verdict, or release approval. I inspected the unchanged source and did not rerun the recorded failing suite.

F1 must change both `Roles.Protected` and the inline `protected` map in `Evaluate`. `Evaluate` must request succession only for role losses backed by `provider-exhaustion.v1`. An ordinary fallback or global-default failure must never create a role transition. Before a pin, those identities must retain transient retries, ordinary settlement, and global-default fallthrough. Explicit designee, pin, organizer, facilitator, and started-drafter protection must remain unchanged. Whole-batch tests must settle every loss before selection. They must cover mixed ordinary-plus-exhaustion losses, simultaneous role losses, successor usability, immutable history, and global-default fallthrough. Tests must also cover readiness, explicit designee/pin gates, and ordinary settlement.

F2-F4 and F6 address their recorded defects directly. F5's active-record wording fixes the stale display without a new digest field. Its tests must cover absent records and handovers before implementation starts. The AC-05 and AC-08 resolutions engage the contrary primary evidence and correctly limit the untagged results under section 15. No vote proves either criterion. Corrupt authority must continue to block, without automatic deletion or fallback. The next review must inspect the fixed complete diff and record current-source PRIMARY evidence. The full Go suite must pass at the fixed HEAD before close.
