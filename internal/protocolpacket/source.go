package protocolpacket

// Authority resolution.
//
// Source-role decks (this repository) are the protocol's upstream: the live deck file IS the
// authority, and the expected drift against an older global core must never substitute stale
// core text (ALT-07). Consumer decks render from the pinned, hash-attested core plus lock plus
// overlay; that chain must verify BEFORE any optimization, and a verified render that differs
// from the on-disk view is drift that cannot pass as a valid optimized packet. Missing or
// unprovable authority is an error, never a bundled snapshot.

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"parley-deck-cli/internal/fsutil"
	"parley-deck-cli/internal/protocol"
	"parley-deck-cli/internal/protocolcore"
)

// Publication permissions. A packet body is protocol text a launch is about to be given: the
// directory is private so another local user cannot substitute a body between publication and
// launch, and the file is owner-read/write only.
const (
	runtimeDirPerm = 0o700
	bodyFilePerm   = 0o600
)

// ErrAuthority wraps every "cannot prove what the protocol is" failure.
var ErrAuthority = errors.New("protocolpacket: authority")

// Source is the resolved authoritative protocol text.
type Source struct {
	Path      string
	Role      string
	Authority string
	Raw       string
	Transport string // parsed from the `**Transport:**` header line
	Drift     string // consumer only: non-empty when the on-disk view is not the verified render
}

// DeckProtocolPath is the deck's protocol view.
func DeckProtocolPath(root string) string {
	return filepath.Join(root, protocol.DeckDir, "COOPERATION.md")
}

// MapPath is the applicability map location.
func MapPath(root string) string {
	return filepath.Join(root, protocol.DeckDir, "meta", "packet-applicability.yaml")
}

// RuntimeDir is where generated bodies go. The path is ignored by git (the ignore entry is
// integration-owned); a packet is never committed because a committed packet is a stale copy.
func RuntimeDir(root string) string {
	return filepath.Join(root, ".parley-runtime", "protocol-packets")
}

// ResolveSource resolves the deck's authoritative protocol text according to its role.
func ResolveSource(root string, store protocolcore.Store) (Source, error) {
	role, err := protocolRole(root)
	if err != nil {
		return Source{}, err
	}
	path := DeckProtocolPath(root)
	raw, err := os.ReadFile(path)
	if err != nil {
		return Source{}, fmt.Errorf("%w: cannot read %s: %v", ErrAuthority, path, err)
	}
	src := Source{Path: path, Role: role, Raw: string(raw), Transport: headerTransport(string(raw))}
	switch role {
	case "source":
		src.Authority = "live source file (protocolRole: source); global core drift is expected and never substituted"
		return src, nil
	case "consumer":
		lockPath := filepath.Join(root, protocol.DeckDir, "meta", "protocol-lock.yaml")
		lb, err := os.ReadFile(lockPath)
		if err != nil {
			return Source{}, fmt.Errorf("%w: consumer deck has no readable %s; the deck view cannot be attested to a core release", ErrAuthority, lockPath)
		}
		lock, err := protocolcore.ParseLock(string(lb))
		if err != nil {
			return Source{}, fmt.Errorf("%w: %v", ErrAuthority, err)
		}
		rel, err := store.Load(lock.CoreVersion)
		if err != nil {
			return Source{}, fmt.Errorf("%w: %w", ErrAuthority, err)
		}
		if rel.SHA256 != lock.CoreBodySHA256 {
			return Source{}, fmt.Errorf("%w: lock attests core %s body %s but the installed release hashes %s",
				ErrAuthority, lock.CoreVersion, protocolcore.ShortHash(lock.CoreBodySHA256), protocolcore.ShortHash(rel.SHA256))
		}
		ovPath := filepath.Join(root, protocol.DeckDir, protocolcore.OverlayFileName)
		var ov *protocolcore.Overlay
		ob, oerr := os.ReadFile(ovPath)
		switch {
		case oerr == nil:
			ov, err = protocolcore.ReconcileOverlay(lock, true, string(ob), nil)
		case errors.Is(oerr, os.ErrNotExist):
			ov, err = protocolcore.ReconcileOverlay(lock, false, "", nil)
		default:
			ov, err = protocolcore.ReconcileOverlay(lock, true, "", oerr)
		}
		if err != nil {
			return Source{}, fmt.Errorf("%w: %v", ErrAuthority, err)
		}
		res, err := protocolcore.Render(rel, src.Raw, ov)
		if err != nil {
			return Source{}, fmt.Errorf("%w: %v", ErrAuthority, err)
		}
		src.Authority = fmt.Sprintf("core %s (%s) + lock + overlay, verified", lock.CoreVersion, protocolcore.ShortHash(rel.SHA256))
		if res.Body != src.Raw {
			src.Drift = "deck view differs from the verified core render (protocol check: hand-edited-or-stale)"
		}
		return src, nil
	default:
		return Source{}, fmt.Errorf("%w: meta/version.json protocolRole %q is unknown; confirm source|consumer (parley preflight)", ErrAuthority, role)
	}
}

func protocolRole(root string) (string, error) {
	path := filepath.Join(root, protocol.DeckDir, "meta", "version.json")
	b, err := os.ReadFile(path)
	if err != nil {
		return "", fmt.Errorf("%w: cannot read %s: %v", ErrAuthority, path, err)
	}
	var meta struct {
		ProtocolRole string `json:"protocolRole"`
	}
	if err := json.Unmarshal(b, &meta); err != nil {
		return "", fmt.Errorf("%w: %s is not valid JSON: %v", ErrAuthority, path, err)
	}
	return strings.ToLower(strings.TrimSpace(meta.ProtocolRole)), nil
}

// headerTransport reads the deck's `**Transport:** `value“ line. An absent or placeholder
// value is returned as-is and rejected by Build as unknown-transport.
func headerTransport(body string) string {
	for _, l := range strings.Split(body, "\n") {
		if !strings.HasPrefix(l, "**Transport:**") {
			continue
		}
		v := strings.TrimSpace(strings.TrimPrefix(l, "**Transport:**"))
		return strings.Trim(v, "` ")
	}
	return ""
}

// WriteBody publishes the emitted body under RuntimeDir as an immutable, content-addressed
// file and records the path on the context. A refused context has no body and nothing is
// written.
//
// Publication is exclusive and never truncates. The earlier form — a name carrying a 12-hex
// short hash, written with os.WriteFile — combined two failures: create-truncate lets a
// concurrent reader observe a partial body under an attestation that names the whole one, and
// a truncated digest lets two different bodies address one path. Here the name carries the
// FULL body digest, the bytes are staged in a sibling temporary file and linked into place, and
// an existing path is accepted only when its bytes are exactly the attested body. A mismatched
// body, a symlink, or a non-regular file is an error: launching from a body that is not the
// attested one is precisely what this renderer exists to prevent.
//
// Every name component is either an int or safeName-sanitised, so a published path cannot
// escape the runtime directory.
func WriteBody(root string, ctx *Context) error {
	if ctx.ContextMode == ModeRefused {
		return fmt.Errorf("protocolpacket: refused context has no body (%s)", ctx.FallbackReason)
	}
	if ctx.Body == "" {
		return errors.New("protocolpacket: empty body is not publishable")
	}
	digest := Hash(ctx.Body)
	if ctx.PacketSHA256 != digest {
		// The attestation is the promise that the launched text is this text. A body that does
		// not hash to it was mutated after Build and must not be published under it.
		return fmt.Errorf("protocolpacket: body hashes %s but the attestation says %s; refusing to publish an unattested body",
			protocolcore.ShortHash(digest), protocolcore.ShortHash(ctx.PacketSHA256))
	}
	dir, err := publicationDir(root)
	if err != nil {
		return err
	}
	name := fmt.Sprintf("%s-phase%d-%s-%s.md", safeName(ctx.ContextMode), ctx.Request.Phase, safeName(ctx.Request.Track), digest)
	path := filepath.Join(dir, name)
	if err := publish(dir, path, ctx.Body); err != nil {
		return err
	}
	ctx.BodyPath = path
	return nil
}

// publicationDir creates or validates the private publication directory and its parent. A
// symlinked or non-directory component, or a directory other local users can write, is
// refused: a publication directory that a third party can rewrite carries no immutability.
func publicationDir(root string) (string, error) {
	dir := RuntimeDir(root)
	if err := ensurePrivateDir(filepath.Dir(dir)); err != nil {
		return "", err
	}
	if err := ensurePrivateDir(dir); err != nil {
		return "", err
	}
	return dir, nil
}

func ensurePrivateDir(path string) error {
	fi, err := os.Lstat(path)
	switch {
	case errors.Is(err, os.ErrNotExist):
		mkErr := os.Mkdir(path, runtimeDirPerm)
		switch {
		case mkErr == nil:
			// Mkdir's mode is masked by umask; set the mode we actually require.
			if err := os.Chmod(path, runtimeDirPerm); err != nil {
				return fmt.Errorf("protocolpacket: cannot restrict %s: %w", path, err)
			}
			return nil
		case errors.Is(mkErr, os.ErrExist):
			// The create race was lost: something appeared between the Lstat and the Mkdir —
			// usually a concurrent publisher, but EEXIST equally describes a symlink or a file
			// planted there. It proves an entry exists, never that the entry is the private
			// directory this call asked for, so what actually got created is re-read and put
			// through the same checks as a pre-existing path instead of being trusted.
			fi, err = os.Lstat(path)
			if err != nil {
				return fmt.Errorf("protocolpacket: %s appeared while it was being created but cannot be stat-ed: %w", path, err)
			}
		default:
			return fmt.Errorf("protocolpacket: cannot create %s: %w", path, mkErr)
		}
	case err != nil:
		return fmt.Errorf("protocolpacket: cannot stat %s: %w", path, err)
	}
	return checkPrivateDir(path, fi)
}

// checkPrivateDir accepts only a real directory that no other local user can write. fi is the
// caller's Lstat of path, so a symlink is judged as a symlink and never as its target.
func checkPrivateDir(path string, fi os.FileInfo) error {
	if fi.Mode()&os.ModeSymlink != 0 {
		return fmt.Errorf("protocolpacket: %s is a symlink; the packet runtime path must be a real directory", path)
	}
	if !fi.IsDir() {
		return fmt.Errorf("protocolpacket: %s exists and is not a directory", path)
	}
	// Windows reports synthetic 0o666/0o777 bits regardless of the ACL, so enforcing POSIX
	// permissions there would fail every publication without adding any protection.
	if runtime.GOOS == "windows" {
		return nil
	}
	if fi.Mode().Perm()&0o077 == 0 {
		return nil
	}
	if err := os.Chmod(path, runtimeDirPerm); err != nil {
		return fmt.Errorf("protocolpacket: %s is group/world accessible (%#o) and cannot be restricted: %w", path, fi.Mode().Perm(), err)
	}
	again, err := os.Lstat(path)
	if err != nil {
		return fmt.Errorf("protocolpacket: cannot stat %s: %w", path, err)
	}
	if again.Mode().Perm()&0o077 != 0 {
		return fmt.Errorf("protocolpacket: %s stays group/world accessible (%#o); refusing to publish protocol bodies there", path, again.Mode().Perm())
	}
	return nil
}

// linkFile is os.Link. It is a seam so the fail-closed behaviour of a filesystem that cannot
// hard link is exercised deterministically in tests; production always uses os.Link.
var linkFile = os.Link

// publish puts body at path without ever truncating, clobbering, or following a symlink.
//
// The body is staged in a sibling temporary file and hard-linked into place: link() either
// publishes the complete file in one step or fails with EEXIST, so no reader can ever open the
// attested path and find an empty or partial protocol body. There is deliberately NO direct-write
// fallback. An O_WRONLY|O_CREATE|O_EXCL write at the target does not clobber, but it does publish
// an empty file and fill it afterwards, which reintroduces exactly the partial-body window this
// renderer exists to close — a launch reading between create and write gets protocol text that is
// not the attested text. Where link() is unsupported, publication fails closed and says so;
// claiming support by weakening publication would be worse than an honest error.
func publish(dir, path, body string) error {
	published, err := existingBody(path)
	if err != nil {
		return err
	}
	if published {
		return verifyPublished(path, body)
	}

	tmp, err := os.CreateTemp(dir, ".staged-packet-*")
	if err != nil {
		return fmt.Errorf("protocolpacket: cannot stage a packet body in %s: %w", dir, err)
	}
	tmpName := tmp.Name()
	defer os.Remove(tmpName)
	if err := writeAndClose(tmp, body); err != nil {
		return err
	}

	switch lerr := linkFile(tmpName, path); {
	case lerr == nil:
		return nil
	case errors.Is(lerr, os.ErrExist):
		// Another publisher won the race for this content-addressed name — or something else
		// was planted there. This path bypassed existingBody's checks, so verifyPublished
		// re-validates the entry itself rather than assuming a well-behaved racer.
		return verifyPublished(path, body)
	default:
		return fmt.Errorf("protocolpacket: cannot publish %s atomically: %w; %s does not support the hard link that makes publication indivisible, and a direct write would expose a partial body at the attested path",
			path, lerr, dir)
	}
}

// existingBody reports whether path is already a regular file. A symlink or any other
// non-regular entry is an error rather than a target, so nothing is ever written through it.
func existingBody(path string) (bool, error) {
	fi, err := os.Lstat(path)
	switch {
	case errors.Is(err, os.ErrNotExist):
		return false, nil
	case err != nil:
		return false, fmt.Errorf("protocolpacket: cannot stat %s: %w", path, err)
	case fi.Mode()&os.ModeSymlink != 0:
		return false, fmt.Errorf("protocolpacket: %s is a symlink; a published packet body must be a regular file", path)
	case !fi.Mode().IsRegular():
		return false, fmt.Errorf("protocolpacket: %s exists and is not a regular file", path)
	}
	return true, nil
}

// verifyPublished accepts an already-published path only when its bytes are the attested body,
// which makes republication idempotent and a mismatch an error instead of an overwrite.
//
// It must not follow a symlink: os.ReadFile does, so a swapped link could have made an attacker's
// file "prove" the publication and let a launch read arbitrary text under the attested name. Every
// caller — including the link-EEXIST race, which never went through existingBody — validates the
// entry here, on the descriptor it actually read.
func verifyPublished(path, body string) error {
	f, err := openPublished(path)
	if err != nil {
		return err
	}
	defer f.Close()
	have, err := io.ReadAll(f)
	if err != nil {
		return fmt.Errorf("protocolpacket: %s exists but cannot be read: %w", path, err)
	}
	if string(have) == body {
		return nil
	}
	return fmt.Errorf("protocolpacket: %s already holds a different body (sha256 %s, attested %s); refusing to overwrite a published packet",
		path, protocolcore.ShortHash(Hash(string(have))), protocolcore.ShortHash(Hash(body)))
}

// openPublished opens an existing published body for validation without trusting the path. The
// entry is Lstat-ed (a symlink or non-regular entry is refused as itself, never as its target),
// then opened, then the OPEN DESCRIPTOR is confirmed to be that very same regular file. If a
// symlink is swapped in after the Lstat, open() follows it and the descriptor identifies a
// different file, which is refused. Bytes are therefore only ever read from a file that was, at
// one moment, the regular file at this path.
func openPublished(path string) (*os.File, error) {
	fi, err := os.Lstat(path)
	if err != nil {
		return nil, fmt.Errorf("protocolpacket: cannot stat %s: %w", path, err)
	}
	if fi.Mode()&os.ModeSymlink != 0 {
		return nil, fmt.Errorf("protocolpacket: %s is a symlink; a published packet body must be a regular file", path)
	}
	if !fi.Mode().IsRegular() {
		return nil, fmt.Errorf("protocolpacket: %s exists and is not a regular file", path)
	}
	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("protocolpacket: %s exists but cannot be read: %w", path, err)
	}
	opened, err := f.Stat()
	if err != nil {
		f.Close()
		return nil, fmt.Errorf("protocolpacket: cannot stat the opened %s: %w", path, err)
	}
	if !opened.Mode().IsRegular() || !os.SameFile(fi, opened) {
		f.Close()
		return nil, fmt.Errorf("protocolpacket: %s was replaced while it was being validated; refusing to trust a published packet that changed under the check", path)
	}
	return f, nil
}

func writeAndClose(f *os.File, body string) error {
	name := f.Name()
	if _, err := f.WriteString(body); err != nil {
		f.Close()
		return fmt.Errorf("protocolpacket: cannot write %s: %w", name, err)
	}
	// CreateTemp's 0600 is masked by umask; set the mode on the descriptor, before the bytes are
	// reachable under the published name, so the path is never briefly readable by anyone else.
	if runtime.GOOS != "windows" {
		if err := f.Chmod(bodyFilePerm); err != nil {
			f.Close()
			return fmt.Errorf("protocolpacket: cannot restrict %s: %w", name, err)
		}
	}
	// Darwin's File.Sync issues F_FULLFSYNC, which some shared volumes reject with ENOTTY even
	// though ordinary fsync works there. fsutil.SyncFile falls back for exactly that errno and
	// still requires the fsync to succeed: an unflushed body is a publication failure.
	if err := fsutil.SyncFile(f); err != nil {
		f.Close()
		return fmt.Errorf("protocolpacket: cannot flush %s: %w", name, err)
	}
	if err := f.Close(); err != nil {
		return fmt.Errorf("protocolpacket: cannot close %s: %w", name, err)
	}
	return nil
}

func safeName(s string) string {
	var b strings.Builder
	for _, r := range s {
		switch {
		case r >= 'a' && r <= 'z', r >= '0' && r <= '9', r == '-':
			b.WriteRune(r)
		default:
			b.WriteRune('_')
		}
	}
	if b.Len() == 0 {
		return "unknown"
	}
	return b.String()
}

// Render is the one-call entry: resolve authority, load the map, build, and write the body.
// The returned error is an authority error (nothing emitted) or an I/O error; a fallback or a
// refusal is reported in the context, not as an error, so callers can attest it.
func Render(root string, store protocolcore.Store, req Request) (Context, error) {
	src, err := ResolveSource(root, store)
	if err != nil {
		return Context{}, err
	}
	m, err := LoadMap(MapPath(root))
	if err != nil && !errors.Is(err, ErrNoMap) {
		// A malformed map is not "no map": it is reported as the fallback reason so the
		// operator sees it, while the body stays full context.
		m = nil
		req.Flags = append([]string(nil), req.Flags...)
		c := Build(src, nil, req)
		if c.ContextMode != ModeRefused {
			c.Problems = append([]string{"malformed-map:" + err.Error()}, c.Problems...)
			if c.ContextMode == ModeFullFallback {
				c.FallbackReason = "malformed-map:" + err.Error() + "; " + c.FallbackReason
			} else if c.Shadow != nil {
				c.Shadow.WouldFallbackReason = "malformed-map:" + err.Error() + "; " + c.Shadow.WouldFallbackReason
			}
			if werr := WriteBody(root, &c); werr != nil {
				return c, werr
			}
		}
		return c, nil
	}
	c := Build(src, m, req)
	if c.ContextMode == ModeRefused {
		return c, nil
	}
	if err := WriteBody(root, &c); err != nil {
		return c, err
	}
	return c, nil
}
