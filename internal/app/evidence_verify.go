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
	"path/filepath"
	"reflect"
	"strings"

	"parley-deck-cli/internal/agents"
	"parley-deck-cli/internal/driver"
	"parley-deck-cli/internal/evidence"
	"parley-deck-cli/internal/fsutil"
	"parley-deck-cli/internal/protocol"
	"parley-deck-cli/internal/runner"
)

// This runtime command is invoked BY the independently selected CLI agent.
// Inherited process markers provide attribution, not authentication against a
// malicious same-UID caller. The parent pins the request/report bytes and accepts
// only the actual helper's retained executions, never a model's textual PASS.
type evidenceVerificationRequest struct {
	Version      int    `json:"version"`
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

func verificationBindings(req evidenceVerificationRequest, report *evidence.Report, ideaDir string) ([]driver.CheckCriterion, error) {
	if report.Idea != req.Idea || report.TreeSHA256 != req.TreeSHA256 {
		return nil, errors.New("report idea or tested tree changed")
	}
	criteria, isList, err := driver.ReadChecksContract(ideaDir)
	if err != nil {
		return nil, err
	}
	if !isList || len(criteria) == 0 || len(criteria) != len(report.Records) {
		return nil, errors.New("named criterion scope is absent or incomplete")
	}
	byName := make(map[string]evidence.CriterionRecord, len(report.Records))
	for _, rec := range report.Records {
		if _, exists := byName[rec.Name]; exists {
			return nil, errors.New("duplicate criterion evidence")
		}
		byName[rec.Name] = rec
		if rec.Provenance.Executor == "" || rec.Provenance.Executor == req.Verifier {
			return nil, errors.New("missing executor or self-verification")
		}
	}
	for _, criterion := range criteria {
		rec, found := byName[criterion.Name]
		if !found || rec.Command.Command != criterion.Command || rec.Command.CommandSHA256 != sha256Hex(criterion.Command) {
			return nil, fmt.Errorf("criterion %q command or scope differs from the original execution", criterion.Name)
		}
	}
	exclusions, err := definedEvidenceArtifacts(req.Root, ideaDir)
	if err != nil {
		return nil, err
	}
	tree, err := evidence.TreeDigest(req.Root, exclusions...)
	if err != nil {
		return nil, err
	}
	if tree != req.TreeSHA256 {
		return nil, errors.New("tested tree changed before or during independent verification")
	}
	rest, rel, err := implementationRestDigest(req.Root, ideaDir)
	if err != nil {
		return nil, err
	}
	if rest != req.RestSHA256 || report.ExtraDigests[rel] != rest {
		return nil, errors.New("non-evidence implementation content changed")
	}
	return criteria, nil
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
	if req.RunID == "" || req.Verifier == "" || os.Getenv("PARLEY_RUN_ID") != req.RunID || os.Getenv("PARLEY_AGENT_ID") != req.Verifier || os.Getenv("PARLEY_PROC_MARKER") == "" {
		return errors.New("helper must run inside the selected independent verifier invocation")
	}
	receipt := evidenceVerificationReceipt{Version: 1, RequestSHA256: requestSHA,
		OriginalSHA256: req.ReportSHA256, Agent: os.Getenv("PARLEY_AGENT_ID"),
		RunID: os.Getenv("PARLEY_RUN_ID"), ProcessMarker: os.Getenv("PARLEY_PROC_MARKER"), HelperPID: os.Getpid()}
	defer func() {
		if resultErr != nil {
			receipt.Error = resultErr.Error()
		}
		resultErr = errors.Join(resultErr, writeVerificationJSON(filepath.Join(filepath.Dir(requestAbs), "result.json"), receipt))
	}()
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
	criteria, err := verificationBindings(req, &report, ideaDir)
	if err != nil {
		return err
	}
	for _, rec := range report.Records {
		if rec.Provenance.Verifier != "" || rec.Provenance.VerifierRerun != nil {
			return errors.New("fresh execution report required; an earlier verifier cannot be replayed")
		}
	}
	for _, criterion := range criteria {
		if err := ctx.Err(); err != nil {
			return err
		}
		if _, err := verificationBindings(req, &report, ideaDir); err != nil {
			return err
		}
		rec := evidence.RunCriterion(ctx, root, criterion.Name, criterion.Command, req.Verifier)
		receipt.Executions = append(receipt.Executions, rec)
		if _, err := verificationBindings(req, &report, ideaDir); err != nil {
			return err
		}
		if rec.Status != evidence.StatusPass {
			return fmt.Errorf("independent criterion %q is %s", criterion.Name, rec.Status)
		}
		rerun := evidence.VerifierExecution{Command: rec.Command, TreeBeforeSHA256: req.TreeSHA256, TreeAfterSHA256: req.TreeSHA256}
		if err := evidence.AttestExecution(&report, criterion.Name, req.Verifier, rerun); err != nil {
			return err
		}
	}
	if err := ctx.Err(); err != nil {
		return err
	}
	if _, err := verificationBindings(req, &report, ideaDir); err != nil {
		return err
	}
	// Detect a replacement observed before save. This snapshot check is not
	// atomic with Save; cooperative report-write serialization is still required.
	currentRaw, err := os.ReadFile(evidence.ReportPath(ideaDir))
	if err != nil || !bytes.Equal(originalRaw, currentRaw) {
		return errors.New("original evidence was replaced during verification")
	}
	if err := evidence.Save(ideaDir, &report); err != nil {
		return err
	}
	updated, err := os.ReadFile(evidence.ReportPath(ideaDir))
	if err != nil {
		return err
	}
	receipt.UpdatedSHA256 = sha256Hex(string(updated))
	return nil
}

func (o driverImplOps) VerifyCompletionEvidence(ctx context.Context) (bool, string) {
	fail := func(err error) (bool, string) { return false, err.Error() }
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
	var original evidence.Report
	raw, err := readVerificationJSON(evidence.ReportPath(o.ideaDir), &original)
	if err != nil {
		return fail(err)
	}
	rest, _, err := implementationRestDigest(o.root, o.ideaDir)
	if err != nil {
		return fail(err)
	}
	req := evidenceVerificationRequest{Version: 1, Root: root, Idea: o.ideaSlug,
		RunID: o.base.RunID, Verifier: o.drafter, ReportSHA256: sha256Hex(string(raw)),
		TreeSHA256: original.TreeSHA256, RestSHA256: rest}
	if _, err := verificationBindings(req, &original, o.ideaDir); err != nil {
		return fail(err)
	}
	base := filepath.Join(root, ".parley-runtime", "evidence-verification")
	if err := os.MkdirAll(base, 0o700); err != nil {
		return fail(err)
	}
	dir, err := os.MkdirTemp(base, "attempt-")
	if err != nil {
		return fail(err)
	}
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
	requestSHA := sha256Hex(string(requestRaw))
	// Refuse if runtime creation itself changed the tested tree (e.g. runtime
	// paths are not ignored in this deck). Never expand exclusions to hide it.
	if _, err := verificationBindings(req, &original, o.ideaDir); err != nil {
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
	ctx = runner.WithLaunchInfo(ctx, runner.LaunchInfo{RunID: o.base.RunID, Idea: o.ideaSlug, Phase: "evidence-verification", Store: o.base.Store})
	res := runner.RunConsult(ctx, runner.ConsultOptions{Root: root, Agent: agent, Prompt: prompt,
		Timeout: o.base.Timeout, StdoutPath: filepath.Join(dir, "agent.stdout.log"), StderrPath: filepath.Join(dir, "agent.stderr.log"), Progress: o.out})
	if res.ExitError != "" || res.AgentExit != 0 || res.InvocationID == "" {
		return false, "independent verifier process failed or its launch was not observed"
	}
	var receipt evidenceVerificationReceipt
	if _, err := readVerificationJSON(filepath.Join(dir, "result.json"), &receipt); err != nil {
		return fail(fmt.Errorf("no valid helper execution receipt: %w", err))
	}
	if receipt.Version != 1 || receipt.RequestSHA256 != requestSHA || receipt.OriginalSHA256 != req.ReportSHA256 || receipt.Agent != o.drafter || receipt.RunID != o.base.RunID || receipt.ProcessMarker != res.InvocationID || receipt.HelperPID <= 0 || receipt.HelperPID == os.Getpid() {
		return false, "independent helper receipt is not bound to this request and verifier"
	}
	if receipt.Error != "" {
		return false, receipt.Error
	}
	var updated evidence.Report
	updatedRaw, err := readVerificationJSON(evidence.ReportPath(o.ideaDir), &updated)
	if err != nil {
		return fail(err)
	}
	if receipt.UpdatedSHA256 != sha256Hex(string(updatedRaw)) || len(receipt.Executions) != len(original.Records) {
		return false, "helper receipt does not cover the persisted report and complete criterion scope"
	}
	// Independent execution may add only the verifier provenance. The original
	// criterion status, command, output, scope and tree observations stay intact.
	if len(updated.Records) != len(original.Records) {
		return false, "verifier changed original criterion scope"
	}
	for i := range updated.Records {
		retained := updated.Records[i].Provenance.VerifierRerun
		if retained == nil || receipt.Executions[i].Name != updated.Records[i].Name || receipt.Executions[i].Status != evidence.StatusPass || !reflect.DeepEqual(retained.Command, receipt.Executions[i].Command) || receipt.Executions[i].Provenance.Executor != o.drafter {
			return false, "retained attestation differs from actual helper execution"
		}
		updated.Records[i].Provenance = original.Records[i].Provenance
	}
	if !reflect.DeepEqual(updated, original) {
		return false, "verifier changed the original evidence instead of only attesting it"
	}
	gate := o.EvidenceCloseGate(o.drafter)
	if !gate.Allowed {
		return false, strings.Join(gate.Reasons, "; ")
	}
	return true, "independent criterion execution and current-tree evidence passed"
}
