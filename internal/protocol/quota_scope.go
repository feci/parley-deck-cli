package protocol

import (
	"fmt"
	"os"
	"path/filepath"
)

// QuotaLeaseRoot resolves workspace ancestors, then requires the idea scope to
// resolve to that exact join, as evidence verification does. Never discover a
// different deck by searching the resolved path's ancestors.
func QuotaLeaseRoot(ideaDir string) (string, error) {
	dir, err := filepath.Abs(ideaDir)
	if err != nil {
		return "", err
	}
	ideas := filepath.Dir(dir)
	deck := filepath.Dir(ideas)
	if filepath.Base(ideas) != "ideas" || filepath.Base(deck) != DeckDir {
		return "", fmt.Errorf("quota lease requires parley-deck/ideas/<idea>")
	}
	return quotaScopeRoot(filepath.Dir(deck), filepath.Base(dir))
}

// An empty idea checks creation before any kickoff writes. Missing directories
// may be created, but an existing alias (including a dangling link) is refused.
func quotaScopeRoot(root, idea string) (string, error) {
	root, err := filepath.Abs(root)
	if err != nil {
		return "", err
	}
	root, err = filepath.EvalSymlinks(root)
	if err != nil {
		return "", err
	}
	parts := []string{DeckDir, "ideas"}
	if idea != "" {
		parts = append(parts, idea)
	}
	expected := root
	for _, part := range parts {
		expected = filepath.Join(expected, part)
		st, err := os.Lstat(expected)
		if os.IsNotExist(err) && idea == "" {
			continue
		}
		if err != nil {
			return "", err
		}
		actual, err := filepath.EvalSymlinks(expected)
		if err != nil || actual != expected || !st.IsDir() {
			return "", fmt.Errorf("quota leases do not support symlinked deck or idea scopes: %s; use a physical deck or disable quota_auto_exclude for ordinary driving", expected)
		}
	}
	return root, nil
}
