package main

import (
	"flag"
	"fmt"
	"log"
	"os"
)

func main() {
	path := flag.String("C", "", "template's directory (by default is ~/"+GetDefaultTemplateDir()+")")
	outputFileName := flag.String("o", "", "output file name")
	defaultName := flag.Bool("d", false, "use default name for template")

	flag.Parse()
	names := flag.Args()

	storage, err := NewStorage(*path)
	if err != nil {
		log.Fatal(err)
	}

	if len(names) > 0 {
		var file *os.File = os.Stdout

		if !(*defaultName) && *outputFileName != "" {
			file, err = os.Create(*outputFileName)
			if err != nil {
				log.Fatal(err)
			}
			defer file.Close()
		}

		for _, name := range names {
			templateFile, ok := storage.templates[name]
			if !ok {
				log.Println("Can't find template", name)
				continue
			}

			content, err := os.ReadFile(templateFile.Path)
			if err != nil {
				log.Fatal(err)
			}

			template := NewTemplate(&templateFile, content)
			values := make(map[string]string)
			for _, variable := range template.Variables {
				values[variable] = os.Getenv(variable)
			}
			if *defaultName {
				if template.Filename == "" {
					log.Fatalf("Template %s has no default name, but -d flag was set", template.Name)
				}
				file, err = os.Create(template.Filename)
				if err != nil {
					log.Fatal(err)
				}
				defer file.Close()
			}
			fmt.Fprint(file, template.fillTemplate(values))
		}
		return
	}

	// if no args, print all template names
	for _, templateFile := range storage.templates {
		fmt.Println(templateFile)
	}
}
