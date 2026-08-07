//go:build windows

package protocolcore

// Windows has no O_NOFOLLOW. O_EXCL alone still refuses to create over an existing name, and
// Publish separately refuses when the release directory exists, so the write-once guarantee holds;
// only the symlink-specific hardening is unavailable.
const noFollow = 0
