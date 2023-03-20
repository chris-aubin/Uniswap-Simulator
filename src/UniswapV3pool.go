package pool

import (
	"math/big"
	"time"

	"github.com/chris-aubin/Uniswap-Simulator/src/libraries/position"
	"github.com/chris-aubin/Uniswap-Simulator/src/libraries/tick"
)

type Slot0 struct {
	// The current price
	SqrtPriceX96 *big.Int
	// The current tick
	Tick int
	// The most-recently updated index of the observations array
	ObservationIndex int
	// The current maximum number of observations that are being stored
	ObservationCardinality int
	// The next maximum number of observations to store, triggered in observations.write
	ObservationCardinalityNext int
	// The current protocol fee as a percentage of the swap fee taken on withdrawal
	// Represented as an integer denominator (1/x)%
	FeeProtocol int
}

// Accumulated protocol fees in token0/token1 units (fees that could be collected by Uniswap governance)
type ProtocolFees struct {
	Token0 *big.Int
	Token1 *big.Int
}

// Pool state
type Pool struct {
	Slot0                Slot0
	FeeGrowthGlobal0X128 *big.Int
	FeeGrowthGlobal1X128 *big.Int
	ProtocolFees         ProtocolFees
	Liquidity            *big.Int
	// Tick-indexed state, see section 6.3 in Uniswap V3 Whitepaper
	Ticks				 t.Ticks
	// Keeps track of which ticks have been initialised, see section 6.2 in Uniswap V3 Whitepaper
	TickBitmap 	         map[int]bool
	// Position-indexed state, see section 6.4 in Uniswap V3 Whitepaper
	Positions			 map[int]*position.Position
	// Balance of token0 and token1 held by the pool. Not part of state in the deployed contract
	// (the deployed contract checks the balance of the token owned by the pool address)
	Balance0 	         *big.Int
	Balance1 	         *big.Int
}

// Common checks for valid tick inputs
func checkTicks(tickLower int, tickUpper int) {
	if tickLower >= tickUpper {
		panic("Pool.checkTicks: TLU")
	}
	if tickLower < tickMath.MIN_TICK {
		panic("Pool.checkTicks: TLM")
	}
	if tickUpper > tickMath.MAX_TICK {
		panic("Pool.checkTicks: TUM")
	}
}

func _blockTimestamp () int {
	return int(time.Now().Unix())
}
// Returns a snapshot of the tick cumulative, seconds per liquidity and seconds inside a tick range
func (p *Pool) snapshotCumulativesInside(tickLower int, tickUpper int) (big.Int, big.Int, int) {
	checkTicks(tickLower, tickUpper)

	tickCumulativeLower := new(big.Int)
	tickCumulativeUpper := new(big.Int)
	secondsPerLiquidityOutsideLowerX128 := new(big.Int)
	secondsPerLiquidityOutsideUpperX128 := new(big.Int)
	var secondsOutsideLower int
	var secondsOutsideUpper int

	lower := p.Ticks.TickData[tickLower]
	upper := p.Ticks.TickData[tickUpper]

	var initialisedLower bool

	tickCumulativeLower = lower.TickCumulativeOutside
	secondsPerLiquidityOutsideLowerX128 = lower.SecondsPerLiquidityOutsideX128
	secondsOutsideLower = lower.SecondsOutside
	initialisedLower = lower.Initialised
	if !initialisedLower {
		panic()
	}

	tickCumulativeUpper = upper.TickCumulativeOutside
	secondsPerLiquidityOutsideUpperX128 = upper.SecondsPerLiquidityOutsideX128
	secondsOutsideUpper = upper.SecondsOutside
	initialisedUpper := upper.Initialised
	if !initialisedUpper {
		panic()
	}

	_slot0 := p.Slot0

	///////////////////////////
	if (_slot0.tick < tickLower) {
		return (
			new(big.Int).Sub(tickCumulativeLower, tickCumulativeUpper),
			new(big.Int).Sub(secondsPerLiquidityOutsideLowerX128, secondsPerLiquidityOutsideUpperX128),
			secondsOutsideLower - secondsOutsideUpper
		);
	} else if (_slot0.tick < tickUpper) {
		uint32 time = _blockTimestamp();
		(int56 tickCumulative, uint160 secondsPerLiquidityCumulativeX128) = observations.observeSingle(
			time,
			0,
			_slot0.Tick,
			_slot0.ObservationIndex,
			liquidity,
			_slot0.ObservationCardinality
		);
		return (
			tickCumulative - tickCumulativeLower - tickCumulativeUpper,
			secondsPerLiquidityCumulativeX128 -
				secondsPerLiquidityOutsideLowerX128 -
				secondsPerLiquidityOutsideUpperX128,
			time - secondsOutsideLower - secondsOutsideUpper
		);
	} else {
		return (
			tickCumulativeUpper - tickCumulativeLower,
			secondsPerLiquidityOutsideUpperX128 - secondsPerLiquidityOutsideLowerX128,
			secondsOutsideUpper - secondsOutsideLower
		);
	}

}

type modifyPositionParams struct {
	// the address that owns the position
	owner int
	// the lower and upper tick of the position
	tickLower int
	tickUpper int
	// any change in liquidity
	liquidityDelta *big.Int
}

// Effect some changes to a position
// Accepts params,  an instance of the modifyPositionParams type that contains the position details
// and the change to the position's liquidity to effect
// position is a representation of the position with the given owner and tick range
// amount0 is the amount of token0 owed to the pool (it is n/ amount1 is the amount of token1 owed to the pool (it is negative if the pool should pay the recipient)
func (p *Pool) modifyPosition(params modifyPositionParams) {
	checkTicks(params.tickLower, params.tickUpper)

}

// TODO
//
func (p *Pool) Mint(recipient int, tickLower int, tickUpper int, amount int) (int, int) {

}

// TODO
//
func (p *Pool) Burn() {

}

// TODO
//
func (p *Pool) Swap() {

}

// TODO
//
func (p *Pool) Collect() {

}

f
}

