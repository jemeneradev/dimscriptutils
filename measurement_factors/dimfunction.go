package measurement_factors

import (
	list "container/list"
	"fmt"
)

type DimFunction struct {
	Name    string
	Args    *list.List
	secName string
}

type componentRep interface {
	ToComponentMap() map[string]interface{}
}

type referenceable interface {
	GetReferenceName() string
}

func NewDimFunctionCall(fn string, args *list.List, sec interface{}) *DimFunction {
	nf := new(DimFunction)
	nf.Name = fn
	nf.Args = args

	switch sections := sec.(type) {
	case *list.List:
		{
			backValue := sections.Back().Value
			sectionName, ok := backValue.(referenceable)
			if ok {
				nf.secName = sectionName.GetReferenceName()
			}
		}
	}
	//fmt.Fprintf(os.Stderr, "nf : %v\n", nf)
	return nf
}

type sectionEvaluator interface {
	EvaluateSectionAt(sectionName string, methodName string, args interface{}, cntx interface{}) interface{}
}

type valuegetterforfunc interface {
	GetTableValue(k string) interface{}
	EraseValueFromStore(k string) interface{}
}

func (fc *DimFunction) DetermineResults(context interface{}) interface{} {
	//fmt.Fprintf(os.Stderr, "dealing with a function call %T %v,context:%v\n", fc, fc, context)

	/* for e := fc.Args.Front(); e != nil; e = e.Next() {
		// do something with e.Value
		fmt.Fprintf(os.Stderr, "arg: %v\n", e)
	} */

	sectionRunner, ok := context.(sectionEvaluator)
	if ok {
		results := sectionRunner.EvaluateSectionAt(fc.secName, fc.Name, fc.Args, context)
		/* fmt.Fprintf(os.Stderr, "after dealing with func:\n")
		fmt.Fprintf(os.Stderr, "\t\tfunc: %v\n", fc)
		fmt.Fprintf(os.Stderr, "\t\tcontext: %v\n", context)
		fmt.Fprintf(os.Stderr, "\t\tresults: %v\n", results) */

		return results
	}
	return nil
}

func (fc *DimFunction) String() string {
	var argslice []interface{}

	for e := fc.Args.Front(); e != nil; e = e.Next() {
		// do something with e.Value
		argslice = append(argslice, e.Value)
	}

	return fmt.Sprintf("func:%v(%v) @%v ", fc.Name, argslice, fc.secName)
}
