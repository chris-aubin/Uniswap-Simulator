package sqrtPriceMath

import (
	"fmt"
	"math/big"
	"testing"

	"github.com/chris-aubin/Uniswap-Simulator/src/libraries/constants"
	"github.com/chris-aubin/Uniswap-Simulator/src/libraries/utilities"
)

var (
	// 10^18
	TEN18 = new(big.Int)
	// 10^17
	TEN17 = new(big.Int)
)

// Initialize big.Ints
func init() {
	TEN18.SetString("1000000000000000000", 10)
	TEN17.SetString("100000000000000000", 10)
}

// Fails if price is zero
func TestGetNextSqrtPriceFromInput1(t *testing.T) {
	fmt.Println("GetNextSqrtPriceFromInput: Fails if price is zero")
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("GetNextSqrtPriceFromInput did not panic when price was 0.")
		}
	}()

	GetNextSqrtPriceFromInput(big.NewInt(0), big.NewInt(0), TEN17, false)
}

// Fails if liquidity is zero
func TestGetNextSqrtPriceFromInput2(t *testing.T) {
	fmt.Println("GetNextSqrtPriceFromInput: Fails if liquidity is zero")
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("GetNextSqrtPriceFromInput did not panic when liquidity was 0.")
		}
	}()

	GetNextSqrtPriceFromInput(big.NewInt(1), big.NewInt(0), TEN17, true)
}

// Fails if input amount overflows the price
func TestGetNextSqrtPriceFromInput3(t *testing.T) {
	fmt.Println("GetNextSqrtPriceFromInput: Fails if input amount overflows the price")
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("GetNextSqrtPriceFromInput did not panic when input amount overflows the price.")
		}
	}()

	GetNextSqrtPriceFromInput(big.NewInt(1).Lsh(big.NewInt(1), 160).Sub(big.NewInt(1), big.NewInt(1)), big.NewInt(1024), big.NewInt(1024), false)
}

// Fails if input amount underflows the price
func TestGetNextSqrtPriceFromInput4(t *testing.T) {
	fmt.Println("GetNextSqrtPriceFromInput: Any input amount cannot underflow the price")

	result := GetNextSqrtPriceFromInput(big.NewInt(1), big.NewInt(1), big.NewInt(1).Lsh(big.NewInt(1), 255), true)
	expected := big.NewInt(1)
	if result.Cmp(expected) != 0 {
		t.Errorf("GetNextSqrtPriceFromInput: Got %v; want %v", result, expected)
	}
}

// Returns input price if amount in is zero and zeroForOne = true
func TestGetNextSqrtPriceFromInput5(t *testing.T) {
	fmt.Println("GetNextSqrtPriceFromInput: Returns input price if amount in is zero and zeroForOne = true")
	price := utilities.EncodePriceSqrt(big.NewInt(1), big.NewInt(1))
	result := GetNextSqrtPriceFromInput(price, TEN17, big.NewInt(0), true)
	if result.Cmp(price) != 0 {
		t.Errorf("GetNextSqrtPriceFromInput: Got %v; want %v", result, price)
	}
}

// Returns input price if amount in is zero and zeroForOne = false
func TestGetNextSqrtPriceFromInput6(t *testing.T) {
	fmt.Println("GetNextSqrtPriceFromInput: Returns input price if amount in is zero and zeroForOne = false")
	price := utilities.EncodePriceSqrt(big.NewInt(1), big.NewInt(1))
	result := GetNextSqrtPriceFromInput(price, TEN17, big.NewInt(0), false)
	if result.Cmp(price) != 0 {
		t.Errorf("GetNextSqrtPriceFromInput: Got %v; want %v", result, price)
	}
}

// Returns the minimum price for max inputs
func TestGetNextSqrtPriceFromInput7(t *testing.T) {
	fmt.Println("GetNextSqrtPriceFromInput: Returns the minimum price for max inputs")
	liquidity := constants.MaxUint128
	temp := new(big.Int).Lsh(liquidity, 96)
	maxAmountNoOverflow := new(big.Int).Div(temp, constants.MaxUint160)
	result := GetNextSqrtPriceFromInput(constants.MaxUint160, liquidity, maxAmountNoOverflow, true)
	if result.Cmp(big.NewInt(1)) != 0 {
		t.Errorf("GetNextSqrtPriceFromInput: Got %v; want %v", result, 1)
	}
}

// Input amount of 0.1 token1
func TestGetNextSqrtPriceFromInput8(t *testing.T) {
	fmt.Println("GetNextSqrtPriceFromInput: Input amount of 0.1 token1")
	sqrtP := utilities.EncodePriceSqrt(big.NewInt(1), big.NewInt(1))
	liquidity := TEN18
	amountIn := TEN17
	expected := new(big.Int)
	expected.SetString("87150978765690771352898345369", 10)
	result := GetNextSqrtPriceFromInput(sqrtP, liquidity, amountIn, false)
	if result.Cmp(expected) != 0 {
		t.Errorf("GetNextSqrtPriceFromInput: Got %v; want %v", result, expected)
	}
}

// Input amount of 0.1 token0
func TestGetNextSqrtPriceFromInput9(t *testing.T) {
	fmt.Println("GetNextSqrtPriceFromInput: Input amount of 0.1 token0")
	sqrtP := utilities.EncodePriceSqrt(big.NewInt(1), big.NewInt(1))
	liquidity := TEN18
	amountIn := TEN17
	expected := new(big.Int)
	expected.SetString("72025602285694852357767227579", 10)
	result := GetNextSqrtPriceFromInput(sqrtP, liquidity, amountIn, true)
	if result.Cmp(expected) != 0 {
		t.Errorf("GetNextSqrtPriceFromInput: Got %v; want %v", result, expected)
	}
}

// amountIn > type(uint96).max and zeroForOne = true
func TestGetNextSqrtPriceFromInput10(t *testing.T) {
	fmt.Println("GetNextSqrtPriceFromInput: amountIn > type(uint96).max and zeroForOne = true")
	sqrtP := utilities.EncodePriceSqrt(big.NewInt(1), big.NewInt(1))
	liquidity := TEN18
	amountIn := new(big.Int).Lsh(big.NewInt(1), 100)
	expected := new(big.Int)
	expected.SetString("624999999995069620", 10)
	result := GetNextSqrtPriceFromInput(sqrtP, liquidity, amountIn, true)
	if result.Cmp(expected) != 0 {
		t.Errorf("GetNextSqrtPriceFromInput: Got %v; want %v", result, expected)
	}
}

// Can return 1 with enough amountIn and zeroForOne = true
func TestGetNextSqrtPriceFromInput11(t *testing.T) {
	fmt.Println("GetNextSqrtPriceFromInput: Can return 1 with enough amountIn and zeroForOne = true")
	sqrtP := utilities.EncodePriceSqrt(big.NewInt(1), big.NewInt(1))
	liquidity := big.NewInt(1)
	amountIn := new(big.Int).Div(constants.MaxUint256, big.NewInt(2))
	result := GetNextSqrtPriceFromInput(sqrtP, liquidity, amountIn, true)
	if result.Cmp(big.NewInt(1)) != 0 {
		t.Errorf("GetNextSqrtPriceFromInput: Got %v; want %v", result, 1)
	}
}

// Fails if price is zero
func TestGetNextSqrtPriceFromOutput1(t *testing.T) {
	fmt.Println("GetNextSqrtPriceFromOutput: Fails if price is zero")
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("GetNextSqrtPriceFromOutput did not panic when price was 0.")
		}
	}()

	GetNextSqrtPriceFromOutput(big.NewInt(0), big.NewInt(0), TEN17, false)
}

// Fails if liquidity is zero
func TestGetNextSqrtPriceFromOutput2(t *testing.T) {
	fmt.Println("GetNextSqrtPriceFromOutput: Fails if liquidity is zero")
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("GetNextSqrtPriceFromOutput did not panic when liquidity was 0.")
		}
	}()

	GetNextSqrtPriceFromOutput(big.NewInt(1), big.NewInt(0), TEN17, false)
}

// Fails if output amount is exactly the virtual reserves of token0
func TestGetNextSqrtPriceFromOutput3(t *testing.T) {
	fmt.Println("GetNextSqrtPriceFromOutput: Fails if output amount is exactly the virtual reserves of token0")
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("GetNextSqrtPriceFromOutput did not panic when output amount was exactly the virtual reserves of token0.")
		}
	}()

	price := new(big.Int)
	price.SetString("20282409603651670423947251286016", 10)
	liquidity := big.NewInt(1024)
	amountOut := big.NewInt(4)
	result := GetNextSqrtPriceFromOutput(price, liquidity, amountOut, false)
	fmt.Printf("Got %v; want %v", result, "panic")
}

// Fails if output amount is greater than virtual reserves of token0
func TestGetNextSqrtPriceFromOutput4(t *testing.T) {
	fmt.Println("GetNextSqrtPriceFromOutput: Fails if output amount is greater than virtual reserves of token0")
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("GetNextSqrtPriceFromOutput did not panic when output amount was greater than the virtual reserves of token0.")
		}
	}()

	price := new(big.Int)
	price.SetString("20282409603651670423947251286016", 10)
	liquidity := big.NewInt(1024)
	amountOut := big.NewInt(5)
	result := GetNextSqrtPriceFromOutput(price, liquidity, amountOut, false)
	fmt.Printf("Got %v; want %v", result, "panic")
}

// Fails if output amount is greater than virtual reserves of token1
func TestGetNextSqrtPriceFromOutput5(t *testing.T) {
	fmt.Println("GetNextSqrtPriceFromOutput: Fails if output amount is greater than virtual reserves of token1")
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("GetNextSqrtPriceFromOutput did not panic when output amount was greater than the virtual reserves of token1.")
		}
	}()

	price := new(big.Int)
	price.SetString("20282409603651670423947251286016", 10)
	liquidity := big.NewInt(1024)
	amountOut := big.NewInt(262145)
	result := GetNextSqrtPriceFromOutput(price, liquidity, amountOut, true)
	fmt.Printf("Got %v; want %v", result, "panic")
}

// Fails if output amount is exactly the virtual reserves of token1
func TestGetNextSqrtPriceFromOutput6(t *testing.T) {
	fmt.Println("GetNextSqrtPriceFromOutput: Fails if output amount is exactly the virtual reserves of token1")
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("GetNextSqrtPriceFromOutput: GetNextSqrtPriceFromInput did not panic when output amount was exactly the virtual reserves of token1.")
		}
	}()

	price := new(big.Int)
	price.SetString("20282409603651670423947251286016", 10)
	liquidity := big.NewInt(1024)
	amountOut := big.NewInt(262144)
	result := GetNextSqrtPriceFromOutput(price, liquidity, amountOut, true)
	fmt.Printf("Got %v; want %v", result, "panic")
}

// Succeeds if output amount is just less than the virtual reserves of token1
func TestGetNextSqrtPriceFromOutput7(t *testing.T) {
	fmt.Println("GetNextSqrtPriceFromOutput: Succeeds if output amount is just less than the virtual reserves of token1")
	price := new(big.Int)
	price.SetString("20282409603651670423947251286016", 10)
	liquidity := big.NewInt(1024)
	amountOut := big.NewInt(262143)
	expected := new(big.Int)
	expected.SetString("77371252455336267181195264", 10)
	result := GetNextSqrtPriceFromOutput(price, liquidity, amountOut, true)
	if result.Cmp(expected) != 0 {
		t.Errorf("GetNextSqrtPriceFromOutput: Got %v; want %v", result, expected)
	}
}

// Puzzling echidna test
func TestGetNextSqrtPriceFromOutput8(t *testing.T) {
	fmt.Println("GetNextSqrtPriceFromOutput: Puzzling echidna test")
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("GetNextSqrtPriceFromOutput did not panic.")
		}
	}()

	price := new(big.Int)
	price.SetString("20282409603651670423947251286016", 10)
	liquidity := big.NewInt(1024)
	amountOut := big.NewInt(4)
	GetNextSqrtPriceFromOutput(price, liquidity, amountOut, false)
}

// Returns input price if amount in is zero and zeroForOne = true
func TestGetNextSqrtPriceFromOutput9(t *testing.T) {
	fmt.Println("GetNextSqrtPriceFromOutput: Returns input price if amount in is zero and zeroForOne = true")
	price := utilities.EncodePriceSqrt(big.NewInt(1), big.NewInt(1))
	result := GetNextSqrtPriceFromOutput(price, TEN17, big.NewInt(0), true)
	if result.Cmp(price) != 0 {
		t.Errorf("GetNextSqrtPriceFromOutput: Got %v; want %v", result, price)
	}
}

// Returns input price if amount in is zero and zeroForOne = false
func TestGetNextSqrtPriceFromOutput10(t *testing.T) {
	fmt.Println("GetNextSqrtPriceFromOutput: Returns input price if amount in is zero and zeroForOne = false")
	price := utilities.EncodePriceSqrt(big.NewInt(1), big.NewInt(1))
	result := GetNextSqrtPriceFromOutput(price, TEN17, big.NewInt(0), false)
	if result.Cmp(price) != 0 {
		t.Errorf("GetNextSqrtPriceFromOutput: Got %v; want %v", result, price)
	}
}

// Output amount of 0.1 token1
func TestGetNextSqrtPriceFromOutput11(t *testing.T) {
	fmt.Println("GetNextSqrtPriceFromOutput: Output amount of 0.1 token1")
	price := utilities.EncodePriceSqrt(big.NewInt(1), big.NewInt(1))
	result := GetNextSqrtPriceFromOutput(price, TEN18, TEN17, false)
	expected := new(big.Int)
	expected.SetString("88031291682515930659493278152", 10)
	if result.Cmp(expected) != 0 {
		t.Errorf("GetNextSqrtPriceFromOutput: Got %v; want %v", result, expected)
	}
}

// Output amount of 0.1 token1
func TestGetNextSqrtPriceFromOutput12(t *testing.T) {
	fmt.Println("GetNextSqrtPriceFromOutput: Output amount of 0.1 token1")
	price := utilities.EncodePriceSqrt(big.NewInt(1), big.NewInt(1))
	result := GetNextSqrtPriceFromOutput(price, TEN18, TEN17, true)
	expected := new(big.Int)
	expected.SetString("71305346262837903834189555302", 10)
	if result.Cmp(expected) != 0 {
		t.Errorf("GetNextSqrtPriceFromOutput: Got %v; want %v", result, expected)
	}
}

// Reverts if amountOut is impossible in zero for one direction
func TestGetNextSqrtPriceFromOutput13(t *testing.T) {
	fmt.Println("GetNextSqrtPriceFromOutput: Reverts if amountOut is impossible in zero for one direction")
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("GetNextSqrtPriceFromOutput did not panic when output amount was greater than the virtual reserves of token1.")
		}
	}()

	price := utilities.EncodePriceSqrt(big.NewInt(1), big.NewInt(1))
	GetNextSqrtPriceFromOutput(price, big.NewInt(1), constants.MaxUint256, true)
}

// Reverts if amountOut is impossible in one for zero direction
func TestGetNextSqrtPriceFromOutput14(t *testing.T) {
	fmt.Println("GetNextSqrtPriceFromOutput: Reverts if amountOut is impossible in one for zero direction")
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("GetNextSqrtPriceFromOutput did not panic when output amount was greater than the virtual reserves of token1.")
		}
	}()

	price := utilities.EncodePriceSqrt(big.NewInt(1), big.NewInt(1))
	GetNextSqrtPriceFromOutput(price, big.NewInt(1), constants.MaxUint256, false)
}

// Returns 0 if liquidity is 0
func TestGetAmount0Delta1(t *testing.T) {
	fmt.Println("GetAmount0Delta: Returns 0 if liquidity is 0")
	result := GetAmount0Delta(utilities.EncodePriceSqrt(big.NewInt(1), big.NewInt(1)), utilities.EncodePriceSqrt(big.NewInt(2), big.NewInt(1)), big.NewInt(0), true)
	if result.Cmp(big.NewInt(0)) != 0 {
		t.Errorf("GetAmount0Delta did not return 0 when liquidity was 0.")
	}
}

// Returns 0 if prices are equal
func TestGetAmount0Delta2(t *testing.T) {
	fmt.Println("GetAmount0Delta: Returns 0 if prices are equal")
	result := GetAmount0Delta(utilities.EncodePriceSqrt(big.NewInt(1), big.NewInt(1)), utilities.EncodePriceSqrt(big.NewInt(1), big.NewInt(1)), big.NewInt(0), true)
	if result.Cmp(big.NewInt(0)) != 0 {
		t.Errorf("GetAmount0Delta did not return 0 when prices were equal.")
	}
}

// Returns 0.1 amount1 for price of 1 to 1.21
func TestGetAmount0Delta3(t *testing.T) {
	fmt.Println("GetAmount0Delta: Returns 0.1 amount1 for price of 1 to 1.21")
	amount0 := GetAmount0Delta(utilities.EncodePriceSqrt(big.NewInt(1), big.NewInt(1)), utilities.EncodePriceSqrt(big.NewInt(121), big.NewInt(100)), TEN18, true)
	expected := new(big.Int)
	expected.SetString("90909090909090909", 10)
	if amount0.Cmp(expected) != 0 {
		t.Errorf("GetAmount0Delta: Got %v; want %v", amount0, expected)
	}

	amount0RoundedDown := GetAmount0Delta(utilities.EncodePriceSqrt(big.NewInt(1), big.NewInt(1)), utilities.EncodePriceSqrt(big.NewInt(121), big.NewInt(100)), TEN18, false)
	if amount0RoundedDown.Cmp(new(big.Int).Sub(amount0, big.NewInt(1))) != 0 {
		t.Errorf("GetAmount0Delta: Got %v; want %v", amount0RoundedDown, new(big.Int).Sub(amount0, big.NewInt(1)))
	}
}

// Works for prices that overflow
func TestGetAmount0Delta4(t *testing.T) {
	fmt.Println("GetAmount0Delta: Works for prices that overflow")
	Q90 := new(big.Int).Exp(big.NewInt(2), big.NewInt(90), nil)
	amount0Up := GetAmount0Delta(utilities.EncodePriceSqrt(Q90, big.NewInt(1)), utilities.EncodePriceSqrt(constants.Q96, big.NewInt(1)), TEN18, true)
	amount0Down := GetAmount0Delta(utilities.EncodePriceSqrt(Q90, big.NewInt(1)), utilities.EncodePriceSqrt(constants.Q96, big.NewInt(1)), TEN18, false)
	if amount0Up.Cmp(new(big.Int).Add(amount0Down, big.NewInt(1))) != 0 {
		t.Errorf("GetAmount0Delta: Got %v; want %v", amount0Up, new(big.Int).Add(amount0Down, big.NewInt(1)))
	}
}

// Returns 0 if liquidity is 0
func TestGetAmount1Delta1(t *testing.T) {
	fmt.Println("GetAmount1Delta: Returns 0 if liquidity is 0")
	result := GetAmount1Delta(utilities.EncodePriceSqrt(big.NewInt(1), big.NewInt(1)), utilities.EncodePriceSqrt(big.NewInt(2), big.NewInt(1)), big.NewInt(0), true)
	if result.Cmp(big.NewInt(0)) != 0 {
		t.Errorf("GetAmount1Delta did not return 0 when liquidity was 0.")
	}
}

// Returns 0 if prices are equal
func TestGetAmount1Delta2(t *testing.T) {
	fmt.Println("GetAmount1Delta: Returns 0 if prices are equal")
	result := GetAmount1Delta(utilities.EncodePriceSqrt(big.NewInt(1), big.NewInt(1)), utilities.EncodePriceSqrt(big.NewInt(1), big.NewInt(1)), big.NewInt(0), true)
	if result.Cmp(big.NewInt(0)) != 0 {
		t.Errorf("GetAmount1Delta did not return 0 when prices were equal.")
	}
}

// Returns 0.1 amount1 for price of 1 to 1.21
func TestGetAmount1Delta3(t *testing.T) {
	fmt.Println("GetAmount1Delta: Returns 0.1 amount1 for price of 1 to 1.21")
	amount1 := GetAmount1Delta(utilities.EncodePriceSqrt(big.NewInt(1), big.NewInt(1)), utilities.EncodePriceSqrt(big.NewInt(121), big.NewInt(100)), TEN18, true)
	if amount1.Cmp(TEN17) != 0 {
		t.Errorf("GetAmount1Delta: Got %v; want %v", amount1, TEN17)
	}

	amount1RoundedDown := GetAmount1Delta(utilities.EncodePriceSqrt(big.NewInt(1), big.NewInt(1)), utilities.EncodePriceSqrt(big.NewInt(121), big.NewInt(100)), TEN18, false)
	if amount1RoundedDown.Cmp(new(big.Int).Sub(amount1, big.NewInt(1))) != 0 {
		t.Errorf("GetAmount1Delta: Got %v; want %v", amount1RoundedDown, new(big.Int).Sub(amount1, big.NewInt(1)))
	}
}

// Swap computation
func TestSwapComputation(t *testing.T) {
	fmt.Println("Swap computation")

	sqrtP := new(big.Int)
	sqrtP.SetString("1025574284609383690408304870162715216695788925244", 10)
	liquidity := new(big.Int)
	liquidity.SetString("50015962439936049619261659728067971248", 10)
	zeroForOne := true
	amountIn := new(big.Int)
	amountIn.SetString("406", 10)
	expected := new(big.Int)
	expected.SetString("1025574284609383582644711336373707553698163132913", 10)

	sqrtQ := GetNextSqrtPriceFromInput(sqrtP, liquidity, amountIn, zeroForOne)
	if sqrtQ.Cmp(expected) != 0 {
		t.Errorf("SwapComputation: Got %v; want %v", sqrtQ, expected)
	}

	amount0Delta := GetAmount0Delta(sqrtP, sqrtQ, liquidity, true)
	if amount0Delta.Cmp(amountIn) != 0 {
		t.Errorf("SwapComputation: Got %v; want %v", amount0Delta, amountIn)
	}
}
