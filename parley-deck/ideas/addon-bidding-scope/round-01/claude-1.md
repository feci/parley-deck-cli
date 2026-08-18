---
agent: claude-1
idea: addon-bidding-scope
round: 1
date: 2026-08-17
---

## Summary

**FLIP THE DEFAULT.** Keep the add-on, keep it in this package, make it opt-in. Rank order:
FLIP > SPLIT > KEEP > CUT.

The decisive argument is not size, not jurisdiction, and not the 25 fix-up cycles. It is that
**the README's own defence is wrong about the thing it is defending.**

> *"What expands is availability, not permission."*

For a human tool that distinction holds. For an agent runtime it does not. Skills are surfaced to
models **by description**, and a model chooses what to reach for by matching a task against those
descriptions. Availability *is* behavioural surface. `parley-bidding`'s description is 120 words
naming tenders, suppliers, pricing, portals, submissions, amendments and withdrawals — sitting in the
skill listing of every agent this installer touches, including on people's machines who never asked
for a procurement tool.

On a machine that never bids, its expected value is zero and its variance is not. That asymmetry is
the whole case, and it does not require anything to have gone wrong yet.

## Proposed approach

### FLIP THE DEFAULT — recommended

`--no-addons` and `--only` already exist, so the mechanism is built. The change is which side of
the line this one add-on starts on.

**What it costs:** nothing to build. **Who breaks:** nobody — existing installs keep what they have;
the installed marker already records which add-ons were installed and from which package version, so
an upgrade can leave a previously-installed `parley-bidding` in place while stopping new default
installs. That is the same marker mechanism the manifest fork produced, used for a second purpose.
**Upgrade path:** `install --only parley-bidding` for anyone who wants it.

**Why it is the right size of change:** it aligns the default with a concern the package *already
documents*, discards none of the 25 cycles of work, and is reversible in one line if the argument is
wrong.

**Cost of being wrong:** someone who would have benefited from bidding never discovers it. That is a
discoverability cost, and it is recoverable — the README entry stays, the skill stays, one flag adds
it. Compare with the failure it prevents, which is silent and lands on third parties.

### SPLIT — correct in principle, not yet worth its price

A DE-jurisdiction procurement vertical with a Python runtime genuinely is a different product from
"a protocol for making agents deliberate". If this package ever ships to third parties at scale,
splitting is where it ends up.

I do not propose it now for one reason, and it is a repository fact rather than a preference: this
deck's standing release discipline is that **every release goes through all channels** — npm,
Homebrew, winget — with an independent verifier per channel. A second package doubles that tax
permanently. The one thing bidding demonstrably bought the package (the generic manifest mechanism,
`lib/installer.js:2388`, now used by all six skills) is already banked and does not depend on
bidding staying here.

Revisit if the flip proves insufficient, or when external adoption makes the vertical's presence a
credibility question rather than a footprint question.

### KEEP — the steelman, and why I still reject it

Honestly stated: every gate in the skill binds, it handles no credentials, and it performs no portal
action without action-specific human approval. Nothing has gone wrong. The install is 252 KB. On
those facts, "leave it" is defensible and I want that recorded rather than strawmanned.

I reject it because the defence answers a question nobody asked. The risk was never that the skill
would submit a tender by itself. It is that a procurement vertical is *in the listing* on machines
whose owners never chose it — and the package's answer to that is a README paragraph. **A README
warning is not a gate**, and this deck has now recorded three instances of that exact class. Whether
this is the fourth is the question I would most like the other three to attack, because I am the one
who put the framing in the prompt and I may be pattern-matching.

### CUT — no

It works, it is paid for, the domain material exists in this workspace (`BYTE/`, `IHK_PFALZ/`, both
with their own decks), and no measured harm has occurred. Removing a working, gated, human-controlled
tool because it is unlike its siblings is tidiness, not engineering.

## Concerns / open questions

1. **The Python dependency is a separate defect from the default, and I could not find its
   fallback.** PRIMARY: `SKILL.md:98-114` invokes `python3 scripts/manifest.py`,
   `release_lint.py`, `completeness_lint.py` in the documented workflow. A grep across `SKILL.md`
   and every file in `references/` for a no-interpreter path returned nothing relevant. So: on a
   runtime without `python3`, what does an agent following this skill do — fail loudly, or fail
   confusingly? I did not establish this and it should not be assumed. If it fails confusingly,
   that is worth fixing **whichever** way the default goes.
2. **Two of six add-ons advertise "zero runtime dependencies" as a selling point.** A package where
   one member carries an interpreter requirement makes that claim harder to state cleanly. This
   argues mildly for SPLIT and not at all for CUT.
3. **My usage evidence is weak and I want it treated as weak.** My searches were depth-limited
   (`-maxdepth 5–6`) and this shell's `grep` honours `.gitignore`. "No artifact in this skill's own
   shape was found" is absence of evidence. If someone wants to argue KEEP on the grounds that it is
   used, widen the search — `find` without depth limits, or `rg -uuu` — and say what you ran.
4. **Does the flip actually reduce the surface, or just move it?** If an agent's skill listing is
   built from what is *installed*, the flip works. If some runtimes discover skills from the package
   rather than the install, it does not. I have not verified how each of the seven runtime targets
   builds its listing, and the argument depends on it.

## Risks

- **Framing capture.** I wrote the prompt, including the "fourth instance of a printed rule" framing
  and the 71%-of-payload table. Both push toward change. The check on me is whether anyone reaches
  KEEP independently; if all four of us land on FLIP, that is a shared prior and `consensus.md`
  should say so under `COOPERATION.md:1356-1360`.
- **Solving a hypothetical.** Nothing has gone wrong. The whole case rests on an asymmetry argument,
  not an incident. That is a weaker footing than the last two ideas had, and I would rather name it
  than let it pass as measured.
- **Sunk cost in both directions.** 25 fix-up cycles is a reason to neither keep nor remove. If any
  argument in this round leans on how much work it was, it should be struck.
- **Scope creep into `parley-design-check`.** It is 615 KB — half the payload, more than double
  bidding — and it is explicitly out of scope here. If the answer to bidding is "flip the default",
  the obvious next question is whether the same applies to a 615 KB enforcement engine. That is its
  own recon and its own idea; do not settle it here by implication.
