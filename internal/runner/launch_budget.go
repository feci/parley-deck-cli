package runner

import (
	"context"
	"errors"
	"fmt"
	"math"
	"time"

	"parley-deck-cli/internal/budget"
)

// LaunchBudget is trusted orchestration policy. It is not read from participant
// artifacts. The caller must resolve the same durable Store and frozen Limits
// for every entrypoint/resume; this context API does not perform that resolution.
type LaunchBudget struct {
	Store         budget.Store
	Limits        budget.Limits
	ReserveMicros *int64
}

type launchBudgetKey struct{}
type launchIntent bool

const launchHandoff launchIntent = true

// WithLaunchBudget keeps policy independent of replaceable launch metadata.
// Detach mutable limits so a caller cannot accidentally change in-flight policy.
func WithLaunchBudget(ctx context.Context, policy LaunchBudget) context.Context {
	copyMap := func(in map[budget.Kind]int) map[budget.Kind]int {
		out := make(map[budget.Kind]int, len(in))
		for k, v := range in {
			out[k] = v
		}
		return out
	}
	policy.Limits.Actions = copyMap(policy.Limits.Actions)
	denied := make(map[budget.Kind]bool, len(policy.Limits.Denied))
	for k, v := range policy.Limits.Denied {
		denied[k] = v
	}
	policy.Limits.Denied = denied
	if policy.ReserveMicros != nil {
		value := *policy.ReserveMicros
		policy.ReserveMicros = &value
	}
	return context.WithValue(ctx, launchBudgetKey{}, policy)
}

type launchBudgetError struct{ cause error }

func (e *launchBudgetError) Error() string { return fmt.Sprintf("launch budget refused: %v", e.cause) }
func (e *launchBudgetError) Unwrap() error { return e.cause }

func (l *launchEvidence) reserveBudget(ctx context.Context, handoff bool) error {
	policy, enabled := ctx.Value(launchBudgetKey{}).(LaunchBudget)
	if !enabled || handoff {
		return nil
	}
	_, err := policy.Store.Reserve(ctx, budget.Request{
		ID: l.invocation.ID, Kind: budget.Launch, ReserveMicros: policy.ReserveMicros,
	}, policy.Limits)
	if err != nil {
		return &launchBudgetError{cause: err}
	}
	l.budget = &policy
	return nil
}

// A failed start/telemetry write does not refund a reservation. Settlement must
// run even when the process deadline/cancellation or terminal publication failed.
// A CLI estimate remains an estimate in telemetry; rounding up here is solely
// conservative budget accounting, not an invoice or an invented provider price.
func (l *launchEvidence) settleBudget(costUSD *float64) error {
	if l.budget == nil {
		return nil
	}
	var micros *int64
	if costUSD != nil {
		cost := math.Ceil(*costUSD * 1_000_000)
		if math.IsNaN(cost) || math.IsInf(cost, 0) || cost < 0 || cost >= math.Exp2(63) {
			return errors.New("invalid terminal cost; launch reservation retained")
		}
		value := int64(cost)
		micros = &value
	}
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	_, err := l.budget.Store.Settle(ctx, l.invocation.ID, micros)
	return err
}
