package evidence

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

func refusalFixture() VerificationRefusal {
	return VerificationRefusal{Version: 1, ObservationID: strings.Repeat("a", 64), VerificationID: strings.Repeat("b", 64), ObservedAt: time.Date(2026, 9, 11, 0, 0, 0, 0, time.UTC), Observer: "helper", Stage: "execution", Executions: []RefusalExecution{{Status: StatusSkipped}}}
}

func TestRefusalReadOnlyInspectionAndImmutableRecovery(t *testing.T) {
	dir := t.TempDir()
	entries, err := InspectVerificationRefusals(dir)
	if err != nil || len(entries) != 0 {
		t.Fatalf("pristine: %+v %v", entries, err)
	}
	files, _ := os.ReadDir(dir)
	if len(files) != 0 {
		t.Fatal("inspection created state")
	}
	r := refusalFixture()
	sum, err := RetainVerificationRefusal(dir, r)
	if err != nil {
		t.Fatal(err)
	}
	// Interrupted commit retains both the pending and canonical exact bytes.
	sentinel := errors.New("commit failed")
	if err := PublishVerificationRefusal(context.Background(), dir, sum, func(string, []byte) error { return sentinel }); !errors.Is(err, sentinel) {
		t.Fatal(err)
	}
	if _, err := RetainVerificationRefusal(dir, r); err != nil {
		t.Fatal(err)
	}
	entries, err = InspectVerificationRefusals(dir)
	if err != nil || len(entries) != 1 || !entries[0].Canonical || len(entries[0].PendingFiles) != 2 {
		t.Fatalf("lost or duplicated observation: %+v %v", entries, err)
	}
	original, _ := encodeRefusal(r)
	commits := 0
	commit := func(path string, data []byte) error {
		commits++
		if filepath.Base(path) != sum+".json" || !bytes.Equal(data, original) {
			t.Fatal("wrong recovery bytes")
		}
		return nil
	}
	if err := PublishVerificationRefusal(context.Background(), dir, sum, commit); err != nil {
		t.Fatal(err)
	}
	if commits != 1 {
		t.Fatal("commit was not verified")
	}
	path := filepath.Join(dir, VerificationRefusalDirectory, sum+".json")
	if err := os.WriteFile(path, []byte("different"), 0600); err != nil {
		t.Fatal(err)
	}
	if err := PublishVerificationRefusal(context.Background(), dir, sum, commit); err == nil {
		t.Fatal("conflicting canonical bytes overwritten")
	}
	actual, _ := os.ReadFile(path)
	if string(actual) != "different" {
		t.Fatal("conflict was destroyed")
	}
}

func TestRefusalStrictSafeSchema(t *testing.T) {
	raw, _ := encodeRefusal(refusalFixture())
	for _, data := range [][]byte{
		bytes.Replace(raw, []byte(`"version": 1`), []byte(`"version": 1, "version": 1`), 1),
		bytes.Replace(raw, []byte(`"version": 1`), []byte(`"Version": 1`), 1),
		bytes.Replace(raw, []byte(`"idea": null,`), nil, 1),
		append(append([]byte{}, raw...), []byte("{}")...),
		bytes.Replace(raw, []byte(`"idea": null`), []byte(`"idea": "Bearer-a-secret"`), 1),
		bytes.Replace(raw, []byte(`"stage": "execution"`), []byte(`"stage": "raw-provider-error"`), 1),
	} {
		if _, err := decodeRefusal(data); err == nil {
			t.Fatalf("unsafe/ambiguous schema accepted: %s", data)
		}
	}
}

func TestRefusalIncompleteEntriesAreVisibleAndDoNotEraseValidRecovery(t *testing.T) {
	dir := t.TempDir()
	sum, err := RetainVerificationRefusal(dir, refusalFixture())
	if err != nil {
		t.Fatal(err)
	}
	pending, _, _ := refusalDirs(dir)
	if err := os.WriteFile(filepath.Join(pending, "observation-123.json"), []byte(`{"version":`), 0600); err != nil {
		t.Fatal(err)
	}
	entries, err := InspectVerificationRefusals(dir)
	if err != nil || len(entries) != 2 || entries[1].Problem == "" {
		t.Fatalf("lost interrupted write: %+v %v", entries, err)
	}
	if err := PublishVerificationRefusal(context.Background(), dir, sum, func(string, []byte) error { return nil }); err != nil {
		t.Fatal(err)
	}
	entries, _ = InspectVerificationRefusals(dir)
	if len(entries) != 2 || entries[1].Problem == "" {
		t.Fatal("recovery erased incomplete observation")
	}
}

func TestRefusalConcurrentObserversAndExactRecovery(t *testing.T) {
	dir := t.TempDir()
	const n = 8
	sums := make(chan string, n)
	errs := make(chan error, n)
	var group sync.WaitGroup
	for i := 0; i < n; i++ {
		group.Add(1)
		go func(i int) {
			defer group.Done()
			r := refusalFixture()
			r.ObservationID = fmt.Sprintf("%064x", i)
			sum, err := RetainVerificationRefusal(dir, r)
			if err == nil {
				sums <- sum
			}
			errs <- err
		}(i)
	}
	group.Wait()
	close(errs)
	close(sums)
	for err := range errs {
		if err != nil {
			t.Fatal(err)
		}
	}
	entries, err := InspectVerificationRefusals(dir)
	if err != nil || len(entries) != n {
		t.Fatalf("lost concurrent observers: %d %v", len(entries), err)
	}
	first := <-sums
	var calls, active atomic.Int32
	for i := 0; i < n; i++ {
		group.Add(1)
		go func() {
			defer group.Done()
			err := PublishVerificationRefusal(context.Background(), dir, first, func(string, []byte) error {
				// The race detector cannot infer happens-before from the kernel lock.
				// Atomic instrumentation also checks that callbacks never overlap.
				calls.Add(1)
				if active.Add(1) != 1 {
					t.Error("guard admitted overlapping publication callbacks")
				}
				defer active.Add(-1)
				time.Sleep(5 * time.Millisecond)
				return nil
			})
			if err != nil {
				t.Error(err)
			}
		}()
	}
	group.Wait()
	if calls.Load() != n {
		t.Fatal("not every exact replay verified publication")
	}
	entries, err = InspectVerificationRefusals(dir)
	if err != nil || len(entries) != n {
		t.Fatal("replay lost later observations")
	}
}

func TestRefusalStorageAliasesAndChangedRecordRefuse(t *testing.T) {
	for _, which := range []string{"pending", "canonical"} {
		t.Run(which, func(t *testing.T) {
			dir := t.TempDir()
			pending, canonical, err := refusalDirs(dir)
			if err != nil {
				t.Fatal(err)
			}
			path := canonical
			if which == "pending" {
				path = pending
			}
			if err := os.MkdirAll(filepath.Dir(path), 0700); err != nil {
				t.Fatal(err)
			}
			if err := os.Symlink(t.TempDir(), path); err != nil {
				t.Fatal(err)
			}
			if _, err := InspectVerificationRefusals(dir); err == nil {
				t.Fatal("aliased storage read")
			}
			if which == "pending" {
				if _, err := RetainVerificationRefusal(dir, refusalFixture()); err == nil {
					t.Fatal("aliased pending storage written")
				}
			}
		})
	}
	dir := t.TempDir()
	sum, err := RetainVerificationRefusal(dir, refusalFixture())
	if err != nil {
		t.Fatal(err)
	}
	pending, _, _ := refusalDirs(dir)
	files, _ := os.ReadDir(pending)
	r := refusalFixture()
	r.Stage = "publication"
	changed, _ := json.MarshalIndent(r, "", "  ")
	if err := os.WriteFile(filepath.Join(pending, files[0].Name()), changed, 0600); err != nil {
		t.Fatal(err)
	}
	if err := PublishVerificationRefusal(context.Background(), dir, sum, func(string, []byte) error { return nil }); err == nil {
		t.Fatal("changed observation recovered under old digest")
	}
}

func TestRefusalRetentionSurvivesMissingGuardOriginLock(t *testing.T) {
	dir, err := filepath.EvalSymlinks(t.TempDir())
	if err != nil {
		t.Fatal(err)
	}
	if err := WithReportWriter(context.Background(), dir, func(*ReportWriter) error { return nil }); err != nil {
		t.Fatal(err)
	}
	guard, err := reportGuardDirectory(dir)
	if err != nil {
		t.Fatal(err)
	}
	origin, err := os.ReadFile(filepath.Join(guard, "lock-origin"))
	if err != nil {
		t.Fatal(err)
	}
	parts := strings.Split(string(origin), "\n")
	if len(parts) != 5 {
		t.Fatal("unexpected fixture lock origin")
	}
	// No writer holds this test-only inode. Simulate loss without recreating it.
	if err := os.Remove(parts[2]); err != nil {
		t.Fatal(err)
	}
	sum, err := RetainVerificationRefusal(dir, refusalFixture())
	if err != nil {
		t.Fatal(err)
	}
	if err := PublishVerificationRefusal(context.Background(), dir, sum, func(string, []byte) error { t.Fatal("missing lock was bypassed"); return nil }); err == nil {
		t.Fatal("missing lock allowed publication")
	}
	entries, err := InspectVerificationRefusals(dir)
	if err != nil || len(entries) != 1 || entries[0].Canonical {
		t.Fatalf("guard failure lost pending observation: %+v %v", entries, err)
	}
	if _, err := os.Stat(parts[2]); !os.IsNotExist(err) {
		t.Fatal("recovery recreated missing lock")
	}
}

func TestRefusalConflictingObservationIdentityCannotRecover(t *testing.T) {
	dir := t.TempDir()
	r := refusalFixture()
	first, err := RetainVerificationRefusal(dir, r)
	if err != nil {
		t.Fatal(err)
	}
	r.Stage = "publication"
	second, err := RetainVerificationRefusal(dir, r)
	if err != nil {
		t.Fatal(err)
	}
	entries, err := InspectVerificationRefusals(dir)
	if err != nil || len(entries) != 2 {
		t.Fatal(err)
	}
	for _, entry := range entries {
		if entry.Problem != "conflicting observation identity" {
			t.Fatal("conflicting identity hidden")
		}
	}
	for _, sum := range []string{first, second} {
		if err := PublishVerificationRefusal(context.Background(), dir, sum, func(string, []byte) error { return nil }); err == nil {
			t.Fatal("conflicting identity recovered")
		}
	}
}

func TestRefusalRecoversOnlyMatchingInterruptedCanonicalStaging(t *testing.T) {
	for _, conflict := range []bool{false, true} {
		t.Run(fmt.Sprint(conflict), func(t *testing.T) {
			dir := t.TempDir()
			sum, err := RetainVerificationRefusal(dir, refusalFixture())
			if err != nil {
				t.Fatal(err)
			}
			_, canonical, _ := refusalDirs(dir)
			if err := os.MkdirAll(canonical, 0700); err != nil {
				t.Fatal(err)
			}
			data, _ := encodeRefusal(refusalFixture())
			staged := data[:len(data)/2]
			if conflict {
				staged = []byte("divergent abandoned write")
			}
			path := filepath.Join(canonical, ".refusal-"+sum+"-123.tmp")
			if err := os.WriteFile(path, staged, 0600); err != nil {
				t.Fatal(err)
			}
			err = PublishVerificationRefusal(context.Background(), dir, sum, func(string, []byte) error { return nil })
			if conflict {
				if err == nil {
					t.Fatal("divergent staging silently recovered")
				}
				actual, e := os.ReadFile(path)
				if e != nil || !bytes.Equal(actual, staged) {
					t.Fatal("divergent staging erased")
				}
			} else {
				if err != nil {
					t.Fatal(err)
				}
				if _, err := os.Stat(path); !os.IsNotExist(err) {
					t.Fatal("attributable stale staging still blocks recovery")
				}
				entries, err := InspectVerificationRefusals(dir)
				if err != nil || len(entries) != 1 || entries[0].Problem != "" {
					t.Fatalf("staging poisoned valid recovery: %+v %v", entries, err)
				}
			}
		})
	}
}
