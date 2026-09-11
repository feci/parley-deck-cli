package budget

import (
	"errors"
	"fmt"
	"math"
	"math/big"
	"strconv"
)

// USDLimitMicros maps a configured dollar ceiling to whole microdollars without
// increasing it. Use its shortest decimal representation rather than binary
// multiplication, which can move an exact decimal boundary by one microdollar.
// Zero means no new requested ceiling; it never disables a persisted policy.
func USDLimitMicros(usd float64) (int64, error) {
	if math.IsNaN(usd) || math.IsInf(usd, 0) || usd < 0 {
		return 0, errors.New("configured monetary ceiling must be finite and nonnegative")
	}
	if usd == 0 {
		return 0, nil
	}
	r, ok := new(big.Rat).SetString(strconv.FormatFloat(usd, 'g', -1, 64))
	if !ok {
		return 0, errors.New("cannot represent configured monetary ceiling")
	}
	r.Mul(r, big.NewRat(1_000_000, 1))
	n := new(big.Int).Quo(r.Num(), r.Denom())
	if !n.IsInt64() || n.Sign() <= 0 {
		return 0, errors.New("configured monetary ceiling is below one microdollar or overflows microdollars")
	}
	return n.Int64(), nil
}

// RequireMonetaryBinding maps loop configuration to the launch boundary. A
// dollar limit is not a per-launch price: without an existing finite policy and
// explicit conservative reservation, stop instead of inventing one or allowing
// a first unbounded call. Configuration never changes a frozen operator policy.
// The caller loads the binding for the actual shared idea/live-origin scope.
func RequireMonetaryBinding(b *LaunchBinding, usd float64) error {
	micros, err := USDLimitMicros(usd)
	if err != nil || micros == 0 {
		return err
	}
	if b == nil {
		return fmt.Errorf("%w: configured max_cost_usd requires an attended budget configure policy with max-cost-micros %d and an explicit conservative reserve-micros before execution", ErrUnknownCost, micros)
	}
	if err := b.Policy.validate(); err != nil {
		return err
	}
	if b.Policy.MaxCostMicros != micros {
		return errors.New("configured monetary ceiling differs from the frozen launch policy; configuration is not an operator extension or migration")
	}
	return nil
}
