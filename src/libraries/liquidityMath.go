package liquidityMath

import "math/big"

func AddDelta(x, y *big.Int) *big.Int {
	if y.Cmp(big.NewInt(0)) == -1 {
		result := new(big.Int).Add(x, y)
		if result.Cmp(x) != -1 {
			panic("liquidityMath.AddDelta: LS")
		}
		return result
	} else {
		result := new(big.Int).Add(x, y)
		if result.Cmp(x) == -1 {
			panic("liquidityMath.AddDelta: LA")
		}
		return result
	}
}
