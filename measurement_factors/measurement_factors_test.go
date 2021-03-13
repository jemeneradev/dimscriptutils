package measurement_factors

import (
	"fmt"
	"testing"

	"github.com/jemeneradev/dimscriptutils/constants"
)

func TestWholeNummberMeasurementInches(t *testing.T) {
	want := float64(1.0)
	input := NewDimMeasurement("1", constants.DimInch, constants.DimWholeNumber)
	if got := input.ToInches(); got != want {
		t.Errorf("CategorizeNumeric(%v) = %f, want %f", input, got, want)
	}
}

func TestWholeNummberMeasurementFeet(t *testing.T) {
	want := float64(12.0)
	input := NewDimMeasurement("1", constants.DimFoot, constants.DimWholeNumber)
	if got := input.ToInches(); got != want {
		t.Errorf("CategorizeNumeric(%v) = %f, want %f", input, got, want)
	}
}

func TestRealNummberMeasurementInches(t *testing.T) {
	want := float64(1.0)
	input := NewDimMeasurement("1.0", constants.DimInch, constants.DimReal)
	if got := input.ToInches(); got != want {
		t.Errorf("CategorizeNumeric(%v) = %f, want %f", input, got, want)
	}
}

func TestRealNumberMeasurementFeet(t *testing.T) {
	want := float64(12.0)
	input := NewDimMeasurement("1.0", constants.DimFoot, constants.DimReal)
	if got := input.ToInches(); got != want {
		t.Errorf("CategorizeNumeric(%v) = %f, want %f", input, got, want)
	}
}

func TestMeasurementEncoding(t *testing.T) {
	want := fmt.Sprintf("<dimM{ %v %v %v }>", "1", constants.DimFoot, constants.DimReal)
	input := NewDimMeasurement("1", constants.DimFoot, constants.DimReal)
	if got := input.Encode(); got != want {
		t.Errorf("Encode(%q) got: %v, want %v", *input, got, want)
	}

	want = fmt.Sprintf("<dimM{ %v %v %v }>", "1", constants.DimInch, constants.DimWholeNumber)
	input = NewDimMeasurement("1", constants.DimInch, constants.DimWholeNumber)
	if got := input.Encode(); got != want {
		t.Errorf("Encode(%q) got: %v, want %v", *input, got, want)
	}
}

func TestMeasurementDecoding(t *testing.T) {

	want := NewDimMeasurement("4", constants.DimFoot, constants.DimReal).Encode()
	input := want
	if got := DimMeasurementDecode(input).Encode(); got != want {
		t.Errorf("DimMeasureDecode(%s) got: %s, want %s", input, got, want)
	}

}

func TestMeasurementListEncoding(t *testing.T) {
	tList := NewMeasurementList()
	tList.PushBack(NewDimMeasurement("4", constants.DimInch, constants.DimWholeNumber))
	tList.PushBack(NewDimMeasurement("40", constants.DimFoot, constants.DimWholeNumber))
	tList.PushBack(NewDimMeasurement("2.0", constants.DimInch, constants.DimReal))
	tList.PushBack(NewDimMeasurement("3.5", constants.DimFoot, constants.DimReal))
	want := EncodeList(tList)
	input := want
	if got := EncodeList(DecodeList(input)); got != want {
		t.Errorf("DimMeasureDecode(%s) got: %s, want %s", input, got, want)
	}
}

func TestCharacteristicEncodingAndDecoding(t *testing.T) {
	tList := NewMeasurementList()
	tList.PushBack(NewDimMeasurement("4", constants.DimInch, constants.DimWholeNumber))
	tList.PushBack(NewDimMeasurement("40", constants.DimFoot, constants.DimWholeNumber))
	tList.PushBack(NewDimMeasurement("3.5", constants.DimNonMeasurement, constants.DimReal))
	tCharacteristic := NewDimCharacteristic("trim", tList)

	want := tCharacteristic.Encode()
	input := want
	if got := DimCharacteristicDecode(input).Encode(); got != want {
		t.Errorf("DimMeasureDecode(%s) got: %s, want %s", input, got, want)
	}
}
