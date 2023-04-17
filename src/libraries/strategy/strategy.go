package strategy

import (
	"fmt"
	"math/big"

	"github.com/chris-aubin/Uniswap-Simulator/src/libraries/constants"
	"github.com/chris-aubin/Uniswap-Simulator/src/libraries/pool"
)

var strategies map[string]func(p *pool.Pool, s *Strategy)

func init() {
	strategies = make(map[string]func(p *pool.Pool, s *Strategy))
	strategies["v2"] = V2StrategyRebalance
	strategies["v2Reinvesting"] = V2StrategyReinvestingRebalance
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
	MintGas *big.Int `json:"mintAv"`
	// Av. gas required to burn a position
	BurnGas *big.Int `json:"burnAv"`
	// Av. gas required to swap
	SwapGas *big.Int `json:"swapAv"`
	// Av. gas required to flash (likely unnecessary for strategies)
	FlashGas *big.Int `json:"flashAv"`
	// Av. gas required to collect fees from a position
	CollectGas *big.Int `json:"collectAv"`
}

type Strategy struct {
	Address        string
	Amount0        *big.Int
	Amount1        *big.Int
	GasUsed        *big.Int
	GasAvs         *GasAvs
	UpdateInterval int
	Positions      []*StrategyPosition
	Rebalance      func(p *pool.Pool, s *Strategy)
}

func (s *Strategy) Init(p *pool.Pool) {
	s.Rebalance(p, s)
}

func (s *Strategy) BurnAll(p *pool.Pool) (amount0, amount1 *big.Int) {
	fmt.Println("HERE")
	for _, stratPos := range s.Positions {
		fmt.Println("IN LOOP")
		p.Burn(s.Address, stratPos.TickLower, stratPos.TickUpper, stratPos.Amount)
		s.GasUsed = new(big.Int).Add(s.GasUsed, s.GasAvs.BurnGas)
		amount0, amount1 := p.Collect(s.Address, stratPos.TickLower, stratPos.TickUpper, constants.MaxUint256, constants.MaxUint256)
		fmt.Println("BURN ALL")
		fmt.Println("amount0", amount0)
		fmt.Println("amount1", amount1)
		s.GasUsed = new(big.Int).Add(s.GasUsed, s.GasAvs.CollectGas)
		s.Amount0 = new(big.Int).Add(s.Amount0, amount0)
		s.Amount1 = new(big.Int).Add(s.Amount1, amount1)
	}
	amount0 = new(big.Int).Set(s.Amount0)
	amount1 = new(big.Int).Set(s.Amount1)
	s.Positions = *new([]*StrategyPosition)
	return
}

func (s *Strategy) Results(p *pool.Pool) (*big.Int, *big.Int, *big.Int) {
	amount0temp, amount1temp := s.BurnAll(p)
	s.Amount0 = new(big.Int).Set(amount0temp)
	s.Amount1 = new(big.Int).Set(amount1temp)
	return s.Amount0, s.Amount1, s.GasUsed
}

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
