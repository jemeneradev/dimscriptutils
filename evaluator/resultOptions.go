package evaluator

/*
ResultOutputOptions holds output format(json,csv,table,etc) and instruction strings (expand later).
*/
type ResultOutputOptions struct {
	responseType string
	//option:	include	all
	//					only <section list>
	//					except <sections to exclude>
	//TODO: 			exec <execute calculation where section resultsare referenced as variables>
	includeInResponce string
}

/*
NewResultOutputOptions Constructor
*/
func NewResultOutputOptions(passedInOptions interface{}) *ResultOutputOptions {
	resultOptions := new(ResultOutputOptions)
	if passedInOptions == nil {
		resultOptions.responseType = "json"
		resultOptions.includeInResponce = "all"
	} else {
		optionsMap, ok := passedInOptions.([]interface{})
		if ok {
			evalOptions, rok := optionsMap[0].(map[string]string)
			if rok {
				responseType, doesResponseTypeExist := evalOptions["outformat"]
				if doesResponseTypeExist {
					resultOptions.responseType = responseType
				} else {
					resultOptions.responseType = "json"
				}

				includeInResponce, doesincludeInResponceExist := evalOptions["include"]
				if doesincludeInResponceExist {
					resultOptions.includeInResponce = includeInResponce
				} else {
					resultOptions.includeInResponce = "all"
				}
			}

		}
	}
	return resultOptions
}
