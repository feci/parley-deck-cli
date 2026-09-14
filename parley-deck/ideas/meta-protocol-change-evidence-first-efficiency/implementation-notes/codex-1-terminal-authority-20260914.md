---
agent: codex-1
idea: meta-protocol-change-evidence-first-efficiency
date: 2026-09-14
source-base: 83d0475fd81acc7a33a66b2a2f1588f1b870aa04
status: integrated-source-validation-passed
---

# Retain the actual live terminal without rereading older content

Run.Finish now uses the existing guarded control-authority reader to retain
an actual outcome despite unavailable older archive, reservation-intent or
resolution-result bytes. It still requires canonical policy, the complete
published ledger, the original charge/launch/invocation and the exact canonical
state frozen into its runtime-only handle when Begin published the launch.
Changed original state or charges refuse; terminal replay cannot replace bytes.

The current post-source is still observed and archived through the existing
full capture code. If that observation/capture fails, the existing explicit
source-unavailable/archive-unavailable terminal is retained with its error.
No old source or missing terminal is reconstructed from today's tree. Full
source/result validation remains on execution, reconciliation and acceptance.
The runtime binding has no persisted-handle or historical-replay API.

## Direct production-boundary tests

terminal_authority_test.go adds four cases over real fixture reservations,
telemetry and local exit-7 child processes:

- Removing the historical baseline after the actual process exit does not lose
  its failed terminal or changed post-source archive. Original ledger bytes stay
  equal; ordinary Inspect and completion still refuse incomplete history.
- A first real unchanged invocation is reconciled and explicitly continued.
  Removing its terminal after a second actual process exit still permits the
  second original terminal, without making the older history acceptable.
- Rewriting the test-owned attempt's reservation amount must fail the exact
  original-ledger predicate and leave state unchanged.
- Rewriting a structurally valid precharge-intent hash after Begin must fail
  the frozen-state predicate and leave state unchanged. Original policy and
  ledger alone cannot authorize a rewritten runtime lineage.

Both successful retention cases require exact terminal replay refusal without
changing the retained state. These are execution fixtures, not independent
model acceptance or a demonstration of historical recovery.

## Native prototype provenance

.parley-runtime/terminal-authority-prototype-20260914/ retains the original
route-only candidate, its initial unused-import failure, a corrected focused
pass (6.065s), surrounding checks (66.108s), two old-route negatives and its
4.1848765s full-size publication observation. The refined candidate freezes
Begin's exact state; four authority cases passed (7.791s), removal of that
binding failed at the intended assertion (2.243s), and the same full-size
fixture published in 4.020197291s. Each successful synthetic state is retained
separately; the test restored the exact original pending input after each probe.

The synthetic envelope uses 128 real low-level charges/local lifecycle records,
127 fixture-assembled resolutions/continuations and an exact 256 MiB archive.
The Run handle in that probe is fixture-constructed, not a normal end-to-end
Begin/runner launch. Other Go validation ran concurrently; no uncontended
throughput, cold-cache or universal deadline bound is claimed. Original
production Inspect and Finish deadline failures remain preserved independently.

## Validation and remaining limits

Integrated focused/negative/envelope evidence will be recorded after termination
in .parley-runtime/terminal-authority-development-20260914/. Full suite, six-package
race, vet, Windows CLI/trajectory cross-builds and shared compiled selections are
required against a new frozen Go/module manifest before publication.

This prevents one prospective terminal-loss boundary. Historical absent/abnormal
terminals, consumed/incomplete helper tickets, custody, escaped descendants and
workflow effects still require explicit recovery. This responds to F6's publication-coupling finding. General history inspection
still revalidates old archives/results and is not claimed faster; a fresh
participant review must decide the finding's disposition.
Current-source participant acceptance/signatures, live launch/concurrency/closure
checks, packet/full-six experiments, final offline HTML and delivery follow-ups
remain open. Existing participant artifacts and historical HTML are unchanged.


## Completed integrated validation

Manifest 2d4316bf1d6e43873bd0a508529081fb983ea4f14e6fa521cf2b9a12f094ed2f pins 418 Go/module files.
The focused 22 top-level tests passed in 46.546s. Three integrated negatives
failed at the intended old archive route (1.558s), old result route (2.152s)
and removed state binding (1.831s) assertions. The production-source envelope
with only a test-file overlay published in 5.192811625s under the original
30-second context, retaining original charges and 127 resolutions. The exact
original pending input was restored after storing successful-state bytes.

Full Go suite PASS 412.84s (all 32 package terminals);
six-package race PASS 463.156s; vet PASS. Full/race ran
concurrently over the same frozen source with separate native test processes
and fixtures. Windows CLI and trajectory cross-builds passed; PE amd64 checked,
runtime unverified. Compiled shared trajectory selection PASS
76.562s. All 22 focused tests and the four
new authority cases are accounted for in full/race/shared evidence. Exact
source/log hashes and unchanged participant review/historical HTML are checked.

The earlier disk-exhaustion failure belongs to the cleanup checkpoint and
remains preserved there. Native prototype timings and original 30-second
Inspect/Finish failures retain their own source provenance. No retry was used
to convert those historical failures into passes. Current-source participant
review must still decide the F6 disposition. The separate owned abnormal-recovery
proposal is discussion input, not an implemented operator or custody grant.
