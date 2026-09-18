package app

import (
	"context"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"

	"parley-deck-cli/internal/agents"
	"parley-deck-cli/internal/config"
	"parley-deck-cli/internal/runner"
	"parley-deck-cli/internal/store"
	"parley-deck-cli/internal/telemetry"
	"parley-deck-cli/internal/trajectory"
)

// Attended recovery of exactly one retained pre-start budget-refused verifier
// launch. Preview and apply never launch a process, never spend, refund or
// reset a budget, and never resolve a trajectory. The subsequent explicit
// relaunch is a real, separately checked and separately charged verifier launch
// under the invocation the immutable recovery pinned; its budget admission is
// the ordinary operator budget binding, unchanged.

// readRetainedVerifierRun loads the exact private parent request the original
// verify retained for this run and refuses aliases and mismatched roots, runs
// or ideas: the embedded ticket must be well-formed and bound to the canonical
// worktree and the requested idea/run.
func readRetainedVerifierRun(root, idea, runID string) (trajectory.HelperRequest, []byte, string, error) {
	var req trajectory.HelperRequest
	if !trajectoryPathID(idea) || !trajectoryPathID(runID) {
		return req, nil, "", errors.New("trajectory verifier recovery requires exact idea and run identities")
	}
	abs, err := filepath.Abs(root)
	if err == nil {
		abs, err = filepath.EvalSymlinks(abs)
	}
	if err != nil {
		return req, nil, "", err
	}
	path := filepath.Join(abs, ".parley-runtime", "trajectory-verification", runID, "request.json")
	raw, err := readVerificationJSON(path, &req)
	if err != nil {
		return req, nil, "", errors.New("retained verifier run request is unavailable")
	}
	if req.Version != 1 || req.Ticket.Root != abs || req.Ticket.RunID != runID || req.Ticket.Request.Idea != idea {
		return req, nil, "", errors.New("retained verifier run request is aliased or mismatched")
	}
	if _, err = req.Ticket.SHA256(); err != nil {
		return req, nil, "", err
	}
	return req, raw, abs, nil
}

// runTrajectoryVerifierRecovery previews or applies the single immutable
// recovery of the retained budget-refused verifier launch. The apply rechecks
// the exact inspected preview digest against live evidence before writing.
func runTrajectoryVerifierRecovery(ctx context.Context, args []string, out, errout io.Writer) int {
	f := flag.NewFlagSet("trajectory recover-verifier", flag.ContinueOnError)
	f.SetOutput(errout)
	root := f.String("dir", ".", "original trajectory worktree")
	idea := f.String("idea", "", "idea identity")
	run := f.String("run", "", "original independent verifier run")
	replacement := f.String("replacement", "", "selected replacement verifier invocation identity")
	expected := f.String("sha256", "", "exact inspected recovery preview digest")
	yes := f.Bool("yes", false, "apply only this exact inspected recovery; executes no work")
	if err := f.Parse(args); err != nil || f.NArg() != 0 || !trajectoryPathID(*idea) || !trajectoryPathID(*run) ||
		(!*yes && *expected != "") || (*yes && *expected == "") || (*replacement != "" && !trajectoryPathID(*replacement)) {
		fmt.Fprintln(errout, "usage: parley trajectory recover-verifier --dir DIR --idea ID --run RUN [--replacement ID] [--sha256 SHA --yes]")
		return 2
	}
	req, _, _, err := readRetainedVerifierRun(*root, *idea, *run)
	if err != nil {
		fmt.Fprintf(errout, "trajectory verifier recovery: %v\n", err)
		return 1
	}
	selected := *replacement
	if selected == "" {
		if selected, err = telemetry.NewID(); err != nil {
			fmt.Fprintf(errout, "trajectory verifier recovery: %v\n", err)
			return 1
		}
	}
	var p trajectory.CapturedRecoveryPreview
	if *yes {
		p, err = trajectory.ApplyCapturedVerificationRecovery(ctx, req.Ticket, selected, *expected)
	} else {
		p, err = trajectory.PreviewCapturedVerificationRecovery(ctx, req.Ticket, selected)
	}
	if err != nil {
		fmt.Fprintf(errout, "trajectory verifier recovery: %v\n", err)
		return 1
	}
	if err = json.NewEncoder(out).Encode(struct {
		Preview trajectory.CapturedRecoveryPreview `json:"preview"`
		SHA256  string                             `json:"sha256"`
		Applied bool                               `json:"applied"`
	}{p, p.SHA256(), *yes}); err != nil {
		fmt.Fprintln(errout, "cannot write verifier recovery result; preserve the record and repeat the exact request")
		return 1
	}
	return 0
}

// verifierRelaunchPlan is the exact inspected relaunch authority: the ticket,
// the validated immutable recovery identity (which names the bound replacement
// invocation and the refused one it repairs), the ticket-bound verifier and the
// run identity. The apply rechecks this plan's digest before any launch.
type verifierRelaunchPlan struct {
	Version      int                                 `json:"version"`
	Root         string                              `json:"root"`
	Idea         string                              `json:"idea"`
	RunID        string                              `json:"run_id"`
	TicketSHA256 string                              `json:"ticket_sha256"`
	Recovery     trajectory.CapturedRecoveryIdentity `json:"recovery"`
	Verifier     string                              `json:"verifier"`
}

func (p verifierRelaunchPlan) SHA256() string {
	data, err := json.MarshalIndent(p, "", "  ")
	if err != nil {
		return ""
	}
	return sha256Hex(string(append(data, '\n')))
}

type verifierRelaunchOutput struct {
	Plan            verifierRelaunchPlan               `json:"plan"`
	SHA256          string                             `json:"sha256"`
	Recovered       *trajectory.RecoveredParentPreview `json:"recovered,omitempty"`
	RecoveredSHA256 string                             `json:"recovered_sha256,omitempty"`
	Applied         bool                               `json:"applied"`
}

// buildVerifierRelaunchPlan validates everything the relaunch depends on
// WITHOUT launching anything: the read-only recovery identity, the retained
// original failed parent (still exactly the pre-start budget refusal of the
// recovery's refused invocation), the live quorum/material scope, and the
// configured roster resolution of the ticket-bound verifier.
func buildVerifierRelaunchPlan(ctx context.Context, root string, req trajectory.HelperRequest, reqRaw []byte) (verifierRelaunchPlan, agents.Discovery, error) {
	var plan verifierRelaunchPlan
	var agent agents.Discovery
	identity, err := trajectory.ReadCapturedVerificationRecovery(ctx, req.Ticket)
	if err != nil {
		return plan, agent, err
	}
	base := filepath.Join(root, ".parley-runtime", "trajectory-verification", req.Ticket.RunID)
	var original trajectoryVerificationResult
	if _, err = readVerificationJSON(filepath.Join(base, "parent-result.json"), &original); err != nil {
		return plan, agent, errors.New("recovered verifier relaunch requires the retained original refused parent result")
	}
	if original.Version != 1 || original.RunID != req.Ticket.RunID ||
		original.RequestPath != filepath.Join(base, "request.json") || original.RequestSHA256 != sha256Hex(string(reqRaw)) ||
		original.InvocationID != identity.RefusedInvocationID || original.TerminalSHA256 != identity.RefusedTerminalSHA256 ||
		original.ReceiptSHA256 != "" || original.Assessment != nil || original.FailureStage != "launch" || !original.TrajectoryPending {
		return plan, agent, errors.New("original retained result is not the exact pre-start budget refusal")
	}
	if err = checkTrajectoryHelperScope(ctx, req); err != nil {
		return plan, agent, err
	}
	discovered, err := discoverConfigured(ctx, root)
	if err != nil {
		return plan, agent, errors.New("cannot load configured participants")
	}
	mapping, err := config.LoadRosterAdapters(root)
	if err != nil {
		return plan, agent, errors.New("invalid roster mapping")
	}
	agent, err = resolveMeasuredAgent(req.Ticket.Request.Verifier, discovered, mapping)
	if err != nil {
		return plan, agent, errors.New("selected verifier is unavailable")
	}
	if agents.LaunchModeOrDefault(agent.LaunchMode) != agents.LaunchHeadless {
		return plan, agent, errors.New("trajectory verification currently requires a headless participant CLI")
	}
	ticketSHA, err := req.Ticket.SHA256()
	if err != nil {
		return plan, agent, err
	}
	plan = verifierRelaunchPlan{Version: 1, Root: root, Idea: req.Ticket.Request.Idea, RunID: req.Ticket.RunID,
		TicketSHA256: ticketSHA, Recovery: identity, Verifier: req.Ticket.Request.Verifier}
	return plan, agent, nil
}

// runTrajectoryVerifierRelaunch previews the exact relaunch plan, or — with
// --yes and the exact inspected plan digest — launches the bound verifier once
// under the recovered invocation and publishes the immutable recovered parent
// observation over the complete validated lineage.
func runTrajectoryVerifierRelaunch(ctx context.Context, args []string, out, errout io.Writer) int {
	executable, err := os.Executable()
	if err != nil {
		fmt.Fprintf(errout, "trajectory verifier relaunch: %v\n", err)
		return 1
	}
	return runTrajectoryVerifierRelaunchWith(ctx, args, out, errout, executable)
}

// runTrajectoryVerifierRelaunchWith takes the helper executable explicitly so
// the exact app relaunch path is exercisable in tests (same convention as
// verifyTrajectoryWithAgent); production passes os.Executable.
func runTrajectoryVerifierRelaunchWith(ctx context.Context, args []string, out, errout io.Writer, executable string) int {
	f := flag.NewFlagSet("trajectory relaunch-verifier", flag.ContinueOnError)
	f.SetOutput(errout)
	root := f.String("dir", ".", "original trajectory worktree")
	idea := f.String("idea", "", "idea identity")
	run := f.String("run", "", "original independent verifier run")
	expected := f.String("sha256", "", "exact inspected relaunch plan digest")
	timeout := f.Duration("timeout", 0, "hard agent timeout; default is its configured timeout")
	yes := f.Bool("yes", false, "launch the bound verifier once under the recovered invocation")
	if err := f.Parse(args); err != nil || f.NArg() != 0 || !trajectoryPathID(*idea) || !trajectoryPathID(*run) ||
		(!*yes && *expected != "") || (*yes && *expected == "") || *timeout < 0 {
		fmt.Fprintln(errout, "usage: parley trajectory relaunch-verifier --dir DIR --idea ID --run RUN [--sha256 SHA --timeout D --yes]")
		return 2
	}
	req, reqRaw, abs, err := readRetainedVerifierRun(*root, *idea, *run)
	if err != nil {
		fmt.Fprintf(errout, "trajectory verifier relaunch: %v\n", err)
		return 1
	}
	plan, agent, err := buildVerifierRelaunchPlan(ctx, abs, req, reqRaw)
	if err != nil {
		fmt.Fprintf(errout, "trajectory verifier relaunch: %v\n", err)
		return 1
	}
	result := verifierRelaunchOutput{Plan: plan, SHA256: plan.SHA256()}
	if *yes {
		if plan.SHA256() != *expected {
			fmt.Fprintln(errout, "trajectory verifier relaunch: plan changed since preview; inspect again")
			return 1
		}
		recovered, err := relaunchRecoveredVerifier(ctx, abs, req, reqRaw, plan.Recovery, agent, *timeout, executable, errout)
		if err != nil {
			fmt.Fprintf(errout, "trajectory verifier relaunch refused: %v\n", err)
			return 1
		}
		result.Recovered, result.RecoveredSHA256, result.Applied = &recovered, recovered.SHA256(), true
	}
	if err := json.NewEncoder(out).Encode(result); err != nil {
		fmt.Fprintln(errout, "cannot write verifier relaunch result; inspect recover-verifier-parent for the retained replacement evidence")
		return 1
	}
	return 0
}

// relaunchRecoveredVerifier executes the one bound verifier relaunch through
// the real runner and publishes the recovered parent observation. The launch
// passes the protocol precheck, the trusted recovery context, the complete
// recovered-lineage revalidation at the reservation boundary, and the ordinary
// operator budget admission and settlement — the recovery is never a bypass.
// The original failed parent result, the original request and the refused
// invocation's records are never modified; the recovered logs are separate
// files. A failed relaunch publishes nothing and stays unresolved.
func relaunchRecoveredVerifier(ctx context.Context, root string, req trajectory.HelperRequest, reqRaw []byte, identity trajectory.CapturedRecoveryIdentity, agent agents.Discovery, timeout time.Duration, executable string, progress io.Writer) (trajectory.RecoveredParentPreview, error) {
	var preview trajectory.RecoveredParentPreview
	runID := req.Ticket.RunID
	idea := req.Ticket.Request.Idea
	base := filepath.Join(root, ".parley-runtime", "trajectory-verification", runID)
	var observed telemetry.Record
	launchInfo := runner.LaunchInfo{RunID: runID, Idea: idea, Phase: runner.CapturedVerificationPhase,
		Store: store.New(filepath.Join(base, "run")), Observe: func(record telemetry.Record) { observed = record }}
	// Refuse a known protocol-context refusal before the pinned invocation is
	// spent; such a refusal is recorded under a fresh ID, never the bound one.
	if err := runner.PrecheckProtocolLaunch(ctx, root, agent, launchInfo); err != nil {
		return preview, err
	}
	bound, err := runner.WithCapturedVerificationRecovery(ctx, req.Ticket, identity.InvocationID)
	if err != nil {
		return preview, err
	}
	bound = runner.WithLaunchInfo(bound, launchInfo)
	requestPath := filepath.Join(base, "request.json")
	prompt := trajectoryVerifierPrompt(agent.ID, idea, trajectoryVerifierHelperCommand(executable, requestPath, sha256Hex(string(reqRaw))))
	res := runner.RunConsult(bound, runner.ConsultOptions{Root: root, Agent: agent, Prompt: prompt,
		Timeout: timeout, StdoutPath: filepath.Join(base, "agent-recovered.stdout.log"), StderrPath: filepath.Join(base, "agent-recovered.stderr.log"), Progress: progress})
	if res.InvocationID != "" && res.InvocationID != identity.InvocationID {
		return preview, errors.New("recovered verifier launch escaped its pinned invocation")
	}
	if res.ExitError != "" || res.AgentExit != 0 ||
		observed.SchemaVersion != telemetry.SchemaVersion || observed.Type != "invocation.terminal" ||
		observed.InvocationID != identity.InvocationID || observed.Metadata.RunID != runID || observed.Metadata.Idea != idea ||
		observed.Metadata.Agent != agent.ID || observed.Metadata.Phase != runner.CapturedVerificationPhase ||
		observed.StartedAt == nil || observed.CompletedAt == nil || observed.PID == nil || *observed.PID <= 0 ||
		observed.Outcome == nil || observed.Outcome.Status != "process-exited" || observed.Outcome.FailureClass != nil || observed.Outcome.ExitCode == nil || *observed.Outcome.ExitCode != 0 {
		return preview, errors.New("recovered independent trajectory verifier has no matching successful observed terminal")
	}
	var retained trajectory.HelperRequest
	current, err := readVerificationJSON(requestPath, &retained)
	if err != nil || sha256Hex(string(current)) != sha256Hex(string(reqRaw)) || !sameVerificationJSON(retained, req) {
		return preview, errors.New("trajectory helper request changed after recovered launch")
	}
	if err = checkTrajectoryHelperScope(ctx, req); err != nil {
		return preview, err
	}
	receipt, observation, err := trajectory.ReadCapturedVerification(ctx, req.Ticket, identity.InvocationID)
	if err != nil {
		return preview, err
	}
	if receipt.HelperPID == os.Getpid() || receipt.HelperPID <= 0 || receipt.FinishedAt.Before(*observed.StartedAt) || receipt.FinishedAt.After(*observed.CompletedAt) {
		return preview, errors.New("trajectory helper receipt is outside the observed verifier invocation")
	}
	if _, err = trajectory.AssessCaptured(req.Ticket.Request, observation); err != nil {
		return preview, err
	}
	p, err := trajectory.PreviewRecoveredParent(ctx, root, idea, runID)
	if err != nil {
		return preview, err
	}
	return trajectory.PublishRecoveredParent(ctx, root, idea, runID, p.SHA256())
}

// A completed replacement may survive failed observation publication. Expose
// the existing exact preview/apply API without attempting another verifier
// launch or changing its original charged invocation and failed parent result.
func runTrajectoryVerifierParentRecovery(ctx context.Context, args []string, out, errout io.Writer) int {
	f := flag.NewFlagSet("trajectory recover-verifier-parent", flag.ContinueOnError)
	f.SetOutput(errout)
	root := f.String("dir", ".", "original trajectory worktree")
	idea := f.String("idea", "", "idea identity")
	run := f.String("run", "", "original independent verifier run")
	expected := f.String("sha256", "", "exact inspected recovered parent preview digest")
	yes := f.Bool("yes", false, "publish the retained replacement observation without launching work")
	if err := f.Parse(args); err != nil || f.NArg() != 0 || !trajectoryPathID(*idea) || !trajectoryPathID(*run) || (!*yes && *expected != "") || (*yes && *expected == "") {
		fmt.Fprintln(errout, "usage: parley trajectory recover-verifier-parent --dir DIR --idea ID --run RUN [--sha256 SHA --yes]")
		return 2
	}
	var p trajectory.RecoveredParentPreview
	var err error
	if *yes {
		p, err = trajectory.PublishRecoveredParent(ctx, *root, *idea, *run, *expected)
	} else {
		p, err = trajectory.PreviewRecoveredParent(ctx, *root, *idea, *run)
	}
	if err != nil {
		fmt.Fprintf(errout, "trajectory verifier parent recovery: %v\n", err)
		return 1
	}
	if err = json.NewEncoder(out).Encode(struct {
		Preview trajectory.RecoveredParentPreview `json:"preview"`
		SHA256  string                            `json:"sha256"`
		Applied bool                              `json:"applied"`
	}{p, p.SHA256(), *yes}); err != nil {
		fmt.Fprintln(errout, "cannot write recovered parent result; preserve the record and repeat the exact request")
		return 1
	}
	return 0
}
