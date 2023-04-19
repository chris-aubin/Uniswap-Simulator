# Strategies

Strategies are represented using the following struct
```
    type Strategy struct {
    	Address        string
    	Amount0        *big.Int
    	Amount1        *big.Int
    	GasUsed        *big.Int
    	GasAvs         *GasAvs
    	UpdateInterval int
    	Positions      []*StrategyPosition
    	Rebalance      func(p *pool.Pool, s *Strategy)
    }
```

- `Amount0` is the amount of `token0` that the strategy has available to provide liquidity. 
- `Amount1` is the amount of `token1` that the strategy has available to provide liquidity. 
- `GasUsed` is the amount of gas the strategy has used in GETH.
- `GasAvs` is the average cost of each pool operation in GETH.
- `UpdateInterval` is how often, in blocks, the `Rebalance` function should be called (assuming that every block contains at least one transaction). In the case that there are no transactions in a block, `Rebalance` will not be called until there is a new transaction, regardless of the `UpdateInterval`.
- `Positions` is a slice of the strategy's positions (for a given position the slice stores the `TickLower`, `TickUpper` (so that the position can be identified in the pool's position-indexed state) and the `Liquidity`).
- `Rebalance` is the function that mints or burns liquidity based upon the state of the pool. This is what distinguishes different strategies.



All strategies have a `BurnAll` function that burns all of the strategy's positions and calculates the tokens owed to the strategy, a `Results` function that returns the tokens that the strategy has accumulated and the total amount of gas that the strategy has spent and a 
`Make` function that initialises a strategy. 

The only field that differs significantly from strategy to strategy is the `Rebalance` function. The function is of type `func(p *pool.Pool, s *Strategy)`. It takes in a `Pool` and a  `Strategy`. It can call any of the `Pool` methods and it has access to all of the `Pool` and `Strategy` state. It make use of any number of helper functions. For example, the `Rebalance` function for a Uniswap v2 style strategy would look like:

```
    func V2StrategyRebalance(p *pool.Pool, s *Strategy) {
    	if len(s.Positions) == 0 {
    		V2StrategyMintPosition(p, s)
    	}
    	return
    }
    
    func V2StrategyMintPosition(p *pool.Pool, s *Strategy) {
    	s.GasUsed = new(big.Int).Add(s.GasUsed, s.GasAvs.MintGas)
    	tickLower := constants.MinTick
    	tickUpper := constants.MaxTick
    	sqrtRatioAX96 := tickMath.GetSqrtRatioAtTick(tickLower)
    	sqrtRatioBX96 := tickMath.GetSqrtRatioAtTick(tickUpper)
    
    	liquidity := liquidityAmounts.GetLiquidityForAmounts(p.Slot0.SqrtPriceX96, sqrtRatioAX96, sqrtRatioBX96, s.Amount0, s.Amount1)
    	 sqrtRatioAX96, sqrtRatioBX96, liquidity)
    	if liquidity.Cmp(big.NewInt(0)) <= 0 {
    		return
    	}
    	p.Mint(s.Address, tickLower, tickUpper, liquidity)
    	s.Positions = append(s.Positions, &StrategyPosition{
    		TickLower: tickLower,
    		TickUpper: tickUpper,
    		Liquidity: liquidity,
    	})
    }
```

This design makes it possible to create and test far more complicated, dynamic than the above `v2` strategy. Each strategy's rebalance function must be added to the `strategies` map in `strategy.go` before it can be used.