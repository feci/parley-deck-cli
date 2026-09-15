package trajectory

import (
	"archive/tar"
	"bytes"
	"context"
	"encoding/json"
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"parley-deck-cli/internal/evidence"
)

func snapshotWrite(t *testing.T, root, name string, data []byte, mode os.FileMode) {
	t.Helper()
	file := filepath.Join(root, filepath.FromSlash(name))
	if err := os.MkdirAll(filepath.Dir(file), 0700); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(file, data, mode); err != nil {
		t.Fatal(err)
	}
	if err := os.Chmod(file, mode); err != nil {
		t.Fatal(err)
	}
}

func snapshotStoreFixture(t *testing.T) string {
	t.Helper()
	dir := t.TempDir()
	if err := os.Chmod(dir, 0700); err != nil {
		t.Fatal(err)
	}
	return dir
}

func snapshotFixture(t *testing.T) string {
	t.Helper()
	root := t.TempDir()
	gitFixture(t, root, "init", "-q")
	snapshotWrite(t, root, ".gitignore", []byte("ignored/\n"), 0600)
	snapshotWrite(t, root, "source", []byte("original\n"), 0600)
	snapshotWrite(t, root, "removed", []byte("retained at baseline\n"), 0600)
	gitFixture(t, root, "add", ".")
	gitFixture(t, root, "commit", "-qm", "Snapshot fixture")
	return root
}

func captureFixture(t *testing.T, root, dir string) (Source, SnapshotRef) {
	t.Helper()
	source, err := Observe(context.Background(), root)
	if err != nil {
		t.Fatal(err)
	}
	ref, err := CaptureSnapshot(context.Background(), root, dir, source)
	if err != nil {
		t.Fatal(err)
	}
	return source, ref
}

func snapshotRead(t *testing.T, path string) []byte {
	t.Helper()
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	return data
}

func TestSnapshotRoundTripRetainsDirtySourceAndDoesNotWriteGit(t *testing.T) {
	root, dir := snapshotFixture(t), snapshotStoreFixture(t)
	baseline, baselineRef := captureFixture(t, root, dir)
	if !baseline.Clean {
		t.Fatal("fixture baseline is dirty")
	}
	snapshotWrite(t, root, "source", []byte("changed\n"), 0700)
	if err := os.Remove(filepath.Join(root, "removed")); err != nil {
		t.Fatal(err)
	}
	binary := []byte{0, 1, 2, 255, '\n', 0}
	snapshotWrite(t, root, "new/binary", binary, 0640)
	snapshotWrite(t, root, "ignored/private", []byte("omitted by original Git inventory\n"), 0600)
	indexPath := filepath.Join(root, ".git", "index")
	index := snapshotRead(t, indexPath)
	head := snapshotRead(t, filepath.Join(root, ".git", "HEAD"))
	refs := gitFixture(t, root, "show-ref")
	source, ref := captureFixture(t, root, dir)
	if source.Clean || source.Tree.Commit != baseline.Tree.Commit || source.Tree.SHA256 == baseline.Tree.SHA256 {
		t.Fatal("dirty source identity was not retained")
	}
	if !bytes.Equal(index, snapshotRead(t, indexPath)) || !bytes.Equal(head, snapshotRead(t, filepath.Join(root, ".git", "HEAD"))) || refs != gitFixture(t, root, "show-ref") {
		t.Fatal("read-only capture changed the user's Git state")
	}
	_, replay := captureFixture(t, root, dir)
	if ref != replay {
		t.Fatal("unchanged source archive is not deterministic")
	}
	// The archive, not the current live files, must supply every restored byte.
	snapshotWrite(t, root, "source", []byte("edited again\n"), 0600)
	snapshotWrite(t, root, "new/binary", []byte("replaced\n"), 0600)
	parent := t.TempDir()
	snapshotWrite(t, parent, "existing", []byte("preserve\n"), 0600)
	for _, tc := range []struct {
		source Source
		ref    SnapshotRef
		dirty  bool
	}{{baseline, baselineRef, false}, {source, ref, true}} {
		restored, err := RestoreSnapshot(context.Background(), dir, tc.ref, tc.source, parent)
		if err != nil {
			t.Fatal(err)
		}
		actual, err := evidence.TreeDigest(restored)
		if err != nil || actual != tc.source.Tree.SHA256 {
			t.Fatalf("actual restored filesystem differs: %s %v", actual, err)
		}
		if _, err = os.Lstat(filepath.Join(restored, ".git")); !os.IsNotExist(err) {
			t.Fatal("restore manufactured Git metadata")
		}
		if info, err := os.Stat(restored); err != nil || info.Mode().Perm()&0077 != 0 {
			t.Fatal("restore directory is not private")
		}
		if tc.dirty {
			if string(snapshotRead(t, filepath.Join(restored, "source"))) != "changed\n" || !bytes.Equal(snapshotRead(t, filepath.Join(restored, "new/binary")), binary) {
				t.Fatal("restored current live files instead of archived contents")
			}
			info, err := os.Stat(filepath.Join(restored, "source"))
			if err != nil || info.Mode().Perm() != 0700 {
				t.Fatal("executable mode was lost")
			}
			if _, err = os.Lstat(filepath.Join(restored, "removed")); !os.IsNotExist(err) {
				t.Fatal("tracked deletion was resurrected")
			}
		} else if string(snapshotRead(t, filepath.Join(restored, "removed"))) != "retained at baseline\n" {
			t.Fatal("baseline contents were lost")
		}
		if _, err = os.Lstat(filepath.Join(restored, "ignored")); !os.IsNotExist(err) {
			t.Fatal("ignored files entered the source snapshot")
		}
	}
	if string(snapshotRead(t, filepath.Join(parent, "existing"))) != "preserve\n" || string(snapshotRead(t, filepath.Join(root, "source"))) != "edited again\n" {
		t.Fatal("restore changed pre-existing files")
	}
}

func TestSnapshotRelativeLinksRoundTrip(t *testing.T) {
	root, dir := snapshotFixture(t), snapshotStoreFixture(t)
	long := strings.Repeat("long-segment-", 12) + "/file"
	snapshotWrite(t, root, long, []byte("long path retained\n"), 0600)
	snapshotWrite(t, root, "nested/file", []byte("nested\n"), 0600)
	links := map[string]string{
		"alias": "nested", "chain": "alias/file", "nested/back": "../source",
		"long-link": long, "via-parent": "alias/../source",
	}
	for name, target := range links {
		if err := os.Symlink(target, filepath.Join(root, name)); err != nil {
			t.Skipf("symlinks unavailable: %v", err)
		}
	}
	source, ref := captureFixture(t, root, dir)
	restored, err := RestoreSnapshot(context.Background(), dir, ref, source, t.TempDir())
	if err != nil {
		t.Fatal(err)
	}
	for name, target := range links {
		got, err := os.Readlink(filepath.Join(restored, name))
		if err != nil || got != target {
			t.Fatalf("raw link target lost: %s %q %v", name, got, err)
		}
	}
	if string(snapshotRead(t, filepath.Join(restored, "via-parent"))) != "original\n" {
		t.Fatal("relative parent-link semantics changed")
	}
}

func TestSnapshotRefusesUnretainedOrEscapingLinks(t *testing.T) {
	for _, kind := range []string{"absolute-in-root", "escape", "dangling", "ignored", "cycle", "control-component"} {
		t.Run(kind, func(t *testing.T) {
			root, dir := snapshotFixture(t), snapshotStoreFixture(t)
			target := ""
			switch kind {
			case "absolute-in-root":
				target = filepath.Join(root, "source")
			case "escape":
				outside := t.TempDir()
				snapshotWrite(t, outside, "private", []byte("outside\n"), 0600)
				var err error
				target, err = filepath.Rel(root, filepath.Join(outside, "private"))
				if err != nil {
					t.Fatal(err)
				}
			case "dangling":
				target = "missing"
			case "ignored":
				snapshotWrite(t, root, "ignored/private", []byte("excluded\n"), 0600)
				target = "ignored/private"
			case "cycle":
				target = "link"
			case "control-component":
				if err := os.Mkdir(filepath.Join(root, "ignored", "bad\t"), 0700); os.IsNotExist(err) {
					if err = os.MkdirAll(filepath.Join(root, "ignored", "bad\t"), 0700); err != nil {
						t.Fatal(err)
					}
				} else if err != nil {
					t.Fatal(err)
				}
				target = "ignored/bad\t/../../source"
			}
			if err := os.Symlink(target, filepath.Join(root, "link")); err != nil {
				t.Skipf("symlinks unavailable: %v", err)
			}
			source, err := Observe(context.Background(), root)
			if err == nil {
				_, err = CaptureSnapshot(context.Background(), root, dir, source)
			}
			if err == nil {
				t.Fatal("unsupported source was archived")
			}
			archives, _ := filepath.Glob(filepath.Join(dir, "*.tar"))
			if len(archives) != 0 {
				t.Fatal("failed source capture published an archive")
			}
		})
	}
}

type archiveFixtureMember struct {
	header *tar.Header
	data   []byte
}

func archiveFixtureBytes(t *testing.T, source Source, members []archiveFixtureMember) []byte {
	t.Helper()
	var out bytes.Buffer
	w := tar.NewWriter(&out)
	meta, err := canonical(snapshotMetadata{Version: 1, Source: source})
	if err != nil {
		t.Fatal(err)
	}
	all := append([]archiveFixtureMember{{snapshotHeader("snapshot.json", tar.TypeReg, 0600, int64(len(meta)), ""), meta}}, members...)
	for _, m := range all {
		if err := w.WriteHeader(m.header); err != nil {
			t.Fatal(err)
		}
		if _, err := w.Write(m.data); err != nil {
			t.Fatal(err)
		}
	}
	if err := w.Close(); err != nil {
		t.Fatal(err)
	}
	return out.Bytes()
}

func archiveFixtureStore(t *testing.T, dir string, data []byte) SnapshotRef {
	t.Helper()
	ref := SnapshotRef{Version: 1, SHA256: digest(data), Bytes: int64(len(data))}
	if err := os.WriteFile(snapshotPath(dir, ref), data, 0600); err != nil {
		t.Fatal(err)
	}
	return ref
}

func TestSnapshotRejectsMalformedArchivesBeforeRestoration(t *testing.T) {
	root, dir := snapshotFixture(t), snapshotStoreFixture(t)
	source, ref := captureFixture(t, root, dir)
	original := snapshotRead(t, snapshotPath(dir, ref))
	reader := tar.NewReader(bytes.NewReader(original))
	if _, err := reader.Next(); err != nil {
		t.Fatal(err)
	}
	var members []archiveFixtureMember
	for {
		h, err := reader.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			t.Fatal(err)
		}
		b, err := io.ReadAll(reader)
		if err != nil {
			t.Fatal(err)
		}
		members = append(members, archiveFixtureMember{h, b})
	}
	for _, kind := range []string{
		"duplicate", "unordered", "traversal", "git-metadata", "hardlink", "device",
		"directory", "mode", "mtime", "pax", "wrong-source", "omitted-file",
		"bad-checksum", "truncated-header", "truncated-body", "missing-footer", "one-footer",
		"trailing-zero", "trailing-data", "oversize-header", "symlink-escape", "symlink-parent",
	} {
		t.Run(kind, func(t *testing.T) {
			copyMembers := make([]archiveFixtureMember, len(members))
			for i, m := range members {
				h := *m.header
				copyMembers[i] = archiveFixtureMember{&h, append([]byte{}, m.data...)}
			}
			want := source
			switch kind {
			case "duplicate":
				copyMembers = append(copyMembers[:1], copyMembers...)
			case "unordered":
				copyMembers[0], copyMembers[1] = copyMembers[1], copyMembers[0]
			case "traversal":
				copyMembers[0].header.Name = "files/../../outside"
			case "git-metadata":
				copyMembers[0].header.Name = "files/.git/config"
			case "hardlink", "device", "directory":
				kindMap := map[string]byte{"hardlink": tar.TypeLink, "device": tar.TypeChar, "directory": tar.TypeDir}
				copyMembers[0].header.Typeflag = kindMap[kind]
				copyMembers[0].header.Size = 0
				copyMembers[0].data = nil
				if kind == "hardlink" {
					copyMembers[0].header.Linkname = "files/source"
				}
			case "mode":
				copyMembers[0].header.Mode = 04755
			case "mtime":
				copyMembers[0].header.ModTime = time.Unix(1, 0)
			case "pax":
				copyMembers[0].header.Format = tar.FormatPAX
				copyMembers[0].header.PAXRecords = map[string]string{"comment": "not permitted"}
			case "wrong-source":
				want.Tree.SHA256 = digest([]byte("another source"))
			case "omitted-file":
				copyMembers = copyMembers[1:]
			case "symlink-escape":
				copyMembers[0].header = snapshotHeader("files/.escape", tar.TypeSymlink, 0777, 0, "../outside")
				copyMembers[0].data = nil
			case "symlink-parent":
				copyMembers = append(copyMembers, archiveFixtureMember{snapshotHeader("files/source/child", tar.TypeReg, 0600, 0, ""), nil})
			}
			data := archiveFixtureBytes(t, want, copyMembers)
			switch kind {
			case "bad-checksum":
				data[100] ^= 1
			case "truncated-header":
				data = data[:511]
			case "truncated-body":
				data = data[:len(data)-1537]
			case "missing-footer":
				data = data[:len(data)-1024]
			case "one-footer":
				data = data[:len(data)-512]
			case "trailing-zero":
				data = append(data, make([]byte, 512)...)
			case "trailing-data":
				data = append(data, bytes.Repeat([]byte{'x'}, 512)...)
			case "oversize-header":
				var out bytes.Buffer
				w := tar.NewWriter(&out)
				meta, _ := canonical(snapshotMetadata{Version: 1, Source: source})
				if err := w.WriteHeader(snapshotHeader("snapshot.json", tar.TypeReg, 0600, int64(len(meta)), "")); err != nil {
					t.Fatal(err)
				}
				if _, err := w.Write(meta); err != nil {
					t.Fatal(err)
				}
				if err := w.WriteHeader(snapshotHeader("files/oversize", tar.TypeReg, 0600, MaxSnapshotFileBytes+1, "")); err != nil {
					t.Fatal(err)
				}
				data = out.Bytes() // Declared body intentionally absent.
			}
			ref := archiveFixtureStore(t, dir, data) // Even a freshly matching hash must refuse.
			if err := CheckSnapshot(context.Background(), dir, ref, source); err == nil {
				t.Fatal("malformed archive passed verification")
			}
			parent := t.TempDir()
			snapshotWrite(t, parent, "existing", []byte("untouched\n"), 0600)
			if restored, err := RestoreSnapshot(context.Background(), dir, ref, source, parent); err == nil || restored != "" {
				t.Fatalf("malformed archive restored: %q %v", restored, err)
			}
			entries, err := os.ReadDir(parent)
			if err != nil || len(entries) != 1 || entries[0].Name() != "existing" {
				t.Fatal("failed extraction left a directory or changed existing content")
			}
		})
	}
}

func TestSnapshotRejectsLossDriftAndUnsafeStorage(t *testing.T) {
	root, dir := snapshotFixture(t), snapshotStoreFixture(t)
	source, ref := captureFixture(t, root, dir)
	original := snapshotRead(t, snapshotPath(dir, ref))
	snapshotWrite(t, root, "source", []byte("new\n"), 0600)
	if _, err := CaptureSnapshot(context.Background(), root, dir, source); err == nil {
		t.Fatal("source drift admitted")
	}
	snapshotWrite(t, root, "source", []byte("original\n"), 0600)
	corrupt := append([]byte{}, original...)
	corrupt[0] ^= 1
	if err := os.WriteFile(snapshotPath(dir, ref), corrupt, 0600); err != nil {
		t.Fatal(err)
	}
	if _, err := CaptureSnapshot(context.Background(), root, dir, source); err == nil {
		t.Fatal("existing corrupted content address replaced or accepted")
	}
	if !bytes.Equal(corrupt, snapshotRead(t, snapshotPath(dir, ref))) {
		t.Fatal("capture overwrote retained corrupt evidence")
	}
	if err := os.Remove(snapshotPath(dir, ref)); err != nil {
		t.Fatal(err)
	}
	if err := CheckSnapshot(context.Background(), dir, ref, source); err == nil {
		t.Fatal("lost archive accepted")
	}
	// A directory at the immutable content address is a real publication failure.
	if err := os.Mkdir(snapshotPath(dir, ref), 0700); err != nil {
		t.Fatal(err)
	}
	if _, err := CaptureSnapshot(context.Background(), root, dir, source); err == nil {
		t.Fatal("unpublishable archive reported captured")
	}
	public := t.TempDir()
	if err := os.Chmod(public, 0755); err != nil {
		t.Fatal(err)
	}
	if _, err := CaptureSnapshot(context.Background(), root, public, source); err == nil {
		t.Fatal("public source archive directory accepted")
	}
	cancelled, cancel := context.WithCancel(context.Background())
	cancel()
	if _, err := CaptureSnapshot(cancelled, root, snapshotStoreFixture(t), source); err == nil {
		t.Fatal("cancelled capture succeeded")
	}
	// Metadata must never contain raw source, paths or commands.
	data, err := json.Marshal(ref)
	if err != nil || bytes.Contains(data, []byte(root)) || bytes.Contains(data, []byte("original")) {
		t.Fatal("private source escaped through the public reference")
	}
}

func TestSnapshotRejectsOversizedSourceBeforeDigest(t *testing.T) {
	root := snapshotFixture(t)
	f, err := os.Create(filepath.Join(root, "oversize"))
	if err != nil {
		t.Fatal(err)
	}
	err = f.Truncate(MaxSnapshotFileBytes + 1)
	closeErr := f.Close()
	if err != nil || closeErr != nil {
		t.Fatal(err, closeErr)
	}
	if _, err := Observe(context.Background(), root); err == nil {
		t.Fatal("oversized source accepted before bounded capture")
	}
}

func TestSnapshotLinkResolutionPreservesComponentSemantics(t *testing.T) {
	for _, tc := range []struct {
		name    string
		members map[string]snapshotMember
	}{
		{"file-parent", map[string]snapshotMember{"a": {kind: tar.TypeReg}, "a/child": {kind: tar.TypeReg}}},
		{"link-parent", map[string]snapshotMember{"a": {kind: tar.TypeSymlink, target: "."}, "a/child": {kind: tar.TypeReg}}},
		{"missing-target", map[string]snapshotMember{"a": {kind: tar.TypeSymlink, target: "missing"}}},
		{"cycle", map[string]snapshotMember{"a": {kind: tar.TypeSymlink, target: "b"}, "b": {kind: tar.TypeSymlink, target: "a"}}},
		{"escape-after-link", map[string]snapshotMember{"nested/a": {kind: tar.TypeSymlink, target: ".."}, "b": {kind: tar.TypeSymlink, target: "nested/a/../source"}, "source": {kind: tar.TypeReg}}},
		{"file-before-parent", map[string]snapshotMember{"a": {kind: tar.TypeReg}, "b": {kind: tar.TypeSymlink, target: "a/../source"}, "source": {kind: tar.TypeReg}}},
	} {
		t.Run(tc.name, func(t *testing.T) {
			if err := validateSnapshotLinks(tc.members); err == nil {
				t.Fatal("invalid filesystem link semantics accepted")
			}
		})
	}
}

func TestSnapshotConcurrentCapturePublishesOneExactArchive(t *testing.T) {
	root, dir := snapshotFixture(t), snapshotStoreFixture(t)
	source, err := Observe(context.Background(), root)
	if err != nil {
		t.Fatal(err)
	}
	type result struct {
		ref SnapshotRef
		err error
	}
	done := make(chan result, 2)
	for i := 0; i < 2; i++ {
		go func() {
			ref, err := CaptureSnapshot(context.Background(), root, dir, source)
			done <- result{ref, err}
		}()
	}
	a, b := <-done, <-done
	if a.err != nil || b.err != nil || a.ref != b.ref {
		t.Fatalf("concurrent publication differs: %+v %+v", a, b)
	}
	files, err := filepath.Glob(filepath.Join(dir, "*.tar"))
	if err != nil || len(files) != 1 {
		t.Fatal("same content published more than once", err)
	}
}

func TestSnapshotRejectsSourceChangedDuringArchiveWrite(t *testing.T) {
	root, dir := snapshotFixture(t), snapshotStoreFixture(t)
	f, err := os.Create(filepath.Join(root, "a-large-source"))
	if err != nil {
		t.Fatal(err)
	}
	err = f.Truncate(32 << 20)
	closeErr := f.Close()
	if err != nil || closeErr != nil {
		t.Fatal(err, closeErr)
	}
	source, err := Observe(context.Background(), root)
	if err != nil {
		t.Fatal(err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	changed := make(chan error, 1)
	go func() {
		tick := time.NewTicker(time.Millisecond)
		defer tick.Stop()
		for {
			select {
			case <-ctx.Done():
				changed <- ctx.Err()
				return
			case <-tick.C:
				staged, err := filepath.Glob(filepath.Join(dir, ".snapshot-*.tmp"))
				if err != nil {
					changed <- err
					return
				}
				if len(staged) > 0 {
					changed <- os.WriteFile(filepath.Join(root, "source"), []byte("changed during capture\n"), 0600)
					return
				}
			}
		}
	}()
	ref, err := CaptureSnapshot(ctx, root, dir, source)
	cancel()
	if changeErr := <-changed; changeErr != nil {
		t.Fatalf("fixture did not change a live archive write: %v", changeErr)
	}
	if err == nil || ref != (SnapshotRef{}) {
		t.Fatal("mixed-source archive was accepted")
	}
	files, globErr := filepath.Glob(filepath.Join(dir, "*.tar"))
	if globErr != nil || len(files) != 0 {
		t.Fatal("source drift published an archive", globErr)
	}
}
