package app

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"parley-deck-cli/internal/config"
	"parley-deck-cli/internal/protocol"
)

// Fleet migration: move each deck's membership out of the hand-edited §2 table and into
// [roster.<id>] blocks in its own parley-deck/agents.toml.
//
// This mutates OTHER people's repositories, several of them years old, so every step is
// built to be refusable and reversible: a deck that is not obviously clean is SKIPPED and
// reported rather than guessed at, every write is backed up first and verified after, a
// failed verification rolls that deck back, and a re-run is a no-op on decks already done.

type migrateDisposition string

const (
	migrateApplied   migrateDisposition = "applied"
	migrateSkipped   migrateDisposition = "skipped"
	migrateUnchanged migrateDisposition = "unchanged"
	migrateFailed    migrateDisposition = "failed-and-restored"
)

type migrateDeck struct {
	Deck        string             `json:"deck"`
	Disposition migrateDisposition `json:"disposition"`
	Reason      string             `json:"reason,omitempty"`
	Members     []string           `json:"members,omitempty"`
	Inactive    []string           `json:"inactive,omitempty"`
	BeforeHash  string             `json:"before_hash,omitempty"`
	AfterHash   string             `json:"after_hash,omitempty"`
	Backup      string             `json:"backup,omitempty"`
}

type migrateReport struct {
	SchemaVersion int           `json:"schema_version"`
	Root          string        `json:"root"`
	DryRun        bool          `json:"dry_run"`
	Decks         []migrateDeck `json:"decks"`
	Applied       int           `json:"applied"`
	Skipped       int           `json:"skipped"`
	Unchanged     int           `json:"unchanged"`
	Failed        int           `json:"failed"`
}

// rosterMigrate walks every deck under root and moves its §2 membership into config.
func rosterMigrate(root, backupDir string, dryRun, yes, jsonOut, confirmBreaking bool, stdout, stderr io.Writer) int {
	// INTERIM GUARD (review DF-1). The ratified migration contract also requires
	// compare-and-swap between preview and apply, per-deck confirmation honoring
	// roster_change_policy, a foreign-deck version gate and a fuller inventory. Those are
	// deferred; until they land, a fleet-wide apply must not be reachable on --yes alone,
	// because --yes here means "rewrite the roster of every repository under this root".
	if yes && !dryRun && !confirmBreaking {
		fmt.Fprintln(stderr, "roster migrate: --yes rewrites the roster of EVERY deck under this root.\n"+
			"The full migration contract (compare-and-swap, per-deck confirmation, version gate) is not\n"+
			"implemented yet, so this operation is attended-only. Re-run with --confirm-breaking as well as --yes.")
		return 2
	}
	decks, err := findDecks(root)
	if err != nil {
		fmt.Fprintf(stderr, "roster migrate: %v\n", err)
		return 1
	}
	if len(decks) == 0 {
		fmt.Fprintf(stderr, "roster migrate: no parley-deck/COOPERATION.md found under %s\n", root)
		return 1
	}

	report := migrateReport{SchemaVersion: 1, Root: root, DryRun: dryRun || !yes}
	for _, deck := range decks {
		// A deck with uncommitted changes is migrated in place with git history as its
		// only rollback — and a dirty tree makes that rollback ambiguous. Skip and report
		// rather than write.
		if !dryRun && yes && deckTreeDirty(deck) {
			report.Decks = append(report.Decks, migrateDeck{
				Deck:        deck,
				Disposition: migrateSkipped,
				Reason:      "working tree has uncommitted changes; commit or stash first so the migration is reversible",
			})
			continue
		}
		report.Decks = append(report.Decks, migrateOneDeck(deck, backupDir, dryRun || !yes))
	}
	for _, d := range report.Decks {
		switch d.Disposition {
		case migrateApplied:
			report.Applied++
		case migrateSkipped:
			report.Skipped++
		case migrateUnchanged:
			report.Unchanged++
		case migrateFailed:
			report.Failed++
		}
	}

	if jsonOut {
		enc := json.NewEncoder(stdout)
		enc.SetIndent("", "  ")
		if err := enc.Encode(report); err != nil {
			fmt.Fprintf(stderr, "roster migrate: %v\n", err)
			return 1
		}
	} else {
		printMigrateReport(stdout, report)
	}
	if report.Failed > 0 {
		return 1
	}
	return 0
}

func printMigrateReport(w io.Writer, r migrateReport) {
	mode := "APPLY"
	if r.DryRun {
		mode = "DRY RUN — nothing was written"
	}
	fmt.Fprintf(w, "roster migrate (%s)\n  root: %s\n\n", mode, r.Root)
	for _, d := range r.Decks {
		fmt.Fprintf(w, "%-12s %s\n", string(d.Disposition), d.Deck)
		if d.Reason != "" {
			fmt.Fprintf(w, "             %s\n", d.Reason)
		}
		if len(d.Members) > 0 {
			fmt.Fprintf(w, "             members: %s\n", strings.Join(d.Members, ", "))
		}
		if len(d.Inactive) > 0 {
			fmt.Fprintf(w, "             retained inactive: %s\n", strings.Join(d.Inactive, ", "))
		}
	}
	fmt.Fprintf(w, "\napplied=%d skipped=%d unchanged=%d failed=%d\n",
		r.Applied, r.Skipped, r.Unchanged, r.Failed)
	if r.DryRun {
		fmt.Fprintln(w, "\nRe-run with --yes to apply. Every write is backed up first and verified after.")
	}
}

// findDecks enumerates workspace roots that contain parley-deck/COOPERATION.md.
func findDecks(root string) ([]string, error) {
	var out []string
	err := filepath.WalkDir(root, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return nil // an unreadable subtree is not fatal to the sweep
		}
		if !d.IsDir() {
			return nil
		}
		name := d.Name()
		// Never descend into build output, VCS internals or scratch copies: a deck found
		// there is a copy, not a project, and migrating it would be meaningless churn.
		if name == "node_modules" || name == ".git" || strings.HasPrefix(name, ".") {
			return filepath.SkipDir
		}
		if _, statErr := os.Stat(filepath.Join(path, "parley-deck", "COOPERATION.md")); statErr == nil {
			out = append(out, path)
		}
		return nil
	})
	sort.Strings(out)
	return out, err
}

func migrateOneDeck(deck, backupDir string, dryRun bool) migrateDeck {
	res := migrateDeck{Deck: deck}
	tomlPath := filepath.Join(deck, "parley-deck", "agents.toml")

	// Already migrated? A re-run must be a no-op, so a crash mid-sweep is resumable.
	existing, err := config.RosterEntriesInFile(tomlPath)
	if err != nil {
		res.Disposition, res.Reason = migrateSkipped, "unclean: "+tomlPath+" does not parse: "+err.Error()
		return res
	}
	if len(existing) > 0 {
		res.Disposition, res.Reason = migrateUnchanged, "already declares [roster.*]"
		for id := range existing {
			res.Members = append(res.Members, id)
		}
		sort.Strings(res.Members)
		return res
	}

	active, inactive, ok := protocol.ReadRosterIDs(deck)
	if !ok || len(active) == 0 {
		res.Disposition, res.Reason = migrateSkipped, "unclean: no readable §2 roster table to migrate"
		return res
	}

	ids := make([]string, 0, len(active))
	for id := range active {
		ids = append(ids, id)
	}
	sort.Strings(ids)

	// Resolve each ID to an adapter using the same proposal rule the resolver uses. An ID
	// we cannot map is the definition of unclean: writing a guessed adapter would create a
	// roster that launches the wrong CLI.
	specs, err := config.LoadAgentSpecs(deck)
	if err != nil {
		res.Disposition, res.Reason = migrateSkipped, "unclean: cannot load agent specs: "+err.Error()
		return res
	}
	byFamily := map[string]struct{}{}
	byFamilySpecs := map[string]any{}
	for _, s := range specs {
		byFamily[s.ID] = struct{}{}
		byFamilySpecs[s.ID] = s
	}
	mapping, _ := config.LoadRosterAdapters(deck)

	// Adapters still current on this machine, taken from the machine roster rather than a
	// hardcoded list. A deck row whose adapter is no longer rostered here (agy, gemini,
	// antigravity) is migrated as `active = false`.
	//
	// NOTE: this makes the migration change MEANING, not only storage — it retires agents
	// the §2 table still listed as active. That is a deliberate, user-authorized choice
	// (2026-08-06); the alternative was to carry rows across verbatim and retire them
	// afterwards with `roster set --state inactive`. Every such row is named in the report.
	current := map[string]bool{}
	if machineFile, err := rosterScopeFile(deck, "machine"); err == nil {
		if machineEntries, err := config.RosterEntriesInFile(machineFile); err == nil {
			for _, e := range machineEntries {
				if e.Adapter != "" {
					current[e.Adapter] = true
				}
			}
		}
	}

	type member struct {
		id, adapter string
		active      bool
	}
	var members []member
	var unresolved []string
	for _, id := range ids {
		adapter := mapping[id]
		if adapter == "" {
			adapter = proposeFamilyByName(id, byFamily)
		}
		if adapter == "" {
			unresolved = append(unresolved, id)
			continue
		}
		isActive := !inactive[id]
		if isActive && len(current) > 0 && !current[adapter] {
			isActive = false // adapter retired from the machine roster
		}
		members = append(members, member{id, adapter, isActive})
	}
	if len(unresolved) > 0 {
		res.Disposition = migrateSkipped
		res.Reason = "unclean: no adapter resolves for " + strings.Join(unresolved, ", ") +
			" — map them with `parley roster set <id> --adapter <family>` first"
		return res
	}

	var b strings.Builder
	b.WriteString("\n# Roster membership, migrated from the §2 table of COOPERATION.md on " +
		time.Now().UTC().Format("2006-01-02") + ".\n" +
		"# §2 is now a generated, non-authoritative view; this file is the deck's roster authority.\n" +
		"# Retired agents are kept with active = false rather than deleted, so past ideas stay readable.\n")
	for _, m := range members {
		b.WriteString(fmt.Sprintf("\n[roster.%s]\nadapter = %q\n", m.id, m.adapter))
		if !m.active {
			b.WriteString("active = false\n")
			res.Inactive = append(res.Inactive, m.id)
		}
		res.Members = append(res.Members, m.id)
	}

	before, _ := os.ReadFile(tomlPath)
	res.BeforeHash = shortHash(before)
	if dryRun {
		res.Disposition = migrateApplied // what WOULD happen
		return res
	}

	backup, err := backupFile(tomlPath, backupDir, deck)
	if err != nil {
		res.Disposition, res.Reason = migrateFailed, "backup failed: "+err.Error()
		return res
	}
	res.Backup = backup

	updated := strings.TrimRight(string(before), "\n") + "\n" + b.String()
	if strings.TrimSpace(string(before)) == "" {
		updated = strings.TrimLeft(b.String(), "\n")
	}
	if err := writeRosterFileAtomic(tomlPath, []byte(updated)); err != nil {
		res.Disposition, res.Reason = migrateFailed, "write failed: "+err.Error()
		return res
	}

	// Verify by re-reading through the real loader, and roll back if it does not come
	// back as expected. A migration that "succeeds" into an unparseable file is worse
	// than one that refuses.
	check, err := config.RosterEntriesInFile(tomlPath)
	if err != nil || len(check) != len(members) {
		_ = restoreFile(backup, tomlPath)
		res.Disposition = migrateFailed
		res.Reason = "post-write validation failed; deck restored from backup"
		return res
	}
	after, _ := os.ReadFile(tomlPath)
	res.AfterHash = shortHash(after)
	res.Disposition = migrateApplied
	return res
}

// proposeFamilyByName mirrors the resolver's proposal rule without needing spec values.
func proposeFamilyByName(id string, byFamily map[string]struct{}) string {
	if _, ok := byFamily[id]; ok {
		return id
	}
	stem := rosterInstanceSuffix.ReplaceAllString(id, "")
	if _, ok := byFamily[stem]; ok {
		return stem
	}
	if alias, ok := familyAliases[stem]; ok {
		if _, ok := byFamily[alias]; ok {
			return alias
		}
	}
	return ""
}

func backupFile(path, backupDir, deck string) (string, error) {
	rel := strings.ReplaceAll(strings.TrimPrefix(deck, "/"), "/", "__")
	dir := filepath.Join(backupDir, rel)
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return "", err
	}
	dst := filepath.Join(dir, "agents.toml")
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			// Record the absence explicitly so rollback can delete a file we created.
			return dst + ".absent", os.WriteFile(dst+".absent", []byte(""), 0o644)
		}
		return "", err
	}
	return dst, os.WriteFile(dst, data, 0o644)
}

func restoreFile(backup, path string) error {
	if strings.HasSuffix(backup, ".absent") {
		return os.Remove(path)
	}
	data, err := os.ReadFile(backup)
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0o644)
}

// shortHash is a stable content fingerprint used to record before/after state per deck.
func shortHash(b []byte) string {
	if len(b) == 0 {
		return "empty"
	}
	sum := sha256.Sum256(b)
	return hex.EncodeToString(sum[:])[:16]
}

// deckTreeDirty reports whether a deck's repository has uncommitted changes. A
// non-repository, or a git that will not answer, is treated as CLEAN: the guard exists to
// protect a reversible rollback path, and outside a repository there was never one to
// protect, so refusing there would block migration for no gain.
func deckTreeDirty(deck string) bool {
	cmd := exec.Command("git", "-C", deck, "status", "--porcelain")
	out, err := cmd.Output()
	if err != nil {
		return false
	}
	return len(strings.TrimSpace(string(out))) > 0
}
