package runner

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"path/filepath"
	"time"

	"parley-deck-cli/internal/agents"
	"parley-deck-cli/internal/telemetry"
	"parley-deck-cli/internal/trajectory"
)

const CapturedVerificationPhase = "trajectory-verification"

type capturedVerificationKey struct{}

// WithCapturedVerification binds a prepared independent verification request to
// the actual launch boundary. It detaches caller-owned slices and pointers; a
// label in LaunchInfo alone cannot supply this authority or allow a retry.
func WithCapturedVerification(ctx context.Context, ticket trajectory.VerificationTicket) (context.Context, error) {
	frozen, err := freezeVerificationTicket(ticket)
	if err != nil {
		return nil, err
	}
	return context.WithValue(ctx, capturedVerificationKey{}, frozen), nil
}

func freezeVerificationTicket(ticket trajectory.VerificationTicket) (trajectory.VerificationTicket, error) {
	if _, err := ticket.SHA256(); err != nil {
		return trajectory.VerificationTicket{}, err
	}
	data, err := json.Marshal(ticket)
	if err != nil {
		return trajectory.VerificationTicket{}, err
	}
	var frozen trajectory.VerificationTicket
	if err = json.Unmarshal(data, &frozen); err != nil {
		return trajectory.VerificationTicket{}, err
	}
	return frozen, nil
}

type capturedVerificationRecovery struct {
	ticket     trajectory.VerificationTicket
	invocation string
}

type capturedVerificationRecoveryKey struct{}

// WithCapturedVerificationRecovery is the narrow trusted bridge that lets ONE
// real headless verifier relaunch execute under the invocation ID an immutable
// captured-verification recovery already pinned. The caller supplies the ticket
// plus the intended invocation exactly as recorded in the journal's
// recovery.json — the authority that applied the recovery chose that ID; this
// API never chooses or accepts a caller-invented one as authority.
//
// Before any launch exists, it validates the recovery authority through the
// genuinely read-only trajectory recovery reader: the full guarded authority is
// rechecked, the retained recovery is revalidated against the live launch
// reservation and the exact refused-launch evidence (bound digests and
// chronology), and the returned bound identity must name the intended
// invocation. The read performs no reservation and no publication of any kind —
// unlike the former existence probes plus reservation call, a missing or
// disappearing launch/recovery authority is a refusal here, never a recreated
// launch reservation. The validation is repeated at the actual launch's
// verification reservation, so a recovery removed or changed between this call
// and the real launch fails closed before any process or budget action.
//
// The bound context also binds the ticket exactly as WithCapturedVerification;
// the launch still passes every ordinary scoped check (origin, run, idea,
// verifier, phase, headless), still pays its own budget admission and
// settlement, and never refunds, resets or bypasses anything. Ordinary
// launches without this context always generate a fresh invocation ID:
// LaunchInfo fields, retry_of and phase labels never select one. A bound ID
// can launch exactly once — telemetry exclusively creates the invocation
// directory — so concurrent or replayed launches with one recovery ID cannot
// execute twice.
func WithCapturedVerificationRecovery(ctx context.Context, ticket trajectory.VerificationTicket, intendedInvocation string) (context.Context, error) {
	frozen, err := freezeVerificationTicket(ticket)
	if err != nil {
		return nil, err
	}
	if !telemetry.ValidBoundInvocationID(intendedInvocation) {
		return nil, errors.New("bound verifier recovery requires a safe single-element invocation identifier")
	}
	origin, absErr := filepath.Abs(ticket.Root)
	if absErr == nil {
		origin, absErr = filepath.EvalSymlinks(origin)
	}
	if absErr != nil || origin != ticket.Root {
		return nil, errors.New("bound verifier recovery origin is unavailable or no longer canonical")
	}
	identity, err := trajectory.ReadCapturedVerificationRecovery(ctx, ticket)
	if err != nil {
		return nil, fmt.Errorf("bound verifier recovery authority refused the intended invocation: %w", err)
	}
	switch intendedInvocation {
	case identity.InvocationID:
	case identity.RefusedInvocationID:
		return nil, errors.New("bound verifier recovery authority refused the intended invocation: verification invocation was superseded by an explicit recovery; it retains no control over the recovered launch")
	default:
		return nil, errors.New("bound verifier recovery authority refused the intended invocation: verification invocation differs from its durable reservation")
	}
	ctx = context.WithValue(ctx, capturedVerificationKey{}, frozen)
	return context.WithValue(ctx, capturedVerificationRecoveryKey{}, capturedVerificationRecovery{ticket: frozen, invocation: intendedInvocation}), nil
}

// boundRecoveryInvocation resolves the pinned invocation a trusted recovery
// context authorizes for THIS launch, or "" for every ordinary launch. It
// refuses — before any telemetry directory is created — a bound launch that
// does not already carry the exact verifier shape the reservation boundary
// will require, so a mis-bound or handoff launch cannot consume the one-shot
// pinned invocation. The authoritative recheck stays in
// reserveCapturedVerification.
func boundRecoveryInvocation(ctx context.Context, info LaunchInfo, agent agents.Discovery, handoff bool) (string, error) {
	bound, ok := ctx.Value(capturedVerificationRecoveryKey{}).(capturedVerificationRecovery)
	if !ok {
		return "", nil
	}
	ticket, boundTicket := ctx.Value(capturedVerificationKey{}).(trajectory.VerificationTicket)
	want, wantErr := ticket.SHA256()
	pinned, pinnedErr := bound.ticket.SHA256()
	if !boundTicket || wantErr != nil || pinnedErr != nil || want != pinned {
		return "", errors.New("bound verifier recovery launch lost its prepared ticket binding")
	}
	if handoff || info.Phase != CapturedVerificationPhase || info.Idea != bound.ticket.Request.Idea ||
		info.RunID != bound.ticket.RunID || agent.ID != bound.ticket.Request.Verifier ||
		agents.LaunchModeOrDefault(agent.LaunchMode) != agents.LaunchHeadless {
		return "", errors.New("bound verifier recovery launch changed its prepared origin, run, participant or mode")
	}
	return bound.invocation, nil
}

func (l *launchEvidence) reserveCapturedVerification(ctx context.Context, root string, handoff bool) error {
	ticket, required := ctx.Value(capturedVerificationKey{}).(trajectory.VerificationTicket)
	if !required {
		if l.info.Phase == CapturedVerificationPhase {
			return errors.New("trajectory verifier launch requires a prepared ticket")
		}
		return nil
	}
	origin, err := filepath.Abs(root)
	if err == nil {
		origin, err = filepath.EvalSymlinks(origin)
	}
	record := l.invocation.Snapshot()
	if err != nil || origin != ticket.Root || handoff ||
		l.info.Phase != CapturedVerificationPhase || l.info.Idea != ticket.Request.Idea ||
		l.info.RunID != ticket.RunID || record.Metadata.Agent != ticket.Request.Verifier ||
		record.Metadata.LaunchMode != agents.LaunchHeadless ||
		record.Metadata.RunID != ticket.RunID || record.Metadata.Idea != ticket.Request.Idea ||
		record.Metadata.Phase != CapturedVerificationPhase {
		return errors.New("trajectory verifier launch changed its prepared origin, run, participant or mode")
	}
	return trajectory.ReserveCapturedVerificationLaunch(ctx, ticket, record.InvocationID)
}

func (l *launchEvidence) stopCapturedVerification(ctx context.Context) error {
	ticket, required := ctx.Value(capturedVerificationKey{}).(trajectory.VerificationTicket)
	if !required {
		return nil
	}
	cleanup, cancel := context.WithTimeout(context.WithoutCancel(ctx), 20*time.Second)
	defer cancel()
	return trajectory.StopCapturedVerification(cleanup, ticket, l.invocation.ID)
}
