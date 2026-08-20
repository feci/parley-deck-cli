package protocol

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

// codex-1/F19: `init` substituted only the transport, so every CLI-created deck carried
// "**Workspace:** `<workspace-name>`" and "**Created:** `<date> …`" literally in the protocol's
// own header — false provenance on line 3, while the command reported success.
func TestInitFillsWorkspaceAndCreatedDate(t *testing.T) {
	root := filepath.Join(t.TempDir(), "my-project")
	if err := os.MkdirAll(root, 0o755); err != nil {
		t.Fatal(err)
	}
	out := cooperationForInitAt("local-dir", root, time.Date(2026, 8, 20, 0, 0, 0, 0, time.UTC))

	for _, unwanted := range []string{"`<workspace-name>`", "`<date> — created by parley init`"} {
		if strings.Contains(out, unwanted) {
			t.Errorf("init left the placeholder %s in the header", unwanted)
		}
	}
	for _, want := range []string{"**Workspace:** `my-project`", "**Created:** `2026-08-20 — created by parley init`"} {
		if !strings.Contains(out, want) {
			t.Errorf("init did not write %q", want)
		}
	}
}

// The static bootstrap template must keep its placeholders: the drift guard pins it and the skill
// ships the same shape as a vendor-neutral reference.
func TestStaticBootstrapTemplateKeepsItsPlaceholders(t *testing.T) {
	out := defaultCooperationForInit()
	for _, want := range []string{"**Workspace:** `<workspace-name>`", "created by parley init"} {
		if !strings.Contains(out, want) {
			t.Errorf("static template lost %q", want)
		}
	}
}
