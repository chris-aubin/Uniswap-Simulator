// Package liquidityMath simulates the Uniswap liquidityMath library.
//
// In production Uniswap uses the liquidityMath library to calculate the change
// in liquidity when a user adds or removes liquidity in such a way as to
// detect overflow and underflow.
package liquidityMath

import (
	"math/big"

	"github.com/chris-aubin/Uniswap-Simulator/src/libraries/constants"
)

// Adds a signed liquidity delta to liquidity and reverts if it overflows or
// underflows.
func AddDelta(x, y *big.Int) (delta *big.Int) {
	// Ensure that x > 0
	if x.Cmp(big.NewInt(0)) <= -1 {
		panic("liquidityMath.AddDelta: x must be greater than 0")
	}

	delta = new(big.Int).Add(x, y)

	// Check whether result could fit in a uint128
	if delta.Cmp(constants.MaxUint128) >= 1 {
		panic("liquidityMath.AddDelta: Overflow")
	}
	if delta.Cmp(big.NewInt(0)) <= -1 {
		panic("liquidityMath.AddDelta: Underflow")
	}
	return delta
}
