package app

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"parley-deck-cli/internal/agents"
	"parley-deck-cli/internal/driver"
	"parley-deck-cli/internal/membership"
	"parley-deck-cli/internal/pipeline"
	"parley-deck-cli/internal/protocol"
	"parley-deck-cli/internal/quota"
	"parley-deck-cli/internal/quotatest"
	"parley-deck-cli/internal/runmanifest"
	"parley-deck-cli/internal/runner"
	"parley-deck-cli/internal/store"
	"parley-deck-cli/internal/tui"
)

func fixupAppFixture(t *testing.T, p quota.Policy) (string, protocol.IdeaStatus, string) {
	t.Helper()
	root := t.TempDir()
	t.Setenv("PARLEY_HOME", t.TempDir())
	if e := protocol.InitWorkspace(root); e != nil {
		t.Fatal(e)
	}
	declareAppTestSource(t, root)
	ids := []string{"alpha", "beta", "gamma"}
	run := "recorded-run"
	idea, k, e := protocol.CreateIdeaWithQuota(root, "fixup app history", ids, nil, "deliberation", "", run, &p, nil)
	if e != nil {
		t.Fatal(e)
	}
	m := runmanifest.New(runmanifest.Options{Root: root, RunID: run, IdeaSlug: idea.Slug, Participants: ids, QuotaKickoff: k})
	if e = runmanifest.Write(root, run, m); e != nil {
		t.Fatal(e)
	}
	if e = store.New(filepath.Join(root, protocol.DeckDir, "runs", run)).AppendDurable(store.Event{Type: "run.created", Data: map[string]any{"idea": idea.Slug, "participants": ids, "quota_kickoff": k}}); e != nil {
		t.Fatal(e)
	}
	return root, idea, run
}
func TestQuotaFixupCreateRunOnFilesystem(t *testing.T) {
	for _, enabled := range []string{"true", "false"} {
		t.Run(enabled, func(t *testing.T) {
			root := t.TempDir()
			t.Setenv("PARLEY_HOME", t.TempDir())
			if e := protocol.InitWorkspace(root); e != nil {
				t.Fatal(e)
			}
			declareAppTestSource(t, root)
			bin := t.TempDir()
			writeFakeRoundAgentCLI(t, bin, "codex", "local-stub 1.0")
			writeAgentsLocalConfig(t, root, fakeAgentConfig{ID: "codex", Path: filepath.Join(bin, "codex")})
			fixupIsolateDiscovery(t, root)
			var out, errs bytes.Buffer
			args := []string{"run", "--dir", root, "--no-tui", "--no-auto", "--no-ping", "--no-preflight", "--yes", "--quota-auto-exclude=" + enabled, "--participants", "codex", "Shared filesystem acceptance"}
			code := Run(args, &out, &errs)
			t.Logf("command: parley %q\nexit=%d\nstdout:\n%sstderr:\n%s", args, code, out.String(), errs.String())
			if code != 0 {
				t.FailNow()
			}
			paths, e := filepath.Glob(filepath.Join(root, protocol.DeckDir, "ideas", "*"))
			ideas := []protocol.IdeaStatus{}
			for _, path := range paths {
				ideas = append(ideas, protocol.IdeaStatus{Path: path, Slug: filepath.Base(path)})
			}
			if e != nil || len(ideas) != 1 {
				t.Fatal(ideas, e)
			}
			h, e := quota.ReadHistory(ideas[0].Path)
			if e != nil || h.Policy().Enabled != (enabled == "true") {
				t.Fatal(h, e)
			}
			if e = protocol.ValidateParticipantRoundOneArtifact(filepath.Join(ideas[0].Path, "round-01/codex.md"), "codex", ideas[0].Slug); e != nil {
				t.Fatal(e)
			}
		})
	}
}
func TestQuotaFixupCLIRevisionRecovery(t *testing.T) {
	root, idea, run := fixupAppFixture(t, quota.NewPolicy(nil, nil))
	ids := []string{"alpha", "beta"}
	p := quota.NewPolicy(nil, nil)
	a := quotatest.Authority(t, root, idea.Slug, "cli-change", (quota.RevisionDirective{Participants: ids, Policy: p}).Text())
	req := map[string]any{"participants": ids, "policy": p, "authority": a}
	raw, _ := json.Marshal(req)
	path := filepath.Join(root, "request.json")
	os.WriteFile(path, raw, 0600)
	var out, errs bytes.Buffer
	args := []string{"quota", "revise", "--dir", root, "--idea", idea.Slug, "--run", run, "--request", path}
	if c := Run(args, &out, &errs); c != 0 {
		t.Fatal(c, out.String(), errs.String())
	}
	t.Logf("parley %q\n%s", args, out.String())
	h, e := quota.ReadHistory(idea.Path)
	if e != nil || h.Revision != 1 {
		t.Fatal(h, e)
	}
	receipt := filepath.Join(idea.Path, "quota-applied", h.Batches[0].ID)
	os.WriteFile(receipt, nil, 0600)
	if _, _, e = protocol.QuotaMembers(idea.Path, nil); e == nil {
		t.Fatal("truncated receipt admitted action")
	}
	for i := 0; i < 2; i++ {
		out.Reset()
		errs.Reset()
		args = []string{"quota", "recover", "--dir", root, "--idea", idea.Slug, "--run", run}
		if c := Run(args, &out, &errs); c != 0 {
			t.Fatal(c, errs.String())
		}
		t.Logf("parley %q\n%s", args, out.String())
	}
	for _, bad := range []string{strings.Replace(string(raw), `"participants":`, `"participants":["gamma"],"participants":`, 1), string(raw) + " {}"} {
		os.WriteFile(path, []byte(bad), 0600)
		out.Reset()
		errs.Reset()
		if c := Run([]string{"quota", "revise", "--dir", root, "--idea", idea.Slug, "--run", run, "--request", path}, &out, &errs); c == 0 {
			t.Fatal("ambiguous request accepted")
		}
	}
	if c := Run([]string{"quota", "recover", "--dir", root, "--idea", "..", "--run", run}, &out, &errs); c != 2 {
		t.Fatal("invalid slug", c)
	}
}
func TestQuotaFixupManualAndInteractiveHandoffKeepLease(t *testing.T) {
	root, idea, run := fixupAppFixture(t, quota.NewPolicy(nil, nil))
	prompt, _ := os.ReadFile(filepath.Join(idea.Path, "00-prompt.md"))
	writeConsensusIdea(t, root, idea.Slug, idea.Participants, false, map[string]string{"alpha": "accept"})
	os.WriteFile(filepath.Join(idea.Path, "00-prompt.md"), prompt, 0600)
	ctx, release, e := membership.Acquire(context.Background(), idea.Path, run)
	if e != nil {
		t.Fatal(e)
	}
	defer release()
	// Record a mid-idea owner revision first; a kickoff-only fixture would miss G5.
	ids := []string{"alpha", "beta"}
	p := quota.NewPolicy(nil, nil)
	a := quotatest.Authority(t, root, idea.Slug, "handoff-history", (quota.RevisionDirective{Participants: ids, Policy: p}).Text())
	if _, e = membership.Revise(ctx, root, idea.Path, run, ids, p, a, nil); e != nil {
		t.Fatal(e)
	}
	path := filepath.Join(idea.Path, "consensus.md")
	before, _ := os.ReadFile(path)
	for _, mode := range []string{agents.LaunchManual, agents.LaunchInteractive} {
		agent := agents.Discovery{Spec: agents.Spec{ID: "alpha", LaunchMode: mode}, Found: true}
		got, e := runSignoffAgent(ctx, root, run, agent, "Local handoff fixture", path, string(before), io.Discard, io.Discard)
		if e != nil {
			t.Fatalf("%s: %v", mode, e)
		}
		if mode == agents.LaunchManual && !got.Pending {
			t.Fatal("manual handoff not pending")
		}
		if _, r, e := membership.Acquire(context.Background(), idea.Path, "competing-run"); e == nil {
			r()
			t.Fatal("handoff dropped lease")
		}
		if _, e = runSignoffAgent(context.Background(), root, run, agent, "No owner context", path, string(before), io.Discard, io.Discard); e == nil {
			t.Fatal("unleased handoff accepted")
		}
		t.Logf("%s with history revision 1: handoff succeeded, competing run and unleased handoff refused", mode)
	}
}
func TestQuotaFixupTUIExitDrainsDriverBeforeRelease(t *testing.T) {
	root := t.TempDir()
	t.Setenv("PARLEY_HOME", t.TempDir())
	if e := protocol.InitWorkspace(root); e != nil {
		t.Fatal(e)
	}
	declareAppTestSource(t, root)
	bin := t.TempDir()
	writeFakeRoundAgentCLI(t, bin, "codex", "local-stub 1.0")
	writeAgentsLocalConfig(t, root, fakeAgentConfig{ID: "codex", Path: filepath.Join(bin, "codex")})
	fixupIsolateDiscovery(t, root)
	oldLive, oldDrive := runTaskLive, runTaskDrive
	defer func() { runTaskLive = oldLive; runTaskDrive = oldDrive }()
	started := make(chan struct{})
	cleanup := make(chan struct{})
	finish := make(chan struct{})
	ui := make(chan tui.LiveOptions, 1)
	runTaskDrive = func(ctx context.Context, d *driver.Driver) error {
		close(started)
		<-ctx.Done()
		close(cleanup)
		<-finish
		return ctx.Err()
	}
	runTaskLive = func(o tui.LiveOptions) error {
		ui <- o
		select {
		case <-started:
		case <-time.After(10 * time.Second):
			return fmt.Errorf("stub driver did not start")
		}
		return nil
	}
	done := make(chan int, 1)
	var out, errs bytes.Buffer
	go func() {
		done <- Run([]string{"run", "--dir", root, "--auto", "--no-ping", "--no-preflight", "--yes", "--participants", "codex", "TUI lifetime fixture"}, &out, &errs)
	}()
	var opts tui.LiveOptions
	select {
	case opts = <-ui:
	case <-time.After(10 * time.Second):
		t.Fatal("TUI never started")
	}
	select {
	case <-cleanup:
	case <-time.After(10 * time.Second):
		t.Fatal("driver not cancelled")
	}
	_, r, e := membership.Acquire(context.Background(), opts.Idea.Path, "during-cleanup")
	if e == nil {
		r()
		close(finish)
		t.Fatal("TUI released lease before driver cleanup")
	}
	select {
	case <-done:
		close(finish)
		t.Fatal("run returned while driver active")
	default:
	}
	close(finish)
	select {
	case c := <-done:
		if c != 0 {
			t.Fatal(c, out.String(), errs.String())
		}
	case <-time.After(10 * time.Second):
		t.Fatal("run never joined driver")
	}
	_, r, e = membership.Acquire(context.Background(), opts.Idea.Path, "after-cleanup")
	if e != nil {
		t.Fatal(e)
	}
	r()
	t.Log("UI returned; driver cancellation observed; competing acquisition refused until cleanup joined; acquired after run returned")
}
func TestQuotaFixupPipelineLifetimeAcrossStages(t *testing.T) {
	root, idea, _ := fixupAppFixture(t, quota.NewPolicy(nil, nil))
	ctx, release, e := quotaPipelineStart(context.Background(), root, idea.Path)
	if e != nil {
		t.Fatal(e)
	}
	run := membership.DrivingRunID(ctx)
	if run == "" {
		t.Fatal("missing pipeline owner")
	}
	for _, stage := range []string{"round", "implementation", "review", "signoff", "fixup"} {
		if quotaPipelineRun(ctx, stage) != run {
			t.Fatal("pipeline stage changed identity")
		}
		nested, r, e := membership.Acquire(ctx, idea.Path, run)
		if e != nil {
			t.Fatal(e)
		}
		if e = membership.CheckLease(nested, idea.Path, run); e != nil {
			t.Fatal(e)
		}
		r()
		if _, r, e := membership.Acquire(context.Background(), idea.Path, "between-"+stage); e == nil {
			r()
			t.Fatal("between-stage ownership lost")
		}
	}
	release()
	if _, r, e := membership.Acquire(context.Background(), idea.Path, "after-pipeline"); e != nil {
		t.Fatal(e)
	} else {
		r()
	}
	t.Log("all pipeline stages use one live idea-bound run; between-stage competitors refused; release permits next driver")
}
func TestQuotaFixupPreflightBare503AndNoise(t *testing.T) {
	for _, s := range []string{"503", " 503 \n", "503 Service Unavailable", "HTTP/2 503", "statusCode: 503", "API error 503"} {
		if got := providerFailureClass(s); got != "overloaded" {
			t.Fatalf("%q: %q", s, got)
		}
	}
	for _, s := range []string{"size=503 bytes", "at source.js:503:1", "processed 503 records", "HTTP/2 200 (503 bytes)", "line 503", "build 1503"} {
		if got := providerFailureClass(s); got != "" {
			t.Fatalf("noise %q: %q", s, got)
		}
	}
}

func fixupIsolateDiscovery(t *testing.T, root string) {
	t.Helper()
	path := filepath.Join(root, protocol.DeckDir, "agents.local.toml")
	raw, _ := os.ReadFile(path)
	f, e := os.OpenFile(path, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0600)
	if e != nil {
		t.Fatal(e)
	}
	defer f.Close()
	for _, s := range agents.DefaultSpecs() {
		if !strings.Contains(string(raw), "[agents."+s.ID+"]") {
			fmt.Fprintf(f, "\n[agents.%s]\ncommand = %q\n", s.ID, filepath.Join(root, "absent-local-test-agent"))
		}
	}
}
func TestQuotaFixupOffModeManualExclusionThroughSignoffCLI(t *testing.T) {
	root, idea, _ := fixupAppFixture(t, quota.Policy{Enabled: false, Scope: quota.KickoffAndMidIdea})
	path := filepath.Join(idea.Path, "00-prompt.md")
	raw, _ := os.ReadFile(path)
	writeConsensusIdea(t, root, idea.Slug, idea.Participants, false, nil)
	raw = bytes.Replace(raw, []byte("participants: [alpha, beta, gamma]"), []byte("participants: [alpha, beta]\nexcluded: [gamma — unavailable — confirmed 2026-10-04]"), 1)
	os.WriteFile(path, raw, 0600)
	for _, id := range []string{"alpha", "beta"} {
		var out, errs bytes.Buffer
		args := []string{"consensus", "signoff", "--dir", root, "--agent", id, "--status", "accept", idea.Slug}
		if c := Run(args, &out, &errs); c != 0 {
			t.Fatal(c, out.String(), errs.String())
		}
		t.Logf("parley %q\n%s", args, out.String())
	}
	h, e := quota.ReadHistory(idea.Path)
	if e != nil || h.Revision != 1 || strings.Join(h.Current, ",") != "alpha,beta" {
		t.Fatal(h, e)
	}
}
func TestQuotaFixupAgentsExecCanonicalTargetWithoutIdea(t *testing.T) {
	root, idea, _ := fixupAppFixture(t, quota.Policy{Enabled: false, Scope: quota.KickoffAndMidIdea})
	path := filepath.Join(idea.Path, "00-prompt.md")
	raw, _ := os.ReadFile(path)
	raw = bytes.Replace(raw, []byte("participants: [alpha, beta, gamma]"), []byte("participants: [alpha, beta]\nexcluded: [gamma — unavailable — confirmed 2026-10-04]"), 1)
	os.WriteFile(path, raw, 0600)
	if e := membership.RecordManual(root, idea.Path); e != nil {
		t.Fatal(e)
	}
	bin := t.TempDir()
	script := filepath.Join(bin, "gamma")
	marker := filepath.Join(root, "dispatched")
	os.WriteFile(script, []byte("#!/bin/sh\nif [ \"$1\" = \"--version\" ]; then echo local-stub; exit 0; fi\ncat >/dev/null\nprintf called > '"+strings.ReplaceAll(marker, "'", "'\\''")+"'\n"), 0700)
	writeAgentsLocalConfig(t, root, fakeAgentConfig{ID: "gamma", Path: script})
	fixupIsolateDiscovery(t, root)
	prompt := filepath.Join(root, "prompt.txt")
	os.WriteFile(prompt, []byte("local synthetic work"), 0600)
	var out, errs bytes.Buffer
	args := []string{"agents", "exec", "--dir", root, "--agent", "gamma", "--prompt-file", prompt, "--artifact", filepath.Join(idea.Path, "round-02/gamma.md"), "--yes"}
	if c := Run(args, &out, &errs); c == 0 {
		t.Fatal("excluded canonical launch accepted")
	}
	if _, e := os.Stat(marker); !os.IsNotExist(e) {
		t.Fatal("excluded process launched")
	}
	t.Logf("parley %q\n%s", args, errs.String())
	out.Reset()
	errs.Reset()
	args = []string{"agents", "exec", "--dir", root, "--agent", "gamma", "--prompt-file", prompt, "--yes"}
	if c := Run(args, &out, &errs); c != 0 {
		t.Fatal("unbound work rejected", c, errs.String())
	}
	if _, e := os.Stat(marker); e != nil {
		t.Fatal(e)
	}
	t.Logf("unbound generic exec allowed: %s", out.String())
}

func TestQuotaFixupPipelineBlockDriverOwnsAllStages(t *testing.T) {
	root := t.TempDir()
	t.Setenv("PARLEY_HOME", t.TempDir())
	if e := protocol.InitWorkspace(root); e != nil {
		t.Fatal(e)
	}
	declareAppTestSource(t, root)
	slug, blockID := "pipeline-fixture", "build"
	ideaSlug := slug + "__" + blockID
	deck := filepath.Join(root, protocol.DeckDir)
	dir := filepath.Join(deck, "ideas", ideaSlug)
	os.MkdirAll(dir, 0700)
	p := quota.NewPolicy(nil, nil)
	ids := []string{"alpha", "beta", "gamma"}
	k := quota.NewKickoff(ideaSlug, "kickoff", p, ids, nil, time.Now())
	if e := quota.WriteKickoff(dir, k); e != nil {
		t.Fatal(e)
	}
	os.WriteFile(filepath.Join(dir, "00-prompt.md"), []byte(fmt.Sprintf("---\nidea: %s\nstatus: round-01\nparticipants: [alpha, beta, gamma]\nquota_auto_exclude: true\nquota_auto_exclude_scope: kickoff-and-mid-idea\n---\n", ideaSlug)), 0600)
	m := runmanifest.New(runmanifest.Options{Root: root, RunID: "kickoff", IdeaSlug: ideaSlug, Participants: ids, QuotaKickoff: &k})
	if e := runmanifest.Write(root, "kickoff", m); e != nil {
		t.Fatal(e)
	}
	store.New(filepath.Join(deck, "runs", "kickoff")).AppendDurable(store.Event{Type: "run.created", Data: map[string]any{"idea": ideaSlug, "participants": ids, "quota_kickoff": k}})
	bin := t.TempDir()
	stub := filepath.Join(bin, "version-only")
	os.WriteFile(stub, []byte("#!/bin/sh\nif [ \"$1\" = \"--version\" ]; then echo stub; exit 0; fi\nexit 99\n"), 0700)
	writeAgentsLocalConfig(t, root, fakeAgentConfig{ID: "alpha", Path: stub}, fakeAgentConfig{ID: "beta", Path: stub}, fakeAgentConfig{ID: "gamma", Path: stub})
	fixupIsolateDiscovery(t, root)
	oldI, oldR, oldC, oldF := pipelineRunImplementation, pipelineRunReviewRound, pipelineRunReviewConsensus, pipelineRunFixup
	defer func() {
		pipelineRunImplementation = oldI
		pipelineRunReviewRound = oldR
		pipelineRunReviewConsensus = oldC
		pipelineRunFixup = oldF
	}()
	run := ""
	var stages []string
	check := func(ctx context.Context, o runner.Options, stage string) {
		t.Helper()
		if run == "" {
			run = o.RunID
		}
		if o.RunID != run || membership.DrivingRunID(ctx) != run {
			t.Fatal("stage identity changed", stage, o.RunID, run)
		}
		if e := membership.CheckLease(ctx, dir, run); e != nil {
			t.Fatal(e)
		}
		if _, r, e := membership.Acquire(context.Background(), dir, "competitor-"+stage); e == nil {
			r()
			t.Fatal("pipeline stage not exclusively owned", stage)
		}
		stages = append(stages, stage)
	}
	pipelineRunImplementation = func(ctx context.Context, o runner.Options) runner.Result {
		check(ctx, o, "implementation")
		return runner.Result{AgentID: "alpha", ArtifactOK: true}
	}
	pipelineRunReviewRound = func(ctx context.Context, o runner.Options) []runner.Result {
		check(ctx, o, fmt.Sprintf("review-%d", o.Round))
		return []runner.Result{{AgentID: "beta", ArtifactOK: true}, {AgentID: "gamma", ArtifactOK: true}}
	}
	pipelineRunReviewConsensus = func(ctx context.Context, o runner.Options) runner.Result {
		check(ctx, o, fmt.Sprintf("consensus-%d", o.Round))
		os.MkdirAll(filepath.Join(dir, "review"), 0700)
		count := 1
		if o.Round == 2 {
			count = 0
		}
		os.WriteFile(filepath.Join(dir, "review/consensus.md"), []byte(fmt.Sprintf("---\noutstanding_agreed_fixes: %d\n---\n", count)), 0600)
		return runner.Result{AgentID: "alpha", ArtifactOK: true}
	}
	pipelineRunFixup = func(ctx context.Context, o runner.Options) runner.Result {
		check(ctx, o, "fixup")
		return runner.Result{AgentID: "alpha", ArtifactOK: true}
	}
	var out, errs bytes.Buffer
	code := autoDriveImplementationBlock(context.Background(), root, deck, slug, pipeline.Block{ID: blockID}, strings.Join(ids, ","), "alpha", true, &out, &errs)
	if code != 0 {
		t.Fatal(code, out.String(), errs.String())
	}
	if strings.Join(stages, ",") != "implementation,review-1,consensus-1,fixup,review-2,consensus-2" {
		t.Fatal(stages, out.String())
	}
	if _, r, e := membership.Acquire(context.Background(), dir, "after-return"); e != nil {
		t.Fatal(e)
	} else {
		r()
	}
	t.Logf("actual block driver, locally stubbed stage dispatch; stages=%v run=%s\n%s", stages, run, out.String())
}
