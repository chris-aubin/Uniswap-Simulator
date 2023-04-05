package oracle
// Provides price and liquidity data useful for a wide variety of system designs 
// (in this simulator the data is useful for strategies). Instances of stored 
// oracle data, "observations", are collected in the oracle array. Every pool 
// is initialized with an oracle array length of 1. In production (i.e. when 
// using the deployed Uniswap v3 smart contracts), anyone can pay to increase 
// the maximum length of the oracle array (allowing the oracle to store more 
// observations, which are accessible to all). New slots will be added when the 
// array is fully populated. Observations are overwritten when the length of the
// oracle array is populated. The most recent observation is available, 
// independent  of the length of the oracle array, by passing 0 to observe().

type Oracle struct {
	// A map from the time of an observation to the observation itself.
	ObservationsMap   [int]Observation
	// A list of observations, sorted by time.
	ObservationsSlice []Observation
}

type Observation struct {
	// The block timestamp of the observation.
	blockTimestamp int
	// The tick accumulator, i.e. tick * time elapsed since the pool was first 
	// initialized.
	tickCumulative int
	// The seconds per liquidity, i.e. seconds elapsed / max(1, liquidity) since
	// the pool was first initialized.
	secondsPerLiquidityCumulativeX128 *big.Int
	// Whether or not the observation is initialized.
	initialized bool
}

/// @notice Transforms a previous observation into a new observation, given the passage of time and the current tick and liquidity values
/// @dev blockTimestamp _must_ be chronologically equal to or greater than last.blockTimestamp, safe for 0 or 1 overflows
/// @param last The specified observation to be transformed
/// @param blockTimestamp The timestamp of the new observation
/// @param tick The active tick at the time of the new observation
/// @param liquidity The total in-range liquidity at the time of the new observation
/// @return Observation The newly populated observation
func (o *Oracle) transform(
	blockTimestamp,
	tick int,
	liquidity *big.Int,
) *Observation {
	last := o.ObservationsSlice[len(o.ObservationsSlice) - 1]
	delta = blockTimestamp - last.blockTimestamp;

	// Calculate seconds per liquidity
	var secondsPerLiquidityCumulativeX128 = big.NewInt(0)
	deltaShifted := big.NewInt(0).Lsh(big.NewInt(delta), 128)
	if liquidity.Cmp(big.NewInt(0)) > 0 {
		deltaShifted.Div(deltaShifted, liquidity)	
	}
	secondsPerLiquidityCumulativeX128.Add(last.secondsPerLiquidityCumulativeX128, deltaShifted)

	return
		&Observation{
			blockTimestamp: blockTimestamp,
			tickCumulative: last.tickCumulative + tick * delta,
			secondsPerLiquidityCumulativeX128: secondsPerLiquidityCumulativeX128,
			initialized: true,
		};
}

/// @notice Fetches the observations beforeOrAt and atOrAfter a given target, i.e. where [beforeOrAt, atOrAfter] is satisfied
/// @dev Assumes there is at least 1 initialized observation.
/// Used by observeSingle() to compute the counterfactual accumulator values as of a given block timestamp.
/// @param self The stored oracle array
/// @param time The current block.timestamp
/// @param target The timestamp at which the reserved observation should be for
/// @param tick The active tick at the time of the returned or simulated observation
/// @param index The index of the observation that was most recently written to the observations array
/// @param liquidity The total pool liquidity at the time of the call
/// @param cardinality The number of populated elements in the oracle array
/// @return beforeOrAt The observation which occurred at, or before, the given timestamp
/// @return atOrAfter The observation which occurred at, or after, the given timestamp
func (o *Oracle) getSurroundingObservations(
	time,
	target,
	tick int,
	liquidity *big.Int,
) (beforeOrAt, atOrAfter *Observation){

}

/// @dev Reverts if an observation at or before the desired observation timestamp does not exist.
/// 0 may be passed as `secondsAgo' to return the current cumulative values.
/// If called with a timestamp falling between two observations, returns the counterfactual accumulator values
/// at exactly the timestamp between the two observations.
/// @param self The stored oracle array
/// @param time The current block timestamp
/// @param secondsAgo The amount of time to look back, in seconds, at which point to return an observation
/// @param tick The current tick
/// @param index The index of the observation that was most recently written to the observations array
/// @param liquidity The current in-range pool liquidity
/// @param cardinality The number of populated elements in the oracle array
/// @return tickCumulative The tick * time elapsed since the pool was first initialized, as of `secondsAgo`
/// @return secondsPerLiquidityCumulativeX128 The time elapsed / max(1, liquidity) since the pool was first initialized, as of `secondsAgo`
func (o *Oracle) observeSingle(
	time,
	secondsAgo,
	tick int,
	liquidity *big.Int,
) (tickCumulative int, secondsPerLiquidityCumulativeX128 *big.Int) {
	if (secondsAgo == 0) {
		// Get the most recent observation
		last := o.ObservationsSlice[len(o.ObservationsSlice)-1]
		if (last.blockTimestamp != time) {
			last = transform(last, time, tick, liquidity)
		}
		o.ObservationsMap[time] = last
		o.ObservationsSlice = append(o.ObservationsSlice, last)
		return last.tickCumulative, last.secondsPerLiquidityCumulativeX128;
	}
	target := time - secondsAgo;

	////////////
	(Observation memory beforeOrAt, Observation memory atOrAfter) = getSurroundingObservations(
		self,
		time,
		target,
		tick,
		index,
		liquidity,
		cardinality
	);

	if (target == beforeOrAt.blockTimestamp) {
		// we're at the left boundary
		return (beforeOrAt.tickCumulative, beforeOrAt.secondsPerLiquidityCumulativeX128);
	} else if (target == atOrAfter.blockTimestamp) {
		// we're at the right boundary
		return (atOrAfter.tickCumulative, atOrAfter.secondsPerLiquidityCumulativeX128);
	} else {
		// we're in the middle
		uint32 observationTimeDelta = atOrAfter.blockTimestamp - beforeOrAt.blockTimestamp;
		uint32 targetDelta = target - beforeOrAt.blockTimestamp;
		return (
			beforeOrAt.tickCumulative +
				((atOrAfter.tickCumulative - beforeOrAt.tickCumulative) / observationTimeDelta) *
				targetDelta,
			beforeOrAt.secondsPerLiquidityCumulativeX128 +
				uint160(
					(uint256(
						atOrAfter.secondsPerLiquidityCumulativeX128 - beforeOrAt.secondsPerLiquidityCumulativeX128
					) * targetDelta) / observationTimeDelta
				)
		);
	}
}

/// @notice Returns the accumulator values as of each time seconds ago from the given time in the array of `secondsAgos`
/// @dev Reverts if `secondsAgos` > oldest observation
/// @param self The stored oracle array
/// @param time The current block.timestamp
/// @param secondsAgos Each amount of time to look back, in seconds, at which point to return an observation
/// @param tick The current tick
/// @param index The index of the observation that was most recently written to the observations array
/// @param liquidity The current in-range pool liquidity
/// @param cardinality The number of populated elements in the oracle array
/// @return tickCumulatives The tick * time elapsed since the pool was first initialized, as of each `secondsAgo`
/// @return secondsPerLiquidityCumulativeX128s The cumulative seconds / max(1, liquidity) since the pool was first initialized, as of each `secondsAgo`
function observe(
	Observation[65535] storage self,
	uint32 time,
	uint32[] memory secondsAgos,
	int24 tick,
	uint16 index,
	uint128 liquidity,
	uint16 cardinality
) internal view returns (int56[] memory tickCumulatives, uint160[] memory secondsPerLiquidityCumulativeX128s) {
	require(cardinality > 0, 'I');

	tickCumulatives = new int56[](secondsAgos.length);
	secondsPerLiquidityCumulativeX128s = new uint160[](secondsAgos.length);
	for (uint256 i = 0; i < secondsAgos.length; i++) {
		(tickCumulatives[i], secondsPerLiquidityCumulativeX128s[i]) = observeSingle(
			self,
			time,
			secondsAgos[i],
			tick,
			index,
			liquidity,
			cardinality
		);
	}
}

func Make (time int) *Oracle {
	observations := make([]Observation, 1)
	observations[0] = Observation{
		blockTimestamp: time,
		tickCumulative: 0,
		secondsPerLiquidityCumulativeX128: big.NewInt(0),
		initialized: true,
	}
	return &Oracle{
		Observations: observations,
	}
}