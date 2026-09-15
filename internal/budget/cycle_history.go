package budget

import (
	"bufio"
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

// CycleKindForPhase only classifies the runner's typed operation. It cannot
// infer whether arbitrary natural-language prompt text asks for a code fix.
func CycleKindForPhase(phase string) (Kind, bool) {
	if phase == "fixup" {
		return Fixup, true
	}
	if n, ok := cycleRound(phase); ok && n > 1 {
		return CrossReview, true
	}
	return "", false
}

// An unknown/manual historical phase could have attempted either kind of work.
// Only the runner's known non-cycle operations may be excluded at bootstrap.
func knownNonCyclePhase(phase string) bool {
	if n, ok := cycleRound(phase); ok && n == 1 {
		return true
	}
	switch phase {
	case "preflight", "runtime-probe", "implementation", "review", "review-consensus",
		"consensus", "final", "consensus-signoff", "review-signoff", "goal-check",
		"evidence-verification", "consult", "steer", "handoff":
		return true
	}
	return false
}

func cycleRound(name string) (int, bool) {
	if !strings.HasPrefix(name, "round-") {
		return 0, false
	}
	digits := strings.TrimPrefix(name, "round-")
	if digits == "" {
		return 0, false
	}
	for _, r := range digits {
		if r < '0' || r > '9' {
			return 0, false
		}
	}
	n, err := strconv.Atoi(digits)
	return n, err == nil && n > 0
}

// LegacyCycleFloor reads the existing driver marker/round boundary. It is only
// a carried lower bound. Other cursor/request history is checked separately;
// this floor cannot authorize forgetting a failed attempt with no marker.
func LegacyCycleFloor(ideaDir string, kind Kind) (int, error) {
	dir := ideaDir
	if kind == Fixup {
		dir = filepath.Join(dir, "review")
	}
	info, err := os.Lstat(dir)
	if os.IsNotExist(err) {
		return 0, nil
	}
	if err != nil {
		return 0, err
	}
	if !info.IsDir() {
		return 0, errors.New("cycle history must be a real directory")
	}
	entries, err := os.ReadDir(dir)
	if err != nil {
		return 0, err
	}
	count := 0
	for _, e := range entries {
		n, round := cycleRound(e.Name())
		if !round {
			continue
		}
		if !e.IsDir() {
			return 0, errors.New("cycle round history must be a real directory")
		}
		if kind == CrossReview {
			if n-1 > count {
				count = n - 1
			}
			continue
		}
		marker := filepath.Join(dir, e.Name(), ".fixup-done")
		info, err := os.Lstat(marker)
		if os.IsNotExist(err) {
			continue
		}
		if err != nil {
			return 0, err
		}
		if !info.Mode().IsRegular() {
			return 0, errors.New("fixup marker must be a regular file")
		}
		count++
	}
	return count, nil
}

func refuseUnmigratedCycles(roots []string, idea, relative string, kind Kind, carried int, currentRun string) error {
	for _, root := range roots {
		floor, err := LegacyCycleFloor(filepath.Join(root, filepath.FromSlash(relative)), kind)
		if err != nil {
			return err
		}
		if floor > carried {
			return errors.New("historical cycle markers exceed the carried count; migration is required")
		}
		runs := filepath.Join(root, "parley-deck", "runs")
		entries, err := stepHistoryDirs(runs)
		if err != nil {
			return err
		}
		for _, entry := range entries {
			dir := filepath.Join(runs, entry.Name())
			identity, eventCount, err := cycleRunEvents(filepath.Join(dir, "events.jsonl"), kind)
			if err != nil {
				return err
			}
			current := sameCycleDirectory(dir, currentRun)
			if identity != "" && identity != idea {
				if current {
					return errors.New("current driver run belongs to another idea")
				}
				continue
			}
			data, err := readLockOrigin(filepath.Join(dir, "driver.json"))
			count := eventCount
			if err == nil {
				var fields map[string]json.RawMessage
				if validateHistoryJSON(json.NewDecoder(bytes.NewReader(data)), 0) != nil || json.Unmarshal(data, &fields) != nil {
					return errors.New("unknown historical cycle cursor")
				}
				field := "fixup_cycles_published"
				if kind == CrossReview {
					field = "rounds_run"
				}
				var value *int
				if json.Unmarshal(fields[field], &value) != nil || value == nil || *value < 0 {
					return errors.New("historical cursor has an unknown charged cycle count")
				}
				n := *value
				if kind == CrossReview && n > 0 {
					n--
				}
				if n > count {
					count = n
				}
			} else if !os.IsNotExist(err) {
				return err
			}
			if count > 0 && (!current || count > carried) {
				return errors.New("another or insufficiently accounted driver run requires cycle migration")
			}
		}
		history := filepath.Join(root, ".parley-runtime", "invocations")
		if info, err := os.Lstat(filepath.Dir(history)); err != nil && !os.IsNotExist(err) {
			return err
		} else if err == nil && !info.IsDir() {
			return errors.New("historical runtime must be a real directory")
		}
		invocations, err := stepHistoryDirs(history)
		if err != nil {
			return err
		}
		for _, entry := range invocations {
			data, err := readLockOrigin(filepath.Join(history, entry.Name(), "requested.json"))
			if err != nil {
				return fmt.Errorf("historical cycle request unavailable: %w", err)
			}
			var record struct {
				Metadata *struct {
					Idea  *string `json:"idea"`
					Phase *string `json:"phase"`
				} `json:"metadata"`
			}
			if validateHistoryJSON(json.NewDecoder(bytes.NewReader(data)), 0) != nil || json.Unmarshal(data, &record) != nil || record.Metadata == nil || record.Metadata.Idea == nil {
				return errors.New("unknown historical cycle invocation identity")
			}
			if *record.Metadata.Idea != idea {
				continue
			}
			if record.Metadata.Phase == nil {
				return errors.New("unknown historical cycle invocation phase")
			}
			if recorded, ok := CycleKindForPhase(*record.Metadata.Phase); ok {
				if recorded == kind {
					return errors.New("historical cycle requests require explicit accounting migration before first binding")
				}
			} else if !knownNonCyclePhase(*record.Metadata.Phase) {
				return errors.New("unknown historical cycle invocation phase requires migration")
			}
		}
	}
	return nil
}

func sameCycleDirectory(a, b string) bool {
	if b == "" {
		return false
	}
	ai, err := os.Stat(a)
	if err != nil {
		return false
	}
	bi, err := os.Stat(b)
	return err == nil && ai.IsDir() && bi.IsDir() && os.SameFile(ai, bi)
}

func cycleRunEvents(path string, kind Kind) (string, int, error) {
	data, err := readStepHistoryFile(path, 16<<20)
	if os.IsNotExist(err) {
		return "", 0, nil
	}
	if err != nil {
		return "", 0, err
	}
	identity := ""
	count := 0
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
		if validateHistoryJSON(json.NewDecoder(bytes.NewReader(line)), 0) != nil || json.Unmarshal(line, &event) != nil {
			return "", 0, errors.New("malformed historical cycle event")
		}
		if event.Type != "run.created" && event.Type != "run.phase" {
			continue
		}
		var name string
		if json.Unmarshal(event.Data["idea"], &name) != nil || name == "" {
			return "", 0, errors.New("historical cycle event lacks idea identity")
		}
		if identity != "" && name != identity {
			return "", 0, errors.New("conflicting historical cycle event identities")
		}
		identity = name
		if event.Type == "run.phase" {
			var action string
			if json.Unmarshal(event.Data["action"], &action) != nil || action == "" {
				return "", 0, errors.New("historical phase event lacks action identity")
			}
			if (kind == Fixup && action == "fixup") || (kind == CrossReview && (action == "promoted" || action == "reopened")) {
				count++
			}
		}
	}
	return identity, count, scan.Err()
}
