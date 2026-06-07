package protocol

import (
	"bufio"
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

func InitWorkspace(root string) error {
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
	if _, err := os.Stat(cooperation); err == nil {
		return nil
	} else if !errors.Is(err, os.ErrNotExist) {
		return err
	}

	return os.WriteFile(cooperation, []byte(defaultCooperationForInit()), 0o644)
}

func defaultCooperationForInit() string {
	return strings.Replace(defaultCooperation, "**Transport:** `github-pr`", "**Transport:** `local-dir`", 1)
}

func CreateIdea(root, task string, participants []string) (IdeaStatus, error) {
	now := time.Now()
	slug := uniqueSlug(filepath.Join(root, DeckDir, "ideas"), timestampedSlug(task, now))
	ideaDir := filepath.Join(root, DeckDir, "ideas", slug)
	if err := fsutil.MkdirAllResilient(filepath.Join(ideaDir, "round-01"), 0o755); err != nil {
		return IdeaStatus{}, err
	}

	prompt := fmt.Sprintf(`---
idea: %s
author: user
created: %s
participants: [%s]
status: round-01
---

## Problem / idea

%s

## Constraints

- Local-directory transport for this initial run.
- Human-in-the-loop mode by default unless the run was started with auto mode.

## Non-goals

- Do not make unrelated repository changes.
`, slug, now.Format("2006-01-02"), strings.Join(participants, ", "), task)
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
