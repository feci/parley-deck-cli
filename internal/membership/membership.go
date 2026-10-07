package membership

import (
	"context"
	"crypto/sha256"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"strconv"
	"strings"
	"time"

	"parley-deck-cli/internal/protocol"
	"parley-deck-cli/internal/quota"
	"parley-deck-cli/internal/runmanifest"
	"parley-deck-cli/internal/store"
	"parley-deck-cli/internal/telemetry"
)

// Before is the only mutating recovery entry. A caller must hold the lifetime
// lease. Read surfaces use protocol.InspectQuota and never invoke this function.
func Before(ctx context.Context, root, ideaDir, runID string) (history *quota.History, resultErr error) {
	defer func() {
		if resultErr != nil && !IsBlocked(resultErr) {
			resultErr = IntegrityBlock(root, ideaDir, runID, history, resultErr)
		}
	}()

	v, err := protocol.InspectQuota(ideaDir)
	history = v.History
	if os.IsNotExist(err) && v.History == nil {
		return nil, nil
	}
	if err != nil {
		return v.History, err
	}
	h := v.History
	if len(v.Catchup) > 0 {
		return h, fmt.Errorf("%s", v.Pending)
	}
	if v.Manual != nil {
		if err := RecordManual(root, ideaDir); err != nil {
			return h, err
		}
		h, err = quota.ReadHistory(ideaDir)
		if err != nil {
			return h, err
		}
	}
	history = h
	if h == nil {
		return nil, nil
	}
	if !h.MidIdea() && h.Revision == 0 {
		return h, nil
	}
	if h.MidIdea() {
		if err = requireLease(ctx, ideaDir); err != nil {
			return nil, err
		}
	}
	if err = RequireStopped(root, h.Kickoff.Idea); err != nil {
		return h, err
	}
	if err = Reconcile(ctx, root, ideaDir, runID, h); err != nil {
		return h, err
	}
	return h, nil
}

func Reconcile(ctx context.Context, root, ideaDir, runID string, h *quota.History) error {
	release, err := ProjectionLock(ideaDir)
	if err != nil {
		return err
	}
	defer release()
	return reconcileLocked(ctx, root, ideaDir, runID, h)
}
func reconcileLocked(ctx context.Context, root, ideaDir, runID string, h *quota.History) error {
	if h.MidIdea() {
		if err := requireLease(ctx, ideaDir); err != nil {
			return err
		}
	}
	for i := range h.Batches {
		if err := quota.SyncPath(filepath.Join(ideaDir, quota.HistoryDir, fmt.Sprintf("%06d.json", i+1))); err != nil {
			return err
		}
	}
	manifests, err := validateManifests(root, runID, h)
	if err != nil {
		return err
	}

	if h.Revision > 0 {
		if err := projectionFault("prompt"); err != nil {
			return err
		}
	}
	if err := protocol.ReconcileQuotaPrompt(ideaDir, h); err != nil {
		return err
	}
	for id, m := range manifests {
		if m.QuotaRevision == h.Revision {
			if err := quota.SyncPath(runmanifest.Path(root, id)); err != nil {
				return err
			}
			continue
		}
		m.Participants = append([]string(nil), h.Current...)
		m.QuotaRevision = h.Revision
		// The roster snapshot is historical launch identity, not quorum; retain it.
		for i := range m.ActiveSteps {
			if !Has(h.Current, m.ActiveSteps[i].AgentID) {
				m.ActiveSteps[i].Status = "incomplete-excluded"
			}
		}
		m.NextActions = nil
		if err := projectionFault("manifest"); err != nil {
			return err
		}
		if err := runmanifest.Write(root, id, m); err != nil {
			return err
		}
	}
	for _, b := range h.Batches {
		receipt := filepath.Join(ideaDir, "quota-applied", b.ID)
		applied, err := quota.ReadApplied(ideaDir, b.ID)
		if err != nil {
			return err
		}
		if err := projectionFault("evaluation"); err != nil {
			return err
		}
		if err := evaluateRound(root, ideaDir, b); err != nil {
			return err
		}
		if err := publishNotice(root, ideaDir, b, applied); err != nil {
			return err
		}
		if err := projectionFault("applied"); err != nil {
			return err
		}
		if err := quota.DurableWrite(receipt, []byte(b.ID+"\n"), false); err != nil {
			return err
		}
	}
	return nil
}

func validateManifests(root, runID string, h *quota.History) (map[string]runmanifest.Manifest, error) {
	// Validate every touched projection before writing any of them. Missing original
	// run identity is an integrity gate, not an opportunity to repair legacy gap 11.
	runs := map[string]bool{runID: true, h.Kickoff.RunID: true}
	for _, b := range h.Batches {
		runs[b.RunID] = true
	}
	paths, err := filepath.Glob(filepath.Join(root, protocol.DeckDir, "runs", "*", "run.json"))
	if err != nil {
		return nil, err
	}
	for _, path := range paths {
		id := filepath.Base(filepath.Dir(path))
		m, e := runmanifest.Load(root, id)
		if e == nil && m.IdeaSlug == h.Kickoff.Idea {
			runs[id] = true
		}
	}
	manifests := map[string]runmanifest.Manifest{}
	for id := range runs {
		m, err := runmanifest.Load(root, id)
		if err != nil {
			return nil, fmt.Errorf("quota manifest %s: %w", id, err)
		}
		if m.IdeaSlug != h.Kickoff.Idea || m.RunID != id || m.QuotaKickoff == nil || !reflect.DeepEqual(m.QuotaKickoff, h.Kickoff) {
			return nil, fmt.Errorf("contradictory quota manifest %s", id)
		}
		if m.QuotaRevision < 0 || m.QuotaRevision > h.Revision {
			return nil, fmt.Errorf("contradictory quota manifest revision")
		}
		expected := h.Kickoff.Participants
		if m.QuotaRevision > 0 {
			expected = h.Batches[m.QuotaRevision-1].Decision.After
		}
		if !reflect.DeepEqual(m.Participants, expected) {
			return nil, fmt.Errorf("contradictory quota manifest membership")
		}
		manifests[id] = m
	}
	return manifests, nil
}

var projectionFault = func(stage string) error { return nil }

func ValidateRound(ideaDir, slug, label string, expected []string) error {
	review := strings.HasPrefix(filepath.ToSlash(label), "review/")
	n, err := strconv.Atoi(strings.TrimPrefix(filepath.Base(label), "round-"))
	if err != nil || n < 1 {
		return fmt.Errorf("invalid quota terminal round %s", label)
	}
	for _, id := range expected {
		path := filepath.Join(ideaDir, label, id+".md")
		if review {
			err = protocol.ValidateReviewArtifact(path, id, slug, n)
		} else {
			err = protocol.ValidateParticipantRoundArtifact(path, id, slug, n)
		}
		if err != nil {
			return fmt.Errorf("survivor artifact incomplete: %w", err)
		}
		if !review && n > 1 {
			meta, e := protocol.ReadFrontmatter(path)
			if e != nil {
				return e
			}
			if strings.Trim(meta["responding-to"], " []\"'") == "" {
				return fmt.Errorf("survivor cross-review missing responding-to")
			}
			raw, e := os.ReadFile(path)
			if e != nil {
				return e
			}
			for _, other := range expected {
				if other != id && !strings.Contains(string(raw), "### @"+other) {
					return fmt.Errorf("survivor cross-review missing @%s", other)
				}
			}
		}
	}
	if len(expected) == 0 {
		return fmt.Errorf("no surviving artifacts")
	}
	return nil
}
func evaluateRound(root, ideaDir string, b quota.Batch) error {
	if b.Round == "" || !strings.HasPrefix(filepath.Base(b.Round), "round-") {
		return nil
	}
	survivors := Intersect(b.Expected, b.Decision.After)
	if err := ValidateRound(ideaDir, b.Idea, b.Round, survivors); err != nil {
		return err
	}
	s := store.New(filepath.Join(root, protocol.DeckDir, "runs", b.RunID))
	evs, err := s.Load()
	if err != nil {
		return err
	}
	for _, e := range evs {
		if e.Type == "round.completed" && e.Data["quota_transition"] == b.ID {
			return s.Sync()
		}
	}
	return s.AppendDurable(store.Event{Type: "round.completed", Data: map[string]any{"idea": b.Idea, "round": b.Round, "participants": survivors, "completed": len(survivors), "total": len(survivors), "quota_transition": b.ID, "membership_revision": b.PriorRevision + 1, "reconciled": true}})
}

func Has(ids []string, id string) bool {
	for _, p := range ids {
		if p == id {
			return true
		}
	}
	return false
}
func Intersect(requested, current []string) []string {
	out := []string{}
	for _, id := range quota.FilterConfirmed(requested, nil) {
		if Has(current, id) {
			out = append(out, id)
		}
	}
	return out
}

// Settle receives execution facts only after ALL supervised writers have stopped.
// The launch layer validates success; content is never a source of candidacy.
func Settle(ctx context.Context, root, ideaDir, runID, round string, expected []string, members []quota.Member) (batch *quota.Batch, resultErr error) {
	var attempted *quota.Decision
	defer func() {
		if resultErr == nil || IsBlocked(resultErr) || attempted == nil {
			return
		}
		d := *attempted
		d.Block = "integrity/recovery gate: " + resultErr.Error()
		resultErr = Block(root, ideaDir, runID, round, d)
	}()

	h, err := Before(ctx, root, ideaDir, runID)
	if err != nil {
		return nil, err
	}
	if !h.MidIdea() {
		return nil, nil
	}
	roles, roleErr := Roles(ideaDir)
	d := quota.Evaluate(h.Policy(), h.Current, members, roles)
	if roleErr != nil && len(d.Candidates) > 0 {
		d.Applied = false
		d.Block = roleErr.Error()
		return nil, Block(root, ideaDir, runID, round, d)
	}
	if d.Block != "" {
		return nil, Block(root, ideaDir, runID, round, d)
	}
	if !d.Applied {
		return nil, nil
	}
	attempted = &d
	if err := CheckGates(root, ideaDir, runID, d.After); err != nil {
		d.Applied = false
		d.Block = err.Error()
		return nil, Block(root, ideaDir, runID, round, d)
	}
	if strings.HasPrefix(filepath.Base(round), "round-") {
		if err = ValidateRound(ideaDir, h.Kickoff.Idea, round, Intersect(expected, d.After)); err != nil {
			d.Applied = false
			d.Block = err.Error()
			return nil, Block(root, ideaDir, runID, round, d)
		}
	}
	release, lockErr := ProjectionLock(ideaDir)
	if lockErr != nil {
		return nil, lockErr
	}
	defer release()
	retained, err := protocol.CaptureQuotaObligations(ideaDir, quota.CandidateIDs(d.Candidates))
	if err != nil {
		return nil, err
	}
	if _, err = validateManifests(root, runID, h); err != nil {
		return nil, err
	}
	b := quota.NewBatch(h, runID, round, expected, d, time.Now()).BindRetained(retained)
	if err = quota.CommitBatch(ideaDir, b); err != nil {
		return nil, err
	}
	h, err = quota.ReadHistory(ideaDir)
	if err != nil {
		return &b, err
	}
	if err = reconcileLocked(ctx, root, ideaDir, runID, h); err != nil {
		return &b, err
	}
	return &b, nil
}

func Roles(ideaDir string) (quota.Roles, error) {
	role := protocol.ReadFacilitatorRole(ideaDir)
	r := quota.Roles{Facilitator: role.Facilitator}
	des := protocol.ReadImplementerDesignation(ideaDir)
	if des.State == protocol.DesignationSet {
		r.Designee = des.ID
	}
	if m, err := protocol.ReadFrontmatter(filepath.Join(ideaDir, "IMPLEMENTATION.md")); err == nil {
		r.PinnedImplementer = strings.Trim(m["implementer"], "\"'")
	} else if !os.IsNotExist(err) {
		return r, err
	}
	for _, name := range []string{"consensus.md", "FINAL.md", "review/consensus.md"} {
		path := filepath.Join(ideaDir, name)
		if _, err := os.Stat(path); os.IsNotExist(err) {
			continue
		} else if err != nil {
			return r, err
		}
		m, err := protocol.ReadFrontmatter(path)
		if err != nil {
			return r, err
		}
		id := strings.Trim(m["drafted-by"], "\"'")
		if id == "" {
			id = strings.Trim(m["author"], "\"'")
		}
		if id == "" {
			return r, fmt.Errorf("started canonical draft %s lacks attributable drafter", name)
		}
		r.StartedDrafters = append(r.StartedDrafters, id)
	}
	return r, nil
}
func Block(root, ideaDir, runID, phase string, d quota.Decision) error {
	blob, _ := json.Marshal(struct {
		Idea, Run, Phase string
		Decision         quota.Decision
	}{filepath.Base(ideaDir), runID, phase, d})
	keyPhase := phase
	if strings.HasPrefix(d.Block, "integrity/recovery gate:") {
		keyPhase = "quota-integrity"
	}
	key, _ := json.Marshal(struct {
		Idea, Run, Phase string
		Before           []string
		Candidates       []quota.Candidate
	}{filepath.Base(ideaDir), runID, keyPhase, d.Before, d.Candidates})
	id := fmt.Sprintf("quota-block-%x", sha256.Sum256(key))
	path := filepath.Join(root, protocol.DeckDir, "inbox", "parley-to-user_"+id+".md")
	if old, err := os.ReadFile(path); err == nil {
		if !strings.Contains(string(old), "transition: "+id+"\n") {
			return fmt.Errorf("contradictory quota escalation")
		}
		return &BlockedError{Reason: d.Block}
	} else if !os.IsNotExist(err) {
		return err
	}
	outcome := "Nothing applied."
	if strings.HasPrefix(d.Block, "integrity/recovery gate:") {
		outcome = "No further dispatch authorized; inspect immutable history for a committed pending transition."
	}
	arithmetic := fmt.Sprintf("%d distinct before; %d candidate identities; %d usable non-facilitator survivors; fixed floor %d", len(d.Before), len(quota.CandidateIDs(d.Candidates)), d.UsableSurvivors, quota.Floor)
	if len(d.Before) == 0 {
		arithmetic = fmt.Sprintf("unavailable because immutable membership could not be validated; fixed floor %d", quota.Floor)
	}
	text := fmt.Sprintf("---\nfrom: parley\nto: user\nidea: %s\nphase: %s\nblocking: yes\ntransition: %s\n---\n\nQuota batch blocked: %s\n\nCandidates: %v\nBefore: %v\nArithmetic: %s. %s\n\nEvidence:\n```json\n%s\n```\n", filepath.Base(ideaDir), phase, id, d.Block, quota.CandidateIDs(d.Candidates), d.Before, arithmetic, outcome, blob)
	if err := quota.DurableWrite(path, []byte(text), true); err != nil {
		return err
	}
	return &BlockedError{Reason: d.Block}
}

type BlockedError struct{ Reason string }

func (e *BlockedError) Error() string { return "quota batch blocked: " + e.Reason }
func IsBlocked(err error) bool        { var e *BlockedError; return errors.As(err, &e) }

// RequireStopped refuses unresolved started invocations from any run of this
// idea. A crashed driver's lease may disappear while a supervised child still
// writes; absence of terminal evidence is never treated as stopped.
func RequireStopped(root, idea string) error {
	paths, err := filepath.Glob(filepath.Join(root, ".parley-runtime", "invocations", "*", "started.json"))
	if err != nil {
		return err
	}
	for _, path := range paths {
		raw, err := os.ReadFile(path)
		if err != nil {
			return err
		}
		var started telemetry.Record
		if err = json.Unmarshal(raw, &started); err != nil {
			return fmt.Errorf("unreadable writer identity: %w", err)
		}
		if started.Metadata.Idea != idea {
			continue
		}
		terminal, err := os.ReadFile(filepath.Join(filepath.Dir(path), "terminal.json"))
		if os.IsNotExist(err) {
			if err = settleCrashedWriter(path, raw, started, idea); err != nil {
				return err
			}
			continue
		}
		if err != nil {
			return fmt.Errorf("unreadable terminal writer evidence: %w", err)
		}
		var result telemetry.Record
		if json.Unmarshal(terminal, &result) != nil || result.Type != "invocation.terminal" || result.InvocationID != started.InvocationID || result.InvocationID != filepath.Base(filepath.Dir(path)) || result.Metadata.Idea != idea || result.Outcome == nil || result.Outcome.Status == "" || result.CompletedAt == nil || result.StartedAt == nil || result.CompletedAt.Before(*result.StartedAt) {
			return fmt.Errorf("invalid terminal writer evidence")
		}
	}
	return nil
}

// MarkDraftStarted attributes a canonical draft before dispatch. An incomplete
// scaffold is preserved on failure; closed artifacts are never patched here.
func MarkDraftStarted(ideaDir, path, id string) error {
	release, err := ProjectionLock(ideaDir)
	if err != nil {
		return err
	}
	defer release()
	raw, err := os.ReadFile(path)
	if os.IsNotExist(err) {
		return quota.DurableWrite(path, []byte(fmt.Sprintf("---\nidea: %s\nstatus: draft\ndrafted-by: %s\n---\n", filepath.Base(ideaDir), id)), true)
	}
	if err != nil {
		return err
	}
	m, err := protocol.ReadFrontmatter(path)
	if err != nil {
		return err
	}
	if strings.Trim(m["drafted-by"], "\"'") != "" {
		return nil
	}
	if m["status"] == "final" || m["status"] == "complete" {
		return fmt.Errorf("closed canonical draft lacks attributed drafter")
	}
	if !strings.HasPrefix(string(raw), "---\n") {
		return fmt.Errorf("canonical draft lacks frontmatter")
	}
	out := strings.Replace(string(raw), "---\n", "---\ndrafted-by: "+id+"\n", 1)
	return quota.DurableWrite(path, []byte(out), false)
}

// IntegrityBlock is also used before a lease can be acquired when immutable
// authority itself is unreadable. It never guesses membership from exclusions.
func IntegrityBlock(root, ideaDir, run string, h *quota.History, cause error) error {
	scoped := h.MidIdea()
	if !scoped {
		if k, e := quota.ReadKickoff(ideaDir); e == nil && k != nil {
			scoped = k.Policy.Enabled && k.Policy.Scope == quota.KickoffAndMidIdea
		}
		if !scoped {
			if m, e := protocol.ReadFrontmatter(filepath.Join(ideaDir, "00-prompt.md")); e == nil {
				scoped = m["quota_auto_exclude"] == "true" && m["quota_auto_exclude_scope"] == quota.KickoffAndMidIdea
			}
		}
	}
	if !scoped {
		return cause
	}
	d := quota.Decision{Block: "integrity/recovery gate: " + cause.Error()}
	if h != nil {
		d.Before = append([]string(nil), h.Current...)
		d.After = append([]string(nil), h.Current...)
		if n := len(h.Batches); n > 0 {
			last := h.Batches[n-1]
			run = last.RunID
			d.Before = last.Decision.Before
			d.After = last.Decision.After
			d.Candidates = last.Decision.Candidates
			d.UsableSurvivors = last.Decision.UsableSurvivors
		}
	}
	return Block(root, ideaDir, run, "quota-recovery", d)
}
