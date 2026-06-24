---
agent: hermes-1
idea: automation-outer-loop
review-round: 1
date: 2026-06-24
---

## Summary

The implementation delivers LE-8 (§14 human brake, both COOPERATION.md copies,
drift guard green) and LE-9 (`internal/loop` + `parley loop tick`). I attempted
to break the §14 security boundary and the fail-closed / dedupe / sanitization
guarantees across every named vector. The loop's *code path* is genuinely
draft-only: there is no `run`/`exec`/`push`/`merge`/`finalize` call anywhere in
`internal/loop`, `--enable` does nothing beyond `cfg.Enabled = true`, disabled /
absent / malformed configs all fail closed, `os.Stat` unexpected errors fail
closed, and path traversal via the slug is neutralized by `sanitize`.

However I found a CRITICAL defect: `writeCandidate` interpolates the raw,
unsanitized `c.Source` and `c.ID` directly into the YAML frontmatter of the
drafted `00-prompt.md`. A malicious signal carrying newlines in its `source` or
`id` field injects arbitrary frontmatter keys — including `participants:` (the
quorum claim §14 forbids) and a second `status: round-01` (the promotion the
brake must prevent). The slug is sanitized; the frontmatter values are not. I
confirmed this against the real `loop.Tick` function, not just a format-string
PoC. Because downstream consumers (`protocol.ReadFrontmatter`,
`driver_impl.go`) parse these keys and one of them executes `sh -c` on a
`checks:` field, the poisoned candidate is also a latent RCE vector. The loop
itself does not auto-run, so this is not an *immediate* escalation, but the
drafted artifact violates the §14 output invariant — the brake produces exactly
the quorum-staffed, round-01-flipped candidate it exists to prevent.

## Refutation attempts

1. §14 boundary — can tick do more than discover-and-draft?
   Searched `internal/loop` for `run`/`exec`/`Run`/`push`/`merge`/`finalize`/
   `FINAL`/`quorum`/`participants`/`round-01`. The only hits are doc comments
   saying it does NOT do these things. `Tick` (loop.go:131-151) only calls
   `SlugFor`, `os.Stat`, `writeCandidate`, and appends to result slices.
   `writeCandidate` only calls `os.MkdirAll` + `os.WriteFile`. No subprocess,
   no git, no roster edit. The code path is clean. Could NOT break the
   *action* boundary — but broke the *output* boundary (Finding 1).

2. Disabled-by-default / fail-closed config.
   `ReadConfig` (loop.go:55-68): absent → `{Enabled:false}, nil`; read error →
   `{Enabled:false}, err`; malformed JSON → `{Enabled:false}, error`. The
   caller (`runLoopTick`, loop_cmd.go:50-54) returns exit 1 on any error, and
   `Tick` returns immediately when `!cfg.Enabled` (loop.go:133-134). An absent
   or malformed config can never silently enable the loop. Could not break.

3. Dedupe — slug stability, double-draft, os.Stat fail-open.
   `SlugFor` is deterministic: same `(source, id, fingerprint)` → same slug.
   `Tick` stats the idea dir; if it exists → `Skipped` + continue; if
   `fs.ErrNotExist` → write; any other stat error → return err (fail closed,
   loop.go:142-143). The `os.Stat` error does NOT fail open. Could not break
   the fail-open vector. Two distinct signals sharing the same *explicit*
   fingerprint collide to one slug and the second is silently skipped — but
   this is by-design (fingerprint is the dedupe key); noted as MINOR (Finding 4).

4. `--enable` scope.
   loop_cmd.go:55-59: `if *enable { cfg.Enabled = true }`. That is the entire
   effect. It does not set participants, status, call run, or alter the
   candidate shape. `Tick` with `Enabled=true` still only drafts. Could not
   break.

5. Path traversal via slug (fingerprint/source/id).
   Built a Go harness replicating `sanitize` + `SlugFor` + `filepath.Join` and
   tested `../../../../etc/passwd`, `../../../tmp/evil`, `..`, `...`,
   `../../etc` as source. `sanitize` maps `/` → `-` and drops `.`, so `..`
   becomes `x` and `../` becomes `-`. Every resulting `dir` stays under
   `ideas/` — no escape. `filepath.Rel` confirms no `..` prefix. Could not
   break.

6. Frontmatter injection via signal fields (the break).
   `writeCandidate` (loop.go:168-199) builds the prompt with `fmt.Sprintf` and
   interpolates `c.Source` (line 173, 181, 184), `c.ID` (line 174, 185), and
   `fingerprintOf(c)` (line 175) into the template. The slug and fingerprint
   are sanitized; `c.Source` and `c.ID` are NOT. I wrote a test in
   `internal/loop` that calls the real `Tick` with
   `Source: "commit\nparticipants: [evil-agent]\nstatus: round-01"`. The
   resulting `00-prompt.md` frontmatter contains `participants: [evil-agent]`
   and a second `status: round-01`. Test FAILED (injection confirmed), then I
   removed it. This is Finding 1.

7. Downstream impact of the injected frontmatter.
   Traced two frontmatter parsers. `protocol.ReadFrontmatter`
   (workspace.go:296-325) builds a map; on duplicate keys the LAST wins, so
   injected `status: round-01` / `participants:` override the real values.
   `readFrontmatterField` (cursor.go:283-313) returns the FIRST match, so
   `readIdeaStatus` would still read `candidate`. `driver_impl.go:205-209`
   reads `checks:` via `ReadFrontmatter` and runs `sh -c $checks` — so an
   injected `checks:` key is a latent RCE when a human later promotes+runs.
   `readIdeas` (workspace.go:260-294) would list the poisoned idea as
   `status: round-01` with `participants: [evil-agent]`. The loop does not
   auto-trigger any of this, but the artifact is poisoned. Finding 1.

8. ReadSignals / ReadConfig missing vs malformed.
   `ReadSignals` (loop.go:72-85): missing → `nil, nil`; read error → `nil,
   err`; malformed JSON → `nil, error`. `ReadConfig` as above. Both handle the
   missing-vs-malformed distinction correctly. Could not break.

9. Drift guard (both COOPERATION.md copies).
   The diff adds identical §14 text to both `parley-deck/COOPERATION.md` and
   `internal/protocol/defaults/COOPERATION.md`. `go test ./internal/protocol/`
   (the drift guard, `TestEmbeddedDefaultMatchesLiveDeck`) passes. The §14
   text is outside the five allowlisted normalization zones, so it is compared
   byte-for-byte. Could not break.

## Findings

### CRITICAL — F1: Frontmatter injection via unsanitized Source/ID breaks the §14 output invariant

What: `writeCandidate` (loop.go:173-174) interpolates `c.Source` and `c.ID`
raw into the YAML frontmatter of the drafted candidate. Neither is passed
through `sanitize` or stripped of newlines. A signal with
`source: "commit\nparticipants: [evil-agent]\nstatus: round-01"` produces a
`00-prompt.md` whose frontmatter contains a `participants:` quorum claim and a
second `status: round-01` — exactly the two things §14.2 says an automated
loop must NEVER write. I confirmed this by calling the real `loop.Tick`.

Why it matters: The §14 brake is the security boundary. Its invariant is not
just "the loop doesn't call parley run" — it is "a loop-drafted idea is always
a non-active `status: candidate` with no `participants:` claim" (§14.1,
FINAL.md LE-8). This implementation can produce a candidate that *claims* a
quorum and *claims* `round-01` status. Downstream, `protocol.ReadFrontmatter`
is last-wins on duplicate keys, so `readIdeas` / `parley status` / the TUI
would surface the poisoned idea as an active `round-01` idea with
participants. Worse, `driver_impl.go:205-209` reads a `checks:` frontmatter
key and passes it to `sh -c`, so an attacker who controls the signals file can
plant a `checks: <arbitrary command>` key that executes when a human later
promotes and runs the idea. The loop stays draft-only, but the draft is a
trojan horse.

Concrete fix: Sanitize every field interpolated into the frontmatter (and
reject/escape newlines). At minimum, before formatting, strip control chars
and newlines from `c.Source`, `c.ID`, `c.Title`, `c.Detail`:
```go
func cleanField(s string) string {
    s = strings.TrimSpace(s)
    s = strings.Map(func(r rune) rune {
        if r == '\n' || r == '\r' || r == '\t' || r == 0 {
            return ' '
        }
        return r
    }, s)
    return s
}
```
Apply to `c.Source`, `c.ID`, `c.Title`, `c.Detail` before `fmt.Sprintf`. This
preserves readability (spaces, not dropped) while making frontmatter injection
impossible. Add a test: signal with newlines in source/id → frontmatter has
exactly one `status:` line (`candidate`) and zero `participants:` lines.

### MAJOR — F2: No test exercises the frontmatter-injection vector

What: The existing tests (`TestTickEnabledDraftsCandidate`,
`TestLoopTickEnableDraftsCandidateOnly`) assert no `participants:` key, but
only feed clean signals (no newlines). The `hasFrontmatterKey` helper would
catch an injection if one were tested, but no test provides a hostile signal.

Why it matters: F1 passed all tests because the tests never probe the
boundary with adversarial input. A security-boundary test suite that only
tests benign inputs is theater.

Concrete fix: Add a test that ticks with a signal whose `source`/`id`/`title`/
`detail` contain `\n`, `participants:`, `status: round-01`, `checks:`, and
assert the output frontmatter contains exactly the expected keys and no
injected ones.

### MINOR — F3: Empty `## Constraints` / `## Non-goals` sections in candidate prompt

What: loop.go:196-197 emits `## Constraints\n## Non-goals\n` with no body
under either header. Every loop-drafted candidate ships with two empty
sections.

Why it matters: Purely cosmetic — downstream parsers do not break on empty
sections, but it is untidy and inconsistent with hand-authored prompts that
populate these.

Concrete fix: Either omit the empty headers or leave a placeholder
`(to be filled on promotion)`. Low priority.

### MINOR — F4: Explicit-fingerprint collision silently drops a distinct signal

What: Two genuinely distinct signals that share the same explicit
`fingerprint` (e.g. `ci`/`build-1` and `ci`/`build-2` both with
`fingerprint: "flaky"`) produce the same slug; the second is added to
`Skipped` and its content (title/detail/id) is lost.

Why it matters: This is documented as by-design ("Fingerprint dedupes",
FINAL.md LE-9), so it is not a bug per se. But a signals-file author who
reuses a fingerprint across distinct signals gets silent data loss with no
warning. For the MVP this is acceptable; flagging for awareness.

Concrete fix: None required for the MVP. If desired later, log a warning when
a skip is due to a fingerprint match on a signal with a different `id` than
the one that created the dir (requires persisting provenance, out of scope).

## Open questions

1. Is the signals file (`loop/signals.json`) ever fed by an untrusted source
   in the intended deployment? If a CI hook or MCP trigger can write to it,
   F1 is an active exploit path, not just a latent invariant violation. If it
   is always human-authored, F1 is still a correctness bug (the brake's output
   invariant is violated) but the blast radius is smaller. The FINAL.md
   "Out of scope" lists live connectors as a follow-up, which suggests
   untrusted sources are a future concern — making F1 a fix-now item.

2. Should `ReadFrontmatter` (last-wins) and `readFrontmatterField`
   (first-wins) be reconciled? They disagree on duplicate keys. F1's
   `status: round-01` injection is read as `round-01` by `readIdeas`
   (last-wins) but as `candidate` by `readIdeaStatus` (first-wins), so the
   driver and the status list would disagree on the same idea. This is a
   pre-existing inconsistency surfaced by F1, not introduced by this change —
   but it amplifies the finding.
