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
	Transactions []*transaction.Transaction
	// StartBlock			    int
	// EndBlock					int
	// UpdateInterval		    int
}

func Make(pool *pool.Pool, transactions []*transaction.Transaction) *Simulation {
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
		case "MINT":
			methodData := t.MethodData
			mintMethodData := methodData.(transaction.MintMethodData)
			pool.Mint(mintMethodData.Owner, mintMethodData.TickLower, mintMethodData.TickUpper, mintMethodData.Amount)
		case "BURN":
			methodData := t.MethodData
			burnMethodData := methodData.(transaction.BurnMethodData)
			pool.Burn(burnMethodData.Owner, burnMethodData.TickLower, burnMethodData.TickUpper, burnMethodData.Amount)
		case "SWAP":
			zeroForOne := true
			methodData := t.MethodData
			swapMethodData := methodData.(transaction.SwapMethodData)
			amount := swapMethodData.Amount1
			if swapMethodData.Amount0.Cmp(big.NewInt(0)) >= 1 {
				zeroForOne = false
				amount = swapMethodData.Amount0
			}
			pool.Swap(swapMethodData.Sender, swapMethodData.Recipient, zeroForOne, amount, constants.MaxUint160)
		}
	}
}
