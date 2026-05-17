package sessionstore

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"sort"
	"time"
)

const (
	SchemaVersion = 1
	homeEnv       = "PARLEY_HOME"
)

type Session struct {
	WorkspaceRoot string    `json:"workspace_root"`
	RunID         string    `json:"run_id"`
	IdeaSlug      string    `json:"idea_slug"`
	Task          string    `json:"task,omitempty"`
	Participants  []string  `json:"participants,omitempty"`
	CreatedAt     time.Time `json:"created_at,omitempty"`
	UpdatedAt     time.Time `json:"updated_at,omitempty"`
	LastEventAt   time.Time `json:"last_event_at,omitempty"`
	Terminal      bool      `json:"terminal"`
}

type Store struct {
	path string
}

type document struct {
	SchemaVersion int       `json:"schema_version"`
	Sessions      []Session `json:"sessions"`
}

func DefaultDir() (string, error) {
	if override := os.Getenv(homeEnv); override != "" {
		return override, nil
	}
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(home, ".parley-deck"), nil
}

func Default() (Store, error) {
	dir, err := DefaultDir()
	if err != nil {
		return Store{}, err
	}
	return New(filepath.Join(dir, "sessions.json")), nil
}

func New(path string) Store {
	return Store{path: path}
}

func (s Store) Path() string {
	return s.path
}

func (s Store) List() ([]Session, error) {
	doc, err := s.read()
	if err != nil {
		return nil, err
	}
	sessions := append([]Session(nil), doc.Sessions...)
	sort.SliceStable(sessions, func(i, j int) bool {
		left := sessions[i].UpdatedAt
		if left.IsZero() {
			left = sessions[i].CreatedAt
		}
		right := sessions[j].UpdatedAt
		if right.IsZero() {
			right = sessions[j].CreatedAt
		}
		return left.After(right)
	})
	return sessions, nil
}

func (s Store) Upsert(session Session) error {
	if session.WorkspaceRoot == "" || session.RunID == "" {
		return errors.New("workspace root and run id are required")
	}
	if abs, err := filepath.Abs(session.WorkspaceRoot); err == nil {
		session.WorkspaceRoot = abs
	}
	now := time.Now().UTC()
	if session.UpdatedAt.IsZero() {
		session.UpdatedAt = now
	}

	doc, err := s.read()
	if err != nil {
		return err
	}
	replaced := false
	for i, existing := range doc.Sessions {
		if existing.WorkspaceRoot != session.WorkspaceRoot || existing.RunID != session.RunID {
			continue
		}
		if session.Task == "" {
			session.Task = existing.Task
		}
		if len(session.Participants) == 0 {
			session.Participants = existing.Participants
		}
		if session.IdeaSlug == "" {
			session.IdeaSlug = existing.IdeaSlug
		}
		if session.CreatedAt.IsZero() {
			session.CreatedAt = existing.CreatedAt
		}
		doc.Sessions[i] = session
		replaced = true
		break
	}
	if !replaced {
		if session.CreatedAt.IsZero() {
			session.CreatedAt = now
		}
		doc.Sessions = append(doc.Sessions, session)
	}
	return s.write(doc)
}

func (s Store) read() (document, error) {
	doc := document{SchemaVersion: SchemaVersion}
	data, err := os.ReadFile(s.path)
	if errors.Is(err, os.ErrNotExist) {
		return doc, nil
	}
	if err != nil {
		return document{}, err
	}
	if err := json.Unmarshal(data, &doc); err != nil {
		return document{}, err
	}
	if doc.SchemaVersion == 0 {
		doc.SchemaVersion = SchemaVersion
	}
	return doc, nil
}

func (s Store) write(doc document) error {
	if doc.SchemaVersion == 0 {
		doc.SchemaVersion = SchemaVersion
	}
	if err := os.MkdirAll(filepath.Dir(s.path), 0o755); err != nil {
		return err
	}
	data, err := json.MarshalIndent(doc, "", "  ")
	if err != nil {
		return err
	}
	data = append(data, '\n')
	tmp, err := os.CreateTemp(filepath.Dir(s.path), ".sessions.*.tmp")
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
	return os.Rename(tmpName, s.path)
}
