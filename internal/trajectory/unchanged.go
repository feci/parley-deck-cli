package trajectory

import (
	"context"
	"errors"
	"os"
	"path/filepath"
	"slices"
	"time"

	"parley-deck-cli/internal/budget"
	"parley-deck-cli/internal/fsutil"
	"parley-deck-cli/internal/telemetry"
)

// UnchangedEvidence is an observation of equal captured material bytes and an
// actually exited process. It is not a verifier run or evidence of acceptance.
// AttemptSHA256 retains the complete original sources, archives and charge,
// including different commit/status metadata on materially identical sources.
// Hashes and lifecycle records do not authenticate actors against same-UID edits.
type UnchangedEvidence struct {
	Version         int    `json:"version"`
	PolicySHA256    string `json:"policy_sha256"`
	AttemptSHA256   string `json:"attempt_sha256"`
	ScopeSHA256     string `json:"scope_sha256"`
	InvocationID    string `json:"invocation_id"`
	RequestedSHA256 string `json:"requested_sha256"`
	StartedSHA256   string `json:"started_sha256"`
	TerminalSHA256  string `json:"terminal_sha256"`
}

func unchangedAssessment(p Policy) Assessment {
	a := Assessment{Outcome: Inconclusive}
	for _, c := range p.Criteria {
		a.Unresolved = append(a.Unresolved, c.Name)
	}
	return a
}

func unchangedAttempt(s State, sequence int) (Attempt, error) {
	if sequence < 1 || sequence > len(s.Attempts) {
		return Attempt{}, errors.New("no original charged attempt at this sequence")
	}
	a := s.Attempts[sequence-1]
	if a.Launch == nil || !runtimeID(a.Launch.InvocationID) || a.Terminal == nil || a.Terminal.SnapshotError != "" || a.After == nil || a.AfterArchive == nil {
		return Attempt{}, errors.New("unchanged observation needs the original launch, terminal and source archives")
	}
	if a.Before.Tree.SHA256 != a.After.Tree.SHA256 {
		return Attempt{}, errors.New("changed material source requires independent patch verification")
	}
	return a, nil
}

func validateUnchangedPreview(s State, p ReconciliationPreview) error {
	u := p.Unchanged
	if u == nil || u.Version != 1 || p.RunID != "" || p.ParentSHA256 != "" || p.RecoverySHA256 != "" || !validHash(u.ScopeSHA256) || !validHash(u.RequestedSHA256) || !validHash(u.StartedSHA256) || !validHash(u.TerminalSHA256) {
		return errors.New("unchanged observation has invalid or mixed verifier provenance")
	}
	a, err := unchangedAttempt(s, p.Sequence)
	if err != nil {
		return err
	}
	policySHA, err := s.Policy.SHA256()
	if err != nil {
		return err
	}
	data, err := canonical(a)
	if err != nil {
		return err
	}
	if u.PolicySHA256 != policySHA || u.AttemptSHA256 != digest(data) || u.InvocationID != a.Launch.InvocationID || !sameJSON(p.Assessment, unchangedAssessment(s.Policy)) {
		return errors.New("unchanged observation changed its original attempt, scope or unresolved assessment")
	}
	return nil
}

// Read the original scope from retained before-source bytes, not current
// frontmatter. Current contract order, commands and quorum must still match.
func unchangedScope(ctx context.Context, b budget.CycleBinding, s State, a Attempt, dir *os.Root) (string, error) {
	name := "parley-deck/ideas/" + s.Policy.Idea + "/00-prompt.md"
	raw, err := readSnapshotMember(ctx, snapshotDirectory(b), a.BeforeArchive, a.Before, name, 1<<20)
	if err != nil {
		return "", err
	}
	original, err := parseReconciliationScope(raw)
	if err != nil {
		return "", err
	}
	if len(original.Participants) < 2 || !slices.Contains(original.Participants, s.Policy.Implementer) || len(original.Criteria) != len(s.Policy.Criteria) {
		return "", errors.New("archived original scope differs from the frozen policy")
	}
	seen := map[string]bool{}
	for _, id := range original.Participants {
		if !safeLabel(id) || seen[id] {
			return "", errors.New("invalid archived original quorum")
		}
		seen[id] = true
	}
	for i, c := range original.Criteria {
		if c.Name != s.Policy.Criteria[i].Name || digest([]byte(c.Command)) != s.Policy.Criteria[i].CommandSHA256 {
			return "", errors.New("archived material criterion differs from the frozen policy")
		}
	}
	current, err := readReconciliationScope(dir, s.Policy.Idea)
	if err != nil {
		return "", err
	}
	if !sameJSON(current, original) {
		return "", errors.New("current quorum or material scope differs from the original archive")
	}
	return digest(raw), nil
}

func unchangedPreview(ctx context.Context, b budget.CycleBinding, s State, root string, sequence int) (ReconciliationPreview, error) {
	var p ReconciliationPreview
	origin, err := canonicalRoot(root)
	if err != nil || origin != root {
		return p, errors.New("unchanged observation origin changed")
	}
	originBinding, err := budget.LoadCycleBinding(ctx, root, s.Policy.Idea, budget.Fixup)
	if err != nil || originBinding == nil || originBinding.Store.Dir != b.Store.Dir || originBinding.Policy.TrajectorySHA256 != b.Policy.TrajectorySHA256 {
		return p, errors.New("unchanged observation origin differs from shared trajectory authority")
	}
	a, err := unchangedAttempt(s, sequence)
	if err != nil {
		return p, err
	}
	dir, err := os.OpenRoot(root)
	if err != nil {
		return p, err
	}
	defer dir.Close()
	u := &UnchangedEvidence{Version: 1, InvocationID: a.Launch.InvocationID}
	u.ScopeSHA256, err = unchangedScope(ctx, b, s, a, dir)
	if err != nil {
		return p, err
	}
	base := filepath.Join(".parley-runtime", "invocations", u.InvocationID)
	var terminal, requested, started telemetry.Record
	u.TerminalSHA256, err = readReconciliationJSON(dir, filepath.Join(base, "terminal.json"), &terminal)
	if err != nil || terminal.SchemaVersion != telemetry.SchemaVersion || terminal.Type != "invocation.terminal" || terminal.InvocationID != u.InvocationID || terminal.Metadata.Idea != s.Policy.Idea || terminal.Metadata.Agent != a.Launch.Agent || terminal.Metadata.Phase != "fixup" || !runtimeID(terminal.Metadata.RunID) || terminal.RequestedAt.IsZero() || terminal.StartedAt == nil || terminal.CompletedAt == nil || terminal.PID == nil || *terminal.PID <= 0 || terminal.DurationMS == nil || *terminal.DurationMS < 0 || terminal.Outcome == nil || terminal.Outcome.ExitCode == nil || a.Terminal.ExitCode == nil {
		return p, errors.New("unchanged attempt lacks a complete observed process terminal")
	}
	o := terminal.Outcome
	// Missing starts, deadlines, cancellation, provider failures and signal exits
	// cannot establish a normally terminated unchanged attempt. Retain them for
	// explicit process recovery instead of inferring inactivity from silence.
	completed := o.Status == "process-exited" && o.FailureClass == nil && *o.ExitCode == 0
	failed := o.Status == "failed" && o.FailureClass != nil && *o.FailureClass == "process_failure" && *o.ExitCode > 0
	if (!completed && !failed) || o.Status != a.Terminal.Status || *o.ExitCode != *a.Terminal.ExitCode || terminal.RequestedAt.Before(s.Policy.StartedAt) || terminal.StartedAt.Before(terminal.RequestedAt) || terminal.StartedAt.Before(a.Charge.ReservedAt) || a.Terminal.At.Before(*terminal.StartedAt) || terminal.CompletedAt.Before(a.Terminal.At) {
		return p, errors.New("unchanged attempt does not have a matching normally exited process lifecycle")
	}
	u.RequestedSHA256, err = readReconciliationJSON(dir, filepath.Join(base, "requested.json"), &requested)
	if err != nil || requested.SchemaVersion != terminal.SchemaVersion || requested.Type != "invocation.requested" || requested.InvocationID != u.InvocationID || !sameJSON(requested.Metadata, terminal.Metadata) || !requested.RequestedAt.Equal(terminal.RequestedAt) || requested.StartedAt != nil || requested.CompletedAt != nil || requested.PID != nil || requested.Outcome != nil || requested.DurationMS != nil {
		return p, errors.New("unchanged attempt requested lifecycle changed")
	}
	u.StartedSHA256, err = readReconciliationJSON(dir, filepath.Join(base, "started.json"), &started)
	if err != nil || started.SchemaVersion != terminal.SchemaVersion || started.Type != "invocation.started" || started.InvocationID != u.InvocationID || !sameJSON(started.Metadata, terminal.Metadata) || !started.RequestedAt.Equal(terminal.RequestedAt) || started.StartedAt == nil || !started.StartedAt.Equal(*terminal.StartedAt) || started.PID == nil || *started.PID != *terminal.PID || started.CompletedAt != nil || started.Outcome != nil || started.DurationMS != nil {
		return p, errors.New("unchanged attempt started lifecycle changed")
	}
	u.PolicySHA256, err = s.Policy.SHA256()
	if err != nil {
		return p, err
	}
	attemptBytes, err := canonical(a)
	if err != nil {
		return p, err
	}
	u.AttemptSHA256 = digest(attemptBytes)
	stateBytes, err := canonical(s)
	if err != nil {
		return p, err
	}
	p = ReconciliationPreview{Version: 1, Root: root, StateSHA256: digest(stateBytes), Sequence: sequence, ChargeKey: a.Charge.EntryKey, Assessment: unchangedAssessment(s.Policy), Unchanged: u}
	return p, validateUnchangedPreview(s, p)
}

func PreviewUnchanged(ctx context.Context, root, idea string, sequence int) (ReconciliationPreview, error) {
	var p ReconciliationPreview
	root, err := canonicalRoot(root)
	if err != nil {
		return p, err
	}
	err = withState(ctx, root, idea, func(b budget.CycleBinding, _ budget.Snapshot, s State) error {
		var err error
		p, err = unchangedPreview(ctx, b, s, root, sequence)
		if err != nil {
			return err
		}
		if sequence <= len(s.Resolutions) {
			old := s.Resolutions[sequence-1].Preview
			if err = compareReconciledParent(p, old); err != nil {
				return err
			}
			p = old
		} else if sequence != len(s.Resolutions)+1 {
			return errors.New("unchanged reconciliation cannot omit an earlier charged attempt")
		}
		return nil
	})
	if err == nil && p.Version == 0 {
		err = errors.New("trajectory is not active")
	}
	return p, err
}

// ReconcileUnchanged publishes a deterministic observation under the existing
// cycle guard. It spends nothing, executes nothing and grants no continuation.
func ReconcileUnchanged(ctx context.Context, root, idea string, sequence int, expected string) (ReconciliationPreview, error) {
	return reconcileUnchanged(ctx, root, idea, sequence, expected, writeState)
}

func reconcileUnchanged(ctx context.Context, root, idea string, sequence int, expected string, persist func(string, State) error) (ReconciliationPreview, error) {
	var applied ReconciliationPreview
	root, err := canonicalRoot(root)
	if err != nil {
		return applied, err
	}
	err = withState(ctx, root, idea, func(b budget.CycleBinding, ledger budget.Snapshot, s State) error {
		p, err := unchangedPreview(ctx, b, s, root, sequence)
		if err != nil {
			return err
		}
		if sequence <= len(s.Resolutions) {
			old := s.Resolutions[sequence-1]
			if err = compareReconciledParent(p, old.Preview); err != nil {
				return err
			}
			if old.SHA256 != expected {
				return errors.New("unchanged reconciliation replay changed its original decision")
			}
			// Complete a failed durability barrier without changing history, even
			// if a later legitimate attempt or continuation is already retained.
			if err = syncUnchangedState(statePath(b)); err != nil {
				return err
			}
			applied = old.Preview
			return nil
		}
		if sequence != len(s.Resolutions)+1 || p.SHA256() != expected {
			return errors.New("trajectory changed since unchanged reconciliation preview")
		}
		s.Version = 3
		s.Resolutions = append(s.Resolutions, Resolution{p, expected, time.Now().UTC()})
		if err = validateState(s, b, ledger); err != nil {
			return err
		}
		if err = persist(statePath(b), s); err != nil {
			return err
		}
		applied = p
		return nil
	})
	if err == nil && applied.Version == 0 {
		err = errors.New("trajectory is not active")
	}
	return applied, err
}

func syncUnchangedState(path string) error {
	f, err := os.Open(path)
	if err != nil {
		return err
	}
	if err = errors.Join(fsutil.SyncFile(f), f.Close()); err != nil {
		return err
	}
	dir, err := os.Open(filepath.Dir(path))
	if err != nil {
		return err
	}
	return errors.Join(fsutil.SyncFile(dir), dir.Close())
}
