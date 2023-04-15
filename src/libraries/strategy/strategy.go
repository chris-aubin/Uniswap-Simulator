package strategy

import (
	"math/big"

	"github.com/chris-aubin/Uniswap-Simulator/src/libraries/pool"
	"github.com/chris-aubin/Uniswap-Simulator/src/libraries/simulation"
)

// StrategyPosition represents a position held by a strategy.
// Positions are indexed using the owner's address, tickLower, and tickUpper.
type StrategyPosition struct {
	TickLower int
	TickUpper int
	Amount    *big.Int
}

type Strategy struct {
	// Address of the strategy
	Address string
	// Amount of token0 in the strategy
	Amount0 *big.Int
	// Amount of token1 in the strategy
	Amount1 *big.Int
	// Positions held by the strategy
	Positions []*StrategyPosition
	// Average gas price for mints, burns and swaps
	GasAvs *simulation.GasAvs
	// Initialize the strategy with the given amounts
	Init func(*big.Int, *big.Int, *pool.Pool, *simulation.GasAvs) *Strategy
	// Adjust strategy positions positions (called at update interval)
	Rebalance func(*pool.Pool)
	// Burn all positions and collect all tokens
	BurnAll func(*pool.Pool) (*big.Int, *big.Int)
}
