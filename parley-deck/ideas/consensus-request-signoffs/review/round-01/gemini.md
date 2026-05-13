---
agent: gemini
idea: consensus-request-signoffs
review-round: 1
date: 2026-05-13
reviewed-commit: d9337399c238a1e7b72b2d75cf4bbc24590c22ea
---

## Findings

- `MAJOR` internal/app/consensus_request_signoffs.go:196 - `validateRequestedSignoff` only verifies that existing signoff blocks are unchanged, but it does not check if other parts of the consensus file (frontmatter, headings, or the "Agreed decisions" section) were modified. This fails to fully enforce the "edit no other line" protocol rule. An agent could technically modify the agreed plan while signing off without being detected by the current validation logic. Suggested fix: Read and compare the raw file content before and after invocation, ensuring that the only difference is the appended signoff block.

- `MINOR` internal/app/consensus_request_signoffs.go:214 - `buildConsensusSignoffPrompt` omits the `Counter-proposal` field from the canonical block shape example. While the text instructions mention it, headless agents frequently prioritize the visual example, which may lead to malformed `BLOCK` signoffs that lack the required counter-proposal field. Suggested fix: Update the prompt example to include an optional `Counter-proposal` line or a note on where to place it for `BLOCK` statuses.

- `NIT` internal/app/consensus_request_signoffs.go:88 - When a signoff sequence is interrupted by a `BLOCK` or a child process error, the command does not print a final recap of which signoffs were successfully appended before the failure. While individual successes are streamed to stdout, a final "Partial progress" summary would improve clarity for the user. Suggested fix: Catch the error in `requestConsensusSignoffs` and print the `successes` slice before returning the error.

## Open questions

- Should the command automatically commit the consensus file after each successful signoff? Currently, it preserves partial progress on disk but doesn't interact with Git, which is consistent with other `parley` commands but might lead to messy unstaged changes in long signoff rounds.

## Summary

The implementation is solid and handles the core requirements well, including dry-runs, hosted backend gating, and sequential invocation with strict per-agent validation. The identified `MAJOR` finding regarding preamble/section integrity is the most significant gap in protocol enforcement. The tests are comprehensive and correctly cover failure modes like `BLOCK` stops and non-zero exit codes after valid appends.
