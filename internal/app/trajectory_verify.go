package app

import (
	"context"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"slices"
	"strings"
	"time"

	"parley-deck-cli/internal/agents"
	"parley-deck-cli/internal/config"
	"parley-deck-cli/internal/driver"
	"parley-deck-cli/internal/fsutil"
	"parley-deck-cli/internal/protocol"
	"parley-deck-cli/internal/runner"
	"parley-deck-cli/internal/store"
	"parley-deck-cli/internal/telemetry"
	"parley-deck-cli/internal/trajectory"
)

type trajectoryHelperRequest struct {
	Version      int                           `json:"version"`
	Ticket       trajectory.VerificationTicket `json:"ticket"`
	Participants []string                      `json:"participants"`
	Criteria     []trajectory.Criterion        `json:"criteria"`
}

type trajectoryVerificationResult struct {
	Version           int                    `json:"version"`
	RunID             string                 `json:"run_id"`
	RequestPath       string                 `json:"request_path"`
	RequestSHA256     string                 `json:"request_sha256"`
	InvocationID      string                 `json:"invocation_id"`
	TerminalSHA256    string                 `json:"terminal_sha256"`
	ReceiptSHA256     string                 `json:"receipt_sha256"`
	Assessment        *trajectory.Assessment `json:"assessment"`
	FailureStage      string                 `json:"failure_stage"`
	TrajectoryPending bool                   `json:"trajectory_pending"`
}

func runTrajectoryVerify(ctx context.Context, args []string, out, errout io.Writer) int {
	f := flag.NewFlagSet("trajectory verify", flag.ContinueOnError)
	f.SetOutput(errout)
	root := f.String("dir", ".", "original trajectory worktree")
	idea := f.String("idea", "", "idea identity")
	verifier := f.String("verifier", "", "selected independent participant")
	timeout := f.Duration("timeout", 0, "hard agent timeout; default is its configured timeout")
	yes := f.Bool("yes", false, "invoke this configured verifier once")
	if err := f.Parse(args); err != nil || f.NArg() != 0 || !*yes || *timeout < 0 || !trajectoryPathID(*idea) || !trajectoryPathID(*verifier) {
		fmt.Fprintln(errout, "usage: parley trajectory verify --dir DIR --idea ID --verifier ID [--timeout D] --yes")
		return 2
	}
	discovered, err := discoverConfigured(ctx, *root)
	if err != nil {
		fmt.Fprintln(errout, "trajectory verify: cannot load configured participants")
		return 1
	}
	mapping, err := config.LoadRosterAdapters(*root)
	if err != nil {
		fmt.Fprintln(errout, "trajectory verify: invalid roster mapping")
		return 1
	}
	agent, err := resolveMeasuredAgent(*verifier, discovered, mapping)
	if err != nil {
		fmt.Fprintln(errout, "trajectory verify: selected participant is unavailable")
		return 1
	}
	executable, err := os.Executable()
	if err != nil {
		return 1
	}
	result, err := verifyTrajectoryWithAgent(ctx, *root, *idea, agent, *timeout, executable, errout)
	if encodeErr := json.NewEncoder(out).Encode(result); encodeErr != nil {
		return 1
	}
	if err != nil {
		fmt.Fprintf(errout, "trajectory verification refused: %v\n", err)
		return 1
	}
	return 0
}

func trajectoryPathID(value string) bool {
	return telemetry.SafeLabel(value) != nil && len(value) <= 128 && !strings.ContainsAny(value, `/\\`) && !strings.Contains(value, "..")
}

func trajectoryMaterialScope(root, idea, implementer, verifier string) ([]string, []trajectory.Criterion, error) {
	if !trajectoryPathID(idea) || !trajectoryPathID(implementer) || !trajectoryPathID(verifier) || implementer == verifier {
		return nil, nil, errors.New("trajectory verification requires distinct original and selected participants")
	}
	status, err := protocol.ReadWorkspaceStatus(root)
	if err != nil {
		return nil, nil, err
	}
	current, found := findIdeaStatus(status, idea)
	if !found || !slices.Contains(current.Participants, implementer) || !slices.Contains(current.Participants, verifier) {
		return nil, nil, errors.New("trajectory author or selected verifier is outside the idea's existing quorum")
	}
	checks, isList, err := driver.ReadChecksContract(current.Path)
	if err != nil || !isList || len(checks) == 0 {
		return nil, nil, errors.New("trajectory verification requires the original named material checks")
	}
	criteria := make([]trajectory.Criterion, len(checks))
	for i, c := range checks {
		criteria[i] = trajectory.Criterion{Name: c.Name, Command: c.Command}
	}
	return current.Participants, criteria, nil
}

func checkTrajectoryHelperScope(req trajectoryHelperRequest) error {
	if req.Version != 1 {
		return errors.New("unsupported trajectory helper request")
	}
	if _, err := req.Ticket.SHA256(); err != nil {
		return err
	}
	r := req.Ticket.Request
	participants, criteria, err := trajectoryMaterialScope(req.Ticket.Root, r.Idea, r.Implementer, r.Verifier)
	if err != nil {
		return err
	}
	if !slices.Equal(participants, req.Participants) || !sameVerificationJSON(criteria, req.Criteria) || len(criteria) != len(r.Criteria) {
		return errors.New("trajectory helper quorum or material scope changed")
	}
	for i, c := range criteria {
		if c.Name != r.Criteria[i].Name || sha256Hex(c.Command) != r.Criteria[i].CommandSHA256 {
			return errors.New("trajectory helper criteria differ from the original frozen scope")
		}
	}
	return nil
}

func trajectoryRuntime(root string) (string, error) {
	if runtime.GOOS == "windows" {
		return "", errors.New("trajectory verification requires a POSIX execution host")
	}
	if err := exec.Command("git", "--no-optional-locks", "-C", root, "check-ignore", "-q", "--", ".parley-runtime/").Run(); err != nil {
		return "", errors.New("trajectory verification requires ignored private .parley-runtime storage")
	}
	tracked, err := exec.Command("git", "--no-optional-locks", "-C", root, "ls-files", "-z", "--", ".parley-runtime/").Output()
	if err != nil || len(tracked) != 0 {
		return "", errors.New("trajectory runtime cannot contain tracked files")
	}
	base := filepath.Join(root, ".parley-runtime", "trajectory-verification")
	for _, path := range []string{filepath.Dir(base), base} {
		if err = os.Mkdir(path, 0700); err != nil && !os.IsExist(err) {
			return "", err
		}
		info, err := os.Lstat(path)
		if err != nil || !info.IsDir() {
			return "", errors.New("trajectory runtime requires real directories")
		}
		if err = syncTrajectoryRuntimeParent(path); err != nil {
			return "", err
		}
	}
	return base, nil
}

func syncTrajectoryRuntimeParent(path string) error {
	dir, err := os.Open(filepath.Dir(path))
	if err != nil {
		return err
	}
	defer dir.Close()
	return fsutil.SyncFile(dir)
}

func writeTrajectoryRuntimeJSON(path string, value any) error {
	data, err := json.MarshalIndent(value, "", "  ")
	if err != nil || len(data) > 1<<20 {
		return errors.New("invalid or excessive trajectory runtime artifact")
	}
	f, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_EXCL, 0600)
	if err != nil {
		return err
	}
	_, writeErr := f.Write(data)
	if err = errors.Join(writeErr, fsutil.SyncFile(f), f.Close()); err != nil {
		return err
	}
	return syncTrajectoryRuntimeParent(path)
}

func verifyTrajectoryWithAgent(ctx context.Context, root, idea string, agent agents.Discovery, timeout time.Duration, executable string, progress io.Writer) (result trajectoryVerificationResult, resultErr error) {
	result.Version, result.TrajectoryPending = 1, true
	if agents.LaunchModeOrDefault(agent.LaunchMode) != agents.LaunchHeadless {
		return result, errors.New("trajectory verification currently requires a headless participant CLI")
	}
	root, err := filepath.Abs(root)
	if err == nil {
		root, err = filepath.EvalSymlinks(root)
	}
	if err != nil {
		return result, err
	}
	request, err := trajectory.FreezeCaptured(ctx, root, idea, agent.ID)
	if err != nil {
		return result, err
	}
	participants, criteria, err := trajectoryMaterialScope(root, idea, request.Implementer, agent.ID)
	if err != nil {
		return result, err
	}
	runID, err := telemetry.NewID()
	if err != nil {
		return result, err
	}
	req := trajectoryHelperRequest{1, trajectory.VerificationTicket{Version: 1, Root: root, RunID: runID, Request: request}, participants, criteria}
	if err = checkTrajectoryHelperScope(req); err != nil {
		return result, err
	}
	base, err := trajectoryRuntime(root)
	if err != nil {
		return result, err
	}
	dir := filepath.Join(base, runID)
	if err = os.Mkdir(dir, 0700); err != nil {
		return result, err
	}
	if err = syncTrajectoryRuntimeParent(dir); err != nil {
		return result, err
	}
	result.RunID, result.RequestPath = runID, filepath.Join(dir, "request.json")
	stage := "preparation"
	defer func() {
		if resultErr != nil {
			result.FailureStage, result.Assessment = stage, nil
		}
		if err := writeTrajectoryRuntimeJSON(filepath.Join(dir, "parent-result.json"), result); err != nil {
			result.FailureStage, result.Assessment = "parent-publication", nil
			resultErr = errors.Join(resultErr, err)
		}
	}()
	req.Ticket, err = trajectory.PrepareCapturedVerification(ctx, root, idea, agent.ID, runID)
	if err != nil {
		return result, err
	}
	if !sameVerificationJSON(req.Ticket.Request, request) {
		return result, errors.New("trajectory changed during request preparation")
	}
	if err = writeTrajectoryRuntimeJSON(result.RequestPath, req); err != nil {
		return result, err
	}
	raw, err := os.ReadFile(result.RequestPath)
	if err != nil {
		return result, err
	}
	result.RequestSHA256 = sha256Hex(string(raw))
	ctx, err = runner.WithCapturedVerification(ctx, req.Ticket)
	if err != nil {
		return result, err
	}
	var observed telemetry.Record
	ctx = runner.WithLaunchInfo(ctx, runner.LaunchInfo{RunID: runID, Idea: idea,
		Phase: runner.CapturedVerificationPhase, Store: store.New(filepath.Join(dir, "run")),
		Observe: func(record telemetry.Record) { observed = record }})
	quote := func(s string) string { return "'" + strings.ReplaceAll(s, "'", "'\\''") + "'" }
	command := quote(executable) + " trajectory verify-helper --request " + quote(result.RequestPath) + " --request-sha256 " + quote(result.RequestSHA256)
	prompt := fmt.Sprintf(`You are %s, the selected independent execution verifier for idea %s.
Invoke the exact helper command below once from this invocation. It runs the original material checks against the archived before/after sources in AB/BA order and retains actual observations. Do not edit code, scope, requests, journals, receipts or signatures. Do not substitute a written PASS or manually populated evidence. Preserve inherited PARLEY_RUN_ID, PARLEY_AGENT_ID and PARLEY_PROC_MARKER. If the helper fails, report the actual failure and stop; do not retry. After it returns, report its outcome briefly and stop. This records an observation; it does not grant trajectory continuation or completion. Process attribution is not human authentication.
Verifier command: %s
`, agent.ID, idea, command)
	stage = "launch"
	res := runner.RunConsult(ctx, runner.ConsultOptions{Root: root, Agent: agent, Prompt: prompt,
		Timeout: timeout, StdoutPath: filepath.Join(dir, "agent.stdout.log"), StderrPath: filepath.Join(dir, "agent.stderr.log"), Progress: progress})
	result.InvocationID = res.InvocationID
	if trajectoryPathID(res.InvocationID) {
		var terminal telemetry.Record
		terminalRaw, err := readVerificationJSON(filepath.Join(root, ".parley-runtime", "invocations", res.InvocationID, "terminal.json"), &terminal)
		if err == nil && sameVerificationJSON(terminal, observed) {
			result.TerminalSHA256 = sha256Hex(string(terminalRaw))
		}
	}
	if res.ExitError != "" || res.AgentExit != 0 || result.TerminalSHA256 == "" ||
		observed.SchemaVersion != telemetry.SchemaVersion || observed.Type != "invocation.terminal" ||
		observed.InvocationID != res.InvocationID || observed.Metadata.RunID != runID || observed.Metadata.Idea != idea ||
		observed.Metadata.Agent != agent.ID || observed.Metadata.Phase != runner.CapturedVerificationPhase ||
		observed.StartedAt == nil || observed.CompletedAt == nil || observed.PID == nil || *observed.PID <= 0 ||
		observed.Outcome == nil || observed.Outcome.Status != "process-exited" || observed.Outcome.FailureClass != nil || observed.Outcome.ExitCode == nil || *observed.Outcome.ExitCode != 0 {
		return result, errors.New("independent trajectory verifier has no matching successful observed terminal")
	}
	stage = "acceptance"
	var retained trajectoryHelperRequest
	current, err := readVerificationJSON(result.RequestPath, &retained)
	if err != nil || sha256Hex(string(current)) != result.RequestSHA256 || !sameVerificationJSON(retained, req) {
		return result, errors.New("trajectory helper request changed after launch")
	}
	if err = checkTrajectoryHelperScope(req); err != nil {
		return result, err
	}
	receipt, observation, err := trajectory.ReadCapturedVerification(ctx, req.Ticket, res.InvocationID)
	if err != nil {
		return result, err
	}
	if receipt.HelperPID == os.Getpid() || receipt.HelperPID <= 0 || receipt.FinishedAt.Before(*observed.StartedAt) || receipt.FinishedAt.After(*observed.CompletedAt) {
		return result, errors.New("trajectory helper receipt is outside the observed verifier invocation")
	}
	assessment, err := trajectory.AssessCaptured(req.Ticket.Request, observation)
	if err != nil {
		return result, err
	}
	encoded, err := json.MarshalIndent(receipt, "", "  ")
	if err != nil {
		return result, err
	}
	result.ReceiptSHA256 = sha256Hex(string(append(encoded, '\n')))
	result.Assessment = &assessment
	return result, nil
}

func runTrajectoryVerifyHelper(ctx context.Context, args []string, out, errout io.Writer) int {
	f := flag.NewFlagSet("trajectory verify-helper", flag.ContinueOnError)
	f.SetOutput(errout)
	path := f.String("request", "", "frozen private parent request")
	sha := f.String("request-sha256", "", "exact original request digest")
	if err := f.Parse(args); err != nil || f.NArg() != 0 || *path == "" || *sha == "" {
		return 2
	}
	if err := executeTrajectoryHelper(ctx, *path, *sha); err != nil {
		fmt.Fprintf(errout, "trajectory helper refused: %v\n", err)
		return 1
	}
	fmt.Fprintln(out, "Trajectory helper execution receipt recorded; trajectory remains pending.")
	return 0
}

func executeTrajectoryHelper(ctx context.Context, path, expectedSHA string) error {
	var req trajectoryHelperRequest
	raw, err := readVerificationJSON(path, &req)
	if err != nil {
		return err
	}
	if sha256Hex(string(raw)) != expectedSHA || req.Version != 1 {
		return errors.New("trajectory helper request digest or version mismatch")
	}
	canonical, err := filepath.EvalSymlinks(path)
	if err == nil {
		canonical, err = filepath.Abs(canonical)
	}
	if err != nil || !trajectoryPathID(req.Ticket.RunID) || canonical != filepath.Join(req.Ticket.Root, ".parley-runtime", "trajectory-verification", req.Ticket.RunID, "request.json") {
		return errors.New("trajectory helper request is outside its original runtime directory")
	}
	invocation := os.Getenv("PARLEY_PROC_MARKER")
	if os.Getenv("PARLEY_RUN_ID") != req.Ticket.RunID || os.Getenv("PARLEY_AGENT_ID") != req.Ticket.Request.Verifier || !trajectoryPathID(invocation) {
		return errors.New("trajectory helper must run inside its selected verifier invocation")
	}
	if err = checkTrajectoryHelperScope(req); err != nil {
		return err
	}
	_, err = trajectory.ExecuteCapturedVerification(ctx, req.Ticket, invocation, req.Criteria, "")
	return err
}
