# codex-1 cycle-3 frozen host evidence

Producer evidence only, not independent acceptance. Verified 2026-10-06T16:09:29.377755+00:00.
CLI product commit `fac40aaa1bba5350351bab0310ea8ede30ed3e6e`;
skill `e976f7c9515f3250761380c1e09295e0f6c59985`. Inputs/outputs hash-identical; roster hash set unchanged. All
provider commands are shadowed by local guards; both the cycle-3 and reused cycle-2
guard-denial logs are absent. The original sandbox failure of native boot proof remains historical evidence.

| Check | Exit | Seconds | Raw host log |
| --- | ---: | ---: | --- |
| native-boot | 0 | 0.006 | `.parley-runtime/quota-implementation/fixup-3/host/native-boot.log` |
| native-lease | 0 | 0.759 | `.parley-runtime/quota-implementation/fixup-3/host/native-lease.log` |
| full | 0 | 567.41 | `.parley-runtime/quota-implementation/fixup-3/host/full.log` |
| build | 0 | 3.373 | `.parley-runtime/quota-implementation/fixup-3/host/build.log` |
| vet | 0 | 1.745 | `.parley-runtime/quota-implementation/fixup-3/host/vet.log` |
| windows-build | 0 | 3.413 | `.parley-runtime/quota-implementation/fixup-3/host/windows-build.log` |
| race | 0 | 22.441 | `.parley-runtime/quota-implementation/fixup-3/host/race.log` |
| shared-volume | 0 | 9.472 | `.parley-runtime/quota-implementation/fixup-3/host/shared-volume.log` |
| local-control | 0 | 11.761 | `.parley-runtime/quota-implementation/fixup-3/host/local-control.log` |
| reviewer-zadv | 0 | 0.346 | `.parley-runtime/quota-implementation/fixup-3/host/reviewer-zadv.log` |
| reviewer-changing-countdowns | 0 | 3.518 | `.parley-runtime/quota-implementation/fixup-3/host/reviewer-changing-countdowns.log` |

Formatting: 93 changed Go files, none unformatted. Full skill suite exited 0
in 163.923 seconds: 399 Node tests, 54 Python tests, six addon manifests. Core
SKILL.md is 19992 bytes; unchanged cap 20000.

Deck and skill protocol copies are byte-identical at 73613f95aabff8ed57372a535913985554014b71911cc8126ec3f7890a9d4f4e.
The embedded bootstrap hash is a92d8123ff2e924fff699497636e300fa176899fe5f56f7371a7399aaa5aef66;
its intentional workspace/date/host-handle normalization passes the unmodified drift gate.
Phases 0/5/8 render full live source, no fallback; packet/source hashes match the deck.
No cycle-3 protocol hunk changed. Installed skill/CLI remain 2.14.0/1.50.0; source-deck
metadata is valid but stale, with read-only sync dry-run retained. Release/install work
has not occurred. The project source is never replaced by the older packaged protocol.

Identical real CLI baseline/current D1–D4 probes pass both filesystems. See parent
supplement for the explicit Missing-signoffs diagnostic difference on a pending decline.
Native actual-parley crash exercise ran independently of the unit helper on each volume:
one real supervisor and two verified local sleeping workers were killed, then recovery
ran twice. Both recoveries exit 0; exactly two immutable crash settlements, unchanged
on replay; zero fabricated normal terminal records. Same-host/boot proof executes on
this macOS host. Windows runtime/native recovery and foreign/container PID namespaces
remain unverified; cross-build is compilation only.

The original parent harness failures (raw protocol equality; shell-before-exec race;
manifest.json versus run.json filename) are disclosed in the parent supplement. No
product/assertion changed to make the host checks pass. Frozen source hashes bind the
successful checks to the product snapshot. Native AC2 remains NOT MET / owner-waived,
R5-MAJOR-2 owner-accepted and deferred. No real zcode invocation or capture occurred.

Selected bytes and hashes: `fixup-cycle-3-20261006/manifest.json`. Private full logs remain
at the paths above. The independent reviewer must supply its own verdict and identify
which executions it repeats; this report never substitutes for that review.
