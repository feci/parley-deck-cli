package app

import (
	"bytes"
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"time"

	"parley-deck-cli/internal/agents"
	"parley-deck-cli/internal/fsutil"
	"parley-deck-cli/internal/procctl"
	"parley-deck-cli/internal/protocol"
	"parley-deck-cli/internal/runner"
)

// defaultPingTimeout bounds each hosted-PONG probe. hermes inits ~40s; 90s gives
// headroom (consensus §"Roster ping").
const defaultPingTimeout = 90 * time.Second

// pongSentinel is the exact token a probed agent must echo to count as available.
const pongSentinel = "PONG"

// pongPrompt is the probe sent through the agent's real configured invocation.
const pongPrompt = "Reply with exactly the single token: PONG"

// preflightOptions configures a readiness check shared by the standalone command
// and the `parley run` pre-check.
type preflightOptions struct {
	Root        string
	JSON        bool
	Yes         bool
	PingTimeout time.Duration
	NoPing      bool
}

// gateKind classifies a pending gate so callers can name it and print the matching
// confirm command.
type gateKind string

const (
	gateBreakingFreshness gateKind = "breaking-freshness"
	gateUnknownRole       gateKind = "unknown-role"
	gateExcludeAgent      gateKind = "exclude-agent"
)

// gate is a pending readiness gate that requires explicit user confirmation. A
// gate never carries a low-risk default, so the auto-answerer never resolves it.
type gate struct {
	Kind    gateKind `json:"kind"`
	Detail  string   `json:"detail"`
	Confirm string   `json:"confirm"`
}

// freshness summarizes the protocol freshness classification.
type freshness struct {
	Role        string `json:"role"`
	DeckVersion string `json:"deckVersion"`
	LiveSha      string `json:"protocolSha256"`
	PackagedSha  string `json:"packagedProtocolSha256"`
	SkillVersion string `json:"skillVersion,omitempty"`
	// Classification: "source-advisory" | "in-sync" | "additive" | "breaking" |
	// "unknown-role" | "no-version-json".
	Classification string `json:"classification"`
	Summary        string `json:"summary"`
	// Synced is set when an additive consumer sync was written.
	Synced bool `json:"synced,omitempty"`
}

// rosterEntry is one row of the readiness roster table.
type rosterEntry struct {
	RosterID  string `json:"rosterId"`
	Runtime   string `json:"runtime"`
	Version   string `json:"version"`
	Available bool   `json:"available"`
	Reason    string `json:"reason,omitempty"`
}

// preflightReport is the full readiness result.
type preflightReport struct {
	Freshness freshness     `json:"freshness"`
	Roster    []rosterEntry `json:"roster"`
	Pinged    bool          `json:"pinged"`
	Gates     []gate        `json:"gates"`
}

// versionMeta is the subset of meta/version.json preflight reads.
type versionMeta struct {
	ProtocolRole           string `json:"protocolRole"`
	ProtocolSha256         string `json:"protocolSha256"`
	PackagedProtocolSha256 string `json:"packagedProtocolSha256"`
	DeckVersion            string `json:"deckVersion"`
}

// probeFunc is the seam for the hosted-PONG probe so tests do not shell out.
// It returns (available, reason). reason is empty when available.
type probeFunc func(ctx context.Context, root string, agent agents.Discovery, timeout time.Duration) (bool, string)

// pingProbe is the production probe; tests swap it for a fake.
var pingProbe probeFunc = hostedPONG

func runPreflight(ctx context.Context, args []string, stdout, stderr io.Writer) int {
	fs := flag.NewFlagSet("preflight", flag.ContinueOnError)
	fs.SetOutput(stderr)
	root := fs.String("dir", ".", "workspace directory")
	jsonOut := fs.Bool("json", false, "print JSON output")
	yes := fs.Bool("yes", false, "confirm excluding unavailable agents (records the exclusion)")
	pingTimeout := fs.Duration("ping-timeout", defaultPingTimeout, "per-agent hosted-PONG probe timeout")
	noPing := fs.Bool("no-ping", false, "presence-only check (skip the hosted-PONG round-trip)")
	if err := fs.Parse(args); err != nil {
		return 2
	}
	if fs.NArg() != 0 {
		fmt.Fprintln(stderr, "usage: parley preflight [--dir DIR] [--json] [--yes] [--ping-timeout D] [--no-ping]")
		return 2
	}

	opts := preflightOptions{
		Root:        *root,
		JSON:        *jsonOut,
		Yes:         *yes,
		PingTimeout: *pingTimeout,
		NoPing:      *noPing,
	}

	discovered, err := discoverConfigured(ctx, opts.Root)
	if err != nil {
		fmt.Fprintf(stderr, "agent config failed: %v\n", err)
		return 1
	}
	report, code, err := preflight(ctx, opts, selectedParticipants(discovered), stdout, stderr)
	if err != nil {
		fmt.Fprintf(stderr, "preflight failed: %v\n", err)
		return 1
	}
	if opts.JSON {
		return printJSON(stdout, report, stderr)
	}
	printPreflightReport(stdout, report)
	return code
}

// selectedParticipants is the default participant set preflight evaluates: the
// installed agents (mirroring selectedParticipantIDs' default), minus the
// deprecated legacy gemini. Uninstalled optional catalog entries are
// not_selected — never an availability gate. A selected agent that is present
// but fails its hosted PONG IS a gate.
func selectedParticipants(discovered []agents.Discovery) []agents.Discovery {
	out := make([]agents.Discovery, 0, len(discovered))
	for _, agent := range discovered {
		if agent.Found && agent.ID != "gemini" {
			out = append(out, agent)
		}
	}
	return out
}

// attendedRun reports whether a `parley run` invocation is attended. A run is
// attended when it stops for round-01 (--no-auto) without a blanket --yes, or
// when stdin is a terminal. Unattended runs (--auto with no TTY) must hard-stop
// on a gate and never read stdin.
func attendedRun(auto, yes bool) bool {
	if !auto && !yes {
		return true
	}
	return isTerminal(os.Stdin)
}

// runTaskPreflight is the `parley run` pre-check. It reuses the already-discovered
// agents and runs the shared readiness check before any idea is created. Per the
// operator ruling, the default is a full hosted-PONG roster ping every idea
// (`--no-ping` falls back to a presence-only Tier-0 check for speed/CI; `--no-preflight`
// skips it entirely). It returns (exitCode, stop): stop=true means abort the run before
// any idea is created. On a genuine gate it prints the gate + confirm command. Attended
// and unattended both stop on a gate (we never auto-answer the new gates); the only
// difference is unattended must never block on stdin, which this path honors because it
// never reads stdin.
func runTaskPreflight(ctx context.Context, root string, discovered []agents.Discovery, attended, noPing bool, stdout, stderr io.Writer) (int, bool) {
	opts := preflightOptions{Root: root, NoPing: noPing}
	report, code, err := preflight(ctx, opts, selectedParticipants(discovered), stdout, stderr)
	if err != nil {
		fmt.Fprintf(stderr, "preflight failed: %v\n", err)
		return 1, true
	}
	if code == 0 {
		if report.Freshness.Synced {
			fmt.Fprintf(stdout, "preflight: %s\n", report.Freshness.Summary)
		}
		return 0, false
	}
	// A gate or a hard failure: stop before runcontrol.Create. No half-open idea.
	if code == 1 {
		fmt.Fprintln(stderr, "preflight: hard stop — excluding unavailable agents would leave fewer than 2 participants (§1 non-solo).")
	} else {
		w := stdout
		if !attended {
			w = stderr
			fmt.Fprintln(w, "preflight: pending gate in an unattended run — hard stop (gates are never auto-answered).")
		} else {
			fmt.Fprintln(w, "preflight: pending gate — confirm before this idea can start:")
		}
		for _, g := range report.Gates {
			fmt.Fprintf(w, "  [%s] %s\n", g.Kind, g.Detail)
			fmt.Fprintf(w, "    confirm: %s\n", g.Confirm)
		}
	}
	return code, true
}

// preflight is the shared readiness check used by `parley preflight` and
// `parley run`. It computes freshness (writing an additive consumer sync when
// safe), pings the roster (unless NoPing), and returns the report plus an exit
// code: 0 ready · 1 hard failure (<2 participants after exclusion) · 3 pending
// gate. It never reads stdin and never auto-answers a gate.
func preflight(ctx context.Context, opts preflightOptions, discovered []agents.Discovery, stdout, stderr io.Writer) (preflightReport, int, error) {
	report := preflightReport{}

	fr, freshGates, err := classifyAndSyncFreshness(ctx, opts)
	if err != nil {
		return report, 1, err
	}
	report.Freshness = fr
	report.Gates = append(report.Gates, freshGates...)

	report.Pinged = !opts.NoPing
	report.Roster = checkRoster(ctx, opts, discovered)

	// Availability exclusion is a gate: an unavailable agent the user must
	// confirm excluding. Count would-be participants to honor the §1 non-solo
	// hard-stop.
	available := 0
	for _, entry := range report.Roster {
		if entry.Available {
			available++
			continue
		}
		report.Gates = append(report.Gates, gate{
			Kind:    gateExcludeAgent,
			Detail:  fmt.Sprintf("%s unavailable (%s) — confirm excluding it from this idea", entry.RosterID, entry.Reason),
			Confirm: fmt.Sprintf("parley preflight --dir %s --yes", opts.Root),
		})
	}

	// Hard failure: excluding the unavailable agents would leave < 2 participants
	// (the §1 non-solo hard-stop). This outranks the pending gate.
	if len(report.Roster) > 0 && available < 2 {
		return report, 1, nil
	}

	if len(report.Gates) > 0 {
		return report, 3, nil
	}
	return report, 0, nil
}

// classifyAndSyncFreshness reads meta/version.json + the installed skill status,
// classifies drift, and performs the one safe write (additive consumer sync).
// It returns any freshness gate (breaking bump / unknown role).
func classifyAndSyncFreshness(ctx context.Context, opts preflightOptions) (freshness, []gate, error) {
	root := strings.TrimSpace(opts.Root)
	if root == "" {
		root = "."
	}

	skillStatus, _ := parleyDeckSkillStatus(ctx, root)
	fr := freshness{SkillVersion: skillStatusVersion(skillStatus)}

	meta, err := readVersionMeta(root)
	if err != nil {
		if os.IsNotExist(err) {
			fr.Classification = "no-version-json"
			fr.Summary = "no meta/version.json — freshness check skipped (advisory)"
			return fr, nil, nil
		}
		return fr, nil, err
	}
	fr.Role = meta.ProtocolRole
	fr.DeckVersion = meta.DeckVersion
	fr.LiveSha = meta.ProtocolSha256
	fr.PackagedSha = meta.PackagedProtocolSha256

	switch strings.ToLower(strings.TrimSpace(meta.ProtocolRole)) {
	case "source":
		fr.Classification = "source-advisory"
		fr.Summary = "source — ahead of / independent of packaged skill; advisory only, never writes COOPERATION.md"
		return fr, nil, nil
	case "consumer":
		// fall through to drift classification below
	default:
		fr.Classification = "unknown-role"
		fr.Summary = "protocolRole absent/unknown — will not auto-write; confirm role + backfill"
		return fr, []gate{{
			Kind:    gateUnknownRole,
			Detail:  "meta/version.json has no protocolRole; confirm role (source|consumer) and backfill the field",
			Confirm: fmt.Sprintf("parley preflight --dir %s --yes", opts.Root),
		}}, nil
	}

	// Consumer path.
	if meta.ProtocolSha256 == meta.PackagedProtocolSha256 {
		fr.Classification = "in-sync"
		fr.Summary = "consumer — protocol matches packaged skill (in sync)"
		return fr, nil, nil
	}

	bump := classifyBump(meta.DeckVersion, skillStatusPackagedDeckVersion(skillStatus))
	if bump == bumpMajor {
		fr.Classification = "breaking"
		fr.Summary = fmt.Sprintf("consumer — packaged protocol differs and the deck bump is breaking (major); pending confirmation")
		return fr, []gate{{
			Kind:    gateBreakingFreshness,
			Detail:  "packaged protocol is a breaking (major) bump; review and confirm before adopting",
			Confirm: fmt.Sprintf("parley preflight --dir %s --yes", opts.Root),
		}}, nil
	}

	// Additive (minor/patch): auto-sync, preserving project zones, unless no-op.
	fr.Classification = "additive"
	synced, summary, err := syncConsumerProtocol(ctx, root, meta, skillStatus)
	if err != nil {
		return fr, nil, err
	}
	fr.Synced = synced
	fr.Summary = summary
	return fr, nil, nil
}

type bumpKind int

const (
	bumpUnknown bumpKind = iota
	bumpPatchMinor
	bumpMajor
)

// classifyBump compares the project deckVersion against the packaged deckVersion
// by semver. A bump to a higher major is breaking; minor/patch is additive.
// Unparseable versions are treated conservatively as additive (the safer
// auto-sync path is still zone-preserving, recorded, and git-reversible).
func classifyBump(projectVersion, packagedVersion string) bumpKind {
	pj, ok1 := parseMajor(projectVersion)
	pk, ok2 := parseMajor(packagedVersion)
	if !ok1 || !ok2 {
		return bumpPatchMinor
	}
	if pk > pj {
		return bumpMajor
	}
	return bumpPatchMinor
}

func parseMajor(version string) (int, bool) {
	v := strings.TrimSpace(strings.TrimPrefix(strings.TrimSpace(version), "v"))
	if v == "" {
		return 0, false
	}
	head, _, _ := strings.Cut(v, ".")
	n, err := strconv.Atoi(head)
	if err != nil {
		return 0, false
	}
	return n, true
}

// syncConsumerProtocol rewrites the project COOPERATION.md from the packaged
// protocol body, preserving project-specific zones, then records the sync. It is
// a no-op (returns synced=false) when the merge would not change the file. NOTE:
// this never runs in a source repo (gated by protocolRole earlier).
func syncConsumerProtocol(ctx context.Context, root string, meta versionMeta, skillStatus map[string]any) (bool, string, error) {
	packagedBody, ok := skillStatusPackagedProtocol(skillStatus)
	if !ok || strings.TrimSpace(packagedBody) == "" {
		// No packaged body available from the skill status: cannot merge. Report
		// without writing (fail-safe; the source repo never reaches here anyway).
		return false, "consumer — additive bump detected, but no packaged protocol body available to merge; no write performed", nil
	}
	deckPath := filepath.Join(root, protocol.DeckDir, "COOPERATION.md")
	projectBody, err := os.ReadFile(deckPath)
	if err != nil {
		return false, "", err
	}

	merged, err := mergePreservingZones(string(projectBody), packagedBody)
	if err != nil {
		return false, "", err
	}
	merged = refreshProtocolSyncedLine(merged, meta.DeckVersion, skillStatusVersion(skillStatus))
	if merged == string(projectBody) {
		return false, "consumer — additive bump, merge is a no-op (already current); no write performed", nil
	}

	if err := fsutil.WriteFileAtomic(deckPath, []byte(merged), 0o644); err != nil {
		return false, "", err
	}

	oldSha := sha256Hex(string(projectBody))
	newSha := sha256Hex(merged)
	if err := writeSyncRecord(root, meta, oldSha, newSha, projectBody, []byte(merged)); err != nil {
		return false, "", err
	}
	return true, fmt.Sprintf("consumer — additive protocol sync applied (preserved project zones); recorded under %s/meta/", protocol.DeckDir), nil
}

// mergePreservingZones produces the new COOPERATION.md body: project-specific
// zones (everything from the start of the file through the end of `## 2. Active
// agents (roster)`) are kept from the PROJECT file; everything from `## 3.`
// onward is taken from the PACKAGED file. Pure function — unit-tested in
// isolation.
func mergePreservingZones(projectBody, packagedBody string) (string, error) {
	projectHead, err := zonesThroughRoster(projectBody)
	if err != nil {
		return "", fmt.Errorf("project COOPERATION.md: %w", err)
	}
	packagedTail, err := sectionThreeOnward(packagedBody)
	if err != nil {
		return "", fmt.Errorf("packaged protocol: %w", err)
	}
	return projectHead + packagedTail, nil
}

const sectionThreeAnchor = "## 3."

// zonesThroughRoster returns everything from the start of body up to (but not
// including) the `## 3.` section line — i.e. the header + §0 + §1 + §2.
func zonesThroughRoster(body string) (string, error) {
	idx := sectionThreeIndex(body)
	if idx < 0 {
		return "", fmt.Errorf("no %q section line found", sectionThreeAnchor)
	}
	return body[:idx], nil
}

// sectionThreeOnward returns everything from the `## 3.` section line onward.
func sectionThreeOnward(body string) (string, error) {
	idx := sectionThreeIndex(body)
	if idx < 0 {
		return "", fmt.Errorf("no %q section line found", sectionThreeAnchor)
	}
	return body[idx:], nil
}

// sectionThreeIndex returns the byte offset of the start of the `## 3.` section
// line, or -1 if not present. The anchor must start a line.
func sectionThreeIndex(body string) int {
	for offset := 0; offset < len(body); {
		lineEnd := strings.IndexByte(body[offset:], '\n')
		var line string
		if lineEnd < 0 {
			line = body[offset:]
		} else {
			line = body[offset : offset+lineEnd]
		}
		if strings.HasPrefix(line, sectionThreeAnchor) {
			return offset
		}
		if lineEnd < 0 {
			break
		}
		offset += lineEnd + 1
	}
	return -1
}

// refreshProtocolSyncedLine updates (or, if absent, inserts after the title) the
// `**Protocol synced:**` provenance line in the header zone.
func refreshProtocolSyncedLine(body, deckVersion, skillVersion string) string {
	stamp := time.Now().Format("2006-01-02")
	skill := strings.TrimSpace(skillVersion)
	if skill == "" {
		skill = deckVersion
	}
	syncedLine := fmt.Sprintf("**Protocol synced:** %s — parley-deck-skill %s (preflight)", stamp, skill)

	lines := strings.Split(body, "\n")
	for i, line := range lines {
		if strings.HasPrefix(line, "**Protocol synced:**") {
			lines[i] = syncedLine
			return strings.Join(lines, "\n")
		}
	}
	// Not present: insert after the first non-empty title line.
	for i, line := range lines {
		if strings.TrimSpace(line) != "" {
			out := append([]string{}, lines[:i+1]...)
			out = append(out, syncedLine)
			out = append(out, lines[i+1:]...)
			return strings.Join(out, "\n")
		}
	}
	return body
}

// writeSyncRecord drops a maintenance-sync record under meta/.
func writeSyncRecord(root string, meta versionMeta, oldSha, newSha string, oldBody, newBody []byte) error {
	stamp := time.Now().UTC().Format("2006-01-02T150405Z")
	recordPath := filepath.Join(root, protocol.DeckDir, "meta", "protocol-sync_"+stamp+".md")
	record := fmt.Sprintf(`# Protocol sync — %s

- type: additive (consumer auto-sync via preflight §9.0)
- deckVersion: %s
- old protocolSha256: %s
- new protocolSha256: %s
- preserved zones: header + ## 0 (transport) + ## 1 (scope) + ## 2 (roster) kept from project; ## 3. onward from packaged protocol
- diff summary: %s
`, stamp, meta.DeckVersion, oldSha, newSha, diffSummary(oldBody, newBody))
	if err := fsutil.MkdirAllResilient(filepath.Dir(recordPath), 0o755); err != nil {
		return err
	}
	return fsutil.WriteFileAtomic(recordPath, []byte(record), 0o644)
}

func diffSummary(oldBody, newBody []byte) string {
	oldLines := bytes.Count(oldBody, []byte("\n"))
	newLines := bytes.Count(newBody, []byte("\n"))
	return fmt.Sprintf("%d -> %d lines (%+d)", oldLines, newLines, newLines-oldLines)
}

// checkRoster builds the roster table. With NoPing it is a presence-only (Tier-0)
// check; otherwise each present agent gets a bounded hosted-PONG probe, all run
// concurrently behind a single global deadline.
func checkRoster(ctx context.Context, opts preflightOptions, discovered []agents.Discovery) []rosterEntry {
	entries := make([]rosterEntry, len(discovered))
	timeout := opts.PingTimeout
	if timeout <= 0 {
		timeout = defaultPingTimeout
	}

	// Global deadline bounds all concurrent probes together.
	probeCtx := ctx
	if !opts.NoPing {
		var cancel context.CancelFunc
		probeCtx, cancel = context.WithTimeout(ctx, timeout)
		defer cancel()
	}

	var wg sync.WaitGroup
	for i := range discovered {
		agent := discovered[i]
		entries[i] = rosterEntry{
			RosterID: agent.ID,
			Runtime:  runtimeCLI(agent),
			Version:  agent.Version,
		}
		// Presence pre-check (command -v): a missing CLI is never PONG-probed.
		if !agent.Found {
			entries[i].Available = false
			entries[i].Reason = "unavailable:missing"
			continue
		}
		if opts.NoPing {
			// Tier-0 fallback: presence is the availability signal.
			entries[i].Available = true
			continue
		}
		wg.Add(1)
		go func(idx int, a agents.Discovery) {
			defer wg.Done()
			ok, reason := pingProbe(probeCtx, opts.Root, a, timeout)
			entries[idx].Available = ok
			if !ok {
				entries[idx].Reason = reason
			}
		}(i, agent)
	}
	wg.Wait()
	return entries
}

// hostedPONG runs the agent's real configured invocation with the PONG prompt,
// bounded by timeout, and kills the process group on timeout so no children
// leak. Available = exits in time AND stdout contains the sentinel.
func hostedPONG(ctx context.Context, root string, agent agents.Discovery, timeout time.Duration) (bool, string) {
	probeCtx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	cmd, cleanup, err := runner.CommandFor(probeCtx, root, agent, pongPrompt)
	if cleanup != nil {
		defer cleanup()
	}
	if err != nil {
		return false, "unavailable:command-build-error"
	}
	cmd.Dir = root
	// Spawn into its own process group so a timeout kill reaps the whole tree.
	procctl.SetNewProcessGroup(cmd)
	var out bytes.Buffer
	var errOut bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errOut
	if agent.PromptMode == agents.PromptStdin {
		cmd.Stdin = strings.NewReader(pongPrompt)
	}

	if err := cmd.Start(); err != nil {
		return false, "unavailable:start-error"
	}
	sp := procctl.Capture(cmd, "")
	waitErr := make(chan error, 1)
	go func() { waitErr <- cmd.Wait() }()

	select {
	case <-probeCtx.Done():
		_ = procctl.KillGroup(sp)
		<-waitErr // reap
		return false, "unavailable:timeout"
	case err := <-waitErr:
		if err != nil {
			return false, "unavailable:exit-error"
		}
		if strings.Contains(out.String(), pongSentinel) {
			return true, ""
		}
		return false, "unavailable:no-pong"
	}
}

func runtimeCLI(agent agents.Discovery) string {
	if len(agent.Commands) > 0 {
		return agent.Commands[0]
	}
	return ""
}

func readVersionMeta(root string) (versionMeta, error) {
	path := filepath.Join(root, protocol.DeckDir, "meta", "version.json")
	data, err := os.ReadFile(path)
	if err != nil {
		return versionMeta{}, err
	}
	var meta versionMeta
	if err := json.Unmarshal(data, &meta); err != nil {
		return versionMeta{}, fmt.Errorf("meta/version.json is not valid JSON: %w", err)
	}
	return meta, nil
}

// skillStatusVersion extracts the installed skill version from the status payload.
func skillStatusVersion(status map[string]any) string {
	installer := mapValue(status, "installer")
	return stringValue(installer, "version")
}

// skillStatusPackagedDeckVersion extracts the packaged deck version the installed
// skill carries, if the status payload exposes it. Falls back to "".
func skillStatusPackagedDeckVersion(status map[string]any) string {
	if v := stringValue(mapValue(status, "packaged"), "deckVersion"); v != "" {
		return v
	}
	if v := stringValue(mapValue(status, "protocol"), "deckVersion"); v != "" {
		return v
	}
	return skillStatusVersion(status)
}

// skillStatusPackagedProtocol extracts the packaged protocol body from the status
// payload if present. This source repo never reaches the consumer-sync path, so a
// missing body degrades to a no-write report rather than an error.
func skillStatusPackagedProtocol(status map[string]any) (string, bool) {
	for _, key := range []string{"packagedProtocol", "protocolBody"} {
		if v := stringValue(status, key); v != "" {
			return v, true
		}
	}
	if v := stringValue(mapValue(status, "packaged"), "protocol"); v != "" {
		return v, true
	}
	return "", false
}

func printPreflightReport(stdout io.Writer, report preflightReport) {
	fmt.Fprintln(stdout, "Parley Deck preflight readiness check")
	fr := report.Freshness
	fmt.Fprintf(stdout, "Freshness: %s\n", fr.Summary)
	if fr.Role != "" || fr.DeckVersion != "" {
		fmt.Fprintf(stdout, "  role=%s deckVersion=%s classification=%s\n", valueOr(fr.Role, "(none)"), valueOr(fr.DeckVersion, "(none)"), fr.Classification)
	}

	fmt.Fprintln(stdout, "Roster:")
	fmt.Fprintln(stdout, "  ROSTER-ID        RUNTIME      VERSION                 AVAILABLE  REASON")
	for _, entry := range report.Roster {
		fmt.Fprintf(stdout, "  %-16s %-12s %-23s %-10s %s\n",
			entry.RosterID,
			valueOr(entry.Runtime, "-"),
			truncatePreflight(valueOr(entry.Version, "-"), 23),
			yesNo(entry.Available),
			entry.Reason,
		)
	}
	if !report.Pinged {
		fmt.Fprintln(stdout, "  (presence-only; hosted PONG skipped via --no-ping)")
	}

	if len(report.Gates) == 0 {
		fmt.Fprintln(stdout, "Ready: no pending gates.")
		return
	}
	fmt.Fprintln(stdout, "Pending gates (require user confirmation):")
	for _, g := range report.Gates {
		fmt.Fprintf(stdout, "  [%s] %s\n", g.Kind, g.Detail)
		fmt.Fprintf(stdout, "    confirm: %s\n", g.Confirm)
	}
}

func truncatePreflight(value string, width int) string {
	if len(value) <= width {
		return value
	}
	if width <= 1 {
		return value[:width]
	}
	return value[:width-1] + "…"
}
