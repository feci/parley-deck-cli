package runner

import (
	"context"
	"crypto/sha256"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"sort"
	"time"

	"parley-deck-cli/internal/agents"
	"parley-deck-cli/internal/membership"
	"parley-deck-cli/internal/pidlease"
	"parley-deck-cli/internal/quota"
	"parley-deck-cli/internal/telemetry"
)

const ParticipantRetryDelay = 5 * time.Second

type StepValidation struct {
	Valid  bool
	Reason string
	// Integrity distinguishes shared-file corruption from invalid OWN output.
	Integrity error
}

type ParticipantStepOptions struct {
	Root, Idea, Step string
	Agent            agents.Discovery
	// Files are preserved privately before dispatch and after each failed
	// attempt, before a retry can overwrite a log or incomplete own artifact.
	Files    []string
	Validate func() StepValidation
	// Restore only a structurally invalid OWN suffix after its private copy and
	// validation receipt are durable. Shared-file integrity failures never call it.
	AfterFailure func() error
	RecoverOnly  bool // settle/preserve an interrupted suffix without a new child
}

type ParticipantStepResult struct {
	Valid    bool
	Replayed bool
	Record   telemetry.Record
	Evidence *quota.Evidence
}

type participantStepKey struct{}

func ParticipantStepActive(ctx context.Context) bool {
	_, ok := ctx.Value(participantStepKey{}).(*participantStepLaunch)
	return ok
}

type participantStepLaunch struct {
	idea      string
	key       string
	ordinal   int
	retry     string
	observe   func(telemetry.Record)
	integrity error
}

type participantValidationRecord struct {
	InvocationID   string `json:"invocation_id"`
	Step           string `json:"step"`
	TerminalSHA256 string `json:"terminal_sha256"`
	Valid          bool   `json:"valid"`
	Reason         string `json:"reason,omitempty"`
	Integrity      string `json:"integrity,omitempty"`
}

// RunParticipantStep bounds one logical step across runs/restarts. Immutable
// requested/terminal telemetry is the attempt ledger; the small validation
// receipt supplies the phase-specific structural result. There is no second
// membership reducer. Callers settle the whole batch through membership.Settle.
func RunParticipantStep(ctx context.Context, opts ParticipantStepOptions, run func(context.Context, int, string) error) (out ParticipantStepResult, returnedErr error) {
	if opts.Idea == "" || opts.Step == "" || opts.Agent.ID == "" || opts.Validate == nil {
		return out, fmt.Errorf("incomplete participant step identity/validator")
	}
	if origin, ok := ctx.Value(launchOriginKey{}).(string); ok {
		opts.Root = origin
	}
	key := participantKey(opts.Idea, opts.Agent.ID, opts.Step)
	stateDir := filepath.Join(opts.Root, ".parley-runtime", "participant-steps", key)
	lease, err := pidlease.TryAcquire(filepath.Join(stateDir, "step.lease"), opts.Idea+"/"+opts.Agent.ID+"/"+opts.Step)
	if err != nil {
		return out, err
	}
	defer lease.Release()
	if err := preserveStepFiles(filepath.Join(stateDir, "initial"), opts.Files); err != nil {
		return out, err
	}
	records, err := participantRecords(opts.Root, key, opts.Idea, opts.Agent.ID)
	if err != nil {
		return out, err
	}
	var attempts []quota.FailedAttempt
	for i, r := range records {
		out.Record = r
		validation, err := readParticipantValidation(opts.Root, key, r)
		if err != nil {
			return out, err
		}
		if validation == nil {
			// A crash after terminal but before structural validation is safely
			// recoverable only for the last invocation. A later launch without
			// the earlier receipt is contradictory, never a fresh attempt budget.
			if i != len(records)-1 {
				return out, fmt.Errorf("missing prior participant validation receipt")
			}
			if err := preserveStepFiles(filepath.Join(opts.Root, ".parley-runtime", "invocations", r.InvocationID, "participant-files"), opts.Files); err != nil {
				return out, err
			}
			v := opts.Validate()
			validation, err = writeParticipantValidation(opts.Root, key, r, v)
			if err != nil {
				return out, err
			}
		}
		if validation.Integrity != "" {
			return out, fmt.Errorf("preserved participant integrity failure: %s", validation.Integrity)
		}
		a, err := telemetry.ParticipantAttempt(r, validation.Valid, validation.Reason)
		if err != nil {
			// A repaired control-plane refusal can be explicitly dispatched
			// again. It consumed no CHILD attempt; it never triggers a retry
			// inside the invocation that observed the refusal.
			if r.Outcome != nil && r.Outcome.FailureClass != nil && !quota.ParticipantFailureClass(*r.Outcome.FailureClass) {
				continue
			}
			return out, err
		}
		if a == nil {
			if v := opts.Validate(); !v.Valid || v.Integrity != nil {
				return out, fmt.Errorf("previously valid participant output changed; no new attempt authorized")
			}
			out.Valid, out.Replayed = true, true
			return out, nil
		}
		if err := appendParticipantAttempt(&attempts, *a, r.Metadata.AttemptOrdinal); err != nil {
			return out, err
		}
		if i == len(records)-1 && opts.AfterFailure != nil {
			if err := opts.AfterFailure(); err != nil {
				return out, err
			}
		}
	}
	if opts.RecoverOnly && len(attempts) < 2 {
		return out, nil
	}
	for len(attempts) < 2 {
		if err := ctx.Err(); err != nil {
			return out, err
		}
		if err := lease.Check(); err != nil {
			return out, err
		}
		if len(attempts) == 1 {
			delay := time.Until(attempts[0].ObservedAt.Add(ParticipantRetryDelay))
			if delay > 0 {
				if err := participantRetryWait(ctx, delay); err != nil {
					return out, err
				}
			}
		}
		retry := ""
		if len(attempts) > 0 {
			retry = attempts[len(attempts)-1].InvocationID
		}
		var observed telemetry.Record
		binding := &participantStepLaunch{idea: opts.Idea, key: key, ordinal: len(attempts) + 1, retry: retry,
			observe: func(r telemetry.Record) { observed = r }}
		childCtx := context.WithValue(ctx, participantStepKey{}, binding)
		runErr := run(childCtx, binding.ordinal, retry)
		out.Record = observed
		if err := ctx.Err(); err != nil {
			return out, errors.Join(runErr, err)
		}
		if observed.Outcome == nil || observed.CompletedAt == nil {
			return out, errors.Join(runErr, fmt.Errorf("participant invocation has no terminal evidence"))
		}
		v := opts.Validate()
		if binding.integrity != nil {
			v.Integrity = binding.integrity
		}
		if err := preserveStepFiles(filepath.Join(opts.Root, ".parley-runtime", "invocations", observed.InvocationID, "participant-files"), opts.Files); err != nil {
			return out, err
		}
		if _, err := writeParticipantValidation(opts.Root, key, observed, v); err != nil {
			return out, err
		}
		if v.Integrity != nil {
			return out, v.Integrity
		}
		a, err := telemetry.ParticipantAttempt(observed, v.Valid, v.Reason)
		if err != nil {
			return out, errors.Join(runErr, err)
		}
		if a == nil {
			out.Valid = true
			return out, nil
		}
		if err := appendParticipantAttempt(&attempts, *a, observed.Metadata.AttemptOrdinal); err != nil {
			return out, err
		}
		if opts.AfterFailure != nil {
			if err := opts.AfterFailure(); err != nil {
				return out, err
			}
		}
	}
	e, err := telemetry.ParticipantEvidence(opts.Idea, opts.Agent.ID, opts.Step, opts.Agent.Adapter(), attempts)
	if err != nil {
		return out, err
	}
	out.Evidence = &e
	out.Replayed = len(records) > 0
	return out, fmt.Errorf("participant step failed twice: %s", opts.Step)
}

func participantKey(idea, agent, step string) string {
	return fmt.Sprintf("%x", sha256.Sum256([]byte(idea+"\x00"+agent+"\x00"+step)))
}

// ParticipantBaseline returns the privately preserved first artifact, if any.
// The caller must hold the idea lease before using it for shared-file recovery.
func ParticipantBaseline(root, idea, agent, step, name string) ([]byte, error) {
	return os.ReadFile(filepath.Join(root, ".parley-runtime", "participant-steps", participantKey(idea, agent, step), "initial", "00-"+filepath.Base(name)))
}

func ParticipantStepPending(root, idea, agent, step string) (bool, error) {
	key := participantKey(idea, agent, step)
	records, err := participantRecords(root, key, idea, agent)
	if err != nil {
		return false, err
	}
	for i, r := range records {
		v, err := readParticipantValidation(root, key, r)
		if err != nil {
			return false, err
		}
		if v != nil && v.Integrity != "" {
			return false, fmt.Errorf("preserved participant integrity failure: %s", v.Integrity)
		}
		if i == len(records)-1 {
			return v == nil || !v.Valid, nil
		}
	}
	return false, nil
}

func appendParticipantAttempt(attempts *[]quota.FailedAttempt, a quota.FailedAttempt, ordinal int) error {
	prior := ""
	if len(*attempts) > 0 {
		prior = (*attempts)[len(*attempts)-1].InvocationID
	}
	if len(*attempts) >= 2 || ordinal != len(*attempts)+1 || a.RetryOf != prior {
		return fmt.Errorf("contradictory participant attempt count/linkage")
	}
	*attempts = append(*attempts, a)
	return nil
}

var participantRetryWait = func(ctx context.Context, delay time.Duration) error {
	timer := time.NewTimer(delay)
	defer timer.Stop()
	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-timer.C:
		return nil
	}
}

func participantRecords(root, key, idea, agent string) ([]telemetry.Record, error) {
	paths, err := filepath.Glob(filepath.Join(root, ".parley-runtime", "invocations", "*", "requested.json"))
	if err != nil {
		return nil, err
	}
	var records []telemetry.Record
	for _, path := range paths {
		raw, err := os.ReadFile(path)
		if err != nil {
			return nil, err
		}
		var requested telemetry.Record
		if json.Unmarshal(raw, &requested) != nil {
			return nil, fmt.Errorf("unreadable invocation request")
		}
		if requested.Metadata.ParticipantStep != key {
			continue
		}
		if requested.Metadata.Idea != idea || requested.Metadata.Agent != agent || requested.InvocationID != filepath.Base(filepath.Dir(path)) {
			return nil, fmt.Errorf("participant step identity conflict")
		}
		raw, err = os.ReadFile(filepath.Join(filepath.Dir(path), "terminal.json"))
		if os.IsNotExist(err) {
			// This validates/reuses the existing host/boot/PID crash proof;
			// it does not invent a provider outcome or terminal telemetry.
			r, err := membership.ParticipantCrash(root, requested)
			if err != nil {
				return nil, err
			}
			records = append(records, r)
			continue
		}
		if err != nil {
			return nil, err
		}
		var terminal telemetry.Record
		if json.Unmarshal(raw, &terminal) != nil || terminal.Type != "invocation.terminal" || terminal.InvocationID != requested.InvocationID || !reflect.DeepEqual(terminal.Metadata, requested.Metadata) || terminal.Outcome == nil || terminal.CompletedAt == nil || terminal.CompletedAt.Before(requested.RequestedAt) {
			return nil, fmt.Errorf("contradictory participant terminal evidence")
		}
		records = append(records, terminal)
	}
	sort.Slice(records, func(i, j int) bool { return records[i].RequestedAt.Before(records[j].RequestedAt) })
	return records, nil
}

func participantTerminalHash(r telemetry.Record) string {
	raw, _ := json.Marshal(r)
	return fmt.Sprintf("%x", sha256.Sum256(raw))
}

func readParticipantValidation(root, key string, r telemetry.Record) (*participantValidationRecord, error) {
	path := filepath.Join(root, ".parley-runtime", "invocations", r.InvocationID, "participant-result.json")
	raw, err := os.ReadFile(path)
	if os.IsNotExist(err) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	var v participantValidationRecord
	if quota.DecodeRevisionRequest(raw, &v) != nil || v.InvocationID != r.InvocationID || v.Step != key || v.TerminalSHA256 != participantTerminalHash(r) {
		return nil, fmt.Errorf("contradictory participant validation receipt")
	}
	return &v, nil
}

func writeParticipantValidation(root, key string, r telemetry.Record, v StepValidation) (*participantValidationRecord, error) {
	record := participantValidationRecord{InvocationID: r.InvocationID, Step: key, TerminalSHA256: participantTerminalHash(r), Valid: v.Valid, Reason: v.Reason}
	if v.Integrity != nil {
		record.Integrity = v.Integrity.Error()
	}
	raw, err := json.Marshal(record)
	if err != nil {
		return nil, err
	}
	path := filepath.Join(root, ".parley-runtime", "invocations", r.InvocationID, "participant-result.json")
	if err := quota.DurableWrite(path, raw, true); err != nil {
		return nil, err
	}
	return &record, nil
}

func preserveStepFiles(dir string, paths []string) error {
	for i, path := range paths {
		st, err := os.Lstat(path)
		if os.IsNotExist(err) {
			continue
		}
		if err != nil {
			return err
		}
		if !st.Mode().IsRegular() {
			return fmt.Errorf("participant evidence path is not a regular file: %s", path)
		}
		name := filepath.Join(dir, fmt.Sprintf("%02d-%s", i, filepath.Base(path)))
		if _, err := os.Lstat(name); err == nil {
			continue
		} else if !os.IsNotExist(err) {
			return err
		}
		raw, err := os.ReadFile(path)
		if err != nil {
			return err
		}
		if err := quota.DurableWrite(name, raw, true); err != nil {
			return err
		}
	}
	return nil
}
