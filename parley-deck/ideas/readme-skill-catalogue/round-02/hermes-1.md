---
idea: readme-skill-catalogue
round: 02
agent: hermes-1
date: 2026-07-29
---

## 0. What this round is

A cross-review, not a second proposal. I re-read the real README.md and the
shipped files, checked every truth claim the other three made, and then grafted
the best sentences into one finished copy. The forks F1 (catalogue before
Install), F2 (prose not table), F3 (one load-bearing position per add-on), F5
(failure-mode hook), and F6 (separate idea for the CLI README) are unanimous
across round 1 and I agree with all five; I do not re-litigate them below. The
only structural disagreement left is inside F1 — kimi-1 splits the difference
(fastest-path install above the catalogue, full apparatus below) — and I address
that in §4 because it is the one ordering question still open.

---

## 1. Per-participant review

### claude-1

What they got right that I did not:

- The packaging-defect discovery (their §"One thing outside my lens"): they ran
  `npx -y skills@latest add feci/parley-deck-skill --list` and got 1 of 5
  skills, because that CLI shadows nested skills behind a root SKILL.md. I did
  not find this. It is out of scope for the README change but it is the kind of
  thing a reviewer should know before approving a sentence that promises five
  skills through every path.
- The "selection, never averaging" entry for `parley-design` is the tightest in
  the field: "Average four directions together and you get the average of
  four, which is the look everything already has." That is the single best
  sentence in any round-1 entry. I graft it in §3.
- Their cut table is the most specific: it names line ranges and gives a per-
  section saved estimate that sums to ~146, which is the right order.

What is wrong or unsupported:

- The "15 agent runtimes" claim. They hedged it themselves and asked codex-1 to
  check. I checked: the `--help` output lists 15 *targets*
  (codex, claude, agy, gemini, hermes, qwen, codebuddy, goose, kimi, droid,
  vibe, cursor, opencode, aionrs, generic) — 14 named runtimes plus `generic`.
  "15 runtimes" is wrong because `generic` is not a runtime; it is a
  destination. The correct phrasing is "fourteen runtimes plus a generic
  directory." claude-1's instinct to flag it was right; the number should not
  ship as written.
- "Repository Layout … also omits `test/`, `NOTICE.md`, `RELEASING.md`"
  (their cut table). I checked the on-disk tree: `test/`, `NOTICE.md`, and
  `RELEASING.md` all exist. But the current README's layout block (lines
  148–166) does not list them — so the omission claim is accurate *as a
  description of the README*, and the fix is to add them. Not wrong, just worth
  separating "the file omits X" from "the repo lacks X."

### codex-1

What they got right that I did not:

- Exit code 4. The brief lists codes 0–3; codex-1 read the shipped
  `addons/parley-design-check/SKILL.md` (line 61) and found code 4 for
  `UNJUDGEABLE`. I verified it: SKILL.md line 61 and engine.js line 2069 both
  define exit 4. My round-1 entry omitted it. This is a genuine factual gap in
  my file and codex-1 caught it.
- The tracker nuance: "live tracker writes require a separate opt-in
  connector." I verified in `addons/parley-tracker/SKILL.md` line 386: "Actual
  create/update against a live tracker (MCP, REST, CLI) is an **opt-in**
  … connector." My round-1 entry said "mirroring them into Jira, Linear …"
  without that qualifier, which overstates what the skill does out of the box.
- The "eight prompt blocks" correction: the brief says "seven near-identical
  prompt blocks." I counted the fenced blocks in README lines 69–117: there
  are eight (design, implement, review, continue, quick-architecture, compare,
  ship-e2e, github-pr-transport). codex-1 is right and the brief is stale.
- Line 371 audit: "the default does not simply use 'all discovered installed
  CLI agents.'" Verified — SKILL.md line 102 and line 283 say "default to 2–4
  active participants … at least one non-facilitator." The README line 371
  claim is materially false. I missed this.

What is wrong or unsupported:

- "The current file is 401 lines, not 402." I ran `wc -l`: it is 401 lines.
  The brief says 402. codex-1 is correct and the brief is off by one. This is
  not a claim *about* the README's content, but it is the most-cited number in
  the round and codex-1 is the only one who got it right.
- codex-1's `parley-design` entry says "PDS/1.0 makes participants diverge on
  directions, critique them, choose one whole, bind it as a contract, apply
  it, and audit the result." This is accurate but it is a list of verbs with no
  load-bearing refusal stated until the second paragraph. claude-1's entry
  states the refusal first ("selection, never averaging") and is better
  structured for a README.
- codex-1's `parley-design-check` entry says "a run that judged nothing
  reportable, or failed to verify a level claim, exits `4`, not `0`." This is
  correct per SKILL.md line 61, but the sentence is inside the catalogue entry
  and may be too much detail for the README. The exit-4 fact belongs; the
  full condition list may not.

### kimi-1

What they got right that I did not:

- The three-reader model. kimi-1's framing of the search reader, the
  colleague-sent reader, and the installed reader is the clearest analytical
  lens in the round. My round-1 file argued placement from intuition; kimi-1
  argued it from a model of who arrives and what they bounce on. The "fastest-
  path install above the catalogue, full apparatus below" split is a genuine
  third position on F1 that I did not consider.
- The scan list. Five one-clause bullets above the prose entries is the best
  compromise between "a table scans" and "prose persuades." I did not propose
  an index and claude-1's thin table is less readable than kimi-1's bulleted
  list. I graft kimi-1's scan list in §3.
- The "Reach for it when…" closing line per entry. This is the single best
  structural innovation in the round. It tells the reader *when* without
  becoming a feature list. I graft it for every entry.

What is wrong or unsupported:

- kimi-1's `parley-deck` entry says "consensus rates confidence by agreement
  and surfaces blind spots instead." This repeats the README's line 27 claim,
  which codex-1 correctly flagged as false. The shipped COOPERATION.md (lines
  324–327, 345) says the `## Comparison & blind spots` section asks for
  "contradictions, partial coverage, unique insights, and blind spots" and is
  explicitly labeled "advisory drafting discipline, not a gate." It does not
  "rate confidence by agreement." kimi-1 copied a stale README claim into the
  new copy; this must not ship.
- kimi-1's hook says "fourteen targets in all." This is correct (14 named
  runtimes). But kimi-1's F7 section says the README line 9 lists "five
  targets" — I count six in line 9 (Codex, Claude Code, Antigravity, Gemini,
  Hermes, custom skill directory). Off by one in the *accusation*, though the
  accusation itself (the intro under-sells) is correct.
- kimi-1's cut table says "Keep three blocks (design, implement, continue)"
  for the prompt section. That is a defensible position, but codex-1's "two
  blocks plus one-line substitutions" is tighter and the difference is only
  ~8 lines. I take codex-1's version in §5.

---

## 2. Fact-check: the "false claims in the current README" audit

I checked each claim against the shipped files. Line numbers are from the real
README.md (401 lines, `wc -l` verified).

### Genuinely false (must fix)

1. **Line 9 runtime list.** README line 9: "Codex, Claude Code, Antigravity,
   Gemini, Hermes, or a custom skill directory." The installer `--help`
   (verified by running `node bin/parley-deck-skill.js --help`) lists 14 named
   targets plus `generic`: codex, claude, agy, gemini, hermes, qwen, codebuddy,
   goose, kimi, droid, vibe, cursor, opencode, aionrs. The opening
   under-sells by more than half AND contradicts the README's own line 186
   which lists the full set. claude-1, codex-1, hermes-1, and kimi-1 all
   flagged this. **Verdict: genuinely false.**

2. **Line 27 "rates confidence by agreement."** README line 27: "a consensus
   'Comparison & blind spots' lens that rates confidence by agreement and
   surfaces blind spots." COOPERATION.md lines 324–327 defines the section as
   asking for "contradictions, partial coverage, unique insights, and blind
   spots." Line 345 explicitly says it is "an advisory drafting discipline,
   not a gate." There is no confidence rating. codex-1 flagged this.
   **Verdict: genuinely false.** kimi-1 reproduced this claim in new copy —
   must be removed before shipping.

3. **Line 239 version string.** README line 239: `v1.2.1`. package.json line 3:
   `"version": "1.5.0"`. All four participants flagged this.
   **Verdict: genuinely false.**

4. **Line 242 WinGet.** README line 242: "Until the WinGet manifest is
   accepted, download the `.exe`." The repo ships manifests for 1.0.4, 1.2.0,
   1.3.0, and 1.3.1 under `packaging/winget/manifests/f/Feci/ParleyDeckSkill/`.
   However, `packaging/winget/README.md` line 3 says "This directory contains
   a *draft* manifest for submitting …" and line 34 says "Open a pull request
   to `microsoft/winget-pkgs`." The manifests exist in-repo but there is no
   evidence in the repo that a PR to `microsoft/winget-pkgs` was merged.
   claude-1 claimed "a WinGet manifest was accepted; `Feci.ParleyDeckSkill` is
   published." I checked: the in-repo manifests prove the *packaging* exists,
   not that it is *published* to the community repo. **Verdict: claude-1's
   claim that the manifest "was accepted" is unsupported by the repo. The
   README's "until accepted" is likely still accurate, but the wording should
   not promise a `winget install` command unless we can verify publication.
   Treat as imprecise, not false — the line is not wrong, just stale in tone.**

5. **Line 371 "all discovered installed CLI agents."** README line 371: "By
   default it uses all discovered installed CLI agents." SKILL.md line 102:
   "default to 2-4 active participants." Line 283: "a bounded set … normally
   2-4 active participants … MUST include at least one non-facilitator." The
   README claim is materially false — the default is a bounded set, not "all."
   codex-1 flagged this. **Verdict: genuinely false.**

6. **Lines 148–166 Repository Layout.** The layout tree shows no `addons/`
   directory. Verified: `addons/` exists on disk with four subdirectories
   (`parley-design`, `parley-design-check`, `parley-tracker`,
   `parley-worktrees`). The tree also omits `test/`, `NOTICE.md`,
   `RELEASING.md`, `dist/`, `scripts/`, and `packaging/`. **Verdict: false by
   omission.** The brief already named this; claude-1 added the specific
   omissions.

### Imprecise but not false

7. **Line 22 "append-only."** README line 22: "append-only and resumable."
   codex-1 flagged this as false because `00-prompt.md` status and
   `IMPLEMENTATION.md` are updated in place. I checked COOPERATION.md: line
   684 says "consensus.md signoffs are append-only"; line 103 says "append-only
   signoffs." The signoffs are append-only; the overall lifecycle is not.
   **Verdict: imprecise.** The word "append-only" applied to the whole
   lifecycle overstates it; applied to signoffs it is correct. Fix by
   scoping the word, not deleting it.

8. **Line 119 "any capable tier-1 model."** codex-1 flagged this as an
   unsupported universal claim. I checked: no shipped file defines "tier-1" or
   gives a compatibility criterion. **Verdict: unsupported, not false.** There
   is no evidence it is wrong, but there is no evidence it is right. It should
   be cut or softened to "any model that can follow structured Markdown
   instructions."

### Wrong accusations (claims that are actually true)

9. **codex-1 line 371 audit — the "all discovered" claim.** This is a correct
   accusation, not a wrong one. Listed above as genuinely false.

10. **claude-1's "15 runtimes."** This is wrong as written (should be 14 + 1
    generic), but claude-1 self-flagged it and asked for verification. Not a
    wrong accusation — a wrong claim that was correctly self-identified.

No participant made a false accusation against the README. Every "this claim is
false" I checked was either genuinely false or imprecise. The quality of the
truth audit across the round is high.

---

## 3. The merged copy — hook and all five catalogue entries

This is selection and grafting, not averaging. For each entry I name the base
participant and the grafted lines.

### (a) The hook

Base: **kimi-1**. kimi-1's blockquote ("One agent's answer is a first draft.
Several agents' recorded agreement is a decision.") is the strongest single
line in the round — it is a framing claim, not a product claim, and it earns
the second paragraph. I graft claude-1's "files you can read, diff, and resume"
phrasing (which kimi-1 already used) and codex-1's artifact-trail sentence.

```markdown
# Parley Deck Skill

> One agent's answer is a first draft. Several agents' recorded agreement is a
> decision.

`parley-deck-skill` installs the Parley Deck cooperation protocol into your
agent runtimes — fourteen named targets plus a generic directory — so that a
design, a plan, or a code review is *worked* by several agents instead of
performed by one. Each participant writes its own files: independent round-1
analyses, cross-reviews, a gated consensus, then implementation and review as
separate phases. Every step is a file in your repository you can read, diff,
and resume — not a chat log you have to trust.

One install gives you five skills: the core protocol plus add-ons for design
doctrine, design enforcement, ticketing, and parallel worktrees.
```

Attribution:
- "One agent's answer is a first draft. Several agents' recorded agreement is a
  decision." — kimi-1, verbatim.
- "fourteen named targets plus a generic directory" — corrected from kimi-1's
  "fourteen targets in all" and claude-1's "15 agent runtimes" (which was
  wrong; see §2).
- "Each participant writes its own files … Every step is a file in your
  repository you can read, diff, and resume — not a chat log you have to
  trust." — kimi-1, verbatim, which itself grafted the "read, diff, and resume"
  phrase from claude-1.
- "One install gives you five skills" — kimi-1, verbatim.

### (b) The catalogue section header + scan list

Base: **kimi-1** for the scan list; **claude-1** for the "What's in the box"
heading and intro framing.

```markdown
## What's in the box

Installing this package installs five skills. The first is the protocol; the
other four build on it and are opt-in. All five install by default —
`--no-addons` takes just the core skill, `--only <name>` picks specific
add-ons.

- **`parley-deck`** — the multi-agent cooperation protocol.
- **`parley-design`** — collaborative design that refuses to read as machine-made.
- **`parley-design-check`** — that doctrine's rules, enforced against files on disk.
- **`parley-tracker`** — tickets a stakeholder, a reviewer, and an agent can all read.
- **`parley-worktrees`** — parallel agents over one repo, without silent corruption.
```

Attribution:
- "Installing this package installs five skills. The first is the protocol;
  the other four build on it and are opt-in." — claude-1, verbatim.
- "All five install by default — `--no-addons` takes just the core skill,
  `--only <name>` picks specific add-ons." — kimi-1, verbatim.
- The five-bullet scan list — kimi-1, verbatim.

### (c) Entry 1: `parley-deck`

Base: **kimi-1**. kimi-1's entry is the most complete (lifecycle, non-solo,
tracks, transports) and has the "Reach for it when…" close. I graft
claude-1's "write before you read" load-bearing rule, which is the one thing
kimi-1's entry lacks. I remove kimi-1's "consensus rates confidence by
agreement" (false — see §2) and replace with the accurate framing from
COOPERATION.md.

```markdown
### `parley-deck` — the core skill

Teaches your agents the 8-phase idea lifecycle: kickoff, independent round-1
analyses, cross-review, a recorded consensus, `FINAL.md`,
`IMPLEMENTATION.md`, code review, fix-up. It is non-solo by design — one
canonical file per agent per round, so no participant overwrites another —
and the load-bearing rule is **write before you read**: a round-1 file
drafted after reading the others is not an independent position, and the
protocol treats it as a failed round. Canonical files live in your
repository; GitHub PRs or GitLab MRs are mirrors, not the source of truth.
Ideas run at `fast`, `standard`, or `deliberation` rigor depending on what
is at stake. Reach for it whenever a decision, design, plan, or review is
worth more than one model's first answer.
```

Attribution:
- First sentence (lifecycle list) — kimi-1, verbatim.
- "non-solo by design — one canonical file per agent per round, so no
  participant overwrites another" — kimi-1, verbatim.
- "the load-bearing rule is write before you read: a round-1 file drafted
  after reading the others is not an independent position, and the protocol
  treats it as a failed round." — claude-1, verbatim (their "load-bearing
  rule" framing, which is the best structural innovation for the entries).
- "Canonical files live in your repository; GitHub PRs or GitLab MRs are
  mirrors, not the source of truth." — kimi-1, verbatim.
- "Reach for it whenever …" — kimi-1, verbatim.
- Removed: kimi-1's "consensus rates confidence by agreement and surfaces
  blind spots" — false per §2.

### (d) Entry 2: `parley-design`

Base: **kimi-1** for the structure (refusals as the point, budget, "Reach for
it when"). Graft **claude-1**'s "Average four directions together and you get
the average of four, which is the look everything already has" — the best
single sentence in the round.

```markdown
### `parley-design` — design doctrine with refusals built in

Ships the PDS/1.0 protocol as pure markdown — zero runtime dependencies, the
whole thing held under a 64 KiB budget that a test enforces. Its positions
are refusals, and they are the point: one direction wins whole (selection,
never averaging — average four directions together and you get the average
of four, which is the look everything already has); it emits a findings
ledger and never a numeric aesthetic score; and it declines to hand you "a
good default aesthetic," because a look guessable from the category is
exactly the failure it exists to prevent. Reach for it when a Parley idea
creates a new visual world, changes a ratified design rule, or needs an
interface audited against a contract instead of against taste.
```

Attribution:
- "Ships the PDS/1.0 protocol as pure markdown — zero runtime dependencies,
  the whole thing held under a 64 KiB budget that a test enforces." — kimi-1,
  verbatim.
- "Its positions are refusals, and they are the point" — kimi-1, verbatim.
- "one direction wins whole (selection, never averaging" — kimi-1, verbatim.
- "average four directions together and you get the average of four, which is
  the look everything already has" — claude-1, verbatim (grafted inside
  kimi-1's parenthetical).
- "it emits a findings ledger and never a numeric aesthetic score" — kimi-1,
  verbatim.
- "it declines to hand you 'a good default aesthetic,' because a look
  guessable from the category is exactly the failure it exists to prevent" —
  kimi-1, verbatim (slightly reworded from "declines to" for sentence flow).
- "Reach for it when …" — kimi-1, verbatim.

### (e) Entry 3: `parley-design-check`

Base: **kimi-1** for structure and the "Reach for it when" close. Graft
**codex-1**'s exit-code-4 fact (which kimi-1 and I both missed) and
**claude-1**'s "it would rather refuse than pass" framing for the load-bearing
rule.

```markdown
### `parley-design-check` — the doctrine, enforced by a tool

Runs the rule registry from `parley-design` against files on disk and emits
findings in one shape — `rule-id — violation — remedy` — stable across runs
and diffable in review. Node built-ins only, no network at check time. Its
load-bearing rule: **it would rather refuse than pass.** With no registry
found it refuses rule checks outright (exit 3) rather than fall back to a
bundled copy; a rule it has no detector for is reported `UNJUDGEABLE` by
name, never silently skipped; and a run that judged nothing reportable exits
4, not 0. Its capability declaration is generated by scanning its own
detectors, so it cannot claim coverage it does not have. Reach for it in a
pre-commit hook or CI step, to verify a conformance-level claim instead of
accepting it, or to get the mechanical findings into the ledger before a
review round argues about the rest.
```

Attribution:
- "Runs the rule registry from `parley-design` against files on disk and emits
  findings in one shape — `rule-id — violation — remedy` — stable across runs
  and diffable in review." — kimi-1, verbatim.
- "Node built-ins only, no network at check time." — kimi-1, verbatim.
- "it would rather refuse than pass" — claude-1, verbatim (their load-bearing-
  rule framing).
- "With no registry found it refuses rule checks outright (exit 3) rather than
  fall back to a bundled copy" — kimi-1, verbatim.
- "a run that judged nothing reportable exits 4, not 0" — codex-1, verbatim
  (the exit-4 fact codex-1 discovered and the rest of us missed).
- "Its capability declaration is generated by scanning its own detectors, so
  it cannot claim coverage it does not have." — kimi-1, verbatim.
- "Reach for it when …" — kimi-1, verbatim.

### (f) Entry 4: `parley-tracker`

Base: **kimi-1** for structure, the mirror/canonical framing, and the "Reach
for it when" close. Graft **codex-1**'s "live tracker writes require a
separate opt-in connector" (which kimi-1 and I both omitted) and
**claude-1**'s "migrate from Jira to Linear and you lose a projection, not a
requirement" — the best single sentence for this skill.

```markdown
### `parley-tracker` — tickets three audiences can read

Authors epics, stories, and subtasks as canonical markdown — one file with
`## [B] Business`, `## [T] Technical`, `## [A] Agent directives`, and a
mandatory `## At a glance` — then mirrors them into Jira, Linear, GitHub
Issues, GitLab, Trello, or a plain board. The tracker is a mirror; the file
is canonical. Sync is one-way (file → tracker) by default, `--pull` writes
back only fields the file flags `mirror-owned`, and a field the tracker
lacks is dropped from the mirror, never from the file — migrate from Jira to
Linear and you lose a projection, not a requirement. Live tracker writes
require a separate opt-in connector. A tool-enforced gap-scan refuses to let
an agent claim a ticket that still has holes. Reach for it when a
`FINAL.md` has to become a backlog — or whenever tickets must survive a
tracker migration.
```

Attribution:
- "Authors epics, stories, and subtasks as canonical markdown — one file
  with `## [B] Business`, `## [T] Technical`, `## [A] Agent directives`, and
  a mandatory `## At a glance`" — kimi-1, verbatim.
- "The tracker is a mirror; the file is canonical. Sync is one-way
  (file → tracker) by default, `--pull` writes back only fields the file flags
  `mirror-owned`, and a field the tracker lacks is dropped from the mirror,
  never from the file" — kimi-1, verbatim.
- "migrate from Jira to Linear and you lose a projection, not a requirement"
  — claude-1, verbatim (grafted as the closing example).
- "Live tracker writes require a separate opt-in connector." — codex-1,
  verbatim (the fact codex-1 discovered and the rest of us omitted).
- "A tool-enforced gap-scan refuses to let an agent claim a ticket that still
  has holes." — kimi-1, verbatim.
- "Reach for it when …" — kimi-1, verbatim.

### (g) Entry 5: `parley-worktrees`

Base: **kimi-1** (the most complete and the only one with the "Reach for it
when" close). Graft **claude-1**'s "declare your file set, and it must be
disjoint" as the load-bearing rule label, which kimi-1's entry states but
does not name as the rule.

```markdown
### `parley-worktrees` — parallel agents that don't corrupt each other

The branch + worktree + file-set discipline that turns invisible
concurrent-write corruption into an ordinary git merge conflict. The
load-bearing rule: **declare your file set, and it must be disjoint.** A
claim manifest in `IMPLEMENTATION.md` is the lock layer; the disjointness
check refuses — or demands a recorded override — when two concurrent sessions
claim intersecting files; runtime state, ports, caches, and env are isolated
per worktree; worktrees live in a sibling directory, never inside `.git/`.
Reach for it when two or more sessions or Phase-5 implementers touch the same
repository at once.
```

Attribution:
- "The branch + worktree + file-set discipline that turns invisible
  concurrent-write corruption into an ordinary git merge conflict." — kimi-1,
  verbatim.
- "declare your file set, and it must be disjoint" — claude-1, verbatim (their
  load-bearing-rule label).
- "A claim manifest in `IMPLEMENTATION.md` is the lock layer; the
  disjointness check refuses — or demands a recorded override — when two
  concurrent sessions claim intersecting files; runtime state, ports, caches,
  and env are isolated per worktree; worktrees live in a sibling directory,
  never inside `.git/`." — kimi-1, verbatim.
- "Reach for it when …" — kimi-1, verbatim.

---

## 4. The one open ordering question (F1)

kimi-1 proposes: hook → fastest-path Install (~15 lines) → catalogue → full
install apparatus → use → rest. The rest of us put the catalogue before any
install. kimi-1's argument is the colleague-sent reader: someone who got a
link in a chat wants the `npx` line immediately, and 80 lines of catalogue
before the command is friction for exactly that reader.

I think kimi-1 is right about the reader but wrong about the remedy. The
colleague-sent reader scrolls to the first `npx` line regardless of where it
is — they Ctrl-F "npx" or "install." The search reader and the installed
reader both need the catalogue before they care about the command. Putting a
15-line install block above the catalogue serves the colleague-sent reader
(who skips it) and slows the search reader (who reads top-down). The
catalogue-before-install order serves two of three readers fully and does not
harm the third.

But kimi-1's split has a real merit: the *full* install apparatus (per-target
commands, manual paths, Windows, Homebrew, updating) should go below the
catalogue. That part is not contested. So the converged order is:

hook → catalogue → fastest-path Install (~8 lines: one `npx`, one `doctor`)
→ Use (2 prompt blocks) → Why This Exists → full install apparatus →
contract → transports → relationships → status → license.

This is claude-1's and my round-1 position, with kimi-1's insight that the
install *apparatus* belongs lower. The fastest-path install is inside the
catalogue's shadow, not above it.

---

## 5. Final line budget and ordered section list

Target: **≤ 320 lines** (down from 401). The catalogue + scan list above is
~90 lines. The cuts below remove ~170.

| Section | Lines now | Target | Action |
|---|---:|---:|---|
| Title + hook | 1–17 | 16 | Replace with merged hook (§3a) |
| What the protocol gives your agents | 19–29 | 0 | Cut — the catalogue's core-skill entry does this better |
| Inspired by — adopted & adapted | 31–42 | 12 | Keep, move below relationships |
| Install (fastest path) | 44–63 | 8 | One `npx` block + one `doctor` block |
| **What's in the box** (catalogue) | — | 90 | New — the merged copy from §3 |
| Use Parley Deck | 65–119 | 30 | 2 prompt blocks (design, implement) + the tier-1 sentence, softened |
| Why This Exists | 121–131 | 11 | Keep as-is |
| What The Skill Does | 133–146 | 0 | Cut — restates the catalogue |
| Repository Layout | 148–176 | 20 | Rewrite: add `addons/` + 4 subdirs, `test/`, `NOTICE.md`, `RELEASING.md`; cut per-file prose |
| Installation Details + Installer Commands + Updating | 178–359 | 50 | Merge into one reference section; one install, one update, one per-target example, manual paths, flags |
| Local Agent Contract | 361–371 | 10 | Keep, fix line 371 ("all discovered" → "2–4 participants, at least one non-facilitator") |
| Transports | 373–381 | 9 | Keep |
| Relationship to other repos | 383–391 | 9 | Keep |
| Inspired by (moved here) | — | 12 | Moved from top |
| Status | 393–397 | 4 | Cut "should be obvious" (line 397); keep the rest |
| License | 399–401 | 3 | Keep |
| **Total** | 401 | **~284** | |

The finished README lands around 284 lines — well under the 320 target, with
the catalogue paid for by cuts and the stale claims fixed. Every factual claim
in the merged copy traces to a shipped file I read this round.
