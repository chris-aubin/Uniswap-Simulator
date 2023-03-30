package constants

import "math/big"

const (
	MinTick = -887272
	// The maximum tick that may be passed to #getSqrtRatioAtTick computed from log base 1.0001 of 2**128
	MaxTick = 887272
	// The minimum value that can be returned from #getSqrtRatioAtTick. Equivalent to getSqrtRatioAtTick(MIN_TICK)
	MinSqrtRatio = 4295128739
)

var (
	// The minimum tick that may be passed to #getSqrtRatioAtTick computed from log base 1.0001 of 2**-128
	MinTickBig = big.NewInt(-887272)
	// The maximum tick that may be passed to #getSqrtRatioAtTick computed from log base 1.0001 of 2**128
	MaxTickBig = big.NewInt(887272)
	// The minimum value that can be returned from #getSqrtRatioAtTick. Equivalent to getSqrtRatioAtTick(MIN_TICK)
	MinSqrtRatioBig = big.NewInt(4295128739)
	// The maximum value that can be returned from #getSqrtRatioAtTick. Equivalent to getSqrtRatioAtTick(MAX_TICK)
	MaxSqrtRatio, _ = new(big.Int).SetString("1461446703485210103287273052203988822378723970342", 10)
	Q96, _          = new(big.Int).SetString("0x1000000000000000000000000", 16)
	// The maximum _ bit unsigned integer
	MaxUint256, _ = new(big.Int).SetString("0xffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff", 16)
	MaxUint160, _ = new(big.Int).SetString("0xffffffffffffffffffffffffffffffffffffffff", 16)
	MaxUint128, _ = new(big.Int).SetString("0xffffffffffffffffffffffffffffffff", 16)
	MaxUint64, _  = new(big.Int).SetString("0xffffffffffffffff", 16)
	MaxUint32, _  = new(big.Int).SetString("0xffffffff", 16)
	MaxUint16, _  = new(big.Int).SetString("0xffff", 16)
	MaxUint8, _   = new(big.Int).SetString("0xff", 16)
	MaxUint4, _   = new(big.Int).SetString("0xf", 16)
	MaxUint2, _   = new(big.Int).SetString("0x3", 16)
	MaxUint1, _   = new(big.Int).SetString("0x1", 16)
	MaxUints      = map[int]*big.Int{
		8: MaxUint256,
		7: MaxUint128,
		6: MaxUint64,
		5: MaxUint32,
		4: MaxUint16,
		3: MaxUint8,
		2: MaxUint4,
		1: MaxUint2,
		0: MaxUint1,
	}
)
