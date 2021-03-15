package constants

const (
	/*
		Indicates Operand has been computed and results have been determined.
	*/
	StatusResolved = iota

	/*
		Indicates Operand as Vector
		Format: {meaurements_list},{characteristics_list}[[@]{classification}][[!]{count_scalar}]
		Example: 6",5"@dim!5
	*/
	DimVectorOperand

	/*
		Indicates Operand as Assignment Statement
		Format: {description}[:]{measurements_vector_expression}
		Example; simple calc: 1',4' + 5',64'
	*/
	DimAssignmentOperand
	/*
		Indicates Operand as Section Statement
	*/
	DimSectionOperand
	/*
		Indicates Operand as FunctionCall
	*/
	DimFunctionCallOperand
	/*
		Indicates Operand as Number
	*/
	DimNumOperand
	/*
		Indicates Operand as Variable
	*/
	DimVarOperand
)

const (
	AddOperator = iota + 20
	SubOperator
	MulOperator
	DivOperator
)

const (
	/*
		Indicates Whole Number
		Format: [0-9]+, 0+ is treated as 0
	*/
	DimWholeNumber = iota + 30 //0
	/*
		Indicates Real Number
		Format: [0-9]+[.][0-9]+, 0+ is treated as 0
	*/
	DimReal
	/*
		Indicates Fraction Number
		Format: {whole}_{Numerator}/{Denominator}, where each is a whole number
		Examples: _1/2, 3_2/7
	*/
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
