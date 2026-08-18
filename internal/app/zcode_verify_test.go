package app

import (
	"bytes"
	"context"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"parley-deck-cli/internal/agents"
)

// zcodeSpec returns the shipped built-in spec, so these tests exercise the real argv
// rather than a copy that could drift from it.
func zcodeSpec(t *testing.T) agents.Spec {
	t.Helper()
	for _, s := range agents.DefaultSpecs() {
		if s.ID == "zcode" {
			return s
		}
	}
	t.Fatal("built-in spec zcode is missing")
	return agents.Spec{}
}

// writeFakeZcode installs a stub `zcode` whose task behaviour is chosen by mode:
//
//	"help-exit0" — prints top-level help and exits 0 without writing anything. This is the
//	               shape a rejected-flag or degraded launch would take, and full verification
//	               MUST NOT accept it.
//	"honest"     — accepts exactly `--prompt <text> --mode yolo --cwd <root>`, extracts the
//	               output path and sentinel from the probe prompt, and writes the file.
func writeFakeZcode(t *testing.T, dir, mode string) string {
	t.Helper()
	path := filepath.Join(dir, "zcode")
	body := `#!/bin/sh
if [ "$1" = "--version" ]; then echo 'zcode-app-cli 0.0.0-test'; exit 0; fi

# Assert the exact token order the built-in spec promises. A spec that appends or reorders
# an option zcode does not accept must fail here rather than in production.
[ "$1" = "--prompt" ] || { echo "arg1=$1 want --prompt" >&2; exit 64; }
PROMPT="$2"
[ "$3" = "--mode" ] || { echo "arg3=$3 want --mode" >&2; exit 64; }
[ "$4" = "yolo" ]   || { echo "arg4=$4 want yolo" >&2; exit 64; }
[ "$5" = "--cwd" ]  || { echo "arg5=$5 want --cwd" >&2; exit 64; }
[ -n "$6" ]         || { echo "arg6 (root) empty" >&2; exit 64; }
[ "$#" -eq 6 ]      || { echo "argc=$# want 6" >&2; exit 64; }

if [ "MODE_PLACEHOLDER" = "help-exit0" ]; then
  echo "Usage: zcode [options]"
  exit 0
fi

OUT=$(printf '%s\n' "$PROMPT" | sed -n '2p')
SENTINEL=$(printf '%s\n' "$PROMPT" | sed -n '5p')
printf '%s\nheadless probe ok\n' "$SENTINEL" > "$OUT"
exit 0
`
	body = strings.Replace(body, "MODE_PLACEHOLDER", mode, 1)
	if err := os.WriteFile(path, []byte(body), 0o755); err != nil {
		t.Fatal(err)
	}
	return path
}

func verifyZcode(t *testing.T, mode string) error {
	t.Helper()
	root := t.TempDir()
	bin := t.TempDir()
	spec := zcodeSpec(t)
	d := agents.Discovery{Spec: spec, Path: writeFakeZcode(t, bin, mode), Found: true, Version: "test"}
	var out bytes.Buffer
	return runFullVerification(context.Background(), root, []agents.Discovery{d}, &out)
}

// A launch that exits 0 while writing nothing must FAIL full verification. zcode's own
// failure modes have this shape, so exit status alone can never be the acceptance signal.
func TestZcodeFullVerifyRejectsHelpExitZero(t *testing.T) {
	if err := verifyZcode(t, "help-exit0"); err == nil {
		t.Fatal("full verification accepted a zcode that exited 0 and wrote no artifact")
	}
}

// The honest stub, driven by the real built-in argv, must pass.
func TestZcodeFullVerifyAcceptsHonestLaunch(t *testing.T) {
	if err := verifyZcode(t, "honest"); err != nil {
		t.Fatalf("full verification rejected a correct zcode launch: %v", err)
	}
}

// The `--explain` trailer must name the model SOURCE without reading its current value:
// a live read would answer "which model will deliberate?" with a snapshot that can change
// before launch, one command away from the column that refuses exactly that staleness.
func TestZcodeExplainTrailerIsStatic(t *testing.T) {
	trailer := modelSourceTrailer("zcode")
	if trailer == "" {
		t.Fatal("zcode has no model-source trailer; an operator reading `model unknown` has no way to learn where the model comes from")
	}
	for _, want := range []string{"~/.zcode/cli/config.json", "model.main", "never passed by parley"} {
		if !strings.Contains(trailer, want) {
			t.Errorf("trailer missing %q: %s", want, trailer)
		}
	}
	// Must not carry a concrete model id — that would be the live value the design rejects.
	for _, forbidden := range []string{"glm-5.3", "zai/glm", "glm-5-turbo"} {
		if strings.Contains(trailer, forbidden) {
			t.Errorf("trailer carries live model value %q; it must name the source only", forbidden)
		}
	}
	if modelSourceTrailer("claude") != "" {
		t.Error("adapters that can bind a model must not get a source trailer")
	}
}
