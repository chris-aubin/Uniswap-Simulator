package strategy

import (
	"math/big"

	"github.com/chris-aubin/Uniswap-Simulator/src/libraries/pool"
)

// StrategyPosition represents a position held by a strategy.
// Positions are indexed using the owner's address, tickLower, and tickUpper.
type StrategyPosition struct {
	TickLower int
	TickUpper int
	Amount    *big.Int
}

type GasAvs struct {
	// Av. gas required to mint a position
	MintGas *big.Int
	// Av. gas required to burn a position
	BurnGas *big.Int
	// Av. gas required to swap
	SwapGas *big.Int
	// Av. gas required to flash (likely unnecessary for strategies)
	FlashGas *big.Int
	// Av. gas required to collect fees from a position
	CollectGas *big.Int
}

type Strategy struct {
	// Address of the strategy
	Address string
	// Amount of token0 in the strategy
	Amount0 *big.Int
	// Amount of token1 in the strategy
	Amount1 *big.Int
	// Update interval for the strategy (how regularly to rebalance) in blocks
	UpdateInterval int
	// Positions held by the strategy
	Positions []*StrategyPosition
	// Average gas price for mints, burns and swaps
	GasAvs *GasAvs
	// Initialize the strategy with the given amounts
	Make func(*big.Int, *big.Int, *pool.Pool, *GasAvs) *Strategy
	// Adjust strategy positions positions (called at update interval)
	Rebalance func(*pool.Pool)
	// Burn all positions and collect all tokens
	BurnAll func(*pool.Pool) (*big.Int, *big.Int)
}
