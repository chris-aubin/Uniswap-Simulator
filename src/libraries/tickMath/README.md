# TickMath contract

## `getSqrtRationAtTick`
The `getSqrtRationAtTick` function accepts a tick and calculates $\sqrt{1.0001^{tick}} \cdot 2^{96}$.  The $2^{96}$ is necessary because square root prices are represented as Q64.96 fixed point numbers. This means that, for some fixed point number $x$, $x = \frac{y}{2^{96}}$ is calculated and $x$ is represented using $y$. However, rather than calculating this value through exponentiation, `getSqrtRationAtTick` performs a complex combination of bit manipulation and multiplication by hardcoded numbers. Furthermore, Uniswap provides no explanation of what those hardcoded numbers represent or how the function uses them. In fact, at the time of writing, there didn't seem to be any explanations of `getSqrtRationAtTick` published online.

`getSqrtRationAtTick` processes that tick one bit at a time (`0x1 = b1`, `0x2 = b10`, `0x4 = b100`, etc. so each if statement is checking whether a bit in tick is `1`. Each of hard coded hex values is a pre-computed power of $\sqrt{\frac{1}{1.0001}}$:
$$
    \texttt{0xfffcb933bd6fad37aa2d162d1a594001} = \sqrt{\frac{1}{1.0001}}^1 \\
    \texttt{0xfff97272373d413259a46990580e213a} = \sqrt{\frac{1}{1.0001}}^2 \\
    \texttt{0xfff2e50f5f656932ef12357cf3c7fdcc} = \sqrt{\frac{1}{1.0001}}^4 \\
    \texttt{0xfff2e50f5f656932ef12357cf3c7fdcc} = \sqrt{\frac{1}{1.0001}}^8 \\
    \cdots
$$
This can be confirmed trivially in most programming languages, for example, in Python one could compare:
```
    int("0xfff2e50f5f656932ef12357cf3c7fdcc", 16)/2**128
```
to:
```
    math.sqrt((1/1.0001))**8
```
These pre-computed values are used to compute $\sqrt{\frac{1}{1.0001}^{tick}}$ using the fact that $x^c = x^a*x^b$ where $a + b = c$ as follows: \\ 
$$
    tick = 2^0 \cdot x_0 \in \{0,1\} \; + \\
    2^1 \cdot x_1 \in \{0,1\} \; + \\
    2^2 \cdot x_2 \in \{0,1\} \; + \\
    2^3 \cdot x_3 \in \{0,1\} \; + \\ 
    2^4 \cdot x_4 \in \{0,1\} \; + \\
    2^5 \cdot x_5 \in \{0,1\} \; + \\
    \cdots
$$
Therefore:
$$
    \sqrt{(1/1.0001)^{tick}} = \sqrt{(1/1.0001)^{2^0}} \cdot x_0 \in \{0,1\} \cdot \\
    \sqrt{(1/1.0001)^{2^1}} \cdot x_1 \in \{0,1\} \; \cdot \\
    \sqrt{(1/1.0001)^{2^2}} \cdot x_2 \in \{0,1\} \; \cdot \\
    \sqrt{(1/1.0001)^{2^3}} \cdot x_3 \in \{0,1\} \; \cdot \\
    \sqrt{(1/1.0001)^{2^4}} \cdot x_4 \in \{0,1\} \; \cdot \\
    \sqrt{(1/1.0001)^{2^5}} \cdot x_5 \in \{0,1\} \; \cdot \\
    \cdots
$$

After $\sqrt{\frac{1}{1.0001}^{tick}}$ has been calculated the function merely needs to calculate the reciprocal in the case that `tick` $>0$ and then convert the value from Q128.128 to Q64.96. This approach is used in the name of gas-optimisation, to avoid the repeated multiplication that exponentiation typically involves.

Go supports exponentiation, and the simulator code is not concerned with gas efficiency, so it could simply calculate $\sqrt{1.0001^{tick}} \cdot 2^{96}$. However, in the name of simulating the smart contracts as accurately as possible, the `getSqrtRationAtTick` function was implemented in the simulator as it was in the smart contracts.

## `getTickAtSqrtRatio`
The `getTickAtSqrtRatio` function accepts `sqrtPriceX96`, the square root price for which to compute the tick, and returns the greatest tick for which the ratio is less than or equal to the input ratio. This means that it calculates $\log_{\sqrt{1.0001}}\texttt{sqrtPriceX96}$

The maths behind the approximation algorithm used by `getTickAtSqrtRatio` is explained in a 2 year old blog post [here](https://hackmd.io/@abdk/SkVJeHK9v). First, $\log_2\texttt{sqrtPriceX96}$ is approximated using the following iterative algorithm:

- Start with an initial approximation $l_0(x) = \left\lfloor\log _2 x\right\rfloor$
- Define $l_{i+1}(x)$ as follows:
        $$
            l_{i+1}(x)=\left\lfloor\log _2 x\right\rfloor+\frac{1}{2} l_i\left(\left\lfloor\left(\left\lfloor\frac{x}{2^{\left\lfloor\log _2 x\right\rfloor}}\right\rfloor_{2^{-127}}\right)^2\right]_{2^{-127}}\right)
        $$
        Where $\left\lfloor a \right\rfloor_b$ is $a$ rounded down to the closest factor of $b$. 


Notice that the most significant bit is equal to $\left\lfloor\log _2 x\right\rfloor$. `getTickAtSqrtRatio` starts by finding the most significant bit of `sqrtPriceX96`. `getTickAtSqrtRatio` then proceeds to refine that approximation according to the above formula. It does this by calculating one fixed point place at a time. Consider the following code from the `tickMath` contract (or alternative lines 163-167 `tickMath.go `):

```
    if (msb >= 128) r = ratio >> (msb - 127);
    else r = ratio << (127 - msb);
```

Notice that this is setting: 
$$
\texttt{r} = \left\lfloor\frac{x}{2^{\left\lfloor\log _2 x\right\rfloor}}\right\rfloor_{2^{-127}}
$$

Consider the first of the inline assembly blocks in the `tickMath` contract (or alternatively the second for loop in `tickMath.go`):

```
   assembly {
       r := shr(127, mul(r, r))
       let f := shr(128, r)
       log_2 := or(log_2, shl(63, f))
       r := shr(f, r)
   }
```

- `r := shr(127, mul(r, r))` calculates the square of `r` and then shifts the result 127 bits to the right to divide it by $2^{127}$, equivalent to $\left(\left\lfloor\frac{x}{2^{\left\lfloor\log _2 x\right\rfloor}}\right\rfloor_{2^{-127}}\right)^2$
- `let f := shr(128, r)` shifts the value of `r` 128 bits to the right, extracting the most significant bit of `r`, and assigns the result to a temporary variable `f`.
- `log_2 := or(log_2, shl(n, f))` performs a bitwise OR operation on the current value of `log_2` and the result of shifting `f` `n` bits to the left. Because `f` is the most significant bit of `r`, this effectively sets the `n`th bit in the binary representation of `log_2` to the most sifnificant bit of `r`. This is done iteratively with decreasing values of `n` in subsequent assembly blocks, starting at `63`, which corresponds to the first bit after the fixed point in a Q64.96 number. Thus, the inline assembly blocks are each used to approximate one bit after the fixed point in `log_2`.
- `r := shr(f, r)` shifts the value of `r` `f` bits to the right, dividing `r` by $2^f$. This step is performed to prepare `r` for the next iteration of the loop.

After $\log_2 \texttt{sqrtPriceX96}$ has been approximated, the change of base rule can be used to approximate $\log_{\sqrt{1.0001}} \texttt{sqrtPriceX96}$:
$$
    \log_{\sqrt{1.0001}}x = \frac{\log_2 x}{\log_2\sqrt{1.0001}}
$$
And:
$$
    \log_2{\sqrt{1.0001}} = \frac{\log_{\sqrt{1.0001}}\sqrt{1.0001}}{\log_{\sqrt{1.0001}}2} \\ 
    = \frac{1}{\log_{\sqrt{1.0001}}2}
$$
So, substituting that into the change of base rule:
$$
    \log_{\sqrt{1.0001}}x = \frac{\log_2 x}{\log_2\sqrt{1.0001}} \\
    = \frac{\log_2 x}{\frac{1}{\log_{\sqrt{1.0001}}2}} \\
    = \log_2 x \cdot \log_{\sqrt{1.0001}}2
$$
Notice, $\log_{\sqrt{1.0001}}2$ is a constant. It is approximately:
$$
1.33058... \cdot 10^{15}
$$
Which, in Q64.96 is:
$$
1.33058... \cdot 10^{15} \cdot 2^{64} \approx 2.55739... \cdot 10^{23} \approx 255738958999603826347141
$$

This change of base is performed by the following line the `tickMath` contract:

```
    int256 log_sqrt10001 = log_2 * 255738958999603826347141;
```

The last part of `getTickAtSqrtRatio` uses the absolute error of the approximation to ensure that the function returns the greatest tick for which the ratio is less than or equal to the input ratio.