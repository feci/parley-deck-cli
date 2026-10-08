---
agent: zcode-1
idea: meta-protocol-change-participant-dropout
kind: goal-check
date: 2026-10-09
goal-check-commit: 9b6edba81149a07d9a2372d118ec5477f76b7b4c
product-commit: cb78e9f43e70cb59284a21ba6e0e8a331cb447a1
reviewed-commit: f8f4f1f17bc99832aeb46c0a8ee44c6331325d1b
skill-commit: efe296c7acf13a147ab820ce6cbf8e6705b68691
---

# Goal-done check — zcode-1 (independent non-implementer, LE-7/LE-11)

## Protocol context attestation

`context_mode=full`,
`source_sha256=packet_sha256=091e6fb841685c85fa153f7e2f05e2329c3f28dbd9b0b4c72831bc88dc458bbf`,
`fallback_reason` absent. Verified by `shasum -a 256` on
`.parley-runtime/protocol-packets/full-phase-1-deliberation-091e6fb841685c85fa153f7e2f05e2329c3f28dbd9b0b4c72831bc88dc458bbf.md`
(PRIMARY, executed 2026-10-09); the hash equals the live launch attestation's declared
source/packet hash. The launch's shadow-packet audit line is an unapplied diagnostic and does
not describe the supplied full protocol. Active transport: `github-pr`.

## Product identity at this check

- HEAD `9b6edba` is deck-only over `f8f4f1f` (IMPLEMENTATION.md, source-context/validation-cycle1.md,
  one inbox note): `git diff --stat cb78e9f..9b6edba -- . ':(exclude)parley-deck'` is empty.
- Working tree product paths are clean: `git status --porcelain -- . ':(exclude)parley-deck'` empty.
- Therefore the product at this check (`cb78e9f` content) is byte-identical to the reviewed
  `f8f4f1f` tree and to the tree my round-02 probes ran against (archive of `f8f4f1f`).
- Skill unchanged at `efe296c`.

## Evidence inspected (commands executed by me today unless stated)

1. Archived fix-consensus integrity: `shasum -a 256 review/consensus-cycle-01.md` =
   `4db5a60edc68c89a3194328200e8c973363a7bb53e74d9618a8071c88a7a07b3` — byte-identical archive of
   the signed cycle-1 fix consensus; both signoffs (codex-1 ✅, zcode-1 ✅, 2026-10-09) retained
   verbatim. No signoff content was edited; final consensus is not yet created and I do not sign
   it here.
2. Qualifying full-host run: `full-host-tests-cycle1-local-result.json` — exit 0 at `f8f4f1f`,
   731.06 s, `go test ./... -json -count=1 -timeout 45m`, local TMPDIR/GOTMPDIR. My independent
   recount of the raw `full-host-tests-cycle1-local.jsonl` (python3 over the event stream): 37
   packages = 34 pass + 3 no-test skip, 0 package FAILs; 3262 test-level pass events, 0 fail
   events, 4 skips. Log SHA256 `131f8ef7c53971e1cedcfa0402db7c118cd7244fbd764425e22c26deeb37e503`
   matches `source-context/validation-cycle1.md`'s claim.
3. Failed external run preserved as history, not evidence: `full-host-tests-cycle1-result.json`
   (exit 1) with 7 test failures in internal/budget ×2, internal/driver ×1, internal/runner ×1,
   internal/app ×3 — all packages untouched by this idea's diff; consistent with my round-02
   adjudication (external-mount/git/lock/path-space environment breakage). It stays failed
   history; nothing relabels it as a pass.
4. Fresh current-tree spot checks (local `TMPDIR=GOTMPDIR=/private/var/tmp/zc-goalcheck-20261009`,
   warm external GOCACHE): `go build ./cmd/parley` → `parley 1.52.0`; `go vet ./internal/runner/
   ./internal/app/` → PASS; `go test ./internal/runner/ -run TestDropout -count=1` → ok 5.873 s;
   `go test ./internal/app/ -run 'TestDropout|TestGoalCheck|TestReadiness' -count=1` → ok 33.318 s.
   I did not rerun the full suite: the qualifying exit-0 full-host run already covers it at a
   git-verified-identical product tree, and no product change has occurred since (no new evidence
   warrants a rerun).
5. Artifact-format repair of my own review: `review/round-02/zcode-1.md` heading normalized to the
   exact `## Refutation attempts` (scope text moved into a paragraph; nothing else changed —
   trailing correction comment in that file). `parley wait --dir . --idea
   meta-protocol-change-participant-dropout --for review --timeout 5s --json` now returns the
   artifact `valid: true` and exits 0 at the review-round boundary.

## AC assessment (AC1–AC13; AC14 out of scope here)

- AC1–AC12: evidenced by my independent round-02 verification (source reading, committed tests,
  four adversarial probes, all PASS) on the `f8f4f1f` archive tree, whose product content is
  byte-identical to current HEAD (verified today by git, item 1 of identity); the code paths are
  unchanged since `cb78e9f`. Today's focused reruns (item 4) reconfirm the dropout/goal-check/
  readiness surfaces on the live tree. Detailed per-AC attempts: `review/round-02/zcode-1.md`.
- AC13: met — full-host `go test ./... -count=1` exit 0 at `f8f4f1f` (34 packages pass / 3
  no-test / 3262 test-pass events, recounted by me), vet/build exit 0 (producer's
  `vet-build-cycle1-result.json` at `cb78e9f` plus my own fresh build/vet), and the skill suite
  (my round-02 full rerun on unchanged `efe296c`: 399 Node + 54 Python + 6 manifests, zero
  failures). No open CRITICAL/MAJOR anywhere in the review lifecycle: round-02 requests zero
  agreed fixes; Z1/Z3 verified fixed, Z2 withdrawn with commitments landed, Z4 a recorded
  deferral. Windows limits disclosed, not passed.
- AC14: explicitly post-close delivery (releases, formulae, WinGet, installs, channel
  verification, core staging) — not started, not claimed, and not assessed here.

## Limits

No real provider 429/503/400 was exercised against the new trigger (evidence is local supervisor
facts by design; the test matrix fabricates the classes). Cross-host/PID-namespace crash
recovery, real hosted error timing, and native Windows remain disclosed untested surfaces, as
FINAL states. My probes ran on the archived, git-verified-identical tree, and today's live-tree
checks were the focused suites above, not a fresh full suite. Per LE-7/LE-11 this check can only
withhold a close, never establish one: an attended close additionally requires the final review
consensus (zero Agreed fixes) to be drafted and signed by both codex-1 and zcode-1, and a textual
pass never replaces that current-tree review evidence.

## Verdict

**GOAL-CHECK: PASS** — implementation acceptance criteria AC1–AC13 are met on the current
product tree (HEAD `9b6edba`, product `cb78e9f`), ready for attended close AFTER the final review
consensus is drafted and signed by both participants. AC14 remains post-close owner+implementer
delivery work.
