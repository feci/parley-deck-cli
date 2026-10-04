package pidlease

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"
	"testing"
)

func TestLeaseProcess(t *testing.T) {
	acquire := TryAcquire
	if os.Getenv("PARLEY_PIDLEASE_SYNTHETIC_IDENTITY") == "1" {
		acquire = func(path, id string) (*Lease, error) {
			return tryAcquire(path, id, "synthetic-test-host", "synthetic-test-boot")
		}
	}
	if path := os.Getenv("PARLEY_PIDLEASE_CHILD"); path != "" {
		l, e := acquire(path, "idea/child-run")
		if e != nil {
			fmt.Println(e)
			os.Exit(23)
		}
		fmt.Println("OWNED")
		var b [1]byte
		os.Stdin.Read(b[:])
		l.Release()
		os.Exit(0)
	}
	path := filepath.Join(t.TempDir(), "driver.lease")
	cmd := exec.Command(os.Args[0], "-test.run=^TestLeaseProcess$")
	cmd.Env = append(os.Environ(), "PARLEY_PIDLEASE_CHILD="+path)
	in, _ := cmd.StdinPipe()
	out, _ := cmd.StdoutPipe()
	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	if e := cmd.Start(); e != nil {
		t.Fatal(e)
	}
	defer func() { in.Close(); cmd.Process.Kill(); cmd.Wait() }()
	data := make([]byte, 6)
	if n, e := out.Read(data); e != nil || n != 6 || string(data) != "OWNED\n" {
		t.Fatalf("child: %d %v %s %s", n, e, data, stderr.String())
	}
	if l, e := acquire(path, "idea/parent-run"); e == nil {
		l.Release()
		t.Fatal("two processes own same lease")
	}
	if e := cmd.Process.Kill(); e != nil {
		t.Fatal(e)
	}
	cmd.Wait()
	l, e := acquire(path, "idea/takeover")
	if e != nil {
		t.Fatal("stale takeover", e)
	}
	defer l.Release()
	if e = l.Check(); e != nil {
		t.Fatal(e)
	}
	t.Logf("%s: different-process refusal, proven-dead local takeover, ownership intact", path)
}
func TestLeaseProcessSyntheticIdentity(t *testing.T) {
	t.Setenv("PARLEY_PIDLEASE_SYNTHETIC_IDENTITY", "1")
	TestLeaseProcess(t)
}
func TestLeaseSamePIDPartialForeignAndReleaseIdentity(t *testing.T) {
	path := filepath.Join(t.TempDir(), "lease")
	l, e := TryAcquire(path, "idea/run-a")
	if e != nil {
		t.Fatal(e)
	}
	old := append([]byte(nil), l.raw...)
	if rival, e := TryAcquire(path, "idea/run-b"); e == nil {
		rival.Release()
		t.Fatal("same PID competing run admitted")
	}
	os.Remove(path)
	winner, e := TryAcquire(path, "idea/run-c")
	if e != nil {
		t.Fatal(e)
	}
	l.Release()
	if e = winner.Check(); e != nil {
		t.Fatal("old release removed winner", e)
	}
	winner.Release()
	for _, kind := range []string{"empty", "partial", "foreign-host", "foreign-boot", "unknown-host", "unknown-boot"} {
		t.Run(kind, func(t *testing.T) {
			var o Owner
			json.Unmarshal(old, &o)
			o.PID = 99999999
			raw := []byte{}
			switch kind {
			case "partial":
				raw = []byte(`{"version":1`)
			case "foreign-host":
				o.Host += "-foreign"
			case "foreign-boot":
				o.Boot += "-foreign"
			case "unknown-host":
				o.Host = ""
			case "unknown-boot":
				o.Boot = ""
			}
			if strings.Contains(kind, "host") || strings.Contains(kind, "boot") {
				raw, _ = json.Marshal(o)
			}
			os.WriteFile(path, raw, 0600)
			defer os.Remove(path)
			if got, e := TryAcquire(path, "idea/new"); e == nil {
				got.Release()
				t.Fatal("unattributed owner reclaimed")
			}
			after, _ := os.ReadFile(path)
			if !bytes.Equal(after, raw) {
				t.Fatal("owner changed")
			}
		})
	}
}
func TestLeaseStaleReaperGenerationClaim(t *testing.T) {
	path := filepath.Join(t.TempDir(), "lease")
	l, e := TryAcquire(path, "idea/dead")
	if e != nil {
		t.Fatal(e)
	}
	var o Owner
	json.Unmarshal(l.raw, &o)
	old := l.raw
	// Reclamation itself is exercised with two delayed readers of the same bytes.
	// Production calls this only after proving local liveness; no live process is
	// signalled or killed by this race test.
	var wg sync.WaitGroup
	wg.Add(2)
	start := make(chan struct{})
	results := make(chan error, 2)
	for i := 0; i < 2; i++ {
		go func() { defer wg.Done(); <-start; results <- reap(path, old, o.Token) }()
	}
	close(start)
	wg.Wait()
	close(results)
	successes := 0
	for e := range results {
		if e == nil {
			successes++
		}
	}
	if successes != 1 {
		t.Fatalf("reapers succeeded=%d", successes)
	}
	winner, e := TryAcquire(path, "idea/winner")
	if e != nil {
		t.Fatal(e)
	}
	defer winner.Release()
	if e = reap(path, old, o.Token); e == nil {
		t.Fatal("delayed stale reader reaped twice")
	}
	if e = winner.Check(); e != nil {
		t.Fatal("delayed reaper removed winner", e)
	}
}
