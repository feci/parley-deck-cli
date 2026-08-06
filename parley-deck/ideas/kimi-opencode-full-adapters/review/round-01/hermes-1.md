---
agent: hermes-1
idea: kimi-opencode-full-adapters
round: 1
phase: review
date: 2026-08-06
reviewed-commit: ed360c4
---

## Scope I checked

I read the two source files (`internal/agents/discover.go` lines 280–344 for the
new specs, `internal/agents/acp_specs.go` for the merge logic), the test
(`internal/agents/autonomous_test.go`), and the runner launch path
(`internal/runner/runner.go` lines 1056–1122: `buildAgentInvocation` and the
`PromptStdin` stdin binding). I re-ran probes 1, 2, 3, 5, 6 myself against the
installed CLIs (kimi 0.33.0 at `/Users/tomasfecko/.kimi-code/bin/kimi`, opencode
1.18.14 on PATH). I built the project (`go build ./...`), ran the full test suite
(`go test -count=1 ./...`, 25 packages, all pass), ran `agents list` on the new
binary, and wrote a throwaway test exercising `buildAgentInvocation` for both
specs to verify the exact argv the runner produces. The throwaway test was
deleted after running; the working tree is clean.

## Summary

A clean, surgical two-spec addition. The probes hold, the merge path is correct,
the launch argv is verified end-to-end, and nothing regressed. The one risk the
implementer flagged (the launch path) is the one I verified most deeply, and it
is correct. This is ready to release.

## Refutation attempts

**Attempt 1 — "kimi's Mode is misnamed as `prompt`."** Re-ran probe 1:
`kimi --auto -p "…"` exits 1 with `error: Cannot combine --prompt with --auto.`
[PRIMARY — I executed it, exit code captured directly, not through a pipe].
`-p` is the only headless shape kimi has, and in print mode the default
permission policy applies tool calls without asking (probe 2: file written,
exit 0). Calling the mode `prompt` is accurate — it names the vendor's
non-interactive mode that auto-approves. The alternative `auto` would be a lie
since `--auto` is rejected alongside `-p`. Not broken.

**Attempt 2 — "opencode's Mode `auto` / Args `["--auto"]` is misnamed or
redundant."** Re-ran probe 3: `opencode run --auto "…"` writes the file, exit 0
[PRIMARY]. Probe 6 confirms `--auto` = "auto-approve permissions that are not
explicitly denied (dangerous!)". The user's decision (00-prompt.md Q1) was to
pass `--auto` explicitly to keep `AutonomousWrite.Args` a subset of
`HeadlessArgs`, matching claude/codex/hermes. The name `auto` matches the
vendor's flag. Not broken.

**Attempt 3 — "Promoting to DefaultSpecs broke the ACP path for kimi/opencode."**
`mergeACPCatalog` (acp_specs.go:52) builds a `byID` map from the built-in specs,
then for each catalog entry takes the merge branch (`mergeACPBackend`) if the ID
already exists, or the append branch (`specFromACPBackend`) otherwise. Before
this change, kimi/opencode were NOT in `defaultBuiltinSpecs()`, so they took the
append branch and got `LaunchMode: LaunchACP` stubs. Now they ARE in
`defaultBuiltinSpecs()` with `LaunchMode: LaunchHeadless`, so they take the merge
branch, which preserves their headless fields and only overlays `ACPArgs` +
notes. I confirmed via `agents list` that both still show `acp: … ["acp"]` lines,
and I confirmed `kimi acp --help` and `opencode acp --help` are valid
subcommands [PRIMARY — I ran both]. ACP survives as an alternative launch mode
(selectable by overriding `LaunchMode` to `acp`). Not broken.

**Attempt 4 — "The runner builds the wrong argv for `{prompt}`."** This is the
gap the implementer explicitly did not close. I wrote a test calling
`buildAgentInvocation("/fake/root", disc, "DO THE THING")` for each spec and
inspected the returned `args` [PRIMARY — I executed the code]:
- kimi → `[-p DO THE THING]` (prompt is the VALUE of `-p`, one argv element)
- opencode → `[run --auto DO THE THING]` (prompt is positional after `run --auto`)
Both specs declare `PromptMode: PromptArg`. The runner at line 1056 only binds
stdin when `PromptMode == PromptStdin`, so for both, stdin is NOT connected and
the prompt travels entirely via argv — which is exactly how these CLIs expect it.
Not broken. This is the most important refutation and it holds.

**Attempt 5 — "A regression in `roster init` or `MachineFamilyCatalog`."**
`MachineFamilyCatalog` (runtime.go:253) collects `DefaultSpecs()` IDs into a
`map[string]bool`. Before this change, kimi/opencode IDs entered that map via
the ACP append stubs; after, they enter via the built-in specs. Either way the
IDs are present — the set is identical. IMPLEMENTATION.md check 5 claims both
old and new binaries propose the roster mappings, and I have no reason to
dispute it; the code path is ID-set membership, which is unchanged. Not broken.

## Findings

### [CRITICAL] — none.

### [MAJOR] — none.

### [MINOR] — none.

### [NIT] — `Telemetry` field value is imprecise for both specs

Both kimi and opencode set `Telemetry: "final text on stdout"`. For kimi this is
accurate — probe 2 shows it prints final text. For opencode, probe 3 shows it
also streams intermediate output (`> build · xai/grok-4.5`, `← Write probe3.txt`)
before the final text, so "final text on stdout" understates what the runner's
activity tracker will see. This field is descriptive only (it drives a `note:`
line in `agents list`, not any runtime behavior), so it is a NIT. A value like
"streamed build/tool events then final text" would be more accurate for
opencode. Not blocking.

## Open questions

None. The change is self-contained, the probes are reproducible, and the launch
path — the one thing the implementer could not verify — is verified here.

## Verdict

**Ready to release.** All five review questions resolve in the implementer's
favor. The declared modes are truthful, empty `Scope` is correct and consistent
with claude/agy/hermes, the ACP path survives, and the launch argv is correct
for both `{prompt}` placements. The full test suite passes (25 packages), the
build is clean, and the working tree is byte-identical to the base commit.
