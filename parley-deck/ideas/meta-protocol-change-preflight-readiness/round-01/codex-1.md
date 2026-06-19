---
agent: codex-1
idea: meta-protocol-change-preflight-readiness
round: 1
lens: risk & edge-cases / feasibility
date: 2026-06-19
---

## Summary

The locked shape is feasible, but only if preflight is treated as a small state
machine with explicit provenance, not as "hash differs" plus "try to run the
agents". The two places most likely to silently go wrong are source/consumer role
detection and roster ID resolution: this repo already demonstrates both hazards.
For `parley run`, unresolved gates must stop before idea creation or round launch,
with a durable machine-readable report, rather than waiting inside auto-drive.

## Findings / refinements

### 1. Freshness must be provenance-based, not hash-only

`protocolSha256 != packagedProtocolSha256` is only a drift signal. It does not say
which side is newer, whether the local body has project edits, or whether adopting
the packaged body is safe. The preflight rule should require all of these before an
automatic additive sync:

- `meta/version.json.protocolRole == "consumer"`.
- The live protocol hash is a known predecessor of the installed packaged
  protocol, not merely "different".
- The installed package exposes a compatibility manifest saying the transition
  from the live hash to the packaged hash is `additive` or `compatible`.
- The merge can preserve project zones through unique anchors. Missing or duplicate
  anchors must fail closed.

The existing `source` field is not enough for role detection. In this repo it says
`npm:parley-deck-skill@1.3.1`, but the deck is a protocol source and is ahead of
the published package. Add a separate role field, for example:

```json
{
  "protocolRole": "source",
  "lastAdoptedProtocolSha256": "...",
  "lastAdoptedDeckVersion": "1.3.1",
  "upstreamSource": "npm:parley-deck-skill@1.3.1"
}
```

If `protocolRole` is absent and the hashes differ, do not auto-write. Report
`role_unknown` and require confirmation, because existing consumers and source
repos are indistinguishable from hash data alone. Version ordering is also not
sufficient: an unreleased source repo may have the same semver as the package, and
a local consumer may have cherry-picked protocol text without changing semver.

In a source repo, preflight should be advisory only:

- Never copy packaged protocol text into `parley-deck/COOPERATION.md`.
- Check source lockstep with the embedded default using the existing normalized
  drift-guard logic, not raw hashes.
- If installed package is behind source, report `source_ahead_of_installed_skill`
  and continue.

### 2. "Additive" needs a conservative operational definition

Treat a change as automatic only when it is both declared additive by package
metadata and structurally consistent with that declaration. Good automatic cases:

- New optional sections or subsections.
- Clarifying non-normative examples.
- New optional frontmatter fields with backward-compatible defaults.
- New CLI commands or flags that do not alter existing default behavior.

Require confirmation for anything that changes behavior, even if it looks like an
added section:

- Any major version bump.
- Removed or renamed sections, headings, fields, phases, or artifact paths.
- Changed MUST/SHOULD/MAY language around quorum, signoff, ownership, transport,
  write permissions, review gates, or auto-drive.
- A new default-on gate that can stop unattended runs.
- Any change to project-specific zones outside the allowlist.
- Any body change when the current live hash is not a known predecessor.
- Missing compatibility metadata.

The auto-sync implementation should write atomically, then record a
`meta/protocol-sync_<timestamp>.md` entry with old/new hashes, installed package
version, role, change class, preserved zones, and a short diff summary. Same-date
collisions are likely, so use an ISO-ish timestamp or include the new hash prefix,
not just the date.

For this source repo, the implementation phase must update both
`parley-deck/COOPERATION.md` and `internal/protocol/defaults/COOPERATION.md`. The
consumer auto-sync path should not blindly assume every downstream repo has the
embedded default copy.

### 3. Roster ping depends on a stable roster-ID to runtime-ID map

The live protocol roster uses stable participant IDs like `codex-1` and
`antigravity-1`, while the current CLI runtime specs use IDs like `codex` and
`agy`. A preflight that records `agy` as excluded while `00-prompt.md` lists
`antigravity-1` has already lost auditability.

Do not parse the human prose in the roster Role column to infer commands. Add or
reuse a machine-readable mapping, for example:

```toml
[roster.codex-1]
runtime = "codex"

[roster.antigravity-1]
runtime = "agy"
```

The preflight report and `00-prompt.md` must use roster IDs for quorum and
canonical artifacts. Runtime IDs are implementation details in the probe evidence.
If a roster ID has no runtime mapping, classify it as `not_configured`; do not
silently skip it or substitute a similarly named runtime.

When `--participants` is supplied, distinguish operator selection from
availability exclusion. My preference: ping all active roster entries for the
readiness table, but apply exclude/reinclude gates only to the effective
participant set for this idea. Omitted roster agents should be recorded as
`not_selected`, not as unavailable exclusions.

### 4. The PONG probe must reuse the real invocation path

The probe should exercise the same configured launch mechanics as a real round,
but with a narrow sentinel contract. Avoid ad-hoc calls; this repo has known CLI
edge cases where ad-hoc invocations produce false hangs:

- Codex can hang if stdin is mishandled.
- Claude can swallow a prompt behind variadic flags.
- Agy buffers stdout until process exit.
- Hermes may take roughly 40 seconds to initialize and can silently die while
  `--version` still works.

The probe should:

- Use the same command construction as the runner, including prompt mode, root,
  isolated homes, and configured args.
- Write a unique sentinel file under a non-canonical probe directory, not under an
  idea round path.
- Treat success as "process exits within timeout and sentinel file has exact
  content".
- Classify failures as `not_found`, `version_failed`, `pong_timeout`,
  `pong_exit_error`, `pong_bad_sentinel`, or `probe_setup_error`.
- Kill the whole process group on timeout so wrappers do not leak children.
- Run probes concurrently with per-agent timeouts and a global deadline.
- Preserve enough logs for diagnosis, but do not let probe artifacts count as
  canonical Parley artifacts.

`parley agents verify --full` is too heavy to reuse directly as the idea-start
probe. It writes runtime probe files and, for Codex, performs a Git smoke test
that can write into `.git`. Preflight should be cheaper and side-effect bounded.

Timeouts should be configurable per runtime, with a default that avoids false
negatives for Hermes. A 60 second default may be marginal; 90 seconds with a hard
global cap is safer. "Slow" and "hung" should both be unavailable for quorum, but
the reason should record whether the process was killed on timeout.

### 5. The two user-confirm gates need explicit outcomes

There are four common roster cases:

- Available, no recent exclusion: include without a gate.
- Unavailable: ask to exclude for this idea. If the user declines, abort rather
  than launching a doomed round.
- Previously excluded and now available: ask to re-include. If the user declines,
  exclude for this idea and record that as operator-deferred reinclude.
- Previously excluded and still unavailable: ask to exclude again for this idea,
  because the exclusion is temporary and per-idea.

The "previously excluded" detector should scan prior `00-prompt.md` records for an
agent exclusion that has no later matching `reincluded` record. To avoid a
forever-confirm loop, a confirmed reinclude in a later idea should close that
prior exclusion chain. If parsing prior records fails, ask rather than guessing.

`00-prompt.md` should carry both the effective quorum and the readiness audit:

```yaml
participants: [claude-1, codex-1, hermes-1]
excluded: ["antigravity-1 | pong_timeout | confirmed_by=user | confirmed_at=2026-06-19"]
reincluded: ["hermes-1 | prior_idea=previous-slug | confirmed_by=user | confirmed_at=2026-06-19"]
```

The body can include a richer `## Readiness` table for humans, but the frontmatter
should stay simple enough for existing line-oriented parsers. If the last
non-facilitator would be excluded, the exclusion gate can be confirmed, but the run
still cannot proceed without the existing recorded solo exception.

### 6. `parley preflight` should be usable alone and embeddable by `parley run`

Suggested surface:

```text
parley preflight [--dir DIR] [--participants IDS] [--json] [--check-only]
                 [--ping-timeout D]
                 [--confirm-exclude IDS] [--confirm-reinclude IDS]
```

Output should include:

- Protocol role, live hash, packaged hash, installed version, and proposed action.
- Agent readiness table keyed by roster ID, with runtime ID, command path, version,
  PONG result, duration, and reason.
- Pending gates, if any.
- Effective participants after confirmed decisions.
- The exact records to place in `00-prompt.md`.

Exit codes should be boring and scriptable:

- `0`: ready; any allowed additive sync has completed or `--check-only` verified it.
- `1`: fatal runtime/check error not representable as a user gate.
- `2`: usage error.
- `3`: pending user confirmation or solo-exception requirement.

Do not overload `--yes` to answer exclusion or reinclude gates. `--yes` can mean
"it is OK to contact hosted backends", but roster changes are semantic and need
specific confirmation. For unattended runs, `--confirm-exclude` and
`--confirm-reinclude` should only apply if the current probe result still matches
the reported agent and reason; stale confirmations should fail.

### 7. Wiring into `parley run` must stop before idea creation

Today `parley run` discovers agents, selects participants, then calls
`runcontrol.Create`, which creates `00-prompt.md` and the run before round-01
starts. Preflight should run before `runcontrol.Create` so that:

- The protocol sync can update `COOPERATION.md`, then workspace status can be
  re-read before participants are resolved.
- The effective participants and exclusion/reinclude records can be written into
  the initial `00-prompt.md`.
- An unresolved gate can return exit code `3` without creating a half-open idea.
- Auto-drive never starts with a pending human gate.

For TUI mode, do not prompt from a background goroutine after the live run starts.
Resolve preflight first, then launch round-01. For noninteractive or auto-drive
runs, pending gates are a clean stop, not a deadlock: print JSON or text guidance
and exit before any agent is launched.

## Risks

- **Protocol regression:** a consumer-style auto-sync in a source repo would replace
  newer canonical text with older packaged text. This is the highest-risk failure
  and should be impossible when `protocolRole=source` or role is unknown.
- **False quorum:** runtime IDs and roster IDs can drift. Without a machine-readable
  mapping, the tool may probe one identity and record another.
- **False availability:** version probes pass when a backend cannot answer. The
  sentinel PONG reduces this, but timeouts, buffered stdout, isolated homes, and
  leaked child processes need careful handling.
- **Unattended deadlock:** any unresolved gate that waits on stdin, the TUI, or
  auto-answer will hang automation. Gates need deterministic exit/report behavior.
- **Hidden cost:** PONG probes may call hosted backends. The report should make that
  explicit, and standalone preflight should not surprise users with broad expensive
  probes.
- **Partial writes:** protocol sync plus idea creation can leave a dirty or mixed
  state if the process is interrupted. Use atomic writes, re-read state after sync,
  and avoid creating the idea until gates are resolved.
- **Perpetual reinclude prompts:** if prior exclusions are not closed by a
  `reincluded` record, the same agent may require confirmation forever.
- **Solo collapse by workflow accident:** confirming exclusions is not enough when
  no non-facilitator remains. The existing solo-exception rule must still fire.
