package protocol

import (
	"crypto/sha256"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"strings"

	"parley-deck-cli/internal/quota"
)

var retainedSignoff = regexp.MustCompile(`(?m)^### Signoff: ([A-Za-z0-9._-]+) [—-] .*$`)
var retainedDispute = regexp.MustCompile(`\bDISPUTED\b`)
var retainedNoDispute = regexp.MustCompile(`(?i)\b(no|not|zero) (?:unresolved )?DISPUTED\b`)
var retainedFinding = regexp.MustCompile(`(?mi)^[ \t]*#{2,6}[ \t]*\[[ \t]*(CRITICAL|MAJOR|MINOR|NIT)[ \t]*\][^\n]*`)

// CaptureQuotaObligations records only already-filed work from the removed
// identities. It runs AFTER the provider-based decision, never as its trigger.
func CaptureQuotaObligations(ideaDir string, removed []string) ([]quota.Obligation, error) {
	wanted := map[string]bool{}
	for _, id := range removed {
		wanted[id] = true
	}
	meta, err := ReadFrontmatter(filepath.Join(ideaDir, "00-prompt.md"))
	if err != nil {
		return nil, err
	}
	strict := strings.EqualFold(meta["strict_gate"], "true")
	out := []quota.Obligation{}
	add := func(kind, agent, path, raw, text string) {
		sum := fmt.Sprintf("%x", sha256.Sum256([]byte(raw)))
		id := fmt.Sprintf("obligation-%x", sha256.Sum256([]byte(kind+"\x00"+agent+"\x00"+path+"\x00"+text)))
		out = append(out, quota.Obligation{ID: id, Kind: kind, Agent: agent, Path: path, SHA256: sum, Text: text})
	}
	err = filepath.WalkDir(ideaDir, func(path string, entry os.DirEntry, e error) error {
		if e != nil {
			return e
		}
		if entry.IsDir() {
			if path != ideaDir && (entry.Name() == quota.HistoryDir || entry.Name() == "quota-applied") {
				return filepath.SkipDir
			}
			return nil
		}
		if !strings.HasSuffix(path, ".md") {
			return nil
		}
		rel, _ := filepath.Rel(ideaDir, path)
		rel = filepath.ToSlash(rel)
		// Only canonical participant rounds/reviews and signoff artifacts.
		isConsensus := strings.Contains(filepath.Base(path), "consensus") || filepath.Base(path) == "FINAL.md"
		agent := strings.TrimSuffix(filepath.Base(path), ".md")
		isParticipant := wanted[agent] && (strings.HasPrefix(rel, "round-") || strings.HasPrefix(rel, "review/round-"))
		if !isConsensus && !isParticipant {
			return nil
		}
		data, e := os.ReadFile(path)
		if e != nil {
			return e
		}
		raw := string(data)
		if isConsensus {
			matches := retainedSignoff.FindAllStringSubmatchIndex(raw, -1)
			for i, m := range matches {
				a := raw[m[2]:m[3]]
				if !wanted[a] {
					continue
				}
				end := len(raw)
				if i+1 < len(matches) {
					end = matches[i+1][0]
				}
				block := raw[m[0]:end]
				if strings.Contains(block, "❌") || strings.Contains(block, "Status: BLOCK") {
					add("veto", a, rel, raw, strings.TrimSpace(block))
				}
			}
		}
		if isParticipant {
			for _, m := range retainedFinding.FindAllStringSubmatch(raw, -1) {
				if strict || strings.EqualFold(m[1], "CRITICAL") || strings.EqualFold(m[1], "MAJOR") {
					add("finding", agent, rel, raw, m[0])
				}
			}
			for _, line := range strings.Split(raw, "\n") {
				if retainedDispute.MatchString(line) && !retainedNoDispute.MatchString(line) {
					add("dispute", agent, rel, raw, line)
				}
			}
		}
		return nil
	})
	return out, err
}

// UnresolvedQuotaObligations is a conservative close veto, never an automatic
// evidence verdict. A disposition names the immutable obligation, rationale,
// authority and an independent evidence artifact. A veto specifically needs an
// owner ruling quoted into that next artifact; exclusion cannot withdraw it.
// Records use the existing disposition vocabulary in canonical consensus prose.
func UnresolvedQuotaObligations(ideaDir, consensusPath string) ([]quota.Obligation, error) {
	h, err := quota.ReadHistory(ideaDir)
	if err != nil || h == nil {
		return nil, err
	}
	raw, err := os.ReadFile(consensusPath)
	if err != nil {
		return nil, err
	}
	unresolved := []quota.Obligation{}
	seen := map[string]bool{}
	for _, b := range h.Batches {
		for _, ob := range b.Retained {
			if seen[ob.ID] {
				continue
			}
			seen[ob.ID] = true
			if !quotaDisposition(ideaDir, string(raw), ob) {
				unresolved = append(unresolved, ob)
			}
		}
	}
	return unresolved, nil
}
func quotaDisposition(ideaDir, raw string, ob quota.Obligation) bool {
	// Each disposition is a contiguous paragraph, so fields from unrelated claims
	// cannot accidentally satisfy one another. Paths are local canonical evidence.
	for _, paragraph := range strings.Split(raw, "\n\n") {
		if !strings.Contains(paragraph, ob.ID) {
			continue
		}
		fields := map[string]string{}
		duplicate := false
		for _, line := range strings.Split(paragraph, "\n") {
			key, value, ok := strings.Cut(strings.TrimSpace(line), ":")
			if ok {
				if _, exists := fields[key]; exists {
					duplicate = true
				}
				fields[key] = strings.TrimSpace(value)
			}
		}
		if duplicate || fields["Rationale"] == "" || fields["Authority"] == "" || fields["Evidence"] == "" {
			continue
		}
		d := fields["Disposition"]
		if d != "resolved" && d != "withdrawn" && d != "operator-ruling" && d != "no-dependency" {
			continue
		}
		if ob.Kind == "veto" && d != "operator-ruling" && d != "withdrawn" {
			continue
		}
		path := filepath.Clean(fields["Evidence"])
		if filepath.IsAbs(path) || path == ".." || strings.HasPrefix(path, ".."+string(filepath.Separator)) {
			continue
		}
		path = filepath.Join(ideaDir, path)
		if path == filepath.Join(ideaDir, "consensus.md") || path == filepath.Join(ideaDir, "review", "consensus.md") {
			continue
		}
		data, err := os.ReadFile(path)
		if err != nil {
			continue
		}
		meta, err := ReadFrontmatter(path)
		if err != nil {
			continue
		}
		author := strings.Trim(meta["agent"], "\"'")
		if author == "" {
			author = strings.Trim(meta["author"], "\"'")
		}
		rel, _ := filepath.Rel(ideaDir, path)
		canonical := regexp.MustCompile(`^(?:review/)?round-[0-9]{2}/` + regexp.QuoteMeta(author) + `\.md$`).MatchString(filepath.ToSlash(rel))
		if author == "" || !canonical || meta["idea"] != filepath.Base(ideaDir) || !strings.Contains(string(data), ob.ID) {
			continue
		}
		history, e := quota.ReadHistory(ideaDir)
		if e != nil || history == nil || !quotaContains(history.Known, author) {
			continue
		}
		if ob.Kind != "veto" {
			if author == ob.Agent {
				continue
			}
			pin, _ := ReadFrontmatter(filepath.Join(ideaDir, "IMPLEMENTATION.md"))
			if author == strings.Trim(pin["implementer"], "\"'") {
				continue
			}
		}
		if ob.Kind == "veto" {
			if d == "withdrawn" {
				h, e := quota.ReadHistory(ideaDir)
				revision, e2 := strconv.Atoi(meta["quota-revision"])
				if e != nil || e2 != nil || h == nil || revision < 1 || revision > len(h.Batches) || author != ob.Agent || fields["Authority"] != ob.Agent || !strings.Contains(string(data), "Withdrawal: "+ob.ID) {
					continue
				}
				r := h.Batches[revision-1]
				if meta["reinclusion"] != r.ID || r.Owner == nil || (r.Owner.Authority == nil && r.Owner.ManualPrompt == "") || !quotaContains(r.Decision.After, author) || quotaContains(r.Decision.Before, author) || !quotaContains(h.Current, author) {
					continue
				}
				oldRound := quotaArtifactRound(ob.Path)
				if quotaArtifactRound(filepath.ToSlash(rel)) <= oldRound {
					continue
				}
				return true
			}
			if author == ob.Agent || fields["Authority"] != "owner" {
				continue
			}
			a := quota.Authority{Path: fields["Owner-answer"], Commit: fields["Owner-commit"], Blob: fields["Owner-blob"], SHA256: fields["Owner-sha256"], Quote: fields["Owner-quote"]}
			_, direction, hasDirection := strings.Cut(string(data), "## User direction\n")
			direction = strings.SplitN(direction, "\n## ", 2)[0]
			if !strings.Contains(a.Quote, ob.ID) || !hasDirection || !strings.Contains(direction, a.Quote) || quota.ValidateAuthority(quota.IdeaRoot(ideaDir), filepath.Base(ideaDir), a) != nil {
				continue
			}
		}
		if strings.Contains(string(data), "Evidence:") || strings.Contains(string(data), "PRIMARY") || ob.Kind == "veto" {
			return true
		}
	}
	return false
}

// ResolvedQuotaVeto reports only an explicit owner ruling linked to EVERY retained
// veto by this historical signer. The filed signoff bytes remain untouched.
func ResolvedQuotaVeto(ideaDir, consensusPath, agent string) bool {
	return ResolvedQuotaVetoBlock(ideaDir, consensusPath, agent, "")
}

func ResolvedQuotaVetoBlock(ideaDir, consensusPath, agent, block string) bool {
	_, ok := ResolvedQuotaVetoIdentity(ideaDir, consensusPath, agent, block)
	return ok
}

// Identity binds the exact preserved original signoff bytes, including blank
// lines, to retained obligation IDs. A later identical veto can consume that
// historical identity only once; altered/new vetoes do not match it.
func ResolvedQuotaVetoIdentity(ideaDir, consensusPath, agent, suffix string) (string, bool) {
	h, err := quota.ReadHistory(ideaDir)
	if err != nil || h == nil {
		return "", false
	}
	raw, err := os.ReadFile(consensusPath)
	if err != nil {
		return "", false
	}
	ids := []string{}
	seen := map[string]bool{}
	for _, b := range h.Batches {
		for _, ob := range b.Retained {
			if ob.Kind != "veto" || ob.Agent != agent || suffix != "" && !strings.HasPrefix(suffix, ob.Text) {
				continue
			}
			if !quotaDisposition(ideaDir, string(raw), ob) {
				return "", false
			}
			if !seen[ob.ID] {
				ids = append(ids, ob.ID)
				seen[ob.ID] = true
			}
		}
	}
	sort.Strings(ids)
	return strings.Join(ids, ","), len(ids) > 0
}

func quotaContains(ids []string, id string) bool {
	for _, x := range ids {
		if x == id {
			return true
		}
	}
	return false
}
func quotaArtifactRound(path string) int {
	for _, part := range strings.Split(filepath.ToSlash(path), "/") {
		if strings.HasPrefix(part, "round-") {
			n, _ := strconv.Atoi(strings.TrimPrefix(part, "round-"))
			return n
		}
	}
	return 0
}
