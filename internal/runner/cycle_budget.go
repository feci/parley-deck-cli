package runner

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v3"
	"parley-deck-cli/internal/budget"
	"parley-deck-cli/internal/track"
)

// groupProtocolCycle freezes/loads a policy before a grouped operation starts.
// A failure is deferred to the common instrumented launch boundary, so it
// cannot disappear behind a high-level early return. Saved limits are reused.
func groupProtocolCycle(ctx context.Context, root, idea, ideaDir, runID string, kind budget.Kind) (context.Context, func()) {
	noop := func() {}
	fixups, cross, explicit, err := cycleCeilings(ideaDir)
	cap := fixups
	if kind == budget.CrossReview {
		cap = cross
	}
	var b *budget.CycleBinding
	if err == nil {
		b, err = budget.LoadCycleBinding(ctx, root, idea, kind)
	}
	if err == nil {
		floor, e := budget.LegacyCycleFloor(ideaDir, kind)
		err = e
		if err == nil && b == nil {
			b, err = budget.EnsureCycleBinding(ctx, root, idea, kind, cap, floor, filepath.Join(root, "parley-deck", "runs", runID), ideaDir)
		} else if err == nil {
			state, e := b.Store.Inspect(ctx)
			err = e
			if err == nil && floor > b.Count(state) && !budget.CycleSessionMatches(ctx, b) {
				err = errors.New("visible cycle history exceeds the shared charged count")
			}
			if err == nil && explicit && cap < b.Policy.InitialMaximum() {
				err = errors.New("current track is stricter than the frozen cycle policy; reconcile before continuing")
			}
		}
	}
	if err == nil {
		// A disposable clone gets the private live origin before reaching here.
		rel, e := filepath.Rel(root, ideaDir)
		if e != nil || filepath.ToSlash(rel) != b.Policy.IdeaPath {
			err = errors.New("runner idea path differs from the frozen cycle scope")
		}
	}
	if err != nil {
		return budget.DeferCycleRefusal(ctx, kind, err), noop
	}
	scoped, finish, err := budget.OpenCycleSession(ctx, b)
	if err != nil {
		return budget.DeferCycleRefusal(ctx, kind, err), noop
	}
	return scoped, finish
}

// Standalone typed CLI launches use the same boundary as grouped protocol
// runners. No natural-language prompt is inspected to guess an operation kind.
func prepareLaunchCycle(ctx context.Context, root, idea, phase, runID string) (context.Context, func()) {
	kind, ok := budget.CycleKindForPhase(phase)
	if !ok {
		return ctx, func() {}
	}
	ideaDir := filepath.Join(root, "parley-deck", "ideas", idea)
	if b, err := budget.LoadCycleBinding(ctx, root, idea, kind); err == nil && b != nil {
		ideaDir = filepath.Join(root, filepath.FromSlash(b.Policy.IdeaPath))
	}
	return groupProtocolCycle(ctx, root, idea, ideaDir, runID, kind)
}

// Preserve the existing driver defaults and track cells. Parse only the actual
// frontmatter document, rejecting ambiguous YAML keys/types instead of treating
// a malformed track as a larger legacy grant. Pipeline blocks also have seeded
// idea prompts; missing authority is never a reason to freeze a legacy grant.
func cycleCeilings(ideaDir string) (int, int, bool, error) {
	const legacyFixups, legacyCross = 3, 4
	path := filepath.Join(ideaDir, "00-prompt.md")
	info, err := os.Lstat(path)
	if os.IsNotExist(err) {
		return 0, 0, false, errors.New("cycle policy prompt is missing")
	}
	if err != nil {
		return 0, 0, false, err
	}
	if !info.Mode().IsRegular() || info.Size() > 1<<20 {
		return 0, 0, false, errors.New("cycle policy prompt must be a bounded regular file")
	}
	f, err := os.Open(path)
	if err != nil {
		return 0, 0, false, err
	}
	defer f.Close()
	opened, err := f.Stat()
	if err != nil || !os.SameFile(info, opened) {
		return 0, 0, false, errors.New("cycle policy prompt changed during open")
	}
	data, err := io.ReadAll(io.LimitReader(f, (1<<20)+1))
	if err != nil {
		return 0, 0, false, err
	}
	if len(data) > 1<<20 {
		return 0, 0, false, errors.New("cycle policy prompt exceeds limit")
	}
	lines := strings.Split(strings.ReplaceAll(string(data), "\r\n", "\n"), "\n")
	if len(lines) == 0 || lines[0] != "---" {
		return 0, 0, false, errors.New("cycle policy prompt lacks frontmatter")
	}
	end := 0
	for i := 1; i < len(lines); i++ {
		if lines[i] == "---" {
			end = i
			break
		}
	}
	if end == 0 {
		return 0, 0, false, errors.New("cycle policy frontmatter is unterminated")
	}
	var doc yaml.Node
	decoder := yaml.NewDecoder(bytes.NewBufferString(strings.Join(lines[1:end], "\n")))
	if err := decoder.Decode(&doc); err != nil {
		return 0, 0, false, err
	}
	if len(doc.Content) != 1 || doc.Content[0].Kind != yaml.MappingNode {
		return 0, 0, false, errors.New("cycle policy frontmatter must be a mapping")
	}
	fields := map[string]*yaml.Node{}
	for i := 0; i < len(doc.Content[0].Content); i += 2 {
		k, v := doc.Content[0].Content[i], doc.Content[0].Content[i+1]
		if k.Kind != yaml.ScalarNode || k.Tag != "!!str" || k.Value != strings.ToLower(k.Value) || fields[k.Value] != nil || v.Kind == yaml.AliasNode {
			return 0, 0, false, errors.New("ambiguous cycle policy frontmatter")
		}
		fields[k.Value] = v
	}
	raw := ""
	if n := fields["track"]; n != nil {
		if n.Kind != yaml.ScalarNode || n.Tag != "!!str" || strings.TrimSpace(n.Value) == "" {
			return 0, 0, false, errors.New("invalid cycle policy track")
		}
		raw = n.Value
	}
	t, present, err := track.NormalizeStrict(raw)
	if err != nil {
		return 0, 0, false, err
	}
	boolField := func(name string) (bool, error) {
		n := fields[name]
		if n == nil {
			return false, nil
		}
		if n.Kind != yaml.ScalarNode || n.Tag != "!!bool" {
			return false, fmt.Errorf("invalid %s in cycle policy", name)
		}
		return strings.EqualFold(n.Value, "true"), nil
	}
	auto, err := boolField("auto_implement")
	if err != nil {
		return 0, 0, false, err
	}
	strict, err := boolField("strict_gate")
	if err != nil {
		return 0, 0, false, err
	}
	participants := map[string]bool{}
	if n := fields["participants"]; n != nil {
		if n.Kind != yaml.SequenceNode {
			return 0, 0, false, errors.New("cycle policy participants must be a sequence")
		}
		for _, v := range n.Content {
			if v.Kind != yaml.ScalarNode || v.Tag != "!!str" || v.Value == "" {
				return 0, 0, false, errors.New("invalid cycle policy participant")
			}
			participants[v.Value] = true
		}
	}
	policy, err := track.PolicyFor(t, present, len(participants)-1, auto, strict)
	if err != nil {
		return 0, 0, false, err
	}
	fixups, cross := legacyFixups, legacyCross
	if policy.ApplyOverrides {
		if policy.MaxFixupCycles > 0 {
			fixups = policy.MaxFixupCycles
		}
		if policy.CapCrossReviewRounds > 0 {
			cross = policy.CapCrossReviewRounds
		}
		if policy.CrossReviewRounds == 0 {
			cross = 0
		}
	}
	return fixups, cross, present, nil
}
