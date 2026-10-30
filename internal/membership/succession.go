package membership

import (
	"crypto/sha256"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"parley-deck-cli/internal/config"
	"parley-deck-cli/internal/protocol"
	"parley-deck-cli/internal/quota"
	"parley-deck-cli/internal/telemetry"
)

func resolveSuccessionRoles(ideaDir string, r quota.Roles) (quota.Roles, error) {
	h, err := quota.ReadHistory(ideaDir)
	if err != nil || h == nil || h.RoleSnapshot() == nil {
		return r, err
	}
	r.Snapshot = h.RoleSnapshot()
	r.Facilitator = r.Snapshot.Facilitator
	if previous := h.LastRoleTransition(); previous != nil {
		r.Assignment = previous.After
		r.Designee, r.PinnedImplementer = "", previous.After.Implementer
		return r, nil
	}
	r.Assignment.Organizer = r.Snapshot.Organizer
	meta, err := protocol.ReadFrontmatter(filepath.Join(ideaDir, "00-prompt.md"))
	if err != nil {
		return r, err
	}
	organizer := strings.Trim(meta["organizer"], "\"'")
	if organizer == "" {
		organizer = strings.Trim(meta["facilitator"], "\"'")
	}
	if organizer != r.Snapshot.Organizer || strings.Trim(meta["facilitator"], "\"'") != r.Snapshot.Facilitator {
		return r, fmt.Errorf("prompt role identity differs from frozen authority")
	}
	des := protocol.ImplementerDesignationFromMeta(meta)
	if des.State == protocol.DesignationEmpty {
		return r, fmt.Errorf("incomplete implementer designation")
	}
	eligible := append([]string(nil), h.Current...)
	if r.Facilitator != "" && !protocol.ReadFacilitatorRole(ideaDir).Participates {
		eligible = nil
		for _, id := range h.Current {
			if id != r.Facilitator {
				eligible = append(eligible, id)
			}
		}
	}
	if des.State == protocol.DesignationSet && !Has(eligible, des.ID) {
		return r, fmt.Errorf("invalid per-idea implementer designation")
	}
	if r.PinnedImplementer != "" && !Has(eligible, r.PinnedImplementer) {
		return r, fmt.Errorf("implementation pin is outside eligible membership")
	}
	if r.PinnedImplementer != "" && des.State == protocol.DesignationSet && des.ID != r.PinnedImplementer {
		if !protocol.ParseImplementerReassignment(meta, r.PinnedImplementer, des.ID) {
			return r, fmt.Errorf("implementer pin conflicts with designation")
		}
		r.PinnedImplementer = ""
	}
	chain := []protocol.ImplementerCandidate{}
	if r.PinnedImplementer != "" {
		chain = append(chain, protocol.ImplementerCandidate{Source: protocol.SourceImplementerPin, ID: r.PinnedImplementer})
	}
	if des.State == protocol.DesignationSet && !protocol.ImplementerWaived(meta, des.ID) {
		chain = append(chain, protocol.ImplementerCandidate{Source: protocol.SourceImplementerDesignation, ID: des.ID})
	} else if des.State == protocol.DesignationAbsent {
		defs, err := config.LoadDefaults(quota.IdeaRoot(ideaDir))
		if err != nil {
			return r, err
		}
		g := defs.DefaultImplementer
		if g != "" && !strings.EqualFold(g, "none") {
			if protocol.ImplementerDesignationFromMeta(map[string]string{"implementer": g}).ID != g || strings.ContainsAny(g, " \t\n#") {
				return r, fmt.Errorf("malformed global implementer designation")
			}
			chain = append(chain, protocol.ImplementerCandidate{Source: protocol.SourceImplementerGlobalDefault, ID: g})
		}
	}
	chain = append(chain, protocol.LegacyImplementerCandidates()...)
	id, source, ok := protocol.ResolveImplementerChain(ideaDir, chain, eligible)
	if !ok && len(eligible) > 0 {
		id, source = eligible[0], protocol.SourceImplementerLegacy
	}
	r.Assignment.Implementer, r.Assignment.ImplementerSource = id, string(source)
	return r, nil
}

func bindSuccession(root, ideaDir, phase string, d quota.Decision) (quota.Decision, error) {
	if d.RoleTransition == nil {
		return d, nil
	}
	a, err := protocol.ResolveActiveExecution(ideaDir)
	if err != nil {
		return d, err
	}
	base, hash := "none", ""
	raw, err := os.ReadFile(a.Path)
	if err == nil {
		meta, err := protocol.ReadFrontmatter(a.Path)
		if err != nil {
			return d, err
		}
		if meta["status"] == "complete" {
			return d, fmt.Errorf("closed execution record stays frozen")
		}
		rel, err := filepath.Rel(ideaDir, a.Path)
		if err != nil {
			return d, err
		}
		base, hash = filepath.ToSlash(rel), fmt.Sprintf("%x", sha256.Sum256(raw))
	} else if !os.IsNotExist(err) {
		return d, err
	}
	gitValue := func(args ...string) string {
		cmd := exec.Command("git", append([]string{"-C", root}, args...)...)
		out, err := cmd.Output()
		if err != nil {
			return "n/a"
		}
		return strings.TrimSpace(string(out))
	}
	return d.BindHandover(phase, base, hash, gitValue("branch", "--show-current"), gitValue("rev-parse", "HEAD")), nil
}

// Valid phase artifacts supply positive usability without a provider probe.
func SupplementUsability(ideaDir string, ids []string, observed []quota.Member) ([]quota.Member, error) {
	seen := map[string]bool{}
	for _, m := range observed {
		if m.Usable || m.ValidArtifact || m.LaterSuccess {
			if err := telemetry.ExhaustionAdmission(quota.IdeaRoot(ideaDir), filepath.Base(ideaDir), m.ID); err != nil {
				return nil, err
			}
		}
		seen[m.ID] = true
	}
	for _, id := range ids {
		if seen[id] {
			continue
		}
		if err := telemetry.ExhaustionAdmission(quota.IdeaRoot(ideaDir), filepath.Base(ideaDir), id); err != nil {
			return nil, err
		}
		paths, _ := filepath.Glob(filepath.Join(ideaDir, "round-*", id+".md"))
		for i := len(paths) - 1; i >= 0; i-- {
			var n int
			if _, err := fmt.Sscanf(filepath.Base(filepath.Dir(paths[i])), "round-%d", &n); err == nil && protocol.ValidateParticipantRoundArtifact(paths[i], id, filepath.Base(ideaDir), n) == nil {
				observed = append(observed, quota.Member{ID: id, Usable: true, ValidArtifact: true})
				break
			}
		}
	}
	return observed, nil
}

func reconcileDrafts(ideaDir string, b quota.Batch, applied bool) error {
	for _, a := range b.Decision.DraftArchives {
		source, archive := filepath.Join(ideaDir, filepath.FromSlash(a.Path)), filepath.Join(ideaDir, filepath.FromSlash(a.Archive))
		archived, err := os.ReadFile(archive)
		if os.IsNotExist(err) {
			if applied {
				return fmt.Errorf("committed draft archive is missing")
			}
			raw, e := os.ReadFile(source)
			if e != nil || fmt.Sprintf("%x", sha256.Sum256(raw)) != a.SHA256 {
				return fmt.Errorf("draft archive source changed")
			}
			if e = quota.DurableWrite(archive, raw, true); e != nil {
				return e
			}
			archived = raw
		} else if err != nil {
			return err
		}
		if fmt.Sprintf("%x", sha256.Sum256(archived)) != a.SHA256 {
			return fmt.Errorf("draft archive hash mismatch")
		}
		if applied {
			continue
		}
		raw, err := os.ReadFile(source)
		if os.IsNotExist(err) {
			continue
		}
		if err != nil {
			return err
		}
		if fmt.Sprintf("%x", sha256.Sum256(raw)) != a.SHA256 {
			return fmt.Errorf("draft changed before archive transaction completed")
		}
		if err = os.Remove(source); err != nil {
			return err
		}
		if err = quota.SyncPath(filepath.Dir(source)); err != nil {
			return err
		}
	}
	return nil
}
