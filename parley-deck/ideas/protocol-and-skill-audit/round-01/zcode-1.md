---
agent: zcode-1
idea: protocol-and-skill-audit
round: 1
date: 2026-08-20
---

Lens: COOPERATION.md against itself (contradictions, dead rules, stale references, §4.0-override
consistency). All experiments ran in a copy of the deck at `/tmp/pdaudit/`; the shared tree was not
modified. CLI under test: `dist/parley-v1.45.0-darwin-arm64` (the version the deck header names).
All line numbers refer to the live `parley-deck/COOPERATION.md` unless another file is named.

## Findings

### F1 — §2's roster-authority model is false for the deck it governs: membership does NOT live in `[roster.<id>]` blocks in `parley-deck/agents.toml`
severity: MAJOR
tag: PRIMARY
command: `cat parley-deck/agents.toml` and `parley roster show --dir <deck-copy>`
output: the deck's `agents.toml` contains zero `[roster.*]` blocks and states: "ROSTER MEMBERSHIP
IS NOT DECLARED HERE — BY OWNER INSTRUCTION (2026-08-19). … With no [roster.*] block in this file
the deck INHERITS ~/.parley/agents.toml at read time … `parley roster show` marks this state
INHERITED and every row reports `inherited-roster`." `roster show` (exit 0) lists 6 active agents,
every row with status `inherited-roster`.
contradicts: `parley-deck/COOPERATION.md:103-106` — "The roster's authority is
`parley-deck/agents.toml`, not this table. Membership and each agent's adapter, model, effort and
speed live in `[roster.<id>]` blocks there". §2:124-129 documents exactly one alternative state
(`legacy-roster`, hand-written table); the `inherited-roster` state this deck is actually in appears
nowhere in the protocol.
why it matters: §2 is the section every participant reads to learn who is on the roster. A reader
following it verbatim would open `parley-deck/agents.toml`, find no `[roster.*]` blocks, and
conclude the deck has no roster — or "fix" it by declaring one, which the owner explicitly forbade
on 2026-08-19 (same day the header says the protocol was synced). The protocol text lags the
deck's own ratified state and gives no path to it.

### F2 — §2's "generated view" is permanently empty: `roster render` refuses to write it under the deck's ratified configuration
severity: MAJOR
tag: PRIMARY
command: `parley roster render --dir /tmp/pdaudit/parley-deck` (deck copy, parley 1.45.0)
output: exit 0, no write: "roster render: this deck declares no roster of its own; the 6 rows shown
come from ~/.parley/agents.toml. Writing them into §2 would commit a machine-local roster into a
shared file. Declare the roster with `parley roster set <agent> --scope deck --adapter <family>`,
or re-run with --adopt-inherited to accept the inherited roster as this deck's own"
contradicts: `COOPERATION.md:131-135` — "The generated view:" followed by a table with headers and
zero data rows, in the live deck. Also §2:106 "the edit will not take effect and will be
overwritten on the next render" — under inheritance there is no next render; and Appendix A:1092
"then run `parley roster render` to generate the §2 view".
why it matters: §2 was redesigned because "17 decks carried no roster at all" (§2:110-111). The
flagship deck now carries an empty roster view by design: the one command §2 gives for generating
the view declines to run without `--adopt-inherited`, and adopting was rejected by the owner (see
F1). So the section's human-readable surface shows nothing while `roster show` knows 6 members —
the exact failure mode the redesign was written to eliminate, produced by the redesign itself.

### F3 — Header claims "Protocol synced: 2026-08-19 — parley-deck-skill 2.9.0", but the §9.0-mandated record of that sync does not exist and meta/version.json still says 2.8.0
severity: MAJOR
tag: PRIMARY
command: `head -7 parley-deck/COOPERATION.md; cat parley-deck/meta/version.json; ls parley-deck/meta/`
output: header line 7: "**Protocol synced:** 2026-08-19 — parley-deck-skill 2.9.0 / parley-deck-cli
1.45.0". `meta/version.json`: `"deckVersion": "2.8.0"`, `"source": "npm:parley-deck-skill@2.8.0"`,
`"updatedAt": "2026-08-12T16:55:24.209Z"`. `meta/` contains exactly one protocol-sync record:
`protocol-sync_2026-06-13_v1.3.1.md`. (The skill package really is 2.9.0 and the CLI VERSION file
is 1.45.0, so the header's versions are right — the recorded sync trail is not.)
contradicts: `COOPERATION.md:843-846` (§9.0): a consumer sync "updates the `Protocol synced:`
header line **and records `meta/protocol-sync_<ISO-timestamp>.md`**". If the header was updated by
a sync on 2026-08-19, its record is missing; if the header was hand-edited, the line misstates a
sync that never ran. Either way §9.0's audit trail and the header disagree.
why it matters: §9.0's freshness mechanism is the deck's defense against protocol drift, and its
own bookkeeping is the evidence. An auditor running the §9.0 check today compares the header
(2.9.0, 2026-08-19) against version.json (2.8.0, 2026-08-12) and must conclude the deck's sync
metadata is stale or false — the mechanism cannot verify its own last action.

### F4 — The deck declares `Transport: github-pr` but §11.B's mechanics have never been used: zero `idea/*` branches, 25 of 34 close transactions are direct commits
severity: MAJOR
tag: PRIMARY
command: `git branch -a | grep -c idea/` → 0 (of 72 branches); `git log --all --oneline | grep -c
"FINAL.md + close idea"` → 34; `git log --all --oneline --merges | grep -c "FINAL.md + close
idea"` → 9; `git log --oneline -15 -- parley-deck/` — all direct commits to main, including the
two newest Phase-0 kickoffs: `91dbad1 [claude-1] protocol-and-skill-audit: 00-prompt.md + round-01/
(Phase 0)` and `f06ccbc [claude-1] agents-verify-hermes-probe, preflight-liveness-false-negative:
00-prompt.md (Phase 0)`.
contradicts: header line 5 "**Transport:** `github-pr`" read together with §11.B:939-943 (Phase 0:
create branch `idea/<slug>`, push, open a **Draft PR**), §11.B:920-923 (one long-lived design PR
per idea on branch `idea/<slug>`), and §11.B:970-973 ("**Merges the PR** using **Merge commit** …
The merge IS the close-idea transaction").
why it matters: §0:45 says the transport choice "determines how phase transitions, cross-review,
and signoffs are mechanically performed". In practice the deck operates §11.A local-directory
mechanics (direct commits, no PRs) while printing `github-pr` in the header every participant
reads at session start (§9 step 1: "note the active `Transport:`"). The printed rule binds
nothing. (Charitably, branch deletions per §11.B step 7 could hide old `idea/*` branches — but
they cannot hide merge commits, and only 9 exist; the other 25 closes, and every recent kickoff,
bypassed §11.B entirely.)

### F5 — §9.0 keys its freshness behavior on `meta/version.json` `protocolRole`; the field does not exist in this deck's version.json
severity: MINOR
tag: PRIMARY
command: `grep -c protocolRole parley-deck/meta/version.json` → 0; `grep -rl protocolRole
parley-deck/meta/` → (no files); `cat parley-deck/meta/version.json` shows keys: schemaVersion,
deckVersion, protocolSchema, projectMetadataSchema, source, protocolSha256, skillSha256,
packagedProtocolSha256, compatibilityManifestSha256, updatedAt, updatedBy.
contradicts: `COOPERATION.md:838-848` (§9.0): "Behaviour depends on `meta/version.json`
`protocolRole`: `source` → …; `consumer` → …; `protocolRole` missing/unknown → **do not
auto-write**; ask the user once and backfill the field."
why it matters: the three-way switch §9.0 builds on is undefined for this deck, and the prescribed
fallback (ask once, backfill) evidently never ran — the field is still absent. Whatever `parley
preflight` is deciding, it is not deciding it via the documented field, so the documented rule
describes a mechanism the deck's own metadata does not carry. (The other two §9.0 fields,
`protocolSha256` and `packagedProtocolSha256`, do exist verbatim — this one was never written.)

### F6 — §12.12 cites its ratifying idea by a slug that exists nowhere: `meta-protocol-change-end-to-end-pipeline`
severity: MINOR
tag: PRIMARY
command: `ls parley-deck/ideas/ | grep -i "end-to-end\|pipeline"` → (empty); `git log --all
--oneline -- "*end-to-end-pipeline*"` → (empty — never existed under that name in history either)
output: no idea directory, branch, or historical path matches the slug. The artifact exists only
as `parley-deck/ideas/2026-06-02T12-07-14-meta-protocol-ch/`, whose 00-prompt.md reads
"Meta-protocol-change: evolve Parley Deck into a full automatic idea-to-monitoring pipeline" and
whose FINAL.md contains §12's text.
contradicts: `COOPERATION.md:1155` — "This section was ratified by idea
`meta-protocol-change-end-to-end-pipeline` (2026-06-02)." Also §3:190 "Idea slug rules:
`kebab-case`" and §7:745 "Open an idea under `ideas/meta-protocol-change-<topic>/`" — the actual
dir carries a timestamp + truncated name that violates both.
why it matters: §1's purpose #4 is an auditable trail where every decision "lives in a file".
§12's ratification citation cannot be followed to its evidence; a reader must already know which
truncated directory to open. It also makes the idea unreachable by tooling that validates slugs
(see F10c).

### F7 — The Quickstart's core/reference map omits §15 (binding on every track) and §10, and §10 sits physically between §8 and §9
severity: MINOR
tag: PRIMARY
command: `grep -n "^## " parley-deck/COOPERATION.md`
output: §8 at line 787, **§10 at line 811**, §9 at line 827, §11 at line 877; §15 at line 1230.
Quickstart lines 36-39: "The **core** every participant needs is §0–§8. The rest are **reference
appendices** — skip them until a task needs them: **§9** session-start checklist, **§11**
transport mechanics, **§12** pipelines & action stages, **§13** retrospective optimization,
**§14** automated outer loop."
contradicts: itself — the map claims to enumerate "the rest" but §15 (Verification integrity,
1371-line document's most cross-referenced late addition: blockquoted at the top of Phases 1, 2,
3 and 6, lines 305, 326, 353, 500) and §10 are missing from it. §15.7:1360-1370 binds §15.1–15.5
on **every** track including `fast`.
why it matters: the Quickstart is the declared entry point ("You do not need to read all of this",
line 14) and the progressive-disclosure contract is its promise. A compliant newcomer who reads
only what the map lists will never open the section that defines the provenance rules
(PRIMARY/SECONDARY/RECALL) their own round file is graded against — the tags this very audit is
required to carry.

### F8 — §3's canonical directory layout omits `agents.toml` (§2's authority file) and `runs/` (present on disk, referenced by §12.12)
severity: MINOR
tag: PRIMARY
command: `ls parley-deck/` → `COOPERATION.md  agents.toml  ideas  inbox  meta  runs`
contradicts: `COOPERATION.md:164-188` — the §3 tree lists only COOPERATION.md, ideas/, inbox/,
meta/ (plus per-idea files). `agents.toml` is missing although §2:104 names it the roster's
**authority**; `runs/` is missing although it exists and §12.12:1153 references it
("`ideas/`, `inbox/`, `meta/`, `runs/` are unchanged"); `parley init`'s own help says it creates
"parley-deck/COOPERATION.md, ideas/, inbox/, meta/, and runs/".
why it matters: §3 is the layout section a newcomer uses to build a mental model of the deck. It
predates the roster redesign and the runs store and was never updated; two of the six real
top-level entries are invisible in it, one of them the file §2 declares most authoritative.

### F9 — Phase 8 has a paragraph spliced into the wrong place: the fix-up template is separated from its introducing sentence
severity: MINOR
tag: PRIMARY
command: `sed -n '588,605p' parley-deck/COOPERATION.md`
output: line 588 ends "On completion, they append a new section to `IMPLEMENTATION.md`:" — followed
not by the template but by an unrelated paragraph (line 590: "When `checks:` is a list (LE-4
completion contract), closing additionally requires the latest driver run to be all-pass at the
current HEAD: …") and only then, at line 592, by the promised `## Fix-up cycle N` block. The LE-4
paragraph's content ("closing additionally requires…") semantically belongs with the closing rules
at line 605.
contradicts: the document's own structure — §4 Phase 8 (lines 588-605); every other template in
§4 (00-prompt at 273-292, round at 310-315, consensus at 360-372, IMPLEMENTATION at 443-491)
immediately follows its introducing colon.
why it matters: this reads as a mis-merged edit, and it lands in the one phase the auto-driver
implements mechanically. A reader (or a driver author) scanning Phase 8 for "what goes in the
fix-up section" gets a close-veto rule spliced into the middle of the instruction, and the close
rule it actually belongs to appears 15 lines later without it.

### F10 — `parley learn` (§13.5) deviates from its documented contract in three measurable ways
severity: MINOR
tag: PRIMARY
command + output (deck copy, parley 1.45.0):
(a) `parley learn automation-outer-loop --dir /tmp/pdaudit/parley-deck` → exit 1, "learn: no idea
at /tmp/pdaudit/parley-deck/**parley-deck**/ideas/automation-outer-loop" — the deck-dir `--dir`
form that `roster show`, `roster render`, and `retro scan` all accept (verified exit 0) is doubled
by `learn`; only the workspace root works (`--dir /tmp/pdaudit` → exit 0, "Wrote advisory playbook
parley-deck/playbooks/automation-outer-loop.md").
(b) the successful run distilled `automation-outer-loop`, whose `00-prompt.md` says
`status: round-01` — not a completed idea.
(c) `parley learn 2026-06-02T12-07-14-meta-protocol-ch --dir /tmp/pdaudit` → exit 2, "learn:
invalid idea slug … (lowercase kebab-case)" — the deck's three completed 2026-06-02 ideas are
unusable with `learn` forever.
contradicts: `COOPERATION.md:1185` (§13.5): "`parley learn <closed-idea-slug>` scaffolds … from a
COMPLETED idea" (b); and the CLI's own uniform `--dir` convention (a).
why it matters: §13.5's whole design is "deterministic skeleton from a completed idea"; the tool
neither checks completion nor accepts the deck directory every sibling command takes, and the
deck's oldest completed ideas fail its slug validation outright. The advisory output itself is
fine (exactly one file written, advisory banner printed — verified), but three of the section's
four concrete claims about the command do not hold.

### F11 — §11.B prescribes committing IMPLEMENTATION.md "directly … (small, no PR needed)" while its own recommended branch protection requires PRs for all changes to `ideas/`
severity: MINOR
tag: PRIMARY
command: `sed -n '979p;1013p' parley-deck/COOPERATION.md`
output: line 979 (Phase 5, step 3): "commits `IMPLEMENTATION.md` directly to the integration
branch of the parley-deck repo (small, no PR needed)"; line 1013 (Branch protection,
recommended): "Require PRs for all changes to `ideas/`."
contradicts: itself — `IMPLEMENTATION.md` lives at `ideas/<slug>/IMPLEMENTATION.md` (§3:177), so
under §11.B's recommended protection the prescribed Phase 5 push is rejected. Phase 6 (line 985)
hedges the equivalent review-file commit with "(or via a small PR if branch protection requires)";
Phase 5 asserts the opposite.
why it matters: a deck that follows both the flow and the hardening advice gets a blocked
implementer at Phase 5. The two paragraphs were written against different assumptions and were
never reconciled.

### F12 — Appendix A instructs filling in header fields that do not exist in the header
severity: MINOR
tag: PRIMARY
command: `sed -n '3,8p;1091p' parley-deck/COOPERATION.md`
output: header (lines 3-8) carries: Workspace, Parley deck, Transport, Created, Protocol synced,
Status. Appendix A step 3 (line 1091): "Fill in the header: workspace name, **shared channel
path**, transport, creation date, **bootstrapping agent ID**."
contradicts: the header template at lines 3-8 — there is no shared-channel field and no
bootstrapping-agent field (nor does any other header line carry them).
why it matters: Appendix A is the adoption path for new projects; two of its five fill-in
instructions point at fields the document's own header does not have. The instructions predate a
header redesign and were not updated.

### F13 — §9.0 requires the readiness result to be recorded in the new idea's 00-prompt.md; the deck's three newest ideas contain no such record
severity: MINOR
tag: PRIMARY
command: `grep -n "excluded\|readiness\|preflight\|liveness" parley-deck/ideas/{agents-verify-hermes-probe,preflight-liveness-false-negative,protocol-and-skill-audit}/00-prompt.md`
output: the only matches are the preflight idea *describing* a preflight defect (its evidence, not
a kickoff record) and the audit idea's defect table quoting it. No `excluded:` field, no
available/unavailable table, no readiness result in any of the three frontmatters or bodies.
contradicts: `COOPERATION.md:831-834` (§9.0): "the facilitator runs a readiness check … **and
records the result in the new idea's `00-prompt.md`**"; §9.0:856 specifies the recording shape
(`excluded: [<roster-id> — reason — confirmed <date>]`).
why it matters: §5:720-721 says "Quorum is set at the §9.0 pre-idea readiness check and locks
once Phase 0 completes" — but for these three ideas (all standard-track, where the ping is
"full") there is no recorded result to lock from, and no evidence any agent was excluded or
confirmed. The quorum-locking step is unlinkable to its input. (Compliance measurement across the
whole deck is hermes-1's lens; I report only the three ideas opened alongside this audit.)

### F14 — 27% of closed ideas carry a stale `status:` in 00-prompt.md, feeding §6 rule 5's stale-round guard false data
severity: MINOR
tag: PRIMARY
command: loop over `parley-deck/ideas/*/`: 78 dirs contain `FINAL.md`; for each, read `status:`
from `00-prompt.md`.
output: 21 of 78 closed ideas have a status other than final/complete/abandoned — 16 say
`round-01`, one `open` (meta-protocol-change-fusion-execplans), one `kickoff`
(launch-mkdir-resilience), one empty (launch-orphan-hardening), one `round-01` with a leading
space variance. Five of the six ideas cited in COOPERATION.md as having *ratified* sections are in
this set (automation-outer-loop §14, rho-retrospective-optimization §13, fusion-execplans §13,
parley-learn-playbooks §13.5, global-core-protocol §7); only verification-integrity (§15) was
bumped to `final`.
contradicts: `COOPERATION.md:422` (Phase 4: "Update `00-prompt.md` `status:` to `final`") and
§6.5:740 ("Before working on an idea, re-read `00-prompt.md` `status:` to avoid writing into a
closed round").
why it matters: the status field is the protocol's own collision-avoidance input, and for one in
four closed ideas it reports the idea as sitting in round 1. An agent following §6 rule 5 on, say,
`automation-outer-loop` (status `round-01`, FINAL.md shipped, §14 ratified) would conclude the
idea is open and write into a closed round. Nothing enforces the bump — a small instance of the
audit's target class: a printed rule with no enforcement.

### F15 — Two functional commands are invisible in `parley --help`: `learn` (documented in §13.5) and `preset list` (cross-referenced by `run`'s own flag help)
severity: NIT
tag: PRIMARY
command: `parley --help | grep -ci learn` → 0; `parley --help | grep -ci preset` → 0; then
`parley learn <missing>` → exit 1 "learn: no idea at …" (works); `parley preset list` → exit 0
"No roster presets defined. Add a [rosters.<name>] block to ~/.parley/agents.toml or
parley-deck/agents.toml." (works).
contradicts: the help text's implicit claim to be the command surface; `-preset`'s own help line
says "see `parley preset list`", a command the same binary never lists.
why it matters: §13.5 sends users to `parley learn`; anyone discovering the CLI through `parley
help` finds no trace of it (or of `preset`). Minor, but it is the CLI misdescribing its own
surface.

## What I checked and found clean

- **Every `parley` command COOPERATION.md references exists in 1.45.0** [PRIMARY — help-tree
  diff]: init (§0:57), protocol status|render|check|publish (§7:761), roster show|set|sync|render|migrate
  (§2:115-129), consensus status/draft/signoff/request-signoffs/finalize/reopen, context repo-map,
  sessions list|inspect, consult + consults list (§8:806), retro scan|select|diagnose|propose
  (§13:1187), preflight (§9.0:834), run, loop tick (§14:1196) — all present in `parley --help`.
- **`parley run --max-driver-steps` / `--max-wall-clock` exist exactly as §14:668-669 describes**
  [PRIMARY — `run --help`], including "explicit 0 = unlimited; omit to use ~/.parley
  [defaults.loop]".
- **`roster show` prints one fixed table** (§2:115 "the canonical answer, one fixed table")
  [PRIMARY] — and its `inherited-roster` per-row status matches what agents.toml's comment
  promises (the CLI knows the true state; only §2's prose doesn't — see F1).
- **`meta/version.json` contains `protocolSha256` and `packagedProtocolSha256`** with the exact
  field names §9.0:838 uses [PRIMARY]. (`protocolRole` does not — F5.)
- **`meta/protocol-changelog.md` exists** (§9 step 1 references it) [PRIMARY].
- **Header version numbers match reality**: skill `package.json` = 2.9.0, CLI `VERSION` = 1.45.0
  [PRIMARY]. Only the sync *trail* is inconsistent (F3).
- **`parley learn`'s write path honors §13.5's output contract** [PRIMARY]: wrote exactly one
  file, `parley-deck/playbooks/automation-outer-loop.md`, with the advisory banner — in a copy.
- **Ratification citations that resolve** [PRIMARY]: §13/§14/§15/§7's cited ideas
  (rho-retrospective-optimization, fusion-execplans, parley-learn-playbooks,
  verification-integrity, global-core-protocol, automation-outer-loop) all have directories with
  FINAL.md; only §12.12's slug is dead (F6).
- **§4.0's fail-safe bands are internally consistent** [PRIMARY — read]: the "~300–~1000 LOC gap"
  and "6–14-file band" (§4.0:213-214) match the classifier thresholds (fast ≤ ~300 LOC / 3-5
  files at line 209; deliberation > ~15 files / ~1000 LOC at line 209). §15.7's table matches
  §15.6's text (fast exempt, standard = round-02 section, deliberation = assigned artifact).
- **§4.0's override statement (lines 231-235) vs the phases**: the per-track table's deltas
  (fast: 1 reviewer, skipped cross-review, collapsed FINAL; standard: 2 reviewers) are correctly
  flagged as overrides; I found no phase text that contradicts the table *without* being covered
  by the override clause.

## What I could not check, and why

- **§11.C GitLab mechanics** — no GitLab project exists to measure against; text-only review is
  SECONDARY and the audit bar requires running commands.
- **`parley-deck-skill status`, installer, doctor, install targets** — kimi-1's assigned lens; not
  run to avoid duplicating that audit.
- **The three COOPERATION.md copies against each other** (skill package's, embedded default's,
  live deck's) — claude-1's assigned lens; I audited only the live deck copy.
- **`parley protocol publish` terminal-refusal behavior** (§7:761-778) — it requires an attended
  terminal and would create a release in `~/.parley`; not run on a live store, and building a
  fake HOME for it exceeded this round's scope.
- **`parley preflight` live ping behavior** — probes the real rostered agents; the false-negative
  is already the subject of the open idea `preflight-liveness-false-negative` and re-probing live
  agents from this audit would duplicate it.
