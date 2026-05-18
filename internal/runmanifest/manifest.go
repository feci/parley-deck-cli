package runmanifest

import (
	"encoding/json"
	"os"
	"path/filepath"
	"time"

	"parley-deck-cli/internal/protocol"
)

const (
	SchemaVersion = 1
	FileName      = "run.json"

	StatusRunning = "running"
)

type Manifest struct {
	SchemaVersion int       `json:"schema_version"`
	RunID         string    `json:"run_id"`
	WorkspaceRoot string    `json:"workspace_root"`
	IdeaSlug      string    `json:"idea_slug"`
	Task          string    `json:"task,omitempty"`
	Mode          string    `json:"mode,omitempty"`
	Transport     string    `json:"transport,omitempty"`
	Status        string    `json:"status,omitempty"`
	Participants  []string  `json:"participants,omitempty"`
	CreatedAt     time.Time `json:"created_at,omitempty"`
	UpdatedAt     time.Time `json:"updated_at,omitempty"`
}

type Options struct {
	Root         string
	RunID        string
	IdeaSlug     string
	Task         string
	Mode         string
	Transport    string
	Status       string
	Participants []string
	CreatedAt    time.Time
	UpdatedAt    time.Time
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
	return Manifest{
		SchemaVersion: SchemaVersion,
		RunID:         opts.RunID,
		WorkspaceRoot: root,
		IdeaSlug:      opts.IdeaSlug,
		Task:          opts.Task,
		Mode:          opts.Mode,
		Transport:     opts.Transport,
		Status:        status,
		Participants:  append([]string(nil), opts.Participants...),
		CreatedAt:     createdAt,
		UpdatedAt:     updatedAt,
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
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
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
