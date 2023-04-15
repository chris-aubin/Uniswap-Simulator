package simulation

// package executor

import (
	"github.com/chris-aubin/Uniswap-Simulator/src/libraries/pool"
	"github.com/chris-aubin/Uniswap-Simulator/src/libraries/strategy"
	"github.com/chris-aubin/Uniswap-Simulator/src/libraries/transaction"
)

type Simulation struct {
	Strategy     *strategy.Strategy
	Pool         *pool.Pool
	Transactions []transaction.Transaction
	StartBlock   int
	EndBlock     int
}

func Make(pool *pool.Pool, transactions []transaction.Transaction, strategy *strategy.Strategy, startBlock, endBlock, updateInterval int) *Simulation {
	return &Simulation{
		Strategy:     strategy,
		Pool:         pool,
		Transactions: transactions,
		StartBlock:   startBlock,
		EndBlock:     endBlock,
	}
}

func (s *Simulation) Simulate() {
	for _, t := range s.Transactions {

		if t.BlockNo < s.StartBlock {
			continue
		}

		if t.BlockNo > s.EndBlock {
			// End simulation (must calculate returns, etc)
			break
		}

		if (t.BlockNo-s.StartBlock)%s.Strategy.UpdateInterval == 0 {
			// Call rebalance function
			s.Strategy.Rebalance(s.Pool)
		}

		transaction.Execute(t, s.Pool)
	}
}
