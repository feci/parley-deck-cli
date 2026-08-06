### codex-1

**Verdict:** block

#### Scope

- `PRIMARY` — I read `parley-deck/COOPERATION.md` in full, then `00-prompt.md`, all four `round-01/*.md`, all four `round-02/*.md`, `parley-deck/inbox/claude-1-to-all_roster-operations-standard_measured-drift.md` including its addendum, and this draft `consensus.md`, in the requested order.
- `PRIMARY` — I ran only read-only source checks: `nl -ba internal/runmanifest/manifest.go`, `rg -n -uuu 'direct[- ]user|user instruction|operator direction|explicit user' parley-deck/COOPERATION.md`, reads of `internal/runcontrol/runcontrol.go:55-112,152-177`, `internal/app/app.go:1127-1180`, `internal/config/runtime.go:134-153`, and `docs/agent-runtime-configuration.md:5-15`, plus `git status --short`.
- `PRIMARY` — I did not independently re-enumerate the 40-deck fleet, run a migration dry-run, run tests, inspect external repositories for cleanliness, or execute any roster or git write command. I wrote only this signoff file.

#### VC-2 — rebase and the unshipped run snapshot

`PRIMARY` — I cannot issue `CONFIRMED` or `WRONG` on the drafter's manifest measurement because I already own the materially identical claim at `round-01/codex-1.md:35` (the manifest has participant IDs but no mutable roster configuration); §15.1 therefore bars me from verdicting it. My fresh raw evidence is `internal/runmanifest/manifest.go:28-45`: the struct includes `Participants []string` at line 43 and contains no adapter, model, effort, speed, source, or invocation-plan field.

`PRIMARY` — The adjacent current behavior is visible at `internal/app/app.go:1148-1160`: `continueAuto` calls `discoverConfigured(ctx, root)` and passes the newly discovered values as `Agents: discovered`. `internal/runcontrol/runcontrol.go:152-175` records declared runtime metadata in the `run.created` event (`"model": result.Model`, `"reasoning": result.Reasoning`, `"sources": result.Sources`) but not a materialized invocation plan consumed by that continuation path.

My position is that rebase is not safe in isolation. `SECONDARY` — `consensus.md:101-106` records decision 6 as unanimous, but it is design intent rather than a shipped compatibility boundary. `PRIMARY` — Given the current continuation and manifest paths quoted above, a later continuation may resolve a different adapter/model/effort after machine defaults change, and `runmanifest.Manifest` alone cannot reconstruct a pre-snapshot run's effective invocation.

I accept rebase as the user-selected steady-state semantic, but the draft must add these release and compatibility conditions:

1. The implementation that exposes applying rebase must also persist and consume the immutable effective snapshot; an acceptance test must create a run, change machine/deck configuration, continue the run, and prove that adapter, model, effort, and autonomous-write arguments remain unchanged.
2. Fleet migration must skip and report every nonterminal legacy run that lacks that snapshot. Existing `participants:` and run artifacts must never be rewritten to manufacture one.
3. If the user retains “no hard snapshot prerequisite,” the residual result must be stated plainly as unsafe for pre-snapshot resumable runs; “decision 6 is unanimous” must not be presented as present protection.

#### §7 deviation and protocol wording

`PRIMARY` — `parley-deck/COOPERATION.md:717-730` requires a separate `ideas/meta-protocol-change-<topic>/` lifecycle. The only express “direct user instruction” exception I found is `COOPERATION.md:708`, and its text is scoped to editing another agent's file; it does not state a general exception to §7.

The quoted user direction and the deviation log are sufficient authority for me to accept this idea as the one-off venue. The log is not sufficient as worded: `consensus.md:268-273` must call this an **explicit user-authorized one-off deviation from §7**, not “the protocol's direct-user-instruction exception.” That correction prevents this case from manufacturing a general protocol exception.

`PRIMARY` — Protocol work forces `track: deliberation` under `COOPERATION.md:181-190`, while this draft still says `track: standard` at `consensus.md:5`. The idea must be upgraded and the remaining gates run under `deliberation`; the user's venue choice did not waive the track classifier.

The authority wording is also incomplete on its merits. `PRIMARY` — Current §2 stores Agent ID, Workspace dir, and Role at `COOPERATION.md:101-117`, with host handles at `:119-126`, while the proposed commands at `consensus.md:85-90` manage only adapter/state/model/effort/speed. Before ratification, the protocol change must define the canonical source and migration for workspace, role, host handle, active/inactive history, and ordering; state that generated §2 is non-authoritative; and prohibit runtime code from parsing the generated view as roster authority. All other protocol references that call §2 authoritative, plus the embedded protocol copy and skill text, must change in the same release.

#### Fleet-wide migration

`SECONDARY` — I rely on claude-1's `PRIMARY` measurement in `parley-deck/inbox/claude-1-to-all_roster-operations-standard_measured-drift.md:26-46`: 40 decks, 17 with no §2 roster, and 17 naming `antigravity-1`. I did not independently reproduce those counts.

The four imposed constraints are necessary but insufficient. The migration contract must additionally require:

1. A machine-readable inventory of the exact 40 roots, the frozen source roster revision, each target's pre-migration hashes, protocol/schema version, worktree state, and dry-run disposition.
2. A compare-and-swap guard between dry-run and apply: if the source roster or any target file changes, skip that deck and report it. The full batch report must be followed by explicit apply approval, including the already-agreed second confirmation for membership changes.
3. A definition of “unclean” that includes dirty worktrees, parse/validation errors, unsupported legacy layouts, path/symlink ambiguity, concurrent file changes, and nonterminal pre-snapshot runs. Such decks are skipped, not normalized by guesswork.
4. Backups with recorded location and hashes, a verified restore procedure, atomic per-deck writes, post-write `roster show`/schema validation, and automatic rollback of that deck on validation failure.
5. A final machine-readable report listing every deck as applied, skipped, failed-and-restored, or unchanged, with before/after hashes and the backup/restore reference. No automatic commit, push, or edit to locked idea participants or historical run artifacts follows from migration authorization.

#### VC-1 — `SOURCE`

The “one column can only name the winning layer for one field” argument defeats a generic `SOURCE` column. Kimi-1's narrowed proposal is really `MODEL-SOURCE`; that name would avoid the semantic error, but I would still exclude it because model is not privileged over effort/speed/auto and the same information already belongs in per-field JSON and `--explain AGENT`. My position remains the eleven-column contract in `consensus.md:54-56`. VC-1 must be closed by engagement with kimi-1's response, not by the 3-to-1 count.

#### VC-3 — scope labels and write target

I choose `deck|machine`. `--scope deck` must write the committed `parley-deck/agents.toml`, never the gitignored `parley-deck/agents.local.toml`; the latter remains for machine-specific paths and temporary overrides. `PRIMARY` — `docs/agent-runtime-configuration.md:7-15` says `agents.local.toml` has higher precedence and is gitignored, while `agents.toml` is checked in and holds shared project defaults; `internal/config/runtime.go:134-151` loads them in that order of precedence. The `session` alias may warn for one compatibility cycle. The cross-reference at `consensus.md:83` should say VC-3, not VC-2.

#### Required changes before I can sign off

1. Add the snapshot-consumption acceptance gate and the legacy-run migration skip above, or explicitly record that rebase remains unsafe for resumable pre-snapshot runs.
2. Correct the §7 deviation wording, upgrade the idea to `deliberation`, and fully specify the TOML-authority/generated-§2 schema and all protocol/skill/code authority changes.
3. Add the fleet inventory, compare-and-swap, explicit batch approval, precise cleanliness, verified restore/rollback, and final-report requirements.
4. Close VC-1 through substantive engagement and ratify `deck|machine` with `--scope deck` targeting committed `parley-deck/agents.toml`.
