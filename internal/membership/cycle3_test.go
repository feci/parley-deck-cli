package membership

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"parley-deck-cli/internal/protocol"
	"parley-deck-cli/internal/quota"
	"parley-deck-cli/internal/store"
)

func TestQuotaCycle3NoticeDeliveryRecovery(t *testing.T) {
	for _, mode := range []string{"applied-archive", "applied-delete", "before-receipt-archive", "before-receipt-delete", "before-notice", "corrupt-receipt-archive", "corrupt-receipt-delete"} {
		t.Run(mode, func(t *testing.T) {
			root, dir, run := fixture(t, quota.NewPolicy(nil, nil))
			ctx := leaseFixture(t, dir, run)
			old := projectionFault
			t.Cleanup(func() { projectionFault = old })
			interrupt := strings.HasPrefix(mode, "before-")
			if interrupt {
				projectionFault = func(stage string) error {
					if stage == "applied" || (mode == "before-notice" && stage == "notice") {
						return fmt.Errorf("checked interruption %s", stage)
					}
					return nil
				}
			}
			b, err := Settle(ctx, root, dir, run, "round-01", []string{"a", "b", "c", "d"}, candidates())
			if b == nil || (err != nil) != interrupt {
				t.Fatal(b, err)
			}
			projectionFault = old
			notice := filepath.Join(root, protocol.DeckDir, "inbox", "parley-to-user_"+b.ID+".md")
			receipt := filepath.Join(dir, "quota-applied", b.ID)
			if strings.HasPrefix(mode, "corrupt-receipt") {
				if err := os.WriteFile(receipt, []byte("other-transition\n"), 0600); err != nil {
					t.Fatal(err)
				}
			}
			if strings.HasSuffix(mode, "archive") {
				archive := filepath.Join(filepath.Dir(notice), "archived")
				if err := os.MkdirAll(archive, 0700); err != nil {
					t.Fatal(err)
				}
				if err := os.Rename(notice, filepath.Join(archive, filepath.Base(notice))); err != nil {
					t.Fatal(err)
				}
			} else if strings.HasSuffix(mode, "delete") {
				if err := os.Remove(notice); err != nil {
					t.Fatal(err)
				}
			}
			for i := 0; i < 3; i++ {
				if _, err := Before(ctx, root, dir, run); err != nil {
					t.Fatal(err)
				}
			}
			_, err = os.Stat(notice)
			wantLive := mode == "before-receipt-delete" || mode == "before-notice" || mode == "corrupt-receipt-delete"
			if (err == nil) != wantLive {
				t.Fatal("wrong publication after inbox action", mode, err)
			}
			if ok, err := quota.ReadApplied(dir, b.ID); err != nil || !ok {
				t.Fatal(ok, err)
			}
			events, err := store.New(filepath.Join(root, protocol.DeckDir, "runs", run)).Load()
			if err != nil {
				t.Fatal(err)
			}
			n := 0
			for _, e := range events {
				if e.Type == "round.completed" && e.Data["quota_transition"] == b.ID {
					n++
				}
			}
			if n != 1 {
				t.Fatal("terminal evaluation count", n)
			}
		})
	}
}

func TestQuotaCycle3HistoricalManualClarificationSurvivesInboxActions(t *testing.T) {
	for _, action := range []string{"archive", "delete"} {
		t.Run(action, func(t *testing.T) {
			root, dir, run := fixture(t, quota.Policy{Enabled: false, Scope: quota.KickoffAndMidIdea})
			prompt := filepath.Join(dir, "00-prompt.md")
			raw, err := os.ReadFile(prompt)
			if err != nil {
				t.Fatal(err)
			}
			raw = []byte(strings.Replace(string(raw), "participants: [a, b, c, d]", "participants: [a, b, c]\nexcluded: d — unavailable — confirmed 2026-10-04", 1))
			if err := os.WriteFile(prompt, raw, 0600); err != nil {
				t.Fatal(err)
			}
			if err := RecordManual(root, dir); err != nil {
				t.Fatal(err)
			}
			h, err := quota.ReadHistory(dir)
			if err != nil {
				t.Fatal(err)
			}
			b := h.Batches[0]
			inbox := filepath.Join(root, protocol.DeckDir, "inbox")
			name := "parley-to-user_" + b.ID + ".md"
			// Exact historical producer notice, never an owner answer or catch-up proof.
			if err := os.WriteFile(filepath.Join(inbox, name), []byte(b.LegacyManualNotice()), 0600); err != nil {
				t.Fatal(err)
			}
			archive := filepath.Join(inbox, "archived")
			if err := os.MkdirAll(archive, 0700); err != nil {
				t.Fatal(err)
			}
			if err := os.Rename(filepath.Join(inbox, name), filepath.Join(archive, name)); err != nil {
				t.Fatal(err)
			}
			if _, err := Before(context.Background(), root, dir, run); err != nil {
				t.Fatal(err)
			}
			correction := "parley-to-user_" + b.ID + "-manual-authority.md"
			if action == "archive" {
				err = os.Rename(filepath.Join(inbox, correction), filepath.Join(archive, correction))
			} else {
				err = os.Remove(filepath.Join(inbox, correction))
			}
			if err != nil {
				t.Fatal(err)
			}
			for i := 0; i < 2; i++ {
				if _, err := Before(context.Background(), root, dir, run); err != nil {
					t.Fatal(err)
				}
			}
			if _, err := os.Stat(filepath.Join(inbox, name)); !os.IsNotExist(err) {
				t.Fatal("historical notice republished")
			}
			if _, err := os.Stat(filepath.Join(inbox, correction)); !os.IsNotExist(err) {
				t.Fatal("clarification republished")
			}
			h, err = quota.ReadHistory(dir)
			if err != nil || h.Batches[0].Owner.Authority != nil {
				t.Fatal("manual display gained authority", h, err)
			}
			if action == "archive" {
				if err := os.WriteFile(filepath.Join(archive, correction), []byte("false authority"), 0600); err != nil {
					t.Fatal(err)
				}
				if _, err := Before(context.Background(), root, dir, run); err != nil {
					t.Fatal("owner clarification edit gated replay", err)
				}
			}
		})
	}
}
