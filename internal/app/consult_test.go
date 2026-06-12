package app

import "testing"

func TestConsultSlug(t *testing.T) {
	cases := map[string]string{
		"How does the TUI ribbon work?":      "how-does-the-tui-ribbon-work",
		"  Weird?? punctuation!! everywhere": "weird-punctuation-everywhere",
		"":                                   "question",
	}
	for in, want := range cases {
		if got := consultSlug(in); got != want {
			t.Errorf("consultSlug(%q)=%q, want %q", in, got, want)
		}
	}
}
