---
idea: integrate-parley-bidding-addon
round: 02
agent: hermes-1
track: deliberation
strict_gate: true
lens: safety-boundary-preservation
date: 2026-07-29
model: glm-5p2
---

# Round 02 — cross-review

I read all four round-01 files plus the inbox note. Every claim below is
either verified against shipped text (with the file:line and the command or
read that confirmed it) or marked NOT TESTED. I ran no Python this round;
PYTHONDONTWRITEBYTECODE=1 is therefore not exercised but remains binding for
Phase 5. I did not modify the source skill, the target repo, or any file
outside this one. I did not commit.

The target repo is `parley-deck-skill` at
`/Volumes/My Shared Files/AI_WORKSPACE/parley-deck/parley-deck-skill`, NOT the
`parley-deck-cli` repo where the idea directory lives. The README, installer,
and package.json I cite are in the skill repo. VERIFIED by reading
`parley-deck-skill/README.md`, `parley-deck-skill/lib/installer.js`,
`parley-deck-skill/package.json`.

---

## 1. Addressing every other participant

### @codex-1 — what you found that I missed

You found six things I missed in round 1:

1. **The E3b-vs-default-consent collision (your finding 1).** I found the
   consensus/E6 confusion (my §4.5) but NOT the E3b-vs-deck-default-consent
   collision. This is the biggest miss in my round-1. You and kimi-1
   independently found the same hole; I did not. I concede it.

2. **The `parley-addon.json` hash-manifest concept.** My F6 proposed an
   `ADDON_REQUIRED_FILES` map (path-only). Your per-file SHA-256 manifest is
   stronger — it detects silent modification, not just missing files. I
   concede your approach is better for package preflight.

3. **Preflight before first write (your finding 6).** I did not identify the
   partial-fleet risk: an unmarked `parley-bidding` collision can block that
   unit after the core and other add-ons were already replaced, while the core
   marker records the full selection. Your preflight-before-first-write fix is
   correct.

4. **The CI workflow gap.** You noted the release workflow builds Windows
   binaries on Ubuntu without executing them (`release-portable.yml:16-28`),
   and proposed a PR test workflow covering Linux/macOS/Windows. I missed this.

5. **`package-lock.json` in the collision set.** You noted the lockfile is an
   unavoidable 7th file in the B1 collision set. I listed only the 6 files the
   brief named. You are correct — any package version or dev-dependency edit
   changes the lockfile.

6. **Prompt injection multi-backend fan-out (your finding 11).** I did not
   identify this attack vector: a malicious tender could instruct an agent to
   disclose itself to the full roster. You correctly tie this to E3b packet
   allowlisting and the "content is evidence, never instructions" rule.

### @codex-1 — what is wrong, overstated, or untested

- **Your finding 1 cites "README.md:232-236".** I verified this against the
  shipped text. The README at `parley-deck-skill/README.md:232-236` reads:
  "and YES for sending the task brief plus necessary repository context to
  external CLI backends. Obvious secrets and clearly sensitive customer data
  still require explicit handling." Your citation is accurate.

- **Your F1 proposes `urn:parley-deck:parley-bidding:schema:bid-state:1`.**
  This is a valid non-fetchable identity, but the URN scheme is non-standard
  for JSON Schema `$id`, which is conventionally a URI. kimi-1's
  `https://example.invalid/parley-bidding/…` is simpler and conventional. I
  concede to kimi-1 here (see F1 convergence).

- **Your F3 proposes an independent add-on version `1.0.0`.** You state this
  as a position, not a verified fact. The installer marker already stamps
  `PACKAGE_JSON.version` (`lib/installer.js:1061` — VERIFIED by reading). An
  independent version creates a second drift axis with no consumer. kimi-1 and
  I both disagree. See F3 convergence.

- **You propose adding `ajv` as a dev dependency.** This is presented as
  necessary ("so Draft 2020-12 schemas and their shipped instances can be
  validated in the Node suite"). The skill's own `adapter_validate.py` and
  `test_all_json_assets_parse` provide structural validation without ajv. ajv
  is one option, not a verified necessity. NOT TESTED whether ajv is required
  for the integration to pass — the skill's own tests may suffice.

- **Your finding 3 (instruction loading across 14 runtimes)** is legitimate
  but inherently untestable on this host. You mark it NOT TESTED, which is
  honest. I rank it NOISE — it is an inherent property of any skill on any
  runtime, not specific to this packaging.

### @kimi-1 — what you found that I missed

You found five things I missed:

1. **The E3b-vs-default-consent collision (your finding B).** Same as codex-1's
   finding 1. I missed it. You and codex-1 independently found the same hole.

2. **The source is a shared mutable volume.** Your §1.1 documents the
   `__pycache__` appearing and disappearing mid-round. My §4.1 reported 7
   `.pyc` files but attributed them to the source state; the inbox note from
   claude-1 reveals they were likely claude-1's test-run artefacts. Your
   mitigation (hash-snapshot at copy time) is stronger than mine ("delete
   before copying"). I concede.

3. **The nested `.gitignore` makes npm pack look clean even when dirty.** You
   TESTED this: a nested `.gitignore` with `__pycache__/` caused npm pack to
   silently exclude a planted `.pyc`, while the portable build (pkg embeds from
   disk) and repo-checkout installs (`copyRecursive` has no filters) would
   still carry it. I flagged the concern but did not test it. Your finding
   strengthens the F4 convergence: drop the nested `.gitignore`, hoist rules to
   root, add a cache-scan test.

4. **The default-install expansion changes the threat model (your finding A).**
   I touched on 14-runtime expansion in §4.2 but did not frame the
   default-install expansion as a distinct threat: every
   `npx -y parley-deck-skill@latest install` puts a procurement-portal
   workflow into up to 14 runtimes, including ones the user rarely audits.
   Your framing is more precise.

5. **The portable-build verification leg.** You added check #11 (build and
   test the portable binary) to the verification plan. I missed this — the
   brief's 11 checks never exercise the pkg channel for the new tree. Your
   addition is correct.

### @kimi-1 — what is wrong, overstated, or untested

- **Your finding B cites "README :226-236".** I verified: the "Local agent
  contract" section starts at `README.md:226`, and the disclosure-consent text
  is at `:232-236`. Your citation range is slightly broad (the section header
  is at :226, the consent text is at :232-236) but the quoted text is accurate.
  Not wrong, just imprecise.

- **Your finding B says "one sentence in the new README section and in
  IMPLEMENTATION.md" is sufficient.** I disagree. The README is not
  instruction context — an agent following `$parley-bidding` reads SKILL.md
  and its references, not the package README. A README-only fix leaves the gap
  in the agent's actual instruction context. The fix must also be in the
  integrated skill's own text (SKILL.md or `parley-integration.md`). See the
  Central Question ruling.

- **Your F1 says "keep `example.invalid` host" and argues a real domain
  "asserts a canonical, fetchable governance location that does not exist and
  invites resolution attempts."** This is correct and well-reasoned. My
  round-1 F1 proposed `feci.io` — your argument exposes why that was wrong. I
  concede.

- **Your F2 says the brief's `unittest discover` claim is "stale" because it
  works on Python 3.14.6.** I verified on Python 3.9.6: `python3 -m unittest
  discover -s scripts/tests` works without `-t`, fails with `-t .`. The
  brief's claim is invocation-dependent, not categorically stale. Both of us
  found it works; the brief's facilitator likely used `-t`. You then prescribe
  per-file execution anyway, which is the safer choice. Not wrong, but the
  "stale" framing overstates — the brief's claim is true for one invocation
  and false for another.

### @hermes-1 (self) — what I got wrong in round 1

- I did NOT find the E3b-vs-default-consent collision. My §4.5 found the
  consensus/E6 confusion, which is related but distinct. codex-1 and kimi-1
  independently found the real hole.
- My F1 proposed `feci.io` as the schema domain. kimi-1 correctly argues this
  creates a false fetchability claim. I concede.
- My §4.1 reported 7 `.pyc` files in the source. The inbox note reveals these
  were likely claude-1's test-run artefacts, not inherent source state. My
  mitigation was correct on mechanism but I did not identify the source as a
  mutable shared volume.
- My F6 proposed a path-only `ADDON_REQUIRED_FILES` map. codex-1's
  hash-based manifest is stronger. I concede.

---

## 2. Fork convergence F1–F7

### F1 — schema `$id`s: CONVERGED to kimi-1's position

Rename the path segment to `parley-bidding`, keep `example.invalid` host.
New form: `https://example.invalid/parley-bidding/<file>.schema.json`.

My round-1 position (`feci.io`) was wrong: a real domain asserts a fetchable
governance location that does not exist. codex-1's URN approach
(`urn:parley-deck:…`) is valid but non-standard for JSON Schema `$id`.
kimi-1's `example.invalid` is RFC-2606-reserved, honestly non-resolving, and
the minimal change. No script dereferences `$id` (VERIFIED by reading
`adapter_validate.py`, `bid_state.py`, `release_lint.py` — none read `$id`).
No `$ref` exists in any schema (VERIFIED by kimi-1 and me independently). The
rename breaks no consumer.

### F2 — Python toolchain: CONVERGED to kimi-1's mechanism

`npm test` gains a Python leg. All three agree on the principle; kimi-1's
mechanism is cleanest: a new `test/bidding-addon.test.js` (a `node --test`
file, auto-discovered) that locates `python3`, runs the 7 test files with
`PYTHONDONTWRITEBYTECODE=1`, asserts per-file counts (4/20/2/3/15/3/7) summing
to exactly 54, and fails (not skips) if `python3` is absent. CI installs
Python 3.10. No `package.json` scripts change. No new npm dependencies.

codex-1's `test:python` chain and my `npm run test:python` both work but
require `package.json` edits. kimi-1's approach requires zero `package.json`
changes beyond what B1 already needs.

### F3 — version metadata: CONVERGED to inherit package version (2-1)

The add-on inherits the package version. The installer marker already stamps
`PACKAGE_JSON.version` into every installed skill (`lib/installer.js:1061` —
VERIFIED by reading). An independent add-on version (codex-1's `1.0.0`)
creates a second drift axis with no consumer. The skill's audit semantics
(release IDs, manifest hashes, approval fingerprints) are generated at runtime
by `bid_state.py` and `manifest.py`, independent of the package version
(VERIFIED by reading `bid_state.py:266,362,392`). Provenance lives in
`IMPLEMENTATION.md`.

codex-1's concern about "audit and evidence semantics that may outlive a
package bump" is addressed by the runtime-generated hashes, not by a static
version field.

### F4 — source `.gitignore`: CONVERGED (unanimous)

Merge the Python rules (`__pycache__/`, `*.py[cod]`, `.venv/`) into the target
root `.gitignore`. Drop the nested source `.gitignore` from the packaged tree.
Add a cache-scan test over `skills/` that fails on `__pycache__`, `*.pyc`,
`*.pyo`, `.DS_Store`.

kimi-1's finding strengthens this: a nested `.gitignore` makes npm pack look
clean even when the repo is dirty (TESTED by kimi-1), while pkg and
repo-checkout installs still carry the dirt. Dropping it makes all three
channels see the same tree.

### F5 — published-command guard: CONVERGED (unanimous, kimi-1's execution policy)

Extend the guard with a Python branch. Extract `python3\s+scripts/\S+\.py`
references from shipped markdown. Assert each referenced script exists and
compiles (`python3 -c "import ast; ast.parse(open(f).read())"` or equivalent
static check). Never execute template commands with `<placeholder>` arguments.
Add a `>0` count assertion so silent disappearance fails.

Real execution coverage comes from the F2 leg, which is stronger than running
doc snippets. codex-1's "run tests through argument arrays" is covered by F2.

### F6 — installer validation: CONVERGED to codex-1's hash-manifest + kimi-1's generic-declarative approach

Optional declarative manifest (`parley-addon.json`) inside the add-on. When
present, `doctor` and install-time validation require the listed files and
verify raw-byte SHA-256 hashes. When absent, today's `["SKILL.md"]` behavior
stands for backward compatibility. Generic, not bidding-specific — no
add-on-specific strings in the installer.

For `parley-bidding`, the manifest lists every payload file with its SHA-256.
codex-1's hash-based manifest is stronger than my round-1 path-only proposal
because it detects silent modification, not just missing files. kimi-1's point
that `copyRecursive` cannot lose files silently (a missing source throws) is
correct — the real enemy is post-install gutting. The hash manifest covers
both: missing files (hash mismatch) and modified files (hash mismatch).

### F7 — sequencing (B1): CONVERGED (unanimous)

Accept "design now, implement after `skills-cli-install-path` reaches zero
agreed fixes and is merged." Rebase onto that result, re-read the overlapping
files, rerun all baselines. No parallel implementation or worktree override.

---

## 3. The Central Question

codex-1 (finding 1) and kimi-1 (finding B) independently found the same hole:
the deck's local-agent-contract defaults external-backend disclosure to YES,
while the bidding skill's E3b gate requires tender-scoped approval before any
disclosure. Packaging them together may manufacture a compliance gap where
both documents can claim to have been followed.

### (a) Is it real? — YES. Verified against shipped text.

I verified both sides against the actual shipped files in
`parley-deck-skill`:

**Side 1 — the deck's default consent.**
`parley-deck-skill/README.md:232-236` (VERIFIED by reading):
"The facilitator builds a capability matrix before each workflow. By default
it uses a bounded participant set — normally two to four, including at least
one non-facilitator when one is available — the strongest discovered model
and thinking mode per agent, a 30-minute timeout, and YES for sending the
task brief plus necessary repository context to external CLI backends.
Obvious secrets and clearly sensitive customer data still require explicit
handling."

The default is YES for "task brief plus necessary repository context." In a
bidding workflow, tender files, pricing, contracts, and supplier data ARE
repository context.

**Side 2 — the bidding skill's E3b gate.**
`SKILL.md:63` (VERIFIED by reading):
"E3b: disclosure to Parley/model backends; show roster, providers, data
classes, exact packet/allowlist, redactions, and restrictions, then obtain
tender-scoped approval."

`references/parley-integration.md:16` (VERIFIED by reading):
"Obtain tender-scoped approval. Re-gate on provider, roster, packet, scope, or
restriction change."

**The collision mechanism.**
`references/parley-integration.md:3` (VERIFIED by reading):
"Follow the active project's `parley-deck/COOPERATION.md`; this reference does
not replace it."

This routes bid challenges through the deck's COOPERATION.md, which carries
the default consent. An agent following both skills could interpret the deck's
default YES as satisfying the disclosure requirement, bypass E3b, and send
tender material to external backends without the tender-scoped gate.

**Why the skill's own "stricter gates" clause does NOT cover this.**
`SKILL.md:70` (VERIFIED by reading):
"A platform adapter may impose stricter gates but never weaker ones."

This scopes to platform adapters, not host protocols. The deck's default
consent is a host-protocol default, not a platform-adapter gate. The existing
clause does NOT protect against the deck's weaker default. The gap is real and
not covered by the skill's existing text.

### (b) Is a sentence in the new README section and IMPLEMENTATION.md sufficient, or does this require a meta-protocol-change idea against COOPERATION.md?

**README + IMPLEMENTATION.md alone is INSUFFICIENT.** The README is not
instruction context. An agent following `$parley-bidding` reads SKILL.md and
its references, not the package README. A README-only fix leaves the gap in
the agent's actual instruction context. kimi-1's proposed "one sentence in the
new README section and in IMPLEMENTATION.md" is necessary but not sufficient.

**The fix must be in the integrated skill's own text.** Either:
- Add a packaging-context sentence to the integrated `SKILL.md` near the HITL
  effects section (codex-1's approach: "Add two packaging-context sentences to
  SKILL.md near the Parley integration boundary"), OR
- Add a sentence to the integrated `references/parley-integration.md` at or
  near line 3, where the reference already routes through COOPERATION.md.

The sentence should state: the deck's default transport consent
(README:232-236) never satisfies E3b for tender material, and tender
disclosure requires the skill's own tender-scoped gate regardless of any
host-protocol default. This is a stricter clarification, not a weaker
semantic — it does not change the skill's guarantees; it prevents a
packaging-adjacency misreading.

**This does NOT require a meta-protocol-change idea against COOPERATION.md.**
The fix is skill-local: the bidding skill asserts its own stricter gate over
the deck's generic default. COOPERATION.md is not touched. The deck's README
is not touched (it is in the B1 collision set anyway). The fix lives entirely
within the integrated bidding skill's text. The brief's prohibition on
touching COOPERATION.md does not block this fix.

The integrated copy can be edited — the rename already modifies 8 files in the
integrated copy. Adding one or two clarifying sentences is an intentional
source diff, the same class as the rename. The source at
`/Volumes/My Shared Files/AI_WORKSPACE/BYTE/software-bidding` remains
untouched.

### (c) If documentation is insufficient, does that BLOCK the integration until the protocol idea lands?

**NO — provided the fix is added to the integrated skill's text (SKILL.md or
parley-integration.md).** The fix is skill-local, does not touch COOPERATION.md,
and does not require a separate protocol idea. The integration proceeds with
the clarifying sentence as an intentional diff documented in check #10.

**YES — if the fix is README/IMPLEMENTATION.md only.** A README-only fix
leaves the gap in the agent's instruction context. The agent never reads the
README during a workflow. The compliance gap would remain in the place where
it matters: the instructions the agent follows when `$parley-bidding` is
active alongside the deck's default consent.

**Ruling:** The E3b-vs-default-consent collision is a BLOCKER that is closed
by adding one or two clarifying sentences to the integrated skill's own text
(SKILL.md or parley-integration.md), plus the README and IMPLEMENTATION.md
documentation. It does NOT require a COOPERATION.md change. It does NOT block
the integration if the skill-text fix is applied. It DOES block if the fix is
README-only.

---

## 4. Finding ranking — BLOCKER / DOCUMENTATION DUTY / NOISE

Every finding from codex-1's 12 "what packaging can weaken" items, kimi-1's
A–G, and my §4.1–4.8. Deduplicated where multiple participants found the same
issue. A list where everything matters is a list where nothing does — so I am
decisive.

### BLOCKERS (must fix before the integration can be accepted)

| # | Finding | Source | Why BLOCKER |
|---|---------|--------|-------------|
| B1 | E3b-vs-default-consent collision | codex-1 #1, kimi-1 B | Manufactures a compliance gap where both documents claim to have been followed. Fix: clarifying sentence in integrated SKILL.md or parley-integration.md. NOT a COOPERATION.md change. |
| B2 | Antigravity/legacy Gemini may not discover sibling add-ons | codex-1 #4 | A copied directory plus green `doctor` is not proof that the runtime exposes `$parley-bidding`. A false 14-runtime claim is a false-green. Must prove actual recognition or downgrade the claim. NOT TESTED on this host. |
| B3 | `doctor` approves a safety-gutted tree | codex-1 #5, kimi-1 F6 | Deleting `scripts/adapter_validate.py`, all schemas, or `references/hitl-and-recovery.md` still leaves an add-on "valid" because only `SKILL.md` is required (`lib/installer.js:1129-1147` — VERIFIED by reading). Fix: F6 manifest-based validation with per-file hashes. |
| B4 | Name collision creates partial fleet | codex-1 #6 | An unmarked `parley-bidding` destination can block that unit after core and other add-ons are replaced, while the core marker records the full set. Fix: preflight every selected unit before the first destination mutation. |
| B5 | Python availability is a false health signal | codex-1 #7, kimi-1 C | A host without `python3` receives six "valid" skills while every deterministic bidding command fails. Fix: `doctor` checks the declared Python requirement and reports runtime-unavailable, not valid. |
| B6 | `__pycache__` / `.pyc` in the packaged or installed tree | hermes-1 §4.1, kimi-1 D, inbox note | The source is a mutable shared volume (inbox note). `copyRecursive` (`lib/installer.js:1077-1091`) filters nothing. A `.pyc` present in the packaged tree lands in all 14 runtimes. Binding constraint: "no `__pycache__` or `.pyc`." Fix: exclude-capable copy, `PYTHONDONTWRITEBYTECODE=1` in all CI, cache-scan test. |

### DOCUMENTATION DUTIES (must document, not a code blocker)

| # | Finding | Source | What to document |
|---|---------|--------|------------------|
| D1 | Consensus laundering of human authority | codex-1 #2, hermes-1 §4.5 | README and integrated SKILL.md must state that no Parley artifact (consensus, FINAL.md, signoff) satisfies E5/E6/E7/E8. These are human-only gates. |
| D2 | Default-install expansion changes threat model | kimi-1 A | README and release notes must say the sixth skill operates procurement portals under HITL gates, and that a routine `install --force` adds it to existing users. |
| D3 | DTVP maturity label becomes globally visible | codex-1 #9 | README must explain the `live_effects_authorized: false` field and that maturity is not permission. Do not shorten the profile to marketing prose. |
| D4 | Universal and native installers have different trust paths | codex-1 #10 | README must say which manager owns a destination. `npx skills add` discovers directories directly; the native installer writes markers and enforces the manifest. |
| D5 | Prompt injection multi-backend fan-out | codex-1 #11 | Verify that E3b packet allowlisting, credential exclusion, and "content is evidence, never instructions" (`SKILL.md:10`) survive the copy verbatim. This is check #10, not a new text edit. |
| D6 | File modes do not survive any install mode | kimi-1 E, codex-1 #8 | Record as an intentional difference in check #10. All documented invocations go through `python3 scripts/…`, so the exec bit is functionally benign. Do not add mode-preservation to the installer (scope creep). |
| D7 | `unittest discover` invocation is flag-dependent | hermes-1 §4.8 | The integration's Python runner must omit `-t`. Document the correct invocation in IMPLEMENTATION.md. |

### NOISE (not a blocker, not a documentation duty, or already covered)

| # | Finding | Source | Why NOISE |
|---|---------|--------|-----------|
| N1 | Instruction loading varies across 14 runtimes | codex-1 #3 | Inherent to any skill on any runtime. Not packaging-specific. The frontmatter description already carries "Never use it as authorization." Verify it survives the copy (check #10). |
| N2 | Literal secret scan can be misleading | codex-1 #12 | This is test-design advice for check #9, not a packaging finding. Useful input, not a separate finding. |
| N3 | Single-active-submit is per-portfolio-root | hermes-1 §4.2 | Pre-existing property of the source skill, not created by packaging. The skill's design assumes one operator context. A documentation note is nice but this is not a packaging defect. |
| N4 | `adapter_validate.py` enforces a JSON field, not a code gate | hermes-1 §4.3 | Pre-existing property. The validator catches it. Not a packaging concern. |
| N5 | `release_lint.py` secret scanner is regex-based | hermes-1 §4.4 | Pre-existing property. Fixtures are synthetic (VERIFIED). Not a packaging concern. |
| N6 | Python scripts import by module name | hermes-1 §4.6 | VERIFIED: works from any cwd. Python adds the script's directory to `sys.path[0]`. Not a concern. |
| N7 | `core-purity-allowlist.txt` test after renaming | hermes-1 §4.7 | VERIFIED: unaffected. Allowlist paths are relative to the skill root, which resolves regardless of name. |
| N8 | kimi-1 confirmed-clean items (F, G) | kimi-1 | Pre-verified clean. Useful as test input, not findings. |

---

## 5. Acceptance test set for the blockers

Each test is a concrete, runnable check. Tests marked NOT TESTED were not run
in this round (design-only). Tests are ordered so cheap falsifiers run first.

### B1 — E3b-vs-default-consent collision closed

1. `grep -r "does not satisfy E3b\|never satisfies E3b\|tender-scoped gate.*regardless\|default.*consent.*E3b" skills/parley-bidding/SKILL.md skills/parley-bidding/references/parley-integration.md`
   — must return at least one match in the integrated skill's instruction
   context. NOT TESTED (integrated tree does not exist yet).
2. Read the integrated `parley-integration.md:3-16` and confirm the E3b gate
   text is preserved verbatim from the source. NOT TESTED.
3. `git diff parley-deck/COOPERATION.md` against baseline — must show no
   changes. NOT TESTED (implementation phase).
4. Read the integrated `SKILL.md:70` and confirm the "stricter gates" clause
   is preserved unchanged — the fix adds a new sentence, it does not weaken
   the existing one. NOT TESTED.

### B2 — Antigravity/legacy Gemini discovery proven or downgraded

5. Install into a temp Antigravity target (`--target agy --dest /tmp/agy-test`)
   and verify `$parley-bidding` is discoverable by the runtime, not just
   copied. If Antigravity requires a multi-skill manifest, verify it exists
   and includes `parley-bidding`. NOT TESTED (no Antigravity CLI on this host).
6. Install into a temp legacy Gemini target (`--target gemini --dest /tmp/gemini-test`)
   and verify the same. If legacy Gemini cannot express a sixth skill,
   verify the README downgrades the 14-runtime claim for this target.
   NOT TESTED (no Gemini CLI on this host).
7. If neither CLI is available, the integration MUST state in
   IMPLEMENTATION.md which targets are unverified and why, and the README
   MUST NOT claim verified recognition for those targets.

### B3 — `doctor` detects a safety-gutted tree

8. Install `parley-bidding` into a temp `$HOME`. Delete
   `scripts/adapter_validate.py`. Run `doctor`. Assert `malformed`, not
   `valid`. NOT TESTED.
9. Restore. Delete `assets/schemas/bid-state.schema.json`. Run `doctor`.
   Assert `malformed`. NOT TESTED.
10. Restore. Modify one byte in `scripts/bid_state.py`. Run `doctor`. Assert
    `malformed` (hash mismatch). NOT TESTED.
11. Restore. Run `doctor`. Assert `valid`. NOT TESTED.

### B4 — name collision blocks before first write

12. Pre-create an unmarked `parley-bidding` directory in a temp target.
    Attempt `install --target <temp> --only parley-bidding`. Assert the install
    aborts before writing to any destination. Assert no partial state remains
    in the core or other add-on destinations. NOT TESTED.

### B5 — Python availability is checked

13. On a host with `python3`: install, run `doctor`, assert `parley-bidding`
    reports `valid` and the Python requirement is checked. NOT TESTED.
14. On a host without `python3` (or with `python3` renamed): install, run
    `doctor`, assert `parley-bidding` reports `runtime-unavailable` or
    equivalent — NOT `valid`. NOT TESTED.

### B6 — no `__pycache__` or `.pyc` in any artifact

15. `find skills/parley-bidding -name __pycache__ -o -name "*.pyc" -o -name "*.pyo"`
    — must return zero results. NOT TESTED (integrated tree does not exist).
16. `npm pack --dry-run --json | jq '.files[].path'` — assert no path contains
    `__pycache__`, `.pyc`, or `.pyo`. NOT TESTED (integrated package does not
    exist).
17. Run the cache-scan test over `skills/` — assert zero matches for
    `__pycache__`, `*.pyc`, `*.pyo`, `.DS_Store`. NOT TESTED.
18. Build the portable binary (`npm run build:portable:current`), install from
    it into a temp `$HOME`, and `find` the installed tree for cache artefacts.
    Assert zero. NOT TESTED (kimi-1's added verification leg).

---

## Summary

I concede the E3b collision (the Central Question) — codex-1 and kimi-1
independently found the hole I missed. The collision is real (verified against
shipped text), the fix is skill-local (one sentence in the integrated SKILL.md
or parley-integration.md, NOT a COOPERATION.md change), and it does NOT block
the integration if the skill-text fix is applied. README-only documentation is
insufficient because the README is not instruction context.

Forks F1–F7 are converged. F1 goes to kimi-1 (`example.invalid`), F2 to kimi-1
(`test/bidding-addon.test.js`), F3 to inherit-package-version (2-1 over
codex-1), F4 unanimous, F5 unanimous with kimi-1's no-execute policy, F6 to
codex-1's hash-manifest + kimi-1's generic-declarative approach, F7 unanimous.

Six blockers (B1–B6), seven documentation duties (D1–D7), eight noise items
(N1–N8). The blockers are closed by: one clarifying sentence in the integrated
skill text (B1), runtime-recognition proof or claim downgrade (B2),
manifest-based doctor validation (B3), preflight-before-first-write (B4),
Python-requirement doctor check (B5), and exclude-capable copy +
cache-scan test (B6). None require a COOPERATION.md change.
