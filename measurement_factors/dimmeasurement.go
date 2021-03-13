package measurement_factors

import (
	"fmt"
	"regexp"
	"strconv"

	"github.com/jemeneradev/dimscriptutils/constants"
)

/*
DimMeasurement represents measurements or numeric
value
Indicator, i,e. what is it? Inch, Foot? Future representation of yard,mm,cm,meter
Precision: whole, float, or fraction, but again it's generic
*/
type DimMeasurement struct {
	value     string
	Indicator int //inch or feet
	precision int //whole or real
}

var splitterRegExp = regexp.MustCompile(`[_/]`)

/*
NewDimMeasurement constructor return a measurement allocation
*/
func NewDimMeasurement(str string, indicator int, precision int) *DimMeasurement {
	n := new(DimMeasurement)
	n.value = str
	n.Indicator = indicator
	n.precision = precision
	return n
}

func (measurement *DimMeasurement) String() string {
	return fmt.Sprintf("dim: value(%v) indicator(%v) precision(%v)", measurement.value, measurement.Indicator, measurement.precision)
}

/*
UpdateIndicator updates Indicator
*/
func (measurement *DimMeasurement) UpdateIndicator(indicator int) {
	measurement.Indicator = indicator
}

/*
ToInches returns measurement representation in inches.
*/
func (measurement *DimMeasurement) ToInches() float64 {
	var result float64
	result = 0.0
	val := measurement.value
	switch measurement.precision {
	case constants.DimWholeNumber:
		{
			switch measurement.Indicator {
			case constants.DimInch:
				{
					if s, err := strconv.Atoi(val); err == nil {
						result = float64(s)
					}
				}
			case constants.DimFoot:
				{
					if s, err := strconv.Atoi(val); err == nil {
						result = float64(s * 12)
					}
				}
			}
		}
	case constants.DimFraction:
		{
			//representation format whole_numerator/denominator
			numval := splitterRegExp.Split(val, -1)
			whole, _ := strconv.ParseFloat(numval[0], 64)
			num, _ := strconv.ParseFloat(numval[1], 64)
			denom, _ := strconv.ParseFloat(numval[2], 64)

			switch measurement.Indicator {
			case constants.DimInch:
				{
					result = whole + (num / denom)
				}
			case constants.DimFoot:
				{
					result = (whole + (num / denom)) * float64(12)
				}
			}
		}
	case constants.DimReal:
		{
			switch measurement.Indicator {
			case constants.DimInch:
				{
					if s, err := strconv.ParseFloat(val, 64); err == nil {
						result = float64(s)
					}
				}
			case constants.DimFoot:
				{
					if s, err := strconv.ParseFloat(val, 64); err == nil {
						result = float64(s * 12)
					}
				}
			}
		}
	}
	return result
}

/* type workingContext interface {
	IsCountMutable() bool
	IncrementCount(float64)
	GetHandler(string) interface{}
} */

func (measurement *DimMeasurement) DetermineResults(context interface{}) interface{} {
	res := make(map[string]float64)
	//res["val"] = measurement.ToNumber()
	//fmt.Printf("%v\n", measurement)
	switch wc := context.(type) {
	case workingContext:
		{
			switch measurement.Indicator {
			case constants.DimReference:
				{
					switch measurement.precision {
					case constants.DimVarOperand:
						{
							lookres := wc.LookUpVar(measurement.value)
							if lookres != nil {
								switch m := lookres.(type) {
								case map[string]float64:
									{
										val, ok := m["value"]
										//fmt.Fprintf(os.Stderr, "lookup val: %v %T %v\n", m, val, ok)
										if ok {
											res["value"] = val
										}
									}
								}
							}
							//fmt.Fprintf(os.Stderr, "lookup res: %v %T\n", lookres, lookres)
						}
					}
				}
			default:
				{
					res["value"] = measurement.ToNumber()
				}
			}
		}
	}
	//fmt.Printf("measurement %v %v\n", measurement, res)
	return res
}

/*
ToNumber returns numerical representation if not measurement
*/
func (measurement *DimMeasurement) ToNumber() float64 {
	if measurement != nil {
		switch measurement.Indicator {
		case constants.DimNonMeasurement:
			switch measurement.precision {
			case constants.DimWholeNumber:
				{
					if s, err := strconv.Atoi(measurement.value); err == nil {
						return float64(s)
					}
				}
			case constants.DimReal:
				{
					if s, err := strconv.ParseFloat(measurement.value, 64); err == nil {
						return s
					}
				}
			case constants.DimFraction:
				{
					numval := splitterRegExp.Split(measurement.value, -1)
					whole, _ := strconv.ParseFloat(numval[0], 64)
					num, _ := strconv.ParseFloat(numval[1], 64)
					denom, _ := strconv.ParseFloat(numval[2], 64)

					return whole + (num / denom)
				}
			}
		default:
			return measurement.ToInches()
		}
	}
	return 1
}

/*
Encode debug string
*/
func (measurement *DimMeasurement) Encode() string {
	if measurement != nil {
		return fmt.Sprintf("<dimM{ %s %d %d }>", measurement.value, measurement.Indicator, measurement.precision)
	}
	return ""
}

/*
CombineDimMeasurements given a set of measurement (ex: 16'2"), return the sum of the components
*/
func CombineDimMeasurements(measurements []*DimMeasurement) *DimMeasurement {
	var combinedValue float64 = 0.0
	for _, measurement := range measurements {
		if measurement != nil {
			combinedValue += measurement.ToInches()
		}
	}
	return NewDimMeasurement(strconv.FormatFloat(combinedValue, 'f', -1, 64), constants.DimInch, constants.DimReal)
}
