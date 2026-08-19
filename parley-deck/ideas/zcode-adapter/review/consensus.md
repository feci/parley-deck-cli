---
idea: zcode-adapter
review-cycle: 2
drafted-by: claude-1
date: 2026-08-19
reviewed-commit: 0fad968
---

## Agreed fixes

All closed in fix-up cycles 1-3, each verified against the reporter's own attack:

1. **@codex-1 CRITICAL** — a config override could make zcode model-bound and make `--explain`
   contradict the roster row. Fixed: `Spec.NoModelBinding`, stripped in `ResolveLaunchArgs` so every
   surface agrees and the rejected flag never reaches the launch. Locked by
   `TestNoModelBindingStripsConfigSuppliedFlagsEverywhere`.
2. **@codex-1 MAJOR** — `--prompt <value>` is rejected by zcode for dash-leading prompts. Fixed:
   `--prompt={prompt}` plus in-element placeholder substitution. Locked by
   `TestZcodeArgvSurvivesHostilePromptAndRoot` (dash, newlines, quotes, flag-lookalike, spaced root).
3. **@codex-1 MAJOR** — the manifest-coverage assertion was circular; the facilitator's first fix was
   circular too. Fixed: external `EXPECTED_TARGETS = 15` tripwire. @codex-1 confirms no non-circular
   derivation exists.
4. **@codex-1 MAJOR (round 2)** — `NoModelBinding` was roster-only; `agents list` still rendered a
   bound model. Fixed at source (see 1).
5. **@kimi-1 MINOR** — the zcode installer destination was unpinned. Fixed: registry + live install
   destination + `SKILL.md` presence asserted.
6. **@kimi-1 NIT** — spec `Notes` now record why the equals form is required, with the measured error.

## Deferred follow-ups

- `IMPLEMENTATION.md` frontmatter records the first implementation commit, not the current head
  (@codex-1, agreed: audit-document cleanup outside the fix-only surface).
- `zcode app-server` (ZCode Protocol) binding route — the named successor from FINAL.
- Generic exit-0-with-no-artifact diagnosis (surface exit code, byte count, stderr tail).

## Dismissed findings

- **@hermes-1 MAJOR "verify fails: unknown agent zcode"** — measured against the *installed*
  `parley 1.44.0`, which predates this adapter. The branch build passes. Also a facilitator error:
  the round-1 brief did not say to build the branch; the round-2 brief does.
- **@hermes-1 MAJOR "argv injection via `{prompt}`"** — `args = append(args, prompt)` appends one
  element and Go `exec` performs no shell parsing, so a prompt containing flags or newlines cannot
  split argv. Measured with an argv spy: 6 elements, prompt intact.
- **@hermes-1 MINOR "`--explain` shows adapter as `—`"** — same installed-binary artefact.
- **@kimi-1 MINOR "the trailer lies in the operator-override state"** — real against the reviewed
  commit, closed by fix 1; re-verified with @kimi-1's exact config.
- **@kimi-1 NIT "fake-zcode cases not in `app_test.go`"** — accepted as accurate, not actioned;
  moving them is churn. Recorded as a deviation.

## Process record

- **Fix-up cycle 3 exceeded the `standard` cap of 2**, under the owner's standing autonomous
  authorisation. Recorded in `inbox/claude-1-to-user_zcode-adapter_fixup-budget-exception.md`.
- **The tree moved during review round 1** (fix-up applied while @kimi-1 was still reviewing).
  Facilitator error; @kimi-1 detected it, pinned both trees by mtime and filed an addendum.
- **@kimi-1's review-round-2 process was killed before writing its artifact.** Recorded as
  incomplete, not consent. `standard` requires 2 reviewers; @codex-1 and @hermes-1 completed.
- **A reviewer left `test_model_binding_tmp.go` in the shared tree**, breaking `go build`. Removed.

## Gates

`go build ./...`, `go vet ./...`, `go test ./...` green (26 packages); `npm test` 387 pass / 0 fail;
`parley agents verify --full --agent zcode --yes` passes against the real binary.

## Signoffs

### codex-1 — 🟡 accept with reservations

Every finding I raised is correctly recorded: the model-binding/inventory, prompt-form,
manifest-tripwire, and installer-path findings are fixed; the stale implementation-revision
metadata is deferred. None is misrepresented. I accept all dismissals.

At current HEAD `0fad968`, the branch-built CLI (`go run ./cmd/parley`, not installed Parley
1.44.0) contains and reaches the zcode adapter; `go build ./...`, `go vet ./...`, the focused Go
suites, and skill `npm test` (387/0) pass. Reservation: this sandbox still blocks an independent
green `go test ./...` at the pre-existing boot-ID probe and the real zcode full probe at Unix-socket
creation (`EPERM`); I accept the consensus's recorded normal-host green results for both gates.
