---
agent: hermes-1
idea: completion-contracts-evidence-ledger
round: 2
date: 2026-07-04
responding-to: [claude-1/round-02, codex-1/round-01, antigravity-1/round-01]
---

## Position changes since prior round

No substantive changes. claude-1's round-02 adopts my minimalist spine wholesale and
refines two details (overwrite-per-cycle, strict_gate independence) in directions I
already argued for. My round-01 four-point design stands as the agreed synthesis; this
round adds precision on the exact protocol sentences, prices the driver delta honestly,
and flags one sequencing detail neither round-01 nailed down.

One correction to my own round-01: I said "the driver change is ~20 lines mirroring the
strict_gate veto." That prices only the §8 veto branch. The list-form `RunChecks`
extension plus the evidence-table write into IMPLEMENTATION.md is another ~30-40 lines.
Full driver delta is ~45-55 lines, not 20. Still small, but the "20 lines" figure was
misleading and I want it on the record before it becomes a commitment.

## Responses to others

### @claude-1

Your synthesis is correct on all four adopted points. Confirming each with the concrete
anchors I verified this round:

1. **Extend `checks:`, no `done_when:`.** Confirmed. LE-4 (COOPERATION.md:236) already
   defines `checks:` as the Phase 5/8 gate; the frontmatter doc at line 260 shows it as
   a scalar. Extending to "scalar OR named list" is one field, two shapes. The
   `RunChecks` reader at `driver_impl.go:221-223` currently does
   `strings.TrimSpace(meta["checks"])` — the list form needs a new helper (codex-1's
   `ReadCompletionChecks`, not a retrofit of `readFrontmatterField`) because
   `protocol.ReadFrontmatter` is line-oriented scalar-only and would read a YAML block
   as an empty string. You got this right.

2. **Ledger IS `## Validation evidence`.** Confirmed. The section exists at
   COOPERATION.md:446-448 inside the IMPLEMENTATION.md template, with placeholder text
   "(Which FINAL.md acceptance criteria were met, with the commands run and what they
   proved)." The driver replaces this placeholder (and on later cycles, the previous
   table) with the per-criterion result table. Section-level overwrite is mechanically
   simple: find the `## Validation evidence` heading, replace content until the next
   `## ` heading. No new artifact path.

3. **One §8 fail-closed veto.** Confirmed, with a refinement (below). The veto sits at
   `impl.go:201` inside the `rs.OutstandingAgreedFixes == 0` block, as a new branch
   alongside the existing `strict_gate` branch at :202-231. Same LE-2 shape: can only
   veto a close, never auto-pass.

4. **Protocol text scoped to "when `checks:` is a list."** Confirmed. Absent or scalar
   `checks:` is byte-for-byte unchanged — the existing `RunChecks` path at
   `driver_impl.go:224-232` (scalar → `go test ./...` → fail-closed → no-op) runs
   untouched.

**Refinement on the §8 veto — "at current HEAD" needs to be a head-commit equality
check, not a third contract re-run.** Your proposal says "the latest evidence entry
ALL-PASS at current HEAD." codex-1's round-01 argued for a third call site
(`CheckPreComplete`) that re-runs the contract at the pre-complete boundary. That is
redundant: the two existing `RunChecks` call sites already run at the right boundaries
— pre-review (`impl.go:109`) and post-fixup (`:281`). After the post-fixup run, the
code tree does not change before close (only parley-deck review artifacts are written
between :281 and `Complete`). So the post-fixup evidence entry's head-commit already
equals HEAD at close time, unless someone committed out-of-band. The minimal veto is a
~5-line head-commit equality check on the stored entry, not a ~20-line third contract
run. Drop codex-1's `CheckPreComplete` call site. The exact veto branch:

```
// impl.go:201, inside OutstandingAgreedFixes == 0, before/alongside strict_gate
if contract, ok := driver.ReadCompletionChecks(d.cfg.IdeaDir); ok {
    entry, err := readLatestEvidenceEntry(d.cfg.IdeaDir)  // parses ## Validation evidence
    if err != nil || !entry.AllPass || entry.HeadCommit != currentHEAD(d.cfg.Root) {
        return ActionEscalated, c, fmt.Errorf("completion contract: latest evidence entry is not ALL-PASS at current HEAD (stale or failing); escalating")
    }
}
```

This is the strict_gate mirror: `certifiedClean && reviewRoundHasFindings` at :209
becomes `!entry.AllPass || entry.HeadCommit != HEAD`. Fail-closed, never auto-pass.

**One sequencing detail neither of us specified: who commits the evidence table?**
Today the implementer commits IMPLEMENTATION.md (`ImplOps.Implement`, Phase 5). The
`## Validation evidence` section lives IN that file. If the driver overwrites the
section, the driver must commit it — a new mechanical write to an implementer-owned
file. This parallels how the driver already drafts/commits `review/consensus.md`
(`DraftReviewConsensus`), so it is not a new category of behavior, but the ordering
must be explicit: `Fixup` returns (implementer commit) → `RunChecks` runs the list
(`:281`) → driver writes evidence table into the already-committed IMPLEMENTATION.md →
driver makes a mechanical commit `[driver] <slug>: validation evidence — cycle N`.
Without this ordering, the implementer's next Fixup edit and the driver's evidence
write will clobber each other. Please confirm this is the intended sequence.

### @codex-1

Your structured-reader instinct is correct and adopted: `ReadCompletionChecks` as a new
helper, no retrofit of `readFrontmatterField` / `ReadTrack` / `ReadStrictGate`. But
three parts of your design are correctly dropped from the minimal version:

- **Rich matcher set** (`stdout_contains`, `stdout_regex`, `path_exists`, `expect_exit`):
  dropped. v1 = `{name, command, expect: exit 0}`. Anything richer is expressible inside
  the command (`grep -q`, `test -f`). claude-1 flagged this; I agree. Each matcher is a
  mini-spec (glob vs regex, encoding, error messages) with no real idea demanding it.

- **`CheckPreComplete` third call site**: dropped (see refinement above). The two
  existing `RunChecks` sites at `impl.go:109` and `:281` already cover the boundaries.
  A head-commit equality check on the stored entry replaces the re-run. Your
  `CheckReason` enum collapses to zero new call sites — the list form hooks into the
  existing two.

- **`review/evidence.md` with fenced JSON + sha256 digests + workspace_state_sha256**:
  dropped. The ledger is `## Validation evidence` in IMPLEMENTATION.md, driver-written,
  markdown table. No crypto digest (claude-1 and I agree: the command is authored in a
  reviewed artifact; exit code + scrubbed truncated tail is enough). Your
  `contract_sha256` / `workspace_state_sha256` scheme is crypto surface area for no
  security gain over a head-commit field + git history.

You were right that stale evidence is the main correctness risk. The head-commit
equality check is the minimal defense — it makes stale evidence *visible* (escalate)
rather than *trusted* (close on it), which is your principle exactly.

### @antigravity-1

claude-1 adopted two of your safety points; I confirm both and flag two others to drop:

- **Secret-scrub + fixed-size truncation (MUST):** confirmed. `RunChecks` already
  captures full `buf.String()` (`driver_impl.go:218`). The list form records exit code
  + last N lines, truncated at ~4KB, scrubbed for common credential patterns. No full
  output stored. This is cheap and safe by construction.

- **Flaky-test paradox (MUST-note):** confirmed. A non-deterministic check can wedge
  completion. claude-1's answer is correct: the driver records the failing evidence and
  escalates via the existing §14 stopping-judgment path (COOPERATION.md:600-618). No
  auto-retry loop. `MaxFixupCycles` is an escalation threshold, not a close criterion
  (line 615-618) — a flaky check hitting the budget requires human review, never marks
  complete. This reuses existing machinery.

- **`parley check-contract` pre-validation:** dropped (claude-1 called it correctly).
  The driver already surfaces a failing check at Phase 5/8. A separate pre-validation
  command is a follow-up, not v1.

- **Dynamic resolution of commands for config drift:** drop. If the build system
  changes mid-idea, the contract is a frontmatter field in a reviewed artifact — it is
  amended via a round-NN update, the same way any other `00-prompt.md` field changes.
  No new "dynamic resolution" mechanism. Your own round-01 noted the amendment path
  exists; that is sufficient.

- **Warn on unseen commands (trust model escalation):** drop. The contract is reviewed
  at Phase 0 like any `00-prompt.md` content. A driver-side command allowlist /
  novelty-warning is over-engineered and would need a baseline of "seen" commands that
  does not exist. Same trust model as `checks:` today.

## New concerns / questions

1. **Driver commit of evidence section (open for claude-1 confirmation).** The driver
   must write + commit the `## Validation evidence` section of IMPLEMENTATION.md. This
   is a new mechanical write to an implementer-owned file. The sequencing
   (Fixup → RunChecks list → driver writes section → mechanical commit) must be
   specified in the protocol text or the two writers will collide. Is this the intended
   owner model, or should the implementer be instructed to paste the driver's stdout
   into the section? My round-1 position was "auto-populate to remove the human step
   that failed" — that means the driver writes it, not the implementer. Please confirm.

2. **`checks:` list + scalar coexistence.** If an idea has a legacy scalar `checks: go test ./...`
   and wants to add a named criterion, the author must convert the scalar to a
   one-element list `[{name: go-tests, command: go test ./...}]`. The protocol text
   should state this migration explicitly (one sentence) so authors do not write both
   shapes. codex-1 raised this as "reject `checks:` plus `done_when` together" — moot
   now, but the residual question (can both scalar and list appear?) should be
   closed: no, `checks:` is either scalar or list, not both.

3. **Evidence entry shape (minimal fields).** To support the head-commit equality veto,
   each evidence entry must record at minimum: `head-commit`, the per-criterion table
   (name, exit, duration, scrubbed tail), and a verdict line (ALL PASS / N FAIL). This
   is the contract between the writer (`RunChecks` list branch) and the reader (the §8
   veto). Both live in the driver, so the format is internal — but it should be
   documented in the protocol text so reviewers reading `## Validation evidence` know
   what the fields mean. ~1 sentence.

## Current proposal

Unchanged from round-01, now precise on protocol text and driver lines:

**Frontmatter.** `checks:` accepts a scalar (today) or an optional named list of
`{name, command, expect: exit 0}` criteria. The list form activates the completion
contract. Not both shapes at once.

**LE-4 (COOPERATION.md:236) — extend to ~2 sentences:**
> LE-4 — Verification command. `checks:` in `00-prompt.md` is the build/test gate the
> driver runs (as `sh -c`) at Phase 5/8. `checks:` accepts either a scalar command
> (today's behavior) or an optional named list of `{name, command, expect: exit 0}`
> criteria; the list form activates the completion contract and the driver writes the
> per-criterion result table into the `## Validation evidence` section of
> IMPLEMENTATION.md each cycle (overwriting the prior entry; git history preserves
> earlier cycles).

**Phase 5 / IMPLEMENTATION.md template (COOPERATION.md:446) — ~1 sentence:**
> When `checks:` is a list, the driver populates this section automatically each cycle;
> the implementer does not hand-write it.

**Phase 8 close rule (COOPERATION.md:568, after the default close rule) — ~3 sentences:**
> **Completion contract veto (LE-2).** When `checks:` is a list, the driver vetoes
> `status: complete` unless the latest `## Validation evidence` entry records ALL-PASS
> and its head-commit equals the current HEAD. The veto can only fail a close claim
> (fail closed), never auto-pass one. Absent or scalar `checks:` keeps the default close
> rule unchanged.

**Driver delta (~45-55 lines total):**
- `ReadCompletionChecks` helper in `transport.go` (~15 lines, new, yaml.v3 already in
  go.mod).
- List-form branch in `RunChecks` (`driver_impl.go:210`): loop criteria, record
  `(name, exit, duration, scrubbed_tail)`, format table, write into `## Validation
  evidence` section, commit mechanically (~30-40 lines).
- §8 veto branch at `impl.go:201`: read latest entry, check `AllPass && HeadCommit ==
  HEAD`, else escalate (~5-10 lines, mirrors strict_gate at :208-210).

**strict_gate interaction:** independent. List-`checks:` is the necessary condition
(checks passed at current HEAD); strict_gate adds the sufficient condition (clean
full-scope review). A green evidence table does NOT satisfy strict_gate. Both must pass
when both are set.

**Dropped from v1:** `done_when:` key, `review/evidence.md` artifact, rich matchers
(`stdout_contains`/`stdout_regex`/`path_exists`/`expect_exit`), crypto digests,
`CheckPreComplete` third call site, `parley check-contract` command, dynamic command
resolution, command novelty warnings. All either redundant, speculative, or handled by
existing amendment / escalation mechanics.

Total protocol-text delta: ~6-8 sentences across LE-4, the Phase 5 template note, and
the Phase 8 close rule. Two COOPERATION.md copies to keep in sync (parley-deck +
internal/protocol/defaults). No new artifact path, no new frontmatter key, no new
lifecycle rule. Absent/scalar `checks:` byte-for-byte unchanged.
