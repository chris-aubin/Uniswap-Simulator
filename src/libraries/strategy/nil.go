// The nil strategy is a strategy that does nothing.
package strategy

import (
	"github.com/chris-aubin/Uniswap-Simulator/src/libraries/pool"
)

// Does nothing.
func NilStrategyRebalance(p *pool.Pool, s *Strategy) {
	return
}
