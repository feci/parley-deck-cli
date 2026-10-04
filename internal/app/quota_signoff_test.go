package app

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"parley-deck-cli/internal/consensus"
	"parley-deck-cli/internal/protocol"
	"parley-deck-cli/internal/quota"
	"parley-deck-cli/internal/runmanifest"
	"parley-deck-cli/internal/store"
)

func TestQuotaSignoffProcessSettlesBatchAndKeepsCurrentRequiredSet(t *testing.T) {
	root := t.TempDir()
	t.Setenv("PARLEY_HOME", t.TempDir())
	if err := protocol.InitWorkspace(root); err != nil {
		t.Fatal(err)
	}
	declareAppTestSource(t, root)
	p := quota.NewPolicy(nil, nil)
	ids := []string{"alpha", "beta", "zcode"}
	idea, k, err := protocol.CreateIdeaWithQuota(root, "signoff quota", ids, nil, "deliberation", "", "kickoff-run", &p, nil)
	if err != nil {
		t.Fatal(err)
	}
	prompt, _ := os.ReadFile(filepath.Join(idea.Path, "00-prompt.md"))
	writeConsensusIdea(t, root, idea.Slug, ids, false, nil)
	os.WriteFile(filepath.Join(idea.Path, "00-prompt.md"), prompt, 0600)
	path := filepath.Join(idea.Path, "consensus.md")
	raw, _ := os.ReadFile(path)
	raw = bytes.Replace(raw, []byte("---\n"), []byte("---\ndrafted-by: alpha\n"), 1)
	os.WriteFile(path, raw, 0600)
	m := runmanifest.New(runmanifest.Options{Root: root, RunID: k.RunID, IdeaSlug: idea.Slug, Participants: ids, QuotaKickoff: k})
	if err = runmanifest.Write(root, k.RunID, m); err != nil {
		t.Fatal(err)
	}
	store.New(filepath.Join(root, protocol.DeckDir, "runs", k.RunID)).AppendDurable(store.Event{Type: "run.created", Data: map[string]any{"idea": idea.Slug, "participants": ids, "quota_kickoff": k}})
	bin := t.TempDir()
	alpha := writeFakeSignoffCLI(t, bin, "alpha", "accept", 0)
	beta := writeFakeSignoffCLI(t, bin, "beta", "accept", 0)
	failed := filepath.Join(bin, "zcode")
	script := fmt.Sprintf("#!/bin/sh\ncat >/dev/null\ncat >&2 <<'QUOTA'\nAPICallError [AI_APICallError]: Weekly Limit Exhausted\n    at fixture (stub.js:1:1) {\n  cause: undefined,\n  url: 'https://provider.invalid',\n  requestBodyValues: undefined,\n  statusCode: 429,\n  responseHeaders: {},\n  responseBody: '{\"error\":{\"message\":\"Weekly Limit Exhausted\",\"reset_at\":\"%s\"}}',\n  isRetryable: true,\n  data: undefined,\n  Symbol(vercel.ai.error): true,\n  Symbol(vercel.ai.error.AI_APICallError): true\n}\nError: Turn execution failed\nQUOTA\nexit 1\n", time.Now().UTC().Add(3*time.Hour).Format(time.RFC3339Nano))
	if err = os.WriteFile(failed, []byte(script), 0700); err != nil {
		t.Fatal(err)
	}
	writeAgentsLocalConfig(t, root, fakeAgentConfig{ID: "alpha", Path: alpha}, fakeAgentConfig{ID: "beta", Path: beta}, fakeAgentConfig{ID: "zcode", Path: failed})
	var out, errs bytes.Buffer
	if err = requestConsensusSignoffs(context.Background(), requestSignoffsOptions{Root: root, IdeaSlug: idea.Slug, Yes: true}, &out, &errs); err != nil {
		t.Fatal(err, out.String(), errs.String())
	}
	h, err := quota.ReadHistory(idea.Path)
	if err != nil || h.Revision != 1 || strings.Join(h.Current, ",") != "alpha,beta" {
		t.Fatal(h, err, out.String(), errs.String())
	}
	summary, err := consensus.Status(root, idea.Slug, false)
	if err != nil || summary.Triage != consensus.TriageReady || len(summary.Signoffs) != 2 {
		t.Fatal(summary, err)
	}
	notes, _ := filepath.Glob(filepath.Join(root, protocol.DeckDir, "inbox", "parley-to-user_batch-*.md"))
	if len(notes) != 1 {
		t.Fatal(notes)
	}
}
