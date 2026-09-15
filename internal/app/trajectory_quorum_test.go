package app

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"parley-deck-cli/internal/budget"
	"parley-deck-cli/internal/trajectory"
)

func TestTrajectoryLiveQuorumRefusesBeforeTicketAndLaunch(t *testing.T) {
	binary := trajectoryHelperBinary(t)
	for _, kind := range []string{"widened", "reordered", "yaml-comment"} {
		t.Run(kind, func(t *testing.T) {
			members := "[builder, reviewer]"
			if kind == "yaml-comment" {
				members = "[builder, reviewer, other] # original quorum"
			}
			root, trace, journal, agent := trajectoryHelperFixtureScope(t, "actual-helper", nil, members)
			ctx := context.Background()
			before, err := trajectory.Inspect(ctx, root, "idea-x")
			if err != nil {
				t.Fatal(err)
			}
			binding, err := budget.LoadCycleBinding(ctx, root, "idea-x", budget.Fixup)
			if err != nil {
				t.Fatal(err)
			}
			ledger, err := binding.Store.Inspect(ctx)
			if err != nil {
				t.Fatal(err)
			}
			prompt := filepath.Join(root, "parley-deck", "ideas", "idea-x", "00-prompt.md")
			original, err := os.ReadFile(prompt)
			if err != nil {
				t.Fatal(err)
			}
			changed := string(original)
			switch kind {
			case "widened":
				changed = strings.Replace(changed, members, "[builder, reviewer, other]", 1)
			case "reordered":
				changed = strings.Replace(changed, members, "[reviewer, builder]", 1)
			}
			if err = os.WriteFile(prompt, []byte(changed), 0600); err != nil {
				t.Fatal(err)
			}
			result, err := verifyTrajectoryWithAgent(ctx, root, "idea-x", agent, 30*time.Second, binary, io.Discard)
			if err == nil || !strings.Contains(err.Error(), "original activation quorum") || result.InvocationID != "" || result.RunID != "" {
				t.Fatal("live membership mismatch did not refuse before one-time authority", result, err)
			}
			for _, path := range []string{journal, trace, filepath.Join(root, ".parley-runtime", "verifier-starts")} {
				if _, err = os.Lstat(path); !os.IsNotExist(err) {
					t.Fatal("refused scope consumed a ticket or executed work", path, err)
				}
			}
			after, err := trajectory.Inspect(ctx, root, "idea-x")
			if err != nil {
				t.Fatal(err)
			}
			afterLedger, err := binding.Store.Inspect(ctx)
			if err != nil {
				t.Fatal(err)
			}
			oldState, _ := json.Marshal(before)
			newState, _ := json.Marshal(after)
			oldLedger, _ := json.Marshal(ledger)
			newLedger, _ := json.Marshal(afterLedger)
			if !bytes.Equal(oldState, newState) || !bytes.Equal(oldLedger, newLedger) {
				t.Fatal("early refusal rewrote original state or charges")
			}
			if kind == "widened" {
				if err = os.WriteFile(prompt, original, 0600); err != nil {
					t.Fatal(err)
				}
				result, err = verifyTrajectoryWithAgent(ctx, root, "idea-x", agent, 60*time.Second, binary, io.Discard)
				if err != nil || result.Assessment == nil || result.Assessment.Outcome != trajectory.Regression || result.InvocationID == "" {
					t.Fatal("exact restored membership lost the original helper opportunity", result, err)
				}
				actual, err := os.ReadFile(trace)
				if err != nil || string(actual) != "reviewer:original\nreviewer:broken\nreviewer:broken\nreviewer:original\n" {
					t.Fatal("restored helper did not execute the exact independent AB/BA checks", string(actual), err)
				}
			}
		})
	}
}

func TestTrajectoryHelperRejectsCoordinatedLiveQuorumBeforeClaim(t *testing.T) {
	root, trace, journal, _ := trajectoryHelperFixture(t, "actual-helper")
	ctx := context.Background()
	ticket, err := trajectory.PrepareCapturedVerification(ctx, root, "idea-x", "reviewer", "helper-scope")
	if err != nil {
		t.Fatal(err)
	}
	if err = trajectory.ReserveCapturedVerificationLaunch(ctx, ticket, "helper-scope-invocation"); err != nil {
		t.Fatal(err)
	}
	root = ticket.Root
	prompt := filepath.Join(root, "parley-deck", "ideas", "idea-x", "00-prompt.md")
	raw, err := os.ReadFile(prompt)
	if err != nil {
		t.Fatal(err)
	}
	if err = os.WriteFile(prompt, []byte(strings.Replace(string(raw), "[builder, reviewer]", "[builder, reviewer, other]", 1)), 0600); err != nil {
		t.Fatal(err)
	}
	participants, criteria, err := trajectoryMaterialScope(root, "idea-x", "builder", "reviewer")
	if err != nil {
		t.Fatal(err)
	}
	req := trajectoryHelperRequest{Version: 1, Ticket: ticket, Participants: participants, Criteria: criteria}
	base, err := trajectoryRuntime(root)
	if err != nil {
		t.Fatal(err)
	}
	dir := filepath.Join(base, ticket.RunID)
	if err = os.Mkdir(dir, 0700); err != nil {
		t.Fatal(err)
	}
	path := filepath.Join(dir, "request.json")
	if err = writeTrajectoryRuntimeJSON(path, req); err != nil {
		t.Fatal(err)
	}
	raw, err = os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	t.Setenv("PARLEY_RUN_ID", ticket.RunID)
	t.Setenv("PARLEY_AGENT_ID", "reviewer")
	t.Setenv("PARLEY_PROC_MARKER", "helper-scope-invocation")
	if err = executeTrajectoryHelper(ctx, path, sha256Hex(string(raw))); err == nil || !strings.Contains(err.Error(), "original activation quorum") {
		t.Fatal("coordinated live/request membership reached helper execution", err)
	}
	for _, name := range []string{"claim.json", "prepared.json", "process-001.json", "step-001.json", "receipt.json"} {
		if _, err = os.Lstat(filepath.Join(journal, name)); !os.IsNotExist(err) {
			t.Fatal("scope refusal consumed helper execution", name, err)
		}
	}
	if _, err = os.Lstat(trace); !os.IsNotExist(err) {
		t.Fatal("refused helper ran a criterion", err)
	}
}
