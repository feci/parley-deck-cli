package budget

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"math"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"
	"testing"
	"time"
)

func runtimeFixture(t *testing.T, root string, kind Kind) runtimeBinding {
	t.Helper()
	ctx := context.Background()
	if kind == Launch {
		reserve := int64(50)
		b, err := ConfigureLaunchBudget(ctx, root, "idea", LaunchPolicy{MaxLaunches: 1, MaxCostMicros: 100, WallClockMS: 3600000, ReserveMicros: &reserve})
		if err != nil {
			t.Fatal(err)
		}
		return runtimeBinding{b.Policy, b.Store}
	}
	b, err := EnsureStepBinding(ctx, root, "idea", 1, time.Hour)
	if err != nil {
		t.Fatal(err)
	}
	return runtimeBinding{b.Policy, b.Store}
}

func reserveRuntime(ctx context.Context, b runtimeBinding, id string) (Snapshot, error) {
	if p, ok := b.policy.(LaunchPolicy); ok {
		return (&LaunchBinding{p, b.store}).Reserve(ctx, id)
	}
	return (&StepBinding{b.policy.(StepPolicy), b.store}).reserve(ctx, id)
}

func runtimeGrant(t *testing.T, root string, kind Kind, id string, ceilings PolicyCeilings) PolicyExtensionRequest {
	t.Helper()
	s, err := InspectRuntimeBudget(context.Background(), root, "idea", kind)
	if err != nil {
		t.Fatal(err)
	}
	return PolicyExtensionRequest{id, s.PolicySHA256, "Explicit finite fixture decision", ceilings}
}

func TestRuntimeExtensionPreservesLedgerClockAndReplay(t *testing.T) {
	ctx := context.Background()
	for _, kind := range []Kind{Launch, DriverStep} {
		t.Run(string(kind), func(t *testing.T) {
			root := t.TempDir()
			b := runtimeFixture(t, root, kind)
			if _, err := reserveRuntime(ctx, b, "first"); err != nil {
				t.Fatal(err)
			}
			ledger := filepath.Join(b.store.Dir, "ledger.json")
			before, err := os.ReadFile(ledger)
			if err != nil {
				t.Fatal(err)
			}
			ceilings := b.policy.ceilings()
			ceilings.Actions = 2
			ceilings.WallClockNS = int64(2 * time.Hour)
			if kind == Launch {
				ceilings.CostMicros = 200
			}
			r := runtimeGrant(t, root, kind, "first-grant", ceilings)
			s, err := ExtendRuntimeBudget(ctx, root, "idea", kind, r)
			if err != nil {
				t.Fatal(err)
			}
			after, _ := os.ReadFile(ledger)
			if !bytes.Equal(before, after) || s.Spent != 1 {
				t.Fatal("grant rewrote spent ledger")
			}
			current, err := b.current()
			if err != nil {
				t.Fatal(err)
			}
			v, original, history := current.policy.extensionMetadata()
			if v != 2 || original == nil || *original != b.policy.ceilings() || len(history) != 1 || history[0].StartedAt != s.StartedAt {
				t.Fatalf("lost original policy/history: %+v", current.policy)
			}
			if _, err := reserveRuntime(ctx, b, "second"); err != nil {
				t.Fatalf("cached binding ignored grant: %v", err)
			}
			if _, err := reserveRuntime(ctx, b, "third-denied"); !errors.Is(err, ErrLimit) {
				t.Fatalf("extension became unlimited: %v", err)
			}
			ceilings.Actions = 3
			if kind == Launch {
				ceilings.CostMicros = 300
			}
			r2 := runtimeGrant(t, root, kind, "second-grant", ceilings)
			latest, err := ExtendRuntimeBudget(ctx, root, "idea", kind, r2)
			if err != nil {
				t.Fatal(err)
			}
			replay, err := ExtendRuntimeBudget(ctx, root, "idea", kind, r)
			if err != nil || replay.PolicySHA256 != latest.PolicySHA256 || replay.Spent != 2 {
				t.Fatalf("old exact replay: %+v %v", replay, err)
			}
			conflict := r
			conflict.Reason = "different"
			if _, err := ExtendRuntimeBudget(ctx, root, "idea", kind, conflict); err == nil {
				t.Fatal("conflicting replay accepted")
			}
			stale := r
			stale.DecisionID = "stale-new"
			if _, err := ExtendRuntimeBudget(ctx, root, "idea", kind, stale); err == nil {
				t.Fatal("stale new decision accepted")
			}
			if kind == Launch {
				originalPolicy := b.policy.(LaunchPolicy)
				replayed, err := ConfigureLaunchBudget(ctx, root, "idea", originalPolicy)
				if err != nil || replayed.Policy.MaxLaunches != 3 {
					t.Fatalf("original configure replay lost grant: %+v %v", replayed, err)
				}
				if err := RequireMonetaryBinding(replayed, 0.0001); err != nil {
					t.Fatal(err)
				}
				if err := RequireMonetaryBinding(replayed, 0.0002); err != nil {
					t.Fatal(err)
				}
				if err := RequireMonetaryBinding(replayed, 0.00025); err == nil {
					t.Fatal("unrecorded monetary cap accepted")
				}
			} else {
				for _, limits := range []struct {
					steps int
					wall  time.Duration
				}{{1, time.Hour}, {2, 2 * time.Hour}, {0, 0}} {
					if _, err := EnsureStepBinding(ctx, root, "idea", limits.steps, limits.wall); err != nil {
						t.Fatal(err)
					}
				}
				if _, err := EnsureStepBinding(ctx, root, "idea", 4, 0); err == nil {
					t.Fatal("unrecorded step limit accepted")
				}
			}
		})
	}
}

func TestRuntimeExtensionPublicationRecoveryAndConcurrentPreview(t *testing.T) {
	ctx := context.Background()
	for _, kind := range []Kind{Launch, DriverStep} {
		t.Run(string(kind), func(t *testing.T) {
			root := t.TempDir()
			b := runtimeFixture(t, root, kind)
			ceilings := b.policy.ceilings()
			ceilings.Actions = 2
			r := runtimeGrant(t, root, kind, "grant", ceilings)
			sentinel := errors.New("publication failure")
			if _, err := b.extend(ctx, r, func(string, []byte) error { return sentinel }); !errors.Is(err, sentinel) {
				t.Fatal(err)
			}
			old, err := InspectRuntimeBudget(ctx, root, "idea", kind)
			if err != nil || old.PolicySHA256 != r.ExpectedPolicySHA256 {
				t.Fatalf("prepublication failure changed policy: %v", err)
			}
			if _, err := b.extend(ctx, r, func(path string, data []byte) error {
				if err := writeSynced(path, data); err != nil {
					return err
				}
				return sentinel
			}); !errors.Is(err, sentinel) {
				t.Fatal(err)
			}
			if _, err := ExtendRuntimeBudget(ctx, root, "idea", kind, r); err != nil {
				t.Fatalf("exact recovery: %v", err)
			}
			ceilings.Actions = 3
			r = runtimeGrant(t, root, kind, "next-a", ceilings)
			start := make(chan struct{})
			results := make(chan error, 2)
			for _, id := range []string{"next-a", "next-b"} {
				go func(id string) {
					<-start
					req := r
					req.DecisionID = id
					_, err := ExtendRuntimeBudget(ctx, root, "idea", kind, req)
					results <- err
				}(id)
			}
			close(start)
			wins := 0
			for i := 0; i < 2; i++ {
				if <-results == nil {
					wins++
				}
			}
			if wins != 1 {
				t.Fatalf("preview winners=%d", wins)
			}
		})
	}
}

func TestRuntimeExtensionRefusesCorruptionAndLostCharges(t *testing.T) {
	ctx := context.Background()
	for _, kind := range []Kind{Launch, DriverStep} {
		t.Run(string(kind), func(t *testing.T) {
			root := t.TempDir()
			b := runtimeFixture(t, root, kind)
			if _, err := reserveRuntime(ctx, b, "first"); err != nil {
				t.Fatal(err)
			}
			ceilings := b.policy.ceilings()
			ceilings.Actions = 2
			if _, err := ExtendRuntimeBudget(ctx, root, "idea", kind, runtimeGrant(t, root, kind, "grant", ceilings)); err != nil {
				t.Fatal(err)
			}
			path := filepath.Join(filepath.Dir(b.store.Dir), "policy.json")
			raw, err := os.ReadFile(path)
			if err != nil {
				t.Fatal(err)
			}
			mutations := []func(map[string]any){
				func(m map[string]any) { delete(m, "original") },
				func(m map[string]any) { m["original"] = nil },
				func(m map[string]any) { m["version"] = 1 },
				func(m map[string]any) { m["original"].(map[string]any)["actions"] = 0 },
				func(m map[string]any) { delete(m["extensions"].([]any)[0].(map[string]any), "spent") },
				func(m map[string]any) {
					m["extensions"].([]any)[0].(map[string]any)["previous_sha256"] = strings.Repeat("f", 64)
				},
				func(m map[string]any) {
					m["extensions"].([]any)[0].(map[string]any)["ceilings"].(map[string]any)["actions"] = 0
				},
				func(m map[string]any) { m["unknown-field"] = true },
			}
			for _, mutate := range mutations {
				var m map[string]any
				if err := json.Unmarshal(raw, &m); err != nil {
					t.Fatal(err)
				}
				mutate(m)
				bad, _ := json.Marshal(m)
				if err := os.WriteFile(path, bad, 0600); err != nil {
					t.Fatal(err)
				}
				if _, err := reserveRuntime(ctx, b, "bad"); err == nil {
					t.Fatalf("corrupt grant allowed work: %s", bad)
				}
				if _, err := InspectRuntimeBudget(ctx, root, "idea", kind); err == nil {
					t.Fatal("corrupt grant inspected as valid")
				}
			}
			for _, bad := range []string{strings.Replace(string(raw), `"spent": 1`, `"spent": 1, "spent": 0`, 1), strings.Replace(string(raw), `"spent": 1`, `"Spent": 1`, 1)} {
				if err := os.WriteFile(path, []byte(bad), 0600); err != nil {
					t.Fatal(err)
				}
				if _, err := InspectRuntimeBudget(ctx, root, "idea", kind); err == nil {
					t.Fatal("aliased/duplicate grant field accepted")
				}
			}
			if err := os.WriteFile(path, raw, 0600); err != nil {
				t.Fatal(err)
			}
			if _, err := b.store.update(ctx, func(s *Snapshot, _ time.Time) error { delete(s.Entries, key("first")); return nil }); err != nil {
				t.Fatal(err)
			}
			if _, err := reserveRuntime(ctx, b, "after-loss"); err == nil {
				t.Fatal("grant ignored a lost charged attempt")
			}
		})
	}
}

func TestRuntimeExtensionFiniteAxesAndUnknownExposure(t *testing.T) {
	ctx := context.Background()
	root := t.TempDir()
	b := runtimeFixture(t, root, Launch)
	valid := b.policy.ceilings()
	valid.Actions = 2
	r := runtimeGrant(t, root, Launch, "grant", valid)
	for _, change := range []func(*PolicyExtensionRequest){
		func(r *PolicyExtensionRequest) { r.Ceilings.Actions = 0 },
		func(r *PolicyExtensionRequest) { r.Ceilings.Actions = math.MaxInt },
		func(r *PolicyExtensionRequest) { r.Ceilings.Actions = 1 },
		func(r *PolicyExtensionRequest) { r.Ceilings.CostMicros = 0 },
		func(r *PolicyExtensionRequest) { r.Ceilings.CostMicros = math.MaxInt64 },
		func(r *PolicyExtensionRequest) { r.Ceilings.WallClockNS = 1 },
		func(r *PolicyExtensionRequest) { r.DecisionID = "" },
		func(r *PolicyExtensionRequest) { r.Reason = strings.Repeat("x", 1025) },
	} {
		req := r
		change(&req)
		if _, err := ExtendRuntimeBudget(ctx, root, "idea", Launch, req); err == nil {
			t.Fatalf("invalid extension accepted: %+v", req)
		}
	}
	// Historical unknown exposure remains unknown until an explicit ceiling is
	// reconciled. A count-only grant cannot erase it or authorize a priced call.
	if _, err := b.store.Reserve(ctx, Request{ID: "unknown", Kind: Launch}, Limits{}); err != nil {
		t.Fatal(err)
	}
	r.Ceilings.CostMicros = 200
	if _, err := ExtendRuntimeBudget(ctx, root, "idea", Launch, r); !errors.Is(err, ErrUnknownCost) {
		t.Fatalf("unknown cost extension: %v", err)
	}
	r.Ceilings.CostMicros = 100
	if _, err := ExtendRuntimeBudget(ctx, root, "idea", Launch, r); err != nil {
		t.Fatal(err)
	}
	if _, err := reserveRuntime(ctx, b, "still-unknown"); !errors.Is(err, ErrUnknownCost) {
		t.Fatalf("count extension erased unknown cost: %v", err)
	}
	if _, err := b.store.ReconcileUnknown(ctx, "unknown", "bound", 100, "Explicit fixture bound"); err != nil {
		t.Fatal(err)
	}
	r = runtimeGrant(t, root, Launch, "cost-grant", PolicyCeilings{Actions: 2, CostMicros: 200, WallClockNS: int64(time.Hour)})
	if _, err := ExtendRuntimeBudget(ctx, root, "idea", Launch, r); err != nil {
		t.Fatal(err)
	}
	if _, err := reserveRuntime(ctx, b, "priced"); err != nil {
		t.Fatal(err)
	}
	if _, err := ExtendRuntimeBudget(ctx, t.TempDir(), "idea", Launch, r); err == nil {
		t.Fatal("extension initialized missing policy")
	}
}

func TestRuntimeExtensionNestedSessionAndClock(t *testing.T) {
	ctx := context.Background()
	root := t.TempDir()
	rb := runtimeFixture(t, root, DriverStep)
	b := &StepBinding{rb.policy.(StepPolicy), rb.store}
	scoped, finish, err := OpenStepSession(ctx, b)
	if err != nil {
		t.Fatal(err)
	}
	defer finish()
	if err := ChargeStep(scoped); err != nil {
		t.Fatal(err)
	}
	c := b.Policy.ceilings()
	c.Actions = 2
	r := runtimeGrant(t, root, DriverStep, "grant", c)
	if _, err := ExtendRuntimeBudget(ctx, root, "idea", DriverStep, r); err != nil {
		t.Fatal(err)
	}
	current, err := LoadStepBinding(ctx, root, "idea")
	if err != nil {
		t.Fatal(err)
	}
	nested, done, err := OpenStepSession(scoped, current)
	if err != nil {
		t.Fatal(err)
	}
	defer done()
	if err := ChargeStep(nested); err != nil {
		t.Fatal(err)
	}
	s, err := InspectRuntimeBudget(ctx, root, "idea", DriverStep)
	if err != nil || s.Spent != 1 {
		t.Fatalf("grant split a nested transition: %+v %v", s, err)
	}
	current.Store.now = func() time.Time { return s.StartedAt.Add(-time.Second) }
	if _, err := current.reserve(ctx, "clock-back"); !errors.Is(err, ErrClockSkew) {
		t.Fatalf("clock rollback: %v", err)
	}
	cancelled, cancel := context.WithCancel(ctx)
	cancel()
	c.Actions = 3
	r = runtimeGrant(t, root, DriverStep, "cancelled", c)
	if _, err := ExtendRuntimeBudget(cancelled, root, "idea", DriverStep, r); !errors.Is(err, context.Canceled) {
		t.Fatal(err)
	}
}

func TestRuntimeExtensionSharesWorktreeScopeAndConcurrentCharges(t *testing.T) {
	ctx := context.Background()
	for _, kind := range []Kind{Launch, DriverStep} {
		t.Run(string(kind), func(t *testing.T) {
			root := t.TempDir()
			git := func(args ...string) {
				t.Helper()
				if out, err := exec.Command("git", append([]string{"-C", root}, args...)...).CombinedOutput(); err != nil {
					t.Fatalf("git: %v %s", err, out)
				}
			}
			git("init", "-q")
			git("-c", "user.name=t", "-c", "user.email=t@t", "commit", "--allow-empty", "-qm", "fixture")
			other := filepath.Join(t.TempDir(), "other")
			git("worktree", "add", "-q", "-b", "other", other)
			b := runtimeFixture(t, root, kind)
			c := b.policy.ceilings()
			c.Actions = 3
			if kind == Launch {
				c.CostMicros = 500
			}
			if _, err := ExtendRuntimeBudget(ctx, other, "idea", kind, runtimeGrant(t, root, kind, "grant", c)); err != nil {
				t.Fatal(err)
			}
			var wg sync.WaitGroup
			results := make(chan error, 8)
			for i := 0; i < 8; i++ {
				wg.Add(1)
				go func(i int) { defer wg.Done(); _, err := reserveRuntime(ctx, b, string(rune('a'+i))); results <- err }(i)
			}
			wg.Wait()
			close(results)
			success := 0
			for err := range results {
				if err == nil {
					success++
				} else if !errors.Is(err, ErrLimit) {
					t.Fatal(err)
				}
			}
			s, err := InspectRuntimeBudget(ctx, other, "idea", kind)
			if err != nil || success != 3 || s.Spent != 3 {
				t.Fatalf("shared grant: successes=%d status=%+v err=%v", success, s, err)
			}
		})
	}
}

func TestRuntimeExtensionExpiredClockUsesOriginalActivation(t *testing.T) {
	ctx := context.Background()
	for _, kind := range []Kind{Launch, DriverStep} {
		t.Run(string(kind), func(t *testing.T) {
			root := t.TempDir()
			b := runtimeFixture(t, root, kind)
			s, err := b.inspect(ctx)
			if err != nil {
				t.Fatal(err)
			}
			b.store.now = func() time.Time { return s.StartedAt.Add(2 * time.Hour) }
			if _, err := reserveRuntime(ctx, b, "expired"); !errors.Is(err, ErrLimit) {
				t.Fatalf("original clock did not expire: %v", err)
			}
			c := b.policy.ceilings()
			c.WallClockNS = int64(90 * time.Minute)
			r := runtimeGrant(t, root, kind, "time", c)
			if _, err := b.extend(ctx, r, writeSynced); err == nil {
				t.Fatal("time grant did not exceed elapsed lifetime")
			}
			r.Ceilings.WallClockNS = int64(3 * time.Hour)
			granted, err := b.extend(ctx, r, writeSynced)
			if err != nil || granted.StartedAt != s.StartedAt {
				t.Fatalf("absolute time grant: %+v %v", granted, err)
			}
			if _, err := reserveRuntime(ctx, b, "after-time-grant"); err != nil {
				t.Fatalf("extended lifetime unusable: %v", err)
			}
			b.store.now = func() time.Time { return s.StartedAt.Add(4 * time.Hour) }
			if _, err := reserveRuntime(ctx, b, "expired-again"); !errors.Is(err, ErrLimit) {
				t.Fatalf("extension reset the origin: %v", err)
			}
		})
	}
}

func TestRuntimeExtensionLocksSettlementsThroughPublication(t *testing.T) {
	ctx := context.Background()
	root := t.TempDir()
	b := runtimeFixture(t, root, Launch)
	if _, err := reserveRuntime(ctx, b, "first"); err != nil {
		t.Fatal(err)
	}
	c := b.policy.ceilings()
	c.CostMicros = 200
	r := runtimeGrant(t, root, Launch, "cost", c)
	held, release := make(chan struct{}), make(chan struct{})
	result := make(chan error, 1)
	go func() {
		_, err := b.extend(ctx, r, func(path string, data []byte) error { close(held); <-release; return writeSynced(path, data) })
		result <- err
	}()
	select {
	case <-held:
	case err := <-result:
		t.Fatalf("grant failed before publication: %v", err)
	}
	bounded, cancel := context.WithTimeout(ctx, 100*time.Millisecond)
	zero := int64(0)
	_, err := b.store.Settle(bounded, "first", &zero)
	cancel()
	close(release)
	if grantErr := <-result; grantErr != nil {
		t.Fatal(grantErr)
	}
	if !errors.Is(err, ErrLockContention) {
		t.Fatalf("settlement raced the recorded grant exposure: %v", err)
	}
	if _, err := b.store.Settle(ctx, "first", &zero); err != nil {
		t.Fatal(err)
	}
	if _, err := InspectRuntimeBudget(ctx, root, "idea", Launch); err != nil {
		t.Fatalf("later lower actual cost invalidated historical reservation: %v", err)
	}
}

func TestRuntimeExtensionHistoryBoundAndReadOnlyMissingPolicy(t *testing.T) {
	ctx := context.Background()
	root := t.TempDir()
	b := runtimeFixture(t, root, DriverStep)
	s, err := b.inspect(ctx)
	if err != nil {
		t.Fatal(err)
	}
	p := b.policy
	for i := 0; i < maxCycleExtensions; i++ {
		c := p.ceilings()
		c.Actions++
		p = p.withGrant(PolicyExtension{ID: strings.Repeat("a", i+1), At: s.StartedAt, PreviousSHA256: policyDigest(p), Ceilings: c, Spent: 0, ExposureMicros: new(int64), StartedAt: s.StartedAt, Reason: "Explicit fixture decision"})
	}
	if err := validateRuntimePolicy(p); err != nil {
		t.Fatal(err)
	}
	raw, err := json.MarshalIndent(p, "", "  ")
	if err != nil {
		t.Fatal(err)
	}
	if len(raw) <= 16<<10 {
		t.Fatal("fixture did not exceed the old origin reader limit")
	}
	path := filepath.Join(filepath.Dir(b.store.Dir), "policy.json")
	if err := writeSynced(path, raw); err != nil {
		t.Fatal(err)
	}
	current, err := InspectRuntimeBudget(ctx, root, "idea", DriverStep)
	if err != nil {
		t.Fatal(err)
	}
	c := p.ceilings()
	c.Actions++
	if _, err := ExtendRuntimeBudget(ctx, root, "idea", DriverStep, PolicyExtensionRequest{"full", current.PolicySHA256, "Explicit fixture decision", c}); err == nil {
		t.Fatal("history bound ignored")
	}
	absent := t.TempDir()
	if _, err := InspectRuntimeBudget(ctx, absent, "idea", DriverStep); err == nil {
		t.Fatal("missing policy inspected")
	}
	if _, err := os.Stat(filepath.Join(absent, ".parley-runtime")); !os.IsNotExist(err) {
		t.Fatal("read-only inspection initialized state")
	}
}
