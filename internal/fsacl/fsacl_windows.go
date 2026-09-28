//go:build windows

package fsacl

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"syscall"
	"unsafe"

	"golang.org/x/sys/windows"
)

// repairInstruction is the §D.1 refuse-and-instruct text for pre-existing
// stores: the user moves the store aside and re-runs so it is re-created
// under the protected DACL. No in-place DACL rewrite of user data happens.
const repairInstruction = "existing snapshot store left untouched; to repair, move it aside and re-run the same command so the store is re-created with an owner-only access policy on an ACL-capable (NTFS) volume"

// EnsurePrivateStore is the product boundary guard for a snapshot store: an
// owner-only protected DACL — exactly one access-allowed ACE for the token
// user, PROTECTED_DACL (non-inherited). A store that already satisfies the
// contract is used as-is; a pre-existing store that fails verification is user
// data and is refused with instructions, never rewritten. A missing store is
// created ATOMICALLY with its policy: CreateDirectory carries the owner-only
// security descriptor in SECURITY_ATTRIBUTES, so no concurrent first capture
// can ever observe the store without its owner-only DACL (the ledger row-24
// create→set-DACL window); ERROR_ALREADY_EXISTS means a concurrent creator
// won — and its directory was created under the same atomic policy — so the
// loser verifies rather than refuses. Volumes without ACL support (FAT/exFAT)
// fail creation-with-policy or the read-back and refuse rather than weaken.
func EnsurePrivateStore(dir string) error {
	_, statErr := os.Lstat(dir)
	switch {
	case statErr == nil:
		return verifyOrRefuseStore(dir)
	case !errors.Is(statErr, os.ErrNotExist):
		return statErr
	}
	if err := os.MkdirAll(filepath.Dir(dir), 0o700); err != nil {
		return err
	}
	sd, err := ownerOnlySD()
	if err != nil {
		return err
	}
	sa := &windows.SecurityAttributes{
		Length:             uint32(unsafe.Sizeof(windows.SecurityAttributes{})),
		SecurityDescriptor: sd,
	}
	path, err := syscall.UTF16PtrFromString(dir)
	if err != nil {
		return err
	}
	if err := windows.CreateDirectory(path, sa); err != nil {
		if errors.Is(err, windows.ERROR_ALREADY_EXISTS) {
			return verifyOrRefuseStore(dir)
		}
		return fmt.Errorf("%w: creating store directory with owner-only DACL: %w", ErrNotPrivate, err)
	}
	return VerifyPrivateStore(dir)
}

// verifyOrRefuseStore accepts a store that passes verification and turns any
// privacy failure on a present directory into the refuse-and-instruct refusal.
func verifyOrRefuseStore(dir string) error {
	err := VerifyPrivateStore(dir)
	if err == nil {
		return nil
	}
	if !errors.Is(err, ErrNotPrivate) {
		return err
	}
	return fmt.Errorf("%w (%s)", err, repairInstruction)
}

// ProtectPrivateStore applies the §D.1 creation policy to a directory the
// product itself just created through a non-atomic allocator (restore's
// os.MkdirTemp): the owner-only protected DACL is set on the product-owned
// fresh directory, then verified. There is deliberately no refuse-and-instruct
// path — the caller guarantees product creation — which is what distinguishes
// creation policy from the EnsurePrivateStore guard.
func ProtectPrivateStore(dir string) error {
	acl, err := ownerOnlyDACL()
	if err != nil {
		return err
	}
	if err := windows.SetNamedSecurityInfo(dir, windows.SE_FILE_OBJECT,
		windows.DACL_SECURITY_INFORMATION|windows.PROTECTED_DACL_SECURITY_INFORMATION,
		nil, nil, acl, nil); err != nil {
		return fmt.Errorf("%w: setting owner-only DACL: %w", ErrNotPrivate, err)
	}
	return VerifyPrivateStore(dir)
}

// VerifyPrivateStore walks the effective DACL and refuses any non-owner grant,
// inherited or explicit, naming the trustee (§D.1: the refusal is
// sentence-stable with the trustee appended).
func VerifyPrivateStore(dir string) error {
	info, err := os.Lstat(dir)
	if err != nil {
		return err
	}
	if !info.IsDir() || info.Mode()&os.ModeSymlink != 0 || hasReparsePoint(info) {
		return ErrNotPrivate
	}
	return verifyOwnerOnlyDACL(dir)
}

// ProtectPrivateFile applies the owner-only policy to a freshly created file
// (the snapshot archive temp file). The file is product-owned, so setting its
// DACL is creation policy, not a rewrite of user data.
func ProtectPrivateFile(file string) error {
	acl, err := ownerOnlyDACL()
	if err != nil {
		return err
	}
	if err := windows.SetNamedSecurityInfo(file, windows.SE_FILE_OBJECT,
		windows.DACL_SECURITY_INFORMATION|windows.PROTECTED_DACL_SECURITY_INFORMATION,
		nil, nil, acl, nil); err != nil {
		return fmt.Errorf("%w: setting owner-only file DACL: %w", ErrNotPrivate, err)
	}
	return nil
}

// VerifyPrivateFile checks the §D.1 policy on an archive file; info is the
// caller's Lstat result (regularity is checked by the caller).
func VerifyPrivateFile(file string, info os.FileInfo) error {
	if info.Mode()&os.ModeSymlink != 0 || hasReparsePoint(info) || !info.Mode().IsRegular() {
		return ErrNotPrivate
	}
	return verifyOwnerOnlyDACL(file)
}

func hasReparsePoint(info os.FileInfo) bool {
	sys, ok := info.Sys().(*syscall.Win32FileAttributeData)
	return ok && sys.FileAttributes&windows.FILE_ATTRIBUTE_REPARSE_POINT != 0
}

func verifyOwnerOnlyDACL(path string) error {
	sd, err := windows.GetNamedSecurityInfo(path, windows.SE_FILE_OBJECT, windows.DACL_SECURITY_INFORMATION)
	if err != nil {
		return fmt.Errorf("%w: reading DACL: %w", ErrNotPrivate, err)
	}
	dacl, _, err := sd.DACL()
	if err != nil {
		return fmt.Errorf("%w: reading DACL: %w", ErrNotPrivate, err)
	}
	if dacl == nil {
		return fmt.Errorf("%w: no DACL", ErrNotPrivate)
	}
	control, _, err := sd.Control()
	if err != nil {
		return fmt.Errorf("%w: reading control: %w", ErrNotPrivate, err)
	}
	if control&windows.SE_DACL_PROTECTED == 0 {
		return fmt.Errorf("%w: DACL not protected (inheriting)", ErrNotPrivate)
	}
	owner, err := tokenUserSID()
	if err != nil {
		return err
	}
	const aclHeaderSize = uintptr(8) // revision, sbz1, aclSize, aceCount, sbz2
	off := aclHeaderSize
	for i := uint16(0); i < dacl.AceCount; i++ {
		hdr := (*windows.ACE_HEADER)(unsafe.Pointer(uintptr(unsafe.Pointer(dacl)) + off))
		switch hdr.AceType {
		case windows.ACCESS_ALLOWED_ACE_TYPE:
			ace := (*windows.ACCESS_ALLOWED_ACE)(unsafe.Pointer(uintptr(unsafe.Pointer(dacl)) + off))
			sid := (*windows.SID)(unsafe.Pointer(uintptr(unsafe.Pointer(dacl)) + off + unsafe.Offsetof(ace.SidStart)))
			if sid.String() != owner {
				return fmt.Errorf("%w: grants %s", ErrNotPrivate, sid.String())
			}
		default:
			return fmt.Errorf("%w: unexpected ACE type %d", ErrNotPrivate, hdr.AceType)
		}
		if hdr.AceFlags&windows.INHERITED_ACE != 0 {
			return fmt.Errorf("%w: inherited ACE present", ErrNotPrivate)
		}
		off += uintptr(hdr.AceSize)
	}
	if dacl.AceCount != 1 {
		return fmt.Errorf("%w: %d ACEs, want exactly the owner allow", ErrNotPrivate, dacl.AceCount)
	}
	return nil
}

// DenyRead makes path unreadable for tests: a deny-read ACE for the current
// user merged into the existing DACL (§D.9; shares this module's mechanism).
func DenyRead(path string) error {
	owner, err := tokenUser()
	if err != nil {
		return err
	}
	sd, err := windows.GetNamedSecurityInfo(path, windows.SE_FILE_OBJECT, windows.DACL_SECURITY_INFORMATION)
	if err != nil {
		return err
	}
	existing, _, err := sd.DACL()
	if err != nil {
		return err
	}
	entries := []windows.EXPLICIT_ACCESS{{
		AccessPermissions: windows.GENERIC_READ | windows.GENERIC_EXECUTE,
		AccessMode:        windows.DENY_ACCESS,
		Inheritance:       windows.NO_INHERITANCE,
		Trustee: windows.TRUSTEE{
			TrusteeForm:  windows.TRUSTEE_IS_SID,
			TrusteeType:  windows.TRUSTEE_IS_USER,
			TrusteeValue: windows.TrusteeValueFromSID(owner),
		},
	}}
	acl, err := windows.ACLFromEntries(entries, existing)
	if err != nil {
		return err
	}
	return windows.SetNamedSecurityInfo(path, windows.SE_FILE_OBJECT,
		windows.DACL_SECURITY_INFORMATION, nil, nil, acl, nil)
}

// DenyWrite makes a directory (or file) unwritable for tests: a deny ACE
// for the current user covering GENERIC_WRITE, merged into the existing
// DACL (§D.9; the write-side sibling of DenyRead). GENERIC_WRITE on a
// directory maps to adding files/subdirectories and does NOT include
// FILE_DELETE_CHILD or DELETE, so cleanup removals still work.
func DenyWrite(path string) error {
	// Hosted evidence (444d930's organizer leg): a GENERIC_WRITE deny also
	// blocks directory ENUMERATION (the walk's dir open got Access denied).
	// Narrow to FILE_WRITE_DATA|FILE_APPEND_DATA (add-file/add-subdirectory for directories) — write-into
	// is denied, listing/reading stays possible, so write-detection walks
	// still function.
	return denyAccess(path, windows.FILE_WRITE_DATA|windows.FILE_APPEND_DATA)
}

// DenyWriteTree makes a whole TREE unwritable for tests (claude-1
// readonly-walk consult 4.4): the deny ACE carries
// SUB_CONTAINERS_AND_OBJECTS_INHERIT so it propagates to existing children
// (their DACLs are unprotected) — existing files' bytes genuinely cannot be
// rewritten — while the narrow mask keeps reads and cleanup removals
// possible. Distinct from DenyWrite (one level) so the other call sites'
// semantics are unchanged.
func DenyWriteTree(path string) error {
	return denyAccessTree(path, windows.FILE_WRITE_DATA|windows.FILE_APPEND_DATA, windows.SUB_CONTAINERS_AND_OBJECTS_INHERIT)
}

// AllowWriteTree removes a DenyWriteTree (the same rebuilt-allow restore,
// protection flag preserved as-is so inheritance semantics stay intact).
func AllowWriteTree(path string) error {
	sd, err := windows.GetNamedSecurityInfo(path, windows.SE_FILE_OBJECT, windows.DACL_SECURITY_INFORMATION)
	if err != nil {
		return err
	}
	dacl, _, err := sd.DACL()
	if err != nil {
		return err
	}
	if dacl == nil {
		return nil
	}
	var kept []windows.EXPLICIT_ACCESS
	const aclHeaderSize = uintptr(8)
	off := aclHeaderSize
	for i := uint16(0); i < dacl.AceCount; i++ {
		hdr := (*windows.ACE_HEADER)(unsafe.Pointer(uintptr(unsafe.Pointer(dacl)) + off))
		if hdr.AceType == windows.ACCESS_ALLOWED_ACE_TYPE {
			ace := (*windows.ACCESS_ALLOWED_ACE)(unsafe.Pointer(uintptr(unsafe.Pointer(dacl)) + off))
			sid := (*windows.SID)(unsafe.Pointer(uintptr(unsafe.Pointer(dacl)) + off + unsafe.Offsetof(ace.SidStart)))
			kept = append(kept, windows.EXPLICIT_ACCESS{
				AccessPermissions: ace.Mask,
				AccessMode:        windows.GRANT_ACCESS,
				Inheritance:       windows.SUB_CONTAINERS_AND_OBJECTS_INHERIT,
				Trustee: windows.TRUSTEE{
					TrusteeForm:  windows.TRUSTEE_IS_SID,
					TrusteeType:  windows.TRUSTEE_IS_USER,
					TrusteeValue: windows.TrusteeValueFromSID(sid),
				},
			})
		}
		off += uintptr(hdr.AceSize)
	}
	acl, err := windows.ACLFromEntries(kept, nil)
	if err != nil {
		return err
	}
	return windows.SetNamedSecurityInfo(path, windows.SE_FILE_OBJECT,
		windows.DACL_SECURITY_INFORMATION, nil, nil, acl, nil)
}

func denyAccessTree(path string, mask windows.ACCESS_MASK, inherit uint32) error {
	owner, err := tokenUser()
	if err != nil {
		return err
	}
	sd, err := windows.GetNamedSecurityInfo(path, windows.SE_FILE_OBJECT, windows.DACL_SECURITY_INFORMATION)
	if err != nil {
		return err
	}
	existing, _, err := sd.DACL()
	if err != nil {
		return err
	}
	entries := []windows.EXPLICIT_ACCESS{{
		AccessPermissions: mask,
		AccessMode:        windows.DENY_ACCESS,
		Inheritance:       inherit,
		Trustee: windows.TRUSTEE{
			TrusteeForm:  windows.TRUSTEE_IS_SID,
			TrusteeType:  windows.TRUSTEE_IS_USER,
			TrusteeValue: windows.TrusteeValueFromSID(owner),
		},
	}}
	acl, err := windows.ACLFromEntries(entries, existing)
	if err != nil {
		return err
	}
	return windows.SetNamedSecurityInfo(path, windows.SE_FILE_OBJECT,
		windows.DACL_SECURITY_INFORMATION, nil, nil, acl, nil)
}

// AllowWrite removes a DenyWrite deny ACE (tests restore writability before
// their cleanup). Unix restores mode 0755; Windows drops the deny ACE.
func AllowWrite(path string) error {
	return removeDeny(path, windows.FILE_WRITE_DATA|windows.FILE_APPEND_DATA)
}

func denyAccess(path string, mask windows.ACCESS_MASK) error {
	owner, err := tokenUser()
	if err != nil {
		return err
	}
	sd, err := windows.GetNamedSecurityInfo(path, windows.SE_FILE_OBJECT, windows.DACL_SECURITY_INFORMATION)
	if err != nil {
		return err
	}
	existing, _, err := sd.DACL()
	if err != nil {
		return err
	}
	entries := []windows.EXPLICIT_ACCESS{{
		AccessPermissions: mask,
		AccessMode:        windows.DENY_ACCESS,
		Inheritance:       windows.NO_INHERITANCE,
		Trustee: windows.TRUSTEE{
			TrusteeForm:  windows.TRUSTEE_IS_SID,
			TrusteeType:  windows.TRUSTEE_IS_USER,
			TrusteeValue: windows.TrusteeValueFromSID(owner),
		},
	}}
	acl, err := windows.ACLFromEntries(entries, existing)
	if err != nil {
		return err
	}
	return windows.SetNamedSecurityInfo(path, windows.SE_FILE_OBJECT,
		windows.DACL_SECURITY_INFORMATION, nil, nil, acl, nil)
}

func removeDeny(path string, mask windows.ACCESS_MASK) error {
	sd, err := windows.GetNamedSecurityInfo(path, windows.SE_FILE_OBJECT, windows.DACL_SECURITY_INFORMATION)
	if err != nil {
		return err
	}
	dacl, _, err := sd.DACL()
	if err != nil {
		return err
	}
	if dacl == nil {
		return nil
	}
	var kept []windows.EXPLICIT_ACCESS
	const aclHeaderSize = uintptr(8)
	off := aclHeaderSize
	for i := uint16(0); i < dacl.AceCount; i++ {
		hdr := (*windows.ACE_HEADER)(unsafe.Pointer(uintptr(unsafe.Pointer(dacl)) + off))
		if hdr.AceType == windows.ACCESS_ALLOWED_ACE_TYPE {
			// Keep every allow (inherited or explicit) as an explicit grant;
			// deny ACEs — ours — are dropped by not being carried into kept.
			ace := (*windows.ACCESS_ALLOWED_ACE)(unsafe.Pointer(uintptr(unsafe.Pointer(dacl)) + off))
			sid := (*windows.SID)(unsafe.Pointer(uintptr(unsafe.Pointer(dacl)) + off + unsafe.Offsetof(ace.SidStart)))
			kept = append(kept, windows.EXPLICIT_ACCESS{
				AccessPermissions: ace.Mask,
				AccessMode:        windows.GRANT_ACCESS,
				Inheritance:       windows.NO_INHERITANCE,
				Trustee: windows.TRUSTEE{
					TrusteeForm:  windows.TRUSTEE_IS_SID,
					TrusteeType:  windows.TRUSTEE_IS_USER,
					TrusteeValue: windows.TrusteeValueFromSID(sid),
				},
			})
		}
		off += uintptr(hdr.AceSize)
	}
	acl, err := windows.ACLFromEntries(kept, nil)
	if err != nil {
		return err
	}
	// NO PROTECTED flag: setting it strips the directory's inherited ACEs
	// and bricks cleanup traversal (the hosted 36447708947 regression root
	// cause — the DACL became an empty protected list). Without it the
	// rebuilt grants merge with parent inheritance as before.
	return windows.SetNamedSecurityInfo(path, windows.SE_FILE_OBJECT,
		windows.DACL_SECURITY_INFORMATION,
		nil, nil, acl, nil)
}

func tokenUserSID() (string, error) {
	sid, err := tokenUser()
	if err != nil {
		return "", err
	}
	return sid.String(), nil
}

func tokenUser() (*windows.SID, error) {
	tu, err := windows.GetCurrentProcessToken().GetTokenUser()
	if err != nil {
		return nil, fmt.Errorf("%w: reading process token user: %w", ErrNotPrivate, err)
	}
	return tu.User.Sid, nil
}

func ownerOnlyDACL() (*windows.ACL, error) {
	sid, err := tokenUser()
	if err != nil {
		return nil, err
	}
	entries := []windows.EXPLICIT_ACCESS{{
		AccessPermissions: windows.GENERIC_ALL,
		AccessMode:        windows.GRANT_ACCESS,
		Inheritance:       windows.NO_INHERITANCE,
		Trustee: windows.TRUSTEE{
			TrusteeForm:  windows.TRUSTEE_IS_SID,
			TrusteeType:  windows.TRUSTEE_IS_USER,
			TrusteeValue: windows.TrusteeValueFromSID(sid),
		},
	}}
	return windows.ACLFromEntries(entries, nil)
}

// ownerOnlySD is the atomic-creation form of the same single-ACE policy:
// SDDL "D:P" = a protected (non-inherited) DACL holding one allow-all ACE for
// the token user with no inheritance flags — the DACL ownerOnlyDACL produces,
// carried inside CreateDirectory's SECURITY_ATTRIBUTES so the directory never
// exists without its policy (ledger row 24).
func ownerOnlySD() (*windows.SECURITY_DESCRIPTOR, error) {
	sid, err := tokenUserSID()
	if err != nil {
		return nil, err
	}
	return windows.SecurityDescriptorFromString("D:P(A;;GA;;;" + sid + ")")
}
