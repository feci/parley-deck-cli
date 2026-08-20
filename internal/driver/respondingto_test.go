package driver

import (
	"os"
	"path/filepath"
	"testing"
)

// codex-1/F16: the cross-review gate accepted the mere presence of `responding-to`, so an empty
// list satisfied a gate whose whole purpose is proving the reviewer read its peers.
func TestRespondingToMustNameSomebody(t *testing.T) {
	cases := []struct {
		name  string
		front string
		want  bool
	}{
		{"named peers", "responding-to: [a/round-01, b/round-01]", true},
		{"single peer", "responding-to: a/round-01", true},
		{"empty list", "responding-to: []", false},
		{"bare key", "responding-to:", false},
		{"whitespace only", "responding-to:    ", false},
		{"commas only", "responding-to: [ , ]", false},
		{"absent", "agent: a", false},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			path := filepath.Join(t.TempDir(), "r.md")
			body := "---\nagent: a\n" + tc.front + "\n---\n\n## Findings\n\nsomething\n"
			if err := os.WriteFile(path, []byte(body), 0o644); err != nil {
				t.Fatal(err)
			}
			if got := hasRespondingTo(path); got != tc.want {
				t.Fatalf("hasRespondingTo(%q) = %v, want %v", tc.front, got, tc.want)
			}
		})
	}
}
