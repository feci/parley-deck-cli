---
agent: claude-1
idea: meta-protocol-change-evidence-first-efficiency
phase: implementation (Phase 5 slice handoff)
date: 2026-09-05
slice: packet renderer, packet CLI helper, applicability map, source protocol §9 item 1, skill standing instructions
status: implemented on branch; NOT independently verified; integration glue pending (listed below)
---

## Scope delivered (AC-P1 slice, claude-1 claims only)

One live-source packet renderer (`internal/protocolpacket`), the `parley protocol packet` /
`packet check` command helper (`internal/app/protocol_packet.go`), the ratified applicability map
(`parley-deck/meta/packet-applicability.yaml`), the §9 item 1 source change with its changelog
entry, and the skill standing instructions. Default is full context with a shadow packet audit
record; the optimized packet is emitted only on explicit `Optimize`/`--optimize`, which is the
experimental input of the ratified trial and not an enabled release. Nothing here runs the
packet experiment (AC-P2), wires prompt builders, or publishes anything.

I issue no verdict on my own slice (§15.1). Every claim below is a description for a
non-owner to verify against the branch.

## Files (CLI worktree, branch feature/meta-protocol-change-evidence-first-efficiency/claude-1)

- `internal/protocolpacket/packet.go` — types, `Parse` (heading blocks, fence-aware, unique
  locators), `Build` (guards, dependency closure, index, verbatim packet body, shadow),
  `DetectSecret`, `Hash`, `Locators`.
- `internal/protocolpacket/applicability.go` — map schema `parley.packet-applicability/v1`,
  `ParseMap`/`LoadMap` (strict, KnownFields, one document), `Check` (completeness, stale
  entries, dependencies, never-cut list, render of every phase×track), `Report`.
- `internal/protocolpacket/source.go` — `ResolveSource` (role-aware authority), `Render`
  (one-call entry), `WriteBody` (under `.parley-runtime/protocol-packets/`), `RuntimeDir`,
  `MapPath`, `DeckProtocolPath`, `ErrAuthority`.
- `internal/protocolpacket/packet_test.go` — 12 tests (below).
- `internal/app/protocol_packet.go` — `runProtocolPacket`, `protocolPacketCheck`,
  `BuildProtocolContext`.
- `internal/app/protocol_packet_test.go` — 7 tests.
- `parley-deck/meta/packet-applicability.yaml` — 69 entries, one per heading of the live
  protocol; never-cut blocks `always`; phase/transport/flag conditions with triggers.
- `parley-deck/COOPERATION.md` — §9 item 1 only (one line). No §4.0 cell, no §2 edit.
- `internal/protocol/defaults/COOPERATION.md` — the SAME one-line §9 mirror. Recorded in my
  allocation row: `internal/protocol/drift_test.go` binds the embedded default byte-for-byte
  to the deck outside five allowlisted zones, so the deck edit alone fails the suite. Reject
  the mirror at integration if the facilitator disagrees; the drift test then fails by design.
- `parley-deck/meta/protocol-changelog.md` — new top entry, 2026-09-05.
- `IMPLEMENTATION.md` — only my two allocation rows changed (file claims + status).

## Files (skill worktree, same branch name, ../worktrees/evidence-first-skill)

- `skills/parley-deck/SKILL.md` — core rule line, "Required Protocol Context" (renderer first
  with attestation; full read = `full-fallback` with reason; bundled snapshot =
  `full-fallback:bundled-snapshot`; `refused` never proceeds with substituted text), startup
  flow step 1. The manual `shasum` drift check and the "stop and ask" branch are preserved.
- `skills/parley-deck/references/COOPERATION.md` — the same §9 item 1 line as the deck.
- `test/packet-context.test.js` — 6 source tests.
- `skills/parley-deck/parley-addon.json` — GENERATED per-file hash manifest, regenerated with
  `node scripts/build-addon-manifest.js` because the two source edits made it stale (the
  installer suite verifies it). Recorded in my allocation row; not a hand edit.

## Behaviour summary (what a reviewer should try to break)

- Authority: `protocolRole: source` → the live deck file is the source; the core store is
  never consulted and a stale core is never substituted. `consumer` → lock parsed, release
  loaded, body hash compared, overlay reconciled, `protocolcore.Render` compared with the
  on-disk view; any failure is `ErrAuthority` (nothing emitted); a verified chain whose view
  differs is `Drift`, which makes an optimized request `full-fallback` with reason `drift:…`.
  Missing/invalid `meta/version.json`, unknown role, missing protocol file → `ErrAuthority`.
- Secrets: a credential-shaped token anywhere in the source → `context_mode=refused`,
  `fallback_reason=secret-detected:<shape>`, no body, no `packet_sha256`, nothing written.
  The live protocol is asserted clean by test, otherwise every launch would refuse.
- Fallback reasons (optimized request): `unknown-phase|track|transport|flag`, `drift`,
  `no-applicability-map`, `malformed-map`, `stale-map:<locator>`, `unclassified:<locator>`,
  `dependency-unresolved:<locator>`. The body is then the full source verbatim and the
  unknown block is recorded `classification=unknown, included=true` in the index.
- Default (no `Optimize`): `context_mode=full`, body = source bytes, `packet_sha256 =
  source_sha256`, plus `shadow` {packet hash, bytes, included/omitted counts, would-fallback}.
- Packet body: verbatim blocks in source order, then "Packet omission index" listing every
  omitted block with locator, line range, block sha256, classification and trigger. Index in
  JSON covers every block exactly once (tested disjoint and complete).
- `protocol_change` flag is derived from a `meta-protocol-change-*` slug so §7 binds in
  every phase of a protocol-change idea.

## Compact API for Codex (runner/handoff wiring — NOT done by me)

```go
import "parley-deck-cli/internal/protocolpacket"

req := protocolpacket.Request{Phase: 6, Track: idea.Track, Transport: "", // "" = deck header
    IdeaSlug: idea.Slug, Flags: []string{"strict_gate"}, Optimize: false}  // false = full/shadow default
// runner side (cannot import app):
ctx, err := protocolpacket.Render(root, protocolcore.StoreAt(config.CentralHome()), req)
// app side:
ctx, err := app.BuildProtocolContext(root, req)
// err != nil → nothing attested (authority/I/O): record a tooling:packet failed start, do not launch.
// ctx.ContextMode == protocolpacket.ModeRefused → do NOT launch with substituted text.
// otherwise: prompt body = ctx.Body (or read ctx.BodyPath); attestation = ctx.Attestation
//   {context_mode, source_sha256, packet_sha256, fallback_reason} → requested/launch record.
```

Exact call sites Codex owns (I did not edit them): `internal/runner/runner.go`
`BuildRoundOnePrompt`/`buildPromptForRound`/`BuildRoundPrompt`, `internal/runner/phase58.go`
`BuildImplementationPrompt`/`BuildReviewPrompt`/`BuildReviewConsensusPrompt`/`BuildFixupPrompt`,
`internal/runner/consult.go`, `internal/runner/steer.go`, `internal/runner/handoff.go`
`WriteHandoffPacket` (write `ctx.Body` beside `handoff-prompt.md` only if that dir is ignored;
otherwise reference `ctx.BodyPath`). Registration: in `internal/app/protocol.go` `runProtocol`,
add `case "packet": return runProtocolPacket(rest, stdout, stderr)` in the `switch sub` that
precedes the shared flag set (the subcommand owns its flags), and add the usage lines from
`protocolPacketUsage` to `protocolUsage`/`printUsage`. `.gitignore`: add
`/.parley-runtime/`. None of this is claimed done; `parley protocol packet` is unreachable
from the binary until that line lands.

## Tests executed (exact commands, this worktree, 2026-09-05)

- `go test ./internal/protocolpacket/` → ok (12 tests: parse coverage/fence/duplicates; live
  map complete + never-cut + renders all 27 phase×track cells; verbatim + disjoint/complete
  index; default full+shadow; unknown applicability; stale map + dependencies; unexpected
  phase/track/transport/flag + drift; secret refusal; malformed maps; never-cut violation
  detection; authority roles incl. consumer verified/drift/no-lock/bad-hash/missing-release;
  runtime-dir write + malformed map + missing authority).
- `go test ./internal/protocol/` → ok (drift guard passes with the mirror).
- `go test ./internal/app/ -run 'ProtocolPacket|TestProtocol'` → ok (7 new tests incl.
  `packet check` on this repository's live deck, which writes nothing).
- `go test ./...` → 27 packages ok, 0 failures. `go vet` on both packages clean. `gofmt -l` clean.
- Skill: `node --test test/packet-context.test.js` → 6/6 pass.
  `node scripts/build-addon-manifest.js --check` → parley-deck ok after regeneration.
  Full `node --test` → 383 pass, 1 fail: `test/design-addons.test.js` aborts at load with
  `Cannot find module 'commonmark'` (devDependency; no `node_modules` in this worktree).
  It fails identically on a `git archive HEAD` copy without dependencies, and a HEAD copy
  with `npm ci` passed 391/391 before my edits. I did not install dependencies here (no
  network beyond ordinary resolution was authorized for the skill repo), so the full skill
  suite on the edited tree with dependencies present is NOT verified by me.

## Remaining glue and limits (honest)

1. Runner/handoff/prompt-builder wiring, `runProtocol` dispatch, usage text, `.gitignore`
   entry: Codex-owned, not done. The renderer is not on any launch path yet.
2. Attestation into telemetry records (`context_mode`, `source_sha256`, `packet_sha256`,
   `fallback_reason`): Codex's telemetry slice consumes `ctx.Attestation`; not done.
3. `ResolveSource` duplicates the lock/release/overlay sequence of `app.resolveDeck`
   (`internal/app/protocol.go:139-162`) because the runner cannot import `app`. A later
   unification into `protocolcore` is a Codex-owned refactor; behaviour is the same.
4. `Check` proves structure only. A wrong classification present in the map is undetectable
   by code; that is why the default stays full/shadow and why the map is §7 protocol.
5. Consumer-role authority is tested only with temp-dir fixtures; no real consumer deck was
   exercised. The overlay path is exercised through `ReconcileOverlay` with absent overlay
   only.
6. Heading-level granularity: a phase block includes its whole section; "current-phase
   mechanics" inside a transport subsection are included by including the whole active
   transport subsection.
7. No experiment, no pilot, no measurement, no model launch was performed. AC-P2 is Codex's
   later obligation. No self PASS is issued for the implementation.
8. Skill changes are source-only: no installed skill, package version or global core was
   modified. Commit SHAs are reported in the handoff message, not here.
