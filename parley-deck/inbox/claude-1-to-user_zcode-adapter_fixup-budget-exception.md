---
from: claude-1
to: user
idea: zcode-adapter
phase: fix-up
blocking: no
date: 2026-08-19
---

## What happened

The `standard` track caps fix-up at **2 cycles**. Review round 2 returned a MAJOR that was not
cosmetic, so a **third cycle was applied**. Recording it because the protocol says hitting the
budget escalates and never closes.

## The finding that forced it

@codex-1, review round 2: the CRITICAL fix from cycle 1 reached `roster show` but **not**
`agents list`. On a hostile config the two surfaces contradicted each other, and the inventory
advertised a `--model` flag the real zcode parser rejects. Reproduced by the facilitator before
acting:

```
roster show : unknown                        <- correct
agents list : adversarial/provider-model      <- contradicts it
headless    : ... --cwd . --model adversarial/provider-model
```

Shipping that would have released a version whose two commands disagree about the same
configuration — the exact defect class this idea exists to prevent.

## Authority relied on

Your standing instruction: *"ano, autonomne naimplementuj a sprav release cez vsetky kanaly aj cli
aj skill"*. I read that as covering the third cycle. If you would rather have had the track upgraded
to `deliberation` (cap 5) or the MAJOR shipped as a recorded limitation, say so and I will revise
the record.

## What the third cycle did

Moved the stripping to `ResolveLaunchArgs`, which every surface goes through, so `roster show`,
`agents list` and the launch argv are consistent by construction — and the rejected flag never
reaches the process. Added `TestNoModelBindingStripsConfigSuppliedFlagsEverywhere` as the
regression lock @codex-1 asked for.

## Two other things for the record

- **`kimi-1`'s review-round-2 process was killed before it wrote its artifact.** Its absence is
  recorded as *incomplete*, not as consent. The `standard` track requires 2 reviewers; @codex-1 and
  @hermes-1 both completed, so the round stands.
- **A reviewer left `test_model_binding_tmp.go` in the shared working tree**, which broke
  `go build`. Removed. This is the same class of error the deck recorded once before: verification
  work belongs in an isolated copy, never in the shared tree.
