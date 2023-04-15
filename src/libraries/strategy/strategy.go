package strategy

import (
	"math/big"

	"github.com/chris-aubin/Uniswap-Simulator/src/libraries/pool"
)

strategies := map[string]func(){
	"v2": MakeV2Strategy

}
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

type StrategyMethods interface {
	// Initialize the strategy with the given amounts
	Init (*big.Int, *big.Int, *pool.Pool, *GasAvs)
	// Adjust strategy positions positions (called at update interval)
	Rebalance (*pool.Pool)
	// Burn all positions and collect all tokens
	BurnAll (*pool.Pool) (*big.Int, *big.Int)
}

type Strategy struct {
	Address      string
	Amount0      *big.Int
	Amount1      *big.Int
	GasAllowance *big.Int
	GasAvs       *GasAvs
	Positions    []*StrategyPosition
	StrategyMethods
}
