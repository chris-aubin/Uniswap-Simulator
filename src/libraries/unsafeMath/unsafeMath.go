package unsafeMath

// In production Uniswap uses these functions that do not check inputs or
// outputs only when overflow is unavoidable. Need to adjust these functions to
// simulate overflow, because when using the Go big.Int library, overflow is
// not an issue.

import (
	"math/big"
)

// Calculates ceil(numerator/denominator) for *big.Ints
// numerator is the numerator
// denominator is the denominator
func DivRoundingUp(
	numerator,
	denominator *big.Int,
) *big.Int {
	remainder := big.NewInt(0)
	quotient, remainder := new(big.Int).DivMod(numerator, denominator, remainder)
	if remainder.Cmp(big.NewInt(0)) == 1 {
		quotient.Add(quotient, big.NewInt(1))
	}
	return quotient
}
