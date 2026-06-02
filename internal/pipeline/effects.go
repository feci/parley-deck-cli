package pipeline

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// EffectStatus is the §12.6 per-effect lifecycle. The semantics are append-only:
// attempts accumulate and status only moves forward.
type EffectStatus string

const (
	EffectPlanned    EffectStatus = "planned"
	EffectDryRunOK   EffectStatus = "dry_run_ok"
	EffectExecuting  EffectStatus = "executing"
	EffectSucceeded  EffectStatus = "succeeded"
	EffectFailed     EffectStatus = "failed"
	EffectReconciled EffectStatus = "reconciled"
	EffectAbandoned  EffectStatus = "abandoned"
)

// Attempt is one recorded try at performing an effect.
type Attempt struct {
	At          time.Time `json:"at"`
	Status      EffectStatus `json:"status"`
	ExternalRef string    `json:"external_ref,omitempty"`
	Note        string    `json:"note,omitempty"`
}

// Effect is one side-effect ledger record at
// pipelines/<slug>/effects/<digest>.json (§12.6/§12.7). The driver is its only
// writer; agents never mutate it.
type Effect struct {
	IdempotencyKey string       `json:"idempotency_key"`
	PipelineSlug   string       `json:"pipeline_slug"`
	BlockID        string       `json:"block_id"`
	Provider       string       `json:"provider"`
	Capability     string       `json:"capability"`
	Target         string       `json:"target"`
	RequestHash    string       `json:"request_hash"`
	Risk           Risk         `json:"risk"`
	Status         EffectStatus `json:"status"`
	GateID         string       `json:"gate_id,omitempty"`
	DryRunResult   string       `json:"dry_run_result,omitempty"`
	ExternalRef    string       `json:"external_ref,omitempty"`
	Attempts       []Attempt    `json:"attempts"`
	CreatedAt      time.Time    `json:"created_at"`
	UpdatedAt      time.Time    `json:"updated_at"`
}

// IdempotencyKey is the canonical §12.7 key string for an action. Re-planning
// the same action body yields the same key (dedupe); any change in
// target/capability/body yields a new key.
func IdempotencyKey(slug, blockID, provider, capability, target, requestHash string) string {
	return strings.Join([]string{slug, blockID, provider, capability, target, requestHash}, "|")
}

// KeyDigest is the stable filename prefix derived from the full key.
func KeyDigest(key string) string {
	sum := sha256.Sum256([]byte(key))
	return hex.EncodeToString(sum[:])[:24]
}

// HashRequest is the §12.7 request_hash over a normalized action body.
func HashRequest(body []byte) string {
	sum := sha256.Sum256(body)
	return hex.EncodeToString(sum[:])[:16]
}

// NewEffect builds a planned effect with its idempotency key computed.
func NewEffect(slug, blockID, provider, capability, target, requestHash string, risk Risk, now time.Time) Effect {
	key := IdempotencyKey(slug, blockID, provider, capability, target, requestHash)
	return Effect{
		IdempotencyKey: key,
		PipelineSlug:   slug,
		BlockID:        blockID,
		Provider:       provider,
		Capability:     capability,
		Target:         target,
		RequestHash:    requestHash,
		Risk:           risk,
		Status:         EffectPlanned,
		Attempts:       []Attempt{{At: now, Status: EffectPlanned}},
		CreatedAt:      now,
		UpdatedAt:      now,
	}
}

// EffectPath is the on-disk location for an effect, keyed by digest.
func EffectPath(deckDir, slug, key string) string {
	return filepath.Join(PipelineDir(deckDir, slug), "effects", KeyDigest(key)+".json")
}

// Advance records a forward status transition with an attempt entry.
func (e *Effect) Advance(status EffectStatus, externalRef, note string, now time.Time) {
	e.Status = status
	if externalRef != "" {
		e.ExternalRef = externalRef
	}
	e.Attempts = append(e.Attempts, Attempt{At: now, Status: status, ExternalRef: externalRef, Note: note})
	e.UpdatedAt = now
}

// SaveEffect writes an effect atomically. The same idempotency key always maps
// to the same file, so re-planning never creates a duplicate record.
func SaveEffect(deckDir string, e Effect) error {
	path := EffectPath(deckDir, e.PipelineSlug, e.IdempotencyKey)
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return fmt.Errorf("create effects dir: %w", err)
	}
	data, err := json.MarshalIndent(e, "", "  ")
	if err != nil {
		return fmt.Errorf("encode effect: %w", err)
	}
	data = append(data, '\n')
	tmp := path + ".tmp"
	if err := os.WriteFile(tmp, data, 0o644); err != nil {
		return fmt.Errorf("write effect: %w", err)
	}
	if err := os.Rename(tmp, path); err != nil {
		return fmt.Errorf("commit effect: %w", err)
	}
	return nil
}

// LoadEffect reads an effect by its full idempotency key; the bool reports
// whether it exists.
func LoadEffect(deckDir, slug, key string) (Effect, bool, error) {
	path := EffectPath(deckDir, slug, key)
	data, err := os.ReadFile(path)
	if os.IsNotExist(err) {
		return Effect{}, false, nil
	}
	if err != nil {
		return Effect{}, false, fmt.Errorf("read effect: %w", err)
	}
	var e Effect
	if err := json.Unmarshal(data, &e); err != nil {
		return Effect{}, false, fmt.Errorf("parse effect: %w", err)
	}
	return e, true, nil
}

// BlockHasSucceededEffect reports whether the block has at least one effect
// recorded as succeeded. An action block is "complete" only when this is true —
// a finalized plan alone must NOT advance the pipeline past an unexecuted side
// effect (§12.10).
func BlockHasSucceededEffect(deckDir, slug, blockID string) (bool, error) {
	dir := filepath.Join(PipelineDir(deckDir, slug), "effects")
	entries, err := os.ReadDir(dir)
	if os.IsNotExist(err) {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	for _, e := range entries {
		if e.IsDir() || !strings.HasSuffix(e.Name(), ".json") {
			continue
		}
		data, err := os.ReadFile(filepath.Join(dir, e.Name()))
		if err != nil {
			return false, err
		}
		var eff Effect
		if json.Unmarshal(data, &eff) == nil && eff.BlockID == blockID && eff.Status == EffectSucceeded {
			return true, nil
		}
	}
	return false, nil
}

// NeedsReconcile reports whether an effect is in an ambiguous state that the
// driver MUST reconcile against external state before retrying (§12.7).
func (e Effect) NeedsReconcile() bool {
	return e.Status == EffectExecuting || e.Status == EffectFailed
}

// LoadEffectByDigest reads an effect by its filename digest (as printed by the
// execute command); the bool reports whether it exists.
func LoadEffectByDigest(deckDir, slug, digest string) (Effect, bool, error) {
	path := filepath.Join(PipelineDir(deckDir, slug), "effects", digest+".json")
	data, err := os.ReadFile(path)
	if os.IsNotExist(err) {
		return Effect{}, false, nil
	}
	if err != nil {
		return Effect{}, false, fmt.Errorf("read effect: %w", err)
	}
	var e Effect
	if err := json.Unmarshal(data, &e); err != nil {
		return Effect{}, false, fmt.Errorf("parse effect: %w", err)
	}
	return e, true, nil
}
