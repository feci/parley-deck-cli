---
idea: skills-cli-install-path
review-round: 16
agent: codex-1
date: 2026-07-30
reviewed-commit: c45601f
---

## Summary

Reviewed `c45601f` in isolation. Cycle 16 fixes the exact round-15 defect: all 22 supplied
probes classify as claimed, the deliberately broken single-line control is run rather than
universally refused, and the full suite passes 253/253.

The guard is still unsound. It searches Markdown source for the literal source spelling
`node --test`, but a reader copies rendered text. A CommonMark escape, emphasis delimiter, or
character reference can render that exact command without the literal spelling appearing in
the source. I measured three such fixtures: each left the guard green at 12/12 while the
rendered command exited 1. This directly refutes the cycle-16 claim that reconstruction now
equals what the reader copies.

## What I verified and how

### Supplied battery and full suite

I ran the supplied harness unchanged in its isolated `wt-codex` worktree:

```text
$ zsh /private/tmp/claude-501/-Volumes-My-Shared-Files-AI-WORKSPACE-parley-deck/5dc331bd-5ddf-45e0-b6c2-d519d8c05128/scratchpad/probe-codex.sh
```

```text
baseline                     pass=12 fail=0
P1-P3, P5-P6, P9-P18,
P20-P22                     pass=11 fail=1  REFUSED
P4, P7, P19                 pass=12 fail=0  GREEN
P8 genuinely-broken-path    pass=11 fail=1  RAN-AND-FAILED
```

That is the claimed 18 refusals, three green controls, and one command that genuinely ran and
failed. The last category confirms the guard has not degraded into refusing everything.

I then archived exact commit `c45601f` into a fresh private temporary copy and ran:

```text
$ npm test
ℹ tests 253
ℹ pass 253
ℹ fail 0
```

### Rendered-text probes

I added one temporary shipped Markdown file per case and ran
`node --test test/design-addons.test.js`:

```text
baseline                     pass=12 fail=0
escaped-flag                 pass=12 fail=0 GREEN
emphasized-flag              pass=12 fail=0 GREEN
html-space-entity            pass=12 fail=0 GREEN
```

A local Marked 17.0.6 render measured the source-to-render transformation:

```text
"node \\--test no/such/dir" => <p>node --test no/such/dir</p>
"node --**test** no/such/dir" => <p>node --<strong>test</strong> no/such/dir</p>
"node&#32;--test no/such/dir" => <p>node&#32;--test no/such/dir</p>
```

The first rendered line is already the exact copied command; the second has the same browser
text after the `<strong>` element is flattened, and the third character reference renders as
a space. Running that copied text independently produced:

```text
$ /bin/sh -c 'node --test no/such/dir'
exit=1
Could not find 'no/such/dir'
```

### Named trade-off

Refusing a legitimate backslash-continued command is acceptable here. P22 uses a valid target
but fails by the named refusal (`11 pass / 1 fail`), no shipped file uses that form, and a
single-line publication rule is reasonable. The defect below is different: it is silent
success, not conservative refusal.

All temporary files and private copies were removed. I did not modify tracked implementation
files.

## Findings

### [MAJOR] Markdown rendering can synthesize a broken command that the source scanner never sees

`publishedTestCommands` prefilters raw source with `/node\s+--test/` at
`test/design-addons.test.js:273`. `logicalLines` models continuations and blockquote markers,
but Markdown itself also removes or replaces inline syntax before the reader sees the page.
The prefilter silently skips those rendered commands, so neither `SUPPORTED_COMMAND` nor the
shell runner gets a chance to refuse or execute them.

Reproduction — publish this standalone Markdown line:

```markdown
node \--test no/such/dir
```

The backslash is a Markdown escape and is absent from the rendered page. The reader sees and
copies `node --test no/such/dir`, which exits 1. The guard reports `12 pass / 0 fail`. These
equivalent probes also remained green:

```markdown
node --**test** no/such/dir
node&#32;--test no/such/dir
```

This is an approach defect, not a request for three more regex branches. Sixteen cycles have
shown that reconstructing rendered Markdown with a hand-written line scanner does not close
an open syntax. Narrow the publication contract instead: require every verification command
to occupy a complete code node on one physical line (a whole inline code span or one fenced
code line), reject continued forms, and use a real CommonMark/GFM parser to extract those
nodes and reject command-bearing rendered text outside them. An even smaller alternative is
a machine-readable command registry from which the documentation is generated. In either
design, execute the exact registered/code-node text; do not approximate Markdown rendering.

### Signoff: codex-1 — 2026-07-30
Status: ❌ BLOCK
