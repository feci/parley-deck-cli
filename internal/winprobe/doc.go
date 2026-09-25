// Package winprobe holds the FINAL §G early hosted probe bundle
// (windows-portability idea): H1 ACL round-trip and parent dumps, H2
// FlushFileBuffers on a read-only directory handle, H3 rename mechanics and the
// NTFS volume assertion, H4 toolchain pin and ReadDir-on-regular-file error
// class, H6 git config provenance dump, and H7 read-only-attribute
// rename-over. It contains no product code: the probes are test-only
// diagnostics whose outcome→behavior mappings were fixed in FINAL §G before
// execution — mechanics only, never durability, never acceptance evidence, and
// no probe result upgrades an inference.
package winprobe
