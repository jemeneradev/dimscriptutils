package evaluator

import (
	"testing"
)

func TestDimContextString(t *testing.T) {
	emptyDimContext := NewDimContext(NewDimensionHandlers(), nil)
	want := `{"Totals":{"Count":0}}`
	if got := emptyDimContext.String(); got != want {
		t.Errorf("Context: got %v, want %v", got, want)
	}
}

func TestDimContextAccumulatorHandling(t *testing.T) {
	dimContext := NewDimContext(NewDimensionHandlers(), nil)

	//dimContext.StoreAccumulator("main")
	mainSection := make(map[string]float64)
	//add to stack
	dimContext.AccumulatePush(mainSection)

	results := make(map[string]float64)
	results["Area"] = float64(5)
	results["Sum"] = float64(9)

	dimContext.Accumulate(results)
	dimContext.StoreAccumulator("main")
	/* want := "{\"Totals\":{\"Area\":2,\"Count\":0,\"Sum\":8}}"
	if got := dimContext.String(); got != want {
		t.Errorf("Context: got %v, want %v", got, want)
	} */

	want := `{"Totals":{"Area":5,"Count":0,"Sum":9}}`
	if got := dimContext.String(); got != want {
		t.Errorf("Context: got %v, want %v", got, want)
	}

	dimContext.Accumulate(results)
	dimContext.StoreAccumulator("main")
	/* want := "{\"Totals\":{\"Area\":2,\"Count\":0,\"Sum\":8}}"
	if got := dimContext.String(); got != want {
		t.Errorf("Context: got %v, want %v", got, want)
	} */

	want = `{"Totals":{"Area":10,"Count":0,"Sum":18}}`
	if got := dimContext.String(); got != want {
		t.Errorf("Context: got %v, want %v", got, want)
	}

}
