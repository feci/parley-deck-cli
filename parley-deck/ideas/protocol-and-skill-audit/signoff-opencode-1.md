### Signoff: opencode-1 — 2026-08-21
Status: ✅ ACCEPT
Notes: Prior reservations (33/34 arithmetic; F3 STATUS path; F20 Fixed omission; F15/F5 disposition mismatch; then the correction-pass 36-vs-37 off-by-one and stale §2 bullets) are all satisfied at HEAD `6cf39f63b613bbaaac5b3f68987e3095d4ff912b` / skill `40ae6f5417dd0275dec8c5384139972ebce46c5e`. Re-derived corpus confirmed = 36; fix list = 37 with Q2 separate; §2 bullets and §5 `24+8+4+1=37` match; claude-1/F1 cites two PARTIAL verifiers; kimi-1/F4 lives in the skill repo at the cited SHA with `recommendedActions` branching on `protocolRole`. Both suites exit 0 unsandboxed.

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

## Final round

Status: ✅ ACCEPT
Notes: Residual 36-vs-37 / stale-bullet reservation is closed. Corpus confirmed and fix-list totals are now stated separately and both re-derive; no open objection remains.

### Own objections vs this amendment

| Prior residual | Now | Evidence |
| --- | --- | --- |
| (R1) claimed 36 with components summing to 37; Q2 folded into corpus | **Satisfied** — §2: "36 from the round-02 corpus … Plus @hermes-1/Q2 … fix-list total is **37**"; CORRECTED TWICE names the prior off-by-one | PRIMARY re-derive below; SECONDARY `consensus.md` §2 L42–55 |
| (R2) §2 bullets still (21)/(3), omitting F15/F23/F4 | **Satisfied** — `@codex-1 (23)` lists F15 and F23; `@kimi-1 (4)` lists F4 | SECONDARY `consensus.md` §2 L60–81 |
| (R3) §5 surface sum must be 37 | **Satisfied** — `24 + 8 + 4 + 1 = 37` | PRIMARY arithmetic; SECONDARY §5 L139–149 |
| (new claims this round) claude-1/F1 two-PARTIAL citations; skill@40ae6f5 holds F4; Seatbelt fact | **True** | PRIMARY citations + `git show` + suites below |

### 1. Re-derived figures from `round-02/*.md` (PRIMARY)

```bash
python3 - <<'PY'
# verdict on each ### heading; confirmed = ≥1 CONFIRMED and 0 REFUTED across codex/kimi/zcode
import re, pathlib
from collections import Counter, defaultdict
root = pathlib.Path("parley-deck/ideas/protocol-and-skill-audit/round-02")
for name in ["codex-1","kimi-1","zcode-1"]:
    text = (root/f"{name}.md").read_text()
    hits = [re.search(r"\b(CONFIRMED|PARTIAL|REFUTED)\b", m.group(1)).group(1)
            for m in re.finditer(r"^###\s+(.+)$", text, re.M)
            if re.search(r"\b(CONFIRMED|PARTIAL|REFUTED)\b", m.group(1))]
    print(name, len(hits), Counter(hits))
verdicts = defaultdict(dict)
for name in ["codex-1","kimi-1","zcode-1"]:
    for m in re.finditer(r"^###\s+(.+)$", (root/f"{name}.md").read_text(), re.M):
        line = m.group(1).strip()
        fm = re.search(r"((?:codex|kimi|zcode|claude)-1/F\d+)", line)
        vm = re.search(r"\b(CONFIRMED|PARTIAL|REFUTED)\b", line)
        if fm and vm: verdicts[fm.group(1)][name] = vm.group(1)
conf = [f for f,vs in verdicts.items() if "CONFIRMED" in vs.values() and "REFUTED" not in vs.values()]
by = defaultdict(list)
for f in conf: by[f.split("/")[0]].append(f.split("/")[1])
print("corpus_confirmed", len(conf), {a: len(by[a]) for a in sorted(by)})
print("partial_only", [f for f,vs in verdicts.items() if set(vs.values())=={"PARTIAL"} or ( "PARTIAL" in vs.values() and "CONFIRMED" not in vs.values() and "REFUTED" not in vs.values())])
PY
```

At HEAD `6cf39f63b613bbaaac5b3f68987e3095d4ff912b`:

| verifier | assessed | CONFIRMED | PARTIAL | REFUTED | UNREPRO |
| --- | --- | --- | --- | --- | --- |
| @codex-1 | 23 | 6 | 8 | **9** | 0 |
| @kimi-1 | 42 | **36** | 6 | 0 | 0 |
| @zcode-1 | 32 | 30 | 2 | 0 | 0 |

- Corpus confirmed (≥1 CONFIRMED, 0 REFUTED): **codex 23 + zcode 7 + kimi 4 + claude 2 = 36**. Matches §2.
- PARTIAL-only (no CONFIRMED, no REFUTED): **codex-1/F12, zcode-1/F4** — two, matches §3.
- Contested with any REFUTED: **9** (claude-1/F1, kimi-1/F5, zcode F2/F3/F5/F9/F10/F12/F13). §3's "11 contested" = those 9 + the 2 PARTIAL-only — framing, not a count error.
- All nine REFUTED are @codex-1's. 66/74 → **89%**. Matches §1.
- @hermes-1/Q2: CONFIRMED only in `round-02/claude-1.md:12` (`hermes-1/F-Q2 … CONFIRMED`); **no codex/kimi/zcode assessment** — correctly excluded from the 36 and added as supplemental → fix list **37**.
- §5: CLI 23+F3=24, protocol 7+F2=8, skill 4, deck Q2=1 → `24+8+4+1=37`. PRIMARY: `python3 -c 'print(24+8+4+1)'` → `37`.

### 2. Newly asserted claims this round (PRIMARY / SECONDARY)

- **claude-1/F1 two PARTIAL citations** — PRIMARY line hits: `round-02/kimi-1.md:295` `### claude-1/F1 — PARTIAL`; `round-02/zcode-1.md:428` `### claude-1/F1 — … — PARTIAL`; codex REFUTED. Matches §3 L97–99.
- **`parley-deck-skill@40ae6f5` contains kimi-1/F4** — PRIMARY: `git -C ../parley-deck-skill rev-parse HEAD` → `40ae6f5417dd0275dec8c5384139972ebce46c5e`; `git -C ../parley-deck-skill show --stat 40ae6f5` → `lib/installer.js` + `test/source-role-advice.test.js` (+80/−1). Diff shows `protocolRole` resolved in `projectStatus` and `recommendedActions` branching: source decks get "do NOT adopt the packaged protocol" rather than the consumer adopt-after-review string. SECONDARY: consensus L226–239 correctly records the prior false citation of CLI `815c93a`.
- **Seatbelt / durable-kill fact** — SECONDARY `consensus.md` L205–223. This host is **unsandboxed**: PRIMARY `sysctl kern.boottime` → exit 0 (boot id readable). No skip required.

### 3. Suites (PRIMARY)

- `go test ./...` in `parley-deck-cli` at `6cf39f6` → all packages `ok` or `[no test files]`; `echo "GO_EXIT=$?"` → **`GO_EXIT=0`**. Did not skip `TestDurableKillEndToEndRealProcess` (sysctl permitted).
- `npm test` in `../parley-deck-skill` at `40ae6f5` → `ℹ pass 391` / `fail 0`, python 54 OK, six addon manifests ok; `echo "NPM_EXIT=$?"` → **`NPM_EXIT=0`**. Includes `a source-role deck is never told to adopt the packaged protocol`.

### 4. What is true / not re-opened

- True: §1 table; 36/37 split; §2 bullets (23)/(7)/(4)/(2)/(1) with F15, F23, F4; §3 two-PARTIAL F1 line; §5 37; dispositions; skill SHA; Seatbelt note; cross-repo citation rule.
- Not re-opened: review/consensus.md; deferred F6/F8/F14/F1; owner release gate; contested set adjudication.

No open reservation. Flip to ✅ ACCEPT.
