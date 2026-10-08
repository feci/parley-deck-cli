# Review-01 validation observations — codex-1

Candidate CLI e4681cf5b9144ed786b269d8ab14b8d65b23d581; skill efe296c7acf13a147ab820ce6cbf8e6705b68691. These are implementer observations, not independent verdicts. Private full logs remain in `.parley-runtime/participant-dropout-launches/`. No acceptance is inferred from an unfinished process.

- Full host `go test ./... -json -count=1 -timeout 45m` remains running in `full-host-tests-final.jsonl`. One observed failure: `TestEvidenceVerifierProductionClosure/report-persistence-failure-recovery` at evidence_verify_test.go:292, “fresh real verification did not recover: escalated independent completion evidence refused: independent verifier process failed or its launch was not observed; refusal committed; fresh checks are required”.
- Exact focused rerun `go test ./internal/app -run 'TestEvidenceVerifierProductionClosure/report-persistence-failure-recovery' -count=3 -timeout 5m` passed (13.590s), `evidence-recovery-rerun.log`. Cause is not established.
- Full skill `npm test` had 398/399 Node tests pass; its fleet-wide immovable-destination case failed at bidding-addon.test.js:2480 during the initial installation, before the immutability fixture. Python/manifests did not run after the failed Node leg. `skill-tests-final.log`.
- Exact focused skill rerun `node --test --test-name-pattern='install is fleet-wide too' test/bidding-addon.test.js` passed 1/1, 14.072s (`skill-fleet-rerun.log`). Cause is not established. A fresh unmodified full `npm test` runs in `skill-tests-rerun.log`.
- CLI PR CI run 37854001955 Linux job 113573680699 failed `TestDropoutExecWatchdogTimeoutAndACPShareTwoSlots/acp-started-exit`: dropout_test.go:104 expected process_failure; printed pointer-only failure class. Whole subtest elapsed 0.31s with two attempts using a 150ms ceiling. This may be a timeout-sensitive fixture but that is a hypothesis. `ci-linux-pr-failure.log`; focused repetition in `acp-exit-repeat.log`. The push CI run 37853995105 Linux job passed. Neither result cancels the other.
- Skill push and PR CI both pass Python 3.10/3.13 matrix jobs (runs 37853996652 and 37854130747). Native macOS suite remains required.
- Final-candidate vet/build passed (`vet-candidate.log` and parley-1.52.0 binary). Focused behavior and protocol checks passed as documented in IMPLEMENTATION.md, but independent refutation and current-tree full validation remain open.

No product or test edits were made during review pending its signed fix consensus. Re-runs retain failed evidence and do not relabel it as a pass.

## Completed follow-ups

- Unchanged full skill rerun `npm test` passed: all 399 Node tests, 54 Python tests across seven files, and all six payload manifests. `skill-tests-rerun.log`. No source/test change or dependency version change between failure and pass. Initial fleet-install cause remains unestablished.
- ACP exact focused test repeated 30 times passed in 12.989s (`acp-exit-repeat.log`). Proposed review disposition: give the started-exit fixture a realistic non-timeout ceiling, retaining the separate deliberate 150ms timeout case and original-plus-retry equality; improve diagnostics to print the actual class. This is a test proposal, not an applied fix or a product verdict.
- Full host run finished exit 1: 33 package passes, three no-test packages, one failed app package; only the already reported recovery subtest failed. Fresh full host rerun uses the unchanged tree in `full-host-tests-review01-rerun.jsonl`. A private Go overlay adds diagnostics only for a three-repeat run of the complete evidence verifier test; it modifies no product/test file on disk. `evidence-diagnostic-repeat.log`.

- Diagnostic-only overlay repetition of all 19 verifier cases completed 3/3, exit 0, 156.779s. No diagnostic failure was captured. This is supplemental investigation, not substitution for the unmodified full-host check.
