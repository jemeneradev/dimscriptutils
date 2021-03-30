package expression_handling

import (
	"container/list"
)

type DimExpList struct {
	list  interface{}
	Value interface{}
}

func NewDimExpList(l interface{}) *DimExpList {
	dimlist := new(DimExpList)
	dimlist.list = l
	return dimlist
}

type biOperatorHandler interface {
	HandleBiOperator(left interface{}, right interface{}, context interface{}) interface{}
}

func (explist *DimExpList) Compute(context interface{}, evalFunc func(interface{}, interface{}) map[string]float64) {

	switch v := (explist.list).(type) {
	case *list.List:
		{
			exp := make(map[int]interface{}, v.Len())
			i := 0
			for e := v.Back(); e != nil; e = e.Prev() {
				switch x := e.Value.(type) {
				case biOperatorHandler:
					{
						exp[i] = x.HandleBiOperator(exp[i-2], exp[i-1], context)
						//fmt.Fprintf(os.Stderr, "\nr(%v):%v", i, exp[i])
						i--
						//i = 0
					}
					////fmt.Fprintf(os.Stderr, "bi:%v\n", item)
				default:
					{
						exp[i] = evalFunc(x, context)
						//fmt.Fprintf(os.Stderr, "\ni:%v", i)
						i++
					}
				}
			}
			explist.Value = exp[0]
		}
	}
}

func (explist *DimExpList) Results() map[string]float64 {
	return explist.Value.(map[string]float64)
}
