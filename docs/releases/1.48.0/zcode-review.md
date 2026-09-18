---
agent: zcode
idea: meta-protocol-change-evidence-first-efficiency
review-round: 1
date: 2026-09-18
reviewed-commit: cli head 646254cd39d1261ccaeaf0e92eaaf5fae04c095b plus release edits (candidate tree /private/tmp/parley-release-20260918-3jxrzfsf; not a git checkout, so the head sha is taken from the launch brief)
reviewer-role: release reviewer (Phase 6, deliberation track, github-pr transport)
protocol-context: full (attestation evidence/review-context.json: context_mode=full, source_sha256 == packet_sha256 4519258c…, body read at cli/.parley-runtime/protocol-packets/full-phase6-deliberation-4519258c….md)
---

# ZCode release review — parley-deck-cli 1.48.0 + parley-deck-skill 2.12.0

## Summary

The release candidate is in good shape. The renderer (`parley protocol packet`) is genuinely fail-closed where it matters — full context is the verified code default, refusals emit nothing, fallbacks are visible with reasons — and every release document (CHANGELOG, release notes, SKILL.md) consistently refuses to claim the optimized packet or the empirical audit as proven. I found one MAJOR defect: the skill's unchanged Core Rule still tells every facilitator to "Always read `parley-deck/COOPERATION.md` first", which contradicts the new renderer-first, attestation-required launch flow this same release teaches — a two-line doc fix that should land before tagging. Everything else is MINOR or cosmetic. Recommendation: **release (GO)**, ideally with MAJOR-1 fixed first.

## Refutation attempts

- **Try to break "refused leaks a packet hash":** traced `Build` early returns (packet.go:244-254) — both refusal paths return before `ctx.PacketSHA256 = Hash(ctx.Body)` at packet.go:397, so a refused attestation carries no body hash; `WriteBody` independently refuses refused/empty bodies (source.go:169-175). Could not break it.
- **Try to break "full is only the documented default":** searched for any caller that sets `Optimize` by default. `runProtocolPacket` defaults `--optimize` to false (protocol_packet.go:56); `Build` emits `ModeFull` + shadow audit record unless `Request.Optimize` (packet.go:369-395); the shadow record is counters and would-fallback reason only, never a body swap. The live attestation for this very review (review-context.json) shows `context_mode: full` with packet sha equal to source sha. Could not break it.
- **Try to break the skill's new instructions by finding a contradiction:** succeeded — see MAJOR-1 (SKILL.md:12 vs SKILL.md:22-50).
- **Try to break "the hard-coded `Transport: github-pr` in bundled templates is a regression":** partially succeeded in the CLI's favor — `parley init` uses that exact line as its substitution anchor when writing a new deck (workspace.go:111, workspace.go:138), so in the CLI default template the value is deliberate machinery. The skill snapshot finding survives in reduced form (MINOR-2) because the portable reference has no substitution step behind it.
- **Try to find a packaging blocker:** `prepack` runs the addon-manifest sha256 check (package.json:66) and the recorded release evidence shows `npm pack` succeeding for 2.12.0 with the aggregate `7ffd138c…` matching parley-addon.json, and the portable build running all five pkg targets (evidence/npm-pack.json, evidence/portable-build.log — recorded by the release tooling, not executed by me). No blocker found in package.json `files`/`bin`/`engines` or in build-portable.js target handling; note the portable set has no linux-arm64 target (NIT-6).

## Findings

### [MAJOR] SKILL.md Core Rule contradicts the new renderer-first protocol-context flow

`skill/skills/parley-deck/SKILL.md:12` (unchanged by this release) says: "Always read `parley-deck/COOPERATION.md` first." The new Required Protocol Context section (`skill/skills/parley-deck/SKILL.md:22-50`) says an official launch obtains context from `parley protocol packet` with an attestation, and that reading the raw live file directly is the disclosed `full-fallback` used only when no renderer is reachable (and the updated Startup Flow step 1, SKILL.md:103 in the diff, points at that section). An agent obeying line 12 literally reads the raw file first, never visits the renderer, and produces a launch with no attestation — exactly the unattested-launch state the bundled COOPERATION.md §9 checklist now calls "unresolved". This is the release's own focus area (manual skill launch instructions), and the fix is trivial: reword line 12 to defer to Required Protocol Context. Should fix before tagging.

### [MINOR] Manual-launch instructions don't classify a reachable renderer that exits non-zero without refusing

SKILL.md:44-47 defines "renderer unreachable" (example: older CLI without the command) and `refused` (a stop). But `Render` also returns plain errors that are neither: authority failures (e.g. consumer deck with unreadable/missing lock, source.go:81-110) and publication I/O failures — including the deliberate fail-closed no-direct-write path on filesystems that can't hard-link (source.go:312-315). At the CLI these print `protocol packet: <error>` with exit 1 and no attestation (protocol_packet.go:76-79), distinguishable from `REFUSED — …` (protocol_packet.go:80-86), but the skill never tells the facilitator which class maps to `full-fallback` versus stop. The text's own example (an older CLI erroring out) implies non-refusal failures are treated as unavailability, which lands on the safe live-authority fallback — but that mapping is implicit. One sentence in SKILL.md would close it.

### [MINOR] Bundled portable snapshot now presents a chosen transport where it had a placeholder

`skill/skills/parley-deck/references/COOPERATION.md:5` changed from `` **Transport:** `<transport-choice>` (pick one of local-dir | github-pr | gitlab-mr at deck bootstrap — see §0) `` to `` **Transport:** `github-pr` `` (skill-release-diff.patch:136-139). The byte-identical sync with the CLI default (both blob c12e523) is a defensible goal, and §0 line 54 still instructs choosing/replacing the line; but the portable reference is the file agents without repository context read, and it now shows a specific transport as if already selected while `**Workspace:**` on line 3 remains a placeholder. In the CLI template the value is a substitution sentinel (workspace.go:111); in the static skill snapshot nothing substitutes it. Consider keeping the placeholder in the snapshot (or a sync that restores it).

### [MINOR] Applicability-rule conditions combine as first-match OR, undocumented for map authors

`Rule.applies` (cli/internal/protocolpacket/applicability.go:160-177) returns on the first matching dimension — phases, then tracks, then transports, then flags. A rule listing both `phases: [6]` and `tracks: [deliberation]` includes the block for ANY phase-6 launch on any track, not the intersection. Over-inclusion is the safe direction and the default stays full context, but a map author writing intersection intent gets silently different packet content, and `packet check` proves structure only (stated honestly at protocol_packet.go:174). Worth one line of schema documentation or an intersection semantic.

### [NIT] `headerTransport` scans all lines, fence-blind

cli/internal/protocolpacket/source.go:143-152 takes the first `**Transport:**`-prefixed line anywhere, including inside fenced code blocks. The header precedes everything in practice, so this is theoretical; fence-aware parsing (as `Parse` already does) would remove the edge.

### [NIT] Published packet bodies accumulate with no housekeeping

Every non-refused render writes a content-addressed body under `.parley-runtime/protocol-packets/` (source.go:187-192); names key on mode/phase/track plus body digest, so each distinct protocol revision and launch shape adds a new ~100 KB file with no GC. Harmless for a long while; worth a prune rule eventually.

### [NIT] Windows skips directory-permission enforcement by design

checkPrivateDir returns early on windows (source.go:250-254) because the mode bits are synthetic; the 0700/0600 hardening is effectively POSIX-only. Disclosed in-comment; acceptable, but release notes don't mention that packet-body privacy rests on default ACLs on Windows.

### [NIT] CHANGELOG version gaps consistent with unreleased bases, but unverifiable here

CLI CHANGELOG jumps 1.46.0 → 1.48.0 (no 1.47.0 entry) and skill CHANGELOG 2.10.0 → 2.12.0 (no 2.11.0), while the diff bases were VERSION 1.47.0 / package 2.11.0 — coherent if those were unreleased development bumps. The candidate tree is not a git checkout, so whether 1.47.0/2.11.0 ever shipped is unknown history I cannot resolve and do not assume.

### [NIT] Portable build target set lacks linux-arm64

`skill/scripts/build-portable.js:21-27` builds windows x64/arm64, linux x64, macos x64/arm64 — no linux-arm64, although the CLI evidence set includes a linux-arm64 Go build log. If ARM Linux consumers matter, add `node24-linux-arm64`.

## Verified release claims (PRIMARY, code- and doc-located)

- **Full context remains the default:** packet.go:11-13, packet.go:369-395 (default path emits `ModeFull`, body = raw source, packet sha = source sha, plus a shadow audit record); `--optimize` defaults false at protocol_packet.go:56; matches usage text (protocol_packet.go:24-25), both CHANGELOGs, both release notes, and SKILL.md:34-35. The live attestation for this review confirms it end-to-end.
- **Optimized packets are not claimed proven:** SKILL.md:34-35 ("not a default or a proven efficiency improvement"), cli CHANGELOG "shadow diagnostics and optimized packets are experimental", release notes make "no measured performance claim". No doc in the inspected set claims trial completion or measured efficiency.
- **Refusal behavior:** empty source and detected secrets refuse with no body and no packet hash (packet.go:244-254); refusal is printed as `REFUSED`, exits 1, and still emits the attestation JSON when `--json` (protocol_packet.go:80-86); `WriteBody` refuses refused contexts (source.go:169-175); publication is atomic hard-link with no direct-write fallback, symlink-checked, idempotent-verify (source.go:275-357); consumer decks must verify core+lock+overlay before anything else and drift blocks optimization (source.go:80-120). Bundled COOPERATION.md §9 (cli diff, line 98 of cli-release-diff.patch) matches this behavior in prose.
- **Honest limitations preserved:** both release notes explicitly do not claim the empirical audit, its amendment, the 12-task pilot, or the packet experiment as complete; missing historical worktrees remain unknown; the R3 timing-window NIT is disclosed rather than hidden.
- **Packaging:** package.json `files` covers the add-on payloads; `prepack` enforces the sha256 manifest (package.json:66); recorded pack/portable/test evidence in evidence/ is consistent with the pinned hashes (recorded by release tooling — SECONDARY; I did not execute builds, tests, or packs, per my constraints).

## Open questions

- Was CLI 1.47.0 / skill 2.11.0 ever published anywhere? If yes, the missing CHANGELOG entries become a real docs defect rather than a NIT (see NIT above; unknown history, not resolvable from the candidate tree).
- Is the intended long-term contract that `references/COOPERATION.md` stays byte-identical to `internal/protocol/defaults/COOPERATION.md`? If yes, the placeholder question in MINOR-3 should be settled once, in the sync script, rather than per-release.

## Inspected scope

- evidence/review-context.json and its body_path packet (full-phase6-deliberation-4519258c….md, read at §Quickstart, §0, §1, Phase 6 + dispositions, §15.2-15.3).
- evidence/cli-release-diff.patch (all 101 lines), evidence/skill-release-diff.patch (all 195 lines), evidence/cli-release-notes.md, evidence/skill-release-notes.md.
- CLI: internal/protocolpacket/packet.go, source.go, applicability.go (full read); internal/app/protocol_packet.go (full read); internal/protocol/defaults/COOPERATION.md (header + §9 checklist + §0 line 54); targeted search across internal/ for Transport handling (workspace.go:111,138 located).
- Skill: skills/parley-deck/SKILL.md (lines 1-60 + startup-flow diff hunk), references/COOPERATION.md (header), parley-addon.json (via diff), package.json, scripts/build-portable.js (full read).
- Recorded evidence consulted as SECONDARY artifacts only: npm-pack.json, portable-build.log.
- Not inspected (out of bounded scope / not executable under my constraints): the driver/telemetry/budget code behind the CHANGELOG "Added/Fixed" bullets other than the protocol-packet surface, the Go and Node test suites (tests are being run separately; I claim no execution), git history (candidate is not a checkout).

## Recommendation

**GO — approve the release** of parley-deck-cli 1.48.0 and parley-deck-skill 2.12.0 as authorized (direct merge/release, no development PR). Fix MAJOR-1 (the SKILL.md:12 Core Rule wording) before tagging if the release train allows — it is a two-line documentation edit with no code risk; the MINOR/NIT items are safe as fast-follows. I found no CRITICAL issue, no evidence-integrity violation, and no overclaim in the shipped documents; absent experiment results and unknown history remain recorded as such.
