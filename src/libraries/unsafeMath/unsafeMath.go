// Package unsafeMath simulates the Uniswap unsafeMath library.
//
// In production Uniswap uses these functions that do not check inputs or
// outputs only when overflow is unavoidable. These functions simulate overflow,
// because when using the Go big.Int library, overflow is not an issue.
package unsafeMath

import (
	"math/big"

	"github.com/chris-aubin/Uniswap-Simulator/src/libraries/constants"
)

// Calculates ceil(numerator/denominator) for *big.Ints
func DivRoundingUp(numerator, denominator *big.Int) *big.Int {
	remainder := big.NewInt(0)
	quotient, remainder := new(big.Int).DivMod(numerator, denominator, remainder)
	if remainder.Cmp(big.NewInt(0)) == 1 {
		quotient.Add(quotient, big.NewInt(1))
	}

	// Check whether result could fit in a uint256, if not try to simulate
	// overflow.
	if quotient.Cmp(constants.MaxUint256) >= 1 {
		return new(big.Int).Mod(quotient, constants.MaxUint256)
	}
	return quotient
}
