package evidence

import (
	"net"
	"os"
	"os/exec"
	"path/filepath"
	"syscall"
	"testing"
)

// scratchGitRepo returns a fresh git repository in a temp dir with the given
// files committed, so the digest exercises the production `git ls-files` path
// (a bare temp dir inside this worktree would be swallowed by the parent
// repo's .gitignore and is not a valid stand-in).
func scratchGitRepo(t *testing.T, files map[string]string) string {
	t.Helper()
	root := t.TempDir()
	git(t, root, "init", "-q")
	for rel, content := range files {
		abs := filepath.Join(root, filepath.FromSlash(rel))
		if err := os.MkdirAll(filepath.Dir(abs), 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(abs, []byte(content), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	git(t, root, "add", "-A")
	git(t, root, "-c", "user.email=t@t", "-c", "user.name=t", "commit", "-qm", "init")
	return root
}

func git(t *testing.T, dir string, args ...string) {
	t.Helper()
	cmd := exec.Command("git", append([]string{"-C", dir}, args...)...)
	if out, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("git %v: %v\n%s", args, err, out)
	}
}

// Adversarial (probe 3): retargeting an in-root symlink a→b MUST change the
// digest even though no file content changed.
func TestTreeDigestSymlinkRetargetChanges(t *testing.T) {
	root := scratchGitRepo(t, map[string]string{"a.txt": "same\n", "b.txt": "same\n"})
	link := filepath.Join(root, "link")
	if err := os.Symlink("a.txt", link); err != nil {
		t.Fatal(err)
	}
	d1, err := TreeDigest(root)
	if err != nil {
		t.Fatal(err)
	}
	if err := os.Remove(link); err != nil {
		t.Fatal(err)
	}
	if err := os.Symlink("b.txt", link); err != nil {
		t.Fatal(err)
	}
	d2, err := TreeDigest(root)
	if err != nil {
		t.Fatal(err)
	}
	if d1 == d2 {
		t.Fatal("symlink retarget a→b must change the tree digest")
	}
}

// Adversarial (probe 3): a mode-only change (0600→0700) MUST change the digest.
func TestTreeDigestModeChangeChanges(t *testing.T) {
	root := scratchGitRepo(t, map[string]string{"x.sh": "#!/bin/sh\n"})
	f := filepath.Join(root, "x.sh")
	if err := os.Chmod(f, 0o600); err != nil {
		t.Fatal(err)
	}
	d1, err := TreeDigest(root)
	if err != nil {
		t.Fatal(err)
	}
	if err := os.Chmod(f, 0o700); err != nil {
		t.Fatal(err)
	}
	d2, err := TreeDigest(root)
	if err != nil {
		t.Fatal(err)
	}
	if d1 == d2 {
		t.Fatal("mode change 0600→0700 must change the tree digest")
	}
}

// Adversarial: an unsupported entry (fifo, or a unix socket where fifos are
// not supported) is an error, never silently skipped — including when the
// entry is UNTRACKED and git's file inventory omits it entirely.
func TestTreeDigestUnsupportedEntryFails(t *testing.T) {
	root := scratchGitRepo(t, map[string]string{"a.go": "package a\n"})
	if err := syscall.Mkfifo(filepath.Join(root, "pipe"), 0o644); err != nil {
		// This host's shared test volume cannot host FIFOs, and its long
		// mandated TMPDIR paths exceed the AF_UNIX sun_path limit. Fall back
		// to a unix socket created through a short RELATIVE path; the socket
		// file lands inside the repo and must fail the digest the same way.
		makeUnsupportedSocket(t, root)
	}
	if _, err := TreeDigest(root); err == nil {
		t.Fatal("unsupported entry type must fail the digest")
	}
}

// makeUnsupportedSocket creates a unix socket file inside root via a relative
// bind path (cwd-relative binds sidestep the AF_UNIX length limit). It stays
// open until cleanup so the socket file is present during the digest.
func makeUnsupportedSocket(t *testing.T, root string) {
	t.Helper()
	prev, err := os.Getwd()
	if err != nil {
		t.Skipf("cannot set up unsupported-entry probe: %v", err)
	}
	if err := os.Chdir(root); err != nil {
		t.Skipf("cannot set up unsupported-entry probe: %v", err)
	}
	ln, err := net.Listen("unix", "parley-evidence-probe.sock")
	if cherr := os.Chdir(prev); cherr != nil {
		t.Fatalf("cannot restore cwd: %v", cherr)
	}
	if err != nil {
		t.Skipf("no unsupported-entry probe available on this filesystem: %v", err)
	}
	t.Cleanup(func() { _ = ln.Close() })
}

// Positive: a symlink through a /var-style path alias inside the canonical
// root is accepted (root is canonicalized before the escape check).
func TestTreeDigestCanonicalRootSymlink(t *testing.T) {
	root := scratchGitRepo(t, map[string]string{"a.txt": "x\n"})
	if err := os.Symlink("a.txt", filepath.Join(root, "ok-link")); err != nil {
		t.Fatal(err)
	}
	if _, err := TreeDigest(root); err != nil {
		t.Fatalf("in-root symlink must digest cleanly: %v", err)
	}
}

// Positive: digest is deterministic for unchanged content.
func TestTreeDigestDeterministic(t *testing.T) {
	root := scratchGitRepo(t, map[string]string{
		"a.go":     "package a\n",
		"sub/b.go": "package sub\n",
	})

	d1, err := TreeDigest(root)
	if err != nil {
		t.Fatal(err)
	}
	d2, err := TreeDigest(root)
	if err != nil {
		t.Fatal(err)
	}
	if d1 != d2 || d1 == "" {
		t.Fatalf("digest not deterministic: %q vs %q", d1, d2)
	}
	if ReviewedCommit(root) == "" {
		t.Fatal("committed scratch repo must report a reviewed commit")
	}
	if TreeDirty(root) {
		t.Fatal("clean committed tree must not be dirty")
	}
}

// Positive: the non-git fallback walk covers commitless directories.
func TestTreeDigestCommitlessFallback(t *testing.T) {
	root := t.TempDir()
	git(t, root, "init", "-q") // a repo with no commits: no HEAD → fallback walk
	os.WriteFile(filepath.Join(root, "a.go"), []byte("package a\n"), 0o644)
	if ReviewedCommit(root) != "" {
		t.Fatal("commitless repo must not report a reviewed commit")
	}
	if !TreeDirty(root) {
		t.Fatal("no commit baseline must report dirty")
	}
	if _, err := TreeDigest(root); err != nil {
		t.Fatalf("fallback walk must digest the tree: %v", err)
	}
}

// Adversarial: any code change invalidates the tested-tree digest (staleness
// is detectable at close time).
func TestTreeDigestChangesWithContent(t *testing.T) {
	root := scratchGitRepo(t, map[string]string{"a.go": "package a\n"})
	d1, err := TreeDigest(root)
	if err != nil {
		t.Fatal(err)
	}
	os.WriteFile(filepath.Join(root, "a.go"), []byte("package a // changed\n"), 0o644)
	if !TreeDirty(root) {
		t.Fatal("modified tracked file must mark the tree dirty")
	}
	d2, err := TreeDigest(root)
	if err != nil {
		t.Fatal(err)
	}
	if d1 == d2 {
		t.Fatal("content change must change the digest")
	}
	// Adding a file also changes it.
	os.WriteFile(filepath.Join(root, "c.go"), []byte("package a\n"), 0o644)
	d3, err := TreeDigest(root)
	if err != nil {
		t.Fatal(err)
	}
	if d3 == d2 {
		t.Fatal("adding a file must change the digest")
	}
}

// Positive: only the explicitly excluded evidence artifacts are left out.
func TestTreeDigestExcludesOnlyDefinedArtifacts(t *testing.T) {
	root := scratchGitRepo(t, map[string]string{
		"code.go":       "package a\n",
		"EVIDENCE.json": "{}",
	})
	withEvidence, err := TreeDigest(root)
	if err != nil {
		t.Fatal(err)
	}
	excluded, err := TreeDigest(root, "EVIDENCE.json")
	if err != nil {
		t.Fatal(err)
	}
	if withEvidence == excluded {
		t.Fatal("exclusion of EVIDENCE.json must change the digest")
	}
	// And rewriting ONLY the excluded artifact leaves the digest unchanged —
	// persisting evidence must not invalidate itself.
	os.WriteFile(filepath.Join(root, "EVIDENCE.json"), []byte("{\"records\":[]}"), 0o644)
	excluded2, err := TreeDigest(root, "EVIDENCE.json")
	if err != nil {
		t.Fatal(err)
	}
	if excluded != excluded2 {
		t.Fatal("changes to the excluded evidence artifact must not change the digest")
	}
}

// Adversarial: a symlink escaping the root is an error, never followed.
func TestTreeDigestRejectsEscapingSymlink(t *testing.T) {
	outside := t.TempDir()
	os.WriteFile(filepath.Join(outside, "secret.go"), []byte("package s\n"), 0o644)
	root := scratchGitRepo(t, map[string]string{"a.go": "package a\n"})
	if err := os.Symlink(filepath.Join(outside, "secret.go"), filepath.Join(root, "link.go")); err != nil {
		t.Skipf("symlinks unavailable: %v", err)
	}
	if _, err := TreeDigest(root); err == nil {
		t.Fatal("escaping symlink must fail the digest")
	}
}

// Adversarial: an empty tree fails rather than certifying nothing.
func TestTreeDigestEmptyTreeFails(t *testing.T) {
	root := scratchGitRepo(t, map[string]string{".gitkeep": ""})
	os.Remove(filepath.Join(root, ".gitkeep"))
	if _, err := TreeDigest(root); err == nil {
		t.Fatal("empty tree must fail")
	}
}

// Report persistence: roundtrip, and failure modes are errors.
func TestSaveLoadRoundtrip(t *testing.T) {
	dir := t.TempDir()
	r := positiveReport()
	if err := Save(dir, r); err != nil {
		t.Fatal(err)
	}
	got, err := Load(dir)
	if err != nil {
		t.Fatal(err)
	}
	if got.ReviewedCommit != r.ReviewedCommit || len(got.Records) != len(r.Records) {
		t.Fatalf("roundtrip mismatch: %+v", got)
	}
}

// Adversarial: missing and corrupt reports are errors (fail closed).
func TestLoadMissingOrCorruptFails(t *testing.T) {
	dir := t.TempDir()
	if _, err := Load(dir); err == nil {
		t.Fatal("missing report must fail to load")
	}
	os.WriteFile(ReportPath(dir), []byte("{not json"), 0o644)
	if _, err := Load(dir); err == nil {
		t.Fatal("corrupt report must fail to load")
	}
}

// Adversarial: a write failure surfaces as an error — callers must treat it
// as a failure of the whole completion attempt.
func TestSaveUnwritableDirFails(t *testing.T) {
	dir := t.TempDir()
	os.Chmod(dir, 0o555)
	t.Cleanup(func() { os.Chmod(dir, 0o755) })
	if err := Save(dir, positiveReport()); err == nil {
		t.Fatal("save into an unwritable dir must fail")
	}
}

// Adversarial: a failed save must not leave a readable half-written report.
func TestFailedSaveLeavesNoReport(t *testing.T) {
	dir := t.TempDir()
	os.Chmod(dir, 0o555)
	t.Cleanup(func() { os.Chmod(dir, 0o755) })
	_ = Save(dir, positiveReport())
	if _, err := os.Stat(ReportPath(dir)); !os.IsNotExist(err) {
		t.Fatalf("failed save left a report behind: %v", err)
	}
}
