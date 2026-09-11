// Package evidence implements typed completion evidence for the
// meta-protocol-change-evidence-first-efficiency slice owned by kimi-1.
//
// A Report binds a reviewed commit plus a tested code-tree digest to typed
// per-criterion records (status, provenance, command/output hashes, executed
// case counts when the format supports them, bounded secret-safe diagnostics).
// Evaluate is the fail-closed closure gate: missing/self/stale/skipped/
// no-execution/package-failure/partial evidence cannot close a whole
// implementation.
//
// Trust boundary: Provenance.Executor/Verifier are ASSERTED runtime identities
// (agent ids from the orchestrator), not cryptographic authentication of a
// human. They exist so a self verdict is detectable and rejectable, not to
// prove who a human operator is.
package evidence

import (
	"fmt"
	"strings"
	"time"
)

// Status is the typed per-criterion outcome.
type Status string

const (
	StatusPass    Status = "pass"
	StatusFail    Status = "fail"
	StatusSkipped Status = "skipped"
	StatusNotRun  Status = "not-run"
)

// Output formats we can attach semantics to. Anything else is unsupported:
// a plain shell exit code plus stdout text is NOT universally semantically
// certified by a regex, so FormatShell evidence cannot prove executed cases.
const (
	FormatGoTestJSON = "gotest-json" // `go test -json` structured events
	FormatEnvelope   = "envelope"    // explicit PARLEY-EVIDENCE {...} JSON lines
	FormatShell      = "shell"       // opaque command output: exit code only
)

// Provenance records who executed and who verified, as asserted runtime
// identity (NOT authentication). Executor==Verifier is a self verdict.
// Verifier is empty until an independent verifier attests via AttestExecution;
// VerifierRerun RETAINS the verifier's actual independent execution record —
// command hash, output hash, counts, exit status and the before/after tested
// tree digests it ran against — so Evaluate can reconcile the attestation
// instead of trusting a supplied name.
type Provenance struct {
	Executor string `json:"executor"`
	Verifier string `json:"verifier"` // empty until an independent verifier attests
	// VerifierRerun is the retained independent execution the attestation is
	// based on. An attestation recorded against a different command, different
	// counts, or a different tree than the report's is not evidence for this
	// report. Callers MUST NOT hand-populate Provenance fields: the only
	// legitimate write path is AttestExecution.
	VerifierRerun *VerifierExecution `json:"verifier_rerun,omitempty"`
}

// VerifierExecution is the retained record of the verifier's independent
// re-execution of one criterion. TreeBeforeSHA256/TreeAfterSHA256 are the
// tested-tree digests the verifier computed immediately before and after its
// re-run; both must exist, be equal (the re-run tested a stable tree), and
// match the report's tested tree.
type VerifierExecution struct {
	Command          CommandEvidence `json:"command"`
	TreeBeforeSHA256 string          `json:"tree_before_sha256"`
	TreeAfterSHA256  string          `json:"tree_after_sha256"`
}

// CommandEvidence binds one executed command to hashes of its exact input and
// bounded scrubbed output, plus executed-case counts when the format supports
// them. ExecutedCases is -1 when the format cannot prove a count.
// FailedPackages is -1 unless the format can prove package-scope outcomes
// (gotest-json): it counts package-level fail/build-fail EVENTS — structured
// proof a package failed at package scope (a failed build runs no test cases).
// One failed build normally emits both a build-fail and a package-level fail,
// so FailedPackages is a failure signal to be read as >0, not a
// distinct-package tally, and it is never folded into the test-case counts.
type CommandEvidence struct {
	Command        string `json:"command"`
	CommandSHA256  string `json:"command_sha256"`
	OutputSHA256   string `json:"output_sha256"` // hash of the RAW captured output
	ExitCode       int    `json:"exit_code"`
	DurationMillis int64  `json:"duration_ms"`
	Format         string `json:"format"`
	ExecutedCases  int    `json:"executed_cases"`  // -1 = unknown/unsupported
	FailedCases    int    `json:"failed_cases"`    // -1 = unknown/unsupported
	SkippedCases   int    `json:"skipped_cases"`   // -1 = unknown/unsupported
	FailedPackages int    `json:"failed_packages"` // -1 = unknown/unsupported (shell/envelope)
	Diagnostics    string `json:"diagnostics"`     // bounded, secret-scrubbed tail
}

// CriterionRecord is the typed record for one named completion criterion.
type CriterionRecord struct {
	Name       string          `json:"name"`
	Status     Status          `json:"status"`
	Command    CommandEvidence `json:"command"`
	Provenance Provenance      `json:"provenance"`
}

// Report is the typed evidence artifact for one close attempt. It is persisted
// as EVIDENCE.json inside the idea directory (an explicit evidence artifact,
// excluded from the tree digest so persisting it does not invalidate itself).
type Report struct {
	Idea           string            `json:"idea"`
	ReviewedCommit string            `json:"reviewed_commit"` // HEAD at review time
	TreeSHA256     string            `json:"tree_sha256"`     // digest of the tested code tree (taken BEFORE execution, re-verified after)
	TreeDirty      bool              `json:"tree_dirty"`      // working tree differs from ReviewedCommit
	GeneratedAt    time.Time         `json:"generated_at"`
	Records        []CriterionRecord `json:"records"`
	// ExtraDigests binds content the tree digest's exclusion set would
	// otherwise hide: path → digest of the file with ONLY the generated
	// evidence content removed (e.g. IMPLEMENTATION.md minus its
	// driver-written ## Validation evidence section). Closure recomputes and
	// compares, so a non-evidence edit to an excluded file invalidates the
	// report instead of hiding behind the exclusion.
	ExtraDigests map[string]string `json:"extra_digests,omitempty"`
	// CompletionTransition, when present, is the single authorized `status:
	// complete` transition for one ExtraDigests-bound path, recorded by the
	// independent verifier via AuthorizeCompletionTransition before the
	// transition is applied. Absent on older reports: those remain valid only
	// in their ORIGINAL state and are never silently upgraded.
	CompletionTransition *CompletionTransition `json:"completion_transition,omitempty"`
}

// ClosureOptions parameterizes Evaluate.
type ClosureOptions struct {
	// RequiredScope is the full original criterion scope; partial evidence
	// (records covering fewer criteria) cannot close.
	RequiredScope []string
	// CurrentTreeSHA256 is the tree digest recomputed at close time. A mismatch
	// with the report's tested digest is stale evidence.
	CurrentTreeSHA256 string
	// Verifier is the independent runtime identity closing the idea. It must be
	// set and must differ from every record's Executor (no self verdict).
	Verifier string
	// RequireStructured requires every record to carry a format that can prove
	// executed cases (gotest-json or envelope). A code-writing contract sets
	// this: scalar/absent shell checks must not bypass it.
	RequireStructured bool
}

// Evaluate is the fail-closed closure gate. It returns a list of reasons; an
// empty list means the evidence may close the implementation.
//
// Fail-closed surface (all must hold):
//   - non-empty required scope, a current tree digest supplied, a reviewed
//     commit and a tested tree digest recorded, tested == current;
//   - no duplicate criterion records; every required criterion has a record;
//   - every record passes, with non-contradictory exit/status, valid counts
//     (structured formats must prove >0 executed, 0 failed, non-negative
//     counts, and a gotest-json record must carry its failed-package count
//     with zero failed package-level events), and command/output hash binding;
//   - every record carries a PERSISTED independent verifier attestation
//     (Provenance.Verifier plus the retained Provenance.VerifierRerun, both set
//     only by AttestExecution after an independent re-execution): the verifier
//     name is non-empty, not the executor, matches the closing verifier
//     identity, and the retained re-run reconciles with the record — same
//     exact command hash, non-empty output hash, pass-consistent exit/status,
//     valid counts on the re-run itself (nothing below -1; a gotest-json
//     re-run must carry its skipped count), equal executed/skipped counts,
//     zero failures, and before/after tested tree digests equal to each other
//     and to this report's tested tree. A
//     caller supplying a different name in ClosureOptions.Verifier conjures
//     nothing — the attestation and its retained execution must already exist
//     in the report.
//
// Trusted-caller boundary: Evaluate and AttestExecution verify the INTERNAL
// CONSISTENCY AND BINDING of the persisted evidence (a supplied name without a
// consistent retained execution, a different command/tree/counts, or a self
// verdict all fail closed). They trust the calling orchestrator to have
// ACTUALLY executed the re-run it attests — a same-UID caller able to
// fabricate mutually consistent hashes is outside this boundary. Runtime
// attribution is not cryptographic authentication of a human (FINAL D3).
func Evaluate(r *Report, opts ClosureOptions) []string {
	var reasons []string
	if r == nil {
		return []string{"no evidence report"}
	}
	if opts.Verifier == "" {
		reasons = append(reasons, "no independent verifier identity")
	}
	if len(opts.RequiredScope) == 0 {
		reasons = append(reasons, "empty required scope — an unbounded criterion set cannot close")
	}
	if r.ReviewedCommit == "" {
		reasons = append(reasons, "no reviewed commit recorded")
	}
	if opts.CurrentTreeSHA256 == "" {
		reasons = append(reasons, "no current tree digest supplied — staleness cannot be ruled out")
	}
	if r.TreeSHA256 == "" {
		reasons = append(reasons, "no tested tree digest recorded")
	} else if opts.CurrentTreeSHA256 != "" && r.TreeSHA256 != opts.CurrentTreeSHA256 {
		reasons = append(reasons, "stale evidence: tested tree digest does not match the current tree")
	}
	byName := map[string]CriterionRecord{}
	for _, rec := range r.Records {
		if _, dup := byName[rec.Name]; dup {
			reasons = append(reasons, fmt.Sprintf("duplicate evidence record for criterion %q", rec.Name))
			continue
		}
		byName[rec.Name] = rec
	}
	// Missing and partial scope: every required criterion must have a record.
	for _, name := range opts.RequiredScope {
		if _, ok := byName[name]; !ok {
			reasons = append(reasons, fmt.Sprintf("missing evidence for required criterion %q (partial scope cannot close)", name))
		}
	}
	for _, rec := range r.Records {
		p := fmt.Sprintf("criterion %q: ", rec.Name)
		switch rec.Status {
		case StatusPass:
		case StatusSkipped:
			reasons = append(reasons, p+"skipped evidence cannot close")
			continue
		case StatusNotRun:
			reasons = append(reasons, p+"no execution recorded")
			continue
		default:
			reasons = append(reasons, p+"status is "+string(rec.Status))
			continue
		}
		// Contradictory exit/status: a pass claim with a non-zero exit is
		// internally inconsistent and never admissible.
		if rec.Command.ExitCode != 0 {
			reasons = append(reasons, p+fmt.Sprintf("contradictory record: status pass with exit code %d", rec.Command.ExitCode))
		}
		// Self verdict: the executor of the check cannot also be the verifier
		// (asserted runtime identity comparison), whether the sameness comes
		// from the closing identity or is persisted in the record itself.
		if opts.Verifier != "" && rec.Provenance.Executor != "" && rec.Provenance.Executor == opts.Verifier {
			reasons = append(reasons, p+"self verdict — executor "+rec.Provenance.Executor+" is the closing verifier")
		}
		if rec.Provenance.Executor == "" {
			reasons = append(reasons, p+"no executor provenance")
		}
		// Persisted independent attestation: an empty Verifier, one equal to
		// the executor, one that does not match the closing identity, or one
		// not bound to THIS tested tree all fail closed.
		switch {
		case rec.Provenance.Verifier == "":
			reasons = append(reasons, p+"no persisted independent verifier attestation (AttestExecution after an independent re-run is required)")
		case rec.Provenance.Verifier == rec.Provenance.Executor:
			reasons = append(reasons, p+"self verdict — persisted verifier equals executor "+rec.Provenance.Executor)
		case opts.Verifier != "" && rec.Provenance.Verifier != opts.Verifier:
			reasons = append(reasons, p+"closing verifier "+opts.Verifier+" does not match persisted attestation by "+rec.Provenance.Verifier)
		}
		if rec.Provenance.Verifier != "" {
			if rec.Provenance.VerifierRerun == nil {
				reasons = append(reasons, p+"attestation retains no independent execution record (AttestExecution is the only legitimate write path)")
			} else {
				reasons = append(reasons, reconcileRerun(p, rec, rec.Provenance.VerifierRerun, r.TreeSHA256)...)
			}
		}
		// Count validity: anything below -1 is malformed in every format.
		if rec.Command.ExecutedCases < -1 || rec.Command.FailedCases < -1 || rec.Command.SkippedCases < -1 || rec.Command.FailedPackages < -1 {
			reasons = append(reasons, p+"invalid (negative) case count")
		}
		// No-execution: a format that can count cases must report at least one,
		// and its counts must be present and consistent.
		if rec.Command.Format == FormatGoTestJSON || rec.Command.Format == FormatEnvelope {
			if rec.Command.ExecutedCases <= 0 {
				reasons = append(reasons, p+fmt.Sprintf("zero executed cases in %s output", rec.Command.Format))
			}
			if rec.Command.FailedCases < 0 {
				reasons = append(reasons, p+"failed-case count absent in structured output")
			}
			if rec.Command.FailedCases > 0 {
				reasons = append(reasons, p+fmt.Sprintf("%d failed cases", rec.Command.FailedCases))
			}
			if rec.Command.Format == FormatGoTestJSON && rec.Command.SkippedCases < 0 {
				reasons = append(reasons, p+"skipped-case count absent in test2json output")
			}
			if rec.Command.Format == FormatGoTestJSON && rec.Command.FailedPackages < 0 {
				reasons = append(reasons, p+"failed-package count absent in test2json output")
			}
			if rec.Command.Format == FormatGoTestJSON && rec.Command.FailedPackages > 0 {
				reasons = append(reasons, p+fmt.Sprintf("%d failed package-level event(s) in test2json output", rec.Command.FailedPackages))
			}
			if rec.Command.FailedCases > rec.Command.ExecutedCases && rec.Command.ExecutedCases >= 0 {
				reasons = append(reasons, p+"conflicting counts: more failed than executed cases")
			}
		} else if opts.RequireStructured {
			reasons = append(reasons, p+"unsupported evidence format "+formatOrNone(rec.Command.Format)+
				" — shell output is not semantically certified; use go test -json or an explicit evidence envelope")
		}
		if rec.Command.CommandSHA256 == "" || rec.Command.OutputSHA256 == "" {
			reasons = append(reasons, p+"missing command/output hash binding")
		}
	}
	return reasons
}

// AttestExecution records an independent verifier attestation for one criterion
// AFTER the verifier has independently executed the check itself. rerun is the
// RETAINED record of that independent execution: its CommandEvidence plus the
// tested-tree digests taken immediately before and after the re-run. The
// attestation is refused unless the re-run is bound to this criterion and this
// report's tested tree: same exact command hash as the criterion record, a
// non-empty output hash, exit 0, a case-proving format, >0 executed and 0
// failed cases, zero failed package-level events, valid counts on the re-run
// itself (a gotest-json re-run with a missing or malformed skipped or
// failed-package count is refused — only the envelope format may legitimately
// carry skipped_cases=-1 and failed_packages=-1), counts equal to the
// original record's, a pass-typed original record, and equal before/after tree
// digests matching the report's tested tree. An opaque shell exit-0 is not an
// attestation basis.
//
// The rerun record is persisted inside the report (Provenance.VerifierRerun),
// not discarded, and Evaluate re-reconciles every one of these bindings at
// close time. This API exists so the only way to place a verifier name on a
// record is to supply evidence of an independent execution; there is
// deliberately no API that stamps a name without one, and a supplied name is
// never itself a verdict.
//
// Trusted-caller boundary: the caller is trusted to have ACTUALLY executed
// the re-run and to supply the real digests it produced; this function
// guarantees that anything LESS than a fully bound, consistent, independent
// execution is refused. It does not and cannot prove execution to a party
// beyond that boundary — runtime attribution is not authentication of a human.
func AttestExecution(r *Report, criterion, verifier string, rerun VerifierExecution) error {
	if r == nil {
		return fmt.Errorf("evidence: cannot attest on a nil report")
	}
	if r.TreeSHA256 == "" {
		return fmt.Errorf("evidence: report has no tested tree digest to bind the attestation to")
	}
	if verifier == "" {
		return fmt.Errorf("evidence: empty verifier identity")
	}
	for i, rec := range r.Records {
		if rec.Name != criterion {
			continue
		}
		if rec.Provenance.Executor == verifier {
			return fmt.Errorf("evidence: %s cannot attest criterion %q it executed (self verdict)", verifier, criterion)
		}
		if rec.Status != StatusPass {
			return fmt.Errorf("evidence: criterion %q is %q — only a pass-typed record can be attested", criterion, rec.Status)
		}
		if errs := reconcileRerun("", rec, &rerun, r.TreeSHA256); len(errs) > 0 {
			return fmt.Errorf("evidence: independent re-run of %q is not an attestation basis: %s", criterion, strings.Join(errs, "; "))
		}
		r.Records[i].Provenance.Verifier = verifier
		r.Records[i].Provenance.VerifierRerun = &rerun
		return nil
	}
	return fmt.Errorf("evidence: no record for criterion %q", criterion)
}

// reconcileRerun verifies that a retained independent execution is bound to
// the criterion record and the report's tested tree: same exact command hash,
// non-empty output hash, exit 0, a case-proving format, >0 executed and 0
// failed cases, zero failed package-level events (a masked exit cannot
// launder a build/package failure into an attestation), valid counts on the
// retained execution itself (nothing below
// -1 in any format; a gotest-json re-run must carry its skipped and
// failed-package counts — only the envelope format may legitimately lack them,
// so format compatibility cannot smuggle an invalid count past
// reconciliation), executed/skipped/failed-package
// counts equal to the original record's, and non-empty before/after tree
// digests that are equal to each other (a stable tree during verification)
// and to the report's tested tree. p prefixes each reason ("" at attestation
// time, the criterion prefix at close time).
func reconcileRerun(p string, rec CriterionRecord, vr *VerifierExecution, treeSHA256 string) []string {
	var reasons []string
	rc := vr.Command
	if rc.ExitCode != 0 {
		reasons = append(reasons, p+fmt.Sprintf("independent re-run exited %d", rc.ExitCode))
	}
	if rc.Format != FormatGoTestJSON && rc.Format != FormatEnvelope {
		reasons = append(reasons, p+"independent re-run used unstructured "+formatOrNone(rc.Format)+" output")
	}
	if rc.ExecutedCases <= 0 {
		reasons = append(reasons, p+"independent re-run executed no cases")
	}
	if rc.FailedCases != 0 {
		reasons = append(reasons, p+fmt.Sprintf("independent re-run recorded %d failed cases", rc.FailedCases))
	}
	if rc.FailedPackages > 0 {
		reasons = append(reasons, p+fmt.Sprintf("independent re-run recorded %d failed package-level event(s)", rc.FailedPackages))
	}
	// Count validity applies to the retained execution exactly as to the
	// original record: anything below -1 is malformed in every format, and a
	// gotest-json re-run must carry real skipped and failed-package counts — a
	// true `go test -json` parse always yields them, so a negative value there
	// means the retained record did not come from one. -1 remains legitimate
	// only for the envelope format, which really lacks those fields.
	if rc.ExecutedCases < -1 || rc.FailedCases < -1 || rc.SkippedCases < -1 || rc.FailedPackages < -1 {
		reasons = append(reasons, p+"independent re-run has an invalid (negative) case count")
	}
	if rc.Format == FormatGoTestJSON && rc.SkippedCases < 0 {
		reasons = append(reasons, p+"independent re-run's test2json output carries no skipped-case count — not a real go test -json parse")
	}
	if rc.Format == FormatGoTestJSON && rc.FailedPackages < 0 {
		reasons = append(reasons, p+"independent re-run's test2json output carries no failed-package count — not a real go test -json parse")
	}
	if rec.Command.CommandSHA256 == "" || rc.CommandSHA256 != rec.Command.CommandSHA256 {
		reasons = append(reasons, p+"independent re-run command hash missing or differs from the criterion record — not the same exact command")
	}
	if rc.OutputSHA256 == "" {
		reasons = append(reasons, p+"independent re-run binds no output hash")
	}
	if vr.TreeBeforeSHA256 == "" || vr.TreeAfterSHA256 == "" {
		reasons = append(reasons, p+"independent re-run is not bound to before/after tested-tree digests")
	} else {
		if vr.TreeBeforeSHA256 != vr.TreeAfterSHA256 {
			reasons = append(reasons, p+"tested tree changed during the independent re-run")
		}
		if treeSHA256 != "" && vr.TreeBeforeSHA256 != treeSHA256 {
			reasons = append(reasons, p+"independent re-run executed against a different tree than the report tested")
		}
	}
	// Result reconciliation: the same exact command on the same tested tree
	// must reproduce the same structured counts; divergence means the two
	// executions are not evidence of the same thing.
	if rc.ExecutedCases >= 0 && rec.Command.ExecutedCases >= 0 && rc.ExecutedCases != rec.Command.ExecutedCases {
		reasons = append(reasons, p+fmt.Sprintf("executed-case count differs between original (%d) and independent re-run (%d)", rec.Command.ExecutedCases, rc.ExecutedCases))
	}
	if rc.SkippedCases >= 0 && rec.Command.SkippedCases >= 0 && rc.SkippedCases != rec.Command.SkippedCases {
		reasons = append(reasons, p+fmt.Sprintf("skipped-case count differs between original (%d) and independent re-run (%d)", rec.Command.SkippedCases, rc.SkippedCases))
	}
	if rc.FailedPackages >= 0 && rec.Command.FailedPackages >= 0 && rc.FailedPackages != rec.Command.FailedPackages {
		reasons = append(reasons, p+fmt.Sprintf("failed-package count differs between original (%d) and independent re-run (%d)", rec.Command.FailedPackages, rc.FailedPackages))
	}
	return reasons
}

func formatOrNone(s string) string {
	if s == "" {
		return "(none)"
	}
	return strings.TrimSpace(s)
}
