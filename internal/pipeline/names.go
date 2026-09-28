package pipeline

import (
	"fmt"
	"net/url"
	"strings"
)

// Gate-file name policy (FINAL §D.4 / AC-NAME-1..4): one universal encoding on
// every OS applied inside GatePath, so the .tmp staging name inherits it by
// construction (the hosted ERROR_INVALID_NAME failures were on the staging
// name). Windows is the binding constraint; per-OS naming is rejected because
// decks cross OSes.

// maxBlockIDLen bounds raw block IDs at creation. Allowlisted characters never
// grow under the encoding, so a from->to edge file stays far inside the 255
// byte component limit even at the bound.
const maxBlockIDLen = 64

// encodeEdgeID encodes an arbitrary edge ID into one filename component that
// is valid on every OS: net/url.PathEscape base (every separator, '>' and '<'
// included, becomes %XX), an explicit ':' escape (PathEscape keeps ':' for
// path segments; on Windows a colon opens an alternate data stream), then the
// reserved-device-name and trailing-dot guards. PathEscape also escapes '%',
// so the encoding is injective: two edge IDs never share a gate file, and a
// component that was guard-escaped (%43ON for CON) cannot be produced by any
// other input. It is total — no input can escape the gates directory, because
// no separator survives PathEscape.
func encodeEdgeID(edgeID string) string {
	enc := url.PathEscape(edgeID)
	enc = strings.ReplaceAll(enc, ":", "%3A")
	if reservedDeviceName(enc) {
		// Reserved even with an extension: escape the first character so the
		// component no longer names a DOS device (CON -> %43ON).
		enc = fmt.Sprintf("%%%02X%s", enc[0], enc[1:])
	}
	// PathEscape already turns a trailing space into %20; a trailing dot is
	// left alone by PathEscape and is itself stripped by some Windows APIs.
	if strings.HasSuffix(enc, ".") {
		enc = enc[:len(enc)-1] + "%2E"
	}
	return enc
}

// reservedDeviceName reports whether component names a DOS device on Windows:
// the part before the first dot, case-insensitively, equal to CON, PRN, AUX,
// NUL, COM1-9 or LPT1-9 ("reserved even with extension").
func reservedDeviceName(component string) bool {
	head, _, _ := strings.Cut(component, ".")
	head = strings.ToUpper(head)
	switch head {
	case "CON", "PRN", "AUX", "NUL",
		"COM1", "COM2", "COM3", "COM4", "COM5", "COM6", "COM7", "COM8", "COM9",
		"LPT1", "LPT2", "LPT3", "LPT4", "LPT5", "LPT6", "LPT7", "LPT8", "LPT9":
		return true
	}
	return false
}

// checkBlockID refuses the unsafe classes on the RAW block ID (FINAL §D.4: the
// executed IsLocal table showed validating the concatenated slug__+id path is
// unsound — the prefix absorbs ".." inside Clean and wrongly accepts
// traversals, landing in another idea's directory). No legacy exemption:
// traversal, separators, absolute/drive paths, ADS colons, reserved device
// names (even with extension) and trailing dot/space are refused on every OS.
// Safe-but-non-allowlisted IDs (spaces, unicode) are NOT refused here — they
// are cosmetic grandfathering and stay encodable.
func checkBlockID(id string) error {
	if id == "" {
		return fmt.Errorf("block id is empty")
	}
	if strings.ContainsAny(id, `/\`) {
		return fmt.Errorf("block id %q contains a path separator", id)
	}
	if strings.Contains(id, ":") {
		return fmt.Errorf("block id %q contains a colon (drive/ADS path)", id)
	}
	if id == "." || id == ".." || strings.HasPrefix(id, "../") || strings.HasPrefix(id, "..\\") {
		return fmt.Errorf("block id %q is a traversal", id)
	}
	if reservedDeviceName(id) {
		return fmt.Errorf("block id %q is a reserved device name", id)
	}
	if strings.HasSuffix(id, ".") || strings.HasSuffix(id, " ") {
		return fmt.Errorf("block id %q ends in a dot or space", id)
	}
	return nil
}

// validNewBlockID is the strict charset allowlist applied to raw IDs at NEW
// pipeline creation (AC-NAME-3): [A-Za-z0-9._-], bounded length, no leading or
// trailing dot, not a reserved device name. Everything else — including the
// unsafe classes checkBlockID refuses — is refused for new pipelines; existing
// decks only grandfather the safe-but-cosmetic remainder at load.
func ValidNewBlockID(id string) error {
	if err := checkBlockID(id); err != nil {
		return err
	}
	if len(id) > maxBlockIDLen {
		return fmt.Errorf("block id %q longer than %d characters", id, maxBlockIDLen)
	}
	for i := 0; i < len(id); i++ {
		c := id[i]
		switch {
		case c >= 'a' && c <= 'z', c >= 'A' && c <= 'Z', c >= '0' && c <= '9':
		case c == '.' || c == '_' || c == '-':
		default:
			return fmt.Errorf("block id %q contains character %q outside the [A-Za-z0-9._-] allowlist", id, string(c))
		}
	}
	if strings.HasPrefix(id, ".") {
		return fmt.Errorf("block id %q starts with a dot", id)
	}
	return nil
}
