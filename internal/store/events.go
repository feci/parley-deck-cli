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

func NewRunID(t time.Time) string {
	return t.UTC().Format("20060102T150405.000000000Z")
}

func (s Store) Append(event Event) error {
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
	if _, err := file.Write(append(encoded, '\n')); err != nil {
		return err
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
