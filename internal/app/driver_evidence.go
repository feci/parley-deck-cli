package app

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"path/filepath"
	"strings"

	"parley-deck-cli/internal/driver"
	"parley-deck-cli/internal/evidence"
)

// driver_evidence.go — typed-evidence close gate for the completion contract
// (idea meta-protocol-change-evidence-first-efficiency, kimi-1 claim).
//
// runChecksContract (driver_checks.go) is the per-cycle Phase 5/8 gate: it runs
// the `checks:` list and RECORDS typed evidence. This file is the closure gate:
// before the driver may accept `status: complete` for an idea with a list-form
// `checks:` contract, it must call EvidenceCloseGate with the runtime identity
// of an independent (non-implementer, non-executor) verifier. The gate is
// fail-closed and veto-only: it can reject a close claim, it can never create
// one. It performs no writes; the decision is computed from the persisted
// report plus a freshly recomputed tree digest.
//
// Wiring (Codex-owned driver state machine; this file is the API, the call
// site is theirs):
//
//	if criteria, isList, err := driver.ReadChecksContract(o.ideaDir); err == nil && isList && len(criteria) > 0 {
//	    gate := o.EvidenceCloseGate(independentVerifierID) // e.g. the goal-check / review-consensus non-implementer identity
//	    if !gate.Allowed {
//	        // veto the `status: complete` transition; escalate per §14 stopping
//	        // judgment — never auto-retry, never close:
//	        return fmt.Errorf("typed evidence gate: %s", strings.Join(gate.Reasons, "; "))
//	    }
//	}
//
// The verifier identity is an ASSERTED runtime identity (attribution), not
// authentication of a human; its purpose is that a self verdict is detectable
// and rejected.

// EvidenceGateResult is the close-gate decision.
type EvidenceGateResult struct {
	Allowed bool
	Reasons []string
}

// definedEvidenceArtifacts enumerates exactly which paths the tested-tree
// digest excludes — the driver-authored evidence artifacts for this idea (the
// IMPLEMENTATION.md evidence section carrier and the typed EVIDENCE.json
// report). Nothing else is excluded: any change to code, tests, config, or
// other deck files invalidates the tested-tree digest.
func definedEvidenceArtifacts(root, ideaDir string) ([]string, error) {
	rel, err := filepath.Rel(root, ideaDir)
	if err != nil || rel == ".." || strings.HasPrefix(rel, ".."+string(filepath.Separator)) || filepath.IsAbs(rel) {
		return nil, fmt.Errorf("idea directory %q is not inside the tree root %q", ideaDir, root)
	}
	return []string{
		filepath.ToSlash(filepath.Join(rel, "IMPLEMENTATION.md")),
		filepath.ToSlash(filepath.Join(rel, evidence.ReportFileName)),
	}, nil
}

// EvidenceCloseGate evaluates the persisted typed evidence against the
// contract's known criterion scope and the CURRENT tree. Missing/self/stale/
// skipped/zero-execution/partial evidence, an unknown (shell) output format,
// or a missing/corrupt report all deny closure.
func (o driverImplOps) EvidenceCloseGate(verifier string) EvidenceGateResult {
	deny := func(reasons ...string) EvidenceGateResult {
		return EvidenceGateResult{Allowed: false, Reasons: reasons}
	}
	criteria, isList, err := driver.ReadChecksContract(o.ideaDir)
	if err != nil {
		return deny("checks contract unreadable: " + err.Error())
	}
	if !isList || len(criteria) == 0 {
		return deny("no named checks contract — the criterion scope is unknown and closure cannot be evidenced")
	}
	scope := make([]string, 0, len(criteria))
	for _, c := range criteria {
		scope = append(scope, c.Name)
	}
	report, err := evidence.Load(o.ideaDir)
	if err != nil {
		return deny("typed evidence report unavailable: " + err.Error())
	}
	excl, err := definedEvidenceArtifacts(o.root, o.ideaDir)
	if err != nil {
		return deny(err.Error())
	}
	current, err := evidence.TreeDigest(o.root, excl...)
	if err != nil {
		return deny("current tree digest failed: " + err.Error())
	}
	// The tree digest excludes IMPLEMENTATION.md wholesale (the driver rewrites
	// its evidence section each cycle), so the report must separately bind the
	// NON-evidence remainder. A report without that binding, or with a binding
	// that no longer matches, is denied — a scope edit must not hide behind the
	// exclusion. The ONLY admissible mismatch is the recorded completion
	// transition: the attesting verifier's exact, pre-authorized `status:
	// complete` flip, independently recomputed against the current content
	// (evidence.VerifyCompletionTransition). Status/frontmatter stay bound —
	// they are reconciled, never excluded.
	restContent, implRel, err := implementationRestContent(o.root, o.ideaDir)
	if err != nil {
		return deny("non-evidence implementation digest failed: " + err.Error())
	}
	restSum := sha256.Sum256(restContent)
	restDigest := hex.EncodeToString(restSum[:])
	bound, ok := report.ExtraDigests[implRel]
	if !ok || bound == "" {
		return deny("report does not bind the non-evidence IMPLEMENTATION.md content — cannot rule out hidden scope edits")
	}
	if bound != restDigest {
		if trReasons := evidence.VerifyCompletionTransition(report, implRel, bound, restContent, verifier); len(trReasons) > 0 {
			return deny("IMPLEMENTATION.md non-evidence content changed after the evidence was recorded (" + strings.Join(trReasons, "; ") + ")")
		}
	}
	reasons := evidence.Evaluate(report, evidence.ClosureOptions{
		RequiredScope:     scope,
		CurrentTreeSHA256: current,
		Verifier:          verifier,
		// The completion contract is a code gate: opaque shell exit-0 output is
		// not semantically certified and cannot close a whole implementation.
		RequireStructured: true,
	})
	if len(reasons) > 0 {
		return deny(reasons...)
	}
	return EvidenceGateResult{Allowed: true}
}
