Parley Deck now uses ASD-STE100 Issue 9 as the default for new English protocol prose. Appendix B supplies the binding working summary.

- Preserve exact quotations, code, commands, identifiers, frontmatter, required syntax, and evidence.
- Use `ste_output: false` with an English reason for a per-idea exception at kickoff. Later changes need owner direction or participant consensus.
- Existing decks adopt the rule through their normal authority update. A newer CLI does not override an older live authority.
- Thirteen model-facing prompt paths now carry the same writing-policy instruction.
- Headless prompt arguments at or above 120 KiB use a private file with exact text and a SHA256 read instruction. This fixes the Linux argument-limit regression caused by the larger full protocol.

Existing artifacts remain unchanged. Style feedback remains advisory. The summary does not certify full ASD conformance.

The companion skill and staged core are 2.19.0. npm and global-core publication remain separate attended owner actions.
CLI Windows assets remain experimental. CLI WinGet remains held.
