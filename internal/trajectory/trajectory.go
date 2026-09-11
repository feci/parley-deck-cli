// Package trajectory derives patch-regression observations from actual paired
// criterion executions. It does not grant completion or authenticate actors.
package trajectory

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"parley-deck-cli/internal/evidence"
	"parley-deck-cli/internal/telemetry"
)

const MaxPatches = 128
const MaxCriteria = 64

type Criterion struct {
	Name    string
	Command string // runtime input only; requests retain its exact hash
}

type CriterionBinding struct {
	Name          string `json:"name"`
	CommandSHA256 string `json:"command_sha256"`
}

type Tree struct {
	Commit string `json:"commit"`
	SHA256 string `json:"sha256"`
}

// Request is frozen by the orchestrator before independent verification.
// Sequence/PreviousSHA256 must describe every patch in the opted-in trajectory;
// Evaluate requires that complete expected list, separately from observations.
// Every criterion is an original material acceptance obligation, not a new
// severity classification chosen after a failure.
type Request struct {
	Version        int                `json:"version"`
	Idea           string             `json:"idea"`
	ID             string             `json:"id"`
	Sequence       int                `json:"sequence"`
	PreviousSHA256 string             `json:"previous_sha256"`
	Implementer    string             `json:"implementer"`
	Verifier       string             `json:"verifier"`
	Before         Tree               `json:"before"`
	After          Tree               `json:"after"`
	PatchSHA256    string             `json:"patch_sha256"`
	Criteria       []CriterionBinding `json:"criteria"`
}

type Execution struct {
	Complete         bool                     `json:"complete"`
	Record           evidence.CriterionRecord `json:"record"`
	TreeBeforeSHA256 string                   `json:"tree_before_sha256"`
	TreeAfterSHA256  string                   `json:"tree_after_sha256"`
}

type Pair struct {
	Name   string       `json:"name"`
	Before [2]Execution `json:"before"`
	After  [2]Execution `json:"after"`
}

// Observation contains execution facts, never a caller-supplied pass/fail vote.
// The caller must bind the selected runtime verifier to its actual invocation;
// names and mutually consistent data are not authentication against same-UID
// fabrication. Ordinary completion still uses evidence.AttestExecution.
type Observation struct {
	Version       int    `json:"version"`
	RequestSHA256 string `json:"request_sha256"`
	Verifier      string `json:"verifier"`
	Pairs         []Pair `json:"pairs"`
}

type Outcome string

const (
	Regression   Outcome = "confirmed-regression"
	NoRegression Outcome = "confirmed-no-regression"
	Inconclusive Outcome = "inconclusive"
)

type Assessment struct {
	Outcome    Outcome  `json:"outcome"`
	Regressed  []string `json:"regressed_criteria"`
	Unresolved []string `json:"unresolved_criteria"`
}

type Decision struct {
	ReviewRequired  bool         `json:"review_required"`
	TriggerSequence int          `json:"trigger_sequence"`
	Consecutive     int          `json:"consecutive"`
	Pending         bool         `json:"pending"`
	Assessments     []Assessment `json:"assessments"`
}

func digest(b []byte) string { h := sha256.Sum256(b); return hex.EncodeToString(h[:]) }
func validHash(s string) bool {
	b, err := hex.DecodeString(s)
	return err == nil && len(b) == 32 && s == strings.ToLower(s)
}
func validCommit(s string) bool {
	b, err := hex.DecodeString(s)
	return err == nil && (len(b) == 20 || len(b) == 32) && s == strings.ToLower(s)
}
func safeLabel(s string) bool {
	label := telemetry.SafeLabel(s)
	return s != "" && len(s) <= 128 && strings.TrimSpace(s) == s &&
		!strings.ContainsAny(s, "\r\n\x00") && label != nil && *label == s
}

func (r Request) validate() error {
	if r.Version != 1 || !safeLabel(r.Idea) || !safeLabel(r.ID) || r.Sequence < 1 || r.Sequence > MaxPatches ||
		!safeLabel(r.Implementer) || !safeLabel(r.Verifier) || r.Implementer == r.Verifier ||
		!validCommit(r.Before.Commit) || !validCommit(r.After.Commit) || r.Before.Commit == r.After.Commit ||
		!validHash(r.Before.SHA256) || !validHash(r.After.SHA256) || r.Before.SHA256 == r.After.SHA256 ||
		!validHash(r.PatchSHA256) || len(r.Criteria) == 0 || len(r.Criteria) > MaxCriteria {
		return errors.New("invalid patch request or non-independent verifier")
	}
	if (r.Sequence == 1 && r.PreviousSHA256 != "") || (r.Sequence > 1 && !validHash(r.PreviousSHA256)) {
		return errors.New("invalid patch sequence ancestry")
	}
	seen := map[string]bool{}
	for _, c := range r.Criteria {
		if !safeLabel(c.Name) || !validHash(c.CommandSHA256) || seen[c.Name] {
			return errors.New("invalid or duplicate material criterion binding")
		}
		seen[c.Name] = true
	}
	return nil
}

func (r Request) SHA256() (string, error) {
	if err := r.validate(); err != nil {
		return "", err
	}
	data, err := json.Marshal(r)
	if err != nil {
		return "", err
	}
	return digest(data), nil
}

// classify only accepts complete structured execution observations. An exit
// failure without a structured failing test/package is not a proven regression.
// Existing failure on BOTH trees is inconclusive: aggregate counts cannot prove
// that the same individual failures persisted, or exclude a new failure.
func classify(e Execution, c CriterionBinding, tree Tree, verifier string) (evidence.Status, error) {
	r, cmd := e.Record, e.Record.Command
	if !e.Complete || r.Name != c.Name || r.Provenance.Executor != verifier || r.Provenance.Verifier != "" || r.Provenance.VerifierRerun != nil ||
		e.TreeBeforeSHA256 != tree.SHA256 || e.TreeAfterSHA256 != tree.SHA256 ||
		cmd.CommandSHA256 != c.CommandSHA256 || !validHash(cmd.OutputSHA256) || cmd.DurationMillis < 0 ||
		cmd.Command == "" || len(cmd.Command) > 4096 || len(cmd.Diagnostics) > 8192 {
		return "", errors.New("execution is not bound to the exact criterion, stable tree and selected verifier")
	}
	if cmd.Format != evidence.FormatGoTestJSON && cmd.Format != evidence.FormatEnvelope {
		return "", errors.New("opaque command output cannot confirm a patch outcome")
	}
	if cmd.ExecutedCases < 0 || cmd.FailedCases < 0 || cmd.FailedCases > cmd.ExecutedCases || cmd.SkippedCases < -1 || cmd.FailedPackages < -1 {
		return "", errors.New("missing or contradictory structured execution counts")
	}
	if cmd.Format == evidence.FormatGoTestJSON && (cmd.SkippedCases < 0 || cmd.FailedPackages < 0) {
		return "", errors.New("missing Go test package/skip counts")
	}
	if cmd.Format == evidence.FormatEnvelope && (cmd.SkippedCases != -1 || cmd.FailedPackages != -1) {
		return "", errors.New("envelope carries unsupported package/skip counts")
	}
	failed := cmd.FailedCases > 0 || (cmd.Format == evidence.FormatGoTestJSON && cmd.FailedPackages > 0)
	if failed {
		if r.Status != evidence.StatusFail || cmd.ExitCode < 0 {
			return "", errors.New("contradictory structured failure or unavailable process outcome")
		}
		return evidence.StatusFail, nil
	}
	if r.Status != evidence.StatusPass || cmd.ExitCode != 0 || cmd.ExecutedCases == 0 {
		return "", errors.New("no complete executed pass or structured failure")
	}
	return evidence.StatusPass, nil
}

func consistent(a, b Execution) bool {
	x, y := a.Record.Command, b.Record.Command
	return a.Record.Status == b.Record.Status && x.Format == y.Format && x.CommandSHA256 == y.CommandSHA256 &&
		x.ExecutedCases == y.ExecutedCases && x.FailedCases == y.FailedCases && x.SkippedCases == y.SkippedCases &&
		x.FailedPackages == y.FailedPackages && x.ExitCode == y.ExitCode
}

func Assess(r Request, o Observation) (Assessment, error) {
	sha, err := r.SHA256()
	if err != nil {
		return Assessment{}, err
	}
	if o.Version != 1 || o.RequestSHA256 != sha || o.Verifier != r.Verifier || len(o.Pairs) != len(r.Criteria) {
		return Assessment{}, errors.New("missing or mismatched independent patch observation")
	}
	result := Assessment{Outcome: NoRegression}
	for i, c := range r.Criteria {
		p := o.Pairs[i]
		if p.Name != c.Name {
			return Assessment{}, errors.New("patch observation changed or reordered criterion scope")
		}
		var before, after evidence.Status
		for j := 0; j < 2; j++ {
			before, err = classify(p.Before[j], c, r.Before, r.Verifier)
			if err != nil {
				return Assessment{}, fmt.Errorf("baseline %q: %w", c.Name, err)
			}
			after, err = classify(p.After[j], c, r.After, r.Verifier)
			if err != nil {
				return Assessment{}, fmt.Errorf("patched %q: %w", c.Name, err)
			}
		}
		if !consistent(p.Before[0], p.Before[1]) || !consistent(p.After[0], p.After[1]) || p.Before[0].Record.Command.Format != p.After[0].Record.Command.Format {
			result.Unresolved = append(result.Unresolved, c.Name)
			continue
		}
		if before == evidence.StatusPass && after == evidence.StatusFail {
			result.Regressed = append(result.Regressed, c.Name)
		} else if after != evidence.StatusPass ||
			p.Before[0].Record.Command.ExecutedCases != p.After[0].Record.Command.ExecutedCases ||
			p.Before[0].Record.Command.SkippedCases != p.After[0].Record.Command.SkippedCases {
			result.Unresolved = append(result.Unresolved, c.Name)
		}
	}
	// One confirmed material regression suffices even if another criterion is
	// unresolved. An unresolved criterion can never establish a clean patch.
	if len(result.Regressed) > 0 {
		result.Outcome = Regression
	} else if len(result.Unresolved) > 0 {
		result.Outcome = Inconclusive
	}
	return result, nil
}

// Evaluate needs the full frozen patch list, not a filtered list of failures.
// A clean intervening patch resets the consecutive count. An inconclusive patch
// cannot supply a second confirmation or prove a clean interval. The first
// two-confirmation trigger stays recorded even if later input includes a clean
// patch; deciding to continue after that trigger is outside this API.
func Evaluate(expected []Request, observations []Observation) (Decision, error) {
	if len(expected) > MaxPatches || len(observations) != len(expected) {
		return Decision{}, errors.New("trajectory coverage is missing or exceeds its bound")
	}
	result := Decision{}
	seenIDs, seenPatches := map[string]bool{}, map[string]bool{}
	var previous string
	for i, r := range expected {
		sha, err := r.SHA256()
		if err != nil {
			return Decision{}, err
		}
		if r.Sequence != i+1 || r.PreviousSHA256 != previous || seenIDs[r.ID] || seenPatches[r.Before.Commit+":"+r.After.Commit] {
			return Decision{}, errors.New("trajectory has a gap, changed ancestry or replayed patch")
		}
		if i > 0 {
			p := expected[i-1]
			if r.Idea != p.Idea || r.Before != p.After || len(r.Criteria) != len(p.Criteria) {
				return Decision{}, errors.New("trajectory changed idea, patch continuity or material scope")
			}
			for j := range r.Criteria {
				if r.Criteria[j] != p.Criteria[j] {
					return Decision{}, errors.New("trajectory changed frozen material criteria")
				}
			}
		}
		seenIDs[r.ID], seenPatches[r.Before.Commit+":"+r.After.Commit] = true, true
		previous = sha
		a, err := Assess(r, observations[i])
		if err != nil {
			return Decision{}, fmt.Errorf("patch %d: %w", i+1, err)
		}
		result.Assessments = append(result.Assessments, a)
		switch a.Outcome {
		case Regression:
			result.Consecutive++
		case NoRegression:
			result.Consecutive = 0
		case Inconclusive:
			result.Pending, result.Consecutive = true, 0
		}
		if result.Consecutive >= 2 && !result.ReviewRequired {
			result.ReviewRequired, result.TriggerSequence = true, i+1
		}
	}
	return result, nil
}
