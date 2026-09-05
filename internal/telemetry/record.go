// Package telemetry records content-free launch evidence. Provider estimates are
// not invoices, and absent observations are represented by JSON nulls.
package telemetry

import (
	"crypto/rand"
	"encoding/json"
	"errors"
	"fmt"
	"math"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"sync"
	"time"

	"parley-deck-cli/internal/fsutil"
)

const SchemaVersion = 1

type Context struct {
	Mode           string  `json:"context_mode"`
	SourceSHA256   *string `json:"source_sha256"`
	PacketSHA256   *string `json:"packet_sha256"`
	FallbackReason *string `json:"fallback_reason"`
}

type Metadata struct {
	RunID           string  `json:"run_id"`
	SegmentID       string  `json:"segment_id"`
	Idea            string  `json:"idea"`
	Phase           string  `json:"phase"`
	Agent           string  `json:"agent"`
	Adapter         string  `json:"adapter"`
	LaunchMode      string  `json:"launch_mode"`
	AttemptOrdinal  int     `json:"attempt_ordinal"`
	RetryOf         *string `json:"retry_of"`
	RequestedModel  *string `json:"requested_model"`
	RequestedEffort *string `json:"requested_effort"`
	RequestedSpeed  *string `json:"requested_speed"`
	Context         Context `json:"context"`
}

type Usage struct {
	InputTokens      *int64   `json:"input_tokens"`
	OutputTokens     *int64   `json:"output_tokens"`
	CacheReadTokens  *int64   `json:"cache_read_tokens"`
	CacheWriteTokens *int64   `json:"cache_write_tokens"`
	TotalTokens      *int64   `json:"total_tokens"`
	CostUSD          *float64 `json:"cost_usd"`
	ReportedModel    *string  `json:"reported_model"`
	ReportedModels   []string `json:"reported_models"`
	Source           string   `json:"source"`
	CostBasis        string   `json:"cost_basis"`
	Coverage         string   `json:"coverage"`
}

type Observation struct {
	FirstActivityMS *int64 `json:"first_activity_ms"`
	StdoutBytes     int64  `json:"stdout_bytes"`
	StderrBytes     int64  `json:"stderr_bytes"`
	TruncatedInput  bool   `json:"parser_input_truncated"`
}

type Outcome struct {
	Status         string      `json:"status"`
	ExitCode       *int        `json:"exit_code"`
	FailureClass   *string     `json:"failure_class"`
	ArtifactSHA256 *string     `json:"artifact_sha256"`
	Usage          Usage       `json:"usage"`
	Observation    Observation `json:"observation"`
}

type Record struct {
	SchemaVersion int        `json:"schema_version"`
	Type          string     `json:"type"`
	InvocationID  string     `json:"invocation_id"`
	Metadata      Metadata   `json:"metadata"`
	RequestedAt   time.Time  `json:"requested_at"`
	StartedAt     *time.Time `json:"started_at"`
	CompletedAt   *time.Time `json:"completed_at"`
	DurationMS    *int64     `json:"duration_ms"`
	PID           *int       `json:"pid"`
	Outcome       *Outcome   `json:"outcome"`
	Warnings      []string   `json:"warnings"`
}

// Invocation owns one unique directory. It never appends competing process data
// to a shared file, and cannot silently overwrite an earlier terminal result.
type Invocation struct {
	mu        sync.Mutex
	Dir       string
	ID        string
	record    Record
	began     time.Time
	finished  bool
	finishErr error
}

var labelPattern = regexp.MustCompile(`^[A-Za-z0-9][A-Za-z0-9._:/+@-]{0,179}$`)
var secretPattern = regexp.MustCompile(`(?i)((^|[/.:@])(npm_|gh[pousr]_|github_pat_|sk-|xox[baprs]-|AKIA[0-9A-Z]|eyJ[A-Za-z0-9_-]+\.)|bearer|password|api[_-]?key|authorization)`)
var opaqueSecretPattern = regexp.MustCompile(`^(?:[a-fA-F0-9]{32,}|[A-Za-z0-9_]{48,})$`)
var hashPattern = regexp.MustCompile(`^[0-9a-f]{64}$`)

func String(value string) *string {
	if value == "" {
		return nil
	}
	return &value
}

// SafeLabel accepts identifiers, not arbitrary diagnostic prose or URLs.
func SafeLabel(value string) *string {
	if !labelPattern.MatchString(value) || strings.Contains(value, "://") || secretPattern.MatchString(value) || opaqueSecretPattern.MatchString(value) {
		return nil
	}
	return &value
}

func safeHash(value *string) *string {
	if value == nil || !hashPattern.MatchString(*value) {
		return nil
	}
	return clone(value)
}

func clone[T any](value *T) *T {
	if value == nil {
		return nil
	}
	copy := *value
	return &copy
}

func NewID() (string, error) {
	var raw [16]byte
	if _, err := rand.Read(raw[:]); err != nil {
		return "", err
	}
	raw[6] = raw[6]&0x0f | 0x40
	raw[8] = raw[8]&0x3f | 0x80
	return fmt.Sprintf("%x-%x-%x-%x-%x", raw[0:4], raw[4:6], raw[6:8], raw[8:10], raw[10:16]), nil
}

func Begin(directory string, metadata Metadata) (*Invocation, error) {
	id, err := NewID()
	if err != nil {
		return nil, err
	}
	if err := fsutil.MkdirAllResilient(directory, 0o700); err != nil {
		return nil, fmt.Errorf("create telemetry directory: %w", err)
	}
	info, err := os.Lstat(directory)
	if err != nil || !info.IsDir() || info.Mode()&os.ModeSymlink != 0 {
		return nil, errors.New("telemetry directory must be a real directory")
	}
	dir := filepath.Join(directory, id)
	if err := os.Mkdir(dir, 0o700); err != nil {
		return nil, fmt.Errorf("reserve invocation directory: %w", err)
	}
	warnings := []string{}
	metadata = cleanMetadata(metadata, &warnings)
	now := time.Now()
	i := &Invocation{Dir: dir, ID: id, began: now, record: Record{
		SchemaVersion: SchemaVersion, Type: "invocation.requested", InvocationID: id,
		Metadata: metadata, RequestedAt: now.UTC(), Warnings: warnings,
	}}
	if err := i.write("requested.json"); err != nil {
		return nil, fmt.Errorf("persist requested invocation before launch: %w", err)
	}
	return i, nil
}

func cleanMetadata(m Metadata, warnings *[]string) Metadata {
	fields := []struct {
		name   string
		target *string
	}{
		{"run_id", &m.RunID}, {"segment_id", &m.SegmentID}, {"idea", &m.Idea},
		{"phase", &m.Phase}, {"agent", &m.Agent}, {"adapter", &m.Adapter},
		{"launch_mode", &m.LaunchMode}, {"context_mode", &m.Context.Mode},
	}
	for _, field := range fields {
		if *field.target != "" && SafeLabel(*field.target) == nil {
			*field.target = "unknown"
			*warnings = append(*warnings, "unsafe_"+field.name+"_omitted")
		}
	}
	optional := []struct {
		name   string
		target **string
	}{
		{"requested_model", &m.RequestedModel}, {"requested_effort", &m.RequestedEffort},
		{"requested_speed", &m.RequestedSpeed}, {"retry_of", &m.RetryOf},
		{"fallback_reason", &m.Context.FallbackReason},
	}
	for _, field := range optional {
		if *field.target != nil {
			clean := SafeLabel(**field.target)
			if clean == nil {
				*warnings = append(*warnings, "unsafe_"+field.name+"_omitted")
			}
			*field.target = clean
		}
	}
	m.Context.SourceSHA256 = safeHash(m.Context.SourceSHA256)
	m.Context.PacketSHA256 = safeHash(m.Context.PacketSHA256)
	if m.AttemptOrdinal < 1 {
		m.AttemptOrdinal = 1
	}
	return m
}

func (i *Invocation) Started(pid int) error {
	i.mu.Lock()
	defer i.mu.Unlock()
	if i.finished || i.record.StartedAt != nil {
		return errors.New("invocation already started or completed")
	}
	now := time.Now().UTC()
	i.record.Type = "invocation.started"
	i.record.StartedAt, i.record.PID = &now, &pid
	return i.write("started.json")
}

func CleanUsage(usage Usage) Usage {
	for _, ptr := range []**int64{&usage.InputTokens, &usage.OutputTokens, &usage.CacheReadTokens, &usage.CacheWriteTokens, &usage.TotalTokens} {
		if *ptr != nil && (**ptr < 0 || **ptr > 1_000_000_000_000) {
			*ptr = nil
		}
		*ptr = clone(*ptr)
	}
	if usage.CostUSD != nil && (math.IsNaN(*usage.CostUSD) || math.IsInf(*usage.CostUSD, 0) || *usage.CostUSD < 0 || *usage.CostUSD > 1_000_000) {
		usage.CostUSD = nil
	}
	usage.CostUSD = clone(usage.CostUSD)
	if usage.ReportedModel != nil {
		usage.ReportedModel = SafeLabel(*usage.ReportedModel)
	}
	models := []string{}
	for _, model := range usage.ReportedModels {
		if safe := SafeLabel(model); safe != nil {
			models = append(models, *safe)
		}
	}
	usage.ReportedModels = models
	for _, field := range []*string{&usage.Source, &usage.CostBasis, &usage.Coverage} {
		if SafeLabel(*field) == nil {
			*field = "unavailable"
		}
	}
	return usage
}

func (i *Invocation) Finish(outcome Outcome) error {
	i.mu.Lock()
	defer i.mu.Unlock()
	if i.finished {
		return i.finishErr
	}
	i.finished = true
	if SafeLabel(outcome.Status) == nil {
		outcome.Status = "unknown"
	}
	if outcome.FailureClass != nil {
		outcome.FailureClass = SafeLabel(*outcome.FailureClass)
	}
	outcome.ArtifactSHA256 = safeHash(outcome.ArtifactSHA256)
	outcome.ExitCode = clone(outcome.ExitCode)
	outcome.Observation.FirstActivityMS = clone(outcome.Observation.FirstActivityMS)
	outcome.Usage = CleanUsage(outcome.Usage)
	now := time.Now().UTC()
	duration := time.Since(i.began).Milliseconds()
	i.record.Type = "invocation.terminal"
	i.record.CompletedAt, i.record.DurationMS, i.record.Outcome = &now, &duration, &outcome
	i.finishErr = i.write("terminal.json")
	return i.finishErr
}

func (i *Invocation) Snapshot() Record {
	i.mu.Lock()
	defer i.mu.Unlock()
	// All persisted fields have already been reduced to JSON-safe scalar types.
	data, _ := json.Marshal(i.record)
	var copy Record
	_ = json.Unmarshal(data, &copy)
	return copy
}

func (i *Invocation) write(name string) error {
	data, err := json.MarshalIndent(i.record, "", "  ")
	if err != nil {
		return err
	}
	path := filepath.Join(i.Dir, name)
	file, err := os.OpenFile(path, os.O_CREATE|os.O_EXCL|os.O_WRONLY, 0o600)
	if err != nil {
		return err
	}
	if _, err := file.Write(append(data, '\n')); err != nil {
		file.Close()
		return err
	}
	if err := file.Sync(); err != nil {
		file.Close()
		return err
	}
	return file.Close()
}
