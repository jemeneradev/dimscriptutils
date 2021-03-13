package measurement_factors

import (
	list "container/list"
)

//import "fmt"

//DimCharacteristic is a place holder to further describe
//measurements. Can serve to condition measurement computation.
//example: 1',2',avg()@sum
type DimCharacteristic struct {
	Name string
	Args *list.List
}

//NewCharacteristicList returns a initialized list.
//Mask off representation.
func NewCharacteristicList() *list.List {
	return list.New()
}

type numRep interface {
	ToNumber() float64
}

//LoadToMap turns a list into a map representation for ease of use.
//Instead of making user traverse a list, more error prone, just look up a map.
func (dc *DimCharacteristic) LoadToMap(m map[string]interface{}) {
	if dc.Args != nil {
		characteristicArgs := make([]float64, dc.Args.Len())
		i := 0
		for e := dc.Args.Front(); e != nil; e = e.Next() {
			switch n := e.Value.(type) {
			case numRep:
				characteristicArgs[i] = n.ToNumber()
			}
			i++
		}
		m[dc.Name] = characteristicArgs
	} else {
		m[dc.Name] = nil
	}
}

func (dc *DimCharacteristic) ToComponentMap() map[string]interface{} {
	components := make(map[string]interface{})
	components["name"] = dc.Name
	if dc.Args != nil {
		components["args"] = dc.Args
	}
	return components
}
