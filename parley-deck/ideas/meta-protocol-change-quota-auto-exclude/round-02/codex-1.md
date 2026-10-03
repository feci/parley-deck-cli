---
agent: codex-1
idea: meta-protocol-change-quota-auto-exclude
round: 2
date: 2026-10-03
responding-to: [round-01/claude-1.md, round-01/codex-1.md]
---

## Summary

We agree on temporary per-idea exclusion, a simultaneous batch decision, preserved dissent, protected implementer assignments, and unchanged review/close gates. I adopt your 60-minute threshold and a stable kickoff participant list with an explicitly recorded effective quorum. I retain two substantive constraints: a timer alone is not quota evidence, and the automatic floor must count two non-facilitator participants as the brief specifies. The remaining recovery and recording details below make those shared decisions implementable.

Protocol attestation: `context_mode=full`; `source_sha256=b273af1e0649a365bf384083d27c319ddb9cb62e0efe1b757eb5e0ccfc22f388`; `packet_sha256=b273af1e0649a365bf384083d27c319ddb9cb62e0efe1b757eb5e0ccfc22f388`; `fallback_reason` absent. The shadow packet was not applied. Transport is `github-pr` with the kickoff's local-files/no-commit-by-participants/no-push/no-merge override. This relaunch replaces only my interrupted round-02 scaffold; the organizer reports that the prior process was killed when its launching session exited, which supplies no quota-exclusion evidence.

I read both prior-round artifacts, the current `status: round-02` kickoff, both required source-context documents, the readiness JSON, the protocol changelog, and this idea's worktree-prune/resolved-driver-error notes. Inspected HEAD: `8000034a3ef0cca2edcfa573f9ccf5042f58da09`. A targeted `git diff 04b22e265cc8bb96975bd59cddaf11b737d4cde9 HEAD -- <inspected source paths>` produced no output for the app/preflight, classifier, telemetry, driver, consensus, runcontrol/runstate, wait, protocol and applicability files discussed below. OpenViking tools are unavailable; local sources supply this review. The existing Graphify graph has eight organizer-document nodes: querying vocabulary tokens `[driver, execution, launch, phase]` returned organizer notes, not implementation relationships, so source inspection supplies the code evidence.

## Responses to other participants

### @claude-1

**Where we agree.** Your batch rule, per-idea scope, no timer-based rejoin, default-on recommendation, fixed floor, and preservation of already-filed objections match my preferred direction. I support covering both preflight and mid-idea participant failures; preflight-only is a partial delivery. I also accept the implementer carve-out and keeping organizer self-exhaustion outside this change. An incident that still requires intervention must remain visible in FINAL's incident table.

**1. Narrow the authority, then apply the threshold.** Your B.1 permits a parsed long reset as a positive form independently of explicit exhaustion semantics. I cannot accept that as written. The supplied incident-4 text is exactly `Unavailable (reset after 5h 51m 11s)`: its timer does not state that account credits or a usage allowance have been exhausted. Please retain that incident as unclassified for automatic exclusion unless a located provider contract establishes the meaning of that particular native error.

Likewise, stderr is a channel, not a sufficient provenance rule. Accept an adapter-recognized terminal provider-error envelope or a specifically parsed CLI terminal error on stderr/exit error; do not authorize by searching arbitrary stderr or all `error` events for a phrase. Reuse the existing parsers, with one shared typed refinement carrying the matched rule, invocation and decisive source. A successful retry, complete valid artifact, quoted example, tool error, local budget failure or interrupted organizer launch must not become an exclusion.

I change my round-01 recommendation from **15 to 60 minutes**. This is a conservative policy choice, not a provider fact. Known reset below 60 minutes takes the existing path even if the message contains a weekly/monthly phrase; otherwise a phrase fallback would defeat the threshold. Unknown reset may qualify only for explicit insufficient account credits/account quota, or a reached/exhausted hourly-or-longer usage allowance. Bare `credit balance`, generic `quota exceeded`, bare 429 and long `Retry-After` without account/allowance semantics are insufficient. Contradictory, past or unparseable reset evidence cannot be relabeled “unknown” to gain eligibility.

**2. Retain the brief's prospective floor.** `00-prompt.md:64–65` says: “**The minimum of 2.** Count only non-facilitator participants.” The same wording is in `source-context/ORGANIZER-BRIEF.md`, item 4. The explicit choice of `codex-1 + claude-1` at kickoff authorizes this particular design run; I would not generalize it into permission for an automatic exclusion elsewhere. Count distinct non-facilitator identities for this automatic rule, even if `facilitator_participates: true`. A facilitator still has its normal participation/signoff rights; it simply cannot fill this particular automatic-reduction floor. This proposal neither retroactively restaffs nor invalidates our owner-approved run.

Evaluate the entire candidate batch once. At least two usable non-facilitators must survive; unresolved failures do not count as usable. A protected designee/pin in the candidate batch also blocks automatic application of the batch. Do not silently exclude a convenient subset. Other readiness or review gates remain gates even when the numerical floor passes.

**3. Adopt your stable base membership, but not an unqualified `excluded:` reader.** My round-01 recording proposal said: “Remove the ID from current `participants:` at the same logical transition.” I replace that recommendation with your distinction: kickoff `participants:` records the original locked quorum; a validated exclusion transition yields the effective quorum for later work. Pre-kickoff exclusions must still be filtered before the initial `participants:` list, creation event, manifest and first dispatch are built.

For mid-idea changes, define `P0` as that original quorum and `A` as the applied, still-active automatic exclusions. Without an intervening owner-approved membership change, `E = P0 − A`; any separately authorized removal or re-inclusion must be reconciled into the same effective view. The readable `excluded:` marker must identify the corresponding structured transition. Do not subtract every historical free-form `excluded:` line: older records, pending transitions and manually edited prose must not silently become new authorization. Use the original identities for historical artifact/signoff validation and `E` for future dispatch and required signers. Continue/replay, role selection, captured app adapters and status/digest readers must use the same validated reduction.

This resolves our data-model disagreement in favor of your stable base, while retaining crash recovery. Reusing an atomic frontmatter writer is useful but not sufficient for a transition also represented in events and a manifest. Persist an intent with before/after membership, invocation IDs and a transition ID; serialize by idea; apply the canonical marker and projections; then record application. Before dispatch or close, recover a matching interrupted transition idempotently, or stop on mismatched state. Read-only views report inconsistency and never repair it.

**Executed evidence for your V12.** I ran the real `app.Run` entry point using a temporary Go overlay test, an isolated `PARLEY_HOME`, fake local CLI wrappers, and the existing `pingProbe` test seam. Inputs: selected `[codex, claude, kimi]`; codex/claude probes ready; kimi probe `ClassProcessFailure`, exit 1; flags `--no-tui --no-auto --yes`. All three fake round CLIs deliberately exit 1, so no hosted participant or provider was invoked. The repository's implementation files were unchanged.

Command:

```text
GOPROXY=off GOTOOLCHAIN=local GOCACHE=<temporary>/cache go test \
  -overlay <temporary>/overlay.json ./internal/app \
  -run '^TestCodexQuotaR2RunExclusionProbe$' -count=1 -v
```

The overlay and diagnostic source are at `/var/folders/yt/p2sr23f12_qcfx_w2z5c1p4r0000gn/T/codex-quota-r2-2_6m3qte/{overlay.json,probe_test.go}`. Relevant output:

```text
run exit=1 (fake round CLIs deliberately exit 1)
participants: [codex, claude, kimi]
excluded: kimi — process-failure:exit-1 — confirmed 2026-10-03
replayed run participants=[codex claude kimi]
manifest participants=[codex claude kimi]
run.created participants=[codex claude kimi]
agent.started agent=claude
agent.started agent=kimi
agent.started agent=codex
--- PASS: TestCodexQuotaR2RunExclusionProbe (1.97s)
ok parley-deck-cli/internal/app 2.294s
```

Source locator: `internal/app/app.go:1922–1943` assigns `preflightExcluded = excluded`, then passes both `Participants: participants` and `Excluded: preflightExcluded` into `runcontrol.Create`. This executed output is PRIMARY evidence for your independent disposition of V12. Because I already owned the related assertion in round 1, I do **not** issue myself a CONFIRMED verdict. The diagnostic's PASS means the probe ran and observed the exclusion record; it does not mean the behavior meets the proposed rule.

**4. Re-evaluate the failed round explicitly.** `internal/driver/driver.go:467–468` reads `case "round.incomplete": return false, nil // authoritative block`. After membership changes, a helper that filters dispatch alone would leave this terminal event in place. Revalidate all surviving artifacts and append a new terminal round evaluation bound to the exclusion revision; preserve the failed invocation and earlier incomplete event. Do not manufacture acceptance for the excluded participant. Existing valid artifacts win over the provider-failure candidate; partial/invalid files remain preserved and identified as incomplete, rather than deleted, overwritten or promoted by filename existence.

For signoffs, reuse `internal/consensus/consensus.go:523–525`: “separates who may SIGN (known) from who is AWAITED (required).” Its `StatusBlock` path sets `hasBlock = true` at lines 571–572 and produces `TriageBlocked` at 589–590. Keep that historical-block path active while changing only the required set. Engaging a counter-proposal is necessary but does not by itself erase the old BLOCK; resolution must meet the protocol's normal evidence/signoff or explicit owner-ruling process. Exclusion supplies neither a withdrawal nor a verdict on a DISPUTED claim. All strict-gate findings retain their existing force, including low-severity findings.

**5. Preserve `wait`'s existing exit meanings.** Your B.7 explicitly proposes returning at every exclusion as degradation. I recommend a non-blocking digest annotation for a successfully applied reduction, and normal waiting for the requested boundary. PRIMARY source evidence: `internal/app/wait.go:10–19` defines exit 0 for a reached boundary and exit 4 for an invalid artifact or newly arrived blocking escalation/driver error; `wait.go:163–186` implements those branches. A `blocking: no` exclusion notice is not one of those exit-4 conditions. Add effective-membership refresh inside the wait loop, since its current `BuildPhaseDigest(..., ideaStatus.Participants)` call at line 156 uses the startup participant list. An actual floor/role/integrity failure still produces a blocking notice and follows the existing degradation behavior.

**6. Bound the rollout and the size claim.** I agree with `[defaults].quota_auto_exclude` and a per-idea false override. Record the resolved, presence-aware policy when a new idea/run starts. Legacy or already-running ideas without that recorded authorization retain confirmation behavior; this is a feature rule, not a claim that protocol-version pinning exists. For an idea continuing in another run, reuse its recorded policy instead of silently activating a newer machine default. Missing evidence, malformed policy or an ambiguous record cannot authorize reduction.

Your “knob off → today's behavior byte-for-byte” should instead promise **no automatic quorum reduction**. Fixing V12 or the preflight 503 misclassification necessarily changes the affected existing behavior; those corrections need explicit regression expectations, not an incompatible byte-identity promise. Separate generic provider-failure classification from the stronger predicate that authorizes quota exclusion. Unrecognized exhaustion must not fall through into automatic process-failure exclusion.

Your packet-map question has a small answer: PRIMARY source evidence at `parley-deck/meta/packet-applicability.yaml:98–99,114–115` places both `## 5. Quorum and async participation` and `### 9.0 Pre-idea readiness check ...` under `include: always`. Keep the new core rule inside those existing blocks and check packet rendering; no new map entry is needed solely for a paragraph inside them. New headings or moved obligations would require a fresh map review. I withdraw no substantive recovery requirement to fit the estimated 300–450 LOC; that estimate should remain unverified until implementation is scoped.

## Refined position

The following is my proposed consensus basis, superseding my round-01 choices where stated:

1. **Authorize prospectively.** Default-on for newly created ideas after ratification and delivery, with presence-aware machine/deck configuration and a per-idea opt-out. Record the effective policy; keep older ideas on confirmation unless the owner explicitly enables it. No permanent roster edit and no silent change to `roster_change_policy`. Document this as the narrow per-idea exception, leaving roster-file membership changes under their existing policy.
2. **Require positive terminal evidence.** Failed invocation, no valid completed artifact, adapter-recognized terminal provider error, explicit account-credit/account-quota or exhausted usage-allowance semantics, and no later successful retry. A reliable known reset must be at least **60 minutes** away; without one, only the explicit no-reset forms above qualify. Timer/status/substring alone never authorizes exclusion.
3. **Apply a whole batch at a quiescent boundary.** At least two distinct usable non-facilitators survive. A per-idea designation or implementation pin belonging to a candidate prevents automatic application. Floor or protected-role failure changes no membership and escalates with all candidates and the arithmetic.
4. **Keep original and effective membership distinct.** Filter before kickoff; after kickoff preserve `P0`, record each applied exclusion and derive `E`. All mutable consumers use the same reduction; historical identities, signoffs and findings remain interpretable. Reconcile interrupted transitions and failed-round evaluations before proceeding. Reuse existing store/manifest/atomic-write primitives, with idea-level serialization, rather than introducing another service.
5. **Preserve obligations.** No automatic implementer or partial-draft ownership transfer. The existing global-default fallback remains available before a pin exists. Recompute reviewer counts, diversity and independent goal-check eligibility; preserve BLOCK, material DISPUTED claims, open findings and strict-gate conditions. If lawful progress remains impossible, stop rather than waive a gate.
6. **Record and display honestly.** One structured transition bound to the invocation/batch, a linked automatic `excluded:` marker, and one deduplicated non-blocking owner notice with evidence, survivors and reset hint. Scrub evidence before persistence. `status`, `wait` and organizer brief show the effective quorum and pending inconsistencies; a successful reduction alone does not change `wait`'s exit contract. No same-idea timer rejoin; ordinary next-idea probing uses the unchanged roster, and same-idea re-inclusion remains owner-confirmed.

Alternative dispositions recommended for consensus: adopt ALT-1's shared native-error refinement, ALT-2's standard-library parsing without a timer, ALT-3's temporary exclusion recording, ALT-4's existing persistence primitives, and ALT-5's known/required signer distinction. Retain ALT-6's manual escalation for ambiguous failures. Reject permanent roster mutation, generic regex authority, a new quota-polling subsystem and ALT-7 preflight-only as the complete solution. Our proposals are one architectural family around the existing Parley driver; their agreement is not independent evidence that its failure paths are safe.

For FINAL, retain the CLI/skill file groups and tests from my round-01 inventory, narrowed to these choices. The central rule belongs in §5/§9.0, with cross-references where templates, implementer eligibility and verification consume membership; mirror all three COOPERATION.md copies and record the ratified change. The acceptance set must include: quoted/tool/stderr false positives; the ambiguous incident-4 503; threshold and contradictory resets; later success/artifact-wins; simultaneous reductions and facilitator floor; kickoff filtering (the diagnostic above); resume and interrupted writes; failed-round reconciliation; preserved dissent; designation/reviewer/goal-check gates; legacy policy and explicit false; and consistent notices/status/wait. Preflight-only or classifier-only delivery must not be described as completed mid-idea support.

## Remaining disagreements

- **Evidence semantics:** I would block a FINAL that accepts the supplied 503 or a long reset alone as exhaustion, or treats arbitrary stderr as trusted provider evidence. Counter-proposal: the terminal-evidence predicate in items 1–2 above, plus a provider-specific rule only after its meaning is evidenced.
- **Automatic floor:** I would block silently counting the declared facilitator toward the brief's two-non-facilitator floor. Counter-proposal: retain that prospective floor and describe this owner-approved run separately. If you still prefer changing the brief, surface that exact policy choice for an owner ruling rather than infer it from kickoff.
- **Contract details awaiting your agreement:** stable original membership with validated transitions, no automatic activation for legacy ideas, and non-blocking `wait` annotation rather than a new exit condition. These are concrete reconciliation proposals, not claims that you have already accepted them.

The 15-versus-60 threshold and my former requirement to rewrite mid-idea `participants:` are withdrawn as disagreements. No missing human input blocks this cross-review artifact; remaining differences can be addressed in the next round before consensus. This file is a completed round-02 position, not a consensus signoff or an authorization to implement/publish.
