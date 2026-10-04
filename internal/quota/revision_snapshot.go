package quota

import (
	"fmt"
	"path/filepath"
	"reflect"
	"regexp"
	"strings"
	"time"

	"gopkg.in/yaml.v3"
)

// Snapshots are revalidated on every immutable-history read. A hash proves byte
// identity, not that those bytes satisfy the original confirmation contract.
func snapshotMeta(raw string) (map[string]string, error) {
	lines := strings.Split(raw, "\n")
	if len(lines) < 3 || lines[0] != "---" {
		return nil, fmt.Errorf("missing snapshot frontmatter")
	}
	meta := map[string]string{}
	for _, line := range lines[1:] {
		if line == "---" {
			return meta, nil
		}
		k, v, ok := strings.Cut(line, ":")
		if !ok {
			continue
		}
		k = strings.TrimSpace(k)
		if _, exists := meta[k]; exists && k != "excluded" && k != "included" {
			return nil, fmt.Errorf("duplicate snapshot field %s", k)
		}
		meta[k] = strings.TrimSpace(v)
	}
	return nil, fmt.Errorf("incomplete snapshot frontmatter")
}
func ConfirmedInPrompt(raw, key, id string) bool {
	lines := strings.Split(raw, "\n")
	if len(lines) == 0 || lines[0] != "---" {
		return false
	}
	for _, line := range lines[1:] {
		if line == "---" {
			break
		}
		prefix := key + ":"
		if !strings.HasPrefix(line, prefix) {
			continue
		}
		value := strings.TrimSpace(strings.TrimPrefix(line, prefix))
		value = strings.Trim(value, "[]")
		parts := strings.Split(value, " — ")
		if len(parts) != 3 || parts[0] != id || strings.TrimSpace(parts[1]) == "" || !strings.HasPrefix(parts[2], "confirmed ") {
			continue
		}
		if _, e := time.Parse("2006-01-02", strings.TrimPrefix(parts[2], "confirmed ")); e == nil {
			return true
		}
	}
	return false
}
func validateManualSnapshot(b Batch, h *History) error {
	m, e := snapshotMeta(b.Owner.ManualPrompt)
	if e != nil {
		return e
	}
	var ids []string
	if yaml.Unmarshal([]byte(m["participants"]), &ids) != nil || !reflect.DeepEqual(ids, b.Decision.After) || m["idea"] != b.Idea || m["quota_auto_exclude"] != "false" || m["quota_auto_exclude_scope"] != b.Policy.Scope {
		return fmt.Errorf("manual snapshot does not authorize recorded membership/policy")
	}
	for _, id := range h.Current {
		if !contains(ids, id) && !ConfirmedInPrompt(b.Owner.ManualPrompt, "excluded", id) {
			return fmt.Errorf("manual exclusion %s lacks recorded confirmation", id)
		}
	}
	for _, id := range ids {
		if !contains(h.Current, id) && !ConfirmedInPrompt(b.Owner.ManualPrompt, "included", id) {
			return fmt.Errorf("manual inclusion %s lacks recorded confirmation", id)
		}
	}
	return nil
}

var priorArtifact = regexp.MustCompile(`^(?:review/)?round-[0-9]{2}/[A-Za-z0-9][A-Za-z0-9._-]*\.md$`)

func ValidateCatchupSnapshot(raw, idea, id string) error {
	m, e := snapshotMeta(raw)
	if e != nil {
		return e
	}
	if strings.Trim(m["agent"], "\"'") != id || strings.Trim(m["idea"], "\"'") != idea || m["round"] != "1" || m["catch-up"] != "true" || m["join-from"] != "round-02" {
		return fmt.Errorf("invalid catch-up attribution or round")
	}
	var priors []string
	if yaml.Unmarshal([]byte(m["read-priors"]), &priors) != nil || len(priors) == 0 {
		return fmt.Errorf("catch-up requires recorded prior reading")
	}
	for _, p := range priors {
		if !priorArtifact.MatchString(p) || filepath.Base(p) == id+".md" {
			return fmt.Errorf("invalid catch-up prior %s", p)
		}
	}
	for _, heading := range []string{"## Summary", "## Proposed approach", "## Concerns / open questions", "## Risks", "## Existing alternatives"} {
		_, tail, ok := strings.Cut(raw, heading+"\n")
		if !ok {
			return fmt.Errorf("missing catch-up section %s", heading)
		}
		tail = strings.SplitN(tail, "\n## ", 2)[0]
		content := false
		for _, line := range strings.Split(tail, "\n") {
			line = strings.TrimSpace(line)
			if line != "" && !strings.HasPrefix(line, "<!--") {
				content = true
			}
		}
		if !content {
			return fmt.Errorf("empty catch-up section %s", heading)
		}
	}
	return nil
}

// The normal off-mode catch-up path reads this recorded answer from the late
// round artifact; it does not require running the quota revision command.
func CatchupAuthority(raw string) (Authority, error) {
	m, e := snapshotMeta(raw)
	if e != nil {
		return Authority{}, e
	}
	return Authority{Path: m["owner-answer"], Commit: m["owner-commit"], Blob: m["owner-blob"], SHA256: m["owner-sha256"], Quote: m["owner-quote"]}, nil
}

func CatchupPriors(raw string) ([]string, error) {
	m, e := snapshotMeta(raw)
	if e != nil {
		return nil, e
	}
	var ids []string
	e = yaml.Unmarshal([]byte(m["read-priors"]), &ids)
	return ids, e
}
