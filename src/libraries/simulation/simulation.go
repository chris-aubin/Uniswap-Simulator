package simulation

import (
	"github.com/chris-aubin/Uniswap-Simulator/src/libraries/pool"
	"github.com/chris-aubin/Uniswap-Simulator/src/libraries/strategy"
	"github.com/chris-aubin/Uniswap-Simulator/src/libraries/transaction"
)

type Simulation struct {
	Strategy     *strategy.Strategy
	Pool         *pool.Pool
	Transactions []transaction.Transaction
}

// func Make(pool *pool.Pool, transactions []transaction.Transaction, strategy *strategy.Strategy, startBlock, endBlock, updateInterval int) *Simulation {
func Make(pool *pool.Pool, transactions []transaction.Transaction, strategy *strategy.Strategy) *Simulation {
	return &Simulation{
		Strategy:     strategy,
		Pool:         pool,
		Transactions: transactions,
	}
}

func (s *Simulation) Simulate() {
	startBlock := s.Transactions[0].BlockNo
	prevBlock := startBlock
	for _, t := range s.Transactions {
		if (t.BlockNo-startBlock)%s.Strategy.UpdateInterval == 0 {
			// Call rebalance function
			s.Strategy.Rebalance(s.Pool, s.Strategy)
		} else if t.BlockNo-prevBlock >= s.Strategy.UpdateInterval {
			// Call rebalance function
			s.Strategy.Rebalance(s.Pool, s.Strategy)
		}

		transaction.Execute(t, s.Pool)
		prevBlock = t.BlockNo
	}
}
