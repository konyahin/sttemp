package main

import (
	"fmt"
	"os"
)

// CliState represents the state of the running app with all options set
// and keep all business logic inside his functions
type CliState struct {
	outputFileName string
	defaultName    bool
	templateNames  []string
	storage        *Storage
	noInput        bool
}

func (cs *CliState) Run() error {
	if err := cs.validateState(); err != nil {
		return err
	}

	// if no templates set, print all templates names
	if len(cs.templateNames) == 0 {
		for _, templateFile := range cs.storage.templates {
			fmt.Println(templateFile)
		}
		return nil
	}

	for _, name := range cs.templateNames {
		templateFile := cs.storage.templates[name]
		content, err := os.ReadFile(templateFile.Path)
		if err != nil {
			return err
		}

		template := NewTemplate(&templateFile, content)

		file, err := cs.getOutputFile(template)
		if err != nil {
			return err
		}

		values, err := getVariableValues(template, cs.noInput)
		if err != nil {
			return err
		}

		fmt.Fprint(file, template.fillTemplate(values))
		file.Close()
	}

	return nil
}

func (cs *CliState) validateState() error {
	if cs.defaultName && cs.outputFileName != "" {
		return fmt.Errorf("both -d and -o flags were set, but only one of them can be used at the same time")
	}

	for _, name := range cs.templateNames {
		templateFile, ok := cs.storage.templates[name]
		if !ok {
			return fmt.Errorf("template %s not found", name)
		}

		if cs.defaultName && templateFile.Filename == "" {
			return fmt.Errorf("template %s has no default name, but -d flag was set", templateFile.Name)
		}
	}

	return nil
}

func (cs *CliState) getOutputFile(template *Template) (OutputFile, error) {
	if cs.defaultName {
		return os.Create(template.Filename)
	}

	if cs.outputFileName != "" {
		return os.Create(cs.outputFileName)
	}

	return StdoutInstance(), nil
}
