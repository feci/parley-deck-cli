package runner

import (
	"context"
	"encoding/json"
	"errors"
	"path/filepath"
	"time"

	"parley-deck-cli/internal/agents"
	"parley-deck-cli/internal/trajectory"
)

const CapturedVerificationPhase = "trajectory-verification"

type capturedVerificationKey struct{}

// WithCapturedVerification binds a prepared independent verification request to
// the actual launch boundary. It detaches caller-owned slices and pointers; a
// label in LaunchInfo alone cannot supply this authority or allow a retry.
func WithCapturedVerification(ctx context.Context, ticket trajectory.VerificationTicket) (context.Context, error) {
	if _, err := ticket.SHA256(); err != nil {
		return nil, err
	}
	data, err := json.Marshal(ticket)
	if err != nil {
		return nil, err
	}
	var frozen trajectory.VerificationTicket
	if err = json.Unmarshal(data, &frozen); err != nil {
		return nil, err
	}
	return context.WithValue(ctx, capturedVerificationKey{}, frozen), nil
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
