package main

import (
	"flag"
	"log"
)

func main() {
	path := flag.String("C", "", "template's directory (by default is ~/"+GetDefaultTemplateDir()+")")
	outputFileName := flag.String("o", "", "output file name")
	defaultName := flag.Bool("d", false, "use default name for template")

	flag.Parse()

	storage, err := NewStorage(*path)
	if err != nil {
		log.Fatal(err)
	}

	runState := CliState{
		*outputFileName,
		*defaultName,
		flag.Args(),
		storage,
	}

	if err := runState.Run(); err != nil {
		log.Fatal(err)
	}
}
