---
idea: embedded-default-protocol-resync
author: claude
created: 2026-06-13
participants: [claude, codex, agy, hermes]
roles:
  claude: facilitator + drift/maintenance lens
  codex: build/test-guard + Go-consumer lens
  agy: protocol-content + genericization lens
  hermes: bootstrap-UX + portability lens
status: final
---

## Problem / idea

`internal/protocol/defaults/COOPERATION.md` (the **embedded default**, baked into
the binary via `//go:embed` in `internal/protocol/workspace.go`) has drifted from
the live project deck `parley-deck/COOPERATION.md`. This was flagged as a deferred
follow-up in the 1.24.0 kindly-adoption inbox note
(`claude-to-all_review-gate-honesty_external-skill-snapshot-sync.md`). We just
synced the project deck to parley-deck-skill 1.3.1; now resolve the embedded copy.

### Established facts (verified 2026-06-13)

- The embedded default is the **bootstrap template**: `InitWorkspace(root)` writes
  it to a new project's `parley-deck/COOPERATION.md` *iff* none exists, via
  `defaultCooperationForInit()`, which makes exactly ONE transform:
  `**Transport:** \`github-pr\`` → `**Transport:** \`local-dir\``. Nothing else is
  parameterized.
- `diff internal/protocol/defaults/COOPERATION.md parley-deck/COOPERATION.md`
  yields EXACTLY two differences:
  1. The project deck has a header line
     `**Protocol synced:** 2026-06-13 — parley-deck-skill 1.3.1 (claude)`
     (added during the 1.3.1 sync); the embedded default does not.
  2. The project deck has the full **§12 "Pipeline blocks & action stages"**
     (45 lines, ratified by idea `meta-protocol-change-end-to-end-pipeline`,
     2026-06-02); the embedded default ends at §11.
  Everything else — including the Phase 6 "Review briefs and dispositions",
  Phase 8 strict-gate + "Stopping judgment", and §8 "Consults" amendments from
  1.24.0 — is already byte-identical in both copies (those WERE applied to both
  copies during 1.24.0; only §12 was missed).
- The embedded default's §2 roster table is filled with THIS project's specific
  agents (`codex`, `claude`, `agy`, `hermes`) and the header says
  `**Workspace:** \`parley-deck\``. So `parley init` currently injects
  parley-deck's own roster/workspace into every newly-initialized project, even
  though §2's prose says "The roster is project-specific."

### Design axes to decide

1. **§12 propagation.** Add §12 to the embedded default. It is already-ratified,
   protocol-general content (no project specifics) — presumably a verbatim carry.
   Confirm and specify exactly how (anchor, trailing newline, the closing
   "ratified by idea …" provenance line — keep or drop in the template?).
2. **Genericization scope (the core question).** The embedded default is a
   *template for new projects*, yet it currently hard-codes parley-deck's roster
   rows, workspace name, host-handle table, and now potentially a `Protocol
   synced:` line. Decide, field by field, what is **protocol-general** (carry
   verbatim) vs **project-specific** (should be genericized to placeholders or
   omitted in the template):
   - header `Workspace:` / `Created:` / `Protocol synced:` lines
   - §2 roster table rows + host-handle table rows
   - any other parley-deck-specific string
   Trade-off: a verbatim copy is trivial to keep in sync and to test, but ships
   a misleading roster to new projects; a genericized template is correct for
   bootstrap but reintroduces a structured difference that the anti-drift guard
   must model.
3. **Anti-drift mechanism (highest-value outcome).** How do we stop this from
   silently recurring? Options to weigh: a Go test that asserts
   embedded-default ≡ project-deck modulo a documented allowlist of intended
   differences; making one file the generated source of the other; a
   `parley protocol check` / lint subcommand; or a documented manual checklist
   in §7. Specify what actually ships in THIS idea.

## Constraints

- Do NOT change the semantics of the live `parley-deck/COOPERATION.md`. §12 is
  already ratified; this idea propagates a derived copy and adds a guard — it is
  not a meta-protocol-change to protocol content.
- The embedded default is consumed by `InitWorkspace` with the single
  `github-pr`→`local-dir` transport swap; any genericization must keep that
  consumer correct (or extend it deliberately, with the change recorded).
- English-only for all `parley-deck/` files. Transport is `github-pr`.
- Keep the change set minimal and shippable; this is maintenance + one guard,
  not a redesign.

## Non-goals

- Re-ratifying or editing §12's content.
- Genericizing the out-of-repo packaged skill reference
  (`references/COOPERATION.md`) or the `parley-deck-skill` package — that is a
  separate concern tracked in the original inbox note.
- Adding pipeline (§12) runtime/driver code — §12 is SPEC-only in 1.3.1.
