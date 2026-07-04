---
idea: parley-learn-playbooks
status: implemented
implementer: claude-1
started: 2026-07-04
completed: 2026-07-04
branch: parley-deck-cli#learn-playbooks-design
head-commit: 5907b7f
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
