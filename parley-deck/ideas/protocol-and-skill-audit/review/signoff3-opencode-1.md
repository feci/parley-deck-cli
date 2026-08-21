### Signoff: opencode-1 — 2026-08-21
Status: ✅ ACCEPT
Notes: Cycle-2 block conditions from peers are satisfied at HEAD `39dbc778ba56ec20f531a9da249e0cb018e0a381`. I verified in a disposable clone (not the shared tree): `go test ./...` exit 0 (26 ok + 1 no-test package), `npm test` in `parley-deck-skill` 388/0 exit 0, all four previously undispositioned MAJORs hold with mutation-failing tests, `internal/driver` no longer owns a second section list, the two corrected consensus claims match a three-binary deck sweep, and `addon-manifest-coverage` (track: deliberation) is `partial` again. Prior signoff2 slice (no test weakening) still stands; flipping to accept on this amended consensus.

## Evidence

### 1. `go test ./...` — PRIMARY
- Tree: git clone `--no-hardlinks` of the shared repo → `/var/folders/.../T/opencode/signoff3-v3` at `39dbc778ba56ec20f531a9da249e0cb018e0a381`.
- Command (foreground; exit from the test process only):
  ```
  set +e
  go test ./... > /tmp/gotest-final.out 2>&1
  GO_EXIT=$?
  set -e
  echo "GO_TEST_EXIT_CODE=${GO_EXIT}"
  ```
- Observed: **`GO_TEST_EXIT_CODE=0`**. Summary: 26 `ok` lines, 0 `FAIL` lines, 1 `?` (`cmd/parley`, no test files). Packages include `internal/app` (53.5s), `internal/consensus`, `internal/driver`, `internal/runner` — no hang.
- How exit was read: `$?` of the `go test` process in the same foreground shell immediately after it returned (not a background job, not the last pipeline stage of something else).

### 2. `npm test` in `parley-deck-skill/` — PRIMARY
- Command: `cd "/Volumes/My Shared Files/AI_WORKSPACE/parley-deck/parley-deck-skill" && npm test > /tmp/npmtest-signoff3.out 2>&1; echo NPM_TEST_EXIT_CODE=$?`
- Observed: **`NPM_TEST_EXIT_CODE=0`**. Tail: `tests 388`, `pass 388`, `fail 0`; python adapter suites OK; packaged skills hash checks OK.

### 3. Four dispositioned MAJORs — genuinely fixed

**3a. Manual `consensus finalize` rejects FINAL with wrong `idea:` / `status:`** — PRIMARY + SECONDARY  
- SECONDARY: `internal/consensus/consensus.go` Finalize path calls `protocol.ValidateFinal(string(existing), idea.Slug)` and returns `cannot close this idea: …` when non-empty (`4c43200`).  
- SECONDARY: `protocol.ValidateFinal` requires `status: final` and matching `idea:` before content scaffold checks (`internal/protocol/finalsections.go`).  
- PRIMARY: `go test ./internal/driver/ -run 'TestFinalDeclaringAnotherIdeaIsRejected|TestFinalWithNoSlugIsRejected|TestFinalScaffoldReason' -count=1 -v` → all PASS (gate rejects wrong idea / missing slug / non-final status via the shared validator).  
- PRIMARY mutation: stripped status+slug checks from `ValidateFinal` (content-only) → `TestFinalScaffoldReason` FAIL: `non-final status should be rejected`; restore → PASS.

**3b. `consensus draft --review` rejects artifacts with no `reviewed-commit`** — PRIMARY  
- PRIMARY: `go test ./internal/consensus/ -run TestManualReviewDraftRejectsAnArtifactWithNoReviewedCommit -count=1` → PASS at HEAD.  
- PRIMARY mutation: removed the empty-`reviewed-commit` branch in `protocol.ValidateReviewArtifact` → same test FAIL: `manual review draft accepted a review with no reviewed-commit`; restore → PASS.

**3c. Deliberation review quorum awaits implementer; standard/fast do not** — PRIMARY  
- PRIMARY:  
  - `TestDeliberationReviewConsensusStillAwaitsTheImplementersSignoff` PASS  
  - `TestStandardReviewConsensusDoesNotAwaitTheImplementersSignoff` PASS  
  - `TestAbsentTrackReviewConsensusUsesTheStandardQuorum` PASS  
- PRIMARY mutation: `reviewConsensusVoters` always returns `expectedRoundParticipants` (pre-fix) → deliberation test FAIL: `does not await the implementer: missing=[] triage=ready`; restore → PASS.  
- PRIMARY live: `parley consensus status --review --json addon-manifest-coverage` → `triage: "partial"`, `missing: ["claude-1"]` (implementer), `track: deliberation` in `00-prompt.md`.

**3d. `preflight --yes` clears `unknown-freshness` on a fresh deck** — PRIMARY  
- PRIMARY: `go test ./internal/app/ -run TestConfirmedFreshnessClearsTheGateOnAFreshDeck -count=1` → PASS.  
- PRIMARY mutation: deleted the `if opts.Yes { return confirmProtocolHashes(...) }` branch → FAIL: `--yes left a gate standing … Kind:unknown-freshness`; restore → PASS.

### 4. Mutation discipline (at least one fix reverted, test at HEAD) — PRIMARY
Performed in the clone for deliberation quorum, `reviewed-commit` draft gate, preflight `--yes`, and `ValidateFinal` status/slug — each FAIL under mutation and PASS after restore (commands above). No shared-tree writes.

### 5. Driver no longer owns a second section list — PRIMARY
- Command: `rg -n "requiredFinalSections|missingFinalSections" internal/driver` → no matches (`driver clean`).  
- SECONDARY: `finalScaffoldReason` is a thin wrapper over `protocol.ValidateFinal` (`internal/driver/consensus.go:167-180`). Canonical list is only `protocol.RequiredFinalSections` / `MissingFinalSections`.

### 6. Corrected claims in `review/consensus.md` — PRIMARY
Three binaries built from the same clone (`go build -o … ./cmd/parley` at `a1926ae`, `0bb9903`, `39dbc77`), then `consensus status --review --json` over all **66** review consensuses on the HEAD deck tree:

| comparison | flips | shape |
| --- | --- | --- |
| reviewed `0bb9903` vs base `a1926ae` | **30** | 6 `partial→ready`, **24 `→malformed`** (17 ready→malformed + 7 partial→malformed) |
| HEAD vs base | **5** | **5 `partial→ready`**, 0 malformed flips |

- Matches the consensus table as written.  
- `addon-manifest-coverage`: base=`partial`, reviewed=`ready`, HEAD=`partial` (deliberation quorum restored).  
- Anti-drift: PRIMARY `TestAFinalBuiltFromThePromptSatisfiesTheProductionGate` exists in `internal/app/finalprompt_test.go` and is in the green `go test ./...` run; it builds a FINAL from prompt headings only and feeds `protocol.ValidateFinal`.

### 7. `addon-manifest-coverage` partial at HEAD — PRIMARY
- `/tmp/parley-head consensus status --review --json --dir <clone> addon-manifest-coverage` → `"triage": "partial"`, missing implementer `claude-1`.  
- SECONDARY: `00-prompt.md` has `track: deliberation`, `IMPLEMENTATION.md` `implementer: claude-1`.

### Commits read — SECONDARY
`7112e03` (pipeline hang fixtures), `4c43200` (four MAJORs), `1f3d971` (prompt-to-gate contract test), `39dbc77` (consensus + IMPLEMENTATION amendments). Reviewed commit matches `git rev-parse HEAD` at verification start: `39dbc778ba56ec20f531a9da249e0cb018e0a381`.
