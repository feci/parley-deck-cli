package app

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"

	"parley-deck-cli/internal/consensus"
	"parley-deck-cli/internal/driver"
	"parley-deck-cli/internal/protocol"
	"parley-deck-cli/internal/runner"
)

// driverImplOps is the production driver.ImplOps adapter (driver-impl-phase). It
// reuses the runner Phase 5-8 helpers + review-mode consensus; the driver depends
// only on the driver.ImplOps interface, so internal/driver never imports
// internal/app.
type driverImplOps struct {
	base        runner.Options // the round-01 runOpts (Root, RunID, Store, Agents, Idea, Timeout)
	root        string
	ideaSlug    string
	ideaDir     string
	implementer string   // FINAL drafter / first participant
	reviewers   []string // non-implementer participants
	drafter     string   // review-consensus drafter (facilitator)
	out         io.Writer
}

func newDriverImplOps(base runner.Options, root, ideaSlug, ideaDir string, participants []string, out io.Writer) driver.ImplOps {
	implementer := resolveImplementer(ideaDir, participants)
	var reviewers []string
	for _, p := range participants {
		if p != implementer {
			reviewers = append(reviewers, p)
		}
	}
	// The review-consensus drafter MUST be a non-implementer so the implementer
	// cannot filter reviewer findings out of the consensus (AF3). Fall back to the
	// implementer only if there are somehow no reviewers.
	drafter := implementer
	if len(reviewers) > 0 {
		drafter = reviewers[0]
	}
	return driverImplOps{
		base: base, root: root, ideaSlug: ideaSlug, ideaDir: ideaDir,
		implementer: implementer, reviewers: reviewers, drafter: drafter, out: out,
	}
}

// resolveImplementer picks the implementer from durable role metadata (D10/AF6):
// IMPLEMENTATION.md `implementer` (on re-entry), else FINAL.md `implementer` /
// `drafted-by`, validated against participants; otherwise participants[0].
func resolveImplementer(ideaDir string, participants []string) string {
	isParticipant := func(id string) bool {
		for _, p := range participants {
			if p == id {
				return true
			}
		}
		return false
	}
	for _, src := range []struct {
		file string
		keys []string
	}{
		{"IMPLEMENTATION.md", []string{"implementer"}},
		{"FINAL.md", []string{"implementer", "drafted-by"}},
	} {
		meta, err := protocol.ReadFrontmatter(filepath.Join(ideaDir, src.file))
		if err != nil {
			continue
		}
		for _, k := range src.keys {
			id := strings.Trim(strings.TrimSpace(meta[k]), `"'`)
			if id != "" && isParticipant(id) {
				return id
			}
		}
	}
	if len(participants) > 0 {
		return participants[0]
	}
	return ""
}

// withParticipants clones the base run options with a narrowed participant set.
func (o driverImplOps) withParticipants(ids ...string) runner.Options {
	opts := o.base
	opts.Idea.Participants = append([]string(nil), ids...)
	return opts
}

func (o driverImplOps) Implement(ctx context.Context) error {
	fmt.Fprintf(o.out, "driver: implementing via %s ...\n", o.implementer)
	r := runner.RunImplementation(ctx, o.withParticipants(o.implementer))
	if !r.Success() {
		return fmt.Errorf("implementer %s: %s", r.AgentID, r.ExitError)
	}
	return nil
}

func (o driverImplOps) ImplementationStatus() (string, error) {
	meta, err := protocol.ReadFrontmatter(filepath.Join(o.ideaDir, "IMPLEMENTATION.md"))
	if err != nil {
		return "", err
	}
	return strings.Trim(strings.TrimSpace(meta["status"]), `"'`), nil
}

// RunChecks runs `go test ./...` in the workspace when it is a Go module; a
// non-Go workspace has nothing to check (treated as passing).
func (o driverImplOps) RunChecks(ctx context.Context) (bool, string) {
	if _, err := os.Stat(filepath.Join(o.root, "go.mod")); err != nil {
		return true, "no go.mod in workspace; no checks to run"
	}
	fmt.Fprintf(o.out, "driver: running checks (go test ./...) ...\n")
	cmd := exec.CommandContext(ctx, "go", "test", "./...")
	cmd.Dir = o.root
	var buf bytes.Buffer
	cmd.Stdout = &buf
	cmd.Stderr = &buf
	err := cmd.Run()
	return err == nil, buf.String()
}

func (o driverImplOps) OpenReviewRound(ctx context.Context, round int) error {
	if len(o.reviewers) == 0 {
		return fmt.Errorf("no non-implementer reviewers available")
	}
	fmt.Fprintf(o.out, "driver: opening review round %d (reviewers: %s) ...\n", round, strings.Join(o.reviewers, ", "))
	// AF5: drop any reviewer artifact that exists but fails validation, so
	// RunReviewRound (Overwrite=false) regenerates it instead of skipping a
	// malformed file forever (which would spin the driver to the deadline).
	dir := filepath.Join(o.ideaDir, "review", roundDirLabel(round))
	for _, reviewer := range o.reviewers {
		path := filepath.Join(dir, reviewer+".md")
		if _, err := os.Stat(path); err == nil {
			if runner.ValidateReviewArtifact(path, reviewer, o.ideaSlug, round) != nil {
				_ = os.Remove(path)
			}
		}
	}
	opts := o.withParticipants(o.reviewers...)
	opts.Round = round
	results := runner.RunReviewRound(ctx, opts)
	failed := 0
	for _, r := range results {
		if r.ExitError != "" {
			failed++
		}
	}
	if failed == len(results) && failed > 0 {
		return fmt.Errorf("all reviewers failed in review round %d", round)
	}
	return nil
}

func (o driverImplOps) ReviewRoundComplete(round int) (bool, error) {
	dir := filepath.Join(o.ideaDir, "review", roundDirLabel(round))
	for _, reviewer := range o.reviewers {
		path := filepath.Join(dir, reviewer+".md")
		if _, err := os.Stat(path); err != nil {
			return false, nil
		}
		if err := runner.ValidateReviewArtifact(path, reviewer, o.ideaSlug, round); err != nil {
			return false, nil
		}
	}
	return len(o.reviewers) > 0, nil
}

func (o driverImplOps) DraftReviewConsensus(ctx context.Context, round int) error {
	fmt.Fprintf(o.out, "driver: drafting review consensus via %s ...\n", o.drafter)
	opts := o.withParticipants(o.drafter)
	opts.Round = round
	r := runner.RunReviewConsensus(ctx, opts)
	if !r.Success() {
		return fmt.Errorf("review-consensus drafter %s: %s", r.AgentID, r.ExitError)
	}
	path := filepath.Join(o.ideaDir, "review", "consensus.md")
	return runner.ValidateReviewConsensusArtifact(path)
}

func (o driverImplOps) ReviewStatus() (driver.ReviewStatus, error) {
	summary, err := consensus.Status(o.root, o.ideaSlug, true)
	if err != nil {
		return driver.ReviewStatus{}, err
	}
	meta, err := protocol.ReadFrontmatter(filepath.Join(o.ideaDir, "review", "consensus.md"))
	if err != nil {
		return driver.ReviewStatus{}, err
	}
	rawFixes := strings.Trim(strings.TrimSpace(meta["outstanding_agreed_fixes"]), `"'`)
	fixes, err := strconv.Atoi(rawFixes)
	if err != nil || fixes < 0 {
		return driver.ReviewStatus{}, fmt.Errorf("review/consensus.md outstanding_agreed_fixes=%q is not a non-negative integer", rawFixes)
	}
	blocked := strings.EqualFold(strings.Trim(strings.TrimSpace(meta["blocked"]), `"'`), "true")
	return driver.ReviewStatus{Summary: summary, OutstandingAgreedFixes: fixes, Blocked: blocked}, nil
}

func (o driverImplOps) RequestReviewSignoffs(ctx context.Context, missing []string) error {
	return requestConsensusSignoffs(ctx, requestSignoffsOptions{
		Root:            o.root,
		IdeaSlug:        o.ideaSlug,
		Review:          true,
		ParticipantsRaw: strings.Join(missing, ","),
		Yes:             true,
	}, o.out, o.out)
}

func (o driverImplOps) Fixup(ctx context.Context, cycle int) error {
	fmt.Fprintf(o.out, "driver: running fix-up cycle %d via %s ...\n", cycle, o.implementer)
	r := runner.RunFixup(ctx, o.withParticipants(o.implementer))
	if !r.Success() {
		return fmt.Errorf("fix-up implementer %s: %s", r.AgentID, r.ExitError)
	}
	return nil
}

// Complete writes IMPLEMENTATION.md status=complete. This is a deterministic file
// write by the orchestrator (NOT an implementer agent), so an implementer cannot
// short-circuit review (consensus D5).
func (o driverImplOps) Complete(ctx context.Context) error {
	path := filepath.Join(o.ideaDir, "IMPLEMENTATION.md")
	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	lines := strings.Split(string(data), "\n")
	inFrontmatter := false
	replaced := false
	for i, line := range lines {
		trimmed := strings.TrimSpace(line)
		if trimmed == "---" {
			if !inFrontmatter {
				inFrontmatter = true
				continue
			}
			if !replaced {
				lines = append(lines[:i], append([]string{"status: complete"}, lines[i:]...)...)
				replaced = true
			}
			break
		}
		if inFrontmatter && strings.HasPrefix(trimmed, "status:") {
			lines[i] = "status: complete"
			replaced = true
		}
	}
	if !replaced {
		return fmt.Errorf("%s has no frontmatter status field", path)
	}
	tmp := path + ".tmp"
	if err := os.WriteFile(tmp, []byte(strings.Join(lines, "\n")), 0o644); err != nil {
		return err
	}
	return os.Rename(tmp, path)
}

func roundDirLabel(n int) string { return fmt.Sprintf("round-%02d", n) }
