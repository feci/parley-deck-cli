package protocol

import (
	"errors"
	"fmt"
	"os"
	"parley-deck-cli/internal/fsutil"
	"path/filepath"
	"reflect"
	"strings"
	"time"

	"parley-deck-cli/internal/quota"
)

// ReadQuotaState is read-only. Absence of both the kickoff record and policy
// fields means legacy confirmation. Partial or contradictory records fail closed.
type QuotaView struct {
	History *quota.History `json:"history,omitempty"`
	Manual  *quota.Batch   `json:"-"`
	Pending string         `json:"pending,omitempty"`
	// Catchup is an uncommitted policy-off edit, distinct from a committed transition.
	Catchup []string `json:"-"`
}

// InspectQuota reads immutable authority and mutable projections without repairing
// anything. Policy-off prompt edits may preview Current; an incomplete catch-up
// never widens Known or commits history.
func InspectQuota(ideaDir string) (QuotaView, error) {
	h, err := quota.ReadHistory(ideaDir)
	if err != nil {
		return QuotaView{}, err
	}
	v := QuotaView{History: h}
	meta, err := ReadFrontmatter(filepath.Join(ideaDir, "00-prompt.md"))
	if err != nil {
		return v, err
	}
	raw, err := os.ReadFile(filepath.Join(ideaDir, "00-prompt.md"))
	if err != nil {
		return v, err
	}
	fields := map[string]int{}
	for i, line := range strings.Split(string(raw), "\n") {
		if i > 0 && strings.TrimSpace(line) == "---" {
			break
		}
		key, _, ok := strings.Cut(line, ":")
		if ok {
			fields[key]++
		}
	}
	for _, key := range []string{"quota_auto_exclude", "quota_auto_exclude_scope", "participants", "quota_revision"} {
		if fields[key] > 1 && (h != nil || fields["quota_auto_exclude"] > 0 || fields["quota_auto_exclude_scope"] > 0) {
			return v, fmt.Errorf("ambiguous duplicate quota authority field %s", key)
		}
	}
	value, hasPolicy := meta["quota_auto_exclude"]
	scope, hasScope := meta["quota_auto_exclude_scope"]
	if h == nil {
		if hasPolicy && value != "false" || hasScope || fields["quota_revision"] > 0 {
			return v, fmt.Errorf("quota policy has no immutable kickoff record")
		}
		return v, nil
	}
	k := h.Kickoff
	current := parseList(meta["participants"])
	matches := func(ids []string) bool {
		return reflect.DeepEqual(quota.Unique(current), quota.Unique(ids)) && len(current) == len(ids)
	}
	allApplied, latestApplied := true, true
	for i, b := range h.Batches {
		applied, e := quota.ReadApplied(ideaDir, b.ID)
		applied = applied && e == nil
		allApplied = allApplied && applied
		if i == len(h.Batches)-1 {
			latestApplied = applied
		}
		if !applied {
			v.Pending = "committed quota transition pending; checked recovery required: " + b.ID
		}
	}
	// Only the latest transition's durable receipt determines whether a stale
	// prompt can be interrupted projection work. Matching old content is no proof.
	before := k.Participants
	beforePolicy := k.Policy
	if len(h.Batches) > 0 {
		before = h.Batches[len(h.Batches)-1].Decision.Before
		for _, b := range h.Batches[:len(h.Batches)-1] {
			beforePolicy = b.Policy
		}
	}
	policyMatches := hasPolicy && hasScope && value == fmt.Sprint(h.Policy().Enabled) && scope == h.Policy().Scope
	priorPolicy := hasPolicy && hasScope && value == fmt.Sprint(beforePolicy.Enabled) && scope == beforePolicy.Scope
	if meta["idea"] != k.Idea || !policyMatches && (latestApplied || !priorPolicy) {
		return v, fmt.Errorf("contradictory immutable quota policy; restore the recorded prompt, then use parley quota revise for an owner-authorized change")
	}
	if !matches(h.Current) && !h.Policy().Enabled && policyMatches && allApplied {
		b, e := manualQuotaRevision(ideaDir, h, string(raw), current)
		if e != nil {
			var pending *manualCatchupPending
			if !errors.As(e, &pending) {
				return v, e
			}
			v.Catchup = pending.ids
			v.Pending = pending.Error()
			// Preview the requested quorum without adding any historical signer.
			// It stays incomplete until the late artifact validates and is imported.
			h.Current = append([]string(nil), current...)
			return v, nil
		}
		v.Manual = &b
		h.Current = append([]string(nil), current...)
		h.Known = quota.FilterConfirmed(append(h.Known, current...), nil)
	}
	if !matches(h.Current) && (latestApplied || !matches(before)) {
		return v, fmt.Errorf("contradictory quota membership projection; preserve the edit for the owner, restore recorded participants, then use parley quota revise; checked recovery only replays a pending transition's before/after set")
	}

	return v, nil
}

func ReadQuotaState(ideaDir string) (*quota.Kickoff, error) {
	v, err := InspectQuota(ideaDir)
	if err != nil {
		return nil, err
	}
	if v.Pending != "" {
		return nil, fmt.Errorf("%s", v.Pending)
	}
	if v.History == nil {
		return nil, nil
	}
	return v.History.Kickoff, nil
}

// QuotaMembers is the signoff/close read gate. It never repairs a pending commit.
func QuotaMembers(ideaDir string, fallback []string) (current, known []string, err error) {
	v, err := InspectQuota(ideaDir)
	if os.IsNotExist(err) && v.History == nil {
		return fallback, fallback, nil
	}
	if err != nil {
		return nil, nil, err
	}
	if v.Pending != "" && len(v.Catchup) == 0 {
		return nil, nil, fmt.Errorf("%s", v.Pending)
	}
	if v.History == nil {
		return fallback, fallback, nil
	}
	return v.History.Current, v.History.Known, nil
}

// ReconcileQuotaPrompt touches only membership and automatic display markers.
// It accepts an earlier recorded projection, never an unrelated manual edit.
func ReconcileQuotaPrompt(ideaDir string, h *quota.History) error {
	v, err := InspectQuota(ideaDir)
	if err != nil {
		return err
	}
	if v.History == nil || v.History.Revision != h.Revision {
		return fmt.Errorf("quota history changed")
	}
	path := filepath.Join(ideaDir, "00-prompt.md")
	raw, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	lines := strings.Split(string(raw), "\n")
	end := 0
	for i := 1; i < len(lines); i++ {
		if strings.TrimSpace(lines[i]) == "---" {
			end = i
			break
		}
		if strings.HasPrefix(lines[i], "quota_auto_exclude:") {
			lines[i] = fmt.Sprintf("quota_auto_exclude: %t", h.Policy().Enabled)
		}
		if strings.HasPrefix(lines[i], "quota_auto_exclude_scope:") {
			lines[i] = "quota_auto_exclude_scope: " + h.Policy().Scope
		}
		if strings.HasPrefix(lines[i], "participants:") {
			lines[i] = "participants: [" + strings.Join(h.Current, ", ") + "]"
		}
	}
	if end == 0 {
		return fmt.Errorf("missing prompt frontmatter")
	}
	var markers []string
	for _, b := range h.Batches {
		for _, m := range b.Markers() {
			line := "excluded: " + m
			if !strings.Contains(string(raw), line+"\n") {
				markers = append(markers, line)
			}
		}
	}
	lines = append(lines[:end], append(markers, lines[end:]...)...)
	out := strings.Join(lines, "\n")
	if out == string(raw) {
		return quota.SyncPath(path)
	}
	meta, err := ReadFrontmatter(path)
	if err != nil {
		return err
	}
	if meta["status"] == "final" || meta["status"] == "closed" {
		return fmt.Errorf("closed quota prompt is frozen")
	}
	return quota.DurableWrite(path, []byte(out), false)
}

// writeQuotaPrompt is used only while creating the initial prompt. It never
// derives membership from excluded lines or widens a saved policy on resume.
func writeQuotaPrompt(ideaDir string, k quota.Kickoff) error {
	path := filepath.Join(ideaDir, "00-prompt.md")
	b, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	fields := fmt.Sprintf("quota_auto_exclude: %t\nquota_auto_exclude_scope: %s\n", k.Policy.Enabled, k.Policy.Scope)
	for _, m := range k.Markers() {
		fields += "excluded: " + m + "\n"
	}
	text := strings.Replace(string(b), "\nstatus:", "\n"+fields+"status:", 1)
	return os.WriteFile(path, []byte(text), 0644)
}

// CreateIdeaWithQuota preserves the existing creation API for legacy callers.
// The first prompt's membership is already filtered before any bytes are written.
func CreateIdeaWithQuota(root, task string, participants, excluded []string, track, provenance, run string, p *quota.Policy, d *quota.Decision) (IdeaStatus, *quota.Kickoff, error) {
	if p == nil {
		idea, err := CreateIdeaFull(root, task, participants, excluded, track, provenance)
		return idea, nil, err
	}
	if err := p.Validate(); err != nil {
		return IdeaStatus{}, nil, err
	}
	if p.Enabled && p.Scope == quota.KickoffAndMidIdea {
		if _, err := quotaScopeRoot(root, ""); err != nil {
			return IdeaStatus{}, nil, err
		}
	}
	// Prepare in runtime storage, outside the visible idea inventory. Publish the
	// complete kickoff (evidence, frozen policy and filtered prompt) by one rename.
	stagingParent := filepath.Join(root, ".parley-runtime")
	if err := os.MkdirAll(stagingParent, 0700); err != nil {
		return IdeaStatus{}, nil, err
	}
	stage, err := os.MkdirTemp(stagingParent, "quota-kickoff-")
	if err != nil {
		return IdeaStatus{}, nil, err
	}
	defer os.RemoveAll(stage)
	idea, err := CreateIdeaFull(stage, task, participants, excluded, track, provenance)
	if err != nil {
		return idea, nil, err
	}
	finalSlug := uniqueSlug(filepath.Join(root, DeckDir, "ideas"), idea.Slug)
	promptPath := filepath.Join(idea.Path, "00-prompt.md")
	if finalSlug != idea.Slug {
		b, e := os.ReadFile(promptPath)
		if e != nil {
			return idea, nil, e
		}
		b = []byte(strings.Replace(string(b), "idea: "+idea.Slug+"\n", "idea: "+finalSlug+"\n", 1))
		if e = os.WriteFile(promptPath, b, 0644); e != nil {
			return idea, nil, e
		}
		idea.Slug = finalSlug
	}
	k := quota.NewKickoff(idea.Slug, run, *p, idea.Participants, d, time.Now().UTC())
	if err = quota.WriteKickoff(idea.Path, k); err != nil {
		return idea, nil, err
	}
	if err = writeQuotaPrompt(idea.Path, k); err != nil {
		return idea, &k, err
	}
	f, err := os.OpenFile(promptPath, os.O_RDWR, 0)
	if err != nil {
		return idea, &k, err
	}
	err = fsutil.SyncFile(f)
	closeErr := f.Close()
	if err != nil {
		return idea, &k, err
	}
	if closeErr != nil {
		return idea, &k, closeErr
	}
	parent := filepath.Join(root, DeckDir, "ideas")
	if err = os.MkdirAll(parent, 0755); err != nil {
		return idea, &k, err
	}
	target := filepath.Join(parent, idea.Slug)
	if err = os.Rename(idea.Path, target); err != nil {
		return idea, &k, err
	}
	idea.Path = target
	dir, err := os.Open(parent)
	if err != nil {
		return idea, &k, err
	}
	err = fsutil.SyncFile(dir)
	dir.Close()
	if err != nil {
		return idea, &k, err
	}
	return idea, &k, nil
}
