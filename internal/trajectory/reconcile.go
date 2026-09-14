package trajectory

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"time"

	"gopkg.in/yaml.v3"
	"parley-deck-cli/internal/budget"
	"parley-deck-cli/internal/telemetry"
)

// HelperRequest and ParentResult preserve the existing private CLI wire format.
// The parent result is retained observation evidence, never a resolution grant.
type HelperRequest struct {
	Version      int                `json:"version"`
	Ticket       VerificationTicket `json:"ticket"`
	Participants []string           `json:"participants"`
	Criteria     []Criterion        `json:"criteria"`
}

type ParentResult struct {
	Version           int         `json:"version"`
	RunID             string      `json:"run_id"`
	RequestPath       string      `json:"request_path"`
	RequestSHA256     string      `json:"request_sha256"`
	InvocationID      string      `json:"invocation_id"`
	TerminalSHA256    string      `json:"terminal_sha256"`
	ReceiptSHA256     string      `json:"receipt_sha256"`
	Assessment        *Assessment `json:"assessment"`
	FailureStage      string      `json:"failure_stage"`
	TrajectoryPending bool        `json:"trajectory_pending"`
}

type ReconciliationPreview struct {
	Version        int                `json:"version"`
	Root           string             `json:"root"`
	StateSHA256    string             `json:"state_sha256"`
	Sequence       int                `json:"sequence"`
	ChargeKey      string             `json:"charge_key"`
	RunID          string             `json:"run_id"`
	ParentSHA256   string             `json:"parent_sha256"`
	Assessment     Assessment         `json:"assessment"`
	RecoverySHA256 string             `json:"recovery_sha256,omitempty"`
	Unchanged      *UnchangedEvidence `json:"unchanged,omitempty"`
}

func (p ReconciliationPreview) SHA256() string { data, _ := canonical(p); return digest(data) }

type Resolution struct {
	Preview ReconciliationPreview `json:"preview"`
	SHA256  string                `json:"sha256"`
	At      time.Time             `json:"at"`
}

func runtimeID(s string) bool {
	return safeLabel(s) && !strings.ContainsAny(s, `/\`) && s != "." && s != ".."
}

func canonicalRoot(root string) (string, error) {
	root, err := filepath.Abs(root)
	if err != nil {
		return "", err
	}
	return filepath.EvalSymlinks(root)
}

// Private runtime records have either of the two existing exact encodings:
// indented JSON (parent files) or that JSON plus a newline (telemetry/journals).
// Duplicates, aliases, unknown fields and other ambiguous encodings refuse.
func readReconciliationJSON(dir *os.Root, path string, value any) (string, error) {
	info, err := dir.Lstat(path)
	if err != nil {
		return "", err
	}
	if !info.Mode().IsRegular() || info.Size() > 16<<20 {
		return "", errors.New("reconciliation requires a bounded regular artifact")
	}
	f, err := dir.Open(path)
	if err != nil {
		return "", err
	}
	defer f.Close()
	opened, err := f.Stat()
	if err != nil || !os.SameFile(info, opened) {
		return "", errors.New("reconciliation artifact changed during open")
	}
	data, err := io.ReadAll(io.LimitReader(f, (16<<20)+1))
	if err != nil || len(data) > 16<<20 {
		return "", errors.New("reconciliation artifact exceeds its bound")
	}
	dec := json.NewDecoder(bytes.NewReader(data))
	dec.DisallowUnknownFields()
	if err = dec.Decode(value); err != nil {
		return "", err
	}
	want, err := json.MarshalIndent(value, "", "  ")
	if err != nil || (!bytes.Equal(data, want) && !bytes.Equal(data, append(want, '\n'))) {
		return "", errors.New("reconciliation artifact is ambiguous or noncanonical")
	}
	return digest(data), nil
}

func sameJSON(a, b any) bool {
	x, _ := json.Marshal(a)
	y, _ := json.Marshal(b)
	return bytes.Equal(x, y)
}

// Recheck the original material scope on every durable resolution read. This
// parser requires named checks and exact quorum/order, never scalar fallback.
func checkReconciliationScope(dir *os.Root, req HelperRequest) error {
	if !runtimeID(req.Ticket.Request.Idea) {
		return errors.New("invalid reconciliation idea path")
	}
	scope, err := readReconciliationScope(dir, req.Ticket.Request.Idea)
	if err != nil {
		return err
	}
	r := req.Ticket.Request
	if !slices.Equal(scope.Participants, req.Participants) || !slices.Contains(scope.Participants, r.Implementer) || !slices.Contains(scope.Participants, r.Verifier) || len(scope.Criteria) != len(r.Criteria) || len(req.Criteria) != len(r.Criteria) {
		return errors.New("reconciliation quorum or material scope changed")
	}
	for i, c := range scope.Criteria {
		if c.Name != r.Criteria[i].Name || digest([]byte(c.Command)) != r.Criteria[i].CommandSHA256 || req.Criteria[i].Name != c.Name || req.Criteria[i].Command != c.Command {
			return errors.New("reconciliation changed original criterion command or order")
		}
	}
	return nil
}

type reconciliationScope struct {
	Participants []string
	Criteria     []Criterion
}

func readReconciliationScope(dir *os.Root, idea string) (reconciliationScope, error) {
	var scope reconciliationScope
	if !runtimeID(idea) {
		return scope, errors.New("invalid reconciliation idea path")
	}
	f, err := dir.Open(filepath.Join("parley-deck", "ideas", idea, "00-prompt.md"))
	if err != nil {
		return scope, err
	}
	defer f.Close()
	info, err := f.Stat()
	if err != nil || !info.Mode().IsRegular() || info.Size() > 1<<20 {
		return scope, errors.New("invalid reconciliation scope file")
	}
	raw, err := io.ReadAll(io.LimitReader(f, (1<<20)+1))
	if err != nil || len(raw) > 1<<20 {
		return scope, errors.New("reconciliation scope exceeds its bound")
	}
	return parseReconciliationScope(raw)
}

func reconciliationFrontmatter(raw []byte) ([]byte, error) {
	lines := strings.Split(strings.ReplaceAll(string(raw), "\r\n", "\n"), "\n")
	if len(lines) < 3 || lines[0] != "---" {
		return nil, errors.New("reconciliation scope lacks frontmatter")
	}
	end := 1
	for end < len(lines) && lines[end] != "---" {
		end++
	}
	if end == len(lines) {
		return nil, errors.New("reconciliation scope has unterminated frontmatter")
	}
	return []byte(strings.Join(lines[1:end], "\n")), nil
}

func parseReconciliationScope(raw []byte) (reconciliationScope, error) {
	var scope reconciliationScope
	frontmatter, err := reconciliationFrontmatter(raw)
	if err != nil {
		return scope, err
	}
	var fm struct {
		Participants []string `yaml:"participants"`
		Checks       []struct {
			Name    string `yaml:"name"`
			Command string `yaml:"command"`
		} `yaml:"checks"`
	}
	if err := yaml.Unmarshal(frontmatter, &fm); err != nil {
		return scope, err
	}
	scope.Participants = fm.Participants
	for _, c := range fm.Checks {
		scope.Criteria = append(scope.Criteria, Criterion{strings.TrimSpace(c.Name), strings.TrimSpace(c.Command)})
	}
	return scope, nil
}

// Derivation reuses the actual requested/started/terminal lifecycle and the
// complete helper journal. It never executes a command or adopts a self verdict.
type ParentDerivation struct {
	CompletedAt     time.Time    `json:"completed_at"`
	Sequence        int          `json:"sequence"`
	ChargeKey       string       `json:"charge_key"`
	Result          ParentResult `json:"result"`
	RequestedSHA256 string       `json:"requested_sha256"`
	StartedSHA256   string       `json:"started_sha256"`
	TicketSHA256    string       `json:"ticket_sha256"`
	LaunchSHA256    string       `json:"launch_sha256"`
}

func deriveParentEvidence(ctx context.Context, b budget.CycleBinding, s State, root, runID string) (ParentDerivation, error) {
	var derived ParentDerivation
	if !runtimeID(runID) {
		return derived, errors.New("invalid verification run identity")
	}
	origin, err := canonicalRoot(root)
	if err != nil || origin != root {
		return derived, errors.New("verification origin changed")
	}
	originBinding, err := budget.LoadCycleBinding(ctx, root, s.Policy.Idea, budget.Fixup)
	if err != nil || originBinding == nil || originBinding.Store.Dir != b.Store.Dir || originBinding.Policy.TrajectorySHA256 != b.Policy.TrajectorySHA256 {
		return derived, errors.New("verification origin differs from shared trajectory authority")
	}
	dir, err := os.OpenRoot(root)
	if err != nil {
		return derived, err
	}
	defer dir.Close()
	base := filepath.Join(".parley-runtime", "trajectory-verification", runID)
	var req HelperRequest
	requestSHA, err := readReconciliationJSON(dir, filepath.Join(base, "request.json"), &req)
	if err != nil || req.Version != 1 || req.Ticket.Root != root || req.Ticket.RunID != runID {
		return derived, errors.New("original parent request is unavailable or changed")
	}
	r, err := capturedRequestAt(s, req.Ticket.Request.Verifier, req.Ticket.Request.Sequence)
	if err != nil || !sameJSON(r, req.Ticket.Request) {
		return derived, errors.New("parent evidence differs from its original charged patch")
	}
	members, err := checkCapturedActivationQuorum(ctx, b, s, r)
	if err != nil {
		return derived, err
	}
	if !slices.Equal(req.Participants, members) {
		return derived, errors.New("helper request differs from the original activation quorum")
	}
	if err = checkReconciliationScope(dir, req); err != nil {
		return derived, err
	}
	journal, err := openVerificationDirectory(b, r.Charge.EntryKey, false)
	if err != nil {
		return derived, err
	}
	defer journal.Close()
	var ticket VerificationTicket
	derived.TicketSHA256, err = readVerificationArtifact(journal, "request.json", &ticket)
	want, hashErr := req.Ticket.SHA256()
	if err != nil || hashErr != nil || derived.TicketSHA256 != want {
		return derived, errors.New("shared verification ticket differs from parent request")
	}
	var launch verificationLaunch
	derived.LaunchSHA256, err = readVerificationArtifact(journal, "launch.json", &launch)
	if err != nil || launch.Version != 1 || launch.TicketSHA256 != want || !runtimeID(launch.InvocationID) || launch.At.IsZero() {
		return derived, errors.New("original verifier launch is unavailable")
	}
	invocation := launch.InvocationID
	var terminal, requested, started telemetry.Record
	invBase := filepath.Join(".parley-runtime", "invocations", invocation)
	terminalSHA, err := readReconciliationJSON(dir, filepath.Join(invBase, "terminal.json"), &terminal)
	if err != nil || terminal.SchemaVersion != telemetry.SchemaVersion || terminal.Type != "invocation.terminal" || terminal.InvocationID != invocation || terminal.Metadata.RunID != runID || terminal.Metadata.Idea != s.Policy.Idea || terminal.Metadata.Agent != r.Verifier || terminal.Metadata.Phase != "trajectory-verification" || terminal.Metadata.LaunchMode != "headless" || terminal.RequestedAt.IsZero() || terminal.StartedAt == nil || terminal.CompletedAt == nil || terminal.PID == nil || *terminal.PID <= 0 || terminal.StartedAt.Before(terminal.RequestedAt) || terminal.CompletedAt.Before(*terminal.StartedAt) || terminal.Outcome == nil || terminal.Outcome.Status != "process-exited" || terminal.Outcome.FailureClass != nil || terminal.Outcome.ExitCode == nil || *terminal.Outcome.ExitCode != 0 {
		return derived, errors.New("successful observed parent terminal binding is unavailable")
	}
	derived.RequestedSHA256, err = readReconciliationJSON(dir, filepath.Join(invBase, "requested.json"), &requested)
	if err != nil || requested.SchemaVersion != terminal.SchemaVersion || requested.Type != "invocation.requested" || requested.InvocationID != invocation || !sameJSON(requested.Metadata, terminal.Metadata) || !requested.RequestedAt.Equal(terminal.RequestedAt) || requested.StartedAt != nil || requested.CompletedAt != nil || requested.PID != nil || requested.Outcome != nil {
		return derived, errors.New("original requested lifecycle differs from the verifier terminal")
	}
	derived.StartedSHA256, err = readReconciliationJSON(dir, filepath.Join(invBase, "started.json"), &started)
	if err != nil || started.SchemaVersion != terminal.SchemaVersion || started.Type != "invocation.started" || started.InvocationID != invocation || !sameJSON(started.Metadata, terminal.Metadata) || !started.RequestedAt.Equal(terminal.RequestedAt) || started.StartedAt == nil || !started.StartedAt.Equal(*terminal.StartedAt) || started.PID == nil || *started.PID != *terminal.PID || started.CompletedAt != nil || started.Outcome != nil || launch.At.Before(terminal.RequestedAt) || launch.At.After(*terminal.StartedAt) {
		return derived, errors.New("original started lifecycle differs from the verifier terminal")
	}
	receipt, observation, err := readCapturedVerificationJournal(journal, ticket, invocation)
	if err != nil {
		return derived, err
	}
	encoded, err := canonical(receipt)
	if err != nil || receipt.FinishedAt.Before(*terminal.StartedAt) || receipt.FinishedAt.After(*terminal.CompletedAt) {
		return derived, errors.New("helper receipt is outside its observed invocation")
	}
	assessment, err := AssessCaptured(r, observation)
	if err != nil {
		return derived, err
	}
	derived.CompletedAt = *terminal.CompletedAt
	derived.Sequence, derived.ChargeKey = r.Sequence, r.Charge.EntryKey
	derived.Result = ParentResult{Version: 1, RunID: runID, RequestPath: filepath.Join(root, base, "request.json"), RequestSHA256: requestSHA, InvocationID: invocation, TerminalSHA256: terminalSHA, ReceiptSHA256: digest(encoded), Assessment: &assessment, TrajectoryPending: true}
	return derived, nil
}

func parentPreview(s State, root, runID, parentSHA, recoverySHA string, derived ParentDerivation) ReconciliationPreview {
	stateBytes, _ := canonical(s)
	sequence := derived.Sequence
	return ReconciliationPreview{Version: 1, Root: root, StateSHA256: digest(stateBytes), Sequence: sequence, ChargeKey: s.Attempts[sequence-1].Charge.EntryKey, RunID: runID, ParentSHA256: parentSHA, Assessment: *derived.Result.Assessment, RecoverySHA256: recoverySHA}
}

// Called under the cycle guard after state/archive validation. Recovery retains
// a separate provenance record; an ordinary parent continues to bind its bytes.
func readParentEvidence(ctx context.Context, b budget.CycleBinding, s State, root, runID string) (ReconciliationPreview, error) {
	derived, err := deriveParentEvidence(ctx, b, s, root, runID)
	if err != nil {
		return ReconciliationPreview{}, err
	}
	dir, err := os.OpenRoot(root)
	if err != nil {
		return ReconciliationPreview{}, err
	}
	defer dir.Close()
	base := filepath.Join(".parley-runtime", "trajectory-verification", runID)
	recovery, recoverySHA, err := readParentRecovery(dir, base)
	if err == nil {
		if err = validateParentRecovery(dir, base, root, s, derived, recovery); err != nil {
			return ReconciliationPreview{}, err
		}
		return parentPreview(s, root, runID, recovery.Preview.ParentSHA256, recoverySHA, derived), nil
	}
	if !os.IsNotExist(err) {
		return ReconciliationPreview{}, err
	}
	var result ParentResult
	parentSHA, err := readReconciliationJSON(dir, filepath.Join(base, "parent-result.json"), &result)
	if err != nil {
		return ReconciliationPreview{}, err
	}
	if !sameJSON(result, derived.Result) {
		return ReconciliationPreview{}, errors.New("parent did not retain the independently derived successful comparison")
	}
	return parentPreview(s, root, runID, parentSHA, "", derived), nil
}

func compareReconciledParent(actual, original ReconciliationPreview) error {
	actual.StateSHA256 = original.StateSHA256
	// Pre-recovery resolutions already pin the original complete parent bytes.
	// A derived replacement must reproduce those exact bytes and every fact.
	// New resolutions additionally pin their recovery provenance record.
	if original.RecoverySHA256 == "" {
		actual.RecoverySHA256 = ""
	}
	if !sameJSON(actual, original) {
		return errors.New("reconciled patch facts changed")
	}
	return nil
}

func checkResolutions(ctx context.Context, b budget.CycleBinding, s State) error {
	for _, resolution := range s.Resolutions {
		p := resolution.Preview
		actual, err := readResolutionEvidence(ctx, b, s, p)
		if err != nil {
			return fmt.Errorf("reconciled patch %d lost evidence: %w", p.Sequence, err)
		}
		// StateSHA256 records the exact pre-transition preview, not today's
		// later ledger. Every other fact must still be independently derived.
		if err := compareReconciledParent(actual, p); err != nil {
			return err
		}
	}
	return nil
}

func readResolutionEvidence(ctx context.Context, b budget.CycleBinding, s State, p ReconciliationPreview) (ReconciliationPreview, error) {
	if p.Unchanged != nil {
		return unchangedPreview(ctx, b, s, p.Root, p.Sequence)
	}
	return readParentEvidence(ctx, b, s, p.Root, p.RunID)
}

func PreviewReconciliation(ctx context.Context, root, idea, runID string) (ReconciliationPreview, error) {
	var p ReconciliationPreview
	root, err := canonicalRoot(root)
	if err != nil {
		return p, err
	}
	err = withState(ctx, root, idea, func(b budget.CycleBinding, _ budget.Snapshot, s State) error {
		var err error
		p, err = readParentEvidence(ctx, b, s, root, runID)
		if err != nil {
			return err
		}
		if p.Sequence <= len(s.Resolutions) {
			old := s.Resolutions[p.Sequence-1].Preview
			if old.RunID != runID || old.Root != root {
				return errors.New("charge was reconciled from a different parent")
			}
			p = old
		} else if p.Sequence != len(s.Resolutions)+1 {
			return errors.New("reconciliation cannot omit an earlier charged attempt")
		}
		return nil
	})
	if err == nil && p.Version == 0 {
		err = errors.New("trajectory is not active")
	}
	return p, err
}

func Reconcile(ctx context.Context, root, idea, runID, expected string) error {
	root, err := canonicalRoot(root)
	if err != nil {
		return err
	}
	observed := false
	err = withState(ctx, root, idea, func(b budget.CycleBinding, ledger budget.Snapshot, s State) error {
		observed = true
		p, err := readParentEvidence(ctx, b, s, root, runID)
		if err != nil {
			return err
		}
		if p.Sequence <= len(s.Resolutions) {
			old := s.Resolutions[p.Sequence-1]
			if old.Preview.Root == root && old.Preview.RunID == runID && old.SHA256 == expected {
				return nil
			}
			return errors.New("reconciliation replay changed its original decision")
		}
		if p.Sequence != len(s.Resolutions)+1 || expected != p.SHA256() {
			return errors.New("trajectory changed since reconciliation preview")
		}
		s.Version = 3
		s.Resolutions = append(s.Resolutions, Resolution{p, expected, time.Now().UTC()})
		if err = validateState(s, b, ledger); err != nil {
			return err
		}
		return writeState(statePath(b), s)
	})
	if err == nil && !observed {
		return errors.New("trajectory is not active")
	}
	return err
}
