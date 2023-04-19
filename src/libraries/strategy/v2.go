// The v2 strategy is the simplest strategy. It simply mints a Uniswap v2 style
// position that provides the maximum amount of liquidity over the entire tick
// range.
package strategy

import (
	"math/big"

	"github.com/chris-aubin/Uniswap-Simulator/src/libraries/constants"
	"github.com/chris-aubin/Uniswap-Simulator/src/libraries/liquidityAmounts"
	"github.com/chris-aubin/Uniswap-Simulator/src/libraries/pool"
	"github.com/chris-aubin/Uniswap-Simulator/src/libraries/tickMath"
)

func V2StrategyRebalance(p *pool.Pool, s *Strategy) {
	// Only rebalance once, when the strategy is first created.
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

	// Calculate the amount of liquidity to mint.
	liquidity := liquidityAmounts.GetLiquidityForAmounts(p.Slot0.SqrtPriceX96, sqrtRatioAX96, sqrtRatioBX96, s.Amount0, s.Amount1)
	if liquidity.Cmp(big.NewInt(0)) <= 0 {
		return
	}

	// Mint the position.
	p.Mint(s.Address, tickLower, tickUpper, liquidity)
	// Add the position to the strategy's positions slice.
	s.Positions = append(s.Positions, &StrategyPosition{
		TickLower: tickLower,
		TickUpper: tickUpper,
		Liquidity: liquidity,
	})
}
