package app

import (
	"bytes"
	"context"
	"strings"
	"testing"

	"parley-deck-cli/internal/budget"
)

func TestBudgetConfigureUsesAttendedExplicitCeilings(t *testing.T) {
	root := t.TempDir()
	args := []string{"configure", "--dir", root, "--idea", "fixture", "--max-launches", "1", "--max-cost-micros", "1000000", "--wall-clock", "1h", "--reserve-micros", "1000000", "--yes"}
	var out, stderr bytes.Buffer
	if code := runBudgetControl(context.Background(), args, &out, &stderr, false); code != 2 {
		t.Fatalf("unattended configure: %d %s", code, stderr.String())
	}
	if b, err := budget.LoadLaunchBinding(context.Background(), root, "fixture"); err != nil || b != nil {
		t.Fatal("unattended command created policy")
	}
	stderr.Reset()
	if code := runBudgetControl(context.Background(), args, &out, &stderr, true); code != 0 {
		t.Fatalf("attended configure: %d %s", code, stderr.String())
	}
	b, err := budget.LoadLaunchBinding(context.Background(), root, "fixture")
	if err != nil || b == nil || b.Policy.MaxLaunches != 1 || b.Policy.WallClockMS != 3600000 {
		t.Fatalf("wrong persisted policy: %+v %v", b, err)
	}
}

func TestBudgetConfigureUnsupportedPlatformDoesNotActivate(t *testing.T) {
	root := t.TempDir()
	var stdout, stderr bytes.Buffer
	code := runBudgetPlatformControl(context.Background(), []string{"configure", "--dir", root}, &stdout, &stderr, false, true)
	if code != 2 || !strings.Contains(stderr.String(), "budget configure: attended policy activation is unavailable") {
		t.Fatalf("unsupported configure: %d %s", code, stderr.String())
	}
	if b, err := budget.LoadLaunchBinding(context.Background(), root, ""); err != nil || b != nil {
		t.Fatalf("unsupported configuration mutated state: %+v %v", b, err)
	}
}
