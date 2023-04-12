package pool

import (
	"encoding/json"
	"math/big"
	"strconv"

	"github.com/chris-aubin/Uniswap-Simulator/src/libraries/position"
	"github.com/chris-aubin/Uniswap-Simulator/src/libraries/tick"
)

type ProtocolFeesInput struct {
	Token0 *big.Int `json:"token0"`
	Token1 *big.Int `json:"token1"`
}

type Slot0Input struct {
	SqrtPriceX96 *big.Int `json:"sqrtPriceX96"`
	Tick         int      `json:"tick"`
	// The simulator doesn't currently simulate the oracle, but these fields are
	// included for completeness, because they may be useful in the future if
	// anyone decides to add oracle simulation.
	ObservationIndex       int `json:"observationIndex"`
	ObservationCardinality int `json:"observationCardinality"`
	FeeProtocol            int `json:"feeProtocol"`
}

type ObservationInput struct {
	BlockTimestamp                    *big.Int `json:"blockTimestamp"`
	TickCumulative                    *big.Int `json:"tickCumulative"`
	SecondsPerLiquidityCumulativeX128 *big.Int `json:"secondsPerLiquidityCumulativeX128"`
	Initialized                       bool     `json:"initialized"`
}

type TickInput struct {
	LiquidityGross        *big.Int `json:"liquidityGross"`
	LiquidityNet          *big.Int `json:"liquidityNet"`
	FeeGrowthOutside0X128 *big.Int `json:"feeGrowthOutside0X128"`
	FeeGrowthOutside1X128 *big.Int `json:"feeGrowthOutside1X128"`
	// The simulator doesn't currently simulate the oracle, but these fields are
	// included for completeness, because they may be useful in the future if
	// anyone decides to add oracle simulation.
	TickCumulativeOutside          *big.Int `json:"tickCumulativeOutside"`
	SecondsPerLiquidityOutsideX128 *big.Int `json:"secondsPerLiquidityOutsideX128"`
	SecondsOutside                 *big.Int `json:"secondsOutside"`
	Initialized                    bool     `json:"initialized"`
}

type PositionInput struct {
	Liquidity                *big.Int `json:"liquidity"`
	FeeGrowthInside0LastX128 *big.Int `json:"feeGrowthInside0LastX128"`
	FeeGrowthInside1LastX128 *big.Int `json:"feeGrowthInside1LastX128"`
	TokensOwed0              *big.Int `json:"tokensOwed0"`
	TokensOwed1              *big.Int `json:"tokensOwed1"`
}

type PoolInput struct {
	TickSpacing          int                       `json:"tickSpacing"`
	FeeGrowthGlobal0X128 *big.Int                  `json:"feeGrowthGlobal0X128"`
	FeeGrowthGlobal1X128 *big.Int                  `json:"feeGrowthGlobal1X128"`
	ProtocolFees         *ProtocolFeesInput        `json:"protocolFees"`
	Liquidity            *big.Int                  `json:"liquidity"`
	Slot0                *Slot0Input               `json:"amount1"`
	Observations         []*ObservationInput       `json:"observations"`
	Ticks                map[string]*TickInput     `json:"ticks"`
	Positions            map[string]*PositionInput `json:"positions"`
	Balance0             *big.Int                  `json:"balance0"`
	Balance1             *big.Int                  `json:"balance1"`
}

func Slot0InputToSlot0(si *Slot0Input) *Slot0 {
	return &Slot0{
		SqrtPriceX96: si.SqrtPriceX96,
		Tick:         si.Tick,
		FeeProtocol:  si.FeeProtocol,
	}
}

func ProtocolFeesInputToProtocolFees(pfi *ProtocolFeesInput) *ProtocolFees {
	return &ProtocolFees{
		Token0: pfi.Token0,
		Token1: pfi.Token1,
	}
}

func TicksInputToTicks(ti map[string]*TickInput) *tick.Ticks {
	ticks := make(map[int]*tick.Tick)
	for k, v := range ti {
		tick_idx, _ := strconv.ParseInt(k, 10, 64)
		ticks[int(tick_idx)] = &tick.Tick{
			LiquidityGross:        v.LiquidityGross,
			LiquidityNet:          v.LiquidityNet,
			FeeGrowthOutside0X128: v.FeeGrowthOutside0X128,
			FeeGrowthOutside1X128: v.FeeGrowthOutside1X128,
			Initialized:           v.Initialized,
		}
	}
	return &tick.Ticks{TickData: ticks}
}

func PositionsInputToPositions(pi map[string]*PositionInput) map[string]*position.Position {
	positions := make(map[string]*position.Position)
	for k, v := range pi {
		positions[k] = &position.Position{
			Liquidity:                v.Liquidity,
			FeeGrowthInside0LastX128: v.FeeGrowthInside0LastX128,
			FeeGrowthInside1LastX128: v.FeeGrowthInside1LastX128,
			TokensOwed0:              v.TokensOwed0,
		}
	}
	return positions
}

func PoolInputToPool(pi *PoolInput) *Pool {
	return &Pool{
		TickSpacing:          pi.TickSpacing,
		Slot0:                Slot0InputToSlot0(pi.Slot0),
		FeeGrowthGlobal0X128: pi.FeeGrowthGlobal0X128,
		FeeGrowthGlobal1X128: pi.FeeGrowthGlobal1X128,
		ProtocolFees:         ProtocolFeesInputToProtocolFees(pi.ProtocolFees),
		Liquidity:            pi.Liquidity,
		Ticks:                TicksInputToTicks(pi.Ticks),
		Positions:            PositionsInputToPositions(pi.Positions),
	}
}

func (p PoolInput) MarshalJSON() ([]byte, error) {
	return json.Marshal(&PoolInput{
		TickSpacing:          p.TickSpacing,
		FeeGrowthGlobal0X128: p.FeeGrowthGlobal0X128,
		FeeGrowthGlobal1X128: p.FeeGrowthGlobal1X128,
		ProtocolFees:         p.ProtocolFees,
		Liquidity:            p.Liquidity,
		Slot0:                p.Slot0,
		Ticks:                p.Ticks,
		Positions:            p.Positions,
	})
}
