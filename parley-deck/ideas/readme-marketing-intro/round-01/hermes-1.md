---
agent: hermes-1
idea: readme-marketing-intro
round: 1
date: 2026-06-20
---

## Summary

My lens is positioning and differentiation. The current README opens with a
flat "local CLI for Parley Deck multi-agent cooperation" line that answers
*what* before a newcomer knows *why they should care* or *why this instead of
just spawning three agents in three terminals*. The intro's job is to land the
one comparison a technical reader is already making in their head — "isn't this
just ad-hoc multi-agent, or worse, one model wearing committee costumes?" — and
answer it in the first screenful, with claims that map to real protocol sections.

The core hook: **Parley Deck is the protocol the ad-hoc approach reinvents badly,
every time.** Spawning N agents by hand gives you no audit trail, no conflict
discipline, no consensus mechanism, no resume, no supervision, and no way to
trust the result. Asking one model to roleplay a committee is the opposite of
multi-agent — it's solo reasoning in costume, and it reintroduces exactly the
single-model self-preference that real quorum-gated review exists to defeat.

## Proposed approach

### Framing (the story the intro tells)

Lead with the contrast, then earn it with three concrete "things you'd have to
build yourself" that ship in the protocol. Structure, roughly:

1. **One-line hook** — name the category and the anti-pattern in the same breath.
2. **The two things it replaces** — ad-hoc multi-agent + one-model-as-committee — each killed in one bullet with a real section citation.
3. **What you get instead** — a compact four-to-five bullet list of durable protocol guarantees (audit trail, conflict avoidance, consensus, supervision, transport-agnostic), each citable.
4. **One-line closer** pointing at the 8-phase lifecycle and the `parley` CLI as the thing that makes it usable, not just specified.

### Concrete intro bullets (each maps to a real section — no invention)

- **Real multi-agent, or it doesn't count (§1).** Stable agent IDs, per-agent canonical artifacts, one-file-per-agent-per-round (§6) — so agents never stomp each other and every voice is recoverable. This is the line that separates Parley Deck from "one model with three personas."
- **A protocol, not a prompt.** The 8-phase idea lifecycle (§4) — kickoff → independent analysis → cross-review → consensus → `FINAL.md` → `IMPLEMENTATION.md` → code review → fix-up — is a durable, append-only audit trail you can resume from the docs alone.
- **Compare, don't merge.** The consensus "Comparison & blind spots" lens (adopted from OpenRouter Fusion) rates confidence by agreement and surfaces contradictions and blind spots instead of averaging them away. One model role-playing a committee can't do this — it has no second voice to disagree.
- **Supervised, not fire-and-forget.** First-output watchdog, stall guard, failure classification, and "validated artifact beats nonzero exit code" (§agent supervision) mean a hung agent doesn't silently poison the round.
- **Transport-agnostic (§0, §11).** `local-dir`, `github-pr`, `gitlab-mr` — same protocol whether agents share a filesystem or review each other through a PR.
- **Belt and suspenders.** §9.0 `parley preflight` pings roster liveness before an idea starts; §13 `parley retro` runs an advisory, quorum-gated pass over the deck's own history — quorum, not single-model self-preference (the RHO twist).

### Draft intro text for README.md (drop-in, before the existing `## Install`)

> # parley-deck-cli
>
> `parley` runs **Parley Deck** — a transport-agnostic protocol for real
> multi-agent cooperation, with a CLI that makes it usable instead of just
> specified.
>
> **Why not just spawn three agents in three terminals?** Because ad-hoc
> multi-agent gives you no audit trail, no conflict discipline, no consensus
> step, and no way to resume. Parley Deck ships those as protocol: an 8-phase
> idea lifecycle (§4) with an append-only trail, one-file-per-agent-per-round
> conflict avoidance (§6), and a compare-not-merge consensus lens that rates
> confidence by agreement and surfaces blind spots instead of averaging them
> away.
>
> **Why not ask one model to role-play a committee?** Because that's solo
> reasoning in a costume — there is no second voice to actually disagree, and
> it reintroduces the single-model self-preference that quorum-gated review
> exists to defeat. Parley Deck is non-solo by design (§1): stable agent IDs,
> per-agent canonical artifacts, any CLI/model on the roster (§2).
>
> What you get on top: `parley preflight` (§9.0) pings the roster before an
> idea starts; agent supervision (watchdog, stall guard, validated-artifact-
> beats-exit-code) catches hung agents; `parley retro` (§13) runs an advisory,
> quorum-gated pass over the deck's own history. Runs over `local-dir`,
> `github-pr`, or `gitlab-mr` (§0, §11).
>
> Inspired by — and adapted from — OpenRouter Fusion (consensus lens),
> OpenAI ExecPlans / PLANS.md (`FINAL.md` + living `IMPLEMENTATION.md`), RHO
> (advisory, quorum-gated retro), and kindly (review gates, supervision
> knobs, artifact-wins). We didn't invent these; we wired them into one
> protocol with a CLI.

Keep it to roughly that length — one screenful, punchy, every claim citable.
The Install section follows immediately after.

## Concerns / open questions

- **How much protocol jargon in the intro?** I've leaned on `§`-citations to stay truthful, but a first-time reader doesn't know what §4 is. Decision needed: keep citations as honesty markers (and let the protocol doc own the detail), or strip them and trust the bullets. My vote: keep them — they signal "this is spec'd, not vaporware" and cost one character each.
- **Does the adapted-for-skill README drop the transport-agnostic line?** The skill version may not ship `github-pr`/`gitlab-mr` transports the same way. The "transport-agnostic" bullet is a strong differentiator for the CLI; for the skill it may need softening or a different anchor (the 8-phase lifecycle + consensus lens travel fine).
- **Positioning vs. "workflow framework" framing.** I've pitched this as "the protocol ad-hoc reinvents badly." An alternative is "multi-agent with a durable audit trail." The former is sharper and more confrontational; the latter is safer. I chose sharp — confirm that's the brand voice, or dial it back.
- **Antigravity-1's attribution pass will need to validate** that I've credited inspirations as adopted-and-adapted, not claimed invention. The closing line is designed to hand off cleanly to that lens.

## Risks

- **Over-claiming on supervision.** "Catches hung agents" is strong; if the watchdog/stall guard has known gaps, this bullet oversells. Codex-1's technical-accuracy pass should verify the supervision features map to shipped behavior, not just spec text.
- **The "one model in a costume" line is provocative.** It's accurate and it's the sharpest differentiation we have, but it may read as throwing shade at a common workflow. If the project wants a friendlier tone, this is the first line to soften — but I'd argue the edge is what makes the hook land.
- **Bullet count creep.** Six bullets is already at the edge of "tight intro." If claude-1's narrative pass or antigravity-1's attribution pass each want to add one, we risk a wall of text. The intro should stay ≤ one screenful; anything that doesn't fit moves to a later section.
- **"Compare, don't merge" needs the Fusion credit visible.** If the intro uses the consensus lens as a selling point but buries the OpenRouter Fusion attribution in the closer, we risk the exact "we invented this" failure mode the prompt warns against. The credit line must stay in the intro, not get pushed to a footer.
