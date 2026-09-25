# Skill gate evidence logs — kimi-1, 2026-09-25

Evidence retention for the release-preparation audit's **OBS-2** (claude-1,
`release-preparation-audit-claude-1.md`, verdict PASS — preparation): the handoff narrated the
skill gate results but the delivery `logs/` held no `npm test` / `manifest --check` /
`npm pack --dry-run` output, unlike the 2.13.0 delivery. The gates themselves were satisfied
(the auditor re-ran them green); this is durability of channel evidence, not a gate failure.

## What was recoverable, what was not

The original preflight runs (session of 2026-09-25, before candidate commit `a5664d8`) wrote to
the terminal; my captured session log retains **tails only**, and I do not label excerpts
complete or reconstruct output. Surviving originals: the doctor JSON (already persisted at
`logs/skill-doctor-2.14.0.json`) and the install dry-run log. The three outputs above were gone,
so per the audit's durability observation I made **one bounded capture run** — evidence
retention only, no feature/code change, no broad Go suite.

## Capture run provenance and results

Clean detached checkout of the frozen skill candidate
`a5664d803f1fb6156ef95ff5921dcf13a3dd2031` (porcelain 0, removed after), `npm ci` from the
committed lock (exit 0, untracked environment; same install-scripts guard warning as the
original run), node v26.9.0 / npm 11.19.1 / python 3.14.7 / macOS arm64. Full stdout+stderr and
real exit codes, in `release-delivery/2026-09-25-designated-implementer/logs/`:

| command | log | exit | result |
|---|---|---|---|
| `npm test` | `npm-test-2.14.0.log` (426 lines) | 0 | tests 399 / pass 399 / fail 0 (Node + python leg + manifest `--check`) |
| `npm run manifest:check` | `manifest-check-2.14.0.log` (11 lines) | 0 | all 6 add-ons ok, parley-deck sha256 `57d3c535…f29` |
| `npm pack --dry-run` | `npm-pack-dry-run-2.14.0.log` (235 lines) | 0 | `parley-deck-skill-2.14.0.tgz`, 210 files, prepack gate |

Provenance companion: `logs/skill-gates-2.14.0-provenance.txt`. These match the results
narrated in `release-preparation-kimi-1.md` (399/399; 210 files) — no drift between narration
and capture.

## Unchanged artifact status

Candidates stay frozen: CLI `01fc49526ca66b771b549a01cfbbf4addeaa886b`, skill `a5664d8`. The
capture used `--dry-run` only — no tarball was produced; the delivery tarball
(`skill/parley-deck-skill-2.14.0.tgz`, sha256 `681b9076a017e67935403ad6cd9779cafa09a117d168331daca300404336135b`)
and all CLI assets are byte-unchanged. The organizer's external CLI release notes update
(OBS-1/3: candidate CI run, explicit Windows non-execution caveat) touched no source or
canonical peer file, and this report edits none — no peer review, signature, or record was
modified. Nothing published, versioned, installed, or changed in global/roster settings.
