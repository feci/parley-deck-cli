---
idea: meta-protocol-change-participant-dropout-review-gate-timing
drafted-by: codex-1
date: 2026-10-09
---

## Agreed decisions

Role concentration: codex-1 is organizer, participant, consensus/FINAL drafter and implementer under the controlling brief; zcode-1 independently reviews. Procedural convergence is provisional until both signoffs. No Claude participation. Kimi is owner-excluded, with no attributed position.

D1 — **Causal proof, not reviewer-seat credit.** Derive one shared predicate from validated immutable quota history, with an optional settled prospective Decision before its commit. The latest membership-changing automatic transition must have After equal current/prospective membership, contain the current implementer and sole independent reviewer in both Before and After, and reduce at least two independent non-facilitator reviewers to exactly one. The removed IDs must exactly equal eligible candidates with valid typed evidence. Reject unknown/duplicate/mismatched IDs, arbitrary RuleIDs, stale causes, corrupt/missing/pending history, manual/excluded markers and unrelated exclusions. A prospective decision must extend the actual current history; no committed-only check may deadlock Settle before CommitBatch. Kickoff readiness transitions qualify on their recorded proposed Before set, without turning excluded IDs into historical signers. A policy-only revision retains cause; any later owner/manual membership edit invalidates an earlier cause.

Accepted RuleIDs are frozen constants: `participant-failure.v1` validated by the existing paired-attempt validator, or `quota.named-reset-ge-60m.v1`, `quota.account-exhausted-no-reset.v1`, `quota.allowance-ge-24h-no-reset.v1` with the existing legacy provenance and reset constraints. The only production legacy recognizer currently carries `zcode.stderr.owner-deviation.2026-10-04`; accepting these existing rules does not invent provenance for other adapters. Named reset requires >=60 minutes; no-reset variants require no contradictory reset; arbitrary `test` or prose gives no authority.

D2 — **One numeric exception, remaining obligations intact.** Only D1 permits minimum one reviewer at precommit, review, review-consensus, goal and auto-close. The sole reviewer must be independent of the implementer and have a known distinct model in the run's roster snapshot; trim and compare case-insensitively, treating empty, `unknown` and `cli-default` as unknown. This requirement applies even if require_model_diversity=false. The shared helper supplies all consumers; no writable frontmatter assertion or general track-policy change. Roster snapshots remain configured authority where NoModelBinding prevents model observation (not proof of actual native model use); preserve normal snapshot binding and do not invent observation. Missing snapshot does not qualify.

Current signers remain required under the track; deliberation requires every current participant, implementer included. Retained BLOCK/disputes/findings keep their shipped force. Reservations still stop auto-close. strict_gate still requires a fresh full-scope round with zero findings of any severity, NIT included, plus clean consensus. Same remaining reviewer may perform the goal check in a fresh one-shot process; keep the two-minute product ceiling and existing failed/missing/self/inconclusive checks. Current-tree independent AC evidence remains required. Protected-role, usable-floor, permanently-dropped-author, serialization/recovery and all-or-none batch rules remain unchanged. Without D1 existing gates remain unchanged, including standard's already-shipped two-member behavior.

D3 — **Scoped existing supervision and truthful output contracts.** Eligible participant-failure steps use first-output120s, default stall300s after activity, heartbeat60s; explicit overrides/disables and buffers_stdout remain honored. Heartbeats never count as activity. Reuse shared supervisor classification/cleanup and the durable original+one-retry-after5s path at the same effective hard ceiling; no third attempt or new retry service. Preserve hard/operation ceilings, including readiness90s and goal120s. Wire headless signoffs inside their existing RunParticipantStep context to the shared supervisor, persisting watchdog classification after child/process-group cleanup. Protected/non-step operations retain old defaults. Manual agents exec and interactive remain hard-ceiling-only; ACP activity still uses protocol events.

Correct the proven Zcode final-text contract to BuffersStdout=true. Audit the four active roster adapters (codex, claude, kimi, zcode) without launching prohibited Claude task processes; declare default Claude --output-format text buffered as well. Codex emits progress/tool stderr or events; Kimi uses stream-json and stays streaming. Existing agy buffering remains. Custom args/streaming modes may explicitly override buffers_stdout; document that responsibility. Do not broaden to a generic adapter redesign or change models/rosters/credentials. Buffered transports disable both soft guards and retain hard bounds; healthy silent output cannot justify automatic dropout. Local fixture tests prove supervision, not actual provider latency. This is a necessary safety completion of the new signoff wiring.

D4 — **Batch consumers.** Keep current-vs-known signer/dispatch semantics, terminal survivor validation and same-tick rebind. A failure during a round, consensus or signoff batch settles only when all writers stop; remaining required signers proceed, while an earlier filed BLOCK remains binding. Regress these seams rather than redesigning shipped membership/retry accounting.

D5 — **Exact scope, compaction, delivery.** The 14 exact sequential protocol substitutions in `source-context/protocol-hunks.json` (SHA256 54d0dd07d747cf54ce9d0e234b9d65c34a8282b525bd8667090aaf19cbffcc39) are the ratified protocol changes: P1 §4.0; P2 Phase6; P3 Phase7; P4 Phase8 LE-7/11; P5 §9.0 cause/precommit; P6 §9.0 timing; C1/C5 §0; C2 §2; C3/C4/C6/C7/C8 §9 compaction. Apply identically to live, embedded and packaged COOPERATION.md. Stage core2.17.0 by these same substitutions on staged2.16.0, retaining generic zones. No applicability-map or budget change. The phase1/deliberation/github-pr facilitator guard remains70,000B; SKILL.md remains20,000B. Preview targets ~68.8KB for the existing unflagged guard; remeasure the actual renderer and do not imply every flag/phase has that size. Compact SKILL.md startup/quality/repeated exclusion prose, preserving every duty, required heading, description/Core Rule, references and drift tests; update its exclusion summary and ROSTER_AND_PROTOCOL/HEADLESS_LAUNCH guidance.

Behavior inventory: internal/membership/gates.go, membership.go and a small cause helper; internal/driver/impl.go; internal/app/driver_impl.go and consensus_request_signoffs.go; internal/runner/supervision.go, runner.go, acp.go, consult.go and launch.go (shared supervised-command seam); internal/agents/discover.go. Adjacent tests in membership/driver/app/runner/agents; existing consensus/rebind/regression suites. Update obsolete gate wording in internal/quota/record.go notices, README.md, docs/quota-membership.md, docs/agent-cli-mechanics.md, docs/agent-runtime-configuration.md, changelogs/version/compatibility manifests and protocol metadata. No generic policy schema, track, accounting, roster, credentials or Windows promotion change.

Verification inventory: cause-positive kickoff/mid/prospective commits; malformed/missing/manual/stale/history and unknown rule negatives; same/unknown model negatives (including flag=false), distinct and normalized models; standard/unrecorded behavior; driver auto-close solely with proof; reservations/strict findings/failed goal/retained BLOCK unchanged; round and signoff-batch rebind; no-first-output/stalled/timeout two-attempt behavior, cleanup and restart cap; buffered late success and hard timeout; signoff watchdog terminal classification, valid BLOCK wins; overrides/disables/clamping and out-of-scope defaults. Full Go suite `go test ./... -count=1 -timeout 45m`, build/vet, packet/drift guards; full skill npm test including Node/Python/manifests. Product tests must use real subprocess fixtures where lifecycle is the behavior at issue.

After both design ACCEPTs, publish frozen FINAL with observable ACs, merge design PR77 by merge commit, write IMPLEMENTATION plan before product edits, implement, and obtain independent Zcode review. Maximum five fix-up cycles. The controlling brief's attended close requires final review with no open CRITICAL/MAJOR, both review-consensus ACCEPTs, a fresh Zcode goal process PASS and current-tree AC evidence. This run has only an owner marker and no snapshot/history: it does not claim the product exception. Quote owner intent: “participanti nie su nevyhnutne potrebny obaja, staci jeden a to by nemalo zaseknut parley-deck” — one participant reviewer is enough and should not stall Parley. Exact controlling authority is copied in source-context/ORGANIZER-BRIEF.md.

Release after close: merge main, CLI1.53.0, skill/core2.17.0; GitHub, both Homebrew formulae, skill-only WinGet PR, 15 managed plus four named generic destinations and every SKILL.md hash. Separate Zcode channel process. npm/core publication stays owner-only, listing only unpublished commands in2.15→2.16→2.17 order in the final released note. No worktree pruning/declaration or D6 repair.

## Agreed trade-offs

- One model-diverse reviewer reduces redundancy; causality and retained obligations limit when it is allowed. Proactive manual exclusions without typed evidence stay attended, even when legitimate. This is decided scope, not an unresolved owner decision.
- A roster model that a native CLI cannot expose is configured authority, not an independently observed binding. This run additionally lacks a snapshot and therefore uses its attended bridge.
- Truthful final-text buffering prevents false permanent drops at the cost of waiting up to the existing hard ceiling. Streaming defaults are tighter; no universal prompt-hang detector is claimed.
- Goal120s remains potentially tight. Record this run's fresh healthy goal duration; if it exceeds120s, preserve that measurement as the basis for a separate tuning follow-up.

## Open items deferred to implementation

No design blocker. Method names/adjacent fixture files may follow the smallest implementation seam without changing D1–D5. Compaction exactness is independently checked before close. D6 accounting, Windows sync failures/CLI WinGet, native-positive quota provenance, alias/manual-edit guidance and persistent pre-idea proposal lifecycle remain separate prior follow-ups.

## Comparison & blind spots

Both participants independently found the split precommit/close count checks and favored existing history, reducer and retry machinery. Codex supplied the precise causal and signoff-supervisor seams; Zcode challenged marker compatibility and the buffer-safety assumption. Round02 resolves those differences explicitly. No genuine provider hang, cross-host cleanup or native Windows installation was reproduced. Unanimity among related agent systems is a shared prior, not independent evidence: both proposals remain in the existing-reducer/supervisor family. A recorded but causally unrelated exclusion earning credit, inconsistent model/count gates, an early kill of declared buffered success, loss of retained dissent, or a compaction that drops a duty would falsify the design. Refutation fixtures and independent code review must try those cases.

## Drafter position changes

- Prior `round-02/codex-1.md` T6: “No adapter/roster setting changes are proposed from this one observation.” New: correct Zcode and default Claude final-text buffering declarations and audit the four active adapter output contracts. Zcode round02 T6 combines the measured775.745s silent-success with the adapter's own final-text contract and the new signoff seam, establishing a concrete false-drop path. This changes adapter declarations only, not roster/model settings.
- Prior `round-02/codex-1.md`: “snapshot and effective-launch models cannot silently diverge.” Clarified: the run's roster snapshot is the shared gate basis; a CLI with NoModelBinding supplies no observed native model, so no observed equality is invented. Normal snapshot launch binding stays; missing/unknown snapshot models fail this exception. Source: zcode-1/round-02 T4/T5 and runmanifest's shipped authority.
- All other substantive D1–D5 positions are unchanged from codex round02. D1a's enumerated constants formalize that file's explicit recognized-rule-only constraint.

## Alternatives disposition

- ALT-1 adopt existing immutable reducer/history and prospective Decision: causality is available without another schema.
- ALT-2 reject blanket/general-track count degradation: it would relax two-person-by-design cases.
- ALT-3 adopt existing supervisor/two-attempt ledger plus the missing headless signoff seam and scoped300s stall: avoids a retry service and shortens the actual banner-then-silence case.
- ALT-4 adopt fresh-process same independent reviewer: process freshness and AC evidence are preserved without inventing a reviewer identity.
- ALT-5 adopt attended bridge for this run/manual evidence: authority is explicit and does not masquerade as product proof.
- ALT-6 reject larger packet/skill limits: compact duties in place.
- ALT-7 reject marker-based seat credit: grammar proves neither cause nor former membership.
- ALT-8 reject mandatory refusal for explicitly disabled soft guards: documented existing hard ceilings and explicit operator configuration remain binding.
- ALT-9 reject the ten-minute product goal ceiling in this scope: gather real evidence first.
- ALT-10 adopt truthful buffering declarations: shipped final-text contracts plus measured Zcode behavior refute the prior assumption that every undeclared adapter streams.

## Protocol attestation

Full live context read; phase3 attestation context_mode=full, source_sha256=packet_sha256=091e6fb841685c85fa153f7e2f05e2329c3f28dbd9b0b4c72831bc88dc458bbf; no fallback. Source claims above are PRIMARY inspection by their respective authors, not an implementer-issued code-review verdict. No product files changed before consensus.

## Signoffs

<!-- Each agent APPENDS their signoff block. Do NOT edit others' blocks. -->

### Signoff: codex-1 — 2026-10-09
Status: ✅ ACCEPT
Notes: I accept D1-D5 and the exact protocol hunks, including the disclosed buffering correction and snapshot-authority clarification. No product self-review is claimed; Zcode remains the independent reviewer.
