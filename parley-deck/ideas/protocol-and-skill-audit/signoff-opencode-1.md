### Signoff: opencode-1 — 2026-08-21
Status: 🟡 ACCEPT-WITH-RESERVATIONS
Notes: Read consensus.md against my round-01 skeleton (no findings) and the absent round-02 file; §4's account of my participation is accurate and does not invent findings I never filed. §1's verification-split table and structural asymmetry explanation are a fair account of the round-02 record as written. I accept the confirmed/contested split and the three-surface fix framing as the audit ledger, with reservations that do not overturn the bulk of the document: (1) §2's title and §5 say "33 confirmed" while the agent-finding pairs listed under §2 sum to 34 when hermes-1/Q2 is included, and §5's three-surface arithmetic (21+1+7+1+3) silently drops Q2; (2) claude-1/F3 is on the §5 CLI fix list but has no Fixed/Deferred row in IMPLEMENTATION.md and still has no `addStatus("masked-by-env")` path; (3) codex-1/F20 is confirmed and fixed in code/tests yet omitted from IMPLEMENTATION Fixed tables; (4) IMPLEMENTATION applied codex-1/F15 (PARTIAL-only in §3) and deferred kimi-1/F5 (REFUTED in §3), so the shipped disposition set is not identical to the confirmed set. §5's owner rule that no complete release ships until command defects are fixed remains an owner policy, not a claim I re-adjudicate; review/consensus.md is already unanimous ready with five deferred design items recorded. I am not re-verifying the code in this signoff.

## Evidence
- SECONDARY — `parley-deck/ideas/protocol-and-skill-audit/round-01/opencode-1.md`: frontmatter only; empty Findings / clean / could-not-check sections; zero findings filed.
- PRIMARY — `test ! -f parley-deck/ideas/protocol-and-skill-audit/round-02/opencode-1.md && echo ABSENT` → `ABSENT`; `ls parley-deck/ideas/protocol-and-skill-audit/round-02/` lists claude-1, codex-1, hermes-1, kimi-1, zcode-1 only.
- SECONDARY — `consensus.md` §4 lines 82–84: states I created a skeleton and filed no finding; matches the two bullets above; does not misattribute work to me.
- SECONDARY — `consensus.md` §1 lines 17–31: verifier table (codex 6/23 confirmed, kimi 37/42, zcode 30/32) and the asymmetry that codex wrote 24 findings and could not assess them; fair as an account of the filed verification design. I did not re-run those verifiers this session.
- PRIMARY — count of agent-finding IDs listed under `consensus.md` §2: codex 21 + zcode 7 + kimi 3 + claude 2 + hermes 1 = **34**; without hermes = **33**; §5 three-surface sum 21+1+7+1+3 = **33**. Command: inline Python over the ID lists quoted in §2.
- SECONDARY — `consensus.md` §5 lines 88–91: 33 confirmed as fix list; surfaces CLI / protocol / skill; owner release gate on command defects.
- SECONDARY — `IMPLEMENTATION.md` Fixed tables: 28 code/doc rows (22 first batch excluding hermes-1/Q2 data-only + 7 second batch); Deferred: codex-1/F6, F8, F14; kimi-1/F1, F5 (5). Matches claude-1 design signoff's "28 applied / 5 deferred" claim.
- PRIMARY — confirmed vs IMPLEMENTATION disposition set difference: confirmed missing from Fixed+Deferred+hermes-Q2 → `claude-1/F3`, `codex-1/F20`; Fixed not in confirmed → `codex-1/F15`; Deferred not in confirmed → `kimi-1/F5`.
- SECONDARY — `internal/consensus/consensus.go:462-482` and `internal/consensus/roundgate_test.go:120-141` (`TestSignoffsBeforeTheSignoffsHeadingDoNotCount`): codex-1/F20 fixed in tree; IMPLEMENTATION mentions F20 only under mistakes/reversion (lines 132–133, 162), not in Fixed tables.
- PRIMARY — `rg -n 'addStatus\("masked-by-env"\)|masked-by-env' internal/app --glob '*.go'`: hits only `roster_set.go:83` (comment) and `:88` (stderr advice string); no `addStatus("masked-by-env")`. claude-1/F3 as STATUS emission remains open despite CHANGELOG.md:358 wording.
- SECONDARY — `review/consensus.md`: `outstanding_agreed_fixes: 0`, `blocked: false`, five non-implementer ✅ ACCEPT; deferred follow-ups include the same five design deferrals; my prior review signoff is `review/signoff3-opencode-1.md` (separate phase; not re-opened here).
- SECONDARY — `consensus.md` §3: codex-1/F15 listed under PARTIAL-only; kimi-1/F5 REFUTED; neither is in the §2 confirmed set that §5 calls the fix list.

## Re-verification after the ledger correction

Status: 🟡 ACCEPT-WITH-RESERVATIONS
Notes: All four of my prior reservations are addressed in substance (F3 emits STATUS; F20 and F15 dispositioned; F5 dismissed not deferred; §1 table and all-nine-REFUTED-are-codex attribution corrected; Q2 restored to a four-surface split). I still reserve because **the correction itself re-asserts a false confirmed count**: header and §2 title say **36**, the CORRECTED prose claims "true count … is 36: @codex-1 ×23, @zcode-1 ×7, @kimi-1 ×4, @claude-1 ×2, plus @hermes-1/Q2", and those components sum to **37**; the live §2 bullets still read **(21)** and **(3)** and omit F15/F23/F4 from the enumerated lists even while the CORRECTED note says those three belong. §5's four-surface split (codex 23 + claude F3) + (zcode 7 + claude F2) + (kimi 4) + (hermes Q2) = **37**, not 36. Same class of defect I reserved on: a verification-integrity audit cannot close on an unchecked ledger. Not a BLOCK — the substance and my four original points hold; the residual is arithmetic/list hygiene on the correction pass.

### Own objections vs correction

| Prior reservation | Now | Evidence |
| --- | --- | --- |
| (1) 33 vs 34 / three-surface drops Q2 | **Partially fixed** — Q2 and four surfaces present; count still wrong (36 claimed, 37 derived) and bullets stale | PRIMARY re-derive below |
| (2) claude-1/F3 no STATUS path / no Fixed row | **Satisfied** | PRIMARY `rg` + mutation; SECONDARY dispositions table |
| (3) codex-1/F20 omitted from Fixed | **Satisfied** | SECONDARY `consensus.md` dispositions row; tests still present |
| (4) F15 applied though PARTIAL-only; F5 deferred though contested | **Satisfied** | SECONDARY dispositions: F15 CONFIRMED ALREADY FIXED; F5 DISMISSED not deferred |

### 1. Re-derived figures from `round-02/*.md` (PRIMARY)

Command (verdict = trailing `CONFIRMED|PARTIAL|REFUTED` on each `###` heading):

```bash
python3 - <<'PY'
import re, pathlib
from collections import Counter
root = pathlib.Path("parley-deck/ideas/protocol-and-skill-audit/round-02")
for name in ["codex-1","kimi-1","zcode-1"]:
    text = (root/f"{name}.md").read_text()
    hits = []
    for m in re.finditer(r"^###\s+(.+)$", text, re.M):
        line = m.group(1).strip()
        vm = re.search(r"\b(CONFIRMED|PARTIAL|REFUTED)\b(?:\s+as\b.*)?$", line)
        if vm: hits.append(vm.group(1))
    print(name, len(hits), Counter(hits))
PY
```

Result at HEAD `2a5ea3ad781e76a09b3e67bf1cc43144bd02bdfd`:

| verifier | assessed | CONFIRMED | PARTIAL | REFUTED | UNREPRO |
| --- | --- | --- | --- | --- | --- |
| @codex-1 | 23 | 6 | 8 | **9** | 0 |
| @kimi-1 | 42 | **36** | 6 | 0 | 0 |
| @zcode-1 | 32 | 30 | 2 | 0 | 0 |

- §1 table **matches** (was wrongly 37/8/2/1 for kimi).
- All nine REFUTED are @codex-1's, including **kimi-1/F5** and claude-1/F1 and seven zcode findings — matches CORRECTED prose; old "every REFUTED but one" / "kimi authored F5 refutation" is gone.
- 66 of 74 = kimi 36 + zcode 30 over 42+32 → **~89%** (PRIMARY: `python3 -c 'print(round(66/74*100))'` → `89`). Matches CORRECTED; old ~92% was wrong.
- Confirmed set under "≥1 CONFIRMED and 0 REFUTED across round-02" (including claude-1 on hermes-1/Q2): **codex 23 + zcode 7 + kimi 4 + claude 2 + hermes 1 = 37**. PARTIAL-only with no confirm: **codex-1/F12, zcode-1/F4** (and hermes-1/Q1, not in §3). Contested with any REFUTED among the §3 set: seven zcode + claude-1/F1 + kimi-1/F5 + two PARTIAL-only = **11** contested bucket as §3 frames it.
- Therefore: corrected §1 is right; corrected §2/§5 **header count 36 is false** (should be 37 if F15+F23+F4 join the old 34); live bullets `@codex-1 (21)` / `@kimi-1 (3)` **contradict** the CORRECTED note that adds F15, F23, F4.

### 2. Tests (PRIMARY)

- `go test ./...` in repo root → last line `GO_EXIT=0` (all packages `ok` or `[no test files]`; exit code read via `echo "GO_EXIT=$?"` immediately after the test command in the same shell).
- `npm test` in `../parley-deck-skill` → `ℹ pass 391` / `fail 0`, python 54 OK, addon manifest checks OK; `NPM_EXIT=0` via the same `echo` pattern.

### 3. Mutation check — @codex-1/F23 (PRIMARY, own copy)

Shared tree left untouched. Copied the module to
`/var/folders/yt/p2sr23f12_qcfx_w2z5c1p4r0000gn/T/opencode/f23-mut-21nX/cli`, reverted only
`ValidateImplementationArtifact` in `internal/runner/phase58.go` to the pre-fix behavior (any
non-empty status; `strings.Contains` for `## Summary of work`), left tests in place, ran:

```bash
go test ./internal/runner/ -run 'TestValidateImplementationArtifactRejects' -count=1
```

→ **FAIL** `MUT_EXIT=1`:
- `TestValidateImplementationArtifactRejectsAnUnknownStatus`: `` `status: banana` validated as a finished Phase 5 artifact ``
- `TestValidateImplementationArtifactRejectsAnEmptyOrMentionedSummary`: `an empty '## Summary of work' validated`

Fix is load-bearing. (Also confirmed live tree has `addStatus("masked-by-env")` at `internal/app/roster.go:326` and `protocolRole: "source"` in `parley-deck/meta/version.json` — SECONDARY/PRIMARY read.)

### 4. Newly asserted falsehoods in the correction (PRIMARY)

- **False:** "true count under this section's own criterion is **36**" with components 23+7+4+2+1. Sum is **37**.
- **False / stale:** §2 bullets still `@codex-1 (21): … F16–F22, F24` and `@kimi-1 (3): F1…F3` after the CORRECTED note says F15, F23, F4 were mis-recorded as PARTIAL-only and belong in the confirmed set.
- **False if 36 is kept:** §5 four-surface arithmetic implied by the corrected prose is 24+8+4+1 = **37**.
- **True (checked):** §1 table; all nine REFUTED = @codex-1; ~89%; PARTIAL-only list is two (F12, zcode-1/F4); F14 number-struck doubt RETRACTED; dispositions for F23/F4/F2/F3/F15/F20/F5 and `protocolRole: source`; deferred set without F5.

I do not re-open review/consensus.md. Residual reservation is ledger hygiene on the design consensus only.
