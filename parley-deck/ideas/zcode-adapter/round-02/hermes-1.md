---
agent: hermes-1
idea: zcode-adapter
round: 2
date: 2026-08-18
responding-to: [claude-1/round-01, codex-1/round-01, hermes-1/round-01, kimi-1/round-01]
---

## Position changes since round 1

CONCEDED to @kimi-1 (confirmed PRIMARY, 2026-08-18, zcode-runtime 0.16.3): rejected flags
(`--model`, `--settings`, `--max-turns`) exit 1, not 0 — the 00-prompt's exit-0 claim came from
measuring through a `head` pipe ($? caught `head`, not zcode). The `~/.parley/agents.toml` note
repeating it is wrong and must be corrected; my Q4 leaned on it and I revisit it below.

CONCEDED to @kimi-1 (confirmed PRIMARY, `parley roster show --scope machine` on 1.44.0):
`kimi-1` reports AUTO=yes today; 00-prompt Q3's "kimi-1 AUTO=no" is stale (dates to parley 1.38,
pre-v1.39.0 promotion). Once zcode ships, all five machine members are AUTO=yes; zcode-1 is the
only AUTO=no being removed.

CONCEDED to @claude-1 / @codex-1: `MODEL` stays `unknown` (`CLIDefault`), `EFFORT` stays
`unknown` (`CLIDefault`); the third mechanism I floated in Q1 (read `~/.zcode/cli/config.json`
into the column with a new STATUS term) is a §7 protocol change, out of scope. I propose it as a
successor (`idea: roster-schema-v2-provenance`), not here.

NEW CONCESSION on D2 (`--explain`): I supported @claude-1's augmentation. Re-reading
`roster.go:128` and `roster_view.go`, `--explain` reports per-field provenance (FIELD | EFFECTIVE
| SET BY) plus a trailing `status:` line; it is real, already works (`parley roster show --explain
zcode-1` prints the `model unknown / SET BY ~/.parley/agents.toml` row), and adding the agent-side
`~/.zcode/cli/config.json → model.main` line labelled READ, NOT PASSED fits cleanly. It ships
with this change.

CHANGE on D1 (app-server): I agree with @codex-1 — it is the only known binding mechanism but it is
ZCode Protocol (`session/create` carries `{providerId,modelId,variant?}` + thought level), not ACP
(`Spec.ACPArgs` at `discover.go:69-73`), requires a dedicated runner/client, and needs a decision
about how the roster proves a non-argv binding. Defer to a named successor (`idea: zcode-app-server`).
It must NOT be smuggled into this additive adapter.

CHANGE on D4 (test files): I listed six candidates in R1 (autonomous_test, modelmeta_test,
acp_specs_test, naming_test, app_test, preflight_test) but did NOT verify which enumerate adapters.
PRIMARY verification this round (see below): autonomous_test.go enumerates autonomous modes
(:11-49) and full-contract cases (:60-67) — ADD zcode; modelmeta_test.go enumerates producers
(:20) — ADD `zai/glm-5.3`; acp_specs_test.go is required-subset assertions, zcode not ACP, NO
CHANGE; naming_test.go is spec-generic, optional case only; app_test pins fake-codex behavior,
NOT adapter enumeration; preflight_test pins protocol freshness. Confirmed: only the first two
need the adapter assertion; the others are verified unchanged.

CHANGE on D4 (agents.toml clean-up): delete the redundant wholesale `[agents.zcode]` override in
`~/.parley/agents.toml` (same mechanism that stripped hermes's `--yolo` once) and drop `model`/`effort`
from `[roster.zcode-1]` to prevent the display/column divergence @claude-1 flagged (display name
would read `zcode_glm-5.3_max` while MODEL reports `unknown`). Confirmed with `ls` and `cat`.

CHANGE on D3: the exit-0 hazard does not reproduce — no zcode-specific post-launch parser needed.
The generic artifact-validation path (`runner.go:627-714`) already decides correctness (artifact
missing → `agent.failed`, regardless of exit code). The real exit-0 failure from this very round
was hermes-1 (model/harness error, not a flag), caught by the artifact check, not by exit code.
Keep artifact validation; add the `verify --full --yes` gate from Q4; no extra machinery.

## Responses to others

### @claude-1
Right on Q1 (`unknown` / model-unbound; `--explain` instead of a column lie) and Q2 (Scope empty
per `discover.go:86-92`). I adopt your `--explain` proposal fully: `parley roster show --explain
zcode-1` will name `~/.zcode/cli/config.json → model.main` as READ, NOT PASSED; it ships. On D2 I
agree it ships here, not separately. On Q4 I concede: the exit-0 claim doesn't reproduce on
0.16.3, the generic artifact/`verify --full --yes` checks are sufficient, and a zcode-specific stdout
heuristic is brittle. On the display-name divergence (roster-entry `model` vs column `unknown`) I
will delete `model`/`effort` from `[roster.zcode-1]` in `~/.parley/agents.toml` in the same change
so display falls back to a `cli-default` form — explicit, not hidden.

One correction to your list: `naming_test.go` is spec-generic — no adapter enumeration, so no
required change there. `acp_specs_test.go` is subset/no-dup — zcode is not ACP, so no change.
`app_test.go` and `preflight_test.go` are also not adapter-enumerating; verified by grep. Only
`autonomous_test.go` and `modelmeta_test.go` must be extended.

### @codex-1
Conceded on app-server (D1): it is the only binding mechanism but needs a dedicated runner/client,
needs the roster to prove non-argv binding, and is NOT this adapter's scope. I propose the named
successor (`idea: zcode-app-server`) rather than pulling it in. Right about `--allowed-tools`/
`--disallowed-tools`: they filter *tools*, not filesystem — not a sandbox — so Scope stays `""`. The
release evidence (Q6) is adopted as stated: CLI 1.45.0 + skill 2.9.0, per-channel independent proof,
live full-probe with the real 3.7.7-13 binary, clean-profile smoke, and the post-cleanup AUTO check.
Confirmed `version.go` is 1.44.0 and `VERSION` matches; minor bump is correct.

On Q1: agreed, MODEL/EFFORT stay `unknown`; adding `"zai": "Zhipu AI"` to `modelmeta.go`
does NOT change today's `metadata-unknown` (derived from `row.Model` in `roster.go:386-390`), but
it prepares the registry for any future effective `zai/` reference (future schema or a model flag).
That is exactly the reason, not overclaim.

### @hermes-1
Conceded your exit-1 measurement (my exit-0 was a `head`-pipe artifact); corrected the agents.toml
note. Adopted your `--explain` and scope-empty arguments; adopted your `--mode yolo` explicit
(in case config returns to `build`). On Q5 your file list is largely right; my verification narrows
it (see Position changes). Adopted your `naming_test.go` observation (optional only). On D6, adopted
your `fireworks/inkling` record — the roster records the same `glm-5.3` model previously failing as
"no answer" (fixed by write-first prompting); 1 failure / 1 success this round, intermittent, not
categorical. This is a roster/prompt-shape signal, not a zcode design decision; I record it in D6.
Your secret-redaction discipline adopted (no key content in any artifact; `~/.hermes/config.yaml` and
`~/.zcode/cli/config.json` referenced by path only, with key values redacted).

### @kimi-1
Confirmed your measurements PRIMARY: `--model` not in help at all; `--settings` and `--max-turns`
listed but rejected by parser (all exit 1); no env knob across the 13 MB runtime bundle; no model
id in `--json`; `permission.mode=build`; base `headless_args` must stay minimal. Your locator
`roster.go:188-192` is the binding contract against putting `zai/glm-5.3` in the MODEL cell.
Conceded Q3 (kimi-1 AUTO=yes) — verified by `roster show`. Confirmed the `zai` (no hyphen) prefix
divergence explains `metadata-unknown`; adding to `modelmeta.go` is correct. Confirmed your
app-server probe: `session/create` carries `{providerId,modelId,variant?}` + thought level; it IS
ZCode Protocol, not ACP. Confirmed clean-profile smoke and the release-channel checklist.

## New concerns / questions

- `--allowed-tools` and `--disallowed-tools`: listed in `--help` but I did NOT verify their runtime
  behavior (they may also be listed-but-rejected, given `--settings`/`--max-turns`). Before anyone
  proposes adding them to `HeadlessArgs`, a probe is required — and it must be run with `echo $?`
  directly (not through a pipe) to avoid the exit-0 trap.
- `parley roster show --explain zcode-1` works today but names `~/.../agents.toml` as the source
  for the unbound fields; once the adapter ships, `--explain` should name `discover.go` spec +
  `~/.zcode/cli/config.json` for MODEL provenance. Confirm the explanation renderer picks up
  the new source after the adapter is registered.
- `zcode --json` does NOT include `modelId`; if a future version adds it, the adapter's Notes
  must be updated but MODEL stays `unknown` (the cell is about argv binding, not telemetry).
- The release matrix (D5/D6) needs the actual live binary on this machine; it is present
  (`/opt/homebrew/bin/zcode`, 3.7.7-13, 0.16.3); the full-verify integration case from Q4
  must execute against it before calling the adapter shipped.

## Current proposal

**D1 — Defer `zcode app-server`.** Confirmed: `zcode --help` names `app-server` as "ZCode Protocol
stdio app server"; the `zr` method table in the installed `vendor/zcode.cjs` has `session/create`,
`session/setModel`, `session/setThoughtLevel`; the `session/create` schema (`$l`) carries
`{providerId,modelId,variant?}` plus thought level. The transport's requests/lifecycle are
ZCode-specific, not JSON-RPC 2.0 ACP (`discover.go:69-73`). It is NOT `Spec.ACPArgs`. Ship this
adapter (one-shot argv) and open `idea: zcode-app-server` as the successor: derive a protocol spec,
test session lifecycle, decide how the roster proves non-argv binding, then adopt. I would sign.

**D2 — `--explain` augmentation SHIPS here.** `parley roster show --explain AGENT` is verified
live in the CLI (`internal/app/roster.go:66`, `roster_view.go`). Proposal: for zcode-1's MODEL
field, `--explain` reports `~/.zcode/cli/config.json → model.main` (current value, e.g. `zai/glm-5.3`)
labelled **READ, NOT PASSED**; the EFFECTIVE column stays `unknown`; STATUS stays `model-unbound`;
no new vocabulary. I would sign this; it answers the operator's real question ("where does the model
come from?") without violating the contract.

**D3 — Q3 premise dead; no adapter-level change for launch-failure detection.** Confirmed: rejected
flags exit 1 (three shapes, measured directly, not through a pipe); the generic artifact-validation
path (`runner.go:627-714`) catches both exit-0-empty (the hermes-1 failure this round: model error,
40 bytes, exit 0, no artifact) and exit-1-empty failures; `agents verify --full --yes` (
`internal/app/app.go:2108-2157`) is the adapter-qualification gate (sentinel file + exact argv).
I sign: no zcode-specific post-launch parser; rely on the existing artifact/verify gates. The exit-0
hazard from 00-prompt does not reproduce; correct the agents.toml note.

**D4 — Confirmed file list.** PRIMARY verification (grep + read):
- `discover.go`: ADD zcode `Spec` (`LaunchHeadless`, `HeadlessArgs` exact token order, `Scope: ""`,
  `AutonomousWrite{Mode: "yolo", Args: ["--mode","yolo"], Scope: ""}`, `Model: CLIDefault`, `Notes`
  with no-model-flag + exit-1 facts + 0.16.3 pin).
- `modelmeta.go`: ADD `"zai": "Zhipu AI"` (namespace form of the existing `"z-ai"` entry at :43).
- `autonomous_test.go`: ADD zcode (`wantMode` + full-contract case) — VERIFIED enumerates.
- `modelmeta_test.go`: ADD `zai/glm-5.3`, `zai/glm-5-turbo` — VERIFIED enumerates producers.
- `acp_specs_test.go`: NO CHANGE (subset/no-dup, zcode not ACP) — VERIFIED by reading.
- `naming_test.go`: OPTIONAL composite case (`zcode_glm-5.3_max`) — NOT required; spec-generic.
- `app_test.go`, `preflight_test.go`: NO CHANGE — NOT adapter-enumerating (verified).
- Skill `SKILL.md`: ADD `zcode` autonomous-write row.
- `docs/cli-reference.md`: ADD verify example + MODEL/EFFORT `unknown` is expected, not a fault.
- `README.md`: no structural change needed (`:33` already covers "and more"); optional mention.
- `CHANGELOG.md`: v1.45.0 entry.
- `VERSION` / `version.go`: 1.44.0 → 1.45.0.
- `~/.parley/agents.toml`: DELETE redundant `[agents.zcode]` block (wholesale override landmine);
  DELETE `model`/`effort` from `[roster.zcode-1]` (display divergence); correct stale AUTO note.

**D5 — Release: CLI 1.45.0 + skill 2.9.0.** Confirmed: `VERSION` and `version.go` both 1.44.0;
precedent (v1.39.0 for kimi/opencode promotion) is minor. Skill `package.json` is 2.8.0 → 2.9.0.
Per `RELEASING.md` (read): npm (`npm test`, `npm pack --dry-run`, `npm publish`), GitHub (`tag`
`v2.9.0`, portable assets from release, not local builds), WinGet (manifest hashes from GitHub
release assets, NOT local builds — `winget validate`), Homebrew (formula `url`/`sha256` from release
tarball, `brew audit --strict --online`). For CLI (README docs `brew install` but no repo-local
`.github/release*` automation): tag + binary upload + checksum publication must be written/proven.
New-adapter-specific gates: (a) `parley agents list` shows zcode headless argv; (b) `parley agents verify --full --agent zcode --yes` passes with real 3.7.7-13 binary; (c) clean-profile smoke
(no old `agents.toml` override) — AUTO stays yes; (d) `--explain zcode-1` reports READ source;
(e) artifact-check success in a scratch round. I sign.

**D6 — hermes-1 (`fireworks/inkling`) record.** Confirmed from round-01 file: 1 failure (`Model
generated invalid tool call: bash`, 40 bytes, exit 0, no artifact — caught by artifact check) + 1
success this round (same prompt, same flags); a short tool-call probe passed on inkling and on the
old `glm-5p2`; the roster records the *same `glm-5.3` model* previously failing under opencode as
"no answer, not wrong answer" (fixed by write-first prompting). This is a roster/prompt-shape signal,
intermittent, NOT a zcode adapter decision — but it is this round's evidence and belongs in the
artifact (see this section). I do NOT propose changing the adapter for it; it stays recorded.

## Round-1 files, in full

See `round-01/claude-1.md`, `codex-1.md`, `hermes-1.md`, `kimi-1.md`. All four independently
reached MODEL/EFFORT `unknown`; @claude-1 proposed `--explain`; @codex-1 scoped app-server as
successor; @kimi-1 confirmed exit-1, AUTO=yes, and the 13 MB bundle probes; @hermes-1 recorded the
exit-0 confusion and the hermes-1 failure.

The complete proposal I sign: (1) adapter spec in `discover.go` with exactly the verified argv, Scope
`""`, MODEL/EFFORT `CLIDefault`; (2) `modelmeta.go` adds `"zai"`; (3) `autonomous_test.go` and
`modelmeta_test.go` extended; others unchanged; (4) `SKILL.md` updated; (5) `docs/cli-reference.md`
updated; (6) `CHANGELOG.md` / `VERSION` / version = 1.45.0; (7) skill 2.9.0; (8) `--explain` reports
agent-side config source, READ NOT PASSED; (9) `~/.parley/agents.toml` override deleted; (10) no
protocol change (COOPERATION.md untouched); (11) app-server deferred to named successor; (12) no
new STATUS vocabulary. Confirmed PRIMARY measurements (zcode binary, `--version`, `--model` rejected,
`--settings`/`--max-turns` exit 1, `--json` lacks modelId, `discover.go` rules read, `autonomous_test.go`
and `modelmeta_test.go` enumerating, others not, `agents.toml` contents). No fabricated results.
