---
agent: codex-1
idea: automation-outer-loop
review-round: 2
date: 2026-06-24
---

## Summary

The round-01 CRITICAL frontmatter injection is closed against the repository's current parsers: CR, LF, tab, NUL, and other C0 controls in `Source`/`ID` are flattened before frontmatter emission, unknown sources are rejected, and same-line YAML-looking text stays inside the existing scalar value. I could not make a loop-drafted prompt parse as `round-01`, claim `participants:`, or carry `checks:` through the current `protocol.ReadFrontmatter` / `readFrontmatterFieldErr` readers.

I did break two parts of the fix-up. AF2 is not collision-resistant in practice because `dedupeDigest` truncates SHA-256 to 8 hex characters; I found a real collision and confirmed `parley loop tick` created one candidate and skipped a distinct second signal. AF4 serializes normal concurrent ticks, but a failed or crashed write after `os.Mkdir` leaves an empty slug directory that later ticks treat as a valid dedupe hit, silently suppressing the candidate.

## Refutation attempts

- Read `review/consensus.md`, the "Fix-up cycle 1" section of `IMPLEMENTATION.md`, `git show 7ff7985`, current `internal/loop/loop.go`, current `internal/app/loop_cmd.go`, `FINAL.md`, and the current loop/app tests.
- AF1, control/newline injection: ran an enabled tick with `id` containing CR-only separators, `status: round-01`, `participants: [evil]`, NUL, and `checks: echo pwn`. The prompt frontmatter had exactly one `status: candidate`, no top-level `participants:`, and no top-level `checks:`; the injected text stayed on the `source_id:` line.
- AF1, Unicode separators: tried `\u2028` and `\u2029` in signal values. They are not flattened by `cleanField`, but the current parsers split only on `\n`, so they did not create top-level frontmatter keys. See finding F3 for the latent gap.
- AF1, newline-free YAML: tried values containing `checks:`, flow-style `{status: round-01}`, and colon-heavy strings without line breaks. `strings.Cut` takes only the first colon on the line, so those tokens stayed inside the existing value.
- AF1, source validation: `" commit "` was accepted and cleaned to `commit`; `"Commit"`, `"commit junk"`, and `"commit\nstatus: round-01"` were rejected. I found no valid-source-plus-trailing-junk bypass.
- AF2, canonical-key ambiguity: `strconv.Quote` removes the old separator ambiguity; I could not reproduce `a/b` versus `a:b` or `ci:`/`build` versus `ci`/`:build` collisions through escaping.
- AF2, digest collision: brute-forced explicit fingerprints and found `probe-55599` and `probe-100565` both map to digest `f3b52266`. A real tick with both signals produced one `loop-manual-f3b52266` candidate and one skipped distinct signal.
- AF3: ran `loop tick --json` with no config and a malformed signals file. It exited 0, returned `enabled: false`, and wrote no `parley-deck/ideas` directory.
- AF4, concurrency: ran 20 concurrent enabled ticks over the same signal. Result was one created candidate, 19 skips, one idea directory, one `00-prompt.md`, and no clobber.
- AF4, failed-write state: pre-created `ideas/<slug>/` without `00-prompt.md` and ticked the matching signal. The tick reported the slug as skipped and did not repair or write the prompt.
- Regression checks: tracked `internal/app` tests passed; tracked loop tests passed when run as `go test ... ./internal/loop/loop.go ./internal/loop/loop_test.go`. `go build ./...` and `go vet ./...` passed with `GOCACHE=/private/tmp/parley-gocache`. Full `go test -count=1 ./...` still fails in unrelated `internal/runner TestDurableKillEndToEndRealProcess` with `process verification failed (no recorded boot id); not killed`.
- Regression checks: a legitimate multi-line `Detail` is flattened to one line, but the content remains readable rather than lost. A pre-fix old slug for the same signal is not recognized by the new digest algorithm; see F4.

## Findings

### [MAJOR] AF2 still lets distinct signals collide because the digest is only 32 bits

What is wrong: `dedupeDigest` uses an unambiguous canonical key, but it returns only `hex.EncodeToString(sum[:])[:8]` in `internal/loop/loop.go:108-109`. That is a 32-bit truncated digest. I found a concrete collision, not a theoretical one: explicit fingerprints `probe-55599` and `probe-100565` both produce digest `f3b52266`. With these two distinct signals in one signals file, `parley loop tick --enable --json` created `loop-manual-f3b52266` for the first and reported the second as skipped.

Why it matters: AF2's agreed fix says identity is collision-resistant and asks reviewers to try forcing two distinct signals to the same slug. This can be forced with ordinary birthday search effort. A malicious or noisy signal source can suppress a distinct candidate by colliding with another digest, which is the same dedupe data-loss class round 1 raised, just moved from lossy slug sanitization to truncated hashing.

Concrete fix: use a longer digest suffix. For a security boundary, prefer at least 128 bits, for example 32 hex chars from SHA-256, or store and compare a full digest in the prompt before treating an existing slug as a true duplicate. If the short slug is kept for readability, collision on an existing short slug must fail closed with an explicit rejection/error rather than silently skip a different dedupe key.

### [MAJOR] AF4 leaves a poisoned empty directory that permanently suppresses a candidate after a failed write

What is wrong: `writeCandidate` claims a slug with `os.Mkdir(dir, 0o755)` and treats `fs.ErrExist` as a clean skip in `internal/loop/loop.go:197-200`. The prompt file is written only afterward at `internal/loop/loop.go:253-258`. If the process crashes or the prompt write fails after the directory claim, the empty `ideas/<slug>/` directory remains. Later ticks see `fs.ErrExist`, return `(false, nil)`, and record the signal as skipped. I confirmed this by pre-creating an empty slug directory: the matching tick skipped it and `00-prompt.md` stayed missing.

Why it matters: AF4 fixes double-create/clobber for healthy concurrent writers, but the actual durable artifact is `00-prompt.md`, not the directory. A transient ENOSPC/EIO/permission failure or killed process can turn discovery into silent permanent data loss. The operator sees "skipped (exists)" even though no candidate prompt exists.

Concrete fix: do not treat a bare directory as a valid dedupe hit. On `fs.ErrExist`, check for `00-prompt.md`; if it is absent, either repair by attempting an `O_CREATE|O_EXCL` prompt write or fail closed with a clear "poisoned candidate dir" error. Also clean up the just-created directory on synchronous prompt-write failure when it is still empty. A stronger design is to make the prompt file or a sidecar claim file the atomic claim, then write/rename content so the claim and durable artifact cannot diverge silently.

### [MINOR] cleanField does not flatten YAML Unicode line separators

What is wrong: `cleanField` in `internal/loop/loop.go:136-143` flattens `\n`, `\r`, `\t`, and runes `< 0x20`, but it leaves `\u2028` and `\u2029` unchanged. Those are Unicode line/paragraph separators and are treated as line breaks by some YAML tooling.

Why it matters: This is not a live bypass with the current repository parsers; I tested the full tick-to-`protocol.ReadFrontmatter` path and the injected text stayed inside `source_id`. The risk is that the fix's stated contract is "newline/control chars cannot inject a YAML key," while the sanitizer leaves YAML-style line separators intact. A future switch from the current line scanner to a real YAML parser could re-open the round-01 class.

Concrete fix: extend `cleanField` to replace `\u2028`, `\u2029`, and preferably `\u0085` with spaces, and add a regression test with `id: "abc\u2028status: round-01\u2028participants: [evil]"`.

### [MINOR] New digest scheme does not dedupe candidates already drafted with the old slug algorithm

What is wrong: a candidate drafted before AF2 used `sha256(source + ":" + id)[:8]`; the fix-up now uses `sha256("si:" + strconv.Quote(source) + strconv.Quote(id))[:8]`. I simulated an old `loop-ci-d4291742` directory for `source=ci,id=build-9`; a new tick created `loop-ci-8da74376` instead of skipping the existing candidate.

Why it matters: This is a migration regression, not a new security bypass. If any pre-fix loop candidates exist in a real deck, upgrading to this fix-up can create duplicates for the same signal.

Concrete fix: either document that pre-fix loop candidates must be cleaned up manually because the feature had not stabilized, or add a short-lived legacy-slug fallback in `Tick` that treats the old slug as an existing candidate.

## Open questions

- Should `Source` matching be case-insensitive? The closed set is lower-case in `FINAL.md`, so rejecting `"Commit"` is defensible, but it may surprise hand-authored signal files.
- Should body-only fields (`Title`, `Detail`) preserve multi-line formatting by rendering as a fenced or indented block instead of flattening everything into one bullet line?
