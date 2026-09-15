package evidence

import (
	"bytes"
	"context"
	"errors"
	"os/exec"
	"slices"
	"strings"
	"time"
)

const MaxSourceInventoryEntries = 100000
const MaxSourceInventoryBytes = 16 << 20

// ErrLocalSourceExcludes means an untracked path's inclusion depends on local
// or global Git exclusions outside the source inventory. Do not read/capture
// those hidden files or silently certify their omission. The error deliberately
// contains no private path, ignore pattern or Git diagnostic.
var ErrLocalSourceExcludes = errors.New("source inventory depends on local or global Git exclusions; make source scope explicit in project .gitignore or remove the local-only exclusion")

// GitSourceInventory returns sorted tracked and untracked paths using Git's
// project .gitignore rules. The ordinary --exclude-standard inventory must agree:
// private info/exclude and core.excludesFile rules may not silently hide source.
// Project-ignored files remain outside scope, even if a build can consume them.
// This is an observation of a cooperative worktree, not an atomic filesystem or
// Git-configuration snapshot; callers that need stability must recheck it.
func GitSourceInventory(ctx context.Context, root string) ([]string, error) {
	ctx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()
	project, err := readGitSourceInventory(ctx, root, "--exclude-per-directory=.gitignore", MaxSourceInventoryEntries, MaxSourceInventoryBytes)
	if err != nil {
		return nil, err
	}
	standard, err := readGitSourceInventory(ctx, root, "--exclude-standard", MaxSourceInventoryEntries, MaxSourceInventoryBytes)
	if err != nil {
		return nil, err
	}
	if !slices.Equal(project, standard) {
		return nil, ErrLocalSourceExcludes
	}
	return project, nil
}

// Keep the buffer private: embedding it promotes ReadFrom, letting io.Copy
// bypass Write and its bound when os/exec drains stdout.
type sourceInventoryOutput struct {
	buf      bytes.Buffer
	limit    int
	overflow bool
	cancel   context.CancelFunc
}

func (w *sourceInventoryOutput) Write(p []byte) (int, error) {
	if len(p) > w.limit-w.buf.Len() {
		w.overflow = true
		w.cancel()
		return 0, errors.New("source inventory output exceeds bound")
	}
	return w.buf.Write(p)
}

func readGitSourceInventory(ctx context.Context, root, exclusion string, maxEntries, maxBytes int) ([]string, error) {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()
	cmd := exec.CommandContext(ctx, "git", "-C", root, "--no-pager", "--no-optional-locks", "ls-files", "-c", "-o", exclusion, "-z")
	out := sourceInventoryOutput{limit: maxBytes, cancel: cancel}
	cmd.Stdout = &out
	// No stderr copy: diagnostics can contain private paths or source snippets.
	cmd.WaitDelay = time.Second
	err := cmd.Run()
	if out.overflow {
		return nil, errors.New("source inventory output exceeds bound")
	}
	if err != nil {
		if ctx.Err() != nil {
			return nil, ctx.Err()
		}
		return nil, errors.New("Git source inventory failed")
	}
	if err := ctx.Err(); err != nil {
		return nil, err
	}
	if out.buf.Len() == 0 {
		return []string{}, nil
	}
	if out.buf.Bytes()[out.buf.Len()-1] != 0 {
		return nil, errors.New("unterminated source inventory")
	}
	names := strings.Split(string(out.buf.Bytes()[:out.buf.Len()-1]), "\x00")
	if len(names) > maxEntries {
		return nil, errors.New("source inventory entry count exceeds bound")
	}
	slices.Sort(names)
	for i, name := range names {
		if name == "" || (i > 0 && name == names[i-1]) {
			return nil, errors.New("empty or duplicate source inventory path")
		}
	}
	return names, nil
}
