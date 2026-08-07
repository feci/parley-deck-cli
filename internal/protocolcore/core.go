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
	if strings.ContainsAny(version, "/\\") || strings.Contains(version, "..") {
		return Release{}, fmt.Errorf("protocolcore: unsafe version %q", version)
	}
	dir := s.ReleaseDir(version)
	path := filepath.Join(dir, CoreFileName)
	if _, err := os.Stat(path); err == nil {
		return Release{}, fmt.Errorf(
			"protocolcore: release %s already exists and releases are write-once; publish a new version instead", version)
	}
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return Release{}, err
	}
	// 0444: a release is frozen the moment it exists. This is a second layer, not the boundary
	// (D9) — it converts an accidental write into a visible error.
	if err := os.WriteFile(path, []byte(body), 0o444); err != nil {
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
