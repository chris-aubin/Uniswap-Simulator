package utilities

import (
	"math/big"

	"github.com/chris-aubin/Uniswap-Simulator/src/libraries/constants"
	"github.com/chris-aubin/Uniswap-Simulator/src/libraries/fullMath"
)

func EncodePriceSqrt(reserve1, reserve0 *big.Int) *big.Int {
	r1 := new(big.Float).SetInt(reserve1)
	r0 := new(big.Float).SetInt(reserve0)
	price := new(big.Float).Quo(r1, r0)
	priceSqrt := new(big.Float).Sqrt(price)
	// Convert to Q96
	priceSqrtQ96 := new(big.Float).Mul(priceSqrt, new(big.Float).SetInt(constants.Q96))
	priceSqrtQ96Int, _ := priceSqrtQ96.Int(nil)
	// Check rounding
	priceSqrtQ96IntFloat := new(big.Float).SetInt(priceSqrtQ96Int)
	if priceSqrtQ96.Cmp(priceSqrtQ96IntFloat) >= 1 {
		priceSqrtQ96Int.Sub(priceSqrtQ96Int, big.NewInt(1))
	}

	return priceSqrtQ96Int
}

func getLiquidityForAmount0(sqrtRatioAX96, sqrtRatioBX96, amount0 *big.Int) *big.Int {
	if sqrtRatioAX96.Cmp(sqrtRatioBX96) > 0 {
		sqrtRatioAX96, sqrtRatioBX96 = sqrtRatioBX96, sqrtRatioAX96
	}
	intermediate := fullMath.MulDiv(sqrtRatioAX96, sqrtRatioBX96, constants.Q96)
	return fullMath.MulDiv(amount0, intermediate, new(big.Int).Sub(sqrtRatioBX96, sqrtRatioAX96))
}

func getLiquidityForAmount1(sqrtRatioAX96, sqrtRatioBX96, amount1 *big.Int) *big.Int {
	if sqrtRatioAX96.Cmp(sqrtRatioBX96) > 0 {
		sqrtRatioAX96, sqrtRatioBX96 = sqrtRatioBX96, sqrtRatioAX96
	}
	return fullMath.MulDiv(amount1, constants.Q96, new(big.Int).Sub(sqrtRatioBX96, sqrtRatioAX96))
}

func GetLiquidityForAmount(sqrtRatioX96, sqrtRatioAX96, sqrtRatioBX96, amount0, amount1 *big.Int) (liquidity *big.Int) {
	if sqrtRatioAX96.Cmp(sqrtRatioBX96) > 0 {
		sqrtRatioAX96, sqrtRatioBX96 = sqrtRatioBX96, sqrtRatioAX96
	}
	if sqrtRatioX96.Cmp(sqrtRatioAX96) <= 0 {
		liquidity = getLiquidityForAmount0(sqrtRatioAX96, sqrtRatioBX96, amount0)
	} else if sqrtRatioX96.Cmp(sqrtRatioBX96) < 0 {
		liquidity0 := getLiquidityForAmount0(sqrtRatioX96, sqrtRatioBX96, amount0)
		liquidity1 := getLiquidityForAmount1(sqrtRatioAX96, sqrtRatioX96, amount1)

		if liquidity0.Cmp(liquidity1) < 0 {
			liquidity = liquidity0
		} else {
			liquidity = liquidity1
		}
	} else {
		liquidity = getLiquidityForAmount1(sqrtRatioAX96, sqrtRatioBX96, amount1)
	}
	return liquidity
}

func getAmount0ForLiquidity(sqrtRatioAX96, sqrtRatioBX96, liquidity *big.Int) *big.Int {
	if sqrtRatioAX96.Cmp(sqrtRatioBX96) > 0 {
		sqrtRatioAX96, sqrtRatioBX96 = sqrtRatioBX96, sqrtRatioAX96
	}
	temp1 := new(big.Int).Lsh(liquidity, 96)
	temp2 := new(big.Int).Sub(sqrtRatioBX96, sqrtRatioAX96)
	temp3 := fullMath.MulDiv(temp1, temp2, sqrtRatioBX96)
	return new(big.Int).Div(temp3, sqrtRatioAX96)
}

func getAmount1ForLiquidity(sqrtRatioAX96, sqrtRatioBX96, liquidity *big.Int) *big.Int {
	if sqrtRatioAX96.Cmp(sqrtRatioBX96) > 0 {
		sqrtRatioAX96, sqrtRatioBX96 = sqrtRatioBX96, sqrtRatioAX96
	}
	temp1 := new(big.Int).Sub(sqrtRatioBX96, sqrtRatioAX96)
	return fullMath.MulDiv(liquidity, temp1, constants.Q96)
}

func GetAmountsForLiquidity(sqrtRatioX96, sqrtRatioAX96, sqrtRatioBX96, liquidity *big.Int) (*big.Int, *big.Int) {
	if sqrtRatioAX96.Cmp(sqrtRatioBX96) > 0 {
		sqrtRatioAX96, sqrtRatioBX96 = sqrtRatioBX96, sqrtRatioAX96
	}
	amount0, amount1 := new(big.Int), new(big.Int)
	if sqrtRatioX96.Cmp(sqrtRatioAX96) <= 0 {
		amount0 = getAmount0ForLiquidity(sqrtRatioAX96, sqrtRatioBX96, liquidity)
	} else if sqrtRatioX96.Cmp(sqrtRatioBX96) < 0 {
		amount0 = getAmount0ForLiquidity(sqrtRatioX96, sqrtRatioBX96, liquidity)
		amount1 = getAmount1ForLiquidity(sqrtRatioAX96, sqrtRatioX96, liquidity)
	} else {
		amount1 = getAmount1ForLiquidity(sqrtRatioAX96, sqrtRatioBX96, liquidity)
	}
	return amount0, amount1
}
