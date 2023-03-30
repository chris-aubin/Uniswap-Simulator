package sqrtPriceMath

import (
	"math/big"

	"github.com/chris-aubin/Uniswap-Simulator/src/libraries/constants"
	"github.com/chris-aubin/Uniswap-Simulator/src/libraries/fullMath"
	"github.com/chris-aubin/Uniswap-Simulator/src/libraries/unsafeMath"
)

// Gets the next sqrt price given a delta of token0
// Always rounds up, because in the exact output case (increasing price) we need
// to move the price at least far enough to get the desired output amount, and
// in the exact input case (decreasing price) we need to move the price less in
// order to not send too much output.

// The most precise formula for this is:
//         liquidity * sqrtPX96 / (liquidity +- amount * sqrtPX96),
// If this is impossible because of overflow, Uniswap calculates:
//         liquidity / (liquidity / sqrtPX96 +- amount).

// sqrtPX96 is the starting price, i.e. before accounting for the token0 delta
// liquidity si the amount of usable liquidity
// amount is how much of token0 to add or remove from virtual reserves
// add is whether to add or remove the amount of token0
// Returns the price after adding or removing amount, depending on add
func GetNextSqrtPriceFromAmount0RoundingUp(
	sqrtPX96,
	liquidity,
	amount *big.Int,
	add bool,
) *big.Int {
	// We short circuit amount == 0 because the result is otherwise not guaranteed to equal the input price
	if amount.Cmp(big.NewInt(0)) == 0 {
		return sqrtPX96
	}
	numerator1 := new(big.Int).Lsh(liquidity, 96)

	if add {
		product := new(big.Int).Mul(amount, sqrtPX96)
		if new(big.Int).Div(product, sqrtPX96).Cmp(amount) == 0 {
			denominator := new(big.Int).Add(numerator1, product)
			if denominator.Cmp(numerator1) >= 0 {
				// In this case the result always fits in 160 bits (this is not
				// an issue in Go when using big.Ints, but in the name of
				// accuracy we will simulate this issue).
				return fullMath.MulDivRoundingUp(numerator1, sqrtPX96, denominator)
			}
		}

		denominator := new(big.Int).Div(numerator1, sqrtPX96)
		denominator.Add(denominator, amount)
		return unsafeMath.DivRoundingUp(numerator1, denominator)
	} else {
		product := new(big.Int).Mul(amount, sqrtPX96)
		// If the product overflows, we know the denominator underflows.
		// We must also check that the denominator does not underflow.
		// (this is not an issue in Go when using big.Ints, but in the name of
		// accuracy we will simulate this issue).
		if (new(big.Int).Div(product, sqrtPX96).Cmp(amount) == 0) && (numerator1.Cmp(product) >= 1) {
			denominator := new(big.Int).Sub(numerator1, product)
			return fullMath.MulDivRoundingUp(numerator1, sqrtPX96, denominator)
		} else {
			panic("SQRTPRICE: Overflow")
		}
	}
}

// Gets the next sqrt price given a delta of token1
// Always rounds down, because in the exact output case (decreasing price) we
// need to move the price at least far enough to get the desired output amount,
// and in the exact input case (increasing price) we need to move the price less
// in order to not send too much output. The formula we compute is within <1 wei
// of the lossless version: sqrtPX96 +- amount / liquidity
// sqrtPX96 is the starting price, i.e., before accounting for the token1 delta
// liquidity is the amount of usable liquidity
// amount is how much of token1 to add, or remove, from virtual reserves
// add is a bool that indicates whether to add, or remove, the amount of token1
// Returns the price after adding or removing `amount`
func GetNextSqrtPriceFromAmount1RoundingDown(
	sqrtPX96,
	liquidity,
	amount *big.Int,
	add bool,
) *big.Int {
	// if we're adding (subtracting), rounding down requires rounding the quotient down (up)
	// in both cases, avoid a mulDiv for most inputs
	if add {
		quotient := new(big.Int)
		if amount.Cmp(constants.MaxUint160) >= 0 {
			quotient = quotient.Div(new(big.Int).Lsh(amount, 96), liquidity)
		} else {
			quotient = fullMath.MulDiv(amount, big.NewInt(96), liquidity)
		}

		return new(big.Int).Add(sqrtPX96, quotient)
	} else {
		quotient := new(big.Int)
		if amount.Cmp(constants.MaxUint160) >= 0 {
			quotient = unsafeMath.DivRoundingUp(new(big.Int).Lsh(amount, 96), liquidity)
		} else {
			quotient = fullMath.MulDivRoundingUp(amount, big.NewInt(96), liquidity)
		}

		if sqrtPX96.Cmp(quotient) <= 0 {
			panic("SQRTPRICE: Underflow")
		}
		// always fits 160 bits
		return new(big.Int).Sub(sqrtPX96, quotient)
	}
}

// Gets the next sqrt price given an input amount of token0 or token1
// Panics if price or liquidity are 0, or if the next price is out of bounds
// sqrtPX96 is the starting price, i.e., before accounting for the input amount
// liquidity is the amount of usable liquidity
// amountIn is how much of token0, or token1, is being swapped in
// zeroForOne is a bool that indicates whether the amount in is token0 or token1
// Returns the price after adding the input amount to token0 or token1
func GetNextSqrtPriceFromInput(
	sqrtPX96,
	liquidity,
	amountIn *big.Int,
	zeroForOne bool,
) (sqrtQX96 *big.Int) {
	if (sqrtPX96.Cmp(big.NewInt(0)) <= 0) || (liquidity.Cmp(big.NewInt(0)) <= 0) {
		panic("SQRTPRICE: Invalid input")
	}

	// Round to make sure that we pass the target price
	if zeroForOne {
		return GetNextSqrtPriceFromAmount0RoundingUp(sqrtPX96, liquidity, amountIn, true)
	} else {
		return GetNextSqrtPriceFromAmount1RoundingDown(sqrtPX96, liquidity, amountIn, true)
	}
}

// Gets the next sqrt price given an output amount of token0 or token1
// Panics if price or liquidity are 0 or the next price is out of bounds
// sqrtPX96 is the starting price before accounting for the output amount
// liquidity is the amount of usable liquidity
// amountOut is how much of token0, or token1, is being swapped out
// zeroForOne is a bool that indicates whether the amount out is token0 or token1
// Returns the price after removing the output amount of token0 or token1
func GetNextSqrtPriceFromOutput(
	sqrtPX96,
	liquidity,
	amountOut *big.Int,
	zeroForOne bool,
) (sqrtQX96 *big.Int) {
	if (sqrtPX96.Cmp(big.NewInt(0)) <= 0) || (liquidity.Cmp(big.NewInt(0)) <= 0) {
		panic("SQRTPRICE: Invalid input")
	}

	// Round to make sure that we pass the target price
	if zeroForOne {
		return GetNextSqrtPriceFromAmount0RoundingUp(sqrtPX96, liquidity, amountOut, false)
	} else {
		return GetNextSqrtPriceFromAmount1RoundingDown(sqrtPX96, liquidity, amountOut, false)
	}
}

// Gets the amount0 delta between two prices and calculates:
//         liquidity / sqrt(lower) - liquidity / sqrt(upper)
// i.e. liquidity * (sqrt(upper) - sqrt(lower)) / (sqrt(upper) * sqrt(lower))
// sqrtRatioAX96 is a sqrt price
// sqrtRatioBX96 is another sqrt price
// liquidity is the amount of usable liquidity
// roundUp is a bool that indicates whether to round the amount up or down
// Returns amount0 Amount of token0 required to cover a position of size liquidity between the two passed prices
func GetAmount0Delta(
	sqrtRatioAX96,
	sqrtRatioBX96,
	liquidity *big.Int,
	roundUp bool,
) (amount0 *big.Int) {
	if sqrtRatioAX96.Cmp(sqrtRatioBX96) >= 1 {
		temp := sqrtRatioAX96
		sqrtRatioAX96 = sqrtRatioBX96
		sqrtRatioBX96 = temp
	}

	numerator1 := new(big.Int).Lsh(liquidity, 96)
	numerator2 := new(big.Int).Sub(sqrtRatioBX96, sqrtRatioAX96)

	if sqrtRatioAX96.Cmp(big.NewInt(0)) <= 0 {
		panic("SQRTPRICE: Invalid input")
	}

	if roundUp {
		return unsafeMath.DivRoundingUp(fullMath.MulDivRoundingUp(numerator1, numerator2, sqrtRatioBX96), sqrtRatioAX96)
	} else {
		return new(big.Int).Div(fullMath.MulDiv(numerator1, numerator2, sqrtRatioBX96), sqrtRatioAX96)
	}
}

// Gets the amount1 delta between two prices. Calculates:
//         liquidity * (sqrt(upper) - sqrt(lower))
// sqrtRatioAX96 is a sqrt price
// sqrtRatioBX96 is another sqrt price
// liquidity is the amount of usable liquidity
// roundUp is whether to round the amount up, or down
// Returns amount1, the amount of token1 required to cover a position of size
// liquidity between the two passed prices
func GetAmount1Delta(
	sqrtRatioAX96,
	sqrtRatioBX96,
	liquidity *big.Int,
	roundUp bool,
) (amount1 *big.Int) {
	if sqrtRatioAX96.Cmp(sqrtRatioBX96) >= 1 {
		temp := sqrtRatioAX96
		sqrtRatioAX96 = sqrtRatioBX96
		sqrtRatioBX96 = temp
	}

	if roundUp {
		return fullMath.MulDivRoundingUp(liquidity, new(big.Int).Sub(sqrtRatioBX96, sqrtRatioAX96), constants.Q96)
	} else {
		return fullMath.MulDiv(liquidity, new(big.Int).Sub(sqrtRatioBX96, sqrtRatioAX96), constants.Q96)
	}
}

// Helper that gets signed token0 delta
// sqrtRatioAX96 is a sqrt price
// sqrtRatioBX96 is another sqrt price
// liquidity is the change in liquidity for which to compute the amount0 delta
// Returns amount0, the amount of token0 corresponding to the passed
// liquidityDelta between the two prices.
func getAmount0Delta(
	sqrtRatioAX96,
	sqrtRatioBX96,
	liquidity *big.Int,
) (amount0 *big.Int) {
	if liquidity.Cmp(big.NewInt(0)) < 0 {
		return new(big.Int).Neg(GetAmount0Delta(sqrtRatioAX96, sqrtRatioBX96, new(big.Int).Neg(liquidity), false))
	} else {
		return GetAmount0Delta(sqrtRatioAX96, sqrtRatioBX96, liquidity, true)
	}
}

// Helper that gets signed token1 delta
// sqrtRatioAX96 is a sqrt price
// sqrtRatioBX96 is another sqrt price
// liquidity is the change in liquidity for which to compute the amount1 delta
// Returns amount0, the amount of token1 corresponding to the passed
// liquidityDelta between the two prices.
func getAmount1Delta(
	sqrtRatioAX96,
	sqrtRatioBX96,
	liquidity *big.Int,
) (amount1 *big.Int) {
	if liquidity.Cmp(big.NewInt(0)) < 0 {
		return new(big.Int).Neg(GetAmount1Delta(sqrtRatioAX96, sqrtRatioBX96, new(big.Int).Neg(liquidity), false))
	} else {
		return GetAmount1Delta(sqrtRatioAX96, sqrtRatioBX96, liquidity, true)
	}
}
