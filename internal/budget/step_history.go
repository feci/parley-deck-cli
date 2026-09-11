package budget

import (
	"bufio"
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

// Old Run counters reset on re-entry, so their last value is not a lifetime
// total. A cursor has no native idea identity: resolve it from validated events
// before deciding whether its charges belong to the idea being configured.
// Missing history is not authenticated proof that execution never occurred.
func refuseUnmigratedSteps(roots []string, idea string) error {
	for _, root := range roots {
		runs := filepath.Join(root, "parley-deck", "runs")
		entries, err := stepHistoryDirs(runs)
		if err != nil {
			return err
		}
		for _, entry := range entries {
			if err := inspectStepRun(filepath.Join(runs, entry.Name()), idea); err != nil {
				return err
			}
		}
		if err := inspectStepInvocations(root, idea); err != nil {
			return err
		}
	}
	return nil
}

// Failed actions may never have reached cursor/event publication. Requests at
// those boundaries still establish prior work. The initial readiness/round-one
// calls precede Driver.Run in the existing application and are not old driver
// transitions; all later/unknown phases for this idea require migration.
func inspectStepInvocations(root, idea string) error {
	runtimeDir := filepath.Join(root, ".parley-runtime")
	if info, err := os.Lstat(runtimeDir); err != nil && !os.IsNotExist(err) {
		return err
	} else if err == nil && !info.IsDir() {
		return errors.New("historical runtime storage must be a real directory")
	}
	history := filepath.Join(runtimeDir, "invocations")
	entries, err := stepHistoryDirs(history)
	if err != nil {
		return err
	}
	for _, entry := range entries {
		data, err := readLockOrigin(filepath.Join(history, entry.Name(), "requested.json"))
		if err != nil {
			return fmt.Errorf("historical step invocation is unavailable; migration is required: %w", err)
		}
		var rec struct {
			Metadata *struct {
				Idea  *string `json:"idea"`
				Phase *string `json:"phase"`
			} `json:"metadata"`
		}
		if validateHistoryJSON(json.NewDecoder(bytes.NewReader(data)), 0) != nil || json.Unmarshal(data, &rec) != nil || rec.Metadata == nil || rec.Metadata.Idea == nil {
			return errors.New("unknown historical step invocation identity requires migration")
		}
		if *rec.Metadata.Idea != idea {
			continue
		}
		if rec.Metadata.Phase == nil || (*rec.Metadata.Phase != "preflight" && *rec.Metadata.Phase != "round-01") {
			return errors.New("historical protocol invocation requires step accounting migration")
		}
	}
	return nil
}

func stepHistoryDirs(path string) ([]os.DirEntry, error) {
	info, err := os.Lstat(path)
	if os.IsNotExist(err) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	if !info.IsDir() {
		return nil, errors.New("driver history must be a real directory")
	}
	entries, err := os.ReadDir(path)
	if err != nil {
		return nil, err
	}
	for _, e := range entries {
		if !e.IsDir() {
			return nil, errors.New("driver run history contains a non-directory")
		}
	}
	return entries, nil
}

func inspectStepRun(dir, idea string) error {
	identity := ""
	driverHistory := false
	// Every available identity must agree. Never let the latest run.created
	// reassign earlier charges, or one missing identity hide a cursor.
	bind := func(raw json.RawMessage, required bool) error {
		if raw == nil && !required {
			return nil
		}
		var name string
		if json.Unmarshal(raw, &name) != nil || strings.TrimSpace(name) == "" {
			return errors.New("historical driver idea identity is missing or malformed")
		}
		if identity != "" && identity != name {
			return errors.New("conflicting historical driver idea identities require migration")
		}
		identity = name
		return nil
	}
	cursor, err := readLockOrigin(filepath.Join(dir, "driver.json"))
	if err == nil {
		var fields map[string]json.RawMessage
		if validateHistoryJSON(json.NewDecoder(bytes.NewReader(cursor)), 0) != nil || json.Unmarshal(cursor, &fields) != nil {
			return errors.New("unknown historical driver cursor")
		}
		driverHistory = true
		for _, name := range []string{"idea", "idea_slug"} {
			if err := bind(fields[name], false); err != nil {
				return err
			}
		}
	} else if !os.IsNotExist(err) {
		return err
	}
	data, err := readStepHistoryFile(filepath.Join(dir, "events.jsonl"), 16<<20)
	if err != nil && !os.IsNotExist(err) {
		return err
	}
	scan := bufio.NewScanner(bytes.NewReader(data))
	scan.Buffer(make([]byte, 4096), 1<<20)
	for scan.Scan() {
		line := scan.Bytes()
		if len(bytes.TrimSpace(line)) == 0 {
			continue
		}
		var event struct {
			Type string                     `json:"type"`
			Data map[string]json.RawMessage `json:"data"`
		}
		if validateHistoryJSON(json.NewDecoder(bytes.NewReader(line)), 0) != nil || json.Unmarshal(line, &event) != nil || event.Type == "" {
			return errors.New("unknown historical driver event")
		}
		isDriver := strings.HasPrefix(event.Type, "driver.") || event.Type == "run.phase" || event.Type == "phase.changed" || event.Type == "loop.budget"
		if isDriver || event.Type == "run.created" {
			if err := bind(event.Data["idea"], event.Type == "run.created"); err != nil {
				return err
			}
			driverHistory = driverHistory || isDriver
		}
	}
	if err := scan.Err(); err != nil {
		return err
	}
	if driverHistory && (identity == "" || identity == idea) {
		return errors.New("historical driver activity requires step accounting migration")
	}
	return nil
}

// Read from one bounded regular file. Reject replacement and in-place changes
// observed during the read; never follow an alias into unrelated history.
func readStepHistoryFile(path string, limit int64) ([]byte, error) {
	before, err := os.Lstat(path)
	if err != nil {
		return nil, err
	}
	if !before.Mode().IsRegular() || before.Size() > limit {
		return nil, errors.New("driver history is not a bounded regular file")
	}
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	opened, err := f.Stat()
	if err != nil || !os.SameFile(before, opened) {
		return nil, errors.New("driver history changed during open")
	}
	data, err := io.ReadAll(io.LimitReader(f, limit+1))
	if err != nil {
		return nil, err
	}
	after, statErr := f.Stat()
	current, pathErr := os.Lstat(path)
	if statErr != nil || pathErr != nil || !os.SameFile(opened, current) || !current.Mode().IsRegular() || opened.Size() != after.Size() || !opened.ModTime().Equal(after.ModTime()) {
		return nil, errors.New("driver history changed during read")
	}
	if int64(len(data)) > limit {
		return nil, fmt.Errorf("driver history exceeds %d bytes", limit)
	}
	return data, nil
}
