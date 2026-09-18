package app

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"testing"

	"parley-deck-cli/internal/agents"
	"parley-deck-cli/internal/budget"
	"parley-deck-cli/internal/config"
	"parley-deck-cli/internal/protocol"
	"parley-deck-cli/internal/protocolcore"
	"parley-deck-cli/internal/protocolpacket"
	"parley-deck-cli/internal/runner"
	"parley-deck-cli/internal/store"
	"parley-deck-cli/internal/telemetry"
	"parley-deck-cli/internal/trajectory"
)

func refuseAppProtocolCache(t *testing.T, root string, phase int) func() {
	t.Helper()
	prepared, err := protocolpacket.Render(root, protocolcore.StoreAt(config.CentralHome()), protocolpacket.Request{Phase: phase, Track: "deliberation", IdeaSlug: "idea-x"})
	if err != nil || prepared.BodyPath == "" {
		t.Fatal("cannot prepare exact phase context", err)
	}
	path := prepared.BodyPath
	raw, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	if err = os.WriteFile(path, []byte("tampered cached context"), 0600); err != nil {
		t.Fatal(err)
	}
	return func() {
		if err := os.WriteFile(path, raw, 0600); err != nil {
			t.Fatal(err)
		}
	}
}

func TestTrajectoryCLIProtocolPrecheckPreservesTicketOpportunity(t *testing.T) {
	binary := trajectoryHelperBinary(t)
	root, trace, journal, _ := trajectoryHelperFixture(t, "actual-helper")
	ctx := context.Background()
	before, err := trajectory.Inspect(ctx, root, "idea-x")
	if err != nil {
		t.Fatal(err)
	}
	b, err := budget.LoadCycleBinding(ctx, root, "idea-x", budget.Fixup)
	if err != nil {
		t.Fatal(err)
	}
	ledger, err := b.Store.Inspect(ctx)
	if err != nil {
		t.Fatal(err)
	}
	restore := refuseAppProtocolCache(t, root, -1)
	run := func() (trajectoryVerificationResult, error) {
		t.Helper()
		cmd := exec.Command(binary, "trajectory", "verify", "--dir", root, "--idea", "idea-x", "--verifier", "reviewer", "--timeout", "90s", "--yes")
		var out, errout bytes.Buffer
		cmd.Stdout, cmd.Stderr = &out, &errout
		runErr := cmd.Run()
		var result trajectoryVerificationResult
		if err := json.Unmarshal(out.Bytes(), &result); err != nil {
			t.Fatal("CLI result missing", err, out.String(), errout.String())
		}
		return result, runErr
	}
	result, err := run()
	if err == nil || result.FailureStage != "protocol-precheck" || result.InvocationID == "" || result.TerminalSHA256 == "" || result.Assessment != nil || result.RequestSHA256 != "" {
		t.Fatal("CLI consumed ticket before known protocol refusal", result, err)
	}
	for _, p := range []string{journal, trace, filepath.Join(root, ".parley-runtime", "verifier-starts"), result.RequestPath} {
		if _, err = os.Lstat(p); !os.IsNotExist(err) {
			t.Fatal("known protocol refusal consumed helper authority or executed", p, err)
		}
	}
	terminalPath := filepath.Join(root, ".parley-runtime", "invocations", result.InvocationID, "terminal.json")
	var terminal telemetry.Record
	raw, err := readVerificationJSON(terminalPath, &terminal)
	if err != nil || sha256Hex(string(raw)) != result.TerminalSHA256 {
		t.Fatal("refusal terminal hash missing", err)
	}
	if terminal.Metadata.RunID != result.RunID || terminal.Metadata.Agent != "reviewer" || terminal.Metadata.Phase != runner.CapturedVerificationPhase || terminal.Metadata.Idea != "idea-x" || terminal.StartedAt != nil || terminal.PID != nil || terminal.Outcome == nil || terminal.Outcome.FailureClass == nil || *terminal.Outcome.FailureClass != "protocol_context_refused" {
		t.Fatal("CLI refusal provenance invalid", terminal)
	}
	var parent trajectoryVerificationResult
	parentRaw, err := readVerificationJSON(filepath.Join(filepath.Dir(result.RequestPath), "parent-result.json"), &parent)
	if err != nil || !sameVerificationJSON(parent, result) {
		t.Fatal("CLI lost failed parent result", parent, err)
	}
	after, err := trajectory.Inspect(ctx, root, "idea-x")
	if err != nil {
		t.Fatal(err)
	}
	afterLedger, err := b.Store.Inspect(ctx)
	if err != nil {
		t.Fatal(err)
	}
	if !sameVerificationJSON(before, after) || !sameVerificationJSON(ledger, afterLedger) {
		t.Fatal("CLI refusal rewrote original state or ledger")
	}
	restore()
	next, err := run()
	if err != nil || next.Assessment == nil || next.Assessment.Outcome != trajectory.Regression || next.InvocationID == result.InvocationID || next.RunID == result.RunID {
		t.Fatal("restored protocol lost original verification opportunity", next, err)
	}
	actual, err := os.ReadFile(trace)
	if err != nil || string(actual) != "reviewer:original\nreviewer:broken\nreviewer:broken\nreviewer:original\n" {
		t.Fatal("restored CLI did not execute original AB/BA checks", string(actual), err)
	}
	for p, original := range map[string][]byte{terminalPath: raw, filepath.Join(filepath.Dir(result.RequestPath), "parent-result.json"): parentRaw} {
		current, err := os.ReadFile(p)
		if err != nil || !bytes.Equal(current, original) {
			t.Fatal("restored CLI rewrote refusal history", p, err)
		}
	}
}

func TestDriverAdapterProtocolPrecheckUsesSelectedImplementer(t *testing.T) {
	root, _, _, reviewer := trajectoryHelperFixture(t, "actual-helper")
	builder := reviewer
	builder.ID = "builder"
	runDir := filepath.Join(root, "parley-deck", "runs", "adapter-precheck")
	ideaDir := filepath.Join(root, "parley-deck", "ideas", "idea-x")
	// The base participant list deliberately differs from the implementer. Both
	// actual Fixup and its precheck must use withParticipants(implementer).
	o := driverImplOps{root: root, ideaSlug: "idea-x", ideaDir: ideaDir, implementer: "builder", out: io.Discard,
		base: runner.Options{Root: root, RunID: "adapter-precheck", Idea: protocol.IdeaStatus{Slug: "idea-x", Path: ideaDir, Participants: []string{"reviewer"}}, Agents: []agents.Discovery{reviewer, builder}, Store: store.New(runDir)}}
	refuseAppProtocolCache(t, root, 8)
	if err := o.PrecheckFixup(context.Background()); err == nil {
		t.Fatal("production adapter omitted protocol precheck")
	}
	paths, err := filepath.Glob(filepath.Join(root, ".parley-runtime", "invocations", "*", "terminal.json"))
	if err != nil {
		t.Fatal(err)
	}
	found := 0
	for _, p := range paths {
		var r telemetry.Record
		if _, err := readVerificationJSON(p, &r); err != nil {
			t.Fatal(err)
		}
		if r.Metadata.RunID == "adapter-precheck" {
			found++
			if r.Metadata.Agent != "builder" || r.Metadata.Idea != "idea-x" || r.Metadata.Phase != "fixup" || r.StartedAt != nil || r.Outcome == nil || r.Outcome.FailureClass == nil || *r.Outcome.FailureClass != "protocol_context_refused" {
				t.Fatal("precheck used a different implementer/scope", r)
			}
		}
	}
	if found != 1 {
		t.Fatal("adapter did not retain exactly one actual refusal", found)
	}
}
