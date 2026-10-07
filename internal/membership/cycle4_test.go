package membership

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"parley-deck-cli/internal/protocol"
	"parley-deck-cli/internal/quota"
	"parley-deck-cli/internal/store"
)

func cycle4Read(t *testing.T, path string) []byte {
	t.Helper()
	raw, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	return raw
}
func cycle4Write(t *testing.T, path string, raw []byte) {
	t.Helper()
	if err := os.MkdirAll(filepath.Dir(path), 0700); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, raw, 0600); err != nil {
		t.Fatal(err)
	}
}
func cycle4TerminalCount(t *testing.T, root, run, id string) {
	t.Helper()
	events, err := store.New(filepath.Join(root, protocol.DeckDir, "runs", run)).Load()
	if err != nil {
		t.Fatal(err)
	}
	n := 0
	for _, e := range events {
		if e.Type == "round.completed" && e.Data["quota_transition"] == id {
			n++
		}
	}
	if n != 1 {
		t.Fatal("terminal evaluation count", n)
	}
}

func TestQuotaCycle4OwnerNoticeCopies(t *testing.T) {
	for _, receipt := range []string{"applied", "before-receipt"} {
		for _, copy := range []string{"live", "archived", "both"} {
			for _, edit := range []string{"annotation", "replacement"} {
				t.Run(receipt+"/"+copy+"/"+edit, func(t *testing.T) {
					root, dir, run := fixture(t, quota.NewPolicy(nil, nil))
					ctx := leaseFixture(t, dir, run)
					old := projectionFault
					t.Cleanup(func() { projectionFault = old })
					if receipt == "before-receipt" {
						projectionFault = func(stage string) error {
							if stage == "applied" {
								return fmt.Errorf("interrupt before receipt")
							}
							return nil
						}
					}
					b, err := Settle(ctx, root, dir, run, "round-01", []string{"a", "b", "c", "d"}, candidates())
					if b == nil || (err != nil) != (receipt == "before-receipt") {
						t.Fatal(b, err)
					}
					projectionFault = old
					history := filepath.Join(dir, quota.HistoryDir, "000001.json")
					before := cycle4Read(t, history)
					notice := filepath.Join(root, protocol.DeckDir, "inbox", "parley-to-user_"+b.ID+".md")
					paths := []string{notice}
					if copy != "live" {
						archive := filepath.Join(filepath.Dir(notice), "archived", filepath.Base(notice))
						cycle4Write(t, archive, cycle4Read(t, notice))
						paths = append(paths, archive)
						if copy == "archived" {
							if err := os.Remove(notice); err != nil {
								t.Fatal(err)
							}
							paths = paths[1:]
						}
					}
					want := map[string][]byte{}
					for _, path := range paths {
						raw := []byte("Owner's arbitrary replacement\r\n")
						if edit == "annotation" {
							raw = append(cycle4Read(t, path), []byte("\n## Owner answer\nKeep the current quorum.\n")...)
						}
						cycle4Write(t, path, raw)
						want[path] = raw
					}
					if receipt == "before-receipt" {
						if _, _, err := protocol.QuotaMembers(dir, nil); err == nil {
							t.Fatal("missing receipt did not gate")
						}
					}
					for i := 0; i < 3; i++ {
						if _, err := Before(ctx, root, dir, run); err != nil {
							t.Fatal(err)
						}
						if _, _, err := protocol.QuotaMembers(dir, nil); err != nil {
							t.Fatal(err)
						}
					}
					for path, raw := range want {
						if !bytes.Equal(raw, cycle4Read(t, path)) {
							t.Fatal("owner bytes overwritten", path)
						}
					}
					if !bytes.Equal(before, cycle4Read(t, history)) {
						t.Fatal("immutable history rewritten")
					}
					if copy == "archived" {
						if _, err := os.Lstat(notice); !os.IsNotExist(err) {
							t.Fatal("archive republished", err)
						}
					}
					cycle4TerminalCount(t, root, run, b.ID)
				})
			}
		}
	}
}

func TestQuotaCycle4UnsafePublicationIsDiagnostic(t *testing.T) {
	for _, mode := range []string{"symlink", "directory", "archived-symlink", "inbox-symlink", "archive-dir-symlink", "inbox-file", "write-failure"} {
		t.Run(mode, func(t *testing.T) {
			root, dir, run := fixture(t, quota.NewPolicy(nil, nil))
			ctx := leaseFixture(t, dir, run)
			h, err := quota.ReadHistory(dir)
			if err != nil {
				t.Fatal(err)
			}
			d := quota.Evaluate(h.Policy(), h.Current, candidates(), quota.Roles{})
			b := quota.NewBatch(h, run, "round-01", h.Current, d, time.Now())
			if err := quota.CommitBatch(dir, b); err != nil {
				t.Fatal(err)
			}
			inbox := filepath.Join(root, protocol.DeckDir, "inbox")
			notice := filepath.Join(inbox, "parley-to-user_"+b.ID+".md")
			outside := t.TempDir()
			target := filepath.Join(outside, "owner.txt")
			cycle4Write(t, target, []byte("owner bytes"))
			if err := os.MkdirAll(inbox, 0700); err != nil {
				t.Fatal(err)
			}
			switch mode {
			case "symlink":
				err = os.Symlink(target, notice)
			case "directory":
				err = os.Mkdir(notice, 0700)
			case "archived-symlink":
				archive := filepath.Join(inbox, "archived")
				if err := os.Mkdir(archive, 0700); err != nil {
					t.Fatal(err)
				}
				err = os.Symlink(target, filepath.Join(archive, filepath.Base(notice)))
			case "inbox-symlink", "inbox-file":
				if err := os.Remove(inbox); err != nil {
					t.Fatal(err)
				}
				if mode == "inbox-symlink" {
					err = os.Symlink(outside, inbox)
				} else {
					err = os.WriteFile(inbox, []byte("owner inbox placeholder"), 0600)
				}
			case "archive-dir-symlink":
				err = os.Symlink(outside, filepath.Join(inbox, "archived"))
			case "write-failure":
				old := projectionFault
				t.Cleanup(func() { projectionFault = old })
				projectionFault = func(stage string) error {
					if stage == "notice" {
						return fmt.Errorf("injected publication failure")
					}
					return nil
				}
			}
			if err != nil {
				t.Fatal(err)
			}
			log, err := os.CreateTemp(t.TempDir(), "diagnostic")
			if err != nil {
				t.Fatal(err)
			}
			oldErr := os.Stderr
			os.Stderr = log
			defer func() { os.Stderr = oldErr; log.Close() }()
			for i := 0; i < 3; i++ {
				if _, err := Before(ctx, root, dir, run); err != nil {
					t.Fatal(err)
				}
				if _, _, err := protocol.QuotaMembers(dir, nil); err != nil {
					t.Fatal(err)
				}
			}
			output := string(cycle4Read(t, log.Name()))
			if strings.Count(output, "non-blocking; delivery unconfirmed") != 1 || strings.Contains(output, "corrupt") {
				t.Fatal(output)
			}
			if ok, err := quota.ReadApplied(dir, b.ID); !ok || err != nil {
				t.Fatal(ok, err)
			}
			if string(cycle4Read(t, target)) != "owner bytes" {
				t.Fatal("followed unsafe destination")
			}
			entries, err := os.ReadDir(outside)
			if err != nil || len(entries) != 1 {
				t.Fatal("published through symlink", entries, err)
			}
			cycle4TerminalCount(t, root, run, b.ID)
		})
	}
}

func TestQuotaCycle4ReceiptAndHistoryStillGate(t *testing.T) {
	for _, mode := range []string{"symlink-receipt", "directory-receipt", "changed-receipt", "history"} {
		t.Run(mode, func(t *testing.T) {
			root, dir, run := fixture(t, quota.NewPolicy(nil, nil))
			ctx := leaseFixture(t, dir, run)
			b, err := Settle(ctx, root, dir, run, "round-01", []string{"a", "b", "c", "d"}, candidates())
			if err != nil {
				t.Fatal(err)
			}
			path := filepath.Join(dir, "quota-applied", b.ID)
			if mode == "history" {
				path = filepath.Join(dir, quota.HistoryDir, "000001.json")
			}
			switch mode {
			case "symlink-receipt":
				if err := os.Rename(path, path+".saved"); err != nil {
					t.Fatal(err)
				}
				err = os.Symlink(path+".saved", path)
			case "directory-receipt":
				if err := os.Remove(path); err != nil {
					t.Fatal(err)
				}
				err = os.Mkdir(path, 0700)
			default:
				err = os.WriteFile(path, []byte("changed immutable evidence"), 0600)
			}
			if err != nil {
				t.Fatal(err)
			}
			if _, _, err := protocol.QuotaMembers(dir, nil); err == nil {
				t.Fatal("integrity gate removed")
			}
			if mode == "changed-receipt" {
				// Contradictory receipt bytes still require full checked replay.
				if _, err := Before(ctx, root, dir, run); err != nil {
					t.Fatal(err)
				}
			} else if _, err := Before(ctx, root, dir, run); err == nil {
				t.Fatal("unsafe receipt/history repaired")
			}
			cycle4TerminalCount(t, root, run, b.ID)
		})
	}
}

func TestQuotaCycle4ManualClarificationOwnerCopiesAndReceipt(t *testing.T) {
	for _, mode := range []string{"before-live", "before-archive", "after-live", "after-archive", "unsafe-notice", "unsafe-receipt"} {
		t.Run(mode, func(t *testing.T) {
			root, dir, _ := fixture(t, quota.Policy{Enabled: false, Scope: quota.KickoffAndMidIdea})
			h, err := quota.ReadHistory(dir)
			if err != nil {
				t.Fatal(err)
			}
			b := quota.NewRevision(h, h.Kickoff.RunID, h.Current, h.Policy(), quota.Revision{ManualPrompt: "historical fixture"}, time.Now())
			inbox := filepath.Join(root, protocol.DeckDir, "inbox")
			original := filepath.Join(inbox, "parley-to-user_"+b.ID+".md")
			legacy := []byte(b.LegacyManualNotice() + "\n## Owner note\nThis was a manual edit.\n")
			cycle4Write(t, original, legacy)
			id := b.ID + "-manual-authority"
			correction := filepath.Join(inbox, "parley-to-user_"+id+".md")
			if strings.HasPrefix(mode, "after-") || mode == "unsafe-receipt" {
				if err := publishNotice(root, dir, b, true); err != nil {
					t.Fatal(err)
				}
				if string(cycle4Read(t, correction)) != b.Notice() {
					t.Fatal("inaccurate clarification")
				}
			}
			path := correction
			if strings.HasSuffix(mode, "archive") {
				path = filepath.Join(inbox, "archived", filepath.Base(correction))
				if strings.HasPrefix(mode, "after-") {
					if err := os.Remove(correction); err != nil {
						t.Fatal(err)
					}
				}
			}
			if mode == "unsafe-notice" {
				if err := os.Mkdir(correction, 0700); err != nil {
					t.Fatal(err)
				}
			} else {
				cycle4Write(t, path, []byte("Owner replacement of clarification."))
			}
			if mode == "unsafe-receipt" {
				receipt := filepath.Join(dir, "quota-applied", id)
				if err := os.Rename(receipt, receipt+".saved"); err != nil {
					t.Fatal(err)
				}
				if err := os.Symlink(receipt+".saved", receipt); err != nil {
					t.Fatal(err)
				}
			}
			for i := 0; i < 2; i++ {
				err := publishNotice(root, dir, b, true)
				if (err != nil) != (mode == "unsafe-receipt") {
					t.Fatal(mode, err)
				}
			}
			if !bytes.Equal(legacy, cycle4Read(t, original)) {
				t.Fatal("historical notice changed")
			}
			if mode != "unsafe-notice" && string(cycle4Read(t, path)) != "Owner replacement of clarification." {
				t.Fatal("owner clarification changed")
			}
			if mode != "unsafe-receipt" {
				if ok, err := quota.ReadApplied(dir, id); !ok || err != nil {
					t.Fatal("clarification attempt missing receipt", ok, err)
				}
			}
		})
	}
}
