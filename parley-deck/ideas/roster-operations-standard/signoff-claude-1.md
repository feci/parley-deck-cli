### claude-1

**Verdict:** accept — with one reservation I am recording against a decision I benefit from.

**Role disclosure (§15.5).** I am facilitator, participant and drafter. I also made the VC-2
measurement, and I was in the rebase camp that the user's direction then selected. That is the
configuration §15.1 exists to constrain, and it applies twice here.

**Scope declared.** I read all four round-01 files, all four round-02 files, `00-prompt.md`, and my
own inbox measurement plus its addendum. Ran fresh this session: `parley roster show`,
`parley agents list`, `parley --help`, `/usr/bin/grep` over `~/.parley/agents.toml`, a `find`-based
enumeration of §2 rosters across 47 `COOPERATION.md` files, `codex exec --help` and `kimi --help`
(model flags), and reads of `runtime.go:588-616`, `runner.go:1097-1108`,
`runmanifest/manifest.go:28-55`. I did **not** run a live agent launch to observe the model a
process actually receives, and I issue no verdict on any claim I own.

### Reservation — rebase is being adopted before the thing that makes it safe exists

I measured that `runmanifest.Manifest` records `participants []string` and nothing about model,
effort or adapter. Under rebase, a deck stops carrying its own pins, so **the only place a run's
actual configuration could be recorded is the snapshot — and today the snapshot does not record
it.** Between those two facts there is a window in which neither the deck nor the run state answers
"what did this agent run".

The consensus closes that window because decision 6 is unanimous and ships in the same change. My
reservation is narrow and specific: **implementation must not land the rebase behaviour before the
snapshot captures the effective row.** If they land in that order, every run in between is
unauditable. I would rather this be an explicit ordering constraint in `FINAL.md` than an assumed
one, and I am flagging it precisely because the outcome favours my own earlier position and I do
not want that to make me lenient about its precondition.

### On the user's three directions

**Rebase — accepted, with the ordering constraint above.** The user chose it over additive-pin. I
note for the record that I was in the rebase camp before the user decided, so my agreement here
carries no independent weight; hermes-1 and kimi-1 argued the other side and their signoffs are the
ones that should be read on this point.

**§2 protocol change in this idea — accepted, and the deviation is correctly logged.** §7 asks for a
separate meta idea. The user directed otherwise. What matters is that the *ratification* is not
skipped along with the venue: the protocol text still needs every participant's signoff, and a
signer can still block the wording. The consensus says that. I would block if it did not.

**Mass migration of 40 decks — accepted, and this is the part I am least comfortable with.** The
user authorized the outcome; I imposed four constraints on the method (CLI-executed, backed up,
dry-run-all-first, skip-and-report on anything unclean). I want one more, and I am adding it here
rather than silently: **the dry-run diff goes to the user before a single deck is written**, and
decks belonging to other projects are reported by name in that diff. Seventeen of these decks name
a retired agent and several are months or years old; a migration that "succeeds" on all 40 without
anyone looking at the diff is not a success I can verify.

### VC-1 — `SOURCE` column: I still exclude it

My position changed in round 2 and I hold it. The argument that decided me is codex-1's, not the
count: a single `SOURCE` cell can only name the winning layer for **one** field, so it silently
privileges `MODEL` and misinforms about `EFFORT`, `SPEED` and `AUTO`, whose winning layers can
differ. `--explain AGENT` plus the JSON `sources` object answers the same question without that
defect. **kimi-1 and I were the two who wanted `SOURCE` in round 1**; I changed against my own prior
position and against the only participant who agreed with me, so this is not majority drift.

### VC-3 — `deck|machine`, and `--scope deck` writes the committed file

I adopt `deck|machine` over `local|global`. hermes-1's reason is the right one: `local` is ambiguous
between machine-local and project-local, while `deck` names an actual directory.

**`--scope deck` must write the committed `parley-deck/agents.toml`, never the gitignored
`agents.local.toml`.** A roster change that is invisible to the repository is precisely how a deck
diverges from its own history — which is the failure this whole idea exists to end. If someone wants
a machine-private override they can still edit `agents.local.toml` by hand; the standard verb should
not default to invisibility.

### What I got wrong in this idea

My round-1 file treated the opencode inconsistency as an undefined promotion path. **hermes-1 found
the actual mechanism** — two stores, §2 versus `[roster.*]` — and my contribution was to measure how
far it had spread, not to find it. I also proposed a `SOURCE` column and a copy-style `sync`, and
both were beaten by codex-1's arguments. Four position changes are recorded in
`## Drafter position changes`; all four were forced by another participant.
