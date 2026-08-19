---
agent: zcode-1
idea: protocol-and-skill-audit
round: 2
date: 2026-08-20
responding-to: [codex-1/round-01, zcode-1/round-01, kimi-1/round-01, claude-1/round-01]
---

# Round 2 — verification of round-01 findings (zcode-1)

Assessed: codex-1 (24), kimi-1 (5), claude-1 (3). My own round-01 findings are out of scope.
Every experiment ran in a scratch copy (`/tmp/pd-audit-zcode1/`, repo copied at HEAD 2b36075,
`parley` built from source: `parley 1.45.0`; skill package copied from
`parley-deck-skill` 2.9.0; HOME/PARLEY_HOME redirected). The shared tree was not modified.

Notable round-2 facts about round-1 method: codex-1's Go-level findings (F13–F18, F22, F23) were
backed by `TestAudit*` probes and `cmd/audit-track-probe` that codex created in its own session —
they are NOT in the tree (verified: no `TestAudit*`, `cmd/` contains only `parley`). I therefore
re-derived every one of them with my own probes (`internal/driver/zz_audit_probe_test.go`,
`internal/runner/zz_audit_probe_test.go`, scratch only). All reproduced.

Protocol sections codex cited were read against `parley-deck/COOPERATION.md` in the copy; every
quoted line range says what the findings claim, with two nuance cases (F12, F3 — see verdicts).

## Verdicts

### codex-1/F1 — blank participant files pass `consensus draft` — CONFIRMED
what I ran: fixture idea `audit-empty` (participants [alice, bob], standard) with round-01
files containing only a newline; `parley consensus draft --dir WS --round 1 --by alice audit-empty`.
what I got: `Drafted consensus at …/audit-empty/consensus.md`, `Consensus: partial`, exit 0,
and `00-prompt.md` flipped to `status: consensus`.
verdict reasoning: exact reproduction. The round gate is `missingRoundArtifacts`
(consensus.go:486-493) — pure `os.Stat` per participant. Phase 1 (COOPERATION.md:307-322) requires
frontmatter (incl. `date:`) and four sections; the manual command checks none of it, while the
auto-driver's own validator (`ValidateRoundOneArtifact`) does check shape — the asymmetry shows the
check exists and was simply not wired into the manual path. The command also advances idea status.
corrected severity: — (MAJOR stands: the CLI's only manual round gate advances an idea with zero
participant analysis).
is this a duplicate of another finding? No. (F17 is the weaker-driver-side variant, different gate.)

### codex-1/F2 — `consensus draft --review` demands the implementer's review file — CONFIRMED
what I ran: fixture `audit-review` (participants [impl, reviewer-a, reviewer-b]); protocol-compliant
`review/round-01/` containing only reviewer-a.md and reviewer-b.md;
`parley consensus draft --dir WS --review --round 1 --by reviewer-a audit-review`.
what I got: `consensus draft failed: review/round-01 is incomplete; missing impl.md`, exit 1.
verdict reasoning: exact reproduction. Phase 6 states "The implementer does not write a review-round
file" (COOPERATION.md:514) — on ANY track. `Draft` applies `idea.Participants` to the review round
(consensus.go:122). A protocol-compliant Phase 6 can never draft review consensus through the CLI.
corrected severity: — (MAJOR stands).
is this a duplicate of another finding? Same root cause as F3 (review phases reuse the full idea
participant list); one fix (reviewer-scoped participant set) closes both. Different gate — I would
merge them into one fix item, not drop either.

### codex-1/F3 — review consensus stays partial after every reviewer accepts — CONFIRMED
what I ran: after fabricating `review/round-01/impl.md` (as codex did) to get past F2: draft, then
`consensus signoff --review --agent reviewer-a --status accept`, same for reviewer-b, then
`consensus status --review audit-review`.
what I got: after both reviewers accepted: `Review consensus: partial … Missing signoffs: impl`, exit 0.
verdict reasoning: reproduced. `validateDocument` is called with `idea.Participants` for review
consensus too (consensus.go:102,167), i.e. the deliberation row ("all participants sign off") is
applied to every track. The §4.0 per-track table — which declares itself "the single authoritative
per-track gate" overriding later-phase "every participant" language for fast/standard — says
standard Phase 7 is "reviewers who reviewed sign off". Wrinkle worth recording: Phase 7's inline
comment ("Each active participant (implementer included) APPENDS") agrees with the CODE; the
protocol contradicts itself, and the table claims precedence. The auto-driver is equally affected
(`ReviewStatus` in driver_impl.go:328+ calls the same `consensus.Status`), so standard-track Phase 7
cannot close via CLI or auto-drive without a protocol-prohibited implementer review signoff.
corrected severity: — (MAJOR stands).
is this a duplicate of another finding? Twin of F2 (same root, different gate). Merge as one fix.

### codex-1/F4 — `consensus finalize --by` accepts a non-participant — CONFIRMED
what I ran: fixture `audit-final` (author alice, participants [alice, bob], consensus ready);
`parley consensus finalize --dir WS --by mallory audit-final`.
what I got: `Finalized consensus and created …/FINAL.md`, exit 0; FINAL.md frontmatter
`author: mallory`.
verdict reasoning: behavior reproduced; `Finalize` never checks `opts.By` against participants or
the 00-prompt author. But the consequence is narrower than MAJOR: the close itself is still
authorized by the signoff triage (alice+bob had accepted); only the authorship label on FINAL.md is
wrong. In a files-canonical local-dir protocol a rogue agent could hand-write FINAL.md anyway — the
CLI adds no security boundary here (unlike `signoff`, which does validate participants). A cheap
sanity check (`--by` ∈ participants ∪ {author}) is still warranted.
corrected severity: MINOR (provenance mislabeling, not an authorization bypass).
is this a duplicate of another finding? No.

### codex-1/F5 — finalize closes the idea around an empty FINAL.md scaffold — CONFIRMED
what I ran: same fixture; inspected FINAL.md and 00-prompt.md after the F4 finalize.
what I got: FINAL.md is the empty template — every section blank; the section set itself
(Goal/Scope/Implementation details/Tests/Non-goals/Verification) differs from the protocol's FINAL
template (Purpose/Context/Observable acceptance criteria/Idempotence/Known risks); `00-prompt.md`
already flipped to `status: final`.
verdict reasoning: reproduced. `Finalize` writes `finalTemplate` and immediately commits
`status: final` (consensus.go:233-237) with no content step in the manual flow. The protocol's
"sections may be N/A for trivial ideas" allowance covers emptiness of the added sections for
trivial ideas, not a wholly empty plan/spec on every idea. Telling contrast: the auto-driver
explicitly refuses scaffold FINALs and re-drafts (`finalScaffoldReason`, driver/consensus.go:50-62)
— the manual command is the unguarded path, and it permanently closes the idea.
corrected severity: — (MAJOR stands).
is this a duplicate of another finding? No (F22 is the driver-side heuristic's own limit).

### codex-1/F6 — §15.6 adversarial-alternative close condition is decorative — CONFIRMED
what I ran: fixture `audit-final`: 00-prompt says "this is a judgment call", both round-1 files
independently pick blue with no disagreement, only round-01 exists, no `## Adversarial alternative`
anywhere (grep); finalize succeeded (F4 run). Source-checked the driver: `advanceConsensus` +
`finalScaffoldReason` contain no §15.6 logic; `grep -rn '15\.6|Adversarial' internal/{driver,
consensus,runner}` (non-test) → no hits.
what I got: idea closed (`status: final`) with zero §15.6 artifacts, manual and auto paths alike.
verdict reasoning: §15.6(a) (COOPERATION.md:1339-1358) says on standard "Consensus MUST NOT close
unless at least one existing round-02 artifact contains [## Adversarial alternative]" when round 1
closes without substantive disagreement on a judgment output; §15.6(b) requires the
related-models-unanimity note in consensus.md — none of it is implemented anywhere. Caveat: the
trigger ("primarily a judgment") is semantic, so full enforcement is not mechanical; but the
standard-track clause (a) is mechanically approximable, and not even a warning exists.
corrected severity: — (MAJOR stands; weakest enforceability of the confirmed MAJORs).
is this a duplicate of another finding? No.

### codex-1/F7 — `continue` tells a fast-track idea to open the cross-review round fast skips — CONFIRMED
what I ran: hand-crafted run `audit-fast-run` (events.jsonl `run.created` → idea `audit-fast-plan`,
`track: fast`, no `cross_review_rounds` key, round-01 complete); `parley continue --dir WS audit-fast-run`.
what I got: byte-identical to codex: `Recommended: Open round-02 (cross-review) before drafting
consensus`; `Command: parley run --auto --dir . "continue audit-fast-plan"`.
verdict reasoning: `nextCrossReviewRound` (runplan.go:233-239) reads only `cross_review_rounds`
(default 1) and never `track:`. The §4.0 table (authoritative) says fast SKIPS cross-review.
corrected severity: — (MAJOR stands as a pair with F11).
is this a duplicate of another finding? Related to F8 (same track-blind planner) and F11 (same
output block); distinct claims.

### codex-1/F8 — fast track routed through ordinary consensus instead of collapsed FINAL — CONFIRMED
what I ran: second fixture `cross_review_rounds: 0` + `track: fast`; `parley continue … audit-fast-run2`.
what I got: `Recommended: Draft consensus from completed round artifacts`;
`Command: parley consensus draft --round 1 audit-fast-plan2`.
verdict reasoning: reproduced; and the gap is not just the planner: the auto-driver also has no
collapsed-final branch — `advanceRound`/`advanceConsensus` always draft ordinary `consensus.md` +
`FINAL.md` on every track (PolicyFor only sets `CrossReviewRounds: 0` for fast). The fast row says
"collapsed: one FINAL.md with embedded signoffs". Every fast-track run through the tools produces
the wrong artifact shape.
corrected severity: — (MAJOR stands).
is this a duplicate of another finding? No (F7 is the extra-round symptom of the same blindness).

### codex-1/F9 — reservation need not be logged; filler text passes finalize — CONFIRMED
what I ran: fixture `audit-reserve`: alice 🟡 with note "Rollback is unresolved and must be designed
during implementation", bob ✅. Two finalize attempts: (a) with the Open-items section empty,
(b) with `None.` as its only content.
what I got: (a) blocked — `reserved consensus requires open items deferred to implementation
before finalize`; (b) `Finalized consensus and created …/FINAL.md`, exit 0.
verdict reasoning: exact reproduction including the nuance: a gate EXISTS (`hasSectionContent`,
consensus.go:215,626-644) and blocks the empty case — but any non-comment line satisfies it,
including "None.", which directly contradicts the outstanding reservation. The protocol permits 🟡
only if the reservation is logged under open items (COOPERATION.md:384-385).
corrected severity: — (MAJOR stands, though it is the weakest of the confirmed MAJORs: the check is
a heuristic over semantic content; the mechanically detectable sub-case — 🟡 with notes +
section ∈ {None, N/A} — is the fixable core).
is this a duplicate of another finding? No.

### codex-1/F10 — manual review-consensus template violates the schema its own auto gate requires — CONFIRMED
what I ran: inspected the `review/consensus.md` produced by F2/F3's draft; read
`ValidateReviewConsensusArtifact` (phase58.go:378-386) and `ReviewStatus` (driver_impl.go:328-355).
what I got: frontmatter `cycle: 1` (protocol Phase 7 template says `review-cycle: N`), empty
`reviewed-commit:`, and no `outstanding_agreed_fixes` (grep count 0). `ReviewStatus` hard-fails
without a non-negative-integer `outstanding_agreed_fixes`.
verdict reasoning: confirmed on both halves. The manual CLI writes an artifact the auto Phase 7/8
gate rejects, and `cycle:` is the wrong key per the protocol. Nuance: `outstanding_agreed_fixes`
IS protocol-documented (COOPERATION.md:680) — just not in the Phase 7 frontmatter template, so the
protocol's own template is also incomplete.
corrected severity: — (MAJOR stands: manual+auto flows cannot interoperate at Phase 7/8).
is this a duplicate of another finding? No.

### codex-1/F11 — the printed "Open round-02" command creates a new idea at round 1 — CONFIRMED
what I ran: reproduced the printed command (F7 output). Verified the execution semantics from
source: `runTask` (app.go:1761+) has no "continue"-task special case anywhere; `runcontrol.Create`
→ `CreateIdeaFull` → `uniqueSlug(timestampedSlug(task, now))` unconditionally creates a NEW
timestamped idea from the task text at round-01; `parley run`'s own help says "Create a new idea
from TASK and start round-01"; the correct advance command is `parley continue --auto RUN`
(continueAuto, app.go:1170+).
what I got: `Command: parley run --auto --dir . "continue audit-fast-plan"` printed verbatim; the
run path provably treats that whole string as a new idea's task.
verdict reasoning: I did not execute the printed command (it would launch agents); the creation
semantics are established from PRIMARY source reading of the exact call chain, plus the command's
own help text. The comment in runaction/action.go:50-52 ("the surfaced command re-runs the task to
advance the idea") is simply false. Following the CLI's recommendation forks an unrelated idea and
leaves the original stalled.
corrected severity: — (MAJOR stands).
is this a duplicate of another finding? No (F7 is the wrong recommendation; F11 is the broken
command printed for it).

### codex-1/F12 — `init` reports completed bootstrap without the mandatory confirmation — PARTIAL
what I ran: `PARLEY_HOME=… parley init --dir WS` on fresh dirs; checked for deck agents.toml, any
prompting, output text.
what I got: exit 0, `Initialized Parley Deck workspace at …`, no roster/model/effort interaction,
no `parley-deck/agents.toml` created (only the central `~/.parley/agents.toml`).
verdict reasoning: the behavior reproduces, but the finding mis-attributes the duty. COOPERATION.md:57
assigns the confirmation to "the facilitator" as a step at deck creation and says "See the skill for
the interactive list-roster → confirm → list-models-and-effort → pick flow"; the skill's "Deck
bootstrap" section (SKILL.md:166+) makes it the facilitating agent's mandatory step "before the
first idea". `parley init` is the structure scaffolder, never claims the bootstrap confirmation
happened ("Initialized … workspace" ≠ "bootstrap complete"), and the un-bootstrapped state is
detectable (no deck agents.toml). The real gap is UX: init prints no pointer to the pending
mandatory gate, so a CLI-only consumer can believe setup is done.
corrected severity: MINOR (missing next-step surface, not a skipped mandatory CLI gate).
is this a duplicate of another finding? Related in theme to claude-1/F1 (skill coverage of
bootstrap/readiness duties); different claims.

### codex-1/F13 — driver accepts a standard+auto_implement track its classifier rejects — CONFIRMED
what I ran: `go run ./cmd/parley classify --auto-implement --declared standard` (and `--json`);
my in-package driver probe constructing `driver.New` over a `track: standard` + `auto_implement:
true` 00-prompt with 4 participants.
what I got: classify → `deliberation`, `declared track "standard" is under-tiered; the classifier
floor is "deliberation" (auto_implement)`, `valid: false` (exit 4). Driver probe →
`trackErr=<nil>` with standard caps applied (reviewers=2 cross=2 cap=2 fixups=2).
verdict reasoning: both surfaces are PRIMARY. `track.Classify` is called only by the advisory
`classify` command (grep: sole call site app/classify.go:61); `PolicyFor` rejects only fast +
auto_implement/strict_gate (track.go:143-148) and waves explicit standard + auto_implement through
into the standard row. §4.0's ordering is "normative and fail-safe" and down-tiering below the
floor "requires a recorded user OK" — a run is executed under the smaller caps the project's own
classifier calls unsafe.
corrected severity: — (MAJOR stands).
is this a duplicate of another finding? No.

### codex-1/F14 — omitted `track:` does not apply the documented standard default — CONFIRMED
what I ran: my driver probe with otherwise-identical 4-participant configs, `CrossReviewRounds: 9`
requested, varying only absent vs explicit `track: standard`.
what I got: `absent: track="standard" reviewers=0 min=2 cross=9 cap=0 fixups=3` vs
`explicit-standard: track="standard" reviewers=2 min=2 cross=2 cap=2 fixups=2`.
verdict reasoning: numbers reproduce exactly. The protocol says omitted track defaults to standard
and the table is binding; the code deliberately preserves legacy behaviour (track.go:130-135
comment; driver.go:119-124 comment). Two points temper severity: (1) the error direction is
uniformly MORE ceremony, not less — unlimited reviewers, 9 rounds, 3 fix-ups are all ≥ the standard
row, so this is loop-budget/escalation-valve drift, not under-review; (2) it is an acknowledged
backward-compat decision that the protocol text simply does not document. The printed-vs-actual
drift is real (the most common case — no explicit track — binds none of the printed caps, and
`advanceRound` promotion consults no MaxRounds), but the consequence is inefficiency.
corrected severity: MINOR (policy/documentation decision; fail-direction is stricter, not looser).
is this a duplicate of another finding? F15 is the same mechanism's typo variant — merge them.

### codex-1/F15 — invalid explicit track silently disables every standard cap — CONFIRMED
what I ran: my driver probe with `track: standart`.
what I got: `unknown-standart: track="standard" reviewers=0 min=2 cross=9 cap=0 fixups=3
trackErr=<nil>` — identical to the absent case.
verdict reasoning: `Normalize` returns (Standard, false) for unknown values (track.go:36-48), so a
typo is silently relabeled and treated as absent-legacy, no error. §4.0's fail-safe rule says doubt
fails closed to the stricter track; a garbage value is doubt. Same severity tempering as F14.
corrected severity: MINOR.
is this a duplicate of another finding? Yes — merge with F14 (one normalization/absence policy fix).

### codex-1/F16 — cross-review gate accepts an empty `responding-to:` field — CONFIRMED
what I ran: my driver probe: round-02 artifacts with `responding-to:` (empty value) plus all
required `### @<other>` headings; `d.roundComplete(2)`.
what I got: `done=true err=<nil>`.
verdict reasoning: reproduced; `hasRespondingTo` (driver.go:466-469) checks key presence only —
`readFrontmatterField` returns ok=true for an empty value. However the claim's severity leans on
lost provenance, and the OTHER half of the same gate, `validateCrossReviewBody`, still requires a
`### @<other>` heading for every other participant — the substantive addressing evidence. The
unchecked part is the machine-readable mirror of it.
corrected severity: MINOR (frontmatter provenance metadata unchecked; the per-agent heading gate
holds).
is this a duplicate of another finding? No.

### codex-1/F17 — driver completes round 1 when every required section is empty — CONFIRMED
what I ran: my driver probe: round-01 artifacts with correct agent/idea/round frontmatter, all four
headings, no `date:`, zero body bytes; `roundComplete(1)` then one `Advance` tick.
what I got: `done=true err=<nil>`; `advance action="consensus-ready"` (stops only because my probe
left Consensus ops unwired; production wiring drafts).
verdict reasoning: reproduced. `ValidateRoundOneArtifact` (validation.go:42-72) requires identity
frontmatter + heading SUBSTRINGS — no date, no content. The driver reconstructs a terminal
`round.completed` event from these artifacts (driver.go:385-398), so header-only scaffolds are
"independent analyses" and auto-advance proceeds.
corrected severity: — (MAJOR stands).
is this a duplicate of another finding? No (F1 is the manual CLI's even-weaker existence-only gate).

### codex-1/F18 — review with no `reviewed-commit` passes the Phase 6 validator — CONFIRMED
what I ran: my probe calling the real `runner.ValidateReviewArtifact` on a review artifact with
identity frontmatter + non-empty Refutation attempts + Findings, but no `reviewed-commit`, `date`,
Summary, or Open questions.
what I got: `err=<nil>`.
verdict reasoning: reproduced; the validator (phase58.go:413-442) checks identity frontmatter,
`## Findings`, and a non-empty `## Refutation attempts` only. The protocol's Phase 6/7 templates
carry `reviewed-commit: <sha>` and nothing validates it anywhere (the review-consensus validator
checks only `outstanding_agreed_fixes`). The stale-review consequence is real for the fix-up loop:
nothing ties a re-review to the post-fix commit, so a pre-fix review can satisfy the close gate.
corrected severity: — (MAJOR stands).
is this a duplicate of another finding? No.

### codex-1/F19 — init leaves workspace identity and creation date as placeholders — CONFIRMED
what I ran: `parley init` on fresh workspace; inspected COOPERATION.md header.
what I got: `**Workspace:** `<workspace-name>`` and `**Created:** `<date> — created by parley
init>`` remain unsubstituted; only `**Transport:**` is replaced (workspace.go:94-105 substitutes
exactly one line).
verdict reasoning: reproduced; the embedded template's own text ("created by parley init") implies
init should set it.
corrected severity: — (MINOR as filed).
is this a duplicate of another finding? No.

### codex-1/F20 — signoff-looking headings outside `## Signoffs` satisfy the consensus gate — CONFIRMED
what I ran: fixture `audit-signoff-scope`: two `### Signoff:` blocks under `## Agreed decisions`
(quoted-example framing), real `## Signoffs` section containing only the comment;
`parley consensus status`.
what I got: `Consensus: ready`, both signoffs counted, exit 0.
verdict reasoning: reproduced; `parseDocument` (consensus.go:331-368) scans the whole file with no
section tracking. This is not purely theoretical: BLOCK→reopen→redraft cycles explicitly seed the
new draft from the prior consensus, so quoting prior signoff blocks into `## Agreed decisions` is
a plausible accident that manufactures a false-ready.
corrected severity: — (MAJOR stands; this is the consensus gate itself).
is this a duplicate of another finding? No.

### codex-1/F21 — consensus status ignores a frontmatter slug naming a different idea — CONFIRMED
what I ran: same fixture with `idea: unrelated-idea` in consensus.md frontmatter;
`parley consensus status … audit-signoff-scope`.
what I got: `Consensus: ready`, no warning or error.
verdict reasoning: reproduced; parsing never reads frontmatter. But the precondition is itself an
operator/agent error (misplaced or hand-copied file) — the CLI's own draft always writes the
correct slug, and nothing in a normal flow produces a mismatch. The check is cheap and worth
adding; the consequence is conditional.
corrected severity: MINOR.
is this a duplicate of another finding? No.

### codex-1/F22 — driver FINAL gate accepts three arbitrary padded lines — CONFIRMED
what I ran: my probe calling the real `finalScaffoldReason` on a FINAL.md with `status: final`,
wrong idea slug, `## Final plan / specification`, two one-word lines plus a 260-byte padding line.
what I got: `finalScaffoldReason=""` (accepted).
verdict reasoning: reproduced; the gate checks >250 bytes, status, one heading, no placeholder
tokens, ≥3 content lines (driver/consensus.go:164-207). It ignores the idea slug and every
protocol-required section (Purpose/Context/criteria/Idempotence/Risks — which may be "N/A" for
trivial ideas but must presumably still exist). Temper: this function's documented design purpose
(consensus D7/AF1) is rejecting failed-draft scaffolds behind a real agent drafter, not validating
specification schema; a lazy drafter writing 3 real lines SHOULD pass. The actionable fix is
checking the required headings exist.
corrected severity: MINOR (heuristic backstop by design; heading-presence check is the fix).
is this a duplicate of another finding? No (F5 is the manual finalize path, which has no gate at all).

### codex-1/F23 — implementation gate accepts an unknown status and empty artifact — CONFIRMED
what I ran: my probe calling the real `runner.ValidateImplementationArtifact` on an
IMPLEMENTATION.md with `status: banana`, an empty `## Summary of work`, and no implementer/dates/
branch/head-commit/plan/deviations.
what I got: `err=<nil>`.
verdict reasoning: reproduced; the validator (phase58.go:390-410) checks the idea slug, ANY
non-empty status, and a Summary substring. The protocol documents a closed status vocabulary
(`implemented | fix-up-cycle-N | complete`, COOPERATION.md:445). The validator is the production
backstop on the implementer's output (phase58.go:313) and advanceImpl then opens review. Temper:
the close decision itself is carried by other gates (checks contract, strict gate, review loop),
reviewers read the actual IMPLEMENTATION.md, and status transitions are rewritten by the driver at
completion — so the bogus status mostly wastes a review round rather than closing garbage.
corrected severity: MINOR.
is this a duplicate of another finding? No.

### codex-1/F24 — fresh-deck preflight calls an altered protocol "in sync" — CONFIRMED
what I ran: fresh `parley init` workspace (version.json = protocolRole/deckVersion/created only);
genuinely mutated a body invariant of COOPERATION.md (first attempt's sed matched nothing and was
discarded — the confirmed run used a verified mutation); `parley preflight --no-ping`.
what I got: `Freshness: consumer — protocol matches packaged skill (in sync)`,
`classification=in-sync`, `Ready: no pending gates.`, exit 0.
verdict reasoning: reproduced. The consumer path compares only the two metadata hashes
(preflight.go:418-421); a fresh init writes neither (`"" == ""` → in-sync), and no code hashes the
live file on this path. The only writer of those hashes is the skill package's `sync-project`
(see kimi-1/F2) — which simultaneously deletes `protocolRole`. So the §9.0 freshness gate is
guaranteed fail-open from `parley init` until a skill sync runs.
corrected severity: — (MAJOR stands).
is this a duplicate of another finding? No, but it interlocks with kimi-1/F2 (see Duplicates).

### kimi-1/F1 — `doctor` does not byte-verify the managed core skill — CONFIRMED
what I ran: scratch HOME, `install --target kimi` (6 units, exit 0); length-preserving single-byte
flip at the midpoint of the managed core's `references/COOPERATION.md`; `doctor --target kimi`.
Controls: same flip on managed add-on `parley-design/SKILL.md`.
what I got: core flipped → `kimi/parley-deck: valid`, exit 0 (clean rc, no pipe). Add-on flipped →
`kimi/parley-design: malformed`, `integrity: modified: SKILL.md`, exit 1 (after reinstall-restore
of the core, both states observed). Installed core has NO `parley-addon.json`; the add-on HAS one;
the package ships `skills/parley-deck/parley-addon.json`; `PAYLOAD_ENTRIES` never installs it and
`manifestProblems` runs only for `kind === "addon"` (installer.js:148-154, 2397). README:79-82
says "**Every** packaged skill ships a parley-addon.json integrity manifest, so doctor can verify
any of them byte for byte".
verdict reasoning: core claim, asymmetry control, root cause, and README quote all verified
PRIMARY. The one control I could NOT independently reproduce is the foreign-copy (`valid-unmanaged`
→ malformed on flip): my attempts to construct a markerless tree the doctor enumerates failed
(`--dest` requires `--target generic`; `--only parley-deck` rejected) — that control stays
SECONDARY on kimi's own run and does not affect the verdict.
corrected severity: — (MAJOR stands: the health gate is green for tampering in the one file every
installed agent actually reads).
is this a duplicate of another finding? No.

### kimi-1/F2 — `sync-project --yes` deletes `protocolRole` from meta/version.json — CONFIRMED
what I ran: sandbox project with COOPERATION.md + version.json containing
`"protocolRole": "source"` (plus protocolSha256/packagedProtocolSha256/created);
`sync-project --project … --yes`; also `grep -rc protocolRole` across the skill package.
what I got: exit 0, `wrote …/version.json`; the new file has the 11-key shape with
protocolSha256/packagedProtocolSha256 populated and `protocolRole` (and `created`) GONE. Zero
`protocolRole` references in the skill package.
verdict reasoning: reproduced. §9.0 (COOPERATION.md:835-838) keys behavior on `protocolRole` and
says missing/unknown → "do not auto-write; ask the user once and backfill" — the tool's own
recommended refresh (`status` prints exactly that sync-project command on this very repo) erases
the confirmed role and re-raises the gate, or with `--yes` backfills `consumer`, which is the
wrong role on a source repo. Two packages write incompatible shapes for one file.
corrected severity: — (MAJOR stands).
is this a duplicate of another finding? No — but it is the other half of codex-1/F24: sync-project
is the only populator of the freshness hashes, and it simultaneously breaks the role field.

### kimi-1/F3 — README says "fourteen" and omits `zcode`; installer ships 15 — CONFIRMED
what I ran: `node -e require(installer).TARGETS` (15: …,aionrs,zcode); `grep -n fourteen README.md`;
README --target enumeration line; CHANGELOG 2.9.0; the CLI usage text.
what I got: 15 targets; "fourteen" at README lines 75, 161, 206, 262; enumeration
`…|opencode|aionrs|generic` with no `zcode`; CHANGELOG: "bringing the target count to 15"; CLI
usage DOES include zcode.
verdict reasoning: fully reproduced; only the README lagged.
corrected severity: — (MINOR as filed).
is this a duplicate of another finding? No.

### kimi-1/F4 — `status` recommends adopting packaged updates on the source repo — CONFIRMED
what I ran: live read-only `status` and `status --json` on this repo (the protocol's upstream).
what I got: `compatibility: warning`; action line verbatim: `Review the local COOPERATION.md
changes before adopting packaged protocol updates.`; live protocolSha256 74c8470b… ≠ packaged
254521eb… (live is newer). Also noted: this repo's own version.json currently has NO protocolRole
field (consistent with F2's erasure), so "source role" is semantic here, not recorded.
verdict reasoning: reproduced; the message direction (adopt the older packaged copy onto the newer
live protocol) is backwards for the upstream repo, and the skill has no role concept to gate it
(same root as F2). `status` itself writes nothing, hence MINOR.
corrected severity: — (MINOR as filed).
is this a duplicate of another finding? Symptom of kimi-1/F2's root cause (no protocolRole concept
in the skill package).

### kimi-1/F5 — schema-2 core marker omits `manifest`, violating the installer's own invariant — CONFIRMED
what I ran: `cat` of the core and add-on install markers in my scratch HOME; installer.js:14-16.
what I got: core marker `markerSchema: 2, skill: parley-deck, addon: false`, no `manifest` key;
add-on marker same schema WITH `manifest`; the comment says "a marker at this schema that omits it
is malformed, never treated as legacy", enforced only via `kind === "addon"`.
verdict reasoning: reproduced exactly; a comment-stated universal invariant its own writer
violates, no runtime effect today.
corrected severity: — (NIT as filed).
is this a duplicate of another finding? Shares its root surface with kimi-1/F1 (core vs add-on
asymmetry in the same installer); different claim.

### claude-1/F1 — SKILL.md never names preflight/retro/loop tick/consult/repo-map/tui — PARTIAL
what I ran: the same grep loop over `parley-deck-skill/skills/parley-deck/SKILL.md`; read
SKILL.md's structure, its references to COOPERATION.md, and which `parley <cmd>` words it does name.
what I got: all six counts 0 (exact reproduction). But: SKILL.md's first Core Rule (line 12) is
"**Always read `parley-deck/COOPERATION.md` first**"; lines 26-27 make the bundled
`references/COOPERATION.md` the fallback; line 34 says "`SKILL.md` **and
`references/COOPERATION.md`** are the vendor-neutral instructions for all agents" — claude's quote
("the skill's own text calls it 'the vendor-neutral instructions…'") elides the second half of that
sentence. COOPERATION.md names `parley preflight`/`retro`/`loop tick`, i.e. the commands are one
directed hop away, inside the same skill package. SKILL.md DOES name some commands itself
(init/roster/agents/protocol).
verdict reasoning: the counts are real, but the finding's premise — an agent that "reads only the
skill" sees obligations with no commands — contradicts the skill's own text, which defines the
instruction set as SKILL.md plus the bundled protocol and mandates reading the protocol first. The
surviving substance is narrower: SKILL.md's self-contained "Protocol Coverage Checklist"
paraphrases ping/deadline and phase obligations without naming the CLI commands (preflight, retro,
consult, loop tick), so an agent working from the checklist alone misses them.
corrected severity: MINOR (coverage-checklist completeness gap, not "no way to discharge").
is this a duplicate of another finding? Thematically adjacent to codex-1/F12 (bootstrap duties
surface); different claims.

### claude-1/F2 — the protocol's own bootstrap instruction fails the repo's drift guard — CONFIRMED
what I ran: in the repo copy: baseline `go test ./internal/protocol -run TestEmbeddedDefaultMatchesLiveDeck`
(PASS); `parley roster render --dir <copy> --yes --adopt-inherited` (exit 0, "Regenerated §2");
re-ran the test.
what I got: after render, §2 header is the four-column `| Agent ID | Workspace dir | Role | State |`;
test FAILS: `live deck: anchor "| Agent ID       | Workspace dir                       | Role          |"
appears 0 times, want exactly 1 (drift guard fails closed)`.
verdict reasoning: fully reproduced, including the failure being the fail-closed kind (loud, not
silent). COOPERATION.md:57 mandates this exact render at deck bootstrap; inside this repository the
mandated command breaks the build. The renderer's output shape and the embedded default's §2 shape
disagree — i.e. even a fresh `parley init` deck holds a §2 the mandated render immediately rewrites
into a different shape.
corrected severity: — (MAJOR stands for this repo's workflow).
is this a duplicate of another finding? No.

### claude-1/F3 — `masked-by-env` is documented in the closed STATUS vocabulary but never emitted — CONFIRMED
what I ran: `grep -rn masked-by-env internal/ | grep -v _test`; SKILL.md vocabulary line; the
`addStatus(` call sites.
what I got: only roster_set.go:83-88 — a comment acknowledging the history and literal advice text
in a printf; SKILL.md:284 lists it in the STATUS vocabulary; addStatus emits legacy-roster,
inherited-roster, inactive, unmapped, not-installed, model-drift, … but never masked-by-env.
verdict reasoning: reproduced. Doc-vs-behavior mismatch; the comment shows it was noticed and
half-resolved (advice text kept, vocabulary not corrected). The fix inverts to documentation.
corrected severity: — (MINOR as filed).
is this a duplicate of another finding? No.

## Findings I could not assess, and why

None — all 32 findings were assessed. Two sub-claims within assessed findings rest on less than
full primary reproduction and are flagged in their verdicts:
- codex-1/F11: the printed command's new-idea semantics were established from PRIMARY source
  reading of the exact call chain (runTask → runcontrol.Create → CreateIdeaFull →
  timestampedSlug) plus `run`'s own help text; I did not execute the command because doing so
  would launch hosted agents.
- kimi-1/F1's foreign-copy control (valid-unmanaged → malformed on flip): my attempts to build a
  markerless tree that doctor enumerates failed on CLI constraints (`--dest` requires
  `--target generic`; `--only` rejects the core). The control remains SECONDARY on kimi's run;
  the finding's core claim was independently confirmed without it.

## Duplicates I found across authors

- codex-1/F2 + codex-1/F3 — one root cause: review phases reuse the full idea participant list
  (artifact-existence gate AND signoff quorum). One reviewer-scoped fix closes both; keep as one
  fix item with two symptoms.
- codex-1/F14 + codex-1/F15 — one mechanism: absent ≡ unknown track → legacy behaviour labeled
  standard (track.Normalize returns present=false for garbage). Merge.
- codex-1/F7 + codex-1/F8 — same track-blind continuation planner; distinct manifestations
  (wrong extra round vs wrong consensus shape). Fix together but both are real.
- codex-1/F24 + kimi-1/F2 — complementary halves of one version.json ownership conflict: the CLI
  writes {protocolRole, deckVersion, created} with no hashes (freshness fails open), the skill's
  sync-project writes 11 keys with hashes but no protocolRole (role gate re-raised / wrongly
  backfilled). Nobody writes a correct shape; fix as one item.
- kimi-1/F2 + kimi-1/F4 (+ the README claim in kimi-1/F1's neighborhood) — F4's backwards advice
  and F2's field erasure share the root "skill package has no protocolRole concept"; F2 is the
  substantive defect, F4 a message symptom.
- codex-1/F5 vs codex-1/F22 — NOT duplicates (manual finalize has no content gate; the driver's
  FINAL gate is a heuristic backstop); noted because they read similarly.
- claude-1/F1 vs codex-1/F12 — thematically adjacent (where bootstrap/readiness duties are
  surfaced), different defects; no merge.

## Anything round 1 missed that these findings led me to

1. The manual `consensus finalize` template's section set (Goal/Scope/Implementation details/
   Tests/Non-goals/Verification) is not the protocol's FINAL template (Purpose / user-visible
   outcome, Context & orientation, Observable acceptance criteria, Idempotence & recovery, Known
   risks / de-risking, References). Even a diligent drafter who fills the CLI's scaffold produces
   a FINAL.md missing every protocol-required heading — and the driver's own FINAL gate checks
   only one of them. (Found while verifying codex F5/F22.)
2. `Draft` also silently omits the protocol's advisory `## Comparison & blind spots` section from
   the consensus template (COOPERATION.md:370-373) — minor, same template-drift family as F10.
   (Found while verifying codex F1/F10.)
3. The manual review-consensus template writes an EMPTY `reviewed-commit:` even though DraftOptions
   carries a ReviewedCommit field — the flag exists but nothing surfaced it in my run; combined
   with F18 nothing in the review chain records the reviewed revision end-to-end. (Found while
   verifying codex F10/F18.)
4. `advanceRound` promotion consults no MaxRounds bound (only the consensus-BLOCK path does,
   driver.go:100); with an absent track and a large `cross_review_rounds:` the driver will run
   1+N rounds with no circuit breaker beyond the LE-5 ceilings (which default to unlimited when
   unconfigured). Aggravates codex F14. (Found while verifying F14.)
5. `parley roster render` output shape (4-column compact) vs embedded default §2 (3-column padded)
   disagree everywhere, not just under the drift guard: every fresh `parley init` deck holds a §2
   that the first mandated render rewrites into a different shape — the drift guard only catches
   it in repos that run the Go tests. Extends claude-1/F2. (Found while verifying it.)
6. `parley continue`'s review/impl phases: the planner has no review-phase or implementation-phase
   actions at all (kinds end at finalize) — `continue` on a run whose idea is past FINAL surfaces
   nothing useful. Not filed by round 1; observed while hand-crafting run fixtures.
