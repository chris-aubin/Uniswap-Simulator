// Package liquidityAmounts simulates the Uniswap (periphery) liquidityAmounts library.
//
// In production Uniswap uses the liquidityAmounts library to give users a way
// to compute liquidity amounts from token amounts and prices (it is not part
// of the core Uniswap contracts).
package liquidityAmounts

import (
	"math/big"

	"github.com/chris-aubin/Uniswap-Simulator/src/libraries/constants"
	"github.com/chris-aubin/Uniswap-Simulator/src/libraries/fullMath"
)

// Calculates the amount of liquidity received for a given amount of token0 and
// price range. This is done by calculating:
//     amount0 * (sqrt(upper) * sqrt(lower)) / (sqrt(upper) - sqrt(lower))
//
// Arguments:
// sqrtRatioAX96 -- A sqrt price representing the first tick boundary
// sqrtRatioBX96 -- A sqrt price representing the second tick boundary
// amount0       -- The amount0 being sent in
//
// Returns:
// liquidity     -- The amount of returned liquidity
func getLiquidityForAmount0(sqrtRatioAX96, sqrtRatioBX96, amount0 *big.Int) *big.Int {
	if sqrtRatioAX96.Cmp(sqrtRatioBX96) >= 1 {
		sqrtRatioAX96, sqrtRatioBX96 = sqrtRatioBX96, sqrtRatioAX96
	}
	intermediate := fullMath.MulDiv(sqrtRatioAX96, sqrtRatioBX96, constants.Q96)
	return fullMath.MulDiv(amount0, intermediate, new(big.Int).Sub(sqrtRatioBX96, sqrtRatioAX96))
}

// Calculates the amount of liquidity received for a given amount of token1 and
// price range. This is done by calculating:
//     amount1 / (sqrt(upper) - sqrt(lower)).
//
// Arguments:
// sqrtRatioAX96 -- A sqrt price representing the first tick boundary
// sqrtRatioBX96 -- A sqrt price representing the second tick boundary
// amount1       -- The amount1 being sent in
//
// Returns:
// liquidity     -- The amount of returned liquidity
func getLiquidityForAmount1(sqrtRatioAX96, sqrtRatioBX96, amount1 *big.Int) *big.Int {
	if sqrtRatioAX96.Cmp(sqrtRatioBX96) >= 1 {
		sqrtRatioAX96, sqrtRatioBX96 = sqrtRatioBX96, sqrtRatioAX96
	}
	return fullMath.MulDiv(amount1, constants.Q96, new(big.Int).Sub(sqrtRatioBX96, sqrtRatioAX96))
}

// Calculates the maximum amount of liquidity received for a given amount of
// token0, token1, the current pool prices and the prices at the tick boundaries.
//
// Arguments:
// sqrtRatioX96  -- A sqrt price representing the current pool prices
// sqrtRatioAX96 -- A sqrt price representing the first tick boundary
// sqrtRatioBX96 -- A sqrt price representing the second tick boundary
// amount0       -- The amount of token0 being sent in
// amount1       -- The amount of token1 being sent in
//
// Returns:
// liquidity     -- The maximum amount of liquidity received
func GetLiquidityForAmounts(sqrtRatioX96, sqrtRatioAX96, sqrtRatioBX96, amount0, amount1 *big.Int) (liquidity *big.Int) {
	if sqrtRatioAX96.Cmp(sqrtRatioBX96) >= 1 {
		sqrtRatioAX96, sqrtRatioBX96 = sqrtRatioBX96, sqrtRatioAX96
	}
	if sqrtRatioX96.Cmp(sqrtRatioAX96) <= 0 {
		liquidity = getLiquidityForAmount0(sqrtRatioAX96, sqrtRatioBX96, amount0)
	} else if sqrtRatioX96.Cmp(sqrtRatioBX96) <= -1 {
		liquidity0 := getLiquidityForAmount0(sqrtRatioX96, sqrtRatioBX96, amount0)
		liquidity1 := getLiquidityForAmount1(sqrtRatioAX96, sqrtRatioX96, amount1)
		if liquidity0.Cmp(liquidity1) <= -1 {
			liquidity = liquidity0
		} else {
			liquidity = liquidity1
		}
	} else {
		liquidity = getLiquidityForAmount1(sqrtRatioAX96, sqrtRatioBX96, amount1)
	}
	return liquidity
}

/// @notice Computes the amount of token0 for a given amount of liquidity and a price range
/// @param sqrtRatioAX96 A sqrt price representing the first tick boundary
/// @param sqrtRatioBX96 A sqrt price representing the second tick boundary
/// @param liquidity The liquidity being valued
/// @return amount0 The amount of token0
func getAmount0ForLiquidity(sqrtRatioAX96, sqrtRatioBX96, liquidity *big.Int) *big.Int {
	if sqrtRatioAX96.Cmp(sqrtRatioBX96) >= 1 {
		sqrtRatioAX96, sqrtRatioBX96 = sqrtRatioBX96, sqrtRatioAX96
	}

	temp := fullMath.MulDiv(new(big.Int).Lsh(liquidity, 96), new(big.Int).Sub(sqrtRatioBX96, sqrtRatioAX96), sqrtRatioBX96)
	return new(big.Int).Div(temp, sqrtRatioAX96)
}

// Calculates the amount of token1 for a given amount of liquidity and a price
// range.
//
// Arguments:
// sqrtRatioAX96 -- A sqrt price representing the first tick boundary
// sqrtRatioBX96 -- A sqrt price representing the second tick boundary
// liquidity     -- The liquidity being valued
//
// Returns:
// amount1       -- The amount of token1
func getAmount1ForLiquidity(sqrtRatioAX96, sqrtRatioBX96, liquidity *big.Int) *big.Int {
	if sqrtRatioAX96.Cmp(sqrtRatioBX96) >= 1 {
		sqrtRatioAX96, sqrtRatioBX96 = sqrtRatioBX96, sqrtRatioAX96
	}

	return fullMath.MulDiv(liquidity, new(big.Int).Sub(sqrtRatioBX96, sqrtRatioAX96), constants.Q96)
}

// Calculates the amount of token0 and token1 for a given amount of liquidity
// and a price range.
//
// Arguments:
// sqrtRatioX96  -- A sqrt price representing the current pool prices
// sqrtRatioAX96 -- A sqrt price representing the first tick boundary
// sqrtRatioBX96 -- A sqrt price representing the second tick boundary
// liquidity     -- The liquidity being valued
//
// Returns:
// amount0       -- The amount of token0
// amount1       -- The amount of token1
func GetAmountsForLiquidity(sqrtRatioX96, sqrtRatioAX96, sqrtRatioBX96, liquidity *big.Int) (*big.Int, *big.Int) {
	if sqrtRatioAX96.Cmp(sqrtRatioBX96) >= 1 {
		sqrtRatioAX96, sqrtRatioBX96 = sqrtRatioBX96, sqrtRatioAX96
	}
	amount0, amount1 := big.NewInt(0), big.NewInt(0)
	if sqrtRatioX96.Cmp(sqrtRatioAX96) <= 0 {
		amount0 = getAmount0ForLiquidity(sqrtRatioAX96, sqrtRatioBX96, liquidity)
	} else if sqrtRatioX96.Cmp(sqrtRatioBX96) <= -1 {
		amount0 = getAmount0ForLiquidity(sqrtRatioX96, sqrtRatioBX96, liquidity)
		amount1 = getAmount1ForLiquidity(sqrtRatioAX96, sqrtRatioX96, liquidity)
	} else {
		amount1 = getAmount1ForLiquidity(sqrtRatioAX96, sqrtRatioBX96, liquidity)
	}
	return amount0, amount1
}
