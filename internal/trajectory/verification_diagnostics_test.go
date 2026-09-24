package trajectory

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// dumpRetainedVerificationSteps prints the retained captured-verification
// journal artifacts (step/process/receipt envelopes) for hosted diagnosis of
// "captured criterion execution is incomplete" failures. Test-only, called
// from failure paths only: these envelopes are the product's own retained
// records — Diagnostics was scrubbed (credential-shaped tokens redacted) and
// bounded (last 100 lines / 4 KiB) before it was ever persisted — so printing
// them leaks nothing the product has not already made safe to retain.
func dumpRetainedVerificationSteps(t *testing.T, journal string) {
	t.Helper()
	entries, err := os.ReadDir(journal)
	if err != nil {
		t.Logf("diagnostics: verification journal %s unavailable: %v", journal, err)
		return
	}
	for _, entry := range entries {
		name := entry.Name()
		if entry.IsDir() || (name != "receipt.json" && !strings.HasPrefix(name, "step-") && !strings.HasPrefix(name, "process-")) {
			continue
		}
		raw, err := os.ReadFile(filepath.Join(journal, name))
		if err != nil {
			t.Logf("diagnostics: %s unreadable: %v", name, err)
			continue
		}
		t.Logf("diagnostics: retained %s: %s", name, strings.TrimSpace(string(raw)))
	}
}

// dumpAllRetainedVerifications sweeps every captured-verification journal
// under a runtime root (fixtures create exactly one) and dumps each.
func dumpAllRetainedVerifications(t *testing.T, root string) {
	t.Helper()
	journals, err := filepath.Glob(filepath.Join(root, ".parley-runtime", "trajectory-verifications", "*"))
	if err != nil || len(journals) == 0 {
		t.Logf("diagnostics: no retained verification journals under %s (glob err: %v)", root, err)
		return
	}
	for _, journal := range journals {
		dumpRetainedVerificationSteps(t, journal)
	}
}
