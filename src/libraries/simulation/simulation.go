package simulation

// package executor

import (
	"math/big"

	"github.com/chris-aubin/Uniswap-Simulator/src/libraries/constants"
	"github.com/chris-aubin/Uniswap-Simulator/src/libraries/pool"
	"github.com/chris-aubin/Uniswap-Simulator/src/libraries/transaction"
)

type Simulation struct {
	// Strategy                *strategy.Strategy
	Pool         *pool.Pool
	Transactions []transaction.Transaction
	// StartBlock			    int
	// EndBlock					int
	// UpdateInterval		    int
}

func Make(pool *pool.Pool, transactions []transaction.Transaction) *Simulation {
	return &Simulation{
		// Strategy:                strategy,
		Pool:         pool,
		Transactions: transactions,
		// StartBlock:              startTime,
		// EndBlock:                endTime,
		// UpdateInterval:          updateInterval, // In blocks, default 1
	}
}

func (s *Simulation) Simulate() {
	// strategy := s.Strategy
	pool := s.Pool
	transactions := s.Transactions
	// startBlock := s.StartBlock
	// endBlock := s.EndBlock
	// updateInterval := s.UpdateInterval

	for _, t := range transactions {

		// if transaction.BlockNo < startBlock {
		// 	continue
		// }

		// if transaction.BlockNo > s.EndBlock {
		// 	break
		// }

		switch t.Method {
		case "Mint":
			pool.Mint(t.Owner, t.TickLower, t.TickUpper, t.Amount)
		case "Burn":
			pool.Burn(t.Owner, t.TickLower, t.TickUpper, t.Amount)
		case "Swap":
			zeroForOne := true
			amount := t.Amount1
			if t.Amount0.Cmp(big.NewInt(0)) >= 1 {
				zeroForOne = false
				amount = t.Amount0
			}
			pool.Swap(t.Sender, t.Recipient, zeroForOne, amount, constants.MaxUint160)
		}
	}
}
