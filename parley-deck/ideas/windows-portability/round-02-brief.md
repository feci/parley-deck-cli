# Round 2 organizer routing

All three round-01 artifacts pass the shared validator. This is a procedural opening of cross-review, not acceptance of any proposal or code verdict. Each participant must address both peers and resolve disagreements with located evidence in its own round-02 artifact.

Specific disagreements to engage:
- ACL acceptance and creation: owner only versus owner+SYSTEM versus owner+SYSTEM+Administrators, inheritance, readback/adversarial tests, and the conflicting claim whether x/sys v0.36.0 already wraps ACL construction (ACLFromEntries et al.).
- File-sharing failures: existing LockFileEx versus actual open/share/delete/rename mechanisms. Distinguish demonstrated causes from hypotheses and bounded-retry proposals.
- Directory durability: Windows no-op versus a real flush mechanism versus an explicit user-visible refusal. Apply the owner's product-behavior constraint; documentation alone must not silently replace the required real implementation or refusal.
- Process scope: truthful liveness with reviewed durable-kill refusal versus full Job Objects/attribution. Preserve the reviewed Linux repairs; resolve unsupported claims and scope tradeoffs through peer evidence.
- Product shell execution and CRLF: the new product sh dependency correction, honest missing-shell contract, native fixtures, preservation of content integrity, and actual CI image assumptions.
- Gate filenames: universal encoding plus compatibility fallback versus per-OS names, invalid IDs/path safety, legacy read compatibility and collision handling.
- Test mappings: FIFO/Windows socket/reparse alternatives, pipe-capacity-independent ACP guard, unreadable ACL fixtures; every assertion must continue to test its stated product guarantee. Reconcile unexplained failures, rather than equate a green fixture with correct behavior.
- Ownership: Claude's historical claims and newly raised claims need non-owner verdicts. Claims in conflict remain disputed until the evidence is engaged. Do not self-verdict.

Owner-scope reminders (not new owner decisions): Windows assets are requested. The specified execution acceptance is hosted windows-latest, macOS and Ubuntu; no Windows ARM64 or local-native validation was requested or may be claimed. Honest coverage/prerequisite/refusal documentation is required. Do not silently add an ARM64-run release gate, or a requirement that external winget maintainers merge a submitted PR. Raise any substantive safety disagreement through the existing process. Version selection happens after the prerequisite releases, above the actual then-released version. Both predecessor handoffs remain the release gate.

Nominate a single participant to draft and implement, with the other two independently reviewing. codex-1 never drafts participant decisions, implements, verifies code or signs off. No product code edits in this round. No protocol change is proposed; if that changes, flag explicitly.
