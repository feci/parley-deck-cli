package runner

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"io"
	"os"
	"path/filepath"

	"parley-deck-cli/internal/agents"
	"parley-deck-cli/internal/budget"
)

type actionFileInput struct {
	Path    string
	Missing bool
	SHA256  string
}

func hashActionFile(path string) (actionFileInput, error) {
	row := actionFileInput{Path: path}
	st, e := os.Lstat(path)
	if os.IsNotExist(e) {
		row.Missing = true
		return row, nil
	}
	if e != nil {
		return row, e
	}
	if !st.Mode().IsRegular() || st.Size() > 4<<20 {
		return row, io.ErrUnexpectedEOF
	}
	f, e := os.Open(path)
	if e != nil {
		return row, e
	}
	defer f.Close()
	opened, e := f.Stat()
	if e != nil || !os.SameFile(st, opened) {
		return row, io.ErrUnexpectedEOF
	}
	b, e := io.ReadAll(io.LimitReader(f, (4<<20)+1))
	if e != nil {
		return row, e
	}
	after, e := f.Stat()
	current, nameErr := os.Lstat(path)
	if e != nil || nameErr != nil || len(b) > 4<<20 || !os.SameFile(opened, current) || int64(len(b)) != opened.Size() || after.Size() != opened.Size() || after.ModTime() != opened.ModTime() {
		return row, io.ErrUnexpectedEOF
	}
	sum := sha256.Sum256(b)
	row.SHA256 = hex.EncodeToString(sum[:])
	return row, nil
}

// Hash the declared runner request and protocol inputs. This is not a source-tree
// attestation or evidence that a model completed the requested operation.
func withRunnerActionInput(ctx context.Context, opts Options, operation string) context.Context {
	type agentInput struct {
		Spec          agents.Spec
		Path, Version string
	}
	agentsInput := make([]agentInput, 0, len(opts.Agents))
	for _, a := range opts.Agents {
		agentsInput = append(agentsInput, agentInput{a.Spec, a.Path, a.Version})
	}
	files := []actionFileInput{}
	for _, path := range []string{filepath.Join(opts.Root, "parley-deck", "COOPERATION.md"), filepath.Join(opts.Idea.Path, "00-prompt.md"), filepath.Join(opts.Idea.Path, "FINAL.md"), filepath.Join(opts.Idea.Path, "consensus.md"), filepath.Join(opts.Idea.Path, "review", "consensus.md"), filepath.Join(opts.Idea.Path, "IMPLEMENTATION.md")} {
		row, e := hashActionFile(path)
		if e != nil {
			return budget.InheritActionInput(ctx, budget.ActionInput{Operation: operation, Basis: "runtime-input", RunID: opts.RunID})
		}
		files = append(files, row)
	}
	recipe := struct {
		Operation, Root, RunID, Idea, Task, Phase, RoundLabel, Artifact string
		Round                                                           int
		Participants                                                    []string
		Agents                                                          []agentInput
		Mapping                                                         map[string]string
		TimeoutNS                                                       int64
		Overwrite, Strict                                               bool
		Files                                                           []actionFileInput
	}{operation, opts.Root, opts.RunID, opts.Idea.Slug, opts.Task, opts.Phase, opts.RoundLabel, opts.ArtifactName, opts.Round, append([]string(nil), opts.Idea.Participants...), agentsInput, opts.RosterMapping, int64(opts.Timeout), opts.Overwrite, opts.StrictGate, files}
	digest, _ := budget.ActionInputDigest(recipe) // empty on encoding/bounds error; reserve refuses
	return budget.InheritActionInput(ctx, budget.ActionInput{Operation: operation, Basis: "runtime-input", InputSHA256: digest, RunID: opts.RunID})
}

// Low-level/manual launches lack the complete parent recipe. Keep their basis
// explicit instead of treating launch metadata as a prompt/source attestation.
func withLaunchActionInput(ctx context.Context, info LaunchInfo, agent agents.Discovery) context.Context {
	recipe := struct {
		RunID, Idea, Phase, Artifact string
		Spec                         agents.Spec
	}{info.RunID, info.Idea, info.Phase, info.ArtifactPath, agent.Spec}
	digest, _ := budget.ActionInputDigest(recipe)
	return budget.InheritActionInput(ctx, budget.ActionInput{Operation: "manual-launch", Basis: "launch-metadata", InputSHA256: digest, RunID: info.RunID})
}
