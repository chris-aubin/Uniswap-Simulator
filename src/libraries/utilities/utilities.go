package utilities

import (
	"math/big"

	"github.com/chris-aubin/Uniswap-Simulator/src/libraries/constants"
)

func EncodePriceSqrt(reserve1, reserve0 *big.Int) *big.Int {

	price := new(big.Int).Div(reserve1, reserve0)
	priceSqrt := new(big.Int).Sqrt(price)

	// .Sqrt()
	return new(big.Int).Mul(priceSqrt, constants.Q96)
}
