---
agent: hermes
idea: tui-agent-output
phase: review
round: 2
date: 2026-06-07
---

## Summary
Round-02 re-review of fix-up cycle 1. All round-01 blocking items (codex + agy) have been resolved in the CR ingestion state machine, offset tracking, steer event ordering, artifact tail rendering, UTF-8 safety, relArtifact, and error persistence. New tests cover the high-risk cases. No regressions or new issues found. Build/vet/full suite green.

## Verification
- split \r\n across ticks (codex/agy): FIXED — ingestTranscriptBytes now carries per-stream crPending; trailing \r is deferred and resolved on next chunk (\r\n → commit, lone \r → rewrite). TestIngestSplitCRLFAcrossTicks passes.
- readAppendedChunk offset duplication (agy): FIXED — c.offset = start + int64(len(data)) (actual bytes read). TestReadAppendedChunkNoDuplication passes.
- steer reply lost same-tick (codex): FIXED — eventsMsg does refreshBuffers (drain) BEFORE appendSteerEvents; appendSteerEvents calls advanceBuffer then commits trailing partial before clearing steer cursor/partial/crPending. TestSteerReplyTextAndMarkerBothKept passes; no double-commit.
- artifact view head not tail (codex/agy): FIXED — renderArtifactView slices lines[len(lines)-rows:].
- UTF-8 cap split (codex): FIXED — capSeg advances to utf8.RuneStart after trimming. TestPartialCapIsRuneSafe passes.
- relArtifact relative ideas/ (agy): FIXED — explicit HasPrefix("ideas/") branch.
- /artifact outside tab error cleared (codex): FIXED — error retained when no active agent tab.

## New findings
None.

## Verdict
ACCEPT. All blocking items resolved with no regressions.