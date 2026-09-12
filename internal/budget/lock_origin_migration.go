package budget

import (
	"bytes"
	"context"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"reflect"
	"strings"

	"parley-deck-cli/internal/fsutil"
)

const originMigrationDirectory = "lock-origin-migrations"
const originMigrationVersion = "parley-budget-lock/v3"

type pinnedOrigin struct{ Host, Path, Token, Migration string }

func parsePinnedOrigin(raw []byte) (pinnedOrigin, error) {
	p := strings.Split(string(raw), "\n")
	if len(p) != 5 && len(p) != 6 || p[len(p)-1] != "" {
		return pinnedOrigin{}, errors.New("malformed budget lock origin")
	}
	o := pinnedOrigin{Host: p[1], Path: p[2], Token: p[3]}
	if len(p) == 5 && p[0] != "parley-budget-lock/v2" || len(p) == 6 && p[0] != originMigrationVersion {
		return o, errors.New("unsupported budget lock origin version")
	}
	if len(p) == 6 {
		o.Migration = p[4]
		if !originHash(o.Migration) {
			return o, errors.New("invalid migration reference")
		}
	}
	host, err := os.Hostname()
	if err != nil || host == "" || o.Host != host {
		return o, errors.New("budget lock hostname changed; cross-host migration is unsupported")
	}
	if !filepath.IsAbs(o.Path) || filepath.Clean(o.Path) != o.Path || strings.ContainsAny(o.Path, "\r\x00") || !originHash(o.Token) {
		return o, errors.New("invalid budget lock path or token")
	}
	return o, nil
}
func originBindsResource(dir string, o pinnedOrigin) error {
	if filepath.Base(o.Path) != key(strings.ToLower(dir))+".lock" {
		return errors.New("origin belongs to a different resource path; cache relocation cannot migrate repository or accounting scope")
	}
	return nil
}

func originHash(s string) bool {
	b, e := hex.DecodeString(s)
	return e == nil && len(b) == 32 && strings.ToLower(s) == s
}
func canonicalOriginDirectory(dir string) (string, error) {
	p, e := filepath.Abs(dir)
	if e != nil {
		return "", e
	}
	p, e = filepath.EvalSymlinks(p)
	if e != nil {
		return "", e
	}
	st, e := os.Stat(p)
	if e != nil {
		return "", e
	}
	if !st.IsDir() {
		return "", errors.New("origin resource must be an existing directory")
	}
	return p, nil
}

// Only hashes and file metadata are retained; protected file bodies never enter
// an operator record. Inventory is bounded and excludes only this migration's
// authority files. Descendant locks/scopes are not migrated by a parent move.
type OriginMaterial struct {
	Path   string `json:"path"`
	Mode   uint32 `json:"mode"`
	SHA256 string `json:"sha256"`
}
type LockOriginPreview struct {
	Version          int              `json:"version"`
	Dir              string           `json:"dir"`
	Origin           string           `json:"origin"`
	GuardEstablished bool             `json:"guard_established"`
	TargetPath       string           `json:"target_path"`
	TargetToken      string           `json:"target_token"`
	Material         []OriginMaterial `json:"material"`
	SHA256           string           `json:"sha256"`
}

func (p LockOriginPreview) digest() string { p.SHA256 = ""; return migrationDigest(p) }

type LockOriginRequest struct {
	ExpectedSHA256 string `json:"expected_sha256"`
	DecisionID     string `json:"decision_id"`
	Reason         string `json:"reason"`
}
type originMigration struct {
	Version     int               `json:"version"`
	Request     LockOriginRequest `json:"request"`
	Preview     LockOriginPreview `json:"preview"`
	TargetToken string            `json:"target_token"`
}

func (m originMigration) reference() string { return key(m.Request.DecisionID) }
func (m originMigration) origin() []byte {
	o, _ := parsePinnedOrigin([]byte(m.Preview.Origin))
	return []byte(originMigrationVersion + "\n" + o.Host + "\n" + m.Preview.TargetPath + "\n" + m.TargetToken + "\n" + m.reference() + "\n")
}
func originRecordPath(dir, reference string) string {
	return filepath.Join(dir, originMigrationDirectory, reference+".json")
}
func originCompletePath(dir, reference string) string {
	return filepath.Join(dir, originMigrationDirectory, reference+".complete")
}

func originMaterial(ctx context.Context, dir string) ([]OriginMaterial, error) {
	var total int64
	result := []OriginMaterial{}
	err := filepath.WalkDir(dir, func(path string, entry fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if err = ctx.Err(); err != nil {
			return err
		}
		if path == dir {
			return nil
		}
		rel, err := filepath.Rel(dir, path)
		if err != nil {
			return err
		}
		if rel == originMigrationDirectory {
			if !entry.IsDir() {
				return errors.New("migration journal directory is aliased")
			}
			return filepath.SkipDir
		}
		if rel == "lock-origin" || rel == resourceGuardWitness {
			return nil
		}
		if len(result) >= 10000 {
			return errors.New("origin migration inventory exceeds 10000 entries")
		}
		info, err := entry.Info()
		if err != nil {
			return err
		}
		row := OriginMaterial{Path: filepath.ToSlash(rel), Mode: uint32(info.Mode())}
		if !entry.IsDir() {
			if !info.Mode().IsRegular() {
				return errors.New("origin migration refuses aliased or special protected files")
			}
			if info.Size() > 16<<20 || info.Size() < 0 || total+info.Size() > 64<<20 {
				return errors.New("origin migration inventory exceeds bounded 16 MiB member / 64 MiB total")
			}
			raw, err := readStepHistoryFile(path, 16<<20)
			if err != nil {
				return err
			}
			total += int64(len(raw))
			row.SHA256 = key(string(raw))
		}
		result = append(result, row)
		return nil
	})
	return result, err
}
func originWitness(dir string, from, to []byte) (bool, error) {
	b, e := readLockOrigin(filepath.Join(dir, resourceGuardWitness))
	if os.IsNotExist(e) {
		return false, nil
	}
	if e != nil {
		return false, e
	}
	if !bytes.Equal(b, from) && (to == nil || !bytes.Equal(b, to)) {
		return false, errors.New("resource guard witness differs from migration authority")
	}
	return true, nil
}
func originTarget(dir, targetDirectory string) (string, error) {
	target, e := canonicalOriginDirectory(targetDirectory)
	if e != nil {
		return "", e
	}
	// Targeting the protected tree would make the migration modify its own input.
	rel, e := filepath.Rel(dir, target)
	if e != nil {
		return "", e
	}
	if rel == "." || rel != ".." && !strings.HasPrefix(rel, ".."+string(filepath.Separator)) {
		return "", errors.New("destination lock directory must be outside the protected resource")
	}
	return filepath.Join(target, key(strings.ToLower(dir))+".lock"), nil
}
func originTargetIdentity(path string) (string, error) {
	if _, e := os.Lstat(path); os.IsNotExist(e) {
		return "", nil
	} else if e != nil {
		return "", e
	}
	return lockIdentity(path, false)
}
func inspectOriginLocked(ctx context.Context, dir, target string, raw []byte) (LockOriginPreview, error) {
	p := LockOriginPreview{Version: 1, Dir: dir, Origin: string(raw), TargetPath: target}
	var e error
	p.GuardEstablished, e = originWitness(dir, raw, nil)
	if e != nil {
		return p, e
	}
	p.Material, e = originMaterial(ctx, dir)
	if e != nil {
		return p, e
	}
	p.TargetToken, e = originTargetIdentity(target)
	if os.IsNotExist(e) {
		p.TargetToken = ""
		e = nil
	}
	if e != nil {
		return p, e
	}
	if e = verifyExactLockOrigin(dir, raw); e != nil {
		return p, e
	}
	p.SHA256 = p.digest()
	return p, nil
}

// Preview is read-only and takes the existing kernel lock. It cannot initialize
// a missing lock, charge work, or repair missing accounting. A later apply binds
// the exact result and repeats its inventory after waiting for exclusion.
func PreviewLockOriginMigration(ctx context.Context, dir, targetDirectory string) (LockOriginPreview, error) {
	var zero LockOriginPreview
	dir, e := canonicalOriginDirectory(dir)
	if e != nil {
		return zero, e
	}
	target, e := originTarget(dir, targetDirectory)
	if e != nil {
		return zero, e
	}
	raw, e := readLockOrigin(filepath.Join(dir, "lock-origin"))
	if e != nil {
		return zero, e
	}
	o, e := parsePinnedOrigin(raw)
	if e != nil {
		return zero, e
	}
	if e = originBindsResource(dir, o); e != nil {
		return zero, e
	}
	if o.Path == target {
		return zero, errors.New("destination equals current lock path")
	}
	check := func() error {
		if e := verifyExactLockOrigin(dir, raw); e != nil {
			return e
		}
		if o.Migration != "" {
			_, e := resolveOriginChainLimit(dir, raw, true, 127)
			return e
		}
		return nil
	}
	release, e := acquirePinnedKernelLock(ctx, o.Path, o.Token, tryLock, unlock, check, nil)
	if e != nil {
		return zero, e
	}
	defer release()
	return inspectOriginLocked(ctx, dir, target, raw)
}
func validOriginRequest(r LockOriginRequest) bool {
	return originHash(r.ExpectedSHA256) && strings.TrimSpace(r.DecisionID) == r.DecisionID && r.DecisionID != "" && len(r.DecisionID) <= 200 && strings.TrimSpace(r.Reason) != "" && len(r.Reason) <= 2000 && !strings.ContainsAny(r.DecisionID+r.Reason, "\x00\r\n")
}
func readOriginMigration(dir, reference string) (originMigration, error) {
	var m originMigration
	if !originHash(reference) {
		return m, errors.New("invalid origin migration reference")
	}
	if e := migrationDirectory(dir, originMigrationDirectory); e != nil {
		return m, e
	}
	raw, e := readStepHistoryFile(originRecordPath(dir, reference), 4<<20)
	if e != nil {
		return m, e
	}
	if e = migrationJSON(raw, &m); e != nil {
		return m, e
	}
	if m.Version != 1 || m.Preview.Version != 1 || m.Preview.Dir != dir || m.reference() != reference || !validOriginRequest(m.Request) || m.Request.ExpectedSHA256 != m.Preview.SHA256 || m.Preview.SHA256 != m.Preview.digest() || !originHash(m.TargetToken) {
		return m, errors.New("origin migration record binding mismatch")
	}
	o, e := parsePinnedOrigin([]byte(m.Preview.Origin))
	if e != nil {
		return m, e
	}
	if e = originBindsResource(dir, o); e != nil {
		return m, e
	}
	target, e := originTarget(dir, filepath.Dir(m.Preview.TargetPath))
	if e != nil {
		return m, e
	}
	if target != m.Preview.TargetPath || target == o.Path || m.Preview.TargetToken != "" && m.Preview.TargetToken != m.TargetToken {
		return m, errors.New("origin migration target binding mismatch")
	}
	return m, nil
}
func originMigrationComplete(dir string, m originMigration) (bool, error) {
	b, e := readLockOrigin(originCompletePath(dir, m.reference()))
	if os.IsNotExist(e) {
		return false, nil
	}
	if e != nil {
		return false, e
	}
	if string(b) != migrationDigest(m)+"\n" {
		return false, errors.New("origin migration completion differs from retained decision")
	}
	return true, nil
}

// Validate the entire completed migration chain, never an arbitrary v3 path.
// Historical entries certify cutover, not today's ledger state or human identity.
func resolveMigratedOrigin(dir string, raw []byte) (pinnedOrigin, error) {
	return resolveOriginChain(dir, raw, true)
}
func resolveOriginChain(dir string, raw []byte, checkWitness bool) (pinnedOrigin, error) {
	return resolveOriginChainLimit(dir, raw, checkWitness, 128)
}
func resolveOriginChainLimit(dir string, raw []byte, checkWitness bool, limit int) (pinnedOrigin, error) {
	first, e := parsePinnedOrigin(raw)
	if e != nil {
		return first, e
	}
	current := raw
	seen := map[string]bool{}
	for depth := 0; depth < limit; depth++ {
		o, e := parsePinnedOrigin(current)
		if e != nil {
			return first, e
		}
		if e = originBindsResource(dir, o); e != nil {
			return first, e
		}
		if o.Migration == "" {
			return first, nil
		}
		if seen[o.Migration] {
			return first, errors.New("cyclic origin migration chain")
		}
		seen[o.Migration] = true
		m, e := readOriginMigration(dir, o.Migration)
		if e != nil {
			return first, e
		}
		done, e := originMigrationComplete(dir, m)
		if e != nil {
			return first, e
		}
		if !done {
			return first, errors.New("origin migration is incomplete; replay its exact attended apply decision before work")
		}
		if !bytes.Equal(current, m.origin()) {
			return first, errors.New("origin differs from completed migration")
		}
		if depth == 0 && checkWitness {
			witness, e := originWitness(dir, current, nil)
			if e != nil {
				return first, e
			}
			if m.Preview.GuardEstablished && !witness {
				return first, errors.New("migrated resource guard witness is missing")
			}
		}
		current = []byte(m.Preview.Origin)
	}
	return first, errors.New("origin migration chain has no capacity for this operation; preserve existing authority")
}
func publishOriginExclusive(path string, data []byte) error {
	f, e := os.CreateTemp(filepath.Dir(path), ".origin-migration-*")
	if e != nil {
		return e
	}
	defer os.Remove(f.Name())
	defer f.Close()
	if _, e = f.Write(data); e != nil {
		return e
	}
	if e = fsutil.SyncFile(f); e != nil {
		return e
	}
	if e = f.Close(); e != nil {
		return e
	}
	if e = publishExclusive(f.Name(), path); e != nil && !os.IsExist(e) {
		return e
	}
	b, e := readStepHistoryFile(path, 4<<20)
	if e != nil {
		return e
	}
	if !bytes.Equal(b, data) {
		return errors.New("origin migration publication conflicts with retained bytes")
	}
	return nil
}

// ApplyLockOriginMigration requires the platform's attended operator entrypoint.
// Its library request is decision metadata, not authentication. Exact replay
// returns the retained preview even after later writes/migrations. It never
// resets ledger bytes, reservations, ceilings or the accounting epoch.
func ApplyLockOriginMigration(ctx context.Context, dir, targetDirectory string, r LockOriginRequest) (LockOriginPreview, error) {
	return applyLockOriginMigration(ctx, dir, targetDirectory, r, tryLock, func(string) error { return nil })
}

// after is a publication fault seam for deterministic crash-boundary fixtures;
// takeTarget permits testing a filesystem that falsely reports exclusion.
func applyLockOriginMigration(ctx context.Context, dir, targetDirectory string, r LockOriginRequest, takeTarget func(*os.File) (bool, error), after func(string) error) (LockOriginPreview, error) {
	var zero LockOriginPreview
	if !validOriginRequest(r) {
		return zero, errors.New("migration requires exact preview hash, decision identity and reason")
	}
	dir, e := canonicalOriginDirectory(dir)
	if e != nil {
		return zero, e
	}
	target, e := originTarget(dir, targetDirectory)
	if e != nil {
		return zero, e
	}
	reference := key(r.DecisionID)
	m, e := readOriginMigration(dir, reference)
	existing := e == nil
	if e != nil && !os.IsNotExist(e) {
		return zero, e
	}
	if existing {
		if m.Request != r || m.Preview.TargetPath != target {
			return zero, errors.New("decision identity already belongs to a different origin migration")
		}
		done, e := originMigrationComplete(dir, m)
		if e != nil {
			return zero, e
		}
		if done {
			return m.Preview, nil
		}
	} else {
		p, e := PreviewLockOriginMigration(ctx, dir, targetDirectory)
		if e != nil {
			return zero, e
		}
		if p.SHA256 != r.ExpectedSHA256 {
			return zero, errors.New("stale origin migration preview; inspect again and use a new decision")
		}
		m = originMigration{Version: 1, Request: r, Preview: p, TargetToken: p.TargetToken}
		if m.TargetToken == "" {
			var nonce [32]byte
			if _, e = rand.Read(nonce[:]); e != nil {
				return zero, e
			}
			m.TargetToken = hex.EncodeToString(nonce[:])
		}
	}
	from := []byte(m.Preview.Origin)
	to := m.origin()
	o, e := parsePinnedOrigin(from)
	if e != nil {
		return zero, e
	}
	check := func() error {
		raw, e := readLockOrigin(filepath.Join(dir, "lock-origin"))
		if e != nil {
			return e
		}
		if !bytes.Equal(raw, from) && !bytes.Equal(raw, to) {
			return errors.New("migration no longer owns this resource origin")
		}
		// An ancestor's own witness was already superseded at this migration's
		// cutover; verify its immutable chain without comparing the current witness.
		if o.Migration != "" {
			_, e = resolveOriginChainLimit(dir, from, false, 127)
			return e
		}
		return nil
	}
	var oldFile *os.File
	takeOld := func(f *os.File) (bool, error) {
		if oldFile == nil {
			oldFile = f
		}
		return tryLock(f)
	}
	release, e := acquirePinnedKernelLock(ctx, o.Path, o.Token, takeOld, unlock, check, nil)
	if e != nil {
		return zero, e
	}
	defer release()
	// A competing exact invocation may have completed while this one waited.
	if stored, err := readOriginMigration(dir, reference); err == nil {
		if stored.Request != r || stored.Preview.TargetPath != target {
			return zero, errors.New("conflicting migration won publication")
		}
		m = stored
		to = m.origin()
		if done, err := originMigrationComplete(dir, m); err != nil {
			return zero, err
		} else if done {
			return m.Preview, nil
		}
		existing = true
	} else if !os.IsNotExist(err) {
		return zero, err
	}
	materialCheck := func() error {
		if e := check(); e != nil {
			return e
		}
		present, e := originWitness(dir, from, to)
		if e != nil {
			return e
		}
		if present != m.Preview.GuardEstablished {
			return errors.New("guard witness presence changed since migration preview")
		}
		current, e := originMaterial(ctx, dir)
		if e != nil {
			return e
		}
		if !reflect.DeepEqual(current, m.Preview.Material) {
			return errors.New("protected material changed since migration preview; preserve the decision; a pre-cutover attempt may be superseded by a fresh preview")
		}
		return nil
	}
	if e = materialCheck(); e != nil {
		return zero, e
	}
	if !existing {
		// Recheck destination state after the old-lock wait, before publishing any
		// journal or identity. Destination acquisition below independently checks it.
		token, err := originTargetIdentity(target)
		if os.IsNotExist(err) {
			token = ""
			err = nil
		}
		if err != nil {
			return zero, err
		}
		if token != m.Preview.TargetToken {
			return zero, errors.New("destination identity changed since migration preview")
		}
		journalDir := filepath.Join(dir, originMigrationDirectory)
		if e = fsutil.MkdirAllResilient(journalDir, 0700); e != nil {
			return zero, e
		}
		if e = migrationDirectory(dir, originMigrationDirectory); e != nil {
			return zero, e
		}
		raw, e := json.Marshal(m)
		if e != nil {
			return zero, e
		}
		if e = publishOriginExclusive(originRecordPath(dir, reference), raw); e != nil {
			return zero, e
		}
		if e = after("journal"); e != nil {
			return zero, e
		}
	}
	if e = publishOriginIdentity(target, m.TargetToken); e != nil {
		return zero, e
	}
	if e = after("identity"); e != nil {
		return zero, e
	}
	// Do not wait holding one lock for an arbitrary second resource. Contention
	// leaves the immutable journal available for explicit exact retry.
	take := func(f *os.File) (bool, error) {
		held, e := takeTarget(f)
		if e == nil && !held {
			return false, fmt.Errorf("%w: destination lock is busy; exact migration remains resumable", ErrLockContention)
		}
		return held, e
	}
	// The probe MUST observe contention as a normal false, not as a take error.
	calls := 0
	var targetFile *os.File
	targetTake := func(f *os.File) (bool, error) {
		calls++
		if calls == 1 {
			targetFile = f
			return take(f)
		}
		return takeTarget(f)
	}
	releaseTarget, e := acquirePinnedKernelLock(ctx, target, m.TargetToken, targetTake, unlock, check, nil)
	if e != nil {
		return zero, e
	}
	defer releaseTarget()
	heldCheck := func() error {
		if e := verifyLockIdentity(oldFile, o.Path, o.Token); e != nil {
			return e
		}
		if e := verifyLockIdentity(targetFile, target, m.TargetToken); e != nil {
			return e
		}
		return materialCheck()
	}
	if e = heldCheck(); e != nil {
		return zero, e
	}
	if e = after("target-held"); e != nil {
		return zero, e
	}
	if e = ctx.Err(); e != nil {
		return zero, e
	}
	if e = heldCheck(); e != nil {
		return zero, e
	}
	if e = writeSynced(filepath.Join(dir, "lock-origin"), to); e != nil {
		return zero, e
	}
	if e = after("origin"); e != nil {
		return zero, e
	}
	if e = heldCheck(); e != nil {
		return zero, e
	}
	if m.Preview.GuardEstablished {
		if e = writeSynced(filepath.Join(dir, resourceGuardWitness), to); e != nil {
			return zero, e
		}
	}
	if e = after("witness"); e != nil {
		return zero, e
	}
	if e = heldCheck(); e != nil {
		return zero, e
	}
	if e = publishOriginExclusive(originCompletePath(dir, reference), []byte(migrationDigest(m)+"\n")); e != nil {
		return zero, e
	}
	if e = after("complete"); e != nil {
		return zero, e
	}
	if e = heldCheck(); e != nil {
		return zero, e
	}
	if _, e = resolveMigratedOrigin(dir, to); e != nil {
		return zero, e
	}
	return m.Preview, nil
}
func publishOriginIdentity(path, token string) error {
	prior, e := originTargetIdentity(path)
	if e != nil {
		return e
	}
	if prior != "" {
		if prior != token {
			return errors.New("destination identity differs from retained migration")
		}
		return nil
	}
	return publishOriginExclusive(path, []byte(token+"\n"))
}
