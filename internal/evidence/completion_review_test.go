package evidence

import (
	"strings"
	"testing"
)

func TestCompletionTransitionMatchedTamperingAndMalformedBinding(t *testing.T) {
	const path = "parley-deck/ideas/idea-x/IMPLEMENTATION.md"
	t.Run("short-binding-does-not-panic", func(t *testing.T) {
		defer func() {
			if r := recover(); r != nil {
				t.Errorf("malformed binding panicked: %v", r)
			}
		}()
		r := attestedTransitionReport(transitionDoc, path)
		r.ExtraDigests[path] = "bad"
		if err := AuthorizeCompletionTransition(r, path, []byte(transitionDoc), "codex-1"); err == nil {
			t.Fatal("malformed binding was accepted")
		}
	})
	t.Run("matching-after-digest-still-requires-complete", func(t *testing.T) {
		r := attestedTransitionReport(transitionDoc, path)
		if err := AuthorizeCompletionTransition(r, path, []byte(transitionDoc), "codex-1"); err != nil {
			t.Fatal(err)
		}
		current := []byte(strings.Replace(transitionDoc, "status: implemented", "status: abandoned", 1))
		r.CompletionTransition.AfterSHA256 = sha256Hex(current)
		if reasons := VerifyCompletionTransition(r, path, r.ExtraDigests[path], current, "codex-1"); len(reasons) == 0 {
			t.Fatal("a matched forged digest authorized a non-complete state")
		}
	})
	for name, line := range map[string]string{"quoted-key": "\"status\": abandoned", "spaced-key": "status : abandoned", "merge-key": "<<: {status: abandoned}"} {
		t.Run(name, func(t *testing.T) {
			doc := strings.Replace(transitionDoc, "status: implemented", "status: implemented\n"+line, 1)
			if _, _, err := TransitionStatusToComplete([]byte(doc)); err == nil {
				t.Fatal("ambiguous YAML status transition accepted")
			}
		})
	}
}
