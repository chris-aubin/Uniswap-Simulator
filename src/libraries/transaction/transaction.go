package transaction

import (
	"math/big"
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

// func (t Transaction) MarshalJSON() ([]byte, error) {
// 	switch t.Method {
// 	case "MINT":
// 		return json.Marshal(&Transaction{
// 			BlockNo:   t.BlockNo,
// 			Timestamp: t.Timestamp,
// 			GasPrice:  t.GasPrice,
// 			GasUsed:   t.GasUsed,
// 			GasTotal:  t.GasTotal,
// 			Method:    t.Method,
// 			MethodData: MintMethodData{
// 				Sender:    t.MethodData.(MintMethodData).Sender,
// 				Owner:     t.MethodData.(MintMethodData).Owner,
// 				TickLower: t.MethodData.(MintMethodData).TickLower,
// 				TickUpper: t.MethodData.(MintMethodData).TickUpper,
// 				Amount:    t.MethodData.(MintMethodData).Amount,
// 				Amount0:   t.MethodData.(MintMethodData).Amount0,
// 				Amount1:   t.MethodData.(MintMethodData).Amount1,
// 			},
// 		})
// 	case "BURN":
// 		return json.Marshal(&Transaction{
// 			BlockNo:   t.BlockNo,
// 			Timestamp: t.Timestamp,
// 			GasPrice:  t.GasPrice,
// 			GasUsed:   t.GasUsed,
// 			GasTotal:  t.GasTotal,
// 			Method:    t.Method,
// 			MethodData: BurnMethodData{
// 				Owner:     t.MethodData.(BurnMethodData).Owner,
// 				TickLower: t.MethodData.(BurnMethodData).TickLower,
// 				TickUpper: t.MethodData.(BurnMethodData).TickUpper,
// 				Amount:    t.MethodData.(BurnMethodData).Amount,
// 				Amount0:   t.MethodData.(BurnMethodData).Amount0,
// 				Amount1:   t.MethodData.(BurnMethodData).Amount1,
// 			},
// 		})
// 	case "SWAP":
// 		return json.Marshal(&Transaction{
// 			BlockNo:   t.BlockNo,
// 			Timestamp: t.Timestamp,
// 			GasPrice:  t.GasPrice,
// 			GasUsed:   t.GasUsed,
// 			GasTotal:  t.GasTotal,
// 			Method:    t.Method,
// 			MethodData: SwapMethodData{
// 				Sender:       t.MethodData.(SwapMethodData).Sender,
// 				Recipient:    t.MethodData.(SwapMethodData).Recipient,
// 				Amount0:      t.MethodData.(SwapMethodData).Amount0,
// 				Amount1:      t.MethodData.(SwapMethodData).Amount1,
// 				SqrtPriceX96: t.MethodData.(SwapMethodData).SqrtPriceX96,
// 				Liquidity:    t.MethodData.(SwapMethodData).Liquidity,
// 				Tick:         t.MethodData.(SwapMethodData).Tick,
// 			},
// 		})
// 	}
// 	panic("unreachable")
// }
