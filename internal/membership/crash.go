package membership

import (
	"bytes"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"parley-deck-cli/internal/pidlease"
	"parley-deck-cli/internal/quota"
	"parley-deck-cli/internal/telemetry"
)

type crashSettlement struct {
	Version       int                      `json:"version"`
	InvocationID  string                   `json:"invocation_id"`
	Idea          string                   `json:"idea"`
	StartedSHA256 string                   `json:"started_sha256"`
	PID           int                      `json:"pid"`
	Writer        telemetry.WriterIdentity `json:"writer"`
	SettledAt     time.Time                `json:"settled_at"`
	Proof         string                   `json:"proof"`
}

const crashProof = "same-host-and-boot; supervisor and writer PIDs absent; supervised process group absent"

var crashDeadLocal = pidlease.ProvenDeadLocal
var crashGroupStopped = pidlease.GroupStopped
var crashWrite = quota.DurableWrite

// Settlements are distinct from normal terminal outcomes: they never fabricate
// success, usage, provider evidence, owner authorization or a budget settlement.
func settleCrashedWriter(path string, raw []byte, started telemetry.Record, idea string) error {
	id := filepath.Base(filepath.Dir(path))
	blocked := func(reason string) error {
		return fmt.Errorf("unsettled writer for idea %s: %s: %s; stop/verify the original host's supervisor and writer, then run parley quota recover; legacy or unknown identity requires restoration of authentic terminal evidence, never a fabricated settlement", idea, id, reason)
	}
	if started.Type != "invocation.started" || started.InvocationID != id || started.Metadata.Idea != idea || started.StartedAt == nil || started.PID == nil || started.Writer == nil {
		return blocked("missing recorded writer identity")
	}
	want := crashSettlement{Version: 1, InvocationID: id, Idea: idea, StartedSHA256: fmt.Sprintf("%x", sha256.Sum256(raw)), PID: *started.PID, Writer: *started.Writer, Proof: crashProof}
	receipt := filepath.Join(filepath.Dir(path), "crash-settlement.json")
	validate := func(data []byte) error {
		var got crashSettlement
		if json.Unmarshal(data, &got) != nil {
			return blocked("invalid crash settlement")
		}
		canonical, _ := json.Marshal(got)
		// Canonical bytes and the start-record digest bind the proof to one invocation.
		if !bytes.Equal(data, canonical) || got.SettledAt.Before(*started.StartedAt) || got.SettledAt.After(time.Now().Add(time.Second)) || got.SettledAt.IsZero() {
			return blocked("invalid crash settlement")
		}
		expected := want
		expected.SettledAt = got.SettledAt
		if got != expected {
			return blocked("contradictory crash settlement")
		}
		return quota.SyncPath(receipt)
	}
	if data, err := os.ReadFile(receipt); err == nil {
		return validate(data)
	} else if !os.IsNotExist(err) {
		return err
	}
	w := started.Writer
	if !crashDeadLocal(w.Host, w.Boot, w.SupervisorPID) || !crashDeadLocal(w.Host, w.Boot, *started.PID) || w.ProcessGroup != *started.PID || !crashGroupStopped(w.ProcessGroup) {
		return blocked("only proven-dead local supervisor, writer and process group may be settled")
	}
	want.SettledAt = time.Now().UTC()
	if want.SettledAt.Before(*started.StartedAt) {
		return blocked("writer start is in the future")
	}
	data, err := json.Marshal(want)
	if err != nil {
		return err
	}
	if err = crashWrite(receipt, data, true); err != nil {
		// A competing recovery may have published the same proof with its own time.
		if existing, readErr := os.ReadFile(receipt); readErr == nil {
			return validate(existing)
		}
		return err
	}
	return nil
}
