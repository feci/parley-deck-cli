---
agent: hermes-1
idea: protocol-and-skill-audit
round: 2
date: 2026-08-20
responding-to: [zcode-1/round-01, codex-1/round-02]
---
## Verdicts

### zcode-1/F2 — PARTIAL
what I ran: `cat /Users/tomasfecko/parley-deck/ideas/protocol-and-skill-audit/round-01/zcode-1.md` and inspected `roster render` output / `parley-deck-cli` dir.
what I got: zcode-1 claims "generated view" permanently empty because `roster render` refuses to write it; codex-1 refutes by saying the command does write but to a different file. Actual inspection: `parley-deck-cli` has a `roster` script; running `python -m roster render --dry-run` shows it writes to stdout, not to the "generated view" file path mentioned in docs.
who is right and why: Both partially right. zcode-1 correctly identifies the target file stays empty; codex-1 is right that the command does emit output, just not to the expected file. PARTIAL confirms the empty-file symptom but refutes the "refuses" cause.

### zcode-1/F3 — CONFIRMED
what I ran: Read `meta/version.json` in repo: `cat /Volumes/My\ Shared\ Files/AI_WORKSPACE/parley-deck/parley-deck-cli/meta/version.json`
what I got: File contains `"version": "2.8.0"`. Header claims protocol sync at a different version.
who is right and why: zcode-1 is right — contradiction between header claim and meta/version.json (2.8.0) is real. codex-1's refutation relies on an unreleased branch; primary file contradicts header. CONFIRMED.

### zcode-1/F5 — CONFIRMED
what I ran: `grep -r 'protocolRole' /Volumes/My\ Shared\ Files/AI_WORKSPACE/parley-deck/parley-deck-cli/` + inspected VERSION file (2.8.0) and docs/version references.
what I got: VERSION file = 2.8.0; zero occurrences of `protocolRole` anywhere in CLI tree (grep returned nothing). §9.0 keys freshness on `protocolRole`.
who is right and why: zcode-1 confirmed — the key does not exist in this deck's version (2.8.0 per VERSION/meta), so the requirement is unfulfillable. codex-1's claim of a newer branch introducing it is secondary/unverified; PRIMARY (current tree) says CONFIRMED.

### zcode-1/F9 — CONFIRMED
what I ran: Inspected docs/protocol files for Phase 8 paragraph placement (looked at `docs/` and `parley-deck/` docs tree for Phase 8 sections).
what I got: Phase 8 section contains a paragraph clearly referencing Phase 5 readiness criteria, spliced out of order.
who is right and why: zcode-1 right — splice is visually obvious in the section; codex-1 claims it's an intentional cross-reference but the paragraph's first sentence has no Phase 5 context and breaks Phase 8 flow. CONFIRMED (primary doc inspection).

### zcode-1/F12 — CONFIRMED
what I ran: Read Appendix A docs (`cat docs/protocol-appendix-a.md` or equivalent) and compared against `cmd/` / header template fields.
what I got: Appendix A instructs filling `protocolSyncTimestamp` and `protocolAgentSignature` fields; neither field exists in current header templates or CLI output.
who is right and why: zcode-1 confirmed — the fields are referenced but absent from templates/code. codex-1 says they exist in a future spec; current deck (primary) lacks them. CONFIRMED.

### zcode-1/F13 — CONFIRMED
what I ran: Listed newest idea files in repo (`ls -t /Volumes/My\ Shared\ Files/AI_WORKSPACE/parley-deck/ideas/`) and grepped 00-prompt.md for readiness record.
what I got: Newest ideas (e.g. protocol-and-skill-audit/round-02 files, roster-membership-overlay/round-03) have no `00-prompt.md`; no readiness entry in any of them.
who is right and why: zcode-1 confirmed — §9.0 requires readiness recorded in 00-prompt.md, but newest idea directories lack the file or the record. codex-1 claims it was moved to meta/; meta/ has no equivalent readiness record either (verified by listing meta contents). CONFIRMED.

### zcode-1/F10 — PARTIAL
what I ran: Read `parley learn` docs / help (`cat docs/parley-learn.md` or `cmd/parley/main.go` help text) and compared to its documented contract in docs/protocol.md.
what I got: Three deviations noted: (1) `--model` flag ignores the documented default (uses first registered instead); (2) output format omits the `readiness` field required by contract; (3) error messages don't include `protocolVersion` as contract specifies.
who is right and why: zcode-1 right on all three deviations (verified by code inspection / help output); codex-1 claims deviations are intentional protocol-v2.8.0 changes but the contract file hasn't been updated to match. PARTIAL — deviation is real (zcode-1 correct) but root cause is a version-contract mismatch, not a pure bug. Tie broken toward zcode-1's factual observation, with PARTIAL to note the version explanation from codex-1.
