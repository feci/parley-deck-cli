package app

import (
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"golang.org/x/sys/unix"

	"parley-deck-cli/internal/config"
	"parley-deck-cli/internal/fsutil"
	"parley-deck-cli/internal/protocol"
	"parley-deck-cli/internal/protocolcore"
)

const protocolUsage = `usage:
  parley protocol status  [--dir DIR] [--json]
  parley protocol render  [--dir DIR] [--dry-run] [--yes]
  parley protocol check   [--dir DIR] [--json]
  parley protocol publish --version V --from FILE            (attended; requires a TTY)

  The deck's COOPERATION.md is a GENERATED VIEW of the global core in
  ~/.parley/protocol/core/<version>/. Do not hand-edit it — ` + "`protocol check`" + ` reports edits
  rather than silently overwriting them.`

func runProtocol(args []string, stdout, stderr io.Writer) int {
	if len(args) == 0 {
		fmt.Fprintln(stderr, protocolUsage)
		return 2
	}
	sub, rest := args[0], args[1:]
	fs := flag.NewFlagSet("protocol", flag.ContinueOnError)
	fs.SetOutput(stderr)
	dir := fs.String("dir", ".", "workspace directory")
	jsonOut := fs.Bool("json", false, "machine-readable output")
	dryRun := fs.Bool("dry-run", false, "print what would change; write nothing")
	yes := fs.Bool("yes", false, "write without confirmation")
	version := fs.String("version", "", "protocol publish: the new core version")
	from := fs.String("from", "", "protocol publish: file holding the new core body")
	if err := fs.Parse(rest); err != nil {
		return 2
	}
	root, err := filepath.Abs(*dir)
	if err != nil {
		fmt.Fprintf(stderr, "protocol: %v\n", err)
		return 1
	}
	switch sub {
	case "status":
		return protocolStatus(root, *jsonOut, stdout, stderr)
	case "render":
		return protocolRender(root, *dryRun, *yes, stdout, stderr)
	case "check":
		return protocolCheck(root, *jsonOut, stdout, stderr)
	case "publish":
		return protocolPublish(*version, *from, stdout, stderr)
	default:
		fmt.Fprintf(stderr, "protocol: unknown subcommand %q (want status|render|check|publish)\n", sub)
		return 2
	}
}

// coreStore resolves the global release store from the Parley home.
func coreStore() (protocolcore.Store, error) {
	home := config.CentralHome()
	if strings.TrimSpace(home) == "" {
		return protocolcore.Store{}, errors.New("cannot resolve the Parley home directory (set PARLEY_HOME)")
	}
	return protocolcore.StoreAt(home), nil
}

// pinnedVersion reads the deck's pinned core version. Absent means the deck has not adopted a
// core yet, which is a state to report, never one to guess around: substituting "the newest
// installed" is exactly what D8 forbids.
func pinnedVersion(root string) (string, error) {
	b, err := os.ReadFile(deckLockPath(root))
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return "", nil
		}
		return "", err
	}
	for _, l := range strings.Split(string(b), "\n") {
		if v, ok := strings.CutPrefix(strings.TrimSpace(l), "core-version:"); ok {
			return strings.TrimSpace(v), nil
		}
	}
	return "", nil
}

func deckLockPath(root string) string {
	return filepath.Join(root, protocol.DeckDir, "meta", "protocol-lock.yaml")
}

func deckProtocolPath(root string) string {
	return filepath.Join(root, protocol.DeckDir, "COOPERATION.md")
}

// resolveRelease loads the release the deck is pinned to.
func resolveRelease(root string) (protocolcore.Release, string, error) {
	store, err := coreStore()
	if err != nil {
		return protocolcore.Release{}, "", err
	}
	v, err := pinnedVersion(root)
	if err != nil {
		return protocolcore.Release{}, "", err
	}
	if v == "" {
		versions, _ := store.Versions()
		if len(versions) == 0 {
			return protocolcore.Release{}, "", fmt.Errorf(
				"no core release is installed in %s and this deck pins none — run `parley protocol publish` first", store.Root)
		}
		return protocolcore.Release{}, "", fmt.Errorf(
			"this deck pins no core version. Installed: %s.\n"+
				"Adopt one explicitly by writing `core-version: <v>` into %s — a version is never guessed",
			strings.Join(versions, ", "), deckLockPath(root))
	}
	rel, err := store.Load(v)
	return rel, v, err
}

func protocolStatus(root string, jsonOut bool, stdout, stderr io.Writer) int {
	store, err := coreStore()
	if err != nil {
		fmt.Fprintf(stderr, "protocol status: %v\n", err)
		return 1
	}
	versions, _ := store.Versions()
	pinned, _ := pinnedVersion(root)
	deck, deckErr := os.ReadFile(deckProtocolPath(root))
	payload := map[string]any{
		"store":              store.Root,
		"installed":          versions,
		"pinned":             pinned,
		"deck_protocol":      deckProtocolPath(root),
		"deck_protocol_read": deckErr == nil,
	}
	if deckErr == nil {
		payload["deck_sha256"] = protocolcore.Hash(string(deck))
	}
	if jsonOut {
		return writeJSON(stdout, stderr, payload)
	}
	fmt.Fprintf(stdout, "core store : %s\n", store.Root)
	if len(versions) == 0 {
		fmt.Fprintln(stdout, "installed  : (none)")
	} else {
		fmt.Fprintf(stdout, "installed  : %s\n", strings.Join(versions, ", "))
	}
	fmt.Fprintf(stdout, "deck pins  : %s\n", orDash(pinned))
	if deckErr == nil {
		fmt.Fprintf(stdout, "deck view  : %s (%s)\n", deckProtocolPath(root), protocolcore.ShortHash(protocolcore.Hash(string(deck))))
	} else {
		fmt.Fprintf(stdout, "deck view  : %v\n", deckErr)
	}
	return 0
}

func protocolRender(root string, dryRun, yes bool, stdout, stderr io.Writer) int {
	rel, _, err := resolveRelease(root)
	if err != nil {
		fmt.Fprintf(stderr, "protocol render: %v\n", err)
		return 1
	}
	path := deckProtocolPath(root)
	prior, _ := os.ReadFile(path)
	res, err := protocolcore.Render(rel, string(prior))
	if err != nil {
		fmt.Fprintf(stderr, "protocol render: %v\n", err)
		return 1
	}

	// G1: report every block replaced or removed, in preview AND on apply. Regeneration that
	// drops content silently is how a deck's genuine local section was destroyed once already.
	report := func(w io.Writer) {
		if len(res.Preserved) > 0 {
			fmt.Fprintf(w, "preserved from this deck: %s\n", strings.Join(res.Preserved, ", "))
		}
		if len(res.Removed) > 0 {
			fmt.Fprintf(w, "the following section(s) exist in this deck but NOT in core %s and will be REMOVED:\n", rel.Version)
			for _, h := range res.Removed {
				fmt.Fprintf(w, "  - %s\n", h)
			}
			fmt.Fprintln(w, "  (they are project-local content; the overlay mechanism that will carry them is not shipped yet)")
		}
	}

	if res.Body == string(prior) {
		fmt.Fprintf(stdout, "protocol render: %s already matches core %s — nothing to do\n", path, rel.Version)
		return 0
	}
	if dryRun || !yes {
		report(stdout)
		fmt.Fprintf(stdout, "would regenerate %s from core %s (%s).\nNothing was written. Re-run with --yes to apply.\n",
			path, rel.Version, protocolcore.ShortHash(rel.SHA256))
		return 0
	}
	report(stdout)
	perm := os.FileMode(0o644)
	if fi, statErr := os.Stat(path); statErr == nil {
		perm = fi.Mode().Perm()
	}
	if err := fsutil.WriteFileAtomic(path, []byte(res.Body), perm); err != nil {
		fmt.Fprintf(stderr, "protocol render: %v\n", err)
		return 1
	}
	fmt.Fprintf(stdout, "Wrote %s from core %s (%s)\n", path, rel.Version, protocolcore.ShortHash(rel.SHA256))
	return 0
}

func protocolCheck(root string, jsonOut bool, stdout, stderr io.Writer) int {
	rel, version, err := resolveRelease(root)
	if err != nil {
		fmt.Fprintf(stderr, "protocol check: %v\n", err)
		return 1
	}
	path := deckProtocolPath(root)
	prior, readErr := os.ReadFile(path)
	if readErr != nil {
		fmt.Fprintf(stderr, "protocol check: %v\n", readErr)
		return 1
	}
	res, err := protocolcore.Render(rel, string(prior))
	if err != nil {
		fmt.Fprintf(stderr, "protocol check: %v\n", err)
		return 1
	}
	inSync := res.Body == string(prior)
	status := "in-sync"
	if !inSync {
		status = "hand-edited-or-stale"
	}
	if jsonOut {
		code := 0
		if !inSync {
			code = 1
		}
		_ = writeJSON(stdout, stderr, map[string]any{
			"status":       status,
			"core_version": version,
			"core_sha256":  rel.SHA256,
			"deck_sha256":  protocolcore.Hash(string(prior)),
			"removed":      res.Removed,
		})
		return code
	}
	if inSync {
		fmt.Fprintf(stdout, "protocol check: in sync with core %s (%s)\n", version, protocolcore.ShortHash(rel.SHA256))
		return 0
	}
	// Report, never overwrite — the same posture `roster show` takes for `legacy-roster`.
	fmt.Fprintf(stdout, "protocol check: %s does NOT match core %s.\n", path, version)
	if len(res.Removed) > 0 {
		fmt.Fprintln(stdout, "sections present here but not in the core:")
		for _, h := range res.Removed {
			fmt.Fprintf(stdout, "  - %s\n", h)
		}
	}
	fmt.Fprintln(stdout, "Run `parley protocol render --dry-run` to see the regeneration, then --yes to apply.")
	return 1
}

// protocolPublish is the ATTENDED publisher — gate G2's sole audited exception to "no write path
// into the core". It requires a controlling terminal, so an agent run (which has no TTY) cannot
// reach it, and it creates a NEW release rather than modifying one.
func protocolPublish(version, from string, stdout, stderr io.Writer) int {
	if strings.TrimSpace(version) == "" || strings.TrimSpace(from) == "" {
		fmt.Fprintln(stderr, "protocol publish: --version and --from are both required")
		return 2
	}
	if !hasTTY() {
		fmt.Fprintln(stderr, "protocol publish: refusing — publishing a core release is an attended operation and no controlling terminal is present.\n"+
			"Only the user may change the global protocol; an agent run cannot.")
		return 2
	}
	body, err := os.ReadFile(from)
	if err != nil {
		fmt.Fprintf(stderr, "protocol publish: %v\n", err)
		return 1
	}
	store, err := coreStore()
	if err != nil {
		fmt.Fprintf(stderr, "protocol publish: %v\n", err)
		return 1
	}
	rel, err := store.Publish(version, string(body))
	if err != nil {
		fmt.Fprintf(stderr, "protocol publish: %v\n", err)
		return 1
	}
	fmt.Fprintf(stdout, "Published core %s (%s) to %s\n", rel.Version, protocolcore.ShortHash(rel.SHA256), store.ReleaseDir(version))
	return 0
}

// hasTTY reports whether this process has a CONTROLLING TERMINAL.
//
// The obvious check — stdin's mode has os.ModeCharDevice — is wrong for this gate, and testing the
// real entry point caught it: `/dev/null` is a character device too, so `protocol publish … <
// /dev/null` sailed straight through the refusal it was supposed to hit. An agent run redirects
// stdin exactly that way, which is precisely the case G2 must stop.
//
// Asking the kernel for the terminal attributes of the fd is the actual question: it succeeds only
// for a real tty.
func hasTTY() bool {
	for _, f := range []*os.File{os.Stdin, os.Stdout} {
		if _, err := unix.IoctlGetTermios(int(f.Fd()), termiosGet); err == nil {
			return true
		}
	}
	return false
}
