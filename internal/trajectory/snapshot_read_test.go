package trajectory

import (
	"archive/tar"
	"bytes"
	"context"
	"crypto/sha256"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

func snapshotOpenRoot(t *testing.T, dir string) *os.Root {
	t.Helper()
	root, err := os.OpenRoot(dir)
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { root.Close() })
	return root
}

func snapshotStat(t *testing.T, name string) os.FileInfo {
	t.Helper()
	info, err := os.Lstat(name)
	if err != nil {
		t.Fatal(err)
	}
	return info
}

func snapshotMemberBytes(t *testing.T, archive []byte, want []byte, mode os.FileMode) {
	t.Helper()
	reader := tar.NewReader(bytes.NewReader(archive))
	header, err := reader.Next()
	if err != nil {
		t.Fatal(err)
	}
	data, err := io.ReadAll(reader)
	if err != nil {
		t.Fatal(err)
	}
	if header.Name != "files/source" || header.Typeflag != tar.TypeReg || header.Mode != int64(mode.Perm()) || header.Size != int64(len(want)) || !bytes.Equal(data, want) {
		t.Fatalf("captured member differs: header=%+v data=%q want=%q mode=%v", header, data, want, mode)
	}
	if _, err := reader.Next(); err != io.EOF {
		t.Fatalf("unexpected additional archive member: %v", err)
	}
}

// The stale metadata comes from a real earlier stat of the same inode. No
// fabricated FileInfo or race timing is needed to exercise this boundary.
func TestSnapshotReadUsesDescriptorAfterStaleMetadata(t *testing.T) {
	dir := t.TempDir()
	snapshotWrite(t, dir, "source", []byte("old"), 0600)
	name := filepath.Join(dir, "source")
	prior := snapshotStat(t, name)
	want := []byte("complete replacement contents, longer than the preceding stat\n")
	snapshotWrite(t, dir, "source", want, 0700)
	stamp := prior.ModTime().Add(2 * time.Hour)
	if err := os.Chtimes(name, stamp, stamp); err != nil {
		t.Fatal(err)
	}
	actual := snapshotStat(t, name)
	if !os.SameFile(prior, actual) || prior.Size() == actual.Size() || prior.ModTime().Equal(actual.ModTime()) {
		t.Fatal("fixture did not retain the inode with changed size and mtime")
	}
	var out bytes.Buffer
	w := tar.NewWriter(&out)
	if err := captureSnapshotRegular(context.Background(), snapshotOpenRoot(t, dir), w, "source", prior); err != nil {
		t.Fatal(err)
	}
	if err := w.Close(); err != nil {
		t.Fatal(err)
	}
	snapshotMemberBytes(t, out.Bytes(), want, actual.Mode())
}

func TestSnapshotReadRefusesReplacedRoot(t *testing.T) {
	parent := t.TempDir()
	dir := filepath.Join(parent, "root")
	snapshotWrite(t, dir, "source", []byte("original"), 0600)
	origin := snapshotStat(t, dir)
	if err := os.Rename(dir, filepath.Join(parent, "retained-root")); err != nil {
		t.Fatal(err)
	}
	snapshotWrite(t, dir, "source", []byte("replacement"), 0600)
	var out bytes.Buffer
	err := captureSnapshotMember(context.Background(), dir, origin, tar.NewWriter(&out), "source")
	if err == nil || !strings.Contains(err.Error(), "source root changed") || out.Len() != 0 {
		t.Fatalf("replacement root was not refused before writing: %v, bytes=%d", err, out.Len())
	}
}

func TestSnapshotReadRefusesReplacedFile(t *testing.T) {
	dir := t.TempDir()
	snapshotWrite(t, dir, "source", []byte("original"), 0600)
	name := filepath.Join(dir, "source")
	prior := snapshotStat(t, name)
	// Retain the old inode so an immediate reuse cannot make the fixture pass.
	if err := os.Rename(name, filepath.Join(dir, "retained-source")); err != nil {
		t.Fatal(err)
	}
	snapshotWrite(t, dir, "source", []byte("replacement"), 0600)
	var out bytes.Buffer
	err := captureSnapshotRegular(context.Background(), snapshotOpenRoot(t, dir), tar.NewWriter(&out), "source", prior)
	if err == nil || !strings.Contains(err.Error(), "changed during open") || out.Len() != 0 {
		t.Fatalf("replacement inode was not refused before writing: %v, bytes=%d", err, out.Len())
	}
}

func TestSnapshotReadBoundsDescriptorAfterSmallStat(t *testing.T) {
	dir := t.TempDir()
	snapshotWrite(t, dir, "source", []byte("small"), 0600)
	name := filepath.Join(dir, "source")
	prior := snapshotStat(t, name)
	if err := os.Truncate(name, MaxSnapshotFileBytes+1); err != nil {
		t.Fatal(err)
	}
	if !os.SameFile(prior, snapshotStat(t, name)) {
		t.Fatal("fixture replaced the source inode")
	}
	var out bytes.Buffer
	err := captureSnapshotRegular(context.Background(), snapshotOpenRoot(t, dir), tar.NewWriter(&out), "source", prior)
	if err == nil || !strings.Contains(err.Error(), "exceeds snapshot bound") || out.Len() != 0 {
		t.Fatalf("descriptor size bound was not enforced before writing: %v, bytes=%d", err, out.Len())
	}
}

type snapshotOnHeaderWriter struct {
	bytes.Buffer
	once       func()
	skipWrites int
}

func (w *snapshotOnHeaderWriter) Write(data []byte) (int, error) {
	if len(data) > 0 && w.once != nil {
		if w.skipWrites > 0 {
			w.skipWrites--
		} else {
			f := w.once
			w.once = nil
			f()
		}
	}
	return w.Buffer.Write(data)
}

func TestSnapshotReadRefusesMutationAfterDescriptorStat(t *testing.T) {
	for _, kind := range []string{"size", "content"} {
		t.Run(kind, func(t *testing.T) {
			dir := t.TempDir()
			snapshotWrite(t, dir, "source", []byte("original"), 0600)
			name := filepath.Join(dir, "source")
			prior := snapshotStat(t, name)
			w := &snapshotOnHeaderWriter{once: func() {
				// A header is emitted after the descriptor stat and before
				// copying file bytes, making the mutation deterministic.
				if kind == "size" {
					if err := os.Truncate(name, prior.Size()+1); err != nil {
						t.Fatal(err)
					}
				} else {
					if err := os.WriteFile(name, []byte("replaced"), 0600); err != nil {
						t.Fatal(err)
					}
					stamp := prior.ModTime().Add(2 * time.Hour)
					if err := os.Chtimes(name, stamp, stamp); err != nil {
						t.Fatal(err)
					}
				}
			}}
			if kind == "content" {
				// Change the source only after its original bytes have been
				// read and are being written to the tar. Revalidation must
				// compare against those original bytes, not today's source.
				w.skipWrites = 1
			}
			err := captureSnapshotRegular(context.Background(), snapshotOpenRoot(t, dir), tar.NewWriter(w), "source", prior)
			if w.once != nil || err == nil {
				t.Fatalf("post-open %s change escaped stability check: %v", kind, err)
			}
		})
	}
}

// Running this compiled fixture with TMPDIR on AppleVirtIOFS exercises the
// production member reader on the affected mount, retaining the original root
// handle while rewriting the same inode. Every observation must be exact.
func TestSnapshotReadRefreshesRootAcrossWrites(t *testing.T) {
	dir := t.TempDir()
	snapshotWrite(t, dir, "source", []byte("original"), 0600)
	held := snapshotOpenRoot(t, dir)
	origin, err := held.Stat(".")
	if err != nil {
		t.Fatal(err)
	}
	if _, err := held.Lstat("source"); err != nil {
		t.Fatal(err)
	}
	for i := 0; i < 128; i++ {
		want := []byte(fmt.Sprintf("write %03d: %s\n", i, strings.Repeat("x", i%37)))
		snapshotWrite(t, dir, "source", want, 0600)
		var out bytes.Buffer
		w := tar.NewWriter(&out)
		if err := captureSnapshotMember(context.Background(), dir, origin, w, "source"); err != nil {
			t.Fatalf("write %d: %v", i, err)
		}
		if err := w.Close(); err != nil {
			t.Fatal(err)
		}
		snapshotMemberBytes(t, out.Bytes(), want, 0600)
	}
}

func TestSnapshotReadRetainsExactBytesAcrossDelayedMetadata(t *testing.T) {
	dir := t.TempDir()
	want := []byte("stable source with no subsequent write\n")
	// Match an ordinary content write without a subsequent chmod/stat setter.
	if err := os.WriteFile(filepath.Join(dir, "source"), want, 0600); err != nil {
		t.Fatal(err)
	}
	prior := snapshotStat(t, filepath.Join(dir, "source"))
	w := &snapshotOnHeaderWriter{once: func() {
		// A delayed sink crosses the observed AppleVirtIOFS attribute-cache
		// refresh interval after descriptor stat, without changing source.
		time.Sleep(1200 * time.Millisecond)
	}}
	tw := tar.NewWriter(w)
	if err := captureSnapshotRegular(context.Background(), snapshotOpenRoot(t, dir), tw, "source", prior); err != nil {
		t.Fatal(err)
	}
	if err := tw.Close(); err != nil {
		t.Fatal(err)
	}
	snapshotMemberBytes(t, w.Bytes(), want, 0600)
}

func TestSnapshotRevalidationRejectsDifferentMaterial(t *testing.T) {
	for _, kind := range []string{"bytes", "root", "inode", "symlink", "size", "mode", "oversize", "cancel"} {
		t.Run(kind, func(t *testing.T) {
			parent := t.TempDir()
			dir := filepath.Join(parent, "root")
			want := []byte("original")
			snapshotWrite(t, dir, "source", want, 0600)
			name := filepath.Join(dir, "source")
			prior := snapshotStat(t, name)
			root := snapshotOpenRoot(t, dir)
			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()
			switch kind {
			case "bytes":
				snapshotWrite(t, dir, "source", []byte("replaced"), 0600)
			case "root":
				if err := os.Rename(dir, filepath.Join(parent, "old-root")); err != nil {
					t.Fatal(err)
				}
				snapshotWrite(t, dir, "source", want, 0600)
			case "inode", "symlink":
				if err := os.Rename(name, filepath.Join(dir, "old-source")); err != nil {
					t.Fatal(err)
				}
				if kind == "symlink" {
					if err := os.Symlink("old-source", name); err != nil {
						t.Skipf("symlinks unavailable: %v", err)
					}
				} else {
					snapshotWrite(t, dir, "source", want, 0600)
				}
			case "size":
				if err := os.Truncate(name, prior.Size()-1); err != nil {
					t.Fatal(err)
				}
			case "mode":
				if err := os.Chmod(name, 0400); err != nil {
					t.Fatal(err)
				}
			case "oversize":
				if err := os.Truncate(name, MaxSnapshotFileBytes+1); err != nil {
					t.Fatal(err)
				}
				prior = snapshotStat(t, name)
			case "cancel":
				cancel()
			}
			hash := sha256.Sum256(want)
			err := verifySnapshotRegular(ctx, root, "source", prior, hash[:])
			if err == nil {
				t.Fatal("different material or canceled verification was accepted")
			}
			if kind == "cancel" && !errors.Is(err, context.Canceled) {
				t.Fatal("cancellation cause was lost", err)
			}
		})
	}
}

func TestSnapshotRevalidationRetainsIdenticalMaterialAfterTimestampChange(t *testing.T) {
	dir := t.TempDir()
	want := []byte("original")
	snapshotWrite(t, dir, "source", want, 0600)
	name := filepath.Join(dir, "source")
	prior := snapshotStat(t, name)
	w := &snapshotOnHeaderWriter{skipWrites: 1, once: func() {
		stamp := prior.ModTime().Add(2 * time.Hour)
		if err := os.Chtimes(name, stamp, stamp); err != nil {
			t.Fatal(err)
		}
	}}
	tw := tar.NewWriter(w)
	if err := captureSnapshotRegular(context.Background(), snapshotOpenRoot(t, dir), tw, "source", prior); err != nil {
		t.Fatal(err)
	}
	if err := tw.Close(); err != nil {
		t.Fatal(err)
	}
	if w.once != nil {
		t.Fatal("timestamp transition was not exercised")
	}
	snapshotMemberBytes(t, w.Bytes(), want, 0600)
}

func TestSnapshotRevalidationCannotAdmitDifferentExpectedTree(t *testing.T) {
	dir := t.TempDir()
	gitFixture(t, dir, "init", "-q")
	snapshotWrite(t, dir, "source", []byte("original\n"), 0600)
	gitFixture(t, dir, "add", ".")
	gitFixture(t, dir, "commit", "-qm", "Original source")
	want, err := Observe(context.Background(), dir)
	if err != nil {
		t.Fatal(err)
	}
	name := filepath.Join(dir, "source")
	prior := snapshotStat(t, name)
	w := &snapshotOnHeaderWriter{}
	tw := tar.NewWriter(w)
	meta, err := canonical(snapshotMetadata{Version: 1, Source: want})
	if err != nil {
		t.Fatal(err)
	}
	if err := tw.WriteHeader(snapshotHeader("snapshot.json", tar.TypeReg, 0600, int64(len(meta)), "")); err != nil {
		t.Fatal(err)
	}
	if _, err := tw.Write(meta); err != nil {
		t.Fatal(err)
	}
	snapshotWrite(t, dir, "source", []byte("replaced\n"), 0600)
	// Both reads agree on the same new bytes. That only establishes a stable
	// member; it cannot authorize replacing the original expected whole tree.
	if err := captureSnapshotRegular(context.Background(), snapshotOpenRoot(t, dir), tw, "source", prior); err != nil {
		t.Fatal(err)
	}
	if err := tw.Close(); err != nil {
		t.Fatal(err)
	}
	reader := tar.NewReader(bytes.NewReader(w.Bytes()))
	if _, err := reader.Next(); err != nil {
		t.Fatal(err)
	}
	if _, err := reader.Next(); err != nil {
		t.Fatal(err)
	}
	data, err := io.ReadAll(reader)
	if err != nil || string(data) != "replaced\n" {
		t.Fatalf("fixture did not archive the changed material: %q %v", data, err)
	}
	store := snapshotStoreFixture(t)
	ref := archiveFixtureStore(t, store, w.Bytes())
	if err := CheckSnapshot(context.Background(), store, ref, want); err == nil {
		t.Fatal("stable reread overrode the original expected source tree")
	}
}

// Err is checked on entry and then by contextReader immediately before reading
// data. Interpose at that read boundary to change a real file without a timing
// race or a production-only hook.
type snapshotDuringReadContext struct {
	context.Context
	checks int
	change func()
}

func (c *snapshotDuringReadContext) Err() error {
	c.checks++
	if c.checks == 2 {
		c.change()
	}
	return c.Context.Err()
}

func TestSnapshotRevalidationRefusesChangeDuringRead(t *testing.T) {
	for _, kind := range []string{"mtime", "size", "mode", "root", "inode", "cancel"} {
		t.Run(kind, func(t *testing.T) {
			parent := t.TempDir()
			dir := filepath.Join(parent, "root")
			want := []byte("original")
			snapshotWrite(t, dir, "source", want, 0600)
			name := filepath.Join(dir, "source")
			prior := snapshotStat(t, name)
			root := snapshotOpenRoot(t, dir)
			base, cancel := context.WithCancel(context.Background())
			defer cancel()
			ctx := &snapshotDuringReadContext{Context: base, change: func() {
				var err error
				switch kind {
				case "mtime":
					stamp := prior.ModTime().Add(2 * time.Hour)
					err = os.Chtimes(name, stamp, stamp)
				case "size":
					err = os.Truncate(name, prior.Size()-1)
				case "mode":
					err = os.Chmod(name, 0400)
				case "root":
					err = os.Rename(dir, filepath.Join(parent, "old-root"))
					if err == nil {
						snapshotWrite(t, dir, "source", want, 0600)
					}
				case "inode":
					err = os.Rename(name, filepath.Join(dir, "old-source"))
					if err == nil {
						snapshotWrite(t, dir, "source", want, 0600)
					}
				case "cancel":
					cancel()
				}
				if err != nil {
					t.Fatal(err)
				}
			}}
			hash := sha256.Sum256(want)
			err := verifySnapshotRegular(ctx, root, "source", prior, hash[:])
			if ctx.checks < 2 || err == nil {
				t.Fatalf("change during verification was not observed: checks=%d err=%v", ctx.checks, err)
			}
			if kind == "cancel" && !errors.Is(err, context.Canceled) {
				t.Fatal("cancellation cause was lost", err)
			}
		})
	}
}
