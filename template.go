package main

import (
	"log"
	"path/filepath"
	"strings"
)

type Template struct {
	Name      string
	Title     string
	Path      string
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
		template.Title = name
	} else {
		_, title := filepath.Split(parent)
		template.Title = title
	}

	template.Variables = make(map[string]string)

	return template
}

func (t Template) String() string {
	var buf strings.Builder
	buf.WriteString(t.Name)
	if t.Name != t.Title {
		buf.WriteString(" - ")
		buf.WriteString(t.Title)
	}
	return buf.String()
}
