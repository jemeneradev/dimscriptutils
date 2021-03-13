package evaluator

import "fmt"

/*
NewDimensionHandlers return default dim,sum,profile component functions as a map
*/
func NewDimensionHandlers() *map[string]interface{} {

	m := make(map[string]interface{})
	m["dim"] = func(dimensions []float64, characteristicMap map[string]interface{}) map[string]float64 {
		results := make(map[string]float64)
		//fmt.Printf("chars: %v", characteristicMap)
		switch len(dimensions) {
		case 0:
			{
				results["Area"] = 0
				_, lengthIsPresent := characteristicMap["l"]
				if lengthIsPresent {
					results["Length"] = float64(0)
				}
				_, widthIsPresent := characteristicMap["w"]
				if widthIsPresent {
					results["Width"] = float64(0)
				}
			}
		case 1:
			{
				results["Length"] = dimensions[0]
			}
		default:
			{
				results["Area"] = dimensions[0] * dimensions[1]

				_, lengthIsPresent := characteristicMap["l"]
				if lengthIsPresent {
					results["Length"] = dimensions[0]
				}
				_, widthIsPresent := characteristicMap["w"]
				if widthIsPresent {
					results["Width"] = dimensions[1]
				}
			}
		}
		return results
	}
	m["sum"] = func(dimensions []float64, characteristicMap map[string]interface{}) map[string]float64 {

		results := make(map[string]float64)
		sum := float64(0)
		for _, x := range dimensions {
			sum += x
		}
		results["Sum"] = sum
		dimLength := float64(len(dimensions))

		_, lenIsPresent := characteristicMap["len"]
		if lenIsPresent {
			results["Len"] = float64(dimLength)
		}

		_, avgIsPresent := characteristicMap["avg"]
		if avgIsPresent {
			if dimLength > 0 {
				results["Avg"] = float64(sum / dimLength)
			} else {
				results["Avg"] = float64(0)
			}
		}

		return results
	}
	m["profile"] = func(dimensions []float64, characteristicMap map[string]interface{}) map[string]float64 {
		results := make(map[string]float64)
		//length := 0
		var startIndex int
		sectionMeasurement := float64(0)
		if characteristicMap["xsection"] != nil {
			startIndex = 0
			sectionMeasurement = float64(0) //characteristicMap["xsection"]
			switch v := characteristicMap["xsection"].(type) {
			case []float64:
				{
					if len(v) > 0 {
						sectionMeasurement = v[0]
					}
				}
			}
		} else {
			startIndex = 1
			sectionMeasurement = 0
			if len(dimensions) > 0 {
				sectionMeasurement = dimensions[0]
			}

		}
		sum := float64(0)
		for i := startIndex; i < len(dimensions); i++ {
			sum += dimensions[i]
		}

		results["Area"] = sum * sectionMeasurement
		results["Length"] = sum
		return results
	}
	m["frame"] = func(dimensions []float64, characteristicMap map[string]interface{}) map[string]float64 {
		results := make(map[string]float64)
		var height = float64(0)
		var width = float64(0)
		var depth = float64(0)

		switch len(dimensions) {
		case 0: //this should never happen, but include anyway
			{
				height = float64(12)
				width = float64(12)
				depth = float64(.75)
			}
		case 1:
			{
				height = dimensions[0]
				width = float64(12)
				depth = float64(.75)
			}
		case 2:
			{
				height = dimensions[0]
				width = dimensions[1]
				depth = float64(.75)
			}
		default:
			{
				height = dimensions[0]
				width = dimensions[1]
				depth = dimensions[2]
			}
		}

		offsets := [4]float64{0, 0, 0, 0}
		offsetsValue, offsetsIsPresent := characteristicMap["offsets"]
		if offsetsIsPresent {
			switch offsets_args := offsetsValue.(type) {
			case []float64:
				{
					switch len(offsets_args) {
					case 1:
						{
							offsets[0] = offsets_args[0]
							offsets[1] = offsets_args[0]
							offsets[2] = offsets_args[0]
							offsets[3] = offsets_args[0]
						}
					case 2:
						{
							offsets[0] = offsets_args[0]
							offsets[1] = offsets_args[1]
							offsets[2] = offsets_args[0]
							offsets[3] = offsets_args[1]
						}
					case 4:
						{
							offsets[0] = offsets_args[0]
							offsets[1] = offsets_args[1]
							offsets[2] = offsets_args[2]
							offsets[3] = offsets_args[3]
						}
					default:
						{
							offsets[0] = float64(1)
							offsets[1] = offsets[0]
							offsets[2] = offsets[0]
							offsets[3] = offsets[0]
						}
					}
				}
			}
		} else {
			offsets[0] = float64(1)
			offsets[1] = offsets[0]
			offsets[2] = offsets[0]
			offsets[3] = offsets[0]
		}

		faceNumber := float64(1)
		faceValue, faceIsPresent := characteristicMap["face"]
		if faceIsPresent {
			switch face_args := faceValue.(type) {
			case []float64:
				{
					if len(face_args) > 0 {
						faceNumber = face_args[0]
					}
				}
			}
		}

		//offset[0] top
		//offset[1] right
		//offset[2] bottom
		//offset[3] left

		heightPlusOffsets := height + offsets[0] + offsets[2]
		widthPlusOffsets := width + offsets[1] + offsets[3]
		face := heightPlusOffsets*widthPlusOffsets - width*height

		fmt.Printf("face:%v", face)
		cir := 2 * depth * (heightPlusOffsets + widthPlusOffsets + height + width)
		results["Area"] = (faceNumber * face) + cir
		return results
	}
	return &m
}
