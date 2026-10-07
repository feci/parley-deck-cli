package membership

import (
	"context"
	"crypto/sha256"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
	"time"

	"parley-deck-cli/internal/protocol"
	"parley-deck-cli/internal/quota"
)

func cycle5Tree(t *testing.T, root string) map[string]string {
	t.Helper()
	out := map[string]string{}
	err := filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		rel, err := filepath.Rel(root, path)
		if err != nil {
			return err
		}
		if d.Type()&os.ModeSymlink != 0 {
			target, e := os.Readlink(path)
			out[rel] = "link:" + target
			return e
		}
		if d.IsDir() {
			out[rel] = "dir"
			return nil
		}
		b, e := os.ReadFile(path)
		out[rel] = fmt.Sprintf("%x", sha256.Sum256(b))
		return e
	})
	if err != nil {
		t.Fatal(err)
	}
	return out
}

func TestQuotaCycle5AliasedScopeRefusesWithoutWrites(t *testing.T) {
	for _, part := range []string{"deck", "ideas", "idea"} {
		for _, broken := range []bool{false, true} {
			t.Run(fmt.Sprintf("%s/broken=%v", part, broken), func(t *testing.T) {
				root, dir, run := fixture(t, quota.NewPolicy(nil, nil))
				// Deliberately misleading outer deck name. An ancestor search must never
				// place lease state in its parent, outside this workspace.
				sandbox := t.TempDir()
				moved := filepath.Join(sandbox, "parley-deck", "workspace")
				if err := os.MkdirAll(filepath.Dir(moved), 0700); err != nil {
					t.Fatal(err)
				}
				slug := filepath.Base(dir)
				if err := os.Rename(root, moved); err != nil {
					t.Fatal(err)
				}
				root = moved
				dir = filepath.Join(root, protocol.DeckDir, "ideas", slug)
				path := dir
				if part == "ideas" {
					path = filepath.Dir(dir)
				}
				if part == "deck" {
					path = filepath.Dir(filepath.Dir(dir))
				}
				if err := os.Rename(path, path+"-physical"); err != nil {
					t.Fatal(err)
				}
				if err := os.Symlink(path+"-physical", path); err != nil {
					t.Fatal(err)
				}
				if broken {
					cycle4Write(t, filepath.Join(dir, "quota-kickoff.json"), []byte("broken history"))
				}
				before := cycle5Tree(t, sandbox)
				if _, release, err := Acquire(context.Background(), dir, run); err == nil {
					release()
					t.Fatal("aliased scope acquired")
				} else if !strings.Contains(err.Error(), "symlinked deck or idea scopes") {
					t.Fatal(err)
				}
				for _, kind := range []string{"driver", "projection", "revision"} {
					if _, err := lockPath(dir, kind); err == nil || !strings.Contains(err.Error(), "symlinked deck or idea scopes") {
						t.Fatal(kind, err)
					}
				}
				if !broken {
					if release, err := ProjectionLock(dir); err == nil {
						release()
						t.Fatal("aliased projection acquired")
					}
				}
				if _, err := Before(context.Background(), root, dir, run); err == nil {
					t.Fatal("aliased scoped recovery accepted")
				}
				if !reflect.DeepEqual(before, cycle5Tree(t, sandbox)) {
					t.Fatal("refusal wrote runtime or inbox state")
				}
			})
		}
	}
}

func TestQuotaCycle5AncestorAliasSharesLease(t *testing.T) {
	root, dir, run := fixture(t, quota.NewPolicy(nil, nil))
	alias := filepath.Join(t.TempDir(), "workspace-alias")
	if err := os.Symlink(root, alias); err != nil {
		t.Fatal(err)
	}
	aliased := filepath.Join(alias, protocol.DeckDir, "ideas", filepath.Base(dir))
	original, err := lockPath(dir, "driver")
	if err != nil {
		t.Fatal(err)
	}
	viaAlias, err := lockPath(aliased, "driver")
	if err != nil || viaAlias != original {
		t.Fatal(original, viaAlias, err)
	}
	physical, err := filepath.EvalSymlinks(root)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.HasPrefix(original, filepath.Join(physical, ".parley-runtime", "membership")+string(filepath.Separator)) {
		t.Fatal(original)
	}
	ctx, release, err := Acquire(context.Background(), dir, run)
	if err != nil {
		t.Fatal(err)
	}
	defer release()
	nested, done, err := Acquire(ctx, aliased, run)
	if err != nil {
		t.Fatal(err)
	}
	done()
	if err := CheckLease(nested, aliased, run); err != nil {
		t.Fatal(err)
	}
	if _, r, err := Acquire(context.Background(), aliased, "rival"); err == nil {
		r()
		t.Fatal("alias bypassed lifetime lease")
	}
}

func TestQuotaCycle5AliasedOffScopeDrivingUnchanged(t *testing.T) {
	for _, mode := range []string{"legacy", "off", "kickoff-only"} {
		t.Run(mode, func(t *testing.T) {
			policy := quota.Policy{Enabled: false, Scope: quota.KickoffAndMidIdea}
			if mode == "kickoff-only" {
				policy = quota.Policy{Enabled: true, Scope: quota.KickoffOnly}
			}
			var root, dir, run string
			if mode == "legacy" {
				root = t.TempDir()
				t.Setenv("PARLEY_HOME", t.TempDir())
				if err := protocol.InitWorkspace(root); err != nil {
					t.Fatal(err)
				}
				idea, err := protocol.CreateIdeaFull(root, "legacy fixture", []string{"a", "b"}, nil, "deliberation", "")
				if err != nil {
					t.Fatal(err)
				}
				dir, run = idea.Path, "legacy"
			} else {
				root, dir, run = fixture(t, policy)
			}
			deck := filepath.Join(root, protocol.DeckDir)
			if err := os.Rename(deck, deck+"-physical"); err != nil {
				t.Fatal(err)
			}
			if err := os.Symlink(deck+"-physical", deck); err != nil {
				t.Fatal(err)
			}
			before := cycle5Tree(t, root)
			_, a, err := Acquire(context.Background(), dir, run)
			if err != nil {
				t.Fatal(err)
			}
			defer a()
			_, b, err := Acquire(context.Background(), dir, "rival")
			if err != nil {
				t.Fatal(err)
			}
			b()
			if !reflect.DeepEqual(before, cycle5Tree(t, root)) {
				t.Fatal("off-scope driving wrote a lease")
			}
		})
	}
}

func TestQuotaCycle5ManualClarificationUnsafePath(t *testing.T) {
	for _, applied := range []bool{false, true} {
		t.Run(fmt.Sprint(applied), func(t *testing.T) {
			root, dir, _ := fixture(t, quota.Policy{Enabled: false, Scope: quota.KickoffAndMidIdea})
			h, err := quota.ReadHistory(dir)
			if err != nil {
				t.Fatal(err)
			}
			b := quota.NewRevision(h, h.Kickoff.RunID, h.Current, h.Policy(), quota.Revision{ManualPrompt: "historical fixture"}, time.Now())
			inbox := filepath.Join(root, protocol.DeckDir, "inbox")
			if err := os.Remove(inbox); err != nil {
				t.Fatal(err)
			}
			target := t.TempDir()
			if err := os.Symlink(target, inbox); err != nil {
				t.Fatal(err)
			}
			log, err := os.CreateTemp(t.TempDir(), "diagnostic")
			if err != nil {
				t.Fatal(err)
			}
			old := os.Stderr
			os.Stderr = log
			defer func() { os.Stderr = old; log.Close() }()
			for i := 0; i < 2; i++ {
				if err := publishNotice(root, dir, b, applied); err != nil {
					t.Fatal(err)
				}
			}
			perCall := 1
			if !applied {
				perCall = 2
			}
			output := string(cycle4Read(t, log.Name()))
			if strings.Count(output, "non-blocking; delivery unconfirmed") != 2*perCall {
				t.Fatal(output)
			}
			receipt := filepath.Join(dir, "quota-applied", b.ID+"-manual-authority")
			if _, err := os.Lstat(receipt); !os.IsNotExist(err) {
				t.Fatal("unneeded clarification receipt", err)
			}
			if entries, err := os.ReadDir(target); err != nil || len(entries) != 0 {
				t.Fatal("followed inbox alias", entries, err)
			}
			// A real receipt path error is still returned before the publication-only exit.
			cycle4Write(t, receipt+".saved", []byte(b.ID+"-manual-authority\n"))
			if err := os.Symlink(receipt+".saved", receipt); err != nil {
				t.Fatal(err)
			}
			if err := publishNotice(root, dir, b, true); err == nil {
				t.Fatal("unsafe receipt ignored")
			}
			if err := os.Remove(receipt); err != nil {
				t.Fatal(err)
			}
			// Restoring a safe inbox reveals and corrects a historical label exactly once.
			if err := os.Remove(inbox); err != nil {
				t.Fatal(err)
			}
			cycle4Write(t, filepath.Join(inbox, "parley-to-user_"+b.ID+".md"), []byte(b.LegacyManualNotice()))
			for i := 0; i < 2; i++ {
				if err := publishNotice(root, dir, b, true); err != nil {
					t.Fatal(err)
				}
			}
			correction := filepath.Join(inbox, "parley-to-user_"+b.ID+"-manual-authority.md")
			if string(cycle4Read(t, correction)) != b.Notice() {
				t.Fatal("missing historical clarification")
			}
			if ok, err := quota.ReadApplied(dir, b.ID+"-manual-authority"); !ok || err != nil {
				t.Fatal(ok, err)
			}
			cycle4Write(t, receipt, []byte("malformed receipt"))
			if err := publishNotice(root, dir, b, true); err != nil {
				t.Fatal(err)
			}
			if ok, err := quota.ReadApplied(dir, b.ID+"-manual-authority"); !ok || err != nil {
				t.Fatal("receipt not checked/repaired", ok, err)
			}
		})
	}
}
