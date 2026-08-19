---
agent: claude-1
idea: protocol-and-skill-audit
round: 1
date: 2026-08-20
---

## Findings

### F1 — `SKILL.md` never names the commands the protocol requires
severity: MAJOR

command:
```
$ for w in preflight retro "loop tick" consult repo-map tui; do
    echo "$w: $(grep -ci "$w" parley-deck-skill/skills/parley-deck/SKILL.md)"; done
```
output:
```
preflight: 0
retro: 0
loop tick: 0
consult: 0
repo-map: 0
tui: 0
```
contradicts: `COOPERATION.md:222` (§4.0 per-track table: the §9.0 readiness ping is **`full`** on
`standard` and `deliberation`), and `COOPERATION.md` naming `parley preflight` (×1), `parley retro`
(×2) and `parley loop tick` (×2).

why it matters: `SKILL.md` is the vendor-neutral instruction set a foreign agent loads to run this
protocol — the skill's own text calls it "the vendor-neutral instructions for all agents". An agent
that reads only the skill is told the readiness check is mandatory on two of three tracks and is
never told the command that performs it. Same for §13's retrospective and §14's human brake. This
is not a missing-nicety: it is the skill describing obligations it gives no way to discharge.

[PRIMARY] All counts re-run against `parley-deck-skill/skills/parley-deck/SKILL.md` at HEAD.

### F2 — The protocol instructs the bootstrap to run a command that fails the protocol's own guard
severity: MAJOR

command:
```
$ grep -n "parley roster render" parley-deck/COOPERATION.md
$ parley roster render --dir <iso> --yes --adopt-inherited && go test ./internal/protocol/...
```
output:
```
COOPERATION.md:57  ... (then regenerate the §2 view with `parley roster render`)

Regenerated §2 in <iso>/parley-deck/COOPERATION.md
--- FAIL: TestEmbeddedDefaultMatchesLiveDeck
    drift_test.go:60: live deck: anchor "| Agent ID       | Workspace dir  … | Role          |"
        appears 0 times, want exactly 1 (drift guard fails closed)
```
contradicts: `COOPERATION.md:57` versus `internal/protocol/drift_test.go:28,60`.

why it matters: §57 is the **deck-bootstrap** instruction — the first thing a new deck does. It
tells the operator to regenerate §2 with `roster render`, and `roster render` emits a four-column
compact header while the drift guard anchors the three-column padded one and fails closed. Anyone
who follows the protocol's own setup text on a deck inside this repository breaks the build. This
is D-B, but located one level up: not a tool that disagrees with a guard, a **protocol that
mandates the disagreement**.

[PRIMARY] Reproduced in isolated copies of both the deck and the Go tree; the shared tree was not
used.

### F3 — `masked-by-env` is in the documented closed STATUS vocabulary and never reaches STATUS
severity: MINOR

command:
```
$ grep -rn "masked-by-env" internal/ | grep -v _test
```
output:
```
internal/app/roster_set.go:83:  // `masked-by-env` was in the frozen STATUS vocabulary with nothing to emit
internal/app/roster_set.go:88:  "  (status `masked-by-env`; see `parley roster show --explain %s`)\n"
```
contradicts: `SKILL.md` — *"`STATUS` carries a closed vocabulary: … `masked-by-env` …"*.

why it matters: it is documented as a value the STATUS column carries; it is only ever printed as
advice text by `roster set`. `addStatus("masked-by-env")` exists nowhere. A reader filtering rows on
STATUS will never match it. The source comment at `:83` shows this was already noticed and
half-fixed — the vocabulary was not corrected to match.

## What I checked and found clean

- **The three `COOPERATION.md` copies are in sync.** 105050 / 104805 / 104895 bytes; `diff` between
  the embedded default and the skill snapshot is **one hunk**, and it is the intentional
  `**Transport:**` / `**Created:**` placeholder difference between a bootstrap template and a
  vendor-neutral reference. Section headings are identical. [PRIMARY]
- **The skill snapshot IS guarded**, contrary to what I believed going in:
  `internal/protocol/drift_test.go:276` covers `../../../parley-deck-skill/skills/parley-deck/references/COOPERATION.md`
  as "bundled skill protocol". My prior belief that it was unguarded was stale. [PRIMARY]
- **No undocumented STATUS terms.** Every term the CLI can emit via `addStatus(` appears in
  `SKILL.md`'s closed vocabulary. The mismatch runs only the other way (F3). [PRIMARY]
- **Install target count is honest.** `SKILL.md` says 15; `installer.TARGETS` has 15 keys. [PRIMARY]

## What I could not check, and why

- Whether `SKILL.md`'s **procedural** instructions match runner behaviour end to end. I checked
  command existence and vocabulary, not the phase choreography.
- The `consensus *` verb family — assigned to @opencode-1, and @codex-1 has filed eight findings
  against it. I deliberately did not duplicate that surface.
- `agents/openai.yaml`, `agents/manifest.yaml`, `references/compatibility.json` and
  `references/WORKED_EXAMPLES.md` — unread this round.
- Whether any **other** deck in the fleet has been broken by F2 already. Unmeasured.
