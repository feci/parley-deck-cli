package app

import (
	"bytes"
	"context"
	"encoding/json"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"parley-deck-cli/internal/evidence"
)

func appRefusal(t *testing.T, dir string) string {
	t.Helper()
	r, err := newVerificationRefusal("helper", "execution", "idea-x", "run", "reviewer", sha256Hex("request"), sha256Hex("report"), "attempt-invocation")
	if err != nil {
		t.Fatal(err)
	}
	sum, err := evidence.RetainVerificationRefusal(dir, r)
	if err != nil {
		t.Fatal(err)
	}
	return sum
}

func TestRefusalRecoveryCommitsOnlyExactPathAndInvalidatesOldTree(t *testing.T) {
	root, dir := gateScratchRepo(t, twoCriterionContract())
	gateGit(t, root, "config", "user.name", "fixture")
	gateGit(t, root, "config", "user.email", "fixture@example.invalid")
	if err := os.WriteFile(filepath.Join(root, "unrelated"), []byte("staged"), 0600); err != nil {
		t.Fatal(err)
	}
	gateGit(t, root, "add", "unrelated")
	if err := os.WriteFile(filepath.Join(root, "unrelated"), []byte("working"), 0600); err != nil {
		t.Fatal(err)
	}
	before, err := evidence.TreeDigest(root)
	if err != nil {
		t.Fatal(err)
	}
	sum := appRefusal(t, dir)
	pendingTree, err := evidence.TreeDigest(root)
	if err != nil || pendingTree != before {
		t.Fatalf("pending observation changed source: %s %v", pendingTree, err)
	}
	if err := requireCommittedRefusals(context.Background(), root, dir); err == nil {
		t.Fatal("uncommitted refusal permits verification")
	}
	var out, stderr bytes.Buffer
	if code := runEvidenceVerify(context.Background(), []string{"refusals", "recover", "--dir", root, "--idea", "idea-x", "--expected-sha256", sum}, &out, &stderr); code != 0 {
		t.Fatalf("recovery: %d %s", code, stderr.String())
	}
	after, err := evidence.TreeDigest(root)
	if err != nil || after == before {
		t.Fatal("canonical refusal was excluded from source binding")
	}
	if err := requireCommittedRefusals(context.Background(), root, dir); err != nil {
		t.Fatal(err)
	}
	changed, err := refusalGit(context.Background(), root, "diff-tree", "--no-commit-id", "--name-only", "-r", "HEAD")
	want := "parley-deck/ideas/idea-x/verification-refusals/" + sum + ".json\n"
	if err != nil || string(changed) != want {
		t.Fatalf("committed unrelated change: %q %v", changed, err)
	}
	staged, err := refusalGit(context.Background(), root, "show", ":unrelated")
	if err != nil || string(staged) != "staged" {
		t.Fatal("unrelated index changed")
	}
	working, _ := os.ReadFile(filepath.Join(root, "unrelated"))
	if string(working) != "working" {
		t.Fatal("unrelated work overwritten")
	}
	head, _ := refusalGit(context.Background(), root, "rev-parse", "HEAD")
	if err := recoverVerificationRefusal(context.Background(), root, dir, sum); err != nil {
		t.Fatal(err)
	}
	replayHead, _ := refusalGit(context.Background(), root, "rev-parse", "HEAD")
	if !bytes.Equal(head, replayHead) {
		t.Fatal("exact recovery created a duplicate commit")
	}
	// A refusal and its recovery never create successful evidence or acceptance.
	if _, err := os.Stat(evidence.ReportPath(dir)); !os.IsNotExist(err) {
		t.Fatal("recovery created evidence")
	}
}

func TestRefusalCommitFailureRemainsInspectableAndRecoverable(t *testing.T) {
	root, dir := gateScratchRepo(t, twoCriterionContract())
	gateGit(t, root, "config", "user.name", "fixture")
	gateGit(t, root, "config", "user.email", "fixture@example.invalid")
	hook := filepath.Join(root, ".git", "hooks", "pre-commit")
	if err := os.WriteFile(hook, []byte("#!/bin/sh\nexit 1\n"), 0700); err != nil {
		t.Fatal(err)
	}
	sum := appRefusal(t, dir)
	if err := recoverVerificationRefusal(context.Background(), root, dir, sum); err == nil {
		t.Fatal("failed commit claimed durability")
	}
	var out, stderr bytes.Buffer
	if code := runEvidenceVerify(context.Background(), []string{"refusals", "inspect", "--dir", root, "--idea", "idea-x"}, &out, &stderr); code != 0 {
		t.Fatalf("inspect: %d %s", code, stderr.String())
	}
	if !strings.Contains(out.String(), `"committed": false`) || !strings.Contains(out.String(), sum) {
		t.Fatalf("false commit claim: %s", out.String())
	}
	if err := requireCommittedRefusals(context.Background(), root, dir); err == nil {
		t.Fatal("failed commit allowed verification")
	}
	if err := os.Remove(hook); err != nil {
		t.Fatal(err)
	}
	if err := recoverVerificationRefusal(context.Background(), root, dir, sum); err != nil {
		t.Fatal(err)
	}
	if err := requireCommittedRefusals(context.Background(), root, dir); err != nil {
		t.Fatal(err)
	}
}

func TestRefusalRecoveryCLIRejectsUnboundAndWrongScopes(t *testing.T) {
	root, dir := gateScratchRepo(t, twoCriterionContract())
	sum := appRefusal(t, dir)
	for _, args := range [][]string{
		{"recover", "--idea", "idea-x"},
		{"recover", "--idea", "idea-x", "--expected-sha256", "wrong"},
		{"inspect", "--idea", "idea-x", "--expected-sha256", sum},
		{"recover", "--idea", "../idea-x", "--expected-sha256", sum},
		{"recover", "--idea", "idea-x", "--expected-sha256", sha256Hex("different")},
	} {
		var out, stderr bytes.Buffer
		args = append([]string{"refusals"}, append(args, "--dir", root)...)
		if code := runEvidenceVerify(context.Background(), args, &out, &stderr); code == 0 {
			t.Fatalf("invalid recovery accepted: %v", args)
		}
	}
	if err := requireCommittedRefusals(context.Background(), root, dir); err == nil {
		t.Fatal("invalid recoveries granted verification")
	}
}

func TestRefusalLabelsDoNotPublishProviderSecrets(t *testing.T) {
	r, err := newVerificationRefusal("driver", "launch", "https://example.invalid/private", "Bearer-secret", "sk-secret", "bad digest", sha256Hex("original"), "password-secret")
	if err != nil {
		t.Fatal(err)
	}
	raw, _ := json.Marshal(r)
	for _, secret := range []string{"https://", "Bearer-secret", "sk-secret", "password-secret", "bad digest"} {
		if bytes.Contains(raw, []byte(secret)) {
			t.Fatalf("unsafe canonical label: %s", raw)
		}
	}
}

// This child runs only the production helper, with no driver remaining to
// observe its return. The test's shell parent is explicitly killed first.
func TestRefusalOrphanHelperChild(t *testing.T) {
	path := os.Getenv("PARLEY_REFUSAL_CHILD_REQUEST")
	if path == "" {
		return
	}
	release := os.Getenv("PARLEY_REFUSAL_CHILD_RELEASE")
	ready := os.Getenv("PARLEY_REFUSAL_CHILD_READY")
	done := os.Getenv("PARLEY_REFUSAL_CHILD_DONE")
	if err := os.WriteFile(ready, []byte("ready"), 0600); err != nil {
		t.Fatal(err)
	}
	deadline := time.Now().Add(10 * time.Second)
	for {
		if _, err := os.Stat(release); err == nil {
			break
		}
		if time.Now().After(deadline) {
			t.Fatal("release never arrived")
		}
		time.Sleep(10 * time.Millisecond)
	}
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	if err := executeIndependentVerification(context.Background(), path, sha256Hex(string(data))); err == nil {
		t.Fatal("identity-less helper accepted")
	}
	if err := os.WriteFile(done, []byte("refused"), 0600); err != nil {
		t.Fatal(err)
	}
}

func TestRefusalHelperSurvivesStoppedParent(t *testing.T) {
	root, dir := gateScratchRepo(t, twoCriterionContract())
	root, dir, err := verificationRefusalScope(root, "idea-x")
	if err != nil {
		t.Fatal(err)
	}
	attempt := filepath.Join(root, ".parley-runtime", "evidence-verification", "attempt-orphan")
	if err := os.MkdirAll(attempt, 0700); err != nil {
		t.Fatal(err)
	}
	request := filepath.Join(attempt, "request.json")
	if err := writeVerificationJSON(request, evidenceVerificationRequest{Version: 1, AttemptID: "attempt-orphan", Root: root, Idea: "idea-x", RunID: "orphan", Verifier: "reviewer"}); err != nil {
		t.Fatal(err)
	}
	control := t.TempDir()
	ready, release, done := filepath.Join(control, "ready"), filepath.Join(control, "release"), filepath.Join(control, "done")
	cmd := exec.Command("sh", "-c", `"$1" -test.run=^TestRefusalOrphanHelperChild$ -test.count=1 & child=$!; wait "$child"`, "parent", os.Args[0])
	cmd.Env = append(os.Environ(), "PARLEY_REFUSAL_CHILD_REQUEST="+request, "PARLEY_REFUSAL_CHILD_READY="+ready, "PARLEY_REFUSAL_CHILD_RELEASE="+release, "PARLEY_REFUSAL_CHILD_DONE="+done, "PARLEY_RUN_ID=", "PARLEY_AGENT_ID=", "PARLEY_PROC_MARKER=")
	// Use a file, not parent-owned pipes; an orphan must not keep Cmd.Wait blocked.
	log, err := os.Create(filepath.Join(control, "child.log"))
	if err != nil {
		t.Fatal(err)
	}
	defer log.Close()
	cmd.Stdout, cmd.Stderr = log, log
	if err := cmd.Start(); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = cmd.Process.Kill() })
	waitFile := func(path string) {
		t.Helper()
		deadline := time.Now().Add(10 * time.Second)
		for {
			if _, err := os.Stat(path); err == nil {
				return
			}
			if time.Now().After(deadline) {
				raw, _ := os.ReadFile(log.Name())
				t.Fatalf("no helper signal: %s", raw)
			}
			time.Sleep(10 * time.Millisecond)
		}
	}
	waitFile(ready)
	if err := cmd.Process.Kill(); err != nil {
		t.Fatal(err)
	}
	if err := cmd.Wait(); err == nil {
		t.Fatal("parent did not terminate")
	}
	if err := os.WriteFile(release, nil, 0600); err != nil {
		t.Fatal(err)
	}
	waitFile(done)
	entries, err := evidence.InspectVerificationRefusals(dir)
	if err != nil || len(entries) != 1 || entries[0].Record == nil || entries[0].Record.Observer != "helper" || entries[0].Canonical {
		t.Fatalf("orphan lost pending refusal: %+v %v", entries, err)
	}
	if err := requireCommittedRefusals(context.Background(), root, dir); err == nil {
		t.Fatal("orphan's pending failure allowed verification")
	}
	gateGit(t, root, "config", "user.name", "fixture")
	gateGit(t, root, "config", "user.email", "fixture@example.invalid")
	if err := recoverVerificationRefusal(context.Background(), root, dir, entries[0].SHA256); err != nil {
		t.Fatal(err)
	}
}

func TestRefusalRecoveryChild(t *testing.T) {
	root := os.Getenv("PARLEY_REFUSAL_RECOVERY_ROOT")
	if root == "" {
		return
	}
	sum := os.Getenv("PARLEY_REFUSAL_RECOVERY_SHA")
	dir := filepath.Join(root, "parley-deck", "ideas", "idea-x")
	barrier := os.Getenv("PARLEY_REFUSAL_RECOVERY_BARRIER")
	ready := os.Getenv("PARLEY_REFUSAL_RECOVERY_READY")
	if err := os.WriteFile(ready, nil, 0600); err != nil {
		t.Fatal(err)
	}
	deadline := time.Now().Add(10 * time.Second)
	for {
		if _, err := os.Stat(barrier); err == nil {
			break
		}
		if time.Now().After(deadline) {
			t.Fatal("recovery barrier timeout")
		}
		time.Sleep(10 * time.Millisecond)
	}
	if err := recoverVerificationRefusal(context.Background(), root, dir, sum); err != nil {
		t.Fatal(err)
	}
}

func TestRefusalConcurrentRecoveryProcesses(t *testing.T) {
	root, dir := gateScratchRepo(t, twoCriterionContract())
	gateGit(t, root, "config", "user.name", "fixture")
	gateGit(t, root, "config", "user.email", "fixture@example.invalid")
	sum := appRefusal(t, dir)
	control := t.TempDir()
	barrier := filepath.Join(control, "release")
	children := make([]*exec.Cmd, 2)
	logs := make([]bytes.Buffer, 2)
	for i := range children {
		ready := filepath.Join(control, []string{"ready-a", "ready-b"}[i])
		cmd := exec.Command(os.Args[0], "-test.run=^TestRefusalRecoveryChild$", "-test.count=1")
		cmd.Env = append(os.Environ(), "PARLEY_REFUSAL_RECOVERY_ROOT="+root, "PARLEY_REFUSAL_RECOVERY_SHA="+sum, "PARLEY_REFUSAL_RECOVERY_BARRIER="+barrier, "PARLEY_REFUSAL_RECOVERY_READY="+ready)
		cmd.Stdout, cmd.Stderr = &logs[i], &logs[i]
		if err := cmd.Start(); err != nil {
			t.Fatal(err)
		}
		children[i] = cmd
		t.Cleanup(func() { _ = cmd.Process.Kill() })
	}
	deadline := time.Now().Add(10 * time.Second)
	for {
		ready, _ := filepath.Glob(filepath.Join(control, "ready-*"))
		if len(ready) == 2 {
			break
		}
		if time.Now().After(deadline) {
			t.Fatal("recovery children never reached barrier")
		}
		time.Sleep(10 * time.Millisecond)
	}
	if err := os.WriteFile(barrier, nil, 0600); err != nil {
		t.Fatal(err)
	}
	for i, cmd := range children {
		if err := cmd.Wait(); err != nil {
			t.Fatalf("child %d: %v %s", i, err, logs[i].String())
		}
	}
	count, err := refusalGit(context.Background(), root, "rev-list", "--count", "HEAD")
	if err != nil || strings.TrimSpace(string(count)) != "2" {
		t.Fatalf("duplicate refusal commit: %s %v", count, err)
	}
	if err := requireCommittedRefusals(context.Background(), root, dir); err != nil {
		t.Fatal(err)
	}
}
