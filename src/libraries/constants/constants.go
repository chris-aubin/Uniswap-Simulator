package constants

import "math/big"

const (
	// The maximum tick that may be passed to getSqrtRatioAtTick.
	// log base 1.0001 of 2**128
	MaxTick = 887272
	// The minimum tick that may be passed to getSqrtRatioAtTick (-MaxTick).
	MinTick = -887272
	// The minimum value that can be returned by getSqrtRatioAtTick.
	// Equivalent to getSqrtRatioAtTick(MinTick).
	MinSqrtRatio = 4295128739
)

// Declaring big.Ints as constants proved challenging, so they are declared as
// variables that are just never changed after they are initialised.
var (
	// MinTick as a big.Int for use in calculations.
	MinTickBig = big.NewInt(-887272)
	// MaxTick as a big.Int for use in calculations.
	MaxTickBig = big.NewInt(887272)
	// MinSqrtRatio as a big.Int for use in calculations.
	MinSqrtRatioBig = big.NewInt(4295128739)
	// The maximum value that can be returned by getSqrtRatioAtTick.
	// Equivalent to getSqrtRatioAtTick(MaxTick).
	MaxSqrtRatio = new(big.Int)
	// For handling binary fixed point numbers, see:
	// https://en.wikipedia.org/wiki/Q_(number_format)
	Q128 = new(big.Int)
	Q96  = new(big.Int)
	// Maximum unsigned integers for given number of bits.
	MaxUint256 = new(big.Int)
	MaxUint160 = new(big.Int)
	MaxUint128 = new(big.Int)
	MaxUint64  = new(big.Int)
	MaxUint32  = new(big.Int)
	MaxUint16  = new(big.Int)
	MaxUint8   = new(big.Int)
	MaxUint4   = new(big.Int)
	MaxUint2   = new(big.Int)
	MaxUint1   = new(big.Int)
	// Map of maximum unsigned integers for given number of bits.
	MaxUints = map[int]*big.Int{
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

// Initialise the big.Ints.
func init() {
	MaxSqrtRatio.SetString("1461446703485210103287273052203988822378723970342", 10)
	Q128.SetString("100000000000000000000000000000000", 16)
	Q96.SetString("1000000000000000000000000", 16)
	MaxUint256.SetString("ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff", 16)
	MaxUint160.SetString("ffffffffffffffffffffffffffffffffffffffff", 16)
	MaxUint128.SetString("ffffffffffffffffffffffffffffffff", 16)
	MaxUint64.SetString("ffffffffffffffff", 16)
	MaxUint32.SetString("ffffffff", 16)
	MaxUint16.SetString("ffff", 16)
	MaxUint8.SetString("ff", 16)
	MaxUint4.SetString("f", 16)
	MaxUint2.SetString("3", 16)
	MaxUint1.SetString("1", 16)
}
