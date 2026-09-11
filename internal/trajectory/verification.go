package trajectory

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"runtime"
	"time"

	"parley-deck-cli/internal/budget"
	"parley-deck-cli/internal/evidence"
	"parley-deck-cli/internal/fsutil"
	"parley-deck-cli/internal/procctl"
)

// VerificationTicket is private runtime data, not public telemetry. One ticket
// is retained per original charge, before launch. Changing run, verifier or
// request cannot create a silent retry of that charge's verification.
type VerificationTicket struct {
	Version int             `json:"version"`
	Root    string          `json:"root"`
	RunID   string          `json:"run_id"`
	Request CapturedRequest `json:"request"`
}

func (t VerificationTicket) SHA256() (string, error) {
	if t.Version != 1 || !filepath.IsAbs(t.Root) || filepath.Clean(t.Root) != t.Root || !safeLabel(t.RunID) {
		return "", errors.New("invalid captured verification ticket")
	}
	if _, err := t.Request.SHA256(); err != nil {
		return "", err
	}
	data, err := canonical(t)
	return digest(data), err
}

type verificationLaunch struct {
	Version      int       `json:"version"`
	TicketSHA256 string    `json:"ticket_sha256"`
	InvocationID string    `json:"invocation_id"`
	At           time.Time `json:"at"`
}

type verificationClaim struct {
	Version      int       `json:"version"`
	TicketSHA256 string    `json:"ticket_sha256"`
	LaunchSHA256 string    `json:"launch_sha256"`
	HelperPID    int       `json:"helper_pid"`
	At           time.Time `json:"at"`
}

type verificationStep struct {
	Version        int       `json:"version"`
	ClaimSHA256    string    `json:"claim_sha256"`
	Ordinal        int       `json:"ordinal"`
	PreviousSHA256 string    `json:"previous_sha256"`
	Execution      Execution `json:"execution"`
	ProcessSHA256  string    `json:"process_sha256,omitempty"`
}

type verificationProcess struct {
	Version     int             `json:"version"`
	ClaimSHA256 string          `json:"claim_sha256"`
	Ordinal     int             `json:"ordinal"`
	Identity    procctl.Spawned `json:"identity"`
}

type verificationStop struct {
	Version      int       `json:"version"`
	TicketSHA256 string    `json:"ticket_sha256"`
	LaunchSHA256 string    `json:"launch_sha256"`
	At           time.Time `json:"at"`
}

func verificationNotStopped(dir *os.Root) error {
	if _, err := dir.Lstat("stop.json"); os.IsNotExist(err) {
		return nil
	}
	return errors.New("captured verification is stopped or its stop state is unavailable")
}

func readVerificationProcess(dir *os.Root, claimSHA string, ordinal int) (verificationProcess, string, error) {
	var process verificationProcess
	sha, err := readVerificationArtifact(dir, fmt.Sprintf("process-%03d.json", ordinal), &process)
	if err != nil {
		return process, "", err
	}
	sp := process.Identity
	if process.Version != 1 || process.ClaimSHA256 != claimSHA || process.Ordinal != ordinal ||
		sp.PID <= 0 || sp.PGID != sp.PID || sp.BootID == "" || sp.ProcStart == "" || sp.Command == "" || sp.StartedAt.IsZero() {
		return process, "", errors.New("invalid captured criterion process identity")
	}
	return process, sha, nil
}

// StopCapturedVerification shares the exact guard with criterion creation and
// registration. The durable stop prevents any later start. Live groups are
// signalled only after strict attribution; an unrelated or reused PID refuses.
// Call before killing the enclosing CLI/helper group, with a cleanup context
// independent of the cancelled invocation. Repeated stops preserve the first.
func StopCapturedVerification(ctx context.Context, ticket VerificationTicket, invocation string) error {
	return withVerification(ctx, ticket, func(dir *os.Root, sha string) error {
		launchSHA, err := readVerificationLaunch(dir, sha, invocation)
		if err != nil {
			return err
		}
		var stop verificationStop
		if _, err = readVerificationArtifact(dir, "stop.json", &stop); os.IsNotExist(err) {
			stop = verificationStop{1, sha, launchSHA, time.Now().UTC()}
			_, err = writeVerificationArtifact(dir, "stop.json", stop)
		}
		if err != nil {
			return err
		}
		if stop.Version != 1 || stop.TicketSHA256 != sha || stop.LaunchSHA256 != launchSHA || stop.At.IsZero() {
			return errors.New("captured verification stop binding changed")
		}
		var claim verificationClaim
		claimSHA, err := readVerificationArtifact(dir, "claim.json", &claim)
		if os.IsNotExist(err) {
			return nil // stop won before the helper could claim or start anything
		}
		if err != nil || claim.Version != 2 || claim.TicketSHA256 != sha || claim.LaunchSHA256 != launchSHA || claim.HelperPID <= 0 {
			return errors.New("captured verification process control claim unavailable")
		}
		for ordinal := 1; ordinal <= 4*len(ticket.Request.Criteria); ordinal++ {
			process, _, err := readVerificationProcess(dir, claimSHA, ordinal)
			if os.IsNotExist(err) {
				continue
			}
			if err != nil {
				return err
			}
			if !procctl.Alive(process.Identity) {
				continue
			}
			killed, reason, err := procctl.KillTreeAttributed(process.Identity)
			if err != nil {
				return err
			}
			if !killed {
				return fmt.Errorf("captured criterion cleanup refused: %s", reason)
			}
		}
		return nil
	})
}

type verificationPrepared struct {
	Version     int    `json:"version"`
	ClaimSHA256 string `json:"claim_sha256"`
	BeforeRoot  string `json:"before_root"`
	AfterRoot   string `json:"after_root"`
}

// VerificationReceipt retains success or failure of the helper, not acceptance
// of a trajectory or authentication of a model. Failure is a bounded stage label;
// raw errors and command output are never copied into receipt metadata.
type VerificationReceipt struct {
	Version        int       `json:"version"`
	TicketSHA256   string    `json:"ticket_sha256"`
	LaunchSHA256   string    `json:"launch_sha256"`
	ClaimSHA256    string    `json:"claim_sha256"`
	PreparedSHA256 string    `json:"prepared_sha256"`
	InvocationID   string    `json:"invocation_id"`
	HelperPID      int       `json:"helper_pid"`
	BeforeRoot     string    `json:"before_root"`
	AfterRoot      string    `json:"after_root"`
	Steps          int       `json:"steps"`
	LastSHA256     string    `json:"last_sha256"`
	FailureStage   string    `json:"failure_stage"`
	FinishedAt     time.Time `json:"finished_at"`
}

func syncVerificationDirectory(dir *os.Root) error {
	f, err := dir.Open(".")
	if err != nil {
		return err
	}
	defer f.Close()
	return fsutil.SyncFile(f)
}

// Exclusive writes leave any partial file in place. Its existence prevents
// replay; a torn/noncanonical artifact is an unresolved failure, never a retry.
func writeVerificationArtifact(dir *os.Root, name string, value any) (string, error) {
	data, err := canonical(value)
	if err != nil {
		return "", err
	}
	if len(data) > 1<<20 {
		return "", errors.New("verification artifact exceeds its bound")
	}
	f, err := dir.OpenFile(name, os.O_WRONLY|os.O_CREATE|os.O_EXCL, 0600)
	if err != nil {
		return "", err
	}
	_, writeErr := f.Write(data)
	err = errors.Join(writeErr, fsutil.SyncFile(f), f.Close())
	if err != nil {
		return "", err
	}
	if err = syncVerificationDirectory(dir); err != nil {
		return "", err
	}
	return digest(data), nil
}

func readVerificationArtifact(dir *os.Root, name string, value any) (string, error) {
	info, err := dir.Lstat(name)
	if err != nil {
		return "", err
	}
	if !info.Mode().IsRegular() || info.Size() > 1<<20 {
		return "", errors.New("verification artifact must be a bounded regular file")
	}
	f, err := dir.Open(name)
	if err != nil {
		return "", err
	}
	defer f.Close()
	opened, err := f.Stat()
	if err != nil || !os.SameFile(info, opened) {
		return "", errors.New("verification artifact changed during open")
	}
	data, err := io.ReadAll(io.LimitReader(f, (1<<20)+1))
	if err != nil || len(data) > 1<<20 {
		return "", errors.New("verification artifact cannot be read within its bound")
	}
	dec := json.NewDecoder(bytes.NewReader(data))
	dec.DisallowUnknownFields()
	if err = dec.Decode(value); err != nil {
		return "", err
	}
	expected, err := canonical(value)
	if err != nil || !bytes.Equal(data, expected) {
		return "", errors.New("verification artifact is incomplete, ambiguous or noncanonical")
	}
	return digest(data), nil
}

func openVerificationDirectory(b budget.CycleBinding, charge string, create bool) (*os.Root, error) {
	if runtime.GOOS == "windows" {
		return nil, errors.New("captured verification journals require a POSIX execution host")
	}
	if !validHash(charge) {
		return nil, errors.New("invalid verification charge identity")
	}
	parent, err := os.OpenRoot(filepath.Dir(b.Store.Dir))
	if err != nil {
		return nil, err
	}
	defer parent.Close()
	const name = "trajectory-verifications"
	if create {
		if err = parent.Mkdir(name, 0700); err != nil && !os.IsExist(err) {
			return nil, err
		}
		if err = syncVerificationDirectory(parent); err != nil {
			return nil, err
		}
	}
	info, err := parent.Lstat(name)
	if err != nil || !info.IsDir() {
		return nil, errors.New("verification storage requires a real directory")
	}
	base, err := parent.OpenRoot(name)
	if err != nil {
		return nil, err
	}
	defer base.Close()
	if create {
		// Never remove or reuse an existing reservation, even if request writing
		// was interrupted. Explicit recovery must account for that history.
		if err = base.Mkdir(charge, 0700); err != nil {
			return nil, errors.New("verification already reserved or unavailable; preserve it for explicit recovery")
		}
		if err = syncVerificationDirectory(base); err != nil {
			return nil, err
		}
	}
	info, err = base.Lstat(charge)
	if err != nil || !info.IsDir() {
		return nil, errors.New("verification reservation requires a real directory")
	}
	return base.OpenRoot(charge)
}

// PrepareCapturedVerification is called before selecting a process invocation.
// A successful return means the exact ticket reached the file and directory
// persistence barriers. It grants no permission to repeat a model call.
func PrepareCapturedVerification(ctx context.Context, root, idea, verifier, runID string) (VerificationTicket, error) {
	var ticket VerificationTicket
	root, err := filepath.Abs(root)
	if err != nil {
		return ticket, err
	}
	root, err = filepath.EvalSymlinks(root)
	if err != nil {
		return ticket, err
	}
	err = withState(ctx, root, idea, func(b budget.CycleBinding, _ budget.Snapshot, s State) error {
		r, err := capturedRequest(s, verifier)
		if err != nil {
			return err
		}
		ticket = VerificationTicket{Version: 1, Root: root, RunID: runID, Request: r}
		if _, err = ticket.SHA256(); err != nil {
			return err
		}
		dir, err := openVerificationDirectory(b, r.Charge.EntryKey, true)
		if err != nil {
			return err
		}
		defer dir.Close()
		_, err = writeVerificationArtifact(dir, "request.json", ticket)
		return err
	})
	if err == nil {
		_, err = ticket.SHA256()
	}
	if err != nil {
		return VerificationTicket{}, err
	}
	return ticket, nil
}

// withVerification rechecks the complete charge/state/archive authority and the
// exact prepared ticket under the shared cycle guard. Caller-provided metadata
// cannot substitute an unpersisted request or a different origin worktree.
func withVerification(ctx context.Context, ticket VerificationTicket, fn func(*os.Root, string) error) error {
	want, err := ticket.SHA256()
	if err != nil {
		return err
	}
	root, err := filepath.EvalSymlinks(ticket.Root)
	if err != nil || root != ticket.Root {
		return errors.New("verification origin is unavailable or no longer canonical")
	}
	observed := false
	err = withState(ctx, root, ticket.Request.Idea, func(b budget.CycleBinding, _ budget.Snapshot, s State) error {
		observed = true
		r, err := capturedRequest(s, ticket.Request.Verifier)
		if err != nil {
			return err
		}
		current := ticket
		current.Request = r
		actual, err := current.SHA256()
		if err != nil || actual != want {
			return errors.New("prepared verification original authority changed")
		}
		dir, err := openVerificationDirectory(b, r.Charge.EntryKey, false)
		if err != nil {
			return err
		}
		defer dir.Close()
		var retained VerificationTicket
		actual, err = readVerificationArtifact(dir, "request.json", &retained)
		if err != nil || actual != want {
			return errors.New("prepared verification ticket changed or disappeared")
		}
		return fn(dir, want)
	})
	if err == nil && !observed {
		return errors.New("prepared verification authority disappeared")
	}
	return err
}

// ReserveCapturedVerificationLaunch must be called by the instrumented runner
// with its actual reserved invocation ID before spawning the selected verifier.
// This API does not authenticate that ID. Integration must enforce its lineage.
func ReserveCapturedVerificationLaunch(ctx context.Context, ticket VerificationTicket, invocation string) error {
	if !safeLabel(invocation) || invocation == ticket.Request.InvocationID {
		return errors.New("independent verifier requires a distinct invocation")
	}
	return withVerification(ctx, ticket, func(dir *os.Root, sha string) error {
		_, err := writeVerificationArtifact(dir, "launch.json", verificationLaunch{1, sha, invocation, time.Now().UTC()})
		return err
	})
}

func readVerificationLaunch(dir *os.Root, ticketSHA, invocation string) (string, error) {
	var launch verificationLaunch
	sha, err := readVerificationArtifact(dir, "launch.json", &launch)
	if err != nil {
		return "", err
	}
	if launch.Version != 1 || launch.TicketSHA256 != ticketSHA || launch.InvocationID != invocation || !safeLabel(invocation) || launch.At.IsZero() {
		return "", errors.New("verification invocation differs from its durable reservation")
	}
	return sha, nil
}

// ExecuteCapturedVerification is the execution half of the independent helper.
// The orchestration layer must authenticate its inherited process attribution.
// A durable claim precedes preparation/commands and is never removed. On crash,
// completed step records survive and no fresh process can replay this ticket.
func ExecuteCapturedVerification(ctx context.Context, ticket VerificationTicket, invocation string, criteria []Criterion, parent string) (receipt VerificationReceipt, resultErr error) {
	var retained *os.Root
	stage := "claim"
	resultErr = withVerification(ctx, ticket, func(dir *os.Root, sha string) error {
		if err := verificationNotStopped(dir); err != nil {
			return err
		}
		launchSHA, err := readVerificationLaunch(dir, sha, invocation)
		if err != nil {
			return err
		}
		claim := verificationClaim{2, sha, launchSHA, os.Getpid(), time.Now().UTC()}
		claimSHA, err := writeVerificationArtifact(dir, "claim.json", claim)
		if err != nil {
			return errors.New("verification helper already claimed or claim unavailable; no replay is allowed")
		}
		retained, err = dir.OpenRoot(".")
		if err != nil {
			return err
		}
		receipt = VerificationReceipt{Version: 2, TicketSHA256: sha, LaunchSHA256: launchSHA,
			ClaimSHA256: claimSHA, InvocationID: invocation, HelperPID: claim.HelperPID}
		return nil
	})
	if resultErr != nil {
		return receipt, resultErr
	}
	defer retained.Close()
	defer func() {
		if resultErr != nil {
			receipt.FailureStage = stage
		}
		receipt.FinishedAt = time.Now().UTC()
		_, err := writeVerificationArtifact(retained, "receipt.json", receipt)
		resultErr = errors.Join(resultErr, err)
	}()
	checkAuthority := func(dir *os.Root, sha string) error {
		if err := verificationNotStopped(dir); err != nil {
			return err
		}
		launchSHA, err := readVerificationLaunch(dir, sha, invocation)
		if err != nil || launchSHA != receipt.LaunchSHA256 {
			return errors.New("verification launch authority changed")
		}
		var claim verificationClaim
		claimSHA, err := readVerificationArtifact(dir, "claim.json", &claim)
		if err != nil || claimSHA != receipt.ClaimSHA256 {
			return errors.New("verification helper claim changed")
		}
		var prepared verificationPrepared
		preparedSHA, err := readVerificationArtifact(dir, "prepared.json", &prepared)
		if err != nil || preparedSHA != receipt.PreparedSHA256 {
			return errors.New("verification source preparation record changed")
		}
		return nil
	}
	guard := func() error {
		return withVerification(ctx, ticket, checkAuthority)
	}
	stage = "preparation"
	w, err := OpenCaptured(ctx, ticket.Root, ticket.Request, parent)
	if err != nil {
		return receipt, err
	}
	defer func() { resultErr = errors.Join(resultErr, w.Close()) }()
	receipt.BeforeRoot, receipt.AfterRoot = w.BeforeRoot(), w.AfterRoot()
	receipt.PreparedSHA256, err = writeVerificationArtifact(retained, "prepared.json", verificationPrepared{1, receipt.ClaimSHA256, w.BeforeRoot(), w.AfterRoot()})
	if err != nil {
		return receipt, err
	}
	stage = "execution"
	processSHA := ""
	control := func(ordinal int) evidence.CriterionStartControl {
		return func(start func() (procctl.Spawned, error), release func() error) error {
			processSHA = ""
			return withVerification(ctx, ticket, func(dir *os.Root, sha string) error {
				if err := checkAuthority(dir, sha); err != nil {
					return err
				}
				sp, err := start()
				if err != nil {
					return err
				}
				processSHA, err = writeVerificationArtifact(dir, fmt.Sprintf("process-%03d.json", ordinal), verificationProcess{1, receipt.ClaimSHA256, ordinal, sp})
				if err != nil {
					return err
				}
				return release()
			})
		}
	}
	_, resultErr = verifyCaptured(ctx, w, criteria, guard, func(ordinal int, execution Execution) error {
		step := verificationStep{2, receipt.ClaimSHA256, ordinal, receipt.LastSHA256, execution, processSHA}
		sha, err := writeVerificationArtifact(retained, fmt.Sprintf("step-%03d.json", ordinal), step)
		if err == nil {
			receipt.Steps, receipt.LastSHA256 = ordinal, sha
		}
		return err
	}, control)
	if resultErr == nil {
		stage = "authority"
		resultErr = guard()
	}
	return receipt, resultErr
}

// ReadCapturedVerification recovers retained partial observations without
// executing anything. Missing terminal receipts remain explicit errors. A
// caller must also validate actual runner terminal/identity before using the
// returned observations; this reader never resolves a trajectory attempt.
func ReadCapturedVerification(ctx context.Context, ticket VerificationTicket, invocation string) (receipt VerificationReceipt, observation Observation, resultErr error) {
	requestSHA, err := ticket.Request.SHA256()
	if err != nil {
		return receipt, observation, err
	}
	observation = Observation{Version: 1, RequestSHA256: requestSHA, Verifier: ticket.Request.Verifier}
	resultErr = withVerification(ctx, ticket, func(dir *os.Root, ticketSHA string) error {
		launchSHA, err := readVerificationLaunch(dir, ticketSHA, invocation)
		if err != nil {
			return err
		}
		var claim verificationClaim
		claimSHA, err := readVerificationArtifact(dir, "claim.json", &claim)
		if err != nil {
			return err
		}
		if (claim.Version != 1 && claim.Version != 2) || claim.TicketSHA256 != ticketSHA || claim.LaunchSHA256 != launchSHA || claim.HelperPID <= 0 || claim.At.IsZero() {
			return errors.New("invalid verification helper claim")
		}
		previous := ""
		steps := 0
		for ordinal := 1; ordinal <= 4*len(ticket.Request.Criteria); ordinal++ {
			var step verificationStep
			sha, err := readVerificationArtifact(dir, fmt.Sprintf("step-%03d.json", ordinal), &step)
			if os.IsNotExist(err) {
				break
			}
			if err != nil {
				return err
			}
			if step.Version != claim.Version || step.Ordinal != ordinal || step.ClaimSHA256 != claimSHA || step.PreviousSHA256 != previous {
				return errors.New("verification execution journal has a gap or changed lineage")
			}
			if claim.Version == 1 && step.ProcessSHA256 != "" {
				return errors.New("legacy verification contains mixed process-control history")
			}
			if claim.Version == 2 && (step.ProcessSHA256 != "" || step.Execution.Complete) {
				_, processSHA, err := readVerificationProcess(dir, claimSHA, ordinal)
				if err != nil || processSHA != step.ProcessSHA256 {
					return errors.New("verification execution lost its registered process")
				}
			}
			index, offset := (ordinal-1)/4, (ordinal-1)%4
			if offset == 0 {
				observation.Pairs = append(observation.Pairs, Pair{Name: ticket.Request.Criteria[index].Name})
			}
			pair := &observation.Pairs[index]
			switch offset {
			case 0:
				pair.Before[0] = step.Execution
			case 1:
				pair.After[0] = step.Execution
			case 2:
				pair.After[1] = step.Execution
			case 3:
				pair.Before[1] = step.Execution
			}
			steps, previous = ordinal, sha
		}
		// Refuse extra or skipped numbered observations, including files after
		// the first missing step. Do not hide ambiguous history behind a receipt.
		entries, err := dir.Open(".")
		if err != nil {
			return err
		}
		names, readErr := entries.Readdirnames(7 + 8*MaxCriteria)
		closeErr := entries.Close()
		if errors.Is(readErr, io.EOF) {
			readErr = nil
		}
		if err = errors.Join(readErr, closeErr); err != nil {
			return err
		}
		allowed := map[string]bool{"request.json": true, "launch.json": true, "claim.json": true, "prepared.json": true, "receipt.json": true, "stop.json": true}
		for i := 1; i <= steps; i++ {
			allowed[fmt.Sprintf("step-%03d.json", i)] = true
		}
		if claim.Version == 2 {
			for i := 1; i <= steps+1 && i <= 4*len(ticket.Request.Criteria); i++ {
				allowed[fmt.Sprintf("process-%03d.json", i)] = true
			}
		}
		if len(names) > 6+8*MaxCriteria {
			return errors.New("verification journal exceeds its inventory bound")
		}
		for _, name := range names {
			if !allowed[name] {
				return errors.New("verification journal contains unexpected or out-of-order artifacts")
			}
			for ordinal := 1; ordinal <= steps+1 && ordinal <= 4*len(ticket.Request.Criteria); ordinal++ {
				if name == fmt.Sprintf("process-%03d.json", ordinal) {
					if _, _, err := readVerificationProcess(dir, claimSHA, ordinal); err != nil {
						return err
					}
				}
			}
		}
		if _, err = readVerificationArtifact(dir, "receipt.json", &receipt); err != nil {
			return fmt.Errorf("verification terminal receipt unavailable; retain partial observations: %w", err)
		}
		if receipt.Version != claim.Version || receipt.TicketSHA256 != ticketSHA || receipt.LaunchSHA256 != launchSHA || receipt.ClaimSHA256 != claimSHA || receipt.InvocationID != invocation || receipt.HelperPID != claim.HelperPID || receipt.FinishedAt.Before(claim.At) || receipt.Steps != steps || receipt.LastSHA256 != previous {
			return errors.New("verification receipt differs from the retained execution journal")
		}
		if receipt.FailureStage != "" {
			return errors.New("verification retained a helper failure; no acceptance or retry is authorized")
		}
		if err := verificationNotStopped(dir); err != nil {
			return err
		}
		if !filepath.IsAbs(receipt.BeforeRoot) || !filepath.IsAbs(receipt.AfterRoot) || receipt.BeforeRoot == receipt.AfterRoot || filepath.Dir(receipt.BeforeRoot) != filepath.Dir(receipt.AfterRoot) {
			return errors.New("verification receipt lacks its distinct prepared source roots")
		}
		var prepared verificationPrepared
		preparedSHA, err := readVerificationArtifact(dir, "prepared.json", &prepared)
		if err != nil || preparedSHA != receipt.PreparedSHA256 || prepared.Version != 1 || prepared.ClaimSHA256 != claimSHA || prepared.BeforeRoot != receipt.BeforeRoot || prepared.AfterRoot != receipt.AfterRoot {
			return errors.New("verification receipt differs from the original prepared source roots")
		}
		_, err = AssessCaptured(ticket.Request, observation)
		return err
	})
	return receipt, observation, resultErr
}
