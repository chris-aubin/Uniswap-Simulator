package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"github.com/chris-aubin/Uniswap-Simulator/src/libraries/pool"
	"github.com/chris-aubin/Uniswap-Simulator/src/libraries/simulation"
	"github.com/chris-aubin/Uniswap-Simulator/src/libraries/transaction"
)

func getTransactions(transactionsRaw []byte) []*transaction.Transaction {
	type getTransactionsInput struct {
		Data []*transaction.Transaction
	}
	var transactionsInput getTransactionsInput
	var transactions []*transaction.Transaction

	json.Unmarshal(transactionsRaw, &transactionsInput)
	transactions = transactionsInput.Data
	return transactions
}

func getPoolState(poolRaw []byte) *pool.Pool {
	type getPoolStateInput struct {
		Data pool.Pool
	}
	var poolInput getPoolStateInput
	var pool pool.Pool

	json.Unmarshal(poolRaw, &poolInput)
	pool = poolInput.Data
	return &pool
}

func getGasAvs(gasAvsRaw []byte) *map[string]float64 {
	type getGasAvsInput struct {
		Data map[string]float64
	}
	var gasAvsInput getGasAvsInput
	var gasAvs map[string]float64

	json.Unmarshal(gasAvsRaw, &gasAvsInput)
	gasAvs = gasAvsInput.Data
	return &gasAvs
}

func main() {
	relPathToTransactions := flag.String("transactions", "../data/transactions.txt", "Path to file containing transactions for simulation")
	relPathToPoolState := flag.String("pool", "../data/poolState.txt", "Path to file containing pool state for simulation")
	// relPathToGasAvs := flag.String("gasAvs", "../data/gasAvs.txt", "Path to file containing gas averages for simulation")
	flag.Parse()

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

	// absPathToGasAvs, err := filepath.Abs(*relPathToGasAvs)
	// if err != nil {
	// 	message := fmt.Sprint("Failed to get absolute path to file containing gas averages:", err)
	// 	panic(message)
	// }

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

	// gasAvsRaw, err := os.ReadFile(absPathToGasAvs)
	// if err != nil {
	// 	message := fmt.Sprintf("Error reading gas averages file at path (relative path, absolute path): %s, %s, %v", *relPathToGasAvs, absPathToGasAvs, err)
	// 	panic(message)
	// }

	transactions := getTransactions(transactionsRaw)
	poolState := getPoolState(poolStateRaw)
	// gasAvs := getGasAvs(gasAvsRaw)

	s := simulation.Make(poolState, transactions)
	s.Simulate()
	fmt.Println("Pool state after simulation:")
	fmt.Printf("%+v", s.Pool)
}
