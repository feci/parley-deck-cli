package runner

// lean-organizer D.4: the kimi stream-json envelope unwrapper — text consumers read
// ASSISTANT CONTENT, not envelopes; non-JSON passthrough; never a guess.

import (
	"testing"
)

func TestUnwrapKimiStreamJSON(t *testing.T) {
	in := "{\"role\":\"meta\",\"type\":\"system.version\",\"version\":\"0.42.0\"}\n" +
		"{\"role\":\"assistant\",\"content\":\"PONG\"}\n" +
		"{\"role\":\"meta\",\"type\":\"session.resume_hint\",\"session_id\":\"s\",\"command\":\"kimi -r s\",\"content\":\"To resume\"}\n"
	if got := UnwrapKimiStreamJSON(in); got != "PONG" {
		t.Errorf("unwrap must yield the assistant content, got %q", got)
	}
	multi := "{\"role\":\"assistant\",\"content\":\"a\"}\n{\"role\":\"assistant\",\"content\":\"b\"}\n"
	if got := UnwrapKimiStreamJSON(multi); got != "a\nb" {
		t.Errorf("multi-line contents join with newline, got %q", got)
	}
	plain := "just text\nmore text\n"
	if got := UnwrapKimiStreamJSON(plain); got != plain {
		t.Errorf("non-JSON input passes through unchanged, got %q", got)
	}
	metaOnly := "{\"role\":\"meta\",\"type\":\"system.version\"}\n"
	if got := UnwrapKimiStreamJSON(metaOnly); got != metaOnly {
		t.Errorf("meta-only input is returned unchanged (never a guess), got %q", got)
	}
}
