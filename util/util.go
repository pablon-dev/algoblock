package util

import (
	"encoding/json"
	"fmt"
)

//Pretty print look JSON
func Pretty(data interface{}) {
	p, err := json.MarshalIndent(data, "", "\t")
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Printf("%s \n", p)
}
