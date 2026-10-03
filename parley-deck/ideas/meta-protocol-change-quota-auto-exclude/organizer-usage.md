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
