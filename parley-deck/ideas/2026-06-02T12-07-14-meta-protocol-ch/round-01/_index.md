---
idea: 2026-06-02T12-07-14-meta-protocol-ch
round: round-01
artifact: round-index
derived: true
generated-by: parley
token-heuristic: bytes_div_4
---

# Round Index: round-01

This is a runner-owned derived artifact. Source participant artifacts are not modified.

- Sanitizer: context-only hidden-reasoning fence removal, not secret redaction.
- Supported fences: `<think>`, `<thought>`, `<thinking>`.
- Approx tokens heuristic: `(sanitized_bytes + 3) / 4`.

| Agent | Status | Approx tokens | H2 sections | Artifact |
| --- | --- | ---: | --- | --- |
| agy | missing | 0 | none | agy.md |
| codex | ok | 2960 | Summary; Proposed approach; Concerns / open questions; Risks | codex.md |
| hermes | ok | 410 | Summary; Proposed approach; Concerns / open questions; Risks | hermes.md |

## agy

- Status: missing
- Artifact: `agy.md`
- Approx tokens: 0
- Note: artifact read failed: open /Volumes/My Shared Files/AI_WORKSPACE/parley-deck/parley-deck-cli/parley-deck/ideas/2026-06-02T12-07-14-meta-protocol-ch/round-01/agy.md: no such file or directory
- Sections: no recognized H2 sections

## codex

- Status: ok
- Artifact: `codex.md`
- Approx tokens: 2960
- Sections:
  - Summary: The proposed architecture is mostly sound: reusing the existing Phase 0-8 engine as the unit of composition preserves the current consensus mechanics, file ownership model, and non...
  - Proposed approach: Add `COOPERATION.md` §12, "Pipeline blocks & action stages", as an additive opt-in layer:
  - Concerns / open questions: The term "FINAL.md-shaped artifact" is useful for reuse but too vague for automation. Each stage needs a small typed contract so the next block can consume it without scraping pros...
  - Risks: - Automation may blur deliberation and execution. Mitigation: keep consensus, gate approval, and execute as separate recorded states.

## hermes

- Status: ok
- Artifact: `hermes.md`
- Approx tokens: 410
- Sections:
  - Summary: Proposed architecture is sound on reuse of existing round engine and additive design. Risks: driver complexity for durable resume/reconcile; side-effect gate semantics may leak int...
  - Proposed approach: Adopt the spine + gate + effects-ledger model. §12 text: blocks are Phase 0-8 invocations producing typed FINAL.md artifacts (BUSINESS_SPEC.md etc.); action blocks add `execute` su...
  - Concerns / open questions: - How does block N seed block N+1 prompt without facilitator writing another agent's file?
  - Risks: 1. Restart double-execute despite ledger (mitigate: external idempotency + dry-run mandatory before gate).

