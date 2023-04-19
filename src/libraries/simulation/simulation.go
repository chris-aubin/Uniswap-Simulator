// Package simulation is used to run a simulation of a Uniswap pool
//
// It contains the definition of the simulation type and the methods necessary
// to run a simulation.
package simulation

import (
	"github.com/chris-aubin/Uniswap-Simulator/src/libraries/pool"
	"github.com/chris-aubin/Uniswap-Simulator/src/libraries/strategy"
	"github.com/chris-aubin/Uniswap-Simulator/src/libraries/transaction"
)

// Simulation represents a simulation of a Uniswap pool. It contains the
// pool, the transactions to be run, and the strategy to be tested in the
// simulation.
type Simulation struct {
	Strategy     *strategy.Strategy
	Pool         *pool.Pool
	Transactions []transaction.Transaction
}

// Make returns a new simulation struct.
func Make(pool *pool.Pool, transactions []transaction.Transaction, strategy *strategy.Strategy) *Simulation {
	return &Simulation{
		Strategy:     strategy,
		Pool:         pool,
		Transactions: transactions,
	}
}

// Simulate runs the simulation.
func (s *Simulation) Simulate() {
	startBlock := s.Transactions[0].BlockNo
	prevBlock := startBlock
	for _, t := range s.Transactions {
		// Rebalance the pool if the update interval has been reached.
		if (t.BlockNo-startBlock)%s.Strategy.UpdateInterval == 0 {
			s.Strategy.Rebalance(s.Pool, s.Strategy)
		} else if t.BlockNo-prevBlock >= s.Strategy.UpdateInterval {
			s.Strategy.Rebalance(s.Pool, s.Strategy)
		}

		// Execute the transaction.
		transaction.Execute(t, s.Pool)
		prevBlock = t.BlockNo
	}
}
