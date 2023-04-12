package pool

import (
	"fmt"
	"math/big"

	"github.com/chris-aubin/Uniswap-Simulator/src/libraries/constants"
	"github.com/chris-aubin/Uniswap-Simulator/src/libraries/fullMath"
	"github.com/chris-aubin/Uniswap-Simulator/src/libraries/liquidityMath"
	"github.com/chris-aubin/Uniswap-Simulator/src/libraries/position"
	"github.com/chris-aubin/Uniswap-Simulator/src/libraries/sqrtPriceMath"
	"github.com/chris-aubin/Uniswap-Simulator/src/libraries/swapMath"
	"github.com/chris-aubin/Uniswap-Simulator/src/libraries/tick"
	"github.com/chris-aubin/Uniswap-Simulator/src/libraries/tickMath"
)

// Part of pool state.
type Slot0 struct {
	// The current price.
	SqrtPriceX96 *big.Int
	// The current tick.
	Tick int
	// The current protocol fee as a percentage of the swap fee taken on
	// withdrawal. Represented as an integer denominator (1/x)%
	FeeProtocol int
}

// Accumulated protocol fees in token0/token1 units (fees that could be
// collected by Uniswap governance).
type ProtocolFees struct {
	Token0 *big.Int
	Token1 *big.Int
}

// Pool state.
type Pool struct {
	Token0              string
	Token1              string
	Fee                 int
	TickSpacing         int
	MaxLiquidityPerTick *big.Int
	Slot0               *Slot0
	// FeeGrowthGlobal0X128 and FeeGrowthGlobal1X128 represent the total amount
	// of fees that have been earned per unit of virtual liquidity (L), over the
	// entire history of the contract. This is the same as the total amount of
	// fees that would have been earned by 1 unit of unbounded liquidity that
	// was deposited when the contract was first initialized.
	FeeGrowthGlobal0X128 *big.Int
	FeeGrowthGlobal1X128 *big.Int
	ProtocolFees         *ProtocolFees
	Liquidity            *big.Int
	// Tick-indexed state, as per section 6.3 in Uniswap V3 Whitepaper. This is
	// a mapping from tick index to a Tick struct that contains information
	// about that tick (see the tick package for more).
	Ticks *tick.Ticks
	// Position-indexed state, as per section 6.4 in Uniswap V3 Whitepaper. In
	// the deployed contract, this is a mapping from the hash of a position's
	// owner's address, tickLower, and tickUpper (in byte form) to a Position.
	// In this simulator it is implemented it as a mapping from a string, which
	// is the concatenation of the owner's address, tickUpper and tickLower
	// concatenate, to a Position (see the position package for more).
	Positions map[string]*position.Position
	// Balance of token0 and token1 held by the pool. Not part of state in the
	// deployed contract (the deployed contract checks the balance of the token
	// owned by the pool address).
	Balance0 *big.Int
	Balance1 *big.Int
}

// Common checks for valid tick inputs.
func checkTicks(tickLower int, tickUpper int) {
	if tickLower >= tickUpper {
		panic("Pool.checkTicks: tickLower > tickUpper")
	}
	if tickLower < constants.MinTick {
		panic("Pool.checkTicks: tickLower < MINTICK")
	}
	if tickUpper > constants.MaxTick {
		panic("Pool.checkTicks: tickUpper > MAXTICK")
	}
}

// Finds the next initialized tick contained in the same word (or adjacent
// word i.e. within 256 ticks) as the tick that is either to the left (less
// than or equal to) or right (greater than) of the given tick. This function is
// used because the deployed contract uses a bitMap to efficiently store and
// check which ticks are initialized. The bitMap is made up of multiple words.
// This is unnecessary in the simulator, but we use it to avoid making too many
// changes to the deployed contract.
//
// Arguments:
// tick         -- the starting tick
// tickSpacing  -- the spacing between usable ticks for the pool
// lte          -- a bool that indicates whether to search for the next
//                 initialized tick to the left (less than or equal to the
//                 starting tick)
// Returns:
// next         -- the next initialized or uninitialized tick up to 256 ticks
//                 away from the current tick
// initialized  -- a bool that indicates whether or not next is initialized
//                 (because the function only searches within up to 256 ticks)
func (p *Pool) nextInitializedTickWithinOneWord(tick, tickSpacing int, lte bool) (next int, initialized bool) {
	// Find the boundaries of the word that would contain the tick.
	wordLowerBound := tick - tick%256
	wordUpperBound := wordLowerBound + 255
	// Adjust for the tickSpacing.
	compressed := tick / tickSpacing
	if tick < 0 && tick%tickSpacing != 0 {
		compressed = compressed - 1
	}

	if lte {
		// Search for the closest initialized tick, within the word, with
		// tick_idx less than or equal to tick.
		for i := compressed; i >= wordLowerBound; i-- {
			tick_idx := i * tickSpacing
			tick, err := p.Ticks.TickData[tick_idx]
			if err {
				if i == wordLowerBound {
					return tick_idx, false
				}
				continue
			} else if tick.Initialized {
				return tick_idx, true
			} else {
				if i == wordLowerBound {
					return tick_idx, false
				}
				continue
			}
		}
	} else {
		// Search for the closest initialized tick, within the word, with
		// tick_idx greater than tick.
		for i := compressed + 1; i <= wordUpperBound; i++ {
			tick_idx := i * tickSpacing
			tick, err := p.Ticks.TickData[tick_idx]
			if err {
				if i == wordUpperBound {
					return tick_idx, false
				}
				continue
			} else if tick.Initialized {
				return tick_idx, true
			} else {
				if i == wordUpperBound {
					return tick_idx, false
				}
				continue
			}
		}
	}
	panic("Pool.nextInitializedTickWithinOneWord: Unreachable")
}

type modifyPositionParams struct {
	// the address that owns the position
	Owner string
	// the lower and upper tick of the position
	TickLower int
	TickUpper int
	// any change in liquidity
	LiquidityDelta *big.Int
}

// Effect some changes to a position.
//
// Arguments:
// params   --  An modifyPositionParams type that contains the position details
//              and the changes to the position's liquidity to effect
//
// Returns:
// position -- The updated position
// amount0  -- the amount of token0 owed to the pool (negative if the pool
//             should pay the recipient)
// amount1  -- the amount of token1 owed to the pool (negative if the pool
//             should pay the recipient)
func (p *Pool) modifyPosition(params *modifyPositionParams) (position *position.Position, amount0 *big.Int, amount1 *big.Int) {
	checkTicks(params.TickLower, params.TickUpper)
	slot0 := p.Slot0

	position = p.updatePosition(params.Owner, params.TickLower, params.TickUpper, slot0.Tick, params.LiquidityDelta)

	if params.LiquidityDelta.Cmp(big.NewInt(0)) != 0 {
		if slot0.Tick < params.TickLower {
			// Current tick is below the passed range; liquidity can only become in range by crossing from left to
			// right, when we'll need _more_ token0 (it's becoming more valuable) so user must provide it
			amount0 = sqrtPriceMath.GetAmount0DeltaNoBool(
				tickMath.GetSqrtRatioAtTick(params.TickLower),
				tickMath.GetSqrtRatioAtTick(params.TickUpper),
				params.LiquidityDelta,
			)
		} else if slot0.Tick < params.TickUpper {
			// Current tick is inside the passed range
			liquidityBefore := p.Liquidity

			amount0 = sqrtPriceMath.GetAmount0DeltaNoBool(
				slot0.SqrtPriceX96,
				tickMath.GetSqrtRatioAtTick(params.TickUpper),
				params.LiquidityDelta,
			)
			amount1 = sqrtPriceMath.GetAmount1DeltaNoBool(
				tickMath.GetSqrtRatioAtTick(params.TickLower),
				slot0.SqrtPriceX96,
				params.LiquidityDelta,
			)

			p.Liquidity = liquidityMath.AddDelta(liquidityBefore, params.LiquidityDelta)
		} else {
			// Current tick is above the passed range; liquidity can only become in range by crossing from right to
			// left, when we'll need _more_ token1 (it's becoming more valuable) so user must provide it
			amount1 = sqrtPriceMath.GetAmount1DeltaNoBool(
				tickMath.GetSqrtRatioAtTick(params.TickLower),
				tickMath.GetSqrtRatioAtTick(params.TickUpper),
				params.LiquidityDelta,
			)
		}
	}
	return
}

// Gets and updates a position with the given liquidity delta
//
// Arguments:
// owner     -- the owner of the position
// tickLower -- the lower tick of the position's tick range
// tickUpper -- the upper tick of the position's tick range
//
// Returns:
// position  -- the updated position
func (p *Pool) updatePosition(owner string, tickLower, tickUpper, tick int, liquidityDelta *big.Int) (position *position.Position) {
	position_key := fmt.Sprintf("%s%d%d", owner, tickLower, tickUpper)
	position, err := p.Positions[position_key]
	if err {
		message := fmt.Sprintf("pool.updatePosition: Position %s does not exist", position_key)
		panic(message)
	}

	feeGrowthGlobal0X128 := p.FeeGrowthGlobal0X128
	feeGrowthGlobal1X128 := p.FeeGrowthGlobal1X128

	// Used to determine if we need to clear tickLower/ tickUpper after the
	// position is updated/
	var flippedLower bool
	var flippedUpper bool
	if liquidityDelta.Cmp(big.NewInt(0)) != 0 {

		flippedLower = p.Ticks.Update(
			tickLower,
			tick,
			liquidityDelta,
			feeGrowthGlobal0X128,
			feeGrowthGlobal1X128,
			p.MaxLiquidityPerTick,
			false,
		)

		flippedUpper = p.Ticks.Update(
			tickUpper,
			tick,
			liquidityDelta,
			feeGrowthGlobal0X128,
			feeGrowthGlobal1X128,
			p.MaxLiquidityPerTick,
			true,
		)
	}

	feeGrowthInside0X128, feeGrowthInside1X128 := p.Ticks.GetFeeGrowthInside(
		tickLower,
		tickUpper,
		tick,
		feeGrowthGlobal0X128,
		feeGrowthGlobal1X128,
	)

	position.Update(liquidityDelta, feeGrowthInside0X128, feeGrowthInside1X128)

	// Clear any tick data that is no longer needed
	if liquidityDelta.Cmp(big.NewInt(0)) <= -1 {
		if flippedLower {
			p.Ticks.Clear(tickLower)
		}
		if flippedUpper {
			p.Ticks.Clear(tickUpper)
		}
	}

	return position
}

// Mints liquidity for the given recipient in the given tick range (either
// a new position or increases the liquidity of an existing position)
//
// Arguments:
// recipient -- the recipient of the minted liquidity
// tickLower -- the lower tick of the position's tick range
// tickUpper -- the upper tick of the position's tick range
// amount    -- the amount of liquidity to mint
//
// Returns:
// amount0   -- the amount of token0 to transfer to the recipient
// amount1   -- the amount of token1 to transfer to the recipient
func (p *Pool) Mint(recipient string, tickLower, tickUpper int, amount *big.Int) (amount0, amount1 *big.Int) {
	// Quick sanity checks
	checkTicks(tickLower, tickUpper)
	if amount.Cmp(big.NewInt(0)) <= 0 {
		message := fmt.Sprintf("pool.Mint: Amount %s must be greater than 0", amount)
		panic(message)
	}

	// Get the amount of token0 and token1 owed to the pool (negative if the
	// pool owes the recipient)
	_, amount0, amount1 = p.modifyPosition(
		&modifyPositionParams{
			Owner:          recipient,
			TickLower:      tickLower,
			TickUpper:      tickUpper,
			LiquidityDelta: amount,
		})

	// balance0Before := new(big.Int)
	// balance1Before := new(big.Int)
	// if amount0.Cmp(big.NewInt(0)) >= 1 {
	// 	balance0Before = p.Balance0
	// }
	// if amount1.Cmp(big.NewInt(0)) >= 1 {
	// 	balance1Before = p.Balance1
	// }

	// Let minter know how much token0 and token1 they must pay to mint the position
	// (This is what the IUniswapV3MintCallback does in the UniswapV3Pool contract)
	// If payment was made mint position
	// if (amount0 > 0) require(balance0Before.add(amount0) <= balance0(), 'M0');
	// if (amount1 > 0) require(balance1Before.add(amount1) <= balance1(), 'M1');
	return
}

// Collects tokens owed to the given position
//
// Arguments:
// owner            -- the owner of the position
// tickLower        -- the lower tick of the position's tick range
// tickUpper        -- the upper tick of the position's tick range
// amount0Requested -- the amount of token0 to collect
// amount1Requested -- the amount of token1 to collect
//
// Returns:
// amount0          -- the amount of token0 collected
// amount1          -- the amount of token1 collected
func (p *Pool) Collect(owner string, tickLower, tickUpper int, amount0Requested, amount1Requested *big.Int) (amount0, amount1 *big.Int) {
	// We don't need to checkTicks here, because invalid positions will never
	// have non-zero tokensOwed
	position_key := fmt.Sprintf("%s%d%d", owner, tickLower, tickUpper)
	position, err := p.Positions[position_key]
	if err {
		message := fmt.Sprintf("pool.Collect: Position %s does not exist", position_key)
		panic(message)
	}

	// If more tokens are requests than are owed to the position then just
	// collect all the tokens owed.
	amount0 = new(big.Int)
	if amount0Requested.Cmp(position.TokensOwed0) >= 1 {
		amount0 = position.TokensOwed0
	} else {
		amount0 = amount0Requested
	}
	amount1 = new(big.Int)
	if amount1Requested.Cmp(position.TokensOwed1) >= 1 {
		amount1 = position.TokensOwed1
	} else {
		amount1 = amount1Requested
	}

	// Update the positions' tokensOwed
	if amount0.Cmp(big.NewInt(0)) >= 1 {
		position.TokensOwed0 = new(big.Int).Sub(position.TokensOwed0, amount0)
	}
	if amount1.Cmp(big.NewInt(0)) >= 1 {
		position.TokensOwed1 = new(big.Int).Sub(position.TokensOwed1, amount1)
	}
	return
}

//
func (p *Pool) Burn(owner string, tickLower, tickUpper int, amount *big.Int) (amount0, amount1 *big.Int) {
	position, amount0, amount1 := p.modifyPosition(
		&modifyPositionParams{
			Owner:          owner,
			TickLower:      tickLower,
			TickUpper:      tickUpper,
			LiquidityDelta: new(big.Int).Neg(amount),
		},
	)

	amount0 = new(big.Int).Neg(amount0)
	amount1 = new(big.Int).Neg(amount1)

	// if (amount0 > 0 || amount1 > 0) {
	if amount0.Cmp(big.NewInt(0)) >= 1 || amount1.Cmp(big.NewInt(0)) >= 1 {
		position.TokensOwed0 = new(big.Int).Add(position.TokensOwed0, amount0)
		position.TokensOwed1 = new(big.Int).Add(position.TokensOwed1, amount1)
	}
	return
}

type SwapCache struct {
	// The protocol fee for the input token.
	FeeProtocol int
	// The liquidity at the beginning of the swap.
	LiquidityStart *big.Int
}

// The top level state of the swap, the results of which are recorded in storage
// at the end.
type SwapState struct {
	// The amount remaining to be swapped in/out of the input/output asset.
	AmountSpecifiedRemaining *big.Int
	// The amount already swapped out/in of the output/input asset.
	AmountCalculated *big.Int
	// The current sqrt(price).
	SqrtPriceX96 *big.Int
	// The tick associated with the current price.
	Tick int
	// The global fee growth of the input token.
	FeeGrowthGlobalX128 *big.Int
	// The amount of the input token paid as protocol fee.
	ProtocolFee *big.Int
	// The current liquidity in range.
	Liquidity *big.Int
}

type StepComputations struct {
	// the price at the beginning of the step
	SqrtPriceStartX96 *big.Int
	// the next tick to swap to from the current tick in the swap direction
	TickNext int
	// whether tickNext is initialized or not
	Initialized bool
	// sqrt(price) for the next tick (1/0)
	SqrtPriceNextX96 *big.Int
	// how much is being swapped in in this step
	AmountIn *big.Int
	// how much is being swapped out
	AmountOut *big.Int
	// how much fee is being paid in
	FeeAmount *big.Int
}

/// @inheritdoc IUniswapV3PoolActions
func (p *Pool) Swap(sender, recipient string, zeroForOne bool, amountSpecified, sqrtPriceLimitX96 *big.Int) (amount0, amount1 *big.Int) {
	if amountSpecified.Cmp(big.NewInt(0)) == 0 {
		message := fmt.Sprintf("pool.Swap: amountSpecified %v must not be 0", amountSpecified)
		panic(message)
	}

	slot0Start := p.Slot0

	var cacheFeeProtocol int
	var stateFeeGrowthGlobalX128 *big.Int
	if zeroForOne {
		if !((sqrtPriceLimitX96.Cmp(slot0Start.SqrtPriceX96) <= -1) && (sqrtPriceLimitX96.Cmp(constants.MinTickBig) >= 1)) {
			panic("pool.Swap: Invalid price limit")
		}
		cacheFeeProtocol = slot0Start.FeeProtocol % 16
		stateFeeGrowthGlobalX128 = p.FeeGrowthGlobal0X128
	} else {
		if !((sqrtPriceLimitX96.Cmp(slot0Start.SqrtPriceX96) >= 1) && (sqrtPriceLimitX96.Cmp(constants.MaxTickBig) <= -1)) {
			panic("pool.Swap: Invalid price limit")
		}
		cacheFeeProtocol = slot0Start.FeeProtocol >> 4
		stateFeeGrowthGlobalX128 = p.FeeGrowthGlobal1X128
	}

	cache := &SwapCache{
		LiquidityStart: p.Liquidity,
		FeeProtocol:    cacheFeeProtocol,
	}

	var exactInput bool
	if amountSpecified.Cmp(big.NewInt(0)) >= 1 {
		exactInput = true
	} else {
		exactInput = false
	}

	state := &SwapState{
		AmountSpecifiedRemaining: amountSpecified,
		AmountCalculated:         big.NewInt(0),
		SqrtPriceX96:             slot0Start.SqrtPriceX96,
		Tick:                     slot0Start.Tick,
		FeeGrowthGlobalX128:      stateFeeGrowthGlobalX128,
		ProtocolFee:              big.NewInt(0),
		Liquidity:                cache.LiquidityStart,
	}

	// Continue swapping as long as we haven't used the entire input/output and
	// haven't reached the price limit
	for state.AmountSpecifiedRemaining.Cmp(big.NewInt(0)) >= 1 && state.SqrtPriceX96.Cmp(sqrtPriceLimitX96) <= -1 {
		step := new(StepComputations)
		step.SqrtPriceStartX96 = state.SqrtPriceX96
		step.TickNext, step.Initialized = p.nextInitializedTickWithinOneWord(
			state.Tick,
			p.TickSpacing,
			zeroForOne,
		)

		// Ensure that we do not overshoot the min/max tick (likely unnecessary
		// in this simulator)
		if step.TickNext < constants.MinTick {
			step.TickNext = constants.MinTick
		} else if step.TickNext > constants.MaxTick {
			step.TickNext = constants.MaxTick
		}

		// Get the price for the next tick
		step.SqrtPriceNextX96 = tickMath.GetSqrtRatioAtTick(step.TickNext)

		// Compute values to swap to the target tick, price limit, or point
		// where input/output amount is exhausted
		var sqrtRatioTargetX96 *big.Int
		if zeroForOne {
			// step.sqrtPriceNextX96 < sqrtPriceLimitX96
			if step.SqrtPriceNextX96.Cmp(sqrtPriceLimitX96) <= -1 {
				sqrtRatioTargetX96 = sqrtPriceLimitX96
			} else {
				sqrtRatioTargetX96 = step.SqrtPriceNextX96
			}
		} else {
			if step.SqrtPriceNextX96.Cmp(sqrtPriceLimitX96) >= 1 {
				sqrtRatioTargetX96 = sqrtPriceLimitX96
			} else {
				sqrtRatioTargetX96 = step.SqrtPriceNextX96
			}
		}

		// compute values to swap to the target tick, price limit, or point where input/output amount is exhausted
		state.SqrtPriceX96, step.AmountIn, step.AmountOut, step.FeeAmount = swapMath.ComputeSwapStep(
			state.SqrtPriceX96,
			sqrtRatioTargetX96,
			state.Liquidity,
			state.AmountSpecifiedRemaining,
			p.Fee,
		)

		if exactInput {
			state.AmountSpecifiedRemaining = new(big.Int).Sub(state.AmountSpecifiedRemaining, new(big.Int).Add(step.AmountIn, step.FeeAmount))
			state.AmountCalculated = new(big.Int).Sub(state.AmountCalculated, step.AmountOut)
		} else {
			state.AmountSpecifiedRemaining = new(big.Int).Add(state.AmountSpecifiedRemaining, step.AmountOut)
			state.AmountCalculated = new(big.Int).Add(state.AmountCalculated, new(big.Int).Add(step.AmountIn, step.FeeAmount))
		}

		// If the protocol fee is on, calculate how much is owed, decrement
		// feeAmount, and increment protocolFee
		if cache.FeeProtocol > 0 {
			delta := new(big.Int).Div(step.FeeAmount, big.NewInt(int64(cache.FeeProtocol)))
			step.FeeAmount = new(big.Int).Sub(step.FeeAmount, delta)
			state.ProtocolFee = new(big.Int).Add(state.ProtocolFee, delta)
		}

		// Update global fee tracker
		// if (state.liquidity > 0)
		if state.Liquidity.Cmp(big.NewInt(0)) >= 1 {
			state.FeeGrowthGlobalX128 = new(big.Int).Add(state.FeeGrowthGlobalX128, fullMath.MulDiv(step.FeeAmount, constants.Q128, state.Liquidity))
		}

		// Shift tick if we reached the next price
		if state.SqrtPriceX96.Cmp(step.SqrtPriceNextX96) == 0 {
			// If the tick is initialized, run the tick transition
			if step.Initialized {
				// Check for the placeholder value, which we replace with the
				// actual value the first time the swap crosses an initialized
				/// tick

				tempFeeGrowthGlobal0X128 := p.FeeGrowthGlobal0X128
				tempFeeGrowthGlobal1X128 := state.FeeGrowthGlobalX128
				if zeroForOne {
					tempFeeGrowthGlobal0X128 = state.FeeGrowthGlobalX128
					tempFeeGrowthGlobal1X128 = p.FeeGrowthGlobal1X128
				}

				liquidityNet := p.Ticks.Cross(
					step.TickNext,
					tempFeeGrowthGlobal0X128,
					tempFeeGrowthGlobal1X128,
				)

				if zeroForOne {
					// if we're moving leftward, we interpret liquidityNet as
					// the opposite sign
					liquidityNet = new(big.Int).Neg(liquidityNet)
				}

				state.Liquidity = liquidityMath.AddDelta(state.Liquidity, liquidityNet)
			}

			if zeroForOne {
				state.Tick = step.TickNext - 1
			} else {
				state.Tick = step.TickNext
			}
		} else if state.SqrtPriceX96.Cmp(step.SqrtPriceStartX96) != 0 {
			// Recompute unless we're on a lower tick boundary (i.e. already
			// transitioned ticks), and haven't moved
			state.Tick = tickMath.GetTickAtSqrtRatio(state.SqrtPriceX96)
		}
	}

	// Update tick if the tick change
	if state.Tick != slot0Start.Tick {
		p.Slot0.SqrtPriceX96 = state.SqrtPriceX96
		p.Slot0.Tick = state.Tick
	} else {
		// Otherwise just update the price
		p.Slot0.SqrtPriceX96 = state.SqrtPriceX96
	}

	// Update liquidity if it changed
	if cache.LiquidityStart.Cmp(state.Liquidity) != 0 {
		p.Liquidity = state.Liquidity
	}

	// Update fee growth global and, if necessary, protocol fees
	if zeroForOne {
		p.FeeGrowthGlobal0X128 = state.FeeGrowthGlobalX128
		if state.ProtocolFee.Cmp(big.NewInt(0)) >= 1 {
			p.ProtocolFees.Token0 = new(big.Int).Add(p.ProtocolFees.Token0, state.ProtocolFee)
		}
	} else {
		p.FeeGrowthGlobal1X128 = state.FeeGrowthGlobalX128
		if state.ProtocolFee.Cmp(big.NewInt(0)) >= 1 {
			p.ProtocolFees.Token1 = new(big.Int).Add(p.ProtocolFees.Token1, state.ProtocolFee)
		}
	}

	if zeroForOne == exactInput {
		amount0 = new(big.Int).Sub(amountSpecified, state.AmountSpecifiedRemaining)
		amount1 = state.AmountCalculated
	} else {
		amount0 = state.AmountCalculated
		amount1 = new(big.Int).Sub(amountSpecified, state.AmountSpecifiedRemaining)
	}

	// // do the transfers and collect payment
	// if (zeroForOne) {
	// 	// if (amount1 < 0) TransferHelper.safeTransfer(token1, recipient, uint256(-amount1));

	// 	balance0Before := p.Balance0
	// 	// IUniswapV3SwapCallback(msg.sender).uniswapV3SwapCallback(amount0, amount1, data);
	// 	// require(balance0Before.add(uint256(amount0)) <= balance0(), 'IIA');
	// } else {
	// 	// if (amount0 < 0) TransferHelper.safeTransfer(token0, recipient, uint256(-amount0));

	// 	balance1Before := p.Balance1
	// 	// IUniswapV3SwapCallback(msg.sender).uniswapV3SwapCallback(amount0, amount1, data);
	// 	// require(balance1Before.add(uint256(amount1)) <= balance1(), 'IIA');
	// }
	return
}
