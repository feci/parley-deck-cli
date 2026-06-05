package app

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"parley-deck-cli/internal/agents"
	"parley-deck-cli/internal/consensus"
	"parley-deck-cli/internal/driver"
)

// driverConsensusOps is the production driver.ConsensusOps adapter (slice 2.2). It
// lives in internal/app so it can reuse the existing request-signoffs path and the
// headless agent invoker; the driver depends only on the driver.ConsensusOps
// interface, so internal/driver never imports internal/app (the import-direction
// guarantee D9 sought, achieved by injection instead of extraction).
type driverConsensusOps struct {
	root       string
	ideaSlug   string
	ideaDir    string
	discovered []agents.Discovery
	out        io.Writer
}

func newDriverConsensusOps(root, ideaSlug, ideaDir string, discovered []agents.Discovery, out io.Writer) driver.ConsensusOps {
	return driverConsensusOps{root: root, ideaSlug: ideaSlug, ideaDir: ideaDir, discovered: discovered, out: out}
}

func (o driverConsensusOps) Status() (consensus.Summary, error) {
	return consensus.Status(o.root, o.ideaSlug, false)
}

// Draft creates the consensus.md scaffold (+ sets idea status=consensus) if absent,
// then invokes a drafter agent to author the real synthesis into consensus.md.
func (o driverConsensusOps) Draft(ctx context.Context) error {
	path := filepath.Join(o.ideaDir, "consensus.md")
	if _, err := os.Stat(path); errors.Is(err, os.ErrNotExist) {
		if _, derr := consensus.Draft(o.root, o.ideaSlug, consensus.DraftOptions{}); derr != nil {
			return derr
		}
	} else if err != nil {
		return err
	}
	return o.runDrafter(ctx, "consensus", buildConsensusDraftPrompt(o.ideaDir, path))
}

// RequestSignoffs invokes each missing participant through the existing
// request-signoffs path so every agent authors its own signoff.
func (o driverConsensusOps) RequestSignoffs(ctx context.Context, missing []string) error {
	return requestConsensusSignoffs(ctx, requestSignoffsOptions{
		Root:            o.root,
		IdeaSlug:        o.ideaSlug,
		ParticipantsRaw: strings.Join(missing, ","),
		Yes:             true,
	}, o.out, o.out)
}

// DraftFinal writes the FINAL.md scaffold (+ status=final) via consensus.Finalize
// if absent, then invokes a drafter agent to author the real FINAL content.
func (o driverConsensusOps) DraftFinal(ctx context.Context) error {
	path := filepath.Join(o.ideaDir, "FINAL.md")
	if _, err := os.Stat(path); errors.Is(err, os.ErrNotExist) {
		if _, _, ferr := consensus.Finalize(o.root, o.ideaSlug, consensus.FinalizeOptions{By: "driver"}); ferr != nil {
			return ferr
		}
	} else if err != nil {
		return err
	}
	return o.runDrafter(ctx, "FINAL", buildFinalDraftPrompt(o.ideaDir, path))
}

func (o driverConsensusOps) Reopen(ctx context.Context, reason string) error {
	_, err := consensus.Reopen(o.root, o.ideaSlug, consensus.ReopenOptions{Reason: reason})
	return err
}

// runDrafter invokes the first available headless agent to author the target file
// per the given prompt. Drafting is a single-agent facilitator action (D6).
func (o driverConsensusOps) runDrafter(ctx context.Context, kind, prompt string) error {
	drafter, ok := firstHeadlessAgent(o.discovered)
	if !ok {
		return fmt.Errorf("no headless agent available to draft %s", kind)
	}
	rootAbs, err := filepath.Abs(o.root)
	if err != nil {
		return err
	}
	fmt.Fprintf(o.out, "driver: drafting %s via %s ...\n", kind, drafter.ID)
	return runHeadlessSignoffAgent(ctx, rootAbs, drafter, prompt, o.out, o.out)
}

func firstHeadlessAgent(discovered []agents.Discovery) (agents.Discovery, bool) {
	for _, agent := range discovered {
		if agent.Found && agent.ID != "gemini" && agents.LaunchModeOrDefault(agent.LaunchMode) == agents.LaunchHeadless {
			return agent, true
		}
	}
	return agents.Discovery{}, false
}

func buildConsensusDraftPrompt(ideaDir, path string) string {
	return fmt.Sprintf(`You are the Parley Deck facilitator drafting the consensus.

Read EVERY round artifact under %s/round-*/ (each participant's round files).
The file %s currently holds a scaffold. OVERWRITE it with the real synthesis of the
deliberation, keeping this exact structure and the YAML frontmatter:

## Agreed decisions
(concrete decisions the rounds converged on)
## Trade-offs accepted
## Deferred follow-ups
## Dismissed findings
## Signoffs

Under "## Signoffs", leave one HTML comment placeholder line per participant exactly
like "<!-- <agent-id> appends its signoff below -->" for every participant in the
idea's 00-prompt.md, and NOTHING else under that heading (the agents append their own
✅/🟡/❌ blocks later). Be concrete and concise. English only. Write the file now and
report only the path.`, ideaDir, path)
}

func buildFinalDraftPrompt(ideaDir, path string) string {
	return fmt.Sprintf(`You are the Parley Deck facilitator drafting FINAL.md.

Read %s/consensus.md (the accepted consensus + signoffs) and the round artifacts.
OVERWRITE %s with the final specification, keeping YAML frontmatter that includes
"status: final", and a populated section:

## Final plan / specification

with at least several concrete lines describing the agreed design (no placeholders,
no unexpanded <...> tokens). Be concrete. English only. Write the file now and report
only the path.`, ideaDir, path)
}
