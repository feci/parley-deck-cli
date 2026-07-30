# Round 17 review — kimi-1

## Summary

The rewrite is real, and it is better than anything in cycles 03–16. The parser earns its
place on the container distinction it was bought for: I measured a fence *inside* a blockquote
yielding the bare command and **running**, and a `>` *inside* a fence staying in the literal
and being **refused** — the asymmetry that defeated cycle 16's stripper. Every round-16
finding is genuinely fixed (dispositions below, each measured, none taken on trust). The suite
is 253/253 and all 31 probes in my battery behave as the author claims. The guard still
verifies: I confirmed independently that a genuinely broken path **runs and fails** through
two different publication forms, and that valid forms execute real tests.

But the class is not closed, and it is not closed by construction. The guard's prose net scans
a string that is neither the source nor the rendering: text-node literals with **raw
`html_inline`/`html_block` literals spliced back in**. Renderers do not show that string. They
delete comments entirely and render allowed tags (`<ins>`, `<del>`, `<s>`, `<strike>`) as
formatting. So inline HTML becomes an invisible splice point that breaks detection for the
guard and not for the reader — and the same is true of GFM strikethrough, which the chosen
parser does not implement at all. One line of markdown puts a runnable
`node --test no/such/dir` in front of a GitHub/npm/editor reader while the guard sees nothing,
refuses nothing, and goes green. That refutes, in writing, the contract ("anything else that
renders as such a command is refused by name") and the header comment ("rendered text outside
code nodes is collected separately, with escapes, entities and emphasis already resolved by
the parser"). One MAJOR, one MINOR.

## Dispositions of the round-16 findings — verified, not trusted

| finding | verdict | evidence |
|---|---|---|
| `agy-1`: node flag between tokens was skipped, and skipping reads as success | **FIXED** | P23/P24/P25 now `REFUSED` — detected by the broadened `mentionsATestCommand`, then refused by the grammar. Fixture asserts both halves (`design-addons.test.js:458-459` detection, `:507-508` refusal). |
| `codex-1`: rendering synthesizes commands (escape, emphasis, entity) | **FIXED** — for the transforms CommonMark resolves | P27/P28/P29 all `REFUSED-PROSE`. AST check: `\--` → text `--`, `--**test**` → text `--` + `test`, `&#32;` → text ` `. But see MAJOR: the class survives through channels the parser does not resolve. |
| `kimi-1` + `hermes-1`: `>` stripped, mutated text executed | **FIXED** | P30 (`> node --test …` inside a fence) → `REFUSED`; the `>` stays in the `code_block` literal and `SUPPORTED_COMMAND` rejects it. P19 proves the other arm: fence inside a blockquote yields the bare command and runs green. |
| `kimi-1`: zero-width continuation boundary | **FIXED** | P31 → `REFUSED`; the spliced unit is emitted as `node--test cont/zero-width.test.js \` and the grammar refuses the backslash (fixture `:473`). |
| `hermes-1`: `$ ` prompt repaired instead of refused | **FIXED** | fixture captures `"$ node --test b/dir"` whole (`:444`) and asserts the grammar refuses it (`:498`). |

The author's claim about the *old* fixture also checks out: `git show c45601f` shows the
`` ```not-a-closing-fence `` line — the blockquote cases sat inside an outer fence that never
closed where it looked closed, and the line scanner could not tell. The new fixture is
well-formed; I re-parsed it and every block closes where it appears to.

## What I verified, and how

- **Suite claim**: `npm test` in my worktree at `15ed1ad` → `pass 253, fail 0`. ✓
- **Probe claim**: `zsh probe-kimi.sh` → all 31 probes behave as intended: every refused shape
  `REFUSED`/`REFUSED-PROSE`, P4/P7/P19/P26 `GREEN`, P8 `RAN-AND-FAILED`. ✓
- **The guard still verifies** (independent of the author's battery): ad-hoc probe
  `> ```bash / > node --test no/such/dir / > ``` ` → `RAN-AND-FAILED` — the run path is real
  through a container, not just the refusal path. P4/P19 exercise the green run path, and the
  executor asserts `pass > 0` from the child's own summary, so green cannot be faked by zero
  tests. The two shipped commands are `skills/parley-tracker/templates/subtask.md` and
  `skills/parley-design-check/SKILL.md:372`; `published.size >= 2` anchors extraction.
- **Provenance buckets**: walked every node type the `commonmark` walker emits. Text-bearing
  leaves are exactly `text`, `code`, `code_block`, `html_inline`, `html_block` — all five are
  collected, so nothing lands in *neither* bucket. The failure is the opposite: see MINOR.
- **GFM containers around code nodes** (tables, task lists, footnotes, autolinks): the parser
  reports the code node with the same literal the GFM reader copies; broken targets
  `RAN-AND-FAILED` honestly (probes D1–D4). No divergence there. CRLF in fences is normalized
  by the parser (`code_block:"node --test x\n"`). The code-node side is tight because code
  nodes are literal in both dialects.
- **The dependency**: `commonmark@0.31.2` pinned exactly in `package-lock.json`; CI runs
  `npm ci` (`.github/workflows/release-portable.yml:26`); the extractor fixture pins the parse
  behavior the guard relies on. Tolerable trade — but note the MAJOR is a *dialect* limitation
  no version bump of this package will fix: `commonmark` implements CommonMark, and its
  conformance target is precisely what makes it blind here.

## Findings

### [MAJOR] The prose net scans a string no renderer produces — inline HTML and GFM strikethrough are invisible splice points

The guard builds its "rendered" prose line by joining `text` literals **with raw
`html_inline`/`html_block` literals spliced in** (`design-addons.test.js:320-323`), then runs
`mentionsATestCommand` over the join. That join is not what any renderer shows. Renderers
delete HTML comments entirely, and GitHub/npm/editor previews pass `<ins>`, `<del>`, `<s>`,
`<strike>` through sanitization as formatting — the tag *text* vanishes from the page while
its content stays. So a tag or comment placed between the tokens breaks `--test` for the
guard's regex and not for the reader's eye. This needs no GFM at all — inline HTML is core
CommonMark.

Reproduction (each is the entire body of a one-file probe dropped into `skills/`, run through
the guard exactly as `probe-kimi.sh` does):

```
node --t<!-- -->est no/such/dir
```

AST (measured): `text:"node --t"` | `html_inline:"<!-- -->"` | `text:"est no/such/dir"`.
The joined "rendered" line is `node --t<!-- -->est no/such/dir`, which contains no `--test` —
so the unit is **not detected at all: not run, not refused**. Measured result:
`pass=12 fail=0` — GREEN. On GitHub, npm, and in editor preview the comment renders as
nothing and the page shows a pristine `node --test no/such/dir`. The reader copies a broken
command that the guard certified by silence.

Same mechanism, through tags GitHub renders as formatting (copy includes the styled text —
the exact mechanism of `codex-1`'s ratified round-16 emphasis finding):

```
node --t<ins>e</ins>st no/such/dir                         → pass=12 fail=0  GREEN
node --t<del>e</del>st "skills/parley-tracker/bin/*.test.js" → pass=12 fail=0  GREEN
```

The third line is the sharpest: the reader gets the **canonical, working** command out of
prose, and it is neither refused (as the contract demands) nor run — the provenance-before-form
rule never fires because detection never fires.

And the same root cause through the channel the brief pointed at: this parser is CommonMark,
the readers are GFM. Strikethrough (`~~`, GFM §6.5 — rendered by GitHub, npm's
markdown-it-based renderer, and VS Code preview) is not implemented by `commonmark@0.31.2`:

```
node --t~~e~~st no/such/dir    → text:"node --t~~e~~st no/such/dir" → pass=12 fail=0  GREEN
node ~~--~~test no/such/dir    → pass=12 fail=0  GREEN
```

GitHub renders the first as `node --t` + struck `e` + `st no/such/dir`; copying the line
yields `node --test no/such/dir`. The parser resolves nothing, the regex sees nothing.

What it refutes, in writing: the contract ("anything else that renders as such a command is
refused by name" — these render as such a command and are not even seen), the header comment
("rendered text outside code nodes is collected separately, with escapes, entities and
emphasis already resolved by the parser" — resolved, then re-contaminated by raw HTML
literals; and strikethrough is GFM's emphasis extension), and the design note justifying the
splice ("Raw HTML is not rendered away — it reaches the reader" — comments are rendered away
entirely; formatting tags are rendered away *as text*). Cycle 18 closed `codex-1`'s class only
for the transforms this parser happens to resolve; the rendered line it compares against is
still an approximation — now an approximation of HTML rendering instead of markdown rendering.

Direction, not prescription: build the prose line the way a renderer does — drop comment and
tag literals so fragments join (`--t` + `est` → detected → refused as prose), while still
policing whole commands whose text lives inside visible raw HTML (`<pre>`, `<code>`); and
either parse a GFM dialect or normalize `~~…~~` in prose lines. Both channels must turn
*detected-and-refused*, not merely undocumented.

### [MINOR] The same wrong premise fails the build on text no reader can see — and runs a command out of image alt text

The flip side of "raw HTML reaches the reader": text that is genuinely invisible is policed
as prose. Measured:

```
<!-- maintenance note: the check is node --test skills/foo/bin/*.test.js -->   → REFUSED-PROSE (fail)
<script> var cmd = "node --test no/such/dir"; </script>                        → REFUSED-PROSE (fail)
```

No renderer shows either to any reader, yet each fails the suite with "must be published as
its own code span". And on the code side, `![`node --test no/such/dir`](img.png)` reports the
alt text's `code` node (measured AST: `image` → `code "node --test no/such/dir"`), so the
guard would *execute* a command that renders as a picture. All three fail closed, no shipped
file trips them today (the suite is green, which proves absence), and the fix for the MAJOR
likely subsumes the first two — but the premise in the code comment is false in both
directions, and the image-alt arm executes invisible text, which is the direction this guard
exists to distrust.

### Signoff: kimi-1 — 2026-07-30
Status: ❌ BLOCK
