package budget

import (
	"bytes"
	"context"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"testing"
	"time"
)

func TestConcurrentBootstrapRefreshesOnlyAbsentOrigin(t *testing.T) {
	for _, mode := range []string{"published", "origin-removed", "identity-removed", "origin-conflict"} {
		t.Run(mode, func(t *testing.T) {
			s := testStore(t)
			ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
			defer cancel()
			observed, resume := make(chan struct{}), make(chan struct{})
			var once sync.Once
			unpause := func() { once.Do(func() { close(resume) }) }
			defer unpause()
			done := make(chan error, 1)
			go func() {
				release, err := lockWithBootstrapObserved(ctx, filepath.Join(s.Dir, "ledger.json"), tryLock, unlock, nil, func() {
					close(observed)
					select {
					case <-resume:
					case <-ctx.Done():
					}
				})
				if err == nil {
					release()
				}
				done <- err
			}()
			select {
			case <-observed:
			case <-ctx.Done():
				t.Fatal("bootstrap did not reach absence boundary")
			}
			if _, err := s.Reserve(ctx, Request{ID: "first", Kind: Launch, ReserveMicros: micros(3)}, Limits{CostMicros: 15}); err != nil {
				t.Fatal(err)
			}
			ledger, err := os.ReadFile(filepath.Join(s.Dir, "ledger.json"))
			if err != nil {
				t.Fatal(err)
			}
			originPath := filepath.Join(s.Dir, "lock-origin")
			origin, err := os.ReadFile(originPath)
			if err != nil {
				t.Fatal(err)
			}
			parts := strings.Split(string(origin), "\n")
			if len(parts) != 5 {
				t.Fatal("unexpected origin fixture")
			}
			identity := parts[2]
			switch mode {
			case "origin-removed":
				err = os.Remove(originPath)
			case "identity-removed":
				err = os.Remove(identity)
			case "origin-conflict":
				err = os.WriteFile(originPath, []byte("foreign origin\n"), 0600)
			}
			if err != nil {
				t.Fatal(err)
			}
			unpause()
			select {
			case err = <-done:
			case <-ctx.Done():
				t.Fatal("bootstrap did not terminate")
			}
			if mode == "published" && err != nil {
				t.Fatal("concurrent initial publication was misclassified as orphaned", err)
			}
			if mode != "published" && err == nil {
				t.Fatal("missing/conflicting established origin granted a lock")
			}
			retained, readErr := os.ReadFile(filepath.Join(s.Dir, "ledger.json"))
			if readErr != nil || !bytes.Equal(retained, ledger) {
				t.Fatal("bootstrap changed original charges", readErr)
			}
			if mode == "origin-removed" {
				if _, err := os.Lstat(originPath); !os.IsNotExist(err) {
					t.Fatal("bootstrap recreated a lost established origin")
				}
			} else if mode == "identity-removed" {
				if _, err := os.Lstat(identity); !os.IsNotExist(err) {
					t.Fatal("bootstrap recreated a lost established identity")
				}
			} else if mode == "published" {
				retained, err = os.ReadFile(originPath)
				if err != nil || !bytes.Equal(retained, origin) {
					t.Fatal("bootstrap changed the published origin", err)
				}
			}
		})
	}
}
