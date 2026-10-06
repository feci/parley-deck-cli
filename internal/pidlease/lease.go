// Package pidlease provides exclusive PID/token leases on the owning filesystem.
// It uses the driver's exclusive-file primitive rather than advisory flock,
// which is not exclusive on some shared filesystems.
package pidlease

import (
	"bytes"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"sync"

	"parley-deck-cli/internal/fsutil"
	"parley-deck-cli/internal/procctl"
)

var ErrHeld = errors.New("lease held")
var tokenPattern = regexp.MustCompile(`^[a-f0-9]{32}$`)

type Owner struct {
	Version  int    `json:"version"`
	PID      int    `json:"pid"`
	Token    string `json:"token"`
	Host     string `json:"host"`
	Boot     string `json:"boot"`
	Identity string `json:"identity"`
}

type Lease struct {
	path string
	raw  []byte
	once sync.Once
}

func identity() (string, string) {
	host, _ := os.Hostname()
	return host, procctl.CurrentBootID()
}

// ProvenDeadLocal is the shared crash-recovery predicate. Unknown/foreign host
// or boot and ambiguous liveness (including permission failures) never prove death.
func ProvenDeadLocal(host, boot string, pid int) bool {
	localHost, localBoot := identity()
	return deadLocal(host, boot, pid, localHost, localBoot)
}
func deadLocal(host, boot string, pid int, localHost, localBoot string) bool {
	return localHost != "" && localBoot != "" && host == localHost && boot == localBoot && pid > 0 && pid <= 2147483647 && pid != os.Getpid() && definitelyDead(pid)
}

// TryAcquire publishes a fully written owner by exclusive hard link. A partial,
// unreadable or foreign owner is never stale. Unknown host/boot identity permits
// ownership but disables automatic stale takeover.
func TryAcquire(path, id string) (*Lease, error) {
	host, boot := identity()
	return tryAcquire(path, id, host, boot)
}

func tryAcquire(path, id, host, boot string) (*Lease, error) {
	if id == "" {
		return nil, fmt.Errorf("missing lease identity")
	}
	if err := os.MkdirAll(filepath.Dir(path), 0700); err != nil {
		return nil, err
	}
	token := make([]byte, 16)
	if _, err := rand.Read(token); err != nil {
		return nil, err
	}
	owner := Owner{1, os.Getpid(), hex.EncodeToString(token), host, boot, id}
	raw, _ := json.Marshal(owner)
	for attempt := 0; attempt < 2; attempt++ {
		err := publish(path, raw)
		if err == nil {
			return &Lease{path: path, raw: raw}, nil
		}
		if !errors.Is(err, os.ErrExist) {
			return nil, err
		}
		old, e := os.ReadFile(path)
		if e != nil {
			return nil, fmt.Errorf("%w: unreadable owner at %s: %v", ErrHeld, path, e)
		}
		var o Owner
		parseErr := json.Unmarshal(old, &o)
		canonical, _ := json.Marshal(o)
		if parseErr != nil || !bytes.Equal(canonical, old) || o.Version != 1 || o.PID <= 0 || o.PID > 2147483647 || o.Identity == "" || !tokenPattern.MatchString(o.Token) {
			return nil, fmt.Errorf("%w: partial or invalid owner at %s; owner-visible recovery required", ErrHeld, path)
		}
		if !deadLocal(o.Host, o.Boot, o.PID, host, boot) {
			return nil, fmt.Errorf("%w: %s pid=%d host=%q boot=%q at %s; only a proven dead local owner may be reclaimed", ErrHeld, o.Identity, o.PID, o.Host, o.Boot, path)
		}
		if e = reap(path, old, o.Token); e != nil {
			return nil, e
		}
	}
	return nil, fmt.Errorf("%w: another contender won %s", ErrHeld, path)
}

// A permanent per-generation claim elects exactly one stale reaper. Never remove
// these claims: a delayed second reader of the old owner must not unlink a new
// winner. A crash after claiming is deliberately fail-closed to an owner-visible
// repair, rather than guessing the state of an interrupted reclamation.
func reap(path string, old []byte, token string) error {
	claim := path + ".reaped-" + token
	f, err := os.OpenFile(claim, os.O_WRONLY|os.O_CREATE|os.O_EXCL, 0600)
	if err != nil {
		return fmt.Errorf("%w: stale reclamation already claimed at %s: %v", ErrHeld, claim, err)
	}
	if err = fsutil.SyncFile(f); err != nil {
		f.Close()
		return err
	}
	if err = f.Close(); err != nil {
		return err
	}
	if err = syncDir(filepath.Dir(path)); err != nil {
		return err
	}
	current, err := os.ReadFile(path)
	if err != nil || !bytes.Equal(current, old) {
		return fmt.Errorf("%w: owner changed during stale reclamation", ErrHeld)
	}
	if err = os.Remove(path); err != nil {
		return err
	}
	return syncDir(filepath.Dir(path))
}

func publish(path string, raw []byte) error {
	f, err := os.CreateTemp(filepath.Dir(path), ".lease-")
	if err != nil {
		return err
	}
	defer os.Remove(f.Name())
	if _, err = f.Write(raw); err == nil {
		err = fsutil.SyncFile(f)
	}
	closeErr := f.Close()
	if err != nil {
		return err
	}
	if closeErr != nil {
		return closeErr
	}
	if err = os.Link(f.Name(), path); err != nil {
		return err
	}
	return syncDir(filepath.Dir(path))
}

func syncDir(path string) error {
	f, err := os.Open(path)
	if err != nil {
		return err
	}
	defer f.Close()
	return fsutil.SyncFile(f)
}

func (l *Lease) Check() error {
	raw, err := os.ReadFile(l.path)
	if err != nil || !bytes.Equal(raw, l.raw) {
		return fmt.Errorf("lease ownership lost at %s", l.path)
	}
	return nil
}

// Release cannot delete another generation (including another run in this PID).
// Compliant contenders never reap a live PID, so the identity check and unlink
// cannot race a legitimate stale takeover of this generation.
func (l *Lease) Release() {
	l.once.Do(func() {
		if l.Check() == nil {
			_ = os.Remove(l.path)
		}
	})
}
