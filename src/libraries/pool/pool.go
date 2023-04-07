package pool

import (
	"math/big"
	"time"

	"github.com/chris-aubin/Uniswap-Simulator/src/libraries/position"
	"github.com/chris-aubin/Uniswap-Simulator/src/libraries/tick"
	"github.com/chris-aubin/Uniswap-Simulator/src/libraries/constants"
	"github.com/chris-aubin/Uniswap-Simulator/src/libraries/swapMath"
)

// Pool state.
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
	Slot0                Slot0
	// FeeGrowthGlobal0X128 and FeeGrowthGlobal1X128 represent the total amount 
	// of fees that have been earned per unit of virtual liquidity (L), over the
	// entire history of the contract. This is the same as the total amount of 
	// fees that would have been earned by 1 unit of unbounded liquidity that 
	// was deposited when the contract was first initialized.
	FeeGrowthGlobal0X128 *big.Int
	FeeGrowthGlobal1X128 *big.Int
	ProtocolFees         ProtocolFees
	Liquidity            *big.Int
	// Tick-indexed state, as per section 6.3 in Uniswap V3 Whitepaper. This is 
	// a mapping from tick index to a Tick struct that contains information 
	// about that tick (see the tick package for more).
	Ticks				 tick.Ticks
	// Position-indexed state, as per section 6.4 in Uniswap V3 Whitepaper. In 
	// the deployed contract, this is a mapping from the hash of a position's
	// owner's address, tickLower, and tickUpper (in byte form) to a Position. 
	// In this simulator it is implemented it as a mapping from a string, which 
	// is the concatenation of the owner's address, tickUpper and tickLower 
	// concatenate, to a Position (see the position package for more).
	Positions			 map[string]*position.Position
	// Balance of token0 and token1 held by the pool. Not part of state in the 
	// deployed contract (the deployed contract checks the balance of the token 
	// owned by the pool address).
	Balance0 	         *big.Int
	Balance1 	         *big.Int
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
func (p *Pool) nextInitializedTickWithinOneWord(tick, tickSpacing int, lte	bool) (next int, initialized bool) {
	// Find the boundaries of the word that would contain the tick.
	wordLowerBound := tick - tick % 256
	wordUpperBound := wordLowerBound + 255
	// Adjust for the tickSpacing.
	compressed := tick / tickSpacing
	if (tick < 0 && tick % tickSpacing != 0) {
		compressed = compressed - 1
	}

	if (lte) {
		for i := compressed - 1; i >= wordLowerBound; i-- {
			// !!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!
			// INCORRECT - WRONG KIND OF POSITION
			position, err := p.Positions[i*tickSpacing]
			// !!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!
			if err != nil {
				if i == wordLowerBound {
					return wordLowerBound*tickSpacing, false
				}
				continue
			} else if position.Initialized {
				return i*tickSpacing, true
			} else {
				if i == wordLowerBound {
					return wordLowerBound*tickSpacing, false
				}
				continue
			}
		}
	} else {
		for i := compressed + 1; i <= wordUpperBound; i++ {
			// !!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!
			// INCORRECT - WRONG KIND OF POSITION
			position, err := p.Positions[i*tickSpacing]
			// !!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!
			if err != nil {
				if i == wordUpperBound {
					return wordUpperBound*tickSpacing, false
				}
				continue
			} else if position.Initialized {
				return i*tickSpacing, true
			} else {
				if i == wordUpperBound {
					return wordUpperBound*tickSpacing, false
				}
				continue
			}
		}
	}
}

type modifyPositionParams struct {
	// the address that owns the position
	owner string
	// the lower and upper tick of the position
	tickLower int
	tickUpper int
	// any change in liquidity
	liquidityDelta *big.Int
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
//             should pay the recipient
// amount1  -- the amount of token1 owed to the pool (negative if the pool 
//             should pay the recipient
func (p *Pool) modifyPosition(params *modifyPositionParams) (position *position.Position, amount0 *big.Int, amount1 *big.Int) {
	checkTicks(params.tickLower, params.tickUpper)
	slot0 = p.Slot0

	position = updatePosition(
		params.owner,
		params.tickLower,
		params.tickUpper,
		params.liquidityDelta,
		slot0.tick,
	)

	amount0 := new(big.Int)
	amount1 := new(big.Int)

	if (params.liquidityDelta != 0) {
		if (slot0.tick < params.tickLower) {
			// Current tick is below the passed range; liquidity can only become in range by crossing from left to
			// right, when we'll need _more_ token0 (it's becoming more valuable) so user must provide it
			amount0 = SqrtPriceMath.getAmount0Delta(
				TickMath.getSqrtRatioAtTick(params.tickLower),
				TickMath.getSqrtRatioAtTick(params.tickUpper),
				params.liquidityDelta,
			)
		} else if (_slot0.tick < params.tickUpper) {
			// Current tick is inside the passed range
			liquidityBefore := liquidity

			amount0 = SqrtPriceMath.getAmount0Delta(
				slot0.sqrtPriceX96,
				TickMath.getSqrtRatioAtTick(params.tickUpper),
				params.liquidityDelta,
			)
			amount1 = SqrtPriceMath.getAmount1Delta(
				TickMath.getSqrtRatioAtTick(params.tickLower),
				slot0.sqrtPriceX96,
				params.liquidityDelta,
			)

			liquidity = LiquidityMath.addDelta(liquidityBefore, params.liquidityDelta);
		} else {
			// Current tick is above the passed range; liquidity can only become in range by crossing from right to
			// left, when we'll need _more_ token1 (it's becoming more valuable) so user must provide it
			amount1 = SqrtPriceMath.getAmount1Delta(
				TickMath.getSqrtRatioAtTick(params.tickLower),
				TickMath.getSqrtRatioAtTick(params.tickUpper),
				params.liquidityDelta,
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
func (p *Pool) updatePosition(owner string, tickLower, tickUpper, tick int, liquidityDelta *big.Int,) (position *position.Position) {
	position, err = positions[fmt.Sprintf("%s%d%d", owner, tickLower, tickUpper)]
	if err != nil {
		panic()
	}

	feeGrowthGlobal0X128 := feeGrowthGlobal0X128
	feeGrowthGlobal1X128 := feeGrowthGlobal1X128

	// Used to determine if we need to clear tickLower/ tickUpper after the 
	// position is updated/
	var flippedLower bool
	var flippedUpper bool
	if (liquidityDelta != 0) {
		time := blockTimestamp()

		flippedLower = ticks.update(
			tickLower,
			tick,
			liquidityDelta,
			feeGrowthGlobal0X128,
			feeGrowthGlobal1X128,
			maxLiquidityPerTick,
			false,
		)

		flippedUpper = ticks.update(
			tickUpper,
			tick,
			liquidityDelta,
			feeGrowthGlobal0X128,
			feeGrowthGlobal1X128,
			maxLiquidityPerTick,
			true,
		)
	}

	feeGrowthInside0X128, feeGrowthInside1X128 := ticks.getFeeGrowthInside(
		tickLower,
		tickUpper,
		tick,
		feeGrowthGlobal0X128,
		feeGrowthGlobal1X128,
	)

	position.update(liquidityDelta, feeGrowthInside0X128, feeGrowthInside1X128);

	// Clear any tick data that is no longer needed
	if (liquidityDelta < 0) {
		if (flippedLower) {
			ticks.clear(tickLower);
		}
		if (flippedUpper) {
			ticks.clear(tickUpper);
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
		panic()
	}

	_, amount0, amount1 := modifyPosition(
		&ModifyPositionParams{
			owner:          recipient,
			tickLower:      tickLower,
			tickUpper:      tickUpper,
			liquidityDelta: amount,
		})

	balance0Before := new(big.Int)
	balance1Before := new(big.Int)
	if (amount0 > 0) {
		balance0Before = p.Balance0
	}
	if (amount1 > 0) {
		balance1Before = p.Balance1
	}
	// IUniswapV3MintCallback(msg.sender).uniswapV3MintCallback(amount0, amount1, data);
	// if (amount0 > 0) require(balance0Before.add(amount0) <= balance0(), 'M0');
	// if (amount1 > 0) require(balance1Before.add(amount1) <= balance1(), 'M1');

	// emit Mint(msg.sender, recipient, tickLower, tickUpper, amount, amount0, amount1);
}

//
func (p *Pool) Collect(
	recipient string,
	tickLower,
	tickUpper int,
	amount0Requested,
	amount1Requested *big.Int,
) (amount0, amount1 *big.Int) {
	// we don't need to checkTicks here, because invalid positions will never have non-zero tokensOwed{0,1}
	position, err = positions[fmt.Sprintf("%s%d%d", owner, tickLower, tickUpper)]
	if err != nil {
		panic()
	}

	amount0 = new(big.Int)
	if amount0Requested.Cmp(position.tokensOwed0) >= 1 {
		amount0 = position.tokensOwed0
	} else {
		amount0 = amount0Requested
	}
	amount1 = new(big.Int)
	if amount1Requested.Cmp(position.tokensOwed1) >= 1 {
		amount1 = position.tokensOwed1
	} else {
		amount1 = amount1Requested
	}

	if (amount0.Cmp(big.NewInt(0)) >= 1) {
		position.tokensOwed0 = new(big.Int).Sub(position.tokensOwed0, amount0)
		// TransferHelper.safeTransfer(token0, recipient, amount0);
	}

	if (amount1.Cmp(big.NewInt(0)) >= 1) {
		position.tokensOwed1 = new(big.Int).Sub(position.tokensOwed1, amount1)
		// TransferHelper.safeTransfer(token1, recipient, amount1);
	}
}

// 
func (p *Pool) Burn(tickLower, tickUpper int, amount *big.Int) (amount0, amount1 *big.Int) {
	position, amount0, amount1 := _modifyPosition(
		&ModifyPositionParams{
			owner: msg.sender,
			tickLower: tickLower,
			tickUpper: tickUpper,
			liquidityDelta: new(big.Int).Neg(amount),
		}
	)

	amount0 := new(big.Int).Neg(amount0)
	amount1 := new(big.Int).Neg(amount1)

	// if (amount0 > 0 || amount1 > 0) {
	if (amount0.Cmp(big.NewInt(0)) >= 1 || amount1.Cmp(big.NewInt(0)) >= 1) {
		position[tokensOwed0] = new(big.Int).Add(position[tokensOwed0], amount0)
		position[tokensOwed1] = new(big.Int).Add(position[tokensOwed1], amount1)
	}
}

struct SwapCache {
	// The protocol fee for the input token.
	FeeProtocol int
	// The liquidity at the beginning of the swap.
	LiquidityStart *big.Int
}

// The top level state of the swap, the results of which are recorded in storage
// at the end.
struct SwapState {
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

struct StepComputations {
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
func (p *Pool) Swap(
	recipient string,
	zeroForOne bool,
	amountSpecified,
	sqrtPriceLimitX96 *big.Int,
	// bytes calldata data
) (amount0, amount1 *big.Int) {
	// require(amountSpecified != 0, 'AS');
	if (amountSpecified.Cmp(big.NewInt(0)) == 0) {
		panic()
	}

	slot0Start = p.slot0;

	var cacheFeeProtocol int
	var stateFeeGrowthGlobalX128 *big.Int
	if (zeroForOne) {
		if !((sqrtPriceLimitX96.Cmp(slot0Start.sqrtPriceX96) >= -1) && (sqrtPriceLimitX96.Cmp(TickMath.MIN_SQRT_RATIO) >= 1)) {
			panic()
		}
		cacheFeeProtocol := slot0Start.feeProtocol % 16
		stateFeeGrowthGlobalX128 := slot0Start.feeGrowthGlobal0X128
	} else {		
		if !((sqrtPriceLimitX96.Cmp(slot0Start.sqrtPriceX96) >= 1) && (sqrtPriceLimitX96.Cmp(TickMath.MIN_SQRT_RATIO) <= 1)) {
			panic()
		}
		cacheFeeProtocol := slot0Start.feeProtocol >> 4
		stateFeeGrowthGlobalX128 := slot0Start.feeGrowthGlobal1X128
	}

	cache := SwapCache{
		liquidityStart: liquidity,
		feeProtocol: cacheFeeProtocol,
	}
 
	if (amountSpecified.Cmp(big.NewInt(0)) >= 1) {
		exactInput := true
	} else {
		exactInput := false
	}

	state := SwapState{
		amountSpecifiedRemaining: amountSpecified,
		amountCalculated: 0,
		sqrtPriceX96: slot0Start.sqrtPriceX96,
		tick: slot0Start.tick,
		feeGrowthGlobalX128: stateFeeGrowthGlobalX128,
		protocolFee: 0,
		liquidity: cache.liquidityStart
	};

	// Continue swapping as long as we haven't used the entire input/output and 
	// haven't reached the price limit
	for (state.amountSpecifiedRemaining.Cmp(big.NewInt(0)) >= 1 && state.sqrtPriceX96.Cmp(sqrtPriceLimitX96) >= 1) {
		// StepComputations memory step;
		step := new(StepComputations)

		step[SqrtPriceStartX96] = state.sqrtPriceX96;

		step.tickNext, step.initialized = nextInitializedTickWithinOneWord(
			state.tick,
			tickSpacing,
			zeroForOne
		)

		// Ensure that we do not overshoot the min/max tick, as the tick bitmap is not aware of these bounds
		if (step.tickNext < constants.MinTick) {
			step.tickNext = constants.MinTick
		} else if (step.tickNext > constants.MaxTick) {
			step.tickNext = constants.MaxTick
		}

		// Get the price for the next tick
		step.sqrtPriceNextX96 = tickMath.getSqrtRatioAtTick(step.tickNext);

		// Compute values to swap to the target tick, price limit, or point where input/output amount is exhausted
		var sqrtRatioTargetX96 *big.Int
		if (zeroForOne) {
			// step.sqrtPriceNextX96 < sqrtPriceLimitX96
			if (step.sqrtPriceNextX96.Cmp(sqrtPriceLimitX96) <= -1) {
				sqrtRatioTargetX96 = sqrtPriceLimitX96
			} else {
				sqrtRatioTargetX96 = step.sqrtPriceNextX96
			}
		} else {
			if (step.sqrtPriceNextX96.Cmp(sqrtPriceLimitX96) >= 1) {
				sqrtRatioTargetX96 = sqrtPriceLimitX96
			} else {
				sqrtRatioTargetX96 = step.sqrtPriceNextX96
			}
		}

		if (exactInput) {
			state.amountSpecifiedRemaining = new(big.Int).Sub(state.amountSpecifiedRemaining, new(big.Int).Add(step.amountIn step.feeAmount))
			state.amountCalculated = state.amountCalculated.Sub(step.amountOut)
		} else {
			state.amountSpecifiedRemaining = new(big.Int).Add(state.amountSpecifiedRemaining, step.amountOut)
			state.amountCalculated = new(big.Int)(state.amountCalculated, new(big.Int).Add(step.amountIn, step.feeAmount))
		}

		// If the protocol fee is on, calculate how much is owed, decrement 
		// feeAmount, and increment protocolFee
		if (cache.feeProtocol > 0) {
			delta := new(big.Int).Div(step.feeAmount, cache.feeProtocol)
			step.feeAmount = new(big.Int).Sub(step.feeAmount, delta)
			state.protocolFee = new(big.Int).Add(state.protocolFee, delta);
		}

		// Update global fee tracker
		// if (state.liquidity > 0)
		if (state.liquidity.Cmp(big.NewInt(0)) >= 1) {
			state.feeGrowthGlobalX128 = new(big.Int).Add(state.feeGrowthGlobalX128, FullMath.mulDiv(step.feeAmount, FixedPoint128.Q128, state.liquidity))
		}

		///////////////////////////////////////////////////////////////////////
		// Shift tick if we reached the next price
		if (state.sqrtPriceX96.Cmp(step.sqrtPriceNextX96) == 0) {
			// If the tick is initialized, run the tick transition
			if (step.initialized) {
				// Check for the placeholder value, which we replace with the 
				// actual value the first time the swap crosses an initialized 
				/// tick
				
				tempFeeGrowthGlobal0X128 := state.feeGrowthGlobalX128
				tempFeeGrowthGlobal1X128 := p.feeGrowthGlobalX128
				if (zeroForOne) {
					tempFeeGrowthGlobal0X128 = p.feeGrowthGlobalX128
					tempFeeGrowthGlobal1X128 = state.feeGrowthGlobalX128
					// if we're moving leftward, we interpret liquidityNet as 
					// the opposite sign
					liquidityNet = new(big.Int).Neg(liquidityNet)
				}

				liquidityNet := tick.Cross(
					step.tickNext,
					tempFeeGrowthGlobal0X128,
					tempFeeGrowthGlobal1X128,
				)

				state.liquidity = liquidityMath.AddDelta(state.liquidity, liquidityNet);
			}

			if (zeroForOne) {
				state.tick = step.tickNext - 1
			} else {
				state.tick = step.tickNext
			}
		} else if (state.sqrtPriceX96.Cmp(step.sqrtPriceStartX96) != 0) {
			// Recompute unless we're on a lower tick boundary (i.e. already 
			// transitioned ticks), and haven't moved
			state.tick = tickMath.GetTickAtSqrtRatio(state.sqrtPriceX96);
		}
	}

	// update tick and write an oracle entry if the tick change
	if (state.tick != slot0Start.tick) {
		
		slot0.sqrtPriceX96 = state.sqrtPriceX96
		slot0.tick = state.tick
	} else {
		// Otherwise just update the price
		slot0.sqrtPriceX96 = state.sqrtPriceX96;
	}

	// Update liquidity if it changed
	if (cache.liquidityStart.Cmp(state.liquidity) != 0) {
		liquidity = state.liquidity;
	}

	// Update fee growth global and, if necessary, protocol fees
	if (zeroForOne) {
		feeGrowthGlobal0X128 = state.feeGrowthGlobalX128
		if (state.protocolFee > 0) {
			protocolFees.token0 = new(big.Int).Add(protocolFees.token0, state.protocolFee)
		}
	} else {
		feeGrowthGlobal1X128 = state.feeGrowthGlobalX128
		if (state.protocolFee > 0) {
			protocolFees.token1 = new(big.Int).Add(protocolFees.token1, state.protocolFee)
		}
	}

	if (zeroForOne == exactInput) {
		amount0 = new(big.Int).Sub(amountSpecified, state.amountSpecifiedRemaining)
		amount1 = state.amountCalculated
	} else {
		amount0 = state.amountCalculated
		amount1 = new(big.Int).Sub(amountSpecified, state.amountSpecifiedRemaining)
	}

	// do the transfers and collect payment
	if (zeroForOne) {
		// if (amount1 < 0) TransferHelper.safeTransfer(token1, recipient, uint256(-amount1));

		balance0Before := p.balance0
		// IUniswapV3SwapCallback(msg.sender).uniswapV3SwapCallback(amount0, amount1, data);
		// require(balance0Before.add(uint256(amount0)) <= balance0(), 'IIA');
	} else {
		// if (amount0 < 0) TransferHelper.safeTransfer(token0, recipient, uint256(-amount0));

		balance1Before := p.balance1
		// IUniswapV3SwapCallback(msg.sender).uniswapV3SwapCallback(amount0, amount1, data);
		// require(balance1Before.add(uint256(amount1)) <= balance1(), 'IIA');
	}
}

// func Make (tick int, liquidity []*position.Position) *Pool {

// }

