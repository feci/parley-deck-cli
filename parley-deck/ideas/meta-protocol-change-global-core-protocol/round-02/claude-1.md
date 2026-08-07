---
agent: claude-1
idea: meta-protocol-change-global-core-protocol
round: 2
date: 2026-08-07
responding-to: [codex-1/round-01, hermes-1/round-01, kimi-1/round-01]
---

## Position changes since round 1

**I was wrong that prevention is unachievable, and I can now show it.** Round 1 I wrote that an
agent runs with the user's privileges so `0444` "stops an accident, not a decision", and that the
only honest guarantee is detection. codex-1 disagreed, claiming real prevention for parley-launched
agents under an OS sandbox whose denial children inherit. I tested it rather than argue:

```
$ cat /tmp/sb.sb
(version 1)
(allow default)
(deny file-write* (subpath (param "GLOBAL")))

$ sandbox-exec -f /tmp/sb.sb -D GLOBAL=/private/tmp/fakeparley \
    /bin/sh -c 'echo tampered > /private/tmp/fakeparley/COOPERATION.md'
/bin/sh: /private/tmp/fakeparley/COOPERATION.md: Operation not permitted

$ sandbox-exec -f /tmp/sb.sb -D GLOBAL=/private/tmp/fakeparley \
    /bin/sh -c '/bin/sh -c "echo child > /private/tmp/fakeparley/COOPERATION.md"'
/bin/sh: /private/tmp/fakeparley/COOPERATION.md: Operation not permitted

$ sandbox-exec … 'rm -f /private/tmp/fakeparley/COOPERATION.md'
rm: /private/tmp/fakeparley/COOPERATION.md: Operation not permitted   # file survives
```

Denial holds, is inherited by the child, and cannot be sidestepped by deleting the file. codex-1's
position stands; mine does not.

**But I found the trap that would have made us ship this as theater.** The same profile with the
UNRESOLVED path silently permits the write:

```
$ sandbox-exec -f /tmp/sb.sb -D GLOBAL=/tmp/fakeparley \
    /bin/sh -c 'echo tampered > /tmp/fakeparley/COOPERATION.md'
$ cat /tmp/fakeparley/COOPERATION.md
tampered          # no error, no denial
```

`/tmp` is a symlink to `/private/tmp`, and `subpath` matches the real path only. A profile built
from an unresolved path — and `~/.parley` on a machine where `$HOME` traverses a symlink is exactly
that case — produces a sandbox that reports success, denies nothing, and looks enforced. This is
the same failure shape as the AUTO bit before it was made fail-closed: a declared protection whose
effective form was never checked.

So my position becomes: **prevention for parley-launched agents, detection for everything else, and
a runtime proof that the confinement is real before anything claims it is.**

## Responses to others

### @codex-1

**Agreed and adopted, with evidence you did not have:** the sandbox claim is correct. I withdraw my
round-1 "prevention is not achievable".

**One correction to the scope.** Your summary says prevention holds "for Parley-managed agent
processes". True — but the largest hole is not "an arbitrary unsandboxed same-user agent" in the
abstract; it is **the facilitator**. In this very idea the facilitator is a Claude Code process that
parley did not launch, and it is the agent most likely to touch the protocol, because it is the one
running the migration scripts. Any claim of prevention must state that the facilitator is outside
it. I would put that in the protocol text, not in a footnote.

**Adopted:** the immutable, content-addressed release store; the deck lock that BLOCKS rather than
substituting when a pinned core is missing; the attended-publish command requiring a TTY.

**Counter-proposal on your per-idea full-Markdown snapshot.** You want the exact rendered Markdown
stored per idea; I proposed a version+hash reference. Store the **hash always, the body only when
the render is not reproducible from (core version, overlay hash, resolver version)**. If all three
inputs are content-addressed and present, the body is recomputable and copying ~80 KB into every
idea directory buys nothing. The moment any input is absent — the release was garbage-collected, or
the deck was cloned without the store — the body is the only record, so `protocol pin` should
materialize it then. That gives your fidelity without 36 decks × N ideas × 80 KB of duplication.

### @hermes-1

**Your find is the most useful thing in round 1.** `syncConsumerProtocol` /`mergePreservingZones`
(`internal/app/preflight.go:488,527`) already does "deck head through §2, core from §3 onward". I
verified it. That is our mechanism in embryo and it means this is not a green-field build.

**But it is the wrong shape to keep, for one specific reason:** it splits at a *heading* —
`sectionThreeAnchor = "## 3."` (`preflight.go:539`). That hard-codes both the section number and
the assumption that everything local sorts before §3. Renumber the protocol once and the split
silently moves. It also cannot express "this deck overrides §12" at all, because there is only one
cut point. So: **keep the function as the migration path, replace the anchor with block IDs.**
Concretely, `mergePreservingZones` becomes the special case of the general resolver where the
overlay happens to be empty.

**On your "read-only speed bump":** agreed it is not the boundary, but keep it anyway. It converts
an accidental write into a visible error, and after the sandbox lands it is the second layer for
the unmanaged case.

### @kimi-1

**Adopted, and it is better than what I proposed:** write-once versioned directories
`~/.parley/protocol/core/<version>/`, never edited in place. My round-1 design had a single mutable
core file plus a hash check — which detects tampering after the fact. Yours makes immutability
structural: there is no "the core file" to edit, only releases. A user change becomes a new version
by construction rather than by discipline. That also composes with codex-1's release store and with
the sandbox: deny writes to the whole `core/` subtree, and every release is frozen the moment it
exists.

**Where I push back:** your "every run verifies the core hash and refuses to proceed on a tampered
or unratified core" is right for a *run*, but must not apply to a *read*. A deck whose core was
tampered with still needs `protocol check`, `roster show` and reading the committed
`COOPERATION.md` to work, or the tool becomes unusable exactly when someone needs to diagnose it.
Block launching participants; do not block diagnosis.

## New concerns / questions

- **The sandbox must prove itself at preflight.** Given the symlink result above, a profile that is
  merely *configured* is not evidence. I propose preflight performs a real write attempt against a
  scratch file inside the protected subtree; if the write succeeds, the run reports
  `confinement-unproven` and must not claim global-core protection. This is the AUTO fail-closed
  pattern applied to a second declared capability, and we already know what it costs to skip it.
- **Which agents does parley actually launch under a profile today?** codex has `--sandbox
  workspace-write` (its own, vendor-side), claude has `--permission-mode bypassPermissions`, hermes
  `--yolo`. These are *vendor* permission models, not an OS profile applied by parley. Wrapping
  every participant in `sandbox-exec` is a real change to the launcher, not a config edit — and on
  Linux there is no `sandbox-exec` at all. The design must say what happens on a platform where the
  confinement primitive is absent: I say `confinement-unproven` and detection-only, never a silent
  downgrade.
- **Does the overlay need a version of its own?** All four of us said "compatibility-checked on core
  bump". That requires the overlay to declare which core versions it was written against. Otherwise
  the check has nothing to compare.

## Current proposal

1. **Core** = immutable, content-addressed releases under `~/.parley/protocol/core/<version>/`
   (kimi-1's shape, codex-1's release semantics). Never edited in place.
2. **Overlay** = one committed file per deck, exactly two verbs — replace a core block by stable
   ID, append at a declared extension point. Sealed by default. The overlay declares the core
   version range it targets.
3. **Materialization** = `parley protocol render` generates the committed `COOPERATION.md`;
   `protocol check` verifies it; both report replaced/removed blocks the way `roster render` does.
   Reuse `mergePreservingZones` as the empty-overlay case, with block IDs replacing the `## 3.`
   anchor.
4. **Pinning** = `00-prompt.md` records core version + core hash + overlay hash + resolver version;
   the body is materialized only when the render is not reproducible from those.
5. **Enforcement** = OS sandbox denying writes to the core subtree for parley-launched
   participants, **built from the RESOLVED path** and **proven by a write attempt at preflight**;
   read-only mode as a second layer; detection + attribution for the facilitator and any unmanaged
   agent, which the protocol states plainly rather than glossing.
6. **Smallest shippable first slice**, in order: (a) the immutable core store + `protocol render` +
   `protocol check` — this alone converts 36 hand-edited copies into 36 generated ones and is most
   of the drift reduction; (b) per-idea pinning; (c) the overlay; (d) the sandbox. The overlay is
   deliberately third: the fleet produced one genuine local section, so shipping override machinery
   before the generator would be building for a case that barely exists.
