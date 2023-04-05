// Package fullMath simulates the Uniswap fullMath library.
//
// In production Uniswap uses the fullMath library to perform
// floor(a×b÷denominator) and ceil(a×b÷denominator) in such a way that overflow
// of an intermediate value is handled without any loss of precision. The
// functions compute the product mod 2**256 and mod 2**256 - 1 and then use the
// Chinese Remainder Theorem to reconstruct the 512 bit result. The result is
// stored in two 256 variables such that product = prod1 * 2**256 + prod0. This
// package uses the Go big.Int library to simulate the fullMath library.
package fullMath

import (
	"math/big"

	"github.com/chris-aubin/Uniswap-Simulator/src/libraries/constants"
)

// Calculates floor(a×b÷denominator) for *big.Ints. Panics if any of the
// arguments are negative, if the result would overflow a uint256 or if
// denominator == 0.
//
// multiplicand is the multiplicand.
// multiplier is the multiplier.
// denominator is the denominator.
func MulDiv(
	multiplicand,
	multiplier,
	denominator *big.Int,
) (result *big.Int) {
	// Check arguments
	checkArgs(multiplicand, multiplier, denominator)

	product := new(big.Int).Mul(multiplicand, multiplier)
	result = new(big.Int).Div(product, denominator)

	// Check whether result could fit in a uint256
	if result.Cmp(constants.MaxUint256) >= 1 {
		panic("fullMath.MulDiv: Overflow")
	}
	return result
}

// Calculates ceil(a×b÷denominator) for *big.Ints. Panics if any of the
// arguments are negative, if the result would overflow a uint256 or if
// denominator == 0.
//
// multiplicand is the multiplicand.
// multiplier is the multiplier.
// denominator is the denominator.
func MulDivRoundingUp(
	multiplicand,
	multiplier,
	denominator *big.Int,
) (result *big.Int) {
	// Check arguments
	checkArgs(multiplicand, multiplier, denominator)

	product := new(big.Int).Mul(multiplicand, multiplier)
	remainder := big.NewInt(0)
	quotient, remainder := new(big.Int).DivMod(product, denominator, remainder)
	if remainder.Cmp(big.NewInt(0)) == 1 {
		quotient.Add(quotient, big.NewInt(1))
	}

	// Check whether result could fit in a uint256
	if quotient.Cmp(constants.MaxUint256) >= 1 {
		panic("fullMath.MulDivRoundingUp: Overflow")
	}
	return quotient
}

// Checks whether any of the arguments are negative or if the denominator is
// zero. Panics if any of these conditions are true.
func checkArgs(
	multiplicand,
	multiplier,
	denominator *big.Int,
) {
	// Check for division by zero
	if denominator.Cmp(big.NewInt(0)) == 0 {
		panic("fullMath: Division by zero")
	}

	// Check for negative arguments (not supported by Uniswap)
	if multiplicand.Cmp(big.NewInt(0)) <= -1 {
		panic("fullMath: multiplicand < 0")
	}
	if multiplier.Cmp(big.NewInt(0)) <= -1 {
		panic("fullMath: multiplier < 0")
	}
	if denominator.Cmp(big.NewInt(0)) <= -1 {
		panic("fullMath: denominator < 0")
	}
}
