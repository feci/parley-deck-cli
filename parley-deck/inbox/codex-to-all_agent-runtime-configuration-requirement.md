---
from: codex
to: all
idea: parley-deck-cli-plan
phase: post-merge-product-requirement
blocking: no
date: 2026-05-11
---

## Summary

The user asked to record a project requirement that Parley Deck CLI should expect and recommend effective runtime configuration for cooperating agents.

## Requirement

Codex participation should assume and recommend `sandbox_mode = "workspace-write"` and `approval_policy = "on-failure"` unless the user chooses a different safety profile. Git write capability should be verified directly in the target repository. If a Git write command fails because of sandbox or permission restrictions, the exact failed command should be retried through the runtime approval path rather than silently using an alternate work directory.

The CLI should also guide users when adding different agents to a cooperation run: stable agent ID, CLI path, headless prompt mode, narrow write mode, model/reasoning defaults, timeout, isolated-home needs, and external-backend disclosure should be visible before the agent becomes part of quorum.

## Project note

The user-facing note is captured in `docs/agent-runtime-configuration.md`.
