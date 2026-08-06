### kimi-1 — revision 2

**Verdict:** accept with reservations — two non-blocking reservations (R2, R3), both confined to
the VC-1 measurement appendix. Every required change from revision 1 is met, and the six
decisions stand.

#### codex-1's three required changes — all met

**(i) Claim (b) marked `UNVERIFIED` and removed from the closure — MET, and now exceeded.**
The draft demotes the 23-deck aggregate to illustration, marks it `UNVERIFIED` by a non-owner,
and closes VC-1 on the design argument (`consensus.md:209-213`, `:241-247`; position change 7,
`:278`). That satisfies the condition as codex-1 stated it. Additionally, I have now verified
the aggregate myself as a non-owner (owner: claude-1) — **CONFIRMED [PRIMARY]**:

- `find "/Volumes/My Shared Files/AI_WORKSPACE" -maxdepth 5 -name headless-agents.local.json
  -path "*/parley-deck/meta/*" | sort` → **23 files** (2m26s on this volume; an unbounded
  re-run of the drafter's exact find returned the same 23).
- A `while IFS= read -r f; do grep -l writeModeArgs "$f"; done` loop over those 23 → **12
  files**, the same 12 the draft names (`consensus.md:201-204`), BYTE and
  `ecb-meeting-2026.05` included.
- A `python3` JSON pass comparing each agent entry's `writeModeArgs` against its
  `headlessArgs` → **10 non-empty, 9 exclusive, 2 empty**; the empty two are exactly `BYTE`
  and `ecb-meeting-2026.05`; the exclusive 9 match the draft's list; this repo's own deck is
  the one non-exclusive case (flags duplicated in `headlessArgs`).
- `igm-app` worked example spot-checked verbatim: claude-1
  `["--permission-mode","bypassPermissions","--add-dir","parley-deck"]`, codex-1
  `["--sandbox","workspace-write","-c","approval_policy=\"never\""]`, hermes-1 `["--yolo"]`,
  none present in the respective `headlessArgs` (`consensus.md:205-207`).

So (b) enters the record doubly: still correctly excluded from the closure per codex-1's
condition, and now also non-owner-`PRIMARY` on the numbers.

**(ii) Claim (a) narrowed, three inferences deleted — MET.** The scoped claim is now "no Go
code declares a `writeModeArgs`/`write_mode_args` field or reads
`meta/headless-agents.local.json`" (`consensus.md:153-156`); "no consumer", "zero behaviour
change" and "the protected cost does not exist" are explicitly withdrawn and the manual
facilitator is acknowledged as the consumer (`:175-179`; position change 6, `:277`). Nothing
of the revision-1 overreach survives in the closure.

**(iii) Legacy migration rule — MET.** Adopted verbatim in substance at `consensus.md:254-262`:
when an existing `headless-agents.local.json` contains `writeModeArgs`, merge its arguments
into that agent's `headlessArgs` and remove the field; no deck is edited by this idea; the
1.39.0 `AUTO=no` surfacing is named. This is codex-1's condition as filed.

#### kimi-1's reservation R1 — correctly recorded

MET. Decision 4's revision-2 correction (`consensus.md:83-96`) states that `prepack` runs only
`node scripts/build-addon-manifest.js --check`, that a test-file assertion runs on `npm test`
only, and that a direct `npm publish` would bypass it; implementation MUST close the half
explicitly via option (a) or (b), and "shipping only the test-file assertion does **not**
satisfy this decision". I verified the underlying facts **[PRIMARY]** against
`parley-deck-skill/package.json`: `"test": "node --test && node scripts/run-python-tests.js &&
node scripts/build-addon-manifest.js --check"` and `"prepack": "node
scripts/build-addon-manifest.js --check"`. The draft's `node --test && …` ellipsis
(`consensus.md:85`) is accurate. Position-change entry 9 (`:280`) records the correction with
the right attribution (hermes-1 raised it, kimi-1 restated it as R1, drafter verified).

#### The defective-tool disclosure — assessment

**The corrected measurement is now adequately evidenced.** Every load-bearing number has been
re-derived with `find`-based tooling and, as of this signoff, carries a non-owner `PRIMARY`
(my reproduction above). Claim (a) does not rest on the suspect tool at all: the revision-2
re-run uses `find … -print0 | xargs -0 grep` over all 1663 `*.go` files (`consensus.md:158-168`),
and codex-1, hermes-1 and I each confirmed it independently with `rg`, which shows no defect
(`:170-173`).

**(a1) Zero Go references to `writeModeArgs`/`write_mode_args`/`WriteModeArgs` — CONFIRMED
[PRIMARY], fresh run:** `rg -l -g '*.go'
'headless-agents|headless_agents|writeModeArgs|write_mode_args|WriteModeArgs' .` from the
repository root exited 1 with no output; `internal/config/runtime.go:97-132` (`agentOverride`)
has `HeadlessArgs` and no write-mode field. **(a2) The CLI never reads the JSON file — no
verdict from me; I co-own it** (`round-01/kimi-1.md:99-103`, §15.1). Raw evidence only: the
same `rg` run covered `headless-agents|headless_agents` (zero hits), and
`internal/config/runtime.go:134-154` (`configLayers`) enumerates `~/.parley/agents.toml`
(`:141`), `parley-deck/agents.toml` (`:144`), `parley-deck/agents.local.toml` (`:145`),
`$PARLEY_HEADLESS_AGENT_CONFIG` (`:147-151`) — no JSON path. The standing non-owner verdicts
are codex-1's and hermes-1's `PRIMARY`s.

**However, the tooling-defect claim itself I could not reproduce — Reservation R3
(non-blocking).** The drafter claims `grep -r` saw 172 of 1663 `*.go` files and missed
`BYTE/parley-deck/meta/headless-agents.local.json` entirely (`consensus.md:215-221`). My runs
**[PRIMARY]**, same volume, minutes before this signoff: `grep -rl '' --include='*.go' .` in
this repository returned **1663** files — identical to `find . -name '*.go' -type f | wc -l`
(1663) — and `grep -rl writeModeArgs` over the BYTE tree **did** return the meta file (4
matches on a direct grep). As a non-owner I therefore verdict the two specific defect
measurements **WRONG / not reproduced**; I concede the behaviour may be intermittent or
mount-state-dependent — the volume is demonstrably slow (an unbounded `find` over
`AI_WORKSPACE` took minutes) — but I have no evidence of under-reporting, and my revision-1
signoff's numbers were gathered with `rg`, not `grep -r`. Per §15.3 this is a contradictory
verdict: the defect claim is `DISPUTED` unless resolved. It does not block, because **no
decision or acceptance criterion depends on it**: claim (a) stands on `rg`-based non-owner
evidence, claim (b) stands on my `find`-based reproduction, and revision-1's undercount (11/8
vs the correct 12/9) is a fact whatever its cause. `FINAL.md` MUST record this dependency
check and this dissent (`§15.3`). The draft's consequence 1 — re-derive with
`find … -print0`, distrust `grep -r` on this volume — is sound practice regardless of the
verdict on the explanation.

**Reservation R2 (non-blocking) — the published claim-(b) command is itself malformed on this
volume.** The appendix prints `find "/Volumes/My Shared Files/AI_WORKSPACE" … | sort` then
`… | xargs grep -l writeModeArgs` (`consensus.md:185-189`). The volume path contains spaces,
and `xargs` without `-0` splits on blanks: my run of exactly that pipeline shape returned
**0** matching files, not 12 **[PRIMARY]**. The *numbers* are correct — I reproduced 12/10/9
with a `while read` loop, quoted above — but the *printed command* cannot have produced them
and will mislead the next measurer, which is precisely the failure mode the tooling-defect
section exists to prevent. Fix the record to `find … -print0 | xargs -0 grep -l
writeModeArgs`. Non-blocking because (b) is illustration-only and a working non-owner command
is now on the record in this signoff.

**Does anything else rest on the defective tool? Checked — no.** (a) was re-run `find`-based
and confirmed with `rg` (above). The decision-4 gap claim "nothing under `test/` or `scripts/`
reads `skillVersion`" (`consensus.md:102-105`) — re-run by me with `rg -n skillVersion test/
scripts/` in `parley-deck-skill`: zero hits, rc=1; `skillVersion` occurs in the whole skill
repo only in `skills/parley-deck/references/compatibility.json` itself **CONFIRMED [PRIMARY]**.
The `references/COOPERATION.md` currency claim is `git diff`-based, not `grep`-based;
facilitator-owned; I did not re-run it and nothing in this signoff depends on it (as in
revision 1).

#### §15.5 drafter position changes — ratified

Entries 6-9 match what I can check **[PRIMARY on the comparisons]**: 6 and 7 against the
embedded revision-1 text (`consensus.md:325-338`) versus the rewritten `:153-179` and
`:241-247`; 8 against my own reproduction (11/8 → 12/10/9 confirmed as the correct figures);
9 against `package.json` (above). Entries 1-5 were ratified in my revision-1 signoff and are
unchanged. The table's self-assessment — "the drafter's revision-1 text did not survive
review" (`:282-285`) — is accurate.

#### §15.6

Unchanged from revision 1; my ratification stands [PRIMARY against the round files read for
revision 1]: round 1 was not unanimous, so no steelman section was required; the
correlated-agreement caveat, the what-would-have-to-be-true paragraph and the one-family note
are present at `consensus.md:287-307`.

I authored the adopted `:251` text and the opencode row wording; per §15.1 I issue no verdict
on their quality. This signoff is a signoff decision, not a self-verification.

#### Scope declaration

Read in full for this signoff: `COOPERATION.md` §15 (lines 1176-1316); `consensus.md` revision
2 (all 548 lines, including the four embedded revision-1 signoffs); my own revision-1 signoff
(embedded at `consensus.md:340-442`). Executed fresh: the repo-wide `rg` over all `*.go`
quoted above (zero hits, rc=1); read `internal/config/runtime.go:95-159`; read the
`parley-deck-skill/package.json` scripts block; `rg -n skillVersion` over `parley-deck-skill`;
the bounded `find` (maxdepth 5) plus an unbounded re-run over `AI_WORKSPACE` (23 deck configs
both times); the `while read`/`grep -l` pass (12 files) and the `python3` JSON comparison
(10/9/2) quoted above; the `igm-app` spot-check; `grep -rl` probes on this repository and on
the BYTE tree for the tooling-defect claim. **Not done:** re-ran no git diffs; ran no test
suites and no binaries; did not re-read `runner.go:1094-1108` for this signoff (that claim
rests on the three participants' round readings, and the draft's blind-spot paragraph
discloses the shared-tool caveat, `consensus.md:298-302`); did not read decks outside the 23
measured files beyond the quoted JSON passes. I ran no git write commands and edited no file
other than this signoff.
