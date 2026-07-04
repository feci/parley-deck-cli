---
agent: antigravity-1
idea: parley-learn-playbooks
round: 1
date: 2026-07-04
---

## Summary

Distilled playbooks represent a critical bridge for "consumer decks" (projects using Parley for internal governance/process) between "what was decided" and "how we execute it again." For a developer running Parley in their own project, these playbooks serve as a project-specific memory. However, their value is strictly tied to their freshness; without a strong lifecycle management mechanism, they risk becoming "process debt"—stale markdown files that agents and humans eventually ignore.

## Proposed approach

- **Actionable Skeletons**: Playbooks should go beyond narrative summaries and provide importable fragments (e.g., a `checklist.md` or a Phase 1 plan template). A developer should be able to "instantiate" a playbook directly into a new idea.
- **Contextual Backlinking**: Every distilled playbook must link back to the specific closed idea(s) it originated from. This allows future participants to "drill down" into the original evidence and discussion if a playbook step seems confusing or outdated.
- **Discovery Integration**: Playbook suggestions should be integrated into the `parley new` or `Phase 0` flow. Using simple keyword or similarity matching on the idea brief can prompt the facilitator: "I found a playbook for 'release-burst'; would you like to use it?"
- **Lifecycle Gates**: The §13 retro should not just *propose* playbooks but also *audit* them. If a playbook was ignored by three subsequent ideas of the same topic, it should be flagged for revision or archival.

## Concerns / open questions

- **Maintenance Overhead**: Who owns the playbook? If it's "advisory" and "non-canonical," there is a risk no one feels responsible for updating it when the underlying project structure changes.
- **Protocol Drift**: How do we handle playbooks created under older versions of the Parley protocol? A playbook that expects a specific phase structure might become technically "broken" if the protocol is updated.
- **Granularity**: What is the "Goldilocks" size for a playbook? Too broad (e.g., "Feature Work") and it's useless; too narrow (e.g., "Fixing CSS in Login Page") and it's clutter.

## Risks

- **Stale Clutter (Process Graveyard)**: The most significant risk to consumer adoption. If the `playbooks/` directory becomes a list of 50+ outdated files, it creates mental friction for the developer and noise for the LLM.
- **Process Calcification**: Playbooks might discourage innovation. A developer might follow a "proven" playbook instead of proposing a more efficient but "un-playbooked" process.
- **Context Poisoning**: Injecting too many advisory playbooks into an agent's context can lead to "instruction fatigue," where the agent prioritizes old patterns over the specific requirements of the current brief.
