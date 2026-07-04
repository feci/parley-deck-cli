---
idea: parley-learn-playbooks
status: fix-up-cycle-1
implementer: claude-1
started: 2026-07-04
completed: 2026-07-04
branch: parley-deck-cli#learn-playbooks-design
head-commit: ae84237
design-pr: https://github.com/feci/parley-deck-cli/pull/70
implementation-pr: same
---

## Summary of work

`parley learn <closed-idea-slug>` distills a completed idea into an advisory playbook at
`parley-deck/playbooks/<topic>.md`, mirroring `parley retro propose`'s safety. Protocol
§13.5 defines playbooks as advisory (beside consults). No new phase, gate, or quorum.

## Implementation plan / checklist

- [x] `internal/app/learn.go` (new): `runLearn` — strict slug + `--topic` (parsed with
      the slug pulled out first so `learn <slug> --dir X` works), `status: complete`
      precondition, `Lstat` fail-closed target guard, `distillPlaybook` deterministic
      skeleton from the idea's frontmatter + lifecycle.
- [x] `internal/app/app.go`: `case "learn": runLearn(...)`.
- [x] Protocol §13.5 (BOTH `parley-deck/COOPERATION.md` and
      `internal/protocol/defaults/COOPERATION.md`, byte-identical) + §13 changelog line;
      skill fallback `references/COOPERATION.md` re-synced (body-identical from line 7).
- [x] `meta/version.json` protocolSha256 refreshed; `meta/protocol-changelog.md` entry.
- [x] Tests `internal/app/learn_test.go`: writes playbook; rejects incomplete; fails
      closed on existing target; flag-after-slug. Drift guard green.
- [x] Checks: `go build ./...`, `go vet`, `gofmt -l` clean.

## Deviations from FINAL.md

None. v1 distillation is the deterministic skeleton agreed in consensus (a human refines
the transferable prose before committing the playbook).

## Notes for reviewers

- `runLearn` pulls the first non-flag arg (slug) out before `flag.Parse` so flags after
  the slug are honored (Go's flag stops at the first positional) — `TestLearnFlagAfterSlug`.
- Fail-closed target guard uses `Lstat` (a symlink cannot redirect the write) — matches
  `parley retro propose`.
- The playbook is advisory: `runLearn` prints that it is never quorum and never overrides
  the protocol.

## Fix-up cycle 1
status: complete
completed: 2026-07-04

### Fixes applied (review round-01)
- [MAJOR, codex-1] Hardened the write boundary: Lstat the `playbooks/` PARENT and refuse
  a symlink (a symlinked parent could redirect the write outside the workspace even when
  the final target Lstat sees it absent); create the dir only when genuinely absent; and
  create the playbook with `os.OpenFile(O_CREATE|O_EXCL|O_WRONLY)` (atomic exclusive
  create closes the check-then-write TOCTOU and refuses an existing target). Test
  `TestLearnRejectsSymlinkedPlaybooksDir` (symlinked parent fails and does not write into
  the target).
- [MINOR, hermes-1] §13.5 softened: "scaffolds … a deterministic skeleton … that the
  author refines into transferable prose before committing" — no longer oversells the v1
  distillation as full auto-capture.
- [MINOR, hermes-1] `playbooks/` absent from the §3 directory-layout tree: dismissed as
  consistent with `consults/` (also advisory, also not in the tree).
- [NIT, hermes-1] Wording: the skill fallback lives in the sibling `parley-deck-skill`
  repo (not this repo); it is re-synced there as part of the release.

### Deviations from agreed fixes
None.
