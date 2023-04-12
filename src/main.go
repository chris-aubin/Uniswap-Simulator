package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"

	"github.com/chris-aubin/Uniswap-Simulator/src/libraries/pool"
	"github.com/chris-aubin/Uniswap-Simulator/src/libraries/transaction"
)

func getTransactions(transactionsRaw []byte) *[]transaction.Transaction {
	type getTransactionsInput struct {
		Data []transaction.Transaction
	}
	var transactionsInput getTransactionsInput
	var transactions []transaction.Transaction

	json.Unmarshal(transactionsRaw, &transactionsInput)
	transactions = transactionsInput.Data
	return &transactions
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
	pathToTransactions := flag.String("transactions", "./data/transactions.txt", "Path to file containing transactions for simulation")
	pathToPoolState := flag.String("pool", "./data/poolState.txt", "Path to file containing pool state for simulation")
	pathToGasAvs := flag.String("gasAvs", "./data/gasAvs.txt", "Path to file containing gas averages for simulation")
	flag.Parse()

	transactionsRaw, err := os.ReadFile(*pathToTransactions)
	if err != nil {
		message := fmt.Sprintf("Error reading transactions file at path: %s", *pathToTransactions)
		panic(message)
	}

	poolStateRaw, err := os.ReadFile(*pathToPoolState)
	if err != nil {
		message := fmt.Sprintf("Error reading pool state file at path: %s", *pathToPoolState)
		panic(message)
	}

	gasAvsRaw, err := os.ReadFile(*pathToGasAvs)
	if err != nil {
		message := fmt.Sprintf("Error reading gas averages file at path: %s", *pathToGasAvs)
		panic(message)
	}

	transactions := getTransactions(transactionsRaw)
	poolState := getPoolState(poolStateRaw)
	gasAvs := getGasAvs(gasAvsRaw)

	fmt.Println(transactions)
	fmt.Println()
	fmt.Println(poolState)
	fmt.Println()
	fmt.Println(gasAvs)
	fmt.Println()
}
