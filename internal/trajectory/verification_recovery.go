package trajectory

import (
	"context"
	"errors"
	"io"
	"os"
	"path/filepath"
	"time"

	"parley-deck-cli/internal/telemetry"
)

// The runner reserves a captured-verification launch before its budget check
// (internal/runner/telemetry.go beginLaunch). A denied verifier launch therefore
// leaves the ticket pinned to an invocation that never started: request.json and
// launch.json are retained, no helper claim can exist, and every pre-recovery
// public path refuses the reservation. Recovery below is the explicit repair of
// exactly that state. It is never invoked by a launch path, executes no model
// and no criterion, and spends, refunds or resets nothing: the refused
// invocation's records, the launch reservation and every budget counter stay
// byte-for-byte as the refusal left them. The subsequent verifier launch is
// separately checked against the recovery binding and separately charged by its
// own budget boundary.

const capturedVerificationPhaseLabel = "trajectory-verification"

var errSupersededVerificationInvocation = errors.New("verification invocation was superseded by an explicit recovery; it retains no control over the recovered launch")

// CapturedRecoveryIdentity is the exact validated immutable recovery authority:
// the single intended invocation the retained recovery binds, the refused
// invocation it repairs, the bound refused-launch digests, the prior launch
// reservation digest, the recovery artifact's own digest and its timestamp.
// Every field was revalidated against the live journal and the retained
// refused-launch bytes by the common authority reader; nothing here is trusted
// from recovery.json alone.
type CapturedRecoveryIdentity struct {
	InvocationID           string    `json:"invocation_id"`
	RefusedInvocationID    string    `json:"refused_invocation_id"`
	RefusedRequestedSHA256 string    `json:"refused_requested_sha256"`
	RefusedTerminalSHA256  string    `json:"refused_terminal_sha256"`
	PriorLaunchSHA256      string    `json:"prior_launch_sha256"`
	RecoverySHA256         string    `json:"recovery_sha256"`
	At                     time.Time `json:"at"`
}

func capturedRecoveryIdentity(effective verificationEffectiveLaunch) CapturedRecoveryIdentity {
	return CapturedRecoveryIdentity{
		InvocationID:           effective.Recovery.InvocationID,
		RefusedInvocationID:    effective.Recovery.RefusedInvocationID,
		RefusedRequestedSHA256: effective.Recovery.RefusedRequestedSHA256,
		RefusedTerminalSHA256:  effective.Recovery.RefusedTerminalSHA256,
		PriorLaunchSHA256:      effective.Recovery.PriorLaunchSHA256,
		RecoverySHA256:         effective.AuthoritySHA256,
		At:                     effective.Recovery.At,
	}
}

// ReadCapturedVerificationRecovery is a genuinely READ-ONLY resolution of the
// immutable recovery authority for a captured-verification ticket. It runs
// under the same guarded complete authority as every other ticket handle
// (charge/state/archive, quorum and the exact retained ticket), re-reads the
// actual launch reservation, and resolves the recovery through the common
// authority reader — so the recovery's binding to the live launch.json bytes,
// the full refused-launch evidence with its bound digests, and the
// refusal/recovery chronology are all revalidated against the bytes retained
// now. It performs no reservation and no publication: nothing is written,
// created, removed or reset, and a missing reservation or missing recovery is
// a refusal — never a recreated authority. Callers compare the returned
// identity; they must not treat file presence alone as authority.
func ReadCapturedVerificationRecovery(ctx context.Context, ticket VerificationTicket) (CapturedRecoveryIdentity, error) {
	var identity CapturedRecoveryIdentity
	err := withVerification(ctx, ticket, func(dir *os.Root, sha string) error {
		var launch verificationLaunch
		_, err := readVerificationArtifact(dir, "launch.json", &launch)
		if os.IsNotExist(err) {
			return errors.New("captured verification recovery requires a retained launch reservation and immutable recovery record")
		}
		if err != nil {
			return err
		}
		if launch.Version != 1 || launch.TicketSHA256 != sha || !runtimeID(launch.InvocationID) || launch.At.IsZero() {
			return errors.New("verification launch reservation is not a valid original authority")
		}
		effective, err := resolveVerificationEffectiveLaunch(dir, ticket, launch)
		if err != nil {
			return err
		}
		if !effective.Recovered {
			return errors.New("captured verification recovery requires a retained launch reservation and immutable recovery record")
		}
		identity = capturedRecoveryIdentity(effective)
		return nil
	})
	return identity, err
}

// CapturedRecoveryPreview is the exact inspected recovery authority an attended
// apply must reproduce before any write: the ticket/run identity, the original
// refused invocation with its bound requested/terminal digests and refusal
// completion, and the selected replacement invocation. The preview grants no
// process, budget or resolution; it only describes the single immutable
// recovery the apply may write.
type CapturedRecoveryPreview struct {
	Version                int       `json:"version"`
	Root                   string    `json:"root"`
	Idea                   string    `json:"idea"`
	RunID                  string    `json:"run_id"`
	TicketSHA256           string    `json:"ticket_sha256"`
	RefusedInvocationID    string    `json:"refused_invocation_id"`
	RefusedRequestedSHA256 string    `json:"refused_requested_sha256"`
	RefusedTerminalSHA256  string    `json:"refused_terminal_sha256"`
	RefusedCompletedAt     time.Time `json:"refused_completed_at"`
	InvocationID           string    `json:"invocation_id"`
}

func (p CapturedRecoveryPreview) SHA256() string { raw, _ := canonical(p); return digest(raw) }

// capturedVerificationRecoveryPreview derives the exact recovery preview from
// live evidence under the full guarded authority. Without a retained recovery
// the journal must be exactly the retained refusal (prepared request plus its
// refused launch reservation) — the same admission discipline as
// RecoverCapturedVerificationLaunch. With a retained recovery the derivation is
// an exact replay: the immutable record must bind these same facts, otherwise
// the preview refuses. The bool reports whether the recovery was already
// applied (a replayed preview writes nothing).
func capturedVerificationRecoveryPreview(ctx context.Context, ticket VerificationTicket, replacement string) (CapturedRecoveryPreview, bool, error) {
	var p CapturedRecoveryPreview
	applied := false
	if !safeLabel(replacement) || replacement == ticket.Request.InvocationID {
		return p, false, errors.New("verification recovery requires a distinct intended verifier invocation")
	}
	err := withVerification(ctx, ticket, func(dir *os.Root, sha string) error {
		var launch verificationLaunch
		launchSHA, err := readVerificationArtifact(dir, "launch.json", &launch)
		if err != nil {
			return errors.New("verification recovery requires the original launch reservation")
		}
		if launch.Version != 1 || launch.TicketSHA256 != sha || !runtimeID(launch.InvocationID) || launch.At.IsZero() {
			return errors.New("verification recovery contradicts the original launch reservation")
		}
		refused := launch.InvocationID
		if replacement == refused {
			return errors.New("verification recovery requires a distinct intended verifier invocation")
		}
		terminalSHA, requestedSHA, completedAt, err := refusedVerificationEvidence(ticket, refused, launch.At)
		if err != nil {
			return err
		}
		p = CapturedRecoveryPreview{Version: 1, Root: ticket.Root, Idea: ticket.Request.Idea, RunID: ticket.RunID,
			TicketSHA256: sha, RefusedInvocationID: refused, RefusedRequestedSHA256: requestedSHA,
			RefusedTerminalSHA256: terminalSHA, RefusedCompletedAt: completedAt, InvocationID: replacement}
		var recovery verificationRecovery
		_, err = readVerificationArtifact(dir, "recovery.json", &recovery)
		if os.IsNotExist(err) {
			// Fresh preview: the journal must be exactly the retained refusal.
			return unrecoveredVerificationJournal(dir)
		}
		if err != nil {
			return err
		}
		// Retained replay: the immutable recovery must bind exactly these facts.
		if recovery.Version != 1 || recovery.TicketSHA256 != sha || recovery.PriorLaunchSHA256 != launchSHA ||
			recovery.RefusedInvocationID != refused || recovery.RefusedRequestedSHA256 != requestedSHA ||
			recovery.RefusedTerminalSHA256 != terminalSHA || recovery.InvocationID != replacement ||
			recovery.At.IsZero() || recovery.At.Before(launch.At) || recovery.At.Before(completedAt) {
			return errors.New("verification recovery replay contradicts the retained immutable recovery")
		}
		applied = true
		return nil
	})
	return p, applied, err
}

// PreviewCapturedVerificationRecovery inspects the exact recovery an attended
// apply would write for the selected replacement invocation. It writes nothing
// and launches nothing; on a ticket whose recovery is already retained it
// replays only for the recovery's bound replacement.
func PreviewCapturedVerificationRecovery(ctx context.Context, ticket VerificationTicket, replacement string) (CapturedRecoveryPreview, error) {
	p, _, err := capturedVerificationRecoveryPreview(ctx, ticket, replacement)
	return p, err
}

// ApplyCapturedVerificationRecovery rechecks the exact inspected preview
// against live evidence and then performs the single immutable recovery write.
// Any drift since the preview — changed state, changed refused evidence, a
// different replacement — refuses before any write. An exact replay against the
// retained recovery (same replacement, same expected digest) returns the
// retained preview without writing; a conflicting replacement is refused. The
// admission itself (RecoverCapturedVerificationLaunch) revalidates everything
// again under its own guard, so a race between preview and apply fails closed.
// Recovery executes no model, no criterion and no budget action.
func ApplyCapturedVerificationRecovery(ctx context.Context, ticket VerificationTicket, replacement, expected string) (CapturedRecoveryPreview, error) {
	p, applied, err := capturedVerificationRecoveryPreview(ctx, ticket, replacement)
	if err != nil {
		return CapturedRecoveryPreview{}, err
	}
	if !validHash(expected) || expected != p.SHA256() {
		return CapturedRecoveryPreview{}, errors.New("verification recovery evidence or state changed since preview; inspect again")
	}
	if applied {
		return p, nil
	}
	if err = RecoverCapturedVerificationLaunch(ctx, ticket, p.RefusedInvocationID, p.InvocationID); err != nil {
		return CapturedRecoveryPreview{}, err
	}
	return p, nil
}

// verificationRecovery is immutable once written. launch.json stays the original
// reservation, byte-identical and pinned to the refused invocation; this
// artifact binds the exact retained bytes plus the single intended invocation a
// later verifier launch must carry. Replay is deliberately unsupported: a
// repeated or conflicting apply fails closed on the retained artifact instead of
// rewriting it, and no expected-preview authority exists for an exact-replay
// variant (adopting one would be a separate reviewed change).
type verificationRecovery struct {
	Version                int       `json:"version"`
	TicketSHA256           string    `json:"ticket_sha256"`
	PriorLaunchSHA256      string    `json:"prior_launch_sha256"`
	RefusedInvocationID    string    `json:"refused_invocation_id"`
	RefusedRequestedSHA256 string    `json:"refused_requested_sha256"`
	RefusedTerminalSHA256  string    `json:"refused_terminal_sha256"`
	InvocationID           string    `json:"invocation_id"`
	At                     time.Time `json:"at"`
}

// RecoverCapturedVerificationLaunch repairs one pre-start budget-refused
// verifier launch reservation. Admission requires the exact retained evidence:
// the original canonical request.json under the full charge/state authority, the
// launch reservation pinned to refusedInvocation, an unexecuted and unstopped
// journal, and the refused invocation's requested and terminal records proving
// a pre-start budget refusal (no started lifecycle, no PID, no exit code, the
// verifier launch's own phase binding). On success exactly one immutable
// recovery.json is written, binding those bytes plus newInvocation; nothing else
// is written, removed or reset. Every contradiction — a started lifecycle, an
// exit code, a foreign failure class or metadata, helper history, a stop, a
// wrong or missing launch binding, or an existing recovery — refuses before any
// write.
func RecoverCapturedVerificationLaunch(ctx context.Context, ticket VerificationTicket, refusedInvocation, newInvocation string) error {
	if !runtimeID(refusedInvocation) || !safeLabel(newInvocation) || newInvocation == refusedInvocation ||
		refusedInvocation == ticket.Request.InvocationID || newInvocation == ticket.Request.InvocationID {
		return errors.New("verification recovery requires the refused launch and a distinct intended verifier invocation")
	}
	return withVerification(ctx, ticket, func(dir *os.Root, sha string) error {
		if _, err := dir.Lstat("recovery.json"); !os.IsNotExist(err) {
			return errors.New("verification launch was already recovered or its recovery state is unavailable; the retained artifact stays authoritative")
		}
		var launch verificationLaunch
		launchSHA, err := readVerificationArtifact(dir, "launch.json", &launch)
		if err != nil {
			return errors.New("verification recovery requires the original launch reservation")
		}
		if launch.Version != 1 || launch.TicketSHA256 != sha || launch.InvocationID != refusedInvocation || launch.At.IsZero() {
			return errors.New("verification recovery contradicts the original launch reservation")
		}
		if err = unrecoveredVerificationJournal(dir); err != nil {
			return err
		}
		terminalSHA, requestedSHA, _, err := refusedVerificationEvidence(ticket, refusedInvocation, launch.At)
		if err != nil {
			return err
		}
		_, err = writeVerificationArtifact(dir, "recovery.json", verificationRecovery{Version: 1, TicketSHA256: sha,
			PriorLaunchSHA256: launchSHA, RefusedInvocationID: refusedInvocation,
			RefusedRequestedSHA256: requestedSHA, RefusedTerminalSHA256: terminalSHA,
			InvocationID: newInvocation, At: time.Now().UTC()})
		return err
	})
}

// A recovered launch is never a retry of verification work. Admission requires
// the journal to be exactly the retained refusal — the prepared request plus its
// refused launch reservation. Any claim, receipt, prepared record, stop, step,
// process or unknown record — including a stop written before recovery — keeps
// the reservation bricked instead of being overwritten or reused.
func unrecoveredVerificationJournal(dir *os.Root) error {
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
	if len(names) != 2 {
		return errors.New("verification recovery requires an unexecuted, unstopped journal")
	}
	for _, name := range names {
		if name != "request.json" && name != "launch.json" {
			return errors.New("verification recovery requires an unexecuted, unstopped journal")
		}
	}
	return nil
}

// refusedVerificationEvidence enforces the exact retained budget-refused-before-
// start lifecycle: the original requested and terminal records agree, no started
// lifecycle exists, and the terminal carries the verifier launch's phase with
// status failed, failure class budget_refused, and no PID or exit code. It runs
// both at recovery admission and, through readVerificationLaunch, on every
// recovered read/reserve/execute/stop/receipt — where its returned digests are
// additionally compared against the recovery's bound ones, so evidence changed
// after recovery fails closed instead of being trusted from recovery.json, and
// where the returned refusal completion time bounds the recovery's own At: an
// explicit recovery cannot predate the refusal it repairs. This binds retained
// bytes only. It claims no custody: hashes do not authenticate writers against
// same-UID edits, and no descendant, silence or inactivity inference is used
// or implied.
func refusedVerificationEvidence(ticket VerificationTicket, refusedInvocation string, reservedAt time.Time) (terminalSHA, requestedSHA string, refusedCompletedAt time.Time, err error) {
	root, err := os.OpenRoot(ticket.Root)
	if err != nil {
		return "", "", time.Time{}, err
	}
	defer root.Close()
	base := filepath.Join(".parley-runtime", "invocations", refusedInvocation)
	var requested, terminal telemetry.Record
	requestedSHA, err = readReconciliationJSON(root, filepath.Join(base, "requested.json"), &requested)
	if err != nil || requested.SchemaVersion != telemetry.SchemaVersion || requested.Type != "invocation.requested" ||
		requested.InvocationID != refusedInvocation || requested.RequestedAt.IsZero() || requested.StartedAt != nil ||
		requested.PID != nil || requested.CompletedAt != nil || requested.Outcome != nil ||
		requested.Metadata.RunID != ticket.RunID || requested.Metadata.Idea != ticket.Request.Idea ||
		requested.Metadata.Agent != ticket.Request.Verifier || requested.Metadata.Phase != capturedVerificationPhaseLabel ||
		requested.Metadata.LaunchMode != "headless" || reservedAt.Before(requested.RequestedAt) {
		return "", "", time.Time{}, errors.New("refused verifier launch request is unavailable or changed")
	}
	terminalSHA, err = readReconciliationJSON(root, filepath.Join(base, "terminal.json"), &terminal)
	if err != nil || terminal.SchemaVersion != telemetry.SchemaVersion || terminal.Type != "invocation.terminal" ||
		terminal.InvocationID != refusedInvocation || !sameJSON(terminal.Metadata, requested.Metadata) ||
		!terminal.RequestedAt.Equal(requested.RequestedAt) || terminal.StartedAt != nil || terminal.PID != nil ||
		terminal.CompletedAt == nil || terminal.Outcome == nil || terminal.Outcome.Status != "failed" ||
		terminal.Outcome.FailureClass == nil || *terminal.Outcome.FailureClass != "budget_refused" ||
		terminal.Outcome.ExitCode != nil || terminal.CompletedAt.Before(reservedAt) {
		return "", "", time.Time{}, errors.New("refused verifier launch lacks its exact pre-start budget-refusal terminal")
	}
	if _, err = root.Lstat(filepath.Join(base, "started.json")); !os.IsNotExist(err) {
		return "", "", time.Time{}, errors.New("refused verifier launch must not retain any started lifecycle")
	}
	return terminalSHA, requestedSHA, *terminal.CompletedAt, nil
}
