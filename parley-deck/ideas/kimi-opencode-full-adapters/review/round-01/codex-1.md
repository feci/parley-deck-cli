---
agent: codex-1
idea: kimi-opencode-full-adapters
round: 1
phase: review
date: 2026-08-06
---

## Summary

I found one MAJOR effective-launch mismatch and one MINOR regression-test gap. The implementation
is not ready to release. Detailed §15 verdicts on the implementer's existing factual claims, and
the evidence for the reviewer-owned findings on which I cannot issue self-verdicts, follow below.

## Scope checked

I reviewed commit `ed360c4c43246e6235081a0edaa4fc439c15ad4d` and the implementation commits it
contains. I read the requested prompt, implementation record, `discover.go`, `acp_specs.go`, and
`autonomous_test.go`; I also traced configuration layering and the exec/ACP runner paths in
`internal/config/runtime.go`, `internal/runner/runner.go`, and `internal/runner/acp.go`, plus their
relevant tests. I compared the repository against an archive made with
`git archive ed360c4 | tar -x -C <temp>` and ran checks only in that extracted tree.

Runtime scope was Kimi Code 0.33.0 at `/Users/tomasfecko/.kimi-code/bin/kimi` and OpenCode
1.18.14 at `/opt/homebrew/bin/opencode`. I ran help/version probes, both incompatible Kimi flag
shapes, both OpenCode headless shapes, ACP initialization probes, `agents list` with and without
the machine configuration, and the Go suite. The positive Kimi write probe was prevented by this
reviewer's macOS sandbox with `EMFILE` from `FSWatcher`; the two positive OpenCode write probes
reached the configured model but produced no tool result before I stopped the disposable probes.
Accordingly, I rely on the implementer's named primary evidence for those three write outcomes
rather than claiming an independent reproduction. Kimi ACP initialization completed over stdio;
OpenCode ACP initialization returned `ServeError` under this sandbox, so I checked its structural
preservation and authoritative CLI contract but not a successful local OpenCode ACP session.

Checks in the archived tree:

```text
$ go test -count=1 ./internal/agents/...
ok  parley-deck-cli/internal/agents

$ go test -count=1 -skip TestDurableKillEndToEndRealProcess ./...
all listed packages: ok
```

The unskipped suite produced one failure in
`TestDurableKillEndToEndRealProcess` because process verification could not obtain a boot ID; it
was excluded from the adapter conclusions in this review.

## Refutation attempts and verdicts on existing claims

### Declared modes

- **[PRIMARY] CONFIRMED** — Kimi's vendor name `prompt` is not a fabricated permission-mode name.
  `/Users/tomasfecko/.kimi-code/bin/kimi --help` describes `-p/--prompt` as running one prompt
  non-interactively and calls its output setting "prompt mode". The live commands
  `kimi --auto -p <prompt>` and `kimi -y -p <prompt>` both exited 1 with, respectively,
  `Cannot combine --prompt with --auto` and `Cannot combine --prompt with --yolo`. Thus
  `Mode: "prompt", Args: ["-p"]` names the only compatible headless mode rather than pretending
  that `auto` or `yolo` is enabled.
- **[SECONDARY] CONFIRMED** — the pre-existing claim that plain `kimi -p` writes unattended relies
  on claude-1's non-RECALL `PRIMARY` verdict and quoted command/result in
  `IMPLEMENTATION.md:29-43` (`exit 0`, file written). My independent positive probe was
  environment-blocked as declared above.
- **[PRIMARY] CONFIRMED** — OpenCode's `Mode: "auto", Args: ["--auto"]` exactly matches the
  vendor flag. `opencode run --help` says `--auto` auto-approves permissions that are not
  explicitly denied. This is the mode the user explicitly selected in `00-prompt.md:18-22`.
- **[SECONDARY] CONFIRMED** — the pre-existing claims that OpenCode writes with and without
  `--auto` rely on claude-1's `PRIMARY` verdicts in `IMPLEMENTATION.md:31-46`; my two positive
  probes did not finish within the bounded review window.

### Empty scope and `agents list`

- **[PRIMARY] CONFIRMED** — empty `Scope` is correct for both built-ins.
  `internal/agents/discover.go:61-66,86-108` reserves `workspace` for an enforced sandbox and
  separates `Declared()` from `Confined()`. Neither live CLI help describes these permission
  modes as a filesystem sandbox. With built-ins isolated from local overrides, `agents list`
  displayed `SANDBOX=cli-default` next to `AUTO=yes` for each, so the empty scope itself is not
  hidden or falsely promoted to confinement.

### ACP preservation

- **[PRIMARY] CONFIRMED** — the merge branch preserves the alternative launch configuration.
  `internal/agents/acp_specs.go:52-89` copies non-nil catalog `ACPArgs` into an existing spec;
  the built binary printed `acp: kimi ["acp"]` and `acp: opencode ["acp"]`. A direct Kimi
  `initialize` request to `kimi acp` returned JSON-RPC 2.0 with protocol version 1 and agent info
  for Kimi Code CLI 0.33.0. OpenCode's official ACP documentation states that `opencode acp`
  is an ACP subprocess using JSON-RPC over stdio:
  `https://dev.opencode.ai/docs/acp/`.
- **[PRIMARY] UNVERIFIED** — a successful local OpenCode ACP handshake in this reviewer sandbox.
  The exact stdio initialize check returned `ServeError`; the result does not distinguish an
  adapter failure from the sandbox's filesystem/socket restrictions. This does not contradict
  the source-level merge result above, but it limits the runtime scope of this review.

### Prompt placement in the runner

- **[PRIMARY] CONFIRMED** — the built-in exec path constructs both commands correctly.
  `internal/runner/runner.go:1094-1107` replaces an argument equal to `{prompt}` with one argv
  element, and `runner.go:1056-1058` attaches stdin only for `PromptStdin`. Applied to
  `discover.go:292-296,324-328`, the resulting argument vectors are exactly
  `[-p, <prompt>]` for Kimi and `[run, --auto, <prompt>]` for OpenCode. The prompt is therefore
  the value of Kimi's `-p` and the positional message after OpenCode's `run --auto`.

## Findings

### [CRITICAL] None

### [MAJOR] The effective OpenCode launch silently drops `--auto`

This is a reviewer-owned claim, so §15 forbids me from issuing a verification verdict on it. The
reviewable evidence is:

```text
~/.parley/agents.toml:
headless_args = ["run", "-m", "litellm/xai/grok-4.5", "{prompt}"]

internal/config/runtime.go:542-544:
if len(override.HeadlessArgs) > 0 {
    spec.HeadlessArgs = expandSlice(override.HeadlessArgs, root, tempdir)
}

internal/runner/runner.go:1098-1107:
for _, arg := range agent.HeadlessArgs { ... }
```

The central override wins over `DefaultSpecs`, so a real round in this repository executes
`opencode run -m litellm/xai/grok-4.5 <prompt>`, not the user-decided
`opencode run --auto <prompt>`. At the same time, `AutonomousWrite` remains the built-in
`Mode: "auto", Args: ["--auto"]`. The current `agents list` compounds this by printing the
built-in `HeadlessMode` summary `opencode run --auto` and `AUTO=yes`, even though it does not print
or validate the effective `HeadlessArgs` that the runner uses.

This matters even though the implementer's no-flag probe wrote successfully: the user explicitly
chose `--auto`, and the declared mode's enabling argument is absent from the launched command.
Before release, remove the stale local `headless_args` override or add `--auto` to it, then run one
actual Parley round with OpenCode. Product hardening should also derive the displayed headless
command from effective arguments and fail closed (or at least warn and report `AUTO=no`) when an
effective configuration drops any declared `AutonomousWrite.Args`.

### [MINOR] The autonomous contract test checks only `Mode`

This is a reviewer-owned claim and therefore carries no self-verdict. In
`internal/agents/autonomous_test.go:16-38`, the new cases assert only that the spec exists,
`Declared()` is true, and `Mode` matches. The test would still pass if `Args`, `HeadlessArgs`,
`PromptMode`, prompt placement, `Scope`, or merged `ACPArgs` regressed—the fields that make these
adapters usable and truthful.

Extend the table-driven contract to assert, for both IDs, the exact autonomous args, exact
headless argv template, `PromptArg`, empty scope, and `ACPArgs == ["acp"]`. Add a layered-config
case showing what happens when an override removes a declared autonomous argument; it should not
continue to report an enabled mode silently.

### [NIT] None

## Open questions

None. The OpenCode ACP sandbox limitation is declared scope, not a blocker to the source-level ACP
merge finding. The effective OpenCode headless mismatch is the release blocker.

## Release verdict

**Not ready to release.** Resolve the MAJOR effective-launch mismatch and re-run a full
Parley-launched OpenCode round. The MINOR contract coverage should be added in the same fix because
it directly guards the fields changed by this idea.
