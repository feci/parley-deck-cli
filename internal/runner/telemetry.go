package runner

import (
	"context"
	"crypto/sha256"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sync"
	"time"

	"parley-deck-cli/internal/agents"
	"parley-deck-cli/internal/store"
	"parley-deck-cli/internal/telemetry"
	"parley-deck-cli/internal/trajectory"
)

// LaunchInfo describes orchestration, never prompt content or command arguments.
// An observer receives a detached record; it cannot alter persisted evidence.
type LaunchInfo struct {
	RunID, SegmentID, Idea, Phase string
	AttemptOrdinal                int
	RetryOf                       string
	Context                       telemetry.Context
	Store                         store.Store
	Observe                       func(telemetry.Record)
	ArtifactPath                  string
}

type launchInfoKey struct{}
type launchOriginKey struct{}

// The protocol runner captures its live origin before a disposable review
// checkout replaces the execution cwd. This is runtime lineage, never a
// participant-provided telemetry field. Budgets and evidence survive cleanup.
func withLaunchOrigin(ctx context.Context, root string) context.Context {
	return context.WithValue(ctx, launchOriginKey{}, root)
}

func WithLaunchInfo(ctx context.Context, info LaunchInfo) context.Context {
	return context.WithValue(ctx, launchInfoKey{}, info)
}

type launchEvidence struct {
	invocation *telemetry.Invocation
	collector  *telemetry.Collector
	info       LaunchInfo
	once       sync.Once
	finishErr  error
	budget     *LaunchBudget
	trajectory *trajectory.Run
	// A terminal's file descriptors must reach the child unchanged. Such a
	// process has lifecycle evidence but no captured output stream evidence.
	directTerminal bool
}

// Integrity errors deliberately do not unwrap an ordinary process exit: a
// validated artifact must not override a lost audit record or provider error.
type launchIntegrityError struct{ reason string }

func (e *launchIntegrityError) Error() string { return e.reason }

func launchRequestInfo(ctx context.Context, root, runID string) (string, LaunchInfo) {
	if origin, ok := ctx.Value(launchOriginKey{}).(string); ok {
		root = origin
	}
	info, _ := ctx.Value(launchInfoKey{}).(LaunchInfo)
	if info.RunID == "" {
		info.RunID = runID
	}
	if info.Phase == "" {
		info.Phase = "unspecified"
	}
	if info.Context.Mode == "" {
		info.Context.Mode = "unattested"
	}
	return root, info
}

func beginLaunch(ctx context.Context, root, runID string, agent agents.Discovery, intent ...launchIntent) (*launchEvidence, error) {
	root, info := launchRequestInfo(ctx, root, runID)
	ctx = withLaunchActionInput(ctx, info, agent)
	handoff := len(intent) == 1 && intent[0] == launchHandoff
	if !handoff {
		var finishCycle func()
		ctx, finishCycle = prepareLaunchCycle(ctx, root, info.Idea, info.Phase, info.RunID)
		defer finishCycle()
	}
	boundInvocation, err := boundRecoveryInvocation(ctx, info, agent, handoff)
	if err != nil {
		return nil, err
	}
	var l *launchEvidence
	if boundInvocation != "" {
		l, err = recordBoundLaunchRequest(root, agent, info, boundInvocation)
	} else {
		l, err = recordLaunchRequest(root, agent, info)
	}
	if err != nil {
		return nil, err
	}
	if err := l.reserveCapturedVerification(ctx, root, handoff); err != nil {
		return nil, errors.Join(&launchIntegrityError{reason: err.Error()}, l.finish(err, ctx.Err(), nil))
	}
	if err := l.reserveBudget(ctx, root, handoff); err != nil {
		return nil, errors.Join(err, l.finish(err, ctx.Err(), nil))
	}
	return l, nil
}

// Records a request only. It prepares no cycle and reserves no execution.
func recordLaunchRequest(root string, agent agents.Discovery, info LaunchInfo) (*launchEvidence, error) {
	return recordLaunch(root, agent, info, "")
}

// recordBoundLaunchRequest is the one trusted recovery variant: the invocation
// ID was pinned by an immutable captured-verification recovery before this
// launch existed, so telemetry must allocate exactly that ID, exclusively, or
// refuse. Every ordinary launch keeps generating a fresh ID.
func recordBoundLaunchRequest(root string, agent agents.Discovery, info LaunchInfo, boundInvocation string) (*launchEvidence, error) {
	return recordLaunch(root, agent, info, boundInvocation)
}

func recordLaunch(root string, agent agents.Discovery, info LaunchInfo, boundInvocation string) (*launchEvidence, error) {
	metadata := telemetry.Metadata{
		RunID: info.RunID, SegmentID: info.SegmentID, Idea: info.Idea, Phase: info.Phase,
		Agent: agent.ID, Adapter: agent.Adapter(), LaunchMode: agents.LaunchModeOrDefault(agent.LaunchMode),
		AttemptOrdinal: info.AttemptOrdinal, RetryOf: telemetry.String(info.RetryOf), Context: info.Context,
		RequestedModel: requestedSetting(agent.Model), RequestedEffort: requestedSetting(agent.Reasoning),
		RequestedSpeed: requestedSetting(agent.Speed),
	}
	directory := filepath.Join(root, ".parley-runtime", "invocations")
	var invocation *telemetry.Invocation
	var err error
	if boundInvocation != "" {
		invocation, err = telemetry.BeginBound(directory, boundInvocation, metadata)
	} else {
		invocation, err = telemetry.Begin(directory, metadata)
	}
	if err != nil {
		return nil, &launchIntegrityError{reason: "cannot persist required launch request"}
	}
	args, _ := agent.Spec.ResolveLaunchArgs()
	structured := metadata.LaunchMode != agents.LaunchACP && telemetry.StructuredArgs(args)
	l := &launchEvidence{invocation: invocation, collector: telemetry.NewCollector(agent.Adapter(), structured), info: info}
	l.notify()
	return l, nil
}

func requestedSetting(value string) *string {
	if value == agents.CLIDefault {
		return nil
	}
	return telemetry.String(value)
}

func (l *launchEvidence) notify() {
	if l.info.Observe != nil {
		l.info.Observe(l.invocation.Snapshot())
	}
}

func (l *launchEvidence) started(pid int) error {
	if err := l.invocation.Started(pid); err != nil {
		return &launchIntegrityError{reason: "cannot persist required process start"}
	}
	l.notify()
	return nil
}

func (l *launchEvidence) finish(runErr, ctxErr error, exitCode *int) error {
	l.once.Do(func() {
		usage, observation, providerFailure := l.collector.Result()
		defer func() {
			if err := l.settleBudget(usage.CostUSD); err != nil {
				l.finishErr = errors.Join(l.finishErr, &launchIntegrityError{reason: "cannot settle required launch budget; reservation retained"})
			}
		}()
		if l.directTerminal {
			observation.StreamCoverage = "not-observed-terminal"
		}
		if l.directTerminal || l.invocation.Snapshot().StartedAt == nil {
			observation.StdoutBytes, observation.StderrBytes = nil, nil
		}
		status, failure := "process-exited", ""
		if runErr == nil && l.invocation.Snapshot().StartedAt == nil {
			status = "unobserved-handoff"
		}
		switch {
		case errors.Is(runErr, errNoFirstOutput):
			status, failure = "failed", "no_first_output"
		case errors.Is(runErr, errStalled):
			status, failure = "failed", "stalled"
		case errors.Is(ctxErr, context.DeadlineExceeded):
			status, failure = "failed", "timeout"
		case errors.Is(ctxErr, context.Canceled):
			status, failure = "failed", "cancelled"
		case providerFailure != "":
			status, failure = "failed", providerFailure
		case runErr != nil:
			status, failure = "failed", "process_failure"
			if l.invocation.Snapshot().StartedAt == nil {
				failure = "start_failure"
			}
		}
		var integrity *launchIntegrityError
		if errors.As(runErr, &integrity) {
			status, failure = "failed", "telemetry_failure"
		}
		var contextFailure *protocolContextError
		if errors.As(runErr, &contextFailure) {
			status, failure = "failed", "protocol_context_refused"
		}
		var budgetFailure *launchBudgetError
		if errors.As(runErr, &budgetFailure) {
			status, failure = "failed", "budget_refused"
		}
		if l.trajectory != nil {
			captureCtx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
			err := l.trajectory.Finish(captureCtx, status, exitCode)
			cancel()
			if err != nil {
				l.finishErr = &launchIntegrityError{reason: "cannot retain trajectory post-state; charged attempt remains unresolved"}
				status, failure = "failed", "trajectory_failure"
			}
		}
		outcome := telemetry.Outcome{Status: status, ExitCode: exitCode,
			FailureClass: telemetry.String(failure), Usage: usage, Observation: observation,
			ArtifactSHA256: observedArtifactHash(l.info.ArtifactPath)}
		if err := l.invocation.Finish(outcome); err != nil {
			l.finishErr = &launchIntegrityError{reason: "cannot persist required terminal telemetry"}
			return
		}
		l.notify()
		if l.info.Store.Enabled() {
			// Project only the already-sanitized terminal record, never raw CLI JSON.
			record := l.invocation.Snapshot()
			encoded, _ := json.Marshal(record.Outcome.Usage)
			data := map[string]any{}
			_ = json.Unmarshal(encoded, &data)
			data["invocation_id"], data["agent"] = record.InvocationID, record.Metadata.Agent
			data["segment_id"], data["phase"] = record.Metadata.SegmentID, record.Metadata.Phase
			data["status"], data["duration_ms"] = record.Outcome.Status, record.DurationMS
			if err := l.info.Store.Append(store.Event{Time: time.Now().UTC(), Type: "agent.usage", Data: data}); err != nil {
				l.finishErr = &launchIntegrityError{reason: "cannot publish required normalized usage event"}
				return
			}
		}
		if providerFailure != "" {
			l.finishErr = &launchIntegrityError{reason: fmt.Sprintf("structured provider failure: %s", providerFailure)}
		} else if ctxErr != nil {
			l.finishErr = ctxErr
		}
	})
	return l.finishErr
}

// This identifies observed bytes, not semantic validity or independent review.
func observedArtifactHash(path string) *string {
	if path == "" {
		return nil
	}
	info, err := os.Lstat(path)
	if err != nil || !info.Mode().IsRegular() {
		return nil
	}
	f, err := os.Open(path)
	if err != nil {
		return nil
	}
	defer f.Close()
	hash := sha256.New()
	if _, err := io.Copy(hash, f); err != nil {
		return nil
	}
	return telemetry.String(fmt.Sprintf("%x", hash.Sum(nil)))
}
