---
idea: addon-manifest-coverage
review-cycle: 2
drafted-by: claude-1
date: 2026-08-02
reviewed-commit: e46f661
---

## Where the two review rounds got to

Round 1 produced three MAJOR findings, six MINOR and three NIT against `205416d`. Fix-up cycle 1
applied all twelve. Round 2 confirmed every MAJOR closed — each reviewer checking its own
findings against `e46f661` by execution, not by reading the fix-up notes — and left no MAJOR
standing.

Every round-1 MAJOR was a defect in the implementation, not in the ratified design. Two of the
three were in code written to close a false green and had opened a new hole of their own.

## Agreed fixes

1. **`package-lock.json` still identifies the release as 2.1.0.** Found independently by all
   three reviewers. `package.json`, the executable and the markers say 2.2.0; the lockfile's two
   root `version` fields do not. kimi-1 adds the consequence: `npm ci` on a clean archive exits 0
   despite the mismatch, so nothing breaks — but a plain `npm install --package-lock-only`
   rewrites both entries, so the first install anyone runs dirties the release tree.

2. **`managed: false` contradicts a present and valid marker** when the packaged source cannot be
   enumerated. Raised by hermes-1, reproduced and co-signed by kimi-1 at reduced severity. The
   marker is the evidence for ownership; a source-enumeration failure says nothing about who
   installed the tree. `managed` must be computed from marker evidence and installed-payload
   validity, excluding source-enumeration problems. kimi-1 established the blast radius by
   measurement: it is reporting-only — an unforced `uninstall --only parley-deck` still succeeds
   against the same damaged package, because `installerOwnsDestination` reads the marker from
   disk rather than this field.

3. **The F4 base regression fails at `23a9856` during fixture construction, not on F4.** Raised
   by kimi-1 in round 1, carried forward, and independently raised by codex-1 in round 2. It is
   listed as fix-proving, but at the base commit it cannot even build its fixture, so its failure
   there proves nothing about the migration branch. Either make it discriminate on F4 itself or
   state plainly that it is not runnable at the base commit — silently counting it as
   fix-proving is the same overclaim Amendment 1.1 corrected.

4. **The "every target shape" foreign-copy regression covers two of the three shapes.** kimi-1.
   It asserts codex and gemini; the antigravity shape — the one with a second `skills/SKILL.md` —
   is not exercised, and it is the shape whose staging differs most.

5. **`safeSourceFiles` returns an empty file list for an empty source directory**, so an empty
   packaged skill directory yields no requirement at all. hermes-1. A floor that can be zero is
   not a floor.

## Dismissed findings

**hermes-1's R2-2(a) — that the fail-closed red is a false red — is dismissed**, with hermes-1's
reasoning preserved rather than paraphrased away. Both other reviewers dissented explicitly:

- codex-1: *"That is the deliberate fail-closed trade-off my round-1 fix requested, not a false
  red from a healthy package."*
- kimi-1: *"A managed tree whose packaged source can no longer be enumerated is precisely the
  tree doctor can no longer certify; red with a message naming the exact unreadable source is
  the fail-closed semantics all three of us demanded in round 1, and the alternative (green on
  the shrunken list) is the false green this idea exists to kill."*

hermes-1's own file agrees the direction is right and rates the item MINOR on the strength of the
`managed` half, which is agreed fix 2. What is dismissed is only the claim that a complete
install reporting red against a damaged package is itself a defect. The trade-off is recorded in
`FINAL.md`'s deferred list rather than hidden.

## Deferred follow-ups

1. **Byte-level integrity for the core's managed path** — carried from the design consensus. The
   core catches deleted files but not modified ones.
2. **A natively installed core whose marker was deleted still reports `malformed`** although its
   bytes are correct.
3. **The repair path for a damaged CLI package is not stated in the `doctor` message.** hermes-1
   observed that the message describes a verification limitation rather than an instruction, and
   that re-running install from the damaged package fails at `validatePayload`. kimi-1 calls this
   a documentation nuance rather than a finding. Recorded; not fixed here.
4. **A foreign tree from a different package revision stays `malformed`** — carried from the
   design consensus.

## Unverified claims carried forward

None outstanding. Both claims quarantined earlier were settled by measurement in round 2:
hermes-1's "all existing tests pass with the anchor patch alone" was about the anchor patch in
isolation and is consistent with the six failures the manifests caused; the implementer's
"378/378" was PATH-dependent and is superseded by 382/382 measured under two PATHs with the exact
`npm test` invocation.

## Signoffs

### Signoff: codex-1 — 2026-08-02
Status: ✅ ACCEPT

### Signoff: hermes-1 — 2026-08-02
Status: ✅ ACCEPT

### Signoff: kimi-1 — 2026-08-02
Status: ✅ ACCEPT
