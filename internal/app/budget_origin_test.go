package app

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"os"
	"path/filepath"
	"reflect"
	"testing"

	"parley-deck-cli/internal/budget"
)

func TestBudgetOriginAttendanceAndExactOutputRecovery(t *testing.T) {
	ctx := context.Background()
	dir := t.TempDir()
	target := t.TempDir()
	release, e := budget.AcquireResourceGuard(ctx, dir)
	if e != nil {
		t.Fatal(e)
	}
	release()
	inspect := []string{"origin", "inspect", "--resource", dir, "--target-lock-dir", target}
	var out, errout bytes.Buffer
	if code := runBudgetPlatformControl(ctx, inspect, &out, &errout, false, false); code != 0 {
		t.Fatalf("inspect %d %s", code, errout.String())
	}
	var preview budget.LockOriginPreview
	if e := json.Unmarshal(out.Bytes(), &preview); e != nil {
		t.Fatal(e)
	}
	before, e := os.ReadFile(filepath.Join(dir, "lock-origin"))
	if e != nil {
		t.Fatal(e)
	}
	args := []string{"origin", "apply", "--resource", dir, "--target-lock-dir", target, "--expected-sha256", preview.SHA256, "--decision-id", "fixture", "--reason", "synthetic relocation", "--yes"}
	for _, state := range []struct{ supported, attended bool }{{false, false}, {false, true}, {true, false}} {
		if code := runBudgetPlatformControl(ctx, args, io.Discard, io.Discard, state.supported, state.attended); code != 2 {
			t.Fatalf("unattended control returned %d", code)
		}
		after, e := os.ReadFile(filepath.Join(dir, "lock-origin"))
		if e != nil || !bytes.Equal(before, after) {
			t.Fatal("attendance refusal mutated origin")
		}
	}
	for _, extra := range []string{"--attended", "--operator", "--writers-stopped"} {
		if code := runBudgetPlatformControl(ctx, append(append([]string{}, args...), extra), io.Discard, io.Discard, true, true); code != 2 {
			t.Fatalf("accepted fabricated control %s", extra)
		}
	}
	for _, extra := range []string{"--yes=false", "--reason="} {
		if code := runBudgetPlatformControl(ctx, append(append([]string{}, inspect...), extra), io.Discard, io.Discard, false, false); code != 2 {
			t.Fatal("inspect accepted mutation flag")
		}
	}
	if code := runBudgetPlatformControl(ctx, args, migrationOutputFailure{}, &errout, true, true); code != 1 {
		t.Fatalf("output fault %d %s", code, errout.String())
	}
	out.Reset()
	errout.Reset()
	if code := runBudgetPlatformControl(ctx, args, &out, &errout, true, true); code != 0 {
		t.Fatalf("exact replay %d %s", code, errout.String())
	}
	var got budget.LockOriginPreview
	if e := json.Unmarshal(out.Bytes(), &got); e != nil {
		t.Fatal(e)
	}
	if !reflect.DeepEqual(got, preview) {
		t.Fatal("lost exact applied preview after output failure")
	}
	release, e = budget.AcquireResourceGuard(ctx, dir)
	if e != nil {
		t.Fatal(e)
	}
	release()
}
