### Signoff: codex-1 — 2026-08-21
Status: ❌ BLOCK
Notes: I verified the fix-up at HEAD `a27a3b619d24bb42d518dbc30b0bbfb2eb04c892` in an isolated copy. The narrow behaviors covered by the three new focused tests pass, but the review consensus is not safe to close: one fix still miscomputes deliberation quorum, the FINAL prompt and production gate still have independent authorities, the pipeline change regresses an existing integration test, and four MAJOR findings from my filed review have no disposition in the consensus.

#### 1. KNOWN versus REQUIRED: the split works, but REQUIRED is still wrong for a live deck

PRIMARY — `validateDocumentAwaiting` now correctly accepts signoffs from every `known` participant while deriving `missing` only from `required` (`internal/consensus/consensus.go:445-503`). The focused implementer-signoff regression test passes.

PRIMARY — its caller still obtains `required` from `expectedRoundParticipants`, which always removes the implementer for every review and never reads `track:` (`internal/consensus/consensus.go:649-666`). That contradicts the authoritative Phase-7 table: fast awaits one reviewer, standard awaits the reviewers who reviewed, and deliberation awaits all participants (`parley-deck/COOPERATION.md:227-241`). Phase-6 artifact authors and Phase-7 voters are different sets.

There is a concrete current-deck counterexample: `addon-manifest-coverage` is `track: deliberation`, has participants `[claude-1, codex-1, hermes-1, kimi-1]`, and records `claude-1` as implementer. Its review consensus contains only the other three signoffs. The pre-F3 binary from `a1926ae2` reports `partial`, missing `claude-1`; both the reviewed binary from `0bb9903` and HEAD report `ready`. Therefore the sweep's six `partial → ready` flips are not all intended: at least this flip weakens a deliberation gate. The sentence “zero regressions” in `review/consensus.md:27-29` is false.

#### 2. FINAL drafting prompt: the text improved, but the anti-drift claim does not hold

PRIMARY — `buildFinalDraftPrompt` does iterate `protocol.RequiredFinalSections` and does name the directory slug (`internal/app/driver_consensus.go:133-166`). That part of the attributed @codex-1 finding was repaired.

PRIMARY — the production gate does not use that authority. `internal/driver/consensus.go:162-183,201-258` still owns a second `requiredFinalSections` list and a second `missingFinalSections` implementation. In a copy I added one heading only to the production list: `TestFinalDraftPromptDescribesWhatTheGateRequires` still passed, while a FINAL containing every `protocol.RequiredFinalSections` heading was rejected with `FINAL.md is missing required section(s): Audit-only production-gate section`. Thus the prompt can drift from the gate today, and `review/consensus.md:35-36` overstates the fix.

This also leaves my filed MAJOR about manual `consensus finalize` unresolved: that path calls only the deliberately content-only `protocol.FinalIsScaffold` and still accepts a substantive FINAL with the wrong `idea:` and non-final `status:` (`internal/consensus/consensus.go:257-282`).

#### 3. Pipeline scaffold completion: runtime predicate fixed, integration suite regressed

PRIMARY — `isFinalized` now combines `status: final` with `protocol.FinalIsScaffold` (`internal/app/pipeline_cmd.go:1335-1353`), and the focused scaffold test passes. `pipeline auto` also stops after writing the scaffold and explicitly reports the block NOT complete. The @codex-1/@kimi-1 attribution and this narrow behavior in `review/consensus.md:38-42` are correct.

PRIMARY — the existing `startAndFinalize` integration helper still writes `---\nstatus: final\n---\n\ndone\n` as a supposedly pre-finalized artifact (`internal/app/pipeline_cmd_test.go:36-54`). With a deterministic PATH containing no installed agents, `TestPipelineAutoWalksToDoneUnderAutoLeft` passes at reviewed commit `0bb9903` and fails at HEAD because the new predicate correctly rejects that fixture and the test unexpectedly enters real participant selection. Consequently the current Go suite is not green after this fix-up; a full run attempted hosted-agent launches and did not complete. This is a test/fixture regression, not a reason to revert the stricter predicate.

#### 4. Attribution audit and missing dispositions

I audited every sentence in `review/consensus.md` that names @codex-1:

- Lines 31-36 correctly attribute my automatic-drafter MAJOR and accurately describe the old mismatch, but “cannot drift” is false for the reason above.
- Lines 38-42 correctly attribute the pipeline MAJOR jointly and accurately describe the new runtime behavior, but omit the resulting integration regression.
- Lines 79-80 correctly attribute F6, F8, and F14 to me; I independently rechecked their recorded deferral reasons and uphold them.

The consensus does not disposition four other MAJOR findings actually filed in `review/round-01/codex-1.md`: manual finalization accepts wrong slug/status; manual review consensus accepts artifacts without `reviewed-commit`; deliberation review quorum drops the implementer; and a fresh initialized deck's `unknown-freshness` gate cannot be cleared by its displayed `--yes` command. Source inspection and the original focused reproductions still support all four. Phase 7 requires findings to be agreed, deferred, or dismissed; silence is not a disposition.

#### 5. Dismissal and status repairs

PRIMARY — @hermes-1 Finding 1 was correctly dismissed as filed. At reviewed commit `0bb9903`, all three named `00-prompt.md` files have `status: final`; the report used FINAL.md values instead, and `protocol-overlay-local-extension/FINAL.md` itself also has `status: final` rather than `**`.

PRIMARY — the separate, valid observation was handled at the frontmatter level: the other two FINAL.md files now say `status: final`. The wording must remain narrow, however. Running the current production `finalScaffoldReason` on both repaired files still rejects each because it lacks `## Final plan / specification`. Their status defect is repaired; the broader statement that the current gate rejects the files remains true.

#### 6. Deferred findings

PRIMARY — none of the five stated deferral reasons is wrong:

- codex-1/F6 needs a designed semantic signal for an adversarial alternative; substring inference would not safely enforce §15.6.
- codex-1/F8 is a missing collapsed fast-track close path, not a local turn correction.
- codex-1/F14 needs per-knob precedence: blindly applying standard defaults overwrites explicit idea configuration (the documented `5 → 2` case remains real).
- kimi-1/F1 needs a manifest for the installed core shape, which includes package-root payload absent from the source skill's add-on manifest.
- kimi-1/F5 legitimately shares F1's prerequisite because the core's schema-2 marker cannot truthfully point at a nonexistent installed-shape manifest.

Counter-proposal (required if ❌): (1) Replace `expectedRoundParticipants` as the Phase-7 voter source with a track-aware function: deliberation requires all active participants; standard requires the validated review authors (with the two-reviewer/model-diversity rule enforced separately); fast requires its one validated reviewer. Add explicit tests for all three tracks plus `addon-manifest-coverage`. (2) Put one complete FINAL validator in `internal/protocol`, accepting the expected slug and checking status, slug, sections, placeholders, and content floor; use it from the driver, manual finalize, and pipeline completion. Delete the driver-local section list and add a prompt-output-to-production-gate contract test. (3) Update every pre-finalized pipeline fixture to construct a real seven-section FINAL through one shared test helper, keep agent discovery stubbed, and make `go test ./...` pass without launching hosted agents. (4) Apply or explicitly defer/dismiss, with evidence, the four omitted MAJORs above; in particular validate `reviewed-commit` on manual review drafting and make fresh-init confirmation actually persist truthful hashes. Then revise `review/consensus.md`, correct its “zero regressions” and “cannot drift” claims, open review round 2, and obtain fresh signoffs.
