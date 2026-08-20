---
idea: agents-verify-hermes-probe
status: open
track: standard
initiator: claude-1
date: 2026-08-20
participants: [claude-1, codex-1, kimi-1, zcode-1]
rounds: 1
---

# `agents verify --full --agent hermes` fails while hermes itself works

## The observation

[PRIMARY, 2026-08-19, parley 1.45.0]

```
$ parley agents verify --full --agent hermes --yes
hermes: installed version=Hermes Agent v0.20.4 (2026.8.18)
probe dir: parley-deck/meta/runtime-probes/20260819T110813.911589000Z
hermes: full probe failed: hermes headless probe did not create
  parley-deck/meta/runtime-probes/.../hermes.md: no such file or directory
verify --full failed: 1 full verification probe(s) failed: hermes
```

The same binary, invoked directly, answers correctly:

```
$ hermes --yolo --oneshot "Reply with exactly: HOK" --model fireworks/inkling --reasoning high --accept-hooks
HOK        (exit 0)
```

`zcode` passes the same command shape: `parley agents verify --full --agent zcode --yes` →
`zcode: headless probe passed`.

## What is already known, and what is only suspected

**Known (PRIMARY).** A hermes process launched with `HERMES_HOME` pointing at an **empty** directory
fails every API call with `API call failed after 3 retries: Gemini returned HTTP 404` — three runs,
with and without `--model`, identical. With the real `~/.hermes` it returns `HOK` at exit 0.

**Known (PRIMARY, source read).** parley does **not** hand hermes an empty home:
`isolatedHermesHome` (`internal/runner/runner.go:1259-1289`) creates the temp base, makes
`logs/ sessions/ home/`, and copies `config.yaml`, `.env`, `auth.json` and `SOUL.md` from
`~/.hermes`. So the naive empty-home explanation does **not** apply to parley's own probe.

**Suspected, UNVERIFIED — do not repeat as a cause.** That the probe failure is a *file-writing*
failure rather than an API failure. Evidence pointing that way and nothing more: in a separate
Parley round, hermes produced a complete, correct analysis on stdout and **wrote no file**, twice,
once stating it had run out of tool calls. `config.yaml` sets `agent.max_turns: 90`.

## What this idea must establish

1. **Where the probe dies.** Capture hermes's stdout and stderr from inside the probe, not from a
   hand-rolled invocation. Does it error, or does it succeed and simply not write `hermes.md`?
2. **Whether the seeded isolated home is actually complete.** `isolatedHermesHome` copies four
   named files. Verify against a working `~/.hermes` that nothing else is load-bearing — the
   `isolated_home_env` map in `parley-deck/agents.toml` sets `HERMES_HOME={tempdir}` and the deck
   also sets `approval_policy = "accept-hooks"` and `isolate_home = true`.
3. **Whether hermes writes files at all under `--yolo --oneshot`.** This is the highest-value
   question: if it cannot, hermes cannot produce a Parley artifact, and every round it has
   "participated" in needs re-examination. Test it directly and cheaply.
4. **Whether `isolate_home` is worth keeping for hermes.** It exists for writable logs and
   parallel-run safety. If it is the cause, the options are seed it correctly, drop it, or make
   the probe report the underlying error instead of only the missing file.

## Constraints

- Track `standard`: tooling defect, no protocol change, no security surface.
- Repository is READ-ONLY to participants except each participant's own round file. Verification
  runs in a COPY.
- The probe writes under `parley-deck/meta/runtime-probes/`; do not delete other runs' probe dirs.
- English only. No secrets — `~/.hermes/config.yaml` contains a live API key. **Never print, quote
  or copy it.** Redact any output that would include it.

## Why this is filed separately

Found while running `roster-membership-overlay`. It is unrelated to that idea's subject and would
have distorted it. The owner asked for it to be its own idea.

---

## RESOLVED — 2026-08-20, root cause and fix

**Root cause: the probe handed the agent a RELATIVE path.**
`runFullVerification` built `probeDir` by joining a possibly-relative `root`, so the prompt said
`parley-deck/meta/runtime-probes/<id>/hermes.md`. hermes resolves relative paths against `$HOME`
whatever the process cwd is — [PRIMARY] measured three ways against Hermes Agent v0.20.4:

```
$ hermes --yolo --oneshot "run pwd, reply with only its output" …
/Users/tomasfecko                       # cwd was the repo

$ hermes --yolo --no-restore-cwd  …     → /Users/tomasfecko
$ hermes --yolo --in <repo>       …     → /Users/tomasfecko

$ hermes … "create hermes-path-test.txt with a relative path, reply with the absolute path"
Wrote to relative path `hermes-path-test.txt`. Absolute path: `/Users/tomasfecko/hermes-path-test.txt`
```

**Neither `--in` nor `--no-restore-cwd` changes it**, so no flag on the hermes side fixes this.

**Fix (`internal/app/app.go`):** `probeDirFor(root, runID)` resolves the probe directory with
`filepath.Abs`. Every adapter is now told exactly where to write instead of being assumed to share
parley's cwd.

```
$ parley agents verify --full --agent hermes --yes
probe dir: /Volumes/.../parley-deck/meta/runtime-probes/20260820T082615.468009000Z
hermes: headless probe passed
```

**Question 3 of this prompt — "does hermes write files at all?" — answered: yes, always.** The
write-failure framing was wrong from the start. Three artifacts believed lost were recovered from
`~/parley-deck/` and moved into the repository, including hermes's round-3 artifact for
`roster-membership-overlay`, whose result had been carried into `consensus.md` as SECONDARY quoted
from a log because the file "did not exist".

**Guarded by** `internal/app/probe_abs_test.go`. Reversion check run in an isolated copy of the Go
tree with the fix removed and the helper kept: the test fails on `""` and on `some/relative/deck`,
as it must.

**Not fixed here, deliberately:** the verifier's message. *"did not create <path>"* is literally
true and materially misleading — it never says where it looked. Same class as D-A and D-C, and it
belongs with those.

**Separate defect, filed on its own:** hermes's tie-break in `protocol-and-skill-audit/round-02`
reports fabricated command output under a `PRIMARY` tag — a `python -m roster` module and a
`meta/version.json` path that both do not exist. That is not a path bug and must not be folded in.
