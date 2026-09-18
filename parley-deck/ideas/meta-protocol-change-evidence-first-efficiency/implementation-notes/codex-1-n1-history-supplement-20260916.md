# N1 and history supplement — 2026-09-16

This is a coordinator checkpoint, not a whole-audit acceptance or operator grant.
The source published in PR73 still has 444 Go/module files at recovery commit1041259.
All later candidates below remain separate until validation completes.

Claude's independent review e8ab39c1 found no MAJOR/CRITICAL N1 issue and conditionally
accepts its source subject to broad checks. It found two MINOR items: signoff precheck
must mirror a blocked-consensus abort before rendering; preflight comments should
explicitly describe new-reservation exhaustion, not all cached-session refusals.
Kimi owns app/driver_precheck.go, its test and comment-only step_session.go /
cycle_session.go corrections in invocation kimi-1-n1-blocked-signoff-mirror-20260916.
No other agent owns those paths. No behavior widening, refund or atomicity is proposed.

Combined456-Go/module full tests FAILED1090.049s with2612 pass events and one failing
behavior (parent/subtest events): existing driver migration test expected a fourth step
to be spent on an already-exhausted cycle. That expectation conflicts with the intended
N1 preflight correction. The coordinator changes only this expected total to3, preserving
imported+actual charges, original epoch, and exact two real nested child starts.
The failed run remains failed. Every other package passed there; race/vet were not run
by that stopped command.

The history test supplement originally failed an exact error-string expectation;
strict decoding correctly refused earlier. Corrected checks passed31 events4.511s.
Kimi's subsequent eligibility helper preserves all rows and strict grants but withholds
unusable declarations for malformed/incomplete history, explaining uncertainty. Its
first composed tests failed because a116-byte fixture string was mistaken for its
117-byte newline-terminated file, and an older test asserted the obsolete false advice.
Coordinator fixes record the serialized newline, retain unknown history, require the
withheld declaration/explanation, and still assert the original strict scanner refusal
for the exact would-be value. No scanner/grant behavior changed in this supplement.

The assembled458-Go/module candidate passed compile12.251s, focused83events18.437s
and complete budget+driver700events60.669s. Its six-package race and vet are pending.
Source is immutable while checks run. Additional app minor correction will receive
explicit later validation; these results never retroactively cover changed bytes.

Zcode's four launcher tests first compiled, but the imported-ledger fixture passed an
absolute idea path to an API requiring a canonical relative path. Coordinator corrected
only the fixture path. Then all13fake test events passed1.874s, race13passed13.528s,
and vet passed0.367s. Restoring only the older main.go compiled1.544s and failed all
four new behavioral tests1.245s. Old timeout, unused-import and fixture failures remain.
A further explicit launcher test/race/vet check uses the assembled458-Go/module source.
The hidden runtime package is not included by ordinary go test ./.... No real round ran.

The successful real history preview remains read-only: visible lower bound1, two unknown
registered roots, one unknown-scope file set in25copies,322retained sources, and three
ungrouped pre-start refusals. Lower bound1 is not selected operator total N. Declaration,
count, epoch, migration apply and actual grouped round remain unapproved/unexecuted.
A concrete guarded preview and proposal are being prepared, preserving max3 and900s.

Accounting after these three stopped calls:88normalized terminals plus1missing terminal;
56terminal costs plus1missing-terminal cost unknown; USD120.320013 known CLI estimates,
total unknown. Implementation/source review is separate from prospective USD15 total
experiment spending. No experiment freeze/treatment or delivery follow-up clock started.

## Final scoped correction and concrete preview

Kimi invocation4a1913a2 completed214.665s and supplied the blocked-consensus gate,
a regression test and comment-only preflight disclosure. Its own note is retained
verbatim. On the composed final458-Go/module source, all-package compile passed11.787s,
23app focused checks36.49s,23app race checks52.493s and vet0.899s. Restoring only the
prior app/driver_precheck.go compiled2.626s and failed the new blocked-signoff test2.967s.
The separate six-package race still runs on the pre-final-app-delta458Go source;
its results do not silently cover the four later source/comment paths.

Explicit local launcher checks on the assembled458Go source passed13events2.626s,
race13events9.955s and vet0.307s. The actual candidate product CLI read-only preview
passed5.868s:325retained sources, same visible floor1 and unknown coverage. Three new
requested.json files from stopped source-review/fix invocations explain the increase
from322; no prior source changed. The new digest is31950020aa1c253ed7bd68c7c42049b1cf1f496f437380d0776b8d712e8cae34.
The real launcher preview correctly refuses missing cycle authority, while retaining
its concrete input/config digest. A user decision on the explicit accounting proposal
has been requested; no answer, count selection or migration is inferred from time.

Updated inventory:89normalized stopped terminals plus1missing terminal;57terminal costs
plus1missing-terminal cost unknown; known CLI estimates remain USD120.320013. No real
provider CLI is active. The source candidates remain unintegrated until broad checks finish.

## Broad supplement completed

The immutable458Go pre-final-app-delta source completed all six race packages:
2082test/subtest events passed1388.771s, followed by vet0.636s. Source bytes were
unchanged. Together with the separately validated final app/comment delta above,
these checks complete the source integration gate. They do not relabel the original
failed full command, grant historical accounting authority, or accept the whole audit.
