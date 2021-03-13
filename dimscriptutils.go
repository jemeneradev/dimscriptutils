package dimscriptutils

import "fmt"

func Hello() []byte {
	fmt.Printf("%v", "hi")
	return []byte{'J', 'e', 's', 'u', 's'}
}
