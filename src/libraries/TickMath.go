package tickmath

import (
	"math/big"
)

const (
	// The minimum tick that may be passed to #getSqrtRatioAtTick computed from log base 1.0001 of 2**-128
	MinTick int = -887272
	// The maximum tick that may be passed to #getSqrtRatioAtTick computed from log base 1.0001 of 2**128
	MaxTick int = 887272
	// The minimum value that can be returned from #getSqrtRatioAtTick. Equivalent to getSqrtRatioAtTick(MIN_TICK)
	MinSqrtRatio int = 4295128739
)

var (
	// The maximum value that can be returned from #getSqrtRatioAtTick. Equivalent to getSqrtRatioAtTick(MAX_TICK)
	MaxSqrtRatio, _ = new(big.Int).SetString("1461446703485210103287273052203988822378723970342", 10)
	// The maximum 256 bit unsigned integer
	MaxUint256, _ = new(big.Int).SetString("0xffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff", 16)
)

func getSqrtRatioAtTick(tick int) *big.Int {
	absTick := tick
	if tick < 0 {
		absTick = -tick
	}

	if absTick > MaxTick {
		panic("INVALID_TICK")
	}

	ratio := new(big.Int)
	if absTick&0x1 != 0 {
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
		ratio = new(big.Int).Div(MaxUint256, ratio)
	}

	// this divides by 1<<32 rounding up to go from a Q128.128 to a Q128.96.
	// we then downcast because we know the result always fits within 160 bits due to our tick input constraint
	// we round up in the division so getTickAtSqrtRatio of the output price is always consistent
	rounding := big.NewInt(0)
	if ratio.Mod(ratio, big.NewInt(1<<32)).Cmp(big.NewInt(0)) != 0 {
		rounding = big.NewInt(1)
	}
	sqrtPriceX96 := ratio.Rsh(ratio, 32).Add(ratio, rounding)
	return sqrtPriceX96
}

func mulShift(multiplier *big.Int, multiplicand string) *ui.Int {
	multiplicandBig, _ := new(big.Int).SetString(multiplicand, 16)
	productBig := new(big.Int).Mul(multiplier, multiplicandBig)
	result := productBig.Rsh(productBig, 128)
	return result
}
