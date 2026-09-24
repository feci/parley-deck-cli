package app

// `parley wait` (lean-organizer B): one blocking read that returns when the awaited
// phase boundary is reached, on timeout, or loudly on degradation — an alternative to
// polling, never a phase actor (§14 human-brake boundary: it only observes; it never
// advances, writes, or edits anything).
//
// Exit codes (ratified 0/3/4/1 map with the missing≠invalid split):
//
//	0  boundary reached
//	3  timeout — partial digest printed, outstanding agents named
//	4  a PRESENT artifact fails the shared validator (validator error verbatim), or a
//	   NEW blocking escalation / driver.error event ARRIVES after the wait started —
//	   immediate (fix-up F2, FINAL B.3's "a **new** unanswered `to-user` escalation …
//	   arrives"): an escalation qualifies only when its frontmatter `idea:` matches the
//	   awaited slug, `blocking:` is not `no`, and `status:` is not answered/resolved;
//	   a `driver.error` counts only when appended to the event log after wait start.
//	   Pre-existing qualifying notes and historical errors are reported as digest
//	   annotations, never as exit 4.
//	1  usage / IO error
//
// Missing ≠ invalid: a not-yet-filed artifact keeps waiting; an invalid one exits
// loudly. The timeout default is configuration-first (per-user
// `[defaults.timeouts] wait_ms`, seeded 25 m — a heuristic bounded below the
// owner-requested 30-minute guidance, never asserted as a provider fact) and a hard
// ceiling rejects --timeout values above the active track's §4.0 agent timeout.

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"parley-deck-cli/internal/config"
	"parley-deck-cli/internal/driver"
	"parley-deck-cli/internal/protocol"
	"parley-deck-cli/internal/runstate"
	"parley-deck-cli/internal/store"
)

const (
	waitExitOK      = 0
	waitExitTimeout = 3
	waitExitInvalid = 4
	waitExitUsage   = 1

	waitDefaultMS = 25 * 60 * 1000 // 25 m; heuristic, NOT a provider fact
	waitMinPoll   = 10 * time.Second
)

// waitPollInterval is a test seam; the shipped interval never drops below 10 s.
var waitPollInterval = waitMinPoll

// trackTimeoutCeiling returns the §4.0 per-track agent timeout that bounds --timeout.
func trackTimeoutCeiling(track string) time.Duration {
	switch track {
	case "fast":
		return 5 * time.Minute
	case "standard":
		return 15 * time.Minute
	default: // deliberation and unknown fall to the most permissive protocol figure
		return 30 * time.Minute
	}
}

func runWait(args []string, stdout, stderr interface{ Write([]byte) (int, error) }) int {
	fs := flag.NewFlagSet("wait", flag.ContinueOnError)
	fs.SetOutput(stderr)
	dirFlag := fs.String("dir", ".", "workspace root")
	idea := fs.String("idea", "", "idea slug to wait on")
	forWhat := fs.String("for", "any", "round | consensus | review | implementation | any")
	timeout := fs.String("timeout", "", "bounded wait, e.g. 25m / 90s (default: [defaults.timeouts] wait_ms, ceiling: the track's §4.0 agent timeout)")
	jsonOut := fs.Bool("json", false, "print the digest as JSON")
	if err := fs.Parse(args); err != nil {
		return waitExitUsage
	}
	out, ok := stdout.(interface{ Write([]byte) (int, error) })
	if !ok {
		return waitExitUsage
	}
	if strings.TrimSpace(*idea) == "" {
		fmt.Fprintln(stderr, "wait: --idea is required")
		return waitExitUsage
	}
	scope := strings.TrimSpace(*forWhat)
	switch scope {
	case "round", "consensus", "review", "implementation", "any":
	default:
		fmt.Fprintf(stderr, "wait: --for must be round|consensus|review|implementation|any, got %q\n", scope)
		return waitExitUsage
	}

	root := *dirFlag
	ws, err := protocol.ReadWorkspaceStatus(root)
	if err != nil {
		fmt.Fprintf(stderr, "wait: %v\n", err)
		return waitExitUsage
	}
	var ideaStatus protocol.IdeaStatus
	for _, i := range ws.Ideas {
		if i.Slug == *idea {
			ideaStatus = i
			break
		}
	}
	if ideaStatus.Slug == "" {
		fmt.Fprintf(stderr, "wait: idea %q not found\n", *idea)
		return waitExitUsage
	}
	ideaDir := ideaStatus.Path
	track := "standard"
	if t, present := driver.ReadTrack(ideaDir); present && t != "" {
		track = string(t)
	}

	defaults, _ := config.LoadDefaults(root)
	waitMS := defaults.WaitMS
	if waitMS <= 0 {
		waitMS = waitDefaultMS
	}
	budget := time.Duration(waitMS) * time.Millisecond
	if raw := strings.TrimSpace(*timeout); raw != "" {
		parsed, err := time.ParseDuration(raw)
		if err != nil || parsed <= 0 {
			fmt.Fprintf(stderr, "wait: invalid --timeout %q (use e.g. 25m or 90s)\n", raw)
			return waitExitUsage
		}
		if ceiling := trackTimeoutCeiling(track); parsed > ceiling {
			fmt.Fprintf(stderr, "wait: --timeout %s exceeds the active track's §4.0 agent timeout ceiling (%s for track %q)\n", raw, ceiling, track)
			return waitExitUsage
		}
		budget = parsed
	}
	if ceiling := trackTimeoutCeiling(track); budget > ceiling {
		// The configured default is bounded by the same per-track ceiling.
		budget = ceiling
	}

	// Fix-up F2: the wait-start snapshot. Escalations and driver errors count as
	// blocking only when they ARRIVE after this point; earlier ones are reported as
	// annotations on the digest (the digest itself stays tree-derived and
	// deterministic — the annotations ride the wait output, not PhaseDigest).
	waitStart := time.Now()
	events := waitEventStore(root, *idea)
	startEventCount := eventLogLength(events)
	annotations := preExistingAnnotations(root, *idea, events, startEventCount, waitStart)
	deadline := time.Now().Add(budget)
	poll := waitPollInterval
	if poll < waitMinPoll {
		poll = waitMinPoll
	}
	for {
		digest := driver.BuildPhaseDigest(root, *idea, ideaDir, ideaStatus.Participants)
		if reason, bad := invalidArtifact(digest, scope); bad {
			printWaitDigest(out, digest, annotations, *jsonOut)
			fmt.Fprintf(stderr, "wait: present-but-invalid artifact: %s\n", reason)
			return waitExitInvalid
		}
		if src := arrivedBlockingEscalation(root, *idea, waitStart); src != "" {
			printWaitDigest(out, digest, annotations, *jsonOut)
			fmt.Fprintf(stderr, "wait: blocking escalation (new unanswered to-user inbox note for this idea): %s\n", src)
			return waitExitInvalid
		}
		if drvErr := driverErrorEventSince(events, startEventCount); drvErr != "" {
			printWaitDigest(out, digest, annotations, *jsonOut)
			fmt.Fprintf(stderr, "wait: %s\n", drvErr)
			return waitExitInvalid
		}
		if reached, why := boundaryReached(digest, scope); reached {
			printWaitDigest(out, digest, annotations, *jsonOut)
			fmt.Fprintf(out, "wait: boundary reached (%s)\n", why)
			return waitExitOK
		}
		if time.Now().After(deadline) {
			printWaitDigest(out, digest, annotations, *jsonOut)
			fmt.Fprintf(out, "wait: timeout after %s; outstanding: %s\n", budget, outstandingAgents(digest, scope))
			return waitExitTimeout
		}
		time.Sleep(minDuration(poll, time.Until(deadline)))
	}
}

func minDuration(a, b time.Duration) time.Duration {
	if b > 0 && b < a {
		return b
	}
	return a
}

// waitJSON is the --json envelope: the digest plus the wait's F2 annotations
// (pre-existing escalations / historical errors — reported, never exit-4).
type waitJSON struct {
	Notes  []string           `json:"notes,omitempty"`
	Digest driver.PhaseDigest `json:"digest"`
}

func printWaitDigest(out interface{ Write([]byte) (int, error) }, digest driver.PhaseDigest, annotations []string, asJSON bool) {
	if asJSON {
		data, _ := json.MarshalIndent(waitJSON{Notes: annotations, Digest: digest}, "", "  ")
		fmt.Fprintln(out, string(data))
		return
	}
	fmt.Fprintf(out, "PhaseDigest — %s\n", digest.Idea)
	if digest.Round != nil {
		printRoundSection(out, digest.Round)
	}
	if digest.Review != nil {
		printRoundSection(out, digest.Review)
	}
	if digest.Consensus != nil {
		c := digest.Consensus
		fmt.Fprintf(out, "consensus: present=%v triage=%s missing=%v signed=%v reservations=%v blocks=%v\n",
			c.Present, c.Triage, c.Missing, c.Signed, c.Reservations, c.Blocks)
	}
	if digest.Implementation != nil {
		i := digest.Implementation
		fmt.Fprintf(out, "implementation: present=%v status=%s implementer=%s\n", i.Present, i.Status, i.Implementer)
	}
	fmt.Fprintf(out, "next: %s\n", digest.Next)
	for _, note := range annotations {
		fmt.Fprintf(out, "note: %s\n", note)
	}
}

func printRoundSection(out interface{ Write([]byte) (int, error) }, sec *driver.PhaseRoundSection) {
	fmt.Fprintf(out, "%s: %d/%d filed-and-valid\n", sec.Label, sec.Completed, sec.Total)
	for _, row := range sec.Rows {
		validity := "ok"
		if !row.Valid {
			validity = row.Validity
		}
		fmt.Fprintf(out, "  %-12s filed=%-5v bytes=%-7d owner=%-12s valid=%v validity=%q stance(block/counter/accept/escalate)=%d/%d/%d/%d unparsed=%v fell_back=%v\n",
			row.Agent, row.Filed, row.Bytes, row.Owner, row.Valid, validity,
			row.StanceFlags.Block, row.StanceFlags.Counter, row.StanceFlags.Accept, row.StanceFlags.Escalate,
			row.Unparsed, row.FellBack)
		fmt.Fprintf(out, "    path: %s\n", row.Path)
	}
}

// invalidArtifact reports a PRESENT-but-invalid artifact in the awaited scope
// (missing ≠ invalid — absence keeps waiting).
func invalidArtifact(d driver.PhaseDigest, scope string) (string, bool) {
	secs := map[string][]*driver.PhaseRoundSection{}
	if d.Round != nil {
		secs["round"] = append(secs["round"], d.Round)
		secs["any"] = append(secs["any"], d.Round)
	}
	if d.Review != nil {
		secs["review"] = append(secs["review"], d.Review)
		secs["any"] = append(secs["any"], d.Review)
	}
	if scope == "consensus" || scope == "any" {
		if d.Consensus != nil && len(d.Consensus.Errors) > 0 {
			return strings.Join(d.Consensus.Errors, "; "), true
		}
	}
	for _, sec := range secs[scope] {
		for _, row := range sec.Rows {
			if row.Filed && !row.Valid {
				return row.Agent + ": " + row.Validity, true
			}
		}
	}
	return "", false
}

// escalationNote is one `to-user` inbox note that QUALIFIES for blocking under
// FINAL B.3 as read at fix-up F2: frontmatter `idea:` matches the awaited slug,
// `blocking:` is not `no`, and `status:` is not answered/resolved. Non-qualifying
// notes (other ideas, non-blocking, answered) are simply not this wait's business.
type escalationNote struct {
	Name  string
	MTime time.Time
}

// qualifyingEscalationNotes scans the deck inbox for qualifying `to-user` notes
// (answered notes are moved to inbox/archived/ or deleted per §4; a stale `status:`
// field saying answered/resolved counts as answered too).
func qualifyingEscalationNotes(root, idea string) []escalationNote {
	inbox := filepath.Join(root, protocol.DeckDir, "inbox")
	entries, err := os.ReadDir(inbox)
	if err != nil {
		return nil
	}
	var notes []escalationNote
	for _, e := range entries {
		name := e.Name()
		if e.IsDir() || !strings.HasSuffix(name, ".md") {
			continue
		}
		if strings.Index(name, "-to-user_") <= 0 {
			continue
		}
		path := filepath.Join(inbox, name)
		meta, merr := protocol.ReadFrontmatter(path)
		if merr != nil {
			continue
		}
		if id := strings.Trim(strings.TrimSpace(meta["idea"]), `"'`); id != idea {
			continue
		}
		if strings.EqualFold(strings.TrimSpace(meta["blocking"]), "no") {
			continue
		}
		if status := strings.ToLower(strings.TrimSpace(meta["status"])); strings.HasPrefix(status, "answered") || strings.HasPrefix(status, "resolved") {
			continue
		}
		info, ierr := os.Stat(path)
		if ierr != nil {
			continue
		}
		notes = append(notes, escalationNote{Name: name, MTime: info.ModTime()})
	}
	return notes
}

// arrivedBlockingEscalation names the qualifying escalation that ARRIVED after the
// wait started (mtime is the arrival signal), or "".
func arrivedBlockingEscalation(root, idea string, waitStart time.Time) string {
	for _, n := range qualifyingEscalationNotes(root, idea) {
		if n.MTime.After(waitStart) {
			return n.Name
		}
	}
	return ""
}

// preExistingAnnotations reports, once at wait start, what the older any-note
// semantics would have blocked on: qualifying escalations that predate the wait,
// and historical driver errors. Reported in the digest output, never exit 4.
func preExistingAnnotations(root, idea string, events store.Store, startEventCount int, waitStart time.Time) []string {
	var notes []string
	for _, n := range qualifyingEscalationNotes(root, idea) {
		if !n.MTime.After(waitStart) {
			notes = append(notes, "pre-existing unanswered to-user escalation for this idea (arrived before this wait; reported, not blocking): "+n.Name)
		}
	}
	if h := firstDriverErrorBefore(events, startEventCount); h != "" {
		notes = append(notes, "historical "+h+" (before this wait; reported, not blocking)")
	}
	return notes
}

// waitEventStore resolves the idea's latest run for event watching; without a run it
// returns a disabled store and event-based exits simply never fire (tree state still
// drives every other exit).
func waitEventStore(root, idea string) store.Store {
	runsDir := filepath.Join(root, protocol.DeckDir, "runs")
	// ResolveRun takes the WORKSPACE root (it joins parley-deck/runs itself).
	summary, err := runstate.ResolveRun(root, idea)
	if err == nil && summary.RunID != "" {
		return store.New(filepath.Join(runsDir, summary.RunID))
	}
	return store.Store{}
}

// eventLogLength snapshots the event count at wait start; only entries appended
// after that index can block (fix-up F2 — FINAL B.3's "arrives").
func eventLogLength(events store.Store) int {
	if !events.Enabled() {
		return 0
	}
	evs, err := events.Load()
	if err != nil {
		return 0
	}
	return len(evs)
}

func driverErrorEventSince(events store.Store, startIndex int) string {
	if !events.Enabled() {
		return ""
	}
	evs, err := events.Load()
	if err != nil {
		return ""
	}
	if startIndex > len(evs) {
		startIndex = len(evs)
	}
	for _, ev := range evs[startIndex:] {
		switch ev.Type {
		case "driver.error", "run.error":
			if detail, _ := ev.Data["error"].(string); detail != "" {
				return "driver.error event: " + detail
			}
			return "driver.error event at " + ev.Time.Format(time.RFC3339)
		}
	}
	return ""
}

// firstDriverErrorBefore names the first historical driver error (annotation only).
func firstDriverErrorBefore(events store.Store, startIndex int) string {
	if !events.Enabled() {
		return ""
	}
	evs, err := events.Load()
	if err != nil || startIndex > len(evs) {
		return ""
	}
	for _, ev := range evs[:startIndex] {
		switch ev.Type {
		case "driver.error", "run.error":
			if detail, _ := ev.Data["error"].(string); detail != "" {
				return "driver.error event: " + detail
			}
			return "driver.error event at " + ev.Time.Format(time.RFC3339)
		}
	}
	return ""
}

// boundaryReached evaluates the awaited condition from tree state only.
func boundaryReached(d driver.PhaseDigest, scope string) (bool, string) {
	roundDone := d.Round != nil && d.Round.Completed == d.Round.Total && d.Round.Total > 0
	reviewDone := d.Review != nil && d.Review.Completed == d.Review.Total && d.Review.Total > 0
	consensusDone := d.Consensus != nil && d.Consensus.Triage == "ready"
	implDone := d.Implementation != nil && d.Implementation.ReadyForReview
	finalDone := d.Implementation != nil && d.Implementation.Status == "complete"
	switch scope {
	case "round":
		return roundDone, "round complete"
	case "consensus":
		return consensusDone, "consensus ready"
	case "review":
		return reviewDone, "review round complete"
	case "implementation":
		return implDone || finalDone, "implementation published"
	default: // any
		if finalDone {
			return true, "implementation complete"
		}
		if implDone {
			return true, "implementation published"
		}
		if consensusDone {
			return true, "consensus ready"
		}
		if reviewDone {
			return true, "review round complete"
		}
		if roundDone {
			return true, "round complete"
		}
		return false, ""
	}
}

// outstandingAgents names who/what is still missing at timeout.
func outstandingAgents(d driver.PhaseDigest, scope string) string {
	var missing []string
	if d.Round != nil && (scope == "round" || scope == "any") {
		for _, row := range d.Round.Rows {
			if !row.Filed {
				missing = append(missing, row.Agent+" (round artifact)")
			}
		}
	}
	if d.Review != nil && (scope == "review" || scope == "any") {
		for _, row := range d.Review.Rows {
			if !row.Filed {
				missing = append(missing, row.Agent+" (review artifact)")
			}
		}
	}
	if (scope == "consensus" || scope == "any") && d.Consensus != nil {
		missing = append(missing, d.Consensus.Missing...)
	}
	// Name the implementation's blocking condition instead of "none named" (fix-up
	// F3): the timeout line must never conceal WHY the boundary is unreachable.
	if scope == "implementation" || scope == "any" {
		switch {
		case d.Implementation == nil || !d.Implementation.Present:
			missing = append(missing, "IMPLEMENTATION.md not filed")
		case d.Implementation.Status == "complete", d.Implementation.ReadyForReview:
			// not blocking on this scope
		case d.Implementation.Status == "":
			missing = append(missing, "implementation status is empty — not a recognised ready state")
		default:
			missing = append(missing, "implementation status `"+d.Implementation.Status+"` is not a recognised ready state")
		}
	}
	if len(missing) == 0 {
		missing = append(missing, "awaited condition not yet reached — inspect the digest")
	}
	return strings.Join(missing, ", ")
}

// waitUsage contributes to printUsage.
func waitUsage() string {
	return "" +
		"  wait --idea <slug> --for round|consensus|review|implementation|any [--timeout 25m] [--json]\n" +
		"    Blocking read: exit 0 boundary reached; 3 timeout (partial digest, outstanding named);\n" +
		"    4 present-but-invalid artifact, or a NEW blocking escalation (idea-matching, unanswered,\n" +
		"    blocking) or driver.error arriving after wait start — pre-existing ones are digest\n" +
		"    notes; 1 usage/IO. Observes only.\n"
}
