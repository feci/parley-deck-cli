package quota

import (
	"bytes"
	"crypto/sha256"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"

	"gopkg.in/yaml.v3"
)

// Authority binds a verbatim user answer to committed bytes. The working inbox
// may subsequently be archived or deleted; the Git object remains the evidence.
// This checks attribution and content, not human identity or truth of testimony.
type Authority struct {
	Path   string `json:"path"`
	Commit string `json:"commit"`
	Blob   string `json:"blob"`
	SHA256 string `json:"sha256"`
	Quote  string `json:"quote"`
}

var objectID = regexp.MustCompile(`^(?:[a-f0-9]{40}|[a-f0-9]{64})$`)
var identifier = regexp.MustCompile(`^[A-Za-z0-9][A-Za-z0-9._-]*$`)

func git(root string, args ...string) ([]byte, error) {
	cmd := exec.Command("git", append([]string{"-C", root}, args...)...)
	out, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("committed authority unavailable: %w", err)
	}
	return out, nil
}

func authorityBytes(root string, a Authority) ([]byte, error) {
	clean := filepath.ToSlash(filepath.Clean(a.Path))
	if clean != a.Path || filepath.IsAbs(a.Path) || strings.HasPrefix(a.Path, "../") || !strings.HasPrefix(a.Path, "parley-deck/inbox/") || !strings.HasPrefix(filepath.Base(a.Path), "user-to-") || !strings.HasSuffix(a.Path, ".md") || !objectID.MatchString(a.Commit) {
		return nil, fmt.Errorf("invalid committed user-answer path or commit")
	}
	kind, err := git(root, "cat-file", "-t", a.Commit)
	if err != nil {
		return nil, err
	}
	if strings.TrimSpace(string(kind)) != "commit" {
		return nil, fmt.Errorf("owner answer requires a commit object")
	}
	blob, err := git(root, "rev-parse", "--verify", a.Commit+":"+a.Path)
	if err != nil {
		return nil, err
	}
	if a.Blob != "" && strings.TrimSpace(string(blob)) != a.Blob {
		return nil, fmt.Errorf("user-answer blob changed")
	}
	data, err := git(root, "show", a.Commit+":"+a.Path)
	if err != nil {
		return nil, err
	}
	// An existing contradictory copy is never silently bypassed with old Git bytes.
	for _, path := range []string{a.Path, filepath.ToSlash(filepath.Join("parley-deck/inbox/archived", filepath.Base(a.Path)))} {
		live, err := os.ReadFile(filepath.Join(root, filepath.FromSlash(path)))
		if err == nil && !bytes.Equal(live, data) {
			return nil, fmt.Errorf("user-answer working copy changed: %s", path)
		}
		if err != nil && !os.IsNotExist(err) {
			return nil, err
		}
	}
	return data, nil
}

func BindAuthority(root, idea, path, commit, quote string) (Authority, error) {
	a := Authority{Path: path, Commit: commit, Quote: quote}
	raw, err := authorityBytes(root, a)
	if err != nil {
		return a, err
	}
	blob, err := git(root, "rev-parse", "--verify", commit+":"+path)
	if err != nil {
		return a, err
	}
	a.Blob = strings.TrimSpace(string(blob))
	a.SHA256 = fmt.Sprintf("%x", sha256.Sum256(raw))
	return a, ValidateAuthority(root, idea, a)
}

func ValidateAuthority(root, idea string, a Authority) error {
	if !objectID.MatchString(a.Blob) || len(a.SHA256) != 64 || strings.TrimSpace(a.Quote) == "" {
		return fmt.Errorf("incomplete user authority")
	}
	raw, err := authorityBytes(root, a)
	if err != nil {
		return err
	}
	if fmt.Sprintf("%x", sha256.Sum256(raw)) != a.SHA256 {
		return fmt.Errorf("user-answer digest changed")
	}
	parts := strings.SplitN(string(raw), "---", 3)
	if len(parts) != 3 || parts[0] != "" {
		return fmt.Errorf("missing user-answer frontmatter")
	}
	var meta map[string]string
	if yaml.Unmarshal([]byte(parts[1]), &meta) != nil || meta["from"] != "user" || meta["idea"] != idea || !strings.Contains(parts[2], a.Quote) {
		return fmt.Errorf("wrong user-answer attribution, idea or verbatim quote")
	}
	return nil
}

func IdeaRoot(ideaDir string) string {
	dir := filepath.Clean(ideaDir)
	for filepath.Base(dir) != "parley-deck" && filepath.Dir(dir) != dir {
		dir = filepath.Dir(dir)
	}
	return filepath.Dir(dir)
}
