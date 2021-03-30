package evaluator

type uniOperatorHandler interface {
	HandleUniOperator()
}

type computable interface {
	Compute(interface{})
}

type computableWithFunc interface {
	Compute(interface{}, func(interface{}, interface{}) map[string]float64)
}

type biOperatorHandler interface {
	HandleBiOperator(left interface{}, right interface{}, context interface{}) interface{}
}

type hasResults interface {
	Results() map[string]float64
}

type hasNumberRepresentation interface {
	ToNumber() float64
}

func processMember(item interface{}, context interface{}) {
	//fmt.Fprintf(os.Stderr, "%T %T\n", item, context)
	switch v := item.(type) {
	case computable:
		v.Compute(context)
	case computableWithFunc:
		{
			//fmt.Fprintf(os.Stderr, "func %v %v\n", item, context)
			v.Compute(context, Evaluate)
		}
	}
}

/*
Evaluate passed in item gets calculated and return results. Requires a context object.
*/
func Evaluate(item interface{}, passedContext interface{}) map[string]float64 {
	//fmt.Fprintf(os.Stderr, "eval:: %v %T\n", item, item)

	if item != nil {
		switch passedContext.(type) {
		case *DimContext:
			{
				//fmt.Fprintf(os.Stderr, "here before processMember %v\n", item)
				processMember(item, passedContext)
				//fmt.Fprintf(os.Stderr, "here after processMember %v %T\n", item, item)
			}
		}
		switch v := item.(type) {
		case hasResults:
			{
				return v.Results()
			}
		}
	}
	return nil
}
