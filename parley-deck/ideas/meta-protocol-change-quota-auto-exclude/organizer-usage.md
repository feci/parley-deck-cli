# Organizer usage ledger

These are cumulative client token totals from the organizer's own Claude Code transcript, ingested with
`parley usage ingest --agent claude-1 --source claude-jsonl` at each phase boundary. They are not model
estimates or prices. Attribution is recorded verbatim from the ingest, and no precise per-phase cost is
claimed. The machine-readable record is `usage-ledger.jsonl`. The participants' own usage is not included.

| Boundary | Accounting time (UTC) | Input (uncached) | Cache read | Cache write | Output | Events | Attribution |
|---|---|---:|---:|---:|---:|---:|---|
| Phase 0 (kickoff written, before round 1) | 2026-10-03T10:28:58Z | 316 | 19176303 | 888148 | 134778 | 155 | ambiguous |
| Phase 1 (both round-1 artifacts filed and committed) | 2026-10-03T10:58:02Z | 532 | 42966957 | 1033296 | 246932 | 263 | ambiguous |
| Phase 2, organizer session 1 (pid 1554), final total; the process was killed by the launching session's exit | 2026-10-03T11:01:10Z | 622 | 54766551 | 1090018 | 292155 | 308 | attributed |
| Phase 2, organizer session 2 (pid 57867), final total (round-2 relaunch and the claude-1 quota stop) | 2026-10-03T12:07:39Z | 264 | 14718418 | 458124 | 257691 | 128 | attributed |
| Phase 2, organizer session 3 (pid 1534), final total (claude-1 round-2 relaunch 2; died with the tmux server) | 2026-10-03T16:23:54Z | 216 | 9166220 | 860913 | 122670 | 104 | attributed |
| Phase 4, organizer session 4 (pid 53360), snapshot at the FINAL boundary (rounds 2-3 sweep, round 3, consensus, signoffs, FINAL) | 2026-10-03T19:20:50Z | 882 | 99529951 | 5884383 | 617637 | 436 | ambiguous |

Each row is cumulative over one organizer transcript, so session 1's Phase 2 row includes its earlier
Phase 0 and Phase 1 snapshots. **Organizer total over the four sessions** (the sum of the four final
rows): uncached input 1984, cache read 178181140, cache write 8293438, output 1290153, events 976.
Session 4's row is a snapshot taken just before the owner proposal note, so its last few events are not
included. Attribution is "attributed" when the ingest found the run-record window, and "ambiguous" when it
did not, verbatim from the ingest. The relay's probes and the participants' usage are not included.

## codex-1 implementation organizer (2026-10-04)

The organizer changed to codex-1 at the owner's request. The following rows are client-reported cumulative
Codex usage, not estimates. `input_tokens` includes cached input, and `reasoning_output_tokens` is a subset
of output. Do not sum cumulative snapshots of the same transcript. Separate implementer invocations are
recorded separately and never represented as independent review.

| Boundary | Accounting time UTC | Input | Cached input | Output | Reasoning output | Total | Events | Attribution |
|---|---|---:|---:|---:|---:|---:|---:|---|
| Phase 5 entry / plan recorded, organizer snapshot | 2026-10-03T23:04:28Z | 1933169 | 1571584 | 12337 | 2609 | 1945506 | 23 | ambiguous |

Source: `usage-ledger.jsonl`, codex-rollout/v1, organizer transcript
`rollout-2026-10-04T00-53-58-01a103f9-65ad-7571-8244-83b744f65e63.jsonl`.
The ingest reports no containing run-record window; idea/phase were explicitly supplied. This does not
claim precise attribution to protocol phase 5 or any monetary cost.

Phase-5 stage-1 timeout boundary (2026-10-04):

- Separate codex-1 stage-1 implementer, final before timeout: input 6810745, cached input 6032640, output 39090, reasoning output 7637, total 6849835; 46 events; attribution ambiguous.
- codex-1 organizer, cumulative snapshot (replaces earlier organizer snapshot): input 11688760, cached input 11024384, output 43572, reasoning output 16501, total 11732332; 73 events; attribution ambiguous.

Separate focused claude-1 reviewer (phase 6, final process):

- Input 192, cache read 12478234, cache write 439938, output 90979; client total 91171; 93 events; attribution ambiguous. The client total uses the Claude parser convention and is not comparable to Codex totals without separating caches.

## Supplemental source-review / owner-scope boundary (2026-10-04)

- separate claude-1 supplemental reviewer, final process: input 350, cached input 36169728, cache write 2170376, output 317666, reasoning output 0, client total 318016; 173 events; attribution ambiguous; accounting time 2026-10-04T00:21:56.664655Z.
- codex-1 organizer, cumulative snapshot replacing earlier snapshots: input 22445133, cached input 21107200, cache write 0, output 94816, reasoning output 44684, client total 22539949; 156 events; attribution ambiguous; accounting time 2026-10-04T00:22:54.388125Z.

The same parser/cache conventions and ambiguous run-window qualification above apply. Do not sum the organizer snapshots. These are reported client counters, not a monetary estimate.

## Owner-answer resumption / both-stage handoff (2026-10-04)

- Separate codex-1 both-stage implementer, final: input 20944906, cached input 19082240, output 116912, reasoning output 36555, total 21061818; 134 events; attribution ambiguous; accounting 2026-10-04T02:04:47.060393Z.
- New codex-1 organizer session, cumulative snapshot: input 24776640, cached input 23768960, output 56744, reasoning output 28754, total 24833384; 141 events; attribution ambiguous; accounting 2026-10-04T02:04:47.086291Z.

The organizer transcript is rollout-2026-10-04T02-33-50-01a10454-d471-7521-89b0-aad8df27aea1.jsonl; this is a new session after the prior blocker. Later snapshots replace this row, not add to it. Child is 02-40-20-01a1045a-c604-77e3-9100-eb902188d1db. Input includes cached input; reasoning is a subset of output. The ambiguous run-window qualification remains; no monetary estimate is asserted.

## Final producer checks / full-review dispatch (2026-10-04)

New organizer session cumulative snapshot: input 30507286, cached input 28984832, output 89081, reasoning output 40912, total 30596367; 186 events; attribution ambiguous; accounting 2026-10-04T02:34:01.055346Z. This replaces, rather than adds to, the prior 02:33:50 organizer transcript snapshot. Full-review claude-1 usage will be ingested on its completion.

Separate claude-1 full-review round-03 final: input 292, cache read 25815959, cache write 1746448, output 228716, client total 229008; 146 events; attribution ambiguous; accounting 2026-10-04T02:43:25.136456Z. Claude cache convention applies; this is a separate transcript, not an organizer snapshot.

Separate claude-1 cycle-1 consensus signer final: input 50, cache read 2034509, cache write 236622, output 43727, client total 43777; 25 events; attribution ambiguous; accounting 2026-10-04T03:07:03.066405Z. It recovered from an internal Request timed out event and exited 0; no extra signer process was started.

## Fix-up cycle 1 producer handoff (2026-10-04)

Separate codex-1 fix-up implementer, final: input 15672766, cached input 14806144, output 117021,
reasoning output 34205, total 15789787; 109 events; attribution ambiguous; accounting
2026-10-04T04:33:18.727504Z. Source rollout is
`rollout-2026-10-04T05-07-03-01a104e1-19a8-7121-aeb0-89bf627569a6.jsonl`.
It exited 0 after 5018.3 seconds, without wrapper timeout or provider quota/auth failure. This is a
separate implementer transcript, not independent review and not an organizer snapshot.

Organizer cumulative cycle-1 handoff snapshot: input 65114504, cached input 62962304, output 187168, reasoning output 90146, total 65301672; 398 events; attribution ambiguous; accounting 2026-10-04T04:42:53.749521Z. This replaces earlier snapshots of the 02:33:50 organizer transcript; it is not added to them.

## Round-04 quota-stop boundary (2026-10-04)

Separate claude-1 incomplete round-04 review: input 96, cached input 6602281, cache write 789671, output 62745, reasoning output 0, client total 62841; 50 events; attribution ambiguous; accounting 2026-10-04T04:57:13.626448Z.

Organizer cumulative snapshot (replaces earlier same-transcript snapshots): input 67065317, cached input 64889856, cache write 0, output 199929, reasoning output 94534, client total 67265246; 413 events; attribution ambiguous; accounting 2026-10-04T04:57:13.649263Z.

The review ended on quota exhaustion without a canonical artifact; partial client usage is still spend, never a completed review. Claude cache and Codex cumulative-snapshot conventions above apply. No monetary estimate is claimed.

## Round-04 owner-answer resumption — 2026-10-05

New organizer transcript: `rollout-2026-10-05T21-28-37-01a10d8a-1d2f-7303-bc7b-d04f6f392530.jsonl`.
Dispatch-boundary ingestion: `ingest appended: codex-rollout/v1 idea=meta-protocol-change-quota-auto-exclude phase=6 agent=codex-1 events=13 total_tokens=817839 attribution=ambiguous`. This is a new session, not the October 4 cumulative snapshot.
Later snapshots replace this session row; do not add them. Reviewer relaunch usage will be recorded
after the process ends. Client attribution is ambiguous, and these counts are not monetary estimates.

Successful single round-04 relaunch (claude-1, new transcript 2a3d2bae-1c17-4c60-afab-6fb90afea1d4):
input 242, cache read 39420645, cache write 1950702, output 269637, reasoning output 0, client total
269879; 121 events; attribution ambiguous; accounting 2026-10-05T19:48:39.685598Z. This is distinct
from the failed October 4 review. Claude parser/cache conventions apply; no monetary estimate.

## Phase-7 quota-stop boundary — 2026-10-05

claude-1 (a23a4ff5-8d32-4f47-a6ab-e019561a0a38.jsonl): input_tokens 62, cached_input_tokens 2648360, cache_write_input_tokens 201791, output_tokens 12117, reasoning_output_tokens 0, total_tokens 12179; 32 events; attribution ambiguous; accounting 2026-10-05T19:59:33.268184Z.

codex-1 (rollout-2026-10-05T21-28-37-01a10d8a-1d2f-7303-bc7b-d04f6f392530.jsonl): input_tokens 5836019, cached_input_tokens 5681408, cache_write_input_tokens 0, output_tokens 37263, reasoning_output_tokens 19085, total_tokens 5873282; 57 events; attribution ambiguous; accounting 2026-10-05T19:59:33.279095Z.

The organizer row replaces the earlier October 5 dispatch snapshot; do not sum those snapshots.
The failed signer is a distinct invocation and still consumed tokens. Claude cache conventions apply.
No monetary estimate or completed signoff is claimed.

## Standing-permission resumption — 2026-10-06

New organizer transcript: rollout-2026-10-06T12-18-44-01a110b9-0ae8-7111-839c-8db293ba4be8.jsonl.
Phase-7 dispatch snapshot ingested: 19 events, total_tokens 1598253, attribution ambiguous.
Later snapshots replace this same-transcript row; never sum cumulative snapshots.
claude-1 signer attempt 1 runs separately with the configured Opus 5.5/max and 1200-second ceiling;
usage will be ingested at its boundary. No monetary estimate is claimed.

## Phase-7 silent-timeout boundary — 2026-10-06

claude-1: input_tokens 90, cached_input_tokens 4745234, cache_write_input_tokens 335143, output_tokens 22362, reasoning_output_tokens 0, total_tokens 22452; 45 events; attribution ambiguous; accounting 2026-10-06T10:41:32.291353Z.

codex-1: input_tokens 6396597, cached_input_tokens 6216192, cache_write_input_tokens 0, output_tokens 23468, reasoning_output_tokens 8785, total_tokens 6420065; 50 events; attribution ambiguous; accounting 2026-10-06T10:44:13.18331Z.

The organizer row replaces the earlier October 6 same-transcript snapshot; do not sum them.
Claude cache reads/writes are reported separately under the parser convention. Its timed-out
process consumed recorded tokens without producing a signoff; no completed phase is inferred.
No monetary estimate. Raw transcripts remain outside canonical evidence; the ledger binds their hashes.


## Timeout-standing-permission resumption — 2026-10-06

New organizer transcript: `rollout-2026-10-06T13-26-57-01a110f7-7cfd-73c1-a192-9eafe182a23e.jsonl`.
Phase-7 dispatch snapshot: 17 events, client total_tokens 1278016, attribution ambiguous.
Later snapshots replace this same-transcript row; never sum them. The separate claude-1 signer
runs with unchanged Opus 5.5/max, first timeout relaunch at 2400 seconds. No monetary estimate.


## Phase-7 cycle-2 signoff complete — 2026-10-06

claude-1: input_tokens 126, cached_input_tokens 8088933, cache_write_input_tokens 952401, output_tokens 119782, reasoning_output_tokens 0, total_tokens 119908; 63 events; attribution ambiguous; accounting 2026-10-06T11:44:38.056593Z.

codex-1: input_tokens 4141063, cached_input_tokens 3983872, cache_write_input_tokens 0, output_tokens 18693, reasoning_output_tokens 8757, total_tokens 4159756; 38 events; attribution ambiguous; accounting 2026-10-06T11:44:38.067722Z.

The organizer row replaces the earlier same-transcript dispatch snapshot; never sum them. Claude cache inputs are separate under the parser convention; no monetary estimate. Successful signer artifacts and raw exit metadata were checked.


## Cycle-2 implementation subprocess boundary — 2026-10-06

Separate codex-1 transcript: `rollout-2026-10-06T13-46-05-01a11109-01db-7f60-91fd-a49427e09c8c.jsonl`. input_tokens 8046859, cached_input_tokens 7445760, cache_write_input_tokens 0, output_tokens 51138, reasoning_output_tokens 8203, total_tokens 8097997; 64 events; attribution ambiguous. It exited 0 in 1754.7 seconds without timeout. This is a distinct implementation session, not another organizer cumulative snapshot; no monetary estimate.


## Cycle-2 organizer validation / review dispatch boundary — 2026-10-06

Organizer cumulative snapshot: input_tokens 16573427, cached_input_tokens 16327040, cache_write_input_tokens 0, output_tokens 62692, reasoning_output_tokens 34452, total_tokens 16636119; 102 events; attribution ambiguous; accounting 2026-10-06T12:30:32.162925Z.
This replaces earlier snapshots of the 13:26:57 organizer transcript, never adds to them.
Child implementation usage is a distinct session. No monetary estimate or independent acceptance.


## Round-05 review / trajectory-stop boundary — 2026-10-06

claude-1: input_tokens 310, cached_input_tokens 57178199, cache_write_input_tokens 1784776, output_tokens 237589, reasoning_output_tokens 0, total_tokens 237899; 155 events; attribution ambiguous; accounting 2026-10-06T12:50:33.820343Z.

codex-1: input_tokens 21288751, cached_input_tokens 20456320, cache_write_input_tokens 0, output_tokens 83929, reasoning_output_tokens 42879, total_tokens 21372680; 143 events; attribution ambiguous; accounting 2026-10-06T12:50:33.910307Z.

claude-1 is the distinct b31f341e round-05 transcript. The codex-1 row is a cumulative
snapshot of the 13:26:57 organizer transcript and replaces its earlier snapshots; never sum them.
Claude cache read/write inputs are separate under the client parser convention. No monetary
estimate or completed implementation/close is inferred. Later organizer snapshots, if present,
supersede this boundary row without changing reviewer usage.

Final durable-stop organizer snapshot (supersedes the row above): input_tokens 22471691, cached_input_tokens 21571200, cache_write_input_tokens 0, output_tokens 99460, reasoning_output_tokens 46334, total_tokens 22571151; 155 events; attribution ambiguous; accounting 2026-10-06T12:59:52.648677Z. This is the same organizer transcript, never additive.
