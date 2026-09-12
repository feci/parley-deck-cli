package budget

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
	"time"
)

func originFixture(t *testing.T) (Store, string, LockOriginPreview, LockOriginRequest) {
	t.Helper()
	ctx := context.Background()
	s := testStore(t)
	if _, e := s.Reserve(ctx, Request{ID: "spent", Kind: Fixup, ReserveMicros: micros(7)}, Limits{}); e != nil {
		t.Fatal(e)
	}
	if _, e := s.Settle(ctx, "spent", nil); e != nil {
		t.Fatal(e)
	}
	if _, e := s.Reserve(ctx, Request{ID: "orphan", Kind: Launch}, Limits{}); e != nil {
		t.Fatal(e)
	}
	release, e := AcquireResourceGuard(ctx, s.Dir)
	if e != nil {
		t.Fatal(e)
	}
	release()
	target := t.TempDir()
	p, e := PreviewLockOriginMigration(ctx, s.Dir, target)
	if e != nil {
		t.Fatal(e)
	}
	return s, target, p, LockOriginRequest{ExpectedSHA256: p.SHA256, DecisionID: "cache-move", Reason: "synthetic same-host relocation"}
}
func originRead(t *testing.T, path string) []byte {
	t.Helper()
	b, e := os.ReadFile(path)
	if e != nil {
		t.Fatal(e)
	}
	return b
}
func originApply(t *testing.T, s Store, target string, r LockOriginRequest) LockOriginPreview {
	t.Helper()
	p, e := ApplyLockOriginMigration(context.Background(), s.Dir, target, r)
	if e != nil {
		t.Fatal(e)
	}
	return p
}
func originGuard(t *testing.T, dir string) {
	t.Helper()
	release, e := AcquireResourceGuard(context.Background(), dir)
	if e != nil {
		t.Fatal(e)
	}
	release()
}

func TestOriginMigrationRetainsAccountingAndHistoricalReplay(t *testing.T) {
	s, target, p, r := originFixture(t)
	ctx := context.Background()
	ledger := originRead(t, filepath.Join(s.Dir, "ledger.json"))
	witness := originRead(t, filepath.Join(s.Dir, "ledger-established"))
	old, _ := parsePinnedOrigin([]byte(p.Origin))
	oldInfo, e := os.Stat(old.Path)
	if e != nil {
		t.Fatal(e)
	}
	got := originApply(t, s, target, r)
	if !reflect.DeepEqual(got, p) {
		t.Fatal("apply lost exact preview")
	}
	for name, want := range map[string][]byte{"ledger.json": ledger, "ledger-established": witness} {
		if !bytes.Equal(originRead(t, filepath.Join(s.Dir, name)), want) {
			t.Fatalf("migration rewrote %s", name)
		}
	}
	after, e := os.Stat(old.Path)
	if e != nil || !os.SameFile(oldInfo, after) {
		t.Fatal("old permanent kernel object was replaced")
	}
	raw := originRead(t, filepath.Join(s.Dir, "lock-origin"))
	if !bytes.Equal(raw, originRead(t, filepath.Join(s.Dir, resourceGuardWitness))) {
		t.Fatal("guard witness was not cut over")
	}
	if e := checkLockOriginLocation(s.Dir, old.Path); e == nil {
		t.Fatal("legacy v2 runtime accepted migrated origin")
	}
	// Normal current runtime follows the journal even with a different cache env.
	t.Setenv("XDG_CACHE_HOME", filepath.Join(t.TempDir(), "unavailable"))
	originGuard(t, s.Dir)
	if _, e := s.Reserve(ctx, Request{ID: "spent", Kind: Fixup}, Limits{}); !errors.Is(e, ErrReserved) {
		t.Fatalf("spent attempt was reset: %v", e)
	}
	if _, e := s.Reserve(ctx, Request{ID: "after", Kind: Fixup}, Limits{Actions: map[Kind]int{Fixup: 1}}); !errors.Is(e, ErrLimit) {
		t.Fatalf("cap was reset: %v", e)
	}
	if _, e := s.Reserve(ctx, Request{ID: "later", Kind: Fixup}, Limits{}); e != nil {
		t.Fatal(e)
	}
	secondTarget := t.TempDir()
	p2, e := PreviewLockOriginMigration(ctx, s.Dir, secondTarget)
	if e != nil {
		t.Fatal(e)
	}
	r2 := LockOriginRequest{ExpectedSHA256: p2.SHA256, DecisionID: "second-move", Reason: "synthetic second relocation"}
	originApply(t, s, secondTarget, r2)
	originGuard(t, s.Dir)
	if replay := originApply(t, s, target, r); !reflect.DeepEqual(replay, p) {
		t.Fatal("historical replay consulted later source")
	}
	changed := r
	changed.Reason = "different decision"
	if _, e := ApplyLockOriginMigration(ctx, s.Dir, target, changed); e == nil {
		t.Fatal("decision identity allowed a changed payload")
	}
	// Removing any ancestor completion invalidates current authority.
	if e := os.Remove(originCompletePath(p.Dir, key(r.DecisionID))); e != nil {
		t.Fatal(e)
	}
	if release, e := AcquireResourceGuard(ctx, s.Dir); e == nil {
		release()
		t.Fatal("missing ancestor completion authorized work")
	}
}

func TestOriginMigrationEveryPublicationBoundaryRecovers(t *testing.T) {
	for _, stage := range []string{"journal", "identity", "target-held", "origin", "witness", "complete"} {
		t.Run(stage, func(t *testing.T) {
			s, target, p, r := originFixture(t)
			ledger := originRead(t, filepath.Join(s.Dir, "ledger.json"))
			injected := errors.New("publication interrupted")
			_, e := applyLockOriginMigration(context.Background(), s.Dir, target, r, tryLock, func(at string) error {
				if at == stage {
					return injected
				}
				return nil
			})
			if !errors.Is(e, injected) {
				t.Fatalf("boundary not reached: %v", e)
			}
			release, e := AcquireResourceGuard(context.Background(), s.Dir)
			shouldBlock := stage == "origin" || stage == "witness"
			if shouldBlock && e == nil {
				release()
				t.Fatal("unfinished cutover granted work")
			}
			if !shouldBlock && e != nil {
				t.Fatal(e)
			}
			if release != nil && !shouldBlock {
				release()
			}
			if replay := originApply(t, s, target, r); !reflect.DeepEqual(replay, p) {
				t.Fatal("recovery changed preview")
			}
			originGuard(t, s.Dir)
			if !bytes.Equal(ledger, originRead(t, filepath.Join(s.Dir, "ledger.json"))) {
				t.Fatal("recovery changed accounting")
			}
		})
	}
}
func TestOriginMigrationStalePreviewAndPreparedDecision(t *testing.T) {
	for _, prepared := range []bool{false, true} {
		t.Run(fmt.Sprint(prepared), func(t *testing.T) {
			s, target, p, r := originFixture(t)
			if prepared {
				_, e := applyLockOriginMigration(context.Background(), s.Dir, target, r, tryLock, func(stage string) error {
					if stage == "journal" {
						return errors.New("stop")
					}
					return nil
				})
				if e == nil {
					t.Fatal("fault not reached")
				}
			}
			if _, e := s.Reserve(context.Background(), Request{ID: "after-preview", Kind: Fixup}, Limits{}); e != nil {
				t.Fatal(e)
			}
			if _, e := ApplyLockOriginMigration(context.Background(), s.Dir, target, r); e == nil {
				t.Fatal("stale accounting accepted")
			}
			if string(originRead(t, filepath.Join(s.Dir, "lock-origin"))) != p.Origin {
				t.Fatal("stale decision changed origin")
			}
			next, e := PreviewLockOriginMigration(context.Background(), s.Dir, target)
			if e != nil {
				t.Fatal(e)
			}
			replacement := LockOriginRequest{ExpectedSHA256: next.SHA256, DecisionID: "fresh-preview", Reason: "old v2 writer legitimately charged after interrupted preparation"}
			originApply(t, s, target, replacement)
			if prepared {
				if _, e := readOriginMigration(p.Dir, key(r.DecisionID)); e != nil {
					t.Fatalf("lost superseded prepared record: %v", e)
				}
			}
		})
	}
}
func TestOriginMigrationMissingIdentityAndTamperingRefuse(t *testing.T) {
	for _, scenario := range []string{"old-missing", "old-host", "target-changed", "record", "completion", "witness", "root-copy", "root-copy-v2", "origin-missing"} {
		t.Run(scenario, func(t *testing.T) {
			s, target, p, r := originFixture(t)
			old, _ := parsePinnedOrigin([]byte(p.Origin))
			ctx := context.Background()
			switch scenario {
			case "root-copy-v2":
				copyDir := t.TempDir()
				if e := os.CopyFS(copyDir, os.DirFS(s.Dir)); e != nil {
					t.Fatal(e)
				}
				s.Dir = copyDir
			case "old-missing":
				if e := os.Remove(old.Path); e != nil {
					t.Fatal(e)
				}
			case "old-host":
				parts := strings.Split(p.Origin, "\n")
				parts[1] = "other-host.invalid"
				if e := os.WriteFile(filepath.Join(s.Dir, "lock-origin"), []byte(strings.Join(parts, "\n")), 0600); e != nil {
					t.Fatal(e)
				}
			case "target-changed":
				if _, e := lockIdentity(p.TargetPath, true); e != nil {
					t.Fatal(e)
				}
			default:
				originApply(t, s, target, r)
				switch scenario {
				case "origin-missing":
					if e := os.Remove(filepath.Join(s.Dir, "lock-origin")); e != nil {
						t.Fatal(e)
					}
					if e := os.Remove(filepath.Join(s.Dir, resourceGuardWitness)); e != nil {
						t.Fatal(e)
					}
				case "record":
					if e := os.WriteFile(originRecordPath(p.Dir, key(r.DecisionID)), []byte("{}"), 0600); e != nil {
						t.Fatal(e)
					}
				case "completion":
					if e := os.WriteFile(originCompletePath(p.Dir, key(r.DecisionID)), []byte(strings.Repeat("0", 64)+"\n"), 0600); e != nil {
						t.Fatal(e)
					}
				case "witness":
					if e := os.Remove(filepath.Join(s.Dir, resourceGuardWitness)); e != nil {
						t.Fatal(e)
					}
				case "root-copy":
					copyDir := t.TempDir()
					if e := os.CopyFS(copyDir, os.DirFS(s.Dir)); e != nil {
						t.Fatal(e)
					}
					s.Dir = copyDir
				}
			}
			var release func()
			var e error
			if scenario == "old-missing" || scenario == "old-host" || scenario == "target-changed" || scenario == "root-copy-v2" {
				_, e = ApplyLockOriginMigration(ctx, s.Dir, target, r)
			} else {
				release, e = AcquireResourceGuard(ctx, s.Dir)
			}
			if e == nil {
				if release != nil {
					release()
				}
				t.Fatal("invalid authority allowed operation")
			}
		})
	}
}
func TestOriginMigrationDestinationExclusionAndChangedInode(t *testing.T) {
	for _, scenario := range []string{"busy", "false-exclusion", "old-replaced", "target-replaced"} {
		t.Run(scenario, func(t *testing.T) {
			s, target, p, r := originFixture(t)
			take := tryLock
			if scenario == "busy" {
				take = func(*os.File) (bool, error) { return false, nil }
			}
			if scenario == "false-exclusion" {
				take = func(*os.File) (bool, error) { return true, nil }
			}
			hook := func(stage string) error {
				if stage != "target-held" || scenario != "old-replaced" && scenario != "target-replaced" {
					return nil
				}
				path := p.TargetPath
				if scenario == "old-replaced" {
					o, _ := parsePinnedOrigin([]byte(p.Origin))
					path = o.Path
				}
				raw := originRead(t, path)
				return writeSynced(path, raw)
			}
			_, e := applyLockOriginMigration(context.Background(), s.Dir, target, r, take, hook)
			if e == nil {
				t.Fatal("invalid kernel exclusion allowed cutover")
			}
			if string(originRead(t, filepath.Join(s.Dir, "lock-origin"))) != p.Origin {
				t.Fatal("unverified exclusion changed origin")
			}
		})
	}
}

// Process fixture: barriers report a genuinely failed take, not an arbitrary
// delay. The legacy path preserves the v2 post-wait origin check explicitly.
func TestOriginMigrationChild(t *testing.T) {
	mode := os.Getenv("PARLEY_ORIGIN_CHILD")
	if mode == "" {
		return
	}
	dir := os.Getenv("PARLEY_ORIGIN_DIR")
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	scanner := bufio.NewScanner(os.Stdin)
	if mode == "target-holder" {
		path := os.Getenv("PARLEY_ORIGIN_TARGET")
		token, e := lockIdentity(path, false)
		if e != nil {
			t.Fatal(e)
		}
		release, e := acquirePinnedKernelLock(ctx, path, token, tryLock, unlock, func() error { return nil }, nil)
		if e != nil {
			t.Fatal(e)
		}
		defer release()
		fmt.Println("held")
		if !scanner.Scan() {
			t.Fatal("missing destination release")
		}
		return
	}
	if mode == "migrate" {
		target := os.Getenv("PARLEY_ORIGIN_TARGET")
		p, e := PreviewLockOriginMigration(ctx, dir, target)
		if e != nil {
			t.Fatal(e)
		}
		r := LockOriginRequest{ExpectedSHA256: p.SHA256, DecisionID: "process-move", Reason: "process fixture"}
		_, e = applyLockOriginMigration(ctx, dir, target, r, tryLock, func(stage string) error {
			if stage == "origin" {
				fmt.Println("origin-held")
				if !scanner.Scan() {
					return errors.New("missing migration release")
				}
			}
			return nil
		})
		if e != nil {
			t.Fatal(e)
		}
		fmt.Println("migrated")
		return
	}
	signaled := false
	var mainFile *os.File
	take := func(f *os.File) (bool, error) {
		if mainFile == nil {
			mainFile = f
		}
		held, e := tryLock(f)
		if f == mainFile && !held && e == nil && !signaled {
			fmt.Println("waiting")
			signaled = true
			if mode == "legacy" && !scanner.Scan() {
				return false, errors.New("missing legacy waiter continuation")
			}
		}
		return held, e
	}
	var release func()
	var e error
	if mode == "legacy" {
		raw := originRead(t, filepath.Join(dir, "lock-origin"))
		o, parseErr := parsePinnedOrigin(raw)
		if parseErr != nil || o.Migration != "" {
			t.Fatal("legacy source did not begin at v2")
		}
		release, e = acquirePinnedKernelLock(ctx, o.Path, o.Token, take, unlock, func() error { return verifyPinnedLockOrigin(dir, o.Path, o.Token) }, nil)
	} else {
		release, e = lockWithReady(ctx, filepath.Join(dir, "resource"), take, unlock, func(c string) error { return establishResourceGuard(c, publishExclusive) })
	}
	if mode == "legacy" && e != nil {
		fmt.Println("legacy-refused")
		return
	}
	if e != nil {
		fmt.Println("refused")
		return
	}
	defer release()
	fmt.Println("held")
	if !scanner.Scan() {
		t.Fatal("missing parent release")
	}
}

type originChild struct {
	cmd    *exec.Cmd
	input  io.WriteCloser
	output *bufio.Scanner
	errors bytes.Buffer
	waited bool
}

func startOriginChild(t *testing.T, mode, dir, target string) *originChild {
	t.Helper()
	c := &originChild{}
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	t.Cleanup(cancel)
	c.cmd = exec.CommandContext(ctx, os.Args[0], "-test.run=^TestOriginMigrationChild$")
	c.cmd.Env = append(os.Environ(), "PARLEY_ORIGIN_CHILD="+mode, "PARLEY_ORIGIN_DIR="+dir, "PARLEY_ORIGIN_TARGET="+target)
	var e error
	c.input, e = c.cmd.StdinPipe()
	if e != nil {
		t.Fatal(e)
	}
	out, e := c.cmd.StdoutPipe()
	if e != nil {
		t.Fatal(e)
	}
	c.output = bufio.NewScanner(out)
	c.cmd.Stderr = &c.errors
	if e = c.cmd.Start(); e != nil {
		t.Fatal(e)
	}
	t.Cleanup(func() {
		if !c.waited {
			fmt.Fprintln(c.input, "release")
			c.input.Close()
			if e := c.cmd.Wait(); e != nil {
				t.Errorf("child cleanup: %v %s", e, c.errors.String())
			}
		}
	})
	return c
}
func (c *originChild) line(t *testing.T, want string) {
	t.Helper()
	if !c.output.Scan() || c.output.Text() != want {
		t.Fatalf("child expected %q got %q", want, c.output.Text())
	}
}
func (c *originChild) finish(t *testing.T) {
	t.Helper()
	fmt.Fprintln(c.input, "release")
	c.input.Close()
	for c.output.Scan() {
	}
	c.waited = true
	if e := c.cmd.Wait(); e != nil {
		t.Fatalf("child: %v %s", e, c.errors.String())
	}
}

func TestOriginMigrationProcessHoldersWaitersAndCutover(t *testing.T) {
	s, target, _, _ := originFixture(t)
	holder := startOriginChild(t, "holder", s.Dir, "")
	holder.line(t, "held")
	legacy := startOriginChild(t, "legacy", s.Dir, "")
	legacy.line(t, "waiting")
	// A read-only preview also must wait on the actual old holder.
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	_, e := PreviewLockOriginMigration(ctx, s.Dir, target)
	cancel()
	if !errors.Is(e, ErrLockContention) {
		t.Fatalf("preview bypassed held old lock: %v", e)
	}
	holder.finish(t)
	migration := startOriginChild(t, "migrate", s.Dir, target)
	migration.line(t, "origin-held")
	fmt.Fprintln(legacy.input, "continue pinned waiter")
	fresh := startOriginChild(t, "new", s.Dir, "")
	fresh.line(t, "refused")
	fresh.finish(t)
	migration.finish(t)
	legacy.line(t, "legacy-refused")
	legacy.finish(t)
	fresh = startOriginChild(t, "new", s.Dir, "")
	fresh.line(t, "held")
	fresh.finish(t)
}

func TestOriginMigrationPinnedWaiterCannotAdoptNewOrigin(t *testing.T) {
	s, target, p, r := originFixture(t)
	old, _ := parsePinnedOrigin([]byte(p.Origin))
	pinned := make(chan struct{})
	resume := make(chan struct{})
	result := make(chan error, 1)
	go func() {
		calls := 0
		take := func(f *os.File) (bool, error) {
			calls++
			if calls == 1 {
				close(pinned)
				<-resume
			}
			return tryLock(f)
		}
		release, e := lockWithOps(context.Background(), filepath.Join(s.Dir, "resource"), take, unlock)
		if release != nil {
			release()
		}
		result <- e
	}()
	<-pinned
	originApply(t, s, target, r)
	close(resume)
	if e := <-result; e == nil {
		t.Fatal("old waiter adopted new origin")
	}
	f, e := os.OpenFile(old.Path, os.O_RDWR, 0600)
	if e != nil {
		t.Fatal(e)
	}
	defer f.Close()
	held, e := tryLock(f)
	if e != nil || !held {
		t.Fatal("failed waiter leaked old lock")
	}
	unlock(f)
}

func TestOriginMigrationRealDestinationHolderRefusesAndResumes(t *testing.T) {
	s, target, _, _ := originFixture(t)
	targetResource := t.TempDir()
	// Build the target identity using its permanent resource key, then ask a
	// different OS process to hold that exact kernel object.
	canonical, e := canonicalOriginDirectory(s.Dir)
	if e != nil {
		t.Fatal(e)
	}
	targetPath, e := originTarget(canonical, target)
	if e != nil {
		t.Fatal(e)
	}
	if _, e := lockIdentity(targetPath, true); e != nil {
		t.Fatal(e)
	}
	holder := startOriginChild(t, "target-holder", targetResource, targetPath)
	holder.line(t, "held")
	p, e := PreviewLockOriginMigration(context.Background(), s.Dir, target)
	if e != nil {
		t.Fatal(e)
	}
	r := LockOriginRequest{ExpectedSHA256: p.SHA256, DecisionID: "busy-destination", Reason: "actual destination exclusion fixture"}
	if _, e := ApplyLockOriginMigration(context.Background(), s.Dir, target, r); !errors.Is(e, ErrLockContention) {
		t.Fatalf("busy target was not refused: %v", e)
	}
	if string(originRead(t, filepath.Join(s.Dir, "lock-origin"))) != p.Origin {
		t.Fatal("destination holder did not prevent cutover")
	}
	holder.finish(t)
	originApply(t, s, target, r)
	originGuard(t, s.Dir)
}

func TestOriginMigrationChainCapacityRefusesBeforeCutover(t *testing.T) {
	s, _, p, _ := originFixture(t)
	// Build a valid completed chain directly; this models retained history, not
	// 127 operator activations. Every synthetic decision binds its predecessor.
	dir := p.Dir
	if e := os.Mkdir(filepath.Join(dir, originMigrationDirectory), 0700); e != nil {
		t.Fatal(e)
	}
	targets := []string{t.TempDir(), t.TempDir()}
	raw := []byte(p.Origin)
	var last originMigration
	for i := 0; i < 127; i++ {
		target, e := originTarget(dir, targets[i%2])
		if e != nil {
			t.Fatal(e)
		}
		token, e := lockIdentity(target, true)
		if e != nil {
			t.Fatal(e)
		}
		preview := p
		preview.Origin = string(raw)
		preview.TargetPath = target
		preview.TargetToken = token
		preview.SHA256 = preview.digest()
		last = originMigration{Version: 1, Preview: preview, TargetToken: token, Request: LockOriginRequest{ExpectedSHA256: preview.SHA256, DecisionID: fmt.Sprintf("chain-%d", i), Reason: "synthetic retained chain"}}
		b, e := json.Marshal(last)
		if e != nil {
			t.Fatal(e)
		}
		if e = publishOriginExclusive(originRecordPath(dir, last.reference()), b); e != nil {
			t.Fatal(e)
		}
		if e = publishOriginExclusive(originCompletePath(dir, last.reference()), []byte(migrationDigest(last)+"\n")); e != nil {
			t.Fatal(e)
		}
		raw = last.origin()
	}
	if e := writeSynced(filepath.Join(dir, "lock-origin"), raw); e != nil {
		t.Fatal(e)
	}
	if e := writeSynced(filepath.Join(dir, resourceGuardWitness), raw); e != nil {
		t.Fatal(e)
	}
	originGuard(t, s.Dir)
	target := t.TempDir()
	if _, e := PreviewLockOriginMigration(context.Background(), dir, target); e == nil {
		t.Fatal("inspection offered an unreadable next chain")
	}
	if !bytes.Equal(raw, originRead(t, filepath.Join(dir, "lock-origin"))) {
		t.Fatal("capacity refusal changed authority")
	}
	originGuard(t, s.Dir)
}

func TestOriginMigrationJournalAlonePreventsFreshGuardBootstrap(t *testing.T) {
	dir := t.TempDir()
	target := t.TempDir()
	originGuard(t, dir)
	p, e := PreviewLockOriginMigration(context.Background(), dir, target)
	if e != nil {
		t.Fatal(e)
	}
	r := LockOriginRequest{ExpectedSHA256: p.SHA256, DecisionID: "guard-only", Reason: "retained journal continuity fixture"}
	if _, e = ApplyLockOriginMigration(context.Background(), dir, target, r); e != nil {
		t.Fatal(e)
	}
	for _, name := range []string{"lock-origin", resourceGuardWitness} {
		if e = os.Remove(filepath.Join(dir, name)); e != nil {
			t.Fatal(e)
		}
	}
	if release, e := AcquireResourceGuard(context.Background(), dir); e == nil {
		release()
		t.Fatal("retained migration history became a fresh guard")
	}
	if _, e = os.Lstat(filepath.Join(dir, "lock-origin")); !os.IsNotExist(e) {
		t.Fatal("lost migrated origin was recreated")
	}
}
