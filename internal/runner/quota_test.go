package runner

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"parley-deck-cli/internal/agents"
	"parley-deck-cli/internal/membership"
	"parley-deck-cli/internal/protocol"
	"parley-deck-cli/internal/quota"
	"parley-deck-cli/internal/runmanifest"
	"parley-deck-cli/internal/store"
)

func quotaRunnerFixture(t *testing.T) Options {
	t.Helper()
	root := t.TempDir()
	t.Setenv("PARLEY_HOME", t.TempDir())
	if err := protocol.InitWorkspace(root); err != nil {
		t.Fatal(err)
	}
	p := quota.NewPolicy(nil, nil)
	idea, k, err := protocol.CreateIdeaWithQuota(root, "quota runner", []string{"a", "b", "c"}, nil, "deliberation", "", "quota-test-run", &p, nil)
	if err != nil {
		t.Fatal(err)
	}
	m := runmanifest.New(runmanifest.Options{Root: root, RunID: k.RunID, IdeaSlug: idea.Slug, Participants: idea.Participants, QuotaKickoff: k})
	if err = runmanifest.Write(root, k.RunID, m); err != nil {
		t.Fatal(err)
	}
	sink := store.New(filepath.Join(root, protocol.DeckDir, "runs", k.RunID))
	if err = sink.AppendDurable(store.Event{Type: "run.created", Data: map[string]any{"idea": idea.Slug, "participants": idea.Participants, "quota_kickoff": k}}); err != nil {
		t.Fatal(err)
	}
	for _, id := range []string{"a", "b"} {
		raw := fmt.Sprintf("---\nagent: %s\nidea: %s\nround: 1\n---\n", id, idea.Slug)
		for _, h := range []string{"Summary", "Proposed approach", "Concerns / open questions", "Risks", "Existing alternatives"} {
			raw += "\n## " + h + "\nConcrete evidence.\n"
		}
		if err = os.WriteFile(filepath.Join(idea.Path, "round-01", id+".md"), []byte(raw), 0600); err != nil {
			t.Fatal(err)
		}
	}
	return Options{Root: root, RunID: k.RunID, Idea: idea, Store: sink, Round: 1, RoundLabel: "round-01", Timeout: 10 * time.Second}
}

func quotaRunnerResults(t *testing.T, opts Options) []Result {
	t.Helper()
	now := time.Now().UTC()
	reset := now.Add(2 * time.Hour)
	e := quota.Evidence{Eligible: true, InvocationID: "fixture-c", Adapter: "zcode", Provenance: "synthetic-test", RuleID: "test", ObservedAt: now, ResetAt: &reset, Excerpt: "weekly limit exhausted"}
	return []Result{{AgentID: "a", ArtifactOK: true}, {AgentID: "b", ArtifactOK: true}, {AgentID: "c", ExitError: "exit status 1", QuotaEvidence: &e}}
}

func TestQuotaRunnerRebindAndStaleLowLevelDispatchRefusal(t *testing.T) {
	opts := quotaRunnerFixture(t)
	ctx, live, release, err := quotaBefore(context.Background(), opts)
	if err != nil {
		t.Fatal(err)
	}
	defer release()
	results := quotaSettle(ctx, live, quotaRunnerResults(t, opts))
	if len(results) != 3 || !results[2].QuotaExcluded || results[2].ExitError == "" {
		t.Fatal(results)
	}
	_, again, done, err := quotaBefore(ctx, opts)
	if err != nil {
		t.Fatal(err)
	}
	defer done()
	if strings.Join(again.Idea.Participants, ",") != "a,b" {
		t.Fatal(again.Idea.Participants)
	}
	launch := WithLaunchInfo(ctx, LaunchInfo{RunID: opts.RunID, Idea: opts.Idea.Slug, Phase: "round-01"})
	if _, err = beginLaunch(launch, opts.Root, opts.RunID, agents.Discovery{Spec: agents.Spec{ID: "c"}}); err == nil || !strings.Contains(err.Error(), "excluded participant") {
		t.Fatal(err)
	}
	if _, _, err = membership.Acquire(context.Background(), opts.Idea.Path, "different-run"); err == nil {
		t.Fatal("competing driver admitted")
	}
	invocations, _ := filepath.Glob(filepath.Join(opts.Root, ".parley-runtime", "invocations", "*"))
	if len(invocations) != 0 {
		t.Fatal("stale dispatch created an invocation", invocations)
	}
}

func TestQuotaRunnerValidArtifactAndLaterSuccessCancelCandidate(t *testing.T) {
	for _, kind := range []string{"artifact", "later-success"} {
		t.Run(kind, func(t *testing.T) {
			opts := quotaRunnerFixture(t)
			ctx, live, release, err := quotaBefore(context.Background(), opts)
			if err != nil {
				t.Fatal(err)
			}
			defer release()
			results := quotaRunnerResults(t, opts)
			if kind == "artifact" {
				results[2].ArtifactOK = true
			} else {
				results = append(results, Result{AgentID: "c", ArtifactOK: true})
			}
			results = quotaSettle(ctx, live, results)
			h, err := quota.ReadHistory(opts.Idea.Path)
			if err != nil || h.Revision != 0 {
				t.Fatal(h, err)
			}
			for _, r := range results {
				if r.QuotaExcluded {
					t.Fatal(r)
				}
			}
		})
	}
}

// This uses only local shell fixtures. No provider or roster setting is touched.
func TestQuotaRunnerProcessFailureTransitionsAfterWritersStop(t *testing.T) {
	opts := quotaRunnerFixture(t)
	writeLaunchProtocol(t, opts.Root)
	reset := time.Now().UTC().Add(3 * time.Hour).Format(time.RFC3339Nano)
	stderr := "statusCode: 429\nresponseBody: '{\"error\":{\"message\":\"Weekly Limit Exhausted\",\"reset_at\":\"" + reset + "\"}}'\nError: Turn execution failed\n"
	capture := filepath.Join(opts.Root, "quota.stderr")
	os.WriteFile(capture, []byte(stderr), 0600)
	stub := filepath.Join(opts.Idea.Path, "round-01", "c.md")
	script := "cat > /dev/null\nprintf 'incomplete' > '" + strings.ReplaceAll(stub, "'", "'\\''") + "'\ncat '" + strings.ReplaceAll(capture, "'", "'\\''") + "' >&2\nexit 1"
	for _, id := range opts.Idea.Participants {
		opts.Agents = append(opts.Agents, agents.Discovery{Spec: agents.Spec{ID: id, AdapterID: "zcode", HeadlessArgs: []string{"-c", script}, PromptMode: agents.PromptStdin}, Found: true, Path: "/bin/sh"})
	}
	results := RunRoundOne(context.Background(), opts)
	h, err := quota.ReadHistory(opts.Idea.Path)
	if err != nil || h.Revision != 1 {
		for _, r := range results {
			t.Logf("agent=%s quota=%+v", r.AgentID, r.QuotaEvidence)
		}
		t.Fatalf("history=%+v err=%v results=%+v", h, err, results)
	}
	raw, _ := os.ReadFile(stub)
	if string(raw) != "incomplete" {
		t.Fatal("failed artifact lost", string(raw))
	}
	events, err := opts.Store.Load()
	if err != nil {
		t.Fatal(err)
	}
	failed, oldRound, newRound := false, false, false
	for _, e := range events {
		switch e.Type {
		case "agent.failed":
			failed = true
		case "round.incomplete":
			oldRound = true
		case "round.completed":
			if e.Data["quota_transition"] != nil {
				newRound = true
			}
		}
	}
	if !failed || !oldRound || !newRound {
		t.Fatal("terminal history missing", events)
	}
	for _, record := range terminalRecords(t, opts.Root) {
		if record.Metadata.Agent != "c" {
			t.Fatal("valid survivor relaunched", record.Metadata)
		}
	}
}

func TestQuotaRunnerLostRoundEventBlocksReduction(t *testing.T) {
	opts := quotaRunnerFixture(t)
	ctx, live, release, err := quotaBefore(context.Background(), opts)
	if err != nil {
		t.Fatal(err)
	}
	defer release()
	results := append(quotaRunnerResults(t, opts), Result{AgentID: "runner", ExitError: "round event append failed: injected"})
	results = quotaSettle(ctx, live, results)
	h, err := quota.ReadHistory(opts.Idea.Path)
	if err != nil || h.Revision != 0 || !results[len(results)-1].QuotaBlocked {
		t.Fatal(h, err, results)
	}
}
