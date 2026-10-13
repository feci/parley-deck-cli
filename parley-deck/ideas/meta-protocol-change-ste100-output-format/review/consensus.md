---
idea: meta-protocol-change-ste100-output-format
review-cycle: 1
drafted-by: codex-1
date: 2026-10-13
reviewed-commit: 4aab9f339e81646b5f2bc2325562cfc1f89e3fcf
reviewed-skill-commit: 92e3dc78da2f1c55d35c905cabe57fc69ede0670
outstanding_agreed_fixes: 3
blocked: false
---

## Agreed fixes

1. From zcode-1/review/round-01 [CRITICAL], repair oversized headless argument delivery.
   Preserve the exact prepared prompt, including the full protocol, attestation, and task, in the private invocation directory.
   Use a unique file created with mode 0600. Write and sync its complete bytes before child launch.
   For a prompt argument at or above 120 KiB, substitute a short file instruction with its absolute path and SHA256.
   The instruction requires the agent to verify the hash and read the complete file before the task.
   A missing, unreadable, or changed file requires a blocker. File delivery does not prove that the model read it.
   Keep the original stdin bytes. Do not switch provider modes, add native flags, trim the protocol, or enable optimized packets.
   Apply the same helper at both measured and supervised headless launch paths.
   Preserve embedded and standalone placeholders. Refuse if an expanded argument remains over the bound after substitution.
   Retain the private file for audit. Do not expose its raw text in public telemetry.
   Test real child execution above Linux's limit, exact bytes and hashes, file permissions, short prompts, stdin, and failed publication.
   Use the candidate CLI for the next real Zcode review, so that review exercises the file delivery path.
2. From zcode-1/review/round-01 [MAJOR], update the live map test's expected count from 69 to 70.
   Keep the structural map and never-cut checks. Add a focused assertion for Appendix B's on-demand status and the retained section 6 pointer.
   Preserve the original 70,000-byte guard and the 20,000-byte skill guard.
3. From zcode-1/review/round-01 [MINOR], record the failed local suite and both failed Linux jobs in IMPLEMENTATION.md.
   Record command exit states explicitly. Empty successful build output is not the evidence by itself.
   Add the complete packet request parameters to source-sizes.json.
   Re-run the required CLI suite and relevant checks after the fixes.

The transport repair is an unavoidable implementation deviation from FINAL's initial file list.
It addresses a regression caused by the approved protocol growth. It does not change the signed H1/H2 bytes or writing policy.
New or changed runner files and their focused tests must appear in the next IMPLEMENTATION record.
The existing stdin fake fixture may receive a short file argument, but stdin remains byte-exact.
Do not make a fixture-only fix that hides the production argument limit.

## Deferred follow-ups

Broad native Windows baseline repair remains deferred, with follow-up slug TBD. CLI WinGet remains held.
Full ASD vocabulary verification, a prose linter, and historical re-authoring remain outside this idea.

## Dismissed findings

None. The implementer accepts all three findings.

## Coverage & blind spots

Zcode independently checked the full CLI and skill diffs and AC1–AC8.
It confirmed AC1–AC6 and AC8, and broke AC7 through local tests and Linux CI.
Codex's own checks found the same failures. Codex is not a code reviewer.

Concur with the exclusion of probePrompt: it asks for a liveness sentinel, not new protocol prose.
The thirteen-path writing-policy boundary stays as FINAL states.
The next full source review must inspect the added delivery helper as well as the complete original diff.
No model read claim follows merely from a successful file write or child start.

## Signoffs

<!-- codex-1 appends its signoff below -->
<!-- zcode-1 appends its signoff below -->

### Signoff: codex-1 — 2026-10-13
Status: ✅ ACCEPT
Notes: I accept all three fixes and the bounded private-file delivery mechanism. H1 and H2 stay exact. No source close is authorized by this fix consensus.

### Signoff: zcode-1 — 2026-10-13
Status: ✅ ACCEPT
Notes: I compared this consensus with my round-01 review. The three fixes match my CRITICAL, MAJOR, and MINOR findings. The private-file delivery stays below the kernel argument limit and keeps stdin bytes exact. My next review must verify the delivery helper in the real tree. This signoff authorizes fix-up cycle 1 only. It does not close the idea.
