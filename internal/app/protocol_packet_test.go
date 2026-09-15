package app

import (
	"bytes"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"parley-deck-cli/internal/protocolpacket"
)

const packetTestProtocol = "# Fixture protocol\n\n**Transport:** `local-dir`\n\n## 1. Always\n\nalways text\n\n### Phase 1 — one\n\nphase one\n\n### Phase 2 — two\n\nphase two\n"

const packetTestMap = `schema: parley.packet-applicability/v1
source: fixture
blocks:
  - locator: "# Fixture protocol"
    include: always
  - locator: "## 1. Always"
    include: always
  - locator: "### Phase 1 — one"
    include: when
    phases: [1]
    trigger: "phase one work"
  - locator: "### Phase 2 — two"
    include: when
    phases: [2]
    trigger: "phase two work"
`

// packetFixture builds a source-role deck in a temp root. PARLEY_HOME points at an empty home
// so no real core store is touched.
func packetFixture(t *testing.T, withMap bool) string {
	t.Helper()
	t.Setenv("PARLEY_HOME", t.TempDir())
	root := t.TempDir()
	meta := filepath.Join(root, "parley-deck", "meta")
	if err := os.MkdirAll(meta, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(meta, "version.json"), []byte(`{"protocolRole":"source"}`), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(root, "parley-deck", "COOPERATION.md"), []byte(packetTestProtocol), 0o644); err != nil {
		t.Fatal(err)
	}
	if withMap {
		if err := os.WriteFile(filepath.Join(meta, "packet-applicability.yaml"), []byte(packetTestMap), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	return root
}

func runPacket(t *testing.T, args ...string) (int, string, string) {
	t.Helper()
	var out, errb bytes.Buffer
	code := runProtocolPacket(args, &out, &errb)
	return code, out.String(), errb.String()
}

func TestProtocolPacketDefaultsToFullContextAndWritesUnderRuntimeDir(t *testing.T) {
	root := packetFixture(t, true)
	code, out, errOut := runPacket(t, "--dir", root, "--phase", "1", "--track", "fast", "--json")
	if code != 0 {
		t.Fatalf("exit %d: %s", code, errOut)
	}
	var c protocolpacket.Context
	if err := json.Unmarshal([]byte(out), &c); err != nil {
		t.Fatalf("json: %v\n%s", err, out)
	}
	if c.ContextMode != protocolpacket.ModeFull || c.SourceSHA256 != protocolpacket.Hash(packetTestProtocol) || c.PacketSHA256 != c.SourceSHA256 {
		t.Fatalf("attestation: %+v", c.Attestation)
	}
	if c.Shadow == nil || c.Shadow.OmittedBlocks != 1 || c.Shadow.IncludedBlocks != 3 {
		t.Fatalf("shadow: %+v", c.Shadow)
	}
	want := filepath.Join(root, ".parley-runtime", "protocol-packets")
	if !strings.HasPrefix(c.BodyPath, want) {
		t.Fatalf("body at %s, want under %s", c.BodyPath, want)
	}
	if b, err := os.ReadFile(c.BodyPath); err != nil || string(b) != packetTestProtocol {
		t.Fatalf("full body must be the source verbatim: %v", err)
	}
	if c.Source.Role != "source" || strings.Contains(c.Source.Authority, "verified") {
		t.Fatalf("a source-role deck must not be attributed to a core release: %s", c.Source.Authority)
	}
}

func TestProtocolPacketOptimizeIsExplicitAndPrintsAVerbatimPacket(t *testing.T) {
	root := packetFixture(t, true)
	code, out, errOut := runPacket(t, "--dir", root, "--phase", "2", "--track", "fast", "--optimize", "--print")
	if code != 0 {
		t.Fatalf("exit %d: %s", code, errOut)
	}
	if !strings.Contains(out, "context_mode    : packet") || !strings.Contains(out, "phase two") || strings.Contains(out, "phase one\n") {
		t.Fatalf("unexpected output:\n%s", out)
	}
	if !strings.Contains(out, "| `### Phase 1 — one` |") {
		t.Fatalf("omission index missing:\n%s", out)
	}
}

func TestProtocolPacketFallsBackWithoutAMapAndCheckReportsIt(t *testing.T) {
	root := packetFixture(t, false)
	code, out, _ := runPacket(t, "--dir", root, "--phase", "1", "--track", "fast", "--optimize", "--json")
	if code != 0 {
		t.Fatalf("a fallback still delivers full context; exit %d", code)
	}
	var c protocolpacket.Context
	if err := json.Unmarshal([]byte(out), &c); err != nil {
		t.Fatal(err)
	}
	if c.ContextMode != protocolpacket.ModeFullFallback || c.FallbackReason != "no-applicability-map" {
		t.Fatalf("%+v", c.Attestation)
	}
	code, out, _ = runPacket(t, "check", "--dir", root)
	if code != 1 || !strings.Contains(out, "ABSENT") || !strings.Contains(out, "FAILED") {
		t.Fatalf("check exit %d:\n%s", code, out)
	}
}

func TestProtocolPacketCheckPassesOnTheFixtureAndFailsOnUnknownApplicability(t *testing.T) {
	root := packetFixture(t, true)
	code, out, errOut := runPacket(t, "check", "--dir", root, "--json")
	if code != 0 {
		t.Fatalf("exit %d: %s %s", code, out, errOut)
	}
	var rep protocolpacket.Report
	if err := json.Unmarshal([]byte(out), &rep); err != nil || !rep.OK || rep.Blocks != 4 {
		t.Fatalf("%v %+v", err, rep)
	}
	// Add a heading the map does not know: unknown applicability fails the check.
	p := filepath.Join(root, "parley-deck", "COOPERATION.md")
	if err := os.WriteFile(p, []byte(packetTestProtocol+"\n## 99. New rule\n\nnew\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	code, out, _ = runPacket(t, "check", "--dir", root)
	if code != 1 || !strings.Contains(out, "## 99. New rule") {
		t.Fatalf("check exit %d:\n%s", code, out)
	}
	code, out, _ = runPacket(t, "--dir", root, "--phase", "1", "--track", "fast", "--optimize")
	if code != 0 || !strings.Contains(out, "context_mode    : full-fallback") || !strings.Contains(out, "unclassified:## 99. New rule") {
		t.Fatalf("exit %d:\n%s", code, out)
	}
}

func TestProtocolPacketRefusesOnMissingAuthorityAndSecrets(t *testing.T) {
	root := packetFixture(t, true)
	if err := os.Remove(filepath.Join(root, "parley-deck", "meta", "version.json")); err != nil {
		t.Fatal(err)
	}
	code, _, errOut := runPacket(t, "--dir", root, "--phase", "1", "--track", "fast")
	if code != 1 || !strings.Contains(errOut, "authority") {
		t.Fatalf("exit %d: %s", code, errOut)
	}
	if _, err := os.Stat(filepath.Join(root, ".parley-runtime")); !os.IsNotExist(err) {
		t.Fatal("nothing may be written when authority is missing")
	}

	root = packetFixture(t, true)
	p := filepath.Join(root, "parley-deck", "COOPERATION.md")
	if err := os.WriteFile(p, []byte(packetTestProtocol+"\nAKIAABCDEFGHIJKLMNOP\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	code, out, errOut := runPacket(t, "--dir", root, "--phase", "1", "--track", "fast", "--json")
	if code != 1 || !strings.Contains(errOut, "REFUSED") || !strings.Contains(errOut, "secret-detected:aws-access-key") {
		t.Fatalf("exit %d: %s", code, errOut)
	}
	if strings.Contains(out, "AKIAABCDEFGHIJKLMNOP") {
		t.Fatal("a refused context leaked the source")
	}
	if _, err := os.Stat(filepath.Join(root, ".parley-runtime")); !os.IsNotExist(err) {
		t.Fatal("a refused context must not write a body")
	}
}

// Publication is content-addressed and immutable: the same request resolves to the same path,
// and a body changed on disk is refused instead of silently replaced. Handing a launch
// protocol text that is not the attested text is the failure this command exists to prevent.
func TestProtocolPacketRepublicationIsIdempotentAndRefusesATamperedBody(t *testing.T) {
	root := packetFixture(t, true)
	args := []string{"--dir", root, "--phase", "1", "--track", "fast", "--json"}

	code, out, errOut := runPacket(t, args...)
	if code != 0 {
		t.Fatalf("exit %d: %s", code, errOut)
	}
	var first protocolpacket.Context
	if err := json.Unmarshal([]byte(out), &first); err != nil {
		t.Fatal(err)
	}
	if !strings.HasSuffix(filepath.Base(first.BodyPath), first.PacketSHA256+".md") {
		t.Fatalf("%s does not carry the full body digest %s", filepath.Base(first.BodyPath), first.PacketSHA256)
	}

	code, out, errOut = runPacket(t, args...)
	if code != 0 {
		t.Fatalf("exit %d: %s", code, errOut)
	}
	var second protocolpacket.Context
	if err := json.Unmarshal([]byte(out), &second); err != nil {
		t.Fatal(err)
	}
	if second.BodyPath != first.BodyPath {
		t.Fatalf("republication moved the body: %s vs %s", second.BodyPath, first.BodyPath)
	}
	if entries, err := os.ReadDir(filepath.Join(root, ".parley-runtime", "protocol-packets")); err != nil || len(entries) != 1 {
		t.Fatalf("runtime dir holds %d entries: %v", len(entries), err)
	}

	const tampered = "not the protocol\n"
	if err := os.WriteFile(first.BodyPath, []byte(tampered), 0o600); err != nil {
		t.Fatal(err)
	}
	code, _, errOut = runPacket(t, args...)
	if code != 1 || !strings.Contains(errOut, "refusing to overwrite") {
		t.Fatalf("exit %d: %s", code, errOut)
	}
	if b, _ := os.ReadFile(first.BodyPath); string(b) != tampered {
		t.Fatal("the tampered body was overwritten instead of failing closed")
	}
}

func TestProtocolPacketRequiresPhaseAndPrintsHelp(t *testing.T) {
	code, _, errOut := runPacket(t, "--dir", t.TempDir())
	if code != 2 || !strings.Contains(errOut, "--phase is required") {
		t.Fatalf("exit %d: %s", code, errOut)
	}
	code, out, _ := runPacket(t, "--help")
	if code != 0 || !strings.Contains(out, "packet check") {
		t.Fatalf("exit %d: %s", code, out)
	}
}

// The repository's own deck is a source-role deck with a committed map: the check must pass
// against the live protocol, and it writes nothing.
func TestProtocolPacketCheckPassesOnThisRepository(t *testing.T) {
	t.Setenv("PARLEY_HOME", t.TempDir())
	code, out, errOut := runPacket(t, "check", "--dir", "../..")
	if code != 0 {
		t.Fatalf("exit %d:\n%s%s", code, out, errOut)
	}
}
