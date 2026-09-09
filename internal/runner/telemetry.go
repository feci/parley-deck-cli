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

func WithLaunchInfo(ctx context.Context, info LaunchInfo) context.Context {
	return context.WithValue(ctx, launchInfoKey{}, info)
}

type launchEvidence struct {
	invocation *telemetry.Invocation
	collector  *telemetry.Collector
	info       LaunchInfo
	once       sync.Once
	finishErr  error
}

// Integrity errors deliberately do not unwrap an ordinary process exit: a
// validated artifact must not override a lost audit record or provider error.
type launchIntegrityError struct{ reason string }

func (e *launchIntegrityError) Error() string { return e.reason }

func beginLaunch(ctx context.Context, root, runID string, agent agents.Discovery) (*launchEvidence, error) {
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
	metadata := telemetry.Metadata{
		RunID: info.RunID, SegmentID: info.SegmentID, Idea: info.Idea, Phase: info.Phase,
		Agent: agent.ID, Adapter: agent.Adapter(), LaunchMode: agents.LaunchModeOrDefault(agent.LaunchMode),
		AttemptOrdinal: info.AttemptOrdinal, RetryOf: telemetry.String(info.RetryOf), Context: info.Context,
		RequestedModel: requestedSetting(agent.Model), RequestedEffort: requestedSetting(agent.Reasoning),
		RequestedSpeed: requestedSetting(agent.Speed),
	}
	invocation, err := telemetry.Begin(filepath.Join(root, ".parley-runtime", "invocations"), metadata)
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
