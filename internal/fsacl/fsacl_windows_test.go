//go:build windows

package fsacl

import (
	"errors"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"sync"
	"testing"
	"unsafe"

	"golang.org/x/sys/windows"
)

// Well-known trustees (FINAL §D.1/AC-PRIV-3): canonical SID strings are
// locale-independent on hosted runners, so grants and refusal assertions both
// use the SID form; the refusal names the trustee it finds in the DACL.
const (
	builtinUsersSID         = "S-1-5-32-545" // BUILTIN\Users
	builtinAdministratorsID = "S-1-5-32-544" // BUILTIN\Administrators
)

// grantTrustees replaces the DACL on path with a protected DACL granting
// GENERIC_ALL to the token user plus each extra SID (test scaffolding).
func grantTrustees(t *testing.T, path string, extraSIDs ...string) {
	t.Helper()
	owner, err := tokenUser()
	if err != nil {
		t.Fatal(err)
	}
	entries := []windows.EXPLICIT_ACCESS{ownerEntry(owner)}
	for _, sidStr := range extraSIDs {
		sid, err := windows.StringToSid(sidStr)
		if err != nil {
			t.Fatalf("StringToSid(%s): %v", sidStr, err)
		}
		entries = append(entries, windows.EXPLICIT_ACCESS{
			AccessPermissions: windows.GENERIC_ALL,
			AccessMode:        windows.GRANT_ACCESS,
			Inheritance:       windows.NO_INHERITANCE,
			Trustee: windows.TRUSTEE{
				TrusteeForm:  windows.TRUSTEE_IS_SID,
				TrusteeType:  windows.TRUSTEE_IS_USER,
				TrusteeValue: windows.TrusteeValueFromSID(sid),
			},
		})
	}
	acl, err := windows.ACLFromEntries(entries, nil)
	if err != nil {
		t.Fatalf("ACLFromEntries: %v", err)
	}
	if err := windows.SetNamedSecurityInfo(path, windows.SE_FILE_OBJECT,
		windows.DACL_SECURITY_INFORMATION|windows.PROTECTED_DACL_SECURITY_INFORMATION,
		nil, nil, acl, nil); err != nil {
		t.Fatalf("SetNamedSecurityInfo: %v", err)
	}
}

func ownerEntry(sid *windows.SID) windows.EXPLICIT_ACCESS {
	return windows.EXPLICIT_ACCESS{
		AccessPermissions: windows.GENERIC_ALL,
		AccessMode:        windows.GRANT_ACCESS,
		Inheritance:       windows.NO_INHERITANCE,
		Trustee: windows.TRUSTEE{
			TrusteeForm:  windows.TRUSTEE_IS_SID,
			TrusteeType:  windows.TRUSTEE_IS_USER,
			TrusteeValue: windows.TrusteeValueFromSID(sid),
		},
	}
}

// walkDACL is the adversarial requery used independently of the product's
// verify path: it reports the protection flag and every ACE (type, flags,
// trustee) straight from GetNamedSecurityInfo.
func walkDACL(t *testing.T, path string) (protected bool, aces []string) {
	t.Helper()
	sd, err := windows.GetNamedSecurityInfo(path, windows.SE_FILE_OBJECT, windows.DACL_SECURITY_INFORMATION)
	if err != nil {
		t.Fatalf("GetNamedSecurityInfo: %v", err)
	}
	dacl, _, err := sd.DACL()
	if err != nil {
		t.Fatalf("DACL: %v", err)
	}
	control, _, err := sd.Control()
	if err != nil {
		t.Fatalf("Control: %v", err)
	}
	protected = control&windows.SE_DACL_PROTECTED != 0
	const aclHeaderSize = uintptr(8)
	off := aclHeaderSize
	for i := uint16(0); i < dacl.AceCount; i++ {
		hdr := (*windows.ACE_HEADER)(unsafe.Pointer(uintptr(unsafe.Pointer(dacl)) + off))
		aceType := hdr.AceType
		inherited := hdr.AceFlags&windows.INHERITED_ACE != 0
		trustee := "?"
		if aceType == windows.ACCESS_ALLOWED_ACE_TYPE {
			ace := (*windows.ACCESS_ALLOWED_ACE)(unsafe.Pointer(uintptr(unsafe.Pointer(dacl)) + off))
			sid := (*windows.SID)(unsafe.Pointer(uintptr(unsafe.Pointer(dacl)) + off + unsafe.Offsetof(ace.SidStart)))
			trustee = sid.String()
		}
		aces = append(aces, trustee)
		t.Logf("ACE type=%d inherited=%v trustee=%s", aceType, inherited, trustee)
		off += uintptr(hdr.AceSize)
	}
	return protected, aces
}

// AC-PRIV-1: own creation passes and the create-then-requery readback shows
// exactly one access-allowed ACE for the token user, protected, non-inherited.
func TestOwnerOnlyStoreRoundTrip(t *testing.T) {
	store := filepath.Join(t.TempDir(), "store")
	if err := EnsurePrivateStore(store); err != nil {
		t.Fatalf("EnsurePrivateStore: %v", err)
	}
	owner, err := tokenUserSID()
	if err != nil {
		t.Fatal(err)
	}
	protected, aces := walkDACL(t, store)
	if !protected {
		t.Fatal("DACL not protected after own creation")
	}
	if len(aces) != 1 {
		t.Fatalf("got %d ACEs, want exactly the owner allow: %v", len(aces), aces)
	}
	if aces[0] != owner {
		t.Fatalf("ACE trustee %s, want token user %s", aces[0], owner)
	}
	if err := VerifyPrivateStore(store); err != nil {
		t.Fatalf("VerifyPrivateStore after own creation: %v", err)
	}
}

// AC-PRIV-3: BUILTIN\Users grant → refuse, naming the trustee.
func TestVerifyRefusesBuiltinUsersGrant(t *testing.T) {
	store := filepath.Join(t.TempDir(), "store")
	if err := os.MkdirAll(store, 0o700); err != nil {
		t.Fatal(err)
	}
	grantTrustees(t, store, builtinUsersSID)
	err := VerifyPrivateStore(store)
	if err == nil {
		t.Fatal("BUILTIN\\Users grant accepted")
	}
	if !errors.Is(err, ErrNotPrivate) {
		t.Fatalf("refusal does not wrap ErrNotPrivate: %v", err)
	}
	if !strings.HasPrefix(err.Error(), ErrNotPrivate.Error()) {
		t.Fatalf("refusal lost the sentence-stable prefix: %v", err)
	}
	if !strings.Contains(err.Error(), builtinUsersSID) {
		t.Fatalf("refusal does not name the trustee: %v", err)
	}
}

// AC-PRIV-3: Administrators grant → refuse, naming the trustee.
func TestVerifyRefusesAdministratorsGrant(t *testing.T) {
	store := filepath.Join(t.TempDir(), "store")
	if err := os.MkdirAll(store, 0o700); err != nil {
		t.Fatal(err)
	}
	grantTrustees(t, store, builtinAdministratorsID)
	err := VerifyPrivateStore(store)
	if err == nil {
		t.Fatal("Administrators grant accepted")
	}
	if !errors.Is(err, ErrNotPrivate) || !strings.Contains(err.Error(), builtinAdministratorsID) {
		t.Fatalf("refusal must wrap ErrNotPrivate and name the trustee: %v", err)
	}
}

// AC-PRIV-3: inherited-only → refuse. The parent carries an explicit
// inheritable owner ACE (not protected), so the child's DACL is inherited-only.
func TestVerifyRefusesInheritedOnlyStore(t *testing.T) {
	base := t.TempDir()
	parent := filepath.Join(base, "parent")
	if err := os.MkdirAll(parent, 0o700); err != nil {
		t.Fatal(err)
	}
	owner, err := tokenUser()
	if err != nil {
		t.Fatal(err)
	}
	inheritable := ownerEntry(owner)
	inheritable.Inheritance = windows.SUB_CONTAINERS_AND_OBJECTS_INHERIT
	acl, err := windows.ACLFromEntries([]windows.EXPLICIT_ACCESS{inheritable}, nil)
	if err != nil {
		t.Fatal(err)
	}
	// Deliberately NOT protected: the ACE must flow down to the child.
	if err := windows.SetNamedSecurityInfo(parent, windows.SE_FILE_OBJECT,
		windows.DACL_SECURITY_INFORMATION, nil, nil, acl, nil); err != nil {
		t.Fatal(err)
	}
	child := filepath.Join(parent, "child")
	if err := os.MkdirAll(child, 0o700); err != nil {
		t.Fatal(err)
	}
	err = VerifyPrivateStore(child)
	if err == nil {
		t.Fatal("inherited-only store accepted")
	}
	if !errors.Is(err, ErrNotPrivate) {
		t.Fatalf("refusal does not wrap ErrNotPrivate: %v", err)
	}
}

// AC-PRIV-5: a pre-existing store that fails verification is refused with the
// documented user-invoked repair, and its DACL is NOT rewritten in place.
func TestPreexistingStoreRefusedAndNotRewritten(t *testing.T) {
	store := filepath.Join(t.TempDir(), "store")
	if err := os.MkdirAll(store, 0o700); err != nil {
		t.Fatal(err)
	}
	marker := filepath.Join(store, "existing-archive.tar")
	if err := os.WriteFile(marker, []byte("user data"), 0o600); err != nil {
		t.Fatal(err)
	}
	// Row 26: grantTrustees' protected-DACL replacement on the store
	// auto-propagates inheritance removal to children, which strips a marker
	// carrying only inherited access. Give the marker an EXPLICIT owner allow
	// so its readability measures what the PRODUCT did, not the scaffolding.
	grantTrustees(t, marker)
	grantTrustees(t, store, builtinUsersSID) // permissive pre-existing policy
	protectedBefore, acesBefore := walkDACL(t, store)

	err := EnsurePrivateStore(store)
	if err == nil {
		t.Fatal("permissive pre-existing store accepted")
	}
	if !errors.Is(err, ErrNotPrivate) {
		t.Fatalf("refusal does not wrap ErrNotPrivate: %v", err)
	}
	msg := err.Error()
	if !strings.HasPrefix(msg, ErrNotPrivate.Error()) {
		t.Fatalf("refusal lost the sentence-stable prefix: %q", msg)
	}
	if !strings.Contains(msg, builtinUsersSID) {
		t.Fatalf("refusal does not name the trustee: %q", msg)
	}
	for _, want := range []string{"move it aside", "re-run"} {
		if !strings.Contains(msg, want) {
			t.Fatalf("refusal does not instruct repair (%q missing): %q", want, msg)
		}
	}

	// No in-place rewrite: the DACL is exactly the state the test scaffolding
	// itself installed (grantTrustees deliberately writes a PROTECTED DACL
	// granting Users — protection and permissiveness are orthogonal), and the
	// user data is intact.
	protected, aces := walkDACL(t, store)
	if protected != protectedBefore || !slices.Equal(aces, acesBefore) {
		t.Fatalf("pre-existing store DACL was rewritten (before: protected=%v aces=%v; after: protected=%v aces=%v)",
			protectedBefore, acesBefore, protected, aces)
	}
	grantsUsers := slices.Contains(aces, builtinUsersSID)
	if !grantsUsers {
		t.Fatalf("pre-existing store lost its permissive grant: %v", aces)
	}
	if data, err := os.ReadFile(marker); err != nil || string(data) != "user data" {
		t.Fatalf("user data disturbed: err=%v data=%q", err, data)
	}
}

// AC-PRIV-6: !IsDir() and reparse/symlink rejection also hold on Windows. If
// symlink creation itself fails on the runner, that is a reportable hosted
// fact (fatalf carries the raw error), never a skip.
func TestVerifyRefusesSymlinkAndNonDirectoryStore(t *testing.T) {
	base := t.TempDir()
	file := filepath.Join(base, "file")
	if err := os.WriteFile(file, []byte("x"), 0o600); err != nil {
		t.Fatal(err)
	}
	if err := VerifyPrivateStore(file); !errors.Is(err, ErrNotPrivate) {
		t.Fatalf("non-directory store: err=%v", err)
	}
	real := filepath.Join(base, "real")
	if err := os.Mkdir(real, 0o700); err != nil {
		t.Fatal(err)
	}
	link := filepath.Join(base, "link")
	if err := os.Symlink(real, link); err != nil {
		t.Fatalf("os.Symlink on hosted runner: %v", err)
	}
	if err := VerifyPrivateStore(link); !errors.Is(err, ErrNotPrivate) {
		t.Fatalf("symlink store: err=%v", err)
	}
}

// File half of the policy (archive temp file + read-back): own protection
// passes verification; a granted trustee makes verification refuse.
func TestPrivateFilePolicyRoundTrip(t *testing.T) {
	dir := t.TempDir()
	file := filepath.Join(dir, "archive.tar")
	if err := os.WriteFile(file, []byte("archive"), 0o600); err != nil {
		t.Fatal(err)
	}
	if err := ProtectPrivateFile(file); err != nil {
		t.Fatalf("ProtectPrivateFile: %v", err)
	}
	info, err := os.Lstat(file)
	if err != nil {
		t.Fatal(err)
	}
	if err := VerifyPrivateFile(file, info); err != nil {
		t.Fatalf("VerifyPrivateFile after protection: %v", err)
	}
	grantTrustees(t, file, builtinUsersSID)
	info, err = os.Lstat(file)
	if err != nil {
		t.Fatal(err)
	}
	err = VerifyPrivateFile(file, info)
	if !errors.Is(err, ErrNotPrivate) || !strings.Contains(err.Error(), builtinUsersSID) {
		t.Fatalf("granted file must refuse naming the trustee: %v", err)
	}
}

// Ledger row 24 / §D.1: concurrent first creators race on CreateDirectory.
// The old create→set-DACL window let a loser misclassify the winner's
// mid-creation store as pre-existing and refuse it; with atomic creation the
// loser must verify the winner's policy-complete store and succeed.
func TestConcurrentFirstCreationNeverRefuses(t *testing.T) {
	base := t.TempDir()
	store := filepath.Join(base, "store")
	const creators = 8
	errs := make([]error, creators)
	var start sync.WaitGroup
	start.Add(1)
	var done sync.WaitGroup
	for i := 0; i < creators; i++ {
		done.Add(1)
		go func(i int) {
			defer done.Done()
			start.Wait() // maximize the number of racers on the same missing dir
			errs[i] = EnsurePrivateStore(store)
		}(i)
	}
	start.Done()
	done.Wait()
	for i, err := range errs {
		if err != nil {
			t.Fatalf("concurrent first creator %d refused: %v", i, err)
		}
	}
	if err := VerifyPrivateStore(store); err != nil {
		t.Fatalf("store after concurrent creation: %v", err)
	}
}

// Restore shape (§D.1 creation policy): a product-created MkdirTemp dir has an
// inheriting default DACL — the store guard would (correctly) refuse it as
// foreign — while ProtectPrivateStore applies the owner-only policy to the
// product-owned dir, after which verification passes.
func TestProtectPrivateStoreAppliesPolicyToProductCreatedDir(t *testing.T) {
	parent := t.TempDir()
	dir, err := os.MkdirTemp(parent, "trajectory-source-")
	if err != nil {
		t.Fatal(err)
	}
	if err := VerifyPrivateStore(dir); !errors.Is(err, ErrNotPrivate) {
		t.Fatalf("fresh MkdirTemp dir unexpectedly private (guard premise broken): %v", err)
	}
	if err := ProtectPrivateStore(dir); err != nil {
		t.Fatalf("ProtectPrivateStore: %v", err)
	}
	if err := VerifyPrivateStore(dir); err != nil {
		t.Fatalf("VerifyPrivateStore after protection: %v", err)
	}
}
