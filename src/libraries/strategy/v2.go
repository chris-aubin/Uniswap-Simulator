package strategy

import (
	"math/big"

	"github.com/chris-aubin/Uniswap-Simulator/src/libraries/constants"
	"github.com/chris-aubin/Uniswap-Simulator/src/libraries/liquidityAmounts"
	"github.com/chris-aubin/Uniswap-Simulator/src/libraries/pool"
	"github.com/chris-aubin/Uniswap-Simulator/src/libraries/tickMath"
)

func V2StrategyRebalance(p *pool.Pool, s *Strategy) {
	if len(s.Positions) == 0 {
		V2StrategyMintPosition(p, s)
	}
	return
}

func V2StrategyMintPosition(p *pool.Pool, s *Strategy) {
	s.GasUsed = new(big.Int).Add(s.GasUsed, s.GasAvs.MintGas)
	tickLower := constants.MinTick
	tickUpper := constants.MaxTick
	sqrtRatioAX96 := tickMath.GetSqrtRatioAtTick(tickLower)
	sqrtRatioBX96 := tickMath.GetSqrtRatioAtTick(tickUpper)

	liquidity := liquidityAmounts.GetLiquidityForAmounts(p.Slot0.SqrtPriceX96, sqrtRatioAX96, sqrtRatioBX96, s.Amount0, s.Amount1)
	// amount0, amount1 := liquidityAmounts.GetAmountsForLiquidity(p.Slot0.SqrtPriceX96, sqrtRatioAX96, sqrtRatioBX96, liquidity)
	// fmt.Println("V2StrategyMintPosition")
	// fmt.Println("liquidity", liquidity)
	// fmt.Println("amount0", amount0)
	// fmt.Println("amount1", amount1)
	if liquidity.Cmp(big.NewInt(0)) <= 0 {
		return
	}
	// amount0mint, amount1mint := p.Mint(s.Address, tickLower, tickUpper, liquidity)
	// fmt.Println("amount0mint", amount0mint)
	// fmt.Println("amount1mint", amount1mint)
	p.Mint(s.Address, tickLower, tickUpper, liquidity)
	s.Positions = append(s.Positions, &StrategyPosition{
		TickLower: tickLower,
		TickUpper: tickUpper,
		Liquidity: liquidity,
	})
}
