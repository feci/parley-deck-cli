---
agent: codex-1
idea: meta-protocol-change-evidence-first-efficiency
date: 2026-09-11
source-commit: 78a7b13
status: tested-slice-independent-acceptance-pending
---

# Configured monetary ceilings at the shared launch boundary

## Reproduced problem and resulting behavior

The executed TestMonetaryDefaultRefusesUnreservedManualLaunch counterexample at
0a872c5 configured max_cost_usd = 1 and ran a real local shell through RunMeasured.
It failed because the unreserved call returned nil instead of a budget refusal.
The previous monetary check lived in the driver loop and did not prevent that
first manual launch. Negative log SHA256:
`840d367b61876b4317d4f009d41ddab5d6498d647a65eaef0527e0d996a9861d`.

Source 78a7b13 reads layered runtime defaults at the common process boundary.
A positive monetary default requires an already configured launch policy for
the same shared idea or auxiliary scope, with a matching finite microdollar cap
and explicit conservative per-launch reservation. The driver also applies this
requirement to a nonzero directly supplied MaxCostUSD before dispatch. A first
unreserved call refuses visibly and retains requested/terminal telemetry without
starting a child. An absent policy is not automatically initialized or funded.

Existing frozen policy and ledger accounting remain authoritative across new
runs and worktrees. Zero or omitted defaults do not disable them. Different
nonzero defaults do not extend or silently tighten the saved cap. A caller must
resolve that mismatch through the explicit policy controls; generic runtime
metadata and participant frontmatter cannot supply a grant.

Default precedence remains machine, deck, local, explicit environment-selected
configuration. An explicit deck zero may override a machine seed before policy
activation. Missing explicitly selected configuration, unreadable files and
parser errors no longer silently disable this check. Parser excerpts are not
copied into launch refusal diagnostics. Disposable review execution loads the
default from its captured live origin; handoff preparation activates no policy
and incurs no process charge.

Dollar ceilings convert from the configured float's shortest decimal spelling
to whole microdollars, rounding down without inflating the ceiling. Negative,
nonfinite, positive sub-microdollar and overflowing inputs refuse. An initial new
unit fixture incorrectly assumed a very large decimal literal survived float64
conversion exactly; the fixture was corrected to its representable value before
the final suites. This was a test-expectation correction, not a product regression.

## Executed validation

| Check | Observed result |
| --- | --- |
| go test -count=1 ./... | PASS; app 102.204s, budget 10.856s, driver 11.656s, runner 32.633s; wall 104.835s |
| go test -race -count=1 ./internal/budget ./internal/config ./internal/driver ./internal/runner | PASS; 14.247s / 1.692s / 10.728s / 60.329s; wall 62.376s |
| Scoped budget/config/driver/runner/app vet | PASS |
| Windows amd64 app test cross-build | PASS; Windows execution untested |
| Shared-volume budget/config/driver/runner fixtures | PASS; wall 0.776s / 0.338s / 0.347s / 5.835s |

The six process surfaces are manual, round-process, consult, preflight probe,
interactive execution and ACP. The monetary case has no launch-count cap, so its
second refusal depends on the monetary reservation. New fixtures cover known and
unknown terminal cost, a second run, removed defaults, explicit zero, invalid
amounts/configuration, direct driver refusal, live-origin context and handoffs.
Unknown observed costs remain null; a conservative reservation is not a price.

The 317-file Go/module source manifest was checked unchanged after validation.
Manifest SHA256:
`8841a7b8de2765cdfefb87acbafcca9c614282307dd40629533ad5141c9815f6`.
Full-suite log: 1405 bytes, SHA256
`56717330ec6785441f4986c731c924e61629cebb965d76a24a5522bfcc6b5998`.
Race log SHA256:
`dffa13040375af1d0351fd78532eefd30e57aea98674f4efe33f45567c2fc625`.

Shared fixtures ran under
`/Volumes/My Shared Files/AI_WORKSPACE/parley-monetary-fixtures-3pvu88p5`, verified
outside an enclosing Git repository. Binaries and captured output stayed in
local temporary storage; test output was captured through pipes and written
once, not streamed to the shared volume. Fourteen logs/manifests/metadata files
and verified checksums are retained in the integration worktree's ignored
`.parley-runtime/monetary-default-validation-20260911/`. Original storage is
`/var/folders/yt/p2sr23f12_qcfx_w2z5c1p4r0000gn/T/parley-monetary-default-validation-8oqj0piv/`.

## Remaining obligations and limits

This implements monetary-default enforcement through the existing operator
policy; it does not auto-choose a reservation or authorize funding. The declared
reservation is not a guarantee about actual provider billing. The older driver
usage check may still stop on incomplete observed cost even when reservations
exist; neither it nor this change certifies unknown usage as known.

Existing history, including retained refused requests, can prevent first policy
configuration and require the still-pending explicit legacy migration. Preserve
that evidence; deleting it is not recovery. Launch/step policy extensions, safe
guard recovery, durable semantic action replay, canonical refusal publication
and recovery, and independently confirmed patch-regression trajectory remain D6
work. Cycle extensions are already implemented separately at 321aef9.

Fresh independent source acceptance, actual real-model concurrency/closure, exact
phase-packet trials, twelve-task solo/duo/full-six comparison with frozen equal
ceilings and rotation, two blind nonauthor graders, owned reviews/signatures,
final populated HTML QA and real 14/30-day observations remain required. Pending
Claude recovery and historical quorum/pilot/funding decisions are unchanged.
No participant model call or real policy grant occurred: 34 terminal model
attempts, 17 unknown costs, USD 46.1887585 known CLI estimates and unknown total.
This note is the implementer's execution record, not an independent review.

## Updated report

The English offline HTML was rebuilt for source 78a7b13 at
2026-09-11T15:19:03.645782+00:00 and passes all ten report tests.
It is 605008 bytes, SHA256
`152cf59410db1c0097d46035ef61110ffa1d0d9a884fb5c85e20d059b916415f`.
Earlier browser evidence addresses older report content; this checkpoint does
not claim fresh exact-file browser QA or final populated-report validation.
