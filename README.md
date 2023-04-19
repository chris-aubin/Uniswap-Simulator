# Uniswap v3 Simulator
From the `src` folder where `path_to_simulation_data` is the path to a folder containing, `transactions.txt`, a file that contains the transactions for the backtesting period (i.e. the output of the `transactions.py` script in [chris-aubin/Uniswap-Simulator-Data-Collection](https://github.com/chris-aubin/Uniswap-Simulator-Data-Collection), `pool.txt`, a file tha contains the initial pool state (i.e. the output of the `pool_state.py`  script in [chris-aubin/Uniswap-Simulator-Data-Collection](https://github.com/chris-aubin/Uniswap-Simulator-Data-Collection) run at the date corresponding to the start of the backtesting period), `gas.txt`, a file that contains the average gas fees for each pool method called during the backtesting period (i.e. the output of the `gas_estimates.py` script run on the above transactions in [chris-aubin/Uniswap-Simulator-Data-Collection](https://github.com/chris-aubin/Uniswap-Simulator-Data-Collection)) and `strategy.txt`, a file that contains the details for the strategy that is to be tested. An example of such a file is shown below:

```
{
    "strategy": "v2",
    "amount0": 33,
    "amount1": 480000000000000,
    "updateInterval": 1
}
```

This indicates that a v2 style strategy should be tested, that it should be allocated `33` satoshis and `480000000000000` GETH, and that rebalance should be called in every block that the pool state changes.