package app

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"testing"
	"time"

	"parley-deck-cli/internal/agents"
	"parley-deck-cli/internal/config"
	"parley-deck-cli/internal/protocol"
	"parley-deck-cli/internal/runner"
	"parley-deck-cli/internal/store"
)

// Designated-implementer tests (idea meta-protocol-change-designated-implementer).
// Every test is written to FAIL if the rule it pins is removed.

func designationRoot(t *testing.T) string {
	t.Helper()
	root := t.TempDir()
	if err := protocol.InitWorkspace(root); err != nil {
		t.Fatal(err)
	}
	// The protocol-task launch path requires a declared authority role.
	meta := filepath.Join(root, protocol.DeckDir, "meta")
	if err := os.MkdirAll(meta, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(meta, "version.json"), []byte(`{"protocolRole":"source"}`), 0o600); err != nil {
		t.Fatal(err)
	}
	return root
}

func writeDesignationPrompt(t *testing.T, root, slug, frontmatter string) string {
	t.Helper()
	ideaDir := filepath.Join(root, protocol.DeckDir, "ideas", slug)
	if err := os.MkdirAll(ideaDir, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(ideaDir, "00-prompt.md"), []byte(frontmatter), 0o644); err != nil {
		t.Fatal(err)
	}
	return ideaDir
}

func designationPrompt(slug string, participants []string, extra string) string {
	return "---\nidea: " + slug + "\nauthor: user\ntrack: deliberation\nparticipants: [" + strings.Join(participants, ", ") + "]\n" + extra + "status: final\n---\n"
}

// setCentralDefaultImplementer writes a central layer (PARLEY_HOME/agents.toml) with
// the given [defaults] body and isolates the env layer.
func setCentralDefaults(t *testing.T, defaultsBody string) {
	t.Helper()
	home := t.TempDir()
	t.Setenv(config.EnvParleyHome, home)
	t.Setenv(config.EnvAgentConfig, "")
	if defaultsBody == "" {
		return
	}
	if err := os.WriteFile(filepath.Join(home, "agents.toml"), []byte("[defaults]\n"+defaultsBody), 0o644); err != nil {
		t.Fatal(err)
	}
}

func designationOps(t *testing.T, root, slug, ideaDir string, participants, discovered []string, out io.Writer) driverImplOps {
	t.Helper()
	var agentsList []agents.Discovery
	for _, id := range discovered {
		agentsList = append(agentsList, agents.Discovery{
			Spec:  agents.Spec{ID: id, HeadlessArgs: []string{"-test.run=TestDesignationFakeImplHelper", "--", "parley-fake-designation-impl"}, PromptMode: agents.PromptStdin},
			Path:  os.Args[0],
			Found: true,
		})
	}
	ops := newDriverImplOps(runner.Options{
		Root:    root,
		RunID:   "designation-test",
		Idea:    protocol.IdeaStatus{Slug: slug, Path: ideaDir, Participants: participants},
		Agents:  agentsList,
		Timeout: 30 * time.Second,
	}, root, slug, ideaDir, participants, out)
	impl, ok := ops.(driverImplOps)
	if !ok {
		t.Fatalf("newDriverImplOps returned %T, want driverImplOps", ops)
	}
	return impl
}

func TestDesignationFakeImplHelper(t *testing.T) {
	helper := false
	for _, a := range os.Args {
		if a == "parley-fake-designation-impl" {
			helper = true
		}
	}
	if !helper {
		return
	}
	input, _ := io.ReadAll(os.Stdin)
	out := regexp.MustCompile(`(?m)^- Create exactly this file: (.+)$`).FindStringSubmatch(string(input))
	idea := regexp.MustCompile(`(?m)^idea: (\S+)$`).FindStringSubmatch(string(input))
	if len(out) != 2 || len(idea) != 2 {
		fmt.Fprintf(os.Stderr, "impl prompt missing path/idea:\n%s", string(input))
		os.Exit(2)
	}
	body := "---\nidea: " + idea[1] + "\nstatus: implemented\n---\n\n## Summary of work\nbuilt by the fake implementer\n"
	if err := os.WriteFile(out[1], []byte(body), 0o644); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(2)
	}
	os.Exit(0)
}

// T-8 / AC-4: an undesignated deck is byte-identical — the legacy positional
// dispatch, the legacy stdout line with NOTHING added, and no
// agent.implementer_resolved event in the stream.
func TestUnsetPathIsByteIdentical(t *testing.T) {
	setCentralDefaults(t, "")
	root := designationRoot(t)
	ideaDir := writeDesignationPrompt(t, root, "demo", designationPrompt("demo", []string{"zz-first", "aa-rev"}, ""))
	var out bytes.Buffer
	ops := designationOps(t, root, "demo", ideaDir, []string{"zz-first", "aa-rev"}, []string{"zz-first", "aa-rev"}, &out)
	if ops.implDesignated || ops.implSource != "" || ops.roleErr != "" || ops.dispatchErr != "" {
		t.Fatalf("unset path must carry zero designation state, got %+v", ops)
	}
	if ops.implementer != "zz-first" {
		t.Fatalf("unset dispatch must stay positional (eligible[0]), got %q", ops.implementer)
	}

	// Drive a real dispatch through the fake implementer and inspect the record.
	st := store.New(filepath.Join(root, protocol.DeckDir, "runs", "designation-test"))
	ops.base.Store = st
	out.Reset()
	if err := ops.Implement(context.Background()); err != nil {
		t.Fatalf("unset dispatch failed: %v", err)
	}
	if got := out.String(); got != "driver: implementing via zz-first ...\n" {
		t.Fatalf("unset stdout must be byte-identical, got %q", got)
	}
	events, err := st.Load()
	if err != nil {
		t.Fatal(err)
	}
	for _, e := range events {
		if e.Type == "agent.implementer_resolved" {
			t.Fatalf("no agent.implementer_resolved event may be emitted on the unset path: %+v", e)
		}
	}
}

// AC-18a: a per-idea `implementer:` designation dispatches the designee end to end.
func TestTier2DesignationDispatchesDesignee(t *testing.T) {
	setCentralDefaults(t, "")
	root := designationRoot(t)
	ideaDir := writeDesignationPrompt(t, root, "demo", designationPrompt("demo", []string{"aa-first", "zz-impl"}, "implementer: zz-impl\n"))
	var out bytes.Buffer
	ops := designationOps(t, root, "demo", ideaDir, []string{"aa-first", "zz-impl"}, []string{"aa-first", "zz-impl"}, &out)
	if ops.implementer != "zz-impl" || ops.implSource != protocol.SourceImplementerDesignation || !ops.implDesignated {
		t.Fatalf("tier-2 designation must dispatch the designee, got implementer=%q source=%q", ops.implementer, ops.implSource)
	}
	if ops.roleErr != "" || ops.dispatchErr != "" {
		t.Fatalf("valid designation must not gate: roleErr=%q dispatchErr=%q", ops.roleErr, ops.dispatchErr)
	}
	st := store.New(filepath.Join(root, protocol.DeckDir, "runs", "designation-test"))
	ops.base.Store = st
	out.Reset()
	if err := ops.Implement(context.Background()); err != nil {
		t.Fatalf("designated dispatch failed: %v", err)
	}
	if !strings.Contains(out.String(), "driver: implementing via zz-impl ...\n") ||
		!strings.Contains(out.String(), "driver: implementer resolved: zz-impl (source: designation)\n") {
		t.Fatalf("designation stdout lines missing: %q", out.String())
	}
	events, err := st.Load()
	if err != nil {
		t.Fatal(err)
	}
	found := false
	for _, e := range events {
		if e.Type == "agent.implementer_resolved" {
			found = true
			if e.Data["idea"] != "demo" || e.Data["implementer"] != "zz-impl" || e.Data["source"] != "designation" {
				t.Fatalf("event payload wrong: %+v", e.Data)
			}
		}
	}
	if !found {
		t.Fatal("designation dispatch must record agent.implementer_resolved")
	}
}

// AC-18b: the same key at two config layers resolves to the higher layer's id and
// records the source.
func TestTwoLayerGlobalDefaultDispatchesHigherLayer(t *testing.T) {
	setCentralDefaults(t, "default_implementer = \"zz-low\"\n")
	root := designationRoot(t)
	deck := filepath.Join(root, protocol.DeckDir)
	if err := os.WriteFile(filepath.Join(deck, "agents.toml"), []byte("[defaults]\ndefault_implementer = \"zz-high\"\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	ideaDir := writeDesignationPrompt(t, root, "demo", designationPrompt("demo", []string{"aa-first", "zz-low", "zz-high"}, ""))
	var out bytes.Buffer
	ops := designationOps(t, root, "demo", ideaDir, []string{"aa-first", "zz-low", "zz-high"}, []string{"aa-first", "zz-low", "zz-high"}, &out)
	if ops.implementer != "zz-high" || ops.implSource != protocol.SourceImplementerGlobalDefault {
		t.Fatalf("higher layer must win dispatch, got implementer=%q source=%q", ops.implementer, ops.implSource)
	}
	st := store.New(filepath.Join(root, protocol.DeckDir, "runs", "designation-test"))
	ops.base.Store = st
	if err := ops.Implement(context.Background()); err != nil {
		t.Fatalf("global-default dispatch failed: %v", err)
	}
	events, err := st.Load()
	if err != nil {
		t.Fatal(err)
	}
	found := false
	for _, e := range events {
		if e.Type == "agent.implementer_resolved" && e.Data["implementer"] == "zz-high" && e.Data["source"] == "global-default" {
			found = true
		}
	}
	if !found {
		t.Fatal("the event must record the higher layer's id with source global-default")
	}
}

// T-2 / AC-6: `implementer: none` suppresses tier 3 and is distinguishable from an
// absent key (without the pin the opt-out silently collapses into "absent").
func TestNoneSuppressesTier3(t *testing.T) {
	setCentralDefaults(t, "default_implementer = \"zz-global\"\n")
	root := designationRoot(t)
	participants := []string{"aa-first", "zz-global"}

	withNone := writeDesignationPrompt(t, root, "optout", designationPrompt("optout", participants, "implementer: none\n"))
	opsNone := designationOps(t, root, "optout", withNone, participants, participants, io.Discard)
	if opsNone.implementer != "aa-first" || opsNone.implSource != protocol.SourceImplementerNone {
		t.Fatalf("none must suppress tier 3 and fall to today's chain, got implementer=%q source=%q", opsNone.implementer, opsNone.implSource)
	}
	if !opsNone.implDesignated {
		t.Fatal("a present opt-out is still a present designation (the event records it)")
	}
	if opsNone.roleErr != "" || opsNone.dispatchErr != "" {
		t.Fatal("the opt-out never gates")
	}

	absent := writeDesignationPrompt(t, root, "plain", designationPrompt("plain", participants, ""))
	opsAbsent := designationOps(t, root, "plain", absent, participants, participants, io.Discard)
	if opsAbsent.implementer != "zz-global" || opsAbsent.implSource != protocol.SourceImplementerGlobalDefault {
		t.Fatalf("absent key lets tier 3 fire, got implementer=%q source=%q", opsAbsent.implementer, opsAbsent.implSource)
	}
}

// T-4 / AC-8: the tier-2 unavailability gate blocks dispatch with the three exits
// pre-built, and each exit clears it.
func TestTier2UnavailabilityGateAndExits(t *testing.T) {
	setCentralDefaults(t, "")
	participants := []string{"aa-first", "zz-impl"}

	newOps := func(t *testing.T, extra string) driverImplOps {
		root := designationRoot(t)
		ideaDir := writeDesignationPrompt(t, root, "demo", designationPrompt("demo", participants, extra))
		// zz-impl is NOT discovered — unavailable at the ping.
		return designationOps(t, root, "demo", ideaDir, participants, []string{"aa-first"}, &bytes.Buffer{})
	}

	ops := newOps(t, "implementer: zz-impl\n")
	if ops.dispatchErr == "" {
		t.Fatal("a valid but unavailable tier-2 designee must gate dispatch")
	}
	for _, want := range []string{"zz-impl", "implementer_waived", "implementer: none", "re-designate"} {
		if !strings.Contains(ops.dispatchErr, want) {
			t.Fatalf("gate must pre-build the three exits; missing %q in %q", want, ops.dispatchErr)
		}
	}
	if ops.roleErr != "" {
		t.Fatalf("availability is not a validity gate; roleErr must stay empty, got %q", ops.roleErr)
	}
	if err := ops.Implement(context.Background()); err == nil || !strings.Contains(err.Error(), "unavailable") {
		t.Fatalf("Implement must escalate the availability gate, got %v", err)
	}
	if err := ops.Fixup(context.Background(), 1); err == nil {
		t.Fatal("Fixup dispatches the implementer too and must escalate the same gate")
	}

	// Exit 1: re-designate to an available participant.
	if ops := newOps(t, "implementer: aa-first\n"); ops.dispatchErr != "" || ops.implementer != "aa-first" {
		t.Fatalf("re-designation must clear the gate, got implementer=%q dispatchErr=%q", ops.implementer, ops.dispatchErr)
	}
	// Exit 2: a confirmed waiver falls through to today's chain, loudly.
	ops = newOps(t, "implementer: zz-impl\nimplementer_waived: zz-impl — offline — confirmed 2026-09-25\n")
	if ops.dispatchErr != "" || ops.implementer != "aa-first" || ops.implSource != protocol.SourceImplementerFallThroughUnavailable {
		t.Fatalf("the confirmed waiver must clear the gate into a recorded fall-through, got %+v", ops)
	}
	// An UNCONFIRMED waiver line is not an exit.
	if ops := newOps(t, "implementer: zz-impl\nimplementer_waived: zz-impl — offline\n"); ops.dispatchErr == "" {
		t.Fatal("an unconfirmed waiver must not clear the gate")
	}
	// Exit 3: implementer: none.
	if ops := newOps(t, "implementer: none\n"); ops.dispatchErr != "" || ops.implementer != "aa-first" || ops.implSource != protocol.SourceImplementerNone {
		t.Fatalf("implementer: none must clear the gate, got %+v", ops)
	}
}

// T-5 / AC-11: pin-vs-designation disagreement escalates, never silently honours
// either, and the implementer_reassigned: path re-pins only with a confirmed record.
func TestPinDesignationConflictEscalates(t *testing.T) {
	setCentralDefaults(t, "")
	participants := []string{"aa-first", "zz-impl"}
	newOps := func(t *testing.T, extra string) driverImplOps {
		root := designationRoot(t)
		ideaDir := writeDesignationPrompt(t, root, "demo", designationPrompt("demo", participants, extra))
		pin := "---\nidea: demo\nstatus: implemented\nimplementer: aa-first\n---\n\n## Summary of work\npartial\n"
		if err := os.WriteFile(filepath.Join(ideaDir, "IMPLEMENTATION.md"), []byte(pin), 0o644); err != nil {
			t.Fatal(err)
		}
		return designationOps(t, root, "demo", ideaDir, participants, participants, &bytes.Buffer{})
	}

	ops := newOps(t, "implementer: zz-impl\n")
	if ops.roleErr == "" || !strings.Contains(ops.roleErr, "implementer conflict") {
		t.Fatalf("pin/designation disagreement must escalate, got roleErr=%q", ops.roleErr)
	}
	if !strings.Contains(ops.roleErr, "implementer_reassigned") {
		t.Fatalf("the escalation must pre-build the recorded exit, got %q", ops.roleErr)
	}
	if err := ops.Implement(context.Background()); err == nil {
		t.Fatal("the conflict must block dispatch")
	}

	// An unconfirmed reassignment line does not re-pin.
	if ops := newOps(t, "implementer: zz-impl\nimplementer_reassigned: aa-first to zz-impl — rotation\n"); ops.roleErr == "" {
		t.Fatal("an unconfirmed reassignment record must not clear the conflict")
	}
	// A record naming a different pair does not re-pin.
	if ops := newOps(t, "implementer: zz-impl\nimplementer_reassigned: aa-first to zz-other — rotation — confirmed 2026-09-25\n"); ops.roleErr == "" {
		t.Fatal("a reassignment record for a different pair must not clear this conflict")
	}
	// The confirmed record re-pins: the designation proceeds.
	ops = newOps(t, "implementer: zz-impl\nimplementer_reassigned: aa-first to zz-impl — rotation — confirmed 2026-09-25\n")
	if ops.roleErr != "" || ops.implementer != "zz-impl" {
		t.Fatalf("the confirmed reassignment must re-pin to the designation, got implementer=%q roleErr=%q", ops.implementer, ops.roleErr)
	}

	// Agreement never escalates (R27).
	ops = newOps(t, "implementer: aa-first\n")
	if ops.roleErr != "" || ops.implementer != "aa-first" || ops.implSource != protocol.SourceImplementerPin {
		t.Fatalf("pin and designation agreeing must resolve as the pin, got %+v", ops)
	}
}

// T-7 / AC-9: every dormant tier-3 state falls through to today's chain WITHOUT
// gating, with the source value distinguishing the cause where a designation was
// present; the absent case emits nothing.
func TestDormantTier3StatesFallThrough(t *testing.T) {
	participants := []string{"aa-first", "zz-impl"}

	newOps := func(t *testing.T, defaultsBody, extra string, discovered []string) driverImplOps {
		setCentralDefaults(t, defaultsBody)
		root := designationRoot(t)
		ideaDir := writeDesignationPrompt(t, root, "demo", designationPrompt("demo", participants, extra))
		return designationOps(t, root, "demo", ideaDir, participants, discovered, &bytes.Buffer{})
	}

	// (a) key absent, nothing set anywhere: the pure unset path.
	ops := newOps(t, "", "", participants)
	if ops.implDesignated || ops.implementer != "aa-first" || ops.roleErr != "" || ops.dispatchErr != "" {
		t.Fatalf("absent: got %+v", ops)
	}

	// (b) tier-3 names a non-participant: inapplicable, not invalid (R22).
	ops = newOps(t, "default_implementer = \"zz-ghost\"\n", "", participants)
	if ops.roleErr != "" || ops.dispatchErr != "" || ops.implementer != "aa-first" ||
		ops.implSource != protocol.SourceImplementerFallThroughInapplicable || !ops.implDesignated || ops.implLine == "" {
		t.Fatalf("non-participant tier-3 must fall through with a notice, got %+v", ops)
	}

	// (c) tier-3 is "none": the documented suppressor (R9).
	ops = newOps(t, "default_implementer = \"none\"\n", "", participants)
	if ops.implementer != "aa-first" || ops.implSource != protocol.SourceImplementerNone || !ops.implDesignated || ops.roleErr != "" {
		t.Fatalf("tier-3 none: got %+v", ops)
	}

	// (d) tier-3 names the idea's own declared non-participating facilitator (R20):
	// inapplicable at tier 3, exactly like a non-participant.
	setCentralDefaults(t, "default_implementer = \"zz-fac\"\n")
	root := designationRoot(t)
	ideaDir := writeDesignationPrompt(t, root, "demo", designationPrompt("demo", participants, "facilitator: zz-fac\n"))
	ops = designationOps(t, root, "demo", ideaDir, participants, participants, &bytes.Buffer{})
	if ops.roleErr != "" || ops.dispatchErr != "" || ops.implementer != "aa-first" ||
		ops.implSource != protocol.SourceImplementerFallThroughInapplicable || !ops.implDesignated {
		t.Fatalf("own declared facilitator at tier 3 must fall through, not gate: %+v", ops)
	}

	// (e) tier-3 names an eligible participant that fails the ping (R19): one-line
	// notice and fall-through, no gate — and the notice names the three exits.
	ops = newOps(t, "default_implementer = \"zz-impl\"\n", "", []string{"aa-first"}) // zz-impl undiscovered
	if ops.roleErr != "" || ops.dispatchErr != "" || ops.implementer != "aa-first" ||
		ops.implSource != protocol.SourceImplementerFallThroughUnavailable || !ops.implDesignated {
		t.Fatalf("unavailable tier-3 must fall through, got %+v", ops)
	}
	for _, want := range []string{"zz-impl", "implementer: none", "implementer_waived"} {
		if !strings.Contains(ops.implLine, want) {
			t.Fatalf("the fall-through notice must name the designee and the exits; missing %q in %q", want, ops.implLine)
		}
	}

	// Same id at tier 2 stays a hard gate (R20's tier split).
	ops = newOps(t, "", "facilitator: zz-fac\nimplementer: zz-fac\n", participants)
	if ops.roleErr == "" {
		t.Fatal("the idea's own declared facilitator at tier 2 must hard-gate")
	}
}

// AC-10: a malformed tier-3 value keeps its hard failure and does not fall through.
func TestMalformedTier3HardFails(t *testing.T) {
	setCentralDefaults(t, "default_implementer = \"zz-impl # copied from a comment\"\n")
	root := designationRoot(t)
	participants := []string{"aa-first", "zz-impl"}
	ideaDir := writeDesignationPrompt(t, root, "demo", designationPrompt("demo", participants, ""))
	ops := designationOps(t, root, "demo", ideaDir, participants, participants, &bytes.Buffer{})
	if ops.roleErr == "" {
		t.Fatal("a malformed tier-3 value must hard-fail, never fall through")
	}
	if err := ops.Implement(context.Background()); err == nil {
		t.Fatal("the malformed tier-3 value must block dispatch")
	}
}

// AC-5/R4 at the gate: a malformed tier-2 value is invalid, never repaired into an id.
func TestMalformedTier2Gates(t *testing.T) {
	setCentralDefaults(t, "")
	root := designationRoot(t)
	participants := []string{"aa-first", "zz-impl"}
	ideaDir := writeDesignationPrompt(t, root, "demo", designationPrompt("demo", participants, "implementer: zz-impl  # from the global default\n"))
	ops := designationOps(t, root, "demo", ideaDir, participants, participants, &bytes.Buffer{})
	if ops.roleErr == "" || !strings.Contains(ops.roleErr, "not an eligible participant") {
		t.Fatalf("the malformed value must fail closed, got roleErr=%q", ops.roleErr)
	}

	// Present-empty is an incomplete designation, on any run.
	ideaDir = writeDesignationPrompt(t, root, "demo2", designationPrompt("demo2", participants, "implementer:\n"))
	ops = designationOps(t, root, "demo2", ideaDir, participants, participants, &bytes.Buffer{})
	if ops.roleErr == "" || !strings.Contains(ops.roleErr, "incomplete designation") {
		t.Fatalf("present-empty must gate, got roleErr=%q", ops.roleErr)
	}
}

// T-9 / AC-13: a designation defect in idea A does not gate a run of idea B.
func TestDesignationGateIsScopedToTheIdea(t *testing.T) {
	setCentralDefaults(t, "")
	root := designationRoot(t)
	participants := []string{"aa-first", "zz-impl"}
	defective := writeDesignationPrompt(t, root, "idea-a", designationPrompt("idea-a", participants, "implementer: zz-ghost\n"))
	clean := writeDesignationPrompt(t, root, "idea-b", designationPrompt("idea-b", participants, ""))
	opsA := designationOps(t, root, "idea-a", defective, participants, participants, &bytes.Buffer{})
	if opsA.roleErr == "" {
		t.Fatal("idea A's defective designation must gate idea A")
	}
	opsB := designationOps(t, root, "idea-b", clean, participants, participants, &bytes.Buffer{})
	if opsB.roleErr != "" || opsB.dispatchErr != "" {
		t.Fatalf("idea A's defect must not gate idea B: roleErr=%q dispatchErr=%q", opsB.roleErr, opsB.dispatchErr)
	}
}

// T-10 / AC-14: a two-participant designated idea warns at kickoff and proceeds.
func TestTwoParticipantDesignatedKickoffWarnsNotBlocks(t *testing.T) {
	setCentralDefaults(t, "")
	root := designationRoot(t)
	participants := []string{"zz-impl", "aa-rev"}
	ideaDir := writeDesignationPrompt(t, root, "demo", designationPrompt("demo", participants, "implementer: zz-impl\n"))
	var out bytes.Buffer
	ops := designationOps(t, root, "demo", ideaDir, participants, participants, &out)
	if ops.roleErr != "" || ops.dispatchErr != "" {
		t.Fatalf("a two-participant designated idea must not block: %q / %q", ops.roleErr, ops.dispatchErr)
	}
	if !strings.Contains(out.String(), "single non-implementer reviewer") {
		t.Fatalf("the kickoff warning must fire, got %q", out.String())
	}

	// The same kickoff warning fires under a tier-3 designation.
	setCentralDefaults(t, "default_implementer = \"zz-impl\"\n")
	root = designationRoot(t)
	ideaDir = writeDesignationPrompt(t, root, "demo", designationPrompt("demo", participants, ""))
	out.Reset()
	ops = designationOps(t, root, "demo", ideaDir, participants, participants, &out)
	if !strings.Contains(out.String(), "single non-implementer reviewer") {
		t.Fatalf("the tier-3 kickoff warning must fire, got %q", out.String())
	}
}

// R37: under a designation the unchanged model-diversity check also runs at kickoff —
// the warn path prints early, and the required path surfaces its error early.
func TestKickoffModelDiversityUnderDesignation(t *testing.T) {
	setCentralDefaults(t, "")
	root := designationRoot(t)
	participants := []string{"zz-impl", "aa-rev"}
	ideaDir := writeDesignationPrompt(t, root, "demo", designationPrompt("demo", participants, "implementer: zz-impl\n"))
	var out bytes.Buffer
	sameModel := []agents.Discovery{
		{Spec: agents.Spec{ID: "zz-impl", Model: "same-model"}, Found: true},
		{Spec: agents.Spec{ID: "aa-rev", Model: "same-model"}, Found: true},
	}
	newDriverImplOps(runner.Options{Root: root, Idea: protocol.IdeaStatus{Slug: "demo", Path: ideaDir, Participants: participants}, Agents: sameModel},
		root, "demo", ideaDir, participants, &out)
	if !strings.Contains(out.String(), "model-diversity") {
		t.Fatalf("the kickoff diversity check must warn early under a designation, got %q", out.String())
	}

	// The required variant surfaces the escalation text at kickoff.
	root = designationRoot(t)
	ideaDir = writeDesignationPrompt(t, root, "demo", designationPrompt("demo", participants, "implementer: zz-impl\nrequire_model_diversity: true\n"))
	out.Reset()
	newDriverImplOps(runner.Options{Root: root, Idea: protocol.IdeaStatus{Slug: "demo", Path: ideaDir, Participants: participants}, Agents: sameModel},
		root, "demo", ideaDir, participants, &out)
	if !strings.Contains(out.String(), "kickoff model-diversity check") || !strings.Contains(out.String(), "require_model_diversity") {
		t.Fatalf("the required variant must surface at kickoff, got %q", out.String())
	}

	// Without a designation nothing is printed at construction (inert preference).
	root = designationRoot(t)
	ideaDir = writeDesignationPrompt(t, root, "demo", designationPrompt("demo", participants, ""))
	out.Reset()
	newDriverImplOps(runner.Options{Root: root, Idea: protocol.IdeaStatus{Slug: "demo", Path: ideaDir, Participants: participants}, Agents: sameModel},
		root, "demo", ideaDir, participants, &out)
	if out.String() != "" {
		t.Fatalf("undesignated construction must print nothing, got %q", out.String())
	}
}

// R29: a live tier-2/tier-3 change against a recorded designation dispatch escalates
// rather than reassigning silently; an unchanged designation re-dispatches cleanly.
func TestReentryComparisonEscalatesOnDesignationChange(t *testing.T) {
	setCentralDefaults(t, "")
	participants := []string{"aa-first", "zz-impl", "zz-other"}
	record := func(t *testing.T, st store.Store, implementer string) {
		t.Helper()
		if err := st.Append(store.Event{Type: "agent.implementer_resolved",
			Data: map[string]any{"idea": "demo", "implementer": implementer, "source": "designation"}}); err != nil {
			t.Fatal(err)
		}
	}
	newOps := func(t *testing.T, extra string) (driverImplOps, store.Store) {
		root := designationRoot(t)
		ideaDir := writeDesignationPrompt(t, root, "demo", designationPrompt("demo", participants, extra))
		st := store.New(filepath.Join(root, protocol.DeckDir, "runs", "designation-test"))
		var out bytes.Buffer
		ops := designationOps(t, root, "demo", ideaDir, participants, participants, &out)
		ops.base.Store = st
		return ops, st
	}

	// Recorded dispatch of zz-impl; the designation now names zz-other → escalate.
	ops, st := newOps(t, "implementer: zz-other\n")
	record(t, st, "zz-impl")
	if err := ops.Implement(context.Background()); err == nil || !strings.Contains(err.Error(), "designation changed after a recorded dispatch") {
		t.Fatalf("a changed designation against a recorded dispatch must escalate, got %v", err)
	}

	// Unchanged designation re-dispatches and records again.
	ops, st = newOps(t, "implementer: zz-impl\n")
	record(t, st, "zz-impl")
	if err := ops.Implement(context.Background()); err != nil {
		t.Fatalf("an unchanged designation must re-dispatch cleanly, got %v", err)
	}

	// A recorded pin-sourced event never triggers the comparison.
	ops, st = newOps(t, "implementer: zz-impl\n")
	if err := st.Append(store.Event{Type: "agent.implementer_resolved",
		Data: map[string]any{"idea": "demo", "implementer": "zz-impl", "source": "pin"}}); err != nil {
		t.Fatal(err)
	}
	if err := ops.checkImplementerReentry(); err != nil {
		t.Fatalf("recorded source=pin must not trigger the comparison, got %v", err)
	}

	// An existing pin governs every resume: no comparison, no escalation (R30).
	ops, st = newOps(t, "implementer: zz-impl\n")
	record(t, st, "zz-impl")
	pin := "---\nidea: demo\nstatus: implemented\nimplementer: zz-impl\n---\n\n## Summary of work\npartial\n"
	if err := os.WriteFile(filepath.Join(ops.ideaDir, "IMPLEMENTATION.md"), []byte(pin), 0o644); err != nil {
		t.Fatal(err)
	}
	ops = designationOps(t, ops.root, "demo", ops.ideaDir, participants, participants, &bytes.Buffer{})
	ops.base.Store = st
	if ops.implSource != protocol.SourceImplementerPin {
		t.Fatalf("the pin must govern the resume, got source=%q", ops.implSource)
	}
	if err := ops.checkImplementerReentry(); err != nil {
		t.Fatalf("a pinned resolution never compares, got %v", err)
	}
}

// T-6 / AC-12 + R36: the drafter preference avoids the designee where an alternative
// exists and falls back to the designee when it is the only eligible drafter.
func TestDrafterSeparationPreference(t *testing.T) {
	discovered := []agents.Discovery{
		{Spec: agents.Spec{ID: "zz-impl", LaunchMode: agents.LaunchHeadless}, Found: true},
		{Spec: agents.Spec{ID: "aa-drafter", LaunchMode: agents.LaunchHeadless}, Found: true},
	}
	// Preference: an eligible non-designee drafts.
	agent, ok := firstEligibleHeadlessAgentPreferring(discovered, []string{"zz-impl", "aa-drafter"}, nil, protocol.FacilitatorRole{}, "zz-impl")
	if !ok || agent.ID != "aa-drafter" {
		t.Fatalf("must prefer the non-designee drafter, got %q ok=%v", agent.ID, ok)
	}
	// Degenerate fallback: the designee is the only eligible drafter — draft anyway.
	agent, ok = firstEligibleHeadlessAgentPreferring(discovered, []string{"zz-impl"}, nil, protocol.FacilitatorRole{}, "zz-impl")
	if !ok || agent.ID != "zz-impl" {
		t.Fatalf("must fall back to the designee rather than fail, got %q ok=%v", agent.ID, ok)
	}
	// No preference ("" = no designation): today's first-eligible behavior.
	agent, ok = firstEligibleHeadlessAgentPreferring(discovered, []string{"zz-impl", "aa-drafter"}, nil, protocol.FacilitatorRole{}, "")
	if !ok || agent.ID != "zz-impl" {
		t.Fatalf("empty preferNot must preserve today's order, got %q ok=%v", agent.ID, ok)
	}

	// The preference source: tier-2 wins; none/empty suppress tier 3; tier-3 applies
	// when eligible and is skipped when inapplicable.
	setCentralDefaults(t, "default_implementer = \"zz-global\"\n")
	root := designationRoot(t)
	participants := []string{"zz-impl", "zz-global"}
	ideaDir := writeDesignationPrompt(t, root, "demo", designationPrompt("demo", participants, "implementer: zz-impl\n"))
	if got := designatedImplementerPreference(root, ideaDir, participants, protocol.FacilitatorRole{}); got != "zz-impl" {
		t.Fatalf("tier-2 preference, got %q", got)
	}
	ideaDir = writeDesignationPrompt(t, root, "demo2", designationPrompt("demo2", participants, "implementer: none\n"))
	if got := designatedImplementerPreference(root, ideaDir, participants, protocol.FacilitatorRole{}); got != "" {
		t.Fatalf("none must suppress the tier-3 preference, got %q", got)
	}
	ideaDir = writeDesignationPrompt(t, root, "demo3", designationPrompt("demo3", participants, ""))
	if got := designatedImplementerPreference(root, ideaDir, participants, protocol.FacilitatorRole{}); got != "zz-global" {
		t.Fatalf("tier-3 preference, got %q", got)
	}
}
