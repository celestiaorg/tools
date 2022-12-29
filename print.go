package tools

import (
	"encoding/json"
	"fmt"

	"github.com/davecgh/go-spew/spew"
)

// useful for debugging
func PrettyPrint(v any, prefix ...string) {

	if len(prefix) > 0 {
		fmt.Printf("%s: ", prefix)
	}

	spew.Config.MaxDepth = 1
	spew.Config.Indent = "\t"
	spew.Dump(v)
}

func PrintJson(v any, prefix ...string) {

	if len(prefix) > 0 {
		fmt.Printf("%s: ", prefix)
	}

	b, _ := json.MarshalIndent(v, "", "\t")
	fmt.Printf("\n%s\n", b)
}
