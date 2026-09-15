package protocolpacket

import (
	"errors"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
	"testing"

	"parley-deck-cli/internal/protocolcore"
)

// The live protocol and its map, relative to this package directory (go test runs there).
const (
	liveProtocolPath = "../../parley-deck/COOPERATION.md"
	liveMapPath      = "../../parley-deck/meta/packet-applicability.yaml"
)

func liveSource(t *testing.T) Source {
	t.Helper()
	b, err := os.ReadFile(liveProtocolPath)
	if err != nil {
		t.Fatalf("live protocol: %v", err)
	}
	return Source{Path: "parley-deck/COOPERATION.md", Role: "source", Authority: "test", Raw: string(b), Transport: headerTransport(string(b))}
}

func liveMap(t *testing.T) *Map {
	t.Helper()
	m, err := LoadMap(liveMapPath)
	if err != nil {
		t.Fatalf("live map: %v", err)
	}
	return m
}

const fixtureSource = "# Title\n\n**Transport:** `github-pr`\n\n## A. Always\n\ntext a\n\n```\n## not a heading\n```\n\n### Phase 1 — one\n\nphase one text\n\n### Phase 2 — two\n\nphase two text\n\n## Dup\n\nfirst\n\n## Dup\n\nsecond\n"

const fixtureMap = `schema: parley.packet-applicability/v1
source: fixture
blocks:
  - locator: "# Title"
    include: always
  - locator: "## A. Always"
    include: always
  - locator: "### Phase 1 — one"
    include: when
    phases: [1]
    trigger: "phase one work"
  - locator: "### Phase 2 — two"
    include: when
    phases: [2]
    requires: ["## Dup #2"]
    trigger: "phase two work"
  - locator: "## Dup"
    include: on-demand
    trigger: "first dup"
  - locator: "## Dup #2"
    include: on-demand
    trigger: "second dup"
`

func fixture(t *testing.T) (Source, *Map) {
	t.Helper()
	m, err := ParseMap(fixtureMap)
	if err != nil {
		t.Fatal(err)
	}
	return Source{Path: "fixture.md", Role: "source", Raw: fixtureSource, Transport: "github-pr"}, m
}

func TestParseBlocksAreContiguousAndFenceHeadingsAreNotBlocks(t *testing.T) {
	blocks := Parse(fixtureSource)
	var locs []string
	next := 1
	for _, b := range blocks {
		if b.StartLine != next {
			t.Fatalf("block %q starts at %d, want %d (every line belongs to exactly one block)", b.Locator, b.StartLine, next)
		}
		next = b.EndLine + 1
		locs = append(locs, b.Locator)
	}
	if next != len(strings.Split(fixtureSource, "\n"))+1 {
		t.Fatalf("blocks end at %d, want the last line", next-1)
	}
	want := []string{"# Title", "## A. Always", "### Phase 1 — one", "### Phase 2 — two", "## Dup", "## Dup #2"}
	if strings.Join(locs, "|") != strings.Join(want, "|") {
		t.Fatalf("locators %q, want %q", locs, want)
	}
	if strings.Join(Locators(fixtureSource), "|") != strings.Join(want, "|") {
		t.Fatalf("Locators disagrees with Parse")
	}
}

func TestLiveMapIsCompleteHonorsNeverCutAndRendersEveryPhase(t *testing.T) {
	src := liveSource(t)
	if src.Transport != "github-pr" {
		t.Fatalf("live transport %q", src.Transport)
	}
	rep := Check(src, liveMap(t))
	if !rep.OK {
		t.Fatalf("live map check failed: unclassified=%v stale=%v never-cut=%v fallbacks=%v secret=%q",
			rep.Unclassified, rep.Stale, rep.NeverCut, rep.Fallbacks, rep.Secret)
	}
	if rep.Blocks != 69 {
		t.Fatalf("live protocol parsed into %d blocks; the map was written for 69 — reconcile both", rep.Blocks)
	}
}

func TestOptimizedPacketIsVerbatimAndIndexIsCompleteAndDisjoint(t *testing.T) {
	src, m := liveSource(t), liveMap(t)
	c := Build(src, m, Request{Phase: 6, Track: "deliberation", Optimize: true})
	if c.ContextMode != ModePacket {
		t.Fatalf("mode %s (%s)", c.ContextMode, c.FallbackReason)
	}
	if c.Request.Transport != "github-pr" {
		t.Fatalf("transport not taken from the header: %q", c.Request.Transport)
	}
	if c.SourceSHA256 != Hash(src.Raw) || c.PacketSHA256 != Hash(c.Body) || c.PacketSHA256 == c.SourceSHA256 {
		t.Fatal("hash attestation is wrong")
	}
	blocks := Parse(src.Raw)
	if len(c.Index) != len(blocks) {
		t.Fatalf("index has %d rows for %d blocks", len(c.Index), len(blocks))
	}
	included, omitted := 0, 0
	for _, r := range c.Index {
		if r.Included {
			included++
			if !strings.Contains(c.Body, strings.TrimRight(r.Text, "\n")) {
				t.Fatalf("included block %q is not reproduced verbatim", r.Locator)
			}
		} else {
			omitted++
			if !strings.Contains(c.Body, "| `"+r.Locator+"` |") {
				t.Fatalf("omitted block %q is missing from the omission index", r.Locator)
			}
			if r.Trigger == "" {
				t.Fatalf("omitted block %q has no trigger", r.Locator)
			}
		}
	}
	if included+omitted != len(blocks) || included == 0 || omitted == 0 {
		t.Fatalf("included %d + omitted %d != %d", included, omitted, len(blocks))
	}
	state := map[string]bool{}
	for _, r := range c.Index {
		state[r.Locator] = r.Included
	}
	for loc, want := range map[string]bool{
		"## 14. Automated outer loop (loop engineering) — the human brake": true,
		"## 6. Conflict-avoidance mechanics":                               true,
		"### Phase 6 — Code review rounds":                                 true,
		"#### Review briefs and dispositions":                              true,
		"### 15.1 Scope, ownership, location":                              true,
		"### 11.B — GitHub Pull Requests":                                  true,
		"### 11.A — Local directory":                                       false,
		"### Phase 1 — Round 1 (independent analysis)":                     false,
		"## 7. Changing this protocol":                                     false,
		"## 12. Pipeline blocks & action stages":                           false,
	} {
		if got, ok := state[loc]; !ok || got != want {
			t.Fatalf("%q included=%v (present=%v), want %v", loc, got, ok, want)
		}
	}
	// §7 in every phase of a protocol-change idea, derived from the slug.
	c2 := Build(src, m, Request{Phase: 6, Track: "deliberation", IdeaSlug: "meta-protocol-change-x", Optimize: true})
	for _, r := range c2.Index {
		if r.Locator == "## 7. Changing this protocol" && !r.Included {
			t.Fatal("§7 must be included for a meta-protocol-change idea")
		}
	}
}

func TestDefaultIsFullContextWithShadowPacket(t *testing.T) {
	src, m := liveSource(t), liveMap(t)
	c := Build(src, m, Request{Phase: 1, Track: "standard"})
	if c.ContextMode != ModeFull || c.Body != src.Raw || c.PacketSHA256 != c.SourceSHA256 || c.FallbackReason != "" {
		t.Fatalf("default is not full context: mode=%s reason=%q", c.ContextMode, c.FallbackReason)
	}
	if c.Shadow == nil || c.Shadow.PacketSHA256 == "" || c.Shadow.OmittedBlocks == 0 || c.Shadow.PacketBytes >= c.Shadow.SourceBytes {
		t.Fatalf("shadow record missing or implausible: %+v", c.Shadow)
	}
	p := Build(src, m, Request{Phase: 1, Track: "standard", Optimize: true})
	if p.PacketSHA256 != c.Shadow.PacketSHA256 {
		t.Fatal("shadow packet hash must equal the optimized packet hash for the same request")
	}
}

func TestUnknownApplicabilityIncludesTheBlockAndFallsBack(t *testing.T) {
	src, m := fixture(t)
	m2 := &Map{Schema: m.Schema, index: map[string]int{}}
	for _, r := range m.Blocks {
		if r.Locator == "## A. Always" {
			continue
		}
		m2.index[r.Locator] = len(m2.Blocks)
		m2.Blocks = append(m2.Blocks, r)
	}
	c := Build(src, m2, Request{Phase: 1, Track: "fast", Optimize: true})
	if c.ContextMode != ModeFullFallback || !strings.Contains(c.FallbackReason, "unclassified:## A. Always") || c.Body != src.Raw {
		t.Fatalf("mode=%s reason=%q", c.ContextMode, c.FallbackReason)
	}
	for _, r := range c.Index {
		if r.Locator == "## A. Always" && (r.Classification != ClassUnknown || !r.Included) {
			t.Fatalf("unknown block must be recorded unknown and included: %+v", r)
		}
	}
	rep := Check(src, m2)
	if rep.OK || len(rep.Unclassified) != 1 {
		t.Fatalf("check must fail on unknown applicability: %+v", rep)
	}
	if c := Build(src, nil, Request{Phase: 1, Track: "fast", Optimize: true}); c.ContextMode != ModeFullFallback || c.FallbackReason != "no-applicability-map" {
		t.Fatalf("missing map: %s %q", c.ContextMode, c.FallbackReason)
	}
}

func TestStaleMapAndUnresolvedDependencyFallBack(t *testing.T) {
	src, m := fixture(t)
	stale, _ := ParseMap(fixtureMap + "  - locator: \"## Gone\"\n    include: always\n")
	c := Build(src, stale, Request{Phase: 1, Track: "fast", Optimize: true})
	if c.ContextMode != ModeFullFallback || !strings.Contains(c.FallbackReason, "stale-map:## Gone") {
		t.Fatalf("stale entry: %s %q", c.ContextMode, c.FallbackReason)
	}
	dep := strings.Replace(fixtureMap, `requires: ["## Dup #2"]`, `requires: ["## Missing"]`, 1)
	m3, err := ParseMap(dep)
	if err != nil {
		t.Fatal(err)
	}
	c = Build(src, m3, Request{Phase: 2, Track: "fast", Optimize: true})
	if c.ContextMode != ModeFullFallback || !strings.Contains(c.FallbackReason, "dependency-unresolved:## Missing") {
		t.Fatalf("unresolved dependency: %s %q", c.ContextMode, c.FallbackReason)
	}
	// A resolvable dependency is pulled in.
	c = Build(src, m, Request{Phase: 2, Track: "fast", Optimize: true})
	if c.ContextMode != ModePacket {
		t.Fatalf("%s %q", c.ContextMode, c.FallbackReason)
	}
	for _, r := range c.Index {
		if r.Locator == "## Dup #2" && (!r.Included || !strings.Contains(r.Why, "required by")) {
			t.Fatalf("dependency not pulled in: %+v", r)
		}
		if r.Locator == "## Dup" && r.Included {
			t.Fatal("on-demand block without a dependent must stay omitted")
		}
	}
}

func TestUnexpectedPhaseTrackTransportAndFlagFallBack(t *testing.T) {
	src, m := fixture(t)
	for _, tc := range []struct {
		req  Request
		want string
	}{
		{Request{Phase: 9, Track: "fast", Optimize: true}, "unknown-phase:9"},
		{Request{Phase: 1, Track: "turbo", Optimize: true}, "unknown-track:turbo"},
		{Request{Phase: 1, Track: "fast", Transport: "carrier-pigeon", Optimize: true}, "unknown-transport:carrier-pigeon"},
		{Request{Phase: 1, Track: "fast", Flags: []string{"yolo"}, Optimize: true}, "unknown-flag:yolo"},
	} {
		c := Build(src, m, tc.req)
		if c.ContextMode != ModeFullFallback || !strings.Contains(c.FallbackReason, tc.want) || c.Body != src.Raw {
			t.Fatalf("%+v: mode=%s reason=%q", tc.req, c.ContextMode, c.FallbackReason)
		}
	}
	drift := src
	drift.Drift = "hand-edited"
	if c := Build(drift, m, Request{Phase: 1, Track: "fast", Optimize: true}); c.ContextMode != ModeFullFallback || !strings.Contains(c.FallbackReason, "drift:hand-edited") {
		t.Fatalf("drift: %s %q", c.ContextMode, c.FallbackReason)
	}
}

func TestDetectedSecretRefusesWithoutBodyOrPacketHash(t *testing.T) {
	src, m := fixture(t)
	src.Raw += "\ntoken ghp_ABCDEFGHIJKLMNOPQRSTUVWXYZ012345\n"
	for _, opt := range []bool{false, true} {
		c := Build(src, m, Request{Phase: 1, Track: "fast", Optimize: opt})
		if c.ContextMode != ModeRefused || c.Body != "" || c.PacketSHA256 != "" || c.FallbackReason != "secret-detected:github-token" || c.Index != nil {
			t.Fatalf("optimize=%v: %+v", opt, c.Attestation)
		}
		if err := WriteBody(t.TempDir(), &c); err == nil {
			t.Fatal("a refused context must not be written")
		}
	}
	if name := DetectSecret(liveSource(t).Raw); name != "" {
		t.Fatalf("the live protocol trips the %s detector; the renderer would refuse every launch", name)
	}
}

func TestMalformedMapsAreRefused(t *testing.T) {
	for name, doc := range map[string]string{
		"duplicate":       fixtureMap + "  - locator: \"# Title\"\n    include: always\n",
		"no-trigger":      strings.Replace(fixtureMap, "    trigger: \"phase one work\"\n", "", 1),
		"unknown-include": strings.Replace(fixtureMap, "include: on-demand", "include: sometimes", 1),
		"unknown-key":     strings.Replace(fixtureMap, "include: always", "include: always\n    extra: 1", 1),
		"bad-schema":      strings.Replace(fixtureMap, "parley.packet-applicability/v1", "v0", 1),
		"bad-phase":       strings.Replace(fixtureMap, "phases: [1]", "phases: [11]", 1),
		"empty":           "",
	} {
		if _, err := ParseMap(doc); err == nil {
			t.Fatalf("%s: malformed map accepted", name)
		}
	}
}

func TestNeverCutViolationIsReported(t *testing.T) {
	src, m := liveSource(t), liveMap(t)
	weakened := &Map{Schema: m.Schema, index: map[string]int{}}
	for _, r := range m.Blocks {
		if strings.HasPrefix(r.Locator, "## 14.") {
			r = Rule{Locator: r.Locator, Include: IncludeOnDemand, Trigger: "x"}
		}
		if strings.HasPrefix(r.Locator, "### 15.1") {
			r = Rule{Locator: r.Locator, Include: IncludeWhen, Phases: []int{1, 2, 3, 6, 7}, Trigger: "x"}
		}
		if strings.HasPrefix(r.Locator, "## 7.") {
			r = Rule{Locator: r.Locator, Include: IncludeOnDemand, Trigger: "x"}
		}
		weakened.index[r.Locator] = len(weakened.Blocks)
		weakened.Blocks = append(weakened.Blocks, r)
	}
	rep := Check(src, weakened)
	if rep.OK || len(rep.NeverCut) != 3 {
		t.Fatalf("never-cut violations not reported: %+v", rep.NeverCut)
	}
	for _, want := range []string{"## 14.", "### 15.1", "## 7."} {
		found := false
		for _, v := range rep.NeverCut {
			if strings.HasPrefix(v, want) {
				found = true
			}
		}
		if !found {
			t.Fatalf("missing violation for %s: %v", want, rep.NeverCut)
		}
	}
}

// deck builds a temp root with the given role, protocol body and optional lock.
func deck(t *testing.T, role, body string) string {
	t.Helper()
	root := t.TempDir()
	meta := filepath.Join(root, "parley-deck", "meta")
	if err := os.MkdirAll(meta, 0o755); err != nil {
		t.Fatal(err)
	}
	if role != "" {
		if err := os.WriteFile(filepath.Join(meta, "version.json"), []byte(`{"protocolRole": "`+role+`"}`), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	if body != "" {
		if err := os.WriteFile(DeckProtocolPath(root), []byte(body), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	return root
}

const testCore = "# Core\n\n**Workspace:** `<workspace-name>`\n**Transport:** `<transport-choice>`\n**Created:** `<date>`\n\n## 1. Rule\n\nbody\n"

func TestResolveSourceRoles(t *testing.T) {
	store := protocolcore.StoreAt(t.TempDir())
	rel, err := store.Publish("1.0.0", testCore)
	if err != nil {
		t.Fatal(err)
	}

	// Source role: the live file is the authority, no store consulted, no drift concept.
	src, err := ResolveSource(deck(t, "source", "# X\n\n**Transport:** `local-dir`\n\n## 1. Rule\n\nlocal text\n"), protocolcore.Store{Root: "/nonexistent"})
	if err != nil || src.Role != "source" || src.Transport != "local-dir" || src.Drift != "" || !strings.Contains(src.Raw, "local text") {
		t.Fatalf("source role: %+v %v", src, err)
	}

	// Consumer role, verified chain, view equals the render: no drift.
	rendered, err := protocolcore.Render(rel, "", nil)
	if err != nil {
		t.Fatal(err)
	}
	root := deck(t, "consumer", rendered.Body)
	lock := protocolcore.Lock{CoreVersion: "1.0.0", CoreBodySHA256: rel.SHA256, Overlay: protocolcore.OverlayNone}
	if err := os.WriteFile(filepath.Join(root, "parley-deck", "meta", "protocol-lock.yaml"), []byte(lock.Render()), 0o644); err != nil {
		t.Fatal(err)
	}
	src, err = ResolveSource(root, store)
	if err != nil || src.Drift != "" || !strings.Contains(src.Authority, "verified") {
		t.Fatalf("consumer verified: %+v %v", src, err)
	}
	// Hand-edited view: drift is reported and an optimized packet falls back visibly.
	if err := os.WriteFile(DeckProtocolPath(root), []byte(rendered.Body+"\nhand edit\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	src, err = ResolveSource(root, store)
	if err != nil || src.Drift == "" {
		t.Fatalf("consumer drift not detected: %+v %v", src, err)
	}
	m, _ := ParseMap("schema: parley.packet-applicability/v1\nsource: x\nblocks:\n  - locator: \"# Core\"\n    include: always\n  - locator: \"## 1. Rule\"\n    include: always\n")
	if c := Build(src, m, Request{Phase: 1, Track: "fast", Transport: "local-dir", Optimize: true}); c.ContextMode != ModeFullFallback || !strings.HasPrefix(c.FallbackReason, "drift:") {
		t.Fatalf("drift must fall back: %s %q", c.ContextMode, c.FallbackReason)
	}

	// Authority failures: missing lock, missing release, wrong body hash, missing metadata,
	// unknown role, missing protocol file. None substitutes a snapshot.
	for name, root := range map[string]string{
		"consumer-no-lock": deck(t, "consumer", testCore),
		"no-version-json":  deck(t, "", testCore),
		"unknown-role":     deck(t, "upstream", testCore),
		"no-protocol":      deck(t, "source", ""),
	} {
		if _, err := ResolveSource(root, store); !errors.Is(err, ErrAuthority) {
			t.Fatalf("%s: want ErrAuthority, got %v", name, err)
		}
	}
	bad := deck(t, "consumer", testCore)
	badLock := protocolcore.Lock{CoreVersion: "1.0.0", CoreBodySHA256: strings.Repeat("0", 64), Overlay: protocolcore.OverlayNone}
	if err := os.WriteFile(filepath.Join(bad, "parley-deck", "meta", "protocol-lock.yaml"), []byte(badLock.Render()), 0o644); err != nil {
		t.Fatal(err)
	}
	if _, err := ResolveSource(bad, store); !errors.Is(err, ErrAuthority) || !strings.Contains(err.Error(), "installed release hashes") {
		t.Fatalf("body hash mismatch: %v", err)
	}
	missing := protocolcore.Lock{CoreVersion: "9.9.9", CoreBodySHA256: rel.SHA256, Overlay: protocolcore.OverlayNone}
	if err := os.WriteFile(filepath.Join(bad, "parley-deck", "meta", "protocol-lock.yaml"), []byte(missing.Render()), 0o644); err != nil {
		t.Fatal(err)
	}
	if _, err := ResolveSource(bad, store); !errors.Is(err, ErrAuthority) || !errors.Is(err, protocolcore.ErrNoRelease) {
		t.Fatalf("missing release: %v", err)
	}
}

// packetContext builds a publishable context from the fixture. Publication tests need a real
// Context because WriteBody binds the body to its attestation.
func packetContext(t *testing.T, req Request) Context {
	t.Helper()
	src, m := fixture(t)
	c := Build(src, m, req)
	if c.ContextMode == ModeRefused || c.Body == "" {
		t.Fatalf("fixture context is not publishable: %+v", c.Attestation)
	}
	return c
}

func TestPublishedBodyIsFullDigestAddressedImmutableAndIdempotent(t *testing.T) {
	root := t.TempDir()
	c := packetContext(t, Request{Phase: 1, Track: "fast"})
	if err := WriteBody(root, &c); err != nil {
		t.Fatal(err)
	}
	name := filepath.Base(c.BodyPath)
	if !strings.HasSuffix(name, c.PacketSHA256+".md") || len(c.PacketSHA256) != 64 {
		t.Fatalf("published name %q must carry the full 64-hex body digest %q", name, c.PacketSHA256)
	}
	if b, err := os.ReadFile(c.BodyPath); err != nil || string(b) != c.Body {
		t.Fatalf("published body is not the attested body: %v", err)
	}

	// Republishing identical content is idempotent: same path, no error, no staged leftovers.
	again := packetContext(t, Request{Phase: 1, Track: "fast"})
	if err := WriteBody(root, &again); err != nil || again.BodyPath != c.BodyPath {
		t.Fatalf("republication: %v %s", err, again.BodyPath)
	}
	entries, err := os.ReadDir(RuntimeDir(root))
	if err != nil {
		t.Fatal(err)
	}
	if len(entries) != 1 {
		t.Fatalf("runtime dir holds %d entries, want exactly the published body", len(entries))
	}

	// A path that already holds different bytes is refused, never truncated.
	const tampered = "tampered protocol text\n"
	if err := os.WriteFile(c.BodyPath, []byte(tampered), 0o600); err != nil {
		t.Fatal(err)
	}
	third := packetContext(t, Request{Phase: 1, Track: "fast"})
	if err := WriteBody(root, &third); err == nil || !strings.Contains(err.Error(), "different body") {
		t.Fatalf("a mismatched published body must be refused, got %v", err)
	}
	if b, _ := os.ReadFile(c.BodyPath); string(b) != tampered {
		t.Fatal("the mismatched file was overwritten instead of refused")
	}
}

func TestWriteBodyRefusesUnattestedAndEmptyBodies(t *testing.T) {
	root := t.TempDir()
	c := packetContext(t, Request{Phase: 1, Track: "fast"})
	c.Body += "\nappended after Build\n"
	if err := WriteBody(root, &c); err == nil || !strings.Contains(err.Error(), "unattested") {
		t.Fatalf("a body that does not hash to its attestation must be refused, got %v", err)
	}
	empty := packetContext(t, Request{Phase: 1, Track: "fast"})
	empty.Body, empty.PacketSHA256 = "", Hash("")
	if err := WriteBody(root, &empty); err == nil {
		t.Fatal("an empty body must not be published")
	}
	if _, err := os.Stat(filepath.Join(root, ".parley-runtime")); !os.IsNotExist(err) {
		t.Fatal("a refused publication must not create the runtime directory")
	}
}

func TestWriteBodyRefusesSymlinkedRuntimePathsAndTargets(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("creating symlinks needs privileges on Windows")
	}
	elsewhere := t.TempDir()

	// A symlinked .parley-runtime parent is refused and nothing is written through it.
	root := t.TempDir()
	if err := os.Symlink(elsewhere, filepath.Join(root, ".parley-runtime")); err != nil {
		t.Fatal(err)
	}
	c := packetContext(t, Request{Phase: 1, Track: "fast"})
	if err := WriteBody(root, &c); err == nil || !strings.Contains(err.Error(), "symlink") {
		t.Fatalf("symlinked runtime parent: %v", err)
	}
	if entries, _ := os.ReadDir(elsewhere); len(entries) != 0 {
		t.Fatal("wrote through the symlinked runtime parent")
	}

	// A symlinked protocol-packets directory is refused.
	root = t.TempDir()
	if err := os.Mkdir(filepath.Join(root, ".parley-runtime"), 0o700); err != nil {
		t.Fatal(err)
	}
	if err := os.Symlink(elsewhere, RuntimeDir(root)); err != nil {
		t.Fatal(err)
	}
	c = packetContext(t, Request{Phase: 1, Track: "fast"})
	if err := WriteBody(root, &c); err == nil || !strings.Contains(err.Error(), "symlink") {
		t.Fatalf("symlinked runtime dir: %v", err)
	}

	// A regular file where the runtime directory belongs is refused.
	root = t.TempDir()
	if err := os.WriteFile(filepath.Join(root, ".parley-runtime"), []byte("x"), 0o600); err != nil {
		t.Fatal(err)
	}
	c = packetContext(t, Request{Phase: 1, Track: "fast"})
	if err := WriteBody(root, &c); err == nil || !strings.Contains(err.Error(), "not a directory") {
		t.Fatalf("runtime parent is a file: %v", err)
	}

	// A symlinked TARGET is refused, and the file it points at is untouched.
	scratch := t.TempDir()
	probe := packetContext(t, Request{Phase: 1, Track: "fast"})
	if err := WriteBody(scratch, &probe); err != nil {
		t.Fatal(err)
	}
	root = t.TempDir()
	if err := os.MkdirAll(RuntimeDir(root), 0o700); err != nil {
		t.Fatal(err)
	}
	victim := filepath.Join(elsewhere, "victim.md")
	if err := os.WriteFile(victim, []byte("victim\n"), 0o600); err != nil {
		t.Fatal(err)
	}
	if err := os.Symlink(victim, filepath.Join(RuntimeDir(root), filepath.Base(probe.BodyPath))); err != nil {
		t.Fatal(err)
	}
	c = packetContext(t, Request{Phase: 1, Track: "fast"})
	if err := WriteBody(root, &c); err == nil || !strings.Contains(err.Error(), "symlink") {
		t.Fatalf("symlinked target: %v", err)
	}
	if b, _ := os.ReadFile(victim); string(b) != "victim\n" {
		t.Fatal("wrote through a symlinked target")
	}
}

func TestPublicationDirectoryAndBodyArePrivate(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("POSIX permission bits are not meaningful on Windows")
	}
	root := t.TempDir()
	c := packetContext(t, Request{Phase: 2, Track: "standard"})
	if err := WriteBody(root, &c); err != nil {
		t.Fatal(err)
	}
	for _, p := range []string{filepath.Join(root, ".parley-runtime"), RuntimeDir(root)} {
		fi, err := os.Stat(p)
		if err != nil {
			t.Fatal(err)
		}
		if fi.Mode().Perm() != 0o700 {
			t.Fatalf("%s has mode %04o, want 0700", p, fi.Mode().Perm())
		}
	}
	fi, err := os.Stat(c.BodyPath)
	if err != nil {
		t.Fatal(err)
	}
	if fi.Mode().Perm() != 0o600 {
		t.Fatalf("%s has mode %04o, want 0600", c.BodyPath, fi.Mode().Perm())
	}

	// A pre-existing group/world-accessible publication directory is restricted before use;
	// a directory another local user can rewrite carries no immutability guarantee.
	loose := t.TempDir()
	if err := os.MkdirAll(RuntimeDir(loose), 0o700); err != nil {
		t.Fatal(err)
	}
	if err := os.Chmod(RuntimeDir(loose), 0o777); err != nil {
		t.Fatal(err)
	}
	c2 := packetContext(t, Request{Phase: 2, Track: "standard"})
	if err := WriteBody(loose, &c2); err != nil {
		t.Fatal(err)
	}
	fi, err = os.Stat(RuntimeDir(loose))
	if err != nil {
		t.Fatal(err)
	}
	if fi.Mode().Perm()&0o077 != 0 {
		t.Fatalf("a world-writable publication directory was used as is (%04o)", fi.Mode().Perm())
	}
}

func TestConcurrentPublicationPublishesOneCompleteBody(t *testing.T) {
	root := t.TempDir()
	src, m := fixture(t) // built once: t.Fatal must not be called from a spawned goroutine
	const writers = 8

	var wg sync.WaitGroup
	errs := make([]error, writers)
	paths := make([]string, writers)
	start := make(chan struct{})
	for i := 0; i < writers; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			c := Build(src, m, Request{Phase: 1, Track: "fast"})
			<-start
			errs[i] = WriteBody(root, &c)
			paths[i] = c.BodyPath
		}(i)
	}
	close(start)
	wg.Wait()

	want := Build(src, m, Request{Phase: 1, Track: "fast"}).Body
	for i, err := range errs {
		if err != nil {
			t.Fatalf("writer %d: %v", i, err)
		}
		if paths[i] != paths[0] {
			t.Fatalf("writer %d published %s, writer 0 published %s", i, paths[i], paths[0])
		}
	}
	entries, err := os.ReadDir(RuntimeDir(root))
	if err != nil {
		t.Fatal(err)
	}
	if len(entries) != 1 {
		var names []string
		for _, e := range entries {
			names = append(names, e.Name())
		}
		t.Fatalf("concurrent publication left %d files (%v); want one complete body and no staged leftovers", len(entries), names)
	}
	// Whatever a concurrent reader would have opened must be the whole body, never a prefix.
	b, err := os.ReadFile(filepath.Join(RuntimeDir(root), entries[0].Name()))
	if err != nil {
		t.Fatal(err)
	}
	if string(b) != want {
		t.Fatalf("published body is %d bytes, want the complete %d", len(b), len(want))
	}

	// Concurrent publication of DIFFERENT bodies keeps both, each complete.
	var wg2 sync.WaitGroup
	derr := make([]error, 2)
	dpaths := make([]string, 2)
	for i, phase := range []int{1, 2} {
		wg2.Add(1)
		go func(i, phase int) {
			defer wg2.Done()
			c := Build(src, m, Request{Phase: phase, Track: "fast", Optimize: true})
			derr[i] = WriteBody(root, &c)
			dpaths[i] = c.BodyPath
		}(i, phase)
	}
	wg2.Wait()
	for i, err := range derr {
		if err != nil {
			t.Fatalf("distinct writer %d: %v", i, err)
		}
	}
	if dpaths[0] == dpaths[1] {
		t.Fatalf("phase 1 and phase 2 packets share the path %s", dpaths[0])
	}
	entries, err = os.ReadDir(RuntimeDir(root))
	if err != nil {
		t.Fatal(err)
	}
	if len(entries) != 3 {
		t.Fatalf("runtime dir holds %d files, want the full body plus two distinct packets", len(entries))
	}
	for _, e := range entries {
		p := filepath.Join(RuntimeDir(root), e.Name())
		b, err := os.ReadFile(p)
		if err != nil {
			t.Fatal(err)
		}
		if !strings.HasSuffix(e.Name(), Hash(string(b))+".md") {
			t.Fatalf("%s does not hold the body its name addresses", e.Name())
		}
	}
}

// TestConcurrentPublicationCreatesThePrivateRuntimeDirectoryOnce covers the create race that
// broke repeated concurrent publication: two publishers both Lstat a missing .parley-runtime,
// both call Mkdir, and the loser gets EEXIST. Failing on EEXIST breaks a legitimate concurrent
// publisher; trusting EEXIST would accept whatever entry actually appeared, including a symlink
// or a file planted in the window. Every success must therefore be backed by a real, private
// directory. The storm is repeated on a fresh root so the window is hit rather than assumed.
func TestConcurrentPublicationCreatesThePrivateRuntimeDirectoryOnce(t *testing.T) {
	const rounds, creators = 10, 12
	for round := 0; round < rounds; round++ {
		root := t.TempDir()
		var wg sync.WaitGroup
		errs := make([]error, creators)
		dirs := make([]string, creators)
		start := make(chan struct{})
		for i := 0; i < creators; i++ {
			wg.Add(1)
			go func(i int) {
				defer wg.Done()
				<-start
				dirs[i], errs[i] = publicationDir(root)
			}(i)
		}
		close(start)
		wg.Wait()

		for i, err := range errs {
			if err != nil {
				t.Fatalf("round %d creator %d: %v", round, i, err)
			}
			if dirs[i] != RuntimeDir(root) {
				t.Fatalf("round %d creator %d resolved %s, want %s", round, i, dirs[i], RuntimeDir(root))
			}
		}
		for _, p := range []string{filepath.Join(root, ".parley-runtime"), RuntimeDir(root)} {
			fi, err := os.Lstat(p)
			if err != nil {
				t.Fatalf("round %d: every publicationDir succeeded but %s: %v", round, p, err)
			}
			if fi.Mode()&os.ModeSymlink != 0 || !fi.IsDir() {
				t.Fatalf("round %d: %s is %s, not a real directory", round, p, fi.Mode())
			}
			if runtime.GOOS != "windows" && fi.Mode().Perm()&0o077 != 0 {
				t.Fatalf("round %d: %s is group/world accessible (%04o)", round, p, fi.Mode().Perm())
			}
		}
	}
}

// TestVerifyPublishedRefusesASwappedSymlinkAndDifferentBytes pins the validation both EEXIST
// paths funnel into. The link-EEXIST caller never went through existingBody, so verifyPublished
// itself must refuse a symlink: reading with os.ReadFile followed one, which let a file swapped
// in during the race "prove" the publication and hand a launch text that is not the attested text.
func TestVerifyPublishedRefusesASwappedSymlinkAndDifferentBytes(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("creating symlinks needs privileges on Windows")
	}
	dir := t.TempDir()
	const body = "attested protocol body\n"
	path := filepath.Join(dir, "packet.md")

	// The attested bytes: idempotent republication, no error.
	if err := os.WriteFile(path, []byte(body), 0o600); err != nil {
		t.Fatal(err)
	}
	if err := verifyPublished(path, body); err != nil {
		t.Fatalf("identical bytes must validate: %v", err)
	}

	// Different bytes: refused, and nothing is written.
	const other = "someone else's body\n"
	if err := os.WriteFile(path, []byte(other), 0o600); err != nil {
		t.Fatal(err)
	}
	if err := verifyPublished(path, body); err == nil || !strings.Contains(err.Error(), "different body") {
		t.Fatalf("mismatched bytes must be refused, got %v", err)
	}
	if b, _ := os.ReadFile(path); string(b) != other {
		t.Fatal("verifyPublished modified the file it was asked to validate")
	}

	// A symlink is refused even when its target holds EXACTLY the attested bytes: content read
	// through a link says nothing about what stands at the published path.
	attacker := filepath.Join(t.TempDir(), "attacker.md")
	if err := os.WriteFile(attacker, []byte(body), 0o600); err != nil {
		t.Fatal(err)
	}
	link := filepath.Join(dir, "linked.md")
	if err := os.Symlink(attacker, link); err != nil {
		t.Fatal(err)
	}
	if err := verifyPublished(link, body); err == nil || !strings.Contains(err.Error(), "symlink") {
		t.Fatalf("a body accepted through a symlink: %v", err)
	}

	// A directory (or any non-regular entry) is not a published body either.
	if err := os.Mkdir(filepath.Join(dir, "dir.md"), 0o700); err != nil {
		t.Fatal(err)
	}
	if err := verifyPublished(filepath.Join(dir, "dir.md"), body); err == nil || !strings.Contains(err.Error(), "not a regular file") {
		t.Fatalf("non-regular entry: %v", err)
	}
}

// TestPublicationFailsClosedWhenHardLinksAreUnsupported pins the removal of the direct-write
// fallback. An O_CREATE|O_EXCL write at the target never clobbers, but it publishes an empty file
// and fills it afterwards, so a launch reading in that window gets a partial body under a complete
// attestation. Where the link is unavailable, publication must fail and leave nothing behind.
func TestPublicationFailsClosedWhenHardLinksAreUnsupported(t *testing.T) {
	root := t.TempDir()
	c := packetContext(t, Request{Phase: 1, Track: "fast"})
	restore := linkFile
	linkFile = func(string, string) error { return errors.New("operation not supported") }
	defer func() { linkFile = restore }()

	err := WriteBody(root, &c)
	if err == nil || !strings.Contains(err.Error(), "atomically") {
		t.Fatalf("an unsupported hard link must fail publication, got %v", err)
	}
	if c.BodyPath != "" {
		t.Fatalf("a failed publication recorded %s as the body path", c.BodyPath)
	}
	entries, rerr := os.ReadDir(RuntimeDir(root))
	if rerr != nil {
		t.Fatal(rerr)
	}
	if len(entries) != 0 {
		var names []string
		for _, e := range entries {
			names = append(names, e.Name())
		}
		t.Fatalf("failed publication left %v; it must leave neither a partial body nor a staged file", names)
	}
}

// TestStagedBodyIsFlushedOnThisFilesystem exercises the real filesystem behind t.TempDir().
// Darwin's File.Sync issues F_FULLFSYNC, which shared volumes (a TMPDIR under /Volumes/...)
// reject with ENOTTY even though ordinary fsync succeeds there; publication must neither fail
// on such a volume nor stop flushing. Run it with TMPDIR pointed at one to see the difference.
func TestStagedBodyIsFlushedOnThisFilesystem(t *testing.T) {
	dir := t.TempDir()
	f, err := os.CreateTemp(dir, ".staged-packet-*")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(f.Name())
	const body = "packet body\n"
	if err := writeAndClose(f, body); err != nil {
		t.Fatalf("staging a body in %s failed: %v", dir, err)
	}
	if b, err := os.ReadFile(f.Name()); err != nil || string(b) != body {
		t.Fatalf("staged body is %q (%v), want %q", b, err, body)
	}
	if runtime.GOOS != "windows" {
		fi, err := os.Stat(f.Name())
		if err != nil {
			t.Fatal(err)
		}
		if fi.Mode().Perm() != 0o600 {
			t.Fatalf("staged body has mode %04o, want 0600 before it is ever linked into place", fi.Mode().Perm())
		}
	}
}

func TestMismatchedFenceMarkerDoesNotHideHeadings(t *testing.T) {
	// A ~~~ line inside a ``` block is content. Toggling on it would swallow every following
	// heading into the preceding block, and an omitted block would then carry text the
	// omission index attributes to the wrong locator.
	raw := "# T\n\n## A\n\n```\n~~~\n## Not a heading\n```\n\n## B\n\nb text\n"
	var locs []string
	for _, b := range Parse(raw) {
		locs = append(locs, b.Locator)
	}
	if got, want := strings.Join(locs, "|"), "# T|## A|## B"; got != want {
		t.Fatalf("locators %q, want %q", got, want)
	}
}

func TestRenderWritesTheBodyUnderTheRuntimeDir(t *testing.T) {
	root := deck(t, "source", fixtureSource)
	if err := os.WriteFile(MapPath(root), []byte(fixtureMap), 0o644); err != nil {
		t.Fatal(err)
	}
	c, err := Render(root, protocolcore.Store{Root: "/nonexistent"}, Request{Phase: 1, Track: "fast"})
	if err != nil || c.ContextMode != ModeFull {
		t.Fatalf("%v %+v", err, c.Attestation)
	}
	if !strings.HasPrefix(c.BodyPath, RuntimeDir(root)) {
		t.Fatalf("body written to %s, want under %s", c.BodyPath, RuntimeDir(root))
	}
	b, err := os.ReadFile(c.BodyPath)
	if err != nil || string(b) != fixtureSource {
		t.Fatalf("full body not written verbatim: %v", err)
	}
	p, err := Render(root, protocolcore.Store{Root: "/nonexistent"}, Request{Phase: 1, Track: "fast", Optimize: true})
	if err != nil || p.ContextMode != ModePacket || !strings.Contains(filepath.Base(p.BodyPath), "packet-phase1-fast-") {
		t.Fatalf("%v %+v %s", err, p.Attestation, p.BodyPath)
	}
	// A malformed map is visible as the fallback reason, never silently "no map".
	if err := os.WriteFile(MapPath(root), []byte("schema: nope\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	f, err := Render(root, protocolcore.Store{Root: "/nonexistent"}, Request{Phase: 1, Track: "fast", Optimize: true})
	if err != nil || f.ContextMode != ModeFullFallback || !strings.HasPrefix(f.FallbackReason, "malformed-map:") {
		t.Fatalf("%v %+v", err, f.Attestation)
	}
	// Missing authority is an error and nothing is written.
	if _, err := Render(deck(t, "", fixtureSource), protocolcore.Store{Root: "/nonexistent"}, Request{Phase: 1, Track: "fast"}); !errors.Is(err, ErrAuthority) {
		t.Fatalf("want ErrAuthority, got %v", err)
	}
}
