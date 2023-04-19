package liquidityAmounts

import (
	"fmt"
	"math/big"
	"testing"

	"github.com/chris-aubin/Uniswap-Simulator/src/libraries/utilities"
)

func TestGetLiquidityForAmounts1(t *testing.T) {
	fmt.Println("Amounts for price inside")
	sqrtPriceX96 := utilities.EncodePriceSqrt(big.NewInt(1), big.NewInt(1))
	sqrtPriceAX96 := utilities.EncodePriceSqrt(big.NewInt(100), big.NewInt(110))
	sqrtPriceBX96 := utilities.EncodePriceSqrt(big.NewInt(110), big.NewInt(100))
	liquidity := GetLiquidityForAmounts(sqrtPriceX96, sqrtPriceAX96, sqrtPriceBX96, big.NewInt(100), big.NewInt(200))
	expected := big.NewInt(2148)
	if liquidity.Cmp(expected) != 0 {
		t.Errorf("Expected %v, got %v", expected, liquidity)
	}
}

func TestGetLiquidityForAmounts2(t *testing.T) {
	fmt.Println("Amounts for price below")
	sqrtPriceX96 := utilities.EncodePriceSqrt(big.NewInt(99), big.NewInt(110))
	sqrtPriceAX96 := utilities.EncodePriceSqrt(big.NewInt(100), big.NewInt(110))
	sqrtPriceBX96 := utilities.EncodePriceSqrt(big.NewInt(110), big.NewInt(100))
	liquidity := GetLiquidityForAmounts(sqrtPriceX96, sqrtPriceAX96, sqrtPriceBX96, big.NewInt(100), big.NewInt(200))
	expected := big.NewInt(1048)
	if liquidity.Cmp(expected) != 0 {
		t.Errorf("Expected %v, got %v", expected, liquidity)
	}
}

func TestGetLiquidityForAmounts3(t *testing.T) {
	fmt.Println("Amounts for price above")
	sqrtPriceX96 := utilities.EncodePriceSqrt(big.NewInt(111), big.NewInt(100))
	sqrtPriceAX96 := utilities.EncodePriceSqrt(big.NewInt(100), big.NewInt(110))
	sqrtPriceBX96 := utilities.EncodePriceSqrt(big.NewInt(110), big.NewInt(100))
	liquidity := GetLiquidityForAmounts(sqrtPriceX96, sqrtPriceAX96, sqrtPriceBX96, big.NewInt(100), big.NewInt(200))
	expected := big.NewInt(2097)
	if liquidity.Cmp(expected) != 0 {
		t.Errorf("Expected %v, got %v", expected, liquidity)
	}
}

func TestGetLiquidityForAmounts4(t *testing.T) {
	fmt.Println("Amounts for price equal to lower boundary")
	sqrtPriceAX96 := utilities.EncodePriceSqrt(big.NewInt(100), big.NewInt(110))
	sqrtPriceX96 := sqrtPriceAX96
	sqrtPriceBX96 := utilities.EncodePriceSqrt(big.NewInt(110), big.NewInt(100))
	liquidity := GetLiquidityForAmounts(sqrtPriceX96, sqrtPriceAX96, sqrtPriceBX96, big.NewInt(100), big.NewInt(200))
	expected := big.NewInt(1048)
	if liquidity.Cmp(expected) != 0 {
		t.Errorf("Expected %v, got %v", expected, liquidity)
	}
}

func TestGetLiquidityForAmounts5(t *testing.T) {
	fmt.Println("Amounts for price equal to upper boundary")
	sqrtPriceAX96 := utilities.EncodePriceSqrt(big.NewInt(100), big.NewInt(110))
	sqrtPriceBX96 := utilities.EncodePriceSqrt(big.NewInt(110), big.NewInt(100))
	sqrtPriceX96 := sqrtPriceBX96
	liquidity := GetLiquidityForAmounts(sqrtPriceX96, sqrtPriceAX96, sqrtPriceBX96, big.NewInt(100), big.NewInt(200))
	expected := big.NewInt(2097)
	if liquidity.Cmp(expected) != 0 {
		t.Errorf("Expected %v, got %v", expected, liquidity)
	}
}

func TestGetAmountsForLiquidity1(t *testing.T) {
	fmt.Println("Amounts for price inside")
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

func TestGetAmountsForLiquidity2(t *testing.T) {
	fmt.Println("Amounts for price below")
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

func TestGetAmountsForLiquidity3(t *testing.T) {
	fmt.Println("Amounts for price above")
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

func TestGetAmountsForLiquidity4(t *testing.T) {
	fmt.Println("Amounts for price on lower boundary")
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

func TestGetAmountsForLiquidity5(t *testing.T) {
	fmt.Println("Amounts for price on upper boundary")
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
