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

// codex-1/F15: a misspelled track was indistinguishable from an absent one, and "absent" means
// "apply nothing" — so a typo silently switched off every standard-track cap.
func TestMisspelledTrackIsRefusedNotDefaulted(t *testing.T) {
	dir := t.TempDir()
	write := func(v string) {
		body := "---\nidea: demo\ntrack: " + v + "\n---\n\n## Prompt\n\nx\n"
		if err := os.WriteFile(filepath.Join(dir, "00-prompt.md"), []byte(body), 0o644); err != nil {
			t.Fatal(err)
		}
	}

	write("standrd")
	if _, _, err := ReadTrackStrict(dir); err == nil {
		t.Fatal("a misspelled track must be reported, not silently defaulted")
	}

	write("deliberation")
	tr, present, err := ReadTrackStrict(dir)
	if err != nil || !present || string(tr) != "deliberation" {
		t.Fatalf("valid track mis-read: %v %v %v", tr, present, err)
	}

	if err := os.WriteFile(filepath.Join(dir, "00-prompt.md"), []byte("---\nidea: demo\n---\n\n## Prompt\n\nx\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	if _, present, err := ReadTrackStrict(dir); err != nil || present {
		t.Fatalf("an absent track is not an error: present=%v err=%v", present, err)
	}
}
