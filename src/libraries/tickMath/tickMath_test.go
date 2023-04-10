package tickMath

import (
	"math/big"
	"testing"
)

var (
	// Test inputs for the getSqrtRatioAtTick function. Generated using python.
	// Just random numbers between -MaxTick and MaxTick.
	getSqrtRatioAtTickInputs = [...]int{
		-605490,
		306435,
		795199,
		714479,
		-754928,
		-352354,
		-639979,
		151730,
		708581,
		55379,
		103098,
		479217,
		-91162,
		766800,
		-446965,
		-690435,
		135837,
		-827824,
		865454,
		-69784,
	}
	// Test outputs for the getSqrtRatioAtTick function. To be converted to
	// big.Ints. Generated outputs by redeploying the TickMath contract on
	// remix, changing the functions to public, and calling them with the
	// inputs above.
	getSqrtRatioAtTickOutputsStrings = [...]string{
		"5642721942098613",
		"357024793347176040098758073744472598",
		"14640069913673683022118270373753038941211362350",
		"258713002688033585140685610797793303133093091",
		"3211041733340",
		"1770077280900890172857",
		"1006020832475117",
		"156131288181857195385307075932564",
		"192642105850343801639492348401608095410822540",
		"1262871511813689486900606186164",
		"13724997649175287121978152899130",
		"2015692370062636309172520184117127683240",
		"830655370063815383322100753",
		"3539138290718651736992881117902513786070984607",
		"15618638830384903494",
		"80727895949757",
		"70533443559749144576888219566265",
		"83909003404",
		"490947235291057771851277398888386102982029559644",
		"2418883986714335308443808732",
	}
	// Test inputs for the getTickAtSqrtRatio function. To be converted to
	// big.Ints. Generated using python. Just random numbers between
	// MinSqrtRatio and MaxSqrtRatio.
	getTickAtSqrtRatioInputsStrings = [...]string{
		"3077995146750009554232817326544692445012838697",
		"1171870553120907383840441369152020973416515670267",
		"1317437710319558993890513334665996507864544485332",
		"949312744573864037520356406667233232129045661750",
		"557130615467757883687391181149441695989205072065",
		"1060834765302389972855932676552438198788165721756",
		"639973425120451021526991715212955122066196114986",
		"505241719688214432377574809931101430802396722834",
		"671711214394215808463855237755835919160357718146",
		"608247661852599656319514316922434872207288615322",
		"938714268764441452424871518440175100351912737089",
		"1057599593845233224481936657087383856016124548322",
		"486234556123076272274164831514246158348217477984",
		"731919579266908753537933265034320657780973167735",
		"1213098519448920482978784267975667732511098472607",
		"415299379856527523183447912761904495345604516351",
		"1236718913011199338974364919957735163867821909316",
		"1247081739193507112526927893700064868810793036616",
		"1165407617432778077899261458006446412740509115898",
		"536552080601665919177922730362914296449667818336",
	}
	// Test outputs for the getTickAtSqrtRatio function. Generated outputs by
	// redeploying the TickMath contract on remix, changing the functions to
	// public, and calling them with the inputs above.
	getTickAtSqrtRatioOutputs = [...]int{
		764007,
		882855,
		885197,
		878642,
		867983,
		880864,
		870756,
		866028,
		871724,
		869739,
		878418,
		880803,
		865261,
		873441,
		883546,
		862107,
		883932,
		884099,
		882744,
		867230,
	}
	// Slice to store the big.Int outputs for the getSqrtRatioAtTick function.
	getSqrtRatioAtTickOutputs []*big.Int
	// Slice to store the big.Int inputs for the getTickAtSqrtRatioInputs.
	getTickAtSqrtRatioInputs []*big.Int
)

// init is used to initialize the big.Int slices used in the tests.
func init() {
	for _, input := range getTickAtSqrtRatioInputsStrings {
		temp := new(big.Int)
		temp.SetString(input, 10)
		getTickAtSqrtRatioInputs = append(getTickAtSqrtRatioInputs, temp)
	}

	for _, output := range getSqrtRatioAtTickOutputsStrings {
		temp := new(big.Int)
		temp.SetString(output, 10)
		getSqrtRatioAtTickOutputs = append(getSqrtRatioAtTickOutputs, temp)
	}
}

// TestGetSqrtRatioAtTick tests the GetSqrtRatioAtTick function.
func TestGetSqrtRatioAtTick(t *testing.T) {
	for i, input := range getSqrtRatioAtTickInputs {
		output := GetSqrtRatioAtTick(input)
		if output.Cmp(getSqrtRatioAtTickOutputs[i]) != 0 {
			t.Errorf("GetSqrtRatioAtTick(%d) = %v, want %v", input, output, getSqrtRatioAtTickOutputs[i])
		}
	}
}

// TestGetTickAtSqrtRatio tests the GetTickAtSqrtRatio function.
func TestGetTickAtSqrtRatio(t *testing.T) {
	for i, input := range getTickAtSqrtRatioInputs {
		output := GetTickAtSqrtRatio(input)
		if output != getTickAtSqrtRatioOutputs[i] {
			t.Errorf("GetTickAtSqrtRatio(%s) = %d, want %d", input, output, getTickAtSqrtRatioOutputs[i])
		}
	}
}
