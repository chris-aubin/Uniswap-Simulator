package fullMath

import (
	"math/big"
	"testing"

	"github.com/chris-aubin/Uniswap-Simulator/src/libraries/constants"
)

// Some basic tests for the fullMath package. Because the package is so simple
// these tests are very superficial.

// Check that division by zero causes panic.
func TestDivideByZeroMulDiv(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("MulDiv did not panic on division by zero.")
		}
	}()

	MulDiv(big.NewInt(1), big.NewInt(1), big.NewInt(0))
}

// Check that division by zero causes panic.
func TestDivideByZeroMulDivRoundingUp(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("MulDivRoundingUp did not panic on division by zero.")
		}
	}()

	MulDivRoundingUp(big.NewInt(1), big.NewInt(1), big.NewInt(0))
}

// Check that negative arguments cause panic.
func TestNegativeArgsMulDiv(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("MulDiv did not panic on division by zero.")
		}
	}()

	MulDiv(big.NewInt(-1), big.NewInt(1), big.NewInt(1))

	defer func() {
		if r := recover(); r == nil {
			t.Errorf("MulDiv did not panic on division by zero.")
		}
	}()

	MulDiv(big.NewInt(1), big.NewInt(-1), big.NewInt(1))

	defer func() {
		if r := recover(); r == nil {
			t.Errorf("MulDiv did not panic on division by zero.")
		}
	}()

	MulDiv(big.NewInt(1), big.NewInt(1), big.NewInt(-1))

	defer func() {
		if r := recover(); r == nil {
			t.Errorf("MulDiv did not panic on division by zero.")
		}
	}()

	MulDiv(big.NewInt(-1), big.NewInt(1), big.NewInt(-1))
}

// Check that negative arguments cause panic.
func TestNegativeArgsMulDivRoundingUp(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("MulDivRoundingUp did not panic on division by zero.")
		}
	}()

	MulDivRoundingUp(big.NewInt(-1), big.NewInt(1), big.NewInt(1))

	defer func() {
		if r := recover(); r == nil {
			t.Errorf("MulDivRoundingUp did not panic on division by zero.")
		}
	}()

	MulDivRoundingUp(big.NewInt(1), big.NewInt(-1), big.NewInt(1))

	defer func() {
		if r := recover(); r == nil {
			t.Errorf("MulDivRoundingUp did not panic on division by zero.")
		}
	}()

	MulDivRoundingUp(big.NewInt(1), big.NewInt(1), big.NewInt(-1))

	defer func() {
		if r := recover(); r == nil {
			t.Errorf("MulDivRoundingUp did not panic on division by zero.")
		}
	}()

	MulDivRoundingUp(big.NewInt(-1), big.NewInt(1), big.NewInt(-1))
}

// Check that results > 2^256-1 cause panic.
func TestOverflowMulDiv(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("MulDiv did not panic when result > 2^256-1")
		}
	}()

	MulDiv(constants.MaxUint256, big.NewInt(2), big.NewInt(1))

	defer func() {
		if r := recover(); r == nil {
			t.Errorf("MulDiv did not panic when result > 2^256-1")
		}
	}()

	MulDiv(new(big.Int).Add(constants.MaxUint256, big.NewInt(1)), big.NewInt(1), big.NewInt(1))

	defer func() {
		if r := recover(); r == nil {
			t.Errorf("MulDiv did not panic when result > 2^256-1")
		}
	}()

	MulDiv(constants.MaxUint256, constants.MaxUint256, new(big.Int).Sub(constants.MaxUint256, big.NewInt(1)))
}

// Check that results > 2^256-1 cause panic.
func TestOverflowMulDivRoundingUp(t *testing.T) {
	// Check that overflow panics
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("MulDivRoundingUp did not panic when result > 2^256-1")
		}
	}()

	MulDivRoundingUp(constants.MaxUint256, big.NewInt(2), big.NewInt(1))

	defer func() {
		if r := recover(); r == nil {
			t.Errorf("MulDivRoundingUp did not panic when result > 2^256-1")
		}
	}()

	MulDivRoundingUp(new(big.Int).Add(constants.MaxUint256, big.NewInt(1)), big.NewInt(1), big.NewInt(1))

	defer func() {
		if r := recover(); r == nil {
			t.Errorf("MulDivRoundingUp did not panic when result > 2^256-1")
		}
	}()

	MulDivRoundingUp(constants.MaxUint256, constants.MaxUint256, new(big.Int).Sub(constants.MaxUint256, big.NewInt(1)))
}

// Check that MulDiv is accurate without intermediate overflow and without
// rounding.
func TestNoOverflowMulDiv(t *testing.T) {
	expected := big.NewInt(0)
	got := MulDiv(big.NewInt(0), big.NewInt(1), big.NewInt(1))
	if expected.Cmp(got) != 0 {
		t.Errorf("MulDiv(0, 1, 1) = %v; want %v", got, expected)
	}

	expected = big.NewInt(1)
	got = MulDiv(big.NewInt(1), big.NewInt(1), big.NewInt(1))
	if expected.Cmp(got) != 0 {
		t.Errorf("MulDiv(1, 1, 1) = %v; want %v", got, expected)
	}

	expected = big.NewInt(5)
	got = MulDiv(big.NewInt(15), big.NewInt(2), big.NewInt(6))
	if expected.Cmp(got) != 0 {
		t.Errorf("MulDiv(15, 2, 6) = %v; want %v", got, expected)
	}
}

// Check that MulDiv is accurate without intermediate overflow and without
// rounding.
func TestNoOverflowMulDivRoundingUp(t *testing.T) {
	expected := big.NewInt(0)
	got := MulDivRoundingUp(big.NewInt(0), big.NewInt(1), big.NewInt(1))
	if expected.Cmp(got) != 0 {
		t.Errorf("MulDivRoundingUp(0, 1, 1) = %v; want %v", got, expected)
	}

	expected = big.NewInt(1)
	got = MulDivRoundingUp(big.NewInt(1), big.NewInt(1), big.NewInt(1))
	if expected.Cmp(got) != 0 {
		t.Errorf("MulDivRoundingUp(1, 1, 1) = %v; want %v", got, expected)
	}

	expected = big.NewInt(5)
	got = MulDivRoundingUp(big.NewInt(15), big.NewInt(2), big.NewInt(6))
	if expected.Cmp(got) != 0 {
		t.Errorf("MulDivRoundingUp(15, 2, 6) = %v; want %v", got, expected)
	}
}

// Check that MulDiv is accurate without intermediate overflow and with
// rounding.
func TestRoundingMulDiv(t *testing.T) {
	expected := big.NewInt(2)
	got := MulDiv(big.NewInt(5), big.NewInt(3), big.NewInt(6))
	if expected.Cmp(got) != 0 {
		t.Errorf("MulDiv(5, 3, 6) = %v; want %v", got, expected)
	}
}

// Check that MulDivRoundingUp is accurate without intermediate overflow and
// with rounding.
func TestRoundingMulDivRoundingUp(t *testing.T) {
	expected := big.NewInt(3)
	got := MulDivRoundingUp(big.NewInt(5), big.NewInt(3), big.NewInt(6))
	if expected.Cmp(got) != 0 {
		t.Errorf("MulDivRoundingUp(5, 3, 6) = %v; want %v", got, expected)
	}
}

// Check that MulDiv is accurate with intermediate overflow and without
// rounding.
func TestIntermediateOverflowMulDiv(t *testing.T) {
	expected := constants.MaxUint256
	got := MulDiv(constants.MaxUint256, constants.MaxUint256, constants.MaxUint256)
	if expected.Cmp(got) != 0 {
		t.Errorf("MulDiv(2^256-1, 2^256-1, 2^256-1) = %v; want %v", got, expected)
	}
}

// Check that MulDivRoundingUp is accurate with intermediate overflow and
// without rounding.
func TestIntermediateOverflowMulDivRoundingUp(t *testing.T) {
	// Check that MulDivRoundingUp works correctly
	expected := constants.MaxUint256
	got := MulDivRoundingUp(constants.MaxUint256, constants.MaxUint256, constants.MaxUint256)
	if expected.Cmp(got) != 0 {
		t.Errorf("MulDivRoundingUp(2^256-1, 2^256-1, 2^256-1) = %v; want %v", got, expected)
	}
}
