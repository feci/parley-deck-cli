---
agent: hermes-1
idea: completion-contracts-evidence-ledger
review-round: 1
date: 2026-07-04
reviewed-commit: e0f2b45
---

## Summary

I reviewed the implementation at e0f2b45 against FINAL.md + consensus.md, focusing on
whether it matches the minimalist design I championed (extend `checks:` to a list, ledger
= existing `## Validation evidence` section, one Phase-8 fail-closed veto, v1 exit-0 only,
strict_gate independent, protocol text scoped to the list shape).

The design faithfulness is high. The implementation extends `checks:` via
`ReadChecksContract` (internal/driver/checks.go) using `yaml.Node.Kind` to distinguish
scalar/absent (legacy) from list (contract) — exactly the structured-reader approach
consensus item 8 specified. The ledger is `writeValidationEvidence` + `replaceSection`
overwriting the existing `## Validation evidence` section (internal/app/driver_checks.go).
The veto runs through the existing `RunChecks` gate (internal/app/driver_impl.go:223-227),
which is the documented deviation from FINAL.md item 4 (a separate impl.go branch) — same
fail-closed outcome, smaller surface. No `done_when:` key, no separate `review/evidence.md`
artifact, no append-only rule, no matcher grammar, no auto-close. Protocol text (LE-4 +
Phase-5 template + Phase-8) is byte-identical across all three COOPERATION.md copies (live
deck, embedded default, skill fallback) and is scoped to "when `checks:` is a list." Drift
guard `TestEmbeddedDefaultMatchesLiveDeck` passes.

Tests: `go test ./internal/driver ./internal/app ./internal/protocol` — all green. `go vet`
clean. The implementation is minimal and stays within the agreed surface.

However, I found one MAJOR integration bug (the evidence write is never committed, which
will dirty the tree and trip the `gitTreeClean` guard on the next fix-up cycle), one MAJOR
test-quality issue (the secret-scrub test is vacuous), plus minor items. Details below.

## Findings

### MAJOR-1 — Evidence write is never committed; dirties the tree and trips gitTreeClean

**Location:** internal/app/driver_checks.go:106-132 (writeValidationEvidence), interacting
with internal/driver/impl.go:273, 281.

`writeValidationEvidence` does `os.WriteFile` on `IMPLEMENTATION.md` (driver_checks.go:131)
but the driver never commits that change. In the fix-up loop the sequence is:

1. impl.go:273 — `gitTreeClean(d.cfg.Root)` must be true to proceed
2. impl.go:276 — `Fixup(ctx, cycle)` — implementer runs, commits its work
3. impl.go:281 — `RunChecks(ctx)` → for list-form, `runChecksContract` →
   `writeValidationEvidence` writes uncommitted changes to IMPLEMENTATION.md
4. impl.go:286-294 — marker, archive, open next review round

On the next tick, if review produces more agreed fixes, impl.go:273 `gitTreeClean` fires
again — and the tree is now dirty (the evidence write from step 3 was never committed). The
driver escalates: "git working tree is dirty; refusing to run a fix-up." This breaks the
fix-up loop for any list-form idea that needs more than one fix-up cycle.

The same uncommitted write happens at the pre-review gate (impl.go:109 in advanceImpl):
Implement commits → RunChecks writes evidence (uncommitted) → OpenReviewRound. The dirty
tree then persists into the first fix-up cycle's gitTreeClean check.

This is exactly the open question I raised in round-02 item 1 ("who commits the evidence
table?"): the driver must write AND commit the section, or the two writers (implementer +
driver) collide via the tree-clean guard. The implementation addressed the write but not
the commit. The `runChecksContract` warning at driver_checks.go:64 ("could not write
validation evidence") also silently continues on write error — the evidence is advisory to
the veto (which keys on exit codes, not the written table), so the veto itself still fires
correctly, but the uncommitted-dirty-tree interaction is a hard break.

**Fix:** after `writeValidationEvidence`, the driver must commit the IMPLEMENTATION.md
change mechanically (e.g. `git add <ideaDir>/IMPLEMENTATION.md && git commit -m
"[driver] <slug>: validation evidence — cycle N"`), mirroring how `DraftReviewConsensus`
already commits a driver-authored artifact. Or the evidence write must be deferred to the
implementer's own commit cycle (but that reintroduces the human step the idea set out to
remove). The commit approach is consistent with the existing driver-artifact pattern.

### MAJOR-2 — Secret-scrub test is vacuous (assertion can never fail)

**Location:** internal/app/driver_checks_test.go:14-17.

```go
if got := scrubAndTruncate("api_key=***\nok"); strings.Contains(got, "supersecret") {
    t.Fatalf("secret not scrubbed: %q", got)
}
```

The input `"api_key=***\nok"` does not contain the string "supersecret" anywhere, so
`strings.Contains(got, "supersecret")` is always false regardless of whether scrubbing
worked. The test passes trivially and proves nothing about the scrub regex. I verified the
regex itself is correct: `api_key=supersecret` → `api_key=«redacted»` (manual check
confirms the pattern `(token|secret|password|api[_-]?key|bearer|authorization)[=:\s]+\S+`
redacts all five credential shapes). The production code is safe; the test is not.

This is a safety-critical claim (FINAL.md verification criteria: "no raw unbounded dump, no
secrets"). A vacuous test here is worse than no test — it creates false confidence. The
test must use an input that actually contains a secret-shaped token and assert it is
absent from the output, e.g. `scrubAndTruncate("api_key=supersecret123\nok")` and check
`!strings.Contains(got, "supersecret123")`.

### MINOR-1 — gofmt drift on checks_test.go; IMPLEMENTATION.md claims clean

**Location:** internal/driver/checks_test.go:50-53; IMPLEMENTATION.md:45.

`gofmt -l internal/driver/checks_test.go` reports the file as unformatted (comment-column
alignment in the malformed-cases slice). IMPLEMENTATION.md:45 claims "gofmt -l clean." The
delta is cosmetic (tab-aligned trailing comments), but the claim is inaccurate and `go
test` does not enforce gofmt, so this will not be caught downstream.

**Fix:** `gofmt -w internal/driver/checks_test.go`.

### MINOR-2 — version.json metadata stale (updatedAt / updatedBy not refreshed)

**Location:** parley-deck/meta/version.json:12-13.

The `protocolSha256` was correctly refreshed (d9e060 → b33d2f) to match the new protocol
text. But `updatedAt` is still `2026-07-03T00:00:00.000Z` and `updatedBy` still says
`protocol-restructure-appendices (claude-1)` — neither was updated to reflect this idea.
IMPLEMENTATION.md:41 claims "`meta/version.json` protocolSha256 refreshed" (true) but a
reader checking `updatedAt`/`updatedBy` will misattribute the change.

**Fix:** set `updatedAt` to `2026-07-04` and `updatedBy` to
`completion-contracts-evidence-ledger (claude-1)`.

### MINOR-3 — No end-to-end test that RunChecks dispatches to the contract path

**Location:** internal/app/driver_checks_test.go (tests runChecksContract directly);
internal/app/driver_impl_le_test.go (tests scalar RunChecks paths).

`TestRunChecksContractWritesEvidenceAndVetoes` tests `runChecksContract` directly
(bypassing `RunChecks`). `TestRunChecksHonorsChecksCommand` tests the scalar path through
`RunChecks`. No test exercises `RunChecks` with a list-form `checks:` frontmatter to verify
the dispatch at driver_impl.go:223-227 (ReadChecksContract → isList → runChecksContract).
The dispatch logic is simple (3 lines), but it is the integration seam between the reader
and the runner, and a regression there (e.g. isList check inverted) would silently fall
through to the scalar path. A test that writes a list-form `00-prompt.md` and calls
`RunChecks` would close this gap.

### NIT-1 — IMPLEMENTATION.md head-commit (f0878b6) is one commit behind reviewed HEAD (e0f2b45)

**Location:** IMPLEMENTATION.md:8 (`head-commit: f0878b6`).

f0878b6 and e0f2b45 are the same content (the diff is a single line: the head-commit field
itself, going from "pending-commit" to "f0878b6"). e0f2b45 is the commit that recorded
f0878b6 as the head-commit. This is the normal "commit the head-commit field" bootstrap and
is not a defect — but a reviewer checksumming e0f2b45 should note the IMPLEMENTATION.md
frontmatter names f0878b6, not the commit they are reviewing. No action needed; documenting
for traceability.

### NIT-2 — Phase-8 veto has no head-commit equality check (accepted deviation)

**Location:** internal/app/driver_impl.go:223-227 (veto via RunChecks); consensus.md item 3
("latest driver run ALL-PASS at current HEAD"); my round-02 refinement (head-commit
equality check on the stored entry).

My round-02 explicitly proposed a ~5-line head-commit equality check at the completion path
(impl.go:201) to guard against stale evidence: read the latest evidence entry, verify
`AllPass && HeadCommit == HEAD`. The implementation instead delivers the veto through
RunChecks at the two existing call sites (pre-review impl.go:109, post-fixup impl.go:281),
relying on the invariant that the tree does not change between the post-fixup RunChecks and
Complete. This is the documented deviation in IMPLEMENTATION.md:47-52. The approach is
valid for the common path (no out-of-band commits between post-fixup and close), and is
simpler than a third call site. The protocol text says "at the current HEAD" which is
satisfied by "the latest RunChecks ran at the current HEAD." I accept this as a reasonable
simplification — the head-commit check was defense-in-depth against a narrow out-of-band
window, not a primary safety gate. Noting it here so reviewers are aware the explicit
stale-evidence guard from round-02 was not implemented.

## Open questions

1. **MAJOR-1 fix approach.** Should the evidence commit be a new mechanical git commit by
   the driver (mirroring DraftReviewConsensus), or should the implementer be instructed to
   include the evidence section in its own commit? The former is consistent with the
   driver-artifact pattern and removes the human step; the latter avoids a new
   driver-writes-to-implementer-file category. I lean toward the driver commit, but this
   touches the owner model I flagged in round-02 and was never explicitly resolved in
   consensus.

2. **Evidence write on pre-review RunChecks (impl.go:109).** Even if MAJOR-1 is fixed for
   the post-fixup path, the pre-review RunChecks also writes evidence (uncommitted) before
   OpenReviewRound. Is the intent to write evidence at both call sites, or only post-fixup?
   Writing at both is fine if both are committed, but the pre-review evidence is overwritten
   by the post-fixup evidence in the same cycle, so it may be redundant. Clarify whether
   the pre-review write should be suppressed for list-form (run the checks but skip the
   table write until post-fixup).

3. **Flaky escalation path — is it actually reachable?** Consensus item 6 says a failing
   criterion "escalates via §14 stopping-judgment." In the implementation, a failing
   RunChecks returns `(false, detail)` → advanceImpl/advanceReview returns ActionEscalated
   → the driver stops and surfaces the error. This is correct (no auto-retry loop). But
   MaxFixupCycles bounds the fix-up loop, not the RunChecks escalation — a persistently
   flaky check that passes pre-review but fails post-fixup would escalate immediately
   (ActionEscalated), not consume a fix-up cycle. Is this the intended behavior, or should
   a post-fixup RunChecks failure feed back into a fix-up cycle (giving the implementer a
   chance to fix the flaky command) rather than escalating on the first failure? The
   protocol text ("escalates via stopping judgment rather than auto-retrying") suggests
   immediate escalation is correct, but it means a single flaky post-fixup failure stops
   the loop with no retry budget. Confirm this is intended.
