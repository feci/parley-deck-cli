# claude-1 — independent execution of the concurrency witness and selected close-gate cases (2026-09-16)

Author: claude-1
Date: 2026-09-16
Kind: implementation note — executed observations only

## What this note is, and what it is not

This is a **bounded evidence-preparation note**, written for later canonical review by
whoever holds the Phase 6 / close-gate role on this idea. It is not a review artifact.

It contains:

- the exact command I executed, once, and the raw artifacts it produced;
- identity of the immutable source snapshot the command ran against;
- observed process/test counts, the fixture's serial-vs-barrier output line, timing,
  log paths, and content hashes;
- my own critique of what the execution does and does not establish.

It deliberately contains **no verification verdict**. Under §15.1 a verdict is an
assignment of `CONFIRMED` / `WRONG` / `UNVERIFIED` (or equivalent truth-status language)
to a claim; raw command output reported without such a classification is *evidence, not a
verdict*. Everything below is offered as evidence. I assign no status to any claim about
the implementation, and I do not sign off, concur, or dispute anything here.

I did **not** author the concurrency fixture. The `kimi-1` strings inside it
(`evidence.RunCriterion(..., "kimi-1")`, `implementer: "kimi-1"`) are fixture
attribution labels baked into the test data; they are not my identity and I make no
claim about their authorship beyond "I did not write them".

## Scope of the execution

In scope — exactly two things, both pre-existing tests I did not modify:

1. `TestSerialVsBarrierConcurrencyFixture` in `internal/app/driver_evidence_test.go` —
   the two-process witness: a barrier criterion whose halves *must* run concurrently.
2. `TestEvidenceCloseGate.*` in the same package — the selected close-gate cases.

Out of scope, and **not** supported by anything in this note:

- current-tree evidence covering every original acceptance criterion;
- independent semantic proof that an arbitrary criterion command's output means what it
  claims (I observed a test binary's own exit statuses and its own printed line — that is
  not a semantic proof of arbitrary command output);
- whole-branch acceptance;
- AC-E2 closure;
- whole-audit completion;
- any deployment, experiment, freeze, or budget action.

Constraints I operated under: no source edits, no Git mutations, no new provider calls,
no subagents, no migrations, no retries, single execution, preserve failed output.

## Pre-execution reading (PRIMARY, source snapshot)

Read before executing, from the immutable snapshot at `/private/tmp/pd-check-b1bsdntx/source`:

- `internal/evidence/execute.go:118` —
  `cmd := exec.CommandContext(ctx, "sh", "-c", command)`.
  So `RunCriterion` crosses a **real OS process boundary**: `sh -c` per criterion. The
  fixture's criterion command re-execs the *test binary itself*
  (`os.Executable()` + `-test.run '^TestHelperProcess$'`, gated by
  `PARLEY_BARRIER_HELPER=1`), so each half is a separate child process tree, not an
  in-process goroutine.
- `internal/app/driver_evidence_test.go:269` — `TestSerialVsBarrierConcurrencyFixture`
  runs the pair twice through the real `evidence.RunCriterion`:
  - **serial**: `dial` first (polls for a port file for 2s, then `os.Exit(1)`), then
    `serve` under a 500ms context (no client ever arrives → deadline kill). The test
    asserts *both* halves are `StatusFail`.
  - **barrier**: both launched concurrently under a 10s context; the test asserts each
    record is `StatusPass` **and** `Command.Format == FormatGoTestJSON` **and**
    `Command.ExecutedCases == 1`.
  - the two halves synchronise over a real loopback TCP socket (`net.Listen` on
    `127.0.0.1:0`, port published via atomic `WriteFile` + `Rename`, client `net.Dial`s
    it). That is genuine inter-process rendezvous, not a simulated one.
- `internal/app/driver_evidence_test.go:314` — `TestHelperProcess` is inert without
  `PARLEY_BARRIER_HELPER=1`, so it contributes a trivial pass when run as an ordinary test.

I read this to know *what* I was executing. It is source reading, not a claim that the
design is correct.

## Execution record

Executed **once**. No rerun, no retry, no edit to source or fixture.

### Command I invoked

```
python3 /private/tmp/pd-independent-w358qi00/r/execute-checks.py
```

Wrapper exit code: `0`.

### Command the wrapper executed (from `work/started.json`, not from coordinator text)

```
go test -count=1 -timeout 8m -p 1 -json ./internal/app \
  -run '^(TestSerialVsBarrierConcurrencyFixture|TestEvidenceCloseGate.*)$'
```

- cwd: `/private/tmp/pd-check-b1bsdntx/source`
- started: `2026-09-16T05:18:28.360472+00:00`
- wrapper pid `9955`, `go test` pid `9956`
- env pinned by the wrapper: `GOPROXY=off`, `GOFLAGS=-mod=readonly`,
  `GOTOOLCHAIN=local`, `TMPDIR`/`GOTMPDIR` redirected into the scratch `work/tmp`
- `-count=1` (uncached), `-p 1` (one package at a time), no `-race`, no full suite —
  the disk headroom was ~362 MiB, which the run did not exhaust
- toolchain present on this machine: `go version go1.27.1 darwin/arm64`

### Source identity

- manifest: `/private/tmp/pd-independent-w358qi00/r/source-manifest.json`
- **470** entries — I counted them myself: 456 `.go`, 9 `.md`, `go.mod`, `go.sum`,
  1 `.json`, 1 `.yaml`, 1 extensionless (`VERSION`); 112 entries under `internal/app/`
- manifest `sha256` = `98af53bea7d1c3e320bf3bb302e44b7e1a8f18a6779d3364f2c50f90eb15251f`
  (recomputed by me; equals the value the wrapper reported)
- the wrapper verified every one of the 470 files' `sha256` **before** the run (it asserts
  on mismatch) and **again after** it: `source_mismatches_after: []`. So the tested tree
  was byte-identical before and after execution — the run mutated no tested source.
- files I read, hash-checked by hand against the manifest, both matching:
  - `internal/app/driver_evidence_test.go` →
    `b973100606b9334eee05f2c59c3ff70d982b678c2fba695c7c3ae21387c5f3a4`
  - `internal/evidence/execute.go` →
    `c31f9b559b67af953c9d2c2a024b664aac883d7d6ea3ea0ee50fab4c8ed757d9`
- `/private/tmp/pd-check-b1bsdntx/source/.git` does **not** exist (confirmed).

### Observed results — parsed by me from the raw JSON log, not copied from the summary

Raw stream: 65 JSON events — 33 `output`, 16 `pass`, 15 `run`, 1 `start`; **zero** `fail`,
`skip`, or `pause`/`cont` events.

15 top-level tests ran, **all pass**, no subtests:

| Test | Result | Elapsed |
|---|---|---|
| `TestEvidenceCloseGatePositive` | pass | 0.28s |
| `TestEvidenceCloseGateNoReport` | pass | 0.04s |
| `TestEvidenceCloseGateSelfVerdict` | pass | 0.21s |
| `TestEvidenceCloseGateStaleTree` | pass | 0.22s |
| `TestEvidenceCloseGateAllSkip` | pass | 0.20s |
| `TestEvidenceCloseGateZeroExecution` | pass | 0.18s |
| `TestEvidenceCloseGatePartialScope` | pass | 0.18s |
| `TestEvidenceCloseGateShellNotCertified` | pass | 0.19s |
| `TestEvidenceCloseGateCorruptReport` | pass | 0.20s |
| **`TestSerialVsBarrierConcurrencyFixture`** | **pass** | **2.68s** |
| `TestEvidenceCloseGateCompletionTransition` | pass | 0.39s |
| `TestEvidenceCloseGateStatusFlipWithoutTransitionDenied` | pass | 0.35s |
| `TestEvidenceCloseGateCompletionTransitionWrongVerifier` | pass | 0.26s |
| `TestEvidenceCloseGateCompletionTransitionExtraEditDenied` | pass | 0.34s |
| `TestEvidenceCloseGateCompletionTransitionUnapplied` | pass | 0.25s |

The 16th `pass` is the package event: `ok parley-deck-cli/internal/app 6.395s`.
Wrapper wall-clock including build: **8.34s**. `stderr` was **empty** (0 bytes) — no build
diagnostics, no vet output, no panic text.

`TestHelperProcess` is **not** in the 15: the `-run` regex excludes it, so it executed only
as the fixture's spawned child, never as a test of its own. That is the expected shape.

### Serial / barrier output (the witness line)

Exactly one such line appeared in the stream, attributed to
`TestSerialVsBarrierConcurrencyFixture`:

```
fixture: serial=(client:fail,server:fail) barrier=(pass,pass) — serial execution cannot certify a barrier criterion
```

Because the fixture prints this line *after* its assertions, and the test passed, the
following assertions held during this run (they are in the fixture source at
`driver_evidence_test.go:269`, which I read):

- serial: `recC.Status == StatusFail` **and** `recS.Status == StatusFail`;
- barrier: for both records, `Status == StatusPass`, `Command.Format ==
  FormatGoTestJSON`, and `Command.ExecutedCases == 1` — i.e. one executed case per
  concurrent half, as described.

### Process boundary — what I can and cannot say

The boundary is real: `evidence.RunCriterion` → `RunCriterionControlled` →
`exec.CommandContext(ctx, "sh", "-c", command)` (`execute.go:118`), with
`procctl.SetNewProcessGroup(cmd)` and a `cmd.Cancel` that kills the whole group. The
fixture passes `control == nil`, so it takes the plain `sh -c` path, **not** the
`criterionSupervisor` gated-start path at `execute.go:122`. The criterion command is the
test binary re-execing itself under `PARLEY_BARRIER_HELPER=1`, so each half is a child
process, and the two halves rendezvous over a real loopback TCP socket (`net.Listen` on
`127.0.0.1:0` → atomic port-file publish → `net.Dial`).

Structurally this run made **4** `RunCriterion` calls (2 serial, 2 barrier), each spawning
its own `sh` in its own process group, each running one helper child.

**I did not directly count live OS processes.** I did not sample `ps` during the run and I
am not permitted to rerun, so the process count above is read off the source and the call
sites, not off an observed process table. The observed facts are the test outcomes and the
printed line; treat the process arithmetic as source-derived, not measured.

### Artifacts (preserved, not modified)

| Path | sha256 |
|---|---|
| `/private/tmp/pd-independent-w358qi00/r/work/go-test.stdout` (12 667 B) | `a61727d6f8f858c07ad95d9e9f420b5731b9587f0ffa3ef6894a494aaf51f80a` |
| `/private/tmp/pd-independent-w358qi00/r/work/go-test.stderr` (0 B) | `e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855` |
| `/private/tmp/pd-independent-w358qi00/r/work/result.json` | wrapper summary |
| `/private/tmp/pd-independent-w358qi00/r/work/started.json` | pre-run command/pid record |

I recomputed both log hashes myself; they match the wrapper's. The `stderr` digest is the
SHA-256 of the empty input, consistent with the observed 0-byte file.
`work/tmp` is empty after the run — the scratch git repos and `t.TempDir()` trees the
gate tests create were cleaned up, and nothing leaked outside the provided scratch dir.

## My critique — what this execution does not settle

Offered as critique, not as verdicts on anyone's claims.

1. **The snapshot is not bound to the published commit by anything I can check.** The
   coordinator states the canonical integration HEAD is
   `d86b48b3e4b18b0f4704bc1b1691644145b7c8c3`. The tested tree has no `.git`, so I have no
   object store, no commit tree, and no way to bind these 470 files to that SHA. What I
   verified is self-consistency: the tree matches the manifest, before and after. Anyone
   relying on "this is HEAD" is relying on the coordinator's assertion, not on my run.
   The 470-file manifest also covers only `.go`/config/`.md` — a divergence in any file
   outside it would be invisible here.
2. **Re-running an author's tests is not independent definition of correctness.** All 15
   tests, including the close-gate cases, are the implementation's own. Executing them
   independently rules out "the suite was never run" and "it only passes from cache"; it
   does not rule out a gate whose *assertions encode the wrong rule*. Reading whether the
   close-gate semantics are the right semantics is review work I was scoped out of, and it
   remains open.
3. **The barrier proof is self-reported.** `ExecutedCases == 1` comes from parsing
   test2json that the helper itself printed (`pass(...)` in `TestHelperProcess`). That
   exercises the executor's parsing and record-keeping faithfully, but it is a fixture
   emitting a well-formed envelope about itself. It is not evidence that arbitrary
   real-world criterion commands emit envelopes whose contents can be trusted.
4. **Timing sensitivity is asymmetric — and in the safe direction.** The serial failure is
   structural (the counterpart process does not exist at all: the client polls 2s with no
   server, the server is deadline-killed at 500ms with no client), so it cannot pass by
   luck. The barrier half *is* timing-dependent (2s port-file poll, 10s outer deadline);
   under heavy load it could flake to FAIL. So the witness fails closed, not open — but a
   single green run on an idle machine is weak evidence about its behaviour under load,
   and I ran it exactly once.
5. **One run, one selection, one machine.** `-count=1` with no `-race` and no repetition on
   `darwin/arm64`, go1.27.1. Nothing here speaks to data races, to other packages, or to
   any criterion outside these 15 tests.

## Limits of this note (explicit)

This note supports **only** independent execution of this one concurrency witness and the
selected `TestEvidenceCloseGate*` cases, on the stated immutable snapshot.

It does **not** support, and must not be cited as: current-tree evidence covering every
original acceptance criterion; independent semantic proof of arbitrary command output;
whole-branch acceptance; AC-E2 closure; or whole-audit completion. It is not a Phase 6
review, a signoff, a consensus contribution, an amendment round, or a verification verdict,
and it carries no confidence rating. No deployment, experiment, or freeze was performed or
is implied.

No source was edited, no git state was changed, no provider was called, no subagent was
launched, and no history was rewritten. Writes were confined to this file and the provided
scratch `work/` directory.
