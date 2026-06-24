---
agent: antigravity-1
idea: automation-outer-loop
review-round: 1
date: 2026-06-24
---

## Summary

The Tier 4 outer loop implementation (LE-8 + LE-9) introduces a human-braked automated discovery layer that limits scheduler-driven operations (cron, CI, MCP) to drafting non-active `status: candidate` idea prompts. 

A thorough refutation analysis shows that the security boundary (the §14 human brake) is mostly well-isolated in the Go codebase. There are no code paths in `loop.Tick`, `runLoop`, or CLI wiring that allow automated staff quorum promotion, roster edits, consensus overrides, merges, pushes, or runs. However, a **critical YAML frontmatter injection vulnerability** exists because external signal data is written directly to the draft prompt's YAML frontmatter without newline sanitization. This allows malicious signals to bypass the human brake and promote themselves or staff quorums.

Additionally, some minor issues regarding signal collision/deduplication under boundary-shifting inputs were identified.

## Refutation attempts

I systematically attempted to break the implementation through the following verification strategies:

1. **Security Boundary Bypass Checks**: Traced [loop.go](file:///Volumes/My%20Shared%20Files/AI_WORKSPACE/parley-deck/parley-deck-cli/internal/loop/loop.go) and [loop_cmd.go](file:///Volumes/My%20Shared%20Files/AI_WORKSPACE/parley-deck/parley-deck-cli/internal/app/loop_cmd.go) to confirm no execution of deliberations, git pushes/merges, or writing of `participants:`/`status: round-01` occurs. Tested whether a candidate prompt could be loaded or run automatically.
2. **YAML/Frontmatter Injection Attack**: Attempted to inject carriage returns/newlines and YAML key-value pairs via `c.Source` and `c.ID` to see if they could write arbitrary keys like `status: round-01` into `00-prompt.md`.
3. **Fail-Closed Verification**: Tested if the loop could be enabled by malformed config. Verified behavior under:
   - Missing `loop/config.json` (returns `Enabled=false`, exits 0, writes nothing).
   - Malformed `loop/config.json` (unmarshal error, CLI exits 1).
   - Empty/0-byte `config.json` (unmarshal error, CLI exits 1).
4. **Path Traversal via Slug**: Evaluated whether directory traversal sequences (e.g. `../../etc/passwd` or `..`) inside `Source`, `ID`, or `Fingerprint` could escape the `ideas/` directory. Verified that the `sanitize` whitelist strictly allows `[a-z0-9-]` only, turning directory traversals into safe alphanumeric strings like `etc-passwd` or `x`.
5. **Deduplication & Error Handling**: Analyzed `os.Stat` error handling. Verified that if `os.Stat(dir)` encounters an unexpected error (like permission denied), it halts and returns the error (fail-closed), rather than skipping or overwriting. Checked if two distinct signals could collision-merge under default hash-generation inputs.
6. **Command Flag Isolation**: Assessed if the `--enable` flag on `parley loop tick` has side-effects beyond setting `cfg.Enabled = true`. Verified it only permits a candidate-only discovery run.
7. **Test Executions**: Executed the test suite using `go test -count=1 ./...` to verify that all existing tests and the embedded default drift guards are functional and green.

## Findings

### CRITICAL: YAML Frontmatter Injection in Candidate Prompts
- **What is wrong**: In [loop.go](file:///Volumes/My%20Shared%20Files/AI_WORKSPACE/parley-deck/parley-deck-cli/internal/loop/loop.go#L153-L198), `writeCandidate` formats raw `c.Source` and `c.ID` values directly into the YAML frontmatter of `00-prompt.md` without any sanitization or newline stripping:
  ```go
  prompt := fmt.Sprintf(`---
idea: %s
author: loop
created: %s
status: candidate
source: %s
source_id: %s
fingerprint: %s
---`
  ```
- **Why it matters**: If a signals source file is malicious or contains untrusted input (e.g., a git commit message or issue description mapped to `Source` or `ID`), it can inject newlines and YAML keys:
  ```json
  {
    "source": "commit\nstatus: round-01\nparticipants: [claude-1]",
    "id": "123"
  }
  ```
  This generates:
  ```yaml
  ---
  idea: loop-commit-xxx
  author: loop
  created: 2026-06-24
  status: candidate
  source: commit
  status: round-01
  participants: [claude-1]
  source_id: 123
  fingerprint: xxx
  ---
  ```
  YAML parsers reading this file may parse the duplicate/override `status: round-01` and the newly-injected `participants:` key. This bypasses the §14 human brake by allowing a loop-drafted candidate to immediately self-promote to `round-01` and pretend it has a staffed quorum.
- **Concrete fix**: Sanitize or escape `c.Source` and `c.ID` in `writeCandidate` to replace or remove carriage returns and newlines, or serialize the frontmatter using a safe YAML parser/encoder.
  ```go
  cleanSource := strings.ReplaceAll(strings.ReplaceAll(c.Source, "\n", " "), "\r", " ")
  cleanID := strings.ReplaceAll(strings.ReplaceAll(c.ID, "\n", " "), "\r", " ")
  ```

### MINOR: Boundary-Shift Collision in Default Fingerprints
- **What is wrong**: In [loop.go](file:///Volumes/My%20Shared%20Files/AI_WORKSPACE/parley-deck/parley-deck-cli/internal/loop/loop.go#L87-L94), `fingerprintOf` derives stable fingerprints via `sha256.Sum256([]byte(c.Source + ":" + c.ID))`.
- **Why it matters**: If a signal's source or ID contains a colon, the boundaries can shift and create identical hash inputs. For example:
  - Signal A: `Source = "ci:"`, `ID = "build"` -> `"ci::build"`
  - Signal B: `Source = "ci"`, `ID = ":build"` -> `"ci::build"`
  These distinct signals will hash to the exact same fingerprint and resolve to the same slug. The second signal will be incorrectly deduped and skipped.
- **Concrete fix**: Escape colons inside `c.Source` and `c.ID` before joining them with a colon, or use a structured encoding (like JSON serialization) for the hash input.
  ```go
  sum := sha256.Sum256([]byte(strings.ReplaceAll(c.Source, ":", "\\:") + ":" + strings.ReplaceAll(c.ID, ":", "\\:")))
  ```

### NIT: Unicode/Emoji Fingerprints Collapse to Generic "x" Slug
- **What is wrong**: The `sanitize` function only preserves `[a-z0-9-]`. If an explicit fingerprint contains only non-ASCII characters or emojis, `sanitize` returns an empty string which is replaced by `"x"`.
- **Why it matters**: If multiple signals specify different emoji fingerprints (e.g. `🚀` and `🔥`) under a shared source, they will all map to the exact same slug (e.g., `loop-signal-x`), causing them to collision-dedupe.
- **Concrete fix**: If `sanitize(fp)` returns `"x"` but the original fingerprint string was not empty, append a short hash of the original unsanitized fingerprint (e.g. `x-<hash>`) to guarantee slug uniqueness.

## Open questions

1. **Concurrency and File Locking**: Should `parley loop tick` acquire a file lock (e.g. `loop.lock` or similar) to prevent race conditions if multiple schedulers/cron-jobs invoke the command concurrently on the same workspace?
2. **0-Byte signals/config files**: Currently, an empty (0-byte) `signals.json` or `config.json` file returns a JSON unmarshalling error and exits 1 (fail-closed). Should empty files be treated as missing (safe/no-op) or is strict validation preferred?
3. **Draft Mutation Checks**: If a candidate directory exists but the `00-prompt.md` was manually deleted, `os.Stat` on the directory still returns `nil` and the signal is skipped. Should the deduplication check verify the existence of the `00-prompt.md` file rather than just the directory?
