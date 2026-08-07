// Package protocolcore implements the global core protocol store and the renderer that
// materializes a deck's COOPERATION.md from it.
//
// The model, from FINAL.md of idea meta-protocol-change-global-core-protocol: the deck's
// COOPERATION.md stops being a hand-edited store and becomes a GENERATED VIEW, exactly as §2's
// roster table did. Authority is an immutable, versioned, content-addressed core under
// ~/.parley/protocol/core/<version>/.
//
// Releases are WRITE-ONCE. That is what makes "an agent may not change the global protocol"
// structural rather than merely a rule: there is no "the core file" to edit, only releases, so a
// change by the user becomes a new version by construction.
package protocolcore

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

// CoreFileName is the protocol body inside a release directory.
const CoreFileName = "COOPERATION.md"

// ErrNoRelease reports that the requested core release is not installed.
//
// Callers must distinguish this from a read error: a missing release blocks ADOPTION and
// RENDERING (D8), but must never block continuation of an already-pinned idea, which reads its
// own materialized snapshot (D7).
var ErrNoRelease = errors.New("protocolcore: release not installed")

// Store is the on-disk release store.
type Store struct{ Root string }

// StoreAt returns the store under a Parley home directory (~/.parley, or $PARLEY_HOME).
func StoreAt(parleyHome string) Store {
	return Store{Root: filepath.Join(parleyHome, "protocol", "core")}
}

// ValidVersion reports whether a version string is safe to use as a path element.
//
// The version comes from a COMMITTED deck lock, i.e. from a file any contributor can edit, and it
// is joined onto the store root. Without this, `core-version: ../../../etc` reads and (through
// Publish) writes wherever it likes.
func ValidVersion(v string) bool {
	// Deliberately NOT trimmed: validating a trimmed copy while the untrimmed string is what gets
	// joined onto a path means the check and the use disagree.
	if v != strings.TrimSpace(v) {
		return false
	}
	if v == "" || v == "." || v == ".." {
		return false
	}
	if strings.ContainsAny(v, "/\\") || strings.Contains(v, "..") {
		return false
	}
	if strings.HasPrefix(v, ".") {
		return false
	}
	for _, r := range v {
		switch {
		case r >= 'a' && r <= 'z', r >= 'A' && r <= 'Z', r >= '0' && r <= '9':
		case r == '.' || r == '-' || r == '_' || r == '+':
		default:
			return false
		}
	}
	return true
}

// ReleaseDir is the directory holding one version's frozen bytes.
func (s Store) ReleaseDir(version string) string { return filepath.Join(s.Root, version) }

// Release is one immutable core version.
type Release struct {
	Version string
	Body    string
	SHA256  string
}

// Load reads a release. It deliberately verifies nothing beyond existence: only the caller knows
// which hash it expected, so the byte-comparison against a deck lock (G8) belongs there.
func (s Store) Load(version string) (Release, error) {
	if !ValidVersion(version) {
		return Release{}, fmt.Errorf("protocolcore: unsafe version %q", version)
	}
	path := filepath.Join(s.ReleaseDir(version), CoreFileName)
	b, err := os.ReadFile(path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return Release{}, fmt.Errorf("%w: %s", ErrNoRelease, version)
		}
		return Release{}, err
	}
	return Release{Version: version, Body: string(b), SHA256: Hash(string(b))}, nil
}

// Versions lists installed releases in byte order. Versions are not parsed or compared
// semantically here — the deck lock names an exact version, and guessing "the newest" is exactly
// the substitution D8 forbids.
func (s Store) Versions() ([]string, error) {
	entries, err := os.ReadDir(s.Root)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return nil, nil
		}
		return nil, err
	}
	out := make([]string, 0, len(entries))
	for _, e := range entries {
		if e.IsDir() {
			out = append(out, e.Name())
		}
	}
	sort.Strings(out)
	return out, nil
}

// Publish writes a NEW release.
//
// It refuses to modify an existing one. That refusal is the write-once guarantee (D1) and the
// structural half of gate G2: this is the only function in the package that writes into the core
// store, and it is reachable only from the attended publisher.
func (s Store) Publish(version, body string) (Release, error) {
	if strings.TrimSpace(version) == "" {
		return Release{}, errors.New("protocolcore: empty version")
	}
	if !ValidVersion(version) {
		return Release{}, fmt.Errorf("protocolcore: unsafe version %q", version)
	}
	dir := s.ReleaseDir(version)
	path := filepath.Join(dir, CoreFileName)
	// No component of the store path may be a symlink. Hardening only the final create still lets
	// an intermediate link redirect the whole store elsewhere.
	if err := assertNoSymlinkComponents(s.Root); err != nil {
		return Release{}, err
	}
	// Write-once is per RELEASE, not per file: a directory that already exists is an existing
	// release namespace even when its body is missing, and populating it would let a second
	// publish complete a half-written one. Refuse.
	if _, err := os.Lstat(dir); err == nil {
		return Release{}, fmt.Errorf(
			"protocolcore: release %s already exists and releases are write-once; publish a new version instead", version)
	} else if !errors.Is(err, os.ErrNotExist) {
		return Release{}, err
	}
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return Release{}, err
	}
	// O_EXCL|O_NOFOLLOW: never write THROUGH a symlink, and never overwrite. os.WriteFile follows
	// symlinks, so a planted link at the release path would have redirected a "write-once" release
	// to an arbitrary file — the guarantee inverted into an arbitrary-write primitive.
	// 0444 freezes the release the moment it exists; a second layer, not the boundary (D9).
	f, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_EXCL|noFollow, 0o444)
	if err != nil {
		// A failed publish must not brick the version label: the directory we just created would
		// otherwise make every later attempt at this version report "already exists".
		os.Remove(dir)
		return Release{}, err
	}
	if _, err := f.WriteString(body); err != nil {
		f.Close()
		os.Remove(path)
		os.Remove(dir)
		return Release{}, err
	}
	if err := f.Close(); err != nil {
		os.Remove(path)
		os.Remove(dir)
		return Release{}, err
	}
	return Release{Version: version, Body: body, SHA256: Hash(body)}, nil
}

// Hash is the content address used for core bytes, overlay bytes and the effective render.
func Hash(s string) string {
	sum := sha256.Sum256([]byte(s))
	return hex.EncodeToString(sum[:])
}

// ShortHash is the display form.
func ShortHash(h string) string {
	if len(h) <= 12 {
		return h
	}
	return h[:12]
}

// assertNoSymlinkComponents refuses if any component OF THE STORE ITSELF is a symlink.
//
// O_NOFOLLOW protects only the file being opened. Someone who can plant a link at
// ~/.parley/protocol or ~/.parley/protocol/core redirects the whole release store, and every
// "write-once" guarantee then applies somewhere else entirely.
//
// Scope matters: only the components the store owns are checked, not the whole path to the
// filesystem root. On macOS /var is itself a symlink to /private/var, so walking to the root
// rejects every temp directory on the machine — the first version of this check did exactly that
// and failed eight of its own tests.
func assertNoSymlinkComponents(root string) error {
	abs, err := filepath.Abs(root)
	if err != nil {
		return err
	}
	// root is <parley-home>/protocol/core; check core, protocol, and the home itself.
	p := abs
	for i := 0; i < 3; i++ {
		fi, err := os.Lstat(p)
		if err == nil && fi.Mode()&os.ModeSymlink != 0 {
			return fmt.Errorf("protocolcore: %s is a symlink; the core store must not be reached through one", p)
		}
		parent := filepath.Dir(p)
		if parent == p {
			break
		}
		p = parent
	}
	return nil
}
