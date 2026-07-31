---
idea: integrate-parley-bidding-addon
review-round: 4
agent: kimi-1
date: 2026-07-30
reviewed-commit: b180127
---

## Verdict
ACCEPT

Every round-3 condition is closed by my own re-measurement at `b180127`, the claimed counts
reproduce on both interpreters, and the deferred residual is exactly the docs option my
round-3 ruling offered — I agree with deferring it. New findings are three NITs, none
blocking. All execution happened in a `git clone` under `/tmp` (`/tmp/pds-r4` @ `b180127`)
plus `/tmp` `HOME` sandboxes; the reviewed tree was never mutated. Note: the repository's
HEAD advanced to cycle 7 (`49fc3ec`) during this review — claude-1's own commits, not mine;
everything below is measured at the reviewed commit.

## Round-3 findings — closed or not

### The (b) ruling — CLOSED, all four requirements measured
Sandbox: fresh `install --target codex --yes` per case, then `doctor --target codex --json`
with the host's working python3 3.14.6 in PATH unless stated. Measured:

| case (ruling row) | measured verdict | exit |
|---|---|---|
| marker deleted, manifest fully verifies | `valid-unmanaged`, `managed:false`, `marker:null`, `problems:[]` | **0** |
| same, but `PATH:""` | `valid-unmanaged`, `runtime.ok:false`, "python3 is not available" | 1 |
| marker present, unreadable (`{ not json`) | `malformed`, "unreadable or is not valid JSON" | 1 |
| marker naming another installer | `malformed`, "belongs to \"other-installer\"" | 1 |
| marker + manifest both deleted | `malformed`, "no parley-deck-skill install marker" (gutting signal) | 1 |
| marker + one payload script deleted | `malformed` | 1 |
| marker deleted + stray `scripts/__pycache__/*.pyc` | `malformed` (undeclared file fails the verify) | 1 |
| core unit, marker deleted | `malformed` — the core ships no manifest, so it can never qualify | 1 |
| legacy `markerSchema`-less marker (2.0.0 shape) | `valid`, `managed:true`, exit 0 — **unchanged even after I appended a byte to a payload script** (carve-out untouched, soft edge as disclosed) | 0 |

The four ruling requirements, point by point:

1. **Distinct status string**: `valid-unmanaged` in text and JSON, with `managed` as the
   boolean codex-1 asked for. Different name than my example, same properties — as agreed,
   the requirements were the contract, not the names.
2. **Only entirely-absent qualifies**: rows 3–5 above. A present marker of any kind never
   reaches the new path (`unmanagedButVerified` runs only under `!state.present`).
3. **Does not fail health**: row 1 — doctor **exit 0** with a working interpreter when
   `valid-unmanaged` is the only finding. Their own new test asserts only `!= malformed`
   under an empty PATH; my measurement is the stronger one the ruling asked for.
4. **Runtime probe still applies**: row 2 — exit 1, `runtime.ok:false`, reporting unchanged.

Plus the two guards I demanded: `marker: null` still distinguishes the verdict in JSON, and
the legacy carve-out is byte-for-byte behavior-identical (row 9).

### NIT — `status` explains the core unit — CLOSED
Measured on a full foreign install: `status` now prints the core line followed by
`  integrity: no parley-deck-skill install marker: …`, then per-addon verdicts with their own
`integrity:` lines, exit 0. The unified `explain` printer covers core and add-ons.

### NIT — probe cache key narrower than the spawn — CLOSED (with a new NIT, below)
Measured in one process, both orders: env A = real PATH + `PYTHONHOME=/nonexistent`
(verified: python3 cannot start, exit 1), env B = real PATH. Both orders give the
per-environment answer (A: unavailable; B: `python3 3.14`). The realistic gap is closed.

### codex-1's MINOR — broken shim satisfied the floor — CLOSED (verified, not mine but cheap)
Stub `python3` printing `4.not-a-version` → `runtime.ok:false`, "not available", doctor
exit 1. Stubs `3.12` and `4.0` → ok; `3.9` → below-floor failure with the right detail. The
anchored `^(\d+)\.(\d+)$` parse fails closed.

### hermes-1's MINOR — `status` always exits 0 — CLOSED as documented
`status` exit 0 measured on a tree full of malformed units; the CHANGELOG now names `doctor`
as the health gate. Contract kept, documented — what was asked.

### Claimed counts — reproduce exactly
- `node --test` in the clone after `npm ci`: **296 pass, 0 fail** with Homebrew python3
  3.14.6 first in PATH, and **296 pass, 0 fail** with `/usr/bin/python3` 3.9.6 first.
- Python leg: green on 3.14; on 3.9.6 it prints "python3 is 3.9, but the add-on declares
  >=3.10" and exits 1 — the intended below-floor contract, same precision as round 3.
- Manifest check: `parley-bidding: ok (47 files, sha256:7854adf1…)` — digest unchanged
  across rounds 2, 3, and 4. (296 = 290 + 6: one old test replaced, seven added.)

### The deferred residual — measured; I agree with the deferral
Faithful README-first replication (all six skills copied foreign, none tool-installed):
`parley-deck: malformed`, `parley-bidding: valid-unmanaged`, the other four add-ons
`malformed`, **doctor exit 1** — precisely the "one `valid-unmanaged` and five `malformed`"
the CHANGELOG states. Deferring the remaining five manifests is the docs option my ruling
offered ("or the docs state plainly"); the limit is stated plainly in `CHANGELOG.md` and
referenced from `README.md:78-81`, and B3.11 holding the other add-ons unaffected is a
ratified boundary I will not ask a fix-up cycle to widen. Record the manifest-for-every-unit
follow-up in the consensus, as promised.

## New findings

### [NIT] The widened probe-cache key is no longer injective — demonstrated collision, both directions
**Where:** `lib/installer.js:1341–1343` — `["PATH","PYTHONHOME","PYTHONPATH","PYTHONEXECUTABLE","VIRTUAL_ENV"].map(…).join("\u0001")`.
**What:** `\x01` is legal inside environment values on POSIX, so the join does not determine
the tuple. Two envs, both with all five vars defined — A: `PATH="/tmp/pds-r4-coll\x01PY"`,
`PYTHONHOME=""`; B: `PATH="/tmp/pds-r4-coll"`, `PYTHONHOME="PY\x01"` — produce the identical
key `"/tmp/pds-r4-coll\x01PY\x01\x01\x01\x01"` (verified by direct computation) but different
real answers: A's single PATH entry does not exist (no python3), B has a working stub.
**Why it matters:** measured in one process against a real install: order B→A makes A report
`python3 3.12`, `runtime.ok:true` — **fail-open**, a health-gate green for an environment
with no interpreter; order A→B makes B report unavailable — fail-closed. The old PATH-only
key was injective (env values cannot contain `\u0000`, so `"\u0000inherit"` was unforgeable);
the cycle-3 widening traded "incomplete" for "non-injective". NIT rather than MINOR because
the trigger needs a library caller passing two crafted envs with control characters in one
process — the CLI is single-env per process, and my first two collision constructions failed
on invariants (undefined-set NUL count, embedded-separator count) that real envs never meet.
**Evidence:** above; harness `probe.js` modes `coll-BA` / `coll-AB`, stub dir with `python3`
printing `3.12`.
**Fix:** `JSON.stringify` the five-element array as the key, or join with `\u0000` (values
cannot contain NUL, so the concatenation determines the tuple). One line.

### [NIT] `managed` is absent on `missing` units
**Where:** `skillUnitStatus`'s early return for a nonexistent dest (`lib/installer.js:1222`).
**What:** that shape carries `status/marker/missing/problems/runtime` but no `managed`;
every other status (`valid`, `valid-unmanaged`, `malformed`) now carries the boolean.
**Why it matters:** automation iterating `doctor --json` units and reading `.managed` gets
`undefined` for exactly the units most likely to be acted on. Measured: deleted add-on dir
→ `missing`, `hasOwnProperty("managed") === false`.
**Fix:** add `managed: false` to the early return.

### [NIT] An unmarked malformed tree loses the manifest detail the managed path reports
**Where:** the `!state.present` branch of `skillUnitStatus`.
**What:** when `unmanagedButVerified` fails, only the "no parley-deck-skill install marker"
problem is pushed; the underlying `verifyPayload` problems are discarded. A managed tree
with a stray file reports `unexpected: scripts/__pycache__/x.cpython-314.pyc` (measured);
an unmarked tree missing `scripts/adapter_validate.py` reports only the marker line
(measured) — the manifest module's own principle is telling the user *which* files are
wrong, and the README-first path is exactly where a partial copy needs that list.
**Why it matters:** diagnosability, not verdicts — every such case is still `malformed`,
exit 1.
**Fix:** have `unmanagedButVerified` return the problems instead of a boolean and append
them when it fails (it already computed them).

## What I verified and found correct

- **Reviewed commit as stated**: clone `HEAD` is `b180127`; the original worktree is clean
  and was never written by me (it advanced to `49fc3ec` on its own during the review).
- **No consumer of `skill.status` missed**: the only `=== "valid"`/`!== "valid"` sites are
  `doctorCommand.ok` and `writeResult`'s `broken` (both updated to include
  `valid-unmanaged`), plus `compatibilitySummary`/`recommendedActions`, which consume the
  *target-level* status — that mirrors the core unit, and the core cannot be
  `valid-unmanaged` (no `addon.root`, measured row 9). Complete audit, no fourth path.
- **Interactions**: `install` without `--force` over a foreign `valid-unmanaged` tree is
  `blocked` (exit 1) with the tree **byte-untouched** (fingerprinted before/after) and an
  honest "Re-run with --force" message; `install --force` replaces it, writes the marker,
  and the unit becomes `valid`, `managed:true`, exit 0 — the upgrade path from unmanaged to
  managed exists and works. `uninstall` without `--force` refuses to remove an unmarked
  tree (blocked, exit 1, intact); `--force` removes it. The tool never deletes or overwrites
  what it did not install unless told to.
- **`--no-addons` / `--only`**: after a `--no-addons` install, a later foreign copy of
  `parley-bidding` is invisible to plain `doctor` (expected set comes from the core marker's
  recorded selection — pre-existing semantics, exit 0) and visible as `valid-unmanaged`
  under `doctor --only parley-bidding` (exit 0 with a tool-installed core).
- **No laundering short of the ratified boundary**: unmarked + stray undeclared file →
  `malformed`; unmarked + deleted payload file → `malformed`; unmarked + deleted manifest →
  `malformed`; the only admissible case is full self-consistent verification, and the
  managed path decays identically (a stray `__pycache__` turns a *managed* tree `malformed`
  with the same `unexpected:` problem) — no new asymmetry. The self-consistent-swap +
  marker-deletion case reaching `valid-unmanaged` is the boundary round 3 ratified: the same
  writer can forge a schema-2 marker for a full `valid`, so (b) hands them nothing they did
  not already have.
- **Docs**: `CHANGELOG.md` rewritten from (a) to (b) including the narrow marker rules and
  the stated residual with exit-1 consequence; `README.md` references it from the install
  section; `status`'s always-0 contract documented with `doctor` named as the gate.

## Open questions for the implementer

1. HEAD moved to `49fc3ec` (cycles 4–7, titles mention `valid-unselected` and selection
   semantics) while I measured `b180127`. Are those commits the agreed manifest follow-up,
   or new scope? Round 5 should state which commit it reviews and re-run this matrix against
   it — my harnesses (`/tmp/pds-r4-harness/`) are disposable but trivially re-creatable from
   this file.
2. If you take the probe-key fix (NIT 1), note the choice in `IMPLEMENTATION.md`; the
   collision constructions above are the re-test.
