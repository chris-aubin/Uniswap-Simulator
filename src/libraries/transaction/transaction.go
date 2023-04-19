package transaction

import (
	"fmt"
	"math/big"

	"github.com/chris-aubin/Uniswap-Simulator/src/libraries/constants"
	"github.com/chris-aubin/Uniswap-Simulator/src/libraries/pool"
)

type Transaction struct {
	BlockNo      int      `json:"blockNo"`
	Timestamp    int      `json:"timestamp"`
	GasPrice     int      `json:"gasPrice"`
	GasUsed      int      `json:"gasUsed"`
	GasTotal     int      `json:"gasTotal"`
	Method       string   `json:"method"`
	Sender       string   `json:"sender"`
	Recipient    string   `json:"recipient"`
	Owner        string   `json:"owner"`
	TickLower    int      `json:"tickLower"`
	TickUpper    int      `json:"tickUpper"`
	Amount       *big.Int `json:"amount"`
	Amount0      *big.Int `json:"amount0"`
	Amount1      *big.Int `json:"amount1"`
	SqrtPriceX96 *big.Int `json:"sqrtPriceX96"`
	Liquidity    *big.Int `json:"liquidity"`
	Tick         int      `json:"tick"`
	Paid0        *big.Int `json:"paid0"`
	Paid1        *big.Int `json:"paid1"`
}

func Execute(t Transaction, p *pool.Pool) {
	fmt.Println()
	fmt.Println()
	fmt.Println("Pool liquidity: ", p.Liquidity)
	fmt.Println("Pool current tick: ", p.Slot0.Tick)
	fmt.Println("Pool current sqrt price: ", p.Slot0.SqrtPriceX96)
	fmt.Printf("Transaction: %+v", t)
	fmt.Println()
	switch t.Method {
	case "MINT":
		if t.Amount.Cmp(big.NewInt(0)) == 0 {
			return
		}
		p.Mint(t.Owner, t.TickLower, t.TickUpper, t.Amount)
	case "BURN":
		if t.Amount.Cmp(big.NewInt(0)) == 0 {
			return
		}
		p.Burn(t.Owner, t.TickLower, t.TickUpper, t.Amount)
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
		if t.Amount0.Cmp(big.NewInt(0)) >= 1 {
			zeroForOne = true
			amount = t.Amount0
		}
		p.Swap(t.Sender, t.Recipient, zeroForOne, amount, new(big.Int).Sub(constants.MaxSqrtRatio, big.NewInt(1)))
	case "FLASH":
		p.Flash(t.Paid0, t.Paid1)
	}
}
