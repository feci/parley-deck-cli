---
idea: zcode-adapter
status: implemented
implementer: claude-1
started: 2026-08-19
completed: 2026-08-19
branch: parley-deck-cli#zcode-adapter-impl
head-commit: 9a6f8be
design-pr: n/a
implementation-pr: n/a
---

## Summary of work

`zcode` is now a built-in adapter. `parley agents list` prints its `headless:` argv line,
`roster show` reports AUTO=yes, and the acceptance command from FINAL passes against the real
binary:

```
$ parley agents verify --full --agent zcode --yes
zcode: installed version=zcode-app-cli 3.7.7-13
probe dir: parley-deck/meta/runtime-probes/20260818T232832.862162000Z
zcode: headless probe passed
```

MODEL and EFFORT stay `unknown` / `model-unbound`, as designed.

## Implementation plan / checklist

- [x] `internal/agents/discover.go` — `Spec{ID: "zcode"}`, argv
      `--prompt {prompt} --mode yolo --cwd {root}`, `AutonomousWrite{Mode:"yolo",
      Args:["--mode","yolo"], Scope:""}`, `Model/Reasoning: CLIDefault`.
- [x] `internal/agents/modelmeta.go` — `"zai": "Zhipu AI"` producer (zcode emits `zai/<model>`,
      no hyphen; `z-ai`/`zhipu` did not match).
- [x] `internal/app/roster_view.go` — static `model source:` trailer on `--explain`.
- [x] `internal/agents/launchargs_test.go` — lock: no `{model}`/`{effort}` placeholder, exact
      argv, AutonomousWrite args present in argv, Scope not "workspace".
- [x] `internal/agents/modelmeta_test.go` — `zai/glm-5.3` → GLM / Zhipu AI.
- [x] `internal/app/zcode_verify_test.go` — @codex-1's two fake-zcode full-verify cases plus a
      trailer test.
- [x] `internal/agents/acp_specs_test.go`, `internal/app/roster_test.go` — deliberate presence locks.
- [x] Skill: `lib/installer.js` zcode target, `skills/parley-deck/SKILL.md` autonomous-write row,
      manifests rebuilt.
- [x] `~/.parley/agents.toml` — `[agents.zcode]` removed, roster `model`/`effort` removed,
      stale exit-0 note corrected.
- [x] Checks run: `go build ./...`, `go test ./...` (all green), `npm test` (386 pass / 0 fail),
      real `agents verify --full --agent zcode --yes`.

## Deviations from FINAL.md

**D-1 — `.zcode/skills` was verified, not assumed.** FINAL required a native installer target but
not a path. `~/.zcode/` has no user skills directory on this machine (all 14 skills come from the
plugin cache), so the path could not be read off the filesystem. It was taken from the literal
`.zcode/skills` present in the `zcode-app-cli` runtime bundle, corroborated by
`~/.zcode/cli/config.json` exposing `skills.roots`. Recorded because an unverified install path
would fail silently.

**D-2 — one extra file changed, in the same defect class the idea is about.**
`test/manifest-coverage.test.js` asserted `cores.length === 14` — a hardcoded target count whose
own comment claimed it was "derived from the installer rather than a hand-picked pair". Adding a
15th target failed it. It is now derived from the install result's own target set, so a future
target extends the check instead of breaking an arithmetic assertion nothing keeps in step. Bumping
14→15 would have preserved the defect.

**D-3 — `--explain` shipped in @kimi-1's static form**, as FINAL's Phase-5 default.
@codex-1's labelled live read is not implemented; its reservation stands on the record.

## Notes for reviewers

- The fake-zcode stub asserts **exact token order and `argc == 6`**, so a spec that appends or
  reorders an option zcode does not accept fails in test rather than in production.
- The help-exit-0 case is the important one: it proves full verification cannot be satisfied by a
  process that exits 0 while writing nothing. Exit status alone is never the acceptance signal.
- `agents list` shows a model for zcode when a config layer sets one; `roster show` still reports
  `model-unbound`. That is the contract working — the roster reports what the *launch* carries.
- Not done, deferred by FINAL: the `zcode app-server` (ZCode Protocol) binding route, and generic
  exit-0-with-no-artifact diagnosis.
