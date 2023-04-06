package transaction

import (
	"encoding/json"
)

type MintMethodDataInput struct {
	Sender    string `json:"sender"`
	Owner     string `json:"owner"`
	TickLower int    `json:"tickLower"`
	TickUpper int    `json:"tickUpper"`
	Amount    string `json:"amount"`
	Amount0   string `json:"amount0"`
	Amount1   string `json:"amount1"`
}

type BurnMethodDataInput struct {
	Owner     string `json:"owner"`
	TickLower int    `json:"tickLower"`
	TickUpper int    `json:"tickUpper"`
	Amount    string `json:"amount"`
	Amount0   string `json:"amount0"`
	Amount1   string `json:"amount1"`
}

type SwapMethodDataInput struct {
	Sender       string `json:"sender"`
	Recipient    string `json:"recipient"`
	Amount0      string `json:"amount0"`
	Amount1      string `json:"amount1"`
	SqrtPriceX96 string `json:"sqrtPriceX96"`
	Liquidity    string `json:"liquidity"`
	Tick         int    `json:"tick"`
}

type TransactionInput struct {
	BlockNo    string      `json:"blockNo"`
	Timestamp  int         `json:"timestamp"`
	GasPrice   int         `json:"gasPrice"`
	GasUsed    int         `json:"gasUsed"`
	GasTotal   int         `json:"gasTotal"`
	Method     string      `json:"method"`
	MethodData interface{} `json:"methodData"`
}

type Transaction struct {
	BlockNo    string
	Timestamp  int
	GasPrice   int
	GasUsed    int
	GasTotal   int
	Method     string
	MethodData interface{}
}

func (t Transaction) MarshalJSON() ([]byte, error) {
	switch t.Method {
	case "MINT":
		return json.Marshal(&TransactionInput{
			BlockNo:   t.BlockNo,
			Timestamp: t.Timestamp,
			GasPrice:  t.GasPrice,
			GasUsed:   t.GasUsed,
			GasTotal:  t.GasTotal,
			Method:    t.Method,
			MethodData: MintMethodDataInput{
				Sender:    t.MethodData.(MintMethodDataInput).Sender,
				Owner:     t.MethodData.(MintMethodDataInput).Owner,
				TickLower: t.MethodData.(MintMethodDataInput).TickLower,
				TickUpper: t.MethodData.(MintMethodDataInput).TickUpper,
				Amount:    t.MethodData.(MintMethodDataInput).Amount,
				Amount0:   t.MethodData.(MintMethodDataInput).Amount0,
				Amount1:   t.MethodData.(MintMethodDataInput).Amount1,
			},
		})
	case "BURN":
		return json.Marshal(&TransactionInput{
			BlockNo:   t.BlockNo,
			Timestamp: t.Timestamp,
			GasPrice:  t.GasPrice,
			GasUsed:   t.GasUsed,
			GasTotal:  t.GasTotal,
			Method:    t.Method,
			MethodData: BurnMethodDataInput{
				Owner:     t.MethodData.(BurnMethodDataInput).Owner,
				TickLower: t.MethodData.(BurnMethodDataInput).TickLower,
				TickUpper: t.MethodData.(BurnMethodDataInput).TickUpper,
				Amount:    t.MethodData.(BurnMethodDataInput).Amount,
				Amount0:   t.MethodData.(BurnMethodDataInput).Amount0,
				Amount1:   t.MethodData.(BurnMethodDataInput).Amount1,
			},
		})
	case "SWAP":
		return json.Marshal(&TransactionInput{
			BlockNo:   t.BlockNo,
			Timestamp: t.Timestamp,
			GasPrice:  t.GasPrice,
			GasUsed:   t.GasUsed,
			GasTotal:  t.GasTotal,
			Method:    t.Method,
			MethodData: SwapMethodDataInput{
				Sender:       t.MethodData.(SwapMethodDataInput).Sender,
				Recipient:    t.MethodData.(SwapMethodDataInput).Recipient,
				Amount0:      t.MethodData.(SwapMethodDataInput).Amount0,
				Amount1:      t.MethodData.(SwapMethodDataInput).Amount1,
				SqrtPriceX96: t.MethodData.(SwapMethodDataInput).SqrtPriceX96,
				Liquidity:    t.MethodData.(SwapMethodDataInput).Liquidity,
				Tick:         t.MethodData.(SwapMethodDataInput).Tick,
			},
		})
	}
	panic("unreachable")
}
