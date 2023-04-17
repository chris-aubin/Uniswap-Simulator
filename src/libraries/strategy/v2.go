package strategy

import (
	"fmt"
	"math/big"

	"github.com/chris-aubin/Uniswap-Simulator/src/libraries/constants"
	"github.com/chris-aubin/Uniswap-Simulator/src/libraries/pool"
	"github.com/chris-aubin/Uniswap-Simulator/src/libraries/tickMath"
	"github.com/chris-aubin/Uniswap-Simulator/src/libraries/utilities"
)

func V2StrategyRebalance(p *pool.Pool, s *Strategy) {
	fmt.Println("V2StrategyRebalance")
	if len(s.Positions) == 0 {
		fmt.Println("ABOUT TO MINT")
		V2StrategyMintPosition(p, s)
	}
	return
}

func V2StrategyMintPosition(p *pool.Pool, s *Strategy) {
	s.GasUsed = new(big.Int).Add(s.GasUsed, s.GasAvs.MintGas)
	tickLower := constants.MinTick
	tickUpper := constants.MaxTick
	sqrtRatioAX96 := tickMath.GetSqrtRatioAtTick(tickLower)
	sqrtRatioBX96 := tickMath.GetSqrtRatioAtTick(tickUpper)

	amount := utilities.GetLiquidityForAmount(p.Slot0.SqrtPriceX96, sqrtRatioAX96, sqrtRatioBX96, s.Amount0, s.Amount1)
	fmt.Println("V2StrategyMintPosition")
	fmt.Println("amount", amount)
	if amount.Cmp(big.NewInt(0)) <= 0 {
		return
	}
	s.Positions = append(s.Positions, &StrategyPosition{
		TickLower: tickLower,
		TickUpper: tickUpper,
		Amount:    amount,
	})

}
