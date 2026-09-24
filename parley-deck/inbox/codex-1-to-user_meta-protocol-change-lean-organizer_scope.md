---
from: codex-1
to: user
idea: meta-protocol-change-lean-organizer
phase: release-remediation
blocking: yes
date: 2026-09-24
---

## Question

Please choose whether to defer the newly exposed native-Windows architecture work with Windows explicitly labelled experimental/unvalidated, or authorize a separate Windows-portability idea and reviewed implementation track.

## Context and concrete proposal

A–D passed independent review and goal-done, with unanimous zero-fix closure. FINAL remains unchanged and drops none of A–D. Post-deploy hosted CI exposed older product failures outside that design: Windows snapshot privacy checks, invalid gate filenames, cross-process locking, and test portability; Linux stderr handling and unresolved captured-child execution failures. The independent findings are release-ci-revalidation-claude-1.md; the concrete participant plan is release-repair-plan-zcode-1.md in the idea. The plan separates small Linux channel repairs from Windows ACL/locking/filename design work and roughly a dozen fixture files. No Windows architecture implementation has begun. Published tags are preserved.

Your original brief requires post-deploy findings fixed before reporting completion. Deferring Windows would amend that release-completion expectation; adding a native-Windows architecture track exceeds A–D. Therefore neither choice is inferred. This is a post-release scope escalation, not a claim that the original FINAL changed or that the pre-Phase-5 scope check failed.

## Work continuing independently

The mutable v1.49.0 release notes now accurately disclose failed Linux/Windows validation and correct the unsupported Windows-coverage sentence. Published assets and tags remain unchanged. Narrow Linux channel repair/diagnostic work can be prepared and independently reviewed without deciding the Windows contract. npm separately awaits owner authentication. CLI winget PR is held while validation remains unresolved; skill winget PR #440360 is open. No all-channels-complete claim is made.

## What I need from you

Choose defer-and-label-experimental versus authorize a separate reviewed Windows track. Also complete npm login as requested in codex-1-to-user_meta-protocol-change-lean-organizer_npm-login.md. No credentials should be sent to an agent.
