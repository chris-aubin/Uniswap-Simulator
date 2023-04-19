// Package tickMath simulates the Uniswap TickMath library.
// 	
// 	Contains methods to compute the sqrt price for ticks of size 1.0001, 
// i.e. sqrt(1.0001^tick) as fixed point Q64.96 numbers and the tick for a given
// sqrt price.
package tickMath

import (
	"math/big"

	"github.com/chris-aubin/Uniswap-Simulator/src/libraries/constants"
)

// Computes  the sqrt price for ticks of size 1.0001, i.e. sqrt(1.0001^tick)
// Supports prices between 2**-128 and 2**128.
// 
// Arguments:
// tick -- The tick for which to compute the sqrt price
//
// Returns:
// The sqrt price as a fixed point Q64.96 number
func GetSqrtRatioAtTick(tick int) *big.Int {
	absTick := tick
	if tick < 0 {
		absTick = -tick
	}

	if absTick > constants.MaxTick {
		panic("tickMath.getSqrtRatioAtTick: INVALID_TICK")
	}

	// The tick is processed one bit at a time (0x1 = 1, 0x2 = 10, 0x4 = 100, etc.)
	// Each of these magic values is sqrt(1/1.0001)^2**bit, i.e. the square root
	// of the price ratio. The magic values are used to compute
	// sqrt(1.0001^tick) by using the fact that x^c = x^a*x^b where a + b = c.
	// ratio := new(big.Float)
	ratio := new(big.Int)
	if absTick&0x1 != 0 {
		// 0xfffcb933bd6fad37aa2d162d1a594001 == sqrt(1/1.0001)^1
		// Can check by comparing
		// int('0xfffcb933bd6fad37aa2d162d1a594001', 16)/2**128
		// to
		// math.sqrt(1/1.0001)**0

		ratio.SetString("fffcb933bd6fad37aa2d162d1a594001", 16)
	} else {
		// 0x100000000000000000000000000000000 == sqrt(1/1.0001)^0
		ratio.SetString("100000000000000000000000000000000", 16)
	}
	if (absTick & 0x2) != 0 {
		// 0xfff97272373d413259a46990580e213a == sqrt(1/1.0001)^2
		ratio = mulShift(ratio, "fff97272373d413259a46990580e213a")
	}
	if (absTick & 0x4) != 0 {
		// 0xfff2e50f5f656932ef12357cf3c7fdcc == sqrt(1/1.0001)^4
		ratio = mulShift(ratio, "fff2e50f5f656932ef12357cf3c7fdcc")
	}
	if (absTick & 0x8) != 0 {
		// 0xfff2e50f5f656932ef12357cf3c7fdcc == sqrt(1/1.0001)^8
		ratio = mulShift(ratio, "ffe5caca7e10e4e61c3624eaa0941cd0")
	}
	if (absTick & 0x10) != 0 {
		ratio = mulShift(ratio, "ffcb9843d60f6159c9db58835c926644")
	}
	if (absTick & 0x20) != 0 {
		ratio = mulShift(ratio, "ff973b41fa98c081472e6896dfb254c0")
	}
	if (absTick & 0x40) != 0 {
		ratio = mulShift(ratio, "ff2ea16466c96a3843ec78b326b52861")
	}
	if (absTick & 0x80) != 0 {
		ratio = mulShift(ratio, "fe5dee046a99a2a811c461f1969c3053")
	}
	if (absTick & 0x100) != 0 {
		ratio = mulShift(ratio, "fcbe86c7900a88aedcffc83b479aa3a4")
	}
	if (absTick & 0x200) != 0 {
		ratio = mulShift(ratio, "f987a7253ac413176f2b074cf7815e54")
	}
	if (absTick & 0x400) != 0 {
		ratio = mulShift(ratio, "f3392b0822b70005940c7a398e4b70f3")
	}
	if (absTick & 0x800) != 0 {
		ratio = mulShift(ratio, "e7159475a2c29b7443b29c7fa6e889d9")
	}
	if (absTick & 0x1000) != 0 {
		ratio = mulShift(ratio, "d097f3bdfd2022b8845ad8f792aa5825")
	}
	if (absTick & 0x2000) != 0 {
		ratio = mulShift(ratio, "a9f746462d870fdf8a65dc1f90e061e5")
	}
	if (absTick & 0x4000) != 0 {
		ratio = mulShift(ratio, "70d869a156d2a1b890bb3df62baf32f7")
	}
	if (absTick & 0x8000) != 0 {
		ratio = mulShift(ratio, "31be135f97d08fd981231505542fcfa6")
	}
	if (absTick & 0x10000) != 0 {
		ratio = mulShift(ratio, "9aa508b5b7a84e1c677de54f3e99bc9")
	}
	if (absTick & 0x20000) != 0 {
		ratio = mulShift(ratio, "5d6af8dedb81196699c329225ee604")
	}
	if (absTick & 0x40000) != 0 {
		ratio = mulShift(ratio, "2216e584f5fa1ea926041bedfe98")
	}
	if (absTick & 0x80000) != 0 {
		ratio = mulShift(ratio, "48a170391f7dc42444e8fa2")
	}
	// Because each of these magic values is sqrt(1/1.0001)^2**bit, if tick > 0
	// we need to invert the ratio (calculate 1/ratio). When using Q128.128
	// fixed point numbers we do this by dividing by 1<<128.
	if tick > 0 {
		ratio = new(big.Int).Div(constants.MaxUint256, ratio)
	}

	// This divides by 1<<32 to convert from Q128.128 to Q64.96
	// We then downcast because we know the result always fits within 160 bits due to our tick input constraint
	// We round up in the division so getTickAtSqrtRatio of the output price is always consistent
	rounding := big.NewInt(0)
	if new(big.Int).Mod(ratio, big.NewInt(1<<32)).Cmp(big.NewInt(0)) != 0 {
		rounding = big.NewInt(1)
	}
	sqrtPriceX96 := new(big.Int).Add(new(big.Int).Rsh(ratio, 32), rounding)
	return sqrtPriceX96
}

// Calculates the greatest tick value such that getRatioAtTick(tick) <= ratio.
// 
// Arguments:
// sqrtPriceX96 -- The sqrt ratio for which to compute the tick
//
// Returns:
// The greatest tick for which the ratio is less than or equal to the input ratio
func GetTickAtSqrtRatio(sqrtPriceX96 *big.Int) int {
	if (sqrtPriceX96.Cmp(constants.MaxSqrtRatio) != -1) || (sqrtPriceX96.Cmp(constants.MinSqrtRatioBig) != 1) {
		panic("tickMath.getTickAtSqrtRatio: INVALID_SQRT_RATIO")
	}

	ratio := new(big.Int).Lsh(sqrtPriceX96, 32)

	// Find the most significant bit (msb is an approximation of log2(ratio))
	r := new(big.Int).Lsh(sqrtPriceX96, 32)
	msb := big.NewInt(0)
	for i := 7; i > 0; i-- {
		cmp := r.Cmp(constants.MaxUints[i])
		if cmp <= -1 {
			cmp = 0
		}
		f := new(big.Int).Lsh(big.NewInt(int64(cmp)), uint(i))
		msb = new(big.Int).Or(msb, f)
		r = new(big.Int).Rsh(r, uint(f.Int64()))
	}
	// Calculates l0 in the iterative approximation algorithm described in the
	// README.
	cmp := r.Cmp(constants.MaxUints[0])
	if cmp == -1 {
		cmp = 0
	}
	f := new(big.Int).Lsh(big.NewInt(int64(cmp)), 0)
	msb = new(big.Int).Or(msb, f)
	dif := new(big.Int).Sub(msb, big.NewInt(127))
	if msb.Cmp(big.NewInt(128)) >= 0 {
		r = new(big.Int).Rsh(ratio, uint(dif.Int64()))
	} else {
		r = new(big.Int).Lsh(ratio, uint(dif.Int64()))
	}

	log_2_temp := new(big.Int).Sub(msb, big.NewInt(128))
	log_2 := new(big.Int).Lsh(log_2_temp, 64)

	// Iteratively approximates a bit after the fixed point.
	for i := 0; i < 14; i++ {
		r = new(big.Int).Rsh(new(big.Int).Mul(r, r), 127)
		f := new(big.Int).Rsh(r, 128)
		log_2 = new(big.Int).Or(log_2, new(big.Int).Lsh(f, uint(63-i)))
		r = new(big.Int).Rsh(r, uint(f.Uint64()))
	}

	// Change of base formula to convert from log_2 to log_sqrt(1.0001)
	log_sqrt10001_multiplicand, _ := new(big.Int).SetString("255738958999603826347141", 10)
	log_sqrt10001 := new(big.Int).Mul(log_2, log_sqrt10001_multiplicand)

	// Adjust for the absolute error of the approximation to ensure that the
	// tick is always less than or equal to the input ratio.
	tickLow_sub, _ := new(big.Int).SetString("3402992956809132418596140100660247210", 10)
	tickLow := new(big.Int).Rsh(new(big.Int).Sub(log_sqrt10001, tickLow_sub), 128)

	tickHigh_add, _ := new(big.Int).SetString("291339464771989622907027621153398088495", 10)
	tickHigh := new(big.Int).Rsh(new(big.Int).Add(log_sqrt10001, tickHigh_add), 128)

	if tickLow == tickHigh {
		return int(tickLow.Int64())
	}

	sqrtRatio := GetSqrtRatioAtTick(int(tickHigh.Int64()))
	if sqrtRatio.Cmp(sqrtPriceX96) <= 0 {
		return int(tickHigh.Int64())
	} else {
		return int(tickLow.Int64())
	}
}

// Helper function to multiply two Q128.128 fixed point numbers
func mulShift(multiplier *big.Int, multiplicand string) *big.Int {
	multiplicandBig, _ := new(big.Int).SetString(multiplicand, 16)
	productBig := new(big.Int).Mul(multiplier, multiplicandBig)
	result := new(big.Int).Rsh(productBig, 128)
	return result
}
