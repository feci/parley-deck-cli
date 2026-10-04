package quota_test

import (
	"bytes"
	"crypto/sha256"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"parley-deck-cli/internal/quota"
	"parley-deck-cli/internal/quotatest"
)

func TestQuotaFixupAuthorityCommittedAttributionAndObjects(t *testing.T) {
	root := t.TempDir()
	idea := "authority-fixture"
	quote := "Owner ruling on obligation-fixture."
	a := quotatest.Authority(t, root, idea, "good", quote)
	for _, kind := range []string{"tree-not-commit", "committed-self-author", "committed-wrong-idea", "fabricated-path", "unquoted-content", "changed-archive"} {
		t.Run(kind, func(t *testing.T) {
			bad := a
			switch kind {
			case "tree-not-commit":
				bad.Commit = quotatest.Git(t, root, "rev-parse", a.Commit+"^{tree}")
			case "committed-self-author", "committed-wrong-idea":
				data, _ := os.ReadFile(filepath.Join(root, a.Path))
				if kind == "committed-self-author" {
					data = bytes.Replace(data, []byte("from: user"), []byte("from: codex-1"), 1)
				} else {
					data = bytes.Replace(data, []byte("idea: "+idea), []byte("idea: wrong"), 1)
				}
				path := "parley-deck/inbox/user-to-all_" + kind + ".md"
				os.WriteFile(filepath.Join(root, path), data, 0600)
				quotatest.Git(t, root, "add", "--", path)
				quotatest.Git(t, root, "commit", "-qm", "synthetic negative fixture", "--", path)
				bad.Path = path
				bad.Commit = quotatest.Git(t, root, "rev-parse", "HEAD")
				bad.Blob = quotatest.Git(t, root, "rev-parse", bad.Commit+":"+path)
				bad.SHA256 = fmt.Sprintf("%x", sha256.Sum256(data))
			case "fabricated-path":
				bad.Path = "parley-deck/inbox/../inbox/user-to-all_good.md"
			case "unquoted-content":
				bad.Quote = "Unrecorded authorization"
			case "changed-archive":
				path := filepath.Join(root, "parley-deck/inbox/archived", filepath.Base(a.Path))
				os.MkdirAll(filepath.Dir(path), 0700)
				os.WriteFile(path, []byte("altered"), 0600)
				defer os.Remove(path)
			}
			if e := quota.ValidateAuthority(root, idea, bad); e == nil {
				t.Fatal("invalid committed evidence accepted", kind)
			}
		})
	}
	os.Remove(filepath.Join(root, a.Path))
	if e := quota.ValidateAuthority(root, idea, a); e != nil {
		t.Fatal("permitted deletion invalidated committed answer", e)
	}
	if quota.HasDirective("Do not apply "+quote, quote) || !quota.HasDirective("Recorded answer\n"+quote, quote) {
		t.Fatal("directive not bound to complete line")
	}
}
func TestQuotaFixupNoticeNamesRemainingGates(t *testing.T) {
	k := quota.Kickoff{Participants: []string{"a", "b"}, Transition: &quota.Transition{ID: "test"}}
	notice := k.Notice()
	for _, name := range []string{"reviewer count", "LE-7/LE-11", "goal-checker", "require_model_diversity", "strict_gate", "vetoes", "DISPUTED", "findings"} {
		if !strings.Contains(notice, name) {
			t.Fatal("missing remaining gate", name, notice)
		}
	}
}

func TestQuotaFixupHistoryRevalidatesManualAndCatchupSnapshots(t *testing.T) {
	p := quota.Policy{Enabled: false, Scope: quota.KickoffAndMidIdea}
	k := quota.NewKickoff("snapshot", "run", p, []string{"a", "b", "c"}, nil, time.Now())
	h := &quota.History{Kickoff: &k, Current: k.Participants, Known: k.Participants}
	good := "---\nidea: snapshot\nparticipants: [a, b]\nquota_auto_exclude: false\nquota_auto_exclude_scope: kickoff-and-mid-idea\nexcluded: [c — unavailable — confirmed 2026-10-04]\n---\n"
	for _, tc := range []struct {
		name, raw string
		want      bool
	}{
		{"valid-recorded-confirmation", good, true},
		{"fabricated-body-confirmation", strings.Replace(good, "excluded:", "---\nexcluded:", 1), false},
		{"wrong-current-set", strings.Replace(good, "[a, b]", "[a, c]", 1), false},
		{"wrong-idea", strings.Replace(good, "idea: snapshot", "idea: unrelated", 1), false},
		{"invalid-date", strings.Replace(good, "2026-10-04", "2026-99-99", 1), false},
		{"policy-widening", strings.Replace(good, "exclude: false", "exclude: true", 1), false},
	} {
		t.Run(tc.name, func(t *testing.T) {
			b := quota.NewRevision(h, "run", []string{"a", "b"}, p, quota.Revision{ManualPrompt: tc.raw}, time.Now())
			if e := b.Validate(h); (e == nil) != tc.want {
				t.Fatal(e)
			}
		})
	}
	b := quota.NewRevision(h, "run", []string{"a", "b", "new"}, p, quota.Revision{ManualPrompt: strings.Replace(good, "[a, b]", "[a, b, new]\nincluded: [new — join — confirmed 2026-10-04]", 1), Catchup: map[string]string{"new": "claimed read priors"}}, time.Now())
	if e := b.Validate(h); e == nil {
		t.Fatal("nonempty fabricated catch-up accepted")
	}
}
