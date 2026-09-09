package evidence

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strings"
)

// tree.go computes the tested code-tree identity (evidence-first-efficiency).
//
// The digest covers the exact files a reviewer means by "the code tree":
// tracked files at their working-tree content PLUS real untracked files
// (git ls-files -c -o --exclude-standard), minus only the explicit evidence
// artifacts passed in `excludeRel` (e.g. the idea's EVIDENCE.json) so that
// persisting evidence does not invalidate itself. It does NOT exclude the
// whole deck, all Markdown, or any other broad class. Symlinks that escape
// the root are an error, never silently followed.

// ReviewedCommit returns HEAD's commit id, or "" when the root is not a git
// work tree (the digest still works; the report records the absence).
func ReviewedCommit(root string) string {
	out, err := exec.Command("git", "-C", root, "rev-parse", "HEAD").Output()
	if err != nil {
		return ""
	}
	return strings.TrimSpace(string(out))
}

// TreeDirty reports whether the working tree differs from HEAD for tracked
// files. Honest dirty handling: a dirty tree still gets a digest (the tested
// content is what it is), but the report flags it so closure can decide.
func TreeDirty(root string) bool {
	if ReviewedCommit(root) == "" {
		return true // no commit baseline → cannot claim clean
	}
	err := exec.Command("git", "-C", root, "diff", "--quiet", "HEAD", "--").Run()
	return err != nil
}

// TreeDigest hashes the working-tree identity of every in-scope file under
// root. excludeRel is an exact-match set of slash-separated paths relative to
// root (explicit evidence artifacts only).
//
// The digest binds MORE than content: each entry contributes its type, its
// permission mode, and — for symlinks — its raw link target, so retargeting a
// symlink (a→b) or flipping 0600→0700 changes the digest even when the byte
// content is identical. Vanishing/unreadable/unsupported entries (devices,
// sockets, fifos) are errors, never silently skipped: a tree that cannot be
// fully identified has no valid digest. Symlinks are hashed as links, never
// followed; one whose resolved target escapes the (canonicalized) root is an
// error. The root is canonicalized (macOS /var vs /private/var) before the
// escape check so a symlink resolving through a path alias is judged against
// the real tree location.
func TreeDigest(root string, excludeRel ...string) (string, error) {
	canonicalRoot, err := filepath.EvalSymlinks(root)
	if err != nil {
		return "", fmt.Errorf("tree digest: cannot canonicalize root %s: %w", root, err)
	}
	files, err := listTreeFiles(root)
	if err != nil {
		return "", err
	}
	excluded := map[string]bool{}
	for _, e := range excludeRel {
		excluded[filepath.ToSlash(e)] = true
	}
	h := sha256.New()
	count := 0
	for _, rel := range files {
		if excluded[rel] {
			continue
		}
		abs := filepath.Join(root, filepath.FromSlash(rel))
		fi, err := os.Lstat(abs)
		if err != nil {
			return "", fmt.Errorf("tree digest: cannot lstat %s: %w", rel, err)
		}
		switch {
		case fi.Mode()&os.ModeSymlink != 0:
			// Hash the link itself: type + raw target. Retargeting changes the
			// digest. The resolved target is used ONLY for the escape check,
			// against the canonical root.
			rawTarget, err := os.Readlink(abs)
			if err != nil {
				return "", fmt.Errorf("tree digest: cannot read symlink %s: %w", rel, err)
			}
			target, err := filepath.EvalSymlinks(abs)
			if err != nil {
				return "", fmt.Errorf("tree digest: cannot resolve symlink %s: %w", rel, err)
			}
			r, err := filepath.Rel(canonicalRoot, target)
			if err != nil || r == ".." || strings.HasPrefix(r, ".."+string(filepath.Separator)) || filepath.IsAbs(r) {
				return "", fmt.Errorf("tree digest: unsafe symlink %s escapes the tree root", rel)
			}
			fmt.Fprintf(h, "link:%o:%s\n%s\n", fi.Mode().Perm(), rel, rawTarget)
			count++
		case fi.Mode().IsRegular():
			f, err := os.Open(abs)
			if err != nil {
				return "", fmt.Errorf("tree digest: %s: %w", rel, err)
			}
			content, err := io.ReadAll(f)
			_ = f.Close()
			if err != nil {
				return "", fmt.Errorf("tree digest: %s: %w", rel, err)
			}
			fmt.Fprintf(h, "file:%o:%d:%s\n", fi.Mode().Perm(), len(content), rel)
			h.Write(content)
			count++
		default:
			return "", fmt.Errorf("tree digest: unsupported entry type %s (%s)", rel, fi.Mode().Type())
		}
	}
	if count == 0 {
		return "", fmt.Errorf("tree digest: no in-scope files under %s", root)
	}
	return hex.EncodeToString(h.Sum(nil)), nil
}

// listTreeFiles returns slash-separated paths relative to root. In a git work
// tree it uses `git ls-files -c -o --exclude-standard` (tracked + real
// untracked, honoring .gitignore so dependency/cache directories stay out of
// scope without broad content-class exclusions). Outside git it falls back to
// a filesystem walk that skips only .git.
func listTreeFiles(root string) ([]string, error) {
	if ReviewedCommit(root) != "" {
		out, err := exec.Command("git", "-C", root, "ls-files", "-c", "-o", "--exclude-standard", "-z").Output()
		if err != nil {
			return nil, fmt.Errorf("git ls-files: %w", err)
		}
		var files []string
		for _, p := range strings.Split(string(out), "\x00") {
			if p != "" {
				files = append(files, p)
			}
		}
		sort.Strings(files)
		return files, nil
	}
	var files []string
	err := filepath.WalkDir(root, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			if d.Name() == ".git" {
				return filepath.SkipDir
			}
			return nil
		}
		rel, err := filepath.Rel(root, path)
		if err != nil {
			return err
		}
		files = append(files, filepath.ToSlash(rel))
		return nil
	})
	if err != nil {
		return nil, err
	}
	sort.Strings(files)
	return files, nil
}
