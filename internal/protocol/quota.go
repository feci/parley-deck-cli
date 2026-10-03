package protocol

import (
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"time"

	"parley-deck-cli/internal/quota"
)

// ReadQuotaState is read-only. Absence of both the kickoff record and policy
// fields means legacy confirmation. Partial or contradictory records fail closed.
func ReadQuotaState(ideaDir string) (*quota.Kickoff, error) {
	k, err := quota.ReadKickoff(ideaDir)
	if err != nil {
		return nil, err
	}
	meta, err := ReadFrontmatter(filepath.Join(ideaDir, "00-prompt.md"))
	if err != nil {
		return nil, err
	}
	// Duplicate authority keys are ambiguous even when the loose frontmatter reader would collapse them.
	raw, readErr := os.ReadFile(filepath.Join(ideaDir, "00-prompt.md"))
	if readErr != nil {
		return nil, readErr
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
	for _, key := range []string{"quota_auto_exclude", "quota_auto_exclude_scope", "participants"} {
		if fields[key] > 1 && (k != nil || fields["quota_auto_exclude"] > 0 || fields["quota_auto_exclude_scope"] > 0) {
			return nil, fmt.Errorf("ambiguous duplicate quota authority field %s", key)
		}
	}
	value, hasPolicy := meta["quota_auto_exclude"]
	scope, hasScope := meta["quota_auto_exclude_scope"]
	if k == nil {
		// An owner may explicitly disable a legacy idea; true cannot upgrade one.
		if hasPolicy && value != "false" {
			return nil, fmt.Errorf("quota policy has no immutable kickoff record")
		}
		if hasScope {
			return nil, fmt.Errorf("quota scope has no immutable kickoff record")
		}
		return nil, nil
	}
	if !hasPolicy || !hasScope || value != fmt.Sprint(k.Policy.Enabled) || scope != k.Policy.Scope || meta["idea"] != k.Idea || !reflect.DeepEqual(quota.Unique(parseList(meta["participants"])), quota.Unique(k.Participants)) {
		return nil, fmt.Errorf("pending or contradictory quota kickoff projections")
	}
	return k, nil
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
	err = f.Sync()
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
	err = dir.Sync()
	dir.Close()
	if err != nil {
		return idea, &k, err
	}
	return idea, &k, nil
}
