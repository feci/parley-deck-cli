---
agent: codex-1
idea: parley-learn-playbooks
review-round: 1
date: 2026-07-04
reviewed-commit: 553ef14
---

## Summary

The implementation largely matches FINAL.md and consensus.md: `learn` dispatch is wired in `internal/app/app.go`; slug/topic validation is strict kebab-case; the command is read-only over the source idea; `IMPLEMENTATION.md` frontmatter `status: complete` is enforced; `learn <slug> --topic ... --dir ...` works despite Go flag parsing; and protocol 13.5 is advisory/non-canonical and byte-identical in both protocol copies. Required checks passed: `go build ./...`, `go vet ./internal/...`, and `go test ./internal/app ./internal/protocol`.

I found one MAJOR write-boundary issue: a preexisting `parley-deck/playbooks` symlink can redirect the new playbook outside the workspace.

## Findings

### [MAJOR] Parent-directory symlink bypasses the playbook write boundary

`runLearn` builds the target at `internal/app/learn.go:72`, `Lstat`s only the final `<topic>.md` path at `internal/app/learn.go:75`, accepts/creates the parent with `MkdirAll` at `internal/app/learn.go:82`, and writes with `os.WriteFile` at `internal/app/learn.go:88`. If `parley-deck/playbooks` already exists as a symlink to another directory, `Lstat(playbook)` sees the final file as absent through that symlink and `WriteFile` creates `<topic>.md` outside the workspace while reporting the normal `parley-deck/playbooks/<topic>.md` path.

I reproduced this in a temp workspace by symlinking `root/parley-deck/playbooks` to an outside directory and running `go run ./cmd/parley learn demo-idea --dir root`; the command succeeded and created `outside/demo-idea.md`. This violates the "writes exactly one new playbook file" / fail-closed write boundary and leaves a path-escape hole despite strict topic validation. There is also a final-path TOCTOU window because the checked path is not opened with exclusive create.

Suggested fix: reject a symlinked `parley-deck/playbooks` parent with `Lstat` before writing; create the directory only when absent; then create the playbook with `os.OpenFile(playbook, os.O_CREATE|os.O_EXCL|os.O_WRONLY, 0o644)` and write through that descriptor. Add tests mirroring `retroPropose`'s symlink coverage: a symlinked `playbooks` directory must fail and must not write into the symlink target, and the final target should be created with an exclusive-open path.

## Open questions

- Should `parley learn --dir DIR <slug>` be accepted too? The current parser is intentionally slug-first and passes the required `learn <slug> --dir DIR` case, but it will treat a value-taking flag's value as the slug if flags are placed before the slug. I did not file this as a finding because FINAL.md specifies slug first.
