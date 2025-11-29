package main

import (
	"fmt"
	"os"
)

func main() {
	var state *State
	if len(os.Args) > 2 && os.Args[1] == "-C" {
		state = NewState(os.Args[2])
	} else {
		state = NewStateWithDefaultDir()
	}

	for _, template := range state.Templates() {
		fmt.Println(template)
	}
}
