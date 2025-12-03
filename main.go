package main

import (
	"flag"
	"fmt"
	"os"
)

func main() {
	path := flag.String("C", "", "template's directory")
	flag.Parse()

	state := NewState()
	state.SetTemplateDir(*path)

	templates := state.Templates()
	requestedTemplates := flag.Args()
	if len(requestedTemplates) > 0 {
		for _, request := range requestedTemplates {
			template, ok := templates[request]
			if ok {
				// todo print real template with substituion
				fmt.Println(template)
			} else {
				fmt.Fprintf(os.Stderr, "Template %s not found\n", request)
			}
		}
		return
	}

	for _, template := range templates {
		fmt.Println(template)
	}
}
