// Package fsacl implements the FINAL §D.1 snapshot-store privacy contract for
// the windows-portability idea: on Unix the historical mode-based semantics,
// byte-identical to the previous inline guard; on Windows an owner-only
// protected DACL (exactly one access-allowed ACE for the token user in a
// non-inherited DACL) via the already-direct golang.org/x/sys — zero new
// dependencies. Verification walks effective ACEs and refuses any non-owner
// grant, inherited or explicit, naming the trustee. FAT/exFAT/no-ACL volumes
// refuse, never weaken. It also hosts the shared denyRead test helper (§D.9).
package fsacl

import "errors"

// ErrNotPrivate is the sentence-stable refusal (the hosted 71x signature).
// Windows refusals wrap it with the offending trustee appended.
var ErrNotPrivate = errors.New("snapshot store must be a private real directory")
