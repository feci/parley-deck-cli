package runner

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"
	"testing"
	"time"
)

// §D.9 test-binary re-exec shell for the runner fixture family: on Windows
// /bin/sh does not exist, so a copy of this binary named "sh" emulates the
// exact -c scripts the family uses. Unix is untouched (/bin/sh). The
// emulation covers: exit N, exec sleep N (long-running), cat "$1" (print the
// first positional arg), printf-style child output, and the protocol-context
// fixture scripts.
var shellExeOnce = sync.OnceValue(func() string {
	exe, err := os.Executable()
	if err != nil {
		panic(err)
	}
	dir := filepath.Dir(exe)
	dst := filepath.Join(dir, "sh.exe")
	out, err := os.OpenFile(dst, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0o755)
	if err != nil {
		panic(err)
	}
	in, _ := os.Open(exe)
	if _, err := out.ReadFrom(in); err != nil {
		panic(err)
	}
	out.Close()
	in.Close()
	return dst
})

// TestMain dispatches the "sh" copy before the testing framework (the same
// recursion-guard discipline as the app package: a renamed copy that cannot
// resolve its role exits loudly, never m.Run()).
func TestMain(m *testing.M) {
	if exe, err := os.Executable(); err == nil && strings.TrimSuffix(filepath.Base(exe), ".exe") == "sh" {
		os.Exit(runShellRole(os.Args[1:]))
	}
	os.Exit(m.Run())
}

func runShellRole(argv []string) int {
	// argv: -c <script> <name> [extra...]
	if len(argv) < 2 || argv[0] != "-c" {
		fmt.Fprintln(os.Stderr, "sh role: want -c <script>")
		return 70
	}
	script := argv[1]
	name := "fixture"
	args := argv[2:]
	if len(args) > 0 {
		name = args[0]
		args = args[1:]
	}
	switch {
	case script == "exit 7":
		return 7
	case strings.HasPrefix(script, "exec sleep "):
		secs := strings.TrimSpace(script[len("exec sleep "):])
		d, err := time.ParseDuration(secs + "s")
		if err != nil {
			return 70
		}
		time.Sleep(d)
		return 0
	case script == fmt.Sprintf("cat %q; printf '\\nchild output\\n'", name) || script == `cat "$1"; printf '\nchild output\n'`:
		if len(args) > 0 {
			data, err := os.ReadFile(args[0])
			if err == nil {
				os.Stdout.Write(data)
			}
		}
		fmt.Print("\nchild output\n")
		return 0
	case script == "cat >/dev/null":
		return 0
	}
	fmt.Fprintf(os.Stderr, "sh role: unhandled script %q\n", script)
	return 70
}

var _ = exec.Command // keep exec referenced for future fixtures
