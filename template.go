package main

import (
	"path/filepath"
	"strings"

	"konyahin.xyz/sttemp/parser"
)

type Template struct {
	Name      string
	Filename  string
	Content   []byte
	Variables []string
}

func NewTemplate(name string, content []byte) *Template {
	template := new(Template)

	template.Content = content

	parent, name := filepath.Split(name)
	template.Name = name

	if parent != "" {
		template.Filename = filepath.Clean(parent)
	}

	template.Variables = parser.FindVariables(content)

	return template
}

func (t Template) String() string {
	var buf strings.Builder
	buf.WriteString(t.Name)
	if t.Filename != "" {
		buf.WriteString(" - ")
		buf.WriteString(t.Filename)
	}
	return buf.String()
}
