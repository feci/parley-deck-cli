---
agent: agy
idea: antigravity-agent-migration
review-round: 2
date: 2026-05-27
reviewed-commit: 697cb66
---

# Summary

I have reviewed the fix-up cycle 1 changes for the `antigravity-agent-migration` idea. The changes in the current diff were verified against the agreed fixes in [consensus.md](file:///Volumes/My%20Shared%20Files/AI_WORKSPACE/parley-deck/parley-deck-cli/parley-deck/ideas/antigravity-agent-migration/review/consensus.md). 

All five agreed fixes have been successfully implemented and verified:
1. The `installedAgentIDs` function and TUI default participant logic are correctly scoped to only exclude the legacy `gemini` agent, without filtering out other installed headless or ACP-capable agents.
2. Hermes `--yolo` was completely removed from the built-in catalog in `internal/agents/discover.go` and the shared default in `parley-deck/agents.toml`, retaining the standard `--accept-hooks` policy.
3. Documentation in `docs/agent-runtime-configuration.md` has been updated to note that `agy` does not enable home isolation by default, instructing users to provide a custom `isolated_home_env` mapping if they do so.
4. A code comment has been added in `internal/agents/discover.go` noting that the `{prompt}` argument for `agy` is intentionally positioned last because `--print` is a value-taking flag.
5. The `gemini` host mapping rows in both `internal/protocol/defaults/COOPERATION.md` and `parley-deck/COOPERATION.md` have been updated to explicitly label the legacy status.

All unit tests compile and pass successfully. No findings remain.

# Findings

No findings remain. All previously identified issues have been resolved as agreed.
