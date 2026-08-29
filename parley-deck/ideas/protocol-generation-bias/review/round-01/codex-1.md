---
agent: codex-1
idea: protocol-generation-bias
review-round: 1
date: 2026-08-29
reviewed-commit: 59eb663
---

## Summary

Changes are required. The byte reduction and the direct unit validator are genuine, and the three
present §15.6 copies contain the same replacement. The implementation nevertheless leaves all three
behavioral carriers incomplete: the acquisition validator is not called by the runtime, the exchange
does not exist, and the alternatives-disposition scanner does not exist. `PRIMARY`: `git diff
--name-status 9d4f45c..59eb663`, the call-site searches and focused tests recorded below.

**D1 ruling — exchange fidelity.** The implementer was right not to silently replace the ratified
one-packet/no-Decide design with a different protocol. The primary paper did measure two Exchange
rounds plus one Decide pass, with four agents, and Table 7 reports GPT-4.1 `0.037 -> 0.800`.
`PRIMARY`: [HiddenBench v4 §6.4 and Table 7](https://arxiv.org/html/2505.11556v4#S6.SS4), which
states `Exchange stage (2 rounds)`, `Decide stage (1 pass)`, and `4 agents`. That evidence means the
ratified one-packet form must remain labelled as an unverified transfer; it does not grant an
implementer authority to rewrite a frozen `FINAL.md`. Implement the ratified form here, or ratify a
successor design before adopting the measured two-round-plus-Decide form.

**D4 ruling — withholding the exchange clause.** Withholding prose until its runner carrier exists
is the correct application of the carrier thesis. Withholding the ratified behavior and then marking
the package implemented is not. The valid conclusion is “land the runner and protocol duty together
in this implementation,” not “move leg 2 to an unspecified later idea.”

**Late `opencode-1` ruling.** The disposition leg is not complete without reading and explicitly
disposing the late artifact. `PRIMARY`: `round-01/opencode-1.md` proposes a parallel `REFRAME`
class and a `## Frames considered` destination, while `FINAL.md:66,136-140` both retains the four
severity names and makes that artifact a live implementation input. The neutral `ALT-` route ratified
in `FINAL.md` should be implemented now. Adopting `REFRAME` or `## Frames considered` would change
the frozen design and therefore needs explicit review consensus or a successor idea, not an
implementation-side invention.

## Refutation attempts

1. **Byte claim.** `PRIMARY` — executed the ratified command against the tag and then extracted the
   same heading-bounded section from the reviewed commit:

   ```text
   git show protocol-generation-bias-baseline:parley-deck/COOPERATION.md | sed -n '1346,1368p' | wc -c
   -> 1372
   git show 59eb663:parley-deck/COOPERATION.md | sed -n '/^### 15\.6 /,/^### 15\.7 /p' | sed '$d' | wc -c
   -> 896
   1372 - 896 -> 476
   ```

   I failed to refute the arithmetic: the current text is 896 B and the literal reduction is 476 B.
   This does not validate the full ratified package because the exchange text was omitted to obtain
   that number.

2. **Three protocol copies and guard coverage.** `PRIMARY` — extracted §15.6 from
   `internal/protocol/defaults/COOPERATION.md`, `parley-deck/COOPERATION.md`, and
   `../parley-deck-skill/skills/parley-deck/references/COOPERATION.md`; all three produced SHA-256
   `9914c4d1eccebec63ab10d3f0866397762d230df20443943184219038e93e3a1`. Their baseline/current
   character counts were respectively `104885/104413`, `105131/104659`, and `104974/104502`, so
   every copy moved by the same 472 characters and the pre-existing offsets remain deck `+246` and
   skill `+89`—not the claimed skill `+90`. `internal/protocol/drift_test.go:84-136` compares only
   the embedded default and live deck. Its four-surface test at lines 323-377 checks a short banned
   phrase list, not byte parity. The third copy is consistent now but is not protected by that drift
   comparison.

3. **Mutation check.** `PRIMARY` — first ran
   `go test ./internal/protocol -run '^TestValidateRoundOneArtifactRequiresExistingAlternatives$'
   -count=1` successfully. I then changed only the condition to `if len(raw) < 0 {`, retained the
   test, and reran it. It failed with exit 1 at `roundartifact_test.go:25`: “a round-1 artifact without
   an Existing alternatives section must be rejected.” After restoring the condition, the same test
   passed and `git diff --exit-code -- internal/protocol/roundartifact.go` was clean. I failed to
   refute the direct unit test's mutation sensitivity.

4. **Runtime reachability of the new gate.** `PRIMARY` — `rg` found
   `protocol.ValidateRoundOneArtifact` only in its definition and protocol-package tests. The actual
   path is `internal/driver/driver.go:367 -> runner.ValidateRoundArtifact ->
   runner.ValidateRoundOneArtifact` (`internal/runner/validation.go:15-18,63-91`), and that runner
   validator never calls the protocol validator or checks `## Existing alternatives`. As a concrete
   counterexample,
   `go test ./internal/runner -run '^TestValidateRoundOneArtifactAcceptsMissingOpeningFence$'
   -count=1` passes even though its fixture contains no `## Existing alternatives`. This broke the
   claimed runtime gate.

5. **Exchange acceptance criteria.** `PRIMARY` — an `rg --glob '*.go'` search for `Evidence
   exchange`, `transfer unverified`, the HiddenBench label, `Exchange stage`, and `Decide stage`
   over `internal/` and `cmd/` returned exit 1. The same required label is absent from
   `parley-deck/COOPERATION.md`. The only changed runner file has two added prompt-template lines for
   acquisition. Thus criteria 3-5 cannot be exercised: there is no exchange phase sequence, prompt,
   instrumentation, or protocol label to test.

6. **Disposition acceptance criterion.** `PRIMARY` — `git diff --name-only
   9d4f45c..59eb663 -- internal/consensus internal/driver internal/runner internal/protocol` contains
   no consensus or driver change, and a Go-only search for `Alternatives disposition`, `REFRAME`,
   `Frames considered`, or `recorded adoption` returned exit 1. A recorded adoption therefore has
   no contradiction check, signoff block, owner escalation, or no-auto-halt test. Criterion 6 failed.

7. **Protocol-table consistency.** `PRIMARY` — all three copies say “Unconditional on every track”
   at the start of §15.6, but their unchanged §15.7 row says `15.6 correlated agreement | no | yes
   (section in an existing round-02 file) | yes (assigned round artifact)`. I failed to reconcile the
   new all-track rule with that obsolete per-track rule.

8. **Build and broader tests.** `PRIMARY` — `go build ./...` passed. Focused protocol and runner
   tests passed after restoration. `go test ./...` was not wholly green in this sandbox because
   unchanged `TestDurableKillEndToEndRealProcess` failed twice with “no recorded boot id”; `git diff
   --quiet 9d4f45c..59eb663 -- internal/runner/durablekill.go internal/runner/durablekill_test.go`
   returned 0, so I do not attribute that platform-sensitive failure to this change.

## Findings

### [CRITICAL] The acquisition gate is dead on the runtime path

What is wrong: the new `protocol.ValidateRoundOneArtifact` is tested in isolation but is never
called by the driver or runner. The live runner validator still accepts a round-one artifact that
omits `## Existing alternatives`. `PRIMARY`: call path and passing counterexample in refutation
attempt 4.

Why it matters: acceptance criterion 2 requires rejection, not a prompt warning. The proposal's
central carrier thesis is defeated if the only enforcing function is unreachable; a participant can
ignore the prompt and the driver will still close the round.

Concrete fix: compose `protocol.ValidateRoundOneArtifact(path)` into
`runner.ValidateRoundOneArtifact` (or consolidate the two validators into one shared runtime gate),
then add runner- and driver-level tests proving a missing or empty section prevents round completion.
Keep the protocol-package mutation test, and update the existing runner fixture so it no longer
silently demonstrates the bypass.

### [CRITICAL] D4 leaves the entire ratified exchange leg unimplemented

What is wrong: commit `59eb663` has no sealed packet collection, simultaneous release, `## Evidence
exchange`, next-decision injection, recall instrumentation, or required `transfer unverified`
protocol label. `IMPLEMENTATION.md:124` nevertheless says “I have implemented one packet,” while
its own status table and D4 say the exchange was withheld. `PRIMARY`: refutation attempt 5 and
`IMPLEMENTATION.md:99-104,124-127,145-159`.

Why it matters: acceptance criteria 3-5 fail and the behavior targeting HiddenBench's hidden-profile
failure is absent. Avoiding an unenforced prose rule is good staging, but it cannot turn an omitted
ratified leg into a completed implementation.

Concrete fix: implement the ratified one sealed packet per participant, with simultaneous release
into the already scheduled next decision prompt on every track, self-echo under `## Evidence
exchange`, no asymmetry assertion, bounded cost reporting, and the specified pre/post recall
instrumentation. Add the exact protocol label and tests for zero new rounds/files/agents/artifacts,
simultaneous release, prompt contents, and instrumentation. Do not silently add a second exchange
round or Decide pass; if the group wants the measured HiddenBench form, ratify a successor design.

### [CRITICAL] The disposition carrier and its required late-input decision are absent

What is wrong: the protocol now prints an `ALT-` duty, but no consensus template, parser, validator,
FINAL-consistency scanner, signoff block, owner escalation, or no-auto-halt test implements it. The
implementation record expressly calls the scanner “remaining code work” and records no decision on
the late `REFRAME`/`## Frames considered` proposal. `PRIMARY`: refutation attempt 6,
`IMPLEMENTATION.md:82-104,181-183`, `FINAL.md:50-66,131-140`, and
`round-01/opencode-1.md`.

Why it matters: acceptance criterion 6 fails, so an adopted alternative can still disappear from
`FINAL.md`—the exact B1 failure this leg exists to prevent. Leaving the deferred vocabulary input
undispositioned also repeats that absorption failure inside this implementation.

Concrete fix: add `## Alternatives disposition` to every consensus carrier, enforce the exact
`ALT-<agent>-R<round>-<index>` form plus adopt/reject and decisive reason, and compare recorded
adoptions before signoff. A contradiction must block signoff and produce the specified owner
escalation while the scanner itself never auto-halts; test both halves. Record an explicit decision
on the late artifact. Under this `FINAL.md`, use the neutral `ALT-` route and keep the four severity
names; adopt `REFRAME`/`## Frames considered` only through a successor design or an expressly
ratified change.

### [MAJOR] The per-track table still negates the new unconditional rule

What is wrong: §15.6 is now unconditional on every track, but §15.7 still exempts `fast` and still
describes the deleted standard/deliberation artifact forms. The contradiction appears identically in
all three protocol copies. `PRIMARY`: refutation attempt 7.

Why it matters: a manual participant can reasonably follow the per-track binding table and skip the
new acquisition/disposition duties on `fast`, while the intended runner prompt applies acquisition
unconditionally. The protocol and its carrier therefore define different obligations.

Concrete fix: replace the stale row in all three copies with accurate all-track rows for the new
§15.6 duties (or one unambiguous all-`yes` row), and add a focused textual assertion so future
rewrites cannot retain obsolete per-track language while the section says unconditional.

### [MAJOR] The canonical implementation record does not identify a reproducible reviewed state

What is wrong: `IMPLEMENTATION.md` records `head-commit: 9d4f45c`, which is the baseline parent,
not reviewed commit `59eb663`; it marks the work `implemented` while two ratified legs and the
runtime acquisition gate remain incomplete; and it contradicts itself about whether one exchange
packet was implemented. The reviewed commit changes only six files, while the skill's `SKILL.md`,
manifest, and third protocol copy remain uncommitted in the sibling repository. The stated skill
offset is also `+90`, while the independently reproduced offset is `+89`. `PRIMARY`: `git show -s
59eb663`, `git diff --name-status 9d4f45c..59eb663`, `git -C ../parley-deck-skill status --short`,
and refutation attempt 2.

Why it matters: a reviewer or future publisher cannot reconstruct the claimed three-surface
implementation from the named commit, and uncommitted carrier changes can be lost or changed after
review. Incorrect status and provenance obscure which acceptance criteria remain open.

Concrete fix: keep the implementation status non-terminal until all ratified legs pass; update its
head commit after the fix-up; remove the contradictory completion claims; correct the `+89` offset;
commit the sibling skill changes and record that repository's commit plus manifest hashes. If the
2.11.0 core is only staged for attended publication, record its exact source path/hash and the owner
command rather than saying an unlocatable core is prepared.

## Open questions

- Where is the claimed prepared `~/.parley/protocol/core/2.11.0/COOPERATION.md`? It was absent at
  review time. Is publication expected to materialize it from one of the reviewed copies, and if so,
  which exact path and hash is the attended input? `PRIMARY`: `ls
  /Users/tomasfecko/.parley/protocol/core/2.11.0/COOPERATION.md` returned “No such file or directory.”
- Does the group want a successor design that adopts HiddenBench's measured two-round Exchange plus
  Decide protocol, or does it intentionally want to ship and instrument the ratified one-packet
  transfer first? Either is reviewable; silently crossing between them is not.
- Can the implementer record the environment in which `go test ./...` was green? The only failure in
  this review environment was an unchanged macOS/process-identity test, so it is not a finding on
  commit `59eb663`, but the check claim is not reproducible here as written.
