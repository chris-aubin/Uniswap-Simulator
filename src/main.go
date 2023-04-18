package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"math/big"
	"os"
	"path/filepath"

	"github.com/chris-aubin/Uniswap-Simulator/src/libraries/pool"
	"github.com/chris-aubin/Uniswap-Simulator/src/libraries/simulation"
	"github.com/chris-aubin/Uniswap-Simulator/src/libraries/strategy"
	"github.com/chris-aubin/Uniswap-Simulator/src/libraries/transaction"
)

func getTransactions(transactionsRaw []byte) []transaction.Transaction {
	type getTransactionsInput struct {
		Data []transaction.Transaction
	}
	var transactionsInput getTransactionsInput

	json.Unmarshal(transactionsRaw, &transactionsInput)
	transactions := transactionsInput.Data
	return transactions
}

func getPoolState(poolRaw []byte) *pool.Pool {
	type getPoolStateInput struct {
		Data pool.PoolTemp
	}
	var poolInput getPoolStateInput
	var poolTemp pool.PoolTemp
	var p *pool.Pool

	json.Unmarshal(poolRaw, &poolInput)
	poolTemp = poolInput.Data
	p = pool.PoolTempToPool(&poolTemp)
	return p
}

func getGasAvs(gasAvsRaw []byte) *strategy.GasAvs {
	type getGasAvsInput struct {
		Data strategy.GasAvs
	}
	var gasAvsInput getGasAvsInput
	var gasAvs strategy.GasAvs

	json.Unmarshal(gasAvsRaw, &gasAvsInput)
	gasAvs = gasAvsInput.Data

	if gasAvs.MintGas.Cmp(big.NewInt(-1)) == 0 {
		gasAvs.MintGas = big.NewInt(35000)
	}
	if gasAvs.BurnGas.Cmp(big.NewInt(-1)) == 0 {
		gasAvs.BurnGas = big.NewInt(20000)
	}
	if gasAvs.SwapGas.Cmp(big.NewInt(-1)) == 0 {
		gasAvs.SwapGas = big.NewInt(20000)
	}
	if gasAvs.FlashGas.Cmp(big.NewInt(-1)) == 0 {
		gasAvs.FlashGas = big.NewInt(20000)
	}
	if gasAvs.CollectGas.Cmp(big.NewInt(-1)) == 0 {
		gasAvs.CollectGas = big.NewInt(20000)
	}

	return &gasAvs
}

func getStratInput(stratRaw []byte) *strategy.StrategyInput {
	var stratInput strategy.StrategyInput

	json.Unmarshal(stratRaw, &stratInput)

	return &stratInput
}

func main() {
	// Get command line arguments
	relPathToData := flag.String("data", "../data/testNil1", "Path to file containing data for simulation")
	flag.Parse()

	// Relative paths to files containing data for simulation
	relPathToTransactions := *relPathToData + "/transactions.txt"
	relPathToPoolState := *relPathToData + "/poolBefore.txt"
	relPathToGas := *relPathToData + "/gas.txt"
	relPathToStrat := *relPathToData + "/strategy.txt"

	// Relative paths to files containing results of simulation
	relPathToResults := "../results"
	relPathToStratBefore := relPathToResults + "/strategyBefore.txt"
	relPathToStratAfter := relPathToResults + "/strategyAfter.txt"
	relPathToPoolStateBefore := relPathToResults + "/poolBefore.txt"
	relPathToPoolStateAfter := relPathToResults + "/poolAfter.txt"

	// Get absolute paths to files containing data for simulation
	absPathToTransactions, err := filepath.Abs(relPathToTransactions)
	if err != nil {
		message := fmt.Sprint("Failed to get absolute path to file containing transactions:", err)
		panic(message)
	}

	absPathToPoolState, err := filepath.Abs(relPathToPoolState)
	if err != nil {
		message := fmt.Sprint("Failed to get absolute path to file containing pool state:", err)
		panic(message)
	}

	absPathToGas, err := filepath.Abs(relPathToGas)
	if err != nil {
		message := fmt.Sprint("Failed to get absolute path to file containing gas averages:", err)
		panic(message)
	}

	absPathToStrat, err := filepath.Abs(relPathToStrat)
	if err != nil {
		message := fmt.Sprint("Failed to get absolute path to file containing strategy information:", err)
		panic(message)
	}

	// Get absolute paths to files containing results of simulation
	absPathToStratBefore, _ := filepath.Abs(relPathToStratBefore)
	absPathToStratAfter, _ := filepath.Abs(relPathToStratAfter)
	absPathToPoolStateBefore, _ := filepath.Abs(relPathToPoolStateBefore)
	absPathToPoolStateAfter, _ := filepath.Abs(relPathToPoolStateAfter)

	// Read data for simulation from files
	transactionsRaw, err := os.ReadFile(absPathToTransactions)
	if err != nil {
		message := fmt.Sprintf("Error reading transactions file at path (relative path, absolute path): %s, %s, %v", relPathToTransactions, absPathToTransactions, err)
		panic(message)
	}

	poolStateRaw, err := os.ReadFile(absPathToPoolState)
	if err != nil {
		message := fmt.Sprintf("Error reading pool state file at path (relative path, absolute path): %s, %s, %v", relPathToPoolState, absPathToPoolState, err)
		panic(message)
	}

	gasRaw, err := os.ReadFile(absPathToGas)
	if err != nil {
		message := fmt.Sprintf("Error reading gas averages file at path (relative path, absolute path): %s, %s, %v", relPathToGas, absPathToGas, err)
		panic(message)
	}

	stratRaw, err := os.ReadFile(absPathToStrat)
	if err != nil {
		message := fmt.Sprintf("Error reading strategy file at path (relative path, absolute path): %s, %s, %v", relPathToStrat, absPathToStrat, err)
		panic(message)
	}
	fmt.Println("stratRaw", string(stratRaw))

	// Create simulation
	t := getTransactions(transactionsRaw)
	p := getPoolState(poolStateRaw)
	g := getGasAvs(gasRaw)
	stratInput := getStratInput(stratRaw)

	fmt.Println("stratInput:", stratInput)
	fmt.Println("Amount0:", stratInput.Amount0)
	fmt.Println("Amount1:", stratInput.Amount1)
	fmt.Println("Strategy:", stratInput.Strategy)
	fmt.Println("UpdateInterval:", stratInput.UpdateInterval)

	strat := strategy.Make(stratInput.Amount0, stratInput.Amount1, p, g, stratInput.Strategy, stratInput.UpdateInterval)

	s := simulation.Make(p, t, strat)

	// Save pool state before simulation
	poolJSON, _ := json.MarshalIndent(s.Pool, "", "    ")
	f, _ := os.Create(absPathToPoolStateBefore)
	f.Write(poolJSON)
	f.Close()

	// Save strategy before to simulation
	stratJSON, _ := json.MarshalIndent(s.Strategy, "", "    ")
	f, _ = os.Create(absPathToStratBefore)
	f.Write(stratJSON)
	f.Close()

	s.Simulate()

	// Save pool state after simulation
	poolJSON, _ = json.MarshalIndent(s.Pool, "", "    ")
	f, _ = os.Create(absPathToPoolStateAfter)
	f.Write(poolJSON)
	f.Close()

	// Save strategy after simulation
	stratJSON, _ = json.MarshalIndent(s.Strategy, "", "    ")
	f, _ = os.Create(absPathToStratAfter)
	f.Write(stratJSON)
	f.Close()

	// amount0, amount1, gasUsed := s.Strategy.Results(p)
	// fmt.Println(amount0)
	// fmt.Println(amount1)
	// fmt.Println(gasUsed)
}
