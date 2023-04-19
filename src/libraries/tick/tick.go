// Package tick simulates the Uniswap tick library.
//
// Contains functions for managing tick processes and relevant calculations.
package tick

import (
	"math/big"
	"strconv"

	"github.com/chris-aubin/Uniswap-Simulator/src/libraries/constants"
	"github.com/chris-aubin/Uniswap-Simulator/src/libraries/liquidityMath"
)

// State stored for each initialized tick.
type Tick struct {
	// The total position liquidity that references this tick.
	LiquidityGross *big.Int
	// Amount of net liquidity added when tick is crossed from left to right (or
	// subtracted when the tick is crossed from right to left).
	LiquidityNet *big.Int
	// Fee growth per unit of liquidity on the _other_ side of this tick
	// (relative to the current tick). Only has relative meaning, not absolute —
	// the value depends on when the tick is initialized.
	FeeGrowthOutside0X128 *big.Int
	FeeGrowthOutside1X128 *big.Int
	// True iff the tick is initialized, i.e. the value is exactly equivalent to
	// the expression liquidityGross != 0. These 8 bits are set to prevent fresh
	// stores when crossing newly initialized ticks.
	Initialized bool
}

// Contains all all tick information for a set of initialized ticks.
type Ticks struct {
	// Maps tick index to tick data.
	TickData map[int]*Tick
}

// Takes in a map from string to Tick and returns a Ticks struct (i.e. a map
// from int to Tick). Used when loading tick data from JSON.
func TicksTempToTicks(ticks *map[string]Tick) *Ticks {
	tickData := make(map[int]*Tick)
	for k, v := range *ticks {
		tickIdx, _ := strconv.ParseInt(k, 10, 64)
		tick := &Tick{
			LiquidityGross:        v.LiquidityGross,
			LiquidityNet:          v.LiquidityNet,
			FeeGrowthOutside0X128: v.FeeGrowthOutside0X128,
			FeeGrowthOutside1X128: v.FeeGrowthOutside1X128,
			Initialized:           v.Initialized,
		}
		tickData[int(tickIdx)] = tick
	}
	return &Ticks{
		TickData: tickData,
	}
}

// Calculates max liquidity per tick from given tick spacing.
//
// Arguments:
// tickSpacing -- The amount of required tick separation
//
// Returns:
// The max liquidity per tick
func tickSpacingToMaxLiquidityPerTick(tickSpacing int) *big.Int {
	minTick := (constants.MinTick / tickSpacing) * tickSpacing
	maxTick := (constants.MaxTick / tickSpacing) * tickSpacing
	numTicks := ((maxTick - minTick) / tickSpacing) + 1
	result := new(big.Int).Div(constants.MaxUint128, big.NewInt(int64(numTicks)))
	return result
}

// Solidity automatically initializes all values in maps, so this simulates that
// behavior (in the case that a tick doesn't exist at a particular index this
// function initialises a new tick, associates it with that index and returns it).
//
// Arguments:
// tick -- The tick index
//
// Returns:
// The tick data at the given index
func (t *Ticks) Get(tick int) *Tick {
	tickInfo, found := t.TickData[tick]
	if found {
		return tickInfo
	} else {
		t.TickData[tick] = &Tick{
			LiquidityGross:        big.NewInt(0),
			LiquidityNet:          big.NewInt(0),
			FeeGrowthOutside0X128: big.NewInt(0),
			FeeGrowthOutside1X128: big.NewInt(0),
			Initialized:           false,
		}
		return t.TickData[tick]
	}
}

// Retrieves fee growth data.
//
// Arguments:
// tickLower            -- The lower tick boundary of the position
// tickUpper            -- The upper tick boundary of the position
// tickCurrent          -- The current tick
// feeGrowthGlobal0X128 -- The all-time global fee growth, per unit of liquidity,
//                         in token0
// feeGrowthGlobal1X128 -- The all-time global fee growth, per unit of liquidity,
//                         in token1
//
// Returns:
// feeGrowthInside0X128 -- The all-time fee growth in token0, per unit of
//                         liquidity, inside the position's tick boundaries
// feeGrowthInside1X128 -- The all-time fee growth in token1, per unit of
//                         liquidity, inside the position's tick boundaries
func (t *Ticks) GetFeeGrowthInside(tickLower, tickUpper, tickCurrent int, feeGrowthGlobal0X128, feeGrowthGlobal1X128 *big.Int) (*big.Int, *big.Int) {
	lower := t.Get(tickLower)
	upper := t.Get(tickUpper)

	// Calculate fee growth below
	feeGrowthBelow0X128 := new(big.Int)
	feeGrowthBelow1X128 := new(big.Int)
	if tickCurrent >= tickLower {
		feeGrowthBelow0X128 = lower.FeeGrowthOutside0X128
		feeGrowthBelow1X128 = lower.FeeGrowthOutside1X128
	} else {
		feeGrowthBelow0X128 = new(big.Int).Sub(feeGrowthGlobal0X128, lower.FeeGrowthOutside0X128)
		feeGrowthBelow1X128 = new(big.Int).Sub(feeGrowthGlobal1X128, lower.FeeGrowthOutside1X128)
	}

	// Calculate fee growth above
	feeGrowthAbove0X128 := new(big.Int)
	feeGrowthAbove1X128 := new(big.Int)
	if tickCurrent < tickUpper {
		feeGrowthAbove0X128 = upper.FeeGrowthOutside0X128
		feeGrowthAbove1X128 = upper.FeeGrowthOutside1X128
	} else {
		feeGrowthAbove0X128 = new(big.Int).Sub(feeGrowthGlobal0X128, upper.FeeGrowthOutside0X128)
		feeGrowthAbove1X128 = new(big.Int).Sub(feeGrowthGlobal1X128, upper.FeeGrowthOutside1X128)
	}

	// Calculate fee growth inside
	feeGrowthInside0X128 := new(big.Int).Sub(feeGrowthGlobal0X128, new(big.Int).Add(feeGrowthBelow0X128, feeGrowthAbove0X128))
	feeGrowthInside1X128 := new(big.Int).Sub(feeGrowthGlobal1X128, new(big.Int).Add(feeGrowthBelow1X128, feeGrowthAbove1X128))

	// Simulate solidity underflow
	if feeGrowthInside0X128.Cmp(big.NewInt(0)) <= -1 {
		feeGrowthInside0X128 = big.NewInt(0).Add(feeGrowthInside0X128, constants.Q256)
	}
	if feeGrowthInside1X128.Cmp(big.NewInt(0)) <= -1 {
		feeGrowthInside1X128 = big.NewInt(0).Add(feeGrowthInside1X128, constants.Q256)
	}

	return feeGrowthInside0X128, feeGrowthInside1X128
}

// Updates a tick and returns true if the tick was flipped from initialized to
// uninitialized, or vice versa.
//
// Arguments:
// tick                 -- The index of the tick that will be updated
// tickCurrent          -- The index of the current tick
// liquidityDelta       -- The (new) amount of liquidity to be added (subtracted)
//                         when tick is crossed from left to right (right to left)
// feeGrowthGlobal0X128 -- The all-time global fee growth, per unit of liquidity,
//                         in token0
// feeGrowthGlobal1X128 -- The all-time global fee growth, per unit of liquidity,
//                         in token1
// upper                -- A boolean that is true for updating a position's upper
//                         tick, or false for updating a position's lower tick
// maxLiquidity         -- The maximum liquidity allocation for a single tick
//
// Returns:
// flipped              -- A boolean that indicates whether the tick was flipped
func (t *Ticks) Update(tick, tickCurrent int, liquidityDelta, feeGrowthGlobal0X128, feeGrowthGlobal1X128, maxLiquidity *big.Int, upper bool) bool {
	info := t.Get(tick)
	liquidityGrossBefore := info.LiquidityGross
	liquidityGrossAfter := liquidityMath.AddDelta(liquidityGrossBefore, liquidityDelta)

	if liquidityGrossAfter.Cmp(maxLiquidity) >= 1 {
		panic("tick.Update: Tick liquidity cannot exceed max liquidity per tick")
	}

	flipped := (liquidityGrossAfter.Cmp(big.NewInt(0)) == 0) != (liquidityGrossBefore.Cmp(big.NewInt(0)) == 0)

	if liquidityGrossBefore.Cmp(big.NewInt(0)) == 0 {
		// By convention, Uniswap assumes that all growth before a tick was
		// initialized happened below the tick
		if tick <= tickCurrent {
			info.FeeGrowthOutside0X128 = feeGrowthGlobal0X128
			info.FeeGrowthOutside1X128 = feeGrowthGlobal1X128
		}
		info.Initialized = true
	}

	info.LiquidityGross = liquidityGrossAfter

	// When the lower (upper) tick is crossed left to right (right to left),
	// liquidity must be added (removed)
	if upper {
		info.LiquidityNet = new(big.Int).Sub(info.LiquidityNet, liquidityDelta)
	} else {
		info.LiquidityNet = new(big.Int).Add(info.LiquidityNet, liquidityDelta)
	}

	return flipped
}

// Clears data for a particular tick
//
// Arguments:
// tick -- The tick index of the tick that will be cleared
func (t *Ticks) Clear(tick int) {
	delete(t.TickData, tick)
}

// Transitions to next tick as needed by price movement.
//
// Arguments:
// tick                 -- The destination tick of the transition
// feeGrowthGlobal0X128 -- The all-time global fee growth, per unit of liquidity,
//                         in token0
// feeGrowthGlobal1X128 -- The all-time global fee growth, per unit of liquidity,
//                         in token1
// Returns:
// liquidityNet         -- The amount of liquidity added (subtracted) when tick
//                         is crossed from left to right (right to left)

func (t *Ticks) Cross(tick int, feeGrowthGlobal0X128, feeGrowthGlobal1X128 *big.Int) *big.Int {
	info := t.Get(tick)
	info.FeeGrowthOutside0X128 = new(big.Int).Sub(feeGrowthGlobal0X128, info.FeeGrowthOutside0X128)
	info.FeeGrowthOutside1X128 = new(big.Int).Sub(feeGrowthGlobal1X128, info.FeeGrowthOutside1X128)
	liquidityNet := info.LiquidityNet
	return liquidityNet
}
