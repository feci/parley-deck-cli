package budget

import (
	"context"
	"encoding/json"
	"strings"
	"testing"
	"time"
)

func TestCycleIntentPinsActualRequestBeforeOriginalCharge(t *testing.T) {
	ctx := WithActionInput(context.Background(), ActionInput{Operation: "fixup", Basis: "runtime-input", InputSHA256: key("original input"), RunID: "original-run"})
	b, err := EnsureCycleBinding(ctx, t.TempDir(), "fixture", Fixup, 5, 0, "", "")
	if err != nil {
		t.Fatal(err)
	}
	b.Policy.TrajectorySHA256 = strings.Repeat("b", 64)
	before, err := b.Store.Inspect(ctx)
	if err != nil {
		t.Fatal(err)
	}
	zero := int64(0)
	req, limits := Request{ID: "original-attempt", Kind: Fixup, ReserveMicros: &zero}, cycleLimits(b.Policy)
	i, err := newCycleReservationIntent(ctx, *b, before, req, limits)
	if err != nil || i.Action == nil || i.Action.ReservationSHA256 != "" {
		t.Fatal("precharge invented a reservation", err)
	}
	charged, err := b.Store.Reserve(ctx, req, limits)
	if err != nil {
		t.Fatal(err)
	}
	charge, err := i.CheckCharge(charged)
	if err != nil || charge.ActionSHA256 == "" {
		t.Fatal("original charge not matched", err)
	}
	if err = i.CheckBefore(b.Policy, charged); err == nil {
		t.Fatal("spent intent became uncharged")
	}
	for name, mutate := range map[string]func(*Snapshot){
		"entry": func(s *Snapshot) { s.Entries[key("different")] = s.Entries[i.EntryKey]; delete(s.Entries, i.EntryKey) },
		"epoch": func(s *Snapshot) { s.StartedAt = s.StartedAt.Add(time.Second) },
		"scope": func(s *Snapshot) { s.Scope += "-different" },
		"timestamp": func(s *Snapshot) {
			e := s.Entries[i.EntryKey]
			e.ReservedAt = e.ReservedAt.Add(time.Second)
			s.Entries[i.EntryKey] = e
		},
		"kind": func(s *Snapshot) { e := s.Entries[i.EntryKey]; e.Kind = CrossReview; s.Entries[i.EntryKey] = e },
		"cost": func(s *Snapshot) {
			e := s.Entries[i.EntryKey]
			x := int64(1)
			e.ReserveMicros = &x
			s.Entries[i.EntryKey] = e
		},
		"input":          func(s *Snapshot) { s.Entries[i.EntryKey].Action.Input.InputSHA256 = key("different") },
		"limits":         func(s *Snapshot) { s.Entries[i.EntryKey].Action.LimitsSHA256 = key("different") },
		"missing-action": func(s *Snapshot) { e := s.Entries[i.EntryKey]; e.Action = nil; s.Entries[i.EntryKey] = e },
	} {
		t.Run(name, func(t *testing.T) {
			raw, _ := json.Marshal(charged)
			var changed Snapshot
			if err := json.Unmarshal(raw, &changed); err != nil {
				t.Fatal(err)
			}
			mutate(&changed)
			if _, err := i.CheckCharge(changed); err == nil {
				t.Fatal("changed original accounting accepted")
			}
		})
	}
	settled, err := b.Store.Settle(ctx, req.ID, &zero)
	if err != nil {
		t.Fatal(err)
	}
	if again, err := i.CheckCharge(settled); err != nil || !reflectChargeEqual(charge, again) {
		t.Fatal("settlement rewrote original charge", err)
	}
}

func reflectChargeEqual(a, b CycleCharge) bool {
	x, _ := json.Marshal(a)
	y, _ := json.Marshal(b)
	return string(x) == string(y)
}

func TestCycleIntentPreservesOriginalAndExtendedCeilings(t *testing.T) {
	ctx := context.Background()
	root := t.TempDir()
	b, err := EnsureCycleBinding(ctx, root, "fixture", Fixup, 2, 0, "", "")
	if err != nil {
		t.Fatal(err)
	}
	if err = ActivateCycleTrajectory(ctx, root, "fixture", CyclePolicyDigest(b.Policy), strings.Repeat("a", 64), func(CycleBinding, Snapshot) error { return nil }); err != nil {
		t.Fatal(err)
	}
	b, err = LoadCycleBinding(ctx, root, "fixture", Fixup)
	if err != nil {
		t.Fatal(err)
	}
	var intents []CycleReservationIntent
	for n := 0; n < 2; n++ {
		state, err := b.Store.Inspect(ctx)
		if err != nil {
			t.Fatal(err)
		}
		zero := int64(0)
		id := []string{"original", "extended"}[n]
		input := WithActionInput(ctx, ActionInput{Operation: "fixup", Basis: "policy-only", InputSHA256: CyclePolicyDigest(b.Policy)})
		request := Request{ID: id, Kind: Fixup, ReserveMicros: &zero}
		i, err := newCycleReservationIntent(input, *b, state, request, cycleLimits(b.Policy))
		if err != nil {
			t.Fatal("valid extended intent refused", err)
		}
		intents = append(intents, i)
		if _, err = b.Store.Reserve(input, request, cycleLimits(b.Policy)); err != nil {
			t.Fatal(err)
		}
		_, err = ExtendCycleBudget(ctx, root, "fixture", Fixup, CycleExtensionRequest{DecisionID: id, ExpectedPolicySHA256: CyclePolicyDigest(b.Policy), Reason: "fixture original grant", Maximum: b.Policy.Maximum + 2})
		if err != nil {
			t.Fatal(err)
		}
		b, err = LoadCycleBinding(ctx, root, "fixture", Fixup)
		if err != nil {
			t.Fatal(err)
		}
		state, err = b.Store.Inspect(ctx)
		if err != nil {
			t.Fatal(err)
		}
		for _, original := range intents {
			if err = original.CheckPolicy(b.Policy); err != nil {
				t.Fatal("later grant invalidated original policy", err)
			}
			if _, err = original.CheckCharge(state); err != nil {
				t.Fatal("later grant changed original charge", err)
			}
			changed := cloneCyclePolicy(b.Policy)
			changed.IdeaPath += "-other"
			if err = original.CheckPolicy(changed); err == nil {
				t.Fatal("changed policy authority accepted")
			}
		}
	}
}
