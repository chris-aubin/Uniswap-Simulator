package swapMath

import (
	"fmt"
	"math/big"
	"testing"

	"github.com/chris-aubin/Uniswap-Simulator/src/libraries/sqrtPriceMath"
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

func TestComputeSwapStep1(t *testing.T) {
	fmt.Println("Exact amount in that gets capped at price target in one for zero")

	price := utilities.EncodePriceSqrt(big.NewInt(1), big.NewInt(1))
	priceTarget := utilities.EncodePriceSqrt(big.NewInt(101), big.NewInt(100))
	liquidity := new(big.Int).Mul(big.NewInt(2), TEN18)
	amount := new(big.Int).Set(TEN18)
	fee := 600
	zeroForOne := false

	sqrtQ, amountIn, amountOut, feeAmount := ComputeSwapStep(price, priceTarget, liquidity, amount, fee)

	amountInExpected := new(big.Int)
	amountInExpected.SetString("9975124224178055", 10)
	feeAmountExpected := new(big.Int)
	feeAmountExpected.SetString("5988667735148", 10)
	amountOutExpected := new(big.Int)
	amountOutExpected.SetString("9925619580021728", 10)
	dif := new(big.Int).Sub(amount, new(big.Int).Add(amountIn, feeAmount))
	if (amountInExpected.Cmp(amountIn) != 0) ||
		(amountOutExpected.Cmp(amountOut) != 0) ||
		(feeAmountExpected.Cmp(feeAmount) != 0) {
		t.Errorf("TestComputeSwapStep1 failed, got (amountIn, amountOut, feeAmount): %v, %v, %v, expected: %v, %v, %v", amountIn, amountOut, feeAmount, amountInExpected, amountOutExpected, feeAmountExpected)
	}
	if dif.Cmp(big.NewInt(0)) <= 0 {
		t.Errorf("TestComputeSwapStep1 failed, entire amount should not be used, dif should be >0. dif: %v", dif)
	}

	priceAfterWholeInputAmount := sqrtPriceMath.GetNextSqrtPriceFromInput(price, liquidity, amount, zeroForOne)

	if (sqrtQ.Cmp(priceTarget) != 0) || (sqrtQ.Cmp(priceAfterWholeInputAmount) >= 0) {
		t.Errorf("TestComputeSwapStep1 failed, got: %v, %v, expected: %v, %v", sqrtQ, priceAfterWholeInputAmount, priceTarget, priceAfterWholeInputAmount)
	}
}

func TestComputeSwapStep2(t *testing.T) {
	fmt.Println("Exact amount out that gets capped at price target in one for zero")

	price := utilities.EncodePriceSqrt(big.NewInt(1), big.NewInt(1))
	priceTarget := utilities.EncodePriceSqrt(big.NewInt(101), big.NewInt(100))
	liquidity := new(big.Int).Mul(big.NewInt(2), TEN18)
	amount := new(big.Int).Mul(big.NewInt(-1), TEN18)
	fee := 600
	zeroForOne := false

	sqrtQ, amountIn, amountOut, feeAmount := ComputeSwapStep(price, priceTarget, liquidity, amount, fee)

	amountInExpected := new(big.Int)
	amountInExpected.SetString("9975124224178055", 10)
	feeAmountExpected := new(big.Int)
	feeAmountExpected.SetString("5988667735148", 10)
	amountOutExpected := new(big.Int)
	amountOutExpected.SetString("9925619580021728", 10)
	dif := new(big.Int).Sub(new(big.Int).Neg(amount), amountOut)
	if (amountInExpected.Cmp(amountIn) != 0) ||
		(feeAmountExpected.Cmp(feeAmount) != 0) ||
		(amountOutExpected.Cmp(amountOut) != 0) ||
		(dif.Cmp(big.NewInt(0)) <= 0) {
		t.Errorf("TestComputeSwapStep2 failed, got: %v, %v, %v, %v, expected: %v, %v, %v, %v", amountIn, amountOut, sqrtQ, feeAmount, amountInExpected, amountOutExpected, priceTarget, feeAmountExpected)
	}

	priceAfterWholeOutputAmount := sqrtPriceMath.GetNextSqrtPriceFromOutput(price, liquidity, new(big.Int).Neg(amount), zeroForOne)

	if (sqrtQ.Cmp(priceTarget) != 0) || (sqrtQ.Cmp(priceAfterWholeOutputAmount) >= 0) {
		t.Errorf("TestComputeSwapStep2 failed, got: %v, %v, expected: %v, %v", sqrtQ, priceAfterWholeOutputAmount, priceTarget, priceAfterWholeOutputAmount)
	}
}

func TestComputeSwapStep3(t *testing.T) {
	fmt.Println("Exact amount in that is fully spent in one for zero")

	price := utilities.EncodePriceSqrt(big.NewInt(1), big.NewInt(1))
	priceTarget := utilities.EncodePriceSqrt(big.NewInt(1000), big.NewInt(100))
	liquidity := new(big.Int).Mul(big.NewInt(2), TEN18)
	amount := new(big.Int).Set(TEN18)
	fee := 600
	zeroForOne := false

	sqrtQ, amountIn, amountOut, feeAmount := ComputeSwapStep(price, priceTarget, liquidity, amount, fee)

	amountInExpected := new(big.Int)
	amountInExpected.SetString("999400000000000000", 10)
	feeAmountExpected := new(big.Int)
	feeAmountExpected.SetString("600000000000000", 10)
	amountOutExpected := new(big.Int)
	amountOutExpected.SetString("666399946655997866", 10)
	dif := new(big.Int).Sub(new(big.Int).Add(feeAmount, amountIn), amount)
	if (amountInExpected.Cmp(amountIn) != 0) ||
		(feeAmountExpected.Cmp(feeAmount) != 0) ||
		(amountOutExpected.Cmp(amountOut) != 0) {
		t.Errorf("TestComputeSwapStep3 failed, got: %v, %v, %v, expected: %v, %v, %v", amountIn, feeAmount, amountOut, amountInExpected, feeAmountExpected, amountOutExpected)
	}
	if dif.Cmp(big.NewInt(0)) != 0 {
		t.Errorf("TestComputeSwapStep3 failed, the entire amount was not used. dif: %v", dif)
	}

	priceAfterWholeInputAmountLessFee := sqrtPriceMath.GetNextSqrtPriceFromInput(price, liquidity, new(big.Int).Sub(amount, feeAmount), zeroForOne)

	if (sqrtQ.Cmp(priceTarget) >= 0) || (sqrtQ.Cmp(priceAfterWholeInputAmountLessFee) != 0) {
		t.Errorf("TestComputeSwapStep3 failed, got: %v, %v, expected: %v, %v", sqrtQ, priceAfterWholeInputAmountLessFee, priceTarget, priceAfterWholeInputAmountLessFee)
	}
}

func TestComputeSwapStep4(t *testing.T) {
	fmt.Println("Exact amount out that is fully received in one for zero")

	price := utilities.EncodePriceSqrt(big.NewInt(1), big.NewInt(1))
	priceTarget := utilities.EncodePriceSqrt(big.NewInt(10000), big.NewInt(100))
	liquidity := new(big.Int).Mul(big.NewInt(2), TEN18)
	amount := new(big.Int).Mul(big.NewInt(-1), TEN18)
	fee := 600
	zeroForOne := false

	sqrtQ, amountIn, amountOut, feeAmount := ComputeSwapStep(price, priceTarget, liquidity, amount, fee)

	amountInExpected := new(big.Int)
	amountInExpected.SetString("2000000000000000000", 10)
	feeAmountExpected := new(big.Int)
	feeAmountExpected.SetString("1200720432259356", 10)
	amountOutExpected := new(big.Int).Mul(amount, big.NewInt(-1))
	if (amountInExpected.Cmp(amountIn) != 0) ||
		(feeAmountExpected.Cmp(feeAmount) != 0) ||
		(amountOutExpected.Cmp(amountOut) != 0) {
		t.Errorf("TestComputeSwapStep4 failed, got: %v, %v, %v, %v, expected: %v, %v, %v, %v", amountIn, amountOut, sqrtQ, feeAmount, amountInExpected, amountOutExpected, priceTarget, feeAmountExpected)
	}

	priceAfterWholeOutputAmount := sqrtPriceMath.GetNextSqrtPriceFromOutput(price, liquidity, new(big.Int).Mul(amount, big.NewInt(-1)), zeroForOne)

	if (sqrtQ.Cmp(priceTarget) >= 0) || (sqrtQ.Cmp(priceAfterWholeOutputAmount) != 0) {
		t.Errorf("TestComputeSwapStep4 failed, got: %v, %v, expected: %v, %v", sqrtQ, priceAfterWholeOutputAmount, priceTarget, priceAfterWholeOutputAmount)
	}
}

func TestComputeSwapStep5(t *testing.T) {
	fmt.Println("Amount out is capped at the desired amount out")
	price := new(big.Int)
	price.SetString("417332158212080721273783715441582", 10)
	priceTarget := new(big.Int)
	priceTarget.SetString("1452870262520218020823638996", 10)
	liquidity := new(big.Int)
	liquidity.SetString("159344665391607089467575320103", 10)
	amount := new(big.Int)
	amount.SetString("-1", 10)
	fee := 1
	sqrtQ, amountIn, amountOut, feeAmount := ComputeSwapStep(price, priceTarget, liquidity, amount, fee)

	amountInExpected := big.NewInt(1)
	feeAmountExpected := big.NewInt(1)
	amountOutExpected := big.NewInt(1)
	sqrtQExpected := new(big.Int)
	sqrtQExpected.SetString("417332158212080721273783715441581", 10)
	if (amountInExpected.Cmp(amountIn) != 0) ||
		(feeAmountExpected.Cmp(feeAmount) != 0) ||
		(amountOutExpected.Cmp(amountOut) != 0) ||
		(sqrtQExpected.Cmp(sqrtQ) != 0) {
		t.Errorf("TestComputeSwapStep5 failed, got: %v, %v, %v, %v, expected: %v, %v, %v, %v", amountIn, amountOut, sqrtQ, feeAmount, amountInExpected, amountOutExpected, sqrtQExpected, feeAmountExpected)
	}
}

func TestComputeSwapStep6(t *testing.T) {
	fmt.Println("Target price of 1 uses partial input amount")

	amountOutExpected := new(big.Int)
	amountOutExpected.SetString("3915081100057732413702495386755767", 10)
	sqrtQ, amountIn, amountOut, feeAmount := ComputeSwapStep(big.NewInt(2), big.NewInt(1), big.NewInt(1), amountOutExpected, 1)

	amountInExpected := new(big.Int)
	amountInExpected.SetString("39614081257132168796771975168", 10)
	feeAmountExpected := new(big.Int)
	feeAmountExpected.SetString("39614120871253040049813", 10)
	amountOutExpected = big.NewInt(0)
	sqrtQExpected := big.NewInt(1)
	temp := new(big.Int)
	temp.SetString("3915081100057732413702495386755767", 10)
	dif := new(big.Int).Sub(temp, new(big.Int).Add(amountIn, feeAmount))
	if (amountInExpected.Cmp(amountIn) != 0) ||
		(feeAmountExpected.Cmp(feeAmount) != 0) ||
		(amountOutExpected.Cmp(amountOut) != 0) ||
		(sqrtQExpected.Cmp(sqrtQ)) != 0 ||
		(dif.Cmp(big.NewInt(0)) <= -1) {
		t.Errorf("TestComputeSwapStep6 failed, got: %v, %v, %v, %v, expected: %v, %v, %v, %v", amountIn, amountOut, sqrtQ, feeAmount, amountInExpected, amountOutExpected, sqrtQExpected, feeAmountExpected)
	}
}

func TestComputeSwapStep7(t *testing.T) {
	fmt.Println("Entire input amount taken as fee")

	price := big.NewInt(2413)
	priceTarget := new(big.Int)
	priceTarget.SetString("79887613182836312", 10)
	liquidity := new(big.Int)
	liquidity.SetString("1985041575832132834610021537970", 10)
	amount := big.NewInt(10)
	fee := 1872

	sqrtQ, amountIn, amountOut, feeAmount := ComputeSwapStep(price, priceTarget, liquidity, amount, fee)

	amountInExpected := big.NewInt(0)
	feeAmountExpected := big.NewInt(10)
	amountOutExpected := big.NewInt(0)
	sqrtQExpected := big.NewInt(2413)
	if (amountInExpected.Cmp(amountIn) != 0) ||
		(feeAmountExpected.Cmp(feeAmount) != 0) ||
		(amountOutExpected.Cmp(amountOut) != 0) ||
		(sqrtQExpected.Cmp(sqrtQ)) != 0 {
		t.Errorf("TestComputeSwapStep7 failed, got: %v, %v, %v, %v, expected: %v, %v, %v, %v", amountIn, amountOut, sqrtQ, feeAmount, amountInExpected, amountOutExpected, sqrtQExpected, feeAmountExpected)
	}
}

func TestComputeSwapStep8(t *testing.T) {
	fmt.Println("Handles intermediate insufficient liquidity in zero for one exact output case")

	sqrtP := new(big.Int)
	sqrtP.SetString("20282409603651670423947251286016", 10)
	sqrtPTarget := new(big.Int).Div(new(big.Int).Mul(sqrtP, big.NewInt(11)), big.NewInt(10))
	liquidity := big.NewInt(1024)
	amountRemaining := big.NewInt(-4)
	feePips := 3000

	sqrtQ, amountIn, amountOut, feeAmount := ComputeSwapStep(sqrtP, sqrtPTarget, liquidity, amountRemaining, feePips)

	amountInExpected := big.NewInt(26215)
	feeAmountExpected := big.NewInt(79)
	amountOutExpected := big.NewInt(0)
	sqrtQExpected := new(big.Int).Set(sqrtPTarget)
	if (amountInExpected.Cmp(amountIn) != 0) ||
		(feeAmountExpected.Cmp(feeAmount) != 0) ||
		(amountOutExpected.Cmp(amountOut) != 0) ||
		(sqrtQExpected.Cmp(sqrtQ)) != 0 {
		t.Errorf("TestComputeSwapStep8 failed, got: %v, %v, %v, %v, expected: %v, %v, %v, %v", amountIn, amountOut, sqrtQ, feeAmount, amountInExpected, amountOutExpected, sqrtQExpected, feeAmountExpected)
	}
}

func TestComputeSwapStep9(t *testing.T) {
	fmt.Println("Handles intermediate insufficient liquidity in one for zero exact output case")

	sqrtP := new(big.Int)
	sqrtP.SetString("20282409603651670423947251286016", 10)
	sqrtPTarget := new(big.Int).Div(new(big.Int).Mul(sqrtP, big.NewInt(9)), big.NewInt(10))
	liquidity := big.NewInt(1024)
	amountRemaining := big.NewInt(-263000)
	feePips := 3000

	sqrtQ, amountIn, amountOut, feeAmount := ComputeSwapStep(sqrtP, sqrtPTarget, liquidity, amountRemaining, feePips)

	amountInExpected := big.NewInt(1)
	feeAmountExpected := big.NewInt(1)
	amountOutExpected := big.NewInt(26214)
	sqrtQExpected := new(big.Int).Set(sqrtPTarget)
	if (amountInExpected.Cmp(amountIn) != 0) ||
		(feeAmountExpected.Cmp(feeAmount) != 0) ||
		(amountOutExpected.Cmp(amountOut) != 0) ||
		(sqrtQExpected.Cmp(sqrtQ)) != 0 {
		t.Errorf("TestComputeSwapStep9 failed, got: %v, %v, %v, %v, expected: %v, %v, %v, %v", amountIn, amountOut, sqrtQ, feeAmount, amountInExpected, amountOutExpected, sqrtQExpected, feeAmountExpected)
	}
}
