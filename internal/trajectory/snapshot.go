package trajectory

import (
	"archive/tar"
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"hash"
	"io"
	"io/fs"
	"os"
	"path"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"parley-deck-cli/internal/budget"
	"parley-deck-cli/internal/evidence"
	"parley-deck-cli/internal/fsutil"
)

const MaxSnapshotBytes int64 = 256 << 20
const MaxSnapshotFileBytes int64 = 64 << 20
const MaxSnapshotEntries = 100000

// SnapshotRef is safe metadata. Tar bodies contain private source material and
// must never be included in public telemetry or protocol artifacts.
type SnapshotRef struct {
	Version int    `json:"version"`
	SHA256  string `json:"sha256"`
	Bytes   int64  `json:"bytes"`
}
type snapshotMetadata struct {
	Version int    `json:"version"`
	Source  Source `json:"source"`
}

func (r SnapshotRef) valid() bool {
	return r.Version == 1 && validHash(r.SHA256) && r.Bytes >= 1024 && r.Bytes <= MaxSnapshotBytes && r.Bytes%512 == 0
}
func snapshotDirectory(b budget.CycleBinding) string {
	return filepath.Join(filepath.Dir(b.Store.Dir), "trajectory-snapshots")
}
func snapshotPath(dir string, r SnapshotRef) string { return filepath.Join(dir, r.SHA256+".tar") }

// Restrict portable archive paths, not source authority. Unsupported names
// refuse capture visibly instead of silently omitting files from the snapshot.
func validSnapshotPath(name string) bool {
	if len(name) == 0 || len(name) > 4096 || !fs.ValidPath(name) || name == "." || strings.ContainsAny(name, "\\:\x00\r\n") {
		return false
	}
	for _, c := range name {
		if c < 32 || c == 127 {
			return false
		}
	}
	for _, part := range strings.Split(name, "/") {
		if strings.EqualFold(part, ".git") {
			return false
		}
	}
	return true
}
func snapshotInventory(ctx context.Context, root string) ([]string, []string, error) {
	raw, err := gitOutput(ctx, root, "ls-files", "-c", "-o", "--exclude-standard", "-z")
	if err != nil {
		return nil, nil, err
	}
	names := strings.Split(string(raw), "\x00")
	if len(names) > MaxSnapshotEntries+1 {
		return nil, nil, errors.New("source inventory exceeds snapshot bound")
	}
	rooted, err := os.OpenRoot(root)
	if err != nil {
		return nil, nil, err
	}
	defer rooted.Close()
	var present, absent []string
	var total int64
	seen := map[string]bool{}
	for _, name := range names {
		if name == "" {
			continue
		}
		if !validSnapshotPath(name) || seen[name] {
			return nil, nil, errors.New("unsupported or duplicate source inventory path")
		}
		seen[name] = true
		info, err := rooted.Lstat(filepath.FromSlash(name))
		if os.IsNotExist(err) {
			absent = append(absent, name)
			continue
		}
		if err != nil {
			return nil, nil, errors.New("source inventory cannot be read within its root")
		}
		if info.Mode().IsRegular() {
			if info.Size() > MaxSnapshotFileBytes || info.Size() > MaxSnapshotBytes-total {
				return nil, nil, errors.New("source exceeds snapshot size bounds")
			}
			total += info.Size()
		} else if info.Mode()&os.ModeSymlink == 0 {
			return nil, nil, errors.New("source inventory contains an unsupported entry type")
		}
		present = append(present, name)
	}
	sort.Strings(present)
	sort.Strings(absent)
	return present, absent, nil
}

// Deleted tracked paths are absent material, not an invented existing file.
// Ordinary evidence.TreeDigest is unchanged; only trajectory's actual worktree
// observation supplies these explicit, rechecked absences as exclusions.
func trajectoryTreeDigest(ctx context.Context, root string) (string, error) {
	_, absent, err := snapshotInventory(ctx, root)
	if err != nil {
		return "", err
	}
	return evidence.TreeDigest(root, absent...)
}

type snapshotWriter struct {
	w io.Writer
	n int64
}

func (w *snapshotWriter) Write(p []byte) (int, error) {
	if int64(len(p)) > MaxSnapshotBytes-w.n {
		return 0, errors.New("snapshot archive exceeds size bound")
	}
	n, err := w.w.Write(p)
	w.n += int64(n)
	return n, err
}

type contextReader struct {
	ctx context.Context
	r   io.Reader
}

func (r contextReader) Read(p []byte) (int, error) {
	if err := r.ctx.Err(); err != nil {
		return 0, err
	}
	return r.r.Read(p)
}
func snapshotHeader(name string, kind byte, mode int64, size int64, link string) *tar.Header {
	return &tar.Header{Name: name, Typeflag: kind, Mode: mode, Size: size, Linkname: link, ModTime: time.Unix(0, 0), Format: tar.FormatPAX}
}
func privateSnapshotDirectory(dir string, create bool) error {
	if create {
		if err := os.MkdirAll(dir, 0700); err != nil {
			return err
		}
	}
	info, err := os.Lstat(dir)
	if err != nil {
		return err
	}
	if !info.IsDir() || info.Mode()&os.ModeSymlink != 0 || info.Mode().Perm()&0077 != 0 {
		return errors.New("snapshot store must be a private real directory")
	}
	return nil
}

// CaptureSnapshot freezes actual bytes under an immutable content address. It
// never stages, commits, stashes or rewrites the source repository. A failed
// publication grants no capture and does not replace an existing archive.
func CaptureSnapshot(ctx context.Context, root, dir string, want Source) (SnapshotRef, error) {
	if !sourceValid(want) {
		return SnapshotRef{}, errors.New("snapshot requires a valid actual source observation")
	}
	if err := privateSnapshotDirectory(dir, true); err != nil {
		return SnapshotRef{}, err
	}
	wait, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()
	release, err := budget.AcquireResourceGuard(wait, dir)
	if err != nil {
		return SnapshotRef{}, err
	}
	defer release()
	actual, err := Observe(ctx, root)
	if err != nil {
		return SnapshotRef{}, err
	}
	if actual != want {
		return SnapshotRef{}, errors.New("source changed before archive capture")
	}
	names, _, err := snapshotInventory(ctx, root)
	if err != nil {
		return SnapshotRef{}, err
	}
	rooted, err := os.OpenRoot(root)
	if err != nil {
		return SnapshotRef{}, err
	}
	defer rooted.Close()
	origin, err := rooted.Stat(".")
	if err != nil || !origin.IsDir() {
		return SnapshotRef{}, errors.New("snapshot source root is unavailable")
	}
	f, err := os.CreateTemp(dir, ".snapshot-*.tmp")
	if err != nil {
		return SnapshotRef{}, err
	}
	defer os.Remove(f.Name())
	defer f.Close()
	digestWriter := sha256.New()
	bounded := &snapshotWriter{w: io.MultiWriter(f, digestWriter)}
	tw := tar.NewWriter(bounded)
	meta, err := canonical(snapshotMetadata{Version: 1, Source: want})
	if err != nil {
		return SnapshotRef{}, err
	}
	if err = tw.WriteHeader(snapshotHeader("snapshot.json", tar.TypeReg, 0600, int64(len(meta)), "")); err != nil {
		return SnapshotRef{}, err
	}
	if _, err = tw.Write(meta); err != nil {
		return SnapshotRef{}, err
	}
	for _, name := range names {
		if err = ctx.Err(); err != nil {
			return SnapshotRef{}, err
		}
		if err = captureSnapshotMember(ctx, root, origin, tw, name); err != nil {
			return SnapshotRef{}, err
		}
	}
	if err = tw.Close(); err != nil {
		return SnapshotRef{}, err
	}
	if err = fsutil.SyncFile(f); err != nil {
		return SnapshotRef{}, err
	}
	if err = f.Close(); err != nil {
		return SnapshotRef{}, err
	}
	ref := SnapshotRef{Version: 1, SHA256: hex.EncodeToString(digestWriter.Sum(nil)), Bytes: bounded.n}
	if err = inspectSnapshotFile(ctx, f.Name(), ref, want, nil); err != nil {
		return SnapshotRef{}, err
	}
	actual, err = Observe(ctx, root)
	if err != nil {
		return SnapshotRef{}, err
	}
	if actual != want {
		return SnapshotRef{}, errors.New("source changed while snapshot was being archived")
	}
	final := snapshotPath(dir, ref)
	if _, err = os.Lstat(final); err == nil {
		if err = CheckSnapshot(ctx, dir, ref, want); err != nil {
			return SnapshotRef{}, err
		}
		return ref, nil
	} else if !os.IsNotExist(err) {
		return SnapshotRef{}, err
	}
	if err = fsutil.ReplaceSyncedFile(f.Name(), final); err != nil {
		return SnapshotRef{}, err
	}
	if err = CheckSnapshot(ctx, dir, ref, want); err != nil {
		return SnapshotRef{}, err
	}
	return ref, nil
}

// AppleVirtIOFS can retain an obsolete view behind a long-lived directory
// handle. Reopen from the original path for each member, then require the same
// pinned directory identity before using that fresh, contained view. Opening
// "." through the old Root can inherit its obsolete view instead.
func captureSnapshotMember(ctx context.Context, root string, origin os.FileInfo, tw *tar.Writer, name string) error {
	entryRoot, err := os.OpenRoot(root)
	if err != nil {
		return err
	}
	defer entryRoot.Close()
	actual, err := entryRoot.Stat(".")
	if err != nil || !actual.IsDir() || !os.SameFile(origin, actual) {
		return errors.New("snapshot source root changed during capture")
	}
	info, err := entryRoot.Lstat(filepath.FromSlash(name))
	if err != nil {
		return errors.New("source entry disappeared during capture")
	}
	switch {
	case info.Mode().IsRegular():
		return captureSnapshotRegular(ctx, entryRoot, tw, name, info)
	case info.Mode()&os.ModeSymlink != 0:
		target, err := entryRoot.Readlink(filepath.FromSlash(name))
		if err != nil {
			return err
		}
		if !validSnapshotLink(name, target) {
			return errors.New("snapshot link is absolute, escaping or unsupported")
		}
		return tw.WriteHeader(snapshotHeader("files/"+name, tar.TypeSymlink, int64(info.Mode().Perm()), 0, target))
	default:
		return errors.New("snapshot source contains an unsupported entry type")
	}
}

func captureSnapshotRegular(ctx context.Context, root *os.Root, tw *tar.Writer, name string, prior os.FileInfo) error {
	src, err := root.Open(filepath.FromSlash(name))
	if err != nil {
		return errors.New("source file cannot be opened within its root")
	}
	opened, err := src.Stat()
	if err != nil || !os.SameFile(prior, opened) || !opened.Mode().IsRegular() {
		src.Close()
		return errors.New("source file changed during open")
	}
	// Lstat identifies the entry; the descriptor supplies the bytes and their
	// metadata. The complete archive must still match the frozen material tree.
	if opened.Size() > MaxSnapshotFileBytes {
		src.Close()
		return errors.New("source file exceeds snapshot bound")
	}
	err = tw.WriteHeader(snapshotHeader("files/"+name, tar.TypeReg, int64(opened.Mode().Perm()), opened.Size(), ""))
	copied := sha256.New()
	if err == nil {
		_, err = io.CopyN(io.MultiWriter(tw, copied), contextReader{ctx, src}, opened.Size())
	}
	after, statErr := src.Stat()
	closeErr := src.Close()
	if err != nil {
		return err
	}
	if statErr != nil || closeErr != nil || after.Size() != opened.Size() || after.Mode() != opened.Mode() {
		return errors.New("source file changed during capture")
	}
	return verifySnapshotRegular(ctx, root, name, opened, copied.Sum(nil))
}

// An open AppleVirtIOFS descriptor can retain old bytes/metadata or receive a
// later timestamp without a byte change. Always require one fresh, stable read
// of the same pinned material to confirm the bytes already written to the tar.
// A transition in that read or any byte/identity/size/mode difference refuses. The
// entire archive still has to match its original expected tree before publish.
func verifySnapshotRegular(ctx context.Context, root *os.Root, name string, copiedInfo os.FileInfo, copiedHash []byte) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	origin, err := root.Stat(".")
	if err != nil {
		return err
	}
	fresh, err := os.OpenRoot(root.Name())
	if err != nil {
		return err
	}
	defer fresh.Close()
	actual, err := fresh.Stat(".")
	if err != nil || !actual.IsDir() || !os.SameFile(origin, actual) {
		return errors.New("snapshot source root changed during timestamp revalidation")
	}
	entry, err := fresh.Lstat(filepath.FromSlash(name))
	if err != nil || !entry.Mode().IsRegular() || !os.SameFile(copiedInfo, entry) {
		return errors.New("source file changed before timestamp revalidation")
	}
	src, err := fresh.Open(filepath.FromSlash(name))
	if err != nil {
		return err
	}
	opened, err := src.Stat()
	if err != nil || !opened.Mode().IsRegular() || !os.SameFile(copiedInfo, opened) || opened.Size() != copiedInfo.Size() || opened.Mode() != copiedInfo.Mode() || opened.Size() > MaxSnapshotFileBytes {
		src.Close()
		return errors.New("source file changed during timestamp revalidation open")
	}
	verified := sha256.New()
	_, err = io.CopyN(verified, contextReader{ctx, src}, opened.Size())
	after, statErr := src.Stat()
	closeErr := src.Close()
	if err != nil {
		return err
	}
	if statErr != nil || closeErr != nil || after.Size() != opened.Size() || after.Mode() != opened.Mode() || !after.ModTime().Equal(opened.ModTime()) {
		return errors.New("source file changed during timestamp revalidation")
	}
	if !bytes.Equal(verified.Sum(nil), copiedHash) {
		return errors.New("source bytes changed during timestamp revalidation")
	}
	// Recheck the names as well as the open descriptor: a replaced root or
	// directory entry must not be hidden by the original still-readable inode.
	current, err := os.OpenRoot(root.Name())
	if err != nil {
		return err
	}
	defer current.Close()
	currentRoot, err := current.Stat(".")
	if err != nil || !os.SameFile(origin, currentRoot) {
		return errors.New("snapshot source root changed during timestamp revalidation")
	}
	currentFile, err := current.Lstat(filepath.FromSlash(name))
	if err != nil || !currentFile.Mode().IsRegular() || !os.SameFile(opened, currentFile) || currentFile.Size() != opened.Size() || currentFile.Mode() != opened.Mode() {
		return errors.New("source file replaced during timestamp revalidation")
	}
	return ctx.Err()
}

func validSnapshotLink(name, target string) bool {
	if target == "" || len(target) > 4096 || path.IsAbs(target) || strings.ContainsAny(target, "\\:\x00\r\n") {
		return false
	}
	for _, part := range strings.Split(target, "/") {
		if part != "" && part != "." && part != ".." && !validSnapshotPath(part) {
			return false
		}
	}
	resolved := path.Clean(path.Join(path.Dir(name), target))
	return resolved == "." || validSnapshotPath(resolved)
}

type snapshotMember struct {
	kind   byte
	target string
}

func validateSnapshotLinks(members map[string]snapshotMember) error {
	dirs := map[string]bool{".": true}
	for name := range members {
		for parent := path.Dir(name); parent != "."; parent = path.Dir(parent) {
			if _, exists := members[parent]; exists {
				return errors.New("snapshot member is nested under a file or symlink")
			}
			dirs[parent] = true
		}
	}
	for name, m := range members {
		if m.kind != tar.TypeSymlink {
			continue
		}
		// Resolve one component at a time: cleaning "link/.." before following
		// the link would erase its filesystem meaning and could hide an escape.
		pending := strings.Split(path.Dir(name)+"/"+m.target, "/")
		var resolved []string
		depth := 0
		for len(pending) > 0 {
			part := pending[0]
			pending = pending[1:]
			if len(resolved) > 0 && !dirs[strings.Join(resolved, "/")] {
				return errors.New("snapshot symlink traverses a regular file")
			}
			switch part {
			case "", ".":
				continue
			case "..":
				if len(resolved) == 0 {
					return errors.New("snapshot symlink chain escapes its source")
				}
				resolved = resolved[:len(resolved)-1]
				continue
			}
			prefix := strings.Join(append(append([]string{}, resolved...), part), "/")
			entry, ok := members[prefix]
			if ok && entry.kind == tar.TypeSymlink {
				depth++
				if depth > 40 {
					return errors.New("snapshot symlink cycle or excessive chain")
				}
				pending = append(strings.Split(entry.target, "/"), pending...)
				continue
			}
			if !ok && !dirs[prefix] {
				return errors.New("snapshot symlink target is not retained source")
			}
			resolved = append(resolved, part)
		}
	}
	return nil
}

type countedHash struct {
	hash.Hash
	n int64
}

func (h *countedHash) Write(p []byte) (int, error) {
	n, err := h.Hash.Write(p)
	h.n += int64(n)
	return n, err
}
func archiveHeaderSupported(h *tar.Header) bool {
	if h.Mode < 0 || h.Mode > 0777 || h.Uid != 0 || h.Gid != 0 || h.Uname != "" || h.Gname != "" || h.Devmajor != 0 || h.Devminor != 0 || !h.ModTime.Equal(time.Unix(0, 0)) || !h.AccessTime.IsZero() || !h.ChangeTime.IsZero() {
		return false
	}
	for key := range h.PAXRecords {
		if key != "path" && key != "linkpath" {
			return false
		}
	}
	return len(h.Xattrs) == 0
}

// inspectSnapshotFile verifies the byte address and independently recomputes
// evidence.TreeDigest's file/link serialization from the retained tar members.
// Restore tests compare this with the actual restored filesystem digest. This
// serialization is part of archive v1's compatibility boundary.
func inspectSnapshotFile(ctx context.Context, file string, ref SnapshotRef, want Source, destination *os.Root) error {
	if !ref.valid() || !sourceValid(want) {
		return errors.New("invalid snapshot reference or source binding")
	}
	info, err := os.Lstat(file)
	if err != nil {
		return err
	}
	if !info.Mode().IsRegular() || info.Mode().Perm()&0077 != 0 || info.Size() != ref.Bytes {
		return errors.New("snapshot must be an exact bounded private regular file")
	}
	f, err := os.Open(file)
	if err != nil {
		return err
	}
	defer f.Close()
	opened, err := f.Stat()
	if err != nil || !os.SameFile(info, opened) {
		return errors.New("snapshot changed during open")
	}
	archiveHash := &countedHash{Hash: sha256.New()}
	reader := tar.NewReader(io.TeeReader(contextReader{ctx, io.LimitReader(f, MaxSnapshotBytes+1)}, archiveHash))
	header, err := reader.Next()
	if err != nil {
		return err
	}
	if header.Name != "snapshot.json" || header.Typeflag != tar.TypeReg || header.Size < 1 || header.Size > 8192 || header.Mode != 0600 || header.Linkname != "" || !archiveHeaderSupported(header) {
		return errors.New("invalid snapshot metadata member")
	}
	data, err := io.ReadAll(reader)
	if err != nil {
		return err
	}
	var meta snapshotMetadata
	dec := json.NewDecoder(bytes.NewReader(data))
	dec.DisallowUnknownFields()
	if err = dec.Decode(&meta); err != nil {
		return err
	}
	canonicalMeta, err := canonical(meta)
	if err != nil || !bytes.Equal(data, canonicalMeta) || meta.Version != 1 || meta.Source != want {
		return errors.New("snapshot metadata differs from the frozen source")
	}
	// Re-encode the accepted members to require the exact versioned archive
	// format, including its footer. archive/tar alone accepts missing trailers
	// and normalizes extension headers that we must not silently authorize.
	canonicalHash := &countedHash{Hash: sha256.New()}
	tw := tar.NewWriter(&snapshotWriter{w: canonicalHash})
	if err = tw.WriteHeader(snapshotHeader("snapshot.json", tar.TypeReg, 0600, int64(len(data)), "")); err != nil {
		return err
	}
	if _, err = tw.Write(data); err != nil {
		return err
	}
	treeHash := sha256.New()
	members := map[string]snapshotMember{}
	previous := ""
	for {
		h, err := reader.Next()
		if errors.Is(err, io.EOF) {
			break
		}
		if err != nil {
			return err
		}
		if err = ctx.Err(); err != nil {
			return err
		}
		if !strings.HasPrefix(h.Name, "files/") || !archiveHeaderSupported(h) {
			return errors.New("unsupported snapshot header")
		}
		name := strings.TrimPrefix(h.Name, "files/")
		if !validSnapshotPath(name) || name <= previous || len(members) >= MaxSnapshotEntries {
			return errors.New("snapshot paths are invalid, duplicated, unordered or excessive")
		}
		previous = name
		switch h.Typeflag {
		case tar.TypeReg:
			if h.Size < 0 || h.Size > MaxSnapshotFileBytes || h.Linkname != "" {
				return errors.New("invalid snapshot regular file")
			}
			fmt.Fprintf(treeHash, "file:%o:%d:%s\n", h.Mode, h.Size, name)
			if err = tw.WriteHeader(snapshotHeader(h.Name, tar.TypeReg, h.Mode, h.Size, "")); err != nil {
				return err
			}
			var target *os.File
			var output io.Writer = io.MultiWriter(treeHash, tw)
			if destination != nil {
				if err = destination.MkdirAll(filepath.FromSlash(path.Dir(name)), 0700); err != nil {
					return err
				}
				target, err = destination.OpenFile(filepath.FromSlash(name), os.O_WRONLY|os.O_CREATE|os.O_EXCL, 0600)
				if err != nil {
					return err
				}
				output = io.MultiWriter(treeHash, tw, target)
			}
			copied, copyErr := io.Copy(output, reader)
			if target != nil {
				modeErr := target.Chmod(os.FileMode(h.Mode))
				syncErr := fsutil.SyncFile(target)
				closeErr := target.Close()
				if copyErr == nil {
					copyErr = errors.Join(modeErr, syncErr, closeErr)
				}
			}
			if copyErr != nil {
				return copyErr
			}
			if copied != h.Size {
				return errors.New("truncated snapshot file")
			}
			members[name] = snapshotMember{kind: tar.TypeReg}
		case tar.TypeSymlink:
			if h.Size != 0 || !validSnapshotLink(name, h.Linkname) {
				return errors.New("invalid snapshot symlink")
			}
			fmt.Fprintf(treeHash, "link:%o:%s\n%s\n", h.Mode, name, h.Linkname)
			if err = tw.WriteHeader(snapshotHeader(h.Name, tar.TypeSymlink, h.Mode, 0, h.Linkname)); err != nil {
				return err
			}
			members[name] = snapshotMember{kind: tar.TypeSymlink, target: h.Linkname}
		default:
			return errors.New("snapshot contains an unsupported member type")
		}
	}
	if err = tw.Close(); err != nil {
		return err
	}
	if canonicalHash.n != ref.Bytes || hex.EncodeToString(canonicalHash.Sum(nil)) != ref.SHA256 {
		return errors.New("snapshot archive is truncated or noncanonical")
	}
	if len(members) == 0 || archiveHash.n != ref.Bytes || hex.EncodeToString(archiveHash.Sum(nil)) != ref.SHA256 || hex.EncodeToString(treeHash.Sum(nil)) != want.Tree.SHA256 {
		return errors.New("snapshot archive or complete source digest mismatch")
	}
	if err = validateSnapshotLinks(members); err != nil {
		return err
	}
	after, err := os.Lstat(file)
	if err != nil || !os.SameFile(info, after) || after.Size() != info.Size() || after.Mode() != info.Mode() || !after.ModTime().Equal(info.ModTime()) {
		return errors.New("snapshot changed during verification")
	}
	if destination != nil {
		// All regular files were written before any link is created. No archive
		// member can direct extraction through a symlink parent.
		var links []string
		for name, m := range members {
			if m.kind == tar.TypeSymlink {
				links = append(links, name)
			}
		}
		sort.Strings(links)
		for _, name := range links {
			if err = destination.MkdirAll(filepath.FromSlash(path.Dir(name)), 0700); err != nil {
				return err
			}
			if err = destination.Symlink(members[name].target, filepath.FromSlash(name)); err != nil {
				return err
			}
		}
	}
	return nil
}
func CheckSnapshot(ctx context.Context, dir string, ref SnapshotRef, want Source) error {
	if err := privateSnapshotDirectory(dir, false); err != nil {
		return err
	}
	return inspectSnapshotFile(ctx, snapshotPath(dir, ref), ref, want, nil)
}

// RestoreSnapshot allocates a new private directory and restores source files
// only. It creates no Git metadata and makes no commit/ancestry claim. The
// orchestrator must separately bind a verifier's isolated Git execution roots.
// On error only this new directory is removed; an existing target is never used.
func RestoreSnapshot(ctx context.Context, dir string, ref SnapshotRef, want Source, parent string) (result string, err error) {
	if err = privateSnapshotDirectory(dir, false); err != nil {
		return "", err
	}
	result, err = os.MkdirTemp(parent, "trajectory-source-")
	if err != nil {
		return "", err
	}
	created := result
	defer func() {
		if err != nil {
			_ = os.RemoveAll(created)
			result = ""
		}
	}()
	rooted, err := os.OpenRoot(created)
	if err != nil {
		return "", err
	}
	defer rooted.Close()
	if err = inspectSnapshotFile(ctx, snapshotPath(dir, ref), ref, want, rooted); err != nil {
		return "", err
	}
	actual, err := evidence.TreeDigest(created)
	if err != nil {
		return "", err
	}
	if actual != want.Tree.SHA256 {
		return "", errors.New("restored filesystem differs from captured source")
	}
	return created, nil
}
