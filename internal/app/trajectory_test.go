package app

import (
	"bytes"
	"context"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"parley-deck-cli/internal/budget"
	"parley-deck-cli/internal/trajectory"
)

func TestTrajectoryCLIExactActivationAndIncompleteClose(t *testing.T) {
	root, ideaDir := gateScratchRepo(t, "checks:\n  - name: unit\n    command: true\n")
	var out, errout bytes.Buffer
	invoke := func(attended bool, args ...string) int {
		out.Reset()
		errout.Reset()
		return runTrajectoryControl(context.Background(), args, &out, &errout, attended)
	}
	init := []string{"initialize", "--dir", root, "--idea", "idea-x", "--max-fixups", "3", "--yes"}
	if code := invoke(false, init...); code != 2 {
		t.Fatalf("unattended initialization: %d %s", code, errout.String())
	}
	b, err := budget.LoadCycleBinding(context.Background(), root, "idea-x", budget.Fixup)
	if err != nil || b != nil {
		t.Fatalf("unattended operation changed state: %+v %v", b, err)
	}
	if code := invoke(true, init...); code != 0 {
		t.Fatalf("initialization: %d %s", code, errout.String())
	}
	preview := []string{"configure", "--dir", root, "--idea", "idea-x", "--implementer", "kimi-1"}
	if code := invoke(false, preview...); code != 0 {
		t.Fatalf("preview: %d %s", code, errout.String())
	}
	var p struct {
		Policy    trajectory.Policy `json:"policy"`
		SHA       string            `json:"policy_sha256"`
		Cycle     string            `json:"cycle_policy_sha256"`
		Activated bool              `json:"activated"`
	}
	if err = json.Unmarshal(out.Bytes(), &p); err != nil || p.Activated {
		t.Fatalf("preview: %+v %v", p, err)
	}
	s, err := trajectory.Inspect(context.Background(), root, "idea-x")
	if err != nil || s != nil {
		t.Fatalf("preview activated policy: %+v %v", s, err)
	}
	apply := append(append([]string{}, preview...), "--policy-sha256", p.SHA, "--cycle-policy-sha256", p.Cycle, "--yes")
	if code := invoke(false, apply...); code != 2 {
		t.Fatalf("unattended activation: %d %s", code, errout.String())
	}
	if code := invoke(true, apply...); code != 0 {
		t.Fatalf("activation: %d %s", code, errout.String())
	}
	if code := invoke(true, apply...); code != 0 {
		t.Fatalf("exact replay: %d %s", code, errout.String())
	}
	b, err = budget.LoadCycleBinding(context.Background(), root, "idea-x", budget.Fixup)
	if err != nil {
		t.Fatal(err)
	}
	ctx := budget.WithCycleObserver(context.Background(), &trajectory.Observer{Root: root})
	ctx, finish, err := budget.OpenCycleSession(ctx, b)
	if err != nil {
		t.Fatal(err)
	}
	defer finish()
	if _, err = budget.ChargeCycle(ctx, budget.Fixup); err != nil {
		t.Fatal(err)
	}
	before, err := os.ReadFile(filepath.Join(ideaDir, "IMPLEMENTATION.md"))
	if err != nil {
		t.Fatal(err)
	}
	ops := gateOps(root, ideaDir)
	err = ops.Complete(context.Background())
	if err == nil || !strings.Contains(err.Error(), "trajectory") {
		t.Fatalf("completion did not enforce canonical idea trajectory: %v", err)
	}
	after, err := os.ReadFile(filepath.Join(ideaDir, "IMPLEMENTATION.md"))
	if err != nil || !bytes.Equal(before, after) {
		t.Fatal("refused completion changed implementation")
	}
	if code := invoke(false, "inspect", "--dir", root, "--idea", "idea-x"); code != 0 {
		t.Fatalf("inspect pending attempt: %d %s", code, errout.String())
	}
	var inspected trajectory.State
	if err = json.Unmarshal(out.Bytes(), &inspected); err != nil || len(inspected.Attempts) != 1 || inspected.Attempts[0].Launch != nil {
		t.Fatalf("pending reservation inspection: %+v %v", inspected, err)
	}
	if code := invoke(true, apply...); code == 0 {
		t.Fatal("activation replay reset a charged trajectory")
	}
	if inspected.Version != 2 || inspected.Policy.Version != 2 || inspected.BaselineArchive.SHA256 == "" || inspected.Attempts[0].BeforeArchive != inspected.BaselineArchive {
		t.Fatal("CLI inspection lacks reconstructible baseline binding")
	}
	archive := filepath.Join(filepath.Dir(b.Store.Dir), "trajectory-snapshots", inspected.BaselineArchive.SHA256+".tar")
	if err = os.Remove(archive); err != nil {
		t.Fatal(err)
	}
	if code := invoke(false, "inspect", "--dir", root, "--idea", "idea-x"); code == 0 {
		t.Fatal("CLI reported intact state after archive loss")
	}
	if err = ops.Complete(context.Background()); err == nil {
		t.Fatal("application completion ignored missing archive")
	}
}
func TestTrajectoryCLIRejectsChangedScopeAndNoChecks(t *testing.T) {
	for _, checks := range []string{"", "checks:\n  - name: unit\n    command: true\n"} {
		t.Run(strings.TrimSpace(checks), func(t *testing.T) {
			root, ideaDir := gateScratchRepo(t, checks)
			_, err := budget.EnsureCycleBinding(context.Background(), root, "idea-x", budget.Fixup, 3, 0, "", ideaDir)
			if err != nil {
				t.Fatal(err)
			}
			var out, errout bytes.Buffer
			args := []string{"configure", "--dir", root, "--idea", "idea-x", "--implementer", "builder"}
			code := runTrajectoryControl(context.Background(), args, &out, &errout, false)
			if checks == "" {
				if code == 0 {
					t.Fatal("absent material scope accepted")
				}
				return
			}
			if code != 0 {
				t.Fatal(errout.String())
			}
			var p map[string]any
			if err = json.Unmarshal(out.Bytes(), &p); err != nil {
				t.Fatal(err)
			}
			if err = os.WriteFile(filepath.Join(root, "src", "a.go"), []byte("package changed\n"), 0644); err != nil {
				t.Fatal(err)
			}
			args = append(args, "--policy-sha256", p["policy_sha256"].(string), "--cycle-policy-sha256", p["cycle_policy_sha256"].(string), "--yes")
			if code = runTrajectoryControl(context.Background(), args, &out, &errout, true); code == 0 {
				t.Fatal("changed preview source activated")
			}
		})
	}
}
