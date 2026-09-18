---
agent: codex-1
idea: meta-protocol-change-evidence-first-efficiency
date: 2026-09-16
status: incomplete-history-diagnosis
---

# Concrete history-preview blocker

The declared-unavailable migration candidate has passed its scoped checks but the actual read-only amendment preview still fails. This note records diagnosis, not permission to declare missing history, choose a carried count, change a cap or run the amendment.

The preview targets `meta-protocol-change-evidence-first-efficiency-v2` in the amendment worktree and names the two retained unavailable registered paths. It exited 1 in 3.562 s with `historical run identity is missing or conflicting`. Git registration bytes and the observed empty shared budget directory were unchanged.

A diagnostic-only error-context overlay located the first refusal at `parley-deck/runs/20260510T194003Z/events.jsonl`. A bounded read-only scan then examined 150 visible run instances across the 27 retained registrations. All 25 identity issue rows have the same 117-byte source SHA256:

`ee4f52b717159e54aa8f144161d14467160cfc6846aed97abbdd6768683c6906`

The sole event is:

```json
{"time":"2026-05-10T19:40:03.126637Z","type":"run.created","data":{"mode":"auto","task":"smoke implementation run"}}
```

There is no driver.json or run.json supplying an identity. The same bytes are tracked from initial CLI commit `3ec10ac`. Copies in different worktrees are diagnostic observations, not 25 independently inferred actions. Conversely, no agent.started event in this surviving file is not proof that no historical work occurred. The scanner read 2,218,000 bytes with 64 MiB aggregate/16 MiB per-file bounds and verified unchanged registration bytes. It does not replace authoritative Go inspection or establish a complete floor.

Evidence is retained under `.parley-runtime/managed-continuation-20260915/declared-history-real-preview/` and `declared-history-preview-diagnosis/`; the latter contains `visible-run-identities.json`. No source or historical evidence was deleted, rewritten, assigned a fabricated identity, classified as zero or silently exempted.

Claude is preparing a source-level remedy that can expose a concrete preview while retaining explicitly unknown scope/count. It must precede a reviewable operator migration payload. The earlier roster/full-four authorization does not authorize an unknown-history count.

# Separate launcher state

Zcode's preparation invocation `1d612975-13a8-46e8-b162-7f3c1658dda0` timed out after 1,800 s (exit143; wrapper1), leaving two owned Go files and no owned note. Its unknown cost stays unknown. An immutable 452-file recovery-plus-launcher snapshot failed compilation in 0.950 s: four references disagree with the report field name and the fake discovery seam omits a context argument. No behavior tests or real round ran.

A bounded owner correction has been assigned, including the exact v2 slug invariant. The eventual launcher must still require an existing admissible cycle binding, preserve Codex's artifact and make exactly one instrumented RunRound call for the four members with a shared900s context. Source preparation is not actual amendment execution.

