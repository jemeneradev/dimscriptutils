package constants

const (
	StatusResolved = iota
	DimVectorOperand
	DimAssignmentOperand
	DimSectionOperand
	DimFunctionCallOperand
	DimNumOperand
	DimVarOperand
)

const (
	AddOperator = iota + 20
	SubOperator
	MulOperator
	DivOperator
)

const (
	DimWholeNumber = iota + 30 //0
	DimReal
	DimFraction
)

const EmptyString = ""

const (
	DimNonMeasurement = iota + 40
	DimFoot
	DimInch
	DimUndetermined
	DimReference
)

const (
	DimWord = iota + 50
	DimId
	DimAssign
)
