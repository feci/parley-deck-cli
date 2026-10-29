# Organizer brief: meta-protocol-change-credit-exhaustion-failover

You are **codex-1, the organizer and the implementer** of a new Parley Deck idea. The owner wants Claude
tokens saved: codex-1 organizes and implements, the owner's Claude Code session only relays, and claude-1
takes no part.

## Workspace

- CLI worktree: `/Users/tomasfecko/AI_WORKSPACE/parley-deck/worktrees/credit-failover` (branch
  `credit-failover`, from origin/main `525febb`, CLI 1.55.0). The deck is its `parley-deck/`.
- Skill worktree: `/Users/tomasfecko/AI_WORKSPACE/parley-deck/worktrees/credit-failover-skill` (branch
  `credit-failover`, from origin/main `f557c5e`, skill 2.19.0).
- **Path change.** The workspace is now at `/Users/tomasfecko/AI_WORKSPACE/`. The old path
  `/Volumes/My Shared Files/AI_WORKSPACE/` does not exist on this machine. Old worktree registrations
  that use the old path show as `prunable`. Do not prune them, and do not declare worktrees.
- `codex` is at `/Users/tomasfecko/.nvm/versions/node/v24.15.0/bin/codex`. `kimi` is at
  `~/.kimi-code/bin/kimi` (not on PATH). `zcode` is on PATH.

## The owner's words (verbatim, Slovak) and translation

> "podme navrhnut update protokolu, ktory vynecha automaticky agenta, ktory ma model ktoremu providerovi
> dosli kredity, bez toho aby ho skusal znova pripojit, ak nejakemu agentovi dojdu tokeny alebo kredity, nech
> je hned vyradeny a dalsi v poradi preberie jeho ulohu, nastartuj /parley-deck a sprav novu verziu, pricom na
> zaver to aj vsade pushni a deployni cez vsetky kanaly a nainstaluj lokalne do vsetkych agentov"

Translation: "Let us design a protocol update that automatically skips an agent whose model's provider has
run out of credits, without trying to connect it again. If an agent runs out of tokens or credits, it is
excluded at once, and the next one in order takes over its task. Start /parley-deck and make a new version.
At the end, push and deploy it on all channels, and install it locally into all agents."

## Current state (verify at HEAD, §15)

Read core 2.19.0 §9.0 "Automatic exclusion" (the CLI worktree's `parley-deck/COOPERATION.md`) and the
prior ideas. Do not redo their work:

- `meta-protocol-change-quota-auto-exclude` (1.51.0): the legacy quota predicate. It is strict: HTTP 429,
  explicit "Limit Exhausted", a machine reset of 60 minutes or more, and native provider provenance per
  adapter.
- `meta-protocol-change-participant-dropout` (1.52.0): any failure of a non-protected participant gets two
  attempts (the original, plus one retry after 5 s). The second failure drops it for the rest of the idea.
- `meta-protocol-change-participant-dropout-review-gate-timing` (1.53.0): the single-reviewer exception.
- `meta-protocol-change-driver-unstall` (1.54.0). The owner applied the D6 legacy declaration, so the
  round-02 driver refusal must not occur now. If it occurs, record it as a driver gap.
- The open follow-ups `quota-kickoff-reporting-and-alias-guidance` and
  `quota-zcode-native-exhaustion-capture`. Include the parts that this idea needs. List the rest as
  deferred.

The gaps against the owner's words:

1. A credit or token exhaustion gets a retry today (participant-failure path), or it does not qualify at all
   (legacy predicate with a short or missing reset). The owner wants no retry and no reconnection.
2. Protected roles (facilitator, implementer designee or pin, started drafters) **block** the batch when
   they fail, and "Organizer failover or replacement agents stay outside this rule". The owner wants the
   next agent in order to take over the role.

## Design target (participants decide the details in FINAL)

1. **Exhaustion class and no retry.** Define one "provider exhaustion" class: the provider or gateway
   states that the account's credits, tokens, quota or usage allowance are exhausted. Use the recorded
   evidence (exit, structured provider error, stderr record) and the adapters that exist.
   - On this class: drop the agent at once in the current idea. Do not retry, do not probe again, and do
     not reconnect.
   - Transient failures (bare 429, 5xx, overloaded, timeout, hang) keep the shipped two-attempt path.
   - Decide the boundary for a short stated reset (an observed gateway text: "All claude accounts have
     exhausted their quota (reset after 5m)"). Decide the role of the 60-minute threshold. The owner intent
     is: an agent out of credits or tokens is not retried.
   - Keep false positives low. Quoted prose or an agent's own text is never provider evidence.
2. **Provider scope and "no reconnect".** The owner says "the agent whose model's provider ran out of
   credits". Decide:
   - whether one exhaustion record also excludes other agents of the same idea that use the same provider
     account or gateway connection;
   - whether a small machine-level record (provider or connection, reset time, evidence) lets later
     kickoffs and preflights skip the agent without a probe until the reset, and probe it afresh after the
     reset. If you add it, keep it small and auditable, and never store secrets in it.
3. **Role succession (owner ruling, binding).** Read
   `inbox/user-to-codex-1_meta-protocol-change-credit-exhaustion-failover_succession-order.md` and quote it in
   `00-prompt.md`. The owner set this order:
   - The organizer drops out: the implementer takes over the organizer role.
   - The implementer drops out: the organizer takes over the implementer role.
   - Both drop out: the first remaining participant in `participants:` order takes over both roles.
   - A participant drops out: no succession. It only drops out.

   Decide only what the ruling leaves open:
   - one agent that holds both roles (as in this idea) counts as "both drop out";
   - the rule "a declared facilitator does not implement" yields to the ruling;
   - the consensus, FINAL and review-consensus drafter roles;
   - what the CLI driver (`parley run`, `continue`) enforces, and what stays a protocol rule for a manual or
     external organizer process;
   - **the handover.** The successor continues from the canonical artifacts and the branch state. It
     records the takeover in `IMPLEMENTATION.md` (or the phase artifact) and in one notice. No artifact of
     the dropped agent is rewritten. Filed dissent, ❌ and findings of the dropped agent stay binding;
   - **independence.** A successor implementer never reviews its own code. Its earlier findings stay. Say
     how the reviewer count, the floor and the single-reviewer exception read after a succession.
4. **Floor and escalation stay.** Two distinct usable non-facilitators must remain. Below the floor, or when
   no eligible successor exists, escalate once with a blocking owner note. The drop of a dissenter must stay
   impossible to use as a way to remove dissent.
5. **Size.** The core skill `SKILL.md` is 19,482 of 20,000 bytes. The unflagged phase packet is about
   68,935 of 70,000 bytes, and the meta-idea packet variant is 71,934 bytes (already over). Do not raise
   the limits. Make room by compaction without the loss of a rule. Do not make the meta-idea variant worse.
   List the exact protocol hunks for all three COOPERATION.md copies, plus the CLI files, the skill files
   and the tests.
6. **Writing standard.** Write the new protocol prose and all artifacts in ASD-STE100 (Appendix B).

## Quorum and roles

- `00-prompt.md`: `author: user`, no `facilitator:` field (codex-1 facilitates as an ordinary participant),
  `participants: [codex-1, kimi-1, zcode-1]`, `implementer: codex-1`, `track: deliberation` (a protocol
  change), `auto_implement: true`, `require_model_diversity: true`. Quote the owner.
- Run `parley preflight`. The relay probed all three at about 19:10 CET on 2026-10-29: codex-1, kimi-1 and zcode-1
  answered PONG. kimi-1 works again (machine roster model `omniroute/k3.1-pro`). If one fails, the shipped
  dropout rules apply. Never touch credentials, keys or gateway settings.
- kimi-1 and zcode-1 are the independent reviewers and goal checkers, in separate processes with focused
  briefs. You never review your own code.
- **Close authority.** Normal gates apply: both reviewers. If one reviewer drops, the shipped single-reviewer
  exception applies. If it does not apply mechanically, the owner pre-authorizes the close when all of these
  hold:
  - the final review of the remaining reviewer has no open CRITICAL or MAJOR;
  - all current participants filed a review-consensus ACCEPT;
  - a fresh goal check by a reviewer process passed;
  - the current-tree AC evidence is recorded.

## Process, retries and release

- **No proposal stop.** The owner said: make the new version and release it at the end. Go from design
  through implementation to release without stopping. Stop only if an owner decision is truly open. In that
  case, write one blocking `codex-1-to-user_<slug>_decision.md` with the options and a recommendation.
- **Fix-up.** At most 5 cycles. If a CRITICAL or MAJOR is still unresolved after cycle 5, escalate.
- **Retries.**
  - A provider 429 or 503, or a timeout: wait 15 minutes and relaunch, at most 8 times per step.
  - A timeout: relaunch with 2400 s, then with 3600 s.
  - An auth or credit error on codex-1: stop with a blocking note.
- **Driver.** Use the driver (`parley run`, `continue`, `wait`, `status`) and record each gap. Never prune or
  declare worktrees.
- **Waiting.** Never end your turn while you wait. Wait only with blocking tool calls. End only when the
  work is done, or deliberately after a blocking owner note.
- **Conventions.** English everywhere, in STE. Commit prefix
  `[codex-1] meta-protocol-change-credit-exhaustion-failover: ...`. Never copy secrets. Never drive Google
  Chrome.
- **Release.** Merge to `main`. Then:
  - Versions: CLI **1.56.0**, skill **2.20.0**, core **2.20.0**. Build core 2.20.0 on the staged
    `~/.parley/staging/COOPERATION-2.19.0.md` plus exactly this idea's reviewed hunks. Stage it as
    `~/.parley/staging/COOPERATION-2.20.0.md`.
  - GitHub releases.
  - Both Homebrew formulae.
  - A skill-only winget PR (one application per PR). The CLI winget stays held while Windows is
    experimental.
  - Install the skill into all managed runtimes **and** the four generic locations
    (`~/.hermes/profiles/{ldx,librade,testprofile}/skills/parley-deck` and
    `~/.config/opencode/skills/parley-deck`) with
    `parley-deck-skill install --target generic --dest <dir> --force`. Find every copy with
    `find ~ -maxdepth 7 -path "*/skills/parley-*/SKILL.md"` and verify every SKILL.md by hash.
  - npm publication and `parley protocol publish` are owner-only (npm login with OTP, and a TTY gate). Do
    not try them, and do not work around the gates. The relay updates the owner script.
  - A reviewer process (zcode-1 or kimi-1) verifies all channels in a separate process.
- **Finish.** Write `inbox/codex-1-to-user_meta-protocol-change-credit-exhaustion-failover_released.md`
  with:
  - the channel evidence;
  - the still-pending owner commands, in version order: npm and core for each still-unpublished version
    from 2.15.0 to 2.20.0;
  - the deferred items;
  - usage.

  Then exit.
