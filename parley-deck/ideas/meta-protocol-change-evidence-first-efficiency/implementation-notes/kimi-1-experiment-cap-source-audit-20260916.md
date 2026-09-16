# kimi-1 — Experiment cap source audit (bounded)

- Author: kimi-1 · Date: 2026-09-16 · Context: protocol attested `full` (source/packet sha 4519258c…)
- Class: source-only audit. No runtime check, provider call, test, or experiment was executed by me. Tags below: [DOC] documentation capture · [CODE] published source @ 6a2d0db · [PROPOSAL] my unexecuted proposal.
- Note: per-file line numbers cite the retained captures.

## Sources read

1. `.parley-runtime/managed-continuation-20260915/experiment-cap-discovery/kimi-config-files.txt` [DOC] — capture of the advertised page `/kimi-code/en/configuration/config-files` (md form `/kimi-code/en/configuration/config-files.md`, per its own line 16).
2. `…/kimi-env-vars.txt` [DOC] — `/kimi-code/en/configuration/env-vars`.
3. `…/kimi-providers.txt` [DOC] — `/kimi-code/en/configuration/providers`.
4. `…/kimi-docs-1.txt` [DOC] — `/kimi-code/en/customization/agents`.
5. `internal/runner/launch_budget.go` and `internal/budget/monetary_binding.go` [CODE] at integration HEAD 6a2d0db.

## Documented controls, what each bounds, provider/transport scope [DOC]

- `loop_control.max_steps_per_turn` (config-files L241; env `KIMI_LOOP_MAX_STEPS_PER_TURN`, env-vars L124, env wins): max agent steps **per turn**; unset/0 = unlimited. Provider-agnostic. Not per-run, not wall-clock, not USD.
- `loop_control.max_attempts_per_step` (L242; env `KIMI_LOOP_MAX_ATTEMPTS_PER_STEP`, L125): total attempts per failing step, first included; default 10; retries only on transient failures (connection, timeout, 429, 5xx); a quota/balance 429 is **not** retried (L246). Bounds duplicate spend per failing step only.
- `models.*.max_input_size` (L126): declared per-request input limit steering compaction/overflow — soft, not a hard refusal.
- `models.*.max_output_size` (L127) and `KIMI_MODEL_MAX_OUTPUT_SIZE` (env L92): per-request output cap mapped to `max_tokens`; **only the anthropic provider reads it** — it must not be assumed for kimi-code/k3.
- `KIMI_MODEL_MAX_COMPLETION_TOKENS` (env L134): hard per-step `max_completion_tokens` clamp; **kimi provider only** — the documented k3-relevant output bound.
- `[background]` table (config L252–264; env L103–107): `print_max_turns` default 100000, `print_wait_ceiling_s` default 2147483 (~24.8 d), `print_background_mode` default `steer`, `bash_task_timeout_s` print default 0 (no timeout). Bounds print-mode turn multiplication and background concurrency; defaults effectively unbounded. Env-vars page names this table `[task]`; config-files page names it `[background]` — a doc naming discrepancy to resolve against the binary.
- `[subagent]`/`[swarm] timeout_ms` (L265–274): 2 h default each; 0 = no timeout.
- Custom main-agent profile (agents L86–88, `--agent-file` L100–103): `tools:` allowlist, `subagents: []` blocks delegation; each sub-agent independently consumes tokens (L40), so blocking them removes an unbounded multiplier. `[tools]` global allowlist exists too (config L293–303).
- `KIMI_CODE_INFINITE_RETRY` (env L126): must stay unset/false or attempt caps are defeated.

**Composition:** a per-turn step cap is not a per-run cap. In print `steer` mode each background completion injects a new turn (up to 100000 by default), so `max_steps_per_turn` bounds cost only when composed with `print_max_turns`/`print_background_mode=exit|drain` and the external 15-minute wrapper. **No documented control is USD-denominated.**

## Published driver code [CODE]

- `launch_budget.go` L98–113: every non-handoff launch loads the frozen `LaunchBinding` and calls `RequireMonetaryBinding(bound, defaults.MaxCostUSD)`; unreadable defaults refuse (L105–110) — a config error never silently drops the ceiling. An explicit policy must match the frozen history (L114–128); the binding overrides programmatic policy (L130–134).
- `monetary_binding.go` L15–32, L39–56: exact decimal→micro conversion; configured `max_cost_usd` without an attended `budget configure` policy (matching `max-cost-micros` + explicit conservative `reserve-micros`) refuses before execution with `ErrUnknownCost`.
- Scope limits: settlement rounds up a **CLI estimate** — "not an invoice or an invented provider price" (`launch_budget.go` L153–156); the gate covers driver-managed launches for the shared idea scope, not ad-hoc launches, and proves reservation behavior, not provider billing.

## Unknowns still blocking proof of USD15 over 12 tasks × 3 arms + 28 packet/control calls + blind grading/retries [DOC-derived]

1. No USD cap exists in the documented Kimi 0.42.0 config/env surface.
2. k3 pricing/metering (input/output/thinking/cached rates; whether failed or retried requests bill) is absent from the captures — pricing must not be invented; observed multipliers are not hard bounds.
3. Turn counts per task/grading run are model-behavior-dependent; docs cannot prove the exact-28 call count.
4. Server-side quota is a documented failure mode (non-retried 429), not a configurable cap.

## Smallest decisive non-treatment verification [PROPOSAL, unexecuted]

Loopback fake-provider fixture, no external traffic, no roster/model/global change (env vars are process-scoped and outrank config.toml; `KIMI_MODEL_NAME`+`KIMI_MODEL_BASE_URL` synthesize an in-memory provider, env L74–96; loopback bypasses proxy, L170). Run `kimi -p` against a local stub that counts requests and injects canned tool-call/429/500/timeout sequences, with `KIMI_LOOP_MAX_STEPS_PER_TURN=20`, `KIMI_LOOP_MAX_ATTEMPTS_PER_STEP=1`, `KIMI_MODEL_MAX_COMPLETION_TOKENS=<k>`, `KIMI_CODE_BACKGROUND_PRINT_MAX_TURNS=0` (or mode `exit`), `KIMI_CODE_INFINITE_RETRY` unset, and an `--agent-file` profile `tools: [Read]`, `subagents: []`. Assert: ≤20 steps/turn; no second attempt after injected failure; wire `max_completion_tokens == k`; no Agent dispatch; no off-allowlist tool execution. This also settles the `[task]` vs `[background]` key naming for the installed 0.42.0 binary.

**It cannot prove** actual managed-service billing, k3 pricing, the real failure/retry mix, or real turn counts for the full 12×3 + 28 + grading denominators; it is not a substitute for the pilot or the packet trial.

## Hard-bound route

Requests-per-run hard bound: viable from documented env composition under the preconditions above. USD15 hard bound: **not source-provable today** — the missing proof is k3 token pricing and server metering semantics, plus provider-verified settlement; the driver monetary binding adds a pre-execution gate but settles on estimates, so it complements rather than completes billing proof. Gates stay open.

No acceptance or completion verdict is issued here.
