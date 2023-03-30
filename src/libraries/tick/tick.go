package tick

import (
	"math/big"

	"github.com/chris-aubin/Uniswap-Simulator/src/libraries/liquidityMath"
	"github.com/chris-aubin/Uniswap-Simulator/src/libraries/tickMath"
)

var (
	MaxUint128, _ = new(big.Int).SetString("0xffffffffffffffffffffffffffffffff", 16)
)

// Info stored for each initialised individual tick
type Tick struct {
	// The total position liquidity that references this tick
	LiquidityGross *big.Int
	// Amount of net liquidity added (subtracted) when tick is crossed from left to right (right to left),
	LiquidityNet *big.Int
	// Fee growth per unit of liquidity on the _other_ side of this tick (relative to the current tick)
	// Only has relative meaning, not absolute — the value depends on when the tick is initialised
	FeeGrowthOutside0X128 *big.Int
	FeeGrowthOutside1X128 *big.Int
	// The cumulative tick value on the other side of the tick
	TickCumulativeOutside *big.Int
	// The seconds per unit of liquidity on the _other_ side of this tick (relative to the current tick)
	// Only has relative meaning, not absolute — the value depends on when the tick is initialised
	SecondsPerLiquidityOutsideX128 *big.Int
	// The seconds spent on the other side of the tick (relative to the current tick)
	// Only has relative meaning, not absolute — the value depends on when the tick is initialised
	SecondsOutside int
	// True iff the tick is initialised, i.e. the value is exactly equivalent to the expression liquidityGross != 0
	// These 8 bits are set to prevent fresh sstores when crossing newly initialised ticks
	Initialised bool
}

// Contains all all tick information for initialised ticks
type Ticks struct {
	// Maps tick index to tick data
	TickData map[int]*Tick
}

// Derives max liquidity per tick from given tick spacing
// Accepts tickSpacing, the amount of required tick separation. A tickSpacing of 3
// requires ticks to be initialised every 3rd tick i.e., ..., -6, -3, 0, 3, 6, ...
// Returns the max liquidity per tick
func (t *Ticks) tickSpacingToMaxLiquidityPerTick(tickSpacing int) *big.Int {
	minTick := (tickMath.MinTick / tickSpacing) * tickSpacing
	maxTick := (tickMath.MaxTick / tickSpacing) * tickSpacing
	numTicks := ((maxTick - minTick) / tickSpacing) + 1
	result := new(big.Int).Div(MaxUint128, big.NewInt(int64(numTicks)))
	return result
}

/// Retrieves fee growth data
/// Accepts tickLower, the lower tick boundary of the position
/// Accepts tickUpper, the upper tick boundary of the position
/// Accepts tickCurrent, the current tick
/// Accepts feeGrowthGlobal0X128, the all-time global fee growth, per unit of liquidity, in token0
/// Accepts feeGrowthGlobal1X128, the all-time global fee growth, per unit of liquidity, in token1
/// Returns feeGrowthInside0X128, the all-time fee growth in token0, per unit of liquidity, inside the position's tick boundaries
/// Returns feeGrowthInside1X128, the all-time fee growth in token1, per unit of liquidity, inside the position's tick boundaries
func (t *Ticks) getFeeGrowthInside(
	tickLower,
	tickUpper,
	tickCurrent int,
	feeGrowthGlobal0X128,
	feeGrowthGlobal1X128 *big.Int,
) (*big.Int, *big.Int) {
	lower := t.TickData[tickLower]
	upper := t.TickData[tickUpper]

	// Calculate fee growth below
	feeGrowthBelow0X128 := new(big.Int)
	feeGrowthBelow1X128 := new(big.Int)
	if tickCurrent >= tickLower {
		feeGrowthBelow0X128 = lower.FeeGrowthOutside0X128
		feeGrowthBelow1X128 = lower.FeeGrowthOutside1X128
	} else {
		feeGrowthBelow0X128.Sub(feeGrowthGlobal0X128, lower.FeeGrowthOutside0X128)
		feeGrowthBelow1X128.Sub(feeGrowthGlobal1X128, lower.FeeGrowthOutside1X128)
	}

	// Calculate fee growth above
	feeGrowthAbove0X128 := new(big.Int)
	feeGrowthAbove1X128 := new(big.Int)
	if tickCurrent < tickUpper {
		feeGrowthAbove0X128 = upper.FeeGrowthOutside0X128
		feeGrowthAbove1X128 = upper.FeeGrowthOutside1X128
	} else {
		feeGrowthAbove0X128.Sub(feeGrowthGlobal0X128, upper.FeeGrowthOutside0X128)
		feeGrowthAbove1X128.Sub(feeGrowthGlobal1X128, upper.FeeGrowthOutside1X128)
	}

	// Calculate fee growth inside
	feeGrowthInside0X128 := new(big.Int)
	feeGrowthInside1X128 := new(big.Int)

	feeGrowthInside0X128.Sub(feeGrowthGlobal0X128, new(big.Int).Sub(feeGrowthBelow0X128, feeGrowthAbove0X128))
	feeGrowthInside1X128.Sub(feeGrowthGlobal1X128, new(big.Int).Sub(feeGrowthBelow1X128, feeGrowthAbove1X128))

	return feeGrowthInside0X128, feeGrowthInside1X128
}

// Updates a tick and returns true if the tick was flipped from initialised to uninitialised, or vice versa
// Accepts tick, the index of the tick that will be updated
// Accepts tickCurrent, the index of the current tick
// Accepts liquidityDelta, the (new) amount of liquidity to be added (subtracted) when tick is crossed from left to right (right to left)
// Accepts feeGrowthGlobal0X128, the all-time global fee growth, per unit of liquidity, in token0
// Accepts feeGrowthGlobal1X128, the all-time global fee growth, per unit of liquidity, in token1
// Accepts secondsPerLiquidityCumulativeX128, the all-time seconds per max(1, liquidity) of the pool
// Accepts tickCumulative, the tick * time elapsed since the pool was first initialised
// Accepts time, the current block timestamp cast to a uint32
// Accepts upper, a boolean that is true for updating a position's upper tick, or false for updating a position's lower tick
// Accepts maxLiquidity, the maximum liquidity allocation for a single tick
// Returns flipped, a boolean that indicates whether the tick was flipped from initialised to uninitialised, or vice versa
func (t *Ticks) update(
	tick,
	tickCurrent,
	time int,
	liquidityDelta,
	feeGrowthGlobal0X128,
	feeGrowthGlobal1X128,
	secondsPerLiquidityCumulativeX128,
	tickCumulative,
	maxLiquidity *big.Int,
	upper bool,
) bool {
	info := t.TickData[tick]

	liquidityGrossBefore := info.LiquidityGross
	liquidityGrossAfter := liquidityMath.AddDelta(liquidityGrossBefore, liquidityDelta)

	if liquidityGrossAfter.Cmp(maxLiquidity) == 1 {
		panic("TICK: LO")
	}

	flipped := (liquidityGrossAfter.Cmp(big.NewInt(0)) == 1) != (liquidityGrossBefore.Cmp(big.NewInt(0)) == 1)

	if liquidityGrossBefore.Cmp(big.NewInt(0)) == 1 {
		// By convention, Uniswap assumes that all growth before a tick was initialised happened _below_ the tick
		if tick <= tickCurrent {
			info.FeeGrowthOutside0X128 = feeGrowthGlobal0X128
			info.FeeGrowthOutside1X128 = feeGrowthGlobal1X128
			info.SecondsPerLiquidityOutsideX128 = secondsPerLiquidityCumulativeX128
			info.TickCumulativeOutside = tickCumulative
			info.SecondsOutside = time
		}
		info.Initialised = true
	}

	info.LiquidityGross = liquidityGrossAfter

	if upper {
		info.LiquidityNet.Sub(info.LiquidityNet, liquidityDelta)
	} else {
		info.LiquidityNet.Add(info.LiquidityNet, liquidityDelta)
	}

	return flipped
}

// Clears data for a particular tick
// Accepts tick, the tick index of the tick that will be cleared
func (t *Ticks) clear(tick int) {
	delete(t.TickData, tick)
}

// Transitions to next tick as needed by price movement
// Accepts tick, the destination tick of the transition
// Accepts feeGrowthGlobal0X128, the all-time global fee growth, per unit of liquidity, in token0
// Accepts feeGrowthGlobal1X128, the all-time global fee growth, per unit of liquidity, in token1
// Accepts secondsPerLiquidityCumulativeX128, the current seconds per liquidity
// Accepts tickCumulative, the tick * time elapsed since the pool was first initialised
// Accepts time, the current block.timestamp
// Returns liquidityNet, the amount of liquidity added (subtracted) when tick is crossed from left to right (right to left)
func (t *Ticks) cross(
	tick,
	time int,
	feeGrowthGlobal0X128,
	feeGrowthGlobal1X128,
	secondsPerLiquidityCumulativeX128,
	tickCumulative *big.Int,
) *big.Int {
	info := t.TickData[tick]
	info.FeeGrowthOutside0X128.Sub(feeGrowthGlobal0X128, info.FeeGrowthOutside0X128)
	info.FeeGrowthOutside1X128.Sub(feeGrowthGlobal1X128, info.FeeGrowthOutside1X128)
	info.SecondsPerLiquidityOutsideX128.Sub(secondsPerLiquidityCumulativeX128, info.SecondsPerLiquidityOutsideX128)
	info.TickCumulativeOutside.Sub(tickCumulative, info.TickCumulativeOutside)
	info.SecondsOutside = time - info.SecondsOutside
	liquidityNet := info.LiquidityNet
	return liquidityNet
}

func Init() *Ticks {
	return &Ticks{
		TickData: make(map[int]*Tick),
	}
}