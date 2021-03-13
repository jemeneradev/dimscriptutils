package statements

import (
	"container/list"
	"fmt"
	"os"
	"regexp"
	"strings"
)

var param_re = regexp.MustCompile(`(?m)(?P<arg_def>([[:blank:]]*(?P<arg_type>[[:alpha:]][a-zA-Z_]*))[[:blank:]]*:[[:blank:]]*(?P<arg_val>[[:alpha:]][[:alnum:]]*)[[:blank:]]*)`)

type computable interface {
	Compute(interface{})
}
type resultsSolver interface {
	DetermineResults(interface{}) interface{}
}

type DimSectionDeclaration struct {
	Name       string
	Statements *list.List
	context    interface{}
	values     map[string]interface{}
	acc        map[string]float64

	methodIndex           map[string]int //Answer where I can find given method/section with Statements list
	eval                  *(func(interface{}, interface{}) map[string]float64)
	params                interface{} //map[string]string
	hasVariableReferences bool
	executing             bool
}

func NewDimSectionDeclaration(name string, ismain bool) *DimSectionDeclaration {
	sec := new(DimSectionDeclaration)
	sec.Name = name
	sec.Statements = list.New().Init()
	sec.params = nil
	sec.hasVariableReferences = false
	sec.executing = ismain
	sec.values = make(map[string]interface{})
	return sec
}

func CreateSectionParams(para_str string) map[string]string {
	//TODO: check for correct format,if not send error
	params := make(map[string]string)
	//fmt.Fprintf(os.Stderr, "all params found:%v\n", param_re.FindAllString(para_str, -1))
	for _, match := range param_re.FindAllString(para_str, -1) {
		//fmt.Fprintf(os.Stderr, "found => %v\n", match)
		parts := strings.Split(match, ":")
		//fmt.Fprintf(os.Stderr, "%v :==> %v\n", parts[0], parts[1])
		params[parts[1]] = parts[0]
		//fmt.Fprintf(os.Stderr, "Params:%v\n", params)
		//TODO: check for duplicates
	}
	//fmt.Fprintf(os.Stderr, "params created:%v\n", params)
	return params
}

func (sec *DimSectionDeclaration) Encode() string {
	//fmt.Printf("statements:{\n")
	sec.Compute(sec.context)
	//fmt.Printf("\ncontext:\n%v\n", sec.context)
	return fmt.Sprintf("%v", sec.context)
}

func (sec *DimSectionDeclaration) SetHasVariableReferences(b bool) {
	sec.hasVariableReferences = b
}

func (sec *DimSectionDeclaration) SetExecuting(b bool) {
	sec.executing = b
}

func (sec *DimSectionDeclaration) GetReferenceName() string {
	return sec.Name
}

type getValue interface {
	GetValue() interface{}
}

func (sec *DimSectionDeclaration) GetStatementAt(mthd string) interface{} {
	methodLoc, found := sec.methodIndex[mthd]
	if found == false {
		return nil
	}
	i := 0
	var item interface{}
	item = nil
	for e := sec.Statements.Back(); e != nil; e = e.Prev() {

		if i == methodLoc {
			fmt.Fprintf(os.Stderr, "\tstatement: %v %T %v\n", e.Value, e.Value, methodLoc)
			hasValue, ok := e.Value.(getValue)
			if ok {
				item = hasValue.GetValue()
			}
			break
		}
		i++
	}

	return item
}

func (sec *DimSectionDeclaration) SetParams(params map[string]string) {
	sec.params = params
}

func (sec *DimSectionDeclaration) SetParamsValues(pv interface{}) {
	//fmt.Fprintf(os.Stderr, "parameters: %v %T\n", sec.params, pv)

	switch passedinValues := pv.(type) {
	case map[int]interface{}:
		{
			i := 0
			switch workingParams := sec.params.(type) {
			case map[string]string:
				{
					for pname, _ := range workingParams { //_ = ptype
						//fmt.Fprintf(os.Stderr, "param[%v]:%v -> %v\n", pname, ptype, passedinValues[i])
						sec.values[pname] = passedinValues[i]
						i++
					}
				}
			}
		}
	}
	//sec.executing = true
	//fmt.Fprintf(os.Stderr, "params vals: %v\n", sec.values)
}

func (sec *DimSectionDeclaration) SetMethod(name string, index int) {
	if sec.methodIndex == nil {
		sec.methodIndex = make(map[string]int)
	}
	sec.methodIndex[name] = index
	//fmt.Fprintf(os.Stderr, "\n\nsection(%v):%v\n\n", sec.Name, sec.methodIndex)
}
func (sec *DimSectionDeclaration) GetStatementCount() int {
	return sec.Statements.Len()
}

func (sec *DimSectionDeclaration) SetEval(eval func(interface{}, interface{}) map[string]float64) {
	sec.eval = &eval
}

type retrievable interface {
	RetrieveStore() interface{}
}
type accStorer interface {
	StoreAccumulator(name string)
}
type accStacker interface {
	AccumulatePush(val interface{})
	AccumulatePop() interface{}
}
type sectionIndexer interface {
	AddSectionToIndex(n string, s interface{})
}

func (sec *DimSectionDeclaration) SetContext(cnt interface{}) {
	sec.context = cnt
}

type canlookup interface {
	PushLookUp(n string)
	PopLookUp()
}
type accumulator interface {
	Accumulate(val interface{})
}

func (sec *DimSectionDeclaration) determineContext(context interface{}) interface{} {
	var workingContext interface{}
	if sec.context == nil {
		//fmt.Fprintf(os.Stderr, "context acc\n")
		workingContext = context
	} else {
		//fmt.Fprintf(os.Stderr, "section acc\n")
		workingContext = sec.context
		contextSectionIndex, ok := context.(sectionIndexer)
		if ok {
			contextSectionIndex.AddSectionToIndex(sec.Name, sec)
		}
	}

	return workingContext
}

type storer interface {
	InsertIntoStore(key string, val interface{})
}

func (sec *DimSectionDeclaration) DetermineResults(context interface{}) interface{} {
	//results := make(map[string]float64))

	//fmt.Fprintf(os.Stderr, "solving %v\n\thasVars%v\n\tcontext:%v\n\tparams:%v\n", sec.Name, sec.hasVariableReferences, context, sec.params)
	//fmt.Fprintf(os.Stderr, "\nDetermining results(%v):\n", sec.Name)

	if sec.executing == true || sec.params == nil { // nil sec.params marks a section division, helpful in dividing output results
		workingContext := sec.determineContext(context)
		//fmt.Fprintf(os.Stderr, "\tworking context: %v\n", workingContext)

		if sec.acc == nil {
			sec.acc = make(map[string]float64)
		}

		//push section into context
		switch acc := workingContext.(type) {
		case accStacker:
			{
				//fmt.Fprintf(os.Stderr, "\nsection: %v", sec.Name)
				acc.AccumulatePush(sec.acc) //push section accumulator
			}
		}

		//make section searchable
		switch lookupPush := context.(type) {
		case canlookup:
			{
				lookupPush.PushLookUp(sec.Name) //push section accumulator
			}
		}

		//insert section vars into scope, if any
		if sec.params != nil {
			storer, ok := context.(storer)
			if ok {
				for k, v := range sec.values {
					//fmt.Fprintf(os.Stderr, "\nvalue: %v -> %v\n\n", k, v)
					storer.InsertIntoStore(k, v)
				}
			}
		}

		//fmt.Fprintf(os.Stderr, "got here %v\n", sec.Name)
		switch acc := workingContext.(type) {
		case accumulator:
			{
				for e := sec.Statements.Back(); e != nil; e = e.Prev() {
					// do something with e.Value
					switch v := e.Value.(type) {
					case computable:
						{
							//fmt.Fprintf(os.Stderr, "\nbefore every statement - statement: %v\n", e)
							//fmt.Fprintf(os.Stderr, "hasOwnContext: %v\n", hasOwnContext)
							v.Compute(workingContext)
							//fmt.Fprintf(os.Stderr, "\nafter every statement %v - context: %v\n", v, workingContext)

							acc.Accumulate(v)

							//fmt.Fprintf(os.Stderr, "\n======\n")
							//fmt.Fprintf(os.Stderr, "after assignment acc[%v] %v\n", results, context)
						}
					}
				}
			}
		}

		//fmt.Fprintf(os.Stderr, "\nEnd of Section: %v Acc:%v context:%v\n\n", sec.Name, sec.acc, workingContext)
		//return sec

		if sec.params == nil && sec.executing == true {
			switch saveAcc := workingContext.(type) {
			case accStorer:
				{
					saveAcc.StoreAccumulator(sec.Name)
				}
			}

		}

		switch acc := workingContext.(type) {
		case accStacker:
			{
				acc.AccumulatePop()
			}
		}
		switch lookupPop := context.(type) {
		case canlookup:
			{
				lookupPop.PopLookUp()
			}
		}
		return sec.acc
	}
	return sec
}

type wipefuncfromStore interface {
	EraseValueFromStore(k string)
}

func (sec *DimSectionDeclaration) Compute(context interface{}) {
	//fmt.Fprintf(os.Stderr, "computing: section (%v,%v) %v\n", sec.Name, sec.hasVariableReferences, context)
	var workingContext interface{}
	if sec.context == nil {
		workingContext = context
	} else {
		workingContext = sec.context
	}
	sec.executing = true

	//fmt.Fprintf(os.Stderr, "\nbefore section compute: %v %v %v\n", sec.params, sec.values, context)
	sec.DetermineResults(workingContext)

	//fmt.Fprintf(os.Stderr, "\n========> after section compute: %v\n", context)
	sec.executing = false

	switch wiper := context.(type) {
	case wipefuncfromStore:
		{
			wiper.EraseValueFromStore(sec.Name)
		}
	}

	//fmt.Printf("after section compute %v\n", context)
}

func (sec *DimSectionDeclaration) AddToStatements(value interface{}) {
	if sec.Statements != nil {
		sec.Statements.PushFront(value)
	}
}

func (sec *DimSectionDeclaration) Results() map[string]float64 {
	//fmt.Fprintf(os.Stderr, "section results called\n")
	if sec.acc != nil {
		temp := sec.acc
		sec.acc = nil
		return temp
	}
	return nil
}

func (sec *DimSectionDeclaration) String() string {
	//fmt.Fprintf(os.Stderr, "\nSection(%v):\n\tstatement:%v\n\tparams:%v\n\tparams-v:%v\n", sec.Name, sec.Statements, sec.params, sec.values)
	return fmt.Sprintf("Section(%v):%v\n", sec.Name, sec.context)
}
