package evaluator

import (
	"encoding/json"
	"testing"
)

func TestDimContext(t *testing.T) {
	dimContext := NewDimContext(NewDimensionHandlers(), nil)
	want := "{\"Totals\":{\"Count\":0}}"
	if got := dimContext.String(); got != want {
		t.Errorf("Context: got %v, want %v", got, want)
	}
}

func TestDimContextAccumulate(t *testing.T) {
	dimContext := NewDimContext(NewDimensionHandlers(), nil)
	dimContext.StoreAccumulator("main")
	results := make(map[string]float64)
	results["Area"] = float64(2)
	results["Sum"] = float64(8)
	dimContext.Accumulate(results)
	/* want := "{\"Totals\":{\"Area\":2,\"Count\":0,\"Sum\":8}}"
	if got := dimContext.String(); got != want {
		t.Errorf("Context: got %v, want %v", got, want)
	} */

	want := "{\"Area\":2,\"Sum\":8}"
	if got, _ := json.Marshal(dimContext.Table["main"]); string(got) != want {
		t.Errorf("Context: got %v, want %v", string(got), want)
	}

	dimContext.Accumulate(results)
	want = "{\"Area\":4,\"Sum\":16}"
	if got, _ := json.Marshal(dimContext.Table["main"]); string(got) != want {
		t.Errorf("Context: got %v, want %v", string(got), want)
	}

}
