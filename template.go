package main

import (
	"maps"
	"path/filepath"
	"slices"
	"strings"
)

type TemplateFile struct {
	Name     string
	Filename string
	Path     string
}

func NewTemplateFile(path string, baseDir string) (*TemplateFile, error) {
	relPath, err := filepath.Rel(baseDir, path)
	if err != nil {
		return nil, err
	}

	parent, name := filepath.Split(relPath)

	filename := ""
	if parent != "" {
		filename = filepath.Clean(parent)
	}

	return &TemplateFile{
		Name:     name,
		Filename: filename,
		Path:     path,
	}, nil
}

func (t TemplateFile) String() string {
	if t.Filename == "" {
		return t.Name
	}
	return t.Name + " - " + t.Filename
}

type Template struct {
	*TemplateFile
	Content   []byte
	Variables []string
	Tokens    []Token
}

func NewTemplate(templateFile *TemplateFile, content []byte) *Template {
	template := new(Template)

	template.TemplateFile = templateFile

	template.Content = content
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
	return t.TemplateFile.String()
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
