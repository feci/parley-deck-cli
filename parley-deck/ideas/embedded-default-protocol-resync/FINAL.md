---
idea: embedded-default-protocol-resync
status: final
drafted-by: claude
date: 2026-06-13
participants: [claude, codex, agy, hermes]
consensus: 4/4 ACCEPT (D7 disposition concurred by codex, hermes, agy)
---

## Decision

Resync the embedded default protocol (`internal/protocol/defaults/COOPERATION.md`,
the `parley init` bootstrap template) with the live project deck, genericize the
project-specific zones it should never have shipped to new projects, and add a Go
drift guard so the two copies cannot silently diverge again.

## What ships

1. **§12 verbatim.** Append `## 12. Pipeline blocks & action stages` to the
   embedded default after §11, byte-identical to `parley-deck/COOPERATION.md`
   (incl. the closing "ratified by idea `meta-protocol-change-end-to-end-pipeline`
   (2026-06-02)" sentence and the exact final newline).

2. **Header genericization** (embedded default only):
   - `**Workspace:** \`<workspace-name>\``
   - `**Created:** \`<date> — created by parley init\``
   - keep `**Transport:** \`github-pr\``
   - do NOT add a `**Protocol synced:**` line.

3. **Empty §2 tables.** Empty the bodies of both §2 tables (roster + host-handle),
   retaining their header and separator rows exactly. New projects start with no
   quorum members.

4. **InitWorkspace unchanged.** `defaultCooperationForInit()` keeps performing only
   the `github-pr` → `local-dir` transport swap. No dynamic workspace/date
   rendering, no roster discovery in this idea.

5. **Two Go tests in `internal/protocol`:**
   - **(A) Drift guard** — embedded default ≡ `parley-deck/COOPERATION.md` after
     normalizing exactly five anchored zones and no others: the deck-only
     `**Protocol synced:**` line, the `**Workspace:**` value, the `**Created:**`
     value, the §2 roster table body, the §2 host-handle table body. Anchored
     normalization (header-line prefixes, the `## 2. Active agents (roster)`
     heading, both table header+separator rows, the prose anchors bounding each
     body) — not broad regexes. **Fail closed** on any missing/duplicate anchor or
     heading, or if the deck file is absent. The allowlist lives in one helper and
     is named in the failure message — the documented home of the "edit both
     copies" invariant.
   - **(B) Init output** — `defaultCooperationForInit()` emits
     `**Transport:** \`local-dir\``, contains §12 incl. its provenance line,
     contains the `<workspace-name>` / `created by parley init` placeholders,
     contains no `**Protocol synced:**` line, and contains none of the parley-deck
     roster rows.

## Explicitly not in scope (deferred follow-ups)

- A `go generate` generator deriving one copy from the other (only if the guard
  proves noisy).
- A `parley protocol check` subcommand wrapping the normalizer.
- Dynamic `parley init` header rendering + roster discovery (a UX idea).
- A §7 protocol-text pointer about the dual-copy invariant (meta-protocol-change;
  D7 — dropped here by 4/4 concurrence).
- Resync of the out-of-repo packaged skill reference (kindly inbox note).

## Implementer

claude (initiator). Branch `feature/embedded-default-resync`, transport github-pr.
Verification gate: `go build ./...` + `go test ./internal/protocol/` green, plus
the full suite before ship.
