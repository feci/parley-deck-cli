package budget

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
)

// LegacyAdoption freezes the actual decision payloads as well as their digest.
// It is separate from charged counts; Unknown history never supplies a floor.
type LegacyAdoption struct {
	Version int                 `json:"version"`
	Scope   string              `json:"scope"`
	History string              `json:"history"`
	Records []LegacyDeclaration `json:"records"`
}

func cycleLegacyHistory(ctx context.Context, root string) (*LegacyAdoption, error) {
	dir, scope, roots, err := legacyScope(ctx, root)
	if err != nil {
		return nil, err
	}
	records, err := readLegacyRecords(dir, scope)
	if err != nil {
		return nil, err
	}
	if len(records) == 0 {
		return nil, nil
	}
	if err := validateLegacyCopies(roots, records); err != nil {
		return nil, err
	}
	return &LegacyAdoption{Version: 1, Scope: scope, History: legacyUnknownHistory, Records: records}, nil
}

func legacyPaths(a *LegacyAdoption) map[string]bool {
	paths := map[string]bool{}
	if a != nil {
		for _, r := range a.Records {
			paths[r.Preview.Path] = true
		}
	}
	return paths
}

func readCycleLegacy(dir, digest string) (*LegacyAdoption, error) {
	raw, err := readStepHistoryFile(filepath.Join(dir, "legacy.json"), 8<<20)
	if digest == "" {
		if !os.IsNotExist(err) {
			return nil, errors.New("cycle policy lost its legacy adoption reference")
		}
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	var a LegacyAdoption
	if err := migrationJSON(raw, &a); err != nil {
		return nil, err
	}
	if a.Version != 1 || a.History != legacyUnknownHistory || len(a.Records) == 0 || len(a.Records) > maxLegacyDecisions || migrationDigest(a) != digest {
		return nil, errors.New("invalid frozen legacy adoption")
	}
	for _, r := range a.Records {
		if err := validateLegacyRecord(r, a.Scope); err != nil {
			return nil, err
		}
	}
	return &a, nil
}

func writeCycleLegacy(dir string, a *LegacyAdoption) (string, error) {
	if a == nil {
		_, err := readCycleLegacy(dir, "")
		return "", err
	}
	digest := migrationDigest(a)
	path := filepath.Join(dir, "legacy.json")
	if _, err := os.Lstat(path); err == nil {
		_, err := readCycleLegacy(dir, digest)
		return digest, err
	} else if !os.IsNotExist(err) {
		return "", err
	}
	raw, err := json.Marshal(a)
	if err != nil {
		return "", err
	}
	if err := writeSynced(path, raw); err != nil {
		return "", err
	}
	saved, err := readStepHistoryFile(path, 8<<20)
	if err != nil || !bytes.Equal(saved, raw) {
		return "", errors.New("legacy adoption readback failed")
	}
	return digest, nil
}

func validateCycleLegacy(ctx context.Context, root, dir, digest string) (*LegacyAdoption, error) {
	a, err := readCycleLegacy(dir, digest)
	if err != nil || a == nil {
		return a, err
	}
	if root == "" {
		return nil, errors.New("legacy adoption requires its runtime origin")
	}
	liveDir, scope, roots, err := legacyScope(ctx, root)
	if err != nil {
		return nil, err
	}
	if a.Scope != scope {
		return nil, errors.New("legacy adoption scope mismatch")
	}
	live, err := readLegacyRecords(liveDir, scope)
	if err != nil {
		return nil, err
	}
	byID := map[string]string{}
	for _, r := range live {
		byID[r.Request.DecisionID] = migrationDigest(r)
	}
	for _, r := range a.Records {
		if byID[r.Request.DecisionID] != migrationDigest(r) {
			return nil, errors.New("adopted legacy decision is missing or changed")
		}
	}
	if err := validateLegacyCopies(roots, a.Records); err != nil {
		return nil, err
	}
	// A bound idea cannot borrow a later declaration without explicit migration.
	// Known scoped charges still belong to their existing ledgers. This scan only
	// retains the original refusal on additional unscoped/corrupt cycle events.
	paths := legacyPaths(a)
	for _, origin := range roots {
		if err := migrationDirectory(origin, "parley-deck/runs"); err != nil {
			return nil, err
		}
		entries, err := stepHistoryDirs(filepath.Join(origin, "parley-deck", "runs"))
		if err != nil {
			return nil, err
		}
		for _, e := range entries {
			if paths[unscopedRunPrefix+e.Name()] {
				continue
			}
			if _, _, err := cycleRunEvents(filepath.Join(origin, "parley-deck", "runs", e.Name(), "events.jsonl"), CrossReview); err != nil {
				return nil, err
			}
		}
	}
	return a, nil
}
