package expression_handling

import (
	"fmt"

	"github.com/jemeneradev/dimscriptutils/constants"
)

/*
Struct definition
*/

type DimOperand struct {
	OPType int
	Value  interface{}
}

func NewOperand(op int, v interface{}) *DimOperand {
	n := new(DimOperand)
	n.OPType = op
	n.Value = v
	return n
}

/*
Interfaces
*/

type resultsSolver interface {
	DetermineResults(interface{}) interface{}
}

/*Operations*/

func (op *DimOperand) Compute(context interface{}) {
	//fmt.Fprintf(os.Stderr, "\nBefore operand compute: %T\n", op.Value)
	//g, ok := op.Value.(resultsSolver)
	//fmt.Fprintf(os.Stderr, "\nop: %v\n", op)
	switch v := op.Value.(type) {
	case resultsSolver:
		//fmt.Printf("results\n")
		var result interface{}
		switch op.OPType {
		case constants.DimNumOperand:
			{
				//fmt.Fprintf(os.Stderr, "\ncompute num operand %v %T %v\n", v, v, context)
				result = v.DetermineResults(context)
				break
			}
		default:
			{
				//fmt.Fprintf(os.Stderr, "\ncompute default %v %T %v\n", v, v, context)
				result = v.DetermineResults(context)
				break
			}
		}

		op.OPType = constants.StatusResolved //encode to Solved
		op.Value = result
		//fmt.Fprintf(os.Stderr, "\nAfter operand compute: %v\n", op)
	}
}

func (op *DimOperand) Results() map[string]float64 {

	//fmt.Printf("Op.Value:%v\n", op)
	if op.OPType == constants.StatusResolved {
		s, ok := op.Value.(map[string]float64)
		if ok {
			return s
		}
	}
	return make(map[string]float64)
}

func (op *DimOperand) Encode() string {
	return fmt.Sprintf("<Operand: %v>", op.Value)
}

type identifier interface {
	Identifier() string
}

func (op *DimOperand) Identifier() string {
	switch name := op.Value.(type) {
	case identifier:
		{
			return name.Identifier()
		}
	}
	return "missing identifier"
}

func (op *DimOperand) GetValue() interface{} {
	return op.Value
}
