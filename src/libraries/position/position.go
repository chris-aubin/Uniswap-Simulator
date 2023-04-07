// Package position contains the definition of the position type and the methods
// necessary to manipulate positions.
//
// A position struct represent an owner address' liquidity between a lower and
// upper tick boundary. They also store additional state for the tracking fees
// owed to the position.
package position

import (
	"math/big"

	"github.com/chris-aubin/Uniswap-Simulator/src/libraries/constants"
	"github.com/chris-aubin/Uniswap-Simulator/src/libraries/liquidityMath"
)

type Position struct {
	// The amount of liquidity owned by this position.
	Liquidity *big.Int
	// Fee growth per unit of liquidity as of the last update to liquidity or
	// fees owed.
	FeeGrowthInside0LastX128 *big.Int
	FeeGrowthInside1LastX128 *big.Int
	// The fees owed to the position owner in token0/token1.
	TokensOwed0 *big.Int
	TokensOwed1 *big.Int
}

// Calculates and credits accumulated fees to a user's position.
//
// Arguments:
// liquidityDelta       -- The change in pool liquidity as a result of the
//                         position update
// feeGrowthInside0X128 -- The all-time fee growth in token0, per unit of
//                         liquidity, inside the position's tick boundaries
// feeGrowthInside1X128 -- The all-time fee growth in token1, per unit of
//                         liquidity, inside the position's tick boundaries
func (p *Position) Update(liquidityDelta, feeGrowthGlobal0X128, feeGrowthGlobal1X128 *big.Int) {
	liquidityNext := new(big.Int)
	if liquidityDelta.Cmp(big.NewInt(0)) == 0 {
		liquidityNext = p.Liquidity
	} else {
		liquidityNext = liquidityMath.AddDelta(p.Liquidity, liquidityDelta)
	}

	// Calculate accumulated fees
	tokensOwed0 := new(big.Int).Div(new(big.Int).Mul(p.Liquidity, new(big.Int).Sub(feeGrowthGlobal0X128, p.FeeGrowthInside0LastX128)), constants.Q128)
	tokensOwed1 := new(big.Int).Div(new(big.Int).Mul(p.Liquidity, new(big.Int).Sub(feeGrowthGlobal1X128, p.FeeGrowthInside1LastX128)), constants.Q128)

	// Update the position
	if liquidityDelta.Cmp(big.NewInt(0)) != 0 {
		p.Liquidity = liquidityNext
	}
	p.FeeGrowthInside0LastX128 = feeGrowthGlobal0X128
	p.FeeGrowthInside1LastX128 = feeGrowthGlobal1X128
	if (tokensOwed0.Cmp(big.NewInt(0)) == 1) || (tokensOwed1.Cmp(big.NewInt(0)) == 1) {
		p.TokensOwed0 = new(big.Int).Add(p.TokensOwed0, tokensOwed0)
		p.TokensOwed1 = new(big.Int).Add(p.TokensOwed1, tokensOwed1)
	}
}

func (p *Position) Make() *Position {
	return &Position{
		Liquidity:                big.NewInt(0),
		FeeGrowthInside0LastX128: big.NewInt(0),
		FeeGrowthInside1LastX128: big.NewInt(0),
		TokensOwed0:              big.NewInt(0),
		TokensOwed1:              big.NewInt(0),
	}
}
