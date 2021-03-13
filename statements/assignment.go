package statements

import (
	"container/list"
	"fmt"
	"regexp"

	"github.com/jemeneradev/dimscriptutils/evaluator"
)

var varnameRegex = regexp.MustCompile(`(?m)((?P<dec_var>\$))*(?P<name>[^-=\n]+)(?P<options>(?P<delim>((-)|(=>){1}))(?P<val>.+))*`)
var optionReg = regexp.MustCompile(`(?P<option>[a-x]){1}(?P<params>\((?P<attr>[^)]+)\))?`)
var splitter = regexp.MustCompile(`,`)

type DimAssignment struct {
	Name    string
	Dimlist interface{}
}

func NewDimAssignment(name string, dimlist interface{}) *DimAssignment {
	assignment := new(DimAssignment)
	assignment.Name = name
	assignment.Dimlist = dimlist
	return assignment
}

func (assignment *DimAssignment) Identifier() string {
	return assignment.Name
}

type identifier interface {
	Identifier() string
}
type storeable interface {
	InsertIntoStore(string, interface{})
}
type countMutable interface {
	DisableCountMutable()
	EnableCountMutable()
}

func (assignment *DimAssignment) DetermineResults(context interface{}) interface{} {
	//fmt.Fprintf(os.Stderr, "\nDetermining assignment results %T %v\n", assignment.Dimlist, context)
	switch l := assignment.Dimlist.(type) {
	case *list.List:
		{
			//fmt.Fprintf(os.Stderr, "val found: %v %v\n", l.Front(), assignment)
			if l.Front() != nil {
				value := l.Front().Value
				match := varnameRegex.FindStringSubmatch(assignment.Name)
				includeCount := true
				//fmt.Printf("value -> %v %v\n", value, match)
				//fmt.Printf("\nexamine and combine results:\n%v\t:\t%v\n", assignment.Name, assignment.Dimlist)
				switch match[5] {
				case "-":
					{
						//fmt.Printf("apply simple option:<%v>\n", match[9])
						for _, ops := range splitter.Split(match[9], -1) {
							if optionReg.FindStringSubmatch(ops)[0] == "x" {
								includeCount = false
							}
							//fmt.Printf("here:%+v\n", optionReg.FindStringSubmatch(ops))
						}
					}
				case "=>":
					{
						fmt.Printf("apply assignment option:<%v>\n", match[9])
					}
				}
				//fmt.Fprintf(os.Stderr, "content in assignment:%v %T\n", context, context)
				var results map[string]float64
				if includeCount == true {
					results = evaluator.Evaluate(value, context)
				} else {
					switch countEnabler := context.(type) {
					case countMutable:
						{
							countEnabler.DisableCountMutable()
							results = evaluator.Evaluate(value, context)
							countEnabler.EnableCountMutable()
						}
					}
				}
				//fmt.Printf("=>result: %v %v\n", results, value)
				//match[1] declare symbol
				//match[3] assignment name
				//match[5] option delim
				//match[9] option value
				//fmt.Printf("Line Results:%v\n", results)
				if match[1] == "$" {
					//fmt.Printf("apply var logic: %+v\n", context)
					switch store := context.(type) {
					case storeable:
						{
							store.InsertIntoStore(match[3], results)
							//fmt.Fprintf(os.Stderr, "store: %v", store)
						}
					}
				}
				/* //fmt.Printf("spot %v %v\n", match, len(match))

				//fmt.Printf("End of Assignment\n") */
				return results
			}
		}
	}
	//fmt.Fprintf(os.Stderr, "results:%v\n", evaluator.Evaluate(assignment.Dimlist.Front().Value, context))
	return nil
}
