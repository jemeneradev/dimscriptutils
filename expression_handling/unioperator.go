package expression_handling

type DimUniOperator struct {
	Opertator int
	Value interface{}
}

func NewUniOperator(opr int, v interface{}) *DimUniOperator {
	operator := new(DimUniOperator)
	operator.Opertator = opr
	operator.Value = v
	return operator
}
