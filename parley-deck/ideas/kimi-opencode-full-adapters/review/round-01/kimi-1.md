---
agent: kimi-1
idea: kimi-opencode-full-adapters
round: 1
phase: review
date: 2026-08-06
reviewed-commit: ed360c4
---

## Summary

Both new specs are truthful, and the one gap the implementer admitted — "the probes do not prove
the runner drives a full round with them" — is now closed one level below a full round: I drove
both adapters through parley's own runner (`agents verify --full`) against the pure built-in
specs, and both passed with validated sentinel files. Every probe row in `IMPLEMENTATION.md`
reproduced under my hands, including the roster-init correction on **both** binaries. No
CRITICAL, no MAJOR. Two MINOR (a missing merge-outcome test; kimi is undiscoverable at its
default install prefix without a `command` override), two NIT. **Ready to release**; MINOR-1 is
a ten-line test worth taking before tagging.

Ownership: I filed no earlier artifact in this idea and own no claim in `IMPLEMENTATION.md` or
the diff, so no self-verdict prohibition is engaged. Tag convention: verdicts on factual claims
carry `PRIMARY`/`SECONDARY`/`RECALL` per §15.2; severity assignments and the release ruling are
positions about what should be and carry no tag (§15.1, last line).

## Scope checked (declared per §15)

- Read in full: `00-prompt.md`; `IMPLEMENTATION.md`; the complete diff `d93e313` (code) and
  `ed360c4` (doc) via `git show`; `internal/agents/discover.go`, `acp_specs.go`,
  `autonomous_test.go`, `acp_specs_test.go`; `internal/runner/runner.go:420-660, 1014-1170`;
  `internal/runner/supervision.go` in full; `internal/app/roster.go:78-397`;
  `internal/app/app.go:420-494, 2067-2140`; `internal/app/preflight.go:690-750`;
  `internal/app/consensus_request_signoffs.go:340-460`; `parley-deck/COOPERATION.md` §15
  (`:1176-1316`) and §2 (`:101-134`); `parley-deck/agents.toml`; `~/.parley/agents.toml`
  (machine-local, read-only).
- Executed (all in `/tmp/kimi1-review/`; the review tree was extracted with
  `git archive ed360c4 | tar -x`; **no git write commands were used** and the only working-tree
  write is this file):
  - The four vendor probes, re-run from scratch with byte-level timing (details under §1 below).
  - `kimi --help`, `opencode run --help`, `opencode models`, `kimi --version` (0.33.0),
    `opencode --version` (1.18.14); `command -v kimi` in a plain and a login shell.
  - `go build ./...`, `go test -count=1 ./...` (exit 0, all packages),
    `go test -count=1 ./internal/agents/...` (ok) on the archived tree; binary built as
    `/tmp/kimi1-review/parley`.
  - `parley agents list` against pure built-ins (`PARLEY_HOME` pointed at an empty dir, kimi's
    bin dir prepended to PATH), plus a deck override flipping `launch_mode = "acp"`.
  - `parley agents verify --agent <kimi|opencode> --full --yes` against scratch decks — the
    runner's own headless probe (`runHeadlessProbe`, `internal/app/app.go:2088`), which uses
    `runner.CommandFor` → `buildAgentInvocation`, the same argv construction the round path uses
    (`internal/runner/runner.go:1024, 1097-1108`).
  - `parley roster init --dry-run` on a disposable deck seeded with `kimi-1` + `opencode-1`, on
    the new build **and** the installed `parley 1.38.0`.
- Not checked (visible, not inferred from silence): a full multi-phase protocol round driven by
  parley with kimi/opencode (frontmatter-validated round artifacts, retry semantics, signoff
  flow) — still untested by anyone; long-round stall behavior (my probes are short);
  Windows (`dist/` ships windows-x64/arm64 binaries; I reviewed and probed on darwin/arm64
  only); an actual ACP session (`kimi acp` speaking JSON-RPC) — I verified its
  *configurability*, not the protocol exchange; `agents.local.toml` layering (none exists here).

## 1. Are the declared modes truthful? — Yes, both. `CONFIRMED`, `PRIMARY`

Re-ran all four probes on 2026-08-06 (kimi 0.33.0 at `/Users/tomasfecko/.kimi-code/bin/kimi`;
opencode 1.18.14 at `/opt/homebrew/bin/opencode`):

| Probe (my command, scratch cwd) | My result |
|---|---|
| `kimi --auto -p "Create a file named should-not-exist.txt …"` | `exit=1`, stderr `error: Cannot combine --prompt with --auto.`, file **not** written |
| `kimi -p "Create a file named probe-kimi.txt …"` | `exit=0`, `probe-kimi.txt` = `kimi-probe-ok` |
| `opencode run --auto "Create a file named probe-opencode.txt …"` | `exit=0`, file = `opencode-probe-ok` |
| `opencode run "… probe-opencode-plain.txt …"` (no `--auto`) | `exit=0`, file = `plain-ok` |

Help text also reproduced (`PRIMARY`): kimi `-p, --prompt <prompt>` = "Run one prompt
non-interactively and print the response", `--auto` = "fully autonomous", `-y, --yolo` =
"Auto-approve regular tool calls"; opencode `--auto` = "auto-approve permissions that are not
explicitly denied (dangerous!)", `--variant` = reasoning effort.

Neither mode is misnamed. `Mode "prompt"` is kimi's long flag name (`--prompt`) and names the
only autonomous headless shape the CLI has — the combination that would be more explicit
(`--auto -p`) is rejected by the vendor, probe 1 above. `Mode "auto"` is opencode's flag name
verbatim. One honest asymmetry, already documented in the spec comment itself: kimi's `-p`
auto-approval is the *print-mode default policy*, not an explicit grant — which is exactly the
"implicit default a vendor may change between versions" risk the user accepted `--auto` to avoid
for opencode. For kimi there is no combinable explicit flag, so the risk is unavoidable; the
spec's `ApprovalPolicy: "print-mode default"` states it accurately rather than dressing it up,
and if a vendor tighten ever lands, a parley round fails loud (stall/hard-timeout under
supervision) rather than silently. Observation, no severity.

## 2. Empty `Scope` + `AUTO=yes` — correct, and not misleading. `CONFIRMED`, `PRIMARY`

Both probes show permission grants, not sandboxes: nothing in `kimi -p` or `opencode run --auto`
confines writes at the OS level — the probe prompts could have written anywhere on disk. Empty
`Scope` is the type's own honesty rule (`discover.go:86-92`), correctly applied.

The reader-of-`agents list` concern does not survive contact with the actual output (my run,
pure built-ins):

```
kimi     yes  0.33.0   headless  configured  cli-default  print-mode default  cli-default  1800000  no  yes  hosted
opencode yes  1.18.14  headless  configured  cli-default  auto                cli-default  1800000  no  yes  hosted
```

`AUTO=yes` sits next to `SANDBOX=cli-default`; codex's row shows `workspace-write` in the same
column. The matrix distinguishes confined from unconfined autonomy exactly as it already does
for claude (`bypassPermissions`, `Scope ""`) and hermes (`yolo`, `Scope ""`). What is true is
narrower: `AutonomousWrite.Scope`/`Confined()` is surfaced *nowhere* — the SANDBOX column reads
a different field (`SandboxMode`) that happens to carry the same signal today. That is a display
gap, not a deception → NIT-1.

## 3. Does `DefaultSpecs` promotion break the ACP path? — No. `CONFIRMED`, `PRIMARY`

For both IDs `mergeACPCatalog` now takes the merge branch (`acp_specs.go:60-61`); the merge
result keeps headless as default and gains ACP as an opt-in:

- `agents list` shows `acp: kimi ["acp"]` and `acp: opencode ["acp"]` — `ACPArgs` survived the
  merge onto the full specs (`mergeACPBackend`, `acp_specs.go:77-79`). No duplicate specs (the
  uniqueness sweep in `TestDefaultSpecsMergesACPCatalog` passes).
- A deck override `[agents.kimi] launch_mode = "acp"` flips the row to `LAUNCH=acp` (verified on
  a scratch deck) — the alternative launch mode is one config line away, satisfying the
  00-prompt constraint.
- Regression sweep: `go test -count=1 ./...` exit 0 across all packages on the archived tree;
  the 11 ACP-only catalog agents still append via `specFromACPBackend`
  (`TestACPSpecsAreMarkedLaunchACP` green).
- The kimi spec's own Notes claim ("ACP remains available as an alternative launch mode via
  `kimi acp`") is accurate and not duplicated by the merge (`joinNotes` dedupes; the catalog
  entries for kimi/opencode carry no Notes of their own).

## 4. The launch path — the admitted gap, closed at probe level. `CONFIRMED`, `PRIMARY`

Code reading first: `buildAgentInvocation` (`runner.go:1097-1108`) substitutes `{prompt}` as a
single argv element with no shell involved; `PromptArg` means stdin is never wired
(`runner.go:1056-1058`). For kimi that yields argv `["-p", <prompt>]` — prompt as the flag's
value; for opencode `["run", "--auto", <prompt>]` — prompt as the positional. Both shapes are
exactly what my probes 2 and 3 executed by hand.

Then the part the implementer did not do — through parley's own runner, against the pure
built-in specs (isolated `PARLEY_HOME`, so no user overrides; kimi's bin dir on PATH):

```
$ parley agents verify --agent kimi --full --yes --dir /tmp/kimi1-review/deck
kimi: installed version=0.33.0
kimi: headless probe passed          (rc=0)
$ parley agents verify --agent opencode --full --yes --dir /tmp/kimi1-review/deck2
opencode: installed version=1.18.14
opencode: headless probe passed      (rc=0)
```

Both sentinel files carry the exact required first line (`# parley-runtime-probe agent=kimi
run=…` / `agent=opencode run=…`; `runHeadlessProbe` validates it). This exercises
`runner.CommandFor` → `buildAgentInvocation` — the same substitution code the round path's
`execAgentProcess` calls (`runner.go:1024`) — plus process spawn, cwd, and output capture.

Supervision interaction, which no one had checked at all: the first-event guard kills an agent
that emits **zero** stdout/stderr bytes within 120 s (`supervision.go:159`, default
`defaultFirstEventTimeoutMS`). My byte-sampled probes (5 s sampling, streams redirected to files
exactly as the runner does): kimi's first bytes arrive at t≤5 s (a `kimi version 0.33.0` banner
on stderr, followed by intermittent thinking output), opencode's at t≤5 s (progress on both
streams). The guard will not false-fire on either CLI, and `BuffersStdout: false` is correct for
both (contrast agy, which genuinely buffers until exit).

Honest residue, `UNVERIFIED` by anyone (mine, `RECALL`-free by construction — it is an absence
of evidence, not a claim): a full multi-phase protocol round — frontmatter-validated artifacts
across phases, retry-once-on-`no_first_output` semantics, the signoff flow. The marginal risk is
low: artifact validation is agent-agnostic, and the argv/stdio path is now probe-verified. A
further pre-existing property, not a defect of this change: round prompts travel as one argv
element, so OS per-arg limits (~128 KiB on Linux) bound prompt size for all `PromptArg` adapters
(agy, hermes, and now these two) alike.

## 5. Anything else

- `IMPLEMENTATION.md` probe rows 1–8: all reproduced `PRIMARY` (probes 1–4 above; help texts
  above; `opencode models` lists `litellm/xai/grok-4.5` and `litellm/glm-5p2`;
  `~/.kimi-code/config.toml` holds `default_model = "kimi-code/k3"` and `[thinking]
  effort = "max"`). The corrections table is also borne out: `[models."kimi-code/k3"]
  default_effort = "high"` is in the config file verbatim (`PRIMARY`), and the roster-init
  "fail-close" warning is `NOT REPRODUCED` under my hands on **both** binaries — new build and
  installed `parley 1.38.0` each propose `[roster.kimi-1] adapter = "kimi"` and
  `[roster.opencode-1] adapter = "opencode"`, drop nothing, exit 0 on a disposable deck.
- The invariant cited in the user's decision (`AutonomousWrite.Args ⊆ HeadlessArgs`) holds for
  both built-ins. Note it is nominal under overrides: this machine's central
  `~/.parley/agents.toml` pins opencode `headless_args` *without* `--auto` (and hermes without
  `--yolo`), so the declared mode describes the built-in shape, not the effective config.
  Pre-existing pattern, out of scope — observation, no severity.

## Findings

CRITICAL: none. MAJOR: none.

### [MINOR-1] No test locks the ACP-merge outcome for the two promoted IDs

`TestDefaultSpecsMergesACPCatalog` (`acp_specs_test.go:59-92`) asserts post-merge `ACPArgs` for
claude and hermes and `nil` for codex — but nothing asserts that kimi/opencode keep
`ACPArgs == ["acp"]`, keep `LaunchMode == headless`, or keep their `HeadlessArgs` after
`mergeACPBackend`. `autonomous_test.go` covers only `AutonomousWrite.Mode`. A future edit to
`mergeACPBackend` or the catalog could silently drop `kimi acp` as a selectable mode — the very
constraint 00-prompt names — with every test green. Fix: extend the existing merge test with the
two IDs (ACPArgs, LaunchMode, one HeadlessArgs element). Facts `PRIMARY` (files read in full);
the gap's consequence is my analysis.

### [MINOR-2] kimi is undiscoverable at its default install prefix without a `command` override

`command -v kimi` returns rc=1 in both my shell and a login shell (`PRIMARY`); the binary lives
only at `~/.kimi-code/bin/kimi` (the official installer's prefix). The built-in
`Commands: ["kimi"]` resolves via `exec.LookPath` → PATH only, so a fresh user who installed
kimi normally gets `INSTALLED=no` and the new adapter never engages — on this machine the
central `~/.parley/agents.toml` `command` field is the only thing that finds it (the spec's
Notes don't mention this). Every other built-in family here is likewise reached through a
central `command` override, so this is a consistent pattern, not a regression — but for a
*newly promoted* adapter it means the feature does not work out of the box for the default
install. Fix direction (either suffices): probe the well-known prefix at discovery time, or say
"set `command` if kimi is not on PATH" in the spec Notes / release notes.

### [NIT-1] `Scope`/`Confined()` is surfaced nowhere

The honesty signal this idea carefully preserves is invisible in `agents list`; the SANDBOX
column carries it only because `SandboxMode` happens to align. One extra token in the matrix (or
a `confined: yes/no` detail line under `acp:`) would make the declaration auditable from the
CLI. `PRIMARY` — `PrintRuntimeMatrix` (`discover.go:440-515`) read in full; no `Confined()` call
exists outside tests (`grep`).

### [NIT-2] In-repo deck config comment goes stale on release

`parley-deck/agents.toml:22-27` still instructs: "parley 1.36 treats kimi as an ACP backend
(AUTO=no) — invoke it manually, and do NOT run `parley roster init`", and pins
`headless_args = ["-p", "{prompt}"]` which now duplicates the built-in. On release the comment
becomes false guidance inside this repository (the roster-init half is already false today —
both binaries propose the mapping, §5 above). The machine-local `~/.parley/agents.toml` carries
the same stale warning plus stale version notes; it is out of the repo but the user should know
the warning's premise is dead. Fix: refresh the comment block (and drop the redundant pin) when
this ships.

## Verified clean (explicit null results)

- **No other adapter touched.** The code diff is exactly two appended specs in
  `defaultBuiltinSpecs()` plus five lines in `autonomous_test.go`; `acp_specs.go` is byte-unchanged
  (`git show d93e313` read in full). `PRIMARY`.
- **`TestBuiltinsAreAutonomous` extension is faithful** — mode strings match the specs
  (`"prompt"`/`"auto"`), and the test would fail at the parent commit (the IDs would be missing
  from `DefaultSpecs`' map). `PRIMARY` (diff + code read; the failing-direction was not executed
  — the suite was run only at the review commit).
- **Build and full test suite green** at `ed360c4` in the archived copy; `gofmt` clean per the
  implementer (not re-run; `SECONDARY` on that one line only).
- **`agents list` output matches `IMPLEMENTATION.md` check 4 exactly** (`LAUNCH=headless`,
  `AUTO=yes`, `headless: kimi -p`, `headless: opencode run --auto`). `PRIMARY`, quoted in §2.
- **No docs drift**: `docs/` contains no mention of kimi or opencode (`grep` null result), so no
  documentation contradicts the promotion. `PRIMARY`.
- **Version claims**: kimi 0.33.0, opencode 1.18.14 — both reproduced via `--version`; the
  "stale note" rows in the corrections table are accurate. `PRIMARY`.

## Ready-to-release statement

**Yes — ready to release.** Nothing found rises above MINOR; the single admitted verification
gap is now closed at runner-probe level with the exact built-in argv shapes; ACP survives as an
opt-in; the suite is green. Recommend taking MINOR-1 (ten lines of test) before tagging, and
MINOR-2 as at least a release-note line. The remaining unverified residue — a full parley-driven
protocol round — is best closed by simply running the next real idea with the shipped binary,
which this deck is about to do anyway.

## Open questions

1. For claude-1 (implementer): the user's central config deliberately runs opencode *without*
   `--auto` ("narrow shape wins"). After release, should the deck-level recommendation distinguish
   "declared capability" from "recommended effective shape", or is the built-in `--auto` default
   the shape you expect this roster to actually run? This affects whether NIT-2's cleanup should
   also drop the central `headless_args` pin.
2. For any reviewer with a Windows machine: the `dist/` artifacts ship windows binaries — worth a
   smoke of `kimi -p` / `opencode run --auto` argv handling there, or is darwin+linux coverage
   deemed sufficient for this release?
