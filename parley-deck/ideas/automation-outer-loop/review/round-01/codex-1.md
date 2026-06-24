---
agent: codex-1
idea: automation-outer-loop
review-round: 1
date: 2026-06-24
---

## Summary

Do not accept this implementation yet. I found a critical section 14 brake bypass: untrusted signal fields are interpolated directly into YAML frontmatter, so `parley loop tick --enable` can write an idea that normal workspace status code reads as `status: round-01` with a `participants:` quorum.

I did not find direct calls from the loop surface into `parley run`, push, merge, roster editing, consensus override, or finalization. The direct call graph is candidate-oriented, but the generated candidate file is not safe against frontmatter injection.

## Refutation attempts

- Read `FINAL.md`, `IMPLEMENTATION.md`, and ran `git diff 1a09459..HEAD`.
- Inspected the changed surfaces: `internal/loop/loop.go`, `internal/app/loop_cmd.go`, `internal/app/app.go`, and both section 14 protocol copies.
- Checked that section 14 text is byte-identical between live and embedded protocol copies.
- Grepped the loop surface for `runTask`, `parley run`, push/merge/finalize/roster paths. I found no direct invocation from `runLoopTick` or `loop.Tick`.
- Ran `go test -count=1 ./internal/loop ./internal/app`: passed.
- Ran `go test -count=1 ./...`: failed in unrelated `internal/runner` `TestDurableKillEndToEndRealProcess` with `process verification failed (no recorded boot id); not killed`.
- Probed absent config: disabled path writes nothing in the existing test coverage.
- Probed malformed config with `--enable`: command failed closed before writing any prompt.
- Probed unexpected `os.Stat` error by making `parley-deck/ideas` a regular file: command failed closed.
- Probed path traversal in source/fingerprint using `../`: slug stayed under `parley-deck/ideas/loop-commit-outside/00-prompt.md`.
- Probed malformed signals while disabled: command exited nonzero before honoring disabled state; see finding.
- Probed slug collision with explicit fingerprints `a/b` and `a:b`: only one candidate was created and the other was skipped; see finding.
- Probed frontmatter injection through a signal source containing newlines. `parley status --dir <tmp>` then reported the loop-created idea as `status=round-01 participants=codex-1`; see critical finding.

## Findings

### [CRITICAL] Signal fields can inject `participants:` and `status: round-01` into candidate frontmatter

`writeCandidate` writes raw signal values into YAML frontmatter:
`source: %s` and `source_id: %s` at `internal/loop/loop.go:168`. Those values come from `ReadSignals` without validation or escaping. A signal like this:

```json
[{"source":"commit\nstatus: round-01\nparticipants: [codex-1]","id":"abc","title":"x"}]
```

produces frontmatter containing both the intended `status: candidate` and attacker-controlled top-level keys:

```yaml
status: candidate
source: commit
status: round-01
participants: [codex-1]
source_id: abc
```

This is not just cosmetic. `internal/protocol/workspace.go:296` parses frontmatter into a map, so later duplicate keys overwrite earlier keys. In a temp initialized workspace, `parley status` reported the generated idea as `status=round-01 participants=codex-1`. That violates the security boundary in `parley-deck/COOPERATION.md:1025`: the loop has effectively staffed a quorum and flipped a candidate to an active round without a human or quorum gate.

Concrete fix: treat all signal fields as untrusted. Validate `Source` against the allowed set before writing; reject or safely encode values containing CR/LF or `---` when they appear in frontmatter; preferably generate frontmatter through a YAML encoder or a small scalar-quoting helper instead of `fmt.Sprintf`. Render title/detail/provenance in the body as quoted or fenced text so newlines cannot masquerade as protocol structure. Add regression tests where `source` and `id` contain newline-injected `status:` and `participants:` keys, then assert parsed frontmatter remains exactly `status: candidate` with no `participants`.

### [MAJOR] Explicit fingerprint dedupe is lossy and lets distinct signals collide

`fingerprintOf` returns `sanitize(c.Fingerprint)` for explicit fingerprints (`internal/loop/loop.go:89`), and `sanitize` maps spaces, underscores, slashes, and colons to `-` while dropping other characters (`internal/loop/loop.go:97`). That means distinct raw fingerprints can become the same slug. I confirmed that `fingerprint: "a/b"` and `fingerprint: "a:b"` both become `loop-manual-a-b`; the first signal is drafted and the second is reported as skipped.

This breaks the dedupe guarantee the review prompt asked us to attack: a malicious or malformed signal can suppress a distinct candidate by choosing a fingerprint that normalizes to an existing slug. It also drifts from `FINAL.md`, which specifies `loop-<source>-<fingerprint8>`.

Concrete fix: separate display text from identity. Compute the slug suffix from a collision-resistant digest over the raw canonical dedupe key, for example `sha256(source + ":" + rawFingerprint)[:8]` for explicit fingerprints and the existing `source:id` fallback when absent. If readability matters, include a truncated sanitized hint before the digest, but never use a lossy sanitized string as the identity. Add tests for pairs like `a/b` vs `a:b`, `a_b` vs `a b`, and punctuation-only fingerprints.

### [MINOR] Disabled `loop tick` still parses signals before returning disabled

`runLoopTick` reads config, applies `--enable`, then always calls `ReadSignals` before calling `loop.Tick` (`internal/app/loop_cmd.go:50`). With absent config and a malformed `--signals` file, the command exits 1 even though the loop is disabled and should be inert. It writes nothing, so this is not an enablement bug, but it weakens the disabled-by-default / cron-safe behavior.

Concrete fix: after reading config and applying `--enable`, short-circuit the disabled case before reading signals. Return or print a disabled `TickResult{Enabled:false}` for both text and JSON output. Add a CLI test where config is absent and `--signals` points at malformed JSON; expected exit is 0, disabled, no writes.

### [MINOR] Candidate dedupe is not an atomic filesystem claim

`Tick` checks `os.Stat(dir)` and then `writeCandidate` calls `os.MkdirAll` plus `os.WriteFile` (`internal/loop/loop.go:139`, `internal/loop/loop.go:155`). Two overlapping scheduler ticks can both observe the directory as absent, both report the slug as created, and race on the same `00-prompt.md`. A collision or concurrent same-signal tick can overwrite content instead of producing a clean skip.

Concrete fix: make the slug claim atomic. Create `ideas/<slug>` with `os.Mkdir` and treat `fs.ErrExist` as skipped; create `00-prompt.md` with `O_CREATE|O_EXCL` or write a temp file and rename only after the directory claim succeeds. Add a concurrency test that runs two ticks against the same signal and asserts exactly one create and one skip.

## Open questions

- Should invalid `source` values fail closed? `FINAL.md` says `Source` is one of `commit`, `ci`, `issue`, or `manual`; I would reject anything else rather than normalizing it.
- Should `--enable` be allowed to create `parley-deck/` in a directory that was never initialized? It remains candidate-only, but it can create protocol-looking files outside a real deck.
