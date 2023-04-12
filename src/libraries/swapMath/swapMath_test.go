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

	// const price = encodePriceSqrt(1, 1)
	// const priceTarget = encodePriceSqrt(101, 100)
	// const liquidity = expandTo18Decimals(2)
	// const amount = expandTo18Decimals(1)
	// const fee = 600
	// const zeroForOne = false
	price := utilities.EncodePriceSqrt(big.NewInt(1), big.NewInt(1))
	priceTarget := utilities.EncodePriceSqrt(big.NewInt(101), big.NewInt(100))
	liquidity := new(big.Int).Mul(big.NewInt(2), TEN18)
	amount := new(big.Int).Set(TEN18)
	fee := 600
	zeroForOne := false

	// const { amountIn, amountOut, sqrtQ, feeAmount } = await swapMath.computeSwapStep(
	//     price,
	//     priceTarget,
	//     liquidity,
	//     amount,
	//     fee
	//   )
	sqrtQ, amountIn, amountOut, feeAmount := ComputeSwapStep(price, priceTarget, liquidity, amount, fee)

	// expect(amountIn).to.eq('9975124224178055')
	// expect(feeAmount).to.eq('5988667735148')
	// expect(amountOut).to.eq('9925619580021728')
	// expect(amountIn.add(feeAmount), 'entire amount is not used').to.lt(amount)
	amountInExpected := new(big.Int)
	amountInExpected.SetString("9975124224178055", 10)
	feeAmountExpected := new(big.Int)
	feeAmountExpected.SetString("5988667735148", 10)
	amountOutExpected := new(big.Int)
	amountOutExpected.SetString("9925619580021728", 10)
	dif := new(big.Int).Sub(amount, new(big.Int).Add(amountIn, feeAmount))
	if (amountInExpected.Cmp(amountIn) != 0) ||
		(feeAmountExpected.Cmp(feeAmount) != 0) ||
		(amountOutExpected.Cmp(amountOut) != 0) ||
		(dif.Cmp(big.NewInt(0)) <= 0) {
		t.Errorf("TestComputeSwapStep1 failed, got: %v, %v, %v, %v, expected: %v, %v, %v, %v", amountIn, amountOut, sqrtQ, feeAmount, amountInExpected, amountOutExpected, priceTarget, feeAmountExpected)
	}

	// const priceAfterWholeInputAmount = await sqrtPriceMath.getNextSqrtPriceFromInput(
	// 	price,
	// 	liquidity,
	// 	amount,
	// 	zeroForOne
	// 	)
	priceAfterWholeInputAmount := sqrtPriceMath.GetNextSqrtPriceFromInput(price, liquidity, amount, zeroForOne)

	// expect(sqrtQ, 'price is capped at price target').to.eq(priceTarget)
	// expect(sqrtQ, 'price is less than price after whole input amount').to.lt(priceAfterWholeInputAmount)
	if (sqrtQ.Cmp(priceTarget) != 0) || (sqrtQ.Cmp(priceAfterWholeInputAmount) <= 0) {
		t.Errorf("TestComputeSwapStep1 failed, got: %v, %v, expected: %v, %v", sqrtQ, priceAfterWholeInputAmount, priceTarget, priceAfterWholeInputAmount)
	}
}

func TestComputeSwapStep2(t *testing.T) {
	fmt.Println("Exact amount out that gets capped at price target in one for zero")

	// it('exact amount out that gets capped at price target in one for zero', async () => {
	//   const price = encodePriceSqrt(1, 1)
	//   const priceTarget = encodePriceSqrt(101, 100)
	//   const liquidity = expandTo18Decimals(2)
	//   const amount = expandTo18Decimals(1).mul(-1)
	//   const fee = 600
	//   const zeroForOne = false
	price := utilities.EncodePriceSqrt(big.NewInt(1), big.NewInt(1))
	priceTarget := utilities.EncodePriceSqrt(big.NewInt(101), big.NewInt(100))
	liquidity := new(big.Int).Mul(big.NewInt(2), TEN18)
	amount := new(big.Int).Mul(big.NewInt(-1), TEN18)
	fee := 600
	zeroForOne := false

	// const { amountIn, amountOut, sqrtQ, feeAmount } = await swapMath.computeSwapStep(
	//     price,
	//     priceTarget,
	//     liquidity,
	//     amount,
	//     fee
	//   )
	sqrtQ, amountIn, amountOut, feeAmount := ComputeSwapStep(price, priceTarget, liquidity, amount, fee)

	// expect(amountIn).to.eq('9975124224178055')
	// expect(feeAmount).to.eq('5988667735148')
	// expect(amountOut).to.eq('9925619580021728')
	// expect(amountOut, 'entire amount out is not returned').to.lt(amount.mul(-1))
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

	//   const priceAfterWholeOutputAmount = await sqrtPriceMath.getNextSqrtPriceFromOutput(
	//     price,
	//     liquidity,
	//     amount.mul(-1),
	//     zeroForOne
	//   )
	priceAfterWholeOutputAmount := sqrtPriceMath.GetNextSqrtPriceFromOutput(price, liquidity, new(big.Int).Neg(amount), zeroForOne)

	//   expect(sqrtQ, 'price is capped at price target').to.eq(priceTarget)
	//   expect(sqrtQ, 'price is less than price after whole output amount').to.lt(priceAfterWholeOutputAmount)
	if (sqrtQ.Cmp(priceTarget) != 0) || (sqrtQ.Cmp(priceAfterWholeOutputAmount) >= 0) {
		t.Errorf("TestComputeSwapStep2 failed, got: %v, %v, expected: %v, %v", sqrtQ, priceAfterWholeOutputAmount, priceTarget, priceAfterWholeOutputAmount)
	}
}

func TestComputeSwapStep3(t *testing.T) {
	fmt.Println("Exact amount in that is fully spent in one for one")

	//   const price = encodePriceSqrt(1, 1)
	//   const priceTarget = encodePriceSqrt(1000, 100)
	//   const liquidity = expandTo18Decimals(2)
	//   const amount = expandTo18Decimals(1)
	//   const fee = 600
	//   const zeroForOne = false
	price := utilities.EncodePriceSqrt(big.NewInt(1), big.NewInt(1))
	priceTarget := utilities.EncodePriceSqrt(big.NewInt(1000), big.NewInt(100))
	liquidity := new(big.Int).Mul(big.NewInt(2), TEN18)
	amount := new(big.Int).Set(TEN18)
	fee := 600
	zeroForOne := false

	//   const { amountIn, amountOut, sqrtQ, feeAmount } = await swapMath.computeSwapStep(
	//     price,
	//     priceTarget,
	//     liquidity,
	//     amount,
	//     fee
	//   )
	sqrtQ, amountIn, amountOut, feeAmount := ComputeSwapStep(price, priceTarget, liquidity, amount, fee)

	//   expect(amountIn).to.eq('999400000000000000')
	//   expect(feeAmount).to.eq('600000000000000')
	//   expect(amountOut).to.eq('666399946655997866')
	//   expect(amountIn.add(feeAmount), 'entire amount is used').to.eq(amount)
	amountInExpected := new(big.Int)
	amountInExpected.SetString("999400000000000000", 10)
	feeAmountExpected := new(big.Int)
	feeAmountExpected.SetString("600000000000000", 10)
	amountOutExpected := new(big.Int)
	amountOutExpected.SetString("666399946655997866", 10)
	dif := new(big.Int).Sub(new(big.Int).Add(feeAmount, amountIn), amount)
	if (amountInExpected.Cmp(amountIn) != 0) ||
		(feeAmountExpected.Cmp(feeAmount) != 0) ||
		(amountOutExpected.Cmp(amountOut) != 0) ||
		(dif.Cmp(big.NewInt(0)) != 0) {
		t.Errorf("TestComputeSwapStep3 failed, got: %v, %v, %v, %v, expected: %v, %v, %v, %v", amountIn, amountOut, sqrtQ, feeAmount, amountInExpected, amountOutExpected, priceTarget, feeAmountExpected)
	}

	//   const priceAfterWholeInputAmountLessFee = await sqrtPriceMath.getNextSqrtPriceFromInput(
	//     price,
	//     liquidity,
	//     amount.sub(feeAmount),
	//     zeroForOne
	//   )
	priceAfterWholeInputAmountLessFee := sqrtPriceMath.GetNextSqrtPriceFromInput(price, liquidity, new(big.Int).Sub(amount, feeAmount), zeroForOne)

	//   expect(sqrtQ, 'price does not reach price target').to.be.lt(priceTarget)
	//   expect(sqrtQ, 'price is equal to price after whole input amount').to.eq(priceAfterWholeInputAmountLessFee)
	if (sqrtQ.Cmp(priceTarget) >= 0) || (sqrtQ.Cmp(priceAfterWholeInputAmountLessFee) != 0) {
		t.Errorf("TestComputeSwapStep3 failed, got: %v, %v, expected: %v, %v", sqrtQ, priceAfterWholeInputAmountLessFee, priceTarget, priceAfterWholeInputAmountLessFee)
	}
}

func TestComputeSwapStep4(t *testing.T) {
	fmt.Println("Exact amount out that is fully received in one for zero")

	//   const price = encodePriceSqrt(1, 1)
	//   const priceTarget = encodePriceSqrt(10000, 100)
	//   const liquidity = expandTo18Decimals(2)
	//   const amount = expandTo18Decimals(1).mul(-1)
	//   const fee = 600
	//   const zeroForOne = false
	price := utilities.EncodePriceSqrt(big.NewInt(1), big.NewInt(1))
	priceTarget := utilities.EncodePriceSqrt(big.NewInt(10000), big.NewInt(100))
	liquidity := new(big.Int).Mul(big.NewInt(2), TEN18)
	amount := new(big.Int).Mul(big.NewInt(-1), TEN18)
	fee := 600
	zeroForOne := false

	//   const { amountIn, amountOut, sqrtQ, feeAmount } = await swapMath.computeSwapStep(
	//     price,
	//     priceTarget,
	//     liquidity,
	//     amount,
	//     fee
	//   )
	sqrtQ, amountIn, amountOut, feeAmount := ComputeSwapStep(price, priceTarget, liquidity, amount, fee)

	//   expect(amountIn).to.eq('2000000000000000000')
	//   expect(feeAmount).to.eq('1200720432259356')
	//   expect(amountOut).to.eq(amount.mul(-1))
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

	//   const priceAfterWholeOutputAmount = await sqrtPriceMath.getNextSqrtPriceFromOutput(
	//     price,
	//     liquidity,
	//     amount.mul(-1),
	//     zeroForOne
	//   )
	priceAfterWholeOutputAmount := sqrtPriceMath.GetNextSqrtPriceFromOutput(price, liquidity, new(big.Int).Mul(amount, big.NewInt(-1)), zeroForOne)

	//   expect(sqrtQ, 'price does not reach price target').to.be.lt(priceTarget)
	//   expect(sqrtQ, 'price is less than price after whole output amount').to.eq(priceAfterWholeOutputAmount)
	if (sqrtQ.Cmp(priceTarget) >= 0) || (sqrtQ.Cmp(priceAfterWholeOutputAmount) != 0) {
		t.Errorf("TestComputeSwapStep4 failed, got: %v, %v, expected: %v, %v", sqrtQ, priceAfterWholeOutputAmount, priceTarget, priceAfterWholeOutputAmount)
	}
}

func TestComputeSwapStep5(t *testing.T) {
	fmt.Println("Amount out is capped at the desired amount out")
	//   const { amountIn, amountOut, sqrtQ, feeAmount } = await swapMath.computeSwapStep(
	//     BigNumber.from('417332158212080721273783715441582'),
	//     BigNumber.from('1452870262520218020823638996'),
	//     '159344665391607089467575320103',
	//     '-1',
	//     1
	//   )
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

	//   expect(amountIn).to.eq('1')
	//   expect(feeAmount).to.eq('1')
	//   expect(amountOut).to.eq('1') // would be 2 if not capped
	//   expect(sqrtQ).to.eq('417332158212080721273783715441581')
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

	// it('target price of 1 uses partial input amount', async () => {
	//   const { amountIn, amountOut, sqrtQ, feeAmount } = await swapMath.computeSwapStep(
	//     BigNumber.from('2'),
	//     BigNumber.from('1'),
	//     '1',
	//     '3915081100057732413702495386755767',
	//     1
	//   )
	amountOutExpected := new(big.Int)
	amountOutExpected.SetString("3915081100057732413702495386755767", 10)
	sqrtQ, amountIn, amountOut, feeAmount := ComputeSwapStep(big.NewInt(2), big.NewInt(1), big.NewInt(1), amountOutExpected, 1)

	//   expect(amountIn).to.eq('39614081257132168796771975168')
	//   expect(feeAmount).to.eq('39614120871253040049813')
	//   expect(amountIn.add(feeAmount)).to.be.lte('3915081100057732413702495386755767')
	//   expect(amountOut).to.eq('0')
	//   expect(sqrtQ).to.eq('1')
	// })
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

	//   const { amountIn, amountOut, sqrtQ, feeAmount } = await swapMath.computeSwapStep(
	//     '2413',
	//     '79887613182836312',
	//     '1985041575832132834610021537970',
	//     '10',
	//     1872
	//   )
	price := big.NewInt(2413)
	priceTarget := new(big.Int)
	priceTarget.SetString("79887613182836312", 10)
	liquidity := new(big.Int)
	liquidity.SetString("1985041575832132834610021537970", 10)
	amount := big.NewInt(10)
	fee := 1872

	sqrtQ, amountIn, amountOut, feeAmount := ComputeSwapStep(price, priceTarget, liquidity, amount, fee)

	//   expect(amountIn).to.eq('0')
	//   expect(feeAmount).to.eq('10')
	//   expect(amountOut).to.eq('0')
	//   expect(sqrtQ).to.eq('2413')
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

	//   const sqrtP = BigNumber.from('20282409603651670423947251286016')
	//   const sqrtPTarget = sqrtP.mul(11).div(10)
	//   const liquidity = 1024
	//   // virtual reserves of one are only 4
	//   // https://www.wolframalpha.com/input/?i=1024+%2F+%2820282409603651670423947251286016+%2F+2**96%29
	//   const amountRemaining = -4
	//   const feePips = 3000
	sqrtP := new(big.Int)
	sqrtP.SetString("20282409603651670423947251286016", 10)
	sqrtPTarget := new(big.Int).Div(new(big.Int).Mul(sqrtP, big.NewInt(11)), big.NewInt(10))
	liquidity := big.NewInt(1024)
	amountRemaining := big.NewInt(-4)
	feePips := 3000

	//   const { amountIn, amountOut, sqrtQ, feeAmount } = await swapMath.computeSwapStep(
	//     sqrtP,
	//     sqrtPTarget,
	//     liquidity,
	//     amountRemaining,
	//     feePips
	//   )
	sqrtQ, amountIn, amountOut, feeAmount := ComputeSwapStep(sqrtP, sqrtPTarget, liquidity, amountRemaining, feePips)

	//   expect(amountOut).to.eq(0)
	//   expect(sqrtQ).to.eq(sqrtPTarget)
	//   expect(amountIn).to.eq(26215)
	//   expect(feeAmount).to.eq(79)
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

	//   const sqrtP = BigNumber.from('20282409603651670423947251286016')
	//   const sqrtPTarget = sqrtP.mul(9).div(10)
	//   const liquidity = 1024
	//   // virtual reserves of zero are only 262144
	//   // https://www.wolframalpha.com/input/?i=1024+*+%2820282409603651670423947251286016+%2F+2**96%29
	//   const amountRemaining = -263000
	//   const feePips = 3000
	sqrtP := new(big.Int)
	sqrtP.SetString("20282409603651670423947251286016", 10)
	sqrtPTarget := new(big.Int).Div(new(big.Int).Mul(sqrtP, big.NewInt(9)), big.NewInt(10))
	liquidity := big.NewInt(1024)
	amountRemaining := big.NewInt(-263000)
	feePips := 3000

	//   const { amountIn, amountOut, sqrtQ, feeAmount } = await swapMath.computeSwapStep(
	//     sqrtP,
	//     sqrtPTarget,
	//     liquidity,
	//     amountRemaining,
	//     feePips
	//   )
	sqrtQ, amountIn, amountOut, feeAmount := ComputeSwapStep(sqrtP, sqrtPTarget, liquidity, amountRemaining, feePips)

	//   expect(amountOut).to.eq(26214)
	//   expect(sqrtQ).to.eq(sqrtPTarget)
	//   expect(amountIn).to.eq(1)
	//   expect(feeAmount).to.eq(1)
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
