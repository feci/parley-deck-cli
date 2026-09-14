package app

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"os/exec"
	"path/filepath"
	"testing"

	"parley-deck-cli/internal/budget"
	"parley-deck-cli/internal/trajectory"
)

type appInterruptedCycleObserver struct{ *trajectory.Observer }

func (*appInterruptedCycleObserver) AfterCycle(context.Context, budget.CycleBinding, budget.Snapshot, string) error {
	return errors.New("injected lost observer publication")
}

func reservationAppFixture(t *testing.T) (string, *budget.CycleBinding, string) {
	t.Helper()
	root := unchangedAppFixture(t)
	b, err := budget.LoadCycleBinding(context.Background(), root, "idea-x", budget.Fixup)
	if err != nil {
		t.Fatal(err)
	}
	ctx := budget.WithCycleObserver(context.Background(), &appInterruptedCycleObserver{&trajectory.Observer{Root: root}})
	ctx, finish, err := budget.OpenCycleSession(ctx, b)
	if err != nil {
		t.Fatal(err)
	}
	defer finish()
	if _, err = budget.ChargeCycle(ctx, budget.Fixup); err == nil {
		t.Fatal("interrupted publication permitted work")
	}
	s, err := b.Store.Inspect(context.Background())
	if err != nil || len(s.Entries) != 1 {
		t.Fatal("original spent charge missing", err)
	}
	for entry := range s.Entries {
		return root, b, entry
	}
	t.Fatal("no original charge")
	return "", nil, ""
}

func TestTrajectoryReservationRecoveryCLIOutputFailureAndNoExecution(t *testing.T) {
	root, b, entry := reservationAppFixture(t)
	ctx := context.Background()
	args := []string{"recover-reservation", "--dir", root, "--idea", "idea-x", "--entry", entry}
	var out, errout bytes.Buffer
	if code := runTrajectory(ctx, args, &out, &errout); code != 0 {
		t.Fatal("preview", code, &errout)
	}
	var p struct {
		Preview trajectory.ReservationRecoveryPreview `json:"preview"`
		SHA256  string                                `json:"sha256"`
		Applied bool                                  `json:"applied"`
	}
	if err := json.Unmarshal(out.Bytes(), &p); err != nil || p.Applied || p.SHA256 != p.Preview.SHA256() {
		t.Fatal("invalid preview", err)
	}
	ledger := parentRecoveryBytes(t, filepath.Join(b.Store.Dir, "ledger.json"))
	apply := append(args, "--sha256", p.SHA256, "--yes")
	if code := runTrajectory(ctx, apply, migrationOutputFailure{}, &errout); code != 1 {
		t.Fatal("output failure hidden", code, &errout)
	}
	path := filepath.Join(filepath.Dir(b.Store.Dir), "trajectory.json")
	original := parentRecoveryBytes(t, path)
	out.Reset()
	if code := runTrajectory(ctx, apply, &out, &errout); code != 0 {
		t.Fatal("exact output recovery", code, &errout)
	}
	if !bytes.Equal(original, parentRecoveryBytes(t, path)) || !bytes.Equal(ledger, parentRecoveryBytes(t, filepath.Join(b.Store.Dir, "ledger.json"))) {
		t.Fatal("output replay changed charge or state")
	}
	s, err := trajectory.Inspect(ctx, root, "idea-x")
	if err != nil || len(s.Attempts) != 1 || s.Attempts[0].Launch != nil || s.Attempts[0].Terminal != nil {
		t.Fatal("accounting recovery invented execution", err)
	}
	if err = trajectory.RequireResolved(ctx, root, "idea-x"); err == nil {
		t.Fatal("incomplete attempt closed")
	}
}

func TestTrajectoryReservationRecoveryConcurrentProductionCLI(t *testing.T) {
	binary := trajectoryHelperBinary(t)
	root, b, entry := reservationAppFixture(t)
	args := []string{"trajectory", "recover-reservation", "--dir", root, "--idea", "idea-x", "--entry", entry}
	raw, err := exec.Command(binary, args...).CombinedOutput()
	if err != nil {
		t.Fatalf("production preview: %v %s", err, raw)
	}
	var p struct {
		SHA256 string `json:"sha256"`
	}
	if err = json.Unmarshal(raw, &p); err != nil || len(p.SHA256) != 64 {
		t.Fatal("production preview digest", err)
	}
	ledger := parentRecoveryBytes(t, filepath.Join(b.Store.Dir, "ledger.json"))
	apply := append(args, "--sha256", p.SHA256, "--yes")
	cmds := []*exec.Cmd{exec.Command(binary, apply...), exec.Command(binary, apply...)}
	out := make([]bytes.Buffer, len(cmds))
	for n, c := range cmds {
		c.Stdout = &out[n]
		c.Stderr = &out[n]
		if err = c.Start(); err != nil {
			t.Fatal(err)
		}
	}
	for n, c := range cmds {
		if err = c.Wait(); err != nil {
			t.Fatalf("concurrent production apply: %v %s", err, &out[n])
		}
	}
	if !bytes.Equal(out[0].Bytes(), out[1].Bytes()) || !bytes.Equal(ledger, parentRecoveryBytes(t, filepath.Join(b.Store.Dir, "ledger.json"))) {
		t.Fatal("concurrent replay changed original output or spend")
	}
	s, err := trajectory.Inspect(context.Background(), root, "idea-x")
	if err != nil || len(s.Attempts) != 1 {
		t.Fatal("concurrent recovery duplicated attempt", err)
	}
}

func TestTrajectoryReservationRecoveryHistoricalReplayAfterLaterWork(t *testing.T) {
	root := unchangedAppFixture(t)
	measuredUnchanged(t, root, "exit 7")
	ctx := context.Background()
	s, err := trajectory.Inspect(ctx, root, "idea-x")
	if err != nil {
		t.Fatal(err)
	}
	entry := s.Attempts[0].Charge.EntryKey
	p, err := trajectory.PreviewReservationRecovery(ctx, root, "idea-x", entry)
	if err != nil {
		t.Fatal(err)
	}
	applyUnchanged(t, root, 1)
	continueTrajectory(t, root, "original-inconclusive", false, true)
	measuredUnchanged(t, root, "exit 0")
	s, err = trajectory.Inspect(ctx, root, "idea-x")
	if err != nil {
		t.Fatal(err)
	}
	before, _ := json.Marshal(s)
	executions := parentRecoveryBytes(t, filepath.Join(root, ".parley-runtime/unchanged-executions"))
	q, err := trajectory.RecoverReservation(ctx, root, "idea-x", entry, p.SHA256())
	if err != nil || q.SHA256() != p.SHA256() {
		t.Fatal("historical replay lost original preview", err)
	}
	s, err = trajectory.Inspect(ctx, root, "idea-x")
	if err != nil {
		t.Fatal(err)
	}
	after, _ := json.Marshal(s)
	if !bytes.Equal(before, after) || !bytes.Equal(executions, parentRecoveryBytes(t, filepath.Join(root, ".parley-runtime/unchanged-executions"))) {
		t.Fatal("historical replay rewrote subsequent history or executed again")
	}
}
