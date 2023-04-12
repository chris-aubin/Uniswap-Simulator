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
	Pool                    *pool.Pool
	Transactions            []*transaction.Transaction
	StartBlock			    int
	EndBlock				int
	UpdateInterval		    int	
}

func Make() *Simulation {
	return &Execution{
		// Strategy:                strategy,
		Pool:                    *pool.Pool,
		Transactions:            []*transaction.Transaction,
		StartBlock:              startTime,
		EndBlock:                endTime,
		UpdateInterval:          updateInterval, // In blocks, default 1
	}
}

func (s *Simulation) Simulate() {
	// strategy := s.Strategy
	pool := s.Pool
	transactions := s.Transactions
	startBlock := s.StartBlock
	endBlock := s.EndBlock
	updateInterval := s.UpdateInterval

	for _, transaction := range transactions {

		if transaction.BlockNo < startBlock {
			continue
		}

		if transaction.BlockNo > s.EndBlock {
			break
		}

		switch transaction.Method {
		case "MINT":
			pool.Mint(transaction.MethodData.Owner, transaction.MethodData.TickLower, transaction.MethodData.TickUpper, transaction.MethodData.Amount)
		case "BURN":
			pool.Burn(transaction.MethodData.Owner, transaction.TickLower, transaction.TickUpper, transaction.Amount)
		case "SWAP":
			zeroForOne := true
			amount := transaction.MethodData.Amount1
			if transaction.MethodData.Amount0.Cmp(big.NewInt(0)) >= 1 {
				zeroForOne = false
				amount = transaction.MethodData.Amount0
			}
			pool.Swap(transaction.MethodData.Sender, transaction.MethodData.Recipient, zeroForOne, amount, constants.MaxUint160)
		}
	}
}
