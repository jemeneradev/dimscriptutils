package evaluator

type uniOperatorHandler interface {
	HandleUniOperator()
}

type computable interface {
	Compute(interface{})
}

type biOperatorHandler interface {
	HandleBiOperator(func(interface{}, interface{}), interface{}) interface{}
}

type hasResults interface {
	Results() map[string]float64
}

type hasNumberRepresentation interface {
	ToNumber() float64
}

func processMember(item interface{}, context interface{}) {
	////fmt.Fprintf(os.Stderr, "%T %T\n", item, context)
	switch v := item.(type) {
	case biOperatorHandler:
		////fmt.Fprintf(os.Stderr, "bi:%v\n", item)
		v.HandleBiOperator(processMember, context)
	case uniOperatorHandler:
		////fmt.Printf("\nx is UniOperator %v\n", v)             // here v has type int
	case computable:
		//fmt.Fprintf(os.Stderr, "op:%v - %T %v\n", item, context, context)
		v.Compute(context)
	}
}

/*
Evaluate passed in item gets calculated and return results. Requires a context object.
*/
func Evaluate(item interface{}, passedContext interface{}) map[string]float64 {
	//fmt.Fprintf(os.Stderr, "eval:: %T %T\n", item, passedContext)

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
	return make(map[string]float64)
}
