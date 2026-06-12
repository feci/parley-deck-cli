# Audit

Audit concrete work with a fresh outside reviewer before calling it done — a pass/fail gate, not advice (that is `consult`). The opposite agent CLI reviews the work (Codex reviews Claude's, Claude reviews Codex's) and catches blind spots a same-model self-review cannot. Resolve `scripts/agent.sh` (in this skill) to an absolute path and run it from the target project with the reviewer explicit: Codex hosts pass `--reviewer claude`, Claude hosts pass `--reviewer codex`; omit `--reviewer` only after verifying host-agent markers are visible in the shell (the script warns and records `same-reviewer: yes` when reviewer and host match). Always run the script — never a hand-rolled `claude`/`codex` call, and never with a spend cap — so its verified flags and guards stay in force. It runs the opposite-model reviewer read-only at full reasoning depth with the operator's curated MCP servers available, and writes a provenance-headed report to `.local/audits/` plus a run-ledger line to `.local/audits/index.jsonl`. Launch it with the host's reliable long-running task mechanism and supervise it until the report path lands; `references/reviewers.md` carries the per-runtime launch and the CLI mechanics.

## Gate Standard

The gate closes only on a fresh full-scope pass — never a resumed session — whose verdict counts zero findings of any severity and any kind: no High/Medium floor, no "remaining Lows accepted", no nitpick allowance. Fix everything each pass reports — Lows, nitpicks, and pre-existing issues included — until a closing pass comes back literally clean or the operator personally rules a finding closed. Briefs carry dispositions as context for the reviewer to weigh openly, never as suppression: no "do not re-raise" carve-outs, no severity floors, no instruction that narrows what a closing pass may report. A disputed finding closes only by the operator's call, relayed verbatim.

## Scope and Snapshot

Scope the gate with `--scope uncommitted` (default), `commit:<sha>`, or `base:<ref>`. Read-only audits review a disposable snapshot checkout the script creates automatically — uncommitted work is captured as a snapshot commit, committed scopes get a clean HEAD checkout — so the live tree stays free: keep working there while the review runs, and judge each finding against the code as it stands when the report lands. If the script falls back to a live-tree review (`snapshot: live-fallback` in the header), keep the tree frozen as it instructs. Add `--exec` only when static reading cannot judge the work (behavioral changes, test suites, infrastructure); an `--exec` audit runs in the live tree under a report-don't-fix contract, so hold that tree frozen until the report lands. Once fix commits land on top of a `commit:` unit, widen the scope to `base:<sha>^` so every further pass covers the reviewed commit and its fixes; the script warns when a `commit:` scope no longer covers HEAD.

## Brief

Brief the reviewer over stdin (`- <<'EOF'`) with neutral intent: what the work should do, its constraints, validation output as labeled raw data quoted inline (never by `.local/` path — the snapshot excludes `.local/`), and any dispositions of earlier findings with their rationale. Never claim the work is correct, tested or safe, and never include secrets. The reviewer adjudicates dispositions in the open — expect the report to state whether it concurs.

## Loop

1. **Open** — `agent.sh --slug <unit> --reviewer <cli> --scope <scope> [BRIEF]`; the first pass of a unit is always a fresh session.
2. **Work the wait** — while the reviewer runs, do an independent adversarial self-review of the same scope; the caller's own pass regularly finds real issues the reviewer misses. Triage your own findings with the same severity and kind taxonomy as the reviewer's.
3. **Triage** — verify every reported finding against the code before fixing (a reviewer can be confidently wrong); relay decision findings to the operator unanswered; record rebuttals you intend to argue in the next brief; record operator-deferred findings in `UPDATE.md`.
4. **Fix** — apply fixes for everything, at every severity; run local validation.
5. **Converge** — mid-gate, resumed passes (`--resume <session>`, paired with `--reviewer` whenever the session is not the default reviewer's) verify the worked findings cheaply: the reviewer re-reads the changed files and reports everything still open, at any severity. Resumed passes converge the gate; they never close it.
6. **Close** — a fresh full-scope pass over the whole unit. It closes the gate only when its verdict is zero findings; anything it reports loops back to triage.

## Weighing the Report

The reviewer is an outside opinion with less context than the caller, not an oracle: it read the repository cold, and an authoritative tone is not evidence. Verify each finding against the actual code, accept what holds, rebut what does not (openly, in the next brief), and never let the report's framing replace your own judgment of the work. The same calibration applies in reverse: do not dismiss a finding because it is inconvenient — the gate exists to surface exactly those.

## Stopping Judgment

Default to two passes — an opening audit and a closing fresh pass — then read the trajectory, not a counter, up to a recommended five passes with no hard cap:

- **Converging** (keep going): findings get fewer, lower-severity, and confined to the newest code.
- **Churning** (stop and escalate with the trajectory): fresh High or Medium findings keep landing on code the fix passes themselves introduced, or the same ground re-litigates despite open rebuttals.
- **Blocked** (escalate immediately): a decision finding pauses its thread until the operator answers; everything else proceeds.

## Panel

For high-stakes gates, close with both reviewers in parallel on the same scope and brief (`--reviewer codex` and `--reviewer claude`). Consensus findings are near-certain; divergence marks judgment calls worth surfacing to the operator. Panels are expensive — reserve them for security-sensitive or architectural closures.

## Provenance

Reports accumulate in `.local/audits/` as `audit-<timestamp>-<slug>.md` with provenance headers (scope, ref, reviewer, model, session, snapshot, timestamp), and `index.jsonl` there records every run — reports, failures, and interruptions. They are the audit trail for how a gate closed and survive session loss: never delete them mid-project, and reference the closing report filename in the commit message where the fixes land.
