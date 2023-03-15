package pool

type Slot0 struct {
	// The current price
	SqrtPriceX96 int
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

// Accumulated protocol fees in token0/token1 units
type ProtocolFees struct {
	Token0 int
	Token1 int
}

type Pool struct {
	Slot0                Slot0
	FeeGrowthGlobal0X128 int
	FeeGrowthGlobal1X128 int
	ProtocolFees         ProtocolFees
	Liquidity            int
}

type modifyPositionParams struct {
	// the address that owns the position
	owner int
	// the lower and upper tick of the position
	tickLower int
	tickUpper int
	// any change in liquidity
	liquidityDelta int
}

func (pool *Pool) modifyPosition(params modifyPositionParams) {
	// TODO
}

// TODO
//
func (pool *Pool) Mint(recipient int, tickLower int, tickUpper int, amount int) (int, int) {

}

// TODO
//
func (pool *Pool) Burn() {

}

// TODO
//
func (pool *Pool) Swap() {

}

// TODO
//
func (pool *Pool) Collect() {

}

func Make() *Pool {

}
