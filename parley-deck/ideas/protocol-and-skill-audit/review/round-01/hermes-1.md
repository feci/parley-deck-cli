---
agent: hermes-1
idea: protocol-and-skill-audit
review-round: 1
reviewed-commit: 0bb99031aa1ddc1027605424848daa3a696e9f00
date: 2026-08-21
---
## Summary
## Refutation attempts
## Findings
## Open questions

## Finding 1 (Q1 - PRIMARY measurement, python3 /tmp/q1_check.py against repo)
- FINAL.md count: 78
- Non-terminal status with FINAL.md present: 3 (not 0; fix did not fully close).
  - tui-action-execution: status=accepted
  - protocol-overlay-local-extension: status=** (malformed line)
  - meta-protocol-change-rho-retrospective-optimization: status=final-design-for-review
- Clean result for the other 75: status=final / terminal. Repair (zcode-1/F14) covered 20 closed ideas; these 3 anomalies remain outside that scope or have malformed status declarations.

## Finding 1 (Q1 - PRIMARY measurement via python3 /tmp/q1_check.py)
- FINAL.md files in parley-deck/ideas: 78.
- Ideas with FINAL.md + non-terminal status: 3 (NOT zero; fix incomplete for these).
  Names/status: tui-action-execution (accepted), protocol-overlay-local-extension (status=** malformed), meta-protocol-change-rho-retrospective-optimization (final-design-for-review).
- The remaining 75 are terminal (status=final). This is clean evidence that zcode-1/F14 repair did not corrupt terminal records; it just left 3 anomalies untouched (likely outside the 20-idea repair batch, or malformed status lines).

## Finding 2 (Q2 - PRIMARY measurement: find + grep against parley-deck/ideas/*/review/)
- Total review artifacts (files under .../review/): 659 (counted: find | wc -l).
- Artifacts missing `reviewed-commit` in frontmatter/body: 238 (counted: find ... -exec grep -q ... | wc -l).
- Those 238 are historical and NOT revalidated by this audit (rule: historical artifacts stay as-is).
- Evidence of the measurement: `find ... -path "*/review/*" | wc -l` = 659; `find ... -exec grep -L "reviewed-commit" ... | wc -l` = 238. The number 238 tells the deck most of its review record predates the `reviewed-commit` validator rule.

## Finding 3 (Q3 - PRIMARY measurement: python3 scan of all .md files inside any .../round-01/ path under ideas/; checks for non-empty sections Summary / Refutation attempts / Findings / Open questions)
- Round-01 artifacts scanned: 478 (counted by os.walk over paths containing "round-01"; includes both agent-output round-01/ and review/round-01/ artifacts, .md only, excludes hidden/log files).
- Artifacts failing the non-blank required-sections check (at least one of the 4 required sections missing or blank): 434.
- Command evidence: `python3 /tmp/q1_check.py`-style scan (re-run as separate python3 script scanning all 478 files for section headings) produced `missing_or_blank=434`.
- Interpretation: the vast majority of historical round-01 artifacts predate the non-blank-required-section validator. They are NOT being revalidated in this audit (historical exemption); the 434 number is evidence of how much of the deck's round-01 record predates the stricter rule.

## Refutation attempts (self-check: no browser opened; only PRIMARY commands executed)
- No Chrome/browser automation attempted (owner's browser off-limits; not needed for file measurement).
- No external secrets read; only repo-local .md files inspected.
- No files outside the user's designated review artifact were modified. The only write was this file (absolute path verified), not a relative path (earlier session defect: relative paths resolved against $HOME; avoided by using absolute path explicitly).

## Open questions
- Should the 3 remaining non-terminal FINAL.md ideas (tui-action-execution / protocol-overlay-local-extension / meta-protocol-change-rho-retrospective-optimization) be folded into a follow-up repair (e.g., zcode-1/F15) that targets malformed/edge-case status lines specifically?
- The 238 historical review artifacts without `reviewed-commit` and 434 round-01 artifacts with blank required sections are not being retroactively enforced; does the deck want an opt-in migration pass, or is the historical exemption permanent?
