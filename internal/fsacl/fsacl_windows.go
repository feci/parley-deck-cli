//go:build windows

package fsacl

import (
	"errors"
	"fmt"
	"os"
	"syscall"
	"unsafe"

	"golang.org/x/sys/windows"
)

// repairInstruction is the §D.1 refuse-and-instruct text for pre-existing
// stores: the user moves the store aside and re-runs so it is re-created
// under the protected DACL. No in-place DACL rewrite of user data happens.
const repairInstruction = "existing snapshot store left untouched; to repair, move it aside and re-run the same command so the store is re-created with an owner-only access policy on an ACL-capable (NTFS) volume"

// EnsurePrivateStore creates dir if missing and enforces the Windows privacy
// contract: an owner-only protected DACL — exactly one access-allowed ACE for
// the token user, PROTECTED_DACL (non-inherited) — set via
// ACLFromEntries/SetNamedSecurityInfo, then re-queried (create-then-requery
// self-check, §D.1). Volumes without ACL support (FAT/exFAT) fail
// SetNamedSecurityInfo and refuse rather than weaken. A store that already
// satisfies the contract is used as-is; a pre-existing store that fails
// verification is user data and is refused with instructions, never rewritten.
func EnsurePrivateStore(dir string) error {
	_, statErr := os.Lstat(dir)
	preexisting := statErr == nil
	if !preexisting && !errors.Is(statErr, os.ErrNotExist) {
		return statErr
	}
	if err := os.MkdirAll(dir, 0o700); err != nil {
		return err
	}
	verifyErr := VerifyPrivateStore(dir)
	if verifyErr == nil {
		return nil
	}
	if !errors.Is(verifyErr, ErrNotPrivate) {
		return verifyErr
	}
	if preexisting {
		return fmt.Errorf("%w (%s)", verifyErr, repairInstruction)
	}
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
