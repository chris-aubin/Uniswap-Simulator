package tickMath

import (
	"math/big"

	"github.com/chris-aubin/Uniswap-Simulator/src/libraries/constants"
)

// Computes  the sqrt price for ticks of size 1.0001, i.e. sqrt(1.0001^tick)
// Supports prices between 2**-128 and 2**128
func getSqrtRatioAtTick(tick int) *big.Int {
	absTick := tick
	if tick < 0 {
		absTick = -tick
	}

	if absTick > constants.MaxTick {
		panic("tickMath.getSqrtRatioAtTick: INVALID_TICK")
	}

	// The tick is processed one bit at a time (0x1 = 1, 0x2 = 10, 0x4 = 100, etc.)
	ratio := new(big.Int)
	if absTick&0x1 != 0 {
		// 0xfffcb933bd6fad37aa2d162d1a594001 = 340265354078544963557816517032075149313
		ratio.SetString("0xfffcb933bd6fad37aa2d162d1a594001", 16)
	} else {
		ratio.SetString("0x100000000000000000000000000000000", 16)
	}
	if (absTick & 0x2) != 0 {
		ratio = mulShift(ratio, "0xfff97272373d413259a46990580e213a")
	}
	if (absTick & 0x4) != 0 {
		ratio = mulShift(ratio, "0xfff2e50f5f656932ef12357cf3c7fdcc")
	}
	if (absTick & 0x8) != 0 {
		ratio = mulShift(ratio, "0xffe5caca7e10e4e61c3624eaa0941cd0")
	}
	if (absTick & 0x10) != 0 {
		ratio = mulShift(ratio, "0xffcb9843d60f6159c9db58835c926644")
	}
	if (absTick & 0x20) != 0 {
		ratio = mulShift(ratio, "0xff973b41fa98c081472e6896dfb254c0")
	}
	if (absTick & 0x40) != 0 {
		ratio = mulShift(ratio, "0xff2ea16466c96a3843ec78b326b52861")
	}
	if (absTick & 0x80) != 0 {
		ratio = mulShift(ratio, "0xfe5dee046a99a2a811c461f1969c3053")
	}
	if (absTick & 0x100) != 0 {
		ratio = mulShift(ratio, "0xfcbe86c7900a88aedcffc83b479aa3a4")
	}
	if (absTick & 0x200) != 0 {
		ratio = mulShift(ratio, "0xf987a7253ac413176f2b074cf7815e54")
	}
	if (absTick & 0x400) != 0 {
		ratio = mulShift(ratio, "0xf3392b0822b70005940c7a398e4b70f3")
	}
	if (absTick & 0x800) != 0 {
		ratio = mulShift(ratio, "0xe7159475a2c29b7443b29c7fa6e889d9")
	}
	if (absTick & 0x1000) != 0 {
		ratio = mulShift(ratio, "0xd097f3bdfd2022b8845ad8f792aa5825")
	}
	if (absTick & 0x2000) != 0 {
		ratio = mulShift(ratio, "0xa9f746462d870fdf8a65dc1f90e061e5")
	}
	if (absTick & 0x4000) != 0 {
		ratio = mulShift(ratio, "0x70d869a156d2a1b890bb3df62baf32f7")
	}
	if (absTick & 0x8000) != 0 {
		ratio = mulShift(ratio, "0x31be135f97d08fd981231505542fcfa6")
	}
	if (absTick & 0x10000) != 0 {
		ratio = mulShift(ratio, "0x9aa508b5b7a84e1c677de54f3e99bc9")
	}
	if (absTick & 0x20000) != 0 {
		ratio = mulShift(ratio, "0x5d6af8dedb81196699c329225ee604")
	}
	if (absTick & 0x40000) != 0 {
		ratio = mulShift(ratio, "0x2216e584f5fa1ea926041bedfe98")
	}
	if (absTick & 0x80000) != 0 {
		ratio = mulShift(ratio, "0x48a170391f7dc42444e8fa2")
	}
	if tick > 0 {
		ratio = new(big.Int).Div(constants.MaxUint256, ratio)
	}

	// This divides by 1<<32
	// We then downcast because we know the result always fits within 160 bits due to our tick input constraint
	// We round up in the division so getTickAtSqrtRatio of the output price is always consistent
	rounding := big.NewInt(0)
	if ratio.Mod(ratio, big.NewInt(1<<32)).Cmp(big.NewInt(0)) != 0 {
		rounding = big.NewInt(1)
	}
	sqrtPriceX96 := ratio.Rsh(ratio, 32).Add(ratio, rounding)
	return sqrtPriceX96
}

// Calculates the greatest tick value such that getRatioAtTick(tick) <= ratio
// Accepts sqrtPriceX96, the sqrt ratio for which to compute the tick
// Returns the greatest tick for which the ratio is less than or equal to the input ratio
func getTickAtSqrtRatio(sqrtPriceX96 *big.Int) *big.Int {
	if (sqrtPriceX96.Cmp(constants.MaxSqrtRatio) != -1) || (sqrtPriceX96.Cmp(constants.MinSqrtRatioBig) != 1) {
		panic("tickMath.getTickAtSqrtRatio: INVALID_SQRT_RATIO")
	}

	ratio := new(big.Int).Lsh(sqrtPriceX96, 32)

	r := new(big.Int).Lsh(sqrtPriceX96, 32)
	msb := big.NewInt(0)
	for i := 7; i > 0; i-- {
		cmp := ratio.Cmp(constants.MaxUints[i])
		if cmp == -1 {
			cmp = 0
		}
		f := new(big.Int).Lsh(big.NewInt(int64(cmp)), uint(i))
		msb = new(big.Int).Or(msb, f)
		r = new(big.Int).Rsh(r, uint(f.Int64()))
	}
	cmp := ratio.Cmp(constants.MaxUints[0])
	if cmp == -1 {
		cmp = 0
	}
	f := new(big.Int).Lsh(big.NewInt(int64(cmp)), 0)
	msb = new(big.Int).Or(msb, f)

	dif := new(big.Int).Sub(big.NewInt(127), msb)
	if msb.Cmp(big.NewInt(128)) != -1 {
		r = new(big.Int).Lsh(ratio, uint(dif.Int64()))
	} else {
		r = new(big.Int).Rsh(ratio, uint(dif.Int64()))
	}

	log_2_temp := new(big.Int).Sub(msb, big.NewInt(128))
	log_2 := new(big.Int).Lsh(log_2_temp, 64)

	for i := 0; i < 14; i++ {
		r = new(big.Int).Rsh(new(big.Int).Mul(r, r), 127)
		f := new(big.Int).Rsh(r, 128)
		log_2 = new(big.Int).Or(log_2, new(big.Int).Lsh(f, uint(63-i)))
		r = new(big.Int).Rsh(r, uint(f.Uint64()))
	}

	log_sqrt10001_multiplicand, _ := new(big.Int).SetString("255738958999603826347141", 10)
	log_sqrt10001 := new(big.Int).Mul(log_2, log_sqrt10001_multiplicand)

	tickLow_multiplicand, _ := new(big.Int).SetString("3402992956809132418596140100660247210", 10)
	tickLow := new(big.Int).Rsh(new(big.Int).Add(log_sqrt10001, tickLow_multiplicand), 128)

	tickHigh_multiplicand, _ := new(big.Int).SetString("291339464771989622907027621153398088495", 10)
	tickHigh := new(big.Int).Rsh(new(big.Int).Add(log_sqrt10001, tickHigh_multiplicand), 128)

	if tickLow == tickHigh {
		return tickLow
	}

	sqrtRatio := getSqrtRatioAtTick(int(tickHigh.Int64()))
	if sqrtRatio.Cmp(sqrtPriceX96) <= 0 {
		return tickHigh
	} else {
		return tickLow
	}
}

func mulShift(multiplier *big.Int, multiplicand string) *big.Int {
	multiplicandBig, _ := new(big.Int).SetString(multiplicand, 16)
	productBig := new(big.Int).Mul(multiplier, multiplicandBig)
	result := new(big.Int).Rsh(productBig, 128)
	return result
}
