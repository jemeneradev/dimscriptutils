package measurement_factors

import (
	list "container/list"
	"fmt"
	"os"
	"strings"

	"github.com/jemeneradev/dimscriptutils/constants"
)

type Encoder interface {
	Encode() string
}

func DimMeasurementDecode(dms string) *DimMeasurement {
	var ts string
	var ti, tp int
	_, err := fmt.Sscanf(dms, "<dimM{ %s %d %d }>", &ts, &ti, &tp)
	if err != nil {
		//panic(err)
		return nil
	} else {
		//fmt.Printf("%s %d %d\n",ts,ti,tp)
		return NewDimMeasurement(ts, ti, tp)
	}
}

func NewMeasurementList() *list.List {
	return list.New()
}

func EncodeList(ref interface{}) string {
	switch l := ref.(type) {
	case *list.List:
		{
			if l != nil {
				var buf strings.Builder
				for e := l.Front(); e != nil; e = e.Next() {
					var i interface{} = e.Value
					switch s := i.(type) {
					case *DimMeasurement:
						{
							fmt.Fprintf(&buf, "%s,", s.Encode())
						}
					case *DimCharacteristic:
						{
							fmt.Fprintf(&buf, "%s,", s.Encode())
						}
					}
				}
				//fmt.Printf("EncodingList: %s\n",buf.String())
				return buf.String()
			}
		}
	}

	return ""
}

type maploader interface {
	LoadToMap(m map[string]interface{})
}

func DecodeList(dms string) *list.List {
	if dms != "" {
		//fmt.Printf("str to Decode: %s\n",dms)
		temp := NewMeasurementList()
		if temp != nil {
			for _, el := range strings.Split(dms, ",") {
				if el != "" {
					temp.PushBack(DimMeasurementDecode(el))
				}
			}
		}
		//fmt.Printf("DecodedList: %s\n",EncodeList(temp))
		return temp
	}
	return nil
}

func SaveMeasurementInList(ref interface{}, el *DimMeasurement) {
	switch l := ref.(type) {
	case *list.List:
		{
			l.PushBack(el)
		}
	}
}
func SaveCharacteristicInList(ref interface{}, el *DimCharacteristic) {
	switch l := ref.(type) {
	case *list.List:
		{
			l.PushBack(el)
		}
	}
}

func (dm *DimCharacteristic) Encode() string {
	if dm != nil {
		return fmt.Sprintf("<dimC{ %s;%s }dimC>", dm.Name, EncodeList(dm.Args))
	} else {
		return ""
	}
}

func DimCharacteristicDecode(dms string) *DimCharacteristic {
	if dms != "" {
		r := strings.NewReplacer(" }dimC>", "", "<dimC{ ", "")
		dimc_els := strings.Split(r.Replace(dms), ";")
		return NewDimCharacteristic(dimc_els[0], DecodeList(dimc_els[1]))
	}
	return nil
}

func NewDimCharacteristic(name interface{}, lst interface{}) *DimCharacteristic {
	nD := new(DimCharacteristic)
	switch n := name.(type) {
	case string:
		{
			nD.Name = n
		}
	}
	if lst != nil {
		switch l := lst.(type) {
		case *list.List:
			{
				nD.Args = l
			}
		}
	}
	return nD
}

func MoveElementsOver(left interface{}, right interface{}) {
	lList, okL := left.(*list.List)
	rList, okR := right.(*list.List)
	//fmt.Printf("%T %T\n", left, right)
	if okL && okR {
		for rList.Len() > 0 {
			lList.PushBack(rList.Remove(rList.Front()))
		}
		right = nil
	}
}

type mappable interface {
	ToComponentMap() map[string]interface{}
}

type hasvarrefs interface {
	SetHasVariableReferences(b bool)
}

type sectionLookup interface {
	LookUpVar(v string) interface{}
}

func DimFuncReferenceCall(name string, parameters interface{}, sec interface{}) interface{} {
	switch top := sec.(type) {
	case *list.List:
		{
			sectionTop := top.Front().Value
			boolSetter, ok := sectionTop.(hasvarrefs)
			if ok {
				boolSetter.SetHasVariableReferences(true) //indicate that current section has a var
			}
			nf := NewDimFunctionCall(name, parameters.(*list.List), sec)
			//fmt.Fprintf(os.Stderr, "func %v %T\n", nf, nf)
			return nf
		}
	}
	return nil
}

func DimVarReferenceCall(name string, sec interface{}) interface{} {
	switch top := sec.(type) {
	case *list.List:
		{
			sectionTop := top.Front().Value
			boolSetter, ok := sectionTop.(hasvarrefs)
			if ok {
				boolSetter.SetHasVariableReferences(true) //indicate that current section has a var
			}
			return NewDimMeasurement(name, constants.DimReference, constants.DimVarOperand)
		}
	}
	return nil
}

func DimCharacteristicAsVariable(clist interface{}, sec interface{}) interface{} {

	switch d := clist.(type) {
	case *list.List:
		{
			//fmt.Fprintf(os.Stderr, "%v %v", d.Front().Value, sec)
			switch d.Len() {
			case 1:
				{
					cm, ok := d.Front().Value.(mappable)
					if ok {

						switch top := sec.(type) {
						case *list.List:
							{
								sectionTop := top.Front().Value
								boolSetter, ok := sectionTop.(hasvarrefs)
								if ok {
									boolSetter.SetHasVariableReferences(true) //indicate that current section has a var
								}

								//slookup, canlookup := sectionTop.(sectionLookup)

								//if canlookup {
								refmap := cm.ToComponentMap()
								args, found := refmap["args"]
								if found == true {
									fname, fok := refmap["name"]
									if fok {
										//if slookup.LookUpVar(fname.(string)) != nil {
										nf := NewDimFunctionCall(fname.(string), args.(*list.List), sec)
										//fmt.Fprintf(os.Stderr, "func %v %T\n", nf, nf)
										return nf
										//}
										//return nil
									}
								} else {
									vname, vok := refmap["name"]
									if vok {
										sname, sok := vname.(string)
										if sok {
											//fmt.Fprintf(os.Stderr, "var %v\n", sname)
											return NewDimMeasurement(sname, constants.DimReference, constants.DimVarOperand)
										}
									}

								}
								//}

							}
						}
					}
				}
			default:
				{
					fmt.Fprintf(os.Stderr, "handle multiple char error")
					return nil
				}
			}
		}
	}
	return nil
}
