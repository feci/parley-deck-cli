package budget

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
	"time"
)

func actionTestInput(t *testing.T, task string) ActionInput {
	t.Helper()
	sha, e := ActionInputDigest(struct{ Task string }{task})
	if e != nil {
		t.Fatal(e)
	}
	return ActionInput{Operation: "fixture-operation", Basis: "runtime-input", InputSHA256: sha, RunID: "fixture-run"}
}
func actionTestReserve(t *testing.T, s Store, id string, in ActionInput) (ActionReceipt, Limits) {
	t.Helper()
	limits := Limits{Actions: map[Kind]int{Launch: 5}, CostMicros: 20}
	state, e := s.Reserve(WithActionInput(context.Background(), in), Request{ID: id, Kind: Launch, ReserveMicros: micros(4)}, limits)
	if e != nil {
		t.Fatal(e)
	}
	r, e := actionReceipt(state, key(id))
	if e != nil {
		t.Fatal(e)
	}
	return r, limits
}

func TestActionIdentityExactReplayAndDistinctSpentAttempts(t *testing.T) {
	s := testStore(t)
	in := actionTestInput(t, "private prompt material")
	original, limits := actionTestReserve(t, s, "first", in)
	before := originRead(t, filepath.Join(s.Dir, "ledger.json"))
	if bytes.Contains(before, []byte("private prompt material")) {
		t.Fatal("raw input leaked into accounting")
	}
	if original.Permission != "none" || original.ExecutionStatus != "not-established" || original.Identity == nil {
		t.Fatal("receipt implied execution or lost identity")
	}
	for _, mutation := range []string{"same", "input", "operation", "run", "cost", "limits", "kind"} {
		want := in
		req := Request{ID: "first", Kind: Launch, ReserveMicros: micros(4)}
		cap := Limits{Actions: map[Kind]int{Launch: 5}, CostMicros: 20}
		switch mutation {
		case "input":
			want = actionTestInput(t, "changed")
		case "operation":
			want.Operation = "other"
		case "run":
			want.RunID = "another"
		case "cost":
			req.ReserveMicros = micros(5)
		case "limits":
			cap.CostMicros = 21
		case "kind":
			req.Kind = Fixup
		}
		_, e := s.Reserve(WithActionInput(context.Background(), want), req, cap)
		if !errors.Is(e, ErrReserved) || mutation != "same" && !errors.Is(e, ErrActionConflict) {
			t.Fatalf("%s duplicate: %v", mutation, e)
		}
		if !bytes.Equal(before, originRead(t, filepath.Join(s.Dir, "ledger.json"))) {
			t.Fatalf("%s duplicate changed ledger", mutation)
		}
	}
	// An actual retry is a distinct attempt, even with identical logical input.
	retry, _ := actionTestReserve(t, s, "retry", in)
	if original.Identity.LogicalSHA256 != retry.Identity.LogicalSHA256 || original.EntryKey == retry.EntryKey || original.IdentitySHA256 == retry.IdentitySHA256 {
		t.Fatal("logical identity and actual attempt were conflated")
	}
	if _, e := s.Settle(context.Background(), "first", nil); e != nil {
		t.Fatal(e)
	}
	if _, e := s.ReconcileUnknown(context.Background(), "first", "cost-review", 8, "conservative synthetic ceiling"); e != nil {
		t.Fatal(e)
	}
	replay, e := s.ReplayAction(context.Background(), original.EntryKey, original.IdentitySHA256)
	if e != nil || !reflect.DeepEqual(replay, original) {
		t.Fatalf("historical receipt changed after settlement: %+v %v", replay, e)
	}
	state, e := s.Inspect(context.Background())
	if e != nil || len(state.Entries) != 2 {
		t.Fatalf("retry charges lost: %v", e)
	}
	if original.Identity.LimitsSHA256 != migrationDigest(limits) {
		t.Fatal("original reservation limits were not frozen")
	}
	if _, e := s.ReplayAction(context.Background(), original.EntryKey, retry.IdentitySHA256); !errors.Is(e, ErrActionConflict) {
		t.Fatal("another attempt's identity replayed")
	}
}
func TestActionIdentityReservationPublicationFailure(t *testing.T) {
	for _, published := range []bool{false, true} {
		t.Run(fmt.Sprint(published), func(t *testing.T) {
			s := testStore(t)
			ctx := WithActionInput(context.Background(), actionTestInput(t, "attempt"))
			injected := errors.New("lost reservation result")
			s.persist = func(path string, b []byte) error {
				if published {
					if e := writeSynced(path, b); e != nil {
						return e
					}
				}
				return injected
			}
			req := Request{ID: "interrupted", Kind: Fixup, ReserveMicros: micros(0)}
			if _, e := s.Reserve(ctx, req, Limits{Actions: map[Kind]int{Fixup: 1}}); !errors.Is(e, injected) {
				t.Fatal(e)
			}
			s.persist = nil
			if published {
				rows, e := s.InspectActions(context.Background())
				if e != nil || len(rows) != 1 {
					t.Fatalf("published charge was lost: %v", e)
				}
				replay, e := s.ReplayAction(context.Background(), rows[0].EntryKey, rows[0].IdentitySHA256)
				if e != nil || replay.Permission != "none" {
					t.Fatal("recovery granted another execution")
				}
				if _, e := s.Reserve(ctx, req, Limits{Actions: map[Kind]int{Fixup: 1}}); !errors.Is(e, ErrReserved) {
					t.Fatalf("published attempt was charged/launched again: %v", e)
				}
			} else {
				if _, e := s.InspectActions(context.Background()); e == nil {
					t.Fatal("unpublished in-memory state became durable evidence")
				}
				if _, e := s.Reserve(ctx, req, Limits{Actions: map[Kind]int{Fixup: 1}}); e != nil {
					t.Fatal(e)
				}
			}
			state, e := s.Inspect(context.Background())
			if e != nil || len(state.Entries) != 1 {
				t.Fatalf("incorrect final count: %v", e)
			}
		})
	}
}

func TestActionIdentityReplayRejectsChangedOriginalCharge(t *testing.T) {
	for _, mutation := range []string{"kind", "reserved-at", "epoch", "reserved-cost", "unknown-cost"} {
		t.Run(mutation, func(t *testing.T) {
			s := testStore(t)
			r, _ := actionTestReserve(t, s, "original", actionTestInput(t, "same-task"))
			state, err := s.Inspect(context.Background())
			if err != nil {
				t.Fatal(err)
			}
			entry := state.Entries[r.EntryKey]
			switch mutation {
			case "kind":
				entry.Kind = Fixup
			case "reserved-at":
				entry.ReservedAt = entry.ReservedAt.Add(time.Nanosecond)
			case "epoch":
				state.StartedAt = state.StartedAt.Add(-time.Second)
			case "reserved-cost":
				entry.ReserveMicros = micros(5)
			case "unknown-cost":
				entry.ReserveMicros = nil
			}
			state.Entries[r.EntryKey] = entry
			publishReservationFixture(t, s, state)
			before := originRead(t, filepath.Join(s.Dir, "ledger.json"))
			if _, err := s.ReplayAction(context.Background(), r.EntryKey, r.IdentitySHA256); err == nil {
				t.Fatal("original identity accepted a changed reservation receipt")
			}
			if !bytes.Equal(before, originRead(t, filepath.Join(s.Dir, "ledger.json"))) {
				t.Fatal("refusal changed accounting")
			}
			// A structurally consistent replacement is still a different receipt.
			// Pinning only its input would let this rewrite pass exact replay.
			entry.Action.ReservationSHA256 = originalChargeDigest(state.Scope, state.StartedAt, r.EntryKey, entry)
			state.Entries[r.EntryKey] = entry
			publishReservationFixture(t, s, state)
			if _, err := s.Inspect(context.Background()); err != nil {
				t.Fatalf("replacement must be structurally valid: %v", err)
			}
			if _, err := s.ReplayAction(context.Background(), r.EntryKey, r.IdentitySHA256); !errors.Is(err, ErrActionConflict) {
				t.Fatalf("changed charge replay did not conflict: %v", err)
			}
		})
	}
}
func TestActionIdentityLegacyAndMalformedBindingRefuseReplay(t *testing.T) {
	s := testStore(t)
	if _, e := s.Reserve(context.Background(), Request{ID: "legacy", Kind: Fixup, ReserveMicros: micros(0)}, Limits{}); e != nil {
		t.Fatal(e)
	}
	rows, e := s.InspectActions(context.Background())
	if e != nil || len(rows) != 1 || rows[0].Binding != "legacy-unbound" || rows[0].Identity != nil {
		t.Fatal("invented semantic identity for old accounting")
	}
	if _, e = s.ReplayAction(context.Background(), rows[0].EntryKey, key("invented")); e == nil {
		t.Fatal("legacy entry was certified retroactively")
	}
	if _, e = s.Reserve(WithActionInput(context.Background(), actionTestInput(t, "retrofit")), Request{ID: "legacy", Kind: Fixup, ReserveMicros: micros(0)}, Limits{}); !errors.Is(e, ErrActionConflict) {
		t.Fatal("new metadata adopted old identity")
	}
	r, _ := actionTestReserve(t, s, "current", actionTestInput(t, "bound"))
	raw := originRead(t, filepath.Join(s.Dir, "ledger.json"))
	for _, mutation := range []string{"remove", "entry", "input", "version", "limits"} {
		var state Snapshot
		if e = json.Unmarshal(raw, &state); e != nil {
			t.Fatal(e)
		}
		entry := state.Entries[r.EntryKey]
		switch mutation {
		case "remove":
			entry.Action = nil
		case "entry":
			entry.Action.EntrySHA256 = key("other")
		case "input":
			entry.Action.Input.InputSHA256 = key("other")
		case "version":
			entry.Action.Version = 2
		case "limits":
			entry.Action.LimitsSHA256 = ""
		}
		state.Entries[r.EntryKey] = entry
		publishReservationFixture(t, s, state)
		if _, e = s.ReplayAction(context.Background(), r.EntryKey, r.IdentitySHA256); e == nil {
			t.Fatalf("%s action mutation passed replay", mutation)
		}
	}
}
func TestActionIdentityActiveSessionPinsSemanticFields(t *testing.T) {
	for _, kind := range []Kind{DriverStep, Fixup, CrossReview} {
		t.Run(string(kind), func(t *testing.T) {
			f := newActiveReservationFixture(t, kind)
			state, e := f.store.Inspect(context.Background())
			if e != nil {
				t.Fatal(e)
			}
			for k, entry := range state.Entries {
				if entry.Action == nil || entry.Action.Input.Basis != "policy-only" {
					t.Fatal("direct session falsely claimed a complete runtime input")
				}
				entry.Action.Input.Operation = "different-operation"
				entry.Action.LogicalSHA256 = actionLogicalDigest(state.Scope, entry.Action.Input)
				state.Entries[k] = entry
			}
			publishReservationFixture(t, f.store, state)
			if _, e = f.store.Inspect(context.Background()); e != nil {
				t.Fatalf("mutation must remain structurally valid: %v", e)
			}
			if e = f.charge(); e == nil {
				t.Fatal("changed semantic identity reused active receipt")
			}
		})
	}
}
func TestActionIdentityInvalidInputDoesNotCreateAccounting(t *testing.T) {
	for _, input := range []ActionInput{{}, {Operation: "prompt with secret text", Basis: "runtime-input", InputSHA256: key("x")}, {Operation: "run", Basis: "claimed-success", InputSHA256: key("x")}} {
		dir := filepath.Join(t.TempDir(), "absent")
		s := Store{Dir: dir, Scope: "scope"}
		if _, e := s.Reserve(WithActionInput(context.Background(), input), Request{ID: "one", Kind: Launch}, Limits{}); e == nil {
			t.Fatal("invalid identity allowed work")
		}
		if _, e := os.Stat(dir); !os.IsNotExist(e) {
			t.Fatal("invalid input created accounting")
		}
	}
	if _, e := ActionInputDigest(make(chan int)); e == nil {
		t.Fatal("unencodable input silently hashed as empty")
	}
}

func TestActionIdentitySessionRejectsChangedContext(t *testing.T) {
	for _, kind := range []Kind{DriverStep, Fixup, CrossReview} {
		for _, charged := range []bool{false, true} {
			t.Run(fmt.Sprintf("%s/charged=%t", kind, charged), func(t *testing.T) {
				original := actionTestInput(t, "original-request")
				ctx := WithActionInput(context.Background(), original)
				var store Store
				var charge func(context.Context) error
				var reopen func(context.Context) error
				var finish func()
				var err error
				if kind == DriverStep {
					var b *StepBinding
					b, err = EnsureStepBinding(ctx, t.TempDir(), "idea", 5, 0)
					if err != nil {
						t.Fatal(err)
					}
					store = b.Store
					ctx, finish, err = OpenStepSession(ctx, b)
					charge = ChargeStep
					reopen = func(c context.Context) error { _, close, e := OpenStepSession(c, b); close(); return e }
				} else {
					var b *CycleBinding
					b, err = EnsureCycleBinding(ctx, t.TempDir(), "idea", kind, 5, 0, "", "")
					if err != nil {
						t.Fatal(err)
					}
					store = b.Store
					ctx, finish, err = OpenCycleSession(ctx, b)
					charge = func(c context.Context) error { _, e := ChargeCycle(c, kind); return e }
					reopen = func(c context.Context) error { _, close, e := OpenCycleSession(c, b); close(); return e }
				}
				if err != nil {
					t.Fatal(err)
				}
				defer finish()
				changed := actionTestInput(t, "changed-request")
				if sessionActionInput(InheritActionInput(ctx, changed)) != original {
					t.Fatal("nested runner replaced its outer operation")
				}
				if charged {
					if err = charge(ctx); err != nil {
						t.Fatal(err)
					}
				}
				altered := WithActionInput(ctx, changed)
				if err = reopen(altered); err == nil {
					t.Fatal("changed input reopened original accounting session")
				}
				if err = charge(altered); err == nil {
					t.Fatal("changed input continued under original reservation")
				}
				if err = charge(ctx); err == nil {
					t.Fatal("restoring input revived refused accounting session")
				}
				rows, err := store.InspectActions(context.Background())
				if err != nil {
					t.Fatal(err)
				}
				want := 0
				if charged {
					want = 1
				}
				if len(rows) != want {
					t.Fatalf("charges=%d want=%d", len(rows), want)
				}
			})
		}
	}
}
func TestActionIdentityProcessChild(t *testing.T) {
	dir := os.Getenv("PARLEY_ACTION_IDENTITY_CHILD")
	if dir == "" {
		return
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	ctx = WithActionInput(ctx, ActionInput{Operation: "process-fixture", Basis: "runtime-input", InputSHA256: key("same-input"), RunID: "same-run"})
	_, e := (Store{Dir: dir, Scope: "process-scope"}).Reserve(ctx, Request{ID: "same-attempt", Kind: Launch, ReserveMicros: micros(1)}, Limits{Actions: map[Kind]int{Launch: 1}})
	if e == nil {
		fmt.Println("granted-once")
		return
	}
	if errors.Is(e, ErrReserved) {
		fmt.Println("duplicate-refused")
		return
	}
	t.Fatal(e)
}
func TestActionIdentitySeparateProcessesCannotDuplicateGrant(t *testing.T) {
	dir := t.TempDir()
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	var output [2]bytes.Buffer
	cmds := make([]*exec.Cmd, 2)
	for i := range cmds {
		cmds[i] = exec.CommandContext(ctx, os.Args[0], "-test.run=^TestActionIdentityProcessChild$")
		cmds[i].Env = append(os.Environ(), "PARLEY_ACTION_IDENTITY_CHILD="+dir)
		cmds[i].Stdout = &output[i]
		cmds[i].Stderr = &output[i]
		if e := cmds[i].Start(); e != nil {
			t.Fatal(e)
		}
	}
	granted, refused := 0, 0
	for i, c := range cmds {
		if e := c.Wait(); e != nil {
			t.Fatalf("child: %v %s", e, output[i].String())
		}
		if strings.Contains(output[i].String(), "granted-once") {
			granted++
		}
		if strings.Contains(output[i].String(), "duplicate-refused") {
			refused++
		}
	}
	if granted != 1 || refused != 1 {
		t.Fatalf("grants=%d duplicates=%d", granted, refused)
	}
	rows, e := (Store{Dir: dir, Scope: "process-scope"}).InspectActions(ctx)
	if e != nil || len(rows) != 1 || rows[0].Identity == nil {
		t.Fatalf("durable process accounting: %v", e)
	}
}
