package utils

import (
	"encoding/json"
	"fmt"
)

func PrintJson(v any) {
	indent, _ := json.MarshalIndent(v, "", "  ")
	fmt.Println(string(indent))
}
