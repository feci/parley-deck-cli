package driver

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
)

// TestGitTreeCleanSetsOptionalLocksOff (consensus D8): every git probe the
// driver spawns must run with GIT_OPTIONAL_LOCKS=0 so read-only status checks
// never write .git on a weakly-coherent mount. A PATH-shimmed fake git records
// the env it saw.
func TestGitTreeCleanSetsOptionalLocksOff(t *testing.T) {
	dir := t.TempDir()
	record := filepath.Join(dir, "seen-env")
	// §D.9 test-binary re-exec port of the PATH-shim shell script (the
	// historical t.Skip on Windows is retired with it): the shim is a copy
	// of this test binary named "git"; TestMain below dispatches the
	// "record-env" role (append GIT_OPTIONAL_LOCKS to the spec's path).
	shim := writeGitShimFixture(t, dir, record)
	_ = shim
	t.Setenv("PATH", dir+string(os.PathListSeparator)+os.Getenv("PATH"))

	_ = gitTreeClean(dir)

	data, err := os.ReadFile(record)
	if err != nil {
		t.Fatalf("the shim never ran: %v", err)
	}
	for i, line := range splitNonEmptyLines(string(data)) {
		if line != "0" {
			t.Fatalf("probe %d ran with GIT_OPTIONAL_LOCKS=%q, want 0", i, line)
		}
	}
}

func splitNonEmptyLines(s string) []string {
	var out []string
	start := 0
	for i := 0; i <= len(s); i++ {
		if i == len(s) || s[i] == '\n' {
			if i > start {
				out = append(out, s[start:i])
			}
			start = i + 1
		}
	}
	return out
}

// writeGitShimFixture installs the re-exec git shim: a copy of this test
// binary named git (.exe on Windows; LookPath resolves PATHEXT) with a
// sibling git.role spec naming the record file.
func writeGitShimFixture(t *testing.T, dir, record string) string {
	t.Helper()
	src, err := os.Executable()
	if err != nil {
		t.Fatal(err)
	}
	name := "git"
	if runtime.GOOS == "windows" {
		name += ".exe"
	}
	installed := filepath.Join(dir, name)
	in, err := os.Open(src)
	if err != nil {
		t.Fatal(err)
	}
	defer in.Close()
	out, err := os.OpenFile(installed, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0o755)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := io.Copy(out, in); err != nil {
		out.Close()
		t.Fatal(err)
	}
	if err := out.Close(); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(dir, "git.role"), []byte("record-env "+record), 0o644); err != nil {
		t.Fatal(err)
	}
	return installed
}

// TestMain dispatches the re-exec git shim role before the testing framework.
// The role is resolved from os.Executable() — the actual shim location —
// because argv[0] may be a bare name resolved through PATH.
//
// RECURSION GUARD (runaway-fixture postmortem, 2026-09-28): a renamed copy
// that cannot resolve its role exits loudly — falling through to m.Run()
// would run the suite inside the shim (unbounded recursion).
func TestMain(m *testing.M) {
	if exe, err := os.Executable(); err == nil {
		exeBase := filepath.Base(exe)
		if strings.TrimSuffix(exeBase, ".exe") == "git" {
			if spec, rerr := os.ReadFile(filepath.Join(filepath.Dir(exe), "git.role")); rerr == nil {
				fields := strings.Fields(string(spec))
				if len(fields) == 2 && fields[0] == "record-env" {
					f, ferr := os.OpenFile(fields[1], os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o644)
					if ferr == nil {
						fmt.Fprintf(f, "%s\n", os.Getenv("GIT_OPTIONAL_LOCKS"))
						f.Close()
					}
					os.Exit(0)
				}
			}
			fmt.Fprintln(os.Stderr, "git shim copy could not resolve its role spec; refusing to enter the test framework")
			os.Exit(70)
		}
	}
	os.Exit(m.Run())
}
