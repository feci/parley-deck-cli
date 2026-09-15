package budget

import (
	"context"
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
)

func TestCycleHistoryRequiresKnownPhaseBeforeFirstBinding(t *testing.T) {
	for _, kind := range []Kind{Fixup, CrossReview} {
		for _, phase := range []string{"", "unspecified", "manual", "round-0", "round-two", "future-phase", "fixup", "round-02", "round-01", "review", "review-consensus", "implementation", "evidence-verification"} {
			t.Run(string(kind)+"/"+phase, func(t *testing.T) {
				root := t.TempDir()
				dir := filepath.Join(root, ".parley-runtime", "invocations", "old")
				if err := os.MkdirAll(dir, 0700); err != nil {
					t.Fatal(err)
				}
				data, err := json.Marshal(map[string]any{"metadata": map[string]string{"idea": "idea", "phase": phase}})
				if err != nil {
					t.Fatal(err)
				}
				if err := os.WriteFile(filepath.Join(dir, "requested.json"), data, 0600); err != nil {
					t.Fatal(err)
				}
				knownFree := phase == "round-01" || phase == "review" || phase == "review-consensus" || phase == "implementation" || phase == "evidence-verification"
				otherKind := kind == Fixup && phase == "round-02" || kind == CrossReview && phase == "fixup"
				_, err = EnsureCycleBinding(context.Background(), root, "idea", kind, 3, 0, "", "")
				if (err == nil) != (knownFree || otherKind) {
					t.Fatalf("phase %q allowed=%v: %v", phase, knownFree || otherKind, err)
				}
			})
		}
	}
}
