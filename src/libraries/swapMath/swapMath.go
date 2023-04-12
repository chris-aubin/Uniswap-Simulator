package swapMath

import (
	"math/big"

	"github.com/chris-aubin/Uniswap-Simulator/src/libraries/fullMath"
	"github.com/chris-aubin/Uniswap-Simulator/src/libraries/sqrtPriceMath"
)

// Computes the result of swapping some amount in, or amount out, given the
// parameters of the swap
// The fee, plus the amount in, will never exceed the amount remaining if the
// swap's `amountSpecified` is positive
// sqrtRatioCurrentX96 is the current sqrt price of the pool
// sqrtRatioTargetX96 is the price that cannot be exceeded, from which the
// direction of the swap is inferred
// liquidity is the usable liquidity
// amountRemaining is how much input or output amount is remaining to be swapped
// in/out
// feePips is the fee taken from the input amount, expressed in hundredths of a bip
// Returns sqrtRatioNextX96, the price after swapping the amount in/out, not to
// exceed the price target
// Returns amountIn, the amount to be swapped in, of either token0 or token1,
// based on the direction of the swap
// Returns amountOut, the amount to be received, of either token0 or token1,
// based on the direction of the swap
// Returns feeAmount, the amount of input that will be taken as a fee
func ComputeSwapStep(sqrtRatioCurrentX96, sqrtRatioTargetX96, liquidity, amountRemaining *big.Int, feePips int) (sqrtRatioNextX96, amountIn, amountOut, feeAmount *big.Int) {
	var zeroForOne bool
	if sqrtRatioCurrentX96.Cmp(sqrtRatioTargetX96) >= 0 {
		zeroForOne = true
	} else {
		zeroForOne = false
	}
	// bool exactIn = amountRemaining >= 0;
	var exactIn bool
	if amountRemaining.Cmp(big.NewInt(0)) >= 0 {
		exactIn = true
	} else {
		exactIn = false
	}

	if exactIn {
		amountRemainingLessFee := fullMath.MulDiv(amountRemaining, big.NewInt(int64(1000000-feePips)), big.NewInt(1000000))
		if zeroForOne {
			amountIn = sqrtPriceMath.GetAmount0Delta(sqrtRatioTargetX96, sqrtRatioCurrentX96, liquidity, true)
		} else {
			amountIn = sqrtPriceMath.GetAmount1Delta(sqrtRatioCurrentX96, sqrtRatioTargetX96, liquidity, true)
		}
		// if (amountRemainingLessFee >= amountIn) sqrtRatioNextX96 = sqrtRatioTargetX96;
		if amountRemainingLessFee.Cmp(amountIn) >= 0 {
			sqrtRatioNextX96 = sqrtRatioTargetX96
		} else {
			sqrtRatioNextX96 = sqrtPriceMath.GetNextSqrtPriceFromInput(
				sqrtRatioCurrentX96,
				liquidity,
				amountRemainingLessFee,
				zeroForOne,
			)
		}
	} else {
		if zeroForOne {
			amountOut = sqrtPriceMath.GetAmount1Delta(sqrtRatioTargetX96, sqrtRatioCurrentX96, liquidity, false)
		} else {
			amountOut = sqrtPriceMath.GetAmount0Delta(sqrtRatioCurrentX96, sqrtRatioTargetX96, liquidity, false)
		}
		if new(big.Int).Neg(amountRemaining).Cmp(amountOut) >= 0 {
			sqrtRatioNextX96 = sqrtRatioTargetX96
		} else {
			sqrtRatioNextX96 = sqrtPriceMath.GetNextSqrtPriceFromOutput(
				sqrtRatioCurrentX96,
				liquidity,
				new(big.Int).Neg(amountRemaining),
				zeroForOne,
			)
		}
	}

	var max bool
	if sqrtRatioTargetX96.Cmp(sqrtRatioNextX96) == 0 {
		max = true
	} else {
		max = false
	}

	// Get the input/output amounts
	if zeroForOne {
		if !(max && exactIn) {
			amountIn = sqrtPriceMath.GetAmount0Delta(sqrtRatioNextX96, sqrtRatioCurrentX96, liquidity, true)
		}
		if !(max && !exactIn) {
			amountOut = sqrtPriceMath.GetAmount1Delta(sqrtRatioNextX96, sqrtRatioCurrentX96, liquidity, false)
		}
	} else {
		if !(max && exactIn) {
			amountIn = sqrtPriceMath.GetAmount1Delta(sqrtRatioCurrentX96, sqrtRatioNextX96, liquidity, true)
		}
		if !(max && !exactIn) {
			amountOut = sqrtPriceMath.GetAmount0Delta(sqrtRatioCurrentX96, sqrtRatioNextX96, liquidity, false)
		}
	}

	// Cap the output amount to not exceed the remaining output amount
	if !exactIn && amountOut.Cmp(new(big.Int).Neg(amountRemaining)) >= 1 {
		amountOut = new(big.Int).Neg(amountRemaining)
	}
	// if (exactIn && sqrtRatioNextX96 != sqrtRatioTargetX96) {
	if exactIn && sqrtRatioNextX96.Cmp(sqrtRatioTargetX96) != 0 {
		// We didn't reach the target, so take the remainder of the maximum input as fee
		feeAmount = new(big.Int).Sub(amountRemaining, amountIn)
	} else {
		feeAmount = fullMath.MulDivRoundingUp(amountIn, big.NewInt(int64(feePips)), big.NewInt(int64(1e6-feePips)))
	}
	return
}
