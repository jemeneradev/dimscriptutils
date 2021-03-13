package measurement_factors

import (
	list "container/list"
	"fmt"
)

/*
DimVector represents a measurement component
*/
type DimVector struct {
	Measurements      *list.List
	Characteristics   *list.List
	Multiplier        *DimMeasurement
	ProcessFilterList *list.List
	Filters           []string
}

/*
NewDimVector Constructor
*/
func NewDimVector(mlist interface{}, clist interface{}, mul interface{}, flist interface{}) *DimVector {
	nv := new(DimVector)

	if mlist != nil {
		switch m := mlist.(type) {
		case *list.List:
			{
				nv.Measurements = m
			}
		}
	}

	if clist != nil {
		switch c := clist.(type) {
		case *list.List:
			{
				nv.Characteristics = c
			}
		}
	}

	if mul != nil {
		switch ml := mul.(type) {
		case *DimMeasurement:
			{
				nv.Multiplier = ml
			}
		}
	}

	if flist != nil {
		switch fl := flist.(type) {
		case []string:
			{
				nv.Filters = fl
			}
		}
	}

	return nv
}

type workingContext interface {
	IsCountMutable() bool
	IncrementCount(float64)
	GetHandler(string) interface{}
	LookUpVar(string) interface{}
}

type loadableToMap interface {
	LoadToMap(map[string]interface{})
}

type metricRepresentor interface {
	ToInches() float64
}

func mapunion(l map[string]float64, r map[string]float64) {
	for k, v := range r {
		_, ok := l[k]
		if ok {
			l[k] += v
		} else {
			l[k] = v
		}
	}
}

/*
String debug output
*/
func (dv *DimVector) String() string {
	return fmt.Sprintf("<Vector:[%v-%v-%v-%v]>", EncodeList(dv.Measurements), EncodeList(dv.Characteristics), dv.Filters, dv.Multiplier)
}

/*
GetCharacteristicMap return vector characteristic map
*/
func GetCharacteristicMap(ref interface{}) map[string]interface{} {
	switch l := ref.(type) {
	case *list.List:
		{
			if l != nil {
				cmap := make(map[string]interface{})
				for e := l.Front(); e != nil; e = e.Next() {
					//fmt.Printf("\nelem:%+v\n", e.Value)
					switch elem := e.Value.(type) {
					case loadableToMap:
						{
							//fmt.Printf("load\n")
							elem.LoadToMap(cmap)
						}
					}
				}
				//fmt.Printf("map:%v\n", cmap)
				return cmap
			}
		}
	}
	return nil
}

/*
DetermineResults solve current measurement vector
*/
func (dv *DimVector) DetermineResults(context interface{}) interface{} {
	//fmt.Fprintf(os.Stderr, "Determining vector: %v\n", dv)
	switch contextValue := context.(type) {
	case workingContext:
		multiplierScalar := float64(1)
		multiplierScalar = dv.Multiplier.ToNumber()
		if contextValue.IsCountMutable() {
			contextValue.IncrementCount(multiplierScalar)
		}

		chmap := GetCharacteristicMap(dv.Characteristics)

		var mSlice []float64
		if dv.Measurements != nil {
			mSlice = make([]float64, dv.Measurements.Len())
			i := 0
			for e := dv.Measurements.Front(); e != nil; e = e.Next() {
				//fmt.Fprintf(os.Stderr, "e: %v\n", e.Value)
				switch c := e.Value.(type) {
				case metricRepresentor:
					mSlice[i] = c.ToInches()
				}
				i++
			}
		}
		var results map[string]float64
		//No filter given, so apply default dim
		if dv.Filters == nil {
			fref := contextValue.GetHandler("dim")
			filter, ok := fref.(func([]float64, map[string]interface{}) map[string]float64)
			if ok {
				results = filter(mSlice, chmap)
			}
		} else {
			for _, fi := range dv.Filters {
				//fmt.Printf("%v--,%v\n", fi, mSlice)
				fref := contextValue.GetHandler(fi)
				//TODO:handle unknown filters
				//TODO:consider lang for filter list substituion. ex: @sum@count@divide = @avg
				if fref != nil {
					switch filter := fref.(type) {
					case func([]float64, map[string]interface{}) map[string]float64:
						{
							if results == nil {
								results = filter(mSlice, chmap)
							} else {
								//mapunion(chmap, results)
								mapunion(results, filter(mSlice, chmap)) //TODO: should I only use mapunion? maybe allow to have access to past results to future filters
							}

						}
					}
				}
				//fmt.Printf("%v\n", results)
				//else throw error
			}
		}
		if multiplierScalar != 1 {
			//fmt.Printf("mul%v", multiplierScalar)
			for k := range results {
				results[k] *= multiplierScalar
			}
		}

		return results
	}

	return nil
}
