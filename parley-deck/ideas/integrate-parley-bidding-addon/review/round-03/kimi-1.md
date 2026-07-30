---
idea: integrate-parley-bidding-addon
review-round: 3
agent: kimi-1
date: 2026-07-30
reviewed-commit: 5c324ef
---

## Verdict
ACCEPT WITH CONDITIONS

All four round-2 fixes are closed by my own re-measurement, the claimed counts reproduce on
both interpreters, and my only new findings are two NITs. The single condition is the
skills-CLI ruling below: implement **(b)** — a third verdict for a markerless tree whose
manifest fully verifies — and revise the `CHANGELOG.md` entry that currently documents (a)
as intended. The NITs may ride along.

## Round-2 findings — closed or not

### F4 (codex-1 MAJOR) — probe answered for the wrong environment — CLOSED
Re-measured with my own harness (fresh `install --target claude --yes` into a controlled
`HOME`, then `run(["doctor", "--json"], { env })` twice **in one process**, both orders):

| order | first env | second env | result |
|---|---|---|---|
| empty → real | `PATH: ""` | real PATH (python3 3.14.6) | first: `runtime.ok:false`, "python3 is not available", exit 1; second: `runtime.ok:true`, "python3 3.14", exit 0 |
| real → empty | real PATH | `PATH: ""` | first: `ok:true`; second: `ok:false`, exit 1 |

Both orders give the per-environment answer, so the memoization is genuinely keyed per PATH —
a process-global cache would have reused the first answer in one of the two directions. The
parent process had a working 3.14.6 throughout, so the empty-PATH `ok:false` cannot be
anything but the caller's environment being honored. Their own regression test asserts the
same property and passes.

### F5 (hermes-1 MAJOR) — environment-dependent test — CLOSED
Measured in a clone at `5c324ef`:

- `node --test` with Homebrew python3 3.14.6 first in PATH: **290 pass, 0 fail**.
- `node --test` with `/usr/bin/python3` 3.9.6 first in PATH: **290 pass, 0 fail**.
- Full `npm test` on 3.14.6: 290 node + 54 Python OK across 7 files + manifest check
  `parley-bidding: ok (47 files, sha256:7854adf1…)` — the same digest as round 2.

One precision on the claim "290/290 with 3.9.6": that is the **node** count. The Python leg
on 3.9.6 exits 1 by design ("python3 is 3.9, but the add-on declares >=3.10" — a below-floor
interpreter fails rather than skips, which is the intended contract and matches the runner's
own comment). So full `npm test` green still requires a ≥3.10 interpreter, as it should; the
node suite no longer depends on the ambient one at all.

### F6 (codex-1 MINOR) — probing leaked into `paths`, `status` discarded it — CLOSED
My own sentinel measurement (stub `python3` that appends to a sentinel file and prints
`3.12`, only PATH entry): `paths` left the sentinel **absent** and reported `runtime: null`;
`doctor` in the same environment tripped the sentinel and reported `python3 3.12`. `status`
now prints the `unavailable:` line under each affected add-on. Call-site audit: exactly
three `targetStatus` callers — `paths` (`probeRuntime:false`), `doctor` (probe + `env`),
`status` (probe + `env`). No fourth path can probe.

### F7 (kimi-1 NIT) — stderr wording — CLOSED
Measured human-mode stderr with an availability-only failure:
`"One or more installs are operationally unavailable.\n"` — the collision is gone.

### codex-1's transient 4-test failure
Not reproduced in any of my runs (two full suites on two interpreters plus a dozen probe
processes). I concur with recording it as unexplained rather than closed; the per-PATH probe
keying removes the one plausible order-dependence mechanism either way.

### Answers to my round-2 open questions — both verified
- **Q1 (Windows, `python3`-only).** The contract is now stated in `CHANGELOG.md` ("The probe
  looks for **`python3` specifically**… On a Windows host where only `python` exists… the
  add-on is reported unavailable. That is the fail-safe direction"). Documentation, which is
  what I asked for; verified present.
- **Q2 (legacy carve-out invisible).** The claimed property is real, measured: a schema-2
  unit's `doctor --json` entry carries `marker.markerSchema: 2` plus the manifest anchor;
  after I rewrote the marker to the 2.0.0 shape (no `markerSchema`, no `manifest` field) the
  same unit reported `valid`, exit 0, with `marker.markerSchema` absent. Grandfathered is
  machine-distinguishable from anchored with no new field. Correct answer, no code change
  needed.

## Ruling on the skills-CLI question

**(b)** — with one precision about which marker states qualify, and one scope note the
option's wording does not cover.

I reproduced the scenario and went one step past it. The published GitHub repo has no
`parley-bidding` (the CLI finds 5 skills there), so the measurement only exists against the
branch — `skills add /tmp/…/repo --skill parley-bidding --agent claude-code --yes`, which is
the faithful form of the README command at this commit. Measured:

- the skills CLI **copies** (does not symlink) the tree via a `.agents/skills` staging step;
  `diff -r` against `skills/parley-bidding` is byte-identical, **including
  `parley-addon.json`**; no `.parley-deck-skill-install.json` anywhere;
- `verifyPayload` on that tree: `{ok: true, problems: []}` — provably byte-intact;
- `doctor --target claude`: `claude/parley-bidding: malformed … no parley-deck-skill install
  marker`, exit 1 — the reported behavior, confirmed;
- the realistic end state is worse than the single-skill repro: a user who follows the
  README's first recommendation fully (`--skill '*'`, all six skills) gets **every unit
  `malformed`, exit 1**, on six byte-perfect payloads.

Why (b) and not (a):

1. **`malformed` is false there, and the tool itself proves it.** The same idea ships the
   manifest mechanism whose whole purpose is byte-integrity evidence. When `doctor`'s
   verdict contradicts its own strongest evidence, users learn that `malformed` can mean
   "perfectly fine" — and the integrity word stops meaning anything. That is the
   untested-and-green failure mode this mechanism exists to kill, inverted into
   tested-and-red.
2. **(a) makes the repository contradict itself at the exact point both artifacts are
   exercised.** `README.md:143-151` recommends the universal installer *first*, this
   package's own installer second. Keeping (a) coherent would mean demoting the skills-CLI
   path in the README to protect a code behavior — a docs retreat, backwards. Note the
   `CHANGELOG.md` "Changed" entry already documents (a) as intended ("This includes trees
   copied by a third-party skill installer… `doctor` reports them as not installed by this
   tool"); under (b) that entry must be revised, not just the code.
3. **The security delta of (b) is not a real boundary.** The one case that changes: payload
   and manifest consistently rewritten *and* the marker deleted moves from `malformed` to
   the third verdict. But the marker is unsigned and writable — an attacker who can rewrite
   the tree can equally write a schema-2 marker with a recomputed manifest hash, which
   passes under (a) too ("if the marker is writable, so is the check", as the implementer's
   own notes concede). The marker's real value is anchoring post-install drift on
   tool-managed trees and distinguishing managed from foreign — and under (b) that signal
   *survives*: the third verdict plus `marker: null` still tells an auditor "this is not
   the tree this tool installed".
4. **Round-1 guarantees are preserved.** Marker absent + manifest present but failing
   (gutted tree) → still `malformed`. Marker absent + manifest absent where the source
   ships one (double deletion) → still `malformed`. The legacy-marker carve-out is
   untouched.

Requirements for the implementation (part of the ruling, not flavor):

- The third verdict applies **only when the marker is entirely absent** and the manifest
  fully verifies. A present-but-unreadable marker, or one naming another installer, stays
  `malformed` — those are corruption/tampering of management metadata, not "never
  installed by this tool".
- A status string distinct from both `valid` and `malformed` (e.g. `valid-foreign`), in
  text and in `doctor --json`, so automation can still require tool-managed installs where
  it cares.
- It does **not** fail health: `doctor` exit 0 when that is the only finding. I could not
  construct a case where a byte-verified, runnable payload is *unhealthy* — only
  *unmanaged*. If you disagree, argue it with a case.
- The runtime probe still applies to it (`unavailable:` reporting unchanged).
- `marker: null` continues to distinguish it in JSON; the legacy carve-out is unchanged.

Scope note the option's wording does not cover: today only `parley-bidding` ships a
manifest. The other five units (core included) have nothing to verify against, so under
(b) as worded the README-first full install still reads **five `malformed`, exit 1** — my
measurement above. That residual is real but follow-up scope: either ship manifests for
every unit (the honest end state; the mechanism already generalizes) or the docs state
plainly that `doctor`'s managed verdicts for manifestless units require this tool's
installer. Do not silently leave it.

## New findings

### [NIT] `status` prints `integrity:` detail for add-on units but not for the core unit
**Where:** `lib/installer.js`, `writeResult` status branch (~1535–1547) — the new detail
block lives inside the `slice(1)` add-on loop.
**What:** A malformed **core** unit prints `claude: malformed <dest>` with no reason, while
every add-on gets its `integrity:` / `unavailable:` lines.
**Why it matters:** Small, but it lands exactly on the skills-CLI path above — the
foreign-installed core is the unit whose malformation a user most needs explained, and
`status` is the first command they will run. `doctor` already prints core problems.
**Evidence:** measured on the skills-CLI full install: `claude: malformed …/parley-deck`
bare, followed by six explained add-on lines.
**Fix:** Print the core unit's `problems`/`runtime` before the add-on loop (or unify the
loop). Rides along with the (b) touch of the same printers.

### [NIT] Probe memoization key is PATH-only, but the probe spawns with the full caller env
**Where:** `lib/installer.js:1297–1317` (`probePython3`).
**What:** The cache key is `PATH` alone, while `spawnSync` receives the entire effective
environment. `PYTHONHOME`, `PYTHONSTARTUP`, etc. can change the probe's answer without
changing the key, so two caller envs in one process that differ only there get the first
one's verdict.
**Why it matters:** Marginal — it needs a library caller passing two such envs within one
process; the CLI path is single-env, and the comment already states the PATH-keyed contract
honestly.
**Evidence:** read from the code; not measured (constructing a `PYTHONHOME`-broken python3
is out of proportion to the stakes).
**Fix:** Either key on the small set of vars that can affect interpreter startup, or leave
as-is — but then say "per PATH" means *only* PATH in the docs. No block either way.

## What I verified and found correct

- **The reviewed commit is as stated**: the real repository's HEAD is `5c324ef`, worktree
  clean. All execution happened in a `git clone` under `/tmp` plus `/tmp` `HOME` sandboxes;
  the reviewed tree was never mutated, nothing under `skills/parley-bidding/` was touched,
  and no Python ran outside the project's own runner (which uses `-B` /
  `PYTHONDONTWRITEBYTECODE=1`).
- **Claimed counts reproduce exactly**: 286→290 node tests (the four new ones), 54 Python,
  manifest 47 files `sha256:7854adf1…` — same digest as round 2.
- **The four new tests are honest**: the memoization test stubs both arms (no dependence on
  the host interpreter, the exact F5 lesson); the `paths` test proves non-execution with a
  sentinel rather than trusting the option wiring; the empty-PATH regression test fails
  hard if the host lacks python3 rather than silently passing.
- **Fail-safe directions unchanged**: empty PATH → `valid` payload + `runtime.ok:false` +
  exit 1; a 3.12 stub satisfies the floor; a below-floor interpreter fails the Python leg
  instead of skipping; a markerless unit never reaches the probe (`ok` gates it).
- **Option wiring is complete**: all three `targetStatus` call sites pass explicit probe
  options; there is no fourth caller.
- **Q1/Q2 answers verified by measurement**, detailed above.
- **No install/status/doctor regressions observed** across the sandboxes: fresh installs
  green, foreign installs behave exactly as analyzed above, `doctor --json` shape stable
  (`status`, `problems`, `marker`, `runtime` per unit).

## Open questions for the implementer

1. The five manifestless units on the README-first path (scope note in the ruling): do you
   want a follow-up idea that ships manifests for every unit, or a docs note for now? I
   recommend the follow-up — the integrity basis for the other five is currently just
   required-file presence, and the skills-CLI path makes that visible to every new user.
2. If you implement (b) with a different status string or exit-code choice than the ruling
   proposes, record the reasoning in `IMPLEMENTATION.md` — the round-4 re-review will check
   the four requirements above, not the names.
