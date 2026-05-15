---
idea: meta-protocol-change-agent-teams-patterns
status: final
author: codex
consensus-date: 2026-05-14
participants: [codex, claude, gemini, hermes]
---

## Final plan / specification

### Goal

Improve Parley Deck's cooperation process by borrowing the useful, vendor-neutral parts of Agent Teams practice while preserving Parley Deck's core invariants:

- participant-owned canonical files;
- explicit rounds and append-only signoffs;
- no facilitator synthesis as a substitute for participant artifacts;
- flat quorum and durable audit trail;
- transport independence.

The conclusion is yes: there are useful Agent Teams patterns, but they should enter Parley Deck as small protocol guidance, not as a runtime dependency or a hierarchical replacement.

### Scope

- Add advisory per-idea roles/lenses.
- Clarify that internal helper agents/tools are allowed but are not protocol participants.
- Add a lightweight participant-count and role-use decision framework.
- Add a Phase 5 planning checklist for non-trivial implementation work.
- Keep peer messages inside existing `inbox/` mechanics and require decisions to be mirrored in canonical round/review/consensus files.

### Implementation details

Recommended future `COOPERATION.md` changes:

1. Add optional `roles:` metadata to `00-prompt.md`.

   Keep `participants:` as a list of IDs for parser compatibility:

   ```yaml
   participants: [codex, claude, gemini]
   roles:
     codex: protocol-facilitator
     claude: red-team/process-skeptic
     gemini: delegation-analysis
   ```

   Role values are free-form lens tags. They are advisory only and do not change quorum, signoff weight, artifact ownership, or drafter eligibility.

2. Add an internal-helper clause.

   A participant may use subagents, retrieval, scratchpads, tools, or other internal mechanisms to produce its own canonical artifact. Those helpers are not Parley Deck participants, do not satisfy the non-solo requirement, do not sign off, and do not own protocol files. The named participant remains fully accountable for its file and signoff.

3. Add a decision framework.

   - Default to 2-4 participants.
   - Use roles/lenses when distinct perspectives materially improve coverage.
   - Add more participants only for cleanly separable modules, review scopes, or competing hypotheses.
   - Avoid multi-agent overhead for sequential same-file work or tightly coupled edits.

4. Add Phase 5 plan-gate guidance.

   Before multi-file implementation work or changes outside `parley-deck/`, the implementer should record a short plan/checklist in `IMPLEMENTATION.md`. Reviewers may use the normal review process to block material divergence from that plan. This is not a new phase, artifact, or automated hook.

5. Keep inbox handoffs non-authoritative.

   Mid-round discoveries and handoffs may use `inbox/`, but substantive decisions must appear in the next round/review file, `consensus.md`, or `FINAL.md`.

### Tests

For a later implementation PR that edits `COOPERATION.md` or CLI behavior:

- Existing protocol parsing still accepts `participants: [id, id]`.
- A prompt with `roles:` renders or validates without changing participant quorum.
- Consensus/signoff logic ignores roles and still requires all listed participants.
- Documentation examples show that internal helper agents do not count as Parley Deck participants.
- Any new Phase 5 plan-gate wording is a checklist/guidance change only; it does not create a new phase or required file.

### Non-goals

- Do not replace Parley Deck with Claude Code Agent Teams.
- Do not add Claude-only protocol requirements.
- Do not add nested sub-ideas, delegated child rounds, task-board state, automatic task claiming, or hook-driven phase transitions in this change.
- Do not treat internal helper/subagent outputs as canonical protocol artifacts.
- Do not edit `COOPERATION.md` until the user explicitly asks for implementation.

### Verification

- Round 01 artifacts exist for `codex`, `claude`, `gemini`, and `hermes`.
- Round 02 artifacts exist for `codex`, `claude`, `gemini`, and `hermes`.
- Consensus reached `reserved` status: `codex`, `claude`, and `gemini` accepted; `hermes` accepted with reservations limited to future `context/` bundle and nested idea proposals.
- The final recommendation is vendor-neutral and does not require Claude Agent Teams.

## References

- Consensus: ./consensus.md
- Rounds: ./round-01/, ./round-02/
- Claude Code Agent Teams docs: https://code.claude.com/docs/en/agent-teams
- Agent Teams Workflow guide: https://github.com/FlorianBruniaux/claude-code-ultimate-guide/blob/main/guide/workflows/agent-teams.md
- Heeki Park Medium article: https://heeki.medium.com/collaborating-with-agents-teams-in-claude-code-f64a465f3c11
