// The v2 reinvesting strategy is the similar to the v2 strategy, but it
// reinvests the fees collected from the position it holds at updateInterval.
// This code is not finished and has not been tested.
package strategy

import (
	"math/big"

	"github.com/chris-aubin/Uniswap-Simulator/src/libraries/pool"
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
