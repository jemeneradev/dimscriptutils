package evaluator

import (
	"container/list"
	"encoding/json"
	"fmt"
	"os"
	"regexp"
	"strings"
)

type evaluatorFunction (func(interface{}, interface{}) map[string]float64)

/*
DimContext context object holding global vars
*/
type DimContext struct {
	Count       float64 //holds project count
	MutateCount bool    //should count be changed in current calculation
	Handlers    *map[string]interface{}
	Table       map[string]interface{} //holds results for each section defined
	Store       map[string]interface{} //holds section vars

	Accumulator []map[string]float64 //Holds scope object, use to save values
	accStackTop int                  //where to store next. accStackTop - 1: working base

	outPutOptions *ResultOutputOptions
	eval          evaluatorFunction

	lookup []string

	sectionIndex map[string]interface{}
}

/*
NewDimContext Constructor
*/
func NewDimContext(handlers *map[string]interface{}, passedInOptions ...interface{}) *DimContext {
	n := new(DimContext)
	n.Count = 0
	n.MutateCount = true
	n.Handlers = handlers
	n.Accumulator = make([]map[string]float64, 3) //TODO:this limits section declarations within section
	n.accStackTop = 0
	n.Table = make(map[string]interface{})
	n.Store = make(map[string]interface{})
	n.sectionIndex = make(map[string]interface{})
	n.lookup = make([]string, 0)
	if passedInOptions != nil {
		n.outPutOptions = NewResultOutputOptions(passedInOptions)
	} else {
		n.outPutOptions = NewResultOutputOptions(nil)
	}
	return n
}

type statementGetter interface {
	GetStatementAt(mthd string) interface{}
}

type parameterSetter interface {
	SetParamsValues(passedinValues interface{})
}

type hasDeterminResults interface {
	DetermineResults(context interface{}) interface{}
}

func (dc *DimContext) EvaluateSectionAt(sectionName string, methodName string, args interface{}, cntx interface{}) interface{} {

	section, ok := dc.sectionIndex[sectionName].(statementGetter)
	if ok {
		funcRef := section.GetStatementAt(methodName)
		if funcRef == nil {
			//!handle section not defined
			return nil
		}
		//fmt.Fprintf(os.Stderr, "Evaluate: %v[%v]:%v %v\n", sectionName, methodName, args, cntx)
		switch ps := funcRef.(type) {
		case parameterSetter:
			{
				//fmt.Fprintf(os.Stderr, "\t\tfunction state: %v \t\twith params: %v %T\n", funcRef, args, args)
				switch plist := args.(type) {
				case *list.List:
					{
						extractedParameters := make(map[int]interface{}, plist.Len())
						i := 0
						for e := plist.Front(); e != nil; e = e.Next() {
							//fmt.Fprintf(os.Stderr, "e.Valuetype:%T", e.Value)
							switch valNode := e.Value.(type) {
							case hasDeterminResults:
								{
									extractedParameters[i] = valNode.DetermineResults(cntx)
								}
							}
							i++
						}
						ps.SetParamsValues(extractedParameters)
						//fmt.Fprintf(os.Stderr, "\n\t\tbefore func eval %v %T %T\n", funcRef, funcRef, cntx)
						//fmt.Fprintf(os.Stderr, "content in Evaluate:%v %T\n", cntx, cntx)
						//dc.PushLookUp(methodName)
						results := dc.eval(funcRef, cntx)
						//fmt.Fprintf(os.Stderr, "func results : %v \n", results)
						//dc.PopLookUp()
						return results
					}
				}
			}
		}

		/* for p := range args {
			//fmt.Fprintf(os.Stderr, "parameter: %v\n", p)
		} */
		////fmt.Fprintf(os.Stderr, "Evaluate function:%v at section(%v)\n -- ref%v\n", f, methodName, funcOp)
		//dc.eval(funcOp, dc)
	}
	return nil
}

func (dc *DimContext) PushLookUp(n string) {
	dc.lookup = append(dc.lookup, n)
	//fmt.Fprintf(os.Stderr, "lookup push:%v %v\n", dc.lookup, len(dc.lookup))
}

func (dc *DimContext) PopLookUp() {
	lookupDepth := len(dc.lookup)
	dc.lookup = dc.lookup[:lookupDepth-1]
	//fmt.Fprintf(os.Stderr, "lookup pop:%v\n", dc.lookup)
}

func (dc *DimContext) LookUpVar(v string) interface{} {
	fmt.Fprintf(os.Stderr, "looking for %v\n", v)
	lookupLen := len(dc.lookup) - 1
	var section string
	for i := 0; i <= lookupLen; i++ {
		//fmt.Fprintf(os.Stderr, "%v\n", lookupLen-i)
		section = dc.lookup[lookupLen-i] //reverse index, last to first
		sectionStore, contextHasSection := dc.Store[section]
		//fmt.Fprintf(os.Stderr, "\t%v - %v\n", section, dc.Store)
		if contextHasSection {
			indexable, ok := sectionStore.(map[string]interface{})
			if ok {
				val, sectionStoreHasVar := indexable[v]
				if sectionStoreHasVar {
					return val
				}
			}

		}
	}
	//fmt.Fprintf(os.Stderr, "\tlook array: %v %v\n", dc.lookup, len(dc.lookup))
	/* for _, section := range dc.lookup {
		sectionStore, contextHasSection := dc.Store[section]
		//fmt.Fprintf(os.Stderr, "\t%v - %v\n", section, dc.Store)
		if contextHasSection {
			indexable, ok := sectionStore.(map[string]interface{})
			if ok {
				val, sectionStoreHasVar := indexable[v]
				if sectionStoreHasVar {
					return val
				}
			}

		}

	} */
	return nil
}

func (dc *DimContext) AddSectionToIndex(n string, s interface{}) {
	dc.sectionIndex[n] = s
	//fmt.Fprintf(os.Stderr, "sectionIndex: %v\n", dc.sectionIndex)
}

func (dc *DimContext) SetEval(e evaluatorFunction) {
	dc.eval = e
}

/*
AccumulatePush DimContext stack interface push
*/
func (dc *DimContext) AccumulatePush(val interface{}) {
	switch fv := val.(type) {
	case map[string]float64:
		{
			//fmt.Fprintf(os.Stderr, "\n\t\tbefore acc: %v %v %v\n", dc.Accumulator, dc.accStackTop, dc.lookup)
			dc.Accumulator[dc.accStackTop] = fv
			dc.accStackTop++
			//fmt.Fprintf(os.Stderr, "\n\t\tafter acc: %v %v %v\n", dc.Accumulator, dc.accStackTop, dc.lookup)
		}
	}

}

/*
AccumulatePop DimContext stack interface pop
*/
func (dc *DimContext) AccumulatePop() interface{} {
	tempAcc := dc.Accumulator[dc.accStackTop-1]
	if dc.accStackTop > 0 {
		dc.Accumulator[dc.accStackTop-1] = nil
		dc.accStackTop--
	}
	return tempAcc
}

type hasresults interface {
	Results() map[string]float64
}

/*
Accumulate DimContext stack interface peek top
*/
func (dc *DimContext) Accumulate(val interface{}) {
	//fmt.Printf("in con acc -> %v\n", val)

	if dc.accStackTop > 0 {
		switch fv := val.(type) {
		case hasResults:
			{
				for k, v := range fv.Results() {
					dc.Accumulator[dc.accStackTop-1][k] += v
					//fmt.Printf("k:%v,v:%v\n", k, v)
				}
			}
		case map[string]float64:
			{
				for k, v := range fv {
					dc.Accumulator[dc.accStackTop-1][k] += v
					//fmt.Printf("k:%v,v:%v\n", k, v)
				}
			}
		}
		//fmt.Fprintf(os.Stderr, "\nacctop:%v acc:%v lookup:%v\n", dc.accStackTop, dc.Accumulator, dc.lookup)
	}
}

/*
StoreAccumulator pop accumulator and save it in table
*/
func (dc *DimContext) StoreAccumulator(name string) {
	if len(dc.Accumulator) >= 0 && dc.accStackTop > 0 {
		//dc.Table[name] = dc.Accumulator[dc.accStackTop-1]
		dc.InsertIntoTable(name, dc.Accumulator[dc.accStackTop-1])
	}
}

/*
InsertIntoTable set key, val into table
*/
func (dc *DimContext) InsertIntoTable(key string, val interface{}) {
	if dc.Table == nil {
		dc.Table = make(map[string]interface{})
	}
	dc.Table[key] = val
}

/*
InsertIntoStore set key,val into store
*/
func (dc *DimContext) InsertIntoStore(key string, val interface{}) {
	if dc.Store != nil {
		i := len(dc.lookup)
		if i > 0 {
			i = i - 1
		}

		_, storeExist := dc.Store[dc.lookup[i]]
		if !storeExist {
			dc.Store[dc.lookup[i]] = make(map[string]interface{})
		}

		m, isMap := dc.Store[dc.lookup[i]].(map[string]interface{})
		if isMap {
			m[key] = val
		}

		//fmt.Fprintf(os.Stderr, "\n\033[38;2;255;255;0madding to store:\033[0m \n\tstore:%v\n\tlookup: %v\n\ti:%v\n\tlookup[i]:%v\n\tstore[lookup[i]]:%v\n\n", dc.Store, dc.lookup, i, dc.lookup[i], dc.Store[dc.lookup[i]])
	}

}

/*
RetrieveStore set store to nil
*/
func (dc *DimContext) RetrieveStore() interface{} {
	if dc.Store != nil {
		temp := dc.Store
		dc.Store = nil
		return temp
	}
	return nil
}
func (dc *DimContext) GetTableValue(k string) interface{} {
	if dc.Table != nil {
		temp, found := dc.Table[k]
		if found {
			delete(dc.Table, k)
			return temp
		}
	}
	return nil
}

func (dc *DimContext) EraseValueFromStore(k string) {
	if dc.Store != nil {
		_, found := dc.Store[k]
		if found {
			//fmt.Fprintf(os.Stderr, "looking for %v\n", k)
			delete(dc.Store, k)
		}
	}
}

/*
IsCountMutable count interface
*/
func (dc *DimContext) IsCountMutable() bool {
	return dc.MutateCount
}

func (dc *DimContext) DisableCountMutable() {
	dc.MutateCount = false
}
func (dc *DimContext) EnableCountMutable() {
	dc.MutateCount = true
}

/*
IncrementCount count interface
*/
func (dc *DimContext) IncrementCount(val float64) {
	dc.Count += val
}

/*
GetHandler handler interface
*/
func (dc *DimContext) GetHandler(k string) interface{} {
	return (*dc.Handlers)[k]
}

func (dc *DimContext) retriveSectionsToInclude() []string {
	//dc.StoreAccumulator("main") //save current accumulator as main in table
	////fmt.Printf("includeOption:%v\n", dc.outPutOptions.includeInResponce)
	includeCmdRegex := regexp.MustCompile(`^(?P<includeVerb>(only|except))(?P<rest>.+)`)
	includeResults := includeCmdRegex.FindStringSubmatch(dc.outPutOptions.includeInResponce)
	if len(includeResults) == 0 { //no include set was passed in or value incorrect format, fallback to default
		////fmt.Printf("%v", dc.Table)
		allSectionKeys := make([]string, 1)
		for k := range dc.Table {
			allSectionKeys = append(allSectionKeys, k)
		}
		////fmt.Printf("%v", allSectionKeys)
		return allSectionKeys
	} else {
		switch includeResults[2] {
		case "only":
			{
				multiSpaces := regexp.MustCompile(`\s+`)
				sectionToAdd := make([]string, 1)
				for _, s := range strings.Split(multiSpaces.ReplaceAllString(includeResults[3], " "), " ") {
					if s != "" {
						////fmt.Printf("add: <%v>", s)
						sectionToAdd = append(sectionToAdd, s)
					}
				}
				return sectionToAdd
			}
		case "except":
			{
				////fmt.Printf("sub")
				//Collect sections to exclude
				multiSpaces := regexp.MustCompile(`\s+`)
				sectionToSub := make(map[string]bool)
				for _, s := range strings.Split(multiSpaces.ReplaceAllString(includeResults[3], " "), " ") {
					if s != "" {
						////fmt.Printf("add: <%v>", s)
						sectionToSub[s] = true
					}
				}

				//get all section keys
				allSectionKeys := make([]string, 1)
				var ok bool
				for key := range dc.Table {
					_, ok = sectionToSub[key]
					if !ok {
						allSectionKeys = append(allSectionKeys, key)
					}
				}

				////fmt.Printf("%v - %v", allSectionKeys, sectionToSub)
				return allSectionKeys
			}
		}
	}
	return []string{""}
}

/*
MarshalJSON testing string for DimContext
*/
func (dc *DimContext) MarshalJSON() ([]byte, error) {
	////fmt.Printf("\nStore:%v \nTable:%v\noptions:%v\n", dc.Store, dc.Table, dc.outPutOptions)
	sectionTotals := make(map[string]float64)
	sectionEntries := dc.retriveSectionsToInclude()

	for _, section := range sectionEntries {
		if section != "" {
			////fmt.Fprintf(os.Stderr, "s: %v\n", dc.Table[section])
			switch sM := dc.Table[section].(type) {
			case map[string]float64:
				for sK, sV := range sM {
					sectionTotals[sK] += sV
				}
			}
		}
	}

	////fmt.Fprintf(os.Stderr, "\n==%v==\n", dc.accStackTop)
	/* for mK, mEntry := range dc.Accumulator[dc.accStackTop-1] {
		sectionTotals[mK] += mEntry
	} */

	results := make(map[string]map[string]float64)
	results["Totals"] = sectionTotals
	results["Totals"]["Count"] = float64(dc.Count)
	return json.Marshal(results)
}

func (dc *DimContext) String() string {
	/* fmt.Fprintf(os.Stderr, "\nmain table: %v\n", dc.Table)
	fmt.Fprintf(os.Stderr, "main store: %v\n", dc.Store)
	fmt.Fprintf(os.Stderr, "main lookup: %v\n", dc.lookup) */
	jr, _ := dc.MarshalJSON()
	return string(jr)
}
