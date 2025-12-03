package main

import (
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"strings"
)

const defaultTemplateDir = ".local/share/sttemp"

type State struct {
	TemplateDir string
}

func NewState() *State {
	home, err := os.UserHomeDir()
	if err != nil {
		log.Fatal(err)
	}

	state := new(State)
	state.TemplateDir = filepath.Join(home, defaultTemplateDir)
	return state
}

func (s *State) SetTemplateDir(templatesPath string) {
	if templatesPath == "" {
		return
	}
	templatesPath, err := filepath.Abs(templatesPath)
	if err != nil {
		log.Fatal(err)
	}

	s.TemplateDir = templatesPath
}

func (s State) Templates() map[string]*Template {
	templates := make(map[string]*Template)

	err := filepath.WalkDir(s.TemplateDir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if d.IsDir() && (strings.Contains(d.Name(), "/.") || strings.HasPrefix(d.Name(), ".")) {
			return fs.SkipDir
		}

		if !d.IsDir() {
			template := NewTemplate(s, path)
			templates[template.Name] = template
		}

		return nil
	})

	if err != nil {
		log.Fatal(err)
	}

	return templates
}
