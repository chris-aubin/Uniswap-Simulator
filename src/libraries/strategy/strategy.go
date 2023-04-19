// The strategy package contains the logic for defining and executing strategies.
//
// All strategies have a BurnAll function that burns all of the strategy's
// positions and calculates the tokens owed to the strategy, a Results function
// that returns the tokens that the strategy has accumulated and the total
// amount of gas that the strategy has spent and a Make function that
// initialises a strategy.
//
// The only field that differs significantly from strategy to strategy is the
// Rebalance function. The function is of type \func(p *pool.Pool, s *Strategy).
// It takes in a Pool and a  Strategy. It can call any of the Pool methods and
// it has access to all of the Pool and Strategy state. It make use of any
// number of helper functions.

package strategy

import (
	"math/big"

	"github.com/chris-aubin/Uniswap-Simulator/src/libraries/constants"
	"github.com/chris-aubin/Uniswap-Simulator/src/libraries/pool"
)

// Map of strategy names to strategy rebalance functions.
var strategies map[string]func(p *pool.Pool, s *Strategy)

// Intialises the strategies map.
func init() {
	strategies = make(map[string]func(p *pool.Pool, s *Strategy))
	strategies["nil"] = NilStrategyRebalance
	strategies["v2"] = V2StrategyRebalance
	strategies["v2Reinvesting"] = V2StrategyReinvestingRebalance
}

// Used to decode strategy input from JSON.
type StrategyInput struct {
	Strategy       string   `json:"strategy"`
	Amount0        *big.Int `json:"amount0"`
	Amount1        *big.Int `json:"amount1"`
	UpdateInterval int      `json:"updateInterval"`
}

// StrategyPosition represents a position held by a strategy.
// Pools index positions using the owner's address, tickLower, and tickUpper, so
// the strategy must keep track of these values (the owners address is the same
// as the strategy address in the Strategy struct).
type StrategyPosition struct {
	TickLower int
	TickUpper int
	Liquidity *big.Int
}

// Used to decode gas averages input from JSON.
type GasAvs struct {
	// Av. gas required to mint a position.
	MintGas *big.Int `json:"mintAv"`
	// Av. gas required to burn a position.
	BurnGas *big.Int `json:"burnAv"`
	// Av. gas required to swap
	SwapGas *big.Int `json:"swapAv"`
	// Av. gas required to flash (likely unnecessary for strategies).
	FlashGas *big.Int `json:"flashAv"`
	// Av. gas required to collect fees from a position.
	CollectGas *big.Int `json:"collectAv"`
}

// Strategy state.
type Strategy struct {
	// Address of the strategy.
	Address string
	// Current amount of token0 and token1 held by the strategy (does NOT
	// include the token0 and token1 deployed in strategies).
	Amount0 *big.Int
	Amount1 *big.Int
	// Total amount of gas used by the strategy.
	GasUsed *big.Int
	// The average gas required to perform each operation during the testing.
	// period
	GasAvs *GasAvs
	// The number of blocks between each rebalance.
	UpdateInterval int
	// The positions held by the strategy
	Positions []*StrategyPosition
	// The function that is called to rebalance the strategy.
	Rebalance func(p *pool.Pool, s *Strategy)
}

// Burns all of the strategy's positions and calculates the tokens owed to the
// strategy.
func (s *Strategy) BurnAll(p *pool.Pool) (amount0, amount1 *big.Int) {
	for _, stratPos := range s.Positions {
		p.Burn(s.Address, stratPos.TickLower, stratPos.TickUpper, stratPos.Liquidity)
		s.GasUsed = new(big.Int).Add(s.GasUsed, s.GasAvs.BurnGas)
		amount0, amount1 := p.Collect(s.Address, stratPos.TickLower, stratPos.TickUpper, constants.MaxUint256, constants.MaxUint256)
		s.GasUsed = new(big.Int).Add(s.GasUsed, s.GasAvs.CollectGas)
		s.Amount0 = new(big.Int).Add(s.Amount0, amount0)
		s.Amount1 = new(big.Int).Add(s.Amount1, amount1)
	}
	amount0 = new(big.Int).Set(s.Amount0)
	amount1 = new(big.Int).Set(s.Amount1)
	s.Positions = *new([]*StrategyPosition)
	return
}

// Returns the tokens that the strategy has accumulated and the total amount of
// gas that the strategy has spent.
func (s *Strategy) Results(p *pool.Pool) (*big.Int, *big.Int, *big.Int) {
	amount0temp, amount1temp := s.BurnAll(p)
	s.Amount0 = new(big.Int).Set(amount0temp)
	s.Amount1 = new(big.Int).Set(amount1temp)
	return s.Amount0, s.Amount1, s.GasUsed
}

// Initialises a strategy.
func Make(amount0, amount1 *big.Int, p *pool.Pool, g *GasAvs, identifier string, updateInterval int) *Strategy {
	s := new(Strategy)
	s.Address = "0x0000000000000000000000000000000000000001"
	s.Amount0 = new(big.Int).Set(amount0)
	s.Amount1 = new(big.Int).Set(amount1)
	s.GasAvs = g
	s.GasUsed = big.NewInt(0)
	s.UpdateInterval = updateInterval
	s.Positions = make([]*StrategyPosition, 0)
	s.Rebalance = strategies[identifier]
	return s
}
