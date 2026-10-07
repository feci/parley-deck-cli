# Cycle-5 frozen host verification — 2026-10-08

CLI 2705a1e / skill 99b3f3f. All 13 host checks exit0. Full Go 805.951 seconds.
Build, binary build, vet, Windows cross-build, native boot/lease, race, shared and
local regression tests (including new runcontrol cycle5 cases), shared runner/driver,
and adversarial recognizer probes pass. 100 changed Go files are gofmt-clean.
Product input/output hashes match; both roster files are unchanged. No provider
invocation is recorded by the guard. Exact commands/durations/paths and hashes are
in source-context/fixup-cycle-5-20261008; raw logs stay under
.parley-runtime/quota-implementation/fixup-5/host.

Phase0/5/8 packets rendered with this current-tree binary: all full, source=packet
73613f95aabff8ed57372a535913985554014b71911cc8126ec3f7890a9d4f4e, no fallback.
Deck/skill protocol bytes equal; embedded normalization guard passes in the full suite.
Skill full suite:399 Node/54 Python/6 manifests, 224.359 seconds, exit0.

Windows cross-build is not runtime validation: hosted Windows tests fail as disclosed
in codex-1-cycle5-platform-ci.md. No Windows durability fix is claimed. AC2 remains
NOT MET / owner-waived. The first independent round08 failed provider502 before an
artifact; the unchanged retry is queued. This producer report grants no acceptance.
