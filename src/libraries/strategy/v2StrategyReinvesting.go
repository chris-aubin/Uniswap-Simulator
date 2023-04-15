package strategy

import (
	"math/big"

	"github.com/chris-aubin/Uniswap-Simulator/src/libraries/constants"
	"github.com/chris-aubin/Uniswap-Simulator/src/libraries/pool"
	"github.com/chris-aubin/Uniswap-Simulator/src/libraries/tickMath"
	"github.com/chris-aubin/Uniswap-Simulator/src/libraries/utilities"
)

type V2StrategyReinvesting struct {
	Address      string
	Amount0      *big.Int
	Amount1      *big.Int
	GasAllowance *big.Int
	GasAvs       *GasAvs
	Positions    []*StrategyPosition
}

func (s *V2StrategyReinvesting) Rebalance(p *pool.Pool) {
	return
}

func (s *V2StrategyReinvesting) BurnAll(p *pool.Pool) (amount0, amount1 *big.Int) {
	poolTemp := p
	for _, stratPos := range s.Positions {
		poolTemp.Burn(s.Address, stratPos.TickLower, stratPos.TickUpper, stratPos.Amount)
		s.GasAllowance = new(big.Int).Sub(s.GasAllowance, s.GasAvs.BurnGas)
		amount0, amount1 := poolTemp.Collect(s.Address, stratPos.TickLower, stratPos.TickUpper, constants.MaxUint256, constants.MaxUint256)
		s.GasAllowance = new(big.Int).Sub(s.GasAllowance, s.GasAvs.CollectGas)
		s.Amount0 = new(big.Int).Add(s.Amount0, amount0)
		s.Amount1 = new(big.Int).Add(s.Amount1, amount1)
	}
	amount0 = new(big.Int).Set(s.Amount0)
	amount1 = new(big.Int).Set(s.Amount1)
	s.Positions = *new([]*StrategyPosition)
	return
}

func (s *V2StrategyReinvesting) mintPosition(p *pool.Pool) {
	s.GasAllowance = new(big.Int).Sub(s.GasAllowance, s.GasAvs.MintGas)
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

	// amount0, amount1 := s.Pool.MintStrategy(tickLower, tickUpper, amount)
	// s.Amount0.Sub(s.Amount0, amount0)
	// s.Amount1.Sub(s.Amount1, amount1)
}

func Make(amount0, amount1 *big.Int, p *pool.Pool, g *GasAvs) {
	s := new(V2StrategyReinvesting)
	s.Address = "0x0000000000000000000000000000000000000001"
	s.Amount1 = new(big.Int).Set(amount0)
	s.Amount1 = new(big.Int).Set(amount1)
	s.GasAvs = g
	s.Positions = make([]*StrategyPosition, 0)

	s.mintPosition(p)
	return
}
