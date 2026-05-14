---
agent: claude
idea: interactive-agent-mode
review-round: 1
date: 2026-05-14
reviewed-commit: 348b693
---

## Summary
The slice implements the locked-in consensus targets: `launch_mode` plus the per-mode command shape, three modes with `headless` default, a shared handoff packet writer, polling + post-validation for interactive signoffs, manual handoff with pending events, stricter Round 1 artifact validation, and `resume` for pending signoffs. Mode/config semantics match FINAL.md and the dry-run / `agents list` surface exposes the resolved invocation. The most load-bearing gap is that the append-only signoff validator runs only on the headless and successful interactive paths — it is bypassed for manual mode and for `parley resume`, which contradicts FINAL.md's "must not rely on file existence alone" rule.

## Findings

### [CRITICAL] Manual + resume bypass append-only signoff validation
FINAL.md and consensus.md require that consensus/review signoff validation "reuses the append-only signoff validator and canonical status rules" and that "Interactive/manual completion must not rely on file existence alone." In `internal/app/consensus_request_signoffs.go`:

- The headless path (`runHeadlessSignoffAgent` → `validateRequestedSignoff` → `validateAppendOnlyContent`, line 181 and 651–667) does enforce append-only against `beforeRaw` captured inside the same invocation.
- The interactive path enforces it the same way after the polling loop exits (line 181).
- The manual path (`runManualSignoffAgent`, line 499–519) returns `Pending: true` *before* any append-only check; the calling loop then `continue`s past `validateRequestedSignoff` (line 169).
- `resumePendingConsensusSignoffs` in `internal/app/app.go:510–577` only walks events and `consensus.Status` (which surfaces parse errors but not unauthorized edits). It cannot run `validateAppendOnlyContent` because no pre-handoff snapshot, hash, or checksum is recorded in the handoff event.

Why it matters: a user driving a manual handoff (or finishing an interrupted interactive one) can edit or delete other participants' existing signoff blocks; `parley resume` will still mark the run completed as long as the file parses and the agent's own signoff header is present. This is the exact failure mode the consensus design called out as the reason for the explicit `resume` verb (see consensus.md "do not rely on artifact exists alone").

Suggested fix: at handoff write time, persist a `consensus_before_sha256` (or just `before_bytes_len` + sha256) on the `agent.handoff.pending` / `agent.handoff.started` event for the consensus path. In `resumePendingConsensusSignoffs`, read the post-handoff file, recompute `validateAppendOnlyContent` by re-reading the prefix bytes and verifying the recorded hash, then call the existing canonical-status checks already in `validateRequestedSignoff`. Failing those should leave the handoff pending with a clear error rather than emit `agent.handoff.completed`.

### [MAJOR] Handoff "Suggested command" leaks unsubstituted placeholders
`internal/runner/handoff.go:54–66` only expands `{prompt_path}` when rendering the suggested command. The real launch path uses `expandInteractiveArgs` (`internal/app/consensus_request_signoffs.go:564–573`) which also substitutes `{root}` and `{target_path}`. If an operator configures `interactive_args = ["--add-dir", "{root}", "--write", "{target_path}"]` and runs in the default `print-only` invoke mode, the user copies a command line that still contains the literal strings `{root}` and `{target_path}`. The substitutions diverge between the printed handoff and the spawned command, which is exactly the kind of "what `parley` is doing technically" mismatch FINAL.md tries to avoid.

Suggested fix: in `handoffInstructions`, expand `{root}` to `opts.Root` and `{target_path}` to `opts.TargetPath` using the same code path as `expandInteractiveArgs` (extract a small helper into one of the two files and share it), so the printed command, the dry-run preview, and the `spawn-tty` exec all see the same string.

### [MAJOR] Manual handoff resume hint points at the wrong command
`runManualSignoffAgent` (`internal/app/consensus_request_signoffs.go:508`) and `writeSignoffHandoff` (`:534`) both tell the operator the resume command is `parley consensus status <slug>`. But `consensus status` does not advance any handoff state — only `parley resume <RUN_ID>` triggers `resumePendingConsensusSignoffs` which is the new validation entry point. A user following the on-screen instructions will see "consensus=partial" and not know that the appended signoff is still flagged as `agent.handoff.pending` until they run `resume`.

Suggested fix: pass the `runID` into `writeSignoffHandoff` and emit `Resume after the signoff is appended: parley resume <runID>` for the manual path; keep the `consensus status` line as a separate "Inspect:" hint if desired.

### [MINOR] `ValidateArtifact` hardcodes Round 1 sections while accepting a round parameter
`internal/runner/validation.go:12–41` takes `round int` for the frontmatter check but the required-section list is the Round 1 shape (`## Summary`, `## Proposed approach`, `## Concerns / open questions`, `## Risks`). `runner.RunRoundOne` is the only current caller and always passes Round 1, so existing headless behavior is not broken. But the signature invites future callers (mixed-mode rounds, review-round validation) to reuse it for rounds with different section contracts, where it will silently reject every valid artifact. IMPLEMENTATION.md acknowledges this as deferred, but the validator surface is the source of the fragility.

Suggested fix: either fold the section list inside a `roundSections(round int) []string` helper that errors for unknown rounds, or rename `ValidateArtifact` to `ValidateRoundOneArtifact` until Phase 2 wires a real round-aware validator. Either makes the contract honest.

### [MINOR] Manual signoffs return exit code 0
`runConsensusRequestSignoffs` returns 0 even when `pending` is non-empty (`internal/app/consensus_request_signoffs.go:199–203`). For CI, an `all-manual` invocation of `parley consensus request-signoffs` looks indistinguishable from a fully signed-off consensus. Consensus round-02 deferred the exact code to implementation, but "deferred" was about the value, not about whether to distinguish it. Today, automation cannot tell "pending — run resume" from "complete."

Suggested fix: pick a small distinct non-zero exit code (e.g. 3) for "pending handoffs outstanding," document it in `docs/agent-runtime-configuration.md`, and route `pending != nil` through it. The "Requested signoffs pending:" line already exists for the human-facing case.

### [MINOR] `validateLaunchModes` validates interactive enums for `launch_mode = "manual"` too
`internal/app/consensus_request_signoffs.go:337–356` rejects any non-headless agent whose `interactive_prompt_mode` is not in `{none, file, arg}` or whose `interactive_invoke` is not in `{print-only, spawn-tty}`. For `launch_mode = "manual"` the manual path (`runManualSignoffAgent`) never uses `interactive_invoke` (no TTY spawn) and only uses `interactive_prompt_mode` indirectly via the handoff packet contents, but a misconfigured manual agent will refuse to launch. This is probably fine, but the check should at least be documented as "manual still respects interactive_* fields because they shape the handoff packet" — otherwise a future operator setting `interactive_invoke = "tmux-attach"` on a `manual` agent will get a confusing rejection.

Suggested fix: either skip the invoke-enum check for `manual` (it's never read), or add a one-line code comment explaining why manual reuses the same structural guard.

### [NIT] Dead fallback in `runInteractiveTTY`
`internal/app/consensus_request_signoffs.go:542–548` walks three fallbacks for `command`: `agent.InteractiveCommand` → `agent.Path` → `agents.InteractiveCommandOrDefault(agent.Spec)`. By the time we reach this function, `requestSignoffAgents` (line 256–258) has already rejected `!agent.Found`, so `agent.Path` is always non-empty and the third fallback is unreachable. Worth removing or asserting.

### [NIT] `IMPLEMENTATION.md` frontmatter key drifts from protocol
COOPERATION.md §11.B / §4 Phase 5 specifies `implementer: <agent-id>`. The published `IMPLEMENTATION.md` uses `implemented-by: codex`. Functionally equivalent and not load-bearing, but the protocol's canonical-trail rule benefits from consistent keys.

## Open questions
- For mixed-mode round execution (the deferred slice), how should the runner capture the pre-handoff state for round artifacts so the same append-only invariant can be enforced on `resume`? Round files aren't append-only the way `consensus.md` is, so this needs a separate validation contract — likely just the strict frontmatter + sections check that already exists, but it should be decided before the next slice lands.
- Is there appetite to store a `before_sha256` field on `agent.handoff.pending` events now so the CRITICAL above can be backfilled without an event-schema migration?
- Should the `--allow-interactive-noninteractive` non-TTY opt-in that round-02 discussed actually ship now? Today `runInteractiveTTY` rejects non-TTY with a hard error; `print-only` interactive does not require a TTY at all, so the override may already be unnecessary — worth confirming intent in consensus.
