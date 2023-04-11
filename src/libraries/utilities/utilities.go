package utilities

import (
	"math/big"

	"github.com/chris-aubin/Uniswap-Simulator/src/libraries/constants"
)

func EncodePriceSqrt(reserve1, reserve0 *big.Int) *big.Int {
	r1 := new(big.Float).SetInt(reserve1)
	r0 := new(big.Float).SetInt(reserve0)
	price := new(big.Float).Quo(r1, r0)
	priceSqrt := new(big.Float).Sqrt(price)
	// Convert to Q96
	priceSqrtQ96 := new(big.Float).Mul(priceSqrt, new(big.Float).SetInt(constants.Q96))
	priceSqrtQ96Int, _ := priceSqrtQ96.Int(nil)
	// Check rounding
	priceSqrtQ96IntFloat := new(big.Float).SetInt(priceSqrtQ96Int)
	if priceSqrtQ96.Cmp(priceSqrtQ96IntFloat) >= 1 {
		priceSqrtQ96Int.Sub(priceSqrtQ96Int, big.NewInt(1))
	}

	return priceSqrtQ96Int
}
