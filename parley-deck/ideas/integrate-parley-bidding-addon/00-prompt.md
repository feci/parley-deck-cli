---
idea: integrate-parley-bidding-addon
author: user
created: 2026-07-29
track: deliberation
strict_gate: true
participants: [claude-1, codex-1, hermes-1, kimi-1]
target-repo: parley-deck-skill
source: /Volumes/My Shared Files/AI_WORKSPACE/BYTE/software-bidding  (READ-ONLY)
status: round-01
blocked-by: skills-cli-install-path   # same file set; see B1
---

## Problem / idea

Import the completed standalone `software-bidding` skill into `parley-deck-skill` as a
**sixth packaged skill**, installed and managed exactly like the existing add-ons.

**The user named it `parley-bidding`.** The originating agent's brief said
`addons/software-bidding/`; both parts of that path are superseded:

| the brief said | it must be | why |
|---|---|---|
| `addons/…` | `skills/…` | `addons/` no longer exists — the `skills-cli-install-path` idea moved every skill under `skills/` today so a generic installer discovers all of them |
| `software-bidding` | **`parley-bidding`** | direct user instruction, 2026-07-29 |

So the target is **`skills/parley-bidding/`**, and the package documents **six skills: the
core protocol plus five add-ons.**

## B1 — the blocking constraint, read this before anything else

`skills-cli-install-path` is **in review round 12 and not merged**. It owns exactly the files
this idea must also touch: `skills/`, `lib/installer.js`, `package.json`,
`test/installer.test.js`, `test/design-addons.test.js`, `README.md`.

**Intersecting file sets are the collision `parley-worktrees` exists to refuse.** This idea
therefore **designs now and implements after** that idea reaches zero agreed fixes. Round 1
and 2 are unaffected — they are design work. Phase 5 must not start until the block clears.
State in your round-1 file whether you accept that sequencing or propose a different one.

## Source of truth

- **Skill:** `/Volumes/My Shared Files/AI_WORKSPACE/BYTE/software-bidding` — **read-only.
  Nothing in it may be modified.**
- **Design/audit evidence:** that workspace's
  `parley-deck/ideas/software-bidding-multiplatform-skill/{FINAL.md,IMPLEMENTATION.md}` and
  `review/round-07/consensus.md`.
- The older `dtvp-bidding-hitl-skill` idea is **historical input only**, never the source.

## Measured ground truth (verified by the facilitator, 2026-07-29)

**48 files, 246 KB, no `__pycache__` or `.pyc`.** Layout: `SKILL.md`, `agents/openai.yaml`,
`assets/{core-purity-allowlist.txt, discovery-sources/×3, jurisdiction-profiles/×1,
platform-adapters/×4, schemas/×4, templates/×4}`, `references/×11`, `scripts/×7` plus
`scripts/tests/×7` and 3 fixtures, and a `.gitignore`.

**The Python baseline is exactly 54 tests, and all pass**, run file-by-file with
`python3 scripts/tests/test_X.py`:

| file | tests |
|---|---:|
| test_adapter_validate.py | 4 |
| test_bid_state.py | 20 |
| test_end_to_end.py | 2 |
| test_init_workspace.py | 3 |
| test_linters.py | 15 |
| test_manifest.py | 3 |
| test_skill_structure.py | 7 |
| **total** | **54** |

`pytest` is **not installed** on this host and `unittest discover` fails against that
directory (`Start directory is not importable`). Whatever runner the integration adopts must
be one that actually works here, and the 54 must be reproduced, not asserted.

**The rename touches 8 files** — and one of them is a test that will fail:

```
SKILL.md:2                          name: software-bidding
agents/openai.yaml:4                "Use $software-bidding to qualify, …"
scripts/common.py:2                 docstring
scripts/tests/test_skill_structure.py:16   assertRegex(skill, r"(?m)^name: software-bidding$")
scripts/tests/test_skill_structure.py:27   assertIn("$software-bidding", metadata)
assets/schemas/*.schema.json:3      "$id": "https://example.invalid/software-bidding/…"  (×4)
```

**Our repo hardcodes the current count in at least these places:**

```
README.md:14, :19, :109                       "five skills"
test/installer.test.js:408, :541, :563        the exact four/five-element add-on lists
lib/installer.js:341, :805                    the --only help and error text (generated)
```

## Open forks — take a position

- **F1 — the schema `$id`s.** They read `https://example.invalid/software-bidding/…`. A `$id`
  is an identity, not a label. Rename them (consistency, but a schema identity change), keep
  them (stable identity, but a name that no longer exists), or re-root them under a
  parley-deck identity? Argue it.
- **F2 — the Python toolchain.** This package is Node-only and dependency-free. This add-on
  ships 7 modules and 54 Python tests. Does `npm test` gain a Python leg? Does CI? Or do the
  Python tests stay a skill-local concern documented but unrun by the package suite? A test
  nobody runs is the false-green class this repo just spent 12 review rounds eliminating.
- **F3 — version and compatibility metadata.** Does the add-on inherit the package version, or
  does it need its own? The source skill has audit and evidence semantics that may outlive a
  package bump.
- **F4 — `.gitignore` in the source.** Copy, merge, or drop?
- **F5 — the published-command guard** added by `skills-cli-install-path` scans only
  `node --test` commands. This skill publishes **Python** commands. Extend the guard, or
  accept the gap and say so in writing?
- **F6 — installer validation.** `ADDON_REQUIRED_FILE` is `SKILL.md` only. Should an add-on
  shipping executable scripts and schemas assert more than that at install time?
- **F7 — sequencing (B1).** Accept "design now, implement after", or propose otherwise?

## Binding constraints

- **Preserve the skill's guarantees exactly:** HITL control, the evidence model, exact-byte
  hashing, the ambiguous-submission no-retry rule, authority and pricing gates, adapter
  maturity ceilings, and the separation between upload and submission.
- **Deterministic scripts stay portal-safe.** They must never browse, log in, upload, message,
  submit, amend or withdraw. Introducing any portal-mutation capability is a blocking defect.
- **Browser work stays adapter-driven.** No Chrome dependency. Where browser work is eventually
  performed, prefer the locally installed `ego-browser`/ego-lite skill; any other browser is an
  explicit fallback only.
- **Platform neutrality is preserved**, including the DTVP, Cosinex NRW, subreport ELViS and
  manual profiles.
- **No BYTE or customer material.** No tender files, customer data, credentials, receipts,
  prices or portal state. The integrated tree must be scanned for them.
- **Do not merge this into the core `parley-deck` SKILL.md.** It keeps its own trigger and its
  own safety boundaries.
- **Do not modify `COOPERATION.md`** to accommodate this add-on. A genuine protocol change is a
  separate meta-protocol-change idea.
- Existing add-ons and core-only installs stay backward compatible.

## Required verification (none may be reported as passing unless run)

1. `npm test` — the full Node suite.
2. `npm pack --dry-run` — every file of the add-on present; no `__pycache__`, no `.pyc`.
3. Default install · `--no-addons` · `--only parley-bidding`.
4. `doctor`, `status`, `paths`, `uninstall` for the new add-on.
5. The skill's own adapter validator.
6. All Python tests — **reproduce 54**, or explain the intentional difference.
7. Python compile checks that leave no cache artefacts behind.
8. Every shipped JSON schema and profile validates.
9. A scan of the integrated tree for BYTE/customer content, credentials, unresolved
   placeholders and generated caches.
10. A source-vs-integrated tree diff with **every** intentional difference documented.
11. `npx skills add <repo> --list` still finds the right set — now six.

## Acceptance criteria

`parley-bidding` is a complete, independently triggerable installed skill · existing install
modes unchanged · npm and portable artifacts carry every file, reference, asset, schema,
template and script · no BYTE-specific or confidential material · no portal-mutation capability
in deterministic tooling · documentation says **six** skills · every verification command run
and passing · the final independent review reports **no finding of any severity** ·
`IMPLEMENTATION.md` records provenance, decisions, checks, differences and release
implications.

## Hard stop

**Do not publish, release, push, merge or install globally.** Phase 5 ends by presenting the
completed diff and the validation evidence to the user for the release decision. This
overrides the standing release-after-every-change rule for this idea.
