// claude-1 round-04 probe: R2 knob-off grammar, R3 foreign/unknown lease owners, R1 committed authority.
// Scratch roots only (argv[1]); no provider, roster, credential or worktree mutation.
package main

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"parley-deck-cli/internal/membership"
	"parley-deck-cli/internal/pidlease"
	"parley-deck-cli/internal/procctl"
	"parley-deck-cli/internal/protocol"
	"parley-deck-cli/internal/quota"
	"parley-deck-cli/internal/runmanifest"
)

func replaceLine(text, prefix, repl string) string {
	lines := strings.Split(text, "\n")
	for i, l := range lines {
		if strings.HasPrefix(l, prefix) {
			lines[i] = repl
			return strings.Join(lines, "\n")
		}
	}
	return text
}

func knobOff(root, label string, lines string) {
	r := filepath.Join(root, label)
	os.MkdirAll(r, 0755)
	protocol.InitWorkspace(r)
	off := false
	p := quota.NewPolicy(&off, nil)
	idea, _, err := protocol.CreateIdeaWithQuota(r, "r2 grammar "+label, []string{"a", "b", "c", "d"}, nil, "deliberation", "", "run-1", &p, nil)
	if err != nil {
		fmt.Println(label, "create err", err)
		return
	}
	if k, _ := quota.ReadKickoff(idea.Path); k != nil {
		m := runmanifest.New(runmanifest.Options{Root: r, RunID: "run-1", IdeaSlug: idea.Slug, Participants: k.Participants, QuotaKickoff: k})
		if e := runmanifest.Write(r, "run-1", m); e != nil {
			fmt.Println("manifest", e)
		}
	}
	path := filepath.Join(idea.Path, "00-prompt.md")
	raw, _ := os.ReadFile(path)
	text := replaceLine(string(raw), "participants:", "participants: [a, b, c]\n"+lines)
	os.WriteFile(path, []byte(text), 0644)
	cur, known, e := protocol.QuotaMembers(idea.Path, nil)
	e2 := membership.RecordManual(r, idea.Path)
	h, e3 := quota.ReadHistory(idea.Path)
	rev := -1
	if h != nil {
		rev = h.Revision
	}
	fmt.Printf("R2 %-34s QuotaMembers current=%v known=%v err=%v | RecordManual err=%v | revision=%d readErr=%v\n", label, cur, known, e, e2, rev, e3)
	if label != "this-idea-unbracketed" || e2 != nil {
		return
	}
	// Known return of d: first without an included: marker, then with one.
	raw, _ = os.ReadFile(path)
	text = replaceLine(string(raw), "participants:", "participants: [a, b, c, d]")
	os.WriteFile(path, []byte(text), 0644)
	cur, _, e = protocol.QuotaMembers(idea.Path, nil)
	fmt.Printf("R2 %-34s QuotaMembers current=%v err=%v\n", "return-d-without-included-marker", cur, e)
	text = replaceLine(text, "participants:", "participants: [a, b, c, d]\nincluded: [d — returned after owner confirmation — confirmed 2026-10-05]")
	os.WriteFile(path, []byte(text), 0644)
	cur, _, e = protocol.QuotaMembers(idea.Path, nil)
	e2 = membership.RecordManual(r, idea.Path)
	h, _ = quota.ReadHistory(idea.Path)
	fmt.Printf("R2 %-34s QuotaMembers current=%v err=%v | RecordManual err=%v revision=%d current=%v\n", "return-d-with-included-marker", cur, e, e2, h.Revision, h.Current)
}

func lease(root string) {
	dir := filepath.Join(root, "lease")
	os.MkdirAll(dir, 0700)
	host, _ := os.Hostname()
	boot := procctl.CurrentBootID()
	fmt.Printf("R3 local host=%q boot=%q\n", host, boot)
	tok := make([]byte, 16)
	for i := range tok {
		tok[i] = byte(i + 7)
	}
	for _, tc := range []struct{ name, host, boot string }{
		{"foreign-host-dead-pid", host + "-other", boot},
		{"foreign-boot-dead-pid", host, boot + "-previous-boot"},
		{"unknown-host-dead-pid", "", boot},
		{"unknown-boot-dead-pid", host, ""},
		{"same-host-boot-dead-pid", host, boot},
	} {
		path := filepath.Join(dir, tc.name+".lease")
		o := pidlease.Owner{Version: 1, PID: 2147480000, Token: hex.EncodeToString(tok), Host: tc.host, Boot: tc.boot, Identity: "demo/run-dead"}
		raw, _ := json.Marshal(o)
		os.WriteFile(path, raw, 0600)
		l, err := pidlease.TryAcquire(path, "demo/run-new")
		after, _ := os.ReadFile(path)
		fmt.Printf("R3 %-26s acquired=%v err=%v ownerUnchanged=%v\n", tc.name, l != nil, err, string(after) == string(raw))
		if l != nil {
			l.Release()
		}
	}
}

func git(root string, args ...string) (string, error) {
	cmd := exec.Command("git", append([]string{"-C", root, "-c", "user.name=claude1-r4 probe", "-c", "user.email=probe@example.invalid", "-c", "commit.gpgsign=false", "-c", "core.hooksPath=/dev/null"}, args...)...)
	out, err := cmd.CombinedOutput()
	return strings.TrimSpace(string(out)), err
}

func authority(root string) {
	r := filepath.Join(root, "authority")
	os.MkdirAll(filepath.Join(r, "parley-deck", "inbox"), 0755)
	git(r, "init", "-q")
	quote := "Quota revision: {\"participants\":[\"a\",\"b\",\"c\"],\"policy\":{\"enabled\":true,\"scope\":\"kickoff-and-mid-idea\"}}"
	rel := "parley-deck/inbox/user-to-all_r4-return.md"
	body := "---\nfrom: user\nto: all\nidea: demo-idea\n---\n\nProbe owner answer (synthetic, not real authorization).\n" + quote + "\n"
	os.WriteFile(filepath.Join(r, rel), []byte(body), 0644)
	git(r, "add", "--", rel)
	git(r, "commit", "-qm", "probe answer", "--", rel)
	commit, _ := git(r, "rev-parse", "HEAD")
	a, err := quota.BindAuthority(r, "demo-idea", rel, commit, quote)
	fmt.Printf("R1 bind committed answer err=%v sha=%s\n", err, a.SHA256[:12])
	check := func(label string, x quota.Authority) {
		fmt.Printf("R1 %-40s validate err=%v\n", label, quota.ValidateAuthority(r, "demo-idea", x))
	}
	check("live inbox copy present", a)
	arch := filepath.Join(r, "parley-deck/inbox/archived", filepath.Base(rel))
	os.MkdirAll(filepath.Dir(arch), 0755)
	os.Rename(filepath.Join(r, rel), arch)
	check("moved to inbox/archived (permitted)", a)
	os.Remove(arch)
	check("deleted from inbox (permitted)", a)
	os.WriteFile(arch, []byte(body+"edited\n"), 0644)
	check("modified archived copy present", a)
	os.Remove(arch)
	w := a
	w.Path = "parley-deck/inbox/user-to-all_unrelated.md"
	check("unrelated/fabricated path", w)
	w = a
	w.Commit = strings.Repeat("a", 40)
	check("fabricated commit", w)
	w = a
	sum := sha256.Sum256([]byte("x"))
	w.SHA256 = hex.EncodeToString(sum[:])
	check("wrong sha256", w)
	w = a
	w.Quote = "Quota revision: {\"participants\":[\"a\",\"b\"],\"policy\":{\"enabled\":true,\"scope\":\"kickoff-and-mid-idea\"}}"
	check("quote not in committed answer", w)
	fmt.Printf("R1 %-40s validate err=%v\n", "wrong idea", quota.ValidateAuthority(r, "other-idea", a))
}

func main() {
	root := os.Args[1]
	os.MkdirAll(root, 0755)
	knobOff(root, "this-idea-unbracketed", "excluded: d — HTTP 429 Weekly/Monthly Limit Exhausted, reset 2026-10-05 06:14 local (preflight class provider-failure:rate-limit) — confirmed 2026-10-03")
	knobOff(root, "documented-bracketed", "excluded: [d — unavailable — confirmed 2026-10-04]")
	knobOff(root, "reason-contains-emdash", "excluded: [d — unavailable — quota exhausted — confirmed 2026-10-04]")
	knobOff(root, "date-with-trailing-note", "excluded: [d — unavailable — confirmed 2026-10-04 via relay]")
	knobOff(root, "inbox-reference-style", "excluded: [d]   # §9.0 — see inbox/claude-1-to-all_x.md")
	knobOff(root, "hyphen-separators", "excluded: [d - unavailable - confirmed 2026-10-04]")
	if len(os.Args) > 2 {
		return
	}
	lease(root)
	authority(root)
}
