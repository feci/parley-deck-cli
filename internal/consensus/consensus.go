package consensus

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"

	"parley-deck-cli/internal/protocol"
	"parley-deck-cli/internal/track"
)

const (
	StatusAccept       = "✅ ACCEPT"
	StatusReservations = "🟡 ACCEPT-WITH-RESERVATIONS"
	StatusBlock        = "❌ BLOCK"

	TriageReady     = "ready"
	TriageReserved  = "reserved"
	TriageBlocked   = "blocked"
	TriagePartial   = "partial"
	TriageMalformed = "malformed"
)

type DraftOptions struct {
	Review         bool
	Round          int
	By             string
	ReviewedCommit string
	Now            time.Time
}

type SignoffOptions struct {
	Review          bool
	Agent           string
	Status          string
	Notes           string
	CounterProposal string
	Now             time.Time
}

type FinalizeOptions struct {
	By  string
	Now time.Time
}

type ReopenOptions struct {
	Review bool
	Reason string
}

type Summary struct {
	Idea         string    `json:"idea"`
	Path         string    `json:"path"`
	Review       bool      `json:"review"`
	Triage       string    `json:"triage"`
	Participants []string  `json:"participants"`
	Signoffs     []Signoff `json:"signoffs"`
	Missing      []string  `json:"missing,omitempty"`
	Errors       []string  `json:"errors,omitempty"`
	// Scaffolded is true when finalize WROTE the FINAL.md template and deliberately left the idea
	// open. It is not an error and not a closure; the drafter fills the scaffold in and re-runs.
	Scaffolded bool `json:"scaffolded,omitempty"`
}

type Signoff struct {
	Agent           string `json:"agent"`
	Date            string `json:"date"`
	Status          string `json:"status"`
	Notes           string `json:"notes,omitempty"`
	CounterProposal string `json:"counter_proposal,omitempty"`
	Line            int    `json:"line"`
}

type document struct {
	Path     string
	Raw      string
	Signoffs []Signoff
}

type schema struct {
	review bool
	rel    string
}

var appendMu sync.Mutex

var signoffHeader = regexp.MustCompile(`^### Signoff:\s*([A-Za-z0-9._-]+)\s+(?:—|-)\s+(\d{4}-\d{2}-\d{2})\s*$`)

func Status(root, ideaSlug string, review bool) (Summary, error) {
	idea, err := findIdea(root, ideaSlug)
	if err != nil {
		return Summary{}, err
	}
	path := consensusPath(idea.Path, schemaFor(review))
	doc, err := parseDocument(path)
	if err != nil {
		return Summary{}, err
	}
	// A REVIEW consensus is signed by the agents who reviewed, not by every participant: §6
	// forbids the implementer reviewing its own work, so demanding its signoff left every
	// standard-track review consensus permanently `partial` no matter how many reviewers accepted
	// (audit finding codex-1/F3). The auto-driver calls this same validator, so the quorum
	// override did not bind at the close gate either.
	//
	// Same rule as the round gate (codex-1/F2), applied at close: exclude the resolved
	// implementer, and fail closed to the full list when it cannot be resolved.
	// Two different questions, and collapsing them broke nine review consensuses (review round 1,
	// @kimi-1 MAJOR):
	//   - WHO IS AWAITED?  not the implementer — §6 forbids it reviewing its own work (F2/F3).
	//   - WHO MAY SIGN?    the implementer certainly may, and does: its "fix-up cycle N applied
	//                      the agreed fixes" report is a signoff this deck writes constantly.
	// Passing the reduced list to validateDocument answered both with the first, so the
	// implementer's own signoff became "unknown participant", and malformed outranks every other
	// triage. Two in-flight ideas flipped to malformed.
	return validateDocumentAwaiting(idea.Slug, idea.Participants,
		reviewConsensusVoters(idea.Path, idea.Participants, review), review, doc), nil
}

// reviewConsensusVoters is who must SIGN a consensus — a different rule from who may AUTHOR a
// review round.
//
// Phase 6 authorship and Phase 7 approval are separate: §6 forbids the implementer reviewing its
// own work, while `deliberation`'s §4.0 row requires **all participants** to sign off. Reusing the
// round-author list for the signoff quorum dropped the implementer's REQUIRED signature on
// deliberation ideas, so a deliberation review consensus read `ready` while a mandatory signoff
// was missing (review round 1, @codex-1 MAJOR).
//
// fast/standard: the reviewers who reviewed. deliberation: everyone.
func reviewConsensusVoters(ideaDir string, participants []string, review bool) []string {
	if !review {
		return participants
	}
	if t, present, err := readIdeaTrack(ideaDir); err == nil && present && t == track.Deliberation {
		return participants
	}
	return expectedRoundParticipants(ideaDir, participants, review)
}

// readIdeaTrack reads §4.0's `track:` from an idea's 00-prompt.md.
func readIdeaTrack(ideaDir string) (track.Track, bool, error) {
	meta, err := protocol.ReadFrontmatter(filepath.Join(ideaDir, "00-prompt.md"))
	if err != nil {
		return track.Standard, false, err
	}
	raw, ok := meta["track"]
	if !ok {
		return track.Standard, false, nil
	}
	return track.NormalizeStrict(raw)
}

func Draft(root, ideaSlug string, opts DraftOptions) (Summary, error) {
	idea, err := findIdea(root, ideaSlug)
	if err != nil {
		return Summary{}, err
	}
	s := schemaFor(opts.Review)
	path := consensusPath(idea.Path, s)
	if _, err := os.Stat(path); err == nil {
		return Summary{}, fmt.Errorf("%s already exists", path)
	} else if !errors.Is(err, os.ErrNotExist) {
		return Summary{}, err
	}

	roundLabel, roundRel, err := selectRound(idea.Path, opts.Review, opts.Round)
	if err != nil {
		return Summary{}, err
	}
	expected := expectedRoundParticipants(idea.Path, idea.Participants, opts.Review)
	roundDir := filepath.Join(roundBaseDir(idea.Path, opts.Review), roundLabel)
	if missing := missingRoundArtifacts(roundDir, expected); len(missing) > 0 {
		return Summary{}, fmt.Errorf("%s is incomplete; missing %s", roundRel, strings.Join(missing, ", "))
	}
	if opts.Review {
		if err := validateReviewRound(roundDir, roundLabel, roundRel, idea.Slug, expected); err != nil {
			return Summary{}, err
		}
	}

	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return Summary{}, err
	}
	if err := writeFileAtomic(path, []byte(draftTemplate(idea, opts, roundLabel, roundRel)), 0o644); err != nil {
		return Summary{}, err
	}
	if !opts.Review {
		if err := updateIdeaStatus(idea.Path, "consensus"); err != nil {
			return Summary{}, err
		}
	}
	return Status(root, ideaSlug, opts.Review)
}

// validateReviewRound applies the review-artifact contract to every expected file before a review
// consensus may be drafted.
//
// The driver validated each artifact as its agent wrote it; `parley consensus draft --review` did
// not, so the manual path accepted review files with no `reviewed-commit` and no
// `## Refutation attempts`. That is this deck's recurring defect class — a printed rule binds only
// where enforcement lives — applied to the rule that makes a review attributable to a tree
// (codex-1, MAJOR, review round-02).
//
// The check binds what is being drafted now. Historical rounds whose consensus already exists are
// never revalidated: Draft refuses outright when consensus.md is present.
func validateReviewRound(roundDir, roundLabel, roundRel, ideaSlug string, expected []string) error {
	number, err := strconv.Atoi(strings.TrimPrefix(roundLabel, "round-"))
	if err != nil {
		return fmt.Errorf("cannot read a round number from %q: %v", roundLabel, err)
	}
	var problems []string
	for _, participant := range expected {
		path := filepath.Join(roundDir, participant+".md")
		if err := protocol.ValidateReviewArtifact(path, participant, ideaSlug, number); err != nil {
			problems = append(problems, err.Error())
		}
	}
	if len(problems) > 0 {
		return fmt.Errorf("%s cannot be drafted; %d review artifact(s) do not meet the review contract:\n  %s",
			roundRel, len(problems), strings.Join(problems, "\n  "))
	}
	return nil
}

func AppendSignoff(root, ideaSlug string, opts SignoffOptions) (Summary, error) {
	idea, err := findIdea(root, ideaSlug)
	if err != nil {
		return Summary{}, err
	}
	if !contains(idea.Participants, opts.Agent) {
		return Summary{}, fmt.Errorf("unknown participant %q", opts.Agent)
	}
	status, err := CanonicalStatus(opts.Status)
	if err != nil {
		return Summary{}, err
	}
	if (status == StatusReservations || status == StatusBlock) && strings.TrimSpace(opts.Notes) == "" {
		return Summary{}, fmt.Errorf("notes are required for %s", status)
	}
	if status == StatusBlock && strings.TrimSpace(opts.CounterProposal) == "" {
		return Summary{}, errors.New("counter-proposal is required for ❌ BLOCK")
	}

	appendMu.Lock()
	defer appendMu.Unlock()

	path := consensusPath(idea.Path, schemaFor(opts.Review))
	doc, err := parseDocument(path)
	if err != nil {
		return Summary{}, err
	}
	current := validateDocument(idea.Slug, idea.Participants, opts.Review, doc)
	if len(current.Errors) > 0 {
		return Summary{}, fmt.Errorf("cannot append to malformed consensus: %s", strings.Join(current.Errors, "; "))
	}
	for _, signoff := range current.Signoffs {
		if signoff.Agent == opts.Agent {
			return Summary{}, fmt.Errorf("participant %s already signed", opts.Agent)
		}
	}

	now := opts.Now
	if now.IsZero() {
		now = time.Now()
	}
	block := signoffBlock(opts.Agent, now.Format("2006-01-02"), status, opts.Notes, opts.CounterProposal)
	file, err := os.OpenFile(path, os.O_WRONLY|os.O_APPEND, 0o644)
	if err != nil {
		return Summary{}, err
	}
	defer file.Close()
	if !strings.HasSuffix(doc.Raw, "\n") {
		block = "\n" + block
	}
	if _, err := file.WriteString(block); err != nil {
		return Summary{}, err
	}
	return Status(root, ideaSlug, opts.Review)
}

func Finalize(root, ideaSlug string, opts FinalizeOptions) (string, Summary, error) {
	idea, err := findIdea(root, ideaSlug)
	if err != nil {
		return "", Summary{}, err
	}
	if idea.Status == "final" {
		return "", Summary{}, fmt.Errorf("idea %s is already final", idea.Slug)
	}
	// FINAL.md records who drafted it, and Phase 4 says the drafter is the initiator or an agreed
	// participant. `--by` was written into the artifact unchecked, so an idea could be closed in
	// the name of an agent that never took part (audit finding codex-1/F4).
	if by := strings.TrimSpace(opts.By); by != "" && len(idea.Participants) > 0 && !containsID(idea.Participants, by) {
		return "", Summary{}, fmt.Errorf("--by %q is not a participant of %s (participants: %s)", by, idea.Slug, strings.Join(idea.Participants, ", "))
	}
	summary, err := Status(root, ideaSlug, false)
	if err != nil {
		return "", Summary{}, err
	}
	switch summary.Triage {
	case TriageReady:
	case TriageReserved:
		doc, err := parseDocument(summary.Path)
		if err != nil {
			return "", Summary{}, err
		}
		if !hasSectionContent(doc.Raw, "## Open items deferred to implementation") {
			return "", Summary{}, errors.New("reserved consensus requires open items deferred to implementation before finalize")
		}
		// Any filler satisfied the section, so a reservation could be carried past finalize
		// without ever being written down (audit finding codex-1/F9). A reservation belongs to an
		// agent; the deferred items must name it, or nobody can tell which reservation was
		// deferred or whether it was addressed.
		if unlogged := unloggedReservations(doc); len(unlogged) > 0 {
			return "", Summary{}, fmt.Errorf("reserved consensus: %s reserved but %s not named in '## Open items deferred to implementation'",
				strings.Join(unlogged, ", "), plural(len(unlogged), "is", "are"))
		}
	default:
		return "", Summary{}, fmt.Errorf("cannot finalize consensus with triage=%s", summary.Triage)
	}

	// `finalize` used to write the FINAL.md scaffold and set `status: final` in the same breath,
	// so an idea was permanently closed around an empty outline while the command printed success
	// (audit finding codex-1/F5). FINAL.md is the single source of truth; closing an idea around a
	// template makes that statement false.
	//
	// It is now two truthful steps, and re-running it is how you take the second:
	//   1. no FINAL.md          → write the scaffold, leave the idea OPEN, say so
	//   2. FINAL.md written     → close the idea
	//   2'. FINAL.md still a scaffold → refuse, and name exactly what is missing
	path := filepath.Join(idea.Path, "FINAL.md")
	existing, statErr := os.ReadFile(path)
	switch {
	case statErr == nil:
		// The same gate the driver applies — slug and status included. Checking content only let
		// the manual path close an idea around an artifact declaring a different idea and a
		// non-final status (review round 1, @codex-1 MAJOR).
		if reason := protocol.ValidateFinal(string(existing), idea.Slug); reason != "" {
			return "", Summary{}, fmt.Errorf("%s cannot close this idea: %s", path, reason)
		}
	case errors.Is(statErr, os.ErrNotExist):
		now := opts.Now
		if now.IsZero() {
			now = time.Now()
		}
		if err := writeFileAtomic(path, []byte(finalTemplate(idea, opts.By, now)), 0o644); err != nil {
			return "", Summary{}, err
		}
		summary.Scaffolded = true
		return path, summary, nil
	default:
		return "", Summary{}, statErr
	}

	if err := updateIdeaStatus(idea.Path, "final"); err != nil {
		return "", Summary{}, err
	}
	updated, err := Status(root, ideaSlug, false)
	if err != nil {
		return path, Summary{}, err
	}
	return path, updated, nil
}

func Reopen(root, ideaSlug string, opts ReopenOptions) (string, error) {
	if strings.TrimSpace(opts.Reason) == "" {
		return "", errors.New("reason is required")
	}
	idea, err := findIdea(root, ideaSlug)
	if err != nil {
		return "", err
	}
	summary, err := Status(root, ideaSlug, opts.Review)
	if err != nil {
		return "", err
	}
	if summary.Triage != TriageBlocked {
		return "", fmt.Errorf("cannot reopen consensus with triage=%s", summary.Triage)
	}

	latestRound, _, err := selectRound(idea.Path, opts.Review, 0)
	if err != nil {
		latestRound = "round-01"
	}
	aborted, err := nextAbortedPath(roundBaseDir(idea.Path, opts.Review), latestRound)
	if err != nil {
		return "", err
	}
	if err := os.Rename(summary.Path, aborted); err != nil {
		return "", err
	}
	reason := fmt.Sprintf("\n\n## Reopen reason\n\n%s\n", strings.TrimSpace(opts.Reason))
	file, err := os.OpenFile(aborted, os.O_WRONLY|os.O_APPEND, 0o644)
	if err != nil {
		return "", err
	}
	_, writeErr := file.WriteString(reason)
	if err := closeWith(writeErr, file); err != nil {
		return "", err
	}
	if !opts.Review {
		if err := updateIdeaStatus(idea.Path, latestRound); err != nil {
			return "", err
		}
	}
	return aborted, nil
}

func CanonicalStatus(raw string) (string, error) {
	value := strings.TrimSpace(raw)
	if before, _, ok := strings.Cut(value, "("); ok {
		value = strings.TrimSpace(before)
	}
	switch strings.ToLower(value) {
	case "accept", strings.ToLower(StatusAccept):
		return StatusAccept, nil
	case "reserve", "reservations", "accept-with-reservations", strings.ToLower(StatusReservations):
		return StatusReservations, nil
	case "block", strings.ToLower(StatusBlock):
		return StatusBlock, nil
	default:
		return "", fmt.Errorf("unknown signoff status %q", raw)
	}
}

func schemaFor(review bool) schema {
	if review {
		return schema{review: true, rel: filepath.Join("review", "consensus.md")}
	}
	return schema{rel: "consensus.md"}
}

func consensusPath(ideaDir string, s schema) string {
	return filepath.Join(ideaDir, s.rel)
}

func findIdea(root, slug string) (protocol.IdeaStatus, error) {
	status, err := protocol.ReadWorkspaceStatus(root)
	if err != nil {
		return protocol.IdeaStatus{}, err
	}
	for _, idea := range status.Ideas {
		if idea.Slug == slug {
			return idea, nil
		}
	}
	return protocol.IdeaStatus{}, fmt.Errorf("idea %q not found", slug)
}

func parseDocument(path string) (document, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return document{}, err
	}
	lines := strings.Split(string(data), "\n")
	doc := document{Path: path, Raw: string(data)}
	// Signoffs count only from the `## Signoffs` heading onward.
	//
	// The parser used to scan every line with no idea which section it was in, so a
	// `### Signoff:` block quoted as an example under an earlier section counted as real
	// append-only approval and could carry the consensus gate on its own (audit finding
	// codex-1/F20).
	//
	// A stricter "must be inside the Signoffs section" rule was measured and rejected: 32 of this
	// deck's 405 signoffs legitimately sit under later headings such as
	// "Cycle 2 (review/round-02 → complete)". Position is the property that separates the two —
	// measured across every consensus.md and review/consensus.md in the deck, 405 of 405 real
	// signoffs appear at or after the heading and none before it.
	signoffsFrom := len(lines)
	for i, line := range lines {
		if rest, ok := strings.CutPrefix(strings.TrimSpace(line), "## "); ok {
			if strings.HasPrefix(strings.ToLower(strings.TrimSpace(rest)), "signoff") {
				signoffsFrom = i
				break
			}
		}
	}
	for i := 0; i < len(lines); i++ {
		if i < signoffsFrom {
			continue
		}
		match := signoffHeader.FindStringSubmatch(strings.TrimSpace(lines[i]))
		if len(match) != 3 {
			continue
		}
		signoff := Signoff{Agent: match[1], Date: match[2], Line: i + 1}
		for j := i + 1; j < len(lines); j++ {
			line := strings.TrimSpace(lines[j])
			if signoffHeader.MatchString(line) {
				break
			}
			if value, ok := strings.CutPrefix(line, "Status:"); ok {
				signoff.Status = strings.TrimSpace(value)
				continue
			}
			if value, ok := strings.CutPrefix(line, "Notes:"); ok {
				signoff.Notes = strings.TrimSpace(value)
				continue
			}
			if strings.HasPrefix(line, "Counter-proposal") {
				_, value, ok := strings.Cut(line, ":")
				if ok {
					signoff.CounterProposal = strings.TrimSpace(value)
				}
				continue
			}
		}
		doc.Signoffs = append(doc.Signoffs, signoff)
	}
	return doc, nil
}

// validateDocument keeps its signature for callers that await exactly the participants they list.
func validateDocument(ideaSlug string, participants []string, review bool, doc document) Summary {
	return validateDocumentAwaiting(ideaSlug, participants, participants, review, doc)
}

// validateDocumentAwaiting separates who may SIGN (known) from who is AWAITED (required).
// A signoff from a known participant is always valid; only `required` drives missing/triage.
func validateDocumentAwaiting(ideaSlug string, known, required []string, review bool, doc document) Summary {
	summary := Summary{
		Idea:         ideaSlug,
		Path:         doc.Path,
		Review:       review,
		Participants: append([]string(nil), required...),
		Signoffs:     append([]Signoff(nil), doc.Signoffs...),
	}
	// The consensus document declares which idea it belongs to, and nothing checked it — a
	// consensus.md whose frontmatter named a DIFFERENT idea was read as this idea's consensus
	// (audit finding codex-1/F21). Copying a consensus between ideas is exactly how a signoff
	// for one decision silently becomes approval of another.
	if declared := strings.Trim(strings.TrimSpace(frontmatterField(doc.Raw, "idea")), `"'`); declared != "" && declared != ideaSlug {
		summary.Errors = append(summary.Errors, fmt.Sprintf("frontmatter idea=%q but this is the consensus for %q", declared, ideaSlug))
	}
	signed := map[string]bool{}
	hasReservations := false
	hasBlock := false
	for _, signoff := range doc.Signoffs {
		switch {
		case !contains(known, signoff.Agent):
			summary.Errors = append(summary.Errors, fmt.Sprintf("line %d: unknown participant %s", signoff.Line, signoff.Agent))
		case signed[signoff.Agent]:
			summary.Errors = append(summary.Errors, fmt.Sprintf("line %d: duplicate signoff for %s", signoff.Line, signoff.Agent))
		default:
			signed[signoff.Agent] = true
		}
		status, err := CanonicalStatus(signoff.Status)
		if err != nil {
			summary.Errors = append(summary.Errors, fmt.Sprintf("line %d: %v", signoff.Line, err))
			continue
		}
		switch status {
		case StatusReservations:
			hasReservations = true
			if strings.TrimSpace(signoff.Notes) == "" {
				summary.Errors = append(summary.Errors, fmt.Sprintf("line %d: notes are required for %s", signoff.Line, StatusReservations))
			}
		case StatusBlock:
			hasBlock = true
			if strings.TrimSpace(signoff.Notes) == "" {
				summary.Errors = append(summary.Errors, fmt.Sprintf("line %d: notes are required for %s", signoff.Line, StatusBlock))
			}
			if strings.TrimSpace(signoff.CounterProposal) == "" {
				summary.Errors = append(summary.Errors, fmt.Sprintf("line %d: counter-proposal is required for %s", signoff.Line, StatusBlock))
			}
		}
	}
	for _, participant := range required {
		if !signed[participant] {
			summary.Missing = append(summary.Missing, participant)
		}
	}
	switch {
	case len(summary.Errors) > 0:
		summary.Triage = TriageMalformed
	case hasBlock:
		summary.Triage = TriageBlocked
	case len(summary.Missing) > 0:
		summary.Triage = TriagePartial
	case hasReservations:
		summary.Triage = TriageReserved
	default:
		summary.Triage = TriageReady
	}
	return summary
}

func selectRound(ideaDir string, review bool, requested int) (string, string, error) {
	base := roundBaseDir(ideaDir, review)
	if requested > 0 {
		label := fmt.Sprintf("round-%02d", requested)
		if info, err := os.Stat(filepath.Join(base, label)); err == nil && info.IsDir() {
			return label, roundRel(review, label), nil
		}
		return "", "", fmt.Errorf("%s does not exist", roundRel(review, label))
	}
	entries, err := os.ReadDir(base)
	if err != nil {
		return "", "", err
	}
	var rounds []roundEntry
	for _, entry := range entries {
		if !entry.IsDir() || !strings.HasPrefix(entry.Name(), "round-") {
			continue
		}
		number, err := strconv.Atoi(strings.TrimPrefix(entry.Name(), "round-"))
		if err == nil {
			rounds = append(rounds, roundEntry{name: entry.Name(), number: number})
		}
	}
	if len(rounds) == 0 {
		return "", "", fmt.Errorf("no round directories found under %s", base)
	}
	sort.Slice(rounds, func(i, j int) bool {
		if rounds[i].number == rounds[j].number {
			return rounds[i].name < rounds[j].name
		}
		return rounds[i].number < rounds[j].number
	})
	label := rounds[len(rounds)-1].name
	return label, roundRel(review, label), nil
}

type roundEntry struct {
	name   string
	number int
}

func roundBaseDir(ideaDir string, review bool) string {
	if review {
		return filepath.Join(ideaDir, "review")
	}
	return ideaDir
}

func roundRel(review bool, label string) string {
	if review {
		return filepath.Join("review", label)
	}
	return label
}

// missingRoundArtifacts reports the expected round files that do not count as a filed artifact.
//
// It used to ask only whether the pathname existed, so a file containing a single newline
// satisfied a round and `consensus draft` would advance an idea with no participant analysis in
// it at all (audit finding codex-1/F1). Existence is not filing: the protocol requires
// frontmatter and named sections, and the cheapest honest floor is a body that is not blank and
// carries at least one heading.
//
// The reason is reported per file rather than collapsed into "missing", because a blank file and
// an absent one call for different actions by whoever reads the error — and this deck has spent a
// day on messages that were true about the wrong thing.
//
// Measured before changing the floor: of 1027 round artifacts in this deck, 0 are blank and 9
// carry no heading, all in closed ideas. Nothing currently draftable is rejected by this.
func missingRoundArtifacts(roundDir string, participants []string) []string {
	var problems []string
	for _, participant := range participants {
		name := participant + ".md"
		data, err := os.ReadFile(filepath.Join(roundDir, name))
		if errors.Is(err, os.ErrNotExist) {
			problems = append(problems, name)
			continue
		}
		if err != nil {
			problems = append(problems, name+" (unreadable)")
			continue
		}
		body := stripFrontmatter(string(data))
		if strings.TrimSpace(body) == "" {
			problems = append(problems, name+" (blank)")
			continue
		}
		if !hasHeading(body) {
			problems = append(problems, name+" (no section headings)")
		}
	}
	return problems
}

// stripFrontmatter removes a leading YAML frontmatter block, if present.
func stripFrontmatter(s string) string {
	if !strings.HasPrefix(s, "---") {
		return s
	}
	rest := strings.TrimPrefix(s, "---")
	if i := strings.Index(rest, "\n---"); i >= 0 {
		return rest[i+len("\n---"):]
	}
	return rest
}

// hasHeading accepts any Markdown heading level. The floor being enforced is "this file has
// structure", not a house style: requiring exactly `## ` rejected a `# Title` artifact, which is
// a heading by any reading and was never the defect being fixed.
func hasHeading(body string) bool {
	for _, line := range strings.Split(body, "\n") {
		line = strings.TrimSpace(line)
		if !strings.HasPrefix(line, "#") {
			continue
		}
		if i := strings.IndexFunc(line, func(r rune) bool { return r != '#' }); i > 0 && strings.TrimSpace(line[i:]) != "" {
			return true
		}
	}
	return false
}

// expectedRoundParticipants is who must have filed for a round to be complete.
//
// For a REVIEW round that is not every participant: §6 forbids the implementer from reviewing its
// own work, so requiring its file made a protocol-compliant Phase 6 unreachable without
// fabricating the one artifact the protocol prohibits (audit finding codex-1/F2). When the
// implementer cannot be resolved the full list is used, which fails closed toward asking for more
// rather than silently accepting a short round.
func expectedRoundParticipants(ideaDir string, participants []string, review bool) []string {
	if !review {
		return participants
	}
	implementer := resolveImplementer(ideaDir, participants)
	if implementer == "" {
		return participants
	}
	out := make([]string, 0, len(participants))
	for _, p := range participants {
		if p != implementer {
			out = append(out, p)
		}
	}
	if len(out) == 0 {
		return participants
	}
	return out
}

// resolveImplementer reads the implementer from IMPLEMENTATION.md, else the FINAL drafter.
func resolveImplementer(ideaDir string, participants []string) string {
	isParticipant := func(id string) bool {
		for _, p := range participants {
			if p == id {
				return true
			}
		}
		return false
	}
	for _, src := range []struct {
		file string
		keys []string
	}{
		{"IMPLEMENTATION.md", []string{"implementer"}},
		{"FINAL.md", []string{"implementer", "drafted-by"}},
	} {
		meta, err := protocol.ReadFrontmatter(filepath.Join(ideaDir, src.file))
		if err != nil {
			continue
		}
		for _, k := range src.keys {
			id := strings.Trim(strings.TrimSpace(meta[k]), `"'`)
			if id != "" && isParticipant(id) {
				return id
			}
		}
	}
	return ""
}

func draftTemplate(idea protocol.IdeaStatus, opts DraftOptions, roundLabel, roundRel string) string {
	now := opts.Now
	if now.IsZero() {
		now = time.Now()
	}
	by := firstNonEmpty(opts.By, "user")
	if opts.Review {
		cycle := strings.TrimPrefix(roundLabel, "round-")
		cycle = strings.TrimLeft(cycle, "0")
		if cycle == "" {
			cycle = "1"
		}
		// The generator wrote `cycle:` while the protocol requires `review-cycle:`, and emitted
		// neither `outstanding_agreed_fixes` nor `blocked` — which the auto-driver requires as a
		// non-negative integer. So the manual command produced an artifact that violated the
		// documented schema and could not be consumed by this tool's own Phase 7/8 gate
		// (audit finding codex-1/F10). `outstanding_agreed_fixes` is seeded as a placeholder the
		// drafter must replace, not as 0: defaulting it to zero would assert "nothing outstanding"
		// on the drafter's behalf, which is the claim the gate exists to check.
		return fmt.Sprintf(`---
idea: %s
review-cycle: %s
outstanding_agreed_fixes: <count the agreed fixes below and replace this>
blocked: false
drafted-by: %s
date: %s
reviewed-commit: %s
---

## Agreed fixes

<!-- Review %s and record fixes that must be implemented. -->

## Deferred follow-ups

## Dismissed findings

## Signoffs

<!-- Each agent APPENDS their signoff block. Do NOT edit others' blocks. -->
`, idea.Slug, cycle, by, now.Format("2006-01-02"), opts.ReviewedCommit, roundRel)
	}
	return fmt.Sprintf(`---
idea: %s
drafted-by: %s
date: %s
---

## Agreed decisions

<!-- Review %s and record the decisions participants are signing off. -->

## Agreed trade-offs

## Open items deferred to implementation

## Signoffs

<!-- Each agent APPENDS their signoff block. Do NOT edit others' blocks. -->
`, idea.Slug, by, now.Format("2006-01-02"), roundRel)
}

// finalTemplate emits the FINAL.md scaffold.
//
// It used to emit `### Goal / Scope / Implementation details / Tests / Non-goals / Verification`
// under a single `## Final plan / specification` heading — a shape the protocol's Phase 4 template
// does not describe. That is why only 4 of this deck's 78 FINAL.md files carry the seven sections
// the protocol requires: the tool wrote different ones, and no gate compared them.
//
// The scaffold now carries the protocol's headings with explicit `<fill …>` placeholders, so
// `protocol.FinalIsScaffold` can tell an unwritten FINAL from a written one, and
// `consensus finalize` refuses to close an idea around it.
func finalTemplate(idea protocol.IdeaStatus, by string, now time.Time) string {
	var b strings.Builder
	fmt.Fprintf(&b, `---
idea: %s
status: final
author: %s
consensus-date: %s
participants: [%s]
---

`, idea.Slug, firstNonEmpty(by, "user"), now.Format("2006-01-02"), strings.Join(idea.Participants, ", "))

	for _, section := range protocol.RequiredFinalSections {
		fmt.Fprintf(&b, "%s\n\n", section)
		switch section {
		case "## Final plan / specification":
			b.WriteString("<fill in the agreed plan: goal, scope, implementation details, tests, non-goals, verification>\n\n")
		case "## References":
			b.WriteString("- Consensus: ./consensus.md\n- Rounds: ./round-01/\n\n")
		default:
			b.WriteString("<fill in, or write N/A if this idea is trivial or design-only>\n\n")
		}
	}
	return b.String()
}

func signoffBlock(agent, date, status, notes, counter string) string {
	var b strings.Builder
	fmt.Fprintf(&b, "\n### Signoff: %s — %s\n", agent, date)
	fmt.Fprintf(&b, "Status: %s\n", status)
	if strings.TrimSpace(notes) != "" {
		fmt.Fprintf(&b, "Notes: %s\n", strings.TrimSpace(notes))
	}
	if strings.TrimSpace(counter) != "" {
		fmt.Fprintf(&b, "Counter-proposal (required if ❌): %s\n", strings.TrimSpace(counter))
	}
	return b.String()
}

func updateIdeaStatus(ideaDir, status string) error {
	path := filepath.Join(ideaDir, "00-prompt.md")
	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	lines := strings.Split(string(data), "\n")
	inFrontmatter := false
	replaced := false
	for i, line := range lines {
		trimmed := strings.TrimSpace(line)
		if trimmed == "---" {
			if !inFrontmatter {
				inFrontmatter = true
				continue
			}
			if !replaced {
				lines = append(lines[:i], append([]string{"status: " + status}, lines[i:]...)...)
				replaced = true
			}
			break
		}
		if inFrontmatter && strings.HasPrefix(trimmed, "status:") {
			lines[i] = "status: " + status
			replaced = true
		}
	}
	if !replaced {
		return fmt.Errorf("%s has no frontmatter status field", path)
	}
	out := strings.Join(lines, "\n")
	return writeFileAtomic(path, []byte(out), 0o644)
}

func hasSectionContent(raw, heading string) bool {
	lines := strings.Split(raw, "\n")
	inSection := false
	var content []string
	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		if trimmed == heading {
			inSection = true
			continue
		}
		if inSection && strings.HasPrefix(trimmed, "## ") {
			break
		}
		if inSection && trimmed != "" && !strings.HasPrefix(trimmed, "<!--") {
			content = append(content, trimmed)
		}
	}
	return len(content) > 0
}

func nextAbortedPath(baseDir, roundLabel string) (string, error) {
	for i := 1; ; i++ {
		path := filepath.Join(baseDir, fmt.Sprintf("%s-consensus-aborted-%02d.md", roundLabel, i))
		if _, err := os.Stat(path); errors.Is(err, os.ErrNotExist) {
			return path, nil
		} else if err != nil {
			return "", err
		}
	}
}

func writeFileAtomic(path string, data []byte, perm os.FileMode) error {
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return err
	}
	tmp, err := os.CreateTemp(filepath.Dir(path), "."+filepath.Base(path)+".tmp-*")
	if err != nil {
		return err
	}
	tmpPath := tmp.Name()
	if _, err := tmp.Write(data); err != nil {
		tmp.Close()
		os.Remove(tmpPath)
		return err
	}
	if err := tmp.Close(); err != nil {
		os.Remove(tmpPath)
		return err
	}
	if err := os.Chmod(tmpPath, perm); err != nil {
		os.Remove(tmpPath)
		return err
	}
	return os.Rename(tmpPath, path)
}

func closeWith(err error, file *os.File) error {
	if closeErr := file.Close(); err == nil {
		return closeErr
	}
	return err
}

func contains(values []string, value string) bool {
	for _, existing := range values {
		if existing == value {
			return true
		}
	}
	return false
}

func firstNonEmpty(values ...string) string {
	for _, value := range values {
		if strings.TrimSpace(value) != "" {
			return strings.TrimSpace(value)
		}
	}
	return ""
}

func containsID(list []string, want string) bool {
	for _, item := range list {
		if strings.EqualFold(strings.TrimSpace(item), strings.TrimSpace(want)) {
			return true
		}
	}
	return false
}

func plural(n int, one, many string) string {
	if n == 1 {
		return one
	}
	return many
}

// unloggedReservations returns the agents that filed a reservation and are not named anywhere in
// the deferred-items section.
func unloggedReservations(doc document) []string {
	section := sectionBody(doc.Raw, "## Open items deferred to implementation")
	if section == "" {
		return nil
	}
	lower := strings.ToLower(section)
	var unlogged []string
	for _, s := range doc.Signoffs {
		if !isReserveStatus(s.Status) {
			continue
		}
		if !strings.Contains(lower, strings.ToLower(s.Agent)) {
			unlogged = append(unlogged, s.Agent)
		}
	}
	return unlogged
}

func isReserveStatus(status string) bool {
	l := strings.ToLower(status)
	return strings.Contains(l, "reserve") || strings.Contains(l, "🟡")
}

// sectionBody returns the text under a heading, up to the next `## ` heading.
func sectionBody(raw, heading string) string {
	idx := strings.Index(raw, heading)
	if idx < 0 {
		return ""
	}
	rest := raw[idx+len(heading):]
	if next := strings.Index(rest, "\n## "); next >= 0 {
		rest = rest[:next]
	}
	return rest
}

// frontmatterField reads a single scalar from a leading YAML frontmatter block.
func frontmatterField(raw, key string) string {
	if !strings.HasPrefix(raw, "---") {
		return ""
	}
	rest := strings.TrimPrefix(raw, "---")
	end := strings.Index(rest, "\n---")
	if end < 0 {
		return ""
	}
	for _, line := range strings.Split(rest[:end], "\n") {
		if value, ok := strings.CutPrefix(strings.TrimSpace(line), key+":"); ok {
			return strings.TrimSpace(value)
		}
	}
	return ""
}
