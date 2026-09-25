//go:build windows

package winprobe

import (
	"bytes"
	"errors"
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"testing"
	"unsafe"

	"golang.org/x/sys/windows"
)

// volumeFileSystem reports the filesystem name of the volume containing path
// (H3/H1 context: fixture placement and the NTFS coverage-envelope claim).
func volumeFileSystem(t *testing.T, path string) string {
	t.Helper()
	root := filepath.VolumeName(path)
	if root == "" {
		root = `\`
	} else {
		root += `\`
	}
	p, err := windows.UTF16PtrFromString(root)
	if err != nil {
		t.Fatalf("H3 volume root %q: %v", root, err)
	}
	var fsName [256]uint16
	if err := windows.GetVolumeInformation(p, nil, 0, nil, nil, nil, &fsName[0], uint32(len(fsName))); err != nil {
		t.Fatalf("H3 GetVolumeInformation(%q): %v", root, err)
	}
	return windows.UTF16ToString(fsName[:])
}

// TestH3RunnerVolumeIsNTFS asserts the FINAL §G H3 volume fact: the runner's
// working volume is NTFS. A non-NTFS volume voids the coverage-envelope claim
// and stops the leg for re-baselining (mapped outcome — never silently
// absorbed).
func TestH3RunnerVolumeIsNTFS(t *testing.T) {
	wd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	got := volumeFileSystem(t, wd)
	t.Logf("H3 working dir %s volume filesystem=%s", wd, got)
	if tmp := os.TempDir(); tmp != "" {
		t.Logf("H3 TEMP %s volume filesystem=%s", tmp, volumeFileSystem(t, tmp))
	}
	if got != "NTFS" {
		t.Fatalf("H3 runner volume is %q, not NTFS — coverage-envelope claim voided; leg stops for re-baselining", got)
	}
}

// TestH3RenameMechanics pins the §G H3 mechanics: a directory move is accepted,
// a no-REPLACE move onto an existing destination is rejected, and
// MOVEFILE_WRITE_THROUGH|REPLACE_EXISTING is accepted. These are mechanics
// only: the write-through guarantee basis remains the documented sentence, and
// the no-REPLACE basis remains contrapositive-plus-probe.
func TestH3RenameMechanics(t *testing.T) {
	base := t.TempDir()

	srcDir := filepath.Join(base, "d1")
	dstDir := filepath.Join(base, "d2")
	if err := os.Mkdir(srcDir, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.Rename(srcDir, dstDir); err != nil {
		t.Fatalf("H3 directory move rejected: %v", err)
	}
	t.Log("H3 directory move (os.Rename): accepted")

	fa, err := windows.UTF16PtrFromString(filepath.Join(base, "a.txt"))
	if err != nil {
		t.Fatal(err)
	}
	fb, err := windows.UTF16PtrFromString(filepath.Join(base, "b.txt"))
	if err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(base, "a.txt"), []byte("a"), 0o666); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(base, "b.txt"), []byte("b"), 0o666); err != nil {
		t.Fatal(err)
	}
	if err := windows.MoveFileEx(fa, fb, 0); err == nil {
		t.Fatalf("H3 no-REPLACE move onto existing destination unexpectedly succeeded")
	} else {
		t.Logf("H3 MoveFileEx flags=0 onto existing destination: rejected err=%v", err)
	}
	if err := windows.MoveFileEx(fa, fb, windows.MOVEFILE_WRITE_THROUGH|windows.MOVEFILE_REPLACE_EXISTING); err != nil {
		t.Fatalf("H3 MOVEFILE_WRITE_THROUGH|REPLACE_EXISTING unexpectedly failed: %v", err)
	}
	t.Log("H3 MoveFileEx WRITE_THROUGH|REPLACE_EXISTING: accepted")
}

// TestH2FlushFileBuffersOnReadOnlyDirHandle records the raw error of
// FlushFileBuffers on a read-only (GENERIC_READ, no FILE_SHARE_DELETE-exclusive)
// directory handle. Diagnostic-only forever: no outcome can enable a mechanism;
// the §B refusal-branch table stands regardless.
func TestH2FlushFileBuffersOnReadOnlyDirHandle(t *testing.T) {
	dir := t.TempDir()
	p, err := windows.UTF16PtrFromString(dir)
	if err != nil {
		t.Fatal(err)
	}
	h, err := windows.CreateFile(p, windows.GENERIC_READ, windows.FILE_SHARE_READ|windows.FILE_SHARE_WRITE|windows.FILE_SHARE_DELETE, nil, windows.OPEN_EXISTING, windows.FILE_FLAG_BACKUP_SEMANTICS, 0)
	if err != nil {
		t.Fatalf("H2 CreateFile(dir, GENERIC_READ, FILE_FLAG_BACKUP_SEMANTICS): %v", err)
	}
	defer windows.CloseHandle(h)
	err = windows.FlushFileBuffers(h)
	var errno windows.Errno
	if errors.As(err, &errno) {
		t.Logf("H2 FlushFileBuffers(read-only dir handle) raw error=%v errno=%d", err, errno)
	} else {
		t.Logf("H2 FlushFileBuffers(read-only dir handle) raw error=%v (non-errno)", err)
	}
}

// TestH4ToolchainAndReadDirErrorClass records runtime.Version() at the
// CI-resolved toolchain and the os.ReadDir(regular file) error class. The §D.5
// fix ships under every outcome; this probe only records the mechanism and
// selects the sweep's scope.
func TestH4ToolchainAndReadDirErrorClass(t *testing.T) {
	t.Logf("H4 runtime.Version()=%s GOOS=%s GOARCH=%s", runtime.Version(), runtime.GOOS, runtime.GOARCH)

	base := t.TempDir()
	file := filepath.Join(base, "regular.txt")
	if err := os.WriteFile(file, []byte("x"), 0o666); err != nil {
		t.Fatal(err)
	}
	_, err := os.ReadDir(file)
	t.Logf("H4 os.ReadDir(regular file) err=%v type=%T os.IsNotExist=%v errors.Is(fs.ErrNotExist)=%v",
		err, err, os.IsNotExist(err), errors.Is(err, fs.ErrNotExist))
}

// TestH1ACLRoundTripAndParentDumps seeds the §D.1 mechanism (owner-only
// protected DACL via ACLFromEntries/SetNamedSecurityInfo in the already-direct
// golang.org/x/sys) and verifies it round-trips through GetNamedSecurityInfo,
// then dumps the parent's and %TEMP%'s effective ACLs (fixture-placement data
// only — never the owner-only policy).
func TestH1ACLRoundTripAndParentDumps(t *testing.T) {
	base := t.TempDir()
	store := filepath.Join(base, "store")
	if err := os.Mkdir(store, 0o700); err != nil {
		t.Fatal(err)
	}

	token := windows.GetCurrentProcessToken()
	tu, err := token.GetTokenUser()
	if err != nil {
		t.Fatalf("H1 GetTokenUser: %v", err)
	}
	userSID := tu.User.Sid

	entries := []windows.EXPLICIT_ACCESS{{
		AccessPermissions: windows.GENERIC_ALL,
		AccessMode:        windows.GRANT_ACCESS,
		Inheritance:       windows.NO_INHERITANCE,
		Trustee: windows.TRUSTEE{
			TrusteeForm:  windows.TRUSTEE_IS_SID,
			TrusteeType:  windows.TRUSTEE_IS_USER,
			TrusteeValue: windows.TrusteeValueFromSID(userSID),
		},
	}}
	acl, err := windows.ACLFromEntries(entries, nil)
	if err != nil {
		t.Fatalf("H1 ACLFromEntries: %v", err)
	}
	if err := windows.SetNamedSecurityInfo(store, windows.SE_FILE_OBJECT, windows.DACL_SECURITY_INFORMATION|windows.PROTECTED_DACL_SECURITY_INFORMATION, nil, nil, acl, nil); err != nil {
		t.Fatalf("H1 SetNamedSecurityInfo(PROTECTED_DACL): %v", err)
	}
	t.Logf("H1 owner SID=%s", userSID.String())

	sd, err := windows.GetNamedSecurityInfo(store, windows.SE_FILE_OBJECT, windows.DACL_SECURITY_INFORMATION)
	if err != nil {
		t.Fatalf("H1 GetNamedSecurityInfo round-trip: %v", err)
	}
	dacl, _, err := sd.DACL()
	if err != nil {
		t.Fatalf("H1 sd.DACL(): %v", err)
	}
	if dacl == nil {
		t.Fatalf("H1 round-trip DACL is nil (SDDL=%s)", sd.String())
	}
	control, _, err := sd.Control()
	if err != nil {
		t.Fatalf("H1 sd.Control(): %v", err)
	}
	if control&windows.SE_DACL_PROTECTED == 0 {
		t.Errorf("H1 round-trip DACL not PROTECTED (control=%#x SDDL=%s)", control, sd.String())
	}

	// Walk the ACEs: expect exactly one non-inherited ACCESS_ALLOWED ACE for
	// the token user, refusing any other grant (naming the trustee).
	const aclHeaderSize = uintptr(8) // revision, sbz1, aclSize, aceCount, sbz2
	off := aclHeaderSize
	for i := uint16(0); i < dacl.AceCount; i++ {
		hdr := (*windows.ACE_HEADER)(unsafe.Pointer(uintptr(unsafe.Pointer(dacl)) + off))
		if hdr.AceType == windows.ACCESS_ALLOWED_ACE_TYPE {
			ace := (*windows.ACCESS_ALLOWED_ACE)(unsafe.Pointer(uintptr(unsafe.Pointer(dacl)) + off))
			sid := (*windows.SID)(unsafe.Pointer(uintptr(unsafe.Pointer(dacl)) + off + unsafe.Offsetof(ace.SidStart)))
			inherited := hdr.AceFlags&windows.INHERITED_ACE != 0
			t.Logf("H1 ACE[%d] ALLOW mask=%#x sid=%s inherited=%v flags=%#x", i, ace.Mask, sid.String(), inherited, hdr.AceFlags)
			if sid.String() != userSID.String() {
				t.Errorf("H1 round-trip grants non-owner trustee %s (want owner %s)", sid.String(), userSID.String())
			}
			if inherited {
				t.Errorf("H1 round-trip owner ACE carries INHERITED_ACE")
			}
		} else {
			t.Logf("H1 ACE[%d] type=%d flags=%#x (non-allow)", i, hdr.AceType, hdr.AceFlags)
			t.Errorf("H1 round-trip DACL has a non-allow ACE (type=%d)", hdr.AceType)
		}
		off += uintptr(hdr.AceSize)
	}
	if dacl.AceCount != 1 {
		t.Errorf("H1 round-trip AceCount=%d, want exactly 1 (SDDL=%s)", dacl.AceCount, sd.String())
	}

	// Parent and %TEMP% effective-ACL dumps: fixture placement only.
	for _, p := range []string{base, os.TempDir()} {
		if psd, err := windows.GetNamedSecurityInfo(p, windows.SE_FILE_OBJECT, windows.DACL_SECURITY_INFORMATION); err != nil {
			t.Logf("H1 dump %s: GetNamedSecurityInfo err=%v", p, err)
		} else {
			pc, _, _ := psd.Control()
			t.Logf("H1 dump %s SDDL=%s protected_dacl=%v", p, psd.String(), pc&windows.SE_DACL_PROTECTED != 0)
		}
	}
}

// TestH6GitConfigShowOriginDump records the runner's effective git config with
// provenance (per-invocation pinning layer confirmation for §D.7) and a byte
// sniff of go.mod for CRLF. Record-only; the split rule is non-negotiable under
// every outcome.
func TestH6GitConfigShowOriginDump(t *testing.T) {
	git, err := exec.LookPath("git")
	if err != nil {
		t.Fatalf("H6 git not on PATH (hosted runners ship Git for Windows): %v", err)
	}
	out, err := exec.Command(git, "config", "--list", "--show-origin", "--show-scope").CombinedOutput()
	if err != nil {
		t.Fatalf("H6 git config --list --show-origin: %v\n%s", err, out)
	}
	t.Logf("H6 git config --list --show-origin --show-scope:\n%s", bytes.TrimRight(out, "\n"))
	if home, err := os.UserHomeDir(); err == nil {
		t.Logf("H6 UserHomeDir=%s", home)
	}

	wd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	if data, err := os.ReadFile(filepath.Join(wd, "..", "..", "go.mod")); err == nil {
		n := 64
		if len(data) < n {
			n = len(data)
		}
		hasCRLF := bytes.Contains(data[:n], []byte("\r\n"))
		t.Logf("H6 go.mod first %d bytes: % x (contains CRLF=%v)", n, data[:n], hasCRLF)
	}
	if _, err := os.Stat(filepath.Join(wd, "..", "..", ".gitattributes")); err != nil {
		t.Logf("H6 .gitattributes absent: %v", err)
	} else {
		t.Log("H6 .gitattributes present")
	}
}

// TestH7ReadOnlyAttributeRenameOver pins the §G H7 / §D.6 clear-before-replace
// mechanics: a FILE_ATTRIBUTE_READONLY destination blocks
// MOVEFILE_REPLACE_EXISTING, and clearing the attribute unblocks it. The
// source-readonly direction is recorded, not asserted (mechanics recording
// only).
func TestH7ReadOnlyAttributeRenameOver(t *testing.T) {
	base := t.TempDir()
	src := filepath.Join(base, "src.txt")
	dst := filepath.Join(base, "dst.txt")
	if err := os.WriteFile(src, []byte("s"), 0o666); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(dst, []byte("d"), 0o666); err != nil {
		t.Fatal(err)
	}

	srcP, err := windows.UTF16PtrFromString(src)
	if err != nil {
		t.Fatal(err)
	}
	dstP, err := windows.UTF16PtrFromString(dst)
	if err != nil {
		t.Fatal(err)
	}

	if err := windows.SetFileAttributes(dstP, windows.FILE_ATTRIBUTE_READONLY); err != nil {
		t.Fatalf("H7 SetFileAttributes(dst, READONLY): %v", err)
	}
	if err := windows.MoveFileEx(srcP, dstP, windows.MOVEFILE_REPLACE_EXISTING); err == nil {
		t.Errorf("H7 rename-over READONLY destination unexpectedly succeeded — clear-before-replace premise needs relabelling")
	} else {
		t.Logf("H7 rename-over READONLY destination: rejected err=%v", err)
	}
	if err := windows.SetFileAttributes(dstP, windows.FILE_ATTRIBUTE_NORMAL); err != nil {
		t.Fatalf("H7 SetFileAttributes(dst, NORMAL): %v", err)
	}
	if err := windows.MoveFileEx(srcP, dstP, windows.MOVEFILE_REPLACE_EXISTING); err != nil {
		t.Errorf("H7 rename-over after clearing READONLY still failed: %v", err)
	} else {
		t.Log("H7 rename-over after clearing READONLY: accepted")
	}

	// Source-readonly direction: record only.
	if err := os.WriteFile(src, []byte("s2"), 0o666); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(dst, []byte("d2"), 0o666); err != nil {
		t.Fatal(err)
	}
	if err := windows.SetFileAttributes(srcP, windows.FILE_ATTRIBUTE_READONLY); err != nil {
		t.Fatalf("H7 SetFileAttributes(src, READONLY): %v", err)
	}
	err = windows.MoveFileEx(srcP, dstP, windows.MOVEFILE_REPLACE_EXISTING)
	t.Logf("H7 readonly SOURCE rename-over (record only): err=%v", err)
}
