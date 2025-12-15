package main

import (
	"maps"
	"path/filepath"
	"slices"
	"strings"
)

type Template struct {
	Name      string
	Filename  string
	Content   []byte
	Variables []string
	Tokens    []Token
}

func NewTemplate(name string, content []byte) *Template {
	template := new(Template)

	template.Content = content

	parent, name := filepath.Split(name)
	template.Name = name

	if parent != "" {
		template.Filename = filepath.Clean(parent)
	}

	template.Tokens = tokens(content)

	template.Variables = findVariables(template.Tokens)

	return template
}

func (t Template) fillTemplate(values map[string]string) string {
	var sb strings.Builder
	for _, token := range t.Tokens {
		switch {
		case token.Type == Text:
			sb.Write(token.Content)
		case token.Type == Variable && len(token.Content) > 0:
			val, ok := values[string(token.Content)]
			if ok {
				sb.WriteString(val)
			}
		}
	}
	return sb.String()
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

func findVariables(tokens []Token) []string {

	vars := make(map[string]struct{})

	for _, token := range tokens {
		if token.Type == Variable && len(token.Content) > 0 {
			vars[string(token.Content)] = struct{}{}
		}
	}

	result := slices.Collect(maps.Keys(vars))
	slices.Sort(result)
	return result
}
