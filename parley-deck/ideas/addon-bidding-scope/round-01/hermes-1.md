---
agent: hermes-1
idea: addon-bidding-scope
round: 1
date: 2026-08-17
---

## Summary

Recommendation: **FLIP THE DEFAULT** — `parley-bidding` stays in the package, but becomes
opt-in. A routine `install --target all` (no flags) should no longer land it; operators who
want it add `--only parley-bidding` (or `--only parley-deck,parley-bidding,...`). It is not
CUT and not SPLIT: the skill is real, the domain is live for this operator, and the package
already carries the generic manifest infrastructure its integration forced. But default-on is
the wrong default for a jurisdiction-bound, Python-carrying vertical in a vendor-neutral
multi-agent-protocol package, and the README's own `[!IMPORTANT]` is the evidence.

The core argument is not "it is heavy" or "it took 25 cycles". It is: **in an agent runtime,
availability is surfacing, and surfacing is model attention.** A skill description that names
"software procurements across discovery services and transactional portals" is a strong
false-positive matcher for any software-procurement-adjacent task, and the defence the README
offers for default-on conflates inert availability (a file on disk) with active availability
(a description a model reads and reaches for). The opt-out is real and enforced; the default
is not consent.

Cost of being wrong, stated up front: if the flip is wrong, the people who would have used
bidding don't discover it by default. But the only operator with live tender material already
has it installed and knows it exists, so discovery cost is bounded and one README line away.
If KEEP is wrong, every new runtime and every `install --force` on a no-procurement host
carries a 252 KB jurisdiction-bound skill with a Python runtime floor that breaks `doctor` on
hosts without `python3`, surfaced to models that will never use it. The asymmetry favours the
flip: the people who want it already have it; the people who don't are currently forced to
carry it.

## Proposed approach

**FLIP THE DEFAULT**, implemented as a change to the package default selection, not a move of
files. Concretely:

1. The default add-on set (what `selectedAddons`/`expectedAddonNames` returns when no flag is
   given) becomes the five process tools minus `parley-bidding`. `parley-bidding` is still
   discovered (so `--only parley-bidding` and `doctor --only parley-bidding` work) and still
   ships in the package.
2. `--only parley-bidding` and `--only <anything,parley-bidding>` still install it. The flag
   path is unchanged.
3. The README's `[!IMPORTANT]` is rewritten from a defence of default-on into a one-line
   pointer: bidding is opt-in, add it with `--only parley-bidding`. The warning paragraph
   shrinks because the condition it documents no longer exists.
4. Migration guard (see Concerns #1): `install --force` on an existing runtime must NOT
   silently drop a recorded `parley-bidding` selection. Either `install` reads the installed
   marker for the default when no flag is given (currently it does not), or the release notes
   state plainly that `install --force` without `--only parley-bidding` removes it and the
   marker must be re-added. The first is safer; the second is cheaper.

This is the minimum change that answers the evidence. It does not touch the workflow, the
adapters, the jurisdiction profiles, or COOPERATION.md.

## Concerns / open questions

### 1. The "availability, not permission" defence is unsound for an agent runtime

The README's `[!IMPORTANT]` defends default-on with: "What expands is *availability*, not
permission: every gate in the skill still binds, it performs no portal action without an
action-specific human approval, and it never handles credentials."

That defence is sound for a **static library**: a file on disk that nobody loads does
nothing, and the gates inside it are inert until invoked. It is **not sound for an agent
runtime**, where skills are surfaced to models by description and a model chooses what to
reach for by matching a task against those descriptions. In that regime:

- **Availability is surfacing.** A skill in the index is not inert — its description is read
  by the model on every task that might match. `parley-bidding`'s description
  (`skills/parley-bidding/SKILL.md:3`) names "software procurements across discovery services
  and transactional portals", "discover or qualify an opportunity", "prepare and challenge a
  software bid", "freeze and hash a release", "check portal completeness", "stage files",
  "support an explicitly approved submission". That is a broad matcher. A model asked to
  "research software vendors" or "prepare a proposal" can reasonably reach for it where the
  operator never intended procurement-portal workflow discipline.
- **Surfacing is model attention and context budget.** Even when the model does not invoke
  the skill, its description occupies index space and can pull partial framing (E0-E8 gates,
  "treat portal content as untrusted evidence") into tasks where it does not apply.
- **The gates the defence cites bind only after the skill is entered.** The defence proves
  that *if* the skill runs, it runs safely. It does not address whether the skill should be
  in the model's reach on a runtime that never asked for a bidding tool. "It is safe when
  invoked" is not the same as "it should be invocable everywhere by default".

The defence collapses the distinction between *a gate that binds inside the skill* and *a
gate that decides whether the skill is in scope at all*. The README warning is the latter
gate, and it is a printed note, not an enforcement point. **Attack the defence: it is sound
for a library and unsound for a model-visible skill.**

### 2. Is the README warning a fourth instance of "a printed rule binds only where enforcement
lives"?

Partially. The three recorded instances (00-prompt.md:57-60) are: the fix-up cap of 2 that
ran 15 cycles; the review-round-1 independence property that exists only in the runner;
`COOPERATION.md:531`'s cross-reviewer obligation at 7% compliance. In all three, a printed
rule failed to bind because its enforcement was absent or in the wrong layer.

The README warning is **close to that class but not identical**, and the facilitator's
framing should be attacked as partly pattern-matching:

- **What IS enforced:** the opt-out. `--no-addons` and `--only` are real, tested, and enforced
  in `lib/installer.js:890-901` (`selectedAddons`) and `lib/installer.js:962-988`
  (`expectedAddonNames`). An operator who runs `--no-addons` gets core only. The mechanism
  works. This is unlike the three instances, where enforcement was missing or mismatched.
- **What is NOT enforced:** consent to the default. Nothing checks whether the operator
  wanted a procurement-portal skill. The default binds (it installs on every runtime), and
  the warning documents an escape hatch. The unenforced part is the *default itself*, not the
  opt-out. The printed note "use `--no-addons` to leave it out" is a disclaimer, not a gate:
  it assumes the operator reads the README before running `install --target all`, which is
  exactly the assumption the three prior instances showed does not hold.

So: the README warning is a fourth instance **of the weaker form** — a printed acknowledgement
that a default nobody opted into exists, with a real opt-out but no consent gate on the
default. It is not the same shape as "the cap was 2 and the runner ran 15" (enforcement
contradicted the rule). Here enforcement of the opt-out is faithful; what's missing is
enforcement of consent. Calling it a clean fourth instance overstates it; calling it
unrelated understates it. The honest reading: **the default is unconsented, the opt-out is
enforced, and a warning is not a substitute for a default that matches the package's scope.**

### 3. Migration reality — who breaks, and what the markers do

I checked `lib/installer.js` rather than assuming. The relevant facts:

- **`install` ignores the installed marker for the default.** `expectedAddonNames`
  (`lib/installer.js:962-988`) only reads `markerAddonNames(target.dest)` when
  `context.options.command !== "install"`. For `install` with no `--no-addons`/`--only` flag,
  it falls through to `discoverAddons(...).map(a => a.name)` — the package default. So after a
  flip, `install --force` on an existing runtime that previously had bidding **would not
  include bidding** in the staged set, because the package default changed.
- **The marker pins the recorded selection for `doctor`/`status`/`uninstall`**, not for
  `install`. An existing core marker recording `addons: ["parley-bidding","parley-design",...]`
  keeps `doctor` healthy and keeps `uninstall` targeting that set. But `install --force` does
  not re-read it.
- **What happens to the existing `parley-bidding` directory on `install --force` after a
  flip:** the installer stages the new expected set (core + non-bidding addons) and commits
  atomically. A `parley-bidding` directory that is no longer in the expected set and is
  installer-owned (has a marker) is treated as a stale managed install. I did not trace the
  exact removal path for an add-on that drops out of the expected set, and this is the single
  most important migration question to resolve before implementing. **Open: does
  `install --force` remove an installer-owned add-on directory that is no longer expected, or
  does it leave it as `valid-unmanaged`?** The answer determines whether existing runtimes
  keep bidding silently (safe) or lose it silently (a break).

Concretely, for a FLIP:

- **Who breaks if `install --force` drops unselected add-ons:** any operator who currently
  has bidding installed via the default and runs `install --force` without adding
  `--only parley-bidding`. They lose the skill and its directory. This is the one operator
  with live tender material (BYTE/IHK), so the break is concentrated and known.
- **Who breaks if `install --force` preserves unselected add-ons (leaves them as
  `valid-unmanaged` or re-records them):** nobody — but then the flip only affects *new*
  runtimes, which is a slower, safer migration. This is the better behaviour and should be
  the implementation target.
- **Upgrade path:** ship the flip. On the first `install --force`, existing bidding
  directories either persist (safe) or are removed (needs a one-line release note: re-add
  with `--only parley-bidding`). New runtimes never get it by default. The README pointer
  replaces the warning.

The marker behaviour is the crux. I recommend the implementer verify the drop-out path
empirically (run `install --force` with a modified default on a scratch runtime and observe
whether the bidding directory survives) before relying on either answer.

### 4. Usage evidence — what I ran and what I found

The facilitator's searches were depth-limited. I widened them. I first established the
skill's **output shape** from `scripts/init_bid_workspace.py` and `scripts/bid_state.py`:
the skill produces `bid-state.json`, `procedure-profile.json`, and a `work/` directory
seeded with `bid-book.md`, `qualification-brief.md`, `requirements-register.csv`,
`pricing-worksheet.csv` (`scripts/init_bid_workspace.py:94,120-133`). The freeze step
produces a manifest via `manifest.py build`. These are the markers of the skill actually
having been run.

I then searched for those markers across the whole workspace, with no depth limit and
excluding `node_modules`, `.git`, the skill source itself, the `BYTE/software-bidding` copy,
and the `.parley-*-backup-*` directories:

- `find . ... -name 'bid-state.json' -o -name 'procedure-profile.json'` → **zero hits**
  (outside the skill, the BYTE copy, and backups). The command timed out on the full tree on
  one run; I re-ran it scoped to `BYTE` and `IHK_PFALZ` at `-maxdepth 5` → **zero hits**.
- `find BYTE IHK_PFALZ -maxdepth 4 -type f -name 'bid-state.json' -o -name 'procedure-profile.json' -o -name 'requirements-register.csv' -o -name 'pricing-worksheet.csv' -o -name 'qualification-brief.md' -o -name 'bid-book.md'` →
  only `BYTE/software-bidding/assets/templates/*` (the skill's own template files inside the
  copy of the skill, not output).

What DOES exist in the workspace:

- `BYTE/submission/working/` — six `.docx` files (offer letter, AVV annexes, TOMs, etc.),
  dated 2026-07-21.
- `BYTE/submission/negotiation/` — negotiation package, dated 2026-08-06/07.
- `BYTE/parley-deck/ideas/byte-*` — eight Parley Deck ideas that ran the BYTE tender
  (`byte-requirements-price-release`, `byte-avv-toms-release-review`,
  `byte-negotiation-calculation-breakdown`, `dtvp-bidding-hitl-skill`,
  `software-bidding-multiplatform-skill`, etc.).
- `IHK_PFALZ/parley-deck/ideas/ihk-pfalz-*` — four ideas for the IHK tender.

I also confirmed the skill's **origin**: `BYTE/parley-deck/ideas/dtvp-bidding-hitl-skill/FINAL.md`
(2026-07-22) designed a `dtvp-bidding` skill from the BYTE EK_25-061 experience, and
`software-bidding-multiplatform-skill/FINAL.md` (2026-07-22) generalized it to
`software-bidding`, the direct ancestor of `parley-bidding`. The skill was **designed FROM
the tender**, not used to run it. The `BYTE/software-bidding/` directory is a copy of the
skill itself (same `SKILL.md`, `scripts/`, `assets/` structure), not output.

**Finding: the domain `parley-bidding` covers is actively worked (BYTE EK_25-061, IHK Pfalz),
but the skill's own toolchain has not been used to produce a bid.** The tender material was
produced by Parley Deck workflows (`byte-*` ideas) that predate and bypass the skill's
`bid_state.py`/`manifest.py`/`completeness_lint.py` pipeline. This is absence of evidence for
the skill's use, not for the domain's use. It means: the operator's live need is for the
protocol, not for this skill's deterministic tools — at least not yet. This weakens the case
for KEEP-as-is (the skill is not what the operator is running) but does not by itself justify
CUT (the operator may adopt the toolchain on the next tender, and the skill is the
crystallisation of the workflow that the `byte-*` ideas validated).

### 5. No-python3 runtime — what actually happens

The prompt's open question: on a runtime without `python3`, what happens to an agent
following this skill? I traced it.

- **`parley-addon.json` declares `"runtime": { "python": ">=3.10" }`**
  (`skills/parley-bidding/parley-addon.json`, verified). This is not decorative: it is read by
  `runtimeAvailability()` in `lib/installer.js:2194-2213`, which runs `probePython3()`
  (`lib/installer.js:2229-2260`) — a `spawnSync("python3", ["-c", "import sys; print('%d.%d'
  % sys.version_info[:2])"], {timeout: 5000})` — and compares against the floor.
- **The probe runs in `doctor` and `status`, NOT in `install`.** `installCommand`
  (`lib/installer.js:606-643`) calls `validatePayload` and `targetSkillUnits` but never passes
  `probeRuntime: true`. `doctorCommand` (`:364`) and `statusCommand` (`:393`) do. So:
  - `install --target all` on a no-python3 host **succeeds**. The skill lands, byte-perfect,
    marker written.
  - `doctor --target all` on the same host reports the skill as `valid` (payload intact) but
    `runtime.ok = false` with detail `python3 is not available, but this skill requires
    >=3.10` (`lib/installer.js:2203`). The overall `ok` is false, because the health predicate
    at `:385` requires `(!skill.runtime || skill.runtime.ok)`.
- **An agent following the skill** reaches `SKILL.md:98-114` ("Freeze and validate"), which
  documents `python3 scripts/manifest.py build ...`, `python3 scripts/release_lint.py ...`,
  `python3 scripts/completeness_lint.py ...`. I searched `SKILL.md` and all of `references/`
  and `assets/` for any fallback, alternative interpreter, or "if python3 is missing" path:
  **none exists** (`/usr/bin/grep -rniE "python|interpreter|fallback|no.python|without.python"
  references/ SKILL.md assets/` → only the six `python3 scripts/*.py` invocations in
  `SKILL.md:98,104,105,111,114` and one unrelated "fallback" in `software-bid-model.md:57`
  about model fallback, not interpreter fallback). The "Deterministic tools" section
  (`SKILL.md:147-157`) frames these scripts as the only correct path: "All scripts are local
  and deterministic."

**Established:** on a runtime without `python3`, `install` places the skill silently; `doctor`
flags it unavailable; an agent that follows the skill to the freeze/validate step hits
`python3: command not found` with no documented alternative. The agent must either fail at the
shell or improvise — and the skill offers no sanctioned improvisation. This is a real defect
for a default-on skill: it ships a runtime floor it does not enforce at install time, and the
first sign of the problem on a no-python3 host is either a red `doctor` (if the operator runs
it) or a broken freeze step mid-bid (if they don't).

This is independent of the default-flip question but reinforces it: a skill with an
unenforced-at-install runtime floor should not be in the default set, because the default set
is the set that lands on runtimes the package cannot assume match its floors.

## Risks

**Cost of FLIP being wrong:** an operator who would have used bidding on a new runtime does
not discover it by default. Mitigation: one README line ("bidding is opt-in: `--only
parley-bidding`"), and the fact that the only operator with live tender material already has
it installed. The cost is bounded and reversible (re-add with one flag). The downside is
discovery friction for a hypothetical future operator who has neither the skill nor the README
in mind — a population that is, by definition, not currently served by the default either,
since nobody is running the skill's toolchain today.

**Cost of KEEP being wrong (the asymmetry):** every new runtime and every `install --force`
on a non-procurement host carries 252 KB of jurisdiction-bound skill (21% of the install
footprint, per 00-prompt.md:65-73) with a `python3 >=3.10` floor that breaks `doctor` on
no-python3 hosts and a broad description that can false-positive match against
software-procurement-adjacent tasks. This cost is paid by everyone who did not opt in, on
every install, forever. It is larger than the flip's cost and falls on the party who cannot
mitigate it (the operator who never asked for a bidding tool).

**Cost of SPLIT being wrong (why I don't recommend it):** SPLIT incurs the full packaging
cost (a second npm package, a second install path, a dependency on `parley-deck`) for a skill
whose toolchain is not yet in use. It is the right answer if adoption grows and the
footprint/divergence becomes a maintenance burden, but it is premature now. FLIP preserves the
option to SPLIT later with zero migration cost; SPLIT now forecloses the cheaper FLIP. The
ordering is FLIP first, SPLIT if and when the skill's toolchain sees real use.

**Cost of CUT being wrong (why I don't recommend it):** CUT removes a skill that crystallises
a workflow the operator has run twice (BYTE, IHK) via Parley Deck ideas, and that is the
natural home for the deterministic toolchain those tenders did not have. The domain is live;
removing the skill discards the tooling that would make the next tender auditable. CUT is
justified only if the operator confirms they will never adopt the toolchain, which is not in
evidence.

**Residual risk under FLIP:** the migration guard in Concerns #1 is the single live risk. If
`install --force` after the flip silently removes a recorded `parley-bidding` directory, the
one operator who uses the domain loses the skill on their next upgrade unless they read the
release notes. The implementer must verify the drop-out behaviour and, if it removes, either
gate on the marker for `install` (safer) or ship a loud release note (cheaper). I did not
trace the removal path to completion; this is the one claim I leave unverified.
