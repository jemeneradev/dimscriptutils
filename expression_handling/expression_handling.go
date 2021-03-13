package expression_handling

import (
	"container/list"
	"fmt"
)

func NewExpressionStack() *list.List {
	return list.New()
}

func PushOperand(x interface{}, oprnType int, v interface{}) {
	switch l := x.(type) {
	case *list.List:
		{
			if l != nil {
				l.PushFront(NewOperand(oprnType, v))
			}
		}
	}
}

func PushBiOperator(x interface{}, op int) {
	switch l := x.(type) {
	case *list.List:
		{
			if l != nil {
				right := l.Remove(l.Front())
				left := l.Remove(l.Front())
				l.PushFront(NewDimBiOperator(op, left, right))
			}
		}
	}
}

func Pop(x interface{}) interface{} {
	switch l := x.(type) {
	case *list.List:
		{
			if l != nil && l.Len() > 0 {
				return l.Remove(l.Front())
			}
		}
	}
	return nil
}

func PushUniOperator(l *list.List, op int) {
	if l != nil {
		v := l.Remove(l.Front())
		l.PushFront(NewUniOperator(op, v))
	}
}

type encoder interface {
	Encode() string
}

func Debug(x interface{}) string {
	switch l := x.(type) {
	case *list.List:
		{
			for e := l.Front(); e != nil; e = e.Next() {
				switch v := e.Value.(type) {
				case encoder:
					{
						fmt.Printf("%v", v.Encode())
					}
				}
				//fmt.Printf("finished %v", e)
			}
		}
	}
	return ""
}
