package app

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"parley-deck-cli/internal/evidence"
	"parley-deck-cli/internal/protocol"
	"parley-deck-cli/internal/telemetry"
)

func verificationRefusalScope(root, idea string) (string, string, error) {
	if idea == "" || filepath.Base(idea) != idea || idea == "." || idea == ".." {
		return "", "", errors.New("invalid refusal idea scope")
	}
	root, err := filepath.Abs(root)
	if err != nil {
		return "", "", err
	}
	root, err = filepath.EvalSymlinks(root)
	if err != nil {
		return "", "", err
	}
	dir := filepath.Join(root, protocol.DeckDir, "ideas", idea)
	actual, err := filepath.EvalSymlinks(dir)
	if err != nil || actual != dir {
		return "", "", errors.New("refusal idea scope is missing or aliased")
	}
	return root, dir, nil
}

func newVerificationRefusal(observer, stage, idea, run, verifier, requestSHA, originalSHA, invocation string) (evidence.VerificationRefusal, error) {
	nonce, err := telemetry.NewID()
	if err != nil {
		return evidence.VerificationRefusal{}, err
	}
	id := sha256Hex(nonce)
	verificationID := id
	if evidence.RefusalHash(requestSHA) != nil {
		verificationID = requestSHA
	}
	return evidence.VerificationRefusal{Version: 1, ObservationID: id, VerificationID: verificationID, ObservedAt: time.Now().UTC(), Observer: observer, Stage: stage,
		Idea: telemetry.SafeLabel(idea), RunID: telemetry.SafeLabel(run), Verifier: telemetry.SafeLabel(verifier), InvocationID: telemetry.SafeLabel(invocation),
		RequestSHA256: evidence.RefusalHash(requestSHA), OriginalSHA256: evidence.RefusalHash(originalSHA)}, nil
}

func refusalGit(ctx context.Context, root string, args ...string) ([]byte, error) {
	cmd := exec.CommandContext(ctx, "git", append([]string{"--literal-pathspecs", "-C", root}, args...)...)
	cmd.Env = append(os.Environ(), "GIT_OPTIONAL_LOCKS=0")
	return cmd.Output()
}

func committedRefusal(ctx context.Context, root, path string, data []byte) (bool, error) {
	rel, err := filepath.Rel(root, path)
	if err != nil || filepath.IsAbs(rel) || rel == ".." || strings.HasPrefix(rel, ".."+string(filepath.Separator)) {
		return false, errors.New("refusal path escapes the worktree")
	}
	rel = filepath.ToSlash(rel)
	listed, err := refusalGit(ctx, root, "ls-tree", "-z", "HEAD", "--", rel)
	if err != nil {
		return false, errors.New("Git history is unavailable for refusal verification")
	}
	if len(listed) == 0 {
		return false, nil
	}
	committed, err := refusalGit(ctx, root, "show", "HEAD:"+rel)
	if err != nil || !bytes.Equal(committed, data) {
		return false, errors.New("committed refusal conflicts with the exact observation")
	}
	return true, nil
}

func commitVerificationRefusal(ctx context.Context, root, path string, data []byte) error {
	committed, err := committedRefusal(ctx, root, path, data)
	if err != nil || committed {
		return err
	}
	rel, err := filepath.Rel(root, path)
	if err != nil {
		return err
	}
	if _, err := refusalGit(ctx, root, "add", "--", rel); err != nil {
		return errors.New("could not stage the exact refusal record")
	}
	// A path-specific commit preserves unrelated staged and working-tree changes.
	if _, err := refusalGit(ctx, root, "commit", "--only", "-m", "[driver] independent verification refusal", "--", rel); err != nil {
		return errors.New("could not commit the exact refusal record; retained recovery is required")
	}
	committed, err = committedRefusal(ctx, root, path, data)
	if err != nil {
		return err
	}
	if !committed {
		return errors.New("refusal commit did not contain the exact record")
	}
	return nil
}

func recoverVerificationRefusal(ctx context.Context, root, dir, sum string) error {
	canonicalRoot, canonicalDir, err := verificationRefusalScope(root, filepath.Base(dir))
	if err != nil {
		return err
	}
	actual, err := filepath.EvalSymlinks(dir)
	if err != nil || actual != canonicalDir {
		return errors.New("refusal recovery scope differs from the requested idea")
	}
	root, dir = canonicalRoot, canonicalDir
	return evidence.PublishVerificationRefusal(ctx, dir, sum, func(path string, data []byte) error { return commitVerificationRefusal(ctx, root, path, data) })
}

func requireCommittedRefusals(ctx context.Context, root, dir string) error {
	entries, err := evidence.InspectVerificationRefusals(dir)
	if err != nil {
		return err
	}
	for _, entry := range entries {
		if entry.Problem != "" || entry.Record == nil {
			return errors.New("incomplete refusal observation requires recovery; completion remains refused")
		}
		data, _ := json.MarshalIndent(entry.Record, "", "  ")
		committed, err := committedRefusal(ctx, root, filepath.Join(dir, evidence.VerificationRefusalDirectory, entry.SHA256+".json"), data)
		if err != nil {
			return err
		}
		if !entry.Canonical || !committed {
			return errors.New("prior verification refusal requires exact-record recovery and fresh checks")
		}
	}
	return nil
}

// Called after the verification guard is released, even if the original context
// was canceled. Persistence has its own short deadline; it cannot launch a model.
func (o driverImplOps) retainDriverRefusal(stage, requestSHA, originalSHA, invocation string) string {
	root, dir, err := verificationRefusalScope(o.root, o.ideaSlug)
	if err != nil {
		return "refusal could not be retained: no canonical idea scope"
	}
	actual, err := filepath.EvalSymlinks(o.ideaDir)
	if err != nil || actual != dir {
		return "refusal could not be retained: idea directory mismatch"
	}
	rec, err := newVerificationRefusal("driver", stage, o.ideaSlug, o.base.RunID, o.drafter, requestSHA, originalSHA, invocation)
	if err != nil {
		return "refusal identity could not be allocated"
	}
	sum, err := evidence.RetainVerificationRefusal(dir, rec)
	if err != nil {
		return "refusal retention failed; completion remains refused"
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	entries, err := evidence.InspectVerificationRefusals(dir)
	if err != nil {
		return "refusal retained; canonical publication needs recovery"
	}
	// Publish this attempt's helper observation as well. Do not silently recover
	// older attempts, and never count two observations as two model invocations.
	for _, entry := range entries {
		if entry.Record == nil || entry.Record.VerificationID != rec.VerificationID {
			continue
		}
		if err := recoverVerificationRefusal(ctx, root, dir, entry.SHA256); err != nil {
			return "refusal retained; use evidence refusals inspect/recover (observation " + sum + ")"
		}
	}
	return "refusal committed; fresh checks are required"
}

func runEvidenceRefusals(ctx context.Context, args []string, stdout, stderr io.Writer) int {
	if len(args) == 0 || (args[0] != "inspect" && args[0] != "recover") {
		fmt.Fprintln(stderr, "usage: parley evidence refusals inspect|recover --dir DIR --idea IDEA [--expected-sha256 SHA]")
		return 2
	}
	flags := flag.NewFlagSet("evidence refusals "+args[0], flag.ContinueOnError)
	flags.SetOutput(stderr)
	rootFlag := flags.String("dir", ".", "repository root")
	idea := flags.String("idea", "", "exact idea slug")
	sum := flags.String("expected-sha256", "", "exact inspected observation digest for recovery")
	if err := flags.Parse(args[1:]); err != nil || flags.NArg() != 0 || *idea == "" {
		return 2
	}
	if (args[0] == "inspect" && *sum != "") || (args[0] == "recover" && evidence.RefusalHash(*sum) == nil) {
		fmt.Fprintln(stderr, "recovery requires an exact inspected digest; inspect accepts no recovery digest")
		return 2
	}
	root, dir, err := verificationRefusalScope(*rootFlag, *idea)
	if err != nil {
		fmt.Fprintln(stderr, err)
		return 1
	}
	if args[0] == "recover" {
		if err := recoverVerificationRefusal(ctx, root, dir, *sum); err != nil {
			fmt.Fprintln(stderr, err)
			return 1
		}
		fmt.Fprintln(stdout, "Exact refusal record verified committed. No checks, model invocation or acceptance were performed; rerun checks before verification.")
		return 0
	}
	entries, err := evidence.InspectVerificationRefusals(dir)
	if err != nil {
		fmt.Fprintln(stderr, err)
		return 1
	}
	type inspected struct {
		evidence.RefusalEntry
		Committed     bool   `json:"committed"`
		CommitProblem string `json:"commit_problem"`
	}
	result := make([]inspected, 0, len(entries))
	for _, entry := range entries {
		row := inspected{RefusalEntry: entry}
		if entry.Record != nil {
			data, _ := json.MarshalIndent(entry.Record, "", "  ")
			row.Committed, err = committedRefusal(ctx, root, filepath.Join(dir, evidence.VerificationRefusalDirectory, entry.SHA256+".json"), data)
			if err != nil {
				row.CommitProblem = err.Error()
			}
		}
		result = append(result, row)
	}
	encoder := json.NewEncoder(stdout)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(result); err != nil {
		return 1
	}
	return 0
}
