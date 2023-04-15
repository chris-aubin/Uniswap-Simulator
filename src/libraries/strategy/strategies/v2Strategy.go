package strategy

import (
	"math/big"

	"github.com/chris-aubin/Uniswap-Simulator/src/libraries/constants"
	"github.com/chris-aubin/Uniswap-Simulator/src/libraries/pool"
	"github.com/chris-aubin/Uniswap-Simulator/src/libraries/tickMath"
)

type V2Strategy struct {
	Address   string
	Amount0   *big.Int
	Amount1   *big.Int
	Positions []*StrategyPosition
}

func (s *V2Strategy) Rebalance(p *pool.Pool) {
	s.Amount0, s.Amount1 = s.BurnAll(p)

	s.mintPosition(constants.MinTick, constants.MaxTick)
	return
}

func NewV2Strategy(amount0, amount1 *big.Int) *V2Strategy {
	return new(V2Strategy)
}

func (s *V2Strategy) BurnAll(p *pool.Pool) (amount0, amount1 *big.Int) {
	poolTemp := p
	for _, stratPos := range s.Positions {
		poolTemp.Burn(s.Address, stratPos.TickLower, stratPos.TickUpper, stratPos.Amount)
		amount0, amount1 := poolTemp.Collect(s.Address, stratPos.TickLower, stratPos.TickUpper, constants.MaxUint256, constants.MaxUint256)
		s.Amount0 = new(big.Int).Add(s.Amount0, amount0)
		s.Amount1 = new(big.Int).Add(s.Amount1, amount1)
	}
	amount0 = new(big.Int).Set(s.Amount0)
	amount1 = new(big.Int).Set(s.Amount1)
	s.Positions = *new([]*StrategyPosition)
	return
}

func (s *V2Strategy) mintPosition(tickLower, tickUpper int, p *pool.Pool) {
	sqrtRatioAX96 := tickMath.GetSqrtRatioAtTick(tickLower)
	sqrtRatioBX96 := tickMath.GetSqrtRatioAtTick(tickUpper)

	// amount := la.GetLiquidityForAmount(s.Pool.SqrtRatioX96, sqrtRatioAX96, sqrtRatioBX96, s.Amount0, s.Amount1)
	// if amount.IsZero() {
	// 	return
	// }
	// s.Positions = append(s.Positions, StrategyPosition{
	// 	amount:    amount,
	// 	tickLower: tickLower,
	// 	tickUpper: tickUpper,
	// })

	// amount0, amount1 := s.Pool.MintStrategy(tickLower, tickUpper, amount)
	// s.Amount0.Sub(s.Amount0, amount0)
	// s.Amount1.Sub(s.Amount1, amount1)
}

func (s *V2Strategy) Init() (amount0, amount1 *big.Int) {
	s.Amount1 = new(big.Int).Set(amount0)
	s.Amount1 = new(big.Int).Set(amount1)

	// New Positions
	tickLower := -887270
	tickUpper := -tickLower
	// s.mintPosition(tickLower, tickUpper)
	return
}
