package main

import (
	"flag"
	"log"
)

func main() {
	path := flag.String("C", "", "template's directory (by default is ~/"+GetDefaultTemplateDir()+")")
	outputFileName := flag.String("o", "", "output file name")
	defaultName := flag.Bool("d", false, "use default name for template")
	noInput := flag.Bool("no-input", false, "use only environment variables")
	editMode := flag.Bool("edit", false, "edit selected template in your console editor")
	listTemplates := flag.Bool("l", false, "list all templates")

	flag.Parse()

	ioh := DefaultIOHandler()

	storage, err := NewStorage(ioh, *path)
	if err != nil {
		log.Fatal(err)
	}

	runState := CliState{
		*outputFileName,
		*defaultName,
		flag.Args(),
		storage,
		*noInput,
		ioh,
		*editMode,
		*listTemplates,
	}

	if err := runState.Run(); err != nil {
		log.Fatal(err)
	}
}
