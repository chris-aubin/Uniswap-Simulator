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
			// Is the swap token0 for token1 or token1 for token0? The value
			// that is greater than 0 is the token that the user provided.
			// There's no way to tell whether the swap was for an exact input
			// or an exact output, so we'll just assume that all swaps are for
			// an exact input (by providing the positive amount). We also set
			// the price limit to the max value of a uint160 to ensure that all
			// swaps are executed in their entirety.
			zeroForOne := true
			amount := t.Amount0
			if t.Amount1.Cmp(big.NewInt(0)) >= 1 {
				zeroForOne = false
				amount = t.Amount1
			}
			pool.Swap(t.Sender, t.Recipient, zeroForOne, amount, constants.MaxUint160)
		}
	}
}
