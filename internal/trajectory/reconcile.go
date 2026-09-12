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
	Version      int        `json:"version"`
	Root         string     `json:"root"`
	StateSHA256  string     `json:"state_sha256"`
	Sequence     int        `json:"sequence"`
	ChargeKey    string     `json:"charge_key"`
	RunID        string     `json:"run_id"`
	ParentSHA256 string     `json:"parent_sha256"`
	Assessment   Assessment `json:"assessment"`
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
	f, err := dir.Open(filepath.Join("parley-deck", "ideas", req.Ticket.Request.Idea, "00-prompt.md"))
	if err != nil {
		return err
	}
	defer f.Close()
	info, err := f.Stat()
	if err != nil || !info.Mode().IsRegular() || info.Size() > 1<<20 {
		return errors.New("invalid reconciliation scope file")
	}
	raw, err := io.ReadAll(io.LimitReader(f, (1<<20)+1))
	if err != nil || len(raw) > 1<<20 {
		return errors.New("reconciliation scope exceeds its bound")
	}
	lines := strings.Split(strings.ReplaceAll(string(raw), "\r\n", "\n"), "\n")
	if len(lines) < 3 || lines[0] != "---" {
		return errors.New("reconciliation scope lacks frontmatter")
	}
	end := 1
	for end < len(lines) && lines[end] != "---" {
		end++
	}
	if end == len(lines) {
		return errors.New("reconciliation scope has unterminated frontmatter")
	}
	var fm struct {
		Participants []string `yaml:"participants"`
		Checks       []struct {
			Name    string `yaml:"name"`
			Command string `yaml:"command"`
		} `yaml:"checks"`
	}
	if err = yaml.Unmarshal([]byte(strings.Join(lines[1:end], "\n")), &fm); err != nil {
		return err
	}
	r := req.Ticket.Request
	if !slices.Equal(fm.Participants, req.Participants) || !slices.Contains(fm.Participants, r.Implementer) || !slices.Contains(fm.Participants, r.Verifier) || len(fm.Checks) != len(r.Criteria) || len(req.Criteria) != len(r.Criteria) {
		return errors.New("reconciliation quorum or material scope changed")
	}
	for i, c := range fm.Checks {
		name, command := strings.TrimSpace(c.Name), strings.TrimSpace(c.Command)
		if name != r.Criteria[i].Name || digest([]byte(command)) != r.Criteria[i].CommandSHA256 || req.Criteria[i].Name != name || req.Criteria[i].Command != command {
			return errors.New("reconciliation changed original criterion command or order")
		}
	}
	return nil
}

// Called under the existing cycle guard with already checked state/snapshots.
// It never invokes a helper or trusts the parent assessment without recomputing.
func readParentEvidence(ctx context.Context, b budget.CycleBinding, s State, root, runID string) (ReconciliationPreview, error) {
	var preview ReconciliationPreview
	if !runtimeID(runID) {
		return preview, errors.New("invalid verification run identity")
	}
	origin, err := canonicalRoot(root)
	if err != nil || origin != root {
		return preview, errors.New("verification origin changed")
	}
	originBinding, err := budget.LoadCycleBinding(ctx, root, s.Policy.Idea, budget.Fixup)
	if err != nil || originBinding == nil || originBinding.Store.Dir != b.Store.Dir || originBinding.Policy.TrajectorySHA256 != b.Policy.TrajectorySHA256 {
		return preview, errors.New("verification origin differs from shared trajectory authority")
	}
	dir, err := os.OpenRoot(root)
	if err != nil {
		return preview, err
	}
	defer dir.Close()
	base := filepath.Join(".parley-runtime", "trajectory-verification", runID)
	var result ParentResult
	parentSHA, err := readReconciliationJSON(dir, filepath.Join(base, "parent-result.json"), &result)
	if err != nil {
		return preview, err
	}
	if result.Version != 1 || result.RunID != runID || !runtimeID(result.InvocationID) || !result.TrajectoryPending || result.FailureStage != "" || result.Assessment == nil || result.RequestPath != filepath.Join(root, base, "request.json") {
		return preview, errors.New("parent did not retain a successful independent comparison")
	}
	var req HelperRequest
	requestSHA, err := readReconciliationJSON(dir, filepath.Join(base, "request.json"), &req)
	if err != nil || requestSHA != result.RequestSHA256 || req.Version != 1 || req.Ticket.Root != root || req.Ticket.RunID != runID {
		return preview, errors.New("parent request binding changed")
	}
	r, err := capturedRequestAt(s, req.Ticket.Request.Verifier, req.Ticket.Request.Sequence)
	if err != nil || !sameJSON(r, req.Ticket.Request) {
		return preview, errors.New("parent evidence differs from its original charged patch")
	}
	if err = checkReconciliationScope(dir, req); err != nil {
		return preview, err
	}
	var terminal telemetry.Record
	terminalSHA, err := readReconciliationJSON(dir, filepath.Join(".parley-runtime", "invocations", result.InvocationID, "terminal.json"), &terminal)
	if err != nil || terminalSHA != result.TerminalSHA256 || terminal.SchemaVersion != telemetry.SchemaVersion || terminal.Type != "invocation.terminal" || terminal.InvocationID != result.InvocationID || terminal.Metadata.RunID != runID || terminal.Metadata.Idea != s.Policy.Idea || terminal.Metadata.Agent != r.Verifier || terminal.Metadata.Phase != "trajectory-verification" || terminal.Metadata.LaunchMode != "headless" || terminal.StartedAt == nil || terminal.CompletedAt == nil || terminal.PID == nil || *terminal.PID <= 0 || terminal.CompletedAt.Before(*terminal.StartedAt) || terminal.Outcome == nil || terminal.Outcome.Status != "process-exited" || terminal.Outcome.FailureClass != nil || terminal.Outcome.ExitCode == nil || *terminal.Outcome.ExitCode != 0 {
		return preview, errors.New("successful observed parent terminal binding is unavailable")
	}
	journal, err := openVerificationDirectory(b, r.Charge.EntryKey, false)
	if err != nil {
		return preview, err
	}
	defer journal.Close()
	var ticket VerificationTicket
	ticketSHA, err := readVerificationArtifact(journal, "request.json", &ticket)
	want, hashErr := req.Ticket.SHA256()
	if err != nil || hashErr != nil || ticketSHA != want {
		return preview, errors.New("shared verification ticket differs from parent request")
	}
	receipt, observation, err := readCapturedVerificationJournal(journal, ticket, result.InvocationID)
	if err != nil {
		return preview, err
	}
	encoded, err := canonical(receipt)
	if err != nil || digest(encoded) != result.ReceiptSHA256 || receipt.FinishedAt.Before(*terminal.StartedAt) || receipt.FinishedAt.After(*terminal.CompletedAt) {
		return preview, errors.New("helper receipt is outside its observed invocation")
	}
	assessment, err := AssessCaptured(r, observation)
	if err != nil || !sameJSON(assessment, *result.Assessment) {
		return preview, errors.New("parent assessment differs from retained executions")
	}
	stateBytes, _ := canonical(s)
	return ReconciliationPreview{1, root, digest(stateBytes), r.Sequence, r.Charge.EntryKey, runID, parentSHA, assessment}, nil
}

func checkResolutions(ctx context.Context, b budget.CycleBinding, s State) error {
	for _, resolution := range s.Resolutions {
		p := resolution.Preview
		actual, err := readParentEvidence(ctx, b, s, p.Root, p.RunID)
		if err != nil {
			return fmt.Errorf("reconciled patch %d lost evidence: %w", p.Sequence, err)
		}
		// StateSHA256 records the exact pre-transition preview, not today's
		// later ledger. Every other fact must still be independently derived.
		actual.StateSHA256 = p.StateSHA256
		if !sameJSON(actual, p) {
			return errors.New("reconciled patch facts changed")
		}
	}
	return nil
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
