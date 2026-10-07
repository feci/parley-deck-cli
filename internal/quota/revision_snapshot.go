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
		if strings.HasPrefix(value, "[") {
			if !strings.HasSuffix(value, "]") {
				continue
			}
			value = value[1 : len(value)-1]
		}
		value = strings.TrimSpace(value)
		prefixID := id + " — "
		if !strings.HasPrefix(value, prefixID) {
			continue
		}
		reason, date, ok := strings.Cut(value[len(prefixID):], " — confirmed ")
		// Reasons can themselves contain em dashes (and the word confirmed).
		if at := strings.LastIndex(value[len(prefixID):], " — confirmed "); at >= 0 {
			reason, date, ok = value[len(prefixID):len(prefixID)+at], value[len(prefixID)+at+len(" — confirmed "):], true
		}
		if !ok || strings.TrimSpace(reason) == "" {
			continue
		}
		if _, e := time.Parse("2006-01-02", date); e == nil {
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
		return fmt.Errorf("manual snapshot does not match recorded membership/policy")
	}
	for _, id := range h.Current {
		if !contains(ids, id) && !ConfirmedInPrompt(b.Owner.ManualPrompt, "excluded", id) {
			return ManualExclusionError(id)
		}
	}
	return nil
}

var priorArtifact = regexp.MustCompile(`^(?:review/)?round-[0-9]{2}/[A-Za-z0-9][A-Za-z0-9._-]*\.md$`)

func ManualExclusionError(id string) error {
	return fmt.Errorf("manual exclusion %s requires excluded: [%s — reason — confirmed YYYY-MM-DD] (brackets optional; em dashes allowed in reason); record the owner's confirmation in that form, with no suffix after the date", id, id)
}

func ValidateCatchupSnapshot(raw, idea, id string) error {
	return validateCatchupSnapshot(raw, idea, id, true)
}

// A policy-off catch-up uses the pre-change late round-1 path. Reading priors and
// joining from round 2 remain protocol duties, not a newly mandated record format.
func ValidateManualCatchupSnapshot(raw, idea, id string) error {
	return validateCatchupSnapshot(raw, idea, id, false)
}

func validateCatchupSnapshot(raw, idea, id string, explicit bool) error {
	m, e := snapshotMeta(raw)
	if e != nil {
		return e
	}
	if strings.Trim(m["agent"], "\"'") != id || strings.Trim(m["idea"], "\"'") != idea || m["round"] != "1" {
		return fmt.Errorf("invalid catch-up attribution or round")
	}
	if explicit {
		if m["catch-up"] != "true" || m["join-from"] != "round-02" {
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

// CatchupAuthority decodes historical explicit authority metadata. The ordinary
// policy-off compatibility path does not require or infer this authority.
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

// A policy-off round-1 participant edit permits a return or a new identity.
// Only the revision's prompt snapshot establishes the round; display markers
// establish neither historical membership nor owner authority.
func ManualRoundOneReturn(raw string) bool {
	m, err := snapshotMeta(raw)
	return err == nil && m["status"] == "round-01"
}
