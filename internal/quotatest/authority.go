// Package quotatest builds synthetic committed owner-answer fixtures in isolated
// test repositories. It never commits or changes the real project repository.
package quotatest

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"parley-deck-cli/internal/quota"
)

func Git(t testing.TB, root string, args ...string) string {
	t.Helper()
	cmd := exec.Command("git", append([]string{"-C", root, "-c", "user.name=Quota test fixture", "-c", "user.email=quota-test@example.invalid", "-c", "commit.gpgsign=false", "-c", "core.hooksPath=/dev/null"}, args...)...)
	raw, e := cmd.CombinedOutput()
	if e != nil {
		t.Fatalf("fixture git %v: %v %s", args, e, raw)
	}
	return strings.TrimSpace(string(raw))
}
func Authority(t testing.TB, root, idea, name, quote string) quota.Authority {
	t.Helper()
	if _, e := os.Stat(filepath.Join(root, ".git")); os.IsNotExist(e) {
		Git(t, root, "init", "-q")
	}
	path := "parley-deck/inbox/user-to-all_" + name + ".md"
	if e := os.MkdirAll(filepath.Dir(filepath.Join(root, path)), 0700); e != nil {
		t.Fatal(e)
	}
	raw := fmt.Sprintf("---\nfrom: user\nto: all\nidea: %s\n---\n\nSynthetic test owner answer (not real authorization).\n%s\n", idea, quote)
	if e := os.WriteFile(filepath.Join(root, path), []byte(raw), 0600); e != nil {
		t.Fatal(e)
	}
	Git(t, root, "add", "--", path)
	Git(t, root, "commit", "-qm", "synthetic authority fixture", "--", path)
	a, e := quota.BindAuthority(root, idea, path, Git(t, root, "rev-parse", "HEAD"), quote)
	if e != nil {
		t.Fatal(e)
	}
	return a
}
func Fields(a quota.Authority) string {
	return fmt.Sprintf("Owner-answer: %s\nOwner-commit: %s\nOwner-blob: %s\nOwner-sha256: %s\nOwner-quote: %s\n", a.Path, a.Commit, a.Blob, a.SHA256, a.Quote)
}
