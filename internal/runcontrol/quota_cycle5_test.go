package runcontrol

import (
	"context"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	"parley-deck-cli/internal/membership"
	"parley-deck-cli/internal/protocol"
	"parley-deck-cli/internal/quota"
	"parley-deck-cli/internal/runmanifest"
	"parley-deck-cli/internal/runstate"
)

func TestQuotaCycle5CreationScope(t *testing.T) {
	for _, alias := range []string{"deck", "ideas", "ancestor", "none"} {
		for _, mode := range []string{"mid", "off", "kickoff", "legacy"} {
			t.Run(alias+"/"+mode, func(t *testing.T) {
				root := t.TempDir()
				t.Setenv("PARLEY_HOME", t.TempDir())
				if err := protocol.InitWorkspace(root); err != nil {
					t.Fatal(err)
				}
				if alias == "deck" || alias == "ideas" {
					path := filepath.Join(root, protocol.DeckDir)
					if alias == "ideas" {
						path = filepath.Join(path, "ideas")
					}
					if err := os.Rename(path, path+"-physical"); err != nil {
						t.Fatal(err)
					}
					if err := os.Symlink(path+"-physical", path); err != nil {
						t.Fatal(err)
					}
				} else if alias == "ancestor" {
					link := filepath.Join(t.TempDir(), "workspace-alias")
					if err := os.Symlink(root, link); err != nil {
						t.Fatal(err)
					}
					root = link
				}
				p := quota.NewPolicy(nil, nil)
				var policy *quota.Policy = &p
				if mode == "off" {
					p.Enabled = false
				}
				if mode == "kickoff" {
					p.Scope = quota.KickoffOnly
				}
				if mode == "legacy" {
					policy = nil
				}
				run, err := Create(CreateOptions{Root: root, Task: "scope fixture", Participants: []string{"a", "b"}, QuotaPolicy: policy})
				refused := mode == "mid" && (alias == "deck" || alias == "ideas")
				if refused {
					if err == nil || !strings.Contains(err.Error(), "symlinked deck or idea scopes") {
						t.Fatal(err)
					}
					for _, dir := range []string{"ideas", "runs", "inbox"} {
						entries, e := os.ReadDir(filepath.Join(root, protocol.DeckDir, dir))
						if e != nil || len(entries) != 0 {
							t.Fatal("creation refusal wrote state", dir, entries, e)
						}
					}
					if _, e := os.Lstat(filepath.Join(root, ".parley-runtime")); !os.IsNotExist(e) {
						t.Fatal("creation refusal wrote staging", e)
					}
					return
				}
				if err != nil {
					t.Fatal(err)
				}
				_, release, err := membership.Acquire(context.Background(), run.Idea.Path, run.RunID)
				if err != nil {
					t.Fatal(err)
				}
				release()
			})
		}
	}
}

func TestQuotaCycle5KickoffNoticeNonBlocking(t *testing.T) {
	for _, mode := range []string{"present", "missing", "read-only", "inbox-symlink", "inbox-file", "archive-symlink"} {
		t.Run(mode, func(t *testing.T) {
			root := t.TempDir()
			t.Setenv("PARLEY_HOME", t.TempDir())
			if err := protocol.InitWorkspace(root); err != nil {
				t.Fatal(err)
			}
			inbox := filepath.Join(root, protocol.DeckDir, "inbox")
			outside := t.TempDir()
			switch mode {
			case "missing", "inbox-symlink", "inbox-file":
				if err := os.Remove(inbox); err != nil {
					t.Fatal(err)
				}
				if mode == "inbox-symlink" {
					if err := os.Symlink(outside, inbox); err != nil {
						t.Fatal(err)
					}
				}
				if mode == "inbox-file" {
					if err := os.WriteFile(inbox, []byte("owner bytes"), 0600); err != nil {
						t.Fatal(err)
					}
				}
			case "read-only":
				if err := os.Chmod(inbox, 0500); err != nil {
					t.Fatal(err)
				}
				t.Cleanup(func() { os.Chmod(inbox, 0700) })
				probe, e := os.CreateTemp(inbox, "permission-probe")
				if e == nil {
					probe.Close()
					os.Remove(probe.Name())
					t.Skip("filesystem does not enforce read-only directories")
				}
			case "archive-symlink":
				if err := os.Symlink(outside, filepath.Join(inbox, "archived")); err != nil {
					t.Fatal(err)
				}
			}
			log, err := os.CreateTemp(t.TempDir(), "diagnostic")
			if err != nil {
				t.Fatal(err)
			}
			old := os.Stderr
			os.Stderr = log
			defer func() { os.Stderr = old; log.Close() }()
			policy := quota.Policy{Enabled: true, Scope: quota.KickoffOnly}
			ev := quota.Evidence{InvocationID: "fixture", Adapter: "synthetic-test", RuleID: "synthetic-test", Provenance: "synthetic-test", Eligible: true}
			decision := quota.Evaluate(policy, []string{"a", "b", "c"}, []quota.Member{{ID: "a", Usable: true}, {ID: "b", Usable: true}, {ID: "c", Evidence: &ev}}, quota.Roles{})
			run, err := Create(CreateOptions{Root: root, Task: "kickoff notice fixture", Participants: decision.After, QuotaPolicy: &policy, QuotaDecision: &decision})
			if err != nil {
				t.Fatal("publication blocked run creation", err)
			}
			k, err := protocol.ReadQuotaState(run.Idea.Path)
			if err != nil || k.Transition == nil {
				t.Fatal(k, err)
			}
			m, err := runmanifest.Load(root, run.RunID)
			if err != nil || !reflect.DeepEqual(m.Participants, []string{"a", "b"}) {
				t.Fatal(m, err)
			}
			events, err := run.Store.Load()
			if err != nil {
				t.Fatal(err)
			}
			count := 0
			for _, event := range events {
				if event.Type == "run.created" {
					count++
				}
			}
			if count != 1 {
				t.Fatal("run.created", count)
			}
			for i := 0; i < 2; i++ {
				s, e := runstate.LoadRun(root, run.RunID)
				if e != nil || s.QuotaPending != "" || !reflect.DeepEqual(s.Participants, []string{"a", "b"}) {
					t.Fatal(s, e)
				}
			}
			safe := mode == "present" || mode == "missing"
			notes, _ := filepath.Glob(filepath.Join(inbox, "parley-to-user_kickoff-*.md"))
			want := 0
			if safe {
				want = 1
			}
			if len(notes) != want {
				t.Fatal(notes)
			}
			output, e := os.ReadFile(log.Name())
			if e != nil {
				t.Fatal(e)
			}
			diagnostics := strings.Count(string(output), "non-blocking; delivery unconfirmed")
			if (safe && diagnostics != 0) || (!safe && diagnostics != 1) {
				t.Fatal(string(output))
			}
			if entries, e := os.ReadDir(outside); e != nil || len(entries) != 0 {
				t.Fatal("published through alias", entries, e)
			}
			if mode == "inbox-file" {
				b, _ := os.ReadFile(inbox)
				if string(b) != "owner bytes" {
					t.Fatal("owner file changed")
				}
			}
		})
	}
}
