package app

import (
	"bytes"
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
	"strings"

	"parley-deck-cli/internal/agents"
	"parley-deck-cli/internal/driver"
	"parley-deck-cli/internal/evidence"
	"parley-deck-cli/internal/fsutil"
	"parley-deck-cli/internal/protocol"
	"parley-deck-cli/internal/runner"
	"parley-deck-cli/internal/telemetry"
)

// This runtime command is invoked BY the independently selected CLI agent.
// Inherited process markers provide attribution, not authentication against a
// malicious same-UID caller. The parent pins the request/report bytes and accepts
// only the actual helper's retained executions, never a model's textual PASS.
type evidenceVerificationRequest struct {
	Version      int    `json:"version"`
	AttemptID    string `json:"attempt_id"`
	Root         string `json:"root"`
	Idea         string `json:"idea"`
	RunID        string `json:"run_id"`
	Verifier     string `json:"verifier"`
	ReportSHA256 string `json:"report_sha256"`
	TreeSHA256   string `json:"tree_sha256"`
	RestSHA256   string `json:"implementation_rest_sha256"`
}

type evidenceVerificationReceipt struct {
	Version        int                        `json:"version"`
	RequestSHA256  string                     `json:"request_sha256"`
	OriginalSHA256 string                     `json:"original_report_sha256"`
	UpdatedSHA256  string                     `json:"updated_report_sha256"`
	Agent          string                     `json:"agent"`
	RunID          string                     `json:"run_id"`
	ProcessMarker  string                     `json:"process_marker"`
	HelperPID      int                        `json:"helper_pid"`
	Executions     []evidence.CriterionRecord `json:"executions"`
	Error          string                     `json:"error"`
}

func runEvidenceVerify(ctx context.Context, args []string, stdout, stderr io.Writer) int {
	if len(args) > 0 && args[0] == "refusals" {
		return runEvidenceRefusals(ctx, args[1:], stdout, stderr)
	}
	if len(args) == 0 || args[0] != "verify" {
		fmt.Fprintln(stderr, "usage: parley evidence verify --request PATH --request-sha256 SHA (independent verifier runtime only)")
		return 2
	}
	flags := flag.NewFlagSet("evidence verify", flag.ContinueOnError)
	flags.SetOutput(stderr)
	path := flags.String("request", "", "frozen runtime request path")
	wantSHA := flags.String("request-sha256", "", "exact request byte digest supplied by the driver")
	if err := flags.Parse(args[1:]); err != nil || flags.NArg() != 0 || *path == "" || *wantSHA == "" {
		return 2
	}
	if err := executeIndependentVerification(ctx, *path, *wantSHA); err != nil {
		fmt.Fprintf(stderr, "independent evidence verification refused: %v\n", err)
		return 1
	}
	fmt.Fprintln(stdout, "Independent criterion executions and bound attestation recorded.")
	return 0
}

func readVerificationJSON(path string, target any) ([]byte, error) {
	info, err := os.Lstat(path)
	if err != nil {
		return nil, err
	}
	if !info.Mode().IsRegular() || info.Size() > 16<<20 {
		return nil, errors.New("verification artifact must be a bounded regular file")
	}
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	opened, err := f.Stat()
	if err != nil || !os.SameFile(info, opened) {
		return nil, errors.New("verification artifact changed during open")
	}
	data, err := io.ReadAll(io.LimitReader(f, (16<<20)+1))
	if err != nil {
		return nil, err
	}
	if len(data) > 16<<20 {
		return nil, errors.New("verification artifact exceeds size limit")
	}
	decoder := json.NewDecoder(bytes.NewReader(data))
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(target); err != nil {
		return nil, err
	}
	var extra any
	if err := decoder.Decode(&extra); !errors.Is(err, io.EOF) {
		return nil, errors.New("trailing verification artifact content")
	}
	return data, nil
}

func writeVerificationJSON(path string, value any) error {
	data, err := json.MarshalIndent(value, "", "  ")
	if err != nil {
		return err
	}
	return writeVerificationBytes(path, data)
}

func writeVerificationBytes(path string, data []byte) error {
	f, err := os.CreateTemp(filepath.Dir(path), ".verification-*.tmp")
	if err != nil {
		return err
	}
	defer os.Remove(f.Name())
	if _, err := f.Write(data); err != nil {
		f.Close()
		return err
	}
	if err := fsutil.SyncFile(f); err != nil {
		f.Close()
		return err
	}
	if err := f.Close(); err != nil {
		return err
	}
	return fsutil.ReplaceSyncedFile(f.Name(), path)
}

func verificationBindings(req evidenceVerificationRequest, report *evidence.Report, ideaDir string) ([]driver.CheckCriterion, string, error) {
	if report.Idea != req.Idea || report.TreeSHA256 != req.TreeSHA256 {
		return nil, "", errors.New("report idea or tested tree changed")
	}
	criteria, isList, err := driver.ReadChecksContract(ideaDir)
	if err != nil {
		return nil, "", err
	}
	if !isList || len(criteria) == 0 || len(criteria) != len(report.Records) {
		return nil, "", errors.New("named criterion scope is absent or incomplete")
	}
	byName := make(map[string]evidence.CriterionRecord, len(report.Records))
	for _, rec := range report.Records {
		if _, exists := byName[rec.Name]; exists {
			return nil, "", errors.New("duplicate criterion evidence")
		}
		byName[rec.Name] = rec
		if rec.Provenance.Executor == "" || rec.Provenance.Executor == req.Verifier {
			return nil, "", errors.New("missing executor or self-verification")
		}
	}
	for _, criterion := range criteria {
		rec, found := byName[criterion.Name]
		if !found || rec.Command.Command != criterion.Command || rec.Command.CommandSHA256 != sha256Hex(criterion.Command) {
			return nil, "", fmt.Errorf("criterion %q command or scope differs from the original execution", criterion.Name)
		}
	}
	exclusions, err := definedEvidenceArtifacts(req.Root, ideaDir)
	if err != nil {
		return nil, "", err
	}
	tree, err := evidence.TreeDigest(req.Root, exclusions...)
	if err != nil {
		return nil, "", err
	}
	if tree != req.TreeSHA256 {
		return nil, "", errors.New("tested tree changed before or during independent verification")
	}
	rest, rel, err := implementationRestDigest(req.Root, ideaDir)
	if err != nil {
		return nil, "", err
	}
	if rest != req.RestSHA256 || report.ExtraDigests[rel] != rest {
		return nil, "", errors.New("non-evidence implementation content changed")
	}
	return criteria, tree, nil
}

func executeIndependentVerification(ctx context.Context, requestPath, requestSHA string) (resultErr error) {
	var req evidenceVerificationRequest
	raw, err := readVerificationJSON(requestPath, &req)
	if err != nil {
		return err
	}
	if sha256Hex(string(raw)) != requestSHA || req.Version != 1 {
		return errors.New("frozen verification request digest or version mismatch")
	}
	root, err := filepath.EvalSymlinks(req.Root)
	if err != nil {
		return err
	}
	root, err = filepath.Abs(root)
	if err != nil || root != req.Root {
		return errors.New("verification root is not canonical")
	}
	requestAbs, err := filepath.EvalSymlinks(requestPath)
	if err != nil {
		return err
	}
	requestAbs, err = filepath.Abs(requestAbs)
	if err != nil {
		return err
	}
	rel, err := filepath.Rel(filepath.Join(root, ".parley-runtime", "evidence-verification"), requestAbs)
	if err != nil || strings.HasPrefix(rel, ".."+string(filepath.Separator)) || filepath.IsAbs(rel) || filepath.Base(rel) != "request.json" || filepath.Dir(rel) == "." {
		return errors.New("request is outside its driver-owned runtime directory")
	}
	if req.Idea == "" || filepath.Base(req.Idea) != req.Idea || req.Idea == "." || req.Idea == ".." {
		return errors.New("invalid verification idea identity")
	}
	stage := "identity"
	var observedExecutions []evidence.CriterionRecord
	defer func() {
		if resultErr == nil {
			return
		}
		_, scope, scopeErr := verificationRefusalScope(root, req.Idea)
		if scopeErr != nil {
			resultErr = errors.Join(resultErr, errors.New("helper refusal has no safe idea scope"))
			return
		}
		rec, err := newVerificationRefusal("helper", stage, req.Idea, req.RunID, req.Verifier, requestSHA, req.ReportSHA256, os.Getenv("PARLEY_PROC_MARKER"))
		if err == nil {
			for _, execution := range observedExecutions {
				rec.Executions = append(rec.Executions, evidence.RefusalExecution{Name: telemetry.SafeLabel(execution.Name), Status: execution.Status, CommandSHA256: evidence.RefusalHash(execution.Command.CommandSHA256)})
			}
			var retained evidenceVerificationReceipt
			if raw, e := readVerificationJSON(filepath.Join(filepath.Dir(requestAbs), "result.json"), &retained); e == nil {
				rec.ReceiptSHA256 = evidence.RefusalHash(sha256Hex(string(raw)))
			}
			_, err = evidence.RetainVerificationRefusal(scope, rec)
		}
		if err != nil {
			resultErr = errors.Join(resultErr, errors.New("helper refusal retention failed"))
		}
	}()
	if req.AttemptID != filepath.Base(filepath.Dir(requestAbs)) {
		return errors.New("verification attempt identity differs from its frozen directory")
	}
	if req.RunID == "" || req.Verifier == "" || os.Getenv("PARLEY_RUN_ID") != req.RunID || os.Getenv("PARLEY_AGENT_ID") != req.Verifier || os.Getenv("PARLEY_PROC_MARKER") == "" {
		return errors.New("helper must run inside the selected independent verifier invocation")
	}
	receipt := evidenceVerificationReceipt{Version: 1, RequestSHA256: requestSHA,
		OriginalSHA256: req.ReportSHA256, Agent: os.Getenv("PARLEY_AGENT_ID"),
		RunID: os.Getenv("PARLEY_RUN_ID"), ProcessMarker: os.Getenv("PARLEY_PROC_MARKER"), HelperPID: os.Getpid()}
	receiptWritten := false
	defer func() {
		if receiptWritten {
			return
		}
		if resultErr != nil {
			receipt.Error = resultErr.Error()
		}
		resultErr = errors.Join(resultErr, writeVerificationJSON(filepath.Join(filepath.Dir(requestAbs), "result.json"), receipt))
	}()
	stage = "original-report"
	ideaDir := filepath.Join(root, protocol.DeckDir, "ideas", req.Idea)
	var report evidence.Report
	originalRaw, err := readVerificationJSON(filepath.Join(filepath.Dir(requestAbs), "original-report.json"), &report)
	if err != nil {
		return err
	}
	if sha256Hex(string(originalRaw)) != req.ReportSHA256 {
		return errors.New("original evidence report changed")
	}
	currentOriginal, err := os.ReadFile(evidence.ReportPath(ideaDir))
	if err != nil || !bytes.Equal(currentOriginal, originalRaw) {
		return errors.New("canonical evidence differs from the frozen original report")
	}
	stage = "bindings"
	criteria, _, err := verificationBindings(req, &report, ideaDir)
	if err != nil {
		return err
	}
	for _, rec := range report.Records {
		if rec.Provenance.Verifier != "" || rec.Provenance.VerifierRerun != nil {
			return errors.New("fresh execution report required; an earlier verifier cannot be replayed")
		}
	}
	stage = "execution"
	for _, criterion := range criteria {
		if err := ctx.Err(); err != nil {
			return err
		}
		_, beforeTree, err := verificationBindings(req, &report, ideaDir)
		if err != nil {
			return err
		}
		rec := evidence.RunCriterion(ctx, root, criterion.Name, criterion.Command, req.Verifier)
		receipt.Executions = append(receipt.Executions, rec)
		observedExecutions = receipt.Executions
		_, afterTree, err := verificationBindings(req, &report, ideaDir)
		if err != nil {
			return err
		}
		if rec.Status != evidence.StatusPass {
			return fmt.Errorf("independent criterion %q is %s", criterion.Name, rec.Status)
		}
		rerun := evidence.VerifierExecution{Command: rec.Command, TreeBeforeSHA256: beforeTree, TreeAfterSHA256: afterTree}
		if err := evidence.AttestExecution(&report, criterion.Name, req.Verifier, rerun); err != nil {
			return err
		}
	}
	if err := ctx.Err(); err != nil {
		return err
	}
	if _, _, err := verificationBindings(req, &report, ideaDir); err != nil {
		return err
	}
	restContent, restPath, err := implementationRestContent(root, ideaDir)
	if err != nil {
		return err
	}
	if err := evidence.AuthorizeCompletionTransition(&report, restPath, restContent, req.Verifier); err != nil {
		return err
	}
	stage = "publication"
	return evidence.WithReportWriter(ctx, ideaDir, func(writer *evidence.ReportWriter) error {
		if err := ctx.Err(); err != nil {
			return err
		}
		if _, _, err := verificationBindings(req, &report, ideaDir); err != nil {
			return err
		}
		updated, err := writer.SaveIfUnchanged(originalRaw, &report)
		if err != nil {
			return err
		}
		receipt.UpdatedSHA256 = sha256Hex(string(updated))
		// A cooperating writer cannot replace the report between its CAS and
		// receipt publication. Failure retains a report but never acceptance.
		if err := writeVerificationJSON(filepath.Join(filepath.Dir(requestAbs), "result.json"), receipt); err != nil {
			return err
		}
		receiptWritten = true
		return nil
	})
}

func (o driverImplOps) VerifyCompletionEvidence(ctx context.Context) (ok bool, detail string) {
	stage, requestSHA, originalSHA, invocation := "preflight", "", "", ""
	defer func() {
		if !ok {
			detail += "; " + o.retainDriverRefusal(stage, requestSHA, originalSHA, invocation)
		}
	}()
	fail := func(err error) (bool, string) { return false, err.Error() }
	// The helper command and criteria currently use the POSIX shell contract.
	// Cross-compilation does not establish an executable Windows contract.
	if runtime.GOOS == "windows" {
		return false, "independent evidence verification requires a POSIX execution host; Windows runtime is not supported"
	}
	if !o.base.Store.Enabled() {
		return false, "independent evidence verification requires a persistent driver run store"
	}
	if o.drafter == "" || o.drafter == o.implementer || o.base.RunID == "" {
		return false, "no independent verifier or run identity"
	}
	agent, err := agents.ResolveParticipant(o.drafter, o.base.Agents, rosterMappingFor(o.root))
	if err != nil {
		return false, "independent evidence verifier unavailable"
	}
	root, err := filepath.EvalSymlinks(o.root)
	if err != nil {
		return fail(err)
	}
	root, err = filepath.Abs(root)
	if err != nil {
		return fail(err)
	}
	ideaDir := filepath.Join(root, protocol.DeckDir, "ideas", o.ideaSlug)
	actualIdeaDir, err := filepath.EvalSymlinks(o.ideaDir)
	if err != nil || actualIdeaDir != ideaDir {
		return false, "idea directory does not match the canonical verification root"
	}
	o.root, o.ideaDir = root, ideaDir
	if err := requireCommittedRefusals(ctx, root, ideaDir); err != nil {
		return fail(err)
	}
	stage = "original-report"
	var original evidence.Report
	raw, err := readVerificationJSON(evidence.ReportPath(o.ideaDir), &original)
	if err != nil {
		return fail(err)
	}
	originalSHA = sha256Hex(string(raw))
	if original.CompletionTransition != nil {
		return false, "fresh report required before a new completion verification"
	}
	rest, _, err := implementationRestDigest(o.root, o.ideaDir)
	if err != nil {
		return fail(err)
	}
	req := evidenceVerificationRequest{Version: 1, Root: root, Idea: o.ideaSlug,
		RunID: o.base.RunID, Verifier: o.drafter, ReportSHA256: sha256Hex(string(raw)),
		TreeSHA256: original.TreeSHA256, RestSHA256: rest}
	stage = "bindings"
	if _, _, err := verificationBindings(req, &original, o.ideaDir); err != nil {
		return fail(err)
	}
	stage = "runtime"
	base := filepath.Join(root, ".parley-runtime", "evidence-verification")
	// Verify the ignore prerequisite before writing any retained runtime files.
	// Do not broaden the tested-tree exclusion to hide a configuration mistake.
	if _, err := exec.LookPath("git"); err != nil {
		return fail(fmt.Errorf("independent evidence verification requires Git: %w", err))
	}
	inGit, err := exec.CommandContext(ctx, "git", "-C", root, "rev-parse", "--is-inside-work-tree").Output()
	if err != nil || strings.TrimSpace(string(inGit)) != "true" {
		return false, "independent evidence verification requires a Git worktree with ignored runtime storage; non-Git closure is not supported"
	}
	if err := exec.CommandContext(ctx, "git", "-C", root, "check-ignore", "-q", "--", ".parley-runtime/").Run(); err != nil {
		return false, "verification runtime must be ignored before execution; add .parley-runtime/ to the repository ignore rules, rerun checks, then retry"
	}
	tracked, err := exec.CommandContext(ctx, "git", "-C", root, "ls-files", "-z", "--", ".parley-runtime/").Output()
	if err != nil {
		return fail(err)
	}
	if len(tracked) != 0 {
		return false, "verification runtime contains tracked files; restore private ignored runtime storage and rerun checks"
	}
	for _, path := range []string{filepath.Join(root, ".parley-runtime"), base} {
		info, err := os.Lstat(path)
		if err != nil && !os.IsNotExist(err) {
			return fail(err)
		}
		if err == nil && !info.IsDir() {
			return false, "verification runtime must use real directories inside the tested root"
		}
	}
	if err := os.MkdirAll(base, 0o700); err != nil {
		return fail(err)
	}
	dir, err := os.MkdirTemp(base, "attempt-")
	if err != nil {
		return fail(err)
	}
	stage = "request"
	req.AttemptID = filepath.Base(dir)
	requestPath := filepath.Join(dir, "request.json")
	if err := writeVerificationBytes(filepath.Join(dir, "original-report.json"), raw); err != nil {
		return fail(err)
	}
	if err := writeVerificationJSON(requestPath, req); err != nil {
		return fail(err)
	}
	requestRaw, err := os.ReadFile(requestPath)
	if err != nil {
		return fail(err)
	}
	requestSHA = sha256Hex(string(requestRaw))
	// Refuse if runtime creation itself changed the tested tree (e.g. runtime
	// paths are not ignored in this deck). Never expand exclusions to hide it.
	if _, _, err := verificationBindings(req, &original, o.ideaDir); err != nil {
		return fail(err)
	}
	executable := o.verificationCLI
	if executable == "" {
		executable, err = os.Executable()
		if err != nil {
			return fail(err)
		}
	}
	quote := func(s string) string { return "'" + strings.ReplaceAll(s, "'", "'\\''") + "'" }
	command := quote(executable) + " evidence verify --request " + quote(requestPath) + " --request-sha256 " + quote(requestSHA)
	prompt := fmt.Sprintf(`You are %s, the independent execution verifier for idea %s.
Execute the exact verifier command below once from this invocation. It reads the frozen named criteria, runs them, records actual independent executions and binds them to the original report and tested tree. Do not edit code, requests, reports, receipts, scope or signatures yourself. Do not substitute a written PASS, manually populated attribution, or commands selected from memory. If execution is unavailable or fails, report the actual failure. After the helper returns, print a short outcome and stop. Runtime attribution is not human authentication.
Verifier command: %s
`, o.drafter, o.ideaSlug, command)
	stage = "launch"
	ctx = runner.WithLaunchInfo(ctx, runner.LaunchInfo{RunID: o.base.RunID, Idea: o.ideaSlug, Phase: "evidence-verification", Store: o.base.Store})
	res := runner.RunConsult(ctx, runner.ConsultOptions{Root: root, Agent: agent, Prompt: prompt,
		Timeout: o.base.Timeout, StdoutPath: filepath.Join(dir, "agent.stdout.log"), StderrPath: filepath.Join(dir, "agent.stderr.log"), Progress: o.out})
	invocation = res.InvocationID
	if res.ExitError != "" || res.AgentExit != 0 || res.InvocationID == "" {
		return false, "independent verifier process failed or its launch was not observed"
	}
	stage = "acceptance"
	err = evidence.WithReportWriter(ctx, o.ideaDir, func(_ *evidence.ReportWriter) error {
		return o.acceptVerification(dir, req, original, requestPath, requestSHA, res.InvocationID)
	})
	if err != nil {
		return fail(err)
	}

	return true, "independent criterion execution and current-tree evidence passed"
}

// The accepted report is specific to this driver run. A competing verification
// or failed receipt write must not silently substitute another report at Complete.
type evidenceVerificationAcceptance struct {
	Version       int    `json:"version"`
	RunID         string `json:"run_id"`
	Verifier      string `json:"verifier"`
	ReportSHA256  string `json:"report_sha256"`
	RequestSHA256 string `json:"request_sha256"`
	RequestPath   string `json:"request_path"`
	ReceiptSHA256 string `json:"receipt_sha256"`
}

func (o driverImplOps) acceptVerification(dir string, req evidenceVerificationRequest, original evidence.Report, requestPath, requestSHA, invocationID string) error {
	var receipt evidenceVerificationReceipt
	receiptRaw, err := readVerificationJSON(filepath.Join(dir, "result.json"), &receipt)
	if err != nil {
		return fmt.Errorf("no valid helper execution receipt: %w", err)
	}
	if receipt.Version != 1 || receipt.RequestSHA256 != requestSHA || receipt.OriginalSHA256 != req.ReportSHA256 || receipt.Agent != o.drafter || receipt.RunID != o.base.RunID || receipt.ProcessMarker != invocationID || receipt.HelperPID <= 0 || receipt.HelperPID == os.Getpid() {
		return errors.New("independent helper receipt is not bound to this request and verifier")
	}
	if receipt.Error != "" {
		return errors.New(receipt.Error)
	}
	var updated evidence.Report
	updatedRaw, err := readVerificationJSON(evidence.ReportPath(o.ideaDir), &updated)
	if err != nil {
		return err
	}
	if receipt.UpdatedSHA256 != sha256Hex(string(updatedRaw)) || len(receipt.Executions) != len(original.Records) {
		return errors.New("helper receipt does not cover the persisted report and complete criterion scope")
	}
	// Independent execution may add only the verifier provenance. The original
	// criterion status, command, output, scope and tree observations stay intact.
	if len(updated.Records) != len(original.Records) {
		return errors.New("verifier changed original criterion scope")
	}
	restContent, restPath, err := implementationRestContent(o.root, o.ideaDir)
	if err != nil {
		return err
	}
	encodedUpdated, err := json.Marshal(updated)
	if err != nil {
		return err
	}
	var expectedTransition evidence.Report
	if err := json.Unmarshal(encodedUpdated, &expectedTransition); err != nil {
		return err
	}
	expectedTransition.CompletionTransition = nil
	if err := evidence.AuthorizeCompletionTransition(&expectedTransition, restPath, restContent, o.drafter); err != nil {
		return err
	}
	if !sameVerificationJSON(updated.CompletionTransition, expectedTransition.CompletionTransition) {
		return errors.New("verifier did not authorize the exact generated completion transition")
	}
	updated.CompletionTransition = original.CompletionTransition
	executions := make(map[string]evidence.CriterionRecord, len(receipt.Executions))
	for _, execution := range receipt.Executions {
		if _, duplicate := executions[execution.Name]; duplicate {
			return errors.New("duplicate criterion in helper receipt")
		}
		executions[execution.Name] = execution
	}
	originalRecords := make(map[string]evidence.CriterionRecord, len(original.Records))
	for _, record := range original.Records {
		originalRecords[record.Name] = record
	}
	for i := range updated.Records {
		originalRecord, known := originalRecords[updated.Records[i].Name]
		if !known {
			return errors.New("verifier added an unknown original criterion")
		}
		execution, found := executions[updated.Records[i].Name]
		retained := updated.Records[i].Provenance.VerifierRerun
		if !found || retained == nil || execution.Status != evidence.StatusPass || !sameVerificationJSON(retained.Command, execution.Command) || execution.Provenance.Executor != o.drafter {
			return errors.New("retained attestation differs from actual helper execution")
		}
		updated.Records[i].Provenance = originalRecord.Provenance
	}
	if !sameVerificationJSON(updated, original) {
		return errors.New("verifier changed the original evidence instead of only attesting it")
	}
	gate := o.EvidenceCloseGate(o.drafter)
	if !gate.Allowed {
		return errors.New(strings.Join(gate.Reasons, "; "))
	}
	accepted := evidenceVerificationAcceptance{Version: 1, RunID: o.base.RunID, Verifier: o.drafter, ReportSHA256: receipt.UpdatedSHA256, RequestSHA256: requestSHA, RequestPath: requestPath, ReceiptSHA256: sha256Hex(string(receiptRaw))}
	return writeVerificationJSON(filepath.Join(o.base.Store.Directory(), "evidence-accepted.json"), accepted)
}

// Called while holding the report writer guard, before the final gate/write.
func (o driverImplOps) requireAcceptedVerification() error {
	if err := requireCommittedRefusals(context.Background(), o.root, o.ideaDir); err != nil {
		return err
	}
	if !o.base.Store.Enabled() {
		return errors.New("completion requires persistent independent verification acceptance")
	}
	var accepted evidenceVerificationAcceptance
	if _, err := readVerificationJSON(filepath.Join(o.base.Store.Directory(), "evidence-accepted.json"), &accepted); err != nil {
		return fmt.Errorf("independent verification acceptance unavailable: %w", err)
	}
	if accepted.Version != 1 || accepted.RunID != o.base.RunID || accepted.Verifier != o.drafter {
		return errors.New("independent verification acceptance belongs to another run or verifier")
	}
	var report evidence.Report
	raw, err := readVerificationJSON(evidence.ReportPath(o.ideaDir), &report)
	if err != nil {
		return err
	}
	if sha256Hex(string(raw)) != accepted.ReportSHA256 {
		return errors.New("evidence report changed after driver acceptance")
	}
	var req evidenceVerificationRequest
	requestRaw, err := readVerificationJSON(accepted.RequestPath, &req)
	if err != nil {
		return err
	}
	if sha256Hex(string(requestRaw)) != accepted.RequestSHA256 || req.RunID != accepted.RunID || req.Verifier != accepted.Verifier {
		return errors.New("accepted verification request changed")
	}
	var receipt evidenceVerificationReceipt
	receiptRaw, err := readVerificationJSON(filepath.Join(filepath.Dir(accepted.RequestPath), "result.json"), &receipt)
	if err != nil {
		return fmt.Errorf("accepted verification receipt unavailable: %w", err)
	}
	if sha256Hex(string(receiptRaw)) != accepted.ReceiptSHA256 || receipt.Error != "" || receipt.UpdatedSHA256 != accepted.ReportSHA256 || receipt.RequestSHA256 != accepted.RequestSHA256 || receipt.RunID != accepted.RunID || receipt.Agent != accepted.Verifier {
		return errors.New("accepted verification receipt changed")
	}
	return nil
}

// Compare persisted JSON semantics instead of time.Time's in-memory location or
// monotonic internals. The records are structs/maps with deterministic encoding.
func sameVerificationJSON(a, b any) bool {
	left, e1 := json.Marshal(a)
	right, e2 := json.Marshal(b)
	return e1 == nil && e2 == nil && bytes.Equal(left, right)
}
