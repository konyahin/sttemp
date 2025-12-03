package main

import (
	"log"
	"os"
	"path/filepath"
	"strings"

	"konyahin.xyz/sttemp/parser"
)

type Template struct {
	Name      string
	Filename  string
	Path      string
	Content   string
	Variables map[string]string
}

func NewTemplate(state State, path string) *Template {
	template := new(Template)
	template.Path = path

	parent, name := filepath.Split(path)
	template.Name = name

	parent, err := filepath.Abs(parent)
	if err != nil {
		log.Fatal(err)
	}

	if parent == state.TemplateDir {
		template.Filename = name
	} else {
		_, title := filepath.Split(parent)
		template.Filename = title
	}

	content, err := os.ReadFile(template.Path)
	if err != nil {
		log.Fatal(err)
	}

	template.Content = string(content)
	template.Variables = make(map[string]string)

	variables := parser.FindVariables(content)
	for _, variable := range variables {
		template.Variables[variable] = ""
	}

	return template
}

func (t Template) String() string {
	var buf strings.Builder
	buf.WriteString(t.Name)
	if t.Name != t.Filename {
		buf.WriteString(" - ")
		buf.WriteString(t.Filename)
	}
	return buf.String()
}
