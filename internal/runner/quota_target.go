package runner

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"parley-deck-cli/internal/protocol"
)

// CanonicalArtifactIdea binds only an explicit artifact target under this deck's
// ideas directory. Prompt text and arbitrary unbound work are never inferred.
func CanonicalArtifactIdea(root, target string) (string, error) {
	if target == "" {
		return "", nil
	}
	base, err := filepath.Abs(filepath.Join(root, protocol.DeckDir, "ideas"))
	if err != nil {
		return "", err
	}
	if !filepath.IsAbs(target) {
		target = filepath.Join(root, target)
	}
	target, err = filepath.Abs(target)
	if err != nil {
		return "", err
	}
	rel, err := filepath.Rel(base, target)
	if err != nil {
		return "", err
	}
	if rel == ".." || strings.HasPrefix(rel, ".."+string(filepath.Separator)) {
		return "", nil
	}
	parts := strings.Split(filepath.ToSlash(rel), "/")
	if len(parts) < 2 || parts[0] == "." || parts[0] == "" {
		return "", nil
	}
	// Symlinked target ancestry can name a different idea. Refuse ambiguous
	// canonical targets instead of assigning a launch to the wrong membership.
	for dir := filepath.Dir(target); dir != base; dir = filepath.Dir(dir) {
		st, e := os.Lstat(dir)
		if e != nil {
			if os.IsNotExist(e) {
				continue
			}
			return "", e
		}
		if st.Mode()&os.ModeSymlink != 0 {
			return "", fmt.Errorf("canonical artifact target has symlink ancestry")
		}
	}
	return parts[0], nil
}
