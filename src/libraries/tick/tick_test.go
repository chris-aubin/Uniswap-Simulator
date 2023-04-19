package tick

import (
	"fmt"
	"math/big"
	"testing"

	"github.com/chris-aubin/Uniswap-Simulator/src/libraries/constants"
)

var (
	testTicks *Ticks
)

func init() {
	testTicks = &Ticks{
		TickData: make(map[int]*Tick),
	}
}

// Test the GetFeeGrowthInside function
func TestGetFeeGrowthInside1(t *testing.T) {
	fmt.Println("Returns all for two uninitialized ticks if tick is inside")
	feeGrowthInside0X128, feeGrowthInside1X128 := testTicks.GetFeeGrowthInside(-2, 2, 0, big.NewInt(15), big.NewInt(15))
	if feeGrowthInside0X128.Cmp(big.NewInt(15)) != 0 {
		t.Errorf("Expected 15 for feeGrowthInside0X128, got %v", feeGrowthInside0X128)
	}
	if feeGrowthInside1X128.Cmp(big.NewInt(15)) != 0 {
		t.Errorf("Expected 15 for feeGrowthInside1X128, got %v", feeGrowthInside1X128)
	}
}

func TestGetFeeGrowthInside2(t *testing.T) {
	fmt.Println("Returns 0 for two uninitialized ticks if tick is above")
	feeGrowthInside0X128, feeGrowthInside1X128 := testTicks.GetFeeGrowthInside(-2, 2, 4, big.NewInt(15), big.NewInt(15))
	if feeGrowthInside0X128.Cmp(big.NewInt(0)) != 0 {
		t.Errorf("Expected 0 for feeGrowthInside0X128, got %v", feeGrowthInside0X128)
	}
	if feeGrowthInside1X128.Cmp(big.NewInt(0)) != 0 {
		t.Errorf("Expected 0 for feeGrowthInside1X128, got %v", feeGrowthInside1X128)
	}

}

func TestGetFeeGrowthInside3(t *testing.T) {
	fmt.Println("Returns 0 for two uninitialized ticks if tick is below")
	feeGrowthInside0X128, feeGrowthInside1X128 := testTicks.GetFeeGrowthInside(-2, 2, -4, big.NewInt(15), big.NewInt(15))
	if feeGrowthInside0X128.Cmp(big.NewInt(0)) != 0 {
		t.Errorf("Expected 0 for feeGrowthInside0X128, got %v", feeGrowthInside0X128)
	}
	if feeGrowthInside1X128.Cmp(big.NewInt(0)) != 0 {
		t.Errorf("Expected 0 for feeGrowthInside1X128, got %v", feeGrowthInside1X128)
	}

}

func TestGetFeeGrowthInside4(t *testing.T) {
	testTicks.TickData[2] = &Tick{
		FeeGrowthOutside0X128: big.NewInt(2),
		FeeGrowthOutside1X128: big.NewInt(3),
		LiquidityGross:        big.NewInt(0),
		LiquidityNet:          big.NewInt(0),
		Initialized:           true,
	}
	fmt.Println("Subtracts upper tick if below")
	feeGrowthInside0X128, feeGrowthInside1X128 := testTicks.GetFeeGrowthInside(-2, 2, 0, big.NewInt(15), big.NewInt(15))
	if feeGrowthInside0X128.Cmp(big.NewInt(13)) != 0 {
		t.Errorf("Expected 13 for feeGrowthInside0X128, got %v", feeGrowthInside0X128)
	}
	if feeGrowthInside1X128.Cmp(big.NewInt(12)) != 0 {
		t.Errorf("Expected 12 for feeGrowthInside1X128, got %v", feeGrowthInside1X128)
	}
	testTicks.Clear(2)
}

func TestGetFeeGrowthInside5(t *testing.T) {
	testTicks.TickData[-2] = &Tick{
		FeeGrowthOutside0X128: big.NewInt(2),
		FeeGrowthOutside1X128: big.NewInt(3),
		LiquidityGross:        big.NewInt(0),
		LiquidityNet:          big.NewInt(0),
		Initialized:           true,
	}
	fmt.Println("Subtracts lower tick if above")
	feeGrowthInside0X128, feeGrowthInside1X128 := testTicks.GetFeeGrowthInside(-2, 2, 0, big.NewInt(15), big.NewInt(15))
	if feeGrowthInside0X128.Cmp(big.NewInt(13)) != 0 {
		t.Errorf("Expected 13 for feeGrowthInside0X128, got %v", feeGrowthInside0X128)
	}
	if feeGrowthInside1X128.Cmp(big.NewInt(12)) != 0 {
		t.Errorf("Expected 12 for feeGrowthInside1X128, got %v", feeGrowthInside1X128)
	}
	testTicks.Clear(-2)
}

func TestGetFeeGrowthInside6(t *testing.T) {
	testTicks.TickData[2] = &Tick{
		FeeGrowthOutside0X128: big.NewInt(4),
		FeeGrowthOutside1X128: big.NewInt(1),
		LiquidityGross:        big.NewInt(0),
		LiquidityNet:          big.NewInt(0),
		Initialized:           true,
	}
	testTicks.TickData[-2] = &Tick{
		FeeGrowthOutside0X128: big.NewInt(2),
		FeeGrowthOutside1X128: big.NewInt(3),
		LiquidityGross:        big.NewInt(0),
		LiquidityNet:          big.NewInt(0),
		Initialized:           true,
	}
	fmt.Println("Subtracts upper and lower tick if inside")
	feeGrowthInside0X128, feeGrowthInside1X128 := testTicks.GetFeeGrowthInside(-2, 2, 0, big.NewInt(15), big.NewInt(15))
	if feeGrowthInside0X128.Cmp(big.NewInt(9)) != 0 {
		t.Errorf("Expected 9 for feeGrowthInside0X128, got %v", feeGrowthInside0X128)
	}
	if feeGrowthInside1X128.Cmp(big.NewInt(11)) != 0 {
		t.Errorf("Expected 11 for feeGrowthInside1X128, got %v", feeGrowthInside1X128)
	}
	testTicks.Clear(2)
	testTicks.Clear(-2)
}

func TestGetFeeGrowthInside7(t *testing.T) {
	testTicks.TickData[2] = &Tick{
		FeeGrowthOutside0X128: big.NewInt(3),
		FeeGrowthOutside1X128: big.NewInt(5),
		LiquidityGross:        big.NewInt(0),
		LiquidityNet:          big.NewInt(0),
		Initialized:           true,
	}
	testTicks.TickData[-2] = &Tick{
		FeeGrowthOutside0X128: new(big.Int).Sub(constants.MaxUint256, big.NewInt(3)),
		FeeGrowthOutside1X128: new(big.Int).Sub(constants.MaxUint256, big.NewInt(2)),
		LiquidityGross:        big.NewInt(0),
		LiquidityNet:          big.NewInt(0),
		Initialized:           true,
	}
	fmt.Println("Works correctly with overflow on inside tick")
	feeGrowthInside0X128, feeGrowthInside1X128 := testTicks.GetFeeGrowthInside(-2, 2, 0, big.NewInt(15), big.NewInt(15))
	if feeGrowthInside0X128.Cmp(big.NewInt(16)) != 0 {
		t.Errorf("Expected 16 for feeGrowthInside0X128, got %v", feeGrowthInside0X128)
	}
	if feeGrowthInside1X128.Cmp(big.NewInt(13)) != 0 {
		t.Errorf("Expected 13 for feeGrowthInside1X128, got %v", feeGrowthInside1X128)
	}
	testTicks.Clear(2)
	testTicks.Clear(-2)
}

func TestUpdate1(t *testing.T) {
	fmt.Println("Flips from zero to nonzero")
	flipped := testTicks.Update(0, 0, big.NewInt(1), big.NewInt(0), big.NewInt(0), big.NewInt(3), false)
	if !flipped {
		t.Errorf("Expected flipped to be true")
	}
	testTicks.Clear(0)
}

func TestUpdate2(t *testing.T) {
	fmt.Println("Does not flip from nonzero to greater nonzero")
	testTicks.Update(0, 0, big.NewInt(1), big.NewInt(0), big.NewInt(0), big.NewInt(3), false)
	flipped := testTicks.Update(0, 0, big.NewInt(1), big.NewInt(0), big.NewInt(0), big.NewInt(3), false)
	if flipped {
		t.Errorf("Expected flipped to be false")
	}
	testTicks.Clear(0)
}

func TestUpdate3(t *testing.T) {
	fmt.Println("Flips from nonzero to zero")
	testTicks.Update(0, 0, big.NewInt(1), big.NewInt(0), big.NewInt(0), big.NewInt(3), false)
	flipped := testTicks.Update(0, 0, big.NewInt(-1), big.NewInt(0), big.NewInt(0), big.NewInt(3), false)
	if !flipped {
		t.Errorf("Expected flipped to be true")
	}
	testTicks.Clear(0)
}

func TestUpdate4(t *testing.T) {
	fmt.Println("Does not flip from nonzero to lesser nonzero")
	testTicks.Update(0, 0, big.NewInt(2), big.NewInt(0), big.NewInt(0), big.NewInt(3), false)
	flipped := testTicks.Update(0, 0, big.NewInt(-1), big.NewInt(0), big.NewInt(0), big.NewInt(3), false)
	if flipped {
		t.Errorf("Expected flipped to be false")
	}
	testTicks.Clear(0)
}

func TestUpdate5(t *testing.T) {
	fmt.Println("Does not flip from nonzero to lesser nonzero")
	testTicks.Update(0, 0, big.NewInt(2), big.NewInt(0), big.NewInt(0), big.NewInt(3), false)
	flipped := testTicks.Update(0, 0, big.NewInt(-1), big.NewInt(0), big.NewInt(0), big.NewInt(3), false)
	if flipped {
		t.Errorf("Expected flipped to be false")
	}
	testTicks.Clear(0)
}

func TestUpdate6(t *testing.T) {
	fmt.Println("Reverts if total liquidity gross is greater than max")
	testTicks.Update(0, 0, big.NewInt(2), big.NewInt(0), big.NewInt(0), big.NewInt(3), false)
	testTicks.Update(0, 0, big.NewInt(1), big.NewInt(0), big.NewInt(0), big.NewInt(3), true)
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("Update did not panic when total liquidity gross was greater than max.")
		}
		testTicks.Clear(0)
	}()

	testTicks.Update(0, 0, big.NewInt(1), big.NewInt(0), big.NewInt(0), big.NewInt(3), false)
}

func TestUpdate7(t *testing.T) {
	fmt.Println("Nets the liquidity based on upper flag")
	testTicks.Update(0, 0, big.NewInt(2), big.NewInt(0), big.NewInt(0), big.NewInt(10), false)
	testTicks.Update(0, 0, big.NewInt(1), big.NewInt(0), big.NewInt(0), big.NewInt(10), true)
	testTicks.Update(0, 0, big.NewInt(3), big.NewInt(0), big.NewInt(0), big.NewInt(10), true)
	testTicks.Update(0, 0, big.NewInt(1), big.NewInt(0), big.NewInt(0), big.NewInt(10), false)
	tick := testTicks.Get(0)
	liquidityGross := tick.LiquidityGross
	liquidityNet := tick.LiquidityNet
	if liquidityGross.Cmp(big.NewInt(2+1+3+1)) != 0 {
		t.Errorf("Expected liquidityGross to be 7, got %v", liquidityGross)
	}
	if liquidityNet.Cmp(big.NewInt(2-1-3+1)) != 0 {
		t.Errorf("Expected liquidityNet to be -1, got %v", liquidityNet)
	}
	testTicks.Clear(0)
}

func TestUpdate8(t *testing.T) {
	fmt.Println("Reverts on overflow liquidity gross")
	testTicks.Update(0, 0, new(big.Int).Div(constants.MaxUint128, big.NewInt(2)), big.NewInt(0), big.NewInt(0), constants.MaxUint128, false)
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("Update did not panic when liquidity gross overflowed.")
		}
		testTicks.Clear(0)
	}()

	testTicks.Update(0, 0, new(big.Int).Add(new(big.Int).Add(constants.MaxUint128, big.NewInt(2)), big.NewInt(1)), big.NewInt(0), big.NewInt(0), constants.MaxUint128, false)
}

func TestUpdate9(t *testing.T) {
	fmt.Println("Assumes all growth happens below ticks lte current tick")
	testTicks.Update(1, 1, big.NewInt(1), big.NewInt(1), big.NewInt(2), constants.MaxUint128, false)
	tick := testTicks.Get(1)
	feeGrowthOutside0X128 := tick.FeeGrowthOutside0X128
	feeGrowthOutside1X128 := tick.FeeGrowthOutside1X128
	initialized := tick.Initialized
	if feeGrowthOutside0X128.Cmp(big.NewInt(1)) != 0 {
		t.Errorf("Expected feeGrowthOutside0X128 to be 1, got %v", feeGrowthOutside0X128)
	}
	if feeGrowthOutside1X128.Cmp(big.NewInt(2)) != 0 {
		t.Errorf("Expected feeGrowthOutside1X128 to be 2, got %v", feeGrowthOutside1X128)
	}
	if initialized != true {
		t.Errorf("Expected initialized to be true, got %v", initialized)
	}
	testTicks.Clear(1)
}

func TestUpdate10(t *testing.T) {
	fmt.Println("Does not set any growth fields if tick is already initialized")
	testTicks.Update(1, 1, big.NewInt(1), big.NewInt(1), big.NewInt(2), constants.MaxUint128, false)
	testTicks.Update(1, 1, big.NewInt(1), big.NewInt(6), big.NewInt(7), constants.MaxUint128, false)
	tick := testTicks.Get(1)
	feeGrowthOutside0X128 := tick.FeeGrowthOutside0X128
	feeGrowthOutside1X128 := tick.FeeGrowthOutside1X128
	initialized := tick.Initialized
	if feeGrowthOutside0X128.Cmp(big.NewInt(1)) != 0 {
		t.Errorf("Expected feeGrowthOutside0X128 to be 1, got %v", feeGrowthOutside0X128)
	}
	if feeGrowthOutside1X128.Cmp(big.NewInt(2)) != 0 {
		t.Errorf("Expected feeGrowthOutside1X128 to be 2, got %v", feeGrowthOutside1X128)
	}
	if initialized != true {
		t.Errorf("Expected initialized to be true, got %v", initialized)
	}
	testTicks.Clear(1)
}

func TestUpdate11(t *testing.T) {
	fmt.Println("Does not set any growth fields for ticks gt current tick")
	testTicks.Update(2, 1, big.NewInt(1), big.NewInt(1), big.NewInt(2), constants.MaxUint128, false)
	tick := testTicks.Get(2)
	feeGrowthOutside0X128 := tick.FeeGrowthOutside0X128
	feeGrowthOutside1X128 := tick.FeeGrowthOutside1X128
	initialized := tick.Initialized
	if feeGrowthOutside0X128.Cmp(big.NewInt(0)) != 0 {
		t.Errorf("Expected feeGrowthOutside0X128 to be 0, got %v", feeGrowthOutside0X128)
	}
	if feeGrowthOutside1X128.Cmp(big.NewInt(0)) != 0 {
		t.Errorf("Expected feeGrowthOutside1X128 to be 0, got %v", feeGrowthOutside1X128)
	}
	if initialized != true {
		t.Errorf("Expected initialized to be true, got %v", initialized)
	}
	testTicks.Clear(2)
}
