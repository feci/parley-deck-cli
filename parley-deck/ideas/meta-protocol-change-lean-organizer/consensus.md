---
idea: meta-protocol-change-lean-organizer
drafted-by: kimi-1
date: 2026-09-23
---

## Protocol context attestation

```json
{"context_mode": "full", "source_sha256": "12e4b31cd3f6c1066a97caa7c7cb84213e16429dd43b26de06d4b1c32977aa18", "packet_sha256": "12e4b31cd3f6c1066a97caa7c7cb84213e16429dd43b26de06d4b1c32977aa18", "fallback_reason": null, "body_path": "/Volumes/My Shared Files/AI_WORKSPACE/parley-deck/worktrees/lean-organizer/.parley-runtime/protocol-packets/full-phase3-deliberation-12e4b31cd3f6c1066a97caa7c7cb84213e16429dd43b26de06d4b1c32977aa18.md"}
```

## Drafting context (process note, not a decision)

- The driver could not launch cross-review: its accounting refused an unavailable historical
  worktree registration (`inbox/claude-to-user_meta-protocol-change-lean-organizer_driver-error.md`;
  `organizer-notes.md`, Phase-2 fallback entry). The owner's standing authorization — "When the
  driver cannot do a step, fall back manually and record why in the idea" — was applied. Round-02
  completed manually (cross-review 1 of the deliberation cap of 3), and this `consensus.md` is
  drafted by a participant under the same recorded fallback. No guard, ledger, worktree history,
  or product code was altered to enable it.
- Quorum (locked at Phase 0): claude-1, kimi-1, zcode-1. codex-1 is organizer only — not a
  participant, never signs off, implements, or verifies code (`00-prompt.md`).
- **Role concentration (§15.5), one line:** the facilitator (codex-1) is a non-participant and did
  not draft this artifact; the drafter (kimi-1) is a participant holding no facilitator role in
  this idea — no facilitator-drafter concentration exists, and the §15.5/§15.6 drafter duties
  below are included per the Phase-3 drafter-facing duty and this idea's `00-prompt.md`
  constraints.

## Agreed decisions

Unless a line says otherwise, every mechanism below was proposed in round-01 and converged in
round-02 with the credited refinements; the round files are the canonical argument record.

### A — pure organizer is the DEFAULT for declared facilitator runs

1. `00-prompt.md` gains an **optional** `facilitator: <agent-id>` field. Absent → byte-identical
   v1.48.0 behavior (regression test). This field is what makes a run a "declared facilitator run".
2. **For a declared facilitator run, the default is the pure organizer:** the declared facilitator
   does not implement and does not verify code; participants own drafting, implementation, tests,
   and code verification; the facilitator reads verdicts and validator output. The driver refuses
   to select the declared facilitator as drafter, implementer, reviewer, or goal-done checker, and
   **escalates rather than silently falling back** to facilitator implementation when no
   participant implementer can be launched (zcode-1's rule, adopted by all).
3. **Compatibility/role exceptions, explicit and enumerated** (the default is not preserved in a
   way that drops owner-approved A): (i) `facilitator_participates: true` opts the facilitator back
   into participation — `parley preflight` fails closed (non-zero, naming both fields) when
   `facilitator:` names an agent also in `participants:` without it; with it, the facilitator keeps
   full protocol context (no audience narrowing — a signing agent never silently loses reading set)
   and role eligibility; (ii) a deck with no `facilitator:` field is untouched, byte-identical;
   (iii) the §1 user-authorized solo-exception path is unchanged; (iv) the owner's per-run
   direction is preserved — A changes the default, not the owner's authority.
4. **Driver drafting-prompt repair** (kimi-1's round-01 finding; confirmed by zcode-1 at HEAD;
   adopted by claude-1): `buildConsensusDraftPrompt` / `buildFinalDraftPrompt`
   (`internal/app/driver_consensus.go:112-131`, `:142-167` — SECONDARY via kimi-1/zcode-1 PRIMARY
   verdicts) must emit the canonical Phase-3 section names and every §15 drafter duty
   (`## Verdict conflicts` when conflicts exist, `## Drafter position changes`,
   `## Alternatives disposition`, `## Comparison & blind spots`), with both section lists derived
   from one `protocol.RequiredConsensusSections` / `RequiredFinalSections` constant and a
   **parity test** proving prompt and gate read the same value. Scope basis: the `00-prompt.md`
   constraint ("Preserve §15 dispute/provenance/alternatives duties… the old driver prompt may
   omit them") and A's "default in the skill and the driver". It adds no duty to any agent —
   §15.5/§15.6 already bind on every track (§15.7).
5. Skill text states the boundary as the default and routes the organizer through the driver
   (`run` / `continue` / `wait` / `status` / `consensus` / `preflight`); hand-launching is the
   recorded fallback. Protocol text: one permissive sentence each in §4 Phase 5, §4 Phase 6,
   §9.0, and the Quickstart facilitator row — in **all three COOPERATION.md copies** (deck view,
   `internal/protocol/defaults/`, skill `references/`).

### B — `parley wait` + `PhaseDigest` (blocking read, never a phase actor)

1. Extend the shipped TUI-only `BuildRoundDigest` (`internal/driver/digest.go:10-48`, built at
   `internal/driver/driver.go:488`, persisted as `round.digest` events — claude-1 PRIMARY) into a
   `PhaseDigest` covering design rounds, review rounds, consensus signoff state, and
   implementation status. Per-agent columns, all mechanically derived: `agent`, `path`, `filed`,
   `bytes`, `owner` (frontmatter `agent:`), `valid` (the shipped validators' verdict verbatim,
   **with the failing check named**, e.g. `missing required section "X"` — this run's heading
   stall is the evidence), `stance_flags`, `unparsed`, `fell_back`. Validity reuses
   `ValidateRoundOneArtifact` / `ValidateReviewArtifact` / `ValidateFinal` — one source of truth.
2. Digest guardrails (all adopted): deterministic and byte-identical over an unchanged tree
   (golden test); a structural test asserts no field is model-written; fail-closed `unparsed`,
   never a guess (zcode-1); the next-action line is a **fixed enumeration**, never generated prose
   (zcode-1); every row carries its raw path; the digest is printed, never persisted as canonical
   (`runs/` telemetry capture excepted as non-canonical); any ❌ / `DISPUTED` / `unparsed` /
   adverse validity sends the organizer to the raw artifact. Raw artifacts stay canonical.
3. Command: `parley wait --idea <slug> --for round|consensus|review|implementation|any
   [--timeout D] [--json]`, blocking on the append-only event log (`internal/store/events.go`),
   with bounded portable polling (interval ≥ 10 s) where no driver event log exists — Windows-safe,
   no fsnotify dependency (zcode-1). **Exit codes (converged map):** `0` boundary reached; `3`
   timeout — partial digest printed, outstanding agents named; `4` a **present** artifact fails
   the shared validator (immediate, validator error verbatim) **or** a blocking escalation /
   `driver.error` event arrives (immediate); `1` usage/IO error. Missing ≠ invalid: a not-yet-filed
   artifact keeps waiting (zcode-1's split, adopted); an invalid one exits loudly (claude-1's K-2,
   adopted). Early return on a new unanswered `to-user` escalation is included (claude-1 N-3;
   kimi-1 and zcode-1 concur).
4. **Timeout: bounded below the owner-requested 30 minutes by default, with no universal provider
   TTL asserted anywhere.** Default 25 minutes, configuration-first (seeded in the per-user
   defaults, overridable per call and per deck); acceptance criterion is exactly
   `default < 30m && configurable`; a hard ceiling rejects `--timeout` values above the active
   track's §4.0 agent timeout (claude-1 B-4). The 30-minute figure is the owner brief's
   *requirement* (testimony; `UNVERIFIED`/`RECALL` as a provider claim — claude-1 claims 16/30,
   adopted by kimi-1 round-02 position change 3 and zcode-1 round-02 position change 6); no shipped
   code, help text, doc, or comment may state it as a provider fact. Exact config key (reuse of the
   shipped `[defaults.timeouts]` keys vs a new seeded key) is deferred to implementation — both
   satisfy the criterion.
5. `wait` only **observes**; it never advances a phase (§14 human-brake boundary, stated in its
   docs and skill section).

### C — audience-scoped protocol packet, slim SKILL.md, generated organizer brief

1. `parley protocol packet` gains `--audience participant|facilitator` (name conceded by kimi-1 and
   zcode-1 after claude-1's verified collision finding: `Role` / `source.role` already denote the
   deck's `protocolRole`). The classification lives in an `audiences:` key in
   `meta/packet-applicability.yaml` — the existing §7-governed protocol file (kimi-1; one §7
   surface, adopted by zcode-1). Hard safety property: **the audience dimension may only omit
   blocks already omittable at that phase/track/transport; it can never cut below the ratified
   never-cut floor** (`internal/protocolpacket/applicability.go:193-231`), and `packet check`
   proves it with a negative test. Attestation gains an additive `audience` field (never a
   top-level `role` key — `source.role` stays unambiguous); `context_mode` values are unchanged;
   `full` remains the default for everyone; an unrecognized audience falls back to full context
   with a stated reason; `facilitator_participates: true` ⇒ full context.
2. **Facilitator verbatim reading set** (whole blocks, verbatim, with a complete omission index
   whose triggers say when to read the full source): **retain** Quickstart, §4, §5, §9, the
   active-transport §11 subsection (the protocol's own facilitator row, `COOPERATION.md:34`),
   **§2** (the same row assigns the facilitator "keep the roster"; a facilitator that cannot read
   the roster-authority rules cannot discharge that duty — claude-1 N-1, drafter-adopted), and
   **§15** — which is not a preference: the ratified never-cut floor pins `## 15.` and §15.1–15.4,
   15.7 at kernel phases {1,2,3,5,6,7,8} and §15.5/§15.6 at {1,2,3,6,7}
   (`applicability.go:190-191,208-215`; drafter re-verified at this HEAD), so the map cannot
   express its omission at deliberation phases. **Omit:** §1, §3, §8, §10, §12, §13, Appendix A,
   and non-active §11 subsections.
3. **Acceptance shape (Z-1 resolved):** the binding criterion is the named retention set + named
   omission set + complete omission index, byte-counted and recorded with the attestation
   (kimi-1's map-derived criterion). Secondary guardrail: **≤ 70,000 B** for the
   phase-1 / deliberation / github-pr facilitator body (zcode-1's converged number, replacing its
   round-01 ≤ 60 KB and kimi-1's ≤ 75 KB; claude-1's ≤ 45,000 B at phase 2 is withdrawn as
   unmeetable by whole-section omission arithmetic). The `--optimize` baseline is measured in the
   same test run — what `--optimize` cuts is unknown to all three participants (claude-1 N-2) and
   is measured, not assumed. If the §2-inclusive map-derived body exceeds the guardrail, the bytes
   are shown to the quorum before changing either the map or the ceiling; the never-cut floor is
   untouchable either way.
4. **Slim SKILL.md: core ≤ 20,000 B** (converged — kimi-1 deferred to consensus and joined 20 KB
   in round-02; claude-1 and zcode-1 independently proposed it). Relocation only, no rule text
   deleted: core keeps frontmatter (description verbatim — trigger reliability, zcode-1), Core
   Rule, Non-Solo, Required Protocol Context, Automation Mode, Startup Flow, a new driver-first
   section, File Ownership, Escalation, Quality Gates; `references/` absorbs headless-launch
   configuration, generic CLI invocation contract, timeout policy, artifact templates, roster /
   global-core / drift detail. The core names `parley run`, `continue`, `wait`, `status`,
   `consensus`, `preflight` at least once each (the direct inverse of the verified v2.12.1 state)
   and links every reference file; a skill-repo script proves every moved heading lands in exactly
   one reference and none was dropped. If 20,000 B genuinely cannot fit without dropping a rule,
   the bytes come back to the quorum — the cap is not silently raised (zcode-1's discipline).
5. **`parley organizer brief --idea <slug>`: computed, never stored** (claude-1's C3; zcode-1
   adopted it as the single re-orientation surface in round-02 position change 7, retiring its own
   `status --organizer` flag). Content: packet attestation + facilitator body path,
   `parley status --idea --json` state, current `PhaseDigest`, driver's next action, phase pointer
   from D. Contract: ≤ 8,192 B, byte-identical across two runs over an unchanged tree, writes no
   file (asserted against a read-only deck). A stored brief is a rejected sidecar pattern; any
   drift toward persisting it is a stop moment. Post-compaction re-orientation = generated brief +
   facilitator packet + `parley status` — documented as the lean loop in the slim skill core.

### D — fresh session per phase + client-accounting measurement

1. **State/handoff FILE form is retained and supported** (owner requirement for this consensus):
   the driver writes a per-phase handoff record under `runs/<run-id>/` at each phase transition,
   built on the shipped handoff machinery (`internal/runner/handoff.go` `WriteHandoffPacket` —
   zcode-1's find, claude-1 PC-6 correction recorded) and sharing the PhaseDigest computation.
   It is non-canonical driver state (§3 `runs/`), advisory, and its schema doc states the
   recomputed view (`organizer brief` / `parley status`) is authoritative on any disagreement.
   **No agent is ever obliged to author or refresh a handoff file by hand** — no new
   manual-refresh obligation (claude-1 K-3 conditions, accepted by kimi-1 and zcode-1). The
   user-facing mandatory acceptance criterion is the computed brief's output contract (C.5); the
   `runs/` record is covered by a driver unit test as repo tooling (zcode-1's refinement).
2. **Organizer usage ledger:** `parley usage ingest --agent <id> --source codex-rollout|claude-jsonl
   --path <file> --idea <slug> --phase <n>` (verb name converged on `ingest`). Binding properties:
   **no flag accepts a token count** — the tool accepts a path and parses the client's own
   accounting file itself (the mechanical answer to the earlier rejection of "self-reported usage
   records as cost evidence": a number typed by a model is `RECALL` and inadmissible; a number the
   tool parses from the client-written file is `PRIMARY` with a locator); streaming with bounded
   memory over the verified 228 MB Codex rollout; ≤ 1 KB stdout; one ledger row per ingest carrying
   idea, phase, agent, source path, parser id, ingest time, and the six `total_token_usage` fields
   verbatim — tokens only, no dollar claim; idempotent re-ingest; round-trip re-parse equality
   test (kimi-1).
3. **Attribution is explicit and bounded — never slug scanning** (claude-1 Z-2, adopted): records
   attribute by the explicit `--idea` / `--phase` arguments plus run-record timestamp windows;
   residue the windows cannot resolve is labeled `ambiguous`, with the method stated in the ledger
   header. A slug appearing in a session file is never attribution evidence — it appears for many
   reasons (discussion, grep, pasted briefs), and this run's own `organizer-usage.md` accounting
   correction (ambiguous source selector picked the wrong rollout; superseded snapshots) is
   same-class operational testimony. zcode-1's round-02 proposal already carries the explicit
   `--idea` / `--phase` arguments and retains the `ambiguous` label.
4. **Participant telemetry:** add the `kimi` case to `internal/telemetry/usage.go` (three `case`
   labels exist at HEAD — `claude:230`, `codex:252`, `opencode:261`; no kimi path; zcode is
   envelope-handled — see `## Verdict conflicts`), built from fixtures captured live this run (the
   stream-json readiness envelope in `source-context/kimi-readiness.jsonl` is testimony that a
   structured envelope exists); structured-output argv defaults only where the adapter supports
   them without breaking `-p` semantics (zcode-1). If kimi emits no machine-readable usage, record
   a visible `coverage: none` — never a fabricated number (all three).
5. Protocol text: one permissive §9 line recording that a facilitator may re-orient from the
   generated brief instead of re-reading SKILL.md and the full COOPERATION.md, and one §11
   advisory line preferring one blocking `parley wait` over repeated short polls — all three
   copies.

### Cross-cutting decisions

- **Repos in scope:** the CLI worktree (`lean-organizer`, base v1.48.0) **and** the sibling skill
  worktree (`lean-organizer-skill`, base v2.12.1). Protocol text lands in **all three
  COOPERATION.md copies** (deck view; `internal/protocol/defaults/`; skill
  `skills/parley-deck/references/COOPERATION.md`). The CLI Go drift test keeps guarding the two
  CLI copies; a skill-repo content-parity check for the third copy is added as **repo tooling, not
  a protocol obligation** (all three converged; closes claude-1 concern 5).
- **Global core:** staged from the previous core TEMPLATE plus the reviewed hunks (placeholder
  header, stub §2 — never a deck view copy), placed in `~/.parley/staging/`, with an escalation to
  the owner carrying the exact command. **`parley protocol publish` is the owner's attended-only
  action** (TTY-gated); no participant or organizer works around the gate.
- **Release (after Phase 8):** CLI 1.49.0, skill 2.13.0, CHANGELOGs, tags, GitHub releases with
  Windows assets, both Homebrew formulae, two winget PRs, npm publish, skill install `--target all
  --force` verified by content hash — per the owner brief; direct-to-main under the owner's
  run-scoped override; the deploy is independently participant-verified against real channel
  artifacts before the release is reported complete.
- **Drafter/implementer ownership (no parallel conflicting writers):** FINAL drafter = **kimi-1**
  — first recorded round-01 writer (organizer's operational statement; uncontested — claude-1 N-5
  notes the mechanical record does not establish ordering but does not contest it) and now an
  explicit volunteer via `inbox/kimi-1-to-all_meta-protocol-change-lean-organizer_drafter-claim.md`;
  kimi-1's signoff block will state `Drafter: yes`. Default Phase-5 implementer = the FINAL
  drafter (kimi-1). **zcode-1 has stated willingness** to implement (and to alternate-draft),
  "willing, not insisting" (round-02 Process): it may perfect an implementer claim via
  `inbox/zcode-1-to-all_meta-protocol-change-lean-organizer_impl-claim.md` before Phase-5 work
  begins; if it does, zcode-1 implements and kimi-1 remains drafter. claude-1 claims neither role
  and is willing if the quorum prefers (round-02 N-5). Exactly one drafter and exactly one
  implementer at any time; a claim posted after signoff completion or after implementation starts
  is void.
- **Driver/runtime gaps observed this run are evidence, not automatic new scope:** the kimi
  stream-json preflight parser rejection, the round-01 heading strictness + `parley status`
  validity concealment (claude-1 claims 21/28 verify the mechanism), the `continue --auto`
  fresh-run startup gap, and the historical-worktree registration blocker (with its mislabeled
  escalation author) are recorded organizer operational testimony. B's validity-with-reason column
  covers the concealment class inside scope B. Any further repair is optional hardening inside the
  release's driver work, re-verified at HEAD during implementation — FINAL must not turn these
  into mandatory obligations, and nothing here adds an unrelated process gate.

### A–D scope & acceptance table (organizer's pre-Phase-5 stop-rule lookup)

Every acceptance line is a command with an exit code, a byte bound on a named file, a determinism
assertion, or a negative test. Product acceptance tests implementing A–D are normal refinements of
this table; nothing here adds a mandatory obligation beyond approved A–D, and no unrelated process
gate is created. "Obligation" = a duty falling on any agent or deck.

| Item | Ships (mechanisms) | Acceptance (observable) | New mandatory obligation? |
|---|---|---|---|
| **A** — organizer does not implement / verify code, default for declared facilitator runs | optional `facilitator:` field; preflight fail-closed + `facilitator_participates:` exception; driver role-ineligibility predicate with escalate-not-fallback; `RequiredConsensusSections` prompt repair + parity test; skill boundary text; permissive protocol lines in all 3 copies | preflight non-zero naming both fields, exit 0 with the exception flag; fixture auto-drive run launches the declared facilitator for no drafter/implementer/reviewer/goal-done role (event-log assertion); absent-field deck byte-identical run plan; prompt contains canonical Phase-3 sections + §15.3/15.5/15.6 strings, derived from the constant (parity test); this run's `IMPLEMENTATION.md` records `implementer:` ≠ facilitator | **No** — optional field, permissive text, repo tests |
| **B** — `parley wait` + PhaseDigest | extend shipped digest; new read-only blocking verb; shipped validators reused; event-log wait with portable polling fallback | exit 0 with per-agent `path/filed/bytes/owner/valid(with failing check)/stance/unparsed/fell_back`; missing file → blocks to timeout exit 3, partial digest names outstanding; present-but-invalid → immediate exit 4 with verbatim validator error; escalation/`driver.error` → exit 4; byte-identical digest on unchanged tree; structural no-model-written-field test; `default < 30m && configurable`, ceiling rejects values above the active track's §4.0 timeout; digest never rewrites or replaces a round file (regression test); round-trip: one organizer tool call between launch and digest under `--no-tui` | **No** — an alternative to polling, never a required step |
| **C** — audience packet + slim skill + generated brief | `--audience` on the existing renderer; `audiences:` key in the §7-governed applicability map; SKILL.md relocation split; computed `parley organizer brief` | `packet check` passes with `audiences:` and fails any audience rule omitting a never-cut block (negative test); facilitator body contains the named retention set verbatim (Quickstart, §4, §5, §9, active §11, §2, §15), omits the named omission set, complete omission index with triggers, `source_sha256` over the full authority, additive `audience` attestation field; map-derived byte count recorded, guardrail ≤ 70,000 B (phase 1 / deliberation / github-pr) with the `--optimize` baseline captured in the same run; unknown audience → full fallback with reason; `facilitator_participates: true` → full context; SKILL.md core ≤ 20,000 B, names the six driver commands, relocation script proves nothing dropped, frontmatter/Core Rule verbatim; brief ≤ 8,192 B, byte-identical ×2, writes no file (read-only-deck test) | **No** — opt-in audience, generated views |
| **D** — fresh session per phase + client-accounting ledger + telemetry | driver-written per-phase handoff records under `runs/` (non-canonical, advisory); `parley usage ingest` path-only streaming ingest; kimi telemetry parser | handoff record written at each phase transition with schema-valid fields (driver unit test); schema doc states recomputation authoritative; no count-accepting flag exists (flag-set test); 228 MB fixture streams with bounded memory, ≤ 1 KB stdout, exactly one ledger row with idea/phase/agent/source-path/parser-id/six-fields-verbatim; idempotent re-ingest; round-trip re-parse equality; attribution by explicit args + timestamp windows, residue labeled `ambiguous`, method stated in ledger header; kimi fixture → non-nil usage or visible `coverage: none`; structured argv only where adapter-supported (test) | **No** — driver-generated records; path-only tool; no manual refresh duty |
| **Cross-cutting** | all three protocol copies; CLI drift test; skill-copy parity check (tooling); staged core from previous template + reviewed hunks in `~/.parley/staging/`; owner-only attended `parley protocol publish` with exact-command escalation; release CLI 1.49.0 / skill 2.13.0 all channels, direct-to-main per owner override; participant-verified deploy | both worktrees' test suites green incl. the new tests; the three copies carry identical hunks; staged core exists with project zones preserved; escalation note contains the exact attended publish command; post-release channel verification findings fixed before completion is reported | **No** — publish is owner-only by existing mechanism; checks are repo tooling |
| **Evidence, not scope** | startup parser gap, heading-strictness/status concealment, stale-worktree registration blocker recorded as testimony; B's validity-reason column covers the concealment class | mentioned in FINAL's risks/context; any repair is optional hardening re-verified at HEAD in implementation | **No** — explicitly not mandatory new scope |

## Agreed trade-offs

- **Savings land mostly at cached-read rates.** Polling tokens are ~54%-of-cost cached reads; the
  token share overstates the cost share (claude-1, stated up front). D's ledger measures the real
  effect per phase instead of asserting it (R6 accepted by all).
- **The largest measured lever is deliberately out of scope.** The ~22 KB skills catalogue in every
  request and the AGENTS.md read-in-full lines (study §7) are owner-environment levers, not
  parley-deck artifacts; A–D cannot reach them. Accepted: if the release is later measured against
  total organizer spend it may look like a shortfall for reasons A–D never claimed to cover; the
  ledger makes that visible rather than arguable.
- **The facilitator packet is bigger than the brief's 43–50 KB hope.** The converged retention set
  (with §2 and §15) measures ~63–70 KB depending on phase; the group declined to cut inside §4 to
  chase the lower figure (poor risk/benefit against the never-cut floor; no one will write an
  acceptance criterion the arithmetic cannot meet). Accepted: content correctness over byte
  minimalism.
- **`wait` is a new verb instead of extending `status`.** Accepted cost: a larger command surface;
  accepted benefit: `status` keeps its stable read-only-fast contract for existing scripted callers
  (zcode-1's argument, adopted by all).
- **The 25-minute default is a heuristic under owner guidance, not a provider fact.** If a real
  cache window is shorter, a wait may straddle expiry — cost: one re-invocation, mitigated because
  exit 3 prints the partial digest and outstanding agents. No artifact asserts a provider TTL.
- **The declared-facilitator default changes driver behavior for decks that opt in.** Accepted with
  the enumerated exceptions (`facilitator_participates: true`, absent-field byte-identity, §1 solo
  exception, owner direction); the preflight failure is loud, never a silent downgrade.
- **The consensus-prompt repair changes driver-emitted headings for every auto-drive deck.**
  Accepted because the old headings drifted from the protocol template; mitigated by a
  heading-consumer search during implementation and the prompt/gate parity test.
- **Digest-first reading carries a behavioral over-trust risk.** Mitigated by fail-closed
  `unparsed`, loud exit-4 degradation, path-carrying rows, and the written adjudication discipline;
  the residual risk is accepted and stated in the slim skill core.
- **Kimi structured usage may be unobtainable.** Accepted: an honest visible `coverage: none`
  beats a fabricated number; the D3 criterion passes either way.
- **Windows CI leg added for `wait`/`usage`** (zcode-1): accepted cost in exchange for the release
  shipping Windows assets honestly.

## Open items deferred to implementation

1. Exact configuration key for the `wait` default (reuse shipped `[defaults.timeouts]` keys vs a
   new seeded key) — must satisfy `default < 30m && configurable` plus the per-track ceiling.
2. Whether the §2-inclusive facilitator map fits the ≤ 70,000 B guardrail at each phase; if not,
   show the bytes to the quorum before changing the map or the ceiling. The never-cut floor is not
   negotiable.
3. What `--optimize` actually cuts (claude-1 N-2): the audience × optimizer interaction is measured
   during implementation; the AC baseline is captured in the same test run, not transcribed from
   any round file.
4. Kimi structured-usage event shape (fixtures captured this run); honest `coverage: none` path if
   unobtainable.
5. Signoff/stance parsing: verify reuse of the consensus-status parser or extract shared parsing —
   do not fork it (zcode-1 concern 7; a code-review checklist line, not a protocol obligation).
6. Heading-consumer search for the consensus-prompt repair (kimi-1's mitigation for downstream
   parsers matching old headings).
7. Skill-copy parity check shape (extend the skill repo's manifest-coverage test pattern vs a
   `compatibility.json` sha — implementation choice).
8. Exact `references/` file names and section allocation for the SKILL.md split (claude-1's
   HEADLESS_LAUNCH / ARTIFACT_TEMPLATES / ROSTER_AND_PROTOCOL sketch is the starting point).
9. Handoff-record schema and exact `runs/` location; schema doc wording stating recomputation is
   authoritative.
10. Machine-readable ledger sibling format alongside `organizer-usage.md`.
11. zcode-1's small question: whether `parley preflight` soft-warns when no `facilitator:` field is
    declared. Default: **no warning** unless a signoff note says otherwise.
12. `## Alternatives disposition` consistency duty at FINAL: FINAL.md must not contradict any
    adoption recorded below (§15.6(c)); a contradiction blocks signoff and escalates to the owner.
13. **R-1 (claude-1 signoff reservation — record correction, §15.1; accepted by the drafter,
    concurred by zcode-1, logged here per the accepted-reservation disposition):** decision C.3's
    withdrawal attribution and arithmetic framing are corrected per kimi-1's `SELF-CORRECTION` in
    the signoffs below. zcode-1's 43–50 KB facilitator-body target is the withdrawn figure
    (`round-02/zcode-1.md:30`); claude-1's ≤ 45,000 B secondary ceiling was never withdrawn and
    remains on record (`round-02/claude-1.md:188`, `:315`); the ≤ 70,000 B guardrail stands as
    conservative headroom above the measured whole-section floor (named omission set 27,420 B;
    ≈ 38.1 KB at phase 2 / ≈ 38.0 KB at phase 1 without §2; ≈ 42.1 KB with §2 retained —
    claude-1 PRIMARY measurements at HEAD `ffa4587`), not as an arithmetic limit. FINAL.md states
    the guardrail in this corrected framing.
14. **R-2 (claude-1 signoff reservation — extends open item 2; accepted by the drafter, concurred
    by zcode-1):** because 70,000 B sits ~28 KB above the measured floor, the byte guardrail alone
    cannot detect a body that still carries an omission-set block. The named-omission-set
    assertion (each named block verified absent) is the **gating** check in C.3; implementation
    records the measured facilitator body against **both** the floor and the guardrail in the same
    test run; a body under 70,000 B that retains any named-omission-set block fails C.3 regardless
    of its size. The never-cut floor stays untouchable either way.

## Comparison & blind spots

<!-- Advisory (not a gate): contradictions not smoothed into vague trade-offs;
     partial coverage (what only one participant covered); unique insights worth
     keeping; and blind spots — what did NO participant address? -->

- **Contradictions kept visible, not smoothed:** (i) the telemetry adapter-list discrepancy —
  resolved by evidence, recorded verbatim in `## Verdict conflicts`, not by the 2-vs-1 phrasing
  count; (ii) facilitator-packet byte ceiling — zcode-1 moved 60 KB → 70 KB, kimi-1 held a 75 KB
  guardrail, claude-1 proposed a named-set criterion with ≤ 45,000 B at phase 2; resolved to the
  named-set criterion with a 70,000 B guardrail and map-derived counting; (iii) D's file form —
  claude-1 round-01 proposed no file at all, kimi-1 proposed driver records, zcode-1 proposed the
  computed surface; resolved to driver records (file form retained) + computed brief (authoritative
  read path) under claude-1's two conditions.
- **Partial coverage (single-participant contributions the group adopted):** kimi-1 alone found the
  driver consensus-prompt §15 omission (`driver_consensus.go:112-131`) — the most consequential
  round-01 finding, and this manual consensus is itself live evidence of it. claude-1 alone
  measured the real `--optimize` output (65,516 B at phase 2, still containing §1/§2/§3/§8/§10/
  §12/§13/Appendix A/§11.C), the never-cut §15 pinning, `round.digest` persistence, the
  `status --json` validity concealment, and the shipped `[defaults.timeouts]` config. zcode-1 alone
  found `WriteHandoffPacket`, the missing/invalid exit-code split, the telemetry structured-flag
  heuristic (`usage.go:380-395` + `internal/runner/telemetry.go:147`), kimi's plain `kimi -p`
  default argv (`internal/agents/discover.go:331-358`), the Windows CI need, and (with kimi-1) the
  study's facilitator-row locator correction (`COOPERATION.md:34`, table spanning `:29-34`, not
  `:26-33` as the study cites).
- **Unique insights worth keeping:** "mechanism over text" — the owner's 2026-09-15 instruction cut
  Go edits 452 → 2 while the hourly input rate *rose*, so every item ships as enforced tooling
  (claude-1); the map-derived byte budget — fix the ratified map first, let the number fall out
  (kimi-1); the fixed-enumeration next-action line as the anti-summarization guardrail, and
  "show the bytes in consensus, not in FINAL" for budget changes (zcode-1).
- **Blind spots — what no participant fully addressed:** (i) **§2 in the facilitator set** was
  raised only by claude-1 (N-1) and answered by no one (round-02 files were written concurrently);
  the drafter's proposal includes §2, and signoffs ratify or contest that inclusion; (ii) §15.5 /
  §15.6 are never-cut at phases {1,2,3,6,7} but **not pinned at phases 5 and 8** — no participant
  analyzed whether a facilitator view at those phases loses the drafter-duty text and whether that
  matters; the view inherits the ratified floor by construction, flagged for implementation-time
  measurement rather than silent assumption; (iii) if the real provider cache window is **shorter**
  than 25 minutes the default wait straddles expiry — nobody measured provider behavior; accepted
  as owner-guidance heuristic; (iv) no participant verified whether the skill repo has a Windows CI
  leg today (zcode-1 raised the need; existence unverified); (v) the organizer-usage ledger's own
  same-day accounting correction is treated as testimony motivating Z-2, not as a verified defect
  in any shipped tool.
- **§15.6(b) correlated-agreement statement:** the three participants are three model families,
  but their unanimity here is a **shared prior, not independent evidence**: all three read the same
  study, the same two worktrees, and share a visible common prior that deterministic Go tooling is
  trustworthy and LLM summarization is suspect (the study's §6 rejected-patterns list, which none
  of us challenged). Convergence on A–D therefore does not independently confirm the design.
  **What would make the agreed position wrong:** the digest becoming a de-facto substitute for
  artifacts (organizers stop opening raw files); the audience map omitting a block an organizer
  operationally needs despite `packet check`; the measured savings failing to materialize because
  the out-of-scope skills catalogue dominates; kimi CLI proving unable to emit structured usage
  (D3 degrades to `coverage: none`); the lean loop not being adopted behaviorally despite the
  skill routing; the 30-minute cache guidance being wrong in the short direction. `FINAL.md` must
  state where nominally independent proposals are one family: the three designs were never
  independent mechanisms — they are one family (route the organizer through existing driver /
  renderer / digest machinery) with per-item variations, and round-02 explicitly rebuilt two of
  them onto claude-1's framing.

## Alternatives disposition

§15.6(c): an `ALT-` id and an adopt or reject with the decisive reason for each alternative raised
in round-01/round-02. FINAL.md may not contradict a recorded adoption.

| ID | Alternative (source) | Disposition | Decisive reason |
|---|---|---|---|
| ALT-1 | Existing state digests `parley status --json` / `parley consensus status --json` (all three) | **Adopt** as components | Reuse their parsing inside the digest; wrap in a blocking verb rather than changing `status`, whose read-only-fast contract is scripted by other tools |
| ALT-2 | Foreground `parley run` as the blocking mechanism (kimi-1/zcode-1) | **Reject** as the wait primitive; adopt as complement | It runs a driver cycle, not a passive read; the manual/fallback and post-compaction paths have no blocking primitive; `wait` must never advance a phase |
| ALT-3 | Existing packet renderer + applicability map machinery (all three) | **Adopt** | One §7-governed surface the drift story already knows; verbatim blocks + omission index already implemented; audience is a dimension on it, not a new renderer |
| ALT-4 | `--optimize` packet as the organizer view (zcode-1, rejected there) | **Reject** | Experimental, omission-based, participant-oriented, unratified as default; measured 65,516 B output still carries non-facilitator sections; what it cuts is unmeasured (open item 3) |
| ALT-5 | `internal/sessionstore` + `WriteHandoffPacket` as D's foundation (kimi-1/zcode-1) | **Adopt** | Existing driver state; extend rather than invent a parallel store; claude-1's round-01 table corrected to include it (PC-6) |
| ALT-6 | Existing telemetry collectors (all three) | **Adopt** | Add the kimi case to the shipped switch; no new telemetry subsystem |
| ALT-7 | Skill `references/` on-demand pattern (all three) | **Adopt** | The split extends an existing, shipped pattern rather than inventing one |
| ALT-8 | Previously rejected directions: LLM/lossy summarization, stored summary sidecars, retrieval-selected protocol text (cognee/graphify/vector), asymmetric reviewer context, silent model swaps, self-reported usage as cost evidence (study §6; kimi-1 ALT-7) | **Reject** (uphold prior rejections) | The design routes around each: digests are deterministic, path-carrying, non-canonical; the audience packet is §7-governed verbatim blocks, not retrieval; usage ingest accepts a path, never a number |
| ALT-9 | Hand-written stored `handoff.md` per phase (claude-1 set-aside; kimi-1 round-01 original D1) | **Reject** the hand-written form; adopt the driver-written non-canonical record + computed brief | A stored hand-written sidecar goes stale, becomes a second source of truth, and creates a manual-refresh obligation that trips the owner's stop rule |
| ALT-10 | Ship the digest as a TUI improvement instead of stdout (claude-1 set-aside) | **Reject** | The organizer is headless by construction; a TUI is unreadable to it — exactly why the existing digest has not helped |
| ALT-11 | Cut inside §4 to reach the brief's 43–50 KB facilitator body (claude-1 set-aside) | **Reject** | Marginal bytes come from phase prose the never-cut floor protects per-phase; poor risk/benefit once phase scoping is enabled |
| ALT-12 | Shorten the driver's internal poll interval (claude-1 set-aside) | **Reject** | Internal polling costs the organizer nothing; only organizer-visible tool calls carry resident context |
| ALT-13 | Enforce A by removing the facilitator from the roster (claude-1 set-aside) | **Reject** | The roster is machine-global and shared across decks; the per-idea `facilitator:` field is the correctly scoped mechanism |
| ALT-14 | Make `parley status` blocking instead of a new verb (zcode-1, rejected there) | **Reject** | `status` is read-only-fast everywhere and scripted by other tools; a blocking default is a behavior change for existing callers |
| ALT-15 | `--role` flag name / top-level `role` attestation key (kimi-1 & zcode-1 round-01) | **Reject**; renamed `--audience` / `audience` | `Role` and `source.role` already denote the deck's `protocolRole` (`internal/protocolpacket/packet.go:113-122`, `source.go`) — a live collision, verified PRIMARY by all three |
| ALT-16 | Usage attribution by slug mentions in session text (zcode-1 round-01) | **Reject** | Correctness: a slug appears for many reasons (discussion, grep, pasted briefs); a 228 MB rollout gives the heuristic enormous mis-attribution surface; explicit `--idea`/`--phase` + run-record timestamp windows with an `ambiguous` residue label adopted instead |
| ALT-17 | Hand-written static organizer brief file (zcode-1 set-aside) | **Reject** | Duplicates protocol content and drifts; the brief is generated from the renderer so the reading set has one source |
| ALT-18 | 25-minute timeout as a baked Go constant (kimi-1 round-01 original) | **Reject** | The 30-minute provider TTL is unverifiable from these worktrees; configuration-first keeps the intent without asserting a provider fact |

## Verdict conflicts

§15.3: contradictory verdicts on the same identified claim are resolved by reviewable evidence and
argument, never by counting participants. One round-01 contradiction existed when consensus opened;
it is recorded here with each verdict, its author, its tag, and its evidence verbatim, plus the
resolution and the dependency check.

**Claim:** which adapter cases the telemetry usage parser's switch contains.

- **kimi-1, round-01, `## Evidence & provenance` table:** "Telemetry parses
  claude/codex/opencode/zcode; no kimi path | `internal/telemetry/usage.go` (`consume` switch) |
  PRIMARY — CONFIRMED; matches study §4's 0/20 kimi capture".
- **zcode-1, round-01, `### D` item on participant telemetry:** "the collector switch has cases for
  `claude`, `codex`, `zcode` and no `kimi` case [PRIMARY: `internal/telemetry/usage.go:229-298`
  @CLI]".
- **claude-1, round-02, claim 26:** "`internal/telemetry/usage.go` has exactly three adapter cases:
  `claude:230`, `codex:252`, `opencode:261` — no `kimi`, and no `zcode` *case* | CONFIRMED |
  `PRIMARY` | `grep -n 'case "' internal/telemetry/usage.go`". (claude-1 round-01 claim 8 had
  already noted "a zcode envelope (`:190`, `:286`)" alongside the three cases.)

**Resolution (by evidence, not by count):** claude-1's round-02 verdict rests on a direct PRIMARY
check of the underlying source — a `grep -n 'case "'` over the file, which returns exactly three
case labels. That evidence applies to and entails the scoped claim (the switch's case labels are
claude/codex/opencode; zcode usage is handled through an envelope path near `:190`/`:286` that both
kimi-1's and claude-1's round-01 files independently cited, not through a `case "zcode"` label).
The contrary readings are explained rather than dismissed: kimi-1's "parses … zcode" and zcode-1's
"cases for … zcode" both used "case" loosely for the envelope handling; neither located a
`case "zcode"` or `case "kimi"` label. The drafter re-ran the resolving command at this HEAD while
drafting (`grep -n 'case "' internal/telemetry/usage.go` → `230: case "claude"`,
`252: case "codex"`, `261: case "opencode"`) and reproduces claude-1's result; this corroboration
is drafting diligence, not a new participant verdict — the resolving verdict remains claude-1's
claim 26.

**Materiality and dependency check (§15.1 / §15.3):** NOT material. The only decision that could
depend on the phrasing is D3's fix, which is "add a kimi case" under every version of the list —
all three readings agree there is no kimi path. No acceptance criterion in this consensus turns on
whether zcode support is a switch case or an envelope path. Therefore the claim is **RESOLVED, not
DISPUTED**, consensus may close over it, and FINAL.md carries no `DISPUTED` heading for it. It is
recorded here so the record shows the conflict was settled by the grep evidence, never by the
2-vs-1 phrasing count.

No other contradictory verdicts between participants are known to the drafter. Differing figures
that are compatible rather than contradictory (per-idea `status --json` 1,459 B vs 1,789 B PRIMARY
at different run states vs ~2.2 KB SECONDARY from the study; §11.C 4,190 vs 4,195 B) are noted, not
conflicts. The owner-brief 30-minute cache figure is UNVERIFIED testimony carried as a requirement,
and all three participants now tag it identically — agreement on its provenance state, not a
conflict.

## Drafter position changes

§15.5: every material change in the drafter's (kimi-1's) position since its most recent round file
(`round-02/kimi-1.md`), each with an exact prior quotation or claim identifier, the prior position,
the new position, and the correct source round path.

1. **Facilitator set composition.** Prior (round-02/kimi-1.md, `## New concerns / questions` item
   1(i)): "facilitator set = literal Quickstart table vs table + §15 — I hold +§15". New position:
   closed — §15 inclusion is forced by the ratified never-cut floor (`## 15.`, §15.1–15.4, §15.7 at
   kernel phases {1,2,3,5,6,7,8}; §15.5/§15.6 at {1,2,3,6,7};
   `internal/protocolpacket/applicability.go:190-191,208-215`; claude-1/round-02 claim 20,
   zcode-1/round-02 position change 3; drafter re-verified at this HEAD), so the audience map
   cannot express its omission at deliberation phases. The change is of basis (construction, not
   preference), and §2 is additionally retained per claude-1/round-02 N-1.
2. **Facilitator packet byte budget.** Prior (round-02/kimi-1.md, `## Current proposal` C, and
   `### @zcode-1` disagreement): "the byte budget is derived from the ratified map … with an
   absolute guardrail of ≤ 75,000 B (claude-1's criterion)". New position: guardrail **≤ 70,000 B**
   (zcode-1/round-02 position change 2 and counter-proposal), secondary to the named-set content
   criterion, with the §2-inclusive fit measured at implementation and bytes shown to the quorum if
   the guardrail fails. Source: zcode-1/round-02 `### @claude-1` ("Byte ceiling 70,000 B, not
   75,000 B").
3. **Wait exit-code map.** Prior (round-02/kimi-1.md, `## Position changes` item 5): "his
   loud-degradation exit semantics (`wait` exits distinctly on invalid artifacts and on escalation
   events, not just on timeout)". New position: adopted zcode-1/round-02's refined map verbatim —
   `0` met / `3` timeout (partial digest) / `4` present-but-invalid artifact or blocking
   escalation·driver-error (immediate) / `1` usage·IO error — with the missing≠invalid split (a
   not-yet-filed artifact keeps waiting). Source: zcode-1/round-02 position change 5 and
   `### @claude-1`.
4. **D file form.** Prior (round-02/kimi-1.md, `## Position changes` item 4 and `## Current
   proposal` D): "computed bootstrap as the organizer surface plus the driver's existing,
   non-canonical per-phase records under `runs/` … satisfying the brief's 'state/handoff file'
   letter". New position: retained, with zcode-1/round-02's refinement accepted — the user-facing
   mandatory acceptance criterion is the computed brief's output contract; the `runs/` record is a
   supported driver feature covered by a driver unit test (repo tooling), its schema doc stating
   recomputation authoritative. The owner's instruction for this drafting round requires the file
   form retained; the middle path satisfies it without a manual-refresh obligation.
5. **Timeout default mechanism.** Prior (round-02/kimi-1.md, `## Current proposal` B): "Default
   bounded wait 25 minutes, seeded in `~/.parley [defaults.loop]` configuration, overridable per
   call, documented as a conservative heuristic under the owner's requested bound". New position:
   kept, and hardened with zcode-1's criterion wording (`default < 30m && configurable`) and
   claude-1's ceiling (reject `--timeout` above the active track's §4.0 agent timeout); the exact
   config key (reuse shipped `[defaults.timeouts]` keys vs new seeded key) is deferred to
   implementation. Source: zcode-1/round-02 position change 6; claude-1/round-02 B-4.
6. **Drafter/implementer stance.** Prior (round-02/kimi-1.md, `## New concerns / questions` item
   4): "I am willing and claim both defaults, unless claude-1 or zcode-1 volunteers by the normal
   Phase 4/5 route". New position: claiming the FINAL drafter role now via
   `inbox/kimi-1-to-all_meta-protocol-change-lean-organizer_drafter-claim.md` (posted with this
   consensus, before signoff completion); my signoff block will state `Drafter: yes`. The
   implementer default stays with the FINAL drafter; zcode-1's stated willingness is preserved
   through the Phase-5 inbox claim route; exactly one writer holds each role. Source: this
   consensus's cross-cutting decision on ownership.

## Signoffs

<!-- Each agent APPENDS their signoff block. Do NOT edit others' blocks.
     Template:

     ### Signoff: <agent-id> — YYYY-MM-DD
     Status: ✅ ACCEPT           (or 🟡 ACCEPT-WITH-RESERVATIONS, or ❌ BLOCK)
     Notes: <required if 🟡 or ❌>
     Counter-proposal (required if ❌): <link or inline>

     A volunteer FINAL drafter's block must additionally state: Drafter: yes
     codex-1 is the organizer only and does not sign. -->

(no signoffs yet — claude-1, kimi-1 and zcode-1 each append their own block)

### Signoff: claude-1 — 2026-09-23
Status: 🟡 ACCEPT-WITH-RESERVATIONS
Notes: I accept A–D as drafted, the cross-cutting decisions, the A–D scope & acceptance table, the
`## Alternatives disposition`, and the `## Verdict conflicts` resolution (which correctly keeps the
resolving verdict as my round-02 claim 26 and labels the drafter's re-run as diligence, not a new
verdict — §15.1). Every counter I raised in round-02 is carried faithfully: §2 in the facilitator
set (N-1), early return on an unanswered `to-user` escalation (N-3), immediate exit 4 on a
present-but-invalid artifact with the validator string verbatim (K-2), the two conditions on D's
`runs/` record — no agent obliged to author or refresh it, recomputation authoritative (K-3),
explicit `--idea`/`--phase` + timestamp-window attribution with an `ambiguous` residue and never
slug scanning (Z-2), no top-level `role` attestation key (PC-5), and the named-omission-set +
named-retention-set + omission-index criterion with the byte ceiling secondary (Z-1). I re-verified
two load-bearing claims PRIMARY at HEAD `ffa4587` while signing: the never-cut pinning
(`internal/protocolpacket/applicability.go:190-191,206-215` — `## 15.`, §15.1–15.4, §15.7 at
`kernel{1,2,3,5,6,7,8}`, §15.5/§15.6 at `fullFifteen{1,2,3,6,7}`), and the driver prompt defect
(`internal/app/driver_consensus.go:112-131` emits `## Trade-offs accepted` / `## Deferred
follow-ups` / `## Dismissed findings`; `grep -c` for the three §15 duty headings in that file
returns 0) — so A.4 is warranted and §15 is unomittable by construction.

Two reservations, both of which fit the existing `## Open items deferred to implementation` slot and
neither of which I would upgrade to ❌:

**R-1 (record correction, §15.1 — material to an acceptance criterion).** Decision C.3 states
"claude-1's ≤ 45,000 B at phase 2 is withdrawn as unmeetable by whole-section omission arithmetic".
I did not withdraw it — `round-02/claude-1.md` `## Current proposal` C-4 still reads "≤ 45,000 B
secondary" — and the stated reason is arithmetically wrong. PRIMARY, measured this session at HEAD
`ffa4587` (`parley protocol packet --track deliberation --transport github-pr --optimize --print |
wc -c`): phase 1 = **65,415 B**, phase 2 = **65,516 B** — 101 B apart, so re-scoping the guardrail
from phase 2 to phase 1 does not change the arithmetic. My named omission set (§1, §3, §8, §10,
§12, §13, Appendix A, non-active §11) totals 27,420 B, giving a whole-section floor of ≈ 38.1 KB at
phase 2 / ≈ 38.0 KB at phase 1, and ≈ 42.1 KB with §2 (3,978 B) retained — under 45,000 B either
way. I am **not** insisting on 45,000 B as the ceiling: 70,000 B is acceptable to me precisely
because the binding criterion is the named retention/omission set, which does the real work. What I
ask is that the record not attribute a withdrawal to me that I did not make, and not rest the looser
guardrail on an arithmetic claim the measurement contradicts. FINAL.md should state the guardrail as
"conservative headroom above a measured ≈ 42 KB whole-section floor", not as the arithmetic limit.

**R-2 (implementation, extends open item 2).** Because 70,000 B sits ~28 KB above that floor, the
byte guardrail cannot by itself detect a body that still carries an omission-set block. Implementation
should therefore record the measured facilitator body against **both** the floor and the guardrail in
the same test run, and the named-omission-set assertion (each named block verified absent) must be the
gating check — a body under 70,000 B that retains any named-omission-set block fails C.3 regardless of
its size. The never-cut floor stays untouchable either way.

**One drafting ask for FINAL (my round-02 N-4, no new obligation on anyone).** The owner's stop rule
inspects **FINAL** before Phase 5, but the A–D scope & acceptance table with its "New mandatory
obligation?" column currently lives only here. FINAL.md should carry it (or an equivalent ledger) so
the organizer's pre-Phase-5 check is a lookup rather than a judgment call.

**Roles.** I claim neither drafter nor implementer (round-02 N-5, unchanged). No objection to kimi-1
as FINAL drafter, nor to zcode-1 implementing via the Phase-5 inbox claim; I remain willing to take
either role if the quorum prefers. I also note for the record that my round-02 N-5 caveat stands: the
mechanical evidence does not establish round-01 write ordering, but kimi-1's explicit volunteer claim
makes the ordering question moot.
Counter-proposal (required if ❌): n/a — not a BLOCK.

### Signoff: kimi-1 — 2026-09-23
Status: ✅ ACCEPT
Drafter: yes
Notes: I drafted this consensus and accept A–D, the cross-cutting decisions, the A–D scope &
acceptance table, the `## Alternatives disposition`, the `## Verdict conflicts` resolution, and the
`## Drafter position changes` record as written. As FINAL drafter (claimed via
`inbox/kimi-1-to-all_meta-protocol-change-lean-organizer_drafter-claim.md` before signoff
completion, per the cross-cutting ownership decision) I accept claude-1's 🟡 reservations R-1/R-2
without upgrading either to ❌: both fit the `## Open items deferred to implementation` slot and
neither changes a decision — R-1 corrects the record, R-2 sharpens an acceptance check. I also
accept claude-1's N-4 drafting ask: FINAL.md will carry the A–D scope & acceptance table (or an
equivalent ledger) so the owner's pre-Phase-5 stop-rule check is a lookup.

**SELF-CORRECTION (§15.1 — I own the corrected claim).** Decision C.3's sentence "claude-1's
≤ 45,000 B at phase 2 is withdrawn as unmeetable by whole-section omission arithmetic" is wrong on
both clauses, and I replace it with: *zcode-1's 43–50 KB facilitator-body target is the figure that
was withdrawn (zcode-1 round-02 position change 2); claude-1's ≤ 45,000 B secondary ceiling was
never withdrawn and remains on record in `round-02/claude-1.md`; the ≤ 70,000 B guardrail stands as
conservative headroom above the measured whole-section floor, not as an arithmetic limit.*
Evidence: PRIMARY — `round-02/claude-1.md:188` ("*Secondary ceiling*: ≤ 45,000 B at phase 2 …") and
`:315` ("with ≤ 45,000 B secondary") read at this HEAD, so no withdrawal exists in claude-1's file;
PRIMARY — `round-02/zcode-1.md:30` ("The 43–50 KB facilitator-body target is withdrawn; replaced by
a measured ≤70,000 B ceiling") read at this HEAD, identifying the actual withdrawn figure. For the
floor arithmetic itself (omission set 27,420 B; ≈ 38.1 KB / ≈ 42.1 KB with §2) my support is
SECONDARY, dependent on claude-1's PRIMARY measurements at HEAD `ffa4587`; I did not re-run the
byte measurement. As drafter I commit FINAL.md to the corrected framing: guardrail = conservative
headroom above a measured ≈ 42 KB whole-section floor, with the withdrawal attribution fixed to
zcode-1's 43–50 KB target. claude-1's R-2 concurred and carried: the named-omission-set assertion
(each named block verified absent) is the gating check in C.3, and the facilitator body is recorded
against both the floor and the guardrail in the same test run — a body under 70,000 B that retains
any named-omission-set block fails regardless of size. This extends open item 2; the never-cut
floor stays untouchable.

My other positions are unchanged from `round-02/kimi-1.md` and the `## Drafter position changes`
section above: §15 in the facilitator set is construction-forced by the never-cut floor, the
map-derived byte budget is primary with 70,000 B as the guardrail, the 0/3/4/1 exit map with the
missing≠invalid split, the driver-record-plus-computed-brief middle path for D, the
configuration-first 25-minute wait default, and exactly one writer per role (drafter kimi-1;
implementer defaulting to the drafter, with zcode-1's Phase-5 inbox claim route preserved).
Counter-proposal (required if ❌): n/a — not a BLOCK.

### Signoff: zcode-1 — 2026-09-23
Status: ✅ ACCEPT
Notes: I accept A–D as drafted, the cross-cutting decisions, the A–D scope & acceptance table, the
`## Alternatives disposition`, the `## Verdict conflicts` resolution, and the `## Drafter position
changes` record. My converged round-02 positions are carried faithfully: the ≤ 70,000 B
facilitator-packet guardrail secondary to the named-set criterion (my round-01 ≤ 60 KB and its
43–50 KB target withdrawn, my position change 2), the 0/3/4/1 exit map with the missing≠invalid
split, path-only `usage ingest` with `ambiguous` attribution labeling and no slug scanning (ALT-16
correctly rejects my round-01 slug-mention form), the ≤ 20,000 B skill core with the
show-the-bytes-to-the-quorum discipline, the computed never-stored organizer brief retiring my
`status --organizer` flag, and D's driver-record + computed-brief middle path with the `runs/`
record as repo tooling. Verified this session at HEAD `ffa4587`: the diff `4004a99..ffa4587`
touches only idea/inbox files, so my round-2 PRIMARY code verdicts remain valid; kimi-1's
SELF-CORRECTION citations check out (`round-02/claude-1.md:188` and `:315` — the ≤ 45,000 B
ceiling never withdrawn; `round-02/zcode-1.md:30` — my 43–50 KB figure is the withdrawn one); and
`inbox/kimi-1-to-all_..._drafter-claim.md` exists and claims the FINAL drafter role per the
cross-cutting decision. Three ratifications for the record: (1) blind spot i — I ratify §2 in the
facilitator retention set: the Quickstart facilitator row assigns "keep the roster", and a
facilitator that cannot read §2's roster-authority rules cannot discharge that duty; the
§2-inclusive body (~66.7 KB by my round-02 arithmetic, SECONDARY on the §2 byte figure) fits the
70,000 B guardrail, and claude-1's R-2 makes the named-omission-set assertion the gating check so
the guardrail's slack cannot mask a retained omission-set block. (2) I concur with claude-1's
R-1/R-2 as corrected — the guardrail is conservative headroom above a measured ≈ 42 KB
whole-section floor, recorded against both floor and guardrail in the same test run; kimi-1's
§15.1 self-correction fixes the withdrawal attribution accurately and a weakening takes effect
immediately. (3) On the verdict-conflict resolution: as the author of the loose "cases for …
zcode" phrasing, I accept claude-1's grep evidence — my "case" meant the envelope handling near
`:190`/`:286`, no `case "zcode"` label exists, and no D decision turns on the distinction. Open
item 11: the no-preflight-warning default is fine — my round-02 lean-yes was advisory and I stated
no answer was needed. Blind spot ii (§15.5/§15.6 not pinned at phases 5/8) is correctly flagged
for implementation-time measurement rather than silent assumption. Roles: I make no drafter claim
and do not object to kimi-1 as FINAL drafter; per the cross-cutting decision I intend to perfect
the Phase-5 implementer claim via
`inbox/zcode-1-to-all_meta-protocol-change-lean-organizer_impl-claim.md` before Phase-5 work
begins (this note states intent; the claim itself is that file, not this block). My round-02
concern-4 mandatory-obligation audit stands against the scope table: nothing here adds a
mandatory obligation beyond owner-approved A–D.
Counter-proposal (required if ❌): n/a — not a BLOCK.
