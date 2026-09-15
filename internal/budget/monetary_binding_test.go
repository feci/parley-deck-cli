package budget

import (
	"context"
	"errors"
	"math"
	"os"
	"path/filepath"
	"testing"
)

func TestUSDLimitMicrosConservativeDecimalBoundaries(t *testing.T) {
	for _, tc := range []struct {
		usd  float64
		want int64
	}{
		{0, 0}, {0.000001, 1}, {0.0000019, 1}, {0.0000099, 9},
		{0.29, 290000}, {12.5, 12500000}, {10.1234569, 10123456},
		{9223372036854.773, 9223372036854773000},
	} {
		got, err := USDLimitMicros(tc.usd)
		if err != nil || got != tc.want {
			t.Fatalf("%.17g: got %d, %v; want %d", tc.usd, got, err, tc.want)
		}
	}
	for _, usd := range []float64{-1, math.NaN(), math.Inf(1), math.Inf(-1), 0.0000009, math.SmallestNonzeroFloat64, math.MaxFloat64, 9223372036854.78} {
		if _, err := USDLimitMicros(usd); err == nil {
			t.Fatalf("invalid ceiling %g became a grant", usd)
		}
	}
}

func TestRequireMonetaryBindingDoesNotInventOrChangeGrant(t *testing.T) {
	if err := RequireMonetaryBinding(nil, 1); !errors.Is(err, ErrUnknownCost) {
		t.Fatalf("unreserved policy: %v", err)
	}
	if err := RequireMonetaryBinding(nil, 0); err != nil {
		t.Fatal(err)
	}
	ctx := context.Background()
	reserve := int64(600000)
	b, err := ConfigureLaunchBudget(ctx, t.TempDir(), "fixture", LaunchPolicy{MaxCostMicros: 1000000, ReserveMicros: &reserve})
	if err != nil {
		t.Fatal(err)
	}
	path := filepath.Join(filepath.Dir(b.Store.Dir), "policy.json")
	before, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	for _, usd := range []float64{1, 0} {
		if err := RequireMonetaryBinding(b, usd); err != nil {
			t.Fatalf("matching/omitted ceiling: %v", err)
		}
	}
	for _, usd := range []float64{0.5, 2, -1, math.NaN()} {
		if err := RequireMonetaryBinding(b, usd); err == nil {
			t.Fatalf("changed/invalid ceiling accepted: %g", usd)
		}
	}
	after, err := os.ReadFile(path)
	if err != nil || string(after) != string(before) {
		t.Fatal("default mapping mutated the frozen policy")
	}
	b.Policy.ReserveMicros = nil
	if err := RequireMonetaryBinding(b, 1); err == nil {
		t.Fatal("missing reservation became an implicit price")
	}
}
