package protocol

import (
	"bufio"
	"crypto/sha256"
	_ "embed"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"time"
	"unicode"

	"parley-deck-cli/internal/fsutil"
)

const DeckDir = "parley-deck"

//go:embed defaults/COOPERATION.md
var defaultCooperation string

type WorkspaceStatus struct {
	Root      string
	Transport string
	Ideas     []IdeaStatus
}

type IdeaStatus struct {
	Slug         string
	Status       string
	Participants []string
	Path         string
}

// InitWorkspace creates a fresh deck with the default (local-dir) transport.
func InitWorkspace(root string) error {
	return InitWorkspaceWithTransport(root, "")
}

// InitWorkspaceWithTransport creates a fresh deck, seeding COOPERATION.md with
// the requested transport. An empty or unknown transport falls back to
// local-dir; valid values are local-dir, github-pr, and gitlab-mr.
func InitWorkspaceWithTransport(root, transport string) error {
	deck := filepath.Join(root, DeckDir)
	for _, dir := range []string{
		deck,
		filepath.Join(deck, "ideas"),
		filepath.Join(deck, "inbox"),
		filepath.Join(deck, "meta"),
		filepath.Join(deck, "runs"),
	} {
		if err := os.MkdirAll(dir, 0o755); err != nil {
			return err
		}
	}

	cooperation := filepath.Join(deck, "COOPERATION.md")
	switch _, err := os.Stat(cooperation); {
	case err == nil:
		// Pre-existing protocol: still record metadata, but never clobber the file.
	case errors.Is(err, os.ErrNotExist):
		body := cooperationForInitAt(transport, root, time.Now())
		if werr := os.WriteFile(cooperation, []byte(body), 0o644); werr != nil {
			return werr
		}
	default:
		return err
	}

	return writeInitVersionMeta(deck, cooperation)
}

// writeInitVersionMeta writes meta/version.json with protocolRole:"consumer" so a
// fresh workspace enters the consumer freshness path (preflight §9.0) instead of
// failing open on absent metadata. It never clobbers an existing version.json.
//
// It records `protocolSha256` for the COOPERATION.md the deck actually starts with, so the
// metadata describes this deck from its first minute. `packagedProtocolSha256` is deliberately
// NOT written: `init` does not read the installed skill, and the deck's protocol is not
// guaranteed to be byte-identical to the packaged one. Writing a guess there would make preflight
// report "in sync" on the strength of a hash nobody computed.
func writeInitVersionMeta(deck, cooperation string) error {
	path := filepath.Join(deck, "meta", "version.json")
	if _, err := os.Stat(path); err == nil {
		return nil
	} else if !errors.Is(err, os.ErrNotExist) {
		return err
	}
	sum := ""
	if body, err := os.ReadFile(cooperation); err == nil {
		sum = fmt.Sprintf("%x", sha256.Sum256(body))
	}
	meta := fmt.Sprintf(`{
  "protocolRole": "consumer",
  "deckVersion": "",
  "protocolSha256": "%s",
  "created": "%s"
}
`, sum, time.Now().Format("2006-01-02"))
	return os.WriteFile(path, []byte(meta), 0o644)
}

// defaultCooperationForInit is the STATIC bootstrap template: transport swapped, every other
// placeholder deliberately intact. The drift guard pins it, and the skill ships the same shape as
// a vendor-neutral reference, so it must not acquire this machine's date or directory name.
func defaultCooperationForInit() string {
	target := "local-dir"
	return strings.Replace(defaultCooperation, "**Transport:** `github-pr`", "**Transport:** `"+target+"`", 1)
}

// cooperationForInit returns the embedded default protocol with its transport
// line set to the requested transport. Empty or unknown transports default to
// local-dir.
func cooperationForInit(transport string) string {
	return cooperationForInitAt(transport, "", time.Now())
}

// cooperationForInitAt fills the embedded template's placeholders.
//
// `init` used to substitute the transport and nothing else, so every CLI-created deck started
// with `**Workspace:** <workspace-name>` and `**Created:** <date> — created by parley init`
// still literally in the header — false provenance on the protocol's own first lines, while the
// command reported initialization complete (audit finding codex-1/F19).
//
// The workspace name is the deck root's directory name, which is what a human would write there;
// an empty or unusable root leaves the placeholder rather than inventing an identity.
func cooperationForInitAt(transport, root string, now time.Time) string {
	target := "local-dir"
	switch transport {
	case "github-pr":
		target = "github-pr"
	case "gitlab-mr":
		target = "gitlab-mr"
	}
	out := strings.Replace(defaultCooperation, "**Transport:** `github-pr`", "**Transport:** `"+target+"`", 1)
	out = strings.Replace(out, "**Created:** `<date> — created by parley init`",
		"**Created:** `"+now.Format("2006-01-02")+" — created by parley init`", 1)
	if name := workspaceName(root); name != "" {
		out = strings.Replace(out, "**Workspace:** `<workspace-name>`", "**Workspace:** `"+name+"`", 1)
	}
	return out
}

func workspaceName(root string) string {
	abs, err := filepath.Abs(root)
	if err != nil {
		return ""
	}
	name := filepath.Base(abs)
	if name == "." || name == string(filepath.Separator) || name == "" {
		return ""
	}
	return name
}

func CreateIdea(root, task string, participants []string) (IdeaStatus, error) {
	return CreateIdeaWithExclusions(root, task, participants, nil)
}

// CreateIdeaWithExclusions creates the idea like CreateIdea but also records any
// confirmed participant exclusions in the frontmatter as `excluded:` lines
// (preflight §9.0: exclusions must be explicit and recorded in the idea).
func CreateIdeaWithExclusions(root, task string, participants, excluded []string) (IdeaStatus, error) {
	return CreateIdeaFull(root, task, participants, excluded, "", "")
}

// CreateIdeaFull is CreateIdeaWithExclusions plus an optional `track:` frontmatter
// line (track-aware-driver) and a roster-preset provenance HTML comment written under
// the participants line (named-roster-presets). Both are advisory: `participants:`
// stays the canonical quorum. Empty track/provenance reproduce the base behavior.
func CreateIdeaFull(root, task string, participants, excluded []string, track, provenance string) (IdeaStatus, error) {
	now := time.Now()
	slug := uniqueSlug(filepath.Join(root, DeckDir, "ideas"), timestampedSlug(task, now))
	ideaDir := filepath.Join(root, DeckDir, "ideas", slug)
	if err := fsutil.MkdirAllResilient(filepath.Join(ideaDir, "round-01"), 0o755); err != nil {
		return IdeaStatus{}, err
	}

	excludedBlock := ""
	for _, line := range excluded {
		if strings.TrimSpace(line) == "" {
			continue
		}
		excludedBlock += "excluded: " + line + "\n"
	}
	trackLine := ""
	if t := strings.TrimSpace(track); t != "" {
		trackLine = "track: " + t + "\n"
	}
	// Provenance is an HTML comment BELOW the frontmatter fence (review fix: inside the
	// fence, ReadFrontmatter's `key: value` split would ingest it as a junk key).
	provenanceBlock := ""
	if p := strings.TrimSpace(provenance); p != "" {
		provenanceBlock = p + "\n\n"
	}

	prompt := fmt.Sprintf(`---
idea: %s
author: user
created: %s
participants: [%s]
%s%sstatus: round-01
---

%s## Problem / idea

%s

## Constraints

- Local-directory transport for this initial run.
- Human-in-the-loop mode by default unless the run was started with auto mode.

## Non-goals

- Do not make unrelated repository changes.
`, slug, now.Format("2006-01-02"), strings.Join(participants, ", "), trackLine, excludedBlock, provenanceBlock, task)
	if err := os.WriteFile(filepath.Join(ideaDir, "00-prompt.md"), []byte(prompt), 0o644); err != nil {
		return IdeaStatus{}, err
	}

	return IdeaStatus{
		Slug:         slug,
		Status:       "round-01",
		Participants: participants,
		Path:         ideaDir,
	}, nil
}

func ReadWorkspaceStatus(root string) (WorkspaceStatus, error) {
	deck := filepath.Join(root, DeckDir)
	transport, err := readTransport(filepath.Join(deck, "COOPERATION.md"))
	if err != nil {
		return WorkspaceStatus{}, err
	}

	ideas, err := readIdeas(filepath.Join(deck, "ideas"))
	if err != nil {
		return WorkspaceStatus{}, err
	}

	return WorkspaceStatus{
		Root:      root,
		Transport: transport,
		Ideas:     ideas,
	}, nil
}

func uniqueSlug(ideasDir, base string) string {
	if base == "" {
		base = "task"
	}
	slug := base
	for i := 2; ; i++ {
		if _, err := os.Stat(filepath.Join(ideasDir, slug)); errors.Is(err, os.ErrNotExist) {
			return slug
		}
		slug = fmt.Sprintf("%s-%d", base, i)
	}
}

func timestampedSlug(task string, now time.Time) string {
	prefix := truncateSlug(slugify(task), 16)
	if prefix == "" {
		prefix = "task"
	}
	return now.Format("2006-01-02T15-04-05") + "-" + prefix
}

func truncateSlug(slug string, maxRunes int) string {
	if maxRunes <= 0 {
		return ""
	}
	var b strings.Builder
	count := 0
	for _, r := range slug {
		if count >= maxRunes {
			break
		}
		b.WriteRune(r)
		count++
	}
	return strings.Trim(b.String(), "-")
}

func slugify(value string) string {
	value = strings.ToLower(value)
	var b strings.Builder
	lastDash := false
	for _, r := range value {
		switch {
		case unicode.IsLetter(r), unicode.IsDigit(r):
			b.WriteRune(r)
			lastDash = false
		case !lastDash:
			b.WriteByte('-')
			lastDash = true
		}
	}
	return strings.Trim(b.String(), "-")
}

func readTransport(path string) (string, error) {
	file, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer file.Close()

	re := regexp.MustCompile(`^\*\*Transport:\*\*\s*` + "`?" + `([^` + "`" + `\s]+)`)
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if match := re.FindStringSubmatch(line); len(match) == 2 {
			return match[1], nil
		}
	}
	if err := scanner.Err(); err != nil {
		return "", err
	}
	return "unknown", nil
}

func readIdeas(path string) ([]IdeaStatus, error) {
	entries, err := os.ReadDir(path)
	if errors.Is(err, os.ErrNotExist) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	var ideas []IdeaStatus
	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}
		prompt := filepath.Join(path, entry.Name(), "00-prompt.md")
		meta, err := ReadFrontmatter(prompt)
		if errors.Is(err, os.ErrNotExist) {
			continue
		}
		if err != nil {
			return nil, fmt.Errorf("%s: %w", prompt, err)
		}
		ideas = append(ideas, IdeaStatus{
			Slug:         first(meta["idea"], entry.Name()),
			Status:       first(meta["status"], "unknown"),
			Participants: parseList(meta["participants"]),
			Path:         filepath.Dir(prompt),
		})
	}

	sort.Slice(ideas, func(i, j int) bool {
		return ideas[i].Slug < ideas[j].Slug
	})
	return ideas, nil
}

func ReadFrontmatter(path string) (map[string]string, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	meta := map[string]string{}
	inFrontmatter := false
	for scanner.Scan() {
		line := scanner.Text()
		if strings.TrimSpace(line) == "---" {
			if !inFrontmatter {
				inFrontmatter = true
				continue
			}
			return meta, nil
		}
		if !inFrontmatter {
			return meta, nil
		}
		key, value, ok := strings.Cut(line, ":")
		if !ok {
			continue
		}
		meta[strings.TrimSpace(key)] = strings.TrimSpace(value)
	}
	return meta, scanner.Err()
}

func parseList(raw string) []string {
	raw = strings.TrimSpace(raw)
	raw = strings.TrimPrefix(raw, "[")
	raw = strings.TrimSuffix(raw, "]")
	if raw == "" {
		return nil
	}
	parts := strings.Split(raw, ",")
	values := make([]string, 0, len(parts))
	for _, part := range parts {
		part = strings.Trim(strings.TrimSpace(part), `"'`)
		if part != "" {
			values = append(values, part)
		}
	}
	return values
}

func first(value, fallback string) string {
	if strings.TrimSpace(value) == "" {
		return fallback
	}
	return strings.TrimSpace(value)
}
