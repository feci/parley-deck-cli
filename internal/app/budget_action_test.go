package app

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	"parley-deck-cli/internal/budget"
)

func TestBudgetActionReadOnlyReplayAndMissingEvidence(t *testing.T) {
	s := budget.Store{Dir: t.TempDir(), Scope: "action-fixture"}
	sha, e := budget.ActionInputDigest("synthetic private task")
	if e != nil {
		t.Fatal(e)
	}
	ctx := budget.WithActionInput(context.Background(), budget.ActionInput{Operation: "fixture", Basis: "runtime-input", InputSHA256: sha, RunID: "run"})
	if _, e = s.Reserve(ctx, budget.Request{ID: "one", Kind: budget.Fixup}, budget.Limits{}); e != nil {
		t.Fatal(e)
	}
	inspect := []string{"action", "inspect", "--ledger", s.Dir, "--scope", s.Scope}
	var out, stderr bytes.Buffer
	if code := runBudgetPlatformControl(context.Background(), inspect, &out, &stderr, false, false); code != 0 {
		t.Fatalf("read-only inspection %d %s", code, stderr.String())
	}
	var rows []budget.ActionReceipt
	if e = json.Unmarshal(out.Bytes(), &rows); e != nil || len(rows) != 1 {
		t.Fatalf("inspect: %v", e)
	}
	before, e := os.ReadFile(filepath.Join(s.Dir, "ledger.json"))
	if e != nil {
		t.Fatal(e)
	}
	args := []string{"action", "replay", "--ledger", s.Dir, "--scope", s.Scope, "--entry-key", rows[0].EntryKey, "--expected-identity-sha256", rows[0].IdentitySHA256}
	if code := runBudgetPlatformControl(context.Background(), args, migrationOutputFailure{}, &stderr, false, false); code != 1 {
		t.Fatal("output failure was ignored")
	}
	out.Reset()
	stderr.Reset()
	if code := runBudgetPlatformControl(context.Background(), args, &out, &stderr, false, false); code != 0 {
		t.Fatalf("accounting replay %d %s", code, stderr.String())
	}
	var got budget.ActionReceipt
	if e = json.Unmarshal(out.Bytes(), &got); e != nil {
		t.Fatal(e)
	}
	if !reflect.DeepEqual(got, rows[0]) || got.Permission != "none" || got.ExecutionStatus != "not-established" {
		t.Fatal("accounting replay changed receipt or asserted execution")
	}
	if code := runBudgetPlatformControl(context.Background(), append(append([]string{}, args...), "--yes"), io.Discard, io.Discard, true, true); code != 2 {
		t.Fatal("read-only replay accepted an activation flag")
	}
	bad := append([]string{}, args...)
	bad[len(bad)-1] = strings.Repeat("f", 64)
	if code := runBudgetPlatformControl(context.Background(), bad, io.Discard, io.Discard, false, false); code == 0 {
		t.Fatal("changed expected identity passed")
	}
	after, e := os.ReadFile(filepath.Join(s.Dir, "ledger.json"))
	if e != nil || !bytes.Equal(before, after) {
		t.Fatal("replay changed original charge")
	}
	if e = os.Remove(filepath.Join(s.Dir, "ledger.json")); e != nil {
		t.Fatal(e)
	}
	if code := runBudgetPlatformControl(context.Background(), args, io.Discard, io.Discard, false, false); code == 0 {
		t.Fatal("missing accounting became replayable")
	}
	if _, e = os.Lstat(filepath.Join(s.Dir, "ledger.json")); !os.IsNotExist(e) {
		t.Fatal("read-only command recreated a ledger")
	}
}
