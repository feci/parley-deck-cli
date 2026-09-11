package evidence

import (
	"encoding/hex"
	"fmt"
	"regexp"
	"strings"

	"gopkg.in/yaml.v3"
)

var completionStatusToken = regexp.MustCompile(`^[A-Za-z0-9][A-Za-z0-9_-]*$`)

func validCompletionDigest(value string) bool {
	decoded, err := hex.DecodeString(value)
	return err == nil && len(decoded) == 32 && value == strings.ToLower(value)
}

// completion.go — the exact, authorized BEFORE/AFTER completion status
// transition for evidence-bound documents (idea
// meta-protocol-change-evidence-first-efficiency, kimi-1 slice).
//
// A report's ExtraDigests binds ALL non-evidence scope, including the
// frontmatter status line. The deterministic `status: complete` write at
// completion would otherwise invalidate the attested report after the fact.
// The fix is NOT to exclude status/frontmatter from the binding (that would
// let a scope edit hide behind the exclusion) but to let the INDEPENDENT
// verifier — and only that verifier — explicitly authorize exactly one
// deterministic status transition, bound to the attested report, before it is
// applied. The gate then verifies the authorized transformation against the
// current content and independently recomputes the pre-transition state.
//
// Digest space: a transition lives in the SAME digest space as the report's
// ExtraDigests binding for its path (for IMPLEMENTATION.md, the content with
// only the generated `## Validation evidence` section removed). The caller
// supplies the exact bound bytes; digests are always recomputed from supplied
// content — caller-supplied digests are never accepted.
//
// Trusted-caller boundary (unchanged from AttestExecution/Evaluate): these
// APIs guarantee internal consistency and binding of the persisted evidence.
// They trust the calling orchestrator to pass the actual current document and
// to be the attesting verifier it names — a same-UID caller able to fabricate
// mutually consistent content is outside this boundary. Runtime attribution
// is not cryptographic authentication of a human (FINAL D3).

// CompletionTransition is the persisted authorization for exactly one
// deterministic frontmatter status transition, written ONLY by
// AuthorizeCompletionTransition. BeforeSHA256 is the digest of the exact
// content the verifier authorized from (which must equal the report's bound
// digest for Path at authorization time); AfterSHA256 is recomputed by
// applying the transformation, never supplied by the caller.
type CompletionTransition struct {
	Path         string `json:"path"`          // report-binding path (an ExtraDigests key)
	FromStatus   string `json:"from_status"`   // exact status value before the transition
	ToStatus     string `json:"to_status"`     // always "complete"
	BeforeSHA256 string `json:"before_sha256"` // digest of the authorized-from content
	AfterSHA256  string `json:"after_sha256"`  // recomputed digest of the transformed content
	AuthorizedBy string `json:"authorized_by"` // the attesting verifier's asserted runtime identity
}

// TransitionFrontmatterStatus rewrites exactly one normalized frontmatter
// `status:` field of doc to `to`, byte-exact everywhere else, and returns the
// prior value. It is the single deterministic transformation both the
// verifier (to compute the authorized AFTER state) and the completion writer
// (to apply it) MUST use, so the applied bytes are exactly the authorized
// bytes. Fail-closed: no/unterminated frontmatter, a missing, empty,
// duplicate, or non-normalized status line (`status: <value>` exactly), a
// multi-token/comment value, or a no-op (value already == to) is an error.
// The normalized-line requirement makes the transformation exactly invertible
// for a given FromStatus, which VerifyCompletionTransition relies on.
func TransitionFrontmatterStatus(doc []byte, to string) (out []byte, from string, err error) {
	if !completionStatusToken.MatchString(to) {
		return nil, "", fmt.Errorf("evidence: invalid transition target status %q", to)
	}
	lines := strings.Split(string(doc), "\n")
	if len(lines) == 0 || lines[0] != "---" {
		return nil, "", fmt.Errorf("evidence: document has no frontmatter block")
	}
	closing := -1
	statusIdx := -1
	value := ""
	for i := 1; i < len(lines); i++ {
		if lines[i] == "---" {
			closing = i
			break
		}
		if trimmed := strings.TrimSpace(lines[i]); strings.HasPrefix(trimmed, "status:") {
			if statusIdx >= 0 {
				return nil, "", fmt.Errorf("evidence: duplicate frontmatter status field")
			}
			statusIdx = i
			value = strings.TrimSpace(trimmed[len("status:"):])
		}
	}
	if closing < 0 {
		return nil, "", fmt.Errorf("evidence: frontmatter block is not terminated")
	}
	if statusIdx < 0 {
		return nil, "", fmt.Errorf("evidence: frontmatter has no status field")
	}
	if !completionStatusToken.MatchString(value) {
		return nil, "", fmt.Errorf("evidence: frontmatter status value %q is not a single plain token", value)
	}
	if lines[statusIdx] != "status: "+value {
		return nil, "", fmt.Errorf("evidence: frontmatter status line %q is not in the normalized form %q", lines[statusIdx], "status: "+value)
	}
	// Check YAML meaning as well as normalized bytes: quoted/spaced duplicate
	// keys, merge keys and aliases must not create a second semantic status.
	var parsed yaml.Node
	if err := yaml.Unmarshal([]byte(strings.Join(lines[1:closing], "\n")), &parsed); err != nil {
		return nil, "", fmt.Errorf("evidence: invalid frontmatter: %w", err)
	}
	if len(parsed.Content) != 1 || parsed.Content[0].Kind != yaml.MappingNode {
		return nil, "", fmt.Errorf("evidence: frontmatter must be a mapping")
	}
	fields, statusCount := parsed.Content[0].Content, 0
	for i := 0; i < len(fields); i += 2 {
		key, node := fields[i], fields[i+1]
		if key.Value == "<<" {
			return nil, "", fmt.Errorf("evidence: merged frontmatter cannot authorize status")
		}
		if key.Value == "status" {
			statusCount++
			if key.Kind != yaml.ScalarNode || node.Kind != yaml.ScalarNode || node.Tag != "!!str" || node.Value != value {
				return nil, "", fmt.Errorf("evidence: ambiguous frontmatter status")
			}
		}
	}
	if statusCount != 1 {
		return nil, "", fmt.Errorf("evidence: duplicate or missing YAML status")
	}
	if value == to {
		return nil, value, fmt.Errorf("evidence: status is already %q — no transition to authorize", to)
	}
	outLines := make([]string, len(lines))
	copy(outLines, lines)
	outLines[statusIdx] = "status: " + to
	return []byte(strings.Join(outLines, "\n")), value, nil
}

// TransitionStatusToComplete is the completion specialization: the only
// transition the evidence gate will ever honor.
func TransitionStatusToComplete(doc []byte) (out []byte, from string, err error) {
	return TransitionFrontmatterStatus(doc, "complete")
}

// AuthorizeCompletionTransition is the adapter the INDEPENDENTLY INVOKED
// verifier calls AFTER attesting every criterion (AttestExecution) and BEFORE
// Save. It authorizes exactly one deterministic `status: complete` transition
// of the document bound at path, starting from the actual current content —
// which MUST still be byte-identical to the state the report bound (drift
// invalidates the evidence; re-record instead). The AFTER digest is recomputed
// by applying the transformation here; a caller-supplied digest is never
// accepted. Refusals (fail-closed): nil report, missing verifier, a second
// transition, an unbound path, content drift, zero/unattested records, a
// verifier that is not the attesting verifier of every record, a self
// authorization (verifier == executor), or a malformed/untransitionable
// document. Old reports simply never carry the field: they remain valid only
// in their ORIGINAL state and are never silently upgraded.
func AuthorizeCompletionTransition(r *Report, path string, currentContent []byte, verifier string) error {
	if r == nil {
		return fmt.Errorf("evidence: cannot authorize a completion transition on a nil report")
	}
	if verifier == "" {
		return fmt.Errorf("evidence: a completion transition requires an authorizing verifier identity")
	}
	if r.CompletionTransition != nil {
		return fmt.Errorf("evidence: report already carries a completion transition — one transition per recorded report; re-record to change it")
	}
	bound, ok := r.ExtraDigests[path]
	if !ok || !validCompletionDigest(bound) {
		return fmt.Errorf("evidence: report binds no non-evidence digest for %q", path)
	}
	if len(r.Records) == 0 {
		return fmt.Errorf("evidence: no criterion records — a completion transition requires an independently attested report")
	}
	before := sha256Hex(currentContent)
	if before != bound {
		return fmt.Errorf("evidence: content for %q drifted since the evidence was recorded (recorded %s…, current %s…) — re-record instead of authorizing a transition",
			path, bound[:12], before[:12])
	}
	for _, rec := range r.Records {
		if rec.Provenance.Verifier == "" || rec.Provenance.VerifierRerun == nil {
			return fmt.Errorf("evidence: criterion %q carries no independent attestation — a completion transition requires a fully attested report", rec.Name)
		}
		if rec.Provenance.Verifier != verifier {
			return fmt.Errorf("evidence: criterion %q was attested by %s, not %s — the transition authorizer must be the attesting verifier", rec.Name, rec.Provenance.Verifier, verifier)
		}
		if rec.Provenance.Executor == verifier {
			return fmt.Errorf("evidence: self completion transition — the executor cannot authorize its own status transition")
		}
	}
	transformed, from, err := TransitionStatusToComplete(currentContent)
	if err != nil {
		return err
	}
	r.CompletionTransition = &CompletionTransition{
		Path:         path,
		FromStatus:   from,
		ToStatus:     "complete",
		BeforeSHA256: before,
		AfterSHA256:  sha256Hex(transformed),
		AuthorizedBy: verifier,
	}
	return nil
}

// VerifyCompletionTransition is the gate-side check for a scope digest
// mismatch: the bound non-evidence content no longer matches, and the ONLY
// admissible reconciliation is the recorded authorized completion transition.
// currentContent MUST be the current bytes in the transition's digest space.
// Every check is recomputed here; nothing in the record is taken on trust:
// the transition must bind this path, target status: complete from a
// well-formed source, be authorized by the SELECTED verifier, start from the
// recorded (bound) state, match the current content's digest, and —
// independently recomputed — invert the current content back to exactly the
// recorded state, proving the status flip is the ONLY change. Any tampering
// with the record (path, statuses, digests, authorizer) fails at least one
// recomputation. An empty result means the transition reconciles the scope.
func VerifyCompletionTransition(r *Report, path, boundDigest string, currentContent []byte, verifier string) []string {
	if r == nil || r.CompletionTransition == nil {
		return []string{"no authorized completion transition is recorded"}
	}
	tr := r.CompletionTransition
	var reasons []string
	if verifier == "" {
		reasons = append(reasons, "no selected verifier identity")
	}
	if tr.Path != path {
		reasons = append(reasons, fmt.Sprintf("transition binds %q, not %q", tr.Path, path))
	}
	if tr.AuthorizedBy == "" {
		reasons = append(reasons, "transition names no authorizing verifier")
	} else if verifier != "" && tr.AuthorizedBy != verifier {
		reasons = append(reasons, fmt.Sprintf("completion transition authorized by %s, not the selected verifier %s", tr.AuthorizedBy, verifier))
	}
	if tr.ToStatus != "complete" {
		reasons = append(reasons, fmt.Sprintf("transition target is %q, not status: complete", tr.ToStatus))
	}
	if tr.FromStatus == "" || tr.FromStatus == "complete" {
		reasons = append(reasons, fmt.Sprintf("malformed transition source status %q", tr.FromStatus))
	}
	if tr.BeforeSHA256 == "" || tr.BeforeSHA256 != boundDigest {
		reasons = append(reasons, "transition does not start from the recorded non-evidence state")
	}
	if !validCompletionDigest(boundDigest) || r.ExtraDigests[path] != boundDigest {
		reasons = append(reasons, "transition is not bound to the report's original content")
	}
	if tr.AfterSHA256 == "" || sha256Hex(currentContent) != tr.AfterSHA256 {
		reasons = append(reasons, "current content is not the authorized post-completion state")
	} else if tr.FromStatus != "" && tr.FromStatus != "complete" {
		before, actualStatus, err := TransitionFrontmatterStatus(currentContent, tr.FromStatus)
		if err != nil {
			reasons = append(reasons, "current content's status field is malformed: "+err.Error())
		} else if actualStatus != "complete" {
			reasons = append(reasons, "current content does not have the authorized complete status")
		} else if sha256Hex(before) != tr.BeforeSHA256 {
			reasons = append(reasons, "recomputed pre-transition content does not match the recorded state — the status flip is not the only change")
		}
	}
	return reasons
}
