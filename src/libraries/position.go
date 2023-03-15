package position

type Info struct {
	// The amount of liquidity owned by this position
	Liquidity int
	// Fee growth per unit of liquidity as of the last update to liquidity or fees owed
	FeeGrowthInside0LastX128 int
	FeeGrowthInside1LastX128 int
	// The fees owed to the position owner in token0/token1
	TokensOwed0 int
	TokensOwed1 int
}

func Get() {

}

func Update() {

}

func Make() {

}
