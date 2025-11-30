package main

import (
	"fmt"
	"os"
)

func main() {
	var state State
	if len(os.Args) > 2 && os.Args[1] == "-C" {
		state.SetTemplateDir(os.Args[2])
	} else {
		state.DefaultTemplateDir()
	}

	for _, template := range state.Templates() {
		fmt.Println(template)
	}
}
