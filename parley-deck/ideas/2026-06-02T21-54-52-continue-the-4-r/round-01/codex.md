---
agent: codex
idea: 2026-06-02T21-54-52-continue-the-4-r
round: 1
date: 2026-06-02
---

## Summary

All four follow-ups are feasible as additive changes except Item 3, which is infra-blocked. The implementation path should preserve the shipped 1.9.0 linear surface by defaulting new cursor fields to empty, keeping `current_block` populated, and treating production gates as non-bypassable in every new DAG and auto-drive path.

Recommended build order:

1. Item 4: agy launch robustness, because it stabilizes agent artifact creation for the rest of the pipeline work.
2. Item 2: unattended Phase 8 loop, because it is mostly runner/review state and can be tested with fake agents before multi-active execution.
3. Item 1: parallel multi-active DAG auto-drive, because it touches cursor semantics, status, scheduling, and gate interactions.
4. Item 3: release infra step only; upstream WinGet PR waits for real immutable assets and hashes.

## Proposed approach

Item 1, parallel multi-active DAG auto-drive:

Assessment: extend the DAG cursor without changing the linear contract. Add `ready_blocks []string` and `active_blocks []ActiveBlock` or `[]string` to `pipeline-run.json`, with zero-value defaulting during load. Keep `current_block` set to the first active block, or first ready block if none are active, in deterministic DAG order for existing status and linear commands.

Minimal implementation: make `advanceDAG` compute all ready blocks whose dependencies are complete and whose inbound gates are resolved. Sort by manifest block order, then ID as a stable tie-breaker. `pipeline auto` should launch all ready blocks up to a configurable cap, defaulting to unbounded or a conservative `PARLEY_PIPELINE_MAX_CONCURRENCY` default of CPU count. Track per-block phase/round state under each block's existing artifacts, not in a shared mutable consensus field. Add a new status view that shows ready, active, blocked-by-gate, complete.

Trade-off: unbounded concurrency maximizes throughput but can oversubscribe agents and produce noisy failures; a default cap is safer but means "all ready" means all ready within capacity. I would expose the cap but keep deterministic scheduling.

Item 2, fully-unattended Phase 8 fix-up loop:

Assessment: the loop needs a small machine contract in `review/consensus.md`, otherwise auto cannot distinguish "no fixes" from "not parsed". The drafter should own the contract because it already synthesizes reviews into consensus.

Minimal implementation: require frontmatter in `review/consensus.md` with `outstanding_agreed_fixes: <int>` and optional `blocked: true|false`. Add `ReviewAgreedFixes(path) (count int, blocked bool, err error)`. In auto mode only, run review round N, draft consensus, request signoffs, stop on BLOCK, stop on zero fixes, otherwise run `RunFixup` by the implementer and then start review round N+1. Cap cycles with a manifest/run option defaulting to 3.

Trade-off: frontmatter is easier to parse and test than fenced YAML but constrains consensus documents. "Applied" should mean the implementer produced its fix-up artifact and the next review consensus reports zero outstanding agreed fixes, not that the driver infers code correctness.

Item 3, WinGet upstream PR:

Assessment: infra-blocked. There are no GitHub releases for CLI v1.6.0, v1.7.0, v1.8.0, or v1.9.0; only git tags exist. The latest GitHub release is v1.5.4. The portable `.exe` assets are built by CI on release publication, so real `InstallerSha256` values do not exist yet.

Minimal implementation: do not open a microsoft/winget-pkgs PR and do not invent hashes. The safe actionable step is to create GitHub releases from the existing tags, for example with `gh release create` using the existing tag metadata, wait for portable-build CI to attach immutable `.exe` assets, compute `InstallerSha256` from those assets, validate the WinGet manifest locally, then open the upstream PR.

Trade-off: creating releases now unblocks packaging but may require release notes/signoff discipline. Waiting avoids accidental public release metadata churn but leaves WinGet blocked.

Item 4, agy headless launch regression:

Assessment: agy 1.0.4 can write artifacts in `--print` mode when given a directive prompt and scoped `--add-dir`; the observed `parley run` failure likely comes from broad repo exploration under a large root plus a long idea prompt. The fix should be launch-spec level and generic enough to avoid hard vendor coupling.

Minimal implementation: add a generic write-first preamble for headless file-producing agents: immediately create the exact requested artifact path, include required frontmatter, avoid repository exploration unless explicitly necessary, and report blocker instead of narrating. Keep an agy launch-spec note only for its known `--print` behavior. Add stdout-capture fallback after process exit: if the artifact is missing and stdout parses as markdown starting with `---`, write stdout to the expected artifact path as agent-authored content and record the fallback in driver logs. Test with fake agents that either write the file, print valid markdown, print invalid markdown, or narrate without frontmatter.

Trade-off: stdout fallback slightly broadens "agents write markdown", but the agent still authored the artifact bytes and the driver only places them at the mandated path. The fallback must reject invalid frontmatter to avoid converting narration into protocol artifacts.

## Concerns / open questions

For Item 1, the main open question is whether concurrency should default to all ready blocks or a bounded cap. I recommend a cap with deterministic scheduling because it is easier to test and safer for real agent quotas.

For Item 2, consensus frontmatter should be mandatory only for auto-driven implementation blocks. Manual review consensus can remain human-readable, but auto must fail closed if the machine fields are absent.

For Item 3, the human decision is whether creating GitHub releases from existing tags is acceptable now. My recommendation is yes, if release ownership approves public release creation; the upstream WinGet PR should wait.

For Item 4, stdout fallback should be generic and guarded by strict validation, while agy-specific behavior should stay in launch discovery/spec metadata rather than scattered runner branches.

## Risks

Item 1 risks race conditions around shared `pipeline-run.json`, especially if multiple active blocks finish simultaneously. Use atomic cursor writes and a single driver-owned scheduler loop.

Item 2 risks infinite or low-value fix-up churn. The max-cycle cap, BLOCK signoff stop, and required `outstanding_agreed_fixes` field are the safety controls.

Item 3 risks shipping incorrect WinGet manifests if hashes are guessed or assets are mutable. Do not proceed without real release assets and verified hashes.

Item 4 risks masking a broken agent by accepting narrative stdout. Strict `---` frontmatter validation and required protocol shape keep this fallback narrow.
