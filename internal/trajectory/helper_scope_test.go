package trajectory

import (
	"context"
	"strings"
	"testing"
)

func TestHelperScopePinsOriginalRequestAndMembership(t *testing.T) {
	root, _, _, request := quorumPatchFixture(t, "[builder, reviewer]")
	root, err := canonicalRoot(root)
	if err != nil {
		t.Fatal(err)
	}
	original := HelperRequest{Version: 1, Ticket: VerificationTicket{1, root, "scope-check", request}, Participants: []string{"builder", "reviewer"}, Criteria: []Criterion{{"material", "true"}}}
	if err := CheckHelperScope(context.Background(), original); err != nil {
		t.Fatal("original scope refused", err)
	}
	for _, kind := range []string{"widened", "reordered", "duplicate", "invocation", "criterion", "version"} {
		t.Run(kind, func(t *testing.T) {
			req := original
			req.Criteria = append([]Criterion(nil), original.Criteria...)
			switch kind {
			case "widened":
				req.Participants = []string{"builder", "reviewer", "other"}
			case "reordered":
				req.Participants = []string{"reviewer", "builder"}
			case "duplicate":
				req.Participants = []string{"builder", "reviewer", "reviewer"}
			case "invocation":
				req.Ticket.Request.InvocationID = "different-original"
			case "criterion":
				req.Criteria[0].Command = "false"
			case "version":
				req.Version = 2
			}
			if err := CheckHelperScope(context.Background(), req); err == nil {
				t.Fatal("changed helper scope accepted", kind)
			}
		})
	}
	missing := original
	missing.Ticket.Root, err = canonicalRoot(t.TempDir())
	if err != nil {
		t.Fatal(err)
	}
	if err := CheckHelperScope(context.Background(), missing); err == nil || !strings.Contains(err.Error(), "authority disappeared") {
		t.Fatal("missing trajectory appeared authorized", err)
	}
}
