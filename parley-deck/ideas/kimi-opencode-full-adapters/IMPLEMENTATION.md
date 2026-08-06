---
idea: kimi-opencode-full-adapters
implementer: claude-1
status: complete
cycle: 1
date: 2026-08-06
---

## What changed

One file of product code: `internal/agents/discover.go` gains two full `Spec` entries in
`DefaultSpecs()`. One test file: `internal/agents/autonomous_test.go` extends its contract.

| Field | `kimi` | `opencode` |
|---|---|---|
| `LaunchMode` | `LaunchHeadless` | `LaunchHeadless` |
| `HeadlessArgs` | `["-p", "{prompt}"]` | `["run", "--auto", "{prompt}"]` |
| `PromptMode` | `PromptArg` | `PromptArg` |
| `ApprovalPolicy` | `print-mode default` | `auto` |
| `Model` / `Reasoning` | `CLIDefault` (config-driven) | `CLIDefault` (config-driven) |
| `AutonomousWrite` | `Mode "prompt", Args ["-p"], Scope ""` | `Mode "auto", Args ["--auto"], Scope ""` |

`mergeACPCatalog` still merges the ACP entries onto these specs, so `kimi acp` / `opencode acp`
remain available as an alternative launch mode rather than being replaced.

## Evidence — every spec field was probed, not recalled

The user's instruction for this change was: *verify all your notes, do not rely on history, agents
may behave differently in new versions.* Every claim below is `PRIMARY`, executed 2026-08-06.

| # | Probe | Result |
|---|---|---|
| 1 | `kimi --auto -p "<write a file>"` | **exit 1** — `error: Cannot combine --prompt with --auto.` File not written |
| 2 | `kimi -p "<write a file>"` | **exit 0**, file written unattended |
| 3 | `opencode run --auto "<write a file>"` | **exit 0**, file written |
| 4 | `opencode run "<write a file>"` | **exit 0**, file written — the default already permits |
| 5 | `kimi --help` | `--auto` = "fully autonomous"; `-y/--yolo` = "auto-approve regular tool calls"; `-p/--prompt` = "run one prompt non-interactively" |
| 6 | `opencode run --help` | `--auto` = "auto-approve permissions that are not explicitly denied (dangerous!)"; `--variant` = reasoning effort |
| 7 | `opencode models` | `litellm/xai/grok-4.5` and `litellm/glm-5p2` both listed |
| 8 | `~/.kimi-code/config.toml` | `default_model = "kimi-code/k3"`, `[thinking] effort = "max"` |

Probes 1 and 2 are why `kimi`'s autonomous mode is declared as `prompt` and not `auto`: **`--auto`
exists but is rejected alongside `-p`**, so `-p` is the only autonomous headless shape kimi has.

Probes 3 and 4 are why the opencode choice was put to the user rather than decided here: both
shapes work, and they differ in permission breadth, not in outcome.

## Corrections to the machine's own notes, forced by this verification

`~/.parley/agents.toml` carried claims that no longer hold. They are **not** fixed by this commit
(that file is machine-local and gitignored) but they are recorded so the record is not left wrong:

| Note | Verdict |
|---|---|
| "no headless autonomous-write contract for kimi" | **WRONG** — kimi ships `--auto` and `-y/--yolo`; they merely cannot combine with `-p` |
| "K3 `default_effort = "max"`" | **WRONG** — the model's own `default_effort` is `"high"`; `max` comes from the `[thinking]` override |
| "kimi is 0.32.0" | **stale** — 0.33.0 |
| "opencode is 1.18.13" | **stale within the hour** — 1.18.14 |
| "`parley roster init` fail-closes and would DROP `[roster.kimi-1]`" | **NOT REPRODUCED** — on a disposable deck seeded with `kimi-1` + `opencode-1`, both the installed 1.38.0 binary **and** the new build proposed both mappings and dropped nothing, exit 0. The warning describes parley 1.36 |

The last row matters beyond this idea: that warning is the reason the roster has been hand-edited
for weeks. It is the same defect the previous idea recorded as `T3 — NOT REPRODUCED at 1.37.0`.

## Verification

| # | Check | Result | Kind |
|---|---|---|---|
| 1 | `go build ./...` | OK | survival |
| 2 | `go test -count=1 ./...` (implementer's environment) | rc=0, 25 packages | survival |
| 3 | `go test ./internal/agents/...` incl. `TestBuiltinsAreAutonomous` | ok, with `kimi`/`opencode` added to the contract | **fix-proving** |
| 4 | `agents list` on the new build | `kimi` → `LAUNCH=headless AUTO=yes`, `headless: kimi -p`; `opencode` → `LAUNCH=headless AUTO=yes`, `headless: opencode run --auto` | **fix-proving** |
| 5 | `roster init --dry-run` on a disposable deck, old vs new binary | both propose `[roster.kimi-1]` + `[roster.opencode-1]`, exit 0 | control |
| 6 | `gofmt -l internal/agents/` | no output | survival |

Check 3 and 4 would fail at the parent commit — that is what makes them fix-proving. Check 5 is a
**control**: it shows the change did not fix the roster-init behaviour, because that behaviour was
not broken to begin with.

**Not verified:** that a `parley`-launched kimi or opencode completes a full protocol round. The
probes prove the CLIs write unattended under these exact argv shapes; they do not prove the runner
drives a whole round with them. That is the main thing review should push on.

## Not in scope

- `~/.parley/agents.toml` corrections (machine-local, gitignored).
- Removing the stale roster-init warning from that file.
- Any other adapter.

---

## Cycle 1 — fix-up after review round 1

All three signoffs on `review/consensus.md` were ✅ and all three accepted **targeted
verification** instead of a second review round.

### AF-1 — the MAJOR, and it was wider than the idea

`internal/config/runtime.go:542` replaces `spec.HeadlessArgs` wholesale when a config layer
supplies them and never touches `AutonomousWrite`. So a spec could declare an autonomous mode
whose enabling flag the launched command never passes, while `agents list` printed `AUTO=yes` and
a `headless:` line taken from the built-in **label** rather than the effective argv.

**Three parts, all applied:**

1. **Fail closed.** New `AutonomousWrite.MissingFrom(headlessArgs)` and `Spec.AutonomousEffective()`.
   The `AUTO` cell in both `agents list` (`discover.go`) and `roster show` (`roster.go`) now uses
   the effective reading. A declared mode the launch does not enable reports `no`.
2. **Show the truth.** The `headless:` detail line is now built from the resolved binary plus the
   **effective** `HeadlessArgs`, and prints an explicit warning naming the missing args.
3. **Repair the config** (machine-local, backed up to `agents.toml.bak-2026-08-06-pre-autofix`):
   `--auto` added to the central `opencode` override, `--yolo` restored to `hermes` in **both**
   the central and the in-repo deck override.

**The defect was live before this idea.** With the fix in place and before the config repair:

```
hermes   ... AUTO=no
  headless: …/hermes --oneshot {prompt} --model glm-5p2 --accept-hooks
  ** autonomous mode "yolo" declared but its enabling args [--yolo] are absent from the launch **
opencode ... AUTO=no
  headless: …/opencode run -m litellm/xai/grok-4.5 {prompt}
  ** autonomous mode "auto" declared but its enabling args [--auto] are absent from the launch **
```

The opencode line is the exact effective argv `codex-1` derived from source. **That is the
measurement that settled the conflict**: the implementer's contrary probe had run inside `go test`,
where `PARLEY_HOME` is isolated and the real central config is never read. Withdrawn.

### AF-2 to AF-6

| AF | Fix |
|---|---|
| AF-2 | `TestPromotedAdaptersFullContract` locks autonomous args, headless argv, `PromptArg`, `LaunchHeadless`, empty `Scope`, `Confined()==false`, and `AutonomousEffective()` for both IDs. `TestAutonomousFailsClosedWhenLaunchDropsTheFlag` is the AF-1 regression, written against the exact shape that was live |
| AF-3 | `TestPromotedAdaptersKeepACPAsAlternative` asserts `ACPArgs == ["acp"]`, `LaunchHeadless`, non-empty `HeadlessArgs` and a declared mode survive `mergeACPBackend` |
| AF-4 | kimi's `Notes` now state that the official installer puts the binary at `~/.kimi-code/bin/kimi` and does **not** add it to PATH, so `command` must be set |
| AF-5 | opencode `Telemetry` → `"streamed build/tool events, then final text on stdout"` |
| AF-6 | The in-repo deck guidance for kimi is rewritten: the "do NOT run `parley roster init`" warning is removed as false on current builds, the redundant `headless_args` pin is dropped (a redundant override is exactly how hermes lost `--yolo`), and the `reasoning` comment now says `max` comes from `[thinking]`, not from the model's own `default_effort`, which is `"high"` |

### Targeted verification

| # | Check | Result | Kind |
|---|---|---|---|
| 1 | `go build ./...` | OK | survival |
| 2 | `go test -count=1 ./...` | rc=0 | survival |
| 3 | `TestAutonomousFailsClosedWhenLaunchDropsTheFlag` | passes; asserts `AUTO=false` with the flag stripped and `true` once restored | **fix-proving for AF-1** |
| 4 | `agents list` with the config **unrepaired** | `hermes` and `opencode` → `AUTO=no` + explicit missing-args warning | **fix-proving for AF-1** |
| 5 | `agents list` after the config repair | all five agents `AUTO=yes`, effective args carry the declared flag | control |
| 6 | `TestPromotedAdaptersFullContract`, `TestPromotedAdaptersKeepACPAsAlternative` | pass | fix-proving for AF-2/AF-3 |

Check 4 is the one that matters: it reproduces the defect on the shipped configuration before the
repair, which is what makes the fail-closed behaviour a fix rather than a claim.

**Still not verified:** that `parley` drives a complete protocol round with either agent. The
reviewers did not require it and the launch path is now asserted at the spec level, but it remains
the gap named in cycle 0.
