---
idea: skills-cli-install-path
review-round: 17
agent: codex-1
date: 2026-07-30
reviewed-commit: 15ed1ad
responding-to:
  - agy-1/review/round-16
  - codex-1/review/round-16
  - hermes-1/review/round-16
  - kimi-1/review/round-16
---

## Summary

Reviewed `15ed1ad` in the supplied isolated `wt-codex` worktree. The full suite passes
253/253, all 31 supplied probes classify as intended, every concrete round-16 finding is
fixed, and a canonical broken command still runs and fails rather than being universally
refused.

The new publication contract is nevertheless not closed by this implementation.
`publishedTestCommands` does not actually collect the complete rendered text of a block:
it omits partial code-node literals from that view and inserts raw HTML markup as though tags
were visible text. A command assembled across prose and a code span, or across an empty HTML
tag, therefore lands in neither bucket. Three probes rendered visible
`node --test no/such/dir`, stayed green at 12/12, and the copied command exited 1. A second,
lower-impact defect is that the Map can retain only one provenance for duplicate command text,
so a code occurrence masks an invalid prose occurrence. `Status: ❌ BLOCK`.

## Disposition of every round-16 finding

| Round-16 finding | Measured disposition at `15ed1ad` |
| --- | --- |
| `agy-1`: Node options between `node` and `--test` were skipped | **Fixed.** P23 `node --no-warnings --test …`, P24 `node -r … --test …`, and P25 with a valid target all produce `pass=11 fail=1 REFUSED`. Detection is now broader than `SUPPORTED_COMMAND`. |
| `codex-1`: Markdown escape, emphasis, and numeric entity synthesized commands unseen in source | **Fixed for all three exact forms.** P27–P29 each produce `pass=11 fail=1 REFUSED-PROSE`. The parser resolves them into rendered text and provenance is checked before execution. |
| `kimi-1` and `hermes-1`: a literal `>` inside a fence was stripped and the mutated command executed | **Fixed.** P30 produces `pass=11 fail=1 REFUSED`; my independent fenced-content probe does the same. P19, where the fence is actually inside a blockquote, remains `pass=12 fail=0 GREEN`. The parser now makes the marker/content distinction correctly. |
| `kimi-1`: zero-width continuation `node\` + `--test …` was skipped | **Fixed.** P31 produces `pass=11 fail=1 REFUSED`. |
| `hermes-1`: a literal `$ ` prompt was stripped and the mutated command executed | **Fixed.** An independent fenced `$ node --test …` probe produces `pass=11 fail=1` with `refuses to interpret it`. |

No round-16 finding remains open in its original form.

## What I verified and how

### Isolation and supplied checks

The worktree was detached at the requested commit:

```text
$ git rev-parse HEAD
15ed1ad8bdf7651e40787b8b5c07472a626d8955

$ git status --short
?? node_modules
```

`node_modules` is the supplied pre-existing link. No tracked implementation file was
modified. Every temporary `skills/__probe_codex_round17__.md` fixture was removed after its
probe; the final status is still only `?? node_modules`.

The unchanged supplied harness reported:

```text
$ zsh …/scratchpad/probe-codex.sh
baseline                     pass=12 fail=0
P19 blockquote-valid-single  pass=12 fail=0  GREEN
P23 agy-exact-no-warnings    pass=11 fail=1  REFUSED
P24 require-hook             pass=11 fail=1  REFUSED
P25 flags-with-valid-target  pass=11 fail=1  REFUSED
P27 md-escape-backslash      pass=11 fail=1  REFUSED-PROSE
P28 md-emphasis-inside       pass=11 fail=1  REFUSED-PROSE
P29 md-numeric-entity        pass=11 fail=1  REFUSED-PROSE
P30 gt-is-content-in-fence   pass=11 fail=1  REFUSED
P31 zero-width-continuation  pass=11 fail=1  REFUSED
```

Across the complete output: 26 probes were refused, P4/P7/P19/P26 were green controls, and
P8 was the one `RAN-AND-FAILED` control. The harness's concurrency check was not removed, but
this sandbox printed `sysmon request failed` / `pgrep: Cannot get process list`; isolation
instead came from the reviewer-specific worktree and fixture name, sequential execution, and
the absence of the fixture before and after the run.

The full suite reproduces the author's number:

```text
$ npm test
ℹ tests 253
ℹ pass 253
ℹ fail 0
ℹ duration_ms 3465.590458
```

### The guard still executes supported commands

I independently compared content, container, and execution cases:

```text
prompt-is-content        pass=11 fail=1  refuses to interpret it
gt-is-content            pass=11 fail=1  refuses to interpret it
blockquote-container     pass=12 fail=0  green
canonical-broken         pass=11 fail=1  published command failed
```

The last result is a temporary inline code node containing `node --test no/such/dir`. It
reaches the shell and fails as a command, rather than failing the provenance or form checks.
The guard therefore still verifies real commands; it is not a universal refuser.

### CommonMark/GFM and dependency probes

Plain commands in a table, link label, and footnote definition were conservatively refused
as prose even though bare CommonMark does not model every GFM container. A code span inside a
table was still extracted and actually run:

```text
gfm-table-prose            pass=11 fail=1  must be published as its own code span
gfm-table-code             pass=11 fail=1  published command failed
link-label-prose           pass=11 fail=1  must be published as its own code span
footnote-definition        pass=11 fail=1  must be published as its own code span
```

Autolinking changes link structure, not its visible label, so it offers no equivalent
token-removal escape. Raw HTML does; see the MAJOR below.

The parser dependency is a tolerable test-only trade:

```text
commonmark-lock=0.31.2
commonmark-range=^0.31.2
.github/workflows/release-portable.yml:26:      - run: npm ci
```

The lockfile gives CI an exact parser version, and the extractor fixture plus the supplied
harness act as behavioral conformance assertions for the cases they name. An exact
`package.json` pin would not fix the finding below: it is caused by the chosen AST-to-visible-
text model and the CommonMark/GFM/HTML surface, not version drift.

## Findings

### [MAJOR] Mixed inline nodes and raw HTML synthesize a runnable command in neither provenance bucket

The contract says every command not wholly contained in one code node is refused. The AST
walk at `test/design-addons.test.js:306-323` does not enforce that:

- a `code` node is considered independently but its visible literal is omitted from the
  surrounding rendered-text buffer;
- `html_inline` and `html_block` append their raw markup, although HTML tags themselves are
  not visible or copied from the rendered document.

Reproduction:

```markdown
node `--test` no/such/dir
```

Actual output:

```text
mixed-code-span        pass=12 fail=0
rendered: <p>node <code>--test</code> no/such/dir</p>
```

The code node contains only `--test`, so it is not a command. The prose buffer becomes
`node  no/such/dir`, so it is not a command either. The rendered paragraph's visible text is
nevertheless exactly `node --test no/such/dir`.

Raw HTML creates the same false green while directly refuting the source comment that raw
HTML “reaches the reader” as copyable text:

```text
source:   no<span></span>de --test no/such/dir
guard:    pass=12 fail=0
rendered: <p>no<span></span>de --test no/such/dir</p>

source:   node --te<span></span>st no/such/dir
guard:    pass=12 fail=0
rendered: <p>node --te<span></span>st no/such/dir</p>

source:   <div>no<span></span>de --test no/such/dir</div>
guard:    pass=12 fail=0
rendered: <div>no<span></span>de --test no/such/dir</div>
```

For every form, the browser-visible/copyable text is `node --test no/such/dir`. Running that
text independently gives:

```text
$ /bin/sh -c 'node --test no/such/dir'
Could not find 'no/such/dir'
exit=1
```

This is a silent false green and breaks the claim that the class is closed by construction.
It also shows why using bare CommonMark against GitHub/npm/editor rendering is a semantic
choice, not merely a dependency-version choice.

Suggested fix: construct a complete visible-text view for each rendered block, including code
node literals, while separately retaining occurrence-level provenance. For HTML, either use a
GFM-compatible rendering plus HTML `textContent` model, or prohibit raw HTML in the scanned
documents; appending raw tags is not rendered-text collection. Refuse every visible command
occurrence unless its complete text comes from exactly one code node, and add all three forms
above as regression tests.

### [MINOR] A code occurrence masks an invalid prose occurrence with identical text

`publishedTestCommands` returns `Map<command, provenance>`. A command string can be published
more than once with different provenance, but a Map retains only one value. `add()` overwrites
`"prose"` with `"code"`, while `flushRendered()` deliberately declines to overwrite an
existing key.

Reproduction in one document:

```markdown
node \--test "skills/parley-tracker/bin/*.test.js"

`node --test "skills/parley-tracker/bin/*.test.js"`
```

The first line renders a canonical command from prose and should be refused by the stated
provenance-first rule. The second is the same command in an allowed code node. Actual result:

```text
duplicate-provenance       pass=12 fail=0  green
```

The exact command is still executed and succeeds, so this does not create the MAJOR's broken-
command false green; it does violate the publication contract. Store occurrences, or at least
a set of provenances per command, and refuse if any occurrence is not `"code"`.

### Signoff: codex-1 — 2026-07-30
Status: ❌ BLOCK
