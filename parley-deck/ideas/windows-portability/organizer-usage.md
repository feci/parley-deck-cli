# Organizer usage ledger

Pinned client accounting source: `/Users/tomasfecko/.codex/sessions/2026/09/24/rollout-2026-09-24T23-04-31-01a0d53b-f5e5-7b32-bfed-de64b1926730.jsonl`. Cumulative token_count totals include kickoff preparation, not model estimates. Cached input is a subset of input; reasoning is a subset of output. No monetary-cost claim. Request events count accounting records, not unique billing requests.

CLI `parley usage ingest` was used for phase 0 and reported attribution=ambiguous (16 events, total_tokens=1185971 at ingestion); this table provides the explicitly pinned session reading, without relabeling the CLI attribution as exact.

| Boundary | Timestamp | Input | Cached input | Output | Reasoning | Accounting events |
|---|---|---:|---:|---:|---:|---:|
| Phase 0 to independent round 1 | 2026-09-24T21:11:07.071Z | 1390844 | 1294592 | 8456 | 412 | 18 |
| Round 1 complete; cross-review launched | 2026-09-24T21:33:03.603Z | 6895335 | 6750464 | 20193 | 5657 | 65 |
| Round 2 complete; focused round 3 active | 2026-09-24T22:05:07.911Z | 16074825 | 15873408 | 38549 | 17594 | 119 |
| Round 3 complete; final scheduled cross-review opened | 2026-09-24T22:10:05.681Z | 17465009 | 17243264 | 44608 | 21675 | 126 |
| Round 4 resumed; awaited participant outputs and recorded current memory availability | 2026-09-24T22:20:17.434Z | 19953018 | 19424000 | 54278 | 25539 | 145 |
| Round 4 complete; durability conflict escalated at the cross-review cap | 2026-09-24T22:38:33.156Z | 23233082 | 22653056 | 59823 | 27034 | 190 |

## Resume after owner authorization — 2026-09-25

Client accounting source for this resume: `/Users/tomasfecko/.codex/sessions/2026/09/25/rollout-2026-09-25T08-30-09-01a0d741-ce48-70e2-a848-0dad94dd4d98.jsonl`. Separate cumulative totals from the prior session; no monetary estimate. Dispatch: three real CLIs, original model selectors, one focused round 5. Driver G2 reproduced; pinned CLI 1.49.0 dispatch/wait reused.

| Boundary | Timestamp | Input | Cached input | Output | Reasoning | Accounting events |
|---|---|---:|---:|---:|---:|---:|
| Round 5 launched under owner authorization | 2026-09-25T06:32:27.994Z | 266663 | 228736 | 3198 | 303 | 8 |
| Round 5: Claude and Zcode complete; Kimi pending at first 20-minute wait timeout | 2026-09-25T06:54:44.289Z | 2429773 | 2366208 | 8282 | 1976 | 46 |
| Round 5 complete: 3/3 valid, conflict escalated, exit handoff written | 2026-09-25T06:55:39.573Z | 2643492 | 2574848 | 9318 | 2590 | 49 |

## Sequential round 6 resume — 2026-09-25

Client source: `/Users/tomasfecko/.codex/sessions/2026/09/25/rollout-2026-09-25T18-25-22-01a0d962-bf99-7581-9bfd-8d8c2fbbd7f4.jsonl`. Separate cumulative totals, no monetary estimate. G2 persists; original selectors dispatched one participant at a time. Current global CLI probes preceded rediscovery of the existing pinned 1.49.0 tool; all round-6 dispatch and validation use that private binary.

| Boundary | Timestamp | Input | Cached input | Output | Reasoning | Accounting events |
|---|---|---:|---:|---:|---:|---:|
| Round 6 opened; Claude alone launched | 2026-09-25T16:28:13.787Z | 351540 | 285696 | 3823 | 456 | 9 |
| Claude round 6 exited 0; artifact valid; Kimi next | 2026-09-25T16:41:35.297Z | 2485593 | 2397952 | 9469 | 2452 | 37 |
| First 20-minute round wait timed out; Claude valid, Kimi active, Zcode unlaunched | 2026-09-25T16:49:52.648Z | 3140452 | 3035264 | 11462 | 3683 | 44 |
| Kimi round 6 exited 0; artifact valid; Zcode next | 2026-09-25T16:59:13.586Z | 4347210 | 4235136 | 14140 | 5167 | 56 |
| Round 6 complete; 3/3 valid; consensus drafting begins | 2026-09-25T17:06:48.896Z | 5619207 | 5483776 | 17103 | 6945 | 67 |
| Consensus draft completed; Claude signoff active | 2026-09-25T17:29:45.461Z | 8538774 | 8372608 | 26326 | 11693 | 88 |
| All three consensus ACCEPT signatures present; FINAL drafting authorized | 2026-09-25T17:41:56.079Z | 11219285 | 11034112 | 32718 | 16400 | 105 |
| FINAL complete; closing design and assigning Phase 5 to Zcode | 2026-09-25T17:52:34.403Z | 13163546 | 12960768 | 39973 | 20979 | 116 |
