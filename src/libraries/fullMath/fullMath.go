package fullMath

// In production Uniswap uses these function to handle multiplication and
// division that can have overflow of an intermediate value without any loss of
// precision (where an intermediate value overflows 256 bits).

import (
	"math/big"
)

// Calculates floor(a×b÷denominator) for *big.Ints
// In production this function computes the product mod 2**256 and mod
// 2**256 - 1 and then uses the Chinese Remainder Theorem to reconstruct
// the 512 bit result. The result is stored in two 256 variables such that
// product = prod1 * 2**256 + prod0. We just use the Go big.Int library.
// multiplicand is the multiplicand
// multiplier is the multiplier
// denominator is the denominator
func MulDiv(
	multiplicand,
	multiplier,
	denominator *big.Int,
) *big.Int {
	product := new(big.Int).Mul(multiplicand, multiplier)
	quotient := new(big.Int).Div(product, denominator)
	return quotient
}

// Calculates ceil(a×b÷denominator) for *big.Ints
// multiplicand is the multiplicand
// multiplier is the multiplier
// denominator is the denominator
func MulDivRoundingUp(
	multiplicand,
	multiplier,
	denominator *big.Int,
) *big.Int {
	product := new(big.Int).Mul(multiplicand, multiplier)
	remainder := big.NewInt(0)
	quotient, remainder := new(big.Int).DivMod(product, denominator, remainder)
	if remainder.Cmp(big.NewInt(0)) == 1 {
		quotient.Add(quotient, big.NewInt(1))
	}
	return quotient
}
