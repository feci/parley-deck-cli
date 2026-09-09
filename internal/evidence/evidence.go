// Package evidence implements typed completion evidence for the
// meta-protocol-change-evidence-first-efficiency slice owned by kimi-1.
//
// A Report binds a reviewed commit plus a tested code-tree digest to typed
// per-criterion records (status, provenance, command/output hashes, executed
// case counts when the format supports them, bounded secret-safe diagnostics).
// Evaluate is the fail-closed closure gate: missing/self/stale/skipped/
// no-execution/partial evidence cannot close a whole implementation.
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
// VerifierTreeSHA256 binds that attestation to the tested tree the verifier
// actually re-ran against.
type Provenance struct {
	Executor string `json:"executor"`
	Verifier string `json:"verifier"` // empty until an independent verifier attests
	// VerifierTreeSHA256 is the tested-tree digest the verifier attestation is
	// bound to. An attestation recorded against a different tree than the
	// report's is not evidence for this report.
	VerifierTreeSHA256 string `json:"verifier_tree_sha256,omitempty"`
}

// CommandEvidence binds one executed command to hashes of its exact input and
// bounded scrubbed output, plus executed-case counts when the format supports
// them. ExecutedCases is -1 when the format cannot prove a count.
type CommandEvidence struct {
	Command       string `json:"command"`
	CommandSHA256 string `json:"command_sha256"`
	OutputSHA256  string `json:"output_sha256"` // hash of the RAW captured output
	ExitCode      int    `json:"exit_code"`
	DurationMillis int64 `json:"duration_ms"`
	Format        string `json:"format"`
	ExecutedCases int    `json:"executed_cases"` // -1 = unknown/unsupported
	FailedCases   int    `json:"failed_cases"`   // -1 = unknown/unsupported
	SkippedCases  int    `json:"skipped_cases"`  // -1 = unknown/unsupported
	Diagnostics   string `json:"diagnostics"`    // bounded, secret-scrubbed tail
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
//     counts), and command/output hash binding;
//   - every record carries a PERSISTED independent verifier attestation
//     (Provenance.Verifier, set by AttestExecution after an independent
//     re-execution): non-empty, not the executor, matching the closing
//     verifier identity, and bound to this report's tested tree. A caller
//     supplying a different name in ClosureOptions.Verifier cannot conjure an
//     independent verdict — the attestation must already exist in the report.
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
			if rec.Provenance.VerifierTreeSHA256 == "" {
				reasons = append(reasons, p+"verifier attestation is not bound to a tested tree digest")
			} else if r.TreeSHA256 != "" && rec.Provenance.VerifierTreeSHA256 != r.TreeSHA256 {
				reasons = append(reasons, p+"verifier attestation is bound to a different tree than the one tested")
			}
		}
		// Count validity: anything below -1 is malformed in every format.
		if rec.Command.ExecutedCases < -1 || rec.Command.FailedCases < -1 || rec.Command.SkippedCases < -1 {
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
// CommandEvidence produced by that independent execution; it must show a
// passing structured run (exit 0, a case-proving format, >0 executed, 0 failed
// — an opaque shell exit-0 is not an attestation basis).
//
// This API exists so the only way to place a verifier name on a record is to
// supply evidence of an independent execution; there is deliberately no API
// that stamps a name without one. Runtime attribution is not authentication:
// the attestation makes a self verdict detectable and rejectable, it does not
// prove who a human operator is.
func AttestExecution(r *Report, criterion, verifier string, rerun CommandEvidence) error {
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
		if rerun.ExitCode != 0 {
			return fmt.Errorf("evidence: independent re-run of %q exited %d — not an attestation basis", criterion, rerun.ExitCode)
		}
		if rerun.Format != FormatGoTestJSON && rerun.Format != FormatEnvelope {
			return fmt.Errorf("evidence: independent re-run of %q has unstructured %q output — not an attestation basis", criterion, formatOrNone(rerun.Format))
		}
		if rerun.ExecutedCases <= 0 || rerun.FailedCases != 0 {
			return fmt.Errorf("evidence: independent re-run of %q executed %d cases with %d failures — not an attestation basis", criterion, rerun.ExecutedCases, rerun.FailedCases)
		}
		r.Records[i].Provenance.Verifier = verifier
		r.Records[i].Provenance.VerifierTreeSHA256 = r.TreeSHA256
		return nil
	}
	return fmt.Errorf("evidence: no record for criterion %q", criterion)
}

func formatOrNone(s string) string {
	if s == "" {
		return "(none)"
	}
	return strings.TrimSpace(s)
}
