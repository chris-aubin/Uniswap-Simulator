package position

import (
	"math/big"

	"github.com/chris-aubin/Uniswap-Simulator/src/libraries/liquidityMath"
)

var (
	Q128, _ = new(big.Int).SetString("0x100000000000000000000000000000000", 16)
)

type Position struct {
	// The amount of liquidity owned by this position
	Liquidity *big.Int
	// Fee growth per unit of liquidity as of the last update to liquidity or fees owed
	FeeGrowthInside0LastX128 *big.Int
	FeeGrowthInside1LastX128 *big.Int
	// The fees owed to the position owner in token0/token1
	TokensOwed0 *big.Int
	TokensOwed1 *big.Int
}

func (p *Position) Update(
	liquidityDelta *big.Int,
	feeGrowthGlobal0X128 *big.Int,
	feeGrowthGlobal1X128 *big.Int,
) {
	liquidityNext := new(big.Int)
	if liquidityDelta.Cmp(big.NewInt(0)) == 0 {
		if p.Liquidity.Cmp(big.NewInt(0)) < 1 {
			panic("position.Update: NP") // disallow pokes for 0 liquidity positions
		}
		liquidityNext = p.Liquidity
	} else {
		liquidityNext = liquidityMath.AddDelta(p.Liquidity, liquidityDelta)
	}

	// Calculate accumulated fees
	tokensOwed0 := new(big.Int).Div(new(big.Int).Mul(p.Liquidity, new(big.Int).Sub(feeGrowthGlobal0X128, p.FeeGrowthInside0LastX128)), Q128)
	tokensOwed1 := new(big.Int).Div(new(big.Int).Mul(p.Liquidity, new(big.Int).Sub(feeGrowthGlobal1X128, p.FeeGrowthInside1LastX128)), Q128)

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
