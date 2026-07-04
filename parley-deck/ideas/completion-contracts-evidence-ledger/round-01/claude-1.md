---
agent: claude-1
idea: completion-contracts-evidence-ledger
round: 1
date: 2026-07-04
---

## Summary

Adopt the Hermes-style completion contract as a thin, additive layer over what the
driver already does. Key design stance: **the contract is data, the ledger is an
artifact, the driver is the only writer of the ledger, and the ledger never closes
anything by itself** — it grounds claims that humans and reviewers still judge.

## Proposed approach

### 1. Contract shape — `done_when:` in 00-prompt.md frontmatter

```yaml
done_when:
  - name: unit-tests
    run: go test ./...
    expect: exit 0            # default; the only v1 expectation type
  - name: version-bumped
    run: test "$(cat VERSION)" != "1.34.0"
```

- Optional list. Absent ⇒ today's behavior byte-for-byte (same rule §4.0 used for
  `track:` — present activates, absent is inert).
- v1 expectation = exit code 0 only. Output-pattern and file-exists checks can be
  expressed *inside* the command (`grep -q`, `test -f`), so a richer `expect:` grammar
  is deferred — YAGNI until a real idea needs it.
- Commands run with the repo root as cwd, same trust level as RunChecks today. The
  contract lives in a reviewed artifact (00-prompt.md), so malicious-command risk is
  the same as any reviewed change; note it in the protocol text, don't build a sandbox.

### 2. Ledger artifact — `review/evidence.md`

Driver-owned canonical file, append-only per verification pass:

```markdown
## Evidence pass 3 — fix-up cycle 2 — 2026-07-04T14:02:11Z
head-commit: abc1234
| criterion | command | exit | duration | output-digest |
|---|---|---|---|---|
| unit-tests | go test ./... | 0 | 41.2s | sha256:9f3a… (last 4 lines: "ok …") |
verdict: ALL PASS (2/2)
```

- Output digest = sha256 of full output + last N lines truncated and secret-scrubbed
  (reuse the existing redaction used for agent transcripts if present; otherwise strip
  lines matching common credential patterns). Full output is NOT stored — digest +
  tail is enough to ground the claim and safe by construction.
- Owner: the driver writes it (facilitator-mechanical, like consensus scaffolds). No
  agent identity problem: it is machine evidence, and §Ownership already lets the
  facilitator create mechanical scaffolds. Protocol text should name the driver/
  facilitator as owner explicitly.

### 3. Driver enforcement points (minimal touch)

`internal/driver/impl.go` already gates on `RunChecks` at two points (post-implement
line ~109, post-fixup line ~281). Extend, don't fork:

- New interface method `RunContract(ctx) (ContractResult, error)` alongside RunChecks;
  implemented in `internal/app/driver_impl.go` by parsing `done_when:` from the idea's
  00-prompt.md (reuse the frontmatter reader that `ReadTrack`/`ReadStrictGate` use).
- Where RunChecks passes today, also run the contract when present; write the evidence
  pass; if any criterion fails ⇒ same path as RunChecks failure (block transition,
  surface detail).
- Phase-8 "complete" additionally requires: latest evidence pass is ALL PASS **and**
  its head-commit == current HEAD (staleness guard — evidence from an older commit
  proves nothing).

### 4. Protocol text (both COOPERATION.md copies + skill fallback)

- §5: define `done_when:` (optional, mechanically verifiable, reviewed like any
  artifact) and the reviewer input rule: reviewers MUST be given `review/evidence.md`
  when it exists, and a prose claim that contradicts the ledger is a finding.
- §8: completion rule — with a contract present, "complete" requires a current-HEAD
  ALL-PASS evidence entry; without one, today's rules apply unchanged.
- One-line cross-ref from §6 (review inputs) — keep it to a sentence to respect the
  progressive-disclosure layout.

## Concerns / open questions

1. **Who runs evidence when there is no driver?** Manual facilitation (skill-only
   decks) has no Go driver. Proposal: protocol text says "the driver or the
   facilitator runs the contract"; `parley evidence run <slug>` as a manual CLI
   subcommand so skill-only flows can still produce the ledger.
2. **Command portability** — contracts written on macOS may not run on a reviewer's
   Linux box. Mitigation: evidence is produced by ONE designated runner (implementer's
   env), recorded with env metadata (os/arch, go version); reviewers read, not re-run.
3. **Timeout** — a hung contract command must not wedge the driver: per-criterion
   timeout (default e.g. 10m, configurable), timeout ⇒ FAIL with reason.

## Risks

- **False confidence**: a weak contract (e.g. `true`) yields a green ledger that means
  nothing. Mitigation: review of 00-prompt.md covers contract quality; protocol notes
  a green ledger is necessary, never sufficient.
- **Secret leakage via output tails**: mitigated by scrub + truncation; state it as a
  MUST in protocol text.
- **Scope creep** toward a CI system: resisted by v1 = exit-code-only, local-only,
  one runner.
