package agents

import "testing"

func TestReviewGateBufferedDefaultContracts(t *testing.T) {
	want := map[string]bool{"codex": false, "claude": true, "kimi": false, "zcode": true, "agy": true}
	for _, s := range DefaultSpecs() {
		if buffered, ok := want[s.ID]; ok {
			if s.BuffersStdout != buffered {
				t.Errorf("%s buffering=%v want %v", s.ID, s.BuffersStdout, buffered)
			}
			delete(want, s.ID)
		}
	}
	if len(want) != 0 {
		t.Fatal("missing adapters", want)
	}
}
