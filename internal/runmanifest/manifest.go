package runmanifest

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"parley-deck-cli/internal/fsutil"
	"parley-deck-cli/internal/protocol"
	"parley-deck-cli/internal/runaction"
)

const (
	SchemaVersion = 1
	FileName      = "run.json"

	StatusRunning        = "running"
	StatusWaiting        = "waiting"
	StatusActionRequired = "action_required"
	StatusIncomplete     = "incomplete"
	StatusFailed         = "failed"
	StatusCompleted      = "completed"
	StatusCancelled      = "cancelled"
	StatusStale          = "stale"
)

type Manifest struct {
	SchemaVersion int          `json:"schema_version"`
	RunID         string       `json:"run_id"`
	WorkspaceRoot string       `json:"workspace_root"`
	IdeaSlug      string       `json:"idea_slug"`
	Task          string       `json:"task,omitempty"`
	Mode          string       `json:"mode,omitempty"`
	Transport     string       `json:"transport,omitempty"`
	Status        string       `json:"status,omitempty"`
	Phase         string       `json:"phase,omitempty"`
	IdeaStatus    string       `json:"idea_status,omitempty"`
	CurrentRound  string       `json:"current_round,omitempty"`
	ActiveSteps   []Step       `json:"active_steps,omitempty"`
	LastActionAt  *time.Time   `json:"last_action_at,omitempty"`
	NextActions   []NextAction `json:"next_actions,omitempty"`
	Participants  []string     `json:"participants,omitempty"`
	// RosterSnapshot freezes what each participant ACTUALLY runs, captured at run
	// creation. Before it existed the manifest recorded participant IDs and nothing
	// else, so a finished run could not tell you which model any agent had used — and
	// `continue` re-discovers configuration, so changing a machine default mid-run could
	// silently continue it on a different model. Every later phase of a run uses this
	// snapshot, never a fresh resolve.
	RosterSnapshot []RosterSnapshotEntry `json:"roster_snapshot,omitempty"`
	// RosterRevision is a content hash of the snapshot. `sessions inspect` compares it
	// with the deck's current roster and reports `stale-snapshot` when they differ.
	RosterRevision string    `json:"roster_revision,omitempty"`
	CreatedAt      time.Time `json:"created_at,omitempty"`
	UpdatedAt      time.Time `json:"updated_at,omitempty"`
}

// RosterSnapshotEntry is one participant's effective launch identity at run creation.
// It carries no credentials and no prompt — only what is needed to answer "what did this
// agent actually run?".
type RosterSnapshotEntry struct {
	Agent     string `json:"agent"`
	Adapter   string `json:"adapter"`
	Model     string `json:"model"`
	Effort    string `json:"effort"`
	Speed     string `json:"speed"`
	Auto      bool   `json:"autonomous"`
	Installed bool   `json:"installed"`
	// LaunchArgs is the RESOLVED headless argv at run creation. Without it, Auto is a
	// number the snapshot reports but cannot enforce: a machine-config change that drops
	// --dangerously-skip-permissions would change a continuation's autonomy posture while
	// the frozen row still claimed AUTO=yes. G1 requires auto-args to be pinned too.
	LaunchArgs []string `json:"launch_args,omitempty"`
}

// RosterRevisionOf hashes a snapshot deterministically. Field order and formatting are
// fixed so the same roster always yields the same revision.
func RosterRevisionOf(entries []RosterSnapshotEntry) string {
	if len(entries) == 0 {
		return ""
	}
	sorted := append([]RosterSnapshotEntry(nil), entries...)
	sort.Slice(sorted, func(i, j int) bool { return sorted[i].Agent < sorted[j].Agent })
	h := sha256.New()
	for _, e := range sorted {
		// LaunchArgs is hashed too: it is part of what the run froze, so a change that
		// alters only the argv (e.g. an auto-approve flag added or removed) must make the
		// revision differ. Omitting it let real autonomy drift report `current`.
		fmt.Fprintf(h, "%s\x00%s\x00%s\x00%s\x00%s\x00%t\x00%t\x00%s\n",
			e.Agent, e.Adapter, e.Model, e.Effort, e.Speed, e.Auto, e.Installed,
			strings.Join(e.LaunchArgs, "\x1f"))
	}
	return hex.EncodeToString(h.Sum(nil))[:16]
}

type Step struct {
	ID           string     `json:"id"`
	Kind         string     `json:"kind,omitempty"`
	AgentID      string     `json:"agent_id,omitempty"`
	ArtifactPath string     `json:"artifact_path,omitempty"`
	Status       string     `json:"status,omitempty"`
	StartedAt    *time.Time `json:"started_at,omitempty"`
	UpdatedAt    *time.Time `json:"updated_at,omitempty"`
}

type NextAction = runaction.NextAction

type Options struct {
	Root           string
	RunID          string
	IdeaSlug       string
	Task           string
	Mode           string
	Transport      string
	Status         string
	Phase          string
	IdeaStatus     string
	CurrentRound   string
	ActiveSteps    []Step
	LastActionAt   time.Time
	NextActions    []NextAction
	Participants   []string
	CreatedAt      time.Time
	UpdatedAt      time.Time
	RosterSnapshot []RosterSnapshotEntry
	RosterRevision string
}

func New(opts Options) Manifest {
	createdAt := opts.CreatedAt.UTC()
	if createdAt.IsZero() {
		createdAt = time.Now().UTC()
	}
	updatedAt := opts.UpdatedAt.UTC()
	if updatedAt.IsZero() {
		updatedAt = createdAt
	}
	root := opts.Root
	if abs, err := filepath.Abs(root); err == nil {
		root = abs
	}
	status := opts.Status
	if status == "" {
		status = StatusRunning
	}
	var lastActionAt *time.Time
	if !opts.LastActionAt.IsZero() {
		value := opts.LastActionAt.UTC()
		lastActionAt = &value
	}
	return Manifest{
		SchemaVersion:  SchemaVersion,
		RunID:          opts.RunID,
		WorkspaceRoot:  root,
		IdeaSlug:       opts.IdeaSlug,
		Task:           opts.Task,
		Mode:           opts.Mode,
		Transport:      opts.Transport,
		Status:         status,
		Phase:          opts.Phase,
		IdeaStatus:     opts.IdeaStatus,
		CurrentRound:   opts.CurrentRound,
		ActiveSteps:    append([]Step(nil), opts.ActiveSteps...),
		LastActionAt:   lastActionAt,
		NextActions:    append([]NextAction(nil), opts.NextActions...),
		RosterSnapshot: opts.RosterSnapshot,
		RosterRevision: opts.RosterRevision,
		Participants:   append([]string(nil), opts.Participants...),
		CreatedAt:      createdAt,
		UpdatedAt:      updatedAt,
	}
}

func Path(root, runID string) string {
	if abs, err := filepath.Abs(root); err == nil {
		root = abs
	}
	return filepath.Join(root, protocol.DeckDir, "runs", runID, FileName)
}

func Write(root, runID string, manifest Manifest) error {
	if manifest.SchemaVersion == 0 {
		manifest.SchemaVersion = SchemaVersion
	}
	if manifest.RunID == "" {
		manifest.RunID = runID
	}
	if manifest.Status == "" {
		manifest.Status = StatusRunning
	}
	if manifest.UpdatedAt.IsZero() {
		manifest.UpdatedAt = time.Now().UTC()
	}
	path := Path(root, runID)
	if err := fsutil.MkdirAllResilient(filepath.Dir(path), 0o755); err != nil {
		return err
	}
	data, err := json.MarshalIndent(manifest, "", "  ")
	if err != nil {
		return err
	}
	data = append(data, '\n')
	tmp, err := os.CreateTemp(filepath.Dir(path), ".run.*.tmp")
	if err != nil {
		return err
	}
	tmpName := tmp.Name()
	if _, err := tmp.Write(data); err != nil {
		_ = tmp.Close()
		_ = os.Remove(tmpName)
		return err
	}
	if err := tmp.Close(); err != nil {
		_ = os.Remove(tmpName)
		return err
	}
	return os.Rename(tmpName, path)
}

func Load(root, runID string) (Manifest, error) {
	data, err := os.ReadFile(Path(root, runID))
	if err != nil {
		return Manifest{}, err
	}
	var manifest Manifest
	if err := json.Unmarshal(data, &manifest); err != nil {
		return Manifest{}, err
	}
	if manifest.SchemaVersion == 0 {
		manifest.SchemaVersion = SchemaVersion
	}
	return manifest, nil
}
