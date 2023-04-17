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

func main() {
	relPathToTransactions := flag.String("transactions", "../data/transactions.txt", "Path to file containing transactions for simulation")
	relPathToPoolState := flag.String("pool", "../data/pool.txt", "Path to file containing pool state for simulation")
	relPathToGas := flag.String("gas", "../data/gas.txt", "Path to file containing gas averages for simulation")
	stratIdentifier := flag.String("strat", "v2", "Strategy identifier")
	amount0String := flag.String("amount", "100", "Strategy amount0")
	amount1String := flag.String("amount1", "100", "Strategy amount1")
	updateInterval := flag.Int("updateInterval", 1, "Update interval for simulation")
	flag.Parse()

	amount0, _ := new(big.Int).SetString(*amount0String, 10)
	amount1, _ := new(big.Int).SetString(*amount1String, 10)

	absPathToTransactions, err := filepath.Abs(*relPathToTransactions)
	if err != nil {
		message := fmt.Sprint("Failed to get absolute path to file containing transactions:", err)
		panic(message)
	}

	absPathToPoolState, err := filepath.Abs(*relPathToPoolState)
	if err != nil {
		message := fmt.Sprint("Failed to get absolute path to file containing pool state:", err)
		panic(message)
	}

	absPathToGas, err := filepath.Abs(*relPathToGas)
	if err != nil {
		message := fmt.Sprint("Failed to get absolute path to file containing gas averages:", err)
		panic(message)
	}

	transactionsRaw, err := os.ReadFile(absPathToTransactions)
	if err != nil {
		message := fmt.Sprintf("Error reading transactions file at path (relative path, absolute path): %s, %s, %v", *relPathToTransactions, absPathToTransactions, err)
		panic(message)
	}

	poolStateRaw, err := os.ReadFile(absPathToPoolState)
	if err != nil {
		message := fmt.Sprintf("Error reading pool state file at path (relative path, absolute path): %s, %s, %v", *relPathToPoolState, absPathToPoolState, err)
		panic(message)
	}

	gasRaw, err := os.ReadFile(absPathToGas)
	if err != nil {
		message := fmt.Sprintf("Error reading gas averages file at path (relative path, absolute path): %s, %s, %v", *relPathToGas, absPathToGas, err)
		panic(message)
	}

	t := getTransactions(transactionsRaw)
	p := getPoolState(poolStateRaw)
	g := getGasAvs(gasRaw)

	strat := strategy.Make(amount0, amount1, p, g, *stratIdentifier, *updateInterval)

	s := simulation.Make(p, t, strat)

	// Save pool state before simulation
	poolJSON, _ := json.MarshalIndent(s.Pool, "", "    ")
	f, _ := os.Create("poolBefore.txt")
	f.Write(poolJSON)
	f.Close()

	s.Simulate()

	// Save pool state after simulation
	poolJSON, _ = json.MarshalIndent(s.Pool, "", "    ")
	f, _ = os.Create("poolAfter.txt")
	f.Write(poolJSON)
	f.Close()

	// Save strategy after to file
	amount0, amount1, gasUsed := s.Strategy.Results(p)
	fmt.Println(amount0)
	fmt.Println(amount1)
	fmt.Println(gasUsed)
}
