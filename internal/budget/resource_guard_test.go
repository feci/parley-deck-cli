package budget

import (
	"bufio"
	"bytes"
	"context"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
	"time"
)

func TestResourceGuardChild(t *testing.T) {
	dir := os.Getenv("PARLEY_TEST_GUARD_CHILD_DIR")
	if dir == "" {
		return
	}
	release, err := AcquireResourceGuard(context.Background(), dir)
	if err != nil {
		t.Fatal(err)
	}
	defer release()
	fmt.Println("guard-held")
	if !bufio.NewScanner(os.Stdin).Scan() {
		t.Fatal("parent did not release the held guard")
	}
}

// The first holder is a different OS process with a still-open descriptor.
// A missing pathname is not evidence that this holder has exited.
func TestResourceGuardMissingOriginCannotSplitLiveLock(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("Windows prevents this POSIX open-inode unlink scenario")
	}
	dir := t.TempDir()
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	t.Cleanup(cancel)
	cmd := exec.CommandContext(ctx, os.Args[0], "-test.run=^TestResourceGuardChild$")
	cmd.Env = append(os.Environ(), "PARLEY_TEST_GUARD_CHILD_DIR="+dir)
	input, err := cmd.StdinPipe()
	if err != nil {
		t.Fatal(err)
	}
	output, err := cmd.StdoutPipe()
	if err != nil {
		t.Fatal(err)
	}
	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	if err := cmd.Start(); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		_, _ = fmt.Fprintln(input, "release")
		_ = input.Close()
		if err := cmd.Wait(); err != nil {
			t.Errorf("holder failed: %v: %s", err, stderr.String())
		}
	})
	ready := make(chan string, 1)
	go func() {
		s := bufio.NewScanner(output)
		if s.Scan() {
			ready <- s.Text()
		} else {
			ready <- ""
		}
	}()
	select {
	case line := <-ready:
		if line != "guard-held" {
			t.Fatalf("holder did not acquire the guard: %q", line)
		}
	case <-ctx.Done():
		t.Fatal("holder readiness deadline")
	}
	lockPath := testPinnedLockPath(t, dir)
	origin := filepath.Join(dir, "lock-origin")
	for _, path := range []string{origin, lockPath} {
		if err := os.Remove(path); err != nil {
			t.Fatal(err)
		}
	}
	second, err := AcquireResourceGuard(ctx, dir)
	if second != nil {
		second()
		t.Fatal("second process acquired a new guard while the first still held its unlinked inode")
	}
	if err == nil || !strings.Contains(err.Error(), "has no lock origin") {
		t.Fatalf("missing origin must refuse with retained continuity: %v", err)
	}
	for _, path := range []string{origin, lockPath} {
		if _, err := os.Lstat(path); !os.IsNotExist(err) {
			t.Fatalf("refusal recreated %s: %v", path, err)
		}
	}
}

func TestLockOriginChangedDuringAcquisitionRefuses(t *testing.T) {
	for _, stage := range []int{1, 2} {
		for _, mutation := range []string{"removed", "token", "host"} {
			t.Run(fmt.Sprintf("take-%d-%s", stage, mutation), func(t *testing.T) {
				dir := t.TempDir()
				origin := filepath.Join(dir, "lock-origin")
				calls := 0
				var lockPath string
				var changed []byte
				take := func(f *os.File) (bool, error) {
					lockPath = f.Name()
					calls++
					if calls == stage {
						if mutation == "removed" {
							if err := os.Remove(origin); err != nil {
								return false, err
							}
						} else {
							data, err := os.ReadFile(origin)
							if err != nil {
								return false, err
							}
							parts := strings.Split(string(data), "\n")
							if mutation == "token" {
								parts[3] = strings.Repeat("f", 64)
							} else {
								parts[1] = "changed-host.invalid"
							}
							changed = []byte(strings.Join(parts, "\n"))
							if err := os.WriteFile(origin, changed, 0600); err != nil {
								return false, err
							}
						}
					}
					return tryLock(f)
				}
				release, err := lockWithOps(context.Background(), filepath.Join(dir, "resource"), take, unlock)
				if release != nil {
					release()
					t.Fatal("changed pinned origin authorized work")
				}
				if err == nil {
					t.Fatal("missing origin refusal")
				}
				data, readErr := os.ReadFile(origin)
				if mutation == "removed" && !os.IsNotExist(readErr) || mutation != "removed" && (readErr != nil || !bytes.Equal(data, changed)) {
					t.Fatalf("refusal repaired or replaced origin: %v", readErr)
				}
				f, err := os.OpenFile(lockPath, os.O_RDWR, 0600)
				if err != nil {
					t.Fatal(err)
				}
				defer f.Close()
				held, err := tryLock(f)
				if err != nil || !held {
					t.Fatalf("refusal leaked kernel ownership: %v", err)
				}
				unlock(f)
			})
		}
	}
}

func TestResourceGuardContinuityPublicationMustComplete(t *testing.T) {
	for _, after := range []bool{false, true} {
		t.Run(fmt.Sprintf("after-publication-%v", after), func(t *testing.T) {
			dir := t.TempDir()
			failure := errors.New("injected continuity publication failure")
			publish := func(staged, path string) error {
				if after {
					if err := publishExclusive(staged, path); err != nil {
						return err
					}
				}
				return failure
			}
			release, err := acquireResourceGuard(context.Background(), dir, publish)
			if release != nil {
				release()
				t.Fatal("incomplete continuity publication granted work")
			}
			if !errors.Is(err, failure) {
				t.Fatalf("lost publication failure: %v", err)
			}
			origin, err := os.ReadFile(filepath.Join(dir, "lock-origin"))
			if err != nil {
				t.Fatal(err)
			}
			witness, err := os.ReadFile(filepath.Join(dir, resourceGuardWitness))
			if after && (err != nil || !bytes.Equal(origin, witness)) || !after && !os.IsNotExist(err) {
				t.Fatalf("unexpected witness publication: %v", err)
			}
			// Recovery uses the same kernel identity and exact origin; it does
			// not reset accounting or recreate a lost established lock.
			ctx, cancel := context.WithTimeout(context.Background(), time.Second)
			defer cancel()
			release, err = AcquireResourceGuard(ctx, dir)
			if err != nil {
				t.Fatalf("same-origin retry failed or lock leaked: %v", err)
			}
			release()
			witness, err = os.ReadFile(filepath.Join(dir, resourceGuardWitness))
			if err != nil || !bytes.Equal(origin, witness) {
				t.Fatalf("retry changed origin: %v", err)
			}
			entries, err := os.ReadDir(dir)
			if err != nil || len(entries) != 2 || entries[0].Name() != resourceGuardWitness || entries[1].Name() != "lock-origin" {
				t.Fatalf("guard created ledger or leaked staging: %v %v", entries, err)
			}
		})
	}
}

func TestResourceGuardWitnessCannotChangeOriginOrAuthorizeCancelledWork(t *testing.T) {
	for _, scenario := range []string{"changed-origin", "cancelled", "conflicting-witness"} {
		t.Run(scenario, func(t *testing.T) {
			dir := t.TempDir()
			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()
			publish := func(staged, path string) error {
				if err := publishExclusive(staged, path); err != nil {
					return err
				}
				if scenario == "cancelled" {
					cancel()
					return nil
				}
				target := filepath.Join(dir, "lock-origin")
				if scenario == "conflicting-witness" {
					target = path
				}
				data, err := os.ReadFile(target)
				if err != nil {
					return err
				}
				return os.WriteFile(target, append(data, '\n'), 0600)
			}
			release, err := acquireResourceGuard(ctx, dir, publish)
			if release != nil {
				release()
				t.Fatal("changed or cancelled publication granted a guard")
			}
			if err == nil || scenario == "cancelled" && !errors.Is(err, context.Canceled) {
				t.Fatalf("publication refusal: %v", err)
			}
			before, err := os.ReadFile(filepath.Join(dir, resourceGuardWitness))
			if err != nil {
				t.Fatal(err)
			}
			fresh, stop := context.WithTimeout(context.Background(), time.Second)
			defer stop()
			release, err = AcquireResourceGuard(fresh, dir)
			if scenario == "cancelled" {
				if err != nil {
					t.Fatal(err)
				}
				release()
			} else if err == nil || release != nil {
				if release != nil {
					release()
				}
				t.Fatal("retry silently repaired conflicting identity")
			}
			after, err := os.ReadFile(filepath.Join(dir, resourceGuardWitness))
			if err != nil || !bytes.Equal(before, after) {
				t.Fatalf("retry rewrote witness: %v", err)
			}
		})
	}
}
