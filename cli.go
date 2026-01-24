package main

import (
	"fmt"
)

// CliState represents the state of the running app with all options set
// and keep all business logic inside its functions
type CliState struct {
	outputFileName string
	defaultName    bool
	templateNames  []string
	storage        *Storage
	noInput        bool
	ioh            *IOHandler
	editMode       bool
}

func (cs *CliState) Run() error {
	if err := cs.validateState(); err != nil {
		return err
	}

	if cs.editMode {
		editor, ok := cs.ioh.LookupEnv("EDITOR")
		if !ok {
			editor = "vi"
		}

		templateFile := cs.storage.templates[cs.templateNames[0]]
		return cs.ioh.executeCommand(editor, templateFile.Path)
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
		content, err := cs.ioh.ReadFile(templateFile.Path)
		if err != nil {
			return err
		}

		template := NewTemplate(&templateFile, content)

		file, err := cs.getOutputFile(template)
		if err != nil {
			return err
		}

		values, err := cs.ioh.getVariableValues(template, cs.noInput)
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

	if cs.editMode && len(cs.templateNames) == 0 {
		return fmt.Errorf("edit mode was set, but no template name were provided")
	}

	if cs.editMode && len(cs.templateNames) != 1 {
		return fmt.Errorf("edit mode was set, but too many template names were provided")
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
		return cs.ioh.create(template.Filename)
	}

	if cs.outputFileName != "" {
		return cs.ioh.create(cs.outputFileName)
	}

	return StdoutInstance(cs.ioh), nil
}
