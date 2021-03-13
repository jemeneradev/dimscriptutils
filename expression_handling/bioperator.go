package expression_handling

import (
	//"fmt"
	"fmt"

	"github.com/jemeneradev/dimscriptutils/constants"
)

/*
Struct Definition
*/

/*
DimBiOperator Type(int to describe content), LValue ({}interface), RValue ({}interface)
*/
type DimBiOperator struct {
	Type   int
	LValue interface{}
	RValue interface{}
}

type TempOperand struct {
	Value interface{}
}

func NewTempOperand(op interface{}) interface{} {
	n := new(TempOperand)
	n.Value = op
	return n
}

func (t *TempOperand) Results() map[string]float64 {
	switch v := t.Value.(type) {
	case hasResults:
		{
			return v.Results()
		}
	}
	return nil
}

/*
NewDimBiOperator Constructor
*/
func NewDimBiOperator(opr int, l interface{}, r interface{}) *DimBiOperator {
	operator := new(DimBiOperator)
	operator.Type = opr
	operator.LValue = l
	operator.RValue = r
	return operator
}

/*Vector addition and subtraction*/
func mapDiff(l map[string]float64, r map[string]float64) {
	for kr, v := range r {
		_, ok := l[kr]
		if ok {
			l[kr] -= v
		} else {
			l[kr] = -v
		}
	}
}

/*Vector addition and subtraction*/
func mapMul(l map[string]float64, r map[string]float64) {
	for kr, v := range r {
		_, ok := l[kr]
		if ok {
			l[kr] *= v
		} else {
			l[kr] = v
		}
	}
}

/*Vector addition and subtraction*/
func mapDiv(l map[string]float64, r map[string]float64) {
	for kr, v := range r {
		_, ok := l[kr]
		if ok {
			l[kr] /= v
		} else {
			l[kr] = v
		}
	}
}

func mapUnion(l map[string]float64, r map[string]float64) {
	for kr, v := range r {
		_, ok := l[kr]
		if ok {
			l[kr] += v
		} else {
			l[kr] = v
		}
	}
}

/*
Interfaces
*/

type hasResults interface {
	Results() map[string]float64
}

/*
Operations
*/

/*
BiOperatorResolved return computed results of bioperator
*/
func (biopr *DimBiOperator) BiOperatorResolved() interface{} {
	return biopr.LValue
}

/*
Compute applies type operation(currently only +-) on value nodes
*/
func (biopr *DimBiOperator) Compute() interface{} {
	var lresults map[string]float64
	switch l := (biopr.LValue).(type) {
	case hasResults:
		{
			lresults = l.Results()
		}
	}
	var rresults map[string]float64
	switch r := (biopr.RValue).(type) {
	case hasResults:
		{
			rresults = r.Results()
		}
	}
	//fmt.Printf("\nOperator Compute: left(%v) right(%v)\n",lresults,rresults)
	if lresults != nil && rresults != nil {
		switch biopr.Type {
		case constants.AddOperator:
			{
				mapUnion(lresults, rresults)
			}
		case constants.SubOperator:
			{
				mapDiff(lresults, rresults)
			}
		case constants.MulOperator:
			{
				mapMul(lresults, rresults)
			}
		case constants.DivOperator:
			{
				//check if r is zero
				mapDiv(lresults, rresults)
			}
		}
		return biopr.LValue
	}
	return nil
}

/*
HandleBiOperator passes each value node through passed in handler function, with context
*/
func (biopr *DimBiOperator) HandleBiOperator(fn func(interface{}, interface{}), context interface{}) interface{} {
	fn(biopr.LValue, context)
	fn(biopr.RValue, context)
	return NewTempOperand(biopr.Compute())
}

/*
Results return computed results
*/
func (biopr *DimBiOperator) Results() map[string]float64 {
	switch v := biopr.LValue.(type) {
	case hasResults:
		{
			return v.Results()
		}
	}
	return nil
}

/*
Encode debug string, outputs formatted content
*/
func (biopr *DimBiOperator) Encode() string {

	var encodeStrings [2]string

	switch lop := biopr.LValue.(type) {
	case encoder:
		{
			encodeStrings[0] = lop.Encode()
		}
	default:
		{
			encodeStrings[0] = "_"
		}
	}

	switch rop := biopr.RValue.(type) {
	case encoder:
		{
			encodeStrings[1] = rop.Encode()
		}
	default:
		{
			encodeStrings[1] = "_"
		}
	}

	switch biopr.Type {
	case constants.AddOperator:
		{
			return fmt.Sprintf("<Operator: %v %v ADD>", encodeStrings[0], encodeStrings[1])
		}
	case constants.SubOperator:
		{
			return fmt.Sprintf("<Operator: %v %v SUB>", encodeStrings[0], encodeStrings[1])
		}
	case constants.MulOperator:
		{
			return fmt.Sprintf("<Operator: %v %v MUL>", encodeStrings[0], encodeStrings[1])
		}
	case constants.DivOperator:
		{
			return fmt.Sprintf("<Operator: %v %v DIV>", encodeStrings[0], encodeStrings[1])
		}
	}
	return ""
}
