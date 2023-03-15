package tick

type Info struct {
	// The total position liquidity that references this tick
	LiquidityGross int
	// Amount of net liquidity added (subtracted) when tick is crossed from left to right (right to left),
	LiquidityNet int
	// Fee growth per unit of liquidity on the _other_ side of this tick (relative to the current tick)
	// Only has relative meaning, not absolute — the value depends on when the tick is initialized
	FeeGrowthOutside0X128 int
	FeeGrowthOutside1X128 int
	// The cumulative tick value on the other side of the tick
	TickCumulativeOutside int
	// The seconds per unit of liquidity on the _other_ side of this tick (relative to the current tick)
	// Only has relative meaning, not absolute — the value depends on when the tick is initialized
	SecondsPerLiquidityOutsideX128 int
	// The seconds spent on the other side of the tick (relative to the current tick)
	// Only has relative meaning, not absolute — the value depends on when the tick is initialized
	SecondsOutside int
	// True iff the tick is initialized, i.e. the value is exactly equivalent to the expression liquidityGross != 0
	// These 8 bits are set to prevent fresh sstores when crossing newly initialized ticks
	Initialized bool
}
