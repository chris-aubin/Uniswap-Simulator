package liquidityAmounts

import (
	"math/big"
	"testing"

	"github.com/chris-aubin/Uniswap-Simulator/src/libraries/utilities"
)

//   describe('#getLiquidityForAmounts', () => {
//     it('amounts for price inside', async () => {
//       const sqrtPriceX96 = encodePriceSqrt(1, 1)
//       const sqrtPriceAX96 = encodePriceSqrt(100, 110)
//       const sqrtPriceBX96 = encodePriceSqrt(110, 100)
//       const liquidity = await liquidityFromAmounts.getLiquidityForAmounts(
//         sqrtPriceX96,
//         sqrtPriceAX96,
//         sqrtPriceBX96,
//         100,
//         200
//       )
//       expect(liquidity).to.eq(2148)
//     })

func TestGetLiquidityForAmounts1(t *testing.T) {
	sqrtPriceX96 := utilities.EncodePriceSqrt(big.NewInt(1), big.NewInt(1))
	sqrtPriceAX96 := utilities.EncodePriceSqrt(big.NewInt(100), big.NewInt(110))
	sqrtPriceBX96 := utilities.EncodePriceSqrt(big.NewInt(110), big.NewInt(100))
	liquidity := GetLiquidityForAmounts(sqrtPriceX96, sqrtPriceAX96, sqrtPriceBX96, big.NewInt(100), big.NewInt(200))
	expected := big.NewInt(2148)
	if liquidity.Cmp(expected) != 0 {
		t.Errorf("Expected %v, got %v", expected, liquidity)
	}
}

//     it('amounts for price below', async () => {
//       const sqrtPriceX96 = encodePriceSqrt(99, 110)
//       const sqrtPriceAX96 = encodePriceSqrt(100, 110)
//       const sqrtPriceBX96 = encodePriceSqrt(110, 100)
//       const liquidity = await liquidityFromAmounts.getLiquidityForAmounts(
//         sqrtPriceX96,
//         sqrtPriceAX96,
//         sqrtPriceBX96,
//         100,
//         200
//       )
//       expect(liquidity).to.eq(1048)
//     })

func TestGetLiquidityForAmounts2(t *testing.T) {
	sqrtPriceX96 := utilities.EncodePriceSqrt(big.NewInt(99), big.NewInt(110))
	sqrtPriceAX96 := utilities.EncodePriceSqrt(big.NewInt(100), big.NewInt(110))
	sqrtPriceBX96 := utilities.EncodePriceSqrt(big.NewInt(110), big.NewInt(100))
	liquidity := GetLiquidityForAmounts(sqrtPriceX96, sqrtPriceAX96, sqrtPriceBX96, big.NewInt(100), big.NewInt(200))
	expected := big.NewInt(1048)
	if liquidity.Cmp(expected) != 0 {
		t.Errorf("Expected %v, got %v", expected, liquidity)
	}
}

//     it('amounts for price above', async () => {
//       const sqrtPriceX96 = encodePriceSqrt(111, 100)
//       const sqrtPriceAX96 = encodePriceSqrt(100, 110)
//       const sqrtPriceBX96 = encodePriceSqrt(110, 100)
//       const liquidity = await liquidityFromAmounts.getLiquidityForAmounts(
//         sqrtPriceX96,
//         sqrtPriceAX96,
//         sqrtPriceBX96,
//         100,
//         200
//       )
//       expect(liquidity).to.eq(2097)
//     })

func TestGetLiquidityForAmounts3(t *testing.T) {
	sqrtPriceX96 := utilities.EncodePriceSqrt(big.NewInt(111), big.NewInt(100))
	sqrtPriceAX96 := utilities.EncodePriceSqrt(big.NewInt(100), big.NewInt(110))
	sqrtPriceBX96 := utilities.EncodePriceSqrt(big.NewInt(110), big.NewInt(100))
	liquidity := GetLiquidityForAmounts(sqrtPriceX96, sqrtPriceAX96, sqrtPriceBX96, big.NewInt(100), big.NewInt(200))
	expected := big.NewInt(2097)
	if liquidity.Cmp(expected) != 0 {
		t.Errorf("Expected %v, got %v", expected, liquidity)
	}
}

//     it('amounts for price equal to lower boundary', async () => {
//       const sqrtPriceAX96 = encodePriceSqrt(100, 110)
//       const sqrtPriceX96 = sqrtPriceAX96
//       const sqrtPriceBX96 = encodePriceSqrt(110, 100)
//       const liquidity = await liquidityFromAmounts.getLiquidityForAmounts(
//         sqrtPriceX96,
//         sqrtPriceAX96,
//         sqrtPriceBX96,
//         100,
//         200
//       )
//       expect(liquidity).to.eq(1048)
//     })

func TestGetLiquidityForAmounts4(t *testing.T) {
	sqrtPriceAX96 := utilities.EncodePriceSqrt(big.NewInt(100), big.NewInt(110))
	sqrtPriceX96 := sqrtPriceAX96
	sqrtPriceBX96 := utilities.EncodePriceSqrt(big.NewInt(110), big.NewInt(100))
	liquidity := GetLiquidityForAmounts(sqrtPriceX96, sqrtPriceAX96, sqrtPriceBX96, big.NewInt(100), big.NewInt(200))
	expected := big.NewInt(1048)
	if liquidity.Cmp(expected) != 0 {
		t.Errorf("Expected %v, got %v", expected, liquidity)
	}
}

//     it('amounts for price equal to upper boundary', async () => {
//       const sqrtPriceAX96 = encodePriceSqrt(100, 110)
//       const sqrtPriceBX96 = encodePriceSqrt(110, 100)
//       const sqrtPriceX96 = sqrtPriceBX96
//       const liquidity = await liquidityFromAmounts.getLiquidityForAmounts(
//         sqrtPriceX96,
//         sqrtPriceAX96,
//         sqrtPriceBX96,
//         100,
//         200
//       )
//       expect(liquidity).to.eq(2097)
//     })

func TestGetLiquidityForAmounts5(t *testing.T) {
	sqrtPriceAX96 := utilities.EncodePriceSqrt(big.NewInt(100), big.NewInt(110))
	sqrtPriceBX96 := utilities.EncodePriceSqrt(big.NewInt(110), big.NewInt(100))
	sqrtPriceX96 := sqrtPriceBX96
	liquidity := GetLiquidityForAmounts(sqrtPriceX96, sqrtPriceAX96, sqrtPriceBX96, big.NewInt(100), big.NewInt(200))
	expected := big.NewInt(2097)
	if liquidity.Cmp(expected) != 0 {
		t.Errorf("Expected %v, got %v", expected, liquidity)
	}
}

//   describe('#getAmountsForLiquidity', () => {
//     it('amounts for price inside', async () => {
//       const sqrtPriceX96 = encodePriceSqrt(1, 1)
//       const sqrtPriceAX96 = encodePriceSqrt(100, 110)
//       const sqrtPriceBX96 = encodePriceSqrt(110, 100)
//       const { amount0, amount1 } = await liquidityFromAmounts.getAmountsForLiquidity(
//         sqrtPriceX96,
//         sqrtPriceAX96,
//         sqrtPriceBX96,
//         2148
//       )
//       expect(amount0).to.eq(99)
//       expect(amount1).to.eq(99)
//     })

func TestGetAmountsForLiquidity1(t *testing.T) {
	sqrtPriceX96 := utilities.EncodePriceSqrt(big.NewInt(1), big.NewInt(1))
	sqrtPriceAX96 := utilities.EncodePriceSqrt(big.NewInt(100), big.NewInt(110))
	sqrtPriceBX96 := utilities.EncodePriceSqrt(big.NewInt(110), big.NewInt(100))
	amount0, amount1 := GetAmountsForLiquidity(sqrtPriceX96, sqrtPriceAX96, sqrtPriceBX96, big.NewInt(2148))
	expected0 := big.NewInt(99)
	expected1 := big.NewInt(99)
	if amount0.Cmp(expected0) != 0 {
		t.Errorf("Expected %v, got %v", expected0, amount0)
	}
	if amount1.Cmp(expected1) != 0 {
		t.Errorf("Expected %v, got %v", expected1, amount1)
	}
}

//     it('amounts for price below', async () => {
//       const sqrtPriceX96 = encodePriceSqrt(99, 110)
//       const sqrtPriceAX96 = encodePriceSqrt(100, 110)
//       const sqrtPriceBX96 = encodePriceSqrt(110, 100)
//       const { amount0, amount1 } = await liquidityFromAmounts.getAmountsForLiquidity(
//         sqrtPriceX96,
//         sqrtPriceAX96,
//         sqrtPriceBX96,
//         1048
//       )
//       expect(amount0).to.eq(99)
//       expect(amount1).to.eq(0)
//     })

func TestGetAmountsForLiquidity2(t *testing.T) {
	sqrtPriceX96 := utilities.EncodePriceSqrt(big.NewInt(99), big.NewInt(110))
	sqrtPriceAX96 := utilities.EncodePriceSqrt(big.NewInt(100), big.NewInt(110))
	sqrtPriceBX96 := utilities.EncodePriceSqrt(big.NewInt(110), big.NewInt(100))
	amount0, amount1 := GetAmountsForLiquidity(sqrtPriceX96, sqrtPriceAX96, sqrtPriceBX96, big.NewInt(1048))
	expected0 := big.NewInt(99)
	expected1 := big.NewInt(0)
	if amount0.Cmp(expected0) != 0 {
		t.Errorf("Expected %v, got %v", expected0, amount0)
	}
	if amount1.Cmp(expected1) != 0 {
		t.Errorf("Expected %v, got %v", expected1, amount1)
	}
}

//     it('amounts for price above', async () => {
//       const sqrtPriceX96 = encodePriceSqrt(111, 100)
//       const sqrtPriceAX96 = encodePriceSqrt(100, 110)
//       const sqrtPriceBX96 = encodePriceSqrt(110, 100)
//       const { amount0, amount1 } = await liquidityFromAmounts.getAmountsForLiquidity(
//         sqrtPriceX96,
//         sqrtPriceAX96,
//         sqrtPriceBX96,
//         2097
//       )
//       expect(amount0).to.eq(0)
//       expect(amount1).to.eq(199)
//     })

func TestGetAmountsForLiquidity3(t *testing.T) {
	sqrtPriceX96 := utilities.EncodePriceSqrt(big.NewInt(111), big.NewInt(100))
	sqrtPriceAX96 := utilities.EncodePriceSqrt(big.NewInt(100), big.NewInt(110))
	sqrtPriceBX96 := utilities.EncodePriceSqrt(big.NewInt(110), big.NewInt(100))
	amount0, amount1 := GetAmountsForLiquidity(sqrtPriceX96, sqrtPriceAX96, sqrtPriceBX96, big.NewInt(2097))
	expected0 := big.NewInt(0)
	expected1 := big.NewInt(199)
	if amount0.Cmp(expected0) != 0 {
		t.Errorf("Expected %v, got %v", expected0, amount0)
	}
	if amount1.Cmp(expected1) != 0 {
		t.Errorf("Expected %v, got %v", expected1, amount1)
	}
}

//     it('amounts for price on lower boundary', async () => {
//       const sqrtPriceAX96 = encodePriceSqrt(100, 110)
//       const sqrtPriceX96 = sqrtPriceAX96
//       const sqrtPriceBX96 = encodePriceSqrt(110, 100)
//       const { amount0, amount1 } = await liquidityFromAmounts.getAmountsForLiquidity(
//         sqrtPriceX96,
//         sqrtPriceAX96,
//         sqrtPriceBX96,
//         1048
//       )
//       expect(amount0).to.eq(99)
//       expect(amount1).to.eq(0)
//     })

func TestGetAmountsForLiquidity4(t *testing.T) {
	sqrtPriceAX96 := utilities.EncodePriceSqrt(big.NewInt(100), big.NewInt(110))
	sqrtPriceX96 := sqrtPriceAX96
	sqrtPriceBX96 := utilities.EncodePriceSqrt(big.NewInt(110), big.NewInt(100))
	amount0, amount1 := GetAmountsForLiquidity(sqrtPriceX96, sqrtPriceAX96, sqrtPriceBX96, big.NewInt(1048))
	expected0 := big.NewInt(99)
	expected1 := big.NewInt(0)
	if amount0.Cmp(expected0) != 0 {
		t.Errorf("Expected %v, got %v", expected0, amount0)
	}
	if amount1.Cmp(expected1) != 0 {
		t.Errorf("Expected %v, got %v", expected1, amount1)
	}
}

//     it('amounts for price on upper boundary', async () => {
//       const sqrtPriceAX96 = encodePriceSqrt(100, 110)
//       const sqrtPriceBX96 = encodePriceSqrt(110, 100)
//       const sqrtPriceX96 = sqrtPriceBX96
//       const { amount0, amount1 } = await liquidityFromAmounts.getAmountsForLiquidity(
//         sqrtPriceX96,
//         sqrtPriceAX96,
//         sqrtPriceBX96,
//         2097
//       )
//       expect(amount0).to.eq(0)
//       expect(amount1).to.eq(199)
//     })

func TestGetAmountsForLiquidity5(t *testing.T) {
	sqrtPriceAX96 := utilities.EncodePriceSqrt(big.NewInt(100), big.NewInt(110))
	sqrtPriceBX96 := utilities.EncodePriceSqrt(big.NewInt(110), big.NewInt(100))
	sqrtPriceX96 := sqrtPriceBX96
	amount0, amount1 := GetAmountsForLiquidity(sqrtPriceX96, sqrtPriceAX96, sqrtPriceBX96, big.NewInt(2097))
	expected0 := big.NewInt(0)
	expected1 := big.NewInt(199)
	if amount0.Cmp(expected0) != 0 {
		t.Errorf("Expected %v, got %v", expected0, amount0)
	}
	if amount1.Cmp(expected1) != 0 {
		t.Errorf("Expected %v, got %v", expected1, amount1)
	}
}
