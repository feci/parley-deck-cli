---
agent: claude-1
idea: windows-portability
round: 1
date: 2026-09-24
---

## Protocol attestation

```json
{"context_mode": "full", "source_sha256": "8ce83cde0b4f1c3019229869433df0c32451fdbaa13a186952fdf9805c3a9db7", "packet_sha256": "8ce83cde0b4f1c3019229869433df0c32451fdbaa13a186952fdf9805c3a9db7"}
```

Packet hash re-computed locally before reading: `shasum -a 256
.parley-runtime/protocol-packets/full-phase1-deliberation-8ce83cde…9db7.md` →
`8ce83cde0b4f1c3019229869433df0c32451fdbaa13a186952fdf9805c3a9db7`, matching both attested
values. **Round-1 independence (Phase 1): I did not read `kimi-1.md` or `zcode-1.md`, and I do
not know their contents.** Accuracy note: both peer files did land on disk while I was writing
(mtimes `zcode-1.md` 23:26:06, `kimi-1.md` 23:30:19, mine 23:30:24 local), so the discipline
here rests on my not having opened them — every source I read is cited inline — and not on
their absence. I state this rather than the stronger, and false, "they did not exist yet".
Work tree `worktrees/windows-portability`, HEAD `993a663`, base `868825f`. I changed no
product code, no workflow, no branch, tag or release, and triggered no hosted run.

### Ownership declaration (§15.1) — read this before weighing my verdicts

`source-context/release-ci-revalidation-claude-1.md` is **my own** prior artifact. Under §15.1
I am an **owner** of every claim that first appears canonically there, so I **may not issue a
verification verdict on any of them**. Below I therefore do two distinct things, and label which:

- For **claude-1-owned** claims (Findings A–F of that report): I present current-HEAD evidence
  as *evidence only*, with no truth-status classification, and flag it
  **NEEDS-NON-OWNER-VERDICT** — kimi-1 or zcode-1 must verdict these. Where my own prior
  statement is wrong I append a `SELF-CORRECTION` (a weakening takes effect immediately).
- For **zcode-1-owned** claims (`release-repair-plan-zcode-1.md`, W1–W10 / U1–U3): I issue
  verdicts with provenance tags.
- For **new** claims I raise here I am the owner, so I supply the locator and the executed
  command and **do not** self-verdict. Peers should treat them as unverified until verdicted.

Source copies re-verified byte-for-byte against `source-context/README.md`:
`release-ci-revalidation-claude-1.md` = `754230021eb2aa3f…9853`,
`release-repair-plan-zcode-1.md` = `82878046d8bd515c…0dd3`. Both match. `PRIMARY`.

---

## Summary

The Windows problem is almost entirely a **runtime semantics** problem, not a build problem:
I cross-compiled every package and the whole module builds clean for `windows/amd64` **and**
`windows/arm64`, with exactly **one** compile failure in the tree — a single test line
(`internal/evidence/tree_report_test.go:97`, `syscall.Mkfifo`). Everything else that fails on
Windows fails because POSIX semantics were assumed at runtime: permission-bit privacy checks
that Windows can never satisfy, a process-identity probe that is a stub, mandatory file-sharing
instead of advisory locking, `>` in a filename, and `sh -c` in product source. The good news is
that the repository already contains a real per-OS porting layer (eight `*_windows.go` files,
including a working `LockFileEx` implementation) and already depends **directly** on
`golang.org/x/sys v0.36.0`, which ships every Win32 API needed to implement genuine ACL-based
privacy verification and genuine process attribution — so the honest fixes need **zero new
dependencies**, not a weakened check. The two defects I would not fix — `sh -c` in the evidence
path and cross-restart durable kill — should become explicit, tested, user-visible refusals.

---

## Proposed approach

### 0. The reframe that should drive the plan

Both source reports treat Windows as "13 red packages". The more actionable decomposition,
from evidence I executed at HEAD:

| Surface | Size | Character |
|---|---|---|
| Compile | **1 line** (`tree_report_test.go:97`) | trivial |
| Product runtime semantics | **~6 classes, ~8 product files** | real design |
| Test portability | **≥33 test files, 8 packages** | mechanical sweep |
| Honest refusals | **2 surfaces** | design + docs + tests |

The product work is small and bounded. The test sweep is the bulk of the *volume* and almost
none of the *risk* — provided it never hides a product defect (see §5).

**Evidence I executed myself (`PRIMARY`; new claim, I am the owner, needs a peer verdict):**

```
$ GOOS=windows GOARCH=amd64 go build ./...   → exit 0
$ GOOS=windows GOARCH=arm64 go build ./...   → exit 0
$ for p in $(go list ./...); do GOOS=windows GOARCH=amd64 go vet "$p"; done
### parley-deck-cli/internal/evidence
vet: internal/evidence/tree_report_test.go:97:20: undefined: syscall.Mkfifo
=== DONE ===        # every other package type-checks clean, test files included
```

Limit, stated plainly: cross-compilation and `vet` prove **compilation**, not behavior. They
say nothing about what happens when the binary runs. Nobody on this run has Windows hardware;
every behavioral claim below must be settled on the hosted `windows-latest` leg.

### 1. W-PRIV — snapshot privacy: implement real ACLs, do not relax the check

At HEAD, `internal/trajectory/snapshot.go:162` (directories) and `:544` (regular files) gate on
`info.Mode().Perm()&0077 != 0`. Go's Windows `FileMode` is synthesized from a single attribute:

> `$GOROOT/src/os/types_windows.go:178-183` —
> `if fs.FileAttributes&syscall.FILE_ATTRIBUTE_READONLY != 0 { m |= 0444 } else { m |= 0666 }`,
> then `m |= ModeDir | 0111` for directories.

> `$GOROOT/src/syscall/syscall_windows.go:759-773` — `Chmod` **only** toggles
> `FILE_ATTRIBUTE_READONLY` from the `S_IWRITE` (0200) bit and discards every other bit.

So a Windows directory reports `0777` or `0555` and a file `0666` or `0444`; `&0077` is
`0077`/`0055`/`0066`/`0044` — never zero, whatever you chmod. *(Evidence for a claude-1-owned
claim — Finding F.1. **NEEDS-NON-OWNER-VERDICT**. `PRIMARY`: Go stdlib at the locators above,
read at `go1.27.1`; the repo pins `go 1.26` in `go.mod` and CI reported `1.26.8` — a peer
should confirm the same lines in the CI toolchain, though both the current and the
`modePreGo1_23` path at `types_windows.go:239-241` carry identical logic.)*

**Design — `privateSnapshotDirectory` and the file check become per-OS, with Windows getting a
stronger check, not a waived one.** All APIs below exist in the already-direct
`golang.org/x/sys v0.36.0` (verified by symbol grep in the module cache, locators in
*Existing alternatives*):

1. On create: build a **protected, non-inherited DACL** with exactly one ACE — the current
   token user (`windows.GetCurrentProcessToken().GetTokenUser()`) granted full control — via
   `windows.ACLFromEntries([]EXPLICIT_ACCESS, nil)`, and apply it with
   `SetNamedSecurityInfo(path, SE_FILE_OBJECT, OWNER_SECURITY_INFORMATION|DACL_SECURITY_INFORMATION|PROTECTED_DACL_SECURITY_INFORMATION, …)`.
   `PROTECTED_DACL` is the load-bearing flag: it blocks inheritance from the parent, which is
   the actual Windows analogue of "0700, and nothing else got in".
2. On verify: `GetNamedSecurityInfo` → `sd.Owner()` must equal the token user; `sd.DACL()` →
   iterate `AceCount` with `GetAce`, and refuse unless **every** ACE is
   `ACCESS_ALLOWED_ACE_TYPE` for that same owner SID. Any additional trustee, any denied-type
   ACE we did not place, any inherited ACE → the same user-visible refusal string the POSIX
   path returns.
3. Keep `!info.IsDir()` and the symlink/reparse rejection on both platforms; Windows reparse
   points surface as `ModeSymlink` under Go's default `winsymlink` behavior
   (`types_windows.go:165-176`).

This is **stricter** than the POSIX check in one respect (it pins the owner identity, which
`0700` does not) and equivalent in the respect that matters (no other principal has access). I
would state that asymmetry in the code comment rather than claim parity.

**What I will not do:** the "smallest fix" both reports floated — make the perm check
Unix-only — is a privacy regression on the platform that most needs the check, and the prompt
forbids it. If the ACL work is judged too large for this idea, the correct fallback is a
**refusal**, not a waiver: snapshot capture on Windows returns a reviewed, documented,
user-visible error. That is worse product but honest product.

### 2. W-PROC — process identity: the stub is a live correctness defect, and it is fixable

This is the class **neither source report identified**, and I think it is the most important
one for a user actually running parley on Windows. At `internal/procctl/procctl_windows.go`:

```go
func (windowsProbe) alive(pid int) bool             { return pid > 0 }   // :19
func (windowsProbe) procStart(pid int) (string, bool) { return "", false } // :20
func (windowsProbe) command(pid int) (string, bool)   { return "", false } // :21
func (windowsProbe) pgid(pid int) (int, bool)         { return 0, false } // :22
func (windowsProbe) supportsDurableKill() bool       { return false }    // :17
```

`Alive` therefore reports **true for every positive PID forever**, including after the process
exits and across reboots. Tracing the callers at HEAD:

- `internal/runner/durablekill.go:29` — `KillAgentDurable`: `Alive` is always true, so the
  `clearStale` branch is unreachable; `Attributed` then always fails closed
  (`procctl.go:110`, "durable kill is unsupported on this platform"), so the function returns
  `"process verification failed (…); not killed"`. **On Windows, `parley`'s durable-kill path
  can neither kill nor clear a stale badge — a run stuck "running" stays stuck, permanently,
  with no recovery through this path.**
- `internal/runner/durablekill.go:63` — `AgentLiveness` returns `"stale"` for a genuinely live
  agent, so the TUI badge is wrong in both directions.
- `internal/trajectory/verification.go:142` — `if !procctl.Alive(…) { continue }` is never
  taken, so every recorded ordinal reaches `KillTreeAttributed`, which always refuses, and the
  function returns `"captured criterion cleanup refused: durable kill is unsupported on this
  platform"`. **Captured-verification cleanup fails unconditionally on Windows**, even when
  nothing is running.
- `internal/runner/reviewsnapshot.go:244` — `ownerAlive` always true (and `bootID()` returns
  `""` on Windows, so the boot check at `:241` is skipped), so `sweepStaleSnapshots` never
  reclaims a crashed run's snapshot.

`internal/driver/proclive_windows.go` makes the same conservative choice and documents its
reason explicitly: *"There is no portable signal-0 liveness probe on Windows without
golang.org/x/sys"*. **That premise is false at HEAD** — `golang.org/x/sys v0.36.0` is a
*direct* require in `go.mod:10`, and `internal/budget/lock_windows.go:7` and
`internal/fsutil/replace_windows.go:7` already import `golang.org/x/sys/windows`.

**Design:**

- `alive(pid)`: `windows.OpenProcess(PROCESS_QUERY_LIMITED_INFORMATION, false, pid)` →
  `GetExitCodeProcess`; alive iff the code is `STILL_ACTIVE` (259). Fail **closed** on
  `ERROR_ACCESS_DENIED` (a PID we cannot inspect is not a PID we may signal).
- `procStart(pid)`: `GetProcessTimes` → `CreationTime` as a 100 ns `FILETIME`. This is the
  Windows analogue of Linux `starttime`, and it is *finer*-grained than the microsecond
  comparison `Attributed` already performs.
- `command(pid)`: `QueryFullProcessImageName` for the image path (cheap, no PEB read). Full
  argv via `NtQueryInformationProcess` + PEB is available but is undocumented-API territory;
  I would start with the image path and keep `commandMatches` (`procctl.go:157`) unchanged.
- `pgid(pid)` → **Job Objects**, not process groups. `SetNewProcessGroup` creates a job with
  `JOBOBJECT_EXTENDED_LIMIT_INFORMATION.BasicLimitInformation.LimitFlags |=
  JOB_OBJECT_LIMIT_KILL_ON_JOB_CLOSE`, assigns the child (`AssignProcessToJobObject`), and
  `KillGroup` calls `TerminateJobObject` instead of shelling out to `taskkill /T /F`
  (`procctl_windows.go:31`). A job both kills the tree reliably *and* is the honest identity
  the `pgid` facet is asking for.
- `bootID()`: derive from boot instant (`GetTickCount64` against wall clock, quantized), or —
  safer — leave `""` and keep `supportsDurableKill() == false`.

**The scope decision I want peers to argue with.** There are two coherent stopping points:

- **(P-A) Liveness only.** Implement `alive` truthfully; leave `supportsDurableKill() == false`.
  This alone fixes the stuck-badge bug (`clearStale` becomes reachable), fixes the TUI, fixes
  snapshot sweeping, and fixes `verification.go:142`. It is ~40 lines in one file and does not
  touch the attribution contract at all.
- **(P-B) Full attribution + Job Objects.** Also enables durable cross-restart kill on Windows.
  Larger, touches `SetNewProcessGroup`/`KillGroup`, and needs its own adversarial tests for
  PID reuse.

**I recommend P-A for this idea and P-B as a separate follow-up.** P-A removes real user-visible
breakage with a small, reviewable change; P-B changes a safety gate that the base commit's
Linux repair just spent an entire review cycle hardening, and it deserves the same scrutiny
rather than being bundled into a portability sweep.

### 3. W-LOCK — it is *not* missing `LockFileEx`; it is missing `FILE_SHARE_DELETE`

**Verdict on zcode-1's W3** ("The budget locking model (`budget/lock.go` et al.) assumes POSIX
advisory semantics. Needs a LockFileEx-style design — architecture-scale"): **WRONG at HEAD**,
in its stated cause and therefore in its size estimate. `PRIMARY` — `internal/budget/lock_windows.go`
at HEAD already implements `tryLock`/`unlock` with
`windows.LockFileEx(…, LOCKFILE_EXCLUSIVE_LOCK|LOCKFILE_FAIL_IMMEDIATELY, …)` (`:16-22`) and
`publishExclusive` with `windows.MoveFileEx(…, MOVEFILE_WRITE_THROUGH)` (`:28-38`), and
`internal/budget/lock.go:84-86` already normalizes lock keys for case-insensitive filesystems
with a comment naming that hazard. The Windows lock design exists and was written deliberately.

The actual mechanism of the *"The process cannot access the file because it is being used by
another process"* failures is different and much smaller:

> `$GOROOT/src/syscall/syscall_windows.go:395` —
> `sharemode := uint32(FILE_SHARE_READ | FILE_SHARE_WRITE)`

`os.Open`/`os.OpenFile` open **without `FILE_SHARE_DELETE`**, so while any handle is open the
file cannot be renamed-over or deleted — `MoveFileEx` fails with a sharing violation. By
contrast `os.Root` *does* pass it:

> `$GOROOT/src/os/root_windows.go:176` —
> `syscall.FILE_SHARE_READ|syscall.FILE_SHARE_WRITE|syscall.FILE_SHARE_DELETE`
> (same in `$GOROOT/src/internal/syscall/windows/at_windows.go:151,223,252,268,333,374,448`)

**Design:** add one `fsutil` open helper, Windows-specific, that opens with
`FILE_SHARE_READ|FILE_SHARE_WRITE|FILE_SHARE_DELETE` (either via `windows.CreateFile` directly
or by routing through `os.Root`, which `internal/trajectory/snapshot.go:194` already uses), and
use it on the publication/ledger read paths. A second, separate trap to handle in the same
change: `os.OpenFile(path, …, 0400)` sets `FILE_ATTRIBUTE_READONLY`
(`syscall_windows.go:401-403`), and a read-only file cannot be renamed over — so any staging
file written with a mode lacking `0200` needs its attribute cleared before replace.

This turns an "architecture-scale" item into a focused change in `internal/fsutil` plus its
call sites. **I would like a peer to attack this claim specifically** — it is the single
biggest scope reduction I am proposing, and if I am wrong the estimate moves a long way.

### 4. W-NAME — gate filenames, and an unvalidated ID that is a path surface on *every* OS

**Verdict on zcode-1's W2** (`gate.go:40` `EdgeID = from+"->"+to`, used in the filename at
`:105`; `>` invalid on Windows): **CONFIRMED**, locators exact. `PRIMARY` — `internal/pipeline/gate.go:40`
`func EdgeID(from, to string) string { return from + "->" + to }` and `:104-106`
`filepath.Join(PipelineDir(deckDir, slug), "gates", edgeID+".gate.json")`.

Two facts that make the fix easier than zcode-1's "compatibility decision" framing suggests,
and one that makes it larger:

- **Easier:** nothing ever parses an edge ID back out of a filename. Every read goes through
  `GatePath(deckDir, slug, edgeID)` with an edge ID the caller already holds
  (`gate.go:110`, `:131`; producers at `executor.go:105,187,255`, `app/pipeline_cmd.go:1118`;
  `dag.go:31` consults it as an opaque key), and the gate JSON carries `Edge`/`ID` in its body
  (`gate.go:72-74`). So the on-disk **name** can change without touching identity. `PRIMARY`,
  verified by exhaustive grep for `.gate.json` / `GatePath` / `EdgeID(` across `internal` and
  `cmd`, with no `TrimSuffix`/`ReadDir`-based reverse parse anywhere.
- **Easier:** I therefore argue against a **per-OS** name (which both source reports assume).
  A deck is routinely shared through git or a network share between a Windows and a macOS
  user; divergent on-disk names would make the same gate invisible from the other machine.
  Use **one universal encoding on every platform**, plus a read-fallback to the legacy
  `a->b.gate.json` name for one release so existing decks keep resolving.
- **Larger:** block IDs are **not validated at all**. `internal/pipeline/manifest.go:156-173`
  checks only non-empty and uniqueness — no character set, no length, no path check — and that
  unvalidated, user-authored ID is interpolated into a path twice: `gate.go:105` and
  `executor.go:30` (`filepath.Join(deckDir, "ideas", slug+"__"+blockID)`). `PRIMARY`; grep for
  `Clean(`/`ContainsAny`/`sanitiz` across `internal/pipeline` (non-test) returns nothing. On
  any OS that is a traversal surface (`blockID = "../../x"`); on Windows it additionally admits
  reserved device names and `:` (alternate data streams). **This is a security-adjacent defect
  that Windows merely exposed, and I think it belongs in this idea** — the filename fix has to
  decide an encoding anyway. *(New claim, I am the owner — needs a peer verdict, and I would
  value a challenge on whether it is in scope or should be split out.)*

**Design:** validate block IDs at manifest load (explicit allowed set, e.g. `[A-Za-z0-9._-]`,
bounded length, no leading/trailing dot, not a Windows reserved device name) **and** encode the
filename defensively. On encoding: `net/url.PathEscape` is the stdlib candidate and covers
`< > " \ | ? * /` and space, but I verified it does **not** escape `:` and does not touch
reserved device names or a trailing dot:

```
$ go run …  // net/url.PathEscape
"block-a->block-b" -> "block-a-%3Eblock-b"
":"  -> ":"        "CON" -> "CON"        "end." -> "end."
```

So `PathEscape` alone is insufficient — I would use it as the base and add an explicit `:`
escape plus a device-name guard, with validation as the primary defense and encoding as
defense in depth. Naming validation as *the* fix and encoding as belt-and-braces is the honest
ordering; an encoder alone would still let a hostile manifest pick a nasty name.

### 5. W-SHELL — `sh -c` is in **product** source, and this is where I must correct myself

**`SELF-CORRECTION` (§15.1) — claude-1, Finding F.5 of `release-ci-revalidation-claude-1.md`.**

> Prior statement (the claim I am replacing): *"10 fixtures hard-code `/bin/sh` as the agent
> command (`launch_test.go:48`, `:133`). Test-side only: the sole `/bin/sh` in product source is
> `internal/procctl/terminal_unix.go:40`, which is build-tagged Unix."*
>
> Corrected statement: **product source is not POSIX-free.** At HEAD, non-test product code
> invokes a POSIX shell in three places outside any build tag:
> - `internal/evidence/execute.go:118` — `exec.CommandContext(ctx, "sh", "-c", command)`
> - `internal/evidence/execute.go:122` — `exec.CommandContext(ctx, "sh", "-c", criterionSupervisor)`
> - `internal/app/driver_impl.go:266` — `run(checks, exec.CommandContext(ctx, "sh", "-c", checks))`
>
> This is a **weakening** of my prior claim and under §15.1 it takes effect immediately; it
> does not require a peer verdict to land. The prior claim was produced by grepping for
> `/bin/sh` with a leading slash, which misses the bare `"sh", "-c"` form. `PRIMARY` —
> `grep -rn '"/bin/sh"\|"sh", "-c"' --include='*.go' internal cmd | grep -v _test.go` at HEAD.

This matters far more than a miscount. `internal/evidence/execute.go:100-114` embeds a whole
POSIX shell program — `criterionSupervisor` — using `trap ':' TERM`, `IFS= read -r`, `$!`,
`wait`, `kill -0`, `</dev/null`. That is the supervisor for **controlled criterion capture**:
the LE-4 verification path, the thing the whole evidence regime rests on. And
`driver_impl.go:266` is the `checks:` gate itself.

**The honesty trap I most want the deck to notice.** GitHub's `windows-latest` image ships Git
for Windows, whose `sh.exe` is on `PATH`. If that is so, these calls may well *succeed on the
hosted runner* — and a green Windows leg would then be read as "shell execution works on
Windows" when what it actually shows is "works on a machine that happens to have Git for
Windows installed." A user on a clean Windows box has no `sh`. *(That `sh.exe` is on the
runner's PATH is `RECALL` — I have no Windows host and did not verify it; it is precisely the
kind of thing the hosted leg must be made to answer explicitly rather than incidentally.)*

**Design — make the runner answer the question instead of hiding it.** Add a Windows test that
asserts the *declared* behavior of the criterion path when `exec.LookPath("sh")` fails, and
have the product **detect the absence explicitly and refuse with a named, user-visible error**
("captured verification requires a POSIX shell; install Git for Windows or set …") rather than
surfacing a bare `exec: "sh": executable file not found in %PATH%`. Then the hosted leg proves
the refusal path, not the accident.

I am **not** proposing to port `criterionSupervisor` to PowerShell or `cmd.exe` in this idea.
Its cancellation semantics were hardened across several reviewed cycles; re-implementing them
in a second shell dialect, unreviewed, on a platform nobody here can run, is exactly the kind
of change that looks like progress and ships a regression. See *Existing alternatives* for
`mvdan.cc/sh/v3`, which is the one option I think could genuinely resolve this later.

### 6. W-CRLF — a content-integrity defect, not a test nit

`internal/runner/hardening_test.go:440` writes `"v2-dirty\n"`; `:455` reads it back out of the
snapshot and got `"v2-dirty\r\n"` on Windows. The path between them is git:
`internal/runner/reviewsnapshot.go:123-127` builds `snapGit` with only
`GIT_OPTIONAL_LOCKS=0` and `GIT_INDEX_FILE`, then runs `read-tree` / `add -A` / `write-tree`
(`:128-134`), `commit-tree` (`:138-140`), and `checkout --detach` (`:150`). With Git for
Windows' customary `core.autocrlf=true`, `add -A` normalizes CRLF→LF into the object store and
`checkout` writes LF→CRLF back out. **The snapshot's bytes are therefore not the working
tree's bytes** — in a product whose entire premise is byte-exact evidence capture. There is no
`.gitattributes` in the repo and no `autocrlf` handling anywhere in Go or workflow source
(`PRIMARY`: `grep -rn "autocrlf\|gitattributes"` over `*.go`/`*.yml` returns nothing outside
`graphify-out/`; `ls -a` shows no `.gitattributes`).

**Design:** pin `-c core.autocrlf=false -c core.eol=lf` on parley's **own** git invocations, so
capture is byte-exact regardless of the user's global config. Precedent already exists two
lines away — `reviewsnapshot.go:138-139` pins `-c user.name=parley -c user.email=parley@localhost`
per command for exactly this "do not depend on ambient config" reason. This is product-side and
small. `hardening_test.go:455` is the test that **detects** the bug; skipping it on Windows
would be the textbook case of converting a product failure into a skipped test.

### 7. W-TEST — the sweep, with one non-negotiable rule

Grounded lower bound, computed at HEAD: **33 test files across 8 packages** carry at least one
POSIX-shaped construct (`"sh"/"bin/sh"` invocation, a `#!/bin/sh` fixture, `Chmod`, or
`t.Setenv("HOME")`), distributed `internal/app` 13, `internal/trajectory` 9, `internal/driver` 4,
`internal/runner` 3, and one each in `protocolpacket`, `procctl`, `evidence`, `agents`. That is
a lower bound, not the total: the hosted log reported ~98 unique failing test names.

Sub-classes and their fixes:

- **Mkfifo** (`tree_report_test.go:97`) — the only compile break. The file already carries a
  fallback (`makeUnsupportedSocket`) for hosts that cannot make FIFOs, so the shape exists;
  split the FIFO branch behind a build tag and give Windows a named-pipe or AF_UNIX-socket
  equivalent. The test's *intent* — an unsupported entry type must fail the digest, never be
  silently skipped — is fully expressible on Windows and must be kept.
- **`#!/bin/sh` fixtures** (9 files; `app_test.go` alone has 9 occurrences, e.g.
  `writeFakeParleyDeckSkill` at `:1742` writing an extension-less script) — write a `.bat`/`.cmd`
  on Windows. `PRIMARY`, locators confirmed at HEAD.
- **The `internal/app` abort.** `app_test.go:155` `payload["parley_deck_skill"].(map[string]any)`
  and `:156` `skill["project"].(map[string]any)` are unchecked type assertions that panic on
  nil and kill the whole package binary; the same pattern recurs at `:129` and `:177`, and
  `grep -c '\.(map\[string\]any)'` counts **8** in that file. Use the comma-ok form so one
  missing fixture fails one test instead of aborting 27 others — which is why `wait`/`usage`
  have no Windows results at all. *(Locators confirmed; the causal chain to the missing
  `wait`/`usage` results is a claude-1-owned claim from Finding F.3 — **NEEDS-NON-OWNER-VERDICT**.)*
- **Chmod fixtures.** **Verdict on zcode-1's W7 locators: CONFIRMED.** `PRIMARY` — I initially
  read them as chmod sites and they are not; they are the **assertion** sites a CI log prints,
  and they are exact: `internal/driver/impl_test.go:807` is the
  `t.Fatalf("an unreadable round must escalate…")` following `os.Chmod(bad, 0o000)` at `:800`,
  and `internal/driver/phase_event_test.go:156/170/186` are the three `t.Fatal` lines after
  `os.Chmod(…, 0)` at `:152/:165/:182`. **Verdict on zcode-1's inclusion of
  `strict_gate_test.go:179` in this class: WRONG.** `PRIMARY` — `grep -c "Chmod"
  internal/driver/strict_gate_test.go` = `0`; `:170-180` writes a *regular file where a
  directory is expected*, which should behave the same on Windows. Whatever fails there has a
  different cause and must be diagnosed separately rather than swept into the chmod bucket.
  Note also the pre-existing escape hatch at `impl_test.go:800-801`
  (`t.Skipf("cannot make a directory unreadable here")`) — on Windows `Chmod` *succeeds*
  (it just sets READONLY), so the skip does not fire and the assertion fails instead.
- **HOME.** `t.Setenv("HOME", …)` in `internal/agents/configmodel_test.go:23` and
  `internal/app/roster_configread_test.go:19`; `grep -rn USERPROFILE` over `internal` and `cmd`
  returns **nothing**. The *product* side is fine — `os.UserHomeDir()` is used at
  `sessionstore.go:43`, `config/runtime.go:585`, `runner.go:1298,1337`,
  `agents/configmodel.go:30` and reads `USERPROFILE` on Windows — so this is genuinely
  test-only. **Verdict on zcode-1's W8: CONFIRMED** (`PRIMARY`, locators above).
  *(I checked one adjacent suspicion and it is a false alarm worth recording so nobody re-files
  it: `internal/acp/shellenv.go:54` does `"HOME="+os.Getenv("HOME")`, but `:44` returns early on
  `runtime.GOOS == "windows"`, so it is unreachable there.)*

**The rule for the whole sweep:** a Windows exclusion is admissible only when the *behavior
under test cannot exist* on Windows — not when it is inconvenient to set up. "Chmod cannot make
a directory unreadable" is **not** such a case once §1 gives us ACLs: the test should assert
the ACL-denied path instead. Each exclusion gets its own reviewed justification naming the
platform mechanism that makes it inapplicable, per the prompt. I would put that list in
`IMPLEMENTATION.md` as an explicit table, one row per exclusion, so review can attack them
individually rather than as a batch.

### 8. Sequencing

1. **Compile + harness** — the one `Mkfifo` tag, plus the `internal/app` comma-ok assertions so
   the package stops aborting. Effect: the Windows leg starts producing real per-test data for
   `wait`/`usage` instead of one panic. This is the highest information-per-line change in the
   whole idea and should land first.
2. **Product semantics** — W-PRIV (ACLs), W-PROC P-A (liveness), W-LOCK (share-delete),
   W-NAME (validation + encoding), W-CRLF (git config pinning).
3. **Refusals** — W-SHELL named refusal + durable-kill refusal, each with a test.
4. **Test sweep** — the remaining fixture classes, with the exclusion table.
5. **Hosted validation** — unfiltered `go build ./...` and `go test ./... -count=1 -timeout 45m`
   on all three legs, no new `continue-on-error`, no `-run`, no `-short`, no `|| true`.

Steps 1 and 2 are separable per file and could be claimed by different implementers; §12 of the
prompt requires one implementation owner at a time, so I would expect serialization by the
organizer rather than parallel claims.

---

## Existing alternatives

Per §15.6(a): for each mechanism this proposal would build by hand, what the platform,
toolchain or an existing dependency already ships, with a locator. Sources consulted: the Go
standard library at `$GOROOT` (`go1.27.1`, `/opt/homebrew/Cellar/go/1.27.1/libexec`), the
module cache copy of `golang.org/x/sys@v0.36.0`
(`/Users/tomasfecko/go/pkg/mod/golang.org/x/sys@v0.36.0/windows`), this repository at HEAD
`993a663`, and `go.mod`. Symbol availability was checked by grep against the module cache.
I did **not** consult upstream issue trackers or package registries — availability and
maintenance status of the third-party options below are `RECALL` and must be verified before
any of them is adopted.

**Already in this repository — reuse, do not reinvent.** The single most important alternative
is the repo's own porting layer. Eight files already exist:
`internal/acp/sysproc_windows.go`, `internal/app/budget_attended_windows.go`,
`internal/budget/lock_windows.go`, `internal/driver/proclive_windows.go`,
`internal/fsutil/replace_windows.go`, `internal/procctl/procctl_windows.go`,
`internal/procctl/terminal_windows.go`, `internal/protocolcore/nofollow_windows.go`. Every
design above should extend that convention rather than introduce a new abstraction.

| Mechanism I would build | What already ships | Locator | Why still needed / decisive reason |
|---|---|---|---|
| ACL privacy check (§1) | `golang.org/x/sys/windows`: `GetNamedSecurityInfo`, `SetNamedSecurityInfo`, `ACLFromEntries`, `EXPLICIT_ACCESS`, `Trustee`, `GetCurrentProcessToken`, `Token.GetTokenUser`, `SECURITY_DESCRIPTOR.DACL()/Owner()`, `GetAce`, `ACCESS_ALLOWED_ACE`, `SE_FILE_OBJECT`, `PROTECTED_DACL_SECURITY_INFORMATION` | `x/sys@v0.36.0/windows/security_windows.go:1209,1241,1418,1481,1483; :658,667,709; :893,1098,1274`; `zsyscall_windows.go:808` | **Adopt.** Complete round-trip available; `go.mod:10` already requires x/sys **directly**. Zero new dependencies. |
| ACL helpers | `github.com/hectane/go-acl` (`Chmod`-style POSIX-mode→ACL shim) | external module, not in `go.mod` | **Reject.** New dependency for what x/sys already covers, and its POSIX-mode emulation is the wrong abstraction — it would re-encode the `0700` model we are trying to stop pretending in. |
| Windows process liveness (§2) | `x/sys/windows`: `OpenProcess`, `GetExitCodeProcess`, `PROCESS_QUERY_LIMITED_INFORMATION`, `GetProcessTimes`, `QueryFullProcessImageName`, `NtQueryInformationProcess` | all present in `x/sys@v0.36.0/windows` (symbol grep; `STILL_ACTIVE` absent, it is the constant 259) | **Adopt.** Directly falsifies the stated premise of `internal/driver/proclive_windows.go` ("no portable probe … without golang.org/x/sys"). |
| Process-tree kill | (a) `taskkill /T /F` — already used, `procctl_windows.go:31`; (b) Job Objects: `CreateJobObject`, `AssignProcessToJobObject`, `SetInformationJobObject`, `TerminateJobObject`, `JOB_OBJECT_LIMIT_KILL_ON_JOB_CLOSE`, `JOBOBJECT_EXTENDED_LIMIT_INFORMATION` | `x/sys@v0.36.0/windows` (all present) | **Adopt (b) if P-B.** `taskkill` is a shell-out with a PID-reuse race between lookup and kill; a Job Object is atomic and is also the honest `pgid` analogue. Keep (a) until P-B lands. |
| Process introspection | `github.com/shirou/gopsutil`, `github.com/mitchellh/go-ps` | external modules, not in `go.mod` | **Reject.** Large dependency surface for four syscalls x/sys already exposes. |
| Cross-process file locking (§3) | Already implemented: `windows.LockFileEx`/`UnlockFileEx` | `internal/budget/lock_windows.go:16-26` | **Already done.** This is why zcode-1's W3 sizing is wrong. |
| Portable lock library | `github.com/gofrs/flock` | external module, not in `go.mod` | **Reject.** Would replace a working, deliberately-designed implementation (note the beyond-EOF offset rationale at `lock_windows.go:11-14`) with a generic one. |
| Delete-shareable file handles (§3) | `os.Root` already passes `FILE_SHARE_DELETE`; already used at `internal/trajectory/snapshot.go:194` (`os.OpenRoot`) | `$GOROOT/src/os/root_windows.go:176`; `$GOROOT/src/internal/syscall/windows/at_windows.go:151,223,252,268,333,374,448` | **Adopt.** Prefer routing through `os.Root` over a bespoke `CreateFile` wrapper where the call site allows it — stdlib, already in use here. |
| Atomic replace | Already implemented: `windows.MoveFileEx(… REPLACE_EXISTING\|WRITE_THROUGH)`; stdlib `os.Rename` also uses `windows.Rename` + `fixLongPath` | `internal/fsutil/replace_windows.go:16-27`; `$GOROOT/src/os/file_windows.go:259-265` | **Already done.** The repo's version exists for `WRITE_THROUGH`, which `os.Rename` does not give. |
| Filename-safe edge IDs (§4) | `net/url.PathEscape` (stdlib) | `net/url`; verified by execution — escapes `< > " \ \| ? *` and space, **does not** escape `:`, `CON`, trailing `.` | **Adopt as base only.** Insufficient alone; must be paired with manifest-level ID validation and a device-name guard. Stated as a partial, not a solution. |
| Long paths on Windows | Go prefixes `\\?\` automatically; CI sets `git config --global core.longpaths true` | `$GOROOT/src/os/path_windows.go:100-131`; `.github/workflows/tests.yml` "Enable long paths (Windows)" | **Partial — and a live honesty gap.** Go's own file ops are covered, but `core.longpaths` is a *global git* setting the **CI applies and the product does not**. Users hit git-side long-path failures the hosted leg cannot see. Belongs in the refusal/docs work. |
| POSIX shell on Windows (§5) | (a) Git for Windows `sh.exe` — incidental, not a declared dependency; (b) `mvdan.cc/sh/v3` — pure-Go POSIX shell interpreter; (c) `cmd.exe /c` / `powershell -Command` | (a) `RECALL`, unverified; (b) external module, not in `go.mod`; (c) OS built-ins | **Defer, do not adopt now.** (a) is what a green hosted leg would silently rely on — the core honesty risk in §5. (b) is the only option that could make `sh -c` genuinely portable with no external install, and is the right subject for a **follow-up** idea: it is a new dependency executing untrusted `checks:` strings, so it needs its own security review. (c) cannot express `criterionSupervisor`'s `trap`/`wait`/`kill -0` semantics without a rewrite. |
| CRLF neutrality (§6) | git's own `core.autocrlf`, `core.eol`, `.gitattributes`; per-command `-c` pinning precedent | `internal/runner/reviewsnapshot.go:138-139` (pins `user.name`/`user.email` this exact way) | **Adopt the existing precedent.** Prefer per-command `-c` over a repo `.gitattributes`, because parley snapshots *arbitrary user repos* whose attributes we do not control. |
| FIFO test fixture (§7) | Windows named pipes (`windows.CreateNamedPipe`); AF_UNIX on Windows 10+ via `net.Listen("unix")`; the file's own existing `makeUnsupportedSocket` fallback | `x/sys@v0.36.0/windows`; `internal/evidence/tree_report_test.go:98-103` | **Adopt the existing fallback shape.** The file already handles "this host cannot make FIFOs"; Windows is one more such host, not a new concept. |
| Per-OS code split | Go build constraints — `_windows.go` suffix (implicit) and `//go:build` | used correctly throughout; note `_unix.go` is **not** an implicit GOOS suffix, which is why `internal/fsutil/replace_unix.go:1` needs its explicit `//go:build !windows` while `replace_windows.go` needs none | **Already the convention.** Extend it; introduce no new platform-abstraction layer. |

**Scoped null:** I found no existing toolchain mechanism for (i) validating user-authored
pipeline block IDs against filesystem-reserved names — this must be hand-written; and (ii) a
Windows analogue of `bootID` that is as cheap and exact as Linux's `/proc/sys/kernel/random/boot_id`
— which is part of why I recommend P-A over P-B.

---

## Concerns / open questions

1. **The §15.1 ownership constraint materially shapes this idea, and the organizer should plan
   for it.** Because `source-context/release-ci-revalidation-claude-1.md` is my own artifact,
   I cannot verdict Findings A–F, which includes the *single largest* defect class (W1/F.1,
   ~69 hosted failures). Every claim I marked **NEEDS-NON-OWNER-VERDICT** above requires
   kimi-1 or zcode-1 to verdict it independently before it can support an acceptance criterion.
   This is the rule working as designed, not a defect — but it means the deck needs at least
   one peer to re-derive the snapshot-privacy class from primary sources, and the organizer
   should confirm that happened rather than assume it. **This is not a protocol-change request;
   I am flagging an operational consequence, per the prompt's instruction to surface protocol
   issues before implementation.**
2. **Does a green hosted `windows-latest` leg license removing the experimental label?** I do
   not think it does, on its own, and the release gate in `00-prompt.md` should be read
   carefully. Three named gaps survive a green leg: (a) `windows-latest` is **x64 only** — the
   published `windows-arm64.exe` would still have zero executed evidence, and my clean
   `GOOS=windows GOARCH=arm64 go build` proves compilation only; (b) the leg runs with
   `core.longpaths true`, which is **not** the Windows default; (c) if `sh.exe` is present on
   the image, the evidence path's shell dependency is satisfied by accident. I would keep the
   label unless the leg is green **and** these three are separately disposed of in writing.
3. **Is GitHub's `windows-11-arm` runner class available to this repository?** If it is, adding
   an arm64 leg would close gap (a) directly. I could not check availability, pricing or
   eligibility from this session — `RECALL`, needs verification by someone with repo settings
   access before anyone plans around it.
4. **P-A vs P-B for process attribution.** I recommend P-A and want to be argued with. The
   counter-argument I find strongest: leaving `supportsDurableKill() == false` means Windows
   users permanently lack reattach-kill, which may be a worse product outcome than the review
   risk of doing P-B now.
5. **Does the W-NAME block-ID validation belong in this idea?** It is a security-adjacent fix
   that Windows exposed rather than caused. Including it is defensible (the encoding decision
   forces the question); excluding it leaves a known traversal surface open on all platforms
   while we ship a Windows fix. I lean include; I would accept a well-argued split.
6. **`strict_gate_test.go:179` has no established cause.** I refuted its placement in the chmod
   class but produced no replacement mechanism. It needs diagnosis from the hosted log, not a
   guess — and it should not be quietly folded into another class.
7. **What is the declared Windows shell contract?** Before the refusal in §5 can be written,
   someone has to decide the product's position: is Git for Windows a *documented prerequisite*
   for `checks:` and captured verification on Windows, or is a shell-free path required? That
   is a product decision, not a technical one. If the deck cannot settle it, it is a legitimate
   §10 escalation to the owner.
8. **How much of the ~98 unique failing test names is left unexplained?** My classes plus the
   two source reports' account for the large clusters, but I have not reconciled the total
   against a per-test list. Until someone does, "we understand the Windows failures" is an
   overstatement — and the step-1 fix (stopping the `internal/app` abort) will itself change
   the denominator by revealing tests that never ran.

---

## Risks

- **R1 — A green Windows leg gets over-read (highest risk in this idea).** The single most
  likely bad outcome is that step 1 plus the test sweep turns the leg green, and green is
  reported as "Windows works". It would not be: it would mean "Windows works on an x64 GitHub
  runner with Git for Windows installed and `core.longpaths` enabled." *Mitigation:* the
  release claim must enumerate the runner's assumptions, and the experimental label must
  survive until concerns 2(a)–(c) are separately disposed of.
- **R2 — The test sweep launders product defects into skips.** 33+ files of mechanical fixture
  edits is exactly the volume at which a reviewer's attention runs out, and every POSIX fixture
  has a tempting `t.Skip` available — one is already sitting there at `impl_test.go:800-801`.
  *Mitigation:* the exclusion table in `IMPLEMENTATION.md`, one reviewed row per exclusion,
  reviewed individually, plus a diff-level check that no `t.Skip` was added without a row.
- **R3 — The ACL implementation is written by people with no Windows machine.** ACL code is
  easy to write plausibly and wrong: inheritance flags, `PROTECTED_DACL`, and owner-SID
  comparison all have non-obvious failure modes, and a mistake here **weakens a privacy
  control** rather than breaking a build. *Mitigation:* the verification half must be tested
  adversarially on the hosted leg — construct a directory with an extra trustee and assert
  refusal — not merely "create then verify our own creation", which passes even if both halves
  are wrong in the same direction.
- **R4 — P-B silently weakens the attribution gate the base commit just hardened.** The Linux
  U2 repair in this base spent a full review cycle on fail-closed refusal semantics. A Windows
  `Attributed` path that returns `true` too readily would re-open exactly what that work
  closed, on a platform nobody can test locally. *Mitigation:* P-A now; P-B as its own idea.
- **R5 — Scope creep past the release gate.** This idea can absorb W-SHELL portability, P-B,
  the `mvdan.cc/sh/v3` question and arm64 coverage, and then ship nothing. *Mitigation:* the
  three named deferrals above are deferrals, recorded as such, not silently dropped.
- **R6 — W-NAME breaks existing decks.** Changing gate filenames orphans every existing
  `a->b.gate.json`. *Mitigation:* universal (not per-OS) encoding plus a legacy read-fallback
  for one release, and an explicit persisted-name compatibility note in `FINAL.md` as the
  prompt requires.
- **R7 — My own W-LOCK reduction is the load-bearing scope claim and it is unverified
  behaviorally.** I argue W3 shrinks from "architecture-scale" to a focused `fsutil` change
  based on reading Go's share-mode constants — not on observing a Windows process. If the
  hosted failures have an additional cause, the estimate moves substantially. *Mitigation:*
  I am explicitly asking peers to attack this one first.
- **R8 — Hosted flakiness is misattributed to Windows.** The unresolved Linux trajectory/runner
  family (U2 in zcode-1's plan) has never had a root cause established. If similar shapes
  appear on the Windows leg they could be miscounted as portability defects, or real Windows
  defects could be dismissed as that family. *Mitigation:* keep per-run claims separate, as
  `source-context/README.md` instructs for runs `35987916696` and `36009912946`, and never
  classify by resemblance.

---

## Verdict summary

Verdicts I am **permitted** to issue (zcode-1-owned claims), all `PRIMARY` at HEAD `993a663`:

| Claim | Verdict | Basis |
|---|---|---|
| W2 — `gate.go:40` EdgeID `>` in filename at `:105` | **CONFIRMED** | source at both locators |
| W3 — "needs a LockFileEx-style design, architecture-scale" | **WRONG** | `lock_windows.go:16-26` already implements it; real cause is the `FILE_SHARE_DELETE` omission at `syscall_windows.go:395` |
| W7 locators — `impl_test.go:807`, `phase_event_test.go:156/170/186` | **CONFIRMED** | exact assertion sites following the chmod calls |
| W7 — inclusion of `strict_gate_test.go:179` in the chmod class | **WRONG** | that file contains zero `Chmod` calls |
| W8 — `agents/configmodel_test.go` HOME-only | **CONFIRMED** | `:23`; no `USERPROFILE` anywhere in `internal`/`cmd` |
| Provenance `4e39609` (W1), `d387de6` (W2), `a637628` (U1 file), `7f2676f` (W4) | **CONFIRMED** (all four) | `git log -S` / `--diff-filter=A` per locator |

**No `## Verdict conflicts` section is warranted.** The one apparent conflict — U1 provenance,
claude-1 `5712a36` vs zcode-1 `a637628` — **dissolves on inspection rather than resolving**:
the two statements use different predicates and are both true. `git log -- internal/acp/spawn.go`
gives `a637628` (2026-05-24, added) … `5712a36` (2026-09-05, *last touched* before the base
commit's repair) … `a2db799` (2026-09-24). zcode-1 said "introduced"; claude-1 said "last touched
by". `PRIMARY`. Recording this so nobody re-litigates it as a contradiction.

**Also noted, not a defect:** `internal/acp/spawn.go` at HEAD already drains before reaping
(`Stop` at `:130-137` runs `p.wg.Wait()` before `done <- p.cmd.Wait()`; `Wait` at `:163-166`
likewise). U1 is fixed in this base and is **not** part of the Windows scope.
