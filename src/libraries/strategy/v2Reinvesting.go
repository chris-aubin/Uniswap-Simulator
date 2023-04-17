package strategy

import (
	"math/big"

	"github.com/chris-aubin/Uniswap-Simulator/src/libraries/constants"
	"github.com/chris-aubin/Uniswap-Simulator/src/libraries/pool"
	"github.com/chris-aubin/Uniswap-Simulator/src/libraries/tickMath"
	"github.com/chris-aubin/Uniswap-Simulator/src/libraries/utilities"
)

func V2StrategyReinvestingRebalance(p *pool.Pool, s *Strategy) {
	// Probably better to set aside a little bit of the pool's liquidity for
	// and, instead of burning all liquidity, mint a little bit (to recalculate
	// tokens owed) and then collect and reinvest the rest.
	amount0temp, amount1temp := s.BurnAll(p)
	s.Amount0 = new(big.Int).Add(s.Amount0, amount0temp)
	s.Amount1 = new(big.Int).Add(s.Amount1, amount1temp)
	V2StrategyMintPosition(p, s)
}

func V2StrategyReinvestingMintPosition(p *pool.Pool, s *Strategy) {
	s.GasUsed = new(big.Int).Add(s.GasUsed, s.GasAvs.MintGas)
	tickLower := constants.MinTick
	tickUpper := constants.MaxTick
	sqrtRatioAX96 := tickMath.GetSqrtRatioAtTick(tickLower)
	sqrtRatioBX96 := tickMath.GetSqrtRatioAtTick(tickUpper)

	amount := utilities.GetLiquidityForAmount(p.Slot0.SqrtPriceX96, sqrtRatioAX96, sqrtRatioBX96, s.Amount0, s.Amount1)
	if amount.Cmp(big.NewInt(0)) <= 0 {
		return
	}
	s.Positions = append(s.Positions, &StrategyPosition{
		TickLower: tickLower,
		TickUpper: tickUpper,
		Amount:    amount,
	})
}
