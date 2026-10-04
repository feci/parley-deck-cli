package store

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"

	"parley-deck-cli/internal/fsutil"
)

type Event struct {
	Time time.Time      `json:"time"`
	Type string         `json:"type"`
	Data map[string]any `json:"data,omitempty"`
}

type Store struct {
	dir string
}

var appendMu sync.Mutex

func New(dir string) Store {
	return Store{dir: dir}
}

// Enabled distinguishes a real run store from an optional zero-value sink.
func (s Store) Enabled() bool { return s.dir != "" }

// Directory exposes the actual run directory to cooperating runtime gates.
func (s Store) Directory() string { return s.dir }

func NewRunID(t time.Time) string {
	return t.UTC().Format("20060102T150405.000000000Z")
}

func (s Store) Append(event Event) error { return s.appendEvent(event, false) }

func (s Store) AppendDurable(event Event) error { return s.appendEvent(event, true) }

func (s Store) appendEvent(event Event, durable bool) error {
	appendMu.Lock()
	defer appendMu.Unlock()

	if event.Time.IsZero() {
		event.Time = time.Now().UTC()
	}
	if err := fsutil.MkdirAllResilient(s.dir, 0o755); err != nil {
		return err
	}
	file, err := os.OpenFile(filepath.Join(s.dir, "events.jsonl"), os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0o644)
	if err != nil {
		return err
	}
	defer file.Close()

	encoded, err := json.Marshal(event)
	if err != nil {
		return err
	}
	if n, err := file.Write(append(encoded, '\n')); err != nil {
		return err
	} else if n != len(encoded)+1 {
		return fmt.Errorf("short event write")
	}
	if durable {
		if err := fsutil.SyncFile(file); err != nil {
			return err
		}
		if err := file.Close(); err != nil {
			return err
		}
		d, err := os.Open(s.dir)
		if err != nil {
			return err
		}
		defer d.Close()
		return fsutil.SyncFile(d)
	}
	return nil
}

func (s Store) Load() ([]Event, error) {
	path := filepath.Join(s.dir, "events.jsonl")
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	scanner := bufio.NewScanner(bytes.NewReader(data))
	var events []Event
	lineNo := 0
	for scanner.Scan() {
		lineNo++
		line := scanner.Bytes()
		if len(bytes.TrimSpace(line)) == 0 {
			continue
		}
		var event Event
		if err := json.Unmarshal(line, &event); err != nil {
			return nil, fmt.Errorf("%s:%d: %w", path, lineNo, err)
		}
		events = append(events, event)
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}
	return events, nil
}

// Sync establishes the durability barrier again when a transition replay finds
// its terminal evaluation already appended before an earlier failed sync.
func (s Store) Sync() error {
	f, err := os.Open(filepath.Join(s.dir, "events.jsonl"))
	if err != nil {
		return err
	}
	err = fsutil.SyncFile(f)
	closeErr := f.Close()
	if err != nil {
		return err
	}
	if closeErr != nil {
		return closeErr
	}
	d, err := os.Open(s.dir)
	if err != nil {
		return err
	}
	defer d.Close()
	return fsutil.SyncFile(d)
}
