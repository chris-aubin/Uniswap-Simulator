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
	// pool := s.Pool
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
		case "MINT":
			if t.Amount.Cmp(big.NewInt(0)) == 0 {
				continue
			}
			s.Pool.Mint(t.Owner, t.TickLower, t.TickUpper, t.Amount)
		case "BURN":
			if t.Amount.Cmp(big.NewInt(0)) == 0 {
				continue
			}
			s.Pool.Burn(t.Owner, t.TickLower, t.TickUpper, t.Amount)
		case "SWAP":
			// Is the swap token0 for token1 or token1 for token0? The value
			// that is greater than 0 is the token that the user provided.
			// There's no way to tell whether the swap was for an exact input
			// or an exact output, so we'll just assume that all swaps are for
			// an exact input (by providing the positive amount). We also set
			// the price limit to the max value of a uint160 to ensure that all
			// swaps are executed in their entirety.
			zeroForOne := false
			amount := t.Amount1
			if t.SqrtPriceX96.Cmp(s.Pool.Slot0.SqrtPriceX96) <= -1 {
				zeroForOne = true
				amount = t.Amount0
			}
			s.Pool.Swap(t.Sender, t.Recipient, zeroForOne, amount, new(big.Int).Sub(constants.MaxSqrtRatio, big.NewInt(1)))
		case "FLASH":
			s.Pool.Flash(t.Paid0, t.Paid1)
		}
	}
}
