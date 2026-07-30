---
idea: skills-cli-install-path
review-cycle: 21
drafted-by: claude-1
date: 2026-07-30
reviewed-commit: a544dcd
---

## Agreed fixes

**None.** Review round 21 was a unanimous accept — `codex-1` ✅, `hermes-1` ✅, `kimi-1` ✅ — and
every finding raised in it was labelled a follow-up by the reviewer who raised it.

`codex-1`'s release decision, quoted because it answers the question this cycle asked:

> nothing in cycles 26–28 changes the shipped product, every round-20 false green is now
> refused, the intended controls remain green, a real bad command still runs and fails, and all
> project-owned verification commands execute successfully.

One NIT from round 21 — raised by `codex-1`, and **confirmed** by `hermes-1` with its own `sh`
and `zsh` measurements (A1: the first draft called this an independent discovery) — was **taken
rather than deferred** in cycle 29: the shared detector comment still claimed both brace forms build a shell
word, when `hermes-1` had measured that `--{test..test}` does not expand. Three lines, and a
comment asserting something the shell does not do is the class of error the last two rounds were
about.

## Deferred follow-ups

1. **Glob characters (`?`, `*`, `[…]`) as word-building constructs** — `hermes-1` (MINOR) and
   `kimi-1` (NIT), independently. With a file named `node` in the working directory,
   `n?de --test x` really runs; detection has no construct for it. Both marked it a follow-up
   because glob expansion is *conditional on filesystem state* where brace expansion is
   unconditional. The second ground — that no shipped file puts a glob character in the binary
   or flag position — is `hermes-1`'s evidence. `kimi-1`'s was the filesystem precondition
   itself, which it measured: no file named `node` at the repo root, and none shipped by
   `npm pack --dry-run`. (A2.)

   **The cost of the obvious fix is measured, not assumed.** Adding `?`, `*` and `[` to the
   predicate immediately refuses a real shipped document:

   ```text
   skills/parley-tracker/templates/epic.md
   AC-E2 [B][T] Measurable: … Verify: `node skills/…/validate.js --strict --dir tickets`
   ```

   `[B][T]` carries a `[`, the line names `node`, and the rule fires on an innocent sentence.
   Any future fix must clear that case.

2. **The occurrence rule is order-dependent** — `kimi-1`, round 18. A visible line where
   `--test` precedes `node` is not judged. It is not a runnable command, so this is a question
   about documenting the rule's scope, not a hole in it.

3. **Detection approximates shell word splitting rather than lexing** — `kimi-1` named lexing as
   the consistent move (execution was delegated to `/bin/sh` in cycle 9); the approximation was
   chosen because lexing a line through a shell risks evaluating the substitutions it exists to
   inspect. `hermes-1` judged the approximation sound in round 19: removing quoting characters
   can only join characters into tokens, never split them, so it can reveal a command but never
   conceal one.

4. **Pre-existing installer defect (round 01, `codex-1`)** — in a fresh `HOME`, the Antigravity
   install is what creates the evidence by which Gemini is then detected, so
   `install --target all` must run **twice** before `doctor` is clean. Verified at release time
   to be that defect and not a regression. It predates this idea and remains open against the
   installer.

## Dismissed findings

None. Every finding raised across 21 review rounds was either fixed or explicitly deferred
above.

## Cycle record

21 review rounds, 29 fix-up cycles, four reviewers, two reviewer outages: `codex-1` was
terminated mid-run twice by an upstream content filter (rounds 14 and 19, both recorded as
outages rather than accepts), and `agy-1` exhausted an account quota after two rounds.

The shipped product — the layout move, the installer, the packaging, the README panel — was
settled by round 02 and never changed after it. Everything from round 03 onward was the
documentation guard.

## Amendments after signoff

Two attribution corrections were made to the body **after** `codex-1` signed. Each was requested
by the signer it concerned, and neither changes a verdict, finding, severity or disposition:

- **A1** — `hermes-1`: the round-21 NIT was `codex-1`'s discovery, which `hermes-1` confirmed;
  the draft had called it independent. The correction credits `codex-1` more, not less.
- **A2** — `kimi-1`: the second ground for deferring the glob finding is `hermes-1`'s evidence,
  not `kimi-1`'s.

Recorded here rather than applied silently, because `codex-1` signed the pre-amendment text and
did not see either change.

## Signoffs

### Signoff: codex-1 — 2026-07-30
Status: ✅ ACCEPT

The zero-item agreed-fixes list matches my round-21 release decision: all round-20 blockers were closed and my only new finding was a non-blocking NIT. The draft correctly records that cycle 29 fixed that stale brace-predicate comment rather than leaving it deferred. The deferred list accurately preserves the pre-existing fresh-HOME Antigravity/Gemini detection-order defect I reported; the other deferred items do not conflict with my review. The draft represents my position faithfully.

### Signoff: hermes-1 — 2026-07-30
Status: ✅ ACCEPT

The zero-item agreed-fixes list and the unanimous-accept characterization match my round-21 review exactly: I had no blocks, only a MINOR (glob patterns) and a NIT (stale brace-predicate comment), both labelled follow-ups. Deferred item #1 faithfully reproduces my glob finding — conditional on filesystem state, no shipped file approaches it. One small caveat: the draft says codex-1 and I found the `--{test..test}` NIT "independently," but my review explicitly states I was confirming codex-1's prior identification ("codex-1 already identified this in their round-21 review. I confirm it."), not discovering it separately. This does not change my accept or my assessment of the NIT's severity, but the record should not overstate the independence of my observation.

### Signoff: kimi-1 — 2026-07-30
Status: ✅ ACCEPT

The zero-item agreed-fixes list matches my round-21 review: I verified all three round-20 findings closed, re-ran both halves of the cycle-26 reproduction myself, and filed exactly one new finding — the glob NIT, explicitly labelled a follow-up, with "nothing I found stops this from being published today." Deferred item #1 represents that finding faithfully on the substance: a glob-built command word is conditional on the reader's filesystem where brace expansion is unconditional, and the guard can never certify one green. One precision for the record: my deferral evidence was the filesystem precondition — I checked that no file named `node` exists at the repo root and that `npm pack --dry-run` ships none, with a loud `command not found` where the precondition fails; I did not myself sweep shipped documentation for glob characters in command position, so the draft's second "because" clause describes hermes-1's evidence, not mine. Deferred items #2 and #3 match positions I took in earlier rounds, and the cycle-29 comment fix closes the flag-arm NIT I also confirmed with my own sh and zsh measurements. The release changes none of this; my accept stands.
