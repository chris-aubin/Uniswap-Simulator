package transaction

import (
	"encoding/json"
)

type Transaction struct {
	BlockNo    string
	Timestamp  int
	GasPrice   int
	GasUsed    int
	GasTotal   int
	Method     string
	MethodData interface{}
}

type MintMethodData struct {
	Sender    string
	Owner     string
	TickLower int
	TickUpper int
	Amount    string
	Amount0   string
	Amount1   string
}

type BurnMethodData struct {
	Owner     string
	TickLower int
	TickUpper int
	Amount    string
	Amount0   string
	Amount1   string
}

type SwapMethodData struct {
	Sender       string
	Recipient    string
	Amount0      string
	Amount1      string
	SqrtPriceX96 string
	Liquidity    string
	Tick         int
}

func (t Transaction) MarshalJSON() ([]byte, error) {
	switch t.Method {
	case "MINT":
		return json.Marshal(&Transaction{
			BlockNo:   t.BlockNo,
			Timestamp: t.Timestamp,
			GasPrice:  t.GasPrice,
			GasUsed:   t.GasUsed,
			GasTotal:  t.GasTotal,
			Method:    t.Method,
			MethodData: MintMethodData{
				Sender:    t.MethodData.(MintMethodData).Sender,
				Owner:     t.MethodData.(MintMethodData).Owner,
				TickLower: t.MethodData.(MintMethodData).TickLower,
				TickUpper: t.MethodData.(MintMethodData).TickUpper,
				Amount:    t.MethodData.(MintMethodData).Amount,
				Amount0:   t.MethodData.(MintMethodData).Amount0,
				Amount1:   t.MethodData.(MintMethodData).Amount1,
			},
		})
	case "BURN":
		return json.Marshal(&Transaction{
			BlockNo:   t.BlockNo,
			Timestamp: t.Timestamp,
			GasPrice:  t.GasPrice,
			GasUsed:   t.GasUsed,
			GasTotal:  t.GasTotal,
			Method:    t.Method,
			MethodData: BurnMethodData{
				Owner:     t.MethodData.(BurnMethodData).Owner,
				TickLower: t.MethodData.(BurnMethodData).TickLower,
				TickUpper: t.MethodData.(BurnMethodData).TickUpper,
				Amount:    t.MethodData.(BurnMethodData).Amount,
				Amount0:   t.MethodData.(BurnMethodData).Amount0,
				Amount1:   t.MethodData.(BurnMethodData).Amount1,
			},
		})
	case "SWAP":
		return json.Marshal(&Transaction{
			BlockNo:   t.BlockNo,
			Timestamp: t.Timestamp,
			GasPrice:  t.GasPrice,
			GasUsed:   t.GasUsed,
			GasTotal:  t.GasTotal,
			Method:    t.Method,
			MethodData: SwapMethodData{
				Sender:       t.MethodData.(SwapMethodData).Sender,
				Recipient:    t.MethodData.(SwapMethodData).Recipient,
				Amount0:      t.MethodData.(SwapMethodData).Amount0,
				Amount1:      t.MethodData.(SwapMethodData).Amount1,
				SqrtPriceX96: t.MethodData.(SwapMethodData).SqrtPriceX96,
				Liquidity:    t.MethodData.(SwapMethodData).Liquidity,
				Tick:         t.MethodData.(SwapMethodData).Tick,
			},
		})
	}
	panic("unreachable")
}
