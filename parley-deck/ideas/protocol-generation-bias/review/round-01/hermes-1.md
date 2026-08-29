---
agent: hermes-1
idea: protocol-generation-bias
review-round: 1
date: 2026-08-29
reviewed-commit: 59eb663
---

## Summary

Implemented: acquire (leg 1, full), disposition protocol duty (leg 3 text, partial scanner), §15.6 deletion (funded), correlation clause (carried), byte accounting verified (
-477 B, better than ratified -237 B). Withheld: leg 2 exchange clause (D4 — see ruling below) and the "transfer unverified" label (absent from `COOPERATION.md`; present in `FINAL.md`/`consensus.md` artifacts only). Gate mutation verified (fails closed, restored). Three `COOPERATION.md` copies located: deck (105661 B) and defaults (105415 B) moved consistently; the skill reference third copy (`parley-deck-skill/skills/parley-deck/references/COOPERATION.md`) does not exist in this workspace — drift guard covers zero of two real copies, not two of three.

Explicit ruling on the two flagged items.
- **D1 (exchange fidelity):** The ratified one-packet form is WRONG to ship as-is. The measured HiddenBench protocol (`arXiv:2505.11556v4`, PRIMARY — two exchange rounds + Decide at 4 agents, producing +76.3pp) is structurally different from one packet with no Decide stage. Implementing the truncated form and citing the full protocol's effect size is the exact defect the carrier thesis names: a rule reaches the artifact but not the runtime. The implementer's refusal to silently upgrade a ratified design is the correct procedural stance; the correct substance is to upgrade it. **Block on D1: the exchange stage must be two rounds + Decide, or `FINAL.md` must be re-ratified to match the one-packet form explicitly (with corrected effect-size claims).**
- **D4 (clause b removed):** The implementer is RIGHT. Withholding an unenforced protocol clause from `COOPERATION.md` is the direct application of the idea's own finding (prose is not a carrier). Printing an unimplemented duty would have been the fourth instance of the named defect class committed by the very idea that named it. The byte consequence (-476 B vs ratified -237 B) is real and verifies independently (`sed` over `protocol-generation-bias-baseline` → 1,372 B; current section → 895 B; difference -477 B). The clause must be restored only when `internal/runner/` carries the exchange orchestration; until then the withheld form is the honest one.

## Refutation attempts

Assumed the implementation was wrong; tried to break it. All attempts and outcomes (verified directly, provenance tagged):

1. **Byte claim falsification.** Tried to show the -237 B figure was fabricated. `sed -n '1346,1368p' parley-deck/COOPERATION.md | wc -c` from tag `protocol-generation-bias-baseline` = 1,372 B (PRIMARY, git-object read). New §15.6 = 895 B (PRIMARY, file read). Net = -477 B, confirming the implementer's corrected -476 B claim. The original -726 B was the inherited zcode-1 error; the -237 B was the corrected ratified figure; -476 B is the actual post-deletion figure after withholding clause (b). Refutation FAILED — the byte math holds.
2. **Gate that does not fail.** Mutated `internal/protocol/roundartifact.go` to `if len(raw) < 0 {` (always false, so the gate never rejects a missing `## Existing alternatives`). Re-ran `TestValidateRoundOneArtifactRequiresExistingAlternatives`: test FAILED with message `"a round-1 artifact without an Existing alternatives section must be rejected"` (PRIMARY, `go test` output). Restored original; gate passes. Confirmed: gate is live, fails closed when neutered, restores to green. Refutation FAILED — the gate works.
3. **Three-copy consistency falsification.** Searched the workspace for all `COOPERATION.md`. Only two copies exist (`parley-deck/` and `internal/protocol/defaults/`); the skill-reference third (`..parley-deck-skill/skills/parley-deck/references/`) is absent. Both existing copies carry the new 15.6 title and lack the deleted `primarily a judgment` text; offsets match (104,413 vs 105,661 ≈ same delta pattern as pre-change). The Go drift guard (`TestDriftAnchorsAcceptTheGeneratedRosterTable`) covers the generated defaults only — zero of the two real copies it should cover, and zero of the three the design claims. Refutation PARTIAL — consistency holds on existing copies, but the missing third copy and limited guard coverage confirm the drift-risk finding, not refute it.
4. **Runtime invocation absence.** Searched all `.go` for `ValidateRoundOneArtifact`: called by `ValidateRoundArtifact` (`runner/validation.go:17`), which is called from `internal/driver/driver.go:367` and `internal/runner/phase58.go:321` (SECONDARY, grep read). It IS invoked in the runtime path. The implementer's claim that "a validator that exists and is never invoked" is the exact defect class — in this case it IS invoked, but its coverage is limited to design rounds; there is no design-round existence validator separate from review rounds. Refutation FAILED — it is called, but the absence of a separate design-round validator (as noted in `opencode-1.md`) is itself a real gap.
5. **Disposition leg completeness without reading `opencode-1.md`.** `FINAL.md` explicitly defers the vocabulary question to `opencode-1.md`. Read `opencode-1.md` (PRIMARY, full file read: proposes `REFRAME` class, `## Frames considered` destination, freeze gate, four-field payload, witness requirement). Without it, the disposition leg's vocabulary is undefined; with it, the vocabulary is defined but NOT adopted (no `REFRAME` in `COOPERATION.md` or any `.go` file). The deferred question therefore has a live input (`opencode-1.md`) but no adopted resolution. The disposition leg is NOT complete without reading `opencode-1.md`, and it is NOT complete after reading it either — it requires adoption of `REFRAME` or an explicit rejection. Refutation FAILED — the leg requires the deferred file, and the deferred file proposes a change the implementation has not adopted.
6. **HiddenBench group-size claim.** Read `arXiv:2505.11556v4` via `pdftotext -layout` (PRIMARY, local extraction). Confirmed: protocol evaluated on "4 agents" (verbatim quote), two exchange rounds + Decide. Confirmed dose-response (Table 4): 3 agents +34.8%, 4 +25.0%, 5 +6.4%, 6 +19.7%, 7 +0.6% — non-monotonic (N=6 beats N=5) and middle cells (N=5/6/7) hold 5 tasks each, only N=4 (58 tasks) is well-powered. Confirmed asymmetry condition: works without telling agents asymmetry exists. Confirmed labels in the file carry the same wording: "Exchange stage (2 rounds)" and "Decide stage (1 pass)". The reference brief's figures are accurate; the design's one-packet simplification is not supported by the source. Refutation FAILED — D1's upgrade argument is evidence-backed.
7. **"Transfer unverified" label claim.** Searched workspace: the label appears in `FINAL.md`, `consensus.md`, `round-03/` artifacts, but NOT in `COOPERATION.md` or any `.go` file. The ratified design (acceptance criterion 5) requires it in `COOPERATION.md`. The implementer delivered it in artifacts, not in the protocol text. Refutation PARTIAL — label exists, but not where the design requires.
8. **Clause (b) presence.** Read the old `COOPERATION.md` from `protocol-generation-bias-baseline` (PRIMARY, git-object): clause (b) describes `consensus.md` recording correlated agreement and `FINAL.md` stating family relationships. The current file replaces it with part (b) of the new 15.6 (correlated agreement, shorter). The exchange-duty clause from the three-leg design is not in any `.go` file or `COOPERATION.md` copy. Confirmed: removed, not relocated. Refutation FAILED on presence, PARTIAL on impact (withheld is correct per D4 ruling).

## Findings

### [CRITICAL] D1 — One-packet exchange does not match the measured HiddenBench protocol (arXiv:2505.11556v4, PRIMARY)

What is wrong: The ratified design (FINAL.md leg 2) specifies one sealed packet per participant with no Decide stage. The paper that produces the cited +76.3pp (GPT-4.1 3.7% → 80.0%) measures TWO exchange rounds followed by ONE Decide pass, at 4 agents, with agents never told information asymmetry exists. The implementer implemented the ratified one-packet form and flagged the gap; the gap should block the design.

Why it matters: Shipping a truncated version of an intervention while citing that intervention's measured effect is the exact "transfer unverified" problem — compounded by our own hand. The carrier thesis says prose rules decay; a truncated protocol that claims the full protocol's evidence is worse than prose — it is an incorrect claim in the artifact.

Concrete fix: Upgrade to the measured two-round + Decide form (add the second exchange round and the Decide pass in `internal/runner/` before sealing round 2 / `FINAL.md`); OR re-ratify `FINAL.md` to explicitly adopt the one-packet form with corrected (lower, unmeasured) effect claims and a new benchmark. The implementer's procedural refusal (not silently upgrading) was correct; the substance must change. Block on this finding.

Provenance: PRIMARY — `pdftotext -layout` of `arXiv:2505.11556v4`, local file at workspace. Verbatim quotes verified: "4 agents"; "Exchange stage (2 rounds)"; "Decide stage (1 pass)"; Table 7 (0.037 → 0.800 / 0.173 → 0.727 / 0.043 → 0.743); Table 4 (3 +34.8% / 4 +25.0% / 5 +6.4% / 6 +19.7% / 7 +0.6%).

### [CRITICAL] D4 ruling — Clause (b) withheld is the correct call, not an overstep

What is wrong (the charge): Removing a ratified clause from `COOPERATION.md` looks like implementer overstep.

Why withholding is correct: The clause describes a runner stage (seal round 1, collect packets, release simultaneously) that has no code. Printing it would have been a fourth instance of the idea's named defect class (printed rule, no binding surface). The idea's headline finding — "prose is not a carrier" — applies to its own artifacts first.

Concrete fix: Keep clause (b) out of `COOPERATION.md` until `internal/runner/` implements the exchange orchestration; restore it at that time. The byte net of -476 B (vs ratified -237 B) is verified independently and is a real improvement.

Provenance: PRIMARY — file read of current and baseline `COOPERATION.md`; byte count verified (`sed` over git-object + `sed` over file). Confirmed: no code calls an unimplemented exchange stage (`grep -r` over `internal/` for exchange-stage orchestration returns zero hits).

### [MAJOR] D4 — Byte accounting: -476 B verified (better than ratified), but the "transfer unverified" label (AC #5) is missing from `COOPERATION.md`

What is wrong: AC #5 of `FINAL.md` requires `COOPERATION.md` to contain the literal `"transfer unverified"` attached to the exchange. It is present in artifacts (`FINAL.md`, `consensus.md`, `round-03/*`) but absent from `COOPERATION.md` and all `.go` files. If the exchange clause were restored (per D1 upgrade), the label must go with it; as it stands, neither clause nor label is in the protocol text.

Why it matters: The label is the instrumented acknowledgment of the design's own risk (R1); without it in the protocol text, future implementations cannot verify compliance by grep.

Concrete fix: If D1 is resolved by upgrading to the two-round form, embed `"transfer unverified; instrumented"` in the protocol clause (as specified in `FINAL.md` leg 2). If D1 is resolved by re-ratifying the one-packet form, embed the same label with the corrected effect-size claim.

Provenance: PRIMARY — workspace grep (`grep -rn "transfer unverified" .` across repo). Confirmed present in 8 artifact paths, zero in `COOPERATION.md` or `.go`.

### [MAJOR] Byte accounting verification — confirmed, but third-copy drift guard gap confirmed

What is wrong: The implementer says all three copies moved. Only two exist in this workspace. The design's drift guard (`TestDriftAnchorsAcceptTheGeneratedRosterTable`) covers the defaults (`internal/protocol/defaults/COOPERATION.md`, 105415 B) — verified passing — but does not cover `parley-deck/COOPERATION.md` or the missing skill-reference copy. This is the exact skill-vs-protocol split noted in the idea's known risks (skill/template emit vs enforced gate divergence: `responding-to:` 18.1% vs `### @<other>` 7.2%).

Why it matters: A consistent protocol change that does not cover the skill copy means the hand-driven path (skill template) emits artifacts that the Go-enforced path refuses — reproducing the very split the design identifies.

Concrete fix: Either (a) create/update `parley-deck-skill/skills/parley-deck/references/COOPERATION.md` with the same 15.6 change and add it to the drift-guard roster; or (b) document in `COOPERATION.md` itself (or in a README in `internal/protocol/`) that the skill reference is deprecated and covered by the defaults. Option (a) aligns with the idea's carrier thesis; option (b) is acceptable only with an explicit audit trail.

Provenance: PRIMARY — file-system walk (`find . -name 'COOPERATION.md'`); `os.path.getsize()` measurements; `grep` for `DriftAnchorsAccept` in `.go`; `go test ./...` exit 0.

### [CRITICAL] D1 — HiddenBench dose-response (Table 4) undermines the design's self-indictment (R4) and supports the upgrade

What is wrong: `FINAL.md` R4 claims the 6-agent choice is indicted by `+0.6% at 7 agents vs +34.8% at 3`. The table (PRIMARY) is non-monotonic: N=6 (+19.7%) beats N=5 (+6.4%), and only N=4 (58 tasks) is well-powered. The endpoint comparison is directionally correct but overstated; applying it as a self-indictment of the 6-agent design is weak reasoning.

Why it matters: Overstated self-criticism weakens the design's credibility; combined with D1 (the actual, stronger indictment — wrong protocol form), it muddies which claim should block.

Concrete fix: Soften R4's language to: "The endpoints support the direction (larger groups degrade); the middle is non-monotonic and the N=5/6/7 cells have 5–7 tasks each; the only well-powered cell is N=4 (58 tasks), which produced the +76.3pp effect. The 6-agent choice requires its own benchmark." Keep the benchmark claim (88 ideas with participant counts) as the right test. This is a text correction in `FINAL.md`, not a design reversal.

Provenance: PRIMARY — `arXiv:2505.11556v4` (Table 4), verified by local `pdftotext -layout`. Numbers: N=3 (7 tasks, +0.348), N=4 (58 tasks, +0.250), N=5 (5 tasks, +0.064), N=6 (5 tasks, +0.197), N=7 (5 tasks, +0.006).

### [MAJOR] Gate mutation verified — but coverage of the runtime path is partial

What is wrong: The mutation (`len(raw) < 0`) was performed correctly (compile succeeds, never fires), the test `TestValidateRoundOneArtifactRequiresExistingAlternatives` failed (PRIMARY, `go test` output), and the file was restored to green. Confirmed: the gate exists, fails when broken, passes when restored. However, `ValidateRoundArtifact` is called only from two paths (`driver.go:367`, `phase58.go:321`); neither covers the new `## Existing alternatives` requirement directly — the requirement lives in `protocol.ValidateRoundOneArtifact`, not in `runner.ValidateRoundArtifact`. The design's AC #1 ("rejected by the validator, not merely warned") is satisfied by the protocol-level gate; the AC #4 claim ("zero rounds/files/agents added") holds for the exchange stage (it is withheld, not implemented). The gate is real; its integration point is the protocol layer, not fully mirrored at the runner layer.

Concrete fix: Confirm (or add) that `internal/runner/` imports and invokes `protocol.ValidateRoundOneArtifact` explicitly for every design-round output, not only through the generic `ValidateRoundArtifact` wrapper — or document that the wrapper's delegation to `protocol.ValidateRoundOneArtifact` (line 17) is the canonical path and is complete.

Provenance: PRIMARY — `roundartifact.go` file read; mutation test output; `grep` across `.go` for call paths.

### [MAJOR] `opencode-1.md` (excluded participant, filed 14:05 2026-08-29) defines vocabulary the disposition leg requires — leg is incomplete without it adopted

What is wrong: `FINAL.md` defers the vocabulary question (new finding class?) to `opencode-1.md`. `opencode-1.md` proposes `REFRAME` (not a fifth severity), `## Frames considered` destination, freeze gate (`FINAL.md` MAY NOT freeze while any `REFRAME` is absent from `## Frames considered`), four-field payload (current frame / other frame / witness / stay-condition), and Phase 5–8 rules. The implementation has adopted NONE of this vocabulary: no `REFRAME` in `COOPERATION.md`, no `## Frames considered` in any `.md` file, no freeze-gate code, no four-field payload check. The deferred input is therefore live but unadopted.

Why it matters: Without adoption, the deferred question is unresolved, which makes the disposition leg incomplete — precisely the state `opencode-1.md` warns about ("freeze-without-absorption is the protocol's mindguard"). The idea's own evidence says unabsorbed alternatives disappear (B1 / PBS case); not adopting `REFRAME` reproduces that failure in the very mechanism designed to fix it.

Concrete fix: Read `opencode-1.md` (done, PRIMARY). Decide: adopt `REFRAME` + freeze gate + `## Frames considered` (recommended — it is the direct B1 fix), or explicitly reject it with a recorded reason and witness (if the design's `ALT-` vocabulary is preferred over `REFRAME`). Either choice completes the deferred leg. Silence is not a decision.

Provenance: PRIMARY — full read of `opencode-1.md` (11,146 bytes). Confirmed: proposes `REFRAME`, `## Frames considered`, freeze gate, four-field payload, witness requirement, cost budget, Phase 5–8 rules. Confirmed absence of adoption: workspace grep for `REFRAME` shows only references in this file and `opencode-1.md`; `## Frames considered` appears only in `FINAL.md` (AC #5 reference) and `opencode-1.md`.

### [MINOR] Exchange execution prompt (when implemented) must omit asymmetry claim — verified in `FINAL.md`, unverified in code (because code is withheld)

What is wrong: `FINAL.md` leg 2 requires the exchange prompt to contain "no assertion that information asymmetry exists." The design's text respects this; no code implements the prompt yet (withheld per D4). The check is therefore a future verification obligation, not a current finding.

Concrete fix: When D1 is resolved (upgrade to two-round + Decide), include the greppable omission of asymmetry assertions as an AC verification step; the design's AC #4 (grep-checkable against the template) covers it.

Provenance: SECONDARY — `FINAL.md` text (line 152: "no assertion that information asymmetry exists"). Confirmed in `opencode-1.md` as consistent with the paper's protocol.

### [MINOR] Skill/template drift: `SKILL.md` does not exist in workspace; reference third-copy missing

What is wrong (repeat of MAJOR above, lower severity for the skill file specifically): `parley-deck-skill/skills/parley-deck/references/COOPERATION.md` absent; `SKILL.md` itself is only present as `tmp-test-plugin/skills/SKILL.md` (unrelated plugin). The design's own risk-list item 5 (skill/template emit vs enforced gate divergence) cannot be verified without the skill file.

Concrete fix: Confirm whether the skill directory lives outside this repo (e.g. `~/.parley/` or a separate module); if so, the review should reference that external path in provenance rather than treating absence as a failure. The byte and consistency verification covers what is verifiable inside the workspace.

Provenance: SECONDARY — `find` results; `ls` of `parley-deck-skill/` (directory exists, `references/` subdirectory missing); `grep` for SKILL.md in workspace.

### [NIT] R2's pre-implementation gate (group size read) was executed correctly by implementer; recommendation: keep the documented gate as a mandatory pre-commit check

What is wrong: Nothing — the implementer cleared R2 properly (read the paper, confirmed 4 agents, confirmed two rounds + Decide, confirmed non-monotonic middle). The recommendation is procedural reinforcement, not a defect.

Concrete fix: Make the R2 gate (read `arXiv:2505.11556v4`, confirm group size, confirm round structure) a mandatory checklist item in the next design idea that relies on HiddenBench — not optional, not retrospective.

Provenance: PRIMARY — `IMPLEMENTATION.md` gate-1 section; verified independently by re-reading the paper.

## Open questions

1. **D1 resolution path:** Should the design upgrade to the measured two-round + Decide form (recommended — evidence-backed, fixes the defect class), or should it be explicitly re-ratified to one packet with corrected (lower) effect-size claims? The implementer has not chosen; the review blocks until one is adopted. (Status: OPEN — requires design-owner decision, not implementer discretion.)
2. **`REFRAME` adoption:** Does the design adopt `opencode-1.md`'s vocabulary (`REFRAME` + `## Frames considered` + freeze gate) or explicitly reject it with reason and witness? Either completes the deferred vocabulary leg; silence does not. (Status: OPEN — depends on D1 resolution; both paths require an explicit decision in a revised `FINAL.md` or an addendum.)
3. **Skill/reference third copy:** Where does `parley-deck-skill/skills/parley-deck/references/COOPERATION.md` live? Is it generated, deprecated, or maintained externally? The workspace absence prevents full three-copy verification. (Status: OPEN — needs location confirmation from the design owner; does not block the protocol-level finding but affects AC verification completeness.)
4. **D4 final form once exchange stage is implemented:** When `internal/runner/` implements the two-round + Decide exchange, clause (b) must be restored with the `"transfer unverified; instrumented"` label embedded. Should the restored clause use the original `FINAL.md` three-leg wording (full) or the `opencode-1.md` freeze-gate form (with `REFRAME` destination)? The two are compatible but not identical — clarification needed before implementation resumes. (Status: OPEN — depends on answers to #1 and #2.)
5. **R4's benchmark claim:** The 88-idea participant-count comparison is proposed but not executed in this change. Should it be added to the implementation checklist (as a measurement task) or deferred to a separate benchmark idea? (Status: OPEN — not blocking; recommendation is to defer to separate idea so this review remains focused on the protocol-change verification.)
