---
idea: roster-operations-standard
phase: 7 — review consensus signoff
agent: codex-1
date: 2026-08-06
verdict: BLOCK
---
# Signoff — codex-1
## Verdict
BLOCK — complete the missing finding dispositions, correct §4, supply DF-1's evidence locator, and repair the Phase-7 frontmatter.

## Answers to 1-6

1. **PRIMARY — No.** The two original CRITICAL bases are acknowledged as fixed in cycle 1 and their residuals map to A1-A5/A10/DF-2 (`review/consensus.md:29-33,74-150,196-202,259-262`); the command/JSON, normalizer, protocol-text, and G5 findings map to A6-A9/A11/A13-A15 (`:152-243`). Two of mine are not disposed accurately:
   - **PRIMARY** — The machine-scope `PARLEY_HOME` finding (`review/round-01/codex-1.md:261-277`) is absent from §§2-3. Current code calls `config.CentralAgentsPath()` and explains the former wrong path (`internal/app/roster_set.go:89-103`), with the regression test at `internal/app/roster_sync_test.go:13-20,82-114`; record it explicitly as fixed in cycle 1.
   - **PRIMARY** — My sync hardening finding named unmatched `--keep` tokens and a reread-without-CAS window (`review/round-01/codex-1.md:279-296`). DF-3 mentions only backup/cleanliness and says preview plus `--keep` is sufficient (`review/consensus.md:263-265`), so those two sub-findings were dropped/misrepresented. Agree, defer, or dismiss them with a reason.
   - **PRIMARY** — A13 downgrades the docs half of my composite MAJOR to MINOR while retaining help discoverability as MAJOR (`review/consensus.md:215-221`). I accept that split; it is not a blocker.

2. **PRIMARY** — I re-read the sources. §7 requires the heading plus `Idea:`, `Drafted by:`, and `Summary:` (`COOPERATION.md:741-749`); the entry instead uses bold `Idea` with a slug, `Change`/`Why`, and no `Drafted by:` (`meta/protocol-changelog.md:119-139`). I support C-1 and A11. The resolution follows the normative template, not reviewer count (`review/consensus.md:35-59`). This is evidence and a signoff position, not a §15 truth verdict on the G5 finding I own.

3. **PRIMARY** — A1's resolution is the right interpretation: the deck file owns membership, while model/effort/speed may inherit values (`consensus.md:353-360,374-378`; `FINAL.md:79-82`). Refusing to render inherited machine membership into committed §2 without an explicit flag is correctly scoped. Add one sentence that a valid legacy §2 roster remains the deck's compatibility membership until migration; “no roster of its own” must mean neither deck TOML nor valid legacy §2 (`COOPERATION.md:124-125`; `00-prompt.md:77-78`).

4. **SECONDARY** — `review/consensus.md:251-258` records 24 applied, 9 skipped, 3 unchanged, 0 failed, with backups and validation/rollback; I did not inspect the fleet report. **PRIMARY** — The ratified contract required CAS, inventory/version/worktree gates, small-batch confirmation, and full field carry (`consensus.md:444-468`), matching kimi-1 M10 (`review/round-01/kimi-1.md:261-286`). Deferring the remaining hardening is acceptable after the completed attended operation if dirty-tree refusal plus `--confirm-breaking` ships before reuse and the remainder gets a named follow-up. Because the fleet outcome is material to that deferral, add the stable report path or quoted command/output: the draft's default `PRIMARY` assertion (`review/consensus.md:17-20,251-258`) currently lacks the locator §15.2 requires (`COOPERATION.md:1240-1248`).

5. **PRIMARY** — Something should be dismissed. Hermes-1's “three copies are not identical in §2” finding expressly says the differing roster bodies are an intentionally normalized project-specific zone and only underscored the then-missing generator (`review/round-01/hermes-1.md:276-297`). With the generator present, that finding is resolved/not-an-issue; it does not corroborate A9's separate stale-instruction defect. Move it to §4 with that rationale. DF-4 may remain deferred as an optional design follow-up.

6. **PRIMARY** — The reviews predate this draft and therefore missed its malformed Phase-7 frontmatter: the protocol requires `drafted-by:` and `reviewed-commit:` (`COOPERATION.md:554-564`), while the draft has `drafter:` and only a free-form `baseline:` (`review/consensus.md:1-10`). Fix both in this cycle. I found no additional source-code blocker beyond the dispositions above.

## Conditions (if any)

1. Add explicit dispositions for the machine-scope fix and both sync-hardening sub-findings; retain the accepted A13 severity split.
2. Clarify A1's legacy fallback and move hermes-1's allowlisted-copy finding to Dismissed.
3. Cite the fleet report/command evidence for DF-1 and replace `drafter`/`baseline` with the required Phase-7 provenance fields.
4. Issue the corrected draft for a new signoff.
