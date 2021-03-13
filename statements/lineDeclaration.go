package statements

import (
	"container/list"
	"fmt"
)

func NewSectionStack() *list.List {
	return list.New().Init()
}

func AddToSectionStack(x interface{}, item interface{}) {
	switch stack := x.(type) {
	case *list.List:
		{
			stack.PushFront(item)
		}
	}
}

func PopSectionStack(x interface{}) interface{} {
	switch stack := x.(type) {
	case *list.List:
		{
			front := stack.Front()
			if front != nil {
				return stack.Remove(front)
			}
		}
	}
	return nil
}

func AppendStatementToSectionStack(x interface{}, item interface{}) {
	switch stack := x.(type) {
	case *list.List:
		{
			//fmt.Fprintf(os.Stderr, "assign to sectionStack: %v\n", stack.Front().Value)
			front := stack.Front()
			if front != nil {
				switch valid := front.Value.(type) {
				case *DimSectionDeclaration:
					{
						//fmt.Fprintf(os.Stderr, "\titem:%v\n", item)
						valid.AddToStatements(item)
					}
				}
			}
		}
	}
}

type methodSetter interface {
	SetMethod(name string, index int)
	GetStatementCount() int
}

func AddMethodToSection(section interface{}, mth string) {

	switch stack := section.(type) {
	case *list.List:
		{
			top := stack.Back()
			switch definingSection := top.Value.(type) {
			case methodSetter:
				{
					//fmt.Fprintf(os.Stderr, "added method %v at int %v\n", mth, definingSection.GetStatementCount())
					definingSection.SetMethod(mth, definingSection.GetStatementCount())
				}
			}
		}
	}
}

func DebugSectionStack(x interface{}) {
	switch l := x.(type) {
	case *list.List:
		{
			fmt.Printf("%v %v", l, l.Len())
		}
	}
}
