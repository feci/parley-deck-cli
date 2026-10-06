package membership

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"
	"testing"
	"time"

	"parley-deck-cli/internal/procctl"
	"parley-deck-cli/internal/protocol"
	"parley-deck-cli/internal/quota"
	"parley-deck-cli/internal/runmanifest"
	"parley-deck-cli/internal/telemetry"
)

func TestQuotaCycle2AppliedEditAndUnknownPendingNeverRewritten(t *testing.T) {
	for _, receipt := range []string{"applied", "missing", "corrupt"} {
		t.Run(receipt, func(t *testing.T) {
			root, dir, run := fixture(t, quota.NewPolicy(nil, nil))
			ctx := leaseFixture(t, dir, run)
			original, _ := os.ReadFile(filepath.Join(dir, "00-prompt.md"))
			b, err := Settle(ctx, root, dir, run, "round-01", []string{"a", "b", "c", "d"}, candidates())
			if err != nil {
				t.Fatal(err)
			}
			edit := original
			if receipt != "applied" {
				path := filepath.Join(dir, "quota-applied", b.ID)
				if receipt == "missing" {
					os.Remove(path)
				} else {
					os.WriteFile(path, []byte("unreadable-receipt"), 0600)
				}
				edit = bytes.Replace(original, []byte("[a, b, c, d]"), []byte("[a, b, stranger]"), 1)
			}
			path := filepath.Join(dir, "00-prompt.md")
			os.WriteFile(path, edit, 0600)
			for i := 0; i < 2; i++ {
				if _, err = Before(ctx, root, dir, run); err == nil || !strings.Contains(err.Error(), "parley quota revise") {
					t.Fatal(err)
				}
				after, _ := os.ReadFile(path)
				if !bytes.Equal(edit, after) {
					t.Fatal("owner edit silently overwritten")
				}
			}
			notes, _ := filepath.Glob(filepath.Join(root, "parley-deck/inbox/parley-to-user_quota-block-*.md"))
			if len(notes) != 1 {
				t.Fatal(notes)
			}
		})
	}
}

func TestQuotaCycle2MissingManifestBeforeManualCommitAndCheckedRecovery(t *testing.T) {
	root, dir, run := fixture(t, quota.Policy{Enabled: false, Scope: quota.KickoffAndMidIdea})
	path := filepath.Join(dir, "00-prompt.md")
	raw, _ := os.ReadFile(path)
	raw = bytes.Replace(raw, []byte("participants: [a, b, c, d]"), []byte("participants: [a, b, c]\nexcluded: [d — unavailable — quota exhausted — confirmed 2026-10-04]"), 1)
	os.WriteFile(path, raw, 0600)
	manifestPath := runmanifest.Path(root, run)
	manifest, _ := os.ReadFile(manifestPath)
	os.Remove(manifestPath)
	if err := RecordManual(root, dir); err == nil || !strings.Contains(err.Error(), "quota manifest") {
		t.Fatal(err)
	}
	h, err := quota.ReadHistory(dir)
	if err != nil || h.Revision != 0 {
		t.Fatal(h, err)
	}
	// Reproduce the pre-fix committed-but-incomplete state using the same import.
	v, err := protocol.InspectQuota(dir)
	if err != nil {
		t.Fatal(err)
	}
	if err = quota.CommitBatch(dir, *v.Manual); err != nil {
		t.Fatal(err)
	}
	if _, err = Before(context.Background(), root, dir, run); err == nil {
		t.Fatal("missing manifest dispatched")
	}
	if _, _, err = protocol.QuotaMembers(dir, nil); err == nil {
		t.Fatal("pending state allowed signoff")
	}
	if err = os.WriteFile(manifestPath, manifest, 0600); err != nil {
		t.Fatal(err)
	}
	for i := 0; i < 2; i++ {
		if _, err = Before(context.Background(), root, dir, run); err != nil {
			t.Fatal(err)
		}
	}
	h, err = quota.ReadHistory(dir)
	if err != nil || h.Revision != 1 {
		t.Fatal(h, err)
	}
}

func TestQuotaCycle2CrashSettlementProofDurabilityAndReplay(t *testing.T) {
	oldDead, oldGroup, oldWrite := crashDeadLocal, crashGroupStopped, crashWrite
	t.Cleanup(func() { crashDeadLocal, crashGroupStopped, crashWrite = oldDead, oldGroup, oldWrite })
	crashDeadLocal = func(host, boot string, pid int) bool {
		return host == "fixture-host" && boot == "fixture-boot" && (pid == 12345 || pid == 12346)
	}
	crashGroupStopped = func(pgid int) bool { return pgid == 12345 }
	for _, mode := range []string{"dead", "foreign-host", "foreign-boot", "unknown-host", "unknown-boot", "legacy", "live-writer", "live-supervisor", "live-group"} {
		t.Run(mode, func(t *testing.T) {
			root, dir, run := fixture(t, quota.NewPolicy(nil, nil))
			ctx := leaseFixture(t, dir, run)
			now := time.Now().Add(-time.Minute)
			pid := 12345
			rec := telemetry.Record{Type: "invocation.started", InvocationID: "orphan", Metadata: telemetry.Metadata{Idea: filepath.Base(dir)}, StartedAt: &now, PID: &pid, Writer: &telemetry.WriterIdentity{Host: "fixture-host", Boot: "fixture-boot", SupervisorPID: 12346, ProcessGroup: pid}}
			switch mode {
			case "foreign-host":
				rec.Writer.Host = "foreign"
			case "foreign-boot":
				rec.Writer.Boot = "foreign"
			case "unknown-host":
				rec.Writer.Host = ""
			case "unknown-boot":
				rec.Writer.Boot = ""
			case "legacy":
				rec.Writer = nil
			case "live-writer":
				pid = 54321
			case "live-supervisor":
				rec.Writer.SupervisorPID = 54321
			case "live-group":
				rec.Writer.ProcessGroup = 54321
			}
			inv := filepath.Join(root, ".parley-runtime/invocations/orphan")
			os.MkdirAll(inv, 0700)
			raw, _ := json.Marshal(rec)
			os.WriteFile(filepath.Join(inv, "started.json"), raw, 0600)
			if mode != "dead" {
				if _, err := Before(ctx, root, dir, run); err == nil {
					t.Fatal("unresolved writer admitted")
				}
				return
			}
			crashWrite = func(string, []byte, bool) error { return fmt.Errorf("injected settlement publication failure") }
			if err := RequireStopped(root, filepath.Base(dir)); err == nil {
				t.Fatal("failed settlement admitted")
			}
			crashWrite = oldWrite
			// Concurrent recovery is idempotent and never creates a normal terminal.
			var wg sync.WaitGroup
			errs := make(chan error, 6)
			for i := 0; i < 6; i++ {
				wg.Add(1)
				go func() { defer wg.Done(); errs <- RequireStopped(root, filepath.Base(dir)) }()
			}
			wg.Wait()
			close(errs)
			for err := range errs {
				if err != nil {
					t.Fatal(err)
				}
			}
			receipt := filepath.Join(inv, "crash-settlement.json")
			first, _ := os.ReadFile(receipt)
			if _, err := Before(ctx, root, dir, run); err != nil {
				t.Fatal(err)
			}
			second, _ := os.ReadFile(receipt)
			if !bytes.Equal(first, second) {
				t.Fatal("settlement rewritten")
			}
			if _, err := os.Stat(filepath.Join(inv, "terminal.json")); !os.IsNotExist(err) {
				t.Fatal("fabricated terminal")
			}
			os.WriteFile(filepath.Join(inv, "started.json"), append(raw, ' '), 0600)
			if err := RequireStopped(root, filepath.Base(dir)); err == nil {
				t.Fatal("changed start record accepted")
			}
		})
	}
}

func TestQuotaCycle2NativeCrashWriter(t *testing.T) {
	if root := os.Getenv("PARLEY_CYCLE2_CRASH_ROOT"); root != "" {
		cmd := exec.Command("/bin/sh", "-c", "sleep 0.2") // local stub, never a provider
		procctl.SetNewProcessGroup(cmd)
		inv, err := telemetry.Begin(filepath.Join(root, ".parley-runtime/invocations"), telemetry.Metadata{Idea: "native-crash"})
		if err != nil {
			t.Fatal(err)
		}
		if err = cmd.Start(); err != nil {
			t.Fatal(err)
		}
		if err = inv.Started(cmd.Process.Pid); err != nil {
			t.Fatal(err)
		}
		if err = cmd.Wait(); err != nil {
			t.Fatal(err)
		}
		os.Exit(0) // crash after writer exit, before normal terminal persistence
	}
	root := t.TempDir()
	cmd := exec.Command(os.Args[0], "-test.run=^TestQuotaCycle2NativeCrashWriter$")
	cmd.Env = append(os.Environ(), "PARLEY_CYCLE2_CRASH_ROOT="+root)
	if raw, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("%v %s", err, raw)
	}
	if err := RequireStopped(root, "native-crash"); err != nil {
		t.Fatal(err)
	}
	if err := RequireStopped(root, "native-crash"); err != nil {
		t.Fatal(err)
	}
	receipts, _ := filepath.Glob(filepath.Join(root, ".parley-runtime/invocations/*/crash-settlement.json"))
	if len(receipts) != 1 {
		t.Fatal(receipts)
	}
	t.Logf("native local-stub crash settled idempotently at %s", receipts[0])
}

func TestQuotaCycle2LegacyManualNoticePreservedWithoutAuthority(t *testing.T) {
	root, dir, run := fixture(t, quota.Policy{Enabled: false, Scope: quota.KickoffAndMidIdea})
	path := filepath.Join(dir, "00-prompt.md")
	raw, _ := os.ReadFile(path)
	raw = bytes.Replace(raw, []byte("participants: [a, b, c, d]"), []byte("participants: [a, b, c]\nexcluded: [d — unavailable — confirmed 2026-10-04]"), 1)
	os.WriteFile(path, raw, 0600)
	if err := RecordManual(root, dir); err != nil {
		t.Fatal(err)
	}
	h, err := quota.ReadHistory(dir)
	if err != nil {
		t.Fatal(err)
	}
	b := h.Batches[0]
	notice := filepath.Join(root, "parley-deck/inbox", "parley-to-user_"+b.ID+".md")
	legacy := fmt.Sprintf("---\nfrom: parley\nto: user\nidea: %s\nblocking: no\ntransition: %s\n---\n\nOwner-confirmed membership/policy revision. Current participants: %v. Policy: %+v.\n", b.Idea, b.ID, b.Decision.After, b.Policy)
	os.WriteFile(notice, []byte(legacy), 0600)
	for i := 0; i < 2; i++ {
		if _, err := Before(context.Background(), root, dir, run); err != nil {
			t.Fatal(err)
		}
	}
	after, _ := os.ReadFile(notice)
	if string(after) != legacy {
		t.Fatal("historical notice rewritten")
	}
	correction, _ := os.ReadFile(strings.TrimSuffix(notice, ".md") + "-manual-authority.md")
	if string(correction) != b.Notice() || b.Owner.Authority != nil {
		t.Fatal("missing manual clarification or invented authority")
	}
	os.WriteFile(notice, []byte(legacy+"tampered"), 0600)
	if _, err := Before(context.Background(), root, dir, run); err == nil {
		t.Fatal("unrecognized notice accepted")
	}
}
